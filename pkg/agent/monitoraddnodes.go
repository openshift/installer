package agent

import (
	"context"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	gatherssh "github.com/openshift/installer/pkg/gather/ssh"
)

type addNodeState int

// The node states that we observe and report.
const (
	InitialState addNodeState = iota
	WaitingForInitialBoot
	NodeBootedAndSSHIsUp
	WaitingForAssistedService
	AssistedServiceIsUp
	AssistedServiceFailedToStart // terminal state
	ClusterIsImported
	ClusterImportedFailed // terminal state
	InfraenvIsRegistered
	InfraenvRegisterFailed // terminal state
	HostConfigsApplied
	HostConfigsApplyFailed // terminal state
	WaitingForHostToBeReady
	HostHasErrors // terminal state, host does not pass validation or has some other error
	HostInstallationStarted
	CoreosInstallerWritingToDisk
	HostPreparingForFirstReboot
	AfterFirstReboot
	AfterSecondReboot
	AfterSecondRebootKubeletIsUp
	FirstCSRPendingApproval
	SecondCSRPendingApproval
	NodeJoinedClusterStatusNotReady
	Finish // terminal state
)

var stateNameMap = map[int]string{
	int(InitialState):                    "Initial State",
	int(WaitingForInitialBoot):           "Waiting for node to boot",
	int(NodeBootedAndSSHIsUp):            "Node booted and SSH is up",
	int(WaitingForAssistedService):       "Waiting for Assisted Service to start",
	int(AssistedServiceIsUp):             "Assisted Service API is available",
	int(AssistedServiceFailedToStart):    "Assisted Service API failed to start",
	int(ClusterIsImported):               "Cluster imported into Assisted Service",
	int(ClusterImportedFailed):           "Cluster import failed, check agent-import-cluster.service",
	int(InfraenvIsRegistered):            "InfraEnv registered into Assisted Service",
	int(InfraenvRegisterFailed):          "Failed to register InfraEnv into Assisted Service",
	int(HostConfigsApplied):              "Host configurations (if any) applied",
	int(HostConfigsApplyFailed):          "Failed to apply host configurations",
	int(WaitingForHostToBeReady):         "Waiting for host validations to complete",
	int(HostHasErrors):                   "Host has errors, node cannot be added until errors are corrected",
	int(HostInstallationStarted):         "Host installation started",
	int(CoreosInstallerWritingToDisk):    "Installer is initializing disk",
	int(HostPreparingForFirstReboot):     "Host preparing to reboot",
	int(AfterFirstReboot):                "Node completed reboot 1 of 2",
	int(AfterSecondReboot):               "Node completed reboot 2 of 2",
	int(AfterSecondRebootKubeletIsUp):    "Kubelet is running",
	int(FirstCSRPendingApproval):         "1st CSR Pending approval",
	int(SecondCSRPendingApproval):        "2nd CSR Pending approval",
	int(NodeJoinedClusterStatusNotReady): "Node joined cluster and status is NotReady",
	int(Finish):                          "Node is Ready",
}

type checkFunction func(nodeMonitor *addNodeMonitor)

type checkFunctions map[int]checkFunction

