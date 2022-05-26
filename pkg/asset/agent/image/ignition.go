package image

import (
	"errors"
	"net"
	"net/url"
	"path/filepath"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/sirupsen/logrus"
)

const manifestPathInIso = "/etc/assisted/manifests"

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
	}
}

// Generate generates the agent installer ignition.
func (a *Ignition) Generate(dependencies asset.Parents) error {

	agentManifests := &manifests.AgentManifests{}
	dependencies.Get(agentManifests)

	pullSecret := agentManifests.PullSecret
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

	if pullSecret.StringData[".dockerconfigjson"] == "" {
		return errors.New("pull secret is missing")
	}
	agentTemplateData := getTemplateData(pullSecret.StringData[".dockerconfigjson"], agentManifests.AgentClusterInstall)
	bootstrap.AddStorageFiles(&config, "/", "agent/files", agentTemplateData)

	for _, file := range agentManifests.FileList {
		manifestFile := ignition.FileFromBytes(filepath.Join(manifestPathInIso, filepath.Base(file.Filename)),
			"root", 0600, file.Data)
		config.Storage.Files = append(config.Storage.Files, manifestFile)
	}

	logrus.Infof("RWSU number of files: %d", len(config.Storage.Files))

	bootstrap.AddSystemdUnits(&config, "agent/systemd/units", agentTemplateData, agentEnabledServices)

	logrus.Infof("RWSU number of service: %d", len(config.Systemd.Units))

	a.Config = &config
	return nil
}

func getTemplateData(pullSecret string, agentClusterInstall *hiveext.AgentClusterInstall) *agentTemplateData {
	// TODO: determine nodeZeroIP from NMStateConfig
	nodeZeroIP := "192.168.122.2"
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
	}
}
