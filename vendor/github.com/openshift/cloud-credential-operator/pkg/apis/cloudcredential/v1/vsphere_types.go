/*
Copyright 2020 The OpenShift Authors.

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

// TODO: these types should eventually be broken out, along with the actuator,
// to a separate repo.

// VSphereProviderSpec contains the required information to create RBAC role
// bindings for VSphere.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VSphereProviderSpec struct {
	metav1.TypeMeta `json:",inline"`

	// Permissions contains a list of groups of privileges that are being requested.
	Permissions []VSpherePermission `json:"permissions"`
}

// VSpherePermission captures the details of the privileges being requested for the list of entities.
type VSpherePermission struct {
	// Privileges is the list of access being requested.
	Privileges []string `json:"privileges"`

	// TODO: when implementing mint-mode will need to figure out how to allow
	// a CredentialsRequest to indicate that the above list of privileges should
	// be bound to a specific scope(s) (eg Storage, Hosts/Clusters, Networking, Global, etc).
	// Entities is the list of entities for which the list of permissions should be granted
	// access to.
	// Entities []string `json:"entities"`

	// Also will need to allow specifying whether permissions should "Propagate to children".
	// Propagate bool `json:"propagate"`
}

// VSphereProviderStatus contains the status of the credentials request in VSphere.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VSphereProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// SecretLastResourceVersion is the resource version of the secret resource
	// that was last synced. Used to determine if the object has changed and
	// requires a sync.
	SecretLastResourceVersion string `json:"secretLastResourceVersion"`
}
