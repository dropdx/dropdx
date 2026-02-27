package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

type UpdateCache struct {
	LastCheck     time.Time `json:"last_check"`
	LatestVersion string    `json:"latest_version"`
}

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
	cachePath := filepath.Join(filepath.Dir(viper.ConfigFileUsed()), ".update_cache.json")
	if viper.ConfigFileUsed() == "" {
		// Fallback to default config path if config not loaded yet
		home, _ := os.UserHomeDir()
		cachePath = filepath.Join(home, ".dropdx", ".update_cache.json")
	}

	var cache UpdateCache
	if data, err := os.ReadFile(cachePath); err == nil {
		_ = json.Unmarshal(data, &cache)
	}

	// Only check once every 24 hours
	if time.Since(cache.LastCheck) < 24*time.Hour && cache.LatestVersion != "" {
		if isNewerVersion(Version, cache.LatestVersion) {
			displayUpdateMessage(cache.LatestVersion)
		}
		return
	}

	// Fetch latest version from GitHub API
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/dropdx/dropdx/releases/latest")
	if err != nil {
		return // Silently fail on network issues
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return
	}

	cache.LastCheck = time.Now()
	cache.LatestVersion = release.TagName
	if data, err := json.Marshal(cache); err == nil {
		_ = os.MkdirAll(filepath.Dir(cachePath), 0755)
		_ = os.WriteFile(cachePath, data, 0644)
	}

	if isNewerVersion(Version, cache.LatestVersion) {
		displayUpdateMessage(cache.LatestVersion)
	}
}

func isNewerVersion(current, latest string) bool {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")
	if current == "" || latest == "" {
		return false
	}
	return current != latest // Simple check, could be improved with semantic versioning parser
}

func displayUpdateMessage(latest string) {
	pterm.Warning.Printf("A new version of dropdx is available: %s (current: %s)\n", pterm.Cyan(latest), pterm.Gray(Version))
	pterm.Info.Printf("Download it from: %s\n\n", pterm.LightBlue("https://github.com/dropdx/dropdx/releases/latest"))
}
