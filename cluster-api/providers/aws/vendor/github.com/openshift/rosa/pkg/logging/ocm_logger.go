/*
Copyright (c) 2020 Red Hat, Inc.

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

// This file contains an OCM logger that uses the logging framework of the project.

package logging

import (
	"context"
	"fmt"

	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/sirupsen/logrus"
)

// OCMLoggerBuilder contains the information and logic needed to create an OCM logger that uses
// the logging framework of the project. Don't create instances of this type directly; use the
// NewOCMLogger function instead.
type OCMLoggerBuilder struct {
	logger *logrus.Logger
}

// OCMLogger is an implementation of the OCM logger interface that uses the logging framework of
// the project. Don't create instances of this type directly; use the NewOCMLogger function instead.
type OCMLogger struct {
	logger *logrus.Logger
}

// Make sure that we implement the OCM logger interface.
var _ sdk.Logger = &OCMLogger{}

// NewOCMLogger creates new builder that can then be used to configure and build an OCM logger that
// uses the logging framework of the project.
func NewOCMLogger() *OCMLoggerBuilder {
	return &OCMLoggerBuilder{}
}

// Logger sets the underlying logger that will be used by the OCM logger to send the messages to the
// log.
func (b *OCMLoggerBuilder) Logger(value *logrus.Logger) *OCMLoggerBuilder {
	b.logger = value
	return b
}

// Build uses the information stored in the builder to create a new OCM logger that uses the logging
// framework of the project.
func (b *OCMLoggerBuilder) Build() (result *OCMLogger, err error) {
	// Check parameters:
	if b.logger == nil {
		err = fmt.Errorf("Logger is mandatory")
		return
	}

	// Create and populate the object:
	result = &OCMLogger{
		logger: b.logger,
	}

	return
}

func (l *OCMLogger) DebugEnabled() bool {
	return l.logger.IsLevelEnabled(logrus.DebugLevel)
}

func (l *OCMLogger) InfoEnabled() bool {
	return l.logger.IsLevelEnabled(logrus.InfoLevel)
}

func (l *OCMLogger) WarnEnabled() bool {
	return l.logger.IsLevelEnabled(logrus.WarnLevel)
}

func (l *OCMLogger) ErrorEnabled() bool {
	return l.logger.IsLevelEnabled(logrus.ErrorLevel)
}

func (l *OCMLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *OCMLogger) Info(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *OCMLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *OCMLogger) Error(ctx context.Context, format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *OCMLogger) Fatal(ctx context.Context, format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}
