package image

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/coreos/stream-metadata-go/stream"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	agentcommon "github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/common"
	"github.com/openshift/installer/pkg/asset/agent/gencrypto"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/agent"
	"github.com/openshift/installer/pkg/version"
)

const addNodesEnvPath = "/etc/assisted/add-nodes.env"
const rendezvousHostEnvPath = "/etc/assisted/rendezvous-host.env"
const manifestPath = "/etc/assisted/manifests"
const hostnamesPath = "/etc/assisted/hostnames"
const nmConnectionsPath = "/etc/assisted/network"
const extraManifestPath = "/etc/assisted/extra-manifests"
const registriesConfPath = "/etc/containers/registries.conf"
const registryCABundlePath = "/etc/pki/ca-trust/source/anchors/domain.crt"

// Ignition is an asset that generates the agent installer ignition file.
type Ignition struct {
	Config       *igntypes.Config
	CPUArch      string
	RendezvousIP string
}

// agentTemplateData is the data used to replace values in agent template
// files.
type agentTemplateData struct {
	ServiceProtocol           string
	PullSecret                string
	ControlPlaneAgents        int
	WorkerAgents              int
	ReleaseImages             string
	ReleaseImage              string
	ReleaseImageMirror        string
	HaveMirrorConfig          bool
	PublicContainerRegistries string
	InfraEnvID                string
	ClusterName               string
	OSImage                   *models.OsImage
	Proxy                     *v1beta1.Proxy
	ConfigImageFiles          string
	ImageTypeISO              string
	PublicKeyPEM              string
	PrivateKeyPEM             string
	Token                     string
	CaBundleMount             string
}

// Name returns the human-friendly name of the asset.
func (a *Ignition) Name() string {
	return "Agent Installer Ignition"
}

// Dependencies returns the assets on which the Ignition asset depends.
func (a *Ignition) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&joiner.ClusterInfo{},
		&joiner.AddNodesConfig{},
		&manifests.AgentManifests{},
		&manifests.ExtraManifests{},
		&tls.KubeAPIServerLBSignerCertKey{},
		&tls.KubeAPIServerLocalhostSignerCertKey{},
		&tls.KubeAPIServerServiceNetworkSignerCertKey{},
		&tls.AdminKubeConfigSignerCertKey{},
		&password.KubeadminPassword{},
		&agentconfig.AgentConfig{},
		&agentconfig.AgentHosts{},
		&mirror.RegistriesConf{},
		&mirror.CaBundle{},
		&gencrypto.AuthConfig{},
		&common.InfraEnvID{},
	}
}

