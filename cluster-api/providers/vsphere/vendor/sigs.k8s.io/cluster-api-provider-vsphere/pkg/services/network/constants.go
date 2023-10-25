/*
Copyright 2021 The Kubernetes Authors.

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

package network

const (
	NSXTTypeNetwork = "nsx-t"
	// NSXTVNetSelectorKey
	// This constant is also defined in VM Operator.
	NSXTVNetSelectorKey = "ncp.vmware.com/virtual-network-name"

	CAPVDefaultNetworkLabel    = "capv.vmware.com/is-default-network"
	NetOpNetworkNameAnnotation = "netoperator.vmware.com/network-name"

	// SystemNamespace is the namespace where supervisor control plane VMs reside.
	SystemNamespace = "kube-system"

	// legacyDefaultNetworkLabel was the label used for default networks.
	// This is deprecated and is introduced only for smoother transitions.
	// This will be released in a future release.
	legacyDefaultNetworkLabel = "capw.vmware.com/is-default-network"
)
