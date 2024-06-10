package joiner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coreos/stream-metadata-go/arch"
	"github.com/coreos/stream-metadata-go/stream"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/types"
)

// ClusterInfo it's an asset used to retrieve config info
// from an already existing cluster. A number of different resources
// are inspected to extract the required configuration.
type ClusterInfo struct {
	Client          kubernetes.Interface
	OpenshiftClient configclient.Interface

	ClusterID                     string
	ClusterName                   string
	Version                       string
	ReleaseImage                  string
	APIDNSName                    string
	PullSecret                    string
	Namespace                     string
	UserCaBundle                  string
	Proxy                         *types.Proxy
	Architecture                  string
	ImageDigestSources            []types.ImageDigestSource
	DeprecatedImageContentSources []types.ImageContentSource
	PlatformType                  hiveext.PlatformType
	SSHKey                        string
	OSImage                       *stream.Stream
	OSImageLocation               string
}

var _ asset.WritableAsset = (*ClusterInfo)(nil)

// Name returns the human-friendly name of the asset.
func (ci *ClusterInfo) Name() string {
	return "Agent Installer ClusterInfo"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*ClusterInfo) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&AddNodesConfig{},
	}
}

// Generate generates the ClusterInfo.
func (ci *ClusterInfo) Generate(_ context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	addNodesConfig := &AddNodesConfig{}
	dependencies.Get(agentWorkflow, addNodesConfig)

	if agentWorkflow.Workflow != workflow.AgentWorkflowTypeAddNodes {
		return nil
	}

	err := ci.initClients(addNodesConfig.Params.Kubeconfig)
	if err != nil {
		return err
	}
	err = ci.retrieveClusterData()
	if err != nil {
		return err
	}
	err = ci.retrieveProxy()
	if err != nil {
		return err
	}

	err = ci.retrievePullSecret()
	if err != nil {
		return err
	}
	err = ci.retrieveUserTrustBundle()
	if err != nil {
		return err
	}
	err = ci.retrieveArchitecture()
	if err != nil {
		return err
	}
	err = ci.retrieveInstallConfigData()
	if err != nil {
		return err
	}
	err = ci.retrieveOsImage()
	if err != nil {
		return err
	}

	ci.Namespace = "cluster0"

	return nil
}

func (ci *ClusterInfo) initClients(kubeconfig string) error {
	if ci.Client != nil && ci.OpenshiftClient != nil {
		return nil
	}

	var err error
	var config *rest.Config
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return err
	}

	openshiftClient, err := configclient.NewForConfig(config)
	if err != nil {
		return err
	}
	ci.OpenshiftClient = openshiftClient

	k8sclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	ci.Client = k8sclientset

	return err
}

func (ci *ClusterInfo) retrieveClusterData() error {
	cv, err := ci.OpenshiftClient.ConfigV1().ClusterVersions().Get(context.Background(), "version", metav1.GetOptions{})
	if err != nil {
		return err
	}
	ci.ClusterID = string(cv.Spec.ClusterID)
	ci.ReleaseImage = cv.Status.History[0].Image
	ci.Version = cv.Status.History[0].Version

	return nil
}

func (ci *ClusterInfo) retrieveProxy() error {
	proxy, err := ci.OpenshiftClient.ConfigV1().Proxies().Get(context.Background(), "cluster", metav1.GetOptions{})
	if err != nil {
		return err
	}
	ci.Proxy = &types.Proxy{
		HTTPProxy:  proxy.Spec.HTTPProxy,
		HTTPSProxy: proxy.Spec.HTTPSProxy,
		NoProxy:    proxy.Spec.NoProxy,
	}

	return nil
}

func (ci *ClusterInfo) retrievePullSecret() error {
	pullSecret, err := ci.Client.CoreV1().Secrets("openshift-config").Get(context.Background(), "pull-secret", metav1.GetOptions{})
	if err != nil {
		return err
	}
	ci.PullSecret = string(pullSecret.Data[".dockerconfigjson"])

	return nil
}

func (ci *ClusterInfo) retrieveUserTrustBundle() error {
	userCaBundle, err := ci.Client.CoreV1().ConfigMaps("openshift-config").Get(context.Background(), "user-ca-bundle", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	ci.UserCaBundle = userCaBundle.Data["ca-bundle.crt"]

	return nil
}

func (ci *ClusterInfo) retrieveArchitecture() error {
	nodes, err := ci.Client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/master",
	})
	if err != nil {
		return err
	}
	ci.Architecture = nodes.Items[0].Status.NodeInfo.Architecture

	return nil
}

func (ci *ClusterInfo) retrieveInstallConfigData() error {
	clusterConfig, err := ci.Client.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "cluster-config-v1", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	data, ok := clusterConfig.Data["install-config"]
	if !ok {
		return fmt.Errorf("cannot find install-config data")
	}

	installConfig := types.InstallConfig{}
	if err = yaml.Unmarshal([]byte(data), &installConfig); err != nil {
		return err
	}

	ci.ImageDigestSources = installConfig.ImageDigestSources
	ci.DeprecatedImageContentSources = installConfig.DeprecatedImageContentSources
	ci.PlatformType = agent.HivePlatformType(installConfig.Platform)
	ci.SSHKey = installConfig.SSHKey
	ci.ClusterName = installConfig.ObjectMeta.Name
	ci.APIDNSName = fmt.Sprintf("api.%s.%s", ci.ClusterName, installConfig.BaseDomain)

	return nil
}

func (ci *ClusterInfo) retrieveOsImage() error {
	clusterConfig, err := ci.Client.CoreV1().ConfigMaps("openshift-machine-config-operator").Get(context.Background(), "coreos-bootimages", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	data, ok := clusterConfig.Data["stream"]
	if !ok {
		return fmt.Errorf("cannot find stream data")
	}

	var st stream.Stream
	if err := json.Unmarshal([]byte(data), &st); err != nil {
		return fmt.Errorf("failed to parse CoreOS stream metadata: %w", err)
	}
	ci.OSImage = &st

	clusterArch := arch.RpmArch(ci.Architecture)
	streamArch, err := st.GetArchitecture(clusterArch)
	if err != nil {
		return err
	}
	metal, ok := streamArch.Artifacts["metal"]
	if !ok {
		return fmt.Errorf("stream data not found for 'metal' artifact")
	}
	format, ok := metal.Formats["iso"]
	if !ok {
		return fmt.Errorf("no ISO found to download for %s", clusterArch)
	}
	ci.OSImageLocation = format.Disk.Location

	return nil
}

// Files returns the files generated by the asset.
func (*ClusterInfo) Files() []*asset.File {
	return []*asset.File{}
}

// Load returns agent config asset from the disk.
func (*ClusterInfo) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
