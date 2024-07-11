package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

func TestIAMRolePermissions(t *testing.T) {
	t.Run("Should include", func(t *testing.T) {
		t.Run("create and delete shared IAM role permissions", func(t *testing.T) {
			t.Run("when role specified for controlPlane", func(t *testing.T) {
				ic := validInstallConfig()
				ic.ControlPlane.Platform.AWS.IAMRole = "custom-master-role"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceRole)
				assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
			})
			t.Run("when role specified for compute", func(t *testing.T) {
				ic := validInstallConfig()
				ic.Compute[0].Platform.AWS.IAMRole = "custom-worker-role"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceRole)
				assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
			})
		})
		t.Run("create IAM role permissions", func(t *testing.T) {
			t.Run("when no existing roles are specified", func(t *testing.T) {
				ic := validInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateInstanceRole)
				assert.NotContains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
			})
		})
	})

	t.Run("Should not include create IAM role permissions", func(t *testing.T) {
		t.Run("when role specified for defaultMachinePlatform", func(t *testing.T) {
			ic := validInstallConfig()
			ic.AWS.DefaultMachinePlatform = &aws.MachinePool{
				IAMRole: "custom-default-role",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateInstanceRole)
			assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
		})
		t.Run("when role specified for controlPlane and compute", func(t *testing.T) {
			ic := validInstallConfig()
			ic.ControlPlane.Platform.AWS.IAMRole = "custom-master-role"
			ic.Compute[0].Platform.AWS.IAMRole = "custom-worker-role"
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateInstanceRole)
			assert.Contains(t, requiredPerms, PermissionDeleteSharedInstanceRole)
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
			ic := validInstallConfig()
			ic.ControlPlane.Platform.AWS.EC2RootVolume = aws.EC2RootVolume{
				KMSKeyARN: "custom-master-key",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionKMSEncryptionKeys)
		})
		t.Run("when KMS key specified for compute", func(t *testing.T) {
			ic := validInstallConfig()
			ic.Compute[0].Platform.AWS.EC2RootVolume = aws.EC2RootVolume{
				KMSKeyARN: "custom-worker-key",
			}
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionKMSEncryptionKeys)
		})
		t.Run("when KMS key specified for defaultMachinePlatform", func(t *testing.T) {
			ic := validInstallConfig()
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
			ic := validInstallConfig()
			ic.ControlPlane = nil
			ic.Compute = nil
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionKMSEncryptionKeys)
		})
		t.Run("when no KMS keys specified", func(t *testing.T) {
			ic := validInstallConfig()
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
				ic := validInstallConfig()
				ic.AWS.Subnets = nil
				ic.AWS.HostedZone = ""
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateNetworking)
			})
			t.Run("for secret regions", func(t *testing.T) {
				ic := validInstallConfig()
				ic.AWS.Region = "us-iso-east-1"
				ic.AWS.Subnets = nil
				ic.AWS.HostedZone = ""
				requiredPerms := RequiredPermissionGroups(ic)
				assert.Contains(t, requiredPerms, PermissionCreateNetworking)
			})
		})
		t.Run("delete network permissions when VPC not specified for standard region", func(t *testing.T) {
			ic := validInstallConfig()
			ic.AWS.Subnets = nil
			ic.AWS.HostedZone = ""
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionDeleteNetworking)
		})
		t.Run("delete shared network permissions when VPC specified for standard region", func(t *testing.T) {
			ic := validInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionDeleteSharedNetworking)
		})
	})
	t.Run("Should not include", func(t *testing.T) {
		t.Run("create network permissions when VPC specified", func(t *testing.T) {
			ic := validInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateNetworking)
		})
		t.Run("delete network permissions", func(t *testing.T) {
			t.Run("when VPC specified", func(t *testing.T) {
				ic := validInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteNetworking)
			})
			t.Run("on secret regions", func(t *testing.T) {
				ic := validInstallConfig()
				ic.AWS.Region = "us-iso-east-1"
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteNetworking)
			})
		})
		t.Run("delete shared network permissions", func(t *testing.T) {
			t.Run("when VPC not specified", func(t *testing.T) {
				ic := validInstallConfig()
				ic.AWS.Subnets = nil
				ic.AWS.HostedZone = ""
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteSharedNetworking)
			})
			t.Run("on secret regions", func(t *testing.T) {
				ic := validInstallConfig()
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
			ic := validInstallConfig()
			ic.AWS.HostedZone = ""
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionCreateHostedZone)
		})
		t.Run("delete hosted zone permissions when PHZ not specified on standard regions", func(t *testing.T) {
			ic := validInstallConfig()
			ic.AWS.HostedZone = ""
			requiredPerms := RequiredPermissionGroups(ic)
			assert.Contains(t, requiredPerms, PermissionDeleteHostedZone)
		})
	})
	t.Run("Should not include", func(t *testing.T) {
		t.Run("create hosted zone permissions when PHZ specified", func(t *testing.T) {
			ic := validInstallConfig()
			requiredPerms := RequiredPermissionGroups(ic)
			assert.NotContains(t, requiredPerms, PermissionCreateHostedZone)
		})
		t.Run("delete hosted zone permissions", func(t *testing.T) {
			t.Run("on secret regions", func(t *testing.T) {
				ic := validInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteHostedZone)
			})
			t.Run("when PHZ specified", func(t *testing.T) {
				ic := validInstallConfig()
				requiredPerms := RequiredPermissionGroups(ic)
				assert.NotContains(t, requiredPerms, PermissionDeleteHostedZone)
			})
		})
	})
}

func TestPublicIPv4PoolPermissions(t *testing.T) {
	t.Run("Should include IPv4Pool permissions when IPv4 pool specified", func(t *testing.T) {
		ic := validInstallConfig()
		ic.AWS.PublicIpv4Pool = "custom-ipv4-pool"
		requiredPerms := RequiredPermissionGroups(ic)
		assert.Contains(t, requiredPerms, PermissionPublicIpv4Pool)
	})
	t.Run("Should not include IPv4Pool permissions when IPv4 pool not specified", func(t *testing.T) {
		ic := validInstallConfig()
		requiredPerms := RequiredPermissionGroups(ic)
		assert.NotContains(t, requiredPerms, PermissionPublicIpv4Pool)
	})
}
