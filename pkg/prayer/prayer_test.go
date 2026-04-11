package prayer

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tz := time.UTC
	date := time.Date(2026, 3, 15, 0, 0, 0, 0, tz)

	tests := []struct {
		name    string
		input   string
		wantH   int
		wantM   int
		wantErr bool
	}{
		{"valid", "05:15", 5, 15, false},
		{"midnight", "00:00", 0, 0, false},
		{"noon", "12:00", 12, 0, false},
		{"late", "23:59", 23, 59, false},
		{"invalid", "abc", 0, 0, true},
		{"empty", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTime(tt.input, date, tz)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil {
				if got.Hour() != tt.wantH || got.Minute() != tt.wantM {
					t.Errorf("ParseTime(%q) = %02d:%02d, want %02d:%02d", tt.input, got.Hour(), got.Minute(), tt.wantH, tt.wantM)
				}
			}
		})
	}
}

func TestGetNextPrayer(t *testing.T) {
	now := time.Date(2026, 3, 15, 13, 0, 0, 0, time.UTC)
	prayers := []Prayer{
		{Name: "Fajr", IsPassed: true},
		{Name: "Sunrise", IsPassed: true},
		{Name: "Dhuhr", IsPassed: true},
		{Name: "Asr", IsPassed: false},
		{Name: "Maghrib", IsPassed: false},
		{Name: "Isha", IsPassed: false},
	}

	next := GetNextPrayer(prayers, now)
	if next == nil {
		t.Fatal("GetNextPrayer() returned nil")
	}
	if next.Name != "Asr" {
		t.Errorf("GetNextPrayer() = %s, want Asr", next.Name)
	}
	if !next.IsNext {
		t.Error("IsNext should be true for the next prayer")
	}
}

func TestGetNextPrayerAllPassed(t *testing.T) {
	now := time.Date(2026, 3, 15, 23, 0, 0, 0, time.UTC)
	prayers := []Prayer{
		{Name: "Fajr", IsPassed: true},
		{Name: "Isha", IsPassed: true},
	}

	next := GetNextPrayer(prayers, now)
	if next != nil {
		t.Errorf("GetNextPrayer() = %v, want nil when all passed", next)
	}
}

func TestCalculateTimeDiff(t *testing.T) {
	base := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		from time.Time
		to   time.Time
		want int
	}{
		{"same", base, base, 0},
		{"30min", base, base.Add(30 * time.Minute), 30},
		{"2hours", base, base.Add(2 * time.Hour), 120},
		{"negative", base.Add(1 * time.Hour), base, -60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTimeDiff(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("CalculateTimeDiff() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		minutes int
		want    string
	}{
		{-1, "passed"},
		{0, "0 min"},
		{30, "30 min"},
		{59, "59 min"},
		{60, "1h"},
		{90, "1h 30m"},
		{120, "2h"},
		{150, "2h 30m"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatDuration(tt.minutes)
			if got != tt.want {
				t.Errorf("FormatDuration(%d) = %q, want %q", tt.minutes, got, tt.want)
			}
		})
	}
}

func TestGetCompassDirection(t *testing.T) {
	tests := []struct {
		degrees float64
		want    string
	}{
		{0, "N"},
		{45, "NE"},
		{90, "E"},
		{135, "SE"},
		{180, "S"},
		{225, "SW"},
		{270, "W"},
		{315, "NW"},
		{360, "N"},
		{-90, "W"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := GetCompassDirection(tt.degrees)
			if got != tt.want {
				t.Errorf("GetCompassDirection(%.0f) = %q, want %q", tt.degrees, got, tt.want)
			}
		})
	}
}

func TestIsPrayerPassed(t *testing.T) {
	now := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	if !IsPrayerPassed(past, now) {
		t.Error("IsPrayerPassed() = false for past time")
	}
	if IsPrayerPassed(future, now) {
		t.Error("IsPrayerPassed() = true for future time")
	}
}

func TestPrayerNameByIndex(t *testing.T) {
	tests := []struct {
		index int
		want  string
	}{
		{FajrIndex, "Fajr"},
		{SunriseIndex, "Sunrise"},
		{DhuhrIndex, "Dhuhr"},
		{AsrIndex, "Asr"},
		{MaghribIndex, "Maghrib"},
		{IshaIndex, "Isha"},
		{MidnightIndex, "Midnight"},
		{-1, ""},
		{99, ""},
	}

	for _, tt := range tests {
		got := PrayerNameByIndex(tt.index)
		if got != tt.want {
			t.Errorf("PrayerNameByIndex(%d) = %q, want %q", tt.index, got, tt.want)
		}
	}
}

func TestGetMethod(t *testing.T) {
	// Valid method
	m := GetMethod(5)
	if m == nil {
		t.Fatal("GetMethod(5) returned nil")
	}
	if m.Name != "Egyptian General Authority of Survey" {
		t.Errorf("GetMethod(5).Name = %q, want Egyptian General Authority of Survey", m.Name)
	}

	// New method (14+)
	m14 := GetMethod(14)
	if m14 == nil {
		t.Fatal("GetMethod(14) returned nil")
	}
	if m14.Name != "Moonsighting Committee Worldwide" {
		t.Errorf("GetMethod(14).Name = %q", m14.Name)
	}

	// Invalid method
	if GetMethod(99) != nil {
		t.Error("GetMethod(99) should return nil")
	}
}

func TestGetAllMethods(t *testing.T) {
	methods := GetAllMethods()
	if len(methods) != 24 {
		t.Errorf("GetAllMethods() returned %d methods, want 24", len(methods))
	}

	// Check first and last
	if methods[0].ID != 0 {
		t.Errorf("first method ID = %d, want 0", methods[0].ID)
	}
	if methods[len(methods)-1].ID != 23 {
		t.Errorf("last method ID = %d, want 23", methods[len(methods)-1].ID)
	}
}

func TestGetDailyDua(t *testing.T) {
	// Should return a dua for any date
	date1 := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	dua1 := GetDailyDua(date1)
	if dua1 == nil {
		t.Fatal("GetDailyDua() returned nil")
	}
	if dua1.Arabic == "" {
		t.Error("Arabic text should not be empty")
	}
	if dua1.Transliteration == "" {
		t.Error("Transliteration should not be empty")
	}
	if dua1.Translation == "" {
		t.Error("Translation should not be empty")
	}
	if dua1.Reference == "" {
		t.Error("Reference should not be empty")
	}

	// Same date should return same dua
	dua1b := GetDailyDua(date1)
	if dua1.Arabic != dua1b.Arabic {
		t.Error("same date should return same dua")
	}

	// Different date should (likely) return different dua
	date2 := time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC)
	dua2 := GetDailyDua(date2)
	if dua2 == nil {
		t.Fatal("GetDailyDua() returned nil for different date")
	}
	// Not testing inequality since modulo could theoretically collide
}

func TestDuasCollection(t *testing.T) {
	if len(Duas) < 20 {
		t.Errorf("expected at least 20 duas, got %d", len(Duas))
	}

	// Verify all entries have required fields
	for i, dua := range Duas {
		if dua.Arabic == "" {
			t.Errorf("Duas[%d] has empty Arabic", i)
		}
		if dua.Transliteration == "" {
			t.Errorf("Duas[%d] has empty Transliteration", i)
		}
		if dua.Translation == "" {
			t.Errorf("Duas[%d] has empty Translation", i)
		}
		if dua.Reference == "" {
			t.Errorf("Duas[%d] has empty Reference", i)
		}
	}
}
