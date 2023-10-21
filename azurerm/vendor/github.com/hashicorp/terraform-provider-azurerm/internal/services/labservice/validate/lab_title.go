package validate

import (
	"fmt"
	"regexp"
)

func LabTitle(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^.{1,100}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only be up to 100 characters in length", k))
	}

	return warnings, errors
}
