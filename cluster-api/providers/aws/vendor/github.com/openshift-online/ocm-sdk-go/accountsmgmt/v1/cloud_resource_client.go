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

// CloudResourceClient is the client of the 'cloud_resource' resource.
//
// Manages a specific cloud resource.
type CloudResourceClient struct {
	transport http.RoundTripper
	path      string
}

// NewCloudResourceClient creates a new client for the 'cloud_resource'
// resource using the given transport to send the requests and receive the
// responses.
func NewCloudResourceClient(transport http.RoundTripper, path string) *CloudResourceClient {
	return &CloudResourceClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the cloud resource.
func (c *CloudResourceClient) Delete() *CloudResourceDeleteRequest {
	return &CloudResourceDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the cloud resource.
func (c *CloudResourceClient) Get() *CloudResourceGetRequest {
	return &CloudResourceGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the cloud resource.
func (c *CloudResourceClient) Update() *CloudResourceUpdateRequest {
	return &CloudResourceUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// CloudResourcePollRequest is the request for the Poll method.
type CloudResourcePollRequest struct {
	request    *CloudResourceGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *CloudResourcePollRequest) Parameter(name string, value interface{}) *CloudResourcePollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *CloudResourcePollRequest) Header(name string, value interface{}) *CloudResourcePollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *CloudResourcePollRequest) Interval(value time.Duration) *CloudResourcePollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *CloudResourcePollRequest) Status(value int) *CloudResourcePollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *CloudResourcePollRequest) Predicate(value func(*CloudResourceGetResponse) bool) *CloudResourcePollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*CloudResourceGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *CloudResourcePollRequest) StartContext(ctx context.Context) (response *CloudResourcePollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &CloudResourcePollResponse{
			response: result.(*CloudResourceGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *CloudResourcePollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// CloudResourcePollResponse is the response for the Poll method.
type CloudResourcePollResponse struct {
	response *CloudResourceGetResponse
}

// Status returns the response status code.
func (r *CloudResourcePollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *CloudResourcePollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *CloudResourcePollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *CloudResourcePollResponse) Body() *CloudResource {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *CloudResourcePollResponse) GetBody() (value *CloudResource, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *CloudResourceClient) Poll() *CloudResourcePollRequest {
	return &CloudResourcePollRequest{
		request: c.Get(),
	}
}

// CloudResourceDeleteRequest is the request for the 'delete' method.
type CloudResourceDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *CloudResourceDeleteRequest) Parameter(name string, value interface{}) *CloudResourceDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *CloudResourceDeleteRequest) Header(name string, value interface{}) *CloudResourceDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *CloudResourceDeleteRequest) Impersonate(user string) *CloudResourceDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *CloudResourceDeleteRequest) Send() (result *CloudResourceDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *CloudResourceDeleteRequest) SendContext(ctx context.Context) (result *CloudResourceDeleteResponse, err error) {
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
	result = &CloudResourceDeleteResponse{}
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

// CloudResourceDeleteResponse is the response for the 'delete' method.
type CloudResourceDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *CloudResourceDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *CloudResourceDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *CloudResourceDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// CloudResourceGetRequest is the request for the 'get' method.
type CloudResourceGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *CloudResourceGetRequest) Parameter(name string, value interface{}) *CloudResourceGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *CloudResourceGetRequest) Header(name string, value interface{}) *CloudResourceGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *CloudResourceGetRequest) Impersonate(user string) *CloudResourceGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *CloudResourceGetRequest) Send() (result *CloudResourceGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *CloudResourceGetRequest) SendContext(ctx context.Context) (result *CloudResourceGetResponse, err error) {
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
	result = &CloudResourceGetResponse{}
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
	err = readCloudResourceGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// CloudResourceGetResponse is the response for the 'get' method.
type CloudResourceGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *CloudResource
}

// Status returns the response status code.
func (r *CloudResourceGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *CloudResourceGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *CloudResourceGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *CloudResourceGetResponse) Body() *CloudResource {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *CloudResourceGetResponse) GetBody() (value *CloudResource, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// CloudResourceUpdateRequest is the request for the 'update' method.
type CloudResourceUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *CloudResource
}

// Parameter adds a query parameter.
func (r *CloudResourceUpdateRequest) Parameter(name string, value interface{}) *CloudResourceUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *CloudResourceUpdateRequest) Header(name string, value interface{}) *CloudResourceUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *CloudResourceUpdateRequest) Impersonate(user string) *CloudResourceUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *CloudResourceUpdateRequest) Body(value *CloudResource) *CloudResourceUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *CloudResourceUpdateRequest) Send() (result *CloudResourceUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *CloudResourceUpdateRequest) SendContext(ctx context.Context) (result *CloudResourceUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeCloudResourceUpdateRequest(r, buffer)
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
	result = &CloudResourceUpdateResponse{}
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
	err = readCloudResourceUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// CloudResourceUpdateResponse is the response for the 'update' method.
type CloudResourceUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *CloudResource
}

// Status returns the response status code.
func (r *CloudResourceUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *CloudResourceUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *CloudResourceUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *CloudResourceUpdateResponse) Body() *CloudResource {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *CloudResourceUpdateResponse) GetBody() (value *CloudResource, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
