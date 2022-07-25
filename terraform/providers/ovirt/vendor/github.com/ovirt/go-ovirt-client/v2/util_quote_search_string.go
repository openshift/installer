package ovirtclient

import (
	"fmt"
	"strings"
)

func quoteSearchString(text string) (string, error) {
	if strings.Contains(text, "\"") {
		return "", newError(EBadArgument, "quotes are not allowed in search strings")
	}
	if strings.Contains(text, "*") {
		return "", newError(EBadArgument, "wildcards are not allowed in search strings")
	}
	return fmt.Sprintf("\"%s\"", text), nil
}
