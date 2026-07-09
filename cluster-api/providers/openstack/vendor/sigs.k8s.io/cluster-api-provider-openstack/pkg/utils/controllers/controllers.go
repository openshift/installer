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

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// ValidateSubnets validates that all subnet CIDRs are parseable.
// Multiple subnets of the same IP version are allowed; use PrimarySubnet to
// specify which subnet should be used for load balancer VIP allocation and
// member registration when multiple subnets are present.
func ValidateSubnets(subnets []infrav1.Subnet) error {
	for _, subnet := range subnets {
		if _, _, err := net.ParseCIDR(subnet.CIDR); err != nil {
			return fmt.Errorf("invalid CIDR %q in subnet %q: %w", subnet.CIDR, subnet.ID, err)
		}
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
