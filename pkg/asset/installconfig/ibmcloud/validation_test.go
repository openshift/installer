package ibmcloud_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type editFunctions []func(ic *types.InstallConfig)

var (
	validRegion                  = "us-south"
	validCIDR                    = "10.0.0.0/16"
	validCISInstanceCRN          = "crn:v1:bluemix:public:internet-svcs:global:a/valid-account-id:valid-instance-id::"
	validClusterName             = "valid-cluster-name"
	validDNSZoneID               = "valid-zone-id"
	validBaseDomain              = "valid.base.domain"
	validPublicSubnetUSSouth1ID  = "public-subnet-us-south-1-id"
	validPublicSubnetUSSouth2ID  = "public-subnet-us-south-2-id"
	validPrivateSubnetUSSouth1ID = "private-subnet-us-south-1-id"
	validPrivateSubnetUSSouth2ID = "private-subnet-us-south-2-id"
	validSubnets                 = []string{
		validPublicSubnetUSSouth1ID,
		validPublicSubnetUSSouth2ID,
		validPrivateSubnetUSSouth1ID,
		validPrivateSubnetUSSouth2ID,
	}
	validSubnetName   = "valid-subnet"
	validVPCID        = "valid-id"
	validVPC          = "valid-vpc"
	validRG           = "valid-resource-group"
	validZoneUSSouth1 = "us-south-1"
	wrongRG           = "wrong-resource-group"
	wrongVPCID        = "wrong-id"
	wrongVPC          = "wrong-vpc"
	anotherValidVPCID = "another-valid-id"
	anotherValidVPC   = "another-valid-vpc"
	anotherValidRG    = "another-valid-resource-group"

	validResourceGroups = []resourcemanagerv2.ResourceGroup{
		{
			Name: &validRG,
			ID:   &validRG,
		},
		{
			Name: &anotherValidRG,
			ID:   &anotherValidRG,
		},
	}
	validVPCs = []vpcv1.VPC{
		{
			Name: &validVPC,
			ID:   &validVPCID,
			ResourceGroup: &vpcv1.ResourceGroupReference{
				Name: &validRG,
				ID:   &validRG,
			},
		},
		{
			Name: &anotherValidVPC,
			ID:   &anotherValidVPCID,
			ResourceGroup: &vpcv1.ResourceGroupReference{
				Name: &anotherValidRG,
				ID:   &anotherValidRG,
			},
		},
	}
	invalidVPC = []vpcv1.VPC{
		{
			Name: &wrongVPC,
			ID:   &wrongVPCID,
			ResourceGroup: &vpcv1.ResourceGroupReference{
				Name: &validRG,
				ID:   &validRG,
			},
		},
	}
	validVPCInvalidRG = []vpcv1.VPC{
		{
			Name: &validVPC,
			ID:   &validVPCID,
			ResourceGroup: &vpcv1.ResourceGroupReference{
				Name: &wrongRG,
				ID:   &wrongRG,
			},
		},
	}
	validSubnet = &vpcv1.Subnet{
		Name: &validRG,
		VPC: &vpcv1.VPCReference{
			Name: &validVPC,
			ID:   &validVPCID,
		},
		ResourceGroup: &vpcv1.ResourceGroupReference{
			Name: &validRG,
			ID:   &validRG,
		},
	}

	validInstanceProfies = []vpcv1.InstanceProfile{{Name: &[]string{"type-a"}[0]}, {Name: &[]string{"type-b"}[0]}}

	machinePoolInvalidType = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.IBMCloud = &ibmcloudtypes.MachinePool{
			InstanceType: "invalid-type",
		}
	}

	existingDNSRecordsResponse = []dnsrecordsv1.DnsrecordDetails{
		{
			ID: core.StringPtr("valid-dns-record-1"),
		},
		{
			ID: core.StringPtr("valid-dns-record-2"),
		},
	}
	noDNSRecordsResponse = []dnsrecordsv1.DnsrecordDetails{}
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: validClusterName,
		},
		BaseDomain: validBaseDomain,
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Publish: types.ExternalPublishingStrategy,
		Platform: types.Platform{
			IBMCloud: validMinimalPlatform(),
		},
		ControlPlane: &types.MachinePool{
			Platform: types.MachinePoolPlatform{
				IBMCloud: validMachinePool(),
			},
		},
		Compute: []types.MachinePool{{
			Platform: types.MachinePoolPlatform{
				IBMCloud: validMachinePool(),
			},
		}},
	}
}

func validMinimalPlatform() *ibmcloudtypes.Platform {
	return &ibmcloudtypes.Platform{
		Region: validRegion,
	}
}

