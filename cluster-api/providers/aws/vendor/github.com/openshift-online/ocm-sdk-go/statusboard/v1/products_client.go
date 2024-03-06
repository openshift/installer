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

// ProductsClient is the client of the 'products' resource.
//
// Manages the collection of products.
type ProductsClient struct {
	transport http.RoundTripper
	path      string
}

// NewProductsClient creates a new client for the 'products'
// resource using the given transport to send the requests and receive the
// responses.
func NewProductsClient(transport http.RoundTripper, path string) *ProductsClient {
	return &ProductsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
func (c *ProductsClient) Add() *ProductsAddRequest {
	return &ProductsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of products.
func (c *ProductsClient) List() *ProductsListRequest {
	return &ProductsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Product returns the target 'product' resource for the given identifier.
func (c *ProductsClient) Product(id string) *ProductClient {
	return NewProductClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ProductsAddRequest is the request for the 'add' method.
type ProductsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *Product
}

// Parameter adds a query parameter.
func (r *ProductsAddRequest) Parameter(name string, value interface{}) *ProductsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ProductsAddRequest) Header(name string, value interface{}) *ProductsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ProductsAddRequest) Impersonate(user string) *ProductsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
func (r *ProductsAddRequest) Body(value *Product) *ProductsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ProductsAddRequest) Send() (result *ProductsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ProductsAddRequest) SendContext(ctx context.Context) (result *ProductsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeProductsAddRequest(r, buffer)
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
	result = &ProductsAddResponse{}
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
	err = readProductsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ProductsAddResponse is the response for the 'add' method.
type ProductsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Product
}

// Status returns the response status code.
func (r *ProductsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ProductsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ProductsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
func (r *ProductsAddResponse) Body() *Product {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
func (r *ProductsAddResponse) GetBody() (value *Product, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ProductsListRequest is the request for the 'list' method.
type ProductsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	orderBy   *string
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *ProductsListRequest) Parameter(name string, value interface{}) *ProductsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ProductsListRequest) Header(name string, value interface{}) *ProductsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ProductsListRequest) Impersonate(user string) *ProductsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// OrderBy sets the value of the 'order_by' parameter.
func (r *ProductsListRequest) OrderBy(value string) *ProductsListRequest {
	r.orderBy = &value
	return r
}

// Page sets the value of the 'page' parameter.
func (r *ProductsListRequest) Page(value int) *ProductsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
func (r *ProductsListRequest) Search(value string) *ProductsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
func (r *ProductsListRequest) Size(value int) *ProductsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ProductsListRequest) Send() (result *ProductsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ProductsListRequest) SendContext(ctx context.Context) (result *ProductsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
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
	result = &ProductsListResponse{}
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
	err = readProductsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ProductsListResponse is the response for the 'list' method.
type ProductsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ProductList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ProductsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ProductsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ProductsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
func (r *ProductsListResponse) Items() *ProductList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
func (r *ProductsListResponse) GetItems() (value *ProductList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
func (r *ProductsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
func (r *ProductsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
func (r *ProductsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
func (r *ProductsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
func (r *ProductsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
func (r *ProductsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
