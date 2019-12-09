package bmc

import (
	"net/url"
)

func init() {
	registerFactory("irmc", newIRMCAccessDetails)
}

func newIRMCAccessDetails(parsedURL *url.URL) (AccessDetails, error) {
	return &iRMCAccessDetails{
		bmcType:  parsedURL.Scheme,
		portNum:  parsedURL.Port(),
		hostname: parsedURL.Hostname(),
	}, nil
}

type iRMCAccessDetails struct {
	bmcType  string
	portNum  string
	hostname string
}

func (a *iRMCAccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *iRMCAccessDetails) NeedsMAC() bool {
	return false
}

func (a *iRMCAccessDetails) Driver() string {
	return "irmc"
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
	}

	if a.portNum != "" {
		result["irmc_port"] = a.portNum
	}

	return result
}

func (a *iRMCAccessDetails) BootInterface() string {
	return "pxe"
}

func (a *iRMCAccessDetails) ManagementInterface() string {
	return ""
}

func (a *iRMCAccessDetails) PowerInterface() string {
	return ""
}

func (a *iRMCAccessDetails) RAIDInterface() string {
	return ""
}

func (a *iRMCAccessDetails) VendorInterface() string {
	return ""
}
