package image

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coreos/ignition/v2/config/merge"
	"github.com/coreos/ignition/v2/config/v3_2"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/types"
)

const (
	trustedBundlePath = "/etc/pki/ca-trust/source/anchors/additional-trust-bundle.pem"

	registriesConfPath = "/etc/containers/registries.conf"

	postDeploymentScriptPath = "/var/tmp/post.sh"
)

// Ignition is an asset that generates the image-based installer ignition file.
type Ignition struct {
	Config *igntypes.Config
}

// Name returns the human-friendly name of the asset.
func (i *Ignition) Name() string {
	return "Image-based Installer Ignition"
}

// Dependencies returns the assets on which the Ignition asset depends.
func (i *Ignition) Dependencies() []asset.Asset {
	return []asset.Asset{
		&ImageBasedInstallationConfig{},
		&RegistriesConf{},
		&PostDeployment{},
	}
}

type ibiConfigurationFile struct {
	ExtraPartitionLabel  string   `json:"extraPartitionLabel,omitempty"`
	ExtraPartitionNumber uint     `json:"extraPartitionNumber,omitempty"`
	ExtraPartitionStart  string   `json:"extraPartitionStart,omitempty"`
	InstallationDisk     string   `json:"installationDisk"`
	ReleaseRegistry      string   `json:"releaseRegistry,omitempty"`
	SeedImage            string   `json:"seedImage"`
	SeedVersion          string   `json:"seedVersion"`
	Shutdown             bool     `json:"shutdown,omitempty"`
	SkipDiskCleanup      bool     `json:"skipDiskCleanup,omitempty"`
	CoreosInstallerArgs  []string `json:"coreosInstallerArgs,omitempty"`
}

type ibiTemplateData struct {
	SeedImage        string
	Proxy            *types.Proxy
	PullSecret       string
	NetworkConfig    string
	IBIConfiguration string
}

// Generate generates the image-based installer ignition.
func (i *Ignition) Generate(_ context.Context, dependencies asset.Parents) error {
	configAsset := &ImageBasedInstallationConfig{}
	registriesConf := &RegistriesConf{}
	postDeployment := &PostDeployment{}

	dependencies.Get(configAsset, registriesConf, postDeployment)

	ibiConfig := configAsset.Config
	if ibiConfig == nil {
		return fmt.Errorf("%s is required", configFilename)
	}

	config := &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
		Passwd: igntypes.Passwd{
			Users: []igntypes.PasswdUser{
				{
					Name: "core",
					SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
						igntypes.SSHAuthorizedKey(ibiConfig.SSHKey),
					},
				},
			},
		},
	}

	ibiConfigFile := ibiConfigurationFile{
		ExtraPartitionLabel:  ibiConfig.ExtraPartitionLabel,
		ExtraPartitionNumber: ibiConfig.ExtraPartitionNumber,
		ExtraPartitionStart:  ibiConfig.ExtraPartitionStart,
		InstallationDisk:     ibiConfig.InstallationDisk,
		ReleaseRegistry:      ibiConfig.ReleaseRegistry,
		SeedVersion:          ibiConfig.SeedVersion,
		SeedImage:            ibiConfig.SeedImage,
		Shutdown:             ibiConfig.Shutdown,
		SkipDiskCleanup:      ibiConfig.SkipDiskCleanup,
		CoreosInstallerArgs:  ibiConfig.CoreosInstallerArgs,
	}
	ibiConfigJSON, err := json.Marshal(ibiConfigFile)
	if err != nil {
		return fmt.Errorf("failed to marshall the ibi-configuration data: %w", err)
	}

	ibiTemplateData := &ibiTemplateData{
		SeedImage:        ibiConfig.SeedImage,
		Proxy:            ibiConfig.Proxy,
		PullSecret:       ibiConfig.PullSecret,
		IBIConfiguration: string(ibiConfigJSON),
	}

	if ibiConfig.NetworkConfig != nil {
		ibiTemplateData.NetworkConfig = ibiConfig.NetworkConfig.String()
	}

	if len(registriesConf.Data) > 0 {
		file := ignition.FileFromString(registriesConfPath, "root", 0o644, string(registriesConf.Data))
		config.Storage.Files = append(config.Storage.Files, file)
	}

	if ibiConfig.AdditionalTrustBundle != "" {
		file := ignition.FileFromString(trustedBundlePath, "root", 0o600, ibiConfig.AdditionalTrustBundle)
		config.Storage.Files = append(config.Storage.Files, file)
	}

	if postDeployment.File != nil {
		file := ignition.FileFromString(postDeploymentScriptPath, "root", 0o755, string(postDeployment.File.Data))
		config.Storage.Files = append(config.Storage.Files, file)
	}

	if ibiConfig.IgnitionConfigOverride != "" {
		if err := setIgnitionConfigOverride(config, ibiConfig.IgnitionConfigOverride); err != nil {
			return fmt.Errorf("failed to override ignition config: %w", err)
		}
	}

	if err := bootstrap.AddStorageFiles(config, "/", "imagebased/files", ibiTemplateData); err != nil {
		return fmt.Errorf("failed to add image-based files to ignition config: %w", err)
	}

	enabledServices := defaultEnabledServices()
	if ibiConfig.NetworkConfig != nil && ibiConfig.NetworkConfig.String() != "" {
		enabledServices = append(enabledServices, "network-config.service")
	}
	if err := bootstrap.AddSystemdUnits(config, "imagebased/systemd/units", ibiTemplateData, enabledServices); err != nil {
		return fmt.Errorf("failed to add image-based systemd units to ignition config: %w", err)
	}

	i.Config = config

	return nil
}

func setIgnitionConfigOverride(config *igntypes.Config, override string) error {
	ignitionConfigOverride, _, err := v3_2.Parse([]byte(override))
	if err != nil {
		return fmt.Errorf("failed to parse ignition config override: %w", err)
	}

	merged, _ := merge.MergeStructTranscribe(*config, ignitionConfigOverride)
	marshaledMerged, err := json.Marshal(merged)
	if err != nil {
		return fmt.Errorf("failed to marshal merged ignition config: %w", err)
	}

	if err := json.Unmarshal(marshaledMerged, config); err != nil {
		return fmt.Errorf("failed to unmarshal merged ignition config: %w", err)
	}
	return nil
}

func defaultEnabledServices() []string {
	return []string{
		"install-rhcos-and-restore-seed.service",
	}
}
