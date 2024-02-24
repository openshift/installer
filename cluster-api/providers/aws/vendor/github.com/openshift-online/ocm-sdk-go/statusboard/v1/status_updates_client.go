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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	time "time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// StatusUpdatesClient is the client of the 'status_updates' resource.
//
// Manages the collection of statuses
type StatusUpdatesClient struct {
	transport http.RoundTripper
	path      string
}

// NewStatusUpdatesClient creates a new client for the 'status_updates'
// resource using the given transport to send the requests and receive the
// responses.
func NewStatusUpdatesClient(transport http.RoundTripper, path string) *StatusUpdatesClient {
	return &StatusUpdatesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
func (c *StatusUpdatesClient) Add() *StatusUpdatesAddRequest {
	return &StatusUpdatesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of statuses.
func (c *StatusUpdatesClient) List() *StatusUpdatesListRequest {
	return &StatusUpdatesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Status returns the target 'status' resource for the given identifier.
func (c *StatusUpdatesClient) Status(id string) *StatusClient {
	return NewStatusClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// StatusUpdatesAddRequest is the request for the 'add' method.
type StatusUpdatesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Status
}

// Parameter adds a query parameter.
func (r *StatusUpdatesAddRequest) Parameter(name string, value interface{}) *StatusUpdatesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *StatusUpdatesAddRequest) Header(name string, value interface{}) *StatusUpdatesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *StatusUpdatesAddRequest) Impersonate(user string) *StatusUpdatesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *StatusUpdatesAddRequest) Body(value *Status) *StatusUpdatesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *StatusUpdatesAddRequest) Send() (result *StatusUpdatesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *StatusUpdatesAddRequest) SendContext(ctx context.Context) (result *StatusUpdatesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeStatusUpdatesAddRequest(r, buffer)
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
	result = &StatusUpdatesAddResponse{}
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
	err = readStatusUpdatesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// StatusUpdatesAddResponse is the response for the 'add' method.
type StatusUpdatesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Status
}

// Status returns the response status code.
func (r *StatusUpdatesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *StatusUpdatesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *StatusUpdatesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *StatusUpdatesAddResponse) Body() *Status {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *StatusUpdatesAddResponse) GetBody() (value *Status, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// StatusUpdatesListRequest is the request for the 'list' method.
type StatusUpdatesListRequest struct {
	transport     http.RoundTripper
	path          string
	query         url.Values
	header        http.Header
	createdAfter  *time.Time
	createdBefore *time.Time
	fullNames     *string
	limitScope    *time.Time
	page          *int
	productIds    *string
	search        *string
	size          *int
}

// Parameter adds a query parameter.
func (r *StatusUpdatesListRequest) Parameter(name string, value interface{}) *StatusUpdatesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *StatusUpdatesListRequest) Header(name string, value interface{}) *StatusUpdatesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *StatusUpdatesListRequest) Impersonate(user string) *StatusUpdatesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// CreatedAfter sets the value of the 'created_after' parameter.
func (r *StatusUpdatesListRequest) CreatedAfter(value time.Time) *StatusUpdatesListRequest {
	r.createdAfter = &value
	return r
}

// CreatedBefore sets the value of the 'created_before' parameter.
func (r *StatusUpdatesListRequest) CreatedBefore(value time.Time) *StatusUpdatesListRequest {
	r.createdBefore = &value
	return r
}

// FullNames sets the value of the 'full_names' parameter.
func (r *StatusUpdatesListRequest) FullNames(value string) *StatusUpdatesListRequest {
	r.fullNames = &value
	return r
}

// LimitScope sets the value of the 'limit_scope' parameter.
func (r *StatusUpdatesListRequest) LimitScope(value time.Time) *StatusUpdatesListRequest {
	r.limitScope = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *StatusUpdatesListRequest) Page(value int) *StatusUpdatesListRequest {
	r.page = &value
	return r
}

// ProductIds sets the value of the 'product_ids' parameter.
func (r *StatusUpdatesListRequest) ProductIds(value string) *StatusUpdatesListRequest {
	r.productIds = &value
	return r
}

// Search sets the value of the 'search' parameter.
func (r *StatusUpdatesListRequest) Search(value string) *StatusUpdatesListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *StatusUpdatesListRequest) Size(value int) *StatusUpdatesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *StatusUpdatesListRequest) Send() (result *StatusUpdatesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *StatusUpdatesListRequest) SendContext(ctx context.Context) (result *StatusUpdatesListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.createdAfter != nil {
		helpers.AddValue(&query, "created_after", *r.createdAfter)
	}
	if r.createdBefore != nil {
		helpers.AddValue(&query, "created_before", *r.createdBefore)
	}
	if r.fullNames != nil {
		helpers.AddValue(&query, "full_names", *r.fullNames)
	}
	if r.limitScope != nil {
		helpers.AddValue(&query, "limit_scope", *r.limitScope)
	}
	if r.page != nil {
		helpers.AddValue(&query, "page", *r.page)
	}
	if r.productIds != nil {
		helpers.AddValue(&query, "product_ids", *r.productIds)
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
	result = &StatusUpdatesListResponse{}
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
	err = readStatusUpdatesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// StatusUpdatesListResponse is the response for the 'list' method.
type StatusUpdatesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *StatusList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *StatusUpdatesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *StatusUpdatesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *StatusUpdatesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *StatusUpdatesListResponse) Items() *StatusList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *StatusUpdatesListResponse) GetItems() (value *StatusList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *StatusUpdatesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *StatusUpdatesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *StatusUpdatesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *StatusUpdatesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *StatusUpdatesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *StatusUpdatesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
