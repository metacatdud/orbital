package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"orbital/config"
	"orbital/pkg/certificate"
	"orbital/pkg/cryptographer"
	"orbital/pkg/prompt"
	"os"
	"path/filepath"
)

// Errors
var (
	ErrInvalidIP         = errors.New("invalid ip address")
	ErrInvalidEd25519Key = errors.New("invalid ed25519 key")
)

// Flags
var (
	forced bool
)
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize orbital node",
	RunE: func(cmd *cobra.Command, args []string) error {

		cmdHeader("init")

		secretKey, _ := cmd.Flags().GetString("sk")
		ip, _ := cmd.Flags().GetString("ip")
		dataPath, _ := cmd.Flags().GetString("datapath")

		if secretKey == "" || ip == "" || dataPath == "" {
			return errors.New("secret key, ip and datapath cannot be empty")
		}

		if err := validateEd25519SecretKey(secretKey); err != nil {
			return err
		}

		if err := validateIp(ip); err != nil {
			return err
		}

		if err := createDataDir(dataPath); err != nil {
			return err
		}

		orbitalCfg := config.Config{
			SecretKey: secretKey,
			BindIP:    ip,
			Datapath:  dataPath,
		}

		cfgPath := "/etc/orbital/config.yaml"
		if _, err := os.Stat(cfgPath); err == nil {
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

		certsPath := filepath.Join(orbitalCfg.Datapath, "orbital", "certs")
		caCert, caKey, err := certificate.GenerateCA(certsPath)
		if err != nil {
			return err
		}

		if err = certificate.GenerateServerCert(caCert, caKey, dataPath, ip); err != nil {
			return err
		}

		prompt.Bold(prompt.ColorGreen, "Details |---------------------------- ")
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
		fmt.Println()
		fmt.Println()
		return nil
	},
}

func init() {
	initCmd.Flags().String("sk", "", "Secret key for node communication. Use keygen command to generate")
	initCmd.Flags().String("ip", "", "Node binding ip")
	initCmd.Flags().String("datapath", "", "Orbital data storage path")
	initCmd.Flags().BoolVarP(&forced, "force", "f", false, "Force overwrite of existing config file")

	rootCmd.AddCommand(initCmd)
}

// createDataDir create orbital folder structure
// Create Orbital data dirs
//   - orbital
//   - data
//     -- migrations
//   - certificates
func createDataDir(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {

		// Create database path
		dir := filepath.Join(p, "orbital", "data", "migrations")
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		dir = filepath.Join(p, "orbital", "certs")
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
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
		return err
	}

	return nil
}

func validateIp(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%w:[ip: %s]", ErrInvalidIP, ip)
	}

	return nil
}
