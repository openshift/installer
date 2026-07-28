package v1

import (
	"fmt"

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
// +kubebuilder:resource:shortName=cdc,scope=Namespaced
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

	// InstallerManifestPatches is a list of patches to be applied to installer-generated manifests.
	InstallerManifestPatches []InstallerManifestPatch `json:"installerManifestPatches,omitempty"`
}

type InstallerManifestPatch struct {
	// ManifestSelector identifies one or more manifests to patch
	ManifestSelector ManifestSelector `json:"manifestSelector"`

	// Patches is a list of RFC6902 patches to apply to manifests identified by manifestSelector.
	Patches []PatchEntity `json:"patches"`
}

type ManifestSelector struct {
	// Glob is a file glob (per https://pkg.go.dev/path/filepath#Glob) identifying one or more
	// manifests. Paths should be relative to the installer's working directory. Examples:
	// - openshift/99_role-cloud-creds-secret-reader.yaml
	// - openshift/99_openshift-cluster-api_worker-machineset-*.yaml
	// - */*secret*
	// It is an error if a glob matches zero manifests.
	Glob string `json:"glob"`
}

// PatchEntity represents a json patch (RFC 6902) to be applied
type PatchEntity struct {
	// Op is the operation to perform.
	// +kubebuilder:validation:Enum=add;remove;replace;move;copy;test
	// +required
	Op string `json:"op"`
	// Path is the json path to the value to be modified
	// +required
	Path string `json:"path"`
	// From is the json path to copy or move the value from
	// +optional
	From string `json:"from,omitempty"`
	// Value is the *string* value to be used in the operation. For more complex values, use
	// ValueJSON.
	// +optional
	Value string `json:"value,omitempty"`
	// ValueJSON is a string representing a JSON object to be used in the operation. As such,
	// internal quotes must be escaped. If nonempty, Value is ignored.
	// +optional
	ValueJSON string `json:"valueJSON,omitempty"`
}

// Encode returns a string representation of the RFC6902 patching operation represented by the
// PatchEntity.
func (pe *PatchEntity) Encode() string {
	var val string
	// Prefer ValueJSON
	if len(pe.ValueJSON) != 0 {
		// ValueJSON will be a raw JSON object in the patch; don't quote it.
		val = pe.ValueJSON
	} else {
		// Value will be a JSON string value in the patch; quote it.
		val = fmt.Sprintf("%q", pe.Value)
	}
	// Is this overkill? Should we just assemble the whole thing and let jsonpatch figure it out?
	switch pe.Op {
	case "copy", "move":
		return fmt.Sprintf(
			`{"op": %q, "from": %q, "path": %q}`,
			pe.Op, pe.From, pe.Path)
	case "remove":
		return fmt.Sprintf(
			`{"op": %q, "path": %q}`,
			pe.Op, pe.Path)
	case "test", "replace", "add":
		return fmt.Sprintf(
			// val is already quoted iff necessary
			`{"op": %q, "path": %q, "value": %s}`,
			pe.Op, pe.Path, val)
	default:
		// This should never happen, but...
		return fmt.Sprintf(
			// val is already quoted iff necessary
			`{"op": %q, "path": %q, "from", %q, "value": %s}`,
			pe.Op, pe.Path, pe.From, val)
	}
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
	Conditions []metav1.Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`
}

const (
	ApplySucceededCondition = "ApplySucceeded"
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
