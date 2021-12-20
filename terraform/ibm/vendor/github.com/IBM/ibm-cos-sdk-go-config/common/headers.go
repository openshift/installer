/**
 * (C) Copyright IBM Corp. 2019.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"fmt"
	"runtime"
)

const (
	// HeaderSDKAnalytics Header
	HeaderSDKAnalytics = "X-IBMCloud-SDK-Analytics"
	// HeaderUserAgent Header
	HeaderUserAgent = "User-Agent"

	// SDKName name of this SDK
	SDKName = "ibm-cos-resource-config-sdk-go"
)

// GetSdkHeaders - returns the set of SDK-specific headers to be included in an outgoing request.
// Parameters:
//   serviceName - the name of the service as defined in the API definition (e.g. "MyService1")
//   serviceVersion - the version of the service as defined in the API definition (e.g. "V1")
//   operationId - the operationId as defined in the API definition (e.g. getContext)
//
// Returns:
//   a Map which contains the set of headers to be included in the REST API request
//
func GetSdkHeaders(serviceName string, serviceVersion string, operationID string) map[string]string {
	sdkHeaders := make(map[string]string)

	sdkHeaders[HeaderSDKAnalytics] = fmt.Sprintf("service_name=%s;service_version=%s;operation_id=%s",
		serviceName, serviceVersion, operationID)

	sdkHeaders[HeaderUserAgent] = GetUserAgentInfo()

	return sdkHeaders
}

var userAgent = fmt.Sprintf("%s-%s %s", SDKName, SDKVersion, GetSystemInfo())

// GetUserAgentInfo returns user agent
func GetUserAgentInfo() string {
	return userAgent
}

var systemInfo = fmt.Sprintf("(arch=%s; os=%s; go.version=%s)", runtime.GOARCH, runtime.GOOS, runtime.Version())

// GetSystemInfo returns system information
func GetSystemInfo() string {
	return systemInfo
}
