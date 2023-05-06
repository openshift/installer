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

func (r *LogBucket) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Parent, "Parent"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	return nil
}
func (r *LogBucket) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://logging.googleapis.com/v2/", params)
}

func (r *LogBucket) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"location": dcl.ValueOrEmptyString(nr.Location),
		"parent":   dcl.ValueOrEmptyString(nr.Parent),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("{{parent}}/locations/{{location}}/buckets/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *LogBucket) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"location": dcl.ValueOrEmptyString(nr.Location),
		"parent":   dcl.ValueOrEmptyString(nr.Parent),
	}
	return dcl.URL("{{parent}}/locations/{{location}}/buckets", nr.basePath(), userBasePath, params), nil

}

func (r *LogBucket) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"location": dcl.ValueOrEmptyString(nr.Location),
		"parent":   dcl.ValueOrEmptyString(nr.Parent),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("{{parent}}/locations/{{location}}/buckets?bucketId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *LogBucket) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"location": dcl.ValueOrEmptyString(nr.Location),
		"parent":   dcl.ValueOrEmptyString(nr.Parent),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("{{parent}}/locations/{{location}}/buckets/{{name}}", nr.basePath(), userBasePath, params), nil
}

// logBucketApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type logBucketApiOperation interface {
	do(context.Context, *LogBucket, *Client) error
}

// newUpdateLogBucketUpdateBucketRequest creates a request for an
// LogBucket resource's UpdateBucket update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateLogBucketUpdateBucketRequest(ctx context.Context, f *LogBucket, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.RetentionDays; !dcl.IsEmptyValueIndirect(v) {
		req["retentionDays"] = v
	}
	if v := f.Locked; !dcl.IsEmptyValueIndirect(v) {
		req["locked"] = v
	}
	return req, nil
}

