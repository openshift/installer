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

// ManagedIdentitiesRequirementsBuilder contains the data and logic needed to build 'managed_identities_requirements' objects.
//
// Representation of managed identities requirements.
// When creating ARO-HCP Clusters, the end-users will need to pre-create the set of Managed Identities
// required by the clusters.
// The set of Managed Identities that the end-users need to precreate is not static and depends on
// several factors:
// (1) The OpenShift version of the cluster being created.
// (2) The functionalities that are being enabled for the cluster. Some Managed Identities are not
// always required but become required if a given functionality is enabled.
// Additionally, the Managed Identities that the end-users will need to precreate will have to have a
// set of required permissions assigned to them which also have to be returned to the end users.
type ManagedIdentitiesRequirementsBuilder = api_v1alpha1.ManagedIdentitiesRequirementsBuilder

// NewManagedIdentitiesRequirements creates a new builder of 'managed_identities_requirements' objects.
var NewManagedIdentitiesRequirements = api_v1alpha1.NewManagedIdentitiesRequirements
