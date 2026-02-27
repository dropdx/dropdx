package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var syncSSHConfigCmd = &cobra.Command{
	Use:   "ssh-config",
	Short: "Sync ~/.ssh/config with the dropdx vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		userHome, _ := os.UserHomeDir()
		localConfig := filepath.Join(userHome, ".ssh", "config")

		options := []string{
			"Save local ~/.ssh/config to vault",
			"Restore ~/.ssh/config from vault",
		}

		selected, _ := pterm.DefaultInteractiveSelect.
			WithDefaultText("Select sync action").
			WithOptions(options).
			Show()

		switch selected {
		case "Save local ~/.ssh/config to vault":
			if _, err := os.Stat(localConfig); os.IsNotExist(err) {
				pterm.Warning.Println("Local ~/.ssh/config not found.")
				return promptCreateNewSSHConfig()
			}
			return saveSSHConfigToVault(cfg, localConfig)

		case "Restore ~/.ssh/config from vault":
			return restoreSSHConfigFromVault(cfg, localConfig)
		}

		return nil
	},
}

func init() {
	syncCmd.AddCommand(syncSSHConfigCmd)
}

func saveSSHConfigToVault(cfg *config.Config, localPath string) error {
	machineName, err := selectMachine(cfg)
	if err != nil {
		return err
	}

	dropdxHome := getDropdxHome()
	vaultPath := filepath.Join(dropdxHome, "machines", machineName, "ssh", "config")
	_ = os.MkdirAll(filepath.Dir(vaultPath), 0700)

	if err := copyFile(localPath, vaultPath, 0600); err != nil {
		return err
	}

	pterm.Success.Printf("~/.ssh/config saved to vault for machine '%s'.\n", info(machineName))
	return nil
}

func restoreSSHConfigFromVault(cfg *config.Config, localPath string) error {
	machineName, err := selectMachine(cfg)
	if err != nil {
		return err
	}

	dropdxHome := getDropdxHome()
	vaultPath := filepath.Join(dropdxHome, "machines", machineName, "ssh", "config")
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		pterm.Warning.Printf("SSH config not found in vault for machine '%s'.\n", machineName)
		return promptCreateNewSSHConfig()
	}

	_ = os.MkdirAll(filepath.Dir(localPath), 0700)
	if err := copyFile(vaultPath, localPath, 0600); err != nil {
		return err
	}

	pterm.Success.Printf("~/.ssh/config restored from vault for machine '%s'.\n", info(machineName))
	return nil
}

func promptCreateNewSSHConfig() error {
	confirmed, _ := pterm.DefaultInteractiveConfirm.
		WithDefaultText("Do you want to create a new SSH configuration?").
		Show()
	if confirmed {
		pterm.Info.Println("Please run: dropdx create ssh-config")
	}
	return nil
}
