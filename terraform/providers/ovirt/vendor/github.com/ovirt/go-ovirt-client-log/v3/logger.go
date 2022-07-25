package ovirtclientlog

import "context"

// Logger provides pluggable logging for oVirt client libraries.
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

	// WithContext returns a logger that adheres to a specific context. This is useful when the backend logger needs
	// access to the context the current operation is running under.
	WithContext(ctx context.Context) Logger
}
