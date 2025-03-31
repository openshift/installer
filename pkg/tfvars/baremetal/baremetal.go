// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"

	"github.com/openshift/installer/pkg/types"
)

type Bridge struct {
	Name string `json:"name"`
	MAC  string `json:"mac"`
}

// Config represents the baremetal platform parts of install config needed for bootstrapping.
type Config struct {
	LibvirtURI  string    `json:"libvirt_uri,omitempty"`
	Bridges     []Bridge  `json:"bridges"`
}

type imageDownloadFunc func(baseURL, applicationName string) (string, error)

var (
	imageDownloader imageDownloadFunc
)

func init() {
	imageDownloader = cache.DownloadImageFile
}

// TFVars generates bare metal specific Terraform variables.
func TFVars(libvirtURI string, bootstrapOSImage, externalBridge, externalMAC, provisioningBridge, provisioningMAC string) ([]byte, error) {
	bootstrapOSImage, err := imageDownloader(bootstrapOSImage, cache.InstallerApplicationName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached bootstrap libvirt image")
	}

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
		LibvirtURI:       libvirtURI,
		BootstrapOSImage: bootstrapOSImage,
		Bridges:          bridges,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
