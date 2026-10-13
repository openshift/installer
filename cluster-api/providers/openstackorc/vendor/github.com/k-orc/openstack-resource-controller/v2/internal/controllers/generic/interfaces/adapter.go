/*
Copyright 2025 The ORC Authors.

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

package interfaces

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
)

type APIObjectAdapter[orcObjectPT any, resourceSpecT any, filterT any] interface {
	// NOTE: the objects we generate which implement this interface all also
	// implement client.Object. Despite this, they fail at runtime when passed
	// to the controller-runtime client. Therefore we deliberately don't declare
	// this interface here to prevent accidental usage.
	metav1.Object

	GetObject() orcObjectPT

	GetManagementPolicy() orcv1alpha1.ManagementPolicy
	GetManagedOptions() *orcv1alpha1.ManagedOptions

	GetStatusID() *string
	GetResourceSpec() *resourceSpecT
	GetImportID() *string
	GetImportFilter() *filterT
}
