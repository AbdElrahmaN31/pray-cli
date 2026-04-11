// Package calendar provides calendar generation and ICS file handling
package calendar

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AbdElrahmaN31/pray-cli/internal/api"
)

const (
	// BaseURL for the prayer times calendar API
	BaseURL = "https://pray.ahmedelywa.com"
)

// GenerateICSURL generates the URL for downloading an ICS calendar file
func GenerateICSURL(params *api.CalendarParams) string {
	query := url.Values{}

	// Location
	if params.Address != "" {
		query.Set("address", params.Address)
	} else if params.Latitude != 0 || params.Longitude != 0 {
		query.Set("latitude", fmt.Sprintf("%f", params.Latitude))
		query.Set("longitude", fmt.Sprintf("%f", params.Longitude))
	}

	// Method
	if params.Method > 0 {
		query.Set("method", fmt.Sprintf("%d", params.Method))
	}

	// Duration
	if params.Duration > 0 {
		query.Set("duration", fmt.Sprintf("%d", params.Duration))
	}

	// Months
	if params.Months > 0 {
		query.Set("months", fmt.Sprintf("%d", params.Months))
	}

	// Alarms
	if params.Alarm != "" {
		query.Set("alarm", params.Alarm)
	}

	// Events
	if params.Events != "" && params.Events != "all" {
		query.Set("events", params.Events)
	}

	// Language
	if params.Language != "" && params.Language != "en" {
		query.Set("lang", params.Language)
	}

	// Color
	if params.Color != "" {
		// Remove # from hex color if present
		color := strings.TrimPrefix(params.Color, "#")
		query.Set("color", color)
	}

	// Hijri
	if params.Hijri != "" && params.Hijri != "none" {
		query.Set("hijri", params.Hijri)
	}

	// Special features
	if params.Jumuah {
		query.Set("jumuah", "true")
		if params.JumuahDuration > 0 {
			query.Set("jumuahDuration", fmt.Sprintf("%d", params.JumuahDuration))
		}
	}

	if params.Qibla {
		query.Set("qibla", "true")
	}

	if params.Dua {
		query.Set("dua", "true")
	}

	if params.Traveler {
		query.Set("traveler", "true")
	}

	if params.Ramadan {
		query.Set("ramadanMode", "true")
		if params.IftarDuration > 0 {
			query.Set("iftarDuration", fmt.Sprintf("%d", params.IftarDuration))
		}
		if params.TaraweehDuration > 0 {
			query.Set("traweehDuration", fmt.Sprintf("%d", params.TaraweehDuration))
		}
		if params.SuhoorDuration > 0 {
			query.Set("suhoorDuration", fmt.Sprintf("%d", params.SuhoorDuration))
		}
	}

	if params.HijriHolidays {
		query.Set("hijriHolidays", "true")
	}

	if params.Iqama != "" {
		query.Set("iqama", params.Iqama)
	}

	return fmt.Sprintf("%s/api/prayer-times.ics?%s", BaseURL, query.Encode())
}
