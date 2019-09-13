package defaults

import (
	"fmt"
)

// BootstrapInstanceType sets the defaults for bootstrap instances.
// Minimum requirements are 4 CPU's, 16GiB of ram, and 120GiB storage.
// D4s v3 gives us 4 CPU's, 16GiB ram and 32GiB of temporary storage
func BootstrapInstanceType(region string) string {
	instanceClass := getInstanceClass(region)
	return fmt.Sprintf("%s_D4s_v3", instanceClass)
}

// ControlPlaneInstanceType sets the defaults for control plane instances.
// Minimum requirements are 4 CPU's, 16GiB of ram, and 120GiB storage.
// D4s v3 gives us 4 CPU's, 16GiB ram and 32GiB of temporary storage
func ControlPlaneInstanceType(region string) string {
	instanceClass := getInstanceClass(region)
	return fmt.Sprintf("%s_D16s_v3", instanceClass)
}

// ComputeInstanceType sets the defaults for compute instances.
// Minimum requirements are 2 CPU's, 8GiB of ram, and 120GiB storage.
// D4s v3 gives us 2 CPU's, 8GiB ram and 16GiB of temporary storage
func ComputeInstanceType(region string) string {
	instanceClass := getInstanceClass(region)
	return fmt.Sprintf("%s_D2s_v3", instanceClass)
}
