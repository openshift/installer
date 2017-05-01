package ipnets

import (
	"fmt"
	"math"
	"net"
)

// SubnetInto wraps SubnetShift and divides a network into at least count-many,
// equal-sized subnets, which are as large as allowed.
func SubnetInto(network *net.IPNet, count int) ([]*net.IPNet, error) {
	maskBits, _ := network.Mask.Size()
	hostBits := 32 - maskBits
	hostCount := 1 << uint(hostBits)

	// divide hosts among subnets
	ideal := float64(hostCount) / float64(count)
	// largest power of 2, not exceeding the ideal (float64 to int conversion
	// truncates toward zero)
	newHostBits := int(math.Log2(ideal))
	shift := hostBits - newHostBits
	return SubnetShift(network, shift)
}

// SubnetShift divides a network into subnets by shifting the given number of bits.
func SubnetShift(network *net.IPNet, bits int) ([]*net.IPNet, error) {
	if bits < 0 {
		return nil, fmt.Errorf("bit shift may not be negative, got %d", bits)
	}
	if bits > 31 {
		return nil, fmt.Errorf("network subnets cannot be divided %d times", bits)
	}
	// network divides into 2^bits subnets
	subnetCount := 1 << uint(bits)
	subnets := make([]*net.IPNet, subnetCount)

	// network info
	start := network.IP
	maskBits, _ := network.Mask.Size()
	hostBits := 32 - maskBits

	if maskBits+bits > 32 {
		return nil, fmt.Errorf("network subnet mask greater than /32, /%d is invalid", maskBits+bits)
	}

	// divide network into subnets
	newMaskBits := maskBits + bits
	newHostBits := hostBits - bits
	// subnet bitmasks are shifted by 'bits' places
	newMask := net.CIDRMask(newMaskBits, 32)

	// hosts per subnet
	hostCount := 1 << uint(newHostBits)

	for i := 0; i < subnetCount; i++ {
		ip := numeric(start) + uint32(i*hostCount)
		subnets[i] = &net.IPNet{
			IP:   bytewise(ip),
			Mask: newMask,
		}
	}

	return subnets, nil
}

// IP <-> integer transforms

// numeric returns a uint32 numeric representation of a net.IP.
func numeric(bytes net.IP) uint32 {
	var ip uint32
	// most significant to least significant
	for i, b := range []byte(bytes) {
		// bitwise or ("append" in this case)
		ip |= uint32(b) << (8 * uint32(3-i))
	}
	return ip
}

// bytewise returns a net.IP byte slice alias representation of an uint32.
// Note that not all uint32 values are valid IP addresses.
func bytewise(numeric uint32) net.IP {
	ip := make([]byte, 4)
	// least significant to most significant
	for i := 3; i >= 0; i-- {
		// AND away all but least significant
		ip[i] = byte(numeric & 0xFF)
		// nuke least significant byte
		numeric >>= 8
	}
	return net.IP(ip)
}
