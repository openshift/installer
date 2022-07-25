/**
 * © Copyright IBM Corporation 2021. All Rights Reserved.
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
	"net/http"
	neturl "net/url"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
)

type BaseService struct {
	serviceUrlPathSegmentsSize int
	*core.BaseService
}

type validationRule struct {
	pathSegmentIndex   int
	errorParameterName string
	operationIds       []string
}

var docIdRule = validationRule{
	pathSegmentIndex:   1,
	errorParameterName: "Document ID",
	operationIds: []string{
		"DeleteDocument",
		"GetDocument",
		"HeadDocument",
		"PutDocument",
		"DeleteAttachment",
		"GetAttachment",
		"HeadAttachment",
		"PutAttachment",
	},
}
var attIdRule = validationRule{
	pathSegmentIndex:   2,
	errorParameterName: "Attachment name",
	operationIds: []string{
		"DeleteAttachment",
		"GetAttachment",
		"HeadAttachment",
		"PutAttachment",
	},
}

var validationRules = [...]*validationRule{&docIdRule, &attIdRule}
var rulesByOperation = make(map[string][]*validationRule)

func init() {
	for _, rule := range validationRules {
		// Build a map of operation name to a list of validations
		for _, operationId := range rule.operationIds {
			ruleExists := false
			rules, ok := rulesByOperation[operationId]
			if ok {
				// There are already some rules for this operationId
				// Check if the current rule already exists
				for _, existingRule := range rules {
					if existingRule == rule {
						ruleExists = true
						break
					}
				}
			}
			if !ruleExists {
				// The rule didn't exist, append it
				rulesByOperation[operationId] = append(rules, rule)
			}
		}
	}
}

func NewBaseService(opts *core.ServiceOptions) (*BaseService, error) {
	baseService, err := core.NewBaseService(opts)
	if err != nil {
		return &BaseService{}, err
	}
	return &BaseService{0, baseService}, nil
}

func (c *BaseService) Clone() *BaseService {
	baseService := c.BaseService.Clone()
	return &BaseService{c.serviceUrlPathSegmentsSize, baseService}
}

func (c *BaseService) Request(req *http.Request, result interface{}) (detailedResponse *core.DetailedResponse, err error) {
	// Extract the operation ID from the request headers.
	var operationId string
	//nolint
	header := req.Header["X-IBMCloud-SDK-Analytics"][0]
	for _, element := range strings.Split(header, ";") {
		if strings.HasPrefix(element, "operation_id") {
			operationId = strings.Split(element, "=")[1]
			break
		}
	}
	if operationId != "" {
		if rulesToApply, ok := rulesByOperation[operationId]; ok {
			requestUrlPathSegments := strings.Split(strings.Trim(req.URL.EscapedPath(), "/"), "/")
			// In the no-path case the result is a slice with an empty string
			// use a nil slice instead in those cases
			if len(requestUrlPathSegments) == 1 && requestUrlPathSegments[0] == "" {
				requestUrlPathSegments = []string{}
			}
			// Check each validation rule that applies to the operation.
			for _, rule := range rulesToApply {
				// Allow for any path segments that might exist in e.g. the URL of a proxy by offseting from the service URL index.
				pathSegmentIndex := c.serviceUrlPathSegmentsSize + rule.pathSegmentIndex
				if len(requestUrlPathSegments) > pathSegmentIndex {
					segmentToValidate := requestUrlPathSegments[pathSegmentIndex]
					if strings.HasPrefix(segmentToValidate, "_") {
						segmentToValidateMessage, err := neturl.PathUnescape(segmentToValidate)
						if err != nil {
							// If we couldn't unescape for some reason, just error with the escaped form
							segmentToValidateMessage = segmentToValidate
						}
						return nil, fmt.Errorf("%v %v starts with the invalid _ character.", rule.errorParameterName, segmentToValidateMessage)
					}
				}
			}
		}
	}
	return c.BaseService.Request(req, result)
}

func (c *BaseService) SetServiceURL(url string) error {
	err := c.BaseService.SetServiceURL(url)
	if err == nil {
		serviceUrl, err := neturl.ParseRequestURI(c.BaseService.GetServiceURL())
		if err == nil {
			serviceUrlPathSegments := strings.Split(strings.Trim(serviceUrl.EscapedPath(), "/"), "/")
			serviceUrlPathSegmentsSize := len(serviceUrlPathSegments)
			// In the no-path case the result is a slice with an empty string
			// set the size to zero in those cases
			if serviceUrlPathSegmentsSize == 1 && serviceUrlPathSegments[0] == "" {
				c.serviceUrlPathSegmentsSize = 0
			} else {
				c.serviceUrlPathSegmentsSize = serviceUrlPathSegmentsSize
			}
		}
	}
	return err
}
