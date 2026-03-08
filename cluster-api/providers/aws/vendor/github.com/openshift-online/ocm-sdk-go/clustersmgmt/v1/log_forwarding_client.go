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

// LogForwardingClient is the client of the 'log_forwarding' resource.
//
// Manages log forwarding related resources.
type LogForwardingClient struct {
	transport http.RoundTripper
	path      string
}

// NewLogForwardingClient creates a new client for the 'log_forwarding'
// resource using the given transport to send the requests and receive the
// responses.
func NewLogForwardingClient(transport http.RoundTripper, path string) *LogForwardingClient {
	return &LogForwardingClient{
		transport: transport,
		path:      path,
	}
}

// Applications returns the target 'log_forwarding_applications' resource.
//
// Reference to the collection of log forwarder applications.
func (c *LogForwardingClient) Applications() *LogForwardingApplicationsClient {
	return NewLogForwardingApplicationsClient(
		c.transport,
		path.Join(c.path, "applications"),
	)
}

// Groups returns the target 'log_forwarding_groups' resource.
//
// Reference to the collection of log forwarding group versions.
func (c *LogForwardingClient) Groups() *LogForwardingGroupsClient {
	return NewLogForwardingGroupsClient(
		c.transport,
		path.Join(c.path, "groups"),
	)
}
