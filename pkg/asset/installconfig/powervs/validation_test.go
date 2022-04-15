package powervs_test

import (
	"fmt"
	"testing"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
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
	validCIDR                    = "10.0.0.0/16"
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

	existingDNSRecordsResponse = []dnsrecordsv1.DnsrecordDetails{
		{
			ID: core.StringPtr("valid-dns-record-1"),
		},
		{
			ID: core.StringPtr("valid-dns-record-2"),
		},
	}
	noDNSRecordsResponse = []dnsrecordsv1.DnsrecordDetails{}
	invalidArchitecture  = func(ic *types.InstallConfig) { ic.ControlPlane.Architecture = "ppc64" }
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

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	powervsClient := mock.NewMockAPI(mockCtrl)

	dnsRecordNames := [...]string{fmt.Sprintf("api.%s.%s", validClusterName, validBaseDomain), fmt.Sprintf("api-int.%s.%s", validClusterName, validBaseDomain)}

	metadata := powervs.NewMetadata(validBaseDomain)
	metadata.SetCISInstanceCRN(validCISInstanceCRN)

	// Mocks: no pre-existing DNS records
	powervsClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain).Return(validDNSZoneID, nil)
	for _, dnsRecordName := range dnsRecordNames {
		powervsClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName).Return(noDNSRecordsResponse, nil)
	}

	// Mocks: pre-existing DNS records
	powervsClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain).Return(validDNSZoneID, nil)
	for _, dnsRecordName := range dnsRecordNames {
		powervsClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName).Return(existingDNSRecordsResponse, nil)
	}

	// Mocks: cannot get zone ID
	powervsClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain).Return("", fmt.Errorf(""))

	// Mocks: cannot get DNS records
	powervsClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), validBaseDomain).Return(validDNSZoneID, nil)
	for _, dnsRecordName := range dnsRecordNames {
		powervsClient.EXPECT().GetDNSRecordsByName(gomock.Any(), validCISInstanceCRN, validDNSZoneID, dnsRecordName).Return(nil, fmt.Errorf(""))
	}

	// Run tests
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			validInstallConfig := validInstallConfig()
			aggregatedErrors := powervs.ValidatePreExistingPublicDNS(powervsClient, validInstallConfig, metadata)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, aggregatedErrors)
			} else {
				assert.NoError(t, aggregatedErrors)
			}
		})
	}
}
