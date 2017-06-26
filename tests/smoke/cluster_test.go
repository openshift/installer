package smoke

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	// nodeCountVar is the environment variable to check for Node count.
	nodeCountVar = "NODE_COUNT"
)

func testCluster(t *testing.T) {
	// wait for all nodes to become available
	t.Run("AllNodesRunning", testAllNodesRunning)
	t.Run("GetIdentityLogs", testGetIdentityLogs)
	t.Run("AllPodsRunning", testAllPodsRunning)
	t.Run("KillAPIServer", testKillAPIServer)
}

func testAllPodsRunning(t *testing.T) {
	err := retry(allPodsRunning, t, 3*time.Second, 10*time.Minute)
	if err != nil {
		t.Fatalf("Timed out waiting for pods to be ready.")
	}
	t.Log("All pods are ready.")
}

func allPodsRunning(t *testing.T) error {
	c := newClient(t)
	pods, err := c.Core().Pods("").List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("could not list pods: %v", err)
	}

	allReady := len(pods.Items) != 0
	for _, p := range pods.Items {
		if p.Status.Phase != v1.PodRunning {
			allReady = false
			t.Logf("pod %s/%s not running", p.Namespace, p.Name)
		}
	}
	if !allReady {
		return errors.New("pods are not all ready")
	}
	return nil
}

func allNodesRunning(expected int) func(t *testing.T) error {
	return func(t *testing.T) error {
		c := newClient(t)
		nodes, err := c.Core().Nodes().List(v1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to list nodes: %v", err)
		}
		allReady := len(nodes.Items) != 0
		for _, node := range nodes.Items {
			if nodeReady(node) {
				t.Logf("node %s ready", node.Name)
				continue
			}
			allReady = false
			t.Logf("node %s not ready", node.Name)
		}
		if got := len(nodes.Items); got != expected {
			return fmt.Errorf("expected %d nodes, got %d", expected, got)
		}
		if !allReady {
			return errors.New("nodes are not all ready")
		}
		return nil
	}
}

func testAllNodesRunning(t *testing.T) {
	nodeCount, err := strconv.Atoi(os.Getenv(nodeCountVar))
	if err != nil {
		t.Fatalf("failed to get number of expected nodes from envvar %s: %v", nodeCountVar, err)
	}

	max := 10 * time.Minute
	err = retry(allNodesRunning(nodeCount), t, 10*time.Second, max)
	if err != nil {
		t.Fatalf("Failed to find %d ready nodes in %v.", nodeCount, max)
	}
	t.Logf("Successfully found %d ready nodes.", nodeCount)
}

func getIdentityLogs(t *testing.T) error {
	c := newClient(t)
	namespace := "tectonic-system"
	podPrefix := "tectonic-identity"
	_, err := validatePodLogging(c, namespace, podPrefix)
	if err != nil {
		return fmt.Errorf("failed to gather logs for %s/%s, %v", namespace, podPrefix, err)
	}
	return nil
}

func testGetIdentityLogs(t *testing.T) {
	max := 10 * time.Minute
	err := retry(getIdentityLogs, t, 15*time.Second, max)
	if err != nil {
		t.Fatalf("Failed to gather identity logs in %v.", max)
	}
	t.Log("Successfully gathered identity logs.")
}

// validatePodLogging verifies that logs can be retrieved for a container in Pod.
func validatePodLogging(c *kubernetes.Clientset, namespace, podPrefix string) ([]byte, error) {
	var logs []byte
	pods, err := c.Pods(namespace).List(v1.ListOptions{})
	if err != nil {
		return logs, fmt.Errorf("could not list pods: %v", err)
	}

	var names string
	for _, p := range pods.Items {
		if len(names) != 0 {
			names += ", "
		}
		names += p.Name

		if !strings.HasPrefix(p.Name, podPrefix) {
			continue
		}
		if len(p.Spec.Containers) == 0 {
			return logs, fmt.Errorf("%s pod has no containers", p.Name)
		}

		opt := v1.PodLogOptions{
			Container: p.Spec.Containers[0].Name,
		}

		result := c.Core().Pods(namespace).GetLogs(p.Name, &opt).Do()
		if err := result.Error(); err != nil {
			return logs, fmt.Errorf("failed to get pod logs: %v", err)
		}

		var statusCode int
		result.StatusCode(&statusCode)
		if statusCode/100 != 2 {
			return logs, fmt.Errorf("expected 200 from log response, got %d", statusCode)
		}
		logs, err := result.Raw()
		if err != nil {
			return logs, fmt.Errorf("failed to read logs: %v", err)
		}
		return logs, nil
	}

	return logs, fmt.Errorf("failed to find pods with prefix %q (found pods in %s: %s)", podPrefix, namespace, names)
}

func testKillAPIServer(t *testing.T) {
	c := newClient(t)
	pods, err := getAPIServers(c)
	if err != nil {
		t.Fatalf("Failed to get API server pod: %v", err)
	}

	oldPod := map[string]bool{}

	// Nuke all API servers.
	for _, pod := range pods.Items {
		if err := c.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
			t.Fatalf("Failed to delete API server pod %s: %v", pod.Name, err)
		}
		oldPod[pod.Name] = true
	}

	// API servers and temp API servers come in and out. Ensure
	// that the API server we detect is running for a couple
	// iterations.
	runningLastTime := false

	apiServerUp := func(t *testing.T) error {
		pods, err := getAPIServers(c)
		if err != nil {
			return fmt.Errorf("failed to get API server pod: %v", err)
		}

		for _, pod := range pods.Items {
			if oldPod[pod.Name] {
				return fmt.Errorf("old API server %s still running", pod.Name)
			}
		}

		allReady := len(pods.Items) != 0
		for _, p := range pods.Items {
			if p.Status.Phase != v1.PodRunning {
				allReady = false
			}
		}

		if allReady {
			if runningLastTime {
				return nil
			}
			runningLastTime = true
		}
		return fmt.Errorf("API server has not yet been running for more than one check")
	}

	max := 6 * time.Minute
	err = retry(apiServerUp, t, 10*time.Second, max)
	if err != nil {
		t.Fatalf("Failed waiting for API server pods to be ready in %v.", max)
	}
	t.Log("API server pods successfully came back up.")
}

func getAPIServers(client *kubernetes.Clientset) (*v1.PodList, error) {
	const (
		apiServerSelector   = "k8s-app=kube-apiserver"
		kubeSystemNamespace = "kube-system"
	)
	pods, err := client.Core().Pods(kubeSystemNamespace).List(v1.ListOptions{LabelSelector: apiServerSelector})
	if err != nil {
		return nil, err
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no pods matched the label selector %q in the %s namespace", apiServerSelector, kubeSystemNamespace)
	}
	return pods, nil
}

// podsStr prints a comma separated list of namespaced Pod names
func podsStr(pods []v1.Pod) (out string) {
	for n, p := range pods {
		// add comma to all entries except first
		if n != 0 {
			out += ", "
		}
		out += fmt.Sprintf("%s/%s", p.GetNamespace(), p.GetName())
	}
	return
}

func nodeReady(node v1.Node) (ok bool) {
	for _, cond := range node.Status.Conditions {
		if cond.Type == v1.NodeReady {
			return cond.Status == v1.ConditionTrue
		}
	}
	return false
}
