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

type Asset struct {
	Name            *string               `json:"name"`
	DisplayName     *string               `json:"displayName"`
	Uid             *string               `json:"uid"`
	CreateTime      *string               `json:"createTime"`
	UpdateTime      *string               `json:"updateTime"`
	Labels          map[string]string     `json:"labels"`
	Description     *string               `json:"description"`
	State           *AssetStateEnum       `json:"state"`
	ResourceSpec    *AssetResourceSpec    `json:"resourceSpec"`
	ResourceStatus  *AssetResourceStatus  `json:"resourceStatus"`
	SecurityStatus  *AssetSecurityStatus  `json:"securityStatus"`
	DiscoverySpec   *AssetDiscoverySpec   `json:"discoverySpec"`
	DiscoveryStatus *AssetDiscoveryStatus `json:"discoveryStatus"`
	Project         *string               `json:"project"`
	Location        *string               `json:"location"`
	Lake            *string               `json:"lake"`
	DataplexZone    *string               `json:"dataplexZone"`
}

func (r *Asset) String() string {
	return dcl.SprintResource(r)
}

// The enum AssetStateEnum.
type AssetStateEnum string

// AssetStateEnumRef returns a *AssetStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func AssetStateEnumRef(s string) *AssetStateEnum {
	v := AssetStateEnum(s)
	return &v
}

func (v AssetStateEnum) Validate() error {
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
		Enum:  "AssetStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum AssetResourceSpecTypeEnum.
type AssetResourceSpecTypeEnum string

// AssetResourceSpecTypeEnumRef returns a *AssetResourceSpecTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func AssetResourceSpecTypeEnumRef(s string) *AssetResourceSpecTypeEnum {
	v := AssetResourceSpecTypeEnum(s)
	return &v
}

func (v AssetResourceSpecTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STORAGE_BUCKET", "BIGQUERY_DATASET"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "AssetResourceSpecTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum AssetResourceStatusStateEnum.
type AssetResourceStatusStateEnum string

// AssetResourceStatusStateEnumRef returns a *AssetResourceStatusStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func AssetResourceStatusStateEnumRef(s string) *AssetResourceStatusStateEnum {
	v := AssetResourceStatusStateEnum(s)
	return &v
}

func (v AssetResourceStatusStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "READY", "ERROR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "AssetResourceStatusStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum AssetSecurityStatusStateEnum.
type AssetSecurityStatusStateEnum string

// AssetSecurityStatusStateEnumRef returns a *AssetSecurityStatusStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func AssetSecurityStatusStateEnumRef(s string) *AssetSecurityStatusStateEnum {
	v := AssetSecurityStatusStateEnum(s)
	return &v
}

func (v AssetSecurityStatusStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "READY", "APPLYING", "ERROR"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "AssetSecurityStatusStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum AssetDiscoveryStatusStateEnum.
type AssetDiscoveryStatusStateEnum string

// AssetDiscoveryStatusStateEnumRef returns a *AssetDiscoveryStatusStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func AssetDiscoveryStatusStateEnumRef(s string) *AssetDiscoveryStatusStateEnum {
	v := AssetDiscoveryStatusStateEnum(s)
	return &v
}

func (v AssetDiscoveryStatusStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "SCHEDULED", "IN_PROGRESS", "PAUSED", "DISABLED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "AssetDiscoveryStatusStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type AssetResourceSpec struct {
	empty bool                       `json:"-"`
	Name  *string                    `json:"name"`
	Type  *AssetResourceSpecTypeEnum `json:"type"`
}

type jsonAssetResourceSpec AssetResourceSpec

func (r *AssetResourceSpec) UnmarshalJSON(data []byte) error {
	var res jsonAssetResourceSpec
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyAssetResourceSpec
	} else {

		r.Name = res.Name

		r.Type = res.Type

	}
	return nil
}

// This object is used to assert a desired state where this AssetResourceSpec is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyAssetResourceSpec *AssetResourceSpec = &AssetResourceSpec{empty: true}

