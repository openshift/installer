package bmc

import (
	"strings"
)

type iDracAccessDetails struct {
	bmcType  string
	portNum  string
	hostname string
	path     string
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

func (a *iDracAccessDetails) BootInterface() string {
	return "ipxe"
}
