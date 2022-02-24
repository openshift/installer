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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	header_name_USER_AGENT = "User-Agent"
	sdk_name               = "ibm-go-sdk-core"
)

// ServiceOptions is a struct of configuration values for a service.
type ServiceOptions struct {
	// This is the base URL associated with the service instance. This value will
	// be combined with the paths for each operation to form the request URL
	// [required].
	URL string

	// Authenticator holds the authenticator implementation to be used by the
	// service instance to authenticate outbound requests, typically by adding the
	// HTTP "Authorization" header.
	Authenticator Authenticator
}

// BaseService implements the common functionality shared by generated services
// to manage requests and responses, authenticate outbound requests, etc.
type BaseService struct {

	// Configuration values for a service.
	Options *ServiceOptions

	// A set of "default" http headers to be included with each outbound request.
	DefaultHeaders http.Header

	// The HTTP Client used to send requests and receive responses.
	Client *http.Client

	// The value to be used for the "User-Agent" HTTP header that is added to each
	// outbound request. If this value is not set, then a default value will be
	// used for the header.
	UserAgent string
}

// NewBaseService constructs a new instance of BaseService. Validation on input
// parameters and service options will be performed before instance creation.
func NewBaseService(options *ServiceOptions) (*BaseService, error) {
	if HasBadFirstOrLastChar(options.URL) {
		return nil, fmt.Errorf(ERRORMSG_PROP_INVALID, "URL")
	}

	if options.Authenticator == nil {
		return nil, fmt.Errorf(ERRORMSG_NO_AUTHENTICATOR)
	}

	if err := options.Authenticator.Validate(); err != nil {
		return nil, err
	}

	service := BaseService{
		Options: options,

		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	// Set a default value for the User-Agent http header.
	service.SetUserAgent(service.buildUserAgent())

	return &service, nil
}

// ConfigureService updates the service with external configuration values.
func (service *BaseService) ConfigureService(serviceName string) error {
	// Try to load service properties from external config.
	serviceProps, err := getServiceProperties(serviceName)
	if err != nil {
		return err
	}

	// If we were able to load any properties for this service, then check to see if the
	// service-level properties were present and set them on the service if so.
	if serviceProps != nil {

		// URL
		if url, ok := serviceProps[PROPNAME_SVC_URL]; ok && url != "" {
			err := service.SetURL(url)
			if err != nil {
				return err
			}
		}

		// DISABLE_SSL
		if disableSSL, ok := serviceProps[PROPNAME_SVC_DISABLE_SSL]; ok && disableSSL != "" {
			// Convert the config string to bool.
			boolValue, err := strconv.ParseBool(disableSSL)
			if err != nil {
				boolValue = false
			}

			// If requested, disable SSL.
			if boolValue {
				service.DisableSSLVerification()
			}
		}
	}
	return nil
}

// SetURL sets the service URL.
//
// Deprecated: use SetServiceURL instead.
func (service *BaseService) SetURL(url string) error {
	return service.SetServiceURL(url)
}

// SetServiceURL sets the service URL.
func (service *BaseService) SetServiceURL(url string) error {
	if HasBadFirstOrLastChar(url) {
		return fmt.Errorf(ERRORMSG_PROP_INVALID, "URL")
	}

	service.Options.URL = url
	return nil
}

// GetServiceURL returns the service URL.
func (service *BaseService) GetServiceURL() string {
	return service.Options.URL
}

// SetDefaultHeaders sets HTTP headers to be sent in every request.
func (service *BaseService) SetDefaultHeaders(headers http.Header) {
	service.DefaultHeaders = headers
}

// SetHTTPClient updates the client handling the requests.
func (service *BaseService) SetHTTPClient(client *http.Client) {
	service.Client = client
}

// DisableSSLVerification skips SSL verification.
func (service *BaseService) DisableSSLVerification() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	service.Client.Transport = tr
}

// buildUserAgent builds the user agent string.
func (service *BaseService) buildUserAgent() string {
	return fmt.Sprintf("%s-%s %s", sdk_name, __VERSION__, SystemInfo())
}

// SetUserAgent sets the user agent value.
func (service *BaseService) SetUserAgent(userAgentString string) {
	if userAgentString == "" {
		service.UserAgent = service.buildUserAgent()
	}
	service.UserAgent = userAgentString
}

