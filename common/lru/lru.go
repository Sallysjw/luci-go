// Copyright (c) 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package lru

import (
	"container/list"
	"sync"
)

// snapshot is a snapshot of the contents of the Cache.
type snapshot map[interface{}]interface{}

// Cache is a goroutine-safe least-recently-used (LRU) cache implementation. The
// cache stores key-value mapping entries up to a size limit. If more items are
// added past that limit, the entries that have have been referenced least
// recently will be evicted.
//
// This cache uses a read-write mutex, allowing multiple simultaneous
// non-mutating readers (Peek), but only one mutating reader/writer (Get, Put,
// Mutate).
type Cache interface {
	// Peek fetches the element associated with the supplied key without updating
	// the element's recently-used standing.
	Peek(key interface{}) interface{}

	// Get fetches the element associated with the supplied key, updating its
	// recently-used standing.
	Get(key interface{}) interface{}

	// Put adds a new value to the cache. The value in the cache will be replaced
	// regardless of whether an item with the same key already existed.
	//
	// Returns whether not a value already existed for the key.
	//
	// The new item will be considered most recently used.
	Put(key, value interface{}) (existed bool)

	// Mutate adds a value to the cache, using a generator to create the value.
	//
	// The generator will recieve the current value, or nil if there is no current
	// value, and will return the new value.
	//
	// The generator is called while the cache's lock is held. This means that
	// the generator MUST NOT call any cache methods during its execution, as
	// doing // so will result in deadlock/panic.
	//
	// Returns the value that was put in the cache, which is the value returned
	// by the generator.
	//
	// The key will be considered most recently used regardless of whether it was
	// put.
	Mutate(key interface{}, gen func(interface{}) interface{}) (value interface{})

	// Remove removes an entry from the cache. If the key is present, its current
	// value will be returned; otherwise, nil will be returned.
	Remove(key interface{}) interface{}

	// Purge clears the full contents of the cache.
	Purge()

	// Size returns the current cache size setting.
	Size() int

	// Len returns the number of entries in the cache.
	Len() int
}

// cacheImpl is a Cache interface implementation.
type cacheImpl struct {
	size int // The maximum number of elements that this cache should hold.

	cacheLock sync.RWMutex                  // Mutex to lock around cache reads/writes.
	cache     map[interface{}]*list.Element // Map of elements.
	lru       list.List                     // List of least-recently-used elements.
}

// New creates a new Cache instance with an initial size.
func New(size int) Cache {
	c := cacheImpl{
		size:  size,
		cache: make(map[interface{}]*list.Element),
	}
	c.lru.Init()
	return &c
}

func (c *cacheImpl) Peek(key interface{}) interface{} {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if e := c.cache[key]; e != nil {
		return e.Value
	}
	return nil
}

func (c *cacheImpl) Get(key interface{}) interface{} {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	if e := c.cache[key]; e != nil {
		c.lru.MoveToFront(e)
		return e.Value
	}
	return nil
}

func (c *cacheImpl) Put(key, value interface{}) (existed bool) {
	c.Mutate(key, func(current interface{}) interface{} {
		existed = (current != nil)
		return value
	})
	return
}

func (c *cacheImpl) Mutate(key interface{}, gen func(interface{}) interface{}) (value interface{}) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	e := c.cache[key]
	if e != nil {
		value = e.Value
	}
	value = gen(value)

	if e == nil {
		// The key doesn't currently exist. Create a new one and place it at the
		// front.
		e = c.lru.PushFront(nil)
		c.cache[key] = e
		c.pruneLocked()
	} else {
		// The element already exists. Visit it.
		c.lru.MoveToFront(e)
	}
	e.Value = value
	return
}

func (c *cacheImpl) Remove(key interface{}) interface{} {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	if e, ok := c.cache[key]; ok {
		delete(c.cache, key)
		c.lru.Remove(e)
		return e.Value
	}
	return nil
}

func (c *cacheImpl) Purge() {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	c.cache = make(map[interface{}]*list.Element)
	c.lru.Init()
}

func (c *cacheImpl) Size() int {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	return c.size
}

func (c *cacheImpl) Len() int {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	return int(len(c.cache))
}

// keys returns a list of keys in the cache.
func (c *cacheImpl) keys() []interface{} {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	var keys []interface{}
	if len(c.cache) > 0 {
		keys = make([]interface{}, 0, len(c.cache))
		for k := range c.cache {
			keys = append(keys, k)
		}
	}
	return keys
}

// snapshot returns a snapshot map of the cache's entries.
func (c *cacheImpl) snapshot() (ss snapshot) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if len(c.cache) > 0 {
		ss = make(snapshot)
		for k, e := range c.cache {
			ss[k] = e.Value
		}
	}
	return
}

// cacheLock's write lock must be held by the caller.
func (c *cacheImpl) pruneLocked() {
	for int(c.lru.Len()) > c.size {
		e := c.lru.Back()
		delete(c.cache, e.Value)
		c.lru.Remove(e)
	}
}