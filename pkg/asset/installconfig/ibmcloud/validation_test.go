package ibmcloud_test

import (
	"fmt"
	"testing"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/stretchr/testify/assert"
)

type editFunctions []func(ic *types.InstallConfig)

var (
	validRegion                  = "us-south"
	validCIDR                    = "10.0.0.0/16"
	validClusterOSImage          = "valid-rhcos-image"
	validDNSZoneID               = "valid-zone-id"
	validBaseDomain              = "valid.base.domain"
	validVPC                     = "valid-vpc"
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
	validZoneUSSouth1 = "us-south-1"

	validInstanceProfies = []vpcv1.InstanceProfile{{Name: &[]string{"type-a"}[0]}, {Name: &[]string{"type-b"}[0]}}

	notFoundBaseDomain             = func(ic *types.InstallConfig) { ic.BaseDomain = "notfound.base.domain" }
	notFoundInRegionClusterOSImage = func(ic *types.InstallConfig) { ic.IBMCloud.Region = "us-east" }
	notFoundClusterOSImage         = func(ic *types.InstallConfig) { ic.IBMCloud.ClusterOSImage = "not-found" }
	validVPCConfig                 = func(ic *types.InstallConfig) {
		ic.IBMCloud.VPC = validVPC
		ic.IBMCloud.Subnets = validSubnets
	}
	notFoundVPC            = func(ic *types.InstallConfig) { ic.IBMCloud.VPC = "not-found" }
	internalErrorVPC       = func(ic *types.InstallConfig) { ic.IBMCloud.VPC = "internal-error-vpc" }
	subnetInvalidZone      = func(ic *types.InstallConfig) { ic.IBMCloud.Subnets = []string{"subnet-invalid-zone"} }
	machinePoolInvalidType = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.IBMCloud = &ibmcloudtypes.MachinePool{
			InstanceType: "invalid-type",
		}
	}

	dnsZoneResponses = []ibmcloud.DNSZoneResponse{
		{
			Name: validBaseDomain,
			ID:   "valid-zone-id-1",
		},
		{
			Name: "another.domain",
			ID:   "valid-zone-id-2",
		},
	}
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
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
		Region:         validRegion,
		ClusterOSImage: validClusterOSImage,
	}
}

func validMachinePool() *ibmcloudtypes.MachinePool {
	return &ibmcloudtypes.MachinePool{}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name     string
		edits    editFunctions
		errorMsg string
	}{
		{
			name:     "Valid install config",
			edits:    editFunctions{},
			errorMsg: "",
		},
		{
			name:     "not found clusterOSImage in region",
			edits:    editFunctions{notFoundInRegionClusterOSImage},
			errorMsg: `^platform\.ibmcloud\.clusterOSImage: Not found: "valid-rhcos-image"$`,
		},
		{
			name:     "not found clusterOSImage",
			edits:    editFunctions{notFoundClusterOSImage},
			errorMsg: `^platform\.ibmcloud\.clusterOSImage: Not found: "not-found"$`,
		},
		{
			name:     "valid vpc config",
			edits:    editFunctions{validVPCConfig},
			errorMsg: "",
		},
		{
			name:     "not found vpc",
			edits:    editFunctions{validVPCConfig, notFoundVPC},
			errorMsg: `^platform\.ibmcloud\.vpc: Not found: \"not-found\"$`,
		},
		{
			name:     "internal error vpc",
			edits:    editFunctions{validVPCConfig, internalErrorVPC},
			errorMsg: `^platform\.ibmcloud\.vpc: Internal error$`,
		},
		{
			name:     "subnet invalid zone",
			edits:    editFunctions{validVPCConfig, subnetInvalidZone},
			errorMsg: `^\Qplatform.ibmcloud.subnets[0]: Invalid value: "subnet-invalid-zone": subnet is not in expected zones: [us-south-1 us-south-2 us-south-3]\E$`,
		},
		{
			name:     "machine pool invalid type",
			edits:    editFunctions{validVPCConfig, machinePoolInvalidType},
			errorMsg: `^\QcontrolPlane.platform.ibmcloud.type: Not found: "invalid-type"\E$`,
		},
		{
			name:     "machine pool invalid type",
			edits:    editFunctions{validVPCConfig, machinePoolInvalidType},
			errorMsg: `^\QcontrolPlane.platform.ibmcloud.type: Not found: "invalid-type"\E$`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	ibmcloudClient.EXPECT().GetCustomImageByName(gomock.Any(), validClusterOSImage, validRegion).Return(&vpcv1.Image{}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetCustomImageByName(gomock.Any(), validClusterOSImage, gomock.Not(validRegion)).Return(nil, fmt.Errorf("")).AnyTimes()
	ibmcloudClient.EXPECT().GetCustomImageByName(gomock.Any(), gomock.Not(validClusterOSImage), validRegion).Return(nil, fmt.Errorf("")).AnyTimes()

	ibmcloudClient.EXPECT().GetVPC(gomock.Any(), validVPC).Return(&vpcv1.VPC{}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetVPC(gomock.Any(), "not-found").Return(nil, &ibmcloud.VPCResourceNotFoundError{})
	ibmcloudClient.EXPECT().GetVPC(gomock.Any(), "internal-error-vpc").Return(nil, fmt.Errorf(""))

	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), validPublicSubnetUSSouth1ID).Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &validZoneUSSouth1}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), validPublicSubnetUSSouth2ID).Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &validZoneUSSouth1}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), validPrivateSubnetUSSouth1ID).Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &validZoneUSSouth1}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), validPrivateSubnetUSSouth2ID).Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &validZoneUSSouth1}}, nil).AnyTimes()
	ibmcloudClient.EXPECT().GetSubnet(gomock.Any(), "subnet-invalid-zone").Return(&vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: &[]string{"invalid"}[0]}}, nil).AnyTimes()

	ibmcloudClient.EXPECT().GetVSIProfiles(gomock.Any()).Return(validInstanceProfies, nil).AnyTimes()

	ibmcloudClient.EXPECT().GetVPCZonesForRegion(gomock.Any(), validRegion).Return([]string{"us-south-1", "us-south-2", "us-south-3"}, nil).AnyTimes()

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
