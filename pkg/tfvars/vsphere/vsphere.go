package vsphere

import (
	"encoding/json"
	"fmt"
	"net"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	vtypes "github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pkg/errors"
)

type folder struct {
	Name       string `json:"name"`
	Datacenter string `json:"vsphere_datacenter"`
}

type config struct {
	OvaFilePath              string                                   `json:"vsphere_ova_filepath"`
	DiskType                 vtypes.DiskType                          `json:"vsphere_disk_type"`
	VCenters                 map[string]vtypes.VCenter                `json:"vsphere_vcenters"`
	FailureDomains           []vtypes.FailureDomain                   `json:"vsphere_failure_domains"`
	NetworksInFailureDomains map[string]string                        `json:"vsphere_networks"`
	ControlPlanes            []*machineapi.VSphereMachineProviderSpec `json:"vsphere_control_planes"`
	ControlPlaneNetworkKargs []string                                 `json:"vsphere_control_plane_network_kargs"`
	BootStrapNetworkKargs    string                                   `json:"vsphere_bootstrap_network_kargs"`
	DatacentersFolders       map[string]*folder                       `json:"vsphere_folders"`
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
	cachedImage, err := cache.DownloadImageFile(sources.ImageURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached vsphere image")
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
		FailureDomains:           sources.InstallConfig.Config.VSphere.FailureDomains,
		NetworksInFailureDomains: sources.NetworksInFailureDomain,
		ControlPlanes:            sources.ControlPlaneConfigs,
		DatacentersFolders:       datacentersFolders,
	}

	if len(sources.InstallConfig.Config.VSphere.Hosts) > 0 {
		processGuestNetworkConfiguration(cfg, sources)
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

func constructKargsFromNetworkConfig(networkConfig *vtypes.NetworkDeviceSpec) (string, error) {
	outKargs := ""
	// if an IPv4 gateway is defined, we'll only handle IPv4 addresses
	if len(networkConfig.Gateway4) > 0 {
		for _, address := range networkConfig.IPAddrs {
			ip, mask, err := net.ParseCIDR(address)
			if err != nil {
				return "", err
			}
			maskParts := mask.Mask
			maskStr := fmt.Sprintf("%d.%d.%d.%d", maskParts[0], maskParts[1], maskParts[2], maskParts[3])
			outKargs = outKargs + fmt.Sprintf("ip=%s::%s:%s:::none ", ip.String(), networkConfig.Gateway4, maskStr)
		}
	}

	for _, nameserver := range networkConfig.Nameservers {
		outKargs = outKargs + fmt.Sprintf("nameserver=%s ", nameserver)
	}

	return outKargs, nil
}

func processGuestNetworkConfiguration(cfg *config, sources TFVarsSources) error {
	platform := sources.InstallConfig.Config.Platform.VSphere
	for _, host := range platform.Hosts {
		switch host.Role {
		case vtypes.BootstrapRole:
			kArgs, err := constructKargsFromNetworkConfig(host.NetworkDevice)
			if err != nil {
				return err
			}
			cfg.BootStrapNetworkKargs = kArgs
		case vtypes.ControlPlaneRole:
			kArgs, err := constructKargsFromNetworkConfig(host.NetworkDevice)
			if err != nil {
				return err
			}
			cfg.ControlPlaneNetworkKargs = append(cfg.ControlPlaneNetworkKargs, kArgs)
		}
	}
	return nil
}
