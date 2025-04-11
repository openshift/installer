// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

var stubMachineInputManagedVpc = &MachineInput{
	Role: "master",
	Pool: &types.MachinePool{
		Name:     "master",
		Replicas: ptr.To(int64(3)),
		Platform: types.MachinePoolPlatform{
			AWS: &awstypes.MachinePool{
				Zones: []string{"A", "B", "C"},
			},
		},
	},
	Subnets:  make(map[string]string, 0),
	Tags:     capa.Tags{},
	PublicIP: false,
	Ignition: &capa.Ignition{
		StorageType: capa.IgnitionStorageTypeOptionUnencryptedUserData,
	},
}

func stubDeepCopyMachineInput(in *MachineInput) *MachineInput {
	out := &MachineInput{
		Role:     in.Role,
		PublicIP: in.PublicIP,
	}
	if in.Pool != nil {
		out.Pool = &types.MachinePool{}
		*out.Pool = *in.Pool
	}
	if len(in.Subnets) > 0 {
		out.Subnets = make(map[string]string, len(in.Subnets))
		for k, v := range in.Subnets {
			out.Subnets[k] = v
		}
	}
	if len(in.Tags) > 0 {
		out.Tags = in.Tags.DeepCopy()
	}
	if in.Ignition != nil {
		out.Ignition = in.Ignition.DeepCopy()
	}
	return out
}

func stubGetMachineManagedVpc() *MachineInput {
	return stubDeepCopyMachineInput(stubMachineInputManagedVpc)
}

