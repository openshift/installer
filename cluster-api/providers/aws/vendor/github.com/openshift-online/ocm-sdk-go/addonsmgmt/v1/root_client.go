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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

import (
	"net/http"
	"path"
)

// Client is the client of the 'root' resource.
//
// Root of the tree of resources of the Addons Management.
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

// Addons returns the target 'addons' resource.
//
// Reference to the resource that manages the collection of Addons.
func (c *Client) Addons() *AddonsClient {
	return NewAddonsClient(
		c.transport,
		path.Join(c.path, "addons"),
	)
}

// Clusters returns the target 'clusters' resource.
func (c *Client) Clusters() *ClustersClient {
	return NewClustersClient(
		c.transport,
		path.Join(c.path, "clusters"),
	)
}
