package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

const (
	// defaultRootVolumeSize is the default roote volume size in GB.
	defaultRootVolumeSize = 120
)

// SetMachinePoolDefaults sets the defaults for the platform.
func SetMachinePoolDefaults(pool *aws.MachinePool, role string) {
	if pool == nil {
		return
	}

	// Set the default volume type for machine pool.
	// The current default is gp3 for control plane and worker pool.
	// For edge pool, the current default is gp2.
	// See: https://github.com/openshift/installer/blob/fd5a518e4951510b82705eee8184b3dd4f2723b2/pkg/asset/machines/worker.go#L102-L117
	if pool.EC2RootVolume.Type == "" {
		defaultEBSType := aws.VolumeTypeGp3
		if role == types.MachinePoolEdgeRoleName {
			defaultEBSType = aws.VolumeTypeGp2
		}
		pool.EC2RootVolume.Type = defaultEBSType
	}

	if pool.EC2RootVolume.Size == 0 {
		pool.EC2RootVolume.Size = defaultRootVolumeSize
	}
}

// Apply sets values from the default machine platform to the machinepool.
func Apply(defaultMachinePlatform, machinePool *aws.MachinePool) {
	// Construct a temporary machine pool so we can set the
	// defaults first, without overwriting the pool-sepcific values,
	// which have precedence.
	tempMP := &aws.MachinePool{}
	tempMP.Set(defaultMachinePlatform)
	tempMP.Set(machinePool)
	machinePool.Set(tempMP)
}
