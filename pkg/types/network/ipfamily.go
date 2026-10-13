// Package network contains types related to network configuration.
package network

// IPFamily specifies the IP address family for the cluster network.
type IPFamily string

const (
	// IPv4 indicates the cluster will use IPv4-only networking.
	// This is the default mode.
	IPv4 IPFamily = "IPv4"

	// DualStackIPv4Primary indicates the cluster will use dual-stack
	// networking with both IPv4 and IPv6 addresses, with IPv4 as the primary address family.
	DualStackIPv4Primary IPFamily = "DualStackIPv4Primary"

	// DualStackIPv6Primary indicates the cluster will use dual-stack
	// networking with both IPv4 and IPv6 addresses, with IPv6 as the primary address family.
	DualStackIPv6Primary IPFamily = "DualStackIPv6Primary"
)

// DualStackEnabled returns true if the IPFamily is configured for dual-stack networking
// (either DualStackIPv4Primary or DualStackIPv6Primary).
func (f IPFamily) DualStackEnabled() bool {
	return f == DualStackIPv4Primary || f == DualStackIPv6Primary
}
