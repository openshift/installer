package vsphere

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	machineapi "github.com/openshift/api/machine/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	vtypes "github.com/openshift/installer/pkg/types/vsphere"
)

type folder struct {
	Name       string `json:"name"`
	Datacenter string `json:"vsphere_datacenter"`
}

type controlplane struct {
	DeploymentZone     string                                 `json:"name"`
	ControlPlaneConfig *machineapi.VSphereMachineProviderSpec `json:"vsphere_control_plane"`
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

	VCenters          map[string]vtypes.VCenter         `json:"vsphere_vcenters"`
	DeploymentZone    map[string]*vtypes.DeploymentZone `json:"vsphere_deployment_zone"`
	FailureDomainZone map[string]vtypes.FailureDomain   `json:"vsphere_failure_zone"`
	NetworkZone       map[string]string                 `json:"vsphere_network_zone"`

	FolderZone map[string]*folder `json:"vsphere_folder_zone"`

	ControlPlanes []controlplane `json:"vsphere_control_planes"`

	ControlPlaneConfigZone map[string]*machineapi.VSphereMachineProviderSpec `json:"vsphere_control_planes_zone"`
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

	NetworkZone   map[string]string
	InstallConfig *installconfig.InstallConfig
	InfraID       string

	ControlPlaneMachines []machineapi.Machine
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

	deploymentZones := convertDeploymentZonesToMap(sources.InstallConfig.Config.VSphere.DeploymentZones)
	vcenterZones := convertVCentersToMap(sources.InstallConfig.Config.VSphere.VCenters)
	failureDomainZones := convertFailureZoneToMap(sources.InstallConfig.Config.VSphere.FailureDomains)
	controlPlaneConfigZone, controlPlanes := convertControlPlaneToMap(sources.ControlPlaneMachines, sources.InstallConfig)
	folderZone := createFolderZoneMap(sources.InfraID, deploymentZones, failureDomainZones)

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

		NetworkZone:            sources.NetworkZone,
		FolderZone:             folderZone,
		VCenters:               vcenterZones,
		DeploymentZone:         deploymentZones,
		FailureDomainZone:      failureDomainZones,
		ControlPlaneConfigZone: controlPlaneConfigZone,

		ControlPlanes: controlPlanes,
	}

	return json.MarshalIndent(cfg, "", "  ")
}

func createFolderZoneMap(infraID string, deploymentZone map[string]*vtypes.DeploymentZone, failureDomainZones map[string]vtypes.FailureDomain) map[string]*folder {
	folders := make(map[string]*folder)

	for k, v := range deploymentZone {
		tempFolder := new(folder)

		tempFolder.Datacenter = failureDomainZones[v.FailureDomain].Topology.Datacenter
		tempFolder.Name = v.PlacementConstraint.Folder

		if tempFolder.Name == "" {
			tempFolder.Name = infraID
			deploymentZone[k].PlacementConstraint.Folder = infraID
		}

		key := fmt.Sprintf("%s-%s", tempFolder.Datacenter, tempFolder.Name)
		folders[key] = tempFolder
	}
	return folders
}

func convertDeploymentZonesToMap(values []vtypes.DeploymentZone) map[string]*vtypes.DeploymentZone {
	deploymentZoneMap := make(map[string]*vtypes.DeploymentZone)

	for i := range values {
		tempValue := &values[i]
		deploymentZoneMap[tempValue.Name] = tempValue
	}
	return deploymentZoneMap
}

func convertControlPlaneToMap(values []machineapi.Machine, installConfig *installconfig.InstallConfig) (map[string]*machineapi.VSphereMachineProviderSpec, []controlplane) {
	controlPlaneZonalConfigs := make(map[string]*machineapi.VSphereMachineProviderSpec)
	var controlPlaneConfigs []controlplane

	var region string
	var zone string
	var deploymentZone vtypes.DeploymentZone
	var failureDomainName string

	for i, v := range values {

		// We need to know the region and zone of each master virtual machine
		// to determine the name of the DeploymentZone
		if val, ok := v.ObjectMeta.Labels["machine.openshift.io/region"]; ok {
			region = val
		}
		if val, ok := v.ObjectMeta.Labels["machine.openshift.io/zone"]; ok {
			zone = val
		}
		// Using failuredomains, region and zone names find the failureDomain
		for _, fd := range installConfig.Config.VSphere.FailureDomains {
			if fd.Region.Name == region {
				if fd.Zone.Name == zone {
					failureDomainName = fd.Name
				}
			}
		}

		// Using the deploymentzones and failuredomainname, find the deploymentzone
		for _, dz := range installConfig.Config.VSphere.DeploymentZones {
			if dz.FailureDomain == failureDomainName {
				deploymentZone = dz
			}
		}
		controlPlaneZonalConfigs[deploymentZone.Name] = values[i].Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec)

		controlPlaneConfigs[i] = controlplane{
			ControlPlaneConfig: values[i].Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec),
			DeploymentZone:     deploymentZone.Name,
		}
	}

	return controlPlaneZonalConfigs, controlPlaneConfigs
}

func convertFailureZoneToMap(values []vtypes.FailureDomain) map[string]vtypes.FailureDomain {
	failureDomainMap := make(map[string]vtypes.FailureDomain)
	for _, v := range values {
		failureDomainMap[v.Name] = v
	}
	return failureDomainMap
}

func convertVCentersToMap(values []vtypes.VCenter) map[string]vtypes.VCenter {
	vcenterMap := make(map[string]vtypes.VCenter)
	for _, v := range values {
		vcenterMap[v.Server] = v
	}
	return vcenterMap
}
