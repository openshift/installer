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
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// SubscriptionsClient is the client of the 'subscriptions' resource.
//
// Manages the collection of subscriptions.
type SubscriptionsClient struct {
	transport http.RoundTripper
	path      string
}

// NewSubscriptionsClient creates a new client for the 'subscriptions'
// resource using the given transport to send the requests and receive the
// responses.
func NewSubscriptionsClient(transport http.RoundTripper, path string) *SubscriptionsClient {
	return &SubscriptionsClient{
		transport: transport,
		path:      path,
	}
}

// List creates a request for the 'list' method.
//
// Retrieves a list of subscriptions.
func (c *SubscriptionsClient) List() *SubscriptionsListRequest {
	return &SubscriptionsListRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Post creates a request for the 'post' method.
//
// Create a new subscription and register a cluster for it.
func (c *SubscriptionsClient) Post() *SubscriptionsPostRequest {
	return &SubscriptionsPostRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// Labels returns the target 'generic_labels' resource.
//
// Reference to the list of labels of a specific subscription.
func (c *SubscriptionsClient) Labels() *GenericLabelsClient {
	return NewGenericLabelsClient(
		c.transport,
		path.Join(c.path, "labels"),
	)
}

// Subscription returns the target 'subscription' resource for the given identifier.
//
// Reference to the service that manages a specific subscription.
func (c *SubscriptionsClient) Subscription(id string) *SubscriptionClient {
	return NewSubscriptionClient(
		c.transport,
		path.Join(c.path, id),
	)
}

// SubscriptionsListRequest is the request for the 'list' method.
type SubscriptionsListRequest struct {
	transport     http.RoundTripper
	path          string
	query         url.Values
	header        http.Header
	fetchAccounts *bool
	fetchLabels   *bool
	fields        *string
	labels        *string
	order         *string
	page          *int
	search        *string
	size          *int
}

// Parameter adds a query parameter.
func (r *SubscriptionsListRequest) Parameter(name string, value interface{}) *SubscriptionsListRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *SubscriptionsListRequest) Header(name string, value interface{}) *SubscriptionsListRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *SubscriptionsListRequest) Impersonate(user string) *SubscriptionsListRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// FetchAccounts sets the value of the 'fetch_accounts' parameter.
//
// If true, includes the account reference information in the output. Could slow request response time.
func (r *SubscriptionsListRequest) FetchAccounts(value bool) *SubscriptionsListRequest {
	r.fetchAccounts = &value
	return r
}

// FetchLabels sets the value of the 'fetch_labels' parameter.
//
// If true, includes the labels on a subscription in the output. Could slow request response time.
func (r *SubscriptionsListRequest) FetchLabels(value bool) *SubscriptionsListRequest {
	r.fetchLabels = &value
	return r
}

// Fields sets the value of the 'fields' parameter.
//
// Projection
// This field contains a comma-separated list of fields you'd like to get in
// a result. No new fields can be added, only existing ones can be filtered.
// To specify a field 'id' of a structure 'plan' use 'plan.id'.
// To specify all fields of a structure 'labels' use 'labels.*'.
func (r *SubscriptionsListRequest) Fields(value string) *SubscriptionsListRequest {
	r.fields = &value
	return r
}

// Labels sets the value of the 'labels' parameter.
//
// Filter subscriptions by a comma separated list of labels:
//
// [source]
// ----
// env=staging,department=sales
// ----
func (r *SubscriptionsListRequest) Labels(value string) *SubscriptionsListRequest {
	r.labels = &value
	return r
}

// Order sets the value of the 'order' parameter.
//
// Order criteria.
//
// The syntax of this parameter is similar to the syntax of the _order by_ clause of
// a SQL statement. For example, in order to sort the
// subscriptions descending by name identifier the value should be:
//
// ```sql
// name desc
// ```
//
// If the parameter isn't provided, or if the value is empty, then the order of the
// results is undefined.
func (r *SubscriptionsListRequest) Order(value string) *SubscriptionsListRequest {
	r.order = &value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *SubscriptionsListRequest) Page(value int) *SubscriptionsListRequest {
	r.page = &value
	return r
}

