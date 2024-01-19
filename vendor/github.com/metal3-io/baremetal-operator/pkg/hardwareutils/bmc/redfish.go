package bmc

import (
	"fmt"
	"net/url"
	"strings"
)

func init() {
	schemes := []string{"http", "https"}
	RegisterFactory("redfish", newRedfishAccessDetails, schemes)
	RegisterFactory("ilo5-redfish", newRedfishAccessDetails, schemes)
	RegisterFactory("idrac-redfish", newRedfishiDracAccessDetails, schemes)
}

func redfishDetails(parsedURL *url.URL, disableCertificateVerification bool) *redfishAccessDetails {
	return &redfishAccessDetails{
		bmcType:                        parsedURL.Scheme,
		host:                           parsedURL.Host,
		path:                           parsedURL.Path,
		disableCertificateVerification: disableCertificateVerification,
	}
}

func newRedfishAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return redfishDetails(parsedURL, disableCertificateVerification), nil
}

func newRedfishiDracAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &redfishiDracAccessDetails{
		*redfishDetails(parsedURL, disableCertificateVerification),
	}, nil
}

type redfishAccessDetails struct {
	bmcType                        string
	host                           string
	path                           string
	disableCertificateVerification bool
}

type redfishiDracAccessDetails struct {
	redfishAccessDetails
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

func (a *redfishAccessDetails) DisableCertificateVerification() bool {
	return a.disableCertificateVerification
}

func getRedfishAddress(bmcType, host string) string {
	redfishAddress := []string{}
	schemes := strings.Split(bmcType, "+")
	if len(schemes) > 1 {
		redfishAddress = append(redfishAddress, schemes[1])
	} else {
		redfishAddress = append(redfishAddress, redfishDefaultScheme)
	}
	redfishAddress = append(redfishAddress, "://")
	redfishAddress = append(redfishAddress, host)
	return strings.Join(redfishAddress, "")
}

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *redfishAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {
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

func (a *redfishAccessDetails) BIOSInterface() string {
	return ""
}

// That can be either pxe or redfish-virtual-media
func (a *redfishAccessDetails) BootInterface() string {
	return "ipxe"
}

func (a *redfishAccessDetails) FirmwareInterface() string {
	return "redfish"
}

func (a *redfishAccessDetails) ManagementInterface() string {
	return ""
}

func (a *redfishAccessDetails) PowerInterface() string {
	return ""
}

func (a *redfishAccessDetails) RAIDInterface() string {
	return "redfish"
}

func (a *redfishAccessDetails) VendorInterface() string {
	return ""
}

func (a *redfishAccessDetails) SupportsSecureBoot() bool {
	return true
}

func (a *redfishAccessDetails) SupportsISOPreprovisioningImage() bool {
	return false
}

func (a *redfishAccessDetails) RequiresProvisioningNetwork() bool {
	return true
}

func (a *redfishAccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
	if firmwareConfig != nil {
		return nil, fmt.Errorf("firmware settings for %s are not supported", a.Driver())
	}
	return nil, nil
}

// iDrac Redfish Overrides
func (a *redfishiDracAccessDetails) Driver() string {
	return "idrac"
}

func (a *redfishiDracAccessDetails) BIOSInterface() string {
	return "idrac-redfish"
}

func (a *redfishiDracAccessDetails) BootInterface() string {
	return "ipxe"
}

func (a *redfishiDracAccessDetails) FirmwareInterface() string {
	return "redfish"
}

func (a *redfishiDracAccessDetails) ManagementInterface() string {
	return "idrac-redfish"
}

func (a *redfishiDracAccessDetails) PowerInterface() string {
	return "idrac-redfish"
}

func (a *redfishiDracAccessDetails) RAIDInterface() string {
	return "idrac-redfish"
}

func (a *redfishiDracAccessDetails) VendorInterface() string {
	// NOTE(dtantsur): the idrac hardware type defaults to WSMAN vendor, we need to use the Redfish implementation.
	return "idrac-redfish"
}

func (a *redfishiDracAccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
	if firmwareConfig != nil {
		return nil, fmt.Errorf("firmware settings for %s are not supported", a.Driver())
	}
	return nil, nil
}
