package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestCache(t *testing.T, opts ...Option) *Cache {
	t.Helper()
	dir := t.TempDir()
	c, err := New(dir, opts...)
	if err != nil {
		t.Fatalf("failed to create cache: %v", err)
	}
	return c
}

func TestNew(t *testing.T) {
	dir := t.TempDir()
	c, err := New(dir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if !c.IsEnabled() {
		t.Error("expected cache to be enabled by default")
	}
	if c.ttl != DefaultTTL {
		t.Errorf("expected default TTL %v, got %v", DefaultTTL, c.ttl)
	}
}

func TestNewWithOptions(t *testing.T) {
	dir := t.TempDir()
	ttl := 1 * time.Hour
	c, err := New(dir, WithTTL(ttl), WithEnabled(false))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if c.IsEnabled() {
		t.Error("expected cache to be disabled")
	}
	if c.ttl != ttl {
		t.Errorf("expected TTL %v, got %v", ttl, c.ttl)
	}
}

func TestSetAndGet(t *testing.T) {
	c := setupTestCache(t)

	key := "test-key"
	data := []byte(`{"prayer":"fajr","time":"05:15"}`)

	if err := c.Set(key, data); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	got, found := c.Get(key)
	if !found {
		t.Fatal("Get() returned not found for existing key")
	}
	if string(got) != string(data) {
		t.Errorf("Get() = %s, want %s", got, data)
	}
}

func TestGetNonExistent(t *testing.T) {
	c := setupTestCache(t)

	_, found := c.Get("nonexistent")
	if found {
		t.Error("Get() returned found for non-existent key")
	}
}

func TestGetExpired(t *testing.T) {
	c := setupTestCache(t, WithTTL(1*time.Millisecond))

	key := "expiring"
	if err := c.Set(key, []byte(`"data"`)); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	_, found := c.Get(key)
	if found {
		t.Error("Get() returned found for expired key")
	}

	// File should have been removed
	path := c.getPath(key)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expired cache file was not removed")
	}
}

func TestGetCorruptedFile(t *testing.T) {
	c := setupTestCache(t)

	key := "corrupted"
	path := c.getPath(key)
	if err := os.WriteFile(path, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	_, found := c.Get(key)
	if found {
		t.Error("Get() returned found for corrupted cache entry")
	}

	// File should have been removed
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("corrupted cache file was not removed")
	}
}

func TestDelete(t *testing.T) {
	c := setupTestCache(t)

	key := "to-delete"
	if err := c.Set(key, []byte(`"data"`)); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	if err := c.Delete(key); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, found := c.Get(key)
	if found {
		t.Error("Get() returned found after Delete()")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	c := setupTestCache(t)
	if err := c.Delete("nonexistent"); err != nil {
		t.Errorf("Delete() non-existent key error = %v", err)
	}
}

func TestClear(t *testing.T) {
	c := setupTestCache(t)

	for i := 0; i < 5; i++ {
		if err := c.Set(GenerateKey("key", i), []byte(`"data"`)); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	if err := c.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	entries, _, err := c.Stats()
	if err != nil {
		t.Fatalf("Stats() error = %v", err)
	}
	if entries != 0 {
		t.Errorf("expected 0 entries after Clear(), got %d", entries)
	}
}

func TestCleanExpired(t *testing.T) {
	c := setupTestCache(t, WithTTL(1*time.Millisecond))

	for i := 0; i < 3; i++ {
		if err := c.Set(GenerateKey("expired", i), []byte(`"data"`)); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	time.Sleep(5 * time.Millisecond)

	removed, err := c.CleanExpired()
	if err != nil {
		t.Fatalf("CleanExpired() error = %v", err)
	}
	if removed != 3 {
		t.Errorf("CleanExpired() removed %d, want 3", removed)
	}
}

func TestStats(t *testing.T) {
	c := setupTestCache(t)

	data := []byte(`{"test":"data"}`)
	for i := 0; i < 3; i++ {
		if err := c.Set(GenerateKey("stats", i), data); err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	entries, totalSize, err := c.Stats()
	if err != nil {
		t.Fatalf("Stats() error = %v", err)
	}
	if entries != 3 {
		t.Errorf("Stats() entries = %d, want 3", entries)
	}
	if totalSize == 0 {
		t.Error("Stats() totalSize = 0, want > 0")
	}
}

func TestExists(t *testing.T) {
	c := setupTestCache(t)

	key := "exists-test"
	if c.Exists(key) {
		t.Error("Exists() = true for non-existent key")
	}

	if err := c.Set(key, []byte(`"data"`)); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	if !c.Exists(key) {
		t.Error("Exists() = false for existing key")
	}
}

func TestDisabledCache(t *testing.T) {
	c := setupTestCache(t, WithEnabled(false))

	if err := c.Set("key", []byte(`"data"`)); err != nil {
		t.Fatalf("Set() on disabled cache error = %v", err)
	}

	_, found := c.Get("key")
	if found {
		t.Error("Get() returned found on disabled cache")
	}
}

func TestSetEnabled(t *testing.T) {
	c := setupTestCache(t)
	c.SetEnabled(false)
	if c.IsEnabled() {
		t.Error("expected disabled after SetEnabled(false)")
	}
	c.SetEnabled(true)
	if !c.IsEnabled() {
		t.Error("expected enabled after SetEnabled(true)")
	}
}

func TestGenerateKey(t *testing.T) {
	key1 := GenerateKey("times", 30.0444, 31.2357, "2026-02-03", 5)
	key2 := GenerateKey("times", 30.0444, 31.2357, "2026-02-03", 5)
	key3 := GenerateKey("times", 30.0444, 31.2357, "2026-02-04", 5)

	if key1 != key2 {
		t.Error("same params should produce same key")
	}
	if key1 == key3 {
		t.Error("different params should produce different key")
	}
	if len(key1) != 16 {
		t.Errorf("key length = %d, want 16", len(key1))
	}
}

func TestNonJsonFilesIgnored(t *testing.T) {
	c := setupTestCache(t)

	// Create a non-json file in cache dir
	nonJsonPath := filepath.Join(c.dir, "readme.txt")
	if err := os.WriteFile(nonJsonPath, []byte("not a cache file"), 0644); err != nil {
		t.Fatalf("failed to write non-json file: %v", err)
	}

	// Set a real cache entry
	if err := c.Set("real-key", []byte(`"data"`)); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Clear should only remove .json files
	if err := c.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	// Non-json file should still exist
	if _, err := os.Stat(nonJsonPath); os.IsNotExist(err) {
		t.Error("Clear() removed non-json file")
	}
}