// Request invokes the specified HTTP request and returns the response.
func (service *BaseService) Request(req *http.Request, result interface{}) (detailedResponse *DetailedResponse, err error) {
	// Add default headers.
	if service.DefaultHeaders != nil {
		for k, v := range service.DefaultHeaders {
			req.Header.Add(k, strings.Join(v, ""))
		}
	}

	// Add the default User-Agent header if not already present.
	userAgent := req.Header.Get(header_name_USER_AGENT)
	if userAgent == "" {
		req.Header.Add(header_name_USER_AGENT, service.UserAgent)
	}

	// Add authentication to the outbound request.
	if service.Options.Authenticator == nil {
		err = fmt.Errorf(ERRORMSG_NO_AUTHENTICATOR)
		return
	}

	err = service.Options.Authenticator.Authenticate(req)
	if err != nil {
		return
	}

	// Invoke the request.
	httpResponse, err := service.Client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), SSL_CERTIFICATION_ERROR) {
			err = fmt.Errorf(ERRORMSG_SSL_VERIFICATION_FAILED + "\n" + err.Error())
		}
		return
	}

	// Start to populate the DetailedResponse.
	detailedResponse = &DetailedResponse{
		StatusCode: httpResponse.StatusCode,
		Headers:    httpResponse.Header,
	}

	contentType := httpResponse.Header.Get(CONTENT_TYPE)

	// If the operation was unsuccessful, then set up the DetailedResponse
	// and error objects appropriately.
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {

		var responseBody []byte

		// First, read the response body into a byte array.
		if httpResponse.Body != nil {
			var readErr error

			defer httpResponse.Body.Close()
			responseBody, readErr = ioutil.ReadAll(httpResponse.Body)
			if readErr != nil {
				err = fmt.Errorf("An error occurred while reading the response body: '%s'", readErr.Error())
				return
			}
		}

		// If the responseBody is empty, then just return a generic error based on the status code.
		if len(responseBody) == 0 {
			err = fmt.Errorf(http.StatusText(httpResponse.StatusCode))
			return
		}

		// For a JSON-based error response body, decode it into a map (generic JSON object).
		if IsJSONMimeType(contentType) {
			// Return the error response body as a map, along with an
			// error object containing our best guess at an error message.
			responseMap, decodeErr := decodeAsMap(responseBody)
			if decodeErr == nil {
				detailedResponse.Result = responseMap
				err = fmt.Errorf(getErrorMessage(responseMap, detailedResponse.StatusCode))
				return
			}
		}

		// For a non-JSON response or if we tripped while decoding the JSON response,
		// just return the response body byte array in the RawResult field along with
		// an error object that contains the generic error message for the status code.
		detailedResponse.RawResult = responseBody
		err = fmt.Errorf(http.StatusText(httpResponse.StatusCode))
		return
	}

	// Operation was successful and we are expecting a response, so process the response.
	if result != nil {
		// For a JSON response, decode it into the response object.
		resultType := reflect.TypeOf(result).String()
		if IsJSONMimeType(contentType) && resultType != "*io.ReadCloser" {

			// First, read the response body into a byte array.
			defer httpResponse.Body.Close()
			responseBody, readErr := ioutil.ReadAll(httpResponse.Body)
			if readErr != nil {
				err = fmt.Errorf("An error occurred while reading the response body: '%s'", readErr.Error())
				return
			}

			// Decode the byte array as JSON.
			decodeErr := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&result)
			if decodeErr != nil {
				// Error decoding the response body.
				// Return the response body in RawResult, along with an error.
				err = fmt.Errorf("An error occurred while unmarshalling the response body: '%s'", decodeErr.Error())
				detailedResponse.RawResult = responseBody
				return
			}

			// Decode step was successful. Return the decoded response object in the Result field.
			detailedResponse.Result = result
			return
		}

		// For a non-JSON response body, just return it as an io.Reader in the Result field.
		detailedResponse.Result = httpResponse.Body
	}

	return
}

// Errors is a struct used to hold an array of errors received in an operation
// response.
type Errors struct {
	Errors []Error `json:"errors,omitempty"`
}

// Error is a struct used to represent a single error received in an operation
// response.
type Error struct {
	Message string `json:"message,omitempty"`
}

// decodeAsMap: Decode the specified JSON byte-stream into a map (akin to a generic JSON object).
// Notes:
// 1) This function will return the map (result of decoding the byte-stream) as well as the raw
// byte buffer.  We return the byte buffer in addition to the decoded map so that the caller can
// re-use (if necessary) the stream of bytes after we've consumed them via the JSON decode step.
// 2) The primary return value of this function will be:
//    a) an instance of map[string]interface{} if the specified byte-stream was successfully
//       decoded as JSON.
//    b) the string form of the byte-stream if the byte-stream could not be successfully
//       decoded as JSON.
// 3) This function will close the io.ReadCloser before returning.
func decodeAsMap(byteBuffer []byte) (result map[string]interface{}, err error) {
	err = json.NewDecoder(bytes.NewReader(byteBuffer)).Decode(&result)
	return
}

// getErrorMessage: try to retrieve an error message from the decoded response body (map).
func getErrorMessage(responseMap map[string]interface{}, statusCode int) string {

	// If the response contained the "errors" field, then try to deserialize responseMap
	// into an array of Error structs, then return the first entry's "Message" field.
	if _, ok := responseMap["errors"]; ok {
		var errors Errors
		responseBuffer, _ := json.Marshal(responseMap)
		if err := json.Unmarshal(responseBuffer, &errors); err == nil {
			return errors.Errors[0].Message
		}
	}

	// Return the "error" field if present and is a string.
	if val, ok := responseMap["error"]; ok {
		errorMsg, ok := val.(string)
		if ok {
			return errorMsg
		}
	}

	// Return the "message" field if present and is a string.
	if val, ok := responseMap["message"]; ok {
		errorMsg, ok := val.(string)
		if ok {
			return errorMsg
		}
	}

	// Finally, return the "errorMessage" field if present and is a string.
	if val, ok := responseMap["errorMessage"]; ok {
		errorMsg, ok := val.(string)
		if ok {
			return errorMsg
		}
	}

	// If we couldn't find an error message above, just return the generic text
	// for the status code.
	return http.StatusText(statusCode)
}