func (r *AssetResourceSpec) Empty() bool {
	return r.empty
}

func (r *AssetResourceSpec) String() string {
	return dcl.SprintResource(r)
}

func (r *AssetResourceSpec) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type AssetResourceStatus struct {
	empty      bool                          `json:"-"`
	State      *AssetResourceStatusStateEnum `json:"state"`
	Message    *string                       `json:"message"`
	UpdateTime *string                       `json:"updateTime"`
}

type jsonAssetResourceStatus AssetResourceStatus

func (r *AssetResourceStatus) UnmarshalJSON(data []byte) error {
	var res jsonAssetResourceStatus
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyAssetResourceStatus
	} else {

		r.State = res.State

		r.Message = res.Message

		r.UpdateTime = res.UpdateTime

	}
	return nil
}

// This object is used to assert a desired state where this AssetResourceStatus is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyAssetResourceStatus *AssetResourceStatus = &AssetResourceStatus{empty: true}

func (r *AssetResourceStatus) Empty() bool {
	return r.empty
}

func (r *AssetResourceStatus) String() string {
	return dcl.SprintResource(r)
}

func (r *AssetResourceStatus) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type AssetSecurityStatus struct {
	empty      bool                          `json:"-"`
	State      *AssetSecurityStatusStateEnum `json:"state"`
	Message    *string                       `json:"message"`
	UpdateTime *string                       `json:"updateTime"`
}

type jsonAssetSecurityStatus AssetSecurityStatus

func (r *AssetSecurityStatus) UnmarshalJSON(data []byte) error {
	var res jsonAssetSecurityStatus
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyAssetSecurityStatus
	} else {

		r.State = res.State

		r.Message = res.Message

		r.UpdateTime = res.UpdateTime

	}
	return nil
}

// This object is used to assert a desired state where this AssetSecurityStatus is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyAssetSecurityStatus *AssetSecurityStatus = &AssetSecurityStatus{empty: true}

func (r *AssetSecurityStatus) Empty() bool {
	return r.empty
}

func (r *AssetSecurityStatus) String() string {
	return dcl.SprintResource(r)
}

func (r *AssetSecurityStatus) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type AssetDiscoverySpec struct {
	empty           bool                           `json:"-"`
	Enabled         *bool                          `json:"enabled"`
	IncludePatterns []string                       `json:"includePatterns"`
	ExcludePatterns []string                       `json:"excludePatterns"`
	CsvOptions      *AssetDiscoverySpecCsvOptions  `json:"csvOptions"`
	JsonOptions     *AssetDiscoverySpecJsonOptions `json:"jsonOptions"`
	Schedule        *string                        `json:"schedule"`
}

type jsonAssetDiscoverySpec AssetDiscoverySpec

func (r *AssetDiscoverySpec) UnmarshalJSON(data []byte) error {
	var res jsonAssetDiscoverySpec
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyAssetDiscoverySpec
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

// This object is used to assert a desired state where this AssetDiscoverySpec is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyAssetDiscoverySpec *AssetDiscoverySpec = &AssetDiscoverySpec{empty: true}

func (r *AssetDiscoverySpec) Empty() bool {
	return r.empty
}

func (r *AssetDiscoverySpec) String() string {
	return dcl.SprintResource(r)
}

func (r *AssetDiscoverySpec) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type AssetDiscoverySpecCsvOptions struct {
	empty                bool    `json:"-"`
	HeaderRows           *int64  `json:"headerRows"`
	Delimiter            *string `json:"delimiter"`
	Encoding             *string `json:"encoding"`
	DisableTypeInference *bool   `json:"disableTypeInference"`
}

type jsonAssetDiscoverySpecCsvOptions AssetDiscoverySpecCsvOptions

func (r *AssetDiscoverySpecCsvOptions) UnmarshalJSON(data []byte) error {
	var res jsonAssetDiscoverySpecCsvOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyAssetDiscoverySpecCsvOptions
	} else {

		r.HeaderRows = res.HeaderRows

		r.Delimiter = res.Delimiter

		r.Encoding = res.Encoding

		r.DisableTypeInference = res.DisableTypeInference

	}
	return nil
}

