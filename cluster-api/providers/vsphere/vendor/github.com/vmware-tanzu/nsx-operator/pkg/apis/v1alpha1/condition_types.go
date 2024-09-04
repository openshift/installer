package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConditionType string

const (
	Ready ConditionType = "Ready"
)

// Condition defines condition of custom resource.
type Condition struct {
	// Type defines condition type.
	Type ConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Last time the condition transitioned from one status to another.
	// This should be when the underlying condition changed. If that is not known, then using the time when
	// the API field changed is acceptable.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason shows a brief reason of condition.
	Reason string `json:"reason,omitempty"`
	// Message shows a human-readable message about condition.
	Message string `json:"message,omitempty"`
}
