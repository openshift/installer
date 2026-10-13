// Copyright (c) 2016-2018 Hewlett Packard Enterprise Development LP

package bmc

import (
	"net/url"
)

func init() {
	RegisterFactory("ilo4", newILOAccessDetails, []string{"https"})
	RegisterFactory("ilo4-virtualmedia", newILOVirtualMediaAccessDetails, []string{"https"})
}

func newILOAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &iLOAccessDetails{
		bmcType:                        parsedURL.Scheme,
		portNum:                        parsedURL.Port(),
		hostname:                       parsedURL.Hostname(),
		disableCertificateVerification: disableCertificateVerification,
	}, nil
}

func newILOVirtualMediaAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &iLOAccessDetails{
		bmcType:                        parsedURL.Scheme,
		portNum:                        parsedURL.Port(),
		hostname:                       parsedURL.Hostname(),
		disableCertificateVerification: disableCertificateVerification,
		useVirtualMedia:                true,
	}, nil
}

type iLOAccessDetails struct {
	bmcType                        string
	portNum                        string
	hostname                       string
	disableCertificateVerification bool
	useVirtualMedia                bool
}

func (a *iLOAccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *iLOAccessDetails) NeedsMAC() bool {
	// For the inspection to work, we need a MAC address
	// https://github.com/metal3-io/baremetal-operator/pull/284#discussion_r317579040
	return true
}

func (a *iLOAccessDetails) Driver() string {
	return "ilo"
}

func (a *iLOAccessDetails) DisableCertificateVerification() bool {
	return a.disableCertificateVerification
}

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *iLOAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {

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

func (a *iLOAccessDetails) BIOSInterface() string {
	return ""
}

func (a *iLOAccessDetails) BootInterface() string {
	if a.useVirtualMedia {
		return "ilo-virtual-media"
	} else {
		return "ilo-ipxe"
	}
}

func (a *iLOAccessDetails) FirmwareInterface() string {
	return ""
}

func (a *iLOAccessDetails) ManagementInterface() string {
	return ""
}

func (a *iLOAccessDetails) PowerInterface() string {
	return ""
}

func (a *iLOAccessDetails) RAIDInterface() string {
	return "no-raid"
}

func (a *iLOAccessDetails) VendorInterface() string {
	return ""
}

func (a *iLOAccessDetails) SupportsSecureBoot() bool {
	return true
}

func (a *iLOAccessDetails) SupportsISOPreprovisioningImage() bool {
	return a.useVirtualMedia
}

func (a *iLOAccessDetails) RequiresProvisioningNetwork() bool {
	return !a.useVirtualMedia
}

func (a *iLOAccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
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
