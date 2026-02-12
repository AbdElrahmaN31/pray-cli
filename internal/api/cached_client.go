// Package api provides HTTP client for the prayer times API
package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anashaat/pray-cli/internal/cache"
)

// CachedClient wraps Client with caching support
type CachedClient struct {
	*Client
	cache  *cache.Cache
	bypass bool
}

// CachedClientOption configures the CachedClient
type CachedClientOption func(*CachedClient)

// WithCache sets the cache instance
func WithCache(c *cache.Cache) CachedClientOption {
	return func(cc *CachedClient) {
		cc.cache = c
	}
}

// WithBypassCache sets whether to bypass the cache
func WithBypassCache(bypass bool) CachedClientOption {
	return func(cc *CachedClient) {
		cc.bypass = bypass
	}
}

// NewCachedClient creates a new CachedClient
func NewCachedClient(client *Client, opts ...CachedClientOption) *CachedClient {
	cc := &CachedClient{
		Client: client,
		bypass: false,
	}

	for _, opt := range opts {
		opt(cc)
	}

	return cc
}

// GetPrayerTimes fetches prayer times with caching support
func (cc *CachedClient) GetPrayerTimes(ctx context.Context, params *PrayerTimesParams) (*PrayerTimesResponse, error) {
	if cc.cache == nil || cc.bypass || !cc.cache.IsEnabled() {
		return cc.Client.GetPrayerTimes(ctx, params)
	}

	// Generate cache key
	key := cache.GenerateKey(
		"times",
		params.Latitude,
		params.Longitude,
		params.GetDateString(),
		params.Method,
	)

	// Try to get from cache
	if data, found := cc.cache.Get(key); found {
		var result PrayerTimesResponse
		if err := json.Unmarshal(data, &result); err == nil {
			return &result, nil
		}
	}

	// Fetch from API
	result, err := cc.Client.GetPrayerTimes(ctx, params)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(result); err == nil {
		cc.cache.Set(key, data)
	}

	return result, nil
}

// GetPrayerTimesByAddress fetches prayer times by address with caching support
func (cc *CachedClient) GetPrayerTimesByAddress(ctx context.Context, params *PrayerTimesParams) (*PrayerTimesResponse, error) {
	if cc.cache == nil || cc.bypass || !cc.cache.IsEnabled() {
		return cc.Client.GetPrayerTimesByAddress(ctx, params)
	}

	// Generate cache key
	key := cache.GenerateKey(
		"addr",
		params.Address,
		params.GetDateString(),
		params.Method,
	)

	// Try to get from cache
	if data, found := cc.cache.Get(key); found {
		var result PrayerTimesResponse
		if err := json.Unmarshal(data, &result); err == nil {
			return &result, nil
		}
	}

	// Fetch from API
	result, err := cc.Client.GetPrayerTimesByAddress(ctx, params)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(result); err == nil {
		cc.cache.Set(key, data)
	}

	return result, nil
}

// GetQibla fetches the Qibla direction with caching support
func (cc *CachedClient) GetQibla(ctx context.Context, latitude, longitude float64) (*QiblaResponse, error) {
	if cc.cache == nil || cc.bypass || !cc.cache.IsEnabled() {
		return cc.Client.GetQibla(ctx, latitude, longitude)
	}

	// Generate cache key (Qibla doesn't change, so use longer TTL key)
	key := cache.GenerateKey(
		"qibla",
		fmt.Sprintf("%.4f", latitude),
		fmt.Sprintf("%.4f", longitude),
	)

	// Try to get from cache
	if data, found := cc.cache.Get(key); found {
		var result QiblaResponse
		if err := json.Unmarshal(data, &result); err == nil {
			return &result, nil
		}
	}

	// Fetch from API
	result, err := cc.Client.GetQibla(ctx, latitude, longitude)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(result); err == nil {
		cc.cache.Set(key, data)
	}

	return result, nil
}

// SetBypass sets whether to bypass the cache
func (cc *CachedClient) SetBypass(bypass bool) {
	cc.bypass = bypass
}

// ClearCache clears the cache
func (cc *CachedClient) ClearCache() error {
	if cc.cache != nil {
		return cc.cache.Clear()
	}
	return nil
}
