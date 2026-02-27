package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dropdx/dropdx/packages/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var syncSSHKeysCmd = &cobra.Command{
	Use:   "ssh-keys",
	Short: "Manage SSH keys within the dropdx vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		options := []string{
			"Save local SSH key to dropdx vault",
			"Sync/Restore SSH key from dropdx vault to this machine",
		}

		selected, _ := pterm.DefaultInteractiveSelect.
			WithDefaultText("What do you want to do with SSH keys?").
			WithOptions(options).
			Show()

		switch selected {
		case "Save local SSH key to dropdx vault":
			return saveSSHKeyToVault(cfg)
		case "Sync/Restore SSH key from dropdx vault to this machine":
			return restoreSSHKeyFromVault(cfg)
		}

		return nil
	},
}

func init() {
	syncCmd.AddCommand(syncSSHKeysCmd)
}

func saveSSHKeyToVault(cfg *config.Config) error {
	// 1. Select machine
	machineName, err := selectMachine(cfg)
	if err != nil {
		return err
	}

	// 2. Find local keys in ~/.ssh
	userHome, _ := os.UserHomeDir()
	sshDir := filepath.Join(userHome, ".ssh")
	files, err := os.ReadDir(sshDir)
	if err != nil {
		return fmt.Errorf("failed to read ~/.ssh directory: %w", err)
	}

	var keyFiles []string
	for _, f := range files {
		if !f.IsDir() && !strings.HasSuffix(f.Name(), ".pub") && f.Name() != "known_hosts" && f.Name() != "config" && f.Name() != "authorized_keys" {
			keyFiles = append(keyFiles, f.Name())
		}
	}

	if len(keyFiles) == 0 {
		pterm.Warning.Println("No private SSH keys found in ~/.ssh")
		return nil
	}

	selectedKey, _ := pterm.DefaultInteractiveSelect.
		WithDefaultText("Select the SSH key to save").
		WithOptions(keyFiles).
		Show()

	// 3. Define vault paths
	dropdxHome := getDropdxHome()
	vaultSSHDir := filepath.Join(dropdxHome, "machines", machineName, "ssh")
	if err := os.MkdirAll(vaultSSHDir, 0700); err != nil {
		return fmt.Errorf("failed to create vault ssh directory: %w", err)
	}

	// Copy private key
	if err := copyFile(filepath.Join(sshDir, selectedKey), filepath.Join(vaultSSHDir, selectedKey), 0600); err != nil {
		return err
	}

	// Copy public key if exists
	pubKey := selectedKey + ".pub"
	if _, err := os.Stat(filepath.Join(sshDir, pubKey)); err == nil {
		if err := copyFile(filepath.Join(sshDir, pubKey), filepath.Join(vaultSSHDir, pubKey), 0644); err != nil {
			return err
		}
	}

	pterm.Success.Printf("SSH key '%s' saved to vault for machine '%s'.\n", info(selectedKey), info(machineName))
	return nil
}

func restoreSSHKeyFromVault(cfg *config.Config) error {
	// 1. Select machine
	machineName, err := selectMachine(cfg)
	if err != nil {
		return err
	}

	// 2. List keys in vault for this machine
	dropdxHome := getDropdxHome()
	vaultSSHDir := filepath.Join(dropdxHome, "machines", machineName, "ssh")
	files, err := os.ReadDir(vaultSSHDir)
	if err != nil || len(files) == 0 {
		pterm.Warning.Printf("No SSH keys found in vault for machine '%s'.\n", machineName)
		return nil
	}

	var keyFiles []string
	for _, f := range files {
		if !f.IsDir() && !strings.HasSuffix(f.Name(), ".pub") {
			keyFiles = append(keyFiles, f.Name())
		}
	}

	if len(keyFiles) == 0 {
		pterm.Warning.Printf("No private SSH keys found in vault for machine '%s'.\n", machineName)
		return nil
	}

	selectedKey, _ := pterm.DefaultInteractiveSelect.
		WithDefaultText("Select the SSH key to restore").
		WithOptions(keyFiles).
		Show()

	// 3. Restore to ~/.ssh
	userHome, _ := os.UserHomeDir()
	sshDir := filepath.Join(userHome, ".ssh")
	_ = os.MkdirAll(sshDir, 0700)

	// Copy private key
	destPrivate := filepath.Join(sshDir, selectedKey)
	if err := copyFile(filepath.Join(vaultSSHDir, selectedKey), destPrivate, 0600); err != nil {
		return err
	}

	// Copy public key if exists
	pubKey := selectedKey + ".pub"
	if _, err := os.Stat(filepath.Join(vaultSSHDir, pubKey)); err == nil {
		if err := copyFile(filepath.Join(vaultSSHDir, pubKey), filepath.Join(sshDir, pubKey), 0644); err != nil {
			return err
		}
	}

	pterm.Success.Printf("SSH key '%s' restored from vault.\n", info(selectedKey))

	// Advice
	fmt.Println()
	pterm.Warning.Prefix.Text = "ADVICE"
	pterm.Warning.Printf("To use this key, run:\n")
	pterm.Info.Printf("  eval \"$(ssh-agent -s)\"\n")
	pterm.Info.Printf("  ssh-add ~/.ssh/%s\n", selectedKey)

	return nil
}

func selectMachine(cfg *config.Config) (string, error) {
	if len(cfg.Machines) == 0 {
		return "", fmt.Errorf("no machines configured. Run 'dropdx create machine' first")
	}

	var machineNames []string
	for k := range cfg.Machines {
		machineNames = append(machineNames, k)
	}
	sort.Strings(machineNames)

	selected, _ := pterm.DefaultInteractiveSelect.
		WithDefaultText("Select a machine").
		WithOptions(machineNames).
		Show()

	return selected, nil
}

func getDropdxHome() string {
	home := os.Getenv("DROPDX_HOME")
	if home == "" {
		uh, _ := os.UserHomeDir()
		home = filepath.Join(uh, ".dropdx")
	}
	return home
}

func copyFile(src, dst string, perm os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
