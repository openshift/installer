// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/rhcos/cache"
)

// Config represents the baremetal platform parts of install config needed for bootstrapping.
type Config struct {
	LibvirtURI       string              `json:"libvirt_uri,omitempty"`
	BootstrapOSImage string              `json:"bootstrap_os_image,omitempty"`
	Bridges          []map[string]string `json:"bridges"`
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

	var bridges []map[string]string

	bridges = append(bridges,
		map[string]string{
			"name": externalBridge,
			"mac":  externalMAC,
		})

	if provisioningBridge != "" {
		bridges = append(
			bridges,
			map[string]string{
				"name": provisioningBridge,
				"mac":  provisioningMAC,
			})
	}

	cfg := &Config{
		LibvirtURI:       libvirtURI,
		BootstrapOSImage: bootstrapOSImage,
		Bridges:          bridges,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
