package roles

import (
	errors "github.com/zgalor/weberr"

	"github.com/openshift/rosa/pkg/arguments"
)

func ValidateSharedVpcInputs(vpcEndpointRoleArn string, route53RoleArn string,
	vpcEndpointRoleArnFlag string, route53RoleArnFlag string) (bool, error) {
	if route53RoleArn != "" && vpcEndpointRoleArn != "" {
		return true, nil
	}

	if vpcEndpointRoleArn != "" && route53RoleArn == "" {
		return false, errors.UserErrorf(
			arguments.MustUseBothFlagsErrorMessage,
			route53RoleArnFlag,
			vpcEndpointRoleArnFlag,
		)
	}

	if route53RoleArn != "" && vpcEndpointRoleArn == "" {
		return false, errors.UserErrorf(
			arguments.MustUseBothFlagsErrorMessage,
			vpcEndpointRoleArnFlag,
			route53RoleArnFlag,
		)
	}

	return false, nil
}
