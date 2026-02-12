// Package calendar provides calendar generation and ICS file handling
package calendar

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		result = replaceAll(result, r.old, r.new)
	}

	// Convert to lowercase
	result = toLowerCase(result)

	// Remove multiple consecutive dashes
	for contains(result, "--") {
		result = replaceAll(result, "--", "-")
	}

	// Trim leading/trailing dashes
	result = trim(result, "-")

	return result
}

func replaceAll(s, old, new string) string {
	for i := 0; i <= len(s)-len(old); i++ {
		if len(s) >= i+len(old) && s[i:i+len(old)] == old {
			s = s[:i] + new + s[i+len(old):]
			if new != "" {
				i += len(new) - 1
			} else {
				i--
			}
		}
	}
	return s
}

func toLowerCase(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func trim(s, cutset string) string {
	for len(s) > 0 {
		found := false
		for i := 0; i < len(cutset); i++ {
			if s[0] == cutset[i] {
				s = s[1:]
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	for len(s) > 0 {
		found := false
		for i := 0; i < len(cutset); i++ {
			if s[len(s)-1] == cutset[i] {
				s = s[:len(s)-1]
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return s
}
