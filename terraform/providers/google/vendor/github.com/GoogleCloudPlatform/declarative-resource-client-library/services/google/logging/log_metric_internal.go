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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *LogMetric) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "filter"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.MetricDescriptor) {
		if err := r.MetricDescriptor.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.BucketOptions) {
		if err := r.BucketOptions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *LogMetricMetricDescriptor) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Metadata) {
		if err := r.Metadata.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *LogMetricMetricDescriptorLabels) validate() error {
	return nil
}
func (r *LogMetricMetricDescriptorMetadata) validate() error {
	return nil
}
func (r *LogMetricBucketOptions) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"LinearBuckets", "ExponentialBuckets", "ExplicitBuckets"}, r.LinearBuckets, r.ExponentialBuckets, r.ExplicitBuckets); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.LinearBuckets) {
		if err := r.LinearBuckets.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ExponentialBuckets) {
		if err := r.ExponentialBuckets.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ExplicitBuckets) {
		if err := r.ExplicitBuckets.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *LogMetricBucketOptionsLinearBuckets) validate() error {
	return nil
}
func (r *LogMetricBucketOptionsExponentialBuckets) validate() error {
	return nil
}
func (r *LogMetricBucketOptionsExplicitBuckets) validate() error {
	return nil
}
func (r *LogMetric) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://logging.googleapis.com/v2/", params)
}

func (r *LogMetric) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/metrics/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *LogMetric) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/metrics", nr.basePath(), userBasePath, params), nil

}

func (r *LogMetric) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/metrics", nr.basePath(), userBasePath, params), nil

}

func (r *LogMetric) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/metrics/{{name}}", nr.basePath(), userBasePath, params), nil
}

// logMetricApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type logMetricApiOperation interface {
	do(context.Context, *LogMetric, *Client) error
}

// newUpdateLogMetricUpdateRequest creates a request for an
// LogMetric resource's update update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateLogMetricUpdateRequest(ctx context.Context, f *LogMetric, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.Filter; !dcl.IsEmptyValueIndirect(v) {
		req["filter"] = v
	}
	if v := f.Disabled; !dcl.IsEmptyValueIndirect(v) {
		req["disabled"] = v
	}
	if v, err := expandLogMetricMetricDescriptor(c, f.MetricDescriptor, res); err != nil {
		return nil, fmt.Errorf("error expanding MetricDescriptor into metricDescriptor: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["metricDescriptor"] = v
	}
	if v := f.ValueExtractor; !dcl.IsEmptyValueIndirect(v) {
		req["valueExtractor"] = v
	}
	if v := f.LabelExtractors; !dcl.IsEmptyValueIndirect(v) {
		req["labelExtractors"] = v
	}
	if v, err := expandLogMetricBucketOptions(c, f.BucketOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding BucketOptions into bucketOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["bucketOptions"] = v
	}
	return req, nil
}

// marshalUpdateLogMetricUpdateRequest converts the update into
// the final JSON request body.
func marshalUpdateLogMetricUpdateRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateLogMetricUpdateOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateLogMetricUpdateOperation) do(ctx context.Context, r *LogMetric, c *Client) error {
	_, err := c.GetLogMetric(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "update")
	if err != nil {
		return err
	}

	req, err := newUpdateLogMetricUpdateRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateLogMetricUpdateRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PUT", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listLogMetricRaw(ctx context.Context, r *LogMetric, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != LogMetricMaxPage {
		m["pageSize"] = fmt.Sprintf("%v", pageSize)
	}

	u, err = dcl.AddQueryParams(u, m)
	if err != nil {
		return nil, err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return nil, err
	}
	defer resp.Response.Body.Close()
	return ioutil.ReadAll(resp.Response.Body)
}

type listLogMetricOperation struct {
	Metrics []map[string]interface{} `json:"metrics"`
	Token   string                   `json:"nextPageToken"`
}

