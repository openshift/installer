package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Annotation constants
const (
	// ClusterIDLabel is the label that a machineset must have to identify the
	// cluster to which it belongs.
	ClusterIDLabel   = "machine.openshift.io/cluster-api-cluster"
	MachineRoleLabel = "machine.openshift.io/cluster-api-machine-role"
	MachineTypeLabel = "machine.openshift.io/cluster-api-machine-type"
)

// LibvirtMachineProviderConfig is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an Libvirt instance.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LibvirtMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	DomainMemory             int        `json:"domainMemory"`
	DomainVcpu               int        `json:"domainVcpu"`
	IgnKey                   string     `json:"ignKey"`
	Ignition                 *Ignition  `json:"ignition"`
	CloudInit                *CloudInit `json:"cloudInit"`
	Volume                   *Volume    `json:"volume"`
	NetworkInterfaceName     string     `json:"networkInterfaceName"`
	NetworkInterfaceHostname string     `json:"networkInterfaceHostname"`
	NetworkInterfaceAddress  string     `json:"networkInterfaceAddress"`
	NetworkUUID              string     `json:"networkUUID"`
	Autostart                bool       `json:"autostart"`
	URI                      string     `json:"uri"`
}

// Ignition contains location of ignition to be run during bootstrapping
type Ignition struct {
	// Ignition config to be run during bootstrapping
	UserDataSecret string `json:"userDataSecret"`
}

// CloudInit contains location of user data to be run during bootstrapping
type CloudInit struct {
	// Bash script to be run during bootstrapping
	UserDataSecret string `json:"userDataSecret"`
	// Allow to ssh into instance
	SSHAccess bool `json:"sshAccess"`
}

// Volume contains the info for the actuator to create a volume
type Volume struct {
	PoolName     string             `json:"poolName"`
	BaseVolumeID string             `json:"baseVolumeID"`
	VolumeName   string             `json:"volumeName"`
	VolumeSize   *resource.Quantity `json:"volumeSize,omitempty"`
}

// LibvirtClusterProviderConfig is the type that will be embedded in a Cluster.Spec.ProviderSpec field.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LibvirtClusterProviderConfig struct {
	metav1.TypeMeta `json:",inline"`
}

// LibvirtMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains Libvirt-specific status information.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LibvirtMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// InstanceID is the instance ID of the machine created in Libvirt
	InstanceID *string `json:"instanceID"`

	// InstanceState is the state of the Libvirt instance for this machine
	InstanceState *string `json:"instanceState"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []LibvirtMachineProviderCondition `json:"conditions"`
}

// LibvirtMachineProviderConditionType is a valid value for LibvirtMachineProviderCondition.Type
type LibvirtMachineProviderConditionType string

// Valid conditions for an Libvirt machine instance
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated LibvirtMachineProviderConditionType = "MachineCreated"
)

// LibvirtMachineProviderCondition is a condition in a LibvirtMachineProviderStatus
type LibvirtMachineProviderCondition struct {
	// Type is the type of the condition.
	Type LibvirtMachineProviderConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message"`
}

// LibvirtClusterProviderStatus is the type that will be embedded in a Cluster.Status.ProviderStatus field.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LibvirtClusterProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LibvirtMachineProviderConfigList contains a list of LibvirtMachineProviderConfig
type LibvirtMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LibvirtMachineProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LibvirtMachineProviderConfig{}, &LibvirtMachineProviderConfigList{}, &LibvirtMachineProviderStatus{})
}
