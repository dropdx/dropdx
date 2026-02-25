package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
		return runSync()
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
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
		fmt.Printf("ℹ %s is not a git repository.\n", home)
		fmt.Println("To enable sync, initialize git in that directory:")
		fmt.Printf("  cd %s\n  git init\n  git remote add origin <your-repo-url>\n", home)
		return nil
	}

	// 2. Perform git pull
	fmt.Println("⬇ Pulling changes...")
	if err := executeGit(home, "pull", "--rebase"); err != nil {
		return fmt.Errorf("failed to pull: %w", err)
	}

	// 3. Perform git push
	fmt.Println("⬆ Pushing changes...")
	if err := executeGit(home, "push"); err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Println("✔ Sync completed successfully.")
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
