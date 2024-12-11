package joiner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/coreos/stream-metadata-go/stream"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/models"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	machineconfigclient "github.com/openshift/client-go/machineconfiguration/clientset/versioned"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	workflowreport "github.com/openshift/installer/pkg/asset/agent/workflow/report"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// ClusterInfo it's an asset used to retrieve config info
// from an already existing cluster. A number of different resources
// are inspected to extract the required configuration.
type ClusterInfo struct {
	Client                       kubernetes.Interface
	OpenshiftClient              configclient.Interface
	OpenshiftMachineConfigClient machineconfigclient.Interface
	addNodesConfig               AddNodesConfig

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
	IgnitionEndpointWorker        *models.IgnitionEndpoint
	FIPS                          bool
	Nodes                         *corev1.NodeList
	ChronyConf                    *igntypes.File
	BootArtifactsBaseURL          string
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
func (ci *ClusterInfo) Generate(ctx context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	dependencies.Get(agentWorkflow, &ci.addNodesConfig)

	if agentWorkflow.Workflow != workflow.AgentWorkflowTypeAddNodes {
		return nil
	}

	if err := workflowreport.GetReport(ctx).Stage(workflow.StageClusterInspection); err != nil {
		return err
	}

	err := ci.initClients()
	if err != nil {
		return err
	}

	for _, f := range []func() error{
		ci.retrieveClusterData,
		ci.retrieveProxy,
		ci.retrievePullSecret,
		ci.retrieveUserTrustBundle,
		ci.retrieveArchitecture,
		ci.retrieveImageDigestMirrorSets,
		ci.retrieveImageContentPolicies,
		ci.retrieveOsImage,
		ci.retrieveIgnitionEndpointWorker,
		ci.retrievePlatformType,
		ci.retrieveAPIDNSName,
		ci.retrieveClusterName,
		ci.retrieveSSHKey,
		ci.retrieveFIPS,
		ci.retrieveWorkerChronyConfig,
		ci.retrieveBootArtifactsBaseURL,
		ci.retrieveNamespace,
	} {
		if err := f(); err != nil {
			return err
		}
	}

	if err = ci.validate().ToAggregate(); err != nil {
		return err
	}
	return ci.reportResult(ctx)
}

func (ci *ClusterInfo) initClients() error {
	if ci.Client != nil && ci.OpenshiftClient != nil {
		return nil
	}

	var err error
	var config *rest.Config
	kubeconfig := ci.addNodesConfig.Params.Kubeconfig
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

	openshiftMachineConfigClient, err := machineconfigclient.NewForConfig(config)
	if err != nil {
		return err
	}
	ci.OpenshiftMachineConfigClient = openshiftMachineConfigClient

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
	nodes, err := ci.Client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	ci.Nodes = nodes

	if ci.addNodesConfig.Config.CPUArchitecture != "" {
		logrus.Infof("CPU architecture set to: %v", ci.addNodesConfig.Config.CPUArchitecture)
		ci.Architecture = ci.addNodesConfig.Config.CPUArchitecture
		return nil
	}

	for _, n := range ci.Nodes.Items {
		if _, found := n.GetLabels()["node-role.kubernetes.io/master"]; found {
			ci.Architecture = n.Status.NodeInfo.Architecture
			return nil
		}
	}

	return fmt.Errorf("unable to determine target cluster architecture")
}

func (ci *ClusterInfo) retrieveFIPS() error {
	workerMachineConfig, err := ci.OpenshiftMachineConfigClient.MachineconfigurationV1().MachineConfigs().Get(context.Background(), "99-worker-fips", metav1.GetOptions{})
	if err != nil {
		// Older oc client may not have sufficient permissions,
		// falling back to previous implementation.
		if errors.IsForbidden(err) {
			installConfig, err := ci.retrieveInstallConfig()
			if err != nil {
				if errors.IsNotFound(err) {
					return nil
				}
				return err
			}
			ci.FIPS = installConfig.FIPS
			return nil
		}
		// This resource is not created when FIPS is not enabled
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	ci.FIPS = workerMachineConfig.Spec.FIPS
	return nil
}

func (ci *ClusterInfo) retrieveSSHKey() error {
	if ci.addNodesConfig.Config.SSHKey != "" {
		ci.SSHKey = ci.addNodesConfig.Config.SSHKey
		return nil
	}

	workerMachineConfig, err := ci.OpenshiftMachineConfigClient.MachineconfigurationV1().MachineConfigs().Get(context.Background(), "99-worker-ssh", metav1.GetOptions{})
	if err != nil {
		// Older oc client may not have sufficient permissions,
		// falling back to previous implementation.
		if errors.IsForbidden(err) {
			installConfig, err := ci.retrieveInstallConfig()
			if err != nil {
				if errors.IsNotFound(err) {
					return nil
				}
				return err
			}
			ci.SSHKey = installConfig.SSHKey
			return nil
		}
		return err
	}
	var ign igntypes.Config
	if err := yaml.Unmarshal(workerMachineConfig.Spec.Config.Raw, &ign); err != nil {
		return err
	}
	if len(ign.Passwd.Users) == 0 {
		return fmt.Errorf("cannot retrieve SSH key from machineconfig/99-worker-ssh: no user found")
	}
	if len(ign.Passwd.Users[0].SSHAuthorizedKeys) == 0 {
		return fmt.Errorf("cannot retrieve SSH key from machineconfig/99-worker-ssh: no SSH key found for user %s", ign.Passwd.Users[0].Name)
	}
	ci.SSHKey = string(ign.Passwd.Users[0].SSHAuthorizedKeys[0])
	return nil
}

func (ci *ClusterInfo) retrieveWorkerChronyConfig() error {
	workerMachineConfig, err := ci.OpenshiftMachineConfigClient.MachineconfigurationV1().MachineConfigs().Get(context.Background(), "50-workers-chrony-configuration", metav1.GetOptions{})
	if err != nil {
		// Older oc client may not have sufficient permissions,
		// falling back to previous implementation.
		if errors.IsForbidden(err) {
			return nil
		}
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	var ign igntypes.Config
	if err := yaml.Unmarshal(workerMachineConfig.Spec.Config.Raw, &ign); err != nil {
		return err
	}
	for _, f := range ign.Storage.Files {
		if f.Path != "/etc/chrony.conf" {
			continue
		}
		chronyConf := f
		ci.ChronyConf = &chronyConf
		break
	}
	return nil
}

func (ci *ClusterInfo) retrieveBootArtifactsBaseURL() error {
	if ci.addNodesConfig.Config.BootArtifactsBaseURL != "" {
		ci.BootArtifactsBaseURL = ci.addNodesConfig.Config.BootArtifactsBaseURL
	}
	return nil
}

func (ci *ClusterInfo) retrieveInstallConfig() (*types.InstallConfig, error) {
	clusterConfig, err := ci.Client.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "cluster-config-v1", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	data, ok := clusterConfig.Data["install-config"]
	if !ok {
		return nil, fmt.Errorf("cannot find install-config data")
	}

	installConfig := types.InstallConfig{}
	if err = yaml.Unmarshal([]byte(data), &installConfig); err != nil {
		return nil, err
	}
	return &installConfig, nil
}

func (ci *ClusterInfo) retrieveImageDigestMirrorSets() error {
	imageDigestMirrorSets, err := ci.OpenshiftClient.ConfigV1().ImageDigestMirrorSets().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		// Older oc client may not have sufficient permissions,
		// falling back to previous implementation.
		if errors.IsForbidden(err) {
			installConfig, err := ci.retrieveInstallConfig()
			if err != nil {
				if errors.IsNotFound(err) {
					return nil
				}
				return err
			}
			ci.ImageDigestSources = installConfig.ImageDigestSources
			return nil
		}
		if !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	for _, idms := range imageDigestMirrorSets.Items {
		for _, digestMirror := range idms.Spec.ImageDigestMirrors {
			digestSource := types.ImageDigestSource{
				Source: digestMirror.Source,
			}
			for _, m := range digestMirror.Mirrors {
				digestSource.Mirrors = append(digestSource.Mirrors, string(m))
			}
			ci.ImageDigestSources = append(ci.ImageDigestSources, digestSource)
		}
	}

	return nil
}

func (ci *ClusterInfo) retrieveImageContentPolicies() error {
	imageContentPolicies, err := ci.OpenshiftClient.ConfigV1().ImageContentPolicies().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		// Older oc client may not have sufficient permissions,
		// falling back to previous implementation.
		if errors.IsForbidden(err) {
			installConfig, err := ci.retrieveInstallConfig()
			if err != nil {
				if errors.IsNotFound(err) {
					return nil
				}
				return err
			}
			ci.DeprecatedImageContentSources = installConfig.DeprecatedImageContentSources
			return nil
		}
		if !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	for _, icp := range imageContentPolicies.Items {
		for _, digestMirror := range icp.Spec.RepositoryDigestMirrors {
			digestSource := types.ImageContentSource{
				Source: digestMirror.Source,
			}
			for _, m := range digestMirror.Mirrors {
				digestSource.Mirrors = append(digestSource.Mirrors, string(m))
			}
			ci.DeprecatedImageContentSources = append(ci.DeprecatedImageContentSources, digestSource)
		}
	}

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

// This method retrieves, if present, the secured ignition endpoint - along with its ca certificate.
// These information will be used to configure subsequently the imported Assisted Service cluster,
// so that the secure port (22623) could be used by the nodes to fetch the worker ignition.
func (ci *ClusterInfo) retrieveIgnitionEndpointWorker() error {
	workerUserDataManaged, err := ci.Client.CoreV1().Secrets("openshift-machine-api").Get(context.Background(), "worker-user-data-managed", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	userData := workerUserDataManaged.Data["userData"]

	config := &igntypes.Config{}
	err = json.Unmarshal(userData, config)
	if err != nil {
		return err
	}

	// Check if there is at least a CA certificate in the ignition
	if len(config.Ignition.Security.TLS.CertificateAuthorities) == 0 {
		return nil
	}

	// Get the first source and ca certificate (and strip the base64 data header)
	ignEndpointURL := config.Ignition.Config.Merge[0].Source
	caCertSource := *config.Ignition.Security.TLS.CertificateAuthorities[0].Source

	hdrIndex := strings.Index(caCertSource, ",")
	if hdrIndex == -1 {
		return fmt.Errorf("error while parsing ignition endpoints ca certificate, invalid data url format")
	}
	caCert := caCertSource[hdrIndex+1:]

	ci.IgnitionEndpointWorker = &models.IgnitionEndpoint{
		URL:           ignEndpointURL,
		CaCertificate: &caCert,
	}

	return nil
}

func (ci *ClusterInfo) getInfrastructure() (*configv1.Infrastructure, error) {
	return ci.OpenshiftClient.ConfigV1().Infrastructures().Get(context.Background(), "cluster", metav1.GetOptions{})
}

func (ci *ClusterInfo) retrievePlatformType() error {
	infra, err := ci.getInfrastructure()
	if err != nil {
		return err
	}
	platform, err := ci.toTypesPlatform(infra.Spec.PlatformSpec.Type)
	if err != nil {
		return err
	}

	ci.PlatformType = agent.HivePlatformType(platform)
	return nil
}

func (ci *ClusterInfo) retrieveAPIDNSName() error {
	infra, err := ci.getInfrastructure()
	if err != nil {
		return err
	}

	apiURL, err := url.Parse(infra.Status.APIServerURL)
	if err != nil {
		return err
	}
	ci.APIDNSName = apiURL.Hostname()
	return nil
}

func (ci *ClusterInfo) retrieveClusterName() error {
	if ci.APIDNSName == "" {
		return fmt.Errorf("cannot get cluster name: API DNS name is empty")
	}

	re := regexp.MustCompile(`^api\.([^.]+)\..*`)
	match := re.FindStringSubmatch(ci.APIDNSName)
	if len(match) > 1 {
		ci.ClusterName = match[1]
		return nil
	}

	return fmt.Errorf("cannot get cluster name from API DNS name: %s", ci.APIDNSName)
}

func (ci *ClusterInfo) retrieveNamespace() error {
	ci.Namespace = "cluster0"
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

func (ci *ClusterInfo) validate() field.ErrorList {
	var allErrs field.ErrorList

	if err := ci.validateSupportedPlatforms(); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (ci *ClusterInfo) validateSupportedPlatforms() field.ErrorList {
	var allErrs field.ErrorList

	infra, err := ci.getInfrastructure()
	if err != nil {
		return append(allErrs, field.InternalError(nil, err))
	}
	platformType, err := ci.toTypesPlatform(infra.Spec.PlatformSpec.Type)
	if err != nil {
		return append(allErrs, field.InternalError(nil, err))
	}
	return agent.ValidateSupportedPlatforms(platformType, ci.Architecture)
}

func (ci *ClusterInfo) toTypesPlatform(platformType configv1.PlatformType) (types.Platform, error) {
	platform := types.Platform{}

	switch platformType {
	case configv1.AWSPlatformType:
		platform.AWS = &aws.Platform{}
	case configv1.AzurePlatformType:
		platform.Azure = &azure.Platform{}
	case configv1.BareMetalPlatformType:
		platform.BareMetal = &baremetal.Platform{}
	case configv1.GCPPlatformType:
		platform.GCP = &gcp.Platform{}
	case configv1.OpenStackPlatformType:
		platform.OpenStack = &openstack.Platform{}
	case configv1.NonePlatformType:
		platform.None = &none.Platform{}
	case configv1.VSpherePlatformType:
		platform.VSphere = &vsphere.Platform{}
	case configv1.OvirtPlatformType:
		platform.Ovirt = &ovirt.Platform{}
	case configv1.IBMCloudPlatformType:
		platform.IBMCloud = &ibmcloud.Platform{}
	case configv1.PowerVSPlatformType:
		platform.PowerVS = &powervs.Platform{}
	case configv1.NutanixPlatformType:
		platform.Nutanix = &nutanix.Platform{}
	case configv1.ExternalPlatformType:
		platform.External = &external.Platform{}
	default:
		return platform, fmt.Errorf("unable to convert platform type %v", platformType)
	}

	return platform, nil
}

func (ci *ClusterInfo) reportResult(ctx context.Context) error {
	results := map[string]string{
		"Version":      ci.Version,
		"ReleaseImage": ci.ReleaseImage,
		"OSImage":      ci.OSImageLocation,
		"PlatformType": string(ci.PlatformType),
		"Architecture": ci.Architecture,
	}

	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	return workflowreport.GetReport(ctx).StageResult(workflow.StageClusterInspection, string(data))
}
