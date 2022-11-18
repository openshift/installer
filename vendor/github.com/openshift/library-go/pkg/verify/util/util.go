package util

import (
	"fmt"
	"strings"
)

// DigestToKeyPrefix changes digest to use the provided newDivider in place of ':',
// {algo}{newDivider}{hash} instead of {algo}:{hash}, because colons are not allowed
// in various places such as ConfigMap keys.
func DigestToKeyPrefix(digest string, newDivider string) (string, error) {
	parts := strings.SplitN(digest, ":", 3)
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return "", fmt.Errorf("the provided digest must be of the form ALGO:HASH")
	}
	algo, hash := parts[0], parts[1]
	return fmt.Sprintf("%s%s%s", algo, newDivider, hash), nil
}
