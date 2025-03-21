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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	iamv1 "sigs.k8s.io/cluster-api-provider-aws/v2/iam/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

var (
	// DefaultEKSFargateRole is the name of the default IAM role to use for fargate
	// profiles if no other role is supplied in the spec and if iam role creation
	// is not enabled. The default can be created using clusterawsadm or created manually.
	DefaultEKSFargateRole = fmt.Sprintf("eks-fargate%s", iamv1.DefaultNameSuffix)
)

// FargateProfileSpec defines the desired state of FargateProfile.
type FargateProfileSpec struct {
	// ClusterName is the name of the Cluster this object belongs to.
	// +kubebuilder:validation:MinLength=1
	ClusterName string `json:"clusterName"`

	// ProfileName specifies the profile name.
	ProfileName string `json:"profileName,omitempty"`

	// SubnetIDs specifies which subnets are used for the
	// auto scaling group of this nodegroup.
	// +optional
	SubnetIDs []string `json:"subnetIDs,omitempty"`

	// AdditionalTags is an optional set of tags to add to AWS resources managed by the AWS provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// RoleName specifies the name of IAM role for this fargate pool
	// If the role is pre-existing we will treat it as unmanaged
	// and not delete it on deletion. If the EKSEnableIAM feature
	// flag is true and no name is supplied then a role is created.
	// +optional
	RoleName string `json:"roleName,omitempty"`

	// Selectors specify fargate pod selectors.
	Selectors []FargateSelector `json:"selectors,omitempty"`
}

// FargateSelector specifies a selector for pods that should run on this fargate pool.
type FargateSelector struct {
	// Labels specifies which pod labels this selector should match.
	Labels map[string]string `json:"labels,omitempty"`

	// Namespace specifies which namespace this selector should match.
	Namespace string `json:"namespace,omitempty"`
}

// FargateProfileStatus defines the observed state of FargateProfile.
type FargateProfileStatus struct {
	// Ready denotes that the FargateProfile is available.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the FargateProfile and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the FargateProfile's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of
	// FargateProfiles can be added as events to the FargateProfile object
	// and/or logged in the controller's output.
	// +optional
	FailureReason *string `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the FargateProfile and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the FargateProfile's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of
	// FargateProfiles can be added as events to the FargateProfile
	// object and/or logged in the controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// Conditions defines current state of the Fargate profile.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=awsfargateprofiles,scope=Namespaced,categories=cluster-api,shortName=awsfp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="AWSFargateProfile ready status"
// +kubebuilder:printcolumn:name="ProfileName",type="string",JSONPath=".spec.profileName",description="EKS Fargate profile name"
// +kubebuilder:printcolumn:name="FailureReason",type="string",JSONPath=".status.failureReason",description="Failure reason"

// AWSFargateProfile is the Schema for the awsfargateprofiles API.
type AWSFargateProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FargateProfileSpec   `json:"spec,omitempty"`
	Status FargateProfileStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the AWSFargateProfile resource.
func (r *AWSFargateProfile) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the AWSFargateProfile to the predescribed clusterv1.Conditions.
func (r *AWSFargateProfile) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// AWSFargateProfileList contains a list of FargateProfiles.
type AWSFargateProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSFargateProfile `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSFargateProfile{}, &AWSFargateProfileList{})
}
