package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	ConfigDir  string
	Editor     string
	InsideTmux bool
}

func New() (*Config, error) {
	editor := os.Getenv("EDITOR")
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(userConfigDir, "thop")
	tmuxSession := os.Getenv("TMUX")

	return &Config{
		ConfigDir:  configPath,
		Editor:     editor,
		InsideTmux: tmuxSession != "",
	}, nil
}

func (c *Config) GetConfigDir() string { return c.ConfigDir }
func (c *Config) GetEditor() string    { return c.Editor }
func (c *Config) IsInsideTmux() bool   { return c.InsideTmux }
