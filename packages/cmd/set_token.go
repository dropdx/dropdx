package cmd

import (
	"fmt"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	tokenName string
	tokenExp  string
)

/**
 * setTokenCmd represents the set-token command.
 */
var setTokenCmd = &cobra.Command{
	Use:   "set-token [provider]",
	Short: "Store a new token with metadata",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("%s failed to load config: %w", errCrit("✖"), err)
		}

		var provider string
		if len(args) > 0 {
			provider = args[0]
		} else {
			// Interactive provider selection
			var providers []string
			// Add configured providers
			for k := range cfg.Providers {
				providers = append(providers, k)
			}
			// Add common defaults if not present
			commonDefaults := []string{"npm", "github", "gh", "gitlab", "pypi"}
			for _, d := range commonDefaults {
				found := false
				for _, p := range providers {
					if p == d {
						found = true
						break
					}
				}
				if !found {
					providers = append(providers, d)
				}
			}
			sort.Strings(providers)
			providers = append(providers, "Other (Type custom provider)")

			selected, _ := pterm.DefaultInteractiveSelect.
				WithDefaultText("Select a provider").
				WithOptions(providers).
				Show()

			if selected == "Other (Type custom provider)" {
				provider, _ = pterm.DefaultInteractiveTextInput.
					WithDefaultText("Enter custom provider name").
					Show()
			} else {
				provider = selected
			}
		}

		if provider == "" {
			return fmt.Errorf("%s provider cannot be empty", errCrit("✖"))
		}

		// 2. If provider is npm, ask for registry
		var registry string
		if provider == "npm" {
			registry, _ = pterm.DefaultInteractiveTextInput.
				WithDefaultText("Enter registry URL").
				WithDefaultValue("https://registry.npmjs.org/").
				Show()

			if registry == "" {
				return fmt.Errorf("%s registry cannot be empty", errCrit("✖"))
			}
		}

		// 3. Get token value securely
		var tokenValue string
		if term.IsTerminal(int(syscall.Stdin)) {
			promptText := fmt.Sprintf(" Enter token for %s", info(provider))
			if registry != "" {
				promptText = fmt.Sprintf(" Enter token for %s (%s)", info(provider), info(registry))
			}
			pterm.Print(info("?"), promptText, ": ")
			byteToken, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("\n%s failed to read token: %w", errCrit("✖"), err)
			}
			pterm.Println() // Print newline after hidden input
			tokenValue = string(byteToken)
		} else {
			// Read from non-terminal (e.g., pipe)
			var input string
			fmt.Scanln(&input)
			tokenValue = input
		}

		if tokenValue == "" {
			return fmt.Errorf("%s token value cannot be empty", errCrit("✖"))
		}

		// 4. Get token name if not provided
		if !cmd.Flags().Changed("name") {
			tokenName, _ = pterm.DefaultInteractiveTextInput.
				WithDefaultText("Enter a descriptive name for this token (optional)").
				Show()
		}

		// 5. Get expiration if not provided or default
		if !cmd.Flags().Changed("exp") {
			options := []string{"7d", "30d", "60d", "90d", "Never", "Custom"}
			selected, _ := pterm.DefaultInteractiveSelect.
				WithDefaultText("Select expiration duration").
				WithOptions(options).
				WithDefaultOption("30d").
				Show()

			if selected == "Custom" {
				tokenExp, _ = pterm.DefaultInteractiveTextInput.
					WithDefaultText("Enter custom expiration (YYYY-MM-DD)").
					Show()
			} else if selected == "Never" {
				tokenExp = "false"
			} else {
				tokenExp = selected
			}
		}

		// 6. Parse expiration
		expiryDate, err := parseExpiration(tokenExp)
		if err != nil {
			return fmt.Errorf("%s invalid expiration format: %w", errCrit("✖"), err)
		}

		// 7. Update config
		if cfg.Tokens == nil {
			cfg.Tokens = make(map[string]config.TokenInfo)
		}

		newToken := config.TokenInfo{
			Value:     tokenValue,
			Name:      tokenName,
			ExpiresAt: expiryDate,
		}

		if provider == "npm" && registry != "" {
			tokenInfo := cfg.Tokens["npm"]
			// If npm is currently a single token or empty, we ensure it's structured
			// In v2, we prefer it to be a list, but registries are a map inside TokenInfo.
			// We'll keep registries as they are but ensure the container is handled correctly.
			if tokenInfo.Registries == nil {
				tokenInfo.Registries = make(map[string]config.TokenInfo)
			}
			tokenInfo.Registries[registry] = newToken
			cfg.Tokens["npm"] = tokenInfo
		} else {
			// Universal list-based storage for all providers
			tokenInfo := cfg.Tokens[provider]

			if len(tokenInfo.Items) > 0 {
				// Already a list, append
				tokenInfo.Items = append(tokenInfo.Items, newToken)
			} else if tokenInfo.Value != "" {
				// Existing single token, convert to list and append new one
				tokenInfo.Items = []config.TokenInfo{
					{
						Value:     tokenInfo.Value,
						Name:      tokenInfo.Name,
						ExpiresAt: tokenInfo.ExpiresAt,
					},
					newToken,
				}
				tokenInfo.Value = "" // Clear single value fields
				tokenInfo.Name = ""
				tokenInfo.ExpiresAt = ""
			} else {
				// New provider or empty, start a list
				tokenInfo.Items = []config.TokenInfo{newToken}
			}
			cfg.Tokens[provider] = tokenInfo
		}

		// 8. Save config
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("%s failed to save config: %w", errCrit("✖"), err)
		}

		expLabel := expiryDate
		if expLabel == "" {
			expLabel = "never"
		}

		if registry != "" {
			pterm.Success.Printf("Token for '%s' (registry: %s) saved successfully. Expiry: %s\n", info(provider), info(registry), info(expLabel))
		} else {
			pterm.Success.Printf("Token for '%s' saved successfully. Expiry: %s\n", info(provider), info(expLabel))
		}
		return nil
	},
}

func init() {
	setTokenCmd.Flags().StringVar(&tokenName, "name", "", "Descriptive name for the token")
	setTokenCmd.Flags().StringVar(&tokenExp, "exp", "30d", "Expiration duration (7d, 30d, 60d, 90d, false) or custom date (YYYY-MM-DD)")
	rootCmd.AddCommand(setTokenCmd)
}

/**
 * parseExpiration converts the exp flag into a YYYY-MM-DD string.
 */
func parseExpiration(exp string) (string, error) {
	now := time.Now()

	switch strings.ToLower(exp) {
	case "false", "none", "never":
		return "", nil
	case "7d":
		return now.AddDate(0, 0, 7).Format("2006-01-02"), nil
	case "30d":
		return now.AddDate(0, 0, 30).Format("2006-01-02"), nil
	case "60d":
		return now.AddDate(0, 0, 60).Format("2006-01-02"), nil
	case "90", "90d":
		return now.AddDate(0, 0, 90).Format("2006-01-02"), nil
	default:
		// Attempt to parse as YYYY-MM-DD
		t, err := time.Parse("2006-01-02", exp)
		if err != nil {
			return "", fmt.Errorf("use shortcuts (7d, 30d, 60d, 90d, false) or format YYYY-MM-DD")
		}
		return t.Format("2006-01-02"), nil
	}
}
