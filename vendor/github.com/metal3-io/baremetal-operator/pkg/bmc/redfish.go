package bmc

import (
	"net/url"
	"strings"
)

func init() {
	registerFactory("redfish", newRedfishAccessDetails)
	registerFactory("redfish+http", newRedfishAccessDetails)
	registerFactory("redfish+https", newRedfishAccessDetails)
}

func newRedfishAccessDetails(parsedURL *url.URL) (AccessDetails, error) {
	return &redfishAccessDetails{
		bmcType:  parsedURL.Scheme,
		host:     parsedURL.Host,
		path:     parsedURL.Path,
	}, nil
}

type redfishAccessDetails struct {
	bmcType  string
	host     string
	path     string
}

const redfishDefaultScheme = "https"

func (a *redfishAccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *redfishAccessDetails) NeedsMAC() bool {
	// For the inspection to work, we need a MAC address
	// https://github.com/metal3-io/baremetal-operator/pull/284#discussion_r317579040
	return true
}

func (a *redfishAccessDetails) Driver() string {
	return "redfish"
}

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *redfishAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {
	redfishAddress := []string{}
	schemes := strings.Split(a.bmcType, "+")
	if len(schemes) > 1 {
		redfishAddress = append(redfishAddress, schemes[1])
	} else {
		redfishAddress = append(redfishAddress, redfishDefaultScheme)
	}
	redfishAddress = append(redfishAddress, "://")
	redfishAddress = append(redfishAddress, a.host)

	result := map[string]interface{}{
		"redfish_system_id":     a.path,
		"redfish_username": bmcCreds.Username,
		"redfish_password": bmcCreds.Password,
		"redfish_address": strings.Join(redfishAddress, ""),
	}

	return result
}

// That can be either pxe or redfish-virtual-media
func (a *redfishAccessDetails) BootInterface() string {
	return "pxe"
}
