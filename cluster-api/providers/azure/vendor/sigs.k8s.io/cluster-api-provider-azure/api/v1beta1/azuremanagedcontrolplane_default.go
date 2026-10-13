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
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"

	utilSSH "sigs.k8s.io/cluster-api-provider-azure/util/ssh"
)

const (
	// defaultAKSVnetCIDR is the default Vnet CIDR.
	defaultAKSVnetCIDR = "10.0.0.0/8"
	// defaultAKSNodeSubnetCIDR is the default Node Subnet CIDR.
	defaultAKSNodeSubnetCIDR = "10.240.0.0/16"
	// defaultAKSVnetCIDRForOverlay is the default Vnet CIDR when Azure CNI overlay is enabled.
	defaultAKSVnetCIDRForOverlay = "10.224.0.0/12"
	// defaultAKSNodeSubnetCIDRForOverlay is the default Node Subnet CIDR when Azure CNI overlay is enabled.
	defaultAKSNodeSubnetCIDRForOverlay = "10.224.0.0/16"
)

// setDefaultResourceGroupName sets the default ResourceGroupName for an AzureManagedControlPlane.
func (m *AzureManagedControlPlane) setDefaultResourceGroupName() {
	if m.Spec.ResourceGroupName == "" {
		if clusterName, ok := m.Labels[clusterv1.ClusterNameLabel]; ok {
			m.Spec.ResourceGroupName = clusterName
		}
	}
}

// setDefaultSSHPublicKey sets the default SSHPublicKey for an AzureManagedControlPlane.
func (m *AzureManagedControlPlane) setDefaultSSHPublicKey() error {
	if sshKey := m.Spec.SSHPublicKey; sshKey != nil && *sshKey == "" {
		_, publicRsaKey, err := utilSSH.GenerateSSHKey()
		if err != nil {
			return err
		}

		m.Spec.SSHPublicKey = ptr.To(base64.StdEncoding.EncodeToString(ssh.MarshalAuthorizedKey(publicRsaKey)))
	}

	return nil
}

// setDefaultNodeResourceGroupName sets the default NodeResourceGroup for an AzureManagedControlPlane.
func (m *AzureManagedControlPlane) setDefaultNodeResourceGroupName() {
	if m.Spec.NodeResourceGroupName == "" {
		m.Spec.NodeResourceGroupName = fmt.Sprintf("MC_%s_%s_%s", m.Spec.ResourceGroupName, m.Name, m.Spec.Location)
	}
}

// setDefaultVirtualNetwork sets the default VirtualNetwork for an AzureManagedControlPlane.
func (m *AzureManagedControlPlane) setDefaultVirtualNetwork() {
	if m.Spec.VirtualNetwork.Name == "" {
		m.Spec.VirtualNetwork.Name = m.Name
	}
	if m.Spec.VirtualNetwork.CIDRBlock == "" {
		m.Spec.VirtualNetwork.CIDRBlock = defaultAKSVnetCIDR
		if ptr.Deref(m.Spec.NetworkPluginMode, "") == NetworkPluginModeOverlay {
			m.Spec.VirtualNetwork.CIDRBlock = defaultAKSVnetCIDRForOverlay
		}
	}
	if m.Spec.VirtualNetwork.ResourceGroup == "" {
		m.Spec.VirtualNetwork.ResourceGroup = m.Spec.ResourceGroupName
	}
}

// setDefaultSubnet sets the default Subnet for an AzureManagedControlPlane.
func (m *AzureManagedControlPlane) setDefaultSubnet() {
	if m.Spec.VirtualNetwork.Subnet.Name == "" {
		m.Spec.VirtualNetwork.Subnet.Name = m.Name
	}
	if m.Spec.VirtualNetwork.Subnet.CIDRBlock == "" {
		m.Spec.VirtualNetwork.Subnet.CIDRBlock = defaultAKSNodeSubnetCIDR
		if ptr.Deref(m.Spec.NetworkPluginMode, "") == NetworkPluginModeOverlay {
			m.Spec.VirtualNetwork.Subnet.CIDRBlock = defaultAKSNodeSubnetCIDRForOverlay
		}
	}
}

// setDefaultFleetsMember sets the default FleetsMember for an AzureManagedControlPlane.
func setDefaultFleetsMember(fleetsMember *FleetsMember, labels map[string]string) *FleetsMember {
	result := fleetsMember.DeepCopy()
	if fleetsMember != nil {
		if clusterName, ok := labels[clusterv1.ClusterNameLabel]; ok && fleetsMember.Name == "" {
			result.Name = clusterName
		}
	}
	return result
}

func setDefaultSku(sku *AKSSku) *AKSSku {
	result := sku.DeepCopy()
	if sku == nil {
		result = new(AKSSku)
		result.Tier = FreeManagedControlPlaneTier
	} else if sku.Tier == PaidManagedControlPlaneTier {
		result.Tier = StandardManagedControlPlaneTier
		ctrl.Log.WithName("AzureManagedControlPlaneWebHookLogger").Info("Paid SKU tier is deprecated and has been replaced by Standard")
	}
	return result
}

func setDefaultVersion(version string) string {
	if version != "" && !strings.HasPrefix(version, "v") {
		normalizedVersion := "v" + version
		version = normalizedVersion
	}
	return version
}

func (m *AzureManagedControlPlane) setDefaultOIDCIssuerProfile() {
	if m.Spec.OIDCIssuerProfile == nil {
		m.Spec.OIDCIssuerProfile = &OIDCIssuerProfile{}
	}
}

func (m *AzureManagedControlPlane) setDefaultDNSPrefix() {
	if m.Spec.DNSPrefix == nil {
		m.Spec.DNSPrefix = ptr.To(m.Name)
	}
}

func (m *AzureManagedControlPlane) setDefaultAKSExtensions() {
	for _, extension := range m.Spec.Extensions {
		if extension.Plan != nil && extension.Plan.Name == "" {
			extension.Plan.Name = fmt.Sprintf("%s-%s", m.Name, extension.Plan.Product)
		}
	}
}
