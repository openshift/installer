package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Common types that can be used by all ClusterInstall implementations.

// ClusterInstallCondition contains details for the current condition of a cluster install.
type ClusterInstallCondition struct {
	// Type is the type of the condition.
	Type ClusterInstallConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

type ClusterInstallConditionType string

// ConditionType satisfies the conditions.Condition interface
func (c ClusterInstallCondition) ConditionType() ConditionType {
	return c.Type
}

// String satisfies the conditions.ConditionType interface
func (t ClusterInstallConditionType) String() string {
	return string(t)
}

const (
	// ClusterInstallRequirementsMet is True when all pre-install requirements have been met.
	ClusterInstallRequirementsMet ClusterInstallConditionType = "RequirementsMet"

	// ClusterInstallCompleted is True when the requested install has been completed successfully.
	ClusterInstallCompleted ClusterInstallConditionType = "Completed"

	// ClusterInstallFailed is True when an attempt to install the cluster has failed.
	// The ClusterInstall controllers may still be retrying if supported, and this condition will
	// go back to False if a later attempt succeeds.
	ClusterInstallFailed ClusterInstallConditionType = "Failed"

	// ClusterInstallStopped is True the controllers are no longer working on this
	// ClusterInstall. Combine with Completed or Failed to know if the overall request was
	// successful or not.
	ClusterInstallStopped ClusterInstallConditionType = "Stopped"
)
