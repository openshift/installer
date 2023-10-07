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

package v1 // github.com/openshift-online/ocm-sdk-go/servicelogs/v1

import (
	"net/http"
	"path"
)

// ClustersClient is the client of the 'clusters' resource.
//
// Manages the collection of clusters for clusters logs.
type ClustersClient struct {
	transport http.RoundTripper
	path      string
}

// NewClustersClient creates a new client for the 'clusters'
// resource using the given transport to send the requests and receive the
// responses.
func NewClustersClient(transport http.RoundTripper, path string) *ClustersClient {
	return &ClustersClient{
		transport: transport,
		path:      path,
	}
}

// Cluster returns the target 'cluster' resource for the given identifier.
//
// Reference to the service that manages a specific Cluster uuid.
func (c *ClustersClient) Cluster(id string) *ClusterClient {
	return NewClusterClient(
		c.transport,
		path.Join(c.path, id),
	)
}
