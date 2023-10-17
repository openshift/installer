/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ttllru

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
)

type (
	// Cache is a TTL LRU cache which caches items with a max time to live and with
	// bounded length.
	Cache struct {
		Cacher
		TimeToLive time.Duration
		mu         sync.Mutex
	}

	// Cacher describes a basic cache.
	Cacher interface {
		Get(key interface{}) (value interface{}, ok bool)
		Add(key interface{}, value interface{}) (evicted bool)
		Remove(key interface{}) (ok bool)
	}

	// PeekingCacher describes a basic cache with the ability to peek.
	PeekingCacher interface {
		Cacher
		Peek(key interface{}) (value interface{}, expiration time.Time, ok bool)
	}

	timeToLiveItem struct {
		LastTouch time.Time
		Value     interface{}
	}
)

// New creates a new TTL LRU cache which caches items with a max time to live and with
// bounded length.
func New(size int, timeToLive time.Duration) (PeekingCacher, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build new LRU cache")
	}

	return newCache(timeToLive, c)
}

func newCache(timeToLive time.Duration, cache Cacher) (PeekingCacher, error) {
	return &Cache{
		Cacher:     cache,
		TimeToLive: timeToLive,
	}, nil
}

// Get returns a value and a bool indicating the value was found for a given key.
func (ttlCache *Cache) Get(key interface{}) (value interface{}, ok bool) {
	ttlItem, ok := ttlCache.peekItem(key)
	if !ok {
		return nil, false
	}

	ttlItem.LastTouch = time.Now()
	return ttlItem.Value, true
}

// Add will add a value for a given key.
func (ttlCache *Cache) Add(key interface{}, val interface{}) bool {
	ttlCache.mu.Lock()
	defer ttlCache.mu.Unlock()

	return ttlCache.Cacher.Add(key, &timeToLiveItem{
		Value:     val,
		LastTouch: time.Now(),
	})
}

// Peek will fetch an item from the cache, but will not update the expiration time.
func (ttlCache *Cache) Peek(key interface{}) (value interface{}, expiration time.Time, ok bool) {
	ttlItem, ok := ttlCache.peekItem(key)
	if !ok {
		return nil, time.Time{}, false
	}

	expirationTime := time.Now().Add(ttlCache.TimeToLive - time.Since(ttlItem.LastTouch))
	return ttlItem.Value, expirationTime, true
}

func (ttlCache *Cache) peekItem(key interface{}) (value *timeToLiveItem, ok bool) {
	ttlCache.mu.Lock()
	defer ttlCache.mu.Unlock()

	val, ok := ttlCache.Cacher.Get(key)
	if !ok {
		return nil, false
	}

	ttlItem, ok := val.(*timeToLiveItem)
	if !ok {
		return nil, false
	}

	if time.Since(ttlItem.LastTouch) > ttlCache.TimeToLive {
		ttlCache.Cacher.Remove(key)
		return nil, false
	}

	return ttlItem, true
}
