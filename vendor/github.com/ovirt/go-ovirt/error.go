package ovirtsdk

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type baseError struct {
	Code int
	Msg  string
}

func (b *baseError) Error() string {
	return b.Msg
}

// AuthError indicates that an authentiation or authorization
// problem happened, like incorrect user name, incorrect password, or
// missing permissions.
type AuthError struct {
	baseError
}

// NotFoundError indicates that an object can't be found.
type NotFoundError struct {
	baseError
}

// CheckFault procoesses error parsing and returns it back
func CheckFault(response *http.Response) error {
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response, reason: %s", err.Error())
	}
	// Process empty response body
	if len(resBytes) == 0 {
		return BuildError(response, nil)
	}

	reader := NewXMLReader(resBytes)
	fault, err := XMLFaultReadOne(reader, nil, "")
	if err != nil {
		// If the XML is not a <fault>, just return nil
		if err, ok := err.(XMLTagNotMatchError); ok {
			return err
		}
		return err
	}
	if fault != nil || response.StatusCode >= 400 {
		return BuildError(response, fault)
	}
	return errors.New("unknown error")
}

// CheckAction checks if response contains an Action instance
func CheckAction(response *http.Response) (*Action, error) {
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response, reason: %s", err.Error())
	}
	// Process empty response body
	if len(resBytes) == 0 {
		return nil, BuildError(response, nil)
	}

	faultreader := NewXMLReader(resBytes)
	fault, err := XMLFaultReadOne(faultreader, nil, "")
	if err != nil {
		// If the tag mismatches, return the err
		if _, ok := err.(XMLTagNotMatchError); !ok {
			return nil, err
		}
	}
	if fault != nil {
		return nil, BuildError(response, fault)
	}

	actionreader := NewXMLReader(resBytes)
	action, err := XMLActionReadOne(actionreader, nil, "")
	if err != nil {
		// If the tag mismatches, return the err
		if _, ok := err.(XMLTagNotMatchError); !ok {
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
