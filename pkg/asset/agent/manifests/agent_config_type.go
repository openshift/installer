package manifests

import (
	metal3v1alpha1 "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AgentConfig is the API for specifying additional configuration for the
// agent-based installer not covered by install-config.
type AgentConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AgentConfigSpec `json:"spec,omitempty"`
	// No status
}

// AgentConfigSpec contains additional configuration for the agent-based installer
type AgentConfigSpec struct {
	// node0 hostname
	Node0 string            `json:"node0,omitempty"`
	Nodes []AgentConfigNode `json:"nodes"`
}

// AgentConfigNode defines per host configurations
type AgentConfigNode struct {
	Hostname        string                         `json:"hostname,omitempty"`
	Role            string                         `json:"role"`
	RootDeviceHints metal3v1alpha1.RootDeviceHints `json:"rootDeviceHints"`
	// list of interfaces and mac addresses
	Interfaces []*aiv1beta1.Interface `json:"interfaces,omitempty"`
}
