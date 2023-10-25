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

// Package logger
package logger

import (
	"context"

	"github.com/go-logr/logr"
)

const (
	logLevelDebug = 2
	logLevelWarn  = 3
	logLevelTrace = 4
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
}

// Logger is a concrete logger using logr underneath.
type Logger struct {
	logger logr.Logger
}

// NewLogger creates a logger with a passed in logr.Logger implementation directly.
func NewLogger(log logr.Logger) *Logger {
	return &Logger{
		logger: log,
	}
}

// FromContext retrieves the logr implementation from Context and uses it as underlying logger.
func FromContext(ctx context.Context) *Logger {
	log := logr.FromContextOrDiscard(ctx)
	return &Logger{
		logger: log,
	}
}

var _ Wrapper = &Logger{}

func (c *Logger) Info(msg string, keysAndValues ...any) {
	c.logger.Info(msg, keysAndValues...)
}

func (c *Logger) Debug(msg string, keysAndValues ...any) {
	c.logger.V(logLevelDebug).Info(msg, keysAndValues...)
}

func (c *Logger) Warn(msg string, keysAndValues ...any) {
	c.logger.V(logLevelWarn).Info(msg, keysAndValues...)
}

func (c *Logger) Trace(msg string, keysAndValues ...any) {
	c.logger.V(logLevelTrace).Info(msg, keysAndValues...)
}

func (c *Logger) Error(err error, msg string, keysAndValues ...any) {
	c.logger.Error(err, msg, keysAndValues...)
}

func (c *Logger) GetLogger() logr.Logger {
	return c.logger
}

func (c *Logger) WithValues(keysAndValues ...any) *Logger {
	c.logger = c.logger.WithValues(keysAndValues...)
	return c
}

func (c *Logger) WithName(name string) *Logger {
	c.logger = c.logger.WithName(name)
	return c
}
