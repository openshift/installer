package agent

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	configv1 "github.com/openshift/api/config/v1"
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
	cvResVersion string
}

const (
	// Need to keep these updated if they change
	consoleNamespace = "openshift-console"
	consoleRouteName = "console"
)

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

// AreClusterOperatorsInitialized Waits for all Openshift cluster operators to initialize
func (ocp *ClusterOpenShiftAPIClient) AreClusterOperatorsInitialized() (bool, error) {

	var lastError string
	failing := configv1.ClusterStatusConditionType("Failing")

	version, err := ocp.ConfigClient.ConfigV1().ClusterVersions().Get(ocp.ctx, "version", metav1.GetOptions{})
	if err != nil {
		return false, errors.Wrap(err, "Getting ClusterVersion object")
	}

	if cov1helpers.IsStatusConditionTrue(version.Status.Conditions, configv1.OperatorAvailable) &&
		cov1helpers.IsStatusConditionFalse(version.Status.Conditions, failing) &&
		cov1helpers.IsStatusConditionFalse(version.Status.Conditions, configv1.OperatorProgressing) {
		return true, nil
	}

	if cov1helpers.IsStatusConditionTrue(version.Status.Conditions, failing) {
		lastError = cov1helpers.FindStatusCondition(version.Status.Conditions, failing).Message
	} else if cov1helpers.IsStatusConditionTrue(version.Status.Conditions, configv1.OperatorProgressing) {
		lastError = cov1helpers.FindStatusCondition(version.Status.Conditions, configv1.OperatorProgressing).Message
	}
	if version.ResourceVersion != ocp.cvResVersion {
		logrus.Debugf("Still waiting for the cluster to initialize: %s", lastError)
		ocp.cvResVersion = version.ResourceVersion
	}

	return false, nil
}

// IsConsoleRouteAvailable Check if the OCP console route is created
func (ocp *ClusterOpenShiftAPIClient) IsConsoleRouteAvailable() (bool, error) {
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
	return false, errors.Wrap(err, "Waiting for openshift-console route")
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
		return false, url, errors.Wrap(err, "Waiting for openshift-console URL")
	}
	return true, url, nil
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
