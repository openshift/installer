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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1"
)

// AccountGroupManagedBy represents the values of the 'account_group_managed_by' enumerated type.
type AccountGroupManagedBy = api_v1.AccountGroupManagedBy

const (
	// Group managed by OCM's API directly
	AccountGroupManagedByOCM AccountGroupManagedBy = api_v1.AccountGroupManagedByOCM
	// Group managed by remote RBAC service, synchronized by job
	AccountGroupManagedByRBAC AccountGroupManagedBy = api_v1.AccountGroupManagedByRBAC
)
