package cfschema

import (
	"regexp"
)

var (
	sanitizePattern           = regexp.MustCompile(`(?m)^(\s+"pattern"\s*:\s*)".*"`)
	sanitizePatternProperties = regexp.MustCompile(`(?m)^(\s+"patternProperties"\s*:\s*{\s*)".*?"`)
)

// Sanitize returns a sanitized copy of the specified JSON Schema document.
// The santized copy rewrites all pattern and patternProperty regexes to the empty string,
// working around any problems with JSON Schema regex validation.
func Sanitize(document string) string {
	document = sanitizePattern.ReplaceAllString(document, `$1""`)
	document = sanitizePatternProperties.ReplaceAllString(document, `$1""`)

	return document
}
