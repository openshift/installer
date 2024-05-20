package cfschema

import (
	"strings"
)

// PropertyTransform represents property transform values.
type PropertyTransform map[string]string

// Value returns the value for a specified property path.
func (p PropertyTransform) Value(path []string) (string, bool) {
	pa := buildPath(path)

	if value, ok := pathClean(p)[pa]; ok {
		return value, true
	}

	return "", false
}

func buildPath(path []string) string {
	if len(path) == 1 {
		return path[0]
	}

	return strings.Join(path, "/")
}

func pathClean(m map[string]string) map[string]string {
	vals := make(map[string]string)

	for k, v := range m {
		vals[strings.TrimPrefix(string(k), PropertiesJsonPointerPrefix+JsonPointerReferenceTokenSeparator)] = v
	}

	return vals
}
