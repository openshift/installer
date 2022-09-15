package agent

import (
	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/types/baremetal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AgentConfigVersion is the version supported by this package.
// If you bump this, you must also update the list of convertable values in
// pkg/types/conversion/agentconfig.go
const AgentConfigVersion = "v1alpha1"

// Config or aka AgentConfig is the API for specifying additional
// configuration for the agent-based installer not covered by
// install-config.
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

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
