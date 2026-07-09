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
	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// PortResolution contains the deduplicated ORC objects and name maps
// produced by resolving CAPO PortOpts into ORC resources.
type PortResolution struct {
	// Networks is the deduplicated list of unmanaged ORC Network objects to create.
	Networks []*orcv1alpha1.Network
	// Subnets is the deduplicated list of unmanaged ORC Subnet objects to create.
	Subnets []*orcv1alpha1.Subnet
	// SecurityGroups is the deduplicated list of unmanaged ORC SecurityGroup objects to create.
	SecurityGroups []*orcv1alpha1.SecurityGroup

	// NetworkNameMap maps a NetworkParam key to the ORC Network object name.
	NetworkNameMap map[string]string
	// SubnetNameMap maps a SubnetParam key to the ORC Subnet object name.
	SubnetNameMap map[string]string
	// SGNameMap maps a SecurityGroupParam key to the ORC SecurityGroup object name.
	SGNameMap map[string]string
}

// ResolvePortsToORC translates CAPO PortOpts into deduplicated ORC
// Network, Subnet, and SecurityGroup objects, along with name maps that
// the port builder uses to reference them.
//
// This replaces the existing networking.ConstructPorts() flow: instead of
// resolving references to OpenStack IDs via Gophercloud, we create
// unmanaged ORC objects that handle the resolution.
func ResolvePortsToORC(
	serverName, namespace string,
	ports []infrav1.PortOpts,
	globalSGs []infrav1.SecurityGroupParam,
	credRef orcv1alpha1.CloudCredentialsReference,
) *PortResolution {
	res := &PortResolution{
		NetworkNameMap: make(map[string]string),
		SubnetNameMap:  make(map[string]string),
		SGNameMap:      make(map[string]string),
	}

	// Collect and deduplicate networks
	for _, port := range ports {
		if port.Network != nil {
			addNetwork(res, serverName, namespace, *port.Network, credRef)
		}
	}

	// Collect and deduplicate subnets from fixed IPs
	for _, port := range ports {
		for _, fixedIP := range port.FixedIPs {
			if fixedIP.Subnet != nil {
				addSubnet(res, serverName, namespace, *fixedIP.Subnet, credRef)
			}
		}
	}

	// Collect and deduplicate security groups (per-port + global)
	for _, port := range ports {
		for _, sg := range port.SecurityGroups {
			addSecurityGroup(res, serverName, namespace, sg, credRef)
		}
	}
	for _, sg := range globalSGs {
		addSecurityGroup(res, serverName, namespace, sg, credRef)
	}

	return res
}

func addNetwork(res *PortResolution, serverName, namespace string, param infrav1.NetworkParam, credRef orcv1alpha1.CloudCredentialsReference) {
	key := NetworkParamKey(param)
	if key == "" {
		return
	}
	if _, exists := res.NetworkNameMap[key]; exists {
		return // deduplicated
	}
	obj := buildNetwork(serverName, namespace, param, credRef)
	res.Networks = append(res.Networks, obj)
	res.NetworkNameMap[key] = obj.Name
}

func addSubnet(res *PortResolution, serverName, namespace string, param infrav1.SubnetParam, credRef orcv1alpha1.CloudCredentialsReference) {
	key := SubnetParamKey(param)
	if key == "" {
		return
	}
	if _, exists := res.SubnetNameMap[key]; exists {
		return // deduplicated
	}
	obj := buildSubnet(serverName, namespace, param, credRef)
	res.Subnets = append(res.Subnets, obj)
	res.SubnetNameMap[key] = obj.Name
}

func addSecurityGroup(res *PortResolution, serverName, namespace string, param infrav1.SecurityGroupParam, credRef orcv1alpha1.CloudCredentialsReference) {
	key := SecurityGroupParamKey(param)
	if key == "" {
		return
	}
	if _, exists := res.SGNameMap[key]; exists {
		return // deduplicated
	}
	obj := buildSecurityGroup(serverName, namespace, param, credRef)
	res.SecurityGroups = append(res.SecurityGroups, obj)
	res.SGNameMap[key] = obj.Name
}
