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

// AccountClient is the client of the 'account' resource.
//
// Manages a specific account.
type AccountClient struct {
	transport http.RoundTripper
	path      string
}

// NewAccountClient creates a new client for the 'account'
// resource using the given transport to send the requests and receive the
// responses.
func NewAccountClient(transport http.RoundTripper, path string) *AccountClient {
	return &AccountClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
func (c *AccountClient) Delete() *AccountDeleteRequest {
	return &AccountDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the account.
func (c *AccountClient) Get() *AccountGetRequest {
	return &AccountGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the account.
func (c *AccountClient) Update() *AccountUpdateRequest {
	return &AccountUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Labels returns the target 'generic_labels' resource.
//
// Reference to the list of labels of a specific account.
func (c *AccountClient) Labels() *GenericLabelsClient {
	return NewGenericLabelsClient(
		c.transport,
		path.Join(c.path, "labels"),
	)
}

// AccountPollRequest is the request for the Poll method.
type AccountPollRequest struct {
	request    *AccountGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *AccountPollRequest) Parameter(name string, value interface{}) *AccountPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *AccountPollRequest) Header(name string, value interface{}) *AccountPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *AccountPollRequest) Interval(value time.Duration) *AccountPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *AccountPollRequest) Status(value int) *AccountPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *AccountPollRequest) Predicate(value func(*AccountGetResponse) bool) *AccountPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*AccountGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *AccountPollRequest) StartContext(ctx context.Context) (response *AccountPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &AccountPollResponse{
			response: result.(*AccountGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *AccountPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// AccountPollResponse is the response for the Poll method.
type AccountPollResponse struct {
	response *AccountGetResponse
}

// Status returns the response status code.
func (r *AccountPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *AccountPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *AccountPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *AccountPollResponse) Body() *Account {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AccountPollResponse) GetBody() (value *Account, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *AccountClient) Poll() *AccountPollRequest {
	return &AccountPollRequest{
		request: c.Get(),
	}
}

// AccountDeleteRequest is the request for the 'delete' method.
type AccountDeleteRequest struct {
	transport                 http.RoundTripper
	path                      string
	query                     url.Values
	header                    http.Header
	deleteAssociatedResources *bool
}

// Parameter adds a query parameter.
func (r *AccountDeleteRequest) Parameter(name string, value interface{}) *AccountDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccountDeleteRequest) Header(name string, value interface{}) *AccountDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccountDeleteRequest) Impersonate(user string) *AccountDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// DeleteAssociatedResources sets the value of the 'delete_associated_resources' parameter.
func (r *AccountDeleteRequest) DeleteAssociatedResources(value bool) *AccountDeleteRequest {
	r.deleteAssociatedResources = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccountDeleteRequest) Send() (result *AccountDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccountDeleteRequest) SendContext(ctx context.Context) (result *AccountDeleteResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.deleteAssociatedResources != nil {
		helpers.AddValue(&query, "deleteAssociatedResources", *r.deleteAssociatedResources)
	}
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
	result = &AccountDeleteResponse{}
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

// AccountDeleteResponse is the response for the 'delete' method.
type AccountDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *AccountDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccountDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccountDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// AccountGetRequest is the request for the 'get' method.
type AccountGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *AccountGetRequest) Parameter(name string, value interface{}) *AccountGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccountGetRequest) Header(name string, value interface{}) *AccountGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccountGetRequest) Impersonate(user string) *AccountGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccountGetRequest) Send() (result *AccountGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccountGetRequest) SendContext(ctx context.Context) (result *AccountGetResponse, err error) {
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
	result = &AccountGetResponse{}
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
	err = readAccountGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AccountGetResponse is the response for the 'get' method.
type AccountGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Account
}

// Status returns the response status code.
func (r *AccountGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccountGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccountGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AccountGetResponse) Body() *Account {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AccountGetResponse) GetBody() (value *Account, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AccountUpdateRequest is the request for the 'update' method.
type AccountUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Account
}

// Parameter adds a query parameter.
func (r *AccountUpdateRequest) Parameter(name string, value interface{}) *AccountUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccountUpdateRequest) Header(name string, value interface{}) *AccountUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccountUpdateRequest) Impersonate(user string) *AccountUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *AccountUpdateRequest) Body(value *Account) *AccountUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccountUpdateRequest) Send() (result *AccountUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccountUpdateRequest) SendContext(ctx context.Context) (result *AccountUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAccountUpdateRequest(r, buffer)
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
	result = &AccountUpdateResponse{}
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
	err = readAccountUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AccountUpdateResponse is the response for the 'update' method.
type AccountUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Account
}

// Status returns the response status code.
func (r *AccountUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccountUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccountUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *AccountUpdateResponse) Body() *Account {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *AccountUpdateResponse) GetBody() (value *Account, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
