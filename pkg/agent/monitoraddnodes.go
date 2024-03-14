package agent

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// MonitorAddNodes waits for the a node to be added to the cluster
// and reports its status until it becomes Ready.
func MonitorAddNodes(cluster *Cluster, nodeIPAddress string) error {
	timeout := 90 * time.Minute
	waitContext, cancel := context.WithTimeout(cluster.Ctx, timeout)
	defer cancel()

	var lastErrOnExit error
	var lastErrStr string
	wait.Until(func() {
		exitOnErr := false
		hasJoined, isReady, err := cluster.HasNodeJoinedClusterAndIsReady(nodeIPAddress)
		if err != nil {
			logrus.Info("error retrieving nodes from cluster")
		}

		hasCSRsPending, err := cluster.HasCSRs(nodeIPAddress)
		if err != nil {
			logrus.Info("error retrieving CSRs from cluster")
		}

		if hasCSRsPending {
			logrus.Info("Cluster has CSRs waiting for approval")
		}

		// TODO: take care of case where 1st CSR is approved
		// node becomes Ready and 2nd CSR is submitted for approval
		// could there be a state where the node is Ready and
		// the command sees there are no CSRs pending and exits before
		// the 2nd CSR is submitted for approval?
		if hasJoined && isReady && !hasCSRsPending && err == nil {
			cancel()
		}

		if err != nil {
			if exitOnErr {
				lastErrOnExit = err
				cancel()
			} else {
				if err.Error() != lastErrStr {
					logrus.Info(err)
					lastErrStr = err.Error()
				}
			}
		}
	}, 15*time.Second, waitContext.Done())

	waitErr := waitContext.Err()
	if waitErr != nil {
		if errors.Is(waitErr, context.Canceled) && lastErrOnExit != nil {
			return errors.Wrap(lastErrOnExit, "monitor-add-nodes process returned error")
		}
		if errors.Is(waitErr, context.DeadlineExceeded) {
			return errors.Wrap(waitErr, "monitor-add-nodes process timed out")
		}
	}

	return nil
}

// HasNodeJoinedClusterAndIsReady checks if the node specified by nodeIPAddress
// has joined the cluster and is in Ready state.
func (c *Cluster) HasNodeJoinedClusterAndIsReady(nodeIPAddress string) (bool, bool, error) {

	nodes, err := c.API.Kube.ListNodes()
	if err != nil {
		logrus.Infof("error getting node list %v", err)
		return false, false, nil
	}

	var joinedNode corev1.Node
	hasJoined := false
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				if address.Address == nodeIPAddress {
					joinedNode = node
					hasJoined = true
				}
			}
		}
	}

	isReady := false
	if hasJoined {
		logrus.Infof("Node %v (%s) has joined cluster", nodeIPAddress, joinedNode.Name)
		for _, cond := range joinedNode.Status.Conditions {
			if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
				isReady = true
			}
		}
		if isReady {
			logrus.Infof("Node %s (%s) is Ready", nodeIPAddress, joinedNode.Name)
		} else {
			logrus.Infof("Node %s (%s) is not Ready", nodeIPAddress, joinedNode.Name)
		}
	} else {
		logrus.Infof("Node %s has not joined cluster", nodeIPAddress)
	}

	return hasJoined, isReady, nil
}

func (c *Cluster) HasCSRs(nodeIPAddress string) (bool, error) {
	csrs, err := c.API.Kube.ListCSRs()
	if err != nil {
		return false, err
	}

	hasCSRsPending := false
	for _, csr := range csrs.Items {
		if len(csr.Status.Conditions) == 0 {
			// CSR is Pending and awaiting approval
			hasCSRsPending = true
			logrus.Infof("CSR %s with signerName %s and username %s is Pending and awaiting approval",
				csr.Name, csr.Spec.SignerName, csr.Spec.Username)
		}
	}

	return hasCSRsPending, nil
}
