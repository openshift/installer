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

type Zone struct {
	Name          *string            `json:"name"`
	DisplayName   *string            `json:"displayName"`
	Uid           *string            `json:"uid"`
	CreateTime    *string            `json:"createTime"`
	UpdateTime    *string            `json:"updateTime"`
	Labels        map[string]string  `json:"labels"`
	Description   *string            `json:"description"`
	State         *ZoneStateEnum     `json:"state"`
	Type          *ZoneTypeEnum      `json:"type"`
	DiscoverySpec *ZoneDiscoverySpec `json:"discoverySpec"`
	ResourceSpec  *ZoneResourceSpec  `json:"resourceSpec"`
	AssetStatus   *ZoneAssetStatus   `json:"assetStatus"`
	Project       *string            `json:"project"`
	Location      *string            `json:"location"`
	Lake          *string            `json:"lake"`
}

func (r *Zone) String() string {
	return dcl.SprintResource(r)
}

// The enum ZoneStateEnum.
type ZoneStateEnum string

// ZoneStateEnumRef returns a *ZoneStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func ZoneStateEnumRef(s string) *ZoneStateEnum {
	v := ZoneStateEnum(s)
	return &v
}

func (v ZoneStateEnum) Validate() error {
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
		Enum:  "ZoneStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ZoneTypeEnum.
type ZoneTypeEnum string

// ZoneTypeEnumRef returns a *ZoneTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func ZoneTypeEnumRef(s string) *ZoneTypeEnum {
	v := ZoneTypeEnum(s)
	return &v
}

func (v ZoneTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TYPE_UNSPECIFIED", "RAW", "CURATED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ZoneTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ZoneResourceSpecLocationTypeEnum.
type ZoneResourceSpecLocationTypeEnum string

// ZoneResourceSpecLocationTypeEnumRef returns a *ZoneResourceSpecLocationTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func ZoneResourceSpecLocationTypeEnumRef(s string) *ZoneResourceSpecLocationTypeEnum {
	v := ZoneResourceSpecLocationTypeEnum(s)
	return &v
}

func (v ZoneResourceSpecLocationTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"LOCATION_TYPE_UNSPECIFIED", "SINGLE_REGION", "MULTI_REGION"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ZoneResourceSpecLocationTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type ZoneDiscoverySpec struct {
	empty           bool                          `json:"-"`
	Enabled         *bool                         `json:"enabled"`
	IncludePatterns []string                      `json:"includePatterns"`
	ExcludePatterns []string                      `json:"excludePatterns"`
	CsvOptions      *ZoneDiscoverySpecCsvOptions  `json:"csvOptions"`
	JsonOptions     *ZoneDiscoverySpecJsonOptions `json:"jsonOptions"`
	Schedule        *string                       `json:"schedule"`
}

type jsonZoneDiscoverySpec ZoneDiscoverySpec

func (r *ZoneDiscoverySpec) UnmarshalJSON(data []byte) error {
	var res jsonZoneDiscoverySpec
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyZoneDiscoverySpec
	} else {

		r.Enabled = res.Enabled

		r.IncludePatterns = res.IncludePatterns

		r.ExcludePatterns = res.ExcludePatterns

		r.CsvOptions = res.CsvOptions

		r.JsonOptions = res.JsonOptions

		r.Schedule = res.Schedule

	}
	return nil
}

// This object is used to assert a desired state where this ZoneDiscoverySpec is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyZoneDiscoverySpec *ZoneDiscoverySpec = &ZoneDiscoverySpec{empty: true}

func (r *ZoneDiscoverySpec) Empty() bool {
	return r.empty
}

func (r *ZoneDiscoverySpec) String() string {
	return dcl.SprintResource(r)
}

func (r *ZoneDiscoverySpec) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ZoneDiscoverySpecCsvOptions struct {
	empty                bool    `json:"-"`
	HeaderRows           *int64  `json:"headerRows"`
	Delimiter            *string `json:"delimiter"`
	Encoding             *string `json:"encoding"`
	DisableTypeInference *bool   `json:"disableTypeInference"`
}

