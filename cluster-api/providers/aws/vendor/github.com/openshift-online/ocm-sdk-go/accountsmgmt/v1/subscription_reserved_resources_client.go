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
	"context"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// SubscriptionReservedResourcesClient is the client of the 'subscription_reserved_resources' resource.
//
// Manages the collection of reserved resources by a subscription.
type SubscriptionReservedResourcesClient struct {
	transport http.RoundTripper
	path      string
}

// NewSubscriptionReservedResourcesClient creates a new client for the 'subscription_reserved_resources'
// resource using the given transport to send the requests and receive the
// responses.
func NewSubscriptionReservedResourcesClient(transport http.RoundTripper, path string) *SubscriptionReservedResourcesClient {
	return &SubscriptionReservedResourcesClient{
		transport: transport,
		path:      path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves items of the collection of reserved resources by the subscription.
func (c *SubscriptionReservedResourcesClient) List() *SubscriptionReservedResourcesListRequest {
	return &SubscriptionReservedResourcesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ReservedResource returns the target 'subscription_reserved_resource' resource for the given identifier.
//
// Reference to the resource that manages the a specific resource reserved by a
// subscription.
func (c *SubscriptionReservedResourcesClient) ReservedResource(id string) *SubscriptionReservedResourceClient {
	return NewSubscriptionReservedResourceClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// SubscriptionReservedResourcesListRequest is the request for the 'list' method.
type SubscriptionReservedResourcesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *SubscriptionReservedResourcesListRequest) Parameter(name string, value interface{}) *SubscriptionReservedResourcesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *SubscriptionReservedResourcesListRequest) Header(name string, value interface{}) *SubscriptionReservedResourcesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *SubscriptionReservedResourcesListRequest) Impersonate(user string) *SubscriptionReservedResourcesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *SubscriptionReservedResourcesListRequest) Page(value int) *SubscriptionReservedResourcesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *SubscriptionReservedResourcesListRequest) Size(value int) *SubscriptionReservedResourcesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *SubscriptionReservedResourcesListRequest) Send() (result *SubscriptionReservedResourcesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *SubscriptionReservedResourcesListRequest) SendContext(ctx context.Context) (result *SubscriptionReservedResourcesListResponse, err error) {
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
	result = &SubscriptionReservedResourcesListResponse{}
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
	err = readSubscriptionReservedResourcesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// SubscriptionReservedResourcesListResponse is the response for the 'list' method.
type SubscriptionReservedResourcesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ReservedResourceList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *SubscriptionReservedResourcesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *SubscriptionReservedResourcesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *SubscriptionReservedResourcesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of reserved resources.
func (r *SubscriptionReservedResourcesListResponse) Items() *ReservedResourceList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of reserved resources.
func (r *SubscriptionReservedResourcesListResponse) GetItems() (value *ReservedResourceList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *SubscriptionReservedResourcesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *SubscriptionReservedResourcesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *SubscriptionReservedResourcesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *SubscriptionReservedResourcesListResponse) GetSize() (value int, ok bool) {
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
func (r *SubscriptionReservedResourcesListResponse) Total() int {
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
func (r *SubscriptionReservedResourcesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
