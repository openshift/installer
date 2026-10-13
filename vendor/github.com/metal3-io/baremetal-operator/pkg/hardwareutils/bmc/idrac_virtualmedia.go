package bmc

import (
	"fmt"
	"net/url"
)

func init() {
	schemes := []string{"http", "https"}
	RegisterFactory("idrac-virtualmedia", newRedfishiDracVirtualMediaAccessDetails, schemes)
}

func newRedfishiDracVirtualMediaAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &redfishiDracVirtualMediaAccessDetails{
		bmcType:                        parsedURL.Scheme,
		host:                           parsedURL.Host,
		path:                           parsedURL.Path,
		disableCertificateVerification: disableCertificateVerification,
	}, nil
}

type redfishiDracVirtualMediaAccessDetails struct {
	bmcType                        string
	host                           string
	path                           string
	disableCertificateVerification bool
}

func (a *redfishiDracVirtualMediaAccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *redfishiDracVirtualMediaAccessDetails) NeedsMAC() bool {
	// For the inspection to work, we need a MAC address
	// https://github.com/metal3-io/baremetal-operator/pull/284#discussion_r317579040
	return true
}

func (a *redfishiDracVirtualMediaAccessDetails) DisableCertificateVerification() bool {
	return a.disableCertificateVerification
}

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *redfishiDracVirtualMediaAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {
	result := map[string]interface{}{
		"redfish_system_id": a.path,
		"redfish_username":  bmcCreds.Username,
		"redfish_password":  bmcCreds.Password,
		"redfish_address":   getRedfishAddress(a.bmcType, a.host),
	}

	if a.disableCertificateVerification {
		result["redfish_verify_ca"] = false
	}

	return result
}

// iDrac Virtual Media Overrides

func (a *redfishiDracVirtualMediaAccessDetails) Driver() string {
	return "idrac"
}

func (a *redfishiDracVirtualMediaAccessDetails) BIOSInterface() string {
	return "idrac-redfish"
}

func (a *redfishiDracVirtualMediaAccessDetails) BootInterface() string {
	return "idrac-redfish-virtual-media"
}

func (a *redfishiDracVirtualMediaAccessDetails) FirmwareInterface() string {
	return "redfish"
}

func (a *redfishiDracVirtualMediaAccessDetails) ManagementInterface() string {
	return "idrac-redfish"
}

func (a *redfishiDracVirtualMediaAccessDetails) PowerInterface() string {
	return "idrac-redfish"
}

func (a *redfishiDracVirtualMediaAccessDetails) RAIDInterface() string {
	return "idrac-redfish"
}

func (a *redfishiDracVirtualMediaAccessDetails) VendorInterface() string {
	// NOTE(dtantsur): the idrac hardware type defaults to WSMAN vendor, we need to use the Redfish implementation.
	return "idrac-redfish"
}

func (a *redfishiDracVirtualMediaAccessDetails) SupportsSecureBoot() bool {
	return true
}

func (a *redfishiDracVirtualMediaAccessDetails) SupportsISOPreprovisioningImage() bool {
	return true
}

func (a *redfishiDracVirtualMediaAccessDetails) RequiresProvisioningNetwork() bool {
	return false
}

func (a *redfishiDracVirtualMediaAccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
	if firmwareConfig != nil {
		return nil, fmt.Errorf("firmware settings for %s are not supported", a.Driver())
	}
	return nil, nil
}
