/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ImageCreatedReason       = "ImageCreated"
	ImageStateCreated        = "Image has been created"
	ImageCreationErrorReason = "ImageCreationError"
	ImageStateFailedToCreate = "Failed to create image"
	InfraEnvNameLabel        = "infraenvs.agent-install.openshift.io"
)

// ClusterReference represents a Cluster Reference. It has enough information to retrieve cluster
// in any namespace
type ClusterReference struct {
	// Name is unique within a namespace to reference a cluster resource.
	// +optional
	Name string `json:"name,omitempty"`
	// Namespace defines the space within which the cluster name must be unique.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

const (
	ImageCreatedCondition conditionsv1.ConditionType = "ImageCreated"
)

type InfraEnvSpec struct {
	// Proxy defines the proxy settings for agents and clusters that use the InfraEnv. If
	// unset, the agents and clusters will not be configured to use a proxy.
	// +optional
	Proxy *Proxy `json:"proxy,omitempty"`

	// AdditionalNTPSources is a list of NTP sources (hostname or IP) to be added to all cluster
	// hosts. They are added to any NTP sources that were configured through other means.
	// +optional
	AdditionalNTPSources []string `json:"additionalNTPSources,omitempty"`

	// SSHAuthorizedKey is a SSH public keys that will be added to all agents for use in debugging.
	// +optional
	SSHAuthorizedKey string `json:"sshAuthorizedKey,omitempty"`

	// PullSecretRef is the reference to the secret to use when pulling images.
	PullSecretRef *corev1.LocalObjectReference `json:"pullSecretRef"`

	// AgentLabels lists labels to apply to Agents that are associated with this InfraEnv upon
	// the creation of the Agents.
	// +optional
	AgentLabels map[string]string `json:"agentLabels,omitempty"`

	// NmstateConfigLabelSelector associates NMStateConfigs for hosts that are considered part
	// of this installation environment.
	// +optional
	NMStateConfigLabelSelector metav1.LabelSelector `json:"nmStateConfigLabelSelector,omitempty"`

	// ClusterRef is the reference to the single ClusterDeployment that will be installed from
	// this InfraEnv.
	// Future versions will allow for multiple ClusterDeployments and this reference will be
	// removed.
	// +optional
	ClusterRef *ClusterReference `json:"clusterRef,omitempty"`

	// Json formatted string containing the user overrides for the initial ignition config
	// +optional
	IgnitionConfigOverride string `json:"ignitionConfigOverride,omitempty"`

	// CpuArchitecture specifies the target CPU architecture. Default is x86_64
	// +kubebuilder:default=x86_64
	// +optional
	CpuArchitecture string `json:"cpuArchitecture,omitempty"`
}

// Proxy defines the proxy settings for agents and clusters that use the InfraEnv.
// At least one of HTTPProxy or HTTPSProxy is required.
type Proxy struct {
	// HTTPProxy is the URL of the proxy for HTTP requests.
	// +optional
	HTTPProxy string `json:"httpProxy,omitempty"`

	// HTTPSProxy is the URL of the proxy for HTTPS requests.
	// +optional
	HTTPSProxy string `json:"httpsProxy,omitempty"`

	// NoProxy is a comma-separated list of domains and CIDRs for which the proxy should not be
	// used.
	// +optional
	NoProxy string `json:"noProxy,omitempty"`
}

type InfraEnvStatus struct {
	// ISODownloadURL specifies an HTTP/S URL that contains a discovery ISO containing the
	// configuration from this InfraEnv.
	ISODownloadURL string                   `json:"isoDownloadURL,omitempty"`
	CreatedTime    *metav1.Time             `json:"createdTime,omitempty"`
	Conditions     []conditionsv1.Condition `json:"conditions,omitempty"`
	// AgentLabelSelector specifies the label that will be applied to Agents that boot from the
	// installation media of this InfraEnv. This is how a user would identify which agents are
	// associated with a particular InfraEnv.
	// +optional
	AgentLabelSelector metav1.LabelSelector `json:"agentLabelSelector,omitempty"`
	// InfraEnvDebugInfo includes information for debugging the installation process.
	// +optional
	InfraEnvDebugInfo InfraEnvDebugInfo `json:"debugInfo"`
}

type InfraEnvDebugInfo struct {
	// EventsURL specifies an HTTP/S URL that contains InfraEnv events
	// +optional
	EventsURL string `json:"eventsURL"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ISO Created At",type="string",JSONPath=".status.createdTime",description="The Discovery ISO creation time"
// +kubebuilder:printcolumn:name="ISO URL",type="string",JSONPath=".status.isoDownloadURL",description="The Discovery ISO download URL",priority=1

type InfraEnv struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InfraEnvSpec   `json:"spec,omitempty"`
	Status InfraEnvStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InfraEnvList contains a list of InfraEnvs
type InfraEnvList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InfraEnv `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InfraEnv{}, &InfraEnvList{})
}
