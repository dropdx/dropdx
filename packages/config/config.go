package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

/**
 * TokenInfo holds the token value and its metadata.
 */
type TokenInfo struct {
	Value      string               `yaml:"value,omitempty" mapstructure:"value"`
	Name       string               `yaml:"name,omitempty" mapstructure:"name"`
	ExpiresAt  string               `yaml:"expires_at,omitempty" mapstructure:"expires_at"`
	Registries map[string]TokenInfo `yaml:"registries,omitempty" mapstructure:"registries"`
}

/**
 * Machine defines the properties of a synchronized machine.
 */
type Machine struct {
	Name string `yaml:"name" mapstructure:"name"`
	OS   string `yaml:"os" mapstructure:"os"`
}

/**
 * Config represents the main structure of the config.yaml file.
 */
type Config struct {
	Tokens    map[string]TokenInfo `mapstructure:"tokens"`
	Providers map[string]Provider  `mapstructure:"providers"`
	Machines  map[string]Machine   `yaml:"machines,omitempty" mapstructure:"machines"`
}

/**
 * Provider defines the template and target for a specific service (e.g., npm).
 */
type Provider struct {
	Template string `mapstructure:"template"`
	Target   string `mapstructure:"target"`
}

/**
 * Load returns the unmarshaled configuration.
 * We prefer reading the file directly with yaml.v3 to avoid Viper's
 * issues with dots in map keys (it interprets them as paths).
 */
func Load() (*Config, error) {
	var cfg Config
	configPath := viper.ConfigFileUsed()
	
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal yaml config: %w", err)
		}
		return &cfg, nil
	}

	// Fallback to viper if no file path is set
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return &cfg, nil
}

/**
 * Save writes the current configuration back to the config.yaml file.
 * We use yaml.v3 directly to have more control over the output if needed.
 */
func Save(cfg *Config) error {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		return fmt.Errorf("no config file found to save to")
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

/**
 * ResolvePath converts a path to an absolute path.
 */
func ResolvePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	path = os.ExpandEnv(path)

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}

		if path == "~" {
			return home, nil
		}

		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(home, path[2:])
		} else {
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
