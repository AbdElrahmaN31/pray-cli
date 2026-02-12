// Package location provides location detection and geocoding functionality
package location

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// Primary IP geolocation service
	IPAPIEndpoint = "http://ip-api.com/json/"

	// Secondary fallback
	IPInfoEndpoint = "https://ipinfo.io/json"

	// Tertiary fallback
	IPAPICoEndpoint = "https://ipapi.co/json/"

	// Default timeout for location detection
	DefaultDetectionTimeout = 10 * time.Second
)

// Detector handles location detection from various sources
type Detector struct {
	httpClient *http.Client
	timeout    time.Duration
}

// NewDetector creates a new location detector
func NewDetector() *Detector {
	return &Detector{
		httpClient: &http.Client{
			Timeout: DefaultDetectionTimeout,
		},
		timeout: DefaultDetectionTimeout,
	}
}

// WithTimeout sets a custom timeout
func (d *Detector) WithTimeout(timeout time.Duration) *Detector {
	d.timeout = timeout
	d.httpClient.Timeout = timeout
	return d
}

// DetectFromIP detects location using IP geolocation services
// It tries multiple services with fallback
func (d *Detector) DetectFromIP(ctx context.Context) (*Location, error) {
	// Try primary service (ip-api.com)
	loc, err := d.detectFromIPAPI(ctx)
	if err == nil && loc.IsValid() {
		loc.Source = "ip"
		loc.DetectedAt = time.Now()
		return loc, nil
	}

	// Try secondary service (ipinfo.io)
	loc, err = d.detectFromIPInfo(ctx)
	if err == nil && loc.IsValid() {
		loc.Source = "ip"
		loc.DetectedAt = time.Now()
		return loc, nil
	}

	// Try tertiary service (ipapi.co)
	loc, err = d.detectFromIPAPICo(ctx)
	if err == nil && loc.IsValid() {
		loc.Source = "ip"
		loc.DetectedAt = time.Now()
		return loc, nil
	}

	return nil, fmt.Errorf("failed to detect location from IP: all services failed")
}

// detectFromIPAPI uses ip-api.com for geolocation
func (d *Detector) detectFromIPAPI(ctx context.Context) (*Location, error) {
	resp, err := d.doRequest(ctx, IPAPIEndpoint)
	if err != nil {
		return nil, err
	}

	var result IPGeoResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse ip-api.com response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("ip-api.com error: %s", result.Message)
	}

	return &Location{
		Latitude:    result.Lat,
		Longitude:   result.Lon,
		City:        result.City,
		Country:     result.Country,
		CountryCode: result.CountryCode,
		Timezone:    result.Timezone,
		Address:     formatAddress(result.City, result.Country),
	}, nil
}

// detectFromIPInfo uses ipinfo.io for geolocation
func (d *Detector) detectFromIPInfo(ctx context.Context) (*Location, error) {
	resp, err := d.doRequest(ctx, IPInfoEndpoint)
	if err != nil {
		return nil, err
	}

	var result IPInfoResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse ipinfo.io response: %w", err)
	}

	// Parse "lat,lon" format
	lat, lon, err := parseLatLon(result.Loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coordinates from ipinfo.io: %w", err)
	}

	return &Location{
		Latitude:    lat,
		Longitude:   lon,
		City:        result.City,
		Country:     result.Country,
		CountryCode: result.Country, // ipinfo.io uses country code in "country" field
		Timezone:    result.Timezone,
		Address:     formatAddress(result.City, result.Region),
	}, nil
}

// detectFromIPAPICo uses ipapi.co for geolocation
func (d *Detector) detectFromIPAPICo(ctx context.Context) (*Location, error) {
	resp, err := d.doRequest(ctx, IPAPICoEndpoint)
	if err != nil {
		return nil, err
	}

	var result IPAPICoResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse ipapi.co response: %w", err)
	}

	if result.Error {
		return nil, fmt.Errorf("ipapi.co error: %s", result.Reason)
	}

	return &Location{
		Latitude:    result.Latitude,
		Longitude:   result.Longitude,
		City:        result.City,
		Country:     result.CountryName,
		CountryCode: result.CountryCode,
		Timezone:    result.Timezone,
		Address:     formatAddress(result.City, result.CountryName),
	}, nil
}

// doRequest performs an HTTP GET request
func (d *Detector) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "pray-cli/1.0.0")
	req.Header.Set("Accept", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

// parseLatLon parses a "lat,lon" string into float64 values
func parseLatLon(loc string) (float64, float64, error) {
	parts := strings.Split(loc, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid location format: %s", loc)
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %w", err)
	}

	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %w", err)
	}

	return lat, lon, nil
}

// formatAddress creates a human-readable address from city and country
func formatAddress(city, country string) string {
	if city == "" && country == "" {
		return ""
	}
	if city == "" {
		return country
	}
	if country == "" {
		return city
	}
	return fmt.Sprintf("%s, %s", city, country)
}

// ValidateLocation validates a location's coordinates
func ValidateLocation(loc *Location) error {
	if loc == nil {
		return fmt.Errorf("location is nil")
	}

	if !loc.IsValid() {
		return fmt.Errorf("invalid coordinates: lat=%f, lon=%f", loc.Latitude, loc.Longitude)
	}

	return nil
}

// FromCoordinates creates a Location from coordinates
func FromCoordinates(lat, lon float64) *Location {
	return &Location{
		Latitude:  lat,
		Longitude: lon,
		Source:    "manual",
	}
}

// FromAddress creates a Location from an address string
// Note: This doesn't geocode the address, it just stores it
// The actual coordinates should be obtained from the API
func FromAddress(address string) *Location {
	return &Location{
		Address: address,
		Source:  "manual",
	}
}
