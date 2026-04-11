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

	"github.com/AbdElrahmaN31/pray-cli/internal/api"
	"github.com/AbdElrahmaN31/pray-cli/internal/config"
	"github.com/AbdElrahmaN31/pray-cli/pkg/prayer"
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

	// Resolve location
	loc, err := resolveLocation()
	if err != nil {
		return err
	}
	if loc == nil {
		fmt.Println("👋 No location configured. Run 'pray init' or 'pray config detect --save'")
		return nil
	}

	methodID := resolveMethod()

	// Create API client with caching
	client, err := createAPIClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.APITimeout)*time.Second)
	defer cancel()

	// Fetch prayer times
	resp, err := fetchPrayerTimes(ctx, client, loc, methodID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to fetch prayer times: %w", err)
	}

	// Get Qibla if enabled and coords available
	var qibla *api.QiblaData
	if ShouldShowQibla() && loc.Latitude != 0 && loc.Longitude != 0 {
		qiblaResp, err := client.GetQibla(ctx, loc.Latitude, loc.Longitude)
		if err == nil {
			qibla = &qiblaResp.Data
		}
	}

	// Load timezone
	var tzLoc *time.Location
	if loc.Timezone != "" {
		tzLoc, err = time.LoadLocation(loc.Timezone)
		if err != nil {
			tzLoc = time.Local
		}
	} else if resp.Data.Meta.Timezone != "" {
		tzLoc, err = time.LoadLocation(resp.Data.Meta.Timezone)
		if err != nil {
			tzLoc = time.Local
		}
	} else {
		tzLoc = time.Local
	}

	// Determine language
	lang := GetLanguage()
	prayerNames := config.PrayerNames
	if lang == "ar" {
		prayerNames = config.PrayerNamesArabic
	}

	// Parse prayer times
	timings := resp.Data.Timings
	prayers := []struct {
		name  string
		time  string
		emoji string
	}{
		{prayerNames[0], cleanTime(timings.Fajr), config.PrayerEmojis["Fajr"]},
		{prayerNames[1], cleanTime(timings.Sunrise), config.PrayerEmojis["Sunrise"]},
		{prayerNames[2], cleanTime(timings.Dhuhr), config.PrayerEmojis["Dhuhr"]},
		{prayerNames[3], cleanTime(timings.Asr), config.PrayerEmojis["Asr"]},
		{prayerNames[4], cleanTime(timings.Maghrib), config.PrayerEmojis["Maghrib"]},
		{prayerNames[5], cleanTime(timings.Isha), config.PrayerEmojis["Isha"]},
	}

	// Hijri info
	hijriFormat := GetHijriFormat()
	var hijriStr string
	if hijriFormat != "none" && hijriFormat != "" {
		h := resp.Data.Date.Hijri
		hijriStr = fmt.Sprintf("%s %s %s", h.Day, h.Month.En, h.Year)
		if lang == "ar" {
			hijriStr = fmt.Sprintf("%s %s %s", h.Day, h.Month.Ar, h.Year)
		}
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
			now := time.Now().In(tzLoc)

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
			if hijriStr != "" {
				fmt.Printf("  📅 %s\n", dim(hijriStr))
			}
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
			fmt.Printf("  %s %s\n", "📍", dim(loc.DisplayName))
			fmt.Printf("  %s %s\n", "⚙️", dim(config.GetMethodName(methodID)))
			fmt.Printf("  %s %s\n", "🕐", dim(now.Format("15:04:05")))
			if qibla != nil {
				compass := prayer.GetCompassDirection(qibla.Direction)
				fmt.Printf("  🧭 %s\n", dim(fmt.Sprintf("Qibla: %s (%.1f°)", compass, qibla.Direction)))
			}
			fmt.Println()
			fmt.Printf("  %s\n", dim("Press Ctrl+C to exit"))
		}
	}
}
