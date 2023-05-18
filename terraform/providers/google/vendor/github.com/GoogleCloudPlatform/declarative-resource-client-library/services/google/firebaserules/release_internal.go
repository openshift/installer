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
package firebaserules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *Release) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "rulesetName"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	return nil
}
func (r *Release) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://firebaserules.googleapis.com/v1/", params)
}

func (r *Release) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/releases/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Release) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/releases", nr.basePath(), userBasePath, params), nil

}

func (r *Release) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/releases", nr.basePath(), userBasePath, params), nil

}

func (r *Release) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/releases/{{name}}", nr.basePath(), userBasePath, params), nil
}

// releaseApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type releaseApiOperation interface {
	do(context.Context, *Release, *Client) error
}

// newUpdateReleaseUpdateReleaseRequest creates a request for an
// Release resource's UpdateRelease update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateReleaseUpdateReleaseRequest(ctx context.Context, f *Release, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := dcl.DeriveField("projects/%s/rulesets/%s", f.RulesetName, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.RulesetName)); err != nil {
		return nil, fmt.Errorf("error expanding RulesetName into rulesetName: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["rulesetName"] = v
	}
	req["name"] = fmt.Sprintf("projects/%s/releases/%s", *f.Project, *f.Name)

	return req, nil
}

// marshalUpdateReleaseUpdateReleaseRequest converts the update into
// the final JSON request body.
func marshalUpdateReleaseUpdateReleaseRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	m = EncodeReleaseUpdateRequest(m)
	return json.Marshal(m)
}

type updateReleaseUpdateReleaseOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateReleaseUpdateReleaseOperation) do(ctx context.Context, r *Release, c *Client) error {
	_, err := c.GetRelease(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateRelease")
	if err != nil {
		return err
	}

	req, err := newUpdateReleaseUpdateReleaseRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateReleaseUpdateReleaseRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listReleaseRaw(ctx context.Context, r *Release, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ReleaseMaxPage {
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

type listReleaseOperation struct {
	Releases []map[string]interface{} `json:"releases"`
	Token    string                   `json:"nextPageToken"`
}

func (c *Client) listRelease(ctx context.Context, r *Release, pageToken string, pageSize int32) ([]*Release, string, error) {
	b, err := c.listReleaseRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listReleaseOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Release
	for _, v := range m.Releases {
		res, err := unmarshalMapRelease(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllRelease(ctx context.Context, f func(*Release) bool, resources []*Release) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteRelease(ctx, res)
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

type deleteReleaseOperation struct{}

func (op *deleteReleaseOperation) do(ctx context.Context, r *Release, c *Client) error {
	r, err := c.GetRelease(ctx, r)
	if err != nil {
		if dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.InfoWithContextf(ctx, "Release not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetRelease checking for existence. error: %v", err)
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
		return fmt.Errorf("failed to delete Release: %w", err)
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createReleaseOperation struct {
	response map[string]interface{}
}

func (op *createReleaseOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createReleaseOperation) do(ctx context.Context, r *Release, c *Client) error {
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

	if _, err := c.GetRelease(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getReleaseRaw(ctx context.Context, r *Release) ([]byte, error) {

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

func (c *Client) releaseDiffsForRawDesired(ctx context.Context, rawDesired *Release, opts ...dcl.ApplyOption) (initial, desired *Release, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Release
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Release); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Release, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetRelease(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Release resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Release resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Release resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeReleaseDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Release: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Release: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractReleaseFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeReleaseInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Release: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeReleaseDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Release: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffRelease(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeReleaseInitialState(rawInitial, rawDesired *Release) (*Release, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeReleaseDesiredState(rawDesired, rawInitial *Release, opts ...dcl.ApplyOption) (*Release, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &Release{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.RulesetName, rawInitial.RulesetName) {
		canonicalDesired.RulesetName = rawInitial.RulesetName
	} else {
		canonicalDesired.RulesetName = rawDesired.RulesetName
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeReleaseNewState(c *Client, rawNew, rawDesired *Release) (*Release, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.RulesetName) && dcl.IsEmptyValueIndirect(rawDesired.RulesetName) {
		rawNew.RulesetName = rawDesired.RulesetName
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.RulesetName, rawNew.RulesetName) {
			rawNew.RulesetName = rawDesired.RulesetName
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

	if dcl.IsEmptyValueIndirect(rawNew.Disabled) && dcl.IsEmptyValueIndirect(rawDesired.Disabled) {
		rawNew.Disabled = rawDesired.Disabled
	} else {
		if dcl.BoolCanonicalize(rawDesired.Disabled, rawNew.Disabled) {
			rawNew.Disabled = rawDesired.Disabled
		}
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffRelease(c *Client, desired, actual *Release, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.RulesetName, actual.RulesetName, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateReleaseUpdateReleaseOperation")}, fn.AddNest("RulesetName")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Disabled, actual.Disabled, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Disabled")); len(ds) != 0 || err != nil {
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

// urlNormalized returns a copy of the resource struct with values normalized
// for URL substitutions. For instance, it converts long-form self-links to
// short-form so they can be substituted in.
func (r *Release) urlNormalized() *Release {
	normalized := dcl.Copy(*r).(Release)
	normalized.Name = r.Name
	normalized.RulesetName = dcl.SelfLinkToName(r.RulesetName)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *Release) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateRelease" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/releases/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Release resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Release) marshal(c *Client) ([]byte, error) {
	m, err := expandRelease(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Release: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalRelease decodes JSON responses into the Release resource schema.
func unmarshalRelease(b []byte, c *Client, res *Release) (*Release, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapRelease(m, c, res)
}

func unmarshalMapRelease(m map[string]interface{}, c *Client, res *Release) (*Release, error) {

	flattened := flattenRelease(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandRelease expands Release into a JSON request object.
func expandRelease(c *Client, f *Release) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveFromPattern("projects/%s/releases/%s", f.Name, dcl.SelfLinkToName(f.Project), f.Name); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v, err := dcl.DeriveField("projects/%s/rulesets/%s", f.RulesetName, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.RulesetName)); err != nil {
		return nil, fmt.Errorf("error expanding RulesetName into rulesetName: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["rulesetName"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenRelease flattens Release from a JSON request object into the
// Release type.
func flattenRelease(c *Client, i interface{}, res *Release) *Release {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Release{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.RulesetName = dcl.FlattenString(m["rulesetName"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Disabled = dcl.FlattenBool(m["disabled"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Release) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalRelease(b, c, r)
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

type releaseDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         releaseApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToReleaseDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]releaseDiff, error) {
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
	var diffs []releaseDiff
	// For each operation name, create a releaseDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := releaseDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToReleaseApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToReleaseApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (releaseApiOperation, error) {
	switch opName {

	case "updateReleaseUpdateReleaseOperation":
		return &updateReleaseUpdateReleaseOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractReleaseFields(r *Release) error {
	return nil
}

func postReadExtractReleaseFields(r *Release) error {
	return nil
}
