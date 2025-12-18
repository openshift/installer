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
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *MonitoredProject) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.MetricsScope, "MetricsScope"); err != nil {
		return err
	}
	return nil
}
func (r *MonitoredProject) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://monitoring.googleapis.com/v1/", params)
}

func (r *MonitoredProject) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"metricsScope": dcl.ValueOrEmptyString(nr.MetricsScope),
	}
	return dcl.URL("locations/global/metricsScopes/{{metricsScope}}", nr.basePath(), userBasePath, params), nil
}

func (r *MonitoredProject) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"metricsScope": dcl.ValueOrEmptyString(nr.MetricsScope),
	}
	return dcl.URL("locations/global/metricsScopes/{{metricsScope}}", nr.basePath(), userBasePath, params), nil

}

func (r *MonitoredProject) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"metricsScope": dcl.ValueOrEmptyString(nr.MetricsScope),
	}
	return dcl.URL("locations/global/metricsScopes/{{metricsScope}}/projects", nr.basePath(), userBasePath, params), nil

}

func (r *MonitoredProject) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"metricsScope": dcl.ValueOrEmptyString(nr.MetricsScope),
		"name":         dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("locations/global/metricsScopes/{{metricsScope}}/projects/{{name}}", nr.basePath(), userBasePath, params), nil
}

// monitoredProjectApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type monitoredProjectApiOperation interface {
	do(context.Context, *MonitoredProject, *Client) error
}

func (c *Client) listMonitoredProjectRaw(ctx context.Context, r *MonitoredProject, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != MonitoredProjectMaxPage {
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

type listMonitoredProjectOperation struct {
	MonitoredProjects []map[string]interface{} `json:"monitoredProjects"`
	Token             string                   `json:"nextPageToken"`
}

func (c *Client) listMonitoredProject(ctx context.Context, r *MonitoredProject, pageToken string, pageSize int32) ([]*MonitoredProject, string, error) {
	b, err := c.listMonitoredProjectRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listMonitoredProjectOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*MonitoredProject
	for _, v := range m.MonitoredProjects {
		res, err := unmarshalMapMonitoredProject(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.MetricsScope = r.MetricsScope
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllMonitoredProject(ctx context.Context, f func(*MonitoredProject) bool, resources []*MonitoredProject) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteMonitoredProject(ctx, res)
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

type deleteMonitoredProjectOperation struct{}

func (op *deleteMonitoredProjectOperation) do(ctx context.Context, r *MonitoredProject, c *Client) error {
	r, err := c.GetMonitoredProject(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "MonitoredProject not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetMonitoredProject checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	resp, err := dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return err
	}

	// wait for object to be deleted.
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		return err
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetMonitoredProject(ctx, r)
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
type createMonitoredProjectOperation struct {
	response map[string]interface{}
}

func (op *createMonitoredProjectOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createMonitoredProjectOperation) do(ctx context.Context, r *MonitoredProject, c *Client) error {
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
	// wait for object to be created.
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetMonitoredProject(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) monitoredProjectDiffsForRawDesired(ctx context.Context, rawDesired *MonitoredProject, opts ...dcl.ApplyOption) (initial, desired *MonitoredProject, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *MonitoredProject
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*MonitoredProject); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected MonitoredProject, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetMonitoredProject(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a MonitoredProject resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve MonitoredProject resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that MonitoredProject resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeMonitoredProjectDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for MonitoredProject: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for MonitoredProject: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractMonitoredProjectFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeMonitoredProjectInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for MonitoredProject: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeMonitoredProjectDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for MonitoredProject: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffMonitoredProject(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeMonitoredProjectInitialState(rawInitial, rawDesired *MonitoredProject) (*MonitoredProject, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeMonitoredProjectDesiredState(rawDesired, rawInitial *MonitoredProject, opts ...dcl.ApplyOption) (*MonitoredProject, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &MonitoredProject{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.NameToSelfLink(rawDesired.MetricsScope, rawInitial.MetricsScope) {
		canonicalDesired.MetricsScope = rawInitial.MetricsScope
	} else {
		canonicalDesired.MetricsScope = rawDesired.MetricsScope
	}
	return canonicalDesired, nil
}

func canonicalizeMonitoredProjectNewState(c *Client, rawNew, rawDesired *MonitoredProject) (*MonitoredProject, error) {

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

	rawNew.MetricsScope = rawDesired.MetricsScope

	return rawNew, nil
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffMonitoredProject(c *Client, desired, actual *MonitoredProject, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.MetricsScope, actual.MetricsScope, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetricsScope")); len(ds) != 0 || err != nil {
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

// urlNormalized returns a copy of the resource struct with values normalized
// for URL substitutions. For instance, it converts long-form self-links to
// short-form so they can be substituted in.
func (r *MonitoredProject) urlNormalized() *MonitoredProject {
	normalized := dcl.Copy(*r).(MonitoredProject)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.MetricsScope = dcl.SelfLinkToName(r.MetricsScope)
	return &normalized
}

func (r *MonitoredProject) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the MonitoredProject resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *MonitoredProject) marshal(c *Client) ([]byte, error) {
	m, err := expandMonitoredProject(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling MonitoredProject: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalMonitoredProject decodes JSON responses into the MonitoredProject resource schema.
func unmarshalMonitoredProject(b []byte, c *Client, res *MonitoredProject) (*MonitoredProject, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapMonitoredProject(m, c, res)
}

func unmarshalMapMonitoredProject(m map[string]interface{}, c *Client, res *MonitoredProject) (*MonitoredProject, error) {

	flattened := flattenMonitoredProject(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandMonitoredProject expands MonitoredProject into a JSON request object.
func expandMonitoredProject(c *Client, f *MonitoredProject) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("locations/global/metricsScopes/%s/projects/%s", f.Name, dcl.SelfLinkToName(f.MetricsScope), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding MetricsScope into metricsScope: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metricsScope"] = v
	}

	return m, nil
}

// flattenMonitoredProject flattens MonitoredProject from a JSON request object into the
// MonitoredProject type.
func flattenMonitoredProject(c *Client, i interface{}, res *MonitoredProject) *MonitoredProject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &MonitoredProject{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.MetricsScope = dcl.FlattenString(m["metricsScope"])

	return resultRes
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *MonitoredProject) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalMonitoredProject(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if nr.MetricsScope == nil && ncr.MetricsScope == nil {
			c.Config.Logger.Info("Both MetricsScope fields null - considering equal.")
		} else if nr.MetricsScope == nil || ncr.MetricsScope == nil {
			c.Config.Logger.Info("Only one MetricsScope field is null - considering unequal.")
			return false
		} else if *nr.MetricsScope != *ncr.MetricsScope {
			return false
		}
		return true
	}
}

type monitoredProjectDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         monitoredProjectApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToMonitoredProjectDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]monitoredProjectDiff, error) {
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
	var diffs []monitoredProjectDiff
	// For each operation name, create a monitoredProjectDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := monitoredProjectDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToMonitoredProjectApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToMonitoredProjectApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (monitoredProjectApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractMonitoredProjectFields(r *MonitoredProject) error {
	return nil
}

func postReadExtractMonitoredProjectFields(r *MonitoredProject) error {
	return nil
}
