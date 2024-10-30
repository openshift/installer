package azure

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	azdns "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/dns/mgmt/dns"
	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	azenc "github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

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
	allErrs = append(allErrs, validateMarketplaceImages(client, ic)...)
	allErrs = append(allErrs, validateBootDiagnostics(client, ic)...)
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

func validateConfidentialDiskEncryptionSet(client API, diskEncryptionSet *aztypes.DiskEncryptionSet, desFieldPath *field.Path) error {
	resp, requestErr := client.GetDiskEncryptionSet(context.TODO(), diskEncryptionSet.SubscriptionID, diskEncryptionSet.ResourceGroup, diskEncryptionSet.Name)
	if requestErr != nil {
		return requestErr
	} else if resp == nil || resp.EncryptionSetProperties == nil || resp.EncryptionSetProperties.EncryptionType != azenc.ConfidentialVMEncryptedWithCustomerKey {
		return fmt.Errorf("the disk encryption set should be created with type %s", azenc.ConfidentialVMEncryptedWithCustomerKey)
	}
	return nil
}

// ValidateSecurityProfileDiskEncryptionSet ensures the security profile disk encryption set exists and is valid.
func ValidateSecurityProfileDiskEncryptionSet(client API, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	if ic.Platform.Azure.DefaultMachinePlatform != nil &&
		ic.Platform.Azure.DefaultMachinePlatform.OSDisk.SecurityProfile != nil &&
		ic.Platform.Azure.DefaultMachinePlatform.OSDisk.SecurityProfile.DiskEncryptionSet != nil {
		desFieldPath := field.NewPath("platform").Child("azure", "defaultMachinePlatform", "osDisk", "securityProfile", "diskEncryptionSet")
		diskEncryptionSet := ic.Platform.Azure.DefaultMachinePlatform.OSDisk.SecurityProfile.DiskEncryptionSet
		err := validateConfidentialDiskEncryptionSet(client, diskEncryptionSet, desFieldPath)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(desFieldPath, diskEncryptionSet, err.Error()))
		}
	}

	if ic.ControlPlane != nil &&
		ic.ControlPlane.Platform.Azure != nil &&
		ic.ControlPlane.Platform.Azure.OSDisk.SecurityProfile != nil &&
		ic.ControlPlane.Platform.Azure.OSDisk.SecurityProfile.DiskEncryptionSet != nil {
		desFieldPath := field.NewPath("platform").Child("azure", "osDisk", "securityProfile", "diskEncryptionSet")
		diskEncryptionSet := ic.ControlPlane.Platform.Azure.OSDisk.SecurityProfile.DiskEncryptionSet
		err := validateConfidentialDiskEncryptionSet(client, diskEncryptionSet, desFieldPath)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(desFieldPath, diskEncryptionSet, err.Error()))
		}
	}

	for idx, compute := range ic.Compute {
		fieldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.Azure != nil &&
			compute.Platform.Azure.OSDisk.SecurityProfile != nil &&
			compute.Platform.Azure.OSDisk.SecurityProfile.DiskEncryptionSet != nil {
			desFieldPath := fieldPath.Child("platform", "azure", "osDisk", "securityProfile", "diskEncryptionSet")
			diskEncryptionSet := compute.Platform.Azure.OSDisk.SecurityProfile.DiskEncryptionSet
			err := validateConfidentialDiskEncryptionSet(client, diskEncryptionSet, desFieldPath)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(desFieldPath, diskEncryptionSet, err.Error()))
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

func validateSecurityType(fieldPath *field.Path, securityType aztypes.SecurityTypes, instanceType string, capabilities map[string]string) field.ErrorList {
	allErrs := field.ErrorList{}

	_, hasTrustedLaunchDisabled := capabilities["TrustedLaunchDisabled"]
	confidentialComputingType, hasConfidentialComputingType := capabilities["ConfidentialComputingType"]
	isConfidentialComputingTypeSNP := confidentialComputingType == "SNP"

	var reason string
	supportedSecurityType := true
	switch securityType {
	case aztypes.SecurityTypesConfidentialVM:
		supportedSecurityType = hasConfidentialComputingType && isConfidentialComputingTypeSNP

		if !hasConfidentialComputingType {
			reason = "no support for Confidential Computing"
		} else if !isConfidentialComputingTypeSNP {
			reason = "no support for AMD-SEV SNP"
		}
	case aztypes.SecurityTypesTrustedLaunch:
		supportedSecurityType = !(hasTrustedLaunchDisabled || hasConfidentialComputingType)

		if hasTrustedLaunchDisabled {
			reason = "no support for Trusted Launch"
		} else if hasConfidentialComputingType {
			reason = "confidential VMs do not support Trusted Launch for VMs"
		}
	}

	if !supportedSecurityType {
		errMsg := fmt.Sprintf("this security type is not supported for instance type %s, %s", instanceType, reason)
		allErrs = append(allErrs, field.Invalid(fieldPath, securityType, errMsg))
	}

	return allErrs
}

func validateFamily(fieldPath *field.Path, instanceType, family string) field.ErrorList {
	windowsVMFamilies := sets.NewString(
		"standardNVSv4Family",
	)
	diskNVMeVMFamilies := sets.NewString(
		"standardEIBDSv5Family",
		"standardEIBSv5Family",
	)
	allErrs := field.ErrorList{}
	if windowsVMFamilies.Has(family) {
		errMsg := fmt.Sprintf("%s is currently only supported on Windows", family)
		allErrs = append(allErrs, field.Invalid(fieldPath, instanceType, errMsg))
	}
	// FIXME: remove when supported has been added to the provider
	// https://github.com/hashicorp/terraform-provider-azurerm/issues/22058
	if diskNVMeVMFamilies.Has(family) {
		errMsg := fmt.Sprintf("%s is not currently supported but might be in a future release", family)
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
func ValidateInstanceType(client API, fieldPath *field.Path, region, instanceType, diskType string, req resourceRequirements, ultraSSDEnabled bool, vmNetworkingType string, icZones []string, architecture types.Architecture, securityType aztypes.SecurityTypes) field.ErrorList {
	allErrs := field.ErrorList{}

	capabilities, err := client.GetVMCapabilities(context.TODO(), instanceType, region)
	if err != nil {
		return append(allErrs, field.Invalid(fieldPath.Child("type"), instanceType, err.Error()))
	}

	allErrs = append(allErrs, validateMininumRequirements(fieldPath.Child("type"), req, instanceType, capabilities)...)
	allErrs = append(allErrs, validateVMArchitecture(fieldPath.Child("type"), instanceType, architecture, capabilities)...)
	allErrs = append(allErrs, validateSecurityType(fieldPath.Child("settings", "securityType"), securityType, instanceType, capabilities)...)

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

	var securityType aztypes.SecurityTypes

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
		if ic.Platform.Azure.DefaultMachinePlatform.Settings != nil {
			securityType = ic.Platform.Azure.DefaultMachinePlatform.Settings.SecurityType
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

		if ic.ControlPlane.Platform.Azure.Settings != nil {
			securityType = ic.ControlPlane.Platform.Azure.Settings.SecurityType
		}
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
		allErrs = append(allErrs, ValidateInstanceType(client, fieldPath, ic.Azure.Region, instanceType, diskType, controlPlaneReq, ultraSSDEnabled, vmNetworkingType, zones, architecture, securityType)...)
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

			if compute.Platform.Azure.Settings != nil {
				securityType = compute.Platform.Azure.Settings.SecurityType
			}
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
				ic.Azure.Region, instanceType, diskType, computeReq, ultraSSDEnabled, vmNetworkingType, zones, architecture, securityType)...)
		}
	}

	// checking here if the defaultInstanceType is present and not used previously. If so, check it now.
	// The assumption here is that since cp and compute arches cannot differ today, it's ok to not check the
	// default instance as long as it is used in any one place.
	if !useDefaultInstanceType && defaultInstanceType != "" {
		architecture := types.Architecture(types.ArchitectureAMD64)
		if ic.ControlPlane != nil {
			architecture = ic.ControlPlane.Architecture
		}
		minReq := computeReq
		if ic.ControlPlane == nil || ic.ControlPlane.Platform.Azure == nil {
			minReq = controlPlaneReq
		}
		fieldPath := field.NewPath("platform", "azure", "defaultMachinePlatform")
		ultraSSDEnabled := strings.EqualFold(defaultUltraSSDCapability, "Enabled")
		allErrs = append(allErrs, ValidateInstanceType(client, fieldPath,
			ic.Azure.Region, defaultInstanceType, defaultDiskType, minReq, ultraSSDEnabled, defaultVMNetworkingType, defaultZones, architecture, securityType)...)
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

	var addressPrefix string
	switch {
	case subnet.AddressPrefix != nil:
		addressPrefix = *subnet.AddressPrefix
	// NOTE: if the subscription has the `AllowMultipleAddressPrefixesOnSubnet` feature, the Azure API will return a
	// `addressPrefixes` field with a slice of addresses instead of a single value via `addressPrefix`.
	case subnet.AddressPrefixes != nil && len(*subnet.AddressPrefixes) > 0:
		addressPrefix = (*subnet.AddressPrefixes)[0]
	default:
		return append(allErrs, field.Invalid(fieldPath, subnetName, "subnet does not have an address prefix"))
	}

	subnetIP, _, err := net.ParseCIDR(addressPrefix)
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
		return field.ErrorList{field.InternalError(fieldPath, fmt.Errorf("failed to retrieve available regions: %w", err))}
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
		return field.ErrorList{field.InternalError(fieldPath, fmt.Errorf("failed to retrieve resource capable regions: %w", err))}
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
		return fmt.Errorf(fmtStr, zoneName, azdns.CNAME, clusterName)
	}

	// Look for an A record
	rs, err = azureDNS.GetDNSRecordSet(rgName, zoneName, record, azdns.A)
	if err == nil && rs.ARecords != nil && len(*rs.ARecords) > 0 {
		return fmt.Errorf(fmtStr, zoneName, azdns.A, clusterName)
	}

	// Look for an AAAA record
	rs, err = azureDNS.GetDNSRecordSet(rgName, zoneName, record, azdns.AAAA)
	if err == nil && rs.AaaaRecords != nil && len(*rs.AaaaRecords) > 0 {
		return fmt.Errorf(fmtStr, zoneName, azdns.AAAA, clusterName)
	}

	return nil
}

