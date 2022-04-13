package azure

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strconv"
	"strings"

	azdns "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/dns/mgmt/dns"
	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	aztypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/azure/defaults"
)

type resourceRequirements struct {
	minimumVCpus  int64
	minimumMemory int64
}

var controlPlaneReq = resourceRequirements{
	minimumVCpus:  4,
	minimumMemory: 16,
}

var computeReq = resourceRequirements{
	minimumVCpus:  2,
	minimumMemory: 8,
}

// Validate executes platform-specific validation.
func Validate(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateNetworks(client, ic.Azure, ic.Networking.MachineNetwork, field.NewPath("platform").Child("azure"))...)
	allErrs = append(allErrs, validateRegion(client, field.NewPath("platform").Child("azure").Child("region"), ic.Azure)...)
	allErrs = append(allErrs, validateInstanceTypes(client, ic)...)
	if ic.Azure.CloudName == aztypes.StackCloud && ic.Azure.ClusterOSImage != "" {
		StorageEndpointSuffix, err := client.GetStorageEndpointSuffix(context.TODO())
		if err != nil {
			return err
		}
		allErrs = append(allErrs, validateAzureStackClusterOSImage(StorageEndpointSuffix, ic.Azure.ClusterOSImage, field.NewPath("platform").Child("azure"))...)
	}
	return allErrs.ToAggregate()
}

// ValidateDiskEncryptionSet ensures the disk encryption set exists and is valid.
func ValidateDiskEncryptionSet(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.Platform.Azure.DefaultMachinePlatform != nil && ic.Platform.Azure.DefaultMachinePlatform.OSDisk.DiskEncryptionSet != nil {
		diskEncryptionSet := ic.Platform.Azure.DefaultMachinePlatform.OSDisk.DiskEncryptionSet
		_, err := client.GetDiskEncryptionSet(context.TODO(), diskEncryptionSet.SubscriptionID, diskEncryptionSet.ResourceGroup, diskEncryptionSet.Name)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("platform").Child("azure", "defaultMachinePlatform", "osDisk", "diskEncryptionSet"), diskEncryptionSet, err.Error()))
		}
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.Azure != nil && ic.ControlPlane.Platform.Azure.OSDisk.DiskEncryptionSet != nil {
		diskEncryptionSet := ic.ControlPlane.Platform.Azure.OSDisk.DiskEncryptionSet
		_, err := client.GetDiskEncryptionSet(context.TODO(), diskEncryptionSet.SubscriptionID, diskEncryptionSet.ResourceGroup, diskEncryptionSet.Name)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("platform").Child("azure", "osDisk", "diskEncryptionSet"), diskEncryptionSet, err.Error()))
		}
	}

	for idx, compute := range ic.Compute {
		fieldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.Azure != nil && compute.Platform.Azure.OSDisk.DiskEncryptionSet != nil {
			diskEncryptionSet := compute.Platform.Azure.OSDisk.DiskEncryptionSet
			_, err := client.GetDiskEncryptionSet(context.TODO(), diskEncryptionSet.SubscriptionID, diskEncryptionSet.ResourceGroup, diskEncryptionSet.Name)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("platform", "azure", "osDisk", "diskEncryptionSet"), diskEncryptionSet, err.Error()))
			}
		}
	}

	return allErrs
}

// ValidateInstanceType ensures the instance type has sufficient Vcpu and Memory.
func ValidateInstanceType(client API, fieldPath *field.Path, region, instanceType string, diskType string, req resourceRequirements, ultraSSDEnabled bool) field.ErrorList {
	allErrs := field.ErrorList{}

	typeMeta, err := client.GetVirtualMachineSku(context.TODO(), instanceType, region)
	if err != nil {
		return append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, err.Error()))
	}

	if typeMeta == nil {
		errMsg := fmt.Sprintf("not found in region %s", region)
		return append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
	}

	ultraSSDAvailable := false
	for _, capability := range *typeMeta.Capabilities {

		if strings.EqualFold(*capability.Name, "vCPUsAvailable") {
			cpus, err := strconv.ParseFloat(*capability.Value, 0)
			if err != nil {
				return append(allErrs, field.InternalError(fieldPath, err))
			}
			if cpus < float64(req.minimumVCpus) {
				errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d vCPUsAvailable", req.minimumVCpus)
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
			}
		} else if strings.EqualFold(*capability.Name, "MemoryGB") {
			memory, err := strconv.ParseFloat(*capability.Value, 0)
			if err != nil {
				return append(allErrs, field.InternalError(fieldPath, err))
			}
			if memory < float64(req.minimumMemory) {
				errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d GB Memory", req.minimumMemory)
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
			}
		} else if diskType == "Premium_LRS" && strings.EqualFold(*capability.Name, "PremiumIO") {
			if strings.EqualFold(*capability.Value, "False") {
				errMsg := fmt.Sprintf("PremiumIO not supported for instance type %s", instanceType)
				allErrs = append(allErrs, field.Invalid(fieldPath.Child("osDisk", "diskType"), diskType, errMsg))
			}
		} else if strings.EqualFold(*capability.Name, "UltraSSDAvailable") {
			ultraSSDAvailable = strings.EqualFold(*capability.Value, "True")
		}
	}

	// The UltraSSDAvailable capability might not be present at all, in which case it must assumed to be false
	if ultraSSDEnabled && !ultraSSDAvailable {
		errMsg := fmt.Sprintf("UltraSSD capability not supported for this instance type in the %s region", region)
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, errMsg))
	}

	return allErrs
}

