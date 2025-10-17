package aws

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

var stubDefaultCIDR = "10.0.0.0/16"

func stubClusterID() *installconfig.ClusterID {
	return &installconfig.ClusterID{
		InfraID: "infra-id",
	}
}

func stubInstallConfig() *installconfig.InstallConfig {
	return &installconfig.InstallConfig{}
}

func stubInstallConfigType() *types.InstallConfig {
	return &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		Publish: types.ExternalPublishingStrategy,
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: *ipnet.MustParseCIDR(stubDefaultCIDR),
				},
			},
		},
	}
}

func stubInstallConfigPoolCompute() []types.MachinePool {
	return []types.MachinePool{
		{
			Name: "worker",
			Platform: types.MachinePoolPlatform{
				AWS: &awstypes.MachinePool{
					Zones: []string{"b", "c"},
				},
			},
		},
	}
}

func stubInstallConfigPoolComputeWithEdge() []types.MachinePool {
	p := stubInstallConfigPoolCompute()
	p = append(p, types.MachinePool{
		Name: "edge",
		Platform: types.MachinePoolPlatform{
			AWS: &awstypes.MachinePool{
				Zones: []string{"edge-b", "edge-c"},
			},
		},
	})
	return p
}

func stubInstallConfigPoolControl() *types.MachinePool {
	return &types.MachinePool{
		Name: "master",
		Platform: types.MachinePoolPlatform{
			AWS: &awstypes.MachinePool{
				Zones: []string{"a", "b"},
			},
		},
	}
}

func tSortCapaSubnetsByID(in capa.Subnets) capa.Subnets {
	subnetIDs := []string{}
	subnetsMap := make(map[string]capa.SubnetSpec, len(in))
	for _, subnet := range in {
		subnetsMap[subnet.ID] = subnet
		subnetIDs = append(subnetIDs, subnet.ID)
	}
	sort.Strings(subnetIDs)
	out := capa.Subnets{}
	for _, sid := range subnetIDs {
		out = append(out, subnetsMap[sid])
	}
	return out
}

