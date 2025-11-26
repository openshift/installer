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

// ImageMirrorClient is the client of the 'image_mirror' resource.
//
// Manages a specific image mirror configuration for a cluster.
type ImageMirrorClient struct {
	transport http.RoundTripper
	path      string
}

// NewImageMirrorClient creates a new client for the 'image_mirror'
// resource using the given transport to send the requests and receive the
// responses.
func NewImageMirrorClient(transport http.RoundTripper, path string) *ImageMirrorClient {
	return &ImageMirrorClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the image mirror configuration.
func (c *ImageMirrorClient) Delete() *ImageMirrorDeleteRequest {
	return &ImageMirrorDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the image mirror.
func (c *ImageMirrorClient) Get() *ImageMirrorGetRequest {
	return &ImageMirrorGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the image mirror configuration.
// Note: Id and Source fields are immutable and cannot be updated.
// The mirrors array is completely replaced, not merged.
func (c *ImageMirrorClient) Update() *ImageMirrorUpdateRequest {
	return &ImageMirrorUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ImageMirrorPollRequest is the request for the Poll method.
type ImageMirrorPollRequest struct {
	request    *ImageMirrorGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ImageMirrorPollRequest) Parameter(name string, value interface{}) *ImageMirrorPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ImageMirrorPollRequest) Header(name string, value interface{}) *ImageMirrorPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ImageMirrorPollRequest) Interval(value time.Duration) *ImageMirrorPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ImageMirrorPollRequest) Status(value int) *ImageMirrorPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ImageMirrorPollRequest) Predicate(value func(*ImageMirrorGetResponse) bool) *ImageMirrorPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ImageMirrorGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ImageMirrorPollRequest) StartContext(ctx context.Context) (response *ImageMirrorPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ImageMirrorPollResponse{
			response: result.(*ImageMirrorGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ImageMirrorPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ImageMirrorPollResponse is the response for the Poll method.
type ImageMirrorPollResponse struct {
	response *ImageMirrorGetResponse
}

// Status returns the response status code.
func (r *ImageMirrorPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ImageMirrorPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ImageMirrorPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ImageMirrorPollResponse) Body() *ImageMirror {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ImageMirrorPollResponse) GetBody() (value *ImageMirror, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ImageMirrorClient) Poll() *ImageMirrorPollRequest {
	return &ImageMirrorPollRequest{
		request: c.Get(),
	}
}

// ImageMirrorDeleteRequest is the request for the 'delete' method.
type ImageMirrorDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ImageMirrorDeleteRequest) Parameter(name string, value interface{}) *ImageMirrorDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ImageMirrorDeleteRequest) Header(name string, value interface{}) *ImageMirrorDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ImageMirrorDeleteRequest) Impersonate(user string) *ImageMirrorDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ImageMirrorDeleteRequest) Send() (result *ImageMirrorDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ImageMirrorDeleteRequest) SendContext(ctx context.Context) (result *ImageMirrorDeleteResponse, err error) {
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
	result = &ImageMirrorDeleteResponse{}
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

// ImageMirrorDeleteResponse is the response for the 'delete' method.
type ImageMirrorDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ImageMirrorDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ImageMirrorDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ImageMirrorDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ImageMirrorGetRequest is the request for the 'get' method.
type ImageMirrorGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ImageMirrorGetRequest) Parameter(name string, value interface{}) *ImageMirrorGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ImageMirrorGetRequest) Header(name string, value interface{}) *ImageMirrorGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ImageMirrorGetRequest) Impersonate(user string) *ImageMirrorGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ImageMirrorGetRequest) Send() (result *ImageMirrorGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ImageMirrorGetRequest) SendContext(ctx context.Context) (result *ImageMirrorGetResponse, err error) {
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
	result = &ImageMirrorGetResponse{}
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
	err = readImageMirrorGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ImageMirrorGetResponse is the response for the 'get' method.
type ImageMirrorGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ImageMirror
}

// Status returns the response status code.
func (r *ImageMirrorGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ImageMirrorGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ImageMirrorGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ImageMirrorGetResponse) Body() *ImageMirror {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ImageMirrorGetResponse) GetBody() (value *ImageMirror, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ImageMirrorUpdateRequest is the request for the 'update' method.
type ImageMirrorUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ImageMirror
}

// Parameter adds a query parameter.
func (r *ImageMirrorUpdateRequest) Parameter(name string, value interface{}) *ImageMirrorUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ImageMirrorUpdateRequest) Header(name string, value interface{}) *ImageMirrorUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ImageMirrorUpdateRequest) Impersonate(user string) *ImageMirrorUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ImageMirrorUpdateRequest) Body(value *ImageMirror) *ImageMirrorUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ImageMirrorUpdateRequest) Send() (result *ImageMirrorUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ImageMirrorUpdateRequest) SendContext(ctx context.Context) (result *ImageMirrorUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeImageMirrorUpdateRequest(r, buffer)
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
	result = &ImageMirrorUpdateResponse{}
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
	err = readImageMirrorUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ImageMirrorUpdateResponse is the response for the 'update' method.
type ImageMirrorUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ImageMirror
}

// Status returns the response status code.
func (r *ImageMirrorUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ImageMirrorUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ImageMirrorUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ImageMirrorUpdateResponse) Body() *ImageMirror {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ImageMirrorUpdateResponse) GetBody() (value *ImageMirror, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
