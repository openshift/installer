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

// CIDRBlockAccess represents the values of the 'CIDR_block_access' type.
//
// Describes the CIDR Block access policy to the Kubernetes API server.
// Currently, only supported for ARO-HCP based clusters.
// The default policy mode is "allow_all" that is, all access is allowed.
type CIDRBlockAccess = api_v1alpha1.CIDRBlockAccess

// CIDRBlockAccessListKind is the name of the type used to represent list of objects of
// type 'CIDR_block_access'.
const CIDRBlockAccessListKind = api_v1alpha1.CIDRBlockAccessListKind

// CIDRBlockAccessListLinkKind is the name of the type used to represent links to list
// of objects of type 'CIDR_block_access'.
const CIDRBlockAccessListLinkKind = api_v1alpha1.CIDRBlockAccessListLinkKind

// CIDRBlockAccessNilKind is the name of the type used to nil lists of objects of
// type 'CIDR_block_access'.
const CIDRBlockAccessListNilKind = api_v1alpha1.CIDRBlockAccessListNilKind

type CIDRBlockAccessList = api_v1alpha1.CIDRBlockAccessList
