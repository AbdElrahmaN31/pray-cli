package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/anashaat/pray-cli/internal/api"
	"github.com/anashaat/pray-cli/internal/config"
)

var diffCmd = &cobra.Command{
	Use:   "diff <location1> <location2>",
	Short: "Compare prayer times between two locations",
	Long: `Compare prayer times between two different locations.

Displays a side-by-side comparison of prayer times showing the time
difference between the two locations.

Examples:
  pray diff "Cairo, Egypt" "London, UK"
  pray diff "New York" "Los Angeles"
  pray diff "Dubai" "Tokyo"`,
	Args: cobra.ExactArgs(2),
	RunE: runDiffCommand,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiffCommand(cmd *cobra.Command, args []string) error {
	location1 := args[0]
	location2 := args[1]

	cfg := GetConfig()
	methodID := cfg.Method
	if method != 0 {
		methodID = method
	}

	// Create API client
	client := api.NewClient(api.WithTimeout(time.Duration(cfg.APITimeout) * time.Second))

	// Fetch prayer times for both locations in parallel
	type result struct {
		resp *api.PrayerTimesResponse
		err  error
	}

	ch1 := make(chan result, 1)
	ch2 := make(chan result, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.APITimeout)*time.Second)
	defer cancel()

	// Fetch location 1
	go func() {
		params := api.NewPrayerTimesParams().
			WithDate(time.Now()).
			WithMethod(methodID).
			WithAddress(location1)
		resp, err := client.GetPrayerTimesByAddress(ctx, params)
		ch1 <- result{resp, err}
	}()

	// Fetch location 2
	go func() {
		params := api.NewPrayerTimesParams().
			WithDate(time.Now()).
			WithMethod(methodID).
			WithAddress(location2)
		resp, err := client.GetPrayerTimesByAddress(ctx, params)
		ch2 <- result{resp, err}
	}()

	// Wait for results
	r1 := <-ch1
	r2 := <-ch2

	if r1.err != nil {
		return fmt.Errorf("failed to fetch prayer times for %s: %w", location1, r1.err)
	}
	if r2.err != nil {
		return fmt.Errorf("failed to fetch prayer times for %s: %w", location2, r2.err)
	}

	resp1 := r1.resp
	resp2 := r2.resp

	// Colors
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	if noColor {
		color.NoColor = true
	}

	// Header
	fmt.Println()
	fmt.Printf("📊 %s\n", cyan("Prayer Times Comparison"))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📅 %s\n", resp1.Data.Date.Readable)
	fmt.Println()

	// Create comparison table
	table := tablewriter.NewTable(os.Stdout)
	table.Header("Prayer", location1, location2, "Difference")

	// Prayer times to compare
	prayers := []struct {
		name  string
		time1 string
		time2 string
	}{
		{"Fajr", cleanTime(resp1.Data.Timings.Fajr), cleanTime(resp2.Data.Timings.Fajr)},
		{"Sunrise", cleanTime(resp1.Data.Timings.Sunrise), cleanTime(resp2.Data.Timings.Sunrise)},
		{"Dhuhr", cleanTime(resp1.Data.Timings.Dhuhr), cleanTime(resp2.Data.Timings.Dhuhr)},
		{"Asr", cleanTime(resp1.Data.Timings.Asr), cleanTime(resp2.Data.Timings.Asr)},
		{"Maghrib", cleanTime(resp1.Data.Timings.Maghrib), cleanTime(resp2.Data.Timings.Maghrib)},
		{"Isha", cleanTime(resp1.Data.Timings.Isha), cleanTime(resp2.Data.Timings.Isha)},
		{"Midnight", cleanTime(resp1.Data.Timings.Midnight), cleanTime(resp2.Data.Timings.Midnight)},
	}

	for _, p := range prayers {
		diff := calculateTimeDiff(p.time1, p.time2)
		diffStr := formatDiff(diff)

		// Color the difference
		var coloredDiff string
		if diff > 0 {
			coloredDiff = red("+" + diffStr)
		} else if diff < 0 {
			coloredDiff = green(diffStr)
		} else {
			coloredDiff = yellow("same")
		}

		table.Append(p.name, p.time1, p.time2, coloredDiff)
	}

	table.Render()

	fmt.Println()
	fmt.Printf("⚙️  Method: %s\n", config.GetMethodName(methodID))
	fmt.Println()
	fmt.Println("Note: Positive difference means location 2 is later")
	fmt.Println()

	return nil
}

// calculateTimeDiff calculates the difference in minutes between two time strings
func calculateTimeDiff(time1, time2 string) int {
	mins1 := parseTimeToMinutes(time1)
	mins2 := parseTimeToMinutes(time2)
	return mins2 - mins1
}

// parseTimeToMinutes converts HH:MM to minutes since midnight
func parseTimeToMinutes(timeStr string) int {
	var hour, minute int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		return 0
	}
	return hour*60 + minute
}

// formatDiff formats a minute difference as a readable string
func formatDiff(mins int) string {
	if mins == 0 {
		return "0m"
	}

	abs := mins
	if abs < 0 {
		abs = -abs
	}

	if abs < 60 {
		return fmt.Sprintf("%dm", mins)
	}

	hours := abs / 60
	remaining := abs % 60

	sign := ""
	if mins < 0 {
		sign = "-"
	}

	if remaining == 0 {
		return fmt.Sprintf("%s%dh", sign, hours)
	}
	return fmt.Sprintf("%s%dh %dm", sign, hours, remaining)
}
