/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta2

import (
	"strings"

	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

// OpenStackMachineTemplateResource describes the data needed to create a OpenStackMachine from a template.
type OpenStackMachineTemplateResource struct {
	// spec is the specification of the desired behavior of the machine.
	// +required
	Spec OpenStackMachineSpec `json:"spec,omitzero"`
}

type ResourceReference struct {
	// name is the name of the referenced resource
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
}

// ImageParam describes a glance image. It can be specified by ID, filter, or a
// reference to an ORC Image.
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type ImageParam struct {
	// id is the uuid of the image. ID will not be validated before use.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// filter describes a query for an image. If specified, the combination
	// of name and tags must return a single matching image or an error will
	// be raised.
	// +optional
	Filter *ImageFilter `json:"filter,omitempty"`

	// imageRef is a reference to an ORC Image in the same namespace as the
	// referring object.
	// +optional
	ImageRef *ResourceReference `json:"imageRef,omitempty"`
}

// ImageFilter describes a query for an image.
// +kubebuilder:validation:MinProperties:=1
type ImageFilter struct {
	// name is the name of the desired image. If specified, the combination of name and tags must return a single matching image or an error will be raised.
	// +optional
	Name optional.String `json:"name,omitempty"`

	// tags are the tags associated with the desired image. If specified, the combination of name and tags must return a single matching image or an error will be raised.
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

// FlavorParam describes a nova flavor. It can be specified by ID or filter
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type FlavorParam struct {
	// id is the uuid of the flavor. ID will not be validated before use.
	// +optional
	ID optional.String `json:"id,omitempty"`

	// filter describes a query for a flavor.
	// +optional
	Filter *FlavorFilter `json:"filter,omitempty"`
}

// FlavorFilter describes a query for a flavor. If defined,
// the combination of attributes should return exactly one
// flavor, if not an error will be raised.
// +kubebuilder:validation:MinProperties:=1
type FlavorFilter struct {
	// name is the name of the desired flavor.
	// +optional
	Name optional.String `json:"name,omitempty"`
}

func (f *FlavorFilter) IsZero() bool {
	if f == nil {
		return true
	}

	return f.Name == nil
}

type ExternalRouterIPParam struct {
	// fixedIP is the FixedIP in the corresponding subnet.
	// +optional
	FixedIP string `json:"fixedIP,omitempty"`
	// subnet is the subnet in which the FixedIP is used for the Gateway of this router.
	// +required
	Subnet SubnetParam `json:"subnet,omitzero"`
}

// NeutronTag represents a tag on a Neutron resource.
// It may not be empty and may not contain commas.
// +kubebuilder:validation:Pattern:="^[^,]+$"
// +kubebuilder:validation:MinLength:=1
type NeutronTag string

type FilterByNeutronTags struct {
	// tags is a list of tags to filter by. If specified, the resource must
	// have all of the tags specified to be included in the result.
	// +listType=set
	// +optional
	Tags []NeutronTag `json:"tags,omitempty"`

	// tagsAny is a list of tags to filter by. If specified, the resource
	// must have at least one of the tags specified to be included in the
	// result.
	// +listType=set
	// +optional
	TagsAny []NeutronTag `json:"tagsAny,omitempty"`

	// notTags is a list of tags to filter by. If specified, resources which
	// contain all of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	NotTags []NeutronTag `json:"notTags,omitempty"`

	// notTagsAny is a list of tags to filter by. If specified, resources
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
	// id is the ID of the security group to use. If ID is provided, the other filters cannot be provided. Must be in UUID format.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// filter specifies a query to select an OpenStack security group. If provided, cannot be empty.
	// +optional
	Filter *SecurityGroupFilter `json:"filter,omitempty"`
}

// SecurityGroupFilter specifies a query to select an OpenStack security group. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type SecurityGroupFilter struct {
	// name filters security groups by name.
	// +optional
	Name string `json:"name,omitempty"`
	// description filters security groups by description.
	// +optional
	Description string `json:"description,omitempty"`
	// projectID filters security groups by project ID.
	// +optional
	ProjectID string `json:"projectID,omitempty"`

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
	// id is the ID of the network to use. If ID is provided, the other filters cannot be provided. Must be in UUID format.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// filter specifies a filter to select an OpenStack network. If provided, cannot be empty.
	// +optional
	Filter *NetworkFilter `json:"filter,omitempty"`
}

// NetworkFilter specifies a query to select an OpenStack network. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type NetworkFilter struct {
	// name filters networks by name.
	// +optional
	Name string `json:"name,omitempty"`
	// description filters networks by description.
	// +optional
	Description string `json:"description,omitempty"`
	// projectID filters networks by project ID.
	// +optional
	ProjectID string `json:"projectID,omitempty"`

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
	// id is the uuid of the subnet. It will not be validated.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// filter specifies a filter to select the subnet. It must match exactly one subnet.
	// +optional
	Filter *SubnetFilter `json:"filter,omitempty"`
}

