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

// OidcConfigClient is the client of the 'oidc_config' resource.
//
// Manages an Oidc Config configuration.
type OidcConfigClient struct {
	transport http.RoundTripper
	path      string
}

// NewOidcConfigClient creates a new client for the 'oidc_config'
// resource using the given transport to send the requests and receive the
// responses.
func NewOidcConfigClient(transport http.RoundTripper, path string) *OidcConfigClient {
	return &OidcConfigClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the OidcConfig.
func (c *OidcConfigClient) Delete() *OidcConfigDeleteRequest {
	return &OidcConfigDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of an OidcConfig.
func (c *OidcConfigClient) Get() *OidcConfigGetRequest {
	return &OidcConfigGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates attributes of an OidcConfig.
func (c *OidcConfigClient) Update() *OidcConfigUpdateRequest {
	return &OidcConfigUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// OidcConfigPollRequest is the request for the Poll method.
type OidcConfigPollRequest struct {
	request    *OidcConfigGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *OidcConfigPollRequest) Parameter(name string, value interface{}) *OidcConfigPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *OidcConfigPollRequest) Header(name string, value interface{}) *OidcConfigPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *OidcConfigPollRequest) Interval(value time.Duration) *OidcConfigPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *OidcConfigPollRequest) Status(value int) *OidcConfigPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *OidcConfigPollRequest) Predicate(value func(*OidcConfigGetResponse) bool) *OidcConfigPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*OidcConfigGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *OidcConfigPollRequest) StartContext(ctx context.Context) (response *OidcConfigPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &OidcConfigPollResponse{
			response: result.(*OidcConfigGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *OidcConfigPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// OidcConfigPollResponse is the response for the Poll method.
type OidcConfigPollResponse struct {
	response *OidcConfigGetResponse
}

// Status returns the response status code.
func (r *OidcConfigPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *OidcConfigPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *OidcConfigPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *OidcConfigPollResponse) Body() *OidcConfig {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *OidcConfigPollResponse) GetBody() (value *OidcConfig, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *OidcConfigClient) Poll() *OidcConfigPollRequest {
	return &OidcConfigPollRequest{
		request: c.Get(),
	}
}

// OidcConfigDeleteRequest is the request for the 'delete' method.
type OidcConfigDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *OidcConfigDeleteRequest) Parameter(name string, value interface{}) *OidcConfigDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OidcConfigDeleteRequest) Header(name string, value interface{}) *OidcConfigDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OidcConfigDeleteRequest) Impersonate(user string) *OidcConfigDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OidcConfigDeleteRequest) Send() (result *OidcConfigDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OidcConfigDeleteRequest) SendContext(ctx context.Context) (result *OidcConfigDeleteResponse, err error) {
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
	result = &OidcConfigDeleteResponse{}
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

// OidcConfigDeleteResponse is the response for the 'delete' method.
type OidcConfigDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *OidcConfigDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OidcConfigDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OidcConfigDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// OidcConfigGetRequest is the request for the 'get' method.
type OidcConfigGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *OidcConfigGetRequest) Parameter(name string, value interface{}) *OidcConfigGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OidcConfigGetRequest) Header(name string, value interface{}) *OidcConfigGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OidcConfigGetRequest) Impersonate(user string) *OidcConfigGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OidcConfigGetRequest) Send() (result *OidcConfigGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OidcConfigGetRequest) SendContext(ctx context.Context) (result *OidcConfigGetResponse, err error) {
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
	result = &OidcConfigGetResponse{}
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
	err = readOidcConfigGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// OidcConfigGetResponse is the response for the 'get' method.
type OidcConfigGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *OidcConfig
}

// Status returns the response status code.
func (r *OidcConfigGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OidcConfigGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OidcConfigGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *OidcConfigGetResponse) Body() *OidcConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *OidcConfigGetResponse) GetBody() (value *OidcConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// OidcConfigUpdateRequest is the request for the 'update' method.
type OidcConfigUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *OidcConfig
}

// Parameter adds a query parameter.
func (r *OidcConfigUpdateRequest) Parameter(name string, value interface{}) *OidcConfigUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OidcConfigUpdateRequest) Header(name string, value interface{}) *OidcConfigUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OidcConfigUpdateRequest) Impersonate(user string) *OidcConfigUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *OidcConfigUpdateRequest) Body(value *OidcConfig) *OidcConfigUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OidcConfigUpdateRequest) Send() (result *OidcConfigUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OidcConfigUpdateRequest) SendContext(ctx context.Context) (result *OidcConfigUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeOidcConfigUpdateRequest(r, buffer)
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
	result = &OidcConfigUpdateResponse{}
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
	err = readOidcConfigUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// OidcConfigUpdateResponse is the response for the 'update' method.
type OidcConfigUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *OidcConfig
}

// Status returns the response status code.
func (r *OidcConfigUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OidcConfigUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OidcConfigUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *OidcConfigUpdateResponse) Body() *OidcConfig {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *OidcConfigUpdateResponse) GetBody() (value *OidcConfig, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
