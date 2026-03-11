package azure

import (
	"fmt"
	"net"

	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils/cidr"
	"github.com/openshift/installer/pkg/ipnet"
)

// AddressFamilySubnet holds an IPv4 and IPv6 subnet pair.
type AddressFamilySubnet struct {
	IPv4Subnet *net.IPNet
	IPv6Subnet *net.IPNet
	SubnetRole *capz.SubnetRole
}

// DeepCopyInto deep copies an AddressFamilySubnet struct.
func (in *AddressFamilySubnet) DeepCopyInto(out *AddressFamilySubnet) {
	*out = *in
	if in.IPv4Subnet != nil {
		out.IPv4Subnet = &net.IPNet{
			IP:   make(net.IP, len(in.IPv4Subnet.IP)),
			Mask: make(net.IPMask, len(in.IPv4Subnet.Mask)),
		}
		copy(out.IPv4Subnet.IP, in.IPv4Subnet.IP)
		copy(out.IPv4Subnet.Mask, in.IPv4Subnet.Mask)
	}
	if in.IPv6Subnet != nil {
		out.IPv6Subnet = &net.IPNet{
			IP:   make(net.IP, len(in.IPv6Subnet.IP)),
			Mask: make(net.IPMask, len(in.IPv6Subnet.Mask)),
		}
		copy(out.IPv6Subnet.IP, in.IPv6Subnet.IP)
		copy(out.IPv6Subnet.Mask, in.IPv6Subnet.Mask)
	}
	if in.SubnetRole != nil {
		in, out := &in.SubnetRole, &out.SubnetRole
		*out = new(capz.SubnetRole)
		**out = **in
	}
}

// AddressFamilySubnets keeps track of IPv4 and IPv6 subnets.
type AddressFamilySubnets struct {
	addressFamilySubnets []AddressFamilySubnet
	length               int
	ipv4Count            int
	ipv6Count            int
	// XXX: leave this for now, fix to use install-config later
	isDualStack bool
}

// Length returns the number of subnets in the list.
func (a AddressFamilySubnets) Length() int {
	return a.length
}

// IPv4Count returns the number of IPv4 subnets in the list.
func (a AddressFamilySubnets) IPv4Count() int {
	return a.ipv4Count
}

// IPv6Count returns the number of IPv6 subnets in the list.
func (a AddressFamilySubnets) IPv6Count() int {
	return a.ipv6Count
}

// IsDualStack determines if we are using single or dual stack networking.
func (a AddressFamilySubnets) IsDualStack() bool {
	return a.isDualStack
}

// GetIPv4Subnets returns all IPv4 subnets in the list.
func (a AddressFamilySubnets) GetIPv4Subnets() []*net.IPNet {
	var ipv4Subnets []*net.IPNet
	for _, ipv4Subnet := range a.addressFamilySubnets {
		if ipv4Subnet.IPv4Subnet != nil {
			tmp := new(AddressFamilySubnet)
			ipv4Subnet.DeepCopyInto(tmp)
			ipv4Subnets = append(ipv4Subnets, tmp.IPv4Subnet)
		}
	}
	return ipv4Subnets
}

// GetIPv6Subnets returns all IPv6 subnets in the list.
func (a AddressFamilySubnets) GetIPv6Subnets() []*net.IPNet {
	var ipv6Subnets []*net.IPNet
	for _, ipv6Subnet := range a.addressFamilySubnets {
		if ipv6Subnet.IPv6Subnet != nil {
			tmp := new(AddressFamilySubnet)
			ipv6Subnet.DeepCopyInto(tmp)
			ipv6Subnets = append(ipv6Subnets, tmp.IPv6Subnet)
		}
	}
	return ipv6Subnets
}

// GetControlPlaneSubnet returns the control plane subnet.
func (a AddressFamilySubnets) GetControlPlaneSubnet() AddressFamilySubnet {
	if a.length > 0 {
		tmp := new(AddressFamilySubnet)
		a.addressFamilySubnets[0].DeepCopyInto(tmp)
		return *tmp
	}
	return AddressFamilySubnet{}
}

// SplitIPv4ComputeSubnets splits the IPv4 compute subnet into the number of subnets specified.
func (a AddressFamilySubnets) SplitIPv4ComputeSubnets(numSubnets int) ([]*net.IPNet, error) {
	if a.ipv4Count > 1 && numSubnets > 0 {
		return cidr.SplitIntoSubnetsIPv4(a.GetIPv4Subnets()[1].String(), numSubnets)
	}
	return []*net.IPNet{}, nil
}

