// Package config provides configuration management for the pray CLI
package config

// CalculationMethod represents a prayer calculation method
type CalculationMethod struct {
	ID          int
	Name        string
	Description string
}

// CalculationMethods contains all available calculation methods
var CalculationMethods = []CalculationMethod{
	{ID: 0, Name: "Shia Ithna-Ashari", Description: "Shia Ithna-Ashari, Leva Institute, Qum"},
	{ID: 1, Name: "University of Islamic Sciences, Karachi", Description: "University of Islamic Sciences, Karachi"},
	{ID: 2, Name: "Islamic Society of North America", Description: "Islamic Society of North America (ISNA)"},
	{ID: 3, Name: "Muslim World League", Description: "Muslim World League (MWL)"},
	{ID: 4, Name: "Umm Al-Qura University, Makkah", Description: "Umm Al-Qura University, Makkah"},
	{ID: 5, Name: "Egyptian General Authority of Survey", Description: "Egyptian General Authority of Survey"},
	{ID: 6, Name: "Institute of Geophysics, University of Tehran", Description: "Institute of Geophysics, University of Tehran"},
	{ID: 7, Name: "Gulf Region", Description: "Gulf Region"},
	{ID: 8, Name: "Kuwait", Description: "Kuwait"},
	{ID: 9, Name: "Qatar", Description: "Qatar"},
	{ID: 10, Name: "Majlis Ugama Islam Singapura", Description: "Majlis Ugama Islam Singapura, Singapore"},
	{ID: 11, Name: "Union Organization Islamic de France", Description: "Union Organization Islamic de France"},
	{ID: 12, Name: "Diyanet İşleri Başkanlığı", Description: "Diyanet İşleri Başkanlığı, Turkey"},
	{ID: 13, Name: "Spiritual Administration of Muslims of Russia", Description: "Spiritual Administration of Muslims of Russia"},
	{ID: 14, Name: "Moonsighting Committee Worldwide", Description: "Moonsighting Committee Worldwide"},
	{ID: 15, Name: "Dubai", Description: "Dubai (experimental)"},
	{ID: 16, Name: "JAKIM", Description: "Jabatan Kemajuan Islam Malaysia (JAKIM)"},
	{ID: 17, Name: "Tunisia", Description: "Ministry of Religious Affairs, Tunisia"},
	{ID: 18, Name: "Algeria", Description: "Ministry of Religious Affairs and Wakfs, Algeria"},
	{ID: 19, Name: "KEMENAG", Description: "Kementerian Agama Republik Indonesia"},
	{ID: 20, Name: "Morocco", Description: "Ministry of Habous and Islamic Affairs, Morocco"},
	{ID: 21, Name: "Comunidade Islamica de Lisboa", Description: "Comunidade Islamica de Lisboa, Portugal"},
	{ID: 22, Name: "MUIS", Description: "Ministry of Religious Affairs of Jordan"},
	{ID: 23, Name: "Custom", Description: "Custom setting"},
}

// GetMethodByID returns a calculation method by its ID
func GetMethodByID(id int) *CalculationMethod {
	for _, method := range CalculationMethods {
		if method.ID == id {
			return &method
		}
	}
	return nil
}

// GetMethodName returns the name of a calculation method by ID
func GetMethodName(id int) string {
	method := GetMethodByID(id)
	if method != nil {
		return method.Name
	}
	return "Unknown"
}

// ValidMethodID checks if the method ID is valid
func ValidMethodID(id int) bool {
	return GetMethodByID(id) != nil
}

// DefaultOutputFormats lists available output formats
var DefaultOutputFormats = []string{
	"table",
	"pretty",
	"json",
	"slack",
	"discord",
	"webhook",
}

// DefaultLanguages lists available languages
var DefaultLanguages = []string{
	"en",
	"ar",
}

// PrayerNames contains the standard prayer names
var PrayerNames = []string{
	"Fajr",
	"Sunrise",
	"Dhuhr",
	"Asr",
	"Maghrib",
	"Isha",
	"Midnight",
}

// PrayerNamesArabic contains the Arabic prayer names
var PrayerNamesArabic = []string{
	"الفجر",
	"الشروق",
	"الظهر",
	"العصر",
	"المغرب",
	"العشاء",
	"منتصف الليل",
}

// PrayerEmojis contains emojis for each prayer
var PrayerEmojis = map[string]string{
	"Fajr":     "🌅",
	"Sunrise":  "🌄",
	"Dhuhr":    "☀️",
	"Asr":      "🌤️",
	"Maghrib":  "🌆",
	"Isha":     "🌙",
	"Midnight": "🌃",
}
