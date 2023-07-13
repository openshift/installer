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
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
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
	if ic.Azure.CloudName == aztypes.StackCloud {
		allErrs = append(allErrs, validateAzureStackDiskType(client, ic)...)
	}
	allErrs = append(allErrs, validateInstanceTypes(client, ic)...)
	if ic.Azure.CloudName == aztypes.StackCloud && ic.Azure.ClusterOSImage != "" {
		StorageEndpointSuffix, err := client.GetStorageEndpointSuffix(context.TODO())
		if err != nil {
			return err
		}
		allErrs = append(allErrs, validateAzureStackClusterOSImage(StorageEndpointSuffix, ic.Azure.ClusterOSImage, field.NewPath("platform").Child("azure"))...)
	}
	allErrs = append(allErrs, validateMarketplaceImage(client, ic)...)
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

func validatePremiumDisk(fieldPath *field.Path, diskType string, instanceType string, capabilities map[string]string) field.ErrorList {
	fldPath := fieldPath.Child("osDisk", "diskType")
	val, ok := capabilities["PremiumIO"]
	if !ok {
		return field.ErrorList{field.Invalid(fldPath, diskType, "capability not found: PremiumIO")}
	}
	if strings.EqualFold(val, "False") {
		errMsg := fmt.Sprintf("PremiumIO not supported for instance type %s", instanceType)
		return field.ErrorList{field.Invalid(fldPath, diskType, errMsg)}
	}
	return field.ErrorList{}
}

func validateVMArchitecture(fieldPath *field.Path, instanceType string, architecture types.Architecture, capabilities map[string]string) field.ErrorList {
	allErrs := field.ErrorList{}

	val, ok := capabilities["CpuArchitectureType"]
	if ok {
		if architecture != types.ArchitectureARM64 && architecture != types.ArchitectureAMD64 {
			errMsg := fmt.Sprintf("invalid install config architecture %s", architecture)
			allErrs = append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
		} else if (architecture == types.ArchitectureARM64 && !strings.EqualFold(val, "Arm64")) || (architecture == types.ArchitectureAMD64 && !strings.EqualFold(val, "x64")) {
			errMsg := fmt.Sprintf("instance type architecture '%s' does not match install config architecture %s", val, architecture)
			allErrs = append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
		}
	} else {
		logrus.Warnf("Could not determine VM's architecture from its capabilities. Assuming it is %v", architecture)
	}

	return allErrs
}

func validateMininumRequirements(fieldPath *field.Path, req resourceRequirements, instanceType string, capabilities map[string]string) field.ErrorList {
	allErrs := field.ErrorList{}

	val, ok := capabilities["vCPUsAvailable"]
	if ok {
		cpus, err := strconv.ParseFloat(val, 0)
		if err != nil {
			return append(allErrs, field.InternalError(fieldPath, err))
		}
		if cpus < float64(req.minimumVCpus) {
			errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d vCPUsAvailable", req.minimumVCpus)
			allErrs = append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
		}
	} else {
		logrus.Warnf("could not find vCPUsAvailable information for instance type %s", instanceType)
	}

	val, ok = capabilities["MemoryGB"]
	if ok {
		memory, err := strconv.ParseFloat(val, 0)
		if err != nil {
			return append(allErrs, field.InternalError(fieldPath, err))
		}
		if memory < float64(req.minimumMemory) {
			errMsg := fmt.Sprintf("instance type does not meet minimum resource requirements of %d GB Memory", req.minimumMemory)
			allErrs = append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
		}
	} else {
		logrus.Warnf("could not find MemoryGB information for instance type %s", instanceType)
	}

	return allErrs
}

func validateFamily(fieldPath *field.Path, instanceType, family string) field.ErrorList {
	confidentialVMFamilies := sets.NewString(
		"standardDCASv5Family",
		"standardDCADSv5Family",
		"standardECASv5Family",
		"standardECADSv5Family",
		"standardECIADSv5Family",
		"standardECIASv5Family",
	)
	windowsVMFamilies := sets.NewString(
		"standardNVSv4Family",
	)
	allErrs := field.ErrorList{}
	if confidentialVMFamilies.Has(family) {
		errMsg := fmt.Sprintf("%s is not currently supported but will be in a future release", family)
		allErrs = append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
	} else if windowsVMFamilies.Has(family) {
		errMsg := fmt.Sprintf("%s is currently only supported on Windows", family)
		allErrs = append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
	}

	return allErrs
}

