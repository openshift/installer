/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package errwrap

type hideError struct {
	err error
	msg string
}

var _ error = &hideError{}

// Hide wraps an error with a message, hiding the original error.
// This is useful when you want to provide a user-friendly message without exposing the underlying error details.
// The original error can still be accessed via the Unwrap method (so errors.Is and errors.As work as expected).
func Hide(err error, msg string) error {
	return &hideError{
		err: err,
		msg: msg,
	}
}

func (w *hideError) Error() string {
	return w.msg
}

func (w *hideError) Unwrap() error {
	return w.err
}