// This object is used to assert a desired state where this AssetDiscoverySpecCsvOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyAssetDiscoverySpecCsvOptions *AssetDiscoverySpecCsvOptions = &AssetDiscoverySpecCsvOptions{empty: true}

func (r *AssetDiscoverySpecCsvOptions) Empty() bool {
	return r.empty
}

func (r *AssetDiscoverySpecCsvOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *AssetDiscoverySpecCsvOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type AssetDiscoverySpecJsonOptions struct {
	empty                bool    `json:"-"`
	Encoding             *string `json:"encoding"`
	DisableTypeInference *bool   `json:"disableTypeInference"`
}

type jsonAssetDiscoverySpecJsonOptions AssetDiscoverySpecJsonOptions

func (r *AssetDiscoverySpecJsonOptions) UnmarshalJSON(data []byte) error {
	var res jsonAssetDiscoverySpecJsonOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyAssetDiscoverySpecJsonOptions
	} else {

		r.Encoding = res.Encoding

		r.DisableTypeInference = res.DisableTypeInference

	}
	return nil
}

// This object is used to assert a desired state where this AssetDiscoverySpecJsonOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyAssetDiscoverySpecJsonOptions *AssetDiscoverySpecJsonOptions = &AssetDiscoverySpecJsonOptions{empty: true}

func (r *AssetDiscoverySpecJsonOptions) Empty() bool {
	return r.empty
}

func (r *AssetDiscoverySpecJsonOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *AssetDiscoverySpecJsonOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type AssetDiscoveryStatus struct {
	empty           bool                           `json:"-"`
	State           *AssetDiscoveryStatusStateEnum `json:"state"`
	Message         *string                        `json:"message"`
	UpdateTime      *string                        `json:"updateTime"`
	LastRunTime     *string                        `json:"lastRunTime"`
	Stats           *AssetDiscoveryStatusStats     `json:"stats"`
	LastRunDuration *string                        `json:"lastRunDuration"`
}

type jsonAssetDiscoveryStatus AssetDiscoveryStatus

func (r *AssetDiscoveryStatus) UnmarshalJSON(data []byte) error {
	var res jsonAssetDiscoveryStatus
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyAssetDiscoveryStatus
	} else {

		r.State = res.State

		r.Message = res.Message

		r.UpdateTime = res.UpdateTime

		r.LastRunTime = res.LastRunTime

		r.Stats = res.Stats

		r.LastRunDuration = res.LastRunDuration

	}
	return nil
}

// This object is used to assert a desired state where this AssetDiscoveryStatus is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyAssetDiscoveryStatus *AssetDiscoveryStatus = &AssetDiscoveryStatus{empty: true}

func (r *AssetDiscoveryStatus) Empty() bool {
	return r.empty
}

func (r *AssetDiscoveryStatus) String() string {
	return dcl.SprintResource(r)
}

func (r *AssetDiscoveryStatus) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type AssetDiscoveryStatusStats struct {
	empty     bool   `json:"-"`
	DataItems *int64 `json:"dataItems"`
	DataSize  *int64 `json:"dataSize"`
	Tables    *int64 `json:"tables"`
	Filesets  *int64 `json:"filesets"`
}

type jsonAssetDiscoveryStatusStats AssetDiscoveryStatusStats

func (r *AssetDiscoveryStatusStats) UnmarshalJSON(data []byte) error {
	var res jsonAssetDiscoveryStatusStats
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyAssetDiscoveryStatusStats
	} else {

		r.DataItems = res.DataItems

		r.DataSize = res.DataSize

		r.Tables = res.Tables

		r.Filesets = res.Filesets

	}
	return nil
}

// This object is used to assert a desired state where this AssetDiscoveryStatusStats is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyAssetDiscoveryStatusStats *AssetDiscoveryStatusStats = &AssetDiscoveryStatusStats{empty: true}

func (r *AssetDiscoveryStatusStats) Empty() bool {
	return r.empty
}

