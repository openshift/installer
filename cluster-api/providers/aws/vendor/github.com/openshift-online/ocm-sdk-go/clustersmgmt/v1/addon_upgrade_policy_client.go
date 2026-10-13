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

// AddonUpgradePolicyClient is the client of the 'addon_upgrade_policy' resource.
//
// Manages a specific addon upgrade policy.
type AddonUpgradePolicyClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddonUpgradePolicyClient creates a new client for the 'addon_upgrade_policy'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddonUpgradePolicyClient(transport http.RoundTripper, path string) *AddonUpgradePolicyClient {
	return &AddonUpgradePolicyClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the addon upgrade policy.
func (c *AddonUpgradePolicyClient) Delete() *AddonUpgradePolicyDeleteRequest {
	return &AddonUpgradePolicyDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the addon upgrade policy.
func (c *AddonUpgradePolicyClient) Get() *AddonUpgradePolicyGetRequest {
	return &AddonUpgradePolicyGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Update the addon upgrade policy.
func (c *AddonUpgradePolicyClient) Update() *AddonUpgradePolicyUpdateRequest {
	return &AddonUpgradePolicyUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// State returns the target 'addon_upgrade_policy_state' resource.
//
// Reference to the state of the addon upgrade policy.
func (c *AddonUpgradePolicyClient) State() *AddonUpgradePolicyStateClient {
	return NewAddonUpgradePolicyStateClient(
		c.transport,
		path.Join(c.path, "state"),
	)
}

// AddonUpgradePolicyPollRequest is the request for the Poll method.
type AddonUpgradePolicyPollRequest struct {
	request    *AddonUpgradePolicyGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AddonUpgradePolicyPollRequest) Parameter(name string, value interface{}) *AddonUpgradePolicyPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AddonUpgradePolicyPollRequest) Header(name string, value interface{}) *AddonUpgradePolicyPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AddonUpgradePolicyPollRequest) Interval(value time.Duration) *AddonUpgradePolicyPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AddonUpgradePolicyPollRequest) Status(value int) *AddonUpgradePolicyPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AddonUpgradePolicyPollRequest) Predicate(value func(*AddonUpgradePolicyGetResponse) bool) *AddonUpgradePolicyPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AddonUpgradePolicyGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AddonUpgradePolicyPollRequest) StartContext(ctx context.Context) (response *AddonUpgradePolicyPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AddonUpgradePolicyPollResponse{
			response: result.(*AddonUpgradePolicyGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AddonUpgradePolicyPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AddonUpgradePolicyPollResponse is the response for the Poll method.
type AddonUpgradePolicyPollResponse struct {
	response *AddonUpgradePolicyGetResponse
}

// Status returns the response status code.
func (r *AddonUpgradePolicyPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AddonUpgradePolicyPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AddonUpgradePolicyPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AddonUpgradePolicyPollResponse) Body() *AddonUpgradePolicy {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonUpgradePolicyPollResponse) GetBody() (value *AddonUpgradePolicy, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AddonUpgradePolicyClient) Poll() *AddonUpgradePolicyPollRequest {
	return &AddonUpgradePolicyPollRequest{
		request: c.Get(),
	}
}

// AddonUpgradePolicyDeleteRequest is the request for the 'delete' method.
type AddonUpgradePolicyDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddonUpgradePolicyDeleteRequest) Parameter(name string, value interface{}) *AddonUpgradePolicyDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonUpgradePolicyDeleteRequest) Header(name string, value interface{}) *AddonUpgradePolicyDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonUpgradePolicyDeleteRequest) Impersonate(user string) *AddonUpgradePolicyDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonUpgradePolicyDeleteRequest) Send() (result *AddonUpgradePolicyDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonUpgradePolicyDeleteRequest) SendContext(ctx context.Context) (result *AddonUpgradePolicyDeleteResponse, err error) {
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
	result = &AddonUpgradePolicyDeleteResponse{}
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

// AddonUpgradePolicyDeleteResponse is the response for the 'delete' method.
type AddonUpgradePolicyDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *AddonUpgradePolicyDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonUpgradePolicyDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonUpgradePolicyDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// AddonUpgradePolicyGetRequest is the request for the 'get' method.
type AddonUpgradePolicyGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddonUpgradePolicyGetRequest) Parameter(name string, value interface{}) *AddonUpgradePolicyGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonUpgradePolicyGetRequest) Header(name string, value interface{}) *AddonUpgradePolicyGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonUpgradePolicyGetRequest) Impersonate(user string) *AddonUpgradePolicyGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonUpgradePolicyGetRequest) Send() (result *AddonUpgradePolicyGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonUpgradePolicyGetRequest) SendContext(ctx context.Context) (result *AddonUpgradePolicyGetResponse, err error) {
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
	result = &AddonUpgradePolicyGetResponse{}
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
	err = readAddonUpgradePolicyGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonUpgradePolicyGetResponse is the response for the 'get' method.
type AddonUpgradePolicyGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddonUpgradePolicy
}

// Status returns the response status code.
func (r *AddonUpgradePolicyGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonUpgradePolicyGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonUpgradePolicyGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddonUpgradePolicyGetResponse) Body() *AddonUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonUpgradePolicyGetResponse) GetBody() (value *AddonUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddonUpgradePolicyUpdateRequest is the request for the 'update' method.
type AddonUpgradePolicyUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddonUpgradePolicy
}

// Parameter adds a query parameter.
func (r *AddonUpgradePolicyUpdateRequest) Parameter(name string, value interface{}) *AddonUpgradePolicyUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonUpgradePolicyUpdateRequest) Header(name string, value interface{}) *AddonUpgradePolicyUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonUpgradePolicyUpdateRequest) Impersonate(user string) *AddonUpgradePolicyUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *AddonUpgradePolicyUpdateRequest) Body(value *AddonUpgradePolicy) *AddonUpgradePolicyUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonUpgradePolicyUpdateRequest) Send() (result *AddonUpgradePolicyUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonUpgradePolicyUpdateRequest) SendContext(ctx context.Context) (result *AddonUpgradePolicyUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddonUpgradePolicyUpdateRequest(r, buffer)
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
	result = &AddonUpgradePolicyUpdateResponse{}
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
	err = readAddonUpgradePolicyUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonUpgradePolicyUpdateResponse is the response for the 'update' method.
type AddonUpgradePolicyUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddonUpgradePolicy
}

// Status returns the response status code.
func (r *AddonUpgradePolicyUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonUpgradePolicyUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonUpgradePolicyUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddonUpgradePolicyUpdateResponse) Body() *AddonUpgradePolicy {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonUpgradePolicyUpdateResponse) GetBody() (value *AddonUpgradePolicy, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
