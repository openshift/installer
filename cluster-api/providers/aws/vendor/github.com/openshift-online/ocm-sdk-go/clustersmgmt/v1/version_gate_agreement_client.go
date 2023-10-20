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

// VersionGateAgreementClient is the client of the 'version_gate_agreement' resource.
//
// Manages a specific version gate agreement.
type VersionGateAgreementClient struct {
	transport http.RoundTripper
	path      string
}

// NewVersionGateAgreementClient creates a new client for the 'version_gate_agreement'
// resource using the given transport to send the requests and receive the
// responses.
func NewVersionGateAgreementClient(transport http.RoundTripper, path string) *VersionGateAgreementClient {
	return &VersionGateAgreementClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
//
// Deletes the version gate agreement.
func (c *VersionGateAgreementClient) Delete() *VersionGateAgreementDeleteRequest {
	return &VersionGateAgreementDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the version gate agreement.
func (c *VersionGateAgreementClient) Get() *VersionGateAgreementGetRequest {
	return &VersionGateAgreementGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// VersionGateAgreementPollRequest is the request for the Poll method.
type VersionGateAgreementPollRequest struct {
	request    *VersionGateAgreementGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *VersionGateAgreementPollRequest) Parameter(name string, value interface{}) *VersionGateAgreementPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *VersionGateAgreementPollRequest) Header(name string, value interface{}) *VersionGateAgreementPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *VersionGateAgreementPollRequest) Interval(value time.Duration) *VersionGateAgreementPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *VersionGateAgreementPollRequest) Status(value int) *VersionGateAgreementPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *VersionGateAgreementPollRequest) Predicate(value func(*VersionGateAgreementGetResponse) bool) *VersionGateAgreementPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*VersionGateAgreementGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *VersionGateAgreementPollRequest) StartContext(ctx context.Context) (response *VersionGateAgreementPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &VersionGateAgreementPollResponse{
			response: result.(*VersionGateAgreementGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *VersionGateAgreementPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// VersionGateAgreementPollResponse is the response for the Poll method.
type VersionGateAgreementPollResponse struct {
	response *VersionGateAgreementGetResponse
}

// Status returns the response status code.
func (r *VersionGateAgreementPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *VersionGateAgreementPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *VersionGateAgreementPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *VersionGateAgreementPollResponse) Body() *VersionGateAgreement {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *VersionGateAgreementPollResponse) GetBody() (value *VersionGateAgreement, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *VersionGateAgreementClient) Poll() *VersionGateAgreementPollRequest {
	return &VersionGateAgreementPollRequest{
		request: c.Get(),
	}
}

// VersionGateAgreementDeleteRequest is the request for the 'delete' method.
type VersionGateAgreementDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *VersionGateAgreementDeleteRequest) Parameter(name string, value interface{}) *VersionGateAgreementDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *VersionGateAgreementDeleteRequest) Header(name string, value interface{}) *VersionGateAgreementDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *VersionGateAgreementDeleteRequest) Impersonate(user string) *VersionGateAgreementDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *VersionGateAgreementDeleteRequest) Send() (result *VersionGateAgreementDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *VersionGateAgreementDeleteRequest) SendContext(ctx context.Context) (result *VersionGateAgreementDeleteResponse, err error) {
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
	result = &VersionGateAgreementDeleteResponse{}
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

// VersionGateAgreementDeleteResponse is the response for the 'delete' method.
type VersionGateAgreementDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *VersionGateAgreementDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *VersionGateAgreementDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *VersionGateAgreementDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// VersionGateAgreementGetRequest is the request for the 'get' method.
type VersionGateAgreementGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *VersionGateAgreementGetRequest) Parameter(name string, value interface{}) *VersionGateAgreementGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *VersionGateAgreementGetRequest) Header(name string, value interface{}) *VersionGateAgreementGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *VersionGateAgreementGetRequest) Impersonate(user string) *VersionGateAgreementGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *VersionGateAgreementGetRequest) Send() (result *VersionGateAgreementGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *VersionGateAgreementGetRequest) SendContext(ctx context.Context) (result *VersionGateAgreementGetResponse, err error) {
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
	result = &VersionGateAgreementGetResponse{}
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
	err = readVersionGateAgreementGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// VersionGateAgreementGetResponse is the response for the 'get' method.
type VersionGateAgreementGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *VersionGateAgreement
}

// Status returns the response status code.
func (r *VersionGateAgreementGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *VersionGateAgreementGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *VersionGateAgreementGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *VersionGateAgreementGetResponse) Body() *VersionGateAgreement {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *VersionGateAgreementGetResponse) GetBody() (value *VersionGateAgreement, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
