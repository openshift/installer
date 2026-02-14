package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

func basicInstallConfig() types.InstallConfig {
	return types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ClusterMetaName",
		},
		Platform: types.Platform{
			AWS: &aws.Platform{},
		},
	}
}

// validBYOSubnetsInstallConfig returns a valid install config for BYO subnets use case.
// Test cases can unset fields if necessary.
func validBYOSubnetsInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		BaseDomain: validDomainName,
		Publish:    types.ExternalPublishingStrategy,
		Platform: types.Platform{
			AWS: &aws.Platform{
				Region: "us-east-1",
				VPC: aws.VPC{
					Subnets: []aws.Subnet{
						{ID: "subnet-valid-private-a"},
						{ID: "subnet-valid-private-b"},
						{ID: "subnet-valid-private-c"},
						{ID: "subnet-valid-public-a"},
						{ID: "subnet-valid-public-b"},
						{ID: "subnet-valid-public-c"},
					},
				},
				HostedZone: validHostedZoneName,
			},
		},
		ControlPlane: &types.MachinePool{
			Architecture: types.ArchitectureAMD64,
			Replicas:     ptr.To[int64](3),
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{
					Zones: []string{"a", "b", "c"},
				},
			},
		},
		Compute: []types.MachinePool{{
			Name:         types.MachinePoolComputeRoleName,
			Architecture: types.ArchitectureAMD64,
			Replicas:     ptr.To[int64](3),
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{
					Zones: []string{"a", "b", "c"},
				},
			},
		}},
		ObjectMeta: metav1.ObjectMeta{
			Name: metaName,
		},
	}
}

func TestIncludesCreateInstanceRole(t *testing.T) {
	t.Run("Should be true when", func(t *testing.T) {
		t.Run("no machine types specified", func(t *testing.T) {
			ic := basicInstallConfig()
			assert.True(t, includesCreateInstanceRole(&ic))
		})
		t.Run("no IAM roles specified", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{},
					},
				},
			}
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{}
			assert.True(t, includesCreateInstanceRole(&ic))
		})
		t.Run("IAM role specified for controlPlane", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						IAMRole: "custom-master-role",
					},
				},
			}
			assert.True(t, includesCreateInstanceRole(&ic))
		})
		t.Run("IAM role specified for compute", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							IAMRole: "custom-master-role",
						},
					},
				},
			}
			assert.True(t, includesCreateInstanceRole(&ic))
		})
	})

	t.Run("Should be false when", func(t *testing.T) {
		t.Run("IAM role specified for defaultMachinePlatform", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMRole: "custom-default-role",
			}
			assert.False(t, includesCreateInstanceRole(&ic))
		})
		t.Run("IAM roles specified for controlPlane and compute", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						IAMRole: "custom-master-role",
					},
				},
			}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							IAMRole: "custom-master-role",
						},
					},
				},
			}
			assert.False(t, includesCreateInstanceRole(&ic))
		})
		t.Run("IAM roles specified for controlPlane and defaultMachinePlatform, compute is nil", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						IAMRole: "custom-master-role",
					},
				},
			}
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMRole: "custom-default-role",
			}
			assert.False(t, includesCreateInstanceRole(&ic))
		})
		t.Run("IAM roles specified for controlPlane and defaultMachinePlatform, compute is not nil", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						IAMRole: "custom-master-role",
					},
				},
			}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{},
					},
				},
			}
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMRole: "custom-default-role",
			}
			assert.False(t, includesCreateInstanceRole(&ic))
		})
		t.Run("IAM roles specified for compute and defaultMachinePlatform", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							IAMRole: "custom-master-role",
						},
					},
				},
			}
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMRole: "custom-default-role",
			}
			assert.False(t, includesCreateInstanceRole(&ic))
		})
	})
}

func TestIncludesExistingInstanceRole(t *testing.T) {
	t.Run("Should be true when", func(t *testing.T) {
		t.Run("IAM role specified for defaultMachinePlatform", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMRole: "custom-default-role",
			}
			assert.True(t, includesExistingInstanceRole(&ic))
		})
		t.Run("IAM role specified for controlPlane", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						IAMRole: "custom-master-role",
					},
				},
			}
			assert.True(t, includesExistingInstanceRole(&ic))
		})
		t.Run("IAM role specified for compute", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							IAMRole: "custom-master-role",
						},
					},
				},
			}
			assert.True(t, includesExistingInstanceRole(&ic))
		})
		t.Run("IAM role specified for controlPlane and compute", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						IAMRole: "custom-master-role",
					},
				},
			}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							IAMRole: "custom-master-role",
						},
					},
				},
			}
			assert.True(t, includesExistingInstanceRole(&ic))
		})
	})
	t.Run("Should be false when", func(t *testing.T) {
		t.Run("no machine types specified", func(t *testing.T) {
			ic := basicInstallConfig()
			assert.False(t, includesExistingInstanceRole(&ic))
		})
		t.Run("no IAM roles specified", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{}
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				},
			}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{},
					},
				},
			}
			assert.False(t, includesExistingInstanceRole(&ic))
		})
	})
}

