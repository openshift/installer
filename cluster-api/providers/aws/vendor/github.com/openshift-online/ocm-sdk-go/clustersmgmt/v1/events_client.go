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

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// EventsClient is the client of the 'events' resource.
//
// Manages a collection used to track events reported by external clients.
type EventsClient struct {
	transport http.RoundTripper
	path      string
}

// NewEventsClient creates a new client for the 'events'
// resource using the given transport to send the requests and receive the
// responses.
func NewEventsClient(transport http.RoundTripper, path string) *EventsClient {
	return &EventsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new event to be tracked. When sending a new event request,
// it gets tracked in Prometheus, Pendo, CloudWatch, or whichever
// analytics client is configured as part of clusters service. This
// allows for reporting on events that happen outside of a regular API
// request, but are found to be useful for understanding customer
// needs and possible blockers.
func (c *EventsClient) Add() *EventsAddRequest {
	return &EventsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// EventsAddRequest is the request for the 'add' method.
type EventsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Event
}

// Parameter adds a query parameter.
func (r *EventsAddRequest) Parameter(name string, value interface{}) *EventsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *EventsAddRequest) Header(name string, value interface{}) *EventsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *EventsAddRequest) Impersonate(user string) *EventsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the event.
func (r *EventsAddRequest) Body(value *Event) *EventsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *EventsAddRequest) Send() (result *EventsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *EventsAddRequest) SendContext(ctx context.Context) (result *EventsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeEventsAddRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "POST",
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
	result = &EventsAddResponse{}
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
	err = readEventsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// EventsAddResponse is the response for the 'add' method.
type EventsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Event
}

// Status returns the response status code.
func (r *EventsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *EventsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *EventsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the event.
func (r *EventsAddResponse) Body() *Event {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the event.
func (r *EventsAddResponse) GetBody() (value *Event, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
