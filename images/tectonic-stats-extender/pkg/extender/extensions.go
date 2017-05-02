package extender

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Extensions is a map of key value pairs of extensions to report.
type Extensions map[string]string

// String returns a string representation of the extensions map.
func (e Extensions) String() string {
	var s []string
	for k, v := range e {
		s = append(s, fmt.Sprintf("%s:%s", k, v))
	}
	sort.Strings(s)
	return strings.Join(s, ", ")
}

// Set adds a new extension to the extensions map. This method splits the flag
// string on the first ":" character, using the first string as the key and the
// second string as the value. If the flag does not split on ":" then an error
// is raised.
func (e Extensions) Set(value string) error {
	s := strings.SplitN(value, ":", 2)

	if len(s) != 2 {
		return errors.New("extension must be of the form <key>:<value>")
	}

	e[s[0]] = s[1]
	return nil
}
