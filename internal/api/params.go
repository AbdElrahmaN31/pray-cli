// Package api provides HTTP client for the prayer times API
package api

import (
	"fmt"
	"net/url"
	"time"
)

// PrayerTimesParams contains parameters for fetching prayer times
type PrayerTimesParams struct {
	// Location
	Latitude  float64
	Longitude float64
	Address   string

	// Date
	Date time.Time

	// Calculation method (0-23)
	Method int

	// School (0 = Shafi, 1 = Hanafi)
	School int

	// Timezone (e.g., "Africa/Cairo")
	Timezone string

	// Language (en or ar)
	Language string

	// Adjustments
	Adjustment int // Days adjustment (-30 to +30)

	// ISO8601 format for timings
	ISO8601 bool

	// Include Qibla direction
	Qibla bool
}

// NewPrayerTimesParams creates a new PrayerTimesParams with defaults
func NewPrayerTimesParams() *PrayerTimesParams {
	return &PrayerTimesParams{
		Date:     time.Now(),
		Method:   5, // Egyptian General Authority
		School:   0, // Shafi
		Language: "en",
	}
}

// GetDateString returns the date formatted for the API
func (p *PrayerTimesParams) GetDateString() string {
	if p.Date.IsZero() {
		p.Date = time.Now()
	}
	return p.Date.Format("02-01-2006") // DD-MM-YYYY
}

// ToQueryParams converts the parameters to URL query parameters
func (p *PrayerTimesParams) ToQueryParams() url.Values {
	query := url.Values{}

	// Location
	if p.Latitude != 0 || p.Longitude != 0 {
		query.Set("latitude", fmt.Sprintf("%f", p.Latitude))
		query.Set("longitude", fmt.Sprintf("%f", p.Longitude))
	}

	// Method
	query.Set("method", fmt.Sprintf("%d", p.Method))

	// School
	if p.School > 0 {
		query.Set("school", fmt.Sprintf("%d", p.School))
	}

	// Timezone
	if p.Timezone != "" {
		query.Set("timezonestring", p.Timezone)
	}

	// Adjustment
	if p.Adjustment != 0 {
		query.Set("adjustment", fmt.Sprintf("%d", p.Adjustment))
	}

	// ISO8601 format
	if p.ISO8601 {
		query.Set("iso8601", "true")
	}

	return query
}

// CalendarParams contains parameters for calendar generation
type CalendarParams struct {
	// Location
	Latitude  float64
	Longitude float64
	Address   string

	// Time range
	Year  int
	Month int

	// Calculation method
	Method int

	// Event settings
	Duration int    // Event duration in minutes
	Months   int    // Number of months to generate
	Alarm    string // Comma-separated alarm offsets
	Events   string // Events to include ("all" or indices)

	// Display settings
	Language string
	Color    string
	Hijri    string

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

// NewCalendarParams creates a new CalendarParams with defaults
func NewCalendarParams() *CalendarParams {
	now := time.Now()
	return &CalendarParams{
		Year:     now.Year(),
		Month:    int(now.Month()),
		Method:   5,
		Duration: 25,
		Months:   3,
		Alarm:    "5,10,15",
		Events:   "all",
		Language: "en",
		Color:    "#1e90ff",
		Hijri:    "desc",
	}
}

// ToQueryParams converts the parameters to URL query parameters
func (p *CalendarParams) ToQueryParams() url.Values {
	query := url.Values{}

	// Location
	if p.Latitude != 0 || p.Longitude != 0 {
		query.Set("latitude", fmt.Sprintf("%f", p.Latitude))
		query.Set("longitude", fmt.Sprintf("%f", p.Longitude))
	}

	// Method
	query.Set("method", fmt.Sprintf("%d", p.Method))

	return query
}

// WithCoordinates sets latitude and longitude
func (p *PrayerTimesParams) WithCoordinates(lat, lon float64) *PrayerTimesParams {
	p.Latitude = lat
	p.Longitude = lon
	return p
}

// WithAddress sets the address
func (p *PrayerTimesParams) WithAddress(address string) *PrayerTimesParams {
	p.Address = address
	return p
}

// WithMethod sets the calculation method
func (p *PrayerTimesParams) WithMethod(method int) *PrayerTimesParams {
	p.Method = method
	return p
}

// WithDate sets the date
func (p *PrayerTimesParams) WithDate(date time.Time) *PrayerTimesParams {
	p.Date = date
	return p
}

// WithTimezone sets the timezone
func (p *PrayerTimesParams) WithTimezone(tz string) *PrayerTimesParams {
	p.Timezone = tz
	return p
}

// CalendarParams builder methods

// WithCoordinates sets latitude and longitude
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

// WithEvents sets the events to include
func (p *CalendarParams) WithEvents(events string) *CalendarParams {
	p.Events = events
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

// WithTraveler enables traveler mode
func (p *CalendarParams) WithTraveler(enabled bool) *CalendarParams {
	p.Traveler = enabled
	return p
}
