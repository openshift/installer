package aws

import (
	"reflect"
	"sort"
	"testing"

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

var (
	stubDefaultCIDR = "10.0.0.0/16"
)

func stubClusterID() *installconfig.ClusterID {
	return &installconfig.ClusterID{
		InfraID: "infra-id",
	}
}

func stubInstallConfig() *installconfig.InstallConfig {
	ic := &installconfig.InstallConfig{}
	return ic
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

func stubInstallCOnfigPoolCompute() []types.MachinePool {
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
	p := stubInstallCOnfigPoolCompute()
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
	subnetIds := []string{}
	subnetsMap := make(map[string]capa.SubnetSpec, len(in))
	for _, subnet := range in {
		subnetsMap[subnet.ID] = subnet
		subnetIds = append(subnetIds, subnet.ID)
	}
	sort.Strings(subnetIds)
	out := capa.Subnets{}
	for _, sid := range subnetIds {
		out = append(out, subnetsMap[sid])
	}
	return out
}

func Test_extractZonesFromInstallConfig(t *testing.T) {
	type args struct {
		in *zonesInput
	}
	tests := []struct {
		name    string
		args    args
		want    *zonesCAPI
		wantErr bool
	}{
		{
			name: "default zones",
			args: args{
				in: &zonesInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = stubInstallConfigType()
						return ic
					}(),
				},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
				edgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "custom zones control plane pool",
			args: args{
				in: &zonesInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: stubInstallConfigPoolControl(),
							Compute:      nil,
						}
						return ic
					}(),
				},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.New("a", "b"),
				computeZones:      sets.Set[string]{},
				edgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "custom zones compute pool",
			args: args{
				in: &zonesInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: nil,
							Compute:      stubInstallCOnfigPoolCompute(),
						}
						return ic
					}(),
				},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("b", "c"),
				edgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "custom zones control plane and compute pools",
			args: args{
				in: &zonesInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: stubInstallConfigPoolControl(),
							Compute:      stubInstallCOnfigPoolCompute(),
						}
						return ic
					}(),
				},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.New("a", "b"),
				computeZones:      sets.New("b", "c"),
				edgeZones:         sets.Set[string]{},
			},
		},
		{
			name: "custom zones control plane, compute and edge pools",
			args: args{
				in: &zonesInput{
					InstallConfig: func() *installconfig.InstallConfig {
						ic := stubInstallConfig()
						ic.Config = &types.InstallConfig{
							ControlPlane: stubInstallConfigPoolControl(),
							Compute:      stubInstallConfigPoolComputeWithEdge(),
						}
						return ic
					}(),
				},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.New("a", "b"),
				computeZones:      sets.New("b", "c"),
				edgeZones:         sets.New("edge-b", "edge-c"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractZonesFromInstallConfig(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractZonesFromInstallConfig() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractZonesFromInstallConfig() err=%v\ngot : %#v,\nwant: %#v\n", err, got, tt.want)
			}
		})
	}
}

