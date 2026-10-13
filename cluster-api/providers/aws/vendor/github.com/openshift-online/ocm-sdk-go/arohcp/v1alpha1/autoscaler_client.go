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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

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

// AutoscalerClient is the client of the 'autoscaler' resource.
//
// Manages global autoscaler configurations for a cluster.
type AutoscalerClient struct {
	transport http.RoundTripper
	path      string
}

// NewAutoscalerClient creates a new client for the 'autoscaler'
// resource using the given transport to send the requests and receive the
// responses.
func NewAutoscalerClient(transport http.RoundTripper, path string) *AutoscalerClient {
	return &AutoscalerClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the autoscaler of a cluster.
func (c *AutoscalerClient) Get() *AutoscalerGetRequest {
	return &AutoscalerGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the cluster autoscaler.
func (c *AutoscalerClient) Update() *AutoscalerUpdateRequest {
	return &AutoscalerUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AutoscalerPollRequest is the request for the Poll method.
type AutoscalerPollRequest struct {
	request    *AutoscalerGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AutoscalerPollRequest) Parameter(name string, value interface{}) *AutoscalerPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AutoscalerPollRequest) Header(name string, value interface{}) *AutoscalerPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AutoscalerPollRequest) Interval(value time.Duration) *AutoscalerPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AutoscalerPollRequest) Status(value int) *AutoscalerPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AutoscalerPollRequest) Predicate(value func(*AutoscalerGetResponse) bool) *AutoscalerPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AutoscalerGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AutoscalerPollRequest) StartContext(ctx context.Context) (response *AutoscalerPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AutoscalerPollResponse{
			response: result.(*AutoscalerGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AutoscalerPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AutoscalerPollResponse is the response for the Poll method.
type AutoscalerPollResponse struct {
	response *AutoscalerGetResponse
}

// Status returns the response status code.
func (r *AutoscalerPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AutoscalerPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AutoscalerPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AutoscalerPollResponse) Body() *ClusterAutoscaler {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AutoscalerPollResponse) GetBody() (value *ClusterAutoscaler, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AutoscalerClient) Poll() *AutoscalerPollRequest {
	return &AutoscalerPollRequest{
		request: c.Get(),
	}
}

// AutoscalerGetRequest is the request for the 'get' method.
type AutoscalerGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AutoscalerGetRequest) Parameter(name string, value interface{}) *AutoscalerGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AutoscalerGetRequest) Header(name string, value interface{}) *AutoscalerGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AutoscalerGetRequest) Impersonate(user string) *AutoscalerGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AutoscalerGetRequest) Send() (result *AutoscalerGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AutoscalerGetRequest) SendContext(ctx context.Context) (result *AutoscalerGetResponse, err error) {
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
	result = &AutoscalerGetResponse{}
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
	err = readAutoscalerGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AutoscalerGetResponse is the response for the 'get' method.
type AutoscalerGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ClusterAutoscaler
}

// Status returns the response status code.
func (r *AutoscalerGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AutoscalerGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AutoscalerGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AutoscalerGetResponse) Body() *ClusterAutoscaler {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AutoscalerGetResponse) GetBody() (value *ClusterAutoscaler, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AutoscalerUpdateRequest is the request for the 'update' method.
type AutoscalerUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ClusterAutoscaler
}

// Parameter adds a query parameter.
func (r *AutoscalerUpdateRequest) Parameter(name string, value interface{}) *AutoscalerUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AutoscalerUpdateRequest) Header(name string, value interface{}) *AutoscalerUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AutoscalerUpdateRequest) Impersonate(user string) *AutoscalerUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *AutoscalerUpdateRequest) Body(value *ClusterAutoscaler) *AutoscalerUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AutoscalerUpdateRequest) Send() (result *AutoscalerUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AutoscalerUpdateRequest) SendContext(ctx context.Context) (result *AutoscalerUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAutoscalerUpdateRequest(r, buffer)
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
	result = &AutoscalerUpdateResponse{}
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
	err = readAutoscalerUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AutoscalerUpdateResponse is the response for the 'update' method.
type AutoscalerUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ClusterAutoscaler
}

// Status returns the response status code.
func (r *AutoscalerUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AutoscalerUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AutoscalerUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AutoscalerUpdateResponse) Body() *ClusterAutoscaler {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AutoscalerUpdateResponse) GetBody() (value *ClusterAutoscaler, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
