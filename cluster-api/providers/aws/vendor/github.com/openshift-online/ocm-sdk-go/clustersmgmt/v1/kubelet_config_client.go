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

// KubeletConfigClient is the client of the 'kubelet_config' resource.
//
// Manages global KubeletConfig configuration for worker nodes in a cluster.
type KubeletConfigClient struct {
	transport http.RoundTripper
	path      string
}

// NewKubeletConfigClient creates a new client for the 'kubelet_config'
// resource using the given transport to send the requests and receive the
// responses.
func NewKubeletConfigClient(transport http.RoundTripper, path string) *KubeletConfigClient {
	return &KubeletConfigClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the cluster KubeletConfig
func (c *KubeletConfigClient) Delete() *KubeletConfigDeleteRequest {
	return &KubeletConfigDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the KubeletConfig for a cluster
func (c *KubeletConfigClient) Get() *KubeletConfigGetRequest {
	return &KubeletConfigGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Post creates a request for the 'post' method.
//
// Creates a new cluster KubeletConfig
func (c *KubeletConfigClient) Post() *KubeletConfigPostRequest {
	return &KubeletConfigPostRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the existing cluster KubeletConfig
func (c *KubeletConfigClient) Update() *KubeletConfigUpdateRequest {
	return &KubeletConfigUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// KubeletConfigPollRequest is the request for the Poll method.
type KubeletConfigPollRequest struct {
	request    *KubeletConfigGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *KubeletConfigPollRequest) Parameter(name string, value interface{}) *KubeletConfigPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *KubeletConfigPollRequest) Header(name string, value interface{}) *KubeletConfigPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *KubeletConfigPollRequest) Interval(value time.Duration) *KubeletConfigPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *KubeletConfigPollRequest) Status(value int) *KubeletConfigPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *KubeletConfigPollRequest) Predicate(value func(*KubeletConfigGetResponse) bool) *KubeletConfigPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*KubeletConfigGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *KubeletConfigPollRequest) StartContext(ctx context.Context) (response *KubeletConfigPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &KubeletConfigPollResponse{
			response: result.(*KubeletConfigGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *KubeletConfigPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// KubeletConfigPollResponse is the response for the Poll method.
type KubeletConfigPollResponse struct {
	response *KubeletConfigGetResponse
}

// Status returns the response status code.
func (r *KubeletConfigPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *KubeletConfigPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *KubeletConfigPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *KubeletConfigPollResponse) Body() *KubeletConfig {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *KubeletConfigPollResponse) GetBody() (value *KubeletConfig, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *KubeletConfigClient) Poll() *KubeletConfigPollRequest {
	return &KubeletConfigPollRequest{
		request: c.Get(),
	}
}

// KubeletConfigDeleteRequest is the request for the 'delete' method.
type KubeletConfigDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *KubeletConfigDeleteRequest) Parameter(name string, value interface{}) *KubeletConfigDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *KubeletConfigDeleteRequest) Header(name string, value interface{}) *KubeletConfigDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *KubeletConfigDeleteRequest) Impersonate(user string) *KubeletConfigDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *KubeletConfigDeleteRequest) Send() (result *KubeletConfigDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *KubeletConfigDeleteRequest) SendContext(ctx context.Context) (result *KubeletConfigDeleteResponse, err error) {
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
	result = &KubeletConfigDeleteResponse{}
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

// KubeletConfigDeleteResponse is the response for the 'delete' method.
type KubeletConfigDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *KubeletConfigDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *KubeletConfigDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *KubeletConfigDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// KubeletConfigGetRequest is the request for the 'get' method.
type KubeletConfigGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *KubeletConfigGetRequest) Parameter(name string, value interface{}) *KubeletConfigGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *KubeletConfigGetRequest) Header(name string, value interface{}) *KubeletConfigGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *KubeletConfigGetRequest) Impersonate(user string) *KubeletConfigGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *KubeletConfigGetRequest) Send() (result *KubeletConfigGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *KubeletConfigGetRequest) SendContext(ctx context.Context) (result *KubeletConfigGetResponse, err error) {
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
	result = &KubeletConfigGetResponse{}
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
	err = readKubeletConfigGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// KubeletConfigGetResponse is the response for the 'get' method.
type KubeletConfigGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *KubeletConfig
}

// Status returns the response status code.
func (r *KubeletConfigGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *KubeletConfigGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *KubeletConfigGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *KubeletConfigGetResponse) Body() *KubeletConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *KubeletConfigGetResponse) GetBody() (value *KubeletConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// KubeletConfigPostRequest is the request for the 'post' method.
type KubeletConfigPostRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *KubeletConfig
}

// Parameter adds a query parameter.
func (r *KubeletConfigPostRequest) Parameter(name string, value interface{}) *KubeletConfigPostRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *KubeletConfigPostRequest) Header(name string, value interface{}) *KubeletConfigPostRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *KubeletConfigPostRequest) Impersonate(user string) *KubeletConfigPostRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *KubeletConfigPostRequest) Body(value *KubeletConfig) *KubeletConfigPostRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *KubeletConfigPostRequest) Send() (result *KubeletConfigPostResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *KubeletConfigPostRequest) SendContext(ctx context.Context) (result *KubeletConfigPostResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeKubeletConfigPostRequest(r, buffer)
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
	result = &KubeletConfigPostResponse{}
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
	err = readKubeletConfigPostResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// KubeletConfigPostResponse is the response for the 'post' method.
type KubeletConfigPostResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *KubeletConfig
}

// Status returns the response status code.
func (r *KubeletConfigPostResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *KubeletConfigPostResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *KubeletConfigPostResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *KubeletConfigPostResponse) Body() *KubeletConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *KubeletConfigPostResponse) GetBody() (value *KubeletConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// KubeletConfigUpdateRequest is the request for the 'update' method.
type KubeletConfigUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *KubeletConfig
}

// Parameter adds a query parameter.
func (r *KubeletConfigUpdateRequest) Parameter(name string, value interface{}) *KubeletConfigUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *KubeletConfigUpdateRequest) Header(name string, value interface{}) *KubeletConfigUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *KubeletConfigUpdateRequest) Impersonate(user string) *KubeletConfigUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *KubeletConfigUpdateRequest) Body(value *KubeletConfig) *KubeletConfigUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *KubeletConfigUpdateRequest) Send() (result *KubeletConfigUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *KubeletConfigUpdateRequest) SendContext(ctx context.Context) (result *KubeletConfigUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeKubeletConfigUpdateRequest(r, buffer)
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
	result = &KubeletConfigUpdateResponse{}
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
	err = readKubeletConfigUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// KubeletConfigUpdateResponse is the response for the 'update' method.
type KubeletConfigUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *KubeletConfig
}

// Status returns the response status code.
func (r *KubeletConfigUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *KubeletConfigUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *KubeletConfigUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *KubeletConfigUpdateResponse) Body() *KubeletConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *KubeletConfigUpdateResponse) GetBody() (value *KubeletConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
