package powervs_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/asset/installconfig/powervs/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

type editFunctions []func(ic *types.InstallConfig)

var (
	validRegion                  = "dal"
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
	validSubnets                 = []string{
		validPublicSubnetUSSouth1ID,
		validPublicSubnetUSSouth2ID,
		validPrivateSubnetUSSouth1ID,
		validPrivateSubnetUSSouth2ID,
	}
	validUserID = "valid-user@example.com"
	validZone   = "dal10"

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
	noDNSRecordsResponse   = []powervs.DNSRecordResponse{}
	invalidArchitecture    = func(ic *types.InstallConfig) { ic.ControlPlane.Architecture = "ppc64" }
	cidrInvalid, _         = ipnet.ParseCIDR("192.168.0.0/16")
	invalidMachinePoolCIDR = func(ic *types.InstallConfig) { ic.Networking.MachineNetwork[0].CIDR = *cidrInvalid }
	cidrValid, _           = ipnet.ParseCIDR("192.168.0.0/24")
	validMachinePoolCIDR   = func(ic *types.InstallConfig) { ic.Networking.MachineNetwork[0].CIDR = *cidrValid }
	validVPCRegion         = "us-south"
	invalidVPCRegion       = "foo-bah"
	setValidVPCRegion      = func(ic *types.InstallConfig) { ic.Platform.PowerVS.VPCRegion = validVPCRegion }
	validRG                = "valid-resource-group"
	anotherValidRG         = "another-valid-resource-group"
	validVPCID             = "valid-id"
	anotherValidVPCID      = "another-valid-id"
	validVPC               = "valid-vpc"
	setValidVPCName        = func(ic *types.InstallConfig) { ic.Platform.PowerVS.VPCName = validVPC }
	anotherValidVPC        = "another-valid-vpc"
	invalidVPC             = "bogus-vpc"
	validVPCs              = []vpcv1.VPC{
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
	regionWithPER    = "dal10"
	regionWithoutPER = "foo99"
	regionPERUnknown = "bah77"
	mapWithPERFalse  = map[string]bool{
		"disaster-recover-site": true,
		"power-edge-router":     false,
		"vpn-connections":       true,
	}
	mapWithPERTrue = map[string]bool{
		"disaster-recover-site": true,
		"power-edge-router":     true,
		"vpn-connections":       true,
	}
	mapPERUnknown = map[string]bool{
		"disaster-recover-site": true,
		"power-vpn-connections": false,
	}
	defaultSysType           = "s922"
	newSysType               = "s1022"
	invalidRegion            = "foo"
	validServiceInstanceGUID = ""
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
		ServiceInstanceGUID:  validServiceInstanceGUID,
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
		{
			name:     "invalid machine pool CIDR",
			edits:    editFunctions{invalidMachinePoolCIDR},
			errorMsg: `Networking.MachineNetwork.CIDR: Invalid value: "192.168.0.0/16": Machine Pool CIDR must be /24.`,
		},
		{
			name:     "valid machine pool CIDR",
			edits:    editFunctions{validMachinePoolCIDR},
			errorMsg: "",
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
			name: "invalid VPC region supplied alone",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCRegion = invalidVPCRegion
				},
			},
			errorMsg: fmt.Sprintf(`VPC.vpcRegion: Not found: "%s"`, invalidVPCRegion),
		},
		{
			name: "valid VPC region supplied alone",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCRegion = validVPCRegion
				},
			},
			errorMsg: "",
		},
		{
			name: "invalid VPC name supplied, without VPC region, not found near PowerVS region",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCName = invalidVPC
				},
			},
			errorMsg: fmt.Sprintf(`VPC.vpcName: Not found: "%s"`, invalidVPC),
		},
		{
			name: "valid VPC name supplied, without VPC region, but found close to PowerVS region",
			edits: editFunctions{
				setValidVPCName,
			},
			errorMsg: "",
		},
		{
			name: "valid VPC name, with invalid VPC region",
			edits: editFunctions{
				setValidVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCRegion = invalidVPCRegion
				},
			},
			errorMsg: "VPC.vpcRegion: Internal error: unknown region",
		},
		{
			name: "valid VPC name, valid VPC region",
			edits: editFunctions{
				setValidVPCName,
				setValidVPCRegion,
			},
			errorMsg: "",
		},
		{
			name: "VPC subnet supplied, without vpcName",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCSubnets = []string{validVPCSubnet}
				},
			},
			errorMsg: `VPC.vpcSubnets: Invalid value: "null": invalid without vpcName`,
		},
		{
			name: "VPC found, but not subnet",
			edits: editFunctions{
				setValidVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCSubnets = []string{invalidVPCSubnet}
				},
			},
			errorMsg: "VPC.vpcSubnets: Internal error",
		},
		{
			name: "VPC found, subnet found as well, but not attached to the VPC",
			edits: editFunctions{
				setValidVPCName,
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.VPCSubnets = []string{wrongVPCSubnet}
				},
			},
			errorMsg: `VPC.vpcSubnets: Invalid value: "null": not attached to VPC`,
		},
		{
			name: "region specified, VPC found, subnet found, and properly attached",
			edits: editFunctions{
				setValidVPCName,
				setValidVPCRegion,
				func(ic *types.InstallConfig) {
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

	// Mocks: invalid VPC region only
	// nothing to mock

	// Mocks: valid VPC region only
	// nothing to mock

	// Mocks: invalid VPC name results in error
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)

	// Mocks: valid VPC name only, no issues
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)

	// Mocks: valid VPC name, invalid VPC region
	powervsClient.EXPECT().GetVPCs(gomock.Any(), invalidVPCRegion).Return(nil, fmt.Errorf("unknown region"))

	// Mocks: valid VPC name, valid VPC region, all good
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)

	// Mocks: subnet specified, without vpcName, invalid
	// nothing to mock

	// Mocks: valid VPC name, but Subnet not found
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)
	powervsClient.EXPECT().GetSubnetByName(gomock.Any(), invalidVPCSubnet, validVPCRegion).Return(nil, fmt.Errorf(""))

	// Mocks: valid VPC name, but wrong Subnet (present, but not attached)
	powervsClient.EXPECT().GetVPCs(gomock.Any(), validVPCRegion).Return(validVPCs, nil)
	powervsClient.EXPECT().GetSubnetByName(gomock.Any(), wrongVPCSubnet, validVPCRegion).Return(wrongSubnet, nil)

	// Mocks: region specified, valid VPC, valid region, valid Subnet, all good
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

