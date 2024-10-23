package output

import "strings"

const (
	Yes        = "Yes"
	No         = "No"
	EmptySlice = ""
)

// PrintBool returns a prettified version of a boolean. "Yes" for true, or "No" for false
func PrintBool(b bool) string {
	if b {
		return Yes
	}
	return No
}

func PrintStringSlice(in []string) string {
	if len(in) == 0 {
		return EmptySlice
	}
	return strings.Join(in, ", ")
}
