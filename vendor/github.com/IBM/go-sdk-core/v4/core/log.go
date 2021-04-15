package core

// (C) Copyright IBM Corp. 2020.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"log"
	"os"
	"sync"
)

// LogLevel defines a type for logging levels
type LogLevel int

// Log level constants
const (
	LevelNone LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

// Logger is the logging interface implemented and used by the Go core library.
// Users of the library can supply their own implementation by calling SetLogger().
type Logger interface {
	Log(level LogLevel, format string, inserts ...interface{})
	Error(format string, inserts ...interface{})
	Warn(format string, inserts ...interface{})
	Info(format string, inserts ...interface{})
	Debug(format string, inserts ...interface{})
}

// SDKLoggerImpl is the Go core's implementation of the Logger interface.
type SDKLoggerImpl struct {

	// The current log level configured in this logger.
	// Only messages with a log level that is <= this configured log level
	// will be displayed.
	logLevel LogLevel

	// The underlying log.Logger instance that will be used to log each message.
	goLogger *log.Logger

	loggerInit sync.Once
}

// logImpl returns the underlying log.Logger instance to be used to do the actual logging.
func (l *SDKLoggerImpl) logImpl() *log.Logger {
	l.loggerInit.Do(func() {
		if l.goLogger == nil {
			l.goLogger = log.New(os.Stdout, "", log.LstdFlags)
		}
	})

	return l.goLogger
}

// Log will log the specified message if "level" is currently enabled.
func (l *SDKLoggerImpl) Log(level LogLevel, format string, inserts ...interface{}) {
	if level <= l.logLevel {
		l.logImpl().Printf(format, inserts...)
	}
}

// Error logs a message at level "Error"
func (l *SDKLoggerImpl) Error(format string, inserts ...interface{}) {
	l.Log(LevelError, "[Error] "+format, inserts...)
}

// Warn logs a message at level "Warn"
func (l *SDKLoggerImpl) Warn(format string, inserts ...interface{}) {
	l.Log(LevelWarn, "[Warn] "+format, inserts...)
}

// Info logs a message at level "Info"
func (l *SDKLoggerImpl) Info(format string, inserts ...interface{}) {
	l.Log(LevelInfo, "[Info] "+format, inserts...)
}

// Debug logs a message at level "Debug"
func (l *SDKLoggerImpl) Debug(format string, inserts ...interface{}) {
	l.Log(LevelDebug, "[Debug] "+format, inserts...)
}

// NewLogger constructs an SDKLoggerImpl instance with the specified logging level
// enabled and the specified log.Logger instance as the underlying logger to use.
// If "stdLogger" is nil, then a default log.Logger instance will be used.
func NewLogger(level LogLevel, stdLogger *log.Logger) *SDKLoggerImpl {
	return &SDKLoggerImpl{
		logLevel: level,
		goLogger: stdLogger,
	}
}

// sdkLogger holds the Logger implementation used by the Go core library.
var sdkLogger Logger = NewLogger(LevelError, nil)

// SetLogger sets the specified Logger instance as the logger to be used by the Go core library.
func SetLogger(logger Logger) {
	sdkLogger = logger
}

// GetLogger returns the Logger instance currently used by the Go core.
func GetLogger() Logger {
	return sdkLogger
}

// SetLoggingLevel will enable the specified logging level in the Go core library.
func SetLoggingLevel(level LogLevel) {
	SetLogger(NewLogger(level, nil))
}