// Search sets the value of the 'search' parameter.
//
// Search criteria.
//
// The syntax of this parameter is similar to the syntax of the _where_ clause of a
// SQL statement, but using the names of the attributes of the subscription instead
// of the names of the columns of a table. For example, in order to retrieve all the
// subscriptions for managed clusters the value should be:
//
// ```sql
// managed = 't'
// ```
//
// If the parameter isn't provided, or if the value is empty, then all the
// clusters that the user has permission to see will be returned.
func (r *SubscriptionsListRequest) Search(value string) *SubscriptionsListRequest {
	r.search = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *SubscriptionsListRequest) Size(value int) *SubscriptionsListRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *SubscriptionsListRequest) Send() (result *SubscriptionsListResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *SubscriptionsListRequest) SendContext(ctx context.Context) (result *SubscriptionsListResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.fetchAccounts != nil {
		helpers.AddValue(&query, "fetchAccounts", *r.fetchAccounts)
	}
	if r.fetchLabels != nil {
		helpers.AddValue(&query, "fetchLabels", *r.fetchLabels)
	}
	if r.fields != nil {
		helpers.AddValue(&query, "fields", *r.fields)
	}
	if r.labels != nil {
		helpers.AddValue(&query, "labels", *r.labels)
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
	result = &SubscriptionsListResponse{}
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
	err = readSubscriptionsListResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// SubscriptionsListResponse is the response for the 'list' method.
type SubscriptionsListResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *SubscriptionList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *SubscriptionsListResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *SubscriptionsListResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *SubscriptionsListResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of subscriptions.
func (r *SubscriptionsListResponse) Items() *SubscriptionList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of subscriptions.
func (r *SubscriptionsListResponse) GetItems() (value *SubscriptionList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the requested page, where one corresponds to the first page.
func (r *SubscriptionsListResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the requested page, where one corresponds to the first page.
func (r *SubscriptionsListResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Maximum number of items that will be contained in the returned page.
func (r *SubscriptionsListResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Maximum number of items that will be contained in the returned page.
func (r *SubscriptionsListResponse) GetSize() (value int, ok bool) {
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
func (r *SubscriptionsListResponse) Total() int {
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
func (r *SubscriptionsListResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}

// SubscriptionsPostRequest is the request for the 'post' method.
type SubscriptionsPostRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	request   *SubscriptionRegistration
}

// Parameter adds a query parameter.
func (r *SubscriptionsPostRequest) Parameter(name string, value interface{}) *SubscriptionsPostRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *SubscriptionsPostRequest) Header(name string, value interface{}) *SubscriptionsPostRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *SubscriptionsPostRequest) Impersonate(user string) *SubscriptionsPostRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Request sets the value of the 'request' parameter.
func (r *SubscriptionsPostRequest) Request(value *SubscriptionRegistration) *SubscriptionsPostRequest {
	r.request = value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *SubscriptionsPostRequest) Send() (result *SubscriptionsPostResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *SubscriptionsPostRequest) SendContext(ctx context.Context) (result *SubscriptionsPostResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeSubscriptionsPostRequest(r, buffer)
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
	result = &SubscriptionsPostResponse{}
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
	err = readSubscriptionsPostResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// SubscriptionsPostResponse is the response for the 'post' method.
type SubscriptionsPostResponse struct {
	status   int
	header   http.Header
	err      *errors.Error
	response *Subscription
}

// Status returns the response status code.
func (r *SubscriptionsPostResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *SubscriptionsPostResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *SubscriptionsPostResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Response returns the value of the 'response' parameter.
func (r *SubscriptionsPostResponse) Response() *Subscription {
	if r == nil {
		return nil
	}
	return r.response
}

// GetResponse returns the value of the 'response' parameter and
// a flag indicating if the parameter has a value.
func (r *SubscriptionsPostResponse) GetResponse() (value *Subscription, ok bool) {
	ok = r != nil && r.response != nil
	if ok {
		value = r.response
	}
	return
}
