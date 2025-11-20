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

// OidcConfig represents the values of the 'oidc_config' type.
//
// Contains the necessary attributes to support oidc configuration hosting under Red Hat or registering a Customer's byo oidc config.
type OidcConfig = api_v1.OidcConfig

// OidcConfigListKind is the name of the type used to represent list of objects of
// type 'oidc_config'.
const OidcConfigListKind = api_v1.OidcConfigListKind

// OidcConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'oidc_config'.
const OidcConfigListLinkKind = api_v1.OidcConfigListLinkKind

// OidcConfigNilKind is the name of the type used to nil lists of objects of
// type 'oidc_config'.
const OidcConfigListNilKind = api_v1.OidcConfigListNilKind

type OidcConfigList = api_v1.OidcConfigList
