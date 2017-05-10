package aws

import (
	"testing"
)

func TestValidateKubernetesCIDRs(t *testing.T) {
	cidrs := []struct {
		VPCCIDR     string
		PodCIDR     string
		ServiceCIDR string
	}{
		{"10.0.0.0/16", "10.2.0.0/16", "10.3.0.0/24"},
		{"10.0.0.0/16", "10.2.0.0/16", "10.3.0.0/16"},
		// bare-metal
		{"", "10.2.0.0/16", "10.3.0.0/24"},
		{"", "10.2.0.0/16", "10.3.0.0/16"},
		{"", "192.168.1.0/24", "192.168.2.0/24"},
	}

	for _, c := range cidrs {
		if err := ValidateKubernetesCIDRs(c.VPCCIDR, c.PodCIDR, c.ServiceCIDR); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}
}

func TestValidateKubernetesCIDRs_Invalid(t *testing.T) {
	cidrs := []struct {
		VPCCIDR     string
		PodCIDR     string
		ServiceCIDR string
	}{
		// CIDR parsing
		{"10.0.0.0", "10.2.0.0", "10.3.0.0"},
		// Conflicts
		{"10.1.0.0/16", "10.0.2.0/16", "10.0.3.0/16"},
		{"10.1.0.0/16", "10.1.0.0/16", "10.1.0.0/16"},
		{"", "10.2.0.0/16", "10.2.0.0/24"},
	}

	for _, c := range cidrs {
		if err := ValidateKubernetesCIDRs(c.VPCCIDR, c.PodCIDR, c.ServiceCIDR); err == nil {
			t.Errorf("expected error using vpc CIDR %s pod CIDR %s service CIDR %s, got %v", c.VPCCIDR, c.PodCIDR, c.ServiceCIDR, err)
		}
	}
}

func TestValidateSubnets(t *testing.T) {
	correctVPCCIDR := "10.0.0.0/16"
	subnets := []VPCSubnet{
		VPCSubnet{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.0.0/20",
		},
		VPCSubnet{
			AvailabilityZone: "us-west-2b",
			InstanceCIDR:     "10.0.16.0/20",
		},
		VPCSubnet{
			AvailabilityZone: "us-west-2c",
			InstanceCIDR:     "10.0.32.0/20",
		},
		VPCSubnet{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.48.0/20",
		},
		VPCSubnet{
			AvailabilityZone: "us-west-2b",
			InstanceCIDR:     "10.0.64.0/20",
		},
		VPCSubnet{
			AvailabilityZone: "us-west-2c",
			InstanceCIDR:     "10.0.80.0/20",
		},
	}

	if err := ValidateSubnets(correctVPCCIDR, subnets); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateSubnets_InvalidVPCCIDR(t *testing.T) {
	notCIDR := "10.0.0.0/99"
	subnets := []VPCSubnet{}

	if err := ValidateSubnets(notCIDR, subnets); err == nil {
		t.Errorf("expected invalid CIDR error, got %v", err)
	}
}

func TestValidateSubnets_MissingAvailabilityZone(t *testing.T) {
	vpcCIDR := "10.0.0.0/16"
	subnets := []VPCSubnet{
		VPCSubnet{
			AvailabilityZone: "",
			InstanceCIDR:     "10.0.0.0/20",
		},
	}

	if err := ValidateSubnets(vpcCIDR, subnets); err == nil {
		t.Errorf("expected missing Availability zone error, got %v", err)
	}
}

func TestValidateSubnets_InvalidSubnetCIDR(t *testing.T) {
	vpcCIDR := "10.0.0.0/16"
	subnets := []VPCSubnet{
		VPCSubnet{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.0.0/99",
		},
	}

	if err := ValidateSubnets(vpcCIDR, subnets); err == nil {
		t.Errorf("expected missing Availability zone error, got %v", err)
	}
}

func TestValidateSubnets_SubnetOverlaps(t *testing.T) {
	vpcCIDR := "10.0.0.0/16"
	subnets := []VPCSubnet{
		VPCSubnet{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.0.0/20",
		},
		VPCSubnet{
			AvailabilityZone: "us-west-2b",
			InstanceCIDR:     "10.0.0.0/20",
		},
	}

	if err := ValidateSubnets(vpcCIDR, subnets); err == nil {
		t.Errorf("expected subnet overlap error, got %v", err)
	}
}

func TestValidateSubnetsAgainstExistingVPC(t *testing.T) {
	existingVPCBlock := "10.0.0.0/16"
	proposedControllerSubnets := []VPCSubnet{
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.16.0/24",
		},
		{
			AvailabilityZone: "us-west-2b",
			InstanceCIDR:     "10.0.18.0/24",
		},
	}
	proposedWorkerSubnets := []VPCSubnet{
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.20.0/24",
		},
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.22.0/24",
		},
	}
	existingControllerSubnets := []VPCSubnet{}
	existingWorkerSubnets := []VPCSubnet{
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.20.0/24",
		},
	}

	err := validateSubnetsAgainstExistingVPC(
		existingVPCBlock,
		existingControllerSubnets,
		existingWorkerSubnets,
		proposedControllerSubnets,
		proposedWorkerSubnets,
	)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateSubnetsAgainstExistingVPC_MissingAvailabilityZone(t *testing.T) {
	existingVPCBlock := "10.0.0.0/16"
	proposedControllerSubnets := []VPCSubnet{
		{
			AvailabilityZone: "",
			InstanceCIDR:     "10.0.2.0/24",
		},
	}
	proposedWorkerSubnets := []VPCSubnet{}
	existingControllerSubnets := []VPCSubnet{}
	existingWorkerSubnets := []VPCSubnet{}

	err := validateSubnetsAgainstExistingVPC(
		existingVPCBlock,
		existingControllerSubnets,
		existingWorkerSubnets,
		proposedControllerSubnets,
		proposedWorkerSubnets,
	)
	if err == nil {
		t.Errorf("expected missing availability zone error, got %v", err)
	}
}