// SubnetFilter specifies a filter to select a subnet. At least one parameter must be specified.
// +kubebuilder:validation:MinProperties:=1
type SubnetFilter struct {
	// name filters subnets by name.
	// +optional
	Name string `json:"name,omitempty"`
	// description filters subnets by description.
	// +optional
	Description string `json:"description,omitempty"`
	// projectID filters subnets by project ID.
	// +optional
	ProjectID string `json:"projectID,omitempty"`
	// ipVersion filters subnets by IP version.
	// +optional
	IPVersion int32 `json:"ipVersion,omitempty"`
	// gatewayIP filters subnets by gateway IP.
	// +optional
	GatewayIP string `json:"gatewayIP,omitempty"`
	// cidr filters subnets by CIDR.
	// +optional
	CIDR string `json:"cidr,omitempty"`
	// ipv6AddressMode filters subnets by IPv6 address mode.
	// +optional
	IPv6AddressMode string `json:"ipv6AddressMode,omitempty"`
	// ipv6RAMode filters subnets by IPv6 Router Advertisement mode.
	// +optional
	IPv6RAMode string `json:"ipv6RAMode,omitempty"`

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
	// id is the ID of the router to use. If ID is provided, the other filters cannot be provided. Must be in UUID format.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// filter specifies a filter to select an OpenStack router. If provided, cannot be empty.
	// +optional
	Filter *RouterFilter `json:"filter,omitempty"`
}

// RouterFilter specifies a query to select an OpenStack router. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type RouterFilter struct {
	// name filters routers by name.
	// +optional
	Name string `json:"name,omitempty"`
	// description filters routers by description.
	// +optional
	Description string `json:"description,omitempty"`
	// projectID filters routers by project ID.
	// +optional
	ProjectID string `json:"projectID,omitempty"`

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
	// cidr is representing the IP address range used to create the subnet, e.g. 10.0.0.0/24.
	// This field is required when defining a subnet.
	// +required
	// +kubebuilder:validation:MinLength=1
	CIDR string `json:"cidr,omitempty"`

	// dnsNameservers holds a list of DNS server addresses that will be provided when creating
	// the subnet. These addresses need to have the same IP version as CIDR.
	// +listType=atomic
	// +optional
	DNSNameservers []string `json:"dnsNameservers,omitempty"`

	// allocationPools is an array of AllocationPool objects that will be applied to OpenStack Subnet being created.
	// If set, OpenStack will only allocate these IPs for Machines. It will still be possible to create ports from
	// outside of these ranges manually.
	// +listType=atomic
	// +optional
	AllocationPools []AllocationPool `json:"allocationPools,omitempty"`
}

type AllocationPool struct {
	// start represents the start of the AllocationPool, that is the lowest IP of the pool.
	// +required
	// +kubebuilder:validation:MinLength=1
	Start string `json:"start,omitempty"`

	// end represents the end of the AlloctionPool, that is the highest IP of the pool.
	// +required
	// +kubebuilder:validation:MinLength=1
	End string `json:"end,omitempty"`
}

