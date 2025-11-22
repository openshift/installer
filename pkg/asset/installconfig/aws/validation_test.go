package aws

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset/installconfig/aws/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

var (
	metaName   = "ClusterMetaName"
	validCIDR  = "10.0.0.0/16"
	validVPCID = "vpc-valid-id"

	validCallerRef        = "valid-caller-reference"
	validDSId             = "valid-delegation-set-id"
	validHostedZoneName   = "valid-hosted-zone"
	invalidHostedZoneName = "invalid-hosted-zone"
	validNameServers      = []string{"valid-name-server"}

	validDomainName   = "valid-base-domain"
	invalidBaseDomain = "invalid-base-domain"

	// Convert IDs to Subnet type.
	subnetsFromIDs = func(ids []string) []aws.Subnet {
		subnets := make([]aws.Subnet, len(ids))
		for idx, id := range ids {
			subnets[idx] = aws.Subnet{ID: aws.AWSSubnetID(id)}
		}
		// We sort the by ID to ensure a predictable error output in tests
		sort.Slice(subnets, func(i, j int) bool {
			return subnets[i].ID < subnets[j].ID
		})
		return subnets
	}

	// Remove http or https scheme from a URL if any.
	trimURLScheme = func(url string) string {
		if str, found := strings.CutPrefix(url, "https://"); found {
			return str
		}
		if str, found := strings.CutPrefix(url, "http://"); found {
			return str
		}
		return url
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		installConfig *types.InstallConfig
		availRegions  []string
		availZones    []string
		edgeZones     []string
		subnets       SubnetGroups
		subnetsInVPC  *SubnetGroups
		vpcTags       Tags
		instanceTypes map[string]InstanceType
		hosts         map[string]Host
		proxy         string
		publicOnly    bool
		expectErr     string
	}{
		{
			name:          "valid instance types",
			installConfig: icBuild.build(icBuild.withInstanceType("m5.xlarge", "m5.xlarge", "m5.large")),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			instanceTypes: validInstanceTypes(),
		},
		{
			name:          "invalid control plane instance type",
			installConfig: icBuild.build(icBuild.withInstanceType("m5.xlarge", "t2.small", "m5.large")),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			instanceTypes: validInstanceTypes(),
			expectErr:     `^\Q[controlPlane.platform.aws.type: Invalid value: "t2.small": instance type does not meet minimum resource requirements of 4 vCPUs, controlPlane.platform.aws.type: Invalid value: "t2.small": instance type does not meet minimum resource requirements of 16384 MiB Memory]\E$`,
		},
		{
			name:          "invalid compute instance type",
			installConfig: icBuild.build(icBuild.withInstanceType("m5.xlarge", "m5.xlarge", "t2.small")),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			instanceTypes: validInstanceTypes(),
			expectErr:     `^\Q[compute[0].platform.aws.type: Invalid value: "t2.small": instance type does not meet minimum resource requirements of 2 vCPUs, compute[0].platform.aws.type: Invalid value: "t2.small": instance type does not meet minimum resource requirements of 8192 MiB Memory]\E$`,
		},
		{
			name:          "invalid undefined compute instance type",
			installConfig: icBuild.build(icBuild.withInstanceType("m5.xlarge", "m5.xlarge", "m5.dummy")),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			instanceTypes: validInstanceTypes(),
			expectErr:     `^\Qcompute[0].platform.aws.type: Invalid value: "m5.dummy": instance type m5.dummy not found\E$`,
		},
		{
			name: "valid compute pools architectures",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCEdgeSubnetIDs(validSubnets("edge").IDs(), false),
				icBuild.withInstanceArchitecture(types.ArchitectureARM64, types.ArchitectureAMD64, types.ArchitectureAMD64),
			),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			instanceTypes: validInstanceTypes(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				Edge:    validSubnets("edge"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "invalid mismatched instance architecture",
			installConfig: icBuild.build(
				icBuild.withInstanceType("m5.xlarge", "", "m6g.xlarge"),
				icBuild.withInstanceArchitecture(types.ArchitectureARM64, types.ArchitectureAMD64),
			),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			instanceTypes: validInstanceTypes(),
			expectErr:     `^\[controlPlane.platform.aws.type: Invalid value: "m5.xlarge": instance type supported architectures \[amd64\] do not match specified architecture arm64, compute\[0\].platform.aws.type: Invalid value: "m6g.xlarge": instance type supported architectures \[arm64\] do not match specified architecture amd64\]$`,
		},
		{
			name: "invalid mismatched compute pools architectures",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCEdgeSubnetIDs(validSubnets("edge").IDs(), false),
				icBuild.withInstanceArchitecture(types.ArchitectureARM64, types.ArchitectureAMD64, types.ArchitectureARM64),
			),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			instanceTypes: validInstanceTypes(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				Edge:    validSubnets("edge"),
				VpcID:   validVPCID,
			},
			expectErr: `^compute\[1\].architecture: Invalid value: "arm64": all compute machine pools must be of the same architecture$`,
		},
		{
			name: "invalid edge pool, missing zones",
			installConfig: icBuild.build(
				icBuild.withComputeMachinePool([]types.MachinePool{{
					Name: types.MachinePoolEdgeRoleName,
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{},
					},
				}}, true),
				icBuild.withControlPlaneMachinePool(types.MachinePool{}),
			),
			availRegions: validAvailRegions(),
			expectErr:    `^compute\[0\]\.platform\.aws: Required value: zone is required when using edge machine pools$`,
		},
		{
			name: "invalid edge pool, empty zones",
			installConfig: icBuild.build(
				icBuild.withComputeMachinePool([]types.MachinePool{{
					Name: types.MachinePoolEdgeRoleName,
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							Zones: []string{},
						},
					},
				}}, true),
			),
			availRegions: validAvailRegions(),
			expectErr:    `^compute\[0\]\.platform\.aws: Required value: zone is required when using edge machine pools$`,
		},
		{
			name: "invalid edge pool missing platform definition",
			installConfig: icBuild.build(
				icBuild.withComputeMachinePool([]types.MachinePool{{
					Name:     types.MachinePoolEdgeRoleName,
					Platform: types.MachinePoolPlatform{},
				}}, true),
				icBuild.withControlPlaneMachinePool(types.MachinePool{}),
			),
			availRegions: validAvailRegions(),
			expectErr:    `^\[compute\[0\]\.platform\.aws: Required value: edge compute pools are only supported on the AWS platform, compute\[0\].platform.aws: Required value: zone is required when using edge machine pools\]$`,
		},
		{
			name: "valid service endpoints, custom region and no endpoints provided",
			installConfig: icBuild.build(
				icBuild.withPlatformRegion("test-region"),
				icBuild.withPlatformAMIID("dummy-id"),
			),
		},
		{
			name: "valid service endpoints, custom region and some endpoints provided",
			installConfig: icBuild.build(
				icBuild.withPlatformRegion("test-region"),
				icBuild.withPlatformAMIID("dummy-id"),
				icBuild.withServiceEndpoints(validServiceEndpoints()[:3], true),
			),
		},
		{
			name: "valid service endpoints, custom region and all endpoints provided",
			installConfig: icBuild.build(
				icBuild.withPlatformRegion("test-region"),
				icBuild.withPlatformAMIID("dummy-id"),
				icBuild.withServiceEndpoints(validServiceEndpoints(), true),
			),
		},
		{
			name: "invalid service endpoint URLs",
			installConfig: icBuild.build(
				icBuild.withServiceEndpoints(invalidServiceEndpoint(), true),
			),
			availRegions: validAvailRegions(),
			expectErr:    `^\Q[platform.aws.serviceEndpoints[0].url: Invalid value: "bad-aws-endpoint": failed to connect to service endpoint url: Head "bad-aws-endpoint": dial tcp: lookup bad-aws-endpoint: no such host, platform.aws.serviceEndpoints[1].url: Invalid value: "http://bad-aws-endpoint.non": failed to connect to service endpoint url: Head "http://bad-aws-endpoint.non": dial tcp: lookup bad-aws-endpoint.non: no such host]\E$`,
		},
		{
			name: "valid AMI, from platform level",
			installConfig: icBuild.build(
				icBuild.withPlatformRegion("us-gov-east-1"),
				icBuild.withPlatformAMIID("custom-ami"),
			),
			availRegions: validAvailRegions(),
		},
		{
			name: "valid AMI, from default platform machine",
			installConfig: icBuild.build(
				icBuild.withPlatformRegion("us-gov-east-1"),
				icBuild.withDefaultPlatformMachine(aws.MachinePool{AMIID: "custom-ami"}),
			),
			availRegions: validAvailRegions(),
		},
		{
			name: "valid AMIs, from machine pools",
			installConfig: icBuild.build(
				icBuild.withPlatformRegion("us-gov-east-1"),
				icBuild.withControlPlanePlatformAMI("custom-ami"),
				icBuild.withComputePlatformAMI("custom-ami", 0),
			),
			availRegions: validAvailRegions(),
		},
		{
			name: "valid AMI, omitted for compute with no replicas",
			installConfig: icBuild.build(
				icBuild.withPlatformRegion("us-gov-east-1"),
				icBuild.withControlPlanePlatformAMI("custom-ami"),
				icBuild.withComputeReplicas(0, 0),
			),
			availRegions: validAvailRegions(),
		},
		{
			name: "invalid AMI not provided for unknown region",
			installConfig: icBuild.build(
				icBuild.withPlatformRegion("test-region"),
			),
			availRegions: validAvailRegions(),
			expectErr:    `^platform\.aws\.amiID: Required value: AMI must be provided$`,
		},
		{
			name: "invalid proxy URL but valid service endpoint URL",
			installConfig: icBuild.build(
				icBuild.withPlatformAMIID("custom-ami"),
				icBuild.withServiceEndpoints(validServiceEndpoints(), true),
			),
			availRegions: validAvailRegions(),
			proxy:        "proxy",
		},
		{
			name: "invalid proxy URL and invalid service endpoint URL",
			installConfig: icBuild.build(
				icBuild.withPlatformAMIID("custom-ami"),
				icBuild.withServiceEndpoints(invalidServiceEndpoint(), true),
			),
			availRegions: validAvailRegions(),
			proxy:        "http://proxy.com",
			expectErr:    `^\Q[platform.aws.serviceEndpoints[0].url: Invalid value: "bad-aws-endpoint": failed to connect to service endpoint url: Head "bad-aws-endpoint": dial tcp: lookup bad-aws-endpoint: no such host, platform.aws.serviceEndpoints[1].url: Invalid value: "http://bad-aws-endpoint.non": failed to connect to service endpoint url: Head "http://bad-aws-endpoint.non": dial tcp: lookup bad-aws-endpoint.non: no such host]\E$`,
		},
		{
			name: "invalid public ipv4 pool in private installation",
			installConfig: icBuild.build(
				icBuild.withPublish(types.InternalPublishingStrategy),
				icBuild.withPublicIPv4Pool("ipv4pool-ec2-123"),
				icBuild.withVPCSubnets([]aws.Subnet{}, true),
			),
			availRegions: validAvailRegions(),
			expectErr:    `^platform.aws.publicIpv4PoolId: Invalid value: "ipv4pool-ec2-123": publish strategy Internal can't be used with custom Public IPv4 Pools$`,
		},
		{
			name:          "valid no byo subnets, unspecified subnet list",
			installConfig: icBuild.build(),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
		},
		{
			name:          "valid no byo subnets, empty subnet list",
			installConfig: icBuild.build(icBuild.withVPCSubnets([]aws.Subnet{}, true)),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
		},
		{
			name:          "valid no byo subnets, unspecified subnet list for public-only subnets cluster",
			installConfig: icBuild.build(),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			publicOnly:    true,
		},
		{
			name:          "valid byo subnets",
			installConfig: icBuild.build(icBuild.withBaseBYO()),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "valid byo subnets, include edge subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCEdgeSubnetIDs(validSubnets("edge").IDs(), false),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				Edge:    validSubnets("edge"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "valid byo subnets, private subnets only with publish internal",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withPublish(types.InternalPublishingStrategy),
				icBuild.withVPCSubnetIDs(validSubnets("private").IDs(), true),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "invalid byo subnets, no private subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withPublish(types.InternalPublishingStrategy),
				icBuild.withVPCSubnetIDs(validSubnets("public").IDs(), true),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Public: validSubnets("public"),
				VpcID:  validVPCID,
			},
			expectErr: `^\[platform\.aws\.vpc\.subnets: Invalid value: \[{"id":"subnet-valid-public-a"},{"id":"subnet-valid-public-b"},{"id":"subnet-valid-public-c"}\]: no private subnets found, controlPlane\.platform\.aws\.zones: Invalid value: \[\"a\",\"b\",\"c\"]: No subnets provided for zones \[a b c\], compute\[0\]\.platform\.aws\.zones: Invalid value: \[\"a\",\"b\",\"c\"]: No subnets provided for zones \[a b c\]\]$`,
		},
		{
			name: "invalid byo subnets, no public subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnetIDs(validSubnets("private").IDs(), true),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				VpcID:   validVPCID,
			},
			expectErr: `^platform\.aws\.vpc\.subnets: Invalid value: \[{"id":"subnet-valid-private-a"},{"id":"subnet-valid-private-b"},{"id":"subnet-valid-private-c"}\]: No public subnet provided for zones \[a b c\]$`,
		},
		{
			name: "invalid byo subnets, invalid cidr does not belong to machine CIDR",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnetIDs([]string{"invalid-private-cidr-subnet", "invalid-public-cidr-subnet"}, false),
			),
			availRegions: validAvailRegions(),
			availZones:   append(validAvailZones(), "zone-for-invalid-cidr-subnet"),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{"invalid-private-cidr-subnet": Subnet{
					ID:   "invalid-private-cidr-subnet",
					Zone: &Zone{Name: "zone-for-invalid-cidr-subnet"},
					CIDR: "192.168.126.0/24",
				}}),
				Public: mergeSubnets(validSubnets("public"), Subnets{"invalid-public-cidr-subnet": Subnet{
					ID:   "invalid-public-cidr-subnet",
					Zone: &Zone{Name: "zone-for-invalid-cidr-subnet"},
					CIDR: "192.168.127.0/24",
				}}),
				VpcID: validVPCID,
			},
			expectErr: `^\[platform\.aws\.vpc\.subnets\[6\]: Invalid value: \"invalid-private-cidr-subnet\": subnet's CIDR range start 192\.168\.126\.0 is outside of the specified machine networks, platform\.aws\.vpc\.subnets\[7\]: Invalid value: \"invalid-public-cidr-subnet\": subnet's CIDR range start 192\.168\.127\.0 is outside of the specified machine networks\]$`,
		},
		{
			name: "invalid byo subnets, missing public subnet in a zone",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnetIDs([]string{"no-matching-public-private-zone"}, false),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{"no-matching-public-private-zone": Subnet{
					ID:   "no-matching-public-private-zone",
					Zone: &Zone{Name: "f"},
					CIDR: "10.0.7.0/24",
				}}),
				Public: validSubnets("public"),
				VpcID:  validVPCID,
			},
			expectErr: `^platform\.aws\.vpc\.subnets: Invalid value: \[{"id":"subnet-valid-private-a"},{"id":"subnet-valid-private-b"},{"id":"subnet-valid-private-c"},{"id":"subnet-valid-public-a"},{"id":"subnet-valid-public-b"},{"id":"subnet-valid-public-c"},{"id":"no-matching-public-private-zone"}\]: No public subnet provided for zones \[f\]$`,
		},
		{
			name: "invalid byo subnets with no roles, multiple private in same zone",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnetIDs([]string{"valid-private-zone-c-2"}, false),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{"valid-private-zone-c-2": Subnet{
					ID:   "valid-private-zone-c-2",
					Zone: &Zone{Name: "c"},
					CIDR: "10.0.7.0/24",
				}}),
				Public: validSubnets("public"),
				VpcID:  validVPCID,
			},
			expectErr: `^platform\.aws\.vpc\.subnets\[6\]: Invalid value: \"valid-private-zone-c-2\": private subnet subnet-valid-private-c is also in zone c$`,
		},
		{
			name: "invalid byo subnets with no roles, multiple public in same zone",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnetIDs([]string{"valid-public-zone-c-2"}, false),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public: mergeSubnets(validSubnets("public"), Subnets{"valid-public-zone-c-2": Subnet{
					ID:   "valid-public-zone-c-2",
					Zone: &Zone{Name: "c"},
					CIDR: "10.0.7.0/24",
				}}),
				VpcID: validVPCID,
			},
			expectErr: `^platform\.aws\.vpc\.subnets\[6\]: Invalid value: \"valid-public-zone-c-2\": public subnet subnet-valid-public-c is also in zone c$`,
		},
		{
			name: "invalid byo subnets with no roles, multiple public edge in same zone",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCEdgeSubnetIDs(validSubnets("edge").IDs(), false),
				icBuild.withVPCEdgeSubnetIDs([]string{"valid-public-zone-edge-c-2"}, false),
			),
			availRegions: validAvailRegions(),
			availZones:   append(validAvailZones(), validEdgeAvailZones()...),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				Edge: mergeSubnets(validSubnets("edge"), Subnets{
					"valid-public-zone-edge-c-2": Subnet{
						ID:   "valid-public-zone-edge-c-2",
						Zone: &Zone{Name: "edge-c", Type: aws.LocalZoneType},
						CIDR: "10.0.9.0/24",
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `^platform\.aws\.vpc\.subnets\[9\]: Invalid value: \"valid-public-zone-edge-c-2\": edge subnet subnet-valid-public-edge-c is also in zone edge-c$`,
		},
		{
			name: "invalid byo subnets, edge pool missing valid subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCEdgeSubnetIDs(validSubnets("edge").IDs(), false),
			),
			availRegions: validAvailRegions(),
			availZones:   append(validAvailZones(), validEdgeAvailZones()...),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			expectErr: `^compute\[1\]\.platform\.aws: Required value: the provided subnets must include valid subnets for the specified edge zones$`,
		},
		{
			name: "invalid byo subnets, edge pool missing subnets on availability zones",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCEdgeSubnetIDs(validSubnets("edge").IDs(), true),
			),
			availRegions: validAvailRegions(),
			availZones:   validEdgeAvailZones(),
			subnets: SubnetGroups{
				Edge:  validSubnets("edge"),
				VpcID: validVPCID,
			},
			expectErr: `^\[platform\.aws\.vpc\.subnets: Invalid value: \[{"id":"subnet-valid-public-edge-a"},{"id":"subnet-valid-public-edge-b"},{"id":"subnet-valid-public-edge-c"}\]: no private subnets found, controlPlane\.platform\.aws\.zones: Invalid value: \[\"a\",\"b\",\"c\"]: No subnets provided for zones \[a b c\], compute\[0\]\.platform\.aws\.zones: Invalid value: \[\"a\",\"b\",\"c\"]: No subnets provided for zones \[a b c\]\]$`,
		},
		{
			name: "invalid byo subnets, no subnet for control plane zones",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withControlPlanePlatformZones([]string{"d", "e"}, false),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			expectErr: `^controlPlane\.platform\.aws\.zones: Invalid value: \[\"a\",\"b\",\"c\",\"d\",\"e\"\]: No subnets provided for zones \[d e\]$`,
		},
		{
			name: "invalid byo subnets, no subnet for compute[0] zones",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withComputePlatformZones([]string{"d"}, false, 0),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			expectErr: `^compute\[0\]\.platform\.aws\.zones: Invalid value: \[\"a\",\"b\",\"c\",\"d\"\]: No subnets provided for zones \[d\]$`,
		},
		{
			name: "invalid byo subnets, no subnet for compute zone",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withComputePlatformZones([]string{"d"}, false, 0),
				icBuild.withComputeMachinePool([]types.MachinePool{{
					Architecture: types.ArchitectureAMD64,
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							Zones: []string{"a", "b", "e"},
						},
					},
				}}, false),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			expectErr: `^\[compute\[0\]\.platform\.aws\.zones: Invalid value: \[\"a\",\"b\",\"c\",\"d\"\]: No subnets provided for zones \[d\], compute\[1\]\.platform\.aws\.zones: Invalid value: \[\"a\",\"b\",\"e\"\]: No subnets provided for zones \[e\]\]$`,
		},
		{
			name:          "valid byo subnets, private and public subnets provided for public-only subnets cluster",
			installConfig: icBuild.build(icBuild.withBaseBYO()),
			availRegions:  validAvailRegions(),
			availZones:    validAvailZones(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("public"), validSubnets("private")),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			publicOnly: true,
		},
		{
			name: "valid byo subnets, public subnets provided for public-only subnets cluster",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnetIDs(validSubnets("public").IDs(), true),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("public"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			publicOnly: true,
		},
		{
			name: "invalid byo subnets, no public subnets specified for public-only subnets cluster",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnetIDs(validSubnets("private").IDs(), true),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				VpcID:   validVPCID,
			},
			publicOnly: true,
			expectErr:  `^\Q[platform.aws.vpc.subnets: Required value: public subnets are required for a public-only subnets cluster, platform.aws.vpc.subnets: Invalid value: [{"id":"subnet-valid-private-a"},{"id":"subnet-valid-private-b"},{"id":"subnet-valid-private-c"}]: No public subnet provided for zones [a b c]]\E$`,
		},
		{
			name: "invalid byo subnets, internal publish method for public-only subnets install",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withPublish(types.InternalPublishingStrategy),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			publicOnly: true,
			expectErr:  `^publish: Invalid value: \"Internal\": cluster cannot be private with public subnets$`,
		},
		{
			name: "valid byo subnets, no roles and vpc has no untagged subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			subnetsInVPC: &SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), otherTaggedPrivateSubnets()),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "invalid byo subnets, no roles but vpc has untagged subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			subnetsInVPC: &SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), otherUntaggedPrivateSubnets()),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			expectErr: `^platform\.aws\.vpc\.subnets: Forbidden: additional subnets \[subnet-valid-private-a1 subnet-valid-private-b1\] without tag prefix kubernetes\.io/cluster/ are found in vpc vpc-valid-id of provided subnets\. Please add a tag kubernetes\.io/cluster/unmanaged to those subnets to exclude them from cluster installation or explicitly assign roles in the install-config to provided subnets$`,
		},
		{
			name: "valid byo subnets with roles",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "valid byo subnets with roles, include edge subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCEdgeSubnets(byoEdgeSubnetsWithRoles(), false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				Edge:    validSubnets("edge"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "valid byo subnets with roles, ignore other untagged subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			subnetsInVPC: &SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), otherUntaggedPrivateSubnets()),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "valid byo subnets with roles, vpc has other tagged subnets",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			subnetsInVPC: &SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), otherTaggedPrivateSubnets()),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
		},
		{
			name: "valid byo subnets with roles, dedicated subnets for LBs",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withControlPlanePlatformZones([]string{"a"}, true),
				icBuild.withComputePlatformZones([]string{"a"}, true, 0),
				icBuild.withVPCSubnets(
					[]aws.Subnet{
						{
							ID: "subnet-valid-private-a",
							Roles: []aws.SubnetRole{
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-valid-private-a1",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneInternalLBSubnetRole},
							},
						},
						{
							ID: "subnet-valid-public-a",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneExternalLBSubnetRole},
								{Type: aws.BootstrapNodeSubnetRole},
							},
						},
						{
							ID: "subnet-valid-public-a1",
							Roles: []aws.SubnetRole{
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					}, true),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: Subnets{
					"subnet-valid-private-a": {
						ID:     "subnet-valid-private-a",
						Zone:   &Zone{Name: "a"},
						CIDR:   "10.0.1.0/24",
						Public: false,
					},
					"subnet-valid-private-a1": {
						ID:   "subnet-valid-private-a1",
						Zone: &Zone{Name: "a"},
						CIDR: "10.0.2.0/24",
					},
				},
				Public: Subnets{
					"subnet-valid-public-a": {
						ID:     "subnet-valid-public-a",
						Zone:   &Zone{Name: "a"},
						CIDR:   "10.0.3.0/24",
						Public: true,
					},
					"subnet-valid-public-a1": {
						ID:     "subnet-valid-public-a1",
						Zone:   &Zone{Name: "a"},
						CIDR:   "10.0.4.0/24",
						Public: true,
					},
				},
				VpcID: validVPCID,
			},
		},
		{
			name: "invalid byo subnets with roles, subnets of the same role in the same AZs",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID: "subnet-valid-private-a1",
						Roles: []aws.SubnetRole{
							{Type: aws.ClusterNodeSubnetRole},
							{Type: aws.ControlPlaneInternalLBSubnetRole},
						},
					},
					{
						ID: "subnet-valid-public-a1",
						Roles: []aws.SubnetRole{
							{Type: aws.BootstrapNodeSubnetRole},
							{Type: aws.ControlPlaneExternalLBSubnetRole},
							{Type: aws.IngressControllerLBSubnetRole},
						},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{
					"subnet-valid-private-a1": {
						ID:   "subnet-valid-private-a1",
						Zone: &Zone{Name: "a"},
						CIDR: "10.0.6.0/24",
					},
				}),
				Public: mergeSubnets(validSubnets("public"), Subnets{
					"subnet-valid-public-a1": {
						ID:     "subnet-valid-public-a1",
						Zone:   &Zone{Name: "a"},
						CIDR:   "10.0.6.0/24",
						Public: true,
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `^\Q[platform.aws.vpc.subnets[6]: Invalid value: "subnet-valid-private-a1": subnets subnet-valid-private-a and subnet-valid-private-a1 have role ClusterNode and are both in zone a, platform.aws.vpc.subnets[7]: Invalid value: "subnet-valid-public-a1": subnets subnet-valid-public-a and subnet-valid-public-a1 have role BootstrapNode and are both in zone a, platform.aws.vpc.subnets[7]: Invalid value: "subnet-valid-public-a1": subnets subnet-valid-public-a and subnet-valid-public-a1 have role IngressControllerLB and are both in zone a, platform.aws.vpc.subnets[7]: Invalid value: "subnet-valid-public-a1": subnets subnet-valid-public-a and subnet-valid-public-a1 have role ControlPlaneExternalLB and are both in zone a, platform.aws.vpc.subnets[6]: Invalid value: "subnet-valid-private-a1": subnets subnet-valid-private-a and subnet-valid-private-a1 have role ControlPlaneInternalLB and are both in zone a]\E$`,
		},
		{
			name: "invalid byo subnets with roles, public subnet assigned ClusterNode",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-public-a1",
						Roles: []aws.SubnetRole{{Type: aws.ClusterNodeSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public: mergeSubnets(validSubnets("public"), Subnets{
					"subnet-valid-public-a1": {
						ID:     "subnet-valid-public-a1",
						Zone:   &Zone{Name: "a"},
						CIDR:   "10.0.6.0/24",
						Public: true,
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `\Qplatform.aws.vpc.subnets[6]: Invalid value: "subnet-valid-public-a1": subnet subnet-valid-public-a1 has role ClusterNode, but is public, expected to be private\E`,
		},
		{
			name: "invalid byo subnets with roles, private subnet assigned Bootstrap in external cluster",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-private-a1",
						Roles: []aws.SubnetRole{{Type: aws.BootstrapNodeSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{
					"subnet-valid-private-a1": {
						ID:   "subnet-valid-private-a1",
						Zone: &Zone{Name: "a"},
						CIDR: "10.0.6.0/24",
					},
				}),
				Public: validSubnets("public"),
				VpcID:  validVPCID,
			},
			expectErr: `\Qplatform.aws.vpc.subnets[6]: Invalid value: "subnet-valid-private-a1": subnet subnet-valid-private-a1 has role BootstrapNode, but is private, expected to be public\E`,
		},
		{
			name: "invalid byo subnets with roles, private subnet assigned ControlPlaneExternalLB",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-private-a1",
						Roles: []aws.SubnetRole{{Type: aws.ControlPlaneExternalLBSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{
					"subnet-valid-private-a1": {
						ID:   "subnet-valid-private-a1",
						Zone: &Zone{Name: "a"},
						CIDR: "10.0.6.0/24",
					},
				}),
				Public: validSubnets("public"),
				VpcID:  validVPCID,
			},
			expectErr: `\Qplatform.aws.vpc.subnets[6]: Invalid value: "subnet-valid-private-a1": subnet subnet-valid-private-a1 has role ControlPlaneExternalLB, but is private, expected to be public\E`,
		},
		{
			name: "invalid byo subnets with roles, public subnet assigned ControlPlaneInternalLB",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-public-a1",
						Roles: []aws.SubnetRole{{Type: aws.ControlPlaneInternalLBSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public: mergeSubnets(validSubnets("public"), Subnets{
					"subnet-valid-public-a1": {
						ID:     "subnet-valid-public-a1",
						Zone:   &Zone{Name: "a"},
						CIDR:   "10.0.6.0/24",
						Public: true,
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `\Qplatform.aws.vpc.subnets[6]: Invalid value: "subnet-valid-public-a1": subnet subnet-valid-public-a1 has role ControlPlaneInternalLB, but is public, expected to be private\E`,
		},
		{
			name: "valid byo subnets with roles, public subnet assigned ControlPlaneInternalLB, ClusterNode and public-only cluster",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoPublicOnlySubnetsWithRoles(), true),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Public:  validSubnets("public"),
				Private: validSubnets("public"),
				VpcID:   validVPCID,
			},
			publicOnly: true,
		},
		{
			name: "invalid byo subnets with roles, public subnet assigned IngressControllerLB when publish is internal",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withPublish(types.InternalPublishingStrategy),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			expectErr: `platform\.aws\.vpc\.subnets\[3\]: Invalid value: \"subnet-valid-public-a\": subnet subnet-valid-public-a has role IngressControllerLB and is public, which is not allowed when publish is set to Internal`,
		},
		{
			name: "invalid byo subnets with roles, private subnet assigned IngressControllerLB when publish is external",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-private-a1",
						Roles: []aws.SubnetRole{{Type: aws.IngressControllerLBSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{
					"subnet-valid-private-a1": {
						ID:   "subnet-valid-private-a1",
						Zone: &Zone{Name: "a"},
						CIDR: "10.0.6.0/24",
					},
				}),
				Public: validSubnets("public"),
				VpcID:  validVPCID,
			},
			expectErr: `platform\.aws\.vpc\.subnets\[6\]: Invalid value: \"subnet-valid-private-a1\": subnet subnet-valid-private-a1 has role IngressControllerLB and is private, which is not allowed when publish is set to External`,
		},
		{
			name: "invalid byo subnets with roles, subnets assigned IngressControllerLB in the same zone",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-public-a1",
						Roles: []aws.SubnetRole{{Type: aws.IngressControllerLBSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public: mergeSubnets(validSubnets("public"), Subnets{
					"subnet-valid-public-a1": {
						ID:     "subnet-valid-public-a1",
						Zone:   &Zone{Name: "a"},
						CIDR:   "10.0.6.0/24",
						Public: true,
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `^\Qplatform.aws.vpc.subnets[6]: Invalid value: "subnet-valid-public-a1": subnets subnet-valid-public-a and subnet-valid-public-a1 have role IngressControllerLB and are both in zone a\E$`,
		},
		{
			name: "invalid byo subnets with roles, AZs of IngressControllerLB and ControlPlaneLB not match AZs of ClusterNode",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-public-f",
						Roles: []aws.SubnetRole{{Type: aws.ControlPlaneExternalLBSubnetRole}, {Type: aws.IngressControllerLBSubnetRole}},
					},
					{
						ID:    "subnet-valid-private-f",
						Roles: []aws.SubnetRole{{Type: aws.ControlPlaneInternalLBSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{
					"subnet-valid-private-f": {
						ID:   "subnet-valid-private-f",
						Zone: &Zone{Name: "f"},
						CIDR: "10.0.6.0/24",
					},
				}),
				Public: mergeSubnets(validSubnets("public"), Subnets{
					"subnet-valid-public-f": {
						ID:     "subnet-valid-public-f",
						Zone:   &Zone{Name: "f"},
						CIDR:   "10.0.6.0/24",
						Public: true,
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `^\Q[platform.aws.vpc.subnets: Forbidden: zones [f] are enabled for ControlPlaneInternalLB load balancers, but are not used by any nodes, platform.aws.vpc.subnets: Forbidden: zones [f] are enabled for IngressControllerLB load balancers, but are not used by any nodes, platform.aws.vpc.subnets: Forbidden: zones [f] are enabled for ControlPlaneExternalLB load balancers, but are not used by any nodes]\E$`,
		},
		{
			name: "invalid byo subnets with roles, AZs of ClusterNode not match AZs of IngressControllerLB and ControlPlaneLB",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-private-f",
						Roles: []aws.SubnetRole{{Type: aws.ClusterNodeSubnetRole}},
					},
					{
						ID:    "subnet-valid-public-f",
						Roles: []aws.SubnetRole{{Type: aws.BootstrapNodeSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: mergeSubnets(validSubnets("private"), Subnets{
					"subnet-valid-private-f": {
						ID:   "subnet-valid-private-f",
						Zone: &Zone{Name: "f"},
						CIDR: "10.0.6.0/24",
					},
				}),
				Public: mergeSubnets(validSubnets("public"), Subnets{
					"subnet-valid-public-f": {
						ID:     "subnet-valid-public-f",
						Zone:   &Zone{Name: "f"},
						CIDR:   "10.0.6.0/24",
						Public: true,
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `^\Q[platform.aws.vpc.subnets: Forbidden: zones [f] are not enabled for ControlPlaneInternalLB load balancers, nodes in those zones are unreachable, platform.aws.vpc.subnets: Forbidden: zones [f] are not enabled for IngressControllerLB load balancers, nodes in those zones are unreachable, platform.aws.vpc.subnets: Forbidden: zones [f] are not enabled for ControlPlaneExternalLB load balancers, nodes in those zones are unreachable]\E$`,
		},
		{
			name: "invalid byo subnets with roles, subnet assigned EdgeNode but not edge subnet",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCEdgeSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-public-a1",
						Roles: []aws.SubnetRole{{Type: aws.EdgeNodeSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public: mergeSubnets(validSubnets("public"), Subnets{
					"subnet-valid-public-a1": {
						ID:     "subnet-valid-public-a1",
						Zone:   &Zone{Name: "a"},
						CIDR:   "10.0.6.0/24",
						Public: true,
					},
				}),
				Edge:  validSubnets("edge"),
				VpcID: validVPCID,
			},
			expectErr: `platform\.aws\.vpc\.subnets\[6\]: Invalid value: \"subnet-valid-public-a1\": subnet subnet-valid-public-a1 has role EdgeNode, but is not in a Local or WaveLength Zone`,
		},
		{
			name: "invalid byo subnets with roles, edge subnet assigned with other roles",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets(byoSubnetsWithRoles(), true),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID:    "subnet-valid-public-edge-a1",
						Roles: []aws.SubnetRole{{Type: aws.BootstrapNodeSubnetRole}},
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				Edge: mergeSubnets(validSubnets("edge"), Subnets{
					"subnet-valid-public-edge-a1": {
						ID:     "subnet-valid-public-edge-a1",
						Zone:   &Zone{Name: "edge-a"},
						CIDR:   "10.0.6.0/24",
						Public: true,
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `platform\.aws\.vpc\.subnets\[6\]: Invalid value: \"subnet-valid-public-edge-a1\": subnet subnet-valid-public-edge-a1 must only be assigned role EdgeNode since it is in a Local or WaveLength Zone`,
		},
		{
			name: "invalid byo subnets, vpc has cluster-owned tags",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public:  validSubnets("public"),
				VpcID:   validVPCID,
			},
			vpcTags: Tags{
				"kubernetes.io/cluster/another-cluster": "owned",
			},
			expectErr: `^\Qplatform.aws.vpc.subnets: Forbidden: VPC of subnets is owned by other clusters [another-cluster] and cannot be used for new installations, another VPC must be created separately\E$`,
		},
		{
			name: "invalid byo subnets, subnets have cluster-owned tags",
			installConfig: icBuild.build(
				icBuild.withBaseBYO(),
				icBuild.withVPCSubnets([]aws.Subnet{
					{
						ID: "subnet-valid-public-d",
					},
				}, false),
			),
			availRegions: validAvailRegions(),
			subnets: SubnetGroups{
				Private: validSubnets("private"),
				Public: mergeSubnets(validSubnets("public"), Subnets{
					"subnet-valid-public-a1": {
						ID:     "subnet-valid-public-d",
						Zone:   &Zone{Name: "d"},
						CIDR:   "10.0.6.0/24",
						Public: true,
						Tags: Tags{
							"kubernetes.io/cluster/another-cluster": "owned",
						},
					},
				}),
				VpcID: validVPCID,
			},
			expectErr: `^\Qplatform.aws.vpc.subnets: Forbidden: subnet subnet-valid-public-a1 is owned by other clusters [another-cluster] and cannot be used for new installations, another subnet must be created separately\E$`,
		},
		{
			name: "valid dedicated host placement on compute",
			installConfig: icBuild.build(
				icBuild.withComputePlatformZones([]string{"a"}, true, 0),
				icBuild.withComputeHostPlacement([]string{"h-1234567890abcdef0"}, 0),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			hosts: map[string]Host{
				"h-1234567890abcdef0": {ID: "h-1234567890abcdef0", Zone: "a"},
			},
		},
		{
			name: "invalid dedicated host not found",
			installConfig: icBuild.build(
				icBuild.withComputePlatformZones([]string{"a"}, true, 0),
				icBuild.withComputeHostPlacement([]string{"h-aaaaaaaaaaaaaaaaa"}, 0),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			hosts: map[string]Host{
				"h-1234567890abcdef0": {ID: "h-1234567890abcdef0", Zone: "a"},
			},
			expectErr: "dedicated host h-aaaaaaaaaaaaaaaaa not found",
		},
		{
			name: "invalid dedicated host zone not in pool zones",
			installConfig: icBuild.build(
				icBuild.withComputePlatformZones([]string{"a"}, true, 0),
				icBuild.withComputeHostPlacement([]string{"h-bbbbbbbbbbbbbbbbb"}, 0),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			hosts: map[string]Host{
				"h-bbbbbbbbbbbbbbbbb": {ID: "h-bbbbbbbbbbbbbbbbb", Zone: "b"},
			},
			expectErr: "is not available in pool's zone list",
		},
		{
			name: "dedicated host placement on compute but for a zone that pool is not using",
			installConfig: icBuild.build(
				icBuild.withComputePlatformZones([]string{"b"}, true, 0),
				icBuild.withComputeHostPlacementAndZone([]string{"h-1234567890abcdef0"}, "b", 0),
			),
			availRegions: validAvailRegions(),
			availZones:   validAvailZones(),
			hosts: map[string]Host{
				"h-1234567890abcdef0": {ID: "h-1234567890abcdef0", Zone: "a"},
			},
			expectErr: "dedicated host was configured with zone b but expected zone a",
		},
	}

	// Register mock http(s) responses for tests.
	httpmock.Activate()
	t.Cleanup(httpmock.Deactivate)

	for _, endpoint := range validServiceEndpoints() {
		httpmock.RegisterResponder(http.MethodHead, endpoint.URL, func(r *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})
	}

	for _, endpoint := range invalidServiceEndpoint() {
		httpmock.RegisterResponder(http.MethodHead, endpoint.URL, func(r *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("dial tcp: lookup %s: no such host", trimURLScheme(endpoint.URL))
		})
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			meta := &Metadata{
				availabilityZones: test.availZones,
				availableRegions:  test.availRegions,
				edgeZones:         test.edgeZones,
				subnets:           test.subnets,
				vpcSubnets:        test.subnets,
				vpc: VPC{
					ID:   validVPCID,
					CIDR: validCIDR,
					Tags: test.vpcTags,
				},
				instanceTypes:   test.instanceTypes,
				Hosts:           test.hosts,
				ProvidedSubnets: test.installConfig.Platform.AWS.VPC.Subnets,
			}

			if test.subnetsInVPC != nil {
				meta.vpcSubnets = *test.subnetsInVPC
			}
			if test.proxy != "" {
				os.Setenv("HTTP_PROXY", test.proxy)
			} else {
				os.Unsetenv("HTTP_PROXY")
			}
			if test.publicOnly {
				os.Setenv("OPENSHIFT_INSTALL_AWS_PUBLIC_ONLY", "true")
			} else {
				os.Unsetenv("OPENSHIFT_INSTALL_AWS_PUBLIC_ONLY")
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
		icOptions   []icOption
		expectedErr string
	}{{
		// This really should test for nil, as nothing happened, but no errors were provided
		name:      "internal publish strategy no hosted zone",
		icOptions: []icOption{icBuild.withPublish(types.InternalPublishingStrategy), icBuild.withHostedZone("")},
	}, {
		name:        "external publish strategy no hosted zone invalid (empty) base domain",
		icOptions:   []icOption{icBuild.withHostedZone(""), icBuild.withBaseDomain("")},
		expectedErr: "baseDomain: Invalid value: \"\": cannot find base domain",
	}, {
		name:        "external publish strategy no hosted zone invalid base domain",
		icOptions:   []icOption{icBuild.withHostedZone(""), icBuild.withBaseDomain(invalidBaseDomain)},
		expectedErr: "baseDomain: Invalid value: \"invalid-base-domain\": cannot find base domain",
	}, {
		name:      "external publish strategy no hosted zone valid base domain",
		icOptions: []icOption{icBuild.withHostedZone("")},
	}, {
		name:      "internal publish strategy valid hosted zone",
		icOptions: []icOption{icBuild.withPublish(types.InternalPublishingStrategy)},
	}, {
		name:        "internal publish strategy invalid hosted zone",
		icOptions:   []icOption{icBuild.withPublish(types.InternalPublishingStrategy), icBuild.withHostedZone(invalidHostedZoneName)},
		expectedErr: "aws.hostedZone: Invalid value: \"invalid-hosted-zone\": unable to retrieve hosted zone",
	}, {
		name: "external publish strategy valid hosted zone",
	}, {
		name:        "external publish strategy invalid hosted zone",
		icOptions:   []icOption{icBuild.withHostedZone(invalidHostedZoneName)},
		expectedErr: "aws.hostedZone: Invalid value: \"invalid-hosted-zone\": unable to retrieve hosted zone",
	}}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	route53Client := mock.NewMockAPI(mockCtrl)

	validHostedZoneOutput := createValidHostedZoneOutput()
	validDomainOutput := createBaseDomainHostedZoneOutput()

	route53Client.EXPECT().GetBaseDomain(validDomainName).Return(&validDomainOutput, nil).AnyTimes()
	route53Client.EXPECT().GetBaseDomain("").Return(nil, fmt.Errorf("invalid value: \"\": cannot find base domain")).AnyTimes()
	route53Client.EXPECT().GetBaseDomain(invalidBaseDomain).Return(nil, fmt.Errorf("invalid value: \"%s\": cannot find base domain", invalidBaseDomain)).AnyTimes()

	route53Client.EXPECT().ValidateZoneRecords(&validDomainOutput, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(field.ErrorList{}).AnyTimes()
	route53Client.EXPECT().ValidateZoneRecords(gomock.Any(), validHostedZoneName, gomock.Any(), gomock.Any(), gomock.Any()).Return(field.ErrorList{}).AnyTimes()

	// An invalid hosted zone should provide an error
	route53Client.EXPECT().GetHostedZone(validHostedZoneName, gomock.Any()).Return(&validHostedZoneOutput, nil).AnyTimes()
	route53Client.EXPECT().GetHostedZone(gomock.Not(validHostedZoneName), gomock.Any()).Return(nil, fmt.Errorf("invalid value: \"invalid-hosted-zone\": cannot find hosted zone")).AnyTimes()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			ic := icBuild.build(test.icOptions...)
			meta := &Metadata{
				availabilityZones: validAvailZones(),
				subnets: SubnetGroups{
					Private: validSubnets("private"),
					Public:  validSubnets("public"),
					VpcID:   validVPCID,
				},
				instanceTypes: validInstanceTypes(),
				Region:        ic.AWS.Region,
				vpc: VPC{
					ID:   validVPCID,
					CIDR: validCIDR,
				},
				ProvidedSubnets: ic.Platform.AWS.VPC.Subnets,
			}

			err := ValidateForProvisioning(route53Client, ic, meta)
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

	validDomainOutput := createBaseDomainHostedZoneOutput()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	route53Client := mock.NewMockAPI(mockCtrl)

	for _, test := range cases {

		t.Run(test.name, func(t *testing.T) {
			ic := icBuild.build(icBuild.withBaseDomain(test.baseDomain))
			if test.expectedErr != "" {
				if test.problematicRecords == nil {
					route53Client.EXPECT().GetSubDomainDNSRecords(&validDomainOutput, ic, gomock.Any()).Return(nil, fmt.Errorf("%s", test.expectedErr)).AnyTimes()
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
					route53Client.EXPECT().GetSubDomainDNSRecords(&validDomainOutput, ic, gomock.Any()).Return(returnedProblems, fmt.Errorf("%s", test.expectedErr)).AnyTimes()
				}
			} else {
				route53Client.EXPECT().GetSubDomainDNSRecords(&validDomainOutput, ic, gomock.Any()).Return(nil, nil).AnyTimes()
			}

			_, err := route53Client.GetSubDomainDNSRecords(&validDomainOutput, ic, nil)
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
	ic := icBuild.build()
	dottedClusterDomain := ic.ClusterDomain() + "."

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedResult, skipRecord(test.recordName, dottedClusterDomain))
		})
	}
}

func validAvailRegions() []string {
	return []string{"us-east-1", "us-central-1"}
}

func validAvailZones() []string {
	return []string{"a", "b", "c"}
}

func validEdgeAvailZones() []string {
	return []string{"edge-a", "edge-b", "edge-c"}
}

func validSubnets(subnetType string) Subnets {
	switch subnetType {
	case "edge":
		return Subnets{
			"subnet-valid-public-edge-a": {
				ID:     "subnet-valid-public-edge-a",
				Zone:   &Zone{Name: "edge-a"},
				CIDR:   "10.0.7.0/24",
				Public: true,
			},
			"subnet-valid-public-edge-b": {
				ID:     "subnet-valid-public-edge-b",
				Zone:   &Zone{Name: "edge-b"},
				CIDR:   "10.0.8.0/24",
				Public: true,
			},
			"subnet-valid-public-edge-c": {
				ID:     "subnet-valid-public-edge-c",
				Zone:   &Zone{Name: "edge-c"},
				CIDR:   "10.0.9.0/24",
				Public: true,
			},
		}
	case "public":
		return Subnets{
			"subnet-valid-public-a": {
				ID:     "subnet-valid-public-a",
				Zone:   &Zone{Name: "a"},
				CIDR:   "10.0.4.0/24",
				Public: true,
			},
			"subnet-valid-public-b": {
				ID:     "subnet-valid-public-b",
				Zone:   &Zone{Name: "b"},
				CIDR:   "10.0.5.0/24",
				Public: true,
			},
			"subnet-valid-public-c": {
				ID:     "subnet-valid-public-c",
				Zone:   &Zone{Name: "c"},
				CIDR:   "10.0.6.0/24",
				Public: true,
			},
		}
	case "private":
		return Subnets{
			"subnet-valid-private-a": {
				ID:     "subnet-valid-private-a",
				Zone:   &Zone{Name: "a"},
				CIDR:   "10.0.1.0/24",
				Public: false,
			},
			"subnet-valid-private-b": {
				ID:     "subnet-valid-private-b",
				Zone:   &Zone{Name: "b"},
				CIDR:   "10.0.2.0/24",
				Public: false,
			},
			"subnet-valid-private-c": {
				ID:     "subnet-valid-private-c",
				Zone:   &Zone{Name: "c"},
				CIDR:   "10.0.3.0/24",
				Public: false,
			},
		}
	}
	return nil
}

// byoSubnetsWithRoles returns a valid collection of subnets
// with assigned roles.
func byoSubnetsWithRoles() []aws.Subnet {
	return []aws.Subnet{
		{
			ID: "subnet-valid-private-a",
			Roles: []aws.SubnetRole{
				{Type: aws.ClusterNodeSubnetRole},
				{Type: aws.ControlPlaneInternalLBSubnetRole},
			},
		},
		{
			ID: "subnet-valid-private-b",
			Roles: []aws.SubnetRole{
				{Type: aws.ClusterNodeSubnetRole},
				{Type: aws.ControlPlaneInternalLBSubnetRole},
			},
		},
		{
			ID: "subnet-valid-private-c",
			Roles: []aws.SubnetRole{
				{Type: aws.ClusterNodeSubnetRole},
				{Type: aws.ControlPlaneInternalLBSubnetRole},
			},
		},
		{
			ID: "subnet-valid-public-a",
			Roles: []aws.SubnetRole{
				{Type: aws.ControlPlaneExternalLBSubnetRole},
				{Type: aws.IngressControllerLBSubnetRole},
				{Type: aws.BootstrapNodeSubnetRole},
			},
		},
		{
			ID: "subnet-valid-public-b",
			Roles: []aws.SubnetRole{
				{Type: aws.ControlPlaneExternalLBSubnetRole},
				{Type: aws.IngressControllerLBSubnetRole},
			},
		},
		{
			ID: "subnet-valid-public-c",
			Roles: []aws.SubnetRole{
				{Type: aws.ControlPlaneExternalLBSubnetRole},
				{Type: aws.IngressControllerLBSubnetRole},
			},
		},
	}
}

// byoPublicOnlySubnetsWithRoles returns a valid collection of subnets
// with assigned roles for a public-only cluster.
func byoPublicOnlySubnetsWithRoles() []aws.Subnet {
	return []aws.Subnet{
		{
			ID: "subnet-valid-public-a",
			Roles: []aws.SubnetRole{
				{Type: aws.ClusterNodeSubnetRole},
				{Type: aws.ControlPlaneInternalLBSubnetRole},
				{Type: aws.ControlPlaneExternalLBSubnetRole},
				{Type: aws.IngressControllerLBSubnetRole},
			},
		},
		{
			ID: "subnet-valid-public-b",
			Roles: []aws.SubnetRole{
				{Type: aws.ClusterNodeSubnetRole},
				{Type: aws.ControlPlaneInternalLBSubnetRole},
				{Type: aws.ControlPlaneExternalLBSubnetRole},
				{Type: aws.IngressControllerLBSubnetRole},
			},
		},
		{
			ID: "subnet-valid-public-c",
			Roles: []aws.SubnetRole{
				{Type: aws.BootstrapNodeSubnetRole},
			},
		},
	}
}

// byoEdgeSubnetsWithRoles returns a valid collection of edge subnets
// with assigned EdgeNode roles.
func byoEdgeSubnetsWithRoles() []aws.Subnet {
	return []aws.Subnet{
		{
			ID: "subnet-valid-public-edge-a",
			Roles: []aws.SubnetRole{
				{Type: aws.EdgeNodeSubnetRole},
			},
		},
		{
			ID: "subnet-valid-public-edge-b",
			Roles: []aws.SubnetRole{
				{Type: aws.EdgeNodeSubnetRole},
			},
		},
		{
			ID: "subnet-valid-public-edge-c",
			Roles: []aws.SubnetRole{
				{Type: aws.EdgeNodeSubnetRole},
			},
		},
	}
}

func otherTaggedPrivateSubnets() Subnets {
	return Subnets{
		"subnet-valid-private-a1": {
			ID:   "subnet-valid-private-a1",
			Zone: &Zone{Name: "a"},
			CIDR: "10.0.4.0/24",
			Tags: Tags{
				TagNameKubernetesClusterPrefix + "other-cluster": "owned",
			},
		},
		"subnet-valid-private-b1": {
			ID:   "subnet-valid-private-b1",
			Zone: &Zone{Name: "b"},
			CIDR: "10.0.5.0/24",
			Tags: Tags{
				TagNameKubernetesUnmanaged: "true",
			},
		},
	}
}

func otherUntaggedPrivateSubnets() Subnets {
	return Subnets{
		"subnet-valid-private-a1": {
			ID:   "subnet-valid-private-a1",
			Zone: &Zone{Name: "a"},
			CIDR: "10.0.6.0/24",
		},
		"subnet-valid-private-b1": {
			ID:   "subnet-valid-private-b1",
			Zone: &Zone{Name: "b"},
			CIDR: "10.0.7.0/24",
		},
	}
}

func validServiceEndpoints() []aws.ServiceEndpoint {
	return []aws.ServiceEndpoint{{
		Name: "ec2",
		URL:  "custom.ec2.us-east-1.amazonaws.com",
	}, {
		Name: "s3",
		URL:  "custom.s3.us-east-1.amazonaws.com",
	}, {
		Name: "iam",
		URL:  "custom.iam.us-east-1.amazonaws.com",
	}, {
		Name: "elasticloadbalancing",
		URL:  "custom.elasticloadbalancing.us-east-1.amazonaws.com",
	}, {
		Name: "tagging",
		URL:  "custom.tagging.us-east-1.amazonaws.com",
	}, {
		Name: "route53",
		URL:  "custom.route53.us-east-1.amazonaws.com",
	}, {
		Name: "sts",
		URL:  "custom.route53.us-east-1.amazonaws.com",
	}}
}

func invalidServiceEndpoint() []aws.ServiceEndpoint {
	return []aws.ServiceEndpoint{{
		Name: "ec3",
		URL:  "bad-aws-endpoint",
	}, {
		Name: "route55",
		URL:  "http://bad-aws-endpoint.non",
	}}
}

func validInstanceTypes() map[string]InstanceType {
	return map[string]InstanceType{
		"t2.small": {
			DefaultVCpus: 1,
			MemInMiB:     2048,
			Arches:       []string{ec2.ArchitectureTypeX8664},
		},
		"m5.large": {
			DefaultVCpus: 2,
			MemInMiB:     8192,
			Arches:       []string{ec2.ArchitectureTypeX8664},
		},
		"m5.xlarge": {
			DefaultVCpus: 4,
			MemInMiB:     16384,
			Arches:       []string{ec2.ArchitectureTypeX8664},
		},
		"m6g.xlarge": {
			DefaultVCpus: 4,
			MemInMiB:     16384,
			Arches:       []string{ec2.ArchitectureTypeArm64},
		},
	}
}

func createBaseDomainHostedZoneOutput() route53.HostedZone {
	return route53.HostedZone{
		CallerReference: &validCallerRef,
		Id:              &validDSId,
		Name:            &validDomainName,
	}
}

func createValidHostedZoneOutput() route53.GetHostedZoneOutput {
	ptrValidNameServers := []*string{}
	for i := range validNameServers {
		ptrValidNameServers = append(ptrValidNameServers, &validNameServers[i])
	}

	validDelegationSet := route53.DelegationSet{CallerReference: &validCallerRef, Id: &validDSId, NameServers: ptrValidNameServers}
	validHostedZone := route53.HostedZone{CallerReference: &validCallerRef, Id: &validDSId, Name: &validHostedZoneName}
	validVPCs := []*route53.VPC{{VPCId: &validVPCID, VPCRegion: &validAvailRegions()[0]}}

	return route53.GetHostedZoneOutput{
		DelegationSet: &validDelegationSet,
		HostedZone:    &validHostedZone,
		VPCs:          validVPCs,
	}
}

type icOption func(*types.InstallConfig)
type icBuildForAWS struct{}

var icBuild icBuildForAWS

func (icBuild icBuildForAWS) build(opts ...icOption) *types.InstallConfig {
	ic := &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: metaName,
		},
		BaseDomain: validDomainName,
		Publish:    types.ExternalPublishingStrategy,
		ControlPlane: &types.MachinePool{
			Architecture: types.ArchitectureAMD64,
			Replicas:     ptr.To[int64](3),
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{},
			},
		},
		Compute: []types.MachinePool{{
			Name:         types.MachinePoolComputeRoleName,
			Architecture: types.ArchitectureAMD64,
			Replicas:     ptr.To[int64](3),
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{},
			},
		}},
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Platform: types.Platform{
			AWS: &aws.Platform{
				Region: validAvailRegions()[0],
			},
		},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(ic)
		}
	}
	return ic
}

func (icBuild icBuildForAWS) withBaseBYO() icOption {
	return func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.AWS.Zones = validAvailZones()
		ic.Compute[0].Platform.AWS.Zones = validAvailZones()

		subnetIDs := append(validSubnets("public").IDs(), validSubnets("private").IDs()...)
		ic.AWS.VPC.Subnets = append(ic.AWS.VPC.Subnets, subnetsFromIDs(subnetIDs)...)
		ic.AWS.HostedZone = validHostedZoneName
	}
}

func (icBuild icBuildForAWS) withPublish(publish types.PublishingStrategy) icOption {
	return func(ic *types.InstallConfig) {
		ic.Publish = publish
	}
}

func (icBuild icBuildForAWS) withHostedZone(hostedZone string) icOption {
	return func(ic *types.InstallConfig) {
		ic.AWS.HostedZone = hostedZone
	}
}

func (icBuild icBuildForAWS) withBaseDomain(baseDomain string) icOption {
	return func(ic *types.InstallConfig) {
		ic.BaseDomain = baseDomain
	}
}

func (icBuild icBuildForAWS) withVPCSubnets(subnets []aws.Subnet, overwrite bool) icOption {
	return func(ic *types.InstallConfig) {
		if overwrite {
			ic.AWS.VPC.Subnets = subnets
		} else {
			ic.AWS.VPC.Subnets = append(ic.AWS.VPC.Subnets, subnets...)
		}
	}
}

func (icBuild icBuildForAWS) withVPCSubnetIDs(subnetIDs []string, overwrite bool) icOption {
	return func(ic *types.InstallConfig) {
		icBuild.withVPCSubnets(subnetsFromIDs(subnetIDs), overwrite)(ic)
	}
}

func (icBuild icBuildForAWS) withVPCEdgeSubnets(subnets []aws.Subnet, overwrite bool) icOption {
	return func(ic *types.InstallConfig) {
		icBuild.withVPCSubnets(subnets, overwrite)(ic)
		if len(subnets) > 0 {
			icBuild.withComputeMachinePool([]types.MachinePool{{
				Name: types.MachinePoolEdgeRoleName,
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				},
			}}, false)(ic)
		}
	}
}

func (icBuild icBuildForAWS) withVPCEdgeSubnetIDs(subnetIDs []string, overwrite bool) icOption {
	return func(ic *types.InstallConfig) {
		icBuild.withVPCEdgeSubnets(subnetsFromIDs(subnetIDs), overwrite)(ic)
	}
}

func (icBuild icBuildForAWS) withControlPlaneMachinePool(machinePool types.MachinePool) icOption {
	return func(ic *types.InstallConfig) {
		ic.ControlPlane = &machinePool
	}
}

func (icBuild icBuildForAWS) withComputeMachinePool(machinePools []types.MachinePool, overwrite bool) icOption {
	return func(ic *types.InstallConfig) {
		if overwrite {
			ic.Compute = machinePools
		} else {
			ic.Compute = append(ic.Compute, machinePools...)
		}
	}
}

func (icBuild icBuildForAWS) withControlPlanePlatformZones(zones []string, overwrite bool) icOption {
	return func(ic *types.InstallConfig) {
		if overwrite {
			ic.ControlPlane.Platform.AWS.Zones = zones
		} else {
			ic.ControlPlane.Platform.AWS.Zones = append(ic.ControlPlane.Platform.AWS.Zones, zones...)
		}
	}
}

func (icBuild icBuildForAWS) withComputePlatformZones(zones []string, overwrite bool, index int) icOption {
	return func(ic *types.InstallConfig) {
		if overwrite {
			ic.Compute[index].Platform.AWS.Zones = zones
		} else {
			ic.Compute[index].Platform.AWS.Zones = append(ic.Compute[index].Platform.AWS.Zones, zones...)
		}
	}
}

func (icBuild icBuildForAWS) withComputeHostPlacement(hostIDs []string, index int) icOption {
	return func(ic *types.InstallConfig) {
		aff := aws.HostAffinityDedicatedHost
		dhs := make([]aws.DedicatedHost, 0, len(hostIDs))
		for _, id := range hostIDs {
			dhs = append(dhs, aws.DedicatedHost{ID: id})
		}
		ic.Compute[index].Platform.AWS.HostPlacement = &aws.HostPlacement{
			Affinity:      &aff,
			DedicatedHost: dhs,
		}
	}
}

func (icBuild icBuildForAWS) withComputeHostPlacementAndZone(hostIDs []string, zone string, index int) icOption {
	return func(ic *types.InstallConfig) {
		aff := aws.HostAffinityDedicatedHost
		dhs := make([]aws.DedicatedHost, 0, len(hostIDs))
		for _, id := range hostIDs {
			dhs = append(dhs, aws.DedicatedHost{ID: id, Zone: zone})
		}
		ic.Compute[index].Platform.AWS.HostPlacement = &aws.HostPlacement{
			Affinity:      &aff,
			DedicatedHost: dhs,
		}
	}
}

func (icBuild icBuildForAWS) withControlPlanePlatformAMI(amiID string) icOption {
	return func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.AWS.AMIID = amiID
	}
}

func (icBuild icBuildForAWS) withComputePlatformAMI(amiID string, index int) icOption {
	return func(ic *types.InstallConfig) {
		ic.Compute[index].Platform.AWS.AMIID = amiID
	}
}

func (icBuild icBuildForAWS) withComputeReplicas(replicas int64, index int) icOption {
	return func(ic *types.InstallConfig) {
		ic.Compute[index].Replicas = ptr.To(replicas)
	}
}

func (icBuild icBuildForAWS) withInstanceType(defaultInstanceType string, ctrPlaneInstanceType string, computeInstanceTypes ...string) icOption {
	return func(ic *types.InstallConfig) {
		if ic.Platform.AWS.DefaultMachinePlatform == nil {
			icBuild.withDefaultPlatformMachine(aws.MachinePool{})(ic)
		}
		ic.Platform.AWS.DefaultMachinePlatform.InstanceType = defaultInstanceType
		ic.ControlPlane.Platform.AWS.InstanceType = ctrPlaneInstanceType
		for idx, instanceType := range computeInstanceTypes {
			ic.Compute[idx].Platform.AWS.InstanceType = instanceType
		}
	}
}

func (icBuild icBuildForAWS) withInstanceArchitecture(ctrPlaneInstanceArch types.Architecture, computeInstanceArchs ...types.Architecture) icOption {
	return func(ic *types.InstallConfig) {
		ic.ControlPlane.Architecture = ctrPlaneInstanceArch
		for idx, arch := range computeInstanceArchs {
			ic.Compute[idx].Architecture = arch
		}
	}
}

func (icBuild icBuildForAWS) withPlatformRegion(region string) icOption {
	return func(ic *types.InstallConfig) {
		ic.Platform.AWS.Region = region
	}
}

func (icBuild icBuildForAWS) withPlatformAMIID(amiID string) icOption {
	return func(ic *types.InstallConfig) {
		ic.Platform.AWS.AMIID = amiID
	}
}

func (icBuild icBuildForAWS) withServiceEndpoints(endpoints []aws.ServiceEndpoint, overwrite bool) icOption {
	return func(ic *types.InstallConfig) {
		if overwrite {
			ic.Platform.AWS.ServiceEndpoints = endpoints
		} else {
			ic.Platform.AWS.ServiceEndpoints = append(ic.Platform.AWS.ServiceEndpoints, endpoints...)
		}
	}
}

func (icBuild icBuildForAWS) withDefaultPlatformMachine(awsMachine aws.MachinePool) icOption {
	return func(ic *types.InstallConfig) {
		ic.Platform.AWS.DefaultMachinePlatform = &awsMachine
	}
}

func (icBuild icBuildForAWS) withPublicIPv4Pool(publicIPv4Pool string) icOption {
	return func(ic *types.InstallConfig) {
		ic.Platform.AWS.PublicIpv4Pool = publicIPv4Pool
	}
}
