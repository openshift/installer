/**
 * © Copyright IBM Corporation 2020, 2022. All Rights Reserved.
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
)

const (
	headerNameSdkAnalytics = "X-IBMCloud-SDK-Analytics"
)

//
// GetSdkHeaders - returns the set of SDK-specific headers to be included in an outgoing request.
//
// This function is invoked by generated service methods (i.e. methods which implement the REST API operations
// defined within the API definition). The purpose of this function is to give the SDK implementor the opportunity
// to provide SDK-specific HTTP headers that will be sent with an outgoing REST API request.
// This function is invoked for each invocation of a generated service method,
// so the set of HTTP headers could be request-specific.
//
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

	sdkHeaders[headerNameSdkAnalytics] = GetSdkAnalyticsHeader(serviceName, serviceVersion, operationID)

	return sdkHeaders
}

func GetSdkAnalyticsHeader(serviceName string, serviceVersion string, operationID string) string {
	return fmt.Sprintf("service_name=%s;service_version=%s;operation_id=%s",
		serviceName, serviceVersion, operationID)
}
