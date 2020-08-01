package azure

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	azdns "github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	aznetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

// Validate executes platform-specific validation.
func Validate(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateNetworks(client, ic.Azure, ic.Networking.MachineNetwork, field.NewPath("platform").Child("azure"))...)
	allErrs = append(allErrs, validateRegion(client, field.NewPath("platform").Child("azure").Child("region"), ic.Azure)...)
	return allErrs.ToAggregate()
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

// validateRegionForUltraDisks checks that the Ultra SSD disks are available for the user.
func validateRegionForUltraDisks(fldPath *field.Path, client *Client, region string) *field.Error {
	diskType := "UltraSSD_LRS"

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	regionDisks, err := client.GetDiskSkus(ctx, region)
	if err != nil {
		return field.InternalError(fldPath.Child("diskType"), err)
	}

	for _, page := range regionDisks {
		for _, location := range to.StringSlice(page.Locations) {
			if location == diskType {
				return nil
			}
		}
	}

	return field.NotFound(fldPath.Child("diskType"), fmt.Sprintf("%s is not available in specified subscription for region %s", diskType, region))
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
	return allErrs
}
