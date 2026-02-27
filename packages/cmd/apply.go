package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/dropdx/dropdx/packages/config"
	"github.com/dropdx/dropdx/packages/core"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

/**
 * applyCmd represents the apply command.
 */
var applyCmd = &cobra.Command{
	Use:   "apply [provider]",
	Short: "Apply configurations by injecting tokens into templates",
	Long: `Applies the configuration for a specific provider or for all providers
if none is specified. It replaces tokens in templates with actual values.`,
	Example: `  dropdx apply npm
  dropdx apply (applies all providers)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if cfg == nil || (len(cfg.Providers) == 0 && len(cfg.Tokens) == 0) {
			return fmt.Errorf("no configuration found. Run 'dropdx init' first or check your config.yaml")
		}

		engine := core.NewEngine(cfg)

		// 1. Ensure 'github' is in the engine if missing from config
		if _, ok := cfg.Providers["github"]; !ok {
			if cfg.Providers == nil {
				cfg.Providers = make(map[string]config.Provider)
			}
			cfg.Providers["github"] = config.Provider{
				Template: "templates/github.tmpl",
				Target:   "~/.bashrc",
			}
			// Re-create engine to include the dynamic github provider
			engine = core.NewEngine(cfg)
		}

		// 2. Ensure github template exists in the background
		home := os.Getenv("DROPDX_HOME")
		if home == "" {
			uh, _ := os.UserHomeDir()
			home = filepath.Join(uh, ".dropdx")
		}
		tmplPath := filepath.Join(home, "templates", "github.tmpl")
		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			_ = os.MkdirAll(filepath.Dir(tmplPath), 0755)
			_ = os.WriteFile(tmplPath, []byte(`export GITHUB_TOKEN="{{.github}}"`), 0644)
		}

		var providerName string
		if len(args) > 0 {
			providerName = args[0]
		} else {
			// Interactive provider selection
			var options []string
			options = append(options, "All Providers")
			
			// Use a map to avoid duplicates
			seen := make(map[string]bool)
			for k := range cfg.Providers {
				options = append(options, k)
				seen[k] = true
			}
			for k := range cfg.Tokens {
				if !seen[k] {
					options = append(options, k)
					seen[k] = true
				}
			}

			// Force 'github' into options if not seen
			if !seen["github"] {
				options = append(options, "github")
				seen["github"] = true
			}
			
			if len(options) > 1 {
				sort.Strings(options[1:]) // Sort everything except "All Providers"
				selected, _ := pterm.DefaultInteractiveSelect.
					WithDefaultText("Select a provider to apply (v2)").
					WithOptions(options).
					Show()
				
				if selected == "All Providers" {
					providerName = ""
				} else {
					providerName = selected
				}
			}
		}

		if providerName != "" {
			// Check if we have a token for this provider (specifically for github if missing)
			token, hasToken := cfg.Tokens[providerName]
			if (!hasToken || token.Value == "") && providerName == "github" {
				pterm.Warning.Printf("GitHub token is missing. Please enter it now to apply.\n")
				// Call the set-token logic for github (inline for now)
				pterm.Print(info("?"), " Enter token for ", info("github"), ": ")
				byteToken, _ := term.ReadPassword(int(syscall.Stdin))
				pterm.Println()
				tokenValue := string(byteToken)
				
				if tokenValue != "" {
					if cfg.Tokens == nil {
						cfg.Tokens = make(map[string]config.TokenInfo)
					}
					cfg.Tokens["github"] = config.TokenInfo{
						Value: tokenValue,
					}
					// Save updated config
					_ = config.Save(cfg)
					pterm.Success.Println("GitHub token saved.")
				} else {
					return fmt.Errorf("github token cannot be empty")
				}
			}

			// Apply specific provider
			return engine.ApplyProvider(providerName)
		}

		// Apply all with spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Applying all configurations..."
		s.Color("cyan")
		s.Start()
		
		err = engine.ApplyAll()
		s.Stop()
		
		if err == nil {
			fmt.Printf("\n%s All configurations applied successfully.\n", success("✨"))
		}
		
		return err
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
