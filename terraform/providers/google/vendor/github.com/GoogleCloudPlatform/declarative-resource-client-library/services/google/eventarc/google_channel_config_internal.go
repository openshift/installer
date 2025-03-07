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
package eventarc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *GoogleChannelConfig) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	return nil
}
func (r *GoogleChannelConfig) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://eventarc.googleapis.com/v1/", params)
}

func (r *GoogleChannelConfig) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/googleChannelConfig", nr.basePath(), userBasePath, params), nil
}

func (r *GoogleChannelConfig) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/googleChannelConfig", nr.basePath(), userBasePath, params), nil
}

// googleChannelConfigApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type googleChannelConfigApiOperation interface {
	do(context.Context, *GoogleChannelConfig, *Client) error
}

// newUpdateGoogleChannelConfigUpdateGoogleChannelConfigRequest creates a request for an
// GoogleChannelConfig resource's UpdateGoogleChannelConfig update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateGoogleChannelConfigUpdateGoogleChannelConfigRequest(ctx context.Context, f *GoogleChannelConfig, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.CryptoKeyName; !dcl.IsEmptyValueIndirect(v) {
		req["cryptoKeyName"] = v
	}
	return req, nil
}

// marshalUpdateGoogleChannelConfigUpdateGoogleChannelConfigRequest converts the update into
// the final JSON request body.
func marshalUpdateGoogleChannelConfigUpdateGoogleChannelConfigRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateGoogleChannelConfigUpdateGoogleChannelConfigOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateGoogleChannelConfigUpdateGoogleChannelConfigOperation) do(ctx context.Context, r *GoogleChannelConfig, c *Client) error {
	_, err := c.GetGoogleChannelConfig(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateGoogleChannelConfig")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateGoogleChannelConfigUpdateGoogleChannelConfigRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateGoogleChannelConfigUpdateGoogleChannelConfigRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) deleteAllGoogleChannelConfig(ctx context.Context, f func(*GoogleChannelConfig) bool, resources []*GoogleChannelConfig) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteGoogleChannelConfig(ctx, res)
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

type deleteGoogleChannelConfigOperation struct{}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createGoogleChannelConfigOperation struct {
	response map[string]interface{}
}

func (op *createGoogleChannelConfigOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (c *Client) getGoogleChannelConfigRaw(ctx context.Context, r *GoogleChannelConfig) ([]byte, error) {

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

func (c *Client) googleChannelConfigDiffsForRawDesired(ctx context.Context, rawDesired *GoogleChannelConfig, opts ...dcl.ApplyOption) (initial, desired *GoogleChannelConfig, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *GoogleChannelConfig
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*GoogleChannelConfig); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected GoogleChannelConfig, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetGoogleChannelConfig(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a GoogleChannelConfig resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve GoogleChannelConfig resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that GoogleChannelConfig resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeGoogleChannelConfigDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for GoogleChannelConfig: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for GoogleChannelConfig: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractGoogleChannelConfigFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeGoogleChannelConfigInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for GoogleChannelConfig: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeGoogleChannelConfigDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for GoogleChannelConfig: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffGoogleChannelConfig(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeGoogleChannelConfigInitialState(rawInitial, rawDesired *GoogleChannelConfig) (*GoogleChannelConfig, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeGoogleChannelConfigDesiredState(rawDesired, rawInitial *GoogleChannelConfig, opts ...dcl.ApplyOption) (*GoogleChannelConfig, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &GoogleChannelConfig{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.IsZeroValue(rawDesired.CryptoKeyName) || (dcl.IsEmptyValueIndirect(rawDesired.CryptoKeyName) && dcl.IsEmptyValueIndirect(rawInitial.CryptoKeyName)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.CryptoKeyName = rawInitial.CryptoKeyName
	} else {
		canonicalDesired.CryptoKeyName = rawDesired.CryptoKeyName
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	return canonicalDesired, nil
}

func canonicalizeGoogleChannelConfigNewState(c *Client, rawNew, rawDesired *GoogleChannelConfig) (*GoogleChannelConfig, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.CryptoKeyName) && dcl.IsEmptyValueIndirect(rawDesired.CryptoKeyName) {
		rawNew.CryptoKeyName = rawDesired.CryptoKeyName
	} else {
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffGoogleChannelConfig(c *Client, desired, actual *GoogleChannelConfig, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CryptoKeyName, actual.CryptoKeyName, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateGoogleChannelConfigUpdateGoogleChannelConfigOperation")}, fn.AddNest("CryptoKeyName")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
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
func (r *GoogleChannelConfig) urlNormalized() *GoogleChannelConfig {
	normalized := dcl.Copy(*r).(GoogleChannelConfig)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.CryptoKeyName = r.CryptoKeyName
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *GoogleChannelConfig) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateGoogleChannelConfig" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/googleChannelConfig", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the GoogleChannelConfig resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *GoogleChannelConfig) marshal(c *Client) ([]byte, error) {
	m, err := expandGoogleChannelConfig(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling GoogleChannelConfig: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalGoogleChannelConfig decodes JSON responses into the GoogleChannelConfig resource schema.
func unmarshalGoogleChannelConfig(b []byte, c *Client, res *GoogleChannelConfig) (*GoogleChannelConfig, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapGoogleChannelConfig(m, c, res)
}

func unmarshalMapGoogleChannelConfig(m map[string]interface{}, c *Client, res *GoogleChannelConfig) (*GoogleChannelConfig, error) {

	flattened := flattenGoogleChannelConfig(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandGoogleChannelConfig expands GoogleChannelConfig into a JSON request object.
func expandGoogleChannelConfig(c *Client, f *GoogleChannelConfig) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/googleChannelConfig", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.CryptoKeyName; dcl.ValueShouldBeSent(v) {
		m["cryptoKeyName"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}

	return m, nil
}

// flattenGoogleChannelConfig flattens GoogleChannelConfig from a JSON request object into the
// GoogleChannelConfig type.
func flattenGoogleChannelConfig(c *Client, i interface{}, res *GoogleChannelConfig) *GoogleChannelConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &GoogleChannelConfig{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.CryptoKeyName = dcl.FlattenString(m["cryptoKeyName"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *GoogleChannelConfig) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalGoogleChannelConfig(b, c, r)
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
		if nr.Location == nil && ncr.Location == nil {
			c.Config.Logger.Info("Both Location fields null - considering equal.")
		} else if nr.Location == nil || ncr.Location == nil {
			c.Config.Logger.Info("Only one Location field is null - considering unequal.")
			return false
		} else if *nr.Location != *ncr.Location {
			return false
		}
		return true
	}
}

type googleChannelConfigDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         googleChannelConfigApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToGoogleChannelConfigDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]googleChannelConfigDiff, error) {
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
	var diffs []googleChannelConfigDiff
	// For each operation name, create a googleChannelConfigDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := googleChannelConfigDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToGoogleChannelConfigApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToGoogleChannelConfigApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (googleChannelConfigApiOperation, error) {
	switch opName {

	case "updateGoogleChannelConfigUpdateGoogleChannelConfigOperation":
		return &updateGoogleChannelConfigUpdateGoogleChannelConfigOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractGoogleChannelConfigFields(r *GoogleChannelConfig) error {
	return nil
}

func postReadExtractGoogleChannelConfigFields(r *GoogleChannelConfig) error {
	return nil
}