type PortOpts struct {
	// network is a query for an openstack network that the port will be created or discovered on.
	// This will fail if the query returns more than one network.
	// +optional
	Network *NetworkParam `json:"network,omitempty"`

	// description is a human-readable description for the port.
	// +optional
	Description optional.String `json:"description,omitempty"`

	// nameSuffix will be appended to the name of the port if specified. If unspecified, instead the 0-based index of the port in the list is used.
	// +optional
	NameSuffix optional.String `json:"nameSuffix,omitempty"`

	// fixedIPs is a list of pairs of subnet and/or IP address to assign to the port. If specified, these must be subnets of the port's network.
	// +optional
	// +listType=atomic
	FixedIPs []FixedIP `json:"fixedIPs,omitempty"`

	// securityGroups is a list of the names, uuids, filters or any combination these of the security groups to assign to the instance.
	// +optional
	// +listType=atomic
	SecurityGroups []SecurityGroupParam `json:"securityGroups,omitempty"`

	// tags applied to the port (and corresponding trunk, if a trunk is configured.)
	// These tags are applied in addition to the instance's tags, which will also be applied to the port.
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`

	// trunk specifies whether trunking is enabled at the port level. If not
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
	// adminStateUp specifies whether the port should be created in the up (true) or down (false) state. The default is up.
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// macAddress specifies the MAC address of the port. If not specified, the MAC address will be generated.
	// +optional
	MACAddress optional.String `json:"macAddress,omitempty"`

	// allowedAddressPairs is a list of address pairs which Neutron will
	// allow the port to send traffic from in addition to the port's
	// addresses. If not specified, the MAC Address will be the MAC Address
	// of the port. Depending on the configuration of Neutron, it may be
	// supported to specify a CIDR instead of a specific IP address.
	// +listType=atomic
	// +optional
	AllowedAddressPairs []AddressPair `json:"allowedAddressPairs,omitempty"`

	// hostID specifies the ID of the host where the port resides.
	// +optional
	HostID optional.String `json:"hostID,omitempty"`

	// vnicType specifies the type of vNIC which this port should be
	// attached to. This is used to determine which mechanism driver(s) to
	// be used to bind the port. The valid values are normal, macvtap,
	// direct, baremetal, direct-physical, virtio-forwarder, smart-nic and
	// remote-managed, although these values will not be validated in this
	// API to ensure compatibility with future neutron changes or custom
	// implementations. What type of vNIC is actually available depends on
	// deployments. If not specified, the Neutron default value is used.
	// +optional
	VNICType optional.String `json:"vnicType,omitempty"`

	// profile is a set of key-value pairs that are used for binding
	// details. We intentionally don't expose this as a map[string]string
	// because we only want to enable the users to set the values of the
	// keys that are known to work in OpenStack Networking API.  See
	// https://docs.openstack.org/api-ref/network/v2/index.html?expanded=create-port-detail#create-port
	// To set profiles, your tenant needs permissions rule:create_port, and
	// rule:create_port:binding:profile
	// +optional
	Profile *BindingProfile `json:"profile,omitempty"`

	// enablePortSecurity enables or disables the port security when set.
	// When not set, it takes the value of the corresponding field at the network level.
	// +optional
	EnablePortSecurity *bool `json:"enablePortSecurity,omitempty"`

	// propagateUplinkStatus enables or disables the propagate uplink status on the port.
	// +optional
	PropagateUplinkStatus *bool `json:"propagateUplinkStatus,omitempty"`

	// valueSpecs are extra parameters to include in the API request with OpenStack.
	// This is an extension point for the API, so what they do and if they are supported,
	// depends on the specific OpenStack implementation.
	// +optional
	// +listType=map
	// +listMapKey=name
	ValueSpecs []ValueSpec `json:"valueSpecs,omitempty"`
}

// ResolvedPortSpec is a PortOpts with all contained references fully resolved.
type ResolvedPortSpec struct {
	// name is the name of the port.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the port.
	// +optional
	Description string `json:"description"`

	// networkID is the ID of the network the port will be created in.
	// +required
	// +kubebuilder:validation:MinLength=1
	NetworkID string `json:"networkID,omitempty"`

	// tags applied to the port (and corresponding trunk, if a trunk is configured.)
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`

	// trunk specifies whether trunking is enabled at the port level.
	// +optional
	Trunk optional.Bool `json:"trunk,omitempty"`

	// fixedIPs is a list of pairs of subnet and/or IP address to assign to the port. If specified, these must be subnets of the port's network.
	// +optional
	// +listType=atomic
	FixedIPs []ResolvedFixedIP `json:"fixedIPs,omitempty"`

	// securityGroups is a list of security group IDs to assign to the port.
	// +optional
	// +listType=atomic
	SecurityGroups []string `json:"securityGroups,omitempty"`

	ResolvedPortSpecFields `json:",inline"`
}

