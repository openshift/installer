package defaults

import "github.com/openshift/installer/pkg/types/ibmcloud"

// Apply sets values from the default machine platform to the machinepool.
func Apply(defaultMachinePlatform, machinePool *ibmcloud.MachinePool) {
	tempMP := &ibmcloud.MachinePool{}
	tempMP.Set(defaultMachinePlatform)
	tempMP.Set(machinePool)
	machinePool.Set(tempMP)
}
