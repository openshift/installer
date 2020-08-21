package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PacketMachineProviderConfig contains Config for Packet machines.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PacketMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	ProjectID    string   `json:"projectID"`
	Facilities   []string `json:"facility"`
	InstanceType string   `json:"machineType"`
	Tags         []string `json:"tags,omitempty"`
	OS           string   `json:"os,omitempty"`
	BillingCycle string   `json:"billingCycle,omitempty"`
}
