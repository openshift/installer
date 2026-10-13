package nutanix

import (
	"encoding/json"

	"github.com/pkg/errors"

	machinev1 "github.com/openshift/api/machine/v1"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

type config struct {
	PrismCentralAddress            string            `json:"nutanix_prism_central_address"`
	Port                           string            `json:"nutanix_prism_central_port"`
	Username                       string            `json:"nutanix_username"`
	Password                       string            `json:"nutanix_password"`
	MemoryMiB                      int64             `json:"nutanix_control_plane_memory_mib"`
	DiskSizeMiB                    int64             `json:"nutanix_control_plane_disk_mib"`
	NumCPUs                        int64             `json:"nutanix_control_plane_num_cpus"`
	NumCoresPerSocket              int64             `json:"nutanix_control_plane_cores_per_socket"`
	ProjectUUID                    string            `json:"nutanix_control_plane_project_uuid"`
	Categories                     map[string]string `json:"nutanix_control_plane_categories"`
	PrismElementUUIDs              []string          `json:"nutanix_prism_element_uuids"`
	SubnetUUIDs                    []string          `json:"nutanix_subnet_uuids"`
	Image                          string            `json:"nutanix_image"`
	ImageURI                       string            `json:"nutanix_image_uri"`
	BootstrapIgnitionImage         string            `json:"nutanix_bootstrap_ignition_image"`
	BootstrapIgnitionImageFilePath string            `json:"nutanix_bootstrap_ignition_image_filepath"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	PrismCentralAddress   string
	Port                  string
	Username              string
	Password              string
	ImageURI              string
	BootstrapIgnitionData string
	ClusterID             string
	ControlPlaneConfigs   []*machinev1.NutanixMachineProviderConfig
}

// TFVars generate Nutanix-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	bootstrapIgnitionImagePath, err := nutanixtypes.CreateBootstrapISO(sources.ClusterID, sources.BootstrapIgnitionData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bootstrap ignition iso")
	}

	bootstrapIgnitionImageName := nutanixtypes.BootISOImageName(sources.ClusterID)
	controlPlaneConfig := sources.ControlPlaneConfigs[0]
	cpCount := len(sources.ControlPlaneConfigs)
	cfg := &config{
		Port:                           sources.Port,
		PrismCentralAddress:            sources.PrismCentralAddress,
		Username:                       sources.Username,
		Password:                       sources.Password,
		MemoryMiB:                      controlPlaneConfig.MemorySize.Value() / (1024 * 1024),
		DiskSizeMiB:                    controlPlaneConfig.SystemDiskSize.Value() / (1024 * 1024),
		NumCPUs:                        int64(controlPlaneConfig.VCPUSockets),
		NumCoresPerSocket:              int64(controlPlaneConfig.VCPUsPerSocket),
		PrismElementUUIDs:              make([]string, cpCount),
		SubnetUUIDs:                    make([]string, cpCount),
		Image:                          *controlPlaneConfig.Image.Name,
		ImageURI:                       sources.ImageURI,
		BootstrapIgnitionImage:         bootstrapIgnitionImageName,
		BootstrapIgnitionImageFilePath: bootstrapIgnitionImagePath,
	}

	for i, cpcfg := range sources.ControlPlaneConfigs {
		cfg.PrismElementUUIDs[i] = *cpcfg.Cluster.UUID
		cfg.SubnetUUIDs[i] = *cpcfg.Subnets[0].UUID
	}

	if controlPlaneConfig.Project.Type == machinev1.NutanixIdentifierUUID {
		cfg.ProjectUUID = *controlPlaneConfig.Project.UUID
	}
	cfg.Categories = make(map[string]string, len(controlPlaneConfig.Categories))
	for _, category := range controlPlaneConfig.Categories {
		cfg.Categories[category.Key] = category.Value
	}

	return json.MarshalIndent(cfg, "", "  ")
}
