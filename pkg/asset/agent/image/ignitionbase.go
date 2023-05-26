package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

const (
	unconfiguredIgnitionFilename = "unconfigured-agent.ign"
)

// IgnitionBase is an asset that generates the agent installer base ignition filea
// which excludes any cluster configuration files.
type IgnitionBase struct {
	Config                    *igntypes.Config
	CPUArch                   string
	File                      *asset.File
	infraEnvID                string
	archName                  string
	osImage                   models.OsImage
	hasMirrorConfig           bool
	releaseImageMirror        string
	publicContainerRegistries string
}

var agentEnabledServices = []string{
	"agent-interactive-console.service",
	"agent.service",
	"assisted-service-db.service",
	"assisted-service-pod.service",
	"assisted-service.service",
	"create-cluster-and-infraenv.service",
	"node-zero.service",
	"multipathd.service",
	"selinux.service",
	"install-status.service",
	// Services disabled in IgnitionBase
	// "set-hostname.service",
	// "start-cluster-installation.service",
}

// Name returns the human-friendly name of the asset.
func (a *IgnitionBase) Name() string {
	return "Agent Installer Ignition Base"
}

// Dependencies returns the assets on which the Ignition asset depends.
func (a *IgnitionBase) Dependencies() []asset.Asset {
	return []asset.Asset{
		&manifests.InfraEnv{},
		&manifests.AgentPullSecret{},
		&manifests.ClusterImageSet{},
		&manifests.NMStateConfig{},
		&mirror.RegistriesConf{},
		&mirror.CaBundle{},
	}
}

// Generate generates the agent installer base ignition.
func (a *IgnitionBase) Generate(dependencies asset.Parents) error {
	infraEnvAsset := &manifests.InfraEnv{}
	clusterImageSetAsset := &manifests.ClusterImageSet{}
	pullSecretAsset := &manifests.AgentPullSecret{}
	nmStateConfigs := &manifests.NMStateConfig{}
	dependencies.Get(infraEnvAsset, clusterImageSetAsset, pullSecretAsset, nmStateConfigs)

	infraEnv := infraEnvAsset.Config
	clusterImageSet := clusterImageSetAsset.Config

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
					// PasswordHash: &pwdHash,
				},
			},
		},
	}

	// Default to x86_64
	archName := arch.RpmArch(types.ArchitectureAMD64)
	if infraEnv.Spec.CpuArchitecture != "" {
		archName = infraEnv.Spec.CpuArchitecture
	}
	a.archName = archName

	releaseImageList, err := releaseImageList(clusterImageSet.Spec.ReleaseImage, archName)
	if err != nil {
		return err
	}

	registriesConfig := &mirror.RegistriesConf{}
	registryCABundle := &mirror.CaBundle{}
	dependencies.Get(registriesConfig, registryCABundle)

	a.hasMirrorConfig = len(registriesConfig.MirrorConfig) > 0
	a.publicContainerRegistries = getPublicContainerRegistries(registriesConfig)
	a.releaseImageMirror = mirror.GetMirrorFromRelease(clusterImageSet.Spec.ReleaseImage, registriesConfig)

	infraEnvID := uuid.New().String()
	logrus.Debug("Generated random infra-env id ", infraEnvID)
	a.infraEnvID = infraEnvID

	osImage, err := getOSImagesInfo(archName)
	if err != nil {
		return err
	}
	a.osImage = *osImage
	a.CPUArch = *osImage.CPUArchitecture

	agentTemplateData := &agentTemplateData{
		PullSecret:                pullSecretAsset.GetPullSecretData(),
		ReleaseImages:             releaseImageList,
		ReleaseImage:              clusterImageSet.Spec.ReleaseImage,
		ReleaseImageMirror:        a.releaseImageMirror,
		HaveMirrorConfig:          a.hasMirrorConfig,
		PublicContainerRegistries: a.publicContainerRegistries,
		InfraEnvID:                a.infraEnvID,
		OSImage:                   a.osImage,
		Proxy:                     infraEnv.Spec.Proxy,
	}

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
		ReleaseImage: clusterImageSet.Spec.ReleaseImage,
	}
	for _, script := range []string{"release-image.sh", "release-image-download.sh"} {
		if err := bootstrap.AddStorageFiles(&config,
			"/usr/local/bin/"+script,
			"bootstrap/files/usr/local/bin/"+script+".template",
			relImgData); err != nil {
			return err
		}
	}

	err = addStaticNetworkConfig(&config, nmStateConfigs.StaticNetworkConfig)
	if err != nil {
		return err
	}

	// Enable pre-network-manager-config.service only when there are network configs defined
	if len(nmStateConfigs.StaticNetworkConfig) != 0 {
		agentEnabledServices = append(agentEnabledServices, "pre-network-manager-config.service")
	}

	filesToInclude := [...]asset.File{
		*infraEnvAsset.File,
		*clusterImageSetAsset.File,
		*pullSecretAsset.File,
	}

	for _, file := range filesToInclude {
		manifestFile := ignition.FileFromBytes(filepath.Join(manifestPath, filepath.Base(file.Filename)),
			"root", 0600, file.Data)
		config.Storage.Files = append(config.Storage.Files, manifestFile)
	}

	err = bootstrap.AddSystemdUnits(&config, "agent/systemd/units", agentTemplateData, agentEnabledServices)
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

// CopyIgnitionConfig returns a copy of the Ignition config
// The copy is made by marshaling and unmarshaling the ignition yaml.
func (a *IgnitionBase) CopyIgnitionConfig() (copy igntypes.Config, error error) {
	copyData := &igntypes.Config{}
	data, err := yaml.Marshal(a.Config)
	if err != nil {
		return copy, errors.Wrap(err, "copy failed, failed to Marshal IgnitionBase config")
	}

	if err := yaml.UnmarshalStrict(data, copyData); err != nil {
		return copy, errors.Wrap(err, "copy failed, failed to Unmarshal IgnitionBase config")
	}

	return *copyData, nil
}

// PersistToFile writes the unconfigured ignition in the assets folder.
func (a *IgnitionBase) PersistToFile(directory string) error {
	if a.File == nil {
		return errors.New("attempting to persist a IgnitionBase that has not been generated")
	}
	unconfiguredIgnFile := filepath.Join(directory, a.File.Filename)

	err := os.WriteFile(unconfiguredIgnFile, a.File.Data, 0o644) //nolint:gosec // no sensitive info
	if err != nil {
		return err
	}

	return nil
}

func (a *IgnitionBase) generateFile(filename string) error {
	data, err := ignition.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal Ignition config")
	}
	a.File = &asset.File{
		Filename: filename,
		Data:     data,
	}
	return nil
}

// Load returns the IgnitionBase from disk.
func (a *IgnitionBase) Load(f asset.FileFetcher) (bool, error) {
	// The IgnitionBase will not be needed by another asset so load is noop.
	// This is implemented because it is required by WritableAsset
	return false, nil
}

// Files returns the files generated by the asset.
func (a *IgnitionBase) Files() []*asset.File {
	// Return empty array because File will never be loaded.
	return []*asset.File{}
}

func getOSImagesInfo(cpuArch string) (*models.OsImage, error) {
	st, err := rhcos.FetchCoreOSBuild(context.Background())
	if err != nil {
		return nil, err
	}

	osImage := &models.OsImage{
		CPUArchitecture: &cpuArch,
	}

	openshiftVersion, err := version.Version()
	if err != nil {
		return nil, err
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
