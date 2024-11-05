package arguments

import (
	"github.com/spf13/pflag"
)

const (
	deprecatedInstallerRoleArnFlag = "installer-role-arn"
	newInstallerRoleArnFlag        = "role-arn"
	DeprecatedDefaultMPLabelsFlag  = "default-mp-labels"
	NewDefaultMPLabelsFlag         = "worker-mp-labels"
	DeprecatedControlPlaneIAMRole  = "controlplane-iam-role"
	NewControlPlaneIAMRole         = "controlplane-iam-role-arn"
	DeprecatedWorkerIAMRole        = "worker-iam-role"
	NewWorkerIAMRole               = "worker-iam-role-arn"
	DeprecatedEnvFlag              = "env"
	NewEnvFlag                     = "url"
)

func NormalizeFlags(_ *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case deprecatedInstallerRoleArnFlag:
		name = newInstallerRoleArnFlag
	case DeprecatedDefaultMPLabelsFlag:
		name = NewDefaultMPLabelsFlag
	case DeprecatedWorkerIAMRole:
		name = NewWorkerIAMRole
	case DeprecatedControlPlaneIAMRole:
		name = NewControlPlaneIAMRole
	case DeprecatedEnvFlag:
		name = NewEnvFlag
	}
	return pflag.NormalizedName(name)
}