// validateInstanceTypes checks that the user-provided instance types are valid.
func validateInstanceTypes(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	defaultDiskType := aztypes.DefaultDiskType
	defaultInstanceType := ""
	defaultUltraSSDCapability := "Disabled"

	if ic.Platform.Azure.DefaultMachinePlatform != nil {
		if ic.Platform.Azure.DefaultMachinePlatform.OSDisk.DiskType != "" {
			defaultDiskType = ic.Platform.Azure.DefaultMachinePlatform.OSDisk.DiskType
		}
		if ic.Platform.Azure.DefaultMachinePlatform.InstanceType != "" {
			defaultInstanceType = ic.Platform.Azure.DefaultMachinePlatform.InstanceType
		}
		if ic.Platform.Azure.DefaultMachinePlatform.UltraSSDCapability != "" {
			defaultUltraSSDCapability = ic.Platform.Azure.DefaultMachinePlatform.UltraSSDCapability
		}
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.Azure != nil {
		diskType := ic.ControlPlane.Platform.Azure.OSDisk.DiskType
		instanceType := ic.ControlPlane.Platform.Azure.InstanceType
		ultraSSDCapability := ic.ControlPlane.Platform.Azure.UltraSSDCapability

		if diskType == "" {
			diskType = defaultDiskType
		}
		if instanceType == "" {
			instanceType = defaultInstanceType
		}
		if instanceType == "" {
			instanceType = defaults.ControlPlaneInstanceType(ic.Azure.CloudName, ic.Azure.Region)
		}
		if ultraSSDCapability == "" {
			ultraSSDCapability = defaultUltraSSDCapability
		}
		ultraSSDEnabled := strings.EqualFold(ultraSSDCapability, "Enabled")
		allErrs = append(allErrs, ValidateInstanceType(client, field.NewPath("controlPlane", "platform", "azure"), ic.Azure.Region, instanceType, diskType, controlPlaneReq, ultraSSDEnabled)...)
	}

	for idx, compute := range ic.Compute {
		fieldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.Azure != nil {
			diskType := compute.Platform.Azure.OSDisk.DiskType
			instanceType := compute.Platform.Azure.InstanceType
			ultraSSDCapability := compute.Platform.Azure.UltraSSDCapability

			if diskType == "" {
				diskType = defaultDiskType
			}
			if instanceType == "" {
				instanceType = defaultInstanceType
			}
			if instanceType == "" {
				instanceType = defaults.ComputeInstanceType(ic.Azure.CloudName, ic.Azure.Region)
			}
			if ultraSSDCapability == "" {
				ultraSSDCapability = defaultUltraSSDCapability
			}
			ultraSSDEnabled := strings.EqualFold(ultraSSDCapability, "Enabled")
			allErrs = append(allErrs, ValidateInstanceType(client, fieldPath.Child("platform", "azure"),
				ic.Azure.Region, instanceType, diskType, computeReq, ultraSSDEnabled)...)
		}
	}

	return allErrs
}

// validateNetworks checks that the user-provided VNet and subnets are valid.
func validateNetworks(client API, p *aztypes.Platform, machineNetworks []types.MachineNetworkEntry, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.VirtualNetwork != "" {
		_, err := client.GetVirtualNetwork(context.TODO(), p.NetworkResourceGroupName, p.VirtualNetwork)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("virtualNetwork"), p.VirtualNetwork, err.Error()))
		}

		computeSubnet, err := client.GetComputeSubnet(context.TODO(), p.NetworkResourceGroupName, p.VirtualNetwork, p.ComputeSubnet)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("computeSubnet"), p.ComputeSubnet, "failed to retrieve compute subnet"))
		}

		allErrs = append(allErrs, validateSubnet(client, fieldPath.Child("computeSubnet"), computeSubnet, p.ComputeSubnet, machineNetworks)...)

		controlPlaneSubnet, err := client.GetControlPlaneSubnet(context.TODO(), p.NetworkResourceGroupName, p.VirtualNetwork, p.ControlPlaneSubnet)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("controlPlaneSubnet"), p.ControlPlaneSubnet, "failed to retrieve control plane subnet"))
		}

		allErrs = append(allErrs, validateSubnet(client, fieldPath.Child("controlPlaneSubnet"), controlPlaneSubnet, p.ControlPlaneSubnet, machineNetworks)...)
	}

	return allErrs
}

