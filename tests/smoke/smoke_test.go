package smoke

import (
	"errors"
	"flag"
	"os"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// kubeconfigEnv is the environment variable that is checked for a the kubeconfig path to be loaded.
	kubeconfigEnv = "SMOKE_KUBECONFIG"
	// apiServerSelector is the pod label selector for the apiserver.
	apiServerSelector = "k8s-app=kube-apiserver"
	// kubeSystemNamespace is the namespace for k8s.
	kubeSystemNamespace = "kube-system"
	// tectonicSystemNamespace is the namespace for Tectonic.
	tectonicSystemNamespace = "tectonic-system"
)

var (
	// runClusterTests is used as a flag to control whether or not to run cluster tests.
	runClusterTests bool
)

func TestMain(m *testing.M) {
	flag.BoolVar(&runClusterTests, "cluster", false, "run cluster tests (default false)")
	flag.Parse()
	os.Exit(m.Run())
}

// Test is the only test suite run by default. This function will run all common tests,
// and conditionally run other tests suites.
func Test(t *testing.T) {
	t.Run("Common", testCommon)
	if runClusterTests {
		t.Run("Cluster", testCluster)
	}
}

// newClient is a convenient helper for generating a client-go client from a *testing.T.
func newClient(t *testing.T) (*kubernetes.Clientset, clientcmd.ClientConfig) {
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
	return cs, cfg
}

// stopped returns true if the done chan is closed and returns false otherwise.
func stopped(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

// sleepOrDone blocks until either done is closed or the sleep timer runs out.
func sleepOrDone(sleep time.Duration, done <-chan struct{}) {
	select {
	case <-time.After(sleep):
		return
	case <-done:
		return
	}
}

// retriable describes a test function that can be retried using `retry`.
type retriable func(t *testing.T) error

// timeout is a useful helper that returns a chan that is closed after
// the specified duration. This allows selecting the timeout to return
// immediately.
func timeout(d time.Duration) <-chan struct{} {
	t := time.NewTimer(d)
	c := make(chan struct{})
	go func() {
		<-t.C
		close(c)
	}()
	return c
}

// retry abstracts all retry and timeout logic from the tests.
// This function will retry to run the given retriable every interval
// until either it returns no error or the max time is reached.
func retry(r retriable, t *testing.T, interval, max time.Duration) error {
	done := timeout(max)
	for !stopped(done) {
		err := r(t)
		if err == nil {
			return nil
		}
		t.Logf("failed with error: %v", err)
		t.Logf("retrying in %v", interval)
		sleepOrDone(interval, done)
	}
	return errors.New("timed out while retrying")
}
