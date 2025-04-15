/*
Copyright 2022 The Kubernetes Authors.

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

package converters

import (
	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20201101"
	"k8s.io/utils/ptr"
)

// GetSubnetAddresses returns the address prefixes contained in an ASO subnet.
func GetSubnetAddresses(subnet asonetworkv1.VirtualNetworksSubnet) []string {
	var addresses []string
	if subnet.Status.AddressPrefix != nil {
		addresses = []string{ptr.Deref(subnet.Status.AddressPrefix, "")}
	} else if subnet.Status.AddressPrefixes != nil {
		addresses = subnet.Status.AddressPrefixes
	}
	return addresses
}
