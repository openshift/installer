/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

// OpenStackMachineTemplateResource describes the data needed to create a OpenStackMachine from a template.
type OpenStackMachineTemplateResource struct {
	// Spec is the specification of the desired behavior of the machine.
	Spec OpenStackMachineSpec `json:"spec"`
}

type ResourceReference struct {
	// Name is the name of the referenced resource
	Name string `json:"name"`
}

// ImageParam describes a glance image. It can be specified by ID, filter, or a
// reference to an ORC Image.
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type ImageParam struct {
	// ID is the uuid of the image. ID will not be validated before use.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// Filter describes a query for an image. If specified, the combination
	// of name and tags must return a single matching image or an error will
	// be raised.
	// +optional
	Filter *ImageFilter `json:"filter,omitempty"`

	// ImageRef is a reference to an ORC Image in the same namespace as the
	// referring object.
	// +optional
	ImageRef *ResourceReference `json:"imageRef,omitempty"`
}

// ImageFilter describes a query for an image.
// +kubebuilder:validation:MinProperties:=1
type ImageFilter struct {
	// The name of the desired image. If specified, the combination of name and tags must return a single matching image or an error will be raised.
	// +optional
	Name optional.String `json:"name,omitempty"`

	// The tags associated with the desired image. If specified, the combination of name and tags must return a single matching image or an error will be raised.
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`
}

func (f *ImageFilter) IsZero() bool {
	if f == nil {
		return true
	}
	return f.Name == nil && len(f.Tags) == 0
}

type ExternalRouterIPParam struct {
	// The FixedIP in the corresponding subnet
	FixedIP string `json:"fixedIP,omitempty"`
	// The subnet in which the FixedIP is used for the Gateway of this router
	Subnet SubnetParam `json:"subnet"`
}

// NeutronTag represents a tag on a Neutron resource.
// It may not be empty and may not contain commas.
// +kubebuilder:validation:Pattern:="^[^,]+$"
// +kubebuilder:validation:MinLength:=1
type NeutronTag string

type FilterByNeutronTags struct {
	// Tags is a list of tags to filter by. If specified, the resource must
	// have all of the tags specified to be included in the result.
	// +listType=set
	// +optional
	Tags []NeutronTag `json:"tags,omitempty"`

	// TagsAny is a list of tags to filter by. If specified, the resource
	// must have at least one of the tags specified to be included in the
	// result.
	// +listType=set
	// +optional
	TagsAny []NeutronTag `json:"tagsAny,omitempty"`

	// NotTags is a list of tags to filter by. If specified, resources which
	// contain all of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	NotTags []NeutronTag `json:"notTags,omitempty"`

	// NotTagsAny is a list of tags to filter by. If specified, resources
	// which contain any of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	NotTagsAny []NeutronTag `json:"notTagsAny,omitempty"`
}

func (f *FilterByNeutronTags) IsZero() bool {
	return f == nil || (len(f.Tags) == 0 && len(f.TagsAny) == 0 && len(f.NotTags) == 0 && len(f.NotTagsAny) == 0)
}

// SecurityGroupParam specifies an OpenStack security group. It may be specified by ID or filter, but not both.
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type SecurityGroupParam struct {
	// ID is the ID of the security group to use. If ID is provided, the other filters cannot be provided. Must be in UUID format.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// Filter specifies a query to select an OpenStack security group. If provided, cannot be empty.
	// +optional
	Filter *SecurityGroupFilter `json:"filter,omitempty"`
}

// SecurityGroupFilter specifies a query to select an OpenStack security group. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type SecurityGroupFilter struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"projectID,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

func (f *SecurityGroupFilter) IsZero() bool {
	if f == nil {
		return true
	}
	return f.Name == "" &&
		f.Description == "" &&
		f.ProjectID == "" &&
		f.FilterByNeutronTags.IsZero()
}

// NetworkParam specifies an OpenStack network. It may be specified by either ID or Filter, but not both.
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type NetworkParam struct {
	// ID is the ID of the network to use. If ID is provided, the other filters cannot be provided. Must be in UUID format.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// Filter specifies a filter to select an OpenStack network. If provided, cannot be empty.
	// +optional
	Filter *NetworkFilter `json:"filter,omitempty"`
}

// NetworkFilter specifies a query to select an OpenStack network. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type NetworkFilter struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"projectID,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

func (networkFilter *NetworkFilter) IsZero() bool {
	if networkFilter == nil {
		return true
	}
	return networkFilter.Name == "" &&
		networkFilter.Description == "" &&
		networkFilter.ProjectID == "" &&
		networkFilter.FilterByNeutronTags.IsZero()
}

// SubnetParam specifies an OpenStack subnet to use. It may be specified by either ID or filter, but not both.
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type SubnetParam struct {
	// ID is the uuid of the subnet. It will not be validated.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// Filter specifies a filter to select the subnet. It must match exactly one subnet.
	// +optional
	Filter *SubnetFilter `json:"filter,omitempty"`
}

// SubnetFilter specifies a filter to select a subnet. At least one parameter must be specified.
// +kubebuilder:validation:MinProperties:=1
type SubnetFilter struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	ProjectID       string `json:"projectID,omitempty"`
	IPVersion       int    `json:"ipVersion,omitempty"`
	GatewayIP       string `json:"gatewayIP,omitempty"`
	CIDR            string `json:"cidr,omitempty"`
	IPv6AddressMode string `json:"ipv6AddressMode,omitempty"`
	IPv6RAMode      string `json:"ipv6RAMode,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

func (subnetFilter *SubnetFilter) IsZero() bool {
	if subnetFilter == nil {
		return true
	}
	return subnetFilter.Name == "" &&
		subnetFilter.Description == "" &&
		subnetFilter.ProjectID == "" &&
		subnetFilter.IPVersion == 0 &&
		subnetFilter.GatewayIP == "" &&
		subnetFilter.CIDR == "" &&
		subnetFilter.IPv6AddressMode == "" &&
		subnetFilter.IPv6RAMode == "" &&
		subnetFilter.FilterByNeutronTags.IsZero()
}

// RouterParam specifies an OpenStack router to use. It may be specified by either ID or filter, but not both.
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type RouterParam struct {
	// ID is the ID of the router to use. If ID is provided, the other filters cannot be provided. Must be in UUID format.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// Filter specifies a filter to select an OpenStack router. If provided, cannot be empty.
	Filter *RouterFilter `json:"filter,omitempty"`
}

// RouterFilter specifies a query to select an OpenStack router. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type RouterFilter struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"projectID,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

func (f *RouterFilter) IsZero() bool {
	if f == nil {
		return true
	}
	return f.Name == "" &&
		f.Description == "" &&
		f.ProjectID == "" &&
		f.FilterByNeutronTags.IsZero()
}

type SubnetSpec struct {
	// CIDR is representing the IP address range used to create the subnet, e.g. 10.0.0.0/24.
	// This field is required when defining a subnet.
	// +required
	CIDR string `json:"cidr"`

	// DNSNameservers holds a list of DNS server addresses that will be provided when creating
	// the subnet. These addresses need to have the same IP version as CIDR.
	DNSNameservers []string `json:"dnsNameservers,omitempty"`

	// AllocationPools is an array of AllocationPool objects that will be applied to OpenStack Subnet being created.
	// If set, OpenStack will only allocate these IPs for Machines. It will still be possible to create ports from
	// outside of these ranges manually.
	AllocationPools []AllocationPool `json:"allocationPools,omitempty"`
}

type AllocationPool struct {
	// Start represents the start of the AllocationPool, that is the lowest IP of the pool.
	// +required
	Start string `json:"start"`

	// End represents the end of the AlloctionPool, that is the highest IP of the pool.
	// +required
	End string `json:"end"`
}

type PortOpts struct {
	// Network is a query for an openstack network that the port will be created or discovered on.
	// This will fail if the query returns more than one network.
	// +optional
	Network *NetworkParam `json:"network,omitempty"`

	// Description is a human-readable description for the port.
	// +optional
	Description optional.String `json:"description,omitempty"`

	// NameSuffix will be appended to the name of the port if specified. If unspecified, instead the 0-based index of the port in the list is used.
	// +optional
	NameSuffix optional.String `json:"nameSuffix,omitempty"`

	// FixedIPs is a list of pairs of subnet and/or IP address to assign to the port. If specified, these must be subnets of the port's network.
	// +optional
	// +listType=atomic
	FixedIPs []FixedIP `json:"fixedIPs,omitempty"`

	// SecurityGroups is a list of the names, uuids, filters or any combination these of the security groups to assign to the instance.
	// +optional
	// +listType=atomic
	SecurityGroups []SecurityGroupParam `json:"securityGroups,omitempty"`

	// Tags applied to the port (and corresponding trunk, if a trunk is configured.)
	// These tags are applied in addition to the instance's tags, which will also be applied to the port.
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`

	// Trunk specifies whether trunking is enabled at the port level. If not
	// provided the value is inherited from the machine, or false for a
	// bastion host.
	// +optional
	Trunk *bool `json:"trunk,omitempty"`

	ResolvedPortSpecFields `json:",inline"`
}

// ResolvePortSpecFields is a convenience struct containing all fields of a
// PortOpts which don't contain references which need to be resolved, and can
// therefore be shared with ResolvedPortSpec.
type ResolvedPortSpecFields struct {
	// AdminStateUp specifies whether the port should be created in the up (true) or down (false) state. The default is up.
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// MACAddress specifies the MAC address of the port. If not specified, the MAC address will be generated.
	// +optional
	MACAddress optional.String `json:"macAddress,omitempty"`

	// AllowedAddressPairs is a list of address pairs which Neutron will
	// allow the port to send traffic from in addition to the port's
	// addresses. If not specified, the MAC Address will be the MAC Address
	// of the port. Depending on the configuration of Neutron, it may be
	// supported to specify a CIDR instead of a specific IP address.
	// +optional
	AllowedAddressPairs []AddressPair `json:"allowedAddressPairs,omitempty"`

	// HostID specifies the ID of the host where the port resides.
	// +optional
	HostID optional.String `json:"hostID,omitempty"`

	// VNICType specifies the type of vNIC which this port should be
	// attached to. This is used to determine which mechanism driver(s) to
	// be used to bind the port. The valid values are normal, macvtap,
	// direct, baremetal, direct-physical, virtio-forwarder, smart-nic and
	// remote-managed, although these values will not be validated in this
	// API to ensure compatibility with future neutron changes or custom
	// implementations. What type of vNIC is actually available depends on
	// deployments. If not specified, the Neutron default value is used.
	// +optional
	VNICType optional.String `json:"vnicType,omitempty"`

	// Profile is a set of key-value pairs that are used for binding
	// details. We intentionally don't expose this as a map[string]string
	// because we only want to enable the users to set the values of the
	// keys that are known to work in OpenStack Networking API.  See
	// https://docs.openstack.org/api-ref/network/v2/index.html?expanded=create-port-detail#create-port
	// To set profiles, your tenant needs permissions rule:create_port, and
	// rule:create_port:binding:profile
	// +optional
	Profile *BindingProfile `json:"profile,omitempty"`

	// DisablePortSecurity enables or disables the port security when set.
	// When not set, it takes the value of the corresponding field at the network level.
	// +optional
	DisablePortSecurity *bool `json:"disablePortSecurity,omitempty"`

	// PropageteUplinkStatus enables or disables the propagate uplink status on the port.
	// +optional
	PropagateUplinkStatus *bool `json:"propagateUplinkStatus,omitempty"`

	// Value specs are extra parameters to include in the API request with OpenStack.
	// This is an extension point for the API, so what they do and if they are supported,
	// depends on the specific OpenStack implementation.
	// +optional
	// +listType=map
	// +listMapKey=name
	ValueSpecs []ValueSpec `json:"valueSpecs,omitempty"`
}

// ResolvedPortSpec is a PortOpts with all contained references fully resolved.
type ResolvedPortSpec struct {
	// Name is the name of the port.
	Name string `json:"name"`

	// Description is a human-readable description for the port.
	Description string `json:"description"`

	// NetworkID is the ID of the network the port will be created in.
	NetworkID string `json:"networkID"`

	// Tags applied to the port (and corresponding trunk, if a trunk is configured.)
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`

	// Trunk specifies whether trunking is enabled at the port level.
	// +optional
	Trunk optional.Bool `json:"trunk,omitempty"`

	// FixedIPs is a list of pairs of subnet and/or IP address to assign to the port. If specified, these must be subnets of the port's network.
	// +optional
	// +listType=atomic
	FixedIPs []ResolvedFixedIP `json:"fixedIPs,omitempty"`

	// SecurityGroups is a list of security group IDs to assign to the port.
	// +optional
	// +listType=atomic
	SecurityGroups []string `json:"securityGroups,omitempty"`

	ResolvedPortSpecFields `json:",inline"`
}

type PortStatus struct {
	// ID is the unique identifier of the port.
	// +required
	ID string `json:"id"`
}

type BindingProfile struct {
	// OVSHWOffload enables or disables the OVS hardware offload feature.
	// This flag is not required on OpenStack clouds since Yoga as Nova will set it automatically when the port is attached.
	// See: https://bugs.launchpad.net/nova/+bug/2020813
	// +optional
	OVSHWOffload *bool `json:"ovsHWOffload,omitempty"`

	// TrustedVF enables or disables the “trusted mode” for the VF.
	// +optional
	TrustedVF *bool `json:"trustedVF,omitempty"`
}

type FixedIP struct {
	// Subnet is an openstack subnet query that will return the id of a subnet to create
	// the fixed IP of a port in. This query must not return more than one subnet.
	// +optional
	Subnet *SubnetParam `json:"subnet,omitempty"`

	// IPAddress is a specific IP address to assign to the port. If Subnet
	// is also specified, IPAddress must be a valid IP address in the
	// subnet. If Subnet is not specified, IPAddress must be a valid IP
	// address in any subnet of the port's network.
	// +optional
	IPAddress optional.String `json:"ipAddress,omitempty"`
}

// ResolvedFixedIP is a FixedIP with the Subnet resolved to an ID.
type ResolvedFixedIP struct {
	// SubnetID is the id of a subnet to create the fixed IP of a port in.
	// +optional
	SubnetID optional.String `json:"subnet,omitempty"`

	// IPAddress is a specific IP address to assign to the port. If SubnetID
	// is also specified, IPAddress must be a valid IP address in the
	// subnet. If Subnet is not specified, IPAddress must be a valid IP
	// address in any subnet of the port's network.
	// +optional
	IPAddress optional.String `json:"ipAddress,omitempty"`
}

type AddressPair struct {
	// IPAddress is the IP address of the allowed address pair. Depending on
	// the configuration of Neutron, it may be supported to specify a CIDR
	// instead of a specific IP address.
	// +kubebuilder:validation:Required
	IPAddress string `json:"ipAddress"`

	// MACAddress is the MAC address of the allowed address pair. If not
	// specified, the MAC address will be the MAC address of the port.
	// +optional
	MACAddress optional.String `json:"macAddress,omitempty"`
}

type BastionStatus struct {
	ID         string        `json:"id,omitempty"`
	Name       string        `json:"name,omitempty"`
	SSHKeyName string        `json:"sshKeyName,omitempty"`
	State      InstanceState `json:"state,omitempty"`
	IP         string        `json:"ip,omitempty"`
	FloatingIP string        `json:"floatingIP,omitempty"`

	// Resolved contains parts of the bastion's machine spec with all
	// external references fully resolved.
	// +optional
	Resolved *ResolvedMachineSpec `json:"resolved,omitempty"`

	// Resources contains references to OpenStack resources created for the bastion.
	// +optional
	Resources *MachineResources `json:"resources,omitempty"`
}

type RootVolume struct {
	// SizeGiB is the size of the block device in gibibytes (GiB).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum:=1
	SizeGiB int `json:"sizeGiB"`

	BlockDeviceVolume `json:",inline"`
}

// BlockDeviceStorage is the storage type of a block device to create and
// contains additional storage options.
// +union
//
//nolint:godot
type BlockDeviceStorage struct {
	// Type is the type of block device to create.
	// This can be either "Volume" or "Local".
	// +unionDiscriminator
	Type BlockDeviceType `json:"type"`

	// Volume contains additional storage options for a volume block device.
	// +optional
	// +unionMember,optional
	Volume *BlockDeviceVolume `json:"volume,omitempty"`
}

// BlockDeviceVolume contains additional storage options for a volume block device.
type BlockDeviceVolume struct {
	// Type is the Cinder volume type of the volume.
	// If omitted, the default Cinder volume type that is configured in the OpenStack cloud
	// will be used.
	// +optional
	Type string `json:"type,omitempty"`

	// AvailabilityZone is the volume availability zone to create the volume
	// in. If not specified, the volume will be created without an explicit
	// availability zone.
	// +optional
	AvailabilityZone *VolumeAvailabilityZone `json:"availabilityZone,omitempty"`
}

// VolumeAZSource specifies where to obtain the availability zone for a volume.
// +kubebuilder:validation:Enum=Name;Machine
type VolumeAZSource string

const (
	VolumeAZFromName    VolumeAZSource = "Name"
	VolumeAZFromMachine VolumeAZSource = "Machine"
)

// VolumeAZName is the name of a volume availability zone. It may not contain spaces.
// +kubebuilder:validation:Pattern:="^[^ ]+$"
// +kubebuilder:validation:MinLength:=1
type VolumeAZName string

// VolumeAvailabilityZone specifies the availability zone for a volume.
// +kubebuilder:validation:XValidation:rule="!has(self.from) || self.from == 'Name' ? has(self.name) : !has(self.name)",message="name is required when from is 'Name' or default"
type VolumeAvailabilityZone struct {
	// From specifies where we will obtain the availability zone for the
	// volume. The options are "Name" and "Machine". If "Name" is specified
	// then the Name field must also be specified. If "Machine" is specified
	// the volume will use the value of FailureDomain, if any, from the
	// associated Machine.
	// +kubebuilder:default:=Name
	// +optional
	From VolumeAZSource `json:"from,omitempty"`

	// Name is the name of a volume availability zone to use. It is required
	// if From is "Name". The volume availability zone name may not contain
	// spaces.
	// +optional
	Name *VolumeAZName `json:"name,omitempty"`
}

// AdditionalBlockDevice is a block device to attach to the server.
type AdditionalBlockDevice struct {
	// Name of the block device in the context of a machine.
	// If the block device is a volume, the Cinder volume will be named
	// as a combination of the machine name and this name.
	// Also, this name will be used for tagging the block device.
	// Information about the block device tag can be obtained from the OpenStack
	// metadata API or the config drive.
	// Name cannot be 'root', which is reserved for the root volume.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// SizeGiB is the size of the block device in gibibytes (GiB).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum:=1
	SizeGiB int `json:"sizeGiB"`

	// Storage specifies the storage type of the block device and
	// additional storage options.
	Storage BlockDeviceStorage `json:"storage"`
}

// ServerGroupParam specifies an OpenStack server group. It may be specified by ID or filter, but not both.
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type ServerGroupParam struct {
	// ID is the ID of the server group to use.
	// +kubebuilder:validation:Format:=uuid
	ID optional.String `json:"id,omitempty"`

	// Filter specifies a query to select an OpenStack server group. If provided, it cannot be empty.
	Filter *ServerGroupFilter `json:"filter,omitempty"`
}

// ServerGroupFilter specifies a query to select an OpenStack server group. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type ServerGroupFilter struct {
	// Name is the name of a server group to look for.
	Name optional.String `json:"name,omitempty"`
}

