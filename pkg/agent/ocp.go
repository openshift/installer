package agent

import (
	"context"
	"path/filepath"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientwatch "k8s.io/client-go/tools/watch"

	configclient "github.com/openshift/client-go/config/clientset/versioned"
	routeclient "github.com/openshift/client-go/route/clientset/versioned"
	cov1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
	"github.com/openshift/library-go/pkg/route/routeapihelpers"
)

// ClusterOpenShiftAPIClient Kube client using the OpenShift clientset instead of the Kubernetes clientset
type ClusterOpenShiftAPIClient struct {
	ConfigClient *configclient.Clientset
	RouteClient  *routeclient.Clientset
	ctx          context.Context
	config       *rest.Config
	configPath   string
}

const (
	// Need to keep these updated if they change
	consoleNamespace = "openshift-console"
	consoleRouteName = "console"
)

// NewClusterOpenShiftAPIClient Create a kube client with OCP understanding
func NewClusterOpenShiftAPIClient(ctx context.Context, assetDir string) (*ClusterOpenShiftAPIClient, error) {

	ocpClient := &ClusterOpenShiftAPIClient{}

	kubeconfigpath := filepath.Join(assetDir, "auth", "kubeconfig")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigpath)
	if err != nil {
		return nil, errors.Wrap(err, "creating kubeconfig for ocp config client")
	}

	configclient, err := configclient.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating an ocp config client")
	}

	routeclient, err := routeclient.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating an ocp route client")
	}

	ocpClient.ConfigClient = configclient
	ocpClient.RouteClient = routeclient
	ocpClient.ctx = ctx
	ocpClient.config = kubeconfig
	ocpClient.configPath = kubeconfigpath

	return ocpClient, nil

}

// AreClusterOperatorsInitalized Wait for all Openshift cluster operators to initialize
func (ocp *ClusterOpenShiftAPIClient) AreClusterOperatorsInitalized(waitctx context.Context) (bool, error) {
	failing := configv1.ClusterStatusConditionType("Failing")
	var lastError string

	_, err := clientwatch.UntilWithSync(
		waitctx,
		cache.NewListWatchFromClient(ocp.ConfigClient.ConfigV1().RESTClient(), "clusterversions", "", fields.OneTermEqualSelector("metadata.name", "version")),
		&configv1.ClusterVersion{},
		nil,
		func(event watch.Event) (bool, error) {
			switch event.Type {
			case watch.Added, watch.Modified:
				cv, ok := event.Object.(*configv1.ClusterVersion)
				if !ok {
					logrus.Warnf("Expected a ClusterVersion object but got a %q object instead", event.Object.GetObjectKind().GroupVersionKind())
					return false, nil
				}
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorAvailable) &&
					cov1helpers.IsStatusConditionFalse(cv.Status.Conditions, failing) &&
					cov1helpers.IsStatusConditionFalse(cv.Status.Conditions, configv1.OperatorProgressing) {
					logrus.Debug("Cluster operators intitalized")
					return true, nil
				}
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, failing) {
					lastError = cov1helpers.FindStatusCondition(cv.Status.Conditions, failing).Message
				} else if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorProgressing) {
					lastError = cov1helpers.FindStatusCondition(cv.Status.Conditions, configv1.OperatorProgressing).Message
				}
				logrus.Debugf("Still waiting for the cluster to initialize: %s", lastError)
				return false, nil
			}
			logrus.Debug("Still waiting for the cluster to initialize...")
			return false, nil
		},
	)

	if lastError != "" {
		if err == wait.ErrWaitTimeout {
			return false, errors.Errorf("failed to initialize the cluster: %s", lastError)
		}

		return false, errors.Wrapf(err, "failed to initialize the cluster: %s", lastError)
	}

	return false, errors.Wrap(err, "failed to initialize the cluster")
}

// IsConsoleRouteAvaiable Check if the OCP console route is created
func (ocp *ClusterOpenShiftAPIClient) IsConsoleRouteAvaiable() (bool, error) {
	route, err := ocp.RouteClient.RouteV1().Routes(consoleNamespace).Get(ocp.ctx, consoleRouteName, metav1.GetOptions{})
	if err == nil {
		logrus.Debugf("Route found in openshift-console namespace: %s", consoleRouteName)
		if _, _, err2 := routeapihelpers.IngressURI(route, ""); err2 == nil {
			logrus.Debug("OpenShift console route is admitted")
			return true, nil
		} else if err2 != nil {
			err = err2
		}
	}
	return false, errors.Wrap(err, "waiting for openshift-console route")

}

// IsConsoleRouteURLAvailable Check if the console route URL is available
func (ocp *ClusterOpenShiftAPIClient) IsConsoleRouteURLAvailable() (bool, string, error) {
	url := ""
	route, err := ocp.RouteClient.RouteV1().Routes(consoleNamespace).Get(ocp.ctx, consoleRouteName, metav1.GetOptions{})
	if err == nil {
		if uri, _, err2 := routeapihelpers.IngressURI(route, ""); err2 == nil {
			url = uri.String()
		} else {
			err = err2
		}
	}
	if url == "" {
		return false, url, errors.Wrap(err, "waiting for openshift-console URL")
	}
	return true, url, nil
}

// LogClusterOperatorConditions Log OCP cluster operator conditions
func (ocp *ClusterOpenShiftAPIClient) LogClusterOperatorConditions() error {

	operators, err := ocp.ConfigClient.ConfigV1().ClusterOperators().List(ocp.ctx, metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "listing ClusterOperator objects")
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
