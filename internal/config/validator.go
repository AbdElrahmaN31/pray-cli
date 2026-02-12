// Package config provides configuration management for the pray CLI
package config

import (
	"fmt"
	"slices"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateConfig validates the entire configuration
func ValidateConfig(cfg *Config) error {
	// Validate method
	if !ValidMethodID(cfg.Method) {
		return ValidationError{
			Field:   "method",
			Message: fmt.Sprintf("invalid calculation method ID: %d", cfg.Method),
		}
	}

	// Validate language
	if !slices.Contains(DefaultLanguages, cfg.Language) {
		return ValidationError{
			Field:   "language",
			Message: fmt.Sprintf("invalid language: %s (must be 'en' or 'ar')", cfg.Language),
		}
	}

	// Validate output format
	if !slices.Contains(DefaultOutputFormats, cfg.Output.Format) {
		return ValidationError{
			Field:   "output.format",
			Message: fmt.Sprintf("invalid output format: %s", cfg.Output.Format),
		}
	}

	// Validate Hijri display option
	validHijriOptions := []string{"title", "desc", "both", "none"}
	if !slices.Contains(validHijriOptions, cfg.Features.Hijri) {
		return ValidationError{
			Field:   "features.hijri",
			Message: fmt.Sprintf("invalid hijri option: %s (must be title, desc, both, or none)", cfg.Features.Hijri),
		}
	}

	// Validate calendar settings
	if cfg.Calendar.Duration < 1 || cfg.Calendar.Duration > 120 {
		return ValidationError{
			Field:   "calendar.duration",
			Message: "duration must be between 1 and 120 minutes",
		}
	}

	if cfg.Calendar.Months < 1 || cfg.Calendar.Months > 12 {
		return ValidationError{
			Field:   "calendar.months",
			Message: "months must be between 1 and 12",
		}
	}

	// Validate API timeout
	if cfg.APITimeout < 5 || cfg.APITimeout > 120 {
		return ValidationError{
			Field:   "api_timeout",
			Message: "API timeout must be between 5 and 120 seconds",
		}
	}

	// Validate location if set
	if cfg.Location.Latitude != 0 || cfg.Location.Longitude != 0 {
		if cfg.Location.Latitude < -90 || cfg.Location.Latitude > 90 {
			return ValidationError{
				Field:   "location.latitude",
				Message: "latitude must be between -90 and 90",
			}
		}
		if cfg.Location.Longitude < -180 || cfg.Location.Longitude > 180 {
			return ValidationError{
				Field:   "location.longitude",
				Message: "longitude must be between -180 and 180",
			}
		}
	}

	return nil
}

// ValidateLatitude validates a latitude value
func ValidateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %f", lat)
	}
	return nil
}

// ValidateLongitude validates a longitude value
func ValidateLongitude(lon float64) error {
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got %f", lon)
	}
	return nil
}

// ValidateCoordinates validates both latitude and longitude
func ValidateCoordinates(lat, lon float64) error {
	if err := ValidateLatitude(lat); err != nil {
		return err
	}
	return ValidateLongitude(lon)
}
