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
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	time "time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// EventsClient is the client of the 'events' resource.
//
// Manages the collection of events.
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

// List creates a request for the 'list' method.
//
// Retrieves the list of events
func (c *EventsClient) List() *EventsListRequest {
	return &EventsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Event returns the target 'event' resource for the given identifier.
func (c *EventsClient) Event(id string) *EventClient {
	return NewEventClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// EventsListRequest is the request for the 'list' method.
type EventsListRequest struct {
	transport        http.RoundTripper
	path             string
	query            url.Values
	header           http.Header
	createdAfter     *time.Time
	createdBefore    *time.Time
	eventType        *string
	note             *string
	orderBy          *string
	page             *int
	showSystemEvents *bool
	size             *int
}

// Parameter adds a query parameter.
func (r *EventsListRequest) Parameter(name string, value interface{}) *EventsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *EventsListRequest) Header(name string, value interface{}) *EventsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *EventsListRequest) Impersonate(user string) *EventsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// CreatedAfter sets the value of the 'created_after' parameter.
func (r *EventsListRequest) CreatedAfter(value time.Time) *EventsListRequest {
	r.createdAfter = &value
	return r
}

// CreatedBefore sets the value of the 'created_before' parameter.
func (r *EventsListRequest) CreatedBefore(value time.Time) *EventsListRequest {
	r.createdBefore = &value
	return r
}

// EventType sets the value of the 'event_type' parameter.
func (r *EventsListRequest) EventType(value string) *EventsListRequest {
	r.eventType = &value
	return r
}

// Note sets the value of the 'note' parameter.
func (r *EventsListRequest) Note(value string) *EventsListRequest {
	r.note = &value
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *EventsListRequest) OrderBy(value string) *EventsListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *EventsListRequest) Page(value int) *EventsListRequest {
	r.page = &value
	return r
}

// ShowSystemEvents sets the value of the 'show_system_events' parameter.
func (r *EventsListRequest) ShowSystemEvents(value bool) *EventsListRequest {
	r.showSystemEvents = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *EventsListRequest) Size(value int) *EventsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *EventsListRequest) Send() (result *EventsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *EventsListRequest) SendContext(ctx context.Context) (result *EventsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.createdAfter != nil {
		helpers.AddValue(&query, "created_after", *r.createdAfter)
	}
	if r.createdBefore != nil {
		helpers.AddValue(&query, "created_before", *r.createdBefore)
	}
	if r.eventType != nil {
		helpers.AddValue(&query, "event_type", *r.eventType)
	}
	if r.note != nil {
		helpers.AddValue(&query, "note", *r.note)
	}
	if r.orderBy != nil {
		helpers.AddValue(&query, "order_by", *r.orderBy)
	}
	if r.page != nil {
		helpers.AddValue(&query, "page", *r.page)
	}
	if r.showSystemEvents != nil {
		helpers.AddValue(&query, "show_system_events", *r.showSystemEvents)
	}
	if r.size != nil {
		helpers.AddValue(&query, "size", *r.size)
	}
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
	result = &EventsListResponse{}
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
	err = readEventsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// EventsListResponse is the response for the 'list' method.
type EventsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *EventList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *EventsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *EventsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *EventsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *EventsListResponse) Items() *EventList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *EventsListResponse) GetItems() (value *EventList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *EventsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *EventsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *EventsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *EventsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *EventsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *EventsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
