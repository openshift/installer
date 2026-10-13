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

// ExternalAuthClient is the client of the 'external_auth' resource.
//
// Manages a specific external authentication.
type ExternalAuthClient struct {
	transport http.RoundTripper
	path      string
}

// NewExternalAuthClient creates a new client for the 'external_auth'
// resource using the given transport to send the requests and receive the
// responses.
func NewExternalAuthClient(transport http.RoundTripper, path string) *ExternalAuthClient {
	return &ExternalAuthClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the external authentication.
func (c *ExternalAuthClient) Delete() *ExternalAuthDeleteRequest {
	return &ExternalAuthDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of an external authentication.
func (c *ExternalAuthClient) Get() *ExternalAuthGetRequest {
	return &ExternalAuthGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the external authentication.
func (c *ExternalAuthClient) Update() *ExternalAuthUpdateRequest {
	return &ExternalAuthUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ExternalAuthPollRequest is the request for the Poll method.
type ExternalAuthPollRequest struct {
	request    *ExternalAuthGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ExternalAuthPollRequest) Parameter(name string, value interface{}) *ExternalAuthPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ExternalAuthPollRequest) Header(name string, value interface{}) *ExternalAuthPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ExternalAuthPollRequest) Interval(value time.Duration) *ExternalAuthPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ExternalAuthPollRequest) Status(value int) *ExternalAuthPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ExternalAuthPollRequest) Predicate(value func(*ExternalAuthGetResponse) bool) *ExternalAuthPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ExternalAuthGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ExternalAuthPollRequest) StartContext(ctx context.Context) (response *ExternalAuthPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ExternalAuthPollResponse{
			response: result.(*ExternalAuthGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ExternalAuthPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ExternalAuthPollResponse is the response for the Poll method.
type ExternalAuthPollResponse struct {
	response *ExternalAuthGetResponse
}

// Status returns the response status code.
func (r *ExternalAuthPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ExternalAuthPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ExternalAuthPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ExternalAuthPollResponse) Body() *ExternalAuth {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ExternalAuthPollResponse) GetBody() (value *ExternalAuth, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ExternalAuthClient) Poll() *ExternalAuthPollRequest {
	return &ExternalAuthPollRequest{
		request: c.Get(),
	}
}

// ExternalAuthDeleteRequest is the request for the 'delete' method.
type ExternalAuthDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ExternalAuthDeleteRequest) Parameter(name string, value interface{}) *ExternalAuthDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ExternalAuthDeleteRequest) Header(name string, value interface{}) *ExternalAuthDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ExternalAuthDeleteRequest) Impersonate(user string) *ExternalAuthDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ExternalAuthDeleteRequest) Send() (result *ExternalAuthDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ExternalAuthDeleteRequest) SendContext(ctx context.Context) (result *ExternalAuthDeleteResponse, err error) {
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
	result = &ExternalAuthDeleteResponse{}
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

// ExternalAuthDeleteResponse is the response for the 'delete' method.
type ExternalAuthDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ExternalAuthDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ExternalAuthDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ExternalAuthDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ExternalAuthGetRequest is the request for the 'get' method.
type ExternalAuthGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ExternalAuthGetRequest) Parameter(name string, value interface{}) *ExternalAuthGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ExternalAuthGetRequest) Header(name string, value interface{}) *ExternalAuthGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ExternalAuthGetRequest) Impersonate(user string) *ExternalAuthGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ExternalAuthGetRequest) Send() (result *ExternalAuthGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ExternalAuthGetRequest) SendContext(ctx context.Context) (result *ExternalAuthGetResponse, err error) {
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
	result = &ExternalAuthGetResponse{}
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
	err = readExternalAuthGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ExternalAuthGetResponse is the response for the 'get' method.
type ExternalAuthGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ExternalAuth
}

// Status returns the response status code.
func (r *ExternalAuthGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ExternalAuthGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ExternalAuthGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ExternalAuthGetResponse) Body() *ExternalAuth {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ExternalAuthGetResponse) GetBody() (value *ExternalAuth, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ExternalAuthUpdateRequest is the request for the 'update' method.
type ExternalAuthUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ExternalAuth
}

// Parameter adds a query parameter.
func (r *ExternalAuthUpdateRequest) Parameter(name string, value interface{}) *ExternalAuthUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ExternalAuthUpdateRequest) Header(name string, value interface{}) *ExternalAuthUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ExternalAuthUpdateRequest) Impersonate(user string) *ExternalAuthUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ExternalAuthUpdateRequest) Body(value *ExternalAuth) *ExternalAuthUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ExternalAuthUpdateRequest) Send() (result *ExternalAuthUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ExternalAuthUpdateRequest) SendContext(ctx context.Context) (result *ExternalAuthUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeExternalAuthUpdateRequest(r, buffer)
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
	result = &ExternalAuthUpdateResponse{}
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
	err = readExternalAuthUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ExternalAuthUpdateResponse is the response for the 'update' method.
type ExternalAuthUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ExternalAuth
}

// Status returns the response status code.
func (r *ExternalAuthUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ExternalAuthUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ExternalAuthUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ExternalAuthUpdateResponse) Body() *ExternalAuth {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ExternalAuthUpdateResponse) GetBody() (value *ExternalAuth, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
