package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterPoolSpec defines the desired state of the ClusterPool.
type ClusterPoolSpec struct {

	// Platform encompasses the desired platform for the cluster.
	// +required
	Platform Platform `json:"platform"`

	// PullSecretRef is the reference to the secret to use when pulling images.
	// +optional
	PullSecretRef *corev1.LocalObjectReference `json:"pullSecretRef,omitempty"`

	// Size is the default number of clusters that we should keep provisioned and waiting for use.
	// +kubebuilder:validation:Minimum=0
	// +required
	Size int32 `json:"size"`

	// RunningCount is the number of clusters we should keep running. The remainder will be kept hibernated until claimed.
	// By default no clusters will be kept running (all will be hibernated).
	// +kubebuilder:validation:Minimum=0
	// +optional
	RunningCount int32 `json:"runningCount,omitempty"`

	// MaxSize is the maximum number of clusters that will be provisioned including clusters that have been claimed
	// and ones waiting to be used.
	// By default there is no limit.
	// +optional
	MaxSize *int32 `json:"maxSize,omitempty"`

	// MaxConcurrent is the maximum number of clusters that will be provisioned or deprovisioned at an time. This includes the
	// claimed clusters being deprovisioned.
	// By default there is no limit.
	// +optional
	MaxConcurrent *int32 `json:"maxConcurrent,omitempty"`

	// BaseDomain is the base domain to use for all clusters created in this pool.
	// +required
	BaseDomain string `json:"baseDomain"`

	// ImageSetRef is a reference to a ClusterImageSet. The release image specified in the ClusterImageSet will be used
	// by clusters created for this cluster pool.
	ImageSetRef ClusterImageSetReference `json:"imageSetRef"`

	// Labels to be applied to new ClusterDeployments created for the pool. ClusterDeployments that have already been
	// claimed will not be affected when this value is modified.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations to be applied to new ClusterDeployments created for the pool. ClusterDeployments that have already been
	// claimed will not be affected when this value is modified.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// InstallConfigSecretTemplateRef is a secret with the key install-config.yaml consisting of the content of the install-config.yaml
	// to be used as a template for all clusters in this pool.
	// Cluster specific settings (name, basedomain) will be injected dynamically when the ClusterDeployment install-config Secret is generated.
	// +optional
	InstallConfigSecretTemplateRef *corev1.LocalObjectReference `json:"installConfigSecretTemplateRef,omitempty"`

	// HibernateAfter will be applied to new ClusterDeployments created for the pool. HibernateAfter will transition
	// clusters in the clusterpool to hibernating power state after it has been running for the given duration. The time
	// that a cluster has been running is the time since the cluster was installed or the time since the cluster last came
	// out of hibernation.
	// This is a Duration value; see https://pkg.go.dev/time#ParseDuration for accepted formats.
	// Note: due to discrepancies in validation vs parsing, we use a Pattern instead of `Format=duration`. See
	// https://bugzilla.redhat.com/show_bug.cgi?id=2050332
	// https://github.com/kubernetes/apimachinery/issues/131
	// https://github.com/kubernetes/apiextensions-apiserver/issues/56
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	HibernateAfter *metav1.Duration `json:"hibernateAfter,omitempty"`

	// InstallAttemptsLimit is the maximum number of times Hive will attempt to install the cluster.
	// +optional
	InstallAttemptsLimit *int32 `json:"installAttemptsLimit,omitempty"`

	// SkipMachinePools allows creating clusterpools where the machinepools are not managed by hive after cluster creation
	// +optional
	SkipMachinePools bool `json:"skipMachinePools,omitempty"`

	// ClaimLifetime defines the lifetimes for claims for the cluster pool.
	// +optional
	ClaimLifetime *ClusterPoolClaimLifetime `json:"claimLifetime,omitempty"`

	// HibernationConfig configures the hibernation/resume behavior of ClusterDeployments owned by the ClusterPool.
	// +optional
	HibernationConfig *HibernationConfig `json:"hibernationConfig"`

	// Inventory maintains a list of entries consumed by the ClusterPool
	// to customize the default ClusterDeployment.
	// +optional
	Inventory []InventoryEntry `json:"inventory,omitempty"`
}

type HibernationConfig struct {
	// ResumeTimeout is the maximum amount of time we will wait for an unclaimed ClusterDeployment to resume from
	// hibernation (e.g. at the behest of runningCount, or in preparation for being claimed). If this time is
	// exceeded, the ClusterDeployment will be considered Broken and we will replace it. The default (unspecified
	// or zero) means no timeout -- we will allow the ClusterDeployment to continue trying to resume "forever".
	// This is a Duration value; see https://pkg.go.dev/time#ParseDuration for accepted formats.
	// Note: due to discrepancies in validation vs parsing, we use a Pattern instead of `Format=duration`. See
	// https://bugzilla.redhat.com/show_bug.cgi?id=2050332
	// https://github.com/kubernetes/apimachinery/issues/131
	// https://github.com/kubernetes/apiextensions-apiserver/issues/56
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	ResumeTimeout metav1.Duration `json:"resumeTimeout"`
}

// InventoryEntryKind is the Kind of the inventory entry.
// +kubebuilder:validation:Enum="";ClusterDeploymentCustomization
type InventoryEntryKind string

const ClusterDeploymentCustomizationInventoryEntry InventoryEntryKind = "ClusterDeploymentCustomization"

