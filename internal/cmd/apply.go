package cmd

import (
	"fmt"

	"github.com/dcdavidev/dropdx/internal/config"
	"github.com/dcdavidev/dropdx/internal/core"
	"github.com/spf13/cobra"
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

		if len(args) > 0 {
			// Apply specific provider
			providerName := args[0]
			return engine.ApplyProvider(providerName)
		}

		// Apply all
		return engine.ApplyAll()
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
}
