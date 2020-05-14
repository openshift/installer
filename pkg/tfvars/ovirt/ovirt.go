// Package ovirt contains ovirt-specific Terraform-variable logic.
package ovirt

import (
	"encoding/json"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
)

type config struct {
	URL                    string `json:"ovirt_url"`
	Username               string `json:"ovirt_username"`
	Password               string `json:"ovirt_password"`
	Cafile                 string `json:"ovirt_cafile,omitempty"`
	ClusterID              string `json:"ovirt_cluster_id"`
	StorageDomainID        string `json:"ovirt_storage_domain_id"`
	NetworkName            string `json:"ovirt_network_name,omitempty"`
	VNICProfileID          string `json:"ovirt_vnic_profile_id,omitempty"`
	BaseImageName          string `json:"openstack_base_image_name,omitempty"`
	BaseImageLocalFilePath string `json:"openstack_base_image_local_file_path,omitempty"`
}

// TFVars generates ovirt-specific Terraform variables.
func TFVars(
	engineURL string,
	engineUser string,
	enginePass string,
	engineCafile string,
	clusterID string,
	stoarageDomainID string,
	networkName string,
	vnicProfileID string,
	baseImage string,
	infraID string) ([]byte, error) {

	cfg := config{
		URL:             engineURL,
		Username:        engineUser,
		Password:        enginePass,
		Cafile:          engineCafile,
		ClusterID:       clusterID,
		StorageDomainID: stoarageDomainID,
		NetworkName:     networkName,
		VNICProfileID:   vnicProfileID,
		BaseImageName:   baseImage,
	}

	imageName, isURL := rhcos.GenerateOpenStackImageName(baseImage, infraID)
	cfg.BaseImageName = imageName
	if isURL {
		imageFilePath, err := cache.DownloadImageFile(baseImage)
		if err != nil {
			return nil, err
		}
		cfg.BaseImageLocalFilePath = imageFilePath
	}

	return json.MarshalIndent(cfg, "", "  ")
}
