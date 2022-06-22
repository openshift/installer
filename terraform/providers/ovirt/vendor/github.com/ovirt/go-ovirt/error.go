package ovirtsdk

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

type baseError struct {
	// Code contains the HTTP status code that caused this error.
	Code int
	// Msg contains the text message that should be printed to the user.
	Msg string
}

// Error returns the error string.
func (b *baseError) Error() string {
	return b.Msg
}

// AuthError indicates that an authentication or authorization
// problem happened, like incorrect user name, incorrect password, or
// missing permissions.
type AuthError struct {
	baseError
}

// NotFoundError indicates that an object can't be found.
type NotFoundError struct {
	baseError
}

// ResponseParseError indicates that the response from the oVirt Engine could not be parsed.
type ResponseParseError struct {
	baseError

	cause error
	body  []byte
}

// Unwrap returns the root cause of this error.
func (r *ResponseParseError) Unwrap() error {
	return r.cause
}

// Body returns the HTTP response body that caused the parse error.
func (r *ResponseParseError) Body() []byte {
	return r.body
}

// CheckFault takes a failed HTTP response (non-200) and extracts the fault from it.
func CheckFault(resBytes []byte, response *http.Response) error {
	// Process empty response body
	if len(resBytes) == 0 {
		return BuildError(response, nil)
	}

	reader := NewXMLReader(resBytes)
	fault, err := XMLFaultReadOne(reader, nil, "")
	if err != nil {
		return &ResponseParseError{
			baseError{
				Code: response.StatusCode,
				Msg: fmt.Sprintf(
					"failed to parse oVirt Engine fault response: %s (%v)",
					resBytes,
					err,
				),
			},
			err,
			resBytes,
		}
	}
	if fault != nil || response.StatusCode >= 400 {
		return BuildError(response, fault)
	}
	return errors.New("unknown error")
}

// CheckAction checks if response contains an Action instance
func CheckAction(resBytes []byte, response *http.Response) (*Action, error) {
	// Process empty response body
	if len(resBytes) == 0 {
		return nil, BuildError(response, nil)
	}
	var tagNotMatchError XMLTagNotMatchError

	faultreader := NewXMLReader(resBytes)
	fault, err := XMLFaultReadOne(faultreader, nil, "")
	if err != nil {
		// If the tag mismatches, return the err
		if !errors.As(err, &tagNotMatchError) {
			return nil, &ResponseParseError{
				baseError{
					Code: response.StatusCode,
					Msg: fmt.Sprintf(
						"failed to parse oVirt Engine response: %s (%v)",
						resBytes,
						err,
					),
				},
				err,
				resBytes,
			}
		}
	}
	if fault != nil {
		return nil, BuildError(response, fault)
	}

	actionreader := NewXMLReader(resBytes)
	action, err := XMLActionReadOne(actionreader, nil, "")
	if err != nil {
		// If the tag mismatches, return the err
		if errors.As(err, &tagNotMatchError) {
			return nil, err
		}
	}
	if action != nil {
		if afault, ok := action.Fault(); ok {
			return nil, BuildError(response, afault)
		}
		return action, nil
	}
	return nil, nil
}

// BuildError constructs error
func BuildError(response *http.Response, fault *Fault) error {
	var buffer bytes.Buffer
	if fault != nil {
		if reason, ok := fault.Reason(); ok {
			if buffer.Len() > 0 {
				buffer.WriteString(" ")
			}
			buffer.WriteString(fmt.Sprintf("Fault reason is \"%s\".", reason))
		}
		if detail, ok := fault.Detail(); ok {
			if buffer.Len() > 0 {
				buffer.WriteString(" ")
			}
			buffer.WriteString(fmt.Sprintf("Fault detail is \"%s\".", detail))
		}
	}
	if response != nil {
		if buffer.Len() > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(fmt.Sprintf("HTTP response code is \"%d\".", response.StatusCode))
		buffer.WriteString(" ")
		buffer.WriteString(fmt.Sprintf("HTTP response message is \"%s\".", response.Status))

		if Contains(response.StatusCode, []int{401, 403}) {
			return &AuthError{
				baseError{
					response.StatusCode,
					buffer.String(),
				},
			}
		} else if response.StatusCode == 404 {
			return &NotFoundError{
				baseError{
					response.StatusCode,
					buffer.String(),
				},
			}
		}
	}

	return errors.New(buffer.String())
}
