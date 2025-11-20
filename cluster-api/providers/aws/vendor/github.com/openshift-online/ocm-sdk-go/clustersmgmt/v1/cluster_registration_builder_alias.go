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

// ClusterRegistrationBuilder contains the data and logic needed to build 'cluster_registration' objects.
//
// Registration of a new cluster to the service.
//
// For example, to register a cluster that has been provisioned outside
// of this service, send a a request like this:
//
// ```http
// POST /api/clusters_mgmt/v1/register_cluster HTTP/1.1
// ```
//
// With a request body like this:
//
// ```json
//
//	{
//	  "external_id": "d656aecf-11a6-4782-ad86-8f72638449ba",
//	  "subscription_id": "...",
//	  "organization_id": "..."
//	}
//
// ```
type ClusterRegistrationBuilder = api_v1.ClusterRegistrationBuilder

// NewClusterRegistration creates a new builder of 'cluster_registration' objects.
var NewClusterRegistration = api_v1.NewClusterRegistration
