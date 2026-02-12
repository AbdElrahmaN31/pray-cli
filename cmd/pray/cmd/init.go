package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/anashaat/pray-cli/internal/config"
	"github.com/anashaat/pray-cli/internal/ui"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup wizard",
	Long: `Run the interactive setup wizard to configure the pray CLI.

This will guide you through:
  - Location setup (auto-detect or manual)
  - Calculation method selection
  - Language preference
  - Display features
  - Special features (Jumu'ah, Ramadan)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		green := color.New(color.FgGreen).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		// Run the wizard
		wizard := ui.NewWizard()
		newCfg, err := wizard.Run()
		if err != nil {
			return fmt.Errorf("setup wizard failed: %w", err)
		}

		// Save the configuration
		if err := newCfg.Save(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		// Show success message
		path, _ := config.GetConfigPath()
		fmt.Printf("Configuration saved to: %s\n", cyan(path))
		fmt.Println()
		fmt.Println("You can now run " + green("'pray'") + " to see your prayer times!")
		fmt.Println()
		fmt.Println("Commands to try:")
		fmt.Println("  pray              # Show today's prayer times")
		fmt.Println("  pray next         # Show next prayer")
		fmt.Println("  pray calendar url # Generate calendar URL")
		fmt.Println("  pray config show  # View your configuration")
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
