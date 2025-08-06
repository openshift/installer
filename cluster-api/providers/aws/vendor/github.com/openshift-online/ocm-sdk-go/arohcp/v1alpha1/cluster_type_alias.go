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

// ClusterKind is the name of the type used to represent objects
// of type 'cluster'.
const ClusterKind = api_v1alpha1.ClusterKind

// ClusterLinkKind is the name of the type used to represent links
// to objects of type 'cluster'.
const ClusterLinkKind = api_v1alpha1.ClusterLinkKind

// ClusterNilKind is the name of the type used to nil references
// to objects of type 'cluster'.
const ClusterNilKind = api_v1alpha1.ClusterNilKind

// Cluster represents the values of the 'cluster' type.
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
type Cluster = api_v1alpha1.Cluster

// ClusterListKind is the name of the type used to represent list of objects of
// type 'cluster'.
const ClusterListKind = api_v1alpha1.ClusterListKind

// ClusterListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster'.
const ClusterListLinkKind = api_v1alpha1.ClusterListLinkKind

// ClusterNilKind is the name of the type used to nil lists of objects of
// type 'cluster'.
const ClusterListNilKind = api_v1alpha1.ClusterListNilKind

type ClusterList = api_v1alpha1.ClusterList
