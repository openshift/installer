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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// AccessRequestClient is the client of the 'access_request' resource.
//
// Manages a specific access request.
type AccessRequestClient struct {
	transport http.RoundTripper
	path      string
}

// NewAccessRequestClient creates a new client for the 'access_request'
// resource using the given transport to send the requests and receive the
// responses.
func NewAccessRequestClient(transport http.RoundTripper, path string) *AccessRequestClient {
	return &AccessRequestClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the access request.
func (c *AccessRequestClient) Get() *AccessRequestGetRequest {
	return &AccessRequestGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Decisions returns the target 'decisions' resource.
//
// Reference to the resource that manages the collection of decisions.
func (c *AccessRequestClient) Decisions() *DecisionsClient {
	return NewDecisionsClient(
		c.transport,
		path.Join(c.path, "decisions"),
	)
}

// AccessRequestPollRequest is the request for the Poll method.
type AccessRequestPollRequest struct {
	request    *AccessRequestGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AccessRequestPollRequest) Parameter(name string, value interface{}) *AccessRequestPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AccessRequestPollRequest) Header(name string, value interface{}) *AccessRequestPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AccessRequestPollRequest) Interval(value time.Duration) *AccessRequestPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AccessRequestPollRequest) Status(value int) *AccessRequestPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AccessRequestPollRequest) Predicate(value func(*AccessRequestGetResponse) bool) *AccessRequestPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AccessRequestGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AccessRequestPollRequest) StartContext(ctx context.Context) (response *AccessRequestPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AccessRequestPollResponse{
			response: result.(*AccessRequestGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AccessRequestPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AccessRequestPollResponse is the response for the Poll method.
type AccessRequestPollResponse struct {
	response *AccessRequestGetResponse
}

// Status returns the response status code.
func (r *AccessRequestPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AccessRequestPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AccessRequestPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AccessRequestPollResponse) Body() *AccessRequest {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AccessRequestPollResponse) GetBody() (value *AccessRequest, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AccessRequestClient) Poll() *AccessRequestPollRequest {
	return &AccessRequestPollRequest{
		request: c.Get(),
	}
}

// AccessRequestGetRequest is the request for the 'get' method.
type AccessRequestGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AccessRequestGetRequest) Parameter(name string, value interface{}) *AccessRequestGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccessRequestGetRequest) Header(name string, value interface{}) *AccessRequestGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccessRequestGetRequest) Impersonate(user string) *AccessRequestGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccessRequestGetRequest) Send() (result *AccessRequestGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccessRequestGetRequest) SendContext(ctx context.Context) (result *AccessRequestGetResponse, err error) {
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
	result = &AccessRequestGetResponse{}
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
	err = readAccessRequestGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AccessRequestGetResponse is the response for the 'get' method.
type AccessRequestGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AccessRequest
}

// Status returns the response status code.
func (r *AccessRequestGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccessRequestGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccessRequestGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AccessRequestGetResponse) Body() *AccessRequest {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AccessRequestGetResponse) GetBody() (value *AccessRequest, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
