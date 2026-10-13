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

// AWSInfrastructureAccessRoleGrantClient is the client of the 'AWS_infrastructure_access_role_grant' resource.
//
// Manages a specific AWS infrastructure access role grant.
type AWSInfrastructureAccessRoleGrantClient struct {
	transport http.RoundTripper
	path      string
}

// NewAWSInfrastructureAccessRoleGrantClient creates a new client for the 'AWS_infrastructure_access_role_grant'
// resource using the given transport to send the requests and receive the
// responses.
func NewAWSInfrastructureAccessRoleGrantClient(transport http.RoundTripper, path string) *AWSInfrastructureAccessRoleGrantClient {
	return &AWSInfrastructureAccessRoleGrantClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the AWS infrastructure access role grant.
func (c *AWSInfrastructureAccessRoleGrantClient) Delete() *AWSInfrastructureAccessRoleGrantDeleteRequest {
	return &AWSInfrastructureAccessRoleGrantDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the AWS infrastructure access role grant.
func (c *AWSInfrastructureAccessRoleGrantClient) Get() *AWSInfrastructureAccessRoleGrantGetRequest {
	return &AWSInfrastructureAccessRoleGrantGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AWSInfrastructureAccessRoleGrantPollRequest is the request for the Poll method.
type AWSInfrastructureAccessRoleGrantPollRequest struct {
	request    *AWSInfrastructureAccessRoleGrantGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AWSInfrastructureAccessRoleGrantPollRequest) Parameter(name string, value interface{}) *AWSInfrastructureAccessRoleGrantPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AWSInfrastructureAccessRoleGrantPollRequest) Header(name string, value interface{}) *AWSInfrastructureAccessRoleGrantPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AWSInfrastructureAccessRoleGrantPollRequest) Interval(value time.Duration) *AWSInfrastructureAccessRoleGrantPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AWSInfrastructureAccessRoleGrantPollRequest) Status(value int) *AWSInfrastructureAccessRoleGrantPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AWSInfrastructureAccessRoleGrantPollRequest) Predicate(value func(*AWSInfrastructureAccessRoleGrantGetResponse) bool) *AWSInfrastructureAccessRoleGrantPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AWSInfrastructureAccessRoleGrantGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AWSInfrastructureAccessRoleGrantPollRequest) StartContext(ctx context.Context) (response *AWSInfrastructureAccessRoleGrantPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AWSInfrastructureAccessRoleGrantPollResponse{
			response: result.(*AWSInfrastructureAccessRoleGrantGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AWSInfrastructureAccessRoleGrantPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AWSInfrastructureAccessRoleGrantPollResponse is the response for the Poll method.
type AWSInfrastructureAccessRoleGrantPollResponse struct {
	response *AWSInfrastructureAccessRoleGrantGetResponse
}

// Status returns the response status code.
func (r *AWSInfrastructureAccessRoleGrantPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AWSInfrastructureAccessRoleGrantPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AWSInfrastructureAccessRoleGrantPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AWSInfrastructureAccessRoleGrantPollResponse) Body() *AWSInfrastructureAccessRoleGrant {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AWSInfrastructureAccessRoleGrantPollResponse) GetBody() (value *AWSInfrastructureAccessRoleGrant, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AWSInfrastructureAccessRoleGrantClient) Poll() *AWSInfrastructureAccessRoleGrantPollRequest {
	return &AWSInfrastructureAccessRoleGrantPollRequest{
		request: c.Get(),
	}
}

// AWSInfrastructureAccessRoleGrantDeleteRequest is the request for the 'delete' method.
type AWSInfrastructureAccessRoleGrantDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AWSInfrastructureAccessRoleGrantDeleteRequest) Parameter(name string, value interface{}) *AWSInfrastructureAccessRoleGrantDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AWSInfrastructureAccessRoleGrantDeleteRequest) Header(name string, value interface{}) *AWSInfrastructureAccessRoleGrantDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AWSInfrastructureAccessRoleGrantDeleteRequest) Impersonate(user string) *AWSInfrastructureAccessRoleGrantDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AWSInfrastructureAccessRoleGrantDeleteRequest) Send() (result *AWSInfrastructureAccessRoleGrantDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AWSInfrastructureAccessRoleGrantDeleteRequest) SendContext(ctx context.Context) (result *AWSInfrastructureAccessRoleGrantDeleteResponse, err error) {
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
	result = &AWSInfrastructureAccessRoleGrantDeleteResponse{}
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

// AWSInfrastructureAccessRoleGrantDeleteResponse is the response for the 'delete' method.
type AWSInfrastructureAccessRoleGrantDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *AWSInfrastructureAccessRoleGrantDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AWSInfrastructureAccessRoleGrantDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AWSInfrastructureAccessRoleGrantDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// AWSInfrastructureAccessRoleGrantGetRequest is the request for the 'get' method.
type AWSInfrastructureAccessRoleGrantGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AWSInfrastructureAccessRoleGrantGetRequest) Parameter(name string, value interface{}) *AWSInfrastructureAccessRoleGrantGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AWSInfrastructureAccessRoleGrantGetRequest) Header(name string, value interface{}) *AWSInfrastructureAccessRoleGrantGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AWSInfrastructureAccessRoleGrantGetRequest) Impersonate(user string) *AWSInfrastructureAccessRoleGrantGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AWSInfrastructureAccessRoleGrantGetRequest) Send() (result *AWSInfrastructureAccessRoleGrantGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AWSInfrastructureAccessRoleGrantGetRequest) SendContext(ctx context.Context) (result *AWSInfrastructureAccessRoleGrantGetResponse, err error) {
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
	result = &AWSInfrastructureAccessRoleGrantGetResponse{}
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
	err = readAWSInfrastructureAccessRoleGrantGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AWSInfrastructureAccessRoleGrantGetResponse is the response for the 'get' method.
type AWSInfrastructureAccessRoleGrantGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AWSInfrastructureAccessRoleGrant
}

// Status returns the response status code.
func (r *AWSInfrastructureAccessRoleGrantGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AWSInfrastructureAccessRoleGrantGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AWSInfrastructureAccessRoleGrantGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AWSInfrastructureAccessRoleGrantGetResponse) Body() *AWSInfrastructureAccessRoleGrant {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AWSInfrastructureAccessRoleGrantGetResponse) GetBody() (value *AWSInfrastructureAccessRoleGrant, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
