package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/AbdElrahmaN31/pray-cli/internal/config"
	"github.com/AbdElrahmaN31/pray-cli/internal/location"
	"github.com/AbdElrahmaN31/pray-cli/internal/ui"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  `Manage the pray CLI configuration.`,
}

var configShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"list"},
	Short:   "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig()

		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		fmt.Println("Current configuration:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Print(string(data))
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}
		fmt.Println(path)
		return nil
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration in $EDITOR",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}

		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist, run 'pray init' first")
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = os.Getenv("VISUAL")
		}
		if editor == "" {
			editor = "vim"
		}

		execCmd := exec.Command(editor, path)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		return execCmd.Run()
	},
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig()

		if err := cfg.Validate(); err != nil {
			fmt.Printf("❌ Configuration is invalid: %v\n", err)
			return err
		}

		fmt.Println("✅ Configuration is valid")
		return nil
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}

		// Create default config
		defaultCfg := config.DefaultConfig()

		// Ensure directory exists
		if err := ensureConfigDir(); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Marshal and save
		data, err := yaml.Marshal(defaultCfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		fmt.Println("✅ Configuration reset to defaults")
		fmt.Printf("   Saved to: %s\n", path)
		return nil
	},
}

var saveDetected bool

var configDetectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Auto-detect location from IP",
	RunE: func(cmd *cobra.Command, args []string) error {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		// Use spinner for the detection process
		spinner := ui.NewSpinner("Detecting location from IP...")
		spinner.Start()

		detector := location.NewDetector()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		loc, err := detector.DetectFromIP(ctx)
		if err != nil {
			spinner.Fail("Failed to detect location")
			return fmt.Errorf("failed to detect location: %w", err)
		}

		spinner.Stop()
		fmt.Printf("%s Detected: %s\n", green("✓"), cyan(loc.GetDisplayAddress()))
		fmt.Printf("  Coordinates: %.4f°N, %.4f°E\n", loc.Latitude, loc.Longitude)
		fmt.Printf("  Timezone: %s\n", loc.Timezone)
		fmt.Println()

		if saveDetected {
			cfg := GetConfig()
			cfg.Location = *loc

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			path, _ := config.GetConfigPath()
			fmt.Printf("%s Location saved to: %s\n", green("✓"), path)
		} else {
			fmt.Println("Use --save flag to save this location to your config.")
		}

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  address                - City or address (e.g., "Cairo, Egypt")
  latitude               - Latitude in decimal degrees
  longitude              - Longitude in decimal degrees
  city                   - City name (e.g., "Cairo")
  country                - Country name (e.g., "Egypt")
  timezone               - Timezone (e.g., "Africa/Cairo")
  method                 - Calculation method ID (0-23)
  language               - Language: en or ar
  output.format          - Output format: table/pretty/json/slack/discord/webhook
  output.color_enabled   - Enable colored output: true/false
  output.no_emoji        - Disable emojis: true/false
  features.qibla         - Include Qibla direction: true/false
  features.dua           - Include daily Du'a: true/false
  features.hijri         - Hijri date display: title/desc/both/none
  features.hijri_holidays - Include Islamic holidays: true/false
  features.traveler_mode - Enable travel/Qasr mode: true/false
  calendar.duration      - Event duration in minutes (1-120)
  calendar.months        - Months to generate (1-12)
  calendar.alarm         - Alarm offsets, e.g., "5,10,15"
  calendar.events        - Events to include: "all" or "0,2,4"
  calendar.color         - Calendar color (hex or name)
  jumuah.enabled         - Enable Jumu'ah events: true/false
  jumuah.duration        - Jumu'ah duration in minutes
  ramadan.enabled        - Enable Ramadan mode: true/false
  ramadan.iftar_duration - Iftar event duration in minutes
  ramadan.taraweeh_duration - Taraweeh duration in minutes
  ramadan.suhoor_duration - Suhoor duration in minutes
  iqama.enabled          - Enable Iqama times: true/false
  iqama.offsets           - Iqama offsets, e.g., "15,0,10,10,5,10,0"
  cache_enabled          - Enable caching: true/false
  update_check           - Enable update checks: true/false
  api_timeout            - API timeout in seconds (5-120)`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg := GetConfig()
		green := color.New(color.FgGreen).SprintFunc()

		switch key {
		// Location
		case "address":
			cfg.Location.Address = value
			cfg.Location.Latitude = 0
			cfg.Location.Longitude = 0
			cfg.Location.City = ""
			cfg.Location.Country = ""
			cfg.Location.CountryCode = ""
			cfg.Location.Timezone = ""
			cfg.Location.Source = "manual"
		case "latitude":
			var lat float64
			if _, err := fmt.Sscanf(value, "%f", &lat); err != nil {
				return fmt.Errorf("invalid latitude: %s", value)
			}
			cfg.Location.Latitude = lat
			cfg.Location.Address = ""
			cfg.Location.Source = "manual"
		case "longitude":
			var lon float64
			if _, err := fmt.Sscanf(value, "%f", &lon); err != nil {
				return fmt.Errorf("invalid longitude: %s", value)
			}
			cfg.Location.Longitude = lon
			cfg.Location.Address = ""
			cfg.Location.Source = "manual"
		case "city":
			cfg.Location.City = value
			cfg.Location.Address = ""
			cfg.Location.Source = "manual"
		case "country":
			cfg.Location.Country = value
			cfg.Location.Address = ""
			cfg.Location.Source = "manual"
		case "timezone":
			cfg.Location.Timezone = value

		// Calculation
		case "method":
			var m int
			if _, err := fmt.Sscanf(value, "%d", &m); err != nil {
				return fmt.Errorf("invalid method: %s", value)
			}
			if m < 0 || m > 23 {
				return fmt.Errorf("method must be between 0 and 23")
			}
			cfg.Method = m
		case "language":
			if value != "en" && value != "ar" {
				return fmt.Errorf("language must be 'en' or 'ar'")
			}
			cfg.Language = value

		// Output
		case "output.format":
			valid := []string{"table", "pretty", "json", "slack", "discord", "webhook"}
			if !isOneOf(value, valid) {
				return fmt.Errorf("invalid output format: %s", value)
			}
			cfg.Output.Format = value
		case "output.color_enabled":
			cfg.Output.ColorEnabled = value == "true"
		case "output.no_emoji":
			cfg.Output.NoEmoji = value == "true"

		// Features
		case "features.qibla":
			cfg.Features.Qibla = value == "true"
		case "features.dua":
			cfg.Features.Dua = value == "true"
		case "features.hijri":
			valid := []string{"title", "desc", "both", "none"}
			if !isOneOf(value, valid) {
				return fmt.Errorf("invalid hijri option: %s (must be title/desc/both/none)", value)
			}
			cfg.Features.Hijri = value
		case "features.hijri_holidays":
			cfg.Features.HijriHolidays = value == "true"
		case "features.traveler_mode":
			cfg.Features.TravelerMode = value == "true"

		// Calendar
		case "calendar.duration":
			var d int
			if _, err := fmt.Sscanf(value, "%d", &d); err != nil || d < 1 || d > 120 {
				return fmt.Errorf("calendar.duration must be between 1 and 120")
			}
			cfg.Calendar.Duration = d
		case "calendar.months":
			var m int
			if _, err := fmt.Sscanf(value, "%d", &m); err != nil || m < 1 || m > 12 {
				return fmt.Errorf("calendar.months must be between 1 and 12")
			}
			cfg.Calendar.Months = m
		case "calendar.alarm":
			cfg.Calendar.Alarm = value
		case "calendar.events":
			cfg.Calendar.Events = value
		case "calendar.color":
			cfg.Calendar.Color = value

		// Jumu'ah
		case "jumuah.enabled":
			cfg.Jumuah.Enabled = value == "true"
		case "jumuah.duration":
			var d int
			if _, err := fmt.Sscanf(value, "%d", &d); err != nil || d < 1 {
				return fmt.Errorf("invalid jumuah.duration: %s", value)
			}
			cfg.Jumuah.Duration = d

		// Ramadan
		case "ramadan.enabled":
			cfg.Ramadan.Enabled = value == "true"
		case "ramadan.iftar_duration":
			var d int
			if _, err := fmt.Sscanf(value, "%d", &d); err != nil || d < 1 {
				return fmt.Errorf("invalid ramadan.iftar_duration: %s", value)
			}
			cfg.Ramadan.IftarDuration = d
		case "ramadan.taraweeh_duration":
			var d int
			if _, err := fmt.Sscanf(value, "%d", &d); err != nil || d < 1 {
				return fmt.Errorf("invalid ramadan.taraweeh_duration: %s", value)
			}
			cfg.Ramadan.TaraweehDuration = d
		case "ramadan.suhoor_duration":
			var d int
			if _, err := fmt.Sscanf(value, "%d", &d); err != nil || d < 1 {
				return fmt.Errorf("invalid ramadan.suhoor_duration: %s", value)
			}
			cfg.Ramadan.SuhoorDuration = d

		// Iqama
		case "iqama.enabled":
			cfg.Iqama.Enabled = value == "true"
		case "iqama.offsets":
			cfg.Iqama.Offsets = value

		// Advanced
		case "cache_enabled":
			cfg.CacheEnabled = value == "true"
		case "update_check":
			cfg.UpdateCheck = value == "true"
		case "api_timeout":
			var t int
			if _, err := fmt.Sscanf(value, "%d", &t); err != nil || t < 5 || t > 120 {
				return fmt.Errorf("api_timeout must be between 5 and 120")
			}
			cfg.APITimeout = t

		default:
			return fmt.Errorf("unknown config key: %s\nRun 'pray config set --help' to see available keys", key)
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("%s Set %s = %s\n", green("✓"), key, value)
		return nil
	},
}

// isOneOf checks if value is in the list of valid options
func isOneOf(value string, valid []string) bool {
	for _, v := range valid {
		if value == v {
			return true
		}
	}
	return false
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		cfg := GetConfig()

		var value interface{}

		switch key {
		// Location
		case "address":
			value = cfg.Location.Address
		case "latitude":
			value = cfg.Location.Latitude
		case "longitude":
			value = cfg.Location.Longitude
		case "city":
			value = cfg.Location.City
		case "country":
			value = cfg.Location.Country
		case "country_code":
			value = cfg.Location.CountryCode
		case "timezone":
			value = cfg.Location.Timezone
		case "source":
			value = cfg.Location.Source

		// Calculation
		case "method":
			value = cfg.Method
		case "language":
			value = cfg.Language

		// Output
		case "output.format":
			value = cfg.Output.Format
		case "output.color_enabled":
			value = cfg.Output.ColorEnabled
		case "output.no_emoji":
			value = cfg.Output.NoEmoji

		// Features
		case "features.qibla":
			value = cfg.Features.Qibla
		case "features.dua":
			value = cfg.Features.Dua
		case "features.hijri":
			value = cfg.Features.Hijri
		case "features.hijri_holidays":
			value = cfg.Features.HijriHolidays
		case "features.traveler_mode":
			value = cfg.Features.TravelerMode

		// Calendar
		case "calendar.duration":
			value = cfg.Calendar.Duration
		case "calendar.months":
			value = cfg.Calendar.Months
		case "calendar.alarm":
			value = cfg.Calendar.Alarm
		case "calendar.events":
			value = cfg.Calendar.Events
		case "calendar.color":
			value = cfg.Calendar.Color

		// Jumu'ah
		case "jumuah.enabled":
			value = cfg.Jumuah.Enabled
		case "jumuah.duration":
			value = cfg.Jumuah.Duration

		// Ramadan
		case "ramadan.enabled":
			value = cfg.Ramadan.Enabled
		case "ramadan.iftar_duration":
			value = cfg.Ramadan.IftarDuration
		case "ramadan.taraweeh_duration":
			value = cfg.Ramadan.TaraweehDuration
		case "ramadan.suhoor_duration":
			value = cfg.Ramadan.SuhoorDuration

		// Iqama
		case "iqama.enabled":
			value = cfg.Iqama.Enabled
		case "iqama.offsets":
			value = cfg.Iqama.Offsets

		// Advanced
		case "cache_enabled":
			value = cfg.CacheEnabled
		case "update_check":
			value = cfg.UpdateCheck
		case "api_timeout":
			value = cfg.APITimeout

		default:
			return fmt.Errorf("unknown config key: %s\nRun 'pray config set --help' to see available keys", key)
		}

		fmt.Println(value)
		return nil
	},
}

var configLocationCmd = &cobra.Command{
	Use:   "location",
	Short: "Show detailed location information",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig()
		loc := cfg.Location

		fmt.Println("📍 Location Information")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  Address:     %s\n", loc.GetDisplayAddress())
		fmt.Printf("  Latitude:    %.4f\n", loc.Latitude)
		fmt.Printf("  Longitude:   %.4f\n", loc.Longitude)
		fmt.Printf("  Timezone:    %s\n", loc.Timezone)
		fmt.Printf("  Source:      %s\n", loc.Source)
		if !loc.DetectedAt.IsZero() {
			fmt.Printf("  Detected at: %s\n", loc.DetectedAt.Format(time.RFC1123))
		}
		return nil
	},
}

var configRepairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Attempt to fix corrupted configuration",
	Long: `Attempt to repair a corrupted configuration file.

This command will:
  1. Backup the current config file
  2. Try to load and validate the config
  3. Replace invalid values with defaults
  4. Save the repaired config`,
	RunE: func(cmd *cobra.Command, args []string) error {
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		path, err := config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}

		// Check if config exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Println("No config file found. Creating default config...")
			defaultCfg := config.DefaultConfig()
			if err := defaultCfg.Save(); err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}
			fmt.Printf("%s Created default config at: %s\n", green("✓"), path)
			return nil
		}

		// Backup current config
		fmt.Println("📋 Backing up current config...")
		if err := config.Backup(); err != nil {
			fmt.Printf("%s Could not backup config: %v\n", yellow("⚠"), err)
		} else {
			fmt.Printf("%s Backup created: %s.backup\n", green("✓"), path)
		}

		// Try to load current config
		fmt.Println("🔧 Attempting to repair config...")
		currentCfg, err := config.Load()
		if err != nil {
			fmt.Printf("%s Config is corrupted, resetting to defaults\n", yellow("⚠"))
			currentCfg = config.DefaultConfig()
		}

		// Validate and fix issues
		defaultCfg := config.DefaultConfig()
		repaired := false

		// Fix method if invalid
		if !config.ValidMethodID(currentCfg.Method) {
			fmt.Printf("  Fixed: method %d → %d\n", currentCfg.Method, defaultCfg.Method)
			currentCfg.Method = defaultCfg.Method
			repaired = true
		}

		// Fix language if invalid
		if currentCfg.Language != "en" && currentCfg.Language != "ar" {
			fmt.Printf("  Fixed: language '%s' → '%s'\n", currentCfg.Language, defaultCfg.Language)
			currentCfg.Language = defaultCfg.Language
			repaired = true
		}

		// Fix output format if invalid
		validFormats := []string{"table", "pretty", "json", "slack", "discord", "webhook"}
		formatValid := false
		for _, f := range validFormats {
			if currentCfg.Output.Format == f {
				formatValid = true
				break
			}
		}
		if !formatValid {
			fmt.Printf("  Fixed: output.format '%s' → '%s'\n", currentCfg.Output.Format, defaultCfg.Output.Format)
			currentCfg.Output.Format = defaultCfg.Output.Format
			repaired = true
		}

		// Fix calendar settings
		if currentCfg.Calendar.Duration < 1 || currentCfg.Calendar.Duration > 120 {
			fmt.Printf("  Fixed: calendar.duration %d → %d\n", currentCfg.Calendar.Duration, defaultCfg.Calendar.Duration)
			currentCfg.Calendar.Duration = defaultCfg.Calendar.Duration
			repaired = true
		}

		if currentCfg.Calendar.Months < 1 || currentCfg.Calendar.Months > 12 {
			fmt.Printf("  Fixed: calendar.months %d → %d\n", currentCfg.Calendar.Months, defaultCfg.Calendar.Months)
			currentCfg.Calendar.Months = defaultCfg.Calendar.Months
			repaired = true
		}

		// Fix API timeout
		if currentCfg.APITimeout < 5 || currentCfg.APITimeout > 120 {
			fmt.Printf("  Fixed: api_timeout %d → %d\n", currentCfg.APITimeout, defaultCfg.APITimeout)
			currentCfg.APITimeout = defaultCfg.APITimeout
			repaired = true
		}

		// Save repaired config
		if err := currentCfg.Save(); err != nil {
			return fmt.Errorf("failed to save repaired config: %w", err)
		}

		if repaired {
			fmt.Printf("\n%s Configuration repaired and saved!\n", green("✓"))
		} else {
			fmt.Printf("\n%s Configuration is valid, no repairs needed.\n", green("✓"))
		}

		return nil
	},
}

var configExportFile string

var configExportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export configuration to a file",
	Long: `Export the current configuration to a YAML file.

If no file is specified, exports to ./pray-config.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		green := color.New(color.FgGreen).SprintFunc()

		cfg := GetConfig()

		// Determine output file
		outputFile := "pray-config.yaml"
		if len(args) > 0 {
			outputFile = args[0]
		}
		if configExportFile != "" {
			outputFile = configExportFile
		}

		// Export config
		if err := cfg.Export(outputFile); err != nil {
			return fmt.Errorf("failed to export config: %w", err)
		}

		fmt.Printf("%s Configuration exported to: %s\n", green("✓"), outputFile)
		return nil
	},
}

var configImportCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import configuration from a file",
	Long: `Import configuration from a YAML file.

This will replace the current configuration with the imported one.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		inputFile := args[0]

		// Check if file exists
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", inputFile)
		}

		// Backup current config
		if config.Exists() {
			fmt.Println("📋 Backing up current config...")
			if err := config.Backup(); err != nil {
				fmt.Printf("%s Could not backup: %v\n", yellow("⚠"), err)
			} else {
				fmt.Printf("%s Backup created\n", green("✓"))
			}
		}

		// Import config
		importedCfg, err := config.Import(inputFile)
		if err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}

		// Validate imported config
		if err := importedCfg.Validate(); err != nil {
			return fmt.Errorf("imported config is invalid: %w", err)
		}

		// Save to default location
		if err := importedCfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		path, _ := config.GetConfigPath()
		fmt.Printf("%s Configuration imported from: %s\n", green("✓"), inputFile)
		fmt.Printf("   Saved to: %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configDetectCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configLocationCmd)
	configCmd.AddCommand(configRepairCmd)
	configCmd.AddCommand(configExportCmd)
	configCmd.AddCommand(configImportCmd)

	// Add flags for detect command
	configDetectCmd.Flags().BoolVar(&saveDetected, "save", false, "save detected location to config")

	// Add flags for export command
	configExportCmd.Flags().StringVarP(&configExportFile, "file", "f", "", "output file path")
}
