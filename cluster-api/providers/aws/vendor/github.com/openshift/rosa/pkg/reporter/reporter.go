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

// Builder contains the information and logic needed to create a new reporter.
type Builder struct {
	// Empty on purpose.
}

// Object is the reported object used by the tool. It prints the messages to the standard output or
// error streams.
type Object struct {
	errors int
}

// New creates a builder that can then be used to configure and build a reporter.
func New() *Builder {
	return &Builder{}
}

// Build uses the information contained in the builder to create a new reporter.
func (b *Builder) Build() (result *Object, err error) {
	// Create and populate the object:
	result = &Object{}

	return
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
		_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", infoPrefix, message)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "%s%s\n", "INFO: ", message)
	}
}

// Warnf prints an warning message with the given format and arguments.
func (r *Object) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if color.UseColor() {
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", warnPrefix, message)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", "WARN: ", message)
	}
}

// Errorf prints an error message with the given format and arguments. It also return an error
// containing the same information, which will be usually discarded, except when the caller needs to
// report the error and also return it.
func (r *Object) Errorf(format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	if color.UseColor() {
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", errorPrefix, message)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", "ERR: ", message)
	}
	r.errors++
	return errors.New(message)
}

// Errors returns the number of errors that have been reported via this reporter.
func (r *Object) Errors() int {
	return r.errors
}

// Message prefix using ANSI scape sequences to set colors:
const (
	infoPrefix  = "\033[0;36mI:\033[m "
	warnPrefix  = "\033[0;33mW:\033[m "
	errorPrefix = "\033[0;31mE:\033[m "
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

// CreateReporterOrExit creates the reportor instance or exits to the console
// noting the error on failure.
func CreateReporterOrExit() *Object {
	// Create the reporter:
	reporter, err := New().
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create reporter: %v\n", err)
		os.Exit(1)
	}
	return reporter
}
