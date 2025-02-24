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
package cloudresourcemanager

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Folder struct {
	Name        *string          `json:"name"`
	Parent      *string          `json:"parent"`
	DisplayName *string          `json:"displayName"`
	State       *FolderStateEnum `json:"state"`
	CreateTime  *string          `json:"createTime"`
	UpdateTime  *string          `json:"updateTime"`
	DeleteTime  *string          `json:"deleteTime"`
	Etag        *string          `json:"etag"`
}

func (r *Folder) String() string {
	return dcl.SprintResource(r)
}

// The enum FolderStateEnum.
type FolderStateEnum string

// FolderStateEnumRef returns a *FolderStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func FolderStateEnumRef(s string) *FolderStateEnum {
	v := FolderStateEnum(s)
	return &v
}

func (v FolderStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"LIFECYCLE_STATE_UNSPECIFIED", "ACTIVE", "DELETE_REQUESTED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "FolderStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Folder) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "cloud_resource_manager",
		Type:    "Folder",
		Version: "cloudresourcemanager",
	}
}

func (r *Folder) ID() (string, error) {
	if err := extractFolderFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":         dcl.ValueOrEmptyString(nr.Name),
		"parent":       dcl.ValueOrEmptyString(nr.Parent),
		"display_name": dcl.ValueOrEmptyString(nr.DisplayName),
		"state":        dcl.ValueOrEmptyString(nr.State),
		"create_time":  dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":  dcl.ValueOrEmptyString(nr.UpdateTime),
		"delete_time":  dcl.ValueOrEmptyString(nr.DeleteTime),
		"etag":         dcl.ValueOrEmptyString(nr.Etag),
	}
	return dcl.Nprintf("folders/{{name}}", params), nil
}

const FolderMaxPage = -1

type FolderList struct {
	Items []*Folder

	nextToken string

	pageSize int32

	resource *Folder
}

func (l *FolderList) HasNext() bool {
	return l.nextToken != ""
}

func (l *FolderList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listFolder(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListFolder(ctx context.Context, parent string) (*FolderList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListFolderWithMaxResults(ctx, parent, FolderMaxPage)

}

func (c *Client) ListFolderWithMaxResults(ctx context.Context, parent string, pageSize int32) (*FolderList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Folder{
		Parent: &parent,
	}
	items, token, err := c.listFolder(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &FolderList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetFolder(ctx context.Context, r *Folder) (*Folder, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractFolderFields(r)

	b, err := c.getFolderRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalFolder(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeFolderNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractFolderFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteFolder(ctx context.Context, r *Folder) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Folder resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Folder...")
	deleteOp := deleteFolderOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllFolder deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllFolder(ctx context.Context, parent string, filter func(*Folder) bool) error {
	listObj, err := c.ListFolder(ctx, parent)
	if err != nil {
		return err
	}

	err = c.deleteAllFolder(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllFolder(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyFolder(ctx context.Context, rawDesired *Folder, opts ...dcl.ApplyOption) (*Folder, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Folder
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyFolderHelper(c, ctx, rawDesired, opts...)
		resultNewState = newState
		if err != nil {
			// If the error is 409, there is conflict in resource update.
			// Here we want to apply changes based on latest state.
			if dcl.IsConflictError(err) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		return nil, nil
	}, c.Config.RetryProvider)
	return resultNewState, err
}

func applyFolderHelper(c *Client, ctx context.Context, rawDesired *Folder, opts ...dcl.ApplyOption) (*Folder, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyFolder...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractFolderFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.folderDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToFolderDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	var create bool
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		if dcl.HasLifecycleParam(lp, dcl.BlockCreation) {
			return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Creation blocked by lifecycle params: %#v.", desired)}
		}
		create = true
	} else if dcl.HasLifecycleParam(lp, dcl.BlockAcquire) {
		return nil, dcl.ApplyInfeasibleError{
			Message: fmt.Sprintf("Resource already exists - apply blocked by lifecycle params: %#v.", initial),
		}
	} else {
		for _, d := range diffs {
			if d.RequiresRecreate {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) would require recreation", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}

	// 2.4 Imperative Request Planning
	var ops []folderApiOperation
	if create {
		ops = append(ops, &createFolderOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created plan: %#v", ops)

	// 2.5 Request Actuation
	for _, op := range ops {
		c.Config.Logger.InfoWithContextf(ctx, "Performing operation %T %+v", op, op)
		if err := op.do(ctx, desired, c); err != nil {
			c.Config.Logger.InfoWithContextf(ctx, "Failed operation %T %+v: %v", op, op, err)
			return nil, err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Finished operation %T %+v", op, op)
	}
	return applyFolderDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyFolderDiff(c *Client, ctx context.Context, desired *Folder, rawDesired *Folder, ops []folderApiOperation, opts ...dcl.ApplyOption) (*Folder, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetFolder(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createFolderOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapFolder(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeFolderNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeFolderNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeFolderDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractFolderFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractFolderFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffFolder(c, newDesired, newState)
	if err != nil {
		return newState, err
	}

	if len(newDiffs) == 0 {
		c.Config.Logger.InfoWithContext(ctx, "No diffs found. Apply was successful.")
	} else {
		c.Config.Logger.InfoWithContextf(ctx, "Found diffs: %v", newDiffs)
		diffMessages := make([]string, len(newDiffs))
		for i, d := range newDiffs {
			diffMessages[i] = fmt.Sprintf("%v", d)
		}
		return newState, dcl.DiffAfterApplyError{Diffs: diffMessages}
	}
	c.Config.Logger.InfoWithContext(ctx, "Done Apply.")
	return newState, nil
}

func (r *Folder) GetPolicy(basePath string) (string, string, *bytes.Buffer, error) {
	u := r.getPolicyURL(basePath)
	body := &bytes.Buffer{}
	body.WriteString(fmt.Sprintf(`{"options":{"requestedPolicyVersion": %d}}`, r.IAMPolicyVersion()))
	return u, "POST", body, nil
}
