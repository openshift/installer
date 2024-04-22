package sigstore

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// HTTPClient returns a client suitable for retrieving signatures. It is not
// required to be unique per call, but may be called concurrently.
type HTTPClient func() (*http.Client, error)

// DefaultClient creates an http.Client with no configuration.
func DefaultClient() (*http.Client, error) {
	return &http.Client{}, nil
}

// CachedHTTPClientConstructor wraps an HTTPClient implementation so
// that it is not called more frequently than the configured limiter.
type CachedHTTPClientConstructor struct {
	wrapped HTTPClient
	limiter *rate.Limiter

	lock       sync.Mutex
	lastClient *http.Client
	lastError  error
}

// NewCachedHTTPClientConstructor creates a new cached constructor.
// If limiter is not specified it defaults to one call every 30 seconds.
func NewCachedHTTPClientConstructor(wrapped HTTPClient, limiter *rate.Limiter) *CachedHTTPClientConstructor {
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Every(30*time.Second), 1)
	}
	return &CachedHTTPClientConstructor{
		wrapped: wrapped,
		limiter: limiter,
	}
}

func (c *CachedHTTPClientConstructor) HTTPClient() (*http.Client, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	r := c.limiter.Reserve()
	if r.OK() {
		c.lastClient, c.lastError = c.wrapped()
	}
	return c.lastClient, c.lastError
}
