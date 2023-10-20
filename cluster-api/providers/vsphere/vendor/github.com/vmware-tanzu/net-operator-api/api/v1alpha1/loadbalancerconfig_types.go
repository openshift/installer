// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClientSecretReference contains info to locate an object of Kind Secret
// which contains credential specifications for a load balancer.
type ClientSecretReference struct {
	// Name is the name of resource being referenced.
	Name string `json:"name"`
	// Namespace of the resource being referenced. If empty, cluster scoped
	// resource is assumed.
	// +kubebuilder:default:=default
	Namespace string `json:"namespace,omitempty"`
}

// LoadBalancerConfigConditionType is used as a typed string for representing
// LoadBalancerConfig.Status.Conditions.
type LoadBalancerConfigConditionType string

const (
	// LoadBalancerConfigReady is added when the LoadBalancerConfig object has been successfully realized
	LoadBalancerConfigReady LoadBalancerConfigConditionType = "Ready"
	// LoadBalancerConfigFailure is added if any failure is encountered while realizing LoadBalancerConfig object
	LoadBalancerConfigFailure LoadBalancerConfigConditionType = "Failure"
)

// LoadBalancerConfigCondition describes the state of a LoadBalancerConfig at a certain point
type LoadBalancerConfigCondition struct {
	// Type is the type of load balancer condition
	// Can be Ready or Failure
	Type LoadBalancerConfigConditionType `json:"type"`
	// Status is the status of the condition
	// Can be True, False, Unknown
	Status corev1.ConditionStatus `json:"status"`
	// Machine understandable string that gives the reason for the condition's last transition
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition
	// +optional
	Message string `json:"message,omitempty"`
	// Provides a timestamp for when the LoadBalancerConfig object last transitioned from one status to another
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" patchStrategy:"replace"`
}

// LoadBalancerConfigProviderReference represents the specific load balancer instance that needs to be configured
type LoadBalancerConfigProviderReference struct {
	// APIGroup is the group for the resource being referenced
	APIGroup string `json:"apiGroup"`
	// Kind is the type of resource being referenced
	Kind string `json:"kind"`
	// Name is the name of resource being referenced
	Name string `json:"name"`
	// API version of the referent
	APIVersion string `json:"apiVersion,omitempty"`
}

type LoadBalancerConfigType string

const (
	// LoadBalancerConfigTypeHAProxy is the LoadBalancerConfigType for HAProxy.
	LoadBalancerConfigTypeHAProxy LoadBalancerConfigType = "haproxy"

	// LoadBalancerConfigTypeAvi is the LoadBalancerConfigType for Avi.
	LoadBalancerConfigTypeAvi LoadBalancerConfigType = "avi"
)

// LoadBalancerConfigSpec defines the desired state of LoadBalancerConfig
type LoadBalancerConfigSpec struct {
	// Type describes type of load balancer. Supported value is haproxy
	// +kubebuilder:validation:Enum=haproxy;avi
	Type LoadBalancerConfigType `json:"type"`
	// ProviderRef is reference to a load balancer provider object that provides the details for this type of load balancer
	ProviderRef LoadBalancerConfigProviderReference `json:"providerRef"`
}

// LoadBalancerConfigStatus defines the observed state of LoadBalancerConfig
type LoadBalancerConfigStatus struct {
	// Conditions is an array of current observed load balancer conditions
	Conditions []LoadBalancerConfigCondition `json:"conditions,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// LoadBalancerConfig is the Schema for the LoadBalancerConfigs API
type LoadBalancerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoadBalancerConfigSpec   `json:"spec,omitempty"`
	Status LoadBalancerConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoadBalancerConfigList contains a list of LoadBalancerConfig
type LoadBalancerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoadBalancerConfig `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&LoadBalancerConfig{}, &LoadBalancerConfigList{})
}
