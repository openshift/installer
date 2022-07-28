package agent

import (
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/types/baremetal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config or aka AgentConfig is the API for specifying additional
// configuration for the agent-based installer not covered by
// install-config.
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec Spec `json:"spec,omitempty"`
	// No status
}

// Spec contains additional configuration for the agent-based installer
type Spec struct {
	// ip address of node0
	RendezvousIP string `json:"rendezvousIP,omitempty"`
	Hosts        []Host `json:"hosts,omitempty"`
}

// Host defines per host configurations
type Host struct {
	Hostname        string                    `json:"hostname,omitempty"`
	Role            string                    `json:"role,omitempty"`
	RootDeviceHints baremetal.RootDeviceHints `json:"rootDeviceHints,omitempty"`
	// list of interfaces and mac addresses
	Interfaces    []*aiv1beta1.Interface `json:"interfaces,omitempty"`
	NetworkConfig aiv1beta1.NetConfig    `json:"networkConfig,omitempty"`
}
