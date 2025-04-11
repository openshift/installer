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

// HcpKubeletConfigClient is the client of the 'hcp_kubelet_config' resource.
//
// Manages KubeletConfig configuration for Hosted Control Plane clusters. This resource does not support POST operations
// in contrast to the KubeletConfig resource for Classic clusters.
type HcpKubeletConfigClient struct {
	transport http.RoundTripper
	path      string
}

// NewHcpKubeletConfigClient creates a new client for the 'hcp_kubelet_config'
// resource using the given transport to send the requests and receive the
// responses.
func NewHcpKubeletConfigClient(transport http.RoundTripper, path string) *HcpKubeletConfigClient {
	return &HcpKubeletConfigClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the KubeletConfig specified by the id.
func (c *HcpKubeletConfigClient) Delete() *HcpKubeletConfigDeleteRequest {
	return &HcpKubeletConfigDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the KubeletConfig specified by the id.
func (c *HcpKubeletConfigClient) Get() *HcpKubeletConfigGetRequest {
	return &HcpKubeletConfigGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the KubeletConfig specified by the id.
func (c *HcpKubeletConfigClient) Update() *HcpKubeletConfigUpdateRequest {
	return &HcpKubeletConfigUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// HcpKubeletConfigPollRequest is the request for the Poll method.
type HcpKubeletConfigPollRequest struct {
	request    *HcpKubeletConfigGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *HcpKubeletConfigPollRequest) Parameter(name string, value interface{}) *HcpKubeletConfigPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *HcpKubeletConfigPollRequest) Header(name string, value interface{}) *HcpKubeletConfigPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *HcpKubeletConfigPollRequest) Interval(value time.Duration) *HcpKubeletConfigPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *HcpKubeletConfigPollRequest) Status(value int) *HcpKubeletConfigPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *HcpKubeletConfigPollRequest) Predicate(value func(*HcpKubeletConfigGetResponse) bool) *HcpKubeletConfigPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*HcpKubeletConfigGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *HcpKubeletConfigPollRequest) StartContext(ctx context.Context) (response *HcpKubeletConfigPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &HcpKubeletConfigPollResponse{
			response: result.(*HcpKubeletConfigGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *HcpKubeletConfigPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// HcpKubeletConfigPollResponse is the response for the Poll method.
type HcpKubeletConfigPollResponse struct {
	response *HcpKubeletConfigGetResponse
}

// Status returns the response status code.
func (r *HcpKubeletConfigPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *HcpKubeletConfigPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *HcpKubeletConfigPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *HcpKubeletConfigPollResponse) Body() *KubeletConfig {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *HcpKubeletConfigPollResponse) GetBody() (value *KubeletConfig, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *HcpKubeletConfigClient) Poll() *HcpKubeletConfigPollRequest {
	return &HcpKubeletConfigPollRequest{
		request: c.Get(),
	}
}

// HcpKubeletConfigDeleteRequest is the request for the 'delete' method.
type HcpKubeletConfigDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *HcpKubeletConfigDeleteRequest) Parameter(name string, value interface{}) *HcpKubeletConfigDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *HcpKubeletConfigDeleteRequest) Header(name string, value interface{}) *HcpKubeletConfigDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *HcpKubeletConfigDeleteRequest) Impersonate(user string) *HcpKubeletConfigDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *HcpKubeletConfigDeleteRequest) Send() (result *HcpKubeletConfigDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *HcpKubeletConfigDeleteRequest) SendContext(ctx context.Context) (result *HcpKubeletConfigDeleteResponse, err error) {
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
	result = &HcpKubeletConfigDeleteResponse{}
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

// HcpKubeletConfigDeleteResponse is the response for the 'delete' method.
type HcpKubeletConfigDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *HcpKubeletConfigDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *HcpKubeletConfigDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *HcpKubeletConfigDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// HcpKubeletConfigGetRequest is the request for the 'get' method.
type HcpKubeletConfigGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *HcpKubeletConfigGetRequest) Parameter(name string, value interface{}) *HcpKubeletConfigGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *HcpKubeletConfigGetRequest) Header(name string, value interface{}) *HcpKubeletConfigGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *HcpKubeletConfigGetRequest) Impersonate(user string) *HcpKubeletConfigGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *HcpKubeletConfigGetRequest) Send() (result *HcpKubeletConfigGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *HcpKubeletConfigGetRequest) SendContext(ctx context.Context) (result *HcpKubeletConfigGetResponse, err error) {
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
	result = &HcpKubeletConfigGetResponse{}
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
	err = readHcpKubeletConfigGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// HcpKubeletConfigGetResponse is the response for the 'get' method.
type HcpKubeletConfigGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *KubeletConfig
}

// Status returns the response status code.
func (r *HcpKubeletConfigGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *HcpKubeletConfigGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *HcpKubeletConfigGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *HcpKubeletConfigGetResponse) Body() *KubeletConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *HcpKubeletConfigGetResponse) GetBody() (value *KubeletConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// HcpKubeletConfigUpdateRequest is the request for the 'update' method.
type HcpKubeletConfigUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *KubeletConfig
}

// Parameter adds a query parameter.
func (r *HcpKubeletConfigUpdateRequest) Parameter(name string, value interface{}) *HcpKubeletConfigUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *HcpKubeletConfigUpdateRequest) Header(name string, value interface{}) *HcpKubeletConfigUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *HcpKubeletConfigUpdateRequest) Impersonate(user string) *HcpKubeletConfigUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *HcpKubeletConfigUpdateRequest) Body(value *KubeletConfig) *HcpKubeletConfigUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *HcpKubeletConfigUpdateRequest) Send() (result *HcpKubeletConfigUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *HcpKubeletConfigUpdateRequest) SendContext(ctx context.Context) (result *HcpKubeletConfigUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeHcpKubeletConfigUpdateRequest(r, buffer)
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
	result = &HcpKubeletConfigUpdateResponse{}
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
	err = readHcpKubeletConfigUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// HcpKubeletConfigUpdateResponse is the response for the 'update' method.
type HcpKubeletConfigUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *KubeletConfig
}

// Status returns the response status code.
func (r *HcpKubeletConfigUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *HcpKubeletConfigUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *HcpKubeletConfigUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *HcpKubeletConfigUpdateResponse) Body() *KubeletConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *HcpKubeletConfigUpdateResponse) GetBody() (value *KubeletConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