func TestIncludesExistingInstanceProfile(t *testing.T) {
	t.Run("Should be true when", func(t *testing.T) {
		t.Run("instance profile specified for defaultMachinePlatform", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMProfile: "custom-default-profile",
			}
			assert.True(t, includesExistingInstanceProfile(&ic))
		})
		t.Run("instance profile specified for controlPlane", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						IAMProfile: "custom-master-profile",
					},
				},
			}
			assert.True(t, includesExistingInstanceProfile(&ic))
		})
		t.Run("instance profile specified for compute", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							IAMProfile: "custom-worker-profile",
						},
					},
				},
			}
			assert.True(t, includesExistingInstanceProfile(&ic))
		})
		t.Run("instance profile specified for controlPlane and compute", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						IAMProfile: "custom-master-profile",
					},
				},
			}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							IAMProfile: "custom-worker-profile",
						},
					},
				},
			}
			assert.True(t, includesExistingInstanceProfile(&ic))
		})
	})
	t.Run("Should be false when", func(t *testing.T) {
		t.Run("no machine types specified", func(t *testing.T) {
			ic := basicInstallConfig()
			assert.False(t, includesExistingInstanceProfile(&ic))
		})
		t.Run("no instance profiles specified", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{}
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				},
			}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{},
					},
				},
			}
			assert.False(t, includesExistingInstanceProfile(&ic))
		})
	})
}

func TestIAMRolePermissions(t *testing.T) {
	t.Run("Should include", func(t *testing.T) {
		t.Run("create and delete shared IAM role permissions", func(t *testing.T) {
			t.Run("when role specified for controlPlane", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.ControlPlane.Platform.AWS.IAMRole = "custom-master-role"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceRole)
				assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
			})
			t.Run("when instance profile specified for controlPlane", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.ControlPlane.Platform.AWS.IAMProfile = "custom-master-profile"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceRole)
				assert.NotContains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
			})
			t.Run("when role specified for compute", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.Compute[0].Platform.AWS.IAMRole = "custom-worker-role"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceRole)
				assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
			})
			t.Run("when instance profile specified for compute", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.Compute[0].Platform.AWS.IAMProfile = "custom-worker-profile"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceRole)
				assert.NotContains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
			})
		})
		t.Run("create IAM role permissions", func(t *testing.T) {
			t.Run("when no existing roles and instance profiles are specified", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceRole)
				assert.NotContains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
			})
		})
	})

	t.Run("Should not include create IAM role permissions", func(t *testing.T) {
		t.Run("when role specified for defaultMachinePlatform", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMRole: "custom-default-role",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateInstanceRole)
			assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
		})
		t.Run("when role specified for controlPlane and compute", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS.IAMRole = "custom-master-role"
			ic.Compute[0].Platform.AWS.IAMRole = "custom-worker-role"
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateInstanceRole)
			assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
		})
		t.Run("when instance profile specified for defaultMachinePlatform", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMProfile: "custom-default-profile",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateInstanceRole)
			assert.NotContains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
		})
		t.Run("when instance profile specified for controlPlane and compute", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS.IAMProfile = "custom-master-profile"
			ic.Compute[0].Platform.AWS.IAMProfile = "custom-worker-profile"
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateInstanceRole)
			assert.NotContains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
		})
	})
}

func TestIAMProfilePermissions(t *testing.T) {
	t.Run("Should include", func(t *testing.T) {
		t.Run("create and delete shared instance profile permissions", func(t *testing.T) {
			t.Run("when instance profile specified for controlPlane", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.ControlPlane.Platform.AWS.IAMProfile = "custom-master-profile"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceProfile)
				assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceProfile)
			})
			t.Run("when instance profile specified for compute", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.Compute[0].Platform.AWS.IAMProfile = "custom-worker-profile"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceProfile)
				assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceProfile)
			})
		})
		t.Run("create instance profile permissions", func(t *testing.T) {
			t.Run("when no existing instance profiles are specified", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceProfile)
				assert.NotContains(t, requiredPerms, PermissionDeleteSharedInstanceProfile)
			})
		})
	})

	t.Run("Should not include create instance profile permissions", func(t *testing.T) {
		t.Run("when instance profile specified for defaultMachinePlatform", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMProfile: "custom-default-profile",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateInstanceProfile)
			assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceProfile)
		})
		t.Run("when instance profile specified for controlPlane and compute", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS.IAMProfile = "custom-master-profile"
			ic.Compute[0].Platform.AWS.IAMProfile = "custom-worker-profile"
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateInstanceProfile)
			assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceProfile)
		})
	})
}

