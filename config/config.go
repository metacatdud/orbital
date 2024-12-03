package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

//type Client struct {
//	IsClient        bool   `yaml:"isClient"`
//	ServerPublicKey string `yaml:"serverPublicKey"`
//	Protocol        string `yaml:"protocol"`
//	Address         string `yaml:"server"`
//	Port            string `yaml:"port"`
//}

type Config struct {
	SecretKey string `yaml:"secretKey"`
	BindIP    string `yaml:"bindIp"`
	Datapath  string `yaml:"dataPath"`
	//Client    Client `yaml:"client"`
}

// Validate config.
// TODO: Better IP validation
// TODO: Better DataPath validation
// TODO: Better validate server connection string
func (c *Config) Validate() error {
	if len(c.SecretKey) != 64 {
		return fmt.Errorf("%w:[len: %d]", ErrSecretKeyLength, len(c.SecretKey))
	}

	if c.BindIP == "" {
		return fmt.Errorf("%w", ErrIpRequired)
	}

	if c.Datapath == "" {
		return fmt.Errorf("%w", ErrDataPathRequired)
	}

	//if c.Client.IsClient {
	//	if c.Client.ServerPublicKey == "" || c.Client.Protocol == "" || c.Client.Address == "" || c.Client.Port == "" {
	//		return fmt.Errorf("%w:[%s]", ErrConfigClient, "server address not set")
	//	}
	//}

	return nil
}

func (c *Config) Save(cfgPath string) error {
	cfgBytes, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("%w:[%s]", ErrConfigSave, err.Error())
	}

	if err = os.MkdirAll(filepath.Dir(cfgPath), 0755); err != nil {
		if os.IsPermission(err) {
			return ErrConfigWrite
		}
		return fmt.Errorf("%w:[%s]", ErrConfigSave, err.Error())
	}

	if err = os.WriteFile(cfgPath, cfgBytes, 06444); err != nil {
		if os.IsPermission(err) {
			return ErrConfigWrite
		}

		return fmt.Errorf("%w:[%s]", ErrConfigSave, err.Error())
	}

	return nil
}

func LoadConfig(cfgPath string) (*Config, error) {
	cfgBytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("%w:[%s]", ErrConfigRead, err.Error())
	}

	var cfg Config
	if err = yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		return nil, fmt.Errorf("%w:[%s]", ErrConfigRead, err.Error())
	}

	return &cfg, nil
}

func PrintToConsole(config Config) error {
	configData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	fmt.Println("---BEGIN config.yaml---")
	fmt.Println()
	fmt.Println(string(configData))
	fmt.Println("---END config.yaml---")
	fmt.Println()

	return nil
}
