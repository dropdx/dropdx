package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current version of dropdx.
// This can be overridden at build time using ldflags.
var Version = "v0.1.0"

/**
 * versionCmd represents the version command.
 */
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of dropdx",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dropdx version %s\n", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
