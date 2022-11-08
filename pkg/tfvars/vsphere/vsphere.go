package vsphere

import (
	"encoding/json"
	"fmt"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	vtypes "github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pkg/errors"
	"strings"
)

type folder struct {
	Name       string `json:"name"`
	Datacenter string `json:"vsphere_datacenter"`
}

type config struct {
	VSphereURL        string          `json:"vsphere_url"`
	VSphereUsername   string          `json:"vsphere_username"`
	VSpherePassword   string          `json:"vsphere_password"`
	MemoryMiB         int64           `json:"vsphere_control_plane_memory_mib"`
	DiskGiB           int32           `json:"vsphere_control_plane_disk_gib"`
	NumCPUs           int32           `json:"vsphere_control_plane_num_cpus"`
	NumCoresPerSocket int32           `json:"vsphere_control_plane_cores_per_socket"`
	Cluster           string          `json:"vsphere_cluster"`
	ResourcePool      string          `json:"vsphere_resource_pool"`
	Datacenter        string          `json:"vsphere_datacenter"`
	Datastore         string          `json:"vsphere_datastore"`
	Folder            string          `json:"vsphere_folder"`
	Network           string          `json:"vsphere_network"`
	Template          string          `json:"vsphere_template"`
	OvaFilePath       string          `json:"vsphere_ova_filepath"`
	PreexistingFolder bool            `json:"vsphere_preexisting_folder"`
	DiskType          vtypes.DiskType `json:"vsphere_disk_type"`

	// vcenters can still remain a map for easy lookups
	VCenters       map[string]vtypes.VCenter `json:"vsphere_vcenters"`
	FailureDomains []vtypes.FailureDomain    `json:"vsphere_failure_domains"`

	NetworksInFailureDomains map[string]string `json:"vsphere_networks"`

	ControlPlanes []*machineapi.VSphereMachineProviderSpec `json:"vsphere_control_planes"`

	DatacentersFolders map[string]*folder `json:"vsphere_folders"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	ControlPlaneConfigs []*machineapi.VSphereMachineProviderSpec
	Username            string
	Password            string
	Cluster             string
	ImageURL            string
	PreexistingFolder   bool
	DiskType            vtypes.DiskType
	NetworkID           string

	NetworksInFailureDomain map[string]string
	InstallConfig           *installconfig.InstallConfig
	InfraID                 string

	ControlPlaneMachines []machineapi.Machine
}

// TFVars generate vSphere-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	controlPlaneConfig := sources.ControlPlaneConfigs[0]
	cachedImage, err := cache.DownloadImageFile(sources.ImageURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached vsphere image")
	}

	// The vSphere provider needs the relativepath of the folder,
	// so get the relPath from the absolute path. Absolute path is always of the form
	// /<datacenter>/vm/<folder_path> so we can split on "vm/".

	folderPathList := strings.SplitAfterN(controlPlaneConfig.Workspace.Folder, "vm/", 2)

	// This should never happen
	if len(folderPathList) <= 1 {
		return nil, errors.Errorf("control plane folder is not defined as a path %s", controlPlaneConfig.Workspace.Folder)
	}
	folderRelPath := folderPathList[1]

	vcenterZones := convertVCentersToMap(sources.InstallConfig.Config.VSphere.VCenters)
	datacentersFolders := createDatacenterFolderMap(sources.InfraID, sources.InstallConfig.Config.VSphere.FailureDomains)

	cfg := &config{
		VSphereURL:        controlPlaneConfig.Workspace.Server,
		VSphereUsername:   sources.Username,
		VSpherePassword:   sources.Password,
		MemoryMiB:         controlPlaneConfig.MemoryMiB,
		DiskGiB:           controlPlaneConfig.DiskGiB,
		NumCPUs:           controlPlaneConfig.NumCPUs,
		NumCoresPerSocket: controlPlaneConfig.NumCoresPerSocket,
		Cluster:           sources.Cluster,
		ResourcePool:      controlPlaneConfig.Workspace.ResourcePool,
		Datacenter:        controlPlaneConfig.Workspace.Datacenter,
		Datastore:         controlPlaneConfig.Workspace.Datastore,
		Folder:            folderRelPath,
		Network:           sources.NetworkID,
		Template:          controlPlaneConfig.Template,
		OvaFilePath:       cachedImage,
		PreexistingFolder: sources.PreexistingFolder,
		DiskType:          sources.DiskType,

		VCenters:                 vcenterZones,
		FailureDomains:           sources.InstallConfig.Config.VSphere.FailureDomains,
		NetworksInFailureDomains: sources.NetworksInFailureDomain,
		ControlPlanes:            sources.ControlPlaneConfigs,
		DatacentersFolders:       datacentersFolders,
	}

	return json.MarshalIndent(cfg, "", "  ")
}

// createDatacenterFolderMap()
// This function loops over the range of failure domains
// Each failure domain defines the vCenter datacenter and folder
// to be used for the virtual machines within that domain.
// The datacenter could be reused but a folder could be
// unique - the key then becomes a string that contains
// both the datacenter name and the folder to be created.

func createDatacenterFolderMap(infraID string, failureDomains []vtypes.FailureDomain) map[string]*folder {
	folders := make(map[string]*folder)

	for i, fd := range failureDomains {

		tempFolder := new(folder)

		tempFolder.Datacenter = fd.Topology.Datacenter
		tempFolder.Name = fd.Topology.Folder

		if tempFolder.Name == "" {
			tempFolder.Name = infraID
			failureDomains[i].Topology.Folder = infraID
			key := fmt.Sprintf("%s-%s", tempFolder.Datacenter, tempFolder.Name)
			folders[key] = tempFolder
		}
	}
	return folders
}

func convertVCentersToMap(values []vtypes.VCenter) map[string]vtypes.VCenter {
	vcenterMap := make(map[string]vtypes.VCenter)
	for _, v := range values {
		vcenterMap[v.Server] = v
	}
	return vcenterMap
}
