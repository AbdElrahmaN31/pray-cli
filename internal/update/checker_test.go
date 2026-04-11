package update

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCheck_UpdateAvailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"tag_name":"v2.0.0","name":"v2.0.0","body":"New release","html_url":"https://github.com/test/releases/v2.0.0","prerelease":false,"draft":false}`)
	}))
	defer server.Close()

	checker := NewChecker("1.0.0").WithReleasesURL(server.URL)
	result, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.UpdateAvailable {
		t.Error("expected update to be available")
	}
	if result.LatestVersion != "v2.0.0" {
		t.Errorf("expected latest version v2.0.0, got %s", result.LatestVersion)
	}
}

func TestCheck_NoUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"tag_name":"v1.0.0","name":"v1.0.0","body":"Current","html_url":"https://github.com/test/releases/v1.0.0","prerelease":false,"draft":false}`)
	}))
	defer server.Close()

	checker := NewChecker("1.0.0").WithReleasesURL(server.URL)
	result, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UpdateAvailable {
		t.Error("expected no update available")
	}
}

func TestCheck_Prerelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"tag_name":"v3.0.0","name":"v3.0.0","body":"Pre","html_url":"https://github.com/test/releases/v3.0.0","prerelease":true,"draft":false}`)
	}))
	defer server.Close()

	checker := NewChecker("1.0.0").WithReleasesURL(server.URL)
	result, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UpdateAvailable {
		t.Error("prerelease should not trigger update")
	}
}

func TestCheck_Draft(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"tag_name":"v3.0.0","name":"v3.0.0","body":"Draft","html_url":"https://github.com/test/releases/v3.0.0","prerelease":false,"draft":true}`)
	}))
	defer server.Close()

	checker := NewChecker("1.0.0").WithReleasesURL(server.URL)
	result, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UpdateAvailable {
		t.Error("draft should not trigger update")
	}
}

func TestCheck_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	checker := NewChecker("1.0.0").WithReleasesURL(server.URL)
	_, err := checker.Check(context.Background())
	if err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestCheck_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `not json`)
	}))
	defer server.Close()

	checker := NewChecker("1.0.0").WithReleasesURL(server.URL)
	_, err := checker.Check(context.Background())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCheckAsync(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"tag_name":"v2.0.0","name":"v2.0.0","body":"New","html_url":"https://github.com/test","prerelease":false,"draft":false}`)
	}))
	defer server.Close()

	checker := NewChecker("1.0.0").WithReleasesURL(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChan := checker.CheckAsync(ctx)
	result := <-resultChan
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !result.UpdateAvailable {
		t.Error("expected update to be available")
	}
}

func TestCheckAsync_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	checker := NewChecker("1.0.0").WithReleasesURL(server.URL)
	ctx := context.Background()

	resultChan := checker.CheckAsync(ctx)
	result := <-resultChan
	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestFormatUpdateMessage_WithUpdate(t *testing.T) {
	result := &CheckResult{
		UpdateAvailable: true,
		CurrentVersion:  "1.0.0",
		LatestVersion:   "v2.0.0",
		ReleaseURL:      "https://github.com/test/releases/v2.0.0",
	}

	msg := FormatUpdateMessage(result)
	if msg == "" {
		t.Error("expected non-empty message for available update")
	}
	if !strings.Contains(msg, "1.0.0") || !strings.Contains(msg, "v2.0.0") {
		t.Errorf("message should contain version info: %s", msg)
	}
}

func TestFormatUpdateMessage_NoUpdate(t *testing.T) {
	result := &CheckResult{
		UpdateAvailable: false,
		CurrentVersion:  "1.0.0",
	}

	msg := FormatUpdateMessage(result)
	if msg != "" {
		t.Errorf("expected empty message when no update, got: %s", msg)
	}
}

func TestFormatUpdateMessage_Nil(t *testing.T) {
	msg := FormatUpdateMessage(nil)
	if msg != "" {
		t.Errorf("expected empty message for nil result, got: %s", msg)
	}
}
