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
	"path"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// NodePoolClient is the client of the 'node_pool' resource.
//
// Manages a specific nodepool.
type NodePoolClient struct {
	transport http.RoundTripper
	path      string
}

// NewNodePoolClient creates a new client for the 'node_pool'
// resource using the given transport to send the requests and receive the
// responses.
func NewNodePoolClient(transport http.RoundTripper, path string) *NodePoolClient {
	return &NodePoolClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the node pool.
func (c *NodePoolClient) Delete() *NodePoolDeleteRequest {
	return &NodePoolDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the node pool.
func (c *NodePoolClient) Get() *NodePoolGetRequest {
	return &NodePoolGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the node pool.
func (c *NodePoolClient) Update() *NodePoolUpdateRequest {
	return &NodePoolUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// UpgradePolicies returns the target 'node_pool_upgrade_policies' resource.
//
// Reference to the state of the upgrade policy.
func (c *NodePoolClient) UpgradePolicies() *NodePoolUpgradePoliciesClient {
	return NewNodePoolUpgradePoliciesClient(
		c.transport,
		path.Join(c.path, "upgrade_policies"),
	)
}

// NodePoolPollRequest is the request for the Poll method.
type NodePoolPollRequest struct {
	request    *NodePoolGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *NodePoolPollRequest) Parameter(name string, value interface{}) *NodePoolPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *NodePoolPollRequest) Header(name string, value interface{}) *NodePoolPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *NodePoolPollRequest) Interval(value time.Duration) *NodePoolPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *NodePoolPollRequest) Status(value int) *NodePoolPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *NodePoolPollRequest) Predicate(value func(*NodePoolGetResponse) bool) *NodePoolPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*NodePoolGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *NodePoolPollRequest) StartContext(ctx context.Context) (response *NodePoolPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &NodePoolPollResponse{
			response: result.(*NodePoolGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *NodePoolPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// NodePoolPollResponse is the response for the Poll method.
type NodePoolPollResponse struct {
	response *NodePoolGetResponse
}

// Status returns the response status code.
func (r *NodePoolPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *NodePoolPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *NodePoolPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *NodePoolPollResponse) Body() *NodePool {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *NodePoolPollResponse) GetBody() (value *NodePool, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *NodePoolClient) Poll() *NodePoolPollRequest {
	return &NodePoolPollRequest{
		request: c.Get(),
	}
}

// NodePoolDeleteRequest is the request for the 'delete' method.
type NodePoolDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *NodePoolDeleteRequest) Parameter(name string, value interface{}) *NodePoolDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NodePoolDeleteRequest) Header(name string, value interface{}) *NodePoolDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NodePoolDeleteRequest) Impersonate(user string) *NodePoolDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NodePoolDeleteRequest) Send() (result *NodePoolDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NodePoolDeleteRequest) SendContext(ctx context.Context) (result *NodePoolDeleteResponse, err error) {
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
	result = &NodePoolDeleteResponse{}
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

// NodePoolDeleteResponse is the response for the 'delete' method.
type NodePoolDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *NodePoolDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NodePoolDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NodePoolDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// NodePoolGetRequest is the request for the 'get' method.
type NodePoolGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *NodePoolGetRequest) Parameter(name string, value interface{}) *NodePoolGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NodePoolGetRequest) Header(name string, value interface{}) *NodePoolGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NodePoolGetRequest) Impersonate(user string) *NodePoolGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NodePoolGetRequest) Send() (result *NodePoolGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NodePoolGetRequest) SendContext(ctx context.Context) (result *NodePoolGetResponse, err error) {
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
	result = &NodePoolGetResponse{}
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
	err = readNodePoolGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NodePoolGetResponse is the response for the 'get' method.
type NodePoolGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *NodePool
}

// Status returns the response status code.
func (r *NodePoolGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NodePoolGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NodePoolGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *NodePoolGetResponse) Body() *NodePool {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *NodePoolGetResponse) GetBody() (value *NodePool, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// NodePoolUpdateRequest is the request for the 'update' method.
type NodePoolUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *NodePool
}

// Parameter adds a query parameter.
func (r *NodePoolUpdateRequest) Parameter(name string, value interface{}) *NodePoolUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NodePoolUpdateRequest) Header(name string, value interface{}) *NodePoolUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NodePoolUpdateRequest) Impersonate(user string) *NodePoolUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *NodePoolUpdateRequest) Body(value *NodePool) *NodePoolUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NodePoolUpdateRequest) Send() (result *NodePoolUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NodePoolUpdateRequest) SendContext(ctx context.Context) (result *NodePoolUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeNodePoolUpdateRequest(r, buffer)
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
	result = &NodePoolUpdateResponse{}
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
	err = readNodePoolUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NodePoolUpdateResponse is the response for the 'update' method.
type NodePoolUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *NodePool
}

// Status returns the response status code.
func (r *NodePoolUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NodePoolUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NodePoolUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *NodePoolUpdateResponse) Body() *NodePool {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *NodePoolUpdateResponse) GetBody() (value *NodePool, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