func (f *ServerGroupFilter) IsZero() bool {
	if f == nil {
		return true
	}
	return f.Name == nil
}

// BlockDeviceType defines the type of block device to create.
type BlockDeviceType string

const (
	// LocalBlockDevice is an ephemeral block device attached to the server.
	LocalBlockDevice BlockDeviceType = "Local"

	// VolumeBlockDevice is a volume block device attached to the server.
	VolumeBlockDevice BlockDeviceType = "Volume"
)

// NetworkStatus contains basic information about an existing neutron network.
type NetworkStatus struct {
	Name string `json:"name"`
	ID   string `json:"id"`

	//+optional
	Tags []string `json:"tags,omitempty"`
}

// NetworkStatusWithSubnets represents basic information about an existing neutron network and an associated set of subnets.
type NetworkStatusWithSubnets struct {
	NetworkStatus `json:",inline"`

	// Subnets is a list of subnets associated with the default cluster network. Machines which use the default cluster network will get an address from all of these subnets.
	Subnets []Subnet `json:"subnets,omitempty"`
}

// Subnet represents basic information about the associated OpenStack Neutron Subnet.
type Subnet struct {
	Name string `json:"name"`
	ID   string `json:"id"`

	CIDR string `json:"cidr"`

	//+optional
	Tags []string `json:"tags,omitempty"`
}

