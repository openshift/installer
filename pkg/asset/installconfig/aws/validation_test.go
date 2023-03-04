package aws

import (
	"context"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/asset/installconfig/aws/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

var (
	validCIDR             = "10.0.0.0/16"
	validRegion           = "us-east-1"
	validCallerRef        = "valid-caller-reference"
	validDSId             = "valid-delegation-set-id"
	validNameServers      = []string{"valid-name-server"}
	validHostedZoneName   = "valid-private-subnet-a"
	invalidHostedZoneName = "invalid-hosted-zone"
	validDomainName       = "valid-base-domain"
	invalidBaseDomain     = "invalid-base-domain"
	metaName              = "ClusterMetaName"

	publishInternal      = func(ic *types.InstallConfig) { ic.Publish = types.InternalPublishingStrategy }
	clearHostedZone      = func(ic *types.InstallConfig) { ic.AWS.HostedZone = "" }
	invalidateHostedZone = func(ic *types.InstallConfig) { ic.AWS.HostedZone = invalidHostedZoneName }
	invalidateBaseDomain = func(ic *types.InstallConfig) { ic.BaseDomain = invalidBaseDomain }
	clearBaseDomain      = func(ic *types.InstallConfig) { ic.BaseDomain = "" }
	invalidateRegion     = func(ic *types.InstallConfig) { ic.AWS.Region = "us-east4" }
)

type editFunctions []func(ic *types.InstallConfig)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		BaseDomain: validDomainName,
		Publish:    types.ExternalPublishingStrategy,
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
				HostedZone: validHostedZoneName,
			},
		},
		ControlPlane: &types.MachinePool{
			Architecture: types.ArchitectureAMD64,
			Replicas:     pointer.Int64Ptr(3),
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{
					Zones: []string{"a", "b", "c"},
				},
			},
		},
		Compute: []types.MachinePool{{
			Name:         types.MachinePoolComputeRoleName,
			Architecture: types.ArchitectureAMD64,
			Replicas:     pointer.Int64Ptr(3),
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{
					Zones: []string{"a", "b", "c"},
				},
			},
		}},
		ObjectMeta: metav1.ObjectMeta{
			Name: metaName,
		},
	}
}

func validInstallConfigEdge() *types.InstallConfig {
	ic := validInstallConfig()
	edgeSubnets := validEdgeSubnets()
	for subnet := range edgeSubnets {
		ic.Platform.AWS.Subnets = append(ic.Platform.AWS.Subnets, subnet)
	}
	ic.Compute = append(ic.Compute, types.MachinePool{
		Name: types.MachinePoolEdgeRoleName,
		Platform: types.MachinePoolPlatform{
			AWS: &aws.MachinePool{},
		},
	})
	return ic
}

func validAvailZones() []string {
	return []string{"a", "b", "c"}
}

func validAvailZonesWithEdge() []string {
	return []string{"a", "b", "c", "edge-a", "edge-b", "edge-c"}
}

