package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

/**
 * initCmd represents the init command.
 */
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the dropdx home directory and configuration",
	Long: `Creates the base directory (default: ~/.dropdx), 
the templates directory, and default config files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		confirmed, _ := pterm.DefaultInteractiveConfirm.
			WithDefaultText("Do you want to initialize dropdx configuration?").
			Show()
		
		if !confirmed {
			fmt.Println("Initialization cancelled.")
			return nil
		}
		
		return runInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

/**
 * runInit executes the initialization logic with colorful output.
 */
func runInit() error {
	home := os.Getenv("DROPDX_HOME")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("%s failed to get user home directory: %w", errCrit("✖"), err)
		}
		home = filepath.Join(userHome, ".dropdx")
	}

	// 1. Create home directory
	if err := os.MkdirAll(home, 0755); err != nil {
		return fmt.Errorf("%s failed to create home directory: %w", errCrit("✖"), err)
	}
	fmt.Printf("%s Home directory initialized: %s\n", success("✔"), info(home))

	// 2. Create templates directory
	templatesDir := filepath.Join(home, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("%s failed to create templates directory: %w", errCrit("✖"), err)
	}
	fmt.Printf("%s Templates directory created.\n", success("✔"))

	// 2.1 Create machines directory
	machinesDir := filepath.Join(home, "machines")
	if err := os.MkdirAll(machinesDir, 0755); err != nil {
		return fmt.Errorf("%s failed to create machines directory: %w", errCrit("✖"), err)
	}
	fmt.Printf("%s Machines directory created.\n", success("✔"))

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
  github:
    template: "templates/github.tmpl"
    target: "~/.bashrc"
`)
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return fmt.Errorf("%s failed to create default config file: %w", errCrit("✖"), err)
		}
		fmt.Printf("%s Default config file created at: %s\n", success("✔"), info(configPath))
	} else {
		fmt.Printf("%s config.yaml already exists, skipping.\n", warn("ℹ"))
	}

	// 4. Create default templates
	createTemplate(filepath.Join(templatesDir, "github.tmpl"), `export GITHUB_TOKEN="{{.github}}"`)
	createTemplate(filepath.Join(templatesDir, ".pypirc.tmpl"), `[distutils]
index-servers =
    pypi

[pypi]
repository: https://upload.pypi.org/legacy/
username: __token__
password: {{.pypi}}
`)

	createTemplate(filepath.Join(templatesDir, ".npmrc.tmpl"), `//registry.npmjs.org/:_authToken={{.npm}}
`)

	createTemplate(filepath.Join(templatesDir, ".docker-config.json.tmpl"), `{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "{{.docker}}"
		}
	}
}
`)

	fmt.Printf("\n%s Initialization complete. You can now use '%s' and '%s'.\n",
		success("✨"), info("dropdx set-token <provider>"), info("dropdx apply"))
	return nil
}

func createTemplate(path, content string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte(content), 0644)
		fmt.Printf("%s Created template: %s\n", success("📄"), info(filepath.Base(path)))
	}
}
