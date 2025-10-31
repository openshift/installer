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

// StatusKind is the name of the type used to represent objects
// of type 'status'.
const StatusKind = api_v1.StatusKind

// StatusLinkKind is the name of the type used to represent links
// to objects of type 'status'.
const StatusLinkKind = api_v1.StatusLinkKind

// StatusNilKind is the name of the type used to nil references
// to objects of type 'status'.
const StatusNilKind = api_v1.StatusNilKind

// Status represents the values of the 'status' type.
//
// Definition of a Status Board status.
type Status = api_v1.Status

// StatusListKind is the name of the type used to represent list of objects of
// type 'status'.
const StatusListKind = api_v1.StatusListKind

// StatusListLinkKind is the name of the type used to represent links to list
// of objects of type 'status'.
const StatusListLinkKind = api_v1.StatusListLinkKind

// StatusNilKind is the name of the type used to nil lists of objects of
// type 'status'.
const StatusListNilKind = api_v1.StatusListNilKind

type StatusList = api_v1.StatusList