func TestIncludesKMSEncryptionKeys(t *testing.T) {
	t.Run("Should be true when", func(t *testing.T) {
		t.Run("KMS key specified for defaultMachinePlatform", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					KMSKeyARN: "custom-default-key",
				},
			}
			assert.True(t, includesKMSEncryptionKey(&ic))
		})
		t.Run("KMS key specified for controlPlane", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						EC2RootVolume: aws.EC2RootVolume{
							KMSKeyARN: "custom-master-key",
						},
					},
				},
			}
			assert.True(t, includesKMSEncryptionKey(&ic))
		})
		t.Run("KMS key specified for compute", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							EC2RootVolume: aws.EC2RootVolume{
								KMSKeyARN: "custom-worker-key",
							},
						},
					},
				},
			}
			assert.True(t, includesKMSEncryptionKey(&ic))
		})
		t.Run("KMS key specified for controlPlane and compute", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						EC2RootVolume: aws.EC2RootVolume{
							KMSKeyARN: "custom-master-key",
						},
					},
				},
			}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{
							EC2RootVolume: aws.EC2RootVolume{
								KMSKeyARN: "custom-worker-key",
							},
						},
					},
				},
			}
			assert.True(t, includesKMSEncryptionKey(&ic))
		})
	})
	t.Run("Should be false when", func(t *testing.T) {
		t.Run("no machine types specified", func(t *testing.T) {
			ic := basicInstallConfig()
			assert.False(t, includesKMSEncryptionKey(&ic))
		})
		t.Run("no KMS keys specified", func(t *testing.T) {
			ic := basicInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{}
			ic.ControlPlane = &types.MachinePool{
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				},
			}
			ic.Compute = []types.MachinePool{
				{
					Platform: types.MachinePoolPlatform{
						AWS: &aws.MachinePool{},
					},
				},
			}
			assert.False(t, includesKMSEncryptionKey(&ic))
		})
	})
}

func TestKMSKeyPermissions(t *testing.T) {
	t.Run("Should include KMS key permissions", func(t *testing.T) {
		t.Run("when KMS key specified for controlPlane", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS.EC2RootVolume = aws.EC2RootVolume{
				KMSKeyARN: "custom-master-key",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionKMSEncryptionKeys)
		})
		t.Run("when KMS key specified for compute", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.Compute[0].Platform.AWS.EC2RootVolume = aws.EC2RootVolume{
				KMSKeyARN: "custom-worker-key",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionKMSEncryptionKeys)
		})
		t.Run("when KMS key specified for defaultMachinePlatform", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				EC2RootVolume: aws.EC2RootVolume{
					KMSKeyARN: "custom-default-key",
				},
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionKMSEncryptionKeys)
		})
	})

	t.Run("Should not include KMS key permissions", func(t *testing.T) {
		t.Run("when no machine types specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane = nil
			ic.Compute = nil
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionKMSEncryptionKeys)
		})
		t.Run("when no KMS keys specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionKMSEncryptionKeys)
		})
	})
}

func TestVPCPermissions(t *testing.T) {
	t.Run("Should include", func(t *testing.T) {
		t.Run("create network permissions when VPC not specified", func(t *testing.T) {
			t.Run("for standard regions", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.AWS.VPC.Subnets = nil
				ic.AWS.HostedZone = ""
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateNetworking)
			})
			t.Run("for secret regions", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.AWS.Region = "us-iso-east-1"
				ic.AWS.VPC.Subnets = nil
				ic.AWS.HostedZone = ""
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateNetworking)
			})
		})
		t.Run("delete network permissions when VPC not specified for standard region", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.VPC.Subnets = nil
			ic.AWS.HostedZone = ""
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionDeleteNetworking)
		})
		t.Run("delete shared network permissions when VPC specified for standard region", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionDeleteSharedNetworking)
		})
	})
	t.Run("Should not include", func(t *testing.T) {
		t.Run("create network permissions when VPC specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateNetworking)
		})
		t.Run("delete network permissions", func(t *testing.T) {
			t.Run("when VPC specified", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteNetworking)
			})
			t.Run("on secret regions", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.AWS.Region = "us-iso-east-1"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteNetworking)
			})
		})
		t.Run("delete shared network permissions", func(t *testing.T) {
			t.Run("when VPC not specified", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.AWS.VPC.Subnets = nil
				ic.AWS.HostedZone = ""
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteSharedNetworking)
			})
			t.Run("on secret regions", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.AWS.Region = "us-iso-east-1"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteSharedNetworking)
			})
		})
	})
}

