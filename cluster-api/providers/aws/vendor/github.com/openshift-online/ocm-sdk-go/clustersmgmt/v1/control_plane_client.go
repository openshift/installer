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
	"net/http"
	"path"
)

// ControlPlaneClient is the client of the 'control_plane' resource.
//
// Manages a specific upgrade policy.
type ControlPlaneClient struct {
	transport http.RoundTripper
	path      string
}

// NewControlPlaneClient creates a new client for the 'control_plane'
// resource using the given transport to send the requests and receive the
// responses.
func NewControlPlaneClient(transport http.RoundTripper, path string) *ControlPlaneClient {
	return &ControlPlaneClient{
		transport: transport,
		path:      path,
	}
}

// UpgradePolicies returns the target 'control_plane_upgrade_policies' resource.
//
// Reference to the state of the upgrade policy.
func (c *ControlPlaneClient) UpgradePolicies() *ControlPlaneUpgradePoliciesClient {
	return NewControlPlaneUpgradePoliciesClient(
		c.transport,
		path.Join(c.path, "upgrade_policies"),
	)
}
