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

// AddOnClient is the client of the 'add_on' resource.
//
// Manages a specific add-on.
type AddOnClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddOnClient creates a new client for the 'add_on'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddOnClient(transport http.RoundTripper, path string) *AddOnClient {
	return &AddOnClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the add-on.
func (c *AddOnClient) Delete() *AddOnDeleteRequest {
	return &AddOnDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the add-on.
func (c *AddOnClient) Get() *AddOnGetRequest {
	return &AddOnGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the add-on.
func (c *AddOnClient) Update() *AddOnUpdateRequest {
	return &AddOnUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Versions returns the target 'add_on_versions' resource.
//
// Reference to the resource that manages the collection of addon versions.
func (c *AddOnClient) Versions() *AddOnVersionsClient {
	return NewAddOnVersionsClient(
		c.transport,
		path.Join(c.path, "versions"),
	)
}

// AddOnPollRequest is the request for the Poll method.
type AddOnPollRequest struct {
	request    *AddOnGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AddOnPollRequest) Parameter(name string, value interface{}) *AddOnPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AddOnPollRequest) Header(name string, value interface{}) *AddOnPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AddOnPollRequest) Interval(value time.Duration) *AddOnPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AddOnPollRequest) Status(value int) *AddOnPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AddOnPollRequest) Predicate(value func(*AddOnGetResponse) bool) *AddOnPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AddOnGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AddOnPollRequest) StartContext(ctx context.Context) (response *AddOnPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AddOnPollResponse{
			response: result.(*AddOnGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AddOnPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AddOnPollResponse is the response for the Poll method.
type AddOnPollResponse struct {
	response *AddOnGetResponse
}

// Status returns the response status code.
func (r *AddOnPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AddOnPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AddOnPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AddOnPollResponse) Body() *AddOn {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddOnPollResponse) GetBody() (value *AddOn, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AddOnClient) Poll() *AddOnPollRequest {
	return &AddOnPollRequest{
		request: c.Get(),
	}
}

// AddOnDeleteRequest is the request for the 'delete' method.
type AddOnDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddOnDeleteRequest) Parameter(name string, value interface{}) *AddOnDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnDeleteRequest) Header(name string, value interface{}) *AddOnDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnDeleteRequest) Impersonate(user string) *AddOnDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnDeleteRequest) Send() (result *AddOnDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnDeleteRequest) SendContext(ctx context.Context) (result *AddOnDeleteResponse, err error) {
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
	result = &AddOnDeleteResponse{}
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

// AddOnDeleteResponse is the response for the 'delete' method.
type AddOnDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *AddOnDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// AddOnGetRequest is the request for the 'get' method.
type AddOnGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddOnGetRequest) Parameter(name string, value interface{}) *AddOnGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnGetRequest) Header(name string, value interface{}) *AddOnGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnGetRequest) Impersonate(user string) *AddOnGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnGetRequest) Send() (result *AddOnGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnGetRequest) SendContext(ctx context.Context) (result *AddOnGetResponse, err error) {
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
	result = &AddOnGetResponse{}
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
	err = readAddOnGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddOnGetResponse is the response for the 'get' method.
type AddOnGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddOn
}

// Status returns the response status code.
func (r *AddOnGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddOnGetResponse) Body() *AddOn {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddOnGetResponse) GetBody() (value *AddOn, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddOnUpdateRequest is the request for the 'update' method.
type AddOnUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddOn
}

// Parameter adds a query parameter.
func (r *AddOnUpdateRequest) Parameter(name string, value interface{}) *AddOnUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnUpdateRequest) Header(name string, value interface{}) *AddOnUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnUpdateRequest) Impersonate(user string) *AddOnUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *AddOnUpdateRequest) Body(value *AddOn) *AddOnUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnUpdateRequest) Send() (result *AddOnUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnUpdateRequest) SendContext(ctx context.Context) (result *AddOnUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddOnUpdateRequest(r, buffer)
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
	result = &AddOnUpdateResponse{}
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
	err = readAddOnUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddOnUpdateResponse is the response for the 'update' method.
type AddOnUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddOn
}

// Status returns the response status code.
func (r *AddOnUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddOnUpdateResponse) Body() *AddOn {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddOnUpdateResponse) GetBody() (value *AddOn, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
