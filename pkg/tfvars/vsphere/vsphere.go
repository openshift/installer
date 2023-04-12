package vsphere

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
		ControlPlaneNetworkKargs: []string{},
	}

	if len(sources.InstallConfig.Config.VSphere.Hosts) > 0 {
		logrus.Debugf("Applying static IP configs")
		err = processGuestNetworkConfiguration(cfg, sources)
		if err != nil {
			return nil, err
		}
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

// func constructKargsFromNetworkConfig(networkConfig *vtypes.NetworkDeviceSpec) (string, error) {
func constructKargsFromNetworkConfig(ipAddrs, nameservers []string, gateway4 string) (string, error) {
	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		logrus.Tracef("Constructing kargs from IPs [%v] nameservers [%v] gateway4 [%v]", ipAddrs, nameservers, gateway4)
	}
	outKargs := ""
	// if an IPv4 gateway is defined, we'll only handle IPv4 addresses
	if len(gateway4) > 0 {
		for _, address := range ipAddrs {
			ip, mask, err := net.ParseCIDR(address)
			if err != nil {
				return "", err
			}
			maskParts := mask.Mask
			maskStr := fmt.Sprintf("%d.%d.%d.%d", maskParts[0], maskParts[1], maskParts[2], maskParts[3])
			outKargs += fmt.Sprintf("ip=%s::%s:%s:::none ", ip.String(), gateway4, maskStr)
		}
	}

	for _, nameserver := range nameservers {
		outKargs += fmt.Sprintf("nameserver=%s ", nameserver)
	}

	logrus.Debugf("Generated karg: [%v]", outKargs)
	return outKargs, nil
}

// processGuestNetworkConfiguration takes the config and sources data and generates the kernel arguments (kargs)
// needed to boot RHCOS with static IP configurations.
func processGuestNetworkConfiguration(cfg *config, sources TFVarsSources) error {
	platform := sources.InstallConfig.Config.Platform.VSphere

	// Generate bootstrap karg using vsphere platform info from install-config
	for _, host := range platform.Hosts {
		if host.Role == vtypes.BootstrapRole {
			logrus.Debugf("Generating kargs for bootstrap.")
			network := host.NetworkDevice
			kargs, err := constructKargsFromNetworkConfig(network.IPAddrs, network.Nameservers, network.Gateway4)
			if err != nil {
				return err
			}
			cfg.BootStrapNetworkKargs = kargs
		}
	}

	// Generate control plane kargs using info from machine network config
	for _, machine := range sources.ControlPlaneConfigs {
		logrus.Debugf("Generating kargs for control plane %v.", machine.GenerateName)
		network := machine.Network.Devices[0]
		kargs, err := constructKargsFromNetworkConfig(network.IPAddrs, network.Nameservers, network.Gateway4)
		if err != nil {
			return err
		}
		cfg.ControlPlaneNetworkKargs = append(cfg.ControlPlaneNetworkKargs, kargs)
	}
	return nil
}
