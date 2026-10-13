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
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// LimitedSupportReasonClient is the client of the 'limited_support_reason' resource.
//
// Manages a specific reason.
type LimitedSupportReasonClient struct {
	transport http.RoundTripper
	path      string
}

// NewLimitedSupportReasonClient creates a new client for the 'limited_support_reason'
// resource using the given transport to send the requests and receive the
// responses.
func NewLimitedSupportReasonClient(transport http.RoundTripper, path string) *LimitedSupportReasonClient {
	return &LimitedSupportReasonClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the reason.
func (c *LimitedSupportReasonClient) Delete() *LimitedSupportReasonDeleteRequest {
	return &LimitedSupportReasonDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the reason.
func (c *LimitedSupportReasonClient) Get() *LimitedSupportReasonGetRequest {
	return &LimitedSupportReasonGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// LimitedSupportReasonPollRequest is the request for the Poll method.
type LimitedSupportReasonPollRequest struct {
	request    *LimitedSupportReasonGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *LimitedSupportReasonPollRequest) Parameter(name string, value interface{}) *LimitedSupportReasonPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *LimitedSupportReasonPollRequest) Header(name string, value interface{}) *LimitedSupportReasonPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *LimitedSupportReasonPollRequest) Interval(value time.Duration) *LimitedSupportReasonPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *LimitedSupportReasonPollRequest) Status(value int) *LimitedSupportReasonPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *LimitedSupportReasonPollRequest) Predicate(value func(*LimitedSupportReasonGetResponse) bool) *LimitedSupportReasonPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*LimitedSupportReasonGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *LimitedSupportReasonPollRequest) StartContext(ctx context.Context) (response *LimitedSupportReasonPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &LimitedSupportReasonPollResponse{
			response: result.(*LimitedSupportReasonGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *LimitedSupportReasonPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// LimitedSupportReasonPollResponse is the response for the Poll method.
type LimitedSupportReasonPollResponse struct {
	response *LimitedSupportReasonGetResponse
}

// Status returns the response status code.
func (r *LimitedSupportReasonPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *LimitedSupportReasonPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *LimitedSupportReasonPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *LimitedSupportReasonPollResponse) Body() *LimitedSupportReason {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *LimitedSupportReasonPollResponse) GetBody() (value *LimitedSupportReason, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *LimitedSupportReasonClient) Poll() *LimitedSupportReasonPollRequest {
	return &LimitedSupportReasonPollRequest{
		request: c.Get(),
	}
}

// LimitedSupportReasonDeleteRequest is the request for the 'delete' method.
type LimitedSupportReasonDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *LimitedSupportReasonDeleteRequest) Parameter(name string, value interface{}) *LimitedSupportReasonDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *LimitedSupportReasonDeleteRequest) Header(name string, value interface{}) *LimitedSupportReasonDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *LimitedSupportReasonDeleteRequest) Impersonate(user string) *LimitedSupportReasonDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *LimitedSupportReasonDeleteRequest) Send() (result *LimitedSupportReasonDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *LimitedSupportReasonDeleteRequest) SendContext(ctx context.Context) (result *LimitedSupportReasonDeleteResponse, err error) {
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
	result = &LimitedSupportReasonDeleteResponse{}
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

// LimitedSupportReasonDeleteResponse is the response for the 'delete' method.
type LimitedSupportReasonDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *LimitedSupportReasonDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *LimitedSupportReasonDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *LimitedSupportReasonDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// LimitedSupportReasonGetRequest is the request for the 'get' method.
type LimitedSupportReasonGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *LimitedSupportReasonGetRequest) Parameter(name string, value interface{}) *LimitedSupportReasonGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *LimitedSupportReasonGetRequest) Header(name string, value interface{}) *LimitedSupportReasonGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *LimitedSupportReasonGetRequest) Impersonate(user string) *LimitedSupportReasonGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *LimitedSupportReasonGetRequest) Send() (result *LimitedSupportReasonGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *LimitedSupportReasonGetRequest) SendContext(ctx context.Context) (result *LimitedSupportReasonGetResponse, err error) {
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
	result = &LimitedSupportReasonGetResponse{}
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
	err = readLimitedSupportReasonGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// LimitedSupportReasonGetResponse is the response for the 'get' method.
type LimitedSupportReasonGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *LimitedSupportReason
}

// Status returns the response status code.
func (r *LimitedSupportReasonGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *LimitedSupportReasonGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *LimitedSupportReasonGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *LimitedSupportReasonGetResponse) Body() *LimitedSupportReason {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *LimitedSupportReasonGetResponse) GetBody() (value *LimitedSupportReason, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
