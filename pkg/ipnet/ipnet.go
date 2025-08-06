// Package ipnet wraps net.IPNet to get CIDR serialization.
package ipnet

import (
	"encoding/json"
	"net"
	"reflect"

	"github.com/pkg/errors"
	netutils "k8s.io/utils/net"
)

var nullString = "null"
var nullBytes = []byte(nullString)
var emptyIPNet = net.IPNet{}

// IPNet wraps net.IPNet to get CIDR serialization.
// +kubebuilder:validation:Type=string
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

// IPNets represents a list of IPNet.
type IPNets []IPNet

func (ipnets IPNets) String() []string {
	cidrs := make([]string, len(ipnets))
	for i, ipnet := range ipnets {
		cidrs[i] = ipnet.String()
	}
	return cidrs
}

// IPv4Nets returns all IPNet of IPv4 familiy.
func (ipnets IPNets) IPv4Nets() IPNets {
	var ipv4Nets IPNets
	for _, ipnet := range ipnets {
		if netutils.IsIPv4CIDR(&ipnet.IPNet) {
			ipv4Nets = append(ipv4Nets, ipnet)
		}
	}
	return ipv4Nets
}

// IPv6Nets returns all IPNet of IPv6 familiy.
func (ipnets IPNets) IPv6Nets() IPNets {
	var ipv6Nets IPNets
	for _, ipnet := range ipnets {
		if netutils.IsIPv6CIDR(&ipnet.IPNet) {
			ipv6Nets = append(ipv6Nets, ipnet)
		}
	}
	return ipv6Nets
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

// DeepCopyInto copies the receiver into out.  out must be non-nil.
func (ipnet *IPNet) DeepCopyInto(out *IPNet) {
	if ipnet == nil {
		*out = IPNet{}
	} else {
		*out = *ipnet
	}
}

// DeepCopy copies the receiver, creating a new IPNet.
func (ipnet *IPNet) DeepCopy() *IPNet {
	if ipnet == nil {
		return nil
	}
	out := new(IPNet)
	ipnet.DeepCopyInto(out)
	return out
}
