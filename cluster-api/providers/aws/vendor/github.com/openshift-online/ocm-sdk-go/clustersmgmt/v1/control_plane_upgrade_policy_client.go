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

// ControlPlaneUpgradePolicyClient is the client of the 'control_plane_upgrade_policy' resource.
//
// Manages a specific upgrade policy for the control plane.
type ControlPlaneUpgradePolicyClient struct {
	transport http.RoundTripper
	path      string
}

// NewControlPlaneUpgradePolicyClient creates a new client for the 'control_plane_upgrade_policy'
// resource using the given transport to send the requests and receive the
// responses.
func NewControlPlaneUpgradePolicyClient(transport http.RoundTripper, path string) *ControlPlaneUpgradePolicyClient {
	return &ControlPlaneUpgradePolicyClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the upgrade policy for the control plane.
func (c *ControlPlaneUpgradePolicyClient) Delete() *ControlPlaneUpgradePolicyDeleteRequest {
	return &ControlPlaneUpgradePolicyDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the upgrade policy for the control plane.
func (c *ControlPlaneUpgradePolicyClient) Get() *ControlPlaneUpgradePolicyGetRequest {
	return &ControlPlaneUpgradePolicyGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Update the upgrade policy for the control plane.
func (c *ControlPlaneUpgradePolicyClient) Update() *ControlPlaneUpgradePolicyUpdateRequest {
	return &ControlPlaneUpgradePolicyUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ControlPlaneUpgradePolicyPollRequest is the request for the Poll method.
type ControlPlaneUpgradePolicyPollRequest struct {
	request    *ControlPlaneUpgradePolicyGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ControlPlaneUpgradePolicyPollRequest) Parameter(name string, value interface{}) *ControlPlaneUpgradePolicyPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ControlPlaneUpgradePolicyPollRequest) Header(name string, value interface{}) *ControlPlaneUpgradePolicyPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ControlPlaneUpgradePolicyPollRequest) Interval(value time.Duration) *ControlPlaneUpgradePolicyPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ControlPlaneUpgradePolicyPollRequest) Status(value int) *ControlPlaneUpgradePolicyPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ControlPlaneUpgradePolicyPollRequest) Predicate(value func(*ControlPlaneUpgradePolicyGetResponse) bool) *ControlPlaneUpgradePolicyPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ControlPlaneUpgradePolicyGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ControlPlaneUpgradePolicyPollRequest) StartContext(ctx context.Context) (response *ControlPlaneUpgradePolicyPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ControlPlaneUpgradePolicyPollResponse{
			response: result.(*ControlPlaneUpgradePolicyGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ControlPlaneUpgradePolicyPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ControlPlaneUpgradePolicyPollResponse is the response for the Poll method.
type ControlPlaneUpgradePolicyPollResponse struct {
	response *ControlPlaneUpgradePolicyGetResponse
}

// Status returns the response status code.
func (r *ControlPlaneUpgradePolicyPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ControlPlaneUpgradePolicyPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ControlPlaneUpgradePolicyPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ControlPlaneUpgradePolicyPollResponse) Body() *ControlPlaneUpgradePolicy {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ControlPlaneUpgradePolicyPollResponse) GetBody() (value *ControlPlaneUpgradePolicy, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ControlPlaneUpgradePolicyClient) Poll() *ControlPlaneUpgradePolicyPollRequest {
	return &ControlPlaneUpgradePolicyPollRequest{
		request: c.Get(),
	}
}

// ControlPlaneUpgradePolicyDeleteRequest is the request for the 'delete' method.
type ControlPlaneUpgradePolicyDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ControlPlaneUpgradePolicyDeleteRequest) Parameter(name string, value interface{}) *ControlPlaneUpgradePolicyDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ControlPlaneUpgradePolicyDeleteRequest) Header(name string, value interface{}) *ControlPlaneUpgradePolicyDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ControlPlaneUpgradePolicyDeleteRequest) Impersonate(user string) *ControlPlaneUpgradePolicyDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ControlPlaneUpgradePolicyDeleteRequest) Send() (result *ControlPlaneUpgradePolicyDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ControlPlaneUpgradePolicyDeleteRequest) SendContext(ctx context.Context) (result *ControlPlaneUpgradePolicyDeleteResponse, err error) {
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
	result = &ControlPlaneUpgradePolicyDeleteResponse{}
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

// ControlPlaneUpgradePolicyDeleteResponse is the response for the 'delete' method.
type ControlPlaneUpgradePolicyDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *ControlPlaneUpgradePolicyDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ControlPlaneUpgradePolicyDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ControlPlaneUpgradePolicyDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// ControlPlaneUpgradePolicyGetRequest is the request for the 'get' method.
type ControlPlaneUpgradePolicyGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ControlPlaneUpgradePolicyGetRequest) Parameter(name string, value interface{}) *ControlPlaneUpgradePolicyGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ControlPlaneUpgradePolicyGetRequest) Header(name string, value interface{}) *ControlPlaneUpgradePolicyGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ControlPlaneUpgradePolicyGetRequest) Impersonate(user string) *ControlPlaneUpgradePolicyGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ControlPlaneUpgradePolicyGetRequest) Send() (result *ControlPlaneUpgradePolicyGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ControlPlaneUpgradePolicyGetRequest) SendContext(ctx context.Context) (result *ControlPlaneUpgradePolicyGetResponse, err error) {
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
	result = &ControlPlaneUpgradePolicyGetResponse{}
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
	err = readControlPlaneUpgradePolicyGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ControlPlaneUpgradePolicyGetResponse is the response for the 'get' method.
type ControlPlaneUpgradePolicyGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ControlPlaneUpgradePolicy
}

