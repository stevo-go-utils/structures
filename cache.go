package structures

import (
	"sync"
	"time"
)

type Cache[K comparable] struct {
	expiry            time.Duration
	items             map[K]time.Time
	autoDeleteCancels map[K]func()
	mu                sync.RWMutex
	opts              *CacheOpts
}

type CacheOpts struct {
	AutoDelete bool
}

type CacheOpt func(*CacheOpts)

func NewCacheOptions(opts ...CacheOpt) *CacheOpts {
	defaults := &CacheOpts{
		AutoDelete: false,
	}
	for _, o := range opts {
		o(defaults)
	}
	return defaults
}

func AutoDeleteCacheOpt() CacheOpt {
	return func(opts *CacheOpts) {
		opts.AutoDelete = true
	}
}

func NewCache[K comparable](expiry time.Duration, opts ...CacheOpt) *Cache[K] {
	return &Cache[K]{
		expiry:            expiry,
		items:             make(map[K]time.Time),
		autoDeleteCancels: make(map[K]func()),
		mu:                sync.RWMutex{},
		opts:              NewCacheOptions(opts...),
	}
}

func (c *Cache[K]) DeleteExpired() (deleted []K) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items {
		if now.Sub(v) >= c.expiry {
			delete(c.items, k)
			deleted = append(deleted, k)
		}
	}
	return deleted
}

func (c *Cache[K]) Add(keys ...K) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		c.items[key] = now
		if c.opts.AutoDelete {
			canceled := false
			cancelFunc := func() {
				canceled = true
			}
			if cancel, ok := c.autoDeleteCancels[key]; ok {
				cancel()
			}
			c.autoDeleteCancels[key] = cancelFunc
			go func(k K) {
				time.Sleep(c.expiry)
				c.mu.Lock()
				defer c.mu.Unlock()
				if !canceled {
					delete(c.items, k)
				}
			}(key)
		}
	}
}

func (c *Cache[K]) Delete(keys ...K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.items, key)
	}
}

func (c *Cache[K]) Contains(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.items[key]
	return ok
}

func (c *Cache[K]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *Cache[K]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]K, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

func (c *Cache[K]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[K]time.Time)
}
