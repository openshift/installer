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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

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

// AddonInstallationClient is the client of the 'addon_installation' resource.
//
// Manages a specific addon installation.
type AddonInstallationClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddonInstallationClient creates a new client for the 'addon_installation'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddonInstallationClient(transport http.RoundTripper, path string) *AddonInstallationClient {
	return &AddonInstallationClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the addon installation.
func (c *AddonInstallationClient) Delete() *AddonInstallationDeleteRequest {
	return &AddonInstallationDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the addon installation.
func (c *AddonInstallationClient) Get() *AddonInstallationGetRequest {
	return &AddonInstallationGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the addon installation.
func (c *AddonInstallationClient) Update() *AddonInstallationUpdateRequest {
	return &AddonInstallationUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AddonInstallationPollRequest is the request for the Poll method.
type AddonInstallationPollRequest struct {
	request    *AddonInstallationGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AddonInstallationPollRequest) Parameter(name string, value interface{}) *AddonInstallationPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AddonInstallationPollRequest) Header(name string, value interface{}) *AddonInstallationPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AddonInstallationPollRequest) Interval(value time.Duration) *AddonInstallationPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AddonInstallationPollRequest) Status(value int) *AddonInstallationPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AddonInstallationPollRequest) Predicate(value func(*AddonInstallationGetResponse) bool) *AddonInstallationPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AddonInstallationGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AddonInstallationPollRequest) StartContext(ctx context.Context) (response *AddonInstallationPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AddonInstallationPollResponse{
			response: result.(*AddonInstallationGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AddonInstallationPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AddonInstallationPollResponse is the response for the Poll method.
type AddonInstallationPollResponse struct {
	response *AddonInstallationGetResponse
}

// Status returns the response status code.
func (r *AddonInstallationPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AddonInstallationPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AddonInstallationPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AddonInstallationPollResponse) Body() *AddonInstallation {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonInstallationPollResponse) GetBody() (value *AddonInstallation, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AddonInstallationClient) Poll() *AddonInstallationPollRequest {
	return &AddonInstallationPollRequest{
		request: c.Get(),
	}
}

// AddonInstallationDeleteRequest is the request for the 'delete' method.
type AddonInstallationDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddonInstallationDeleteRequest) Parameter(name string, value interface{}) *AddonInstallationDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonInstallationDeleteRequest) Header(name string, value interface{}) *AddonInstallationDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonInstallationDeleteRequest) Impersonate(user string) *AddonInstallationDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonInstallationDeleteRequest) Send() (result *AddonInstallationDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonInstallationDeleteRequest) SendContext(ctx context.Context) (result *AddonInstallationDeleteResponse, err error) {
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
	result = &AddonInstallationDeleteResponse{}
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

// AddonInstallationDeleteResponse is the response for the 'delete' method.
type AddonInstallationDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *AddonInstallationDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonInstallationDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonInstallationDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// AddonInstallationGetRequest is the request for the 'get' method.
type AddonInstallationGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddonInstallationGetRequest) Parameter(name string, value interface{}) *AddonInstallationGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonInstallationGetRequest) Header(name string, value interface{}) *AddonInstallationGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonInstallationGetRequest) Impersonate(user string) *AddonInstallationGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonInstallationGetRequest) Send() (result *AddonInstallationGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonInstallationGetRequest) SendContext(ctx context.Context) (result *AddonInstallationGetResponse, err error) {
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
	result = &AddonInstallationGetResponse{}
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
	err = readAddonInstallationGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonInstallationGetResponse is the response for the 'get' method.
type AddonInstallationGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddonInstallation
}

// Status returns the response status code.
func (r *AddonInstallationGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonInstallationGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonInstallationGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddonInstallationGetResponse) Body() *AddonInstallation {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonInstallationGetResponse) GetBody() (value *AddonInstallation, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddonInstallationUpdateRequest is the request for the 'update' method.
type AddonInstallationUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddonInstallation
	dryRun    *bool
}

// Parameter adds a query parameter.
func (r *AddonInstallationUpdateRequest) Parameter(name string, value interface{}) *AddonInstallationUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddonInstallationUpdateRequest) Header(name string, value interface{}) *AddonInstallationUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddonInstallationUpdateRequest) Impersonate(user string) *AddonInstallationUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *AddonInstallationUpdateRequest) Body(value *AddonInstallation) *AddonInstallationUpdateRequest {
	r.body = value
	return r
}

// DryRun sets the value of the 'dry_run' parameter.
//
// DryRun indicates the request body will not be persisted when dryRun=true.
func (r *AddonInstallationUpdateRequest) DryRun(value bool) *AddonInstallationUpdateRequest {
	r.dryRun = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddonInstallationUpdateRequest) Send() (result *AddonInstallationUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddonInstallationUpdateRequest) SendContext(ctx context.Context) (result *AddonInstallationUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.dryRun != nil {
		helpers.AddValue(&query, "dryRun", *r.dryRun)
	}
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddonInstallationUpdateRequest(r, buffer)
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
	result = &AddonInstallationUpdateResponse{}
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
	err = readAddonInstallationUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddonInstallationUpdateResponse is the response for the 'update' method.
type AddonInstallationUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddonInstallation
}

// Status returns the response status code.
func (r *AddonInstallationUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddonInstallationUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddonInstallationUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddonInstallationUpdateResponse) Body() *AddonInstallation {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddonInstallationUpdateResponse) GetBody() (value *AddonInstallation, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