// Status returns the response status code.
func (r *ControlPlaneUpgradePolicyGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ControlPlaneUpgradePolicyGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ControlPlaneUpgradePolicyGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ControlPlaneUpgradePolicyGetResponse) Body() *ControlPlaneUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ControlPlaneUpgradePolicyGetResponse) GetBody() (value *ControlPlaneUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ControlPlaneUpgradePolicyUpdateRequest is the request for the 'update' method.
type ControlPlaneUpgradePolicyUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ControlPlaneUpgradePolicy
}

// Parameter adds a query parameter.
func (r *ControlPlaneUpgradePolicyUpdateRequest) Parameter(name string, value interface{}) *ControlPlaneUpgradePolicyUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ControlPlaneUpgradePolicyUpdateRequest) Header(name string, value interface{}) *ControlPlaneUpgradePolicyUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ControlPlaneUpgradePolicyUpdateRequest) Impersonate(user string) *ControlPlaneUpgradePolicyUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ControlPlaneUpgradePolicyUpdateRequest) Body(value *ControlPlaneUpgradePolicy) *ControlPlaneUpgradePolicyUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ControlPlaneUpgradePolicyUpdateRequest) Send() (result *ControlPlaneUpgradePolicyUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ControlPlaneUpgradePolicyUpdateRequest) SendContext(ctx context.Context) (result *ControlPlaneUpgradePolicyUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeControlPlaneUpgradePolicyUpdateRequest(r, buffer)
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
	result = &ControlPlaneUpgradePolicyUpdateResponse{}
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
	err = readControlPlaneUpgradePolicyUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ControlPlaneUpgradePolicyUpdateResponse is the response for the 'update' method.
type ControlPlaneUpgradePolicyUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ControlPlaneUpgradePolicy
}

// Status returns the response status code.
func (r *ControlPlaneUpgradePolicyUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ControlPlaneUpgradePolicyUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ControlPlaneUpgradePolicyUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ControlPlaneUpgradePolicyUpdateResponse) Body() *ControlPlaneUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ControlPlaneUpgradePolicyUpdateResponse) GetBody() (value *ControlPlaneUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
