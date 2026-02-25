package cmd

import (
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/dropdx/dropdx/internal/config"
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
	Use:   "set-token <provider>",
	Short: "Store a new token with metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := args[0]
		
		// 1. Get token value securely
		fmt.Printf("Enter token for %s: ", provider)
		
		var tokenValue string
		if term.IsTerminal(int(syscall.Stdin)) {
			byteToken, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("\nfailed to read token: %w", err)
			}
			fmt.Println() // Print newline after hidden input
			tokenValue = string(byteToken)
		} else {
			// Read from non-terminal (e.g., pipe)
			var input string
			fmt.Scanln(&input)
			tokenValue = input
		}

		if tokenValue == "" {
			return fmt.Errorf("token value cannot be empty")
		}

		// 2. Parse expiration
		expiryDate, err := parseExpiration(tokenExp)
		if err != nil {
			return fmt.Errorf("invalid expiration format: %w", err)
		}

		// 3. Load and update config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.Tokens == nil {
			cfg.Tokens = make(map[string]config.TokenInfo)
		}

		cfg.Tokens[provider] = config.TokenInfo{
			Value:     tokenValue,
			Name:      tokenName,
			ExpiresAt: expiryDate,
		}

		// 4. Save config
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✔ Token for '%s' saved successfully. Expiry: %s\n", provider, expiryDate)
		return nil
	},
}

func init() {
	setTokenCmd.Flags().StringVar(&tokenName, "name", "", "Descriptive name for the token")
	setTokenCmd.Flags().StringVar(&tokenExp, "exp", "30d", "Expiration duration (7d, 30d, 60d, 90d, false) or custom date (YYYY-MM-DD)")
	RootCmd.AddCommand(setTokenCmd)
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
