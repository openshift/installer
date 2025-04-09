package validation

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *aws.Platform
		expected string
		credMode types.CredentialsMode
		publish  types.PublishingStrategy
	}{
		{
			name: "minimal",
			platform: &aws.Platform{
				Region: "us-east-1",
			},
		},
		{
			name: "invalid region",
			platform: &aws.Platform{
				Region: "",
			},
			expected: `^test-path\.region: Required value: region must be specified$`,
		},
		{
			name: "hosted zone with subnets",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "subnet-1234567890asdfghj"},
					},
				},
				HostedZone: "test-hosted-zone",
			},
		},
		{
			name: "hosted zone without subnets",
			platform: &aws.Platform{
				Region:     "us-east-1",
				HostedZone: "test-hosted-zone",
			},
			expected: `^test-path\.hostedZone: Invalid value: "test-hosted-zone": may not use an existing hosted zone when not using existing subnets$`,
		},
		{
			name: "invalid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "/path/some",
				}},
			},
			expected: `^test-path\.serviceEndpoints\[0\]\.url: Invalid value: "(.*)": host cannot be empty, empty host provided$`,
		},
		{
			name: "invalid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "https://test-ec2.random.local/path/some",
				}},
			},
			expected: `^test-path\.serviceEndpoints\[0\]\.url: Invalid value: "(.*)": no path or request parameters must be provided, "/path/some" was provided$`,
		},
		{
			name: "invalid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "https://test-ec2.random.local?foo=some",
				}},
			},
			expected: `^test-path\.serviceEndpoints\[0\]\.url: Invalid value: "(.*)": no path or request parameters must be provided, "/\?foo=some" was provided$`,
		},
		{
			name: "valid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "test-ec2.random.local",
				}},
			},
		},
		{
			name: "valid url for service endpoint",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "https://test-ec2.random.local",
				}},
			},
		},
		{
			name: "duplicate service endpoints",
			platform: &aws.Platform{
				Region: "us-east-1",
				ServiceEndpoints: []aws.ServiceEndpoint{{
					Name: "ec2",
					URL:  "test-ec2.random.local",
				}, {
					Name: "s3",
					URL:  "test-ec2.random.local",
				}, {
					Name: "ec2",
					URL:  "test-ec2.random.local",
				}},
			},
			expected: `^test-path\.serviceEndpoints\[2\]\.name: Invalid value: "ec2": duplicate service endpoint not allowed for ec2, service endpoint already defined at test-path\.serviceEndpoints\[0\]$`,
		},
		{
			name: "valid machine pool",
			platform: &aws.Platform{
				Region:                 "us-east-1",
				DefaultMachinePlatform: &aws.MachinePool{},
			},
		},
		{
			name: "invalid machine pool",
			platform: &aws.Platform{
				Region: "us-east-1",
				DefaultMachinePlatform: &aws.MachinePool{
					EC2RootVolume: aws.EC2RootVolume{
						Type: "io1",
						IOPS: -10,
						Size: 128,
					},
				},
			},
			expected: `^test-path.*iops: Invalid value: -10: iops must be a positive number$`,
		},
		{
			name: "invalid userTags, Name key",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"Name": "test-cluster",
				},
			},
			expected: `^\Qtest-path.userTags[Name]: Invalid value: "test-cluster": "Name" key is not allowed for user defined tags\E$`,
		},
		{
			name: "invalid userTags, key with kubernetes.io/cluster/",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"kubernetes.io/cluster/test-cluster": "shared",
				},
			},
			expected: `^\Qtest-path.userTags[kubernetes.io/cluster/test-cluster]: Invalid value: "shared": key is in the kubernetes.io namespace\E$`,
		},
		{
			name: "invalid userTags, value with invalid characters",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"usage-user": "cloud-team-rebase-bot[bot]",
				},
			},
			expected: `^\Qtest-path.userTags[usage-user]: Invalid value: "cloud-team-rebase-bot[bot]": value contains invalid characters`,
		},
		{
			name: "valid userTags, value with spaces",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"test-key": "this test has spaces",
				},
				PropagateUserTag: true,
			},
		},
		{
			name: "valid userTags",
			platform: &aws.Platform{
				Region: "us-east-1",
				UserTags: map[string]string{
					"app": "production",
				},
			},
		},
		{
			name: "too many userTags",
			platform: &aws.Platform{
				Region:   "us-east-1",
				UserTags: generateTooManyUserTags(),
			},
			expected: fmt.Sprintf(`^\Qtest-path.userTags: Too many: %d: must have at most %d items`, userTagLimit+1, userTagLimit),
		},
		{
			name: "hosted zone role without hosted zone should error",
			platform: &aws.Platform{
				Region:         "us-east-1",
				HostedZoneRole: "test-hosted-zone-role",
			},
			expected: `\Qtest-path.hostedZoneRole: Invalid value: "test-hosted-zone-role": may not specify a role to assume for hosted zone operations without also specifying a hosted zone\E`,
		},
		{
			name:     "valid hosted zone & role should not throw an error",
			credMode: types.PassthroughCredentialsMode,
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "subnet-1234567890asdfghj"},
					},
				},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
			},
		},
		{
			name: "hosted zone role without credential mode should error",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "subnet-1234567890asdfghj"},
					},
				},
				HostedZone:     "test-hosted-zone",
				HostedZoneRole: "test-hosted-zone-role",
			},
			expected: `^\Qtest-path.credentialsMode: Forbidden: when specifying a hostedZoneRole, either Passthrough or Manual credential mode must be specified\E$`,
		},
		{
			name: "valid subnets, empty",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC:    aws.VPC{},
			},
		},
		{
			name: "valid subnets, no roles assigned",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "subnet-1234567890asdfghj"},
						{ID: "subnet-asdfghj1234567890"},
					},
				},
			},
		},
		{
			name:    "valid subnets, internal cluster and all required roles assigned",
			publish: types.InternalPublishingStrategy,
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-asdfghj1234567890",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneInternalLBSubnetRole},
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					},
				},
			},
		},
		{
			name:    "valid subnets, external cluster and all required roles assigned",
			publish: types.ExternalPublishingStrategy,
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-asdfghj1234567890",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneInternalLBSubnetRole},
							},
						},
						{
							ID: "subnet-0fcf8e0392f0910d0",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneExternalLBSubnetRole},
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					},
				},
			},
		},
		{
			name: "invalid subnets, some without roles assigned",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-asdfghj1234567890",
						},
					},
				},
			},
			expected: `^\Q[test-path.vpc.subnets: Forbidden: either all subnets must be assigned roles or none of the subnets should have roles assigned, test-path.vpc.subnets: Invalid value: []aws.Subnet{aws.Subnet{ID:"subnet-1234567890asdfghj", Roles:[]aws.SubnetRole{aws.SubnetRole{Type:"BootstrapNode"}, aws.SubnetRole{Type:"ClusterNode"}}}, aws.Subnet{ID:"subnet-asdfghj1234567890", Roles:[]aws.SubnetRole(nil)}}: roles [ControlPlaneExternalLB ControlPlaneInternalLB IngressControllerLB] must be assigned to at least 1 subnet]\E$`,
		},
		{
			name: "invalid subnets, duplicate subnet IDs",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "subnet-1234567890asdfghj"},
						{ID: "subnet-1234567890asdfghj"},
					},
				},
			},
			expected: `^test-path\.vpc\.subnets\[1\]\.id: Duplicate value: \"subnet-1234567890asdfghj\"$`,
		},
		{
			name: "invalid subnets, unsupported roles assigned",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: "UnsupportedRole"},
							},
						},
					},
				},
			},
			expected: `^\Qtest-path.vpc.subnets[0].roles[0].type: Unsupported value: "UnsupportedRole": supported values: "BootstrapNode", "ClusterNode", "ControlPlaneExternalLB", "ControlPlaneInternalLB", "EdgeNode", "IngressControllerLB"\E$`,
		},
		{
			name: "invalid subnets, duplicate roles assigned to a subnet",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-asdfghj1234567890",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneInternalLBSubnetRole},
							},
						},
						{
							ID: "subnet-0fcf8e0392f0910d0",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneExternalLBSubnetRole},
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					},
				},
			},
			expected: `^\Qtest-path.vpc.subnets[0].roles[1].type: Duplicate value: "BootstrapNode"\E$`,
		},
		{
			name: "invalid subnets, ControlPlaneExternalLB and ControlPlaneInternalLB roles assigned to the same subnet",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-0fcf8e0392f0910d0",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneExternalLBSubnetRole},
								{Type: aws.ControlPlaneInternalLBSubnetRole},
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					},
				},
			},
			expected: `^test-path\.vpc\.subnets\[1\]\.roles: Forbidden: must not have both ControlPlaneExternalLB and ControlPlaneInternalLB role$`,
		},
		{
			name: "invalid subnets, EdgeNode and other roles assigned to the same subnet",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.EdgeNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-asdfghj1234567890",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneInternalLBSubnetRole},
							},
						},
						{
							ID: "subnet-0fcf8e0392f0910d0",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneExternalLBSubnetRole},
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					},
				},
			},
			expected: `^test-path\.vpc\.subnets\[0\]\.roles: Forbidden: must not combine EdgeNode role with any other roles$`,
		},
		{
			name: "invalid subnets, IngressControllerLB roles assigned to more than 10 subnets",
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: append(
						[]aws.Subnet{
							{
								ID: "subnet-1234567890asdfghj",
								Roles: []aws.SubnetRole{
									{Type: aws.BootstrapNodeSubnetRole},
									{Type: aws.ClusterNodeSubnetRole},
								},
							},
							{
								ID: "subnet-asdfghj1234567890",
								Roles: []aws.SubnetRole{
									{Type: aws.ControlPlaneInternalLBSubnetRole},
								},
							},
							{
								ID: "subnet-0fcf8e0392f0910d0",
								Roles: []aws.SubnetRole{
									{Type: aws.ControlPlaneExternalLBSubnetRole},
								},
							},
						},
						func() []aws.Subnet {
							subnets := []aws.Subnet{}
							for i := 0; i < 11; i++ {
								subnets = append(subnets, aws.Subnet{
									ID: aws.AWSSubnetID(fmt.Sprintf("subnet-asdfghj392f0910d%d", i)),
									Roles: []aws.SubnetRole{
										{Type: aws.IngressControllerLBSubnetRole},
									},
								})
							}
							return subnets
						}()...,
					),
				},
			},
			expected: `^test-path\.vpc\.subnets: Forbidden: must not include more than 10 subnets with the IngressControllerLB role$`,
		},
		{
			name:    "invalid subnets, internal cluster and ControlPlaneExternalLB role assigned to a subnet",
			publish: types.InternalPublishingStrategy,
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-asdfghj1234567890",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneInternalLBSubnetRole},
							},
						},
						{
							ID: "subnet-0fcf8e0392f0910d0",
							Roles: []aws.SubnetRole{
								{Type: aws.ControlPlaneExternalLBSubnetRole},
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					},
				},
			},
			expected: `^test-path\.vpc\.subnets: Forbidden: must not include subnets with the ControlPlaneExternalLB role in a private cluster$`,
		},
		{
			name:    "invalid subnets, internal cluster and missing required roles",
			publish: types.InternalPublishingStrategy,
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-0fcf8e0392f0910d0",
							Roles: []aws.SubnetRole{
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					},
				},
			},
			expected: `^\Qtest-path.vpc.subnets: Invalid value: []aws.Subnet{aws.Subnet{ID:"subnet-1234567890asdfghj", Roles:[]aws.SubnetRole{aws.SubnetRole{Type:"BootstrapNode"}, aws.SubnetRole{Type:"ClusterNode"}}}, aws.Subnet{ID:"subnet-0fcf8e0392f0910d0", Roles:[]aws.SubnetRole{aws.SubnetRole{Type:"IngressControllerLB"}}}}: roles [ControlPlaneInternalLB] must be assigned to at least 1 subnet\E$`,
		},
		{
			name:    "invalid subnets, external cluster and missing required roles",
			publish: types.ExternalPublishingStrategy,
			platform: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{
							ID: "subnet-1234567890asdfghj",
							Roles: []aws.SubnetRole{
								{Type: aws.BootstrapNodeSubnetRole},
								{Type: aws.ClusterNodeSubnetRole},
							},
						},
						{
							ID: "subnet-0fcf8e0392f0910d0",
							Roles: []aws.SubnetRole{
								{Type: aws.IngressControllerLBSubnetRole},
							},
						},
					},
				},
			},
			expected: `^\Qtest-path.vpc.subnets: Invalid value: []aws.Subnet{aws.Subnet{ID:"subnet-1234567890asdfghj", Roles:[]aws.SubnetRole{aws.SubnetRole{Type:"BootstrapNode"}, aws.SubnetRole{Type:"ClusterNode"}}}, aws.Subnet{ID:"subnet-0fcf8e0392f0910d0", Roles:[]aws.SubnetRole{aws.SubnetRole{Type:"IngressControllerLB"}}}}: roles [ControlPlaneExternalLB ControlPlaneInternalLB] must be assigned to at least 1 subnet\E$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, tc.publish, tc.credMode, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

func TestValidateTag(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		value     string
		expectErr bool
	}{{
		name:  "valid",
		key:   "test-key",
		value: "test-value",
	}, {
		name:      "invalid characters in key",
		key:       "bad-key***",
		value:     "test-value",
		expectErr: true,
	}, {
		name:      "invalid characters in value",
		key:       "test-key",
		value:     "bad-value***",
		expectErr: true,
	}, {
		name:      "empty key",
		key:       "",
		value:     "test-value",
		expectErr: true,
	}, {
		name:      "empty value",
		key:       "test-key",
		value:     "",
		expectErr: true,
	}, {
		name:      "key too long",
		key:       strings.Repeat("a", 129),
		value:     "test-value",
		expectErr: true,
	}, {
		name:      "value too long",
		key:       "test-key",
		value:     strings.Repeat("a", 257),
		expectErr: true,
	}, {
		name:      "key in kubernetes.io namespace",
		key:       "kubernetes.io/cluster/some-cluster",
		value:     "owned",
		expectErr: true,
	}, {
		name:      "key in openshift.io namespace",
		key:       "openshift.io/some-key",
		value:     "some-value",
		expectErr: true,
	}, {
		name:      "key in openshift.io subdomain namespace",
		key:       "other.openshift.io/some-key",
		value:     "some-value",
		expectErr: true,
	}, {
		name:  "key in namespace similar to openshift.io",
		key:   "otheropenshift.io/some-key",
		value: "some-value",
	}, {
		name:  "key with openshift.io in path",
		key:   "some-domain/openshift.io/some-key",
		value: "some-value",
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTag(tc.key, tc.value)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func generateTooManyUserTags() map[string]string {
	tags := map[string]string{}
	for i := 0; i <= userTagLimit; i++ {
		tags[strconv.Itoa(i)] = strconv.Itoa(i)
	}
	return tags
}
