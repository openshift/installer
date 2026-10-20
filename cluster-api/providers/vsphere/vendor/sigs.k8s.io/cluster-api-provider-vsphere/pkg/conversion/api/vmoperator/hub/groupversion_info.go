/*
Copyright 2026 The Kubernetes Authors.

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

package hub

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
)

// GroupName specifies the group name used to register the objects.
const GroupName = "vmoperator.vmware.com"

var (
	// GroupVersion is group version used to register these objects.
	GroupVersion = schema.GroupVersion{Group: GroupName, Version: "hub"}

	// schemeBuilder is used to add go types to the GroupVersionKind scheme.
	schemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = schemeBuilder.AddToScheme

	// converterBuilder is used to add go types to the Converter.
	converterBuilder = conversion.NewConverterBuilder(GroupVersion, addConvertibleTypes)

	// AddToConverter adds the types in this group-version to the given converter.
	AddToConverter = converterBuilder.AddToConverter

	objectTypes = []runtime.Object{}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion, objectTypes...)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}

func addConvertibleTypes(converter *conversion.Converter) error {
	return converter.AddHubTypes(GroupVersion, objectTypes...)
}
