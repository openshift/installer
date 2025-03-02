package image

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/common"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

const (
	unconfiguredIgnitionFilename = "unconfigured-agent.ign"
)

// GetConfigImageFiles returns the list of files or file paths to be included in the config-image.
func GetConfigImageFiles() []string {
	return []string{
		"/etc/assisted/manifests/pull-secret.yaml", //nolint:gosec // not hardcoded credentials
		"/etc/assisted/manifests/cluster-deployment.yaml",
		"/etc/assisted/manifests/cluster-image-set.yaml",
		"/etc/assisted/manifests/agent-cluster-install.yaml",
		"/etc/assisted/manifests/infraenv.yaml",
		"/etc/assisted/manifests",       // optional nmstateconfig.yaml
		"/etc/assisted/extra-manifests", // all files in directory
		"/etc/assisted/hostconfig",      // all files in directory
		"/etc/assisted/hostnames",       // all files in directory
		"/etc/assisted/network",         // all files in directory
		"/etc/NetworkManager/conf.d/clientid.conf",
		"/etc/issue",
		"/etc/systemd/system.conf.d/10-default-env.conf",
		"/root/.docker/config.json",
		"/usr/local/share/start-cluster/start-cluster.env",
		"/usr/local/share/assisted-service/assisted-service.env",
		"/opt/agent/tls/kubeadmin-password.hash", //nolint:gosec // not hardcoded credentials
		"/opt/agent/tls/admin-kubeconfig-signer.key",
		"/opt/agent/tls/admin-kubeconfig-signer.crt",
		"/opt/agent/tls/kube-apiserver-lb-signer.key",
		"/opt/agent/tls/kube-apiserver-lb-signer.crt",
		"/opt/agent/tls/kube-apiserver-localhost-signer.key",
		"/opt/agent/tls/kube-apiserver-localhost-signer.crt",
		"/opt/agent/tls/kube-apiserver-service-network-signer.key",
		"/opt/agent/tls/kube-apiserver-service-network-signer.crt",
		rendezvousHostEnvPath, // This file must be last in the list
	}
}

// UnconfiguredIgnition is an asset that generates the agent installer unconfigured
// ignition file which excludes any cluster configuration.
type UnconfiguredIgnition struct {
	Config  *igntypes.Config
	CPUArch string
	File    *asset.File
}

// Name returns the human-friendly name of the asset.
func (a *UnconfiguredIgnition) Name() string {
	return "Agent Installer Unconfigured Ignition"
}

// Dependencies returns the assets on which the UnconfiguredIgnition asset depends.
func (a *UnconfiguredIgnition) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&agentconfig.AgentConfig{},
		&manifests.InfraEnvFile{},
		&manifests.AgentPullSecret{},
		&manifests.ClusterImageSet{},
		&manifests.NMStateConfig{},
		&mirror.RegistriesConf{},
		&mirror.CaBundle{},
		&common.InfraEnvID{},
	}
}

