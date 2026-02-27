package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync configurations",
	Long:  "Synchronize local files with the dropdx vault or perform git operations on the vault repository.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If a subcommand is called, don't do anything (handled by subcommands)
		if len(args) > 0 {
			return nil
		}

		// Interactive selection
		options := []string{
			"Sync Git Repository (Vault)",
			"Sync SSH Configuration",
			"Sync SSH Keys",
		}

		selected, _ := pterm.DefaultInteractiveSelect.
			WithDefaultText("What do you want to sync?").
			WithOptions(options).
			Show()

		switch selected {
		case "Sync Git Repository (Vault)":
			return syncRepositoryCmd.RunE(syncRepositoryCmd, nil)
		case "Sync SSH Configuration":
			return syncSSHConfigCmd.RunE(syncSSHConfigCmd, nil)
		case "Sync SSH Keys":
			return syncSSHKeysCmd.RunE(syncSSHKeysCmd, nil)
		}

		return nil
	},
}

/**
 * syncRepositoryCmd represents the repository sync subcommand.
 */
var syncRepositoryCmd = &cobra.Command{
	Use:   "repository",
	Short: "Sync the dropdx vault repository using Git",
	Long: `Performs git pull and git push on the dropdx home directory 
to synchronize templates and tokens across different machines.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		confirmed, _ := pterm.DefaultInteractiveConfirm.
			WithDefaultText("Do you want to sync your vault repository with git?").
			Show()

		if !confirmed {
			fmt.Println("Sync cancelled.")
			return nil
		}

		return runSync()
	},
}

func init() {
	syncCmd.AddCommand(syncRepositoryCmd)
	rootCmd.AddCommand(syncCmd)
}

/**
 * runSync executes the synchronization logic using git with autostash.
 */
func runSync() error {
	home := getDropdxHome()

	// 1. Check if .git exists
	gitDir := filepath.Join(home, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		fmt.Printf("%s %s is not a git repository.\n", warn("ℹ"), home)
		fmt.Println("To enable sync, initialize git in that directory:")
		fmt.Printf("  cd %s\n  git init\n  git remote add origin %s\n", home, info("<your-repo-url>"))
		return nil
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Pulling changes (with autostash)..."
	s.Color("cyan")
	s.Start()

	// 2. Perform git pull with autostash
	if err := executeGit(home, "pull", "--rebase", "--autostash"); err != nil {
		s.Stop()
		return fmt.Errorf("failed to pull: %w", err)
	}
	s.Stop()
	fmt.Printf("%s Pulled changes successfully.\n", success("✔"))

	s.Suffix = " Pushing changes..."
	s.Restart()

	// 3. Perform git push
	if err := executeGit(home, "push"); err != nil {
		s.Stop()
		return fmt.Errorf("failed to push: %w", err)
	}
	s.Stop()
	fmt.Printf("%s Pushed changes successfully.\n", success("✔"))

	fmt.Printf("\n%s Repository sync completed successfully.\n", success("✨"))
	return nil
}

/**
 * executeGit runs a git command in the specified directory.
 */
func executeGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
