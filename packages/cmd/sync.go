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

/**
 * syncCmd represents the sync command.
 */
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync configurations using Git",
	Long: `Performs git pull and git push on the dropdx home directory 
to synchronize templates and tokens across different machines.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If a subcommand is called, don't run the default sync
		if cmd.HasAvailableSubCommands() && len(args) > 0 {
			return nil
		}

		confirmed, _ := pterm.DefaultInteractiveConfirm.
			WithDefaultText("Do you want to sync your configurations with git?").
			Show()
		
		if !confirmed {
			fmt.Println("Sync cancelled.")
			return nil
		}
		
		return runSync()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

/**
 * runSync executes the synchronization logic using git.
 */
func runSync() error {
	home := os.Getenv("DROPDX_HOME")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		home = filepath.Join(userHome, ".dropdx")
	}

	// 1. Check if .git exists
	gitDir := filepath.Join(home, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		fmt.Printf("%s %s is not a git repository.\n", warn("ℹ"), home)
		fmt.Println("To enable sync, initialize git in that directory:")
		fmt.Printf("  cd %s\n  git init\n  git remote add origin %s\n", home, info("<your-repo-url>"))
		return nil
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Pulling changes..."
	s.Color("cyan")
	s.Start()

	// 2. Perform git pull
	if err := executeGit(home, "pull", "--rebase"); err != nil {
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

	fmt.Printf("\n%s Sync completed successfully.\n", success("✨"))
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
