package image

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coreos/ignition/v2/config/merge"
	ignutil "github.com/coreos/ignition/v2/config/util"
	"github.com/coreos/ignition/v2/config/v3_2"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/openshift/installer/pkg/asset"
)

const (
	ibiConfigurationPath     = "/var/tmp/ibi-configuration.json"
	installationScriptPath   = "/usr/local/bin/install-rhcos-and-restore-seed.sh"
	installationServiceName  = "install-rhcos-and-restore-seed.service"
	networkConfigServiceName = "network-config.service"
	pullSecretPath           = "/var/tmp/pull-secret.json" //nolint:gosec // not a secret, just a filename

	trustedBundlePath        = "/etc/pki/ca-trust/source/anchors/additional-trust-bundle.pem"
	nmstateConfigPath        = "/var/tmp/network-config.yaml"
	registriesConfPath       = "/etc/containers/registries.conf"
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
		&InstallationService{},
		&RegistriesConf{},
		&PostDeployment{},
	}
}

type ibiConfigurationFile struct {
	ExtraPartitionLabel  string `json:"extraPartitionLabel,omitempty"`
	ExtraPartitionNumber uint   `json:"extraPartitionNumber,omitempty"`
	ExtraPartitionStart  string `json:"extraPartitionStart,omitempty"`
	InstallationDisk     string `json:"installationDisk"`
	SeedImage            string `json:"seedImage"`
	SeedVersion          string `json:"seedVersion"`
	Shutdown             bool   `json:"shutdown,omitempty"`
	SkipDiskCleanup      bool   `json:"skipDiskCleanup,omitempty"`
}

// Generate generates the image-based installer ignition.
func (i *Ignition) Generate(dependencies asset.Parents) error {
	configAsset := &ImageBasedInstallationConfig{}
	installationServiceAsset := &InstallationService{}
	registriesConf := &RegistriesConf{}
	postDeployment := &PostDeployment{}

	dependencies.Get(configAsset, installationServiceAsset, registriesConf, postDeployment)

	ibiConfig := configAsset.Config
	if ibiConfig == nil {
		return errors.New("imagebased-installation-config.yaml is required")
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

	setFileInIgnition(config, pullSecretPath, ibiConfig.PullSecret, 0o600)

	ibiConfigFile := ibiConfigurationFile{
		ExtraPartitionLabel:  ibiConfig.ExtraPartitionLabel,
		ExtraPartitionNumber: ibiConfig.ExtraPartitionNumber,
		ExtraPartitionStart:  ibiConfig.ExtraPartitionStart,
		InstallationDisk:     ibiConfig.InstallationDisk,
		SeedVersion:          ibiConfig.SeedVersion,
		SeedImage:            ibiConfig.SeedImage,
		Shutdown:             ibiConfig.Shutdown,
		SkipDiskCleanup:      ibiConfig.SkipDiskCleanup,
	}
	marshaled, err := json.Marshal(ibiConfigFile)
	if err != nil {
		return fmt.Errorf("failed to marshall image-based installation configuration data: %w", err)
	}
	setFileInIgnition(config, ibiConfigurationPath, string(marshaled), 0o600)

	if len(registriesConf.Data) > 0 {
		setFileInIgnition(config, registriesConfPath, string(registriesConf.Data), 0o600)
	}

	setFileInIgnition(config, installationScriptPath, installationScript, 0o755)
	setUnitInIgnition(config, installationServiceName, installationServiceAsset.Content)

	if ibiConfig.AdditionalTrustBundle != "" {
		setFileInIgnition(config, trustedBundlePath, ibiConfig.AdditionalTrustBundle, 0o600)
	}

	if ibiConfig.NetworkConfig.String() != "" {
		setFileInIgnition(config, nmstateConfigPath, ibiConfig.NetworkConfig.String(), 0o600)
		setUnitInIgnition(config, networkConfigServiceName, networkConfigService)
	}

	if ibiConfig.IgnitionConfigOverride != "" {
		ignitionConfigOverride, _, err := v3_2.Parse([]byte(ibiConfig.IgnitionConfigOverride))
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
	}

	if postDeployment.File != nil {
		setFileInIgnition(config, postDeploymentScriptPath, string(postDeployment.File.Data), 0o755)
	}

	i.Config = config
	return nil
}

func setFileInIgnition(config *igntypes.Config, filePath, fileContents string, mode int) {
	fileContentsEncoded := "data:text/plain;charset=utf-8;base64," + base64.StdEncoding.EncodeToString([]byte(fileContents))
	file := igntypes.File{
		Node: igntypes.Node{
			Path:      filePath,
			Overwrite: ignutil.BoolToPtr(true),
			Group:     igntypes.NodeGroup{},
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Append: []igntypes.Resource{},
			Contents: igntypes.Resource{
				Source: &fileContentsEncoded,
			},
			Mode: &mode,
		},
	}
	config.Storage.Files = append(config.Storage.Files, file)
}

func setUnitInIgnition(config *igntypes.Config, name, contents string) {
	newUnit := igntypes.Unit{
		Contents: &contents,
		Name:     name,
		Enabled:  ignutil.BoolToPtr(true),
	}
	config.Systemd.Units = append(config.Systemd.Units, newUnit)
}
