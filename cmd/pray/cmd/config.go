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

	"github.com/anashaat/pray-cli/internal/config"
	"github.com/anashaat/pray-cli/internal/location"
	"github.com/anashaat/pray-cli/internal/ui"
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
  address         - City or address (e.g., "Cairo, Egypt")
  latitude        - Latitude in decimal degrees
  longitude       - Longitude in decimal degrees
  method          - Calculation method ID (0-23)
  language        - Language: en or ar
  output.format   - Output format: table/pretty/json/slack/discord
  features.qibla  - Include Qibla direction: true/false
  features.hijri  - Hijri date display: title/desc/both/none`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg := GetConfig()
		green := color.New(color.FgGreen).SprintFunc()

		switch key {
		case "address":
			cfg.Location.Address = value
		case "latitude":
			var lat float64
			if _, err := fmt.Sscanf(value, "%f", &lat); err != nil {
				return fmt.Errorf("invalid latitude: %s", value)
			}
			cfg.Location.Latitude = lat
		case "longitude":
			var lon float64
			if _, err := fmt.Sscanf(value, "%f", &lon); err != nil {
				return fmt.Errorf("invalid longitude: %s", value)
			}
			cfg.Location.Longitude = lon
		case "method":
			var method int
			if _, err := fmt.Sscanf(value, "%d", &method); err != nil {
				return fmt.Errorf("invalid method: %s", value)
			}
			if method < 0 || method > 23 {
				return fmt.Errorf("method must be between 0 and 23")
			}
			cfg.Method = method
		case "language":
			if value != "en" && value != "ar" {
				return fmt.Errorf("language must be 'en' or 'ar'")
			}
			cfg.Language = value
		case "output.format":
			valid := []string{"table", "pretty", "json", "slack", "discord", "webhook"}
			isValid := false
			for _, v := range valid {
				if value == v {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid output format: %s", value)
			}
			cfg.Output.Format = value
		case "features.qibla":
			cfg.Features.Qibla = value == "true"
		case "features.dua":
			cfg.Features.Dua = value == "true"
		case "features.hijri":
			valid := []string{"title", "desc", "both", "none"}
			isValid := false
			for _, v := range valid {
				if value == v {
					isValid = true
					break
				}
			}
			if !isValid {
				return fmt.Errorf("invalid hijri option: %s", value)
			}
			cfg.Features.Hijri = value
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("%s Set %s = %s\n", green("✓"), key, value)
		return nil
	},
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
		case "address":
			value = cfg.Location.Address
		case "latitude":
			value = cfg.Location.Latitude
		case "longitude":
			value = cfg.Location.Longitude
		case "method":
			value = cfg.Method
		case "language":
			value = cfg.Language
		case "output.format":
			value = cfg.Output.Format
		case "features.qibla":
			value = cfg.Features.Qibla
		case "features.dua":
			value = cfg.Features.Dua
		case "features.hijri":
			value = cfg.Features.Hijri
		case "timezone":
			value = cfg.Location.Timezone
		default:
			return fmt.Errorf("unknown config key: %s", key)
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
