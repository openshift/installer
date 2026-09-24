/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhooks

import (
	"fmt"
	"net"
	"reflect"
	"regexp"

	valid "github.com/asaskevich/govalidator/v11"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
)

const (
	// can't use: \/"'[]:|<>+=;,.?*@&, Can't start with underscore. Can't end with period or hyphen.
	// not using . in the name to avoid issues when the name is part of DNS name.
	clusterNameRegex = `^[a-z0-9][a-z0-9-]{0,42}[a-z0-9]$`
	// max length of 44 to allow for cluster name to be used as a prefix for VMs and other resources that
	// have limitations as outlined here https://learn.microsoft.com/azure/azure-resource-manager/management/resource-name-rules.
	clusterNameMaxLength = 44
	// obtained from https://learn.microsoft.com/rest/api/resources/resourcegroups/createorupdate#uri-parameters.
	resourceGroupRegex = `^[-\w\._\(\)]+$`
	// described in https://learn.microsoft.com/azure/azure-resource-manager/management/resource-name-rules.
	subnetRegex       = `^[-\w\._]+$`
	loadBalancerRegex = `^[-\w\._]+$`
	// MaxLoadBalancerOutboundIPs is the maximum number of outbound IPs in a Standard LoadBalancer frontend configuration.
	MaxLoadBalancerOutboundIPs = 16
	// MinLBIdleTimeoutInMinutes is the minimum number of minutes for the LB idle timeout.
	MinLBIdleTimeoutInMinutes = 4
	// MaxLBIdleTimeoutInMinutes is the maximum number of minutes for the LB idle timeout.
	MaxLBIdleTimeoutInMinutes = 30
	// Network security rules should be a number between 100 and 4096.
	// https://learn.microsoft.com/azure/virtual-network/network-security-groups-overview#security-rules
	minRulePriority = 100
	maxRulePriority = 4096
	// Must start with 'Microsoft.', then an alpha character, then can include alnum.
	serviceEndpointServiceRegexPattern = `^Microsoft\.[a-zA-Z]{1,42}[a-zA-Z0-9]{0,42}$`
	// Must start with an alpha character and then can include alnum OR be only *.
	serviceEndpointLocationRegexPattern = `^([a-z]{1,42}\d{0,5}|[*])$`
	// described in https://learn.microsoft.com/azure/azure-resource-manager/management/resource-name-rules.
	privateEndpointRegex = `^[-\w\._]+$`
	// resource ID Pattern.
	resourceIDPattern = `(?i)subscriptions/(.+)/resourceGroups/(.+)/providers/(.+?)/(.+?)/(.+)`
)

var (
	serviceEndpointServiceRegex  = regexp.MustCompile(serviceEndpointServiceRegexPattern)
	serviceEndpointLocationRegex = regexp.MustCompile(serviceEndpointLocationRegexPattern)
)

// validateAzureCluster validates a cluster.
func validateAzureCluster(c *infrav1.AzureCluster, old *infrav1.AzureCluster) (admission.Warnings, error) {
	var allErrs field.ErrorList
	allErrs = append(allErrs, validateAzureClusterName(c)...)
	allErrs = append(allErrs, validateAzureClusterSpec(c, old)...)
	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: infrav1.AzureClusterKind},
		c.Name, allErrs)
}

// validateAzureClusterSpec validates a ClusterSpec.
func validateAzureClusterSpec(c *infrav1.AzureCluster, old *infrav1.AzureCluster) field.ErrorList {
	var allErrs field.ErrorList
	var oldNetworkSpec infrav1.NetworkSpec
	if old != nil {
		oldNetworkSpec = old.Spec.NetworkSpec
	}

	allErrs = append(allErrs, validateNetworkSpec(c.Spec.ControlPlaneEnabled, c.Spec.NetworkSpec, oldNetworkSpec, field.NewPath("spec").Child("networkSpec"))...)

	var oldCloudProviderConfigOverrides *infrav1.CloudProviderConfigOverrides
	if old != nil {
		oldCloudProviderConfigOverrides = old.Spec.CloudProviderConfigOverrides
	}
	allErrs = append(allErrs, validateCloudProviderConfigOverrides(c.Spec.CloudProviderConfigOverrides, oldCloudProviderConfigOverrides,
		field.NewPath("spec").Child("cloudProviderConfigOverrides"))...)

	// If ClusterSpec has non-nil ExtendedLocation field but not enable EdgeZone feature gate flag, ClusterSpec validation failed.
	if !feature.Gates.Enabled(feature.EdgeZone) && c.Spec.ExtendedLocation != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "extendedLocation"), "can be set only if the EdgeZone feature flag is enabled"))
	}

	if err := validateBastionSpec(c.Spec.BastionSpec, field.NewPath("spec").Child("azureBastion").Child("bastionSpec")); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := validateIdentityRef(c.Spec.IdentityRef, field.NewPath("spec").Child("identityRef")); err != nil {
		allErrs = append(allErrs, err)
	}

	return allErrs
}

