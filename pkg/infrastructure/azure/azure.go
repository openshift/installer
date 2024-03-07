package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

var _ clusterapi.Provider = (*Provider)(nil)

// Provider implements Azure CAPI installation.
type Provider struct{}

// Name gives the name of the provider, Azure.
func (*Provider) Name() string { return azuretypes.Name }

func (*Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) {
	if in.InstallConfig.Config.Azure.ResourceGroupName == "" {
		createResourceGroup(ctx, in.InstallConfig, in.InfraID)
	}
}

// createResourceGroup creates the resource group required for Azure installation.
func createResourceGroup(ctx context.Context, ic *installconfig.InstallConfig, infraID string) error {
	rgName := fmt.Sprintf("%s-rg", infraID)
	managedBy := ic.Config.Azure.ManagedBy
	session, err := ic.Azure.Session()
	if err != nil {
		return err
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	rgClient, err := armresources.NewResourceGroupsClient(session.Credentials.SubscriptionID, cred, nil)
	if err != nil {
		return err
	}
	param := armresources.ResourceGroup{
		Location:  to.StringPtr(ic.Config.Azure.Region),
		ManagedBy: &managedBy,
	}

	_, err = rgClient.CreateOrUpdate(ctx, rgName, param, nil)
	return err
}
