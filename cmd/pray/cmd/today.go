package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/anashaat/pray-cli/internal/api"
	"github.com/anashaat/pray-cli/internal/config"
	"github.com/anashaat/pray-cli/internal/location"
	"github.com/anashaat/pray-cli/internal/output"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's prayer times",
	Long:  `Display prayer times for today based on your configured location.`,
	RunE:  runTodayCommand,
}

func init() {
	rootCmd.AddCommand(todayCmd)
}

func runTodayCommand(cmd *cobra.Command, args []string) error {
	return fetchAndDisplayPrayerTimes(cmd, time.Now())
}

// fetchAndDisplayPrayerTimes fetches and displays prayer times for a given date
func fetchAndDisplayPrayerTimes(cmd *cobra.Command, date time.Time) error {
	cfg := GetConfig()

	// Determine location
	var lat, lon float64
	var locationStr string
	var tz string
	var detectedLoc *location.Location

	// Priority: flags > config
	if autoDetect {
		// Auto-detect location
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
		detectedLoc = loc
	} else if address != "" {
		// Use address from flag
		locationStr = address
	} else if latitude != 0 || longitude != 0 {
		// Use coordinates from flags
		lat = latitude
		lon = longitude
		locationStr = fmt.Sprintf("%.4f, %.4f", lat, lon)
	} else if cfg.IsConfigured() {
		// Use config
		lat = cfg.Location.Latitude
		lon = cfg.Location.Longitude
		locationStr = cfg.Location.GetDisplayAddress()
		tz = cfg.Location.Timezone
	} else {
		fmt.Println("👋 Welcome! No location configured.")
		fmt.Println()
		fmt.Println("Set your location using one of these options:")
		fmt.Println("  pray config detect --save    Auto-detect from IP")
		fmt.Println("  pray --auto                  Auto-detect (one-time)")
		fmt.Println("  pray -a \"Cairo, Egypt\"       Specify a city")
		fmt.Println("  pray init                    Interactive setup")
		return nil
	}

	// Determine method
	methodID := cfg.Method
	if method != 0 {
		methodID = method
	}

	// Handle --save flag: save current settings to config
	if ShouldSaveConfig() {
		if detectedLoc != nil {
			cfg.Location = *detectedLoc
		} else if address != "" {
			cfg.Location.Address = address
			cfg.Location.Source = "manual"
		} else if latitude != 0 || longitude != 0 {
			cfg.Location.Latitude = latitude
			cfg.Location.Longitude = longitude
			cfg.Location.Source = "manual"
		}
		if method != 0 {
			cfg.Method = method
		}
		if language != "" {
			cfg.Language = language
		}
		if showQibla {
			cfg.Features.Qibla = true
		}
		if showDua {
			cfg.Features.Dua = true
		}
		if hijriFormat != "" {
			cfg.Features.Hijri = hijriFormat
		}
		if travelerMode {
			cfg.Features.TravelerMode = true
		}
		if jumuahMode {
			cfg.Jumuah.Enabled = true
		}
		if ramadanMode {
			cfg.Ramadan.Enabled = true
		}
		if outputFormat != "" {
			cfg.Output.Format = outputFormat
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		if !IsQuiet() {
			fmt.Println("✓ Settings saved to config")
		}
	}

	// Create API client
	client := api.NewClient(api.WithTimeout(time.Duration(cfg.APITimeout) * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.APITimeout)*time.Second)
	defer cancel()

	// Build params
	params := api.NewPrayerTimesParams().
		WithDate(date).
		WithMethod(methodID)

	var resp *api.PrayerTimesResponse
	var err error

	if address != "" {
		// Fetch by address
		params.WithAddress(address)
		resp, err = client.GetPrayerTimesByAddress(ctx, params)
	} else {
		// Fetch by coordinates
		params.WithCoordinates(lat, lon)
		if tz != "" {
			params.WithTimezone(tz)
		}
		resp, err = client.GetPrayerTimes(ctx, params)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch prayer times: %w", err)
	}

	// Get Qibla if enabled (use flag helpers)
	var qibla *api.QiblaData
	qiblaEnabled := ShouldShowQibla() || outputFormat == "json" || outputFormat == "webhook"
	if qiblaEnabled && (lat != 0 && lon != 0) {
		qiblaResp, err := client.GetQibla(ctx, lat, lon)
		if err == nil {
			qibla = &qiblaResp.Data
		}
	}

	// Use flag helpers for display options
	hijri := GetHijriFormat()
	lang := GetLanguage()

	// Prepare output data
	data := &output.PrayerData{
		Response:    resp,
		Location:    locationStr,
		Method:      config.GetMethodName(methodID),
		Qibla:       qibla,
		ShowQibla:   ShouldShowQibla(),
		ShowDua:     ShouldShowDua(),
		ShowHijri:   hijri != "none",
		HijriFormat: hijri,
		Language:    lang,
		NoColor:     noColor,
	}

	// Determine output format
	format := cfg.Output.Format
	if outputFormat != "" {
		format = outputFormat
	}

	// Get formatter
	formatter := output.GetFormatter(format)

	// Determine output destination
	outFile := GetOutputFile()
	if outFile != "" {
		// Write to file
		f, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()

		if err := formatter.Format(f, data); err != nil {
			return err
		}
		if !IsQuiet() {
			fmt.Printf("✓ Output saved to: %s\n", outFile)
		}
		return nil
	}

	// Write to stdout
	return formatter.Format(os.Stdout, data)
}
