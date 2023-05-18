// Copyright 2023 Google LLC. All Rights Reserved.
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
package dataplex

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Lake struct {
	Name            *string              `json:"name"`
	DisplayName     *string              `json:"displayName"`
	Uid             *string              `json:"uid"`
	CreateTime      *string              `json:"createTime"`
	UpdateTime      *string              `json:"updateTime"`
	Labels          map[string]string    `json:"labels"`
	Description     *string              `json:"description"`
	State           *LakeStateEnum       `json:"state"`
	ServiceAccount  *string              `json:"serviceAccount"`
	Metastore       *LakeMetastore       `json:"metastore"`
	AssetStatus     *LakeAssetStatus     `json:"assetStatus"`
	MetastoreStatus *LakeMetastoreStatus `json:"metastoreStatus"`
	Project         *string              `json:"project"`
	Location        *string              `json:"location"`
}

func (r *Lake) String() string {
	return dcl.SprintResource(r)
}

// The enum LakeStateEnum.
type LakeStateEnum string

// LakeStateEnumRef returns a *LakeStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func LakeStateEnumRef(s string) *LakeStateEnum {
	v := LakeStateEnum(s)
	return &v
}

func (v LakeStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "ACTIVE", "CREATING", "DELETING", "ACTION_REQUIRED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "LakeStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum LakeMetastoreStatusStateEnum.
type LakeMetastoreStatusStateEnum string

// LakeMetastoreStatusStateEnumRef returns a *LakeMetastoreStatusStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func LakeMetastoreStatusStateEnumRef(s string) *LakeMetastoreStatusStateEnum {
	v := LakeMetastoreStatusStateEnum(s)
	return &v
}

func (v LakeMetastoreStatusStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "NONE", "READY", "UPDATING", "ERROR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "LakeMetastoreStatusStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type LakeMetastore struct {
	empty   bool    `json:"-"`
	Service *string `json:"service"`
}

type jsonLakeMetastore LakeMetastore

func (r *LakeMetastore) UnmarshalJSON(data []byte) error {
	var res jsonLakeMetastore
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLakeMetastore
	} else {

		r.Service = res.Service

	}
	return nil
}

// This object is used to assert a desired state where this LakeMetastore is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLakeMetastore *LakeMetastore = &LakeMetastore{empty: true}

func (r *LakeMetastore) Empty() bool {
	return r.empty
}

func (r *LakeMetastore) String() string {
	return dcl.SprintResource(r)
}

func (r *LakeMetastore) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type LakeAssetStatus struct {
	empty                        bool    `json:"-"`
	UpdateTime                   *string `json:"updateTime"`
	ActiveAssets                 *int64  `json:"activeAssets"`
	SecurityPolicyApplyingAssets *int64  `json:"securityPolicyApplyingAssets"`
}

type jsonLakeAssetStatus LakeAssetStatus

func (r *LakeAssetStatus) UnmarshalJSON(data []byte) error {
	var res jsonLakeAssetStatus
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLakeAssetStatus
	} else {

		r.UpdateTime = res.UpdateTime

		r.ActiveAssets = res.ActiveAssets

		r.SecurityPolicyApplyingAssets = res.SecurityPolicyApplyingAssets

	}
	return nil
}

// This object is used to assert a desired state where this LakeAssetStatus is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLakeAssetStatus *LakeAssetStatus = &LakeAssetStatus{empty: true}

func (r *LakeAssetStatus) Empty() bool {
	return r.empty
}

func (r *LakeAssetStatus) String() string {
	return dcl.SprintResource(r)
}

func (r *LakeAssetStatus) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type LakeMetastoreStatus struct {
	empty      bool                          `json:"-"`
	State      *LakeMetastoreStatusStateEnum `json:"state"`
	Message    *string                       `json:"message"`
	UpdateTime *string                       `json:"updateTime"`
	Endpoint   *string                       `json:"endpoint"`
}

type jsonLakeMetastoreStatus LakeMetastoreStatus

