// Package cache provides file-based caching for API responses
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// DefaultTTL is the default time-to-live for cache entries
	DefaultTTL = 24 * time.Hour
)

// Entry represents a cached item with metadata
type Entry struct {
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	ExpiresAt time.Time       `json:"expires_at"`
	Key       string          `json:"key"`
}

// Cache provides file-based caching functionality
type Cache struct {
	dir     string
	ttl     time.Duration
	enabled bool
}

// Option configures the Cache
type Option func(*Cache)

// WithTTL sets a custom TTL for cache entries
func WithTTL(ttl time.Duration) Option {
	return func(c *Cache) {
		c.ttl = ttl
	}
}

// WithEnabled enables or disables the cache
func WithEnabled(enabled bool) Option {
	return func(c *Cache) {
		c.enabled = enabled
	}
}

// New creates a new Cache instance
func New(dir string, opts ...Option) (*Cache, error) {
	c := &Cache{
		dir:     dir,
		ttl:     DefaultTTL,
		enabled: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	// Create cache directory if it doesn't exist
	if c.enabled {
		if err := os.MkdirAll(c.dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache directory: %w", err)
		}
	}

	return c, nil
}

// GenerateKey creates a unique cache key from parameters
func GenerateKey(params ...interface{}) string {
	h := sha256.New()
	for _, p := range params {
		fmt.Fprintf(h, "%v", p)
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// Get retrieves a cached entry if it exists and is not expired
func (c *Cache) Get(key string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}

	path := c.getPath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		// Invalid cache file, remove it
		os.Remove(path)
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		os.Remove(path)
		return nil, false
	}

	return entry.Data, true
}

// Set stores data in the cache
func (c *Cache) Set(key string, data []byte) error {
	if !c.enabled {
		return nil
	}

	entry := Entry{
		Data:      data,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(c.ttl),
		Key:       key,
	}

	entryData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	path := c.getPath(key)
	if err := os.WriteFile(path, entryData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// Delete removes a specific entry from the cache
func (c *Cache) Delete(key string) error {
	path := c.getPath(key)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete cache entry: %w", err)
	}
	return nil
}

// Clear removes all entries from the cache
func (c *Cache) Clear() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			path := filepath.Join(c.dir, entry.Name())
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove cache file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// CleanExpired removes all expired entries from the cache
func (c *Cache) CleanExpired() (int, error) {
	if !c.enabled {
		return 0, nil
	}

	entries, err := os.ReadDir(c.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read cache directory: %w", err)
	}

	removed := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(c.dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var cacheEntry Entry
		if err := json.Unmarshal(data, &cacheEntry); err != nil {
			os.Remove(path)
			removed++
			continue
		}

		if time.Now().After(cacheEntry.ExpiresAt) {
			os.Remove(path)
			removed++
		}
	}

	return removed, nil
}

// Stats returns cache statistics
func (c *Cache) Stats() (entries int, totalSize int64, err error) {
	entryList, err := os.ReadDir(c.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, nil
		}
		return 0, 0, fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range entryList {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		entries++
		totalSize += info.Size()
	}

	return entries, totalSize, nil
}

// IsEnabled returns whether caching is enabled
func (c *Cache) IsEnabled() bool {
	return c.enabled
}

// SetEnabled enables or disables the cache
func (c *Cache) SetEnabled(enabled bool) {
	c.enabled = enabled
}

// getPath returns the file path for a cache key
func (c *Cache) getPath(key string) string {
	return filepath.Join(c.dir, key+".json")
}

// Exists checks if a cache entry exists and is not expired
func (c *Cache) Exists(key string) bool {
	_, found := c.Get(key)
	return found
}
