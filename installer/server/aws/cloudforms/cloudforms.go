package cloudforms

type VPCSubnet struct {
	// Identifier of the subnet if already existing
	ID string `json:"id"`
	// Logical name for this subnet
	// ignored if existing
	Name string `json:"name"`
	// Availability zone for this subnet
	// Max one subnet per availability zone
	AvailabilityZone string `json:"availabilityZone"`
	// CIDR for this subnet
	// must be disjoint from other subnets
	// must be contained by VPC CIDR
	InstanceCIDR string `json:"instanceCIDR"`
}
