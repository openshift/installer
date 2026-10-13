package rosa

import (
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	ocmerrors "github.com/openshift-online/ocm-sdk-go/errors"
	errors "github.com/zgalor/weberr"
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

func handleErr(res *ocmerrors.Error, err error) error {
	msg := res.Reason()
	if msg == "" {
		msg = err.Error()
	}
	// Hack to always display the correct terms and conditions message
	if res.Code() == "CLUSTERS-MGMT-451" {
		msg = "You must accept the Terms and Conditions in order to continue.\n" +
			"Go to https://www.redhat.com/wapps/tnc/ackrequired?site=ocm&event=register\n" +
			"Once you accept the terms, you will need to retry the action that was blocked."
	}
	errType := errors.ErrorType(res.Status()) //#nosec G115
	return errType.Set(errors.Errorf("%s", msg))
}
