package agent

import (
	"context"
	"net"
	"net/url"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"

	"github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/client/versions"

	"github.com/openshift/installer/pkg/asset/agent/manifests"
	assetstore "github.com/openshift/installer/pkg/asset/store"
)

type nodeZeroRestClient struct {
	Client     *client.AssistedInstall
	ctx        context.Context
	config     client.Config
	nodeZeroIP string
}

func NewNodeZeroRestClient(ctx context.Context, assetDir string) (*nodeZeroRestClient, error) {
	zeroRestClient := &nodeZeroRestClient{}

	assetStore, err := assetstore.NewStore(assetDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create asset store")
	}
	nmState := &manifests.NMStateConfig{}
	if err := assetStore.Fetch(nmState); err != nil {
		return nil, errors.Wrapf(err, "failed to fetch %s", nmState.Name())
	}

	nodeZeroIP, err := manifests.GetNodeZeroIP(nmState.Config)
	if err != nil {
		return nil, err
	}
	config := client.Config{}
	config.URL = &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(nodeZeroIP, "8090"),
		Path:   client.DefaultBasePath,
	}
	client := client.New(config)

	zeroRestClient.Client = client
	zeroRestClient.ctx = ctx
	zeroRestClient.config = config
	zeroRestClient.nodeZeroIP = nodeZeroIP

	return zeroRestClient, nil
}

func (zerorest *nodeZeroRestClient) IsAgentAPILive() (bool, error) {
	// GET /v2/openshift-versions
	listOpenshiftVersionsParams := versions.NewV2ListSupportedOpenshiftVersionsParams()
	_, err := zerorest.Client.Versions.ListSupportedOpenshiftVersions(zerorest.ctx, (*versions.ListSupportedOpenshiftVersionsParams)(listOpenshiftVersionsParams))

	if err != nil {
		return false, err
	}
	return true, nil
}

func (zerorest *nodeZeroRestClient) GetAgentAPIServiceBaseURL() *url.URL {
	return zerorest.config.URL
}

func (zerorest *nodeZeroRestClient) getAgentClusterClusterID() (*strfmt.UUID, error) {
	// GET /v2/clusters and return first result
	listClusterParams := installer.NewV2ListClustersParams()
	clusterResult, err := zerorest.Client.Installer.V2ListClusters(zerorest.ctx, listClusterParams)
	if err != nil {
		return nil, err
	}
	agentClusterID := clusterResult.Payload[0].ID
	return agentClusterID, nil
}

func (zerorest *nodeZeroRestClient) getAgentClusterInfraEnvID() (*strfmt.UUID, error) {
	// GET /v2/infraenvs and return first result
	listInfraEnvParams := installer.NewListInfraEnvsParams()
	infraenvResult, err := zerorest.Client.Installer.ListInfraEnvs(zerorest.ctx, listInfraEnvParams)
	if err != nil {
		return nil, err
	}
	agentClusterInfraEnvID := infraenvResult.Payload[0].ID
	return agentClusterInfraEnvID, nil
}
