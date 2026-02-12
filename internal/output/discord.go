// Package output provides output formatting for prayer times
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// DiscordFormatter formats output as Discord embed JSON
type DiscordFormatter struct{}

// DiscordMessage represents a Discord message with embeds
type DiscordMessage struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

// DiscordEmbed represents a Discord embed
type DiscordEmbed struct {
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	Color       int            `json:"color"`
	Fields      []DiscordField `json:"fields,omitempty"`
	Footer      *DiscordFooter `json:"footer,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
}

// DiscordField represents a field in a Discord embed
type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// DiscordFooter represents the footer of a Discord embed
type DiscordFooter struct {
	Text string `json:"text"`
}

// Format writes the prayer times as Discord embed JSON
func (f *DiscordFormatter) Format(w io.Writer, data *PrayerData) error {
	if data.Response == nil {
		return fmt.Errorf("no prayer times data")
	}

	resp := data.Response
	timings := resp.Data.Timings
	date := resp.Data.Date
	meta := resp.Data.Meta

	// Get current time for next prayer calculation
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
	}

	// Find next prayer
	nextPrayer := ""
	for _, p := range prayers {
		prayerTime, err := parseTimeToday(p.time, now)
		if err != nil {
			continue
		}
		if now.Before(prayerTime) {
			nextPrayer = p.name
			break
		}
	}

	// Create fields
	fields := make([]DiscordField, 0)
	for _, p := range prayers {
		value := p.time
		if p.name == nextPrayer {
			value = fmt.Sprintf("%s ▶️", p.time)
		}
		fields = append(fields, DiscordField{
			Name:   p.name,
			Value:  value,
			Inline: true,
		})
	}

	// Discord color (blue: 0x1DA1F2 = 1942002)
	message := DiscordMessage{
		Embeds: []DiscordEmbed{
			{
				Title:       "🕌 Prayer Times",
				Description: fmt.Sprintf("**%s**\n%s", data.Location, date.Readable),
				Color:       1942002,
				Fields:      fields,
				Footer: &DiscordFooter{
					Text: fmt.Sprintf("Method: %s", data.Method),
				},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(message)
}

// WebhookFormatter formats output as a detailed webhook JSON
type WebhookFormatter struct{}

// WebhookOutput represents a detailed webhook payload
type WebhookOutput struct {
	Date       DateOutput         `json:"date"`
	Location   LocationOutput     `json:"location"`
	Timings    TimingsOutput      `json:"timings"`
	NextPrayer *WebhookNextPrayer `json:"nextPrayer,omitempty"`
	Qibla      *QiblaOutput       `json:"qibla,omitempty"`
	ServerTime string             `json:"serverTime"`
}

// WebhookNextPrayer includes additional timestamp info
type WebhookNextPrayer struct {
	Name         string `json:"name"`
	Time         string `json:"time"`
	ISO          string `json:"iso"`
	Timestamp    int64  `json:"timestamp"`
	MinutesUntil int    `json:"minutesUntil"`
}

// Format writes the prayer times as detailed webhook JSON
func (f *WebhookFormatter) Format(w io.Writer, data *PrayerData) error {
	if data.Response == nil {
		return fmt.Errorf("no prayer times data")
	}

	resp := data.Response
	timings := resp.Data.Timings
	date := resp.Data.Date
	meta := resp.Data.Meta

	now := time.Now()
	tz := meta.Timezone
	var loc *time.Location
	if tz != "" {
		var err error
		loc, err = time.LoadLocation(tz)
		if err == nil {
			now = time.Now().In(loc)
		}
	}
	if loc == nil {
		loc = time.Local
	}

	output := WebhookOutput{
		Date: DateOutput{
			Gregorian: date.Readable,
		},
		Location: LocationOutput{
			Latitude:  meta.Latitude,
			Longitude: meta.Longitude,
			Timezone:  meta.Timezone,
			Address:   data.Location,
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
		ServerTime: time.Now().UTC().Format(time.RFC3339),
	}

	// Add Hijri date
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

	// Calculate next prayer with full details
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
			output.NextPrayer = &WebhookNextPrayer{
				Name:         p.name,
				Time:         p.time,
				ISO:          prayerTime.UTC().Format(time.RFC3339),
				Timestamp:    prayerTime.Unix(),
				MinutesUntil: mins,
			}
			break
		}
	}

	// Add Qibla
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
