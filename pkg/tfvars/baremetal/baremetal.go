// Package baremetal contains bare metal specific Terraform-variable logic.
package baremetal

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

	baremetalhost "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	"github.com/metal3-io/baremetal-operator/pkg/hardware"
	"github.com/metal3-io/baremetal-operator/pkg/hardwareutils/bmc"
	"github.com/pkg/errors"
	utilsnet "k8s.io/utils/net"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/openshift/installer/pkg/types/baremetal"
)

type config struct {
	LibvirtURI       string              `json:"libvirt_uri,omitempty"`
	IronicURI        string              `json:"ironic_uri,omitempty"`
	InspectorURI     string              `json:"inspector_uri,omitempty"`
	BootstrapOSImage string              `json:"bootstrap_os_image,omitempty"`
	Bridges          []map[string]string `json:"bridges"`

	IronicUsername string `json:"ironic_username"`
	IronicPassword string `json:"ironic_password"`

	DeploySteps []string `json:"deploy_steps"`

	// Data required for control plane deployment - several maps per host, because of terraform's limitations
	Masters       []map[string]interface{} `json:"masters"`
	RootDevices   []map[string]interface{} `json:"root_devices"`
	Properties    []map[string]interface{} `json:"properties"`
	DriverInfos   []map[string]interface{} `json:"driver_infos"`
	InstanceInfos []map[string]interface{} `json:"instance_infos"`
}

type imageDownloadFunc func(baseURL string) (string, error)

var (
	imageDownloader imageDownloadFunc
)

func init() {
	imageDownloader = cache.DownloadImageFile
}

func externalURLs(apiVIPs []string) (externalURLv4 string, externalURLv6 string) {
	if len(apiVIPs) > 1 {
		// IPv6 BMCs may not be able to reach IPv4 servers, use the right callback URL for them.
		externalURL := fmt.Sprintf("http://%s/", net.JoinHostPort(apiVIPs[1], "6180"))
		if utilsnet.IsIPv6String(apiVIPs[1]) {
			externalURLv6 = externalURL
		}
		if utilsnet.IsIPv4String(apiVIPs[1]) {
			externalURLv4 = externalURL
		}
	}

	return
}

// NOTE(dtantsur): this is a verbatim copy of the code from baremetal-operator
// that was not exposed in the version we vendor in 4.12.
func getParsedURL(address string) (parsedURL *url.URL, err error) {
	// Start by assuming "type://host:port"
	parsedURL, err = url.Parse(address)
	if err != nil {
		// We failed to parse the URL, but it may just be a host or
		// host:port string (which the URL parser rejects because ":"
		// is not allowed in the first segment of a
		// path. Unfortunately there is no error class to represent
		// that specific error, so we have to guess.
		if strings.Contains(address, ":") {
			// If we can parse host:port, carry on with those
			// values. Otherwise, report the original parser error.
			_, _, err2 := net.SplitHostPort(address)
			if err2 != nil {
				return nil, errors.Wrap(err, "failed to parse BMC address information")
			}
		}
		parsedURL = &url.URL{
			Scheme: "ipmi",
			Host:   address,
		}
	} else {
		// Successfully parsed the URL
		if parsedURL.Opaque != "" {
			parsedURL, err = url.Parse(strings.Replace(address, ":", "://", 1))
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse BMC address information")
			}
		}
		if parsedURL.Scheme == "" {
			if parsedURL.Hostname() == "" {
				// If there was no scheme at all, the hostname was
				// interpreted as a path.
				parsedURL, err = url.Parse(strings.Join([]string{"ipmi://", address}, ""))
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse BMC address information")
				}
			}
		}
	}
	return parsedURL, nil
}