// InventoryEntry maintains a reference to a custom resource consumed by a clusterpool to customize the cluster deployment.
type InventoryEntry struct {
	// Kind denotes the kind of the referenced resource. The default is ClusterDeploymentCustomization, which is also currently the only supported value.
	// +kubebuilder:default=ClusterDeploymentCustomization
	Kind InventoryEntryKind `json:"kind,omitempty"`
	// Name is the name of the referenced resource.
	// +required
	Name string `json:"name,omitempty"`
}

// ClusterPoolClaimLifetime defines the lifetimes for claims for the cluster pool.
type ClusterPoolClaimLifetime struct {
	// Default is the default lifetime of the claim when no lifetime is set on the claim itself.
	// This is a Duration value; see https://pkg.go.dev/time#ParseDuration for accepted formats.
	// Note: due to discrepancies in validation vs parsing, we use a Pattern instead of `Format=duration`. See
	// https://bugzilla.redhat.com/show_bug.cgi?id=2050332
	// https://github.com/kubernetes/apimachinery/issues/131
	// https://github.com/kubernetes/apiextensions-apiserver/issues/56
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	Default *metav1.Duration `json:"default,omitempty"`

	// Maximum is the maximum lifetime of the claim after it is assigned a cluster. If the claim still exists
	// when the lifetime has elapsed, the claim will be deleted by Hive.
	// The lifetime of a claim is the mimimum of the lifetimes set by the cluster pool and the claim itself.
	// This is a Duration value; see https://pkg.go.dev/time#ParseDuration for accepted formats.
	// Note: due to discrepancies in validation vs parsing, we use a Pattern instead of `Format=duration`. See
	// https://bugzilla.redhat.com/show_bug.cgi?id=2050332
	// https://github.com/kubernetes/apimachinery/issues/131
	// https://github.com/kubernetes/apiextensions-apiserver/issues/56
	// +optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	Maximum *metav1.Duration `json:"maximum,omitempty"`
}

// ClusterPoolStatus defines the observed state of ClusterPool
type ClusterPoolStatus struct {
	// Size is the number of unclaimed clusters that have been created for the pool.
	Size int32 `json:"size"`

	// Standby is the number of unclaimed clusters that are installed, but not running.
	// +optional
	Standby int32 `json:"standby"`

	// Ready is the number of unclaimed clusters that are installed and are running and ready to be claimed.
	Ready int32 `json:"ready"`

	// Conditions includes more detailed status for the cluster pool
	// +optional
	Conditions []ClusterPoolCondition `json:"conditions,omitempty"`
}

// ClusterPoolCondition contains details for the current condition of a cluster pool
type ClusterPoolCondition struct {
	// Type is the type of the condition.
	Type ClusterPoolConditionType `json:"type"`
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

// ClusterPoolConditionType is a valid value for ClusterPoolCondition.Type
type ClusterPoolConditionType string

// ConditionType satisfies the conditions.Condition interface
func (c ClusterPoolCondition) ConditionType() ConditionType {
	return c.Type
}

// String satisfies the conditions.ConditionType interface
func (t ClusterPoolConditionType) String() string {
	return string(t)
}

const (
	// ClusterPoolMissingDependenciesCondition is set when a cluster pool is missing dependencies required to create a
	// cluster. Dependencies include resources such as the ClusterImageSet and the credentials Secret.
	ClusterPoolMissingDependenciesCondition ClusterPoolConditionType = "MissingDependencies"
	// ClusterPoolCapacityAvailableCondition is set to provide information on whether the cluster pool has capacity
	// available to create more clusters for the pool.
	ClusterPoolCapacityAvailableCondition ClusterPoolConditionType = "CapacityAvailable"
	// ClusterPoolAllClustersCurrentCondition indicates whether all unassigned (installing or ready)
	// ClusterDeployments in the pool match the current configuration of the ClusterPool.
	ClusterPoolAllClustersCurrentCondition ClusterPoolConditionType = "AllClustersCurrent"
	// ClusterPoolInventoryValidCondition is set to provide information on whether the cluster pool inventory is valid.
	ClusterPoolInventoryValidCondition ClusterPoolConditionType = "InventoryValid"
	// ClusterPoolDeletionPossibleCondition gives information about a deleted ClusterPool which is pending cleanup.
	// Note that it is normal for this condition to remain Initialized/Unknown until the ClusterPool is deleted.
	ClusterPoolDeletionPossibleCondition ClusterPoolConditionType = "DeletionPossible"
)

const (
	// InventoryReasonValid is used when all ClusterDeploymentCustomization are
	// available and when used the ClusterDeployments are successfully installed.
	InventoryReasonValid = "Valid"
	// InventoryReasonInvalid is used when there is something wrong with ClusterDeploymentCustomization, for example
	// patching issue, provisioning failure, missing, etc.
	InventoryReasonInvalid = "Invalid"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterPool represents a pool of clusters that should be kept ready to be given out to users. Clusters are removed
// from the pool once claimed and then automatically replaced with a new one.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.size,statuspath=.status.size
// +kubebuilder:printcolumn:name="Size",type="string",JSONPath=".spec.size"
// +kubebuilder:printcolumn:name="Standby",type="string",JSONPath=".status.standby"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="BaseDomain",type="string",JSONPath=".spec.baseDomain"
// +kubebuilder:printcolumn:name="ImageSet",type="string",JSONPath=".spec.imageSetRef.name"
// +kubebuilder:resource:path=clusterpools,shortName=cp
type ClusterPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterPoolSpec   `json:"spec"`
	Status ClusterPoolStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterPoolList contains a list of ClusterPools
type ClusterPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterPool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterPool{}, &ClusterPoolList{})
}
