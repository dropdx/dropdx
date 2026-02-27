package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources (machine, etc.)",
}

var createMachineCmd = &cobra.Command{
	Use:   "machine",
	Short: "Create a new machine configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// 1. Get Machine Name
		name, _ := pterm.DefaultInteractiveTextInput.
			WithDefaultText("Enter machine name (e.g., macbook-pro, home-pc)").
			Show()

		if name == "" {
			return fmt.Errorf("machine name cannot be empty")
		}

		// 2. Get OS
		osOptions := []string{"linux", "darwin (macOS)", "windows"}
		selectedOS, _ := pterm.DefaultInteractiveSelect.
			WithDefaultText("Select operating system").
			WithOptions(osOptions).
			Show()

		// 3. Create machine directory
		home := os.Getenv("DROPDX_HOME")
		if home == "" {
			uh, _ := os.UserHomeDir()
			home = filepath.Join(uh, ".dropdx")
		}
		machineDir := filepath.Join(home, "machines", name)
		if err := os.MkdirAll(machineDir, 0755); err != nil {
			return fmt.Errorf("failed to create machine directory: %w", err)
		}

		// 4. Update config
		if cfg.Machines == nil {
			cfg.Machines = make(map[string]config.Machine)
		}
		cfg.Machines[name] = config.Machine{
			Name: name,
			OS:   selectedOS,
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		pterm.Success.Printf("Machine '%s' (%s) created successfully!\n", info(name), info(selectedOS))
		pterm.Info.Printf("Machine directory: %s\n", info(machineDir))
		return nil
	},
}

func init() {
	createCmd.AddCommand(createMachineCmd)
	rootCmd.AddCommand(createCmd)
}
