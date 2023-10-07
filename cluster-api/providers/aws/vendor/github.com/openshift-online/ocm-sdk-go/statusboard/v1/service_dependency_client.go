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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

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

// ServiceDependencyClient is the client of the 'service_dependency' resource.
//
// Provides detailed information about the specified application dependency.
type ServiceDependencyClient struct {
	transport http.RoundTripper
	path      string
}

// NewServiceDependencyClient creates a new client for the 'service_dependency'
// resource using the given transport to send the requests and receive the
// responses.
func NewServiceDependencyClient(transport http.RoundTripper, path string) *ServiceDependencyClient {
	return &ServiceDependencyClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
func (c *ServiceDependencyClient) Delete() *ServiceDependencyDeleteRequest {
	return &ServiceDependencyDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
func (c *ServiceDependencyClient) Get() *ServiceDependencyGetRequest {
	return &ServiceDependencyGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
func (c *ServiceDependencyClient) Update() *ServiceDependencyUpdateRequest {
	return &ServiceDependencyUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ServiceDependencyPollRequest is the request for the Poll method.
type ServiceDependencyPollRequest struct {
	request    *ServiceDependencyGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ServiceDependencyPollRequest) Parameter(name string, value interface{}) *ServiceDependencyPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ServiceDependencyPollRequest) Header(name string, value interface{}) *ServiceDependencyPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ServiceDependencyPollRequest) Interval(value time.Duration) *ServiceDependencyPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ServiceDependencyPollRequest) Status(value int) *ServiceDependencyPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ServiceDependencyPollRequest) Predicate(value func(*ServiceDependencyGetResponse) bool) *ServiceDependencyPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ServiceDependencyGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ServiceDependencyPollRequest) StartContext(ctx context.Context) (response *ServiceDependencyPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ServiceDependencyPollResponse{
			response: result.(*ServiceDependencyGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ServiceDependencyPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ServiceDependencyPollResponse is the response for the Poll method.
type ServiceDependencyPollResponse struct {
	response *ServiceDependencyGetResponse
}

// Status returns the response status code.
func (r *ServiceDependencyPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ServiceDependencyPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ServiceDependencyPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ServiceDependencyPollResponse) Body() *ServiceDependency {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceDependencyPollResponse) GetBody() (value *ServiceDependency, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ServiceDependencyClient) Poll() *ServiceDependencyPollRequest {
	return &ServiceDependencyPollRequest{
		request: c.Get(),
	}
}

// ServiceDependencyDeleteRequest is the request for the 'delete' method.
type ServiceDependencyDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ServiceDependencyDeleteRequest) Parameter(name string, value interface{}) *ServiceDependencyDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServiceDependencyDeleteRequest) Header(name string, value interface{}) *ServiceDependencyDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServiceDependencyDeleteRequest) Impersonate(user string) *ServiceDependencyDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServiceDependencyDeleteRequest) Send() (result *ServiceDependencyDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServiceDependencyDeleteRequest) SendContext(ctx context.Context) (result *ServiceDependencyDeleteResponse, err error) {
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
	result = &ServiceDependencyDeleteResponse{}
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

// ServiceDependencyDeleteResponse is the response for the 'delete' method.
type ServiceDependencyDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ServiceDependencyDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServiceDependencyDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServiceDependencyDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ServiceDependencyGetRequest is the request for the 'get' method.
type ServiceDependencyGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ServiceDependencyGetRequest) Parameter(name string, value interface{}) *ServiceDependencyGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServiceDependencyGetRequest) Header(name string, value interface{}) *ServiceDependencyGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServiceDependencyGetRequest) Impersonate(user string) *ServiceDependencyGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServiceDependencyGetRequest) Send() (result *ServiceDependencyGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServiceDependencyGetRequest) SendContext(ctx context.Context) (result *ServiceDependencyGetResponse, err error) {
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
	result = &ServiceDependencyGetResponse{}
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
	err = readServiceDependencyGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ServiceDependencyGetResponse is the response for the 'get' method.
type ServiceDependencyGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ServiceDependency
}

// Status returns the response status code.
func (r *ServiceDependencyGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServiceDependencyGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServiceDependencyGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ServiceDependencyGetResponse) Body() *ServiceDependency {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceDependencyGetResponse) GetBody() (value *ServiceDependency, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ServiceDependencyUpdateRequest is the request for the 'update' method.
type ServiceDependencyUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ServiceDependency
}

// Parameter adds a query parameter.
func (r *ServiceDependencyUpdateRequest) Parameter(name string, value interface{}) *ServiceDependencyUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServiceDependencyUpdateRequest) Header(name string, value interface{}) *ServiceDependencyUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServiceDependencyUpdateRequest) Impersonate(user string) *ServiceDependencyUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ServiceDependencyUpdateRequest) Body(value *ServiceDependency) *ServiceDependencyUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServiceDependencyUpdateRequest) Send() (result *ServiceDependencyUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServiceDependencyUpdateRequest) SendContext(ctx context.Context) (result *ServiceDependencyUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeServiceDependencyUpdateRequest(r, buffer)
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
	result = &ServiceDependencyUpdateResponse{}
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
	err = readServiceDependencyUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ServiceDependencyUpdateResponse is the response for the 'update' method.
type ServiceDependencyUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ServiceDependency
}

// Status returns the response status code.
func (r *ServiceDependencyUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServiceDependencyUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServiceDependencyUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ServiceDependencyUpdateResponse) Body() *ServiceDependency {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceDependencyUpdateResponse) GetBody() (value *ServiceDependency, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
