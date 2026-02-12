// Package ui provides interactive user interface components
package ui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/anashaat/pray-cli/internal/config"
	"github.com/anashaat/pray-cli/internal/location"
)

// Wizard handles the interactive setup process
type Wizard struct {
	reader io.Reader
	writer io.Writer
	cfg    *config.Config
}

// NewWizard creates a new setup wizard
func NewWizard() *Wizard {
	return &Wizard{
		reader: os.Stdin,
		writer: os.Stdout,
		cfg:    config.DefaultConfig(),
	}
}

// Run executes the setup wizard
func (w *Wizard) Run() (*config.Config, error) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Fprintln(w.writer)
	fmt.Fprintln(w.writer, "🕌 "+cyan("Prayer Times CLI - Initial Setup"))
	fmt.Fprintln(w.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w.writer)

	// Step 1: Location Setup
	fmt.Fprintln(w.writer, yellow("Step 1/5: Location Setup"))
	fmt.Fprintln(w.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w.writer)
	fmt.Fprintln(w.writer, "How would you like to set your location?")
	fmt.Fprintln(w.writer)
	fmt.Fprintln(w.writer, "  [1] Auto-detect from IP address (recommended)")
	fmt.Fprintln(w.writer, "  [2] Enter city or address manually")
	fmt.Fprintln(w.writer, "  [3] Enter coordinates (latitude/longitude)")
	fmt.Fprintln(w.writer)

	choice := w.promptDefault("Choose an option", "1")

	switch choice {
	case "1":
		// Auto-detect
		fmt.Fprintln(w.writer)
		fmt.Fprintln(w.writer, "🔍 Detecting your location...")

		detector := location.NewDetector()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		loc, err := detector.DetectFromIP(ctx)
		if err != nil {
			fmt.Fprintf(w.writer, "❌ Failed to detect location: %v\n", err)
			fmt.Fprintln(w.writer, "Please enter your location manually:")
			w.cfg.Location.Address = w.prompt("City or address")
		} else {
			fmt.Fprintf(w.writer, "%s Detected: %s\n", green("✓"), cyan(loc.GetDisplayAddress()))
			fmt.Fprintf(w.writer, "  Coordinates: %.4f°N, %.4f°E\n", loc.Latitude, loc.Longitude)
			fmt.Fprintf(w.writer, "  Timezone: %s\n", loc.Timezone)
			fmt.Fprintln(w.writer)

			if w.confirmDefault("Is this correct?", true) {
				w.cfg.Location = *loc
			} else {
				w.cfg.Location.Address = w.prompt("Enter correct city or address")
			}
		}

	case "2":
		// Manual address
		w.cfg.Location.Address = w.prompt("Enter city or address (e.g., Cairo, Egypt)")
		w.cfg.Location.Source = "manual"

	case "3":
		// Coordinates
		latStr := w.prompt("Latitude (e.g., 30.0444)")
		lonStr := w.prompt("Longitude (e.g., 31.2357)")

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			fmt.Fprintln(w.writer, "Invalid latitude, using 0")
			lat = 0
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			fmt.Fprintln(w.writer, "Invalid longitude, using 0")
			lon = 0
		}

		w.cfg.Location.Latitude = lat
		w.cfg.Location.Longitude = lon
		w.cfg.Location.Source = "manual"
	}

	fmt.Fprintln(w.writer)

	// Step 2: Calculation Method
	fmt.Fprintln(w.writer, yellow("Step 2/5: Calculation Method"))
	fmt.Fprintln(w.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w.writer)
	fmt.Fprintln(w.writer, "Select your calculation method:")
	fmt.Fprintln(w.writer)
	fmt.Fprintln(w.writer, "  [1] Karachi - Pakistan, Bangladesh, India")
	fmt.Fprintln(w.writer, "  [2] ISNA - North America")
	fmt.Fprintln(w.writer, "  [3] MWL - Europe, Far East")
	fmt.Fprintln(w.writer, "  [4] Umm al-Qura - Saudi Arabia")
	fmt.Fprintln(w.writer, "  [5] Egyptian - Egypt, Africa, Syria (default)")
	fmt.Fprintln(w.writer, "  [12] Diyanet - Turkey")
	fmt.Fprintln(w.writer, "  [0] Other (enter ID manually)")
	fmt.Fprintln(w.writer)

	methodChoice := w.promptDefault("Select method", "5")
	methodID, err := strconv.Atoi(methodChoice)
	if err != nil || methodID < 0 || methodID > 23 {
		methodID = 5
	}
	w.cfg.Method = methodID

	fmt.Fprintln(w.writer)

	// Step 3: Language
	fmt.Fprintln(w.writer, yellow("Step 3/5: Language"))
	fmt.Fprintln(w.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w.writer)
	fmt.Fprintln(w.writer, "Select your preferred language:")
	fmt.Fprintln(w.writer)
	fmt.Fprintln(w.writer, "  [1] English")
	fmt.Fprintln(w.writer, "  [2] العربية (Arabic)")
	fmt.Fprintln(w.writer)

	langChoice := w.promptDefault("Choose language", "1")
	if langChoice == "2" {
		w.cfg.Language = "ar"
	} else {
		w.cfg.Language = "en"
	}

	fmt.Fprintln(w.writer)

	// Step 4: Display Features
	fmt.Fprintln(w.writer, yellow("Step 4/5: Display Features"))
	fmt.Fprintln(w.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w.writer)

	w.cfg.Features.Qibla = w.confirm("Include Qibla direction?")
	w.cfg.Features.Dua = w.confirm("Include daily Du'a (Adhkar)?")

	if w.confirmDefault("Display Hijri date?", true) {
		fmt.Fprintln(w.writer, "  Where should Hijri date appear?")
		fmt.Fprintln(w.writer, "    [1] In description (default)")
		fmt.Fprintln(w.writer, "    [2] In title")
		fmt.Fprintln(w.writer, "    [3] Both")
		hijriChoice := w.promptDefault("  Choose", "1")
		switch hijriChoice {
		case "2":
			w.cfg.Features.Hijri = "title"
		case "3":
			w.cfg.Features.Hijri = "both"
		default:
			w.cfg.Features.Hijri = "desc"
		}
	} else {
		w.cfg.Features.Hijri = "none"
	}

	fmt.Fprintln(w.writer)

	// Step 5: Special Features
	fmt.Fprintln(w.writer, yellow("Step 5/5: Special Features"))
	fmt.Fprintln(w.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w.writer)

	w.cfg.Jumuah.Enabled = w.confirm("Enable Jumu'ah (Friday prayer)?")
	w.cfg.Ramadan.Enabled = w.confirm("Enable Ramadan mode?")
	w.cfg.Features.TravelerMode = w.confirm("Are you traveling (Qasr mode)?")

	fmt.Fprintln(w.writer)
	fmt.Fprintln(w.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w.writer, green("✨ Setup Complete!"))
	fmt.Fprintln(w.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w.writer)

	return w.cfg, nil
}

