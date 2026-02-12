package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/anashaat/pray-cli/internal/config"
)

var filterMethods string

var methodsCmd = &cobra.Command{
	Use:   "methods",
	Short: "List all calculation methods",
	Long: `Display all available prayer time calculation methods.

Each method is used by different regions and organizations to calculate
prayer times based on specific astronomical angles.`,
	Run: func(cmd *cobra.Command, args []string) {
		methods := config.CalculationMethods

		// Apply filter if specified
		if filterMethods != "" {
			filterLower := strings.ToLower(filterMethods)
			var filtered []config.CalculationMethod
			for _, m := range methods {
				if strings.Contains(strings.ToLower(m.Name), filterLower) ||
					strings.Contains(strings.ToLower(m.Description), filterLower) {
					filtered = append(filtered, m)
				}
			}
			methods = filtered
		}

		if len(methods) == 0 {
			fmt.Println("No methods found matching the filter.")
			return
		}

		fmt.Println()
		fmt.Println("📐 Available Calculation Methods")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		// Create table with new API
		table := tablewriter.NewTable(os.Stdout)
		table.Header("ID", "Name", "Description")

		cyan := color.New(color.FgCyan).SprintFunc()

		for _, m := range methods {
			table.Append(cyan(fmt.Sprintf("%d", m.ID)), m.Name, m.Description)
		}

		table.Render()
		fmt.Println()
		fmt.Println("Use -m or --method flag to select a method:")
		fmt.Println("  pray -m 5           Use Egyptian method")
		fmt.Println("  pray --method 2     Use ISNA method")
	},
}

func init() {
	methodsCmd.Flags().StringVar(&filterMethods, "filter", "", "filter methods by name or description")
	rootCmd.AddCommand(methodsCmd)
}
