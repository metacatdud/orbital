package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SecretKey string `yaml:"secretKey"`
	Addr      string `yaml:"addr"`
	Datapath  string `yaml:"dataPath"`
}

// Validate config.
// TODO: Better IP validation
// TODO: Better DataPath validation
func (c *Config) Validate() error {
	if len(c.SecretKey) != 64 {
		return fmt.Errorf("%w:[len: %d]", ErrSecretKeyLength, len(c.SecretKey))
	}

	if c.Addr == "" {
		return validateAddr(c.Addr)
	}

	if c.Datapath == "" {
		return fmt.Errorf("%w", ErrDataPathRequired)
	}

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

func (c *Config) OrbitalRootDir() string {
	return filepath.Join(c.Datapath, "orbital")
}

func LoadConfig() (*Config, error) {

	cfgPath := filepath.Join("/etc/orbital/config.yaml")
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

func validateAddr(addr string) error {
	if addr == "" {
		return ErrAddrIsEmpty
	}

	if ip := net.ParseIP(addr); ip != nil {
		return nil
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		if strings.Count(addr, ":") >= 2 && !strings.Contains(addr, "[") && !strings.Contains(addr, "]") {
			return fmt.Errorf("%w:[%s]", ErrAddrInvalidIPv6, addr)
		}
		return fmt.Errorf("%w:[%v]", ErrAddrInvalidIP, err)
	}

	if net.ParseIP(host) == nil {
		return fmt.Errorf("%w:[%s]", ErrAddrInvalidIP, host)
	}

	if port == "" {
		return ErrAddrPort
	}

	for _, ch := range port {
		if ch < '0' || ch > '9' {
			return ErrAddrPortNaN
		}
	}

	return nil
}
