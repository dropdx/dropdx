package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update existing resources (remote, etc.)",
}

var createSSHConfigCmd = &cobra.Command{
	Use:   "ssh-config",
	Short: "Create a new SSH configuration entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// 1. Collect remote info
		alias, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Server Alias (e.g., prod-db)").Show()
		if alias == "" {
			return fmt.Errorf("alias cannot be empty")
		}

		host, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Host (IP or hostname)").Show()
		portStr, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Port").WithDefaultValue("22").Show()
		port, _ := strconv.Atoi(portStr)
		user, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("User").WithDefaultValue("root").Show()

		// 2. Update config.yaml
		if cfg.Remotes == nil {
			cfg.Remotes = make(map[string]config.Remote)
		}
		cfg.Remotes[alias] = config.Remote{
			Alias: alias,
			Host:  host,
			Port:  port,
			User:  user,
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		pterm.Success.Printf("Remote '%s' saved in config.yaml.\n", info(alias))

		// 3. Ask to sync with machines
		machineName, err := selectMachine(cfg)
		if err == nil {
			if err := writeSSHConfigForMachine(cfg, machineName); err != nil {
				return err
			}
			pterm.Success.Printf("Updated vault ssh/config for machine '%s'.\n", info(machineName))
		}

		return nil
	},
}

var updateRemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Update a remote server configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if len(cfg.Remotes) == 0 {
			pterm.Warning.Println("No remotes configured. Run 'dropdx create ssh-config' first.")
			return nil
		}

		var remoteAliases []string
		for k := range cfg.Remotes {
			remoteAliases = append(remoteAliases, k)
		}
		sort.Strings(remoteAliases)

		selectedRemote, _ := pterm.DefaultInteractiveSelect.
			WithDefaultText("Select remote to update").
			WithOptions(remoteAliases).
			Show()

		remote := cfg.Remotes[selectedRemote]

		// Update fields
		host, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Host").WithDefaultValue(remote.Host).Show()
		portStr, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("Port").WithDefaultValue(strconv.Itoa(remote.Port)).Show()
		port, _ := strconv.Atoi(portStr)
		user, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("User").WithDefaultValue(remote.User).Show()

		cfg.Remotes[selectedRemote] = config.Remote{
			Alias: selectedRemote,
			Host:  host,
			Port:  port,
			User:  user,
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		pterm.Success.Printf("Remote '%s' updated.\n", info(selectedRemote))

		// Sync with machines?
		confirmed, _ := pterm.DefaultInteractiveConfirm.
			WithDefaultText("Do you want to sync this update to a machine in the vault?").
			Show()

		if confirmed {
			machineName, err := selectMachine(cfg)
			if err == nil {
				if err := writeSSHConfigForMachine(cfg, machineName); err != nil {
					return err
				}
				pterm.Success.Printf("Updated vault config for machine '%s'.\n", info(machineName))

				syncLocal, _ := pterm.DefaultInteractiveConfirm.
					WithDefaultText("Do you want to sync with your local ~/.ssh/config?").
					Show()
				if syncLocal {
					home, _ := os.UserHomeDir()
					localPath := filepath.Join(home, ".ssh", "config")
					err = writeSSHConfigToFile(cfg, localPath)
					if err != nil {
						return err
					}
					pterm.Success.Println("Local ~/.ssh/config updated.")
				}
			}
		}

		return nil
	},
}

func init() {
	createCmd.AddCommand(createSSHConfigCmd)
	updateCmd.AddCommand(updateRemoteCmd)
	rootCmd.AddCommand(updateCmd)
}

func writeSSHConfigForMachine(cfg *config.Config, machineName string) error {
	dropdxHome := getDropdxHome()
	vaultPath := filepath.Join(dropdxHome, "machines", machineName, "ssh", "config")
	_ = os.MkdirAll(filepath.Dir(vaultPath), 0700)

	return writeSSHConfigToFile(cfg, vaultPath)
}

func writeSSHConfigToFile(cfg *config.Config, filePath string) error {
	var content string
	content = "# dropdx generated ssh-config\n\n"

	var aliases []string
	for k := range cfg.Remotes {
		aliases = append(aliases, k)
	}
	sort.Strings(aliases)

	for _, alias := range aliases {
		r := cfg.Remotes[alias]
		content += fmt.Sprintf("Host %s\n", r.Alias)
		content += fmt.Sprintf("    HostName %s\n", r.Host)
		content += fmt.Sprintf("    Port %d\n", r.Port)
		content += fmt.Sprintf("    User %s\n\n", r.User)
	}

	return os.WriteFile(filePath, []byte(content), 0600)
}