func TestPrivateZonePermissions(t *testing.T) {
	t.Run("Should include", func(t *testing.T) {
		t.Run("create hosted zone permissions when PHZ not specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.HostedZone = ""
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionCreateHostedZone)
		})
		t.Run("delete hosted zone permissions when PHZ not specified on standard regions", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.HostedZone = ""
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionDeleteHostedZone)
		})
	})
	t.Run("Should not include", func(t *testing.T) {
		t.Run("create hosted zone permissions when PHZ specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateHostedZone)
		})
		t.Run("delete hosted zone permissions", func(t *testing.T) {
			t.Run("on secret regions", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteHostedZone)
			})
			t.Run("when PHZ specified", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteHostedZone)
			})
		})
	})
}

func TestPublicIPv4PoolPermissions(t *testing.T) {
	t.Run("Should include IPv4Pool permissions when IPv4 pool specified", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		ic.AWS.PublicIpv4Pool = "custom-ipv4-pool"
		requiredPerms := RequiredPermissionGroups(ic)
		assert.Contains(t, requiredPerms, PermissionPublicIpv4Pool)
	})
	t.Run("Should not include IPv4Pool permissions when IPv4 pool not specified", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		requiredPerms := RequiredPermissionGroups(ic)
		assert.NotContains(t, requiredPerms, PermissionPublicIpv4Pool)
	})
}

func TestBasePermissions(t *testing.T) {
	t.Run("Should include", func(t *testing.T) {
		t.Run("base create permissions", func(t *testing.T) {
			t.Run("on standard regions", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateBase)
			})
			t.Run("on secret regions", func(t *testing.T) {
				ic := validBYOSubnetsInstallConfig()
				ic.AWS.Region = "us-iso-east-1"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateBase)
			})
		})
		t.Run("base delete permissions on standard regions", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionDeleteBase)
		})
	})
	t.Run("Should not include base delete permissions on secret regions", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		ic.AWS.Region = "us-iso-east-1"
		requiredPerms := RequiredPermissionGroups(ic)
		assert.NotContains(t, requiredPerms, PermissionDeleteBase)
	})
}

func TestDeleteIgnitionPermissions(t *testing.T) {
	t.Run("Should include delete ignition permissions", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		requiredPerms := RequiredPermissionGroups(ic)
		assert.Contains(t, requiredPerms, PermissionDeleteIgnitionObjects)
	})
	t.Run("Should not include delete ignition permission when specified", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		ic.AWS.BestEffortDeleteIgnition = true
		requiredPerms := RequiredPermissionGroups(ic)
		assert.NotContains(t, requiredPerms, PermissionDeleteIgnitionObjects)
	})
}

func TestIncludesInstanceType(t *testing.T) {
	const instanceType = "m7a.2xlarge"
	t.Run("Should be true when instance type specified for", func(t *testing.T) {
		t.Run("defaultMachinePlatform", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				InstanceType: instanceType,
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionValidateInstanceType)
		})
		t.Run("controlPlane", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS.InstanceType = instanceType
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionValidateInstanceType)
		})
		t.Run("compute", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.Compute[0].Platform.AWS.InstanceType = instanceType
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionValidateInstanceType)
		})
	})
	t.Run("Should be false when instance type is not set", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		assert.NotContains(t, RequiredPermissionGroups(ic), PermissionValidateInstanceType)
	})
}

func TestIncludesZones(t *testing.T) {
	t.Run("Should be true when", func(t *testing.T) {
		t.Run("zones specified in defaultMachinePlatform", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS.Zones = []string{}
			ic.Compute[0].Platform.AWS.Zones = []string{}
			ic.AWS.VPC.Subnets = []aws.Subnet{}
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				Zones: []string{"a", "b"},
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionDefaultZones)
		})
		t.Run("zones specified in controlPlane", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.Compute[0].Platform.AWS.Zones = []string{}
			ic.AWS.VPC.Subnets = []aws.Subnet{}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionDefaultZones)
		})
		t.Run("zones specified in compute", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS.Zones = []string{}
			ic.AWS.VPC.Subnets = []aws.Subnet{}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionDefaultZones)
		})
		t.Run("subnets specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS.Zones = []string{}
			ic.Compute[0].Platform.AWS.Zones = []string{}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionDefaultZones)
		})
	})
	t.Run("Should be false when neither zones nor subnets specified", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		ic.AWS.VPC.Subnets = []aws.Subnet{}
		ic.ControlPlane.Platform.AWS.Zones = []string{}
		ic.Compute[0].Platform.AWS.Zones = []string{}
		requiredPerms := RequiredPermissionGroups(ic)
		assert.Contains(t, requiredPerms, PermissionDefaultZones)
	})
}

