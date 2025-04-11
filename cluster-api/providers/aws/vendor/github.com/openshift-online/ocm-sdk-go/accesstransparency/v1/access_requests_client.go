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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// AccessRequestsClient is the client of the 'access_requests' resource.
//
// Manages the collection of access requests.
type AccessRequestsClient struct {
	transport http.RoundTripper
	path      string
}

// NewAccessRequestsClient creates a new client for the 'access_requests'
// resource using the given transport to send the requests and receive the
// responses.
func NewAccessRequestsClient(transport http.RoundTripper, path string) *AccessRequestsClient {
	return &AccessRequestsClient{
		transport: transport,
		path:      path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of access requests.
func (c *AccessRequestsClient) List() *AccessRequestsListRequest {
	return &AccessRequestsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Post creates a request for the 'post' method.
//
// Create a new access request and add it to the collection of access requests.
func (c *AccessRequestsClient) Post() *AccessRequestsPostRequest {
	return &AccessRequestsPostRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AccessRequest returns the target 'access_request' resource for the given identifier.
//
// Returns a reference to the service that manages a specific access request.
func (c *AccessRequestsClient) AccessRequest(id string) *AccessRequestClient {
	return NewAccessRequestClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// AccessRequestsListRequest is the request for the 'list' method.
type AccessRequestsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	order     *string
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *AccessRequestsListRequest) Parameter(name string, value interface{}) *AccessRequestsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccessRequestsListRequest) Header(name string, value interface{}) *AccessRequestsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccessRequestsListRequest) Impersonate(user string) *AccessRequestsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement, but using the names of the attributes of the access request instead of
// the names of the columns of a table. For example, in order to sort the access requests
// descending by created_at the value should be:
//
// ```sql
// created_at desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *AccessRequestsListRequest) Order(value string) *AccessRequestsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AccessRequestsListRequest) Page(value int) *AccessRequestsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of an
// SQL statement, but using the names of the attributes of the access request instead of
// the names of the columns of a table. For example, in order to retrieve all the
// access requests with a requested_by starting with `my` the value should be:
//
// ```sql
// requested_by like 'my%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the access requests
// that the user has permission to see will be returned.
func (r *AccessRequestsListRequest) Search(value string) *AccessRequestsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AccessRequestsListRequest) Size(value int) *AccessRequestsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccessRequestsListRequest) Send() (result *AccessRequestsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccessRequestsListRequest) SendContext(ctx context.Context) (result *AccessRequestsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.order != nil {
		helpers.AddValue(&query, "order", *r.order)
	}
	if r.page != nil {
		helpers.AddValue(&query, "page", *r.page)
	}
	if r.search != nil {
		helpers.AddValue(&query, "search", *r.search)
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
	result = &AccessRequestsListResponse{}
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
	err = readAccessRequestsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AccessRequestsListResponse is the response for the 'list' method.
type AccessRequestsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *AccessRequestList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *AccessRequestsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccessRequestsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccessRequestsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of access requests.
func (r *AccessRequestsListResponse) Items() *AccessRequestList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of access requests.
func (r *AccessRequestsListResponse) GetItems() (value *AccessRequestList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AccessRequestsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AccessRequestsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AccessRequestsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *AccessRequestsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection that match the search criteria,
// regardless of the size of the page.
func (r *AccessRequestsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection that match the search criteria,
// regardless of the size of the page.
func (r *AccessRequestsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}

// AccessRequestsPostRequest is the request for the 'post' method.
type AccessRequestsPostRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AccessRequestPostRequest
}

// Parameter adds a query parameter.
func (r *AccessRequestsPostRequest) Parameter(name string, value interface{}) *AccessRequestsPostRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccessRequestsPostRequest) Header(name string, value interface{}) *AccessRequestsPostRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccessRequestsPostRequest) Impersonate(user string) *AccessRequestsPostRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Input to include new access request.
func (r *AccessRequestsPostRequest) Body(value *AccessRequestPostRequest) *AccessRequestsPostRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccessRequestsPostRequest) Send() (result *AccessRequestsPostResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccessRequestsPostRequest) SendContext(ctx context.Context) (result *AccessRequestsPostResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAccessRequestsPostRequest(r, buffer)
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
	result = &AccessRequestsPostResponse{}
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
	err = readAccessRequestsPostResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AccessRequestsPostResponse is the response for the 'post' method.
type AccessRequestsPostResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AccessRequest
}

// Status returns the response status code.
func (r *AccessRequestsPostResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccessRequestsPostResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccessRequestsPostResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Newly access request.
func (r *AccessRequestsPostResponse) Body() *AccessRequest {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Newly access request.
func (r *AccessRequestsPostResponse) GetBody() (value *AccessRequest, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
