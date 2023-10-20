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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

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

// BillingModelClient is the client of the 'billing_model' resource.
type BillingModelClient struct {
	transport http.RoundTripper
	path      string
}

// NewBillingModelClient creates a new client for the 'billing_model'
// resource using the given transport to send the requests and receive the
// responses.
func NewBillingModelClient(transport http.RoundTripper, path string) *BillingModelClient {
	return &BillingModelClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the billing model
func (c *BillingModelClient) Get() *BillingModelGetRequest {
	return &BillingModelGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// BillingModelPollRequest is the request for the Poll method.
type BillingModelPollRequest struct {
	request    *BillingModelGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *BillingModelPollRequest) Parameter(name string, value interface{}) *BillingModelPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *BillingModelPollRequest) Header(name string, value interface{}) *BillingModelPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *BillingModelPollRequest) Interval(value time.Duration) *BillingModelPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *BillingModelPollRequest) Status(value int) *BillingModelPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *BillingModelPollRequest) Predicate(value func(*BillingModelGetResponse) bool) *BillingModelPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*BillingModelGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *BillingModelPollRequest) StartContext(ctx context.Context) (response *BillingModelPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &BillingModelPollResponse{
			response: result.(*BillingModelGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *BillingModelPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// BillingModelPollResponse is the response for the Poll method.
type BillingModelPollResponse struct {
	response *BillingModelGetResponse
}

// Status returns the response status code.
func (r *BillingModelPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *BillingModelPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *BillingModelPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *BillingModelPollResponse) Body() *BillingModelItem {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *BillingModelPollResponse) GetBody() (value *BillingModelItem, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *BillingModelClient) Poll() *BillingModelPollRequest {
	return &BillingModelPollRequest{
		request: c.Get(),
	}
}

// BillingModelGetRequest is the request for the 'get' method.
type BillingModelGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *BillingModelGetRequest) Parameter(name string, value interface{}) *BillingModelGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *BillingModelGetRequest) Header(name string, value interface{}) *BillingModelGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *BillingModelGetRequest) Impersonate(user string) *BillingModelGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *BillingModelGetRequest) Send() (result *BillingModelGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *BillingModelGetRequest) SendContext(ctx context.Context) (result *BillingModelGetResponse, err error) {
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
	result = &BillingModelGetResponse{}
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
	err = readBillingModelGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// BillingModelGetResponse is the response for the 'get' method.
type BillingModelGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *BillingModelItem
}

// Status returns the response status code.
func (r *BillingModelGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *BillingModelGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *BillingModelGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *BillingModelGetResponse) Body() *BillingModelItem {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *BillingModelGetResponse) GetBody() (value *BillingModelItem, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