// validateAzureClusterName validates ClusterName.
func validateAzureClusterName(c *infrav1.AzureCluster) field.ErrorList {
	var allErrs field.ErrorList
	if len(c.Name) > clusterNameMaxLength {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata").Child("Name"), c.Name,
			fmt.Sprintf("Cluster Name longer than allowed length of %d characters", clusterNameMaxLength)))
	}
	if success, _ := regexp.MatchString(clusterNameRegex, c.Name); !success {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata").Child("name"), c.Name,
			fmt.Sprintf("Cluster Name doesn't match regex %s, can contain only lowercase alphanumeric characters and '-', must start/end with an alphanumeric character",
				clusterNameRegex)))
	}
	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

// validateBastionSpec validates a BastionSpec.
func validateBastionSpec(bastionSpec infrav1.BastionSpec, fldPath *field.Path) *field.Error {
	if bastionSpec.AzureBastion != nil && bastionSpec.AzureBastion.Sku != infrav1.StandardBastionHostSku && bastionSpec.AzureBastion.EnableTunneling {
		return field.Invalid(fldPath.Child("sku"), bastionSpec.AzureBastion.Sku,
			"sku must be Standard if tunneling is enabled")
	}
	return nil
}

// validateIdentityRef validates an IdentityRef.
func validateIdentityRef(identityRef *corev1.ObjectReference, fldPath *field.Path) *field.Error {
	if identityRef == nil {
		return field.Required(fldPath, "identityRef is required")
	}
	if identityRef.Kind != infrav1.AzureClusterIdentityKind {
		return field.NotSupported(fldPath.Child("name"), identityRef.Name, []string{"AzureClusterIdentity"})
	}
	return nil
}

// validateNetworkSpec validates a NetworkSpec.
func validateNetworkSpec(controlPlaneEnabled bool, networkSpec infrav1.NetworkSpec, old infrav1.NetworkSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	// If the user specifies a resourceGroup for vnet, it means
	// that they intend to use a pre-existing vnet. In this case,
	// we need to verify the information they provide
	if networkSpec.Vnet.ResourceGroup != "" {
		if err := validateResourceGroup(networkSpec.Vnet.ResourceGroup,
			fldPath.Child("vnet").Child("resourceGroup")); err != nil {
			allErrs = append(allErrs, err)
		}

		allErrs = append(allErrs, validateVnetCIDR(networkSpec.Vnet.CIDRBlocks, fldPath.Child("cidrBlocks"))...)

		allErrs = append(allErrs, validateSubnets(controlPlaneEnabled, networkSpec.Subnets, networkSpec.Vnet, fldPath.Child("subnets"))...)

		allErrs = append(allErrs, validateVnetPeerings(networkSpec.Vnet.Peerings, fldPath.Child("peerings"))...)
	}

	var cidrBlocks []string
	if controlPlaneEnabled {
		controlPlaneSubnet, err := networkSpec.GetControlPlaneSubnet()
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("subnets"), networkSpec.Subnets, "ControlPlaneSubnet invalid"))
		}

		cidrBlocks = controlPlaneSubnet.CIDRBlocks
		allErrs = append(allErrs, validateAPIServerLB(networkSpec.APIServerLB, old.APIServerLB, cidrBlocks, fldPath.Child("apiServerLB"))...)
	}

	var needOutboundLB bool
	for _, subnet := range networkSpec.Subnets {
		if (subnet.Role == infrav1.SubnetNode || subnet.Role == infrav1.SubnetCluster) && subnet.IsIPv6Enabled() {
			needOutboundLB = true
			break
		}
	}
	if needOutboundLB {
		allErrs = append(allErrs, validateNodeOutboundLB(networkSpec.NodeOutboundLB, old.NodeOutboundLB, networkSpec.APIServerLB, fldPath.Child("nodeOutboundLB"))...)
	}
	if controlPlaneEnabled {
		allErrs = append(allErrs, validateControlPlaneOutboundLB(networkSpec.ControlPlaneOutboundLB, networkSpec.APIServerLB, fldPath.Child("controlPlaneOutboundLB"))...)
	}
	var lbType = infrav1.Internal
	if networkSpec.APIServerLB != nil {
		lbType = networkSpec.APIServerLB.Type
	}
	allErrs = append(allErrs, validatePrivateDNSZoneName(networkSpec.PrivateDNSZoneName, controlPlaneEnabled, lbType, fldPath.Child("privateDNSZoneName"))...)
	allErrs = append(allErrs, validatePrivateDNSZoneResourceGroup(networkSpec.PrivateDNSZoneName, networkSpec.PrivateDNSZoneResourceGroup, fldPath.Child("privateDNSZoneResourceGroup"))...)

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