func TestGenerateMachines(t *testing.T) {
	stubClusterID := "vpc-zr2-m2"
	tests := []struct {
		name           string
		clusterID      string
		input          *MachineInput
		want           []*asset.RuntimeFile
		wantInfraFiles []*asset.RuntimeFile
		wantErr        string
	}{
		{
			name:      "topology ha, managed vpc, default zones region, 2 zones A and B, 3 machines should be in A, B and A private subnet",
			clusterID: stubClusterID,
			input: func() *MachineInput {
				in := stubGetMachineManagedVpc()
				in.Pool.Platform.AWS.Zones = []string{"A", "B"}
				return in
			}(),
			// generate 3 AWSMachine manifests for control plane nodes in two zones
			wantInfraFiles: func() []*asset.RuntimeFile {
				machineZoneMap := map[int]string{0: "A", 1: "B", 2: "A"}
				infraMachineFiles := []*asset.RuntimeFile{}
				for mid := 0; mid < 3; mid++ {
					machineName := fmt.Sprintf("%s-%s-%d", stubClusterID, "master", mid)
					machineZone := machineZoneMap[mid]
					machine := &capa.AWSMachine{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
							Kind:       "AWSMachine",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:   machineName,
							Labels: map[string]string{"cluster.x-k8s.io/control-plane": ""},
						},
						Spec: capa.AWSMachineSpec{
							InstanceMetadataOptions: &capa.InstanceMetadataOptions{
								HTTPEndpoint: capa.InstanceMetadataEndpointStateEnabled,
								HTTPTokens:   capa.HTTPTokensStateOptional,
							},
							AMI: capa.AMIReference{
								ID: ptr.To(""),
							},
							IAMInstanceProfile: fmt.Sprintf("%s-%s-profile", stubClusterID, "master"),
							PublicIP:           ptr.To(false),
							Subnet: &capa.AWSResourceReference{
								Filters: []capa.Filter{{Name: "tag:Name", Values: []string{
									fmt.Sprintf("%s-subnet-private-%s", stubClusterID, machineZone),
								}}},
							},
							SSHKeyName: ptr.To(""),
							RootVolume: &capa.Volume{
								Encrypted: ptr.To(true),
							},
							UncompressedUserData: ptr.To(true),
							Ignition: &capa.Ignition{
								StorageType: capa.IgnitionStorageTypeOptionUnencryptedUserData,
							},
						},
					}
					infraMachineFiles = append(infraMachineFiles, &asset.RuntimeFile{
						File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", machineName)},
						Object: machine,
					})
				}
				return infraMachineFiles
			}(),
		},
		{
			name:      "topology ha, byo vpc, two zones subnets A and B, 3 machines should be in A, B and A private subnets",
			clusterID: stubClusterID,
			input: func() *MachineInput {
				in := stubGetMachineManagedVpc()
				in.Pool.Platform.AWS.Zones = []string{"A", "B"}
				in.Subnets = map[string]string{"A": "subnet-id-A", "B": "subnet-id-B"}
				return in
			}(),
			// generate 3 AWSMachine manifests for control plane nodes in two subnets/zones
			wantInfraFiles: func() []*asset.RuntimeFile {
				machineZoneMap := map[int]string{0: "subnet-id-A", 1: "subnet-id-B", 2: "subnet-id-A"}
				infraMachineFiles := []*asset.RuntimeFile{}
				for mid := 0; mid < 3; mid++ {
					machineName := fmt.Sprintf("%s-%s-%d", stubClusterID, "master", mid)
					machineSubnet := machineZoneMap[mid]
					machine := &capa.AWSMachine{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
							Kind:       "AWSMachine",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:   machineName,
							Labels: map[string]string{"cluster.x-k8s.io/control-plane": ""},
						},
						Spec: capa.AWSMachineSpec{
							InstanceMetadataOptions: &capa.InstanceMetadataOptions{
								HTTPEndpoint: capa.InstanceMetadataEndpointStateEnabled,
								HTTPTokens:   capa.HTTPTokensStateOptional,
							},
							AMI: capa.AMIReference{
								ID: ptr.To(""),
							},
							IAMInstanceProfile: fmt.Sprintf("%s-%s-profile", stubClusterID, "master"),
							PublicIP:           ptr.To(false),
							Subnet: &capa.AWSResourceReference{
								ID: ptr.To(machineSubnet),
							},
							SSHKeyName: ptr.To(""),
							RootVolume: &capa.Volume{
								Encrypted: ptr.To(true),
							},
							UncompressedUserData: ptr.To(true),
							Ignition: &capa.Ignition{
								StorageType: capa.IgnitionStorageTypeOptionUnencryptedUserData,
							},
						},
					}
					infraMachineFiles = append(infraMachineFiles, &asset.RuntimeFile{
						File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", machineName)},
						Object: machine,
					})
				}
				return infraMachineFiles
			}(),
		},
		// Error's scenarios
		{
			name:      "error topology ha, byo vpc, no subnet for zones",
			clusterID: stubClusterID,
			input: func() *MachineInput {
				in := stubGetMachineManagedVpc()
				in.Pool.Platform.AWS.Zones = []string{"A", "B"}
				in.Subnets = map[string]string{"C": "subnet-id-C", "D": "subnet-id-D"}
				return in
			}(),
			wantErr: `no subnet for zone A`,
		},
		{
			name:      "error topology ha, managed vpc, empty subnet zone",
			clusterID: stubClusterID,
			input: func() *MachineInput {
				in := stubGetMachineManagedVpc()
				in.Pool.Platform.AWS.Zones = []string{"A", "B"}
				in.Subnets = map[string]string{"A": "subnet-id-A", "B": ""}
				return in
			}(),
			wantErr: `invalid subnet ID for zone B`,
		},
		// TODO: add more use cases.
		// {
		// 	name: "managed vpc, default zones region, 5 zones A to E, 3 machines should be in A, B and C private subnet",
		// },
		// {
		// 	name: "managed vpc, default zones region, 5 zones A to E, 3 machines should be in A, B and C private subnet",
		// },
		// {
		// 	name: "byo vpc, 2 zones subnets A and B, 3 machines' subnet should be in A, B and A",
		// },
		// {
		// 	name: "managed vpc, default zones region, 5 zones A to E, bootstrap should be in zone A public subnet",
		// },
		// {
		// 	name: "topology ha, managed vpc, two zones subnets A and B, bootstrap node using public subnet A",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := GenerateMachines(tt.clusterID, tt.input)
			if err != nil {
				if len(tt.wantErr) > 0 {
					if got := err.Error(); !cmp.Equal(got, tt.wantErr) {
						t.Errorf("GenerateMachines() unexpected error message: %v", cmp.Diff(got, tt.wantErr))
					}
					return
				}
				t.Errorf("GenerateMachines() unexpected error: %v", err)
				return
			}
			// TODO: support the CAPA v1beta1.Machine manifest check.
			// Support only comparing manifest file for CAPA v1beta2.AWSMachine.
			if len(tt.wantInfraFiles) > 0 {
				got := []*asset.RuntimeFile{}
				for _, file := range files {
					if !strings.HasPrefix(file.Filename, "10_inframachine") {
						continue
					}
					got = append(got, file)
				}
				if !cmp.Equal(got, tt.wantInfraFiles) {
					t.Errorf("GenerateMachines() Got unexpected results:\n%v", cmp.Diff(got, tt.wantInfraFiles))
				}
			}
		})
	}
}
