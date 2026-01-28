package defaults

import (
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/network"
)

var (
	// Overrides
	defaultMachineClass = map[string]string{}

	// AzurestackMinimumDiskSize is the minimum disk size value for azurestack.
	AzurestackMinimumDiskSize int32 = 128
	// AzurestackMaximumDiskSize is the maximum disk size value for azurestack.
	AzurestackMaximumDiskSize int32 = 1023
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *azure.Platform) {
	if p.CloudName == "" {
		p.CloudName = azure.PublicCloud
	}
	if p.OutboundType == "" {
		p.OutboundType = azure.LoadbalancerOutboundType
	}
	if p.IPFamily == "" {
		p.IPFamily = network.IPv4
		logrus.Infof("ipFamily is not specified in install-config; defaulting to %q", network.IPv4)
	}
}

// Apply sets values from the default machine platform to the machinepool.
func Apply(defaultMachinePlatform, machinePool *azure.MachinePool) {
	// Construct a temporary machine pool so we can set the
	// defaults first, without overwriting the pool-sepcific values,
	// which have precedence.
	tempMP := &azure.MachinePool{}
	tempMP.Set(defaultMachinePlatform)
	tempMP.Set(machinePool)
	machinePool.Set(tempMP)
}

// getInstanceClass returns the instance "class" we should use for a given region.
func getInstanceClass(region string) string {
	if class, ok := defaultMachineClass[region]; ok {
		return class
	}

	return "Standard"
}
