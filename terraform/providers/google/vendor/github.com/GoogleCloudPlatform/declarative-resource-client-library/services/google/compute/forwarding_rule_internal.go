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

func (r *ForwardingRule) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	return nil
}
func (r *ForwardingRuleMetadataFilter) validate() error {
	if err := dcl.Required(r, "filterMatchCriteria"); err != nil {
		return err
	}
	if err := dcl.Required(r, "filterLabel"); err != nil {
		return err
	}
	return nil
}
func (r *ForwardingRuleMetadataFilterFilterLabel) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "value"); err != nil {
		return err
	}
	return nil
}
func (r *ForwardingRuleServiceDirectoryRegistrations) validate() error {
	return nil
}
func (r *ForwardingRule) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *ForwardingRule) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/forwardingRules/{{name}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/forwardingRules/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *ForwardingRule) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/forwardingRules", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/forwardingRules", nr.basePath(), userBasePath, params), nil

}

func (r *ForwardingRule) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/forwardingRules", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/forwardingRules", nr.basePath(), userBasePath, params), nil

}

func (r *ForwardingRule) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	if dcl.IsRegion(nr.Location) {
		return dcl.URL("projects/{{project}}/regions/{{location}}/forwardingRules/{{name}}", nr.basePath(), userBasePath, params), nil
	}

	return dcl.URL("projects/{{project}}/global/forwardingRules/{{name}}", nr.basePath(), userBasePath, params), nil
}

// forwardingRuleApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type forwardingRuleApiOperation interface {
	do(context.Context, *ForwardingRule, *Client) error
}

// newUpdateForwardingRuleSetLabelsRequest creates a request for an
// ForwardingRule resource's setLabels update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateForwardingRuleSetLabelsRequest(ctx context.Context, f *ForwardingRule, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	b, err := c.getForwardingRuleRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawLabelFingerprint, err := dcl.GetMapEntry(
		m,
		[]string{"labelFingerprint"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["labelFingerprint"] = rawLabelFingerprint.(string)
	}
	return req, nil
}

// marshalUpdateForwardingRuleSetLabelsRequest converts the update into
// the final JSON request body.
func marshalUpdateForwardingRuleSetLabelsRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateForwardingRuleSetLabelsOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateForwardingRuleSetLabelsOperation) do(ctx context.Context, r *ForwardingRule, c *Client) error {
	_, err := c.GetForwardingRule(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "setLabels")
	if err != nil {
		return err
	}

	req, err := newUpdateForwardingRuleSetLabelsRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateForwardingRuleSetLabelsRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

// newUpdateForwardingRuleSetTargetRequest creates a request for an
// ForwardingRule resource's setTarget update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateForwardingRuleSetTargetRequest(ctx context.Context, f *ForwardingRule, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Target; !dcl.IsEmptyValueIndirect(v) {
		req["target"] = v
	}
	return req, nil
}

// marshalUpdateForwardingRuleSetTargetRequest converts the update into
// the final JSON request body.
func marshalUpdateForwardingRuleSetTargetRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateForwardingRuleSetTargetOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateForwardingRuleSetTargetOperation) do(ctx context.Context, r *ForwardingRule, c *Client) error {
	_, err := c.GetForwardingRule(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "setTarget")
	if err != nil {
		return err
	}

	req, err := newUpdateForwardingRuleSetTargetRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateForwardingRuleSetTargetRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

// newUpdateForwardingRuleUpdateRequest creates a request for an
// ForwardingRule resource's update update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateForwardingRuleUpdateRequest(ctx context.Context, f *ForwardingRule, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := dcl.FalseToNil(f.AllowGlobalAccess); err != nil {
		return nil, fmt.Errorf("error expanding AllowGlobalAccess into allowGlobalAccess: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["allowGlobalAccess"] = v
	}
	return req, nil
}

