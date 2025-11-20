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

// ClusterMigrationsClient is the client of the 'cluster_migrations' resource.
//
// Manages the collection of cluster migrations of a cluster.
type ClusterMigrationsClient struct {
	transport http.RoundTripper
	path      string
}

// NewClusterMigrationsClient creates a new client for the 'cluster_migrations'
// resource using the given transport to send the requests and receive the
// responses.
func NewClusterMigrationsClient(transport http.RoundTripper, path string) *ClusterMigrationsClient {
	return &ClusterMigrationsClient{
		transport: transport,
		path:      path,
	}
}

// Add creates a request for the 'add' method.
//
// Adds a cluster migration to the database.
func (c *ClusterMigrationsClient) Add() *ClusterMigrationsAddRequest {
	return &ClusterMigrationsAddRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// List creates a request for the 'list' method.
func (c *ClusterMigrationsClient) List() *ClusterMigrationsListRequest {
	return &ClusterMigrationsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Migration returns the target 'cluster_migration' resource for the given identifier.
//
// Returns a reference to the service that manages a specific cluster migration.
func (c *ClusterMigrationsClient) Migration(id string) *ClusterMigrationClient {
	return NewClusterMigrationClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// ClusterMigrationsAddRequest is the request for the 'add' method.
type ClusterMigrationsAddRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *ClusterMigration
}

// Parameter adds a query parameter.
func (r *ClusterMigrationsAddRequest) Parameter(name string, value interface{}) *ClusterMigrationsAddRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterMigrationsAddRequest) Header(name string, value interface{}) *ClusterMigrationsAddRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterMigrationsAddRequest) Impersonate(user string) *ClusterMigrationsAddRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Body sets the value of the 'body' parameter.
//
// Description of the cluster migration.
func (r *ClusterMigrationsAddRequest) Body(value *ClusterMigration) *ClusterMigrationsAddRequest {
	r.body = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterMigrationsAddRequest) Send() (result *ClusterMigrationsAddResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterMigrationsAddRequest) SendContext(ctx context.Context) (result *ClusterMigrationsAddResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeClusterMigrationsAddRequest(r, buffer)
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
	result = &ClusterMigrationsAddResponse{}
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
	err = readClusterMigrationsAddResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ClusterMigrationsAddResponse is the response for the 'add' method.
type ClusterMigrationsAddResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ClusterMigration
}

// Status returns the response status code.
func (r *ClusterMigrationsAddResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterMigrationsAddResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterMigrationsAddResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// Description of the cluster migration.
func (r *ClusterMigrationsAddResponse) Body() *ClusterMigration {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// Description of the cluster migration.
func (r *ClusterMigrationsAddResponse) GetBody() (value *ClusterMigration, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}

// ClusterMigrationsListRequest is the request for the 'list' method.
type ClusterMigrationsListRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *ClusterMigrationsListRequest) Parameter(name string, value interface{}) *ClusterMigrationsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClusterMigrationsListRequest) Header(name string, value interface{}) *ClusterMigrationsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClusterMigrationsListRequest) Impersonate(user string) *ClusterMigrationsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the returned page, where one corresponds to the first page.
func (r *ClusterMigrationsListRequest) Page(value int) *ClusterMigrationsListRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items that will be contained in the returned page.
func (r *ClusterMigrationsListRequest) Size(value int) *ClusterMigrationsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClusterMigrationsListRequest) Send() (result *ClusterMigrationsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClusterMigrationsListRequest) SendContext(ctx context.Context) (result *ClusterMigrationsListResponse, err error) {
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
	result = &ClusterMigrationsListResponse{}
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
	err = readClusterMigrationsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ClusterMigrationsListResponse is the response for the 'list' method.
type ClusterMigrationsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *ClusterMigrationList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ClusterMigrationsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClusterMigrationsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClusterMigrationsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of cluster migrations.
func (r *ClusterMigrationsListResponse) Items() *ClusterMigrationList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of cluster migrations.
func (r *ClusterMigrationsListResponse) GetItems() (value *ClusterMigrationList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the returned page, where one corresponds to the first page.
func (r *ClusterMigrationsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the returned page, where one corresponds to the first page.
func (r *ClusterMigrationsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items that will be contained in the returned page.
func (r *ClusterMigrationsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items that will be contained in the returned page.
func (r *ClusterMigrationsListResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection that match the search criteria,
// regardless of the size of the page. As this collection doesn't support paging or
// searching the result will always be the total number of migrations of the cluster.
func (r *ClusterMigrationsListResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection that match the search criteria,
// regardless of the size of the page. As this collection doesn't support paging or
// searching the result will always be the total number of migrations of the cluster.
func (r *ClusterMigrationsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
