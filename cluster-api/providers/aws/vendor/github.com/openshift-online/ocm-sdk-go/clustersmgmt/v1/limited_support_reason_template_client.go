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

// LimitedSupportReasonTemplateClient is the client of the 'limited_support_reason_template' resource.
//
// Manages a specific limited support reason template.
type LimitedSupportReasonTemplateClient struct {
	transport http.RoundTripper
	path      string
}

// NewLimitedSupportReasonTemplateClient creates a new client for the 'limited_support_reason_template'
// resource using the given transport to send the requests and receive the
// responses.
func NewLimitedSupportReasonTemplateClient(transport http.RoundTripper, path string) *LimitedSupportReasonTemplateClient {
	return &LimitedSupportReasonTemplateClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the template.
func (c *LimitedSupportReasonTemplateClient) Get() *LimitedSupportReasonTemplateGetRequest {
	return &LimitedSupportReasonTemplateGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// LimitedSupportReasonTemplatePollRequest is the request for the Poll method.
type LimitedSupportReasonTemplatePollRequest struct {
	request    *LimitedSupportReasonTemplateGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *LimitedSupportReasonTemplatePollRequest) Parameter(name string, value interface{}) *LimitedSupportReasonTemplatePollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *LimitedSupportReasonTemplatePollRequest) Header(name string, value interface{}) *LimitedSupportReasonTemplatePollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *LimitedSupportReasonTemplatePollRequest) Interval(value time.Duration) *LimitedSupportReasonTemplatePollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *LimitedSupportReasonTemplatePollRequest) Status(value int) *LimitedSupportReasonTemplatePollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *LimitedSupportReasonTemplatePollRequest) Predicate(value func(*LimitedSupportReasonTemplateGetResponse) bool) *LimitedSupportReasonTemplatePollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*LimitedSupportReasonTemplateGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *LimitedSupportReasonTemplatePollRequest) StartContext(ctx context.Context) (response *LimitedSupportReasonTemplatePollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &LimitedSupportReasonTemplatePollResponse{
			response: result.(*LimitedSupportReasonTemplateGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *LimitedSupportReasonTemplatePollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// LimitedSupportReasonTemplatePollResponse is the response for the Poll method.
type LimitedSupportReasonTemplatePollResponse struct {
	response *LimitedSupportReasonTemplateGetResponse
}

// Status returns the response status code.
func (r *LimitedSupportReasonTemplatePollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *LimitedSupportReasonTemplatePollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *LimitedSupportReasonTemplatePollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *LimitedSupportReasonTemplatePollResponse) Body() *LimitedSupportReasonTemplate {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *LimitedSupportReasonTemplatePollResponse) GetBody() (value *LimitedSupportReasonTemplate, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *LimitedSupportReasonTemplateClient) Poll() *LimitedSupportReasonTemplatePollRequest {
	return &LimitedSupportReasonTemplatePollRequest{
		request: c.Get(),
	}
}

// LimitedSupportReasonTemplateGetRequest is the request for the 'get' method.
type LimitedSupportReasonTemplateGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *LimitedSupportReasonTemplateGetRequest) Parameter(name string, value interface{}) *LimitedSupportReasonTemplateGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *LimitedSupportReasonTemplateGetRequest) Header(name string, value interface{}) *LimitedSupportReasonTemplateGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *LimitedSupportReasonTemplateGetRequest) Impersonate(user string) *LimitedSupportReasonTemplateGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *LimitedSupportReasonTemplateGetRequest) Send() (result *LimitedSupportReasonTemplateGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *LimitedSupportReasonTemplateGetRequest) SendContext(ctx context.Context) (result *LimitedSupportReasonTemplateGetResponse, err error) {
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
	result = &LimitedSupportReasonTemplateGetResponse{}
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
	err = readLimitedSupportReasonTemplateGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// LimitedSupportReasonTemplateGetResponse is the response for the 'get' method.
type LimitedSupportReasonTemplateGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *LimitedSupportReasonTemplate
}

// Status returns the response status code.
func (r *LimitedSupportReasonTemplateGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *LimitedSupportReasonTemplateGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *LimitedSupportReasonTemplateGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *LimitedSupportReasonTemplateGetResponse) Body() *LimitedSupportReasonTemplate {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *LimitedSupportReasonTemplateGetResponse) GetBody() (value *LimitedSupportReasonTemplate, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
