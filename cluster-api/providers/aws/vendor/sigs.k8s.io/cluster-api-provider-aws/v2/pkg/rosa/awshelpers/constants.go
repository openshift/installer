/*
Copyright 2026 The Kubernetes Authors.

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

// Package awshelpers provides helper functions for AWS resource management in ROSA.
package awshelpers

// Prefix used by all the tag names.
const prefix = "rosa_"

// Role types.
const (
	OCMRoleType = "OCM"
)

// Policy file names.
const (
	OCMRolePolicyFile          = "ocm"
	OCMAdminRolePolicyFile     = "ocm_admin"
	OCMNoConsoleRolePolicyFile = "ocm_no_console"
)

// IAM tag keys.
const (
	// TagRoleType is the name of the tag that will contain the purpose of the role (installer, support, etc.).
	TagRoleType = prefix + "role_type"

	// TagRolePrefix is the name of the tag that will contain the user-set prefix of the role (installer, support, etc.).
	TagRolePrefix = prefix + "role_prefix"

	// TagEnvironment is the name of the tag that will contain the environment of the role (integration/staging/production).
	TagEnvironment = prefix + "environment"

	// TagAdminRole tags the role as admin (true/false).
	TagAdminRole = prefix + "admin_role"

	// TagNoConsoleRole tags the role as no-console role.
	TagNoConsoleRole = prefix + "no_console_role"

	// TagRedHatManaged tags the role as red_hat_managed.
	TagRedHatManaged = "red-hat-managed"
)
