// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPAMDisabledAnnotationKeyName is the name of the annotation added to
// GatewayClass resources that do not participate in net-operator's IPAM.
// The value does not need to be truthy; the presence of the key is what
// disables net-operator's IPAM for that GatewayClass.
const IPAMDisabledAnnotationKeyName = "netoperator.vmware.com/ipam-disabled"

type IPPoolConditionType string

const (
	// IPPoolFull condition is added when no more IPs are free in the pool.
	IPPoolFull IPPoolConditionType = "full"
	// IPPoolReady condition is added when IPPool has been realized.
	IPPoolReady IPPoolConditionType = "ready"
	// IPPoolFail condition is added when an error was encountered in realizing.
	IPPoolFail IPPoolConditionType = "failure"
)

// IPPoolCondition describes the state of a IPPool at a certain point.
type IPPoolCondition struct {
	// Type is the type of IPPool condition.
	Type IPPoolConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Machine understandable string that gives the reason for condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
}

// IPPoolSpec defines the desired state of IPPool
type IPPoolSpec struct {
	// StartingAddress represents the starting IP address of the pool.
	StartingAddress string `json:"startingAddress"`
	// AddressCount represents the number of IP addresses in the pool.
	AddressCount int64 `json:"addressCount"`
}

// IPPoolStatus defines the current state of IPPool.
type IPPoolStatus struct {
	// Conditions is an array of current observed IPPool conditions.
	Conditions []IPPoolCondition `json:"conditions,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// IPPool is the Schema for the ippools API.
// It represents a pool of IP addresses that are owned and managed by the IPPool controller.
// Provider specific networks can associate themselves with IPPool objects to use
// network operator's IPAM implementation.
type IPPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IPPoolSpec   `json:"spec,omitempty"`
	Status IPPoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IPPoolList contains a list of IPPool
type IPPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPPool `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&IPPool{}, &IPPoolList{})
}
