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

package controllers

import (
	"context"
	"fmt"
	"net"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

// ValidateSubnets validates if the amount of IPv4 and IPv6 subnets is allowed by OpenStackCluster.
func ValidateSubnets(subnets []infrav1.Subnet) error {
	isIPv6 := []bool{false, false}
	for i, subnet := range subnets {
		ip, _, err := net.ParseCIDR(subnet.CIDR)
		if err != nil {
			return err
		}

		if ip.To4() == nil {
			isIPv6[i] = true
		}
	}

	if len(subnets) > 1 && isIPv6[0] == isIPv6[1] {
		ethertype := 4
		if isIPv6[0] {
			ethertype = 6
		}
		return fmt.Errorf("multiple IPv%d Subnet not allowed on OpenStackCluster", ethertype)
	}
	return nil
}

func GetInfraCluster(ctx context.Context, c client.Client, cluster *clusterv1.Cluster) (*infrav1.OpenStackCluster, error) {
	openStackCluster := &infrav1.OpenStackCluster{}
	openStackClusterName := client.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := c.Get(ctx, openStackClusterName, openStackCluster); err != nil {
		return nil, err
	}
	return openStackCluster, nil
}
