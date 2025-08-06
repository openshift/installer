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

// ComponentRouteType represents the values of the 'component_route_type' enumerated type.
type ComponentRouteType = api_v1.ComponentRouteType

const (
	//
	ComponentRouteTypeConsole ComponentRouteType = api_v1.ComponentRouteTypeConsole
	//
	ComponentRouteTypeDownloads ComponentRouteType = api_v1.ComponentRouteTypeDownloads
	//
	ComponentRouteTypeOauth ComponentRouteType = api_v1.ComponentRouteTypeOauth
)
