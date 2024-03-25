package ibi

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ImageBasedInstallConfigVersion is the version supported by this package.
const ImageBasedInstallConfigVersion = "v1beta1"

// Config or aka ImageBasedInstallConfig is the API for specifying configuration
// for the image-based installer.
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	SeedImage           string `json:"seedImage,omitempty"`
	SeedVersion         string `json:"seedVersion,omitempty"`
	AuthFile            string `json:"authFile,omitempty"`
	PullSecretFile      string `json:"pullSecretFile,omitempty"`
	SSHPublicKeyFile    string `json:"sshPublicKeyFile,omitempty"`
	LCAImage            string `json:"lcaImage,omitempty"`
	RHCOSLiveISO        string `json:"rhcosLiveIso,omitempty"`
	InstallationDisk    string `json:"installationDisk,omitempty"`
	ExtraPartitionStart string `json:"extraPartitionStart,omitempty"`
	PrecacheBestEffort  bool   `json:"precacheBestEffort,omitempty"`
	PrecacheDisabled    bool   `json:"precacheDisabled,omitempty"`
}
