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

// ExternalAuthClientType represents the values of the 'external_auth_client_type' enumerated type.
type ExternalAuthClientType = api_v1alpha1.ExternalAuthClientType

const (
	// Indicates that the client is confidential.
	//
	// Confidential clients must provide a client secret.
	// For external authentication provider belonging to a ROSA HCP cluster, the secret should be provided
	// in the 'secret' property of the client configuration.
	// For those belonging to an ARO-HCP cluster, the secret should be provided within the cluster itself.
	ExternalAuthClientTypeConfidential ExternalAuthClientType = api_v1alpha1.ExternalAuthClientTypeConfidential
	// Indicates that the client is public
	//
	// Public clients must not provide a client secret
	ExternalAuthClientTypePublic ExternalAuthClientType = api_v1alpha1.ExternalAuthClientTypePublic
)
