package aws

import (
	"sync"
	"time"
)

// CacheEntry is an entry with a value and expiry timer in the cache.
type CacheEntry[T any] struct {
	Value T
	Timer *time.Timer
}

// Cache is a simple thread-safe in-memory cache with TTL support.
type Cache[T any] struct {
	Store sync.Map
}

// NewCache creates and returns a new Cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{}
}

// Set stores an entry in the cache with a specified time-to-live (TTL).
// If an entry with the same key already exists, its value is updated,
// and its TTL is reset.
func (c *Cache[T]) Set(key string, value T, ttl time.Duration) {
	entry := CacheEntry[T]{
		Value: value,
		Timer: time.AfterFunc(ttl, func() {
			c.Store.Delete(key)
		}),
	}

	// Swap the new entry with the old one. If an old entry existed,
	// its timer is stopped to prevent premature deletion.
	if existing, loaded := c.Store.Swap(key, entry); loaded {
		if oldEntry, ok := existing.(CacheEntry[T]); ok {
			oldEntry.Timer.Stop()
		}
	}
}

// Get retrieves an entry from the cache.
// It returns the value and a boolean indicating whether the key was found.
func (c *Cache[T]) Get(key string) (T, bool) {
	var zero T
	existing, ok := c.Store.Load(key)
	if !ok {
		return zero, false
	}

	if entry, ok := existing.(CacheEntry[T]); ok {
		return entry.Value, true
	}

	return zero, false
}

// Delete removes an entry from the cache and stops its expiry timer,
// and returns the existing entry if any.
func (c *Cache[T]) Delete(key string) {
	if existing, loaded := c.Store.LoadAndDelete(key); loaded {
		if oldEntry, ok := existing.(CacheEntry[T]); ok {
			oldEntry.Timer.Stop()
		}
	}
}
