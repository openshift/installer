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
