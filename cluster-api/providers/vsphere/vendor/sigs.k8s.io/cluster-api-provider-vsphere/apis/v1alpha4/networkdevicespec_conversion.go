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

package v1alpha4

import (
	conversion "k8s.io/apimachinery/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

func Convert_v1beta1_NetworkDeviceSpec_To_v1alpha4_NetworkDeviceSpec(in *infrav1.NetworkDeviceSpec, out *NetworkDeviceSpec, s conversion.Scope) error {
	out.NetworkName = in.NetworkName
	out.DeviceName = in.DeviceName
	out.DHCP4 = in.DHCP4
	out.DHCP6 = in.DHCP6
	out.Gateway4 = in.Gateway4
	out.Gateway6 = in.Gateway6
	out.IPAddrs = in.IPAddrs
	out.MTU = in.MTU
	out.MACAddr = in.MACAddr
	out.Nameservers = in.Nameservers
	out.SearchDomains = in.SearchDomains
	if in.Routes != nil {
		inRoutes, outRoutes := &in.Routes, &out.Routes
		*outRoutes = make([]NetworkRouteSpec, len(*inRoutes))
		for i := range *inRoutes {
			if err := Convert_v1beta1_NetworkRouteSpec_To_v1alpha4_NetworkRouteSpec(&(*inRoutes)[i], &(*outRoutes)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Routes = nil
	}
	return nil
}
