// Package update provides version update checking functionality
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const (
	// DefaultReleasesURL is the default API endpoint for GitHub releases
	DefaultReleasesURL = "https://api.github.com/repos/AbdElrahmaN31/pray-cli/releases/latest"

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
	releasesURL    string
	httpClient     *http.Client
	timeout        time.Duration
}

// NewChecker creates a new update checker
func NewChecker(currentVersion string) *Checker {
	return &Checker{
		currentVersion: currentVersion,
		releasesURL:    DefaultReleasesURL,
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

// WithReleasesURL overrides the GitHub releases API endpoint (useful for testing)
func (c *Checker) WithReleasesURL(url string) *Checker {
	c.releasesURL = url
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
	req, err := http.NewRequestWithContext(ctx, "GET", c.releasesURL, nil)
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

	result := &CheckResult{
		CurrentVersion: c.currentVersion,
		LatestVersion:  release.TagName,
		ReleaseURL:     release.HTMLURL,
		ReleaseNotes:   truncateString(release.Body, 500),
		PublishedAt:    release.PublishedAt,
	}

	result.UpdateAvailable = isNewerVersion(c.currentVersion, release.TagName)

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

// canonicalVersion ensures a version string has the "v" prefix required by
// golang.org/x/mod/semver. It also trims whitespace.
func canonicalVersion(version string) string {
	v := strings.TrimSpace(version)
	if v == "" {
		return ""
	}
	if v[0] != 'v' {
		v = "v" + v
	}
	return v
}

// isNewerVersion returns true if latest is a higher semver than current.
// Pre-release versions are compared correctly (e.g. 1.0.0-beta < 1.0.0).
func isNewerVersion(current, latest string) bool {
	if current == "dev" || current == "" {
		return false
	}

	cv := canonicalVersion(current)
	lv := canonicalVersion(latest)

	if !semver.IsValid(cv) || !semver.IsValid(lv) {
		return false
	}

	return semver.Compare(cv, lv) < 0
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatUpdateMessage formats a user-friendly update notification with
// install-method-aware upgrade instructions.
func FormatUpdateMessage(result *CheckResult) string {
	if result == nil || !result.UpdateAvailable {
		return ""
	}

	instruction := detectUpgradeInstruction()

	return fmt.Sprintf(
		"\n📦 A new version of pray is available: %s → %s\n"+
			"   %s\n"+
			"   Or visit: %s\n",
		result.CurrentVersion,
		result.LatestVersion,
		instruction,
		result.ReleaseURL,
	)
}

// detectUpgradeInstruction returns an upgrade command appropriate for how the
// binary was installed. It checks the executable path for hints.
func detectUpgradeInstruction() string {
	exe, err := os.Executable()
	if err != nil {
		return "Run 'go install github.com/AbdElrahmaN31/pray-cli/cmd/pray@latest' to update"
	}
	exe = strings.ToLower(exe)

	switch {
	case strings.Contains(exe, "homebrew") || strings.Contains(exe, "linuxbrew"):
		return "Run 'brew upgrade pray' to update"
	case strings.Contains(exe, "scoop"):
		return "Run 'scoop update pray' to update"
	case strings.Contains(exe, filepath.Join("go", "bin")):
		return "Run 'go install github.com/AbdElrahmaN31/pray-cli/cmd/pray@latest' to update"
	default:
		return "Run 'curl -sSL https://raw.githubusercontent.com/AbdElrahmaN31/pray-cli/main/install.sh | sh' to update"
	}
}
