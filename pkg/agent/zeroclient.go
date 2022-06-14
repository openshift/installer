package agent

import (
	"context"
	"net"
	"net/url"

	"github.com/go-openapi/strfmt"
	"github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/client/versions"
	"github.com/openshift/assisted-service/models"

	"github.com/openshift/installer/pkg/agent/manifests"
)

type NodeZeroClient interface {
	isAgentAPILive(ctx context.Context, zero *nodeZeroClient) (bool, error)
	getClusterZero(ctx context.Context, zero *nodeZeroClient) (*models.Cluster, error)
}

type nodeZeroClient struct {
	restClient          *client.AssistedInstall
	NodeZeroIP          string
	clusterZeroID       *strfmt.UUID
	clusterZeroInfraEnv *strfmt.UUID
}

type clusterZero struct {
	clusterZeroID       *strfmt.UUID
	clusterZeroInfraEnv *strfmt.UUID
}

func NewNodeZeroClient() (*nodeZeroClient, error) {
	ctx := context.Background()
	n := manifests.NewNMConfig()
	nodeZeroIP := n.GetNodeZeroIP()

	restConfig := client.Config{}
	*restConfig.URL = url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(nodeZeroIP, "8090"),
		Path:   client.DefaultBasePath,
	}
	restClient := client.New(restConfig)

	// GET /v2/clusters
	listClusterParams := installer.NewV2ListClustersParams()
	clusterResult, err := restClient.Installer.V2ListClusters(ctx, listClusterParams)
	if err != nil {
		return nil, err
	}
	clusterZeroID := clusterResult.Payload[0].ID

	// GET /v2/infraenvs
	listInfraEnvParams := &installer.ListInfraEnvsParams{ClusterID: clusterZeroID}
	infraenvResult, err := restClient.Installer.ListInfraEnvs(ctx, listInfraEnvParams)
	if err != nil {
		return nil, err
	}
	clusterZeroInfraEnv := infraenvResult.Payload[0].ID

	return &nodeZeroClient{restClient, nodeZeroIP, clusterZeroID, clusterZeroInfraEnv}, nil
}

func isAgentAPILive(ctx context.Context, zero *nodeZeroClient) (bool, error) {

	// GET /v2/openshift-versions
	listOpenshiftVersionsParams := versions.NewV2ListSupportedOpenshiftVersionsParams()
	_, err := zero.restClient.Versions.ListSupportedOpenshiftVersions(ctx, (*versions.ListSupportedOpenshiftVersionsParams)(listOpenshiftVersionsParams))
	if err != nil {
		return false, err
	}

	return true, nil
}

func getClusterZero(ctx context.Context, zero *nodeZeroClient) (*models.Cluster, error) {

	// GET /v2/clusters/${cluster_zero_id}
	getClusterParams := &installer.V2GetClusterParams{ClusterID: *zero.clusterZeroID}
	result, err := zero.restClient.Installer.V2GetCluster(ctx, getClusterParams)
	if err != nil {
		return nil, err
	}
	clusterZero := result.Payload

	return clusterZero, nil
}
