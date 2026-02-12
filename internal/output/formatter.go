// Package output provides output formatting for prayer times
package output

import (
	"io"

	"github.com/anashaat/pray-cli/internal/api"
)

// Formatter is the interface for output formatters
type Formatter interface {
	Format(w io.Writer, data *PrayerData) error
}

// PrayerData contains all the data needed for formatting
type PrayerData struct {
	Response    *api.PrayerTimesResponse
	Location    string
	Method      string
	NextPrayer  *api.NextPrayer
	Qibla       *api.QiblaData
	ShowQibla   bool
	ShowDua     bool
	ShowHijri   bool
	HijriFormat string // "title", "desc", "both", "none"
	Language    string
	NoColor     bool
}

// GetFormatter returns the appropriate formatter for the given format
func GetFormatter(format string) Formatter {
	switch format {
	case "pretty":
		return &PrettyFormatter{}
	case "json":
		return &JSONFormatter{}
	case "slack":
		return &SlackFormatter{}
	case "discord":
		return &DiscordFormatter{}
	case "webhook":
		return &WebhookFormatter{}
	default:
		return &TableFormatter{}
	}
}

// FormatTypes returns all available format types
func FormatTypes() []string {
	return []string{"table", "pretty", "json", "slack", "discord", "webhook"}
}
