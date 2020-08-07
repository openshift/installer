// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"

	"github.com/metal3-io/baremetal-operator/pkg/bmc"
	"github.com/metal3-io/baremetal-operator/pkg/hardware"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/pkg/errors"
)

type config struct {
	LibvirtURI              string              `json:"libvirt_uri,omitempty"`
	BootstrapProvisioningIP string              `json:"bootstrap_provisioning_ip,omitempty"`
	BootstrapOSImage        string              `json:"bootstrap_os_image,omitempty"`
	Bridges                 []map[string]string `json:"bridges"`

	IronicUsername string `json:"ironic_username"`
	IronicPassword string `json:"ironic_password"`

	// Data required for control plane deployment - several maps per host, because of terraform's limitations
	Hosts         []map[string]interface{} `json:"hosts"`
	RootDevices   []map[string]interface{} `json:"root_devices"`
	Properties    []map[string]interface{} `json:"properties"`
	DriverInfos   []map[string]interface{} `json:"driver_infos"`
	InstanceInfos []map[string]interface{} `json:"instance_infos"`
}

// TFVars generates bare metal specific Terraform variables.
func TFVars(libvirtURI, bootstrapProvisioningIP, bootstrapOSImage, externalBridge, externalMAC, provisioningBridge, provisioningMAC string, platformHosts []*baremetal.Host, image, ironicUsername, ironicPassword string) ([]byte, error) {
	bootstrapOSImage, err := cache.DownloadImageFile(bootstrapOSImage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached bootstrap libvirt image")
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
		accessDetails, err := bmc.NewAccessDetails(host.BMC.Address, host.BMC.DisableCertificateVerification)
		if err != nil {
			return nil, err
		}
		credentials := bmc.Credentials{
			Username: host.BMC.Username,
			Password: host.BMC.Password,
		}
		driverInfo := accessDetails.DriverInfo(credentials)
		driverInfo["deploy_kernel"] = fmt.Sprintf("http://%s/images/ironic-python-agent.kernel", net.JoinHostPort(bootstrapProvisioningIP, "80"))
		driverInfo["deploy_ramdisk"] = fmt.Sprintf("http://%s/images/ironic-python-agent.initramfs", net.JoinHostPort(bootstrapProvisioningIP, "80"))

		// Host Details
		hostMap := map[string]interface{}{
			"name":                 host.Name,
			"port_address":         host.BootMACAddress,
			"driver":               accessDetails.Driver(),
			"boot_interface":       accessDetails.BootInterface(),
			"management_interface": accessDetails.ManagementInterface(),
			"power_interface":      accessDetails.PowerInterface(),
			"raid_interface":       accessDetails.RAIDInterface(),
			"vendor_interface":     accessDetails.VendorInterface(),
		}

		// Explicitly set the boot mode to the default "uefi" in case
		// it is not set. We use the capabilities field instead of
		// instance_info to ensure the host is in the right mode for
		// virtualmedia-based introspection.
		bootMode := "boot_mode:uefi"
		if host.BootMode == baremetal.Legacy {
			bootMode = "boot_mode:bios"
		}

		// Properties
		propertiesMap := map[string]interface{}{
			"local_gb":     profile.LocalGB,
			"cpu_arch":     profile.CPUArch,
			"capabilities": bootMode,
		}

		// Root device hints
		rootDevice := make(map[string]interface{})

		// host.RootDeviceHints overrides the root device hint in the profile
		if host.RootDeviceHints != nil {
			rootDeviceStringMap := host.RootDeviceHints.MakeHintMap()
			for key, value := range rootDeviceStringMap {
				rootDevice[key] = value
			}
		} else if profile.RootDeviceHints.HCTL != "" {
			rootDevice["hctl"] = profile.RootDeviceHints.HCTL
		} else {
			rootDevice["name"] = profile.RootDeviceHints.DeviceName
		}

		// Instance Info
		// The machine-os-downloader container downloads the image, compresses it to speed up deployments
		// and then makes it available on bootstrapProvisioningIP via http
		// The image is now formatted with a query string containing the sha256sum, we strip that here
		// and it will be consumed for validation in https://github.com/openshift/ironic-rhcos-downloader
		imageURL, err := url.Parse(image)
		if err != nil {
			return nil, err
		}
		imageURL.RawQuery = ""
		imageURL.Fragment = ""
		// We strip any .gz suffix because ironic-machine-os-downloader unzips the image
		// ref https://github.com/openshift/ironic-rhcos-downloader/pull/12
		imageFilename := path.Base(strings.TrimSuffix(imageURL.String(), ".gz"))
		compressedImageFilename := strings.Replace(imageFilename, "openstack", "compressed", 1)
		cacheImageURL := fmt.Sprintf("http://%s/images/%s/%s", net.JoinHostPort(bootstrapProvisioningIP, "80"), imageFilename, compressedImageFilename)
		cacheChecksumURL := fmt.Sprintf("%s.md5sum", cacheImageURL)
		instanceInfo := map[string]interface{}{
			"root_gb":        25, // FIXME(stbenjam): Needed until https://storyboard.openstack.org/#!/story/2005165
			"image_source":   cacheImageURL,
			"image_checksum": cacheChecksumURL,
		}

		hosts = append(hosts, hostMap)
		properties = append(properties, propertiesMap)
		driverInfos = append(driverInfos, driverInfo)
		rootDevices = append(rootDevices, rootDevice)
		instanceInfos = append(instanceInfos, instanceInfo)
	}

	var bridges []map[string]string

	bridges = append(bridges,
		map[string]string{
			"name": externalBridge,
			"mac":  externalMAC,
		})

	if provisioningBridge != "" {
		bridges = append(
			bridges,
			map[string]string{
				"name": provisioningBridge,
				"mac":  provisioningMAC,
			})
	}

	cfg := &config{
		LibvirtURI:              libvirtURI,
		BootstrapProvisioningIP: bootstrapProvisioningIP,
		BootstrapOSImage:        bootstrapOSImage,
		IronicUsername:          ironicUsername,
		IronicPassword:          ironicPassword,
		Hosts:                   hosts,
		Bridges:                 bridges,
		Properties:              properties,
		DriverInfos:             driverInfos,
		RootDevices:             rootDevices,
		InstanceInfos:           instanceInfos,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
