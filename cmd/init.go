package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"net"
	"orbital/config"
	"orbital/domain"
	"orbital/pkg/cryptographer"
	"orbital/pkg/db"
	"orbital/pkg/prompt"
	"os"
	"path/filepath"
)

// Flags
var (
	forced bool
)

func newInitCmd(deps Dependencies) *cobra.Command {

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize orbital node",
		RunE: func(cmd *cobra.Command, args []string) error {

			cmdHeader("init")

			var err error

			uCfgDir, err := os.UserConfigDir()
			if err != nil {
				return fmt.Errorf("cannot find user config dir: %w", err)
			}

			secretKey, _ := cmd.Flags().GetString("sk")
			ip, _ := cmd.Flags().GetString("ip")
			dataPath, _ := cmd.Flags().GetString("datapath")

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Validating data ]"))

			if secretKey == "" || ip == "" || dataPath == "" {
				return errors.New("secret key, ip and datapath cannot be empty")
			}
			prompt.Bold(prompt.ColorGreen, "        OK")

			isReinit := false
			cfgPath := filepath.Join(uCfgDir, "orbital/config.yaml")
			if _, err = os.Stat(cfgPath); !os.IsNotExist(err) {
				isReinit = true
			}

			orbitalCfg := config.Config{
				SecretKey: secretKey,
				BindIP:    ip,
				Datapath:  dataPath,
			}

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Validating private key ]"))
			if _, err = cryptographer.NewPrivateKeyFromString(secretKey); err != nil {
				return ErrInvalidEd25519Key
			}
			prompt.Bold(prompt.ColorGreen, " OK")

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Validating ip ]"))
			if err = validateIp(ip); err != nil {
				return err
			}
			prompt.Bold(prompt.ColorGreen, "          OK")

			// Skip this step on forced to avoid data corruption or unexpected issues
			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Create data folders ]"))
			if !isReinit {
				if err = createDataDir(orbitalCfg); err != nil {
					return err
				}
				prompt.Bold(prompt.ColorGreen, "    OK")
			} else {
				prompt.Bold(prompt.ColorYellow, "    Skipped")
			}

			if isReinit {
				if !forced {
					prompt.Err(prompt.NewLine("Config file already exists. Use -f, --force to overwrite"))
					fmt.Println()
					return nil
				}

				backupPath := filepath.Join(filepath.Dir(cfgPath), "config.yaml.old")
				if err = os.Rename(cfgPath, backupPath); err != nil {
					// Allow this error to pass to be caught on config file save
					if !os.IsPermission(err) {
						return err
					}
				}
			}

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Creating config file ]"))
			if err = orbitalCfg.Save(cfgPath); err != nil {
				if errors.Is(err, config.ErrConfigWrite) {
					prompt.Warn(prompt.NewLine("Cannot write the file. Use sudo privileges. The config wile will created at: $HOME/.config/orbital/config.yaml"))
					prompt.Info(prompt.NewLine("If you are not comfortable running Orbital with sudo, create the file manually and copy the following contents between the BEGIN and END to it"))
					fmt.Println()
					fmt.Println()

					if err := config.PrintToConsole(orbitalCfg); err != nil {
						return err
					}
				}
				return err
			}
			prompt.Bold(prompt.ColorGreen, "   OK")

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Updating dependencies ]"))
			if err = updateDataDirFromResources(deps.FS, orbitalCfg); err != nil {
				return err
			}
			prompt.Bold(prompt.ColorGreen, prompt.NewLine("OK ----"))

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Config details ]"))
			fmt.Println()
			prompt.OK("Details |---------------------------- ")
			fmt.Println()

			prompt.Info(prompt.NewLine("--- BEGIN %s/orbital/config.yaml -------"), uCfgDir)

			prompt.Err(prompt.NewLine("secretKey: %s"), secretKey)
			prompt.Info(prompt.NewLine("bindIp: %s"), ip)
			prompt.Info(prompt.NewLine("dataPath: %s"), dataPath)

			prompt.Info(prompt.NewLine("--- END %s/orbital/config.yaml ----------"), uCfgDir)

			fmt.Println()
			prompt.Info(prompt.NewLine("Config file location: %s/orbital/config.yaml"), uCfgDir)
			if forced {
				prompt.Warn(prompt.NewLine("Old config backup:    %s/orbital/config.yaml.old"), uCfgDir)
			}

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Migrate database ]"))

			dbPath := filepath.Join(orbitalCfg.OrbitalRootDir(), "data")
			dbConn, err := db.NewDB(dbPath)
			if err != nil {
				return err
			}

			if err = db.AutoMigrate(dbConn, orbitalCfg.OrbitalRootDir()); err != nil {
				return err
			}
			prompt.Bold(prompt.ColorGreen, "   OK")

			// TODO: Move this to a new custom command for creating root user
			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Create root user ]"))
			sk, err := cryptographer.NewPrivateKeyFromString(orbitalCfg.SecretKey)
			if err != nil {
				return err
			}

			userRepo := domain.NewUserRepository(dbConn)
			user := domain.User{
				ID:     uuid.New().String(),
				Name:   "admin",
				PubKey: sk.PublicKey().String(),
				Access: "root",
			}

			found, err := userRepo.ExistsByPublicKey(user.PubKey)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return err
				}
			}

			if !found {
				if err = userRepo.Save(user); err != nil {
					return err
				}
				prompt.Bold(prompt.ColorGreen, "  OK")
			} else {
				prompt.Bold(prompt.ColorYellow, "  Skipped")
			}

			fmt.Println()
			return nil
		},
	}

	initCmd.Flags().String("sk", "", "Secret key for node communication. Use keygen command to generate")
	initCmd.Flags().String("ip", "", "Node binding ip")
	initCmd.Flags().String("datapath", "", "Orbital data storage path")
	initCmd.Flags().BoolVarP(&forced, "force", "f", false, "Force overwrite of existing config file")

	return initCmd
}

