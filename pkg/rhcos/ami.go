package rhcos

import (
	"fmt"
)

const (
	// DefaultChannel is the default RHCOS channel for the cluster.
	DefaultChannel = "tested"
)

// AMI calculates a Red Hat CoreOS AMI.
func AMI(channel, region string) (ami string, err error) {
	if channel != DefaultChannel {
		return "", fmt.Errorf("channel %q is not yet supported", channel)
	}

	if region != "us-east-1" {
		return "", fmt.Errorf("region %q is not yet supported", region)
	}

	return "ami-09b1c714bf5aaa535", nil
}