type PortStatus struct {
	// id is the unique identifier of the port.
	// +required
	// +kubebuilder:validation:MinLength=1
	ID string `json:"id,omitempty"`
}

type BindingProfile struct {
	// ovsHWOffload enables or disables the OVS hardware offload feature.
	// This flag is not required on OpenStack clouds since Yoga as Nova will set it automatically when the port is attached.
	// See: https://bugs.launchpad.net/nova/+bug/2020813
	// +optional
	OVSHWOffload *bool `json:"ovsHWOffload,omitempty"`

	// trustedVF enables or disables the “trusted mode” for the VF.
	// +optional
	TrustedVF *bool `json:"trustedVF,omitempty"`
}

type FixedIP struct {
	// subnet is an openstack subnet query that will return the id of a subnet to create
	// the fixed IP of a port in. This query must not return more than one subnet.
	// +optional
	Subnet *SubnetParam `json:"subnet,omitempty"`

	// ipAddress is a specific IP address to assign to the port. If Subnet
	// is also specified, IPAddress must be a valid IP address in the
	// subnet. If Subnet is not specified, IPAddress must be a valid IP
	// address in any subnet of the port's network.
	// +optional
	IPAddress optional.String `json:"ipAddress,omitempty"`
}

// ResolvedFixedIP is a FixedIP with the Subnet resolved to an ID.
type ResolvedFixedIP struct {
	// subnet is the ID of a subnet to create the fixed IP of a port in.
	// +optional
	SubnetID optional.String `json:"subnet,omitempty"`

	// ipAddress is a specific IP address to assign to the port. If SubnetID
	// is also specified, IPAddress must be a valid IP address in the
	// subnet. If Subnet is not specified, IPAddress must be a valid IP
	// address in any subnet of the port's network.
	// +optional
	IPAddress optional.String `json:"ipAddress,omitempty"`
}

type AddressPair struct {
	// ipAddress is the IP address of the allowed address pair. Depending on
	// the configuration of Neutron, it may be supported to specify a CIDR
	// instead of a specific IP address.
	// +required
	// +kubebuilder:validation:MinLength=1
	IPAddress string `json:"ipAddress,omitempty"`

	// macAddress is the MAC address of the allowed address pair. If not
	// specified, the MAC address will be the MAC address of the port.
	// +optional
	MACAddress optional.String `json:"macAddress,omitempty"`
}

type BastionStatus struct {
	// id is the unique identifier of the bastion.
	// +optional
	ID string `json:"id,omitempty"`
	// name is the name of the bastion.
	// +optional
	Name string `json:"name,omitempty"`
	// sshKeyName is the name of the SSH key used for the bastion.
	// +optional
	SSHKeyName string `json:"sshKeyName,omitempty"`
	// state is the current state of the bastion.
	// +optional
	State InstanceState `json:"state,omitempty"`
	// ip is the IP address of the bastion.
	// +optional
	IP string `json:"ip,omitempty"`
	// floatingIP is the floating IP address of the bastion.
	// +optional
	FloatingIP string `json:"floatingIP,omitempty"`

	// resolved contains parts of the bastion's machine spec with all
	// external references fully resolved.
	// +optional
	Resolved *ResolvedMachineSpec `json:"resolved,omitempty"`

	// resources contains references to OpenStack resources created for the bastion.
	// +optional
	Resources *MachineResources `json:"resources,omitempty"`
}

type RootVolume struct {
	// sizeGiB is the size of the block device in gibibytes (GiB).
	// +required
	// +kubebuilder:validation:Minimum:=1
	SizeGiB int32 `json:"sizeGiB,omitempty"`

	BlockDeviceVolume `json:",inline"`
}

