package vsphere

import (
	corev1 "k8s.io/api/core/v1"
)

// Platform stores any global configuration used for vSphere platforms.
type Platform struct {
	// VCenter is the domain name or IP address of the vCenter.
	VCenter string `json:"vCenter"`

	// CredentialsSecretRef refers to a secret that contains the vSphere account access
	// credentials: GOVC_USERNAME, GOVC_PASSWORD fields.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// CertificatesSecretRef refers to a secret that contains the vSphere CA certificates
	// necessary for communicating with the VCenter.
	CertificatesSecretRef corev1.LocalObjectReference `json:"certificatesSecretRef"`

	// Datacenter is the name of the datacenter to use in the vCenter.
	Datacenter string `json:"datacenter"`

	// DefaultDatastore is the default datastore to use for provisioning volumes.
	DefaultDatastore string `json:"defaultDatastore"`

	// Folder is the name of the folder that will be used and/or created for
	// virtual machines.
	Folder string `json:"folder,omitempty"`

	// Cluster is the name of the cluster virtual machines will be cloned into.
	Cluster string `json:"cluster,omitempty"`

	// Network specifies the name of the network to be used by the cluster.
	Network string `json:"network,omitempty"`
}
