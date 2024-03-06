package rosa

import (
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// IsNodePoolReady checkes whether the nodepool is provisoned and all replicas are available.
// If autosacling is enabled, NodePool must have replicas >= autosacling.MinReplica to be considered ready.
func IsNodePoolReady(nodePool *cmv1.NodePool) bool {
	if nodePool.Status().Message() != "" {
		return false
	}

	if nodePool.Replicas() != 0 {
		return nodePool.Replicas() == nodePool.Status().CurrentReplicas()
	}

	if nodePool.Autoscaling() != nil {
		return nodePool.Status().CurrentReplicas() >= nodePool.Autoscaling().MinReplica()
	}

	return false
}