// Router represents basic information about the associated OpenStack Neutron Router.
type Router struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	//+optional
	Tags []string `json:"tags,omitempty"`
	//+optional
	IPs []string `json:"ips,omitempty"`
}

// LoadBalancer represents basic information about the associated OpenStack LoadBalancer.
type LoadBalancer struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	IP         string `json:"ip"`
	InternalIP string `json:"internalIP"`
	//+optional
	AllowedCIDRs []string `json:"allowedCIDRs,omitempty"`
	//+optional
	Tags []string `json:"tags,omitempty"`
	// LoadBalancerNetwork contains information about network and/or subnets which the
	// loadbalancer is allocated on.
	// If subnets are specified within the LoadBalancerNetwork currently only the first
	// subnet in the list is taken into account.
	// +optional
	LoadBalancerNetwork *NetworkStatusWithSubnets `json:"loadBalancerNetwork,omitempty"`
}

// SecurityGroupStatus represents the basic information of the associated
// OpenStack Neutron Security Group.
type SecurityGroupStatus struct {
	// name of the security group
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// id of the security group
	// +kubebuilder:validation:Required
	ID string `json:"id"`
}

// SecurityGroupRuleSpec represent the basic information of the associated OpenStack
// Security Group Role.
// For now this is only used for the allNodesSecurityGroupRules but when we add
// other security groups, we'll need to add a validation because
// Remote* fields are mutually exclusive.
type SecurityGroupRuleSpec struct {
	// name of the security group rule.
	// It's used to identify the rule so it can be patched and will not be sent to the OpenStack API.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// description of the security group rule.
	// +optional
	Description *string `json:"description,omitempty"`

	// direction in which the security group rule is applied. The only values
	// allowed are "ingress" or "egress". For a compute instance, an ingress
	// security group rule is applied to incoming (ingress) traffic for that
	// instance. An egress rule is applied to traffic leaving the instance.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:enum=ingress;egress
	Direction string `json:"direction"`

	// etherType must be IPv4 or IPv6, and addresses represented in CIDR must match the
	// ingress or egress rules.
	// +kubebuilder:validation:enum=IPv4;IPv6
	// +optional
	EtherType *string `json:"etherType,omitempty"`

	// portRangeMin is a number in the range that is matched by the security group
	// rule. If the protocol is TCP or UDP, this value must be less than or equal
	// to the value of the portRangeMax attribute.
	// +optional
	PortRangeMin *int `json:"portRangeMin,omitempty"`

	// portRangeMax is a number in the range that is matched by the security group
	// rule. The portRangeMin attribute constrains the portRangeMax attribute.
	// +optional
	PortRangeMax *int `json:"portRangeMax,omitempty"`

	// protocol is the protocol that is matched by the security group rule.
	// +optional
	Protocol *string `json:"protocol,omitempty"`

	// remoteGroupID is the remote group ID to be associated with this security group rule.
	// You can specify either remoteGroupID or remoteIPPrefix or remoteManagedGroups.
	// +optional
	RemoteGroupID *string `json:"remoteGroupID,omitempty"`

	// remoteIPPrefix is the remote IP prefix to be associated with this security group rule.
	// You can specify either remoteGroupID or remoteIPPrefix or remoteManagedGroups.
	// +optional
	RemoteIPPrefix *string `json:"remoteIPPrefix,omitempty"`

	// remoteManagedGroups is the remote managed groups to be associated with this security group rule.
	// You can specify either remoteGroupID or remoteIPPrefix or remoteManagedGroups.
	// +optional
	RemoteManagedGroups []ManagedSecurityGroupName `json:"remoteManagedGroups,omitempty"`
}