// Generate generates the agent installer ignition.
func (a *Ignition) Generate(dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	clusterInfo := &joiner.ClusterInfo{}
	addNodesConfig := &joiner.AddNodesConfig{}
	agentManifests := &manifests.AgentManifests{}
	agentConfigAsset := &agentconfig.AgentConfig{}
	agentHostsAsset := &agentconfig.AgentHosts{}
	extraManifests := &manifests.ExtraManifests{}
	keyPairAsset := &gencrypto.AuthConfig{}
	infraEnvAsset := &common.InfraEnvID{}
	dependencies.Get(agentManifests, agentConfigAsset, agentHostsAsset, extraManifests, keyPairAsset, agentWorkflow, clusterInfo, addNodesConfig, infraEnvAsset)

	pwd := &password.KubeadminPassword{}
	dependencies.Get(pwd)
	pwdHash := string(pwd.PasswordHash)

	infraEnv := agentManifests.InfraEnv

	config := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Passwd: igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name: "core",
					SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
						igntypes.SSHAuthorizedKey(infraEnv.Spec.SSHAuthorizedKey),
					},
					PasswordHash: &pwdHash,
				},
			},
		},
	}

	clusterName := ""
	imageTypeISO := "full-iso"
	numMasters := 0
	numWorkers := 0
	enabledServices := getDefaultEnabledServices()
	openshiftVersion := ""
	var err error
	var streamGetter CoreOSBuildFetcher

	switch agentWorkflow.Workflow {
	case workflow.AgentWorkflowTypeInstall:
		// Set rendezvous IP.
		nodeZeroIP, err := RetrieveRendezvousIP(agentConfigAsset.Config, agentHostsAsset.Hosts, agentManifests.NMStateConfigs)
		if err != nil {
			return err
		}
		a.RendezvousIP = nodeZeroIP
		logrus.Infof("The rendezvous host IP (node0 IP) is %s", a.RendezvousIP)
		// Define cluster name and image type.
		clusterName = fmt.Sprintf("%s.%s", agentManifests.ClusterDeployment.Spec.ClusterName, agentManifests.ClusterDeployment.Spec.BaseDomain)
		if agentManifests.AgentClusterInstall.Spec.PlatformType == hiveext.ExternalPlatformType {
			imageTypeISO = "minimal-iso"
		}
		// Fetch the required number of master and worker nodes.
		numMasters = agentManifests.AgentClusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents
		numWorkers = agentManifests.AgentClusterInstall.Spec.ProvisionRequirements.WorkerAgents
		// Enable specific install services
		enabledServices = append(enabledServices, "start-cluster-installation.service")
		// Version is retrieved from the embedded data
		openshiftVersion, err = version.Version()
		if err != nil {
			return err
		}
		streamGetter = DefaultCoreOSStreamGetter

	case workflow.AgentWorkflowTypeAddNodes:
		// In the add-nodes workflow, every node will act independently from the others.
		a.RendezvousIP = "127.0.0.1"
		// Reuse the existing cluster name.
		clusterName = clusterInfo.ClusterName
		// Fetch the required number of master and worker nodes. Currently only adding workers
		// is supported, so forcing the expected number of masters to zero, and assuming implcitly
		// that all the hosts defined are workers.
		numMasters = 0
		numWorkers = len(addNodesConfig.Config.Hosts)
		// Enable add-nodes specific services
		enabledServices = append(enabledServices, "agent-add-node.service")
		// Generate add-nodes.env file
		addNodesEnvFile := ignition.FileFromString(addNodesEnvPath, "root", 0644, getAddNodesEnv(*clusterInfo))
		config.Storage.Files = append(config.Storage.Files, addNodesEnvFile)
		// Version matches the source cluster one
		openshiftVersion = clusterInfo.Version
		streamGetter = func(ctx context.Context) (*stream.Stream, error) {
			return clusterInfo.OSImage, nil
		}

	default:
		return fmt.Errorf("AgentWorkflowType value not supported: %s", agentWorkflow.Workflow)
	}

	// Default to x86_64
	archName := arch.RpmArch(types.ArchitectureAMD64)
	if infraEnv.Spec.CpuArchitecture != "" {
		archName = infraEnv.Spec.CpuArchitecture
	}
	// Examine the release payload to see if its multi
	releaseArch, err := agentcommon.DetermineReleaseImageArch(agentManifests.GetPullSecretData(), agentManifests.ClusterImageSet.Spec.ReleaseImage)
	if err != nil {
		logrus.Warnf("Unable to validate the release image architecture, using infraEnv.Spec.CpuArchitecture for the release image arch")
		releaseArch = archName
	} else {
		releaseArch = arch.RpmArch(releaseArch)
		logrus.Debugf("Found Release Image Architecture: %s", releaseArch)
	}
	releaseArchs := []string{releaseArch}
	if releaseArch == "multi" {
		releaseArchs = []string{arch.RpmArch(types.ArchitectureARM64), arch.RpmArch(types.ArchitectureAMD64), arch.RpmArch(types.ArchitecturePPC64LE), arch.RpmArch(types.ArchitectureS390X)}
	}
	releaseImageList, err := releaseImageListWithVersion(agentManifests.ClusterImageSet.Spec.ReleaseImage, releaseArch, releaseArchs, openshiftVersion)
	if err != nil {
		return err
	}

	registriesConfig := &mirror.RegistriesConf{}
	registryCABundle := &mirror.CaBundle{}
	dependencies.Get(registriesConfig, registryCABundle)

	publicContainerRegistries := getPublicContainerRegistries(registriesConfig)

	releaseImageMirror := mirror.GetMirrorFromRelease(agentManifests.ClusterImageSet.Spec.ReleaseImage, registriesConfig)

	infraEnvID := infraEnvAsset.ID
	logrus.Debug("Generated random infra-env id ", infraEnvID)

	osImage, err := getOSImagesInfo(archName, openshiftVersion, streamGetter)
	if err != nil {
		return err
	}
	a.CPUArch = *osImage.CPUArchitecture

	caBundleMount := defineCABundleMount(registriesConfig, registryCABundle)

	agentTemplateData := getTemplateData(
		clusterName,
		agentManifests.GetPullSecretData(),
		releaseImageList,
		agentManifests.ClusterImageSet.Spec.ReleaseImage,
		releaseImageMirror,
		len(registriesConfig.MirrorConfig) > 0,
		publicContainerRegistries,
		numMasters, numWorkers,
		infraEnvID,
		osImage,
		infraEnv.Spec.Proxy,
		imageTypeISO,
		keyPairAsset.PrivateKey,
		keyPairAsset.PublicKey,
		keyPairAsset.Token,
		caBundleMount)

	err = bootstrap.AddStorageFiles(&config, "/", "agent/files", agentTemplateData)
	if err != nil {
		return err
	}

	rendezvousHostFile := ignition.FileFromString(rendezvousHostEnvPath,
		"root", 0644,
		getRendezvousHostEnv(agentTemplateData.ServiceProtocol, a.RendezvousIP, agentWorkflow.Workflow))
	config.Storage.Files = append(config.Storage.Files, rendezvousHostFile)

	err = addBootstrapScripts(&config, agentManifests.ClusterImageSet.Spec.ReleaseImage)
	if err != nil {
		return err
	}

	// add ZTP manifests to manifestPath
	for _, file := range agentManifests.FileList {
		manifestFile := ignition.FileFromBytes(filepath.Join(manifestPath, filepath.Base(file.Filename)),
			"root", 0600, file.Data)
		config.Storage.Files = append(config.Storage.Files, manifestFile)
	}

	// add AgentConfig if provided
	if agentConfigAsset.Config != nil {
		agentConfigFile := ignition.FileFromBytes(filepath.Join(manifestPath, filepath.Base(agentConfigAsset.File.Filename)),
			"root", 0600, agentConfigAsset.File.Data)
		config.Storage.Files = append(config.Storage.Files, agentConfigFile)
	}

	addMacAddressToHostnameMappings(&config, agentHostsAsset)

	err = addStaticNetworkConfig(&config, agentManifests.StaticNetworkConfigs)
	if err != nil {
		return err
	}

	// Enable pre-network-manager-config.service only when there are network configs defined
	if len(agentManifests.StaticNetworkConfigs) != 0 {
		enabledServices = append(enabledServices, "pre-network-manager-config.service")
	}

	err = bootstrap.AddSystemdUnits(&config, "agent/systemd/units", agentTemplateData, enabledServices)
	if err != nil {
		return err
	}

	addTLSData(&config, dependencies)

	addMirrorData(&config, registriesConfig, registryCABundle)

	err = addHostConfig(&config, agentHostsAsset)
	if err != nil {
		return err
	}

	err = addExtraManifests(&config, extraManifests)
	if err != nil {
		return err
	}

	a.Config = &config
	return nil
}

