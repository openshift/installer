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

// StatusUpdateClient is the client of the 'status_update' resource.
//
// Provides detailed information about the specified status.
type StatusUpdateClient struct {
	transport http.RoundTripper
	path      string
}

// NewStatusUpdateClient creates a new client for the 'status_update'
// resource using the given transport to send the requests and receive the
// responses.
func NewStatusUpdateClient(transport http.RoundTripper, path string) *StatusUpdateClient {
	return &StatusUpdateClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
func (c *StatusUpdateClient) Delete() *StatusUpdateDeleteRequest {
	return &StatusUpdateDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
func (c *StatusUpdateClient) Get() *StatusUpdateGetRequest {
	return &StatusUpdateGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
func (c *StatusUpdateClient) Update() *StatusUpdateUpdateRequest {
	return &StatusUpdateUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// StatusUpdatePollRequest is the request for the Poll method.
type StatusUpdatePollRequest struct {
	request    *StatusUpdateGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *StatusUpdatePollRequest) Parameter(name string, value interface{}) *StatusUpdatePollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *StatusUpdatePollRequest) Header(name string, value interface{}) *StatusUpdatePollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *StatusUpdatePollRequest) Interval(value time.Duration) *StatusUpdatePollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *StatusUpdatePollRequest) Status(value int) *StatusUpdatePollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *StatusUpdatePollRequest) Predicate(value func(*StatusUpdateGetResponse) bool) *StatusUpdatePollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*StatusUpdateGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *StatusUpdatePollRequest) StartContext(ctx context.Context) (response *StatusUpdatePollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &StatusUpdatePollResponse{
			response: result.(*StatusUpdateGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *StatusUpdatePollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// StatusUpdatePollResponse is the response for the Poll method.
type StatusUpdatePollResponse struct {
	response *StatusUpdateGetResponse
}

// Status returns the response status code.
func (r *StatusUpdatePollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *StatusUpdatePollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *StatusUpdatePollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *StatusUpdatePollResponse) Body() *Status {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *StatusUpdatePollResponse) GetBody() (value *Status, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *StatusUpdateClient) Poll() *StatusUpdatePollRequest {
	return &StatusUpdatePollRequest{
		request: c.Get(),
	}
}

// StatusUpdateDeleteRequest is the request for the 'delete' method.
type StatusUpdateDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *StatusUpdateDeleteRequest) Parameter(name string, value interface{}) *StatusUpdateDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *StatusUpdateDeleteRequest) Header(name string, value interface{}) *StatusUpdateDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *StatusUpdateDeleteRequest) Impersonate(user string) *StatusUpdateDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *StatusUpdateDeleteRequest) Send() (result *StatusUpdateDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *StatusUpdateDeleteRequest) SendContext(ctx context.Context) (result *StatusUpdateDeleteResponse, err error) {
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
	result = &StatusUpdateDeleteResponse{}
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

// StatusUpdateDeleteResponse is the response for the 'delete' method.
type StatusUpdateDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *StatusUpdateDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *StatusUpdateDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *StatusUpdateDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// StatusUpdateGetRequest is the request for the 'get' method.
type StatusUpdateGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *StatusUpdateGetRequest) Parameter(name string, value interface{}) *StatusUpdateGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *StatusUpdateGetRequest) Header(name string, value interface{}) *StatusUpdateGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *StatusUpdateGetRequest) Impersonate(user string) *StatusUpdateGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *StatusUpdateGetRequest) Send() (result *StatusUpdateGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *StatusUpdateGetRequest) SendContext(ctx context.Context) (result *StatusUpdateGetResponse, err error) {
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
	result = &StatusUpdateGetResponse{}
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
	err = readStatusUpdateGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// StatusUpdateGetResponse is the response for the 'get' method.
type StatusUpdateGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Status
}

// Status returns the response status code.
func (r *StatusUpdateGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *StatusUpdateGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *StatusUpdateGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *StatusUpdateGetResponse) Body() *Status {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *StatusUpdateGetResponse) GetBody() (value *Status, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// StatusUpdateUpdateRequest is the request for the 'update' method.
type StatusUpdateUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Status
}

// Parameter adds a query parameter.
func (r *StatusUpdateUpdateRequest) Parameter(name string, value interface{}) *StatusUpdateUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *StatusUpdateUpdateRequest) Header(name string, value interface{}) *StatusUpdateUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *StatusUpdateUpdateRequest) Impersonate(user string) *StatusUpdateUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *StatusUpdateUpdateRequest) Body(value *Status) *StatusUpdateUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *StatusUpdateUpdateRequest) Send() (result *StatusUpdateUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *StatusUpdateUpdateRequest) SendContext(ctx context.Context) (result *StatusUpdateUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeStatusUpdateUpdateRequest(r, buffer)
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
	result = &StatusUpdateUpdateResponse{}
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
	err = readStatusUpdateUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// StatusUpdateUpdateResponse is the response for the 'update' method.
type StatusUpdateUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Status
}

// Status returns the response status code.
func (r *StatusUpdateUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *StatusUpdateUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *StatusUpdateUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *StatusUpdateUpdateResponse) Body() *Status {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *StatusUpdateUpdateResponse) GetBody() (value *Status, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
