package bmc

import (
	"net"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// AccessDetailsFactory describes a callable that returns a new
// AccessDetails based on the input parameters.
type AccessDetailsFactory func(parsedURL *url.URL, disableCertificateVerification bool) (AccessDetails, error)

var factories = map[string]AccessDetailsFactory{}

// We could make this function public if we want to support
// out-of-tree factories.
func registerFactory(name string, factory AccessDetailsFactory) {
	factories[name] = factory
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

	// Boot interface to set
	BootInterface() string

	ManagementInterface() string
	PowerInterface() string
	RAIDInterface() string
	VendorInterface() string
}

func getParsedURL(address string) (parsedURL *url.URL, err error) {
	// Start by assuming "type://host:port"
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
	return parsedURL, nil
}

// NewAccessDetails creates an AccessDetails structure from the URL
// for a BMC.
func NewAccessDetails(address string, disableCertificateVerification bool) (AccessDetails, error) {

	parsedURL, err := getParsedURL(address)
	if err != nil {
		return nil, err
	}

	factory, ok := factories[parsedURL.Scheme]
	if !ok {
		return nil, &UnknownBMCTypeError{address, parsedURL.Scheme}
	}

	return factory(parsedURL, disableCertificateVerification)
}
