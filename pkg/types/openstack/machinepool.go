package openstack

import machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"

// MachinePool stores the configuration for a machine pool installed
// on OpenStack.
type MachinePool struct {
	// FlavorName defines the OpenStack Nova flavor.
	// eg. m1.large
	FlavorName string `json:"type"`

	// RootVolume defines the root volume for instances in the machine pool.
	// The instances use ephemeral disks if not set.
	// +optional
	RootVolume *RootVolume `json:"rootVolume,omitempty"`

	// AdditionalNetworkIDs contains IDs of additional networks for machines,
	// where each ID is presented in UUID v4 format.
	// Allowed address pairs won't be created for the additional networks.
	// +optional
	AdditionalNetworkIDs []string `json:"additionalNetworkIDs,omitempty"`

	// AdditionalSecurityGroupIDs contains IDs of additional security groups for machines,
	// where each ID is presented in UUID v4 format.
	// +optional
	AdditionalSecurityGroupIDs []string `json:"additionalSecurityGroupIDs,omitempty"`

	// ServerGroupPolicy will be used to create the Server Group that will contain all the machines of this MachinePool.
	// Defaults to "soft-anti-affinity".
	ServerGroupPolicy ServerGroupPolicy `json:"serverGroupPolicy,omitempty"`

	// Zones is the list of availability zones where the instances should be deployed.
	// If no zones are provided, all instances will be deployed on OpenStack Nova default availability zone
	// +optional
	Zones []string `json:"zones,omitempty"`

	// FailureDomains is a Tech Preview feature for resource placement. It
	// is incompatible with zones and rootVolume zones.
	FailureDomains []FailureDomain `json:"failureDomains,omitempty"`
}

// Set sets the values from `required` to `o`.
func (o *MachinePool) Set(required *MachinePool) {
	if required == nil || o == nil {
		return
	}

	if required.FlavorName != "" {
		o.FlavorName = required.FlavorName
	}

	if required.RootVolume != nil {
		if o.RootVolume == nil {
			o.RootVolume = new(RootVolume)
		}
		o.RootVolume.Size = required.RootVolume.Size
		o.RootVolume.Type = required.RootVolume.Type
		if len(required.RootVolume.Zones) > 0 {
			o.RootVolume.Zones = required.RootVolume.Zones
		}
	}

	if required.AdditionalNetworkIDs != nil {
		o.AdditionalNetworkIDs = append(required.AdditionalNetworkIDs[:0:0], required.AdditionalNetworkIDs...)
	}

	if required.AdditionalSecurityGroupIDs != nil {
		o.AdditionalSecurityGroupIDs = append(required.AdditionalSecurityGroupIDs[:0:0], required.AdditionalSecurityGroupIDs...)
	}

	if required.ServerGroupPolicy != "" {
		o.ServerGroupPolicy = required.ServerGroupPolicy
	}

	if len(required.Zones) > 0 {
		o.Zones = required.Zones
	}

	if len(required.FailureDomains) > 0 {
		o.FailureDomains = required.FailureDomains
	}
}

// RootVolume defines the storage for an instance.
type RootVolume struct {
	// Size defines the size of the volume in gibibytes (GiB).
	// Required
	Size int `json:"size"`
	// Type defines the type of the volume.
	// Required
	Type string `json:"type"`

	// Zones is the list of availability zones where the root volumes should be deployed.
	// If no zones are provided, all instances will be deployed on OpenStack Cinder default availability zone
	// +optional
	Zones []string `json:"zones,omitempty"`
}

// -=**
// FailureDomain and its types are part of a Tech Preview feature, and may change before they're moved to openshift/api
// **=-

// FailureDomain defines a set of correlated fault domains across different
// OpenStack services: compute, storage, and network.
type FailureDomain struct {
	// ComputeAvailabilityZone is the name of a valid nova availability zone. The server will be created in this availability zone.
	// +optional
	ComputeAvailabilityZone string `json:"computeAvailabilityZone"`

	// StorageAvailabilityZone is the name of a valid cinder availability
	// zone. This will be the availability zone of the root volume if one is
	// specified.
	// +optional
	StorageAvailabilityZone string `json:"storageAvailabilityZone"`

	// Ports defines a set of port targets which can be referenced by machines in this failure domain.
	//
	// +optional
	PortTargets []NamedPortTarget `json:"portTargets"`
}

// NamedPortTarget includes an ID to characterize a PortTarget with its
// intended purpose. If ID is set to "control-plane", then the PortTarget will
// replace the default cluster primary network (or the machinesSubnet if
// defined) as the first network for the machine.
type NamedPortTarget struct {
	ID         string `json:"id"`
	PortTarget `json:",inline"`
}

// PortTarget defines, directly or indirectly, one or more subnets where to attach a port.
type PortTarget struct {
	// Network is a query for an openstack network that the port will be created or discovered on.
	// This will fail if the query returns more than one network.
	Network NetworkFilter `json:"network,omitempty"`
	// Specify pairs of subnet and/or IP address. These should be subnets of the network with the given NetworkID.
	FixedIPs []FixedIP `json:"fixedIPs,omitempty"`
}

// NetworkFilter defines a network either by name or by ID.
type NetworkFilter struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// FixedIP defines a subnet.
type FixedIP struct {
	// subnetID specifies the ID of the subnet where the fixed IP will be allocated.
	Subnet machinev1alpha1.SubnetFilter `json:"subnet"`
}
