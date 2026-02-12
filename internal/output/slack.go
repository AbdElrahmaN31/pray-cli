// Package output provides output formatting for prayer times
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
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
		name string
		time string
	}{
		{"Fajr", cleanTime(timings.Fajr)},
		{"Sunrise", cleanTime(timings.Sunrise)},
		{"Dhuhr", cleanTime(timings.Dhuhr)},
		{"Asr", cleanTime(timings.Asr)},
		{"Maghrib", cleanTime(timings.Maghrib)},
		{"Isha", cleanTime(timings.Isha)},
	}

	// Find next prayer
	nextPrayer := ""
	for _, p := range prayers {
		prayerTime, err := parseTimeToday(p.time, now)
		if err != nil {
			continue
		}
		if now.Before(prayerTime) {
			nextPrayer = p.name
			break
		}
	}

	message := SlackMessage{
		Blocks: []SlackBlock{
			{
				Type: "header",
				Text: &SlackText{
					Type:  "plain_text",
					Text:  fmt.Sprintf("🕌 Prayer Times - %s", data.Location),
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
						if p.name == nextPrayer {
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
						Text: fmt.Sprintf("Method: %s", data.Method),
					},
				},
			},
		},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(message)
}
