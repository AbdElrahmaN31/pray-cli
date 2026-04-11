package update

import "testing"

func TestCanonicalVersion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"v1.0.0", "v1.0.0"},
		{"1.0.0", "v1.0.0"},
		{"  v2.3.1  ", "v2.3.1"},
		{"v0.1.0", "v0.1.0"},
		{"", ""},
	}

	for _, tt := range tests {
		got := canonicalVersion(tt.input)
		if got != tt.want {
			t.Errorf("canonicalVersion(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		current string
		latest  string
		want    bool
	}{
		{"1.0.0", "1.0.1", true},
		{"1.0.0", "1.1.0", true},
		{"1.0.0", "2.0.0", true},
		{"1.0.1", "1.0.0", false},
		{"1.0.0", "1.0.0", false},
		{"dev", "1.0.0", false},
		{"", "1.0.0", false},
		{"0.1.0", "0.2.0", true},
		{"v1.0.0", "v1.0.1", true},
		{"v1.0.0", "v2.0.0", true},
		{"1.0.0-beta", "1.0.0", true},
		{"1.0.0", "1.0.0-beta", false},
		{"1.0.0-alpha", "1.0.0-beta", true},
	}

	for _, tt := range tests {
		got := isNewerVersion(tt.current, tt.latest)
		if got != tt.want {
			t.Errorf("isNewerVersion(%q, %q) = %v, want %v", tt.current, tt.latest, got, tt.want)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"this is a long string", 10, "this is..."},
		{"exact", 5, "exact"},
		{"", 10, ""},
	}

	for _, tt := range tests {
		got := truncateString(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}