// BlockDeviceStorage is the storage type of a block device to create and
// contains additional storage options.
// +union
//
//nolint:godot
type BlockDeviceStorage struct {
	// type is the type of block device to create.
	// This can be either "Volume" or "Local".
	// +unionDiscriminator
	// +required
	// +kubebuilder:validation:Enum=Local;Volume
	Type BlockDeviceType `json:"type,omitempty"`

	// volume contains additional storage options for a volume block device.
	// +optional
	// +unionMember,optional
	Volume *BlockDeviceVolume `json:"volume,omitempty"`
}

// BlockDeviceVolume contains additional storage options for a volume block device.
type BlockDeviceVolume struct {
	// type is the Cinder volume type of the volume.
	// If omitted, the default Cinder volume type that is configured in the OpenStack cloud
	// will be used.
	// +optional
	Type string `json:"type,omitempty"`

	// availabilityZone is the volume availability zone to create the volume
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
	// from specifies where we will obtain the availability zone for the
	// volume. The options are "Name" and "Machine". If "Name" is specified
	// then the Name field must also be specified. If "Machine" is specified
	// the volume will use the value of FailureDomain, if any, from the
	// associated Machine.
	// +kubebuilder:default:=Name
	// +optional
	From VolumeAZSource `json:"from,omitempty"`

	// name is the name of a volume availability zone to use. It is required
	// if From is "Name". The volume availability zone name may not contain
	// spaces.
	// +optional
	Name *VolumeAZName `json:"name,omitempty"`
}

// AdditionalBlockDevice is a block device to attach to the server.
type AdditionalBlockDevice struct {
	// name of the block device in the context of a machine.
	// If the block device is a volume, the Cinder volume will be named
	// as a combination of the machine name and this name.
	// Also, this name will be used for tagging the block device.
	// Information about the block device tag can be obtained from the OpenStack
	// metadata API or the config drive.
	// Name cannot be 'root', which is reserved for the root volume.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// sizeGiB is the size of the block device in gibibytes (GiB).
	// +required
	// +kubebuilder:validation:Minimum:=1
	SizeGiB int32 `json:"sizeGiB,omitempty"`

	// storage specifies the storage type of the block device and
	// additional storage options.
	// +required
	Storage BlockDeviceStorage `json:"storage,omitzero"`
}

// ServerGroupParam specifies an OpenStack server group. It may be specified by ID or filter, but not both.
// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type ServerGroupParam struct {
	// id is the ID of the server group to use.
	// +kubebuilder:validation:Format:=uuid
	// +optional
	ID optional.String `json:"id,omitempty"`

	// filter specifies a query to select an OpenStack server group. If provided, it cannot be empty.
	// +optional
	Filter *ServerGroupFilter `json:"filter,omitempty"`
}

// ServerGroupFilter specifies a query to select an OpenStack server group. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type ServerGroupFilter struct {
	// name is the name of a server group to look for.
	// +optional
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
	// name is the name of the network.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// id is the unique identifier of the network.
	// +required
	// +kubebuilder:validation:MinLength=1
	ID string `json:"id,omitempty"`

	// tags is a list of tags on the network.
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// NetworkStatusWithSubnets represents basic information about an existing neutron network and an associated set of subnets.
type NetworkStatusWithSubnets struct {
	NetworkStatus `json:",inline"`

	// subnets is a list of subnets associated with the default cluster network. Machines which use the default cluster network will get an address from all of these subnets.
	// +listType=atomic
	// +optional
	Subnets []Subnet `json:"subnets,omitempty"`
}

// Subnet represents basic information about the associated OpenStack Neutron Subnet.
type Subnet struct {
	// name is the name of the subnet.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// id is the unique identifier of the subnet.
	// +required
	// +kubebuilder:validation:MinLength=1
	ID string `json:"id,omitempty"`

	// cidr is the CIDR of the subnet.
	// +required
	// +kubebuilder:validation:MinLength=1
	CIDR string `json:"cidr,omitempty"`

	// tags is a list of tags on the subnet.
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`
}

// Router represents basic information about the associated OpenStack Neutron Router.
type Router struct {
	// name is the name of the router.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// id is the unique identifier of the router.
	// +required
	// +kubebuilder:validation:MinLength=1
	ID string `json:"id,omitempty"`
	// tags is a list of tags on the router.
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`
	// ips is a list of IP addresses assigned to the router.
	// +listType=set
	// +optional
	IPs []string `json:"ips,omitempty"`
}