// validateResourceGroup validates a ResourceGroup.
func validateResourceGroup(resourceGroup string, fldPath *field.Path) *field.Error {
	if success, _ := regexp.MatchString(resourceGroupRegex, resourceGroup); !success {
		return field.Invalid(fldPath, resourceGroup,
			fmt.Sprintf("resourceGroup doesn't match regex %s", resourceGroupRegex))
	}
	return nil
}

// validateSubnets validates a list of Subnets.
// When configuring a cluster, it is essential to include either a control-plane subnet and a node subnet, or a user can configure a cluster subnet which will be used as a control-plane subnet and a node subnet.
func validateSubnets(controlPlaneEnabled bool, subnets infrav1.Subnets, vnet infrav1.VnetSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	subnetNames := make(map[string]bool, len(subnets))
	requiredSubnetRoles := map[string]bool{
		"node": false,
	}
	if controlPlaneEnabled {
		requiredSubnetRoles["control-plane"] = false
	}
	clusterSubnet := false
	for i, subnet := range subnets {
		if err := validateSubnetName(subnet.Name, fldPath.Index(i).Child("name")); err != nil {
			allErrs = append(allErrs, err)
		}
		if _, ok := subnetNames[subnet.Name]; ok {
			allErrs = append(allErrs, field.Duplicate(fldPath, subnet.Name))
		}
		subnetNames[subnet.Name] = true
		if subnet.Role == infrav1.SubnetCluster {
			clusterSubnet = true
		} else {
			for role := range requiredSubnetRoles {
				if role == string(subnet.Role) {
					requiredSubnetRoles[role] = true
				}
			}
		}

		for j, rule := range subnet.SecurityGroup.SecurityRules {
			if err := validateSecurityRule(
				rule,
				fldPath.Index(i).Child("securityGroup").Child("securityRules").Index(j),
			); err != nil {
				allErrs = append(allErrs, err...)
			}
		}
		allErrs = append(allErrs, validateSubnetCIDR(subnet.CIDRBlocks, vnet.CIDRBlocks, fldPath.Index(i).Child("cidrBlocks"))...)

		if len(subnet.ServiceEndpoints) > 0 {
			allErrs = append(allErrs, validateServiceEndpoints(subnet.ServiceEndpoints, fldPath.Index(i).Child("serviceEndpoints"))...)
		}

		if len(subnet.PrivateEndpoints) > 0 {
			allErrs = append(allErrs, validatePrivateEndpoints(subnet.PrivateEndpoints, subnet.CIDRBlocks, fldPath.Index(i).Child("privateEndpoints"))...)
		}
	}

	// The clusterSubnet is applicable to both the control-plane and node pools.
	// Validation of requiredSubnetRoles is skipped since clusterSubnet is set to true.
	if clusterSubnet {
		return allErrs
	}

	for k, v := range requiredSubnetRoles {
		if !v {
			allErrs = append(allErrs, field.Required(fldPath,
				fmt.Sprintf("required role %s not included in provided subnets", k)))
		}
	}
	return allErrs
}

// validateSubnetName validates the Name of a Subnet.
func validateSubnetName(name string, fldPath *field.Path) *field.Error {
	if success, _ := regexp.Match(subnetRegex, []byte(name)); !success {
		return field.Invalid(fldPath, name,
			fmt.Sprintf("name of subnet doesn't match regex %s", subnetRegex))
	}
	return nil
}

// validateSubnetCIDR validates the CIDR blocks of a Subnet.
func validateSubnetCIDR(subnetCidrBlocks []string, vnetCidrBlocks []string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	var vnetNws []*net.IPNet

	for _, vnetCidr := range vnetCidrBlocks {
		if _, vnetNw, err := net.ParseCIDR(vnetCidr); err == nil {
			vnetNws = append(vnetNws, vnetNw)
		}
	}

	for _, subnetCidr := range subnetCidrBlocks {
		subnetCidrIP, _, err := net.ParseCIDR(subnetCidr)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, subnetCidr, "invalid CIDR format"))
		}

		var found bool
		for _, vnetNw := range vnetNws {
			if vnetNw.Contains(subnetCidrIP) {
				found = true
				break
			}
		}

		if !found {
			allErrs = append(allErrs, field.Invalid(fldPath, subnetCidr, fmt.Sprintf("subnet CIDR not in vnet address space: %s", vnetCidrBlocks)))
		}
	}

	return allErrs
}

// validateVnetCIDR validates the CIDR blocks of a Vnet.
func validateVnetCIDR(vnetCIDRBlocks []string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	for _, vnetCidr := range vnetCIDRBlocks {
		if _, _, err := net.ParseCIDR(vnetCidr); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath, vnetCidr, "invalid CIDR format"))
		}
	}
	return allErrs
}

