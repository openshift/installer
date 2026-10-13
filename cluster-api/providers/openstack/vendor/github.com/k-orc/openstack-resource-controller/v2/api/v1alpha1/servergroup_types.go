/*
Copyright 2024 The ORC Authors.

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

package v1alpha1

// +kubebuilder:validation:Enum:=affinity;anti-affinity;soft-affinity;soft-anti-affinity
type ServerGroupPolicy string

const (
	// ServerGroupPolicyAffinity is a server group policy that restricts instances belonging to the server group to the same host.
	ServerGroupPolicyAffinity ServerGroupPolicy = "affinity"
	// ServerGroupPolicyAntiAffinity is a server group policy that restricts instances belonging to the server group to separate hosts.
	ServerGroupPolicyAntiAffinity ServerGroupPolicy = "anti-affinity"
	// ServerGroupPolicySoftAffinity is a server group policy that attempts to restrict instances belonging to the server group to the same host.
	// Where it is not possible to schedule all instances on one host, they will be scheduled together on as few hosts as possible.
	ServerGroupPolicySoftAffinity ServerGroupPolicy = "soft-affinity"
	// ServerGroupPolicySoftAntiAffinity is a server group policy that attempts to restrict instances belonging to the server group to separate hosts.
	//  Where it is not possible to schedule all instances to separate hosts, they will be scheduled on as many separate hosts as possible.
	ServerGroupPolicySoftAntiAffinity ServerGroupPolicy = "soft-anti-affinity"
)

type ServerGroupRules struct {
	// maxServerPerHost specifies how many servers can reside on a single compute host.
	// It can be used only with the "anti-affinity" policy.
	// +optional
	MaxServerPerHost int32 `json:"maxServerPerHost,omitempty"`
}

// ServerGroupResourceSpec contains the desired state of a servergroup
// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="ServerGroupResourceSpec is immutable"
// +kubebuilder:validation:XValidation:rule="has(self.rules) && self.rules.maxServerPerHost > 0 ? self.policy == 'anti-affinity' : true",message="maxServerPerHost can only be used with the anti-affinity policy"
type ServerGroupResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// policy is the policy to use for the server group.
	// +required
	Policy ServerGroupPolicy `json:"policy,omitempty"`

	// rules is the rules to use for the server group.
	// +optional
	Rules *ServerGroupRules `json:"rules,omitempty"`
}

// ServerGroupFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type ServerGroupFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`
}

type ServerGroupRulesStatus struct {
	// maxServerPerHost specifies how many servers can reside on a single compute host.
	// It can be used only with the "anti-affinity" policy.
	// +optional
	MaxServerPerHost *int32 `json:"maxServerPerHost,omitempty"`
}

// ServerGroupResourceStatus represents the observed state of the resource.
type ServerGroupResourceStatus struct {
	// name is a Human-readable name for the servergroup. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// policy is the policy of the servergroup.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Policy string `json:"policy,omitempty"`

	// projectID is the project owner of the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// userID of the server group.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	UserID string `json:"userID,omitempty"`

	// rules is the rules of the server group.
	// +optional
	Rules *ServerGroupRulesStatus `json:"rules,omitempty"`
}
