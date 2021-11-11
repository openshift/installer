package vsphere

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
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
	networkID, err := getNetworkMoID(controlPlaneConfig.Workspace.Server,
		sources.Username,
		sources.Password,
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

func getNetworkMoID(vcenter, username, password, datacenter, cluster, network string) (string, error) {
	client, _, err := vsphere.CreateVSphereClients(context.TODO(), vcenter, username, password)
	if err != nil {
		return "", errors.Wrap(err, "failed to setup vSphere client")
	}

	finder := find.NewFinder(client)
	path := fmt.Sprintf("/%s/host/*", datacenter)
	ccrs, err := finder.ClusterComputeResourceList(context.TODO(), path)
	if err != nil {
		return "", errors.Wrap(err, "could not list vSphere clusters")
	}

	var ccrmo mo.ClusterComputeResource
	for _, ccr := range ccrs {
		if ccr.Name() == cluster {
			err := ccr.Properties(context.TODO(), ccr.Reference(), []string{"network"}, &ccrmo)
			if err != nil {
				return "", errors.Wrap(err, "could not get properties of cluster")
			}
			for _, net := range ccrmo.Network {
				netObj := object.NewNetwork(client, net)
				name, err := netObj.ObjectName(context.TODO())
				if err != nil {
					return "", errors.Wrap(err, "could not get network name")
				}
				if name == network {
					return net.Value, nil
				}
			}
		}
	}

	return "", errors.Errorf("could not find a network %s in cluster %s", network, cluster)
}
