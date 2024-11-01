package shell

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type SetupResult struct {
	Directory string
}

func (r SetupResult) OutputFile() string {
	return filepath.Join(r.Directory, "output")
}

func SetupCommand(config *Config, command string) (SetupResult, error) {
	store := config.LogDir()
	dest := filepath.Join(store, uuid.New().String())
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return SetupResult{}, err
	}
	cmdFile := filepath.Join(dest, "cmd")
	if err := os.WriteFile(cmdFile, []byte(command), os.ModePerm); err != nil {
		return SetupResult{}, err
	}
	return SetupResult{
		Directory: dest,
	}, nil
}
