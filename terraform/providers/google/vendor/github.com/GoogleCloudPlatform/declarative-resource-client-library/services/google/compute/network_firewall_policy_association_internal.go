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

func (r *NetworkFirewallPolicyAssociation) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "attachmentTarget"); err != nil {
		return err
	}
	if err := dcl.Required(r, "firewallPolicy"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	return nil
}
func (r *NetworkFirewallPolicyAssociation) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *NetworkFirewallPolicyAssociation) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"name":           dcl.ValueOrEmptyString(nr.Name),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}/getAssociation?name={{name}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}/getAssociation?name={{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *NetworkFirewallPolicyAssociation) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}", nr.basePath(), userBasePath, params), nil

}

func (r *NetworkFirewallPolicyAssociation) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}/addAssociation", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}/addAssociation", nr.basePath(), userBasePath, params), nil

}

func (r *NetworkFirewallPolicyAssociation) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":        dcl.ValueOrEmptyString(nr.Project),
		"location":       dcl.ValueOrEmptyString(nr.Location),
		"firewallPolicy": dcl.ValueOrEmptyString(nr.FirewallPolicy),
		"name":           dcl.ValueOrEmptyString(nr.Name),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/firewallPolicies/{{firewallPolicy}}/removeAssociation?name={{name}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/firewallPolicies/{{firewallPolicy}}/removeAssociation?name={{name}}", nr.basePath(), userBasePath, params), nil
}

// networkFirewallPolicyAssociationApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type networkFirewallPolicyAssociationApiOperation interface {
	do(context.Context, *NetworkFirewallPolicyAssociation, *Client) error
}

