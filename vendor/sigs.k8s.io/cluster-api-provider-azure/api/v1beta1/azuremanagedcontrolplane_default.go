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

	"golang.org/x/crypto/ssh"
	"k8s.io/utils/ptr"
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

func (m *AzureManagedControlPlane) setDefaultSku() {
	if m.Spec.SKU == nil {
		m.Spec.SKU = &AKSSku{
			Tier: FreeManagedControlPlaneTier,
		}
	}
}

func (m *AzureManagedControlPlane) setDefaultAutoScalerProfile() {
	if m.Spec.AutoScalerProfile == nil {
		return
	}

	// Default values are from https://learn.microsoft.com/en-us/azure/aks/cluster-autoscaler#using-the-autoscaler-profile
	// If any values are set, they all need to be set.
	if m.Spec.AutoScalerProfile.BalanceSimilarNodeGroups == nil {
		m.Spec.AutoScalerProfile.BalanceSimilarNodeGroups = (*BalanceSimilarNodeGroups)(ptr.To(string(BalanceSimilarNodeGroupsFalse)))
	}
	if m.Spec.AutoScalerProfile.Expander == nil {
		m.Spec.AutoScalerProfile.Expander = (*Expander)(ptr.To(string(ExpanderRandom)))
	}
	if m.Spec.AutoScalerProfile.MaxEmptyBulkDelete == nil {
		m.Spec.AutoScalerProfile.MaxEmptyBulkDelete = ptr.To("10")
	}
	if m.Spec.AutoScalerProfile.MaxGracefulTerminationSec == nil {
		m.Spec.AutoScalerProfile.MaxGracefulTerminationSec = ptr.To("600")
	}
	if m.Spec.AutoScalerProfile.MaxNodeProvisionTime == nil {
		m.Spec.AutoScalerProfile.MaxNodeProvisionTime = ptr.To("15m")
	}
	if m.Spec.AutoScalerProfile.MaxTotalUnreadyPercentage == nil {
		m.Spec.AutoScalerProfile.MaxTotalUnreadyPercentage = ptr.To("45")
	}
	if m.Spec.AutoScalerProfile.NewPodScaleUpDelay == nil {
		m.Spec.AutoScalerProfile.NewPodScaleUpDelay = ptr.To("0s")
	}
	if m.Spec.AutoScalerProfile.OkTotalUnreadyCount == nil {
		m.Spec.AutoScalerProfile.OkTotalUnreadyCount = ptr.To("3")
	}
	if m.Spec.AutoScalerProfile.ScanInterval == nil {
		m.Spec.AutoScalerProfile.ScanInterval = ptr.To("10s")
	}
	if m.Spec.AutoScalerProfile.ScaleDownDelayAfterAdd == nil {
		m.Spec.AutoScalerProfile.ScaleDownDelayAfterAdd = ptr.To("10m")
	}
	if m.Spec.AutoScalerProfile.ScaleDownDelayAfterDelete == nil {
		// Default is the same as the ScanInterval so default to that same value if it isn't set
		m.Spec.AutoScalerProfile.ScaleDownDelayAfterDelete = m.Spec.AutoScalerProfile.ScanInterval
	}
	if m.Spec.AutoScalerProfile.ScaleDownDelayAfterFailure == nil {
		m.Spec.AutoScalerProfile.ScaleDownDelayAfterFailure = ptr.To("3m")
	}
	if m.Spec.AutoScalerProfile.ScaleDownUnneededTime == nil {
		m.Spec.AutoScalerProfile.ScaleDownUnneededTime = ptr.To("10m")
	}
	if m.Spec.AutoScalerProfile.ScaleDownUnreadyTime == nil {
		m.Spec.AutoScalerProfile.ScaleDownUnreadyTime = ptr.To("20m")
	}
	if m.Spec.AutoScalerProfile.ScaleDownUtilizationThreshold == nil {
		m.Spec.AutoScalerProfile.ScaleDownUtilizationThreshold = ptr.To("0.5")
	}
	if m.Spec.AutoScalerProfile.SkipNodesWithLocalStorage == nil {
		m.Spec.AutoScalerProfile.SkipNodesWithLocalStorage = (*SkipNodesWithLocalStorage)(ptr.To(string(SkipNodesWithLocalStorageFalse)))
	}
	if m.Spec.AutoScalerProfile.SkipNodesWithSystemPods == nil {
		m.Spec.AutoScalerProfile.SkipNodesWithSystemPods = (*SkipNodesWithSystemPods)(ptr.To(string(SkipNodesWithSystemPodsTrue)))
	}
}

func (m *AzureManagedControlPlane) setDefaultOIDCIssuerProfile() {
	if m.Spec.OIDCIssuerProfile == nil {
		m.Spec.OIDCIssuerProfile = &OIDCIssuerProfile{}
	}

	if m.Spec.OIDCIssuerProfile.Enabled == nil {
		m.Spec.OIDCIssuerProfile.Enabled = ptr.To(false)
	}
}

func (m *AzureManagedControlPlane) setDefaultDNSPrefix() {
	if m.Spec.DNSPrefix == nil {
		m.Spec.DNSPrefix = ptr.To(m.Name)
	}
}
