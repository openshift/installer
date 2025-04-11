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

// +kubebuilder:validation:Enum:=managed;unmanaged
type ManagementPolicy string

const (
	// ManagementPolicyManaged specifies that the controller will reconcile the
	// state of the referenced OpenStack resource with the state of the ORC
	// object.
	ManagementPolicyManaged ManagementPolicy = "managed"

	// ManagementPolicyUnmanaged specifies that the controller will expect the
	// resource to either exist already or to be created externally. The
	// controller will not make any changes to the referenced OpenStack
	// resource.
	ManagementPolicyUnmanaged ManagementPolicy = "unmanaged"
)

// +kubebuilder:validation:Enum:=delete;detach
type OnDelete string

const (
	// OnDeleteDelete specifies that the OpenStack resource will be deleted
	// when the managed ORC object is deleted.
	OnDeleteDelete OnDelete = "delete"

	// OnDeleteDetach specifies that the OpenStack resource will not be
	// deleted when the managed ORC object is deleted.
	OnDeleteDetach OnDelete = "detach"
)

type ManagedOptions struct {
	// OnDelete specifies the behaviour of the controller when the ORC
	// object is deleted. Options are `delete` - delete the OpenStack resource;
	// `detach` - do not delete the OpenStack resource. If not specified, the
	// default is `delete`.
	// +kubebuilder:default:=delete
	// +optional
	OnDelete OnDelete `json:"onDelete,omitempty"`
}

// GetOnDelete returns the delete behaviour from ManagedOptions. If called on a
// nil receiver it safely returns the default.
func (o *ManagedOptions) GetOnDelete() OnDelete {
	if o == nil {
		return OnDeleteDelete
	}
	return o.OnDelete
}
