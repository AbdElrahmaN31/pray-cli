// Package output provides output formatting for prayer times
package output

import (
	"fmt"
	"strings"

	"github.com/AbdElrahmaN31/pray-cli/internal/api"
)

// prayerKeys are the stable identifiers for prayer entries.
const (
	keyFajr     = "fajr"
	keySunrise  = "sunrise"
	keyDhuhr    = "dhuhr"
	keyAsr      = "asr"
	keyMaghrib  = "maghrib"
	keyIsha     = "isha"
	keyMidnight = "midnight"
)

var translations = map[string]map[string]string{
	"en": {
		keyFajr:             "Fajr",
		keySunrise:          "Sunrise",
		keyDhuhr:            "Dhuhr",
		keyAsr:              "Asr",
		keyMaghrib:          "Maghrib",
		keyIsha:             "Isha",
		keyMidnight:         "Midnight",
		"prayer":            "Prayer",
		"time":              "Time",
		"status":            "Status",
		"iqama":             "Iqama",
		"passed":            "Passed",
		"next":              "Next",
		"next_prayer_in":    "Next prayer in",
		"method":            "Method",
		"qibla":             "Qibla",
		"qibla_direction":   "Qibla Direction",
		"prayer_times_for":  "Prayer Times for",
		"prayer_times_dash": "Prayer Times - %s",
		"todays_dua":        "Today's Du'a",
	},
	"ar": {
		keyFajr:             "الفجر",
		keySunrise:          "الشروق",
		keyDhuhr:            "الظهر",
		keyAsr:              "العصر",
		keyMaghrib:          "المغرب",
		keyIsha:             "العشاء",
		keyMidnight:         "منتصف الليل",
		"prayer":            "الصلاة",
		"time":              "الوقت",
		"status":            "الحالة",
		"iqama":             "الإقامة",
		"passed":            "انقضت",
		"next":              "التالية",
		"next_prayer_in":    "الصلاة التالية خلال",
		"method":            "طريقة الحساب",
		"qibla":             "القبلة",
		"qibla_direction":   "اتجاه القبلة",
		"prayer_times_for":  "مواقيت الصلاة لـ",
		"prayer_times_dash": "مواقيت الصلاة - %s",
		"todays_dua":        "دعاء اليوم",
		"min":               "د",
		"hr":                "س",
		"hr_min":            "س %dد",
	},
}

// Time unit translations. English keeps the existing short forms ("min", "h", "Xh Ym").
var timeUnits = map[string]struct {
	min, hr, hrMin string
}{
	"en": {min: "min", hr: "h", hrMin: "h %dm"},
	"ar": {min: "د", hr: "س", hrMin: "س %dد"},
}

// t looks up a translation key for the given language, falling back to English.
func t(lang, key string) string {
	if lang == "ar" {
		if v, ok := translations["ar"][key]; ok {
			return v
		}
	}
	if v, ok := translations["en"][key]; ok {
		return v
	}
	return key
}

// hijriMonth returns the Arabic or English Hijri month name based on lang.
// Arabic output has tashkeel (diacritics) stripped, since many terminals render
// them as separate grapheme clusters and visually break cursive shaping.
func hijriMonth(lang string, m api.HijriMonthInfo) string {
	if lang == "ar" && m.Ar != "" {
		return stripTashkeel(m.Ar)
	}
	return m.En
}

// stripTashkeel removes Arabic diacritical marks (harakat) from a string.
// Ranges covered: U+064B..U+065F (tashkeel), U+0670 (superscript alef),
// U+06D6..U+06ED (Quranic annotation signs), U+0610..U+061A (honorifics).
func stripTashkeel(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 0x0610 && r <= 0x061A,
			r >= 0x064B && r <= 0x065F,
			r == 0x0670,
			r >= 0x06D6 && r <= 0x06ED:
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// formatDuration formats a duration in minutes using language-appropriate units.
func formatDuration(lang string, mins int) string {
	u, ok := timeUnits[lang]
	if !ok {
		u = timeUnits["en"]
	}
	if mins < 60 {
		return fmt.Sprintf("%d %s", mins, u.min)
	}
	hours := mins / 60
	remaining := mins % 60
	if remaining == 0 {
		return fmt.Sprintf("%d%s", hours, u.hr)
	}
	return fmt.Sprintf("%d"+u.hrMin, hours, remaining)
}
