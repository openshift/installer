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

// AWSInfrastructureAccessRoleGrantsClient is the client of the 'AWS_infrastructure_access_role_grants' resource.
//
// Manages the collection of AWS infrastructure access role grants.
type AWSInfrastructureAccessRoleGrantsClient struct {
	transport http.RoundTripper
	path      string
}

// NewAWSInfrastructureAccessRoleGrantsClient creates a new client for the 'AWS_infrastructure_access_role_grants'
// resource using the given transport to send the requests and receive the
// responses.
func NewAWSInfrastructureAccessRoleGrantsClient(transport http.RoundTripper, path string) *AWSInfrastructureAccessRoleGrantsClient {
	return &AWSInfrastructureAccessRoleGrantsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Create a new AWS infrastructure access role grant and add it to the collection of
// AWS infrastructure access role grants on the cluster.
func (c *AWSInfrastructureAccessRoleGrantsClient) Add() *AWSInfrastructureAccessRoleGrantsAddRequest {
	return &AWSInfrastructureAccessRoleGrantsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of AWS infrastructure access role grants.
func (c *AWSInfrastructureAccessRoleGrantsClient) List() *AWSInfrastructureAccessRoleGrantsListRequest {
	return &AWSInfrastructureAccessRoleGrantsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AWSInfrastructureAccessRoleGrant returns the target 'AWS_infrastructure_access_role_grant' resource for the given identifier.
//
// Returns a reference to the service that manages a specific AWS infrastructure access role grant.
func (c *AWSInfrastructureAccessRoleGrantsClient) AWSInfrastructureAccessRoleGrant(id string) *AWSInfrastructureAccessRoleGrantClient {
	return NewAWSInfrastructureAccessRoleGrantClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// AWSInfrastructureAccessRoleGrantsAddRequest is the request for the 'add' method.
type AWSInfrastructureAccessRoleGrantsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AWSInfrastructureAccessRoleGrant
}

// Parameter adds a query parameter.
func (r *AWSInfrastructureAccessRoleGrantsAddRequest) Parameter(name string, value interface{}) *AWSInfrastructureAccessRoleGrantsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AWSInfrastructureAccessRoleGrantsAddRequest) Header(name string, value interface{}) *AWSInfrastructureAccessRoleGrantsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AWSInfrastructureAccessRoleGrantsAddRequest) Impersonate(user string) *AWSInfrastructureAccessRoleGrantsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the AWS infrastructure access role grant.
func (r *AWSInfrastructureAccessRoleGrantsAddRequest) Body(value *AWSInfrastructureAccessRoleGrant) *AWSInfrastructureAccessRoleGrantsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AWSInfrastructureAccessRoleGrantsAddRequest) Send() (result *AWSInfrastructureAccessRoleGrantsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AWSInfrastructureAccessRoleGrantsAddRequest) SendContext(ctx context.Context) (result *AWSInfrastructureAccessRoleGrantsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAWSInfrastructureAccessRoleGrantsAddRequest(r, buffer)
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
	result = &AWSInfrastructureAccessRoleGrantsAddResponse{}
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
	err = readAWSInfrastructureAccessRoleGrantsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AWSInfrastructureAccessRoleGrantsAddResponse is the response for the 'add' method.
type AWSInfrastructureAccessRoleGrantsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *AWSInfrastructureAccessRoleGrant
}

// Status returns the response status code.
func (r *AWSInfrastructureAccessRoleGrantsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AWSInfrastructureAccessRoleGrantsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AWSInfrastructureAccessRoleGrantsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the AWS infrastructure access role grant.
func (r *AWSInfrastructureAccessRoleGrantsAddResponse) Body() *AWSInfrastructureAccessRoleGrant {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the AWS infrastructure access role grant.
func (r *AWSInfrastructureAccessRoleGrantsAddResponse) GetBody() (value *AWSInfrastructureAccessRoleGrant, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// AWSInfrastructureAccessRoleGrantsListRequest is the request for the 'list' method.
type AWSInfrastructureAccessRoleGrantsListRequest struct {
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
func (r *AWSInfrastructureAccessRoleGrantsListRequest) Parameter(name string, value interface{}) *AWSInfrastructureAccessRoleGrantsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AWSInfrastructureAccessRoleGrantsListRequest) Header(name string, value interface{}) *AWSInfrastructureAccessRoleGrantsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *AWSInfrastructureAccessRoleGrantsListRequest) Impersonate(user string) *AWSInfrastructureAccessRoleGrantsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement, but using the names of the attributes of the AWS infrastructure access role grant
// instead of the names of the columns of a table. For example, in order to sort the
// AWS infrastructure access role grants descending by user ARN the value should be:
//
// ```sql
// user_arn desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *AWSInfrastructureAccessRoleGrantsListRequest) Order(value string) *AWSInfrastructureAccessRoleGrantsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AWSInfrastructureAccessRoleGrantsListRequest) Page(value int) *AWSInfrastructureAccessRoleGrantsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of an
// SQL statement, but using the names of the attributes of the AWS infrastructure access role grant
// instead of the names of the columns of a table. For example, in order to retrieve
// all the AWS infrastructure access role grants with a user ARN starting with `user` the value should be:
//
// ```sql
// user_arn like '%user'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the AWS
// infrastructure access role grants that the user has permission to see will be returned.
func (r *AWSInfrastructureAccessRoleGrantsListRequest) Search(value string) *AWSInfrastructureAccessRoleGrantsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AWSInfrastructureAccessRoleGrantsListRequest) Size(value int) *AWSInfrastructureAccessRoleGrantsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AWSInfrastructureAccessRoleGrantsListRequest) Send() (result *AWSInfrastructureAccessRoleGrantsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AWSInfrastructureAccessRoleGrantsListRequest) SendContext(ctx context.Context) (result *AWSInfrastructureAccessRoleGrantsListResponse, err error) {
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
	result = &AWSInfrastructureAccessRoleGrantsListResponse{}
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
	err = readAWSInfrastructureAccessRoleGrantsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// AWSInfrastructureAccessRoleGrantsListResponse is the response for the 'list' method.
type AWSInfrastructureAccessRoleGrantsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *AWSInfrastructureAccessRoleGrantList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of AWS infrastructure access role grants.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) Items() *AWSInfrastructureAccessRoleGrantList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of AWS infrastructure access role grants.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) GetItems() (value *AWSInfrastructureAccessRoleGrantList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *AWSInfrastructureAccessRoleGrantsListResponse) GetSize() (value int, ok bool) {
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
func (r *AWSInfrastructureAccessRoleGrantsListResponse) Total() int {
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
func (r *AWSInfrastructureAccessRoleGrantsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
