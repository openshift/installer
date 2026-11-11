package nodejoiner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/configimage"
	"github.com/openshift/installer/pkg/asset/agent/image"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	workflowreport "github.com/openshift/installer/pkg/asset/agent/workflow/report"
	"github.com/openshift/installer/pkg/asset/store"
)

const (
	addNodesResultFile = "exit_code"
)

// systemCACertBundle is the path to the system CA bundle on RHEL/CoreOS nodes.
// It is a variable to allow overriding in tests.
var systemCACertBundle = "/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem"

// NewAddNodesCommand creates a new command for add nodes.
func NewAddNodesCommand(directory string, kubeConfig string, generatePXE bool, generateConfigISO bool) error {
	if generatePXE && generateConfigISO {
		return fmt.Errorf("invalid configuration found")
	}

	err := saveParams(directory, kubeConfig)
	if err != nil {
		return err
	}

	if err := setupProxyCACert(directory, kubeConfig); err != nil {
		logrus.Warnf("Failed to setup proxy CA certificate: %v", err)
	}

	assets := []asset.WritableAsset{
		&workflow.AgentWorkflowAddNodes{},
	}
	var targetAsset asset.WritableAsset
	switch {
	case generatePXE:
		targetAsset = &image.AgentPXEFiles{}
	case generateConfigISO:
		targetAsset = &configimage.ConfigImage{}
	default:
		targetAsset = &image.AgentImage{}
	}
	assets = append(assets, targetAsset)

	ctx := workflowreport.Context(string(workflow.AgentWorkflowTypeAddNodes), directory)

	fetcher := store.NewAssetsFetcher(directory)
	err = fetcher.FetchAndPersist(ctx, assets)

	if reportErr := workflowreport.GetReport(ctx).Complete(err); reportErr != nil {
		return reportErr
	}

	// Save the exit code result
	exitCode := "0"
	if err != nil {
		exitCode = "1"
	}
	if err2 := os.WriteFile(filepath.Join(directory, addNodesResultFile), []byte(exitCode), 0600); err2 != nil {
		return err2
	}

	return err
}

// setupProxyCACert fetches the proxy trusted CA from the cluster and configures
// SSL_CERT_FILE so that subsequent image pulls through the proxy succeed.
func setupProxyCACert(directory, kubeConfig string) error {
	var config *rest.Config
	var err error
	if kubeConfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return fmt.Errorf("cannot build cluster config: %w", err)
	}

	configClient, err := configclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("cannot create config client: %w", err)
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("cannot create kubernetes client: %w", err)
	}

	return setupProxyCACertWithClients(directory, configClient, k8sClient)
}

// setupProxyCACertWithClients contains the testable logic for setupProxyCACert.
// The combined bundle includes the system CA certs so public registries remain trusted.
func setupProxyCACertWithClients(directory string, configClient configclient.Interface, k8sClient kubernetes.Interface) error {
	proxy, err := configClient.ConfigV1().Proxies().Get(context.Background(), "cluster", metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("cannot get proxy config: %w", err)
	}

	if proxy.Spec.TrustedCA.Name == "" {
		return nil
	}

	caConfigMap, err := k8sClient.CoreV1().ConfigMaps("openshift-config").Get(context.Background(), proxy.Spec.TrustedCA.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("cannot get CA bundle configmap %q: %w", proxy.Spec.TrustedCA.Name, err)
	}

	proxyCACert, ok := caConfigMap.Data["ca-bundle.crt"]
	if !ok || proxyCACert == "" {
		return nil
	}

	// Read the system CA bundle so public registries remain trusted alongside the proxy CA.
	systemCerts, err := os.ReadFile(systemCACertBundle)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot read system CA bundle: %w", err)
	}

	combined := systemCerts
	if len(combined) > 0 && combined[len(combined)-1] != '\n' {
		combined = append(combined, '\n')
	}
	combined = append(combined, []byte(proxyCACert)...)
	caFile := filepath.Join(directory, "proxy-ca-bundle.crt")
	if err := os.WriteFile(caFile, combined, 0600); err != nil {
		return fmt.Errorf("cannot write combined CA bundle: %w", err)
	}

	logrus.Infof("Proxy CA certificate configured from cluster configmap %q", proxy.Spec.TrustedCA.Name)
	return os.Setenv("SSL_CERT_FILE", caFile)
}

func saveParams(directory, kubeConfig string) error {
	// Store the current parameters into the assets folder, so
	// that they could be retrieved later by the assets
	params := joiner.Params{
		Kubeconfig: kubeConfig,
	}
	return params.Save(directory)
}
