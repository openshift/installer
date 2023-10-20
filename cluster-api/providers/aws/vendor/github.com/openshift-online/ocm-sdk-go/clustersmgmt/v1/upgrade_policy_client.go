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
	"path"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// UpgradePolicyClient is the client of the 'upgrade_policy' resource.
//
// Manages a specific upgrade policy.
type UpgradePolicyClient struct {
	transport http.RoundTripper
	path      string
}

// NewUpgradePolicyClient creates a new client for the 'upgrade_policy'
// resource using the given transport to send the requests and receive the
// responses.
func NewUpgradePolicyClient(transport http.RoundTripper, path string) *UpgradePolicyClient {
	return &UpgradePolicyClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the upgrade policy.
func (c *UpgradePolicyClient) Delete() *UpgradePolicyDeleteRequest {
	return &UpgradePolicyDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the upgrade policy.
func (c *UpgradePolicyClient) Get() *UpgradePolicyGetRequest {
	return &UpgradePolicyGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Update the upgrade policy.
func (c *UpgradePolicyClient) Update() *UpgradePolicyUpdateRequest {
	return &UpgradePolicyUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// State returns the target 'upgrade_policy_state' resource.
//
// Reference to the state of the upgrade policy.
func (c *UpgradePolicyClient) State() *UpgradePolicyStateClient {
	return NewUpgradePolicyStateClient(
		c.transport,
		path.Join(c.path, "state"),
	)
}

// UpgradePolicyPollRequest is the request for the Poll method.
type UpgradePolicyPollRequest struct {
	request    *UpgradePolicyGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *UpgradePolicyPollRequest) Parameter(name string, value interface{}) *UpgradePolicyPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *UpgradePolicyPollRequest) Header(name string, value interface{}) *UpgradePolicyPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *UpgradePolicyPollRequest) Interval(value time.Duration) *UpgradePolicyPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *UpgradePolicyPollRequest) Status(value int) *UpgradePolicyPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *UpgradePolicyPollRequest) Predicate(value func(*UpgradePolicyGetResponse) bool) *UpgradePolicyPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*UpgradePolicyGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *UpgradePolicyPollRequest) StartContext(ctx context.Context) (response *UpgradePolicyPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &UpgradePolicyPollResponse{
			response: result.(*UpgradePolicyGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *UpgradePolicyPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// UpgradePolicyPollResponse is the response for the Poll method.
type UpgradePolicyPollResponse struct {
	response *UpgradePolicyGetResponse
}

// Status returns the response status code.
func (r *UpgradePolicyPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *UpgradePolicyPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *UpgradePolicyPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *UpgradePolicyPollResponse) Body() *UpgradePolicy {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *UpgradePolicyPollResponse) GetBody() (value *UpgradePolicy, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *UpgradePolicyClient) Poll() *UpgradePolicyPollRequest {
	return &UpgradePolicyPollRequest{
		request: c.Get(),
	}
}

// UpgradePolicyDeleteRequest is the request for the 'delete' method.
type UpgradePolicyDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *UpgradePolicyDeleteRequest) Parameter(name string, value interface{}) *UpgradePolicyDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *UpgradePolicyDeleteRequest) Header(name string, value interface{}) *UpgradePolicyDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *UpgradePolicyDeleteRequest) Impersonate(user string) *UpgradePolicyDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *UpgradePolicyDeleteRequest) Send() (result *UpgradePolicyDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *UpgradePolicyDeleteRequest) SendContext(ctx context.Context) (result *UpgradePolicyDeleteResponse, err error) {
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
	result = &UpgradePolicyDeleteResponse{}
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

// UpgradePolicyDeleteResponse is the response for the 'delete' method.
type UpgradePolicyDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *UpgradePolicyDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *UpgradePolicyDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *UpgradePolicyDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// UpgradePolicyGetRequest is the request for the 'get' method.
type UpgradePolicyGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *UpgradePolicyGetRequest) Parameter(name string, value interface{}) *UpgradePolicyGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *UpgradePolicyGetRequest) Header(name string, value interface{}) *UpgradePolicyGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *UpgradePolicyGetRequest) Impersonate(user string) *UpgradePolicyGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *UpgradePolicyGetRequest) Send() (result *UpgradePolicyGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *UpgradePolicyGetRequest) SendContext(ctx context.Context) (result *UpgradePolicyGetResponse, err error) {
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
	result = &UpgradePolicyGetResponse{}
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
	err = readUpgradePolicyGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// UpgradePolicyGetResponse is the response for the 'get' method.
type UpgradePolicyGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *UpgradePolicy
}

// Status returns the response status code.
func (r *UpgradePolicyGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *UpgradePolicyGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *UpgradePolicyGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *UpgradePolicyGetResponse) Body() *UpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *UpgradePolicyGetResponse) GetBody() (value *UpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// UpgradePolicyUpdateRequest is the request for the 'update' method.
type UpgradePolicyUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *UpgradePolicy
}

// Parameter adds a query parameter.
func (r *UpgradePolicyUpdateRequest) Parameter(name string, value interface{}) *UpgradePolicyUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *UpgradePolicyUpdateRequest) Header(name string, value interface{}) *UpgradePolicyUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *UpgradePolicyUpdateRequest) Impersonate(user string) *UpgradePolicyUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *UpgradePolicyUpdateRequest) Body(value *UpgradePolicy) *UpgradePolicyUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *UpgradePolicyUpdateRequest) Send() (result *UpgradePolicyUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *UpgradePolicyUpdateRequest) SendContext(ctx context.Context) (result *UpgradePolicyUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeUpgradePolicyUpdateRequest(r, buffer)
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
	result = &UpgradePolicyUpdateResponse{}
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
	err = readUpgradePolicyUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// UpgradePolicyUpdateResponse is the response for the 'update' method.
type UpgradePolicyUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *UpgradePolicy
}

// Status returns the response status code.
func (r *UpgradePolicyUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *UpgradePolicyUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *UpgradePolicyUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *UpgradePolicyUpdateResponse) Body() *UpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *UpgradePolicyUpdateResponse) GetBody() (value *UpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
