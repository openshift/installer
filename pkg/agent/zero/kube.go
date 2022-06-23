package zero

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type clusterZeroKubeAPIClient struct {
	Client     *kubernetes.Clientset
	ctx        context.Context
	config     *rest.Config
	configPath string
}

func NewClusterZeroKubeAPIClient(ctx context.Context, assetDir string) (*clusterZeroKubeAPIClient, error) {

	zeroKubeClient := &clusterZeroKubeAPIClient{}

	kubeconfigpath := filepath.Join(assetDir, "auth", "kubeconfig")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigpath)
	if err != nil {
		return nil, errors.Wrap(err, "Error loading kubeconfig from assets.")
	}

	kubeclient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "Creating a Kubernetes client from assets failed.")
	}

	zeroKubeClient.Client = kubeclient
	zeroKubeClient.ctx = ctx
	zeroKubeClient.config = kubeconfig
	zeroKubeClient.configPath = kubeconfigpath

	return zeroKubeClient, nil
}

func (zerokube *clusterZeroKubeAPIClient) IsKubeAPILive() (bool, error) {

	discovery := zerokube.Client.Discovery()
	version, err := discovery.ServerVersion()
	if err != nil {
		return false, err
	}
	logrus.Infof("Cluster API is up and running %s", version)
	return true, nil
}

// DEV_NOTES(lranjbar): Potentially redundant? We will fail when making the client if kubeconfig is not around
func (zerokube *clusterZeroKubeAPIClient) DoesKubeConfigExist() (bool, error) {

	_, err := clientcmd.LoadFromFile(zerokube.configPath)
	if err != nil {
		return false, errors.Wrap(err, "Error loading kubeconfig from file.")
	}
	return true, nil
}

func (zerokube *clusterZeroKubeAPIClient) IsBootstrapConfigMapComplete() (bool, error) {

	// Get latest version of bootstrap configmap
	bootstrap, err := zerokube.Client.CoreV1().ConfigMaps("kube-system").Get(zerokube.ctx, "bootstrap", v1.GetOptions{})

	if err != nil {
		logrus.Debug("bootstrap configmap not found")
		return false, err
	}
	// Found a bootstrap configmap need to check its status
	if bootstrap != nil && err == nil {
		status, ok := bootstrap.Data["status"]
		if !ok {
			logrus.Debug("No status found in bootstrap configmap.")
			return false, nil
		}
		if status == "complete" {
			return true, nil
		}
	}
	return false, nil
}
