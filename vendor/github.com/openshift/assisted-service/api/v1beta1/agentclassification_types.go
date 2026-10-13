/*
Copyright 2022.

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

import (
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	QueryErrorsCondition conditionsv1.ConditionType = "QueryErrors"
	QueryNoErrorsReason  string                     = "NoQueryErrors"
	QueryHasErrorsReason string                     = "HasQueryErrors"
)

// AgentClassificationSpec defines the desired state of AgentClassification
type AgentClassificationSpec struct {
	// LabelKey specifies the label key to apply to matched Agents
	//
	// +immutable
	LabelKey string `json:"labelKey"`

	// LabelValue specifies the label value to apply to matched Agents
	//
	// +immutable
	LabelValue string `json:"labelValue"`

	// Query is in gojq format (https://github.com/itchyny/gojq#difference-to-jq)
	// and will be invoked on each Agent's inventory. The query should return a
	// boolean. The operator will apply the label to any Agent for which "true"
	// is returned.
	Query string `json:"query"`
}

// AgentClassificationStatus defines the observed state of AgentClassification
type AgentClassificationStatus struct {
	// MatchedCount shows how many Agents currently match the classification
	MatchedCount int `json:"matchedCount,omitempty"`

	// ErrorCount shows how many Agents encountered errors when matching the classification
	ErrorCount int `json:"errorCount,omitempty"`

	Conditions []conditionsv1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AgentClassification is the Schema for the AgentClassifications API
type AgentClassification struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentClassificationSpec   `json:"spec,omitempty"`
	Status AgentClassificationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AgentClassificationList contains a list of AgentClassification
type AgentClassificationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AgentClassification `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AgentClassification{}, &AgentClassificationList{})
}