type jsonZoneDiscoverySpecCsvOptions ZoneDiscoverySpecCsvOptions

func (r *ZoneDiscoverySpecCsvOptions) UnmarshalJSON(data []byte) error {
	var res jsonZoneDiscoverySpecCsvOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyZoneDiscoverySpecCsvOptions
	} else {

		r.HeaderRows = res.HeaderRows

		r.Delimiter = res.Delimiter

		r.Encoding = res.Encoding

		r.DisableTypeInference = res.DisableTypeInference

	}
	return nil
}

// This object is used to assert a desired state where this ZoneDiscoverySpecCsvOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyZoneDiscoverySpecCsvOptions *ZoneDiscoverySpecCsvOptions = &ZoneDiscoverySpecCsvOptions{empty: true}

func (r *ZoneDiscoverySpecCsvOptions) Empty() bool {
	return r.empty
}

func (r *ZoneDiscoverySpecCsvOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *ZoneDiscoverySpecCsvOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ZoneDiscoverySpecJsonOptions struct {
	empty                bool    `json:"-"`
	Encoding             *string `json:"encoding"`
	DisableTypeInference *bool   `json:"disableTypeInference"`
}

type jsonZoneDiscoverySpecJsonOptions ZoneDiscoverySpecJsonOptions

func (r *ZoneDiscoverySpecJsonOptions) UnmarshalJSON(data []byte) error {
	var res jsonZoneDiscoverySpecJsonOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyZoneDiscoverySpecJsonOptions
	} else {

		r.Encoding = res.Encoding

		r.DisableTypeInference = res.DisableTypeInference

	}
	return nil
}

// This object is used to assert a desired state where this ZoneDiscoverySpecJsonOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyZoneDiscoverySpecJsonOptions *ZoneDiscoverySpecJsonOptions = &ZoneDiscoverySpecJsonOptions{empty: true}

func (r *ZoneDiscoverySpecJsonOptions) Empty() bool {
	return r.empty
}

func (r *ZoneDiscoverySpecJsonOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *ZoneDiscoverySpecJsonOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ZoneResourceSpec struct {
	empty        bool                              `json:"-"`
	LocationType *ZoneResourceSpecLocationTypeEnum `json:"locationType"`
}

type jsonZoneResourceSpec ZoneResourceSpec

func (r *ZoneResourceSpec) UnmarshalJSON(data []byte) error {
	var res jsonZoneResourceSpec
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyZoneResourceSpec
	} else {

		r.LocationType = res.LocationType

	}
	return nil
}

// This object is used to assert a desired state where this ZoneResourceSpec is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyZoneResourceSpec *ZoneResourceSpec = &ZoneResourceSpec{empty: true}

func (r *ZoneResourceSpec) Empty() bool {
	return r.empty
}

func (r *ZoneResourceSpec) String() string {
	return dcl.SprintResource(r)
}

func (r *ZoneResourceSpec) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ZoneAssetStatus struct {
	empty                        bool    `json:"-"`
	UpdateTime                   *string `json:"updateTime"`
	ActiveAssets                 *int64  `json:"activeAssets"`
	SecurityPolicyApplyingAssets *int64  `json:"securityPolicyApplyingAssets"`
}

type jsonZoneAssetStatus ZoneAssetStatus

func (r *ZoneAssetStatus) UnmarshalJSON(data []byte) error {
	var res jsonZoneAssetStatus
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyZoneAssetStatus
	} else {

		r.UpdateTime = res.UpdateTime

		r.ActiveAssets = res.ActiveAssets

		r.SecurityPolicyApplyingAssets = res.SecurityPolicyApplyingAssets

	}
	return nil
}

// This object is used to assert a desired state where this ZoneAssetStatus is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyZoneAssetStatus *ZoneAssetStatus = &ZoneAssetStatus{empty: true}

