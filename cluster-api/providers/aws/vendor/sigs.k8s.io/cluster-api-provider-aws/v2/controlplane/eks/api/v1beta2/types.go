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

package v1beta2

import (
	"fmt"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
)

// ControlPlaneLoggingSpec defines what EKS control plane logs that should be enabled.
type ControlPlaneLoggingSpec struct {
	// APIServer indicates if the Kubernetes API Server log (kube-apiserver) shoulkd be enabled
	// +kubebuilder:default=false
	APIServer bool `json:"apiServer"`
	// Audit indicates if the Kubernetes API audit log should be enabled
	// +kubebuilder:default=false
	Audit bool `json:"audit"`
	// Authenticator indicates if the iam authenticator log should be enabled
	// +kubebuilder:default=false
	Authenticator bool `json:"authenticator"`
	// ControllerManager indicates if the controller manager (kube-controller-manager) log should be enabled
	// +kubebuilder:default=false
	ControllerManager bool `json:"controllerManager"`
	// Scheduler indicates if the Kubernetes scheduler (kube-scheduler) log should be enabled
	// +kubebuilder:default=false
	Scheduler bool `json:"scheduler"`
}

// IsLogEnabled returns true if the log is enabled.
func (s *ControlPlaneLoggingSpec) IsLogEnabled(logName string) bool {
	if s == nil {
		return false
	}

	switch ekstypes.LogType(logName) {
	case ekstypes.LogTypeApi:
		return s.APIServer
	case ekstypes.LogTypeAudit:
		return s.Audit
	case ekstypes.LogTypeAuthenticator:
		return s.Authenticator
	case ekstypes.LogTypeControllerManager:
		return s.ControllerManager
	case ekstypes.LogTypeScheduler:
		return s.Scheduler
	default:
		return false
	}
}

// EKSTokenMethod defines the method for obtaining a client token to use when connecting to EKS.
type EKSTokenMethod string

var (
	// EKSTokenMethodIAMAuthenticator indicates that IAM autenticator will be used to get a token.
	EKSTokenMethodIAMAuthenticator = EKSTokenMethod("iam-authenticator")

	// EKSTokenMethodAWSCli indicates that the AWS CLI will be used to get a token
	// Version 1.16.156 or greater is required of the AWS CLI.
	EKSTokenMethodAWSCli = EKSTokenMethod("aws-cli")
)

var (
	// DefaultEKSControlPlaneRole is the name of the default IAM role to use for the EKS control plane
	// if no other role is supplied in the spec and if iam role creation is not enabled. The default
	// can be created using clusterawsadm or created manually.
	DefaultEKSControlPlaneRole = fmt.Sprintf("eks-controlplane%s", iamv1.DefaultNameSuffix)
)

// IAMAuthenticatorConfig represents an aws-iam-authenticator configuration.
type IAMAuthenticatorConfig struct {
	// RoleMappings is a list of role mappings
	// +optional
	RoleMappings []RoleMapping `json:"mapRoles,omitempty"`
	// UserMappings is a list of user mappings
	// +optional
	UserMappings []UserMapping `json:"mapUsers,omitempty"`
}

// KubernetesMapping represents the kubernetes RBAC mapping.
type KubernetesMapping struct {
	// UserName is a kubernetes RBAC user subject
	UserName string `json:"username"`
	// Groups is a list of kubernetes RBAC groups
	Groups []string `json:"groups"`
}

// RoleMapping represents a mapping from a IAM role to Kubernetes users and groups.
type RoleMapping struct {
	// RoleARN is the AWS ARN for the role to map
	// +kubebuilder:validation:MinLength:=31
	RoleARN string `json:"rolearn"`
	// KubernetesMapping holds the RBAC details for the mapping
	KubernetesMapping `json:",inline"`
}

// UserMapping represents a mapping from an IAM user to Kubernetes users and groups.
type UserMapping struct {
	// UserARN is the AWS ARN for the user to map
	// +kubebuilder:validation:MinLength:=31
	UserARN string `json:"userarn"`
	// KubernetesMapping holds the RBAC details for the mapping
	KubernetesMapping `json:",inline"`
}

// Addon represents a EKS addon.
type Addon struct {
	// Name is the name of the addon
	// +kubebuilder:validation:MinLength:=2
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Version is the version of the addon to use
	Version string `json:"version"`
	// Configuration of the EKS addon
	// +optional
	Configuration string `json:"configuration,omitempty"`
	// ConflictResolution is used to declare what should happen if there
	// are parameter conflicts. Defaults to overwrite
	// +kubebuilder:default=overwrite
	// +kubebuilder:validation:Enum=overwrite;none;preserve
	ConflictResolution *AddonResolution `json:"conflictResolution,omitempty"`
	// ServiceAccountRoleArn is the ARN of an IAM role to bind to the addons service account
	// +optional
	ServiceAccountRoleArn *string `json:"serviceAccountRoleARN,omitempty"`
}

// AddonResolution defines the method for resolving parameter conflicts.
type AddonResolution string

var (
	// AddonResolutionOverwrite indicates that if there are parameter conflicts then
	// resolution will be accomplished via overwriting.
	AddonResolutionOverwrite = AddonResolution("overwrite")

	// AddonResolutionNone indicates that if there are parameter conflicts then
	// resolution will not be done and an error will be reported.
	AddonResolutionNone = AddonResolution("none")

	// AddonResolutionPreserve indicates that if there are parameter conflicts then
	// resolution will result in preserving the existing value
	AddonResolutionPreserve = AddonResolution("preserve")
)