// marshalUpdateForwardingRuleUpdateRequest converts the update into
// the final JSON request body.
func marshalUpdateForwardingRuleUpdateRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateForwardingRuleUpdateOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateForwardingRuleUpdateOperation) do(ctx context.Context, r *ForwardingRule, c *Client) error {
	_, err := c.GetForwardingRule(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "update")
	if err != nil {
		return err
	}

	req, err := newUpdateForwardingRuleUpdateRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateForwardingRuleUpdateRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.ComputeOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listForwardingRuleRaw(ctx context.Context, r *ForwardingRule, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ForwardingRuleMaxPage {
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

type listForwardingRuleOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listForwardingRule(ctx context.Context, r *ForwardingRule, pageToken string, pageSize int32) ([]*ForwardingRule, string, error) {
	b, err := c.listForwardingRuleRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listForwardingRuleOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*ForwardingRule
	for _, v := range m.Items {
		res, err := unmarshalMapForwardingRule(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllForwardingRule(ctx context.Context, f func(*ForwardingRule) bool, resources []*ForwardingRule) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteForwardingRule(ctx, res)
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

type deleteForwardingRuleOperation struct{}

func (op *deleteForwardingRuleOperation) do(ctx context.Context, r *ForwardingRule, c *Client) error {
	r, err := c.GetForwardingRule(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "ForwardingRule not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetForwardingRule checking for existence. error: %v", err)
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
		_, err := c.GetForwardingRule(ctx, r)
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
type createForwardingRuleOperation struct {
	response map[string]interface{}
}

func (op *createForwardingRuleOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createForwardingRuleOperation) do(ctx context.Context, r *ForwardingRule, c *Client) error {
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

	if _, err := c.GetForwardingRule(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getForwardingRuleRaw(ctx context.Context, r *ForwardingRule) ([]byte, error) {

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

func (c *Client) forwardingRuleDiffsForRawDesired(ctx context.Context, rawDesired *ForwardingRule, opts ...dcl.ApplyOption) (initial, desired *ForwardingRule, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *ForwardingRule
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*ForwardingRule); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected ForwardingRule, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetForwardingRule(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a ForwardingRule resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve ForwardingRule resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that ForwardingRule resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeForwardingRuleDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for ForwardingRule: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for ForwardingRule: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractForwardingRuleFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeForwardingRuleInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for ForwardingRule: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeForwardingRuleDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for ForwardingRule: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffForwardingRule(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeForwardingRuleInitialState(rawInitial, rawDesired *ForwardingRule) (*ForwardingRule, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeForwardingRuleDesiredState(rawDesired, rawInitial *ForwardingRule, opts ...dcl.ApplyOption) (*ForwardingRule, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &ForwardingRule{}
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.BoolCanonicalize(rawDesired.AllPorts, rawInitial.AllPorts) {
		canonicalDesired.AllPorts = rawInitial.AllPorts
	} else {
		canonicalDesired.AllPorts = rawDesired.AllPorts
	}
	if dcl.BoolCanonicalize(rawDesired.AllowGlobalAccess, rawInitial.AllowGlobalAccess) {
		canonicalDesired.AllowGlobalAccess = rawInitial.AllowGlobalAccess
	} else {
		canonicalDesired.AllowGlobalAccess = rawDesired.AllowGlobalAccess
	}
	if dcl.StringCanonicalize(rawDesired.BackendService, rawInitial.BackendService) {
		canonicalDesired.BackendService = rawInitial.BackendService
	} else {
		canonicalDesired.BackendService = rawDesired.BackendService
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if canonicalizeIPAddressToReference(rawDesired.IPAddress, rawInitial.IPAddress) {
		canonicalDesired.IPAddress = rawInitial.IPAddress
	} else {
		canonicalDesired.IPAddress = rawDesired.IPAddress
	}
	if dcl.IsZeroValue(rawDesired.IPProtocol) || (dcl.IsEmptyValueIndirect(rawDesired.IPProtocol) && dcl.IsEmptyValueIndirect(rawInitial.IPProtocol)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.IPProtocol = rawInitial.IPProtocol
	} else {
		canonicalDesired.IPProtocol = rawDesired.IPProtocol
	}
	if dcl.IsZeroValue(rawDesired.IPVersion) || (dcl.IsEmptyValueIndirect(rawDesired.IPVersion) && dcl.IsEmptyValueIndirect(rawInitial.IPVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.IPVersion = rawInitial.IPVersion
	} else {
		canonicalDesired.IPVersion = rawDesired.IPVersion
	}
	if dcl.BoolCanonicalize(rawDesired.IsMirroringCollector, rawInitial.IsMirroringCollector) {
		canonicalDesired.IsMirroringCollector = rawInitial.IsMirroringCollector
	} else {
		canonicalDesired.IsMirroringCollector = rawDesired.IsMirroringCollector
	}
	if dcl.IsZeroValue(rawDesired.LoadBalancingScheme) || (dcl.IsEmptyValueIndirect(rawDesired.LoadBalancingScheme) && dcl.IsEmptyValueIndirect(rawInitial.LoadBalancingScheme)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.LoadBalancingScheme = rawInitial.LoadBalancingScheme
	} else {
		canonicalDesired.LoadBalancingScheme = rawDesired.LoadBalancingScheme
	}
	canonicalDesired.MetadataFilter = canonicalizeForwardingRuleMetadataFilterSlice(rawDesired.MetadataFilter, rawInitial.MetadataFilter, opts...)
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Network, rawInitial.Network) {
		canonicalDesired.Network = rawInitial.Network
	} else {
		canonicalDesired.Network = rawDesired.Network
	}
	if dcl.IsZeroValue(rawDesired.NetworkTier) || (dcl.IsEmptyValueIndirect(rawDesired.NetworkTier) && dcl.IsEmptyValueIndirect(rawInitial.NetworkTier)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.NetworkTier = rawInitial.NetworkTier
	} else {
		canonicalDesired.NetworkTier = rawDesired.NetworkTier
	}
	if canonicalizePortRange(rawDesired.PortRange, rawInitial.PortRange) {
		canonicalDesired.PortRange = rawInitial.PortRange
	} else {
		canonicalDesired.PortRange = rawDesired.PortRange
	}
	if dcl.StringArrayCanonicalize(rawDesired.Ports, rawInitial.Ports) {
		canonicalDesired.Ports = rawInitial.Ports
	} else {
		canonicalDesired.Ports = rawDesired.Ports
	}
	if dcl.StringCanonicalize(rawDesired.ServiceLabel, rawInitial.ServiceLabel) {
		canonicalDesired.ServiceLabel = rawInitial.ServiceLabel
	} else {
		canonicalDesired.ServiceLabel = rawDesired.ServiceLabel
	}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Subnetwork, rawInitial.Subnetwork) {
		canonicalDesired.Subnetwork = rawInitial.Subnetwork
	} else {
		canonicalDesired.Subnetwork = rawDesired.Subnetwork
	}
	if dcl.StringCanonicalize(rawDesired.Target, rawInitial.Target) {
		canonicalDesired.Target = rawInitial.Target
	} else {
		canonicalDesired.Target = rawDesired.Target
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
	canonicalDesired.ServiceDirectoryRegistrations = canonicalizeForwardingRuleServiceDirectoryRegistrationsSlice(rawDesired.ServiceDirectoryRegistrations, rawInitial.ServiceDirectoryRegistrations, opts...)
	if dcl.StringArrayCanonicalize(rawDesired.SourceIPRanges, rawInitial.SourceIPRanges) {
		canonicalDesired.SourceIPRanges = rawInitial.SourceIPRanges
	} else {
		canonicalDesired.SourceIPRanges = rawDesired.SourceIPRanges
	}
	return canonicalDesired, nil
}

func canonicalizeForwardingRuleNewState(c *Client, rawNew, rawDesired *ForwardingRule) (*ForwardingRule, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.AllPorts) && dcl.IsEmptyValueIndirect(rawDesired.AllPorts) {
		rawNew.AllPorts = rawDesired.AllPorts
	} else {
		if dcl.BoolCanonicalize(rawDesired.AllPorts, rawNew.AllPorts) {
			rawNew.AllPorts = rawDesired.AllPorts
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.AllowGlobalAccess) && dcl.IsEmptyValueIndirect(rawDesired.AllowGlobalAccess) {
		rawNew.AllowGlobalAccess = rawDesired.AllowGlobalAccess
	} else {
		if dcl.BoolCanonicalize(rawDesired.AllowGlobalAccess, rawNew.AllowGlobalAccess) {
			rawNew.AllowGlobalAccess = rawDesired.AllowGlobalAccess
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.LabelFingerprint) && dcl.IsEmptyValueIndirect(rawDesired.LabelFingerprint) {
		rawNew.LabelFingerprint = rawDesired.LabelFingerprint
	} else {
		if dcl.StringCanonicalize(rawDesired.LabelFingerprint, rawNew.LabelFingerprint) {
			rawNew.LabelFingerprint = rawDesired.LabelFingerprint
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.BackendService) && dcl.IsEmptyValueIndirect(rawDesired.BackendService) {
		rawNew.BackendService = rawDesired.BackendService
	} else {
		if dcl.StringCanonicalize(rawDesired.BackendService, rawNew.BackendService) {
			rawNew.BackendService = rawDesired.BackendService
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreationTimestamp) && dcl.IsEmptyValueIndirect(rawDesired.CreationTimestamp) {
		rawNew.CreationTimestamp = rawDesired.CreationTimestamp
	} else {
		if dcl.StringCanonicalize(rawDesired.CreationTimestamp, rawNew.CreationTimestamp) {
			rawNew.CreationTimestamp = rawDesired.CreationTimestamp
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.IPAddress) && dcl.IsEmptyValueIndirect(rawDesired.IPAddress) {
		rawNew.IPAddress = rawDesired.IPAddress
	} else {
		if canonicalizeIPAddressToReference(rawDesired.IPAddress, rawNew.IPAddress) {
			rawNew.IPAddress = rawDesired.IPAddress
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.IPProtocol) && dcl.IsEmptyValueIndirect(rawDesired.IPProtocol) {
		rawNew.IPProtocol = rawDesired.IPProtocol
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.IPVersion) && dcl.IsEmptyValueIndirect(rawDesired.IPVersion) {
		rawNew.IPVersion = rawDesired.IPVersion
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.IsMirroringCollector) && dcl.IsEmptyValueIndirect(rawDesired.IsMirroringCollector) {
		rawNew.IsMirroringCollector = rawDesired.IsMirroringCollector
	} else {
		if dcl.BoolCanonicalize(rawDesired.IsMirroringCollector, rawNew.IsMirroringCollector) {
			rawNew.IsMirroringCollector = rawDesired.IsMirroringCollector
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.LoadBalancingScheme) && dcl.IsEmptyValueIndirect(rawDesired.LoadBalancingScheme) {
		rawNew.LoadBalancingScheme = rawDesired.LoadBalancingScheme
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.MetadataFilter) && dcl.IsEmptyValueIndirect(rawDesired.MetadataFilter) {
		rawNew.MetadataFilter = rawDesired.MetadataFilter
	} else {
		rawNew.MetadataFilter = canonicalizeNewForwardingRuleMetadataFilterSlice(c, rawDesired.MetadataFilter, rawNew.MetadataFilter)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Network) && dcl.IsEmptyValueIndirect(rawDesired.Network) {
		rawNew.Network = rawDesired.Network
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Network, rawNew.Network) {
			rawNew.Network = rawDesired.Network
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.NetworkTier) && dcl.IsEmptyValueIndirect(rawDesired.NetworkTier) {
		rawNew.NetworkTier = rawDesired.NetworkTier
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PortRange) && dcl.IsEmptyValueIndirect(rawDesired.PortRange) {
		rawNew.PortRange = rawDesired.PortRange
	} else {
		if canonicalizePortRange(rawDesired.PortRange, rawNew.PortRange) {
			rawNew.PortRange = rawDesired.PortRange
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Ports) && dcl.IsEmptyValueIndirect(rawDesired.Ports) {
		rawNew.Ports = rawDesired.Ports
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.Ports, rawNew.Ports) {
			rawNew.Ports = rawDesired.Ports
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Region) && dcl.IsEmptyValueIndirect(rawDesired.Region) {
		rawNew.Region = rawDesired.Region
	} else {
		if dcl.StringCanonicalize(rawDesired.Region, rawNew.Region) {
			rawNew.Region = rawDesired.Region
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.ServiceLabel) && dcl.IsEmptyValueIndirect(rawDesired.ServiceLabel) {
		rawNew.ServiceLabel = rawDesired.ServiceLabel
	} else {
		if dcl.StringCanonicalize(rawDesired.ServiceLabel, rawNew.ServiceLabel) {
			rawNew.ServiceLabel = rawDesired.ServiceLabel
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.ServiceName) && dcl.IsEmptyValueIndirect(rawDesired.ServiceName) {
		rawNew.ServiceName = rawDesired.ServiceName
	} else {
		if dcl.StringCanonicalize(rawDesired.ServiceName, rawNew.ServiceName) {
			rawNew.ServiceName = rawDesired.ServiceName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Subnetwork) && dcl.IsEmptyValueIndirect(rawDesired.Subnetwork) {
		rawNew.Subnetwork = rawDesired.Subnetwork
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Subnetwork, rawNew.Subnetwork) {
			rawNew.Subnetwork = rawDesired.Subnetwork
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Target) && dcl.IsEmptyValueIndirect(rawDesired.Target) {
		rawNew.Target = rawDesired.Target
	} else {
		if dcl.StringCanonicalize(rawDesired.Target, rawNew.Target) {
			rawNew.Target = rawDesired.Target
		}
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	if dcl.IsEmptyValueIndirect(rawNew.ServiceDirectoryRegistrations) && dcl.IsEmptyValueIndirect(rawDesired.ServiceDirectoryRegistrations) {
		rawNew.ServiceDirectoryRegistrations = rawDesired.ServiceDirectoryRegistrations
	} else {
		rawNew.ServiceDirectoryRegistrations = canonicalizeNewForwardingRuleServiceDirectoryRegistrationsSlice(c, rawDesired.ServiceDirectoryRegistrations, rawNew.ServiceDirectoryRegistrations)
	}

	if dcl.IsEmptyValueIndirect(rawNew.PscConnectionId) && dcl.IsEmptyValueIndirect(rawDesired.PscConnectionId) {
		rawNew.PscConnectionId = rawDesired.PscConnectionId
	} else {
		if dcl.StringCanonicalize(rawDesired.PscConnectionId, rawNew.PscConnectionId) {
			rawNew.PscConnectionId = rawDesired.PscConnectionId
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.PscConnectionStatus) && dcl.IsEmptyValueIndirect(rawDesired.PscConnectionStatus) {
		rawNew.PscConnectionStatus = rawDesired.PscConnectionStatus
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.SourceIPRanges) && dcl.IsEmptyValueIndirect(rawDesired.SourceIPRanges) {
		rawNew.SourceIPRanges = rawDesired.SourceIPRanges
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.SourceIPRanges, rawNew.SourceIPRanges) {
			rawNew.SourceIPRanges = rawDesired.SourceIPRanges
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.BaseForwardingRule) && dcl.IsEmptyValueIndirect(rawDesired.BaseForwardingRule) {
		rawNew.BaseForwardingRule = rawDesired.BaseForwardingRule
	} else {
		if dcl.StringCanonicalize(rawDesired.BaseForwardingRule, rawNew.BaseForwardingRule) {
			rawNew.BaseForwardingRule = rawDesired.BaseForwardingRule
		}
	}

	return rawNew, nil
}

func canonicalizeForwardingRuleMetadataFilter(des, initial *ForwardingRuleMetadataFilter, opts ...dcl.ApplyOption) *ForwardingRuleMetadataFilter {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ForwardingRuleMetadataFilter{}

	if dcl.IsZeroValue(des.FilterMatchCriteria) || (dcl.IsEmptyValueIndirect(des.FilterMatchCriteria) && dcl.IsEmptyValueIndirect(initial.FilterMatchCriteria)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.FilterMatchCriteria = initial.FilterMatchCriteria
	} else {
		cDes.FilterMatchCriteria = des.FilterMatchCriteria
	}
	cDes.FilterLabel = canonicalizeForwardingRuleMetadataFilterFilterLabelSlice(des.FilterLabel, initial.FilterLabel, opts...)

	return cDes
}

func canonicalizeForwardingRuleMetadataFilterSlice(des, initial []ForwardingRuleMetadataFilter, opts ...dcl.ApplyOption) []ForwardingRuleMetadataFilter {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ForwardingRuleMetadataFilter, 0, len(des))
		for _, d := range des {
			cd := canonicalizeForwardingRuleMetadataFilter(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ForwardingRuleMetadataFilter, 0, len(des))
	for i, d := range des {
		cd := canonicalizeForwardingRuleMetadataFilter(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewForwardingRuleMetadataFilter(c *Client, des, nw *ForwardingRuleMetadataFilter) *ForwardingRuleMetadataFilter {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ForwardingRuleMetadataFilter while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.FilterLabel = canonicalizeNewForwardingRuleMetadataFilterFilterLabelSlice(c, des.FilterLabel, nw.FilterLabel)

	return nw
}

func canonicalizeNewForwardingRuleMetadataFilterSet(c *Client, des, nw []ForwardingRuleMetadataFilter) []ForwardingRuleMetadataFilter {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ForwardingRuleMetadataFilter
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareForwardingRuleMetadataFilterNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewForwardingRuleMetadataFilter(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewForwardingRuleMetadataFilterSlice(c *Client, des, nw []ForwardingRuleMetadataFilter) []ForwardingRuleMetadataFilter {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ForwardingRuleMetadataFilter
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewForwardingRuleMetadataFilter(c, &d, &n))
	}

	return items
}

func canonicalizeForwardingRuleMetadataFilterFilterLabel(des, initial *ForwardingRuleMetadataFilterFilterLabel, opts ...dcl.ApplyOption) *ForwardingRuleMetadataFilterFilterLabel {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ForwardingRuleMetadataFilterFilterLabel{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Value, initial.Value) || dcl.IsZeroValue(des.Value) {
		cDes.Value = initial.Value
	} else {
		cDes.Value = des.Value
	}

	return cDes
}

func canonicalizeForwardingRuleMetadataFilterFilterLabelSlice(des, initial []ForwardingRuleMetadataFilterFilterLabel, opts ...dcl.ApplyOption) []ForwardingRuleMetadataFilterFilterLabel {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ForwardingRuleMetadataFilterFilterLabel, 0, len(des))
		for _, d := range des {
			cd := canonicalizeForwardingRuleMetadataFilterFilterLabel(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ForwardingRuleMetadataFilterFilterLabel, 0, len(des))
	for i, d := range des {
		cd := canonicalizeForwardingRuleMetadataFilterFilterLabel(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewForwardingRuleMetadataFilterFilterLabel(c *Client, des, nw *ForwardingRuleMetadataFilterFilterLabel) *ForwardingRuleMetadataFilterFilterLabel {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ForwardingRuleMetadataFilterFilterLabel while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewForwardingRuleMetadataFilterFilterLabelSet(c *Client, des, nw []ForwardingRuleMetadataFilterFilterLabel) []ForwardingRuleMetadataFilterFilterLabel {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ForwardingRuleMetadataFilterFilterLabel
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareForwardingRuleMetadataFilterFilterLabelNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewForwardingRuleMetadataFilterFilterLabel(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewForwardingRuleMetadataFilterFilterLabelSlice(c *Client, des, nw []ForwardingRuleMetadataFilterFilterLabel) []ForwardingRuleMetadataFilterFilterLabel {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ForwardingRuleMetadataFilterFilterLabel
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewForwardingRuleMetadataFilterFilterLabel(c, &d, &n))
	}

	return items
}

func canonicalizeForwardingRuleServiceDirectoryRegistrations(des, initial *ForwardingRuleServiceDirectoryRegistrations, opts ...dcl.ApplyOption) *ForwardingRuleServiceDirectoryRegistrations {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ForwardingRuleServiceDirectoryRegistrations{}

	if dcl.StringCanonicalize(des.Namespace, initial.Namespace) || dcl.IsZeroValue(des.Namespace) {
		cDes.Namespace = initial.Namespace
	} else {
		cDes.Namespace = des.Namespace
	}
	if dcl.StringCanonicalize(des.Service, initial.Service) || dcl.IsZeroValue(des.Service) {
		cDes.Service = initial.Service
	} else {
		cDes.Service = des.Service
	}

	return cDes
}

func canonicalizeForwardingRuleServiceDirectoryRegistrationsSlice(des, initial []ForwardingRuleServiceDirectoryRegistrations, opts ...dcl.ApplyOption) []ForwardingRuleServiceDirectoryRegistrations {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ForwardingRuleServiceDirectoryRegistrations, 0, len(des))
		for _, d := range des {
			cd := canonicalizeForwardingRuleServiceDirectoryRegistrations(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ForwardingRuleServiceDirectoryRegistrations, 0, len(des))
	for i, d := range des {
		cd := canonicalizeForwardingRuleServiceDirectoryRegistrations(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewForwardingRuleServiceDirectoryRegistrations(c *Client, des, nw *ForwardingRuleServiceDirectoryRegistrations) *ForwardingRuleServiceDirectoryRegistrations {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ForwardingRuleServiceDirectoryRegistrations while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Namespace, nw.Namespace) {
		nw.Namespace = des.Namespace
	}
	if dcl.StringCanonicalize(des.Service, nw.Service) {
		nw.Service = des.Service
	}

	return nw
}

func canonicalizeNewForwardingRuleServiceDirectoryRegistrationsSet(c *Client, des, nw []ForwardingRuleServiceDirectoryRegistrations) []ForwardingRuleServiceDirectoryRegistrations {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ForwardingRuleServiceDirectoryRegistrations
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareForwardingRuleServiceDirectoryRegistrationsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewForwardingRuleServiceDirectoryRegistrations(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewForwardingRuleServiceDirectoryRegistrationsSlice(c *Client, des, nw []ForwardingRuleServiceDirectoryRegistrations) []ForwardingRuleServiceDirectoryRegistrations {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ForwardingRuleServiceDirectoryRegistrations
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewForwardingRuleServiceDirectoryRegistrations(c, &d, &n))
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
func diffForwardingRule(c *Client, desired, actual *ForwardingRule, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateForwardingRuleSetLabelsOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllPorts, actual.AllPorts, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AllPorts")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowGlobalAccess, actual.AllowGlobalAccess, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateForwardingRuleUpdateOperation")}, fn.AddNest("AllowGlobalAccess")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LabelFingerprint, actual.LabelFingerprint, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LabelFingerprint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BackendService, actual.BackendService, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BackendService")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreationTimestamp, actual.CreationTimestamp, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreationTimestamp")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPAddress, actual.IPAddress, dcl.DiffInfo{ServerDefault: true, CustomDiff: canonicalizeIPAddressToReference, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IPAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPProtocol, actual.IPProtocol, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IPProtocol")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPVersion, actual.IPVersion, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IpVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IsMirroringCollector, actual.IsMirroringCollector, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsMirroringCollector")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoadBalancingScheme, actual.LoadBalancingScheme, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoadBalancingScheme")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetadataFilter, actual.MetadataFilter, dcl.DiffInfo{ObjectFunction: compareForwardingRuleMetadataFilterNewStyle, EmptyObject: EmptyForwardingRuleMetadataFilter, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetadataFilters")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Network, actual.Network, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Network")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NetworkTier, actual.NetworkTier, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkTier")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PortRange, actual.PortRange, dcl.DiffInfo{CustomDiff: canonicalizePortRange, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PortRange")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Ports, actual.Ports, dcl.DiffInfo{Type: "Set", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Ports")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Region, actual.Region, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Region")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SelfLink, actual.SelfLink, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SelfLink")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceLabel, actual.ServiceLabel, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceLabel")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceName, actual.ServiceName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Subnetwork, actual.Subnetwork, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Subnetwork")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Target, actual.Target, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateForwardingRuleSetTargetOperation")}, fn.AddNest("Target")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.ServiceDirectoryRegistrations, actual.ServiceDirectoryRegistrations, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareForwardingRuleServiceDirectoryRegistrationsNewStyle, EmptyObject: EmptyForwardingRuleServiceDirectoryRegistrations, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceDirectoryRegistrations")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PscConnectionId, actual.PscConnectionId, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PscConnectionId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PscConnectionStatus, actual.PscConnectionStatus, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PscConnectionStatus")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SourceIPRanges, actual.SourceIPRanges, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SourceIpRanges")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BaseForwardingRule, actual.BaseForwardingRule, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BaseForwardingRule")); len(ds) != 0 || err != nil {
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
func compareForwardingRuleMetadataFilterNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ForwardingRuleMetadataFilter)
	if !ok {
		desiredNotPointer, ok := d.(ForwardingRuleMetadataFilter)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ForwardingRuleMetadataFilter or *ForwardingRuleMetadataFilter", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ForwardingRuleMetadataFilter)
	if !ok {
		actualNotPointer, ok := a.(ForwardingRuleMetadataFilter)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ForwardingRuleMetadataFilter", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.FilterMatchCriteria, actual.FilterMatchCriteria, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FilterMatchCriteria")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FilterLabel, actual.FilterLabel, dcl.DiffInfo{ObjectFunction: compareForwardingRuleMetadataFilterFilterLabelNewStyle, EmptyObject: EmptyForwardingRuleMetadataFilterFilterLabel, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FilterLabels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareForwardingRuleMetadataFilterFilterLabelNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ForwardingRuleMetadataFilterFilterLabel)
	if !ok {
		desiredNotPointer, ok := d.(ForwardingRuleMetadataFilterFilterLabel)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ForwardingRuleMetadataFilterFilterLabel or *ForwardingRuleMetadataFilterFilterLabel", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ForwardingRuleMetadataFilterFilterLabel)
	if !ok {
		actualNotPointer, ok := a.(ForwardingRuleMetadataFilterFilterLabel)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ForwardingRuleMetadataFilterFilterLabel", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Value, actual.Value, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Value")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareForwardingRuleServiceDirectoryRegistrationsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ForwardingRuleServiceDirectoryRegistrations)
	if !ok {
		desiredNotPointer, ok := d.(ForwardingRuleServiceDirectoryRegistrations)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ForwardingRuleServiceDirectoryRegistrations or *ForwardingRuleServiceDirectoryRegistrations", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ForwardingRuleServiceDirectoryRegistrations)
	if !ok {
		actualNotPointer, ok := a.(ForwardingRuleServiceDirectoryRegistrations)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ForwardingRuleServiceDirectoryRegistrations", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Namespace, actual.Namespace, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Namespace")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Service, actual.Service, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Service")); len(ds) != 0 || err != nil {
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
func (r *ForwardingRule) urlNormalized() *ForwardingRule {
	normalized := dcl.Copy(*r).(ForwardingRule)
	normalized.LabelFingerprint = dcl.SelfLinkToName(r.LabelFingerprint)
	normalized.BackendService = dcl.SelfLinkToName(r.BackendService)
	normalized.CreationTimestamp = dcl.SelfLinkToName(r.CreationTimestamp)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.IPAddress = dcl.SelfLinkToName(r.IPAddress)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Network = dcl.SelfLinkToName(r.Network)
	normalized.PortRange = dcl.SelfLinkToName(r.PortRange)
	normalized.Region = dcl.SelfLinkToName(r.Region)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.ServiceLabel = dcl.SelfLinkToName(r.ServiceLabel)
	normalized.ServiceName = dcl.SelfLinkToName(r.ServiceName)
	normalized.Subnetwork = dcl.SelfLinkToName(r.Subnetwork)
	normalized.Target = dcl.SelfLinkToName(r.Target)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.PscConnectionId = dcl.SelfLinkToName(r.PscConnectionId)
	normalized.BaseForwardingRule = dcl.SelfLinkToName(r.BaseForwardingRule)
	return &normalized
}

func (r *ForwardingRule) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "setLabels" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		if dcl.IsRegion(nr.Location) {
			return dcl.URL("projects/{{project}}/regions/{{location}}/forwardingRules/{{name}}/setLabels", nr.basePath(), userBasePath, fields), nil
		}

		return dcl.URL("projects/{{project}}/global/forwardingRules/{{name}}/setLabels", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "setTarget" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		if dcl.IsRegion(nr.Location) {
			return dcl.URL("projects/{{project}}/regions/{{location}}/forwardingRules/{{name}}/setTarget", nr.basePath(), userBasePath, fields), nil
		}

		return dcl.URL("projects/{{project}}/global/forwardingRules/{{name}}/setTarget", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "update" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		if dcl.IsRegion(nr.Location) {
			return dcl.URL("projects/{{project}}/regions/{{location}}/forwardingRules/{{name}}", nr.basePath(), userBasePath, fields), nil
		}

		return dcl.URL("projects/{{project}}/global/forwardingRules/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the ForwardingRule resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *ForwardingRule) marshal(c *Client) ([]byte, error) {
	m, err := expandForwardingRule(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling ForwardingRule: %w", err)
	}
	m = forwardingRuleEncodeCreateRequest(m)

	return json.Marshal(m)
}

// unmarshalForwardingRule decodes JSON responses into the ForwardingRule resource schema.
func unmarshalForwardingRule(b []byte, c *Client, res *ForwardingRule) (*ForwardingRule, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapForwardingRule(m, c, res)
}

func unmarshalMapForwardingRule(m map[string]interface{}, c *Client, res *ForwardingRule) (*ForwardingRule, error) {

	flattened := flattenForwardingRule(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandForwardingRule expands ForwardingRule into a JSON request object.
func expandForwardingRule(c *Client, f *ForwardingRule) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v, err := dcl.FalseToNil(f.AllPorts); err != nil {
		return nil, fmt.Errorf("error expanding AllPorts into allPorts: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["allPorts"] = v
	}
	if v, err := dcl.FalseToNil(f.AllowGlobalAccess); err != nil {
		return nil, fmt.Errorf("error expanding AllowGlobalAccess into allowGlobalAccess: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["allowGlobalAccess"] = v
	}
	if v := f.BackendService; dcl.ValueShouldBeSent(v) {
		m["backendService"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.IPAddress; dcl.ValueShouldBeSent(v) {
		m["IPAddress"] = v
	}
	if v := f.IPProtocol; dcl.ValueShouldBeSent(v) {
		m["IPProtocol"] = v
	}
	if v := f.IPVersion; dcl.ValueShouldBeSent(v) {
		m["ipVersion"] = v
	}
	if v, err := dcl.FalseToNil(f.IsMirroringCollector); err != nil {
		return nil, fmt.Errorf("error expanding IsMirroringCollector into isMirroringCollector: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["isMirroringCollector"] = v
	}
	if v := f.LoadBalancingScheme; dcl.ValueShouldBeSent(v) {
		m["loadBalancingScheme"] = v
	}
	if v, err := expandForwardingRuleMetadataFilterSlice(c, f.MetadataFilter, res); err != nil {
		return nil, fmt.Errorf("error expanding MetadataFilter into metadataFilters: %w", err)
	} else if v != nil {
		m["metadataFilters"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v, err := dcl.DeriveField("global/networks/%s", f.Network, dcl.SelfLinkToName(f.Network)); err != nil {
		return nil, fmt.Errorf("error expanding Network into network: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["network"] = v
	}
	if v := f.NetworkTier; dcl.ValueShouldBeSent(v) {
		m["networkTier"] = v
	}
	if v := f.PortRange; dcl.ValueShouldBeSent(v) {
		m["portRange"] = v
	}
	if v := f.Ports; v != nil {
		m["ports"] = v
	}
	if v := f.ServiceLabel; dcl.ValueShouldBeSent(v) {
		m["serviceLabel"] = v
	}
	if v, err := dcl.DeriveField("projects/%s/regions/%s/subnetworks/%s", f.Subnetwork, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Subnetwork)); err != nil {
		return nil, fmt.Errorf("error expanding Subnetwork into subnetwork: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subnetwork"] = v
	}
	if v := f.Target; dcl.ValueShouldBeSent(v) {
		m["target"] = v
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
	if v, err := expandForwardingRuleServiceDirectoryRegistrationsSlice(c, f.ServiceDirectoryRegistrations, res); err != nil {
		return nil, fmt.Errorf("error expanding ServiceDirectoryRegistrations into serviceDirectoryRegistrations: %w", err)
	} else if v != nil {
		m["serviceDirectoryRegistrations"] = v
	}
	if v := f.SourceIPRanges; v != nil {
		m["sourceIpRanges"] = v
	}

	return m, nil
}

// flattenForwardingRule flattens ForwardingRule from a JSON request object into the
// ForwardingRule type.
func flattenForwardingRule(c *Client, i interface{}, res *ForwardingRule) *ForwardingRule {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &ForwardingRule{}
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.AllPorts = dcl.FlattenBool(m["allPorts"])
	resultRes.AllowGlobalAccess = dcl.FlattenBool(m["allowGlobalAccess"])
	resultRes.LabelFingerprint = dcl.FlattenString(m["labelFingerprint"])
	resultRes.BackendService = dcl.FlattenString(m["backendService"])
	resultRes.CreationTimestamp = dcl.FlattenString(m["creationTimestamp"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.IPAddress = dcl.FlattenString(m["IPAddress"])
	resultRes.IPProtocol = flattenForwardingRuleIPProtocolEnum(m["IPProtocol"])
	resultRes.IPVersion = flattenForwardingRuleIPVersionEnum(m["ipVersion"])
	resultRes.IsMirroringCollector = dcl.FlattenBool(m["isMirroringCollector"])
	resultRes.LoadBalancingScheme = flattenForwardingRuleLoadBalancingSchemeEnum(m["loadBalancingScheme"])
	resultRes.MetadataFilter = flattenForwardingRuleMetadataFilterSlice(c, m["metadataFilters"], res)
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Network = dcl.FlattenString(m["network"])
	resultRes.NetworkTier = flattenForwardingRuleNetworkTierEnum(m["networkTier"])
	resultRes.PortRange = dcl.FlattenString(m["portRange"])
	resultRes.Ports = dcl.FlattenStringSlice(m["ports"])
	resultRes.Region = dcl.FlattenString(m["region"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.ServiceLabel = dcl.FlattenString(m["serviceLabel"])
	resultRes.ServiceName = dcl.FlattenString(m["serviceName"])
	resultRes.Subnetwork = dcl.FlattenString(m["subnetwork"])
	resultRes.Target = dcl.FlattenString(m["target"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.ServiceDirectoryRegistrations = flattenForwardingRuleServiceDirectoryRegistrationsSlice(c, m["serviceDirectoryRegistrations"], res)
	resultRes.PscConnectionId = dcl.FlattenString(m["pscConnectionId"])
	resultRes.PscConnectionStatus = flattenForwardingRulePscConnectionStatusEnum(m["pscConnectionStatus"])
	resultRes.SourceIPRanges = dcl.FlattenStringSlice(m["sourceIpRanges"])
	resultRes.BaseForwardingRule = dcl.FlattenString(m["baseForwardingRule"])

	return resultRes
}

// expandForwardingRuleMetadataFilterMap expands the contents of ForwardingRuleMetadataFilter into a JSON
// request object.
func expandForwardingRuleMetadataFilterMap(c *Client, f map[string]ForwardingRuleMetadataFilter, res *ForwardingRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandForwardingRuleMetadataFilter(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandForwardingRuleMetadataFilterSlice expands the contents of ForwardingRuleMetadataFilter into a JSON
// request object.
func expandForwardingRuleMetadataFilterSlice(c *Client, f []ForwardingRuleMetadataFilter, res *ForwardingRule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandForwardingRuleMetadataFilter(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenForwardingRuleMetadataFilterMap flattens the contents of ForwardingRuleMetadataFilter from a JSON
// response object.
func flattenForwardingRuleMetadataFilterMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRuleMetadataFilter {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRuleMetadataFilter{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRuleMetadataFilter{}
	}

	items := make(map[string]ForwardingRuleMetadataFilter)
	for k, item := range a {
		items[k] = *flattenForwardingRuleMetadataFilter(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenForwardingRuleMetadataFilterSlice flattens the contents of ForwardingRuleMetadataFilter from a JSON
// response object.
func flattenForwardingRuleMetadataFilterSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRuleMetadataFilter {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRuleMetadataFilter{}
	}

	if len(a) == 0 {
		return []ForwardingRuleMetadataFilter{}
	}

	items := make([]ForwardingRuleMetadataFilter, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRuleMetadataFilter(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandForwardingRuleMetadataFilter expands an instance of ForwardingRuleMetadataFilter into a JSON
// request object.
func expandForwardingRuleMetadataFilter(c *Client, f *ForwardingRuleMetadataFilter, res *ForwardingRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.FilterMatchCriteria; !dcl.IsEmptyValueIndirect(v) {
		m["filterMatchCriteria"] = v
	}
	if v, err := expandForwardingRuleMetadataFilterFilterLabelSlice(c, f.FilterLabel, res); err != nil {
		return nil, fmt.Errorf("error expanding FilterLabel into filterLabels: %w", err)
	} else if v != nil {
		m["filterLabels"] = v
	}

	return m, nil
}

// flattenForwardingRuleMetadataFilter flattens an instance of ForwardingRuleMetadataFilter from a JSON
// response object.
func flattenForwardingRuleMetadataFilter(c *Client, i interface{}, res *ForwardingRule) *ForwardingRuleMetadataFilter {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ForwardingRuleMetadataFilter{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyForwardingRuleMetadataFilter
	}
	r.FilterMatchCriteria = flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnum(m["filterMatchCriteria"])
	r.FilterLabel = flattenForwardingRuleMetadataFilterFilterLabelSlice(c, m["filterLabels"], res)

	return r
}

// expandForwardingRuleMetadataFilterFilterLabelMap expands the contents of ForwardingRuleMetadataFilterFilterLabel into a JSON
// request object.
func expandForwardingRuleMetadataFilterFilterLabelMap(c *Client, f map[string]ForwardingRuleMetadataFilterFilterLabel, res *ForwardingRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandForwardingRuleMetadataFilterFilterLabel(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandForwardingRuleMetadataFilterFilterLabelSlice expands the contents of ForwardingRuleMetadataFilterFilterLabel into a JSON
// request object.
func expandForwardingRuleMetadataFilterFilterLabelSlice(c *Client, f []ForwardingRuleMetadataFilterFilterLabel, res *ForwardingRule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandForwardingRuleMetadataFilterFilterLabel(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenForwardingRuleMetadataFilterFilterLabelMap flattens the contents of ForwardingRuleMetadataFilterFilterLabel from a JSON
// response object.
func flattenForwardingRuleMetadataFilterFilterLabelMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRuleMetadataFilterFilterLabel {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRuleMetadataFilterFilterLabel{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRuleMetadataFilterFilterLabel{}
	}

	items := make(map[string]ForwardingRuleMetadataFilterFilterLabel)
	for k, item := range a {
		items[k] = *flattenForwardingRuleMetadataFilterFilterLabel(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenForwardingRuleMetadataFilterFilterLabelSlice flattens the contents of ForwardingRuleMetadataFilterFilterLabel from a JSON
// response object.
func flattenForwardingRuleMetadataFilterFilterLabelSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRuleMetadataFilterFilterLabel {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRuleMetadataFilterFilterLabel{}
	}

	if len(a) == 0 {
		return []ForwardingRuleMetadataFilterFilterLabel{}
	}

	items := make([]ForwardingRuleMetadataFilterFilterLabel, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRuleMetadataFilterFilterLabel(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandForwardingRuleMetadataFilterFilterLabel expands an instance of ForwardingRuleMetadataFilterFilterLabel into a JSON
// request object.
func expandForwardingRuleMetadataFilterFilterLabel(c *Client, f *ForwardingRuleMetadataFilterFilterLabel, res *ForwardingRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Value; !dcl.IsEmptyValueIndirect(v) {
		m["value"] = v
	}

	return m, nil
}

// flattenForwardingRuleMetadataFilterFilterLabel flattens an instance of ForwardingRuleMetadataFilterFilterLabel from a JSON
// response object.
func flattenForwardingRuleMetadataFilterFilterLabel(c *Client, i interface{}, res *ForwardingRule) *ForwardingRuleMetadataFilterFilterLabel {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ForwardingRuleMetadataFilterFilterLabel{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyForwardingRuleMetadataFilterFilterLabel
	}
	r.Name = dcl.FlattenString(m["name"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandForwardingRuleServiceDirectoryRegistrationsMap expands the contents of ForwardingRuleServiceDirectoryRegistrations into a JSON
// request object.
func expandForwardingRuleServiceDirectoryRegistrationsMap(c *Client, f map[string]ForwardingRuleServiceDirectoryRegistrations, res *ForwardingRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandForwardingRuleServiceDirectoryRegistrations(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandForwardingRuleServiceDirectoryRegistrationsSlice expands the contents of ForwardingRuleServiceDirectoryRegistrations into a JSON
// request object.
func expandForwardingRuleServiceDirectoryRegistrationsSlice(c *Client, f []ForwardingRuleServiceDirectoryRegistrations, res *ForwardingRule) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandForwardingRuleServiceDirectoryRegistrations(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenForwardingRuleServiceDirectoryRegistrationsMap flattens the contents of ForwardingRuleServiceDirectoryRegistrations from a JSON
// response object.
func flattenForwardingRuleServiceDirectoryRegistrationsMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRuleServiceDirectoryRegistrations {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRuleServiceDirectoryRegistrations{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRuleServiceDirectoryRegistrations{}
	}

	items := make(map[string]ForwardingRuleServiceDirectoryRegistrations)
	for k, item := range a {
		items[k] = *flattenForwardingRuleServiceDirectoryRegistrations(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenForwardingRuleServiceDirectoryRegistrationsSlice flattens the contents of ForwardingRuleServiceDirectoryRegistrations from a JSON
// response object.
func flattenForwardingRuleServiceDirectoryRegistrationsSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRuleServiceDirectoryRegistrations {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRuleServiceDirectoryRegistrations{}
	}

	if len(a) == 0 {
		return []ForwardingRuleServiceDirectoryRegistrations{}
	}

	items := make([]ForwardingRuleServiceDirectoryRegistrations, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRuleServiceDirectoryRegistrations(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandForwardingRuleServiceDirectoryRegistrations expands an instance of ForwardingRuleServiceDirectoryRegistrations into a JSON
// request object.
func expandForwardingRuleServiceDirectoryRegistrations(c *Client, f *ForwardingRuleServiceDirectoryRegistrations, res *ForwardingRule) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Namespace; !dcl.IsEmptyValueIndirect(v) {
		m["namespace"] = v
	}
	if v := f.Service; !dcl.IsEmptyValueIndirect(v) {
		m["service"] = v
	}

	return m, nil
}

// flattenForwardingRuleServiceDirectoryRegistrations flattens an instance of ForwardingRuleServiceDirectoryRegistrations from a JSON
// response object.
func flattenForwardingRuleServiceDirectoryRegistrations(c *Client, i interface{}, res *ForwardingRule) *ForwardingRuleServiceDirectoryRegistrations {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ForwardingRuleServiceDirectoryRegistrations{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyForwardingRuleServiceDirectoryRegistrations
	}
	r.Namespace = dcl.FlattenString(m["namespace"])
	r.Service = dcl.FlattenString(m["service"])

	return r
}

// flattenForwardingRuleIPProtocolEnumMap flattens the contents of ForwardingRuleIPProtocolEnum from a JSON
// response object.
func flattenForwardingRuleIPProtocolEnumMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRuleIPProtocolEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRuleIPProtocolEnum{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRuleIPProtocolEnum{}
	}

	items := make(map[string]ForwardingRuleIPProtocolEnum)
	for k, item := range a {
		items[k] = *flattenForwardingRuleIPProtocolEnum(item.(interface{}))
	}

	return items
}

// flattenForwardingRuleIPProtocolEnumSlice flattens the contents of ForwardingRuleIPProtocolEnum from a JSON
// response object.
func flattenForwardingRuleIPProtocolEnumSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRuleIPProtocolEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRuleIPProtocolEnum{}
	}

	if len(a) == 0 {
		return []ForwardingRuleIPProtocolEnum{}
	}

	items := make([]ForwardingRuleIPProtocolEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRuleIPProtocolEnum(item.(interface{})))
	}

	return items
}

// flattenForwardingRuleIPProtocolEnum asserts that an interface is a string, and returns a
// pointer to a *ForwardingRuleIPProtocolEnum with the same value as that string.
func flattenForwardingRuleIPProtocolEnum(i interface{}) *ForwardingRuleIPProtocolEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ForwardingRuleIPProtocolEnumRef(s)
}

// flattenForwardingRuleIPVersionEnumMap flattens the contents of ForwardingRuleIPVersionEnum from a JSON
// response object.
func flattenForwardingRuleIPVersionEnumMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRuleIPVersionEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRuleIPVersionEnum{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRuleIPVersionEnum{}
	}

	items := make(map[string]ForwardingRuleIPVersionEnum)
	for k, item := range a {
		items[k] = *flattenForwardingRuleIPVersionEnum(item.(interface{}))
	}

	return items
}

// flattenForwardingRuleIPVersionEnumSlice flattens the contents of ForwardingRuleIPVersionEnum from a JSON
// response object.
func flattenForwardingRuleIPVersionEnumSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRuleIPVersionEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRuleIPVersionEnum{}
	}

	if len(a) == 0 {
		return []ForwardingRuleIPVersionEnum{}
	}

	items := make([]ForwardingRuleIPVersionEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRuleIPVersionEnum(item.(interface{})))
	}

	return items
}

// flattenForwardingRuleIPVersionEnum asserts that an interface is a string, and returns a
// pointer to a *ForwardingRuleIPVersionEnum with the same value as that string.
func flattenForwardingRuleIPVersionEnum(i interface{}) *ForwardingRuleIPVersionEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ForwardingRuleIPVersionEnumRef(s)
}

// flattenForwardingRuleLoadBalancingSchemeEnumMap flattens the contents of ForwardingRuleLoadBalancingSchemeEnum from a JSON
// response object.
func flattenForwardingRuleLoadBalancingSchemeEnumMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRuleLoadBalancingSchemeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRuleLoadBalancingSchemeEnum{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRuleLoadBalancingSchemeEnum{}
	}

	items := make(map[string]ForwardingRuleLoadBalancingSchemeEnum)
	for k, item := range a {
		items[k] = *flattenForwardingRuleLoadBalancingSchemeEnum(item.(interface{}))
	}

	return items
}

// flattenForwardingRuleLoadBalancingSchemeEnumSlice flattens the contents of ForwardingRuleLoadBalancingSchemeEnum from a JSON
// response object.
func flattenForwardingRuleLoadBalancingSchemeEnumSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRuleLoadBalancingSchemeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRuleLoadBalancingSchemeEnum{}
	}

	if len(a) == 0 {
		return []ForwardingRuleLoadBalancingSchemeEnum{}
	}

	items := make([]ForwardingRuleLoadBalancingSchemeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRuleLoadBalancingSchemeEnum(item.(interface{})))
	}

	return items
}

// flattenForwardingRuleLoadBalancingSchemeEnum asserts that an interface is a string, and returns a
// pointer to a *ForwardingRuleLoadBalancingSchemeEnum with the same value as that string.
func flattenForwardingRuleLoadBalancingSchemeEnum(i interface{}) *ForwardingRuleLoadBalancingSchemeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ForwardingRuleLoadBalancingSchemeEnumRef(s)
}

// flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnumMap flattens the contents of ForwardingRuleMetadataFilterFilterMatchCriteriaEnum from a JSON
// response object.
func flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnumMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRuleMetadataFilterFilterMatchCriteriaEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRuleMetadataFilterFilterMatchCriteriaEnum{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRuleMetadataFilterFilterMatchCriteriaEnum{}
	}

	items := make(map[string]ForwardingRuleMetadataFilterFilterMatchCriteriaEnum)
	for k, item := range a {
		items[k] = *flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnum(item.(interface{}))
	}

	return items
}

// flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnumSlice flattens the contents of ForwardingRuleMetadataFilterFilterMatchCriteriaEnum from a JSON
// response object.
func flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnumSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRuleMetadataFilterFilterMatchCriteriaEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRuleMetadataFilterFilterMatchCriteriaEnum{}
	}

	if len(a) == 0 {
		return []ForwardingRuleMetadataFilterFilterMatchCriteriaEnum{}
	}

	items := make([]ForwardingRuleMetadataFilterFilterMatchCriteriaEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnum(item.(interface{})))
	}

	return items
}

// flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnum asserts that an interface is a string, and returns a
// pointer to a *ForwardingRuleMetadataFilterFilterMatchCriteriaEnum with the same value as that string.
func flattenForwardingRuleMetadataFilterFilterMatchCriteriaEnum(i interface{}) *ForwardingRuleMetadataFilterFilterMatchCriteriaEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ForwardingRuleMetadataFilterFilterMatchCriteriaEnumRef(s)
}

// flattenForwardingRuleNetworkTierEnumMap flattens the contents of ForwardingRuleNetworkTierEnum from a JSON
// response object.
func flattenForwardingRuleNetworkTierEnumMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRuleNetworkTierEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRuleNetworkTierEnum{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRuleNetworkTierEnum{}
	}

	items := make(map[string]ForwardingRuleNetworkTierEnum)
	for k, item := range a {
		items[k] = *flattenForwardingRuleNetworkTierEnum(item.(interface{}))
	}

	return items
}

// flattenForwardingRuleNetworkTierEnumSlice flattens the contents of ForwardingRuleNetworkTierEnum from a JSON
// response object.
func flattenForwardingRuleNetworkTierEnumSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRuleNetworkTierEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRuleNetworkTierEnum{}
	}

	if len(a) == 0 {
		return []ForwardingRuleNetworkTierEnum{}
	}

	items := make([]ForwardingRuleNetworkTierEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRuleNetworkTierEnum(item.(interface{})))
	}

	return items
}

// flattenForwardingRuleNetworkTierEnum asserts that an interface is a string, and returns a
// pointer to a *ForwardingRuleNetworkTierEnum with the same value as that string.
func flattenForwardingRuleNetworkTierEnum(i interface{}) *ForwardingRuleNetworkTierEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ForwardingRuleNetworkTierEnumRef(s)
}

// flattenForwardingRulePscConnectionStatusEnumMap flattens the contents of ForwardingRulePscConnectionStatusEnum from a JSON
// response object.
func flattenForwardingRulePscConnectionStatusEnumMap(c *Client, i interface{}, res *ForwardingRule) map[string]ForwardingRulePscConnectionStatusEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ForwardingRulePscConnectionStatusEnum{}
	}

	if len(a) == 0 {
		return map[string]ForwardingRulePscConnectionStatusEnum{}
	}

	items := make(map[string]ForwardingRulePscConnectionStatusEnum)
	for k, item := range a {
		items[k] = *flattenForwardingRulePscConnectionStatusEnum(item.(interface{}))
	}

	return items
}

// flattenForwardingRulePscConnectionStatusEnumSlice flattens the contents of ForwardingRulePscConnectionStatusEnum from a JSON
// response object.
func flattenForwardingRulePscConnectionStatusEnumSlice(c *Client, i interface{}, res *ForwardingRule) []ForwardingRulePscConnectionStatusEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ForwardingRulePscConnectionStatusEnum{}
	}

	if len(a) == 0 {
		return []ForwardingRulePscConnectionStatusEnum{}
	}

	items := make([]ForwardingRulePscConnectionStatusEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenForwardingRulePscConnectionStatusEnum(item.(interface{})))
	}

	return items
}

// flattenForwardingRulePscConnectionStatusEnum asserts that an interface is a string, and returns a
// pointer to a *ForwardingRulePscConnectionStatusEnum with the same value as that string.
func flattenForwardingRulePscConnectionStatusEnum(i interface{}) *ForwardingRulePscConnectionStatusEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ForwardingRulePscConnectionStatusEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *ForwardingRule) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalForwardingRule(b, c, r)
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

type forwardingRuleDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         forwardingRuleApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToForwardingRuleDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]forwardingRuleDiff, error) {
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
	var diffs []forwardingRuleDiff
	// For each operation name, create a forwardingRuleDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := forwardingRuleDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToForwardingRuleApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToForwardingRuleApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (forwardingRuleApiOperation, error) {
	switch opName {

	case "updateForwardingRuleSetLabelsOperation":
		return &updateForwardingRuleSetLabelsOperation{FieldDiffs: fieldDiffs}, nil

	case "updateForwardingRuleSetTargetOperation":
		return &updateForwardingRuleSetTargetOperation{FieldDiffs: fieldDiffs}, nil

	case "updateForwardingRuleUpdateOperation":
		return &updateForwardingRuleUpdateOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractForwardingRuleFields(r *ForwardingRule) error {
	return nil
}
func extractForwardingRuleMetadataFilterFields(r *ForwardingRule, o *ForwardingRuleMetadataFilter) error {
	return nil
}
func extractForwardingRuleMetadataFilterFilterLabelFields(r *ForwardingRule, o *ForwardingRuleMetadataFilterFilterLabel) error {
	return nil
}
func extractForwardingRuleServiceDirectoryRegistrationsFields(r *ForwardingRule, o *ForwardingRuleServiceDirectoryRegistrations) error {
	return nil
}

func postReadExtractForwardingRuleFields(r *ForwardingRule) error {
	return nil
}
func postReadExtractForwardingRuleMetadataFilterFields(r *ForwardingRule, o *ForwardingRuleMetadataFilter) error {
	return nil
}
func postReadExtractForwardingRuleMetadataFilterFilterLabelFields(r *ForwardingRule, o *ForwardingRuleMetadataFilterFilterLabel) error {
	return nil
}
func postReadExtractForwardingRuleServiceDirectoryRegistrationsFields(r *ForwardingRule, o *ForwardingRuleServiceDirectoryRegistrations) error {
	return nil
}
