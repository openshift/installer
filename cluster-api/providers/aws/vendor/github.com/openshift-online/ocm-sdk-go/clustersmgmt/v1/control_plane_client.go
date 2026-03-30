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

// ControlPlaneClient is the client of the 'control_plane' resource.
//
// Manages control plane resources.
type ControlPlaneClient struct {
	transport http.RoundTripper
	path      string
}

// NewControlPlaneClient creates a new client for the 'control_plane'
// resource using the given transport to send the requests and receive the
// responses.
func NewControlPlaneClient(transport http.RoundTripper, path string) *ControlPlaneClient {
	return &ControlPlaneClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves the details of the control plane
func (c *ControlPlaneClient) Get() *ControlPlaneGetRequest {
	return &ControlPlaneGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
//
// Updates the control plane
func (c *ControlPlaneClient) Update() *ControlPlaneUpdateRequest {
	return &ControlPlaneUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// LogForwarders returns the target 'log_forwarders' resource.
//
// Reference to the collection of log forwarders for the control plane.
func (c *ControlPlaneClient) LogForwarders() *LogForwardersClient {
	return NewLogForwardersClient(
		c.transport,
		path.Join(c.path, "log_forwarders"),
	)
}

// UpgradePolicies returns the target 'control_plane_upgrade_policies' resource.
//
// Reference to the collection of upgrade policies for the control plane.
func (c *ControlPlaneClient) UpgradePolicies() *ControlPlaneUpgradePoliciesClient {
	return NewControlPlaneUpgradePoliciesClient(
		c.transport,
		path.Join(c.path, "upgrade_policies"),
	)
}

// ControlPlanePollRequest is the request for the Poll method.
type ControlPlanePollRequest struct {
	request    *ControlPlaneGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ControlPlanePollRequest) Parameter(name string, value interface{}) *ControlPlanePollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ControlPlanePollRequest) Header(name string, value interface{}) *ControlPlanePollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ControlPlanePollRequest) Interval(value time.Duration) *ControlPlanePollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ControlPlanePollRequest) Status(value int) *ControlPlanePollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ControlPlanePollRequest) Predicate(value func(*ControlPlaneGetResponse) bool) *ControlPlanePollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ControlPlaneGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ControlPlanePollRequest) StartContext(ctx context.Context) (response *ControlPlanePollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ControlPlanePollResponse{
			response: result.(*ControlPlaneGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ControlPlanePollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ControlPlanePollResponse is the response for the Poll method.
type ControlPlanePollResponse struct {
	response *ControlPlaneGetResponse
}

// Status returns the response status code.
func (r *ControlPlanePollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ControlPlanePollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ControlPlanePollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *ControlPlanePollResponse) Body() *ControlPlane {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ControlPlanePollResponse) GetBody() (value *ControlPlane, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ControlPlaneClient) Poll() *ControlPlanePollRequest {
	return &ControlPlanePollRequest{
		request: c.Get(),
	}
}

// ControlPlaneGetRequest is the request for the 'get' method.
type ControlPlaneGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *ControlPlaneGetRequest) Parameter(name string, value interface{}) *ControlPlaneGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ControlPlaneGetRequest) Header(name string, value interface{}) *ControlPlaneGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ControlPlaneGetRequest) Impersonate(user string) *ControlPlaneGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ControlPlaneGetRequest) Send() (result *ControlPlaneGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ControlPlaneGetRequest) SendContext(ctx context.Context) (result *ControlPlaneGetResponse, err error) {
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
	result = &ControlPlaneGetResponse{}
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
	err = readControlPlaneGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ControlPlaneGetResponse is the response for the 'get' method.
type ControlPlaneGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ControlPlane
}

// Status returns the response status code.
func (r *ControlPlaneGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ControlPlaneGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ControlPlaneGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ControlPlaneGetResponse) Body() *ControlPlane {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ControlPlaneGetResponse) GetBody() (value *ControlPlane, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ControlPlaneUpdateRequest is the request for the 'update' method.
type ControlPlaneUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ControlPlane
}

// Parameter adds a query parameter.
func (r *ControlPlaneUpdateRequest) Parameter(name string, value interface{}) *ControlPlaneUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ControlPlaneUpdateRequest) Header(name string, value interface{}) *ControlPlaneUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ControlPlaneUpdateRequest) Impersonate(user string) *ControlPlaneUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ControlPlaneUpdateRequest) Body(value *ControlPlane) *ControlPlaneUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ControlPlaneUpdateRequest) Send() (result *ControlPlaneUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ControlPlaneUpdateRequest) SendContext(ctx context.Context) (result *ControlPlaneUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeControlPlaneUpdateRequest(r, buffer)
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
	result = &ControlPlaneUpdateResponse{}
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
	err = readControlPlaneUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ControlPlaneUpdateResponse is the response for the 'update' method.
type ControlPlaneUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ControlPlane
}

// Status returns the response status code.
func (r *ControlPlaneUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ControlPlaneUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ControlPlaneUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ControlPlaneUpdateResponse) Body() *ControlPlane {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ControlPlaneUpdateResponse) GetBody() (value *ControlPlane, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
