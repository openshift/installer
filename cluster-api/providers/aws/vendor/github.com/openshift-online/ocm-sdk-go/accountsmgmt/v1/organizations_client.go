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

// OrganizationsClient is the client of the 'organizations' resource.
//
// Manages the collection of organizations.
type OrganizationsClient struct {
	transport http.RoundTripper
	path      string
}

// NewOrganizationsClient creates a new client for the 'organizations'
// resource using the given transport to send the requests and receive the
// responses.
func NewOrganizationsClient(transport http.RoundTripper, path string) *OrganizationsClient {
	return &OrganizationsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Creates a new organization.
func (c *OrganizationsClient) Add() *OrganizationsAddRequest {
	return &OrganizationsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves a list of organizations.
func (c *OrganizationsClient) List() *OrganizationsListRequest {
	return &OrganizationsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Organization returns the target 'organization' resource for the given identifier.
//
// Reference to the service that manages a specific organization.
func (c *OrganizationsClient) Organization(id string) *OrganizationClient {
	return NewOrganizationClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// OrganizationsAddRequest is the request for the 'add' method.
type OrganizationsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Organization
}

// Parameter adds a query parameter.
func (r *OrganizationsAddRequest) Parameter(name string, value interface{}) *OrganizationsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OrganizationsAddRequest) Header(name string, value interface{}) *OrganizationsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OrganizationsAddRequest) Impersonate(user string) *OrganizationsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Organization data.
func (r *OrganizationsAddRequest) Body(value *Organization) *OrganizationsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OrganizationsAddRequest) Send() (result *OrganizationsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OrganizationsAddRequest) SendContext(ctx context.Context) (result *OrganizationsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeOrganizationsAddRequest(r, buffer)
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
	result = &OrganizationsAddResponse{}
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
	err = readOrganizationsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// OrganizationsAddResponse is the response for the 'add' method.
type OrganizationsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Organization
}

// Status returns the response status code.
func (r *OrganizationsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OrganizationsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OrganizationsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Organization data.
func (r *OrganizationsAddResponse) Body() *Organization {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Organization data.
func (r *OrganizationsAddResponse) GetBody() (value *Organization, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// OrganizationsListRequest is the request for the 'list' method.
type OrganizationsListRequest struct {
	transport   http.RoundTripper
	path        string
	query       url.Values
	header      http.Header
	fetchLabels *bool
	fields      *string
	page        *int
	search      *string
	size        *int
}

// Parameter adds a query parameter.
func (r *OrganizationsListRequest) Parameter(name string, value interface{}) *OrganizationsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *OrganizationsListRequest) Header(name string, value interface{}) *OrganizationsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *OrganizationsListRequest) Impersonate(user string) *OrganizationsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// FetchLabels sets the value of the 'fetch_labels' parameter.
//
// If true, includes the labels on an organization in the output. Could slow request response time.
func (r *OrganizationsListRequest) FetchLabels(value bool) *OrganizationsListRequest {
	r.fetchLabels = &value
	return r
}

// Fields sets the value of the 'fields' parameter.
//
// Projection
// This field contains a comma-separated list of fields you'd like to get in
// a result. No new fields can be added, only existing ones can be filtered.
// To specify a field 'id' of a structure 'plan' use 'plan.id'.
// To specify all fields of a structure 'labels' use 'labels.*'.
func (r *OrganizationsListRequest) Fields(value string) *OrganizationsListRequest {
	r.fields = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OrganizationsListRequest) Page(value int) *OrganizationsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause
// of an SQL statement, but using the names of the attributes of the organization
// instead of the names of the columns of a table. For example, in order to
// retrieve organizations with name starting with my:
//
// ```sql
// name like 'my%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// items that the user has permission to see will be returned.
func (r *OrganizationsListRequest) Search(value string) *OrganizationsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *OrganizationsListRequest) Size(value int) *OrganizationsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *OrganizationsListRequest) Send() (result *OrganizationsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *OrganizationsListRequest) SendContext(ctx context.Context) (result *OrganizationsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.fetchLabels != nil {
		helpers.AddValue(&query, "fetchLabels", *r.fetchLabels)
	}
	if r.fields != nil {
		helpers.AddValue(&query, "fields", *r.fields)
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
	result = &OrganizationsListResponse{}
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
	err = readOrganizationsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// OrganizationsListResponse is the response for the 'list' method.
type OrganizationsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *OrganizationList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *OrganizationsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *OrganizationsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *OrganizationsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of organizations.
func (r *OrganizationsListResponse) Items() *OrganizationList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of organizations.
func (r *OrganizationsListResponse) GetItems() (value *OrganizationList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OrganizationsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *OrganizationsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *OrganizationsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *OrganizationsListResponse) GetSize() (value int, ok bool) {
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
func (r *OrganizationsListResponse) Total() int {
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
func (r *OrganizationsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
