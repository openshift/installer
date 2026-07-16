/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// ConvertConditionsToV1Beta2 converts CAPI v1beta1 conditions to metav1.Condition slice.
func ConvertConditionsToV1Beta2(src clusterv1beta1.Conditions, observedGeneration int64) []metav1.Condition {
	if src == nil {
		return nil
	}

	dst := make([]metav1.Condition, len(src))
	for i, c := range src {
		dst[i] = metav1.Condition{
			Type:               string(c.Type),
			Status:             metav1.ConditionStatus(c.Status),
			ObservedGeneration: observedGeneration,
			LastTransitionTime: c.LastTransitionTime,
			Reason:             c.Reason,
			Message:            c.Message,
		}
	}
	return dst
}

// ConvertConditionsFromV1Beta2 converts metav1.Condition slice to CAPI v1beta1 conditions.
func ConvertConditionsFromV1Beta2(src []metav1.Condition) clusterv1beta1.Conditions {
	if src == nil {
		return nil
	}

	dst := make(clusterv1beta1.Conditions, len(src))
	for i, c := range src {
		dst[i] = clusterv1beta1.Condition{
			Type:               clusterv1beta1.ConditionType(c.Type),
			Status:             corev1.ConditionStatus(c.Status),
			Severity:           clusterv1beta1.ConditionSeverityNone, // Lost in conversion
			LastTransitionTime: c.LastTransitionTime,
			Reason:             c.Reason,
			Message:            c.Message,
		}
	}
	return dst
}

// IsReady checks if the Ready condition is True.
func IsReady(conditions []metav1.Condition) bool {
	for _, c := range conditions {
		if c.Type == ReadyConditionReason && c.Status == metav1.ConditionTrue {
			return true
		}
	}
	return false
}
