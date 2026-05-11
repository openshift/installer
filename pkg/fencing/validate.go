// Package fencing validates fencing on DualReplica (Two Node with Fencing) clusters
// by fencing both nodes sequentially and verifying recovery.
package fencing

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/openshift/installer/pkg/gather/ssh"
)

const (
	sshPort          = "22"
	pollInterval     = 5 * time.Second
	notReadyTimeout  = 20 * time.Minute
	readyTimeout     = 10 * time.Minute
	pcmkTimeout      = 10 * time.Minute
	etcdTimeout      = 10 * time.Minute
	fenceTimeout     = 600 // seconds, passed to timeout(1)
	fenceWarnSeconds = 60
)

var (
	stonithEnabledRe = regexp.MustCompile(`(?i)stonith-enabled\s*[:=]\s*true`)
	daemonActiveRe   = regexp.MustCompile(`(?i)active.*(running|enabled)`)
)

// NodeInfo holds resolved information about a cluster node.
type NodeInfo struct {
	Name          string
	PacemakerName string
	IP            string
}

// Config holds all inputs needed by the fencing validator.
type Config struct {
	KubeClient   kubernetes.Interface
	ConfigClient configclient.Interface
	SSHUser      string
	SSHKeys      []string
}

// Run executes the full fencing validation sequence.
func Run(ctx context.Context, cfg Config) error {
	nodes, err := discoverNodes(ctx, cfg)
	if err != nil {
		return fmt.Errorf("node discovery failed: %w", err)
	}

	if err := verifyTopology(ctx, cfg); err != nil {
		return err
	}

	logrus.Infof("Connecting to %s (%s)", nodes[0].Name, nodes[0].IP)
	clientA, err := sshConnect(cfg, nodes[0].IP)
	if err != nil {
		return fmt.Errorf("SSH to %s failed: %w", nodes[0].Name, err)
	}
	defer clientA.Close()

	if err := resolvePacemakerNames(clientA, nodes[:]); err != nil {
		return err
	}

	logrus.Info("Running pre-flight checks")
	if err := runPreFlight(clientA, nodes[:]); err != nil {
		return fmt.Errorf("pre-flight check failed: %w", err)
	}
	logrus.Info("Pre-flight checks passed")

	// Fence node B from node A
	if err := fenceAndRecover(ctx, cfg, clientA, nodes[:], 1); err != nil {
		return err
	}

	// Fence node A from node B — need new SSH connection to B
	logrus.Infof("Connecting to %s (%s)", nodes[1].Name, nodes[1].IP)
	clientB, err := sshConnect(cfg, nodes[1].IP)
	if err != nil {
		return fmt.Errorf("SSH to %s failed: %w", nodes[1].Name, err)
	}
	defer clientB.Close()

	if err := fenceAndRecover(ctx, cfg, clientB, nodes[:], 0); err != nil {
		return err
	}

	logrus.Info("Fencing validation passed")
	return nil
}

func fenceAndRecover(ctx context.Context, cfg Config, survivorClient *gossh.Client, nodes []NodeInfo, targetIdx int) error {
	target := nodes[targetIdx]
	survivor := nodes[1-targetIdx]

	logrus.Infof("Fencing %s (pacemaker: %s) from %s", target.Name, target.PacemakerName, survivor.Name)
	fenceStart := time.Now()

	if err := fenceNode(survivorClient, target.PacemakerName); err != nil {
		return fmt.Errorf("failed to fence %s: %w", target.Name, err)
	}

	logrus.Infof("Waiting for %s to become NotReady", target.Name)
	if err := waitNotReady(ctx, cfg.KubeClient, target.Name); err != nil {
		return fmt.Errorf("%s did not become NotReady: %w", target.Name, err)
	}

	fenceDuration := time.Since(fenceStart)
	if fenceDuration.Seconds() > fenceWarnSeconds {
		logrus.Warnf("Fencing %s took %s to power off (threshold is %ds). BMC may be performing graceful shutdown instead of power-off. Check BMC configuration.",
			target.Name, fenceDuration.Round(time.Second), fenceWarnSeconds)
	} else {
		logrus.Infof("Node %s powered off in %s", target.Name, fenceDuration.Round(time.Second))
	}

	logrus.Infof("Waiting for %s to become Ready", target.Name)
	if err := waitReady(ctx, cfg.KubeClient, target.Name); err != nil {
		return fmt.Errorf("%s did not become Ready: %w", target.Name, err)
	}

	logrus.Infof("Waiting for %s to rejoin Pacemaker", target.Name)
	if err := pollPacemakerOnline(ctx, survivorClient, nodes); err != nil {
		return fmt.Errorf("%s did not rejoin Pacemaker: %w", target.Name, err)
	}

	logrus.Info("Waiting for etcd quorum")
	if err := pollEtcdHealth(ctx, survivorClient, nodes); err != nil {
		return fmt.Errorf("etcd quorum not restored after fencing %s: %w", target.Name, err)
	}

	if err := checkDaemons(survivorClient); err != nil {
		return fmt.Errorf("daemon check failed after fencing %s: %w", target.Name, err)
	}

	logrus.Infof("Fencing %s: full recovery completed in %s", target.Name, time.Since(fenceStart).Round(time.Second))
	return nil
}

