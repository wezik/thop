package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	configDir  string
	editor     string
	insideTmux bool
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
		configDir:  configPath,
		editor:     editor,
		insideTmux: tmuxSession != "",
	}, nil
}

func (c *Config) GetConfigDir() string { return c.configDir }
func (c *Config) GetEditor() string    { return c.editor }
func (c *Config) IsInsideTmux() bool   { return c.insideTmux }
