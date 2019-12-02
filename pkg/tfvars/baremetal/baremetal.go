// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"
	"fmt"
	"github.com/metal3-io/baremetal-operator/pkg/bmc"
	"github.com/metal3-io/baremetal-operator/pkg/hardware"
	libvirttfvars "github.com/openshift/installer/pkg/tfvars/libvirt"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/pkg/errors"
	"net/url"
	"path"
	"strings"
)

type config struct {
	LibvirtURI              string `json:"libvirt_uri,omitempty"`
	BootstrapProvisioningIP string `json:"bootstrap_provisioning_ip,omitempty"`
	BootstrapOSImage        string `json:"bootstrap_os_image,omitempty"`
	ExternalBridge          string `json:"external_bridge,omitempty"`
	ProvisioningBridge      string `json:"provisioning_bridge,omitempty"`

	// Data required for control plane deployment - several maps per host, because of terraform's limitations
	Hosts         []map[string]interface{} `json:"hosts"`
	RootDevices   []map[string]interface{} `json:"root_devices"`
	Properties    []map[string]interface{} `json:"properties"`
	DriverInfos   []map[string]interface{} `json:"driver_infos"`
	InstanceInfos []map[string]interface{} `json:"instance_infos"`
}

// TFVars generates bare metal specific Terraform variables.
func TFVars(libvirtURI, bootstrapProvisioningIP, bootstrapOSImage, externalBridge, provisioningBridge string, platformHosts []*baremetal.Host, image string) ([]byte, error) {
	bootstrapOSImage, err := libvirttfvars.CachedImage(bootstrapOSImage)
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
		accessDetails, err := bmc.NewAccessDetails(host.BMC.Address)
		if err != nil {
			return nil, err
		}
		credentials := bmc.Credentials{
			Username: host.BMC.Username,
			Password: host.BMC.Password,
		}
		driverInfo := accessDetails.DriverInfo(credentials)
		driverInfo["deploy_kernel"] = fmt.Sprintf("http://%s/images/ironic-python-agent.kernel", bootstrapProvisioningIP)
		driverInfo["deploy_ramdisk"] = fmt.Sprintf("http://%s/images/ironic-python-agent.initramfs", bootstrapProvisioningIP)

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
		// The rhcos-downloader container downloads the image, compresses it to speed up deployments
		// and then makes it available on bootstrapProvisioningIP via http
		// The image is now formatted with a query string containing the sha256sum, we strip that here
		// and it will be consumed for validation in https://github.com/openshift/ironic-rhcos-downloader
		imageURL, err := url.Parse(image)
		if err != nil {
			return nil, err
		}
		imageURL.RawQuery = ""
		imageURL.Fragment = ""
		imageFilename := path.Base(imageURL.String())
		compressedImageFilename := strings.Replace(imageFilename, "openstack", "compressed", 1)
		cacheImageURL := fmt.Sprintf("http://%s/images/%s/%s", bootstrapProvisioningIP, imageFilename, compressedImageFilename)
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

	cfg := &config{
		LibvirtURI:              libvirtURI,
		BootstrapProvisioningIP: bootstrapProvisioningIP,
		BootstrapOSImage:        bootstrapOSImage,
		ExternalBridge:          externalBridge,
		ProvisioningBridge:      provisioningBridge,
		Hosts:                   hosts,
		Properties:              properties,
		DriverInfos:             driverInfos,
		RootDevices:             rootDevices,
		InstanceInfos:           instanceInfos,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
