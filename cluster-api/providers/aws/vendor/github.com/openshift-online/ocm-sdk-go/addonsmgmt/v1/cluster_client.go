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

// ClusterClient is the client of the 'cluster' resource.
//
// Manages a specific cluster.
type ClusterClient struct {
	transport http.RoundTripper
	path      string
}

// NewClusterClient creates a new client for the 'cluster'
// resource using the given transport to send the requests and receive the
// responses.
func NewClusterClient(transport http.RoundTripper, path string) *ClusterClient {
	return &ClusterClient{
		transport: transport,
		path:      path,
	}
}

// AddonInquiries returns the target 'addon_inquiries' resource.
//
// Reference to the inquiries of addons on a specific cluster
func (c *ClusterClient) AddonInquiries() *AddonInquiriesClient {
	return NewAddonInquiriesClient(
		c.transport,
		path.Join(c.path, "addon_inquiries"),
	)
}

// Addons returns the target 'addon_installations' resource.
//
// Reference to the installations of addon on a specific cluster
func (c *ClusterClient) Addons() *AddonInstallationsClient {
	return NewAddonInstallationsClient(
		c.transport,
		path.Join(c.path, "addons"),
	)
}

// Status returns the target 'addon_statuses' resource.
//
// Reference to the status of addon installation on a specific cluster
func (c *ClusterClient) Status() *AddonStatusesClient {
	return NewAddonStatusesClient(
		c.transport,
		path.Join(c.path, "status"),
	)
}
