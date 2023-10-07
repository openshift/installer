/*
Copyright (c) 2018 Red Hat, Inc.

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

// This file contains a logger that uses the standard output and error streams, or custom writers.

package logging

import (
	"context"
	"fmt"
	"io"
	"os"
)

// StdLoggerBuilder contains the configuration and logic needed to build a logger that uses the
// standard output and error streams, or custom writers.
type StdLoggerBuilder struct {
	debugEnabled bool
	infoEnabled  bool
	warnEnabled  bool
	errorEnabled bool
	outStream    io.Writer
	errStream    io.Writer
}

// StdLogger is a logger that uses the standard output and error streams, or custom writers.
type StdLogger struct {
	debugEnabled bool
	infoEnabled  bool
	warnEnabled  bool
	errorEnabled bool
	outStream    io.Writer
	errStream    io.Writer
}

// NewStdLoggerBuilder creates a builder that knows how to build a logger that uses the standard
// output and error streams, or custom writers. By default these loggers will have enabled the
// information, warning and error levels
func NewStdLoggerBuilder() *StdLoggerBuilder {
	// Allocate the object:
	builder := new(StdLoggerBuilder)

	// Set default values:
	builder.debugEnabled = false
	builder.infoEnabled = true
	builder.warnEnabled = true
	builder.errorEnabled = true

	return builder
}

// Streams sets the standard output and error streams to use. If not used then the logger will use
// os.Stdout and os.Stderr.
func (b *StdLoggerBuilder) Streams(out io.Writer, err io.Writer) *StdLoggerBuilder {
	b.outStream = out
	b.errStream = err
	return b
}

// Debug enables or disables the debug level.
func (b *StdLoggerBuilder) Debug(flag bool) *StdLoggerBuilder {
	b.debugEnabled = flag
	return b
}

// Info enables or disables the information level.
func (b *StdLoggerBuilder) Info(flag bool) *StdLoggerBuilder {
	b.infoEnabled = flag
	return b
}

// Warn enables or disables the warning level.
func (b *StdLoggerBuilder) Warn(flag bool) *StdLoggerBuilder {
	b.warnEnabled = flag
	return b
}

// Error enables or disables the error level.
func (b *StdLoggerBuilder) Error(flag bool) *StdLoggerBuilder {
	b.errorEnabled = flag
	return b
}

// Build creates a new logger using the configuration stored in the builder.
func (b *StdLoggerBuilder) Build() (logger *StdLogger, err error) {
	// Allocate and populate the object:
	logger = new(StdLogger)
	logger.debugEnabled = b.debugEnabled
	logger.infoEnabled = b.infoEnabled
	logger.warnEnabled = b.warnEnabled
	logger.errorEnabled = b.errorEnabled
	logger.outStream = b.outStream
	logger.errStream = b.errStream
	if logger.outStream == nil {
		logger.outStream = os.Stdout
	}
	if logger.errStream == nil {
		logger.errStream = os.Stderr
	}

	return
}

// DebugEnabled returns true iff the debug level is enabled.
func (l *StdLogger) DebugEnabled() bool {
	return l.debugEnabled
}

// InfoEnabled returns true iff the information level is enabled.
func (l *StdLogger) InfoEnabled() bool {
	return l.infoEnabled
}

// WarnEnabled returns true iff the warning level is enabled.
func (l *StdLogger) WarnEnabled() bool {
	return l.warnEnabled
}

// ErrorEnabled returns true iff the error level is enabled.
func (l *StdLogger) ErrorEnabled() bool {
	return l.errorEnabled
}

// Debug sends to the log a debug message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *StdLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	if l.debugEnabled {
		fmt.Fprintf(l.outStream, format+"\n", args...)
	}
}

// Info sends to the log an information message formatted using the fmt.Sprintf function and the
// given format and arguments.
func (l *StdLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.infoEnabled {
		fmt.Fprintf(l.outStream, format+"\n", args...)
	}
}

// Warn sends to the log a warning message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *StdLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.warnEnabled {
		fmt.Fprintf(l.outStream, format+"\n", args...)
	}
}

// Error sends to the log an error message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *StdLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.errorEnabled {
		fmt.Fprintf(l.errStream, format+"\n", args...)
	}
}

// Fatal sends to the log an error message formatted using the fmt.Sprintf function and the given
// format and arguments. After that it will os.Exit(1)
// This level is always enabled
func (l *StdLogger) Fatal(ctx context.Context, format string, args ...interface{}) {
	fmt.Fprintf(l.errStream, format+"\n", args...)
	os.Exit(1)
}
