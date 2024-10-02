/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

// +kubebuilder:object:generate=true
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// RuleAction describes the action to be applied on traffic matching a rule.
type RuleAction string

const (
	// RuleActionAllow describes that the traffic matching the rule must be allowed.
	RuleActionAllow RuleAction = "Allow"
	// RuleActionDrop describes that the traffic matching the rule must be dropped.
	RuleActionDrop RuleAction = "Drop"
	// RuleActionReject indicates that the traffic matching the rule must be rejected and the
	// client will receive a response.
	RuleActionReject RuleAction = "Reject"
)

// RuleDirection specifies the direction of traffic.
type RuleDirection string

const (
	// RuleDirectionIn specifies that the direction of traffic must be ingress, equivalent to "Ingress".
	RuleDirectionIn RuleDirection = "In"
	// RuleDirectionIngress specifies that the direction of traffic must be ingress, equivalent to "In".
	RuleDirectionIngress RuleDirection = "Ingress"
	// RuleDirectionOut specifies that the direction of traffic must be egress, equivalent to "Egress".
	RuleDirectionOut RuleDirection = "Out"
	// RuleDirectionEgress specifies that the direction of traffic must be egress, equivalent to "Out".
	RuleDirectionEgress RuleDirection = "Egress"
)

// SecurityPolicySpec defines the desired state of SecurityPolicy.
type SecurityPolicySpec struct {
	// Priority defines the order of policy enforcement.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1000
	Priority int `json:"priority,omitempty"`
	// AppliedTo is a list of policy targets to apply rules.
	// Policy level 'Applied To' will take precedence over rule level.
	AppliedTo []SecurityPolicyTarget `json:"appliedTo,omitempty"`
	// Rules is a list of policy rules.
	Rules []SecurityPolicyRule `json:"rules,omitempty"`
}

// SecurityPolicyRule defines a rule of SecurityPolicy.
type SecurityPolicyRule struct {
	// Action specifies the action to be applied on the rule.
	Action *RuleAction `json:"action"`
	// AppliedTo is a list of rule targets.
	// Policy level 'Applied To' will take precedence over rule level.
	AppliedTo []SecurityPolicyTarget `json:"appliedTo,omitempty"`
	// Direction is the direction of the rule, including 'In' or 'Ingress', 'Out' or 'Egress'.
	Direction *RuleDirection `json:"direction"`
	// Sources defines the endpoints where the traffic is from. For ingress rule only.
	Sources []SecurityPolicyPeer `json:"sources,omitempty"`
	// Destinations defines the endpoints where the traffic is to. For egress rule only.
	Destinations []SecurityPolicyPeer `json:"destinations,omitempty"`
	// Ports is a list of ports to be matched.
	Ports []SecurityPolicyPort `json:"ports,omitempty"`
	// Name is the display name of this rule.
	Name string `json:"name,omitempty"`
}

// SecurityPolicyTarget defines the target endpoints to apply SecurityPolicy.
type SecurityPolicyTarget struct {
	// VMSelector uses label selector to select VMs.
	VMSelector *metav1.LabelSelector `json:"vmSelector,omitempty"`
	// PodSelector uses label selector to select Pods.
	PodSelector *metav1.LabelSelector `json:"podSelector,omitempty"`
}

// SecurityPolicyPeer defines the source or destination of traffic.
type SecurityPolicyPeer struct {
	// VMSelector uses label selector to select VMs.
	VMSelector *metav1.LabelSelector `json:"vmSelector,omitempty"`
	// PodSelector uses label selector to select Pods.
	PodSelector *metav1.LabelSelector `json:"podSelector,omitempty"`
	// NamespaceSelector uses label selector to select Namespaces.
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
	// IPBlocks is a list of IP CIDRs.
	IPBlocks []IPBlock `json:"ipBlocks,omitempty"`
}

// IPBlock describes a particular CIDR that is allowed or denied to/from the workloads matched by an AppliedTo.
type IPBlock struct {
	// CIDR is a string representing the IP Block.
	// A valid example is "192.168.1.1/24".
	CIDR string `json:"cidr"`
}

// SecurityPolicyPort describes protocol and ports for traffic.
type SecurityPolicyPort struct {
	// Protocol(TCP, UDP) is the protocol to match traffic.
	// It is TCP by default.
	Protocol corev1.Protocol `json:"protocol,omitempty"`
	// Port is the name or port number.
	Port intstr.IntOrString `json:"port,omitempty"`
	// EndPort defines the end of port range.
	EndPort int `json:"endPort,omitempty"`
}

// SecurityPolicyStatus defines the observed state of SecurityPolicy.
type SecurityPolicyStatus struct {
	// Conditions describes current state of security policy.
	Conditions []Condition `json:"conditions"`
}

// +genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion

// SecurityPolicy is the Schema for the securitypolicies API.
type SecurityPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecurityPolicySpec   `json:"spec"`
	Status SecurityPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SecurityPolicyList contains a list of SecurityPolicy.
type SecurityPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecurityPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecurityPolicy{}, &SecurityPolicyList{})
}
