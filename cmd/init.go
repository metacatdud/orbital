package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"orbital/config"
	"orbital/domain"
	"orbital/pkg/cryptographer"
	"orbital/pkg/db"
	"orbital/pkg/prompt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var forced bool

func newInitCmd(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize orbital node",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdHeader("init")

			secretKey := strings.TrimSpace(cmd.Flag("sk").Value.String())
			addr := strings.TrimSpace(cmd.Flag("addr").Value.String())
			dataPath := strings.TrimSpace(cmd.Flag("datapath").Value.String())

			if secretKey == "" || addr == "" || dataPath == "" {
				return errors.New("secret key, addr and dataPath cannot be empty")
			}

			cfgPath := "/etc/orbital/config.yaml"
			isReinit := false
			if _, err := os.Stat(cfgPath); err == nil {
				isReinit = true
				if !forced {
					prompt.Err("\nConfig file already exists. Use -f, --force to overwrite")
					return nil
				}
			}

			orbitalCfg := config.Config{SecretKey: secretKey, Addr: addr, Datapath: dataPath}

			if _, err := cryptographer.NewPrivateKeyFromHex(secretKey); err != nil {
				return ErrInvalidEd25519Key
			}
			if err := orbitalCfg.Validate(); err != nil {
				return err
			}

			if !isReinit {
				if err := createDataDir(orbitalCfg); err != nil {
					return err
				}
			} else {
				backupPath := filepath.Join(filepath.Dir(cfgPath), "config.yaml.old")
				_ = os.Rename(cfgPath, backupPath)
			}

			if err := orbitalCfg.Save(cfgPath); err != nil {
				if errors.Is(err, config.ErrConfigWrite) {
					prompt.Warn("\nCannot write the file. Use sudo privileges. The config will be created at: /etc/orbital/config.yaml")
					prompt.Info("\nIf you prefer manual setup, copy the content below")
					_ = config.PrintToConsole(orbitalCfg)
				}
				return err
			}

			if err := updateDataDirFromResources(deps.FS, orbitalCfg); err != nil {
				return err
			}

			dbPath := filepath.Join(orbitalCfg.OrbitalRootDir(), "data")
			dbConn, err := db.NewDB(dbPath)
			if err != nil {
				return err
			}
			if err = db.AutoMigrate(cmd.Context(), dbConn); err != nil {
				return err
			}

			sk, err := cryptographer.NewPrivateKeyFromHex(orbitalCfg.SecretKey)
			if err != nil {
				return err
			}

			userRepo := domain.NewUserRepository(dbConn)
			user := domain.User{ID: uuid.New().String(), Name: "admin", PubKey: sk.PublicKey().ToHex(), Access: "root"}
			found, err := userRepo.ExistsByPublicKey(user.PubKey)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			if !found {
				if err = userRepo.Save(user); err != nil {
					return err
				}
			}

			prompt.OK("Config file location: %s", cfgPath)
			if forced && isReinit {
				prompt.Warn("Old config backup: %s", filepath.Join(filepath.Dir(cfgPath), "config.yaml.old"))
			}
			return nil
		},
	}

	cmd.Flags().String("sk", "", "Secret key for node communication. Use keygen command to generate")
	cmd.Flags().String("addr", "", "Orbital node binding address")
	cmd.Flags().String("datapath", "", "Orbital data storage path")
	cmd.Flags().BoolVarP(&forced, "force", "f", false, "Force overwrite of existing config file")
	return cmd
}

func createDataDir(orbitalCfg config.Config) error {
	if _, err := os.Stat(orbitalCfg.Datapath); os.IsExist(err) {
		return nil
	}

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

func updateDataDirFromResources(resDir fs.FS, orbitalCfg config.Config) error {
	return fs.WalkDir(resDir, "resources", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel("resources", path)
		if err != nil {
			return fmt.Errorf("%w:[%s]", ErrInvalidFilepath, err.Error())
		}
		if d.IsDir() {
			return nil
		}

		destPath := filepath.Join(orbitalCfg.OrbitalRootDir(), rel)
		prompt.Info("- Copying: %s -> %s", rel, destPath)

		in, err := resDir.Open(path)
		if err != nil {
			return fmt.Errorf("%w:[%s]", ErrReadFile, err.Error())
		}
		defer in.Close()

		if err = os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		out, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("%w:[%s]", ErrCreateFile, err.Error())
		}
		defer out.Close()

		if _, err = io.Copy(out, in); err != nil {
			return fmt.Errorf("%w:[%s]", ErrWriteFile, err.Error())
		}
		return nil
	})
}
