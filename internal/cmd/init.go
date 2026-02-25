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
  # npm: 
  #   value: "npm_..."
  # pypi:
  #   value: "pypi-..."

# Provider configurations
providers:
  npm:
    template: "templates/.npmrc.tmpl"
    target: "~/.npmrc"
  pypi:
    template: "templates/.pypirc.tmpl"
    target: "~/.pypirc"
  docker:
    template: "templates/.docker-config.json.tmpl"
    target: "~/.docker/config.json"
`)
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return fmt.Errorf("failed to create default config file: %w", err)
		}
		fmt.Printf("✔ Default config file created at: %s\n", configPath)
	}

	// 4. Create default templates
	pypiTmplPath := filepath.Join(templatesDir, ".pypirc.tmpl")
	if _, err := os.Stat(pypiTmplPath); os.IsNotExist(err) {
		pypiTmpl := []byte(`[distutils]
index-servers =
    pypi

[pypi]
repository: https://upload.pypi.org/legacy/
username: __token__
password: {{.pypi}}
`)
		os.WriteFile(pypiTmplPath, pypiTmpl, 0644)
	}

	npmTmplPath := filepath.Join(templatesDir, ".npmrc.tmpl")
	if _, err := os.Stat(npmTmplPath); os.IsNotExist(err) {
		npmTmpl := []byte(`//registry.npmjs.org/:_authToken={{.npm}}
`)
		os.WriteFile(npmTmplPath, npmTmpl, 0644)
	}

	dockerTmplPath := filepath.Join(templatesDir, ".docker-config.json.tmpl")
	if _, err := os.Stat(dockerTmplPath); os.IsNotExist(err) {
		dockerTmpl := []byte(`{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "{{.docker}}"
		}
	}
}
`)
		os.WriteFile(dockerTmplPath, dockerTmpl, 0644)
	}

	fmt.Println("\nInitialization complete. You can now use 'dropdx set-token pypi' and 'dropdx apply'.")
	return nil
}
