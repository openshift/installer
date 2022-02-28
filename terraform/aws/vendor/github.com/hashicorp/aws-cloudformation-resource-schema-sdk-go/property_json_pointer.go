package cfschema

import "strings"

const (
	JsonPointerReferenceTokenSeparator  = "/"
	PropertiesJsonPointerReferenceToken = "properties"
	PropertiesJsonPointerPrefix         = JsonPointerReferenceTokenSeparator + PropertiesJsonPointerReferenceToken
)

// PropertyJsonPointer is a simplistic RFC 6901 handler for properties JSON Pointers.
type PropertyJsonPointer string

// EqualsPath returns true if all path parts match.
//
// This automatically handles stripping the /properties prefix.
func (p *PropertyJsonPointer) EqualsPath(other []string) bool {
	if p == nil || *p == "" {
		return false
	}

	path := p.Path()

	if len(path) != len(other) {
		return false
	}

	for i, segment := range path {
		if segment != other[i] {
			return false
		}
	}

	return true
}

// EqualsStringPath returns true if the path string matches.
//
// This automatically handles stripping the /properties prefix.
func (p *PropertyJsonPointer) EqualsStringPath(path string) bool {
	if p == nil || *p == "" {
		return false
	}

	trimmedPath := strings.TrimPrefix(string(*p), PropertiesJsonPointerPrefix)

	return trimmedPath == path
}

// Path returns the path parts.
//
// This automatically handles stripping the /properties path part.
func (p *PropertyJsonPointer) Path() []string {
	if p == nil {
		return nil
	}

	pathParts := strings.Split(strings.TrimPrefix(string(*p), PropertiesJsonPointerPrefix+JsonPointerReferenceTokenSeparator), JsonPointerReferenceTokenSeparator)

	return pathParts
}

// String returns a string representation of the PropertyJsonPointer.
func (p *PropertyJsonPointer) String() string {
	if p == nil {
		return ""
	}

	return string(*p)
}
