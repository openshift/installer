package openstack

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
		o.RootVolume.DeprecatedType = required.RootVolume.DeprecatedType

		if required.RootVolume.DeprecatedType != "" {
			o.RootVolume.DeprecatedType = ""
			o.RootVolume.Types = []string{required.RootVolume.DeprecatedType}
		} else if len(required.RootVolume.Types) > 0 {
			o.RootVolume.Types = required.RootVolume.Types
		}

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
}

// RootVolume defines the storage for an instance.
type RootVolume struct {
	// Size defines the size of the volume in gibibytes (GiB).
	// Required
	Size int `json:"size"`
	// Type defines the type of the volume.
	// Deprecated: Use Types instead.
	// +optional
	DeprecatedType string `json:"type,omitempty"`

	// Types is the list of the volume types of the root volumes.
	// This is mutually exclusive with Type.
	// +required
	Types []string `json:"types"`

	// Zones is the list of availability zones where the root volumes should be deployed.
	// If no zones are provided, all instances will be deployed on OpenStack Cinder default availability zone
	// +optional
	Zones []string `json:"zones,omitempty"`
}

// PortTarget defines, directly or indirectly, one or more subnets where to attach a port.
type PortTarget struct {
	// Network is a query for an openstack network that the port will be discovered on.
	// This will fail if the query returns more than one network.
	Network NetworkFilter `json:"network,omitempty"`
	// Specify subnets of the network where control plane port will be discovered.
	FixedIPs []FixedIP `json:"fixedIPs"`
}

// NetworkFilter defines a network by name and/or ID.
type NetworkFilter struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// FixedIP identifies a subnet defined by a subnet filter.
type FixedIP struct {
	Subnet SubnetFilter `json:"subnet"`
}

// SubnetFilter defines a subnet by ID and/or name.
type SubnetFilter struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
