package ibiconfig

import (
	"fmt"
)

// ImageBasedInstallConfigVersion is the version supported by this package.
const ImageBasedInstallConfigVersion = "v1beta1"

const (
	defaultExtraPartitionLabel = "varlibcontainers"
)

// IBIPrepareConfig or aka ImageBasedInstallConfig is the API for specifying configuration
// for the image-based installer.
type IBIPrepareConfig struct {
	SSHPublicKeyFile string `json:"sshPublicKeyFile,omitempty"`
	RHCOSLiveISO     string `json:"rhcosLiveIso,omitempty"`

	SeedImage            string `json:"seedImage"`
	SeedVersion          string `json:"seedVersion"`
	AuthFile             string `json:"authFile"`
	PullSecretFile       string `json:"pullSecretFile,omitempty"`
	InstallationDisk     string `json:"installationDisk"`
	PrecacheBestEffort   bool   `json:"precacheBestEffort,omitempty"`
	PrecacheDisabled     bool   `json:"precacheDisabled,omitempty"`
	Shutdown             bool   `json:"shutdown,omitempty"`
	UseContainersFolder  bool   `json:"useContainersFolder,omitempty"`
	ExtraPartitionStart  string `json:"extraPartitionStart,omitempty"`
	ExtraPartitionLabel  string `json:"extraPartitionLabel,omitempty"`
	ExtraPartitionNumber uint   `json:"extraPartitionNumber,omitempty"`
	SkipDiskCleanup      bool   `json:"skipDiskCleanup,omitempty"`
}

func (c *IBIPrepareConfig) Validate() error {
	if c.AuthFile == "" {
		return fmt.Errorf("authFile is required")
	}
	if c.SeedImage == "" {
		return fmt.Errorf("seedImage is required")
	}
	if c.SeedVersion == "" {
		return fmt.Errorf("seedVersion is required")
	}
	if c.InstallationDisk == "" {
		return fmt.Errorf("installationDisk is required")
	}
	return nil
}

func (c *IBIPrepareConfig) SetDefaultValues() {

	if c.PullSecretFile == "" {
		c.PullSecretFile = c.AuthFile
	}
	if c.ExtraPartitionStart == "" {
		c.ExtraPartitionStart = "-40G"
	}
	if c.ExtraPartitionNumber == 0 {
		c.ExtraPartitionNumber = 5
	}
	if c.ExtraPartitionLabel == "" {
		c.ExtraPartitionLabel = defaultExtraPartitionLabel
	}
}
