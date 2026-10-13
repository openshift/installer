package diagnostics

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Err wraps diagnostics information for an error.
// Err allows providing information like source, reason and message
// that provides a much better user error reporting capability.
type Err struct {
	Orig error

	// Source defines with entity is generating the error.
	// It allows passing along information about where the error is being
	// generated from. for example, the Asset.
	Source string

	// Reason is a CamelCase string that summarizes the error in one word.
	// It allows easy catgeorizations of known errors.
	Reason string

	// Message is free-form strings which provides important details or
	// diagnostics for the error. When writing messages, make sure to keep in mind
	// that the audience for message is end-users who might not be experts.
	Message string
}

// Unwrap allows the error to be unwrapped.
func (e *Err) Unwrap() error { return e.Orig }

// Error returns a string representation of the Err. The returned value
// is expected to be a single value.
// The format of the error string returned is,
// `error(<Reason>) from <Source>: <Message>: <Cause of Orig>`
func (e *Err) Error() string {
	buf := &bytes.Buffer{}
	if len(e.Source) > 0 {
		fmt.Fprintf(buf, "error(%s) from %s", e.Reason, e.Source)
	} else {
		fmt.Fprintf(buf, "error(%s)", e.Reason)
	}
	if msg := strings.TrimSpace(e.Message); len(msg) > 0 {
		msg = breakre.ReplaceAllString(msg, " ")
		fmt.Fprintf(buf, ": %s", msg)
	}
	if c := errors.Cause(e.Orig); c != nil {
		fmt.Fprintf(buf, ": %s", errors.Cause(e.Orig))
	}
	return buf.String()
}

// Print prints the Err to Writer in a way that is more verbose and
// sectionalized.
// The output looks like:
// Error from <Source>:
// Reason: <reason>
//
// Message:
// <Message>
//
// Original:
// <Orig>
func (e *Err) Print(w io.Writer) {
	fmt.Fprintf(w, "Error from %q\n", e.Source)
	fmt.Fprintf(w, "Reason: %s\n", e.Reason)
	if len(e.Message) > 0 {
		fmt.Fprintf(w, "\nMessage:\n")
		fmt.Fprintln(w, e.Message)
	}
	fmt.Fprintf(w, "\nOriginal error:\n")
	fmt.Fprintln(w, e.Orig)
}

var breakre = regexp.MustCompile(`\r?\n`)
