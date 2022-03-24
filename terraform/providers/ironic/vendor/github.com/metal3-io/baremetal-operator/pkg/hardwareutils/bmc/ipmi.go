package bmc

import (
	"fmt"
	"net/url"
)

func init() {
	RegisterFactory("ipmi", newIPMIAccessDetails, []string{})
	RegisterFactory("libvirt", newIPMIAccessDetails, []string{})
}

func getPrivilegeLevel(rawquery string) string {
	privilegelevel := "ADMINISTRATOR"
	q, err := url.ParseQuery(rawquery)
	if err != nil {
		return privilegelevel
	}
	if val, ok := q["privilegelevel"]; ok {
		return val[0]
	}
	return privilegelevel
}

func newIPMIAccessDetails(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error) {
	return &ipmiAccessDetails{
		bmcType:                        parsedURL.Scheme,
		portNum:                        parsedURL.Port(),
		hostname:                       parsedURL.Hostname(),
		privilegelevel:                 getPrivilegeLevel(parsedURL.RawQuery),
		disableCertificateVerification: disableCertificateVerification,
	}, nil
}

type ipmiAccessDetails struct {
	bmcType                        string
	portNum                        string
	hostname                       string
	privilegelevel                 string
	disableCertificateVerification bool
}

const ipmiDefaultPort = "623"

func (a *ipmiAccessDetails) Type() string {
	return a.bmcType
}

// NeedsMAC returns true when the host is going to need a separate
// port created rather than having it discovered.
func (a *ipmiAccessDetails) NeedsMAC() bool {
	// libvirt-based hosts used for dev and testing require a MAC
	// address, specified as part of the host, but we don't want the
	// provisioner to have to know the rules about which drivers
	// require what so we hide that detail inside this class and just
	// let the provisioner know that "some" drivers require a MAC and
	// it should ask.
	return a.bmcType == "libvirt"
}

func (a *ipmiAccessDetails) Driver() string {
	return "ipmi"
}

func (a *ipmiAccessDetails) DisableCertificateVerification() bool {
	return a.disableCertificateVerification
}

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *ipmiAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {
	result := map[string]interface{}{
		"ipmi_port":       a.portNum,
		"ipmi_username":   bmcCreds.Username,
		"ipmi_password":   bmcCreds.Password,
		"ipmi_address":    a.hostname,
		"ipmi_priv_level": a.privilegelevel,
	}

	if a.disableCertificateVerification {
		result["ipmi_verify_ca"] = false
	}
	if a.portNum == "" {
		result["ipmi_port"] = ipmiDefaultPort
	}
	return result
}

func (a *ipmiAccessDetails) BIOSInterface() string {
	return ""
}

func (a *ipmiAccessDetails) BootInterface() string {
	return "ipxe"
}

func (a *ipmiAccessDetails) ManagementInterface() string {
	return ""
}

func (a *ipmiAccessDetails) PowerInterface() string {
	return ""
}

func (a *ipmiAccessDetails) RAIDInterface() string {
	return "no-raid"
}

func (a *ipmiAccessDetails) VendorInterface() string {
	return ""
}

func (a *ipmiAccessDetails) SupportsSecureBoot() bool {
	return false
}

func (a *ipmiAccessDetails) SupportsISOPreprovisioningImage() bool {
	return false
}

func (a *ipmiAccessDetails) RequiresProvisioningNetwork() bool {
	return true
}

func (a *ipmiAccessDetails) BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error) {
	if firmwareConfig != nil {
		return nil, fmt.Errorf("firmware settings for %s are not supported", a.Driver())
	}
	return nil, nil
}
