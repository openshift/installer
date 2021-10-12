package alibabacloud

import (
	"fmt"

	"github.com/wxnacy/wgo/arrays"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	alibabacloudtypes "github.com/openshift/installer/pkg/types/alibabacloud"
)

// Validate executes platform-specific validation.
func Validate(client *Client, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	platformPath := field.NewPath("platform").Child("alibabacloud")
	allErrs = append(allErrs, validatePlatform(client, ic, platformPath)...)

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.AlibabaCloud != nil {
		allErrs = append(allErrs, validateMachinePool(client, ic, field.NewPath("controlPlane", "platform", "alibabacloud"), ic.ControlPlane.Platform.AlibabaCloud, ic.ControlPlane.Replicas)...)
	}

	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.AlibabaCloud != nil {
			allErrs = append(allErrs, validateMachinePool(client, ic, fldPath.Child("platform", "alibabacloud"), compute.Platform.AlibabaCloud, compute.Replicas)...)
		}
	}

	return allErrs.ToAggregate()
}

func validatePlatform(client *Client, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateResourceGroup(client, ic, path)...)
	if ic.Platform.AlibabaCloud.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, validateMachinePool(client, ic, path.Child("defaultMachinePlatform"), ic.Platform.AlibabaCloud.DefaultMachinePlatform, nil)...)
	}

	return allErrs
}

func validateMachinePool(client *Client, ic *types.InstallConfig, fldPath *field.Path, pool *alibabacloudtypes.MachinePool, replicas *int64) field.ErrorList {
	allErrs := field.ErrorList{}
	var zones []string
	response, err := client.DescribeAvailableResource("Zone")
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}
	for _, zone := range response.AvailableZones.AvailableZone {
		if zone.Status == "Available" {
			zones = append(zones, zone.ZoneId)
		}
	}

	if len(pool.Zones) > 0 {
		for _, zone := range pool.Zones {
			if index := arrays.ContainsString(zones, zone); index == -1 {
				allErrs = append(allErrs, field.Invalid(fldPath, zone, fmt.Sprintf("zone ID is unavailable in region %q", ic.Platform.AlibabaCloud.Region)))
			}

		}
	}

	if pool.InstanceType != "" {
		if len(pool.Zones) > 0 {
			zones = pool.Zones
		} else {
			zones = zones[:3]
		}

		for _, zoneID := range zones {
			response, err := client.DescribeAvailableInstanceType(zoneID, pool.InstanceType)
			if err != nil {
				allErrs = append(allErrs, field.InternalError(fldPath.Child("instanceType"), err))
			}
			if err == nil && response.AvailableZones.AvailableZone == nil {
				allErrs = append(allErrs, field.Invalid(fldPath, pool.InstanceType, fmt.Sprintf("instance type is unavailable in zone %q", zoneID)))
			}
		}
	}
	return allErrs
}

func validateResourceGroup(client *Client, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	resourceGroups, err := client.ListResourceGroups()
	if err != nil {
		return append(allErrs, field.InternalError(path.Child("resourceGroupID"), err))
	}
	for _, rg := range resourceGroups.ResourceGroups.ResourceGroup {
		if rg.Id == ic.AlibabaCloud.ResourceGroupID {
			return allErrs
		}
	}
	return append(allErrs, field.NotFound(path.Child("resourceGroupID"), ic.AlibabaCloud.ResourceGroupID))
}

// ValidateForProvisioning validates if the install config is valid for provisioning the cluster.
func ValidateForProvisioning(client *Client, ic *types.InstallConfig, metadata *Metadata) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateClusterName(client, ic)...)
	return allErrs.ToAggregate()
}

func validateClusterName(client *Client, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	namePath := field.NewPath("matedata").Child("name")

	zoneName := ic.ClusterDomain()
	response, err := client.ListPrivateZones(zoneName)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(namePath, err))
	}
	if response.TotalItems > 0 {
		allErrs = append(allErrs, field.Invalid(namePath, ic.ClusterName, fmt.Sprintf("cluster name is unavailable, private zone name %s already exists", zoneName)))
	}
	return allErrs
}
