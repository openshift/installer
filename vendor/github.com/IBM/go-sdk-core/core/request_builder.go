package core

// (C) Copyright IBM Corp. 2019.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// common HTTP methods
const (
	POST   = http.MethodPost
	GET    = http.MethodGet
	DELETE = http.MethodDelete
	PUT    = http.MethodPut
	PATCH  = http.MethodPatch
)

// common headers
const (
	Accept                  = "Accept"
	APPLICATION_JSON        = "application/json"
	CONTENT_DISPOSITION     = "Content-Disposition"
	CONTENT_TYPE            = "Content-Type"
	FORM_URL_ENCODED_HEADER = "application/x-www-form-urlencoded"

	ERRORMSG_SERVICE_URL_MISSING = "The service URL is required."
	ERRORMSG_SERVICE_URL_INVALID = "There was an error parsing the service URL: %s"
)

// A FormData stores information for form data
type FormData struct {
	fileName    string
	contentType string
	contents    interface{}
}

// A RequestBuilder is an HTTP request to be sent to the service
type RequestBuilder struct {
	Method string
	URL    *url.URL
	Header http.Header
	Body   io.Reader
	Query  map[string][]string
	Form   map[string][]FormData
}

// NewRequestBuilder : Initiates a new request
func NewRequestBuilder(method string) *RequestBuilder {
	return &RequestBuilder{
		Method: method,
		Header: make(http.Header),
		Query:  make(map[string][]string),
		Form:   make(map[string][]FormData),
	}
}

// ConstructHTTPURL creates a properly-encoded URL with path parameters.
// This function returns an error if the serviceURL is "" or is an
// invalid URL string (e.g. ":<badscheme>").
func (requestBuilder *RequestBuilder) ConstructHTTPURL(serviceURL string, pathSegments []string, pathParameters []string) (*RequestBuilder, error) {
	if serviceURL == "" {
		return requestBuilder, fmt.Errorf(ERRORMSG_SERVICE_URL_MISSING)
	}
	var URL *url.URL

	URL, err := url.Parse(serviceURL)
	if err != nil {
		return requestBuilder, fmt.Errorf(ERRORMSG_SERVICE_URL_INVALID, err.Error())
	}

	for i, pathSegment := range pathSegments {
		if pathSegment != "" {
			URL.Path += "/" + pathSegment
		}

		if pathParameters != nil && i < len(pathParameters) {
			URL.Path += "/" + pathParameters[i]
		}
	}
	requestBuilder.URL = URL
	return requestBuilder, nil
}

// AddQuery adds Query name and value
func (requestBuilder *RequestBuilder) AddQuery(name string, value string) *RequestBuilder {
	requestBuilder.Query[name] = append(requestBuilder.Query[name], value)
	return requestBuilder
}

// AddHeader adds header name and value
func (requestBuilder *RequestBuilder) AddHeader(name string, value string) *RequestBuilder {
	requestBuilder.Header[name] = []string{value}
	return requestBuilder
}

// AddFormData makes an entry for Form data
func (requestBuilder *RequestBuilder) AddFormData(fieldName string, fileName string, contentType string,
	contents interface{}) *RequestBuilder {
	if fileName == "" {
		if file, ok := contents.(*os.File); ok {
			if !((os.File{}) == *file) { // if file is not empty
				name := filepath.Base(file.Name())
				fileName = name
			}
		}
	}
	requestBuilder.Form[fieldName] = append(requestBuilder.Form[fieldName], FormData{
		fileName:    fileName,
		contentType: contentType,
		contents:    contents,
	})
	return requestBuilder
}

// SetBodyContentJSON - set the body content from a JSON structure
func (requestBuilder *RequestBuilder) SetBodyContentJSON(bodyContent interface{}) (*RequestBuilder, error) {
	requestBuilder.Body = new(bytes.Buffer)
	err := json.NewEncoder(requestBuilder.Body.(io.Writer)).Encode(bodyContent)
	return requestBuilder, err
}

// SetBodyContentString - set the body content from a string
func (requestBuilder *RequestBuilder) SetBodyContentString(bodyContent string) (*RequestBuilder, error) {
	requestBuilder.Body = strings.NewReader(bodyContent)
	return requestBuilder, nil
}

// SetBodyContentStream - set the body content from an io.Reader instance
func (requestBuilder *RequestBuilder) SetBodyContentStream(bodyContent io.Reader) (*RequestBuilder, error) {
	requestBuilder.Body = bodyContent
	return requestBuilder, nil
}

