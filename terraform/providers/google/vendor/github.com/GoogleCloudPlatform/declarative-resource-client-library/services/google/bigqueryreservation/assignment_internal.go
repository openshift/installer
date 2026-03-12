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
package bigqueryreservation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *Assignment) validate() error {

	if err := dcl.Required(r, "assignee"); err != nil {
		return err
	}
	if err := dcl.Required(r, "jobType"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Reservation, "Reservation"); err != nil {
		return err
	}
	return nil
}
func (r *Assignment) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://bigqueryreservation.googleapis.com/v1/", params)
}

func (r *Assignment) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":     dcl.ValueOrEmptyString(nr.Project),
		"location":    dcl.ValueOrEmptyString(nr.Location),
		"reservation": dcl.ValueOrEmptyString(nr.Reservation),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments", nr.basePath(), userBasePath, params), nil
}

func (r *Assignment) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":     dcl.ValueOrEmptyString(nr.Project),
		"location":    dcl.ValueOrEmptyString(nr.Location),
		"reservation": dcl.ValueOrEmptyString(nr.Reservation),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments", nr.basePath(), userBasePath, params), nil

}

func (r *Assignment) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":     dcl.ValueOrEmptyString(nr.Project),
		"location":    dcl.ValueOrEmptyString(nr.Location),
		"reservation": dcl.ValueOrEmptyString(nr.Reservation),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments", nr.basePath(), userBasePath, params), nil

}

func (r *Assignment) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":     dcl.ValueOrEmptyString(nr.Project),
		"location":    dcl.ValueOrEmptyString(nr.Location),
		"reservation": dcl.ValueOrEmptyString(nr.Reservation),
		"name":        dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments/{{name}}", nr.basePath(), userBasePath, params), nil
}

// assignmentApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type assignmentApiOperation interface {
	do(context.Context, *Assignment, *Client) error
}

