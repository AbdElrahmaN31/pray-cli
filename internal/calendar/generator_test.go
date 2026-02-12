package calendar

import (
	"strings"
	"testing"
)

func TestNewCalendarParams(t *testing.T) {
	params := NewCalendarParams()

	if params == nil {
		t.Fatal("NewCalendarParams returned nil")
	}

	// Check defaults
	if params.Method != 5 {
		t.Errorf("Expected method 5, got %d", params.Method)
	}

	if params.Duration != 25 {
		t.Errorf("Expected duration 25, got %d", params.Duration)
	}

	if params.Months != 3 {
		t.Errorf("Expected months 3, got %d", params.Months)
	}

	if params.Language != "en" {
		t.Errorf("Expected language 'en', got '%s'", params.Language)
	}
}

func TestCalendarParamsBuilders(t *testing.T) {
	params := NewCalendarParams()

	// Test builder methods
	params.WithCoordinates(30.0, 31.0)
	if params.Latitude != 30.0 || params.Longitude != 31.0 {
		t.Error("WithCoordinates did not set values correctly")
	}

	params.WithAddress("Cairo, Egypt")
	if params.Address != "Cairo, Egypt" {
		t.Error("WithAddress did not set value correctly")
	}

	params.WithMethod(2)
	if params.Method != 2 {
		t.Error("WithMethod did not set value correctly")
	}

	params.WithDuration(30)
	if params.Duration != 30 {
		t.Error("WithDuration did not set value correctly")
	}

	params.WithMonths(6)
	if params.Months != 6 {
		t.Error("WithMonths did not set value correctly")
	}

	params.WithAlarm("5,10")
	if params.Alarm != "5,10" {
		t.Error("WithAlarm did not set value correctly")
	}

	params.WithLanguage("ar")
	if params.Language != "ar" {
		t.Error("WithLanguage did not set value correctly")
	}

	params.WithColor("#ff0000")
	if params.Color != "#ff0000" {
		t.Error("WithColor did not set value correctly")
	}

	params.WithJumuah(true, 60)
	if !params.Jumuah || params.JumuahDuration != 60 {
		t.Error("WithJumuah did not set values correctly")
	}

	params.WithRamadan(true)
	if !params.Ramadan {
		t.Error("WithRamadan did not set value correctly")
	}
}

func TestGenerateICSURL(t *testing.T) {
	tests := []struct {
		name     string
		params   *CalendarParams
		contains []string
	}{
		{
			name: "basic with coordinates",
			params: func() *CalendarParams {
				p := NewCalendarParams()
				p.WithCoordinates(30.0, 31.0)
				return p
			}(),
			contains: []string{"latitude=30", "longitude=31", "method=5"},
		},
		{
			name: "basic with address",
			params: func() *CalendarParams {
				p := NewCalendarParams()
				p.WithAddress("Cairo, Egypt")
				return p
			}(),
			contains: []string{"address=Cairo"},
		},
		{
			name: "with jumuah",
			params: func() *CalendarParams {
				p := NewCalendarParams()
				p.WithCoordinates(30.0, 31.0)
				p.WithJumuah(true, 60)
				return p
			}(),
			contains: []string{"jumuah=true", "jumuahDuration=60"},
		},
		{
			name: "with ramadan",
			params: func() *CalendarParams {
				p := NewCalendarParams()
				p.WithCoordinates(30.0, 31.0)
				p.WithRamadan(true)
				p.IftarDuration = 30
				return p
			}(),
			contains: []string{"ramadan=true", "iftarDuration=30"},
		},
		{
			name: "with arabic language",
			params: func() *CalendarParams {
				p := NewCalendarParams()
				p.WithCoordinates(30.0, 31.0)
				p.WithLanguage("ar")
				return p
			}(),
			contains: []string{"lang=ar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := GenerateICSURL(tt.params)

			if !strings.HasPrefix(url, BaseURL) {
				t.Errorf("URL should start with %s, got %s", BaseURL, url)
			}

			for _, contain := range tt.contains {
				if !strings.Contains(url, contain) {
					t.Errorf("URL should contain '%s', got %s", contain, url)
				}
			}
		})
	}
}

func TestGetDefaultFilename(t *testing.T) {
	tests := []struct {
		location string
		want     string
	}{
		{"Cairo, Egypt", "cairo-egypt.ics"},
		{"New York", "new-york.ics"},
		{"London/UK", "london-uk.ics"},
		{"", "prayer-times.ics"},
		{"City With  Spaces", "city-with-spaces.ics"},
	}

	for _, tt := range tests {
		t.Run(tt.location, func(t *testing.T) {
			got := GetDefaultFilename(tt.location)
			if got != tt.want {
				t.Errorf("GetDefaultFilename(%s) = %s, want %s", tt.location, got, tt.want)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Cairo, Egypt", "cairo-egypt"},
		{"Test/File", "test-file"},
		{"File:Name", "file-name"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"--Leading-Dashes--", "leading-dashes"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeFilename(%s) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}
