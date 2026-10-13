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

// NodePoolUpgradePolicyClient is the client of the 'node_pool_upgrade_policy' resource.
//
// Manages a specific upgrade policy for the node pool.
type NodePoolUpgradePolicyClient struct {
	transport http.RoundTripper
	path      string
}

// NewNodePoolUpgradePolicyClient creates a new client for the 'node_pool_upgrade_policy'
// resource using the given transport to send the requests and receive the
// responses.
func NewNodePoolUpgradePolicyClient(transport http.RoundTripper, path string) *NodePoolUpgradePolicyClient {
	return &NodePoolUpgradePolicyClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the upgrade policy for the node pool.
func (c *NodePoolUpgradePolicyClient) Delete() *NodePoolUpgradePolicyDeleteRequest {
	return &NodePoolUpgradePolicyDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the upgrade policy for the node pool.
func (c *NodePoolUpgradePolicyClient) Get() *NodePoolUpgradePolicyGetRequest {
	return &NodePoolUpgradePolicyGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Update the upgrade policy for the node pool.
func (c *NodePoolUpgradePolicyClient) Update() *NodePoolUpgradePolicyUpdateRequest {
	return &NodePoolUpgradePolicyUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// NodePoolUpgradePolicyPollRequest is the request for the Poll method.
type NodePoolUpgradePolicyPollRequest struct {
	request    *NodePoolUpgradePolicyGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *NodePoolUpgradePolicyPollRequest) Parameter(name string, value interface{}) *NodePoolUpgradePolicyPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *NodePoolUpgradePolicyPollRequest) Header(name string, value interface{}) *NodePoolUpgradePolicyPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *NodePoolUpgradePolicyPollRequest) Interval(value time.Duration) *NodePoolUpgradePolicyPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *NodePoolUpgradePolicyPollRequest) Status(value int) *NodePoolUpgradePolicyPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *NodePoolUpgradePolicyPollRequest) Predicate(value func(*NodePoolUpgradePolicyGetResponse) bool) *NodePoolUpgradePolicyPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*NodePoolUpgradePolicyGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *NodePoolUpgradePolicyPollRequest) StartContext(ctx context.Context) (response *NodePoolUpgradePolicyPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &NodePoolUpgradePolicyPollResponse{
			response: result.(*NodePoolUpgradePolicyGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *NodePoolUpgradePolicyPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// NodePoolUpgradePolicyPollResponse is the response for the Poll method.
type NodePoolUpgradePolicyPollResponse struct {
	response *NodePoolUpgradePolicyGetResponse
}

// Status returns the response status code.
func (r *NodePoolUpgradePolicyPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *NodePoolUpgradePolicyPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *NodePoolUpgradePolicyPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *NodePoolUpgradePolicyPollResponse) Body() *NodePoolUpgradePolicy {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *NodePoolUpgradePolicyPollResponse) GetBody() (value *NodePoolUpgradePolicy, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *NodePoolUpgradePolicyClient) Poll() *NodePoolUpgradePolicyPollRequest {
	return &NodePoolUpgradePolicyPollRequest{
		request: c.Get(),
	}
}

// NodePoolUpgradePolicyDeleteRequest is the request for the 'delete' method.
type NodePoolUpgradePolicyDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *NodePoolUpgradePolicyDeleteRequest) Parameter(name string, value interface{}) *NodePoolUpgradePolicyDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NodePoolUpgradePolicyDeleteRequest) Header(name string, value interface{}) *NodePoolUpgradePolicyDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NodePoolUpgradePolicyDeleteRequest) Impersonate(user string) *NodePoolUpgradePolicyDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NodePoolUpgradePolicyDeleteRequest) Send() (result *NodePoolUpgradePolicyDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NodePoolUpgradePolicyDeleteRequest) SendContext(ctx context.Context) (result *NodePoolUpgradePolicyDeleteResponse, err error) {
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
	result = &NodePoolUpgradePolicyDeleteResponse{}
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

// NodePoolUpgradePolicyDeleteResponse is the response for the 'delete' method.
type NodePoolUpgradePolicyDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *NodePoolUpgradePolicyDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NodePoolUpgradePolicyDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NodePoolUpgradePolicyDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// NodePoolUpgradePolicyGetRequest is the request for the 'get' method.
type NodePoolUpgradePolicyGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *NodePoolUpgradePolicyGetRequest) Parameter(name string, value interface{}) *NodePoolUpgradePolicyGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NodePoolUpgradePolicyGetRequest) Header(name string, value interface{}) *NodePoolUpgradePolicyGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NodePoolUpgradePolicyGetRequest) Impersonate(user string) *NodePoolUpgradePolicyGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NodePoolUpgradePolicyGetRequest) Send() (result *NodePoolUpgradePolicyGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NodePoolUpgradePolicyGetRequest) SendContext(ctx context.Context) (result *NodePoolUpgradePolicyGetResponse, err error) {
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
	result = &NodePoolUpgradePolicyGetResponse{}
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
	err = readNodePoolUpgradePolicyGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NodePoolUpgradePolicyGetResponse is the response for the 'get' method.
type NodePoolUpgradePolicyGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *NodePoolUpgradePolicy
}

// Status returns the response status code.
func (r *NodePoolUpgradePolicyGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NodePoolUpgradePolicyGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NodePoolUpgradePolicyGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *NodePoolUpgradePolicyGetResponse) Body() *NodePoolUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *NodePoolUpgradePolicyGetResponse) GetBody() (value *NodePoolUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// NodePoolUpgradePolicyUpdateRequest is the request for the 'update' method.
type NodePoolUpgradePolicyUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *NodePoolUpgradePolicy
}

// Parameter adds a query parameter.
func (r *NodePoolUpgradePolicyUpdateRequest) Parameter(name string, value interface{}) *NodePoolUpgradePolicyUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NodePoolUpgradePolicyUpdateRequest) Header(name string, value interface{}) *NodePoolUpgradePolicyUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NodePoolUpgradePolicyUpdateRequest) Impersonate(user string) *NodePoolUpgradePolicyUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *NodePoolUpgradePolicyUpdateRequest) Body(value *NodePoolUpgradePolicy) *NodePoolUpgradePolicyUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NodePoolUpgradePolicyUpdateRequest) Send() (result *NodePoolUpgradePolicyUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NodePoolUpgradePolicyUpdateRequest) SendContext(ctx context.Context) (result *NodePoolUpgradePolicyUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeNodePoolUpgradePolicyUpdateRequest(r, buffer)
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
	result = &NodePoolUpgradePolicyUpdateResponse{}
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
	err = readNodePoolUpgradePolicyUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NodePoolUpgradePolicyUpdateResponse is the response for the 'update' method.
type NodePoolUpgradePolicyUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *NodePoolUpgradePolicy
}

// Status returns the response status code.
func (r *NodePoolUpgradePolicyUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NodePoolUpgradePolicyUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NodePoolUpgradePolicyUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *NodePoolUpgradePolicyUpdateResponse) Body() *NodePoolUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *NodePoolUpgradePolicyUpdateResponse) GetBody() (value *NodePoolUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
