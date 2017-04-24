package cloudforms

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/coreos/ipnets"
)

// GetDefaultSubnets partitions a CIDR into subnets
func GetDefaultSubnets(sess *session.Session, vpcCIDR string) ([]VPCSubnet, []VPCSubnet, error) {
	zones, err := getAvailabilityZones(sess)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting availability zone %v", err)
	}

	_, vpcNet, err := net.ParseCIDR(vpcCIDR)
	if vpcNet == nil || err != nil {
		return nil, nil, fmt.Errorf("failed parsing VPC CIDR %v", err)
	}

	// Calculate subnetMultipler many times as many subnets as needed to
	// intentionally leave unused IPs for unspecified use. A multipler
	// of 1 divides the VPC among AZs, 2 leaves 50% of the VPC unallocated,
	// 4 leaves 75% unallocated, etc.
	cidrs, err := ipnets.SubnetInto(vpcNet, 2*2*len(zones))
	if err != nil {
		return nil, nil, fmt.Errorf("failed dividing VPC into subnets %v", err)
	}

	controllerSubnets := make([]VPCSubnet, len(zones))
	workerSubnets := make([]VPCSubnet, len(zones))

	// add generated multi-AZ subnets for controllers
	for i, zone := range zones {
		controllerSubnets[i] = VPCSubnet{
			AvailabilityZone: zone,
			InstanceCIDR:     cidrs[i].String(),
		}
	}

	// add generated multi-AZ subnets for workers
	for i, zone := range zones {
		workerSubnets[i] = VPCSubnet{
			AvailabilityZone: zone,
			InstanceCIDR:     cidrs[i+len(zones)].String(),
		}
	}

	return controllerSubnets, workerSubnets, nil
}

// PopulateCIDRs shoves some CIDRs into subnets when we know the IDs
func PopulateCIDRs(sess *session.Session, existingVPCID string, publicSubnets, privateSubnets []VPCSubnet) error {
	existingPublicSubnets, existingPrivateSubnets, err := GetVPCSubnets(sess, existingVPCID)
	if err != nil {
		return err
	}

	existingSubnets := append(existingPublicSubnets, existingPrivateSubnets...)
	for i, subnet := range publicSubnets {
		if subnet.ID == "" || subnet.InstanceCIDR != "" {
			continue
		}
		for _, existing := range existingSubnets {
			if subnet.ID == existing.ID {
				publicSubnets[i].InstanceCIDR = existing.InstanceCIDR
				break
			}
		}
	}
	for i, subnet := range privateSubnets {
		if subnet.ID == "" || subnet.InstanceCIDR != "" {
			continue
		}
		for _, existing := range existingSubnets {
			if subnet.ID == existing.ID {
				privateSubnets[i].InstanceCIDR = existing.InstanceCIDR
				break
			}
		}
	}
	return nil
}

// Does the address space of these networks "a" and "b" overlap?
func cidrOverlap(a, b *net.IPNet) bool {
	return a.Contains(b.IP) || b.Contains(a.IP)
}
