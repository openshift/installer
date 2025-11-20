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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/osdfleetmgmt/v1"
)

// ManagementClusterKind is the name of the type used to represent objects
// of type 'management_cluster'.
const ManagementClusterKind = api_v1.ManagementClusterKind

// ManagementClusterLinkKind is the name of the type used to represent links
// to objects of type 'management_cluster'.
const ManagementClusterLinkKind = api_v1.ManagementClusterLinkKind

// ManagementClusterNilKind is the name of the type used to nil references
// to objects of type 'management_cluster'.
const ManagementClusterNilKind = api_v1.ManagementClusterNilKind

// ManagementCluster represents the values of the 'management_cluster' type.
//
// Definition of an _OpenShift_ cluster.
//
// The `cloud_provider` attribute is a reference to the cloud provider. When a
// cluster is retrieved it will be a link to the cloud provider, containing only
// the kind, id and href attributes:
//
// ```json
//
//	{
//	  "cloud_provider": {
//	    "kind": "CloudProviderLink",
//	    "id": "123",
//	    "href": "/api/clusters_mgmt/v1/cloud_providers/123"
//	  }
//	}
//
// ```
//
// When a cluster is created this is optional, and if used it should contain the
// identifier of the cloud provider to use:
//
// ```json
//
//	{
//	  "cloud_provider": {
//	    "id": "123",
//	  }
//	}
//
// ```
//
// If not included, then the cluster will be created using the default cloud
// provider, which is currently Amazon Web Services.
//
// The region attribute is mandatory when a cluster is created.
//
// The `aws.access_key_id`, `aws.secret_access_key` and `dns.base_domain`
// attributes are mandatory when creation a cluster with your own Amazon Web
// Services account.
type ManagementCluster = api_v1.ManagementCluster

// ManagementClusterListKind is the name of the type used to represent list of objects of
// type 'management_cluster'.
const ManagementClusterListKind = api_v1.ManagementClusterListKind

// ManagementClusterListLinkKind is the name of the type used to represent links to list
// of objects of type 'management_cluster'.
const ManagementClusterListLinkKind = api_v1.ManagementClusterListLinkKind

// ManagementClusterNilKind is the name of the type used to nil lists of objects of
// type 'management_cluster'.
const ManagementClusterListNilKind = api_v1.ManagementClusterListNilKind

type ManagementClusterList = api_v1.ManagementClusterList
