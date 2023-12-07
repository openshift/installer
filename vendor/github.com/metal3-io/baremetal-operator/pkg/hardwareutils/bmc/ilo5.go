// Copyright (c) 2016-2018 Hewlett Packard Enterprise Development LP

package bmc

import (
	"net/url"
)

func init() {
	RegisterFactory("ilo5", newILO5AccessDetails, []string{"https"})
}

func newILO5AccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &iLO5AccessDetails{
		bmcType:                        parsedURL.Scheme,
		portNum:                        parsedURL.Port(),
		hostname:                       parsedURL.Hostname(),
		disableCertificateVerification: disableCertificateVerification,
	}, nil
}

type iLO5AccessDetails struct {
	bmcType                        string
	portNum                        string
	hostname                       string
	disableCertificateVerification bool
}

func (a *iLO5AccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *iLO5AccessDetails) NeedsMAC() bool {
	// For the inspection to work, we need a MAC address
	// https://github.com/metal3-io/baremetal-operator/pull/284#discussion_r317579040
	return true
}

func (a *iLO5AccessDetails) Driver() string {
	return "ilo5"
}

func (a *iLO5AccessDetails) DisableCertificateVerification() bool {
	return a.disableCertificateVerification
}

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *iLO5AccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {

	result := map[string]interface{}{
		"ilo_username": bmcCreds.Username,
		"ilo_password": bmcCreds.Password,
		"ilo_address":  a.hostname,
	}

	if a.disableCertificateVerification {
		result["ilo_verify_ca"] = false
	}

	if a.portNum != "" {
		result["client_port"] = a.portNum
	}

	return result
}

func (a *iLO5AccessDetails) BIOSInterface() string {
	return ""
}

func (a *iLO5AccessDetails) BootInterface() string {
	return "ilo-ipxe"
}

func (a *iLO5AccessDetails) FirmwareInterface() string {
	return ""
}

func (a *iLO5AccessDetails) ManagementInterface() string {
	return ""
}

func (a *iLO5AccessDetails) PowerInterface() string {
	return ""
}

func (a *iLO5AccessDetails) RAIDInterface() string {
	return "ilo5"
}

func (a *iLO5AccessDetails) VendorInterface() string {
	return ""
}

func (a *iLO5AccessDetails) SupportsSecureBoot() bool {
	return true
}

func (a *iLO5AccessDetails) SupportsISOPreprovisioningImage() bool {
	return false
}

func (a *iLO5AccessDetails) RequiresProvisioningNetwork() bool {
	return true
}

func (a *iLO5AccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
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
				"name":  "ProcHyperthreading",
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
				"name":  "Sriov",
				"value": value,
			},
		)
	}

	return
}
