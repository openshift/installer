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

// ClusterRegistryConfig represents the values of the 'cluster_registry_config' type.
//
// ClusterRegistryConfig describes the configuration of registries for the cluster.
// Its format reflects the OpenShift Image Configuration, for which docs are available on
// [docs.openshift.com](https://docs.openshift.com/container-platform/4.16/openshift_images/image-configuration.html)
// ```json
//
//	{
//	   "registry_config": {
//	     "registry_sources": {
//	       "blocked_registries": [
//	         "badregistry.io",
//	         "badregistry8.io"
//	       ]
//	     }
//	   }
//	}
//
// ```
type ClusterRegistryConfig = api_v1.ClusterRegistryConfig

// ClusterRegistryConfigListKind is the name of the type used to represent list of objects of
// type 'cluster_registry_config'.
const ClusterRegistryConfigListKind = api_v1.ClusterRegistryConfigListKind

// ClusterRegistryConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_registry_config'.
const ClusterRegistryConfigListLinkKind = api_v1.ClusterRegistryConfigListLinkKind

// ClusterRegistryConfigNilKind is the name of the type used to nil lists of objects of
// type 'cluster_registry_config'.
const ClusterRegistryConfigListNilKind = api_v1.ClusterRegistryConfigListNilKind

type ClusterRegistryConfigList = api_v1.ClusterRegistryConfigList
