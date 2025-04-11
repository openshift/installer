// Package utils has the common functions that can be used for cluster-api-provider-aws repo.
package utils

import (
	"context"
	"fmt"

	crclient "sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

// GetMachinePools belong to a cluster.
func GetMachinePools(ctx context.Context, client crclient.Client, clusterName string, clusterNS string) ([]expclusterv1.MachinePool, error) {
	machinePoolList := expclusterv1.MachinePoolList{}
	listOptions := []crclient.ListOption{
		crclient.InNamespace(clusterNS),
		crclient.MatchingLabels(map[string]string{clusterv1.ClusterNameLabel: clusterName}),
	}

	if err := client.List(ctx, &machinePoolList, listOptions...); err != nil {
		return []expclusterv1.MachinePool{}, fmt.Errorf("failed to list machine pools for cluster %s: %v", clusterName, err)
	}

	return machinePoolList.Items, nil
}
