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
	"context"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ProductTechnologyPreviewsClient is the client of the 'product_technology_previews' resource.
//
// Manages the collection of product technology previews.
type ProductTechnologyPreviewsClient struct {
	transport http.RoundTripper
	path      string
}

// NewProductTechnologyPreviewsClient creates a new client for the 'product_technology_previews'
// resource using the given transport to send the requests and receive the
// responses.
func NewProductTechnologyPreviewsClient(transport http.RoundTripper, path string) *ProductTechnologyPreviewsClient {
	return &ProductTechnologyPreviewsClient{
		transport: transport,
		path:      path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of product technology previews.
func (c *ProductTechnologyPreviewsClient) List() *ProductTechnologyPreviewsListRequest {
	return &ProductTechnologyPreviewsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// TechnologyPreview returns the target 'product_technology_preview' resource for the given identifier.
//
// Returns a reference to the service that manages a specific product technology preview.
func (c *ProductTechnologyPreviewsClient) TechnologyPreview(id string) *ProductTechnologyPreviewClient {
	return NewProductTechnologyPreviewClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ProductTechnologyPreviewsListRequest is the request for the 'list' method.
type ProductTechnologyPreviewsListRequest struct {
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
func (r *ProductTechnologyPreviewsListRequest) Parameter(name string, value interface{}) *ProductTechnologyPreviewsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ProductTechnologyPreviewsListRequest) Header(name string, value interface{}) *ProductTechnologyPreviewsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ProductTechnologyPreviewsListRequest) Impersonate(user string) *ProductTechnologyPreviewsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement, but using the names of the attributes of the product instead of
// the names of the columns of a table. For example, in order to sort the products
// descending by name the value should be:
//
// ```sql
// name desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *ProductTechnologyPreviewsListRequest) Order(value string) *ProductTechnologyPreviewsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ProductTechnologyPreviewsListRequest) Page(value int) *ProductTechnologyPreviewsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of an
// SQL statement, but using the names of the attributes of the product instead of
// the names of the columns of a table. For example, in order to retrieve all the
// products with a name starting with `my` the value should be:
//
// ```sql
// name like 'my%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the products
// that the user has permission to see will be returned.
func (r *ProductTechnologyPreviewsListRequest) Search(value string) *ProductTechnologyPreviewsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *ProductTechnologyPreviewsListRequest) Size(value int) *ProductTechnologyPreviewsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ProductTechnologyPreviewsListRequest) Send() (result *ProductTechnologyPreviewsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ProductTechnologyPreviewsListRequest) SendContext(ctx context.Context) (result *ProductTechnologyPreviewsListResponse, err error) {
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
	result = &ProductTechnologyPreviewsListResponse{}
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
	err = readProductTechnologyPreviewsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ProductTechnologyPreviewsListResponse is the response for the 'list' method.
type ProductTechnologyPreviewsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ProductTechnologyPreviewList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ProductTechnologyPreviewsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ProductTechnologyPreviewsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ProductTechnologyPreviewsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of product technology previews.
func (r *ProductTechnologyPreviewsListResponse) Items() *ProductTechnologyPreviewList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of product technology previews.
func (r *ProductTechnologyPreviewsListResponse) GetItems() (value *ProductTechnologyPreviewList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ProductTechnologyPreviewsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ProductTechnologyPreviewsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *ProductTechnologyPreviewsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *ProductTechnologyPreviewsListResponse) GetSize() (value int, ok bool) {
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
func (r *ProductTechnologyPreviewsListResponse) Total() int {
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
func (r *ProductTechnologyPreviewsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
