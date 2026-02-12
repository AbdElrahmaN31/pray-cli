// Package cmd contains all CLI commands for the pray application
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/anashaat/pray-cli/internal/config"
	"github.com/anashaat/pray-cli/internal/update"
)

var (
	// Version information
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Global flags
	cfgFile      string
	verbose      bool
	quiet        bool
	noColor      bool
	outputFormat string
	outputFile   string

	// Location flags
	address    string
	latitude   float64
	longitude  float64
	autoDetect bool

	// Calculation flags
	method int

	// Display flags
	language    string
	showQibla   bool
	showDua     bool
	hijriFormat string

	// Feature flags
	travelerMode bool
	jumuahMode   bool
	ramadanMode  bool

	// Config management flags
	saveConfig   bool
	noSaveConfig bool
	noCache      bool

	// Config instance
	cfg *config.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "pray",
	Short: "🕌 Islamic prayer times CLI tool",
	Long: `pray is a command-line tool for fetching Islamic prayer times.

It supports auto-location detection, multiple calculation methods,
various output formats, and calendar integration.

Get started with:
  pray init         Interactive setup wizard
  pray              Show today's prayer times
  pray next         Show the next prayer
  pray --help       Show all available commands`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for init and version commands
		if cmd.Name() == "init" || cmd.Name() == "version" || cmd.Name() == "completion" {
			return nil
		}

		// Handle no-color flag
		if noColor {
			color.NoColor = true
		}

		return initConfig()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Check for updates if enabled in config
		if cfg != nil && cfg.UpdateCheck && !quiet {
			// Skip for certain commands
			if cmd.Name() == "version" || cmd.Name() == "completion" || cmd.Name() == "init" {
				return
			}

			// Async update check with short timeout
			checker := update.NewChecker(version).WithTimeout(3 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			select {
			case result := <-checker.CheckAsync(ctx):
				if result != nil && result.UpdateAvailable {
					fmt.Print(update.FormatUpdateMessage(result))
				}
			case <-ctx.Done():
				// Timeout, skip update notification
			}
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default behavior: show today's prayer times
		return runToday(cmd, args)
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets the version information from build flags
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pray/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output (show debug info)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "minimal output (errors only)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format: table/pretty/json/slack/discord")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "file", "f", "", "save output to file")

	// Location flags
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "city or address (e.g., \"Cairo, Egypt\")")
	rootCmd.PersistentFlags().Float64Var(&latitude, "lat", 0, "latitude in decimal degrees")
	rootCmd.PersistentFlags().Float64Var(&longitude, "lon", 0, "longitude in decimal degrees")
	rootCmd.PersistentFlags().BoolVarP(&autoDetect, "auto", "A", false, "auto-detect location from IP")

	// Calculation flags
	rootCmd.PersistentFlags().IntVarP(&method, "method", "m", 0, "calculation method ID (default: 5)")

	// Display flags
	rootCmd.PersistentFlags().StringVarP(&language, "lang", "l", "", "language: en or ar")
	rootCmd.PersistentFlags().BoolVar(&showQibla, "qibla", false, "include Qibla direction")
	rootCmd.PersistentFlags().BoolVar(&showDua, "dua", false, "include daily Du'a")
	rootCmd.PersistentFlags().StringVar(&hijriFormat, "hijri", "", "Hijri date display: title/desc/both/none")

	// Feature flags
	rootCmd.PersistentFlags().BoolVar(&travelerMode, "traveler", false, "enable travel/Qasr mode")
	rootCmd.PersistentFlags().BoolVar(&jumuahMode, "jumuah", false, "add Jumu'ah (Friday) prayer")
	rootCmd.PersistentFlags().BoolVar(&ramadanMode, "ramadan", false, "enable Ramadan mode")

	// Config management flags
	rootCmd.PersistentFlags().BoolVar(&saveConfig, "save", false, "save current flags as default config")
	rootCmd.PersistentFlags().BoolVar(&noSaveConfig, "no-save", false, "don't save to config (one-time use)")
	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "bypass cache, force fresh data")

	// Bind flags to viper
	viper.BindPFlag("output.format", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("method", rootCmd.PersistentFlags().Lookup("method"))
	viper.BindPFlag("language", rootCmd.PersistentFlags().Lookup("lang"))
}

// initConfig reads in config file and ENV variables
func initConfig() error {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config directory
		configDir, err := config.GetConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %w", err)
		}

		// Search config in config directory
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Read environment variables
	viper.SetEnvPrefix("PRAY")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, use defaults
			cfg = config.DefaultConfig()
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal config
	cfg = config.DefaultConfig()
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// runToday shows today's prayer times (default command behavior)
func runToday(cmd *cobra.Command, args []string) error {
	// Check if configured
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	// Use the shared fetchAndDisplayPrayerTimes from today.go
	return runTodayCommand(cmd, args)
}

// GetConfig returns the current configuration
func GetConfig() *config.Config {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	return cfg
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return verbose
}

// IsQuiet returns whether quiet mode is enabled
func IsQuiet() bool {
	return quiet
}

// GetLanguage returns the language flag or config value
func GetLanguage() string {
	if language != "" {
		return language
	}
	return GetConfig().Language
}

// ShouldShowQibla returns whether to show Qibla direction
func ShouldShowQibla() bool {
	return showQibla || GetConfig().Features.Qibla
}

// ShouldShowDua returns whether to show daily Du'a
func ShouldShowDua() bool {
	return showDua || GetConfig().Features.Dua
}

// GetHijriFormat returns the Hijri date format
func GetHijriFormat() string {
	if hijriFormat != "" {
		return hijriFormat
	}
	return GetConfig().Features.Hijri
}

// IsTravelerMode returns whether traveler mode is enabled
func IsTravelerMode() bool {
	return travelerMode || GetConfig().Features.TravelerMode
}

// IsJumuahMode returns whether Jumu'ah mode is enabled
func IsJumuahMode() bool {
	return jumuahMode || GetConfig().Jumuah.Enabled
}

// IsRamadanMode returns whether Ramadan mode is enabled
func IsRamadanMode() bool {
	return ramadanMode || GetConfig().Ramadan.Enabled
}

// ShouldSaveConfig returns whether to save flags to config
func ShouldSaveConfig() bool {
	return saveConfig
}

// ShouldBypassCache returns whether to bypass cache
func ShouldBypassCache() bool {
	return noCache
}

// GetOutputFile returns the output file path
func GetOutputFile() string {
	return outputFile
}

// ensureConfigDir creates the config directory if it doesn't exist
func ensureConfigDir() error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(configDir, 0755)
}
