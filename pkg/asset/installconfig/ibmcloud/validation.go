package ibmcloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

// Validate executes platform-specific validation.
func Validate(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	platformPath := field.NewPath("platform").Child("ibmcloud")
	allErrs = append(allErrs, validatePlatform(client, ic, platformPath)...)

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.IBMCloud != nil {
		machinePool := ic.ControlPlane.Platform.IBMCloud
		fldPath := field.NewPath("controlPlane").Child("platform").Child("ibmcloud")
		allErrs = append(allErrs, validateMachinePool(client, ic.Platform.IBMCloud, machinePool, fldPath)...)
	}
	for idx, compute := range ic.Compute {
		machinePool := compute.Platform.IBMCloud
		fldPath := field.NewPath("compute").Index(idx).Child("platform").Child("ibmcloud")
		if machinePool != nil {
			allErrs = append(allErrs, validateMachinePool(client, ic.Platform.IBMCloud, machinePool, fldPath)...)
		}
	}

	return allErrs.ToAggregate()
}

func validatePlatform(client API, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.Platform.IBMCloud.ResourceGroupName != "" {
		allErrs = append(allErrs, validateResourceGroup(client, ic, path)...)
	}

	if ic.Platform.IBMCloud.VPC != "" {
		allErrs = append(allErrs, validateNetworking(client, ic, path)...)
	}

	if ic.Platform.IBMCloud.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, validateMachinePool(client, ic.IBMCloud, ic.Platform.IBMCloud.DefaultMachinePlatform, path)...)
	}
	return allErrs
}

func validateMachinePool(client API, platform *ibmcloud.Platform, machinePool *ibmcloud.MachinePool, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if machinePool.InstanceType != "" {
		allErrs = append(allErrs, validateMachinePoolType(client, machinePool.InstanceType, path.Child("type"))...)
	}

	if len(machinePool.Zones) > 0 {
		allErrs = append(allErrs, validateMachinePoolZones(client, platform.Region, machinePool.Zones, path.Child("zones"))...)
	}

	if machinePool.BootVolume != nil {
		allErrs = append(allErrs, validateMachinePoolBootVolume(client, *machinePool.BootVolume, path.Child("bootVolume"))...)
	}

	if len(machinePool.DedicatedHosts) > 0 {
		allErrs = append(allErrs, validateMachinePoolDedicatedHosts(client, machinePool.DedicatedHosts, machinePool.InstanceType, machinePool.Zones, platform.Region, path.Child("dedicatedHosts"))...)
	}

	return allErrs
}

func validateMachinePoolDedicatedHosts(client API, dhosts []ibmcloud.DedicatedHost, machineType string, zones []string, region string, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Get list of supported profiles in region
	dhostProfiles, err := client.GetDedicatedHostProfiles(context.TODO(), region)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(path, err))
	}

	for i, dhost := range dhosts {
		if dhost.Name != "" {
			// Check if host with name exists
			dh, err := client.GetDedicatedHostByName(context.TODO(), dhost.Name, region)
			if err != nil {
				allErrs = append(allErrs, field.InternalError(path.Index(i).Child("name"), err))
			}

			if dh != nil {
				// Check if instance is provisionable on host
				if !*dh.InstancePlacementEnabled || !*dh.Provisionable {
					allErrs = append(allErrs, field.Invalid(path.Index(i).Child("name"), dhost.Name, "dedicated host is unable to provision instances"))
				}

				// Check if host is in zone
				if *dh.Zone.Name != zones[i] {
					allErrs = append(allErrs, field.Invalid(path.Index(i).Child("name"), dhost.Name, fmt.Sprintf("dedicated host not in zone %s", zones[i])))
				}

				// Check if host profile supports machine type
				if !isInstanceProfileInList(machineType, dh.SupportedInstanceProfiles) {
					allErrs = append(allErrs, field.Invalid(path.Index(i).Child("name"), dhost.Name, fmt.Sprintf("dedicated host does not support machine type %s", machineType)))
				}
			}
		} else {
			// Check if host profile is supported in region
			if !isDedicatedHostProfileInList(dhost.Profile, dhostProfiles) {
				allErrs = append(allErrs, field.Invalid(path.Index(i).Child("profile"), dhost.Profile, fmt.Sprintf("dedicated host profile not supported in region %s", region)))
			}

			// Check if host profile supports machine type
			for _, profile := range dhostProfiles {
				if *profile.Name == dhost.Profile {
					if !isInstanceProfileInList(machineType, profile.SupportedInstanceProfiles) {
						allErrs = append(allErrs, field.Invalid(path.Index(i).Child("profile"), dhost.Profile, fmt.Sprintf("dedicated host profile does not support machine type %s", machineType)))
						break
					}
				}
			}
		}
	}

	return allErrs
}

func isInstanceProfileInList(profile string, list []vpcv1.InstanceProfileReference) bool {
	for _, each := range list {
		if *each.Name == profile {
			return true
		}
	}
	return false
}

func isDedicatedHostProfileInList(profile string, list []vpcv1.DedicatedHostProfile) bool {
	for _, each := range list {
		if *each.Name == profile {
			return true
		}
	}
	return false
}

