package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

var (
	validRegion = "us-south"
)

func validMinimalPlatform() *ibmcloud.Platform {
	return &ibmcloud.Platform{
		Region: validRegion,
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
			name: "missing region",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.Region = ""
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
			name: "valid vpc and subnets",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc-subnets"
				p.ControlPlaneSubnets = []string{"cp-1", "cp-2", "cp-3"}
				p.ComputeSubnets = []string{"comp-1", "comp-2"}
				return p
			}(),
			valid: true,
		},
		{
			name: "subnets below ALB limit (1 control plane + 1 compute = 2 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{"cp-1"}
				p.ComputeSubnets = []string{"comp-1"}
				return p
			}(),
			valid: true,
		},
		{
			name: "subnets at ALB limit (12 control plane + 3 compute = 15 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{
					"cp-1", "cp-2", "cp-3", "cp-4", "cp-5", "cp-6",
					"cp-7", "cp-8", "cp-9", "cp-10", "cp-11", "cp-12",
				}
				p.ComputeSubnets = []string{"comp-1", "comp-2", "comp-3"}
				return p
			}(),
			valid: true,
		},
		{
			name: "subnets at ALB limit (8 control plane + 7 compute = 15 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{
					"cp-1", "cp-2", "cp-3", "cp-4", "cp-5", "cp-6", "cp-7", "cp-8",
				}
				p.ComputeSubnets = []string{
					"comp-1", "comp-2", "comp-3", "comp-4", "comp-5", "comp-6", "comp-7",
				}
				return p
			}(),
			valid: true,
		},
		{
			name: "subnets exceed ALB limit (13 control plane + 3 compute = 16 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{
					"cp-1", "cp-2", "cp-3", "cp-4", "cp-5", "cp-6", "cp-7",
					"cp-8", "cp-9", "cp-10", "cp-11", "cp-12", "cp-13",
				}
				p.ComputeSubnets = []string{"comp-1", "comp-2", "comp-3"}
				return p
			}(),
			valid: false,
		},
		{
			name: "subnets exceed ALB limit (8 control plane + 8 compute = 16 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{
					"cp-1", "cp-2", "cp-3", "cp-4", "cp-5", "cp-6", "cp-7", "cp-8",
				}
				p.ComputeSubnets = []string{
					"comp-1", "comp-2", "comp-3", "comp-4", "comp-5", "comp-6", "comp-7", "comp-8",
				}
				return p
			}(),
			valid: false,
		},
		{
			name: "subnets exceed ALB limit (15 control plane + 1 compute = 16 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{
					"cp-1", "cp-2", "cp-3", "cp-4", "cp-5", "cp-6", "cp-7", "cp-8",
					"cp-9", "cp-10", "cp-11", "cp-12", "cp-13", "cp-14", "cp-15",
				}
				p.ComputeSubnets = []string{"comp-1"}
				return p
			}(),
			valid: false,
		},
		{
			name: "subnets exceed ALB limit (1 control plane + 15 compute = 16 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{"cp-1"}
				p.ComputeSubnets = []string{
					"comp-1", "comp-2", "comp-3", "comp-4", "comp-5", "comp-6", "comp-7", "comp-8",
					"comp-9", "comp-10", "comp-11", "comp-12", "comp-13", "comp-14", "comp-15",
				}
				return p
			}(),
			valid: false,
		},
		{
			name: "subnets exceed ALB limit (15 control plane + 15 compute = 30 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{
					"cp-1", "cp-2", "cp-3", "cp-4", "cp-5", "cp-6", "cp-7", "cp-8",
					"cp-9", "cp-10", "cp-11", "cp-12", "cp-13", "cp-14", "cp-15",
				}
				p.ComputeSubnets = []string{
					"comp-1", "comp-2", "comp-3", "comp-4", "comp-5", "comp-6", "comp-7", "comp-8",
					"comp-9", "comp-10", "comp-11", "comp-12", "comp-13", "comp-14", "comp-15",
				}
				return p
			}(),
			valid: false,
		},
		{
			name: "subnets at ALB limit with overlap (8 control plane + 8 compute, 1 shared = 15 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{
					"shared-1", "cp-1", "cp-2", "cp-3", "cp-4", "cp-5", "cp-6", "cp-7",
				}
				p.ComputeSubnets = []string{
					"shared-1", "comp-1", "comp-2", "comp-3", "comp-4", "comp-5", "comp-6", "comp-7",
				}
				return p
			}(),
			valid: true,
		},
		{
			name: "subnets exceed ALB limit with overlap (9 control plane + 9 compute, 2 shared = 16 distinct)",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "valid-vpc"
				p.ControlPlaneSubnets = []string{
					"shared-1", "shared-2", "cp-1", "cp-2", "cp-3", "cp-4", "cp-5", "cp-6", "cp-7",
				}
				p.ComputeSubnets = []string{
					"shared-1", "shared-2", "comp-1", "comp-2", "comp-3", "comp-4", "comp-5", "comp-6", "comp-7",
				}
				return p
			}(),
			valid: false,
		},
		{
			name: "vpc without control plane subnet",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "missing-cp-subnet"
				p.ComputeSubnets = []string{"comp-1", "comp-2"}
				return p
			}(),
			valid: false,
		},
		{
			name: "vpc without compute subnet",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.VPCName = "missing-comp-subnet"
				p.ControlPlaneSubnets = []string{"cp-1", "cp-2"}
				return p
			}(),
			valid: false,
		},
		{
			name: "subnets without vpc",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ControlPlaneSubnets = []string{"cp-1"}
				p.ComputeSubnets = []string{"comp-1"}
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid url (no hostname) for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "/some/path",
				}}
				return p
			}(),
			valid: false,
		},
		{
			name: "invalid url (has path) for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local/some/path",
				}}
				return p
			}(),
			valid: false,
		},
		{
			name: "valid url (has version path, no trailing '/') for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local/v2",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "valid url (has version path and trailing '/') for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local/v35/",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "invalid url (has request) for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local?foo=some",
				}}
				return p
			}(),
			valid: false,
		},
		{
			name: "valid url (no scheme) for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "test-iam.random.local",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "valid url (with scheme) for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "duplicate service endpoints",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "test-iam.random.local",
				}, {
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local",
				}}
				return p
			}(),
			valid: false,
		},
		{
			name: "multiple valid service endpoints",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "test-iam.random.local",
				}, {
					Name: configv1.IBMCloudServiceVPC,
					URL:  "test-vpc.random.local",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "valid url has api and version path, no trailing '/' for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local/api/v22",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "valid url has api and version path, with trailing '/' for service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local/api/v2/",
				}}
				return p
			}(),
			valid: true,
		},
		{
			name: "invalid url has api and version path, with invalid char service endpoint",
			platform: func() *ibmcloud.Platform {
				p := validMinimalPlatform()
				p.ServiceEndpoints = []configv1.IBMCloudServiceEndpoint{{
					Name: configv1.IBMCloudServiceIAM,
					URL:  "https://test-iam.random.local/api/v2e",
				}}
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