// marshalUpdateLogBucketUpdateBucketRequest converts the update into
// the final JSON request body.
func marshalUpdateLogBucketUpdateBucketRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateLogBucketUpdateBucketOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateLogBucketUpdateBucketOperation) do(ctx context.Context, r *LogBucket, c *Client) error {
	_, err := c.GetLogBucket(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateBucket")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateLogBucketUpdateBucketRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateLogBucketUpdateBucketRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listLogBucketRaw(ctx context.Context, r *LogBucket, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != LogBucketMaxPage {
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

type listLogBucketOperation struct {
	Buckets []map[string]interface{} `json:"buckets"`
	Token   string                   `json:"nextPageToken"`
}

func (c *Client) listLogBucket(ctx context.Context, r *LogBucket, pageToken string, pageSize int32) ([]*LogBucket, string, error) {
	b, err := c.listLogBucketRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listLogBucketOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*LogBucket
	for _, v := range m.Buckets {
		res, err := unmarshalMapLogBucket(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Location = r.Location
		res.Parent = r.Parent
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllLogBucket(ctx context.Context, f func(*LogBucket) bool, resources []*LogBucket) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteLogBucket(ctx, res)
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

type deleteLogBucketOperation struct{}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createLogBucketOperation struct {
	response map[string]interface{}
}

func (op *createLogBucketOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createLogBucketOperation) do(ctx context.Context, r *LogBucket, c *Client) error {
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

	if _, err := c.GetLogBucket(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getLogBucketRaw(ctx context.Context, r *LogBucket) ([]byte, error) {

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

func (c *Client) logBucketDiffsForRawDesired(ctx context.Context, rawDesired *LogBucket, opts ...dcl.ApplyOption) (initial, desired *LogBucket, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *LogBucket
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*LogBucket); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected LogBucket, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetLogBucket(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a LogBucket resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve LogBucket resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that LogBucket resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeLogBucketDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for LogBucket: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for LogBucket: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractLogBucketFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeLogBucketInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for LogBucket: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeLogBucketDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for LogBucket: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffLogBucket(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeLogBucketInitialState(rawInitial, rawDesired *LogBucket) (*LogBucket, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeLogBucketDesiredState(rawDesired, rawInitial *LogBucket, opts ...dcl.ApplyOption) (*LogBucket, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &LogBucket{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.IsZeroValue(rawDesired.RetentionDays) || (dcl.IsEmptyValueIndirect(rawDesired.RetentionDays) && dcl.IsEmptyValueIndirect(rawInitial.RetentionDays)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.RetentionDays = rawInitial.RetentionDays
	} else {
		canonicalDesired.RetentionDays = rawDesired.RetentionDays
	}
	if dcl.BoolCanonicalize(rawDesired.Locked, rawInitial.Locked) {
		canonicalDesired.Locked = rawInitial.Locked
	} else {
		canonicalDesired.Locked = rawDesired.Locked
	}
	if dcl.NameToSelfLink(rawDesired.Parent, rawInitial.Parent) {
		canonicalDesired.Parent = rawInitial.Parent
	} else {
		canonicalDesired.Parent = rawDesired.Parent
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	return canonicalDesired, nil
}

func canonicalizeLogBucketNewState(c *Client, rawNew, rawDesired *LogBucket) (*LogBucket, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
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

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.RetentionDays) && dcl.IsEmptyValueIndirect(rawDesired.RetentionDays) {
		rawNew.RetentionDays = rawDesired.RetentionDays
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Locked) && dcl.IsEmptyValueIndirect(rawDesired.Locked) {
		rawNew.Locked = rawDesired.Locked
	} else {
		if dcl.BoolCanonicalize(rawDesired.Locked, rawNew.Locked) {
			rawNew.Locked = rawDesired.Locked
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.LifecycleState) && dcl.IsEmptyValueIndirect(rawDesired.LifecycleState) {
		rawNew.LifecycleState = rawDesired.LifecycleState
	} else {
	}

	rawNew.Parent = rawDesired.Parent

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
func diffLogBucket(c *Client, desired, actual *LogBucket, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogBucketUpdateBucketOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.RetentionDays, actual.RetentionDays, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogBucketUpdateBucketOperation")}, fn.AddNest("RetentionDays")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Locked, actual.Locked, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateLogBucketUpdateBucketOperation")}, fn.AddNest("Locked")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Parent, actual.Parent, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Parent")); len(ds) != 0 || err != nil {
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
func (r *LogBucket) urlNormalized() *LogBucket {
	normalized := dcl.Copy(*r).(LogBucket)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Parent = r.Parent
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *LogBucket) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateBucket" {
		fields := map[string]interface{}{
			"location": dcl.ValueOrEmptyString(nr.Location),
			"parent":   dcl.ValueOrEmptyString(nr.Parent),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("{{parent}}/locations/{{location}}/buckets/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the LogBucket resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *LogBucket) marshal(c *Client) ([]byte, error) {
	m, err := expandLogBucket(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling LogBucket: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalLogBucket decodes JSON responses into the LogBucket resource schema.
func unmarshalLogBucket(b []byte, c *Client, res *LogBucket) (*LogBucket, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapLogBucket(m, c, res)
}

func unmarshalMapLogBucket(m map[string]interface{}, c *Client, res *LogBucket) (*LogBucket, error) {

	flattened := flattenLogBucket(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandLogBucket expands LogBucket into a JSON request object.
func expandLogBucket(c *Client, f *LogBucket) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("%s/locations/%s/buckets/%s", f.Name, f.Parent, dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.RetentionDays; dcl.ValueShouldBeSent(v) {
		m["retentionDays"] = v
	}
	if v := f.Locked; dcl.ValueShouldBeSent(v) {
		m["locked"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Parent into parent: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["parent"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}

	return m, nil
}

// flattenLogBucket flattens LogBucket from a JSON request object into the
// LogBucket type.
func flattenLogBucket(c *Client, i interface{}, res *LogBucket) *LogBucket {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &LogBucket{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.RetentionDays = dcl.FlattenInteger(m["retentionDays"])
	resultRes.Locked = dcl.FlattenBool(m["locked"])
	resultRes.LifecycleState = flattenLogBucketLifecycleStateEnum(m["lifecycleState"])
	resultRes.Parent = dcl.FlattenString(m["parent"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// flattenLogBucketLifecycleStateEnumMap flattens the contents of LogBucketLifecycleStateEnum from a JSON
// response object.
func flattenLogBucketLifecycleStateEnumMap(c *Client, i interface{}, res *LogBucket) map[string]LogBucketLifecycleStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]LogBucketLifecycleStateEnum{}
	}

	if len(a) == 0 {
		return map[string]LogBucketLifecycleStateEnum{}
	}

	items := make(map[string]LogBucketLifecycleStateEnum)
	for k, item := range a {
		items[k] = *flattenLogBucketLifecycleStateEnum(item.(interface{}))
	}

	return items
}

// flattenLogBucketLifecycleStateEnumSlice flattens the contents of LogBucketLifecycleStateEnum from a JSON
// response object.
func flattenLogBucketLifecycleStateEnumSlice(c *Client, i interface{}, res *LogBucket) []LogBucketLifecycleStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []LogBucketLifecycleStateEnum{}
	}

	if len(a) == 0 {
		return []LogBucketLifecycleStateEnum{}
	}

	items := make([]LogBucketLifecycleStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenLogBucketLifecycleStateEnum(item.(interface{})))
	}

	return items
}

// flattenLogBucketLifecycleStateEnum asserts that an interface is a string, and returns a
// pointer to a *LogBucketLifecycleStateEnum with the same value as that string.
func flattenLogBucketLifecycleStateEnum(i interface{}) *LogBucketLifecycleStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return LogBucketLifecycleStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *LogBucket) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalLogBucket(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if nr.Location == nil && ncr.Location == nil {
			c.Config.Logger.Info("Both Location fields null - considering equal.")
		} else if nr.Location == nil || ncr.Location == nil {
			c.Config.Logger.Info("Only one Location field is null - considering unequal.")
			return false
		} else if *nr.Location != *ncr.Location {
			return false
		}
		if nr.Parent == nil && ncr.Parent == nil {
			c.Config.Logger.Info("Both Parent fields null - considering equal.")
		} else if nr.Parent == nil || ncr.Parent == nil {
			c.Config.Logger.Info("Only one Parent field is null - considering unequal.")
			return false
		} else if *nr.Parent != *ncr.Parent {
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

type logBucketDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         logBucketApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToLogBucketDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]logBucketDiff, error) {
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
	var diffs []logBucketDiff
	// For each operation name, create a logBucketDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := logBucketDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToLogBucketApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToLogBucketApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (logBucketApiOperation, error) {
	switch opName {

	case "updateLogBucketUpdateBucketOperation":
		return &updateLogBucketUpdateBucketOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractLogBucketFields(r *LogBucket) error {
	return nil
}

func postReadExtractLogBucketFields(r *LogBucket) error {
	return nil
}
