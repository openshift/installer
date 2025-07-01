package agent

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	certificatesClient "k8s.io/client-go/kubernetes/typed/certificates/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClusterKubeAPIClient is a kube client to interact with the cluster that agent installer is installing.
type ClusterKubeAPIClient struct {
	Client     *kubernetes.Clientset
	csrClient  certificatesClient.CertificateSigningRequestInterface
	ctx        context.Context
	Config     *rest.Config
	configPath string
}

// NewClusterKubeAPIClient Create a new kube client to interact with the cluster under install.
func NewClusterKubeAPIClient(ctx context.Context, kubeconfigPath string) (*ClusterKubeAPIClient, error) {
	kubeClient := &ClusterKubeAPIClient{}

	var kubeconfig *rest.Config
	var err error
	if kubeconfigPath != "" {
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		kubeconfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, errors.Wrap(err, "error loading kubeconfig from assets")
	}

	kubeclient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating a Kubernetes client from assets failed")
	}

	csrClient := kubeclient.CertificatesV1().CertificateSigningRequests()

	kubeClient.Client = kubeclient
	kubeClient.csrClient = csrClient
	kubeClient.ctx = ctx
	kubeClient.Config = kubeconfig
	kubeClient.configPath = kubeconfigPath

	return kubeClient, nil
}

// IsKubeAPILive Determine if the cluster under install has initailized the kubenertes API.
func (kube *ClusterKubeAPIClient) IsKubeAPILive() bool {
	discovery := kube.Client.Discovery()
	_, err := discovery.ServerVersion()
	return err == nil
}

// DoesKubeConfigExist Determine if the kubeconfig for the cluster can be used without errors.
func (kube *ClusterKubeAPIClient) DoesKubeConfigExist() (bool, error) {
	_, err := clientcmd.LoadFromFile(kube.configPath)
	if err != nil {
		return false, errors.Wrap(err, "error loading kubeconfig from file")
	}
	return true, nil
}

// IsBootstrapConfigMapComplete Detemine if the cluster's bootstrap configmap has the status complete.
func (kube *ClusterKubeAPIClient) IsBootstrapConfigMapComplete() (bool, error) {
	// Get latest version of bootstrap configmap
	bootstrap, err := kube.Client.CoreV1().ConfigMaps("kube-system").Get(kube.ctx, "bootstrap", metav1.GetOptions{})

	if err != nil {
		// bootstrap configmap not found
		return false, nil
	}
	// Found a bootstrap configmap need to check its status
	if bootstrap != nil {
		status, ok := bootstrap.Data["status"]
		if !ok {
			logrus.Debug("no status found in bootstrap configmap")
			return false, nil
		}
		if status == "complete" {
			return true, nil
		}
	}
	return false, nil
}

// ListNodes returns a list of nodes that have joined the cluster.
func (kube *ClusterKubeAPIClient) ListNodes() (*corev1.NodeList, error) {
	nodeList, err := kube.Client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return &corev1.NodeList{}, err
	}
	return nodeList, nil
}

// ListCSRs returns a list of this cluster's CSRs.
func (kube *ClusterKubeAPIClient) ListCSRs() (*certificatesv1.CertificateSigningRequestList, error) {
	csrs, err := kube.csrClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return csrs, nil
}
