/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
)

// AzureClusterClassSpec defines the AzureCluster properties that may be shared across several Azure clusters.
type AzureClusterClassSpec struct {
	// +optional
	SubscriptionID string `json:"subscriptionID,omitempty"`

	Location string `json:"location"`

	// ExtendedLocation is an optional set of ExtendedLocation properties for clusters on Azure public MEC.
	// +optional
	ExtendedLocation *ExtendedLocationSpec `json:"extendedLocation,omitempty"`

	// AdditionalTags is an optional set of tags to add to Azure resources managed by the Azure provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// IdentityRef is a reference to an AzureIdentity to be used when reconciling this cluster
	// +optional
	IdentityRef *corev1.ObjectReference `json:"identityRef,omitempty"`

	// AzureEnvironment is the name of the AzureCloud to be used.
	// The default value that would be used by most users is "AzurePublicCloud", other values are:
	// - ChinaCloud: "AzureChinaCloud"
	// - GermanCloud: "AzureGermanCloud"
	// - PublicCloud: "AzurePublicCloud"
	// - USGovernmentCloud: "AzureUSGovernmentCloud"
	// +optional
	AzureEnvironment string `json:"azureEnvironment,omitempty"`

	// CloudProviderConfigOverrides is an optional set of configuration values that can be overridden in azure cloud provider config.
	// This is only a subset of options that are available in azure cloud provider config.
	// Some values for the cloud provider config are inferred from other parts of cluster api provider azure spec, and may not be available for overrides.
	// See: https://cloud-provider-azure.sigs.k8s.io/install/configs
	// Note: All cloud provider config values can be customized by creating the secret beforehand. CloudProviderConfigOverrides is only used when the secret is managed by the Azure Provider.
	// +optional
	CloudProviderConfigOverrides *CloudProviderConfigOverrides `json:"cloudProviderConfigOverrides,omitempty"`
}

// ExtendedLocationSpec defines the ExtendedLocation properties to enable CAPZ for Azure public MEC.
type ExtendedLocationSpec struct {
	// Name defines the name for the extended location.
	Name string `json:"name"`

	// Type defines the type for the extended location.
	// +kubebuilder:validation:Enum=EdgeZone
	Type string `json:"type"`
}

// NetworkClassSpec defines the NetworkSpec properties that may be shared across several Azure clusters.
type NetworkClassSpec struct {
	// PrivateDNSZoneName defines the zone name for the Azure Private DNS.
	// +optional
	PrivateDNSZoneName string `json:"privateDNSZoneName,omitempty"`
}

// VnetClassSpec defines the VnetSpec properties that may be shared across several Azure clusters.
type VnetClassSpec struct {
	// CIDRBlocks defines the virtual network's address space, specified as one or more address prefixes in CIDR notation.
	// +optional
	CIDRBlocks []string `json:"cidrBlocks,omitempty"`

	// Tags is a collection of tags describing the resource.
	// +optional
	Tags Tags `json:"tags,omitempty"`
}

// SubnetClassSpec defines the SubnetSpec properties that may be shared across several Azure clusters.
type SubnetClassSpec struct {
	// Name defines a name for the subnet resource.
	Name string `json:"name"`

	// Role defines the subnet role (eg. Node, ControlPlane)
	// +kubebuilder:validation:Enum=node;control-plane;bastion
	Role SubnetRole `json:"role"`

	// CIDRBlocks defines the subnet's address space, specified as one or more address prefixes in CIDR notation.
	// +optional
	CIDRBlocks []string `json:"cidrBlocks,omitempty"`

	// ServiceEndpoints is a slice of Virtual Network service endpoints to enable for the subnets.
	// +optional
	ServiceEndpoints ServiceEndpoints `json:"serviceEndpoints,omitempty"`

	// PrivateEndpoints defines a list of private endpoints that should be attached to this subnet.
	// +optional
	PrivateEndpoints PrivateEndpoints `json:"privateEndpoints,omitempty"`
}

// LoadBalancerClassSpec defines the LoadBalancerSpec properties that may be shared across several Azure clusters.
type LoadBalancerClassSpec struct {
	// +optional
	SKU SKU `json:"sku,omitempty"`
	// +optional
	Type LBType `json:"type,omitempty"`
	// IdleTimeoutInMinutes specifies the timeout for the TCP idle connection.
	// +optional
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
}

// SecurityGroupClass defines the SecurityGroup properties that may be shared across several Azure clusters.
type SecurityGroupClass struct {
	// +optional
	SecurityRules SecurityRules `json:"securityRules,omitempty"`
	// +optional
	Tags Tags `json:"tags,omitempty"`
}

// FrontendIPClass defines the FrontendIP properties that may be shared across several Azure clusters.
type FrontendIPClass struct {
	// +optional
	PrivateIPAddress string `json:"privateIP,omitempty"`
}

// setDefaults sets default values for AzureClusterClassSpec.
func (acc *AzureClusterClassSpec) setDefaults() {
	if acc.AzureEnvironment == "" {
		acc.AzureEnvironment = DefaultAzureCloud
	}
}

// setDefaults sets default values for VnetClassSpec.
func (vc *VnetClassSpec) setDefaults() {
	if len(vc.CIDRBlocks) == 0 {
		vc.CIDRBlocks = []string{DefaultVnetCIDR}
	}
}

// setDefaults sets default values for SubnetClassSpec.
func (sc *SubnetClassSpec) setDefaults(cidr string) {
	if len(sc.CIDRBlocks) == 0 {
		sc.CIDRBlocks = []string{cidr}
	}
}

// setDefaults sets default values for SecurityGroupClass.
func (sgc *SecurityGroupClass) setDefaults() {
	for i := range sgc.SecurityRules {
		if sgc.SecurityRules[i].Direction == "" {
			sgc.SecurityRules[i].Direction = SecurityRuleDirectionInbound
		}
	}
}
