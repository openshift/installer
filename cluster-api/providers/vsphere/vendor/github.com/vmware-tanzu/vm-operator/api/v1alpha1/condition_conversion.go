// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha1_Condition_To_v1_Condition(in *Condition, out *metav1.Condition, s apiconversion.Scope) error {
	out.Type = string(in.Type)
	out.Status = metav1.ConditionStatus(in.Status)
	out.LastTransitionTime = in.LastTransitionTime
	out.Reason = in.Reason
	out.Message = in.Message

	// The metav1.Condition requires the reason to be non-empty, when it was not in our prior v1a1 Condition.
	// We don't have any great options as to what we can fill this in as.
	if out.Reason == "" {
		out.Reason = string(out.Status)
	}

	// TODO: out.ObservedGeneration =

	return nil
}

func Convert_v1_Condition_To_v1alpha1_Condition(in *metav1.Condition, out *Condition, s apiconversion.Scope) error {
	out.Type = ConditionType(in.Type)
	out.Status = corev1.ConditionStatus(in.Status)
	out.LastTransitionTime = in.LastTransitionTime
	out.Reason = in.Reason
	out.Message = in.Message

	// TODO: out.Severity = ConditionSeverity.

	return nil
}
