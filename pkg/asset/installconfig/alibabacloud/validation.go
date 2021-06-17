package alibabacloud

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
)

// Validate executes platform-specific validation.
func Validate(client *Client, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	platformPath := field.NewPath("platform").Child("alibabacloud")
	allErrs = append(allErrs, validatePlatform(client, ic, platformPath)...)

	return allErrs.ToAggregate()
}

func validatePlatform(client *Client, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.Platform.AlibabaCloud.ResourceGroupID != "" {
		allErrs = append(allErrs, validateResourceGroup(client, ic, path)...)
	}
	return allErrs
}

func validateResourceGroup(client *Client, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.AlibabaCloud.ResourceGroupID == "" {
		return allErrs
	}

	resourceGroups, err := client.ListResourceGroups()
	if err != nil {
		return append(allErrs, field.InternalError(path.Child("resourceGroupID"), err))
	}

	found := false
	for _, rg := range resourceGroups.ResourceGroups.ResourceGroup {
		if rg.Id == ic.AlibabaCloud.ResourceGroupID {
			found = true
		}
	}

	if !found {
		return append(allErrs, field.NotFound(path.Child("resourceGroupID"), ic.AlibabaCloud.ResourceGroupID))
	}

	return allErrs
}

// ValidateForProvisioning validates if the install config is valid for provisioning the cluster.
func ValidateForProvisioning(client *Client, ic *types.InstallConfig, metadata *Metadata) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateInstanceType(client, ic, metadata)...)
	return allErrs.ToAggregate()
}

func validateInstanceType(client *Client, ic *types.InstallConfig, metadata *Metadata) field.ErrorList {
	pool := ic.ControlPlane.Platform.AlibabaCloud
	instanceTypePath := field.NewPath("alibabacloud", "instanceType")
	allErrs := field.ErrorList{}

	if ic.ControlPlane == nil || ic.ControlPlane.Platform.AlibabaCloud == nil {
		return allErrs
	}

	if len(pool.Zones) == 0 {
		return allErrs
	}

	if pool.InstanceType != "" {
		for _, zoneID := range pool.Zones {
			response, err := client.DescribeAvailableInstanceType(zoneID, pool.InstanceType)
			if response.AvailableZones.AvailableZone == nil || err != nil {
				allErrs = append(allErrs, field.Invalid(instanceTypePath, pool.InstanceType, fmt.Sprintf("Instance type is Unavailable in zone(%q)", zoneID)))
			}
		}
	}

	return allErrs
}