var defaultChecks = checkFunctions{
	int(WaitingForInitialBoot): func(mon *addNodeMonitor) {
		if mon.in(InitialState) &&
			!mon.canSSHToNode() {
			mon.setState(WaitingForInitialBoot)
		}
	},
	int(NodeBootedAndSSHIsUp): func(mon *addNodeMonitor) {
		if mon.before(NodeBootedAndSSHIsUp) &&
			mon.canSSHToNode() {
			mon.setState(NodeBootedAndSSHIsUp)
		}
	},
	int(WaitingForAssistedService): func(mon *addNodeMonitor) {
		if mon.before(WaitingForAssistedService) &&
			mon.nodeCommandOutputContains("ls -l /etc/systemd/system/assisted-service.service", "assisted-service.service") {
			mon.setState(WaitingForAssistedService)
		}
	},
	int(AssistedServiceIsUp): func(mon *addNodeMonitor) {
		if mon.in(WaitingForAssistedService) &&
			mon.nodeCommandOutputContains("sudo podman ps -a", "assisted-service") {
			mon.setState(AssistedServiceIsUp)
		}
	},
	int(AssistedServiceFailedToStart): func(mon *addNodeMonitor) {
		if mon.in(WaitingForAssistedService) &&
			mon.nodeCommandOutputContains("systemctl status assisted-service", "FAILURE") {
			mon.setState(AssistedServiceFailedToStart)
		}
	},
	int(ClusterIsImported): func(mon *addNodeMonitor) {
		if mon.in(AssistedServiceIsUp) &&
			mon.nodeCommandOutputContains("systemctl status agent-import-cluster", "Imported cluster with id:") {
			mon.setState(ClusterIsImported)
		}
	},
	int(ClusterImportedFailed): func(mon *addNodeMonitor) {
		// Can simulate this by attempting to add a node using version 4.15 where
		// the ABI CLI does not have the importcluster subcommand
		if mon.in(AssistedServiceIsUp) &&
			mon.nodeCommandOutputContains("systemctl status agent-import-cluster", "FAILURE") {
			mon.setState(ClusterImportedFailed)
		}
	},
	int(InfraenvIsRegistered): func(mon *addNodeMonitor) {
		if mon.in(ClusterIsImported) &&
			mon.nodeCommandOutputContains("systemctl status agent-register-infraenv", "Registered infraenv with id:") {
			mon.setState(InfraenvIsRegistered)
		}
	},
	int(InfraenvRegisterFailed): func(mon *addNodeMonitor) {
		if mon.in(ClusterIsImported) &&
			mon.nodeCommandOutputContains("systemctl status agent-register-infraenv", "FAILURE") {
			mon.setState(InfraenvRegisterFailed)
		}
	},
	int(HostConfigsApplied): func(mon *addNodeMonitor) {
		if mon.in(InfraenvIsRegistered) &&
			mon.nodeCommandOutputContains("systemctl status apply-host-config", "Finished Service") {
			mon.setState(HostConfigsApplied)
		}
	},
	int(HostConfigsApplyFailed): func(mon *addNodeMonitor) {
		if mon.in(InfraenvIsRegistered) &&
			mon.nodeCommandOutputContains("systemctl status apply-host-config", "FAILURE") {
			mon.setState(HostConfigsApplyFailed)
		}
	},
	int(WaitingForHostToBeReady): func(mon *addNodeMonitor) {
		if mon.in(HostConfigsApplied) &&
			mon.nodeCommandOutputContains("sudo podman exec assisted-db psql -d installer -c 'select id, status_info, status from hosts'", "insufficient") {
			mon.setState(WaitingForHostToBeReady)
		}
	},
	int(HostHasErrors): func(mon *addNodeMonitor) {
		if mon.in(WaitingForHostToBeReady) &&
			mon.nodeCommandOutputContains("systemctl status agent-add-node", "FAILURE") {
			mon.setState(HostHasErrors)
		}
	},
	int(HostInstallationStarted): func(mon *addNodeMonitor) {
		if mon.in(WaitingForHostToBeReady) &&
			// Or
			// sudo podman exec assisted-db psql -d installer -c 'select id, status_info, status from hosts'
			// and status will show "installing"
			mon.nodeCommandOutputContains("systemctl status agent-add-node", "Host installation started") {
			mon.setState(HostInstallationStarted)
		}
	},
	int(CoreosInstallerWritingToDisk): func(mon *addNodeMonitor) {
		if mon.in(HostInstallationStarted) &&
			mon.nodeCommandOutputContains("ps -ef | grep coreos-installer", "coreos-installer") {
			mon.setState(CoreosInstallerWritingToDisk)
		}
	},
	int(HostPreparingForFirstReboot): func(mon *addNodeMonitor) {
		if mon.in(CoreosInstallerWritingToDisk) &&
			!mon.cluster.CanSSHToNodeZero() {
			mon.setState(HostPreparingForFirstReboot)
		}
	},
	int(AfterFirstReboot): func(mon *addNodeMonitor) {
		if mon.before(AfterFirstReboot) &&
			mon.nodeCommandOutputContains("ps -ef", "ostree") {
			mon.setState(AfterFirstReboot)
		}
	},
	int(AfterSecondReboot): func(mon *addNodeMonitor) {
		if mon.before(AfterSecondRebootKubeletIsUp) &&
			mon.nodeCommandOutputContains("ps -ef", "coredns") {
			mon.setState(AfterSecondReboot)
		}
	},
	int(AfterSecondRebootKubeletIsUp): func(mon *addNodeMonitor) {
		if mon.onOrAfter(AfterSecondReboot) &&
			mon.before(FirstCSRPendingApproval) &&
			mon.nodeCommandOutputContains("ps -ef", "kubelet") {
			mon.setState(AfterSecondRebootKubeletIsUp)
		}
	},
	int(FirstCSRPendingApproval): func(mon *addNodeMonitor) {
		if mon.onOrAfter(AfterSecondReboot) &&
			mon.onOrBefore(FirstCSRPendingApproval) &&
			mon.nodeCommandOutputContains("sudo KUBECONFIG=/etc/kubernetes/kubeconfig oc get csr | grep Pending", "node-bootstrapper") {
			mon.setState(FirstCSRPendingApproval)
			mon.logCSRsPendingApproval()
		}
	},
	int(SecondCSRPendingApproval): func(mon *addNodeMonitor) {
		if mon.onOrAfter(AfterSecondReboot) &&
			mon.onOrBefore(SecondCSRPendingApproval) &&
			mon.nodeCommandOutputContains("sudo KUBECONFIG=/etc/kubernetes/kubeconfig oc get csr | grep Pending", "system:node") {
			mon.setState(SecondCSRPendingApproval)
			mon.logCSRsPendingApproval()
		}
	},
	int(NodeJoinedClusterStatusNotReady): func(mon *addNodeMonitor) {
		if mon.in(AfterSecondRebootKubeletIsUp) ||
			mon.in(SecondCSRPendingApproval) {
			hasJoined, isReady, err := mon.nodeHasJoinedClusterAndIsReady()
			if err != nil {
				logrus.Debugf("HasNodeJoinedClusterAndIsReady returned err: %v", err)
			}
			if hasJoined && !isReady {
				mon.setState(NodeJoinedClusterStatusNotReady)
			}
		}
	},
	int(Finish): func(mon *addNodeMonitor) {
		if mon.in(AfterSecondRebootKubeletIsUp) ||
			mon.in(SecondCSRPendingApproval) ||
			mon.in(NodeJoinedClusterStatusNotReady) {
			hasJoined, isReady, err := mon.nodeHasJoinedClusterAndIsReady()
			if err != nil {
				logrus.Debugf("HasNodeJoinedClusterAndIsReady returned err: %v", err)
			}
			if hasJoined && isReady {
				mon.setState(Finish)
			}
		}
	},
}