// prompt asks for user input
func (w *Wizard) prompt(question string) string {
	fmt.Fprintf(w.writer, "%s: ", question)

	scanner := bufio.NewScanner(w.reader)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// promptDefault asks for user input with a default value
func (w *Wizard) promptDefault(question, defaultVal string) string {
	fmt.Fprintf(w.writer, "%s [%s]: ", question, defaultVal)

	scanner := bufio.NewScanner(w.reader)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			return defaultVal
		}
		return input
	}
	return defaultVal
}

// confirm asks for yes/no confirmation (default: no)
func (w *Wizard) confirm(question string) bool {
	fmt.Fprintf(w.writer, "%s [y/N]: ", question)

	scanner := bufio.NewScanner(w.reader)
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		return input == "y" || input == "yes"
	}
	return false
}

// confirmDefault asks for yes/no confirmation with a default
func (w *Wizard) confirmDefault(question string, defaultYes bool) bool {
	prompt := "[y/N]"
	if defaultYes {
		prompt = "[Y/n]"
	}
	fmt.Fprintf(w.writer, "%s %s: ", question, prompt)

	scanner := bufio.NewScanner(w.reader)
	if scanner.Scan() {
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if input == "" {
			return defaultYes
		}
		return input == "y" || input == "yes"
	}
	return defaultYes
}