func (c *Client) listAssignmentRaw(ctx context.Context, r *Assignment, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != AssignmentMaxPage {
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

type listAssignmentOperation struct {
	Assignments []map[string]interface{} `json:"assignments"`
	Token       string                   `json:"nextPageToken"`
}

func (c *Client) listAssignment(ctx context.Context, r *Assignment, pageToken string, pageSize int32) ([]*Assignment, string, error) {
	b, err := c.listAssignmentRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listAssignmentOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Assignment
	for _, v := range m.Assignments {
		res, err := unmarshalMapAssignment(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		res.Reservation = r.Reservation
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllAssignment(ctx context.Context, f func(*Assignment) bool, resources []*Assignment) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteAssignment(ctx, res)
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

type deleteAssignmentOperation struct{}

func (op *deleteAssignmentOperation) do(ctx context.Context, r *Assignment, c *Client) error {
	r, err := c.GetAssignment(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Assignment not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetAssignment checking for existence. error: %v", err)
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
		return fmt.Errorf("failed to delete Assignment: %w", err)
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetAssignment(ctx, r)
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
type createAssignmentOperation struct {
	response map[string]interface{}
}

func (op *createAssignmentOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createAssignmentOperation) do(ctx context.Context, r *Assignment, c *Client) error {
	c.Config.Logger.InfoWithContextf(ctx, "Attempting to create %v", r)
	u, err := r.createURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	req, err := r.marshal(c)
	if err != nil {
		return err
	}
	if r.Name != nil {
		// Allowing creation to continue with Name set could result in a Assignment with the wrong Name.
		return fmt.Errorf("server-generated parameter Name was specified by user as %v, should be unspecified", dcl.ValueOrEmptyString(r.Name))
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

	// Include Name in URL substitution for initial GET request.
	m := op.response
	r.Name = dcl.SelfLinkToName(dcl.FlattenString(m["name"]))

	if _, err := c.GetAssignment(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getAssignmentRaw(ctx context.Context, r *Assignment) ([]byte, error) {

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

	b, err = dcl.ExtractElementFromList(b, "assignments", r.matcher(c))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Client) assignmentDiffsForRawDesired(ctx context.Context, rawDesired *Assignment, opts ...dcl.ApplyOption) (initial, desired *Assignment, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Assignment
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Assignment); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Assignment, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	if fetchState.Name == nil {
		// We cannot perform a get because of lack of information. We have to assume
		// that this is being created for the first time.
		desired, err := canonicalizeAssignmentDesiredState(rawDesired, nil)
		return nil, desired, nil, err
	}
	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetAssignment(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Assignment resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Assignment resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Assignment resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeAssignmentDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Assignment: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Assignment: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractAssignmentFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeAssignmentInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Assignment: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeAssignmentDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Assignment: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffAssignment(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeAssignmentInitialState(rawInitial, rawDesired *Assignment) (*Assignment, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeAssignmentDesiredState(rawDesired, rawInitial *Assignment, opts ...dcl.ApplyOption) (*Assignment, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &Assignment{}
	if dcl.IsZeroValue(rawDesired.Name) || (dcl.IsEmptyValueIndirect(rawDesired.Name) && dcl.IsEmptyValueIndirect(rawInitial.Name)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.IsZeroValue(rawDesired.Assignee) || (dcl.IsEmptyValueIndirect(rawDesired.Assignee) && dcl.IsEmptyValueIndirect(rawInitial.Assignee)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Assignee = rawInitial.Assignee
	} else {
		canonicalDesired.Assignee = rawDesired.Assignee
	}
	if dcl.IsZeroValue(rawDesired.JobType) || (dcl.IsEmptyValueIndirect(rawDesired.JobType) && dcl.IsEmptyValueIndirect(rawInitial.JobType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.JobType = rawInitial.JobType
	} else {
		canonicalDesired.JobType = rawDesired.JobType
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
	if dcl.NameToSelfLink(rawDesired.Reservation, rawInitial.Reservation) {
		canonicalDesired.Reservation = rawInitial.Reservation
	} else {
		canonicalDesired.Reservation = rawDesired.Reservation
	}
	return canonicalDesired, nil
}

func canonicalizeAssignmentNewState(c *Client, rawNew, rawDesired *Assignment) (*Assignment, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Assignee) && dcl.IsEmptyValueIndirect(rawDesired.Assignee) {
		rawNew.Assignee = rawDesired.Assignee
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.JobType) && dcl.IsEmptyValueIndirect(rawDesired.JobType) {
		rawNew.JobType = rawDesired.JobType
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	rawNew.Reservation = rawDesired.Reservation

	return rawNew, nil
}

// The differ returns a list of diffs, along with a list of operations that should be taken
// to remedy them. Right now, it does not attempt to consolidate operations - if several
// fields can be fixed with a patch update, it will perform the patch several times.
// Diffs on some fields will be ignored if the `desired` state has an empty (nil)
// value. This empty value indicates that the user does not care about the state for
// the field. Empty fields on the actual object will cause diffs.
// TODO(magic-modules-eng): for efficiency in some resources, add batching.
func diffAssignment(c *Client, desired, actual *Assignment, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Assignee, actual.Assignee, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Assignee")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JobType, actual.JobType, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JobType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Reservation, actual.Reservation, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Reservation")); len(ds) != 0 || err != nil {
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
func (r *Assignment) urlNormalized() *Assignment {
	normalized := dcl.Copy(*r).(Assignment)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Assignee = dcl.SelfLinkToName(r.Assignee)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.Reservation = dcl.SelfLinkToName(r.Reservation)
	return &normalized
}

func (r *Assignment) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Assignment resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Assignment) marshal(c *Client) ([]byte, error) {
	m, err := expandAssignment(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Assignment: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalAssignment decodes JSON responses into the Assignment resource schema.
func unmarshalAssignment(b []byte, c *Client, res *Assignment) (*Assignment, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapAssignment(m, c, res)
}

func unmarshalMapAssignment(m map[string]interface{}, c *Client, res *Assignment) (*Assignment, error) {

	flattened := flattenAssignment(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandAssignment expands Assignment into a JSON request object.
func expandAssignment(c *Client, f *Assignment) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Assignee; dcl.ValueShouldBeSent(v) {
		m["assignee"] = v
	}
	if v := f.JobType; dcl.ValueShouldBeSent(v) {
		m["jobType"] = v
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
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Reservation into reservation: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["reservation"] = v
	}

	return m, nil
}

// flattenAssignment flattens Assignment from a JSON request object into the
// Assignment type.
func flattenAssignment(c *Client, i interface{}, res *Assignment) *Assignment {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Assignment{}
	resultRes.Name = dcl.SelfLinkToName(dcl.FlattenString(m["name"]))
	resultRes.Assignee = dcl.FlattenString(m["assignee"])
	resultRes.JobType = flattenAssignmentJobTypeEnum(m["jobType"])
	resultRes.State = flattenAssignmentStateEnum(m["state"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.Reservation = dcl.FlattenString(m["reservation"])

	return resultRes
}

// flattenAssignmentJobTypeEnumMap flattens the contents of AssignmentJobTypeEnum from a JSON
// response object.
func flattenAssignmentJobTypeEnumMap(c *Client, i interface{}, res *Assignment) map[string]AssignmentJobTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]AssignmentJobTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]AssignmentJobTypeEnum{}
	}

	items := make(map[string]AssignmentJobTypeEnum)
	for k, item := range a {
		items[k] = *flattenAssignmentJobTypeEnum(item.(interface{}))
	}

	return items
}

// flattenAssignmentJobTypeEnumSlice flattens the contents of AssignmentJobTypeEnum from a JSON
// response object.
func flattenAssignmentJobTypeEnumSlice(c *Client, i interface{}, res *Assignment) []AssignmentJobTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []AssignmentJobTypeEnum{}
	}

	if len(a) == 0 {
		return []AssignmentJobTypeEnum{}
	}

	items := make([]AssignmentJobTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenAssignmentJobTypeEnum(item.(interface{})))
	}

	return items
}

// flattenAssignmentJobTypeEnum asserts that an interface is a string, and returns a
// pointer to a *AssignmentJobTypeEnum with the same value as that string.
func flattenAssignmentJobTypeEnum(i interface{}) *AssignmentJobTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return AssignmentJobTypeEnumRef(s)
}

// flattenAssignmentStateEnumMap flattens the contents of AssignmentStateEnum from a JSON
// response object.
func flattenAssignmentStateEnumMap(c *Client, i interface{}, res *Assignment) map[string]AssignmentStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]AssignmentStateEnum{}
	}

	if len(a) == 0 {
		return map[string]AssignmentStateEnum{}
	}

	items := make(map[string]AssignmentStateEnum)
	for k, item := range a {
		items[k] = *flattenAssignmentStateEnum(item.(interface{}))
	}

	return items
}

// flattenAssignmentStateEnumSlice flattens the contents of AssignmentStateEnum from a JSON
// response object.
func flattenAssignmentStateEnumSlice(c *Client, i interface{}, res *Assignment) []AssignmentStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []AssignmentStateEnum{}
	}

	if len(a) == 0 {
		return []AssignmentStateEnum{}
	}

	items := make([]AssignmentStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenAssignmentStateEnum(item.(interface{})))
	}

	return items
}

// flattenAssignmentStateEnum asserts that an interface is a string, and returns a
// pointer to a *AssignmentStateEnum with the same value as that string.
func flattenAssignmentStateEnum(i interface{}) *AssignmentStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return AssignmentStateEnumRef(s)
}

type assignmentDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         assignmentApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToAssignmentDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]assignmentDiff, error) {
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
	var diffs []assignmentDiff
	// For each operation name, create a assignmentDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := assignmentDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToAssignmentApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToAssignmentApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (assignmentApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractAssignmentFields(r *Assignment) error {
	vProject, err := dcl.ValueFromRegexOnField("Project", r.Project, r.Reservation, "projects/([a-z0-9A-Z-]*)/locations/.*")
	if err != nil {
		return err
	}
	r.Project = vProject
	vLocation, err := dcl.ValueFromRegexOnField("Location", r.Location, r.Reservation, "projects/.*/locations/([a-z0-9A-Z-]*)/reservations/.*")
	if err != nil {
		return err
	}
	r.Location = vLocation
	return nil
}

func postReadExtractAssignmentFields(r *Assignment) error {
	return nil
}
