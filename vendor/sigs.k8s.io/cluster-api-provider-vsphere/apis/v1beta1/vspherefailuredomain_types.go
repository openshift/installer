/*
Copyright 2021 The Kubernetes Authors.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FailureDomainType defines the VCenter object the failure domain represents.
type FailureDomainType string

const (
	// HostGroupFailureDomain is a failure domain for a host group.
	HostGroupFailureDomain FailureDomainType = "HostGroup"
	// ComputeClusterFailureDomain is a failure domain for a compute cluster.
	ComputeClusterFailureDomain FailureDomainType = "ComputeCluster"
	// DatacenterFailureDomain is a failure domain for a datacenter.
	DatacenterFailureDomain FailureDomainType = "Datacenter"
)

// VSphereFailureDomainSpec defines the desired state of VSphereFailureDomain.
type VSphereFailureDomainSpec struct {

	// Region defines the name and type of a region
	Region FailureDomain `json:"region"`

	// Zone defines the name and type of a zone
	Zone FailureDomain `json:"zone"`

	// Topology describes a given failure domain using vSphere constructs
	Topology Topology `json:"topology"`
}

// FailureDomain contains data to identify and configure a failure domain.
type FailureDomain struct {
	// Name is the name of the tag that represents this failure domain
	Name string `json:"name"`

	// Type is the type of failure domain, the current values are "Datacenter", "ComputeCluster" and "HostGroup"
	// +kubebuilder:validation:Enum=Datacenter;ComputeCluster;HostGroup
	Type FailureDomainType `json:"type"`

	// TagCategory is the category used for the tag
	TagCategory string `json:"tagCategory"`

	// AutoConfigure tags the Type which is specified in the Topology
	//
	// Deprecated: This field is going to be removed in a future release.
	AutoConfigure *bool `json:"autoConfigure,omitempty"`
}

// Topology describes a given failure domain using vSphere constructs.
type Topology struct {
	// Datacenter as the failure domain.
	// +kubebuilder:validation:Required
	Datacenter string `json:"datacenter"`

	// ComputeCluster as the failure domain
	// +optional
	ComputeCluster *string `json:"computeCluster,omitempty"`

	// Hosts has information required for placement of machines on VSphere hosts.
	// +optional
	Hosts *FailureDomainHosts `json:"hosts,omitempty"`

	// Networks is the list of networks within this failure domain
	// +optional
	Networks []string `json:"networks,omitempty"`

	// NetworkConfigurations is a list of network configurations within this failure domain.
	// +optional
	// +listType=map
	// +listMapKey=networkName
	NetworkConfigurations []NetworkConfiguration `json:"networkConfigurations,omitempty"`

	// Datastore is the name or inventory path of the datastore in which the
	// virtual machine is created/located.
	// +optional
	Datastore string `json:"datastore,omitempty"`
}

// NetworkConfiguration defines a network configuration that should be used when consuming
// a failure domain.
type NetworkConfiguration struct {
	// NetworkName is the network name for this machine's VM.
	// +kubebuilder:validation:Required
	NetworkName string `json:"networkName"`

	// DHCP4 is a flag that indicates whether or not to use DHCP for IPv4.
	// +optional
	DHCP4 *bool `json:"dhcp4,omitempty"`

	// DHCP6 is a flag that indicates whether or not to use DHCP for IPv6.
	// +optional
	DHCP6 *bool `json:"dhcp6,omitempty"`

	// Nameservers is a list of IPv4 and/or IPv6 addresses used as DNS
	// nameservers.
	// Please note that Linux allows only three nameservers (https://linux.die.net/man/5/resolv.conf).
	// +optional
	Nameservers []string `json:"nameservers,omitempty"`

	// SearchDomains is a list of search domains used when resolving IP
	// addresses with DNS.
	// +optional
	SearchDomains []string `json:"searchDomains,omitempty"`

	// DHCP4Overrides allows for the control over several DHCP behaviors.
	// Overrides will only be applied when the corresponding DHCP flag is set.
	// Only configured values will be sent, omitted values will default to
	// distribution defaults.
	// Dependent on support in the network stack for your distribution.
	// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
	// +optional
	DHCP4Overrides *DHCPOverrides `json:"dhcp4Overrides,omitempty"`

	// DHCP6Overrides allows for the control over several DHCP behaviors.
	// Overrides will only be applied when the corresponding DHCP flag is set.
	// Only configured values will be sent, omitted values will default to
	// distribution defaults.
	// Dependent on support in the network stack for your distribution.
	// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
	// +optional
	DHCP6Overrides *DHCPOverrides `json:"dhcp6Overrides,omitempty"`

	// AddressesFromPools is a list of IPAddressPools that should be assigned
	// to IPAddressClaims. The machine's cloud-init metadata will be populated
	// with IPAddresses fulfilled by an IPAM provider.
	// +optional
	AddressesFromPools []corev1.TypedLocalObjectReference `json:"addressesFromPools,omitempty"`
}

// FailureDomainHosts has information required for placement of machines on VSphere hosts.
type FailureDomainHosts struct {
	// VMGroupName is the name of the VM group
	VMGroupName string `json:"vmGroupName"`

	// HostGroupName is the name of the Host group
	HostGroupName string `json:"hostGroupName"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=vspherefailuredomains,scope=Cluster,categories=cluster-api

// VSphereFailureDomain is the Schema for the vspherefailuredomains API.
type VSphereFailureDomain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VSphereFailureDomainSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereFailureDomainList contains a list of VSphereFailureDomain.
type VSphereFailureDomainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereFailureDomain `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VSphereFailureDomain{}, &VSphereFailureDomainList{})
}