// CreateMultipartWriter initializes a new multipart writer
func (requestBuilder *RequestBuilder) createMultipartWriter() *multipart.Writer {
	buff := new(bytes.Buffer)
	requestBuilder.Body = buff
	return multipart.NewWriter(buff)
}

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name and file name and contentType
func createFormFile(formWriter *multipart.Writer, fieldname string, filename string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	contentDisposition := fmt.Sprintf(`form-data; name="%s"`, fieldname)
	if filename != "" {
		contentDisposition += fmt.Sprintf(`; filename="%s"`, filename)
	}

	h.Set(CONTENT_DISPOSITION, contentDisposition)
	if contentType != "" {
		h.Set(CONTENT_TYPE, contentType)
	}

	return formWriter.CreatePart(h)
}

// SetBodyContentForMultipart - sets the body content for a part in a multi-part form
func (requestBuilder *RequestBuilder) SetBodyContentForMultipart(contentType string, content interface{}, writer io.Writer) error {
	var err error
	if stream, ok := content.(io.Reader); ok {
		_, err = io.Copy(writer, stream)
	} else if stream, ok := content.(*io.ReadCloser); ok {
		_, err = io.Copy(writer, *stream)
	} else if IsJSONMimeType(contentType) || IsJSONPatchMimeType(contentType) {
		err = json.NewEncoder(writer).Encode(content)
	} else if str, ok := content.(string); ok {
		_, err = writer.Write([]byte(str))
	} else if strPtr, ok := content.(*string); ok {
		_, err = writer.Write([]byte(*strPtr))
	} else {
		err = fmt.Errorf("Error: unable to determine the type of 'content' provided")
	}
	return err
}

// Build the request
func (requestBuilder *RequestBuilder) Build() (*http.Request, error) {
	// Create multipart form data
	if len(requestBuilder.Form) > 0 {
		// handle both application/x-www-form-urlencoded or multipart/form-data
		contentType := requestBuilder.Header.Get(CONTENT_TYPE)
		if contentType == FORM_URL_ENCODED_HEADER {
			data := url.Values{}
			for fieldName, l := range requestBuilder.Form {
				for _, v := range l {
					data.Add(fieldName, v.contents.(string))
				}
			}
			_, err := requestBuilder.SetBodyContentString(data.Encode())
			if err != nil {
				return nil, err
			}
		} else {
			formWriter := requestBuilder.createMultipartWriter()
			for fieldName, l := range requestBuilder.Form {
				for _, v := range l {
					dataPartWriter, err := createFormFile(formWriter, fieldName, v.fileName, v.contentType)
					if err != nil {
						return nil, err
					}
					if err = requestBuilder.SetBodyContentForMultipart(v.contentType,
						v.contents, dataPartWriter); err != nil {
						return nil, err
					}
				}
			}

			requestBuilder.AddHeader("Content-Type", formWriter.FormDataContentType())
			err := formWriter.Close()
			if err != nil {
				return nil, err
			}
		}
	}

	// Create the request
	req, err := http.NewRequest(requestBuilder.Method, requestBuilder.URL.String(), requestBuilder.Body)
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header = requestBuilder.Header

	// Query
	query := req.URL.Query()
	for k, l := range requestBuilder.Query {
		for _, v := range l {
			query.Add(k, v)
		}
	}
	// Encode query
	req.URL.RawQuery = query.Encode()

	return req, nil
}

// SetBodyContent - sets the body content from one of three different sources
func (requestBuilder *RequestBuilder) SetBodyContent(contentType string, jsonContent interface{}, jsonPatchContent interface{},
	nonJSONContent interface{}) (builder *RequestBuilder, err error) {
	if jsonContent != nil {
		builder, err = requestBuilder.SetBodyContentJSON(jsonContent)
		if err != nil {
			return
		}
	} else if jsonPatchContent != nil {
		builder, err = requestBuilder.SetBodyContentJSON(jsonPatchContent)
		if err != nil {
			return
		}
	} else {
		// Set the non-JSON body content based on the type of value passed in,
		// which should be a "string", "*string" or an "io.Reader"
		if str, ok := nonJSONContent.(string); ok {
			builder, err = requestBuilder.SetBodyContentString(str)
		} else if strPtr, ok := nonJSONContent.(*string); ok {
			builder, err = requestBuilder.SetBodyContentString(*strPtr)
		} else if stream, ok := nonJSONContent.(io.Reader); ok {
			builder, err = requestBuilder.SetBodyContentStream(stream)
		} else if stream, ok := nonJSONContent.(*io.ReadCloser); ok {
			builder, err = requestBuilder.SetBodyContentStream(*stream)
		} else {
			builder = requestBuilder
			err = fmt.Errorf("Invalid type for non-JSON body content: %s", reflect.TypeOf(nonJSONContent).String())
		}
	}
	return
}
