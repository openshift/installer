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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MetadataRequest is the request to retrieve the metadata.
type MetadataRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
}

// MetadataResponse is the response for the metadata request.
type MetadataResponse struct {
	status int
	header http.Header
	err    *errors.Error
	body   *Metadata
}

// Parameter adds a query parameter.
func (r *MetadataRequest) Parameter(name string, value interface{}) *MetadataRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *MetadataRequest) Header(name string, value interface{}) *MetadataRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Send sends the metadata request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *MetadataRequest) Send() (result *MetadataResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends the metadata request, waits for the response, and returns it.
func (r *MetadataRequest) SendContext(ctx context.Context) (result *MetadataResponse, err error) {
	query := helpers.CopyQuery(r.query)
	header := helpers.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: http.MethodGet,
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
	result = &MetadataResponse{
		status: response.StatusCode,
		header: response.Header,
	}
	reader := bufio.NewReader(response.Body)
	_, err = reader.Peek(1)
	if err == io.EOF {
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
	result.body, err = UnmarshalMetadata(reader)
	if err != nil {
		return
	}
	return
}

// Status returns the response status code.
func (r *MetadataResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *MetadataResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *MetadataResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Body returns the response body.
func (r *MetadataResponse) Body() *Metadata {
	return r.body
}
