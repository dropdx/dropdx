package cmd

import (
	"fmt"
	"strings"

	"github.com/dropdx/dropdx/packages/core/config"
	"github.com/fatih/color"
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
			return fmt.Errorf("%s failed to load configuration: %w", errCrit("✖"), err)
		}

		if cfg == nil {
			fmt.Printf("%s No configuration found. Run '%s' first.\n", warn("ℹ"), info("dropdx init"))
			return nil
		}

		header := color.New(color.FgWhite, color.Bold, color.Underline).PrintlnFunc()
		tokenNameStyle := color.New(color.FgMagenta, color.Bold).SprintFunc()
		valStyle := color.New(color.FgHiBlack).SprintFunc()
		expStyle := color.New(color.FgYellow).SprintFunc()

		// 1. List Tokens
		header("--- Tokens ---")
		if len(cfg.Tokens) == 0 {
			fmt.Println("  No tokens defined.")
		} else {
			for name, info := range cfg.Tokens {
				obfuscated := obfuscate(info.Value)
				expiryInfo := ""
				if info.ExpiresAt != "" {
					expiryInfo = expStyle(fmt.Sprintf(" [Exp: %s]", info.ExpiresAt))
				}
				fmt.Printf("  %s %s%s\n", tokenNameStyle(name+":"), valStyle(obfuscated), expiryInfo)
			}
		}

		// 2. List Providers
		fmt.Println()
		header("--- Providers ---")
		if len(cfg.Providers) == 0 {
			fmt.Println("  No providers defined.")
		} else {
			for name, p := range cfg.Providers {
				fmt.Printf("  %s %s %s %s\n", 
					tokenNameStyle(name+":"), 
					info(p.Template), 
					info("→"), 
					info(p.Target))
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
