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

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
)

type nodeZeroClient struct {
	ctx        context.Context
	restClient *client.AssistedInstall
	restConfig client.Config
	nodeZeroIP string
}

func NewNodeZeroClient() (*nodeZeroClient, error) {
	zero := &nodeZeroClient{}
	agentManifests := &manifests.AgentManifests{}
	dependencies := &asset.Parents{}
	dependencies.Get(agentManifests)

	nodeZeroIP, err := manifests.GetNodeZeroIP(agentManifests.NMStateConfigs)
	if err != nil {
		return nil, err
	}

	restConfig := client.Config{}
	*restConfig.URL = url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(nodeZeroIP, "8090"),
		Path:   client.DefaultBasePath,
	}
	restClient := client.New(restConfig)

	zero.ctx = context.Background()
	zero.restClient = restClient
	zero.restConfig = restConfig
	zero.nodeZeroIP = nodeZeroIP

	return zero, nil
}

func (zero *nodeZeroClient) isAgentAPILive() (bool, error) {
	// GET /v2/openshift-versions
	listOpenshiftVersionsParams := versions.NewV2ListSupportedOpenshiftVersionsParams()
	_, err := zero.restClient.Versions.ListSupportedOpenshiftVersions(zero.ctx, (*versions.ListSupportedOpenshiftVersionsParams)(listOpenshiftVersionsParams))

	if err != nil {
		return false, err
	}
	return true, nil
}

func (zero *nodeZeroClient) getAgentAPIServiceBaseURL() *url.URL {
	return zero.restConfig.URL
}

func (zero *nodeZeroClient) getClusterZeroClusterID() (*strfmt.UUID, error) {
	// GET /v2/clusters
	listClusterParams := installer.NewV2ListClustersParams()
	clusterResult, err := zero.restClient.Installer.V2ListClusters(zero.ctx, listClusterParams)
	if err != nil {
		return nil, err
	}
	clusterZeroID := clusterResult.Payload[0].ID
	return clusterZeroID, nil
}

func (zero *nodeZeroClient) getClusterZeroInfraEnvID() (*strfmt.UUID, error) {
	// GET /v2/infraenvs
	listInfraEnvParams := installer.NewListInfraEnvsParams()
	infraenvResult, err := zero.restClient.Installer.ListInfraEnvs(zero.ctx, listInfraEnvParams)
	if err != nil {
		return nil, err
	}
	clusterZeroInfraEnvID := infraenvResult.Payload[0].ID
	return clusterZeroInfraEnvID, nil
}

type clusterZero struct {
	clusterZeroID         *strfmt.UUID
	clusterZeroInfraEnvID *strfmt.UUID
	zeroClient            *nodeZeroClient
}

func NewClusterZero(zeroClient *nodeZeroClient) (*clusterZero, error) {
	czero := &clusterZero{}
	clusterZeroID, err := zeroClient.getClusterZeroClusterID()
	if err != nil {
		return nil, err
	}
	clusterZeroInfraEnvID, err := zeroClient.getClusterZeroInfraEnvID()
	if err != nil {
		return nil, err
	}

	czero.clusterZeroID = clusterZeroID
	czero.clusterZeroInfraEnvID = clusterZeroInfraEnvID
	czero.zeroClient = zeroClient
	return czero, nil
}

func (czero *clusterZero) get() (*models.Cluster, error) {
	// GET /v2/clusters/{cluster_zero_id}
	getClusterParams := &installer.V2GetClusterParams{ClusterID: *czero.clusterZeroID}
	result, err := czero.zeroClient.restClient.Installer.V2GetCluster(czero.zeroClient.ctx, getClusterParams)
	if err != nil {
		return nil, err
	}
	clusterZero := result.Payload
	return clusterZero, nil
}

func (czero *clusterZero) parseValidationInfo(*models.Cluster) (bool, error) {

	return false, nil
}

func (czero *clusterZero) isKubeAPILive() (bool, error) {

	return false, nil
}

func (czero *clusterZero) doesKubeConfigExist() (bool, error) {

	return false, nil
}

func (czero *clusterZero) printInstallStatus() error {

	return nil
}
