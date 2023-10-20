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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

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

// IncidentClient is the client of the 'incident' resource.
//
// Provides detailed information about a specific incident.
type IncidentClient struct {
	transport http.RoundTripper
	path      string
}

// NewIncidentClient creates a new client for the 'incident'
// resource using the given transport to send the requests and receive the
// responses.
func NewIncidentClient(transport http.RoundTripper, path string) *IncidentClient {
	return &IncidentClient{
		transport: transport,
		path:      path,
	}
}

// Delete creates a request for the 'delete' method.
func (c *IncidentClient) Delete() *IncidentDeleteRequest {
	return &IncidentDeleteRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Get creates a request for the 'get' method.
func (c *IncidentClient) Get() *IncidentGetRequest {
	return &IncidentGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Update creates a request for the 'update' method.
func (c *IncidentClient) Update() *IncidentUpdateRequest {
	return &IncidentUpdateRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Events returns the target 'events' resource.
func (c *IncidentClient) Events() *EventsClient {
	return NewEventsClient(
		c.transport,
		path.Join(c.path, "events"),
	)
}

// FollowUps returns the target 'follow_ups' resource.
func (c *IncidentClient) FollowUps() *FollowUpsClient {
	return NewFollowUpsClient(
		c.transport,
		path.Join(c.path, "follow_ups"),
	)
}

// Notifications returns the target 'notifications' resource.
func (c *IncidentClient) Notifications() *NotificationsClient {
	return NewNotificationsClient(
		c.transport,
		path.Join(c.path, "notifications"),
	)
}

// IncidentPollRequest is the request for the Poll method.
type IncidentPollRequest struct {
	request    *IncidentGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *IncidentPollRequest) Parameter(name string, value interface{}) *IncidentPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *IncidentPollRequest) Header(name string, value interface{}) *IncidentPollRequest {
	r.request.Header(name, value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *IncidentPollRequest) Interval(value time.Duration) *IncidentPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *IncidentPollRequest) Status(value int) *IncidentPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *IncidentPollRequest) Predicate(value func(*IncidentGetResponse) bool) *IncidentPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*IncidentGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *IncidentPollRequest) StartContext(ctx context.Context) (response *IncidentPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &IncidentPollResponse{
			response: result.(*IncidentGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *IncidentPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// IncidentPollResponse is the response for the Poll method.
type IncidentPollResponse struct {
	response *IncidentGetResponse
}

// Status returns the response status code.
func (r *IncidentPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *IncidentPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *IncidentPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
func (r *IncidentPollResponse) Body() *Incident {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *IncidentPollResponse) GetBody() (value *Incident, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *IncidentClient) Poll() *IncidentPollRequest {
	return &IncidentPollRequest{
		request: c.Get(),
	}
}

// IncidentDeleteRequest is the request for the 'delete' method.
type IncidentDeleteRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *IncidentDeleteRequest) Parameter(name string, value interface{}) *IncidentDeleteRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *IncidentDeleteRequest) Header(name string, value interface{}) *IncidentDeleteRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *IncidentDeleteRequest) Impersonate(user string) *IncidentDeleteRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *IncidentDeleteRequest) Send() (result *IncidentDeleteResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *IncidentDeleteRequest) SendContext(ctx context.Context) (result *IncidentDeleteResponse, err error) {
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
	result = &IncidentDeleteResponse{}
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

// IncidentDeleteResponse is the response for the 'delete' method.
type IncidentDeleteResponse struct {
	status int
	header http.Header
	err    *errors.Error
}

// Status returns the response status code.
func (r *IncidentDeleteResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *IncidentDeleteResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *IncidentDeleteResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// IncidentGetRequest is the request for the 'get' method.
type IncidentGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// Parameter adds a query parameter.
func (r *IncidentGetRequest) Parameter(name string, value interface{}) *IncidentGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *IncidentGetRequest) Header(name string, value interface{}) *IncidentGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *IncidentGetRequest) Impersonate(user string) *IncidentGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *IncidentGetRequest) Send() (result *IncidentGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *IncidentGetRequest) SendContext(ctx context.Context) (result *IncidentGetResponse, err error) {
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
	result = &IncidentGetResponse{}
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
	err = readIncidentGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// IncidentGetResponse is the response for the 'get' method.
type IncidentGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Incident
}

// Status returns the response status code.
func (r *IncidentGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *IncidentGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *IncidentGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *IncidentGetResponse) Body() *Incident {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *IncidentGetResponse) GetBody() (value *Incident, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// IncidentUpdateRequest is the request for the 'update' method.
type IncidentUpdateRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Incident
}

// Parameter adds a query parameter.
func (r *IncidentUpdateRequest) Parameter(name string, value interface{}) *IncidentUpdateRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *IncidentUpdateRequest) Header(name string, value interface{}) *IncidentUpdateRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *IncidentUpdateRequest) Impersonate(user string) *IncidentUpdateRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *IncidentUpdateRequest) Body(value *Incident) *IncidentUpdateRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *IncidentUpdateRequest) Send() (result *IncidentUpdateResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *IncidentUpdateRequest) SendContext(ctx context.Context) (result *IncidentUpdateResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeIncidentUpdateRequest(r, buffer)
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
	result = &IncidentUpdateResponse{}
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
	err = readIncidentUpdateResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// IncidentUpdateResponse is the response for the 'update' method.
type IncidentUpdateResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Incident
}

// Status returns the response status code.
func (r *IncidentUpdateResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *IncidentUpdateResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *IncidentUpdateResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *IncidentUpdateResponse) Body() *Incident {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *IncidentUpdateResponse) GetBody() (value *Incident, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
