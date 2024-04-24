package agent

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	firstCSRSignerName  = "kubernetes.io/kube-apiserver-client-kubelet"
	secondCSRSignerName = "kubernetes.io/kubelet-serving"
)

type addNodeStatusHistory struct {
	RestAPISeen            bool
	KubeletIsRunningOnNode bool
	FirstCSRSeen           bool
	SecondCSRSeen          bool
	NodeJoinedCluster      bool
	NodeIsReady            bool
}

type addNodeMonitor struct {
	nodeIPAddress string
	cluster       *Cluster
	status        addNodeStatusHistory
}

func newAddNodeMonitor(nodeIP string, cluster *Cluster) *addNodeMonitor {
	mon := addNodeMonitor{
		nodeIPAddress: nodeIP,
		cluster:       cluster,
		status: addNodeStatusHistory{
			RestAPISeen:            false,
			KubeletIsRunningOnNode: false,
			FirstCSRSeen:           false,
			SecondCSRSeen:          false,
			NodeJoinedCluster:      false,
			NodeIsReady:            false,
		},
	}
	return &mon
}

func (mon *addNodeMonitor) logStatus(status string) {
	logrus.Infof("Node %s: %s", mon.nodeIPAddress, status)
}

// MonitorAddNodes waits for the a node to be added to the cluster
// and reports its status until it becomes Ready.
func MonitorAddNodes(cluster *Cluster, nodeIPAddress string) error {
	parsedIPAddress := net.ParseIP(nodeIPAddress)
	if parsedIPAddress == nil {
		return fmt.Errorf("%s is not valid IP Address", nodeIPAddress)
	}

	timeout := 90 * time.Minute
	waitContext, cancel := context.WithTimeout(cluster.Ctx, timeout)
	defer cancel()

	mon := newAddNodeMonitor(nodeIPAddress, cluster)

	wait.Until(func() {
		if !mon.status.RestAPISeen &&
			mon.cluster.API.Rest.IsRestAPILive() {
			mon.status.RestAPISeen = true
			mon.logStatus("Assisted Service API is available")
		}

		if !mon.status.KubeletIsRunningOnNode &&
			mon.isKubeletRunningOnNode() {
			mon.status.KubeletIsRunningOnNode = true
			mon.logStatus("Kubelet is running")
		}

		if mon.status.KubeletIsRunningOnNode &&
			!mon.status.FirstCSRSeen &&
			mon.clusterHasFirstCSRPending() {
			mon.status.FirstCSRSeen = true
			mon.logStatus("First CSR Pending approval")
			mon.logCSRsPendingApproval(firstCSRSignerName)
		}

		if mon.status.KubeletIsRunningOnNode &&
			!mon.status.SecondCSRSeen &&
			mon.clusterHasSecondCSRPending() {
			mon.status.SecondCSRSeen = true
			mon.logStatus("Second CSR Pending approval")
			mon.logCSRsPendingApproval(secondCSRSignerName)
		}

		hasJoined, isReady, err := mon.nodeHasJoinedClusterAndIsReady()
		if err != nil {
			logrus.Debugf("nodeHasJoinedClusterAndIsReady returned err: %v", err)
		}

		if !mon.status.NodeJoinedCluster && hasJoined {
			mon.status.NodeJoinedCluster = true
			mon.logStatus("Node joined cluster")
		}

		if !mon.status.NodeIsReady && isReady {
			mon.status.NodeIsReady = true
			mon.logStatus("Node is Ready")
			// TODO: There appears to be a bug where the node becomes Ready
			// before second CSR is approved. Log Pending CSRs for now, so users
			// are aware there are still some waiting their approval even
			// though the node status is Ready.
			mon.logCSRsPendingApproval(secondCSRSignerName)
			cancel()
		}

		if mon.cluster.API.Rest.IsRestAPILive() {
			_, err = cluster.MonitorStatusFromAssistedService()
			if err != nil {
				logrus.Warnf("Node %s: %s", nodeIPAddress, err)
			}
		}
	}, 5*time.Second, waitContext.Done())

	waitErr := waitContext.Err()
	if waitErr != nil {
		if errors.Is(waitErr, context.Canceled) {
			cancel()
		}
		if errors.Is(waitErr, context.DeadlineExceeded) {
			return errors.Wrap(waitErr, "monitor-add-nodes process timed out")
		}
	}

	return nil
}

