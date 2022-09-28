package ibmcloud

import (
	corev1 "k8s.io/api/core/v1"
)

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// CredentialsSecretRef refers to a secret that contains IBM Cloud account access
	// credentials.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// AccountID is the IBM Cloud Account ID
	AccountID string `json:"accountID"`

	// CISInstanceCRN is the IBM Cloud Internet Services Instance CRN
	CISInstanceCRN string `json:"cisInstanceCRN"`

	// Region specifies the IBM Cloud region where the cluster will be
	// created.
	Region string `json:"region"`
}
