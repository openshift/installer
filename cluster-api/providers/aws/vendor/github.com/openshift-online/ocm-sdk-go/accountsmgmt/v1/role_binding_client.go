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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

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

// RoleBindingClient is the client of the 'role_binding' resource.
//
// Manages a specific role binding.
type RoleBindingClient struct {
	transport http.RoundTripper
	path      string
}

// NewRoleBindingClient creates a new client for the 'role_binding'
// resource using the given transport to send the requests and receive the
// responses.
func NewRoleBindingClient(transport http.RoundTripper, path string) *RoleBindingClient {
	return &RoleBindingClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the role binding.
func (c *RoleBindingClient) Delete() *RoleBindingDeleteRequest {
	return &RoleBindingDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the role binding.
func (c *RoleBindingClient) Get() *RoleBindingGetRequest {
	return &RoleBindingGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the account.
func (c *RoleBindingClient) Update() *RoleBindingUpdateRequest {
	return &RoleBindingUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// RoleBindingPollRequest is the request for the Poll method.
type RoleBindingPollRequest struct {
	request    *RoleBindingGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *RoleBindingPollRequest) Parameter(name string, value interface{}) *RoleBindingPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *RoleBindingPollRequest) Header(name string, value interface{}) *RoleBindingPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *RoleBindingPollRequest) Interval(value time.Duration) *RoleBindingPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *RoleBindingPollRequest) Status(value int) *RoleBindingPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *RoleBindingPollRequest) Predicate(value func(*RoleBindingGetResponse) bool) *RoleBindingPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*RoleBindingGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *RoleBindingPollRequest) StartContext(ctx context.Context) (response *RoleBindingPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &RoleBindingPollResponse{
			response: result.(*RoleBindingGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *RoleBindingPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// RoleBindingPollResponse is the response for the Poll method.
type RoleBindingPollResponse struct {
	response *RoleBindingGetResponse
}

// Status returns the response status code.
func (r *RoleBindingPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *RoleBindingPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *RoleBindingPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *RoleBindingPollResponse) Body() *RoleBinding {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *RoleBindingPollResponse) GetBody() (value *RoleBinding, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *RoleBindingClient) Poll() *RoleBindingPollRequest {
	return &RoleBindingPollRequest{
		request: c.Get(),
	}
}

// RoleBindingDeleteRequest is the request for the 'delete' method.
type RoleBindingDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *RoleBindingDeleteRequest) Parameter(name string, value interface{}) *RoleBindingDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RoleBindingDeleteRequest) Header(name string, value interface{}) *RoleBindingDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RoleBindingDeleteRequest) Impersonate(user string) *RoleBindingDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RoleBindingDeleteRequest) Send() (result *RoleBindingDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RoleBindingDeleteRequest) SendContext(ctx context.Context) (result *RoleBindingDeleteResponse, err error) {
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
	result = &RoleBindingDeleteResponse{}
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

// RoleBindingDeleteResponse is the response for the 'delete' method.
type RoleBindingDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *RoleBindingDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RoleBindingDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RoleBindingDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// RoleBindingGetRequest is the request for the 'get' method.
type RoleBindingGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *RoleBindingGetRequest) Parameter(name string, value interface{}) *RoleBindingGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RoleBindingGetRequest) Header(name string, value interface{}) *RoleBindingGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RoleBindingGetRequest) Impersonate(user string) *RoleBindingGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RoleBindingGetRequest) Send() (result *RoleBindingGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RoleBindingGetRequest) SendContext(ctx context.Context) (result *RoleBindingGetResponse, err error) {
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
	result = &RoleBindingGetResponse{}
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
	err = readRoleBindingGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// RoleBindingGetResponse is the response for the 'get' method.
type RoleBindingGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *RoleBinding
}

// Status returns the response status code.
func (r *RoleBindingGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RoleBindingGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RoleBindingGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *RoleBindingGetResponse) Body() *RoleBinding {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *RoleBindingGetResponse) GetBody() (value *RoleBinding, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// RoleBindingUpdateRequest is the request for the 'update' method.
type RoleBindingUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *RoleBinding
}

// Parameter adds a query parameter.
func (r *RoleBindingUpdateRequest) Parameter(name string, value interface{}) *RoleBindingUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RoleBindingUpdateRequest) Header(name string, value interface{}) *RoleBindingUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RoleBindingUpdateRequest) Impersonate(user string) *RoleBindingUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *RoleBindingUpdateRequest) Body(value *RoleBinding) *RoleBindingUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RoleBindingUpdateRequest) Send() (result *RoleBindingUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RoleBindingUpdateRequest) SendContext(ctx context.Context) (result *RoleBindingUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeRoleBindingUpdateRequest(r, buffer)
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
	result = &RoleBindingUpdateResponse{}
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
	err = readRoleBindingUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// RoleBindingUpdateResponse is the response for the 'update' method.
type RoleBindingUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *RoleBinding
}

// Status returns the response status code.
func (r *RoleBindingUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RoleBindingUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RoleBindingUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *RoleBindingUpdateResponse) Body() *RoleBinding {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *RoleBindingUpdateResponse) GetBody() (value *RoleBinding, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
