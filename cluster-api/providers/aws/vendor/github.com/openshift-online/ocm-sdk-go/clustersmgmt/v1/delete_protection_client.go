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

// DeleteProtectionClient is the client of the 'delete_protection' resource.
//
// Manages delete protection specific parts for a specific cluster.
type DeleteProtectionClient struct {
	transport http.RoundTripper
	path      string
}

// NewDeleteProtectionClient creates a new client for the 'delete_protection'
// resource using the given transport to send the requests and receive the
// responses.
func NewDeleteProtectionClient(transport http.RoundTripper, path string) *DeleteProtectionClient {
	return &DeleteProtectionClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
func (c *DeleteProtectionClient) Get() *DeleteProtectionGetRequest {
	return &DeleteProtectionGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
func (c *DeleteProtectionClient) Update() *DeleteProtectionUpdateRequest {
	return &DeleteProtectionUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// DeleteProtectionPollRequest is the request for the Poll method.
type DeleteProtectionPollRequest struct {
	request    *DeleteProtectionGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *DeleteProtectionPollRequest) Parameter(name string, value interface{}) *DeleteProtectionPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *DeleteProtectionPollRequest) Header(name string, value interface{}) *DeleteProtectionPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *DeleteProtectionPollRequest) Interval(value time.Duration) *DeleteProtectionPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *DeleteProtectionPollRequest) Status(value int) *DeleteProtectionPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *DeleteProtectionPollRequest) Predicate(value func(*DeleteProtectionGetResponse) bool) *DeleteProtectionPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*DeleteProtectionGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *DeleteProtectionPollRequest) StartContext(ctx context.Context) (response *DeleteProtectionPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &DeleteProtectionPollResponse{
			response: result.(*DeleteProtectionGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *DeleteProtectionPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// DeleteProtectionPollResponse is the response for the Poll method.
type DeleteProtectionPollResponse struct {
	response *DeleteProtectionGetResponse
}

// Status returns the response status code.
func (r *DeleteProtectionPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *DeleteProtectionPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *DeleteProtectionPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *DeleteProtectionPollResponse) Body() *DeleteProtection {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *DeleteProtectionPollResponse) GetBody() (value *DeleteProtection, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *DeleteProtectionClient) Poll() *DeleteProtectionPollRequest {
	return &DeleteProtectionPollRequest{
		request: c.Get(),
	}
}

// DeleteProtectionGetRequest is the request for the 'get' method.
type DeleteProtectionGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *DeleteProtectionGetRequest) Parameter(name string, value interface{}) *DeleteProtectionGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DeleteProtectionGetRequest) Header(name string, value interface{}) *DeleteProtectionGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DeleteProtectionGetRequest) Impersonate(user string) *DeleteProtectionGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DeleteProtectionGetRequest) Send() (result *DeleteProtectionGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DeleteProtectionGetRequest) SendContext(ctx context.Context) (result *DeleteProtectionGetResponse, err error) {
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
	result = &DeleteProtectionGetResponse{}
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
	err = readDeleteProtectionGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DeleteProtectionGetResponse is the response for the 'get' method.
type DeleteProtectionGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *DeleteProtection
}

// Status returns the response status code.
func (r *DeleteProtectionGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DeleteProtectionGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DeleteProtectionGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *DeleteProtectionGetResponse) Body() *DeleteProtection {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *DeleteProtectionGetResponse) GetBody() (value *DeleteProtection, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// DeleteProtectionUpdateRequest is the request for the 'update' method.
type DeleteProtectionUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *DeleteProtection
}

// Parameter adds a query parameter.
func (r *DeleteProtectionUpdateRequest) Parameter(name string, value interface{}) *DeleteProtectionUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DeleteProtectionUpdateRequest) Header(name string, value interface{}) *DeleteProtectionUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DeleteProtectionUpdateRequest) Impersonate(user string) *DeleteProtectionUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *DeleteProtectionUpdateRequest) Body(value *DeleteProtection) *DeleteProtectionUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DeleteProtectionUpdateRequest) Send() (result *DeleteProtectionUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DeleteProtectionUpdateRequest) SendContext(ctx context.Context) (result *DeleteProtectionUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeDeleteProtectionUpdateRequest(r, buffer)
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
	result = &DeleteProtectionUpdateResponse{}
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
	err = readDeleteProtectionUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DeleteProtectionUpdateResponse is the response for the 'update' method.
type DeleteProtectionUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *DeleteProtection
}

// Status returns the response status code.
func (r *DeleteProtectionUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DeleteProtectionUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DeleteProtectionUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *DeleteProtectionUpdateResponse) Body() *DeleteProtection {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *DeleteProtectionUpdateResponse) GetBody() (value *DeleteProtection, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
