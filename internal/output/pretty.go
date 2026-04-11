// Package output provides output formatting for prayer times
package output

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/AbdElrahmaN31/pray-cli/pkg/prayer"
)

// PrettyFormatter formats output with colors and emojis
type PrettyFormatter struct{}

// Format writes the prayer times in a pretty format with colors and emojis
func (f *PrettyFormatter) Format(w io.Writer, data *PrayerData) error {
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
	bold := color.New(color.Bold).SprintFunc()

	if data.NoColor {
		color.NoColor = true
	}

	// Header
	fmt.Fprintln(w)
	fmt.Fprintf(w, "🕌 %s\n", bold(fmt.Sprintf("Prayer Times for %s", data.Location)))
	fmt.Fprintf(w, "📅 %s", date.Readable)

	if data.ShowHijri && data.HijriFormat != "none" {
		hijri := date.Hijri
		fmt.Fprintf(w, " | %s %s %s", hijri.Day, hijri.Month.En, hijri.Year)
	}
	fmt.Fprintln(w)

	// Hijri holidays
	if data.HijriHolidays && len(date.Hijri.Holidays) > 0 {
		for _, holiday := range date.Hijri.Holidays {
			fmt.Fprintf(w, "🎉 %s\n", holiday)
		}
	}

	fmt.Fprintln(w)

	// Prayers
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
		{"Midnight", cleanTime(timings.Midnight), "🌃"},
	}

	// Get current time
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
		for _, p := range parts {
			offset, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				offset = 0
			}
			iqamaOffsets = append(iqamaOffsets, offset)
		}
	}

	// Print prayers
	for i, p := range prayers {
		status := ""
		prayerTime, err := parseTimeToday(p.time, now)

		prayerDisplay := fmt.Sprintf("%s %-8s  %s", p.emoji, p.name, p.time)

		// Add iqama time
		if data.IqamaEnabled && i < len(iqamaOffsets) && iqamaOffsets[i] > 0 {
			if err == nil {
				iqamaTime := prayerTime.Add(time.Duration(iqamaOffsets[i]) * time.Minute)
				prayerDisplay += dim(fmt.Sprintf("  (Iqama: %s)", iqamaTime.Format("15:04")))
			}
		}

		if err == nil {
			if now.After(prayerTime) {
				status = dim("✓ Passed")
			} else if i == nextPrayerIdx {
				mins := int(time.Until(prayerTime).Minutes())
				status = yellow(fmt.Sprintf("▶ Next prayer in %s", formatMinutes(mins)))
				prayerDisplay = cyan(prayerDisplay)
			}
		}

		if status != "" {
			fmt.Fprintf(w, "%s  %s\n", prayerDisplay, status)
		} else {
			fmt.Fprintf(w, "%s\n", prayerDisplay)
		}
	}

	fmt.Fprintln(w)

	// Qibla
	if data.ShowQibla && data.Qibla != nil {
		compass := getCompassDirection(data.Qibla.Direction)
		fmt.Fprintf(w, "🧭 Qibla Direction: %s (%.1f°)\n", green(compass), data.Qibla.Direction)
	}

	// Du'a
	if data.ShowDua {
		dua := prayer.GetDailyDua(time.Now())
		if dua != nil {
			fmt.Fprintf(w, "📖 Today's Du'a:\n")
			fmt.Fprintf(w, "   %s\n", dua.Arabic)
			fmt.Fprintf(w, "   %s\n", dim(dua.Transliteration))
			fmt.Fprintf(w, "   %s\n", dim(fmt.Sprintf("\"%s\"", dua.Translation)))
			fmt.Fprintf(w, "   %s\n", dim(fmt.Sprintf("— %s", dua.Reference)))
		}
	}

	// Method
	fmt.Fprintf(w, "⚙️  Method: %s\n", dim(data.Method))
	fmt.Fprintln(w)

	return nil
}
