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

// AddOnInstallationsClient is the client of the 'add_on_installations' resource.
//
// Manages the collection of add-on installations.
type AddOnInstallationsClient struct {
	transport http.RoundTripper
	path      string
}

// NewAddOnInstallationsClient creates a new client for the 'add_on_installations'
// resource using the given transport to send the requests and receive the
// responses.
func NewAddOnInstallationsClient(transport http.RoundTripper, path string) *AddOnInstallationsClient {
	return &AddOnInstallationsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Create a new add-on installation and add it to the collection of add-on installations on the cluster.
func (c *AddOnInstallationsClient) Add() *AddOnInstallationsAddRequest {
	return &AddOnInstallationsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of add-on installations.
func (c *AddOnInstallationsClient) List() *AddOnInstallationsListRequest {
	return &AddOnInstallationsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Addoninstallation returns the target 'add_on_installation' resource for the given identifier.
//
// Returns a reference to the service that manages a specific add-on installation.
func (c *AddOnInstallationsClient) Addoninstallation(id string) *AddOnInstallationClient {
	return NewAddOnInstallationClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// AddOnInstallationsAddRequest is the request for the 'add' method.
type AddOnInstallationsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AddOnInstallation
}

// Parameter adds a query parameter.
func (r *AddOnInstallationsAddRequest) Parameter(name string, value interface{}) *AddOnInstallationsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnInstallationsAddRequest) Header(name string, value interface{}) *AddOnInstallationsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnInstallationsAddRequest) Impersonate(user string) *AddOnInstallationsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the add-on installation.
func (r *AddOnInstallationsAddRequest) Body(value *AddOnInstallation) *AddOnInstallationsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnInstallationsAddRequest) Send() (result *AddOnInstallationsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnInstallationsAddRequest) SendContext(ctx context.Context) (result *AddOnInstallationsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAddOnInstallationsAddRequest(r, buffer)
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
	result = &AddOnInstallationsAddResponse{}
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
	err = readAddOnInstallationsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddOnInstallationsAddResponse is the response for the 'add' method.
type AddOnInstallationsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AddOnInstallation
}

// Status returns the response status code.
func (r *AddOnInstallationsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnInstallationsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnInstallationsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the add-on installation.
func (r *AddOnInstallationsAddResponse) Body() *AddOnInstallation {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the add-on installation.
func (r *AddOnInstallationsAddResponse) GetBody() (value *AddOnInstallation, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AddOnInstallationsListRequest is the request for the 'list' method.
type AddOnInstallationsListRequest struct {
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
func (r *AddOnInstallationsListRequest) Parameter(name string, value interface{}) *AddOnInstallationsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AddOnInstallationsListRequest) Header(name string, value interface{}) *AddOnInstallationsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AddOnInstallationsListRequest) Impersonate(user string) *AddOnInstallationsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement, but using the names of the attributes of the add-on installation
// instead of the names of the columns of a table. For example, in order to sort the
// add-on installations descending by name the value should be:
//
// ```sql
// name desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *AddOnInstallationsListRequest) Order(value string) *AddOnInstallationsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddOnInstallationsListRequest) Page(value int) *AddOnInstallationsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of an
// SQL statement, but using the names of the attributes of the add-on installation
// instead of the names of the columns of a table. For example, in order to retrieve
// all the add-on installations with a name starting with `my` the value should be:
//
// ```sql
// name like 'my%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the add-on
// installations that the user has permission to see will be returned.
func (r *AddOnInstallationsListRequest) Search(value string) *AddOnInstallationsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddOnInstallationsListRequest) Size(value int) *AddOnInstallationsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AddOnInstallationsListRequest) Send() (result *AddOnInstallationsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AddOnInstallationsListRequest) SendContext(ctx context.Context) (result *AddOnInstallationsListResponse, err error) {
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
	result = &AddOnInstallationsListResponse{}
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
	err = readAddOnInstallationsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AddOnInstallationsListResponse is the response for the 'list' method.
type AddOnInstallationsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *AddOnInstallationList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *AddOnInstallationsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AddOnInstallationsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AddOnInstallationsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of add-on installations.
func (r *AddOnInstallationsListResponse) Items() *AddOnInstallationList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of add-on installations.
func (r *AddOnInstallationsListResponse) GetItems() (value *AddOnInstallationList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddOnInstallationsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AddOnInstallationsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddOnInstallationsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *AddOnInstallationsListResponse) GetSize() (value int, ok bool) {
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
func (r *AddOnInstallationsListResponse) Total() int {
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
func (r *AddOnInstallationsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
