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
package logging

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type LogMetric struct {
	Name             *string                    `json:"name"`
	Description      *string                    `json:"description"`
	Filter           *string                    `json:"filter"`
	Disabled         *bool                      `json:"disabled"`
	MetricDescriptor *LogMetricMetricDescriptor `json:"metricDescriptor"`
	ValueExtractor   *string                    `json:"valueExtractor"`
	LabelExtractors  map[string]string          `json:"labelExtractors"`
	BucketOptions    *LogMetricBucketOptions    `json:"bucketOptions"`
	CreateTime       *string                    `json:"createTime"`
	UpdateTime       *string                    `json:"updateTime"`
	Project          *string                    `json:"project"`
}

func (r *LogMetric) String() string {
	return dcl.SprintResource(r)
}

// The enum LogMetricMetricDescriptorLabelsValueTypeEnum.
type LogMetricMetricDescriptorLabelsValueTypeEnum string

// LogMetricMetricDescriptorLabelsValueTypeEnumRef returns a *LogMetricMetricDescriptorLabelsValueTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func LogMetricMetricDescriptorLabelsValueTypeEnumRef(s string) *LogMetricMetricDescriptorLabelsValueTypeEnum {
	v := LogMetricMetricDescriptorLabelsValueTypeEnum(s)
	return &v
}

