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

// ServiceDependenciesClient is the client of the 'service_dependencies' resource.
//
// Manages the collection of service dependencies.
type ServiceDependenciesClient struct {
	transport http.RoundTripper
	path      string
}

// NewServiceDependenciesClient creates a new client for the 'service_dependencies'
// resource using the given transport to send the requests and receive the
// responses.
func NewServiceDependenciesClient(transport http.RoundTripper, path string) *ServiceDependenciesClient {
	return &ServiceDependenciesClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
func (c *ServiceDependenciesClient) Add() *ServiceDependenciesAddRequest {
	return &ServiceDependenciesAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of service dependencies.
func (c *ServiceDependenciesClient) List() *ServiceDependenciesListRequest {
	return &ServiceDependenciesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ServiceDependency returns the target 'service_dependency' resource for the given identifier.
func (c *ServiceDependenciesClient) ServiceDependency(id string) *ServiceDependencyClient {
	return NewServiceDependencyClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ServiceDependenciesAddRequest is the request for the 'add' method.
type ServiceDependenciesAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ServiceDependency
}

// Parameter adds a query parameter.
func (r *ServiceDependenciesAddRequest) Parameter(name string, value interface{}) *ServiceDependenciesAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServiceDependenciesAddRequest) Header(name string, value interface{}) *ServiceDependenciesAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServiceDependenciesAddRequest) Impersonate(user string) *ServiceDependenciesAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ServiceDependenciesAddRequest) Body(value *ServiceDependency) *ServiceDependenciesAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServiceDependenciesAddRequest) Send() (result *ServiceDependenciesAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServiceDependenciesAddRequest) SendContext(ctx context.Context) (result *ServiceDependenciesAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeServiceDependenciesAddRequest(r, buffer)
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
	result = &ServiceDependenciesAddResponse{}
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
	err = readServiceDependenciesAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ServiceDependenciesAddResponse is the response for the 'add' method.
type ServiceDependenciesAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ServiceDependency
}

// Status returns the response status code.
func (r *ServiceDependenciesAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServiceDependenciesAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServiceDependenciesAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ServiceDependenciesAddResponse) Body() *ServiceDependency {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceDependenciesAddResponse) GetBody() (value *ServiceDependency, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ServiceDependenciesListRequest is the request for the 'list' method.
type ServiceDependenciesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	orderBy   *string
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *ServiceDependenciesListRequest) Parameter(name string, value interface{}) *ServiceDependenciesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ServiceDependenciesListRequest) Header(name string, value interface{}) *ServiceDependenciesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ServiceDependenciesListRequest) Impersonate(user string) *ServiceDependenciesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *ServiceDependenciesListRequest) OrderBy(value string) *ServiceDependenciesListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *ServiceDependenciesListRequest) Page(value int) *ServiceDependenciesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *ServiceDependenciesListRequest) Size(value int) *ServiceDependenciesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ServiceDependenciesListRequest) Send() (result *ServiceDependenciesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ServiceDependenciesListRequest) SendContext(ctx context.Context) (result *ServiceDependenciesListResponse, err error) {
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
	result = &ServiceDependenciesListResponse{}
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
	err = readServiceDependenciesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ServiceDependenciesListResponse is the response for the 'list' method.
type ServiceDependenciesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ServiceDependencyList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ServiceDependenciesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ServiceDependenciesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ServiceDependenciesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *ServiceDependenciesListResponse) Items() *ServiceDependencyList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceDependenciesListResponse) GetItems() (value *ServiceDependencyList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *ServiceDependenciesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceDependenciesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *ServiceDependenciesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceDependenciesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *ServiceDependenciesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *ServiceDependenciesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
