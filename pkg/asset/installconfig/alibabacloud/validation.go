package alibabacloud

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// Validate executes platform-specific validation.
func Validate(client *Client, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	platformPath := field.NewPath("platorm").Child("alibabacloud")
	allErrs = append(allErrs, validatePlatform(client, ic, platformPath)...)

	return allErrs.ToAggregate()
}

func validatePlatform(client *Client, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.Platform.AlibabaCloud.ResourceGroupName != "" {
		allErrs = append(allErrs, validateResourceGroup(client, ic, path)...)
	}
	return allErrs
}

func validateResourceGroup(client *Client, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.AlibabaCloud.ResourceGroupName == "" {
		return allErrs
	}

	resourceGroups, err := client.ListResourceGroups()
	if err != nil {
		return append(allErrs, field.InternalError(path.Child("resourceGroupName"), err))
	}

	found := false
	for _, rg := range resourceGroups.ResourceGroups.ResourceGroup {
		if rg.Id == ic.AlibabaCloud.ResourceGroupName || rg.Name == ic.AlibabaCloud.ResourceGroupName {
			found = true
		}
	}

	if !found {
		return append(allErrs, field.NotFound(path.Child("resourceGroupName"), ic.AlibabaCloud.ResourceGroupName))
	}

	return allErrs
}
