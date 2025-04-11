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

// AWSClient is the client of the 'AWS' resource.
//
// Manages AWS specific parts for a specific cluster.
type AWSClient struct {
	transport http.RoundTripper
	path      string
}

// NewAWSClient creates a new client for the 'AWS'
// resource using the given transport to send the requests and receive the
// responses.
func NewAWSClient(transport http.RoundTripper, path string) *AWSClient {
	return &AWSClient{
		transport: transport,
		path:      path,
	}
}

// PrivateLinkConfiguration returns the target 'private_link_configuration' resource.
func (c *AWSClient) PrivateLinkConfiguration() *PrivateLinkConfigurationClient {
	return NewPrivateLinkConfigurationClient(
		c.transport,
		path.Join(c.path, "private_link_configuration"),
	)
}

// RolePolicyBindings returns the target 'role_policy_bindings' resource.
func (c *AWSClient) RolePolicyBindings() *RolePolicyBindingsClient {
	return NewRolePolicyBindingsClient(
		c.transport,
		path.Join(c.path, "role_policy_bindings"),
	)
}