// LoadBalancer represents basic information about the associated OpenStack LoadBalancer.
type LoadBalancer struct {
	// name is the name of the load balancer.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// id is the unique identifier of the load balancer.
	// +required
	// +kubebuilder:validation:MinLength=1
	ID string `json:"id,omitempty"`
	// ip is the IP address of the load balancer.
	// +required
	// +kubebuilder:validation:MinLength=1
	IP string `json:"ip,omitempty"`
	// internalIP is the internal IP address of the load balancer.
	// +required
	// +kubebuilder:validation:MinLength=1
	InternalIP string `json:"internalIP,omitempty"`
	// allowedCIDRs is a list of CIDRs that are allowed to access the load balancer.
	// +listType=set
	// +optional
	AllowedCIDRs []string `json:"allowedCIDRs,omitempty"`
	// tags is a list of tags on the load balancer.
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`
	// loadBalancerNetwork contains information about network and/or subnets which the
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
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// id of the security group
	// +required
	// +kubebuilder:validation:MinLength=1
	ID string `json:"id,omitempty"`
}

// SecurityGroupRuleSpec represent the basic information of the associated OpenStack
// Security Group Role.
// For now this is only used for the clusterNodesSecurityGroupRules but when we add
// other security groups, we'll need to add a validation because
// Remote* fields are mutually exclusive.
type SecurityGroupRuleSpec struct {
	// name of the security group rule.
	// It's used to identify the rule so it can be patched and will not be sent to the OpenStack API.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// description of the security group rule.
	// +optional
	Description *string `json:"description,omitempty"`

	// direction in which the security group rule is applied. The only values
	// allowed are "ingress" or "egress". For a compute instance, an ingress
	// security group rule is applied to incoming (ingress) traffic for that
	// instance. An egress rule is applied to traffic leaving the instance.
	// +required
	// +kubebuilder:validation:Enum=ingress;egress
	Direction string `json:"direction,omitempty"`

	// etherType must be IPv4 or IPv6, and addresses represented in CIDR must match the
	// ingress or egress rules.
	// +kubebuilder:validation:Enum=IPv4;IPv6
	// +optional
	EtherType *string `json:"etherType,omitempty"`

	// portRangeMin is a number in the range that is matched by the security group
	// rule. If the protocol is TCP or UDP, this value must be less than or equal
	// to the value of the portRangeMax attribute.
	// +optional
	PortRangeMin *int32 `json:"portRangeMin,omitempty"`

	// portRangeMax is a number in the range that is matched by the security group
	// rule. The portRangeMin attribute constrains the portRangeMax attribute.
	// +optional
	PortRangeMax *int32 `json:"portRangeMax,omitempty"`

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
	// +listType=set
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

	// InstanceStateSoftDeleted is the string representing an instance in a soft-deleted state.
	// This state occurs when OpenStack is configured with a reclaim_instance_interval > 0,
	// allowing recovery of deleted instances within the reclaim period.
	InstanceStateSoftDeleted = InstanceState("SOFT_DELETED")

	// InstanceStateUndefined is the string representing an undefined instance state.
	InstanceStateUndefined = InstanceState("")
)

// Bastion represents basic information about the bastion node. If you enable bastion, the spec has to be specified.
// +kubebuilder:validation:XValidation:rule="!self.enabled || has(self.spec)",message="spec is required if bastion is enabled"
type Bastion struct {
	// enabled means that bastion is enabled. The bastion is enabled by
	// default if this field is not specified. Set this field to false to disable the
	// bastion.
	//
	// It is not currently possible to remove the bastion from the cluster
	// spec without first disabling it by setting this field to false and
	// waiting until the bastion has been deleted.
	// +kubebuilder:default:=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// spec for the bastion itself
	// +optional
	Spec *OpenStackMachineSpec `json:"spec,omitempty"`

	// availabilityZone is the failure domain that will be used to create the Bastion Spec.
	// +optional
	AvailabilityZone optional.String `json:"availabilityZone,omitempty"`

	// floatingIP which will be associated to the bastion machine. It's the IP address, not UUID.
	// The floating IP should already exist and should not be associated with a port. If FIP of this address does not
	// exist, CAPO will try to create it, but by default only OpenStack administrators have privileges to do so.
	// +optional
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
	// enabled defines whether a load balancer should be created. This value
	// defaults to true if an APIServerLoadBalancer is given.
	//
	// There is no reason to set this to false. To disable creation of the
	// API server loadbalancer, omit the APIServerLoadBalancer field in the
	// cluster spec instead.
	//
	// +optional
	// +kubebuilder:default:=true
	Enabled *bool `json:"enabled,omitempty"`

	// additionalPorts adds additional tcp ports to the load balancer.
	// +optional
	// +listType=set
	AdditionalPorts []int32 `json:"additionalPorts,omitempty"`

	// allowedCIDRs restrict access to all API-Server listeners to the given address CIDRs.
	// +optional
	// +listType=set
	AllowedCIDRs []string `json:"allowedCIDRs,omitempty"`

	// provider specifies name of a specific Octavia provider to use for the
	// API load balancer. The Octavia default will be used if it is not
	// specified.
	// +optional
	Provider optional.String `json:"provider,omitempty"`

	// network defines which network should the load balancer be allocated on.
	// +optional
	Network *NetworkParam `json:"network,omitempty"`

	// subnets define which subnets should the load balancer be allocated on.
	// It is expected that subnets are located on the network specified in this resource.
	// Only the first element is taken into account.
	// +optional
	// +listType=atomic
	// kubebuilder:validation:MaxLength:=2
	Subnets []SubnetParam `json:"subnets,omitempty"`

	// availabilityZone is the failure domain that will be used to create the APIServerLoadBalancer Spec.
	// +optional
	AvailabilityZone optional.String `json:"availabilityZone,omitempty"`

	// flavor is the flavor name that will be used to create the APIServerLoadBalancer Spec.
	// +optional
	Flavor optional.String `json:"flavor,omitempty"`

	// monitor contains configuration for the load balancer health monitor.
	// +optional
	Monitor *APIServerLoadBalancerMonitor `json:"monitor,omitempty"`
}

// APIServerLoadBalancerMonitor contains configuration for the load balancer health monitor.
type APIServerLoadBalancerMonitor struct {
	// delay is the time in seconds between sending probes to members.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default:10
	Delay int32 `json:"delay,omitempty"`

	// timeout is the maximum time in seconds for a monitor to wait for a connection to be established before it times out.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default:5
	Timeout int32 `json:"timeout,omitempty"`

	// maxRetries is the number of successful checks before changing the operating status of the member to ONLINE.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:default:5
	MaxRetries int32 `json:"maxRetries,omitempty"`

	// maxRetriesDown is the number of allowed check failures before changing the operating status of the member to ERROR.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:default:3
	MaxRetriesDown int32 `json:"maxRetriesDown,omitempty"`
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
	// serverGroupID is the ID of the server group the machine should be added to and is calculated based on ServerGroupFilter.
	// +optional
	ServerGroupID string `json:"serverGroupID,omitempty"`

	// imageID is the ID of the image to use for the machine and is calculated based on ImageFilter.
	// +optional
	ImageID string `json:"imageID,omitempty"`

	// flavorID is the ID of the flavor to use.
	// +optional
	FlavorID string `json:"flavorID,omitempty"`

	// ports is the fully resolved list of ports to create for the machine.
	// +listType=atomic
	// +optional
	Ports []ResolvedPortSpec `json:"ports,omitempty"`
}

type MachineResources struct {
	// ports is the status of the ports created for the machine.
	// +listType=atomic
	// +optional
	Ports []PortStatus `json:"ports,omitempty"`
}

// ValueSpec represents a single value_spec key-value pair.
type ValueSpec struct {
	// name is the name of the key-value pair.
	// This is just for identifying the pair and will not be sent to the OpenStack API.
	// +required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`
	// key is the key in the key-value pair.
	// +required
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key,omitempty"`
	// value is the value in the key-value pair.
	// +required
	// +kubebuilder:validation:MinLength=1
	Value string `json:"value,omitempty"`
}

// JoinTags joins a slice of tags into a comma separated list of tags.
func JoinTags(tags []NeutronTag) string {
	var b strings.Builder
	for i := range tags {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(string(tags[i]))
	}
	return b.String()
}