// validateVnetPeerings validates a list of virtual network peerings.
func validateVnetPeerings(peerings infrav1.VnetPeerings, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	vnetIdentifiers := make(map[string]bool, len(peerings))

	for _, peering := range peerings {
		vnetIdentifier := peering.ResourceGroup + "/" + peering.RemoteVnetName
		if _, ok := vnetIdentifiers[vnetIdentifier]; ok {
			allErrs = append(allErrs, field.Duplicate(fldPath, vnetIdentifier))
		}
		vnetIdentifiers[vnetIdentifier] = true
	}
	return allErrs
}

// validateLoadBalancerName validates the Name of a Load Balancer.
func validateLoadBalancerName(name string, fldPath *field.Path) *field.Error {
	if success, _ := regexp.Match(loadBalancerRegex, []byte(name)); !success {
		return field.Invalid(fldPath, name,
			fmt.Sprintf("name of load balancer doesn't match regex %s", loadBalancerRegex))
	}
	return nil
}

// validateInternalLBIPAddress validates a InternalLBIPAddress.
func validateInternalLBIPAddress(address string, cidrs []string, fldPath *field.Path) *field.Error {
	ip := net.ParseIP(address)
	if ip == nil {
		return field.Invalid(fldPath, address,
			"Internal LB IP address isn't a valid IPv4 or IPv6 address")
	}
	for _, cidr := range cidrs {
		_, subnet, _ := net.ParseCIDR(cidr)
		if subnet.Contains(ip) {
			return nil
		}
	}
	return field.Invalid(fldPath, address,
		fmt.Sprintf("Internal LB IP address needs to be in control plane subnet range (%s)", cidrs))
}

// validateSecurityRule validates a SecurityRule.
func validateSecurityRule(rule infrav1.SecurityRule, fldPath *field.Path) (allErrs field.ErrorList) {
	if rule.Priority < minRulePriority || rule.Priority > maxRulePriority {
		allErrs = append(allErrs, field.Invalid(fldPath, rule.Priority, fmt.Sprintf("security rule priorities should be between %d and %d", minRulePriority, maxRulePriority)))
	}

	if rule.Source != nil && rule.Sources != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, rule.Source, "security rule cannot have both source and sources"))
	}

	return allErrs
}

func validateAPIServerLB(lb *infrav1.LoadBalancerSpec, old *infrav1.LoadBalancerSpec, cidrs []string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	lbClassSpec := lb.LoadBalancerClassSpec
	var olLBClassSpec infrav1.LoadBalancerClassSpec
	if old != nil {
		olLBClassSpec = old.LoadBalancerClassSpec
	}

	allErrs = append(allErrs, validateClassSpecForAPIServerLB(lbClassSpec, &olLBClassSpec, fldPath)...)

	// Name should be valid.
	if err := validateLoadBalancerName(lb.Name, fldPath.Child("name")); err != nil {
		allErrs = append(allErrs, err)
	}
	// Name should be immutable.
	if old != nil && old.Name != "" && old.Name != lb.Name {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("name"), "API Server load balancer name should not be modified after AzureCluster creation."))
	}

	publicIPCount, privateIPCount := 0, 0
	privateIP := ""
	for i := range lb.FrontendIPs {
		if lb.FrontendIPs[i].PublicIP != nil {
			publicIPCount++
		}
		if lb.FrontendIPs[i].PrivateIPAddress != "" {
			privateIPCount++
			privateIP = lb.FrontendIPs[i].PrivateIPAddress
		}
	}
	if lb.Type == infrav1.Public {
		// there should be one public IP for public LB.
		if publicIPCount != 1 || ptr.Deref[int32](lb.FrontendIPsCount, 1) != 1 {
			// Note: FrontendIPsCount creates public IPs when set. Therefore, we check for both publicIPCount and FrontendIPsCount to be 1.
			allErrs = append(allErrs, field.Invalid(fldPath.Child("frontendIPConfigs"), lb.FrontendIPs,
				"API Server Load balancer should have 1 Frontend IP"))
		}
		if feature.Gates.Enabled(feature.APIServerILB) {
			if err := validateInternalLBIPAddress(privateIP, cidrs, fldPath.Child("frontendIPConfigs").Index(0).Child("privateIP")); err != nil {
				allErrs = append(allErrs, err)
			}
		} else {
			// API Server LB should not have a Private IP if APIServerILB feature is disabled.
			if privateIPCount > 0 {
				allErrs = append(allErrs, field.Forbidden(fldPath.Child("frontendIPConfigs").Index(0).Child("privateIP"),
					"Public Load Balancers cannot have a Private IP"))
			}
		}
	}

	// internal LB should not have a public IP.
	if lb.Type == infrav1.Internal {
		if publicIPCount != 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("frontendIPConfigs").Index(0).Child("publicIP"),
				"Internal Load Balancers cannot have a Public IP"))
		}
		if privateIPCount != 1 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("frontendIPConfigs"), lb.FrontendIPs,
				"API Server Load balancer of type private should have 1 frontend private IP"))
		} else {
			if err := validateInternalLBIPAddress(lb.FrontendIPs[0].PrivateIPAddress, cidrs,
				fldPath.Child("frontendIPConfigs").Index(0).Child("privateIP")); err != nil {
				allErrs = append(allErrs, err)
			}

			if old != nil && len(old.FrontendIPs) != 0 && old.FrontendIPs[0].PrivateIPAddress != lb.FrontendIPs[0].PrivateIPAddress {
				allErrs = append(allErrs, field.Forbidden(fldPath.Child("name"), "API Server load balancer private IP should not be modified after AzureCluster creation."))
			}
		}
	}
	return allErrs
}

