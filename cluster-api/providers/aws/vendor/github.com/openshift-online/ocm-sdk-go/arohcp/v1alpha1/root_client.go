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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	"net/http"
	"path"
)

// Client is the client of the 'root' resource.
//
// Root of the tree of resources of the aro_hcp service.
type Client struct {
	transport http.RoundTripper
	path      string
}

// NewClient creates a new client for the 'root'
// resource using the given transport to send the requests and receive the
// responses.
func NewClient(transport http.RoundTripper, path string) *Client {
	return &Client{
		transport: transport,
		path:      path,
	}
}

// Creates a new request for the method that retrieves the metadata.
func (c *Client) Get() *MetadataRequest {
	return &MetadataRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Clusters returns the target 'clusters' resource.
//
// Reference to the resource that manages the collection of clusters.
func (c *Client) Clusters() *ClustersClient {
	return NewClustersClient(
		c.transport,
		path.Join(c.path, "clusters"),
	)
}

// ManagedIdentitiesRequirements returns the target 'managed_identities_requirements' resource.
func (c *Client) ManagedIdentitiesRequirements() *ManagedIdentitiesRequirementsClient {
	return NewManagedIdentitiesRequirementsClient(
		c.transport,
		path.Join(c.path, "managed_identities_requirements"),
	)
}

// Versions returns the target 'versions' resource.
//
// Reference to the resource that manage the collection of versions.
func (c *Client) Versions() *VersionsClient {
	return NewVersionsClient(
		c.transport,
		path.Join(c.path, "versions"),
	)
}
