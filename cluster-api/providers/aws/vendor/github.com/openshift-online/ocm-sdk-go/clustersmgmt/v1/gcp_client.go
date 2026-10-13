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

// GCPClient is the client of the 'GCP' resource.
//
// Manages the collection of gcp endpoints.
type GCPClient struct {
	transport http.RoundTripper
	path      string
}

// NewGCPClient creates a new client for the 'GCP'
// resource using the given transport to send the requests and receive the
// responses.
func NewGCPClient(transport http.RoundTripper, path string) *GCPClient {
	return &GCPClient{
		transport: transport,
		path:      path,
	}
}

// WifConfigs returns the target 'wif_configs' resource.
//
// Reference to the resource that manages wif_configs
func (c *GCPClient) WifConfigs() *WifConfigsClient {
	return NewWifConfigsClient(
		c.transport,
		path.Join(c.path, "wif_configs"),
	)
}
