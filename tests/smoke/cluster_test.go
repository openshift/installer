package smoke

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"path/filepath"

	"github.com/coreos/ktestutil/testworkload"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

const (
	// networkingEnv is the environment variable that specifies if calico is running.
	networkingEnv = "SMOKE_NETWORKING"
	// nodeCountEnv is the environment variable that specifies the node count.
	nodeCountEnv = "SMOKE_NODE_COUNT"
	// manifestPathsEnv is the environment variable that defines the paths to the manifests that are deployed on the cluster.
	manifestPathsEnv = "SMOKE_MANIFEST_PATHS"
	// manifestExperimentalEnv is the environment variable that specifies whether or not to test for experimental manifests.
	manifestExperimentalEnv = "SMOKE_MANIFEST_EXPERIMENTAL"
)

var (
	// defaultIgnoredManifests represents the manifests that are ignored by
	// testAllResourcesCreated by default.
	defaultIgnoredManifests = []string{"bootstrap"}

	// experimentalManifests represents the manifests that are ignored by
	// testAllResourcesCreated when manifestExperimentalEnv isn't set to 'true'.
	experimentalManifests = []string{
		// Generated all the time but only deployed when experimental is enabled.
		"tectonic/updater/cluster-config.yaml",
		"tectonic/updater/app_versions/app-version-tectonic-etcd.yaml",
		"tectonic/updater/operators/tectonic-etcd-operator.yaml",
	}

	// equivalentKindRemapping is used by resourceIdentifier to map different
	// Kubernetes object kinds, that can be considered equivalent when checking
	// resource existence, with the same identifier.
	equivalentKindRemapping = map[string]string{
		"extensions/v1beta1:DaemonSet":  "extensions/v1beta1:DeploymentOrDaemonSet",
		"extensions/v1beta1:Deployment": "extensions/v1beta1:DeploymentOrDaemonSet",
	}

	// decodeErrorRegexp defines the format of the error returned by Kubernetes' resource mapper.
	decodeErrorRegexp = regexp.MustCompile(`unable to (?P<Type>decode|recognize) "(?P<Manifest>.*)": (?P<Message>.*)`)
)

func testCluster(t *testing.T) {
	// wait for all nodes to become available
	t.Run("AllNodesRunning", testAllNodesRunning)
	t.Run("GetIdentityLogs", testGetIdentityLogs)
	t.Run("AllPodsRunning", testAllPodsRunning)
	t.Run("KillAPIServer", testKillAPIServer)
	t.Run("AllResourcesCreated", testAllResourcesCreated)

	ne := os.Getenv(networkingEnv)
	if ne == "canal" || ne == "calico" {
		t.Run("NetworkPolicy", testNetworkPolicy)
	}
}

func testAllPodsRunning(t *testing.T) {
	err := retry(allPodsRunning, t, 3*time.Second, 10*time.Minute)
	if err != nil {
		t.Fatalf("Timed out waiting for pods to be ready.")
	}
	t.Log("All pods are ready.")
}

func allPodsRunning(t *testing.T) error {
	c, _ := newClient(t)
	pods, err := c.Core().Pods("").List(meta_v1.ListOptions{})
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
		c, _ := newClient(t)
		nodes, err := c.Core().Nodes().List(meta_v1.ListOptions{})
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
	nodeCount, err := strconv.Atoi(os.Getenv(nodeCountEnv))
	if err != nil {
		t.Fatalf("failed to get number of expected nodes from environment variable %s: %v", nodeCountEnv, err)
	}

	max := 10 * time.Minute
	err = retry(allNodesRunning(nodeCount), t, 10*time.Second, max)
	if err != nil {
		t.Fatalf("Failed to find %d ready nodes in %v.", nodeCount, max)
	}
	t.Logf("Successfully found %d ready nodes.", nodeCount)
}

