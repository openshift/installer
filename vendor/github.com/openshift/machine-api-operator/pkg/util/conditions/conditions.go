package conditions

import (
	machinev1 "github.com/openshift/api/machine/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// GetNodeCondition returns node condition by type
func GetNodeCondition(node *corev1.Node, conditionType corev1.NodeConditionType) *corev1.NodeCondition {
	for _, cond := range node.Status.Conditions {
		if cond.Type == conditionType {
			return &cond
		}
	}
	return nil
}

// GetDeploymentCondition returns node condition by type
func GetDeploymentCondition(deployment *appsv1.Deployment, conditionType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for _, cond := range deployment.Status.Conditions {
		if cond.Type == conditionType {
			return &cond
		}
	}
	return nil
}

func DeepCopyConditions(in []machinev1.Condition) []machinev1.Condition {
	out := make([]machinev1.Condition, 0)
	for _, cond := range in {
		out = append(out, *cond.DeepCopy())
	}
	return out
}