// --- Node discovery ---

func discoverNodes(ctx context.Context, cfg Config) ([2]NodeInfo, error) {
	nodeList, err := cfg.KubeClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/master=",
	})
	if err != nil {
		return [2]NodeInfo{}, fmt.Errorf("listing control-plane nodes: %w", err)
	}

	var nodes []NodeInfo
	for i := range nodeList.Items {
		n := &nodeList.Items[i]
		ip := nodeInternalIP(n)
		if ip == "" {
			return [2]NodeInfo{}, fmt.Errorf("node %s has no InternalIP", n.Name)
		}
		nodes = append(nodes, NodeInfo{Name: n.Name, IP: ip})
	}

	if len(nodes) != 2 {
		return [2]NodeInfo{}, fmt.Errorf("expected 2 control-plane nodes, found %d", len(nodes))
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Name < nodes[j].Name })
	logrus.Infof("Discovered nodes: %s (%s), %s (%s)", nodes[0].Name, nodes[0].IP, nodes[1].Name, nodes[1].IP)
	return [2]NodeInfo{nodes[0], nodes[1]}, nil
}

func nodeInternalIP(node *corev1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}

func verifyTopology(ctx context.Context, cfg Config) error {
	infra, err := cfg.ConfigClient.ConfigV1().Infrastructures().Get(ctx, "cluster", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("reading Infrastructure CR: %w", err)
	}
	if infra.Status.ControlPlaneTopology != configv1.DualReplicaTopologyMode {
		return fmt.Errorf("fencing validation requires DualReplica topology, found %q", infra.Status.ControlPlaneTopology)
	}
	return nil
}

// --- SSH helpers ---

func sshConnect(cfg Config, ip string) (*gossh.Client, error) {
	addr := net.JoinHostPort(ip, sshPort)
	return ssh.NewClient(cfg.SSHUser, addr, cfg.SSHKeys)
}

