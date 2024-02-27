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

// HypershiftClient is the client of the 'hypershift' resource.
//
// Manages a specific Hypershift cluster.
type HypershiftClient struct {
	transport http.RoundTripper
	path      string
}

// NewHypershiftClient creates a new client for the 'hypershift'
// resource using the given transport to send the requests and receive the
// responses.
func NewHypershiftClient(transport http.RoundTripper, path string) *HypershiftClient {
	return &HypershiftClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the Hypershift details for a single cluster.
func (c *HypershiftClient) Get() *HypershiftGetRequest {
	return &HypershiftGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the Hypershift details for a single cluster.
func (c *HypershiftClient) Update() *HypershiftUpdateRequest {
	return &HypershiftUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// HypershiftPollRequest is the request for the Poll method.
type HypershiftPollRequest struct {
	request    *HypershiftGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *HypershiftPollRequest) Parameter(name string, value interface{}) *HypershiftPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *HypershiftPollRequest) Header(name string, value interface{}) *HypershiftPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *HypershiftPollRequest) Interval(value time.Duration) *HypershiftPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *HypershiftPollRequest) Status(value int) *HypershiftPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *HypershiftPollRequest) Predicate(value func(*HypershiftGetResponse) bool) *HypershiftPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*HypershiftGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *HypershiftPollRequest) StartContext(ctx context.Context) (response *HypershiftPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &HypershiftPollResponse{
			response: result.(*HypershiftGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *HypershiftPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// HypershiftPollResponse is the response for the Poll method.
type HypershiftPollResponse struct {
	response *HypershiftGetResponse
}

// Status returns the response status code.
func (r *HypershiftPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *HypershiftPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *HypershiftPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *HypershiftPollResponse) Body() *HypershiftConfig {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *HypershiftPollResponse) GetBody() (value *HypershiftConfig, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *HypershiftClient) Poll() *HypershiftPollRequest {
	return &HypershiftPollRequest{
		request: c.Get(),
	}
}

// HypershiftGetRequest is the request for the 'get' method.
type HypershiftGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *HypershiftGetRequest) Parameter(name string, value interface{}) *HypershiftGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *HypershiftGetRequest) Header(name string, value interface{}) *HypershiftGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *HypershiftGetRequest) Impersonate(user string) *HypershiftGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *HypershiftGetRequest) Send() (result *HypershiftGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *HypershiftGetRequest) SendContext(ctx context.Context) (result *HypershiftGetResponse, err error) {
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
	result = &HypershiftGetResponse{}
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
	err = readHypershiftGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// HypershiftGetResponse is the response for the 'get' method.
type HypershiftGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *HypershiftConfig
}

// Status returns the response status code.
func (r *HypershiftGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *HypershiftGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *HypershiftGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *HypershiftGetResponse) Body() *HypershiftConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *HypershiftGetResponse) GetBody() (value *HypershiftConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// HypershiftUpdateRequest is the request for the 'update' method.
type HypershiftUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *HypershiftConfig
}

// Parameter adds a query parameter.
func (r *HypershiftUpdateRequest) Parameter(name string, value interface{}) *HypershiftUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *HypershiftUpdateRequest) Header(name string, value interface{}) *HypershiftUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *HypershiftUpdateRequest) Impersonate(user string) *HypershiftUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *HypershiftUpdateRequest) Body(value *HypershiftConfig) *HypershiftUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *HypershiftUpdateRequest) Send() (result *HypershiftUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *HypershiftUpdateRequest) SendContext(ctx context.Context) (result *HypershiftUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeHypershiftUpdateRequest(r, buffer)
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
	result = &HypershiftUpdateResponse{}
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
	err = readHypershiftUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// HypershiftUpdateResponse is the response for the 'update' method.
type HypershiftUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *HypershiftConfig
}

// Status returns the response status code.
func (r *HypershiftUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *HypershiftUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *HypershiftUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *HypershiftUpdateResponse) Body() *HypershiftConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *HypershiftUpdateResponse) GetBody() (value *HypershiftConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
