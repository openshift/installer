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
	utilSSH "sigs.k8s.io/cluster-api-provider-azure/util/ssh"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
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
	if clusterName, ok := labels[clusterv1.ClusterNameLabel]; ok && fleetsMember != nil && fleetsMember.Name == "" {
		result.Name = clusterName
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

func setDefaultAutoScalerProfile(autoScalerProfile *AutoScalerProfile) *AutoScalerProfile {
	if autoScalerProfile == nil {
		return nil
	}

	result := autoScalerProfile.DeepCopy()

	// Default values are from https://learn.microsoft.com/en-us/azure/aks/cluster-autoscaler#using-the-autoscaler-profile
	// If any values are set, they all need to be set.
	if autoScalerProfile.BalanceSimilarNodeGroups == nil {
		result.BalanceSimilarNodeGroups = (*BalanceSimilarNodeGroups)(ptr.To(string(BalanceSimilarNodeGroupsFalse)))
	}
	if autoScalerProfile.Expander == nil {
		result.Expander = (*Expander)(ptr.To(string(ExpanderRandom)))
	}
	if autoScalerProfile.MaxEmptyBulkDelete == nil {
		result.MaxEmptyBulkDelete = ptr.To("10")
	}
	if autoScalerProfile.MaxGracefulTerminationSec == nil {
		result.MaxGracefulTerminationSec = ptr.To("600")
	}
	if autoScalerProfile.MaxNodeProvisionTime == nil {
		result.MaxNodeProvisionTime = ptr.To("15m")
	}
	if autoScalerProfile.MaxTotalUnreadyPercentage == nil {
		result.MaxTotalUnreadyPercentage = ptr.To("45")
	}
	if autoScalerProfile.NewPodScaleUpDelay == nil {
		result.NewPodScaleUpDelay = ptr.To("0s")
	}
	if autoScalerProfile.OkTotalUnreadyCount == nil {
		result.OkTotalUnreadyCount = ptr.To("3")
	}
	if autoScalerProfile.ScanInterval == nil {
		result.ScanInterval = ptr.To("10s")
	}
	if autoScalerProfile.ScaleDownDelayAfterAdd == nil {
		result.ScaleDownDelayAfterAdd = ptr.To("10m")
	}
	if autoScalerProfile.ScaleDownDelayAfterDelete == nil {
		// Default is the same as the ScanInterval so default to that same value if it isn't set
		result.ScaleDownDelayAfterDelete = result.ScanInterval
	}
	if autoScalerProfile.ScaleDownDelayAfterFailure == nil {
		result.ScaleDownDelayAfterFailure = ptr.To("3m")
	}
	if autoScalerProfile.ScaleDownUnneededTime == nil {
		result.ScaleDownUnneededTime = ptr.To("10m")
	}
	if autoScalerProfile.ScaleDownUnreadyTime == nil {
		result.ScaleDownUnreadyTime = ptr.To("20m")
	}
	if autoScalerProfile.ScaleDownUtilizationThreshold == nil {
		result.ScaleDownUtilizationThreshold = ptr.To("0.5")
	}
	if autoScalerProfile.SkipNodesWithLocalStorage == nil {
		result.SkipNodesWithLocalStorage = (*SkipNodesWithLocalStorage)(ptr.To(string(SkipNodesWithLocalStorageFalse)))
	}
	if autoScalerProfile.SkipNodesWithSystemPods == nil {
		result.SkipNodesWithSystemPods = (*SkipNodesWithSystemPods)(ptr.To(string(SkipNodesWithSystemPodsTrue)))
	}

	return result
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

func (m *AzureManagedControlPlane) setDefaultAKSExtensions() {
	for _, extension := range m.Spec.Extensions {
		if extension.Plan != nil && extension.Plan.Name == "" {
			extension.Plan.Name = fmt.Sprintf("%s-%s", m.Name, extension.Plan.Product)
		}
		if extension.AutoUpgradeMinorVersion == nil {
			extension.AutoUpgradeMinorVersion = ptr.To(true)
		}
	}
}