func (c *Client) listNetworkFirewallPolicyAssociationRaw(ctx context.Context, r *NetworkFirewallPolicyAssociation, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != NetworkFirewallPolicyAssociationMaxPage {
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

type listNetworkFirewallPolicyAssociationOperation struct {
	Associations []map[string]interface{} `json:"associations"`
	Token        string                   `json:"nextPageToken"`
}

func (c *Client) listNetworkFirewallPolicyAssociation(ctx context.Context, r *NetworkFirewallPolicyAssociation, pageToken string, pageSize int32) ([]*NetworkFirewallPolicyAssociation, string, error) {
	b, err := c.listNetworkFirewallPolicyAssociationRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listNetworkFirewallPolicyAssociationOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*NetworkFirewallPolicyAssociation
	for _, v := range m.Associations {
		res, err := unmarshalMapNetworkFirewallPolicyAssociation(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		res.FirewallPolicy = r.FirewallPolicy
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllNetworkFirewallPolicyAssociation(ctx context.Context, f func(*NetworkFirewallPolicyAssociation) bool, resources []*NetworkFirewallPolicyAssociation) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteNetworkFirewallPolicyAssociation(ctx, res)
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

type deleteNetworkFirewallPolicyAssociationOperation struct{}

func (op *deleteNetworkFirewallPolicyAssociationOperation) do(ctx context.Context, r *NetworkFirewallPolicyAssociation, c *Client) error {
	r, err := c.GetNetworkFirewallPolicyAssociation(ctx, r)
	if err != nil {
		if dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.InfoWithContextf(ctx, "NetworkFirewallPolicyAssociation not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetNetworkFirewallPolicyAssociation checking for existence. error: %v", err)
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
		_, err := c.GetNetworkFirewallPolicyAssociation(ctx, r)
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
type createNetworkFirewallPolicyAssociationOperation struct {
	response map[string]interface{}
}

func (op *createNetworkFirewallPolicyAssociationOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createNetworkFirewallPolicyAssociationOperation) do(ctx context.Context, r *NetworkFirewallPolicyAssociation, c *Client) error {
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

	if _, err := c.GetNetworkFirewallPolicyAssociation(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getNetworkFirewallPolicyAssociationRaw(ctx context.Context, r *NetworkFirewallPolicyAssociation) ([]byte, error) {

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

func (c *Client) networkFirewallPolicyAssociationDiffsForRawDesired(ctx context.Context, rawDesired *NetworkFirewallPolicyAssociation, opts ...dcl.ApplyOption) (initial, desired *NetworkFirewallPolicyAssociation, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *NetworkFirewallPolicyAssociation
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*NetworkFirewallPolicyAssociation); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected NetworkFirewallPolicyAssociation, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetNetworkFirewallPolicyAssociation(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFoundOrCode(err, 400) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a NetworkFirewallPolicyAssociation resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve NetworkFirewallPolicyAssociation resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that NetworkFirewallPolicyAssociation resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeNetworkFirewallPolicyAssociationDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for NetworkFirewallPolicyAssociation: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for NetworkFirewallPolicyAssociation: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractNetworkFirewallPolicyAssociationFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeNetworkFirewallPolicyAssociationInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for NetworkFirewallPolicyAssociation: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeNetworkFirewallPolicyAssociationDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for NetworkFirewallPolicyAssociation: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffNetworkFirewallPolicyAssociation(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeNetworkFirewallPolicyAssociationInitialState(rawInitial, rawDesired *NetworkFirewallPolicyAssociation) (*NetworkFirewallPolicyAssociation, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeNetworkFirewallPolicyAssociationDesiredState(rawDesired, rawInitial *NetworkFirewallPolicyAssociation, opts ...dcl.ApplyOption) (*NetworkFirewallPolicyAssociation, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &NetworkFirewallPolicyAssociation{}
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
	if dcl.IsZeroValue(rawDesired.FirewallPolicy) || (dcl.IsEmptyValueIndirect(rawDesired.FirewallPolicy) && dcl.IsEmptyValueIndirect(rawInitial.FirewallPolicy)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.FirewallPolicy = rawInitial.FirewallPolicy
	} else {
		canonicalDesired.FirewallPolicy = rawDesired.FirewallPolicy
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeNetworkFirewallPolicyAssociationNewState(c *Client, rawNew, rawDesired *NetworkFirewallPolicyAssociation) (*NetworkFirewallPolicyAssociation, error) {

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
	}

	if dcl.IsEmptyValueIndirect(rawNew.ShortName) && dcl.IsEmptyValueIndirect(rawDesired.ShortName) {
		rawNew.ShortName = rawDesired.ShortName
	} else {
		if dcl.StringCanonicalize(rawDesired.ShortName, rawNew.ShortName) {
			rawNew.ShortName = rawDesired.ShortName
		}
	}

	rawNew.Location = rawDesired.Location

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
func diffNetworkFirewallPolicyAssociation(c *Client, desired, actual *NetworkFirewallPolicyAssociation, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
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
func (r *NetworkFirewallPolicyAssociation) urlNormalized() *NetworkFirewallPolicyAssociation {
	normalized := dcl.Copy(*r).(NetworkFirewallPolicyAssociation)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.AttachmentTarget = dcl.SelfLinkToName(r.AttachmentTarget)
	normalized.FirewallPolicy = dcl.SelfLinkToName(r.FirewallPolicy)
	normalized.ShortName = dcl.SelfLinkToName(r.ShortName)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *NetworkFirewallPolicyAssociation) updateURL(userBasePath, updateName string) (string, error) {
	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the NetworkFirewallPolicyAssociation resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *NetworkFirewallPolicyAssociation) marshal(c *Client) ([]byte, error) {
	m, err := expandNetworkFirewallPolicyAssociation(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling NetworkFirewallPolicyAssociation: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalNetworkFirewallPolicyAssociation decodes JSON responses into the NetworkFirewallPolicyAssociation resource schema.
func unmarshalNetworkFirewallPolicyAssociation(b []byte, c *Client, res *NetworkFirewallPolicyAssociation) (*NetworkFirewallPolicyAssociation, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapNetworkFirewallPolicyAssociation(m, c, res)
}

func unmarshalMapNetworkFirewallPolicyAssociation(m map[string]interface{}, c *Client, res *NetworkFirewallPolicyAssociation) (*NetworkFirewallPolicyAssociation, error) {

	flattened := flattenNetworkFirewallPolicyAssociation(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandNetworkFirewallPolicyAssociation expands NetworkFirewallPolicyAssociation into a JSON request object.
func expandNetworkFirewallPolicyAssociation(c *Client, f *NetworkFirewallPolicyAssociation) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.AttachmentTarget; dcl.ValueShouldBeSent(v) {
		m["attachmentTarget"] = v
	}
	if v := f.FirewallPolicy; dcl.ValueShouldBeSent(v) {
		m["firewallPolicyId"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenNetworkFirewallPolicyAssociation flattens NetworkFirewallPolicyAssociation from a JSON request object into the
// NetworkFirewallPolicyAssociation type.
func flattenNetworkFirewallPolicyAssociation(c *Client, i interface{}, res *NetworkFirewallPolicyAssociation) *NetworkFirewallPolicyAssociation {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &NetworkFirewallPolicyAssociation{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.AttachmentTarget = dcl.FlattenString(m["attachmentTarget"])
	resultRes.FirewallPolicy = dcl.FlattenString(m["firewallPolicyId"])
	resultRes.ShortName = dcl.FlattenString(m["shortName"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.Project = dcl.FlattenString(m["project"])

	return resultRes
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *NetworkFirewallPolicyAssociation) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalNetworkFirewallPolicyAssociation(b, c, r)
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

type networkFirewallPolicyAssociationDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         networkFirewallPolicyAssociationApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToNetworkFirewallPolicyAssociationDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]networkFirewallPolicyAssociationDiff, error) {
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
	var diffs []networkFirewallPolicyAssociationDiff
	// For each operation name, create a networkFirewallPolicyAssociationDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := networkFirewallPolicyAssociationDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToNetworkFirewallPolicyAssociationApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToNetworkFirewallPolicyAssociationApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (networkFirewallPolicyAssociationApiOperation, error) {
	switch opName {

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractNetworkFirewallPolicyAssociationFields(r *NetworkFirewallPolicyAssociation) error {
	return nil
}

func postReadExtractNetworkFirewallPolicyAssociationFields(r *NetworkFirewallPolicyAssociation) error {
	return nil
}
