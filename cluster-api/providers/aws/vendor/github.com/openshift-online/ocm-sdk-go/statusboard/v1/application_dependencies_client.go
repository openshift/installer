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

// ApplicationDependenciesClient is the client of the 'application_dependencies' resource.
//
// Manages the collection of applications.
type ApplicationDependenciesClient struct {
	transport http.RoundTripper
	path      string
}

// NewApplicationDependenciesClient creates a new client for the 'application_dependencies'
// resource using the given transport to send the requests and receive the
// responses.
func NewApplicationDependenciesClient(transport http.RoundTripper, path string) *ApplicationDependenciesClient {
	return &ApplicationDependenciesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
func (c *ApplicationDependenciesClient) Add() *ApplicationDependenciesAddRequest {
	return &ApplicationDependenciesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of application dependencies.
func (c *ApplicationDependenciesClient) List() *ApplicationDependenciesListRequest {
	return &ApplicationDependenciesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ApplicationDependency returns the target 'application_dependency' resource for the given identifier.
func (c *ApplicationDependenciesClient) ApplicationDependency(id string) *ApplicationDependencyClient {
	return NewApplicationDependencyClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ApplicationDependenciesAddRequest is the request for the 'add' method.
type ApplicationDependenciesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ApplicationDependency
}

// Parameter adds a query parameter.
func (r *ApplicationDependenciesAddRequest) Parameter(name string, value interface{}) *ApplicationDependenciesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ApplicationDependenciesAddRequest) Header(name string, value interface{}) *ApplicationDependenciesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ApplicationDependenciesAddRequest) Impersonate(user string) *ApplicationDependenciesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ApplicationDependenciesAddRequest) Body(value *ApplicationDependency) *ApplicationDependenciesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ApplicationDependenciesAddRequest) Send() (result *ApplicationDependenciesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ApplicationDependenciesAddRequest) SendContext(ctx context.Context) (result *ApplicationDependenciesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeApplicationDependenciesAddRequest(r, buffer)
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
	result = &ApplicationDependenciesAddResponse{}
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
	err = readApplicationDependenciesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ApplicationDependenciesAddResponse is the response for the 'add' method.
type ApplicationDependenciesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ApplicationDependency
}

// Status returns the response status code.
func (r *ApplicationDependenciesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ApplicationDependenciesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ApplicationDependenciesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ApplicationDependenciesAddResponse) Body() *ApplicationDependency {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationDependenciesAddResponse) GetBody() (value *ApplicationDependency, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ApplicationDependenciesListRequest is the request for the 'list' method.
type ApplicationDependenciesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	orderBy   *string
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *ApplicationDependenciesListRequest) Parameter(name string, value interface{}) *ApplicationDependenciesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ApplicationDependenciesListRequest) Header(name string, value interface{}) *ApplicationDependenciesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ApplicationDependenciesListRequest) Impersonate(user string) *ApplicationDependenciesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *ApplicationDependenciesListRequest) OrderBy(value string) *ApplicationDependenciesListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *ApplicationDependenciesListRequest) Page(value int) *ApplicationDependenciesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *ApplicationDependenciesListRequest) Size(value int) *ApplicationDependenciesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ApplicationDependenciesListRequest) Send() (result *ApplicationDependenciesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ApplicationDependenciesListRequest) SendContext(ctx context.Context) (result *ApplicationDependenciesListResponse, err error) {
	query := helpers.CopyQuery(r.query)
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
	result = &ApplicationDependenciesListResponse{}
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
	err = readApplicationDependenciesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ApplicationDependenciesListResponse is the response for the 'list' method.
type ApplicationDependenciesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ApplicationDependencyList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ApplicationDependenciesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ApplicationDependenciesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ApplicationDependenciesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *ApplicationDependenciesListResponse) Items() *ApplicationDependencyList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationDependenciesListResponse) GetItems() (value *ApplicationDependencyList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *ApplicationDependenciesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationDependenciesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *ApplicationDependenciesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationDependenciesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *ApplicationDependenciesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *ApplicationDependenciesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
