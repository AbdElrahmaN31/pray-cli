package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/AbdElrahmaN31/pray-cli/internal/api"
	"github.com/AbdElrahmaN31/pray-cli/internal/cache"
	"github.com/AbdElrahmaN31/pray-cli/internal/config"
	"github.com/AbdElrahmaN31/pray-cli/internal/location"
)

// ResolvedLocation holds the result of location resolution
type ResolvedLocation struct {
	Latitude    float64
	Longitude   float64
	DisplayName string
	Timezone    string
	Detected    *location.Location // non-nil only when auto-detected
	IsAddress   bool               // true when using address string (no coords)
}

// resolveLocation resolves the user's location from flags or config.
// Priority: -- address > -auto > --lat/--lon > config.
// Returns nil with a helpful message if no location is available.
func resolveLocation() (*ResolvedLocation, error) {
	cfg := GetConfig()

	if address != "" {
		return &ResolvedLocation{
			DisplayName: address,
			IsAddress:   true,
		}, nil
	}

	if autoDetect {
		detector := location.NewDetector()
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		loc, err := detector.DetectFromIP(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to auto-detect location: %w", err)
		}
		return &ResolvedLocation{
			Latitude:    loc.Latitude,
			Longitude:   loc.Longitude,
			DisplayName: loc.GetDisplayAddress(),
			Timezone:    loc.Timezone,
			Detected:    loc,
		}, nil
	}

	if latitude != 0 || longitude != 0 {
		return &ResolvedLocation{
			Latitude:    latitude,
			Longitude:   longitude,
			DisplayName: fmt.Sprintf("%.4f, %.4f", latitude, longitude),
		}, nil
	}

	if cfg.Location.Address != "" && !cfg.Location.IsValid() {
		return &ResolvedLocation{
			DisplayName: cfg.Location.Address,
			IsAddress:   true,
		}, nil
	}

	if cfg.IsConfigured() {
		return &ResolvedLocation{
			Latitude:    cfg.Location.Latitude,
			Longitude:   cfg.Location.Longitude,
			DisplayName: cfg.Location.GetDisplayAddress(),
			Timezone:    cfg.Location.Timezone,
		}, nil
	}

	return nil, nil
}

// resolveMethod returns the calculation method ID from flags or config.
func resolveMethod() int {
	cfg := GetConfig()
	if method != 0 {
		return method
	}
	return cfg.Method
}

// createAPIClient creates an API client with caching support.
// When --no-cache is passed, the cache is bypassed.
func createAPIClient() (*api.CachedClient, error) {
	cfg := GetConfig()
	client := api.NewClient(api.WithTimeout(time.Duration(cfg.APITimeout) * time.Second))

	cacheDir, err := config.GetCacheDir()
	if err != nil {
		// Fall back to uncached client
		return api.NewCachedClient(client, api.WithBypassCache(true)), nil
	}

	c, err := cache.New(cacheDir, cache.WithEnabled(cfg.CacheEnabled))
	if err != nil {
		return api.NewCachedClient(client, api.WithBypassCache(true)), nil
	}

	return api.NewCachedClient(client,
		api.WithCache(c),
		api.WithBypassCache(ShouldBypassCache()),
	), nil
}

// printNoLocationMessage prints a helpful message when no location is configured.
func printNoLocationMessage() {
	fmt.Println("👋 Welcome! No location configured.")
	fmt.Println()
	fmt.Println("Set your location using one of these options:")
	fmt.Println("  pray config detect --save    Auto-detect from IP")
	fmt.Println("  pray --auto                  Auto-detect (one-time)")
	fmt.Println("  pray -a \"Cairo, Egypt\"       Specify a city")
	fmt.Println("  pray init                    Interactive setup")
}

// fetchPrayerTimes fetches prayer times using the resolved location and client.
func fetchPrayerTimes(ctx context.Context, client *api.CachedClient, loc *ResolvedLocation, methodID int, date time.Time) (*api.PrayerTimesResponse, error) {
	params := api.NewPrayerTimesParams().
		WithDate(date).
		WithMethod(methodID)

	// Set school if Hanafi
	if s := GetSchool(); s != 0 {
		params.School = s
	}

	if loc.IsAddress {
		params.WithAddress(loc.DisplayName)
		return client.GetPrayerTimesByAddress(ctx, params)
	}

	params.WithCoordinates(loc.Latitude, loc.Longitude)
	if loc.Timezone != "" {
		params.WithTimezone(loc.Timezone)
	}
	return client.GetPrayerTimes(ctx, params)
}
