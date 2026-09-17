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

package powervs

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

// GetClusterByName finds and return a Cluster object using the specified params.
func GetClusterByName(ctx context.Context, c client.Client, namespace, name string) (*infrav1.IBMPowerVSCluster, error) {
	cluster := &infrav1.IBMPowerVSCluster{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}

	if err := c.Get(ctx, key, cluster); err != nil {
		return nil, fmt.Errorf("failed to get Cluster/%s: %w", name, err)
	}

	return cluster, nil
}

func fetchBucketRegion(cos infrav1.COSInstanceSource, vpc infrav1.VPCStatus) string {
	if cos.BucketRegion != "" {
		return cos.BucketRegion
	}
	return vpc.Region
}

// normalizedVPCSecurityGroupRulePrototype translates the deprecated 'all' protocol value to
// 'icmp_tcp_udp'. This ensures backward compatibility with YAMLs using the old value.
// The prototype is returned unchanged unless it uses 'all', in which case a normalized copy
// is returned.
func normalizedVPCSecurityGroupRulePrototype(prototype *infrav1.VPCSecurityGroupRulePrototype) *infrav1.VPCSecurityGroupRulePrototype {
	if prototype == nil || prototype.Protocol != infrav1.VPCSecurityGroupRuleProtocolAll {
		return prototype
	}
	normalized := prototype.DeepCopy()
	normalized.Protocol = infrav1.VPCSecurityGroupRuleProtocolIcmpTCPUDP
	return normalized
}
