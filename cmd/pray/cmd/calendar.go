package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/anashaat/pray-cli/internal/calendar"
	"github.com/anashaat/pray-cli/internal/config"
	"github.com/anashaat/pray-cli/internal/ui"
)

var (
	calendarFile     string
	calendarMonths   int
	calendarDuration int
	calendarAlarm    string
	calendarColor    string
	calendarEvents   string
)

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar operations",
	Long:  `Generate and manage prayer times calendar files (ICS).`,
}

var calendarGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Download ICS calendar file",
	Long:  `Download an ICS calendar file with prayer times for your location.`,
	RunE:  runCalendarGet,
}

var calendarURLCmd = &cobra.Command{
	Use:   "url",
	Short: "Generate calendar URL",
	Long:  `Generate a subscription URL for your calendar app.`,
	RunE:  runCalendarURL,
}

var calendarSubscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Show subscription instructions",
	Long:  `Display instructions for subscribing to the prayer times calendar.`,
	RunE:  runCalendarSubscribe,
}

func init() {
	rootCmd.AddCommand(calendarCmd)
	calendarCmd.AddCommand(calendarGetCmd)
	calendarCmd.AddCommand(calendarURLCmd)
	calendarCmd.AddCommand(calendarSubscribeCmd)

	// Flags for calendar get
	calendarGetCmd.Flags().StringVarP(&calendarFile, "file", "f", "", "output file path")
	calendarGetCmd.Flags().IntVar(&calendarMonths, "months", 0, "number of months to generate (1-12)")
	calendarGetCmd.Flags().IntVarP(&calendarDuration, "duration", "d", 0, "event duration in minutes")
	calendarGetCmd.Flags().StringVar(&calendarAlarm, "alarm", "", "alarm offsets (e.g., '5,10,15')")
	calendarGetCmd.Flags().StringVar(&calendarColor, "color", "", "calendar color (e.g., '#1e90ff')")
	calendarGetCmd.Flags().StringVarP(&calendarEvents, "events", "e", "", "events to include ('all' or indices)")

	// Flags for calendar url
	calendarURLCmd.Flags().IntVar(&calendarMonths, "months", 0, "number of months to generate (1-12)")
	calendarURLCmd.Flags().IntVarP(&calendarDuration, "duration", "d", 0, "event duration in minutes")
	calendarURLCmd.Flags().StringVar(&calendarAlarm, "alarm", "", "alarm offsets (e.g., '5,10,15')")
	calendarURLCmd.Flags().StringVar(&calendarColor, "color", "", "calendar color (e.g., '#1e90ff')")
	calendarURLCmd.Flags().StringVarP(&calendarEvents, "events", "e", "", "events to include ('all' or indices)")
}

func runCalendarGet(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	green := color.New(color.FgGreen).SprintFunc()

	// Check location
	if !cfg.IsConfigured() && address == "" && latitude == 0 {
		fmt.Println("No location configured. Run 'pray config detect --save' or use -a flag.")
		return nil
	}

	// Build calendar params
	params := buildCalendarParams(cfg)

	// Generate URL
	icsURL := calendar.GenerateICSURL(params)

	// Determine output file
	outputFile := calendarFile
	if outputFile == "" {
		outputFile = calendar.GetDefaultFilename(cfg.Location.GetDisplayAddress())
	}

	// Use spinner for download
	spinner := ui.NewSpinner("Downloading calendar...")
	spinner.Start()

	// Download
	downloader := calendar.NewDownloader()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := downloader.DownloadToFile(ctx, icsURL, outputFile)
	if err != nil {
		spinner.Fail("Failed to download calendar")
		return fmt.Errorf("failed to download calendar: %w", err)
	}

	spinner.Stop()
	fmt.Printf("%s Calendar saved to: %s\n", green("✓"), outputFile)
	fmt.Println()
	fmt.Println("📍 Import this file into your calendar app:")
	fmt.Println("   - Google Calendar: Settings > Import & export > Import")
	fmt.Println("   - Apple Calendar: File > Import")
	fmt.Println("   - Outlook: File > Open > Import")
	fmt.Println()

	return nil
}

