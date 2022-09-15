package image

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types/agent"
	"github.com/pkg/errors"
)

const manifestPath = "/etc/assisted/manifests"
const hostnamesPath = "/etc/assisted/hostnames"
const nmConnectionsPath = "/etc/assisted/network"
const extraManifestPath = "/etc/assisted/extra-manifests"

// Ignition is an asset that generates the agent installer ignition file.
type Ignition struct {
	Config *igntypes.Config
}

// agentTemplateData is the data used to replace values in agent template
// files.
type agentTemplateData struct {
	ServiceProtocol string
	ServiceBaseURL  string
	PullSecret      string
	// PullSecretToken is token to use for authentication when AUTH_TYPE=rhsso
	// in assisted-service
	PullSecretToken     string
	NodeZeroIP          string
	AssistedServiceHost string
	APIVIP              string
	ControlPlaneAgents  int
	WorkerAgents        int
	ReleaseImages       string
	ReleaseImage        string
	ReleaseImageMirror  string
	HaveMirrorConfig    bool
	InfraEnvID          string
}

var (
	agentEnabledServices = []string{
		"agent.service",
		"assisted-service-db.service",
		"assisted-service-pod.service",
		"assisted-service.service",
		"create-cluster-and-infraenv.service",
		"node-zero.service",
		"multipathd.service",
		"pre-network-manager-config.service",
		"selinux.service",
		"set-hostname.service",
		"start-cluster-installation.service",
	}
)

// Name returns the human-friendly name of the asset.
func (a *Ignition) Name() string {
	return "Agent Installer Ignition"
}

// Dependencies returns the assets on which the Ignition asset depends.
func (a *Ignition) Dependencies() []asset.Asset {
	return []asset.Asset{
		&manifests.AgentManifests{},
		&manifests.ExtraManifests{},
		&tls.KubeAPIServerLBSignerCertKey{},
		&tls.KubeAPIServerLocalhostSignerCertKey{},
		&tls.KubeAPIServerServiceNetworkSignerCertKey{},
		&tls.AdminKubeConfigSignerCertKey{},
		&tls.AdminKubeConfigClientCertKey{},
		&agentconfig.AgentConfig{},
		&mirror.RegistriesConf{},
		&mirror.CaBundle{},
	}
}

// Generate generates the agent installer ignition.
func (a *Ignition) Generate(dependencies asset.Parents) error {
	agentManifests := &manifests.AgentManifests{}
	agentConfigAsset := &agentconfig.AgentConfig{}
	extraManifests := &manifests.ExtraManifests{}
	dependencies.Get(agentManifests, agentConfigAsset, extraManifests)

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
				},
			},
		},
	}

	nodeZeroIP, err := RetrieveRendezvousIP(agentConfigAsset.Config, agentManifests.NMStateConfigs)
	if err != nil {
		return err
	}

	// TODO: don't hard-code target arch
	releaseImageList, err := releaseImageList(agentManifests.ClusterImageSet.Spec.ReleaseImage, "x86_64")
	if err != nil {
		return err
	}

	registriesConfig := &mirror.RegistriesConf{}
	registryCABundle := &mirror.CaBundle{}
	dependencies.Get(registriesConfig, registryCABundle)

	releaseImageMirror := getMirrorFromRelease(agentManifests.ClusterImageSet.Spec.ReleaseImage, registriesConfig)

	infraEnvID := uuid.New().String()
	logrus.Debug("Generated random infra-env id ", infraEnvID)

	agentTemplateData := getTemplateData(
		agentManifests.GetPullSecretData(),
		nodeZeroIP,
		releaseImageList,
		agentManifests.ClusterImageSet.Spec.ReleaseImage,
		releaseImageMirror,
		len(registriesConfig.MirrorConfig) > 0,
		agentManifests.AgentClusterInstall,
		infraEnvID)

	err = bootstrap.AddStorageFiles(&config, "/", "agent/files", agentTemplateData)
	if err != nil {
		return err
	}

	// Set up bootstrap service recording
	if err := bootstrap.AddStorageFiles(&config,
		"/usr/local/bin/bootstrap-service-record.sh",
		"bootstrap/files/usr/local/bin/bootstrap-service-record.sh",
		nil); err != nil {
		return err
	}

	// Use bootstrap script to get container images
	relImgData := struct{ ReleaseImage string }{
		ReleaseImage: agentManifests.ClusterImageSet.Spec.ReleaseImage,
	}
	for _, script := range []string{"release-image.sh", "release-image-download.sh"} {
		if err := bootstrap.AddStorageFiles(&config,
			"/usr/local/bin/"+script,
			"bootstrap/files/usr/local/bin/"+script+".template",
			relImgData); err != nil {
			return err
		}
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

	addMacAddressToHostnameMappings(&config, agentConfigAsset)

	err = addStaticNetworkConfig(&config, agentManifests.StaticNetworkConfigs)
	if err != nil {
		return err
	}

	err = bootstrap.AddSystemdUnits(&config, "agent/systemd/units", agentTemplateData, agentEnabledServices)
	if err != nil {
		return err
	}

	addTLSData(&config, dependencies)

	addMirrorData(&config, registriesConfig, registryCABundle)

	addHostConfig(&config, agentConfigAsset)

	addExtraManifests(&config, extraManifests)

	a.Config = &config
	return nil
}

