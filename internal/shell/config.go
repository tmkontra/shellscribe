package shell

import (
	"os"
	"path/filepath"
)

type Config struct {
	homeDir string
}

func NewConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Config{}, err
	}
	return &Config{
		homeDir: homeDir,
	}, nil
}

func (c *Config) LogDir() string {
	return filepath.Join(c.homeDir, ".local", "share", "shellscribe")
}
