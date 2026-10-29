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

package v1beta1

import (
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// SetDefaultsAzureManagedControlPlaneTemplate sets the default values for an AzureManagedControlPlaneTemplate.
func SetDefaultsAzureManagedControlPlaneTemplate(log logr.Logger, mcp *infrav1.AzureManagedControlPlaneTemplate) {
	SetZeroPointerDefault[*string](&mcp.Spec.Template.Spec.NetworkPlugin, ptr.To(infrav1.AzureNetworkPluginName))
	SetZeroPointerDefault[*string](&mcp.Spec.Template.Spec.LoadBalancerSKU, ptr.To("Standard"))
	SetZeroPointerDefault[*bool](&mcp.Spec.Template.Spec.EnablePreviewFeatures, ptr.To(false))

	if mcp.Spec.Template.Spec.Version != "" && !strings.HasPrefix(mcp.Spec.Template.Spec.Version, "v") {
		mcp.Spec.Template.Spec.Version = NormalizeVersion(mcp.Spec.Template.Spec.Version)
	}

	setDefaultAzureManagedControlPlaneTemplateVirtualNetwork(log, mcp)
	setDefaultAzureManagedControlPlaneTemplateSubnet(mcp)
	mcp.Spec.Template.Spec.SKU = DefaultSku(log, mcp.Spec.Template.Spec.SKU)
}

// setDefaultAzureManagedControlPlaneTemplateVirtualNetwork sets the default VirtualNetwork for an AzureManagedControlPlaneTemplate.
func setDefaultAzureManagedControlPlaneTemplateVirtualNetwork(log logr.Logger, mcp *infrav1.AzureManagedControlPlaneTemplate) {
	if mcp.Spec.Template.Spec.VirtualNetwork.Name != "" {
		// Being able to set the vnet name in the template type is a bug, as vnet names cannot be reused across clusters.
		// To avoid a breaking API change, a warning is logged.
		log.Info("WARNING: VirtualNetwork.Name should not be set in the template. Virtual Network names cannot be shared across clusters.")
	}
	if mcp.Spec.Template.Spec.VirtualNetwork.CIDRBlock == "" {
		mcp.Spec.Template.Spec.VirtualNetwork.CIDRBlock = DefaultAKSVnetCIDR
	}
}

// setDefaultAzureManagedControlPlaneTemplateSubnet sets the default Subnet for an AzureManagedControlPlaneTemplate.
func setDefaultAzureManagedControlPlaneTemplateSubnet(mcp *infrav1.AzureManagedControlPlaneTemplate) {
	if mcp.Spec.Template.Spec.VirtualNetwork.Subnet.Name == "" {
		mcp.Spec.Template.Spec.VirtualNetwork.Subnet.Name = mcp.Name
	}
	if mcp.Spec.Template.Spec.VirtualNetwork.Subnet.CIDRBlock == "" {
		mcp.Spec.Template.Spec.VirtualNetwork.Subnet.CIDRBlock = DefaultAKSNodeSubnetCIDR
	}
}
