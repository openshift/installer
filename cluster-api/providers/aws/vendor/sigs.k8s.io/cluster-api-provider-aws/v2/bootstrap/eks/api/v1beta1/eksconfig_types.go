/*
Copyright 2021 The Kubernetes Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// EKSConfigSpec defines the desired state of Amazon EKS Bootstrap Configuration.
type EKSConfigSpec struct {
	// KubeletExtraArgs passes the specified kubelet args into the Amazon EKS machine bootstrap script
	// +optional
	KubeletExtraArgs map[string]string `json:"kubeletExtraArgs,omitempty"`
	// ContainerRuntime specify the container runtime to use when bootstrapping EKS.
	// +optional
	ContainerRuntime *string `json:"containerRuntime,omitempty"`
	//  DNSClusterIP overrides the IP address to use for DNS queries within the cluster.
	// +optional
	DNSClusterIP *string `json:"dnsClusterIP,omitempty"`
	// DockerConfigJson is used for the contents of the /etc/docker/daemon.json file. Useful if you want a custom config differing from the default one in the AMI.
	// This is expected to be a json string.
	// +optional
	DockerConfigJSON *string `json:"dockerConfigJson,omitempty"`
	// APIRetryAttempts is the number of retry attempts for AWS API call.
	// +optional
	APIRetryAttempts *int `json:"apiRetryAttempts,omitempty"`
	// PauseContainer allows customization of the pause container to use.
	// +optional
	PauseContainer *PauseContainer `json:"pauseContainer,omitempty"`
	// UseMaxPods  sets --max-pods for the kubelet when true.
	// +optional
	UseMaxPods *bool `json:"useMaxPods,omitempty"`

	// ServiceIPV6Cidr is the ipv6 cidr range of the cluster. If this is specified then
	// the ip family will be set to ipv6.
	// +optional
	ServiceIPV6Cidr *string `json:"serviceIPV6Cidr,omitempty"`
}

// PauseContainer contains details of pause container.
type PauseContainer struct {
	//  AccountNumber is the AWS account number to pull the pause container from.
	AccountNumber string `json:"accountNumber"`
	// Version is the tag of the pause container to use.
	Version string `json:"version"`
}

// EKSConfigStatus defines the observed state of the Amazon EKS Bootstrap Configuration.
type EKSConfigStatus struct {
	// Ready indicates the BootstrapData secret is ready to be consumed
	Ready bool `json:"ready,omitempty"`

	// DataSecretName is the name of the secret that stores the bootstrap data script.
	// +optional
	DataSecretName *string `json:"dataSecretName,omitempty"`

	// FailureReason will be set on non-retryable errors
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage will be set on non-retryable errors
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// ObservedGeneration is the latest generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions defines current service state of the EKSConfig.
	// +optional
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:resource:path=eksconfigs,scope=Namespaced,categories=cluster-api,shortName=eksc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Bootstrap configuration is ready"
// +kubebuilder:printcolumn:name="DataSecretName",type="string",JSONPath=".status.dataSecretName",description="Name of Secret containing bootstrap data"

// EKSConfig is the schema for the Amazon EKS Machine Bootstrap Configuration API.
type EKSConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSConfigSpec   `json:"spec,omitempty"`
	Status EKSConfigStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the EKSConfig resource.
func (r *EKSConfig) GetConditions() clusterv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the EKSConfig to the predescribed clusterv1beta1.Conditions.
func (r *EKSConfig) SetConditions(conditions clusterv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion

// EKSConfigList contains a list of EKSConfig.
type EKSConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EKSConfig{}, &EKSConfigList{})
}
