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

// This file contains an AWS logger that uses the logging framework of the project.

package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// AWSLoggerBuilder contains the information and logic needed to create an AWS logger that uses
// the logging framework of the project. Don't create instances of this type directly; use the
// NewAWSLogger function instead.
type AWSLoggerBuilder struct {
	logger *logrus.Logger
}

// AWSLogger is an implementation of the OCM logger interface that uses the logging framework of
// the project. Don't create instances of this type directly; use the NewAWSLogger function instead.
type AWSLogger struct {
	logger *logrus.Logger
}

// Logger sets the underlying logger that will be used by the OCM logger to send the messages to the
// log.
func (b *AWSLoggerBuilder) Logger(value *logrus.Logger) *AWSLoggerBuilder {
	b.logger = value
	return b
}

// Build uses the information stored in the builder to create a new OCM logger that uses the logging
// framework of the project.
func (b *AWSLoggerBuilder) Build() (result *AWSLogger, err error) {
	// Check parameters:
	if b.logger == nil {
		err = fmt.Errorf("Logger is mandatory")
		return
	}

	// Create and populate the object:
	result = &AWSLogger{
		logger: b.logger,
	}

	return
}

func (l *AWSLogger) Log(args ...interface{}) {
	l.logger.Info(args...)
}