func getDefaultEnabledServices() []string {
	return []string{
		"agent-interactive-console.service",
		"agent-interactive-console-serial@.service",
		"agent-register-cluster.service",
		"agent-import-cluster.service",
		"agent-register-infraenv.service",
		"agent.service",
		"assisted-service-db.service",
		"assisted-service-pod.service",
		"assisted-service.service",
		"node-zero.service",
		"multipathd.service",
		"selinux.service",
		"install-status.service",
		"set-hostname.service",
	}
}

func addBootstrapScripts(config *igntypes.Config, releaseImage string) (err error) {
	// Set up bootstrap service recording
	if err := bootstrap.AddStorageFiles(config,
		"/usr/local/bin/bootstrap-service-record.sh",
		"bootstrap/files/usr/local/bin/bootstrap-service-record.sh",
		nil); err != nil {
		return err
	}

	// Use bootstrap script to get container images
	relImgData := struct{ ReleaseImage string }{
		ReleaseImage: releaseImage,
	}
	for _, script := range []string{"release-image.sh", "release-image-download.sh"} {
		if err := bootstrap.AddStorageFiles(config,
			"/usr/local/bin/"+script,
			"bootstrap/files/usr/local/bin/"+script+".template",
			relImgData); err != nil {
			return err
		}
	}
	return nil
}

