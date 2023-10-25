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

// ApplicationDependencyClient is the client of the 'application_dependency' resource.
//
// Provides detailed information about the specified application dependency.
type ApplicationDependencyClient struct {
	transport http.RoundTripper
	path      string
}

// NewApplicationDependencyClient creates a new client for the 'application_dependency'
// resource using the given transport to send the requests and receive the
// responses.
func NewApplicationDependencyClient(transport http.RoundTripper, path string) *ApplicationDependencyClient {
	return &ApplicationDependencyClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
func (c *ApplicationDependencyClient) Delete() *ApplicationDependencyDeleteRequest {
	return &ApplicationDependencyDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
func (c *ApplicationDependencyClient) Get() *ApplicationDependencyGetRequest {
	return &ApplicationDependencyGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
func (c *ApplicationDependencyClient) Update() *ApplicationDependencyUpdateRequest {
	return &ApplicationDependencyUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ApplicationDependencyPollRequest is the request for the Poll method.
type ApplicationDependencyPollRequest struct {
	request    *ApplicationDependencyGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ApplicationDependencyPollRequest) Parameter(name string, value interface{}) *ApplicationDependencyPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ApplicationDependencyPollRequest) Header(name string, value interface{}) *ApplicationDependencyPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ApplicationDependencyPollRequest) Interval(value time.Duration) *ApplicationDependencyPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ApplicationDependencyPollRequest) Status(value int) *ApplicationDependencyPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ApplicationDependencyPollRequest) Predicate(value func(*ApplicationDependencyGetResponse) bool) *ApplicationDependencyPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ApplicationDependencyGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ApplicationDependencyPollRequest) StartContext(ctx context.Context) (response *ApplicationDependencyPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ApplicationDependencyPollResponse{
			response: result.(*ApplicationDependencyGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ApplicationDependencyPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ApplicationDependencyPollResponse is the response for the Poll method.
type ApplicationDependencyPollResponse struct {
	response *ApplicationDependencyGetResponse
}

// Status returns the response status code.
func (r *ApplicationDependencyPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ApplicationDependencyPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ApplicationDependencyPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ApplicationDependencyPollResponse) Body() *ApplicationDependency {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationDependencyPollResponse) GetBody() (value *ApplicationDependency, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ApplicationDependencyClient) Poll() *ApplicationDependencyPollRequest {
	return &ApplicationDependencyPollRequest{
		request: c.Get(),
	}
}

// ApplicationDependencyDeleteRequest is the request for the 'delete' method.
type ApplicationDependencyDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ApplicationDependencyDeleteRequest) Parameter(name string, value interface{}) *ApplicationDependencyDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ApplicationDependencyDeleteRequest) Header(name string, value interface{}) *ApplicationDependencyDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ApplicationDependencyDeleteRequest) Impersonate(user string) *ApplicationDependencyDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ApplicationDependencyDeleteRequest) Send() (result *ApplicationDependencyDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ApplicationDependencyDeleteRequest) SendContext(ctx context.Context) (result *ApplicationDependencyDeleteResponse, err error) {
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
	result = &ApplicationDependencyDeleteResponse{}
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

// ApplicationDependencyDeleteResponse is the response for the 'delete' method.
type ApplicationDependencyDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ApplicationDependencyDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ApplicationDependencyDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ApplicationDependencyDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ApplicationDependencyGetRequest is the request for the 'get' method.
type ApplicationDependencyGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ApplicationDependencyGetRequest) Parameter(name string, value interface{}) *ApplicationDependencyGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ApplicationDependencyGetRequest) Header(name string, value interface{}) *ApplicationDependencyGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ApplicationDependencyGetRequest) Impersonate(user string) *ApplicationDependencyGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ApplicationDependencyGetRequest) Send() (result *ApplicationDependencyGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ApplicationDependencyGetRequest) SendContext(ctx context.Context) (result *ApplicationDependencyGetResponse, err error) {
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
	result = &ApplicationDependencyGetResponse{}
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
	err = readApplicationDependencyGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ApplicationDependencyGetResponse is the response for the 'get' method.
type ApplicationDependencyGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ApplicationDependency
}

// Status returns the response status code.
func (r *ApplicationDependencyGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ApplicationDependencyGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ApplicationDependencyGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ApplicationDependencyGetResponse) Body() *ApplicationDependency {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationDependencyGetResponse) GetBody() (value *ApplicationDependency, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ApplicationDependencyUpdateRequest is the request for the 'update' method.
type ApplicationDependencyUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ApplicationDependency
}

// Parameter adds a query parameter.
func (r *ApplicationDependencyUpdateRequest) Parameter(name string, value interface{}) *ApplicationDependencyUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ApplicationDependencyUpdateRequest) Header(name string, value interface{}) *ApplicationDependencyUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ApplicationDependencyUpdateRequest) Impersonate(user string) *ApplicationDependencyUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ApplicationDependencyUpdateRequest) Body(value *ApplicationDependency) *ApplicationDependencyUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ApplicationDependencyUpdateRequest) Send() (result *ApplicationDependencyUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ApplicationDependencyUpdateRequest) SendContext(ctx context.Context) (result *ApplicationDependencyUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeApplicationDependencyUpdateRequest(r, buffer)
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
	result = &ApplicationDependencyUpdateResponse{}
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
	err = readApplicationDependencyUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ApplicationDependencyUpdateResponse is the response for the 'update' method.
type ApplicationDependencyUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ApplicationDependency
}

// Status returns the response status code.
func (r *ApplicationDependencyUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ApplicationDependencyUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ApplicationDependencyUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ApplicationDependencyUpdateResponse) Body() *ApplicationDependency {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationDependencyUpdateResponse) GetBody() (value *ApplicationDependency, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
