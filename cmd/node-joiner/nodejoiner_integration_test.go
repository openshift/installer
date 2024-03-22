package main

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/yaml"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/clusterapi"
)

const (
	testResourcesFolder = "setup"
)

func TestMain(m *testing.M) {

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

	// Prepare the temporary cluster that will be used by node-joiner.
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			// Preload OpenShift specific CRDs so that it will be possible to create specific OpenShift resources.
			filepath.Join(build.Default.GOPATH, "pkg", "mod", "github.com", "openshift", "api@v0.0.0-20240301093301-ce10821dc999", "config", "v1"),
		},
		ErrorIfCRDPathMissing: true,
	}
	t.Cleanup(func() {
		// Ensures that the temporary cluster is stopped when the testscript test is completed.
		assert.NoError(t, testEnv.Stop())
	})

	testscript.Run(t, testscript.Params{
		Dir: "testdata",

		// Uncomment below line to help debug the testcases
		//TestWork: true,

		Deadline: time.Now().Add(10 * time.Minute),

		Setup: func(e *testscript.Env) error {
			// This is required to allow loading the embedded resources.
			e.Cd = filepath.Join(projectDir, "../../data")
			// Set the home dir within the test temporary working directory.
			for i, v := range e.Vars {
				if v == "HOME=/no-home" {
					homeDir := filepath.Join(e.WorkDir, "home")
					if err := os.Mkdir(homeDir, 0777); err != nil {
						return err
					}

					e.Vars[i] = fmt.Sprintf("HOME=%s", homeDir)
					break
				}
			}

			// Reuse clusterapi binaries for testEnv
			if err := os.Chdir(path.Join(projectDir, "../../pkg/clusterapi")); err != nil {
				return err
			}
			defer os.Chdir(projectDir)

			binDir := filepath.Join(e.WorkDir, "bin", "cluster-api")
			if err := clusterapi.UnpackClusterAPIBinary(binDir); err != nil {
				return fmt.Errorf("failed to unpack cluster-api binary: %w", err)
			}
			if err := clusterapi.UnpackEnvtestBinaries(binDir); err != nil {
				return fmt.Errorf("failed to unpack envtest binaries: %w", err)
			}
			testEnv.BinaryAssetsDirectory = binDir

			// Creates a new temporary cluster.
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
			// Setup any required resource specified in the $WORK/setup folder in the cluster.
			return setupInitialResources(testEnv.Config, e.WorkDir)
		},
	})
}

func setupInitialResources(config *rest.Config, workDir string) error {
	csDynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	// Look for any evental resource stored by the current test in the $WORK/setup folder.
	setupPath := filepath.Join(workDir, testResourcesFolder)
	files, err := os.ReadDir(setupPath)
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
		obj := &unstructured.Unstructured{}
		err = yaml.Unmarshal(data, obj)
		if err != nil {
			return fmt.Errorf("%s: %w", fName, err)
		}

		gvr, err := getGVR(obj)
		if err != nil {
			return fmt.Errorf("Error while getting resource gvr from %s: %w", fName, err)
		}
		updObj, err := csDynamic.Resource(gvr).Namespace(obj.GetNamespace()).Create(context.Background(), obj, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Error while creating resource from %s: %w", fName, err)
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
	case "Proxy":
		gvr = v1.SchemeGroupVersion.WithResource("proxies")
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