func Test_extractZonesFromInstallConfig(t *testing.T) {
	type args struct {
		in *networkInput
	}
	tests := []struct {
		name       string
		args       args
		want       *ZonesCAPI
		wantErrMsg string
	}{
		{
			name: "no zones in config, use default from region",
			args: args{
				in: &networkInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = stubInstallConfigType()
						ic.Config.AWS = &awstypes.Platform{
							DefaultMachinePlatform: &awstypes.MachinePool{
								Zones: []string{},
							},
						}
						return ic
					}(),
					ZonesInRegion: []string{"x", "y"},
				},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("x", "y"),
				ComputeZones:      sets.New("x", "y"),
				EdgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "no zones in config pools, use default from platform config",
			args: args{
				in: &networkInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = stubInstallConfigType()
						ic.Config.AWS = &awstypes.Platform{
							DefaultMachinePlatform: &awstypes.MachinePool{
								Zones: []string{"a", "b"},
							},
						}
						return ic
					}(),
					ZonesInRegion: []string{"x", "y"},
				},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("a", "b"),
				ComputeZones:      sets.New("a", "b"),
				EdgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "custom zones control plane pool",
			args: args{
				in: &networkInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: stubInstallConfigPoolControl(),
							Compute:      nil,
						}
						return ic
					}(),
					ZonesInRegion: []string{"x", "y"},
				},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("a", "b"),
				ComputeZones:      sets.New("x", "y"),
				EdgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "custom zones compute pool",
			args: args{
				in: &networkInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: nil,
							Compute:      stubInstallConfigPoolCompute(),
						}
						return ic
					}(),
					ZonesInRegion: []string{"x", "y"},
				},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("x", "y"),
				ComputeZones:      sets.New("b", "c"),
				EdgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "custom zones control plane and compute pools",
			args: args{
				in: &networkInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: stubInstallConfigPoolControl(),
							Compute:      stubInstallConfigPoolCompute(),
						}
						return ic
					}(),
					ZonesInRegion: []string{"x", "y"},
				},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("a", "b"),
				ComputeZones:      sets.New("b", "c"),
				EdgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "custom zones control plane, compute and edge pools",
			args: args{
				in: &networkInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: stubInstallConfigPoolControl(),
							Compute:      stubInstallConfigPoolComputeWithEdge(),
						}
						return ic
					}(),
					ZonesInRegion: []string{"x", "y"},
				},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("a", "b"),
				ComputeZones:      sets.New("b", "c"),
				EdgeZones:         sets.New("edge-b", "edge-c"),
			},
		},
		// errors
		{
			name: "unexpected empty zones on config and metadata",
			args: args{
				in: &networkInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = stubInstallConfigType()
						ic.Config.AWS = &awstypes.Platform{
							DefaultMachinePlatform: &awstypes.MachinePool{
								Zones: []string{},
							},
						}
						return ic
					}(),
					ZonesInRegion: []string{},
				},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
				EdgeZones:         sets.Set[string]{},
			},
			wantErrMsg: `failed to set zones from config, got: []`,
		},
		{
			name: "unexpected empty zones from edge compute pool",
			args: args{
				in: &networkInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: stubInstallConfigPoolControl(),
							Compute: func() []types.MachinePool {
								pools := stubInstallConfigPoolCompute()
								// create empty zones' edge pool to force failures
								pools = append(pools, types.MachinePool{
									Name: "edge",
									Platform: types.MachinePoolPlatform{
										AWS: &awstypes.MachinePool{
											Zones: []string{},
										},
									},
								})
								return pools
							}(),
						}
						return ic
					}(),
					ZonesInRegion: []string{"x", "y"},
				},
			},
			wantErrMsg: `expect one or more zones in the edge compute pool, got: []`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractZonesFromInstallConfig(tt.args.in)
			if err != nil {
				if len(tt.wantErrMsg) > 0 {
					if got := err.Error(); !cmp.Equal(got, tt.wantErrMsg) {
						t.Errorf("extractZonesFromInstallConfig() unexpected error message: %v", cmp.Diff(got, tt.wantErrMsg))
					}
					return
				}
				t.Errorf("extractZonesFromInstallConfig() unexpected error: %v", err)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("extractZonesFromInstallConfig() Got unexpected results:\n%v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_setSubnetsManagedVPC(t *testing.T) {
	type args struct {
		in *networkInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    capa.NetworkSpec
	}{
		{
			name: "default availability zones in the region",
			args: args{
				in: &networkInput{
					ClusterID: stubClusterID(),
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							Publish: types.ExternalPublishingStrategy,
							Networking: &types.Networking{
								MachineNetwork: []types.MachineNetworkEntry{
									{
										CIDR: *ipnet.MustParseCIDR(stubDefaultCIDR),
									},
								},
							},
						}
						return ic
					}(),
					Cluster: &capa.AWSCluster{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "infraId",
							Namespace: capiutils.Namespace,
						},
						Spec: capa.AWSClusterSpec{},
					},
					ZonesInRegion: []string{"a", "b", "c"},
				},
			},
			want: capa.NetworkSpec{
				VPC: capa.VPCSpec{CidrBlock: stubDefaultCIDR},
				Subnets: []capa.SubnetSpec{
					{
						ID:               "infra-id-subnet-private-a",
						AvailabilityZone: "a",
						IsPublic:         false,
						CidrBlock:        "10.0.0.0/19",
					}, {
						ID:               "infra-id-subnet-private-b",
						AvailabilityZone: "b",
						IsPublic:         false,
						CidrBlock:        "10.0.32.0/19",
					}, {
						ID:               "infra-id-subnet-private-c",
						AvailabilityZone: "c",
						IsPublic:         false,
						CidrBlock:        "10.0.64.0/19",
					}, {
						ID:               "infra-id-subnet-public-a",
						AvailabilityZone: "a",
						IsPublic:         true,
						CidrBlock:        "10.0.96.0/21",
					}, {
						ID:               "infra-id-subnet-public-b",
						AvailabilityZone: "b",
						IsPublic:         true,
						CidrBlock:        "10.0.104.0/21",
					}, {
						ID:               "infra-id-subnet-public-c",
						AvailabilityZone: "c",
						IsPublic:         true,
						CidrBlock:        "10.0.112.0/21",
					},
				},
			},
		},
		{
			name: "default availability zones in the region and edge",
			args: args{
				in: &networkInput{
					ClusterID: stubClusterID(),
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							Publish: types.ExternalPublishingStrategy,
							Networking: &types.Networking{
								MachineNetwork: []types.MachineNetworkEntry{
									{
										CIDR: *ipnet.MustParseCIDR(stubDefaultCIDR),
									},
								},
							},
							Compute: []types.MachinePool{
								{
									Name: "edge",
									Platform: types.MachinePoolPlatform{
										AWS: &awstypes.MachinePool{
											Zones: []string{"edge-a"},
										},
									},
								},
							},
						}
						return ic
					}(),
					Cluster: &capa.AWSCluster{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "infraId",
							Namespace: capiutils.Namespace,
						},
						Spec: capa.AWSClusterSpec{},
					},
					ZonesInRegion: []string{"a", "b", "c"},
				},
			},
			want: capa.NetworkSpec{
				VPC: capa.VPCSpec{CidrBlock: stubDefaultCIDR},
				Subnets: []capa.SubnetSpec{
					{
						ID:               "infra-id-subnet-private-a",
						AvailabilityZone: "a",
						IsPublic:         false,
						CidrBlock:        "10.0.0.0/19",
					}, {
						ID:               "infra-id-subnet-private-b",
						AvailabilityZone: "b",
						IsPublic:         false,
						CidrBlock:        "10.0.32.0/19",
					}, {
						ID:               "infra-id-subnet-private-c",
						AvailabilityZone: "c",
						IsPublic:         false,
						CidrBlock:        "10.0.64.0/19",
					}, {
						ID:               "infra-id-subnet-private-edge-a",
						AvailabilityZone: "edge-a",
						IsPublic:         false,
						CidrBlock:        "10.0.128.0/21",
					}, {
						ID:               "infra-id-subnet-public-a",
						AvailabilityZone: "a",
						IsPublic:         true,
						CidrBlock:        "10.0.96.0/21",
					}, {
						ID:               "infra-id-subnet-public-b",
						AvailabilityZone: "b",
						IsPublic:         true,
						CidrBlock:        "10.0.104.0/21",
					}, {
						ID:               "infra-id-subnet-public-c",
						AvailabilityZone: "c",
						IsPublic:         true,
						CidrBlock:        "10.0.112.0/21",
					}, {
						ID:               "infra-id-subnet-public-edge-a",
						AvailabilityZone: "edge-a",
						IsPublic:         true,
						CidrBlock:        "10.0.136.0/21",
					},
				},
			},
		},
		// TODO: error scenarios to review the coverage
		// {
		// 	name: "err: failed to get availability zones: expect one or more zones in the edge compute pool",
		// },
		// {
		// 	name: "err: failed to get availability zones: failed to set zones from config",
		// },
		// {
		// 	name: "err: unable to generate CIDR blocks for all private subnets",
		// },
		// {
		// 	name: "err: unable to generate CIDR blocks for all public subnets",
		// },
		// {
		// 	name: "err: unable to define CIDR blocks to all zones for private subnets",
		// },
		// {
		// 	name: "err: unable to define CIDR blocks to all zones for public subnets",
		// },
		// {
		// 	name: "err: unable to generate CIDR blocks for all edge subnets",
		// },
		// {
		// 	name: "err: unable to define CIDR blocks to all edge zones for private subnets",
		// },
		// {
		// 	name: "err: unable to define CIDR blocks to all edge zones for public subnets",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setSubnetsManagedVPC(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("setSubnetsManagedVPC() #1 error: %+v,\nwantErr %+v\n", err, tt.wantErr)
			}
			if tt.args.in == nil && (err == nil) {
				return
			}
			if tt.args.in.Cluster == nil && (err == nil) {
				return
			}

			if len(tt.args.in.Cluster.Spec.NetworkSpec.Subnets) == 0 {
				if !tt.wantErr {
					t.Errorf("setSubnetsManagedVPC() #2 error: %v, wantErr: %v", err, tt.wantErr)
				}
				return
			}
			tt.args.in.Cluster.Spec.NetworkSpec.Subnets = tSortCapaSubnetsByID(tt.args.in.Cluster.Spec.NetworkSpec.Subnets)
			if got := tt.args.in.Cluster.Spec.NetworkSpec; !cmp.Equal(got, tt.want) {
				t.Errorf("setSubnetsManagedVPC() NetworkSpec.Subnets diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_setSubnetsBYOVPC(t *testing.T) {
	type args struct {
		in *networkInput
	}
	tests := []struct {
		name    string
		args    args
		want    capa.NetworkSpec
		wantErr bool
	}{
		{
			name: "default byo vpc",
			args: args{
				in: &networkInput{
					Cluster: &capa.AWSCluster{},
					Subnets: &subnetsInput{
						vpc: "vpc-id",
						privateSubnets: aws.Subnets{
							"subnetId-a-private": aws.Subnet{
								ID:   "subnetId-a-private",
								CIDR: "10.0.1.0/24",
								Zone: &aws.Zone{
									Name: "a",
								},
								Public: false,
							},
							"subnetId-b-private": aws.Subnet{
								ID:   "subnetId-b-private",
								CIDR: "10.0.2.0/24",
								Zone: &aws.Zone{
									Name: "b",
								},
								Public: false,
							},
							"subnetId-c-private": aws.Subnet{
								ID:   "subnetId-c-private",
								CIDR: "10.0.3.0/24",
								Zone: &aws.Zone{
									Name: "c",
								},
								Public: false,
							},
						},
						publicSubnets: aws.Subnets{
							"subnetId-a-public": aws.Subnet{
								ID:   "subnetId-a-public",
								CIDR: "10.0.4.0/24",
								Zone: &aws.Zone{
									Name: "a",
								},
								Public: true,
							},
							"subnetId-b-public": aws.Subnet{
								ID:   "subnetId-b-public",
								CIDR: "10.0.5.0/24",
								Zone: &aws.Zone{
									Name: "b",
								},
								Public: true,
							},
							"subnetId-c-public": aws.Subnet{
								ID:   "subnetId-c-public",
								CIDR: "10.0.6.0/24",
								Zone: &aws.Zone{
									Name: "c",
								},
								Public: true,
							},
						},
					},
				},
			},
			want: capa.NetworkSpec{
				VPC: capa.VPCSpec{ID: "vpc-id"},
				Subnets: []capa.SubnetSpec{
					{
						ID:               "subnetId-a-private",
						AvailabilityZone: "a",
						IsPublic:         false,
						CidrBlock:        "10.0.1.0/24",
					}, {
						ID:               "subnetId-a-public",
						AvailabilityZone: "a",
						IsPublic:         true,
						CidrBlock:        "10.0.4.0/24",
					}, {
						ID:               "subnetId-b-private",
						AvailabilityZone: "b",
						IsPublic:         false,
						CidrBlock:        "10.0.2.0/24",
					}, {
						ID:               "subnetId-b-public",
						AvailabilityZone: "b",
						IsPublic:         true,
						CidrBlock:        "10.0.5.0/24",
					}, {
						ID:               "subnetId-c-private",
						AvailabilityZone: "c",
						IsPublic:         false,
						CidrBlock:        "10.0.3.0/24",
					}, {
						ID:               "subnetId-c-public",
						AvailabilityZone: "c",
						IsPublic:         true,
						CidrBlock:        "10.0.6.0/24",
					},
				},
			},
		},
		{
			name: "byo vpc only private subnets",
			args: args{
				in: &networkInput{
					Cluster: &capa.AWSCluster{},
					Subnets: &subnetsInput{
						vpc: "vpc-id",
						privateSubnets: aws.Subnets{
							"subnetId-a-private": aws.Subnet{
								ID:   "subnetId-a-private",
								CIDR: "10.0.1.0/24",
								Zone: &aws.Zone{
									Name: "a",
								},
								Public: false,
							},
							"subnetId-b-private": aws.Subnet{
								ID:   "subnetId-b-private",
								CIDR: "10.0.2.0/24",
								Zone: &aws.Zone{
									Name: "b",
								},
								Public: false,
							},
							"subnetId-c-private": aws.Subnet{
								ID:   "subnetId-c-private",
								CIDR: "10.0.3.0/24",
								Zone: &aws.Zone{
									Name: "c",
								},
								Public: false,
							},
						},
					},
				},
			},
			want: capa.NetworkSpec{
				VPC: capa.VPCSpec{
					ID: "vpc-id",
				},
				Subnets: capa.Subnets{
					{
						ID:               "subnetId-a-private",
						AvailabilityZone: "a",
						IsPublic:         false,
						CidrBlock:        "10.0.1.0/24",
					},
					{
						ID:               "subnetId-b-private",
						AvailabilityZone: "b",
						IsPublic:         false,
						CidrBlock:        "10.0.2.0/24",
					},
					{
						ID:               "subnetId-c-private",
						AvailabilityZone: "c",
						IsPublic:         false,
						CidrBlock:        "10.0.3.0/24",
					},
				},
			},
		},
		{
			name: "byo vpc with edge",
			args: args{
				in: &networkInput{
					Cluster: &capa.AWSCluster{},
					Subnets: &subnetsInput{
						vpc: "vpc-id",
						privateSubnets: aws.Subnets{
							"subnetId-a-private": aws.Subnet{
								ID:   "subnetId-a-private",
								CIDR: "10.0.1.0/24",
								Zone: &aws.Zone{
									Name: "a",
								},
								Public: false,
							},
							"subnetId-b-private": aws.Subnet{
								ID:   "subnetId-b-private",
								CIDR: "10.0.2.0/24",
								Zone: &aws.Zone{
									Name: "b",
								},
								Public: false,
							},
							"subnetId-c-private": aws.Subnet{
								ID:   "subnetId-c-private",
								CIDR: "10.0.3.0/24",
								Zone: &aws.Zone{
									Name: "c",
								},
								Public: false,
							},
						},
						publicSubnets: aws.Subnets{
							"subnetId-a-public": aws.Subnet{
								ID:   "subnetId-a-public",
								CIDR: "10.0.4.0/24",
								Zone: &aws.Zone{
									Name: "a",
								},
								Public: true,
							},
							"subnetId-b-public": aws.Subnet{
								ID:   "subnetId-b-public",
								CIDR: "10.0.5.0/24",
								Zone: &aws.Zone{
									Name: "b",
								},
								Public: true,
							},
							"subnetId-c-public": aws.Subnet{
								ID:   "subnetId-c-public",
								CIDR: "10.0.6.0/24",
								Zone: &aws.Zone{
									Name: "c",
								},
								Public: true,
							},
						},
						edgeSubnets: aws.Subnets{
							"subnetId-edge-a-private": aws.Subnet{
								ID:   "subnetId-edge-a-private",
								CIDR: "10.0.7.0/24",
								Zone: &aws.Zone{
									Name: "edge-a",
								},
								Public: false,
							},
							"subnetId-edge-a-public": aws.Subnet{
								ID:   "subnetId-edge-a-public",
								CIDR: "10.0.8.0/24",
								Zone: &aws.Zone{
									Name: "edge-a",
								},
								Public: true,
							},
						},
					},
				},
			},
			want: capa.NetworkSpec{
				VPC: capa.VPCSpec{ID: "vpc-id"},
				Subnets: []capa.SubnetSpec{
					{
						ID:               "subnetId-a-private",
						AvailabilityZone: "a",
						IsPublic:         false,
						CidrBlock:        "10.0.1.0/24",
					}, {
						ID:               "subnetId-a-public",
						AvailabilityZone: "a",
						IsPublic:         true,
						CidrBlock:        "10.0.4.0/24",
					}, {
						ID:               "subnetId-b-private",
						AvailabilityZone: "b",
						IsPublic:         false,
						CidrBlock:        "10.0.2.0/24",
					}, {
						ID:               "subnetId-b-public",
						AvailabilityZone: "b",
						IsPublic:         true,
						CidrBlock:        "10.0.5.0/24",
					}, {
						ID:               "subnetId-c-private",
						AvailabilityZone: "c",
						IsPublic:         false,
						CidrBlock:        "10.0.3.0/24",
					}, {
						ID:               "subnetId-c-public",
						AvailabilityZone: "c",
						IsPublic:         true,
						CidrBlock:        "10.0.6.0/24",
					}, {
						ID:               "subnetId-edge-a-private",
						AvailabilityZone: "edge-a",
						IsPublic:         false,
						CidrBlock:        "10.0.7.0/24",
					}, {
						ID:               "subnetId-edge-a-public",
						AvailabilityZone: "edge-a",
						IsPublic:         true,
						CidrBlock:        "10.0.8.0/24",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setSubnetsBYOVPC(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("setSubnetsBYOVPC() #1 error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if len(tt.args.in.Cluster.Spec.NetworkSpec.Subnets) == 0 {
				if !tt.wantErr {
					t.Errorf("setSubnetsBYOVPC() #2 error: %v, wantErr: %v", err, tt.wantErr)
				}
				return
			}
			tt.args.in.Cluster.Spec.NetworkSpec.Subnets = tSortCapaSubnetsByID(tt.args.in.Cluster.Spec.NetworkSpec.Subnets)
			if got := tt.args.in.Cluster.Spec.NetworkSpec; !cmp.Equal(got, tt.want) {
				t.Errorf("setSubnetsBYOVPC() NetworkSpec.Subnets diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_ZonesCAPI_SetAvailabilityZones(t *testing.T) {
	type fields struct {
		ControlPlaneZones sets.Set[string]
		ComputeZones      sets.Set[string]
	}
	type args struct {
		pool  string
		zones []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ZonesCAPI
	}{
		{
			name: "empty",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:  types.MachinePoolControlPlaneRoleName,
				zones: []string{},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
		},
		{
			name: "set zones for control plane pool",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:  types.MachinePoolControlPlaneRoleName,
				zones: []string{"a", "b"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("a", "b"),
				ComputeZones:      sets.Set[string]{},
			},
		},
		{
			name: "set zones for compute pool",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:  types.MachinePoolComputeRoleName,
				zones: []string{"b", "c"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("b", "c"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo := &ZonesCAPI{
				ControlPlaneZones: tt.fields.ControlPlaneZones,
				ComputeZones:      tt.fields.ComputeZones,
			}
			zo.SetAvailabilityZones(tt.args.pool, tt.args.zones)
			if tt.want != nil {
				assert.EqualValuesf(t, tt.want, zo, "%v failed", tt.name)
			}
		})
	}
}

func Test_ZonesCAPI_SetDefaultConfigZones(t *testing.T) {
	type fields struct {
		AvailabilityZones sets.Set[string]
		ControlPlaneZones sets.Set[string]
		ComputeZones      sets.Set[string]
	}
	type args struct {
		pool      string
		defConfig []string
		defRegion []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ZonesCAPI
	}{
		{
			name: "empty",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{},
				defRegion: []string{},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
		},
		{
			name: "platform defaults when control plane pool exists",
			fields: fields{
				ControlPlaneZones: sets.New("a"),
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{"d"},
				defRegion: []string{"f"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("a"),
				ComputeZones:      sets.Set[string]{},
			},
		},
		{
			name: "platform defaults when control plane pool not exists",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{"d"},
				defRegion: []string{"f"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("d"),
				ComputeZones:      sets.Set[string]{},
			},
		},
		{
			name: "region defaults when control plane pool exists",
			fields: fields{
				ControlPlaneZones: sets.New("a"),
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{},
				defRegion: []string{"f"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("a"),
				ComputeZones:      sets.Set[string]{},
			},
		},
		{
			name: "region defaults when control plane pool not exists",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{},
				defRegion: []string{"f"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.New("f"),
				ComputeZones:      sets.Set[string]{},
			},
		},
		{
			name: "platform defaults when compute pool exists",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("b"),
			},
			args: args{
				pool:      types.MachinePoolComputeRoleName,
				defConfig: []string{"d"},
				defRegion: []string{"f"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("b"),
			},
		},
		{
			name: "platform defaults when compute pool not exists",
			fields: fields{
				AvailabilityZones: sets.Set[string]{},
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolComputeRoleName,
				defConfig: []string{"d"},
				defRegion: []string{"f"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("d"),
			},
		},
		{
			name: "region defaults when compute pool exists",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("b"),
			},
			args: args{
				pool:      types.MachinePoolComputeRoleName,
				defConfig: []string{},
				defRegion: []string{"f"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("b"),
			},
		},
		{
			name: "region defaults when compute pool not exists",
			fields: fields{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolComputeRoleName,
				defConfig: []string{},
				defRegion: []string{"f"},
			},
			want: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("f"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo := &ZonesCAPI{
				ControlPlaneZones: tt.fields.ControlPlaneZones,
				ComputeZones:      tt.fields.ComputeZones,
			}
			zo.SetDefaultConfigZones(tt.args.pool, tt.args.defConfig, tt.args.defRegion)
			if tt.want != nil {
				assert.EqualValuesf(t, tt.want, zo, "%v failed", tt.name)
			}
		})
	}
}

func Test_ZonesCAPI_GetAvailabilityZones(t *testing.T) {
	tests := []struct {
		name  string
		zones *ZonesCAPI
		want  []string
	}{
		{
			name:  "empty",
			zones: &ZonesCAPI{},
			want:  []string{},
		},
		{
			name: "empty az",
			zones: &ZonesCAPI{
				EdgeZones: sets.New("edge-x", "edge-y"),
			},
			want: []string{},
		},
		{
			name: "sorted",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.New("a", "b"),
				ComputeZones:      sets.New("b", "c"),
				EdgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "unsorted",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.New("x", "a"),
				ComputeZones:      sets.New("b", "a"),
				EdgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"a", "b", "x"},
		},
		{
			name: "control planes only",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.New("x", "a"),
				ComputeZones:      sets.Set[string]{},
				EdgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"a", "x"},
		},
		{
			name: "compute only",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("x", "a"),
				EdgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"a", "x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo := tt.zones
			if got := zo.GetAvailabilityZones(); !cmp.Equal(got, tt.want) {
				t.Errorf("ZonesCAPI.GetAvailabilityZones() unexpected results:\n%v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func Test_ZonesCAPI_EdgeZones(t *testing.T) {
	tests := []struct {
		name  string
		zones *ZonesCAPI
		want  []string
	}{
		{
			name:  "empty",
			zones: &ZonesCAPI{},
			want:  []string{},
		},
		{
			name: "empty edge",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.New("x", "y"),
			},
			want: []string{},
		},
		{
			name: "empty only",
			zones: &ZonesCAPI{
				EdgeZones: sets.New("edge-x"),
			},
			want: []string{"edge-x"},
		},
		{
			name: "sorted",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.New("a", "b"),
				ComputeZones:      sets.New("b", "c"),
				EdgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"edge-x", "edge-y"},
		},
		{
			name: "unsorted",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.New("x", "a"),
				ComputeZones:      sets.New("b", "a"),
				EdgeZones:         sets.New("edge-y", "edge-a"),
			},
			want: []string{"edge-a", "edge-y"},
		},
		{
			name: "control planes only",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.New("x", "a"),
				ComputeZones:      sets.Set[string]{},
				EdgeZones:         sets.New("edge-a", "edge-y"),
			},
			want: []string{"edge-a", "edge-y"},
		},
		{
			name: "compute only",
			zones: &ZonesCAPI{
				ControlPlaneZones: sets.Set[string]{},
				ComputeZones:      sets.New("x", "a"),
				EdgeZones:         sets.New("edge-a", "edge-y"),
			},
			want: []string{"edge-a", "edge-y"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo := tt.zones
			if got := zo.GetEdgeZones(); !cmp.Equal(got, tt.want) {
				t.Errorf("ZonesCAPI.EdgeZones() unexpected results:\n %v", cmp.Diff(got, tt.want))
			}
		})
	}
}
