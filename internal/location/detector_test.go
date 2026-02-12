package location

import (
	"context"
	"testing"
	"time"
)

func TestNewDetector(t *testing.T) {
	d := NewDetector()
	if d == nil {
		t.Error("NewDetector returned nil")
	}
	if d.timeout != DefaultDetectionTimeout {
		t.Errorf("Expected timeout %v, got %v", DefaultDetectionTimeout, d.timeout)
	}
}

func TestDetectorWithTimeout(t *testing.T) {
	d := NewDetector().WithTimeout(5 * time.Second)
	if d.timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", d.timeout)
	}
}

func TestParseLatLon(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLat float64
		wantLon float64
		wantErr bool
	}{
		{
			name:    "valid coordinates",
			input:   "30.0444,31.2357",
			wantLat: 30.0444,
			wantLon: 31.2357,
			wantErr: false,
		},
		{
			name:    "valid with spaces",
			input:   "30.0444, 31.2357",
			wantLat: 30.0444,
			wantLon: 31.2357,
			wantErr: false,
		},
		{
			name:    "negative coordinates",
			input:   "-33.8688,151.2093",
			wantLat: -33.8688,
			wantLon: 151.2093,
			wantErr: false,
		},
		{
			name:    "invalid format - no comma",
			input:   "30.0444 31.2357",
			wantErr: true,
		},
		{
			name:    "invalid latitude",
			input:   "abc,31.2357",
			wantErr: true,
		},
		{
			name:    "invalid longitude",
			input:   "30.0444,xyz",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lon, err := parseLatLon(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLatLon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if lat != tt.wantLat {
					t.Errorf("parseLatLon() lat = %v, want %v", lat, tt.wantLat)
				}
				if lon != tt.wantLon {
					t.Errorf("parseLatLon() lon = %v, want %v", lon, tt.wantLon)
				}
			}
		})
	}
}

func TestFormatAddress(t *testing.T) {
	tests := []struct {
		city    string
		country string
		want    string
	}{
		{"Cairo", "Egypt", "Cairo, Egypt"},
		{"London", "", "London"},
		{"", "USA", "USA"},
		{"", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.city+"-"+tt.country, func(t *testing.T) {
			got := formatAddress(tt.city, tt.country)
			if got != tt.want {
				t.Errorf("formatAddress(%q, %q) = %q, want %q", tt.city, tt.country, got, tt.want)
			}
		})
	}
}

func TestLocationIsValid(t *testing.T) {
	tests := []struct {
		name string
		loc  *Location
		want bool
	}{
		{
			name: "valid location",
			loc: &Location{
				Latitude:  30.0444,
				Longitude: 31.2357,
			},
			want: true,
		},
		{
			name: "zero coordinates",
			loc: &Location{
				Latitude:  0,
				Longitude: 0,
			},
			want: false,
		},
		{
			name: "invalid latitude - too high",
			loc: &Location{
				Latitude:  100,
				Longitude: 31.2357,
			},
			want: false,
		},
		{
			name: "invalid latitude - too low",
			loc: &Location{
				Latitude:  -100,
				Longitude: 31.2357,
			},
			want: false,
		},
		{
			name: "invalid longitude - too high",
			loc: &Location{
				Latitude:  30.0444,
				Longitude: 200,
			},
			want: false,
		},
		{
			name: "invalid longitude - too low",
			loc: &Location{
				Latitude:  30.0444,
				Longitude: -200,
			},
			want: false,
		},
		{
			name: "negative valid coordinates",
			loc: &Location{
				Latitude:  -33.8688,
				Longitude: 151.2093,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.loc.IsValid()
			if got != tt.want {
				t.Errorf("Location.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocationGetDisplayAddress(t *testing.T) {
	tests := []struct {
		name string
		loc  *Location
		want string
	}{
		{
			name: "with address",
			loc: &Location{
				Address: "Cairo, Egypt",
				City:    "Cairo",
				Country: "Egypt",
			},
			want: "Cairo, Egypt",
		},
		{
			name: "no address but city and country",
			loc: &Location{
				City:    "London",
				Country: "UK",
			},
			want: "London, UK",
		},
		{
			name: "only city",
			loc: &Location{
				City: "Paris",
			},
			want: "Paris",
		},
		{
			name: "empty",
			loc:  &Location{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.loc.GetDisplayAddress()
			if got != tt.want {
				t.Errorf("Location.GetDisplayAddress() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFromCoordinates(t *testing.T) {
	loc := FromCoordinates(30.0444, 31.2357)
	if loc.Latitude != 30.0444 {
		t.Errorf("Expected latitude 30.0444, got %f", loc.Latitude)
	}
	if loc.Longitude != 31.2357 {
		t.Errorf("Expected longitude 31.2357, got %f", loc.Longitude)
	}
	if loc.Source != "manual" {
		t.Errorf("Expected source 'manual', got %s", loc.Source)
	}
}

func TestFromAddress(t *testing.T) {
	loc := FromAddress("Cairo, Egypt")
	if loc.Address != "Cairo, Egypt" {
		t.Errorf("Expected address 'Cairo, Egypt', got %s", loc.Address)
	}
	if loc.Source != "manual" {
		t.Errorf("Expected source 'manual', got %s", loc.Source)
	}
}

func TestValidateLocation(t *testing.T) {
	tests := []struct {
		name    string
		loc     *Location
		wantErr bool
	}{
		{
			name: "valid location",
			loc: &Location{
				Latitude:  30.0444,
				Longitude: 31.2357,
			},
			wantErr: false,
		},
		{
			name:    "nil location",
			loc:     nil,
			wantErr: true,
		},
		{
			name: "invalid location",
			loc: &Location{
				Latitude:  0,
				Longitude: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLocation(tt.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Integration test - runs against live IP geolocation services
func TestDetectFromIPIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	d := NewDetector()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	loc, err := d.DetectFromIP(ctx)
	if err != nil {
		t.Fatalf("DetectFromIP failed: %v", err)
	}

	if loc == nil {
		t.Fatal("DetectFromIP returned nil location")
	}

	if !loc.IsValid() {
		t.Errorf("DetectFromIP returned invalid location: %+v", loc)
	}

	if loc.Source != "ip" {
		t.Errorf("Expected source 'ip', got %s", loc.Source)
	}

	if loc.Timezone == "" {
		t.Error("Timezone is empty")
	}

	t.Logf("Detected location: %s (%.4f, %.4f) TZ: %s",
		loc.GetDisplayAddress(), loc.Latitude, loc.Longitude, loc.Timezone)
}
