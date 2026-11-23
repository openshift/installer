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

package v1alpha2

import (
	vmoprv1alpha2 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
)

var (
	// converterBuilder is used to add go types to the Converter.
	converterBuilder = conversion.NewConverterBuilder(vmoprv1alpha2.GroupVersion)

	// AddToConverter adds the types in this group-version to the given converter.
	AddToConverter = converterBuilder.AddToConverter
)
