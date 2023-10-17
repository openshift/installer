// Copyright (c) 2021-2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WebConsoleRequestSpec describes the specification for used to request a web console request.
type WebConsoleRequestSpec struct {
	// VirtualMachineName is the VM in the same namespace, for which the web console is requested.
	VirtualMachineName string `json:"virtualMachineName"`
	// PublicKey is used to encrypt the status.response. This is expected to be a RSA OAEP public key in X.509 PEM format.
	PublicKey string `json:"publicKey"`
}

// WebConsoleRequestStatus defines the observed state, which includes the web console request itself.
type WebConsoleRequestStatus struct {
	// Response will be the authenticated ticket corresponding to this web console request.
	Response string `json:"response,omitempty"`
	// ExpiryTime is when the ticket referenced in Response will expire.
	ExpiryTime metav1.Time `json:"expiryTime,omitempty"`
	// ProxyAddr describes the host address and optional port used to access
	// the VM's web console.
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

// WebConsoleRequest allows the creation of a one-time web console ticket that can be used to interact with the VM.
type WebConsoleRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebConsoleRequestSpec   `json:"spec,omitempty"`
	Status WebConsoleRequestStatus `json:"status,omitempty"`
}

func (s *WebConsoleRequest) NamespacedName() string {
	return s.Namespace + "/" + s.Name
}

// +kubebuilder:object:root=true

// WebConsoleRequestList contains a list of WebConsoleRequests.
type WebConsoleRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebConsoleRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebConsoleRequest{}, &WebConsoleRequestList{})
}
