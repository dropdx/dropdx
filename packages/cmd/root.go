package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "dropdx",
	Short: "A cross-platform CLI to sync and update PATs and configurations.",
	Long: `dropdx manages the synchronization and update of Personal Access Tokens (PAT)
and configurations (e.g., .npmrc, environment variables) across different machines.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration
		initConfig()

		if os.Getenv("DROPDX_TEST") == "true" {
			return
		}

		// Don't show banner for version commands or completion
		if cmd.Name() == "version" || cmd.Name() == "completion" {
			return
		}

		CheckForUpdates()

		// Show banner
		pterm.DefaultBigText.WithLetters(
			pterm.NewLettersFromStringWithStyle("dropdx", pterm.NewStyle(pterm.FgYellow)),
		).Render()

		pterm.Info.Println("The secure fortress for your development tokens.")
		fmt.Println()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

/**
 * Execute adds all child commands to the root command and sets flags appropriately.
 * It is called by main.main(). It only needs to happen once to the rootCmd.
 */
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Printf("✖ Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $DROPDX_HOME/config.yaml or ~/.dropdx/config.yaml)")

	// Fancy up the help
	cobra.AddTemplateFunc("accent", func(s string) string {
		return accentStyle.Sprint(s)
	})

	rootCmd.SetHelpTemplate(`
{{.Short}}

{{accent "Usage:"}}
  {{.UseLine}}{{if .HasAvailableSubCommands}} {{accent "[command]"}}{{end}}

{{if .HasAvailableSubCommands}}{{accent "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding}} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableLocalFlags}}{{accent "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasAvailableInheritedFlags}}{{accent "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasHelpSubCommands}}{{accent "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableSubCommands}}Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
	rootCmd.SetVersionTemplate("dropdx {{.Version}}\n")
}

/**
 * initConfig reads in config file and ENV variables if set.
 */
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home := os.Getenv("DROPDX_HOME")
		if home == "" {
			userHome, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting user home dir:", err)
				os.Exit(1)
			}
			home = filepath.Join(userHome, ".dropdx")
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// Config loaded
	}
}

/**
 * CheckForUpdates checks if a newer version of dropdx is available.
 */
func CheckForUpdates() {
	// Placeholder for update check logic
}
