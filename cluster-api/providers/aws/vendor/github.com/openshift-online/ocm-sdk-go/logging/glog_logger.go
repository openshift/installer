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

// This file contains a logger that uses the `glog` package.

package logging

import (
	"context"
	"fmt"
	"os"

	"github.com/golang/glog"
)

// GlogLoggerBuilder contains the configuration and logic needed to build a logger that uses the
// glog V mechanism. Don't create instances of this type directly, use the NewGlogLoggerBuilder
// function instead.
type GlogLoggerBuilder struct {
	debugV glog.Level
	infoV  glog.Level
	warnV  glog.Level
	errorV glog.Level
}

// GlogLogger is a logger that uses the glog V mechanism.
type GlogLogger struct {
	debugV glog.Level
	infoV  glog.Level
	warnV  glog.Level
	errorV glog.Level
}

// NewGlogLoggerBuilder creates a builder that uses the glog V mechanism. By default errors,
// warnings and information messages will be written to the log if the level is 0 or greater, and
// debug messages will be written if the level is 1 or greater. This can be changed using the
// ErrorV, WarnV, InfoV and DebugV methods of the builder. For example, to write errors and warnings
// for level 0, information messages for level 1, and debug messages for level 2, you can create the
// logger like this:
//
//	logger, err := client.NewGlobLoggerBuilder().
//		ErrorV(0).
//		WarnV(0).
//		InfoV(1).
//		DebugV(2).
//		Build()
//
// Once the logger is created these settings can't be changed.
func NewGlogLoggerBuilder() *GlogLoggerBuilder {
	// Allocate the object:
	builder := new(GlogLoggerBuilder)

	// Set default values:
	builder.debugV = 1
	builder.infoV = 0
	builder.warnV = 0
	builder.errorV = 0

	return builder
}

// DebugV sets the V value that will be used for debug messages.
func (b *GlogLoggerBuilder) DebugV(v glog.Level) *GlogLoggerBuilder {
	b.debugV = v
	return b
}

// InfoV sets the V value that will be used for info messages.
func (b *GlogLoggerBuilder) InfoV(v glog.Level) *GlogLoggerBuilder {
	b.infoV = v
	return b
}

// WarnV sets the V value that will be used for warn messages.
func (b *GlogLoggerBuilder) WarnV(v glog.Level) *GlogLoggerBuilder {
	b.warnV = v
	return b
}

// ErrorV sets the V value that will be used for error messages.
func (b *GlogLoggerBuilder) ErrorV(v glog.Level) *GlogLoggerBuilder {
	b.errorV = v
	return b
}

// Build creates a new logger using the configuration stored in the builder.
func (b *GlogLoggerBuilder) Build() (logger *GlogLogger, err error) {
	// Allocate and populate the object:
	logger = new(GlogLogger)
	logger.debugV = b.debugV
	logger.infoV = b.infoV
	logger.warnV = b.warnV
	logger.errorV = b.errorV

	return
}

// DebugEnabled returns true iff the debug level is enabled.
func (l *GlogLogger) DebugEnabled() bool {
	return bool(glog.V(l.debugV))
}

// InfoEnabled returns true iff the information level is enabled.
func (l *GlogLogger) InfoEnabled() bool {
	return bool(glog.V(l.infoV))
}

// WarnEnabled returns true iff the warning level is enabled.
func (l *GlogLogger) WarnEnabled() bool {
	return bool(glog.V(l.warnV))
}

// ErrorEnabled returns true iff the error level is enabled.
func (l *GlogLogger) ErrorEnabled() bool {
	return bool(glog.V(l.errorV))
}

// Debug sends to the log a debug message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *GlogLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	if glog.V(l.debugV) {
		msg := fmt.Sprintf(format, args...)
		glog.InfoDepth(1, msg)
	}
}

// Info sends to the log an information message formatted using the fmt.Sprintf function and the
// given format and arguments.
func (l *GlogLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if glog.V(l.infoV) {
		msg := fmt.Sprintf(format, args...)
		glog.InfoDepth(1, msg)
	}
}

// Warn sends to the log a warning message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *GlogLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if glog.V(l.warnV) {
		msg := fmt.Sprintf(format, args...)
		glog.WarningDepth(1, msg)
	}
}

// Error sends to the log an error message formatted using the fmt.Sprintf function and the given
// format and arguments.
func (l *GlogLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if glog.V(l.errorV) {
		msg := fmt.Sprintf(format, args...)
		glog.ErrorDepth(1, msg)
	}
}

// Fatal sends to the log an error message formatted using the fmt.Sprintf function and the given
// format and arguments. After that it will os.Exit(1)
// This level is always enabled
func (l *GlogLogger) Fatal(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	// #nosec G104
	glog.ErrorDepth(1, msg)
	os.Exit(1)
}
