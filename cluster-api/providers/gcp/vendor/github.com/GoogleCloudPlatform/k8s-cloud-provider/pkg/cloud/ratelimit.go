/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloud

import (
	"context"
	"time"
)

// RateLimitKey is a key identifying the operation to be rate limited. The rate limit
// queue will be determined based on the contents of RateKey.
//
// This type will be removed in a future release. Please change all
// references to CallContextKey.
type RateLimitKey = CallContextKey

// RateLimiter is the interface for a rate limiting policy.
type RateLimiter interface {
	// Accept uses the RateLimitKey to derive a sleep time for the calling
	// goroutine. This call will block until the operation is ready for
	// execution.
	//
	// Accept returns an error if the given context ctx was canceled
	// while waiting for acceptance into the queue.
	Accept(ctx context.Context, key *RateLimitKey) error
	// Observe uses the RateLimitKey to handle response results, which may affect
	// the sleep time for the Accept function.
	Observe(ctx context.Context, err error, key *RateLimitKey)
}

// acceptor is an object which blocks within Accept until a call is allowed to run.
// Accept is a behavior of the flowcontrol.RateLimiter interface.
type acceptor interface {
	// Accept blocks until a call is allowed to run.
	Accept()
}

// AcceptRateLimiter wraps an Acceptor with RateLimiter parameters.
type AcceptRateLimiter struct {
	// Acceptor is the underlying rate limiter.
	Acceptor acceptor
}

// Accept wraps an Acceptor and blocks on Accept or context.Done(). Key is ignored.
func (rl *AcceptRateLimiter) Accept(ctx context.Context, _ *RateLimitKey) error {
	ch := make(chan struct{})
	go func() {
		rl.Acceptor.Accept()
		close(ch)
	}()
	select {
	case <-ch:
		break
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// Observe does nothing.
func (rl *AcceptRateLimiter) Observe(context.Context, error, *RateLimitKey) {
}

// NopRateLimiter is a rate limiter that performs no rate limiting.
type NopRateLimiter struct {
}

// Accept everything immediately.
func (*NopRateLimiter) Accept(context.Context, *RateLimitKey) error {
	return nil
}

// Observe does nothing.
func (*NopRateLimiter) Observe(context.Context, error, *RateLimitKey) {
}

// MinimumRateLimiter wraps a RateLimiter and will only call its Accept until the minimum
// duration has been met or the context is cancelled.
type MinimumRateLimiter struct {
	// RateLimiter is the underlying ratelimiter which is called after the mininum time is reacehd.
	RateLimiter RateLimiter
	// Minimum is the minimum wait time before the underlying ratelimiter is called.
	Minimum time.Duration
}

// Accept blocks on the minimum duration and context. Once the minimum duration is met,
// the func is blocked on the underlying ratelimiter.
func (m *MinimumRateLimiter) Accept(ctx context.Context, key *RateLimitKey) error {
	t := time.NewTimer(m.Minimum)
	select {
	case <-t.C:
		return m.RateLimiter.Accept(ctx, key)
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	}
}

// Observe just passes error to the underlying ratelimiter.
func (m *MinimumRateLimiter) Observe(ctx context.Context, err error, key *RateLimitKey) {
	m.RateLimiter.Observe(ctx, err, key)
}

// TickerRateLimiter uses time.Ticker to spread Accepts over time.
//
// Concurrent calls to Accept will block on the same channel. It is not
// guaranteed what caller will be unblocked first.
type TickerRateLimiter struct {
	ticker *time.Ticker
}

// NewTickerRateLimiter creates a new TickerRateLimiter which will space Accept
// calls at least interval/limit time apart.
//
// For example, limit=4 interval=time.Minute will unblock a single Accept call
// every 15 seconds.
func NewTickerRateLimiter(limit int, interval time.Duration) *TickerRateLimiter {
	return &TickerRateLimiter{
		ticker: time.NewTicker(interval / time.Duration(limit)),
	}
}

// Accept will block until a time, specified when creating TickerRateLimiter,
// passes since the last call to Accept.
func (t *TickerRateLimiter) Accept(ctx context.Context, rlk *RateLimitKey) error {
	select {
	case <-t.ticker.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Observe does nothing.
func (*TickerRateLimiter) Observe(context.Context, error, *RateLimitKey) {
}

// Make sure that TickerRateLimiter implements RateLimiter.
var _ RateLimiter = new(TickerRateLimiter)

// CompositeRateLimiter combines rate limiters based on RateLimitKey.
type CompositeRateLimiter struct {
	// map[resource name]map[operation name]RateLimiter
	rateLimiters map[string]map[string]RateLimiter
	// defaultRL is used when no matching RateLimiter was found.
	defaultRL RateLimiter
}

// NewCompositeRateLimiter creates a new CompositeRateLimiter that will use
// provided default rate limiter if no better match is found. It is intended to
// be used for a single project.
//
// # Example
//
//	bsDefaultRL := /* backend service default rate limiter */
//	bsGetListRL := /* backend service rate limiter for get and list operations */
//
//	rl := NewCompositeRateLimiter(defaultRL)
//	rl.Register("BackendServices", "", bsDefaultRL)
//	rl.Register("BackendServices", "Get", bsGetListRL)
//	rl.Register("BackendServices", "List", bsGetListRL)
//
// This rate limiter is not nesting. Only one rate limiter is used for any
// particular combination of: resource, operation. For the case above, rate
// limiter registered at ("BackendServices", "") won't be applied to operation
// ("BackendServices", "Get"), because a more specific rate limiter was
// registered.
func NewCompositeRateLimiter(defaultRL RateLimiter) *CompositeRateLimiter {
	m := map[string]map[string]RateLimiter{
		"": {
			"": defaultRL,
		},
	}
	return &CompositeRateLimiter{
		rateLimiters: m,
		defaultRL:    defaultRL,
	}
}

// ensureExists creates sub-maps as needed.
func (c *CompositeRateLimiter) ensureExists(service string) {
	if _, ok := c.rateLimiters[service]; !ok {
		c.rateLimiters[service] = map[string]RateLimiter{}
	}
}

// fillMissing finds all combinations where resource and/or operation name
// could be omitted and sets it to defaultRL.
func (c *CompositeRateLimiter) fillMissing() {
	for _, subService := range c.rateLimiters {
		if subService[""] == nil {
			subService[""] = c.defaultRL
		}
	}
}

// Register adds provided rl to the composite rate limiter. Service, operation
// can be omitted by providing an empty string. In this case, the provided rate
// limiter will be used only when there is no other rate limiter matching a
// particular resource, or operation.
//
// It replaces previous rate limiter provided for the same service, operation
// combination. Once a rate limiter is added, it can't be removed.
//
// Same rate limiter can be used for multiple Register calls.
func (c *CompositeRateLimiter) Register(service, operation string, rl RateLimiter) {
	c.ensureExists(service)
	c.rateLimiters[service][operation] = rl
	c.fillMissing()
}

// Accept either calls underlying rate limiter matching rlk or a default rate
// limiter when none is found.
func (c *CompositeRateLimiter) Accept(ctx context.Context, rlk *RateLimitKey) error {
	if rlk == nil {
		return c.defaultRL.Accept(ctx, rlk)
	}
	service := rlk.Service
	if _, ok := c.rateLimiters[service]; !ok {
		service = ""
	}
	operation := rlk.Operation
	if _, ok := c.rateLimiters[service][operation]; !ok {
		operation = ""
	}
	return c.rateLimiters[service][operation].Accept(ctx, rlk)
}

// Observe does nothing.
func (*CompositeRateLimiter) Observe(context.Context, error, *RateLimitKey) {
}
