/*
Copyright 2024 The Kubernetes Authors.

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

package util

import (
	"fmt"
	"net"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// ClusterIPFamily represents the IP topology of the cluster.
type ClusterIPFamily string

const (
	// IPv4SingleStack represents IPv4 single-stack IP topology.
	IPv4SingleStack ClusterIPFamily = "IPv4SingleStack"

	// IPv6SingleStack represents IPv6 single-stack IP topology.
	IPv6SingleStack ClusterIPFamily = "IPv6SingleStack"

	// DualStackIPv4Primary represents dual-stack IP topology with IPv4 as primary.
	DualStackIPv4Primary ClusterIPFamily = "DualStackIPv4Primary"

	// DualStackIPv6Primary represents dual-stack IP topology with IPv6 as primary.
	DualStackIPv6Primary ClusterIPFamily = "DualStackIPv6Primary"
)

// DetermineClusterIPFamily inspects the cluster network to determine the intended IP topology.
func DetermineClusterIPFamily(cluster *clusterv1.Cluster) (ClusterIPFamily, error) {
	if cluster == nil {
		return "", fmt.Errorf("cluster is nil")
	}
	cidrBlocks := cluster.Spec.ClusterNetwork.Pods.CIDRBlocks
	if len(cidrBlocks) == 0 {
		cidrBlocks = cluster.Spec.ClusterNetwork.Services.CIDRBlocks
	}
	if len(cidrBlocks) == 0 {
		return "", fmt.Errorf("cluster %s/%s has no Pod or Service CIDR blocks", cluster.Namespace, cluster.Name)
	}

	var hasIPv4, hasIPv6 bool
	var firstIsIPv6 bool
	for i, cidr := range cidrBlocks {
		ip, _, err := net.ParseCIDR(cidr)
		if err != nil {
			return "", fmt.Errorf("failed to parse CIDR %q: %w", cidr, err)
		}
		isIPv6 := ip.To4() == nil
		if isIPv6 {
			hasIPv6 = true
		} else {
			hasIPv4 = true
		}
		if i == 0 {
			firstIsIPv6 = isIPv6
		}
	}

	if hasIPv4 && hasIPv6 {
		if firstIsIPv6 {
			return DualStackIPv6Primary, nil
		}
		return DualStackIPv4Primary, nil
	}
	if hasIPv6 {
		return IPv6SingleStack, nil
	}
	return IPv4SingleStack, nil
}
