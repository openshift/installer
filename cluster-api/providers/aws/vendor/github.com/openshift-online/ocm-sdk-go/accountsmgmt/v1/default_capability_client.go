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

// DefaultCapabilityClient is the client of the 'default_capability' resource.
//
// Manages a specific default capability.
type DefaultCapabilityClient struct {
	transport http.RoundTripper
	path      string
}

// NewDefaultCapabilityClient creates a new client for the 'default_capability'
// resource using the given transport to send the requests and receive the
// responses.
func NewDefaultCapabilityClient(transport http.RoundTripper, path string) *DefaultCapabilityClient {
	return &DefaultCapabilityClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
func (c *DefaultCapabilityClient) Delete() *DefaultCapabilityDeleteRequest {
	return &DefaultCapabilityDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the default capability.
func (c *DefaultCapabilityClient) Get() *DefaultCapabilityGetRequest {
	return &DefaultCapabilityGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the default capability.
func (c *DefaultCapabilityClient) Update() *DefaultCapabilityUpdateRequest {
	return &DefaultCapabilityUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// DefaultCapabilityPollRequest is the request for the Poll method.
type DefaultCapabilityPollRequest struct {
	request    *DefaultCapabilityGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *DefaultCapabilityPollRequest) Parameter(name string, value interface{}) *DefaultCapabilityPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *DefaultCapabilityPollRequest) Header(name string, value interface{}) *DefaultCapabilityPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *DefaultCapabilityPollRequest) Interval(value time.Duration) *DefaultCapabilityPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *DefaultCapabilityPollRequest) Status(value int) *DefaultCapabilityPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *DefaultCapabilityPollRequest) Predicate(value func(*DefaultCapabilityGetResponse) bool) *DefaultCapabilityPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*DefaultCapabilityGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *DefaultCapabilityPollRequest) StartContext(ctx context.Context) (response *DefaultCapabilityPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &DefaultCapabilityPollResponse{
			response: result.(*DefaultCapabilityGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *DefaultCapabilityPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// DefaultCapabilityPollResponse is the response for the Poll method.
type DefaultCapabilityPollResponse struct {
	response *DefaultCapabilityGetResponse
}

// Status returns the response status code.
func (r *DefaultCapabilityPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *DefaultCapabilityPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *DefaultCapabilityPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *DefaultCapabilityPollResponse) Body() *DefaultCapability {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *DefaultCapabilityPollResponse) GetBody() (value *DefaultCapability, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *DefaultCapabilityClient) Poll() *DefaultCapabilityPollRequest {
	return &DefaultCapabilityPollRequest{
		request: c.Get(),
	}
}

// DefaultCapabilityDeleteRequest is the request for the 'delete' method.
type DefaultCapabilityDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *DefaultCapabilityDeleteRequest) Parameter(name string, value interface{}) *DefaultCapabilityDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DefaultCapabilityDeleteRequest) Header(name string, value interface{}) *DefaultCapabilityDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DefaultCapabilityDeleteRequest) Impersonate(user string) *DefaultCapabilityDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DefaultCapabilityDeleteRequest) Send() (result *DefaultCapabilityDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DefaultCapabilityDeleteRequest) SendContext(ctx context.Context) (result *DefaultCapabilityDeleteResponse, err error) {
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
	result = &DefaultCapabilityDeleteResponse{}
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

// DefaultCapabilityDeleteResponse is the response for the 'delete' method.
type DefaultCapabilityDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *DefaultCapabilityDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DefaultCapabilityDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DefaultCapabilityDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// DefaultCapabilityGetRequest is the request for the 'get' method.
type DefaultCapabilityGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *DefaultCapabilityGetRequest) Parameter(name string, value interface{}) *DefaultCapabilityGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DefaultCapabilityGetRequest) Header(name string, value interface{}) *DefaultCapabilityGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DefaultCapabilityGetRequest) Impersonate(user string) *DefaultCapabilityGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DefaultCapabilityGetRequest) Send() (result *DefaultCapabilityGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DefaultCapabilityGetRequest) SendContext(ctx context.Context) (result *DefaultCapabilityGetResponse, err error) {
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
	result = &DefaultCapabilityGetResponse{}
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
	err = readDefaultCapabilityGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DefaultCapabilityGetResponse is the response for the 'get' method.
type DefaultCapabilityGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *DefaultCapability
}

// Status returns the response status code.
func (r *DefaultCapabilityGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DefaultCapabilityGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DefaultCapabilityGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *DefaultCapabilityGetResponse) Body() *DefaultCapability {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *DefaultCapabilityGetResponse) GetBody() (value *DefaultCapability, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// DefaultCapabilityUpdateRequest is the request for the 'update' method.
type DefaultCapabilityUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *DefaultCapability
}

// Parameter adds a query parameter.
func (r *DefaultCapabilityUpdateRequest) Parameter(name string, value interface{}) *DefaultCapabilityUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DefaultCapabilityUpdateRequest) Header(name string, value interface{}) *DefaultCapabilityUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DefaultCapabilityUpdateRequest) Impersonate(user string) *DefaultCapabilityUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *DefaultCapabilityUpdateRequest) Body(value *DefaultCapability) *DefaultCapabilityUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DefaultCapabilityUpdateRequest) Send() (result *DefaultCapabilityUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DefaultCapabilityUpdateRequest) SendContext(ctx context.Context) (result *DefaultCapabilityUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeDefaultCapabilityUpdateRequest(r, buffer)
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
	result = &DefaultCapabilityUpdateResponse{}
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
	err = readDefaultCapabilityUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DefaultCapabilityUpdateResponse is the response for the 'update' method.
type DefaultCapabilityUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *DefaultCapability
}

// Status returns the response status code.
func (r *DefaultCapabilityUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DefaultCapabilityUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DefaultCapabilityUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *DefaultCapabilityUpdateResponse) Body() *DefaultCapability {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *DefaultCapabilityUpdateResponse) GetBody() (value *DefaultCapability, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
