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

// SubscriptionClient is the client of the 'subscription' resource.
//
// Manages a specific subscription.
type SubscriptionClient struct {
	transport http.RoundTripper
	path      string
}

// NewSubscriptionClient creates a new client for the 'subscription'
// resource using the given transport to send the requests and receive the
// responses.
func NewSubscriptionClient(transport http.RoundTripper, path string) *SubscriptionClient {
	return &SubscriptionClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the subscription by ID.
func (c *SubscriptionClient) Delete() *SubscriptionDeleteRequest {
	return &SubscriptionDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the subscription by ID.
func (c *SubscriptionClient) Get() *SubscriptionGetRequest {
	return &SubscriptionGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Update a subscription
func (c *SubscriptionClient) Update() *SubscriptionUpdateRequest {
	return &SubscriptionUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Labels returns the target 'generic_labels' resource.
//
// Reference to the list of labels of a specific subscription.
func (c *SubscriptionClient) Labels() *GenericLabelsClient {
	return NewGenericLabelsClient(
		c.transport,
		path.Join(c.path, "labels"),
	)
}

// ReservedResources returns the target 'subscription_reserved_resources' resource.
//
// Reference to the resource that manages the collection of resources reserved by the
// subscription.
func (c *SubscriptionClient) ReservedResources() *SubscriptionReservedResourcesClient {
	return NewSubscriptionReservedResourcesClient(
		c.transport,
		path.Join(c.path, "reserved_resources"),
	)
}

// RoleBindings returns the target 'role_bindings' resource.
//
// Reference to the role bindings
func (c *SubscriptionClient) RoleBindings() *RoleBindingsClient {
	return NewRoleBindingsClient(
		c.transport,
		path.Join(c.path, "role_bindings"),
	)
}

// SubscriptionPollRequest is the request for the Poll method.
type SubscriptionPollRequest struct {
	request    *SubscriptionGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *SubscriptionPollRequest) Parameter(name string, value interface{}) *SubscriptionPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *SubscriptionPollRequest) Header(name string, value interface{}) *SubscriptionPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *SubscriptionPollRequest) Interval(value time.Duration) *SubscriptionPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *SubscriptionPollRequest) Status(value int) *SubscriptionPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *SubscriptionPollRequest) Predicate(value func(*SubscriptionGetResponse) bool) *SubscriptionPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*SubscriptionGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *SubscriptionPollRequest) StartContext(ctx context.Context) (response *SubscriptionPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &SubscriptionPollResponse{
			response: result.(*SubscriptionGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *SubscriptionPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// SubscriptionPollResponse is the response for the Poll method.
type SubscriptionPollResponse struct {
	response *SubscriptionGetResponse
}

// Status returns the response status code.
func (r *SubscriptionPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *SubscriptionPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *SubscriptionPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *SubscriptionPollResponse) Body() *Subscription {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *SubscriptionPollResponse) GetBody() (value *Subscription, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *SubscriptionClient) Poll() *SubscriptionPollRequest {
	return &SubscriptionPollRequest{
		request: c.Get(),
	}
}

// SubscriptionDeleteRequest is the request for the 'delete' method.
type SubscriptionDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *SubscriptionDeleteRequest) Parameter(name string, value interface{}) *SubscriptionDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *SubscriptionDeleteRequest) Header(name string, value interface{}) *SubscriptionDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *SubscriptionDeleteRequest) Impersonate(user string) *SubscriptionDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *SubscriptionDeleteRequest) Send() (result *SubscriptionDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *SubscriptionDeleteRequest) SendContext(ctx context.Context) (result *SubscriptionDeleteResponse, err error) {
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
	result = &SubscriptionDeleteResponse{}
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

// SubscriptionDeleteResponse is the response for the 'delete' method.
type SubscriptionDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *SubscriptionDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *SubscriptionDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *SubscriptionDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// SubscriptionGetRequest is the request for the 'get' method.
type SubscriptionGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *SubscriptionGetRequest) Parameter(name string, value interface{}) *SubscriptionGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *SubscriptionGetRequest) Header(name string, value interface{}) *SubscriptionGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *SubscriptionGetRequest) Impersonate(user string) *SubscriptionGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *SubscriptionGetRequest) Send() (result *SubscriptionGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *SubscriptionGetRequest) SendContext(ctx context.Context) (result *SubscriptionGetResponse, err error) {
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
	result = &SubscriptionGetResponse{}
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
	err = readSubscriptionGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// SubscriptionGetResponse is the response for the 'get' method.
type SubscriptionGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Subscription
}

// Status returns the response status code.
func (r *SubscriptionGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *SubscriptionGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *SubscriptionGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *SubscriptionGetResponse) Body() *Subscription {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *SubscriptionGetResponse) GetBody() (value *Subscription, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// SubscriptionUpdateRequest is the request for the 'update' method.
type SubscriptionUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Subscription
}

// Parameter adds a query parameter.
func (r *SubscriptionUpdateRequest) Parameter(name string, value interface{}) *SubscriptionUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *SubscriptionUpdateRequest) Header(name string, value interface{}) *SubscriptionUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *SubscriptionUpdateRequest) Impersonate(user string) *SubscriptionUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Updated subscription data
func (r *SubscriptionUpdateRequest) Body(value *Subscription) *SubscriptionUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *SubscriptionUpdateRequest) Send() (result *SubscriptionUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *SubscriptionUpdateRequest) SendContext(ctx context.Context) (result *SubscriptionUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeSubscriptionUpdateRequest(r, buffer)
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
	result = &SubscriptionUpdateResponse{}
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
	err = readSubscriptionUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// SubscriptionUpdateResponse is the response for the 'update' method.
type SubscriptionUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Subscription
}

// Status returns the response status code.
func (r *SubscriptionUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *SubscriptionUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *SubscriptionUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Updated subscription data
func (r *SubscriptionUpdateResponse) Body() *Subscription {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Updated subscription data
func (r *SubscriptionUpdateResponse) GetBody() (value *Subscription, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
