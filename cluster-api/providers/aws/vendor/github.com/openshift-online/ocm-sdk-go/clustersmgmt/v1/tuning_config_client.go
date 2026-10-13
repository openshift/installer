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

// TuningConfigClient is the client of the 'tuning_config' resource.
//
// Manages a specific tuning config.
type TuningConfigClient struct {
	transport http.RoundTripper
	path      string
}

// NewTuningConfigClient creates a new client for the 'tuning_config'
// resource using the given transport to send the requests and receive the
// responses.
func NewTuningConfigClient(transport http.RoundTripper, path string) *TuningConfigClient {
	return &TuningConfigClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the tuning config.
func (c *TuningConfigClient) Delete() *TuningConfigDeleteRequest {
	return &TuningConfigDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the tuning config.
func (c *TuningConfigClient) Get() *TuningConfigGetRequest {
	return &TuningConfigGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the tuning config.
func (c *TuningConfigClient) Update() *TuningConfigUpdateRequest {
	return &TuningConfigUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// TuningConfigPollRequest is the request for the Poll method.
type TuningConfigPollRequest struct {
	request    *TuningConfigGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *TuningConfigPollRequest) Parameter(name string, value interface{}) *TuningConfigPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *TuningConfigPollRequest) Header(name string, value interface{}) *TuningConfigPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *TuningConfigPollRequest) Interval(value time.Duration) *TuningConfigPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *TuningConfigPollRequest) Status(value int) *TuningConfigPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *TuningConfigPollRequest) Predicate(value func(*TuningConfigGetResponse) bool) *TuningConfigPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*TuningConfigGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *TuningConfigPollRequest) StartContext(ctx context.Context) (response *TuningConfigPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &TuningConfigPollResponse{
			response: result.(*TuningConfigGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *TuningConfigPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// TuningConfigPollResponse is the response for the Poll method.
type TuningConfigPollResponse struct {
	response *TuningConfigGetResponse
}

// Status returns the response status code.
func (r *TuningConfigPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *TuningConfigPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *TuningConfigPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *TuningConfigPollResponse) Body() *TuningConfig {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *TuningConfigPollResponse) GetBody() (value *TuningConfig, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *TuningConfigClient) Poll() *TuningConfigPollRequest {
	return &TuningConfigPollRequest{
		request: c.Get(),
	}
}

// TuningConfigDeleteRequest is the request for the 'delete' method.
type TuningConfigDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *TuningConfigDeleteRequest) Parameter(name string, value interface{}) *TuningConfigDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *TuningConfigDeleteRequest) Header(name string, value interface{}) *TuningConfigDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *TuningConfigDeleteRequest) Impersonate(user string) *TuningConfigDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *TuningConfigDeleteRequest) Send() (result *TuningConfigDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *TuningConfigDeleteRequest) SendContext(ctx context.Context) (result *TuningConfigDeleteResponse, err error) {
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
	result = &TuningConfigDeleteResponse{}
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

// TuningConfigDeleteResponse is the response for the 'delete' method.
type TuningConfigDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *TuningConfigDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *TuningConfigDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *TuningConfigDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// TuningConfigGetRequest is the request for the 'get' method.
type TuningConfigGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *TuningConfigGetRequest) Parameter(name string, value interface{}) *TuningConfigGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *TuningConfigGetRequest) Header(name string, value interface{}) *TuningConfigGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *TuningConfigGetRequest) Impersonate(user string) *TuningConfigGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *TuningConfigGetRequest) Send() (result *TuningConfigGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *TuningConfigGetRequest) SendContext(ctx context.Context) (result *TuningConfigGetResponse, err error) {
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
	result = &TuningConfigGetResponse{}
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
	err = readTuningConfigGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// TuningConfigGetResponse is the response for the 'get' method.
type TuningConfigGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *TuningConfig
}

// Status returns the response status code.
func (r *TuningConfigGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *TuningConfigGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *TuningConfigGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *TuningConfigGetResponse) Body() *TuningConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *TuningConfigGetResponse) GetBody() (value *TuningConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// TuningConfigUpdateRequest is the request for the 'update' method.
type TuningConfigUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *TuningConfig
}

// Parameter adds a query parameter.
func (r *TuningConfigUpdateRequest) Parameter(name string, value interface{}) *TuningConfigUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *TuningConfigUpdateRequest) Header(name string, value interface{}) *TuningConfigUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *TuningConfigUpdateRequest) Impersonate(user string) *TuningConfigUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *TuningConfigUpdateRequest) Body(value *TuningConfig) *TuningConfigUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *TuningConfigUpdateRequest) Send() (result *TuningConfigUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *TuningConfigUpdateRequest) SendContext(ctx context.Context) (result *TuningConfigUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeTuningConfigUpdateRequest(r, buffer)
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
	result = &TuningConfigUpdateResponse{}
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
	err = readTuningConfigUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// TuningConfigUpdateResponse is the response for the 'update' method.
type TuningConfigUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *TuningConfig
}

// Status returns the response status code.
func (r *TuningConfigUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *TuningConfigUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *TuningConfigUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *TuningConfigUpdateResponse) Body() *TuningConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *TuningConfigUpdateResponse) GetBody() (value *TuningConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
