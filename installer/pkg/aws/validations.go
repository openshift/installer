package aws

import (
	"errors"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// CheckKubernetesCIDRs validates an existing VPC, pod, and service CIDRs do
// not overlap.
func CheckKubernetesCIDRs(sess *session.Session, existingVPCID, podCIDR, serviceCIDR string) error {
	vpc, err := getVPC(sess, existingVPCID)
	if err != nil {
		return err
	}
	return ValidateKubernetesCIDRs(aws.StringValue(vpc.CidrBlock), podCIDR, serviceCIDR)
}

// ValidateKubernetesCIDRs validates node, pod, and service CIDRs do not
// overlap. Leave vpcCIDR blank if it is unknown (i.e. bare-metal).
func ValidateKubernetesCIDRs(vpcCIDR, podCIDR, serviceCIDR string) error {
	var cidrs []string
	if vpcCIDR == "" {
		cidrs = []string{podCIDR, serviceCIDR}
	} else {
		cidrs = []string{vpcCIDR, podCIDR, serviceCIDR}
	}
	ipnets := []*net.IPNet{}

	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return fmt.Errorf("invalid CIDR %s: %v", cidr, err)
		}
		ipnets = append(ipnets, ipnet)
	}

	// Ensure no duplicates were given.
	if err := noDuplicateIPNets(ipnets); err != nil {
		return err
	}

	// Verify that no subnets overlap.
	for _, ipnet := range ipnets {
		if err := isCIDROverlappingWith(ipnet, ipnets); err != nil {
			return fmt.Errorf("given subnets are overlapping: %v", err)
		}
	}

	return nil
}

// ValidateSubnets statically validates for correct subnet format, inclusion in
// the VPC CIDR range, and non-overlapping subnets.
func ValidateSubnets(vpcCIDR string, subnets []VPCSubnet) error {
	_, vpcNet, err := net.ParseCIDR(vpcCIDR)
	if err != nil {
		return fmt.Errorf("invalid vpcCIDR: %v", err)
	}

	// all subnets have an AZ and a valid CIDR in the VPC
	var instanceCIDRs = make([]*net.IPNet, 0)
	for i, subnet := range subnets {
		if subnet.AvailabilityZone == "" {
			return fmt.Errorf("availabilityZone must be set for subnet %s: %v", subnet.InstanceCIDR, err)
		}
		_, instanceCIDR, err := net.ParseCIDR(subnet.InstanceCIDR)
		if err != nil {
			return fmt.Errorf("invalid instanceCIDR for subnet %s: %v", subnet.InstanceCIDR, err)
		}
		instanceCIDRs = append(instanceCIDRs, instanceCIDR)
		if !vpcNet.Contains(instanceCIDR.IP) {
			return fmt.Errorf("vpcCIDR (%s) does not contain instanceCIDR (%s) for subnet #%d", vpcCIDR, instanceCIDR, i)
		}
	}

	// subnets do not overlap with one another
	for i, a := range instanceCIDRs {
		for j, b := range instanceCIDRs[:i] {
			if cidrOverlap(a, b) {
				return fmt.Errorf("CIDR of subnet %d (%s) overlaps with CIDR of subnet %d (%s)", i, a, j, b)
			}
		}
	}

	return nil
}

// CheckSubnetsAgainstExistingVPC dynamically checks that the proposed
// subnets are suitable with the given existing VPC and its subnets.
//
// Do not call this method in unit tests. It makes API requests to AWS and
// requires credentials.
func CheckSubnetsAgainstExistingVPC(sess *session.Session, existingVPCID string, controllerSubnets, workerSubnets []VPCSubnet) error {
	if existingVPCID == "" {
		return errors.New("existing VPC ID cannot be empty")
	}

	// Retrieve the existing VPC and its existing subnets.
	existingVPC, err := getVPC(sess, existingVPCID)
	if err != nil {
		return err
	}

	publicSubnets, privateSubnets, err := GetVPCSubnets(sess, existingVPCID)
	if err != nil {
		return err
	}

	// Check that the existing VPC has an available internet gateway.
	_, err = getInternetGateway(sess, existingVPCID)
	if err != nil {
		return err
	}

	err = populateCIDRs(sess, existingVPCID, controllerSubnets, workerSubnets)
	if err != nil {
		return err
	}

	return validateSubnetsAgainstExistingVPC(
		aws.StringValue(existingVPC.CidrBlock),
		publicSubnets,
		privateSubnets,
		controllerSubnets,
		workerSubnets,
	)
}

