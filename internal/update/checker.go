// Package update provides version update checking functionality
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// GitHubReleasesURL is the API endpoint for GitHub releases
	GitHubReleasesURL = "https://api.github.com/repos/anashaat/pray-cli/releases/latest"

	// DefaultTimeout for update checks
	DefaultTimeout = 5 * time.Second
)

// ReleaseInfo contains information about a GitHub release
type ReleaseInfo struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	HTMLURL     string    `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
}

// Checker checks for new versions of the CLI
type Checker struct {
	currentVersion string
	httpClient     *http.Client
	timeout        time.Duration
}

// NewChecker creates a new update checker
func NewChecker(currentVersion string) *Checker {
	return &Checker{
		currentVersion: currentVersion,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		timeout: DefaultTimeout,
	}
}

// WithTimeout sets a custom timeout
func (c *Checker) WithTimeout(timeout time.Duration) *Checker {
	c.timeout = timeout
	c.httpClient.Timeout = timeout
	return c
}

// CheckResult contains the result of an update check
type CheckResult struct {
	UpdateAvailable bool
	CurrentVersion  string
	LatestVersion   string
	ReleaseURL      string
	ReleaseNotes    string
	PublishedAt     time.Time
}

// Check checks for a new version
func (c *Checker) Check(ctx context.Context) (*CheckResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", GitHubReleasesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "pray-cli/"+c.currentVersion)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	// Skip prereleases and drafts
	if release.Prerelease || release.Draft {
		return &CheckResult{
			UpdateAvailable: false,
			CurrentVersion:  c.currentVersion,
		}, nil
	}

	latestVersion := normalizeVersion(release.TagName)
	currentVersion := normalizeVersion(c.currentVersion)

	result := &CheckResult{
		CurrentVersion: c.currentVersion,
		LatestVersion:  release.TagName,
		ReleaseURL:     release.HTMLURL,
		ReleaseNotes:   truncateString(release.Body, 500),
		PublishedAt:    release.PublishedAt,
	}

	// Compare versions
	result.UpdateAvailable = isNewerVersion(currentVersion, latestVersion)

	return result, nil
}

// CheckAsync performs an update check in the background
func (c *Checker) CheckAsync(ctx context.Context) <-chan *CheckResult {
	resultChan := make(chan *CheckResult, 1)

	go func() {
		result, err := c.Check(ctx)
		if err != nil {
			// Silently fail - update checks shouldn't interrupt normal usage
			resultChan <- nil
		} else {
			resultChan <- result
		}
		close(resultChan)
	}()

	return resultChan
}

// normalizeVersion removes the 'v' prefix from version strings
func normalizeVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
}

// isNewerVersion compares two semantic versions
// Returns true if latest is newer than current
func isNewerVersion(current, latest string) bool {
	// Handle development versions
	if current == "dev" || current == "" {
		return false
	}

	currentParts := parseVersion(current)
	latestParts := parseVersion(latest)

	for i := 0; i < 3; i++ {
		var currentPart, latestPart int
		if i < len(currentParts) {
			currentPart = currentParts[i]
		}
		if i < len(latestParts) {
			latestPart = latestParts[i]
		}

		if latestPart > currentPart {
			return true
		}
		if latestPart < currentPart {
			return false
		}
	}

	return false
}

// parseVersion parses a version string into numeric parts
func parseVersion(version string) []int {
	// Remove any suffix after dash (e.g., "1.0.0-beta" -> "1.0.0")
	if idx := strings.Index(version, "-"); idx != -1 {
		version = version[:idx]
	}

	parts := strings.Split(version, ".")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		var num int
		fmt.Sscanf(part, "%d", &num)
		result = append(result, num)
	}

	return result
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatUpdateMessage formats a user-friendly update notification
func FormatUpdateMessage(result *CheckResult) string {
	if result == nil || !result.UpdateAvailable {
		return ""
	}

	return fmt.Sprintf(
		"\n📦 A new version of pray is available: %s → %s\n"+
			"   Run 'go install github.com/anashaat/pray-cli/cmd/pray@latest' to update\n"+
			"   Or visit: %s\n",
		result.CurrentVersion,
		result.LatestVersion,
		result.ReleaseURL,
	)
}
