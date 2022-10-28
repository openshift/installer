package powervs_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/asset/installconfig/powervs/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type editFunctions []func(ic *types.InstallConfig)

var (
	validRegion                  = "lon"
	validCIDR                    = "192.168.0.0/24"
	validCISInstanceCRN          = "crn:v1:bluemix:public:internet-svcs:global:a/valid-account-id:valid-instance-id::"
	validClusterName             = "valid-cluster-name"
	validDNSZoneID               = "valid-zone-id"
	validBaseDomain              = "valid.base.domain"
	validPowerVSResourceGroup    = "valid-resource-group"
	validPublicSubnetUSSouth1ID  = "public-subnet-us-south-1-id"
	validPublicSubnetUSSouth2ID  = "public-subnet-us-south-2-id"
	validPrivateSubnetUSSouth1ID = "private-subnet-us-south-1-id"
	validPrivateSubnetUSSouth2ID = "private-subnet-us-south-2-id"
	validServiceInstanceID       = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	validSubnets                 = []string{
		validPublicSubnetUSSouth1ID,
		validPublicSubnetUSSouth2ID,
		validPrivateSubnetUSSouth1ID,
		validPrivateSubnetUSSouth2ID,
	}
	validUserID = "valid-user@example.com"
	validZone   = "lon04"

	existingDNSRecordsResponse = []powervs.DNSRecordResponse{
		{
			Name: "valid-dns-record-name-1",
			Type: "valid-dns-record-type",
		},
		{
			Name: "valid-dns-record-name-2",
			Type: "valid-dns-record-type",
		},
	}
	noDNSRecordsResponse = []powervs.DNSRecordResponse{}
	invalidArchitecture  = func(ic *types.InstallConfig) { ic.ControlPlane.Architecture = "ppc64" }

	validVPCRegion    = "eu-gb"
	validRG           = "valid-resource-group"
	anotherValidRG    = "another-valid-resource-group"
	validVPCZone      = func(ic *types.InstallConfig) { ic.Platform.PowerVS.VPCZone = "eu-gb-1" }
	invalidVPCZone    = "bogus-vpc-zone"
	validVPCID        = "valid-id"
	anotherValidVPCID = "another-valid-id"
	validVPC          = "valid-vpc"
	anotherValidVPC   = "another-valid-vpc"
	invalidVPC        = "bogus-vpc"
	validVPCs         = []vpcv1.VPC{
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
	validVPCSubnet   = "valid-vpc-subnet"
	invalidVPCSubnet = "invalid-vpc-subnet"
	wrongVPCSubnet   = "wrong-vpc-subnet"
	validSubnet      = &vpcv1.Subnet{
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
	wrongSubnet = &vpcv1.Subnet{
		Name: &validRG,
		VPC: &vpcv1.VPCReference{
			Name: &anotherValidVPC,
			ID:   &anotherValidVPCID,
		},
		ResourceGroup: &vpcv1.ResourceGroupReference{
			Name: &validRG,
			ID:   &validRG,
		},
	}
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
			PowerVS: validMinimalPlatform(),
		},
		ControlPlane: &types.MachinePool{
			Architecture: "ppc64le",
		},
		Compute: []types.MachinePool{{
			Architecture: "ppc64le",
		}},
	}
}

func validMinimalPlatform() *powervstypes.Platform {
	return &powervstypes.Platform{
		PowerVSResourceGroup: validPowerVSResourceGroup,
		Region:               validRegion,
		ServiceInstanceID:    validServiceInstanceID,
		UserID:               validUserID,
		Zone:                 validZone,
	}
}

