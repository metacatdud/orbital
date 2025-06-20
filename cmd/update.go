package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"orbital/config"
	"orbital/domain"
	"orbital/pkg/cryptographer"
	"orbital/pkg/db"
	"orbital/pkg/files"
	"orbital/pkg/prompt"
	"os"
	"path/filepath"
	"strings"
)

func newUpdateCmd(deps Dependencies) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update orbital",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdHeader("update")

			var err error

			prompt.Bold(prompt.ColorYellow, "[ Validating config ]")
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("cannot find user config dir: %w", err)
			}
			prompt.Bold(prompt.ColorGreen, "        OK")
			fmt.Println()

			dbPath := filepath.Join(cfg.OrbitalRootDir(), "data")
			if _, err = os.Stat(dbPath); err != nil {
				return fmt.Errorf("dbPath [%s] does not exist", dbPath)
			}

			dbConn, err := db.NewDB(dbPath)
			if err != nil {
				return err
			}

			prompt.Bold(prompt.ColorYellow, "[ Validating user ]")
			secretKey, _ := cmd.Flags().GetString("sk")
			if secretKey == "" {
				return fmt.Errorf("no secret key provided")
			}

			sk, err := cryptographer.NewPrivateKeyFromString(secretKey)
			if err != nil {
				return ErrInvalidEd25519Key
			}

			// TODO: Validate with database too
			userRepo := domain.NewUserRepository(dbConn)
			user, err := userRepo.GetByPublicKey(sk.PublicKey().String())
			if err != nil {
				return err
			}

			if user.Access != "root" {
				return fmt.Errorf("wrong access level for user")
			}

			prompt.Bold(prompt.ColorGreen, "          OK")
			fmt.Println()

			// TODO: Check URL for new version (Next version)
			//prompt.Bold(prompt.ColorYellow, "[ Check for updates ]")
			//prompt.Bold(prompt.ColorGreen, "        Found:11.22.33")
			//prompt.Warn("        Not implemented yet")
			//fmt.Println()

			// TODO: Backup db: db_enc_timestamp.db
			prompt.Bold(prompt.ColorYellow, "[ Backing up database ]")

			if err = files.Backup(filepath.Join(dbPath, "orbital.db")); err != nil {
				return err
			}

			prompt.Bold(prompt.ColorGreen, "      OK")
			fmt.Println()

			// TODO: Copy migrations (and other files if ever needed: static files)
			prompt.Bold(prompt.ColorYellow, "[ Copying new files ]")
			if err = updateResources(deps.FS, cfg.OrbitalRootDir()); err != nil {
				return err
			}

			// TODO: Apply changes
			prompt.Bold(prompt.ColorGreen, prompt.NewLine("[ Migrate database ]"))
			if err = db.AutoMigrate(dbConn, cfg.OrbitalRootDir()); err != nil {
				return err
			}

			prompt.Bold(prompt.ColorGreen, "         OK")
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().String("sk", "", "Root user secret key")

	return cmd
}

// updateResources copy resources for new version
func updateResources(internalDir fs.FS, userStorage string) error {
	userStorageAbs, err := filepath.Abs(userStorage)
	if err != nil {
		return err
	}
	userStorageAbs = filepath.Clean(userStorageAbs)

	userStorageInfo, err := os.Stat(userStorageAbs)
	if err != nil {
		return err
	}

	if !userStorageInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", userStorageAbs)
	}

	return fs.WalkDir(internalDir, "resources", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}

		rel, err := filepath.Rel("resources", path)
		if err != nil {
			return err
		}

		target := filepath.Join(userStorageAbs, rel)
		target = filepath.Clean(target)
		if target != userStorageAbs && !strings.HasPrefix(target, userStorageAbs+string(filepath.Separator)) {
			return fmt.Errorf("illegal path escape: %s", rel)
		}

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		if _, err = os.Stat(target); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}

		if err = os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		inPath, err := internalDir.Open(path)
		if err != nil {
			return err
		}
		defer inPath.Close()

		outPath, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer outPath.Close()

		if _, err = io.Copy(outPath, inPath); err != nil {
			return err
		}

		return nil
	})
}
