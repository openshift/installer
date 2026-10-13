/*
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

package iamserviceaccount

// ServiceAccountRole represents an IAM role associated with a Kubernetes service account
type ServiceAccountRole struct {
	RoleName            string   `json:"roleName"`
	RoleARN             string   `json:"roleArn"`
	ClusterName         string   `json:"clusterName"`
	ServiceAccountName  string   `json:"serviceAccountName"`
	Namespace           string   `json:"namespace"`
	PolicyARNs          []string `json:"policyArns"`
	InlinePolicy        string   `json:"inlinePolicy,omitempty"`
	PermissionsBoundary string   `json:"permissionsBoundary,omitempty"`
	TrustPolicy         string   `json:"trustPolicy"`
}

// CreateOptions contains options for creating a service account IAM role
type CreateOptions struct {
	ClusterName         string
	ServiceAccountName  string
	Namespace           string
	RoleName            string
	PolicyARNs          []string
	InlinePolicy        string
	PermissionsBoundary string
	Path                string
	Tags                map[string]string
}

// ListOptions contains filtering options for listing service account roles
type ListOptions struct {
	ClusterName string
	Namespace   string
}

// ValidateOptions contains validation configuration
type ValidateOptions struct {
	VerifyServiceAccountExists bool
	CheckExistingRole          bool
}

// ServiceAccountIdentifier represents a service account with its namespace
type ServiceAccountIdentifier struct {
	Name      string
	Namespace string
}
