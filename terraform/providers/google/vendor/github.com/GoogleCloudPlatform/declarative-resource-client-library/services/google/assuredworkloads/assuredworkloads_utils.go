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
// Package assuredworkloads contains support code for the assuredworkload service.
package assuredworkloads

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// Returns the URL of the project resource with the given index in the workload.
func (r *Workload) projectURL(userBasePath string, index int) (string, error) {
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(r.Resources[index].ResourceId),
	}
	return dcl.URL("projects/{{project}}", "https://cloudresourcemanager.googleapis.com/v1/", userBasePath, params), nil
}

// Returns the URL of the folder resource with the given index in the workload.
func (r *Workload) folderURL(userBasePath string, index int) (string, error) {
	params := map[string]interface{}{
		"folder": dcl.ValueOrEmptyString(r.Resources[index].ResourceId),
	}
	return dcl.URL("folders/{{folder}}", "https://cloudresourcemanager.googleapis.com/v2/", userBasePath, params), nil
}

// Returns the lifecycle state of the project or folder resource with the given url.
func lifecycleState(ctx context.Context, client *Client, url string) (string, error) {
	resp, err := dcl.SendRequest(ctx, client.Config, "GET", url, &bytes.Buffer{}, client.Config.RetryProvider)
	if err != nil {
		return "", err
	}
	defer resp.Response.Body.Close()
	b, err := io.ReadAll(resp.Response.Body)
	if err != nil {
		return "", err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	state, ok := m["lifecycleState"].(string)
	if !ok {
		return "", fmt.Errorf("no lifecycle state for resource at %q", url)
	}
	return state, nil
}

// Deletes the resource with the given URL. Returns true if it is already in DELETE_REQUESTED state,
// otherwise returns false.
func deleteResource(ctx context.Context, client *Client, url string) (bool, error) {
	state, err := lifecycleState(ctx, client, url)
	if err != nil {
		return false, err
	}
	if state == "DELETE_REQUESTED" {
		// Do not delete an already deleted resource.
		return true, nil
	}
	// Send delete request for resources not already deleted.
	_, err = dcl.SendRequest(ctx, client.Config, "DELETE", url, &bytes.Buffer{}, client.Config.RetryProvider)
	if err != nil {
		return false, fmt.Errorf("failed to delete resource at %s: %w", url, err)
	}
	return false, nil
}

// Deletes projects and folders owned by the workload prior to workload deletion.
func (r *Workload) deleteResources(ctx context.Context, client *Client) error {
	nr := r.urlNormalized()
	return dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		// First, delete projects
		for i, resource := range nr.Resources {
			if resource.ResourceType == nil {
				return nil, fmt.Errorf("nil resource type in workload %q", dcl.ValueOrEmptyString(nr.Name))
			}
			if *resource.ResourceType == WorkloadResourcesResourceTypeEnum("CONSUMER_PROJECT") || *resource.ResourceType == WorkloadResourcesResourceTypeEnum("ENCRYPTION_KEYS_PROJECT") {
				u, err := nr.projectURL(client.Config.BasePath, i)
				if err != nil {
					return nil, err
				}
				deleted, err := deleteResource(ctx, client, u)
				if err != nil {
					return nil, err
				}
				if !deleted {
					// Retry until all resources are being deleted.
					return &dcl.RetryDetails{}, dcl.OperationNotDone{}
				}
			}
		}
		// Then, delete folders
		for i, resource := range nr.Resources {
			if *resource.ResourceType == WorkloadResourcesResourceTypeEnum("CONSUMER_FOLDER") {
				u, err := nr.folderURL(client.Config.BasePath, i)
				if err != nil {
					return nil, err
				}
				deleted, err := deleteResource(ctx, client, u)
				if err != nil {
					return nil, err
				}
				if !deleted {
					// Retry until all resources are being deleted.
					return &dcl.RetryDetails{}, dcl.OperationNotDone{}
				}
			}
		}
		// All project and folder resources are in DELETE_REQUESTED state.
		return nil, nil
	}, client.Config.RetryProvider)
}