func (v LogMetricMetricDescriptorLabelsValueTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STRING", "BOOL", "INT64", "DOUBLE", "DISTRIBUTION", "MONEY"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "LogMetricMetricDescriptorLabelsValueTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum LogMetricMetricDescriptorMetricKindEnum.
type LogMetricMetricDescriptorMetricKindEnum string

// LogMetricMetricDescriptorMetricKindEnumRef returns a *LogMetricMetricDescriptorMetricKindEnum with the value of string s
// If the empty string is provided, nil is returned.
func LogMetricMetricDescriptorMetricKindEnumRef(s string) *LogMetricMetricDescriptorMetricKindEnum {
	v := LogMetricMetricDescriptorMetricKindEnum(s)
	return &v
}

func (v LogMetricMetricDescriptorMetricKindEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"GAUGE", "DELTA", "CUMULATIVE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "LogMetricMetricDescriptorMetricKindEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum LogMetricMetricDescriptorValueTypeEnum.
type LogMetricMetricDescriptorValueTypeEnum string

// LogMetricMetricDescriptorValueTypeEnumRef returns a *LogMetricMetricDescriptorValueTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func LogMetricMetricDescriptorValueTypeEnumRef(s string) *LogMetricMetricDescriptorValueTypeEnum {
	v := LogMetricMetricDescriptorValueTypeEnum(s)
	return &v
}

func (v LogMetricMetricDescriptorValueTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STRING", "BOOL", "INT64", "DOUBLE", "DISTRIBUTION", "MONEY"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "LogMetricMetricDescriptorValueTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum LogMetricMetricDescriptorLaunchStageEnum.
type LogMetricMetricDescriptorLaunchStageEnum string

// LogMetricMetricDescriptorLaunchStageEnumRef returns a *LogMetricMetricDescriptorLaunchStageEnum with the value of string s
// If the empty string is provided, nil is returned.
func LogMetricMetricDescriptorLaunchStageEnumRef(s string) *LogMetricMetricDescriptorLaunchStageEnum {
	v := LogMetricMetricDescriptorLaunchStageEnum(s)
	return &v
}

func (v LogMetricMetricDescriptorLaunchStageEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"UNIMPLEMENTED", "PRELAUNCH", "EARLY_ACCESS", "ALPHA", "BETA", "GA", "DEPRECATED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "LogMetricMetricDescriptorLaunchStageEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type LogMetricMetricDescriptor struct {
	empty                  bool                                      `json:"-"`
	Name                   *string                                   `json:"name"`
	Type                   *string                                   `json:"type"`
	Labels                 []LogMetricMetricDescriptorLabels         `json:"labels"`
	MetricKind             *LogMetricMetricDescriptorMetricKindEnum  `json:"metricKind"`
	ValueType              *LogMetricMetricDescriptorValueTypeEnum   `json:"valueType"`
	Unit                   *string                                   `json:"unit"`
	Description            *string                                   `json:"description"`
	DisplayName            *string                                   `json:"displayName"`
	Metadata               *LogMetricMetricDescriptorMetadata        `json:"metadata"`
	LaunchStage            *LogMetricMetricDescriptorLaunchStageEnum `json:"launchStage"`
	MonitoredResourceTypes []string                                  `json:"monitoredResourceTypes"`
}

type jsonLogMetricMetricDescriptor LogMetricMetricDescriptor

func (r *LogMetricMetricDescriptor) UnmarshalJSON(data []byte) error {
	var res jsonLogMetricMetricDescriptor
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLogMetricMetricDescriptor
	} else {

		r.Name = res.Name

		r.Type = res.Type

		r.Labels = res.Labels

		r.MetricKind = res.MetricKind

		r.ValueType = res.ValueType

		r.Unit = res.Unit

		r.Description = res.Description

		r.DisplayName = res.DisplayName

		r.Metadata = res.Metadata

		r.LaunchStage = res.LaunchStage

		r.MonitoredResourceTypes = res.MonitoredResourceTypes

	}
	return nil
}

// This object is used to assert a desired state where this LogMetricMetricDescriptor is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLogMetricMetricDescriptor *LogMetricMetricDescriptor = &LogMetricMetricDescriptor{empty: true}

func (r *LogMetricMetricDescriptor) Empty() bool {
	return r.empty
}

func (r *LogMetricMetricDescriptor) String() string {
	return dcl.SprintResource(r)
}

func (r *LogMetricMetricDescriptor) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type LogMetricMetricDescriptorLabels struct {
	empty       bool                                          `json:"-"`
	Key         *string                                       `json:"key"`
	ValueType   *LogMetricMetricDescriptorLabelsValueTypeEnum `json:"valueType"`
	Description *string                                       `json:"description"`
}

type jsonLogMetricMetricDescriptorLabels LogMetricMetricDescriptorLabels

func (r *LogMetricMetricDescriptorLabels) UnmarshalJSON(data []byte) error {
	var res jsonLogMetricMetricDescriptorLabels
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLogMetricMetricDescriptorLabels
	} else {

		r.Key = res.Key

		r.ValueType = res.ValueType

		r.Description = res.Description

	}
	return nil
}

// This object is used to assert a desired state where this LogMetricMetricDescriptorLabels is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLogMetricMetricDescriptorLabels *LogMetricMetricDescriptorLabels = &LogMetricMetricDescriptorLabels{empty: true}

func (r *LogMetricMetricDescriptorLabels) Empty() bool {
	return r.empty
}

func (r *LogMetricMetricDescriptorLabels) String() string {
	return dcl.SprintResource(r)
}

func (r *LogMetricMetricDescriptorLabels) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type LogMetricMetricDescriptorMetadata struct {
	empty        bool    `json:"-"`
	SamplePeriod *string `json:"samplePeriod"`
	IngestDelay  *string `json:"ingestDelay"`
}

type jsonLogMetricMetricDescriptorMetadata LogMetricMetricDescriptorMetadata

func (r *LogMetricMetricDescriptorMetadata) UnmarshalJSON(data []byte) error {
	var res jsonLogMetricMetricDescriptorMetadata
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLogMetricMetricDescriptorMetadata
	} else {

		r.SamplePeriod = res.SamplePeriod

		r.IngestDelay = res.IngestDelay

	}
	return nil
}

// This object is used to assert a desired state where this LogMetricMetricDescriptorMetadata is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLogMetricMetricDescriptorMetadata *LogMetricMetricDescriptorMetadata = &LogMetricMetricDescriptorMetadata{empty: true}

func (r *LogMetricMetricDescriptorMetadata) Empty() bool {
	return r.empty
}

func (r *LogMetricMetricDescriptorMetadata) String() string {
	return dcl.SprintResource(r)
}

