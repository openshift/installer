package nutanix

import (
	"encoding/json"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	"github.com/pkg/errors"
)

type config struct {
	PrismCentralAddress            string `json:"nutanix_prism_central_address"`
	Port                           string `json:"nutanix_prism_central_port"`
	Username                       string `json:"nutanix_username"`
	Password                       string `json:"nutanix_password"`
	MemoryMiB                      int64  `json:"nutanix_control_plane_memory_mib"`
	DiskSizeMiB                    int64  `json:"nutanix_control_plane_disk_mib"`
	NumCPUs                        int64  `json:"nutanix_control_plane_num_cpus"`
	NumCoresPerSocket              int64  `json:"nutanix_control_plane_cores_per_socket"`
	PrismElementUUID               string `json:"nutanix_prism_element_uuid"`
	SubnetUUID                     string `json:"nutanix_subnet_uuid"`
	Image                          string `json:"nutanix_image"`
	ImageFilePath                  string `json:"nutanix_image_filepath"`
	BootstrapIgnitionImage         string `json:"nutanix_bootstrap_ignition_image"`
	BootstrapIgnitionImageFilePath string `json:"nutanix_bootstrap_ignition_image_filepath"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	PrismCentralAddress   string
	Port                  string
	Username              string
	Password              string
	ImageURL              string
	BootstrapIgnitionData string
	ClusterID             string
	ControlPlaneConfigs   []*machinev1.NutanixMachineProviderConfig
}

//TFVars generate Nutanix-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	cachedImage, err := cache.DownloadImageFile(sources.ImageURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached nutanix image")
	}

	bootstrapIgnitionImagePath, err := nutanixtypes.CreateBootstrapISO(sources.ClusterID, sources.BootstrapIgnitionData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bootstrap ignition iso")
	}

	bootstrapIgnitionImageName := nutanixtypes.BootISOImageName(sources.ClusterID)
	controlPlaneConfig := sources.ControlPlaneConfigs[0]
	cfg := &config{
		Port:                           sources.Port,
		PrismCentralAddress:            sources.PrismCentralAddress,
		Username:                       sources.Username,
		Password:                       sources.Password,
		MemoryMiB:                      controlPlaneConfig.MemorySize.Value() / (1024 * 1024),
		DiskSizeMiB:                    controlPlaneConfig.SystemDiskSize.Value() / (1024 * 1024),
		NumCPUs:                        int64(controlPlaneConfig.VCPUSockets),
		NumCoresPerSocket:              int64(controlPlaneConfig.VCPUsPerSocket),
		PrismElementUUID:               *controlPlaneConfig.Cluster.UUID,
		SubnetUUID:                     *controlPlaneConfig.Subnets[0].UUID,
		Image:                          *controlPlaneConfig.Image.Name,
		ImageFilePath:                  cachedImage,
		BootstrapIgnitionImage:         bootstrapIgnitionImageName,
		BootstrapIgnitionImageFilePath: bootstrapIgnitionImagePath,
	}
	return json.MarshalIndent(cfg, "", "  ")
}
