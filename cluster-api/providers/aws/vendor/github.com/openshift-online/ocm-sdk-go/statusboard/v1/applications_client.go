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

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ApplicationsClient is the client of the 'applications' resource.
//
// Manages the collection of applications.
type ApplicationsClient struct {
	transport http.RoundTripper
	path      string
}

// NewApplicationsClient creates a new client for the 'applications'
// resource using the given transport to send the requests and receive the
// responses.
func NewApplicationsClient(transport http.RoundTripper, path string) *ApplicationsClient {
	return &ApplicationsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
func (c *ApplicationsClient) Add() *ApplicationsAddRequest {
	return &ApplicationsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of applications.
func (c *ApplicationsClient) List() *ApplicationsListRequest {
	return &ApplicationsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Application returns the target 'application' resource for the given identifier.
func (c *ApplicationsClient) Application(id string) *ApplicationClient {
	return NewApplicationClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ApplicationsAddRequest is the request for the 'add' method.
type ApplicationsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Application
}

// Parameter adds a query parameter.
func (r *ApplicationsAddRequest) Parameter(name string, value interface{}) *ApplicationsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ApplicationsAddRequest) Header(name string, value interface{}) *ApplicationsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ApplicationsAddRequest) Impersonate(user string) *ApplicationsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ApplicationsAddRequest) Body(value *Application) *ApplicationsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ApplicationsAddRequest) Send() (result *ApplicationsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ApplicationsAddRequest) SendContext(ctx context.Context) (result *ApplicationsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeApplicationsAddRequest(r, buffer)
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
	result = &ApplicationsAddResponse{}
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
	err = readApplicationsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ApplicationsAddResponse is the response for the 'add' method.
type ApplicationsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Application
}

// Status returns the response status code.
func (r *ApplicationsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ApplicationsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ApplicationsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ApplicationsAddResponse) Body() *Application {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationsAddResponse) GetBody() (value *Application, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ApplicationsListRequest is the request for the 'list' method.
type ApplicationsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	fullname  *string
	orderBy   *string
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *ApplicationsListRequest) Parameter(name string, value interface{}) *ApplicationsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ApplicationsListRequest) Header(name string, value interface{}) *ApplicationsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ApplicationsListRequest) Impersonate(user string) *ApplicationsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Fullname sets the value of the 'fullname' parameter.
func (r *ApplicationsListRequest) Fullname(value string) *ApplicationsListRequest {
	r.fullname = &value
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *ApplicationsListRequest) OrderBy(value string) *ApplicationsListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *ApplicationsListRequest) Page(value int) *ApplicationsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
func (r *ApplicationsListRequest) Search(value string) *ApplicationsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *ApplicationsListRequest) Size(value int) *ApplicationsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ApplicationsListRequest) Send() (result *ApplicationsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ApplicationsListRequest) SendContext(ctx context.Context) (result *ApplicationsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.fullname != nil {
		helpers.AddValue(&query, "fullname", *r.fullname)
	}
	if r.orderBy != nil {
		helpers.AddValue(&query, "order_by", *r.orderBy)
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
	result = &ApplicationsListResponse{}
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
	err = readApplicationsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ApplicationsListResponse is the response for the 'list' method.
type ApplicationsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ApplicationList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ApplicationsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ApplicationsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ApplicationsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *ApplicationsListResponse) Items() *ApplicationList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationsListResponse) GetItems() (value *ApplicationList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *ApplicationsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *ApplicationsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *ApplicationsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
