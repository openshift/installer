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

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// AddOnVersionsClient is the client of the 'add_on_versions' resource.
//
// Manages the collection of add-on versions.
type AddOnVersionsClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddOnVersionsClient creates a new client for the 'add_on_versions'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddOnVersionsClient(transport http.RoundTripper, path string) *AddOnVersionsClient {
	return &AddOnVersionsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Create a new add-on version and add it to the collection of add-ons.
func (c *AddOnVersionsClient) Add() *AddOnVersionsAddRequest {
	return &AddOnVersionsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of add-on versions.
func (c *AddOnVersionsClient) List() *AddOnVersionsListRequest {
	return &AddOnVersionsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Version returns the target 'add_on_version' resource for the given identifier.
//
// Returns a reference to the service that manages a specific add-on version.
func (c *AddOnVersionsClient) Version(id string) *AddOnVersionClient {
	return NewAddOnVersionClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// AddOnVersionsAddRequest is the request for the 'add' method.
type AddOnVersionsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddOnVersion
}

// Parameter adds a query parameter.
func (r *AddOnVersionsAddRequest) Parameter(name string, value interface{}) *AddOnVersionsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnVersionsAddRequest) Header(name string, value interface{}) *AddOnVersionsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnVersionsAddRequest) Impersonate(user string) *AddOnVersionsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the add-on version.
func (r *AddOnVersionsAddRequest) Body(value *AddOnVersion) *AddOnVersionsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnVersionsAddRequest) Send() (result *AddOnVersionsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnVersionsAddRequest) SendContext(ctx context.Context) (result *AddOnVersionsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddOnVersionsAddRequest(r, buffer)
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
	result = &AddOnVersionsAddResponse{}
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
	err = readAddOnVersionsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddOnVersionsAddResponse is the response for the 'add' method.
type AddOnVersionsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddOnVersion
}

// Status returns the response status code.
func (r *AddOnVersionsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnVersionsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnVersionsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the add-on version.
func (r *AddOnVersionsAddResponse) Body() *AddOnVersion {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the add-on version.
func (r *AddOnVersionsAddResponse) GetBody() (value *AddOnVersion, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddOnVersionsListRequest is the request for the 'list' method.
type AddOnVersionsListRequest struct {
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
func (r *AddOnVersionsListRequest) Parameter(name string, value interface{}) *AddOnVersionsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnVersionsListRequest) Header(name string, value interface{}) *AddOnVersionsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnVersionsListRequest) Impersonate(user string) *AddOnVersionsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement, but using the names of the attributes of the add-on instead of
// the names of the columns of a table. For example, in order to sort the add-on
// versions descending by id the value should be:
//
// ```sql
// id desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *AddOnVersionsListRequest) Order(value string) *AddOnVersionsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddOnVersionsListRequest) Page(value int) *AddOnVersionsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of an
// SQL statement, but using the names of the attributes of the add-on version instead
// of the names of the columns of a table. For example, in order to retrieve all the
// add-on versions with an id starting with `0.1` the value should be:
//
// ```sql
// id like '0.1.%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the add-on
// versions that the user has permission to see will be returned.
func (r *AddOnVersionsListRequest) Search(value string) *AddOnVersionsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddOnVersionsListRequest) Size(value int) *AddOnVersionsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnVersionsListRequest) Send() (result *AddOnVersionsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnVersionsListRequest) SendContext(ctx context.Context) (result *AddOnVersionsListResponse, err error) {
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
	result = &AddOnVersionsListResponse{}
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
	err = readAddOnVersionsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddOnVersionsListResponse is the response for the 'list' method.
type AddOnVersionsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *AddOnVersionList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *AddOnVersionsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnVersionsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnVersionsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of add-on versions.
func (r *AddOnVersionsListResponse) Items() *AddOnVersionList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of add-on versions.
func (r *AddOnVersionsListResponse) GetItems() (value *AddOnVersionList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddOnVersionsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddOnVersionsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddOnVersionsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddOnVersionsListResponse) GetSize() (value int, ok bool) {
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
func (r *AddOnVersionsListResponse) Total() int {
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
func (r *AddOnVersionsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
