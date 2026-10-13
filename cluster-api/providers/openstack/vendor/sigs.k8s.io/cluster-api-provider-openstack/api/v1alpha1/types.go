/*
Copyright 2024 The Kubernetes Authors.

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

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

// ResolvedServerSpec contains resolved references to resources required by the server.
type ResolvedServerSpec struct {
	// ServerGroupID is the ID of the server group the server should be added to and is calculated based on ServerGroupFilter.
	// +optional
	ServerGroupID string `json:"serverGroupID,omitempty"`

	// ImageID is the ID of the image to use for the server and is calculated based on ImageFilter.
	// +optional
	ImageID string `json:"imageID,omitempty"`

	// FlavorID is the ID of the flavor to use.
	// +optional
	FlavorID string `json:"flavorID,omitempty"`

	// Ports is the fully resolved list of ports to create for the server.
	// +optional
	Ports []infrav1.ResolvedPortSpec `json:"ports,omitempty"`
}

// ServerResources contains references to OpenStack resources created for the server.
type ServerResources struct {
	// Ports is the status of the ports created for the server.
	// +optional
	Ports []infrav1.PortStatus `json:"ports,omitempty"`
}