// Generate generates the agent installer unconfigured ignition.
func (a *UnconfiguredIgnition) Generate(_ context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	infraEnvAsset := &manifests.InfraEnvFile{}
	infraEnvIDAsset := &common.InfraEnvID{}
	clusterImageSetAsset := &manifests.ClusterImageSet{}
	pullSecretAsset := &manifests.AgentPullSecret{}
	nmStateConfigs := &manifests.NMStateConfig{}
	agentConfig := &agentconfig.AgentConfig{}
	dependencies.Get(agentWorkflow, infraEnvAsset, clusterImageSetAsset, pullSecretAsset, nmStateConfigs, infraEnvIDAsset, agentConfig)

	infraEnv := infraEnvAsset.Config
	clusterImageSet := clusterImageSetAsset.Config

	config := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}
	// Add the ssh key for the "core" user only if configured.
	if infraEnv.Spec.SSHAuthorizedKey != "" {
		config.Passwd = igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name: "core",
					SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
						igntypes.SSHAuthorizedKey(infraEnv.Spec.SSHAuthorizedKey),
					},
				},
			},
		}
	}

	// Default to x86_64
	archName := arch.RpmArch(types.ArchitectureAMD64)
	if infraEnv.Spec.CpuArchitecture != "" {
		archName = infraEnv.Spec.CpuArchitecture
	}
	releaseImageList, err := releaseImageList(clusterImageSet.Spec.ReleaseImage, archName, []string{archName})
	if err != nil {
		return err
	}

	registriesConfig := &mirror.RegistriesConf{}
	registryCABundle := &mirror.CaBundle{}
	dependencies.Get(registriesConfig, registryCABundle)

	infraEnvID := infraEnvIDAsset.ID
	logrus.Debug("Generated random infra-env id ", infraEnvID)

	openshiftVersion, err := version.Version()
	if err != nil {
		return err
	}
	osImage, err := getOSImagesInfo(archName, openshiftVersion, DefaultCoreOSStreamGetter)
	if err != nil {
		return err
	}
	a.CPUArch = *osImage.CPUArchitecture

	agentTemplateData := &agentTemplateData{
		PullSecret:                pullSecretAsset.GetPullSecretData(),
		ReleaseImages:             releaseImageList,
		ReleaseImage:              clusterImageSet.Spec.ReleaseImage,
		ReleaseImageMirror:        mirror.GetMirrorFromRelease(clusterImageSet.Spec.ReleaseImage, registriesConfig),
		HaveMirrorConfig:          len(registriesConfig.MirrorConfig) > 0,
		PublicContainerRegistries: getPublicContainerRegistries(registriesConfig),
		InfraEnvID:                infraEnvID,
		OSImage:                   osImage,
		Proxy:                     infraEnv.Spec.Proxy,
	}

	enabledServices := getDefaultEnabledServices()

	switch agentWorkflow.Workflow {
	case workflow.AgentWorkflowTypeInstall:
		agentTemplateData.ConfigImageFiles = strings.Join(GetConfigImageFiles(), ",")

		if len(nmStateConfigs.StaticNetworkConfig) > 0 {
			err = addStaticNetworkConfig(&config, nmStateConfigs.StaticNetworkConfig)
			if err != nil {
				return err
			}

			enabledServices = append(enabledServices, "pre-network-manager-config.service")
		} else {
			// Include the script in case it is needed in config step
			nmStateScriptFilePath := "/usr/local/bin/pre-network-manager-config.sh"
			nmStateScript := ignition.FileFromBytes(nmStateScriptFilePath, "root", 0755, []byte(manifests.PreNetworkConfigScript))
			config.Storage.Files = append(config.Storage.Files, nmStateScript)
		}

		// Enable the agent-check-config-image.service for the current workflow.
		enabledServices = append(enabledServices, "agent-check-config-image.service")

	case workflow.AgentWorkflowTypeInstallInteractiveDisconnected:
		// Add the rendezvous host file. Agent TUI will interact with that file in case
		// the rendezvous IP wasn't previously configured.
		rendezvousIP := ""
		if agentConfig.Config != nil {
			rendezvousIP = agentConfig.Config.RendezvousIP
		}
		rendezvousHostFile := ignition.FileFromString(rendezvousHostEnvPath,
			"root", 0644,
			getRendezvousHostEnv("http", rendezvousIP, "none", "none", agentWorkflow.Workflow))
		config.Storage.Files = append(config.Storage.Files, rendezvousHostFile)

		// Explicitly disable the load-config-iso service, not required in the current flow
		// (even though disabled by default, the udev rule may require it).
		config.Storage.Files = append(config.Storage.Files, ignition.FileFromString("/etc/assisted/no-config-image", "root", 0644, ""))

		// Enable the UI service.
		enabledServices = append(enabledServices, "agent-start-ui.service")

		// Let's disable the assisted-service authentication.
		agentTemplateData.AuthType = "none"

		// Required by assisted-service.
		a.ignAddFolders(&config, "/opt/agent/tls", "/etc/assisted/extra-manifests")
	}

	err = bootstrap.AddStorageFiles(&config, "/", "agent/files", agentTemplateData)
	if err != nil {
		return err
	}

	err = addBootstrapScripts(&config, clusterImageSetAsset.Config.Spec.ReleaseImage)
	if err != nil {
		return err
	}

	ztpManifestsToInclude := [...]asset.File{
		*infraEnvAsset.File,
		*clusterImageSetAsset.File,
		*pullSecretAsset.File,
	}

	for _, file := range ztpManifestsToInclude {
		manifestFile := ignition.FileFromBytes(filepath.Join(manifestPath, filepath.Base(file.Filename)),
			"root", 0600, file.Data)
		config.Storage.Files = append(config.Storage.Files, manifestFile)
	}

	err = bootstrap.AddSystemdUnits(&config, "agent/systemd/units", agentTemplateData, enabledServices)
	if err != nil {
		return err
	}

	addMirrorData(&config, registriesConfig, registryCABundle)

	a.Config = &config

	if err := a.generateFile(unconfiguredIgnitionFilename); err != nil {
		return err
	}

	return nil
}

func (a *UnconfiguredIgnition) ignAddFolders(config *igntypes.Config, paths ...string) {
	for _, p := range paths {
		config.Storage.Directories = append(config.Storage.Directories, igntypes.Directory{
			Node: igntypes.Node{
				Path:      p,
				User:      igntypes.NodeUser{Name: ignutil.StrToPtr("root")},
				Overwrite: ignutil.BoolToPtr(true),
			},
			DirectoryEmbedded1: igntypes.DirectoryEmbedded1{Mode: ignutil.IntToPtr(0644)},
		})
	}
}

// PersistToFile writes the unconfigured ignition in the assets folder.
func (a *UnconfiguredIgnition) PersistToFile(directory string) error {
	if a.File == nil {
		return errors.New("attempting to persist a UnconfiguredIgnition that has not been generated")
	}
	unconfiguredIgnFile := filepath.Join(directory, a.File.Filename)

	err := os.WriteFile(unconfiguredIgnFile, a.File.Data, 0o644) //nolint:gosec // no sensitive info
	if err != nil {
		return err
	}

	logrus.Infof("Unconfigured-Ignition created in: %s", unconfiguredIgnFile)

	return nil
}

func (a *UnconfiguredIgnition) generateFile(filename string) error {
	data, err := ignition.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal UnconfiguredIgnition config")
	}
	a.File = &asset.File{
		Filename: filename,
		Data:     data,
	}
	return nil
}

// Load returns the UnconfiguredIgnition from disk.
func (a *UnconfiguredIgnition) Load(f asset.FileFetcher) (bool, error) {
	// The UnconfiguredIgnition will not be needed by another asset so load is noop.
	// This is implemented because it is required by WritableAsset
	return false, nil
}

// Files returns the files generated by the asset.
func (a *UnconfiguredIgnition) Files() []*asset.File {
	// Return empty array because File will never be loaded.
	return []*asset.File{}
}