// SplitIPv6ComputeSubnets splits the IPv6 compute subnet into the number of subnets specified.
func (a AddressFamilySubnets) SplitIPv6ComputeSubnets(numSubnets int) ([]*net.IPNet, error) {
	if a.ipv6Count > 1 && numSubnets > 0 {
		return cidr.SplitIntoSubnetsIPv6(a.GetIPv6Subnets()[1].String(), numSubnets)
	}
	return []*net.IPNet{}, nil
}

// SplitComputeSubnet splits the compute subnet into the number of subnets specified.
func (a *AddressFamilySubnets) SplitComputeSubnet(numSubnets int) error {
	if numSubnets <= 0 {
		return nil
	}

	ipv4ComputeSubnets, err := a.SplitIPv4ComputeSubnets(numSubnets)
	if err != nil {
		return err
	}

	var ipv6Count int
	var ipv6ComputeSubnets []*net.IPNet
	if a.isDualStack {
		ipv6ComputeSubnets, err = a.SplitIPv6ComputeSubnets(numSubnets)
		if err != nil {
			return err
		}
		ipv6Count = len(ipv6ComputeSubnets)
	}

	addressFamilySubnets := make([]AddressFamilySubnet, 0)
	for i, ipv4ComputeSubnet := range ipv4ComputeSubnets {
		addressFamilySubnet := AddressFamilySubnet{
			IPv4Subnet: ipv4ComputeSubnet,
			SubnetRole: ptr.To(capz.SubnetNode),
		}
		if i < ipv6Count {
			addressFamilySubnet.IPv6Subnet = ipv6ComputeSubnets[i]
		}
		addressFamilySubnets = append(addressFamilySubnets, addressFamilySubnet)
	}

	if a.length > 2 {
		a.addressFamilySubnets = append(a.addressFamilySubnets[:1], append(addressFamilySubnets, a.addressFamilySubnets[2:]...)...)
	} else {
		a.addressFamilySubnets = append(a.addressFamilySubnets[:1], addressFamilySubnets...)
	}

	a.length = len(a.addressFamilySubnets)
	a.ipv4Count = len(a.GetIPv4Subnets())
	if a.isDualStack {
		a.ipv6Count = len(a.GetIPv6Subnets())
	}

	return nil
}

// GetComputeSubnets returns the compute subnets.
func (a AddressFamilySubnets) GetComputeSubnets() []AddressFamilySubnet {
	var computeSubnets []AddressFamilySubnet

	for _, addressFamilySubnet := range a.addressFamilySubnets {
		if addressFamilySubnet.SubnetRole != nil && *addressFamilySubnet.SubnetRole == capz.SubnetNode {
			tmp := new(AddressFamilySubnet)
			addressFamilySubnet.DeepCopyInto(tmp)
			computeSubnets = append(computeSubnets, *tmp)
		}
	}

	return computeSubnets
}

// GetIPv4ComputeSubnets returns the IPv4 compute subnets.
func (a AddressFamilySubnets) GetIPv4ComputeSubnets() []*net.IPNet {
	computeSubnets := a.GetComputeSubnets()

	var ipv4ComputeSubnets []*net.IPNet
	for _, ipv4ComputeSubnet := range computeSubnets {
		if ipv4ComputeSubnet.IPv4Subnet != nil {
			tmp := new(AddressFamilySubnet)
			ipv4ComputeSubnet.DeepCopyInto(tmp)
			ipv4ComputeSubnets = append(ipv4ComputeSubnets, tmp.IPv4Subnet)
		}
	}

	return ipv4ComputeSubnets
}

// GetIPv6ComputeSubnets returns the IPv6 compute subnets.
func (a AddressFamilySubnets) GetIPv6ComputeSubnets() []*net.IPNet {
	computeSubnets := a.GetComputeSubnets()

	var ipv6ComputeSubnets []*net.IPNet
	for _, ipv6ComputeSubnet := range computeSubnets {
		if ipv6ComputeSubnet.IPv6Subnet != nil {
			tmp := new(AddressFamilySubnet)
			ipv6ComputeSubnet.DeepCopyInto(tmp)
			ipv6ComputeSubnets = append(ipv6ComputeSubnets, tmp.IPv6Subnet)
		}
	}

	return ipv6ComputeSubnets
}

