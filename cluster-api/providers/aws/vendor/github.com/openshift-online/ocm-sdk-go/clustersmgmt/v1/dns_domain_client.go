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

// DNSDomainClient is the client of the 'DNS_domain' resource.
//
// Manages DNS domain.
type DNSDomainClient struct {
	transport http.RoundTripper
	path      string
}

// NewDNSDomainClient creates a new client for the 'DNS_domain'
// resource using the given transport to send the requests and receive the
// responses.
func NewDNSDomainClient(transport http.RoundTripper, path string) *DNSDomainClient {
	return &DNSDomainClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Delete the DNS domain.
func (c *DNSDomainClient) Delete() *DNSDomainDeleteRequest {
	return &DNSDomainDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the DNS domain.
func (c *DNSDomainClient) Get() *DNSDomainGetRequest {
	return &DNSDomainGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// DNSDomainPollRequest is the request for the Poll method.
type DNSDomainPollRequest struct {
	request    *DNSDomainGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *DNSDomainPollRequest) Parameter(name string, value interface{}) *DNSDomainPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *DNSDomainPollRequest) Header(name string, value interface{}) *DNSDomainPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *DNSDomainPollRequest) Interval(value time.Duration) *DNSDomainPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *DNSDomainPollRequest) Status(value int) *DNSDomainPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *DNSDomainPollRequest) Predicate(value func(*DNSDomainGetResponse) bool) *DNSDomainPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*DNSDomainGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *DNSDomainPollRequest) StartContext(ctx context.Context) (response *DNSDomainPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &DNSDomainPollResponse{
			response: result.(*DNSDomainGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *DNSDomainPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// DNSDomainPollResponse is the response for the Poll method.
type DNSDomainPollResponse struct {
	response *DNSDomainGetResponse
}

// Status returns the response status code.
func (r *DNSDomainPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *DNSDomainPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *DNSDomainPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *DNSDomainPollResponse) Body() *DNSDomain {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *DNSDomainPollResponse) GetBody() (value *DNSDomain, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *DNSDomainClient) Poll() *DNSDomainPollRequest {
	return &DNSDomainPollRequest{
		request: c.Get(),
	}
}

// DNSDomainDeleteRequest is the request for the 'delete' method.
type DNSDomainDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *DNSDomainDeleteRequest) Parameter(name string, value interface{}) *DNSDomainDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DNSDomainDeleteRequest) Header(name string, value interface{}) *DNSDomainDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DNSDomainDeleteRequest) Impersonate(user string) *DNSDomainDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DNSDomainDeleteRequest) Send() (result *DNSDomainDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DNSDomainDeleteRequest) SendContext(ctx context.Context) (result *DNSDomainDeleteResponse, err error) {
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
	result = &DNSDomainDeleteResponse{}
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

// DNSDomainDeleteResponse is the response for the 'delete' method.
type DNSDomainDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *DNSDomainDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DNSDomainDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DNSDomainDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// DNSDomainGetRequest is the request for the 'get' method.
type DNSDomainGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *DNSDomainGetRequest) Parameter(name string, value interface{}) *DNSDomainGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DNSDomainGetRequest) Header(name string, value interface{}) *DNSDomainGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DNSDomainGetRequest) Impersonate(user string) *DNSDomainGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DNSDomainGetRequest) Send() (result *DNSDomainGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DNSDomainGetRequest) SendContext(ctx context.Context) (result *DNSDomainGetResponse, err error) {
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
	result = &DNSDomainGetResponse{}
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
	err = readDNSDomainGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DNSDomainGetResponse is the response for the 'get' method.
type DNSDomainGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *DNSDomain
}

// Status returns the response status code.
func (r *DNSDomainGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DNSDomainGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DNSDomainGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *DNSDomainGetResponse) Body() *DNSDomain {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *DNSDomainGetResponse) GetBody() (value *DNSDomain, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
