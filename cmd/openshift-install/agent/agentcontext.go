package agent

import (
	"context"
	"time"
)

// AgentContext is a context decorator.
type AgentContext struct {
	ctx context.Context
}

// NewAgentContext creates a new AgentContext instance.
func NewAgentContext(ctx context.Context) *AgentContext {
	return &AgentContext{
		ctx: ctx,
	}
}

// Deadline returns the context Deadline.
func (a *AgentContext) Deadline() (deadline time.Time, ok bool) {
	return a.ctx.Deadline()
}

// Done returns the context Done channel.
func (a *AgentContext) Done() <-chan struct{} {
	return a.ctx.Done()
}

// Err returns the context current error.
func (a *AgentContext) Err() error {
	return a.ctx.Err()
}

// Value returns the context Value.
func (a *AgentContext) Value(key any) any {
	return a.ctx.Value(key)
}

// AddValue allows to add a new value in the context.
func (a *AgentContext) AddValue(key, value interface{}) {
	a.ctx = context.WithValue(a.ctx, key, value)
}
