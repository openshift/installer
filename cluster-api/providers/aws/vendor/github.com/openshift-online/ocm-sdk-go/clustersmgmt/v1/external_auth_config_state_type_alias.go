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

// ExternalAuthConfigState represents the values of the 'external_auth_config_state' enumerated type.
type ExternalAuthConfigState = api_v1.ExternalAuthConfigState

const (
	// Indicates that the cluster does not support configuration of external authentication providers
	ExternalAuthConfigStateDisabled ExternalAuthConfigState = api_v1.ExternalAuthConfigStateDisabled
	// Indicates that the cluster supports configuration of external authentication providers
	ExternalAuthConfigStateEnabled ExternalAuthConfigState = api_v1.ExternalAuthConfigStateEnabled
)
