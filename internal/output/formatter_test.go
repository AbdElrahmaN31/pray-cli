package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/AbdElrahmaN31/pray-cli/internal/api"
)

func TestGetFormatter(t *testing.T) {
	tests := []struct {
		format string
		want   string
	}{
		{"table", "*output.TableFormatter"},
		{"pretty", "*output.PrettyFormatter"},
		{"json", "*output.JSONFormatter"},
		{"slack", "*output.SlackFormatter"},
		{"discord", "*output.DiscordFormatter"},
		{"webhook", "*output.WebhookFormatter"},
		{"unknown", "*output.TableFormatter"}, // Default
		{"", "*output.TableFormatter"},        // Empty default
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			formatter := GetFormatter(tt.format)
			if formatter == nil {
				t.Error("GetFormatter returned nil")
			}
		})
	}
}

func TestFormatTypes(t *testing.T) {
	types := FormatTypes()

	expected := []string{"table", "pretty", "json", "slack", "discord", "webhook"}

	if len(types) != len(expected) {
		t.Errorf("FormatTypes() returned %d types, want %d", len(types), len(expected))
	}

	for _, exp := range expected {
		found := false
		for _, typ := range types {
			if typ == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("FormatTypes() missing '%s'", exp)
		}
	}
}

func createTestPrayerData() *PrayerData {
	return &PrayerData{
		Response: &api.PrayerTimesResponse{
			Code:   200,
			Status: "OK",
			Data: api.Data{
				Timings: api.Timings{
					Fajr:     "05:15",
					Sunrise:  "06:44",
					Dhuhr:    "12:09",
					Asr:      "15:12",
					Maghrib:  "17:34",
					Isha:     "18:54",
					Midnight: "00:09",
				},
				Date: api.Date{
					Readable: "04 Feb 2026",
					Hijri: api.HijriDate{
						Day:  "16",
						Year: "1447",
						Month: api.HijriMonthInfo{
							Number: 8,
							En:     "Sha'ban",
							Ar:     "شعبان",
						},
					},
				},
				Meta: api.Meta{
					Latitude:  30.0,
					Longitude: 31.0,
					Timezone:  "Africa/Cairo",
					Method: api.Method{
						ID:   5,
						Name: "Egyptian General Authority of Survey",
					},
				},
			},
		},
		Location:    "Cairo, Egypt",
		Method:      "Egyptian General Authority of Survey",
		ShowHijri:   true,
		HijriFormat: "desc",
		Language:    "en",
		NoColor:     true,
	}
}

func TestJSONFormatter(t *testing.T) {
	data := createTestPrayerData()

	var buf bytes.Buffer
	formatter := &JSONFormatter{}

	err := formatter.Format(&buf, data)
	if err != nil {
		t.Fatalf("JSONFormatter.Format() error = %v", err)
	}

	output := buf.String()

	// Check that output contains expected JSON fields
	expectedFields := []string{
		`"Fajr"`,
		`"Dhuhr"`,
		`"Isha"`,
		`"gregorian"`,
		`"hijri"`,
		`"latitude"`,
		`"longitude"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("JSON output missing field '%s'", field)
		}
	}
}

func TestSlackFormatter(t *testing.T) {
	data := createTestPrayerData()

	var buf bytes.Buffer
	formatter := &SlackFormatter{}

	err := formatter.Format(&buf, data)
	if err != nil {
		t.Fatalf("SlackFormatter.Format() error = %v", err)
	}

	output := buf.String()

	// Check Slack block format
	if !strings.Contains(output, `"blocks"`) {
		t.Error("Slack output missing 'blocks' array")
	}

	if !strings.Contains(output, `"type": "header"`) {
		t.Error("Slack output missing header block")
	}
}

func TestDiscordFormatter(t *testing.T) {
	data := createTestPrayerData()

	var buf bytes.Buffer
	formatter := &DiscordFormatter{}

	err := formatter.Format(&buf, data)
	if err != nil {
		t.Fatalf("DiscordFormatter.Format() error = %v", err)
	}

	output := buf.String()

	// Check Discord embed format
	if !strings.Contains(output, `"embeds"`) {
		t.Error("Discord output missing 'embeds' array")
	}

	if !strings.Contains(output, `"fields"`) {
		t.Error("Discord output missing 'fields' array")
	}

	if !strings.Contains(output, `"color"`) {
		t.Error("Discord output missing 'color' field")
	}
}

func TestWebhookFormatter(t *testing.T) {
	data := createTestPrayerData()

	var buf bytes.Buffer
	formatter := &WebhookFormatter{}

	err := formatter.Format(&buf, data)
	if err != nil {
		t.Fatalf("WebhookFormatter.Format() error = %v", err)
	}

	output := buf.String()

	// Check webhook format
	if !strings.Contains(output, `"serverTime"`) {
		t.Error("Webhook output missing 'serverTime' field")
	}

	if !strings.Contains(output, `"timings"`) {
		t.Error("Webhook output missing 'timings' field")
	}
}

func TestFormatWithNilResponse(t *testing.T) {
	data := &PrayerData{
		Response: nil,
	}

	formatters := []struct {
		name      string
		formatter Formatter
	}{
		{"table", &TableFormatter{}},
		{"pretty", &PrettyFormatter{}},
		{"json", &JSONFormatter{}},
		{"slack", &SlackFormatter{}},
		{"discord", &DiscordFormatter{}},
		{"webhook", &WebhookFormatter{}},
	}

	for _, f := range formatters {
		t.Run(f.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := f.formatter.Format(&buf, data)
			if err == nil {
				t.Errorf("%s formatter should return error for nil response", f.name)
			}
		})
	}
}

func TestCleanTime(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"05:15", "05:15"},
		{"05:15 (EET)", "05:15"},
		{"12:00 (UTC)", "12:00"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := cleanTime(tt.input)
			if got != tt.want {
				t.Errorf("cleanTime(%s) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetCompassDirection(t *testing.T) {
	tests := []struct {
		degrees float64
		want    string
	}{
		{0, "N"},
		{45, "NE"},
		{90, "E"},
		{135, "SE"},
		{180, "S"},
		{225, "SW"},
		{270, "W"},
		{315, "NW"},
		{360, "N"},
		{-90, "W"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := getCompassDirection(tt.degrees)
			if got != tt.want {
				t.Errorf("getCompassDirection(%f) = %s, want %s", tt.degrees, got, tt.want)
			}
		})
	}
}

func TestFormatMinutes(t *testing.T) {
	tests := []struct {
		mins int
		want string
	}{
		{30, "30 min"},
		{60, "1h"},
		{90, "1h 30m"},
		{120, "2h"},
		{150, "2h 30m"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatMinutes(tt.mins)
			if got != tt.want {
				t.Errorf("formatMinutes(%d) = %s, want %s", tt.mins, got, tt.want)
			}
		})
	}
}
