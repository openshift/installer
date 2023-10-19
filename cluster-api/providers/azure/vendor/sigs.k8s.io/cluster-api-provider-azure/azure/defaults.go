/*
Copyright 2019 The Kubernetes Authors.

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

package azure

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	"sigs.k8s.io/cluster-api-provider-azure/version"
)

const (
	// DefaultUserName is the default username for a created VM.
	DefaultUserName = "capi"
	// DefaultAKSUserName is the default username for a created AKS VM.
	DefaultAKSUserName = "azureuser"
	// PublicCloudName is the name of the Azure public cloud.
	PublicCloudName = "AzurePublicCloud"
	// ChinaCloudName is the name of the Azure China cloud.
	ChinaCloudName = "AzureChinaCloud"
	// USGovernmentCloudName is the name of the Azure US Government cloud.
	USGovernmentCloudName = "AzureUSGovernmentCloud"
)

const (
	// DefaultImageOfferID is the default Azure Marketplace offer ID.
	DefaultImageOfferID = "capi"
	// DefaultWindowsImageOfferID is the default Azure Marketplace offer ID for Windows.
	DefaultWindowsImageOfferID = "capi-windows"
	// DefaultImagePublisherID is the default Azure Marketplace publisher ID.
	DefaultImagePublisherID = "cncf-upstream"
	// LatestVersion is the image version latest.
	LatestVersion = "latest"
)

const (
	// LinuxOS is Linux OS value for OSDisk.OSType.
	LinuxOS = "Linux"
	// WindowsOS is Windows OS value for OSDisk.OSType.
	WindowsOS = "Windows"
)

const (
	// BootstrappingExtensionLinux is the name of the Linux CAPZ bootstrapping VM extension.
	BootstrappingExtensionLinux = "CAPZ.Linux.Bootstrapping"
	// BootstrappingExtensionWindows is the name of the Windows CAPZ bootstrapping VM extension.
	BootstrappingExtensionWindows = "CAPZ.Windows.Bootstrapping"
)

const (
	// DefaultWindowsOsAndVersion is the default Windows Server version to use when
	// generating default images for Windows nodes.
	DefaultWindowsOsAndVersion = "windows-2019"
)

const (
	// Global is the Azure global location value.
	Global = "global"
)

const (
	// PrivateAPIServerHostname will be used as the api server hostname for private clusters.
	PrivateAPIServerHostname = "apiserver"
)

const (
	// ControlPlaneNodeGroup will be used to create availability set for control plane machines.
	ControlPlaneNodeGroup = "control-plane"
)

const (
	// bootstrapExtensionRetries is the number of retries in the BootstrapExtensionCommand.
	// NOTE: the overall timeout will be number of retries * retry sleep, in this case 60 * 5s = 300s.
	bootstrapExtensionRetries = 60
	// bootstrapExtensionSleep is the duration in seconds to sleep before each retry in the BootstrapExtensionCommand.
	bootstrapExtensionSleep = 5
	// bootstrapSentinelFile is the file written by bootstrap provider on machines to indicate successful bootstrapping,
	// as defined by the Cluster API Bootstrap Provider contract (https://cluster-api.sigs.k8s.io/developer/providers/bootstrap.html).
	bootstrapSentinelFile = "/run/cluster-api/bootstrap-success.complete"
)

const (
	// CustomHeaderPrefix is the prefix of annotations that enable additional cluster / node pool features.
	// Whatever follows the prefix will be passed as a header to cluster/node pool creation/update requests.
	// E.g. add `"infrastructure.cluster.x-k8s.io/custom-header-UseGPUDedicatedVHD": "true"` annotation to
	// AzureManagedMachinePool CR to enable creating GPU nodes by the node pool.
	CustomHeaderPrefix = "infrastructure.cluster.x-k8s.io/custom-header-"
)

var (
	// LinuxBootstrapExtensionCommand is the command the VM bootstrap extension will execute to verify Linux nodes bootstrap completes successfully.
	LinuxBootstrapExtensionCommand = fmt.Sprintf("for i in $(seq 1 %d); do test -f %s && break; if [ $i -eq %d ]; then exit 1; else sleep %d; fi; done", bootstrapExtensionRetries, bootstrapSentinelFile, bootstrapExtensionRetries, bootstrapExtensionSleep)
	// WindowsBootstrapExtensionCommand is the command the VM bootstrap extension will execute to verify Windows nodes bootstrap completes successfully.
	WindowsBootstrapExtensionCommand = fmt.Sprintf("powershell.exe -Command \"for ($i = 0; $i -lt %d; $i++) {if (Test-Path '%s') {exit 0} else {Start-Sleep -Seconds %d}} exit -2\"",
		bootstrapExtensionRetries, bootstrapSentinelFile, bootstrapExtensionSleep)
)

// GenerateBackendAddressPoolName generates a load balancer backend address pool name.
func GenerateBackendAddressPoolName(lbName string) string {
	return fmt.Sprintf("%s-%s", lbName, "backendPool")
}

// GenerateOutboundBackendAddressPoolName generates a load balancer outbound backend address pool name.
func GenerateOutboundBackendAddressPoolName(lbName string) string {
	return fmt.Sprintf("%s-%s", lbName, "outboundBackendPool")
}

// GenerateFrontendIPConfigName generates a load balancer frontend IP config name.
func GenerateFrontendIPConfigName(lbName string) string {
	return fmt.Sprintf("%s-%s", lbName, "frontEnd")
}

// GenerateNodeOutboundIPName generates a public IP name, based on the cluster name.
func GenerateNodeOutboundIPName(clusterName string) string {
	return fmt.Sprintf("pip-%s-node-outbound", clusterName)
}

// GenerateNodePublicIPName generates a node public IP name, based on the machine name.
func GenerateNodePublicIPName(machineName string) string {
	return fmt.Sprintf("pip-%s", machineName)
}

// GenerateControlPlaneOutboundLBName generates the name of the control plane outbound LB.
func GenerateControlPlaneOutboundLBName(clusterName string) string {
	return fmt.Sprintf("%s-outbound-lb", clusterName)
}

// GenerateControlPlaneOutboundIPName generates a public IP name, based on the cluster name.
func GenerateControlPlaneOutboundIPName(clusterName string) string {
	return fmt.Sprintf("pip-%s-controlplane-outbound", clusterName)
}

// GeneratePrivateDNSZoneName generates the name of a private DNS zone based on the cluster name.
func GeneratePrivateDNSZoneName(clusterName string) string {
	return fmt.Sprintf("%s.capz.io", clusterName)
}

// GeneratePrivateFQDN generates the FQDN for a private API Server based on the private DNS zone name.
func GeneratePrivateFQDN(zoneName string) string {
	return fmt.Sprintf("%s.%s", PrivateAPIServerHostname, zoneName)
}

// GenerateVNetLinkName generates the name of a virtual network link name based on the vnet name.
func GenerateVNetLinkName(vnetName string) string {
	return fmt.Sprintf("%s-link", vnetName)
}

// GenerateNICName generates the name of a network interface based on the name of a VM.
func GenerateNICName(machineName string, multiNIC bool, index int) string {
	if multiNIC {
		return fmt.Sprintf("%s-nic-%d", machineName, index)
	}
	return fmt.Sprintf("%s-nic", machineName)
}

// GeneratePublicNICName generates the name of a public network interface based on the name of a VM.
func GeneratePublicNICName(machineName string) string {
	return fmt.Sprintf("%s-public-nic", machineName)
}

// GenerateOSDiskName generates the name of an OS disk based on the name of a VM.
func GenerateOSDiskName(machineName string) string {
	return fmt.Sprintf("%s_OSDisk", machineName)
}

// GenerateDataDiskName generates the name of a data disk based on the name of a VM.
func GenerateDataDiskName(machineName, nameSuffix string) string {
	return fmt.Sprintf("%s_%s", machineName, nameSuffix)
}

// GenerateVnetPeeringName generates the name for a peering between two vnets.
func GenerateVnetPeeringName(sourceVnetName string, remoteVnetName string) string {
	return fmt.Sprintf("%s-To-%s", sourceVnetName, remoteVnetName)
}

// GenerateAvailabilitySetName generates the name of a availability set based on the cluster name and the node group.
// node group identifies the set of nodes that belong to this availability set:
// For control plane nodes, this will be `control-plane`.
// For worker nodes, this will be the machine deployment name.
func GenerateAvailabilitySetName(clusterName, nodeGroup string) string {
	return fmt.Sprintf("%s_%s-as", clusterName, nodeGroup)
}

// WithIndex appends the index as suffix to a generated name.
func WithIndex(name string, n int) string {
	return fmt.Sprintf("%s-%d", name, n)
}

// ResourceGroupID returns the azure resource ID for a given resource group.
func ResourceGroupID(subscriptionID, resourceGroup string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionID, resourceGroup)
}

// VMID returns the azure resource ID for a given VM.
func VMID(subscriptionID, resourceGroup, vmName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s", subscriptionID, resourceGroup, vmName)
}

// VNetID returns the azure resource ID for a given VNet.
func VNetID(subscriptionID, resourceGroup, vnetName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s", subscriptionID, resourceGroup, vnetName)
}

// SubnetID returns the azure resource ID for a given subnet.
func SubnetID(subscriptionID, resourceGroup, vnetName, subnetName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s", subscriptionID, resourceGroup, vnetName, subnetName)
}

// PublicIPID returns the azure resource ID for a given public IP.
func PublicIPID(subscriptionID, resourceGroup, ipName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/publicIPAddresses/%s", subscriptionID, resourceGroup, ipName)
}

// RouteTableID returns the azure resource ID for a given route table.
func RouteTableID(subscriptionID, resourceGroup, routeTableName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/routeTables/%s", subscriptionID, resourceGroup, routeTableName)
}

// SecurityGroupID returns the azure resource ID for a given security group.
func SecurityGroupID(subscriptionID, resourceGroup, nsgName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s", subscriptionID, resourceGroup, nsgName)
}

// NatGatewayID returns the azure resource ID for a given NAT gateway.
func NatGatewayID(subscriptionID, resourceGroup, natgatewayName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/natGateways/%s", subscriptionID, resourceGroup, natgatewayName)
}

// NetworkInterfaceID returns the azure resource ID for a given network interface.
func NetworkInterfaceID(subscriptionID, resourceGroup, nicName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkInterfaces/%s", subscriptionID, resourceGroup, nicName)
}

// FrontendIPConfigID returns the azure resource ID for a given frontend IP config.
func FrontendIPConfigID(subscriptionID, resourceGroup, loadBalancerName, configName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/frontendIPConfigurations/%s", subscriptionID, resourceGroup, loadBalancerName, configName)
}

// AddressPoolID returns the azure resource ID for a given backend address pool.
func AddressPoolID(subscriptionID, resourceGroup, loadBalancerName, backendPoolName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/backendAddressPools/%s", subscriptionID, resourceGroup, loadBalancerName, backendPoolName)
}

// ProbeID returns the azure resource ID for a given probe.
func ProbeID(subscriptionID, resourceGroup, loadBalancerName, probeName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/probes/%s", subscriptionID, resourceGroup, loadBalancerName, probeName)
}

// NATRuleID returns the azure resource ID for a inbound NAT rule.
func NATRuleID(subscriptionID, resourceGroup, loadBalancerName, natRuleName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/inboundNatRules/%s", subscriptionID, resourceGroup, loadBalancerName, natRuleName)
}

// AvailabilitySetID returns the azure resource ID for a given availability set.
func AvailabilitySetID(subscriptionID, resourceGroup, availabilitySetName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/availabilitySets/%s", subscriptionID, resourceGroup, availabilitySetName)
}

// PrivateDNSZoneID returns the azure resource ID for a given private DNS zone.
func PrivateDNSZoneID(subscriptionID, resourceGroup, privateDNSZoneName string) string {
	return fmt.Sprintf("subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s", subscriptionID, resourceGroup, privateDNSZoneName)
}

// VirtualNetworkLinkID returns the azure resource ID for a given virtual network link.
func VirtualNetworkLinkID(subscriptionID, resourceGroup, privateDNSZoneName, virtualNetworkLinkName string) string {
	return fmt.Sprintf("subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s/virtualNetworkLinks/%s", subscriptionID, resourceGroup, privateDNSZoneName, virtualNetworkLinkName)
}

// ManagedClusterID returns the azure resource ID for a given managed cluster.
func ManagedClusterID(subscriptionID, resourceGroup, managedClusterName string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s", subscriptionID, resourceGroup, managedClusterName)
}

// GetBootstrappingVMExtension returns the CAPZ Bootstrapping VM extension.
// The CAPZ Bootstrapping extension is a simple clone of https://github.com/Azure/custom-script-extension-linux for Linux or
// https://learn.microsoft.com/azure/virtual-machines/extensions/custom-script-windows for Windows.
// This extension allows running arbitrary scripts on the VM.
// Its role is to detect and report Kubernetes bootstrap failure or success.
func GetBootstrappingVMExtension(osType string, cloud string, vmName string, cpuArchitectureType string) *ExtensionSpec {
	// currently, the bootstrap extension is only available in AzurePublicCloud.
	if osType == LinuxOS && cloud == PublicCloudName {
		// The command checks for the existence of the bootstrapSentinelFile on the machine, with retries and sleep between retries.
		// We set the version to 1.1 (will target 1.1.1) for arm64 machines and 1.0 for x64. This is due to a known issue with newer versions of
		// Go on Ubuntu 20.04. The issue is being tracked here: https://github.com/golang/go/issues/58550
		// TODO: Remove this once the issue is fixed, or when Ubuntu 20.04 is no longer supported.
		// We are using 1.1 instead of 1.1.1 for Arm64 as AzureAPI do not allow us to specify the full version.
		extensionVersion := "1.0"
		if cpuArchitectureType == string(armcompute.ArchitectureTypesArm64) {
			extensionVersion = "1.1"
		}
		return &ExtensionSpec{
			Name:      BootstrappingExtensionLinux,
			VMName:    vmName,
			Publisher: "Microsoft.Azure.ContainerUpstream",
			Version:   extensionVersion,
			ProtectedSettings: map[string]string{
				"commandToExecute": LinuxBootstrapExtensionCommand,
			},
		}
	} else if osType == WindowsOS && cloud == PublicCloudName {
		// This command for the existence of the bootstrapSentinelFile on the machine, with retries and sleep between reties.
		// If the file is not present after the retries are exhausted the extension fails with return code '-2' - ERROR_FILE_NOT_FOUND.
		return &ExtensionSpec{
			Name:      BootstrappingExtensionWindows,
			VMName:    vmName,
			Publisher: "Microsoft.Azure.ContainerUpstream",
			Version:   "1.0",
			ProtectedSettings: map[string]string{
				"commandToExecute": WindowsBootstrapExtensionCommand,
			},
		}
	}

	return nil
}

// UserAgent specifies a string to append to the agent identifier.
func UserAgent() string {
	return fmt.Sprintf("cluster-api-provider-azure/%s", version.Get().String())
}

// ARMClientOptions returns default ARM client options for CAPZ SDK v2 requests.
func ARMClientOptions(azureEnvironment string, extraPolicies ...policy.Policy) (*arm.ClientOptions, error) {
	opts := &arm.ClientOptions{}

	switch azureEnvironment {
	case PublicCloudName:
		opts.Cloud = cloud.AzurePublic
	case ChinaCloudName:
		opts.Cloud = cloud.AzureChina
	case USGovernmentCloudName:
		opts.Cloud = cloud.AzureGovernment
	case "":
		// No cloud name provided, so leave at defaults.
	default:
		return nil, fmt.Errorf("invalid cloud name %q", azureEnvironment)
	}
	opts.PerCallPolicies = []policy.Policy{
		correlationIDPolicy{},
		userAgentPolicy{},
	}
	opts.PerCallPolicies = append(opts.PerCallPolicies, extraPolicies...)
	opts.Retry.MaxRetries = -1 // Less than zero means one try and no retries.

	return opts, nil
}

// correlationIDPolicy adds the "x-ms-correlation-request-id" header to requests.
// It implements the policy.Policy interface.
type correlationIDPolicy struct{}

// Do adds the "x-ms-correlation-request-id" header if a request has a correlation ID in its context.
func (p correlationIDPolicy) Do(req *policy.Request) (*http.Response, error) {
	if corrID, ok := tele.CorrIDFromCtx(req.Raw().Context()); ok {
		req.Raw().Header.Set(string(tele.CorrIDKeyVal), string(corrID))
	}
	return req.Next()
}

// userAgentPolicy extends the "User-Agent" header on requests.
// It implements the policy.Policy interface.
type userAgentPolicy struct{}

// Do extends the "User-Agent" header of a request by appending CAPZ's user agent.
func (p userAgentPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Set("User-Agent", req.Raw().UserAgent()+" "+UserAgent())
	return req.Next()
}

// CustomPutPatchHeaderPolicy adds custom headers to a PUT or PATCH request.
// It implements the policy.Policy interface.
type CustomPutPatchHeaderPolicy struct {
	Headers map[string]string
}

// Do adds any custom headers to a PUT or PATCH request.
func (p CustomPutPatchHeaderPolicy) Do(req *policy.Request) (*http.Response, error) {
	if req.Raw().Method == http.MethodPut || req.Raw().Method == http.MethodPatch {
		for key, element := range p.Headers {
			req.Raw().Header.Set(key, element)
		}
	}

	return req.Next()
}
