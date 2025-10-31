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

// NodeInfo represents the values of the 'node_info' type.
//
// Provides information about a node from specific type in the cluster.
type NodeInfo = api_v1.NodeInfo

// NodeInfoListKind is the name of the type used to represent list of objects of
// type 'node_info'.
const NodeInfoListKind = api_v1.NodeInfoListKind

// NodeInfoListLinkKind is the name of the type used to represent links to list
// of objects of type 'node_info'.
const NodeInfoListLinkKind = api_v1.NodeInfoListLinkKind

// NodeInfoNilKind is the name of the type used to nil lists of objects of
// type 'node_info'.
const NodeInfoListNilKind = api_v1.NodeInfoListNilKind

type NodeInfoList = api_v1.NodeInfoList
