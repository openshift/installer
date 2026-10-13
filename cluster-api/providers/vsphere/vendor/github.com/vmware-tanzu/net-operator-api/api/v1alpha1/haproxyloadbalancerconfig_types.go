// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HAProxyLoadBalancerConfigSpec defines the configuration for an HAProxyLoadBalancerConfig instance.
// The spec is used to configure the HAProxyLoadBalancer instance to correctly route traffic to services.
// This spec supports HAProxyLoadBalancerConfig Dataplane API 2.0+ sidecar
type HAProxyLoadBalancerConfigSpec struct {
	// EndPointURLs is a list of the addresses for the DataPlane API servers used
	// to configure HAProxy.
	// One or more DataPlane API endpoints are possible due to the following topologies:
	// Single Node Topology
	// Multi-Node Active/Passive Topology
	// The strings should include the host, port, and API version, ex.:
	// https://hostname:port/v1
	// +kubebuilder:validation:MinItems=1
	EndPointURLs []string `json:"endPointURLs"`

	// ServerName is used to verify the hostname on the returned
	// certificates. It is also included
	// in the client's handshake to support virtual hosting unless it is
	// an IP address.
	// Defaults to the host part parsed from Server
	// +optional
	ServerName string `json:"serverName,omitempty"`

	// CredentialSecretRef is an object name of kind Secret.
	// It will be used to access and configure the HAProxy load balancer DataPlane API servers.
	// The following fields are optional:
	//
	// * certificateAuthorityData - CertificateAuthorityData contains PEM-encoded certificate authority certificates.
	//
	// * clientCertificateData - ClientCertificateData contains PEM-encoded data from a client cert file.
	//
	// * clientKeyData - ClientKeyData contains PEM-encoded data from a client key file for TLS.
	//
	// * username - Username is the username for basic authentication. Defaults to "client".
	//
	// * password - Password is the password for basic authentication. Defaults to "cert".
	//
	// Sample of a secret:
	//
	// apiVersion: v1
	// kind: Secret
	// metadata:
	// name: haproxy-lb-config
	// namespace: vmware-system-netop
	// data:
	// 	 certificateAuthorityData: <base64_Encoded>
	// 	 clientCertificateData: <base64_Encoded>
	// 	 clientKeyData: <base64_Encoded>
	//   username: <base64_Encoded>
	//   password: <base64_Encoded>
	// +optional
	CredentialSecretRef ClientSecretReference `json:"credentialSecretRef,omitempty"`
}

// HAProxyLoadBalancerConfigStatus is unused. This is because HAProxyLoadBalancerConfig is purely a configuration resource
type HAProxyLoadBalancerConfigStatus struct {
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// HAProxyLoadBalancerConfig is the Schema for the HAProxyLoadBalancerConfigs API
type HAProxyLoadBalancerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HAProxyLoadBalancerConfigSpec   `json:"spec,omitempty"`
	Status HAProxyLoadBalancerConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HAProxyLoadBalancerConfigList contains a list of HAProxyLoadBalancerConfig
type HAProxyLoadBalancerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HAProxyLoadBalancerConfig `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&HAProxyLoadBalancerConfig{}, &HAProxyLoadBalancerConfigList{})
}
