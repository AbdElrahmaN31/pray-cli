package api

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient returned nil")
	}
	if client.baseURL != AlAdhanBaseURL {
		t.Errorf("Expected baseURL %s, got %s", AlAdhanBaseURL, client.baseURL)
	}
	if client.timeout != DefaultTimeout {
		t.Errorf("Expected timeout %v, got %v", DefaultTimeout, client.timeout)
	}
	if client.maxRetries != DefaultMaxRetries {
		t.Errorf("Expected maxRetries %d, got %d", DefaultMaxRetries, client.maxRetries)
	}
}

func TestClientWithOptions(t *testing.T) {
	customTimeout := 60 * time.Second
	customRetries := 5
	customURL := "https://custom.api.com"

	client := NewClient(
		WithTimeout(customTimeout),
		WithMaxRetries(customRetries),
		WithBaseURL(customURL),
	)

	if client.timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.timeout)
	}
	if client.maxRetries != customRetries {
		t.Errorf("Expected maxRetries %d, got %d", customRetries, client.maxRetries)
	}
	if client.baseURL != customURL {
		t.Errorf("Expected baseURL %s, got %s", customURL, client.baseURL)
	}
}

func TestPrayerTimesParams(t *testing.T) {
	params := NewPrayerTimesParams()

	if params.Method != 5 {
		t.Errorf("Expected default method 5, got %d", params.Method)
	}
	if params.Language != "en" {
		t.Errorf("Expected default language 'en', got %s", params.Language)
	}

	// Test builder methods
	params.WithCoordinates(30.0444, 31.2357).WithMethod(3)

	if params.Latitude != 30.0444 {
		t.Errorf("Expected latitude 30.0444, got %f", params.Latitude)
	}
	if params.Longitude != 31.2357 {
		t.Errorf("Expected longitude 31.2357, got %f", params.Longitude)
	}
	if params.Method != 3 {
		t.Errorf("Expected method 3, got %d", params.Method)
	}
}

func TestCalendarParams(t *testing.T) {
	params := NewCalendarParams()

	if params.Method != 5 {
		t.Errorf("Expected default method 5, got %d", params.Method)
	}
	if params.Duration != 25 {
		t.Errorf("Expected default duration 25, got %d", params.Duration)
	}
	if params.Months != 3 {
		t.Errorf("Expected default months 3, got %d", params.Months)
	}
	if params.Color != "#1e90ff" {
		t.Errorf("Expected default color '#1e90ff', got %s", params.Color)
	}
}

func TestValidateParams(t *testing.T) {
	tests := []struct {
		name    string
		params  *PrayerTimesParams
		wantErr bool
	}{
		{
			name: "valid coordinates",
			params: &PrayerTimesParams{
				Latitude:  30.0444,
				Longitude: 31.2357,
				Method:    5,
			},
			wantErr: false,
		},
		{
			name: "valid address",
			params: &PrayerTimesParams{
				Address: "Cairo, Egypt",
				Method:  5,
			},
			wantErr: false,
		},
		{
			name: "missing location",
			params: &PrayerTimesParams{
				Method: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid latitude",
			params: &PrayerTimesParams{
				Latitude:  100, // > 90
				Longitude: 31.2357,
				Method:    5,
			},
			wantErr: true,
		},
		{
			name: "invalid longitude",
			params: &PrayerTimesParams{
				Latitude:  30.0444,
				Longitude: 200, // > 180
				Method:    5,
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			params: &PrayerTimesParams{
				Latitude:  30.0444,
				Longitude: 31.2357,
				Method:    99, // > 23
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCalendarParams(t *testing.T) {
	tests := []struct {
		name    string
		params  *CalendarParams
		wantErr bool
	}{
		{
			name: "valid params",
			params: &CalendarParams{
				Latitude:  30.0444,
				Longitude: 31.2357,
				Method:    5,
				Duration:  25,
				Months:    3,
				Year:      2026,
				Month:     2,
			},
			wantErr: false,
		},
		{
			name: "invalid duration",
			params: &CalendarParams{
				Latitude:  30.0444,
				Longitude: 31.2357,
				Method:    5,
				Duration:  200, // > 120
				Months:    3,
				Year:      2026,
				Month:     2,
			},
			wantErr: true,
		},
		{
			name: "invalid months",
			params: &CalendarParams{
				Latitude:  30.0444,
				Longitude: 31.2357,
				Method:    5,
				Duration:  25,
				Months:    15, // > 12
				Year:      2026,
				Month:     2,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCalendarParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCalendarParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildICSURL(t *testing.T) {
	params := &CalendarParams{
		Address:  "Cairo, Egypt",
		Method:   5,
		Duration: 25,
		Months:   3,
		Language: "en",
		Color:    "#1e90ff",
	}

	url := BuildICSURL(params)

	if url == "" {
		t.Error("BuildICSURL returned empty string")
	}

	// Check that URL contains expected parts
	expectedParts := []string{
		"prayer-times.ics",
		"address=",
		"method=5",
		"duration=25",
		"months=3",
	}

	for _, part := range expectedParts {
		if !containsString(url, part) {
			t.Errorf("URL missing expected part: %s", part)
		}
	}
}

// Integration test - only runs if INTEGRATION_TEST env is set
func TestGetPrayerTimesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := NewClient()
	params := NewPrayerTimesParams().
		WithCoordinates(30.0444, 31.2357).
		WithMethod(5)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.GetPrayerTimes(ctx, params)
	if err != nil {
		t.Fatalf("GetPrayerTimes failed: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("Expected code 200, got %d", resp.Code)
	}

	if resp.Data.Timings.Fajr == "" {
		t.Error("Fajr time is empty")
	}
	if resp.Data.Timings.Dhuhr == "" {
		t.Error("Dhuhr time is empty")
	}
	if resp.Data.Timings.Maghrib == "" {
		t.Error("Maghrib time is empty")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