func createControlPlanes(numControlPlanes int, controlPlane *machinev1.PowerVSMachineProviderConfig) []machinev1beta1.Machine {
	controlPlanes := make([]machinev1beta1.Machine, numControlPlanes)

	for i := range controlPlanes {
		masterName := fmt.Sprintf("rdr-hamzy-test3-syd04-zwmgs-master-%d", i)
		controlPlanes[i].TypeMeta = metav1.TypeMeta{
			Kind:       "Machine",
			APIVersion: "machine.openshift.io/v1beta1",
		}
		controlPlanes[i].ObjectMeta = metav1.ObjectMeta{
			Name:      masterName,
			Namespace: "openshift-machine-api",
			Labels:    make(map[string]string),
		}
		controlPlanes[i].Labels["machine.openshift.io/cluster-api-cluster"] = "rdr-hamzy-test3-syd04-zwmgs"
		controlPlanes[i].Labels["machine.openshift.io/cluster-api-machine-role"] = "master"
		controlPlanes[i].Labels["machine.openshift.io/cluster-api-machine-type"] = "master"

		controlPlanes[i].Spec.ProviderSpec = machinev1beta1.ProviderSpec{
			Value: &runtime.RawExtension{
				Raw:    nil,
				Object: controlPlane,
			},
		}
	}

	return controlPlanes
}