func getIdentityLogs(t *testing.T) error {
	c, _ := newClient(t)
	podPrefix := "tectonic-identity"
	_, err := validatePodLogging(c, tectonicSystemNamespace, podPrefix)
	if err != nil {
		return fmt.Errorf("failed to gather logs for %s/%s, %v", tectonicSystemNamespace, podPrefix, err)
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
	pods, err := c.Pods(namespace).List(meta_v1.ListOptions{})
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
	c, _ := newClient(t)
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

// allResourcesCreated reads all the manifests recursively from the given paths (except the ones that are explictely
// ignored), and ensures that all the associated resources have been properly created.
func allResourcesCreated(manifestsPaths, ignoredManifests []string) func(t *testing.T) error {
	return func(t *testing.T) error {
		t.Logf("looking for resources defined by the provided manifests...")

		_, cc := newClient(t)
		failed := false

		// Walk recursively through the provided folders to find manifests, decode them and
		// ensure they exist on the cluster.
		resourcesToManifests := make(map[string][]string)
		resourcesCreated := make(map[string]bool)
		errs := walkPathForObjects(cc, manifestsPaths, func(info *resource.Info, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			resourceIdentifier := resourceIdentifier(info)
			manifest := stripPathPrefixes(info.Source, manifestsPaths)

			if containsAnyOfStrings(ignoredManifests, manifest) {
				// The manifest is ignored.
				return nil
			}

			resourcesToManifests[resourceIdentifier] = append(resourcesToManifests[resourceIdentifier], manifest)
			resourcesCreated[resourceIdentifier] = resourcesCreated[resourceIdentifier] || (info.Get() == nil)
			if resourcesCreated[resourceIdentifier] {
				t.Logf("OK : %s - %v", resourceIdentifier, manifest)
			}

			return nil
		})

		// Each resource that is not present on the server, but declared by one or
		// multiple manifests is treated as a failure.
		//
		// TODO: This might be ineffective if resources are clearly different types
		// and purpose (i.e. DaemonSet / ConfigMap) are named the same. To work around
		// this, we could have a mapping of types we consider to be similar. For now,
		// this does not seem to be an issue and the effort is not justified.
		for resourceIdentifier, resourceCreated := range resourcesCreated {
			if !resourceCreated {
				t.Logf("MISSING : %s - %v", resourceIdentifier, resourcesToManifests[resourceIdentifier])
				failed = true
			}
		}

		// Each manifests that failed to get decoded (e.g. unknown type) is treated as
		// a failure. This typically means that the TPR kind, or the operator that is
		// responsible for creating the TPR kind, does not exist.
		for _, err := range errs {
			if containsAnyOfStrings(ignoredManifests, err.Error()) {
				// The manifest is ignored.
				continue
			}
			t.Log(err)
			failed = true
		}

		if failed {
			return errors.New("all defined resources were not present")
		}

		t.Logf("all resources defined by the provided manifests are present")
		return nil
	}
}

func testAllResourcesCreated(t *testing.T) {
	// Read configuration from environment.
	manifestPaths := os.Getenv(manifestPathsEnv)
	manifestsPathsSp := strings.Split(manifestPaths, ",")
	if len(manifestsPathsSp) == 0 {
		t.Skipf("no manifest paths in environment variable %s, skipping", manifestPathsEnv)
	}

	ignoredManifests := defaultIgnoredManifests
	if manifestExperimental := os.Getenv(manifestExperimentalEnv); manifestExperimental != "true" {
		ignoredManifests = append(ignoredManifests, experimentalManifests...)
	}

	max := 10 * time.Minute
	err := retry(allResourcesCreated(manifestsPathsSp, ignoredManifests), t, 30*time.Second, max)
	if err != nil {
		t.Fatalf("timed out waiting for all manifests to be created after %v", max)
	}
}

// testNetworkPolicy permforms 3 tests:
// * first ping test to check if network is setup correctly and reachable.
// * second ping test after setting `default-deny` policy on `network-policy-test` namespace
//   to ensure nothing can talk to each other.
// * third ping test after setting `access-nginx` policy to ensure now nginx workload is reachable.
func testNetworkPolicy(t *testing.T) {
	var (
		namespace = "network-policy-test"
		nginx     *testworkload.Nginx
		client    kubernetes.Interface
	)
	client, _ = newClient(t)
	_, err := client.CoreV1().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
	if apierrors.IsAlreadyExists(err) {
		t.Logf("ns already exists")
	} else if err != nil {
		t.Fatalf("failed to create namespace with name %v", namespace)
	}

	if err := wait.Poll(10*time.Second, 2*time.Minute, func() (bool, error) {
		var err error
		if nginx, err = testworkload.NewNginx(client, namespace, testworkload.WithNginxPingJobLabels(map[string]string{"allow": "access"})); err != nil {
			t.Logf("failed to create test nginx: %v", err)
			return false, nil
		}
		return true, nil
	}); err != nil {
		t.Fatalf("failed to create an testworkload: %v", err)
	}
	defer nginx.Delete()

	if err := wait.Poll(10*time.Second, 2*time.Minute, func() (bool, error) {
		if err := nginx.IsReachable(); err != nil {
			t.Logf("error not reachable %s: %v", nginx.Name, err)
			return false, nil
		}
		return true, nil
	}); err != nil {
		t.Fatalf("network not set up correctly: %v", err)
	}

	t.Run("DefaultDeny", func(t *testing.T) { testDefaultDenyNetworkPolicy(t, client, namespace, nginx) })
	t.Run("NetworkPolicy", func(t *testing.T) { testAllowNetworkPolicy(t, client, namespace, nginx) })
}

func testDefaultDenyNetworkPolicy(t *testing.T, client kubernetes.Interface, namespace string, nginx *testworkload.Nginx) {
	var defaultDenyNetworkPolicy = []byte(`kind: NetworkPolicy
apiVersion: extensions/v1beta1
metadata:
  name: default-deny
spec:
  podSelector:
`)

	npi, _, err := api.Codecs.UniversalDecoder().Decode(defaultDenyNetworkPolicy, nil, &v1beta1.NetworkPolicy{})
	if err != nil {
		t.Fatalf("unable to decode network policy manifest: %v", err)
	}
	np, ok := npi.(*v1beta1.NetworkPolicy)
	if !ok {
		t.Fatalf("expected manifest to decode into *api.networkpolicy, got %T", npi)
	}

	httpRestClient := client.ExtensionsV1beta1().RESTClient()
	uri := fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s",
		strings.ToLower("extensions"),
		strings.ToLower("v1beta1"),
		strings.ToLower(namespace),
		strings.ToLower("NetworkPolicies"))

	result := httpRestClient.Post().RequestURI(uri).Body(np).Do()
	if result.Error() != nil {
		t.Fatal(result.Error())
	}
	defer func() {
		uri = fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s/%s",
			strings.ToLower("extensions"),
			strings.ToLower("v1beta1"),
			strings.ToLower(namespace),
			strings.ToLower("NetworkPolicies"),
			strings.ToLower(np.ObjectMeta.Name))

		result = httpRestClient.Delete().RequestURI(uri).Do()
		if result.Error() != nil {
			t.Fatal(result.Error())
		}

	}()

	if err := wait.Poll(10*time.Second, 2*time.Minute, func() (bool, error) {
		if err := nginx.IsUnReachable(); err != nil {
			t.Logf("error still reachable %s: %v", nginx.Name, err)
			return false, nil
		}
		return true, nil
	}); err != nil {
		t.Fatalf("default deny failed: %v", err)
	}
}