// ValidateForProvisioning validates if the install config is valid for provisioning the cluster.
func ValidateForProvisioning(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateResourceGroup(client, field.NewPath("platform").Child("azure"), ic.Azure)...)
	allErrs = append(allErrs, ValidateDiskEncryptionSet(client, ic)...)
	allErrs = append(allErrs, ValidateSecurityProfileDiskEncryptionSet(client, ic)...)
	allErrs = append(allErrs, validateSkipImageUpload(field.NewPath("image"), ic)...)
	if ic.Azure.CloudName == aztypes.StackCloud {
		allErrs = append(allErrs, checkAzureStackClusterOSImageSet(ic.Azure.ClusterOSImage, field.NewPath("platform").Child("azure"))...)
	}
	return allErrs.ToAggregate()
}

func validateSkipImageUpload(fieldPath *field.Path, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	defaultOSImage := aztypes.OSImage{}
	if ic.Azure.DefaultMachinePlatform != nil {
		defaultOSImage = aztypes.OSImage{
			Plan:      ic.Azure.DefaultMachinePlatform.OSImage.Plan,
			Publisher: ic.Azure.DefaultMachinePlatform.OSImage.Publisher,
			SKU:       ic.Azure.DefaultMachinePlatform.OSImage.SKU,
			Version:   ic.Azure.DefaultMachinePlatform.OSImage.Version,
			Offer:     ic.Azure.DefaultMachinePlatform.OSImage.Offer,
		}
	}
	controlPlaneOSImage := defaultOSImage
	if ic.ControlPlane.Platform.Azure != nil {
		controlPlaneOSImage = ic.ControlPlane.Platform.Azure.OSImage
	}
	allErrs = append(allErrs, validateOSImage(fieldPath.Child("controlplane"), controlPlaneOSImage)...)
	computeOSImage := defaultOSImage
	if len(ic.Compute) > 0 && ic.Compute[0].Platform.Azure != nil {
		computeOSImage = ic.Compute[0].Platform.Azure.OSImage
	}
	allErrs = append(allErrs, validateOSImage(fieldPath.Child("compute"), computeOSImage)...)
	return allErrs
}
func validateOSImage(fieldPath *field.Path, osImage aztypes.OSImage) field.ErrorList {
	if _, ok := os.LookupEnv("OPENSHIFT_INSTALL_SKIP_IMAGE_UPLOAD"); ok {
		if len(osImage.SKU) > 0 {
			return nil
		}
		return field.ErrorList{field.Invalid(fieldPath, "image", "cannot skip image upload without marketplace image specified")}
	}
	return nil
}
func validateResourceGroup(client API, fieldPath *field.Path, platform *aztypes.Platform) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(platform.ResourceGroupName) == 0 {
		return allErrs
	}
	group, err := client.GetGroup(context.TODO(), platform.ResourceGroupName)
	if err != nil {
		return append(allErrs, field.InternalError(fieldPath.Child("resourceGroupName"), fmt.Errorf("failed to get resource group: %w", err)))
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
			return append(allErrs, field.InternalError(fieldPath.Child("resourceGroupName"), fmt.Errorf("failed to list resources in the resource group: %w", err)))
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

func validateMarketplaceImages(client API, installConfig *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	region := installConfig.Azure.Region
	cloudName := installConfig.Azure.CloudName

	var defaultInstanceType string
	var defaultOSImage aztypes.OSImage
	if installConfig.Azure.DefaultMachinePlatform != nil {
		defaultInstanceType = installConfig.Azure.DefaultMachinePlatform.InstanceType
		defaultOSImage = installConfig.Azure.DefaultMachinePlatform.OSImage
	}

	// Validate ControlPlane marketplace images
	if installConfig.ControlPlane != nil {
		platform := installConfig.ControlPlane.Platform.Azure
		fldPath := field.NewPath("controlPlane")

		// Determine instance type
		instanceType := ""
		if platform != nil {
			instanceType = platform.InstanceType
		}
		if instanceType == "" {
			instanceType = defaultInstanceType
		}
		if instanceType == "" {
			instanceType = defaults.ControlPlaneInstanceType(cloudName, region, installConfig.ControlPlane.Architecture)
		}

		capabilities, err := client.GetVMCapabilities(context.Background(), instanceType, region)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("platform", "azure", "type"), instanceType, err.Error()))
		}

		generations, err := GetHyperVGenerationVersions(capabilities)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("platform", "azure", "type"), instanceType, err.Error()))
		}

		// If not set, try to use the OS Image definition from the default machine pool
		var osImage aztypes.OSImage
		if platform != nil {
			osImage = platform.OSImage
		}
		if osImage.Publisher == "" {
			osImage = defaultOSImage
		}

		imgErr := validateMarketplaceImage(client, region, generations, &osImage, fldPath)
		if imgErr != nil {
			allErrs = append(allErrs, imgErr)
		}
	}

	// Validate Compute marketplace images
	for i, compute := range installConfig.Compute {
		platform := compute.Platform.Azure
		fldPath := field.NewPath("compute").Index(i)

		// Determine instance type
		instanceType := ""
		if platform != nil {
			instanceType = platform.InstanceType
		}
		if instanceType == "" {
			instanceType = defaultInstanceType
		}
		if instanceType == "" {
			instanceType = defaults.ComputeInstanceType(cloudName, region, compute.Architecture)
		}

		capabilities, err := client.GetVMCapabilities(context.Background(), instanceType, region)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("platform", "azure", "type"), instanceType, err.Error()))
			continue
		}

		generations, err := GetHyperVGenerationVersions(capabilities)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("platform", "azure", "type"), instanceType, err.Error()))
			continue
		}

		// If not set, try to use the OS Image definition from the default machine pool
		var osImage aztypes.OSImage
		if platform != nil {
			osImage = platform.OSImage
		}
		if osImage.Publisher == "" {
			osImage = defaultOSImage
		}
		imgErr := validateMarketplaceImage(client, region, generations, &osImage, fldPath)
		if imgErr != nil {
			allErrs = append(allErrs, imgErr)
		}
	}

	return allErrs
}