// AddonStatus defines the status for an addon.
type AddonStatus string

var (
	// AddonStatusCreating is a status to indicate the addon is creating.
	AddonStatusCreating = "creating"

	// AddonStatusActive is a status to indicate the addon is active.
	AddonStatusActive = "active"

	// AddonStatusCreateFailed is a status to indicate the addon failed creation.
	AddonStatusCreateFailed = "create_failed"

	// AddonStatusUpdating is a status to indicate the addon is updating.
	AddonStatusUpdating = "updating"

	// AddonStatusDeleting is a status to indicate the addon is deleting.
	AddonStatusDeleting = "deleting"

	// AddonStatusDeleteFailed is a status to indicate the addon failed deletion.
	AddonStatusDeleteFailed = "delete_failed"

	// AddonStatusDegraded is a status to indicate the addon is in a degraded state.
	AddonStatusDegraded = "degraded"
)

// AddonState represents the state of an addon.
type AddonState struct {
	// Name is the name of the addon
	Name string `json:"name"`
	// Version is the version of the addon to use
	Version string `json:"version"`
	// ARN is the AWS ARN of the addon
	ARN string `json:"arn"`
	// ServiceAccountRoleArn is the ARN of the IAM role used for the service account
	ServiceAccountRoleArn *string `json:"serviceAccountRoleARN,omitempty"`
	// CreatedAt is the date and time the addon was created at
	CreatedAt metav1.Time `json:"createdAt,omitempty"`
	// ModifiedAt is the date and time the addon was last modified
	ModifiedAt metav1.Time `json:"modifiedAt,omitempty"`
	// Status is the status of the addon
	Status *string `json:"status,omitempty"`
	// Issues is a list of issue associated with the addon
	Issues []AddonIssue `json:"issues,omitempty"`
}

// AddonIssue represents an issue with an addon.
type AddonIssue struct {
	// Code is the issue code
	Code *string `json:"code,omitempty"`
	// Message is the textual description of the issue
	Message *string `json:"message,omitempty"`
	// ResourceIDs is a list of resource ids for the issue
	ResourceIDs []string `json:"resourceIds,omitempty"`
}

const (
	// SecurityGroupCluster is the security group for communication between EKS
	// control plane and managed node groups.
	SecurityGroupCluster = infrav1.SecurityGroupRole("cluster")
)

// OIDCIdentityProviderConfig represents the configuration for an OIDC identity provider.
type OIDCIdentityProviderConfig struct {
	// This is also known as audience. The ID for the client application that makes
	// authentication requests to the OpenID identity provider.
	// +kubebuilder:validation:Required
	ClientID string `json:"clientId,omitempty"`

	// The JWT claim that the provider uses to return your groups.
	// +optional
	GroupsClaim *string `json:"groupsClaim,omitempty"`

	// The prefix that is prepended to group claims to prevent clashes with existing
	// names (such as system: groups). For example, the valueoidc: will create group
	// names like oidc:engineering and oidc:infra.
	// +optional
	GroupsPrefix *string `json:"groupsPrefix,omitempty"`

	// The name of the OIDC provider configuration.
	//
	// IdentityProviderConfigName is a required field
	// +kubebuilder:validation:Required
	IdentityProviderConfigName string `json:"identityProviderConfigName,omitempty"`

	// The URL of the OpenID identity provider that allows the API server to discover
	// public signing keys for verifying tokens. The URL must begin with https://
	// and should correspond to the iss claim in the provider's OIDC ID tokens.
	// Per the OIDC standard, path components are allowed but query parameters are
	// not. Typically the URL consists of only a hostname, like https://server.example.org
	// or https://example.com. This URL should point to the level below .well-known/openid-configuration
	// and must be publicly accessible over the internet.
	//
	// +kubebuilder:validation:Required
	IssuerURL string `json:"issuerUrl,omitempty"`

	// The key value pairs that describe required claims in the identity token.
	// If set, each claim is verified to be present in the token with a matching
	// value. For the maximum number of claims that you can require, see Amazon
	// EKS service quotas (https://docs.aws.amazon.com/eks/latest/userguide/service-quotas.html)
	// in the Amazon EKS User Guide.
	// +optional
	RequiredClaims map[string]string `json:"requiredClaims,omitempty"`

	// The JSON Web Token (JWT) claim to use as the username. The default is sub,
	// which is expected to be a unique identifier of the end user. You can choose
	// other claims, such as email or name, depending on the OpenID identity provider.
	// Claims other than email are prefixed with the issuer URL to prevent naming
	// clashes with other plug-ins.
	// +optional
	UsernameClaim *string `json:"usernameClaim,omitempty"`

	// The prefix that is prepended to username claims to prevent clashes with existing
	// names. If you do not provide this field, and username is a value other than
	// email, the prefix defaults to issuerurl#. You can use the value - to disable
	// all prefixing.
	// +optional
	UsernamePrefix *string `json:"usernamePrefix,omitempty"`

	// tags to apply to oidc identity provider association
	// +optional
	Tags infrav1.Tags `json:"tags,omitempty"`
}
