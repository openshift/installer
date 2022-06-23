package image

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/tls"
)

const manifestPath = "/etc/assisted/manifests"
const nmConnectionsPath = "/etc/assisted/network"
const mirrorPath = "/etc/assisted/mirror"

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
	PullSecretToken       string
	NodeZeroIP            string
	AssistedServiceHost   string
	APIVIP                string
	ControlPlaneAgents    int
	WorkerAgents          int
	ReleaseImages         string
	ReleaseImageMirror    string
	MirrorRegistriesMount string
	CaBundleMount         string
	InfraEnvID            string
}

var (
	agentEnabledServices = []string{
		"agent.service",
		"assisted-service-db.service",
		"assisted-service-pod.service",
		"assisted-service-ui.service",
		"assisted-service.service",
		"create-cluster-and-infraenv.service",
		"node-zero.service",
		"multipathd.service",
		"pre-network-manager-config.service",
		"selinux.service",
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
		&tls.KubeAPIServerLBSignerCertKey{},
		&tls.KubeAPIServerLocalhostSignerCertKey{},
		&tls.KubeAPIServerServiceNetworkSignerCertKey{},
		&tls.AdminKubeConfigSignerCertKey{},
		&tls.AdminKubeConfigClientCertKey{},
		&agentconfig.Asset{},
		&mirror.AgentMirror{},
	}
}

// Generate generates the agent installer ignition.
func (a *Ignition) Generate(dependencies asset.Parents) error {
	agentManifests := &manifests.AgentManifests{}
	agentConfigAsset := &agentconfig.Asset{}
	dependencies.Get(agentManifests, agentConfigAsset)

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

	nodeZeroIP, err := manifests.GetNodeZeroIP(agentManifests.NMStateConfigs)
	if err != nil {
		return err
	}

	// TODO: don't hard-code target arch
	releaseImageList, err := releaseImageList(agentManifests.ClusterImageSet.Spec.ReleaseImage, "x86_64")
	if err != nil {
		return err
	}

	agentMirror := &mirror.AgentMirror{}
	dependencies.Get(agentMirror)

	// Mount files for assisted-service
	mirrorRegistriesMount := ""
	caBundleMount := ""
	for _, file := range agentMirror.FileList {
		if file.Filename == mirror.RegistriesConfFilename {
			mirrorRegistriesMount = fmt.Sprintf("-v %s:/etc/containers/registries.conf:z", filepath.Join("/etc/assisted", file.Filename))
		}
		if file.Filename == mirror.CaBundleFilename {
			caBundleMount = fmt.Sprintf("-v %s:/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem:z", filepath.Join("/etc/assisted", file.Filename))
		}
	}

	// Get the mirror for release image
	releaseImageMirror := ""
	source := strings.Split(agentManifests.ClusterImageSet.Spec.ReleaseImage, ":")
	for _, config := range agentMirror.MirrorConfig {
		if config.Location == source[0] {
			// include the tag with the build release image
			releaseImageMirror = fmt.Sprintf("%s:%s", config.Mirror, source[1])
		}
	}

	infraEnvID := uuid.New().String()
	logrus.Debug("Generated random infra-env id ", infraEnvID)

	agentTemplateData := getTemplateData(
		agentManifests.GetPullSecretData(),
		nodeZeroIP,
		releaseImageList,
		releaseImageMirror,
		mirrorRegistriesMount,
		caBundleMount,
		agentManifests.AgentClusterInstall,
		infraEnvID)

	err = bootstrap.AddStorageFiles(&config, "/", "agent/files", agentTemplateData)
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

	err = addStaticNetworkConfig(&config, agentManifests.StaticNetworkConfigs)
	if err != nil {
		return err
	}

	err = bootstrap.AddSystemdUnits(&config, "agent/systemd/units", agentTemplateData, agentEnabledServices)
	if err != nil {
		return err
	}

	addTLSData(&config, dependencies)

	addMirrorData(&config, agentMirror)

	a.Config = &config
	return nil
}

func getTemplateData(pullSecret string, nodeZeroIP string, releaseImageList string,
	releaseImageMirror string, mirrorRegistriesMount string, caBundleMount string, agentClusterInstall *hiveext.AgentClusterInstall, infraEnvID string) *agentTemplateData {
	serviceBaseURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(nodeZeroIP, "8090"),
		Path:   "/",
	}

	return &agentTemplateData{
		ServiceProtocol:       serviceBaseURL.Scheme,
		ServiceBaseURL:        serviceBaseURL.String(),
		PullSecret:            pullSecret,
		PullSecretToken:       "",
		NodeZeroIP:            serviceBaseURL.Hostname(),
		AssistedServiceHost:   serviceBaseURL.Host,
		APIVIP:                agentClusterInstall.Spec.APIVIP,
		ControlPlaneAgents:    agentClusterInstall.Spec.ProvisionRequirements.ControlPlaneAgents,
		WorkerAgents:          agentClusterInstall.Spec.ProvisionRequirements.WorkerAgents,
		ReleaseImages:         releaseImageList,
		ReleaseImageMirror:    releaseImageMirror,
		MirrorRegistriesMount: mirrorRegistriesMount,
		CaBundleMount:         caBundleMount,
		InfraEnvID:            infraEnvID,
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

func addMirrorData(config *igntypes.Config, agentMirror *mirror.AgentMirror) {

	// add mirror files to ignition
	for _, file := range agentMirror.FileList {
		// These are required for assisted-service to build the ICSP for openshift-install
		mirrorFile := ignition.FileFromBytes(filepath.Join(mirrorPath, filepath.Base(file.Filename)),
			"root", 0600, file.Data)
		config.Storage.Files = append(config.Storage.Files, mirrorFile)

		// This is required for the agent to run the podman commands to the mirror
		if file.Filename == mirror.CaBundleFilename {
			mirrorFile := ignition.FileFromBytes("/etc/pki/ca-trust/source/anchors/domain.crt",
				"root", 0600, file.Data)
			config.Storage.Files = append(config.Storage.Files, mirrorFile)

		}
	}
}
