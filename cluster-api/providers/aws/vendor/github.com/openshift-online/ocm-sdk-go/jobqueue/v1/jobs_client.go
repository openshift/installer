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

package v1 // github.com/openshift-online/ocm-sdk-go/jobqueue/v1

import (
	"net/http"
	"path"
)

// JobsClient is the client of the 'jobs' resource.
//
// Manages status of jobs on a job queue.
type JobsClient struct {
	transport http.RoundTripper
	path      string
}

// NewJobsClient creates a new client for the 'jobs'
// resource using the given transport to send the requests and receive the
// responses.
func NewJobsClient(transport http.RoundTripper, path string) *JobsClient {
	return &JobsClient{
		transport: transport,
		path:      path,
	}
}

// Job returns the target 'job' resource for the given identifier.
//
// jobs' operations (success, failure)
func (c *JobsClient) Job(id string) *JobClient {
	return NewJobClient(
		c.transport,
		path.Join(c.path, id),
	)
}
