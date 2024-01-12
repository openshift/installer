package vsphere

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/rhcos/cache"
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
	ControlPlaneNetworkKargs []string                                 `json:"vsphere_control_plane_network_kargs"`
	BootStrapNetworkKargs    string                                   `json:"vsphere_bootstrap_network_kargs"`
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
		cachedImage, err = cache.DownloadImageFile(sources.ImageURL, cache.InstallerApplicationName)
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
		OvaFilePath:               cachedImage,
		DiskType:                  sources.DiskType,
		VCenters:                  vcenterZones,
		NetworksInFailureDomains:  sources.NetworksInFailureDomain,
		ControlPlanes:             sources.ControlPlaneConfigs,
		DatacentersFolders:        datacentersFolders,
		ImportOvaFailureDomainMap: importOvaFailureDomainMap,
		FailureDomainMap:          failureDomainMap,
		ControlPlaneNetworkKargs:  []string{},
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

func getSubnetMask(prefix netip.Prefix) (string, error) {
	prefixLength := net.IPv4len * 8
	if prefix.Addr().Is6() {
		prefixLength = net.IPv6len * 8
	}
	ipMask := net.CIDRMask(prefix.Masked().Bits(), prefixLength)
	maskBytes, err := hex.DecodeString(ipMask.String())
	if err != nil {
		return "", err
	}
	ip := net.IP(maskBytes)
	maskStr := ip.To16().String()
	return maskStr, nil
}

func constructKargsFromNetworkConfig(ipAddrs []string, nameservers []string, gateway string) (string, error) {
	outKargs := ""

	var gatewayIP netip.Addr
	if len(gateway) > 0 {
		ip, err := netip.ParseAddr(gateway)
		if err != nil {
			return "", err
		}
		if ip.Is6() {
			gateway = fmt.Sprintf("[%s]", gateway)
		}
		gatewayIP = ip
	}

	for _, address := range ipAddrs {
		prefix, err := netip.ParsePrefix(address)
		if err != nil {
			return "", err
		}
		var ipStr, gatewayStr, maskStr string
		addr := prefix.Addr()
		switch {
		case addr.Is6():
			maskStr = fmt.Sprintf("%d", prefix.Bits())
			ipStr = fmt.Sprintf("[%s]", addr.String())
			if len(gateway) > 0 && gatewayIP.Is6() {
				gatewayStr = gateway
			}
		case addr.Is4():
			maskStr, err = getSubnetMask(prefix)
			if err != nil {
				return "", err
			}
			if len(gateway) > 0 && gatewayIP.Is4() {
				gatewayStr = gateway
			}
			ipStr = addr.String()
		default:
			return "", errors.New("IP address must adhere to IPv4 or IPv6 format")
		}
		outKargs += fmt.Sprintf("ip=%s::%s:%s:::none ", ipStr, gatewayStr, maskStr)
	}

	for _, nameserver := range nameservers {
		ip := net.ParseIP(nameserver)
		if ip.To4() == nil {
			nameserver = fmt.Sprintf("[%s]", nameserver)
		}
		outKargs += fmt.Sprintf("nameserver=%s ", nameserver)
	}

	outKargs = strings.Trim(outKargs, " ")
	logrus.Debugf("Generated karg: [%v].", outKargs)
	return outKargs, nil
}

// processGuestNetworkConfiguration takes the config and sources data and generates the kernel arguments (kargs)
// needed to boot RHCOS with static IP configurations.
func processGuestNetworkConfiguration(cfg *config, sources TFVarsSources) error {
	platform := sources.InstallConfig.Config.Platform.VSphere

	// Generate bootstrap karg using vsphere platform info from install-config
	for _, host := range platform.Hosts {
		if host.Role == vtypes.BootstrapRole {
			logrus.Debugf("Generating kargs for bootstrap")
			network := host.NetworkDevice
			kargs, err := constructKargsFromNetworkConfig(network.IPAddrs, network.Nameservers, network.Gateway)
			if err != nil {
				return err
			}
			cfg.BootStrapNetworkKargs = kargs
			break
		}
	}

	// Generate control plane kargs using info from machine network config
	for _, machine := range sources.ControlPlaneConfigs {
		logrus.Debugf("Generating kargs for control plane %v", machine.GenerateName)
		network := machine.Network.Devices[0]
		kargs, err := constructKargsFromNetworkConfig(network.IPAddrs, network.Nameservers, network.Gateway)
		if err != nil {
			return err
		}
		cfg.ControlPlaneNetworkKargs = append(cfg.ControlPlaneNetworkKargs, kargs)
	}
	return nil
}
