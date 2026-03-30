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

// +kubebuilder:validation:Enum:=ingress;egress
type RuleDirection string

// +kubebuilder:validation:Enum:=ah;dccp;egp;esp;gre;icmp;icmpv6;igmp;ipip;ipv6-encap;ipv6-frag;ipv6-icmp;ipv6-nonxt;ipv6-opts;ipv6-route;ospf;pgm;rsvp;sctp;tcp;udp;udplite;vrrp
type Protocol string

const (
	ProtocolAH        Protocol = "ah"
	ProtocolDCCP      Protocol = "dccp"
	ProtocolEGP       Protocol = "egp"
	ProtocolESP       Protocol = "esp"
	ProtocolGRE       Protocol = "gre"
	ProtocolICMP      Protocol = "icmp"
	ProtocolICMPV6    Protocol = "icmpv6"
	ProtocolIGMP      Protocol = "igmp"
	ProtocolIPIP      Protocol = "ipip"
	ProtocolIPV6ENCAP Protocol = "ipv6-encap"
	ProtocolIPV6FRAG  Protocol = "ipv6-frag"
	ProtocolIPV6ICMP  Protocol = "ipv6-icmp"
	ProtocolIPV6NONXT Protocol = "ipv6-nonxt"
	ProtocolIPV6OPTS  Protocol = "ipv6-opts"
	ProtocolIPV6ROUTE Protocol = "ipv6-route"
	ProtocolOSPF      Protocol = "ospf"
	ProtocolPGM       Protocol = "pgm"
	ProtocolRSVP      Protocol = "rsvp"
	ProtocolSCTP      Protocol = "sctp"
	ProtocolTCP       Protocol = "tcp"
	ProtocolUDP       Protocol = "udp"
	ProtocolUDPLITE   Protocol = "udplite"
	ProtocolVRRP      Protocol = "vrrp"
)

// +kubebuilder:validation:Enum:=IPv4;IPv6
type Ethertype string

const (
	EthertypeIPv4 Ethertype = "IPv4"
	EthertypeIPv6 Ethertype = "IPv6"
)

// +kubebuilder:validation:Minimum:=0
// +kubebuilder:validation:Maximum:=65535
type PortNumber int32

type PortRangeSpec struct {
	// min is the minimum port number in the range that is matched by the security group rule.
	// If the protocol is TCP, UDP, DCCP, SCTP or UDP-Lite this value must be less than or equal
	// to the port_range_max attribute value. If the protocol is ICMP, this value must be an ICMP type
	// +required
	Min PortNumber `json:"min"`
	// max is the maximum port number in the range that is matched by the security group rule.
	// If the protocol is TCP, UDP, DCCP, SCTP or UDP-Lite this value must be greater than or equal
	// to the port_range_min attribute value. If the protocol is ICMP, this value must be an ICMP code.
	// +required
	Max PortNumber `json:"max"`
}

type PortRangeStatus struct {
	// min is the minimum port number in the range that is matched by the security group rule.
	// If the protocol is TCP, UDP, DCCP, SCTP or UDP-Lite this value must be less than or equal
	// to the port_range_max attribute value. If the protocol is ICMP, this value must be an ICMP type
	// +optional
	Min int32 `json:"min"`
	// max is the maximum port number in the range that is matched by the security group rule.
	// If the protocol is TCP, UDP, DCCP, SCTP or UDP-Lite this value must be greater than or equal
	// to the port_range_min attribute value. If the protocol is ICMP, this value must be an ICMP code.
	// +optional
	Max int32 `json:"max"`
}

// NOTE: A validation was removed from SecurityGroupRule until we bump minimum k8s to at least v1.31:
// - remoteIPPrefix matches the address family defined in ethertype: PR #336

