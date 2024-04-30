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

// AlibabaCloudProviderSpec contains the required information to create a user policy in AlibabaCloud.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AlibabaCloudProviderSpec struct {
	metav1.TypeMeta `json:",inline"`
	// StatementEntries contains a list of policy statements that should be associated with this credentials access key.
	StatementEntries []AlibabaStatementEntry `json:"statementEntries"`
}

// StatementEntry models an AlibabaCloud policy statement entry.
type AlibabaStatementEntry struct {
	// Effect indicates if this policy statement is to Allow or Deny.
	Effect string `json:"effect"`
	// Action describes the particular AlibabaCloud service actions that should be allowed or denied.
	Action []string `json:"action"`
	// Resource specifies the object(s) this statement should apply to. (or "*" for all)
	Resource string `json:"resource"`
}

// AlibabaCloudProviderStatus containes the status of the credentials request in AlibabaCloud.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AlibabaCloudProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
}
