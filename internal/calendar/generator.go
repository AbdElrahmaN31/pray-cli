// Package calendar provides calendar generation and ICS file handling
package calendar

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	// BaseURL for the prayer times calendar API
	BaseURL = "https://pray.ahmedelywa.com"
)

// CalendarParams contains parameters for generating a calendar
type CalendarParams struct {
	// Location
	Latitude  float64
	Longitude float64
	Address   string

	// Calendar settings
	Method   int
	Duration int    // Event duration in minutes
	Months   int    // Number of months to generate
	Alarm    string // Comma-separated alarm offsets
	Events   string // Events to include

	// Display settings
	Language string
	Color    string
	Hijri    string // "title", "desc", "both", "none"

	// Special features
	Jumuah           bool
	JumuahDuration   int
	Qibla            bool
	Dua              bool
	Traveler         bool
	Ramadan          bool
	IftarDuration    int
	TaraweehDuration int
	SuhoorDuration   int
	HijriHolidays    bool
	Iqama            string
}

// NewCalendarParams creates default calendar parameters
func NewCalendarParams() *CalendarParams {
	return &CalendarParams{
		Method:   5, // Egyptian
		Duration: 25,
		Months:   3,
		Alarm:    "5,10,15",
		Events:   "all",
		Language: "en",
		Color:    "#1e90ff",
		Hijri:    "desc",
	}
}

// GenerateICSURL generates the URL for downloading an ICS calendar file
func GenerateICSURL(params *CalendarParams) string {
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
		query.Set("ramadan", "true")
		if params.IftarDuration > 0 {
			query.Set("iftarDuration", fmt.Sprintf("%d", params.IftarDuration))
		}
		if params.TaraweehDuration > 0 {
			query.Set("taraweehDuration", fmt.Sprintf("%d", params.TaraweehDuration))
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

// WithCoordinates sets the coordinates
func (p *CalendarParams) WithCoordinates(lat, lon float64) *CalendarParams {
	p.Latitude = lat
	p.Longitude = lon
	return p
}

// WithAddress sets the address
func (p *CalendarParams) WithAddress(address string) *CalendarParams {
	p.Address = address
	return p
}

// WithMethod sets the calculation method
func (p *CalendarParams) WithMethod(method int) *CalendarParams {
	p.Method = method
	return p
}

// WithDuration sets the event duration
func (p *CalendarParams) WithDuration(duration int) *CalendarParams {
	p.Duration = duration
	return p
}

// WithMonths sets the number of months
func (p *CalendarParams) WithMonths(months int) *CalendarParams {
	p.Months = months
	return p
}

// WithAlarm sets the alarm offsets
func (p *CalendarParams) WithAlarm(alarm string) *CalendarParams {
	p.Alarm = alarm
	return p
}

// WithLanguage sets the language
func (p *CalendarParams) WithLanguage(lang string) *CalendarParams {
	p.Language = lang
	return p
}

// WithColor sets the calendar color
func (p *CalendarParams) WithColor(color string) *CalendarParams {
	p.Color = color
	return p
}

// WithJumuah enables Jumu'ah prayer
func (p *CalendarParams) WithJumuah(enabled bool, duration int) *CalendarParams {
	p.Jumuah = enabled
	p.JumuahDuration = duration
	return p
}

// WithRamadan enables Ramadan mode
func (p *CalendarParams) WithRamadan(enabled bool) *CalendarParams {
	p.Ramadan = enabled
	return p
}
