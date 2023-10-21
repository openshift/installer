package validate

import (
	"fmt"
	"regexp"
)

func LabDescription(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^.{1,500}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only be up to 500 characters in length", k))
	}

	return warnings, errors
}