func validateMachinePoolType(client API, machineType string, path *field.Path) field.ErrorList {
	vsiProfiles, err := client.GetVSIProfiles(context.TODO())
	if err != nil {
		return field.ErrorList{field.InternalError(path, err)}
	}

	for _, profile := range vsiProfiles {
		if *profile.Name == machineType {
			return nil
		}
	}

	return field.ErrorList{field.NotFound(path, machineType)}
}

func validateMachinePoolZones(client API, region string, zones []string, path *field.Path) field.ErrorList {
	regionalZones, err := client.GetVPCZonesForRegion(context.TODO(), region)
	if err != nil {
		return field.ErrorList{field.InternalError(path, err)}
	}

	for idx, zone := range zones {
		validZones := sets.NewString(regionalZones...)
		if !validZones.Has(zone) {
			return field.ErrorList{field.Invalid(path.Index(idx), zone, fmt.Sprintf("zone must be in region %q", region))}
		}
	}
	return nil
}

func validateMachinePoolBootVolume(client API, bootVolume ibmcloud.BootVolume, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if bootVolume.EncryptionKey == "" {
		return allErrs
	}

	// Make sure the encryptionKey exists
	key, err := client.GetEncryptionKey(context.TODO(), bootVolume.EncryptionKey)
	if err != nil {
		return field.ErrorList{field.InternalError(path.Child("encryptionKey"), err)}
	}

	if key == nil {
		return field.ErrorList{field.NotFound(path.Child("encryptionKey"), bootVolume.EncryptionKey)}
	}

	return allErrs
}

func validateResourceGroup(client API, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.IBMCloud.ResourceGroupName == "" {
		return allErrs
	}

	resourceGroups, err := client.GetResourceGroups(context.TODO())
	if err != nil {
		return append(allErrs, field.InternalError(path.Child("resourceGroupName"), err))
	}

	found := false
	for _, rg := range resourceGroups {
		if *rg.ID == ic.IBMCloud.ResourceGroupName || *rg.Name == ic.IBMCloud.ResourceGroupName {
			found = true
		}
	}

	if !found {
		return append(allErrs, field.NotFound(path.Child("resourceGroupName"), ic.IBMCloud.ResourceGroupName))
	}

	return allErrs
}

func validateNetworking(client API, ic *types.InstallConfig, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	platform := ic.Platform.IBMCloud

	_, err := client.GetVPC(context.TODO(), platform.VPC)
	if err != nil {
		if errors.Is(err, &VPCResourceNotFoundError{}) {
			allErrs = append(allErrs, field.NotFound(path.Child("vpc"), platform.VPC))
		} else {
			allErrs = append(allErrs, field.InternalError(path.Child("vpc"), err))
		}
	}

	allErrs = append(allErrs, validateSubnets(client, ic, platform.Subnets, path)...)

	return allErrs
}

func validateSubnets(client API, ic *types.InstallConfig, subnets []string, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	zones, err := client.GetVPCZonesForRegion(context.TODO(), ic.Platform.IBMCloud.Region)
	if err != nil {
		allErrs = append(allErrs, field.InternalError(path.Child("subnets"), err))
	}
	validZones := sets.NewString(zones...)
	for idx, subnet := range subnets {
		subnetPath := path.Child("subnets").Index(idx)
		allErrs = append(allErrs, validateSubnetZone(client, subnet, validZones, subnetPath)...)
	}

	// TODO: IBM[#80]: additional subnet validation
	return allErrs
}

func validateSubnetZone(client API, subnetID string, validZones sets.String, subnetPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if subnet, err := client.GetSubnet(context.TODO(), subnetID); err == nil {
		zoneName := *subnet.Zone.Name
		if !validZones.Has(zoneName) {
			allErrs = append(allErrs, field.Invalid(subnetPath, subnetID, fmt.Sprintf("subnet is not in expected zones: %s", validZones.List())))
		}
	} else {
		if errors.Is(err, &VPCResourceNotFoundError{}) {
			allErrs = append(allErrs, field.NotFound(subnetPath, subnetID))
		} else {
			allErrs = append(allErrs, field.InternalError(subnetPath, err))
		}
	}
	return allErrs
}

// ValidatePreExitingPublicDNS ensure no pre-existing DNS record exists in the CIS
// DNS zone for cluster's Kubernetes API.
func ValidatePreExitingPublicDNS(client API, ic *types.InstallConfig, metadata *Metadata) error {
	// Get CIS CRN
	crn, err := metadata.CISInstanceCRN(context.TODO())
	if err != nil {
		return err
	}

	// Get CIS zone ID by name
	zoneID, err := client.GetDNSZoneIDByName(context.TODO(), ic.BaseDomain)
	if err != nil {
		return field.InternalError(field.NewPath("baseDomain"), err)
	}

	// Get CIS DNS record by name
	recordName := fmt.Sprintf("api.%s", ic.ClusterDomain())
	records, err := client.GetDNSRecordsByName(context.TODO(), crn, zoneID, recordName)
	if err != nil {
		return field.InternalError(field.NewPath("baseDomain"), err)
	}

	// DNS record exists
	if len(records) != 0 {
		return fmt.Errorf("record %s already exists in CIS zone (%s) and might be in use by another cluster, please remove it to continue", recordName, zoneID)
	}

	return nil
}
