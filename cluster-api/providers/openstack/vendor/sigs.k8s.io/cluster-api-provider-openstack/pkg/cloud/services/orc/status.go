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
	corev1 "k8s.io/api/core/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// MapInstanceState converts an ORC/Nova server status string to a CAPO
// InstanceState.
func MapInstanceState(status string) infrav1.InstanceState {
	switch status {
	case "ACTIVE":
		return infrav1.InstanceStateActive
	case "BUILD":
		return infrav1.InstanceStateBuild
	case "ERROR":
		return infrav1.InstanceStateError
	case "DELETED":
		return infrav1.InstanceStateDeleted
	case "SHUTOFF":
		return infrav1.InstanceStateShutoff
	case "STOPPED":
		return infrav1.InstanceStateStopped
	case "SOFT_DELETED":
		return infrav1.InstanceStateSoftDeleted
	case "":
		return infrav1.InstanceStateUndefined
	default:
		return infrav1.InstanceState(status)
	}
}

// MapAddresses extracts IP addresses from ORC Server interface status
// and converts them to Kubernetes NodeAddresses. All fixed IPs are
// mapped as NodeInternalIP.
func MapAddresses(interfaces []orcv1alpha1.ServerInterfaceStatus) []corev1.NodeAddress {
	var addresses []corev1.NodeAddress
	for _, iface := range interfaces {
		for _, fixedIP := range iface.FixedIPs {
			if fixedIP.IPAddress != "" {
				addresses = append(addresses, corev1.NodeAddress{
					Type:    corev1.NodeInternalIP,
					Address: fixedIP.IPAddress,
				})
			}
		}
	}
	return addresses
}
