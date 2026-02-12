package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	// Check default values
	if cfg.Method != 5 {
		t.Errorf("Expected method 5, got %d", cfg.Method)
	}

	if cfg.Language != "en" {
		t.Errorf("Expected language 'en', got '%s'", cfg.Language)
	}

	if cfg.Output.Format != "table" {
		t.Errorf("Expected output format 'table', got '%s'", cfg.Output.Format)
	}

	if cfg.Calendar.Duration != 25 {
		t.Errorf("Expected calendar duration 25, got %d", cfg.Calendar.Duration)
	}

	if cfg.Calendar.Months != 3 {
		t.Errorf("Expected calendar months 3, got %d", cfg.Calendar.Months)
	}

	if cfg.APITimeout != 30 {
		t.Errorf("Expected API timeout 30, got %d", cfg.APITimeout)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr bool
	}{
		{
			name:    "valid config",
			modify:  func(c *Config) {},
			wantErr: false,
		},
		{
			name:    "invalid method",
			modify:  func(c *Config) { c.Method = 100 },
			wantErr: true,
		},
		{
			name:    "invalid language",
			modify:  func(c *Config) { c.Language = "invalid" },
			wantErr: true,
		},
		{
			name:    "invalid output format",
			modify:  func(c *Config) { c.Output.Format = "invalid" },
			wantErr: true,
		},
		{
			name:    "invalid hijri option",
			modify:  func(c *Config) { c.Features.Hijri = "invalid" },
			wantErr: true,
		},
		{
			name:    "invalid calendar duration - too low",
			modify:  func(c *Config) { c.Calendar.Duration = 0 },
			wantErr: true,
		},
		{
			name:    "invalid calendar duration - too high",
			modify:  func(c *Config) { c.Calendar.Duration = 200 },
			wantErr: true,
		},
		{
			name:    "invalid calendar months",
			modify:  func(c *Config) { c.Calendar.Months = 15 },
			wantErr: true,
		},
		{
			name:    "invalid API timeout - too low",
			modify:  func(c *Config) { c.APITimeout = 2 },
			wantErr: true,
		},
		{
			name:    "invalid latitude",
			modify:  func(c *Config) { c.Location.Latitude = 100 },
			wantErr: true,
		},
		{
			name:    "invalid longitude",
			modify:  func(c *Config) { c.Location.Longitude = 200 },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			tt.modify(cfg)

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "pray-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testPath := filepath.Join(tmpDir, "config.yaml")

	// Create a config with custom values
	cfg := DefaultConfig()
	cfg.Method = 2
	cfg.Language = "ar"
	cfg.Location.Address = "Test City"
	cfg.Location.Latitude = 30.0
	cfg.Location.Longitude = 31.0

	// Save the config
	err = cfg.SaveToFile(testPath)
	if err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	// Load the config
	loadedCfg, err := LoadFromFile(testPath)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}

	// Check loaded values
	if loadedCfg.Method != cfg.Method {
		t.Errorf("Method = %d, want %d", loadedCfg.Method, cfg.Method)
	}
	if loadedCfg.Language != cfg.Language {
		t.Errorf("Language = %s, want %s", loadedCfg.Language, cfg.Language)
	}
	if loadedCfg.Location.Address != cfg.Location.Address {
		t.Errorf("Location.Address = %s, want %s", loadedCfg.Location.Address, cfg.Location.Address)
	}
}

func TestConfigIsConfigured(t *testing.T) {
	tests := []struct {
		name       string
		modify     func(*Config)
		configured bool
	}{
		{
			name:       "not configured - default",
			modify:     func(c *Config) {},
			configured: false,
		},
		{
			name: "configured with coordinates",
			modify: func(c *Config) {
				c.Location.Latitude = 30.0
				c.Location.Longitude = 31.0
			},
			configured: true,
		},
		{
			name: "configured with address only",
			modify: func(c *Config) {
				c.Location.Address = "Cairo"
			},
			configured: false, // Address without coords is not valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			tt.modify(cfg)

			if cfg.IsConfigured() != tt.configured {
				t.Errorf("IsConfigured() = %v, want %v", cfg.IsConfigured(), tt.configured)
			}
		})
	}
}

func TestGetMethodByID(t *testing.T) {
	// Test valid method
	method := GetMethodByID(5)
	if method == nil {
		t.Error("GetMethodByID(5) returned nil")
	} else if method.ID != 5 {
		t.Errorf("GetMethodByID(5).ID = %d, want 5", method.ID)
	}

	// Test invalid method
	method = GetMethodByID(100)
	if method != nil {
		t.Error("GetMethodByID(100) should return nil")
	}
}

func TestGetMethodName(t *testing.T) {
	name := GetMethodName(5)
	if name == "" || name == "Unknown" {
		t.Errorf("GetMethodName(5) = '%s', want a valid name", name)
	}

	name = GetMethodName(100)
	if name != "Unknown" {
		t.Errorf("GetMethodName(100) = '%s', want 'Unknown'", name)
	}
}

func TestValidMethodID(t *testing.T) {
	if !ValidMethodID(5) {
		t.Error("ValidMethodID(5) should return true")
	}

	if ValidMethodID(100) {
		t.Error("ValidMethodID(100) should return false")
	}
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		lat     float64
		lon     float64
		wantErr bool
	}{
		{30.0, 31.0, false},    // Valid
		{-90.0, -180.0, false}, // Boundary valid
		{90.0, 180.0, false},   // Boundary valid
		{91.0, 0.0, true},      // Invalid latitude
		{-91.0, 0.0, true},     // Invalid latitude
		{0.0, 181.0, true},     // Invalid longitude
		{0.0, -181.0, true},    // Invalid longitude
	}

	for _, tt := range tests {
		err := ValidateCoordinates(tt.lat, tt.lon)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateCoordinates(%f, %f) error = %v, wantErr %v",
				tt.lat, tt.lon, err, tt.wantErr)
		}
	}
}
