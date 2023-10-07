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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

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

// AttachmentClient is the client of the 'attachment' resource.
//
// Provides detailed information about a specific attachment.
type AttachmentClient struct {
	transport http.RoundTripper
	path      string
}

// NewAttachmentClient creates a new client for the 'attachment'
// resource using the given transport to send the requests and receive the
// responses.
func NewAttachmentClient(transport http.RoundTripper, path string) *AttachmentClient {
	return &AttachmentClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
func (c *AttachmentClient) Delete() *AttachmentDeleteRequest {
	return &AttachmentDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
func (c *AttachmentClient) Get() *AttachmentGetRequest {
	return &AttachmentGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
func (c *AttachmentClient) Update() *AttachmentUpdateRequest {
	return &AttachmentUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AttachmentPollRequest is the request for the Poll method.
type AttachmentPollRequest struct {
	request    *AttachmentGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AttachmentPollRequest) Parameter(name string, value interface{}) *AttachmentPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AttachmentPollRequest) Header(name string, value interface{}) *AttachmentPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AttachmentPollRequest) Interval(value time.Duration) *AttachmentPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AttachmentPollRequest) Status(value int) *AttachmentPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AttachmentPollRequest) Predicate(value func(*AttachmentGetResponse) bool) *AttachmentPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AttachmentGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AttachmentPollRequest) StartContext(ctx context.Context) (response *AttachmentPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AttachmentPollResponse{
			response: result.(*AttachmentGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AttachmentPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AttachmentPollResponse is the response for the Poll method.
type AttachmentPollResponse struct {
	response *AttachmentGetResponse
}

// Status returns the response status code.
func (r *AttachmentPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AttachmentPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AttachmentPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AttachmentPollResponse) Body() *Attachment {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AttachmentPollResponse) GetBody() (value *Attachment, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AttachmentClient) Poll() *AttachmentPollRequest {
	return &AttachmentPollRequest{
		request: c.Get(),
	}
}

// AttachmentDeleteRequest is the request for the 'delete' method.
type AttachmentDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AttachmentDeleteRequest) Parameter(name string, value interface{}) *AttachmentDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AttachmentDeleteRequest) Header(name string, value interface{}) *AttachmentDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AttachmentDeleteRequest) Impersonate(user string) *AttachmentDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AttachmentDeleteRequest) Send() (result *AttachmentDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AttachmentDeleteRequest) SendContext(ctx context.Context) (result *AttachmentDeleteResponse, err error) {
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
	result = &AttachmentDeleteResponse{}
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

// AttachmentDeleteResponse is the response for the 'delete' method.
type AttachmentDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *AttachmentDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AttachmentDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AttachmentDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// AttachmentGetRequest is the request for the 'get' method.
type AttachmentGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AttachmentGetRequest) Parameter(name string, value interface{}) *AttachmentGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AttachmentGetRequest) Header(name string, value interface{}) *AttachmentGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AttachmentGetRequest) Impersonate(user string) *AttachmentGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AttachmentGetRequest) Send() (result *AttachmentGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AttachmentGetRequest) SendContext(ctx context.Context) (result *AttachmentGetResponse, err error) {
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
	result = &AttachmentGetResponse{}
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
	err = readAttachmentGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AttachmentGetResponse is the response for the 'get' method.
type AttachmentGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Attachment
}

// Status returns the response status code.
func (r *AttachmentGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AttachmentGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AttachmentGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AttachmentGetResponse) Body() *Attachment {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AttachmentGetResponse) GetBody() (value *Attachment, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AttachmentUpdateRequest is the request for the 'update' method.
type AttachmentUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Attachment
}

// Parameter adds a query parameter.
func (r *AttachmentUpdateRequest) Parameter(name string, value interface{}) *AttachmentUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AttachmentUpdateRequest) Header(name string, value interface{}) *AttachmentUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AttachmentUpdateRequest) Impersonate(user string) *AttachmentUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *AttachmentUpdateRequest) Body(value *Attachment) *AttachmentUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AttachmentUpdateRequest) Send() (result *AttachmentUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AttachmentUpdateRequest) SendContext(ctx context.Context) (result *AttachmentUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAttachmentUpdateRequest(r, buffer)
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
	result = &AttachmentUpdateResponse{}
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
	err = readAttachmentUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AttachmentUpdateResponse is the response for the 'update' method.
type AttachmentUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Attachment
}

// Status returns the response status code.
func (r *AttachmentUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AttachmentUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AttachmentUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AttachmentUpdateResponse) Body() *Attachment {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AttachmentUpdateResponse) GetBody() (value *Attachment, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
