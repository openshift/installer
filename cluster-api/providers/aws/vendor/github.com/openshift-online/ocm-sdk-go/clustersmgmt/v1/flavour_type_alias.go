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

// FlavourKind is the name of the type used to represent objects
// of type 'flavour'.
const FlavourKind = api_v1.FlavourKind

// FlavourLinkKind is the name of the type used to represent links
// to objects of type 'flavour'.
const FlavourLinkKind = api_v1.FlavourLinkKind

// FlavourNilKind is the name of the type used to nil references
// to objects of type 'flavour'.
const FlavourNilKind = api_v1.FlavourNilKind

// Flavour represents the values of the 'flavour' type.
//
// Set of predefined properties of a cluster. For example, a _huge_ flavour can be a cluster
// with 10 infra nodes and 1000 compute nodes.
type Flavour = api_v1.Flavour

// FlavourListKind is the name of the type used to represent list of objects of
// type 'flavour'.
const FlavourListKind = api_v1.FlavourListKind

// FlavourListLinkKind is the name of the type used to represent links to list
// of objects of type 'flavour'.
const FlavourListLinkKind = api_v1.FlavourListLinkKind

// FlavourNilKind is the name of the type used to nil lists of objects of
// type 'flavour'.
const FlavourListNilKind = api_v1.FlavourListNilKind

type FlavourList = api_v1.FlavourList
