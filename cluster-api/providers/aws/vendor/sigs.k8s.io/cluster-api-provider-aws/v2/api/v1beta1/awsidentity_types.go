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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSClusterIdentitySpec defines the Spec struct for AWSClusterIdentity types.
type AWSClusterIdentitySpec struct {
	// AllowedNamespaces is used to identify which namespaces are allowed to use the identity from.
	// Namespaces can be selected either using an array of namespaces or with label selector.
	// An empty allowedNamespaces object indicates that AWSClusters can use this identity from any namespace.
	// If this object is nil, no namespaces will be allowed (default behaviour, if this field is not provided)
	// A namespace should be either in the NamespaceList or match with Selector to use the identity.
	//
	// +optional
	// +nullable
	AllowedNamespaces *AllowedNamespaces `json:"allowedNamespaces"`
}

// AllowedNamespaces is a selector of namespaces that AWSClusters can
// use this ClusterPrincipal from. This is a standard Kubernetes LabelSelector,
// a label query over a set of resources. The result of matchLabels and
// matchExpressions are ANDed.
type AllowedNamespaces struct {
	// An nil or empty list indicates that AWSClusters cannot use the identity from any namespace.
	//
	// +optional
	// +nullable
	NamespaceList []string `json:"list"`

	// An empty selector indicates that AWSClusters cannot use this
	// AWSClusterIdentity from any namespace.
	// +optional
	Selector metav1.LabelSelector `json:"selector"`
}

// AWSRoleSpec defines the specifications for all identities based around AWS roles.
type AWSRoleSpec struct {
	// The Amazon Resource Name (ARN) of the role to assume.
	RoleArn string `json:"roleARN"`
	// An identifier for the assumed role session
	SessionName string `json:"sessionName,omitempty"`
	// The duration, in seconds, of the role session before it is renewed.
	// +kubebuilder:validation:Minimum:=900
	// +kubebuilder:validation:Maximum:=43200
	DurationSeconds int32 `json:"durationSeconds,omitempty"`
	// An IAM policy as a JSON-encoded string that you want to use as an inline session policy.
	InlinePolicy string `json:"inlinePolicy,omitempty"`

	// The Amazon Resource Names (ARNs) of the IAM managed policies that you want
	// to use as managed session policies.
	// The policies must exist in the same account as the role.
	PolicyARNs []string `json:"policyARNs,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsclusterstaticidentities,scope=Cluster,categories=cluster-api,shortName=awssi

// AWSClusterStaticIdentity is the Schema for the awsclusterstaticidentities API
// It represents a reference to an AWS access key ID and secret access key, stored in a secret.
type AWSClusterStaticIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec for this AWSClusterStaticIdentity
	Spec AWSClusterStaticIdentitySpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AWSClusterStaticIdentityList contains a list of AWSClusterStaticIdentity.
type AWSClusterStaticIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSClusterStaticIdentity `json:"items"`
}

// AWSClusterStaticIdentitySpec defines the specifications for AWSClusterStaticIdentity.
type AWSClusterStaticIdentitySpec struct {
	AWSClusterIdentitySpec `json:",inline"`
	// Reference to a secret containing the credentials. The secret should
	// contain the following data keys:
	//  AccessKeyID: AKIAIOSFODNN7EXAMPLE
	//  SecretAccessKey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
	//  SessionToken: Optional
	SecretRef string `json:"secretRef"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsclusterroleidentities,scope=Cluster,categories=cluster-api,shortName=awsri

// AWSClusterRoleIdentity is the Schema for the awsclusterroleidentities API
// It is used to assume a role using the provided sourceRef.
type AWSClusterRoleIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec for this AWSClusterRoleIdentity.
	Spec AWSClusterRoleIdentitySpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AWSClusterRoleIdentityList contains a list of AWSClusterRoleIdentity.
type AWSClusterRoleIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSClusterRoleIdentity `json:"items"`
}

// AWSClusterRoleIdentitySpec defines the specifications for AWSClusterRoleIdentity.
type AWSClusterRoleIdentitySpec struct {
	AWSClusterIdentitySpec `json:",inline"`
	AWSRoleSpec            `json:",inline"`
	// A unique identifier that might be required when you assume a role in another account.
	// If the administrator of the account to which the role belongs provided you with an
	// external ID, then provide that value in the ExternalId parameter. This value can be
	// any string, such as a passphrase or account number. A cross-account role is usually
	// set up to trust everyone in an account. Therefore, the administrator of the trusting
	// account might send an external ID to the administrator of the trusted account. That
	// way, only someone with the ID can assume the role, rather than everyone in the
	// account. For more information about the external ID, see How to Use an External ID
	// When Granting Access to Your AWS Resources to a Third Party in the IAM User Guide.
	// +optional
	ExternalID string `json:"externalID,omitempty"`

	// SourceIdentityRef is a reference to another identity which will be chained to do
	// role assumption. All identity types are accepted.
	SourceIdentityRef *AWSIdentityReference `json:"sourceIdentityRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsclustercontrolleridentities,scope=Cluster,categories=cluster-api,shortName=awsci

// AWSClusterControllerIdentity is the Schema for the awsclustercontrolleridentities API
// It is used to grant access to use Cluster API Provider AWS Controller credentials.
type AWSClusterControllerIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec for this AWSClusterControllerIdentity.
	Spec AWSClusterControllerIdentitySpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// AWSClusterControllerIdentityList contains a list of AWSClusterControllerIdentity.
type AWSClusterControllerIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSClusterControllerIdentity `json:"items"`
}

// AWSClusterControllerIdentitySpec defines the specifications for AWSClusterControllerIdentity.
type AWSClusterControllerIdentitySpec struct {
	AWSClusterIdentitySpec `json:",inline"`
}

func init() {
	SchemeBuilder.Register(
		&AWSClusterStaticIdentity{},
		&AWSClusterStaticIdentityList{},
		&AWSClusterRoleIdentity{},
		&AWSClusterRoleIdentityList{},
		&AWSClusterControllerIdentity{},
		&AWSClusterControllerIdentityList{},
	)
}
