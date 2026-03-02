package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

		// 0. Check version and warn
		if cfg.Version != config.CurrentVersion {
			pterm.Warning.Printf("Your configuration version (%s) is outdated. Current version is %s.\n", cfg.Version, config.CurrentVersion)
			pterm.Warning.Printf("Please run 'dropdx migrate' to update it.\n\n")
		}

		// 1. Ensure dynamic providers exist for tokens that don't have one
		home := os.Getenv("DROPDX_HOME")
		if home == "" {
			uh, _ := os.UserHomeDir()
			home = filepath.Join(uh, ".dropdx")
		}

		if cfg.Providers == nil {
			cfg.Providers = make(map[string]config.Provider)
		}

		// List of tokens that should have a default provider if missing
		tokensToAutoProvide := []string{"github", "gh", "gitlab", "pypi"}
		
		for _, name := range tokensToAutoProvide {
			if _, ok := cfg.Providers[name]; !ok {
				// Check if token exists (either as single value or list)
				token, hasToken := cfg.Tokens[name]
				hasContent := hasToken && (token.Value != "" || len(token.Items) > 0)
				
				// Special case: 'gh' can use 'github' tokens if it doesn't have its own
				if name == "gh" && !hasContent {
					gt, hgt := cfg.Tokens["github"]
					hasContent = hgt && (gt.Value != "" || len(gt.Items) > 0)
				}

				// For github, we always add it even if token is missing (to support the interactive prompt)
				if name == "github" || hasContent {
					envVar := strings.ToUpper(name) + "_TOKEN"
					if name == "gh" {
						envVar = "GH_TOKEN"
					}
					
					cfg.Providers[name] = config.Provider{
						Template: "templates/" + name + ".tmpl",
						Target:   "~/.bashrc",
					}
					
					tmplPath := filepath.Join(home, "templates", name + ".tmpl")
					if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
						_ = os.MkdirAll(filepath.Dir(tmplPath), 0755)
						_ = os.WriteFile(tmplPath, []byte(fmt.Sprintf(`export %s="{{.%s}}"`, envVar, name)), 0644)
					}
				}
			}
		}

		engine := core.NewEngine(cfg)

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
			hasValue := token.Value != "" || len(token.Items) > 0
			
			if (!hasToken || !hasValue) && providerName == "github" {
				pterm.Warning.Printf("GitHub token is missing. Please enter it now to apply.\n")
				pterm.Print(info("?"), " Enter token for ", info("github"), ": ")
				byteToken, _ := term.ReadPassword(int(syscall.Stdin))
				pterm.Println()
				tokenValue := string(byteToken)
				
				if tokenValue != "" {
					if cfg.Tokens == nil {
						cfg.Tokens = make(map[string]config.TokenInfo)
					}
					// Maintain list structure for v2
					cfg.Tokens["github"] = config.TokenInfo{
						Items: []config.TokenInfo{
							{Value: tokenValue, Name: "default"},
						},
					}
					// Save updated config
					_ = config.Save(cfg)
					pterm.Success.Println("GitHub token saved.")
					
					// Re-create engine with new token
					engine = core.NewEngine(cfg)
				} else {
					return fmt.Errorf("github token cannot be empty")
				}
			}

			// Apply specific provider
			err := engine.ApplyProvider(providerName)
			if err == nil {
				showSuccessAndAdvice()
			}
			return err
		}

		// Apply all with spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Applying all configurations..."
		s.Color("cyan")
		s.Start()
		
		err = engine.ApplyAll()
		s.Stop()
		
		if err == nil {
			showSuccessAndAdvice()
		}
		
		return err
	},
}

func showSuccessAndAdvice() {
	fmt.Printf("\n%s Configurations applied successfully.\n", success("✨"))

	// Suggest sourcing the RC file
	shell := os.Getenv("SHELL")
	rcFile := ".bashrc"
	if strings.Contains(shell, "zsh") {
		rcFile = ".zshrc"
	}
	fmt.Println()
	pterm.Warning.Prefix.Text = "ADVICE"
	pterm.Warning.Printf("To apply changes to your current shell session, run:\n")
	pterm.Info.Printf("  source ~/%s\n", rcFile)
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