func (r *AssetDiscoveryStatusStats) String() string {
	return dcl.SprintResource(r)
}

func (r *AssetDiscoveryStatusStats) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Asset) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "dataplex",
		Type:    "Asset",
		Version: "dataplex",
	}
}

func (r *Asset) ID() (string, error) {
	if err := extractAssetFields(r); err != nil {
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
		"resource_spec":    dcl.ValueOrEmptyString(nr.ResourceSpec),
		"resource_status":  dcl.ValueOrEmptyString(nr.ResourceStatus),
		"security_status":  dcl.ValueOrEmptyString(nr.SecurityStatus),
		"discovery_spec":   dcl.ValueOrEmptyString(nr.DiscoverySpec),
		"discovery_status": dcl.ValueOrEmptyString(nr.DiscoveryStatus),
		"project":          dcl.ValueOrEmptyString(nr.Project),
		"location":         dcl.ValueOrEmptyString(nr.Location),
		"lake":             dcl.ValueOrEmptyString(nr.Lake),
		"dataplex_zone":    dcl.ValueOrEmptyString(nr.DataplexZone),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{dataplex_zone}}/assets/{{name}}", params), nil
}

const AssetMaxPage = -1

type AssetList struct {
	Items []*Asset

	nextToken string

	pageSize int32

	resource *Asset
}

func (l *AssetList) HasNext() bool {
	return l.nextToken != ""
}

func (l *AssetList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listAsset(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListAsset(ctx context.Context, project, location, dataplexZone, lake string) (*AssetList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListAssetWithMaxResults(ctx, project, location, dataplexZone, lake, AssetMaxPage)

}

func (c *Client) ListAssetWithMaxResults(ctx context.Context, project, location, dataplexZone, lake string, pageSize int32) (*AssetList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Asset{
		Project:      &project,
		Location:     &location,
		DataplexZone: &dataplexZone,
		Lake:         &lake,
	}
	items, token, err := c.listAsset(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &AssetList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetAsset(ctx context.Context, r *Asset) (*Asset, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractAssetFields(r)

	b, err := c.getAssetRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalAsset(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.DataplexZone = r.DataplexZone
	result.Lake = r.Lake
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeAssetNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractAssetFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteAsset(ctx context.Context, r *Asset) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Asset resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Asset...")
	deleteOp := deleteAssetOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllAsset deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllAsset(ctx context.Context, project, location, dataplexZone, lake string, filter func(*Asset) bool) error {
	listObj, err := c.ListAsset(ctx, project, location, dataplexZone, lake)
	if err != nil {
		return err
	}

	err = c.deleteAllAsset(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllAsset(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyAsset(ctx context.Context, rawDesired *Asset, opts ...dcl.ApplyOption) (*Asset, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Asset
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyAssetHelper(c, ctx, rawDesired, opts...)
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

func applyAssetHelper(c *Client, ctx context.Context, rawDesired *Asset, opts ...dcl.ApplyOption) (*Asset, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyAsset...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractAssetFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.assetDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToAssetDiffs(c.Config, fieldDiffs, opts)
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
	var ops []assetApiOperation
	if create {
		ops = append(ops, &createAssetOperation{})
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
	return applyAssetDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyAssetDiff(c *Client, ctx context.Context, desired *Asset, rawDesired *Asset, ops []assetApiOperation, opts ...dcl.ApplyOption) (*Asset, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetAsset(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createAssetOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapAsset(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeAssetNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeAssetNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeAssetDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractAssetFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractAssetFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffAsset(c, newDesired, newState)
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

func (r *Asset) GetPolicy(basePath string) (string, string, *bytes.Buffer, error) {
	u := r.getPolicyURL(basePath)
	body := &bytes.Buffer{}
	u, err := dcl.AddQueryParams(u, map[string]string{"optionsRequestedPolicyVersion": fmt.Sprintf("%d", r.IAMPolicyVersion())})
	if err != nil {
		return "", "", nil, err
	}
	return u, "", body, nil
}