// GetAdditionalSubnets returns any additional subnets.
func (a AddressFamilySubnets) GetAdditionalSubnets() []AddressFamilySubnet {
	var additionalSubnets []AddressFamilySubnet

	for _, addressFamilySubnet := range a.addressFamilySubnets {
		if addressFamilySubnet.SubnetRole == nil {
			tmp := new(AddressFamilySubnet)
			addressFamilySubnet.DeepCopyInto(tmp)
			additionalSubnets = append(additionalSubnets, *tmp)
		}
	}

	return additionalSubnets
}

// GetIPv4AdditionalSubnets returns any additional IPv4 subnets.
func (a AddressFamilySubnets) GetIPv4AdditionalSubnets() []*net.IPNet {
	additionalSubnets := a.GetAdditionalSubnets()

	var ipv4Subnets []*net.IPNet
	for _, ipv4Subnet := range additionalSubnets {
		if ipv4Subnet.IPv4Subnet != nil {
			tmp := new(AddressFamilySubnet)
			ipv4Subnet.DeepCopyInto(tmp)
			ipv4Subnets = append(ipv4Subnets, tmp.IPv4Subnet)
		}
	}

	return ipv4Subnets
}

// GetIPv6AdditionalSubnets returns any additional IPv6 subnets.
func (a AddressFamilySubnets) GetIPv6AdditionalSubnets() []*net.IPNet {
	additionalSubnets := a.GetAdditionalSubnets()

	var ipv6Subnets []*net.IPNet
	for _, ipv6Subnet := range additionalSubnets {
		if ipv6Subnet.IPv6Subnet != nil {
			tmp := new(AddressFamilySubnet)
			ipv6Subnet.DeepCopyInto(tmp)
			ipv6Subnets = append(ipv6Subnets, tmp.IPv6Subnet)
		}
	}

	return ipv6Subnets
}

// GetIPv4CIDRBlocks returns all IPv4 CIDR blocks.
func (a AddressFamilySubnets) GetIPv4CIDRBlocks() []string {
	var ipv4CIDRBlocks []string

	for _, ipv4Subnet := range a.addressFamilySubnets {
		if ipv4Subnet.IPv4Subnet != nil {
			ipv4CIDRBlocks = append(ipv4CIDRBlocks, (*ipv4Subnet.IPv4Subnet).String())
		}
	}

	return ipv4CIDRBlocks
}

// GetIPv6CIDRBlocks returns all IPv6 CIDR blocks.
func (a AddressFamilySubnets) GetIPv6CIDRBlocks() []string {
	var ipv6CIDRBlocks []string

	for _, ipv6Subnet := range a.addressFamilySubnets {
		if ipv6Subnet.IPv6Subnet != nil {
			ipv6CIDRBlocks = append(ipv6CIDRBlocks, (*ipv6Subnet.IPv6Subnet).String())
		}
	}

	return ipv6CIDRBlocks
}

// GetIPv4ControlPlaneCIDRBlocks returns all IPv4 control plane CIDR blocks.
func (a AddressFamilySubnets) GetIPv4ControlPlaneCIDRBlocks() string {
	ipv4Subnet := a.GetControlPlaneSubnet()
	if ipv4Subnet.IPv4Subnet != nil {
		return (*ipv4Subnet.IPv4Subnet).String()
	}
	return ""
}

// GetIPv6ControlPlaneCIDRBlocks returns all IPv6 control plane CIDR blocks.
func (a AddressFamilySubnets) GetIPv6ControlPlaneCIDRBlocks() string {
	ipv6Subnet := a.GetControlPlaneSubnet()
	if ipv6Subnet.IPv6Subnet != nil {
		return (*ipv6Subnet.IPv6Subnet).String()
	}
	return ""
}

// GetIPv4ComputeCIDRBlocks returns all IPv4 compute CIDR blocks.
func (a AddressFamilySubnets) GetIPv4ComputeCIDRBlocks() []string {
	var ipv4CIDRBlocks []string

	for _, computeSubnet := range a.GetComputeSubnets() {
		if computeSubnet.IPv4Subnet != nil {
			ipv4CIDRBlocks = append(ipv4CIDRBlocks, (*computeSubnet.IPv4Subnet).String())
		}
	}

	return ipv4CIDRBlocks
}

// GetIPv6ComputeCIDRBlocks returns all IPv6 compute CIDR blocks.
func (a AddressFamilySubnets) GetIPv6ComputeCIDRBlocks() []string {
	var ipv6CIDRBlocks []string

	for _, computeSubnet := range a.GetComputeSubnets() {
		if computeSubnet.IPv6Subnet != nil {
			ipv6CIDRBlocks = append(ipv6CIDRBlocks, (*computeSubnet.IPv6Subnet).String())
		}
	}

	return ipv6CIDRBlocks
}