func runCalendarURL(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()
	cyan := color.New(color.FgCyan).SprintFunc()

	// Check location
	if !cfg.IsConfigured() && address == "" && latitude == 0 {
		fmt.Println("No location configured. Run 'pray config detect --save' or use -a flag.")
		return nil
	}

	// Build calendar params
	params := buildCalendarParams(cfg)

	// Generate URL
	icsURL := calendar.GenerateICSURL(params)

	fmt.Println()
	fmt.Println("📅 Calendar Subscription URL")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println(cyan(icsURL))
	fmt.Println()
	fmt.Println("Use this URL to subscribe in your calendar app.")
	fmt.Println("Run 'pray calendar subscribe' for detailed instructions.")
	fmt.Println()

	return nil
}

func runCalendarSubscribe(cmd *cobra.Command, args []string) error {
	cfg := GetConfig()

	// Check location
	if !cfg.IsConfigured() && address == "" && latitude == 0 {
		fmt.Println("No location configured. Run 'pray config detect --save' or use -a flag.")
		return nil
	}

	// Build calendar params
	params := buildCalendarParams(cfg)

	// Generate URL
	icsURL := calendar.GenerateICSURL(params)

	// Show instructions
	instructions := calendar.NewInstructions(icsURL)
	instructions.Print(os.Stdout)

	return nil
}

func buildCalendarParams(cfg *config.Config) *calendar.CalendarParams {
	params := calendar.NewCalendarParams()

	// Location
	if address != "" {
		params.WithAddress(address)
	} else if latitude != 0 || longitude != 0 {
		params.WithCoordinates(latitude, longitude)
	} else {
		params.WithCoordinates(cfg.Location.Latitude, cfg.Location.Longitude)
	}

	// Method
	methodID := cfg.Method
	if method != 0 {
		methodID = method
	}
	params.WithMethod(methodID)

	// Calendar settings from flags or config
	if calendarMonths > 0 {
		params.WithMonths(calendarMonths)
	} else if cfg.Calendar.Months > 0 {
		params.WithMonths(cfg.Calendar.Months)
	}

	if calendarDuration > 0 {
		params.WithDuration(calendarDuration)
	} else if cfg.Calendar.Duration > 0 {
		params.WithDuration(cfg.Calendar.Duration)
	}

	if calendarAlarm != "" {
		params.WithAlarm(calendarAlarm)
	} else if cfg.Calendar.Alarm != "" {
		params.WithAlarm(cfg.Calendar.Alarm)
	}

	if calendarColor != "" {
		params.WithColor(calendarColor)
	} else if cfg.Calendar.Color != "" {
		params.WithColor(cfg.Calendar.Color)
	}

	if calendarEvents != "" {
		params.Events = calendarEvents
	} else if cfg.Calendar.Events != "" {
		params.Events = cfg.Calendar.Events
	}

	// Language
	params.WithLanguage(cfg.Language)

	// Hijri
	params.Hijri = cfg.Features.Hijri

	// Features
	params.Qibla = cfg.Features.Qibla
	params.Dua = cfg.Features.Dua
	params.Traveler = cfg.Features.TravelerMode
	params.HijriHolidays = cfg.Features.HijriHolidays

	// Jumu'ah
	if cfg.Jumuah.Enabled {
		params.WithJumuah(true, cfg.Jumuah.Duration)
	}

	// Ramadan
	if cfg.Ramadan.Enabled {
		params.WithRamadan(true)
		params.IftarDuration = cfg.Ramadan.IftarDuration
		params.TaraweehDuration = cfg.Ramadan.TaraweehDuration
		params.SuhoorDuration = cfg.Ramadan.SuhoorDuration
	}

	// Iqama
	if cfg.Iqama.Enabled {
		params.Iqama = cfg.Iqama.Offsets
	}

	return params
}
