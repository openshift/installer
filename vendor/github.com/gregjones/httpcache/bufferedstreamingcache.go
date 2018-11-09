package httpcache

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

// BufferedStreamingCache converts a Cache into a StreamingCache by using in-memory buffers.
type BufferedStreamingCache struct {
	cache Cache
}

// NewBufferedStreamingCache(cache Cache) wraps a Cache in a BufferedStreamingCache.
func NewBufferedStreamingCache(cache Cache) StreamingCache {
	return &BufferedStreamingCache{cache: cache}
}

// Get returns the []byte representation of a cached response and a bool
// set to true if the value isn't empty.
func (c *BufferedStreamingCache) Get(key string) (responseBytes []byte, ok bool) {
	return c.cache.Get(key)
}

// Set stores the []byte representation of a response against a key.
func (c *BufferedStreamingCache) Set(key string, responseBytes []byte) {
	c.cache.Set(key, responseBytes)
}

// Delete removes the value associated with the key.
func (c *BufferedStreamingCache) Delete(key string) {
	c.cache.Delete(key)
}

// GetReader streams data from the cache.  Returns os.ErrNotExist on cache misses.
func (c *BufferedStreamingCache) GetReader(key string) (response io.ReadCloser, err error) {
	data, ok := c.cache.Get(key)
	if !ok {
		return nil, os.ErrNotExist
	}

	reader := bytes.NewBuffer(data)
	return ioutil.NopCloser(reader), nil
}

// SetReader streams data into the cache.
func (c *BufferedStreamingCache) SetReader(key string, reader io.Reader) (err error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	c.cache.Set(key, data)
	return nil
}
