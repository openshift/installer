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
package compute

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

func (r *FirewallPolicyAssociation) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "attachmentTarget"); err != nil {
		return err
	}
	if err := dcl.Required(r, "firewallPolicy"); err != nil {
		return err
	}
	return nil
}
func (r *FirewallPolicyAssociation) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *FirewallPolicyAssociation) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"name":           dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}/getAssociation?name={{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *FirewallPolicyAssociation) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
	}
	return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}", nr.basePath(), userBasePath, params), nil

}

func (r *FirewallPolicyAssociation) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
	}
	return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}/addAssociation", nr.basePath(), userBasePath, params), nil

}

func (r *FirewallPolicyAssociation) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"name":           dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("locations/global/firewallPolicies/{{firewallPolicy}}/removeAssociation?name={{name}}", nr.basePath(), userBasePath, params), nil
}

// firewallPolicyAssociationApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type firewallPolicyAssociationApiOperation interface {
	do(context.Context, *FirewallPolicyAssociation, *Client) error
}

func (c *Client) listFirewallPolicyAssociationRaw(ctx context.Context, r *FirewallPolicyAssociation, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != FirewallPolicyAssociationMaxPage {
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

type listFirewallPolicyAssociationOperation struct {
	Associations []map[string]interface{} `json:"associations"`
	Token        string                   `json:"nextPageToken"`
}

func (c *Client) listFirewallPolicyAssociation(ctx context.Context, r *FirewallPolicyAssociation, pageToken string, pageSize int32) ([]*FirewallPolicyAssociation, string, error) {
	b, err := c.listFirewallPolicyAssociationRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listFirewallPolicyAssociationOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*FirewallPolicyAssociation
	for _, v := range m.Associations {
		res, err := unmarshalMapFirewallPolicyAssociation(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.FirewallPolicy = r.FirewallPolicy
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllFirewallPolicyAssociation(ctx context.Context, f func(*FirewallPolicyAssociation) bool, resources []*FirewallPolicyAssociation) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteFirewallPolicyAssociation(ctx, res)
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

type deleteFirewallPolicyAssociationOperation struct{}

func (op *deleteFirewallPolicyAssociationOperation) do(ctx context.Context, r *FirewallPolicyAssociation, c *Client) error {
	r, err := c.GetFirewallPolicyAssociation(ctx, r)
	if err != nil {
		if dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.InfoWithContextf(ctx, "FirewallPolicyAssociation not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetFirewallPolicyAssociation checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, body, c.Config.RetryProvider)
	if err != nil {
		return err
	}

	// wait for object to be deleted.
	var o operations.ComputeOperation
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
		_, err := c.GetFirewallPolicyAssociation(ctx, r)
		if dcl.IsNotFoundOrCode(err, 400) {
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
type createFirewallPolicyAssociationOperation struct {
	response map[string]interface{}
}

func (op *createFirewallPolicyAssociationOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createFirewallPolicyAssociationOperation) do(ctx context.Context, r *FirewallPolicyAssociation, c *Client) error {
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
	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetFirewallPolicyAssociation(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getFirewallPolicyAssociationRaw(ctx context.Context, r *FirewallPolicyAssociation) ([]byte, error) {

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

func (c *Client) firewallPolicyAssociationDiffsForRawDesired(ctx context.Context, rawDesired *FirewallPolicyAssociation, opts ...dcl.ApplyOption) (initial, desired *FirewallPolicyAssociation, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *FirewallPolicyAssociation
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*FirewallPolicyAssociation); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected FirewallPolicyAssociation, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetFirewallPolicyAssociation(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a FirewallPolicyAssociation resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve FirewallPolicyAssociation resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that FirewallPolicyAssociation resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeFirewallPolicyAssociationDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for FirewallPolicyAssociation: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for FirewallPolicyAssociation: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractFirewallPolicyAssociationFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeFirewallPolicyAssociationInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for FirewallPolicyAssociation: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeFirewallPolicyAssociationDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for FirewallPolicyAssociation: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffFirewallPolicyAssociation(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeFirewallPolicyAssociationInitialState(rawInitial, rawDesired *FirewallPolicyAssociation) (*FirewallPolicyAssociation, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeFirewallPolicyAssociationDesiredState(rawDesired, rawInitial *FirewallPolicyAssociation, opts ...dcl.ApplyOption) (*FirewallPolicyAssociation, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &FirewallPolicyAssociation{}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.IsZeroValue(rawDesired.AttachmentTarget) || (dcl.IsEmptyValueIndirect(rawDesired.AttachmentTarget) && dcl.IsEmptyValueIndirect(rawInitial.AttachmentTarget)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.AttachmentTarget = rawInitial.AttachmentTarget
	} else {
		canonicalDesired.AttachmentTarget = rawDesired.AttachmentTarget
	}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.FirewallPolicy, rawInitial.FirewallPolicy) {
		canonicalDesired.FirewallPolicy = rawInitial.FirewallPolicy
	} else {
		canonicalDesired.FirewallPolicy = rawDesired.FirewallPolicy
	}
	return canonicalDesired, nil
}

func canonicalizeFirewallPolicyAssociationNewState(c *Client, rawNew, rawDesired *FirewallPolicyAssociation) (*FirewallPolicyAssociation, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.AttachmentTarget) && dcl.IsEmptyValueIndirect(rawDesired.AttachmentTarget) {
		rawNew.AttachmentTarget = rawDesired.AttachmentTarget
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.FirewallPolicy) && dcl.IsEmptyValueIndirect(rawDesired.FirewallPolicy) {
		rawNew.FirewallPolicy = rawDesired.FirewallPolicy
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.FirewallPolicy, rawNew.FirewallPolicy) {
			rawNew.FirewallPolicy = rawDesired.FirewallPolicy
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.ShortName) && dcl.IsEmptyValueIndirect(rawDesired.ShortName) {
		rawNew.ShortName = rawDesired.ShortName
	} else {
		if dcl.StringCanonicalize(rawDesired.ShortName, rawNew.ShortName) {
			rawNew.ShortName = rawDesired.ShortName
		}
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
func diffFirewallPolicyAssociation(c *Client, desired, actual *FirewallPolicyAssociation, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.AttachmentTarget, actual.AttachmentTarget, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AttachmentTarget")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FirewallPolicy, actual.FirewallPolicy, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FirewallPolicyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ShortName, actual.ShortName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ShortName")); len(ds) != 0 || err != nil {
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
func (r *FirewallPolicyAssociation) urlNormalized() *FirewallPolicyAssociation {
	normalized := dcl.Copy(*r).(FirewallPolicyAssociation)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.AttachmentTarget = dcl.SelfLinkToName(r.AttachmentTarget)
	normalized.FirewallPolicy = dcl.SelfLinkToName(r.FirewallPolicy)
	normalized.ShortName = dcl.SelfLinkToName(r.ShortName)
	return &normalized
}

func (r *FirewallPolicyAssociation) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the FirewallPolicyAssociation resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *FirewallPolicyAssociation) marshal(c *Client) ([]byte, error) {
	m, err := expandFirewallPolicyAssociation(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling FirewallPolicyAssociation: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalFirewallPolicyAssociation decodes JSON responses into the FirewallPolicyAssociation resource schema.
func unmarshalFirewallPolicyAssociation(b []byte, c *Client, res *FirewallPolicyAssociation) (*FirewallPolicyAssociation, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapFirewallPolicyAssociation(m, c, res)
}

func unmarshalMapFirewallPolicyAssociation(m map[string]interface{}, c *Client, res *FirewallPolicyAssociation) (*FirewallPolicyAssociation, error) {

	flattened := flattenFirewallPolicyAssociation(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandFirewallPolicyAssociation expands FirewallPolicyAssociation into a JSON request object.
func expandFirewallPolicyAssociation(c *Client, f *FirewallPolicyAssociation) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.AttachmentTarget; dcl.ValueShouldBeSent(v) {
		m["attachmentTarget"] = v
	}
	if v, err := dcl.DeriveField("locations/global/firewallPolicies/%s", f.FirewallPolicy, dcl.SelfLinkToName(f.FirewallPolicy)); err != nil {
		return nil, fmt.Errorf("error expanding FirewallPolicy into firewallPolicyId: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["firewallPolicyId"] = v
	}

	return m, nil
}

// flattenFirewallPolicyAssociation flattens FirewallPolicyAssociation from a JSON request object into the
// FirewallPolicyAssociation type.
func flattenFirewallPolicyAssociation(c *Client, i interface{}, res *FirewallPolicyAssociation) *FirewallPolicyAssociation {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &FirewallPolicyAssociation{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.AttachmentTarget = dcl.FlattenString(m["attachmentTarget"])
	resultRes.FirewallPolicy = dcl.FlattenString(m["firewallPolicyId"])
	resultRes.ShortName = dcl.FlattenString(m["shortName"])

	return resultRes
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *FirewallPolicyAssociation) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalFirewallPolicyAssociation(b, c, r)
		if err != nil {
			c.Config.Logger.Warning("failed to unmarshal provided resource in matcher.")
			return false
		}
		nr := r.urlNormalized()
		ncr := cr.urlNormalized()
		c.Config.Logger.Infof("looking for %v\nin %v", nr, ncr)

		if nr.FirewallPolicy == nil && ncr.FirewallPolicy == nil {
			c.Config.Logger.Info("Both FirewallPolicy fields null - considering equal.")
		} else if nr.FirewallPolicy == nil || ncr.FirewallPolicy == nil {
			c.Config.Logger.Info("Only one FirewallPolicy field is null - considering unequal.")
			return false
		} else if *nr.FirewallPolicy != *ncr.FirewallPolicy {
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

type firewallPolicyAssociationDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         firewallPolicyAssociationApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToFirewallPolicyAssociationDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]firewallPolicyAssociationDiff, error) {
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
	var diffs []firewallPolicyAssociationDiff
	// For each operation name, create a firewallPolicyAssociationDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := firewallPolicyAssociationDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToFirewallPolicyAssociationApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToFirewallPolicyAssociationApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (firewallPolicyAssociationApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractFirewallPolicyAssociationFields(r *FirewallPolicyAssociation) error {
	return nil
}

func postReadExtractFirewallPolicyAssociationFields(r *FirewallPolicyAssociation) error {
	return nil
}
