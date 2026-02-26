package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

/**
 * RootCmd represents the base command when called without any subcommands.
 */
var RootCmd = &cobra.Command{
	Use:   "dropdx",
	Short: "A cross-platform CLI to sync and update PATs and configurations.",
	Long: `
      _                 _      
   __| |_ __ ___  _ __| |__  __
  / _` + "`" + ` | '__/ _ \| '_ \ / _` + "`" + ` \/ /
 | (_| | | | (_) | |_) | (_|  >  < 
  \__,_|_|  \___/| .__/ \__,_/_/\_\
                 |_|               

dropdx manages the synchronization and update of Personal Access Tokens (PAT) 
and configurations (e.g., .npmrc, environment variables) across different machines.`,
}

/**
 * Execute adds all child commands to the root command and sets flags appropriately.
 * It is called by main.main(). It only needs to happen once to the RootCmd.
 */
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $DROPDX_HOME/config.yaml or ~/.dropdx/config.yaml)")

	// Fancy up the help
	cobra.AddTemplateFunc("style", color.New(color.FgCyan, color.Bold).SprintFunc())
	RootCmd.SetHelpTemplate(`
{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}{{end}}

{{style "Usage:"}}
  {{.UseLine}}{{if .HasAvailableSubCommands}} {{style "[command]"}}{{end}}

{{if .HasAvailableSubCommands}}{{style "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding}} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableLocalFlags}}{{style "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasAvailableInheritedFlags}}{{style "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasHelpSubCommands}}{{style "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}

{{if .HasAvailableSubCommands}}Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
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
		// Just a subtle hint if config is loaded
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