func validateNodeOutboundLB(lb *infrav1.LoadBalancerSpec, old *infrav1.LoadBalancerSpec, apiserverLB *infrav1.LoadBalancerSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	var lbClassSpec, oldClassSpec *infrav1.LoadBalancerClassSpec
	if lb != nil {
		lbClassSpec = &lb.LoadBalancerClassSpec
	}
	if old != nil {
		oldClassSpec = &old.LoadBalancerClassSpec
	}
	apiserverLBClassSpec := apiserverLB.LoadBalancerClassSpec

	allErrs = append(allErrs, validateClassSpecForNodeOutboundLB(lbClassSpec, oldClassSpec, apiserverLBClassSpec, fldPath)...)

	if lb == nil {
		return allErrs
	}

	if old != nil && old.ID != lb.ID {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("id"), "Node outbound load balancer ID should not be modified after AzureCluster creation."))
	}

	if old != nil && old.Name != lb.Name {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("name"), "Node outbound load balancer Name should not be modified after AzureCluster creation."))
	}

	if old != nil && old.FrontendIPsCount == lb.FrontendIPsCount {
		if len(old.FrontendIPs) != len(lb.FrontendIPs) {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("frontendIPs"), "Node outbound load balancer FrontendIPs cannot be modified after AzureCluster creation."))
		}

		if len(old.FrontendIPs) == len(lb.FrontendIPs) {
			for i, frontEndIP := range lb.FrontendIPs {
				oldFrontendIP := old.FrontendIPs[i]
				if oldFrontendIP.Name != frontEndIP.Name || !reflect.DeepEqual(*oldFrontendIP.PublicIP, *frontEndIP.PublicIP) {
					allErrs = append(allErrs, field.Forbidden(fldPath.Child("frontendIPs").Index(i),
						"Node outbound load balancer FrontendIPs cannot be modified after AzureCluster creation."))
				}
			}
		}
	}

	if lb.FrontendIPsCount != nil && *lb.FrontendIPsCount > MaxLoadBalancerOutboundIPs {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("frontendIPsCount"), *lb.FrontendIPsCount,
			fmt.Sprintf("Max front end ips allowed is %d", MaxLoadBalancerOutboundIPs)))
	}

	return allErrs
}

func validateControlPlaneOutboundLB(lb *infrav1.LoadBalancerSpec, apiserverLB *infrav1.LoadBalancerSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	var lbClassSpec *infrav1.LoadBalancerClassSpec
	if lb != nil {
		lbClassSpec = &lb.LoadBalancerClassSpec
	}
	apiServerLBClassSpec := apiserverLB.LoadBalancerClassSpec

	allErrs = append(allErrs, validateClassSpecForControlPlaneOutboundLB(lbClassSpec, apiServerLBClassSpec, fldPath)...)

	if apiServerLBClassSpec.Type == infrav1.Internal && lb != nil {
		if lb.FrontendIPsCount != nil && *lb.FrontendIPsCount > MaxLoadBalancerOutboundIPs {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("frontendIPsCount"), *lb.FrontendIPsCount,
				fmt.Sprintf("Max front end ips allowed is %d", MaxLoadBalancerOutboundIPs)))
		}
	}

	return allErrs
}

// validatePrivateDNSZoneName validates the PrivateDNSZoneName.
func validatePrivateDNSZoneName(privateDNSZoneName string, controlPlaneEnabled bool, apiserverLBType infrav1.LBType, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if privateDNSZoneName != "" {
		if controlPlaneEnabled && apiserverLBType != infrav1.Internal {
			allErrs = append(allErrs, field.Invalid(fldPath, apiserverLBType,
				"PrivateDNSZoneName is available only if APIServerLB.Type is Internal"))
		}
		if !valid.IsDNSName(privateDNSZoneName) {
			allErrs = append(allErrs, field.Invalid(fldPath, privateDNSZoneName,
				"PrivateDNSZoneName can only contain alphanumeric characters, underscores and dashes, must end with an alphanumeric character",
			))
		}
	}

	return allErrs
}

