package vsphere

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	vsphereapis "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider/v1alpha1"
)

type config struct {
	VSphereURL        string `json:"vsphere_url"`
	VSphereUsername   string `json:"vsphere_username"`
	VSpherePassword   string `json:"vsphere_password"`
	MemoryMiB         int64  `json:"vsphere_control_plane_memory_mib"`
	DiskGiB           int32  `json:"vsphere_control_plane_disk_gib"`
	NumCPUs           int32  `json:"vsphere_control_plane_num_cpus"`
	NumCoresPerSocket int32  `json:"vsphere_control_plane_cores_per_socket"`
	Cluster           string `json:"vsphere_cluster"`
	Datacenter        string `json:"vsphere_datacenter"`
	Datastore         string `json:"vsphere_datastore"`
	Folder            string `json:"vsphere_folder"`
	Network           string `json:"vsphere_network"`
	Template          string `json:"vsphere_template"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	ControlPlaneConfigs []*vsphereapis.VSphereMachineProviderSpec
	Username            string
	Password            string
	Cluster             string
	ImageURI            string
}

//TFVars generate vSphere-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, string, error) {
	controlPlaneConfig := sources.ControlPlaneConfigs[0]

	cachedImage, err := cache.DownloadImageFile(sources.ImageURI)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to use cached vsphere image")
	}

	cfg := &config{
		VSphereURL:        controlPlaneConfig.Workspace.Server,
		VSphereUsername:   sources.Username,
		VSpherePassword:   sources.Password,
		MemoryMiB:         controlPlaneConfig.MemoryMiB,
		DiskGiB:           controlPlaneConfig.DiskGiB,
		NumCPUs:           controlPlaneConfig.NumCPUs,
		NumCoresPerSocket: controlPlaneConfig.NumCoresPerSocket,
		Cluster:           sources.Cluster,
		Datacenter:        controlPlaneConfig.Workspace.Datacenter,
		Datastore:         controlPlaneConfig.Workspace.Datastore,
		Folder:            controlPlaneConfig.Workspace.Folder,
		Network:           controlPlaneConfig.Network.Devices[0].NetworkName,
		Template:          controlPlaneConfig.Template,
	}

	cfgFormatted, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, "", err
	}

	return cfgFormatted, cachedImage, nil
}
