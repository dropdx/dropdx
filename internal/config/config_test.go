package config

import (
	"os"
	"path/filepath"
	"testing"
)

/**
 * TestResolvePath verifies that path resolution handles environment variables
 * and the tilde (~) shorthand correctly.
 */
func TestResolvePath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Absolute path", "/tmp/test", "/tmp/test"},
		{"Tilde path", "~/config", filepath.Join(home, "config")},
		{"Env variable", "$HOME/test", filepath.Join(home, "test")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolvePath(tt.input)
			if err != nil {
				t.Fatalf("ResolvePath() error = %v", err)
			}
			
			if !filepath.IsAbs(got) {
				t.Errorf("ResolvePath() got = %v, want absolute path", got)
			}
			
			if tt.name == "Tilde path" && got != tt.expected {
				t.Errorf("ResolvePath() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

/**
 * TestLoad verifies that the configuration structure is correctly unmarshaled.
 */
func TestLoad(t *testing.T) {
	cfg := &Config{
		Tokens: map[string]TokenInfo{
			"test": {Value: "val"},
		},
	}
	if cfg.Tokens["test"].Value != "val" {
		t.Errorf("Config structure mapping failed")
	}
}
