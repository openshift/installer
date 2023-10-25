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

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// NotificationsClient is the client of the 'notifications' resource.
//
// Manages the collection of notifications.
type NotificationsClient struct {
	transport http.RoundTripper
	path      string
}

// NewNotificationsClient creates a new client for the 'notifications'
// resource using the given transport to send the requests and receive the
// responses.
func NewNotificationsClient(transport http.RoundTripper, path string) *NotificationsClient {
	return &NotificationsClient{
		transport: transport,
		path:      path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of notifications
func (c *NotificationsClient) List() *NotificationsListRequest {
	return &NotificationsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Notification returns the target 'notification' resource for the given identifier.
func (c *NotificationsClient) Notification(id string) *NotificationClient {
	return NewNotificationClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// NotificationsListRequest is the request for the 'list' method.
type NotificationsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	checked   *bool
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *NotificationsListRequest) Parameter(name string, value interface{}) *NotificationsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *NotificationsListRequest) Header(name string, value interface{}) *NotificationsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *NotificationsListRequest) Impersonate(user string) *NotificationsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Checked sets the value of the 'checked' parameter.
func (r *NotificationsListRequest) Checked(value bool) *NotificationsListRequest {
	r.checked = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *NotificationsListRequest) Page(value int) *NotificationsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *NotificationsListRequest) Size(value int) *NotificationsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *NotificationsListRequest) Send() (result *NotificationsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *NotificationsListRequest) SendContext(ctx context.Context) (result *NotificationsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.checked != nil {
		helpers.AddValue(&query, "checked", *r.checked)
	}
	if r.page != nil {
		helpers.AddValue(&query, "page", *r.page)
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
	result = &NotificationsListResponse{}
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
	err = readNotificationsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// NotificationsListResponse is the response for the 'list' method.
type NotificationsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *NotificationList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *NotificationsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *NotificationsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *NotificationsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *NotificationsListResponse) Items() *NotificationList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *NotificationsListResponse) GetItems() (value *NotificationList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *NotificationsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *NotificationsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *NotificationsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *NotificationsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *NotificationsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *NotificationsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
