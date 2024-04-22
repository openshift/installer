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

type addNodeState int

// The node states that we observe and report.
const (
	InitialState addNodeState = iota
	AssistedServiceIsUp
	KubeletIsRunningOnNode
	FirstCSRPendingApproval
	SecondCSRPendingApproval
	NodeJoinedClusterStatusNotReady
	Finish // terminal state
)

var stateNameMap = map[int]string{
	int(InitialState):                    "Initial State",
	int(AssistedServiceIsUp):             "Assisted Service API is available",
	int(KubeletIsRunningOnNode):          "Kubelet is running",
	int(FirstCSRPendingApproval):         "First CSR Pending approval",
	int(SecondCSRPendingApproval):        "Second CSR Pending approval",
	int(NodeJoinedClusterStatusNotReady): "Node joined cluster and status is NotReady",
	int(Finish):                          "Node is Ready",
}

type checkFunction func(nodeMonitor *addNodeMonitor)

type checkFunctions map[int]checkFunction

var defaultChecks = checkFunctions{
	int(AssistedServiceIsUp): func(mon *addNodeMonitor) {
		if mon.in(InitialState) &&
			mon.cluster.API.Rest.IsRestAPILive() {
			mon.setState(AssistedServiceIsUp)
		}
	},
	int(KubeletIsRunningOnNode): func(mon *addNodeMonitor) {
		if mon.onOrBefore(KubeletIsRunningOnNode) &&
			mon.isKubeletRunningOnNode() {
			mon.setState(KubeletIsRunningOnNode)
		}
	},
	int(FirstCSRPendingApproval): func(mon *addNodeMonitor) {
		if mon.onOrAfter(KubeletIsRunningOnNode) &&
			mon.clusterHasFirstCSRPending() {
			mon.setState(FirstCSRPendingApproval)
		}
	},
	int(SecondCSRPendingApproval): func(mon *addNodeMonitor) {
		if mon.onOrAfter(KubeletIsRunningOnNode) &&
			mon.clusterHasSecondCSRPending() {
			mon.setState(SecondCSRPendingApproval)
		}
	},
	int(NodeJoinedClusterStatusNotReady): func(mon *addNodeMonitor) {
		hasJoined, isReady, err := mon.nodeHasJoinedClusterAndIsReady()
		if err != nil {
			logrus.Debugf("HasNodeJoinedClusterAndIsReady returned err: %v", err)
		}
		if hasJoined && !isReady {
			mon.setState(NodeJoinedClusterStatusNotReady)
		}
	},
	int(Finish): func(mon *addNodeMonitor) {
		hasJoined, isReady, err := mon.nodeHasJoinedClusterAndIsReady()
		if err != nil {
			logrus.Debugf("HasNodeJoinedClusterAndIsReady returned err: %v", err)
		}
		if hasJoined && isReady {
			mon.setState(Finish)
		}
	},
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
		if mon.cluster.API.Rest.IsRestAPILive() {
			_, err := cluster.MonitorStatusFromAssistedService()
			if err != nil {
				logrus.Warnf("Node %s: %s", nodeIPAddress, err)
			}
		}

		for _, checkFunc := range defaultChecks {
			lastState := mon.currentState
			checkFunc(mon)
			if !mon.in(lastState) {
				// log state change
				logrus.Infof("Node %s: %s", nodeIPAddress, mon.currentStateName())
			}
		}

		if mon.clusterHasFirstCSRPending() || mon.clusterHasSecondCSRPending() {
			// TODO? The CSRs have no IP address to identify for which
			// node it is for, so it is possible to log CSRs pending for
			// other nodes.
			mon.logCSRsPendingApproval()
		}

		if mon.in(Finish) {
			// TODO: There appears to be a bug where the node becomes Ready
			// before second CSR is approved. Log Pending CSRs for now, so users
			// are aware there are still some waiting their approval even
			// though the node status is Ready.
			mon.logCSRsPendingApproval()
			cancel()
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

type addNodeMonitor struct {
	nodeIPAddress string
	cluster       *Cluster
	currentState  addNodeState
}

func newAddNodeMonitor(nodeIP string, cluster *Cluster) *addNodeMonitor {
	mon := addNodeMonitor{
		nodeIPAddress: nodeIP,
		cluster:       cluster,
		currentState:  InitialState,
	}
	return &mon
}

func (mon *addNodeMonitor) setState(state addNodeState) {
	mon.currentState = state
}

func (mon *addNodeMonitor) in(state addNodeState) bool {
	return mon.currentState == state
}

func (mon *addNodeMonitor) onOrBefore(state addNodeState) bool {
	return mon.currentState <= state
}

func (mon *addNodeMonitor) onOrAfter(state addNodeState) bool {
	return mon.currentState >= state
}

func (mon *addNodeMonitor) currentStateName() string {
	return stateNameMap[int(mon.currentState)]
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

func (mon *addNodeMonitor) logCSRsPendingApproval() {
	csrs, err := mon.cluster.API.Kube.ListCSRs()
	if err != nil {
		logrus.Debugf("error calling listCSRs(): %v ", err)
	}

	for _, csr := range csrs.Items {
		if len(csr.Status.Conditions) == 0 {
			// CSR is Pending and awaiting approval
			logrus.Infof("CSR %s with signerName %s and username %s is Pending and awaiting approval",
				csr.Name, csr.Spec.SignerName, csr.Spec.Username)
		}
	}
}

func (mon *addNodeMonitor) clusterHasFirstCSRPending() bool {
	return len(mon.cluster.API.Kube.getCSRsPendingApproval("kubernetes.io/kube-apiserver-client-kubelet")) > 0
}

func (mon *addNodeMonitor) clusterHasSecondCSRPending() bool {
	return len(mon.cluster.API.Kube.getCSRsPendingApproval("kubernetes.io/kubelet-serving")) > 0
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
