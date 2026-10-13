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
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// OperatorIAMRoleClient is the client of the 'operator_IAM_role' resource.
//
// Manages a list of operator roles for STS clusters.
type OperatorIAMRoleClient struct {
	transport http.RoundTripper
	path      string
}

// NewOperatorIAMRoleClient creates a new client for the 'operator_IAM_role'
// resource using the given transport to send the requests and receive the
// responses.
func NewOperatorIAMRoleClient(transport http.RoundTripper, path string) *OperatorIAMRoleClient {
	return &OperatorIAMRoleClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the operator role.
func (c *OperatorIAMRoleClient) Delete() *OperatorIAMRoleDeleteRequest {
	return &OperatorIAMRoleDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// OperatorIAMRoleDeleteRequest is the request for the 'delete' method.
type OperatorIAMRoleDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *OperatorIAMRoleDeleteRequest) Parameter(name string, value interface{}) *OperatorIAMRoleDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OperatorIAMRoleDeleteRequest) Header(name string, value interface{}) *OperatorIAMRoleDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OperatorIAMRoleDeleteRequest) Impersonate(user string) *OperatorIAMRoleDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OperatorIAMRoleDeleteRequest) Send() (result *OperatorIAMRoleDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OperatorIAMRoleDeleteRequest) SendContext(ctx context.Context) (result *OperatorIAMRoleDeleteResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "DELETE",
		URL:    uri,
		Header: header,
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &OperatorIAMRoleDeleteResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
		err = nil
		return
	}
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(reader, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	return
}

// OperatorIAMRoleDeleteResponse is the response for the 'delete' method.
type OperatorIAMRoleDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *OperatorIAMRoleDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OperatorIAMRoleDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OperatorIAMRoleDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}