func Test_setSubnetsManagedVPC(t *testing.T) {
	type args struct {
		in *zonesInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *capa.NetworkSpec
	}{
		{
			name: "regular zones in the region",
			args: args{
				in: &zonesInput{
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
			want: &capa.NetworkSpec{
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
			name: "regular zones in the region with edge",
			args: args{
				in: &zonesInput{
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
			want: &capa.NetworkSpec{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setSubnetsManagedVPC(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("setSubnetsManagedVPC() #1 error: %+v,\nwantErr %+v\n", err, tt.wantErr)
			}
			var got *capa.NetworkSpec
			if tt.args.in != nil && tt.args.in.Cluster != nil {
				got = &tt.args.in.Cluster.Spec.NetworkSpec
			} else {
				if !tt.wantErr {
					t.Errorf("setSubnetsManagedVPC() #2 error: %v, wantErr: %v", err, tt.wantErr)
				}
				return
			}
			if len(got.Subnets) == 0 {
				if !tt.wantErr {
					t.Errorf("setSubnetsManagedVPC() #2 error: %v, wantErr: %v", err, tt.wantErr)
				}
				return
			}
			got.Subnets = tSortCapaSubnetsByID(got.Subnets)
			if tt.want != nil {
				assert.EqualValuesf(t, tt.want, got, "%v failed.\nWant: %+v\nGot: %+v", tt.name)
			}
		})
	}
}

func Test_setSubnetsBYOVPC(t *testing.T) {
	type args struct {
		in *zonesInput
	}
	tests := []struct {
		name    string
		args    args
		want    *capa.NetworkSpec
		wantErr bool
	}{
		{
			name: "default byo vpc",
			args: args{
				in: &zonesInput{
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
			want: &capa.NetworkSpec{
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
				in: &zonesInput{
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
			want: &capa.NetworkSpec{
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
				in: &zonesInput{
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
			want: &capa.NetworkSpec{
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
			var got *capa.NetworkSpec
			if tt.args.in != nil && tt.args.in.Cluster != nil {
				got = &tt.args.in.Cluster.Spec.NetworkSpec
			} else {
				if !tt.wantErr {
					t.Errorf("setSubnetsBYOVPC() #2 error: %v, wantErr: %v", err, tt.wantErr)
				}
				return
			}
			if len(got.Subnets) == 0 {
				if !tt.wantErr {
					t.Errorf("setSubnetsBYOVPC() #2 error: %v, wantErr: %v", err, tt.wantErr)
				}
				return
			}
			got.Subnets = tSortCapaSubnetsByID(got.Subnets)
			if tt.want != nil {
				assert.EqualValuesf(t, tt.want, got, "%v failed.\nWant: %+v\nGot: %+v", tt.name)
			}
		})
	}
}

func Test_zonesCAPI_SetAvailabilityZones(t *testing.T) {
	type fields struct {
		controlPlaneZones sets.Set[string]
		computeZones      sets.Set[string]
	}
	type args struct {
		pool  string
		zones []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *zonesCAPI
	}{
		{
			name: "empty",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:  types.MachinePoolControlPlaneRoleName,
				zones: []string{},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
		},
		{
			name: "set zones for control plane pool",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:  types.MachinePoolControlPlaneRoleName,
				zones: []string{"a", "b"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.New("a", "b"),
				computeZones:      sets.Set[string]{},
			},
		},
		{
			name: "set zones for compute pool",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:  types.MachinePoolComputeRoleName,
				zones: []string{"b", "c"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("b", "c"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo := &zonesCAPI{
				controlPlaneZones: tt.fields.controlPlaneZones,
				computeZones:      tt.fields.computeZones,
			}
			zo.SetAvailabilityZones(tt.args.pool, tt.args.zones)
			if tt.want != nil {
				assert.EqualValuesf(t, tt.want, zo, "%v failed", tt.name)
			}
		})
	}
}

func Test_zonesCAPI_SetDefaultConfigZones(t *testing.T) {
	type fields struct {
		AvailabilityZones sets.Set[string]
		controlPlaneZones sets.Set[string]
		computeZones      sets.Set[string]
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
		want   *zonesCAPI
	}{
		{
			name: "empty",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{},
				defRegion: []string{},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
		},
		{
			name: "platform defaults when control plane pool exists",
			fields: fields{
				controlPlaneZones: sets.New("a"),
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{"d"},
				defRegion: []string{"f"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.New("a"),
				computeZones:      sets.Set[string]{},
			},
		},
		{
			name: "platform defaults when control plane pool not exists",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{"d"},
				defRegion: []string{"f"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.New("d"),
				computeZones:      sets.Set[string]{},
			},
		},
		{
			name: "region defaults when control plane pool exists",
			fields: fields{
				controlPlaneZones: sets.New("a"),
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{},
				defRegion: []string{"f"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.New("a"),
				computeZones:      sets.Set[string]{},
			},
		},
		{
			name: "region defaults when control plane pool not exists",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolControlPlaneRoleName,
				defConfig: []string{},
				defRegion: []string{"f"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.New("f"),
				computeZones:      sets.Set[string]{},
			},
		},
		{
			name: "platform defaults when compute pool exists",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("b"),
			},
			args: args{
				pool:      types.MachinePoolComputeRoleName,
				defConfig: []string{"d"},
				defRegion: []string{"f"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("b"),
			},
		},
		{
			name: "platform defaults when compute pool not exists",
			fields: fields{
				AvailabilityZones: sets.Set[string]{},
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolComputeRoleName,
				defConfig: []string{"d"},
				defRegion: []string{"f"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("d"),
			},
		},
		{
			name: "region defaults when compute pool exists",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("b"),
			},
			args: args{
				pool:      types.MachinePoolComputeRoleName,
				defConfig: []string{},
				defRegion: []string{"f"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("b"),
			},
		},
		{
			name: "region defaults when compute pool not exists",
			fields: fields{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.Set[string]{},
			},
			args: args{
				pool:      types.MachinePoolComputeRoleName,
				defConfig: []string{},
				defRegion: []string{"f"},
			},
			want: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("f"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo := &zonesCAPI{
				controlPlaneZones: tt.fields.controlPlaneZones,
				computeZones:      tt.fields.computeZones,
			}
			zo.SetDefaultConfigZones(tt.args.pool, tt.args.defConfig, tt.args.defRegion)
			if tt.want != nil {
				assert.EqualValuesf(t, tt.want, zo, "%v failed", tt.name)
			}
		})
	}
}

func Test_zonesCAPI_AvailabilityZones(t *testing.T) {
	tests := []struct {
		name  string
		zones *zonesCAPI
		want  []string
	}{
		{
			name:  "empty",
			zones: &zonesCAPI{},
			want:  []string{},
		},
		{
			name: "empty az",
			zones: &zonesCAPI{
				edgeZones: sets.New("edge-x", "edge-y"),
			},
			want: []string{},
		},
		{
			name: "sorted",
			zones: &zonesCAPI{
				controlPlaneZones: sets.New("a", "b"),
				computeZones:      sets.New("b", "c"),
				edgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "unsorted",
			zones: &zonesCAPI{
				controlPlaneZones: sets.New("x", "a"),
				computeZones:      sets.New("b", "a"),
				edgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"a", "b", "x"},
		},
		{
			name: "control planes only",
			zones: &zonesCAPI{
				controlPlaneZones: sets.New("x", "a"),
				computeZones:      sets.Set[string]{},
				edgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"a", "x"},
		},
		{
			name: "compute only",
			zones: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("x", "a"),
				edgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"a", "x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo := tt.zones
			if got := zo.AvailabilityZones(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("zonesCAPI.AvailabilityZones() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_zonesCAPI_EdgeZones(t *testing.T) {
	tests := []struct {
		name  string
		zones *zonesCAPI
		want  []string
	}{
		{
			name:  "empty",
			zones: &zonesCAPI{},
			want:  []string{},
		},
		{
			name: "empty edge",
			zones: &zonesCAPI{
				controlPlaneZones: sets.New("x", "y"),
			},
			want: []string{},
		},
		{
			name: "empty only",
			zones: &zonesCAPI{
				edgeZones: sets.New("edge-x"),
			},
			want: []string{"edge-x"},
		},
		{
			name: "sorted",
			zones: &zonesCAPI{
				controlPlaneZones: sets.New("a", "b"),
				computeZones:      sets.New("b", "c"),
				edgeZones:         sets.New("edge-x", "edge-y"),
			},
			want: []string{"edge-x", "edge-y"},
		},
		{
			name: "unsorted",
			zones: &zonesCAPI{
				controlPlaneZones: sets.New("x", "a"),
				computeZones:      sets.New("b", "a"),
				edgeZones:         sets.New("edge-y", "edge-a"),
			},
			want: []string{"edge-a", "edge-y"},
		},
		{
			name: "control planes only",
			zones: &zonesCAPI{
				controlPlaneZones: sets.New("x", "a"),
				computeZones:      sets.Set[string]{},
				edgeZones:         sets.New("edge-a", "edge-y"),
			},
			want: []string{"edge-a", "edge-y"},
		},
		{
			name: "compute only",
			zones: &zonesCAPI{
				controlPlaneZones: sets.Set[string]{},
				computeZones:      sets.New("x", "a"),
				edgeZones:         sets.New("edge-a", "edge-y"),
			},
			want: []string{"edge-a", "edge-y"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zo := tt.zones
			if got := zo.EdgeZones(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("zonesCAPI.EdgeZones() = %v, want %v", got, tt.want)
			}
		})
	}
}