func testAllowNetworkPolicy(t *testing.T, client kubernetes.Interface, namespace string, nginx *testworkload.Nginx) {
	var netPolicyTpl = []byte(`kind: NetworkPolicy
apiVersion: extensions/v1beta1
metadata:
  name: access-nginx
spec:
  podSelector:
    matchLabels:
      app: %s
  ingress:
    - from:
      - podSelector:
          matchLabels:
            allow: access
`)

	netPolicy := fmt.Sprintf(string(netPolicyTpl), nginx.Name)
	npi, _, err := api.Codecs.UniversalDecoder().Decode([]byte(netPolicy), nil, &v1beta1.NetworkPolicy{})
	if err != nil {
		t.Fatalf("unable to decode network policy manifest: %v", err)
	}
	np, ok := npi.(*v1beta1.NetworkPolicy)
	if !ok {
		t.Fatalf("expected manifest to decode into *api.networkpolicy, got %T", npi)
	}

	httpRestClient := client.ExtensionsV1beta1().RESTClient()
	uri := fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s",
		strings.ToLower("extensions"),
		strings.ToLower("v1beta1"),
		strings.ToLower(namespace),
		strings.ToLower("NetworkPolicies"))

	result := httpRestClient.Post().RequestURI(uri).Body(np).Do()
	if result.Error() != nil {
		t.Fatal(result.Error())
	}
	defer func() {
		uri = fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s/%s",
			strings.ToLower("extensions"),
			strings.ToLower("v1beta1"),
			strings.ToLower(namespace),
			strings.ToLower("NetworkPolicies"),
			strings.ToLower(np.ObjectMeta.Name))

		result = httpRestClient.Delete().RequestURI(uri).Do()
		if result.Error() != nil {
			t.Fatal(result.Error())
		}

	}()

	if err := wait.Poll(10*time.Second, 2*time.Minute, func() (bool, error) {
		if err := nginx.IsReachable(); err != nil {
			t.Logf("error not reachable %s: %v", nginx.Name, err)
			return false, nil
		}
		return true, nil
	}); err != nil {
		t.Fatalf("allow nginx network policy failed: %v", err)
	}
}