func validateMarketplaceImage(client API, region string, instanceHyperVGenSet sets.Set[string], osImage *aztypes.OSImage, fldPath *field.Path) *field.Error {
	// Marketplace image not specified
	if osImage.Publisher == "" {
		return nil
	}

	osImageFieldPath := fldPath.Child("platform", "azure", "osImage")
	vmImage, err := client.GetMarketplaceImage(
		context.Background(),
		region,
		osImage.Publisher,
		osImage.Offer,
		osImage.SKU,
		osImage.Version,
	)
	if err != nil {
		return field.Invalid(osImageFieldPath, osImage, err.Error())
	}
	imageHyperVGen := string(vmImage.HyperVGeneration)
	if !instanceHyperVGenSet.Has(imageHyperVGen) {
		errMsg := fmt.Sprintf("instance type supports HyperVGenerations %v but the specified image is for HyperVGeneration %s; to correct this issue either specify a compatible instance type or change the HyperVGeneration for the image by using a different SKU", instanceHyperVGenSet.UnsortedList(), imageHyperVGen)
		return field.Invalid(osImageFieldPath, osImage.SKU, errMsg)
	}

	// Image has license terms to be accepted
	osImagePlan := osImage.Plan
	if len(osImagePlan) == 0 {
		// Use the default if not set in the install-config
		osImagePlan = aztypes.ImageWithPurchasePlan
	}
	if plan := vmImage.Plan; plan != nil {
		if osImagePlan == aztypes.ImageNoPurchasePlan {
			return field.Invalid(osImageFieldPath, osImage, "marketplace image requires license terms to be accepted")
		}
		termsAccepted, err := client.AreMarketplaceImageTermsAccepted(context.Background(), osImage.Publisher, osImage.Offer, osImage.SKU)
		if err != nil {
			return field.Invalid(osImageFieldPath, osImage, fmt.Sprintf("could not determine if the license terms for the marketplace image have been accepted: %v", err))
		}
		if !termsAccepted {
			return field.Invalid(osImageFieldPath, osImage, "the license terms for the marketplace image have not been accepted")
		}
	} else if osImagePlan == aztypes.ImageWithPurchasePlan {
		return field.Invalid(osImageFieldPath, osImage, "image has no license terms. Set Plan to \"NoPurchasePlan\" to continue.")
	}

	return nil
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

func validateBootDiagnostics(client API, ic *types.InstallConfig) (allErrs field.ErrorList) {
	if ic.Azure.DefaultMachinePlatform != nil {
		bootDiag := ic.Azure.DefaultMachinePlatform.BootDiagnostics
		if err := checkBootDiagnosticsURI(client, bootDiag, ic.Platform.Azure.Region); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("platform", "azure", "defaultMachinePlatform",
				"bootDiagnostics"), bootDiag, err.Error()))
		}
	}

	if ic.ControlPlane != nil && ic.ControlPlane.Platform.Azure != nil {
		bootDiag := ic.ControlPlane.Platform.Azure.BootDiagnostics
		if err := checkBootDiagnosticsURI(client, bootDiag, ic.Platform.Azure.Region); err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("platform", "azure", "controlPlane",
				"bootDiagnostics"), bootDiag, err.Error()))
		}
	}

	if ic.Compute != nil {
		for inx, compute := range ic.Compute {
			if compute.Platform.Azure == nil {
				continue
			}
			bootDiag := compute.Platform.Azure.BootDiagnostics
			if err := checkBootDiagnosticsURI(client, bootDiag, ic.Platform.Azure.Region); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("platform", "azure", fmt.Sprintf("compute[%d]", inx),
					"bootDiagnostics"), bootDiag, err.Error()))
			}
		}
	}
	return allErrs
}

func checkBootDiagnosticsURI(client API, diag *aztypes.BootDiagnostics, region string) error {
	missingErrorMessage := "missing %s for user managed boot diagnostics"
	errorField := ""
	if diag != nil && diag.Type == capz.UserManagedDiagnosticsStorage {
		if diag.StorageAccountName != "" && diag.ResourceGroup != "" {
			return client.CheckIfExistsStorageAccount(context.TODO(), diag.ResourceGroup, diag.StorageAccountName, region)
		}
		if diag.ResourceGroup == "" {
			errorField += "resource group, "
		}
		if diag.StorageAccountName == "" {
			errorField += "storage account name, "
		}
		return fmt.Errorf(missingErrorMessage, errorField[:len(errorField)-2])
	}
	return nil
}