// validateSubnet checks that the subnet is in the same network as the machine CIDR
func validateSubnet(client API, fieldPath *field.Path, subnet *aznetwork.Subnet, subnetName string, networks []types.MachineNetworkEntry) field.ErrorList {
	allErrs := field.ErrorList{}

	subnetIP, _, err := net.ParseCIDR(*subnet.AddressPrefix)
	if err != nil {
		return append(allErrs, field.Invalid(fieldPath, subnetName, "unable to parse subnet CIDR"))
	}

	allErrs = append(allErrs, validateMachineNetworksContainIP(fieldPath, networks, *subnet.Name, subnetIP)...)
	return allErrs
}

func validateMachineNetworksContainIP(fldPath *field.Path, networks []types.MachineNetworkEntry, subnetName string, ip net.IP) field.ErrorList {
	for _, network := range networks {
		if network.CIDR.Contains(ip) {
			return nil
		}
	}
	return field.ErrorList{field.Invalid(fldPath, subnetName, fmt.Sprintf("subnet %s address prefix is outside of the specified machine networks", ip))}
}

// validateRegion checks that the desired region is valid and available to the user
func validateRegion(client API, fieldPath *field.Path, p *aztypes.Platform) field.ErrorList {
	locations, err := client.ListLocations(context.TODO())
	if err != nil {
		return field.ErrorList{field.InternalError(fieldPath, errors.Wrap(err, "failed to retrieve available regions"))}
	}

	availableRegions := map[string]string{}
	for _, location := range *locations {
		availableRegions[to.String(location.Name)] = to.String(location.DisplayName)
	}

	displayName, ok := availableRegions[p.Region]
	if !ok {
		errMsg := fmt.Sprintf("region %q is not valid or not available for this account", p.Region)

		normalizedRegion := strings.Replace(strings.ToLower(p.Region), " ", "", -1)
		if _, ok := availableRegions[normalizedRegion]; ok {
			errMsg += fmt.Sprintf(", did you mean %q?", normalizedRegion)
		}

		return field.ErrorList{field.Invalid(fieldPath, p.Region, errMsg)}

	}

	provider, err := client.GetResourcesProvider(context.TODO(), "Microsoft.Resources")
	if err != nil {
		return field.ErrorList{field.InternalError(fieldPath, errors.Wrap(err, "failed to retrieve resource capable regions"))}
	}

	for _, resType := range *provider.ResourceTypes {
		if *resType.ResourceType == "resourceGroups" {
			for _, resourceCapableRegion := range *resType.Locations {
				if resourceCapableRegion == displayName {
					return field.ErrorList{}
				}
			}
		}
	}

	return field.ErrorList{field.Invalid(fieldPath, p.Region, fmt.Sprintf("region %q does not support resource creation", p.Region))}
}

