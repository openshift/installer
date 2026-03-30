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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

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

// AccountGroupsClient is the client of the 'account_groups' resource.
//
// Manages the collection of account groups.
type AccountGroupsClient struct {
	transport http.RoundTripper
	path      string
}

// NewAccountGroupsClient creates a new client for the 'account_groups'
// resource using the given transport to send the requests and receive the
// responses.
func NewAccountGroupsClient(transport http.RoundTripper, path string) *AccountGroupsClient {
	return &AccountGroupsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Creates a new account group.
func (c *AccountGroupsClient) Add() *AccountGroupsAddRequest {
	return &AccountGroupsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of account groups.
func (c *AccountGroupsClient) List() *AccountGroupsListRequest {
	return &AccountGroupsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AccountGroup returns the target 'account_group' resource for the given identifier.
//
// Reference to the service that manages a specific account group.
func (c *AccountGroupsClient) AccountGroup(id string) *AccountGroupClient {
	return NewAccountGroupClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// AccountGroupsAddRequest is the request for the 'add' method.
type AccountGroupsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AccountGroup
}

// Parameter adds a query parameter.
func (r *AccountGroupsAddRequest) Parameter(name string, value interface{}) *AccountGroupsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccountGroupsAddRequest) Header(name string, value interface{}) *AccountGroupsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccountGroupsAddRequest) Impersonate(user string) *AccountGroupsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Account group data.
func (r *AccountGroupsAddRequest) Body(value *AccountGroup) *AccountGroupsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccountGroupsAddRequest) Send() (result *AccountGroupsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccountGroupsAddRequest) SendContext(ctx context.Context) (result *AccountGroupsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAccountGroupsAddRequest(r, buffer)
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
	result = &AccountGroupsAddResponse{}
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
	err = readAccountGroupsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AccountGroupsAddResponse is the response for the 'add' method.
type AccountGroupsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AccountGroup
}

// Status returns the response status code.
func (r *AccountGroupsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccountGroupsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccountGroupsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Account group data.
func (r *AccountGroupsAddResponse) Body() *AccountGroup {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Account group data.
func (r *AccountGroupsAddResponse) GetBody() (value *AccountGroup, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AccountGroupsListRequest is the request for the 'list' method.
type AccountGroupsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	fields    *string
	order     *string
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *AccountGroupsListRequest) Parameter(name string, value interface{}) *AccountGroupsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AccountGroupsListRequest) Header(name string, value interface{}) *AccountGroupsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AccountGroupsListRequest) Impersonate(user string) *AccountGroupsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Fields sets the value of the 'fields' parameter.
//
// Projection
// This field contains a comma-separated list of fields you'd like to get in
// a result. No new fields can be added, only existing ones can be filtered.
// To specify a field 'id' of a structure 'organization' use 'organization.id'.
// To specify all fields of a structure use 'field_name.*'.
// To specify multiple fields use 'name,description,created_at'
func (r *AccountGroupsListRequest) Fields(value string) *AccountGroupsListRequest {
	r.fields = &value
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement. For example, in order to sort the
// account groups descending by name the value should be:
//
// ```sql
// name desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *AccountGroupsListRequest) Order(value string) *AccountGroupsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AccountGroupsListRequest) Page(value int) *AccountGroupsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause
// of an SQL statement, but using the names of the attributes of the account group
// instead of the names of the columns of a table. For example, in order to
// retrieve account groups with the name starting with Default:
//
// ```sql
// name like 'Default%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// items that the user has permission to see will be returned.
func (r *AccountGroupsListRequest) Search(value string) *AccountGroupsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AccountGroupsListRequest) Size(value int) *AccountGroupsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AccountGroupsListRequest) Send() (result *AccountGroupsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AccountGroupsListRequest) SendContext(ctx context.Context) (result *AccountGroupsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.fields != nil {
		helpers.AddValue(&query, "fields", *r.fields)
	}
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
	result = &AccountGroupsListResponse{}
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
	err = readAccountGroupsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AccountGroupsListResponse is the response for the 'list' method.
type AccountGroupsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *AccountGroupList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *AccountGroupsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AccountGroupsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AccountGroupsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of account groups.
func (r *AccountGroupsListResponse) Items() *AccountGroupList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of account groups.
func (r *AccountGroupsListResponse) GetItems() (value *AccountGroupList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AccountGroupsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AccountGroupsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AccountGroupsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *AccountGroupsListResponse) GetSize() (value int, ok bool) {
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
func (r *AccountGroupsListResponse) Total() int {
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
func (r *AccountGroupsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