func TestIncludesAssumeRole(t *testing.T) {
	t.Run("Should be true when IAM role specified", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		ic.AWS.HostedZoneRole = "custom-role"
		requiredPerms := RequiredPermissionGroups(ic)
		assert.Contains(t, requiredPerms, PermissionAssumeRole)
	})
	t.Run("Should be false when IAM role not specified", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		requiredPerms := RequiredPermissionGroups(ic)
		assert.NotContains(t, requiredPerms, PermissionAssumeRole)
	})
}

func TestIncludesWavelengthZones(t *testing.T) {
	t.Run("Should be true when edge compute specified with WL zones", func(t *testing.T) {
		ic := validBYOSubnetsInstallConfig()
		ic.Compute = append(ic.Compute, types.MachinePool{
			Name: "edge",
			Platform: types.MachinePoolPlatform{
				AWS: &aws.MachinePool{
					Zones: []string{"us-west-2-pdx-1a", "us-west-2-wl1-sea-wlz-1"},
				},
			},
		})
		requiredPerms := RequiredPermissionGroups(ic)
		assert.Contains(t, requiredPerms, PermissionCarrierGateway)
	})
	t.Run("Should be false when", func(t *testing.T) {
		t.Run("edge compute specified without WL zones", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.Compute = append(ic.Compute, types.MachinePool{
				Name: "edge",
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						Zones: []string{"us-west-1a", "us-west-2-pdx-1a"},
					},
				},
			})
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCarrierGateway)
		})
		t.Run("edge compute not specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCarrierGateway)
		})
	})
}

func TestIncludesEdgeDefaultInstance(t *testing.T) {
	t.Run("Should be true when at least one edge compute pool specified", func(t *testing.T) {
		t.Run("without platform", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.Compute = append(ic.Compute, types.MachinePool{
				Name: "edge",
			})
			ic.Compute = append(ic.Compute, types.MachinePool{
				Name: "edge",
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						Zones:        []string{"us-west-2-pdx-1a", "us-west-2-wl1-sea-wlz-1"},
						InstanceType: "m6a.xlarge",
					},
				},
			})
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionEdgeDefaultInstance)
		})
		t.Run("without instance type", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.Compute = append(ic.Compute, types.MachinePool{
				Name: "edge",
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						Zones: []string{"us-west-2-pdx-1a", "us-west-2-wl1-sea-wlz-1"},
					},
				},
			})
			ic.Compute = append(ic.Compute, types.MachinePool{
				Name: "edge",
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						Zones:        []string{"us-west-2-pdx-1a", "us-west-2-wl1-sea-wlz-1"},
						InstanceType: "m6a.xlarge",
					},
				},
			})
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionEdgeDefaultInstance)
		})
	})
	t.Run("Should be false when", func(t *testing.T) {
		t.Run("edge compute specified with instance type", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.Compute = append(ic.Compute, types.MachinePool{
				Name: "edge",
				Platform: types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						Zones:        []string{"us-west-1a", "us-west-2-pdx-1a"},
						InstanceType: "m6a.xlarge",
					},
				},
			})
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionEdgeDefaultInstance)
		})
		t.Run("edge compute not specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionEdgeDefaultInstance)
		})
	})
}

func TestAMIEncryptionPermissions(t *testing.T) {
	t.Run("Should include AMI encryption permissions", func(t *testing.T) {
		t.Run("when custom AMI is specified for defaultMachinePlatform", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				AMIID: "custom-presumably-encrypted-ami",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionKMSEncryptionKeys)
			assert.Contains(t, requiredPerms, PermissionAMIEncryptionKeys)
		})
		t.Run("when custom AMI specified for control plane", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane.Platform.AWS = &aws.MachinePool{
				AMIID: "custom-presumably-encrypted-ami",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionKMSEncryptionKeys)
			assert.Contains(t, requiredPerms, PermissionAMIEncryptionKeys)
		})
	})

	t.Run("Should not include AMI encryption permissions", func(t *testing.T) {
		t.Run("when no machine types specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.ControlPlane = nil
			ic.Compute = nil
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionAMIEncryptionKeys)
		})
		t.Run("when no custom AMI specified", func(t *testing.T) {
			ic := validBYOSubnetsInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionAMIEncryptionKeys)
		})
	})
}