// +kubebuilder:validation:Enum=bastion;controlplane;worker
type ManagedSecurityGroupName string

func (m ManagedSecurityGroupName) String() string {
	return string(m)
}

// InstanceState describes the state of an OpenStack instance.
type InstanceState string

var (
	// InstanceStateBuild is the string representing an instance in a build state.
	InstanceStateBuild = InstanceState("BUILD")

	// InstanceStateActive is the string representing an instance in an active state.
	InstanceStateActive = InstanceState("ACTIVE")

	// InstanceStateError is the string representing an instance in an error state.
	InstanceStateError = InstanceState("ERROR")

	// InstanceStateStopped is the string representing an instance in a stopped state.
	InstanceStateStopped = InstanceState("STOPPED")

	// InstanceStateShutoff is the string representing an instance in a shutoff state.
	InstanceStateShutoff = InstanceState("SHUTOFF")

	// InstanceStateDeleted is the string representing an instance in a deleted state.
	InstanceStateDeleted = InstanceState("DELETED")

	// InstanceStateUndefined is the string representing an undefined instance state.
	InstanceStateUndefined = InstanceState("")
)

// Bastion represents basic information about the bastion node. If you enable bastion, the spec has to be specified.
// +kubebuilder:validation:XValidation:rule="!self.enabled || has(self.spec)",message="spec is required if bastion is enabled"
type Bastion struct {
	// Enabled means that bastion is enabled. The bastion is enabled by
	// default if this field is not specified. Set this field to false to disable the
	// bastion.
	//
	// It is not currently possible to remove the bastion from the cluster
	// spec without first disabling it by setting this field to false and
	// waiting until the bastion has been deleted.
	// +kubebuilder:default:=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Spec for the bastion itself
	Spec *OpenStackMachineSpec `json:"spec,omitempty"`

	// AvailabilityZone is the failure domain that will be used to create the Bastion Spec.
	//+optional
	AvailabilityZone optional.String `json:"availabilityZone,omitempty"`

	// FloatingIP which will be associated to the bastion machine. It's the IP address, not UUID.
	// The floating IP should already exist and should not be associated with a port. If FIP of this address does not
	// exist, CAPO will try to create it, but by default only OpenStack administrators have privileges to do so.
	//+optional
	//+kubebuilder:validation:Format:=ipv4
	FloatingIP optional.String `json:"floatingIP,omitempty"`
}

