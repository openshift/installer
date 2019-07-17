// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"
	"github.com/metal3-io/baremetal-operator/pkg/bmc"
	"github.com/metal3-io/baremetal-operator/pkg/hardware"
	libvirttfvars "github.com/openshift/installer/pkg/tfvars/libvirt"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/pkg/errors"
)

type config struct {
	LibvirtURI         string `json:"libvirt_uri,omitempty"`
	IronicURI          string `json:"ironic_uri,omitempty"`
	Image              string `json:"os_image,omitempty"`
	ExternalBridge     string `json:"external_bridge,omitempty"`
	ProvisioningBridge string `json:"provisioning_bridge,omitempty"`

	// Data required for control plane deployment - several maps per host, because of terraform's limitations
	Hosts         []map[string]interface{} `json:"hosts"`
	RootDevices   []map[string]interface{} `json:"root_devices"`
	Properties    []map[string]interface{} `json:"properties"`
	DriverInfos   []map[string]interface{} `json:"driver_infos"`
	InstanceInfos []map[string]interface{} `json:"instance_infos"`
}

// TFVars generates bare metal specific Terraform variables.
func TFVars(libvirtURI, ironicURI, osImage, externalBridge, provisioningBridge string, platformHosts []*baremetal.Host, image baremetal.Image) ([]byte, error) {
	osImage, err := libvirttfvars.CachedImage(osImage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached libvirt image")
	}

	var hosts, rootDevices, properties, driverInfos, instanceInfos []map[string]interface{}

	for _, host := range platformHosts {
		// Get hardware profile
		if host.HardwareProfile == "default" {
			host.HardwareProfile = hardware.DefaultProfileName
		}

		profile, err := hardware.GetProfile(host.HardwareProfile)
		if err != nil {
			return nil, err
		}

		// BMC Driver Info
		accessDetails, err := bmc.NewAccessDetails(host.BMC.Address)
		if err != nil {
			return nil, err
		}
		credentials := bmc.Credentials{
			Username: host.BMC.Username,
			Password: host.BMC.Password,
		}
		driverInfo := accessDetails.DriverInfo(credentials)
		driverInfo["deploy_kernel"] = image.DeployKernel
		driverInfo["deploy_ramdisk"] = image.DeployRamdisk

		// Host Details
		hostMap := map[string]interface{}{
			"name":         host.Name,
			"port_address": host.BootMACAddress,
			"driver":       accessDetails.Type(),
		}

		// Properties
		propertiesMap := map[string]interface{}{
			"local_gb": profile.LocalGB,
			"cpu_arch": profile.CPUArch,
		}

		// Root device hints
		rootDevice := make(map[string]interface{})
		if profile.RootDeviceHints.HCTL != "" {
			rootDevice["hctl"] = profile.RootDeviceHints.HCTL
		} else {
			rootDevice["name"] = profile.RootDeviceHints.DeviceName
		}

		// Instance Info
		instanceInfo := map[string]interface{}{
			"root_gb":        25, // FIXME(stbenjam): Needed until https://storyboard.openstack.org/#!/story/2005165
			"image_source":   image.Source,
			"image_checksum": image.Checksum,
		}

		hosts = append(hosts, hostMap)
		properties = append(properties, propertiesMap)
		driverInfos = append(driverInfos, driverInfo)
		rootDevices = append(rootDevices, rootDevice)
		instanceInfos = append(instanceInfos, instanceInfo)
	}

	cfg := &config{
		LibvirtURI:         libvirtURI,
		IronicURI:          ironicURI,
		Image:              osImage,
		ExternalBridge:     externalBridge,
		ProvisioningBridge: provisioningBridge,
		Hosts:              hosts,
		Properties:         properties,
		DriverInfos:        driverInfos,
		RootDevices:        rootDevices,
		InstanceInfos:      instanceInfos,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
