package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/AbdElrahmaN31/pray-cli/internal/api"
	"github.com/AbdElrahmaN31/pray-cli/internal/config"
	"github.com/AbdElrahmaN31/pray-cli/internal/output"
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

	// Resolve location
	loc, err := resolveLocation()
	if err != nil {
		return err
	}
	if loc == nil {
		printNoLocationMessage()
		return nil
	}

	methodID := resolveMethod()

	// Handle --save flag (only if --no-save is not set)
	if ShouldSaveConfig() && !noSaveConfig {
		if loc.Detected != nil {
			cfg.Location = *loc.Detected
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
		if iqamaEnabled {
			cfg.Iqama.Enabled = true
		}
		if hijriHolidays {
			cfg.Features.HijriHolidays = true
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

	// Create API client with caching
	client, err := createAPIClient()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.APITimeout)*time.Second)
	defer cancel()

	// Fetch prayer times
	resp, err := fetchPrayerTimes(ctx, client, loc, methodID, date)
	if err != nil {
		return fmt.Errorf("failed to fetch prayer times: %w", err)
	}

	// Get Qibla if enabled (use flag helpers)
	var qibla *api.QiblaData
	qiblaEnabled := ShouldShowQibla() || outputFormat == "json" || outputFormat == "webhook"
	if qiblaEnabled && (loc.Latitude != 0 && loc.Longitude != 0) {
		qiblaResp, err := client.GetQibla(ctx, loc.Latitude, loc.Longitude)
		if err == nil {
			qibla = &qiblaResp.Data
		}
	}

	// Use flag helpers for display options
	hijri := GetHijriFormat()
	lang := GetLanguage()

	// Prepare output data
	data := &output.PrayerData{
		Response:      resp,
		Location:      loc.DisplayName,
		Method:        config.GetMethodNameLang(methodID, lang),
		Qibla:         qibla,
		ShowQibla:     ShouldShowQibla(),
		ShowDua:       ShouldShowDua(),
		ShowHijri:     hijri != "none",
		HijriFormat:   hijri,
		Language:      lang,
		NoColor:       noColor,
		IqamaEnabled:  IsIqamaEnabled(),
		IqamaOffsets:  GetConfig().Iqama.Offsets,
		HijriHolidays: IsHijriHolidays(),
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
