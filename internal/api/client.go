// Package api provides HTTP client for the prayer times API
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the base URL for the prayer times API
	DefaultBaseURL = "https://pray.ahmedelywa.com"

	// AlAdhanBaseURL is the alternative API (aladhan.com)
	AlAdhanBaseURL = "https://api.aladhan.com/v1"

	// DefaultTimeout is the default HTTP timeout
	DefaultTimeout = 30 * time.Second

	// DefaultMaxRetries is the default number of retries
	DefaultMaxRetries = 3

	// UserAgent identifies the CLI client
	UserAgent = "pray-cli/1.0.0"
)

// Client is the HTTP client for the prayer times API
type Client struct {
	httpClient *http.Client
	baseURL    string
	timeout    time.Duration
	maxRetries int
	userAgent  string
}

// ClientOption configures the Client
type ClientOption func(*Client)

// WithTimeout sets the HTTP timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
		c.httpClient.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries
func WithMaxRetries(retries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = retries
	}
}

// WithBaseURL sets the base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new API client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:    AlAdhanBaseURL, // Using AlAdhan API as it's more reliable
		timeout:    DefaultTimeout,
		maxRetries: DefaultMaxRetries,
		userAgent:  UserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetPrayerTimes fetches prayer times for a specific date and location
func (c *Client) GetPrayerTimes(ctx context.Context, params *PrayerTimesParams) (*PrayerTimesResponse, error) {
	endpoint := fmt.Sprintf("%s/timings/%s", c.baseURL, params.GetDateString())

	// Build query parameters
	query := params.ToQueryParams()
	fullURL := fmt.Sprintf("%s?%s", endpoint, query.Encode())

	resp, err := c.doRequestWithRetry(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prayer times: %w", err)
	}

	var result PrayerTimesResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API error: %s (code: %d)", result.Status, result.Code)
	}

	return &result, nil
}

// GetPrayerTimesByAddress fetches prayer times using an address
func (c *Client) GetPrayerTimesByAddress(ctx context.Context, params *PrayerTimesParams) (*PrayerTimesResponse, error) {
	if params.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	endpoint := fmt.Sprintf("%s/timingsByAddress/%s", c.baseURL, params.GetDateString())

	// Build query parameters
	query := params.ToQueryParams()
	query.Set("address", params.Address)
	fullURL := fmt.Sprintf("%s?%s", endpoint, query.Encode())

	resp, err := c.doRequestWithRetry(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prayer times: %w", err)
	}

	var result PrayerTimesResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API error: %s (code: %d)", result.Status, result.Code)
	}

	return &result, nil
}

// GetQibla fetches the Qibla direction for a location
func (c *Client) GetQibla(ctx context.Context, latitude, longitude float64) (*QiblaResponse, error) {
	endpoint := fmt.Sprintf("%s/qibla/%f/%f", c.baseURL, latitude, longitude)

	resp, err := c.doRequestWithRetry(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Qibla direction: %w", err)
	}

	var result QiblaResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API error: %s (code: %d)", result.Status, result.Code)
	}

	return &result, nil
}

// GetCalendarMonth fetches prayer times for an entire month
func (c *Client) GetCalendarMonth(ctx context.Context, params *CalendarParams) ([]PrayerTimesResponse, error) {
	endpoint := fmt.Sprintf("%s/calendar/%d/%d", c.baseURL, params.Year, params.Month)

	// Build query parameters
	query := params.ToQueryParams()
	fullURL := fmt.Sprintf("%s?%s", endpoint, query.Encode())

	resp, err := c.doRequestWithRetry(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calendar: %w", err)
	}

	var result struct {
		Code   int                   `json:"code"`
		Status string                `json:"status"`
		Data   []PrayerTimesResponse `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}

// doRequestWithRetry performs an HTTP request with retry logic
func (c *Client) doRequestWithRetry(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, err := c.doRequest(ctx, method, url, body)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on context cancellation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries+1, lastErr)
}

// doRequest performs a single HTTP request
func (c *Client) doRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// DownloadICS downloads an ICS calendar file
func (c *Client) DownloadICS(ctx context.Context, icsURL string) ([]byte, error) {
	resp, err := c.doRequestWithRetry(ctx, "GET", icsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download ICS file: %w", err)
	}
	return resp, nil
}

// BuildICSURL builds a URL for the ICS calendar endpoint
func BuildICSURL(params *CalendarParams) string {
	baseURL := DefaultBaseURL
	query := url.Values{}

	// Location
	if params.Address != "" {
		query.Set("address", params.Address)
	} else {
		query.Set("latitude", fmt.Sprintf("%f", params.Latitude))
		query.Set("longitude", fmt.Sprintf("%f", params.Longitude))
	}

	// Method
	if params.Method > 0 {
		query.Set("method", fmt.Sprintf("%d", params.Method))
	}

	// Duration
	if params.Duration > 0 {
		query.Set("duration", fmt.Sprintf("%d", params.Duration))
	}

	// Months
	if params.Months > 0 {
		query.Set("months", fmt.Sprintf("%d", params.Months))
	}

	// Alarms
	if params.Alarm != "" {
		query.Set("alarm", params.Alarm)
	}

	// Events
	if params.Events != "" {
		query.Set("events", params.Events)
	}

	// Language
	if params.Language != "" {
		query.Set("lang", params.Language)
	}

	// Color
	if params.Color != "" {
		query.Set("color", params.Color)
	}

	// Hijri
	if params.Hijri != "" {
		query.Set("hijri", params.Hijri)
	}

	// Special features
	if params.Jumuah {
		query.Set("jumuah", "true")
	}
	if params.JumuahDuration > 0 {
		query.Set("jumuahDuration", fmt.Sprintf("%d", params.JumuahDuration))
	}
	if params.Qibla {
		query.Set("qibla", "true")
	}
	if params.Dua {
		query.Set("dua", "true")
	}
	if params.Traveler {
		query.Set("traveler", "true")
	}
	if params.Ramadan {
		query.Set("ramadan", "true")
	}
	if params.IftarDuration > 0 {
		query.Set("iftarDuration", fmt.Sprintf("%d", params.IftarDuration))
	}
	if params.TaraweehDuration > 0 {
		query.Set("taraweehDuration", fmt.Sprintf("%d", params.TaraweehDuration))
	}
	if params.SuhoorDuration > 0 {
		query.Set("suhoorDuration", fmt.Sprintf("%d", params.SuhoorDuration))
	}
	if params.HijriHolidays {
		query.Set("hijriHolidays", "true")
	}
	if params.Iqama != "" {
		query.Set("iqama", params.Iqama)
	}

	return fmt.Sprintf("%s/api/prayer-times.ics?%s", baseURL, query.Encode())
}