func createComputes(numComputes int32, compute *machinev1.PowerVSMachineProviderConfig) []machinev1beta1.MachineSet {
	computes := make([]machinev1beta1.MachineSet, 1)

	computes[0].Spec.Replicas = &numComputes

	computes[0].Spec.Template.Spec.ProviderSpec = machinev1beta1.ProviderSpec{
		Value: &runtime.RawExtension{
			Raw:    nil,
			Object: compute,
		},
	}

	return computes
}

func TestValidatePERAvailability(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name: "Region without PER",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.Zone = regionWithoutPER
				},
			},
			errorMsg: fmt.Sprintf("power-edge-router is not available at: %s", regionWithoutPER),
		},
		{
			name: "Region with PER",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.Zone = regionWithPER
				},
			},
			errorMsg: "",
		},
		{
			name: "Region with no PER availability info",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.Zone = regionPERUnknown
				},
			},
			errorMsg: fmt.Sprintf("power-edge-router capability unknown at: %s", regionPERUnknown),
		},
	}
	setMockEnvVars()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	powervsClient := mock.NewMockAPI(mockCtrl)

	// Mocks: PER-absent region results in false
	powervsClient.EXPECT().GetDatacenterCapabilities(gomock.Any(), regionWithoutPER).Return(mapWithPERFalse, nil)

	// Mocks: PER-enabled region results in true
	powervsClient.EXPECT().GetDatacenterCapabilities(gomock.Any(), regionWithPER).Return(mapWithPERTrue, nil)

	// Mocks: PER-unknown region results in false
	powervsClient.EXPECT().GetDatacenterCapabilities(gomock.Any(), regionPERUnknown).Return(mapPERUnknown, nil)

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := powervs.ValidatePERAvailability(powervsClient, editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestValidateSystemTypeForRegion(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name: "Unknown Region specified",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.Region = invalidRegion
					ic.ControlPlane.Platform.PowerVS = validMachinePool()
					ic.ControlPlane.Platform.PowerVS.SysType = defaultSysType
				},
			},
			errorMsg: fmt.Sprintf("failed to obtain available SysTypes for: %s", invalidRegion),
		},
		{
			name: "No Platform block",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.ControlPlane.Platform.PowerVS = nil
				},
			},
			errorMsg: "",
		},
		{
			name: "Structure present, but no SysType specified",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.ControlPlane.Platform.PowerVS = validMachinePool()
				},
			},
			errorMsg: "",
		},
		{
			name: "Unavailable SysType specified for Dallas Region",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.Region = validRegion
					ic.ControlPlane.Platform.PowerVS = validMachinePool()
					ic.ControlPlane.Platform.PowerVS.SysType = newSysType
				},
			},
			errorMsg: fmt.Sprintf("%s is not available in: %s", newSysType, validRegion),
		},
		{
			name: "Good Region/SysType combo specified",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.Region = validRegion
					ic.ControlPlane.Platform.PowerVS = validMachinePool()
					ic.ControlPlane.Platform.PowerVS.SysType = defaultSysType
				},
			},
			errorMsg: "",
		},
	}
	setMockEnvVars()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	powervsClient := mock.NewMockAPI(mockCtrl)

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := powervs.ValidateSystemTypeForRegion(powervsClient, editedInstallConfig)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}

func TestValidateServiceInstance(t *testing.T) {
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
			name: "invalid install config",
			edits: editFunctions{
				func(ic *types.InstallConfig) {
					ic.Platform.PowerVS.ServiceInstanceGUID = "invalid-uuid"
				},
			},
			errorMsg: "platform:powervs:serviceInstanceGUID has an invalid guid",
		},
	}
	setMockEnvVars()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	powervsClient := mock.NewMockAPI(mockCtrl)

	// FIX: Unexpected call to *mock.MockAPI.ListServiceInstances([context.TODO.WithDeadline(2023-12-02 08:38:15.542340268 -0600 CST m=+300.012357408 [4m59.999979046s])]) at validation.go:289 because: there are no expected calls of the method "ListServiceInstances" for that receiver
	powervsClient.EXPECT().ListServiceInstances(gomock.Any()).AnyTimes()

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			aggregatedErrors := powervs.ValidateServiceInstance(powervsClient, editedInstallConfig)
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
