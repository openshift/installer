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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// ManagedIdentitiesRequirementsClient is the client of the 'managed_identities_requirements' resource.
//
// Manages the ManagedIdentitiesRequirements resource.
type ManagedIdentitiesRequirementsClient struct {
	transport http.RoundTripper
	path      string
}

// NewManagedIdentitiesRequirementsClient creates a new client for the 'managed_identities_requirements'
// resource using the given transport to send the requests and receive the
// responses.
func NewManagedIdentitiesRequirementsClient(transport http.RoundTripper, path string) *ManagedIdentitiesRequirementsClient {
	return &ManagedIdentitiesRequirementsClient{
		transport: transport,
		path:      path,
	}
}

// Get creates a request for the 'get' method.
//
// Retrieves an ManagedIdentitiesRequirements by version query param.
func (c *ManagedIdentitiesRequirementsClient) Get() *ManagedIdentitiesRequirementsGetRequest {
	return &ManagedIdentitiesRequirementsGetRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// ManagedIdentitiesRequirementsPollRequest is the request for the Poll method.
type ManagedIdentitiesRequirementsPollRequest struct {
	request    *ManagedIdentitiesRequirementsGetRequest
	interval   time.Duration
	statuses   []int
	predicates []func(interface{}) bool
}

// Parameter adds a query parameter to all the requests that will be used to retrieve the object.
func (r *ManagedIdentitiesRequirementsPollRequest) Parameter(name string, value interface{}) *ManagedIdentitiesRequirementsPollRequest {
	r.request.Parameter(name, value)
	return r
}

// Header adds a request header to all the requests that will be used to retrieve the object.
func (r *ManagedIdentitiesRequirementsPollRequest) Header(name string, value interface{}) *ManagedIdentitiesRequirementsPollRequest {
	r.request.Header(name, value)
	return r
}

// Required sets the value of the 'required' parameter for all the requests that
// will be used to retrieve the object.
//
// Get the managed identities requirements depending on when they are required.
// The query parameter is optional, it needs to be either ("always" or "on_enablement").
// When not supplied, this enablement constraint won't be taken into account.
// When supplied and among the accepted values, the query parameter will be used to return all managed identities requirements
// that matches the value given in the query parameter.
// When supplied but the value is invalid, an error is going to be returned.
func (r *ManagedIdentitiesRequirementsPollRequest) Required(value string) *ManagedIdentitiesRequirementsPollRequest {
	r.request.Required(value)
	return r
}

// Version sets the value of the 'version' parameter for all the requests that
// will be used to retrieve the object.
//
// Get the managed identities requirements by OpenShift version.
// The query parameter is optional, but when supplied it needs to be
// in the format X.Y (e.g 4.18) where X and Y are major and minor segments of
// the OpenShift version respectively.
// When supplied, the returned response will include all the control plane
// and data plane operators requirements for the given version.
// If not supplied, the OpenShift version constraint won't be taken into account
// when returning the managed identities requirements.
func (r *ManagedIdentitiesRequirementsPollRequest) Version(value string) *ManagedIdentitiesRequirementsPollRequest {
	r.request.Version(value)
	return r
}

// Interval sets the polling interval. This parameter is mandatory and must be greater than zero.
func (r *ManagedIdentitiesRequirementsPollRequest) Interval(value time.Duration) *ManagedIdentitiesRequirementsPollRequest {
	r.interval = value
	return r
}

// Status set the expected status of the response. Multiple values can be set calling this method
// multiple times. The response will be considered successful if the status is any of those values.
func (r *ManagedIdentitiesRequirementsPollRequest) Status(value int) *ManagedIdentitiesRequirementsPollRequest {
	r.statuses = append(r.statuses, value)
	return r
}

// Predicate adds a predicate that the response should satisfy be considered successful. Multiple
// predicates can be set calling this method multiple times. The response will be considered successful
// if all the predicates are satisfied.
func (r *ManagedIdentitiesRequirementsPollRequest) Predicate(value func(*ManagedIdentitiesRequirementsGetResponse) bool) *ManagedIdentitiesRequirementsPollRequest {
	r.predicates = append(r.predicates, func(response interface{}) bool {
		return value(response.(*ManagedIdentitiesRequirementsGetResponse))
	})
	return r
}

// StartContext starts the polling loop. Responses will be considered successful if the status is one of
// the values specified with the Status method and if all the predicates specified with the Predicate
// method return nil.
//
// The context must have a timeout or deadline, otherwise this method will immediately return an error.
func (r *ManagedIdentitiesRequirementsPollRequest) StartContext(ctx context.Context) (response *ManagedIdentitiesRequirementsPollResponse, err error) {
	result, err := helpers.PollContext(ctx, r.interval, r.statuses, r.predicates, r.task)
	if result != nil {
		response = &ManagedIdentitiesRequirementsPollResponse{
			response: result.(*ManagedIdentitiesRequirementsGetResponse),
		}
	}
	return
}

// task adapts the types of the request/response types so that they can be used with the generic
// polling function from the helpers package.
func (r *ManagedIdentitiesRequirementsPollRequest) task(ctx context.Context) (status int, result interface{}, err error) {
	response, err := r.request.SendContext(ctx)
	if response != nil {
		status = response.Status()
		result = response
	}
	return
}

// ManagedIdentitiesRequirementsPollResponse is the response for the Poll method.
type ManagedIdentitiesRequirementsPollResponse struct {
	response *ManagedIdentitiesRequirementsGetResponse
}

// Status returns the response status code.
func (r *ManagedIdentitiesRequirementsPollResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.response.Status()
}

// Header returns header of the response.
func (r *ManagedIdentitiesRequirementsPollResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.response.Header()
}

// Error returns the response error.
func (r *ManagedIdentitiesRequirementsPollResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.response.Error()
}

// Body returns the value of the 'body' parameter.
//
// ManagedIdentitiesRequirements status response.
func (r *ManagedIdentitiesRequirementsPollResponse) Body() *ManagedIdentitiesRequirements {
	return r.response.Body()
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// ManagedIdentitiesRequirements status response.
func (r *ManagedIdentitiesRequirementsPollResponse) GetBody() (value *ManagedIdentitiesRequirements, ok bool) {
	return r.response.GetBody()
}

// Poll creates a request to repeatedly retrieve the object till the response has one of a given set
// of states and satisfies a set of predicates.
func (c *ManagedIdentitiesRequirementsClient) Poll() *ManagedIdentitiesRequirementsPollRequest {
	return &ManagedIdentitiesRequirementsPollRequest{
		request: c.Get(),
	}
}

// ManagedIdentitiesRequirementsGetRequest is the request for the 'get' method.
type ManagedIdentitiesRequirementsGetRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	required  *string
	version   *string
}

// Parameter adds a query parameter.
func (r *ManagedIdentitiesRequirementsGetRequest) Parameter(name string, value interface{}) *ManagedIdentitiesRequirementsGetRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *ManagedIdentitiesRequirementsGetRequest) Header(name string, value interface{}) *ManagedIdentitiesRequirementsGetRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Impersonate wraps requests on behalf of another user.
// Note: Services that do not support this feature may silently ignore this call.
func (r *ManagedIdentitiesRequirementsGetRequest) Impersonate(user string) *ManagedIdentitiesRequirementsGetRequest {
	helpers.AddImpersonationHeader(&r.header, user)
	return r
}