// GetIPv4AdditionalCIDRBlocks returns all additional IPv4 CIDR blocks.
func (a AddressFamilySubnets) GetIPv4AdditionalCIDRBlocks() []string {
	additionalSubnets := a.GetAdditionalSubnets()

	var ipv4CIDRBlocks []string
	for _, additionalSubnet := range additionalSubnets {
		if additionalSubnet.IPv4Subnet != nil {
			ipv4CIDRBlocks = append(ipv4CIDRBlocks, (*additionalSubnet.IPv4Subnet).String())
		}
	}

	return ipv4CIDRBlocks
}

// GetIPv6AdditionalCIDRBlocks returns all additional IPv6 CIDR blocks.
func (a AddressFamilySubnets) GetIPv6AdditionalCIDRBlocks() []string {
	additionalSubnets := a.GetAdditionalSubnets()

	var ipv6CIDRBlocks []string
	for _, additionalSubnet := range additionalSubnets {
		if additionalSubnet.IPv6Subnet != nil {
			ipv6CIDRBlocks = append(ipv6CIDRBlocks, (*additionalSubnet.IPv6Subnet).String())
		}
	}

	return ipv6CIDRBlocks
}

// GetAddressFamilySubnets takes a list of cidrs and parses it into paired subnets.
func GetAddressFamilySubnets(cidrs []ipnet.IPNet) (AddressFamilySubnets, error) {
	var addressFamilySubnets AddressFamilySubnets
	var err error

	// Split cidrs into IPv4 and IPv6 cidrs
	var ipv4Cidrs, ipv6Cidrs []ipnet.IPNet
	for _, cidr := range cidrs {
		if cidr.IP.To4() != nil {
			ipv4Cidrs = append(ipv4Cidrs, cidr)
		} else {
			ipv6Cidrs = append(ipv6Cidrs, cidr)
		}
	}
	if len(ipv4Cidrs) == 0 {
		return addressFamilySubnets, fmt.Errorf("at least one IPv4 CIDR is required")
	}

	// Split IPv4 cidrs into IPv4 subnets
	var ipv4Subnets []*net.IPNet
	switch len(ipv4Cidrs) {
	case 1:
		ipv4Subnets, err = cidr.SplitIntoSubnetsIPv4(ipv4Cidrs[0].String(), 2)
		if err != nil {
			return addressFamilySubnets, err
		}
	default:
		for _, ipv4Cidr := range ipv4Cidrs {
			ipv4Subnets = append(ipv4Subnets, &net.IPNet{
				IP:   ipv4Cidr.IP,
				Mask: ipv4Cidr.Mask,
			})
		}
	}

	// Split IPv6 cidrs into IPv6 subnets
	var ipv6Subnets []*net.IPNet
	switch len(ipv6Cidrs) {
	case 1:
		ipv6Subnets, err = cidr.SplitIntoSubnetsIPv6(ipv6Cidrs[0].String(), 2)
		if err != nil {
			return addressFamilySubnets, err
		}
	default:
		for _, ipv6Cidr := range ipv6Cidrs {
			ipv6Subnets = append(ipv6Subnets, &net.IPNet{
				IP:   ipv6Cidr.IP,
				Mask: ipv6Cidr.Mask,
			})
		}
	}
	if len(ipv6Subnets) > len(ipv4Subnets) {
		return addressFamilySubnets, fmt.Errorf("number of IPv6 CIDRs can't be more than number of IPv4 CIDRs")
	}

	for i, ipv4Subnet := range ipv4Subnets {
		addressFamilySubnet := AddressFamilySubnet{
			IPv4Subnet: ipv4Subnet,
			SubnetRole: nil,
		}
		switch i {
		case 0:
			addressFamilySubnet.SubnetRole = ptr.To(capz.SubnetControlPlane)
		case 1:
			addressFamilySubnet.SubnetRole = ptr.To(capz.SubnetNode)
		}
		if i < len(ipv6Subnets) {
			addressFamilySubnet.IPv6Subnet = ipv6Subnets[i]
			addressFamilySubnets.isDualStack = true
			addressFamilySubnets.ipv6Count++
		}
		addressFamilySubnets.addressFamilySubnets = append(addressFamilySubnets.addressFamilySubnets, addressFamilySubnet)
		addressFamilySubnets.ipv4Count++
		addressFamilySubnets.length++
	}

	return addressFamilySubnets, err
}
