package ibmcloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/mock"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/responses"
	"github.com/openshift/installer/pkg/types"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

func TestValidatePerms(t *testing.T) {
	ctx := context.Background()
	region := "us-south"
	rgName := "valid-resource-group"
	accountID := "valid-account-id"

	validKey := &iamidentityv1.APIKey{AccountID: &accountID}
	validGroups := []resourcemanagerv2.ResourceGroup{
		{Name: &rgName, ID: &rgName},
	}
	emptyVPCs := []vpcv1.VPC{}
	emptyZones := []responses.DNSZoneResponse{}

	forbiddenVPC := errors.New("vpc list failed: status code: 403 Forbidden")
	forbiddenCIS := errors.New("user is not authorized to access CIS")
	forbiddenDNS := errors.New("status 403: not authorized for DNS Services")
	forbiddenIAM := errors.New("GetAPIKeysDetails unauthorized")
	forbiddenRG := errors.New("ListResourceGroups: access is denied")
	forbiddenCOS := errors.New("failed to list cos instances: Forbidden")

	tests := []struct {
		name        string
		ic          *types.InstallConfig
		setup       func(*mock.MockAPI)
		wantErr     bool
		errContains []string
	}{
		{
			name: "success with empty VPC list and Manual credentials",
			ic:   validIBMInstallConfig(region, "", types.ExternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(validKey, nil)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(validGroups, nil)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(emptyVPCs, nil)
				m.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(emptyZones, nil)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, &COSResourceNotFoundError{})
			},
		},
		{
			name: "success when existing resource group is named",
			ic:   validIBMInstallConfig(region, rgName, types.ExternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(validKey, nil)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(validGroups, nil)
				m.EXPECT().GetResourceGroup(gomock.Any(), rgName).Return(&validGroups[0], nil)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(emptyVPCs, nil)
				m.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(emptyZones, nil)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, &COSResourceNotFoundError{})
			},
		},
		{
			name: "missing VPC distinguished from missing CIS",
			ic:   validIBMInstallConfig(region, "", types.ExternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(validKey, nil)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(validGroups, nil)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(nil, forbiddenVPC)
				m.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(nil, forbiddenCIS)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, &COSResourceNotFoundError{})
			},
			wantErr:     true,
			errContains: []string{serviceVPC, serviceCIS},
		},
		{
			name: "internal publish reports DNS Services not CIS",
			ic:   validIBMInstallConfig(region, "", types.InternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(validKey, nil)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(validGroups, nil)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(emptyVPCs, nil)
				m.EXPECT().GetDNSZones(gomock.Any(), types.InternalPublishingStrategy).Return(nil, forbiddenDNS)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, &COSResourceNotFoundError{})
			},
			wantErr:     true,
			errContains: []string{serviceDNS},
		},
		{
			name: "IAM identity forbidden",
			ic:   validIBMInstallConfig(region, "", types.ExternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(nil, forbiddenIAM)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(validGroups, nil)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(emptyVPCs, nil)
				m.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(emptyZones, nil)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, &COSResourceNotFoundError{})
			},
			wantErr:     true,
			errContains: []string{serviceIAMIdentity},
		},
		{
			name: "resource groups forbidden",
			ic:   validIBMInstallConfig(region, "", types.ExternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(validKey, nil)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(nil, forbiddenRG)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(emptyVPCs, nil)
				m.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(emptyZones, nil)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, &COSResourceNotFoundError{})
			},
			wantErr:     true,
			errContains: []string{serviceResourceGroups},
		},
		{
			name: "named resource group get is forbidden",
			ic:   validIBMInstallConfig(region, rgName, types.ExternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(validKey, nil)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(validGroups, nil)
				m.EXPECT().GetResourceGroup(gomock.Any(), rgName).Return(nil, forbiddenRG)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(emptyVPCs, nil)
				m.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(emptyZones, nil)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, &COSResourceNotFoundError{})
			},
			wantErr:     true,
			errContains: []string{serviceResourceGroups},
		},
		{
			name: "COS forbidden",
			ic:   validIBMInstallConfig(region, "", types.ExternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(validKey, nil)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(validGroups, nil)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(emptyVPCs, nil)
				m.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(emptyZones, nil)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, forbiddenCOS)
			},
			wantErr:     true,
			errContains: []string{serviceCOS},
		},
		{
			name: "non-permission VPC error is ignored",
			ic:   validIBMInstallConfig(region, "", types.ExternalPublishingStrategy),
			setup: func(m *mock.MockAPI) {
				m.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(validKey, nil)
				m.EXPECT().GetResourceGroups(gomock.Any()).Return(validGroups, nil)
				m.EXPECT().GetVPCs(gomock.Any(), region).Return(nil, fmt.Errorf("connection reset"))
				m.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(emptyZones, nil)
				m.EXPECT().GetCOSInstanceByName(gomock.Any(), cosPermsProbeName).Return(nil, &COSResourceNotFoundError{})
			},
		},
		{
			name:        "nil ibmcloud platform",
			ic:          &types.InstallConfig{},
			setup:       func(m *mock.MockAPI) {},
			wantErr:     true,
			errContains: []string{"ibmcloud platform configuration is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := mock.NewMockAPI(ctrl)
			tt.setup(client)

			err := ValidatePerms(ctx, client, tt.ic)
			if tt.wantErr {
				assert.Error(t, err)
				for _, needle := range tt.errContains {
					assert.Contains(t, err.Error(), needle)
				}
				if containsAll(tt.errContains, serviceVPC, serviceCIS) {
					assert.Contains(t, err.Error(), serviceVPC)
					assert.Contains(t, err.Error(), serviceCIS)
					assert.NotContains(t, err.Error(), serviceDNS)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestIsForbidden(t *testing.T) {
	assert.False(t, isForbidden(nil))
	assert.False(t, isForbidden(&COSResourceNotFoundError{}))
	assert.False(t, isForbidden(&VPCResourceNotFoundError{}))
	assert.False(t, isForbidden(errors.New("resource group \"foo\" not found")))
	assert.True(t, isForbidden(errors.New("status code: 403")))
	assert.True(t, isForbidden(errors.New("User is not authorized")))
}

func validIBMInstallConfig(region, resourceGroup string, publish types.PublishingStrategy) *types.InstallConfig {
	return &types.InstallConfig{
		CredentialsMode: types.ManualCredentialsMode,
		Publish:         publish,
		Platform: types.Platform{
			IBMCloud: &ibmcloudtypes.Platform{
				Region:            region,
				ResourceGroupName: resourceGroup,
			},
		},
	}
}

func containsAll(got []string, want ...string) bool {
	set := map[string]struct{}{}
	for _, g := range got {
		set[g] = struct{}{}
	}
	for _, w := range want {
		if _, ok := set[w]; !ok {
			return false
		}
	}
	return true
}
