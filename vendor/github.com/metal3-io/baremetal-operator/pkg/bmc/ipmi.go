package bmc

type ipmiAccessDetails struct {
	bmcType  string
	portNum  string
	hostname string
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

// DriverInfo returns a data structure to pass as the DriverInfo
// parameter when creating a node in Ironic. The structure is
// pre-populated with the access information, and the caller is
// expected to add any other information that might be needed (such as
// the kernel and ramdisk locations).
func (a *ipmiAccessDetails) DriverInfo(bmcCreds Credentials) map[string]interface{} {
	result := map[string]interface{}{
		"ipmi_port":     a.portNum,
		"ipmi_username": bmcCreds.Username,
		"ipmi_password": bmcCreds.Password,
		"ipmi_address":  a.hostname,
	}
	if a.portNum == "" {
		result["ipmi_port"] = ipmiDefaultPort
	}
	return result
}

func (a *ipmiAccessDetails) BootInterface() string {
	return "ipxe"
}
