/*
Copyright 2024 The Kubernetes Authors.

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

package ssa

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfigv1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// ConditionsEqual compares all properties of a ConditionApplyConfiguration with
// the corresponding properties of a Condition, except ObservedGeneration and
// LastTransitionTime.
func ConditionsEqual(previous *metav1.Condition, applyConfig *applyconfigv1.ConditionApplyConfiguration) bool {
	return (applyConfig.Type != nil && previous.Type == *applyConfig.Type) &&
		(applyConfig.Status != nil && previous.Status == *applyConfig.Status) &&
		(applyConfig.Reason != nil && previous.Reason == *applyConfig.Reason) &&
		(applyConfig.Message != nil && previous.Message == *applyConfig.Message)
}
