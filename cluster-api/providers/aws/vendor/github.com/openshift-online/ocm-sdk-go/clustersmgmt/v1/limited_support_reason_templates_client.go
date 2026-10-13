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

// LimitedSupportReasonTemplatesClient is the client of the 'limited_support_reason_templates' resource.
//
// Manages the collection of limited support reason templates.
type LimitedSupportReasonTemplatesClient struct {
	transport http.RoundTripper
	path      string
}

// NewLimitedSupportReasonTemplatesClient creates a new client for the 'limited_support_reason_templates'
// resource using the given transport to send the requests and receive the
// responses.
func NewLimitedSupportReasonTemplatesClient(transport http.RoundTripper, path string) *LimitedSupportReasonTemplatesClient {
	return &LimitedSupportReasonTemplatesClient{
		transport: transport,
		path:      path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of templates.
func (c *LimitedSupportReasonTemplatesClient) List() *LimitedSupportReasonTemplatesListRequest {
	return &LimitedSupportReasonTemplatesListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// LimitedSupportReasonTemplate returns the target 'limited_support_reason_template' resource for the given identifier.
//
// Reference to the service that manages an specific template.
func (c *LimitedSupportReasonTemplatesClient) LimitedSupportReasonTemplate(id string) *LimitedSupportReasonTemplateClient {
	return NewLimitedSupportReasonTemplateClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// LimitedSupportReasonTemplatesListRequest is the request for the 'list' method.
type LimitedSupportReasonTemplatesListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *LimitedSupportReasonTemplatesListRequest) Parameter(name string, value interface{}) *LimitedSupportReasonTemplatesListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *LimitedSupportReasonTemplatesListRequest) Header(name string, value interface{}) *LimitedSupportReasonTemplatesListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *LimitedSupportReasonTemplatesListRequest) Impersonate(user string) *LimitedSupportReasonTemplatesListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *LimitedSupportReasonTemplatesListRequest) Page(value int) *LimitedSupportReasonTemplatesListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *LimitedSupportReasonTemplatesListRequest) Size(value int) *LimitedSupportReasonTemplatesListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *LimitedSupportReasonTemplatesListRequest) Send() (result *LimitedSupportReasonTemplatesListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *LimitedSupportReasonTemplatesListRequest) SendContext(ctx context.Context) (result *LimitedSupportReasonTemplatesListResponse, err error) {
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
	result = &LimitedSupportReasonTemplatesListResponse{}
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
	err = readLimitedSupportReasonTemplatesListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// LimitedSupportReasonTemplatesListResponse is the response for the 'list' method.
type LimitedSupportReasonTemplatesListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *LimitedSupportReasonTemplateList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *LimitedSupportReasonTemplatesListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *LimitedSupportReasonTemplatesListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *LimitedSupportReasonTemplatesListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of template.
func (r *LimitedSupportReasonTemplatesListResponse) Items() *LimitedSupportReasonTemplateList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of template.
func (r *LimitedSupportReasonTemplatesListResponse) GetItems() (value *LimitedSupportReasonTemplateList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *LimitedSupportReasonTemplatesListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *LimitedSupportReasonTemplatesListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items contained in the returned page.
func (r *LimitedSupportReasonTemplatesListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items contained in the returned page.
func (r *LimitedSupportReasonTemplatesListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *LimitedSupportReasonTemplatesListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *LimitedSupportReasonTemplatesListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
