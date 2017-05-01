package cloudforms

import (
	"testing"
)

// Modify the config, return bool indicating whether or not it's valid
type configTwiddler func(c *Config) bool

func testConfig(t *testing.T, description string, twiddler configTwiddler) {
	// TODO(chom): move createRecordSet/hostedZoneId to test cases when that behavior is configurable
	c := &Config{
		Channel:          "alpha",
		ClusterName:      "unit-test-config",
		ELBScheme:        "internet-facing",
		ControllerDomain: "cluster-k8s.staging.core-os.net",
		TectonicDomain:   "cluster.staging.core-os.net",
		KeyName:          "test-key-name",
		Region:           "us-west-1",
		KMSKeyARN:        "arn:aws:kms:us-west-1:xxxxxxxxx:key/xxxxxxxxxxxxxxxxxxx",
		HostedZoneID:     "XXXXXXXXXXX",
		VPCCIDR:          "10.4.0.0/16",
		WorkerSubnets: []VPCSubnet{
			{
				InstanceCIDR:     "10.4.3.0/24",
				AvailabilityZone: "us-west-1c",
			},
			{
				InstanceCIDR:     "10.4.4.0/24",
				AvailabilityZone: "us-west-1a",
			},
		},
		PodCIDR:     "172.4.0.0/16",
		ServiceCIDR: "172.5.0.0/24",
	}
	shouldBeValid := twiddler(c)

	c.SetDefaults()
	err := c.Valid()

	if err != nil && shouldBeValid {
		t.Errorf("%s: unexpected error validating config: %v", description, err)
	} else if err == nil && !shouldBeValid {
		t.Errorf("%s: expected error validating config, got none", description)
	}
}

func TestNetworkConfig_CorrectParams(t *testing.T) {
	testConfig(t, "create vpc", func(c *Config) bool {
		return true
	})

	testConfig(t, "existing vpc, default route table", func(c *Config) bool {
		c.VPCID = "vpc-xxxxx"
		return true
	})

	testConfig(t, "existing vpc, specify route table", func(c *Config) bool {
		c.VPCID = "vpc-xxxxx"
		c.RouteTableID = "rtb-xxxxx"
		return true
	})
}

func TestNetworkConfig_BadAddressLayout(t *testing.T) {
	testConfig(t, "PodCIDR intersects VPCCIDR", func(c *Config) bool {
		c.PodCIDR = "10.4.100.0/20"
		return false
	})

	testConfig(t, "ServiceCIDR intersects VPCCIDR", func(c *Config) bool {
		c.ServiceCIDR = "10.4.100.0/24"
		return false
	})

	testConfig(t, "ServiceCIDR intersects PodCIDR", func(c *Config) bool {
		c.ServiceCIDR = "172.4.10.0/24"
		return false
	})
}

func TestNetworkConfig_BadExistingVPCDefinition(t *testing.T) {
	testConfig(t, "RouteTableID specified without VPCID", func(c *Config) bool {
		c.RouteTableID = "rtb-xxxxx"
		return false
	})
}