// validatePrivateDNSZoneResourceGroup validates the PrivateDNSZoneResourceGroup.
// A private DNS Zone's resource group is valid as long as privateDNSZoneName is provided with the private dns resource group name.
func validatePrivateDNSZoneResourceGroup(privateDNSZoneName string, privateDNSZoneResourceGroup string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if privateDNSZoneResourceGroup != "" {
		if privateDNSZoneName == "" {
			allErrs = append(allErrs, field.Invalid(fldPath, privateDNSZoneName,
				"PrivateDNSZoneResourceGroup can only be used when PrivateDNSZoneName is provided"))
		}
		if err := validateResourceGroup(privateDNSZoneResourceGroup, fldPath); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	return allErrs
}

// validateCloudProviderConfigOverrides validates CloudProviderConfigOverrides.
func validateCloudProviderConfigOverrides(oldConfig, newConfig *infrav1.CloudProviderConfigOverrides, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if !reflect.DeepEqual(oldConfig, newConfig) {
		allErrs = append(allErrs, field.Invalid(fldPath, newConfig, "cannot change cloudProviderConfigOverrides cluster creation"))
	}
	return allErrs
}

func validateClassSpecForAPIServerLB(lb infrav1.LoadBalancerClassSpec, old *infrav1.LoadBalancerClassSpec, apiServerLBPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// SKU should be Standard
	if lb.SKU != infrav1.SKUStandard {
		allErrs = append(allErrs, field.NotSupported(apiServerLBPath.Child("sku"), lb.SKU, []string{string(infrav1.SKUStandard)}))
	}

	// Type should be Public or Internal.
	if lb.Type != infrav1.Internal && lb.Type != infrav1.Public {
		allErrs = append(allErrs, field.NotSupported(apiServerLBPath.Child("type"), lb.Type,
			[]string{string(infrav1.Public), string(infrav1.Internal)}))
	}

	// SKU should be immutable.
	if old != nil && old.SKU != "" && old.SKU != lb.SKU {
		allErrs = append(allErrs, field.Forbidden(apiServerLBPath.Child("sku"), "API Server load balancer SKU should not be modified after AzureCluster creation."))
	}

	// Type should be immutable.
	if old != nil && old.Type != "" && old.Type != lb.Type {
		allErrs = append(allErrs, field.Forbidden(apiServerLBPath.Child("type"), "API Server load balancer type should not be modified after AzureCluster creation."))
	}

	// IdletimeoutInMinutes should be immutable.
	if old != nil && old.IdleTimeoutInMinutes != nil && !ptr.Equal(old.IdleTimeoutInMinutes, lb.IdleTimeoutInMinutes) {
		allErrs = append(allErrs, field.Forbidden(apiServerLBPath.Child("idleTimeoutInMinutes"), "API Server load balancer idle timeout cannot be modified after AzureCluster creation."))
	}

	if lb.IdleTimeoutInMinutes != nil && (*lb.IdleTimeoutInMinutes < MinLBIdleTimeoutInMinutes || *lb.IdleTimeoutInMinutes > MaxLBIdleTimeoutInMinutes) {
		allErrs = append(allErrs, field.Invalid(apiServerLBPath.Child("idleTimeoutInMinutes"), *lb.IdleTimeoutInMinutes,
			fmt.Sprintf("API Server load balancer idle timeout should be between %d and %d minutes", MinLBIdleTimeoutInMinutes, MaxLBIdleTimeoutInMinutes)))
	}

	return allErrs
}

func validateClassSpecForNodeOutboundLB(lb *infrav1.LoadBalancerClassSpec, old *infrav1.LoadBalancerClassSpec, apiserverLB infrav1.LoadBalancerClassSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// LB can be nil when disabled for private clusters.
	if lb == nil && apiserverLB.Type == infrav1.Internal {
		return allErrs
	}

	if lb == nil {
		allErrs = append(allErrs, field.Required(fldPath, "Node outbound load balancer cannot be nil for public clusters."))
		return allErrs
	}

	if old != nil && old.SKU != lb.SKU {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("sku"), "Node outbound load balancer SKU should not be modified after AzureCluster creation."))
	}

	if old != nil && old.Type != lb.Type {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("type"), "Node outbound load balancer Type cannot be modified after AzureCluster creation."))
	}

	if old != nil && !ptr.Equal(old.IdleTimeoutInMinutes, lb.IdleTimeoutInMinutes) {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("idleTimeoutInMinutes"), "Node outbound load balancer idle timeout cannot be modified after AzureCluster creation."))
	}

	if lb.IdleTimeoutInMinutes != nil && (*lb.IdleTimeoutInMinutes < MinLBIdleTimeoutInMinutes || *lb.IdleTimeoutInMinutes > MaxLBIdleTimeoutInMinutes) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("idleTimeoutInMinutes"), *lb.IdleTimeoutInMinutes,
			fmt.Sprintf("Node outbound idle timeout should be between %d and %d minutes", MinLBIdleTimeoutInMinutes, MaxLBIdleTimeoutInMinutes)))
	}

	return allErrs
}

