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

package v1beta1

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// GCPManagedControlPlaneTemplateResourceSpec specifies an GCP managed control plane template resource.
type GCPManagedControlPlaneTemplateResourceSpec struct {
	GCPManagedControlPlaneClassSpec `json:",inline"`
}

// GCPManagedControlPlaneTemplateMachineTemplate is only used to fulfill the CAPI contract which expects a
// MachineTemplate field for any controlplane ref in a topology.
type GCPManagedControlPlaneTemplateMachineTemplate struct{}

// GCPManagedMachinePoolTemplateResourceSpec specifies an GCP managed control plane template resource.
type GCPManagedMachinePoolTemplateResourceSpec struct {
	GCPManagedMachinePoolClassSpec `json:",inline"`
}

// GCPManagedClusterTemplateResourceSpec specifies an GCP managed cluster template resource.
type GCPManagedClusterTemplateResourceSpec struct {
	// Project is the name of the project to deploy the cluster to.
	Project string `json:"project"`

	// The GCP Region the cluster lives in.
	Region string `json:"region"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// NetworkSpec encapsulates all things related to the GCP network.
	// +optional
	Network infrav1.NetworkSpec `json:"network"`

	// AdditionalLabels is an optional set of tags to add to GCP resources managed by the GCP provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalLabels infrav1.Labels `json:"additionalLabels,omitempty"`

	// ResourceManagerTags is an optional set of tags to apply to GCP resources managed
	// by the GCP provider. GCP supports a maximum of 50 tags per resource.
	// +maxItems=50
	// +optional
	ResourceManagerTags infrav1.ResourceManagerTags `json:"resourceManagerTags,omitempty"`

	// CredentialsRef is a reference to a Secret that contains the credentials to use for provisioning this cluster. If not
	// supplied then the credentials of the controller will be used.
	// +optional
	CredentialsRef *infrav1.ObjectReference `json:"credentialsRef,omitempty"`

	// LoadBalancerSpec contains configuration for one or more LoadBalancers.
	// +optional
	LoadBalancer infrav1.LoadBalancerSpec `json:"loadBalancer,omitempty"`

	// ServiceEndpoints contains the custom GCP Service Endpoint urls for each applicable service.
	// For instance, the user can specify a new endpoint for the compute service.
	// +optional
	ServiceEndpoints *infrav1.ServiceEndpoints `json:"serviceEndpoints,omitempty"`
}
