/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package conditions

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// ReadyConditionImpactingError is an error that requires notification in the Ready condition
type ReadyConditionImpactingError struct {
	Severity            ConditionSeverity
	Reason              string
	cause               error
	RetryClassification RetryClassification
}

// NewReadyConditionImpactingError creates a new ReadyConditionImpactingError
func NewReadyConditionImpactingError(cause error, severity ConditionSeverity, reason Reason) *ReadyConditionImpactingError {
	return &ReadyConditionImpactingError{
		cause:               cause,
		Severity:            severity,
		Reason:              reason.Name,
		RetryClassification: reason.RetryClassification,
	}
}

var _ error = &ReadyConditionImpactingError{}

func AsReadyConditionImpactingError(err error) (*ReadyConditionImpactingError, bool) {
	var typedErr *ReadyConditionImpactingError
	if errors.As(err, &typedErr) {
		return typedErr, true
	}

	return nil, false
}

func (e *ReadyConditionImpactingError) WithRetryClassification(classification RetryClassification) *ReadyConditionImpactingError {
	e.RetryClassification = classification
	return e
}

func (e *ReadyConditionImpactingError) Error() string {
	return fmt.Sprintf("Reason: %s, Severity: %s, RetryClassification: %s, Cause: %s",
		e.Reason,
		e.Severity,
		e.RetryClassification,
		e.cause.Error())
}

func (e *ReadyConditionImpactingError) Is(err error) bool {
	var typedErr *ReadyConditionImpactingError
	if errors.As(err, &typedErr) {
		return e.Severity == typedErr.Severity && e.Reason == typedErr.Reason
	}
	return false
}

func (e *ReadyConditionImpactingError) Cause() error {
	return e.cause
}

// This was adapted from the function in errors
func (e *ReadyConditionImpactingError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			n, _ := fmt.Fprintf(s, "%s", e.Cause())
			if n > 0 {
				_, _ = fmt.Fprintf(s, "\n")
			}
			_, _ = io.WriteString(s, e.Error())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}
