package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/dropdx/dropdx/packages/config"
	"github.com/dropdx/dropdx/packages/core"
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
