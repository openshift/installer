/*
Copyright 2024 The ORC Authors.

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

package v1alpha1

// TODO validations:
//
// * IP addresses in CIDR, AllocationPools, Gateway, DNSNameserver(?), and
//   HostRoutes match the version in IPVersion (Spec and SubnetFilter)
// * IPv6 may only be set if IPVersion is 6 (Spec and SubnetFilter)
// * AllocationPools must be in CIDR

// SubnetFilter specifies a filter to select a subnet. At least one parameter must be specified.
// +kubebuilder:validation:MinProperties:=1
type SubnetFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description of the existing resource
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// ipVersion of the existing resource
	// +optional
	IPVersion *IPVersion `json:"ipVersion,omitempty"`

	// gatewayIP is the IP address of the gateway of the existing resource
	// +optional
	GatewayIP *IPvAny `json:"gatewayIP,omitempty"`

	// cidr of the existing resource
	// +optional
	CIDR *CIDR `json:"cidr,omitempty"`

	// ipv6 options of the existing resource
	// +optional
	IPv6 *IPv6Options `json:"ipv6,omitempty"`

	// networkRef is a reference to the ORC Network which this subnet is associated with.
	// +optional
	NetworkRef KubernetesNameRef `json:"networkRef"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

type SubnetResourceSpec struct {
	// name is a human-readable name of the subnet. If not set, the object's name will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// networkRef is a reference to the ORC Network which this subnet is associated with.
	// +required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="networkRef is immutable"
	NetworkRef KubernetesNameRef `json:"networkRef,omitempty"`

	// tags is a list of tags which will be applied to the subnet.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	Tags []NeutronTag `json:"tags,omitempty"`

	// ipVersion is the IP version for the subnet.
	// +required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="ipVersion is immutable"
	IPVersion IPVersion `json:"ipVersion"`

	// cidr is the address CIDR of the subnet. It must match the IP version specified in IPVersion.
	// +required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="cidr is immutable"
	CIDR CIDR `json:"cidr,omitempty"`

	// allocationPools are IP Address pools that will be available for DHCP. IP
	// addresses must be in CIDR.
	// +kubebuilder:validation:MaxItems:=32
	// +listType=atomic
	// +optional
	AllocationPools []AllocationPool `json:"allocationPools,omitempty"`

	// gateway specifies the default gateway of the subnet. If not specified,
	// neutron will add one automatically. To disable this behaviour, specify a
	// gateway with a type of None.
	// +optional
	Gateway *SubnetGateway `json:"gateway,omitempty"`

	// enableDHCP will either enable to disable the DHCP service.
	// +optional
	EnableDHCP *bool `json:"enableDHCP,omitempty"`

	// dnsNameservers are the nameservers to be set via DHCP.
	// +kubebuilder:validation:MaxItems:=16
	// +listType=set
	// +optional
	DNSNameservers []IPvAny `json:"dnsNameservers,omitempty"`

	// dnsPublishFixedIP will either enable or disable the publication of
	// fixed IPs to the DNS. Defaults to false.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="dnsPublishFixedIP is immutable"
	DNSPublishFixedIP *bool `json:"dnsPublishFixedIP,omitempty"`

	// hostRoutes are any static host routes to be set via DHCP.
	// +kubebuilder:validation:MaxItems:=256
	// +listType=atomic
	// +optional
	HostRoutes []HostRoute `json:"hostRoutes,omitempty"`

	// ipv6 contains IPv6-specific options. It may only be set if IPVersion is 6.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="ipv6 is immutable"
	IPv6 *IPv6Options `json:"ipv6,omitempty"`

	// routerRef specifies a router to attach the subnet to
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="routerRef is immutable"
	RouterRef *KubernetesNameRef `json:"routerRef,omitempty"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="projectRef is immutable"
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	// TODO: Support service types
	// TODO: Support subnet pools
}

type SubnetResourceStatus struct {
	// name is the human-readable name of the subnet. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// ipVersion specifies IP version, either `4' or `6'.
	// +optional
	IPVersion *int32 `json:"ipVersion,omitempty"`

	// cidr representing IP range for this subnet, based on IP version.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	CIDR string `json:"cidr,omitempty"`

	// gatewayIP is the default gateway used by devices in this subnet, if any.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	GatewayIP string `json:"gatewayIP,omitempty"`

	// dnsNameservers is a list of name servers used by hosts in this subnet.
	// +kubebuilder:validation:MaxItems:=16
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	DNSNameservers []string `json:"dnsNameservers,omitempty"`

	// dnsPublishFixedIP specifies whether the fixed IP addresses are published to the DNS.
	// +optional
	DNSPublishFixedIP *bool `json:"dnsPublishFixedIP,omitempty"`

	// allocationPools is a list of sub-ranges within CIDR available for dynamic
	// allocation to ports.
	// +kubebuilder:validation:MaxItems:=32
	// +listType=atomic
	// +optional
	AllocationPools []AllocationPoolStatus `json:"allocationPools,omitempty"`

	// hostRoutes is a list of routes that should be used by devices with IPs
	// from this subnet (not including local subnet route).
	// +kubebuilder:validation:MaxItems:=256
	// +listType=atomic
	// +optional
	HostRoutes []HostRouteStatus `json:"hostRoutes,omitempty"`

	// enableDHCP specifies whether DHCP is enabled for this subnet or not.
	// +optional
	EnableDHCP *bool `json:"enableDHCP,omitempty"`

	// networkID is the ID of the network to which the subnet belongs.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	NetworkID string `json:"networkID,omitempty"`

	// projectID is the project owner of the subnet.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// ipv6AddressMode specifies mechanisms for assigning IPv6 IP addresses.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	IPv6AddressMode string `json:"ipv6AddressMode,omitempty"`

	// ipv6RAMode is the IPv6 router advertisement mode. It specifies
	// whether the networking service should transmit ICMPv6 packets.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	IPv6RAMode string `json:"ipv6RAMode,omitempty"`

	// subnetPoolID is the id of the subnet pool associated with the subnet.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	SubnetPoolID string `json:"subnetPoolID,omitempty"`

	// tags optionally set via extensions/attributestags
	// +kubebuilder:validation:MaxItems:=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`

	NeutronStatusMetadata `json:",inline"`
}

// +kubebuilder:validation:Enum:=slaac;dhcpv6-stateful;dhcpv6-stateless
type IPv6AddressMode string

const (
	IPv6AddressModeSLAAC           = "slaac"
	IPv6AddressModeDHCPv6Stateful  = "dhcpv6-stateful"
	IPv6AddressModeDHCPv6Stateless = "dhcpv6-stateless"
)

// +kubebuilder:validation:Enum:=slaac;dhcpv6-stateful;dhcpv6-stateless
type IPv6RAMode string

const (
	IPv6RAModeSLAAC           = "slaac"
	IPv6RAModeDHCPv6Stateful  = "dhcpv6-stateful"
	IPv6RAModeDHCPv6Stateless = "dhcpv6-stateless"
)

// +kubebuilder:validation:MinProperties:=1
type IPv6Options struct {
	// addressMode specifies mechanisms for assigning IPv6 IP addresses.
	// +optional
	AddressMode *IPv6AddressMode `json:"addressMode,omitempty"`

	// raMode specifies the IPv6 router advertisement mode. It specifies whether
	// the networking service should transmit ICMPv6 packets.
	// +optional
	RAMode *IPv6RAMode `json:"raMode,omitempty"`
}

type SubnetGatewayType string

const (
	SubnetGatewayTypeNone      = "None"
	SubnetGatewayTypeAutomatic = "Automatic"
	SubnetGatewayTypeIP        = "IP"
)

type SubnetGateway struct {
	// type specifies how the default gateway will be created. `Automatic`
	// specifies that neutron will automatically add a default gateway. This is
	// also the default if no Gateway is specified. `None` specifies that the
	// subnet will not have a default gateway. `IP` specifies that the subnet
	// will use a specific address as the default gateway, which must be
	// specified in `IP`.
	// +kubebuilder:validation:Enum:=None;Automatic;IP
	// +required
	Type SubnetGatewayType `json:"type,omitempty"`

	// ip is the IP address of the default gateway, which must be specified if
	// Type is `IP`. It must be a valid IP address, either IPv4 or IPv6,
	// matching the IPVersion in SubnetResourceSpec.
	// +optional
	IP *IPvAny `json:"ip,omitempty"`
}

type AllocationPool struct {
	// start is the first IP address in the allocation pool.
	// +required
	Start IPvAny `json:"start,omitempty"`

	// end is the last IP address in the allocation pool.
	// +required
	End IPvAny `json:"end,omitempty"`
}

type AllocationPoolStatus struct {
	// start is the first IP address in the allocation pool.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Start string `json:"start,omitempty"`

	// end is the last IP address in the allocation pool.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	End string `json:"end,omitempty"`
}

type HostRoute struct {
	// destination for the additional route.
	// +required
	Destination CIDR `json:"destination,omitempty"`

	// nextHop for the additional route.
	// +required
	NextHop IPvAny `json:"nextHop,omitempty"`
}

type HostRouteStatus struct {
	// destination for the additional route.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Destination string `json:"destination,omitempty"`

	// nextHop for the additional route.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	NextHop string `json:"nextHop,omitempty"`
}
