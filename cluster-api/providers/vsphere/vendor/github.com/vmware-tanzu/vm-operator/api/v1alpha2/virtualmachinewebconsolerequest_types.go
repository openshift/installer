// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VirtualMachineWebConsoleRequestSpec describes the desired state for a web
// console request to a VM.
type VirtualMachineWebConsoleRequestSpec struct {
	// Name is the name of a VM in the same Namespace as this web console
	// request.
	Name string `json:"name"`
	// PublicKey is used to encrypt the status.response. This is expected to be a RSA OAEP public key in X.509 PEM format.
	PublicKey string `json:"publicKey"`
}

// VirtualMachineWebConsoleRequestStatus describes the observed state of the
// request.
type VirtualMachineWebConsoleRequestStatus struct {
	// Response will be the authenticated ticket corresponding to this web console request.
	Response string `json:"response,omitempty"`
	// ExpiryTime is the time at which access via this request will expire.
	ExpiryTime metav1.Time `json:"expiryTime,omitempty"`

	// ProxyAddr describes the host address and optional port used to access
	// the VM's web console.
	//
	// The value could be a DNS entry, IPv4, or IPv6 address, followed by an
	// optional port. For example, valid values include:
	//
	//     DNS
	//         * host.com
	//         * host.com:6443
	//
	//     IPv4
	//         * 1.2.3.4
	//         * 1.2.3.4:6443
	//
	//     IPv6
	//         * 1234:1234:1234:1234:1234:1234:1234:1234
	//         * [1234:1234:1234:1234:1234:1234:1234:1234]:6443
	//         * 1234:1234:1234:0000:0000:0000:1234:1234
	//         * 1234:1234:1234::::1234:1234
	//         * [1234:1234:1234::::1234:1234]:6443
	//
	// In other words, the field may be set to any value that is parsable
	// by Go's https://pkg.go.dev/net#ResolveIPAddr and
	// https://pkg.go.dev/net#ParseIP functions.
	ProxyAddr string `json:"proxyAddr,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VirtualMachineWebConsoleRequest allows the creation of a one-time, web
// console connection to a VM.
type VirtualMachineWebConsoleRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineWebConsoleRequestSpec   `json:"spec,omitempty"`
	Status VirtualMachineWebConsoleRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineWebConsoleRequestList contains a list of
// VirtualMachineWebConsoleRequests.
type VirtualMachineWebConsoleRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineWebConsoleRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&VirtualMachineWebConsoleRequest{},
		&VirtualMachineWebConsoleRequestList{},
	)
}
