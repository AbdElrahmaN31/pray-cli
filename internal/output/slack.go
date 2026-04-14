// Package output provides output formatting for prayer times
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/AbdElrahmaN31/pray-cli/pkg/prayer"
)

// SlackFormatter formats output as Slack Block Kit JSON
type SlackFormatter struct{}

// SlackMessage represents a Slack message with blocks
type SlackMessage struct {
	Blocks []SlackBlock `json:"blocks"`
}

// SlackBlock represents a Slack block
type SlackBlock struct {
	Type     string         `json:"type"`
	Text     *SlackText     `json:"text,omitempty"`
	Fields   []SlackText    `json:"fields,omitempty"`
	Elements []SlackElement `json:"elements,omitempty"`
}

// SlackText represents text content in Slack
type SlackText struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

// SlackElement represents an element in a context block
type SlackElement struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Format writes the prayer times as Slack Block Kit JSON
func (f *SlackFormatter) Format(w io.Writer, data *PrayerData) error {
	if data.Response == nil {
		return fmt.Errorf("no prayer times data")
	}

	resp := data.Response
	timings := resp.Data.Timings
	date := resp.Data.Date
	meta := resp.Data.Meta

	// Get current time for next prayer calculation
	now := time.Now()
	tz := meta.Timezone
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err == nil {
			now = time.Now().In(loc)
		}
	}

	prayers := []struct {
		key  string
		name string
		time string
	}{
		{keyFajr, t(data.Language, keyFajr), cleanTime(timings.Fajr)},
		{keySunrise, t(data.Language, keySunrise), cleanTime(timings.Sunrise)},
		{keyDhuhr, t(data.Language, keyDhuhr), cleanTime(timings.Dhuhr)},
		{keyAsr, t(data.Language, keyAsr), cleanTime(timings.Asr)},
		{keyMaghrib, t(data.Language, keyMaghrib), cleanTime(timings.Maghrib)},
		{keyIsha, t(data.Language, keyIsha), cleanTime(timings.Isha)},
	}

	// Find next prayer (compared by stable key, not localized name)
	nextPrayer := ""
	for _, p := range prayers {
		prayerTime, err := parseTimeToday(p.time, now)
		if err != nil {
			continue
		}
		if now.Before(prayerTime) {
			nextPrayer = p.key
			break
		}
	}

	message := SlackMessage{
		Blocks: []SlackBlock{
			{
				Type: "header",
				Text: &SlackText{
					Type:  "plain_text",
					Text:  fmt.Sprintf("🕌 %s", fmt.Sprintf(t(data.Language, "prayer_times_dash"), data.Location)),
					Emoji: true,
				},
			},
			{
				Type: "section",
				Text: &SlackText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("📅 *%s*", date.Readable),
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Fields: func() []SlackText {
					fields := make([]SlackText, 0)
					for _, p := range prayers {
						indicator := ""
						if p.key == nextPrayer {
							indicator = " ▶️"
						}
						fields = append(fields, SlackText{
							Type: "mrkdwn",
							Text: fmt.Sprintf("*%s:*\n%s%s", p.name, p.time, indicator),
						})
					}
					return fields
				}(),
			},
			{
				Type: "context",
				Elements: []SlackElement{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("%s: %s", t(data.Language, "method"), data.Method),
					},
				},
			},
		},
	}

	// Add Du'a section if enabled
	if data.ShowDua {
		dua := prayer.GetDailyDua(time.Now())
		if dua != nil {
			message.Blocks = append(message.Blocks,
				SlackBlock{Type: "divider"},
				SlackBlock{
					Type: "section",
					Text: &SlackText{
						Type: "mrkdwn",
						Text: fmt.Sprintf("📖 *%s*\n%s\n_%s_\n\"%s\"\n— %s",
							t(data.Language, "todays_dua"), dua.Arabic, dua.Transliteration, dua.Translation, dua.Reference),
					},
				},
			)
		}
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(message)
}
