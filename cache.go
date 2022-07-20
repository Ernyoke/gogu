package gogu

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

// Item holds the object itself (which could be of any type) and an expiration time.
// The expiration time defines the object lifetime.
type Item[V any] struct {
	Object     V
	Expiration int64
}

type cache[T ~string, V any] struct {
	mu    *sync.RWMutex
	items map[T]*Item[V]
	exp   time.Duration
	ci    time.Duration
	done  chan struct{}
}

// Exported cache struct.
type Cache[T ~string, V any] struct {
	*cache[T, V]
}

// newCache has a local scope only. NewCache will be used for the cache instantiation outside of this package.
func newCache[T ~string, V any](exp, ci time.Duration, item map[T]*Item[V]) *cache[T, V] {
	c := &cache[T, V]{
		mu:    &sync.RWMutex{},
		items: item,
		exp:   exp,
		ci:    ci,
		done:  make(chan struct{}),
	}
	return c
}

// NewCache initializes a new cache struct which request an expiration and a cleanup time.
// The cache will be invalidated once the expiration time is reached.
// A cleanup method is running in the background and removes the expired caches at a predifined interval.
func NewCache[T ~string, V any](expTime, cleanupTime time.Duration) *Cache[T, V] {
	items := make(map[T]*Item[V])
	c := newCache(expTime, cleanupTime, items)

	if expTime != NoExpiration {
		go c.cleanup()
		// Because there is no human interaction happening in the background, we need to have
		// a guarentee that the goroutine responsible for the cache purge stops after the cleanup.
		// This is the reason why runtime.SetFinalizer is used. This will be called when the garbage collector
		// finds an unreachable block ready to be garbage collected.
		runtime.SetFinalizer(c, stopCleanup[T, V])
	}

	return &Cache[T, V]{c}
}

// Set add a new item to the cache. At first verifies if an item with the same key already exists in the cache.
func (c *cache[T, V]) Set(key T, val V, d time.Duration) error {
	item, err := c.Get(key)
	if item != nil && err != nil {
		return fmt.Errorf("item with key '%v' already exists. Use the Update method", key)
	}
	c.add(key, val, d)

	return nil
}

// add place a new item into the cache. This method is not exported.
// It puts the item into the cache among with the expiration time.
func (c *cache[T, V]) add(key T, val V, d time.Duration) error {
	var exp int64

	if d == DefaultExpiration {
		d = c.exp
	}
	if d > 0 {
		exp = time.Now().Add(d).UnixNano()
	}

	item, err := c.Get(key)
	if item != nil && err != nil {
		return fmt.Errorf("item with key '%v' already exists", key)
	}

	c.mu.Lock()
	c.items[key] = &Item[V]{
		Object:     val,
		Expiration: exp,
	}
	c.mu.Unlock()

	return nil
}

// Get returns the item from the cache identifiable by its key.
// If an item expiration time is reached an error is returned instead of the item itself.
// Anyway the item will be purged by the cleanup method at the predifined interval.
func (c *cache[T, V]) Get(key T) (*Item[V], error) {
	c.mu.RLock()
	if item, ok := c.items[key]; ok {
		if item.Expiration > 0 {
			now := time.Now().UnixNano()
			if now > item.Expiration {
				c.mu.RUnlock()
				return nil, fmt.Errorf("item with key '%v' expired", key)
			}
		}
		c.mu.RUnlock()
		return c.items[key], nil
	}
	c.mu.RUnlock()
	return nil, fmt.Errorf("item with key '%v' not found", key)
}

// Update replaces an item from the cache with the new values.
func (c *cache[T, V]) Update(key T, val V, d time.Duration) error {
	c.mu.Lock()
	item, _ := c.Get(key)
	if item == nil {
		c.mu.Unlock()
		return fmt.Errorf("item with key '%v' does not exists", key)
	}
	c.Set(key, val, d)
	c.mu.Unlock()

	return nil
}

// SetDefault put a new item into the cache with the default expiration time.
func (c *cache[T, V]) SetDefault(key T, val V) {
	c.Set(key, val, DefaultExpiration)
}

// Delete deletes an item from the cache.
func (c *cache[T, V]) Delete(key T) error {
	item, _ := c.Get(key)
	if item != nil {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
	}

	return fmt.Errorf("item with key '%v' does not exists", key)
}

// DeleteExpired deletes all the expired items from the cache.
func (c *cache[T, V]) DeleteExpired() error {
	for k, item := range c.items {
		now := time.Now().UnixNano()
		if now > item.Expiration {
			return c.Delete(k)
		}
	}
	return nil
}

// List returns the cache items.
func (c *cache[T, V]) List() map[T]*Item[V] {
	return c.items
}

// MapToCache moves the items from a map into the cache.
func (c *cache[T, V]) MapToCache(m map[T]V) []error {
	errors := []error{}

	for k, v := range m {
		err := c.Set(k, v, DefaultExpiration)
		errors = append(errors, err)
	}

	return errors
}

// IsExpired checks if an item is expired or not.
func (c *cache[T, V]) IsExpired(key T) bool {
	item, err := c.Get(key)
	if item != nil && err != nil {
		if item.Expiration > time.Now().UnixNano() {
			return true
		}
	}
	return false
}

// cleanup runs the cache cleanup function on the specified interval an deletes all the expired cache items.
func (c *cache[T, V]) cleanup() {
	tick := time.NewTicker(c.ci)
	for {
		select {
		case <-tick.C:
			c.DeleteExpired()
		case <-c.done:
			tick.Stop()
		}
	}
}

// stopCleanup stops the cleanup process once the cache item goes out of scope
// and became unreachable and is being marked as garbage collactable by the GC.
func stopCleanup[T ~string, V any](c *cache[T, V]) {
	c.done <- struct{}{}
}