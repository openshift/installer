/*
Copyright 2023 The Kubernetes Authors.

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
	"strings"

	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (mcp *AzureManagedControlPlaneTemplate) setDefaults() {
	setDefault[*string](&mcp.Spec.Template.Spec.NetworkPlugin, ptr.To(AzureNetworkPluginName))
	setDefault[*string](&mcp.Spec.Template.Spec.LoadBalancerSKU, ptr.To("Standard"))
	setDefault[*bool](&mcp.Spec.Template.Spec.EnablePreviewFeatures, ptr.To(false))

	if mcp.Spec.Template.Spec.Version != "" && !strings.HasPrefix(mcp.Spec.Template.Spec.Version, "v") {
		mcp.Spec.Template.Spec.Version = setDefaultVersion(mcp.Spec.Template.Spec.Version)
	}

	mcp.setDefaultVirtualNetwork()
	mcp.setDefaultSubnet()
	mcp.Spec.Template.Spec.SKU = setDefaultSku(mcp.Spec.Template.Spec.SKU)
}

// setDefaultVirtualNetwork sets the default VirtualNetwork for an AzureManagedControlPlaneTemplate.
func (mcp *AzureManagedControlPlaneTemplate) setDefaultVirtualNetwork() {
	if mcp.Spec.Template.Spec.VirtualNetwork.Name != "" {
		// Being able to set the vnet name in the template type is a bug, as vnet names cannot be reused across clusters.
		// To avoid a breaking API change, a warning is logged.
		ctrl.Log.WithName("AzureManagedControlPlaneTemplateWebHookLogger").Info("WARNING: VirtualNetwork.Name should not be set in the template. Virtual Network names cannot be shared across clusters.")
	}
	if mcp.Spec.Template.Spec.VirtualNetwork.CIDRBlock == "" {
		mcp.Spec.Template.Spec.VirtualNetwork.CIDRBlock = defaultAKSVnetCIDR
	}
}

// setDefaultSubnet sets the default Subnet for an AzureManagedControlPlaneTemplate.
func (mcp *AzureManagedControlPlaneTemplate) setDefaultSubnet() {
	if mcp.Spec.Template.Spec.VirtualNetwork.Subnet.Name == "" {
		mcp.Spec.Template.Spec.VirtualNetwork.Subnet.Name = mcp.Name
	}
	if mcp.Spec.Template.Spec.VirtualNetwork.Subnet.CIDRBlock == "" {
		mcp.Spec.Template.Spec.VirtualNetwork.Subnet.CIDRBlock = defaultAKSNodeSubnetCIDR
	}
}

// setDefault sets the default value for a pointer to a value for any comparable type.
func setDefault[T comparable](field *T, value T) {
	if field == nil {
		// shouldn't happen with proper use
		return
	}
	var zero T
	if *field == zero {
		*field = value
	}
}
