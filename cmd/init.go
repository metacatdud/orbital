package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"orbital/config"
	"orbital/pkg/cryptographer"
	"orbital/pkg/prompt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Errors
var (
	ErrInvalidIP         = errors.New("invalid ip address")
	ErrInvalidServer     = errors.New("invalid server address format")
	ErrPortOutOfRange    = errors.New("port out of range")
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
		server, _ := cmd.Flags().GetString("server")

		if secretKey == "" || ip == "" || dataPath == "" {
			return errors.New("secret key, ip and datapath cannot be empty")
		}

		if err := validateSecretKeyEd25519(secretKey); err != nil {
			return err
		}

		if err := validateIp(ip); err != nil {
			return err
		}

		clientCfg := config.Client{}
		if server != "" {
			protocol, address, port, serverPublicKey, err := parseServerAddress(server)
			if err != nil {
				return err
			}

			clientCfg = config.Client{
				IsClient:        true,
				ServerPublicKey: serverPublicKey,
				Protocol:        protocol,
				Address:         address,
				Port:            port,
			}
		}

		orbitalCfg := config.Config{
			SecretKey: secretKey,
			BindIP:    ip,
			Datapath:  dataPath,
			Client:    clientCfg,
		}

		cfgPath := "/etc/orbital/config.yaml"
		if _, err := os.Stat(cfgPath); err == nil {
			if !forced {
				prompt.Warn(prompt.NewLine("Config file already exists. Use -f, --force to overwrite"))
				fmt.Println()
				return nil
			}

			backupPath := filepath.Join(filepath.Dir(cfgPath), "config.yaml.old")
			if err := os.Rename(cfgPath, backupPath); err != nil {
				return err
			}

			prompt.Warn("Old config file backup successfully.")
			fmt.Println()
			fmt.Println()
		}

		if err := orbitalCfg.Save(cfgPath); err != nil {
			if errors.Is(err, config.ErrConfigWrite) {
				prompt.Warn(prompt.NewLine("Cannot write the file. Use sudo privileges. The config wile will created at: /etc/orbital/config.yaml"))
				prompt.Info(prompt.NewLine("If you are not conformable to run in sudo, create the file manually and copy the following contents to it"))
				fmt.Println()
				fmt.Println()

				if err := config.PrintToConsole(orbitalCfg); err != nil {
					return err
				}
			}
			return err
		}

		prompt.Bold(prompt.ColorGreen, "Details |---------------------------- ")
		prompt.Info(prompt.NewLine("- Secret key: %s"), secretKey)
		prompt.Info(prompt.NewLine("- IP: %s"), ip)
		prompt.Info(prompt.NewLine("- Data storage path: %s"), dataPath)
		fmt.Println()

		if server != "" {
			fmt.Println()
			prompt.Bold(prompt.ColorWhite, "Server details |---------------------------- ")
			prompt.Info(prompt.NewLine("- Server Id: %s"), orbitalCfg.Client.ServerPublicKey)
			prompt.Info(prompt.NewLine("- Server IP:Port: %s"), orbitalCfg.Client.Address, orbitalCfg.Client.Port)
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
	initCmd.Flags().String("server", "", "Set this node as a client. Format /protocol/ip/tcp/port/server-public-key")
	initCmd.Flags().BoolVarP(&forced, "force", "f", false, "Force overwrite of existing config file")

	rootCmd.AddCommand(initCmd)
}

// validateSecretKeyEd25519 validates if provided key is ed25519 valid
func validateSecretKeyEd25519(secretKeyHex string) error {
	seedBytes, err := hex.DecodeString(secretKeyHex)
	if err != nil || len(seedBytes) != ed25519.SeedSize {
		return fmt.Errorf("%w:[secret: %s]", ErrInvalidEd25519Key, secretKeyHex)
	}

	if _, err = cryptographer.NewPrivateKeyFromSeed(seedBytes); err != nil {
		return err
	}

	return nil
}

func validateEd25519PublicKey(pubKeyHex string) error {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil || len(pubKeyBytes) != ed25519.PublicKeySize {
		return fmt.Errorf("%w:[public: %s]", ErrInvalidEd25519Key, pubKeyHex)
	}

	return nil
}

func validateIp(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%w:[ip: %s]", ErrInvalidIP, ip)
	}

	return nil
}

func parseServerAddress(server string) (protocol, address, port, serverID string, err error) {
	parts := strings.Split(server, "/")
	if len(parts) != 6 || parts[0] != "" {
		return "", "", "", "", fmt.Errorf("%w:[format: %s]", ErrInvalidServer, "/ipv4/<ip>/tcp/<port>/server-public-key")
	}

	protocol, address, transport, port, serverID := parts[1], parts[2], parts[3], parts[4], parts[5]

	if protocol != "ipv4" && protocol != "ipv6" {
		return "", "", "", "", fmt.Errorf("%w:[protocol: %s]", ErrInvalidServer, protocol)
	}

	ip := net.ParseIP(address)
	if ip == nil {
		return "", "", "", "", fmt.Errorf("%w:[ip: %s]", ErrInvalidServer, address)
	}

	if transport != "tcp" {
		return "", "", "", "", fmt.Errorf("%w:[transport: %s]", ErrInvalidServer, transport)
	}

	var portInt int
	if portInt, err = strconv.Atoi(port); err != nil {
		if portInt < 1 || portInt > 65535 {
			return "", "", "", "", fmt.Errorf("%w:[port: %s]", ErrPortOutOfRange, port)
		}

		return "", "", "", "", fmt.Errorf("%w:[port: %s]", ErrInvalidServer, port)
	}

	if err = validateEd25519PublicKey(serverID); err != nil {
		return "", "", "", "", fmt.Errorf("%w:[server key: %s]", ErrInvalidServer, serverID)
	}

	return protocol, address, port, serverID, nil
}
