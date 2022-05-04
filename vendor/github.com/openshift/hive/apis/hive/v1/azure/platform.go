package azure

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// CredentialsSecretRef refers to a secret that contains the Azure account access
	// credentials.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// Region specifies the Azure region where the cluster will be created.
	Region string `json:"region"`

	// BaseDomainResourceGroupName specifies the resource group where the azure DNS zone for the base domain is found
	BaseDomainResourceGroupName string `json:"baseDomainResourceGroupName,omitempty"`
}

//SetBaseDomain parses the baseDomainID and sets the related fields on azure.Platform
func (p *Platform) SetBaseDomain(baseDomainID string) error {
	parts := strings.Split(baseDomainID, "/")
	p.BaseDomainResourceGroupName = parts[4]
	return nil
}
