package zero

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

func NewClusterZeroKubeAPIClient(ctx context.Context, assertDir string) (*clusterZeroKubeAPIClient, error) {

	zeroKubeClient := &clusterZeroKubeAPIClient{}

	kubeconfigpath := filepath.Join(assertDir, "auth", "kubeconfig")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigpath)
	if err != nil {
		return nil, errors.Wrap(err, "loading kubeconfig")
	}

	kubeclient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating a Kubernetes client")
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
		return false, errors.Wrap(err, "loading kubeconfig")
	}
	return true, nil
}

func (zerokube *clusterZeroKubeAPIClient) DoesBootstrapConfigMapExist() (bool, error) {

	return true, nil
}
