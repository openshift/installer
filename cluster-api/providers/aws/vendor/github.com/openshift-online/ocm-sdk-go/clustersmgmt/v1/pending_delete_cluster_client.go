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
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// PendingDeleteClusterClient is the client of the 'pending_delete_cluster' resource.
//
// Manages a specific pending delete cluster.
type PendingDeleteClusterClient struct {
	transport http.RoundTripper
	path      string
}

// NewPendingDeleteClusterClient creates a new client for the 'pending_delete_cluster'
// resource using the given transport to send the requests and receive the
// responses.
func NewPendingDeleteClusterClient(transport http.RoundTripper, path string) *PendingDeleteClusterClient {
	return &PendingDeleteClusterClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the pending delete cluster.
func (c *PendingDeleteClusterClient) Get() *PendingDeleteClusterGetRequest {
	return &PendingDeleteClusterGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the pending delete cluster entry.
func (c *PendingDeleteClusterClient) Update() *PendingDeleteClusterUpdateRequest {
	return &PendingDeleteClusterUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// PendingDeleteClusterPollRequest is the request for the Poll method.
type PendingDeleteClusterPollRequest struct {
	request    *PendingDeleteClusterGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *PendingDeleteClusterPollRequest) Parameter(name string, value interface{}) *PendingDeleteClusterPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *PendingDeleteClusterPollRequest) Header(name string, value interface{}) *PendingDeleteClusterPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *PendingDeleteClusterPollRequest) Interval(value time.Duration) *PendingDeleteClusterPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *PendingDeleteClusterPollRequest) Status(value int) *PendingDeleteClusterPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *PendingDeleteClusterPollRequest) Predicate(value func(*PendingDeleteClusterGetResponse) bool) *PendingDeleteClusterPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*PendingDeleteClusterGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *PendingDeleteClusterPollRequest) StartContext(ctx context.Context) (response *PendingDeleteClusterPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &PendingDeleteClusterPollResponse{
			response: result.(*PendingDeleteClusterGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *PendingDeleteClusterPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// PendingDeleteClusterPollResponse is the response for the Poll method.
type PendingDeleteClusterPollResponse struct {
	response *PendingDeleteClusterGetResponse
}

// Status returns the response status code.
func (r *PendingDeleteClusterPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *PendingDeleteClusterPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *PendingDeleteClusterPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *PendingDeleteClusterPollResponse) Body() *PendingDeleteCluster {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *PendingDeleteClusterPollResponse) GetBody() (value *PendingDeleteCluster, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *PendingDeleteClusterClient) Poll() *PendingDeleteClusterPollRequest {
	return &PendingDeleteClusterPollRequest{
		request: c.Get(),
	}
}

// PendingDeleteClusterGetRequest is the request for the 'get' method.
type PendingDeleteClusterGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *PendingDeleteClusterGetRequest) Parameter(name string, value interface{}) *PendingDeleteClusterGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PendingDeleteClusterGetRequest) Header(name string, value interface{}) *PendingDeleteClusterGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PendingDeleteClusterGetRequest) Impersonate(user string) *PendingDeleteClusterGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PendingDeleteClusterGetRequest) Send() (result *PendingDeleteClusterGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PendingDeleteClusterGetRequest) SendContext(ctx context.Context) (result *PendingDeleteClusterGetResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "GET",
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
	result = &PendingDeleteClusterGetResponse{}
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
	err = readPendingDeleteClusterGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// PendingDeleteClusterGetResponse is the response for the 'get' method.
type PendingDeleteClusterGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *PendingDeleteCluster
}

// Status returns the response status code.
func (r *PendingDeleteClusterGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PendingDeleteClusterGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PendingDeleteClusterGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *PendingDeleteClusterGetResponse) Body() *PendingDeleteCluster {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *PendingDeleteClusterGetResponse) GetBody() (value *PendingDeleteCluster, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// PendingDeleteClusterUpdateRequest is the request for the 'update' method.
type PendingDeleteClusterUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *PendingDeleteCluster
}

// Parameter adds a query parameter.
func (r *PendingDeleteClusterUpdateRequest) Parameter(name string, value interface{}) *PendingDeleteClusterUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PendingDeleteClusterUpdateRequest) Header(name string, value interface{}) *PendingDeleteClusterUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PendingDeleteClusterUpdateRequest) Impersonate(user string) *PendingDeleteClusterUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *PendingDeleteClusterUpdateRequest) Body(value *PendingDeleteCluster) *PendingDeleteClusterUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PendingDeleteClusterUpdateRequest) Send() (result *PendingDeleteClusterUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PendingDeleteClusterUpdateRequest) SendContext(ctx context.Context) (result *PendingDeleteClusterUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writePendingDeleteClusterUpdateRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "PATCH",
		URL:    uri,
		Header: header,
		Body:   io.NopCloser(buffer),
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &PendingDeleteClusterUpdateResponse{}
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
	err = readPendingDeleteClusterUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// PendingDeleteClusterUpdateResponse is the response for the 'update' method.
type PendingDeleteClusterUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *PendingDeleteCluster
}

// Status returns the response status code.
func (r *PendingDeleteClusterUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PendingDeleteClusterUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PendingDeleteClusterUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *PendingDeleteClusterUpdateResponse) Body() *PendingDeleteCluster {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *PendingDeleteClusterUpdateResponse) GetBody() (value *PendingDeleteCluster, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
