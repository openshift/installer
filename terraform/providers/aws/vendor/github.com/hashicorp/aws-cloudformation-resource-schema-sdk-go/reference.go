package cfschema

import (
	"fmt"
	"strings"
)

const (
	ReferenceAnchor          = "#"
	ReferenceSeparator       = "/"
	ReferenceTypeDefinitions = "definitions"
	ReferenceTypeProperties  = "properties"
)

// Reference is an internal implementation for RFC 6901 JSON Pointer values.
type Reference string

// Field returns the JSON Pointer string path after the type.
func (r Reference) Field() (string, error) {
	referenceParts := strings.Split(strings.TrimPrefix(string(r), ReferenceAnchor), ReferenceSeparator)

	if got, expected := len(referenceParts), 3; got != expected {
		return "", fmt.Errorf("invalid Reference (%s). Expected %d parts, got %d", r, expected, got)
	}

	return referenceParts[2], nil
}

// String returns the string representation of a Reference.
func (r Reference) String() string {
	return string(r)
}

// Type returns the first path part of the JSON Pointer.
//
// In CloudFormation Resources, this should be definitions or properties.
func (r Reference) Type() (string, error) {
	referenceParts := strings.Split(strings.TrimPrefix(string(r), ReferenceAnchor), ReferenceSeparator)

	if got, expected := len(referenceParts), 3; got != expected {
		return "", fmt.Errorf("invalid Reference (%s). Expected %d parts, got %d", r, expected, got)
	}

	return referenceParts[1], nil
}
