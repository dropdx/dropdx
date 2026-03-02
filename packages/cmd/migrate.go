package cmd

import (
	"fmt"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

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

		if cfg.Version == config.CurrentVersion {
			pterm.Success.Printf("Configuration is already at the latest version: %s\n", config.CurrentVersion)
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

		// 2. Set the new version
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
	rootCmd.AddCommand(migrateCmd)
}
