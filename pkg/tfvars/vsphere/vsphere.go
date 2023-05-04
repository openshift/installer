package vsphere

import (
	"encoding/json"
	"fmt"

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

type config struct {
	OvaFilePath              string                                   `json:"vsphere_ova_filepath"`
	DiskType                 vtypes.DiskType                          `json:"vsphere_disk_type"`
	VCenters                 map[string]vtypes.VCenter                `json:"vsphere_vcenters"`
	NetworksInFailureDomains map[string]string                        `json:"vsphere_networks"`
	ControlPlanes            []*machineapi.VSphereMachineProviderSpec `json:"vsphere_control_planes"`
	DatacentersFolders       map[string]*folder                       `json:"vsphere_folders"`

	ImportOvaFailureDomainMap map[string]vtypes.FailureDomain `json:"vsphere_import_ova_failure_domain_map"`
	FailureDomainMap          map[string]vtypes.FailureDomain `json:"vsphere_failure_domain_map"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	ControlPlaneConfigs     []*machineapi.VSphereMachineProviderSpec
	ImageURL                string
	DiskType                vtypes.DiskType
	NetworksInFailureDomain map[string]string
	InstallConfig           *installconfig.InstallConfig
	InfraID                 string
	ControlPlaneMachines    []machineapi.Machine
}

// TFVars generate vSphere-specific Terraform variables
func TFVars(sources TFVarsSources) ([]byte, error) {
	var err error
	cachedImage := ""

	failureDomainMap, importOvaFailureDomainMap := createFailureDomainMaps(sources.InstallConfig.Config.VSphere.FailureDomains, sources.InfraID)

	if len(importOvaFailureDomainMap) > 0 {
		cachedImage, err = cache.DownloadImageFile(sources.ImageURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to use cached vsphere image")
		}
	}

	vcenterZones := convertVCentersToMap(sources.InstallConfig.Config.VSphere.VCenters)
	datacentersFolders, err := createDatacenterFolderMap(sources.InfraID, sources.InstallConfig.Config.VSphere.FailureDomains)
	if err != nil {
		return nil, err
	}

	cfg := &config{
		OvaFilePath:              cachedImage,
		DiskType:                 sources.DiskType,
		VCenters:                 vcenterZones,
		NetworksInFailureDomains: sources.NetworksInFailureDomain,
		ControlPlanes:            sources.ControlPlaneConfigs,
		DatacentersFolders:       datacentersFolders,

		ImportOvaFailureDomainMap: importOvaFailureDomainMap,
		FailureDomainMap:          failureDomainMap,
	}

	return json.MarshalIndent(cfg, "", "  ")
}

func createFailureDomainMaps(failureDomains []vtypes.FailureDomain, infraID string) (map[string]vtypes.FailureDomain, map[string]vtypes.FailureDomain) {
	importOvaFailureDomainMap := make(map[string]vtypes.FailureDomain)
	failureDomainMap := make(map[string]vtypes.FailureDomain)

	for _, fd := range failureDomains {
		if fd.Topology.Folder == "" {
			fd.Topology.Folder = infraID
		}

		if fd.Topology.Template == "" {
			fd.Topology.Template = fmt.Sprintf("%s-rhcos-%s-%s", infraID, fd.Region, fd.Zone)
			importOvaFailureDomainMap[fd.Name] = fd
		}
		failureDomainMap[fd.Name] = fd
	}

	return failureDomainMap, importOvaFailureDomainMap
}

// createDatacenterFolderMap()
// This function loops over the range of failure domains
// Each failure domain defines the vCenter datacenter and folder
// to be used for the virtual machines within that domain.
// The datacenter could be reused but a folder could be
// unique - the key then becomes a string that contains
// both the datacenter name and the folder to be created.

func createDatacenterFolderMap(infraID string, failureDomains []vtypes.FailureDomain) (map[string]*folder, error) {
	folders := make(map[string]*folder)

	for i, fd := range failureDomains {
		tempFolder := new(folder)
		tempFolder.Datacenter = fd.Topology.Datacenter
		tempFolder.Name = fd.Topology.Folder

		// Only if the folder is empty do we create a folder resource
		// If a folder has been provided it means that it already exists
		// and it is to be used.
		if tempFolder.Name == "" {
			tempFolder.Name = infraID
			failureDomains[i].Topology.Folder = infraID
			key := fmt.Sprintf("%s-%s", tempFolder.Datacenter, tempFolder.Name)
			folders[key] = tempFolder
		}
	}
	return folders, nil
}

func convertVCentersToMap(values []vtypes.VCenter) map[string]vtypes.VCenter {
	vcenterMap := make(map[string]vtypes.VCenter)
	for _, v := range values {
		vcenterMap[v.Server] = v
	}
	return vcenterMap
}