func getTemplateData(name, pullSecret, releaseImageList, releaseImage,
	releaseImageMirror string, haveMirrorConfig bool, publicContainerRegistries string,
	numMasters, numWorkers int,
	infraEnvID string,
	osImage *models.OsImage,
	proxy *v1beta1.Proxy,
	imageTypeISO,
	privateKey, publicKey, token string,
	caBundleMount string) *agentTemplateData {
	return &agentTemplateData{
		ServiceProtocol:           "http",
		PullSecret:                pullSecret,
		ControlPlaneAgents:        numMasters,
		WorkerAgents:              numWorkers,
		ReleaseImages:             releaseImageList,
		ReleaseImage:              releaseImage,
		ReleaseImageMirror:        releaseImageMirror,
		HaveMirrorConfig:          haveMirrorConfig,
		PublicContainerRegistries: publicContainerRegistries,
		InfraEnvID:                infraEnvID,
		ClusterName:               name,
		OSImage:                   osImage,
		Proxy:                     proxy,
		ImageTypeISO:              imageTypeISO,
		PrivateKeyPEM:             privateKey,
		PublicKeyPEM:              publicKey,
		Token:                     token,
		CaBundleMount:             caBundleMount,
	}
}

func getRendezvousHostEnv(serviceProtocol, nodeZeroIP string, workflowType workflow.AgentWorkflowType) string {
	serviceBaseURL := url.URL{
		Scheme: serviceProtocol,
		Host:   net.JoinHostPort(nodeZeroIP, "8090"),
		Path:   "/",
	}
	imageServiceBaseURL := url.URL{
		Scheme: serviceProtocol,
		Host:   net.JoinHostPort(nodeZeroIP, "8888"),
		Path:   "/",
	}

	return fmt.Sprintf(`NODE_ZERO_IP=%s
SERVICE_BASE_URL=%s
IMAGE_SERVICE_BASE_URL=%s
WORKFLOW_TYPE=%s
`, nodeZeroIP, serviceBaseURL.String(), imageServiceBaseURL.String(), workflowType)
}

func getAddNodesEnv(clusterInfo joiner.ClusterInfo) string {
	return fmt.Sprintf(`CLUSTER_ID=%s
CLUSTER_NAME=%s
CLUSTER_API_VIP_DNS_NAME=%s
`, clusterInfo.ClusterID, clusterInfo.ClusterName, clusterInfo.APIDNSName)
}

func addStaticNetworkConfig(config *igntypes.Config, staticNetworkConfig []*models.HostStaticNetworkConfig) (err error) {
	if len(staticNetworkConfig) == 0 {
		return nil
	}

	// Get the static network configuration from nmstate and generate NetworkManager ignition files
	filesList, err := manifests.GetNMIgnitionFiles(staticNetworkConfig)
	if err != nil {
		return err
	}

	for i := range filesList {
		nmFilePath := path.Join(nmConnectionsPath, filesList[i].FilePath)
		nmStateIgnFile := ignition.FileFromBytes(nmFilePath, "root", 0600, []byte(filesList[i].FileContents))
		config.Storage.Files = append(config.Storage.Files, nmStateIgnFile)
	}

	nmStateScriptFilePath := "/usr/local/bin/pre-network-manager-config.sh"
	// A local version of the assisted-service internal script is currently used
	nmStateScript := ignition.FileFromBytes(nmStateScriptFilePath, "root", 0755, []byte(manifests.PreNetworkConfigScript))
	config.Storage.Files = append(config.Storage.Files, nmStateScript)

	return nil
}

