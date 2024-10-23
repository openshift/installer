package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"

	v1 "github.com/openshift/api/config/v1"
	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/internal/tshelpers"
)

const (
	testResourcesFolder = "setup"
)

func TestMain(m *testing.M) {
	// Set up the logger for testing
	log.SetLogger(logr.Logger{})

	os.Exit(testscript.RunMain(m, map[string]func() int{
		"node-joiner": func() int {
			if err := nodeJoiner(); err != nil {
				return 1
			}
			return 0
		},
	}))
}

func TestNodeJoinerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	projectDir, err := os.Getwd()
	assert.NoError(t, err)

	testscript.Run(t, testscript.Params{
		Dir: "testdata",

		// Uncomment below line to help debug the testcases
		// TestWork: true,

		Deadline: time.Now().Add(10 * time.Minute),

		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"isoCmp":              tshelpers.IsoCmp,
			"isoCmpRegEx":         tshelpers.IsoCmpRegEx,
			"isoIgnitionContains": tshelpers.IsoIgnitionContains,
			"isoIgnitionUser":     tshelpers.IsoIgnitionUser,
			"isoFileCmpRegEx":     tshelpers.IsoFileCmpRegEx,
		},

		Setup: func(e *testscript.Env) error {
			// This is required for loading properly the embedded resources.
			e.Cd = filepath.Join(projectDir, "../../data")

			// Set the home dir within the test temporary working directory.
			homeDir := filepath.Join(e.WorkDir, "home")
			if err := os.Mkdir(homeDir, 0777); err != nil {
				return err
			}
			for i, v := range e.Vars {
				if v == "HOME=/no-home" {
					e.Vars[i] = fmt.Sprintf("HOME=%s", homeDir)
					break
				}
			}

			// Reuse the current test tmp folder for envTest and fakeRegistry.
			tmpDir := e.Getenv("TMPDIR")

			// Create the fake registry
			fakeRegistry := tshelpers.NewFakeOCPRegistry(tmpDir)

			// Creates a new temporary cluster.
			etcdDataDir, err := os.MkdirTemp(tmpDir, "etcd")
			assert.NoError(t, err)
			apiServerDataDir, err := os.MkdirTemp(tmpDir, "api-server")
			assert.NoError(t, err)

			testEnv := &envtest.Environment{
				CRDDirectoryPaths: []string{
					// Preload OpenShift specific CRDs.
					filepath.Join(projectDir, "testdata", "setup", "crds"),
				},
				ErrorIfCRDPathMissing: true,

				// Uncomment the following line if you wish to run the test without
				// using the hack/go-integration-test-nodejoiner.sh script.
				// BinaryAssetsDirectory: "/tmp/k8s/1.31.0-linux-amd64",

				ControlPlane: envtest.ControlPlane{
					Etcd: &envtest.Etcd{
						DataDir: etcdDataDir,
					},
					APIServer: &envtest.APIServer{
						CertDir: apiServerDataDir,
					},
				},
			}
			// Ensures they are cleaned up on test completion.
			e.Defer(func() {
				assert.NoError(t, testEnv.Stop())
				fakeRegistry.Close()
			})
			// Starts the registry and cluster.
			err = fakeRegistry.Start()
			if err != nil {
				return err
			}
			config, err := testEnv.Start()
			if err != nil {
				return err
			}
			// Creates a valid kubeconfig and store it in the test temporary working dir,
			// so that it could be used by the node-joiner.
			err = createKubeConfig(config, e.WorkDir)
			if err != nil {
				return err
			}

			// TEST_IMAGE env var will be used to replace the OCP release reference in the
			// yaml setup files, so that the one exposed by the fake registry will be used.
			e.Setenv("TEST_IMAGE", fakeRegistry.ReleasePullspec())

			// Setup global resources required for any tests.
			err = setupInitialResources(testEnv.Config, filepath.Join(projectDir, "testdata", "setup", "default"), e.Vars)
			if err != nil {
				return err
			}
			// Setup test specific resources (defined in the $WORK/setup folder).
			return setupInitialResources(testEnv.Config, filepath.Join(e.WorkDir, testResourcesFolder), e.Vars)
		},
	})
}

