package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/hive/apis/hive/v1/aws"
	"github.com/openshift/hive/apis/hive/v1/azure"
	"github.com/openshift/hive/apis/hive/v1/gcp"
	"github.com/openshift/hive/apis/hive/v1/ibmcloud"
	"github.com/openshift/hive/apis/hive/v1/openstack"
	"github.com/openshift/hive/apis/hive/v1/ovirt"
	"github.com/openshift/hive/apis/hive/v1/vsphere"
)

const (
	// MachinePoolImageIDOverrideAnnotation can be applied to MachinePools to control the precise image ID to be used
	// for the MachineSets we reconcile for this pool. This feature is presently only implemented for AWS, and
	// is intended for very limited use cases we do not recommend pursuing regularly. As such it is not currently
	// part of our official API.
	MachinePoolImageIDOverrideAnnotation = "hive.openshift.io/image-id-override"
)

// MachinePoolSpec defines the desired state of MachinePool
type MachinePoolSpec struct {

	// ClusterDeploymentRef references the cluster deployment to which this
	// machine pool belongs.
	ClusterDeploymentRef corev1.LocalObjectReference `json:"clusterDeploymentRef"`

	// Name is the name of the machine pool.
	Name string `json:"name"`

	// Replicas is the count of machines for this machine pool.
	// Replicas and autoscaling cannot be used together.
	// Default is 1, if autoscaling is not used.
	// +optional
	Replicas *int64 `json:"replicas,omitempty"`

	// Autoscaling is the details for auto-scaling the machine pool.
	// Replicas and autoscaling cannot be used together.
	// +optional
	Autoscaling *MachinePoolAutoscaling `json:"autoscaling,omitempty"`

	// Platform is configuration for machine pool specific to the platform.
	Platform MachinePoolPlatform `json:"platform"`

	// Map of label string keys and values that will be applied to the created MachineSet's
	// MachineSpec. This list will overwrite any modifications made to Node labels on an
	// ongoing basis.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// List of taints that will be applied to the created MachineSet's MachineSpec.
	// This list will overwrite any modifications made to Node taints on an ongoing basis.
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`
}

// MachinePoolAutoscaling details how the machine pool is to be auto-scaled.
type MachinePoolAutoscaling struct {
	// MinReplicas is the minimum number of replicas for the machine pool.
	MinReplicas int32 `json:"minReplicas"`

	// MaxReplicas is the maximum number of replicas for the machine pool.
	MaxReplicas int32 `json:"maxReplicas"`
}

// MachinePoolPlatform is the platform-specific configuration for a machine
// pool. Only one of the platforms should be set.
type MachinePoolPlatform struct {
	// AWS is the configuration used when installing on AWS.
	AWS *aws.MachinePoolPlatform `json:"aws,omitempty"`
	// Azure is the configuration used when installing on Azure.
	Azure *azure.MachinePool `json:"azure,omitempty"`
	// GCP is the configuration used when installing on GCP.
	GCP *gcp.MachinePool `json:"gcp,omitempty"`
	// OpenStack is the configuration used when installing on OpenStack.
	OpenStack *openstack.MachinePool `json:"openstack,omitempty"`
	// VSphere is the configuration used when installing on vSphere
	VSphere *vsphere.MachinePool `json:"vsphere,omitempty"`
	// Ovirt is the configuration used when installing on oVirt.
	Ovirt *ovirt.MachinePool `json:"ovirt,omitempty"`
	// IBMCloud is the configuration used when installing on IBM Cloud.
	IBMCloud *ibmcloud.MachinePool `json:"ibmcloud,omitempty"`
}

// MachinePoolStatus defines the observed state of MachinePool
type MachinePoolStatus struct {
	// Replicas is the current number of replicas for the machine pool.
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// MachineSets is the status of the machine sets for the machine pool on the remote cluster.
	MachineSets []MachineSetStatus `json:"machineSets,omitempty"`

	// Conditions includes more detailed status for the cluster deployment
	// +optional
	Conditions []MachinePoolCondition `json:"conditions,omitempty"`
}

// MachineSetStatus is the status of a machineset in the remote cluster.
type MachineSetStatus struct {
	// Name is the name of the machine set.
	Name string `json:"name"`

	// Replicas is the current number of replicas for the machine set.
	Replicas int32 `json:"replicas"`

	// The number of ready replicas for this MachineSet. A machine is considered ready
	// when the node has been created and is "Ready". It is transferred as-is from the
	// MachineSet from remote cluster.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// MinReplicas is the minimum number of replicas for the machine set.
	MinReplicas int32 `json:"minReplicas"`

	// MaxReplicas is the maximum number of replicas for the machine set.
	MaxReplicas int32 `json:"maxReplicas"`

	// In the event that there is a terminal problem reconciling the
	// replicas, both ErrorReason and ErrorMessage will be set. ErrorReason
	// will be populated with a succinct value suitable for machine
	// interpretation, while ErrorMessage will contain a more verbose
	// string suitable for logging and human consumption.
	// +optional
	ErrorReason *string `json:"errorReason,omitempty"`
	// +optional
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// MachinePoolCondition contains details for the current condition of a machine pool
type MachinePoolCondition struct {
	// Type is the type of the condition.
	Type MachinePoolConditionType `json:"type"`
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

// MachinePoolConditionType is a valid value for MachinePoolCondition.Type
type MachinePoolConditionType string

const (
	// NotEnoughReplicasMachinePoolCondition is true when the minReplicas field
	// is set too low for the number of machinesets for the machine pool.
	NotEnoughReplicasMachinePoolCondition MachinePoolConditionType = "NotEnoughReplicas"

	// NoMachinePoolNameLeasesAvailable is true when the cloud provider requires a name lease for the in-cluster MachineSet, but no
	// leases are available.
	NoMachinePoolNameLeasesAvailable MachinePoolConditionType = "NoMachinePoolNameLeasesAvailable"

	// InvalidSubnetsMachinePoolCondition is true when there are missing or invalid entries in the subnet field
	InvalidSubnetsMachinePoolCondition MachinePoolConditionType = "InvalidSubnets"

	// UnsupportedConfigurationMachinePoolCondition is true when the configuration of the MachinePool is unsupported
	// by the cluster.
	UnsupportedConfigurationMachinePoolCondition MachinePoolConditionType = "UnsupportedConfiguration"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachinePool is the Schema for the machinepools API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas
// +kubebuilder:printcolumn:name="PoolName",type="string",JSONPath=".spec.name"
// +kubebuilder:printcolumn:name="ClusterDeployment",type="string",JSONPath=".spec.clusterDeploymentRef.name"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:resource:path=machinepools,scope=Namespaced
type MachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachinePoolSpec   `json:"spec,omitempty"`
	Status MachinePoolStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachinePoolList contains a list of MachinePool
type MachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachinePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachinePool{}, &MachinePoolList{})
}