// ValidatePublicDNS checks DNS for CNAME, A, and AAA records for
// api.zoneName. If a record exists, it's likely a cluster already exists.
func ValidatePublicDNS(ic *types.InstallConfig, azureDNS *DNSConfig) error {
	// If this is an internal cluster, this check is not necessary
	if ic.Publish == types.InternalPublishingStrategy {
		return nil
	}

	clusterName := ic.ObjectMeta.Name
	record := fmt.Sprintf("api.%s", clusterName)
	rgName := ic.Azure.BaseDomainResourceGroupName
	zoneName := ic.BaseDomain
	fmtStr := "api.%s %s record already exists in %s and might be in use by another cluster, please remove it to continue"

	// Look for an existing CNAME first
	rs, err := azureDNS.GetDNSRecordSet(rgName, zoneName, record, azdns.CNAME)
	if err == nil && rs.CnameRecord != nil {
		return errors.New(fmt.Sprintf(fmtStr, zoneName, azdns.CNAME, clusterName))
	}

	// Look for an A record
	rs, err = azureDNS.GetDNSRecordSet(rgName, zoneName, record, azdns.A)
	if err == nil && rs.ARecords != nil && len(*rs.ARecords) > 0 {
		return errors.New(fmt.Sprintf(fmtStr, zoneName, azdns.A, clusterName))
	}

	// Look for an AAAA record
	rs, err = azureDNS.GetDNSRecordSet(rgName, zoneName, record, azdns.AAAA)
	if err == nil && rs.AaaaRecords != nil && len(*rs.AaaaRecords) > 0 {
		return errors.New(fmt.Sprintf(fmtStr, zoneName, azdns.AAAA, clusterName))
	}

	return nil
}

// ValidateForProvisioning validates if the isntall config if valid for provisioning the cluster.
func ValidateForProvisioning(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateResourceGroup(client, field.NewPath("platform").Child("azure"), ic.Azure)...)
	allErrs = append(allErrs, ValidateDiskEncryptionSet(client, ic)...)
	if ic.Azure.CloudName == aztypes.StackCloud {
		allErrs = append(allErrs, checkAzureStackClusterOSImageSet(ic.Azure.ClusterOSImage, field.NewPath("platform").Child("azure"))...)
	}
	return allErrs.ToAggregate()
}

func validateResourceGroup(client API, fieldPath *field.Path, platform *aztypes.Platform) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(platform.ResourceGroupName) == 0 {
		return allErrs
	}
	group, err := client.GetGroup(context.TODO(), platform.ResourceGroupName)
	if err != nil {
		return append(allErrs, field.InternalError(fieldPath.Child("resourceGroupName"), errors.Wrap(err, "failed to get resource group")))
	}

	normalizedRegion := strings.Replace(strings.ToLower(to.String(group.Location)), " ", "", -1)
	if !strings.EqualFold(normalizedRegion, platform.Region) {
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("resourceGroupName"), platform.ResourceGroupName, fmt.Sprintf("expected to in region %s, but found it to be in %s", platform.Region, normalizedRegion)))
	}

	tagKeys := make([]string, 0, len(group.Tags))
	for key := range group.Tags {
		tagKeys = append(tagKeys, key)
	}
	sort.Strings(tagKeys)
	conflictingTagKeys := tagKeys[:0]
	for _, key := range tagKeys {
		if strings.HasPrefix(key, "kubernetes.io_cluster") {
			conflictingTagKeys = append(conflictingTagKeys, key)
		}
	}
	if len(conflictingTagKeys) > 0 {
		allErrs = append(allErrs, field.Invalid(fieldPath.Child("resourceGroupName"), platform.ResourceGroupName, fmt.Sprintf("resource group has conflicting tags %s", strings.Join(conflictingTagKeys, ", "))))
	}

	// ARO provisions Azure resources before resolving the asset graph.
	if !platform.IsARO() {
		ids, err := client.ListResourceIDsByGroup(context.TODO(), platform.ResourceGroupName)
		if err != nil {
			return append(allErrs, field.InternalError(fieldPath.Child("resourceGroupName"), errors.Wrap(err, "failed to list resources in the resource group")))
		}
		if l := len(ids); l > 0 {
			if len(ids) > 2 {
				ids = ids[:2]
			}
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("resourceGroupName"), platform.ResourceGroupName, fmt.Sprintf("resource group must be empty but it has %d resources like %s ...", l, strings.Join(ids, ", "))))
		}
	}
	return allErrs
}

func checkAzureStackClusterOSImageSet(ClusterOSImage string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if ClusterOSImage == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("clusterOSImage"), "clusterOSImage must be set when installing on Azure Stack"))
	}
	return allErrs
}

func validateAzureStackClusterOSImage(StorageEndpointSuffix string, ClusterOSImage string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	imageParsedURL, err := url.Parse(ClusterOSImage)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterOSImage"), ClusterOSImage, fmt.Errorf("clusterOSImage URL is invalid: %w", err).Error()))
	}
	// If the URL for the image isn't in the Azure Stack environment we can't use it.
	if !strings.HasSuffix(imageParsedURL.Host, StorageEndpointSuffix) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("clusterOSImage"), ClusterOSImage, "clusterOSImage must be in the Azure Stack environment"))
	}
	return allErrs
}
