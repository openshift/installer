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

// This file contains aliases for the types and functions that have been moved to the logging
// package.

package sdk

import (
	"github.com/openshift-online/ocm-sdk-go/logging"
)

// Logger has been moved to the logging package.
type Logger = logging.Logger

// GlogLoggerBuilder has been moved to the logging package.
type GlogLoggerBuilder = logging.GlogLoggerBuilder

// GlogLogger has been moved to the logging package.
type GlogLogger = logging.GlogLogger

// NewGlogLoggerBuilder has been moved to the logging package.
func NewGlogLoggerBuilder() *GlogLoggerBuilder {
	return logging.NewGlogLoggerBuilder()
}

// GoLoggerBuilder has been moved to the logging package.
type GoLoggerBuilder = logging.GoLoggerBuilder

// GoLogger has been moved to the logging package.
type GoLogger = logging.GoLogger

// NewGoLoggerBuilder has been moved to the logging package.
func NewGoLoggerBuilder() *GoLoggerBuilder {
	return logging.NewGoLoggerBuilder()
}

// StdLoggerBuilder has been moved to the logging package.
type StdLoggerBuilder = logging.StdLoggerBuilder

// StdLogger has been moved to the logging package.
type StdLogger = logging.StdLogger

// NewStdLoggerBuilder has been moved to the logging package.
func NewStdLoggerBuilder() *StdLoggerBuilder {
	return logging.NewStdLoggerBuilder()
}
