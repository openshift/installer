package agent

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	configclient "github.com/openshift/client-go/config/clientset/versioned"
	routeclient "github.com/openshift/client-go/route/clientset/versioned"
)

// ClusterOpenShiftAPIClient Kube client using the OpenShift clientset instead of the Kubernetes clientset
type ClusterOpenShiftAPIClient struct {
	ConfigClient *configclient.Clientset
	RouteClient  *routeclient.Clientset
	ctx          context.Context
	config       *rest.Config
	configPath   string
}

// NewClusterOpenShiftAPIClient Create a kube client with OCP understanding
func NewClusterOpenShiftAPIClient(ctx context.Context, kubeconfigPath string) (*ClusterOpenShiftAPIClient, error) {
	ocpClient := &ClusterOpenShiftAPIClient{}

	var kubeconfig *rest.Config
	var err error
	if kubeconfigPath != "" {
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		kubeconfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, errors.Wrap(err, "creating kubeconfig for ocp config client")
	}

	configClient, err := configclient.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating an ocp config client")
	}

	routeClient, err := routeclient.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating an ocp route client")
	}

	ocpClient.ConfigClient = configClient
	ocpClient.RouteClient = routeClient
	ocpClient.ctx = ctx
	ocpClient.config = kubeconfig
	ocpClient.configPath = kubeconfigPath

	return ocpClient, nil
}
