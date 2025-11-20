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

package reporter

import (
	"errors"
	"fmt"
	"os"

	"github.com/openshift/rosa/pkg/color"
	"github.com/openshift/rosa/pkg/debug"
)

type Logger interface {
	// Debugf logs a debug message with formatted arguments.
	Debugf(format string, args ...interface{})

	// Errorf logs an error message with formatted arguments and returns an error.
	Errorf(format string, args ...interface{}) error

	// Infof logs an info message with formatted arguments.
	Infof(format string, args ...interface{})

	// IsTerminal checks if the output is a terminal.
	IsTerminal() bool

	// Warnf logs a warning message with formatted arguments.
	Warnf(format string, args ...interface{})
}

// Object is the reported object used by the tool. It prints the messages to the standard output or
// error streams.
type Object struct {
}

// Debugf prints a debug message with the given format and arguments.
func (r *Object) Debugf(format string, args ...interface{}) {
	if !debug.Enabled() {
		return
	}
	r.Infof(format, args...)
}

// Infof prints an informative message with the given format and arguments.
func (r *Object) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if color.UseColor() {
		_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", infoColorPrefix, message)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", infoPrefix, message)
	}
}

// Warnf prints an warning message with the given format and arguments.
func (r *Object) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if color.UseColor() {
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", warnColorPrefix, message)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", warnPrefix, message)
	}
}

// Errorf prints an error message with the given format and arguments. It also return an error
// containing the same information, which will be usually discarded, except when the caller needs to
// report the error and also return it.
//
//nolint:errcheck
func (r *Object) Errorf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	if color.UseColor() {
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", errorColorPrefix, message)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", errorPrefix, message)
	}
	return errors.New(message)
}

// Message prefix using ANSI scape sequences to set colors:
const (
	infoColorPrefix  = "\033[0;36mI:\033[m "
	infoPrefix       = "INFO: "
	warnColorPrefix  = "\033[0;33mW:\033[m "
	warnPrefix       = "WARN: "
	errorColorPrefix = "\033[0;31mE:\033[m "
	errorPrefix      = "ERR: "
)

// Determine whether the reporter output is meant for the terminal
// or whether it's piped or redirected to a file.
func (r *Object) IsTerminal() bool {
	stdout, err := os.Stdout.Stat()
	if err != nil {
		return true
	}
	return (stdout.Mode()&os.ModeDevice != 0) && (stdout.Mode()&os.ModeNamedPipe == 0)
}

// CreateReporter returns a new reporter
func CreateReporter() *Object {
	return &Object{}
}
