package cmd

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"net"
	"orbital/config"
	"orbital/domain"
	"orbital/pkg/certificate"
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

			secretKey, _ := cmd.Flags().GetString("sk")
			ip, _ := cmd.Flags().GetString("ip")
			dataPath, _ := cmd.Flags().GetString("datapath")

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Validating data ]"))

			if secretKey == "" || ip == "" || dataPath == "" {
				return errors.New("secret key, ip and datapath cannot be empty")
			}
			prompt.Bold(prompt.ColorGreen, "        OK")

			isReinit := false
			cfgPath := "/etc/orbital/config.yaml"
			if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
				isReinit = true
			}

			orbitalCfg := config.Config{
				SecretKey: secretKey,
				BindIP:    ip,
				Datapath:  dataPath,
			}

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Validating private key ]"))
			if err := validateEd25519SecretKey(secretKey); err != nil {
				return err
			}
			prompt.Bold(prompt.ColorGreen, " OK")

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Validating ip ]"))
			if err := validateIp(ip); err != nil {
				return err
			}
			prompt.Bold(prompt.ColorGreen, "          OK")

			// Skip this step on forced to avoid data corruption or unexpected issues
			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Create data folders ]"))
			if !isReinit {
				if err := createDataDir(orbitalCfg); err != nil {
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
				if err := os.Rename(cfgPath, backupPath); err != nil {
					// Allow this error to pass to be caught on config file save
					if !os.IsPermission(err) {
						return err
					}
				}
			}

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Creating config file ]"))
			if err := orbitalCfg.Save(cfgPath); err != nil {
				if errors.Is(err, config.ErrConfigWrite) {
					prompt.Warn(prompt.NewLine("Cannot write the file. Use sudo privileges. The config wile will created at: /etc/orbital/config.yaml"))
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
			if err := updateDataDirFromResources(deps.FS, orbitalCfg); err != nil {
				return err
			}
			prompt.Bold(prompt.ColorGreen, prompt.NewLine("OK ----"))

			var (
				caCert *x509.Certificate
				caKey  *ecdsa.PrivateKey
				err    error
			)

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Generate CA Certificate ]"))
			certsPath := filepath.Join(orbitalCfg.OrbitalRootDir(), "certs", "ca")
			if !isReinit {
				caCert, caKey, err = certificate.GenerateCA(certsPath)
				if err != nil {
					return err
				}
				prompt.Bold(prompt.ColorGreen, "OK ----")
			} else {
				prompt.Bold(prompt.ColorYellow, " Skipped")
			}

			if caCert == nil && caKey == nil {
				caCert, caKey, err = certificate.LoadCA(certsPath)
				if err != nil {
					return err
				}
			}

			serverCertsPath := filepath.Join(orbitalCfg.OrbitalRootDir(), "certs")
			if err = certificate.GenerateServerCert(caCert, caKey, serverCertsPath, ip); err != nil {
				return err
			}

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Config details ]"))
			fmt.Println()
			prompt.OK("Details |---------------------------- ")
			fmt.Println()

			prompt.Info(prompt.NewLine("--- BEGIN etc/orbital/config.yaml -------"))

			prompt.Err(prompt.NewLine("secretKey: %s"), secretKey)
			prompt.Info(prompt.NewLine("bindIp: %s"), ip)
			prompt.Info(prompt.NewLine("dataPath: %s"), dataPath)

			prompt.Info(prompt.NewLine("--- END etc/orbital/config.yaml ----------"))

			fmt.Println()
			prompt.Info(prompt.NewLine("Config file location: /etc/orbital/config.yaml"))
			if forced {
				prompt.Warn(prompt.NewLine("Old config backup:    /etc/orbital/config.yaml.old"))
			}

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Migrate database ]"))

			dbPath := filepath.Join(orbitalCfg.OrbitalRootDir(), "data")
			orbitalDB, err := db.NewDB(dbPath)
			if err != nil {
				return err
			}

			if err = db.AutoMigrate(orbitalDB, orbitalCfg.OrbitalRootDir()); err != nil {
				return err
			}
			prompt.Bold(prompt.ColorGreen, "   OK")

			prompt.Bold(prompt.ColorYellow, prompt.NewLine("[ Create first user ]"))
			sk, err := cryptographer.NewPrivateKeyFromString(orbitalCfg.SecretKey)
			if err != nil {
				return err
			}

			userRepo := domain.NewUserRepository(orbitalDB)
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
//     -- migrations
//   - certificates
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

// validateEd25519SecretKey validates if provided key is ed25519 valid
func validateEd25519SecretKey(secretKeyHex string) error {
	seedBytes, err := hex.DecodeString(secretKeyHex)
	if err != nil || len(seedBytes) != ed25519.SeedSize {
		return fmt.Errorf("%w:[secret: %s]", ErrInvalidEd25519Key, secretKeyHex)
	}

	if _, err = cryptographer.NewPrivateKeyFromSeed(seedBytes); err != nil {
		return fmt.Errorf("%w:[secret: %s]", ErrInvalidEd25519Seed, secretKeyHex)
	}

	return nil
}

func validateIp(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%w:[ip: %s]", ErrInvalidIP, ip)
	}

	return nil
}

// updateDataDirFromResources will migrate database
// TODO: use fs.WalkDir to go over embed.FS folders and files
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

	if err != nil {
		return err
	}
	return nil
}