func addTLSData(config *igntypes.Config, dependencies asset.Parents) {
	certKeys := []asset.Asset{
		&tls.KubeAPIServerLBSignerCertKey{},
		&tls.KubeAPIServerLocalhostSignerCertKey{},
		&tls.KubeAPIServerServiceNetworkSignerCertKey{},
		&tls.AdminKubeConfigSignerCertKey{},
	}
	dependencies.Get(certKeys...)

	for _, ck := range certKeys {
		for _, d := range ck.(asset.WritableAsset).Files() {
			f := ignition.FileFromBytes(path.Join("/opt/agent", d.Filename), "root", 0600, d.Data)
			config.Storage.Files = append(config.Storage.Files, f)
		}
	}

	pwd := &password.KubeadminPassword{}
	dependencies.Get(pwd)
	config.Storage.Files = append(config.Storage.Files,
		ignition.FileFromBytes("/opt/agent/tls/kubeadmin-password.hash", "root", 0600, pwd.PasswordHash))
}

func addMirrorData(config *igntypes.Config, registriesConfig *mirror.RegistriesConf, registryCABundle *mirror.CaBundle) {

	// This is required for assisted-service to build the ICSP for openshift-install
	if registriesConfig.File != nil {
		registriesFile := ignition.FileFromBytes(registriesConfPath,
			"root", 0644, registriesConfig.File.Data)
		config.Storage.Files = append(config.Storage.Files, registriesFile)
	}

	// This is required for the agent to run the podman commands to the mirror
	if registryCABundle.File != nil && len(registryCABundle.File.Data) > 0 {
		caFile := ignition.FileFromBytes(registryCABundlePath,
			"root", 0600, registryCABundle.File.Data)
		config.Storage.Files = append(config.Storage.Files, caFile)
	}
}

func defineCABundleMount(registriesConfig *mirror.RegistriesConf, registryCABundle *mirror.CaBundle) string {
	// By default, the current host CA bundle is used (it will also contain eventually a user CA bundle, if
	// defined in the AdditionalTrustBundle field of install-config.yaml).
	hostSourceCABundle := "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"

	// If mirror registry is configured and the user provided a bundle, then let's mount just the user one.
	if len(registriesConfig.MirrorConfig) > 0 && registryCABundle.File != nil && len(registryCABundle.File.Data) > 0 {
		hostSourceCABundle = registryCABundlePath
	}

	return fmt.Sprintf("-v %s:/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem:z", hostSourceCABundle)
}

// Creates a file named with a host's MAC address. The desired hostname
// is the file's content. The files are read by a systemd service that
// sets the hostname using "hostnamectl set-hostname" when the ISO boots.
func addMacAddressToHostnameMappings(
	config *igntypes.Config,
	agentHostsAsset *agentconfig.AgentHosts) {
	if len(agentHostsAsset.Hosts) == 0 {
		return
	}
	for _, host := range agentHostsAsset.Hosts {
		if host.Hostname != "" {
			file := ignition.FileFromBytes(filepath.Join(hostnamesPath,
				strings.ToLower(filepath.Base(host.Interfaces[0].MacAddress))),
				"root", 0600, []byte(host.Hostname))
			config.Storage.Files = append(config.Storage.Files, file)
		}
	}
}

func addHostConfig(config *igntypes.Config, agentHosts *agentconfig.AgentHosts) error {
	confs, err := agentHosts.HostConfigFiles()
	if err != nil {
		return err
	}

	for path, content := range confs {
		hostConfigFile := ignition.FileFromBytes(filepath.Join("/etc/assisted/hostconfig", path), "root", 0644, content)
		config.Storage.Files = append(config.Storage.Files, hostConfigFile)
	}
	return nil
}

