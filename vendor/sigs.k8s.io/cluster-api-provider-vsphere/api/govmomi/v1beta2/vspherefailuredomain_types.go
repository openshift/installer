/*
Copyright 2025 The Kubernetes Authors.

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
	// region defines the name and type of a region
	// +required
	Region FailureDomain `json:"region,omitzero"`

	// zone defines the name and type of a zone
	// +required
	Zone FailureDomain `json:"zone,omitzero"`

	// topology describes a given failure domain using vSphere constructs
	// +required
	Topology Topology `json:"topology,omitzero"`
}

// FailureDomain contains data to identify and configure a failure domain.
type FailureDomain struct {
	// name is the name of the tag that represents this failure domain
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Name string `json:"name,omitempty"`

	// type is the type of failure domain, the current values are "Datacenter", "ComputeCluster" and "HostGroup"
	// +required
	// +kubebuilder:validation:Enum=Datacenter;ComputeCluster;HostGroup
	Type FailureDomainType `json:"type,omitempty"`

	// tagCategory is the category used for the tag
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	TagCategory string `json:"tagCategory,omitempty"`
}

// Topology describes a given failure domain using vSphere constructs.
type Topology struct {
	// datacenter as the failure domain.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Datacenter string `json:"datacenter,omitempty"`

	// computeCluster as the failure domain
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	ComputeCluster string `json:"computeCluster,omitempty"`

	// hosts has information required for placement of machines on VSphere hosts.
	// +optional
	Hosts FailureDomainHosts `json:"hosts,omitempty,omitzero"`

	// networks is the list of networks within this failure domain
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=2048
	Networks []string `json:"networks,omitempty"`

	// networkConfigurations is a list of network configurations within this failure domain.
	// +optional
	// +listType=map
	// +listMapKey=networkName
	// +kubebuilder:validation:MaxItems=128
	NetworkConfigurations []NetworkConfiguration `json:"networkConfigurations,omitempty"`

	// datastore is the name or inventory path of the datastore in which the
	// virtual machine is created/located.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Datastore string `json:"datastore,omitempty"`
}

// NetworkConfiguration defines a network configuration that should be used when consuming
// a failure domain.
type NetworkConfiguration struct {
	// networkName is the network name for this machine's VM.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	NetworkName string `json:"networkName,omitempty"`

	// dhcp4 is a flag that indicates whether or not to use DHCP for IPv4.
	// +optional
	DHCP4 *bool `json:"dhcp4,omitempty"`

	// dhcp6 is a flag that indicates whether or not to use DHCP for IPv6.
	// +optional
	DHCP6 *bool `json:"dhcp6,omitempty"`

	// nameservers is a list of IPv4 and/or IPv6 addresses used as DNS
	// nameservers.
	// Please note that Linux allows only three nameservers (https://linux.die.net/man/5/resolv.conf).
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=64
	Nameservers []string `json:"nameservers,omitempty"`

	// searchDomains is a list of search domains used when resolving IP
	// addresses with DNS.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	// +kubebuilder:validation:items:MinLength=1
	// +kubebuilder:validation:items:MaxLength=1024
	SearchDomains []string `json:"searchDomains,omitempty"`

	// dhcp4Overrides allows for the control over several DHCP behaviors.
	// Overrides will only be applied when the corresponding DHCP flag is set.
	// Only configured values will be sent, omitted values will default to
	// distribution defaults.
	// Dependent on support in the network stack for your distribution.
	// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
	// +optional
	DHCP4Overrides *DHCPOverrides `json:"dhcp4Overrides,omitempty"`

	// dhcp6Overrides allows for the control over several DHCP behaviors.
	// Overrides will only be applied when the corresponding DHCP flag is set.
	// Only configured values will be sent, omitted values will default to
	// distribution defaults.
	// Dependent on support in the network stack for your distribution.
	// For more information see the netplan reference (https://netplan.io/reference#dhcp-overrides)
	// +optional
	DHCP6Overrides *DHCPOverrides `json:"dhcp6Overrides,omitempty"`

	// addressesFromPools is a list of IPAddressPools that should be assigned
	// to IPAddressClaims. The machine's cloud-init metadata will be populated
	// with IPAddresses fulfilled by an IPAM provider.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=128
	AddressesFromPools []IPPoolReference `json:"addressesFromPools,omitempty"`
}

// FailureDomainHosts has information required for placement of machines on VSphere hosts.
type FailureDomainHosts struct {
	// vmGroupName is the name of the VM group
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	VMGroupName string `json:"vmGroupName,omitempty"`

	// hostGroupName is the name of the Host group
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	HostGroupName string `json:"hostGroupName,omitempty"`
}

// IsDefined returns true if the ref is defined.
func (m *FailureDomainHosts) IsDefined() bool {
	return m.VMGroupName != "" || m.HostGroupName != ""
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=vspherefailuredomains,scope=Cluster,categories=cluster-api

// VSphereFailureDomain is the Schema for the vspherefailuredomains API.
type VSphereFailureDomain struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereFailureDomain.
	// +required
	Spec VSphereFailureDomainSpec `json:"spec,omitempty,omitzero"`
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