// SecurityGroupRule defines a Security Group rule
// +kubebuilder:validation:MinProperties:=1
// +kubebuilder:validation:XValidation:rule="(!has(self.portRange)|| !(self.protocol == 'tcp'|| self.protocol == 'udp' || self.protocol == 'dccp' || self.protocol == 'sctp' || self.protocol == 'udplite') || (self.portRange.min <= self.portRange.max))",message="portRangeMax should be equal or greater than portRange.min"
// +kubebuilder:validation:XValidation:rule="!(self.protocol == 'icmp' || self.protocol == 'icmpv6') || !has(self.portRange)|| (self.portRange.min >= 0 && self.portRange.min <= 255)",message="When protocol is ICMP or ICMPv6 portRange.min should be between 0 and 255"
// +kubebuilder:validation:XValidation:rule="!(self.protocol == 'icmp' || self.protocol == 'icmpv6') || !has(self.portRange)|| (self.portRange.max >= 0 && self.portRange.max <= 255)",message="When protocol is ICMP or ICMPv6 portRange.max should be between 0 and 255"
type SecurityGroupRule struct {
	// description is a human-readable description for the resource.
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// direction represents the direction in which the security group rule
	// is applied. Can be ingress or egress.
	// +optional
	Direction *RuleDirection `json:"direction,omitempty"`

	// remoteIPPrefix is an IP address block. Should match the Ethertype (IPv4 or IPv6)
	// +optional
	RemoteIPPrefix *CIDR `json:"remoteIPPrefix,omitempty"`

	// protocol is the IP protocol is represented by a string
	// +optional
	Protocol *Protocol `json:"protocol,omitempty"`

	// ethertype must be IPv4 or IPv6, and addresses represented in CIDR
	// must match the ingress or egress rules.
	// +required
	Ethertype Ethertype `json:"ethertype,omitempty"`

	// portRange sets the minimum and maximum ports range that the security group rule
	// matches. If the protocol is [tcp, udp, dccp sctp,udplite] PortRange.Min must be less than
	// or equal to the PortRange.Max attribute value.
	// If the protocol is ICMP, this PortRamge.Min must be an ICMP code and PortRange.Max
	// should be an ICMP type
	// +optional
	PortRange *PortRangeSpec `json:"portRange,omitempty"`
}

type SecurityGroupRuleStatus struct {
	// id is the ID of the security group rule.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ID string `json:"id,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// direction represents the direction in which the security group rule
	// is applied. Can be ingress or egress.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Direction string `json:"direction,omitempty"`

	// RemoteAddressGroupId (Not in gophercloud)

	// remoteGroupID is the remote group UUID to associate with this security group rule
	// RemoteGroupID
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	RemoteGroupID string `json:"remoteGroupID,omitempty"`

	// remoteIPPrefix is an IP address block. Should match the Ethertype (IPv4 or IPv6)
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	RemoteIPPrefix string `json:"remoteIPPrefix,omitempty"`

	// protocol is the IP protocol can be represented by a string, an
	// integer, or null
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Protocol string `json:"protocol,omitempty"`

	// ethertype must be IPv4 or IPv6, and addresses represented in CIDR
	// must match the ingress or egress rules.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Ethertype string `json:"ethertype,omitempty"`

	// portRange sets the minimum and maximum ports range that the security group rule
	// matches. If the protocol is [tcp, udp, dccp sctp,udplite] PortRange.Min must be less than
	// or equal to the PortRange.Max attribute value.
	// If the protocol is ICMP, this PortRamge.Min must be an ICMP code and PortRange.Max
	// should be an ICMP type
	// +optional
	PortRange *PortRangeStatus `json:"portRange,omitempty"`
	// FIXME(mandre) This field is not yet returned by gophercloud
	// BelongsToDefaultSG bool `json:"belongsToDefaultSG,omitempty"`

	// FIXME(mandre) Technically, the neutron status metadata are returned
	// for SG rules. Should we include them? Gophercloud does not
	// implements this yet.
	// NeutronStatusMetadata `json:",inline"`
}

// SecurityGroupResourceSpec contains the desired state of a security group
type SecurityGroupResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// tags is a list of tags which will be applied to the security group.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	Tags []NeutronTag `json:"tags,omitempty"`

	// stateful indicates if the security group is stateful or stateless.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="stateful is immutable"
	Stateful *bool `json:"stateful,omitempty"`

	// rules is a list of security group rules belonging to this SG.
	// +kubebuilder:validation:MaxItems:=256
	// +listType=atomic
	// +optional
	Rules []SecurityGroupRule `json:"rules,omitempty"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="projectRef is immutable"
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`
}

// SecurityGroupFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type SecurityGroupFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description of the existing resource
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

// SecurityGroupResourceStatus represents the observed state of the resource.
type SecurityGroupResourceStatus struct {
	// name is a Human-readable name for the security group. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// projectID is the project owner of the security group.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems:=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`

	// stateful indicates if the security group is stateful or stateless.
	// +optional
	Stateful bool `json:"stateful,omitempty"`

	// rules is a list of security group rules belonging to this SG.
	// +kubebuilder:validation:MaxItems:=256
	// +listType=atomic
	// +optional
	Rules []SecurityGroupRuleStatus `json:"rules,omitempty"`

	NeutronStatusMetadata `json:",inline"`
}