// createDataDir create orbital folder structure
// Create Orbital data dirs
//   - orbital
//   - data
//   - migrations
//   - certs
func createDataDir(orbitalCfg config.Config) error {
	if _, err := os.Stat(orbitalCfg.Datapath); os.IsExist(err) {
		return nil
	}

	// Create database path
	dir := filepath.Join(orbitalCfg.OrbitalRootDir(), "data", "migrations")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%w:[%s]", ErrCannotCreateDir, dir)
	}

	dir = filepath.Join(orbitalCfg.OrbitalRootDir(), "certs", "ca")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%w:[%s]", ErrCannotCreateDir, dir)
	}

	return nil
}

func validateIp(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%w:[ip: %s]", ErrInvalidIP, ip)
	}

	return nil
}

// updateDataDirFromResources wih resources files from the embed
func updateDataDirFromResources(resDir fs.FS, orbitalCfg config.Config) error {

	err := fs.WalkDir(resDir, "resources", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel("resources", path)
		if err != nil {
			return fmt.Errorf("%w:[%s]", ErrInvalidFilepath, err.Error())
		}

		//Skip dirs?
		if d.IsDir() {
			return nil
		}

		destPath := filepath.Join(orbitalCfg.OrbitalRootDir(), relativePath)

		prompt.Info(prompt.NewLine("- Copying: %s -> %s"), relativePath, destPath)

		resFile, err := resDir.Open(path)
		if err != nil {
			return fmt.Errorf("%w:[%s]", ErrReadFile, err.Error())
		}

		defer resFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("%w:[%s]", ErrCreateFile, err.Error())
		}

		defer destFile.Close()

		if _, err = io.Copy(destFile, resFile); err != nil {
			return fmt.Errorf("%w:[%s]", ErrWriteFile, err.Error())
		}

		return nil
	})

	return err
}
