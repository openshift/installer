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

// DeletedClusterKind is the name of the type used to represent objects
// of type 'deleted_cluster'.
const DeletedClusterKind = api_v1.DeletedClusterKind

// DeletedClusterLinkKind is the name of the type used to represent links
// to objects of type 'deleted_cluster'.
const DeletedClusterLinkKind = api_v1.DeletedClusterLinkKind

// DeletedClusterNilKind is the name of the type used to nil references
// to objects of type 'deleted_cluster'.
const DeletedClusterNilKind = api_v1.DeletedClusterNilKind

// DeletedCluster represents the values of the 'deleted_cluster' type.
//
// Representation of a deleted cluster with its deleted_timestamp and the entire cluster details
type DeletedCluster = api_v1.DeletedCluster

// DeletedClusterListKind is the name of the type used to represent list of objects of
// type 'deleted_cluster'.
const DeletedClusterListKind = api_v1.DeletedClusterListKind

// DeletedClusterListLinkKind is the name of the type used to represent links to list
// of objects of type 'deleted_cluster'.
const DeletedClusterListLinkKind = api_v1.DeletedClusterListLinkKind

// DeletedClusterNilKind is the name of the type used to nil lists of objects of
// type 'deleted_cluster'.
const DeletedClusterListNilKind = api_v1.DeletedClusterListNilKind

type DeletedClusterList = api_v1.DeletedClusterList
