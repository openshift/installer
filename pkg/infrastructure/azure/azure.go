package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/sirupsen/logrus"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/types"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

var _ clusterapi.Provider = (*Provider)(nil)
var _ clusterapi.PostProvider = (*Provider)(nil)

// Provider implements Azure CAPI installation.
type Provider struct{}

// Name gives the name of the provider, Azure.
func (*Provider) Name() string { return azuretypes.Name }

// PostProvision provisions an external Load Balancer (when appropriate), and adds configuration
// for the MCS to the CAPI-provisioned internal LB.
func (*Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	ssn, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("error retrieving Azure session: %w", err)
	}

	region := in.InstallConfig.Config.Azure.Region
	subscriptionID := ssn.Credentials.SubscriptionID
	networkClientFactory, err := armnetwork.NewClientFactory(subscriptionID, ssn.TokenCreds, nil)
	if err != nil {
		return fmt.Errorf("error creating network client factory: %w", err)
	}

	lbClient := networkClientFactory.NewLoadBalancersClient()

	intLoadBalancer, err := updateInternalLoadBalancer(ctx, in.InfraID, region, subscriptionID, lbClient)
	if err != nil {
		return fmt.Errorf("failed to update internal load balancer: %w", err)
	}
	logrus.Debugf("updated internal load balancer: %s", *intLoadBalancer.ID)

	if in.InstallConfig.Config.Publish == types.ExternalPublishingStrategy {
		pipClient := networkClientFactory.NewPublicIPAddressesClient()
		publicIP, err := createPublicIP(ctx, in.InfraID, region, pipClient)
		if err != nil {
			return fmt.Errorf("failed to create public ip: %w", err)
		}
		logrus.Debugf("created public ip: %s", *publicIP.ID)

		loadBalancer, err := createExternalLoadBalancer(ctx, in.InfraID, region, subscriptionID, publicIP, lbClient)
		if err != nil {
			return fmt.Errorf("failed to create load balancer: %w", err)
		}
		logrus.Debugf("created load balancer: %s", *loadBalancer.ID)

		vmClient, err := armcompute.NewVirtualMachinesClient(subscriptionID, ssn.TokenCreds, nil)
		if err != nil {
			return fmt.Errorf("error creating vm client: %w", err)
		}
		nicClient, err := armnetwork.NewInterfacesClient(ssn.Credentials.SubscriptionID, ssn.TokenCreds, nil)
		if err != nil {
			return fmt.Errorf("error creating nic client: %w", err)
		}

		vmIDs, err := getControlPlaneIDs(in.Client, in.InstallConfig.Config.ControlPlane.Replicas, in.InfraID)
		if err != nil {
			return fmt.Errorf("failed to get control plane VM IDs: %w", err)
		}

		bap := loadBalancer.Properties.BackendAddressPools[0]
		if err = associateVMToBackendPool(ctx, in.InfraID, vmIDs, bap, vmClient, nicClient); err != nil {
			return fmt.Errorf("failed to associate control plane VMs with external load balancer: %w", err)
		}
	}

	return nil
}

func getControlPlaneIDs(cl client.Client, replicas *int64, infraID string) ([]string, error) {
	res := []string{}
	total := int64(1)
	if replicas != nil {
		total = *replicas
	}
	for i := int64(0); i < total; i++ {
		machineName := fmt.Sprintf("%s-master-%d", infraID, i)
		key := client.ObjectKey{
			Name:      machineName,
			Namespace: capiutils.Namespace,
		}
		azureMachine := &capz.AzureMachine{}
		if err := cl.Get(context.Background(), key, azureMachine); err != nil {
			return nil, fmt.Errorf("failed to get AzureMahcine: %w", err)
		}
		if vmID := azureMachine.Spec.ProviderID; vmID != nil && len(*vmID) != 0 {
			res = append(res, *azureMachine.Spec.ProviderID)
		} else {
			return nil, fmt.Errorf("%s .Spec.ProviderID is empty", machineName)
		}
	}

	bootstrapName := capiutils.GenerateBoostrapMachineName(infraID)
	key := client.ObjectKey{
		Name:      bootstrapName,
		Namespace: capiutils.Namespace,
	}
	azureMachine := &capz.AzureMachine{}
	if err := cl.Get(context.Background(), key, azureMachine); err != nil {
		return nil, fmt.Errorf("failed to get AzureMachine: %w", err)
	}
	if vmID := azureMachine.Spec.ProviderID; vmID != nil && len(*vmID) != 0 {
		res = append(res, *azureMachine.Spec.ProviderID)
	} else {
		return nil, fmt.Errorf("%s .Spec.ProviderID is empty", bootstrapName)
	}
	return res, nil
}
