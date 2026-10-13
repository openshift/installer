/*
Copyright 2022 The Kubernetes Authors.

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

// Package logger provides a convenient interface to use to log.
package logger

import (
	"context"
	"fmt"

	smithylogging "github.com/aws/smithy-go/logging"
	"github.com/go-logr/logr"
)

// These are the log levels used by the logger.
// See https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/logging.md#what-method-to-use
const (
	logLevelWarn  = 1
	logLevelDebug = 4
	logLevelTrace = 5
)

// Wrapper defines a convenient interface to use to log things.
type Wrapper interface {
	Info(msg string, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Trace(msg string, keysAndValues ...any)
	Error(err error, msg string, keysAndValues ...any)
	WithValues(keysAndValues ...any) *Logger
	WithName(name string) *Logger
	GetLogger() logr.Logger
	GetAWSLogger() smithylogging.Logger
}

// Logger is a concrete logger using logr underneath.
type Logger struct {
	callStackHelper func()
	logger          logr.Logger
}

// NewLogger creates a logger with a passed in logr.Logger implementation directly.
func NewLogger(log logr.Logger) *Logger {
	helper, log := log.WithCallStackHelper()
	return &Logger{
		callStackHelper: helper,
		logger:          log,
	}
}

// FromContext retrieves the logr implementation from Context and uses it as underlying logger.
func FromContext(ctx context.Context) *Logger {
	helper, log := logr.FromContextOrDiscard(ctx).WithCallStackHelper()
	return &Logger{
		callStackHelper: helper,
		logger:          log,
	}
}

var _ Wrapper = &Logger{}

// Info logs a message at the info level.
func (c *Logger) Info(msg string, keysAndValues ...any) {
	c.callStackHelper()
	c.logger.Info(msg, keysAndValues...)
}

// Debug logs a message at the debug level.
func (c *Logger) Debug(msg string, keysAndValues ...any) {
	c.callStackHelper()
	c.logger.V(logLevelDebug).Info(msg, keysAndValues...)
}

// Warn logs a message at the warn level.
func (c *Logger) Warn(msg string, keysAndValues ...any) {
	c.callStackHelper()
	c.logger.V(logLevelWarn).Info(msg, keysAndValues...)
}

// Trace logs a message at the trace level.
func (c *Logger) Trace(msg string, keysAndValues ...any) {
	c.callStackHelper()
	c.logger.V(logLevelTrace).Info(msg, keysAndValues...)
}

// Error logs a message at the error level.
func (c *Logger) Error(err error, msg string, keysAndValues ...any) {
	c.callStackHelper()
	c.logger.Error(err, msg, keysAndValues...)
}

// GetLogger returns the underlying logr.Logger.
func (c *Logger) GetLogger() logr.Logger {
	return c.logger
}

// GetAWSLogger returns a compatible aws-sdk-go-v2 logger.
func (c Logger) GetAWSLogger() smithylogging.Logger {
	return c
}

// WithValues adds some key-value pairs of context to a logger.
func (c *Logger) WithValues(keysAndValues ...any) *Logger {
	return &Logger{
		callStackHelper: c.callStackHelper,
		logger:          c.logger.WithValues(keysAndValues...),
	}
}

// WithName adds a new element to the logger's name.
func (c *Logger) WithName(name string) *Logger {
	return &Logger{
		callStackHelper: c.callStackHelper,
		logger:          c.logger.WithName(name),
	}
}

// Logf allows Logger to satisfy the smithylogging.Logger interface.
func (c Logger) Logf(classification smithylogging.Classification, format string, v ...interface{}) {
	switch classification {
	// The only two classification levels are WARN and DEBUG
	case smithylogging.Debug:
		c.logger.V(logLevelDebug).Info(fmt.Sprintf(format, v...))
	default:
		c.logger.V(logLevelWarn).Info(fmt.Sprintf(format, v...))
	}
}