func (r *LogMetricMetricDescriptorMetadata) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type LogMetricBucketOptions struct {
	empty              bool                                      `json:"-"`
	LinearBuckets      *LogMetricBucketOptionsLinearBuckets      `json:"linearBuckets"`
	ExponentialBuckets *LogMetricBucketOptionsExponentialBuckets `json:"exponentialBuckets"`
	ExplicitBuckets    *LogMetricBucketOptionsExplicitBuckets    `json:"explicitBuckets"`
}

type jsonLogMetricBucketOptions LogMetricBucketOptions

func (r *LogMetricBucketOptions) UnmarshalJSON(data []byte) error {
	var res jsonLogMetricBucketOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLogMetricBucketOptions
	} else {

		r.LinearBuckets = res.LinearBuckets

		r.ExponentialBuckets = res.ExponentialBuckets

		r.ExplicitBuckets = res.ExplicitBuckets

	}
	return nil
}

// This object is used to assert a desired state where this LogMetricBucketOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLogMetricBucketOptions *LogMetricBucketOptions = &LogMetricBucketOptions{empty: true}

func (r *LogMetricBucketOptions) Empty() bool {
	return r.empty
}

func (r *LogMetricBucketOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *LogMetricBucketOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type LogMetricBucketOptionsLinearBuckets struct {
	empty            bool     `json:"-"`
	NumFiniteBuckets *int64   `json:"numFiniteBuckets"`
	Width            *float64 `json:"width"`
	Offset           *float64 `json:"offset"`
}

type jsonLogMetricBucketOptionsLinearBuckets LogMetricBucketOptionsLinearBuckets

func (r *LogMetricBucketOptionsLinearBuckets) UnmarshalJSON(data []byte) error {
	var res jsonLogMetricBucketOptionsLinearBuckets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLogMetricBucketOptionsLinearBuckets
	} else {

		r.NumFiniteBuckets = res.NumFiniteBuckets

		r.Width = res.Width

		r.Offset = res.Offset

	}
	return nil
}

// This object is used to assert a desired state where this LogMetricBucketOptionsLinearBuckets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLogMetricBucketOptionsLinearBuckets *LogMetricBucketOptionsLinearBuckets = &LogMetricBucketOptionsLinearBuckets{empty: true}

func (r *LogMetricBucketOptionsLinearBuckets) Empty() bool {
	return r.empty
}

func (r *LogMetricBucketOptionsLinearBuckets) String() string {
	return dcl.SprintResource(r)
}

func (r *LogMetricBucketOptionsLinearBuckets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type LogMetricBucketOptionsExponentialBuckets struct {
	empty            bool     `json:"-"`
	NumFiniteBuckets *int64   `json:"numFiniteBuckets"`
	GrowthFactor     *float64 `json:"growthFactor"`
	Scale            *float64 `json:"scale"`
}

type jsonLogMetricBucketOptionsExponentialBuckets LogMetricBucketOptionsExponentialBuckets

func (r *LogMetricBucketOptionsExponentialBuckets) UnmarshalJSON(data []byte) error {
	var res jsonLogMetricBucketOptionsExponentialBuckets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLogMetricBucketOptionsExponentialBuckets
	} else {

		r.NumFiniteBuckets = res.NumFiniteBuckets

		r.GrowthFactor = res.GrowthFactor

		r.Scale = res.Scale

	}
	return nil
}

// This object is used to assert a desired state where this LogMetricBucketOptionsExponentialBuckets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLogMetricBucketOptionsExponentialBuckets *LogMetricBucketOptionsExponentialBuckets = &LogMetricBucketOptionsExponentialBuckets{empty: true}

func (r *LogMetricBucketOptionsExponentialBuckets) Empty() bool {
	return r.empty
}

func (r *LogMetricBucketOptionsExponentialBuckets) String() string {
	return dcl.SprintResource(r)
}

func (r *LogMetricBucketOptionsExponentialBuckets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type LogMetricBucketOptionsExplicitBuckets struct {
	empty  bool      `json:"-"`
	Bounds []float64 `json:"bounds"`
}

type jsonLogMetricBucketOptionsExplicitBuckets LogMetricBucketOptionsExplicitBuckets

func (r *LogMetricBucketOptionsExplicitBuckets) UnmarshalJSON(data []byte) error {
	var res jsonLogMetricBucketOptionsExplicitBuckets
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyLogMetricBucketOptionsExplicitBuckets
	} else {

		r.Bounds = res.Bounds

	}
	return nil
}

