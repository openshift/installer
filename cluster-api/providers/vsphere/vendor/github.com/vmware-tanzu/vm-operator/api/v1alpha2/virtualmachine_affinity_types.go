// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// VMAffinityTerm defines the VM affinity/anti-affinity term.
type VMAffinityTerm struct {
	// +optional

	// LabelSelector is a label query over a set of VMs.
	// When omitted, this term matches with no VMs.
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`

	// TopologyKey describes where this VM should be co-located (affinity) or not
	// co-located (anti-affinity).
	// Commonly used values include:
	// `kubernetes.io/hostname` -- The rule is executed in the context of a node/host.
	// `topology.kubernetes.io/zone` -- This rule is executed in the context of a zone.
	//
	// Please note, The following rules apply when specifying the topology key in the context of a zone/host.
	//
	// - When topology key is in the context of a zone, the only supported verbs are
	//   PreferredDuringSchedulingPreferredDuringExecution and RequiredDuringSchedulingPreferredDuringExecution.
	// - When topology key is in the context of a host, the only supported verbs are
	//   PreferredDuringSchedulingPreferredDuringExecution and RequiredDuringSchedulingPreferredDuringExecution
	//   for VM-VM node-level anti-affinity scheduling.
	TopologyKey string `json:"topologyKey"`
}

// VMAffinitySpec defines the affinity requirements for scheduling rules related
// to other VMs.
type VMAffinitySpec struct {
	// +optional
	// +listType=atomic

	// RequiredDuringSchedulingPreferredDuringExecution describes affinity
	// requirements that must be met or the VM will not be scheduled.
	//
	// When there are multiple elements, the lists of nodes corresponding to
	// each term are intersected, i.e. all terms must be satisfied.
	//
	// Note: Any update to this field will replace the entire list rather than
	// merging with the existing elements.
	RequiredDuringSchedulingPreferredDuringExecution []VMAffinityTerm `json:"requiredDuringSchedulingPreferredDuringExecution,omitempty"`

	// +optional
	// +listType=atomic

	// PreferredDuringSchedulingPreferredDuringExecution describes affinity
	// requirements that should be met, but the VM can still be scheduled if
	// the requirement cannot be satisfied. The scheduler will prefer to
	// schedule VMs that satisfy the affinity expressions specified by this
	// field, but it may choose to violate one or more of the expressions.
	//
	// When there are multiple elements, the lists of nodes corresponding to
	// each term are intersected, i.e. all terms must be satisfied.
	//
	// Note: Any update to this field will replace the entire list rather than
	// merging with the existing elements.
	PreferredDuringSchedulingPreferredDuringExecution []VMAffinityTerm `json:"preferredDuringSchedulingPreferredDuringExecution,omitempty"`
}

// VMAntiAffinitySpec defines the anti-affinity requirements for scheduling
// rules related to other VMs.
type VMAntiAffinitySpec struct {
	// +optional
	// +listType=atomic

	// RequiredDuringSchedulingPreferredDuringExecution describes anti-affinity
	// requirements that must be met or the VM will not be scheduled.
	//
	// When there are multiple elements, the lists of nodes corresponding to
	// each term are intersected, i.e. all terms must be satisfied.
	//
	// Note: Any update to this field will replace the entire list rather than
	// merging with the existing elements.
	RequiredDuringSchedulingPreferredDuringExecution []VMAffinityTerm `json:"requiredDuringSchedulingPreferredDuringExecution,omitempty"`

	// +optional
	// +listType=atomic

	// PreferredDuringSchedulingPreferredDuringExecution describes anti-affinity
	// requirements that should be met, but the VM can still be scheduled if
	// the requirement cannot be satisfied. The scheduler will prefer to
	// schedule VMs that satisfy the anti-affinity expressions specified by
	// this field, but it may choose to violate one or more of the expressions.
	// Additionally, it also describes the anti-affinity requirements that
	// should be met during run-time, but the VM can still be run if the
	// requirements cannot be satisfied.
	//
	// When there are multiple elements, the lists of nodes corresponding to
	// each term are intersected, i.e. all terms must be satisfied.
	//
	// Note: Any update to this field will replace the entire list rather than
	// merging with the existing elements.
	PreferredDuringSchedulingPreferredDuringExecution []VMAffinityTerm `json:"preferredDuringSchedulingPreferredDuringExecution,omitempty"`
}

// AffinitySpec defines the group of affinity scheduling rules.
type AffinitySpec struct {
	// +optional

	// VMAffinity describes affinity scheduling rules related to other VMs.
	VMAffinity *VMAffinitySpec `json:"vmAffinity,omitempty"`

	// +optional

	// VMAntiAffinity describes anti-affinity scheduling rules related to other
	// VMs.
	VMAntiAffinity *VMAntiAffinitySpec `json:"vmAntiAffinity,omitempty"`
}
