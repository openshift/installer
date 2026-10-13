package bmc

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// AccessDetailsFactory describes a callable that returns a new
// AccessDetails based on the input parameters.
type AccessDetailsFactory func(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error)

var factories = map[string]AccessDetailsFactory{}

// RegisterFactory maps a BMC type name to an AccessDetailsFactory,
// with optional scheme extensions.
//
// RegisterFactory("bmcname", theFunc, []string{"http", "https"})
// maps "bmcname", "bmcname+http", and "bmcname+https" to theFunc
func RegisterFactory(name string, factory AccessDetailsFactory, schemes []string) {
	factories[name] = factory

	for _, scheme := range schemes {
		factories[fmt.Sprintf("%s+%s", name, scheme)] = factory
	}
}

type FirmwareConfig struct {
	// Supports the virtualization of platform hardware.
	VirtualizationEnabled *bool

	// Allows a single physical processor core to appear as several logical processors.
	SimultaneousMultithreadingEnabled *bool

	// SR-IOV support enables a hypervisor to create virtual instances of a PCI-express device, potentially increasing performance.
	SriovEnabled *bool
}

// AccessDetails contains the information about how to get to a BMC.
//
// NOTE(dhellmann): This structure is very likely to change as we
// adapt it to additional types.
type AccessDetails interface {
	// Type returns the kind of the BMC, indicating the driver that
	// will be used to communicate with it.
	Type() string

	// NeedsMAC returns true when the host is going to need a separate
	// port created rather than having it discovered.
	NeedsMAC() bool

	// The name of the driver to instantiate the BMC with. This may differ
	// from the Type - both the ipmi and libvirt types use the ipmi driver.
	Driver() string

	// DriverInfo returns a data structure to pass as the DriverInfo
	// parameter when creating a node in Ironic. The structure is
	// pre-populated with the access information, and the caller is
	// expected to add any other information that might be needed
	// (such as the kernel and ramdisk locations).
	DriverInfo(bmcCreds Credentials) map[string]interface{}

	BIOSInterface() string

	// Boot interface to set
	BootInterface() string

	ManagementInterface() string
	PowerInterface() string
	RAIDInterface() string
	VendorInterface() string

	// Firmware interface to set
	FirmwareInterface() string

	// Whether the driver supports changing secure boot state.
	SupportsSecureBoot() bool

	// Whether the driver supports booting a preprovisioning image in ISO format
	SupportsISOPreprovisioningImage() bool

	// RequiresProvisioningNetwork checks the driver requires provisioning network
	RequiresProvisioningNetwork() bool

	// Build bios clean steps for ironic
	BuildBIOSSettings(firmwareConfig *FirmwareConfig) (settings []map[string]string, err error)
}

func GetParsedURL(address string) (parsedURL *url.URL, err error) {

	parsedURL, err = url.Parse(address)
	if err != nil {
		// We failed to parse the URL, but it may just be a host or
		// host:port string (which the URL parser rejects because ":"
		// is not allowed in the first segment of a
		// path. Unfortunately there is no error class to represent
		// that specific error, so we have to guess.
		if strings.Contains(address, ":") {
			// If we can parse host:port, carry on with those
			// values. Otherwise, report the original parser error.
			_, _, err2 := net.SplitHostPort(address)
			if err2 != nil {
				return nil, errors.Wrap(err, "failed to parse BMC address information")
			}
		}
		parsedURL = &url.URL{
			Scheme: "ipmi",
			Host:   address,
		}
	} else {
		// Successfully parsed the URL
		if parsedURL.Opaque != "" {
			parsedURL, err = url.Parse(strings.Replace(address, ":", "://", 1))
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse BMC address information")

			}
		}
		if parsedURL.Scheme == "" {
			if parsedURL.Hostname() == "" {
				// If there was no scheme at all, the hostname was
				// interpreted as a path.
				parsedURL, err = url.Parse(strings.Join([]string{"ipmi://", address}, ""))
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse BMC address information")
				}
			}
		}
	}

	// Check for expected hostname format
	if err := checkDNSValid(parsedURL.Hostname()); err != nil {
		return nil, errors.Wrap(err, "failed to parse BMC address information")
	}

	return parsedURL, nil
}

// NewAccessDetails creates an AccessDetails structure from the URL
// for a BMC.
func NewAccessDetails(address string, disableCertificateVerification bool) (AccessDetails, error) {

	if address == "" {
		return nil, errors.New("missing BMC address")
	}

	parsedURL, err := GetParsedURL(address)
	if err != nil {
		return nil, err
	}

	factory, ok := factories[parsedURL.Scheme]
	if !ok {
		return nil, &UnknownBMCTypeError{address, parsedURL.Scheme}
	}

	return factory(parsedURL, disableCertificateVerification)
}

func checkDNSValid(address string) error {

	// Allowing empty BMC address
	if address == "" {
		return nil
	}

	// Check if its a IPv6/IPv4 address
	if net.ParseIP(address) != nil {
		return nil
	}

	// Check if BMC address hostname follows DNS Standard
	valid, _ := regexp.MatchString(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]{0,61}[A-Za-z0-9])$`, address)
	if !valid {
		return fmt.Errorf("BMC address hostname/IP : [%s] is invalid", address)
	}
	return nil
}