func validateClassSpecForControlPlaneOutboundLB(lb *infrav1.LoadBalancerClassSpec, apiserverLB infrav1.LoadBalancerClassSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	switch apiserverLB.Type {
	case infrav1.Public:
		if lb != nil {
			allErrs = append(allErrs, field.Forbidden(fldPath, "Control plane outbound load balancer cannot be set for public clusters."))
		}
	case infrav1.Internal:
		// Control plane outbound lb can be nil when it's disabled for private clusters.
		if lb == nil {
			return nil
		}

		if lb.IdleTimeoutInMinutes != nil && (*lb.IdleTimeoutInMinutes < MinLBIdleTimeoutInMinutes || *lb.IdleTimeoutInMinutes > MaxLBIdleTimeoutInMinutes) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("idleTimeoutInMinutes"), *lb.IdleTimeoutInMinutes,
				fmt.Sprintf("Control plane outbound idle timeout should be between %d and %d minutes", MinLBIdleTimeoutInMinutes, MaxLBIdleTimeoutInMinutes)))
		}
	}

	return allErrs
}

func validateServiceEndpoints(serviceEndpoints []infrav1.ServiceEndpointSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	serviceEndpointsServices := make(map[string]bool, len(serviceEndpoints))
	for i, se := range serviceEndpoints {
		if se.Service == "" {
			allErrs = append(allErrs, field.Required(fldPath.Index(i).Child("service"), "service is required for all service endpoints"))
		} else {
			if err := validateServiceEndpointServiceName(se.Service, fldPath.Index(i).Child("service")); err != nil {
				allErrs = append(allErrs, err)
			}
			if _, ok := serviceEndpointsServices[se.Service]; ok {
				allErrs = append(allErrs, field.Duplicate(fldPath.Index(i).Child("service"), se.Service))
			}
			serviceEndpointsServices[se.Service] = true
		}

		if len(se.Locations) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Index(i).Child("locations"), "locations are required for all service endpoints"))
		} else {
			serviceEndpointsLocations := make(map[string]bool, len(se.Locations))
			for j, locationName := range se.Locations {
				if err := validateServiceEndpointLocationName(locationName, fldPath.Index(i).Child("locations").Index(j)); err != nil {
					allErrs = append(allErrs, err)
				}
				if _, ok := serviceEndpointsLocations[locationName]; ok {
					allErrs = append(allErrs, field.Duplicate(fldPath.Index(i).Child("locations").Index(j), locationName))
				}
				serviceEndpointsLocations[locationName] = true
			}
		}
	}

	return allErrs
}

func validateServiceEndpointServiceName(serviceName string, fldPath *field.Path) *field.Error {
	if success := serviceEndpointServiceRegex.MatchString(serviceName); !success {
		return field.Invalid(fldPath, serviceName, fmt.Sprintf("service name of endpoint service doesn't match regex %s", serviceEndpointServiceRegexPattern))
	}
	return nil
}

func validateServiceEndpointLocationName(location string, fldPath *field.Path) *field.Error {
	if success := serviceEndpointLocationRegex.MatchString(location); !success {
		return field.Invalid(fldPath, location, fmt.Sprintf("location doesn't match regex %s", serviceEndpointLocationRegexPattern))
	}
	return nil
}

func validatePrivateEndpoints(privateEndpointSpecs []infrav1.PrivateEndpointSpec, subnetCIDRs []string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	for i, pe := range privateEndpointSpecs {
		if err := validatePrivateEndpointName(pe.Name, fldPath.Index(i).Child("name")); err != nil {
			allErrs = append(allErrs, err)
		}

		if len(pe.PrivateLinkServiceConnections) == 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i), pe.PrivateLinkServiceConnections, "privateLinkServiceConnections cannot be empty"))
		}

		for j, privateLinkServiceConnection := range pe.PrivateLinkServiceConnections {
			if privateLinkServiceConnection.PrivateLinkServiceID == "" {
				allErrs = append(allErrs, field.Required(fldPath.Index(i).Child("privateLinkServiceConnections").Index(j), "privateLinkServiceID is required for all privateLinkServiceConnections in private endpoints"))
			} else {
				if err := validatePrivateEndpointPrivateLinkServiceConnection(privateLinkServiceConnection, fldPath.Index(i).Child("privateLinkServiceConnections").Index(j)); err != nil {
					allErrs = append(allErrs, err)
				}
			}
		}

		for _, privateIP := range pe.PrivateIPAddresses {
			if err := validatePrivateEndpointIPAddress(privateIP, subnetCIDRs, fldPath.Index(i).Child("privateIPAddresses")); err != nil {
				allErrs = append(allErrs, err)
			}
		}
	}

	return allErrs
}

// validatePrivateEndpointName validates the Name of a Private Endpoint.
func validatePrivateEndpointName(name string, fldPath *field.Path) *field.Error {
	if name == "" {
		return field.Invalid(fldPath, name, "name of private endpoint cannot be empty")
	}

	if success, _ := regexp.MatchString(privateEndpointRegex, name); !success {
		return field.Invalid(fldPath, name,
			fmt.Sprintf("name of private endpoint doesn't match regex %s", privateEndpointRegex))
	}
	return nil
}

