/*
Copyright 2021 The OpenShift Authors.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: these types should eventually be broken out, along with the actuator, to a separate repo.

// IBMCloudProviderSpec is the specification of the credentials request in IBM Cloud.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IBMCloudProviderSpec struct {
	metav1.TypeMeta `json:",inline"`

	// Policies are a list of access policies to create for the generated credentials
	Policies []AccessPolicy `json:"policies"`
}

// AccessPolicy is a definition of an IAM access policy
type AccessPolicy struct {
	// Attributes identify the resources to which this policy applies
	Attributes []ResourceAttribute `json:"attributes"`

	// Roles are the IAM roles assigned to this policy
	Roles []string `json:"roles"`
}

// ResourceAttribute is an attribute associated with a resource.
type ResourceAttribute struct {
	// Name is the name of an attribute.
	Name string `json:"name"`

	// Value is the value of an attribute.
	Value string `json:"value"`

	// Operator is the operator of an attribute.
	Operator string `json:"operator,omitempty"`
}

// IBMCloudProviderStatus contains the status of the IBM Cloud credentials request.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IBMCloudProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
}
