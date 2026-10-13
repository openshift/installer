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
	"path"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// WifConfigClient is the client of the 'wif_config' resource.
//
// Manages a specific wif_config.
type WifConfigClient struct {
	transport http.RoundTripper
	path      string
}

// NewWifConfigClient creates a new client for the 'wif_config'
// resource using the given transport to send the requests and receive the
// responses.
func NewWifConfigClient(transport http.RoundTripper, path string) *WifConfigClient {
	return &WifConfigClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the wif_config.
func (c *WifConfigClient) Delete() *WifConfigDeleteRequest {
	return &WifConfigDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the WifConfig.
func (c *WifConfigClient) Get() *WifConfigGetRequest {
	return &WifConfigGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the WifConfig.
func (c *WifConfigClient) Update() *WifConfigUpdateRequest {
	return &WifConfigUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Status returns the target 'wif_config_status' resource.
func (c *WifConfigClient) Status() *WifConfigStatusClient {
	return NewWifConfigStatusClient(
		c.transport,
		path.Join(c.path, "status"),
	)
}

// WifConfigPollRequest is the request for the Poll method.
type WifConfigPollRequest struct {
	request    *WifConfigGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *WifConfigPollRequest) Parameter(name string, value interface{}) *WifConfigPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *WifConfigPollRequest) Header(name string, value interface{}) *WifConfigPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *WifConfigPollRequest) Interval(value time.Duration) *WifConfigPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *WifConfigPollRequest) Status(value int) *WifConfigPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *WifConfigPollRequest) Predicate(value func(*WifConfigGetResponse) bool) *WifConfigPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*WifConfigGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *WifConfigPollRequest) StartContext(ctx context.Context) (response *WifConfigPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &WifConfigPollResponse{
			response: result.(*WifConfigGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *WifConfigPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// WifConfigPollResponse is the response for the Poll method.
type WifConfigPollResponse struct {
	response *WifConfigGetResponse
}

// Status returns the response status code.
func (r *WifConfigPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *WifConfigPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *WifConfigPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *WifConfigPollResponse) Body() *WifConfig {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *WifConfigPollResponse) GetBody() (value *WifConfig, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *WifConfigClient) Poll() *WifConfigPollRequest {
	return &WifConfigPollRequest{
		request: c.Get(),
	}
}

// WifConfigDeleteRequest is the request for the 'delete' method.
type WifConfigDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	dryRun    *bool
}

// Parameter adds a query parameter.
func (r *WifConfigDeleteRequest) Parameter(name string, value interface{}) *WifConfigDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *WifConfigDeleteRequest) Header(name string, value interface{}) *WifConfigDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *WifConfigDeleteRequest) Impersonate(user string) *WifConfigDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// DryRun sets the value of the 'dry_run' parameter.
//
// Dry run flag is used to check if the operation can be completed, but won't delete.
func (r *WifConfigDeleteRequest) DryRun(value bool) *WifConfigDeleteRequest {
	r.dryRun = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *WifConfigDeleteRequest) Send() (result *WifConfigDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *WifConfigDeleteRequest) SendContext(ctx context.Context) (result *WifConfigDeleteResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.dryRun != nil {
		helpers.AddValue(&query, "dry_run", *r.dryRun)
	}
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
	result = &WifConfigDeleteResponse{}
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

// WifConfigDeleteResponse is the response for the 'delete' method.
type WifConfigDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *WifConfigDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *WifConfigDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *WifConfigDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// WifConfigGetRequest is the request for the 'get' method.
type WifConfigGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *WifConfigGetRequest) Parameter(name string, value interface{}) *WifConfigGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *WifConfigGetRequest) Header(name string, value interface{}) *WifConfigGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *WifConfigGetRequest) Impersonate(user string) *WifConfigGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *WifConfigGetRequest) Send() (result *WifConfigGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *WifConfigGetRequest) SendContext(ctx context.Context) (result *WifConfigGetResponse, err error) {
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
	result = &WifConfigGetResponse{}
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
	err = readWifConfigGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// WifConfigGetResponse is the response for the 'get' method.
type WifConfigGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *WifConfig
}

// Status returns the response status code.
func (r *WifConfigGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *WifConfigGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *WifConfigGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *WifConfigGetResponse) Body() *WifConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *WifConfigGetResponse) GetBody() (value *WifConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// WifConfigUpdateRequest is the request for the 'update' method.
type WifConfigUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *WifConfig
}

// Parameter adds a query parameter.
func (r *WifConfigUpdateRequest) Parameter(name string, value interface{}) *WifConfigUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *WifConfigUpdateRequest) Header(name string, value interface{}) *WifConfigUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *WifConfigUpdateRequest) Impersonate(user string) *WifConfigUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *WifConfigUpdateRequest) Body(value *WifConfig) *WifConfigUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *WifConfigUpdateRequest) Send() (result *WifConfigUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *WifConfigUpdateRequest) SendContext(ctx context.Context) (result *WifConfigUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeWifConfigUpdateRequest(r, buffer)
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
	result = &WifConfigUpdateResponse{}
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
	err = readWifConfigUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// WifConfigUpdateResponse is the response for the 'update' method.
type WifConfigUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *WifConfig
}

// Status returns the response status code.
func (r *WifConfigUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *WifConfigUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *WifConfigUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *WifConfigUpdateResponse) Body() *WifConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *WifConfigUpdateResponse) GetBody() (value *WifConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
