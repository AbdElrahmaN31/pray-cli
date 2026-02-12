package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var dateFlag string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get prayer times for a specific date",
	Long: `Fetch prayer times for a specific date.

Supported date formats:
  - ISO format: "2026-02-10"
  - Relative: "today", "tomorrow", "yesterday"
  - Weekday: "monday", "tuesday", "friday", etc.
  - Offset: "+1", "+7", "-3" (days from today)

Examples:
  pray get --date tomorrow
  pray get --date 2026-02-15
  pray get --date friday
  pray get --date +7`,
	RunE: runGetCommand,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&dateFlag, "date", "today", "date to fetch prayer times for")
}

func runGetCommand(cmd *cobra.Command, args []string) error {
	// Parse the date
	targetDate, err := parseDate(dateFlag)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	return fetchAndDisplayPrayerTimes(cmd, targetDate)
}

// parseDate parses various date formats into time.Time
func parseDate(dateStr string) (time.Time, error) {
	now := time.Now()
	dateStr = strings.ToLower(strings.TrimSpace(dateStr))

	// Handle relative dates
	switch dateStr {
	case "today":
		return now, nil
	case "tomorrow":
		return now.AddDate(0, 0, 1), nil
	case "yesterday":
		return now.AddDate(0, 0, -1), nil
	}

	// Handle offset format (+1, -3, etc.)
	if len(dateStr) > 0 && (dateStr[0] == '+' || dateStr[0] == '-') {
		var days int
		_, err := fmt.Sscanf(dateStr, "%d", &days)
		if err == nil {
			return now.AddDate(0, 0, days), nil
		}
	}

	// Handle weekday names
	weekdays := map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}

	if targetWeekday, ok := weekdays[dateStr]; ok {
		return getNextWeekday(now, targetWeekday), nil
	}

	// Handle ISO format (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}

	// Handle DD-MM-YYYY format
	if t, err := time.Parse("02-01-2006", dateStr); err == nil {
		return t, nil
	}

	// Handle MM/DD/YYYY format
	if t, err := time.Parse("01/02/2006", dateStr); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unrecognized date format: %s", dateStr)
}

// getNextWeekday returns the next occurrence of the given weekday
func getNextWeekday(from time.Time, targetWeekday time.Weekday) time.Time {
	daysUntil := int(targetWeekday) - int(from.Weekday())
	if daysUntil <= 0 {
		daysUntil += 7
	}
	return from.AddDate(0, 0, daysUntil)
}
