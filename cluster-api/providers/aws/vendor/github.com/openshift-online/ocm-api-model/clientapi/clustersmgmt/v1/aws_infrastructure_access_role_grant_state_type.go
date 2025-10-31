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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// AWSInfrastructureAccessRoleGrantState represents the values of the 'AWS_infrastructure_access_role_grant_state' enumerated type.
type AWSInfrastructureAccessRoleGrantState string

const (
	// The access role grant is in the process of being deleted.
	AWSInfrastructureAccessRoleGrantStateDeleting AWSInfrastructureAccessRoleGrantState = "deleting"
	// The attempt to grant access role to user ARN failed.
	AWSInfrastructureAccessRoleGrantStateFailed AWSInfrastructureAccessRoleGrantState = "failed"
	// The access role grant in pending.
	AWSInfrastructureAccessRoleGrantStatePending AWSInfrastructureAccessRoleGrantState = "pending"
	// Access role has been granted to user.
	AWSInfrastructureAccessRoleGrantStateReady AWSInfrastructureAccessRoleGrantState = "ready"
	// This ia a special state intended for the user know
	// that the access role grant has been removed by SRE.
	// The user can delete this grant from the DB.
	AWSInfrastructureAccessRoleGrantStateRemoved AWSInfrastructureAccessRoleGrantState = "removed"
)