func (mon *addNodeMonitor) nodeHasJoinedClusterAndIsReady() (bool, bool, error) {
	nodes, err := mon.cluster.API.Kube.ListNodes()
	if err != nil {
		logrus.Debugf("error getting node list %v", err)
		return false, false, nil
	}

	var joinedNode corev1.Node
	hasJoined := false
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				if address.Address == mon.nodeIPAddress {
					joinedNode = node
					hasJoined = true
				}
			}
		}
	}

	isReady := false
	if hasJoined {
		logrus.Debugf("Node %v (%s) has joined cluster", mon.nodeIPAddress, joinedNode.Name)
		for _, cond := range joinedNode.Status.Conditions {
			if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
				isReady = true
			}
		}
		if isReady {
			logrus.Debugf("Node %s (%s) is Ready", mon.nodeIPAddress, joinedNode.Name)
		} else {
			logrus.Debugf("Node %s (%s) is not Ready", mon.nodeIPAddress, joinedNode.Name)
		}
	} else {
		logrus.Debugf("Node %s has not joined cluster", mon.nodeIPAddress)
	}

	return hasJoined, isReady, nil
}

func (mon *addNodeMonitor) logCSRsPendingApproval(signerName string) {
	// TODO? The CSRs have no IP address to identify for which
	// node it is for, so it is possible to log CSRs pending for
	// other nodes.
	csrs, err := mon.cluster.API.Kube.ListCSRs()
	if err != nil {
		logrus.Debugf("error calling listCSRs(): %v ", err)
	}

	for _, csr := range csrs.Items {
		if len(csr.Status.Conditions) == 0 {
			// CSR is Pending and awaiting approval
			if signerName != "" && signerName != csr.Spec.SignerName {
				continue
			}

			logrus.Infof("CSR %s with signerName %s and username %s is Pending and awaiting approval",
				csr.Name, csr.Spec.SignerName, csr.Spec.Username)
		}
	}
}

func (mon *addNodeMonitor) clusterHasFirstCSRPending() bool {
	return len(mon.cluster.API.Kube.getCSRsPendingApproval(firstCSRSignerName)) > 0
}

func (mon *addNodeMonitor) clusterHasSecondCSRPending() bool {
	// TODO: the csr.Spec.Username contains the node name
	// can we use it as an additional filter to only show
	// those matching mon.nodeIPAddress?
	return len(mon.cluster.API.Kube.getCSRsPendingApproval(secondCSRSignerName)) > 0
}

// isKubeletRunningOnNode checks if kubelet responds
// to http. Even if kubelet responds with error like
// TLS errors, kubelet is considered running.
func (mon *addNodeMonitor) isKubeletRunningOnNode() bool {
	url := fmt.Sprintf("https://%s:10250/metrics", mon.nodeIPAddress)
	// http get without authentication
	resp, err := http.Get(url) //nolint mon.nodeIPAddress is prevalidated to be IP address
	if err != nil {
		logrus.Debugf("kubelet http err: %v", err)
		if strings.Contains(err.Error(), "remote error: tls: internal error") {
			// nodes being added will return this error
			return true
		}
		if strings.Contains(err.Error(), "tls: failed to verify certificate: x509: certificate signed by unknown authority") {
			// existing control plane nodes returns this error
			return true
		}
		if strings.Contains(err.Error(), "connect: no route to host") {
			return false
		}
	} else {
		logrus.Debugf("kubelet http status code: %v", resp.StatusCode)
	}
	return false
}
