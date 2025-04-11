package agent

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// AgentConfigVersion is the version supported by this package.
// If you bump this, you must also update the list of convertable values in
// pkg/types/conversion/agentconfig.go
const AgentConfigVersion = "v1beta1"

// Config or aka AgentConfig is the API for specifying additional
// configuration for the agent-based installer not covered by
// install-config.
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// AdditionalNTPSources is a list of NTP sources (hostname or IP) to be added to all cluster
	// hosts. They are added to any NTP sources that were configured through other means.
	// +optional
	AdditionalNTPSources []string `json:"additionalNTPSources,omitempty"`
	// ip address of node0
	RendezvousIP         string `json:"rendezvousIP,omitempty"`
	BootArtifactsBaseURL string `json:"bootArtifactsBaseURL,omitempty"`
	Hosts                []Host `json:"hosts,omitempty"`
	// When MinimalISO is set to true, a minimal ISO that does not contain the rootfs will be generated.
	// By default a full ISO will be created, unless the platform is External, which generates a minimal ISO.
	MinimalISO bool `json:"minimalISO,omitempty"`
}

// Host defines per host configurations
type Host struct {
	Hostname        string                    `json:"hostname,omitempty"`
	Role            string                    `json:"role,omitempty"`
	RootDeviceHints baremetal.RootDeviceHints `json:"rootDeviceHints,omitempty"`
	// list of interfaces and mac addresses
	Interfaces    []*aiv1beta1.Interface `json:"interfaces,omitempty"`
	NetworkConfig aiv1beta1.NetConfig    `json:"networkConfig,omitempty"`
	BMC           baremetal.BMC
}