func sshRun(client *gossh.Client, cmd string) (string, error) {
	wrapped := fmt.Sprintf("sudo bash -lc %s", shellQuote(cmd))
	stdout, stderr, err := ssh.RunOutput(client, wrapped)
	if err != nil {
		return stdout, fmt.Errorf("command failed: %s\nstderr: %s: %w", cmd, stderr, err)
	}
	return stdout, nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

// --- Pre-flight checks ---

func runPreFlight(client *gossh.Client, nodes []NodeInfo) error {
	if err := checkStonith(client); err != nil {
		return err
	}
	if err := checkPacemakerOnline(client, nodes); err != nil {
		return err
	}
	if err := checkDaemons(client); err != nil {
		return err
	}
	return checkEtcdHealth(client, nodes)
}

func checkStonith(client *gossh.Client) error {
	out, err := sshRun(client, "(pcs stonith config || pcs stonith status || pcs stonith show) 2>&1")
	if err != nil || strings.TrimSpace(out) == "" {
		return fmt.Errorf("no STONITH devices configured")
	}

	prop, err := sshRun(client, "pcs property config stonith-enabled 2>&1 || pcs property list stonith-enabled 2>&1 || pcs property show --all stonith-enabled 2>&1")
	if err != nil {
		return fmt.Errorf("could not read stonith-enabled property: %w", err)
	}
	if !parseStonithEnabled(prop) {
		return fmt.Errorf("stonith-enabled is not set to true")
	}
	logrus.Info("STONITH is configured and enabled")
	return nil
}

func checkPacemakerOnline(client *gossh.Client, nodes []NodeInfo) error {
	out, err := sshRun(client, "pcs status nodes 2>/dev/null || crm_mon -1 2>/dev/null")
	if err != nil {
		return fmt.Errorf("checking pacemaker status: %w", err)
	}

	online := parsePacemakerOnline(out)
	for _, n := range nodes {
		if !nodeInOnlineList(n, online) {
			return fmt.Errorf("node %s is not online in Pacemaker", n.Name)
		}
	}
	logrus.Info("Both nodes are online in Pacemaker")
	return nil
}

func checkDaemons(client *gossh.Client) error {
	out, err := sshRun(client, "pcs status --full 2>/dev/null || pcs status 2>/dev/null")
	if err != nil {
		return fmt.Errorf("checking daemon status: %w", err)
	}

	missing := parseDaemonStatus(out)
	if len(missing) > 0 {
		return fmt.Errorf("daemons not active/running: %s", strings.Join(missing, ", "))
	}
	logrus.Info("All daemons (corosync, pacemaker, pcsd) are active")
	return nil
}

func checkEtcdHealth(client *gossh.Client, nodes []NodeInfo) error {
	endpoints := formatEtcdURL(nodes[0].IP) + "," + formatEtcdURL(nodes[1].IP)

	healthCmd := fmt.Sprintf("podman exec etcd sh -lc 'ETCDCTL_API=3 etcdctl -w json endpoint health --endpoints=%s'", endpoints)
	out, err := sshRun(client, healthCmd)
	if err != nil {
		return fmt.Errorf("etcd endpoint health check failed: %w", err)
	}
	if err := parseEtcdHealth(out); err != nil {
		return err
	}

	memberCmd := "podman exec etcd sh -lc 'ETCDCTL_API=3 etcdctl -w json member list'"
	out, err = sshRun(client, memberCmd)
	if err != nil {
		return fmt.Errorf("etcd member list failed: %w", err)
	}
	if err := parseEtcdMembers(out, nodes[0].IP, nodes[1].IP); err != nil {
		return err
	}

	logrus.Info("etcd has 2 healthy voter members")
	return nil
}

// --- Fencing ---

func fenceNode(client *gossh.Client, pcmkName string) error {
	cmd := fmt.Sprintf("timeout %d pcs stonith fence %s", fenceTimeout, pcmkName)
	_, err := sshRun(client, cmd)
	return err
}

// --- Polling ---

func waitNotReady(ctx context.Context, kube kubernetes.Interface, nodeName string) error {
	pollCtx, cancel := context.WithTimeout(ctx, notReadyTimeout)
	defer cancel()
	start := time.Now()
	return wait.PollUntilContextCancel(pollCtx, pollInterval, true, func(ctx context.Context) (bool, error) {
		if !isNodeReady(ctx, kube, nodeName) {
			return true, nil
		}
		if time.Since(start) > time.Minute {
			logrus.Debugf("Still waiting for %s to become NotReady (%s elapsed)", nodeName, time.Since(start).Round(time.Second))
		}
		return false, nil
	})
}

func waitReady(ctx context.Context, kube kubernetes.Interface, nodeName string) error {
	pollCtx, cancel := context.WithTimeout(ctx, readyTimeout)
	defer cancel()
	start := time.Now()
	return wait.PollUntilContextCancel(pollCtx, pollInterval, true, func(ctx context.Context) (bool, error) {
		if isNodeReady(ctx, kube, nodeName) {
			return true, nil
		}
		if time.Since(start) > time.Minute {
			logrus.Debugf("Still waiting for %s to become Ready (%s elapsed)", nodeName, time.Since(start).Round(time.Second))
		}
		return false, nil
	})
}

func isNodeReady(ctx context.Context, kube kubernetes.Interface, nodeName string) bool {
	node, err := kube.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	for _, c := range node.Status.Conditions {
		if c.Type == corev1.NodeReady {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

func pollPacemakerOnline(ctx context.Context, client *gossh.Client, nodes []NodeInfo) error {
	pollCtx, cancel := context.WithTimeout(ctx, pcmkTimeout)
	defer cancel()
	return wait.PollUntilContextCancel(pollCtx, pollInterval, true, func(ctx context.Context) (bool, error) {
		out, err := sshRun(client, "pcs status nodes 2>/dev/null || crm_mon -1 2>/dev/null")
		if err != nil {
			return false, nil
		}
		online := parsePacemakerOnline(out)
		for _, n := range nodes {
			if !nodeInOnlineList(n, online) {
				return false, nil
			}
		}
		return true, nil
	})
}

func pollEtcdHealth(ctx context.Context, client *gossh.Client, nodes []NodeInfo) error {
	endpoints := formatEtcdURL(nodes[0].IP) + "," + formatEtcdURL(nodes[1].IP)
	pollCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()
	return wait.PollUntilContextCancel(pollCtx, pollInterval, true, func(ctx context.Context) (bool, error) {
		healthCmd := fmt.Sprintf("podman exec etcd sh -lc 'ETCDCTL_API=3 etcdctl -w json endpoint health --endpoints=%s'", endpoints)
		out, err := sshRun(client, healthCmd)
		if err != nil {
			return false, nil
		}
		if parseEtcdHealth(out) != nil {
			return false, nil
		}

		memberCmd := "podman exec etcd sh -lc 'ETCDCTL_API=3 etcdctl -w json member list'"
		out, err = sshRun(client, memberCmd)
		if err != nil {
			return false, nil
		}
		if parseEtcdMembers(out, nodes[0].IP, nodes[1].IP) != nil {
			return false, nil
		}
		return true, nil
	})
}

// --- Parsing ---

func parseStonithEnabled(output string) bool {
	return stonithEnabledRe.MatchString(output)
}

func parsePacemakerOnline(output string) []string {
	var names []string
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "Online:") {
			continue
		}
		rest := strings.TrimPrefix(trimmed, "Online:")
		rest = strings.NewReplacer("[", "", "]", "").Replace(rest)
		for _, name := range strings.Fields(rest) {
			if name != "" {
				names = append(names, name)
			}
		}
		break
	}
	return names
}

func nodeInOnlineList(node NodeInfo, online []string) bool {
	short := strings.SplitN(node.Name, ".", 2)[0]
	pcmkShort := strings.SplitN(node.PacemakerName, ".", 2)[0]
	for _, name := range online {
		nameShort := strings.SplitN(name, ".", 2)[0]
		if name == node.Name || name == node.PacemakerName ||
			nameShort == short || nameShort == pcmkShort {
			return true
		}
	}
	return false
}

func parseDaemonStatus(output string) []string {
	inSection := false
	var missing []string
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "Daemon Status:") {
			inSection = true
			continue
		}
		if !inSection {
			continue
		}
		lower := strings.ToLower(strings.TrimSpace(line))
		for _, svc := range []string{"corosync", "pacemaker", "pcsd"} {
			if strings.HasPrefix(lower, svc+":") && !daemonActiveRe.MatchString(line) {
				missing = append(missing, svc)
			}
		}
	}
	return missing
}

