/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// RegistrySources represents the values of the 'registry_sources' type.
//
// RegistrySources contains configuration that determines how the container runtime should treat individual
// registries when accessing images for builds and pods. For instance, whether or not to allow insecure access.
// It does not contain configuration for the internal cluster registry.
type RegistrySources = api_v1.RegistrySources

// RegistrySourcesListKind is the name of the type used to represent list of objects of
// type 'registry_sources'.
const RegistrySourcesListKind = api_v1.RegistrySourcesListKind

// RegistrySourcesListLinkKind is the name of the type used to represent links to list
// of objects of type 'registry_sources'.
const RegistrySourcesListLinkKind = api_v1.RegistrySourcesListLinkKind

// RegistrySourcesNilKind is the name of the type used to nil lists of objects of
// type 'registry_sources'.
const RegistrySourcesListNilKind = api_v1.RegistrySourcesListNilKind

type RegistrySourcesList = api_v1.RegistrySourcesList