func validateAcceleratedNetworking(fieldPath *field.Path, vmNetworkingType string, instanceType string, capabilities map[string]string) field.ErrorList {
	val, ok := capabilities[string(aztypes.AcceleratedNetworkingEnabled)]
	if ok {
		if !strings.EqualFold(val, "True") {
			errMsg := fmt.Sprintf("vm networking type is not supported for instance type %s", instanceType)
			return field.ErrorList{field.Invalid(fieldPath.Child("vmNetworkingType"), vmNetworkingType, errMsg)}
		}
	} else {
		return field.ErrorList{field.Invalid(fieldPath.Child("type"), instanceType, "capability not found: AcceleratedNetworkingEnabled")}
	}

	return field.ErrorList{}
}

func validateUltraSSD(client API, fieldPath *field.Path, icZones []string, region string, instanceType string, capabilities map[string]string) field.ErrorList {
	allErrs := field.ErrorList{}

	locationInfo, err := client.GetLocationInfo(context.TODO(), region, instanceType)
	if err != nil {
		errMsg := fmt.Sprintf("could not determine Availability Zones support in the %s region: %v", region, err)
		return append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
	}
	// If Availability Zones not supported
	if locationInfo == nil || len(to.StringSlice(locationInfo.Zones)) == 0 {
		errMsg := fmt.Sprintf("UltraSSD capability is not compatible with Availability Sets which are used because region %s does not support Availability Zones", region)
		return append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
	}
	allZones := to.StringSlice(locationInfo.Zones)

	// The UltraSSDAvailable capability might not be present at all, in which case it must assumed to be false
	ultraSSDAvailable := false
	if val, ok := capabilities["UltraSSDAvailable"]; ok {
		ultraSSDAvailable = strings.EqualFold(val, "True")
	}
	for _, zoneDetails := range *locationInfo.ZoneDetails {
		for _, capability := range *zoneDetails.Capabilities {
			if !strings.EqualFold(to.String(capability.Name), "UltraSSDAvailable") {
				continue
			}
			if strings.EqualFold(to.String(capability.Value), "True") {
				zones := icZones
				// If no zones are configured in the install config, then all available zones in the region are used
				if len(zones) == 0 {
					zones = allZones
				}
				capZones := to.StringSlice(zoneDetails.Name)
				ultraSSDZones := sets.NewString(capZones...)
				if !ultraSSDZones.HasAll(zones...) {
					errMsg := fmt.Sprintf("UltraSSD capability only supported in zones %v for this instance type in the %s region", capZones, region)
					return append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
				}
				ultraSSDAvailable = true
			}
		}
	}
	if !ultraSSDAvailable {
		errMsg := fmt.Sprintf("UltraSSD capability not supported for this instance type in the %s region", region)
		allErrs = append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
	}

	return allErrs
}

// ValidateInstanceType ensures the instance type has sufficient Vcpu, Memory, and a valid family type.
func ValidateInstanceType(client API, fieldPath *field.Path, region, instanceType, diskType string, req resourceRequirements, ultraSSDEnabled bool, vmNetworkingType string, icZones []string, architecture types.Architecture) field.ErrorList {
	allErrs := field.ErrorList{}

	capabilities, err := client.GetVMCapabilities(context.TODO(), instanceType, region)
	if err != nil {
		return append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, err.Error()))
	}

	allErrs = append(allErrs, validateMininumRequirements(fieldPath.Child("type"), req, instanceType, capabilities)...)
	allErrs = append(allErrs, validateVMArchitecture(fieldPath.Child("type"), instanceType, architecture, capabilities)...)

	family, _ := client.GetVirtualMachineFamily(context.TODO(), instanceType, region)
	if family != "" {
		allErrs = append(allErrs, validateFamily(fieldPath.Child("type"), instanceType, family)...)
	}

	if diskType == "Premium_LRS" {
		allErrs = append(allErrs, validatePremiumDisk(fieldPath, diskType, instanceType, capabilities)...)
	}

	if vmNetworkingType == string(aztypes.VMnetworkingTypeAccelerated) {
		allErrs = append(allErrs, validateAcceleratedNetworking(fieldPath, vmNetworkingType, instanceType, capabilities)...)
	}

	if ultraSSDEnabled {
		allErrs = append(allErrs, validateUltraSSD(client, fieldPath.Child("type"), icZones, region, instanceType, capabilities)...)
	}

	return allErrs
}

