// Package calendar provides calendar generation and ICS file handling
package calendar

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// DefaultDownloadTimeout for downloading ICS files
	DefaultDownloadTimeout = 60 * time.Second
)

// Downloader handles downloading ICS files
type Downloader struct {
	httpClient *http.Client
	timeout    time.Duration
}

// NewDownloader creates a new ICS file downloader
func NewDownloader() *Downloader {
	return &Downloader{
		httpClient: &http.Client{
			Timeout: DefaultDownloadTimeout,
		},
		timeout: DefaultDownloadTimeout,
	}
}

// WithTimeout sets a custom timeout
func (d *Downloader) WithTimeout(timeout time.Duration) *Downloader {
	d.timeout = timeout
	d.httpClient.Timeout = timeout
	return d
}

// Download downloads an ICS file from the given URL
func (d *Downloader) Download(ctx context.Context, icsURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", icsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "pray-cli/1.0.0")
	req.Header.Set("Accept", "text/calendar,application/ics")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download calendar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download calendar: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// DownloadToFile downloads an ICS file and saves it to disk
func (d *Downloader) DownloadToFile(ctx context.Context, icsURL, filePath string) error {
	data, err := d.Download(ctx, icsURL)
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Write file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetDefaultFilename returns the default filename for the ICS file
func GetDefaultFilename(location string) string {
	// Sanitize location for filename
	safe := sanitizeFilename(location)
	if safe == "" {
		safe = "prayer-times"
	}
	return fmt.Sprintf("%s.ics", safe)
}

// sanitizeFilename removes or replaces invalid filename characters
func sanitizeFilename(name string) string {
	// Simple sanitization - replace common invalid characters
	replacer := []struct {
		old string
		new string
	}{
		{" ", "-"},
		{",", ""},
		{".", ""},
		{"/", "-"},
		{"\\", "-"},
		{":", "-"},
		{"*", ""},
		{"?", ""},
		{"\"", ""},
		{"<", ""},
		{">", ""},
		{"|", ""},
	}

	result := name
	for _, r := range replacer {
		result = strings.ReplaceAll(result, r.old, r.new)
	}

	// Convert to lowercase
	result = strings.ToLower(result)

	// Remove multiple consecutive dashes
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}

	// Trim leading/trailing dashes
	result = strings.Trim(result, "-")

	return result
}
