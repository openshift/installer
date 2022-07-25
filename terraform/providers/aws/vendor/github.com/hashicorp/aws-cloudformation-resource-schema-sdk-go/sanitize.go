package cfschema

import (
	"bufio"
	"encoding/json"
	"regexp"
	"strings"
)

var (
	sanitizePattern           = regexp.MustCompile(`^(\s*"pattern"\s*:\s*)"(.*)"(\s*,?\s*)$`)
	sanitizePatternProperties = regexp.MustCompile(`(?m)^(\s+"patternProperties"\s*:\s*{\s*)".*?"`)
)

// Sanitize returns a sanitized copy of the specified JSON Schema document.
// The sanitized copy works around any problems with JSON Schema regex validation by
//  - Rewriting all patternProperty regexes to the empty string (the regex is never used anyway)
//  - Rewriting all unsupported (valid for ECMA-262 but not for Go) pattern regexes to the empty string
func Sanitize(document string) (string, error) {
	document = sanitizePatternProperties.ReplaceAllString(document, `$1""`)

	var sb strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(document))
	for scanner.Scan() {
		line := scanner.Text()

		if v := sanitizePattern.FindStringSubmatch(line); len(v) == 4 {
			if expr := v[2]; expr != "" && !isSupportedRegexp(expr) {
				line = v[1] + "\"\"" + v[3]
			}
		}
		if _, err := sb.WriteString(line); err != nil {
			return "", err
		}
		if _, err := sb.WriteString("\n"); err != nil {
			return "", err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func isSupportedRegexp(expr string) bool {
	// github.com/xeipuuv/gojsonschema attempts to compile the regex after it has been unmarshaled from JSON.
	var v string
	b := []byte("\"" + expr + "\"")

	if err := json.Unmarshal(b, &v); err != nil {
		return false
	}

	_, err := regexp.Compile(v)
	return err == nil
}