func validAvailZonesOnlyEdge() []string {
	return []string{"edge-a", "edge-b", "edge-c"}
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

func validEdgeSubnets() map[string]Subnet {
	return map[string]Subnet{
		"valid-public-subnet-edge-a": {
			Zone: "edge-a",
			CIDR: "10.0.7.0/24",
		},
		"valid-public-subnet-edge-b": {
			Zone: "edge-b",
			CIDR: "10.0.8.0/24",
		},
		"valid-public-subnet-edge-c": {
			Zone: "edge-c",
			CIDR: "10.0.9.0/24",
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

func invalidServiceEndpoint() []aws.ServiceEndpoint {
	return []aws.ServiceEndpoint{{
		Name: "testing",
		URL:  "testing",
	}, {
		Name: "test",
		URL:  "http://testing.non",
	}}
}

func validInstanceTypes() map[string]InstanceType {
	return map[string]InstanceType{
		"t2.small": {
			DefaultVCpus: 1,
			MemInMiB:     2048,
		},
		"m5.large": {
			DefaultVCpus: 2,
			MemInMiB:     8192,
		},
		"m5.xlarge": {
			DefaultVCpus: 4,
			MemInMiB:     16384,
		},
	}
}

func createBaseDomainHostedZone() route53.HostedZone {
	return route53.HostedZone{
		CallerReference: &validCallerRef,
		Id:              &validDSId,
		Name:            &validDomainName,
	}
}

func createValidHostedZone() route53.GetHostedZoneOutput {
	ptrValidNameServers := []*string{}
	for i := range validNameServers {
		ptrValidNameServers = append(ptrValidNameServers, &validNameServers[i])
	}

	validDelegationSet := route53.DelegationSet{CallerReference: &validCallerRef, Id: &validDSId, NameServers: ptrValidNameServers}
	validHostedZone := route53.HostedZone{CallerReference: &validCallerRef, Id: &validDSId, Name: &validHostedZoneName}
	validVPCs := []*route53.VPC{{VPCId: &validHostedZoneName, VPCRegion: &validRegion}}

	return route53.GetHostedZoneOutput{
		DelegationSet: &validDelegationSet,
		HostedZone:    &validHostedZone,
		VPCs:          validVPCs,
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		installConfig  *types.InstallConfig
		availZones     []string
		privateSubnets map[string]Subnet
		publicSubnets  map[string]Subnet
		edgeSubnets    map[string]Subnet
		instanceTypes  map[string]InstanceType
		proxy          string
		expectErr      string
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
		name:           "valid byo",
		installConfig:  validInstallConfigEdge(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		edgeSubnets:    validEdgeSubnets(),
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
		name: "valid instance types",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS = &aws.Platform{
				Region: "us-east-1",
				DefaultMachinePlatform: &aws.MachinePool{
					InstanceType: "m5.xlarge",
				},
			}
			c.ControlPlane.Platform.AWS.InstanceType = "m5.xlarge"
			c.Compute[0].Platform.AWS.InstanceType = "m5.large"
			return c
		}(),
		availZones:    validAvailZones(),
		instanceTypes: validInstanceTypes(),
	}, {
		name: "invalid control plane instance type",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS = &aws.Platform{Region: "us-east-1"}
			c.ControlPlane.Platform.AWS.InstanceType = "t2.small"
			c.Compute[0].Platform.AWS.InstanceType = "m5.large"
			return c
		}(),
		availZones:    validAvailZones(),
		instanceTypes: validInstanceTypes(),
		expectErr:     `^\Q[controlPlane.platform.aws.type: Invalid value: "t2.small": instance type does not meet minimum resource requirements of 4 vCPUs, controlPlane.platform.aws.type: Invalid value: "t2.small": instance type does not meet minimum resource requirements of 16384 MiB Memory]\E$`,
	}, {
		name: "invalid compute instance type",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS = &aws.Platform{Region: "us-east-1"}
			c.ControlPlane.Platform.AWS.InstanceType = "m5.xlarge"
			c.Compute[0].Platform.AWS.InstanceType = "t2.small"
			return c
		}(),
		availZones:    validAvailZones(),
		instanceTypes: validInstanceTypes(),
		expectErr:     `^\Q[compute[0].platform.aws.type: Invalid value: "t2.small": instance type does not meet minimum resource requirements of 2 vCPUs, compute[0].platform.aws.type: Invalid value: "t2.small": instance type does not meet minimum resource requirements of 8192 MiB Memory]\E$`,
	}, {
		name: "undefined compute instance type",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS = &aws.Platform{Region: "us-east-1"}
			c.Compute[0].Platform.AWS.InstanceType = "m5.dummy"
			return c
		}(),
		availZones:    validAvailZones(),
		instanceTypes: validInstanceTypes(),
		expectErr:     `^\Qcompute[0].platform.aws.type: Invalid value: "m5.dummy": instance type m5.dummy not found\E$`,
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
		expectErr:     `^\[platform\.aws\.subnets: Invalid value: \[\]string{\"valid-public-subnet-a\", \"valid-public-subnet-b\", \"valid-public-subnet-c\"}: No private subnets found, controlPlane\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\"}: No subnets provided for zones \[a b c\], compute\[0\]\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\"}: No subnets provided for zones \[a b c\]\]$`,
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
		expectErr:      `^platform\.aws\.subnets: Invalid value: \[\]string{\"valid-private-subnet-a\", \"valid-private-subnet-b\", \"valid-private-subnet-c\"}: No public subnet provided for zones \[a b c\]$`,
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
		expectErr: `^platform\.aws\.subnets\[6\]: Invalid value: \"invalid-cidr-subnet\": subnet's CIDR range start 192.168.126.0 is outside of the specified machine networks$`,
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
		expectErr: `^\[platform\.aws\.subnets\[6\]: Invalid value: \"invalid-private-cidr-subnet\": subnet's CIDR range start 192.168.126.0 is outside of the specified machine networks, platform\.aws\.subnets\[7\]: Invalid value: \"invalid-public-cidr-subnet\": subnet's CIDR range start 192.168.127.0 is outside of the specified machine networks\]$`,
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
		expectErr:     `^platform\.aws\.subnets: Invalid value: \[\]string{\"valid-private-subnet-a\", \"valid-private-subnet-b\", \"valid-private-subnet-c\", \"valid-public-subnet-a\", \"valid-public-subnet-b\", \"valid-public-subnet-c\", \"no-matching-public-private-zone\"}: No public subnet provided for zones \[f\]$`,
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
		expectErr:     `^platform\.aws\.subnets\[6\]: Invalid value: \"valid-private-zone-c-2\": private subnet valid-private-subnet-c is also in zone c$`,
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
		expectErr: `^platform\.aws\.subnets\[6\]: Invalid value: \"valid-public-zone-c-2\": public subnet valid-public-subnet-c is also in zone c$`,
	}, {
		name: "invalid multiple public edge in same zone",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfigEdge()
			c.Platform.AWS.Subnets = append(c.Platform.AWS.Subnets, "valid-public-zone-edge-c-2")
			return c
		}(),
		availZones:     validAvailZonesWithEdge(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		edgeSubnets: func() map[string]Subnet {
			s := validEdgeSubnets()
			s["valid-public-zone-edge-c-2"] = Subnet{
				Zone:     "edge-c",
				CIDR:     "10.0.9.0/24",
				ZoneType: aws.LocalZoneType,
			}
			return s
		}(),
		expectErr: `^platform\.aws\.subnets\[9\]: Invalid value: \"valid-public-zone-edge-c-2\": edge subnet valid-public-subnet-edge-c is also in zone edge-c$`,
	}, {
		name: "invalid edge pool missing subnets",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfigEdge()
			c.Platform.AWS.Subnets = []string{}
			return c
		}(),
		availZones:     validAvailZonesWithEdge(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		edgeSubnets:    validEdgeSubnets(),
		expectErr:      `^compute\[1\]: Required value: invalid install config\. edge machine pool is valid when installing in existing VPC$`,
	}, {
		name: "invalid edge pool missing edge subnets",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfigEdge()
			return c
		}(),
		availZones:     validAvailZonesWithEdge(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		edgeSubnets:    map[string]Subnet{},
		expectErr:      `^compute\[1\]: Required value: invalid install config\. There is no valid subnets for edge machine pool$`,
	}, {
		name: "invalid edge pool missing subnets on regular zones",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfigEdge()
			c.Platform.AWS.Subnets = []string{}
			edgeSubnets := validEdgeSubnets()
			for subnet := range edgeSubnets {
				c.Platform.AWS.Subnets = append(c.Platform.AWS.Subnets, subnet)
			}
			sort.Strings(c.Platform.AWS.Subnets)
			return c
		}(),
		availZones:     validAvailZonesOnlyEdge(),
		privateSubnets: map[string]Subnet{},
		publicSubnets:  map[string]Subnet{},
		edgeSubnets:    validEdgeSubnets(),
		expectErr:      `^platform\.aws\.subnets: Invalid value: \[\]string{\"valid-public-subnet-edge-a\", \"valid-public-subnet-edge-b\", \"valid-public-subnet-edge-c\"}: edge pool. no subnets configured$`,
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
		expectErr:      `^controlPlane\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\", \"d\"}: No subnets provided for zones \[d\]$`,
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
		expectErr:      `^controlPlane\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\", \"d\", \"e\"}: No subnets provided for zones \[d e\]$`,
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
		expectErr:      `^compute\[0\]\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\", \"d\"}: No subnets provided for zones \[d\]$`,
	}, {
		name: "invalid no subnet for compute zone",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Compute[0].Platform.AWS.Zones = append(c.ControlPlane.Platform.AWS.Zones, "d")
			c.Compute = append(c.Compute, types.MachinePool{
				Architecture: types.ArchitectureAMD64,
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
		expectErr:      `^\[compute\[0\]\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"c\", \"d\"}: No subnets provided for zones \[d\], compute\[1\]\.platform\.aws\.zones: Invalid value: \[\]string{\"a\", \"b\", \"e\"}: No subnets provided for zones \[e\]\]$`,
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
		expectErr:      `^platform\.aws\.serviceEndpoints: Invalid value: (.|\n)*: \[failed to find endpoint for service "ec2": (.|\n)*, failed to find endpoint for service "elasticloadbalancing": (.|\n)*, failed to find endpoint for service "iam": (.|\n)*, failed to find endpoint for service "route53": (.|\n)*, failed to find endpoint for service "s3": (.|\n)*, failed to find endpoint for service "sts": (.|\n)*, failed to find endpoint for service "tagging": (.|\n)*\]$`,
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
		expectErr:      `^platform\.aws\.serviceEndpoints: Invalid value: (.|\n)*: \[failed to find endpoint for service "elasticloadbalancing": (.|\n)*, failed to find endpoint for service "route53": (.|\n)*, failed to find endpoint for service "sts": (.|\n)*, failed to find endpoint for service "tagging": (.|\n)*$`,
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
		name: "AMI omitted for new region in standard partition",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "us-newregion-1"
			c.Platform.AWS.ServiceEndpoints = validServiceEndpoints()
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
	}, {
		name: "accept platform-level AMI",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "us-gov-east-1"
			c.Platform.AWS.AMIID = "custom-ami"
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
	}, {
		name: "accept AMI from default machine platform",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "us-gov-east-1"
			c.Platform.AWS.DefaultMachinePlatform = &aws.MachinePool{AMIID: "custom-ami"}
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
	}, {
		name: "accept AMIs specified for each machine pool",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "us-gov-east-1"
			c.ControlPlane.Platform.AWS.AMIID = "custom-ami"
			c.Compute[0].Platform.AWS.AMIID = "custom-ami"
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
	}, {
		name: "AMI omitted for compute with no replicas",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "us-gov-east-1"
			c.ControlPlane.Platform.AWS.AMIID = "custom-ami"
			c.Compute[0].Replicas = pointer.Int64Ptr(0)
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
		expectErr:      `^platform\.aws\.amiID: Required value: AMI must be provided$`,
	}, {
		name: "invalid endpoint URL",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "us-east-1"
			c.Platform.AWS.ServiceEndpoints = invalidServiceEndpoint()
			c.Platform.AWS.AMIID = "custom-ami"
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		expectErr:      `^\Q[platform.aws.serviceEndpoints[0].url: Invalid value: "testing": Head "testing": unsupported protocol scheme "", platform.aws.serviceEndpoints[1].url: Invalid value: "http://testing.non": Head "http://testing.non": dial tcp: lookup testing.non\E.*: no such host\]$`,
	}, {
		name: "invalid proxy URL but valid URL",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "us-east-1"
			c.Platform.AWS.AMIID = "custom-ami"
			c.Platform.AWS.ServiceEndpoints = []aws.ServiceEndpoint{{Name: "test", URL: "http://testing.com"}}
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		proxy:          "proxy",
	}, {
		name: "invalid proxy URL and invalid URL",
		installConfig: func() *types.InstallConfig {
			c := validInstallConfig()
			c.Platform.AWS.Region = "us-east-1"
			c.Platform.AWS.AMIID = "custom-ami"
			c.Platform.AWS.ServiceEndpoints = []aws.ServiceEndpoint{{Name: "test", URL: "http://test"}}
			return c
		}(),
		availZones:     validAvailZones(),
		privateSubnets: validPrivateSubnets(),
		publicSubnets:  validPublicSubnets(),
		proxy:          "http://proxy.com",
		expectErr:      `^\Qplatform.aws.serviceEndpoints[0].url: Invalid value: "http://test": Head "http://test": dial tcp: lookup test\E.*: no such host$`,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			meta := &Metadata{
				availabilityZones: test.availZones,
				privateSubnets:    test.privateSubnets,
				publicSubnets:     test.publicSubnets,
				edgeSubnets:       test.edgeSubnets,
				instanceTypes:     test.instanceTypes,
			}
			if test.proxy != "" {
				os.Setenv("HTTP_PROXY", test.proxy)
			} else {
				os.Unsetenv("HTTP_PROXY")
			}
			err := Validate(context.TODO(), meta, test.installConfig)
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Regexp(t, test.expectErr, err.Error())
				}
			}
		})
	}
}

