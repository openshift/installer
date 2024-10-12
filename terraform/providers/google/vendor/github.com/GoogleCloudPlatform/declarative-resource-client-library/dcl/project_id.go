// Copyright 2024 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// This matches either the entire string if it contains no forward slashes or just projects/{project_number}/ if it does.
var projectNumberRegex = regexp.MustCompile(`(^\d+$|projects/\d+|metricsScopes/\d+)`)

// This matches either the entire string if it contains no forward slashes or just projects/{project_id}/ if it does.
var projectIDRegex = regexp.MustCompile(`(^[^/]+$|projects/[^/]+|metricsScopes/[^/]+)`)

// ProjectResponse is the response from Cloud Resource Manager.
type ProjectResponse struct {
	ProjectID     string `json:"projectId"`
	ProjectNumber string `json:"projectNumber"`
}

// FlattenProjectNumbersToIDs converts a project number to project ID.
func FlattenProjectNumbersToIDs(config *Config, fromServer *string) *string {
	if fromServer == nil {
		return nil
	}
	// Look for a number somewhere in here.
	editedServer := projectNumberRegex.ReplaceAllStringFunc(*fromServer, func(number string) string {
		config.Logger.Infof("Preparing to use Cloud Resource Manager to convert %s to project id", number)

		p, err := fetchProjectInfo(config, number)
		if err != nil {
			config.Logger.Warning(err)
			return number
		}

		if strings.HasPrefix(number, "projects/") {
			p.ProjectID = "projects/" + p.ProjectID
		}
		if strings.HasPrefix(number, "metricsScopes/") {
			p.ProjectID = "metricsScopes/" + p.ProjectID
		}

		return p.ProjectID
	})
	return &editedServer
}

var fetchProjectInfo = FetchProjectInfo

// ExpandProjectIDsToNumbers converts a project ID to a project number.
func ExpandProjectIDsToNumbers(config *Config, fromConfig *string) (*string, error) {
	if fromConfig == nil {
		return nil, nil
	}

	// Look for a project id somewhere in here.
	editedConfig := projectIDRegex.ReplaceAllStringFunc(*fromConfig, func(id string) string {
		config.Logger.Infof("Preparing to convert %s to project number", id)

		p, err := fetchProjectInfo(config, id)
		if err != nil {
			config.Logger.Warning(err)
			return id
		}

		if strings.HasPrefix(id, "projects/") {
			p.ProjectNumber = "projects/" + p.ProjectNumber
		}
		if strings.HasPrefix(id, "metricsScopes/") {
			p.ProjectNumber = "metricsScopes/" + p.ProjectNumber
		}

		return p.ProjectNumber
	})
	return &editedConfig, nil
}

// FetchProjectInfo returns a ProjectResponse from CloudResourceManager.
func FetchProjectInfo(config *Config, projectIdentifier string) (ProjectResponse, error) {
	var p ProjectResponse
	trimmedIdentifier := strings.TrimPrefix(projectIdentifier, "projects/")
	trimmedIdentifier = strings.TrimPrefix(trimmedIdentifier, "metricsScopes/")
	trimmedIdentifier = strings.TrimSuffix(trimmedIdentifier, "/")
	retryDetails, err := SendRequest(context.TODO(), config, "GET", "https://cloudresourcemanager.googleapis.com/v1/projects/"+trimmedIdentifier, nil, nil)
	if err != nil {
		return p, fmt.Errorf("failed to send request for project info using identifier %q: %s", projectIdentifier, err)
	}
	if err := ParseResponse(retryDetails.Response, &p); err != nil {
		return p, fmt.Errorf("failed to parse response %v for project with identifier %q: %s", retryDetails.Response, projectIdentifier, err)
	}

	return p, nil
}
