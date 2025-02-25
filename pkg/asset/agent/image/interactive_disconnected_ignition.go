package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/common"
	"github.com/openshift/installer/pkg/asset/agent/interactive"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	interactiveDisconnectedIgnitionFilename = "interactive-disconnected-agent.ign"
)

// InteractiveDisconnectedIgnition is the asset responsible for generating the
// ignition file required to support the interactive disconnected setup.
type InteractiveDisconnectedIgnition struct {
	Config *igntypes.Config
	File   *asset.File
}

// Name returns the human-friendly name of the asset.
func (i *InteractiveDisconnectedIgnition) Name() string {
	return "Agent Installer Interactive Disconnected Ignition"
}

// Dependencies returns the assets on which the current asset depends.
func (i *InteractiveDisconnectedIgnition) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&interactive.InstallConfig{},
		&manifests.ClusterImageSet{},
		&manifests.AgentPullSecret{},
		&manifests.InfraEnv{},
		&common.InfraEnvID{},
	}
}

// Generate produces an asset used only for the interactive disconnected workflow.
// In the connected UI, the only mandatory installation detail required is the release version (and architecture).
// Optionally the user may specify the pull secret, an SSH key or the list of the OLM operators to be installed.
// All the remaining configuration will be provided by the user subsequently in the disconnected environment,
// through the assisted UI running on the rendezvous node.
func (i *InteractiveDisconnectedIgnition) Generate(_ context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	interactiveInstallConfig := &interactive.InstallConfig{}
	clusterImageSet := &manifests.ClusterImageSet{}
	agentPullSecret := &manifests.AgentPullSecret{}
	infraEnv := &manifests.InfraEnv{}
	infraEnvID := &common.InfraEnvID{}
	dependencies.Get(agentWorkflow, interactiveInstallConfig, clusterImageSet, agentPullSecret, infraEnv, infraEnvID)

	if agentWorkflow.Workflow != workflow.AgentWorkflowTypeInstallInteractiveDisconnected {
		return fmt.Errorf("AgentWorkflowType value not supported: %s", agentWorkflow.Workflow)
	}

	// By default, no user is added.
	config := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}
	// Add the optional ssh key for the "core" user if configured.
	if interactiveInstallConfig.SSHKey() != "" {
		config.Passwd = igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name: "core",
					SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
						igntypes.SSHAuthorizedKey(interactiveInstallConfig.SSHKey()),
					},
				},
			},
		}
	}

	// Get the current release architecture and version, to define the complete release image list.
	openshiftArch, err := i.getCurrentReleaseArch()
	if err != nil {
		return err
	}
	openshiftVersion, err := version.Version()
	if err != nil {
		return err
	}
	openShiftReleaseImage := clusterImageSet.Config.Spec.ReleaseImage // Note: in the current workflow it's always equivalent to releaseimage.Image{}.PullSpec

	releaseImageList, err := releaseImageListWithVersion(openShiftReleaseImage, openshiftArch, []string{openshiftArch}, openshiftVersion)
	if err != nil {
		return err
	}
	// Get the current OS image.
	osImage, err := getOSImagesInfo(openshiftArch, openshiftVersion, DefaultCoreOSStreamGetter)
	if err != nil {
		return err
	}

	// Prepare (just) the required template input parameters.
	agentTemplateData := &agentTemplateData{
		PullSecret:    interactiveInstallConfig.PullSecret(),
		ReleaseImages: releaseImageList,
		ReleaseImage:  openShiftReleaseImage,
		InfraEnvID:    infraEnvID.ID,
		OSImage:       osImage,
	}

	// Add the agent files to ignition. The list of files added could be improved,
	// by skipping those ones not really required for the current workflow.
	err = bootstrap.AddStorageFiles(&config, "/", "agent/files", agentTemplateData)
	if err != nil {
		return err
	}
	// Add the required bootstrap files.
	err = addBootstrapScripts(&config, openShiftReleaseImage)
	if err != nil {
		return err
	}
	// Add the rendezvous host file. Agent TUI will interact with that file in case
	// the rendezvous IP wasn't previously configured.
	rendezvousHostFile := ignition.FileFromString(rendezvousHostEnvPath,
		"root", 0644,
		getRendezvousHostEnv("http", interactiveInstallConfig.RendezvousIP(), "", "", agentWorkflow.Workflow))
	config.Storage.Files = append(config.Storage.Files, rendezvousHostFile)
	// Add the only required agent manifests.
	for _, file := range []asset.File{
		*infraEnv.File,
		*clusterImageSet.File,
		*agentPullSecret.File,
	} {
		manifestFile := ignition.FileFromBytes(filepath.Join(manifestPath, filepath.Base(file.Filename)),
			"root", 0600, file.Data)
		config.Storage.Files = append(config.Storage.Files, manifestFile)
	}

	// Configure the ignition services.
	enabledServices := getDefaultEnabledServices()
	err = bootstrap.AddSystemdUnits(&config, "agent/systemd/units", agentTemplateData, enabledServices)
	if err != nil {
		return err
	}
	// TBD: Add assisted UI related services.

	// Generate the ignition file.
	i.Config = &config
	return i.generateFile()
}

func (i *InteractiveDisconnectedIgnition) getCurrentReleaseArch() (string, error) {
	releaseArch, err := version.ReleaseArchitecture()
	if err != nil {
		return "", err
	}
	if strings.Contains(releaseArch, "multi") || strings.Contains(releaseArch, "unknown") {
		releaseArch = string(version.DefaultArch())
	}
	return arch.RpmArch(releaseArch), nil
}

func (i *InteractiveDisconnectedIgnition) generateFile() error {
	data, err := ignition.Marshal(i.Config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal InteractiveDisconnectedIgnition config")
	}
	i.File = &asset.File{
		Filename: interactiveDisconnectedIgnitionFilename,
		Data:     data,
	}
	return nil
}

// PersistToFile writes the unconfigured ignition in the assets folder.
func (i *InteractiveDisconnectedIgnition) PersistToFile(directory string) error {
	if i.File == nil {
		return errors.New("attempting to persist a InteractiveDisconnectedIgnition that has not been generated")
	}
	ignitionFile := filepath.Join(directory, i.File.Filename)

	err := os.WriteFile(ignitionFile, i.File.Data, 0o644) //nolint:gosec // no sensitive info
	if err != nil {
		return err
	}

	logrus.Infof("ignition created in: %s", ignitionFile)

	return nil
}

// Load returns the UnconfiguredIgnition from disk.
func (*InteractiveDisconnectedIgnition) Load(f asset.FileFetcher) (bool, error) {
	// The InteractiveDisconnectedIgnition will not be needed by another asset so load is noop.
	// This is implemented because it is required by WritableAsset
	return false, nil
}

// Files returns the files generated by the asset.
func (*InteractiveDisconnectedIgnition) Files() []*asset.File {
	// Return empty array because File will never be loaded.
	return []*asset.File{}
}
