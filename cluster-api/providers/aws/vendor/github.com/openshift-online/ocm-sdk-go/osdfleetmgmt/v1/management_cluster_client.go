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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

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

// ManagementClusterClient is the client of the 'management_cluster' resource.
//
// Manages a specific management cluster.
type ManagementClusterClient struct {
	transport http.RoundTripper
	path      string
}

// NewManagementClusterClient creates a new client for the 'management_cluster'
// resource using the given transport to send the requests and receive the
// responses.
func NewManagementClusterClient(transport http.RoundTripper, path string) *ManagementClusterClient {
	return &ManagementClusterClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the management cluster.
func (c *ManagementClusterClient) Delete() *ManagementClusterDeleteRequest {
	return &ManagementClusterDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the management cluster.
func (c *ManagementClusterClient) Get() *ManagementClusterGetRequest {
	return &ManagementClusterGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Post creates a request for the 'post' method.
//
// Updates the management cluster.
func (c *ManagementClusterClient) Post() *ManagementClusterPostRequest {
	return &ManagementClusterPostRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Labels returns the target 'labels' resource.
//
// Reference to the resource that manages the collection of label
func (c *ManagementClusterClient) Labels() *LabelsClient {
	return NewLabelsClient(
		c.transport,
		path.Join(c.path, "labels"),
	)
}

// ManagementClusterPollRequest is the request for the Poll method.
type ManagementClusterPollRequest struct {
	request    *ManagementClusterGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ManagementClusterPollRequest) Parameter(name string, value interface{}) *ManagementClusterPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ManagementClusterPollRequest) Header(name string, value interface{}) *ManagementClusterPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ManagementClusterPollRequest) Interval(value time.Duration) *ManagementClusterPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ManagementClusterPollRequest) Status(value int) *ManagementClusterPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ManagementClusterPollRequest) Predicate(value func(*ManagementClusterGetResponse) bool) *ManagementClusterPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ManagementClusterGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ManagementClusterPollRequest) StartContext(ctx context.Context) (response *ManagementClusterPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ManagementClusterPollResponse{
			response: result.(*ManagementClusterGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ManagementClusterPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ManagementClusterPollResponse is the response for the Poll method.
type ManagementClusterPollResponse struct {
	response *ManagementClusterGetResponse
}

// Status returns the response status code.
func (r *ManagementClusterPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ManagementClusterPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ManagementClusterPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ManagementClusterPollResponse) Body() *ManagementCluster {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ManagementClusterPollResponse) GetBody() (value *ManagementCluster, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ManagementClusterClient) Poll() *ManagementClusterPollRequest {
	return &ManagementClusterPollRequest{
		request: c.Get(),
	}
}

// ManagementClusterDeleteRequest is the request for the 'delete' method.
type ManagementClusterDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ManagementClusterDeleteRequest) Parameter(name string, value interface{}) *ManagementClusterDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ManagementClusterDeleteRequest) Header(name string, value interface{}) *ManagementClusterDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ManagementClusterDeleteRequest) Impersonate(user string) *ManagementClusterDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ManagementClusterDeleteRequest) Send() (result *ManagementClusterDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ManagementClusterDeleteRequest) SendContext(ctx context.Context) (result *ManagementClusterDeleteResponse, err error) {
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
	result = &ManagementClusterDeleteResponse{}
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

// ManagementClusterDeleteResponse is the response for the 'delete' method.
type ManagementClusterDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ManagementClusterDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ManagementClusterDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ManagementClusterDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ManagementClusterGetRequest is the request for the 'get' method.
type ManagementClusterGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ManagementClusterGetRequest) Parameter(name string, value interface{}) *ManagementClusterGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ManagementClusterGetRequest) Header(name string, value interface{}) *ManagementClusterGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ManagementClusterGetRequest) Impersonate(user string) *ManagementClusterGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ManagementClusterGetRequest) Send() (result *ManagementClusterGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ManagementClusterGetRequest) SendContext(ctx context.Context) (result *ManagementClusterGetResponse, err error) {
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
	result = &ManagementClusterGetResponse{}
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
	err = readManagementClusterGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ManagementClusterGetResponse is the response for the 'get' method.
type ManagementClusterGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ManagementCluster
}

// Status returns the response status code.
func (r *ManagementClusterGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ManagementClusterGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ManagementClusterGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ManagementClusterGetResponse) Body() *ManagementCluster {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ManagementClusterGetResponse) GetBody() (value *ManagementCluster, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ManagementClusterPostRequest is the request for the 'post' method.
type ManagementClusterPostRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	request   *ManagementClusterRequestPayload
}

// Parameter adds a query parameter.
func (r *ManagementClusterPostRequest) Parameter(name string, value interface{}) *ManagementClusterPostRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ManagementClusterPostRequest) Header(name string, value interface{}) *ManagementClusterPostRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ManagementClusterPostRequest) Impersonate(user string) *ManagementClusterPostRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Request sets the value of the 'request' parameter.
func (r *ManagementClusterPostRequest) Request(value *ManagementClusterRequestPayload) *ManagementClusterPostRequest {
	r.request = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ManagementClusterPostRequest) Send() (result *ManagementClusterPostResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ManagementClusterPostRequest) SendContext(ctx context.Context) (result *ManagementClusterPostResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeManagementClusterPostRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "POST",
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
	result = &ManagementClusterPostResponse{}
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
	err = readManagementClusterPostResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ManagementClusterPostResponse is the response for the 'post' method.
type ManagementClusterPostResponse struct {
	status   int
	header   http.Header
	err      *errors.Error
	response *ManagementCluster
}

// Status returns the response status code.
func (r *ManagementClusterPostResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ManagementClusterPostResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ManagementClusterPostResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Response returns the value of the 'response' parameter.
func (r *ManagementClusterPostResponse) Response() *ManagementCluster {
	if r == nil {
		return nil
	}
	return r.response
}

// GetResponse returns the value of the 'response' parameter and
// a flag indicating if the parameter has a value.
func (r *ManagementClusterPostResponse) GetResponse() (value *ManagementCluster, ok bool) {
	ok = r != nil && r.response != nil
	if ok {
		value = r.response
	}
	return
}
