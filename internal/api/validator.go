// Package api provides HTTP client for the prayer times API
package api

import (
	"fmt"
)

// ValidateParams validates the prayer times parameters
func ValidateParams(params *PrayerTimesParams) error {
	// Check for location
	if params.Address == "" && (params.Latitude == 0 && params.Longitude == 0) {
		return fmt.Errorf("location is required: provide either address or coordinates")
	}

	// Validate latitude
	if params.Latitude != 0 {
		if params.Latitude < -90 || params.Latitude > 90 {
			return fmt.Errorf("latitude must be between -90 and 90, got %f", params.Latitude)
		}
	}

	// Validate longitude
	if params.Longitude != 0 {
		if params.Longitude < -180 || params.Longitude > 180 {
			return fmt.Errorf("longitude must be between -180 and 180, got %f", params.Longitude)
		}
	}

	// Validate method
	if params.Method < 0 || params.Method > 23 {
		return fmt.Errorf("method must be between 0 and 23, got %d", params.Method)
	}

	// Validate school
	if params.School < 0 || params.School > 1 {
		return fmt.Errorf("school must be 0 (Shafi) or 1 (Hanafi), got %d", params.School)
	}

	// Validate adjustment
	if params.Adjustment < -30 || params.Adjustment > 30 {
		return fmt.Errorf("adjustment must be between -30 and 30, got %d", params.Adjustment)
	}

	return nil
}

// ValidateCalendarParams validates the calendar parameters
func ValidateCalendarParams(params *CalendarParams) error {
	// Check for location
	if params.Address == "" && (params.Latitude == 0 && params.Longitude == 0) {
		return fmt.Errorf("location is required: provide either address or coordinates")
	}

	// Validate latitude
	if params.Latitude != 0 {
		if params.Latitude < -90 || params.Latitude > 90 {
			return fmt.Errorf("latitude must be between -90 and 90, got %f", params.Latitude)
		}
	}

	// Validate longitude
	if params.Longitude != 0 {
		if params.Longitude < -180 || params.Longitude > 180 {
			return fmt.Errorf("longitude must be between -180 and 180, got %f", params.Longitude)
		}
	}

	// Validate method
	if params.Method < 0 || params.Method > 23 {
		return fmt.Errorf("method must be between 0 and 23, got %d", params.Method)
	}

	// Validate duration
	if params.Duration < 1 || params.Duration > 120 {
		return fmt.Errorf("duration must be between 1 and 120 minutes, got %d", params.Duration)
	}

	// Validate months
	if params.Months < 1 || params.Months > 12 {
		return fmt.Errorf("months must be between 1 and 12, got %d", params.Months)
	}

	// Validate year
	if params.Year < 2000 || params.Year > 2100 {
		return fmt.Errorf("year must be between 2000 and 2100, got %d", params.Year)
	}

	// Validate month
	if params.Month < 1 || params.Month > 12 {
		return fmt.Errorf("month must be between 1 and 12, got %d", params.Month)
	}

	return nil
}

// ValidateCoordinates validates latitude and longitude values
func ValidateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %f", lat)
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got %f", lon)
	}
	return nil
}

// ValidateMethod validates a calculation method ID
func ValidateMethod(method int) error {
	if method < 0 || method > 23 {
		return fmt.Errorf("method must be between 0 and 23, got %d", method)
	}
	return nil
}
