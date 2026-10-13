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

// PeerDependencyClient is the client of the 'peer_dependency' resource.
//
// Provides detailed information about the specified peer dependency.
type PeerDependencyClient struct {
	transport http.RoundTripper
	path      string
}

// NewPeerDependencyClient creates a new client for the 'peer_dependency'
// resource using the given transport to send the requests and receive the
// responses.
func NewPeerDependencyClient(transport http.RoundTripper, path string) *PeerDependencyClient {
	return &PeerDependencyClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
func (c *PeerDependencyClient) Delete() *PeerDependencyDeleteRequest {
	return &PeerDependencyDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
func (c *PeerDependencyClient) Get() *PeerDependencyGetRequest {
	return &PeerDependencyGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
func (c *PeerDependencyClient) Update() *PeerDependencyUpdateRequest {
	return &PeerDependencyUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// PeerDependencyPollRequest is the request for the Poll method.
type PeerDependencyPollRequest struct {
	request    *PeerDependencyGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *PeerDependencyPollRequest) Parameter(name string, value interface{}) *PeerDependencyPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *PeerDependencyPollRequest) Header(name string, value interface{}) *PeerDependencyPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *PeerDependencyPollRequest) Interval(value time.Duration) *PeerDependencyPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *PeerDependencyPollRequest) Status(value int) *PeerDependencyPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *PeerDependencyPollRequest) Predicate(value func(*PeerDependencyGetResponse) bool) *PeerDependencyPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*PeerDependencyGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *PeerDependencyPollRequest) StartContext(ctx context.Context) (response *PeerDependencyPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &PeerDependencyPollResponse{
			response: result.(*PeerDependencyGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *PeerDependencyPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// PeerDependencyPollResponse is the response for the Poll method.
type PeerDependencyPollResponse struct {
	response *PeerDependencyGetResponse
}

// Status returns the response status code.
func (r *PeerDependencyPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *PeerDependencyPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *PeerDependencyPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *PeerDependencyPollResponse) Body() *Service {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *PeerDependencyPollResponse) GetBody() (value *Service, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *PeerDependencyClient) Poll() *PeerDependencyPollRequest {
	return &PeerDependencyPollRequest{
		request: c.Get(),
	}
}

// PeerDependencyDeleteRequest is the request for the 'delete' method.
type PeerDependencyDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *PeerDependencyDeleteRequest) Parameter(name string, value interface{}) *PeerDependencyDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PeerDependencyDeleteRequest) Header(name string, value interface{}) *PeerDependencyDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PeerDependencyDeleteRequest) Impersonate(user string) *PeerDependencyDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PeerDependencyDeleteRequest) Send() (result *PeerDependencyDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PeerDependencyDeleteRequest) SendContext(ctx context.Context) (result *PeerDependencyDeleteResponse, err error) {
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
	result = &PeerDependencyDeleteResponse{}
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

// PeerDependencyDeleteResponse is the response for the 'delete' method.
type PeerDependencyDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *PeerDependencyDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PeerDependencyDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PeerDependencyDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// PeerDependencyGetRequest is the request for the 'get' method.
type PeerDependencyGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *PeerDependencyGetRequest) Parameter(name string, value interface{}) *PeerDependencyGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PeerDependencyGetRequest) Header(name string, value interface{}) *PeerDependencyGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PeerDependencyGetRequest) Impersonate(user string) *PeerDependencyGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PeerDependencyGetRequest) Send() (result *PeerDependencyGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PeerDependencyGetRequest) SendContext(ctx context.Context) (result *PeerDependencyGetResponse, err error) {
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
	result = &PeerDependencyGetResponse{}
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
	err = readPeerDependencyGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// PeerDependencyGetResponse is the response for the 'get' method.
type PeerDependencyGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Service
}

// Status returns the response status code.
func (r *PeerDependencyGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PeerDependencyGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PeerDependencyGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *PeerDependencyGetResponse) Body() *Service {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *PeerDependencyGetResponse) GetBody() (value *Service, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// PeerDependencyUpdateRequest is the request for the 'update' method.
type PeerDependencyUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *PeerDependency
}

// Parameter adds a query parameter.
func (r *PeerDependencyUpdateRequest) Parameter(name string, value interface{}) *PeerDependencyUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *PeerDependencyUpdateRequest) Header(name string, value interface{}) *PeerDependencyUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *PeerDependencyUpdateRequest) Impersonate(user string) *PeerDependencyUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *PeerDependencyUpdateRequest) Body(value *PeerDependency) *PeerDependencyUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *PeerDependencyUpdateRequest) Send() (result *PeerDependencyUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *PeerDependencyUpdateRequest) SendContext(ctx context.Context) (result *PeerDependencyUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writePeerDependencyUpdateRequest(r, buffer)
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
	result = &PeerDependencyUpdateResponse{}
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
	err = readPeerDependencyUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// PeerDependencyUpdateResponse is the response for the 'update' method.
type PeerDependencyUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *PeerDependency
}

// Status returns the response status code.
func (r *PeerDependencyUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *PeerDependencyUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *PeerDependencyUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *PeerDependencyUpdateResponse) Body() *PeerDependency {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *PeerDependencyUpdateResponse) GetBody() (value *PeerDependency, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
