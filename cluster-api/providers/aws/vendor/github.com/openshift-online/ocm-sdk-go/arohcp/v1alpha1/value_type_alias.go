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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	api_v1alpha1 "github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1"
)

// Value represents the values of the 'value' type.
//
// Numeric value and the unit used to measure it.
//
// Units are not mandatory, and they're not specified for some resources. For
// resources that use bytes, the accepted units are:
//
// - 1 B = 1 byte
// - 1 KB = 10^3 bytes
// - 1 MB = 10^6 bytes
// - 1 GB = 10^9 bytes
// - 1 TB = 10^12 bytes
// - 1 PB = 10^15 bytes
//
// - 1 B = 1 byte
// - 1 KiB = 2^10 bytes
// - 1 MiB = 2^20 bytes
// - 1 GiB = 2^30 bytes
// - 1 TiB = 2^40 bytes
// - 1 PiB = 2^50 bytes
type Value = api_v1alpha1.Value

// ValueListKind is the name of the type used to represent list of objects of
// type 'value'.
const ValueListKind = api_v1alpha1.ValueListKind

// ValueListLinkKind is the name of the type used to represent links to list
// of objects of type 'value'.
const ValueListLinkKind = api_v1alpha1.ValueListLinkKind

// ValueNilKind is the name of the type used to nil lists of objects of
// type 'value'.
const ValueListNilKind = api_v1alpha1.ValueListNilKind

type ValueList = api_v1alpha1.ValueList