// validateSubnetAgainstExistingVPC statically validates that the proposed
// subnets are part of the VPC CIDR, do not overlap, and that if they represent
// existing subnets, they can be used for controllers or workers as proposed.
func validateSubnetsAgainstExistingVPC(existingVPCBlock string, existingPublicSubnets, existingPrivateSubnets, proposedPublicSubnets, proposedPrivateSubnets []VPCSubnet) error {
	// Parse CIDRs of every subnets.
	_, existingVPCCIDR, err := net.ParseCIDR(existingVPCBlock)
	if err != nil {
		return err
	}
	proposedPublicSubnetsCIDRs, err := vpcSubnetsToIPNets(proposedPublicSubnets)
	if err != nil {
		return err
	}
	proposedPrivateSubnetsCIDRs, err := vpcSubnetsToIPNets(proposedPrivateSubnets)
	if err != nil {
		return err
	}
	existingPublicSubnetsCIDRs, err := vpcSubnetsToIPNets(existingPublicSubnets)
	if err != nil {
		return err
	}
	existingPrivateSubnetsCIDRs, err := vpcSubnetsToIPNets(existingPrivateSubnets)
	if err != nil {
		return err
	}

	var existingSubnetsCIDRs, proposedSubnetsCIDRs []*net.IPNet
	var allProposedSubnets []VPCSubnet
	existingSubnetsCIDRs = append(existingSubnetsCIDRs, existingPublicSubnetsCIDRs...)
	existingSubnetsCIDRs = append(existingSubnetsCIDRs, existingPrivateSubnetsCIDRs...)
	proposedSubnetsCIDRs = append(proposedSubnetsCIDRs, proposedPublicSubnetsCIDRs...)
	proposedSubnetsCIDRs = append(proposedSubnetsCIDRs, proposedPrivateSubnetsCIDRs...)
	allProposedSubnets = append(allProposedSubnets, proposedPublicSubnets...)
	allProposedSubnets = append(allProposedSubnets, proposedPrivateSubnets...)

	// Verify that the proposed subnets are part of the VPC's CIDR.
	if err := containsIPNets(existingVPCCIDR, proposedSubnetsCIDRs); err != nil {
		return fmt.Errorf("given subnet is not part of the VPC network: %v", err)
	}

	// Make sure all the proposed subnets have an AZ.
	for _, subnet := range allProposedSubnets {
		if subnet.AvailabilityZone == "" {
			return fmt.Errorf("availabilityZone must be set for subnet %s", subnet.InstanceCIDR)
		}
	}

	// Ensure no duplicates were given.
	if err := noDuplicateIPNets(proposedPublicSubnetsCIDRs); err != nil {
		return err
	}

	if err := noDuplicateIPNets(proposedPrivateSubnetsCIDRs); err != nil {
		return err
	}

	// Verify that no subnets overlap.
	allSubnetCIDRs := append(proposedSubnetsCIDRs, existingSubnetsCIDRs...)
	for _, subnetCIDR := range proposedSubnetsCIDRs {
		// isCIDROverlappingWith doesn't consider strictly equal subnets to be overlapping.
		if err := isCIDROverlappingWith(subnetCIDR, allSubnetCIDRs); err != nil {
			return fmt.Errorf("given subnets are overlapping: %v", err)
		}
	}

	// Verify that the proposed subnets that already exist can be used for what
	// they are proposed for: controller or worker, depending on whether the
	// subnet is private or public.
	for _, proposedPrivateSubnetCIDR := range proposedPrivateSubnetsCIDRs {
		if containsIPNet(proposedPrivateSubnetCIDR, existingPublicSubnetsCIDRs) {
			return fmt.Errorf("subnet %s is public", proposedPrivateSubnetCIDR)
		}
	}
	return nil
}

// populateCIDRs shoves some CIDRs into subnets when we know the IDs
func populateCIDRs(sess *session.Session, existingVPCID string, publicSubnets, privateSubnets []VPCSubnet) error {
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

// vpcSubnetsToIPNets returns a slice of *net.IPNet containing the given
// VPCSubnets.
func vpcSubnetsToIPNets(subnets []VPCSubnet) ([]*net.IPNet, error) {
	subnetCIDRs := make([]*net.IPNet, 0, len(subnets))
	for _, subnet := range subnets {
		_, subnetCIDR, err := net.ParseCIDR(subnet.InstanceCIDR)
		if err != nil {
			return subnetCIDRs, fmt.Errorf("invalid CIDR for subnet %s: %v", subnet.InstanceCIDR, err)
		}
		subnetCIDRs = append(subnetCIDRs, subnetCIDR)
	}
	return subnetCIDRs, nil
}

// Does the address space of these networks "a" and "b" overlap?
func cidrOverlap(a, b *net.IPNet) bool {
	return a.Contains(b.IP) || b.Contains(a.IP)
}

// containsIPNets returns an error if the IPNets are outside of the
// the given CIDR.
func containsIPNets(cidr *net.IPNet, ipNets []*net.IPNet) error {
	for _, ipNet := range ipNets {
		if !cidr.Contains(ipNet.IP) {
			return fmt.Errorf("%s is not a subnet of %s", ipNet.String(), cidr.String())
		}
	}

	return nil
}

// noDuplicateIPNets returns an error if the given IPNet list contains the same
// element twice.
func noDuplicateIPNets(ipNets []*net.IPNet) error {
	m := make(map[string]struct{})
	for _, subnet := range ipNets {
		if _, exists := m[subnet.String()]; exists {
			return fmt.Errorf("IP Network %q has been given twice", subnet)
		}
		m[subnet.String()] = struct{}{}
	}
	return nil
}

// isCIDROverlappingWith returns an error if the given CIDR overlaps with
// any of the other IPNets.
//
// This function ignores strictly duplicated IPNets.
func isCIDROverlappingWith(cidr *net.IPNet, ipNets []*net.IPNet) error {
	for _, ipNet := range ipNets {
		if cidr.String() == ipNet.String() {
			continue
		}
		if cidrOverlap(cidr, ipNet) {
			return fmt.Errorf("IP Network %s conflicts with %s", cidr, ipNet)
		}
	}
	return nil
}

// containsIPNet returns true if the given IPNet exists in the list
// of other CIDRs.
func containsIPNet(subnetCIDR *net.IPNet, otherSubnetCIDRs []*net.IPNet) bool {
	for _, otherSubnetCIDR := range otherSubnetCIDRs {
		if subnetCIDR.String() == otherSubnetCIDR.String() {
			return true
		}
	}
	return false
}
