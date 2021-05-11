package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/ibmcloud"
)

var (
	validCRN            = "crn:v1:bluemix:public:internet-svcs:us-south:a/account:instance::"
	validRegion         = "us-south"
	validClusterOSImage = "valid-rhcos-image-name"
)

func validMinimalPlatform() *ibmcloud.Platform {
	return &ibmcloud.Platform{
		Region:         validRegion,
		ClusterOSImage: validClusterOSImage,
	}
}

func validMachinePool() *ibmcloud.MachinePool {
	return &ibmcloud.MachinePool{}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name     string
		platform *ibmcloud.Platform
		valid    bool
	}{
		{
			name:     "minimal",
			platform: validMinimalPlatform(),
			valid:    true,
		},
		{
			name: "invalid region",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.Region = "invalid"
				return p
			}(),
			valid: false,
		},
		{
			name: "missing region",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.Region = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "missing clusterOSImage",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ClusterOSImage = ""
				return p
			}(),
			valid: false,
		},
		{
			name: "valid machine pool",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.DefaultMachinePlatform = validMachinePool()
				return p
			}(),
			valid: true,
		},
		{
			name: "valid vpc config",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPC = "valid-vpc-name"
				p.Subnets = []string{"valid-compute-subnet-id", "valid-control-subnet-id"}
				p.VPCResourceGroupName = "vpc-rg-name"
				return p
			}(),
			valid: true,
		},
		{
			name: "invalid vpc config missing vpc",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.Subnets = []string{"valid-compute-subnet-id", "valid-control-subnet-id"}
				p.VPCResourceGroupName = "vpc-rg-name"
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid vpc config missing subnets",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPC = "valid-vpc-name"
				p.VPCResourceGroupName = "vpc-rg-name"
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid vpc config missing vpcResourceGroupNname",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPC = "valid-vpc-name"
				p.Subnets = []string{"valid-compute-subnet-id", "valid-control-subnet-id"}
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid vpc config missing vpc and subnets",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCResourceGroupName = "vpc-rg-name"
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid vpc config missing vpc and vpcResourceGroupName",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.Subnets = []string{"valid-compute-subnet-id", "valid-control-subnet-id"}
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid vpc config missing subnets and vpcResourceGroupName",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPC = "valid-vpc-name"
				return p
			}(),
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
