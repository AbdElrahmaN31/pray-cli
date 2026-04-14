// Package output provides output formatting for prayer times
package output

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	"github.com/AbdElrahmaN31/pray-cli/pkg/prayer"
)

// TableFormatter formats output as an ASCII table
type TableFormatter struct{}

// Format writes the prayer times as a table
func (f *TableFormatter) Format(w io.Writer, data *PrayerData) error {
	if data.Response == nil {
		return fmt.Errorf("no prayer times data")
	}

	resp := data.Response
	timings := resp.Data.Timings
	date := resp.Data.Date

	// Colors
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	if data.NoColor {
		color.NoColor = true
	}

	// Header
	fmt.Fprintln(w)
	fmt.Fprintf(w, "┌──────────────────────────────────────────────────┐\n")
	fmt.Fprintf(w, "│%s│\n", centerText(fmt.Sprintf(t(data.Language, "prayer_times_dash"), data.Location), 50))
	fmt.Fprintf(w, "│%s│\n", centerText(date.Readable, 50))

	if data.ShowHijri && data.HijriFormat != "none" {
		hijri := date.Hijri
		hijriStr := fmt.Sprintf("%s %s %s", hijri.Day, hijriMonth(data.Language, hijri.Month), hijri.Year)
		fmt.Fprintf(w, "│%s│\n", centerText(hijriStr, 50))
	}

	// Hijri holidays
	if data.HijriHolidays && len(date.Hijri.Holidays) > 0 {
		for _, holiday := range date.Hijri.Holidays {
			fmt.Fprintf(w, "│%s│\n", centerText("🎉 "+holiday, 50))
		}
	}

	// Create prayers list with status
	prayers := []struct {
		key   string
		name  string
		time  string
		emoji string
	}{
		{keyFajr, t(data.Language, keyFajr), cleanTime(timings.Fajr), "🌅"},
		{keySunrise, t(data.Language, keySunrise), cleanTime(timings.Sunrise), "🌄"},
		{keyDhuhr, t(data.Language, keyDhuhr), cleanTime(timings.Dhuhr), "☀️"},
		{keyAsr, t(data.Language, keyAsr), cleanTime(timings.Asr), "🌤️"},
		{keyMaghrib, t(data.Language, keyMaghrib), cleanTime(timings.Maghrib), "🌆"},
		{keyIsha, t(data.Language, keyIsha), cleanTime(timings.Isha), "🌙"},
		{keyMidnight, t(data.Language, keyMidnight), cleanTime(timings.Midnight), "🌃"},
	}

	// Get current time for status
	now := time.Now()
	tz := resp.Data.Meta.Timezone
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err == nil {
			now = time.Now().In(loc)
		}
	}

	// Find next prayer
	var nextPrayerIdx int = -1
	for i, p := range prayers {
		prayerTime, err := parseTimeToday(p.time, now)
		if err != nil {
			continue
		}
		if now.Before(prayerTime) {
			nextPrayerIdx = i
			break
		}
	}

	// Parse iqama offsets
	var iqamaOffsets []int
	if data.IqamaEnabled && data.IqamaOffsets != "" {
		parts := strings.Split(data.IqamaOffsets, ",")
		for _, pt := range parts {
			offset, err := strconv.Atoi(strings.TrimSpace(pt))
			if err != nil {
				offset = 0
			}
			iqamaOffsets = append(iqamaOffsets, offset)
		}
	}

	// Table
	table := tablewriter.NewTable(os.Stdout)
	if data.IqamaEnabled && len(iqamaOffsets) > 0 {
		table.Header(t(data.Language, "prayer"), t(data.Language, "time"), t(data.Language, "iqama"), t(data.Language, "status"))
	} else {
		table.Header(t(data.Language, "prayer"), t(data.Language, "time"), t(data.Language, "status"))
	}

	for i, p := range prayers {
		status := ""
		prayerName := p.name
		prayerTime := p.time

		prayerDateTime, err := parseTimeToday(p.time, now)
		if err == nil {
			if now.After(prayerDateTime) {
				status = dim("✓ " + t(data.Language, "passed"))
			} else if i == nextPrayerIdx {
				mins := int(time.Until(prayerDateTime).Minutes())
				status = yellow(fmt.Sprintf("▶ %s (%s)", t(data.Language, "next"), formatDuration(data.Language, mins)))
				prayerName = cyan(p.name)
				prayerTime = green(p.time)
			}
		}

		if data.IqamaEnabled && len(iqamaOffsets) > 0 {
			iqamaStr := "-"
			if i < len(iqamaOffsets) && iqamaOffsets[i] > 0 && err == nil {
				iqamaTime := prayerDateTime.Add(time.Duration(iqamaOffsets[i]) * time.Minute)
				iqamaStr = iqamaTime.Format("15:04")
			}
			_ = table.Append(prayerName, prayerTime, iqamaStr, status)
		} else {
			_ = table.Append(prayerName, prayerTime, status)
		}
	}

	fmt.Fprintln(w, "├──────────────────────────────────────────────────┤")
	_ = table.Render()

	// Footer with Qibla and Method
	fmt.Fprintf(w, "├──────────────────────────────────────────────────┤\n")
	if data.ShowQibla && data.Qibla != nil {
		compass := getCompassDirection(data.Qibla.Direction)
		qiblaLabel := t(data.Language, "qibla")
		line := fmt.Sprintf(" %s: %.1f° (%s)", qiblaLabel, data.Qibla.Direction, compass)
		fmt.Fprintf(w, "│%s%s│\n", line, strings.Repeat(" ", padWidth(line, 50)))
	}
	methodLine := fmt.Sprintf(" %s: %s", t(data.Language, "method"), data.Method)
	fmt.Fprintf(w, "│%s%s│\n", methodLine, strings.Repeat(" ", padWidth(methodLine, 50)))
	fmt.Fprintf(w, "└──────────────────────────────────────────────────┘\n")

	// Du'a
	if data.ShowDua {
		dua := prayer.GetDailyDua(time.Now())
		if dua != nil {
			dim := color.New(color.Faint).SprintFunc()
			fmt.Fprintln(w)
			fmt.Fprintf(w, "📖 %s\n", dua.Arabic)
			fmt.Fprintf(w, "   %s\n", dim(fmt.Sprintf("\"%s\" — %s", dua.Translation, dua.Reference)))
		}
	}

	return nil
}

