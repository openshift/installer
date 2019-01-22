package google

import (
	"fmt"
)

// Zones retrieves a list of zones for the given region.
func Zones(region string) ([]string, error) {
	switch region {
	case "us-central1":
		return []string{"us-central1-a", "us-central1-b", "us-central1-c"}, nil
	default:
		return nil, fmt.Errorf("cannot fetch zones: not implemented")
	}
}
