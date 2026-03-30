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

// Package v1 contains API Schema definitions for the cloudcredential v1 API group
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential
// +k8s:defaulter-gen=TypeMeta
// +groupName=cloudcredential.openshift.io
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	GroupName     = "cloudcredential.openshift.io"
	GroupVersion  = schema.GroupVersion{Group: GroupName, Version: "v1"}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// Install is a function which adds this version to a scheme
	Install = SchemeBuilder.AddToScheme

	// SchemeGroupVersion generated code relies on this name
	// DEPRECATED
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}
	// AddToScheme exists solely to keep the old generators creating valid code
	// DEPRECATED
	AddToScheme = SchemeBuilder.AddToScheme
)

// Resource generated code relies on this being here, but it logically belongs to the group
// DEPRECATED
func Resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: GroupName, Resource: resource}
}

func addKnownTypes(scheme *runtime.Scheme) error {
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&CredentialsRequest{}, &CredentialsRequestList{},
		&AWSProviderStatus{}, &AWSProviderSpec{},
		&AzureProviderStatus{}, &AzureProviderSpec{},
		&GCPProviderStatus{}, &GCPProviderSpec{},
		&IBMCloudProviderStatus{}, &IBMCloudProviderSpec{},
		&IBMCloudPowerVSProviderStatus{}, &IBMCloudPowerVSProviderSpec{},
		&NutanixProviderStatus{}, &NutanixProviderSpec{},
		&VSphereProviderStatus{}, &VSphereProviderSpec{},
		&KubevirtProviderStatus{}, &KubevirtProviderSpec{},
	)

	return nil
}