type etcdHealthEntry struct {
	Health bool   `json:"health"`
	Error  string `json:"error,omitempty"`
}

func parseEtcdHealth(output string) error {
	var entries []etcdHealthEntry
	if err := json.Unmarshal([]byte(output), &entries); err != nil {
		return fmt.Errorf("failed to parse etcd health output: %w", err)
	}
	for _, e := range entries {
		if !e.Health {
			return fmt.Errorf("unhealthy etcd endpoint: %s", e.Error)
		}
	}
	return nil
}

type etcdMemberList struct {
	Members []etcdMember `json:"members"`
}

type etcdMember struct {
	IsLearner  bool     `json:"isLearner"`
	ClientURLs []string `json:"clientURLs"`
}

func parseEtcdMembers(output, ipA, ipB string) error {
	var list etcdMemberList
	if err := json.Unmarshal([]byte(output), &list); err != nil {
		return fmt.Errorf("failed to parse etcd member list: %w", err)
	}

	foundA, foundB := false, false
	for _, m := range list.Members {
		if m.IsLearner {
			continue
		}
		for _, u := range m.ClientURLs {
			if strings.Contains(u, ipA) {
				foundA = true
			}
			if strings.Contains(u, ipB) {
				foundB = true
			}
		}
	}
	if !foundA || !foundB {
		return fmt.Errorf("etcd does not have 2 voting members (found A=%v, B=%v)", foundA, foundB)
	}
	return nil
}

func resolvePacemakerNames(client *gossh.Client, nodes []NodeInfo) error {
	out, err := sshRun(client, "pcs status nodes 2>/dev/null || crm_mon -1 2>/dev/null")
	if err != nil {
		return fmt.Errorf("resolving pacemaker node names: %w", err)
	}

	online := parsePacemakerOnline(out)
	if len(online) < 2 {
		return fmt.Errorf("expected at least 2 pacemaker nodes online, found %d", len(online))
	}

	for i := range nodes {
		short := strings.SplitN(nodes[i].Name, ".", 2)[0]
		for _, pname := range online {
			pshort := strings.SplitN(pname, ".", 2)[0]
			if pname == nodes[i].Name || pshort == short {
				nodes[i].PacemakerName = pname
				break
			}
		}
		if nodes[i].PacemakerName == "" {
			return fmt.Errorf("could not resolve pacemaker name for node %s", nodes[i].Name)
		}
	}

	logrus.Infof("Pacemaker names: %s → %s, %s → %s",
		nodes[0].Name, nodes[0].PacemakerName, nodes[1].Name, nodes[1].PacemakerName)
	return nil
}

func formatEtcdURL(ip string) string {
	if strings.Contains(ip, ":") {
		return fmt.Sprintf("https://[%s]:2379", ip)
	}
	return fmt.Sprintf("https://%s:2379", ip)
}
