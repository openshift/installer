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

// AddOnInstallationClient is the client of the 'add_on_installation' resource.
//
// Manages a specific add-on installation.
type AddOnInstallationClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddOnInstallationClient creates a new client for the 'add_on_installation'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddOnInstallationClient(transport http.RoundTripper, path string) *AddOnInstallationClient {
	return &AddOnInstallationClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Delete an add-on installation and remove it from the collection of add-on installations on the cluster.
func (c *AddOnInstallationClient) Delete() *AddOnInstallationDeleteRequest {
	return &AddOnInstallationDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the add-on installation.
func (c *AddOnInstallationClient) Get() *AddOnInstallationGetRequest {
	return &AddOnInstallationGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the add-on installation.
func (c *AddOnInstallationClient) Update() *AddOnInstallationUpdateRequest {
	return &AddOnInstallationUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AddOnInstallationPollRequest is the request for the Poll method.
type AddOnInstallationPollRequest struct {
	request    *AddOnInstallationGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AddOnInstallationPollRequest) Parameter(name string, value interface{}) *AddOnInstallationPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AddOnInstallationPollRequest) Header(name string, value interface{}) *AddOnInstallationPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AddOnInstallationPollRequest) Interval(value time.Duration) *AddOnInstallationPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AddOnInstallationPollRequest) Status(value int) *AddOnInstallationPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AddOnInstallationPollRequest) Predicate(value func(*AddOnInstallationGetResponse) bool) *AddOnInstallationPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AddOnInstallationGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AddOnInstallationPollRequest) StartContext(ctx context.Context) (response *AddOnInstallationPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AddOnInstallationPollResponse{
			response: result.(*AddOnInstallationGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AddOnInstallationPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AddOnInstallationPollResponse is the response for the Poll method.
type AddOnInstallationPollResponse struct {
	response *AddOnInstallationGetResponse
}

// Status returns the response status code.
func (r *AddOnInstallationPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AddOnInstallationPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AddOnInstallationPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AddOnInstallationPollResponse) Body() *AddOnInstallation {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddOnInstallationPollResponse) GetBody() (value *AddOnInstallation, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AddOnInstallationClient) Poll() *AddOnInstallationPollRequest {
	return &AddOnInstallationPollRequest{
		request: c.Get(),
	}
}

// AddOnInstallationDeleteRequest is the request for the 'delete' method.
type AddOnInstallationDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddOnInstallationDeleteRequest) Parameter(name string, value interface{}) *AddOnInstallationDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnInstallationDeleteRequest) Header(name string, value interface{}) *AddOnInstallationDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnInstallationDeleteRequest) Impersonate(user string) *AddOnInstallationDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnInstallationDeleteRequest) Send() (result *AddOnInstallationDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnInstallationDeleteRequest) SendContext(ctx context.Context) (result *AddOnInstallationDeleteResponse, err error) {
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
	result = &AddOnInstallationDeleteResponse{}
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

// AddOnInstallationDeleteResponse is the response for the 'delete' method.
type AddOnInstallationDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *AddOnInstallationDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnInstallationDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnInstallationDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// AddOnInstallationGetRequest is the request for the 'get' method.
type AddOnInstallationGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AddOnInstallationGetRequest) Parameter(name string, value interface{}) *AddOnInstallationGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnInstallationGetRequest) Header(name string, value interface{}) *AddOnInstallationGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnInstallationGetRequest) Impersonate(user string) *AddOnInstallationGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnInstallationGetRequest) Send() (result *AddOnInstallationGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnInstallationGetRequest) SendContext(ctx context.Context) (result *AddOnInstallationGetResponse, err error) {
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
	result = &AddOnInstallationGetResponse{}
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
	err = readAddOnInstallationGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddOnInstallationGetResponse is the response for the 'get' method.
type AddOnInstallationGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddOnInstallation
}

// Status returns the response status code.
func (r *AddOnInstallationGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnInstallationGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnInstallationGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddOnInstallationGetResponse) Body() *AddOnInstallation {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddOnInstallationGetResponse) GetBody() (value *AddOnInstallation, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddOnInstallationUpdateRequest is the request for the 'update' method.
type AddOnInstallationUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddOnInstallation
}

// Parameter adds a query parameter.
func (r *AddOnInstallationUpdateRequest) Parameter(name string, value interface{}) *AddOnInstallationUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnInstallationUpdateRequest) Header(name string, value interface{}) *AddOnInstallationUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnInstallationUpdateRequest) Impersonate(user string) *AddOnInstallationUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *AddOnInstallationUpdateRequest) Body(value *AddOnInstallation) *AddOnInstallationUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnInstallationUpdateRequest) Send() (result *AddOnInstallationUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnInstallationUpdateRequest) SendContext(ctx context.Context) (result *AddOnInstallationUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddOnInstallationUpdateRequest(r, buffer)
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
	result = &AddOnInstallationUpdateResponse{}
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
	err = readAddOnInstallationUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddOnInstallationUpdateResponse is the response for the 'update' method.
type AddOnInstallationUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddOnInstallation
}

// Status returns the response status code.
func (r *AddOnInstallationUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnInstallationUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnInstallationUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AddOnInstallationUpdateResponse) Body() *AddOnInstallation {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AddOnInstallationUpdateResponse) GetBody() (value *AddOnInstallation, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
