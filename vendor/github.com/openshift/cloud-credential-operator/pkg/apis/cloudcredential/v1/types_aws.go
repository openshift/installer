/*
Copyright 2018 The OpenShift Authors.

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

// AWSProviderSpec contains the required information to create a user policy in AWS.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AWSProviderSpec struct {
	metav1.TypeMeta `json:",inline"`
	// StatementEntries contains a list of policy statements that should be associated with this credentials access key.
	StatementEntries []StatementEntry `json:"statementEntries"`
	// stsIAMRoleARN is the Amazon Resource Name (ARN) of an IAM Role which was created manually for the associated
	// CredentialsRequest.
	// The presence of an stsIAMRoleARN within the AWSProviderSpec initiates creation of a secret containing IAM
	// Role details necessary for assuming the IAM Role via Amazon's Secure Token Service.
	// +optional
	STSIAMRoleARN string `json:"stsIAMRoleARN,omitempty"`
}

// StatementEntry models an AWS policy statement entry.
type StatementEntry struct {
	// Effect indicates if this policy statement is to Allow or Deny.
	Effect string `json:"effect"`
	// Action describes the particular AWS service actions that should be allowed or denied. (i.e. ec2:StartInstances, iam:ChangePassword)
	Action []string `json:"action"`
	// Resource specifies the object(s) this statement should apply to. (or "*" for all)
	Resource string `json:"resource"`
	// PolicyCondition specifies under which condition StatementEntry will apply
	PolicyCondition IAMPolicyCondition `json:"policyCondition,omitempty"`
}

// AWSProviderStatus containes the status of the credentials request in AWS.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AWSProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
	// User is the name of the User created in AWS for these credentials.
	User string `json:"user"`
	// Policy is the name of the policy attached to the user in AWS.
	Policy string `json:"policy"`
}

// IAMPolicyCondition - map of condition types, with associated key - value mapping
// +k8s:deepcopy-gen=false
type IAMPolicyCondition map[string]IAMPolicyConditionKeyValue

// IAMPolicyConditionKeyValue - mapping of values for the chosen type
// +k8s:deepcopy-gen=false
type IAMPolicyConditionKeyValue map[string]interface{}