// This object is used to assert a desired state where this LogMetricBucketOptionsExplicitBuckets is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyLogMetricBucketOptionsExplicitBuckets *LogMetricBucketOptionsExplicitBuckets = &LogMetricBucketOptionsExplicitBuckets{empty: true}

func (r *LogMetricBucketOptionsExplicitBuckets) Empty() bool {
	return r.empty
}

func (r *LogMetricBucketOptionsExplicitBuckets) String() string {
	return dcl.SprintResource(r)
}

func (r *LogMetricBucketOptionsExplicitBuckets) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *LogMetric) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "logging",
		Type:    "LogMetric",
		Version: "logging",
	}
}

func (r *LogMetric) ID() (string, error) {
	if err := extractLogMetricFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":              dcl.ValueOrEmptyString(nr.Name),
		"description":       dcl.ValueOrEmptyString(nr.Description),
		"filter":            dcl.ValueOrEmptyString(nr.Filter),
		"disabled":          dcl.ValueOrEmptyString(nr.Disabled),
		"metric_descriptor": dcl.ValueOrEmptyString(nr.MetricDescriptor),
		"value_extractor":   dcl.ValueOrEmptyString(nr.ValueExtractor),
		"label_extractors":  dcl.ValueOrEmptyString(nr.LabelExtractors),
		"bucket_options":    dcl.ValueOrEmptyString(nr.BucketOptions),
		"create_time":       dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":       dcl.ValueOrEmptyString(nr.UpdateTime),
		"project":           dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.Nprintf("projects/{{project}}/metrics/{{name}}", params), nil
}

const LogMetricMaxPage = -1

type LogMetricList struct {
	Items []*LogMetric

	nextToken string

	pageSize int32

	resource *LogMetric
}

func (l *LogMetricList) HasNext() bool {
	return l.nextToken != ""
}

func (l *LogMetricList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listLogMetric(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListLogMetric(ctx context.Context, project string) (*LogMetricList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListLogMetricWithMaxResults(ctx, project, LogMetricMaxPage)

}

func (c *Client) ListLogMetricWithMaxResults(ctx context.Context, project string, pageSize int32) (*LogMetricList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &LogMetric{
		Project: &project,
	}
	items, token, err := c.listLogMetric(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &LogMetricList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetLogMetric(ctx context.Context, r *LogMetric) (*LogMetric, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractLogMetricFields(r)

	b, err := c.getLogMetricRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalLogMetric(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeLogMetricNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractLogMetricFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteLogMetric(ctx context.Context, r *LogMetric) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("LogMetric resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting LogMetric...")
	deleteOp := deleteLogMetricOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllLogMetric deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllLogMetric(ctx context.Context, project string, filter func(*LogMetric) bool) error {
	listObj, err := c.ListLogMetric(ctx, project)
	if err != nil {
		return err
	}

	err = c.deleteAllLogMetric(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllLogMetric(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyLogMetric(ctx context.Context, rawDesired *LogMetric, opts ...dcl.ApplyOption) (*LogMetric, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *LogMetric
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyLogMetricHelper(c, ctx, rawDesired, opts...)
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

func applyLogMetricHelper(c *Client, ctx context.Context, rawDesired *LogMetric, opts ...dcl.ApplyOption) (*LogMetric, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyLogMetric...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractLogMetricFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.logMetricDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToLogMetricDiffs(c.Config, fieldDiffs, opts)
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
	var ops []logMetricApiOperation
	if create {
		ops = append(ops, &createLogMetricOperation{})
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
	return applyLogMetricDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyLogMetricDiff(c *Client, ctx context.Context, desired *LogMetric, rawDesired *LogMetric, ops []logMetricApiOperation, opts ...dcl.ApplyOption) (*LogMetric, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetLogMetric(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createLogMetricOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapLogMetric(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeLogMetricNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeLogMetricNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeLogMetricDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractLogMetricFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractLogMetricFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffLogMetric(c, newDesired, newState)
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
