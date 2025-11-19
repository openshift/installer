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

// GCPNetwork represents the values of the 'GCP_network' type.
//
// GCP Network configuration of a cluster.
type GCPNetwork = api_v1alpha1.GCPNetwork

// GCPNetworkListKind is the name of the type used to represent list of objects of
// type 'GCP_network'.
const GCPNetworkListKind = api_v1alpha1.GCPNetworkListKind

// GCPNetworkListLinkKind is the name of the type used to represent links to list
// of objects of type 'GCP_network'.
const GCPNetworkListLinkKind = api_v1alpha1.GCPNetworkListLinkKind

// GCPNetworkNilKind is the name of the type used to nil lists of objects of
// type 'GCP_network'.
const GCPNetworkListNilKind = api_v1alpha1.GCPNetworkListNilKind

type GCPNetworkList = api_v1alpha1.GCPNetworkList
