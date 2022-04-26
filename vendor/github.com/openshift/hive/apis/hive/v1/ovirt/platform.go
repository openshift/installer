package ovirt

import (
	corev1 "k8s.io/api/core/v1"
)

// Platform stores all the global oVirt configuration
type Platform struct {
	// The target cluster under which all VMs will run
	ClusterID string `json:"ovirt_cluster_id"`
	// CredentialsSecretRef refers to a secret that contains the oVirt account access
	// credentials with fields: ovirt_url, ovirt_username, ovirt_password, ovirt_ca_bundle
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
	// CertificatesSecretRef refers to a secret that contains the oVirt CA certificates
	// necessary for communicating with oVirt.
	CertificatesSecretRef corev1.LocalObjectReference `json:"certificatesSecretRef"`
	// The target storage domain under which all VM disk would be created.
	StorageDomainID string `json:"storage_domain_id"`
	// The target network of all the network interfaces of the nodes. Omitting defaults to ovirtmgmt
	// network which is a default network for evert ovirt cluster.
	NetworkName string `json:"ovirt_network_name,omitempty"`
}
