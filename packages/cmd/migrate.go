package cmd

import (
	"fmt"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var forceMigrate bool

/**
 * migrateCmd represents the migrate command.
 */
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the configuration to the latest version",
	Long:  `Analyzes the current config.yaml and upgrades it to the latest schema version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if cfg.Version == config.CurrentVersion && !forceMigrate {
			pterm.Success.Printf("Configuration is already at the latest version: %s (use --force to re-run migration)\n", config.CurrentVersion)
			return nil
		}

		oldVersion := cfg.Version
		if oldVersion == "" {
			oldVersion = "v1"
		}

		pterm.Info.Printf("Migrating configuration from %s to %s...\n", oldVersion, config.CurrentVersion)

		// 1. Perform migrations on tokens
		if cfg.Tokens != nil {
			// Special handling: merge gh into github
			ghToken, hasGh := cfg.Tokens["gh"]
			githubToken, hasGithub := cfg.Tokens["github"]

			if hasGh {
				pterm.Info.Println("Consolidating 'gh' tokens into 'github'...")
				
				// Ensure github is a list
				if !hasGithub {
					githubToken = config.TokenInfo{}
				}
				
				// Convert gh to list if it isn't one
				ghItems := ghToken.Items
				if len(ghItems) == 0 && ghToken.Value != "" {
					ghItems = []config.TokenInfo{{
						Value:     ghToken.Value,
						Name:      "classic",
						ExpiresAt: ghToken.ExpiresAt,
					}}
				}

				// Append gh items to github
				if len(githubToken.Items) == 0 && githubToken.Value != "" {
					githubToken.Items = []config.TokenInfo{{
						Value:     githubToken.Value,
						Name:      githubToken.Name,
						ExpiresAt: githubToken.ExpiresAt,
					}}
					githubToken.Value = ""
					githubToken.Name = ""
					githubToken.ExpiresAt = ""
				}
				
				githubToken.Items = append(githubToken.Items, ghItems...)
				cfg.Tokens["github"] = githubToken
				delete(cfg.Tokens, "gh")
			}

			// Generic: convert ALL providers to lists
			for name, token := range cfg.Tokens {
				if len(token.Items) == 0 && token.Value != "" {
					token.Items = []config.TokenInfo{{
						Value:     token.Value,
						Name:      token.Name,
						ExpiresAt: token.ExpiresAt,
					}}
					token.Value = ""
					token.Name = ""
					token.ExpiresAt = ""
					cfg.Tokens[name] = token
				}
			}
		}

		// 2. Perform migrations on providers
		if cfg.Providers != nil {
			if _, hasGh := cfg.Providers["gh"]; hasGh {
				pterm.Info.Println("Consolidating 'gh' provider into 'github'...")
				delete(cfg.Providers, "gh")
				cfg.Providers["github"] = config.Provider{
					Template: "templates/github.tmpl",
					Target:   "~/.bashrc",
				}
			}
		}

		// 3. Clean up templates
		home := os.Getenv("DROPDX_HOME")
		if home == "" {
			uh, _ := os.UserHomeDir()
			home = filepath.Join(uh, ".dropdx")
		}
		ghTmpl := filepath.Join(home, "templates", "gh.tmpl")
		githubTmpl := filepath.Join(home, "templates", "github.tmpl")

		if _, err := os.Stat(ghTmpl); err == nil {
			pterm.Info.Println("Removing old gh.tmpl and updating github.tmpl...")
			_ = os.Remove(ghTmpl)
			// Always update github.tmpl to include both exports during migration
			_ = os.MkdirAll(filepath.Dir(githubTmpl), 0755)
			content := "export GITHUB_TOKEN=\"{{.github}}\"\nexport GH_TOKEN=\"{{.github}}\""
			_ = os.WriteFile(githubTmpl, []byte(content), 0644)
		}

		// 4. Set the new version
		cfg.Version = config.CurrentVersion

		// 3. Save the updated configuration
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save migrated configuration: %w", err)
		}

		pterm.Success.Printf("Configuration successfully migrated to %s!\n", config.CurrentVersion)
		return nil
	},
}

func init() {
	migrateCmd.Flags().BoolVarP(&forceMigrate, "force", "f", false, "Force migration even if already at latest version")
	rootCmd.AddCommand(migrateCmd)
}
