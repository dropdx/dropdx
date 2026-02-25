package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

/**
 * initCmd represents the init command.
 */
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the dropdx home directory and configuration",
	Long: `Creates the base directory (default: ~/.dropdx), 
the templates directory, and a default config.yaml file if they don't exist.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}

/**
 * runInit executes the initialization logic.
 */
func runInit() error {
	home := os.Getenv("DROPDX_HOME")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		home = filepath.Join(userHome, ".dropdx")
	}

	// 1. Create home directory
	if err := os.MkdirAll(home, 0755); err != nil {
		return fmt.Errorf("failed to create home directory %s: %w", home, err)
	}
	fmt.Printf("✔ Home directory initialized: %s\n", home)

	// 2. Create templates directory
	templatesDir := filepath.Join(home, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}
	fmt.Println("✔ Templates directory created.")

	// 3. Create default config.yaml
	configPath := filepath.Join(home, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := []byte(`# dropdx configuration
# Personal Access Tokens (PAT)
tokens:
  # npm: "npm_your_token_here"
  # github: "ghp_your_token_here"

# Provider configurations
providers:
  # npm:
  #   template: "templates/.npmrc.tmpl"
  #   target: "~/.npmrc"
`)
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return fmt.Errorf("failed to create default config file: %w", err)
		}
		fmt.Printf("✔ Default config file created at: %s\n", configPath)
	} else {
		fmt.Println("ℹ config.yaml already exists, skipping.")
	}

	fmt.Println("\nInitialization complete. You can now edit your config.yaml and start adding templates.")
	return nil
}
