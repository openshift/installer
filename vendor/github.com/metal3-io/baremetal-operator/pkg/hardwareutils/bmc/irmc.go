package bmc

import (
	"net/url"
)

func init() {
	RegisterFactory("irmc", newIRMCAccessDetails, []string{})
}

func newIRMCAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &iRMCAccessDetails{
		bmcType:                        parsedURL.Scheme,
		portNum:                        parsedURL.Port(),
		hostname:                       parsedURL.Hostname(),
		disableCertificateVerification: disableCertificateVerification,
	}, nil
}

type iRMCAccessDetails struct {
	bmcType                        string
	portNum                        string
	hostname                       string
	disableCertificateVerification bool
}

func (a *iRMCAccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *iRMCAccessDetails) NeedsMAC() bool {
	// For the inspection to work, we need a MAC address
	// https://github.com/metal3-io/baremetal-operator/pull/284#discussion_r317579040
	return true
}

func (a *iRMCAccessDetails) Driver() string {
	return "irmc"
}

func (a *iRMCAccessDetails) DisableCertificateVerification() bool {
	return a.disableCertificateVerification
}

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *iRMCAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {
	result := map[string]interface{}{
		"irmc_username": bmcCreds.Username,
		"irmc_password": bmcCreds.Password,
		"irmc_address":  a.hostname,
		"ipmi_username": bmcCreds.Username,
		"ipmi_password": bmcCreds.Password,
		"ipmi_address":  a.hostname,
	}

	if a.disableCertificateVerification {
		result["irmc_verify_ca"] = false
	}

	if a.portNum != "" {
		result["irmc_port"] = a.portNum
	}

	return result
}

func (a *iRMCAccessDetails) BIOSInterface() string {
	return ""
}

func (a *iRMCAccessDetails) BootInterface() string {
	return "ipxe"
}

func (a *iRMCAccessDetails) ManagementInterface() string {
	return ""
}

func (a *iRMCAccessDetails) PowerInterface() string {
	return "ipmitool"
}

func (a *iRMCAccessDetails) RAIDInterface() string {
	return "irmc"
}

func (a *iRMCAccessDetails) VendorInterface() string {
	return ""
}

func (a *iRMCAccessDetails) SupportsSecureBoot() bool {
	return true
}

func (a *iRMCAccessDetails) SupportsISOPreprovisioningImage() bool {
	return false
}

func (a *iRMCAccessDetails) RequiresProvisioningNetwork() bool {
	return true
}

func (a *iRMCAccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
	if firmwareConfig == nil {
		return nil, nil
	}

	var value string

	if firmwareConfig.VirtualizationEnabled != nil {
		value = "False"
		if *firmwareConfig.VirtualizationEnabled {
			value = "True"
		}
		settings = append(settings,
			map[string]string{
				"name":  "cpu_vt_enabled",
				"value": value,
			},
		)
	}

	if firmwareConfig.SimultaneousMultithreadingEnabled != nil {
		value = "False"
		if *firmwareConfig.SimultaneousMultithreadingEnabled {
			value = "True"
		}
		settings = append(settings,
			map[string]string{
				"name":  "hyper_threading_enabled",
				"value": value,
			},
		)
	}

	if firmwareConfig.SriovEnabled != nil {
		value = "False"
		if *firmwareConfig.SriovEnabled {
			value = "True"
		}
		settings = append(settings,
			map[string]string{
				"name":  "single_root_io_virtualization_support_enabled",
				"value": value,
			},
		)
	}

	return
}
