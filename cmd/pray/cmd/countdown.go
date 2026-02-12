package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/anashaat/pray-cli/internal/api"
	"github.com/anashaat/pray-cli/internal/config"
	"github.com/anashaat/pray-cli/internal/location"
)

var countdownCmd = &cobra.Command{
	Use:   "countdown",
	Short: "Live countdown to next prayer",
	Long: `Display a live countdown to the next prayer time.

The countdown updates every second and shows:
  - Next prayer name and time
  - Time remaining (hours, minutes, seconds)
  - Current local time

Press Ctrl+C to exit.`,
	RunE: runCountdownCommand,
}

func init() {
	rootCmd.AddCommand(countdownCmd)
}

func runCountdownCommand(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Determine location
	var lat, lon float64
	var locationStr string
	var tz string

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

	// Fetch prayer times
	client := api.NewClient(api.WithTimeout(time.Duration(cfg.APITimeout) * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.APITimeout)*time.Second)
	defer cancel()

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

	// Load timezone
	var loc *time.Location
	if tz != "" {
		loc, err = time.LoadLocation(tz)
		if err != nil {
			loc = time.Local
		}
	} else if resp.Data.Meta.Timezone != "" {
		loc, err = time.LoadLocation(resp.Data.Meta.Timezone)
		if err != nil {
			loc = time.Local
		}
	} else {
		loc = time.Local
	}

	// Parse prayer times
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

	// Colors
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	if noColor {
		color.NoColor = true
	}

	// Set up signal handling for clean exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create ticker for updates
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[H\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit

	for {
		select {
		case <-sigChan:
			fmt.Print("\033[?25h") // Show cursor
			fmt.Println("\n\n👋 Goodbye!")
			return nil

		case <-ticker.C:
			now := time.Now().In(loc)

			// Find next prayer
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

			// Clear screen and move cursor to top
			fmt.Print("\033[H\033[2J")

			// Header
			fmt.Println()
			fmt.Printf("  %s %s\n", "⏱️", cyan("Prayer Time Countdown"))
			fmt.Println("  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()

			if nextPrayer == nil {
				fmt.Printf("  %s\n", yellow("🌙 All prayers for today have passed"))
				fmt.Printf("  %s\n", dim("Tomorrow's Fajr: "+cleanTime(timings.Fajr)))
			} else {
				remaining := time.Until(nextPrayer.prayerTime)
				hours := int(remaining.Hours())
				minutes := int(remaining.Minutes()) % 60
				seconds := int(remaining.Seconds()) % 60

				fmt.Printf("  %s %s\n", nextPrayer.emoji, cyan(fmt.Sprintf("Next Prayer: %s", nextPrayer.name)))
				fmt.Printf("  %s\n", green(fmt.Sprintf("Time: %s", nextPrayer.time)))
				fmt.Println()
				fmt.Printf("  %s\n", yellow(fmt.Sprintf("    %02d : %02d : %02d", hours, minutes, seconds)))
				fmt.Printf("  %s\n", dim("    hr   min   sec"))
			}

			fmt.Println()
			fmt.Println("  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("  %s %s\n", "📍", dim(locationStr))
			fmt.Printf("  %s %s\n", "⚙️", dim(config.GetMethodName(methodID)))
			fmt.Printf("  %s %s\n", "🕐", dim(now.Format("15:04:05")))
			fmt.Println()
			fmt.Printf("  %s\n", dim("Press Ctrl+C to exit"))
		}
	}
}
