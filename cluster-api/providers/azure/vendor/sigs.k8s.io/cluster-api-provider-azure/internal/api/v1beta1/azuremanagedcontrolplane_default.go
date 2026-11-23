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
	"encoding/base64"
	"fmt"

	"github.com/go-logr/logr"
	"golang.org/x/crypto/ssh"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	utilSSH "sigs.k8s.io/cluster-api-provider-azure/util/ssh"
)

const (
	// DefaultAKSVnetCIDR is the default Vnet CIDR.
	DefaultAKSVnetCIDR = "10.0.0.0/8"
	// DefaultAKSNodeSubnetCIDR is the default Node Subnet CIDR.
	DefaultAKSNodeSubnetCIDR = "10.240.0.0/16"
	// DefaultAKSVnetCIDRForOverlay is the default Vnet CIDR when Azure CNI overlay is enabled.
	DefaultAKSVnetCIDRForOverlay = "10.224.0.0/12"
	// DefaultAKSNodeSubnetCIDRForOverlay is the default Node Subnet CIDR when Azure CNI overlay is enabled.
	DefaultAKSNodeSubnetCIDRForOverlay = "10.224.0.0/16"
)

// DefaultFleetsMember sets the default FleetsMember for an AzureManagedControlPlane.
func DefaultFleetsMember(fleetsMember *infrav1.FleetsMember, labels map[string]string) *infrav1.FleetsMember {
	result := fleetsMember.DeepCopy()
	if fleetsMember != nil {
		if clusterName, ok := labels[clusterv1.ClusterNameLabel]; ok && fleetsMember.Name == "" {
			result.Name = clusterName
		}
	}
	return result
}

// DefaultSku returns the default SKU for an AzureManagedControlPlane.
func DefaultSku(log logr.Logger, sku *infrav1.AKSSku) *infrav1.AKSSku {
	result := sku.DeepCopy()
	if sku == nil {
		result = new(infrav1.AKSSku)
		result.Tier = infrav1.FreeManagedControlPlaneTier
	} else if sku.Tier == infrav1.PaidManagedControlPlaneTier {
		result.Tier = infrav1.StandardManagedControlPlaneTier
		log.Info("Paid SKU tier is deprecated and has been replaced by Standard")
	}
	return result
}

// SetDefaultAzureManagedControlPlaneResourceGroupName sets the default ResourceGroupName for an AzureManagedControlPlane.
func SetDefaultAzureManagedControlPlaneResourceGroupName(m *infrav1.AzureManagedControlPlane) {
	if m.Spec.ResourceGroupName == "" {
		if clusterName, ok := m.Labels[clusterv1.ClusterNameLabel]; ok {
			m.Spec.ResourceGroupName = clusterName
		}
	}
}

// SetDefaultAzureManagedControlPlaneSSHPublicKey sets the default SSHPublicKey for an AzureManagedControlPlane.
func SetDefaultAzureManagedControlPlaneSSHPublicKey(m *infrav1.AzureManagedControlPlane) error {
	if sshKey := m.Spec.SSHPublicKey; sshKey != nil && *sshKey == "" {
		_, publicRsaKey, err := utilSSH.GenerateSSHKey()
		if err != nil {
			return err
		}

		m.Spec.SSHPublicKey = ptr.To(base64.StdEncoding.EncodeToString(ssh.MarshalAuthorizedKey(publicRsaKey)))
	}

	return nil
}

// SetDefaultAzureManagedControlPlaneNodeResourceGroupName sets the default NodeResourceGroup for an AzureManagedControlPlane.
func SetDefaultAzureManagedControlPlaneNodeResourceGroupName(m *infrav1.AzureManagedControlPlane) {
	if m.Spec.NodeResourceGroupName == "" {
		m.Spec.NodeResourceGroupName = fmt.Sprintf("MC_%s_%s_%s", m.Spec.ResourceGroupName, m.Name, m.Spec.Location)
	}
}

// SetDefaultAzureManagedControlPlaneVirtualNetwork sets the default VirtualNetwork for an AzureManagedControlPlane.
func SetDefaultAzureManagedControlPlaneVirtualNetwork(m *infrav1.AzureManagedControlPlane) {
	if m.Spec.VirtualNetwork.Name == "" {
		m.Spec.VirtualNetwork.Name = m.Name
	}
	if m.Spec.VirtualNetwork.CIDRBlock == "" {
		m.Spec.VirtualNetwork.CIDRBlock = DefaultAKSVnetCIDR
		if ptr.Deref(m.Spec.NetworkPluginMode, "") == infrav1.NetworkPluginModeOverlay {
			m.Spec.VirtualNetwork.CIDRBlock = DefaultAKSVnetCIDRForOverlay
		}
	}
	if m.Spec.VirtualNetwork.ResourceGroup == "" {
		m.Spec.VirtualNetwork.ResourceGroup = m.Spec.ResourceGroupName
	}
}

// SetDefaultAzureManagedControlPlaneSubnet sets the default Subnet for an AzureManagedControlPlane.
func SetDefaultAzureManagedControlPlaneSubnet(m *infrav1.AzureManagedControlPlane) {
	if m.Spec.VirtualNetwork.Subnet.Name == "" {
		m.Spec.VirtualNetwork.Subnet.Name = m.Name
	}
	if m.Spec.VirtualNetwork.Subnet.CIDRBlock == "" {
		m.Spec.VirtualNetwork.Subnet.CIDRBlock = DefaultAKSNodeSubnetCIDR
		if ptr.Deref(m.Spec.NetworkPluginMode, "") == infrav1.NetworkPluginModeOverlay {
			m.Spec.VirtualNetwork.Subnet.CIDRBlock = DefaultAKSNodeSubnetCIDRForOverlay
		}
	}
}

// SetDefaultAzureManagedControlPlaneOIDCIssuerProfile sets the default OIDCIssuerProfile for an AzureManagedControlPlane.
func SetDefaultAzureManagedControlPlaneOIDCIssuerProfile(m *infrav1.AzureManagedControlPlane) {
	if m.Spec.OIDCIssuerProfile == nil {
		m.Spec.OIDCIssuerProfile = &infrav1.OIDCIssuerProfile{}
	}
}

// SetDefaultAzureManagedControlPlaneDNSPrefix sets the default DNSPrefix for an AzureManagedControlPlane.
func SetDefaultAzureManagedControlPlaneDNSPrefix(m *infrav1.AzureManagedControlPlane) {
	if m.Spec.DNSPrefix == nil {
		m.Spec.DNSPrefix = ptr.To(m.Name)
	}
}

// SetDefaultAzureManagedControlPlaneAKSExtensions sets the default AKS extensions for an AzureManagedControlPlane.
func SetDefaultAzureManagedControlPlaneAKSExtensions(m *infrav1.AzureManagedControlPlane) {
	for _, extension := range m.Spec.Extensions {
		if extension.Plan != nil && extension.Plan.Name == "" {
			extension.Plan.Name = fmt.Sprintf("%s-%s", m.Name, extension.Plan.Product)
		}
	}
}