// validateInstanceTypes checks that the user-provided instance types are valid.
func validateInstanceTypes(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	defaultDiskType := aztypes.DefaultDiskType
	defaultInstanceType := ""
	defaultUltraSSDCapability := "Disabled"
	defaultVMNetworkingType := ""
	defaultZones := []string{}
	useDefaultInstanceType := false

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
		if ic.Platform.Azure.DefaultMachinePlatform.VMNetworkingType != "" {
			defaultVMNetworkingType = ic.Platform.Azure.DefaultMachinePlatform.VMNetworkingType
		}
		if ic.Platform.Azure.DefaultMachinePlatform.Zones != nil {
			defaultZones = ic.Platform.Azure.DefaultMachinePlatform.Zones
		}
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.Azure != nil {
		fieldPath := field.NewPath("controlPlane", "platform", "azure")
		diskType := ic.ControlPlane.Platform.Azure.OSDisk.DiskType
		instanceType := ic.ControlPlane.Platform.Azure.InstanceType
		ultraSSDCapability := ic.ControlPlane.Platform.Azure.UltraSSDCapability
		vmNetworkingType := ic.ControlPlane.Platform.Azure.VMNetworkingType
		zones := ic.ControlPlane.Platform.Azure.Zones
		architecture := ic.ControlPlane.Architecture

		if diskType == "" {
			diskType = defaultDiskType
		}
		if instanceType == "" {
			instanceType = defaultInstanceType
			useDefaultInstanceType = true
		}
		if instanceType == "" {
			instanceType = defaults.ControlPlaneInstanceType(ic.Azure.CloudName, ic.Azure.Region, architecture)
		}
		if ultraSSDCapability == "" {
			ultraSSDCapability = defaultUltraSSDCapability
		}
		if vmNetworkingType == "" {
			vmNetworkingType = defaultVMNetworkingType
		}
		if len(zones) == 0 {
			zones = defaultZones
		}
		ultraSSDEnabled := strings.EqualFold(ultraSSDCapability, "Enabled")
		allErrs = append(allErrs, ValidateInstanceType(client, fieldPath, ic.Azure.Region, instanceType, diskType, controlPlaneReq, ultraSSDEnabled, vmNetworkingType, zones, architecture)...)
	}

	for idx, compute := range ic.Compute {
		fieldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.Azure != nil {
			diskType := compute.Platform.Azure.OSDisk.DiskType
			instanceType := compute.Platform.Azure.InstanceType
			ultraSSDCapability := compute.Platform.Azure.UltraSSDCapability
			vmNetworkingType := compute.Platform.Azure.VMNetworkingType
			zones := compute.Platform.Azure.Zones
			architecture := compute.Architecture

			if diskType == "" {
				diskType = defaultDiskType
			}
			if instanceType == "" {
				instanceType = defaultInstanceType
				useDefaultInstanceType = true
			}
			if instanceType == "" {
				instanceType = defaults.ComputeInstanceType(ic.Azure.CloudName, ic.Azure.Region, architecture)
			}
			if ultraSSDCapability == "" {
				ultraSSDCapability = defaultUltraSSDCapability
			}
			if vmNetworkingType == "" {
				vmNetworkingType = defaultVMNetworkingType
			}
			if len(zones) == 0 {
				zones = defaultZones
			}
			ultraSSDEnabled := strings.EqualFold(ultraSSDCapability, "Enabled")
			allErrs = append(allErrs, ValidateInstanceType(client, fieldPath.Child("platform", "azure"),
				ic.Azure.Region, instanceType, diskType, computeReq, ultraSSDEnabled, vmNetworkingType, zones, architecture)...)
		}
	}

	// checking here if the defaultInstanceType is present and not used previously. If so, check it now.
	// The assumption here is that since cp and compute arches cannot differ today, it's ok to not check the
	// default instance as long as it is used in any one place.
	if !useDefaultInstanceType && defaultInstanceType != "" {
		if ic.ControlPlane != nil {
			fieldPath := field.NewPath("platform", "azure", "defaultMachinePlatform")
			capabilities, err := client.GetVMCapabilities(context.TODO(), defaultInstanceType, ic.Azure.Region)
			if err != nil {
				return append(allErrs, field.Invalid(fieldPath.Child("type"), defaultInstanceType, err.Error()))
			}
			allErrs = append(allErrs, validateVMArchitecture(fieldPath.Child("type"), defaultInstanceType, ic.ControlPlane.Architecture, capabilities)...)
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

// ValidateForProvisioning validates if the install config is valid for provisioning the cluster.
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

func validateMarketplaceImage(client API, installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList
	for i, compute := range installConfig.Compute {
		platform := compute.Platform.Azure
		if platform == nil {
			continue
		}
		if platform.OSImage.Publisher == "" {
			continue
		}
		osImageFieldPath := field.NewPath("compute").Index(i).Child("platform", "azure", "osImage")
		vmImage, err := client.GetMarketplaceImage(
			context.Background(),
			installConfig.Platform.Azure.Region,
			platform.OSImage.Publisher,
			platform.OSImage.Offer,
			platform.OSImage.SKU,
			platform.OSImage.Version,
		)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(osImageFieldPath, platform.OSImage, err.Error()))
			continue
		}
		instanceType := platform.InstanceType
		if instanceType == "" && installConfig.Platform.Azure.DefaultMachinePlatform != nil {
			instanceType = installConfig.Platform.Azure.DefaultMachinePlatform.InstanceType
		}
		if instanceType == "" {
			instanceType = defaults.ComputeInstanceType(installConfig.Azure.CloudName, installConfig.Azure.Region, compute.Architecture)
		}
		capabilities, err := client.GetVMCapabilities(context.Background(), instanceType, installConfig.Azure.Region)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("compute").Index(i).Child("platform", "azure", "type"), instanceType, err.Error()))
			continue
		}

		generations, err := GetHyperVGenerationVersions(capabilities)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("compute").Index(i).Child("platform", "azure", "type"), instanceType, err.Error()))
			continue
		}
		imageHyperVGen := string(vmImage.HyperVGeneration)
		if !generations.Has(imageHyperVGen) {
			errMsg := fmt.Sprintf("instance type %s supports HyperVGenerations %v but the specified image is for HyperVGeneration %s; to correct this issue either specify a compatible instance type or change the HyperVGeneration for the image by using a different SKU", instanceType, generations.UnsortedList(), imageHyperVGen)
			allErrs = append(allErrs, field.Invalid(osImageFieldPath, platform.OSImage.SKU, errMsg))
			continue
		}

		termsAccepted, err := client.AreMarketplaceImageTermsAccepted(context.Background(), platform.OSImage.Publisher, platform.OSImage.Offer, platform.OSImage.SKU)
		if err == nil {
			if !termsAccepted {
				allErrs = append(allErrs, field.Invalid(osImageFieldPath, platform.OSImage, "the license terms for the marketplace image have not been accepted"))
			}
		} else {
			allErrs = append(allErrs, field.Invalid(osImageFieldPath, platform.OSImage,
				fmt.Sprintf("could not determine if the license terms for the marketplace image have been accepted: %v", err)))
		}
	}
	return allErrs
}

