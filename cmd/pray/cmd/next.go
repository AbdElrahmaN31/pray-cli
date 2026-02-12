package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/AbdElrahmaN31/pray-cli/internal/api"
	"github.com/AbdElrahmaN31/pray-cli/internal/config"
	"github.com/AbdElrahmaN31/pray-cli/internal/location"
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next prayer",
	Long:  `Display information about the next upcoming prayer time.`,
	RunE:  runNextCommand,
}

func init() {
	rootCmd.AddCommand(nextCmd)
}

func runNextCommand(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Determine location
	var lat, lon float64
	var locationStr string
	var tz string

	// Priority: flags > config
	if autoDetect {
		detector := location.NewDetector()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		loc, err := detector.DetectFromIP(ctx)
		if err != nil {
			return fmt.Errorf("failed to auto-detect location: %w", err)
		}
		lat = loc.Latitude
		lon = loc.Longitude
		locationStr = loc.GetDisplayAddress()
		tz = loc.Timezone
	} else if address != "" {
		locationStr = address
	} else if latitude != 0 || longitude != 0 {
		lat = latitude
		lon = longitude
		locationStr = fmt.Sprintf("%.4f, %.4f", lat, lon)
	} else if cfg.IsConfigured() {
		lat = cfg.Location.Latitude
		lon = cfg.Location.Longitude
		locationStr = cfg.Location.GetDisplayAddress()
		tz = cfg.Location.Timezone
	} else {
		fmt.Println("👋 No location configured. Run 'pray init' or 'pray config detect --save'")
		return nil
	}

	// Determine method
	methodID := cfg.Method
	if method != 0 {
		methodID = method
	}

	// Create API client
	client := api.NewClient(api.WithTimeout(time.Duration(cfg.APITimeout) * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.APITimeout)*time.Second)
	defer cancel()

	// Build params
	params := api.NewPrayerTimesParams().
		WithDate(time.Now()).
		WithMethod(methodID)

	var resp *api.PrayerTimesResponse
	var err error

	if address != "" {
		params.WithAddress(address)
		resp, err = client.GetPrayerTimesByAddress(ctx, params)
	} else {
		params.WithCoordinates(lat, lon)
		if tz != "" {
			params.WithTimezone(tz)
		}
		resp, err = client.GetPrayerTimes(ctx, params)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch prayer times: %w", err)
	}

	// Get current time
	now := time.Now()
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err == nil {
			now = time.Now().In(loc)
		}
	}

	// Find next prayer
	timings := resp.Data.Timings
	prayers := []struct {
		name  string
		time  string
		emoji string
	}{
		{"Fajr", cleanTime(timings.Fajr), "🌅"},
		{"Sunrise", cleanTime(timings.Sunrise), "🌄"},
		{"Dhuhr", cleanTime(timings.Dhuhr), "☀️"},
		{"Asr", cleanTime(timings.Asr), "🌤️"},
		{"Maghrib", cleanTime(timings.Maghrib), "🌆"},
		{"Isha", cleanTime(timings.Isha), "🌙"},
	}

	var nextPrayer *struct {
		name       string
		time       string
		emoji      string
		prayerTime time.Time
	}

	for _, p := range prayers {
		prayerTime, err := parseTimeForToday(p.time, now)
		if err != nil {
			continue
		}
		if now.Before(prayerTime) {
			nextPrayer = &struct {
				name       string
				time       string
				emoji      string
				prayerTime time.Time
			}{p.name, p.time, p.emoji, prayerTime}
			break
		}
	}

	// Colors
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	if noColor {
		color.NoColor = true
	}

	// Output based on format
	if outputFormat == "json" {
		if nextPrayer != nil {
			mins := int(time.Until(nextPrayer.prayerTime).Minutes())
			fmt.Printf(`{"name":"%s","time":"%s","minutesUntil":%d,"location":"%s"}%s`,
				nextPrayer.name, nextPrayer.time, mins, locationStr, "\n")
		} else {
			fmt.Println(`{"name":null,"message":"All prayers for today have passed"}`)
		}
		return nil
	}

	// Pretty output
	fmt.Println()
	if nextPrayer == nil {
		fmt.Println("🌙 All prayers for today have passed")
		fmt.Printf("   Tomorrow's Fajr: %s\n", cleanTime(timings.Fajr))
	} else {
		mins := int(time.Until(nextPrayer.prayerTime).Minutes())

		fmt.Printf("%s %s\n", nextPrayer.emoji, cyan(fmt.Sprintf("Next Prayer: %s", nextPrayer.name)))
		fmt.Printf("   Time: %s\n", green(nextPrayer.time))
		fmt.Printf("   In:   %s\n", yellow(formatMinutesLong(mins)))
		fmt.Println()
		fmt.Printf("   %s\n", dim(fmt.Sprintf("Location: %s", locationStr)))
		fmt.Printf("   %s\n", dim(fmt.Sprintf("Method: %s", config.GetMethodName(methodID))))
	}
	fmt.Println()

	return nil
}

// cleanTime removes timezone info from time string
func cleanTime(timeStr string) string {
	for i, c := range timeStr {
		if c == ' ' || c == '(' {
			return timeStr[:i]
		}
	}
	return timeStr
}

// parseTimeForToday parses a time string and returns time.Time for today
func parseTimeForToday(timeStr string, now time.Time) (time.Time, error) {
	var hour, minute int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location()), nil
}

// formatMinutesLong formats minutes in a longer human-readable format
func formatMinutesLong(mins int) string {
	if mins < 0 {
		return "passed"
	}
	if mins < 60 {
		if mins == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", mins)
	}
	hours := mins / 60
	remaining := mins % 60

	hourStr := "hour"
	if hours > 1 {
		hourStr = "hours"
	}

	if remaining == 0 {
		return fmt.Sprintf("%d %s", hours, hourStr)
	}

	minStr := "minute"
	if remaining > 1 {
		minStr = "minutes"
	}
	return fmt.Sprintf("%d %s %d %s", hours, hourStr, remaining, minStr)
}