func (r *LakeMetastoreStatus) UnmarshalJSON(data []byte) error {
	var res jsonLakeMetastoreStatus
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLakeMetastoreStatus
	} else {

		r.State = res.State

		r.Message = res.Message

		r.UpdateTime = res.UpdateTime

		r.Endpoint = res.Endpoint

	}
	return nil
}

// This object is used to assert a desired state where this LakeMetastoreStatus is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLakeMetastoreStatus *LakeMetastoreStatus = &LakeMetastoreStatus{empty: true}

func (r *LakeMetastoreStatus) Empty() bool {
	return r.empty
}

func (r *LakeMetastoreStatus) String() string {
	return dcl.SprintResource(r)
}

func (r *LakeMetastoreStatus) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Lake) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "dataplex",
		Type:    "Lake",
		Version: "dataplex",
	}
}

func (r *Lake) ID() (string, error) {
	if err := extractLakeFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":             dcl.ValueOrEmptyString(nr.Name),
		"display_name":     dcl.ValueOrEmptyString(nr.DisplayName),
		"uid":              dcl.ValueOrEmptyString(nr.Uid),
		"create_time":      dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":      dcl.ValueOrEmptyString(nr.UpdateTime),
		"labels":           dcl.ValueOrEmptyString(nr.Labels),
		"description":      dcl.ValueOrEmptyString(nr.Description),
		"state":            dcl.ValueOrEmptyString(nr.State),
		"service_account":  dcl.ValueOrEmptyString(nr.ServiceAccount),
		"metastore":        dcl.ValueOrEmptyString(nr.Metastore),
		"asset_status":     dcl.ValueOrEmptyString(nr.AssetStatus),
		"metastore_status": dcl.ValueOrEmptyString(nr.MetastoreStatus),
		"project":          dcl.ValueOrEmptyString(nr.Project),
		"location":         dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/lakes/{{name}}", params), nil
}

const LakeMaxPage = -1

type LakeList struct {
	Items []*Lake

	nextToken string

	pageSize int32

	resource *Lake
}

func (l *LakeList) HasNext() bool {
	return l.nextToken != ""
}

func (l *LakeList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listLake(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListLake(ctx context.Context, project, location string) (*LakeList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListLakeWithMaxResults(ctx, project, location, LakeMaxPage)

}

func (c *Client) ListLakeWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*LakeList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Lake{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listLake(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &LakeList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetLake(ctx context.Context, r *Lake) (*Lake, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractLakeFields(r)

	b, err := c.getLakeRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalLake(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeLakeNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractLakeFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteLake(ctx context.Context, r *Lake) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Lake resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Lake...")
	deleteOp := deleteLakeOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllLake deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllLake(ctx context.Context, project, location string, filter func(*Lake) bool) error {
	listObj, err := c.ListLake(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllLake(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllLake(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyLake(ctx context.Context, rawDesired *Lake, opts ...dcl.ApplyOption) (*Lake, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Lake
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyLakeHelper(c, ctx, rawDesired, opts...)
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

func applyLakeHelper(c *Client, ctx context.Context, rawDesired *Lake, opts ...dcl.ApplyOption) (*Lake, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyLake...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractLakeFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.lakeDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToLakeDiffs(c.Config, fieldDiffs, opts)
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
	var ops []lakeApiOperation
	if create {
		ops = append(ops, &createLakeOperation{})
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
	return applyLakeDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyLakeDiff(c *Client, ctx context.Context, desired *Lake, rawDesired *Lake, ops []lakeApiOperation, opts ...dcl.ApplyOption) (*Lake, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetLake(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createLakeOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapLake(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeLakeNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeLakeNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeLakeDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractLakeFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractLakeFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffLake(c, newDesired, newState)
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

func (r *Lake) GetPolicy(basePath string) (string, string, *bytes.Buffer, error) {
	u := r.getPolicyURL(basePath)
	body := &bytes.Buffer{}
	u, err := dcl.AddQueryParams(u, map[string]string{"optionsRequestedPolicyVersion": fmt.Sprintf("%d", r.IAMPolicyVersion())})
	if err != nil {
		return "", "", nil, err
	}
	return u, "", body, nil
}
