package v1

import (
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// CustomizationApplyReasonSucceeded indicates that the customization
	// worked properly on the last applied cluster deployment.
	CustomizationApplyReasonSucceeded = "Succeeded"
	// CustomizationApplyReasonBrokenSyntax indicates that Hive failed to apply
	// customization patches on install-config. More details would be found in
	// ApplySucceded condition message.
	CustomizationApplyReasonBrokenSyntax = "BrokenBySyntax"
	// CustomizationApplyReasonBrokenCloud indicates that cluster deployment provision has failed
	// when using this customization. More details would be found in the ApplySucceeded condition message.
	CustomizationApplyReasonBrokenCloud = "BrokenByCloud"
	// CustomizationApplyReasonInstallationPending indicates that the customization patches have
	// been successfully applied but provisioning is not completed yet.
	CustomizationApplyReasonInstallationPending = "InstallationPending"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterDeploymentCustomization is the Schema for clusterdeploymentcustomizations API.
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced
type ClusterDeploymentCustomization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterDeploymentCustomizationSpec   `json:"spec"`
	Status ClusterDeploymentCustomizationStatus `json:"status,omitempty"`
}

// ClusterDeploymentCustomizationSpec defines the desired state of ClusterDeploymentCustomization.
type ClusterDeploymentCustomizationSpec struct {
	// InstallConfigPatches is a list of patches to be applied to the install-config.
	InstallConfigPatches []PatchEntity `json:"installConfigPatches,omitempty"`
}

// PatchEntity represent a json patch (RFC 6902) to be applied to the install-config
type PatchEntity struct {
	// Op is the operation to perform: add, remove, replace, move, copy, test
	// +required
	Op string `json:"op"`
	// Path is the json path to the value to be modified
	// +required
	Path string `json:"path"`
	// From is the json path to copy or move the value from
	// +optional
	From string `json:"from,omitempty"`
	// Value is the value to be used in the operation
	// +required
	Value string `json:"value"`
}

// ClusterDeploymentCustomizationStatus defines the observed state of ClusterDeploymentCustomization.
type ClusterDeploymentCustomizationStatus struct {
	// ClusterDeploymentRef is a reference to the cluster deployment that this customization is applied on.
	// +optional
	ClusterDeploymentRef *corev1.LocalObjectReference `json:"clusterDeploymentRef,omitempty"`

	// ClusterPoolRef is the name of the current cluster pool the CDC used at.
	// +optional
	ClusterPoolRef *corev1.LocalObjectReference `json:"clusterPoolRef,omitempty"`

	// LastAppliedConfiguration contains the last applied patches to the install-config.
	// The information will retain for reference in case the customization is updated.
	// +optional
	LastAppliedConfiguration string `json:"lastAppliedConfiguration,omitempty"`

	// Conditions describes the state of the operator's reconciliation functionality.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []conditionsv1.Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`
}

const (
	ApplySucceededCondition conditionsv1.ConditionType = "ApplySucceeded"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterDeploymentCustomizationList contains a list of ClusterDeploymentCustomizations.
type ClusterDeploymentCustomizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterDeploymentCustomization `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterDeploymentCustomization{}, &ClusterDeploymentCustomizationList{})
}
