package credentialsrequest

import (
	"k8s.io/apimachinery/pkg/runtime"

	ccov1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	"github.com/openshift/installer/pkg/types/azure"
)

// AzureProviderSpec holds the decoded Azure-specific provider spec
// from a CredentialsRequest.
type AzureProviderSpec struct {
	// RoleBindings contains a list of roles that should be associated with the minted credential.
	RoleBindings []ccov1.RoleBinding `json:"roleBindings"`

	// Permissions is the list of Azure permissions required to create a more fine-grained custom role to
	// satisfy the CredentialsRequest.
	// The Permissions field may be provided in addition to RoleBindings. When both fields are specified,
	// the user-assigned managed identity will have union of permissions defined from both Permissions
	// and RoleBindings.
	Permissions []string `json:"permissions,omitempty"`

	// DataPermissions is the list of Azure data permissions required to create a more fine-grained custom
	// role to satisfy the CredentialsRequest.
	// The DataPermissions field may be provided in addition to RoleBindings. When both fields are specified,
	// the user-assigned managed identity will have union of permissions defined from both DataPermissions
	// and RoleBindings.
	DataPermissions []string `json:"dataPermissions,omitempty"`

	// The following fields are only required for Azure Workload Identity.
	// AzureClientID is the ID of the specific application you created in Azure
	AzureClientID string `json:"azureClientID,omitempty"`

	// AzureRegion is the geographic region of the Azure service.
	AzureRegion string `json:"azureRegion,omitempty"`

	// Each Azure subscription has an ID associated with it, as does the tenant to which a subscription belongs.
	// AzureSubscriptionID is the ID of the subscription.
	AzureSubscriptionID string `json:"azureSubscriptionID,omitempty"`

	// AzureTenantID is the ID of the tenant to which the subscription belongs.
	AzureTenantID string `json:"azureTenantID,omitempty"`
}

func init() {
	registerProviderSpecDecoder(azure.Name, decodeAzureProviderSpec)
}

func decodeAzureProviderSpec(raw *runtime.RawExtension) (interface{}, error) {
	azureSpec := &ccov1.AzureProviderSpec{}
	if err := ccov1.Codec.DecodeProviderSpec(raw, azureSpec); err != nil {
		return nil, err
	}
	return &AzureProviderSpec{
		RoleBindings:        azureSpec.RoleBindings,
		Permissions:         azureSpec.Permissions,
		DataPermissions:     azureSpec.DataPermissions,
		AzureClientID:       azureSpec.AzureClientID,
		AzureRegion:         azureSpec.AzureRegion,
		AzureSubscriptionID: azureSpec.AzureSubscriptionID,
		AzureTenantID:       azureSpec.AzureTenantID,
	}, nil
}
