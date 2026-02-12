// Package prayer provides prayer times calculation helpers and data
package prayer

import (
	"fmt"
	"time"
)

// Prayer represents a single prayer time
type Prayer struct {
	Name     string
	Time     time.Time
	TimeStr  string
	Index    int
	IsPassed bool
	IsNext   bool
}

// PrayerTimes holds all prayer times for a day
type PrayerTimes struct {
	Date     time.Time
	Location string
	Prayers  []Prayer
	Timezone *time.Location
}

// ParseTime parses a time string (HH:MM) into a time.Time for the given date
func ParseTime(timeStr string, date time.Time, tz *time.Location) (time.Time, error) {
	var hour, minute int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	return time.Date(
		date.Year(), date.Month(), date.Day(),
		hour, minute, 0, 0, tz,
	), nil
}

// GetNextPrayer returns the next prayer from the list
func GetNextPrayer(prayers []Prayer, now time.Time) *Prayer {
	for i := range prayers {
		if !prayers[i].IsPassed {
			prayers[i].IsNext = true
			return &prayers[i]
		}
	}
	return nil
}

// CalculateTimeDiff calculates the difference between two times in minutes
func CalculateTimeDiff(from, to time.Time) int {
	diff := to.Sub(from)
	return int(diff.Minutes())
}

// FormatDuration formats a duration in minutes to a human-readable string
func FormatDuration(minutes int) string {
	if minutes < 0 {
		return "passed"
	}
	if minutes < 60 {
		return fmt.Sprintf("%d min", minutes)
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, mins)
}

// GetCompassDirection converts degrees to compass direction
func GetCompassDirection(degrees float64) string {
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

// IsPrayerPassed checks if a prayer time has passed
func IsPrayerPassed(prayerTime, now time.Time) bool {
	return now.After(prayerTime)
}

// PrayerIndex represents the index of prayers
const (
	FajrIndex = iota
	SunriseIndex
	DhuhrIndex
	AsrIndex
	MaghribIndex
	IshaIndex
	MidnightIndex
)

// PrayerNameByIndex returns the prayer name for a given index
func PrayerNameByIndex(index int) string {
	names := []string{"Fajr", "Sunrise", "Dhuhr", "Asr", "Maghrib", "Isha", "Midnight"}
	if index >= 0 && index < len(names) {
		return names[index]
	}
	return ""
}
