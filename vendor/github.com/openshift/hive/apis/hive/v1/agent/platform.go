package agent

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// BareMetalPlatform defines agent based install configuration specific to bare metal clusters.
// Can only be used with spec.installStrategy.agent.
type BareMetalPlatform struct {

	// AgentSelector is a label selector used for associating relevant custom resources with this cluster.
	// (Agent, BareMetalHost, etc)
	AgentSelector metav1.LabelSelector `json:"agentSelector"`
}
