/*
Copyright (c) 2018 Red Hat, Inc.

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

// This file contains the implementation of the request type.

package sdk

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/openshift-online/ocm-sdk-go/internal"
)

// Request contains the information and logic needed to perform an HTTP request.
type Request struct {
	transport http.RoundTripper
	method    string
	path      string
	query     url.Values
	header    http.Header
	body      []byte
}

// GetMethod returns the request method (GET/POST/PATCH/PUT/DELETE).
func (r *Request) GetMethod() string {
	return r.method
}

// Path defines the request path, for example `/api/clusters_mgmt/v1/clusters`. This is mandatory; an
// error will be returned immediately when calling the Send method if this isn't provided.
func (r *Request) Path(value string) *Request {
	r.path = value
	return r
}

// GetPath returns the request path.
func (r *Request) GetPath() string {
	return r.path
}

// Parameter adds a query parameter.
func (r *Request) Parameter(name string, value interface{}) *Request {
	internal.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *Request) Header(name string, value interface{}) *Request {
	internal.AddHeader(&r.header, name, value)
	return r
}

// Bytes sets the request body from an slice of bytes.
func (r *Request) Bytes(value []byte) *Request {
	if value != nil {
		r.body = make([]byte, len(value))
		copy(r.body, value)
	} else {
		r.body = nil
	}
	return r
}

// String sets the request body from an string.
func (r *Request) String(value string) *Request {
	r.body = []byte(value)
	return r
}

// Send sends this request to the server and returns the corresponding response, or an error if
// something fails. Note that any HTTP status code returned by the server is considered a valid
// response, and will not be translated into an error. It is up to the caller to check the status
// code and handle it.
//
// This operation is potentially lengthy, as it requires network communication. Consider using a
// context and the SendContext method.
func (r *Request) Send() (result *Response, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request to the server and returns the corresponding response, or an error
// if something fails. Note that any HTTP status code returned by the server is considered a valid
// response, and will not be translated into an error. It is up to the caller to check the status
// code and handle it.
func (r *Request) SendContext(ctx context.Context) (result *Response, err error) {
	query := internal.CopyQuery(r.query)
	header := internal.CopyHeader(r.header)
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	var body io.ReadCloser
	if r.body != nil {
		body = io.NopCloser(bytes.NewBuffer(r.body))
	}
	request := &http.Request{
		Method: r.method,
		URL:    uri,
		Header: header,
		Body:   body,
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = new(Response)
	result.status = response.StatusCode
	result.header = response.Header
	result.body, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}