// TFVars generates bare metal specific Terraform variables.
func TFVars(numControlPlaneReplicas int64, libvirtURI string, apiVIPs []string, imageCacheIP, bootstrapOSImage, externalBridge, externalMAC, provisioningBridge, provisioningMAC string, platformHosts []*baremetal.Host, hostFiles []*asset.File, image, ironicUsername, ironicPassword, ignition string) ([]byte, error) {
	bootstrapOSImage, err := imageDownloader(bootstrapOSImage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached bootstrap libvirt image")
	}

	externalURLv4, externalURLv6 := externalURLs(apiVIPs)

	var masters, rootDevices, properties, driverInfos, instanceInfos []map[string]interface{}
	var deploySteps []string

	// Select the first N hosts as masters, excluding the workers
	for i, host := range platformHosts {
		if len(masters) >= int(numControlPlaneReplicas) {
			break
		}

		if host.IsWorker() {
			//Skipping workers
			continue
		}

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
			// Some valid BMC addresses
			return nil, err
		}
		bmcURL, err := getParsedURL(host.BMC.Address)
		if err != nil {
			return nil, err
		}

		credentials := bmc.Credentials{
			Username: host.BMC.Username,
			Password: host.BMC.Password,
		}
		driverInfo := accessDetails.DriverInfo(credentials)
		driverInfo["deploy_kernel"] = fmt.Sprintf("http://%s/images/ironic-python-agent.kernel", net.JoinHostPort(imageCacheIP, "6180"))
		driverInfo["deploy_ramdisk"] = fmt.Sprintf("http://%s/%s.initramfs", net.JoinHostPort(imageCacheIP, "8084"), host.Name)
		driverInfo["deploy_iso"] = fmt.Sprintf("http://%s/%s.iso", net.JoinHostPort(imageCacheIP, "8084"), host.Name)
		if externalURLv6 != "" && utilsnet.IsIPv6String(bmcURL.Hostname()) {
			driverInfo["external_http_url"] = externalURLv6
		}
		if externalURLv4 != "" && utilsnet.IsIPv4String(bmcURL.Hostname()) {
			driverInfo["external_http_url"] = externalURLv4
		}

		var raidConfig, bmhFirmwareConfig, biosSettings []byte
		var bmcFirmwareConfig *bmc.FirmwareConfig
		var tmpBiosSettings []map[string]string
		var bmh baremetalhost.BareMetalHost

		err = yaml.Unmarshal(hostFiles[i].Data, &bmh)
		if err != nil {
			return nil, err
		}
		if bmh.Spec.RAID != nil {
			raidConfig, err = json.Marshal(bmh.Spec.RAID)
			if err != nil {
				return nil, err
			}
		}
		if bmh.Spec.Firmware != nil {
			bmhFirmwareConfig, err = json.Marshal(bmh.Spec.Firmware)
			if err != nil {
				return nil, err
			}
			if err = json.Unmarshal(bmhFirmwareConfig, &bmcFirmwareConfig); err != nil {
				return nil, err
			}
			tmpBiosSettings, err = accessDetails.BuildBIOSSettings(bmcFirmwareConfig)
			if err != nil {
				return nil, err
			}
			biosSettings, err = json.Marshal(tmpBiosSettings)
			if err != nil {
				return nil, err
			}
		}

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
			"deploy_interface":     "custom-agent",
			"raid_config":          string(raidConfig),
			"bios_settings":        string(biosSettings),
		}

		// Explicitly set the boot mode to the default "uefi" in case
		// it is not set. We use the capabilities field instead of
		// instance_info to ensure the host is in the right mode for
		// virtualmedia-based introspection.
		var bootMode string
		switch host.BootMode {
		case baremetal.Legacy:
			bootMode = "boot_mode:bios"
		case baremetal.UEFISecureBoot:
			bootMode = "boot_mode:uefi,secure_boot:true"
		default:
			bootMode = "boot_mode:uefi"
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

		// This is the only place where we need to set instance_info capabilities,
		// if we need to add another capabilitie we need merge the values
		// and ensure they are in the `key1:value1,key2:value2` format
		instanceInfo := make(map[string]interface{})
		if host.BootMode == baremetal.UEFISecureBoot {
			instanceInfo["capabilities"] = "secure_boot:true"
		}

		masters = append(masters, hostMap)
		// deploy_steps is set when a custom deployment is desired. We will use ironic's custom deployment
		// interface to use live ISO based installer. Currently this value is static but may be configurable
		// in the future.
		hostDeploySteps := `[{"interface": "deploy", "step": "install_coreos", "priority": 80, "args": {}}]`

		properties = append(properties, propertiesMap)
		driverInfos = append(driverInfos, driverInfo)
		rootDevices = append(rootDevices, rootDevice)
		instanceInfos = append(instanceInfos, instanceInfo)
		deploySteps = append(deploySteps, hostDeploySteps)
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

	ironicIP := apiVIPs[0]
	cfg := &config{
		LibvirtURI:       libvirtURI,
		IronicURI:        fmt.Sprintf("http://%s/v1", net.JoinHostPort(ironicIP, "6385")),
		InspectorURI:     fmt.Sprintf("http://%s/v1", net.JoinHostPort(ironicIP, "5050")),
		BootstrapOSImage: bootstrapOSImage,
		IronicUsername:   ironicUsername,
		IronicPassword:   ironicPassword,
		Masters:          masters,
		Bridges:          bridges,
		Properties:       properties,
		DriverInfos:      driverInfos,
		RootDevices:      rootDevices,
		InstanceInfos:    instanceInfos,
		DeploySteps:      deploySteps,
	}

	return json.MarshalIndent(cfg, "", "  ")
}
