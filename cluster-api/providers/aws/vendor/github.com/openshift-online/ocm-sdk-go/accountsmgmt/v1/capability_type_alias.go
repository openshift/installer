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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1"
)

// Capability represents the values of the 'capability' type.
//
// Capability model that represents internal labels with a key that matches a set list defined in AMS (defined in pkg/api/capability_types.go).
type Capability = api_v1.Capability

// CapabilityListKind is the name of the type used to represent list of objects of
// type 'capability'.
const CapabilityListKind = api_v1.CapabilityListKind

// CapabilityListLinkKind is the name of the type used to represent links to list
// of objects of type 'capability'.
const CapabilityListLinkKind = api_v1.CapabilityListLinkKind

// CapabilityNilKind is the name of the type used to nil lists of objects of
// type 'capability'.
const CapabilityListNilKind = api_v1.CapabilityListNilKind

type CapabilityList = api_v1.CapabilityList
