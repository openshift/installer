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
// Package cloudresourcemanager contains support code for the CRM service.
package cloudresourcemanager

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *Folder) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"parent": dcl.ValueOrEmptyString(nr.Parent),
	}
	return dcl.URL("folders?parent={{parent}}", "https://cloudresourcemanager.googleapis.com/v2", userBasePath, params), nil

}

func (r *Folder) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("folders/{{name}}", "https://cloudresourcemanager.googleapis.com/v2", userBasePath, params), nil
}

func (r *Folder) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"parent": dcl.ValueOrEmptyString(nr.Parent),
	}
	return dcl.URL("folders?parent={{parent}}", "https://cloudresourcemanager.googleapis.com/v2", userBasePath, params), nil
}

func (r *Folder) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "MoveFolder" {
		fields := map[string]interface{}{
			"name": dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("folders/{{name}}:move", "https://cloudresourcemanager.googleapis.com/v2", userBasePath, fields), nil

	}
	if updateName == "UpdateFolder" {
		fields := map[string]interface{}{
			"name": dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("folders/{{name}}?updateMask=displayName", "https://cloudresourcemanager.googleapis.com/v2", userBasePath, fields), nil

	}
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

func (r *Folder) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("folders/{{name}}", "https://cloudresourcemanager.googleapis.com/v2", userBasePath, params), nil
}

func (op *updateFolderMoveFolderOperation) do(ctx context.Context, r *Folder, c *Client) error {
	_, err := c.GetFolder(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "MoveFolder")
	if err != nil {
		return err
	}

	req, err := newUpdateFolderMoveFolderRequest(ctx, r, c)
	if err != nil {
		return err
	}

	if p, ok := req["parent"]; ok {
		req["destinationParent"] = p
		delete(req, "parent")
	}

	c.Config.Logger.Infof("Created update: %#v", req)
	body, err := marshalUpdateFolderMoveFolderRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(ctx, c.Config, "https://cloudresourcemanager.googleapis.com/v1", "GET")

	if err != nil {
		return err
	}

	return nil
}

// The project is already effectively deleted if it's in DELETE_REQUESTED state.
func projectDeletePrecondition(r *Project) bool {
	return *r.LifecycleState == *ProjectLifecycleStateEnumRef("DELETE_REQUESTED")
}

// Project's list endpoint has a custom url method to use the filter query parameters.
func (r *Project) listURL(userBasePath string) (string, error) {
	parentParts := strings.Split(dcl.ValueOrEmptyString(r.Parent), "/")
	var parentType, parentID string
	if len(parentParts) == 2 {
		parentType = strings.TrimSuffix(parentParts[0], "s")
		parentID = parentParts[1]
		u, err := dcl.AddQueryParams("https://cloudresourcemanager.googleapis.com/v1/projects", map[string]string{
			"filter": fmt.Sprintf("parent.type=%s parent.id=%s", parentType, parentID),
		})
		if err != nil {
			return "", err
		}
		return u, nil
	}
	return "https://cloudresourcemanager.googleapis.com/v1/projects", nil
}

// expandProjectParent expands an instance of ProjectParent into a JSON
// request object.
func expandProjectParent(_ *Client, fval *string, _ *Project) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(fval) {
		return nil, nil
	}

	s := strings.Split(*fval, "/")
	m := make(map[string]interface{})
	if len(s) < 2 || dcl.IsEmptyValueIndirect(s[0]) || dcl.IsEmptyValueIndirect(s[1]) || !strings.HasSuffix(s[0], "s") {
		return m, fmt.Errorf("invalid parent argument. got value = %s. should be of the form organizations/org_id or folders/folder_id", *fval)
	}

	m["type"] = s[0][:len(s[0])-1]
	m["id"] = s[1]

	return m, nil
}

// flattenProjectParent flattens an instance of ProjectParent from a JSON
// response object.
func flattenProjectParent(c *Client, i interface{}, _ *Project) *string {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if dcl.IsEmptyValueIndirect(i) {
		return nil
	}
	// Ading s(plural) to change type to type(s). Example: organization/org_id to organizations/ord_id
	parent := fmt.Sprintf("%ss/%s", m["type"], m["id"])
	return &parent
}

func (r *TagKey) createURL(userBasePath string) (string, error) {
	params := make(map[string]interface{})
	return dcl.URL("tagKeys", "https://cloudresourcemanager.googleapis.com/v3", userBasePath, params), nil
}

func (r *TagKey) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("tagKeys/{{name}}", "https://cloudresourcemanager.googleapis.com/v3", userBasePath, params), nil
}

func (r *TagKey) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	fields := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("tagKeys/{{name}}?updateMask=displayName", "https://cloudresourcemanager.googleapis.com/v3", userBasePath, fields), nil
}

func (r *TagKey) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("tagKeys/{{name}}", "https://cloudresourcemanager.googleapis.com/v3", userBasePath, params), nil
}

func (r *TagValue) createURL(userBasePath string) (string, error) {
	params := make(map[string]any)
	return dcl.URL("tagValues", "https://cloudresourcemanager.googleapis.com/v3", userBasePath, params), nil
}

func (r *TagValue) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]any{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("tagValues/{{name}}", "https://cloudresourcemanager.googleapis.com/v3", userBasePath, params), nil
}

func (r *TagValue) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	fields := map[string]any{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("tagValues/{{name}}?updateMask=displayName", "https://cloudresourcemanager.googleapis.com/v3", userBasePath, fields), nil
}

func (r *TagValue) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]any{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("tagValues/{{name}}", "https://cloudresourcemanager.googleapis.com/v3", userBasePath, params), nil
}
