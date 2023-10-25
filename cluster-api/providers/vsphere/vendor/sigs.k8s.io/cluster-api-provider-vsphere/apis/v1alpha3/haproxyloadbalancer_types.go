/*
Copyright 2019 The Kubernetes Authors.

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

//nolint:gocritic,godot
package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// HAProxyLoadBalancerFinalizer allows a reconciler to clean up
	// resources associated with an HAProxyLoadBalancer before removing
	// it from the API server.
	HAProxyLoadBalancerFinalizer = "haproxyloadbalancer.infrastructure.cluster.x-k8s.io"
)

// HAProxyLoadBalancerSpec defines the desired state of HAProxyLoadBalancer.
type HAProxyLoadBalancerSpec struct {
	// VirtualMachineConfiguration is information used to deploy a load balancer
	// VM.
	VirtualMachineConfiguration VirtualMachineCloneSpec `json:"virtualMachineConfiguration"`

	// SSHUser specifies the name of a user that is granted remote access to the
	// deployed VM.
	// +optional
	User *SSHUser `json:"user,omitempty"`
}

// HAProxyLoadBalancerStatus defines the observed state of HAProxyLoadBalancer.
type HAProxyLoadBalancerStatus struct {
	// Ready indicates whether or not the load balancer is ready.
	//
	// This field is required as part of the Portable Load Balancer model and is
	// inspected via an unstructured reader by other controllers to determine
	// the status of the load balancer.
	//
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Address is the IP address or DNS name of the load balancer.
	//
	// This field is required as part of the Portable Load Balancer model and is
	// inspected via an unstructured reader by other controllers to determine
	// the status of the load balancer.
	//
	// +optional
	Address string `json:"address,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:deprecatedversion
// +kubebuilder:resource:path=haproxyloadbalancers,scope=Namespaced
// +kubebuilder:subresource:status

// HAProxyLoadBalancer is the Schema for the haproxyloadbalancers API
//
// Deprecated: This type will be removed in v1alpha4.
type HAProxyLoadBalancer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HAProxyLoadBalancerSpec   `json:"spec,omitempty"`
	Status HAProxyLoadBalancerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HAProxyLoadBalancerList contains a list of HAProxyLoadBalancer
//
// Deprecated: This type will be removed in one of the next releases.
type HAProxyLoadBalancerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HAProxyLoadBalancer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HAProxyLoadBalancer{}, &HAProxyLoadBalancerList{})
}