func (r *ZoneAssetStatus) Empty() bool {
	return r.empty
}

func (r *ZoneAssetStatus) String() string {
	return dcl.SprintResource(r)
}

func (r *ZoneAssetStatus) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Zone) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "dataplex",
		Type:    "Zone",
		Version: "dataplex",
	}
}

func (r *Zone) ID() (string, error) {
	if err := extractZoneFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":           dcl.ValueOrEmptyString(nr.Name),
		"display_name":   dcl.ValueOrEmptyString(nr.DisplayName),
		"uid":            dcl.ValueOrEmptyString(nr.Uid),
		"create_time":    dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":    dcl.ValueOrEmptyString(nr.UpdateTime),
		"labels":         dcl.ValueOrEmptyString(nr.Labels),
		"description":    dcl.ValueOrEmptyString(nr.Description),
		"state":          dcl.ValueOrEmptyString(nr.State),
		"type":           dcl.ValueOrEmptyString(nr.Type),
		"discovery_spec": dcl.ValueOrEmptyString(nr.DiscoverySpec),
		"resource_spec":  dcl.ValueOrEmptyString(nr.ResourceSpec),
		"asset_status":   dcl.ValueOrEmptyString(nr.AssetStatus),
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"lake":           dcl.ValueOrEmptyString(nr.Lake),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{name}}", params), nil
}

const ZoneMaxPage = -1

type ZoneList struct {
	Items []*Zone

	nextToken string

	pageSize int32

	resource *Zone
}

func (l *ZoneList) HasNext() bool {
	return l.nextToken != ""
}

func (l *ZoneList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listZone(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListZone(ctx context.Context, project, location, lake string) (*ZoneList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListZoneWithMaxResults(ctx, project, location, lake, ZoneMaxPage)

}

func (c *Client) ListZoneWithMaxResults(ctx context.Context, project, location, lake string, pageSize int32) (*ZoneList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Zone{
		Project:  &project,
		Location: &location,
		Lake:     &lake,
	}
	items, token, err := c.listZone(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &ZoneList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetZone(ctx context.Context, r *Zone) (*Zone, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractZoneFields(r)

	b, err := c.getZoneRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalZone(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Lake = r.Lake
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeZoneNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractZoneFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteZone(ctx context.Context, r *Zone) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Zone resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Zone...")
	deleteOp := deleteZoneOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllZone deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllZone(ctx context.Context, project, location, lake string, filter func(*Zone) bool) error {
	listObj, err := c.ListZone(ctx, project, location, lake)
	if err != nil {
		return err
	}

	err = c.deleteAllZone(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllZone(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyZone(ctx context.Context, rawDesired *Zone, opts ...dcl.ApplyOption) (*Zone, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Zone
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyZoneHelper(c, ctx, rawDesired, opts...)
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

func applyZoneHelper(c *Client, ctx context.Context, rawDesired *Zone, opts ...dcl.ApplyOption) (*Zone, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyZone...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractZoneFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.zoneDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToZoneDiffs(c.Config, fieldDiffs, opts)
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
	var ops []zoneApiOperation
	if create {
		ops = append(ops, &createZoneOperation{})
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
	return applyZoneDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyZoneDiff(c *Client, ctx context.Context, desired *Zone, rawDesired *Zone, ops []zoneApiOperation, opts ...dcl.ApplyOption) (*Zone, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetZone(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createZoneOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapZone(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeZoneNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeZoneNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeZoneDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractZoneFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractZoneFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffZone(c, newDesired, newState)
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

func (r *Zone) GetPolicy(basePath string) (string, string, *bytes.Buffer, error) {
	u := r.getPolicyURL(basePath)
	body := &bytes.Buffer{}
	u, err := dcl.AddQueryParams(u, map[string]string{"optionsRequestedPolicyVersion": fmt.Sprintf("%d", r.IAMPolicyVersion())})
	if err != nil {
		return "", "", nil, err
	}
	return u, "", body, nil
}