func addExtraManifests(config *igntypes.Config, extraManifests *manifests.ExtraManifests) error {

	user := "root"
	mode := 0644

	config.Storage.Directories = append(config.Storage.Directories, igntypes.Directory{
		Node: igntypes.Node{
			Path: extraManifestPath,
			User: igntypes.NodeUser{
				Name: &user,
			},
			Overwrite: util.BoolToPtr(true),
		},
		DirectoryEmbedded1: igntypes.DirectoryEmbedded1{
			Mode: &mode,
		},
	})

	for _, file := range extraManifests.FileList {

		type unstructured map[string]interface{}

		yamlList, err := manifests.GetMultipleYamls[unstructured](file.Data)
		if err != nil {
			return errors.Wrapf(err, "could not decode YAML for %s", file.Filename)
		}

		for n, manifest := range yamlList {
			m, err := yaml.Marshal(manifest)
			if err != nil {
				return err
			}

			base := filepath.Base(file.Filename)
			ext := filepath.Ext(file.Filename)
			baseWithoutExt := strings.TrimSuffix(base, ext)
			baseFileName := filepath.Join(extraManifestPath, baseWithoutExt)
			fileName := fmt.Sprintf("%s-%d%s", baseFileName, n, ext)

			extraFile := ignition.FileFromBytes(fileName, user, mode, m)
			config.Storage.Files = append(config.Storage.Files, extraFile)
		}
	}

	return nil
}

func getOSImagesInfo(cpuArch string, openshiftVersion string, streamGetter CoreOSBuildFetcher) (*models.OsImage, error) {
	st, err := streamGetter(context.Background())
	if err != nil {
		return nil, err
	}

	osImage := &models.OsImage{
		CPUArchitecture: &cpuArch,
	}
	osImage.OpenshiftVersion = &openshiftVersion

	streamArch, err := st.GetArchitecture(cpuArch)
	if err != nil {
		return nil, err
	}

	artifacts, ok := streamArch.Artifacts["metal"]
	if !ok {
		return nil, fmt.Errorf("failed to retrieve coreos metal info for architecture %s", cpuArch)
	}
	osImage.Version = &artifacts.Release

	isoFormat, ok := artifacts.Formats["iso"]
	if !ok {
		return nil, fmt.Errorf("failed to retrieve coreos ISO info for architecture %s", cpuArch)
	}
	osImage.URL = &isoFormat.Disk.Location

	return osImage, nil
}

// RetrieveRendezvousIP Returns the Rendezvous IP from either AgentConfig or NMStateConfig
func RetrieveRendezvousIP(agentConfig *agent.Config, hosts []agent.Host, nmStateConfigs []*v1beta1.NMStateConfig) (string, error) {
	var err error
	var rendezvousIP string

	if agentConfig != nil && agentConfig.RendezvousIP != "" {
		rendezvousIP = agentConfig.RendezvousIP
		logrus.Debug("RendezvousIP from the AgentConfig ", rendezvousIP)

	} else {
		rendezvousIP, err = manifests.GetNodeZeroIP(hosts, nmStateConfigs)
		if err != nil {
			return "", errors.Wrap(err, "missing rendezvousIP in agent-config, at least one host networkConfig, or at least one NMStateConfig manifest")
		}
		logrus.Debug("RendezvousIP from the NMStateConfig ", rendezvousIP)
	}

	// Convert IPv6 address to canonical to match host format for comparisons
	addr := net.ParseIP(rendezvousIP)
	if addr == nil {
		err = errors.New(fmt.Sprintf("invalid rendezvous IP: %s", rendezvousIP))
		return "", err
	}
	return addr.String(), err
}

func getPublicContainerRegistries(registriesConfig *mirror.RegistriesConf) string {

	if len(registriesConfig.MirrorConfig) > 0 {
		registries := []string{}
		for _, config := range registriesConfig.MirrorConfig {
			location := strings.SplitN(config.Location, "/", 2)[0]

			allRegs := fmt.Sprint(registries)
			if !strings.Contains(allRegs, location) {
				registries = append(registries, location)
			}
		}
		return strings.Join(registries, ",")
	}

	return "quay.io"
}
