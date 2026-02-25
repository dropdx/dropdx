package cmd

import (
	"fmt"
	"strings"

	"github.com/dropdx/dropdx/internal/config"
	"github.com/spf13/cobra"
)

/**
 * listCmd represents the list command.
 */
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured tokens and providers",
	Long:  `Displays all tokens (obfuscated) and providers currently defined in the configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if cfg == nil {
			fmt.Println("ℹ No configuration found. Run 'dropdx init' first.")
			return nil
		}

		// 1. List Tokens
		fmt.Println("--- Tokens ---")
		if len(cfg.Tokens) == 0 {
			fmt.Println("No tokens defined.")
		} else {
			for name, info := range cfg.Tokens {
				obfuscated := obfuscate(info.Value)
				expiryInfo := ""
				if info.ExpiresAt != "" {
					expiryInfo = fmt.Sprintf(" (Exp: %s)", info.ExpiresAt)
				}
				fmt.Printf("  %s: %s%s\n", name, obfuscated, expiryInfo)
			}
		}

		// 2. List Providers
		fmt.Println("\n--- Providers ---")
		if len(cfg.Providers) == 0 {
			fmt.Println("No providers defined.")
		} else {
			for name, p := range cfg.Providers {
				fmt.Printf("  %s:\n", name)
				fmt.Printf("    Template: %s\n", p.Template)
				fmt.Printf("    Target:   %s\n", p.Target)
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

/**
 * obfuscate hides most of the token value for security.
 */
func obfuscate(val string) string {
	if len(val) <= 8 {
		return strings.Repeat("*", len(val))
	}
	return val[:4] + "..." + val[len(val)-4:]
}
