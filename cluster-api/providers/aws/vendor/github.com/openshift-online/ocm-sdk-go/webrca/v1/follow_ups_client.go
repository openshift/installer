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

// FollowUpsClient is the client of the 'follow_ups' resource.
//
// Manages the collection of follow-ups.
type FollowUpsClient struct {
	transport http.RoundTripper
	path      string
}

// NewFollowUpsClient creates a new client for the 'follow_ups'
// resource using the given transport to send the requests and receive the
// responses.
func NewFollowUpsClient(transport http.RoundTripper, path string) *FollowUpsClient {
	return &FollowUpsClient{
		transport: transport,
		path:      path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of follow-ups
func (c *FollowUpsClient) List() *FollowUpsListRequest {
	return &FollowUpsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// FollowUp returns the target 'follow_up' resource for the given identifier.
func (c *FollowUpsClient) FollowUp(id string) *FollowUpClient {
	return NewFollowUpClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// FollowUpsListRequest is the request for the 'list' method.
type FollowUpsListRequest struct {
	transport      http.RoundTripper
	path           string
	query          url.Values
	header         http.Header
	createdAfter   *time.Time
	createdBefore  *time.Time
	followUpStatus *string
	orderBy        *string
	page           *int
	size           *int
}

// Parameter adds a query parameter.
func (r *FollowUpsListRequest) Parameter(name string, value interface{}) *FollowUpsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *FollowUpsListRequest) Header(name string, value interface{}) *FollowUpsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *FollowUpsListRequest) Impersonate(user string) *FollowUpsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// CreatedAfter sets the value of the 'created_after' parameter.
func (r *FollowUpsListRequest) CreatedAfter(value time.Time) *FollowUpsListRequest {
	r.createdAfter = &value
	return r
}

// CreatedBefore sets the value of the 'created_before' parameter.
func (r *FollowUpsListRequest) CreatedBefore(value time.Time) *FollowUpsListRequest {
	r.createdBefore = &value
	return r
}

// FollowUpStatus sets the value of the 'follow_up_status' parameter.
func (r *FollowUpsListRequest) FollowUpStatus(value string) *FollowUpsListRequest {
	r.followUpStatus = &value
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *FollowUpsListRequest) OrderBy(value string) *FollowUpsListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *FollowUpsListRequest) Page(value int) *FollowUpsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *FollowUpsListRequest) Size(value int) *FollowUpsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *FollowUpsListRequest) Send() (result *FollowUpsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *FollowUpsListRequest) SendContext(ctx context.Context) (result *FollowUpsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.createdAfter != nil {
		helpers.AddValue(&query, "created_after", *r.createdAfter)
	}
	if r.createdBefore != nil {
		helpers.AddValue(&query, "created_before", *r.createdBefore)
	}
	if r.followUpStatus != nil {
		helpers.AddValue(&query, "follow_up_status", *r.followUpStatus)
	}
	if r.orderBy != nil {
		helpers.AddValue(&query, "order_by", *r.orderBy)
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
	result = &FollowUpsListResponse{}
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
	err = readFollowUpsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// FollowUpsListResponse is the response for the 'list' method.
type FollowUpsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *FollowUpList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *FollowUpsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *FollowUpsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *FollowUpsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *FollowUpsListResponse) Items() *FollowUpList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *FollowUpsListResponse) GetItems() (value *FollowUpList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *FollowUpsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *FollowUpsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *FollowUpsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *FollowUpsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *FollowUpsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *FollowUpsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