func (b *Bastion) IsEnabled() bool {
	if b == nil {
		return false
	}
	return b.Enabled == nil || *b.Enabled
}

type APIServerLoadBalancer struct {
	// Enabled defines whether a load balancer should be created. This value
	// defaults to true if an APIServerLoadBalancer is given.
	//
	// There is no reason to set this to false. To disable creation of the
	// API server loadbalancer, omit the APIServerLoadBalancer field in the
	// cluster spec instead.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:default:=true
	Enabled *bool `json:"enabled"`

	// AdditionalPorts adds additional tcp ports to the load balancer.
	// +optional
	// +listType=set
	AdditionalPorts []int `json:"additionalPorts,omitempty"`

	// AllowedCIDRs restrict access to all API-Server listeners to the given address CIDRs.
	// +optional
	// +listType=set
	AllowedCIDRs []string `json:"allowedCIDRs,omitempty"`

	// Provider specifies name of a specific Octavia provider to use for the
	// API load balancer. The Octavia default will be used if it is not
	// specified.
	// +optional
	Provider optional.String `json:"provider,omitempty"`

	// Network defines which network should the load balancer be allocated on.
	//+optional
	Network *NetworkParam `json:"network,omitempty"`

	// Subnets define which subnets should the load balancer be allocated on.
	// It is expected that subnets are located on the network specified in this resource.
	// Only the first element is taken into account.
	// +optional
	// +listType=atomic
	// kubebuilder:validation:MaxLength:=2
	Subnets []SubnetParam `json:"subnets,omitempty"`

	// AvailabilityZone is the failure domain that will be used to create the APIServerLoadBalancer Spec.
	//+optional
	AvailabilityZone optional.String `json:"availabilityZone,omitempty"`

	// Flavor is the flavor name that will be used to create the APIServerLoadBalancer Spec.
	//+optional
	Flavor optional.String `json:"flavor,omitempty"`
}

