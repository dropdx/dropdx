package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/pterm/pterm"
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

		// 1. List Tokens
		fmt.Println(header("--- Tokens ---"))
		if len(cfg.Tokens) == 0 {
			fmt.Println("  No tokens defined.")
		} else {
			for name, tInfo := range cfg.Tokens {
				obfuscated := obfuscate(tInfo.Value)
				expiryInfo := ""
				if tInfo.ExpiresAt != "" {
					expiryInfo = warn(fmt.Sprintf(" [Exp: %s]", tInfo.ExpiresAt))
				}
				
				extra := ""
				if len(tInfo.Registries) > 0 {
					extra = info(fmt.Sprintf(" (%d registries)", len(tInfo.Registries)))
				}

				fmt.Printf("  %s %s%s%s\n", tokenStyle(name+":"), muted(obfuscated), expiryInfo, extra)
			}
		}

		// 2. List Providers
		fmt.Println()
		fmt.Println(header("--- Providers ---"))
		if len(cfg.Providers) == 0 {
			fmt.Println("  No providers defined.")
		} else {
			for name, p := range cfg.Providers {
				fmt.Printf("  %s %s %s %s\n",
					tokenStyle(name+":"),
					info(p.Template),
					info("→"),
					info(p.Target))
			}
		}

		// 3. List Machines
		fmt.Println()
		fmt.Println(header("--- Machines ---"))
		if len(cfg.Machines) == 0 {
			fmt.Println("  No machines defined.")
		} else {
			for name, machine := range cfg.Machines {
				fmt.Printf("  %s %s\n",
					tokenStyle(name+":"),
					info(machine.OS))
			}
		}

		// 4. Interactive Selection
		if len(cfg.Tokens) > 0 {
			var tokenNames []string
			for k := range cfg.Tokens {
				tokenNames = append(tokenNames, k)
			}
			sort.Strings(tokenNames)
			tokenNames = append(tokenNames, "Quit")

			for {
				selected, _ := pterm.DefaultInteractiveSelect.
					WithDefaultText("Select a token to see details").
					WithOptions(tokenNames).
					Show()

				if selected == "Quit" {
					break
				}

				tInfo := cfg.Tokens[selected]
				fmt.Println()
				fmt.Printf("%s details:\n", tokenStyle(selected))
				if tInfo.Value != "" {
					fmt.Printf("  Value: %s\n", tInfo.Value)
				}
				if tInfo.Name != "" {
					fmt.Printf("  Name: %s\n", tInfo.Name)
				}
				if tInfo.ExpiresAt != "" {
					fmt.Printf("  Expires: %s\n", tInfo.ExpiresAt)
				}

				if len(tInfo.Registries) > 0 {
					fmt.Println("  Registries:")
					for reg, regInfo := range tInfo.Registries {
						fmt.Printf("    - %s:\n", info(reg))
						fmt.Printf("      Value: %s\n", regInfo.Value)
						if regInfo.Name != "" {
							fmt.Printf("      Name: %s\n", regInfo.Name)
						}
						if regInfo.ExpiresAt != "" {
							fmt.Printf("      Expires: %s\n", regInfo.ExpiresAt)
						}
					}
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
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