func getAPIServers(client *kubernetes.Clientset) (*v1.PodList, error) {
	pods, err := client.Core().Pods(kubeSystemNamespace).List(meta_v1.ListOptions{LabelSelector: apiServerSelector})
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

// walkPathForObjects is a helper that calls the given resource.VisitorFunc function for each decoded Kubernetes
// manifest present in the given paths. Any decoding or parsing errors are aggregated.
func walkPathForObjects(cfg clientcmd.ClientConfig, paths []string, fn resource.VisitorFunc) (errs []error) {
	f := cmdutil.NewFactory(cfg)

	schema, err := f.Validator(false, "")
	if err != nil {
		return []error{err}
	}

	mapper, typer, err := f.UnstructuredObject()
	if err != nil {
		return []error{err}
	}

	// As of 1.7, we need to replace the typer with: f.CategoryExpander()
	result := resource.NewBuilder(mapper, typer, resource.ClientMapperFunc(f.UnstructuredClientForMapping), unstructured.UnstructuredJSONScheme).
		ContinueOnError().
		Schema(schema).
		FilenameParam(false, &resource.FilenameOptions{Recursive: true, Filenames: paths}).
		Flatten().
		Do()

	err = result.Err()
	if err != nil && !strings.HasPrefix(err.Error(), "you must provide one or more resources") {
		return []error{err}
	}

	if err := result.Visit(fn); err != nil {
		for _, err := range err.(utilerrors.Aggregate).Errors() {
			if manifest, message, ok := parseMapperDecodingError(err.Error()); ok {
				errs = append(errs, fmt.Errorf("manifest %q not recognized: %s (this is likely due to a missing TPR kind / Operator)", stripPathPrefixes(manifest, paths), message))
			} else {
				errs = append(errs, fmt.Errorf("failed to parse manifest: %s (syntax?)", err))
			}
		}
	}
	return errs
}

// parseMapperDecodingError extracts information from a Kubernetes' mapper
// error.
func parseMapperDecodingError(err string) (manifest, message string, ok bool) {
	tokens := decodeErrorRegexp.FindStringSubmatch(err)
	if tokens == nil {
		return "", "", false
	}

	for i, name := range decodeErrorRegexp.SubexpNames() {
		if name == "Manifest" {
			manifest = tokens[i]
		} else if name == "Message" {
			message = tokens[i]
		}
	}

	return manifest, message, true
}

// stripPathPrefixes attempts to remove a prefix from the given path
// using the provided lists of potential prefixes. If none of the provided
// prefixes matched, the original path is returned.
func stripPathPrefixes(path string, prefixes []string) string {
	for _, prefix := range prefixes {
		if rel, err := filepath.Rel(prefix, path); err == nil {
			return rel
		}
	}
	return path
}

// containsAnyOfStrings returns whether one of the needles is
// contained within the haystack.
func containsAnyOfStrings(needles []string, haystack string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// resourcesIdentifier returns a string that can be used to identify and map
// a Kubernetes resource easily. Some object kinds are treated equivalently
// (see equivalentKindRemapping) in order to ease executing presence tests.
func resourceIdentifier(info *resource.Info) string {
	kindObject := info.VersionedObject.GetObjectKind().GroupVersionKind()
	kind := fmt.Sprintf("%s/%s:%s", kindObject.Group, kindObject.Version, kindObject.Kind)
	if equivalentKind, ok := equivalentKindRemapping[kind]; ok {
		kind = equivalentKind
	}
	return fmt.Sprintf("[Kind: %s | Namespace: %s | Name: %s]", kind, info.Namespace, info.Name)
}