func validateAzureStackDiskType(_ API, installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	supportedTypes := sets.New("Premium_LRS", "Standard_LRS")
	errMsg := fmt.Sprintf("disk format not supported. Must be one of %v", sets.List(supportedTypes))

	defaultDiskType := aztypes.DefaultDiskType
	if installConfig.Azure.DefaultMachinePlatform != nil &&
		installConfig.Azure.DefaultMachinePlatform.DiskType != "" {
		defaultDiskType = installConfig.Azure.DefaultMachinePlatform.DiskType
	}

	diskType := defaultDiskType
	if installConfig.ControlPlane != nil &&
		installConfig.ControlPlane.Platform.Azure != nil &&
		installConfig.ControlPlane.Platform.Azure.DiskType != "" {
		diskType = installConfig.ControlPlane.Platform.Azure.DiskType
	}
	if !supportedTypes.Has(diskType) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("controlPlane", "platform", "azure", "OSDisk", "diskType"), diskType, errMsg))
	}

	for idx, compute := range installConfig.Compute {
		fieldPath := field.NewPath("compute").Index(idx)
		diskType := defaultDiskType
		if compute.Platform.Azure != nil && compute.Platform.Azure.DiskType != "" {
			diskType = compute.Platform.Azure.DiskType
		}
		if !supportedTypes.Has(diskType) {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("platform", "azure", "OSDisk", "diskType"), diskType, errMsg))
		}
	}

	return allErrs
}
