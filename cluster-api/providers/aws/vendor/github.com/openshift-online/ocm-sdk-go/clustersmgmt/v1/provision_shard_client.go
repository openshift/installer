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

// ProvisionShardClient is the client of the 'provision_shard' resource.
//
// Manages provision shard.
type ProvisionShardClient struct {
	transport http.RoundTripper
	path      string
}

// NewProvisionShardClient creates a new client for the 'provision_shard'
// resource using the given transport to send the requests and receive the
// responses.
func NewProvisionShardClient(transport http.RoundTripper, path string) *ProvisionShardClient {
	return &ProvisionShardClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Delete the provision shard.
func (c *ProvisionShardClient) Delete() *ProvisionShardDeleteRequest {
	return &ProvisionShardDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the provision shard.
func (c *ProvisionShardClient) Get() *ProvisionShardGetRequest {
	return &ProvisionShardGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the details of the provision shard.
func (c *ProvisionShardClient) Update() *ProvisionShardUpdateRequest {
	return &ProvisionShardUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ProvisionShardPollRequest is the request for the Poll method.
type ProvisionShardPollRequest struct {
	request    *ProvisionShardGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ProvisionShardPollRequest) Parameter(name string, value interface{}) *ProvisionShardPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ProvisionShardPollRequest) Header(name string, value interface{}) *ProvisionShardPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ProvisionShardPollRequest) Interval(value time.Duration) *ProvisionShardPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ProvisionShardPollRequest) Status(value int) *ProvisionShardPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ProvisionShardPollRequest) Predicate(value func(*ProvisionShardGetResponse) bool) *ProvisionShardPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ProvisionShardGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ProvisionShardPollRequest) StartContext(ctx context.Context) (response *ProvisionShardPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ProvisionShardPollResponse{
			response: result.(*ProvisionShardGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ProvisionShardPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ProvisionShardPollResponse is the response for the Poll method.
type ProvisionShardPollResponse struct {
	response *ProvisionShardGetResponse
}

// Status returns the response status code.
func (r *ProvisionShardPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ProvisionShardPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ProvisionShardPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ProvisionShardPollResponse) Body() *ProvisionShard {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ProvisionShardPollResponse) GetBody() (value *ProvisionShard, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ProvisionShardClient) Poll() *ProvisionShardPollRequest {
	return &ProvisionShardPollRequest{
		request: c.Get(),
	}
}

// ProvisionShardDeleteRequest is the request for the 'delete' method.
type ProvisionShardDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ProvisionShardDeleteRequest) Parameter(name string, value interface{}) *ProvisionShardDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ProvisionShardDeleteRequest) Header(name string, value interface{}) *ProvisionShardDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ProvisionShardDeleteRequest) Impersonate(user string) *ProvisionShardDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ProvisionShardDeleteRequest) Send() (result *ProvisionShardDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ProvisionShardDeleteRequest) SendContext(ctx context.Context) (result *ProvisionShardDeleteResponse, err error) {
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
	result = &ProvisionShardDeleteResponse{}
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

// ProvisionShardDeleteResponse is the response for the 'delete' method.
type ProvisionShardDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ProvisionShardDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ProvisionShardDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ProvisionShardDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ProvisionShardGetRequest is the request for the 'get' method.
type ProvisionShardGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ProvisionShardGetRequest) Parameter(name string, value interface{}) *ProvisionShardGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ProvisionShardGetRequest) Header(name string, value interface{}) *ProvisionShardGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ProvisionShardGetRequest) Impersonate(user string) *ProvisionShardGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ProvisionShardGetRequest) Send() (result *ProvisionShardGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ProvisionShardGetRequest) SendContext(ctx context.Context) (result *ProvisionShardGetResponse, err error) {
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
	result = &ProvisionShardGetResponse{}
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
	err = readProvisionShardGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ProvisionShardGetResponse is the response for the 'get' method.
type ProvisionShardGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ProvisionShard
}

// Status returns the response status code.
func (r *ProvisionShardGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ProvisionShardGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ProvisionShardGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ProvisionShardGetResponse) Body() *ProvisionShard {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ProvisionShardGetResponse) GetBody() (value *ProvisionShard, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ProvisionShardUpdateRequest is the request for the 'update' method.
type ProvisionShardUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ProvisionShard
}

// Parameter adds a query parameter.
func (r *ProvisionShardUpdateRequest) Parameter(name string, value interface{}) *ProvisionShardUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ProvisionShardUpdateRequest) Header(name string, value interface{}) *ProvisionShardUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ProvisionShardUpdateRequest) Impersonate(user string) *ProvisionShardUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ProvisionShardUpdateRequest) Body(value *ProvisionShard) *ProvisionShardUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ProvisionShardUpdateRequest) Send() (result *ProvisionShardUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ProvisionShardUpdateRequest) SendContext(ctx context.Context) (result *ProvisionShardUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeProvisionShardUpdateRequest(r, buffer)
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
	result = &ProvisionShardUpdateResponse{}
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
	err = readProvisionShardUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ProvisionShardUpdateResponse is the response for the 'update' method.
type ProvisionShardUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ProvisionShard
}

// Status returns the response status code.
func (r *ProvisionShardUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ProvisionShardUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ProvisionShardUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ProvisionShardUpdateResponse) Body() *ProvisionShard {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ProvisionShardUpdateResponse) GetBody() (value *ProvisionShard, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
