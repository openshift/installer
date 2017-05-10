package sanity

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// nodeCountVar is the environment variable to check for Node count.
	nodeCountVar = "NODE_COUNT"

	// kubeconfigEnv is the environment variable that is checked for a the kubeconfig path to be loaded.
	kubeconfigEnv = "TEST_KUBECONFIG"
)

type timer struct {
	timeout time.Time
}

func newTimer(timeout time.Duration) *timer {
	return &timer{time.Now().Add(timeout)}
}

func (t *timer) timedOut() bool {
	return time.Now().After(t.timeout)
}

func newClient(t *testing.T) *kubernetes.Clientset {
	kcfgPath := os.Getenv(kubeconfigEnv)
	if len(kcfgPath) == 0 {
		t.Fatalf("no kubeconfig path in environment variable %s", kubeconfigEnv)
	}

	rules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kcfgPath}
	cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	restConfig, err := cfg.ClientConfig()
	if err != nil {
		t.Fatalf("could not create client config: %v", err)
	}

	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		t.Fatalf("could not create clientset: %v", err)
	}
	return cs
}

func TestCluster(t *testing.T) {
	// verify that api server is up in 10 min
	t.Run("APIAvailable", testAPIAvailable)

	// wait for all nodes to become available
	t.Run("AllNodesRunning", testAllNodesRunning)
	t.Run("GetLogs", testLogs)
	t.Run("AllPodsRunning", testAllPodsRunning)
	t.Run("KillAPIServer", testKillAPIServer)
}

func testAPIAvailable(t *testing.T) {
	// chan signaled when API server found
	done := waitForAPIServer(t)

	// timeout searching for server
	wait := 10 * time.Minute
	t.Logf("Waiting %v for API server to become available", wait)

	timeout := time.After(wait)
	select {
	case <-timeout:
		t.Fatalf("Could not connect to API server in %v, FAILING!", wait)
	case <-done:
		t.Log("API server is available.")
		return
	}
}

func testAllPodsRunning(t *testing.T) {
	c := newClient(t)

	timer := newTimer(10 * time.Minute)

	for {
		if timer.timedOut() {
			t.Fatalf("timed out waiting for pods to be ready.")
		}

		pods, err := c.Core().Pods("").List(v1.ListOptions{})
		if err != nil {
			t.Logf("could not list pods: %v", err)
			pods = &v1.PodList{}
		}

		allReady := len(pods.Items) != 0
		for _, p := range pods.Items {
			if p.Status.Phase != v1.PodRunning {
				allReady = false
				t.Logf("pod %s/%s not running", p.Namespace, p.Name)
			}
		}
		if allReady {
			return
		}

		time.Sleep(3 * time.Second)
	}
}

func testLogs(t *testing.T) {
	c := newClient(t)

	namespace := "tectonic-system"
	podPrefix := "tectonic-identity"

	wait := 3 * time.Minute
	timeout := time.After(wait)
	done := make(chan struct{})
	go func() {
		for {
			err := validatePodLogging(c, namespace, podPrefix)
			if err == nil {
				done <- struct{}{}
				return
			}
			t.Log("Failed to get Pod logs with error: ", err)
			time.Sleep(3 * time.Second)
		}
	}()

	select {
	case <-timeout:
		t.Fatalf("Failed to gather logs for %s/%s* in %v", namespace, podPrefix, wait)
	case <-done:
		return
	}
}

// validatePodLogging verifies that logs can be retrieved for a container in Pod.
func validatePodLogging(c *kubernetes.Clientset, namespace, podPrefix string) error {
	pods, err := c.Pods(namespace).List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("could not list pods: %v", err)
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
			return fmt.Errorf("tectonic identity pod has no containers")
		}

		opt := v1.PodLogOptions{
			Container: p.Spec.Containers[0].Name,
		}

		result := c.Core().Pods(namespace).GetLogs(p.Name, &opt).Do()
		if err := result.Error(); err != nil {
			return fmt.Errorf("failed to get pod logs: %v", err)
		}

		var statusCode int
		result.StatusCode(&statusCode)
		if statusCode/100 != 2 {
			return fmt.Errorf("expected 200 from log response, got %d", statusCode)
		}
		return nil
	}

	return fmt.Errorf("failed to find tectonic-identity pod (found pods in %s: %s)", namespace, names)
}

func testAllNodesRunning(t *testing.T) {
	c := newClient(t)

	expNodeCount, err := strconv.Atoi(os.Getenv(nodeCountVar))
	if err != nil {
		t.Fatalf("failed to get number of expected nodes from envvar %s: %v", nodeCountVar, err)
	}

	timer := newTimer(10 * time.Minute)
	for {
		if timer.timedOut() {
			t.Fatalf("timed out waiting for nodes to be ready.")
		}

		nodes, err := c.Core().Nodes().List(v1.ListOptions{})
		if err != nil {
			t.Fatalf("could not list nodes: %v", err)
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

		if allReady {
			return
		}

		if got := len(nodes.Items); got != expNodeCount {
			t.Logf("expected %d nodes got %d", expNodeCount, got)
		}

		time.Sleep(10 * time.Second)
	}
}

func testKillAPIServer(t *testing.T) {
	c := newClient(t)
	pods, err := getAPIServers(c)
	if err != nil {
		t.Fatalf("get apiserver pod: %v", err)
	}

	oldPod := map[string]bool{}

	// Nuke all API servers.
	for _, pod := range pods.Items {
		if err := c.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
			t.Fatalf("failed to delete pod %s: %v", pod.Name, err)
		}
		oldPod[pod.Name] = true
	}

	// API servers and temp API servers come in and out. Ensure
	// that the API server we detect is running for a couple
	// iterations.
	runningLastTime := false

	apiServerUp := func() bool {
		pods, err := getAPIServers(c)
		if err != nil {
			t.Logf("get apiserver pod: %v", err)
			return false
		}

		for _, pod := range pods.Items {
			if oldPod[pod.Name] {
				t.Logf("old API server %s still running", pod.Name)
				return false
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
				return true
			}
			runningLastTime = true
		}
		return false
	}

	timer := newTimer(6 * time.Minute)
	for {
		if timer.timedOut() {
			t.Fatalf("timed out waiting for pods to be ready.")
		}

		if apiServerUp() {
			return
		}

		time.Sleep(10 * time.Second)
	}
}

func waitForAPIServer(t *testing.T) <-chan struct{} {
	done := make(chan struct{}, 1)
	go func() {
		var client *kubernetes.Clientset
		for {
			client = newClient(t)
			_, err := client.ServerVersion()
			if err == nil {
				done <- struct{}{}
				return
			}
			wait := 10 * time.Second
			t.Logf("Waiting %v after failed attempt to connect to API server. Error was: %v", wait, err)
			time.Sleep(wait)
		}
	}()
	return done
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
