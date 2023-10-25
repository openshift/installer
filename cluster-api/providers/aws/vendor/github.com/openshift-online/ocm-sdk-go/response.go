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

// This file contains the implementation of the response type.

package sdk

import (
	"net/http"
)

// Response contains the information extracted from an HTTP POST response.
type Response struct {
	status int
	header http.Header
	body   []byte
}

// Status returns the response status code.
func (r *Response) Status() int {
	return r.status
}

// Bytes returns an slice of bytes containing the response body. Not that this will never return
// nil; if the response body is empty it will return an empty slice.
func (r *Response) Bytes() []byte {
	return r.body
}

// Bytes returns an string containing the response body.
func (r *Response) String() string {
	return string(r.body)
}

// Header returns the header value. In case there's no value for the header, an empty string ("") will be returned.
func (r *Response) Header(name string) string {
	if r.header == nil {
		return ""
	}
	return r.header.Get(name)
}
