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
		items:             map[K]time.Time{},
		autoDeleteCancels: map[K]func(){},
		mu:                sync.RWMutex{},
		opts:              NewCacheOptions(opts...),
	}
}

func (c *Cache[K]) DeleteExpired() (deleted []K) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items {
		if now.Sub(v) >= 0 {
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
		c.items[key] = now.Add(c.expiry)
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

func (c *Cache[K]) AddWithExpiry(key K, dur time.Duration) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = now.Add(dur)
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
			time.Sleep(dur)
			c.mu.Lock()
			defer c.mu.Unlock()
			if !canceled {
				delete(c.items, k)
			}
		}(key)
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

type CacheMap[K comparable, V any] struct {
	expiry            time.Duration
	items             map[K]V
	itemExpiries      map[K]time.Time
	autoDeleteCancels map[K]func()
	mu                sync.RWMutex
	opts              *CacheOpts
}

func NewCacheMap[K comparable, V any](expiry time.Duration, opts ...CacheOpt) *CacheMap[K, V] {
	return &CacheMap[K, V]{
		expiry:            expiry,
		items:             map[K]V{},
		itemExpiries:      map[K]time.Time{},
		autoDeleteCancels: map[K]func(){},
		mu:                sync.RWMutex{},
		opts:              NewCacheOptions(opts...),
	}
}

func (c *CacheMap[K, V]) DeleteExpired() (deleted []K) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, expiry := range c.itemExpiries {
		if now.Sub(expiry) >= 0 {
			delete(c.items, k)
			delete(c.itemExpiries, k)
			deleted = append(deleted, k)
		}
	}
	return deleted
}

func (c *CacheMap[K, V]) Add(key K, value V) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.itemExpiries[key] = now.Add(c.expiry)
	c.items[key] = value
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
				delete(c.itemExpiries, k)
			}
		}(key)
	}
}

func (c *CacheMap[K, V]) AddWithExpiry(key K, value V, dur time.Duration) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.itemExpiries[key] = now.Add(dur)
	c.items[key] = value
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
			time.Sleep(dur)
			c.mu.Lock()
			defer c.mu.Unlock()
			if !canceled {
				delete(c.items, k)
				delete(c.itemExpiries, k)
			}
		}(key)
	}
}

func (c *CacheMap[K, V]) Delete(keys ...K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.items, key)
		delete(c.itemExpiries, key)
	}
}

func (c *CacheMap[K, V]) Has(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.items[key]
	return ok
}

func (c *CacheMap[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *CacheMap[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]K, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

func (c *CacheMap[K, V]) Vals() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	vals := make([]V, 0, len(c.items))
	for _, v := range c.items {
		vals = append(vals, v)
	}
	return vals
}

func (c *CacheMap[K, V]) Get(key K) (value V, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok = c.items[key]
	return value, ok
}

func (c *CacheMap[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[K]V{}
	c.itemExpiries = map[K]time.Time{}
}
