package defaults

import (
	"github.com/openshift/installer/pkg/types/aws"
)

var (
	defaultMachineClass = map[string][]string{
		// Example region default machine class override:
		// "ap-east-1":      {"m5", "m4"},
	}
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *aws.Platform) {
}

// InstanceClass returns the instance "class" we should use for a given
// region. Default is m5 unless a region override is defined in defaultMachineClass.
func InstanceClass(region string) string {
	if classes, ok := defaultMachineClass[region]; ok {
		return classes[0]
	}
	return "m5"
}

// InstanceClasses returns a list of instance "class", in decreasing priority order, which we should use for a given
// region. Default is m5 then m4 unless a region override is defined in defaultMachineClass.
func InstanceClasses(region string) []string {
	if classes, ok := defaultMachineClass[region]; ok {
		return classes
	}
	return []string{"m5", "m4"}
}
