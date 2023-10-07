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

// HTPasswdUsersClient is the client of the 'HT_passwd_users' resource.
//
// Manages the collection of users in an _HTPasswd_ IDP of a cluster.
type HTPasswdUsersClient struct {
	transport http.RoundTripper
	path      string
}

// NewHTPasswdUsersClient creates a new client for the 'HT_passwd_users'
// resource using the given transport to send the requests and receive the
// responses.
func NewHTPasswdUsersClient(transport http.RoundTripper, path string) *HTPasswdUsersClient {
	return &HTPasswdUsersClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a new user to the _HTPasswd_ file.
func (c *HTPasswdUsersClient) Add() *HTPasswdUsersAddRequest {
	return &HTPasswdUsersAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Import creates a request for the 'import' method.
//
// Adds multiple new users to the _HTPasswd_ file.
func (c *HTPasswdUsersClient) Import() *HTPasswdUsersImportRequest {
	return &HTPasswdUsersImportRequest{
		transport: c.transport,
		path:      path.Join(c.path, "import"),
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of _HTPasswd_ IDP users.
func (c *HTPasswdUsersClient) List() *HTPasswdUsersListRequest {
	return &HTPasswdUsersListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// HtpasswdUser returns the target 'HT_passwd_user' resource for the given identifier.
//
// Reference to the service that manages a specific _HTPasswd_ user.
func (c *HTPasswdUsersClient) HtpasswdUser(id string) *HTPasswdUserClient {
	return NewHTPasswdUserClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// HTPasswdUsersAddRequest is the request for the 'add' method.
type HTPasswdUsersAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *HTPasswdUser
}

// Parameter adds a query parameter.
func (r *HTPasswdUsersAddRequest) Parameter(name string, value interface{}) *HTPasswdUsersAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *HTPasswdUsersAddRequest) Header(name string, value interface{}) *HTPasswdUsersAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *HTPasswdUsersAddRequest) Impersonate(user string) *HTPasswdUsersAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// New user to be added
func (r *HTPasswdUsersAddRequest) Body(value *HTPasswdUser) *HTPasswdUsersAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *HTPasswdUsersAddRequest) Send() (result *HTPasswdUsersAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *HTPasswdUsersAddRequest) SendContext(ctx context.Context) (result *HTPasswdUsersAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeHTPasswdUsersAddRequest(r, buffer)
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
	result = &HTPasswdUsersAddResponse{}
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
	err = readHTPasswdUsersAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// HTPasswdUsersAddResponse is the response for the 'add' method.
type HTPasswdUsersAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *HTPasswdUser
}

// Status returns the response status code.
func (r *HTPasswdUsersAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *HTPasswdUsersAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *HTPasswdUsersAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// New user to be added
func (r *HTPasswdUsersAddResponse) Body() *HTPasswdUser {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// New user to be added
func (r *HTPasswdUsersAddResponse) GetBody() (value *HTPasswdUser, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// HTPasswdUsersImportRequest is the request for the 'import' method.
type HTPasswdUsersImportRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	items     []*HTPasswdUser
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *HTPasswdUsersImportRequest) Parameter(name string, value interface{}) *HTPasswdUsersImportRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *HTPasswdUsersImportRequest) Header(name string, value interface{}) *HTPasswdUsersImportRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *HTPasswdUsersImportRequest) Impersonate(user string) *HTPasswdUsersImportRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Items sets the value of the 'items' parameter.
//
// List of users to add to the IDP.
func (r *HTPasswdUsersImportRequest) Items(value []*HTPasswdUser) *HTPasswdUsersImportRequest {
	r.items = value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *HTPasswdUsersImportRequest) Page(value int) *HTPasswdUsersImportRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *HTPasswdUsersImportRequest) Size(value int) *HTPasswdUsersImportRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *HTPasswdUsersImportRequest) Send() (result *HTPasswdUsersImportResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *HTPasswdUsersImportRequest) SendContext(ctx context.Context) (result *HTPasswdUsersImportResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeHTPasswdUsersImportRequest(r, buffer)
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
	result = &HTPasswdUsersImportResponse{}
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
	err = readHTPasswdUsersImportResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// HTPasswdUsersImportResponse is the response for the 'import' method.
type HTPasswdUsersImportResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  []*HTPasswdUser
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *HTPasswdUsersImportResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *HTPasswdUsersImportResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *HTPasswdUsersImportResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Updated list of users of the IDP.
func (r *HTPasswdUsersImportResponse) Items() []*HTPasswdUser {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Updated list of users of the IDP.
func (r *HTPasswdUsersImportResponse) GetItems() (value []*HTPasswdUser, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *HTPasswdUsersImportResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *HTPasswdUsersImportResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *HTPasswdUsersImportResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *HTPasswdUsersImportResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *HTPasswdUsersImportResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *HTPasswdUsersImportResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}

// HTPasswdUsersListRequest is the request for the 'list' method.
type HTPasswdUsersListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *HTPasswdUsersListRequest) Parameter(name string, value interface{}) *HTPasswdUsersListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *HTPasswdUsersListRequest) Header(name string, value interface{}) *HTPasswdUsersListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *HTPasswdUsersListRequest) Impersonate(user string) *HTPasswdUsersListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *HTPasswdUsersListRequest) Page(value int) *HTPasswdUsersListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *HTPasswdUsersListRequest) Size(value int) *HTPasswdUsersListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *HTPasswdUsersListRequest) Send() (result *HTPasswdUsersListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *HTPasswdUsersListRequest) SendContext(ctx context.Context) (result *HTPasswdUsersListResponse, err error) {
	query := helpers.CopyQuery(r.query)
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
	result = &HTPasswdUsersListResponse{}
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
	err = readHTPasswdUsersListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// HTPasswdUsersListResponse is the response for the 'list' method.
type HTPasswdUsersListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *HTPasswdUserList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *HTPasswdUsersListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *HTPasswdUsersListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *HTPasswdUsersListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of users of the IDP.
func (r *HTPasswdUsersListResponse) Items() *HTPasswdUserList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of users of the IDP.
func (r *HTPasswdUsersListResponse) GetItems() (value *HTPasswdUserList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *HTPasswdUsersListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *HTPasswdUsersListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *HTPasswdUsersListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *HTPasswdUsersListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *HTPasswdUsersListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *HTPasswdUsersListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
