package v2022_09_02_preview

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/fleetmembers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/managedclustersnapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/resolveprivatelinkserviceid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/trustedaccess"
)

type Client struct {
	AgentPools                  *agentpools.AgentPoolsClient
	FleetMembers                *fleetmembers.FleetMembersClient
	Fleets                      *fleets.FleetsClient
	MaintenanceConfigurations   *maintenanceconfigurations.MaintenanceConfigurationsClient
	ManagedClusterSnapshots     *managedclustersnapshots.ManagedClusterSnapshotsClient
	ManagedClusters             *managedclusters.ManagedClustersClient
	PrivateEndpointConnections  *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources        *privatelinkresources.PrivateLinkResourcesClient
	ResolvePrivateLinkServiceId *resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient
	Snapshots                   *snapshots.SnapshotsClient
	TrustedAccess               *trustedaccess.TrustedAccessClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	agentPoolsClient := agentpools.NewAgentPoolsClientWithBaseURI(endpoint)
	configureAuthFunc(&agentPoolsClient.Client)

	fleetMembersClient := fleetmembers.NewFleetMembersClientWithBaseURI(endpoint)
	configureAuthFunc(&fleetMembersClient.Client)

	fleetsClient := fleets.NewFleetsClientWithBaseURI(endpoint)
	configureAuthFunc(&fleetsClient.Client)

	maintenanceConfigurationsClient := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI(endpoint)
	configureAuthFunc(&maintenanceConfigurationsClient.Client)

	managedClusterSnapshotsClient := managedclustersnapshots.NewManagedClusterSnapshotsClientWithBaseURI(endpoint)
	configureAuthFunc(&managedClusterSnapshotsClient.Client)

	managedClustersClient := managedclusters.NewManagedClustersClientWithBaseURI(endpoint)
	configureAuthFunc(&managedClustersClient.Client)

	privateEndpointConnectionsClient := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(endpoint)
	configureAuthFunc(&privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(endpoint)
	configureAuthFunc(&privateLinkResourcesClient.Client)

	resolvePrivateLinkServiceIdClient := resolveprivatelinkserviceid.NewResolvePrivateLinkServiceIdClientWithBaseURI(endpoint)
	configureAuthFunc(&resolvePrivateLinkServiceIdClient.Client)

	snapshotsClient := snapshots.NewSnapshotsClientWithBaseURI(endpoint)
	configureAuthFunc(&snapshotsClient.Client)

	trustedAccessClient := trustedaccess.NewTrustedAccessClientWithBaseURI(endpoint)
	configureAuthFunc(&trustedAccessClient.Client)

	return Client{
		AgentPools:                  &agentPoolsClient,
		FleetMembers:                &fleetMembersClient,
		Fleets:                      &fleetsClient,
		MaintenanceConfigurations:   &maintenanceConfigurationsClient,
		ManagedClusterSnapshots:     &managedClusterSnapshotsClient,
		ManagedClusters:             &managedClustersClient,
		PrivateEndpointConnections:  &privateEndpointConnectionsClient,
		PrivateLinkResources:        &privateLinkResourcesClient,
		ResolvePrivateLinkServiceId: &resolvePrivateLinkServiceIdClient,
		Snapshots:                   &snapshotsClient,
		TrustedAccess:               &trustedAccessClient,
	}
}
