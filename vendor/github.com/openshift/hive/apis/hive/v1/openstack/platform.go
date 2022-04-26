package openstack

import (
	corev1 "k8s.io/api/core/v1"
)

// Platform stores all the global OpenStack configuration
type Platform struct {
	// CredentialsSecretRef refers to a secret that contains the OpenStack account access
	// credentials.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// CertificatesSecretRef refers to a secret that contains CA certificates
	// necessary for communicating with the OpenStack.
	// There is additional configuration required for the OpenShift cluster to trust
	// the certificates provided in this secret.
	// The "clouds.yaml" file included in the credentialsSecretRef Secret must also include
	// a reference to the certificate bundle file for the OpenShift cluster being created to
	// trust the OpenStack endpoints.
	// The "clouds.yaml" file must set the "cacert" field to
	// either "/etc/openstack-ca/<key name containing the trust bundle in credentialsSecretRef Secret>" or
	// "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem".
	//
	// For example,
	// """clouds.yaml
	// clouds:
	//   shiftstack:
	//     auth: ...
	//     cacert: "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"
	// """
	//
	// +optional
	CertificatesSecretRef *corev1.LocalObjectReference `json:"certificatesSecretRef,omitempty"`

	// Cloud will be used to indicate the OS_CLOUD value to use the right section
	// from the clouds.yaml in the CredentialsSecretRef.
	Cloud string `json:"cloud"`

	// TrunkSupport indicates whether or not to use trunk ports in your OpenShift cluster.
	// +optional
	TrunkSupport bool `json:"trunkSupport,omitempty"`
}