func validMachinePool() *ibmcloudtypes.MachinePool {
	return &ibmcloudtypes.MachinePool{}
}

func validResourceGroupName(ic *types.InstallConfig) {
	ic.Platform.IBMCloud.ResourceGroupName = "valid-resource-group"
}

func validVPCName(ic *types.InstallConfig) {
	ic.Platform.IBMCloud.VPCName = "valid-vpc"
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name:     "valid install config",
			edits:    editFunctions{},
			errorMsg: "",
		},
		{
			name: "VPC with no ResourceGroup supplied",
			edits: editFunctions{
				validVPCName,
			},
			errorMsg: `resourceGroupName: Not found: ""$`,
		},
		{
			name: "VPC not found",
			edits: editFunctions{
				validResourceGroupName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.VPCName = "missing-vpc"
				},
			},
			errorMsg: `vpcName: Not found: "missing-vpc"$`,
		},
		{
			name: "VPC not in ResourceGroup",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ResourceGroupName = "wrong-resource-group"
				},
				validVPCName,
			},
			errorMsg: `platform.ibmcloud.vpcName: Invalid value: "valid-vpc": vpc is not in provided ResourceGroup: wrong-resource-group`,
		},
		{
			name: "VPC with no control plane subnets",
			edits: editFunctions{
				validResourceGroupName,
				validVPCName,
			},
			errorMsg: `platform.ibmcloud.controlPlaneSubnets: Invalid value: \[\]string\(nil\): controlPlaneSubnets cannot be empty when providing a vpcName: valid-vpc`,
		},
		{
			name: "control plane subnet not found",
			edits: editFunctions{
				validResourceGroupName,
				validVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ControlPlaneSubnets = []string{"missing-cp-subnet"}
				},
			},
			errorMsg: `platform.ibmcloud.controlPlaneSubnets: Not found: "missing-cp-subnet"`,
		},
		{
			name: "control plane subnet IBM error",
			edits: editFunctions{
				validResourceGroupName,
				validVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ControlPlaneSubnets = []string{"ibm-error-cp-subnet"}
				},
			},
			errorMsg: `platform.ibmcloud.controlPlaneSubnets: Internal error: ibmcloud error`,
		},
		{
			name: "control plane subnet invalid VPC",
			edits: editFunctions{
				validResourceGroupName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.VPCName = "wrong-vpc"
				},
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ControlPlaneSubnets = []string{"valid-subnet"}
				},
			},
			errorMsg: `platform.ibmcloud.controlPlaneSubnets: Invalid value: "valid-subnet": controlPlaneSubnets contains subnet: valid-subnet, not found in expected vpcID: wrong-id`,
		},
		{
			name: "control plane subnet invalid ResourceGroup",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ResourceGroupName = "wrong-resource-group"
				},
				validVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ControlPlaneSubnets = []string{"valid-subnet"}
				},
			},
			errorMsg: `platform.ibmcloud.controlPlaneSubnets: Invalid value: "valid-subnet": controlPlaneSubnets contains subnet: valid-subnet, not found in expected resourceGroupName: wrong-resource-group`,
		},
		{
			name: "VPC with no compute subnets",
			edits: editFunctions{
				validResourceGroupName,
				validVPCName,
			},
			errorMsg: `platform.ibmcloud.computeSubnets: Invalid value: \[\]string\(nil\): computeSubnets cannot be empty when providing a vpcName: valid-vpc`,
		},
		{
			name: "compute subnet not found",
			edits: editFunctions{
				validResourceGroupName,
				validVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ComputeSubnets = []string{"missing-compute-subnet"}
				},
			},
			errorMsg: `platform.ibmcloud.computeSubnets: Not found: "missing-compute-subnet"`,
		},
		{
			name: "compute subnet IBM error",
			edits: editFunctions{
				validResourceGroupName,
				validVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ComputeSubnets = []string{"ibm-error-compute-subnet"}
				},
			},
			errorMsg: `platform.ibmcloud.computeSubnets: Internal error: ibmcloud error`,
		},
		{
			name: "compute subnet invalid VPC",
			edits: editFunctions{
				validResourceGroupName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.VPCName = "wrong-vpc"
				},
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ComputeSubnets = []string{"valid-subnet"}
				},
			},
			errorMsg: `platform.ibmcloud.computeSubnets: Invalid value: "valid-subnet": computeSubnets contains subnet: valid-subnet, not found in expected vpcID: wrong-id`,
		},
		{
			name: "compute subnet invalid ResourceGroup",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ResourceGroupName = "wrong-resource-group"
				},
				validVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.IBMCloud.ComputeSubnets = []string{"valid-subnet"}
				},
			},
			errorMsg: `platform.ibmcloud.computeSubnets: Invalid value: "valid-subnet": computeSubnets contains subnet: valid-subnet, not found in expected resourceGroupName: wrong-resource-group`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	// Mocks: valid install config and all other tests ('AnyTimes()')
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), validPublicSubnetUSSouth1ID).Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &validZoneUSSouth1}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), validPublicSubnetUSSouth2ID).Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &validZoneUSSouth1}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), validPrivateSubnetUSSouth1ID).Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &validZoneUSSouth1}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), validPrivateSubnetUSSouth2ID).Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &validZoneUSSouth1}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), "subnet-invalid-zone").Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &[]string{"invalid"}[0]}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetVSIProfiles(gomock.Any()).Return(validInstanceProfies, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetVPCZonesForRegion(gomock.Any(), validRegion).Return([]string{"us-south-1", "us-south-2", "us-south-3"}, nil).AnyTimes()

	// Mocks: VPC with no ResourceGroup supplied
	// No mocks required

	// Mocks: VPC not found
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCs, nil)

	// Mocks: VPC not in ResourceGroup
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCs, nil)

	// Mocks: VPC with no control plane subnets
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCs, nil)

	// Mocks: control plane subnet not found
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCs, nil)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), "missing-cp-subnet", validRegion).Return(nil, &ibmcloud.VPCResourceNotFoundError{})

	// Mocks: control plane subnet IBM error
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCs, nil)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), "ibm-error-cp-subnet", validRegion).Return(nil, errors.New("ibmcloud error"))

	// Mocks: control plane subnet invalid VPC
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(invalidVPC, nil)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), "valid-subnet", validRegion).Return(validSubnet, nil)

	// Mocks: control plane subnet invalid ResourceGroup
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCInvalidRG, nil)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), "valid-subnet", validRegion).Return(validSubnet, nil)

	// Mocks: VPC with no compute subnets
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCs, nil)

	// Mocks: compute subnet not found
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCs, nil)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), "missing-compute-subnet", validRegion).Return(nil, &ibmcloud.VPCResourceNotFoundError{})

	// Mocks: compute subnet IBM error
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCs, nil)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), "ibm-error-compute-subnet", validRegion).Return(nil, errors.New("ibmcloud error"))

	// Mocks: compute subnet invalid VPC
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(invalidVPC, nil)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), "valid-subnet", validRegion).Return(validSubnet, nil)

	// Mocks: compute subnet invalid ResourceGroup
	ibmcloudClient.EXPECT().GetResourceGroups(gomock.Any()).Return(validResourceGroups, nil)
	ibmcloudClient.EXPECT().GetVPCs(gomock.Any(), validRegion).Return(validVPCInvalidRG, nil)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), "valid-subnet", validRegion).Return(validSubnet, nil)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := ibmcloud.Validate(ibmcloudClient, editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestValidatePreExitingPublicDNS(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name:     "no pre-existing DNS records",
			errorMsg: "",
		},
		{
			name:     "pre-existing DNS records",
			errorMsg: `^record api\.valid-cluster-name\.valid\.base\.domain already exists in CIS zone \(valid-zone-id\) and might be in use by another cluster, please remove it to continue$`,
		},
		{
			name:     "cannot get zone ID",
			errorMsg: `^baseDomain: Internal error$`,
		},
		{
			name:     "cannot get DNS records",
			errorMsg: `^baseDomain: Internal error$`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	dnsRecordName := fmt.Sprintf("api.%s.%s", validClusterName, validBaseDomain)

	metadata := ibmcloud.NewMetadata(validBaseDomain, "us-south", nil, nil)
	metadata.SetCISInstanceCRN(validCISInstanceCRN)

	// Mocks: no pre-existing DNS records
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain).Return(validDNSZoneID, nil)
	ibmcloudClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName).Return(noDNSRecordsResponse, nil)

	// Mocks: pre-existing DNS records
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain).Return(validDNSZoneID, nil)
	ibmcloudClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName).Return(existingDNSRecordsResponse, nil)

	// Mocks: cannot get zone ID
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain).Return("", fmt.Errorf(""))

	// Mocks: cannot get DNS records
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain).Return(validDNSZoneID, nil)
	ibmcloudClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName).Return(nil, fmt.Errorf(""))

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			validInstallConfig := validInstallConfig()
			aggregatedErrors := ibmcloud.ValidatePreExitingPublicDNS(ibmcloudClient, validInstallConfig, metadata)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}
