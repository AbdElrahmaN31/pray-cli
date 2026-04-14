package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/AbdElrahmaN31/pray-cli/internal/api"
	"github.com/AbdElrahmaN31/pray-cli/internal/config"
	"github.com/AbdElrahmaN31/pray-cli/pkg/prayer"
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

	// Get current time in the right timezone
	now := time.Now()
	if loc.Timezone != "" {
		if tzLoc, err := time.LoadLocation(loc.Timezone); err == nil {
			now = time.Now().In(tzLoc)
		}
	} else if resp.Data.Meta.Timezone != "" {
		if tzLoc, err := time.LoadLocation(resp.Data.Meta.Timezone); err == nil {
			now = time.Now().In(tzLoc)
		}
	}

	// Determine language
	lang := GetLanguage()
	prayerNames := config.PrayerNames
	if lang == "ar" {
		prayerNames = config.PrayerNamesArabic
	}

	// Find next prayer
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
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	if noColor {
		color.NoColor = true
	}

	// Determine output writer
	var w *os.File
	outFile := GetOutputFile()
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		w = f
	} else {
		w = os.Stdout
	}

	// Output based on format
	if outputFormat == "json" {
		result := buildNextJSON(nextPrayer, loc, qibla, hijriStr, cleanTime(timings.Fajr))
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return err
		}
		if outFile != "" && !IsQuiet() {
			fmt.Printf("✓ Output saved to: %s\n", outFile)
		}
		return nil
	}

	// Pretty output
	fmt.Fprintln(w)
	if hijriStr != "" {
		fmt.Fprintf(w, "  📅 %s\n", dim(hijriStr))
	}

	if nextPrayer == nil {
		fmt.Fprintln(w, "🌙 All prayers for today have passed")
		fmt.Fprintf(w, "   Tomorrow's Fajr: %s\n", cleanTime(timings.Fajr))
	} else {
		mins := int(time.Until(nextPrayer.prayerTime).Minutes())

		fmt.Fprintf(w, "%s %s\n", nextPrayer.emoji, cyan(fmt.Sprintf("Next Prayer: %s", nextPrayer.name)))
		fmt.Fprintf(w, "   Time: %s\n", green(nextPrayer.time))
		fmt.Fprintf(w, "   In:   %s\n", yellow(formatMinutesLong(mins)))
		fmt.Fprintln(w)
		fmt.Fprintf(w, "   %s\n", dim(fmt.Sprintf("Location: %s", loc.DisplayName)))
		fmt.Fprintf(w, "   %s\n", dim(fmt.Sprintf("Method: %s", config.GetMethodNameLang(methodID, lang))))
	}

	// Qibla
	if qibla != nil && ShouldShowQibla() {
		compass := prayer.GetCompassDirection(qibla.Direction)
		fmt.Fprintf(w, "   🧭 Qibla: %s (%.1f°)\n", compass, qibla.Direction)
	}

	fmt.Fprintln(w)

	if outFile != "" && !IsQuiet() {
		fmt.Printf("✓ Output saved to: %s\n", outFile)
	}

	return nil
}

// buildNextJSON builds a JSON-friendly map for the next prayer output
func buildNextJSON(nextPrayer *struct {
	name       string
	time       string
	emoji      string
	prayerTime time.Time
}, loc *ResolvedLocation, qibla *api.QiblaData, hijriStr, fajrTime string) map[string]interface{} {
	result := map[string]interface{}{}
	if nextPrayer != nil {
		mins := int(time.Until(nextPrayer.prayerTime).Minutes())
		result["name"] = nextPrayer.name
		result["time"] = nextPrayer.time
		result["minutesUntil"] = mins
		result["location"] = loc.DisplayName
	} else {
		result["name"] = nil
		result["message"] = "All prayers for today have passed"
		result["tomorrowFajr"] = fajrTime
	}
	if hijriStr != "" {
		result["hijri"] = hijriStr
	}
	if qibla != nil {
		result["qibla"] = map[string]interface{}{
			"direction": qibla.Direction,
			"compass":   prayer.GetCompassDirection(qibla.Direction),
		}
	}
	return result
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
