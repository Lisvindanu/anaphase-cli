package ai

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Cache provides response caching functionality
type Cache struct {
	directory string
	ttl       time.Duration
	enabled   bool
}

// CachedEntry represents a cached response
type CachedEntry struct {
	Request    *GenerateRequest  `json:"request"`
	Response   *GenerateResponse `json:"response"`
	CachedAt   time.Time         `json:"cached_at"`
	ExpiresAt  time.Time         `json:"expires_at"`
	PromptHash string            `json:"prompt_hash"`
}

// NewCache creates a new cache instance
func NewCache(directory string, ttl time.Duration, enabled bool) *Cache {
	return &Cache{
		directory: directory,
		ttl:       ttl,
		enabled:   enabled,
	}
}

// Get retrieves cached response if available and not expired
func (c *Cache) Get(req *GenerateRequest) (*GenerateResponse, bool) {
	if !c.enabled {
		return nil, false
	}

	hash := c.hashRequest(req)
	cacheFile := filepath.Join(c.directory, hash+".json")

	// Check if cache file exists
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, false
	}

	// Parse cached entry
	var entry CachedEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	// Check expiration
	if time.Now().After(entry.ExpiresAt) {
		os.Remove(cacheFile) // Cleanup expired cache
		return nil, false
	}

	// Mark as cache hit
	entry.Response.CacheHit = true

	return entry.Response, true
}

// Set stores response in cache
func (c *Cache) Set(req *GenerateRequest, resp *GenerateResponse) error {
	if !c.enabled {
		return nil
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(c.directory, 0755); err != nil {
		return fmt.Errorf("create cache directory: %w", err)
	}

	hash := c.hashRequest(req)
	cacheFile := filepath.Join(c.directory, hash+".json")

	entry := CachedEntry{
		Request:    req,
		Response:   resp,
		CachedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(c.ttl),
		PromptHash: hash,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal cache entry: %w", err)
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("write cache file: %w", err)
	}

	return nil
}

// Clear removes all cached entries
func (c *Cache) Clear() error {
	return os.RemoveAll(c.directory)
}

// hashRequest creates a deterministic hash of the request
func (c *Cache) hashRequest(req *GenerateRequest) string {
	// Combine all request parameters into a single string
	combined := fmt.Sprintf("%s|%s|%.2f|%d|%.2f",
		req.SystemPrompt,
		req.UserPrompt,
		req.Temperature,
		req.MaxTokens,
		req.TopP,
	)

	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}
