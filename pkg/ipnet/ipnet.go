// Package ipnet wraps net.IPNet to get CIDR serialization.
package ipnet

import (
	"encoding/json"
	"net"
	"reflect"

	"github.com/pkg/errors"
)

var nullString = "null"
var nullBytes = []byte(nullString)
var emptyIPNet = net.IPNet{}

// IPNet wraps net.IPNet to get CIDR serialization.
type IPNet struct {
	net.IPNet
}

// String returns a CIDR serialization of the subnet, or an empty
// string if the subnet is nil.
func (ipnet *IPNet) String() string {
	if ipnet == nil {
		return ""
	}
	return ipnet.IPNet.String()
}

// Version returns the IP version for an IPNet.
func (ipnet *IPNet) Version() int {
	// Per go documentation, To4() returns nil when using IPv6
	if ipnet.IP.To4() == nil {
		return 6
	}

	return 4
}

// CIDR returns the CIDR notation for a network, e.g. 255.255.255.0 would return 24.
func (ipnet IPNet) CIDR() int {
	cidr, _ := ipnet.IPNet.Mask.Size()
	return cidr
}

// MarshalJSON interface for an IPNet
func (ipnet IPNet) MarshalJSON() (data []byte, err error) {
	if reflect.DeepEqual(ipnet.IPNet, emptyIPNet) {
		return nullBytes, nil
	}

	return json.Marshal(ipnet.String())
}

// UnmarshalJSON interface for an IPNet
func (ipnet *IPNet) UnmarshalJSON(b []byte) (err error) {
	if string(b) == nullString {
		ipnet.IP = net.IP{}
		ipnet.Mask = net.IPMask{}
		return nil
	}

	var cidr string
	err = json.Unmarshal(b, &cidr)
	if err != nil {
		return errors.Wrap(err, "failed to Unmarshal string")
	}

	parsedIPNet, err := ParseCIDR(cidr)
	if err != nil {
		return errors.Wrap(err, "failed to Parse cidr string to net.IPNet")
	}

	*ipnet = *parsedIPNet

	return nil
}

// ParseCIDR parses a CIDR from its string representation.
func ParseCIDR(s string) (*IPNet, error) {
	ip, cidr, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}

	// This check is needed in order to work around a strange quirk in the Go
	// standard library. All of the addresses returned by net.ParseCIDR() are
	// 16-byte addresses. This does _not_ imply that they are IPv6 addresses,
	// which is what some libraries (e.g. github.com/apparentlymart/go-cidr)
	// assume. By forcing the address to be the expected length, we can work
	// around these bugs.
	if ip.To4() != nil {
		ip = ip.To4()
	}

	return &IPNet{
		IPNet: net.IPNet{
			IP:   ip,
			Mask: cidr.Mask,
		},
	}, nil
}

// MustParseCIDR parses a CIDR from its string representation. If the parse fails,
// the function will panic.
func MustParseCIDR(s string) *IPNet {
	cidr, err := ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return cidr
}
