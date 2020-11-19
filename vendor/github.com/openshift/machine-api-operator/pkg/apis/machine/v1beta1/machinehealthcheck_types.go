package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// RemediationStrategyType contains remediation strategy type
type RemediationStrategyType string

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineHealthCheck is the Schema for the machinehealthchecks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=mhc;mhcs
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="MaxUnhealthy",type="string",JSONPath=".spec.maxUnhealthy",description="Maximum number of unhealthy machines allowed"
// +kubebuilder:printcolumn:name="ExpectedMachines",type="integer",JSONPath=".status.expectedMachines",description="Number of machines currently monitored"
// +kubebuilder:printcolumn:name="CurrentHealthy",type="integer",JSONPath=".status.currentHealthy",description="Current observed healthy machines"
type MachineHealthCheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of machine health check policy
	Spec MachineHealthCheckSpec `json:"spec,omitempty"`

	// Most recently observed status of MachineHealthCheck resource
	Status MachineHealthCheckStatus `json:"status,omitempty"`
}

func (m *MachineHealthCheck) GetConditions() Conditions {
	return m.Status.Conditions
}

func (m *MachineHealthCheck) SetConditions(conditions Conditions) {
	m.Status.Conditions = conditions
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MachineHealthCheckList contains a list of MachineHealthCheck
type MachineHealthCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineHealthCheck `json:"items"`
}

// MachineHealthCheckSpec defines the desired state of MachineHealthCheck
type MachineHealthCheckSpec struct {
	// Label selector to match machines whose health will be exercised.
	// Note: An empty selector will match all machines.
	Selector metav1.LabelSelector `json:"selector"`

	// UnhealthyConditions contains a list of the conditions that determine
	// whether a node is considered unhealthy.  The conditions are combined in a
	// logical OR, i.e. if any of the conditions is met, the node is unhealthy.
	//
	// +kubebuilder:validation:MinItems=1
	UnhealthyConditions []UnhealthyCondition `json:"unhealthyConditions"`

	// Any farther remediation is only allowed if at most "MaxUnhealthy" machines selected by
	// "selector" are not healthy.
	// Expects either a postive integer value or a percentage value.
	// Percentage values must be positive whole numbers and are capped at 100%.
	// Both 0 and 0% are valid and will block all remediation.
	// +kubebuilder:default:="100%"
	// +kubebuilder:validation:Pattern="^((100|[0-9]{1,2})%|[0-9]+)$"
	// +kubebuilder:validation:Type:=string
	MaxUnhealthy *intstr.IntOrString `json:"maxUnhealthy,omitempty"`

	// Machines older than this duration without a node will be considered to have
	// failed and will be remediated.
	// Expects an unsigned duration string of decimal numbers each with optional
	// fraction and a unit suffix, eg "300ms", "1.5h" or "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// +optional
	// +kubebuilder:default:="10m"
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	// +kubebuilder:validation:Type:=string
	NodeStartupTimeout metav1.Duration `json:"nodeStartupTimeout,omitempty"`
}

// UnhealthyCondition represents a Node condition type and value with a timeout
// specified as a duration.  When the named condition has been in the given
// status for at least the timeout value, a node is considered unhealthy.
type UnhealthyCondition struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	Type corev1.NodeConditionType `json:"type"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	Status corev1.ConditionStatus `json:"status"`

	// Expects an unsigned duration string of decimal numbers each with optional
	// fraction and a unit suffix, eg "300ms", "1.5h" or "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// +kubebuilder:validation:Pattern="^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	// +kubebuilder:validation:Type:=string
	Timeout metav1.Duration `json:"timeout"`
}

// MachineHealthCheckStatus defines the observed state of MachineHealthCheck
type MachineHealthCheckStatus struct {
	// total number of machines counted by this machine health check
	// +kubebuilder:validation:Minimum=0
	ExpectedMachines *int `json:"expectedMachines"`

	// total number of machines counted by this machine health check
	// +kubebuilder:validation:Minimum=0
	CurrentHealthy *int `json:"currentHealthy" protobuf:"varint,4,opt,name=currentHealthy"`

	// RemediationsAllowed is the number of further remediations allowed by this machine health check before
	// maxUnhealthy short circuiting will be applied
	// +kubebuilder:validation:Minimum=0
	RemediationsAllowed int32 `json:"remediationsAllowed,omitempty"`

	// Conditions defines the current state of the MachineHealthCheck
	Conditions Conditions `json:"conditions,omitempty"`
}
