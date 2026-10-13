package ibmcloud

import (
	corev1 "k8s.io/api/core/v1"
)

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// CredentialsSecretRef refers to a secret that contains IBM Cloud account access
	// credentials.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
	// AccountID is the IBM Cloud Account ID.
	// AccountID is DEPRECATED and is gathered via the IBM Cloud API for the provided
	// credentials. This field will be ignored.
	// +optional
	AccountID string `json:"accountID,omitempty"`
	// CISInstanceCRN is the IBM Cloud Internet Services Instance CRN
	// CISInstanceCRN is DEPRECATED and gathered via the IBM Cloud API for the provided
	// credentials and cluster deployment base domain. This field will be ignored.
	// +optional
	CISInstanceCRN string `json:"cisInstanceCRN,omitempty"`
	// Region specifies the IBM Cloud region where the cluster will be
	// created.
	Region string `json:"region"`
}