func (s *APIServerLoadBalancer) IsZero() bool {
	return s == nil || ((s.Enabled == nil || !*s.Enabled) && len(s.AdditionalPorts) == 0 && len(s.AllowedCIDRs) == 0 && ptr.Deref(s.Provider, "") == "")
}

func (s *APIServerLoadBalancer) IsEnabled() bool {
	// The CRD default value for Enabled is true, so if the field is nil, it should be considered as true.
	return s != nil && (s.Enabled == nil || *s.Enabled)
}

// ResolvedMachineSpec contains resolved references to resources required by the machine.
type ResolvedMachineSpec struct {
	// ServerGroupID is the ID of the server group the machine should be added to and is calculated based on ServerGroupFilter.
	// +optional
	ServerGroupID string `json:"serverGroupID,omitempty"`

	// ImageID is the ID of the image to use for the machine and is calculated based on ImageFilter.
	// +optional
	ImageID string `json:"imageID,omitempty"`

	// FlavorID is the ID of the flavor to use.
	// +optional
	FlavorID string `json:"flavorID,omitempty"`

	// Ports is the fully resolved list of ports to create for the machine.
	// +optional
	Ports []ResolvedPortSpec `json:"ports,omitempty"`
}

type MachineResources struct {
	// Ports is the status of the ports created for the machine.
	// +optional
	Ports []PortStatus `json:"ports,omitempty"`
}

// ValueSpec represents a single value_spec key-value pair.
type ValueSpec struct {
	// Name is the name of the key-value pair.
	// This is just for identifying the pair and will not be sent to the OpenStack API.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Key is the key in the key-value pair.
	// +kubebuilder:validation:Required
	Key string `json:"key"`
	// Value is the value in the key-value pair.
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}
