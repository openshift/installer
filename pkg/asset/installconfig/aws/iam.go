package aws

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// IncludesExistingInstanceRole checks if at least one BYO instance role is included in the install-config.
func IncludesExistingInstanceRole(installConfig *types.InstallConfig) bool {
	{
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		if mp := installConfig.ControlPlane; mp != nil {
			mpool.Set(mp.Platform.AWS)
		}
		if len(mpool.IAMRole) > 0 {
			return true
		}
	}

	for _, compute := range installConfig.Compute {
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		mpool.Set(compute.Platform.AWS)
		if len(mpool.IAMRole) > 0 {
			return true
		}
	}

	// If compute stanza is not defined, we know it'll inherit the value from DefaultMachinePlatform
	mpool := installConfig.AWS.DefaultMachinePlatform
	return mpool != nil && len(mpool.IAMRole) > 0
}

// IncludesCreateInstanceRole checks if at least one instance role will be created by the installer.
func IncludesCreateInstanceRole(installConfig *types.InstallConfig) bool {
	{
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		if mp := installConfig.ControlPlane; mp != nil {
			mpool.Set(mp.Platform.AWS)
		}
		if len(mpool.IAMRole) == 0 {
			return true
		}
	}

	for _, compute := range installConfig.Compute {
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		mpool.Set(compute.Platform.AWS)
		if len(mpool.IAMRole) == 0 {
			return true
		}
	}

	// If compute stanza is not defined, we know it'll inherit the value from DefaultMachinePlatform
	mpool := installConfig.AWS.DefaultMachinePlatform
	return mpool == nil || len(mpool.IAMRole) == 0
}

// IncludesExistingInstanceProfile checks if at least one BYO instance profile is included in the install-config.
func IncludesExistingInstanceProfile(installConfig *types.InstallConfig) bool {
	{
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		if mp := installConfig.ControlPlane; mp != nil {
			mpool.Set(mp.Platform.AWS)
		}
		if len(mpool.IAMProfile) > 0 {
			return true
		}
	}

	for _, compute := range installConfig.Compute {
		mpool := aws.MachinePool{}
		mpool.Set(installConfig.AWS.DefaultMachinePlatform)
		mpool.Set(compute.Platform.AWS)
		if len(mpool.IAMProfile) > 0 {
			return true
		}
	}

	// If compute stanza is not defined, we know it'll inherit the value from DefaultMachinePlatform
	mpool := installConfig.AWS.DefaultMachinePlatform
	return mpool != nil && len(mpool.IAMProfile) > 0
}
