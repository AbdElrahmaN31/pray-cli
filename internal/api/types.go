// Package api provides types and client for the prayer times API
package api

// PrayerTimesResponse represents the JSON response from the prayer times API
type PrayerTimesResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

// Data contains the main prayer times data
type Data struct {
	Timings Timings `json:"timings"`
	Date    Date    `json:"date"`
	Meta    Meta    `json:"meta"`
}

// Timings holds all prayer times
type Timings struct {
	Fajr       string `json:"Fajr"`
	Sunrise    string `json:"Sunrise"`
	Dhuhr      string `json:"Dhuhr"`
	Asr        string `json:"Asr"`
	Sunset     string `json:"Sunset"`
	Maghrib    string `json:"Maghrib"`
	Isha       string `json:"Isha"`
	Imsak      string `json:"Imsak"`
	Midnight   string `json:"Midnight"`
	Firstthird string `json:"Firstthird"`
	Lastthird  string `json:"Lastthird"`
}

// Date contains both Gregorian and Hijri date information
type Date struct {
	Readable  string        `json:"readable"`
	Timestamp string        `json:"timestamp"`
	Gregorian GregorianDate `json:"gregorian"`
	Hijri     HijriDate     `json:"hijri"`
}

// GregorianDate represents the Gregorian calendar date
type GregorianDate struct {
	Date        string      `json:"date"`
	Format      string      `json:"format"`
	Day         string      `json:"day"`
	Weekday     Weekday     `json:"weekday"`
	Month       MonthInfo   `json:"month"`
	Year        string      `json:"year"`
	Designation Designation `json:"designation"`
}

// HijriDate represents the Islamic (Hijri) calendar date
type HijriDate struct {
	Date        string         `json:"date"`
	Format      string         `json:"format"`
	Day         string         `json:"day"`
	Weekday     Weekday        `json:"weekday"`
	Month       HijriMonthInfo `json:"month"`
	Year        string         `json:"year"`
	Designation Designation    `json:"designation"`
	Holidays    []string       `json:"holidays"`
}

// Weekday contains the weekday name in English and Arabic
type Weekday struct {
	En string `json:"en"`
	Ar string `json:"ar,omitempty"`
}

// MonthInfo contains month information for Gregorian dates
type MonthInfo struct {
	Number int    `json:"number"`
	En     string `json:"en"`
}

// HijriMonthInfo contains month information for Hijri dates
type HijriMonthInfo struct {
	Number int    `json:"number"`
	En     string `json:"en"`
	Ar     string `json:"ar"`
}

// Designation contains calendar era designation
type Designation struct {
	Abbreviated string `json:"abbreviated"`
	Expanded    string `json:"expanded"`
}

// Meta contains metadata about the calculation
type Meta struct {
	Latitude                 float64 `json:"latitude"`
	Longitude                float64 `json:"longitude"`
	Timezone                 string  `json:"timezone"`
	Method                   Method  `json:"method"`
	LatitudeAdjustmentMethod string  `json:"latitudeAdjustmentMethod"`
	MidnightMode             string  `json:"midnightMode"`
	School                   string  `json:"school"`
	Offset                   Offset  `json:"offset"`
}

// Method contains calculation method information
type Method struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Params   MethodParams   `json:"params"`
	Location MethodLocation `json:"location,omitempty"`
}

// MethodParams contains the calculation parameters
type MethodParams struct {
	Fajr    interface{} `json:"Fajr"`
	Isha    interface{} `json:"Isha"`
	Maghrib interface{} `json:"Maghrib,omitempty"`
}

// MethodLocation contains location info for the method
type MethodLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Offset contains time offsets for each prayer
type Offset struct {
	Imsak    int `json:"Imsak"`
	Fajr     int `json:"Fajr"`
	Sunrise  int `json:"Sunrise"`
	Dhuhr    int `json:"Dhuhr"`
	Asr      int `json:"Asr"`
	Maghrib  int `json:"Maghrib"`
	Sunset   int `json:"Sunset"`
	Isha     int `json:"Isha"`
	Midnight int `json:"Midnight"`
}

// QiblaResponse represents the Qibla direction response
type QiblaResponse struct {
	Code   int       `json:"code"`
	Status string    `json:"status"`
	Data   QiblaData `json:"data"`
}

// QiblaData contains Qibla direction information
type QiblaData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Direction float64 `json:"direction"`
}

// NextPrayer represents the next upcoming prayer
type NextPrayer struct {
	Name         string `json:"name"`
	Time         string `json:"time"`
	ISO          string `json:"iso,omitempty"`
	Timestamp    int64  `json:"timestamp,omitempty"`
	MinutesUntil int    `json:"minutesUntil"`
}

// PrayerTimesOutput represents the formatted output for display
type PrayerTimesOutput struct {
	Date       DateOutput        `json:"date"`
	Location   LocationOutput    `json:"location"`
	Method     string            `json:"method"`
	Timings    map[string]string `json:"timings"`
	NextPrayer *NextPrayer       `json:"nextPrayer,omitempty"`
	Qibla      *QiblaOutput      `json:"qibla,omitempty"`
}

// DateOutput contains formatted date information
type DateOutput struct {
	Gregorian string       `json:"gregorian"`
	Hijri     *HijriOutput `json:"hijri,omitempty"`
}

// HijriOutput contains formatted Hijri date
type HijriOutput struct {
	Day   string         `json:"day"`
	Month HijriMonthInfo `json:"month"`
	Year  string         `json:"year"`
}

// LocationOutput contains location information for output
type LocationOutput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Address   string  `json:"address,omitempty"`
}

// QiblaOutput contains formatted Qibla direction
type QiblaOutput struct {
	Direction float64 `json:"direction"`
	Compass   string  `json:"compass"`
}
