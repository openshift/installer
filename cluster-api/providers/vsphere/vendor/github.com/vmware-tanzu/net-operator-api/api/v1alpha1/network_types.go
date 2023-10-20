// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkProviderReference contains info to locate a network provider object.
type NetworkProviderReference struct {
	// APIGroup is the group for the resource being referenced.
	APIGroup string `json:"apiGroup"`
	// Kind is the type of resource being referenced.
	Kind string `json:"kind"`
	// Name is the name of resource being referenced.
	Name string `json:"name"`
	// Namespace of the resource being referenced. If empty, cluster scoped resource is assumed.
	Namespace string `json:"namespace,omitempty"`
	// API version of the referent.
	APIVersion string `json:"apiVersion,omitempty"`
}

// NetworkType is used to type the constants describing possible network types.
type NetworkType string

const (
	// NetworkTypeNSXT is the network type describing NSX-T.
	NetworkTypeNSXT = NetworkType("nsx-t")

	// NetworkTypeVDS is the network type describing VSphere Distributed Switch.
	NetworkTypeVDS = NetworkType("vsphere-distributed")
)

// NetworkSpec defines the state of Network.
type NetworkSpec struct {
	// Type describes type of Network. Supported values are nsx-t, vsphere-distributed.
	Type NetworkType `json:"type"`
	// ProviderRef is reference to a network provider object that provides this type of network.
	ProviderRef NetworkProviderReference `json:"providerRef"`
	// DNS is a list of DNS server IPs to associate with network interfaces on this network.
	DNS []string `json:"dns,omitempty"`
	// DNSSearchDomains is a list of DNS search domains to associate with network interfaces on this network.
	DNSSearchDomains []string `json:"dnsSearchDomains,omitempty"`
	// NTP is a list of NTP server DNS names or IP addresses to use on this network.
	NTP []string `json:"ntp,omitempty"`
}

// NetworkStatus is unused. This is because Network is purely a configuration resource.
type NetworkStatus struct {
}

// +genclient
// +kubebuilder:object:root=true

// Network is the Schema for the networks API.
// A Network describes type, class and common attributes of a network available
// in a namespace. A NetworkInterface resource references a Network.
type Network struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkSpec   `json:"spec,omitempty"`
	Status NetworkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkList contains a list of Network
type NetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Network `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&Network{}, &NetworkList{})
}
