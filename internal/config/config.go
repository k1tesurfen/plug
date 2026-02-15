package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	GithubOrg string             `mapstructure:"github_org"`
	Starters  map[string]Starter `mapstructure:"starters"`
}

type Starter struct {
	Items   []string `mapstructure:"items"`
	Scripts []string `mapstructure:"scripts"`
}

// Helper to get the specific ~/.config/plug path consistently
func GetPlugConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Force ~/.config/plug on all platforms (macOS/Linux)
	return filepath.Join(home, ".config", "plug"), nil
}

// LoadConfig reads the config.toml from ~/.config/plug/
func LoadConfig() (*Config, error) {
	configPath, err := GetPlugConfigDir()
	if err != nil {
		return nil, err
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

// GetTemplateDir returns the path where raw templates are stored
func GetTemplateDir() (string, error) {
	configPath, err := GetPlugConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configPath, "templates"), nil
}
