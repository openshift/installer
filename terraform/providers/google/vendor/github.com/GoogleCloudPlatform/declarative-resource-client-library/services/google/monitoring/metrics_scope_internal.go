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
package monitoring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *MetricsScope) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	return nil
}
func (r *MetricsScopeMonitoredProjects) validate() error {
	return nil
}
func (r *MetricsScope) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://monitoring.googleapis.com/v1/", params)
}

func (r *MetricsScope) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("locations/global/metricsScopes/{{name}}", nr.basePath(), userBasePath, params), nil
}

// metricsScopeApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type metricsScopeApiOperation interface {
	do(context.Context, *MetricsScope, *Client) error
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createMetricsScopeOperation struct {
	response map[string]interface{}
}

func (op *createMetricsScopeOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (c *Client) getMetricsScopeRaw(ctx context.Context, r *MetricsScope) ([]byte, error) {

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

func (c *Client) metricsScopeDiffsForRawDesired(ctx context.Context, rawDesired *MetricsScope, opts ...dcl.ApplyOption) (initial, desired *MetricsScope, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *MetricsScope
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*MetricsScope); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected MetricsScope, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetMetricsScope(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a MetricsScope resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve MetricsScope resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that MetricsScope resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeMetricsScopeDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for MetricsScope: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for MetricsScope: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractMetricsScopeFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeMetricsScopeInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for MetricsScope: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeMetricsScopeDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for MetricsScope: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffMetricsScope(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeMetricsScopeInitialState(rawInitial, rawDesired *MetricsScope) (*MetricsScope, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeMetricsScopeDesiredState(rawDesired, rawInitial *MetricsScope, opts ...dcl.ApplyOption) (*MetricsScope, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &MetricsScope{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	return canonicalDesired, nil
}

func canonicalizeMetricsScopeNewState(c *Client, rawNew, rawDesired *MetricsScope) (*MetricsScope, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.MonitoredProjects) && dcl.IsEmptyValueIndirect(rawDesired.MonitoredProjects) {
		rawNew.MonitoredProjects = rawDesired.MonitoredProjects
	} else {
		rawNew.MonitoredProjects = canonicalizeNewMetricsScopeMonitoredProjectsSlice(c, rawDesired.MonitoredProjects, rawNew.MonitoredProjects)
	}

	return rawNew, nil
}

func canonicalizeMetricsScopeMonitoredProjects(des, initial *MetricsScopeMonitoredProjects, opts ...dcl.ApplyOption) *MetricsScopeMonitoredProjects {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &MetricsScopeMonitoredProjects{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}

	return cDes
}

func canonicalizeMetricsScopeMonitoredProjectsSlice(des, initial []MetricsScopeMonitoredProjects, opts ...dcl.ApplyOption) []MetricsScopeMonitoredProjects {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]MetricsScopeMonitoredProjects, 0, len(des))
		for _, d := range des {
			cd := canonicalizeMetricsScopeMonitoredProjects(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]MetricsScopeMonitoredProjects, 0, len(des))
	for i, d := range des {
		cd := canonicalizeMetricsScopeMonitoredProjects(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewMetricsScopeMonitoredProjects(c *Client, des, nw *MetricsScopeMonitoredProjects) *MetricsScopeMonitoredProjects {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for MetricsScopeMonitoredProjects while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}

	return nw
}

func canonicalizeNewMetricsScopeMonitoredProjectsSet(c *Client, des, nw []MetricsScopeMonitoredProjects) []MetricsScopeMonitoredProjects {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []MetricsScopeMonitoredProjects
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareMetricsScopeMonitoredProjectsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewMetricsScopeMonitoredProjects(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewMetricsScopeMonitoredProjectsSlice(c *Client, des, nw []MetricsScopeMonitoredProjects) []MetricsScopeMonitoredProjects {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []MetricsScopeMonitoredProjects
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewMetricsScopeMonitoredProjects(c, &d, &n))
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
func diffMetricsScope(c *Client, desired, actual *MetricsScope, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.MonitoredProjects, actual.MonitoredProjects, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareMetricsScopeMonitoredProjectsNewStyle, EmptyObject: EmptyMetricsScopeMonitoredProjects, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MonitoredProjects")); len(ds) != 0 || err != nil {
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
func compareMetricsScopeMonitoredProjectsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*MetricsScopeMonitoredProjects)
	if !ok {
		desiredNotPointer, ok := d.(MetricsScopeMonitoredProjects)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a MetricsScopeMonitoredProjects or *MetricsScopeMonitoredProjects", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*MetricsScopeMonitoredProjects)
	if !ok {
		actualNotPointer, ok := a.(MetricsScopeMonitoredProjects)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a MetricsScopeMonitoredProjects", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreateTime, actual.CreateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreateTime")); len(ds) != 0 || err != nil {
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
func (r *MetricsScope) urlNormalized() *MetricsScope {
	normalized := dcl.Copy(*r).(MetricsScope)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	return &normalized
}

func (r *MetricsScope) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the MetricsScope resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *MetricsScope) marshal(c *Client) ([]byte, error) {
	m, err := expandMetricsScope(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling MetricsScope: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalMetricsScope decodes JSON responses into the MetricsScope resource schema.
func unmarshalMetricsScope(b []byte, c *Client, res *MetricsScope) (*MetricsScope, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapMetricsScope(m, c, res)
}

func unmarshalMapMetricsScope(m map[string]interface{}, c *Client, res *MetricsScope) (*MetricsScope, error) {

	flattened := flattenMetricsScope(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandMetricsScope expands MetricsScope into a JSON request object.
func expandMetricsScope(c *Client, f *MetricsScope) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.ExpandProjectIDsToNumbers(c.Config, f.Name); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}

	return m, nil
}

// flattenMetricsScope flattens MetricsScope from a JSON request object into the
// MetricsScope type.
func flattenMetricsScope(c *Client, i interface{}, res *MetricsScope) *MetricsScope {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &MetricsScope{}
	resultRes.Name = dcl.FlattenProjectNumbersToIDs(c.Config, dcl.FlattenString(m["name"]))
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.MonitoredProjects = flattenMetricsScopeMonitoredProjectsSlice(c, m["monitoredProjects"], res)

	return resultRes
}

// expandMetricsScopeMonitoredProjectsMap expands the contents of MetricsScopeMonitoredProjects into a JSON
// request object.
func expandMetricsScopeMonitoredProjectsMap(c *Client, f map[string]MetricsScopeMonitoredProjects, res *MetricsScope) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandMetricsScopeMonitoredProjects(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandMetricsScopeMonitoredProjectsSlice expands the contents of MetricsScopeMonitoredProjects into a JSON
// request object.
func expandMetricsScopeMonitoredProjectsSlice(c *Client, f []MetricsScopeMonitoredProjects, res *MetricsScope) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandMetricsScopeMonitoredProjects(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenMetricsScopeMonitoredProjectsMap flattens the contents of MetricsScopeMonitoredProjects from a JSON
// response object.
func flattenMetricsScopeMonitoredProjectsMap(c *Client, i interface{}, res *MetricsScope) map[string]MetricsScopeMonitoredProjects {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]MetricsScopeMonitoredProjects{}
	}

	if len(a) == 0 {
		return map[string]MetricsScopeMonitoredProjects{}
	}

	items := make(map[string]MetricsScopeMonitoredProjects)
	for k, item := range a {
		items[k] = *flattenMetricsScopeMonitoredProjects(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenMetricsScopeMonitoredProjectsSlice flattens the contents of MetricsScopeMonitoredProjects from a JSON
// response object.
func flattenMetricsScopeMonitoredProjectsSlice(c *Client, i interface{}, res *MetricsScope) []MetricsScopeMonitoredProjects {
	a, ok := i.([]interface{})
	if !ok {
		return []MetricsScopeMonitoredProjects{}
	}

	if len(a) == 0 {
		return []MetricsScopeMonitoredProjects{}
	}

	items := make([]MetricsScopeMonitoredProjects, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenMetricsScopeMonitoredProjects(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandMetricsScopeMonitoredProjects expands an instance of MetricsScopeMonitoredProjects into a JSON
// request object.
func expandMetricsScopeMonitoredProjects(c *Client, f *MetricsScopeMonitoredProjects, res *MetricsScope) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}

	return m, nil
}

// flattenMetricsScopeMonitoredProjects flattens an instance of MetricsScopeMonitoredProjects from a JSON
// response object.
func flattenMetricsScopeMonitoredProjects(c *Client, i interface{}, res *MetricsScope) *MetricsScopeMonitoredProjects {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &MetricsScopeMonitoredProjects{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyMetricsScopeMonitoredProjects
	}
	r.Name = dcl.FlattenString(m["name"])
	r.CreateTime = dcl.FlattenString(m["createTime"])

	return r
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *MetricsScope) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalMetricsScope(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

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

type metricsScopeDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         metricsScopeApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToMetricsScopeDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]metricsScopeDiff, error) {
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
	var diffs []metricsScopeDiff
	// For each operation name, create a metricsScopeDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := metricsScopeDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToMetricsScopeApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToMetricsScopeApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (metricsScopeApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractMetricsScopeFields(r *MetricsScope) error {
	return nil
}
func extractMetricsScopeMonitoredProjectsFields(r *MetricsScope, o *MetricsScopeMonitoredProjects) error {
	return nil
}

func postReadExtractMetricsScopeFields(r *MetricsScope) error {
	return nil
}
func postReadExtractMetricsScopeMonitoredProjectsFields(r *MetricsScope, o *MetricsScopeMonitoredProjects) error {
	return nil
}
