package defaults

import (
	"fmt"

	"github.com/openshift/installer/pkg/types/azure"
)

// BootstrapInstanceType sets the defaults for bootstrap instances.
// Minimum requirements are 4 CPU's, 16GiB of ram, and 120GiB storage.
// D4s v3 gives us 4 CPU's, 16GiB ram and 32GiB of temporary storage.
// DS4_v2 gives us 8 CPUs, 28GiB ram, and 56GiB of temporary storage.
func BootstrapInstanceType(cloud azure.CloudEnvironment, region string) string {
	instanceClass := getInstanceClass(region)
	size := "D4s_v3"
	if cloud == azure.StackCloud {
		size = "DS4_v2"
	}
	return instanceType(instanceClass, size)
}

// ControlPlaneInstanceType sets the defaults for control plane instances.
// Minimum requirements are 4 CPU's, 16GiB of ram, and 120GiB storage.
// D8s v3 gives us 8 CPU's, 32GiB ram and 64GiB of temporary storage
// This extra bump is done to prevent etcd from overloading
// DS4_v2 gives us 8 CPUs, 28GiB ram, and 56GiB of temporary storage.
func ControlPlaneInstanceType(cloud azure.CloudEnvironment, region string) string {
	instanceClass := getInstanceClass(region)
	size := "D4s_v3"
	if cloud == azure.StackCloud {
		size = "DS4_v2"
	}
	return instanceType(instanceClass, size)
}

// ComputeInstanceType sets the defaults for compute instances.
// Minimum requirements are 2 CPU's, 8GiB of ram, and 120GiB storage.
// D4s v3 gives us 2 CPU's, 8GiB ram and 16GiB of temporary storage
// DS3_v2 gives us 4 CPUs, 14GiB ram, and 28GiB of temporary storage.
func ComputeInstanceType(cloud azure.CloudEnvironment, region string) string {
	instanceClass := getInstanceClass(region)
	size := "D4s_v3"
	if cloud == azure.StackCloud {
		size = "DS4_v2"
	}
	return instanceType(instanceClass, size)
}

func instanceType(class, size string) string {
	return fmt.Sprintf("%s_%s", class, size)
}