// validatePrivateEndpointServiceID validates the service ID of a Private Endpoint.
func validatePrivateEndpointPrivateLinkServiceConnection(privateLinkServiceConnection infrav1.PrivateLinkServiceConnection, fldPath *field.Path) *field.Error {
	if success, _ := regexp.MatchString(resourceIDPattern, privateLinkServiceConnection.PrivateLinkServiceID); !success {
		return field.Invalid(fldPath, privateLinkServiceConnection.PrivateLinkServiceID,
			fmt.Sprintf("private endpoint privateLinkServiceConnection service ID doesn't match regex %s", resourceIDPattern))
	}
	if privateLinkServiceConnection.Name != "" {
		if success, _ := regexp.MatchString(privateEndpointRegex, privateLinkServiceConnection.Name); !success {
			return field.Invalid(fldPath, privateLinkServiceConnection.Name,
				fmt.Sprintf("private endpoint privateLinkServiceConnection name doesn't match regex %s", privateEndpointRegex))
		}
	}
	return nil
}

// validatePrivateEndpointIPAddress validates a Private Endpoint IP Address.
func validatePrivateEndpointIPAddress(address string, cidrs []string, fldPath *field.Path) *field.Error {
	ip := net.ParseIP(address)
	if ip == nil {
		return field.Invalid(fldPath, address,
			"Private Endpoint IP address isn't a valid IPv4 or IPv6 address")
	}

	for _, cidr := range cidrs {
		_, subnet, _ := net.ParseCIDR(cidr)
		if subnet != nil && subnet.Contains(ip) {
			return nil
		}
	}

	return field.Invalid(fldPath, address,
		fmt.Sprintf("Private Endpoint IP address needs to be in subnet range (%s)", cidrs))
}

// validateAzureClusterSubnetUpdate validates a ClusterSpec.NetworkSpec.Subnets for immutability.
func validateAzureClusterSubnetUpdate(c *infrav1.AzureCluster, old *infrav1.AzureCluster) field.ErrorList {
	var allErrs field.ErrorList

	oldSubnetMap := make(map[string]infrav1.SubnetSpec, len(old.Spec.NetworkSpec.Subnets))
	oldSubnetIndex := make(map[string]int, len(old.Spec.NetworkSpec.Subnets))
	for i, subnet := range old.Spec.NetworkSpec.Subnets {
		oldSubnetMap[subnet.Name] = subnet
		oldSubnetIndex[subnet.Name] = i
	}
	for i, subnet := range c.Spec.NetworkSpec.Subnets {
		if oldSubnet, ok := oldSubnetMap[subnet.Name]; ok {
			// Verify the CIDR blocks haven't changed for an owned Vnet.
			// A non-owned Vnet's CIDR block can change based on what's
			// defined in the spec vs what's been loaded from Azure directly.
			// This technically allows the cidr block to be modified in the brief
			// moments before the Vnet is created (because the tags haven't been
			// set yet) but once the Vnet has been created it becomes immutable.
			if old.Spec.NetworkSpec.Vnet.Tags.HasOwned(old.Name) && !reflect.DeepEqual(subnet.CIDRBlocks, oldSubnet.CIDRBlocks) {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "networkSpec", "subnets").Index(oldSubnetIndex[subnet.Name]).Child("CIDRBlocks"),
						c.Spec.NetworkSpec.Subnets[i].CIDRBlocks, "field is immutable"),
				)
			}
			// RouteTable.Name is immutable once set, but allow defaulting it from
			// empty to a generated name so existing clusters (whose control plane
			// subnet had no route table) can adopt the shared node route table on
			// upgrade. This mirrors the NatGateway.Name handling below.
			if (subnet.RouteTable.Name != oldSubnet.RouteTable.Name) && (oldSubnet.RouteTable.Name != "") {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "networkSpec", "subnets").Index(oldSubnetIndex[subnet.Name]).Child("RouteTable").Child("Name"),
						c.Spec.NetworkSpec.Subnets[i].RouteTable.Name, "field is immutable"),
				)
			}
			if (subnet.NatGateway.Name != oldSubnet.NatGateway.Name) && (oldSubnet.NatGateway.Name != "") {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "networkSpec", "subnets").Index(oldSubnetIndex[subnet.Name]).Child("NatGateway").Child("Name"),
						c.Spec.NetworkSpec.Subnets[i].NatGateway.Name, "field is immutable"),
				)
			}
			if subnet.SecurityGroup.Name != oldSubnet.SecurityGroup.Name {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "networkSpec", "subnets").Index(oldSubnetIndex[subnet.Name]).Child("SecurityGroup").Child("Name"),
						c.Spec.NetworkSpec.Subnets[i].SecurityGroup.Name, "field is immutable"),
				)
			}
		}
	}

	return allErrs
}