// MonitorAddNodes waits for the a node to be added to the cluster
// and reports its status until it becomes Ready.
func MonitorAddNodes(cluster *Cluster, nodeIPAddress string) error {
	timeout := 90 * time.Minute
	waitContext, cancel := context.WithTimeout(cluster.Ctx, timeout)
	defer cancel()

	mon := newAddNodeMonitor(nodeIPAddress, cluster)

	wait.Until(func() {
		if mon.onOrAfter(AssistedServiceIsUp) &&
			mon.onOrBefore(HostInstallationStarted) {
			_, _, err := cluster.LogAssistedServiceStatus()
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

		if mon.in(Finish) {
			// TODO: There appears to be a bug where the node becomes Ready
			// before second CSR is approved. Log Pending CSRs for now, so users
			// are aware there are still some waiting their approval even
			// though the node status is Ready.
			mon.logCSRsPendingApproval()
			cancel()
		}

		if mon.in(AssistedServiceFailedToStart) ||
			mon.in(ClusterImportedFailed) ||
			mon.in(InfraenvRegisterFailed) ||
			mon.in(HostConfigsApplyFailed) ||
			mon.in(HostHasErrors) {
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

func (mon *addNodeMonitor) before(state addNodeState) bool {
	return mon.currentState < state
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

func (mon *addNodeMonitor) canSSHToNode() bool {
	return mon.cluster.CanSSHToNodeZero()
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

func (mon *addNodeMonitor) sshClient() (*ssh.Client, error) {
	ip := mon.nodeIPAddress
	port := 22

	client, err := gatherssh.NewClient("core", net.JoinHostPort(ip, strconv.Itoa(port)), mon.cluster.API.Rest.NodeSSHKey)
	if err != nil {
		logrus.Debugf("Failed to create SSH client to %s: %s", ip, err)
	}
	return client, err
}

func (mon *addNodeMonitor) nodeCommandOutputContains(command, stringToMatchInOutput string) bool {
	sshClient, err := mon.sshClient()
	if err != nil {
		logrus.Debugf("Failed to create SSH client to %s: %s", mon.nodeIPAddress, err)
		return false
	}

	logrus.Debugf("sshCommandOutputContains command: %v ", command)

	// Create an SSH session
	session, err := sshClient.NewSession()
	if err != nil {
		logrus.Debugf("failed to create session: %s", err)
		return false
	}
	defer session.Close()

	// Run the command on the remote machine
	stdout, err := session.StdoutPipe()
	if err != nil {
		logrus.Debugf("failed to get stdout: %s", err)
		return false
	}

	if err := session.Start(command); err != nil {
		logrus.Debugf("failed to start command: %s", err)
		return false
	}

	output, err := io.ReadAll(stdout)
	if err != nil {
		logrus.Debugf("failed to read stdout: %s", err)
		return false
	}

	if err := session.Wait(); err != nil {
		if strings.Contains(err.Error(), "Process exited with status") {
			// command finished with exit code other than 0
			// which could happen in the case where "systemctl status"
			// shows service is in failed state and exits with code 3.
		} else {
			logrus.Debugf("command failed: %s output: %v", err, string(output))
			return false
		}
	}

	logrus.Debugf("sshCommandOutputContains command: %v output: \n %v", command, string(output))

	return strings.Contains(string(output), stringToMatchInOutput)
}
