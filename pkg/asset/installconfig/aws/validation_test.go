package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

var (
	validCIDR = "10.0.0.0/16"
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			AWS: &aws.Platform{
				Region: "us-east-1",
				Subnets: []string{
					"valid-private-subnet-a",
					"valid-private-subnet-b",
					"valid-private-subnet-c",
					"valid-public-subnet-a",
					"valid-public-subnet-b",
					"valid-public-subnet-c",
				},
			},
		},
		ControlPlane: &types.MachinePool{
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{
					Zones: []string{"a", "b", "c"},
				},
			},
		},
		Compute: []types.MachinePool{{
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{
					Zones: []string{"a", "b", "c"},
				},
			},
		}},
	}
}

func validAvailZones() []string {
	return []string{"a", "b", "c"}
}

func validPrivateSubnets() map[string]Subnet {
	return map[string]Subnet{
		"valid-private-subnet-a": {
			Zone: "a",
			CIDR: "10.0.1.0/24",
		},
		"valid-private-subnet-b": {
			Zone: "b",
			CIDR: "10.0.2.0/24",
		},
		"valid-private-subnet-c": {
			Zone: "c",
			CIDR: "10.0.3.0/24",
		},
	}
}

func validPublicSubnets() map[string]Subnet {
	return map[string]Subnet{
		"valid-public-subnet-a": {
			Zone: "a",
			CIDR: "10.0.4.0/24",
		},
		"valid-public-subnet-b": {
			Zone: "b",
			CIDR: "10.0.5.0/24",
		},
		"valid-public-subnet-c": {
			Zone: "c",
			CIDR: "10.0.6.0/24",
		},
	}
}