func getTemplateData(pullSecret, nodeZeroIP, releaseImageList, releaseImage,
	releaseImageMirror string, haveMirrorConfig bool,
	agentClusterInstall *hiveext.AgentClusterInstall,
	infraEnvID string) *agentTemplateData {
	serviceBaseURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(nodeZeroIP, "8090"),
		Path:   "/",
	}

	return &agentTemplateData{
		ServiceProtocol:     serviceBaseURL.Scheme,
		ServiceBaseURL:      serviceBaseURL.String(),
		PullSecret:          pullSecret,
		PullSecretToken:     "",
		NodeZeroIP:          serviceBaseURL.Hostname(),
		AssistedServiceHost: serviceBaseURL.Host,
		APIVIP:              agentClusterInstall.Spec.APIVIP,
		ControlPlaneAgents:  agentClusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents,
		WorkerAgents:        agentClusterInstall.Spec.ProvisionRequirements.WorkerAgents,
		ReleaseImages:       releaseImageList,
		ReleaseImage:        releaseImage,
		ReleaseImageMirror:  releaseImageMirror,
		HaveMirrorConfig:    haveMirrorConfig,
		InfraEnvID:          infraEnvID,
	}
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
		&tls.AdminKubeConfigClientCertKey{},
	}
	dependencies.Get(certKeys...)

	for _, ck := range certKeys {
		for _, d := range ck.(asset.WritableAsset).Files() {
			f := ignition.FileFromBytes(path.Join("/opt/agent", d.Filename), "root", 0600, d.Data)
			config.Storage.Files = append(config.Storage.Files, f)
		}
	}
}

func addMirrorData(config *igntypes.Config, registriesConfig *mirror.RegistriesConf, registryCABundle *mirror.CaBundle) {

	// This is required for assisted-service to build the ICSP for openshift-install
	if registriesConfig.File != nil {
		registriesFile := ignition.FileFromBytes("/etc/containers/registries.conf",
			"root", 0600, registriesConfig.File.Data)
		config.Storage.Files = append(config.Storage.Files, registriesFile)
	}

	// This is required for the agent to run the podman commands to the mirror
	if registryCABundle.File != nil && len(registryCABundle.File.Data) > 0 {
		caFile := ignition.FileFromBytes("/etc/pki/ca-trust/source/anchors/domain.crt",
			"root", 0600, registryCABundle.File.Data)
		config.Storage.Files = append(config.Storage.Files, caFile)
	}
}

// Creates a file named with a host's MAC address. The desired hostname
// is the file's content. The files are read by a systemd service that
// sets the hostname using "hostnamectl set-hostname" when the ISO boots.
func addMacAddressToHostnameMappings(
	config *igntypes.Config,
	agentConfigAsset *agentconfig.AgentConfig) {
	if agentConfigAsset.Config == nil || len(agentConfigAsset.Config.Hosts) == 0 {
		return
	}
	for _, host := range agentConfigAsset.Config.Hosts {
		if host.Hostname != "" {
			file := ignition.FileFromBytes(filepath.Join(hostnamesPath,
				strings.ToLower(filepath.Base(host.Interfaces[0].MacAddress))),
				"root", 0600, []byte(host.Hostname))
			config.Storage.Files = append(config.Storage.Files, file)
		}
	}
}

func addHostConfig(config *igntypes.Config, agentConfig *agentconfig.AgentConfig) error {
	confs, err := agentConfig.HostConfigFiles()
	if err != nil {
		return err
	}

	for path, content := range confs {
		hostConfigFile := ignition.FileFromBytes(filepath.Join("/etc/assisted/hostconfig", path), "root", 0644, content)
		config.Storage.Files = append(config.Storage.Files, hostConfigFile)
	}
	return nil
}

func addExtraManifests(config *igntypes.Config, extraManifests *manifests.ExtraManifests) {

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
		extraFile := ignition.FileFromBytes(filepath.Join(extraManifestPath, filepath.Base(file.Filename)), user, mode, file.Data)
		config.Storage.Files = append(config.Storage.Files, extraFile)
	}
}

// RetrieveRendezvousIP Returns the Rendezvous IP from either AgentConfig or NMStateConfig
func RetrieveRendezvousIP(agentConfig *agent.Config, nmStateConfigs []*v1beta1.NMStateConfig) (string, error) {
	var err error
	var rendezvousIP string

	if agentConfig != nil && agentConfig.RendezvousIP != "" {
		rendezvousIP = agentConfig.RendezvousIP
		logrus.Debug("RendezvousIP from the AgentConfig ", rendezvousIP)

	} else if len(nmStateConfigs) > 0 {
		rendezvousIP, err = manifests.GetNodeZeroIP(nmStateConfigs)
		logrus.Debug("RendezvousIP from the NMStateConfig ", rendezvousIP)
	} else {
		err = errors.New("missing rendezvousIP in agent-config or at least one NMStateConfig manifest")
	}
	return rendezvousIP, err
}

func getMirrorFromRelease(releaseImage string, registriesConfig *mirror.RegistriesConf) string {

	releaseImageMirror := ""
	source := regexp.MustCompile(`^(.+?)(@sha256)?:(.+)`).FindStringSubmatch(releaseImage)
	for _, config := range registriesConfig.MirrorConfig {
		if config.Location == source[1] {
			// include the tag with the build release image
			if len(source) == 4 {
				// Has Sha256
				releaseImageMirror = fmt.Sprintf("%s%s:%s", config.Mirror, source[2], source[3])
			} else if len(source) == 3 {
				releaseImageMirror = fmt.Sprintf("%s:%s", config.Mirror, source[2])
			}
		}
	}

	return releaseImageMirror
}