// centerText centers text within a given width (rune-based, not byte-based).
func centerText(text string, width int) string {
	runeLen := utf8.RuneCountInString(text)
	if runeLen >= width {
		return text
	}
	padding := (width - runeLen) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-padding-runeLen)
}

// padWidth returns the number of spaces needed to pad a string to the given rune width.
// Returns 0 if the string is already wider than the target width.
func padWidth(s string, width int) int {
	runeLen := utf8.RuneCountInString(s)
	if runeLen >= width {
		return 0
	}
	return width - runeLen
}

// cleanTime removes timezone info from time string (e.g., "05:23 (EET)" -> "05:23")
func cleanTime(timeStr string) string {
	parts := strings.Split(timeStr, " ")
	if len(parts) > 0 {
		return parts[0]
	}
	return timeStr
}

// parseTimeToday parses a time string (HH:MM) and returns it as time.Time for today
func parseTimeToday(timeStr string, now time.Time) (time.Time, error) {
	var hour, minute int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location()), nil
}

// formatMinutes formats minutes into a human-readable string
func formatMinutes(mins int) string {
	if mins < 60 {
		return fmt.Sprintf("%d min", mins)
	}
	hours := mins / 60
	remaining := mins % 60
	if remaining == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, remaining)
}

// getCompassDirection converts degrees to compass direction
func getCompassDirection(degrees float64) string {
	// Normalize to 0-360
	for degrees < 0 {
		degrees += 360
	}
	for degrees >= 360 {
		degrees -= 360
	}

	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	index := int((degrees+11.25)/22.5) % 16
	return directions[index]
}
