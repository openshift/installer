/*
Copyright 2022 The Kubernetes Authors.

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

import "cloud.google.com/go/container/apiv1/containerpb"

// TaintEffect is the effect for a Kubernetes taint.
type TaintEffect string

// Taint represents a Kubernetes taint.
type Taint struct {
	// Effect specifies the effect for the taint.
	// +kubebuilder:validation:Enum=NoSchedule;NoExecute;PreferNoSchedule
	Effect TaintEffect `json:"effect"`
	// Key is the key of the taint
	Key string `json:"key"`
	// Value is the value of the taint
	Value string `json:"value"`
}

// Taints is an array of Taints.
type Taints []Taint

func convertToSdkTaintEffect(effect TaintEffect) containerpb.NodeTaint_Effect {
	switch effect {
	case "NoSchedule":
		return containerpb.NodeTaint_NO_SCHEDULE
	case "NoExecute":
		return containerpb.NodeTaint_NO_EXECUTE
	case "PreferNoSchedule":
		return containerpb.NodeTaint_PREFER_NO_SCHEDULE
	default:
		return containerpb.NodeTaint_EFFECT_UNSPECIFIED
	}
}

// ConvertToSdkTaint converts taints to format that is used by GCP SDK.
func ConvertToSdkTaint(taints Taints) []*containerpb.NodeTaint {
	if taints == nil {
		return nil
	}
	res := []*containerpb.NodeTaint{}
	for _, taint := range taints {
		res = append(res, &containerpb.NodeTaint{
			Key:    taint.Key,
			Value:  taint.Value,
			Effect: convertToSdkTaintEffect(taint.Effect),
		})
	}
	return res
}
