package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.baseURL != DefaultBaseURL {
		t.Errorf("Expected baseURL %s, got %s", DefaultBaseURL, client.baseURL)
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

func TestGetPrayerTimes_Success(t *testing.T) {
	mockResp := `{"code":200,"status":"OK","data":{"timings":{"Fajr":"05:15","Sunrise":"06:30","Dhuhr":"12:00","Asr":"15:30","Sunset":"17:45","Maghrib":"17:45","Isha":"19:00","Imsak":"05:05","Midnight":"23:30","Firstthird":"21:00","Lastthird":"01:00"},"date":{"readable":"15 Mar 2026","timestamp":"1773792000","gregorian":{"date":"15-03-2026","format":"DD-MM-YYYY","day":"15","weekday":{"en":"Sunday"},"month":{"number":3,"en":"March"},"year":"2026","designation":{"abbreviated":"AD","expanded":"Anno Domini"}},"hijri":{"date":"15-08-1447","format":"DD-MM-YYYY","day":"15","weekday":{"en":"Al Ahad","ar":"الاحد"},"month":{"number":8,"en":"Sha'ban","ar":"شعبان"},"year":"1447","designation":{"abbreviated":"AH","expanded":"Anno Hegirae"},"holidays":[]}},"meta":{"latitude":30.0444,"longitude":31.2357,"timezone":"Africa/Cairo","method":{"id":5,"name":"Egyptian General Authority of Survey","params":{"Fajr":19.5,"Isha":17.5}},"latitudeAdjustmentMethod":"ANGLE_BASED","midnightMode":"STANDARD","school":"STANDARD","offset":{"Imsak":0,"Fajr":0,"Sunrise":0,"Dhuhr":0,"Asr":0,"Maghrib":0,"Sunset":0,"Isha":0,"Midnight":0}}}}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mockResp)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithMaxRetries(0))
	params := NewPrayerTimesParams().WithCoordinates(30.0444, 31.2357).WithMethod(5)

	resp, err := client.GetPrayerTimes(context.Background(), params)
	if err != nil {
		t.Fatalf("GetPrayerTimes() error = %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if resp.Data.Timings.Fajr != "05:15" {
		t.Errorf("expected Fajr 05:15, got %s", resp.Data.Timings.Fajr)
	}
}

func TestGetPrayerTimes_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "internal error")
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithMaxRetries(0))
	params := NewPrayerTimesParams().WithCoordinates(30.0444, 31.2357).WithMethod(5)

	_, err := client.GetPrayerTimes(context.Background(), params)
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention status 500: %v", err)
	}
}

func TestGetPrayerTimes_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":400,"status":"Bad Request","data":{}}`)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithMaxRetries(0))
	params := NewPrayerTimesParams().WithCoordinates(30.0444, 31.2357).WithMethod(5)

	_, err := client.GetPrayerTimes(context.Background(), params)
	if err == nil {
		t.Fatal("expected error for API code 400")
	}
	if !strings.Contains(err.Error(), "API error") {
		t.Errorf("error should mention API error: %v", err)
	}
}

func TestRetry_EventualSuccess(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":200,"status":"OK","data":{"timings":{"Fajr":"05:15","Sunrise":"06:30","Dhuhr":"12:00","Asr":"15:30","Sunset":"17:45","Maghrib":"17:45","Isha":"19:00","Imsak":"05:05","Midnight":"23:30","Firstthird":"21:00","Lastthird":"01:00"},"date":{"readable":"15 Mar 2026","timestamp":"1773792000","gregorian":{"date":"15-03-2026","format":"DD-MM-YYYY","day":"15","weekday":{"en":"Sunday"},"month":{"number":3,"en":"March"},"year":"2026","designation":{"abbreviated":"AD","expanded":"Anno Domini"}},"hijri":{"date":"15-08-1447","format":"DD-MM-YYYY","day":"15","weekday":{"en":"Al Ahad","ar":"الاحد"},"month":{"number":8,"en":"Sha'ban","ar":"شعبان"},"year":"1447","designation":{"abbreviated":"AH","expanded":"Anno Hegirae"},"holidays":[]}},"meta":{"latitude":30.0444,"longitude":31.2357,"timezone":"Africa/Cairo","method":{"id":5,"name":"Egyptian","params":{"Fajr":19.5,"Isha":17.5}},"latitudeAdjustmentMethod":"ANGLE_BASED","midnightMode":"STANDARD","school":"STANDARD","offset":{"Imsak":0,"Fajr":0,"Sunrise":0,"Dhuhr":0,"Asr":0,"Maghrib":0,"Sunset":0,"Isha":0,"Midnight":0}}}}`)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithMaxRetries(3))
	params := NewPrayerTimesParams().WithCoordinates(30.0444, 31.2357).WithMethod(5)

	resp, err := client.GetPrayerTimes(context.Background(), params)
	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if atomic.LoadInt32(&attempts) < 3 {
		t.Errorf("expected at least 3 attempts, got %d", attempts)
	}
}

func TestRetry_AllFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithMaxRetries(2))
	params := NewPrayerTimesParams().WithCoordinates(30.0444, 31.2357).WithMethod(5)

	_, err := client.GetPrayerTimes(context.Background(), params)
	if err == nil {
		t.Fatal("expected error after all retries fail")
	}
	if !strings.Contains(err.Error(), "attempts") {
		t.Errorf("error should mention attempts: %v", err)
	}
}

func TestRetry_ContextCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithMaxRetries(3), WithTimeout(5*time.Second))
	params := NewPrayerTimesParams().WithCoordinates(30.0444, 31.2357).WithMethod(5)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.GetPrayerTimes(ctx, params)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestGetQibla_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":200,"status":"OK","data":{"latitude":30.0444,"longitude":31.2357,"direction":136.17}}`)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL), WithMaxRetries(0))
	resp, err := client.GetQibla(context.Background(), 30.0444, 31.2357)
	if err != nil {
		t.Fatalf("GetQibla() error = %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("expected code 200, got %d", resp.Code)
	}
	if resp.Data.Direction < 136 || resp.Data.Direction > 137 {
		t.Errorf("expected direction ~136.17, got %f", resp.Data.Direction)
	}
}

func TestBuildICSURL_RamadanParams(t *testing.T) {
	params := &CalendarParams{
		Address:          "Cairo, Egypt",
		Method:           5,
		Duration:         25,
		Months:           3,
		Ramadan:          true,
		IftarDuration:    30,
		TaraweehDuration: 60,
		SuhoorDuration:   30,
	}

	url := BuildICSURL(params)

	// Verify correct parameter names per PrayCalendar API
	expectedParts := []string{
		"ramadanMode=true",
		"traweehDuration=60",
		"iftarDuration=30",
		"suhoorDuration=30",
	}
	notExpected := []string{
		"ramadan=true",
		"taraweehDuration=",
	}

	for _, part := range expectedParts {
		if !containsString(url, part) {
			t.Errorf("URL missing expected param: %s\nURL: %s", part, url)
		}
	}
	for _, part := range notExpected {
		if containsString(url, part) {
			t.Errorf("URL should NOT contain old param: %s\nURL: %s", part, url)
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
