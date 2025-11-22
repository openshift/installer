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

// PortFilter specifies a filter to select a port. At least one parameter must be specified.
// +kubebuilder:validation:MinProperties:=1
type PortFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description of the existing resource
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// networkRef is a reference to the ORC Network which this port is associated with.
	// +optional
	NetworkRef KubernetesNameRef `json:"networkRef"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

type AllowedAddressPair struct {
	// ip contains an IP address which a server connected to the port can
	// send packets with. It can be an IP Address or a CIDR (if supported
	// by the underlying extension plugin).
	// +required
	IP IPvAny `json:"ip,omitempty"`

	// mac contains a MAC address which a server connected to the port can
	// send packets with. Defaults to the MAC address of the port.
	// +optional
	MAC *MAC `json:"mac,omitempty"`
}

type AllowedAddressPairStatus struct {
	// ip contains an IP address which a server connected to the port can
	// send packets with.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	IP string `json:"ip,omitempty"`

	// mac contains a MAC address which a server connected to the port can
	// send packets with.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	MAC string `json:"mac,omitempty"`
}

type Address struct {
	// ip contains a fixed IP address assigned to the port. It must belong
	// to the referenced subnet's CIDR. If not specified, OpenStack
	// allocates an available IP from the referenced subnet.
	// +optional
	IP *IPvAny `json:"ip,omitempty"`

	// subnetRef references the subnet from which to allocate the IP
	// address.
	// +required
	SubnetRef KubernetesNameRef `json:"subnetRef,omitempty"`
}

type FixedIPStatus struct {
	// ip contains a fixed IP address assigned to the port.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	IP string `json:"ip,omitempty"`

	// subnetID is the ID of the subnet this IP is allocated from.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	SubnetID string `json:"subnetID,omitempty"`
}

// +kubebuilder:validation:XValidation:rule="has(self.portSecurity) && self.portSecurity == 'Disabled' ? !has(self.securityGroupRefs) : true",message="securityGroupRefs must be empty when portSecurity is set to Disabled"
// +kubebuilder:validation:XValidation:rule="has(self.portSecurity) && self.portSecurity == 'Disabled' ? !has(self.allowedAddressPairs) : true",message="allowedAddressPairs must be empty when portSecurity is set to Disabled"
type PortResourceSpec struct {
	// name is a human-readable name of the port. If not set, the object's name will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// networkRef is a reference to the ORC Network which this port is associated with.
	// +required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="networkRef is immutable"
	NetworkRef KubernetesNameRef `json:"networkRef,omitempty"`

	// tags is a list of tags which will be applied to the port.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	Tags []NeutronTag `json:"tags,omitempty"`

	// allowedAddressPairs are allowed addresses associated with this port.
	// +kubebuilder:validation:MaxItems:=128
	// +listType=atomic
	// +optional
	AllowedAddressPairs []AllowedAddressPair `json:"allowedAddressPairs,omitempty"`

	// addresses are the IP addresses for the port.
	// +kubebuilder:validation:MaxItems:=128
	// +listType=atomic
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="addresses is immutable"
	Addresses []Address `json:"addresses,omitempty"`

	// securityGroupRefs are the names of the security groups associated
	// with this port.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	SecurityGroupRefs []OpenStackName `json:"securityGroupRefs,omitempty"`

	// vnicType specifies the type of vNIC which this port should be
	// attached to. This is used to determine which mechanism driver(s) to
	// be used to bind the port. The valid values are normal, macvtap,
	// direct, baremetal, direct-physical, virtio-forwarder, smart-nic and
	// remote-managed, although these values will not be validated in this
	// API to ensure compatibility with future neutron changes or custom
	// implementations. What type of vNIC is actually available depends on
	// deployments. If not specified, the Neutron default value is used.
	// +kubebuilder:validation:MaxLength:=64
	// +optional
	VNICType string `json:"vnicType,omitempty"`

	// portSecurity controls port security for this port.
	// When set to Enabled, port security is enabled.
	// When set to Disabled, port security is disabled and SecurityGroupRefs must be empty.
	// When set to Inherit (default), it takes the value from the network level.
	// +kubebuilder:default=Inherit
	// +optional
	// +kubebuilder:validation:XValidation:rule="!(oldSelf != 'Inherit' && self == 'Inherit')",message="portSecurity cannot be changed to Inherit"
	PortSecurity PortSecurityState `json:"portSecurity,omitempty"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="projectRef is immutable"
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`
}

type PortResourceStatus struct {
	// name is the human-readable name of the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// networkID is the ID of the attached network.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	NetworkID string `json:"networkID,omitempty"`

	// projectID is the project owner of the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// status indicates the current status of the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Status string `json:"status,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`

	// adminStateUp is the administrative state of the port,
	// which is up (true) or down (false).
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// macAddress is the MAC address of the port.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	MACAddress string `json:"macAddress,omitempty"`

	// deviceID is the ID of the device that uses this port.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	DeviceID string `json:"deviceID,omitempty"`

	// deviceOwner is the entity type that uses this port.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	DeviceOwner string `json:"deviceOwner,omitempty"`

	// allowedAddressPairs is a set of zero or more allowed address pair
	// objects each where address pair object contains an IP address and
	// MAC address.
	// +kubebuilder:validation:MaxItems=128
	// +listType=atomic
	// +optional
	AllowedAddressPairs []AllowedAddressPairStatus `json:"allowedAddressPairs,omitempty"`

	// fixedIPs is a set of zero or more fixed IP objects each where fixed
	// IP object contains an IP address and subnet ID from which the IP
	// address is assigned.
	// +kubebuilder:validation:MaxItems=128
	// +listType=atomic
	// +optional
	FixedIPs []FixedIPStatus `json:"fixedIPs,omitempty"`

	// securityGroups contains the IDs of security groups applied to the port.
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	SecurityGroups []string `json:"securityGroups,omitempty"`

	// propagateUplinkStatus represents the uplink status propagation of
	// the port.
	// +optional
	PropagateUplinkStatus *bool `json:"propagateUplinkStatus,omitempty"`

	// vnicType is the type of vNIC which this port is attached to.
	// +kubebuilder:validation:MaxLength:=64
	// +optional
	VNICType string `json:"vnicType,omitempty"`

	// portSecurityEnabled indicates whether port security is enabled or not.
	// +optional
	PortSecurityEnabled *bool `json:"portSecurityEnabled,omitempty"`

	NeutronStatusMetadata `json:",inline"`
}

// PortSecurityState represents the security state of a port
// +kubebuilder:validation:Enum=Enabled;Disabled;Inherit
type PortSecurityState string

const (
	// PortSecurityEnabled means port security is enabled
	PortSecurityEnabled PortSecurityState = "Enabled"
	// PortSecurityDisabled means port security is disabled
	PortSecurityDisabled PortSecurityState = "Disabled"
	// PortSecurityInherit means port security settings are inherited from the network
	PortSecurityInherit PortSecurityState = "Inherit"
)
