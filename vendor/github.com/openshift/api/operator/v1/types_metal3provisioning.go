package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Metal3Provisioning contains configuration used by the Provisioning
// service (Ironic) to provision baremetal hosts.
//
// Metal3Provisioning is created by the Openshift installer using admin
// or user provided information about the provisioning network and the NIC
// on the server that can be used to PXE boot it.
//
// This CR is a singleton, created by the installer and currently only
// consumed by the machine-api-operator to bring up and update containers
// in a metal3 cluster.
//
type Metal3Provisioning struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of the desired behavior of the
	// Metal3Provisioning.
	Spec Metal3ProvisioningSpec `json:"spec,omitempty"`

	// status is the most recently observed status of the
	// Metal3Provisioning.
	Status Metal3ProvisioningStatus `json:"status,omitempty"`
}

// ProvisioningDHCP represents just the configuration required to fully
// identify the way DHCP would be handled during baremetal host bringup.
//
// DHCP services could be provided external to the metal3 cluster, in
// which case, IP address assignment for the baremetal servers should
// happen via this external DHCP server and not via a DHCP server started
// within the metal3 cluster.
// If IP address assignment needs to happen via the DHCP server within the
// metal3 cluster, then the CR needs to contain the DHCP address range that
// this internal DHCP server needs to use.
//
type ProvisioningDHCP struct {
	// ManagementState within the OperatorSpec needs to be set to
	// indicate if the DHCP server is internal or external to the
	// metal3 cluster. ManagementState set to "Removed" indicates
	// that the DHCP server is outside the metal3 cluster. And a
	// value of "Managed" indicates that the DHCP services are
	// managed within the metal3 cluster.
	// The other fields of OperatorSpec retain their existing
	// semantics.
	OperatorSpec `json:",inline"`

	// If the ManagementState within the OperatorStatus is set to
	// "Managed", then the DHCPRange represents the range of IP addresses
	// that the DHCP server running within the metal3 cluster can use
	// while provisioning baremetal servers. If the value of ManagementState
	// is set to "Removed", then the value of DHCPRange will be ignored.
	// If the ManagementState is "Managed" and the value of DHCPRange is
	// not set, then the DHCP range is taken to be the default range which
	// goes from .10 to .100 of the ProvisioningNetworkCIDR. This is the only
	// value in all of the provisioning configuration that can be changed
	// after the installer has created the CR.
	DHCPRange string `json:"DHCPRange,omitempty"`
}

// Metal3ProvisioningSpec is the specification of the desired behavior of the
// Metal3Provisioning.
type Metal3ProvisioningSpec struct {
	// provisioningInterface is the name of the network interface on a Baremetal
	// server to the provisioning network. It can have values like "eth1" or "ens3".
	ProvisioningInterface string `json:"provisioningInterface"`

	// provisioningIP is the IP address assigned to the provisioningInterface of
	// the baremetal server. This IP address should be within the provisioning
	// subnet, and outside of the DHCP range.
	ProvisioningIP string `json:"provisioningIP"`

	// provisioningNetworkCIDR is the network on which the baremetal nodes are
	// provisioned. The provisioningIP and the IPs in the dhcpRange all come from
	// within this network.
	ProvisioningNetworkCIDR string `json:"provisioningNetworkCIDR"`

	// provisioningDHCP consists of two parameters that represents whether the DHCP
	// server is internal or external to the metal3 cluster. If it is internal,
	// the second parameter represents the DHCP address range to be provided
	// to the baremetal hosts.
	ProvisioningDHCP ProvisioningDHCP `json:"provisioningDHCP"`
}

// Metal3ProvisioningStatus defines the observed status of Metal3Provisioning.
type Metal3ProvisioningStatus struct {
	OperatorStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// Metal3ProvisioningList contains a list of Metal3Provisioning.
type Metal3ProvisioningList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Metal3Provisioning `json:"items"`
}
