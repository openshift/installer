package image

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/coreos/ignition/v2/config/merge"
	"github.com/coreos/ignition/v2/config/v3_2"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/imagebased/configimage"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/imagebased"
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
		&configimage.InstallConfig{},
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
	installConfig := &configimage.InstallConfig{}

	dependencies.Get(configAsset, registriesConf, postDeployment, installConfig)

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

	// Prepare CoreOS installer arguments, adding dual-stack support if needed
	coreosInstallerArgs := ibiConfig.CoreosInstallerArgs
	if dhcpArgs := getDHCPKernelArgs(ibiConfig, installConfig); dhcpArgs != "" {
		// Add specific DHCP kernel arguments based on which IP versions need DHCP
		coreosInstallerArgs = append(coreosInstallerArgs, "--append-karg", dhcpArgs)
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
		CoreosInstallerArgs:  coreosInstallerArgs,
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

// getDHCPKernelArgs determines which DHCP kernel arguments should be added.
// Returns "ip=dhcp", "ip=dhcp6", "ip=dhcp,dhcp6", or "" based on which IP versions need DHCP.
// Only adds args for dual-stack scenarios and when user hasn't already provided them.
func getDHCPKernelArgs(ibiConfig *imagebased.InstallationConfig, installConfig *configimage.InstallConfig) string {
	// Check if we have machine networks configured
	if installConfig.Config == nil || installConfig.Config.Networking == nil ||
		len(installConfig.Config.Networking.MachineNetwork) == 0 {
		return ""
	}

	machineNetworks := installConfig.Config.Networking.MachineNetwork

	if len(machineNetworks) != 2 {
		return ""
	}

	hasIPv4Network := false
	hasIPv6Network := false

	for _, network := range machineNetworks {
		if network.CIDR.IP.To4() != nil {
			hasIPv4Network = true
		} else {
			hasIPv6Network = true
		}
	}

	// Must be true dual-stack (both IPv4 and IPv6)
	if !hasIPv4Network || !hasIPv6Network {
		return ""
	}

	if userAlreadyProvidedKernelArgs(ibiConfig.CoreosInstallerArgs) {
		return ""
	}

	if ibiConfig.NetworkConfig == nil || ibiConfig.NetworkConfig.String() == "" {
		return "ip=dhcp,dhcp6" // Dual-stack DHCP
	}

	networkConfigStr := ibiConfig.NetworkConfig.String()
	ipv4NeedsDHCP := !isIPVersionStatic(networkConfigStr, "ipv4")
	ipv6NeedsDHCP := !isIPVersionStatic(networkConfigStr, "ipv6")

	switch {
	case ipv4NeedsDHCP && ipv6NeedsDHCP:
		return "ip=dhcp,dhcp6"
	case ipv6NeedsDHCP:
		return "ip=dhcp6"
	case ipv4NeedsDHCP:
		return "ip=dhcp"
	}

	return ""
}

// userAlreadyProvidedKernelArgs checks if the user has already provided IP kernel arguments.
func userAlreadyProvidedKernelArgs(args []string) bool {
	for i, arg := range args {
		// Check if this is --append-karg followed by an ip= argument
		if arg == "--append-karg" && i+1 < len(args) {
			nextArg := args[i+1]
			if strings.HasPrefix(nextArg, "ip=") {
				return true
			}
		}
		// Also check if someone provided ip= directly (though unusual)
		if strings.HasPrefix(arg, "ip=") {
			return true
		}
	}
	return false
}

// isIPVersionStatic checks if a specific IP version is configured for static networking.
func isIPVersionStatic(networkConfigStr, ipVersion string) bool {
	// Look for the IP version section (ipv4: or ipv6:)
	if !strings.Contains(networkConfigStr, ipVersion+":") {
		return false // Not configured = not static
	}

	// Check if this IP version has static addresses and DHCP disabled
	lines := strings.Split(networkConfigStr, "\n")
	inIPSection := false
	hasAddress := false
	dhcpDisabled := false
	currentIndent := 0

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		if strings.Contains(trimmedLine, ipVersion+":") {
			inIPSection = true
			currentIndent = indent
			continue
		}

		if inIPSection && indent <= currentIndent && strings.Contains(trimmedLine, ":") &&
			!strings.HasPrefix(trimmedLine, ipVersion) {
			break
		}

		if inIPSection {
			if strings.Contains(trimmedLine, "address:") || strings.Contains(trimmedLine, "- ip:") {
				hasAddress = true
			}
			if strings.Contains(trimmedLine, "dhcp:") {
				if strings.Contains(trimmedLine, "dhcp: false") {
					dhcpDisabled = true
				} else if strings.Contains(trimmedLine, "dhcp: true") {
					return false
				}
			}
		}
	}

	return hasAddress && dhcpDisabled
}
