/*
Copyright 2025 The Kubernetes Authors.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// OidcProviderType set to Managed or UnManaged
type OidcProviderType string

const (
	// Managed OIDC Provider type
	Managed OidcProviderType = "Managed"

	// Unmanaged OIDC Provider type
	Unmanaged OidcProviderType = "Unmanaged"
)

// Operator Role const
const (
	// IngressOperatorARNSuffix is the suffix for the ingress operator role.
	IngressOperatorARNSuffix = "-openshift-ingress-operator-cloud-credentials"

	// ImageRegistryARNSuffix is the suffix for the image registry operator role.
	ImageRegistryARNSuffix = "-openshift-image-registry-installer-cloud-credentials"

	// StorageARNSuffix is the suffix for the storage operator role.
	StorageARNSuffix = "-openshift-cluster-csi-drivers-ebs-cloud-credentials"

	// NetworkARNSuffix is the suffix for the network operator role.
	NetworkARNSuffix = "-openshift-cloud-network-config-controller-cloud-credentials"

	// KubeCloudControllerARNSuffix is the suffix for the kube cloud controller role.
	KubeCloudControllerARNSuffix = "-kube-system-kube-controller-manager"

	// NodePoolManagementARNSuffix is the suffix for the node pool management role.
	NodePoolManagementARNSuffix = "-kube-system-capa-controller-manager"

	// ControlPlaneOperatorARNSuffix is the suffix for the control plane operator role.
	ControlPlaneOperatorARNSuffix = "-kube-system-control-plane-operator"

	// KMSProviderARNSuffix is the suffix for the kms provider role.
	KMSProviderARNSuffix = "-kube-system-kms-provider"
)

// Account Role const
const (
	// HCPROSAInstallerRole is the suffix for installer account role
	HCPROSAInstallerRole = "-HCP-ROSA-Installer-Role"

	// HCPROSASupportRole is the suffix for support account role
	HCPROSASupportRole = "-HCP-ROSA-Support-Role"

	// HCPROSAWorkerRole is the suffix for worker account role
	HCPROSAWorkerRole = "-HCP-ROSA-Worker-Role"
)

const (
	// RosaRoleConfigReadyCondition condition reports on the successful reconciliation of RosaRoleConfig.
	RosaRoleConfigReadyCondition = "RosaRoleConfigReady"

	// RosaRoleConfigDeletionFailedReason used to report failures while deleting RosaRoleConfig.
	RosaRoleConfigDeletionFailedReason = "DeletionFailed"

	// RosaRoleConfigReconciliationFailedReason used to report reconciliation failures.
	RosaRoleConfigReconciliationFailedReason = "ReconciliationFailed"

	// RosaRoleConfigDeletionStarted used to indicate that the deletion of RosaRoleConfig has started.
	RosaRoleConfigDeletionStarted = "DeletionStarted"

	// RosaRoleConfigCreatedReason used to indicate that the RosaRoleConfig has been created.
	RosaRoleConfigCreatedReason = "Created"
)

// ROSARoleConfigSpec defines the desired state of ROSARoleConfig
type ROSARoleConfigSpec struct {
	// AccountRoleConfig defines account-wide IAM roles before creating your ROSA cluster.
	AccountRoleConfig AccountRoleConfig `json:"accountRoleConfig"`

	// OperatorRoleConfig defines cluster-specific operator IAM roles based on your cluster configuration.
	OperatorRoleConfig OperatorRoleConfig `json:"operatorRoleConfig"`

	// IdentityRef is a reference to an identity to be used when reconciling the ROSA Role Config.
	// If no identity is specified, the default identity for this controller will be used.
	// +optional
	IdentityRef *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`

	// CredentialsSecretRef references a secret with necessary credentials to connect to the OCM API.
	// +optional
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`

	// OIDC provider type values are Managed or UnManaged. When set to Unmanged OperatorRoleConfig OIDCID field must be provided.
	// +kubebuilder:validation:Enum=Managed;Unmanaged
	// +kubebuilder:default=Managed
	OidcProviderType OidcProviderType `json:"oidcProviderType"`
}

// AccountRoleConfig defines account IAM roles before creating your ROSA cluster.
type AccountRoleConfig struct {
	// User-defined prefix for all generated AWS account role
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength:=4
	// +kubebuilder:validation:Pattern:=`^[a-z]([-a-z0-9]*[a-z0-9])?$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="prefix is immutable"
	Prefix string `json:"prefix"`

	// The ARN of the policy that is used to set the permissions boundary for the account roles.
	// +optional
	PermissionsBoundaryARN string `json:"permissionsBoundaryARN,omitempty"`

	// The arn path for the account/operator roles as well as their policies.
	// +optional
	Path string `json:"path,omitempty"`

	// Version of OpenShift that will be used to the roles tag in formate of x.y.z example; "4.19.0"
	// Setting the role OpenShift version tag does not affect the associated ROSAControlplane version.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="version is immutable"
	Version string `json:"version"`

	// SharedVPCConfig is used to set up shared VPC.
	// +optional
	SharedVPCConfig SharedVPCConfig `json:"sharedVPCConfig,omitempty"`
}

// OperatorRoleConfig defines cluster-specific operator IAM roles based on your cluster configuration.
type OperatorRoleConfig struct {
	//  User-defined prefix for generated AWS operator roles.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength:=4
	// +kubebuilder:validation:Pattern:=`^[a-z]([-a-z0-9]*[a-z0-9])?$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="prefix is immutable"
	Prefix string `json:"prefix"`

	// The ARN of the policy that is used to set the permissions boundary for the operator roles.
	// +optional
	PermissionsBoundaryARN string `json:"permissionsBoundaryARN,omitempty"`

	// SharedVPCConfig is used to set up shared VPC.
	// +optional
	SharedVPCConfig SharedVPCConfig `json:"sharedVPCConfig,omitempty"`

	// OIDCID is the ID of the OIDC config that will be used to create the operator roles.
	// Cannot be set when OidcProviderType set to Managed
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="oidcID is immutable"
	OIDCID string `json:"oidcID,omitempty"`
}

// SharedVPCConfig is used to set up shared VPC.
type SharedVPCConfig struct {
	// Role ARN associated with the private hosted zone used for Hosted Control Plane cluster shared VPC, this role contains policies to be used with Route 53
	RouteRoleARN string `json:"routeRoleARN,omitempty"`

	// Role ARN associated with the shared VPC used for Hosted Control Plane clusters, this role contains policies to be used with the VPC endpoint
	VPCEndpointRoleARN string `json:"vpcEndpointRoleArn,omitempty"`
}

// ROSARoleConfigStatus defines the observed state of ROSARoleConfig
type ROSARoleConfigStatus struct {
	// ID of created OIDC config
	OIDCID string `json:"oidcID,omitempty"`

	// Create OIDC provider for operators to authenticate against in an STS cluster.
	OIDCProviderARN string `json:"oidcProviderARN,omitempty"`

	// Created Account roles that can be used to
	AccountRolesRef AccountRolesRef `json:"accountRolesRef,omitempty"`

	// AWS IAM roles used to perform credential requests by the openshift operators.
	OperatorRolesRef rosacontrolplanev1.AWSRolesRef `json:"operatorRolesRef,omitempty"`

	// Conditions specifies the ROSARoleConfig conditions
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// AccountRolesRef defscribes ARNs used as Account roles.
type AccountRolesRef struct {
	// InstallerRoleARN is an AWS IAM role that OpenShift Cluster Manager will assume to create the cluster..
	InstallerRoleARN string `json:"installerRoleARN,omitempty"`

	// SupportRoleARN is an AWS IAM role used by Red Hat SREs to enable
	// access to the cluster account in order to provide support.
	SupportRoleARN string `json:"supportRoleARN,omitempty"`

	// WorkerRoleARN is an AWS IAM role that will be attached to worker instances.
	WorkerRoleARN string `json:"workerRoleARN,omitempty"`
}

// ROSARoleConfig is the Schema for the rosaroleconfigs API
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosaroleconfigs,scope=Namespaced,categories=cluster-api,shortName=rosarole
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type ROSARoleConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ROSARoleConfigSpec   `json:"spec,omitempty"`
	Status ROSARoleConfigStatus `json:"status,omitempty"`
}

// ROSARoleConfigList contains a list of ROSARoleConfig
// +kubebuilder:object:root=true
type ROSARoleConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSARoleConfig `json:"items"`
}

// SetConditions sets the conditions of the ROSARoleConfig.
func (r *ROSARoleConfig) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// GetConditions returns the observations of the operational state of the RosaNetwork resource.
func (r *ROSARoleConfig) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// IsSharedVPC checks if the shared VPC config is set.
func (s SharedVPCConfig) IsSharedVPC() bool {
	return s.VPCEndpointRoleARN != "" || s.RouteRoleARN != ""
}

func init() {
	SchemeBuilder.Register(&ROSARoleConfig{}, &ROSARoleConfigList{})
}
