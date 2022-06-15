package agent

import (
	"context"
	"net"
	"net/url"
	"path/filepath"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift/assisted-service/client"
	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/client/versions"
	"github.com/openshift/assisted-service/models"

	"github.com/openshift/installer/pkg/asset/agent/manifests"
	assetstore "github.com/openshift/installer/pkg/asset/store"
)

type nodeZeroClient struct {
	ctx        context.Context
	restClient *client.AssistedInstall
	restConfig client.Config
	nodeZeroIP string
}

func NewNodeZeroClient(directory string) (*nodeZeroClient, error) {
	zero := &nodeZeroClient{}

	assetStore, err := assetstore.NewStore(directory)
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

func (czero *clusterZero) isKubeAPILive(directory string) (bool, error) {

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(directory, "auth", "kubeconfig"))
	if err != nil {
		return false, errors.Wrap(err, "loading kubeconfig")
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, errors.Wrap(err, "creating a Kubernetes client")
	}

	discovery := kubeClient.Discovery()
	version, err := discovery.ServerVersion()
	if err != nil {
		return false, err
	}
	logrus.Infof("Cluster API is up and running %s", version)
	return true, nil
}

func (czero *clusterZero) doesKubeConfigExist(directory string) (bool, error) {

	kubeconfig := filepath.Join(directory, "auth", "kubeconfig")
	_, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return false, errors.Wrap(err, "loading kubeconfig")
	}
	return true, nil
}

func (czero *clusterZero) printInstallStatus() error {

	return nil
}