func (c *Client) listLogMetric(ctx context.Context, r *LogMetric, pageToken string, pageSize int32) ([]*LogMetric, string, error) {
	b, err := c.listLogMetricRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listLogMetricOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*LogMetric
	for _, v := range m.Metrics {
		res, err := unmarshalMapLogMetric(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllLogMetric(ctx context.Context, f func(*LogMetric) bool, resources []*LogMetric) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteLogMetric(ctx, res)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%v", strings.Join(errors, "\n"))
	} else {
		return nil
	}
}

type deleteLogMetricOperation struct{}

func (op *deleteLogMetricOperation) do(ctx context.Context, r *LogMetric, c *Client) error {
	r, err := c.GetLogMetric(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "LogMetric not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetLogMetric checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	_, err = dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return fmt.Errorf("failed to delete LogMetric: %w", err)
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetLogMetric(ctx, r)
		if dcl.IsNotFound(err) {
			return nil, nil
		}
		if retriesRemaining > 0 {
			retriesRemaining--
			return &dcl.RetryDetails{}, dcl.OperationNotDone{}
		}
		return nil, dcl.NotDeletedError{ExistingResource: r}
	}, c.Config.RetryProvider)
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createLogMetricOperation struct {
	response map[string]interface{}
}

func (op *createLogMetricOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createLogMetricOperation) do(ctx context.Context, r *LogMetric, c *Client) error {
	c.Config.Logger.InfoWithContextf(ctx, "Attempting to create %v", r)
	u, err := r.createURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	req, err := r.marshal(c)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(req), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	o, err := dcl.ResponseBodyAsJSON(resp)
	if err != nil {
		return fmt.Errorf("error decoding response body into JSON: %w", err)
	}
	op.response = o

	if _, err := c.GetLogMetric(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getLogMetricRaw(ctx context.Context, r *LogMetric) ([]byte, error) {

	u, err := r.getURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "GET", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return nil, err
	}
	defer resp.Response.Body.Close()
	b, err := ioutil.ReadAll(resp.Response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) logMetricDiffsForRawDesired(ctx context.Context, rawDesired *LogMetric, opts ...dcl.ApplyOption) (initial, desired *LogMetric, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *LogMetric
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*LogMetric); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected LogMetric, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetLogMetric(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a LogMetric resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve LogMetric resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that LogMetric resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeLogMetricDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for LogMetric: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for LogMetric: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractLogMetricFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeLogMetricInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for LogMetric: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeLogMetricDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for LogMetric: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffLogMetric(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeLogMetricInitialState(rawInitial, rawDesired *LogMetric) (*LogMetric, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeLogMetricDesiredState(rawDesired, rawInitial *LogMetric, opts ...dcl.ApplyOption) (*LogMetric, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.MetricDescriptor = canonicalizeLogMetricMetricDescriptor(rawDesired.MetricDescriptor, nil, opts...)
		rawDesired.BucketOptions = canonicalizeLogMetricBucketOptions(rawDesired.BucketOptions, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &LogMetric{}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.StringCanonicalize(rawDesired.Filter, rawInitial.Filter) {
		canonicalDesired.Filter = rawInitial.Filter
	} else {
		canonicalDesired.Filter = rawDesired.Filter
	}
	if dcl.BoolCanonicalize(rawDesired.Disabled, rawInitial.Disabled) {
		canonicalDesired.Disabled = rawInitial.Disabled
	} else {
		canonicalDesired.Disabled = rawDesired.Disabled
	}
	canonicalDesired.MetricDescriptor = canonicalizeLogMetricMetricDescriptor(rawDesired.MetricDescriptor, rawInitial.MetricDescriptor, opts...)
	if dcl.StringCanonicalize(rawDesired.ValueExtractor, rawInitial.ValueExtractor) {
		canonicalDesired.ValueExtractor = rawInitial.ValueExtractor
	} else {
		canonicalDesired.ValueExtractor = rawDesired.ValueExtractor
	}
	if dcl.IsZeroValue(rawDesired.LabelExtractors) || (dcl.IsEmptyValueIndirect(rawDesired.LabelExtractors) && dcl.IsEmptyValueIndirect(rawInitial.LabelExtractors)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.LabelExtractors = rawInitial.LabelExtractors
	} else {
		canonicalDesired.LabelExtractors = rawDesired.LabelExtractors
	}
	canonicalDesired.BucketOptions = canonicalizeLogMetricBucketOptions(rawDesired.BucketOptions, rawInitial.BucketOptions, opts...)
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeLogMetricNewState(c *Client, rawNew, rawDesired *LogMetric) (*LogMetric, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Filter) && dcl.IsEmptyValueIndirect(rawDesired.Filter) {
		rawNew.Filter = rawDesired.Filter
	} else {
		if dcl.StringCanonicalize(rawDesired.Filter, rawNew.Filter) {
			rawNew.Filter = rawDesired.Filter
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Disabled) && dcl.IsEmptyValueIndirect(rawDesired.Disabled) {
		rawNew.Disabled = rawDesired.Disabled
	} else {
		if dcl.BoolCanonicalize(rawDesired.Disabled, rawNew.Disabled) {
			rawNew.Disabled = rawDesired.Disabled
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.MetricDescriptor) && dcl.IsEmptyValueIndirect(rawDesired.MetricDescriptor) {
		rawNew.MetricDescriptor = rawDesired.MetricDescriptor
	} else {
		rawNew.MetricDescriptor = canonicalizeNewLogMetricMetricDescriptor(c, rawDesired.MetricDescriptor, rawNew.MetricDescriptor)
	}

	if dcl.IsEmptyValueIndirect(rawNew.ValueExtractor) && dcl.IsEmptyValueIndirect(rawDesired.ValueExtractor) {
		rawNew.ValueExtractor = rawDesired.ValueExtractor
	} else {
		if dcl.StringCanonicalize(rawDesired.ValueExtractor, rawNew.ValueExtractor) {
			rawNew.ValueExtractor = rawDesired.ValueExtractor
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.LabelExtractors) && dcl.IsEmptyValueIndirect(rawDesired.LabelExtractors) {
		rawNew.LabelExtractors = rawDesired.LabelExtractors
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.BucketOptions) && dcl.IsEmptyValueIndirect(rawDesired.BucketOptions) {
		rawNew.BucketOptions = rawDesired.BucketOptions
	} else {
		rawNew.BucketOptions = canonicalizeNewLogMetricBucketOptions(c, rawDesired.BucketOptions, rawNew.BucketOptions)
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeLogMetricMetricDescriptor(des, initial *LogMetricMetricDescriptor, opts ...dcl.ApplyOption) *LogMetricMetricDescriptor {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LogMetricMetricDescriptor{}

	cDes.Labels = canonicalizeLogMetricMetricDescriptorLabelsSlice(des.Labels, initial.Labels, opts...)
	if dcl.IsZeroValue(des.MetricKind) || (dcl.IsEmptyValueIndirect(des.MetricKind) && dcl.IsEmptyValueIndirect(initial.MetricKind)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MetricKind = initial.MetricKind
	} else {
		cDes.MetricKind = des.MetricKind
	}
	if canonicalizeLogMetricMetricDescriptorValueType(des.ValueType, initial.ValueType) || dcl.IsZeroValue(des.ValueType) {
		cDes.ValueType = initial.ValueType
	} else {
		cDes.ValueType = des.ValueType
	}
	if dcl.StringCanonicalize(des.Unit, initial.Unit) || dcl.IsZeroValue(des.Unit) {
		cDes.Unit = initial.Unit
	} else {
		cDes.Unit = des.Unit
	}
	if dcl.StringCanonicalize(des.DisplayName, initial.DisplayName) || dcl.IsZeroValue(des.DisplayName) {
		cDes.DisplayName = initial.DisplayName
	} else {
		cDes.DisplayName = des.DisplayName
	}
	cDes.Metadata = canonicalizeLogMetricMetricDescriptorMetadata(des.Metadata, initial.Metadata, opts...)
	if dcl.IsZeroValue(des.LaunchStage) || (dcl.IsEmptyValueIndirect(des.LaunchStage) && dcl.IsEmptyValueIndirect(initial.LaunchStage)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.LaunchStage = initial.LaunchStage
	} else {
		cDes.LaunchStage = des.LaunchStage
	}

	return cDes
}

func canonicalizeLogMetricMetricDescriptorSlice(des, initial []LogMetricMetricDescriptor, opts ...dcl.ApplyOption) []LogMetricMetricDescriptor {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LogMetricMetricDescriptor, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLogMetricMetricDescriptor(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LogMetricMetricDescriptor, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLogMetricMetricDescriptor(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLogMetricMetricDescriptor(c *Client, des, nw *LogMetricMetricDescriptor) *LogMetricMetricDescriptor {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LogMetricMetricDescriptor while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Type, nw.Type) {
		nw.Type = des.Type
	}
	nw.Labels = canonicalizeNewLogMetricMetricDescriptorLabelsSet(c, des.Labels, nw.Labels)
	if canonicalizeLogMetricMetricDescriptorValueType(des.ValueType, nw.ValueType) {
		nw.ValueType = des.ValueType
	}
	if dcl.StringCanonicalize(des.Unit, nw.Unit) {
		nw.Unit = des.Unit
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}
	if dcl.StringCanonicalize(des.DisplayName, nw.DisplayName) {
		nw.DisplayName = des.DisplayName
	}
	nw.Metadata = des.Metadata
	nw.LaunchStage = des.LaunchStage
	if dcl.StringArrayCanonicalize(des.MonitoredResourceTypes, nw.MonitoredResourceTypes) {
		nw.MonitoredResourceTypes = des.MonitoredResourceTypes
	}

	return nw
}

func canonicalizeNewLogMetricMetricDescriptorSet(c *Client, des, nw []LogMetricMetricDescriptor) []LogMetricMetricDescriptor {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LogMetricMetricDescriptor
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLogMetricMetricDescriptorNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLogMetricMetricDescriptor(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLogMetricMetricDescriptorSlice(c *Client, des, nw []LogMetricMetricDescriptor) []LogMetricMetricDescriptor {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LogMetricMetricDescriptor
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLogMetricMetricDescriptor(c, &d, &n))
	}

	return items
}

func canonicalizeLogMetricMetricDescriptorLabels(des, initial *LogMetricMetricDescriptorLabels, opts ...dcl.ApplyOption) *LogMetricMetricDescriptorLabels {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LogMetricMetricDescriptorLabels{}

	if dcl.StringCanonicalize(des.Key, initial.Key) || dcl.IsZeroValue(des.Key) {
		cDes.Key = initial.Key
	} else {
		cDes.Key = des.Key
	}
	if canonicalizeLogMetricMetricDescriptorLabelsValueType(des.ValueType, initial.ValueType) || dcl.IsZeroValue(des.ValueType) {
		cDes.ValueType = initial.ValueType
	} else {
		cDes.ValueType = des.ValueType
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}

	return cDes
}

func canonicalizeLogMetricMetricDescriptorLabelsSlice(des, initial []LogMetricMetricDescriptorLabels, opts ...dcl.ApplyOption) []LogMetricMetricDescriptorLabels {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LogMetricMetricDescriptorLabels, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLogMetricMetricDescriptorLabels(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LogMetricMetricDescriptorLabels, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLogMetricMetricDescriptorLabels(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLogMetricMetricDescriptorLabels(c *Client, des, nw *LogMetricMetricDescriptorLabels) *LogMetricMetricDescriptorLabels {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LogMetricMetricDescriptorLabels while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}
	if canonicalizeLogMetricMetricDescriptorLabelsValueType(des.ValueType, nw.ValueType) {
		nw.ValueType = des.ValueType
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}

	return nw
}

func canonicalizeNewLogMetricMetricDescriptorLabelsSet(c *Client, des, nw []LogMetricMetricDescriptorLabels) []LogMetricMetricDescriptorLabels {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LogMetricMetricDescriptorLabels
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLogMetricMetricDescriptorLabelsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLogMetricMetricDescriptorLabels(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLogMetricMetricDescriptorLabelsSlice(c *Client, des, nw []LogMetricMetricDescriptorLabels) []LogMetricMetricDescriptorLabels {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LogMetricMetricDescriptorLabels
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLogMetricMetricDescriptorLabels(c, &d, &n))
	}

	return items
}

func canonicalizeLogMetricMetricDescriptorMetadata(des, initial *LogMetricMetricDescriptorMetadata, opts ...dcl.ApplyOption) *LogMetricMetricDescriptorMetadata {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LogMetricMetricDescriptorMetadata{}

	if dcl.StringCanonicalize(des.SamplePeriod, initial.SamplePeriod) || dcl.IsZeroValue(des.SamplePeriod) {
		cDes.SamplePeriod = initial.SamplePeriod
	} else {
		cDes.SamplePeriod = des.SamplePeriod
	}
	if dcl.StringCanonicalize(des.IngestDelay, initial.IngestDelay) || dcl.IsZeroValue(des.IngestDelay) {
		cDes.IngestDelay = initial.IngestDelay
	} else {
		cDes.IngestDelay = des.IngestDelay
	}

	return cDes
}

func canonicalizeLogMetricMetricDescriptorMetadataSlice(des, initial []LogMetricMetricDescriptorMetadata, opts ...dcl.ApplyOption) []LogMetricMetricDescriptorMetadata {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LogMetricMetricDescriptorMetadata, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLogMetricMetricDescriptorMetadata(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LogMetricMetricDescriptorMetadata, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLogMetricMetricDescriptorMetadata(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLogMetricMetricDescriptorMetadata(c *Client, des, nw *LogMetricMetricDescriptorMetadata) *LogMetricMetricDescriptorMetadata {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LogMetricMetricDescriptorMetadata while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.SamplePeriod, nw.SamplePeriod) {
		nw.SamplePeriod = des.SamplePeriod
	}
	if dcl.StringCanonicalize(des.IngestDelay, nw.IngestDelay) {
		nw.IngestDelay = des.IngestDelay
	}

	return nw
}

func canonicalizeNewLogMetricMetricDescriptorMetadataSet(c *Client, des, nw []LogMetricMetricDescriptorMetadata) []LogMetricMetricDescriptorMetadata {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LogMetricMetricDescriptorMetadata
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLogMetricMetricDescriptorMetadataNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLogMetricMetricDescriptorMetadata(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLogMetricMetricDescriptorMetadataSlice(c *Client, des, nw []LogMetricMetricDescriptorMetadata) []LogMetricMetricDescriptorMetadata {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LogMetricMetricDescriptorMetadata
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLogMetricMetricDescriptorMetadata(c, &d, &n))
	}

	return items
}

func canonicalizeLogMetricBucketOptions(des, initial *LogMetricBucketOptions, opts ...dcl.ApplyOption) *LogMetricBucketOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.LinearBuckets != nil || (initial != nil && initial.LinearBuckets != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.ExponentialBuckets, des.ExplicitBuckets) {
			des.LinearBuckets = nil
			if initial != nil {
				initial.LinearBuckets = nil
			}
		}
	}

	if des.ExponentialBuckets != nil || (initial != nil && initial.ExponentialBuckets != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.LinearBuckets, des.ExplicitBuckets) {
			des.ExponentialBuckets = nil
			if initial != nil {
				initial.ExponentialBuckets = nil
			}
		}
	}

	if des.ExplicitBuckets != nil || (initial != nil && initial.ExplicitBuckets != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.LinearBuckets, des.ExponentialBuckets) {
			des.ExplicitBuckets = nil
			if initial != nil {
				initial.ExplicitBuckets = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &LogMetricBucketOptions{}

	cDes.LinearBuckets = canonicalizeLogMetricBucketOptionsLinearBuckets(des.LinearBuckets, initial.LinearBuckets, opts...)
	cDes.ExponentialBuckets = canonicalizeLogMetricBucketOptionsExponentialBuckets(des.ExponentialBuckets, initial.ExponentialBuckets, opts...)
	cDes.ExplicitBuckets = canonicalizeLogMetricBucketOptionsExplicitBuckets(des.ExplicitBuckets, initial.ExplicitBuckets, opts...)

	return cDes
}

func canonicalizeLogMetricBucketOptionsSlice(des, initial []LogMetricBucketOptions, opts ...dcl.ApplyOption) []LogMetricBucketOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LogMetricBucketOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLogMetricBucketOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LogMetricBucketOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLogMetricBucketOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLogMetricBucketOptions(c *Client, des, nw *LogMetricBucketOptions) *LogMetricBucketOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LogMetricBucketOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.LinearBuckets = canonicalizeNewLogMetricBucketOptionsLinearBuckets(c, des.LinearBuckets, nw.LinearBuckets)
	nw.ExponentialBuckets = canonicalizeNewLogMetricBucketOptionsExponentialBuckets(c, des.ExponentialBuckets, nw.ExponentialBuckets)
	nw.ExplicitBuckets = canonicalizeNewLogMetricBucketOptionsExplicitBuckets(c, des.ExplicitBuckets, nw.ExplicitBuckets)

	return nw
}

func canonicalizeNewLogMetricBucketOptionsSet(c *Client, des, nw []LogMetricBucketOptions) []LogMetricBucketOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LogMetricBucketOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLogMetricBucketOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLogMetricBucketOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLogMetricBucketOptionsSlice(c *Client, des, nw []LogMetricBucketOptions) []LogMetricBucketOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LogMetricBucketOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLogMetricBucketOptions(c, &d, &n))
	}

	return items
}

func canonicalizeLogMetricBucketOptionsLinearBuckets(des, initial *LogMetricBucketOptionsLinearBuckets, opts ...dcl.ApplyOption) *LogMetricBucketOptionsLinearBuckets {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LogMetricBucketOptionsLinearBuckets{}

	if dcl.IsZeroValue(des.NumFiniteBuckets) || (dcl.IsEmptyValueIndirect(des.NumFiniteBuckets) && dcl.IsEmptyValueIndirect(initial.NumFiniteBuckets)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumFiniteBuckets = initial.NumFiniteBuckets
	} else {
		cDes.NumFiniteBuckets = des.NumFiniteBuckets
	}
	if dcl.IsZeroValue(des.Width) || (dcl.IsEmptyValueIndirect(des.Width) && dcl.IsEmptyValueIndirect(initial.Width)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Width = initial.Width
	} else {
		cDes.Width = des.Width
	}
	if dcl.IsZeroValue(des.Offset) || (dcl.IsEmptyValueIndirect(des.Offset) && dcl.IsEmptyValueIndirect(initial.Offset)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Offset = initial.Offset
	} else {
		cDes.Offset = des.Offset
	}

	return cDes
}

func canonicalizeLogMetricBucketOptionsLinearBucketsSlice(des, initial []LogMetricBucketOptionsLinearBuckets, opts ...dcl.ApplyOption) []LogMetricBucketOptionsLinearBuckets {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LogMetricBucketOptionsLinearBuckets, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLogMetricBucketOptionsLinearBuckets(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LogMetricBucketOptionsLinearBuckets, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLogMetricBucketOptionsLinearBuckets(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLogMetricBucketOptionsLinearBuckets(c *Client, des, nw *LogMetricBucketOptionsLinearBuckets) *LogMetricBucketOptionsLinearBuckets {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LogMetricBucketOptionsLinearBuckets while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewLogMetricBucketOptionsLinearBucketsSet(c *Client, des, nw []LogMetricBucketOptionsLinearBuckets) []LogMetricBucketOptionsLinearBuckets {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LogMetricBucketOptionsLinearBuckets
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLogMetricBucketOptionsLinearBucketsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLogMetricBucketOptionsLinearBuckets(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLogMetricBucketOptionsLinearBucketsSlice(c *Client, des, nw []LogMetricBucketOptionsLinearBuckets) []LogMetricBucketOptionsLinearBuckets {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LogMetricBucketOptionsLinearBuckets
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLogMetricBucketOptionsLinearBuckets(c, &d, &n))
	}

	return items
}

func canonicalizeLogMetricBucketOptionsExponentialBuckets(des, initial *LogMetricBucketOptionsExponentialBuckets, opts ...dcl.ApplyOption) *LogMetricBucketOptionsExponentialBuckets {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LogMetricBucketOptionsExponentialBuckets{}

	if dcl.IsZeroValue(des.NumFiniteBuckets) || (dcl.IsEmptyValueIndirect(des.NumFiniteBuckets) && dcl.IsEmptyValueIndirect(initial.NumFiniteBuckets)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumFiniteBuckets = initial.NumFiniteBuckets
	} else {
		cDes.NumFiniteBuckets = des.NumFiniteBuckets
	}
	if dcl.IsZeroValue(des.GrowthFactor) || (dcl.IsEmptyValueIndirect(des.GrowthFactor) && dcl.IsEmptyValueIndirect(initial.GrowthFactor)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.GrowthFactor = initial.GrowthFactor
	} else {
		cDes.GrowthFactor = des.GrowthFactor
	}
	if dcl.IsZeroValue(des.Scale) || (dcl.IsEmptyValueIndirect(des.Scale) && dcl.IsEmptyValueIndirect(initial.Scale)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Scale = initial.Scale
	} else {
		cDes.Scale = des.Scale
	}

	return cDes
}

func canonicalizeLogMetricBucketOptionsExponentialBucketsSlice(des, initial []LogMetricBucketOptionsExponentialBuckets, opts ...dcl.ApplyOption) []LogMetricBucketOptionsExponentialBuckets {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LogMetricBucketOptionsExponentialBuckets, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLogMetricBucketOptionsExponentialBuckets(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LogMetricBucketOptionsExponentialBuckets, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLogMetricBucketOptionsExponentialBuckets(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLogMetricBucketOptionsExponentialBuckets(c *Client, des, nw *LogMetricBucketOptionsExponentialBuckets) *LogMetricBucketOptionsExponentialBuckets {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LogMetricBucketOptionsExponentialBuckets while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewLogMetricBucketOptionsExponentialBucketsSet(c *Client, des, nw []LogMetricBucketOptionsExponentialBuckets) []LogMetricBucketOptionsExponentialBuckets {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LogMetricBucketOptionsExponentialBuckets
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLogMetricBucketOptionsExponentialBucketsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLogMetricBucketOptionsExponentialBuckets(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLogMetricBucketOptionsExponentialBucketsSlice(c *Client, des, nw []LogMetricBucketOptionsExponentialBuckets) []LogMetricBucketOptionsExponentialBuckets {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LogMetricBucketOptionsExponentialBuckets
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLogMetricBucketOptionsExponentialBuckets(c, &d, &n))
	}

	return items
}

func canonicalizeLogMetricBucketOptionsExplicitBuckets(des, initial *LogMetricBucketOptionsExplicitBuckets, opts ...dcl.ApplyOption) *LogMetricBucketOptionsExplicitBuckets {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &LogMetricBucketOptionsExplicitBuckets{}

	if dcl.IsZeroValue(des.Bounds) || (dcl.IsEmptyValueIndirect(des.Bounds) && dcl.IsEmptyValueIndirect(initial.Bounds)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Bounds = initial.Bounds
	} else {
		cDes.Bounds = des.Bounds
	}

	return cDes
}

func canonicalizeLogMetricBucketOptionsExplicitBucketsSlice(des, initial []LogMetricBucketOptionsExplicitBuckets, opts ...dcl.ApplyOption) []LogMetricBucketOptionsExplicitBuckets {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]LogMetricBucketOptionsExplicitBuckets, 0, len(des))
		for _, d := range des {
			cd := canonicalizeLogMetricBucketOptionsExplicitBuckets(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]LogMetricBucketOptionsExplicitBuckets, 0, len(des))
	for i, d := range des {
		cd := canonicalizeLogMetricBucketOptionsExplicitBuckets(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewLogMetricBucketOptionsExplicitBuckets(c *Client, des, nw *LogMetricBucketOptionsExplicitBuckets) *LogMetricBucketOptionsExplicitBuckets {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for LogMetricBucketOptionsExplicitBuckets while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewLogMetricBucketOptionsExplicitBucketsSet(c *Client, des, nw []LogMetricBucketOptionsExplicitBuckets) []LogMetricBucketOptionsExplicitBuckets {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []LogMetricBucketOptionsExplicitBuckets
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareLogMetricBucketOptionsExplicitBucketsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewLogMetricBucketOptionsExplicitBuckets(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewLogMetricBucketOptionsExplicitBucketsSlice(c *Client, des, nw []LogMetricBucketOptionsExplicitBuckets) []LogMetricBucketOptionsExplicitBuckets {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []LogMetricBucketOptionsExplicitBuckets
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewLogMetricBucketOptionsExplicitBuckets(c, &d, &n))
	}

	return items
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffLogMetric(c *Client, desired, actual *LogMetric, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Filter, actual.Filter, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Filter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Disabled, actual.Disabled, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Disabled")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetricDescriptor, actual.MetricDescriptor, dcl.DiffInfo{ObjectFunction: compareLogMetricMetricDescriptorNewStyle, EmptyObject: EmptyLogMetricMetricDescriptor, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetricDescriptor")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ValueExtractor, actual.ValueExtractor, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("ValueExtractor")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LabelExtractors, actual.LabelExtractors, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("LabelExtractors")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BucketOptions, actual.BucketOptions, dcl.DiffInfo{ObjectFunction: compareLogMetricBucketOptionsNewStyle, EmptyObject: EmptyLogMetricBucketOptions, OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("BucketOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreateTime, actual.CreateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Project")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if len(newDiffs) > 0 {
		c.Config.Logger.Infof("Diff function found diffs: %v", newDiffs)
	}
	return newDiffs, nil
}
func compareLogMetricMetricDescriptorNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LogMetricMetricDescriptor)
	if !ok {
		desiredNotPointer, ok := d.(LogMetricMetricDescriptor)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricMetricDescriptor or *LogMetricMetricDescriptor", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LogMetricMetricDescriptor)
	if !ok {
		actualNotPointer, ok := a.(LogMetricMetricDescriptor)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricMetricDescriptor", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{Type: "Set", ObjectFunction: compareLogMetricMetricDescriptorLabelsNewStyle, EmptyObject: EmptyLogMetricMetricDescriptorLabels, OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetricKind, actual.MetricKind, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetricKind")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ValueType, actual.ValueType, dcl.DiffInfo{Type: "EnumType", CustomDiff: canonicalizeLogMetricMetricDescriptorValueType, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ValueType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Unit, actual.Unit, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Unit")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metadata, actual.Metadata, dcl.DiffInfo{Ignore: true, ObjectFunction: compareLogMetricMetricDescriptorMetadataNewStyle, EmptyObject: EmptyLogMetricMetricDescriptorMetadata, OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Metadata")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LaunchStage, actual.LaunchStage, dcl.DiffInfo{Ignore: true, Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("LaunchStage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MonitoredResourceTypes, actual.MonitoredResourceTypes, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MonitoredResourceTypes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareLogMetricMetricDescriptorLabelsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LogMetricMetricDescriptorLabels)
	if !ok {
		desiredNotPointer, ok := d.(LogMetricMetricDescriptorLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricMetricDescriptorLabels or *LogMetricMetricDescriptorLabels", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LogMetricMetricDescriptorLabels)
	if !ok {
		actualNotPointer, ok := a.(LogMetricMetricDescriptorLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricMetricDescriptorLabels", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Key, actual.Key, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Key")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ValueType, actual.ValueType, dcl.DiffInfo{Type: "EnumType", CustomDiff: canonicalizeLogMetricMetricDescriptorLabelsValueType, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ValueType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareLogMetricMetricDescriptorMetadataNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LogMetricMetricDescriptorMetadata)
	if !ok {
		desiredNotPointer, ok := d.(LogMetricMetricDescriptorMetadata)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricMetricDescriptorMetadata or *LogMetricMetricDescriptorMetadata", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LogMetricMetricDescriptorMetadata)
	if !ok {
		actualNotPointer, ok := a.(LogMetricMetricDescriptorMetadata)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricMetricDescriptorMetadata", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SamplePeriod, actual.SamplePeriod, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("SamplePeriod")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IngestDelay, actual.IngestDelay, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("IngestDelay")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareLogMetricBucketOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LogMetricBucketOptions)
	if !ok {
		desiredNotPointer, ok := d.(LogMetricBucketOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricBucketOptions or *LogMetricBucketOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LogMetricBucketOptions)
	if !ok {
		actualNotPointer, ok := a.(LogMetricBucketOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricBucketOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LinearBuckets, actual.LinearBuckets, dcl.DiffInfo{ObjectFunction: compareLogMetricBucketOptionsLinearBucketsNewStyle, EmptyObject: EmptyLogMetricBucketOptionsLinearBuckets, OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("LinearBuckets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExponentialBuckets, actual.ExponentialBuckets, dcl.DiffInfo{ObjectFunction: compareLogMetricBucketOptionsExponentialBucketsNewStyle, EmptyObject: EmptyLogMetricBucketOptionsExponentialBuckets, OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("ExponentialBuckets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExplicitBuckets, actual.ExplicitBuckets, dcl.DiffInfo{ObjectFunction: compareLogMetricBucketOptionsExplicitBucketsNewStyle, EmptyObject: EmptyLogMetricBucketOptionsExplicitBuckets, OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("ExplicitBuckets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareLogMetricBucketOptionsLinearBucketsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LogMetricBucketOptionsLinearBuckets)
	if !ok {
		desiredNotPointer, ok := d.(LogMetricBucketOptionsLinearBuckets)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricBucketOptionsLinearBuckets or *LogMetricBucketOptionsLinearBuckets", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LogMetricBucketOptionsLinearBuckets)
	if !ok {
		actualNotPointer, ok := a.(LogMetricBucketOptionsLinearBuckets)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricBucketOptionsLinearBuckets", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumFiniteBuckets, actual.NumFiniteBuckets, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("NumFiniteBuckets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Width, actual.Width, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Width")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Offset, actual.Offset, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Offset")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareLogMetricBucketOptionsExponentialBucketsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LogMetricBucketOptionsExponentialBuckets)
	if !ok {
		desiredNotPointer, ok := d.(LogMetricBucketOptionsExponentialBuckets)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricBucketOptionsExponentialBuckets or *LogMetricBucketOptionsExponentialBuckets", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LogMetricBucketOptionsExponentialBuckets)
	if !ok {
		actualNotPointer, ok := a.(LogMetricBucketOptionsExponentialBuckets)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricBucketOptionsExponentialBuckets", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumFiniteBuckets, actual.NumFiniteBuckets, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("NumFiniteBuckets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GrowthFactor, actual.GrowthFactor, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("GrowthFactor")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Scale, actual.Scale, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Scale")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareLogMetricBucketOptionsExplicitBucketsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*LogMetricBucketOptionsExplicitBuckets)
	if !ok {
		desiredNotPointer, ok := d.(LogMetricBucketOptionsExplicitBuckets)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricBucketOptionsExplicitBuckets or *LogMetricBucketOptionsExplicitBuckets", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*LogMetricBucketOptionsExplicitBuckets)
	if !ok {
		actualNotPointer, ok := a.(LogMetricBucketOptionsExplicitBuckets)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a LogMetricBucketOptionsExplicitBuckets", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bounds, actual.Bounds, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogMetricUpdateOperation")}, fn.AddNest("Bounds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

// urlNormalized returns a copy of the resource struct with values normalized
// for URL substitutions. For instance, it converts long-form self-links to
// short-form so they can be substituted in.
func (r *LogMetric) urlNormalized() *LogMetric {
	normalized := dcl.Copy(*r).(LogMetric)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Filter = dcl.SelfLinkToName(r.Filter)
	normalized.ValueExtractor = dcl.SelfLinkToName(r.ValueExtractor)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *LogMetric) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "update" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/metrics/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the LogMetric resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *LogMetric) marshal(c *Client) ([]byte, error) {
	m, err := expandLogMetric(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling LogMetric: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalLogMetric decodes JSON responses into the LogMetric resource schema.
func unmarshalLogMetric(b []byte, c *Client, res *LogMetric) (*LogMetric, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapLogMetric(m, c, res)
}

func unmarshalMapLogMetric(m map[string]interface{}, c *Client, res *LogMetric) (*LogMetric, error) {

	flattened := flattenLogMetric(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandLogMetric expands LogMetric into a JSON request object.
func expandLogMetric(c *Client, f *LogMetric) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Filter; dcl.ValueShouldBeSent(v) {
		m["filter"] = v
	}
	if v := f.Disabled; dcl.ValueShouldBeSent(v) {
		m["disabled"] = v
	}
	if v, err := expandLogMetricMetricDescriptor(c, f.MetricDescriptor, res); err != nil {
		return nil, fmt.Errorf("error expanding MetricDescriptor into metricDescriptor: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metricDescriptor"] = v
	}
	if v := f.ValueExtractor; dcl.ValueShouldBeSent(v) {
		m["valueExtractor"] = v
	}
	if v := f.LabelExtractors; dcl.ValueShouldBeSent(v) {
		m["labelExtractors"] = v
	}
	if v, err := expandLogMetricBucketOptions(c, f.BucketOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding BucketOptions into bucketOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["bucketOptions"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenLogMetric flattens LogMetric from a JSON request object into the
// LogMetric type.
func flattenLogMetric(c *Client, i interface{}, res *LogMetric) *LogMetric {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &LogMetric{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Filter = dcl.FlattenString(m["filter"])
	resultRes.Disabled = dcl.FlattenBool(m["disabled"])
	resultRes.MetricDescriptor = flattenLogMetricMetricDescriptor(c, m["metricDescriptor"], res)
	resultRes.ValueExtractor = dcl.FlattenString(m["valueExtractor"])
	resultRes.LabelExtractors = dcl.FlattenKeyValuePairs(m["labelExtractors"])
	resultRes.BucketOptions = flattenLogMetricBucketOptions(c, m["bucketOptions"], res)
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// expandLogMetricMetricDescriptorMap expands the contents of LogMetricMetricDescriptor into a JSON
// request object.
func expandLogMetricMetricDescriptorMap(c *Client, f map[string]LogMetricMetricDescriptor, res *LogMetric) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLogMetricMetricDescriptor(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLogMetricMetricDescriptorSlice expands the contents of LogMetricMetricDescriptor into a JSON
// request object.
func expandLogMetricMetricDescriptorSlice(c *Client, f []LogMetricMetricDescriptor, res *LogMetric) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLogMetricMetricDescriptor(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLogMetricMetricDescriptorMap flattens the contents of LogMetricMetricDescriptor from a JSON
// response object.
func flattenLogMetricMetricDescriptorMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricMetricDescriptor {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricMetricDescriptor{}
	}

	if len(a) == 0 {
		return map[string]LogMetricMetricDescriptor{}
	}

	items := make(map[string]LogMetricMetricDescriptor)
	for k, item := range a {
		items[k] = *flattenLogMetricMetricDescriptor(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLogMetricMetricDescriptorSlice flattens the contents of LogMetricMetricDescriptor from a JSON
// response object.
func flattenLogMetricMetricDescriptorSlice(c *Client, i interface{}, res *LogMetric) []LogMetricMetricDescriptor {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricMetricDescriptor{}
	}

	if len(a) == 0 {
		return []LogMetricMetricDescriptor{}
	}

	items := make([]LogMetricMetricDescriptor, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricMetricDescriptor(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLogMetricMetricDescriptor expands an instance of LogMetricMetricDescriptor into a JSON
// request object.
func expandLogMetricMetricDescriptor(c *Client, f *LogMetricMetricDescriptor, res *LogMetric) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandLogMetricMetricDescriptorLabelsSlice(c, f.Labels, res); err != nil {
		return nil, fmt.Errorf("error expanding Labels into labels: %w", err)
	} else if v != nil {
		m["labels"] = v
	}
	if v := f.MetricKind; !dcl.IsEmptyValueIndirect(v) {
		m["metricKind"] = v
	}
	if v := f.ValueType; !dcl.IsEmptyValueIndirect(v) {
		m["valueType"] = v
	}
	if v := f.Unit; !dcl.IsEmptyValueIndirect(v) {
		m["unit"] = v
	}
	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		m["displayName"] = v
	}
	if v, err := expandLogMetricMetricDescriptorMetadata(c, f.Metadata, res); err != nil {
		return nil, fmt.Errorf("error expanding Metadata into metadata: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metadata"] = v
	}
	if v := f.LaunchStage; !dcl.IsEmptyValueIndirect(v) {
		m["launchStage"] = v
	}

	return m, nil
}

// flattenLogMetricMetricDescriptor flattens an instance of LogMetricMetricDescriptor from a JSON
// response object.
func flattenLogMetricMetricDescriptor(c *Client, i interface{}, res *LogMetric) *LogMetricMetricDescriptor {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LogMetricMetricDescriptor{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLogMetricMetricDescriptor
	}
	r.Name = dcl.FlattenString(m["name"])
	r.Type = dcl.FlattenString(m["type"])
	r.Labels = flattenLogMetricMetricDescriptorLabelsSlice(c, m["labels"], res)
	r.MetricKind = flattenLogMetricMetricDescriptorMetricKindEnum(m["metricKind"])
	r.ValueType = flattenLogMetricMetricDescriptorValueTypeEnum(m["valueType"])
	r.Unit = dcl.FlattenString(m["unit"])
	r.Description = dcl.FlattenString(m["description"])
	r.DisplayName = dcl.FlattenString(m["displayName"])
	r.Metadata = flattenLogMetricMetricDescriptorMetadata(c, m["metadata"], res)
	r.LaunchStage = flattenLogMetricMetricDescriptorLaunchStageEnum(m["launchStage"])
	r.MonitoredResourceTypes = dcl.FlattenStringSlice(m["monitoredResourceTypes"])

	return r
}

// expandLogMetricMetricDescriptorLabelsMap expands the contents of LogMetricMetricDescriptorLabels into a JSON
// request object.
func expandLogMetricMetricDescriptorLabelsMap(c *Client, f map[string]LogMetricMetricDescriptorLabels, res *LogMetric) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLogMetricMetricDescriptorLabels(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLogMetricMetricDescriptorLabelsSlice expands the contents of LogMetricMetricDescriptorLabels into a JSON
// request object.
func expandLogMetricMetricDescriptorLabelsSlice(c *Client, f []LogMetricMetricDescriptorLabels, res *LogMetric) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLogMetricMetricDescriptorLabels(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLogMetricMetricDescriptorLabelsMap flattens the contents of LogMetricMetricDescriptorLabels from a JSON
// response object.
func flattenLogMetricMetricDescriptorLabelsMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricMetricDescriptorLabels {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricMetricDescriptorLabels{}
	}

	if len(a) == 0 {
		return map[string]LogMetricMetricDescriptorLabels{}
	}

	items := make(map[string]LogMetricMetricDescriptorLabels)
	for k, item := range a {
		items[k] = *flattenLogMetricMetricDescriptorLabels(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLogMetricMetricDescriptorLabelsSlice flattens the contents of LogMetricMetricDescriptorLabels from a JSON
// response object.
func flattenLogMetricMetricDescriptorLabelsSlice(c *Client, i interface{}, res *LogMetric) []LogMetricMetricDescriptorLabels {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricMetricDescriptorLabels{}
	}

	if len(a) == 0 {
		return []LogMetricMetricDescriptorLabels{}
	}

	items := make([]LogMetricMetricDescriptorLabels, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricMetricDescriptorLabels(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLogMetricMetricDescriptorLabels expands an instance of LogMetricMetricDescriptorLabels into a JSON
// request object.
func expandLogMetricMetricDescriptorLabels(c *Client, f *LogMetricMetricDescriptorLabels, res *LogMetric) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Key; !dcl.IsEmptyValueIndirect(v) {
		m["key"] = v
	}
	if v := f.ValueType; !dcl.IsEmptyValueIndirect(v) {
		m["valueType"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}

	return m, nil
}

// flattenLogMetricMetricDescriptorLabels flattens an instance of LogMetricMetricDescriptorLabels from a JSON
// response object.
func flattenLogMetricMetricDescriptorLabels(c *Client, i interface{}, res *LogMetric) *LogMetricMetricDescriptorLabels {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LogMetricMetricDescriptorLabels{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLogMetricMetricDescriptorLabels
	}
	r.Key = dcl.FlattenString(m["key"])
	r.ValueType = flattenLogMetricMetricDescriptorLabelsValueTypeEnum(m["valueType"])
	r.Description = dcl.FlattenString(m["description"])

	return r
}

// expandLogMetricMetricDescriptorMetadataMap expands the contents of LogMetricMetricDescriptorMetadata into a JSON
// request object.
func expandLogMetricMetricDescriptorMetadataMap(c *Client, f map[string]LogMetricMetricDescriptorMetadata, res *LogMetric) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLogMetricMetricDescriptorMetadata(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLogMetricMetricDescriptorMetadataSlice expands the contents of LogMetricMetricDescriptorMetadata into a JSON
// request object.
func expandLogMetricMetricDescriptorMetadataSlice(c *Client, f []LogMetricMetricDescriptorMetadata, res *LogMetric) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLogMetricMetricDescriptorMetadata(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLogMetricMetricDescriptorMetadataMap flattens the contents of LogMetricMetricDescriptorMetadata from a JSON
// response object.
func flattenLogMetricMetricDescriptorMetadataMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricMetricDescriptorMetadata {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricMetricDescriptorMetadata{}
	}

	if len(a) == 0 {
		return map[string]LogMetricMetricDescriptorMetadata{}
	}

	items := make(map[string]LogMetricMetricDescriptorMetadata)
	for k, item := range a {
		items[k] = *flattenLogMetricMetricDescriptorMetadata(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLogMetricMetricDescriptorMetadataSlice flattens the contents of LogMetricMetricDescriptorMetadata from a JSON
// response object.
func flattenLogMetricMetricDescriptorMetadataSlice(c *Client, i interface{}, res *LogMetric) []LogMetricMetricDescriptorMetadata {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricMetricDescriptorMetadata{}
	}

	if len(a) == 0 {
		return []LogMetricMetricDescriptorMetadata{}
	}

	items := make([]LogMetricMetricDescriptorMetadata, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricMetricDescriptorMetadata(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLogMetricMetricDescriptorMetadata expands an instance of LogMetricMetricDescriptorMetadata into a JSON
// request object.
func expandLogMetricMetricDescriptorMetadata(c *Client, f *LogMetricMetricDescriptorMetadata, res *LogMetric) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.SamplePeriod; !dcl.IsEmptyValueIndirect(v) {
		m["samplePeriod"] = v
	}
	if v := f.IngestDelay; !dcl.IsEmptyValueIndirect(v) {
		m["ingestDelay"] = v
	}

	return m, nil
}

// flattenLogMetricMetricDescriptorMetadata flattens an instance of LogMetricMetricDescriptorMetadata from a JSON
// response object.
func flattenLogMetricMetricDescriptorMetadata(c *Client, i interface{}, res *LogMetric) *LogMetricMetricDescriptorMetadata {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LogMetricMetricDescriptorMetadata{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLogMetricMetricDescriptorMetadata
	}
	r.SamplePeriod = dcl.FlattenString(m["samplePeriod"])
	r.IngestDelay = dcl.FlattenString(m["ingestDelay"])

	return r
}

// expandLogMetricBucketOptionsMap expands the contents of LogMetricBucketOptions into a JSON
// request object.
func expandLogMetricBucketOptionsMap(c *Client, f map[string]LogMetricBucketOptions, res *LogMetric) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLogMetricBucketOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLogMetricBucketOptionsSlice expands the contents of LogMetricBucketOptions into a JSON
// request object.
func expandLogMetricBucketOptionsSlice(c *Client, f []LogMetricBucketOptions, res *LogMetric) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLogMetricBucketOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLogMetricBucketOptionsMap flattens the contents of LogMetricBucketOptions from a JSON
// response object.
func flattenLogMetricBucketOptionsMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricBucketOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricBucketOptions{}
	}

	if len(a) == 0 {
		return map[string]LogMetricBucketOptions{}
	}

	items := make(map[string]LogMetricBucketOptions)
	for k, item := range a {
		items[k] = *flattenLogMetricBucketOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLogMetricBucketOptionsSlice flattens the contents of LogMetricBucketOptions from a JSON
// response object.
func flattenLogMetricBucketOptionsSlice(c *Client, i interface{}, res *LogMetric) []LogMetricBucketOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricBucketOptions{}
	}

	if len(a) == 0 {
		return []LogMetricBucketOptions{}
	}

	items := make([]LogMetricBucketOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricBucketOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLogMetricBucketOptions expands an instance of LogMetricBucketOptions into a JSON
// request object.
func expandLogMetricBucketOptions(c *Client, f *LogMetricBucketOptions, res *LogMetric) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandLogMetricBucketOptionsLinearBuckets(c, f.LinearBuckets, res); err != nil {
		return nil, fmt.Errorf("error expanding LinearBuckets into linearBuckets: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["linearBuckets"] = v
	}
	if v, err := expandLogMetricBucketOptionsExponentialBuckets(c, f.ExponentialBuckets, res); err != nil {
		return nil, fmt.Errorf("error expanding ExponentialBuckets into exponentialBuckets: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["exponentialBuckets"] = v
	}
	if v, err := expandLogMetricBucketOptionsExplicitBuckets(c, f.ExplicitBuckets, res); err != nil {
		return nil, fmt.Errorf("error expanding ExplicitBuckets into explicitBuckets: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["explicitBuckets"] = v
	}

	return m, nil
}

// flattenLogMetricBucketOptions flattens an instance of LogMetricBucketOptions from a JSON
// response object.
func flattenLogMetricBucketOptions(c *Client, i interface{}, res *LogMetric) *LogMetricBucketOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LogMetricBucketOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLogMetricBucketOptions
	}
	r.LinearBuckets = flattenLogMetricBucketOptionsLinearBuckets(c, m["linearBuckets"], res)
	r.ExponentialBuckets = flattenLogMetricBucketOptionsExponentialBuckets(c, m["exponentialBuckets"], res)
	r.ExplicitBuckets = flattenLogMetricBucketOptionsExplicitBuckets(c, m["explicitBuckets"], res)

	return r
}

// expandLogMetricBucketOptionsLinearBucketsMap expands the contents of LogMetricBucketOptionsLinearBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsLinearBucketsMap(c *Client, f map[string]LogMetricBucketOptionsLinearBuckets, res *LogMetric) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLogMetricBucketOptionsLinearBuckets(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLogMetricBucketOptionsLinearBucketsSlice expands the contents of LogMetricBucketOptionsLinearBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsLinearBucketsSlice(c *Client, f []LogMetricBucketOptionsLinearBuckets, res *LogMetric) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLogMetricBucketOptionsLinearBuckets(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLogMetricBucketOptionsLinearBucketsMap flattens the contents of LogMetricBucketOptionsLinearBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsLinearBucketsMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricBucketOptionsLinearBuckets {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricBucketOptionsLinearBuckets{}
	}

	if len(a) == 0 {
		return map[string]LogMetricBucketOptionsLinearBuckets{}
	}

	items := make(map[string]LogMetricBucketOptionsLinearBuckets)
	for k, item := range a {
		items[k] = *flattenLogMetricBucketOptionsLinearBuckets(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLogMetricBucketOptionsLinearBucketsSlice flattens the contents of LogMetricBucketOptionsLinearBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsLinearBucketsSlice(c *Client, i interface{}, res *LogMetric) []LogMetricBucketOptionsLinearBuckets {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricBucketOptionsLinearBuckets{}
	}

	if len(a) == 0 {
		return []LogMetricBucketOptionsLinearBuckets{}
	}

	items := make([]LogMetricBucketOptionsLinearBuckets, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricBucketOptionsLinearBuckets(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLogMetricBucketOptionsLinearBuckets expands an instance of LogMetricBucketOptionsLinearBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsLinearBuckets(c *Client, f *LogMetricBucketOptionsLinearBuckets, res *LogMetric) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NumFiniteBuckets; !dcl.IsEmptyValueIndirect(v) {
		m["numFiniteBuckets"] = v
	}
	if v := f.Width; !dcl.IsEmptyValueIndirect(v) {
		m["width"] = v
	}
	if v := f.Offset; !dcl.IsEmptyValueIndirect(v) {
		m["offset"] = v
	}

	return m, nil
}

// flattenLogMetricBucketOptionsLinearBuckets flattens an instance of LogMetricBucketOptionsLinearBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsLinearBuckets(c *Client, i interface{}, res *LogMetric) *LogMetricBucketOptionsLinearBuckets {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LogMetricBucketOptionsLinearBuckets{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLogMetricBucketOptionsLinearBuckets
	}
	r.NumFiniteBuckets = dcl.FlattenInteger(m["numFiniteBuckets"])
	r.Width = dcl.FlattenDouble(m["width"])
	r.Offset = dcl.FlattenDouble(m["offset"])

	return r
}

// expandLogMetricBucketOptionsExponentialBucketsMap expands the contents of LogMetricBucketOptionsExponentialBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsExponentialBucketsMap(c *Client, f map[string]LogMetricBucketOptionsExponentialBuckets, res *LogMetric) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLogMetricBucketOptionsExponentialBuckets(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLogMetricBucketOptionsExponentialBucketsSlice expands the contents of LogMetricBucketOptionsExponentialBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsExponentialBucketsSlice(c *Client, f []LogMetricBucketOptionsExponentialBuckets, res *LogMetric) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLogMetricBucketOptionsExponentialBuckets(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLogMetricBucketOptionsExponentialBucketsMap flattens the contents of LogMetricBucketOptionsExponentialBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsExponentialBucketsMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricBucketOptionsExponentialBuckets {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricBucketOptionsExponentialBuckets{}
	}

	if len(a) == 0 {
		return map[string]LogMetricBucketOptionsExponentialBuckets{}
	}

	items := make(map[string]LogMetricBucketOptionsExponentialBuckets)
	for k, item := range a {
		items[k] = *flattenLogMetricBucketOptionsExponentialBuckets(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLogMetricBucketOptionsExponentialBucketsSlice flattens the contents of LogMetricBucketOptionsExponentialBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsExponentialBucketsSlice(c *Client, i interface{}, res *LogMetric) []LogMetricBucketOptionsExponentialBuckets {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricBucketOptionsExponentialBuckets{}
	}

	if len(a) == 0 {
		return []LogMetricBucketOptionsExponentialBuckets{}
	}

	items := make([]LogMetricBucketOptionsExponentialBuckets, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricBucketOptionsExponentialBuckets(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLogMetricBucketOptionsExponentialBuckets expands an instance of LogMetricBucketOptionsExponentialBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsExponentialBuckets(c *Client, f *LogMetricBucketOptionsExponentialBuckets, res *LogMetric) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NumFiniteBuckets; !dcl.IsEmptyValueIndirect(v) {
		m["numFiniteBuckets"] = v
	}
	if v := f.GrowthFactor; !dcl.IsEmptyValueIndirect(v) {
		m["growthFactor"] = v
	}
	if v := f.Scale; !dcl.IsEmptyValueIndirect(v) {
		m["scale"] = v
	}

	return m, nil
}

// flattenLogMetricBucketOptionsExponentialBuckets flattens an instance of LogMetricBucketOptionsExponentialBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsExponentialBuckets(c *Client, i interface{}, res *LogMetric) *LogMetricBucketOptionsExponentialBuckets {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LogMetricBucketOptionsExponentialBuckets{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLogMetricBucketOptionsExponentialBuckets
	}
	r.NumFiniteBuckets = dcl.FlattenInteger(m["numFiniteBuckets"])
	r.GrowthFactor = dcl.FlattenDouble(m["growthFactor"])
	r.Scale = dcl.FlattenDouble(m["scale"])

	return r
}

// expandLogMetricBucketOptionsExplicitBucketsMap expands the contents of LogMetricBucketOptionsExplicitBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsExplicitBucketsMap(c *Client, f map[string]LogMetricBucketOptionsExplicitBuckets, res *LogMetric) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandLogMetricBucketOptionsExplicitBuckets(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandLogMetricBucketOptionsExplicitBucketsSlice expands the contents of LogMetricBucketOptionsExplicitBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsExplicitBucketsSlice(c *Client, f []LogMetricBucketOptionsExplicitBuckets, res *LogMetric) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandLogMetricBucketOptionsExplicitBuckets(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenLogMetricBucketOptionsExplicitBucketsMap flattens the contents of LogMetricBucketOptionsExplicitBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsExplicitBucketsMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricBucketOptionsExplicitBuckets {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricBucketOptionsExplicitBuckets{}
	}

	if len(a) == 0 {
		return map[string]LogMetricBucketOptionsExplicitBuckets{}
	}

	items := make(map[string]LogMetricBucketOptionsExplicitBuckets)
	for k, item := range a {
		items[k] = *flattenLogMetricBucketOptionsExplicitBuckets(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenLogMetricBucketOptionsExplicitBucketsSlice flattens the contents of LogMetricBucketOptionsExplicitBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsExplicitBucketsSlice(c *Client, i interface{}, res *LogMetric) []LogMetricBucketOptionsExplicitBuckets {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricBucketOptionsExplicitBuckets{}
	}

	if len(a) == 0 {
		return []LogMetricBucketOptionsExplicitBuckets{}
	}

	items := make([]LogMetricBucketOptionsExplicitBuckets, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricBucketOptionsExplicitBuckets(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandLogMetricBucketOptionsExplicitBuckets expands an instance of LogMetricBucketOptionsExplicitBuckets into a JSON
// request object.
func expandLogMetricBucketOptionsExplicitBuckets(c *Client, f *LogMetricBucketOptionsExplicitBuckets, res *LogMetric) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Bounds; v != nil {
		m["bounds"] = v
	}

	return m, nil
}

// flattenLogMetricBucketOptionsExplicitBuckets flattens an instance of LogMetricBucketOptionsExplicitBuckets from a JSON
// response object.
func flattenLogMetricBucketOptionsExplicitBuckets(c *Client, i interface{}, res *LogMetric) *LogMetricBucketOptionsExplicitBuckets {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &LogMetricBucketOptionsExplicitBuckets{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyLogMetricBucketOptionsExplicitBuckets
	}
	r.Bounds = dcl.FlattenFloatSlice(m["bounds"])

	return r
}

// flattenLogMetricMetricDescriptorLabelsValueTypeEnumMap flattens the contents of LogMetricMetricDescriptorLabelsValueTypeEnum from a JSON
// response object.
func flattenLogMetricMetricDescriptorLabelsValueTypeEnumMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricMetricDescriptorLabelsValueTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricMetricDescriptorLabelsValueTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]LogMetricMetricDescriptorLabelsValueTypeEnum{}
	}

	items := make(map[string]LogMetricMetricDescriptorLabelsValueTypeEnum)
	for k, item := range a {
		items[k] = *flattenLogMetricMetricDescriptorLabelsValueTypeEnum(item.(interface{}))
	}

	return items
}

// flattenLogMetricMetricDescriptorLabelsValueTypeEnumSlice flattens the contents of LogMetricMetricDescriptorLabelsValueTypeEnum from a JSON
// response object.
func flattenLogMetricMetricDescriptorLabelsValueTypeEnumSlice(c *Client, i interface{}, res *LogMetric) []LogMetricMetricDescriptorLabelsValueTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricMetricDescriptorLabelsValueTypeEnum{}
	}

	if len(a) == 0 {
		return []LogMetricMetricDescriptorLabelsValueTypeEnum{}
	}

	items := make([]LogMetricMetricDescriptorLabelsValueTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricMetricDescriptorLabelsValueTypeEnum(item.(interface{})))
	}

	return items
}

// flattenLogMetricMetricDescriptorLabelsValueTypeEnum asserts that an interface is a string, and returns a
// pointer to a *LogMetricMetricDescriptorLabelsValueTypeEnum with the same value as that string.
func flattenLogMetricMetricDescriptorLabelsValueTypeEnum(i interface{}) *LogMetricMetricDescriptorLabelsValueTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return LogMetricMetricDescriptorLabelsValueTypeEnumRef(s)
}

// flattenLogMetricMetricDescriptorMetricKindEnumMap flattens the contents of LogMetricMetricDescriptorMetricKindEnum from a JSON
// response object.
func flattenLogMetricMetricDescriptorMetricKindEnumMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricMetricDescriptorMetricKindEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricMetricDescriptorMetricKindEnum{}
	}

	if len(a) == 0 {
		return map[string]LogMetricMetricDescriptorMetricKindEnum{}
	}

	items := make(map[string]LogMetricMetricDescriptorMetricKindEnum)
	for k, item := range a {
		items[k] = *flattenLogMetricMetricDescriptorMetricKindEnum(item.(interface{}))
	}

	return items
}

// flattenLogMetricMetricDescriptorMetricKindEnumSlice flattens the contents of LogMetricMetricDescriptorMetricKindEnum from a JSON
// response object.
func flattenLogMetricMetricDescriptorMetricKindEnumSlice(c *Client, i interface{}, res *LogMetric) []LogMetricMetricDescriptorMetricKindEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricMetricDescriptorMetricKindEnum{}
	}

	if len(a) == 0 {
		return []LogMetricMetricDescriptorMetricKindEnum{}
	}

	items := make([]LogMetricMetricDescriptorMetricKindEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricMetricDescriptorMetricKindEnum(item.(interface{})))
	}

	return items
}

// flattenLogMetricMetricDescriptorMetricKindEnum asserts that an interface is a string, and returns a
// pointer to a *LogMetricMetricDescriptorMetricKindEnum with the same value as that string.
func flattenLogMetricMetricDescriptorMetricKindEnum(i interface{}) *LogMetricMetricDescriptorMetricKindEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return LogMetricMetricDescriptorMetricKindEnumRef(s)
}

// flattenLogMetricMetricDescriptorValueTypeEnumMap flattens the contents of LogMetricMetricDescriptorValueTypeEnum from a JSON
// response object.
func flattenLogMetricMetricDescriptorValueTypeEnumMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricMetricDescriptorValueTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricMetricDescriptorValueTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]LogMetricMetricDescriptorValueTypeEnum{}
	}

	items := make(map[string]LogMetricMetricDescriptorValueTypeEnum)
	for k, item := range a {
		items[k] = *flattenLogMetricMetricDescriptorValueTypeEnum(item.(interface{}))
	}

	return items
}

// flattenLogMetricMetricDescriptorValueTypeEnumSlice flattens the contents of LogMetricMetricDescriptorValueTypeEnum from a JSON
// response object.
func flattenLogMetricMetricDescriptorValueTypeEnumSlice(c *Client, i interface{}, res *LogMetric) []LogMetricMetricDescriptorValueTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricMetricDescriptorValueTypeEnum{}
	}

	if len(a) == 0 {
		return []LogMetricMetricDescriptorValueTypeEnum{}
	}

	items := make([]LogMetricMetricDescriptorValueTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricMetricDescriptorValueTypeEnum(item.(interface{})))
	}

	return items
}

// flattenLogMetricMetricDescriptorValueTypeEnum asserts that an interface is a string, and returns a
// pointer to a *LogMetricMetricDescriptorValueTypeEnum with the same value as that string.
func flattenLogMetricMetricDescriptorValueTypeEnum(i interface{}) *LogMetricMetricDescriptorValueTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return LogMetricMetricDescriptorValueTypeEnumRef(s)
}

// flattenLogMetricMetricDescriptorLaunchStageEnumMap flattens the contents of LogMetricMetricDescriptorLaunchStageEnum from a JSON
// response object.
func flattenLogMetricMetricDescriptorLaunchStageEnumMap(c *Client, i interface{}, res *LogMetric) map[string]LogMetricMetricDescriptorLaunchStageEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogMetricMetricDescriptorLaunchStageEnum{}
	}

	if len(a) == 0 {
		return map[string]LogMetricMetricDescriptorLaunchStageEnum{}
	}

	items := make(map[string]LogMetricMetricDescriptorLaunchStageEnum)
	for k, item := range a {
		items[k] = *flattenLogMetricMetricDescriptorLaunchStageEnum(item.(interface{}))
	}

	return items
}

// flattenLogMetricMetricDescriptorLaunchStageEnumSlice flattens the contents of LogMetricMetricDescriptorLaunchStageEnum from a JSON
// response object.
func flattenLogMetricMetricDescriptorLaunchStageEnumSlice(c *Client, i interface{}, res *LogMetric) []LogMetricMetricDescriptorLaunchStageEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []LogMetricMetricDescriptorLaunchStageEnum{}
	}

	if len(a) == 0 {
		return []LogMetricMetricDescriptorLaunchStageEnum{}
	}

	items := make([]LogMetricMetricDescriptorLaunchStageEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogMetricMetricDescriptorLaunchStageEnum(item.(interface{})))
	}

	return items
}

// flattenLogMetricMetricDescriptorLaunchStageEnum asserts that an interface is a string, and returns a
// pointer to a *LogMetricMetricDescriptorLaunchStageEnum with the same value as that string.
func flattenLogMetricMetricDescriptorLaunchStageEnum(i interface{}) *LogMetricMetricDescriptorLaunchStageEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return LogMetricMetricDescriptorLaunchStageEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *LogMetric) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalLogMetric(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if nr.Project == nil && ncr.Project == nil {
			c.Config.Logger.Info("Both Project fields null - considering equal.")
		} else if nr.Project == nil || ncr.Project == nil {
			c.Config.Logger.Info("Only one Project field is null - considering unequal.")
			return false
		} else if *nr.Project != *ncr.Project {
			return false
		}
		if nr.Name == nil && ncr.Name == nil {
			c.Config.Logger.Info("Both Name fields null - considering equal.")
		} else if nr.Name == nil || ncr.Name == nil {
			c.Config.Logger.Info("Only one Name field is null - considering unequal.")
			return false
		} else if *nr.Name != *ncr.Name {
			return false
		}
		return true
	}
}

type logMetricDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         logMetricApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToLogMetricDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]logMetricDiff, error) {
	opNamesToFieldDiffs := make(map[string][]*dcl.FieldDiff)
	// Map each operation name to the field diffs associated with it.
	for _, fd := range fds {
		for _, ro := range fd.ResultingOperation {
			if fieldDiffs, ok := opNamesToFieldDiffs[ro]; ok {
				fieldDiffs = append(fieldDiffs, fd)
				opNamesToFieldDiffs[ro] = fieldDiffs
			} else {
				config.Logger.Infof("%s required due to diff: %v", ro, fd)
				opNamesToFieldDiffs[ro] = []*dcl.FieldDiff{fd}
			}
		}
	}
	var diffs []logMetricDiff
	// For each operation name, create a logMetricDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := logMetricDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToLogMetricApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToLogMetricApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (logMetricApiOperation, error) {
	switch opName {

	case "updateLogMetricUpdateOperation":
		return &updateLogMetricUpdateOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractLogMetricFields(r *LogMetric) error {
	vMetricDescriptor := r.MetricDescriptor
	if vMetricDescriptor == nil {
		// note: explicitly not the empty object.
		vMetricDescriptor = &LogMetricMetricDescriptor{}
	}
	if err := extractLogMetricMetricDescriptorFields(r, vMetricDescriptor); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetricDescriptor) {
		r.MetricDescriptor = vMetricDescriptor
	}
	vBucketOptions := r.BucketOptions
	if vBucketOptions == nil {
		// note: explicitly not the empty object.
		vBucketOptions = &LogMetricBucketOptions{}
	}
	if err := extractLogMetricBucketOptionsFields(r, vBucketOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBucketOptions) {
		r.BucketOptions = vBucketOptions
	}
	return nil
}
func extractLogMetricMetricDescriptorFields(r *LogMetric, o *LogMetricMetricDescriptor) error {
	vMetadata := o.Metadata
	if vMetadata == nil {
		// note: explicitly not the empty object.
		vMetadata = &LogMetricMetricDescriptorMetadata{}
	}
	if err := extractLogMetricMetricDescriptorMetadataFields(r, vMetadata); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetadata) {
		o.Metadata = vMetadata
	}
	return nil
}
func extractLogMetricMetricDescriptorLabelsFields(r *LogMetric, o *LogMetricMetricDescriptorLabels) error {
	return nil
}
func extractLogMetricMetricDescriptorMetadataFields(r *LogMetric, o *LogMetricMetricDescriptorMetadata) error {
	return nil
}
func extractLogMetricBucketOptionsFields(r *LogMetric, o *LogMetricBucketOptions) error {
	vLinearBuckets := o.LinearBuckets
	if vLinearBuckets == nil {
		// note: explicitly not the empty object.
		vLinearBuckets = &LogMetricBucketOptionsLinearBuckets{}
	}
	if err := extractLogMetricBucketOptionsLinearBucketsFields(r, vLinearBuckets); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinearBuckets) {
		o.LinearBuckets = vLinearBuckets
	}
	vExponentialBuckets := o.ExponentialBuckets
	if vExponentialBuckets == nil {
		// note: explicitly not the empty object.
		vExponentialBuckets = &LogMetricBucketOptionsExponentialBuckets{}
	}
	if err := extractLogMetricBucketOptionsExponentialBucketsFields(r, vExponentialBuckets); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExponentialBuckets) {
		o.ExponentialBuckets = vExponentialBuckets
	}
	vExplicitBuckets := o.ExplicitBuckets
	if vExplicitBuckets == nil {
		// note: explicitly not the empty object.
		vExplicitBuckets = &LogMetricBucketOptionsExplicitBuckets{}
	}
	if err := extractLogMetricBucketOptionsExplicitBucketsFields(r, vExplicitBuckets); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExplicitBuckets) {
		o.ExplicitBuckets = vExplicitBuckets
	}
	return nil
}
func extractLogMetricBucketOptionsLinearBucketsFields(r *LogMetric, o *LogMetricBucketOptionsLinearBuckets) error {
	return nil
}
func extractLogMetricBucketOptionsExponentialBucketsFields(r *LogMetric, o *LogMetricBucketOptionsExponentialBuckets) error {
	return nil
}
func extractLogMetricBucketOptionsExplicitBucketsFields(r *LogMetric, o *LogMetricBucketOptionsExplicitBuckets) error {
	return nil
}

func postReadExtractLogMetricFields(r *LogMetric) error {
	vMetricDescriptor := r.MetricDescriptor
	if vMetricDescriptor == nil {
		// note: explicitly not the empty object.
		vMetricDescriptor = &LogMetricMetricDescriptor{}
	}
	if err := postReadExtractLogMetricMetricDescriptorFields(r, vMetricDescriptor); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetricDescriptor) {
		r.MetricDescriptor = vMetricDescriptor
	}
	vBucketOptions := r.BucketOptions
	if vBucketOptions == nil {
		// note: explicitly not the empty object.
		vBucketOptions = &LogMetricBucketOptions{}
	}
	if err := postReadExtractLogMetricBucketOptionsFields(r, vBucketOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBucketOptions) {
		r.BucketOptions = vBucketOptions
	}
	return nil
}
func postReadExtractLogMetricMetricDescriptorFields(r *LogMetric, o *LogMetricMetricDescriptor) error {
	vMetadata := o.Metadata
	if vMetadata == nil {
		// note: explicitly not the empty object.
		vMetadata = &LogMetricMetricDescriptorMetadata{}
	}
	if err := extractLogMetricMetricDescriptorMetadataFields(r, vMetadata); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetadata) {
		o.Metadata = vMetadata
	}
	return nil
}
func postReadExtractLogMetricMetricDescriptorLabelsFields(r *LogMetric, o *LogMetricMetricDescriptorLabels) error {
	return nil
}
func postReadExtractLogMetricMetricDescriptorMetadataFields(r *LogMetric, o *LogMetricMetricDescriptorMetadata) error {
	return nil
}
func postReadExtractLogMetricBucketOptionsFields(r *LogMetric, o *LogMetricBucketOptions) error {
	vLinearBuckets := o.LinearBuckets
	if vLinearBuckets == nil {
		// note: explicitly not the empty object.
		vLinearBuckets = &LogMetricBucketOptionsLinearBuckets{}
	}
	if err := extractLogMetricBucketOptionsLinearBucketsFields(r, vLinearBuckets); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinearBuckets) {
		o.LinearBuckets = vLinearBuckets
	}
	vExponentialBuckets := o.ExponentialBuckets
	if vExponentialBuckets == nil {
		// note: explicitly not the empty object.
		vExponentialBuckets = &LogMetricBucketOptionsExponentialBuckets{}
	}
	if err := extractLogMetricBucketOptionsExponentialBucketsFields(r, vExponentialBuckets); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExponentialBuckets) {
		o.ExponentialBuckets = vExponentialBuckets
	}
	vExplicitBuckets := o.ExplicitBuckets
	if vExplicitBuckets == nil {
		// note: explicitly not the empty object.
		vExplicitBuckets = &LogMetricBucketOptionsExplicitBuckets{}
	}
	if err := extractLogMetricBucketOptionsExplicitBucketsFields(r, vExplicitBuckets); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExplicitBuckets) {
		o.ExplicitBuckets = vExplicitBuckets
	}
	return nil
}
func postReadExtractLogMetricBucketOptionsLinearBucketsFields(r *LogMetric, o *LogMetricBucketOptionsLinearBuckets) error {
	return nil
}
func postReadExtractLogMetricBucketOptionsExponentialBucketsFields(r *LogMetric, o *LogMetricBucketOptionsExponentialBuckets) error {
	return nil
}
func postReadExtractLogMetricBucketOptionsExplicitBucketsFields(r *LogMetric, o *LogMetricBucketOptionsExplicitBuckets) error {
	return nil
}