func validMachinePool() *powervstypes.MachinePool {
	return &powervstypes.MachinePool{}
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
			name:     "invalid architecture",
			edits:    editFunctions{invalidArchitecture},
			errorMsg: `^controlPlane.architecture\: Unsupported value\: \"ppc64\"\: supported values: \"ppc64le\"`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := powervs.Validate(editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestValidatePreExistingPublicDNS(t *testing.T) {
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
			errorMsg: `^\[baseDomain\: Duplicate value\: \"record api\.valid-cluster-name\.valid\.base\.domain already exists in CIS zone \(valid-zone-id\) and might be in use by another cluster, please remove it to continue\", baseDomain\: Duplicate value\: \"record api-int\.valid-cluster-name\.valid\.base\.domain already exists in CIS zone \(valid-zone-id\) and might be in use by another cluster, please remove it to continue\"\]$`,
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
	setMockEnvVars()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	powervsClient := mock.NewMockAPI(mockCtrl)
	metadata := mock.NewMockMetadataAPI(mockCtrl)

	dnsRecordNames := [...]string{fmt.Sprintf("api.%s.%s", validClusterName, validBaseDomain), fmt.Sprintf("api-int.%s.%s", validClusterName, validBaseDomain)}

	// Mock common to all tests
	metadata.EXPECT().CISInstanceCRN(gomock.Any()).Return(validCISInstanceCRN, nil).AnyTimes()

	// Mocks: no pre-existing DNS records
	powervsClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain, types.ExternalPublishingStrategy).Return(validDNSZoneID, nil)
	for _, dnsRecordName := range dnsRecordNames {
		powervsClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName, types.ExternalPublishingStrategy).Return(noDNSRecordsResponse, nil)
	}

	// Mocks: pre-existing DNS records
	powervsClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain, types.ExternalPublishingStrategy).Return(validDNSZoneID, nil)
	for _, dnsRecordName := range dnsRecordNames {
		powervsClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName, types.ExternalPublishingStrategy).Return(existingDNSRecordsResponse, nil)
	}

	// Mocks: cannot get zone ID
	powervsClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain, types.ExternalPublishingStrategy).Return("", fmt.Errorf(""))

	// Mocks: cannot get DNS records
	powervsClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain, types.ExternalPublishingStrategy).Return(validDNSZoneID, nil)
	for _, dnsRecordName := range dnsRecordNames {
		powervsClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName, types.ExternalPublishingStrategy).Return(nil, fmt.Errorf(""))
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			aggregatedErrors := powervs.ValidatePreExistingDNS(powervsClient, validInstallConfig(), metadata)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestValidateCustomVPCSettings(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name: "invalid VPC name supplied, without zone",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = invalidVPC
				},
			},
			errorMsg: "VPC.vpcName: Internal error: vpcName",
		},
		{
			name: "valid VPC name supplied, without zone",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = validVPC
				},
			},
			errorMsg: "",
		},
		{
			name: "valid VPC name, but invalid subnet, without zone",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = validVPC
					ic.Platform.PowerVS.VPCSubnets = []string{invalidVPCSubnet}
				},
			},
			errorMsg: "VPC.vpcSubnets: Internal error",
		},
		{
			name: "valid VPC name, valid subnet, without zone",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = validVPC
					ic.Platform.PowerVS.VPCSubnets = []string{validVPCSubnet}
				},
			},
			errorMsg: "",
		},
		{
			name: "invalid VPC zone supplied",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCZone = invalidVPCZone
				},
			},
			errorMsg: fmt.Sprintf(`VPC.vpcZone: Invalid value: "null": %s`, invalidVPCZone),
		},
		{
			name: "valid VPC zone supplied",
			edits: editFunctions{
				validVPCZone,
			},
			errorMsg: "",
		},
		{
			name: "VPC not found for the specified zone",
			edits: editFunctions{
				validVPCZone,
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = invalidVPC
				},
			},
			errorMsg: fmt.Sprintf(`VPC.vpcName: Not found: "%s"`, invalidVPC),
		},
		{
			name: "VPC found for the specified zone",
			edits: editFunctions{
				validVPCZone,
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = validVPC
				},
			},
			errorMsg: "",
		},
		{
			name: "VPC found, but not subnet, for the specified zone",
			edits: editFunctions{
				validVPCZone,
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = validVPC
					ic.Platform.PowerVS.VPCSubnets = []string{invalidVPCSubnet}
				},
			},
			errorMsg: "VPC.vpcSubnets: Internal error",
		},
		{
			name: "VPC found, subnet found as well, but not attached to the VPC",
			edits: editFunctions{
				validVPCZone,
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = validVPC
					ic.Platform.PowerVS.VPCSubnets = []string{wrongVPCSubnet}
				},
			},
			errorMsg: `VPC.vpcSubnets: Invalid value: "null": not attached to VPC`,
		},
		{
			name: "VPC and subnet both found for the specified zone",
			edits: editFunctions{
				validVPCZone,
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = validVPC
					ic.Platform.PowerVS.VPCSubnets = []string{validVPCSubnet}
				},
			},
			errorMsg: "",
		},
	}
	setMockEnvVars()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	powervsClient := mock.NewMockAPI(mockCtrl)

	// Mocks: invalid VPC name results in error
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(nil, fmt.Errorf("vpcName: not found"))

	// Mocks: valid VPC name only, no issues
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)

	// Mocks: valid VPC name, invalid VPC subnet
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)
	powervsClient.EXPECT().GetSubnetByName(gomock.Any(), invalidVPCSubnet, validVPCRegion).Return(nil, fmt.Errorf(""))

	// Mocks: valid VPC name, valid VPC subnet, all good
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)
	powervsClient.EXPECT().GetSubnetByName(gomock.Any(), validVPCSubnet, validVPCRegion).Return(validSubnet, nil)

	// Mocks: invalid Zone name results in error
	// nothing to mock

	// Mocks: valid Zone name only, no issues
	// nothing to mock

	// Mocks: valid Zone name, but no matching VPC name
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)

	// Mocks: valid Zone name, valid VPC found
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)

	// Mocks: valid Zone, valid VPC, but invalid subnet
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)
	powervsClient.EXPECT().GetSubnetByName(gomock.Any(), invalidVPCSubnet, validVPCRegion).Return(nil, fmt.Errorf(""))

	// Mocks: valid Zone, valid VPC, but wrong subnet (present, but not attached)
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)
	powervsClient.EXPECT().GetSubnetByName(gomock.Any(), wrongVPCSubnet, validVPCRegion).Return(wrongSubnet, nil)

	// Mocks: valid Zone, valid VPC, valid subnet, all good
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)
	powervsClient.EXPECT().GetSubnetByName(gomock.Any(), validVPCSubnet, validVPCRegion).Return(validSubnet, nil)

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := powervs.ValidateCustomVPCSetup(powervsClient, editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func setMockEnvVars() {
	os.Setenv("POWERVS_AUTH_FILEPATH", "./tmp/powervs/config.json")
	os.Setenv("IBMID", "foo")
	os.Setenv("IC_API_KEY", "foo")
	os.Setenv("IBMCLOUD_REGION", "foo")
	os.Setenv("IBMCLOUD_ZONE", "foo")
}
