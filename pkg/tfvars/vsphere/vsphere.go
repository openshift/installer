package vsphere

import (
	"context"
	"encoding/json"
	"strings"

	vsphereapis "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider/v1beta1"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/openshift/installer/pkg/types/vsphere"
)

type config struct {
	VSphereURL        string           `json:"vsphere_url"`
	VSphereUsername   string           `json:"vsphere_username"`
	VSpherePassword   string           `json:"vsphere_password"`
	MemoryMiB         int64            `json:"vsphere_control_plane_memory_mib"`
	DiskGiB           int32            `json:"vsphere_control_plane_disk_gib"`
	NumCPUs           int32            `json:"vsphere_control_plane_num_cpus"`
	NumCoresPerSocket int32            `json:"vsphere_control_plane_cores_per_socket"`
	Cluster           string           `json:"vsphere_cluster"`
	ResourcePool      string           `json:"vsphere_resource_pool"`
	Datacenter        string           `json:"vsphere_datacenter"`
	Datastore         string           `json:"vsphere_datastore"`
	Folder            string           `json:"vsphere_folder"`
	Network           string           `json:"vsphere_network"`
	Template          string           `json:"vsphere_template"`
	OvaFilePath       string           `json:"vsphere_ova_filepath"`
	PreexistingFolder bool             `json:"vsphere_preexisting_folder"`
	DiskType          vsphere.DiskType `json:"vsphere_disk_type"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	ControlPlaneConfigs []*vsphereapis.VSphereMachineProviderSpec
	Username            string
	Password            string
	Cluster             string
	ImageURL            string
	PreexistingFolder   bool
	DiskType            vsphere.DiskType
}

//TFVars generate vSphere-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	controlPlaneConfig := sources.ControlPlaneConfigs[0]

	cachedImage, err := cache.DownloadImageFile(sources.ImageURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached vsphere image")
	}

	// The vSphere provider needs the relativepath of the folder,
	// so get the relPath from the absolute path. Absolute path is always of the form
	// /<datacenter>/vm/<folder_path> so we can split on "vm/".
	folderRelPath := strings.SplitAfterN(controlPlaneConfig.Workspace.Folder, "vm/", 2)[1]

	// Must use the Managed Object ID for a port group (e.g. dvportgroup-5258)
	// instead of the name since port group names aren't always unique in vSphere.
	// https://bugzilla.redhat.com/show_bug.cgi?id=1918005
	vim25Client, _, err := vsphere.CreateVSphereClients(context.TODO(),
		controlPlaneConfig.Workspace.Server,
		sources.Username,
		sources.Password)
	finder := vsphere.NewFinder(vim25Client)
	networkUtil := vsphere.NewNetworkUtil(vim25Client)
	networkID, err := vsphere.GetNetworkMoID(context.TODO(),
		networkUtil,
		finder,
		controlPlaneConfig.Workspace.Datacenter,
		sources.Cluster,
		controlPlaneConfig.Network.Devices[0].NetworkName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get vSphere network ID")
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
		ResourcePool:      controlPlaneConfig.Workspace.ResourcePool,
		Datacenter:        controlPlaneConfig.Workspace.Datacenter,
		Datastore:         controlPlaneConfig.Workspace.Datastore,
		Folder:            folderRelPath,
		Network:           networkID,
		Template:          controlPlaneConfig.Template,
		OvaFilePath:       cachedImage,
		PreexistingFolder: sources.PreexistingFolder,
		DiskType:          sources.DiskType,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
