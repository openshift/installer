package defaults

import (
	"github.com/openshift/installer/pkg/types/aws"
)

var (
	defaultMachineClass = map[string]string{
		"eu-north-1":    "m5",
		"eu-west-3":     "m5",
		"us-gov-east-1": "m5",
	}
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *aws.Platform) {
}

// InstanceClass returns the instance "class" we should use for a given
// region. We prefer m4 if available (more EBS volumes per node) but will use
// m5 in regions that don't have m4.
func InstanceClass(region string) string {
	if class, ok := defaultMachineClass[region]; ok {
		return class
	}
	return "m4"
}
