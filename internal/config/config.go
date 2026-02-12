// Package config provides configuration management for the pray CLI
package config

import (
	"os"
	"path/filepath"

	"github.com/AbdElrahmaN31/pray-cli/internal/location"
)

// Config represents the application configuration
type Config struct {
	// Location settings
	Location location.Location `yaml:"location"`

	// Calculation settings
	Method   int    `yaml:"method"`   // Calculation method ID (default: 5)
	Language string `yaml:"language"` // Language: "en" or "ar"

	// Display preferences
	Output OutputConfig `yaml:"output"`

	// Features
	Features FeaturesConfig `yaml:"features"`

	// Calendar settings
	Calendar CalendarConfig `yaml:"calendar"`

	// Jumu'ah settings
	Jumuah JumuahConfig `yaml:"jumuah"`

	// Ramadan settings
	Ramadan RamadanConfig `yaml:"ramadan"`

	// Iqama settings
	Iqama IqamaConfig `yaml:"iqama"`

	// Advanced settings
	CacheEnabled bool `yaml:"cache_enabled"`
	UpdateCheck  bool `yaml:"update_check"`
	APITimeout   int  `yaml:"api_timeout"` // Timeout in seconds
}

// OutputConfig contains display/output preferences
type OutputConfig struct {
	Format       string `yaml:"format"` // "table", "pretty", "json", "slack", "discord"
	ColorEnabled bool   `yaml:"color_enabled"`
	NoEmoji      bool   `yaml:"no_emoji"`
}

// FeaturesConfig contains feature toggle settings
type FeaturesConfig struct {
	Qibla         bool   `yaml:"qibla"`
	Dua           bool   `yaml:"dua"`
	Hijri         string `yaml:"hijri"` // "title", "desc", "both", "none"
	HijriHolidays bool   `yaml:"hijri_holidays"`
	TravelerMode  bool   `yaml:"traveler_mode"`
}

// CalendarConfig contains calendar generation settings
type CalendarConfig struct {
	Duration int    `yaml:"duration"` // Default event duration in minutes
	Months   int    `yaml:"months"`   // Number of months to generate
	Alarm    string `yaml:"alarm"`    // Comma-separated alarm offsets
	Events   string `yaml:"events"`   // Events to include ("all" or indices)
	Color    string `yaml:"color"`    // Calendar color
}

// JumuahConfig contains Friday prayer settings
type JumuahConfig struct {
	Enabled  bool `yaml:"enabled"`
	Duration int  `yaml:"duration"` // Duration in minutes
}

// RamadanConfig contains Ramadan mode settings
type RamadanConfig struct {
	Enabled          bool `yaml:"enabled"`
	IftarDuration    int  `yaml:"iftar_duration"`
	TaraweehDuration int  `yaml:"taraweeh_duration"`
	SuhoorDuration   int  `yaml:"suhoor_duration"`
}

// IqamaConfig contains Iqama settings
type IqamaConfig struct {
	Enabled bool   `yaml:"enabled"`
	Offsets string `yaml:"offsets"` // Comma-separated offsets for each prayer
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Location: location.Location{
			Source: "manual",
		},
		Method:   5, // Egyptian General Authority
		Language: "en",
		Output: OutputConfig{
			Format:       "table",
			ColorEnabled: true,
			NoEmoji:      false,
		},
		Features: FeaturesConfig{
			Qibla:         false,
			Dua:           false,
			Hijri:         "desc",
			HijriHolidays: false,
			TravelerMode:  false,
		},
		Calendar: CalendarConfig{
			Duration: 25,
			Months:   3,
			Alarm:    "5,10,15",
			Events:   "all",
			Color:    "#1e90ff",
		},
		Jumuah: JumuahConfig{
			Enabled:  false,
			Duration: 60,
		},
		Ramadan: RamadanConfig{
			Enabled:          false,
			IftarDuration:    30,
			TaraweehDuration: 60,
			SuhoorDuration:   30,
		},
		Iqama: IqamaConfig{
			Enabled: false,
			Offsets: "15,0,10,10,5,10,0",
		},
		CacheEnabled: true,
		UpdateCheck:  true,
		APITimeout:   30,
	}
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, ".config", "pray"), nil
	}
	return filepath.Join(configDir, "pray"), nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

// GetCacheDir returns the cache directory path
func GetCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		// Fallback to config directory
		configDir, err := GetConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(configDir, "cache"), nil
	}
	return filepath.Join(cacheDir, "pray"), nil
}

// IsConfigured checks if the config has a valid location set
func (c *Config) IsConfigured() bool {
	return c.Location.IsValid()
}

// Validate validates the configuration
func (c *Config) Validate() error {
	return ValidateConfig(c)
}
