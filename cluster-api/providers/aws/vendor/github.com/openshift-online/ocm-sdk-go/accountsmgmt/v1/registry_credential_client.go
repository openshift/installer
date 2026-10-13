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

// RegistryCredentialClient is the client of the 'registry_credential' resource.
//
// Manages a specific registry credential.
type RegistryCredentialClient struct {
	transport http.RoundTripper
	path      string
}

// NewRegistryCredentialClient creates a new client for the 'registry_credential'
// resource using the given transport to send the requests and receive the
// responses.
func NewRegistryCredentialClient(transport http.RoundTripper, path string) *RegistryCredentialClient {
	return &RegistryCredentialClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Delete the registry credential
func (c *RegistryCredentialClient) Delete() *RegistryCredentialDeleteRequest {
	return &RegistryCredentialDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the registry credential.
func (c *RegistryCredentialClient) Get() *RegistryCredentialGetRequest {
	return &RegistryCredentialGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// RegistryCredentialPollRequest is the request for the Poll method.
type RegistryCredentialPollRequest struct {
	request    *RegistryCredentialGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *RegistryCredentialPollRequest) Parameter(name string, value interface{}) *RegistryCredentialPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *RegistryCredentialPollRequest) Header(name string, value interface{}) *RegistryCredentialPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *RegistryCredentialPollRequest) Interval(value time.Duration) *RegistryCredentialPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *RegistryCredentialPollRequest) Status(value int) *RegistryCredentialPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *RegistryCredentialPollRequest) Predicate(value func(*RegistryCredentialGetResponse) bool) *RegistryCredentialPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*RegistryCredentialGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *RegistryCredentialPollRequest) StartContext(ctx context.Context) (response *RegistryCredentialPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &RegistryCredentialPollResponse{
			response: result.(*RegistryCredentialGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *RegistryCredentialPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// RegistryCredentialPollResponse is the response for the Poll method.
type RegistryCredentialPollResponse struct {
	response *RegistryCredentialGetResponse
}

// Status returns the response status code.
func (r *RegistryCredentialPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *RegistryCredentialPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *RegistryCredentialPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *RegistryCredentialPollResponse) Body() *RegistryCredential {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *RegistryCredentialPollResponse) GetBody() (value *RegistryCredential, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *RegistryCredentialClient) Poll() *RegistryCredentialPollRequest {
	return &RegistryCredentialPollRequest{
		request: c.Get(),
	}
}

// RegistryCredentialDeleteRequest is the request for the 'delete' method.
type RegistryCredentialDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *RegistryCredentialDeleteRequest) Parameter(name string, value interface{}) *RegistryCredentialDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RegistryCredentialDeleteRequest) Header(name string, value interface{}) *RegistryCredentialDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RegistryCredentialDeleteRequest) Impersonate(user string) *RegistryCredentialDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RegistryCredentialDeleteRequest) Send() (result *RegistryCredentialDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RegistryCredentialDeleteRequest) SendContext(ctx context.Context) (result *RegistryCredentialDeleteResponse, err error) {
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
	result = &RegistryCredentialDeleteResponse{}
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

// RegistryCredentialDeleteResponse is the response for the 'delete' method.
type RegistryCredentialDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *RegistryCredentialDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RegistryCredentialDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RegistryCredentialDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// RegistryCredentialGetRequest is the request for the 'get' method.
type RegistryCredentialGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *RegistryCredentialGetRequest) Parameter(name string, value interface{}) *RegistryCredentialGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *RegistryCredentialGetRequest) Header(name string, value interface{}) *RegistryCredentialGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *RegistryCredentialGetRequest) Impersonate(user string) *RegistryCredentialGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *RegistryCredentialGetRequest) Send() (result *RegistryCredentialGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *RegistryCredentialGetRequest) SendContext(ctx context.Context) (result *RegistryCredentialGetResponse, err error) {
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
	result = &RegistryCredentialGetResponse{}
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
	err = readRegistryCredentialGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// RegistryCredentialGetResponse is the response for the 'get' method.
type RegistryCredentialGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *RegistryCredential
}

// Status returns the response status code.
func (r *RegistryCredentialGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *RegistryCredentialGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *RegistryCredentialGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *RegistryCredentialGetResponse) Body() *RegistryCredential {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *RegistryCredentialGetResponse) GetBody() (value *RegistryCredential, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
