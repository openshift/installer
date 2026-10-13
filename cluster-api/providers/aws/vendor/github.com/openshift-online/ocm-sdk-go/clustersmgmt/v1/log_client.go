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

// LogClient is the client of the 'log' resource.
//
// Manages a specific log.
type LogClient struct {
	transport http.RoundTripper
	path      string
}

// NewLogClient creates a new client for the 'log'
// resource using the given transport to send the requests and receive the
// responses.
func NewLogClient(transport http.RoundTripper, path string) *LogClient {
	return &LogClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the log.
func (c *LogClient) Get() *LogGetRequest {
	return &LogGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// LogPollRequest is the request for the Poll method.
type LogPollRequest struct {
	request    *LogGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *LogPollRequest) Parameter(name string, value interface{}) *LogPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *LogPollRequest) Header(name string, value interface{}) *LogPollRequest {
	r.request.Header(name, value)
	return r
}

// Offset sets the value of the 'offset' parameter for all the requests that
// will be used to retrieve the object.
//
// Line offset to start logs from. if 0 retreive entire log.
// If offset > #lines return an empty log.
func (r *LogPollRequest) Offset(value int) *LogPollRequest {
	r.request.Offset(value)
	return r
}

// Tail sets the value of the 'tail' parameter for all the requests that
// will be used to retrieve the object.
//
// Returns the number of tail lines from the end of the log.
// If there are no line breaks or the number of lines < tail
// return the entire log.
// Either 'tail' or 'offset' can be set. Not both.
func (r *LogPollRequest) Tail(value int) *LogPollRequest {
	r.request.Tail(value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *LogPollRequest) Interval(value time.Duration) *LogPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *LogPollRequest) Status(value int) *LogPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *LogPollRequest) Predicate(value func(*LogGetResponse) bool) *LogPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*LogGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *LogPollRequest) StartContext(ctx context.Context) (response *LogPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &LogPollResponse{
			response: result.(*LogGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *LogPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// LogPollResponse is the response for the Poll method.
type LogPollResponse struct {
	response *LogGetResponse
}

// Status returns the response status code.
func (r *LogPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *LogPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *LogPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
//
// Retreived log.
func (r *LogPollResponse) Body() *Log {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Retreived log.
func (r *LogPollResponse) GetBody() (value *Log, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *LogClient) Poll() *LogPollRequest {
	return &LogPollRequest{
		request: c.Get(),
	}
}

// LogGetRequest is the request for the 'get' method.
type LogGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	offset    *int
	tail      *int
}

// Parameter adds a query parameter.
func (r *LogGetRequest) Parameter(name string, value interface{}) *LogGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *LogGetRequest) Header(name string, value interface{}) *LogGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *LogGetRequest) Impersonate(user string) *LogGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Offset sets the value of the 'offset' parameter.
//
// Line offset to start logs from. if 0 retreive entire log.
// If offset > #lines return an empty log.
func (r *LogGetRequest) Offset(value int) *LogGetRequest {
	r.offset = &value
	return r
}

// Tail sets the value of the 'tail' parameter.
//
// Returns the number of tail lines from the end of the log.
// If there are no line breaks or the number of lines < tail
// return the entire log.
// Either 'tail' or 'offset' can be set. Not both.
func (r *LogGetRequest) Tail(value int) *LogGetRequest {
	r.tail = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *LogGetRequest) Send() (result *LogGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *LogGetRequest) SendContext(ctx context.Context) (result *LogGetResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.offset != nil {
		helpers.AddValue(&query, "offset", *r.offset)
	}
	if r.tail != nil {
		helpers.AddValue(&query, "tail", *r.tail)
	}
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
	result = &LogGetResponse{}
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
	err = readLogGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// LogGetResponse is the response for the 'get' method.
type LogGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Log
}

// Status returns the response status code.
func (r *LogGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *LogGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *LogGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Retreived log.
func (r *LogGetResponse) Body() *Log {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Retreived log.
func (r *LogGetResponse) GetBody() (value *Log, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