// Required sets the value of the 'required' parameter.
//
// Get the managed identities requirements depending on when they are required.
// The query parameter is optional, it needs to be either ("always" or "on_enablement").
// When not supplied, this enablement constraint won't be taken into account.
// When supplied and among the accepted values, the query parameter will be used to return all managed identities requirements
// that matches the value given in the query parameter.
// When supplied but the value is invalid, an error is going to be returned.
func (r *ManagedIdentitiesRequirementsGetRequest) Required(value string) *ManagedIdentitiesRequirementsGetRequest {
	r.required = &value
	return r
}

// Version sets the value of the 'version' parameter.
//
// Get the managed identities requirements by OpenShift version.
// The query parameter is optional, but when supplied it needs to be
// in the format X.Y (e.g 4.18) where X and Y are major and minor segments of
// the OpenShift version respectively.
// When supplied, the returned response will include all the control plane
// and data plane operators requirements for the given version.
// If not supplied, the OpenShift version constraint won't be taken into account
// when returning the managed identities requirements.
func (r *ManagedIdentitiesRequirementsGetRequest) Version(value string) *ManagedIdentitiesRequirementsGetRequest {
	r.version = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *ManagedIdentitiesRequirementsGetRequest) Send() (result *ManagedIdentitiesRequirementsGetResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *ManagedIdentitiesRequirementsGetRequest) SendContext(ctx context.Context) (result *ManagedIdentitiesRequirementsGetResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.required != nil {
		helpers.AddValue(&query, "required", *r.required)
	}
	if r.version != nil {
		helpers.AddValue(&query, "version", *r.version)
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
	result = &ManagedIdentitiesRequirementsGetResponse{}
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
	err = readManagedIdentitiesRequirementsGetResponse(result, reader)
	if err != nil {
		return
	}
	return
}

// ManagedIdentitiesRequirementsGetResponse is the response for the 'get' method.
type ManagedIdentitiesRequirementsGetResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *ManagedIdentitiesRequirements
}

// Status returns the response status code.
func (r *ManagedIdentitiesRequirementsGetResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *ManagedIdentitiesRequirementsGetResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *ManagedIdentitiesRequirementsGetResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the value of the 'body' parameter.
//
// ManagedIdentitiesRequirements status response.
func (r *ManagedIdentitiesRequirementsGetResponse) Body() *ManagedIdentitiesRequirements {
	if r == nil {
		return nil
	}
	return r.body
}

// GetBody returns the value of the 'body' parameter and
// a flag indicating if the parameter has a value.
//
// ManagedIdentitiesRequirements status response.
func (r *ManagedIdentitiesRequirementsGetResponse) GetBody() (value *ManagedIdentitiesRequirements, ok bool) {
	ok = r != nil && r.body != nil
	if ok {
		value = r.body
	}
	return
}
