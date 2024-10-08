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

	// cloudName is the name of the Azure cloud environment which can be used to configure the Azure SDK
	// with the appropriate Azure API endpoints.
	// If empty, the value is equal to "AzurePublicCloud".
	// +optional
	CloudName CloudEnvironment `json:"cloudName,omitempty"`
}

// CloudEnvironment is the name of the Azure cloud environment
// +kubebuilder:validation:Enum="";AzurePublicCloud;AzureUSGovernmentCloud;AzureChinaCloud;AzureGermanCloud
type CloudEnvironment string

const (
	// PublicCloud is the general-purpose, public Azure cloud environment.
	PublicCloud CloudEnvironment = "AzurePublicCloud"

	// USGovernmentCloud is the Azure cloud environment for the US government.
	USGovernmentCloud CloudEnvironment = "AzureUSGovernmentCloud"

	// ChinaCloud is the Azure cloud environment used in China.
	ChinaCloud CloudEnvironment = "AzureChinaCloud"

	// GermanCloud is the Azure cloud environment used in Germany.
	GermanCloud CloudEnvironment = "AzureGermanCloud"
)

// Name returns name that Azure uses for the cloud environment.
// See https://github.com/Azure/go-autorest/blob/ec5f4903f77ed9927ac95b19ab8e44ada64c1356/autorest/azure/environments.go#L13
func (e CloudEnvironment) Name() string {
	return string(e)
}

// SetBaseDomain parses the baseDomainID and sets the related fields on azure.Platform
func (p *Platform) SetBaseDomain(baseDomainID string) error {
	parts := strings.Split(baseDomainID, "/")
	p.BaseDomainResourceGroupName = parts[4]
	return nil
}
