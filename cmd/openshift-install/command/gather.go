package command

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
)

// LogClusterOperatorConditions logs the current status of cluster operators that
// are still becoming Available.
func LogClusterOperatorConditions(ctx context.Context, config *rest.Config) error {
	client, err := configclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("creating a config client: %w", err)
	}

	operators, err := client.ConfigV1().ClusterOperators().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing ClusterOperator objects: %w", err)
	}

	for _, operator := range operators.Items {
		for _, condition := range operator.Status.Conditions {
			switch {
			case condition.Type == configv1.OperatorUpgradeable:
				continue
			case condition.Type == configv1.OperatorAvailable && condition.Status == configv1.ConditionTrue:
				continue
			case (condition.Type == configv1.OperatorDegraded || condition.Type == configv1.OperatorProgressing) && condition.Status == configv1.ConditionFalse:
				continue
			}
			if condition.Type == configv1.OperatorAvailable || condition.Type == configv1.OperatorDegraded {
				logrus.Errorf("Cluster operator %s %s is %s with %s: %s", operator.ObjectMeta.Name, condition.Type, condition.Status, condition.Reason, condition.Message)
			} else {
				logrus.Infof("Cluster operator %s %s is %s with %s: %s", operator.ObjectMeta.Name, condition.Type, condition.Status, condition.Reason, condition.Message)
			}
		}
	}

	return nil
}
