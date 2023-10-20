// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package config

type AuthModeOption string

const (
	PodIdentityAuthMode      AuthModeOption = "podidentity"
	WorkloadIdentityAuthMode AuthModeOption = "workloadidentity"

	// AuthMode enum is used to determine if we're using Pod Identity or Workload Identity
	//authentication for namespace and per-resource scoped credentials
	AuthMode = "AUTH_MODE"
)
