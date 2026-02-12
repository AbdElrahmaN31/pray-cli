// Package output provides output formatting for prayer times
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// JSONOutput represents the JSON output structure
type JSONOutput struct {
	Date       DateOutput        `json:"date"`
	Location   LocationOutput    `json:"location"`
	Method     MethodOutput      `json:"method"`
	Timings    TimingsOutput     `json:"timings"`
	NextPrayer *NextPrayerOutput `json:"nextPrayer,omitempty"`
	Qibla      *QiblaOutput      `json:"qibla,omitempty"`
}

// DateOutput represents date information in JSON
type DateOutput struct {
	Gregorian string       `json:"gregorian"`
	Hijri     *HijriOutput `json:"hijri,omitempty"`
}

// HijriOutput represents Hijri date in JSON
type HijriOutput struct {
	Day   string      `json:"day"`
	Month MonthOutput `json:"month"`
	Year  string      `json:"year"`
}

// MonthOutput represents month information
type MonthOutput struct {
	Number int    `json:"number"`
	En     string `json:"en"`
	Ar     string `json:"ar,omitempty"`
}

// LocationOutput represents location in JSON
type LocationOutput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Address   string  `json:"address,omitempty"`
}

// MethodOutput represents calculation method in JSON
type MethodOutput struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TimingsOutput represents prayer times
type TimingsOutput struct {
	Fajr     string `json:"Fajr"`
	Sunrise  string `json:"Sunrise"`
	Dhuhr    string `json:"Dhuhr"`
	Asr      string `json:"Asr"`
	Maghrib  string `json:"Maghrib"`
	Isha     string `json:"Isha"`
	Midnight string `json:"Midnight"`
}

// NextPrayerOutput represents the next prayer
type NextPrayerOutput struct {
	Name         string `json:"name"`
	Time         string `json:"time"`
	MinutesUntil int    `json:"minutesUntil"`
}

// QiblaOutput represents Qibla direction
type QiblaOutput struct {
	Direction float64 `json:"direction"`
	Compass   string  `json:"compass"`
}

// Format writes the prayer times as JSON
func (f *JSONFormatter) Format(w io.Writer, data *PrayerData) error {
	if data.Response == nil {
		return fmt.Errorf("no prayer times data")
	}

	resp := data.Response
	timings := resp.Data.Timings
	date := resp.Data.Date
	meta := resp.Data.Meta

	output := JSONOutput{
		Date: DateOutput{
			Gregorian: date.Readable,
		},
		Location: LocationOutput{
			Latitude:  meta.Latitude,
			Longitude: meta.Longitude,
			Timezone:  meta.Timezone,
			Address:   data.Location,
		},
		Method: MethodOutput{
			ID:   meta.Method.ID,
			Name: meta.Method.Name,
		},
		Timings: TimingsOutput{
			Fajr:     cleanTime(timings.Fajr),
			Sunrise:  cleanTime(timings.Sunrise),
			Dhuhr:    cleanTime(timings.Dhuhr),
			Asr:      cleanTime(timings.Asr),
			Maghrib:  cleanTime(timings.Maghrib),
			Isha:     cleanTime(timings.Isha),
			Midnight: cleanTime(timings.Midnight),
		},
	}

	// Add Hijri date if enabled
	if data.ShowHijri && data.HijriFormat != "none" {
		hijri := date.Hijri
		output.Date.Hijri = &HijriOutput{
			Day:  hijri.Day,
			Year: hijri.Year,
			Month: MonthOutput{
				Number: hijri.Month.Number,
				En:     hijri.Month.En,
				Ar:     hijri.Month.Ar,
			},
		}
	}

	// Calculate next prayer
	now := time.Now()
	tz := meta.Timezone
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err == nil {
			now = time.Now().In(loc)
		}
	}

	prayers := []struct {
		name string
		time string
	}{
		{"Fajr", cleanTime(timings.Fajr)},
		{"Sunrise", cleanTime(timings.Sunrise)},
		{"Dhuhr", cleanTime(timings.Dhuhr)},
		{"Asr", cleanTime(timings.Asr)},
		{"Maghrib", cleanTime(timings.Maghrib)},
		{"Isha", cleanTime(timings.Isha)},
		{"Midnight", cleanTime(timings.Midnight)},
	}

	for _, p := range prayers {
		prayerTime, err := parseTimeToday(p.time, now)
		if err != nil {
			continue
		}
		if now.Before(prayerTime) {
			mins := int(time.Until(prayerTime).Minutes())
			output.NextPrayer = &NextPrayerOutput{
				Name:         p.name,
				Time:         p.time,
				MinutesUntil: mins,
			}
			break
		}
	}

	// Add Qibla if enabled
	if data.ShowQibla && data.Qibla != nil {
		output.Qibla = &QiblaOutput{
			Direction: data.Qibla.Direction,
			Compass:   getCompassDirection(data.Qibla.Direction),
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
