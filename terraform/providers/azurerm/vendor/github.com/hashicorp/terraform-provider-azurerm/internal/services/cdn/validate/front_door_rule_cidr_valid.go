package validate

import (
	"fmt"
	"net"
)

// Evaluates if the passed CIDR is a valid IPv4 or IPv6 CIDR or not.
func FrontDoorRuleCidrIsValid(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, _, err := net.ParseCIDR(v); err != nil {
		errors = append(errors, err)
	}

	return
}
