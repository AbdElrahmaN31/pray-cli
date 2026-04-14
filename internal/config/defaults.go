// Package config provides configuration management for the pray CLI
package config

// CalculationMethod represents a prayer calculation method
type CalculationMethod struct {
	ID          int
	Name        string
	NameAr      string
	Description string
}

// CalculationMethods contains all available calculation methods
var CalculationMethods = []CalculationMethod{
	{ID: 0, Name: "Shia Ithna-Ashari", NameAr: "الشيعة الإثنا عشرية", Description: "Shia Ithna-Ashari, Leva Institute, Qum"},
	{ID: 1, Name: "University of Islamic Sciences, Karachi", NameAr: "جامعة العلوم الإسلامية، كراتشي", Description: "University of Islamic Sciences, Karachi"},
	{ID: 2, Name: "Islamic Society of North America", NameAr: "الجمعية الإسلامية لأمريكا الشمالية", Description: "Islamic Society of North America (ISNA)"},
	{ID: 3, Name: "Muslim World League", NameAr: "رابطة العالم الإسلامي", Description: "Muslim World League (MWL)"},
	{ID: 4, Name: "Umm Al-Qura University, Makkah", NameAr: "جامعة أم القرى، مكة المكرمة", Description: "Umm Al-Qura University, Makkah"},
	{ID: 5, Name: "Egyptian General Authority of Survey", NameAr: "الهيئة المصرية العامة للمساحة", Description: "Egyptian General Authority of Survey"},
	{ID: 6, Name: "Institute of Geophysics, University of Tehran", NameAr: "معهد الجيوفيزياء، جامعة طهران", Description: "Institute of Geophysics, University of Tehran"},
	{ID: 7, Name: "Gulf Region", NameAr: "منطقة الخليج", Description: "Gulf Region"},
	{ID: 8, Name: "Kuwait", NameAr: "الكويت", Description: "Kuwait"},
	{ID: 9, Name: "Qatar", NameAr: "قطر", Description: "Qatar"},
	{ID: 10, Name: "Majlis Ugama Islam Singapura", NameAr: "مجلس الشؤون الإسلامية، سنغافورة", Description: "Majlis Ugama Islam Singapura, Singapore"},
	{ID: 11, Name: "Union Organization Islamic de France", NameAr: "الاتحاد الإسلامي الفرنسي", Description: "Union Organization Islamic de France"},
	{ID: 12, Name: "Diyanet İşleri Başkanlığı", NameAr: "رئاسة الشؤون الدينية التركية", Description: "Diyanet İşleri Başkanlığı, Turkey"},
	{ID: 13, Name: "Spiritual Administration of Muslims of Russia", NameAr: "الإدارة الدينية لمسلمي روسيا", Description: "Spiritual Administration of Muslims of Russia"},
	{ID: 14, Name: "Moonsighting Committee Worldwide", NameAr: "لجنة رؤية الهلال العالمية", Description: "Moonsighting Committee Worldwide"},
	{ID: 15, Name: "Dubai", NameAr: "دبي", Description: "Dubai (experimental)"},
	{ID: 16, Name: "JAKIM", NameAr: "جاكيم، ماليزيا", Description: "Jabatan Kemajuan Islam Malaysia (JAKIM)"},
	{ID: 17, Name: "Tunisia", NameAr: "وزارة الشؤون الدينية، تونس", Description: "Ministry of Religious Affairs, Tunisia"},
	{ID: 18, Name: "Algeria", NameAr: "وزارة الشؤون الدينية والأوقاف، الجزائر", Description: "Ministry of Religious Affairs and Wakfs, Algeria"},
	{ID: 19, Name: "KEMENAG", NameAr: "وزارة الشؤون الدينية، إندونيسيا", Description: "Kementerian Agama Republik Indonesia"},
	{ID: 20, Name: "Morocco", NameAr: "وزارة الأوقاف والشؤون الإسلامية، المغرب", Description: "Ministry of Habous and Islamic Affairs, Morocco"},
	{ID: 21, Name: "Comunidade Islamica de Lisboa", NameAr: "الجالية الإسلامية بلشبونة، البرتغال", Description: "Comunidade Islamica de Lisboa, Portugal"},
	{ID: 22, Name: "MUIS", NameAr: "وزارة الشؤون الدينية، الأردن", Description: "Ministry of Religious Affairs of Jordan"},
	{ID: 23, Name: "Custom", NameAr: "مخصص", Description: "Custom setting"},
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

// GetMethodNameLang returns the localized name of a calculation method by ID.
// Falls back to English if Arabic is requested but not available.
func GetMethodNameLang(id int, lang string) string {
	method := GetMethodByID(id)
	if method == nil {
		return "Unknown"
	}
	if lang == "ar" && method.NameAr != "" {
		return method.NameAr
	}
	return method.Name
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