func setupInitialResources(config *rest.Config, setupPath string, envArgs []string) error {
	files, err := os.ReadDir(setupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	csDynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	// For any valid yaml file, create the related resource.
	for _, f := range files {
		fName := filepath.Join(setupPath, f.Name())
		if filepath.Ext(fName) != ".yaml" && filepath.Ext(fName) != ".yml" {
			continue
		}

		data, err := os.ReadFile(fName)
		if err != nil {
			return err
		}

		// env vars expansion
		for _, ev := range envArgs {
			parts := strings.Split(ev, "=")
			varName := fmt.Sprintf("$%s", parts[0])
			varValue := parts[1]
			data = bytes.ReplaceAll(data, []byte(varName), []byte(varValue))
		}

		obj := &unstructured.Unstructured{}
		err = yaml.Unmarshal(data, obj)
		if err != nil {
			return fmt.Errorf("%s: %w", fName, err)
		}

		gvr, err := getGVR(obj)
		if err != nil {
			return fmt.Errorf("Error while getting resource gvr from %s: %w", fName, err)
		}
		// Create or update the resource (if it already exists).
		updObj, err := csDynamic.Resource(gvr).Namespace(obj.GetNamespace()).Get(context.Background(), obj.GetName(), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				updObj, err = csDynamic.Resource(gvr).Namespace(obj.GetNamespace()).Create(context.Background(), obj, metav1.CreateOptions{})
				if err != nil {
					return fmt.Errorf("Error while creating resource from %s: %w", fName, err)
				}
			}
		} else {
			obj.SetResourceVersion(updObj.GetResourceVersion())
			updObj, err = csDynamic.Resource(gvr).Namespace(obj.GetNamespace()).Update(context.Background(), obj, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("Error while updating resource from %s: %w", fName, err)
			}
		}
		// Take care of a resource status, in case it was configured.
		if status, ok := obj.Object["status"]; ok {
			updObj.Object["status"] = status
			_, err = csDynamic.Resource(gvr).Namespace(obj.GetNamespace()).UpdateStatus(context.Background(), updObj, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("Error while updating resource status from %s: %w", fName, err)
			}
		}
	}

	return nil
}

func getGVR(obj *unstructured.Unstructured) (schema.GroupVersionResource, error) {
	var gvr schema.GroupVersionResource
	var err error

	kind := obj.GetKind()
	switch kind {
	case "ClusterVersion":
		gvr = v1.SchemeGroupVersion.WithResource("clusterversions")
	case "Infrastructure":
		gvr = v1.GroupVersion.WithResource("infrastructures")
	case "Proxy":
		gvr = v1.SchemeGroupVersion.WithResource("proxies")
	case "ImageDigestMirrorSet":
		gvr = v1.SchemeGroupVersion.WithResource("imagedigestmirrorsets")
	case "MachineConfig":
		gvr = machineconfigv1.SchemeGroupVersion.WithResource("machineconfigs")
	case "Namespace":
		gvr = corev1.SchemeGroupVersion.WithResource("namespaces")
	case "Secret":
		gvr = corev1.SchemeGroupVersion.WithResource("secrets")
	case "Node":
		gvr = corev1.SchemeGroupVersion.WithResource("nodes")
	case "ConfigMap":
		gvr = corev1.SchemeGroupVersion.WithResource("configmaps")
	default:
		err = fmt.Errorf("unsupported object kind: %s", kind)
	}

	return gvr, err
}

func createKubeConfig(config *rest.Config, destPath string) error {
	clusterName := "nodejoiner-cluster"
	clusterContext := "nodejoiner-context"
	clusterUser := "nodejoiner-user"

	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters[clusterName] = &clientcmdapi.Cluster{
		Server:                   config.Host,
		CertificateAuthorityData: config.CAData,
	}
	contexts := make(map[string]*clientcmdapi.Context)
	contexts[clusterContext] = &clientcmdapi.Context{
		Cluster:  clusterName,
		AuthInfo: clusterUser,
	}
	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos[clusterUser] = &clientcmdapi.AuthInfo{
		ClientCertificateData: config.CertData,
		ClientKeyData:         config.KeyData,
	}
	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: clusterContext,
		AuthInfos:      authinfos,
	}

	kubeConfigFile := filepath.Join(destPath, "kubeconfig")
	return clientcmd.WriteToFile(clientConfig, kubeConfigFile)
}
