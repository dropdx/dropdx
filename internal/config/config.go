package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

/**
 * Config represents the main structure of the config.yaml file.
 */
type Config struct {
	Tokens    map[string]string   `mapstructure:"tokens"`
	Providers map[string]Provider `mapstructure:"providers"`
}

/**
 * Provider defines the template and target for a specific service (e.g., npm).
 */
type Provider struct {
	Template string `mapstructure:"template"`
	Target   string `mapstructure:"target"`
}

/**
 * Load returns the unmarshaled configuration from Viper.
 */
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return &cfg, nil
}

/**
 * ResolvePath converts a path to an absolute path.
 * It handles the tilde (~) shorthand for the user's home directory.
 */
func ResolvePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Expand environment variables
	path = os.ExpandEnv(path)

	// Handle ~ shorthand
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		
		if path == "~" {
			return home, nil
		}
		
		// Ensure we handle both ~/path and ~path (though ~/ is standard)
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(home, path[2:])
		} else {
			// fallback for ~path without slash
			path = filepath.Join(home, path[1:])
		}
	}

	return filepath.Abs(path)
}

/**
 * GetBaseDir returns the base directory of the current dropdx configuration.
 */
func GetBaseDir() string {
	used := viper.ConfigFileUsed()
	if used == "" {
		return ""
	}
	return filepath.Dir(used)
}
