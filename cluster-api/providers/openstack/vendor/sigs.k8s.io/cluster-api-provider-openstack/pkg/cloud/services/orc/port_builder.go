/*
Copyright 2026 The Kubernetes Authors.

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

package orc

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// buildPort builds a managed ORC Port. ORC will create the Neutron port.
func buildPort(
	serverName, namespace string,
	index int,
	portOpts infrav1.PortOpts,
	globalSGs []infrav1.SecurityGroupParam,
	defaultTrunk bool,
	serverTags []string,
	networkNameMap, subnetNameMap, sgNameMap map[string]string,
	credRef orcv1alpha1.CloudCredentialsReference,
) *orcv1alpha1.Port {
	portSpec := &orcv1alpha1.PortResourceSpec{
		PortSecurity: convertPortSecurity(portOpts.EnablePortSecurity),
	}

	// Network reference
	if portOpts.Network != nil {
		key := NetworkParamKey(*portOpts.Network)
		if name, ok := networkNameMap[key]; ok {
			portSpec.NetworkRef = orcv1alpha1.KubernetesNameRef(name)
		}
	}

	// Fixed IPs → Addresses
	for _, fixedIP := range portOpts.FixedIPs {
		addr := orcv1alpha1.Address{}
		if fixedIP.Subnet != nil {
			key := SubnetParamKey(*fixedIP.Subnet)
			if name, ok := subnetNameMap[key]; ok {
				addr.SubnetRef = orcv1alpha1.KubernetesNameRef(name)
			}
		}
		if fixedIP.IPAddress != nil {
			ip := orcv1alpha1.IPvAny(*fixedIP.IPAddress)
			addr.IP = &ip
		}
		// Only add if we have a subnet ref (required by ORC)
		if addr.SubnetRef != "" {
			portSpec.Addresses = append(portSpec.Addresses, addr)
		}
	}

	// Security groups (per-port + global, deduplicated)
	seen := make(map[string]bool)
	for _, sg := range portOpts.SecurityGroups {
		key := SecurityGroupParamKey(sg)
		if name, ok := sgNameMap[key]; ok && !seen[name] {
			portSpec.SecurityGroupRefs = append(portSpec.SecurityGroupRefs, orcv1alpha1.KubernetesNameRef(name))
			seen[name] = true
		}
	}
	for _, sg := range globalSGs {
		key := SecurityGroupParamKey(sg)
		if name, ok := sgNameMap[key]; ok && !seen[name] {
			portSpec.SecurityGroupRefs = append(portSpec.SecurityGroupRefs, orcv1alpha1.KubernetesNameRef(name))
			seen[name] = true
		}
	}

	// Allowed address pairs
	portSpec.AllowedAddressPairs = convertAllowedAddressPairs(portOpts.AllowedAddressPairs)

	// VNIC type
	if portOpts.VNICType != nil {
		portSpec.VNICType = *portOpts.VNICType
	}

	// Admin state
	portSpec.AdminStateUp = portOpts.AdminStateUp

	// MAC address
	if portOpts.MACAddress != nil {
		portSpec.MACAddress = *portOpts.MACAddress
	}

	// Host ID
	if portOpts.HostID != nil {
		portSpec.HostID = &orcv1alpha1.HostID{
			ID: *portOpts.HostID,
		}
	}

	// Propagate uplink status
	portSpec.PropagateUplinkStatus = portOpts.PropagateUplinkStatus

	// Value specs
	portSpec.ValueSpecs = convertValueSpecs(portOpts.ValueSpecs)

	// Trusted VIF (from binding profile)
	if portOpts.Profile != nil && ptr.Deref(portOpts.Profile.TrustedVF, false) {
		portSpec.TrustedVIF = portOpts.Profile.TrustedVF
	}

	// Description
	if portOpts.Description != nil {
		portSpec.Description = toNeutronDescription(*portOpts.Description)
	}

	// Tags: combine server tags + per-port tags
	var orcTags []orcv1alpha1.NeutronTag
	for _, t := range serverTags {
		orcTags = append(orcTags, orcv1alpha1.NeutronTag(t))
	}
	for _, t := range portOpts.Tags {
		orcTags = append(orcTags, orcv1alpha1.NeutronTag(t))
	}
	portSpec.Tags = orcTags

	return &orcv1alpha1.Port{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PortORCName(serverName, index),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.PortSpec{
			ManagementPolicy:    orcv1alpha1.ManagementPolicyManaged,
			Resource:            portSpec,
			CloudCredentialsRef: credRef,
		},
	}
}

// buildTrunk builds a managed ORC Trunk for a trunk-enabled port.
func buildTrunk(serverName, namespace string, portIndex int, portORCName string, serverTags []string, credRef orcv1alpha1.CloudCredentialsReference) *orcv1alpha1.Trunk {
	var orcTags []orcv1alpha1.NeutronTag
	for _, t := range serverTags {
		orcTags = append(orcTags, orcv1alpha1.NeutronTag(t))
	}

	return &orcv1alpha1.Trunk{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TrunkORCName(serverName, portIndex),
			Namespace: namespace,
		},
		Spec: orcv1alpha1.TrunkSpec{
			ManagementPolicy: orcv1alpha1.ManagementPolicyManaged,
			Resource: &orcv1alpha1.TrunkResourceSpec{
				PortRef: orcv1alpha1.KubernetesNameRef(portORCName),
				Tags:    orcTags,
			},
			CloudCredentialsRef: credRef,
		},
	}
}

// Conversion helpers for port fields.

func convertPortSecurity(enabled *bool) orcv1alpha1.PortSecurityState {
	if enabled == nil {
		return orcv1alpha1.PortSecurityInherit
	}
	if *enabled {
		return orcv1alpha1.PortSecurityEnabled
	}
	return orcv1alpha1.PortSecurityDisabled
}

func convertAllowedAddressPairs(pairs []infrav1.AddressPair) []orcv1alpha1.AllowedAddressPair {
	if len(pairs) == 0 {
		return nil
	}
	result := make([]orcv1alpha1.AllowedAddressPair, len(pairs))
	for i, pair := range pairs {
		result[i] = orcv1alpha1.AllowedAddressPair{
			IP: orcv1alpha1.IPvAny(pair.IPAddress),
		}
		if pair.MACAddress != nil {
			mac := orcv1alpha1.MAC(*pair.MACAddress)
			result[i].MAC = &mac
		}
	}
	return result
}

func convertValueSpecs(specs []infrav1.ValueSpec) []orcv1alpha1.PortValueSpec {
	if len(specs) == 0 {
		return nil
	}
	result := make([]orcv1alpha1.PortValueSpec, len(specs))
	for i, vs := range specs {
		value := vs.Value
		result[i] = orcv1alpha1.PortValueSpec{
			Key:   vs.Key,
			Value: &value,
		}
	}
	return result
}
