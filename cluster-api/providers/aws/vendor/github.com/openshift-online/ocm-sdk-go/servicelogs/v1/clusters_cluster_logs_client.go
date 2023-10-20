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

package v1 // github.com/openshift-online/ocm-sdk-go/servicelogs/v1

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ClustersClusterLogsClient is the client of the 'clusters_cluster_logs' resource.
//
// Manages the collection of Clusters cluster logs.
type ClustersClusterLogsClient struct {
	transport http.RoundTripper
	path      string
}

// NewClustersClusterLogsClient creates a new client for the 'clusters_cluster_logs'
// resource using the given transport to send the requests and receive the
// responses.
func NewClustersClusterLogsClient(transport http.RoundTripper, path string) *ClustersClusterLogsClient {
	return &ClustersClusterLogsClient{
		transport: transport,
		path:      path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves the list of cluster logs by the cluster_id and/or cluster_uuid parameters.
// Use this endpoint to list service logs for a specific cluster (excluding private logs).
// Any authenticated user is able to use this endpoint without any special permissions.
func (c *ClustersClusterLogsClient) List() *ClustersClusterLogsListRequest {
	return &ClustersClusterLogsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ClustersClusterLogsListRequest is the request for the 'list' method.
type ClustersClusterLogsListRequest struct {
	transport   http.RoundTripper
	path        string
	query       url.Values
	header      http.Header
	clusterID   *string
	clusterUUID *string
	order       *string
	page        *int
	search      *string
	size        *int
}

// Parameter adds a query parameter.
func (r *ClustersClusterLogsListRequest) Parameter(name string, value interface{}) *ClustersClusterLogsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ClustersClusterLogsListRequest) Header(name string, value interface{}) *ClustersClusterLogsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ClustersClusterLogsListRequest) Impersonate(user string) *ClustersClusterLogsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// ClusterID sets the value of the 'cluster_ID' parameter.
//
// cluster_id parameter.
func (r *ClustersClusterLogsListRequest) ClusterID(value string) *ClustersClusterLogsListRequest {
	r.clusterID = &value
	return r
}

// ClusterUUID sets the value of the 'cluster_UUID' parameter.
//
// cluster_uuid parameter.
func (r *ClustersClusterLogsListRequest) ClusterUUID(value string) *ClustersClusterLogsListRequest {
	r.clusterUUID = &value
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement. For example, in order to sort the
// cluster logs descending by name identifier the value should be:
//
// ```sql
// name desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *ClustersClusterLogsListRequest) Order(value string) *ClustersClusterLogsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ClustersClusterLogsListRequest) Page(value int) *ClustersClusterLogsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause
// of an SQL statement, but using the names of the attributes of the cluster logs
// instead of the names of the columns of a table. For example, in order to
// retrieve cluster logs with service_name starting with my:
//
// ```sql
// service_name like 'my%'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// items that the user has permission to see will be returned.
func (r *ClustersClusterLogsListRequest) Search(value string) *ClustersClusterLogsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *ClustersClusterLogsListRequest) Size(value int) *ClustersClusterLogsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ClustersClusterLogsListRequest) Send() (result *ClustersClusterLogsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ClustersClusterLogsListRequest) SendContext(ctx context.Context) (result *ClustersClusterLogsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.clusterID != nil {
		helpers.AddValue(&query, "cluster_id", *r.clusterID)
	}
	if r.clusterUUID != nil {
		helpers.AddValue(&query, "cluster_uuid", *r.clusterUUID)
	}
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
	result = &ClustersClusterLogsListResponse{}
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
	err = readClustersClusterLogsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ClustersClusterLogsListResponse is the response for the 'list' method.
type ClustersClusterLogsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *LogEntryList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *ClustersClusterLogsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ClustersClusterLogsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ClustersClusterLogsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of Cluster logs.
func (r *ClustersClusterLogsListResponse) Items() *LogEntryList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of Cluster logs.
func (r *ClustersClusterLogsListResponse) GetItems() (value *LogEntryList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ClustersClusterLogsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *ClustersClusterLogsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *ClustersClusterLogsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *ClustersClusterLogsListResponse) GetSize() (value int, ok bool) {
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
func (r *ClustersClusterLogsListResponse) Total() int {
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
func (r *ClustersClusterLogsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
