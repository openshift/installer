// Package utils has the common functions that can be used for cluster-api-provider-aws repo.
package utils

import (
	"context"
	"fmt"
	"math"

	"k8s.io/utils/ptr"
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

// ToInt64Pointer converts an int32 pointer to an int64 pointer.
func ToInt64Pointer(to *int32) *int64 {
	if to == nil {
		return nil
	}
	return ptr.To(int64(*to))
}

// ToInt32Pointer converts an int64 pointer to an int32 pointer.
func ToInt32Pointer(to *int64) *int32 {
	if to == nil {
		return nil
	}

	if *to > math.MaxInt32 || *to < math.MinInt32 {
		panic(fmt.Sprintf("Value %d is out of int32 range", *to))
	}

	return ptr.To(int32(*to))
}