func validServiceEndpoints() []aws.ServiceEndpoint {
	return []aws.ServiceEndpoint{{
		Name: "ec2",
		URL:  "e2e.local",
	}, {
		Name: "s3",
		URL:  "e2e.local",
	}, {
		Name: "iam",
		URL:  "e2e.local",
	}, {
		Name: "elasticloadbalancing",
		URL:  "e2e.local",
	}, {
		Name: "tagging",
		URL:  "e2e.local",
	}, {
		Name: "route53",
		URL:  "e2e.local",
	}, {
		Name: "sts",
		URL:  "e2e.local",
	}}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		installConfig  *types.InstallConfig
		availZones     []string
		privateSubnets map[string]Subnet
		publicSubnets  map[string]Subnet
		exptectErr     string
	}{{
		name: "valid no byo",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS = &aws.Platform{Region: "us-east-1"}
			return c
		}(),
		availZones: validAvailZones(),
	}, {
		name: "valid no byo",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = nil
			return c
		}(),
		availZones: validAvailZones(),
	}, {
		name: "valid no byo",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = []string{}
			return c
		}(),
		availZones: validAvailZones(),
	}, {
		name:           "valid byo",
		installConfig:  validInstallConfig(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
	}, {
		name: "valid byo",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Publish = types.InternalPublishingStrategy
			c.Platform.AWS.Subnets = []string{
				"valid-private-subnet-a",
				"valid-private-subnet-b",
				"valid-private-subnet-c",
			}
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
	}, {
		name: "invalid no private subnets",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = []string{
				"valid-public-subnet-a",
				"valid-public-subnet-b",
				"valid-public-subnet-c",
			}
			return c
		}(),
		availZones:    validAvailZones(),
		publicSubnets: validPublicSubnets(),
		exptectErr:    `^\[platform\.aws\.subnets: Invalid value: \[\]string{\"valid-public-subnet-a\", \"valid-public-subnet-b\", \"valid-public-subnet-c\"}: No private subnets found, controlPlane\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\"}: No subnets provided for zones \[a b c\], compute\[0\]\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\"}: No subnets provided for zones \[a b c\]\]$`,
	}, {
		name: "invalid no public subnets",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = []string{
				"valid-private-subnet-a",
				"valid-private-subnet-b",
				"valid-private-subnet-c",
			}
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		exptectErr:     `^platform\.aws\.subnets: Invalid value: \[\]string{\"valid-private-subnet-a\", \"valid-private-subnet-b\", \"valid-private-subnet-c\"}: No public subnet provided for zones \[a b c\]$`,
	}, {
		name: "invalid cidr does not belong to machine CIDR",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = append(c.Platform.AWS.Subnets, "invalid-cidr-subnet")
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets: func() map[string]Subnet {
			s := validPublicSubnets()
			s["invalid-cidr-subnet"] = Subnet{
				CIDR: "192.168.126.0/24",
			}
			return s
		}(),
		exptectErr: `^platform\.aws\.subnets\[6\]: Invalid value: \"invalid-cidr-subnet\": subnet's CIDR range start 192.168.126.0 is outside of the specified machine networks$`,
	}, {
		name: "invalid cidr does not belong to machine CIDR",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = append(c.Platform.AWS.Subnets, "invalid-private-cidr-subnet", "invalid-public-cidr-subnet")
			return c
		}(),
		availZones: validAvailZones(),
		privateSubnets: func() map[string]Subnet {
			s := validPrivateSubnets()
			s["invalid-private-cidr-subnet"] = Subnet{
				CIDR: "192.168.126.0/24",
			}
			return s
		}(),
		publicSubnets: func() map[string]Subnet {
			s := validPublicSubnets()
			s["invalid-public-cidr-subnet"] = Subnet{
				CIDR: "192.168.127.0/24",
			}
			return s
		}(),
		exptectErr: `^\[platform\.aws\.subnets\[6\]: Invalid value: \"invalid-private-cidr-subnet\": subnet's CIDR range start 192.168.126.0 is outside of the specified machine networks, platform\.aws\.subnets\[7\]: Invalid value: \"invalid-public-cidr-subnet\": subnet's CIDR range start 192.168.127.0 is outside of the specified machine networks\]$`,
	}, {
		name: "invalid missing public subnet in a zone",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = append(c.Platform.AWS.Subnets, "no-matching-public-private-zone")
			return c
		}(),
		availZones: validAvailZones(),
		privateSubnets: func() map[string]Subnet {
			s := validPrivateSubnets()
			s["no-matching-public-private-zone"] = Subnet{
				Zone: "f",
				CIDR: "10.0.7.0/24",
			}
			return s
		}(),
		publicSubnets: validPublicSubnets(),
		exptectErr:    `^platform\.aws\.subnets: Invalid value: \[\]string{\"valid-private-subnet-a\", \"valid-private-subnet-b\", \"valid-private-subnet-c\", \"valid-public-subnet-a\", \"valid-public-subnet-b\", \"valid-public-subnet-c\", \"no-matching-public-private-zone\"}: No public subnet provided for zones \[f\]$`,
	}, {
		name: "invalid multiple private in same zone",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = append(c.Platform.AWS.Subnets, "valid-private-zone-c-2")
			return c
		}(),
		availZones: validAvailZones(),
		privateSubnets: func() map[string]Subnet {
			s := validPrivateSubnets()
			s["valid-private-zone-c-2"] = Subnet{
				Zone: "c",
				CIDR: "10.0.7.0/24",
			}
			return s
		}(),
		publicSubnets: validPublicSubnets(),
		exptectErr:    `^platform\.aws\.subnets\[6\]: Invalid value: \"valid-private-zone-c-2\": private subnet valid-private-subnet-c is also in zone c$`,
	}, {
		name: "invalid multiple public in same zone",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Subnets = append(c.Platform.AWS.Subnets, "valid-public-zone-c-2")
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets: func() map[string]Subnet {
			s := validPublicSubnets()
			s["valid-public-zone-c-2"] = Subnet{
				Zone: "c",
				CIDR: "10.0.7.0/24",
			}
			return s
		}(),
		exptectErr: `^platform\.aws\.subnets\[6\]: Invalid value: \"valid-public-zone-c-2\": public subnet valid-public-subnet-c is also in zone c$`,
	}, {
		name: "invalid no subnet for control plane zones",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.ControlPlane.Platform.AWS.Zones = append(c.ControlPlane.Platform.AWS.Zones, "d")
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		exptectErr:     `^controlPlane\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\", \"d\"}: No subnets provided for zones \[d\]$`,
	}, {
		name: "invalid no subnet for control plane zones",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.ControlPlane.Platform.AWS.Zones = append(c.ControlPlane.Platform.AWS.Zones, "d", "e")
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		exptectErr:     `^controlPlane\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\", \"d\", \"e\"}: No subnets provided for zones \[d e\]$`,
	}, {
		name: "invalid no subnet for compute[0] zones",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Compute[0].Platform.AWS.Zones = append(c.ControlPlane.Platform.AWS.Zones, "d")
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		exptectErr:     `^compute\[0\]\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\", \"d\"}: No subnets provided for zones \[d\]$`,
	}, {
		name: "invalid no subnet for compute zone",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Compute[0].Platform.AWS.Zones = append(c.ControlPlane.Platform.AWS.Zones, "d")
			c.Compute = append(c.Compute, types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						Zones: []string{"a", "b", "e"},
					},
				},
			})
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		exptectErr:     `^\[compute\[0\]\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\", \"d\"}: No subnets provided for zones \[d\], compute\[1\]\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"e\"}: No subnets provided for zones \[e\]\]$`,
	}, {
		name: "custom region invalid service endpoints none provided",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "test-region"
			c.Platform.AWS.AMIID = "dummy-id"
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		exptectErr:     `^platform\.aws\.serviceEndpoints: Invalid value: (.|\n)*: \[failed to find endpoint for service "ec2": (.|\n)*, failed to find endpoint for service "elasticloadbalancing": (.|\n)*, failed to find endpoint for service "iam": (.|\n)*, failed to find endpoint for service "route53": (.|\n)*, failed to find endpoint for service "s3": (.|\n)*, failed to find endpoint for service "sts": (.|\n)*, failed to find endpoint for service "tagging": (.|\n)*\]$`,
	}, {
		name: "custom region invalid service endpoints some provided",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "test-region"
			c.Platform.AWS.AMIID = "dummy-id"
			c.Platform.AWS.ServiceEndpoints = validServiceEndpoints()[:3]
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		exptectErr:     `^platform\.aws\.serviceEndpoints: Invalid value: (.|\n)*: \[failed to find endpoint for service "elasticloadbalancing": (.|\n)*, failed to find endpoint for service "route53": (.|\n)*, failed to find endpoint for service "sts": (.|\n)*, failed to find endpoint for service "tagging": (.|\n)*$`,
	}, {
		name: "custom region valid service endpoints",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "test-region"
			c.Platform.AWS.AMIID = "dummy-id"
			c.Platform.AWS.ServiceEndpoints = validServiceEndpoints()
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
	}, {
		name: "AMI not provided for unknown region",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "test-region"
			c.Platform.AWS.ServiceEndpoints = validServiceEndpoints()
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		exptectErr:     `^platform\.aws\.amiID: Required value: AMI must be provided$`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			meta := &Metadata{
				availabilityZones: test.availZones,
				privateSubnets:    test.privateSubnets,
				publicSubnets:     test.publicSubnets,
			}
			err := Validate(context.TODO(), meta, test.installConfig)
			if test.exptectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.exptectErr, err.Error())
			}
		})
	}
}