func TestValidateSubnetsAgainstExistingVPC_InvalidCIDR(t *testing.T) {
	existingVPCBlock := "10.0.0.0/16"
	proposedControllerSubnets := []VPCSubnet{
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "invalid-cidr",
		},
	}
	proposedWorkerSubnets := []VPCSubnet{}
	existingControllerSubnets := []VPCSubnet{}
	existingWorkerSubnets := []VPCSubnet{}

	err := validateSubnetsAgainstExistingVPC(
		existingVPCBlock,
		existingControllerSubnets,
		existingWorkerSubnets,
		proposedControllerSubnets,
		proposedWorkerSubnets,
	)
	if err == nil {
		t.Errorf("expected invalid proposed subnet CIDR error, got %v", err)
	}
}

func TestValidateSubnetsAgainstExistingVPC_OutsideVPC(t *testing.T) {
	existingVPCBlock := "10.0.0.0/16"
	proposedControllerSubnets := []VPCSubnet{
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "192.168.0.0/24",
		},
	}
	proposedWorkerSubnets := []VPCSubnet{}
	existingControllerSubnets := []VPCSubnet{}
	existingWorkerSubnets := []VPCSubnet{}

	err := validateSubnetsAgainstExistingVPC(
		existingVPCBlock,
		existingControllerSubnets,
		existingWorkerSubnets,
		proposedControllerSubnets,
		proposedWorkerSubnets,
	)
	if err == nil {
		t.Errorf("expected existing VPC does not contain proposed subnet error, got %v", err)
	}
}

func TestValidateSubnetsAgainstExistingVPC_Duplicated(t *testing.T) {
	existingVPCBlock := "10.0.0.0/16"
	proposedControllerSubnets := []VPCSubnet{
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.1.0/24",
		},
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.1.5/24",
		},
	}
	proposedWorkerSubnets := []VPCSubnet{}
	existingControllerSubnets := []VPCSubnet{}
	existingWorkerSubnets := []VPCSubnet{}

	err := validateSubnetsAgainstExistingVPC(
		existingVPCBlock,
		existingControllerSubnets,
		existingWorkerSubnets,
		proposedControllerSubnets,
		proposedWorkerSubnets,
	)
	if err == nil {
		t.Errorf("expected duplicated subnet error, got %v", err)
	}
}

func TestValidateSubnetsAgainstExistingVPC_Overlap(t *testing.T) {
	existingVPCBlock := "10.0.0.0/16"
	proposedControllerSubnets := []VPCSubnet{
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.0.0/20",
		},
		{
			AvailabilityZone: "us-west-2a",
			InstanceCIDR:     "10.0.0.0/18",
		},
	}
	proposedWorkerSubnets := []VPCSubnet{}
	existingControllerSubnets := []VPCSubnet{}
	existingWorkerSubnets := []VPCSubnet{}

	err := validateSubnetsAgainstExistingVPC(
		existingVPCBlock,
		existingControllerSubnets,
		existingWorkerSubnets,
		proposedControllerSubnets,
		proposedWorkerSubnets,
	)
	if err == nil {
		t.Errorf("expected duplicated subnet error, got %v", err)
	}
}
