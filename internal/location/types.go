// Package location provides location detection and geocoding functionality
package location

import "time"

// Location represents a geographic location
type Location struct {
	Address     string    `yaml:"address" json:"address"`
	Latitude    float64   `yaml:"latitude" json:"latitude"`
	Longitude   float64   `yaml:"longitude" json:"longitude"`
	City        string    `yaml:"city,omitempty" json:"city,omitempty"`
	Country     string    `yaml:"country,omitempty" json:"country,omitempty"`
	CountryCode string    `yaml:"country_code,omitempty" json:"countryCode,omitempty"`
	Timezone    string    `yaml:"timezone" json:"timezone"`
	DetectedAt  time.Time `yaml:"detected_at,omitempty" json:"detectedAt,omitempty"`
	Source      string    `yaml:"source" json:"source"` // "ip", "manual", "gps"
}

// IPGeoResponse represents the response from ip-api.com
type IPGeoResponse struct {
	Status      string  `json:"status"`
	Message     string  `json:"message,omitempty"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Query       string  `json:"query"`
}

// IPInfoResponse represents the response from ipinfo.io
type IPInfoResponse struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"` // "lat,lon" format
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

// IPAPICoResponse represents the response from ipapi.co
type IPAPICoResponse struct {
	IP            string  `json:"ip"`
	City          string  `json:"city"`
	Region        string  `json:"region"`
	CountryName   string  `json:"country_name"`
	CountryCode   string  `json:"country_code"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Timezone      string  `json:"timezone"`
	UtcOffset     string  `json:"utc_offset"`
	ContinentCode string  `json:"continent_code"`
	InEU          bool    `json:"in_eu"`
	Error         bool    `json:"error,omitempty"`
	Reason        string  `json:"reason,omitempty"`
}

// Coordinates represents a simple lat/lon pair
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// IsValid checks if the location has valid coordinates
func (l *Location) IsValid() bool {
	return l.Latitude >= -90 && l.Latitude <= 90 &&
		l.Longitude >= -180 && l.Longitude <= 180 &&
		(l.Latitude != 0 || l.Longitude != 0)
}

// HasTimezone checks if the location has a timezone set
func (l *Location) HasTimezone() bool {
	return l.Timezone != ""
}

// GetDisplayAddress returns a human-readable address
func (l *Location) GetDisplayAddress() string {
	if l.Address != "" {
		return l.Address
	}
	if l.City != "" && l.Country != "" {
		return l.City + ", " + l.Country
	}
	if l.City != "" {
		return l.City
	}
	return ""
}
