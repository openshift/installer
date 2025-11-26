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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/statusboard/v1"
)

// ServiceDependencyKind is the name of the type used to represent objects
// of type 'service_dependency'.
const ServiceDependencyKind = api_v1.ServiceDependencyKind

// ServiceDependencyLinkKind is the name of the type used to represent links
// to objects of type 'service_dependency'.
const ServiceDependencyLinkKind = api_v1.ServiceDependencyLinkKind

// ServiceDependencyNilKind is the name of the type used to nil references
// to objects of type 'service_dependency'.
const ServiceDependencyNilKind = api_v1.ServiceDependencyNilKind

// ServiceDependency represents the values of the 'service_dependency' type.
//
// Definition of a Status Board service dependency.
type ServiceDependency = api_v1.ServiceDependency

// ServiceDependencyListKind is the name of the type used to represent list of objects of
// type 'service_dependency'.
const ServiceDependencyListKind = api_v1.ServiceDependencyListKind

// ServiceDependencyListLinkKind is the name of the type used to represent links to list
// of objects of type 'service_dependency'.
const ServiceDependencyListLinkKind = api_v1.ServiceDependencyListLinkKind

// ServiceDependencyNilKind is the name of the type used to nil lists of objects of
// type 'service_dependency'.
const ServiceDependencyListNilKind = api_v1.ServiceDependencyListNilKind

type ServiceDependencyList = api_v1.ServiceDependencyList
