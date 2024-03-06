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

// Package conditions provides helper functions for working with conditions.
package conditions

import (
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// ErrorConditionAfterInit returns severity error, if the control plane is initialized; otherwise, returns severity warning.
// Failures after control plane is initialized is likely to be non-transient,
// hence conditions severities should be set to Error.
func ErrorConditionAfterInit(getter conditions.Getter) clusterv1.ConditionSeverity {
	if conditions.IsTrue(getter, clusterv1.ControlPlaneInitializedCondition) {
		return clusterv1.ConditionSeverityError
	}
	return clusterv1.ConditionSeverityWarning
}
