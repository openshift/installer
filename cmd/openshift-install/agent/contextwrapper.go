package agent

import (
	"context"
	"time"
)

// ContextWrapper is a context decorator.
type ContextWrapper struct {
	ctx context.Context //nolint:containedctx
}

// NewContextWrapper creates a new ContextWrapper instance.
func NewContextWrapper(ctx context.Context) *ContextWrapper {
	return &ContextWrapper{
		ctx: ctx,
	}
}

// Deadline returns the context Deadline.
func (a *ContextWrapper) Deadline() (deadline time.Time, ok bool) {
	return a.ctx.Deadline()
}

// Done returns the context Done channel.
func (a *ContextWrapper) Done() <-chan struct{} {
	return a.ctx.Done()
}

// Err returns the context current error.
func (a *ContextWrapper) Err() error {
	return a.ctx.Err()
}

// Value returns the context Value.
func (a *ContextWrapper) Value(key any) any {
	return a.ctx.Value(key)
}

// AddValue allows to add a new value in the context.
func (a *ContextWrapper) AddValue(key, value interface{}) {
	a.ctx = context.WithValue(a.ctx, key, value)
}
