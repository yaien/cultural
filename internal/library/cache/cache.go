package cache

import (
	"sync"
	"time"
)

// Cache is a simple concurrency-safe in-memory cache with TTL support.
// Items expire after the configured ttl (unless overridden per-item) and
// are periodically cleaned by an internal goroutine. Remember to call
// Close() when the cache is no longer needed to stop the background worker.
type Cache[T any] struct {
	mu              sync.RWMutex
	items           map[string]item[T]
	ttl             time.Duration
	cleanupInterval time.Duration
	stop            chan struct{}
}

type item[T any] struct {
	value  T
	expiry time.Time
}

// New creates a Cache with the given default ttl for entries.
// If ttl is zero, entries do not expire by default (but you can still
// set per-item TTL using SetWithTTL).
func New[T any](ttl time.Duration) *Cache[T] {
	c := &Cache[T]{
		items: make(map[string]item[T]),
		ttl:   ttl,
		stop:  make(chan struct{}),
	}
	// default cleanup interval is half the ttl, with a sensible floor
	if ttl > 0 {
		c.cleanupInterval = ttl / 2
		if c.cleanupInterval < time.Millisecond*100 {
			c.cleanupInterval = time.Millisecond * 100
		}
	} else {
		c.cleanupInterval = time.Second
	}
	go c.startCleanup()
	return c
}

func (c *Cache[T]) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stop:
			return
		}
	}
}

func (c *Cache[T]) deleteExpired() {
	now := time.Now()
	c.mu.Lock()
	for k, it := range c.items {
		if !it.expiry.IsZero() && it.expiry.Before(now) {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// Set stores a value under key using the cache's default ttl.
func (c *Cache[T]) Set(key string, v T) {
	c.SetWithTTL(key, v, c.ttl)
}

// SetWithTTL stores a value under key with an explicit ttl. If ttl is zero,
// the item will not expire.
func (c *Cache[T]) SetWithTTL(key string, v T, ttl time.Duration) {
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.mu.Lock()
	c.items[key] = item[T]{value: v, expiry: exp}
	c.mu.Unlock()
}

// Get returns the value for key if present and not expired. If the item
// is expired it will be removed and (nil, false) is returned.
func (c *Cache[T]) Get(key string) (t T, ok bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return t, false
	}
	if !it.expiry.IsZero() && time.Now().After(it.expiry) {
		// expired — remove and report miss
		c.mu.Lock()
		// ensure we don't race with a newer value being set
		if cur, exists := c.items[key]; exists {
			if cur.expiry.Equal(it.expiry) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
		return t, false
	}
	return it.value, true
}

// Delete removes a key from the cache.
func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Close stops the background cleanup goroutine. After Close the cache can
// still be used, but periodic cleanup will no longer occur.
func (c *Cache[T]) Close() {
	select {
	case <-c.stop:
		// already closed
	default:
		close(c.stop)
	}
}
