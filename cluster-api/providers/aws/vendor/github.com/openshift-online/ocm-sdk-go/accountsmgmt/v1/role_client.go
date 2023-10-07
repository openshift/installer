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

// RoleClient is the client of the 'role' resource.
//
// Manages a specific role.
type RoleClient struct {
	transport http.RoundTripper
	path      string
}

// NewRoleClient creates a new client for the 'role'
// resource using the given transport to send the requests and receive the
// responses.
func NewRoleClient(transport http.RoundTripper, path string) *RoleClient {
	return &RoleClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the role.
func (c *RoleClient) Delete() *RoleDeleteRequest {
	return &RoleDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the role.
func (c *RoleClient) Get() *RoleGetRequest {
	return &RoleGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the role.
func (c *RoleClient) Update() *RoleUpdateRequest {
	return &RoleUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// RolePollRequest is the request for the Poll method.
type RolePollRequest struct {
	request    *RoleGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *RolePollRequest) Parameter(name string, value interface{}) *RolePollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *RolePollRequest) Header(name string, value interface{}) *RolePollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *RolePollRequest) Interval(value time.Duration) *RolePollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *RolePollRequest) Status(value int) *RolePollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *RolePollRequest) Predicate(value func(*RoleGetResponse) bool) *RolePollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*RoleGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *RolePollRequest) StartContext(ctx context.Context) (response *RolePollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &RolePollResponse{
			response: result.(*RoleGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *RolePollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// RolePollResponse is the response for the Poll method.
type RolePollResponse struct {
	response *RoleGetResponse
}

// Status returns the response status code.
func (r *RolePollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *RolePollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *RolePollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *RolePollResponse) Body() *Role {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *RolePollResponse) GetBody() (value *Role, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *RoleClient) Poll() *RolePollRequest {
	return &RolePollRequest{
		request: c.Get(),
	}
}

// RoleDeleteRequest is the request for the 'delete' method.
type RoleDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *RoleDeleteRequest) Parameter(name string, value interface{}) *RoleDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RoleDeleteRequest) Header(name string, value interface{}) *RoleDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RoleDeleteRequest) Impersonate(user string) *RoleDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RoleDeleteRequest) Send() (result *RoleDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RoleDeleteRequest) SendContext(ctx context.Context) (result *RoleDeleteResponse, err error) {
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
	result = &RoleDeleteResponse{}
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

// RoleDeleteResponse is the response for the 'delete' method.
type RoleDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *RoleDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RoleDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RoleDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// RoleGetRequest is the request for the 'get' method.
type RoleGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *RoleGetRequest) Parameter(name string, value interface{}) *RoleGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RoleGetRequest) Header(name string, value interface{}) *RoleGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RoleGetRequest) Impersonate(user string) *RoleGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RoleGetRequest) Send() (result *RoleGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RoleGetRequest) SendContext(ctx context.Context) (result *RoleGetResponse, err error) {
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
	result = &RoleGetResponse{}
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
	err = readRoleGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// RoleGetResponse is the response for the 'get' method.
type RoleGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Role
}

// Status returns the response status code.
func (r *RoleGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RoleGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RoleGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *RoleGetResponse) Body() *Role {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *RoleGetResponse) GetBody() (value *Role, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// RoleUpdateRequest is the request for the 'update' method.
type RoleUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Role
}

// Parameter adds a query parameter.
func (r *RoleUpdateRequest) Parameter(name string, value interface{}) *RoleUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RoleUpdateRequest) Header(name string, value interface{}) *RoleUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RoleUpdateRequest) Impersonate(user string) *RoleUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *RoleUpdateRequest) Body(value *Role) *RoleUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RoleUpdateRequest) Send() (result *RoleUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RoleUpdateRequest) SendContext(ctx context.Context) (result *RoleUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeRoleUpdateRequest(r, buffer)
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
	result = &RoleUpdateResponse{}
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
	err = readRoleUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// RoleUpdateResponse is the response for the 'update' method.
type RoleUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Role
}

// Status returns the response status code.
func (r *RoleUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RoleUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RoleUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *RoleUpdateResponse) Body() *Role {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *RoleUpdateResponse) GetBody() (value *Role, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
