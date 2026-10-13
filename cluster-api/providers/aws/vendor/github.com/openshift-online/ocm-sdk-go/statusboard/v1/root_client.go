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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	"net/http"
	"path"
)

// Client is the client of the 'root' resource.
//
// Root of the tree of resources for applications.
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

// ApplicationDependencies returns the target 'application_dependencies' resource.
func (c *Client) ApplicationDependencies() *ApplicationDependenciesClient {
	return NewApplicationDependenciesClient(
		c.transport,
		path.Join(c.path, "application_dependencies"),
	)
}

// Applications returns the target 'applications' resource.
func (c *Client) Applications() *ApplicationsClient {
	return NewApplicationsClient(
		c.transport,
		path.Join(c.path, "applications"),
	)
}

// Errors returns the target 'errors' resource.
func (c *Client) Errors() *ErrorsClient {
	return NewErrorsClient(
		c.transport,
		path.Join(c.path, "errors"),
	)
}

// PeerDependencies returns the target 'peer_dependencies' resource.
func (c *Client) PeerDependencies() *PeerDependenciesClient {
	return NewPeerDependenciesClient(
		c.transport,
		path.Join(c.path, "peer_dependencies"),
	)
}

// Products returns the target 'products' resource.
func (c *Client) Products() *ProductsClient {
	return NewProductsClient(
		c.transport,
		path.Join(c.path, "products"),
	)
}

// Services returns the target 'services' resource.
func (c *Client) Services() *ServicesClient {
	return NewServicesClient(
		c.transport,
		path.Join(c.path, "services"),
	)
}

// StatusUpdates returns the target 'statuses' resource.
func (c *Client) StatusUpdates() *StatusesClient {
	return NewStatusesClient(
		c.transport,
		path.Join(c.path, "status_updates"),
	)
}

// Statuses returns the target 'statuses' resource.
func (c *Client) Statuses() *StatusesClient {
	return NewStatusesClient(
		c.transport,
		path.Join(c.path, "statuses"),
	)
}
