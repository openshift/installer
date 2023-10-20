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

// MachinePoolClient is the client of the 'machine_pool' resource.
//
// Manages a specific machine pool.
type MachinePoolClient struct {
	transport http.RoundTripper
	path      string
}

// NewMachinePoolClient creates a new client for the 'machine_pool'
// resource using the given transport to send the requests and receive the
// responses.
func NewMachinePoolClient(transport http.RoundTripper, path string) *MachinePoolClient {
	return &MachinePoolClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the machine pool.
func (c *MachinePoolClient) Delete() *MachinePoolDeleteRequest {
	return &MachinePoolDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the machine pool.
func (c *MachinePoolClient) Get() *MachinePoolGetRequest {
	return &MachinePoolGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the machine pool.
func (c *MachinePoolClient) Update() *MachinePoolUpdateRequest {
	return &MachinePoolUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// MachinePoolPollRequest is the request for the Poll method.
type MachinePoolPollRequest struct {
	request    *MachinePoolGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *MachinePoolPollRequest) Parameter(name string, value interface{}) *MachinePoolPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *MachinePoolPollRequest) Header(name string, value interface{}) *MachinePoolPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *MachinePoolPollRequest) Interval(value time.Duration) *MachinePoolPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *MachinePoolPollRequest) Status(value int) *MachinePoolPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *MachinePoolPollRequest) Predicate(value func(*MachinePoolGetResponse) bool) *MachinePoolPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*MachinePoolGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *MachinePoolPollRequest) StartContext(ctx context.Context) (response *MachinePoolPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &MachinePoolPollResponse{
			response: result.(*MachinePoolGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *MachinePoolPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// MachinePoolPollResponse is the response for the Poll method.
type MachinePoolPollResponse struct {
	response *MachinePoolGetResponse
}

// Status returns the response status code.
func (r *MachinePoolPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *MachinePoolPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *MachinePoolPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *MachinePoolPollResponse) Body() *MachinePool {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *MachinePoolPollResponse) GetBody() (value *MachinePool, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *MachinePoolClient) Poll() *MachinePoolPollRequest {
	return &MachinePoolPollRequest{
		request: c.Get(),
	}
}

// MachinePoolDeleteRequest is the request for the 'delete' method.
type MachinePoolDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *MachinePoolDeleteRequest) Parameter(name string, value interface{}) *MachinePoolDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *MachinePoolDeleteRequest) Header(name string, value interface{}) *MachinePoolDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *MachinePoolDeleteRequest) Impersonate(user string) *MachinePoolDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *MachinePoolDeleteRequest) Send() (result *MachinePoolDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *MachinePoolDeleteRequest) SendContext(ctx context.Context) (result *MachinePoolDeleteResponse, err error) {
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
	result = &MachinePoolDeleteResponse{}
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

// MachinePoolDeleteResponse is the response for the 'delete' method.
type MachinePoolDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *MachinePoolDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *MachinePoolDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *MachinePoolDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// MachinePoolGetRequest is the request for the 'get' method.
type MachinePoolGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *MachinePoolGetRequest) Parameter(name string, value interface{}) *MachinePoolGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *MachinePoolGetRequest) Header(name string, value interface{}) *MachinePoolGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *MachinePoolGetRequest) Impersonate(user string) *MachinePoolGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *MachinePoolGetRequest) Send() (result *MachinePoolGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *MachinePoolGetRequest) SendContext(ctx context.Context) (result *MachinePoolGetResponse, err error) {
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
	result = &MachinePoolGetResponse{}
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
	err = readMachinePoolGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// MachinePoolGetResponse is the response for the 'get' method.
type MachinePoolGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *MachinePool
}

// Status returns the response status code.
func (r *MachinePoolGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *MachinePoolGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *MachinePoolGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *MachinePoolGetResponse) Body() *MachinePool {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *MachinePoolGetResponse) GetBody() (value *MachinePool, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// MachinePoolUpdateRequest is the request for the 'update' method.
type MachinePoolUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *MachinePool
}

// Parameter adds a query parameter.
func (r *MachinePoolUpdateRequest) Parameter(name string, value interface{}) *MachinePoolUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *MachinePoolUpdateRequest) Header(name string, value interface{}) *MachinePoolUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *MachinePoolUpdateRequest) Impersonate(user string) *MachinePoolUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *MachinePoolUpdateRequest) Body(value *MachinePool) *MachinePoolUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *MachinePoolUpdateRequest) Send() (result *MachinePoolUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *MachinePoolUpdateRequest) SendContext(ctx context.Context) (result *MachinePoolUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeMachinePoolUpdateRequest(r, buffer)
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
	result = &MachinePoolUpdateResponse{}
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
	err = readMachinePoolUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// MachinePoolUpdateResponse is the response for the 'update' method.
type MachinePoolUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *MachinePool
}

// Status returns the response status code.
func (r *MachinePoolUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *MachinePoolUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *MachinePoolUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *MachinePoolUpdateResponse) Body() *MachinePool {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *MachinePoolUpdateResponse) GetBody() (value *MachinePool, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
