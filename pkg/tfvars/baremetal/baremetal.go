// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"

	"github.com/openshift/installer/pkg/types"
)

// Bridge represents a network bridge on the provisioner host.
type Bridge struct {
	Name string `json:"name"`
	MAC  string `json:"mac"`
}

// Config represents the baremetal platform parts of install config needed for bootstrapping.
type Config struct {
	LibvirtURI           string             `json:"libvirt_uri,omitempty"`
	ReleaseImagePullSpec string             `json:"release_image,omitempty"`
	PullSecret           string             `json:"pull_secret,omitempty"`
	MirrorConfig         types.MirrorConfig `json:"mirror_config,omitempty"`
	Bridges              []Bridge           `json:"bridges"`
}

// TFVars generates bare metal specific Terraform variables.
func TFVars(
	libvirtURI string,
	releaseImagePullSpec string,
	pullSecret string,
	mirrorConfig types.MirrorConfig,
	externalBridge, externalMAC, provisioningBridge, provisioningMAC string) ([]byte, error) {
	var bridges []Bridge

	bridges = append(bridges,
		Bridge{
			Name: externalBridge,
			MAC:  externalMAC,
		})

	if provisioningBridge != "" {
		bridges = append(bridges,
			Bridge{
				Name: provisioningBridge,
				MAC:  provisioningMAC,
			})
	}

	cfg := &Config{
		LibvirtURI:           libvirtURI,
		ReleaseImagePullSpec: releaseImagePullSpec,
		PullSecret:           pullSecret,
		MirrorConfig:         mirrorConfig,
		Bridges:              bridges,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
