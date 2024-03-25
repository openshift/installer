package agent

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	configv1 "github.com/openshift/api/config/v1"
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
func NewClusterOpenShiftAPIClient(ctx context.Context, assetDir string) (*ClusterOpenShiftAPIClient, error) {

	ocpClient := &ClusterOpenShiftAPIClient{}

	kubeconfigpath := filepath.Join(assetDir, "auth", "kubeconfig")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigpath)
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
	ocpClient.configPath = kubeconfigpath

	return ocpClient, nil

}

// LogClusterOperatorConditions Log OCP cluster operator conditions
func (ocp *ClusterOpenShiftAPIClient) LogClusterOperatorConditions() error {

	operators, err := ocp.ConfigClient.ConfigV1().ClusterOperators().List(ocp.ctx, metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "Listing ClusterOperator objects")
	}

	for _, operator := range operators.Items {
		for _, condition := range operator.Status.Conditions {
			if condition.Type == configv1.OperatorUpgradeable {
				continue
			} else if condition.Type == configv1.OperatorAvailable && condition.Status == configv1.ConditionTrue {
				continue
			} else if (condition.Type == configv1.OperatorDegraded || condition.Type == configv1.OperatorProgressing) && condition.Status == configv1.ConditionFalse {
				continue
			}
			if condition.Type == configv1.OperatorDegraded {
				logrus.Errorf("Cluster operator %s %s is %s with %s: %s", operator.ObjectMeta.Name, condition.Type, condition.Status, condition.Reason, condition.Message)
			} else {
				logrus.Infof("Cluster operator %s %s is %s with %s: %s", operator.ObjectMeta.Name, condition.Type, condition.Status, condition.Reason, condition.Message)
			}
		}
	}

	return nil
}