func TestIsHostedZoneDomainParentOfClusterDomain(t *testing.T) {
	cases := []struct {
		name             string
		hostedZoneDomain string
		clusterDomain    string
		expected         bool
	}{{
		name:             "same",
		hostedZoneDomain: "c.b.a.",
		clusterDomain:    "c.b.a.",
		expected:         true,
	}, {
		name:             "strict parent",
		hostedZoneDomain: "b.a.",
		clusterDomain:    "c.b.a.",
		expected:         true,
	}, {
		name:             "grandparent",
		hostedZoneDomain: "a.",
		clusterDomain:    "c.b.a.",
		expected:         true,
	}, {
		name:             "not parent",
		hostedZoneDomain: "f.e.d.",
		clusterDomain:    "c.b.a.",
		expected:         false,
	}, {
		name:             "child",
		hostedZoneDomain: "d.c.b.a.",
		clusterDomain:    "c.b.a.",
		expected:         false,
	}, {
		name:             "suffix but not parent",
		hostedZoneDomain: "b.a.",
		clusterDomain:    "cb.a.",
		expected:         false,
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			zone := &route53.HostedZone{Name: &tc.hostedZoneDomain}
			actual := isHostedZoneDomainParentOfClusterDomain(zone, tc.clusterDomain)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestValidateForProvisioning(t *testing.T) {
	cases := []struct {
		name        string
		edits       editFunctions
		expectedErr string
	}{{
		// This really should test for nil, as nothing happened, but no errors were provided
		name:  "internal publish strategy no hosted zone",
		edits: editFunctions{publishInternal, clearHostedZone},
	}, {
		name:        "external publish strategy no hosted zone invalid (empty) base domain",
		edits:       editFunctions{clearHostedZone, clearBaseDomain},
		expectedErr: "baseDomain: Invalid value: \"\": cannot find base domain",
	}, {
		name:        "external publish strategy no hosted zone invalid base domain",
		edits:       editFunctions{clearHostedZone, invalidateBaseDomain},
		expectedErr: "baseDomain: Invalid value: \"invalid-base-domain\": cannot find base domain",
	}, {
		name:  "external publish strategy no hosted zone valid base domain",
		edits: editFunctions{clearHostedZone},
	}, {
		name:  "internal publish strategy valid hosted zone",
		edits: editFunctions{publishInternal},
	}, {
		name:        "internal publish strategy invalid hosted zone",
		edits:       editFunctions{publishInternal, invalidateHostedZone},
		expectedErr: "aws.hostedZone: Invalid value: \"invalid-hosted-zone\": cannot find hosted zone",
	}, {
		name: "external publish strategy valid hosted zone",
	}, {
		name:        "external publish strategy invalid hosted zone",
		edits:       editFunctions{invalidateHostedZone},
		expectedErr: "aws.hostedZone: Invalid value: \"invalid-hosted-zone\": cannot find hosted zone",
	}}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	route53Client := mock.NewMockAPI(mockCtrl)

	validHostedZoneOutput := createValidHostedZone()
	validDomainOutput := createBaseDomainHostedZone()

	route53Client.EXPECT().GetBaseDomain(validDomainName).Return(&validDomainOutput, nil).AnyTimes()
	route53Client.EXPECT().GetBaseDomain("").Return(nil, fmt.Errorf("invalid value: \"\": cannot find base domain")).AnyTimes()
	route53Client.EXPECT().GetBaseDomain(invalidBaseDomain).Return(nil, fmt.Errorf("invalid value: \"%s\": cannot find base domain", invalidBaseDomain)).AnyTimes()

	route53Client.EXPECT().ValidateZoneRecords(&validDomainOutput, gomock.Any(), gomock.Any(), gomock.Any()).Return(field.ErrorList{}).AnyTimes()
	route53Client.EXPECT().ValidateZoneRecords(gomock.Any(), validHostedZoneName, gomock.Any(), gomock.Any()).Return(field.ErrorList{}).AnyTimes()

	// An invalid hosted zone should provide an error
	route53Client.EXPECT().GetHostedZone(validHostedZoneName).Return(&validHostedZoneOutput, nil).AnyTimes()
	route53Client.EXPECT().GetHostedZone(gomock.Not(validHostedZoneName)).Return(nil, fmt.Errorf("invalid value: \"invalid-hosted-zone\": cannot find hosted zone")).AnyTimes()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range test.edits {
				edit(editedInstallConfig)
			}

			meta := &Metadata{
				availabilityZones: validAvailZones(),
				privateSubnets:    validPrivateSubnets(),
				publicSubnets:     validPublicSubnets(),
				instanceTypes:     validInstanceTypes(),
				Region:            editedInstallConfig.AWS.Region,
				vpc:               "valid-private-subnet-a",
			}

			err := ValidateForProvisioning(route53Client, editedInstallConfig, meta)
			if test.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Regexp(t, test.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestGetSubDomainDNSRecords(t *testing.T) {
	cases := []struct {
		name               string
		baseDomain         string
		problematicRecords []string
		expectedErr        string
	}{{
		name:        "empty cluster domain",
		expectedErr: fmt.Sprintf("hosted zone domain %s is not a parent of the cluster domain %s", validDomainName, ""),
	}, {
		name:        "period cluster domain",
		baseDomain:  ".",
		expectedErr: fmt.Sprintf("hosted zone domain %s is not a parent of the cluster domain %s", validDomainName, "."),
	}, {
		name:       "valid dns record no problems",
		baseDomain: validDomainName + ".",
	}, {
		name:               "valid dns record with problems",
		baseDomain:         validDomainName,
		problematicRecords: []string{"test1.ClusterMetaName.valid-base-domain."},
	}, {
		name:               "valid dns record with skipped problems",
		baseDomain:         validDomainName,
		problematicRecords: []string{"test1.ClusterMetaName.valid-base-domain.", "ClusterMetaName.xxxxx-xxxx-xxxxxx."},
	},
	}

	validDomainOutput := createBaseDomainHostedZone()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	route53Client := mock.NewMockAPI(mockCtrl)

	for _, test := range cases {

		t.Run(test.name, func(t *testing.T) {

			ic := validInstallConfig()
			ic.BaseDomain = test.baseDomain

			if test.expectedErr != "" {
				if test.problematicRecords == nil {
					route53Client.EXPECT().GetSubDomainDNSRecords(&validDomainOutput, ic).Return(nil, errors.Errorf(test.expectedErr)).AnyTimes()
				} else {
					// mimic the results of what should happen in the internal function passed to
					// ListResourceRecordSetsPages by GetSubDomainDNSRecords. Skip certain problematicRecords
					returnedProblems := make([]string, 0, len(test.problematicRecords))
					expectedName := ic.ClusterDomain() + "."
					for _, pr := range test.problematicRecords {
						if len(pr) != len(expectedName) {
							returnedProblems = append(returnedProblems, pr)
						}
					}
					route53Client.EXPECT().GetSubDomainDNSRecords(&validDomainOutput, ic).Return(returnedProblems, errors.Errorf(test.expectedErr)).AnyTimes()
				}
			} else {
				route53Client.EXPECT().GetSubDomainDNSRecords(&validDomainOutput, ic).Return(nil, nil).AnyTimes()
			}

			_, err := route53Client.GetSubDomainDNSRecords(&validDomainOutput, ic)
			if test.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Regexp(t, test.expectedErr, err.Error())
				}
			}
		})
	}
}

func TestSkipRecords(t *testing.T) {
	cases := []struct {
		name           string
		recordName     string
		expectedResult bool
	}{{
		name:           "record not part of cluster",
		recordName:     fmt.Sprintf("%s.test.domain.", metaName),
		expectedResult: true,
	}, {
		name:           "record and cluster domain are same",
		recordName:     fmt.Sprintf("%s.%s.", metaName, validDomainName),
		expectedResult: true,
	}, {
		name: "record not part of cluster bad suffix",
		// The parent below does not have a dot following it on purpose - do not Remove
		recordName:     fmt.Sprintf("parent%s.%s.", metaName, validDomainName),
		expectedResult: true,
	}, {
		name: "record part of cluster bad suffix",
		// The parent below does not have a dot following it on purpose - do not Remove
		recordName:     fmt.Sprintf("parent.%s.%s.", metaName, validDomainName),
		expectedResult: false,
	},
	}

	// create the dottedClusterDomain in the same manner that it will be used in GetSubDomainDNSRecords
	ic := validInstallConfig()
	ic.BaseDomain = validDomainName
	dottedClusterDomain := ic.ClusterDomain() + "."

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedResult, skipRecord(test.recordName, dottedClusterDomain))
		})
	}
}
