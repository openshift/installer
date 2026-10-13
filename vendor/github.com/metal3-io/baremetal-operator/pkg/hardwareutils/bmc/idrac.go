package bmc

import (
	"net/url"
	"strings"
)

func init() {
	RegisterFactory("idrac", newIDRACAccessDetails, []string{"http", "https"})
}

func newIDRACAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &iDracAccessDetails{
		bmcType:                        parsedURL.Scheme,
		portNum:                        parsedURL.Port(),
		hostname:                       parsedURL.Hostname(),
		path:                           parsedURL.Path,
		disableCertificateVerification: disableCertificateVerification,
	}, nil
}

type iDracAccessDetails struct {
	bmcType                        string
	portNum                        string
	hostname                       string
	path                           string
	disableCertificateVerification bool
}

func (a *iDracAccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *iDracAccessDetails) NeedsMAC() bool {
	return false
}

func (a *iDracAccessDetails) Driver() string {
	return "idrac"
}

func (a *iDracAccessDetails) DisableCertificateVerification() bool {
	return a.disableCertificateVerification
}

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *iDracAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {
	result := map[string]interface{}{
		"drac_username": bmcCreds.Username,
		"drac_password": bmcCreds.Password,
		"drac_address":  a.hostname,
	}
	if a.disableCertificateVerification {
		result["drac_verify_ca"] = false
	}

	schemes := strings.Split(a.bmcType, "+")
	if len(schemes) > 1 {
		result["drac_protocol"] = schemes[1]
	}
	if a.portNum != "" {
		result["drac_port"] = a.portNum
	}
	if a.path != "" {
		result["drac_path"] = a.path
	}

	return result
}

func (a *iDracAccessDetails) BIOSInterface() string {
	return ""
}

func (a *iDracAccessDetails) BootInterface() string {
	return "ipxe"
}

func (a *iDracAccessDetails) FirmwareInterface() string {
	return ""
}

func (a *iDracAccessDetails) ManagementInterface() string {
	return ""
}

func (a *iDracAccessDetails) PowerInterface() string {
	return ""
}

func (a *iDracAccessDetails) RAIDInterface() string {
	return "idrac-wsman"
}

func (a *iDracAccessDetails) VendorInterface() string {
	return ""
}

// NOTE(dtantsur): change to true if we switch to redfish-based implementations
// by default.
func (a *iDracAccessDetails) SupportsSecureBoot() bool {
	return false
}

func (a *iDracAccessDetails) SupportsISOPreprovisioningImage() bool {
	return false
}

func (a *iDracAccessDetails) RequiresProvisioningNetwork() bool {
	return true
}

func (a *iDracAccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
	if firmwareConfig == nil {
		return nil, nil
	}

	var value string

	if firmwareConfig.VirtualizationEnabled != nil {
		value = "Disabled"
		if *firmwareConfig.VirtualizationEnabled {
			value = "Enabled"
		}
		settings = append(settings,
			map[string]string{
				"name":  "ProcVirtualization",
				"value": value,
			},
		)
	}

	if firmwareConfig.SimultaneousMultithreadingEnabled != nil {
		value = "Disabled"
		if *firmwareConfig.SimultaneousMultithreadingEnabled {
			value = "Enabled"
		}
		settings = append(settings,
			map[string]string{
				"name":  "LogicalProc",
				"value": value,
			},
		)
	}

	if firmwareConfig.SriovEnabled != nil {
		value = "Disabled"
		if *firmwareConfig.SriovEnabled {
			value = "Enabled"
		}
		settings = append(settings,
			map[string]string{
				"name":  "SriovGlobalEnable",
				"value": value,
			},
		)
	}

	return
}
