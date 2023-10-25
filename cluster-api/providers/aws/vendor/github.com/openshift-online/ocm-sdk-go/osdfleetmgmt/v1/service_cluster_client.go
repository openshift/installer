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

// ServiceClusterClient is the client of the 'service_cluster' resource.
//
// Manages a specific service cluster.
type ServiceClusterClient struct {
	transport http.RoundTripper
	path      string
}

// NewServiceClusterClient creates a new client for the 'service_cluster'
// resource using the given transport to send the requests and receive the
// responses.
func NewServiceClusterClient(transport http.RoundTripper, path string) *ServiceClusterClient {
	return &ServiceClusterClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the service cluster.
func (c *ServiceClusterClient) Delete() *ServiceClusterDeleteRequest {
	return &ServiceClusterDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the cluster.
func (c *ServiceClusterClient) Get() *ServiceClusterGetRequest {
	return &ServiceClusterGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Post creates a request for the 'post' method.
//
// Updates the service cluster.
func (c *ServiceClusterClient) Post() *ServiceClusterPostRequest {
	return &ServiceClusterPostRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Labels returns the target 'labels' resource.
//
// Reference to the resource that manages the collection of label
func (c *ServiceClusterClient) Labels() *LabelsClient {
	return NewLabelsClient(
		c.transport,
		path.Join(c.path, "labels"),
	)
}

// ServiceClusterPollRequest is the request for the Poll method.
type ServiceClusterPollRequest struct {
	request    *ServiceClusterGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ServiceClusterPollRequest) Parameter(name string, value interface{}) *ServiceClusterPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ServiceClusterPollRequest) Header(name string, value interface{}) *ServiceClusterPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ServiceClusterPollRequest) Interval(value time.Duration) *ServiceClusterPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ServiceClusterPollRequest) Status(value int) *ServiceClusterPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ServiceClusterPollRequest) Predicate(value func(*ServiceClusterGetResponse) bool) *ServiceClusterPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ServiceClusterGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ServiceClusterPollRequest) StartContext(ctx context.Context) (response *ServiceClusterPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ServiceClusterPollResponse{
			response: result.(*ServiceClusterGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ServiceClusterPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ServiceClusterPollResponse is the response for the Poll method.
type ServiceClusterPollResponse struct {
	response *ServiceClusterGetResponse
}

// Status returns the response status code.
func (r *ServiceClusterPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ServiceClusterPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ServiceClusterPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ServiceClusterPollResponse) Body() *ServiceCluster {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceClusterPollResponse) GetBody() (value *ServiceCluster, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ServiceClusterClient) Poll() *ServiceClusterPollRequest {
	return &ServiceClusterPollRequest{
		request: c.Get(),
	}
}

// ServiceClusterDeleteRequest is the request for the 'delete' method.
type ServiceClusterDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ServiceClusterDeleteRequest) Parameter(name string, value interface{}) *ServiceClusterDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServiceClusterDeleteRequest) Header(name string, value interface{}) *ServiceClusterDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServiceClusterDeleteRequest) Impersonate(user string) *ServiceClusterDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServiceClusterDeleteRequest) Send() (result *ServiceClusterDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServiceClusterDeleteRequest) SendContext(ctx context.Context) (result *ServiceClusterDeleteResponse, err error) {
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
	result = &ServiceClusterDeleteResponse{}
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

// ServiceClusterDeleteResponse is the response for the 'delete' method.
type ServiceClusterDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ServiceClusterDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServiceClusterDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServiceClusterDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ServiceClusterGetRequest is the request for the 'get' method.
type ServiceClusterGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ServiceClusterGetRequest) Parameter(name string, value interface{}) *ServiceClusterGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServiceClusterGetRequest) Header(name string, value interface{}) *ServiceClusterGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServiceClusterGetRequest) Impersonate(user string) *ServiceClusterGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServiceClusterGetRequest) Send() (result *ServiceClusterGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServiceClusterGetRequest) SendContext(ctx context.Context) (result *ServiceClusterGetResponse, err error) {
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
	result = &ServiceClusterGetResponse{}
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
	err = readServiceClusterGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ServiceClusterGetResponse is the response for the 'get' method.
type ServiceClusterGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ServiceCluster
}

// Status returns the response status code.
func (r *ServiceClusterGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServiceClusterGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServiceClusterGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ServiceClusterGetResponse) Body() *ServiceCluster {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceClusterGetResponse) GetBody() (value *ServiceCluster, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ServiceClusterPostRequest is the request for the 'post' method.
type ServiceClusterPostRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	request   *ServiceClusterRequestPayload
}

// Parameter adds a query parameter.
func (r *ServiceClusterPostRequest) Parameter(name string, value interface{}) *ServiceClusterPostRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServiceClusterPostRequest) Header(name string, value interface{}) *ServiceClusterPostRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServiceClusterPostRequest) Impersonate(user string) *ServiceClusterPostRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Request sets the value of the 'request' parameter.
func (r *ServiceClusterPostRequest) Request(value *ServiceClusterRequestPayload) *ServiceClusterPostRequest {
	r.request = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServiceClusterPostRequest) Send() (result *ServiceClusterPostResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServiceClusterPostRequest) SendContext(ctx context.Context) (result *ServiceClusterPostResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeServiceClusterPostRequest(r, buffer)
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
	result = &ServiceClusterPostResponse{}
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
	err = readServiceClusterPostResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ServiceClusterPostResponse is the response for the 'post' method.
type ServiceClusterPostResponse struct {
	status   int
	header   http.Header
	err      *errors.Error
	response *ServiceCluster
}

// Status returns the response status code.
func (r *ServiceClusterPostResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServiceClusterPostResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServiceClusterPostResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Response returns the value of the 'response' parameter.
func (r *ServiceClusterPostResponse) Response() *ServiceCluster {
	if r == nil {
		return nil
	}
	return r.response
}

// GetResponse returns the value of the 'response' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceClusterPostResponse) GetResponse() (value *ServiceCluster, ok bool) {
	ok = r != nil && r.response != nil
	if ok {
		value = r.response
	}
	return
}
