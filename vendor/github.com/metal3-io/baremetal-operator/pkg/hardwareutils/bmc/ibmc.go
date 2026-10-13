package bmc

import (
	"fmt"
	"net/url"
	"strings"
)

func init() {
	RegisterFactory("ibmc", newIbmcAccessDetails, []string{"http", "https"})
}

func newIbmcAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &ibmcAccessDetails{
		bmcType:                        parsedURL.Scheme,
		host:                           parsedURL.Host,
		path:                           parsedURL.Path,
		disableCertificateVerification: disableCertificateVerification,
	}, nil
}

type ibmcAccessDetails struct {
	bmcType                        string
	host                           string
	path                           string
	disableCertificateVerification bool
}

func (a *ibmcAccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *ibmcAccessDetails) NeedsMAC() bool {
	// For the inspection to work, we need a MAC address
	// https://github.com/metal3-io/baremetal-operator/pull/284#discussion_r317579040
	return true
}

func (a *ibmcAccessDetails) Driver() string {
	return "ibmc"
}

func (a *ibmcAccessDetails) DisableCertificateVerification() bool {
	return a.disableCertificateVerification
}

const ibmcDefaultScheme = "https"

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *ibmcAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {

	ibmcAddress := []string{}
	schemes := strings.Split(a.bmcType, "+")
	if len(schemes) > 1 {
		ibmcAddress = append(ibmcAddress, schemes[1])
	} else {
		ibmcAddress = append(ibmcAddress, ibmcDefaultScheme)
	}
	ibmcAddress = append(ibmcAddress, "://")
	ibmcAddress = append(ibmcAddress, a.host)
	ibmcAddress = append(ibmcAddress, a.path)

	result := map[string]interface{}{
		"ibmc_username": bmcCreds.Username,
		"ibmc_password": bmcCreds.Password,
		"ibmc_address":  strings.Join(ibmcAddress, ""),
	}

	if a.disableCertificateVerification {
		result["ibmc_verify_ca"] = false
	}

	return result
}

func (a *ibmcAccessDetails) BIOSInterface() string {
	return ""
}

func (a *ibmcAccessDetails) BootInterface() string {
	return "ipxe"
}

func (a *ibmcAccessDetails) FirmwareInterface() string {
	return ""
}

func (a *ibmcAccessDetails) ManagementInterface() string {
	return "ibmc"
}

func (a *ibmcAccessDetails) PowerInterface() string {
	return "ibmc"
}

func (a *ibmcAccessDetails) RAIDInterface() string {
	return "no-raid"
}

func (a *ibmcAccessDetails) VendorInterface() string {
	return ""
}

func (a *ibmcAccessDetails) SupportsSecureBoot() bool {
	return false
}

func (a *ibmcAccessDetails) SupportsISOPreprovisioningImage() bool {
	return false
}

func (a *ibmcAccessDetails) RequiresProvisioningNetwork() bool {
	return true
}

func (a *ibmcAccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
	if firmwareConfig != nil {
		return nil, fmt.Errorf("firmware settings for %s are not supported", a.Driver())
	}
	return nil, nil
}
