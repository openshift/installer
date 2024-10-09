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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *Project) validate() error {

	return nil
}
func (r *Project) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://cloudresourcemanager.googleapis.com/", params)
}

func (r *Project) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("v1/projects/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Project) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{}
	return dcl.URL("v1/projects", nr.basePath(), userBasePath, params), nil

}

func (r *Project) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name": dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("v1/projects/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Project) SetPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{
		"name": *nr.Name,
	}
	return dcl.URL("v1/projects/{{name}}:setIamPolicy", nr.basePath(), userBasePath, fields)
}

func (r *Project) SetPolicyVerb() string {
	return "POST"
}

func (r *Project) getPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{
		"name": *nr.Name,
	}
	return dcl.URL("v1/projects/{{name}}:getIamPolicy", nr.basePath(), userBasePath, fields)
}

func (r *Project) IAMPolicyVersion() int {
	return 3
}

// projectApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type projectApiOperation interface {
	do(context.Context, *Project, *Client) error
}

// newUpdateProjectUpdateProjectRequest creates a request for an
// Project resource's UpdateProject update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateProjectUpdateProjectRequest(ctx context.Context, f *Project, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	return req, nil
}

// marshalUpdateProjectUpdateProjectRequest converts the update into
// the final JSON request body.
func marshalUpdateProjectUpdateProjectRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateProjectUpdateProjectOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateProjectUpdateProjectOperation) do(ctx context.Context, r *Project, c *Client) error {
	_, err := c.GetProject(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateProject")
	if err != nil {
		return err
	}

	req, err := newUpdateProjectUpdateProjectRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateProjectUpdateProjectRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PUT", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.CRMOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listProjectRaw(ctx context.Context, r *Project, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ProjectMaxPage {
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

type listProjectOperation struct {
	Projects []map[string]interface{} `json:"projects"`
	Token    string                   `json:"nextPageToken"`
}

func (c *Client) listProject(ctx context.Context, r *Project, pageToken string, pageSize int32) ([]*Project, string, error) {
	b, err := c.listProjectRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listProjectOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Project
	for _, v := range m.Projects {
		res, err := unmarshalMapProject(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Parent = r.Parent
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllProject(ctx context.Context, f func(*Project) bool, resources []*Project) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteProject(ctx, res)
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

type deleteProjectOperation struct{}

func (op *deleteProjectOperation) do(ctx context.Context, r *Project, c *Client) error {
	r, err := c.GetProject(ctx, r)
	if err != nil {
		if dcl.IsNotFoundOrCode(err, 403) {
			c.Config.Logger.InfoWithContextf(ctx, "Project not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetProject checking for existence. error: %v", err)
		return err
	}

	if projectDeletePrecondition(r) {
		return nil
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	_, err = dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return fmt.Errorf("failed to delete Project: %w", err)
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createProjectOperation struct {
	response map[string]interface{}
}

func (op *createProjectOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createProjectOperation) do(ctx context.Context, r *Project, c *Client) error {
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
	var o operations.CRMOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetProject(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getProjectRaw(ctx context.Context, r *Project) ([]byte, error) {

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

func (c *Client) projectDiffsForRawDesired(ctx context.Context, rawDesired *Project, opts ...dcl.ApplyOption) (initial, desired *Project, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Project
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Project); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Project, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetProject(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFoundOrCode(err, 403) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Project resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Project resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Project resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeProjectDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Project: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Project: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractProjectFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeProjectInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Project: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeProjectDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Project: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffProject(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeProjectInitialState(rawInitial, rawDesired *Project) (*Project, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeProjectDesiredState(rawDesired, rawInitial *Project, opts ...dcl.ApplyOption) (*Project, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &Project{}
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.StringCanonicalize(rawDesired.DisplayName, rawInitial.DisplayName) {
		canonicalDesired.DisplayName = rawInitial.DisplayName
	} else {
		canonicalDesired.DisplayName = rawDesired.DisplayName
	}
	if dcl.StringCanonicalize(rawDesired.Parent, rawInitial.Parent) {
		canonicalDesired.Parent = rawInitial.Parent
	} else {
		canonicalDesired.Parent = rawDesired.Parent
	}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	return canonicalDesired, nil
}

func canonicalizeProjectNewState(c *Client, rawNew, rawDesired *Project) (*Project, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.LifecycleState) && dcl.IsEmptyValueIndirect(rawDesired.LifecycleState) {
		rawNew.LifecycleState = rawDesired.LifecycleState
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.DisplayName) && dcl.IsEmptyValueIndirect(rawDesired.DisplayName) {
		rawNew.DisplayName = rawDesired.DisplayName
	} else {
		if dcl.StringCanonicalize(rawDesired.DisplayName, rawNew.DisplayName) {
			rawNew.DisplayName = rawDesired.DisplayName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Parent) && dcl.IsEmptyValueIndirect(rawDesired.Parent) {
		rawNew.Parent = rawDesired.Parent
	} else {
		if dcl.StringCanonicalize(rawDesired.Parent, rawNew.Parent) {
			rawNew.Parent = rawDesired.Parent
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.ProjectNumber) && dcl.IsEmptyValueIndirect(rawDesired.ProjectNumber) {
		rawNew.ProjectNumber = rawDesired.ProjectNumber
	} else {
	}

	return rawNew, nil
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffProject(c *Client, desired, actual *Project, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateProjectUpdateProjectOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LifecycleState, actual.LifecycleState, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LifecycleState")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Parent, actual.Parent, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Parent")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ProjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ProjectNumber, actual.ProjectNumber, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ProjectNumber")); len(ds) != 0 || err != nil {
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
func (r *Project) urlNormalized() *Project {
	normalized := dcl.Copy(*r).(Project)
	normalized.DisplayName = dcl.SelfLinkToName(r.DisplayName)
	normalized.Parent = r.Parent
	normalized.Name = dcl.SelfLinkToName(r.Name)
	return &normalized
}

func (r *Project) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateProject" {
		fields := map[string]interface{}{
			"name": dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("v1/projects/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Project resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Project) marshal(c *Client) ([]byte, error) {
	m, err := expandProject(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Project: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalProject decodes JSON responses into the Project resource schema.
func unmarshalProject(b []byte, c *Client, res *Project) (*Project, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapProject(m, c, res)
}

func unmarshalMapProject(m map[string]interface{}, c *Client, res *Project) (*Project, error) {

	flattened := flattenProject(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandProject expands Project into a JSON request object.
func expandProject(c *Client, f *Project) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v, err := expandProjectParent(c, f.Parent, res); err != nil {
		return nil, fmt.Errorf("error expanding Parent into parent: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["parent"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["projectId"] = v
	}

	return m, nil
}

// flattenProject flattens Project from a JSON request object into the
// Project type.
func flattenProject(c *Client, i interface{}, res *Project) *Project {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Project{}
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.LifecycleState = flattenProjectLifecycleStateEnum(m["lifecycleState"])
	resultRes.DisplayName = dcl.FlattenString(m["name"])
	resultRes.Parent = flattenProjectParent(c, m["parent"], res)
	resultRes.Name = dcl.FlattenString(m["projectId"])
	resultRes.ProjectNumber = dcl.FlattenInteger(m["projectNumber"])

	return resultRes
}

// flattenProjectLifecycleStateEnumMap flattens the contents of ProjectLifecycleStateEnum from a JSON
// response object.
func flattenProjectLifecycleStateEnumMap(c *Client, i interface{}, res *Project) map[string]ProjectLifecycleStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ProjectLifecycleStateEnum{}
	}

	if len(a) == 0 {
		return map[string]ProjectLifecycleStateEnum{}
	}

	items := make(map[string]ProjectLifecycleStateEnum)
	for k, item := range a {
		items[k] = *flattenProjectLifecycleStateEnum(item.(interface{}))
	}

	return items
}

// flattenProjectLifecycleStateEnumSlice flattens the contents of ProjectLifecycleStateEnum from a JSON
// response object.
func flattenProjectLifecycleStateEnumSlice(c *Client, i interface{}, res *Project) []ProjectLifecycleStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ProjectLifecycleStateEnum{}
	}

	if len(a) == 0 {
		return []ProjectLifecycleStateEnum{}
	}

	items := make([]ProjectLifecycleStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenProjectLifecycleStateEnum(item.(interface{})))
	}

	return items
}

// flattenProjectLifecycleStateEnum asserts that an interface is a string, and returns a
// pointer to a *ProjectLifecycleStateEnum with the same value as that string.
func flattenProjectLifecycleStateEnum(i interface{}) *ProjectLifecycleStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ProjectLifecycleStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Project) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalProject(b, c, r)
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

type projectDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         projectApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToProjectDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]projectDiff, error) {
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
	var diffs []projectDiff
	// For each operation name, create a projectDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := projectDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToProjectApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToProjectApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (projectApiOperation, error) {
	switch opName {

	case "updateProjectUpdateProjectOperation":
		return &updateProjectUpdateProjectOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractProjectFields(r *Project) error {
	return nil
}

func postReadExtractProjectFields(r *Project) error {
	return nil
}
