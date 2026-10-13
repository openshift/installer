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

// DNSDomainsClient is the client of the 'DNS_domains' resource.
//
// Manages the collection of DNS domains.
type DNSDomainsClient struct {
	transport http.RoundTripper
	path      string
}

// NewDNSDomainsClient creates a new client for the 'DNS_domains'
// resource using the given transport to send the requests and receive the
// responses.
func NewDNSDomainsClient(transport http.RoundTripper, path string) *DNSDomainsClient {
	return &DNSDomainsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a DNS domain.
func (c *DNSDomainsClient) Add() *DNSDomainsAddRequest {
	return &DNSDomainsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
func (c *DNSDomainsClient) List() *DNSDomainsListRequest {
	return &DNSDomainsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// DNSDomain returns the target 'DNS_domain' resource for the given identifier.
//
// Reference to the resource that manages a specific DNS doamin.
func (c *DNSDomainsClient) DNSDomain(id string) *DNSDomainClient {
	return NewDNSDomainClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// DNSDomainsAddRequest is the request for the 'add' method.
type DNSDomainsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *DNSDomain
}

// Parameter adds a query parameter.
func (r *DNSDomainsAddRequest) Parameter(name string, value interface{}) *DNSDomainsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DNSDomainsAddRequest) Header(name string, value interface{}) *DNSDomainsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DNSDomainsAddRequest) Impersonate(user string) *DNSDomainsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the DNS domain.
func (r *DNSDomainsAddRequest) Body(value *DNSDomain) *DNSDomainsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DNSDomainsAddRequest) Send() (result *DNSDomainsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DNSDomainsAddRequest) SendContext(ctx context.Context) (result *DNSDomainsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeDNSDomainsAddRequest(r, buffer)
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
	result = &DNSDomainsAddResponse{}
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
	err = readDNSDomainsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DNSDomainsAddResponse is the response for the 'add' method.
type DNSDomainsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *DNSDomain
}

// Status returns the response status code.
func (r *DNSDomainsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DNSDomainsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DNSDomainsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the DNS domain.
func (r *DNSDomainsAddResponse) Body() *DNSDomain {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the DNS domain.
func (r *DNSDomainsAddResponse) GetBody() (value *DNSDomain, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// DNSDomainsListRequest is the request for the 'list' method.
type DNSDomainsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	search    *string
	size      *int
}

// Parameter adds a query parameter.
func (r *DNSDomainsListRequest) Parameter(name string, value interface{}) *DNSDomainsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *DNSDomainsListRequest) Header(name string, value interface{}) *DNSDomainsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *DNSDomainsListRequest) Impersonate(user string) *DNSDomainsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *DNSDomainsListRequest) Page(value int) *DNSDomainsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of a
// SQL statement, but using the names of the attributes of the dns domain instead of
// the names of the columns of a table. For example, in order to retrieve all the
// dns domains with a ID starting with `02a5` should be:
//
// ```sql
// id like '02a5%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// dns domains that the user has permission to see will be returned.
func (r *DNSDomainsListRequest) Search(value string) *DNSDomainsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *DNSDomainsListRequest) Size(value int) *DNSDomainsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *DNSDomainsListRequest) Send() (result *DNSDomainsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *DNSDomainsListRequest) SendContext(ctx context.Context) (result *DNSDomainsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
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
	result = &DNSDomainsListResponse{}
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
	err = readDNSDomainsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// DNSDomainsListResponse is the response for the 'list' method.
type DNSDomainsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *DNSDomainList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *DNSDomainsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *DNSDomainsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *DNSDomainsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved a list of DNS domains.
func (r *DNSDomainsListResponse) Items() *DNSDomainList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved a list of DNS domains.
func (r *DNSDomainsListResponse) GetItems() (value *DNSDomainList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *DNSDomainsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *DNSDomainsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *DNSDomainsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *DNSDomainsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection.
func (r *DNSDomainsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection.
func (r *DNSDomainsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
