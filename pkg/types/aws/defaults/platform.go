package defaults

import (
	"github.com/openshift/installer/pkg/types/aws"
)

var (
	defaultMachineClass = map[string][]string{
		"ap-east-1":      {"m5", "m4"},
		"ap-northeast-2": {"m5", "m4"},
		"eu-north-1":     {"m5", "m4"},
		"eu-west-3":      {"m5", "m4"},
		"me-south-1":     {"m5", "m4"},
		"us-gov-east-1":  {"m5", "m4"},
		"us-west-2":      {"m5", "m4"},
	}
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *aws.Platform) {
}

// InstanceClass returns the instance "class" we should use for a given
// region. We prefer m4 if available (more EBS volumes per node) but will use
// m5 in regions that don't have m4.
func InstanceClass(region string) string {
	if classes, ok := defaultMachineClass[region]; ok {
		return classes[0]
	}
	return "m4"
}

// InstanceClasses returns a list of instance "class", in decreasing priority order, which we should use for a given
// region. We prefer m4 if available (more EBS volumes per node) but will use
// m5 in regions that don't have m4.
func InstanceClasses(region string) []string {
	if classes, ok := defaultMachineClass[region]; ok {
		return classes
	}
	return []string{"m4", "m5"}
}
