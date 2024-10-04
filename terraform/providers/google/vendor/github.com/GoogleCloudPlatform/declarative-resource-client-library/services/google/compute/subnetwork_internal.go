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

func (r *Subnetwork) validate() error {

	if err := dcl.Required(r, "ipCidrRange"); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "network"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Region, "Region"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.LogConfig) {
		if err := r.LogConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *SubnetworkSecondaryIPRanges) validate() error {
	if err := dcl.Required(r, "rangeName"); err != nil {
		return err
	}
	if err := dcl.Required(r, "ipCidrRange"); err != nil {
		return err
	}
	return nil
}
func (r *SubnetworkLogConfig) validate() error {
	return nil
}
func (r *Subnetwork) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *Subnetwork) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"region":  dcl.ValueOrEmptyString(nr.Region),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{region}}/subnetworks/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Subnetwork) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"region":  dcl.ValueOrEmptyString(nr.Region),
	}
	return dcl.URL("projects/{{project}}/regions/{{region}}/subnetworks", nr.basePath(), userBasePath, params), nil

}

func (r *Subnetwork) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"region":  dcl.ValueOrEmptyString(nr.Region),
	}
	return dcl.URL("projects/{{project}}/regions/{{region}}/subnetworks", nr.basePath(), userBasePath, params), nil

}

func (r *Subnetwork) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"region":  dcl.ValueOrEmptyString(nr.Region),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{region}}/subnetworks/{{name}}", nr.basePath(), userBasePath, params), nil
}

// subnetworkApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type subnetworkApiOperation interface {
	do(context.Context, *Subnetwork, *Client) error
}

// newUpdateSubnetworkExpandIpCidrRangeRequest creates a request for an
// Subnetwork resource's expandIpCidrRange update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateSubnetworkExpandIpCidrRangeRequest(ctx context.Context, f *Subnetwork, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.IPCidrRange; !dcl.IsEmptyValueIndirect(v) {
		req["ipCidrRange"] = v
	}
	return req, nil
}

// marshalUpdateSubnetworkExpandIpCidrRangeRequest converts the update into
// the final JSON request body.
func marshalUpdateSubnetworkExpandIpCidrRangeRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateSubnetworkExpandIpCidrRangeOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateSubnetworkExpandIpCidrRangeOperation) do(ctx context.Context, r *Subnetwork, c *Client) error {
	_, err := c.GetSubnetwork(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "expandIpCidrRange")
	if err != nil {
		return err
	}

	req, err := newUpdateSubnetworkExpandIpCidrRangeRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateSubnetworkExpandIpCidrRangeRequest(c, req)
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

// newUpdateSubnetworkSetPrivateIpGoogleAccessRequest creates a request for an
// Subnetwork resource's setPrivateIpGoogleAccess update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateSubnetworkSetPrivateIpGoogleAccessRequest(ctx context.Context, f *Subnetwork, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.PrivateIPGoogleAccess; !dcl.IsEmptyValueIndirect(v) {
		req["privateIpGoogleAccess"] = v
	}
	return req, nil
}

// marshalUpdateSubnetworkSetPrivateIpGoogleAccessRequest converts the update into
// the final JSON request body.
func marshalUpdateSubnetworkSetPrivateIpGoogleAccessRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateSubnetworkSetPrivateIpGoogleAccessOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateSubnetworkSetPrivateIpGoogleAccessOperation) do(ctx context.Context, r *Subnetwork, c *Client) error {
	_, err := c.GetSubnetwork(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "setPrivateIpGoogleAccess")
	if err != nil {
		return err
	}

	req, err := newUpdateSubnetworkSetPrivateIpGoogleAccessRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateSubnetworkSetPrivateIpGoogleAccessRequest(c, req)
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

// newUpdateSubnetworkUpdateRequest creates a request for an
// Subnetwork resource's update update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateSubnetworkUpdateRequest(ctx context.Context, f *Subnetwork, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Role; !dcl.IsEmptyValueIndirect(v) {
		req["role"] = v
	}
	if v, err := expandSubnetworkSecondaryIPRangesSlice(c, f.SecondaryIPRanges, res); err != nil {
		return nil, fmt.Errorf("error expanding SecondaryIPRanges into secondaryIpRanges: %w", err)
	} else if v != nil {
		req["secondaryIpRanges"] = v
	}
	if v, err := expandSubnetworkLogConfig(c, f.LogConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LogConfig into logConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["logConfig"] = v
	}
	b, err := c.getSubnetworkRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawFingerprint, err := dcl.GetMapEntry(
		m,
		[]string{"fingerprint"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["fingerprint"] = rawFingerprint.(string)
	}
	return req, nil
}

// marshalUpdateSubnetworkUpdateRequest converts the update into
// the final JSON request body.
func marshalUpdateSubnetworkUpdateRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateSubnetworkUpdateOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (c *Client) listSubnetworkRaw(ctx context.Context, r *Subnetwork, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != SubnetworkMaxPage {
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

type listSubnetworkOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listSubnetwork(ctx context.Context, r *Subnetwork, pageToken string, pageSize int32) ([]*Subnetwork, string, error) {
	b, err := c.listSubnetworkRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listSubnetworkOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Subnetwork
	for _, v := range m.Items {
		res, err := unmarshalMapSubnetwork(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Region = r.Region
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllSubnetwork(ctx context.Context, f func(*Subnetwork) bool, resources []*Subnetwork) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteSubnetwork(ctx, res)
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

type deleteSubnetworkOperation struct{}

func (op *deleteSubnetworkOperation) do(ctx context.Context, r *Subnetwork, c *Client) error {
	r, err := c.GetSubnetwork(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Subnetwork not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetSubnetwork checking for existence. error: %v", err)
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
		_, err := c.GetSubnetwork(ctx, r)
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
type createSubnetworkOperation struct {
	response map[string]interface{}
}

func (op *createSubnetworkOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createSubnetworkOperation) do(ctx context.Context, r *Subnetwork, c *Client) error {
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

	if _, err := c.GetSubnetwork(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getSubnetworkRaw(ctx context.Context, r *Subnetwork) ([]byte, error) {

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

func (c *Client) subnetworkDiffsForRawDesired(ctx context.Context, rawDesired *Subnetwork, opts ...dcl.ApplyOption) (initial, desired *Subnetwork, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Subnetwork
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Subnetwork); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Subnetwork, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetSubnetwork(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Subnetwork resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Subnetwork resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Subnetwork resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeSubnetworkDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Subnetwork: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Subnetwork: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractSubnetworkFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeSubnetworkInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Subnetwork: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeSubnetworkDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Subnetwork: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffSubnetwork(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeSubnetworkInitialState(rawInitial, rawDesired *Subnetwork) (*Subnetwork, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeSubnetworkDesiredState(rawDesired, rawInitial *Subnetwork, opts ...dcl.ApplyOption) (*Subnetwork, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.LogConfig = canonicalizeSubnetworkLogConfig(rawDesired.LogConfig, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Subnetwork{}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.StringCanonicalize(rawDesired.IPCidrRange, rawInitial.IPCidrRange) {
		canonicalDesired.IPCidrRange = rawInitial.IPCidrRange
	} else {
		canonicalDesired.IPCidrRange = rawDesired.IPCidrRange
	}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.IsZeroValue(rawDesired.Network) || (dcl.IsEmptyValueIndirect(rawDesired.Network) && dcl.IsEmptyValueIndirect(rawInitial.Network)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Network = rawInitial.Network
	} else {
		canonicalDesired.Network = rawDesired.Network
	}
	if dcl.IsZeroValue(rawDesired.Purpose) || (dcl.IsEmptyValueIndirect(rawDesired.Purpose) && dcl.IsEmptyValueIndirect(rawInitial.Purpose)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Purpose = rawInitial.Purpose
	} else {
		canonicalDesired.Purpose = rawDesired.Purpose
	}
	if dcl.IsZeroValue(rawDesired.Role) || (dcl.IsEmptyValueIndirect(rawDesired.Role) && dcl.IsEmptyValueIndirect(rawInitial.Role)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Role = rawInitial.Role
	} else {
		canonicalDesired.Role = rawDesired.Role
	}
	canonicalDesired.SecondaryIPRanges = canonicalizeSubnetworkSecondaryIPRangesSlice(rawDesired.SecondaryIPRanges, rawInitial.SecondaryIPRanges, opts...)
	if dcl.BoolCanonicalize(rawDesired.PrivateIPGoogleAccess, rawInitial.PrivateIPGoogleAccess) {
		canonicalDesired.PrivateIPGoogleAccess = rawInitial.PrivateIPGoogleAccess
	} else {
		canonicalDesired.PrivateIPGoogleAccess = rawDesired.PrivateIPGoogleAccess
	}
	if dcl.NameToSelfLink(rawDesired.Region, rawInitial.Region) {
		canonicalDesired.Region = rawInitial.Region
	} else {
		canonicalDesired.Region = rawDesired.Region
	}
	canonicalDesired.LogConfig = canonicalizeSubnetworkLogConfig(rawDesired.LogConfig, rawInitial.LogConfig, opts...)
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	if dcl.BoolCanonicalize(rawDesired.EnableFlowLogs, rawInitial.EnableFlowLogs) {
		canonicalDesired.EnableFlowLogs = rawInitial.EnableFlowLogs
	} else {
		canonicalDesired.EnableFlowLogs = rawDesired.EnableFlowLogs
	}
	return canonicalDesired, nil
}

func canonicalizeSubnetworkNewState(c *Client, rawNew, rawDesired *Subnetwork) (*Subnetwork, error) {

	if dcl.IsEmptyValueIndirect(rawNew.CreationTimestamp) && dcl.IsEmptyValueIndirect(rawDesired.CreationTimestamp) {
		rawNew.CreationTimestamp = rawDesired.CreationTimestamp
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.GatewayAddress) && dcl.IsEmptyValueIndirect(rawDesired.GatewayAddress) {
		rawNew.GatewayAddress = rawDesired.GatewayAddress
	} else {
		if dcl.StringCanonicalize(rawDesired.GatewayAddress, rawNew.GatewayAddress) {
			rawNew.GatewayAddress = rawDesired.GatewayAddress
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.IPCidrRange) && dcl.IsEmptyValueIndirect(rawDesired.IPCidrRange) {
		rawNew.IPCidrRange = rawDesired.IPCidrRange
	} else {
		if dcl.StringCanonicalize(rawDesired.IPCidrRange, rawNew.IPCidrRange) {
			rawNew.IPCidrRange = rawDesired.IPCidrRange
		}
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
	}

	if dcl.IsEmptyValueIndirect(rawNew.Fingerprint) && dcl.IsEmptyValueIndirect(rawDesired.Fingerprint) {
		rawNew.Fingerprint = rawDesired.Fingerprint
	} else {
		if dcl.StringCanonicalize(rawDesired.Fingerprint, rawNew.Fingerprint) {
			rawNew.Fingerprint = rawDesired.Fingerprint
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Purpose) && dcl.IsEmptyValueIndirect(rawDesired.Purpose) {
		rawNew.Purpose = rawDesired.Purpose
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Role) && dcl.IsEmptyValueIndirect(rawDesired.Role) {
		rawNew.Role = rawDesired.Role
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.SecondaryIPRanges) && dcl.IsEmptyValueIndirect(rawDesired.SecondaryIPRanges) {
		rawNew.SecondaryIPRanges = rawDesired.SecondaryIPRanges
	} else {
		rawNew.SecondaryIPRanges = canonicalizeNewSubnetworkSecondaryIPRangesSlice(c, rawDesired.SecondaryIPRanges, rawNew.SecondaryIPRanges)
	}

	if dcl.IsEmptyValueIndirect(rawNew.PrivateIPGoogleAccess) && dcl.IsEmptyValueIndirect(rawDesired.PrivateIPGoogleAccess) {
		rawNew.PrivateIPGoogleAccess = rawDesired.PrivateIPGoogleAccess
	} else {
		if dcl.BoolCanonicalize(rawDesired.PrivateIPGoogleAccess, rawNew.PrivateIPGoogleAccess) {
			rawNew.PrivateIPGoogleAccess = rawDesired.PrivateIPGoogleAccess
		}
	}

	rawNew.Region = rawDesired.Region

	if dcl.IsEmptyValueIndirect(rawNew.LogConfig) && dcl.IsEmptyValueIndirect(rawDesired.LogConfig) {
		rawNew.LogConfig = rawDesired.LogConfig
	} else {
		rawNew.LogConfig = canonicalizeNewSubnetworkLogConfig(c, rawDesired.LogConfig, rawNew.LogConfig)
	}

	rawNew.Project = rawDesired.Project

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.EnableFlowLogs) && dcl.IsEmptyValueIndirect(rawDesired.EnableFlowLogs) {
		rawNew.EnableFlowLogs = rawDesired.EnableFlowLogs
	} else {
		if dcl.BoolCanonicalize(rawDesired.EnableFlowLogs, rawNew.EnableFlowLogs) {
			rawNew.EnableFlowLogs = rawDesired.EnableFlowLogs
		}
	}

	return rawNew, nil
}

func canonicalizeSubnetworkSecondaryIPRanges(des, initial *SubnetworkSecondaryIPRanges, opts ...dcl.ApplyOption) *SubnetworkSecondaryIPRanges {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &SubnetworkSecondaryIPRanges{}

	if dcl.StringCanonicalize(des.RangeName, initial.RangeName) || dcl.IsZeroValue(des.RangeName) {
		cDes.RangeName = initial.RangeName
	} else {
		cDes.RangeName = des.RangeName
	}
	if dcl.StringCanonicalize(des.IPCidrRange, initial.IPCidrRange) || dcl.IsZeroValue(des.IPCidrRange) {
		cDes.IPCidrRange = initial.IPCidrRange
	} else {
		cDes.IPCidrRange = des.IPCidrRange
	}

	return cDes
}

func canonicalizeSubnetworkSecondaryIPRangesSlice(des, initial []SubnetworkSecondaryIPRanges, opts ...dcl.ApplyOption) []SubnetworkSecondaryIPRanges {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]SubnetworkSecondaryIPRanges, 0, len(des))
		for _, d := range des {
			cd := canonicalizeSubnetworkSecondaryIPRanges(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]SubnetworkSecondaryIPRanges, 0, len(des))
	for i, d := range des {
		cd := canonicalizeSubnetworkSecondaryIPRanges(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewSubnetworkSecondaryIPRanges(c *Client, des, nw *SubnetworkSecondaryIPRanges) *SubnetworkSecondaryIPRanges {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for SubnetworkSecondaryIPRanges while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.RangeName, nw.RangeName) {
		nw.RangeName = des.RangeName
	}
	if dcl.StringCanonicalize(des.IPCidrRange, nw.IPCidrRange) {
		nw.IPCidrRange = des.IPCidrRange
	}

	return nw
}

func canonicalizeNewSubnetworkSecondaryIPRangesSet(c *Client, des, nw []SubnetworkSecondaryIPRanges) []SubnetworkSecondaryIPRanges {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []SubnetworkSecondaryIPRanges
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareSubnetworkSecondaryIPRangesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewSubnetworkSecondaryIPRanges(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewSubnetworkSecondaryIPRangesSlice(c *Client, des, nw []SubnetworkSecondaryIPRanges) []SubnetworkSecondaryIPRanges {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []SubnetworkSecondaryIPRanges
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewSubnetworkSecondaryIPRanges(c, &d, &n))
	}

	return items
}

func canonicalizeSubnetworkLogConfig(des, initial *SubnetworkLogConfig, opts ...dcl.ApplyOption) *SubnetworkLogConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if dcl.IsZeroValue(des.AggregationInterval) {
		des.AggregationInterval = SubnetworkLogConfigAggregationIntervalEnumRef("INTERVAL_5_SEC")
	}

	if dcl.IsZeroValue(des.FlowSampling) {
		des.FlowSampling = dcl.Float64(0.5)
	}

	if dcl.IsZeroValue(des.Metadata) {
		des.Metadata = SubnetworkLogConfigMetadataEnumRef("INCLUDE_ALL_METADATA")
	}

	if initial == nil {
		return des
	}

	cDes := &SubnetworkLogConfig{}

	if dcl.IsZeroValue(des.AggregationInterval) || (dcl.IsEmptyValueIndirect(des.AggregationInterval) && dcl.IsEmptyValueIndirect(initial.AggregationInterval)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AggregationInterval = initial.AggregationInterval
	} else {
		cDes.AggregationInterval = des.AggregationInterval
	}
	if dcl.IsZeroValue(des.FlowSampling) || (dcl.IsEmptyValueIndirect(des.FlowSampling) && dcl.IsEmptyValueIndirect(initial.FlowSampling)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.FlowSampling = initial.FlowSampling
	} else {
		cDes.FlowSampling = des.FlowSampling
	}
	if dcl.IsZeroValue(des.Metadata) || (dcl.IsEmptyValueIndirect(des.Metadata) && dcl.IsEmptyValueIndirect(initial.Metadata)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Metadata = initial.Metadata
	} else {
		cDes.Metadata = des.Metadata
	}

	return cDes
}

func canonicalizeSubnetworkLogConfigSlice(des, initial []SubnetworkLogConfig, opts ...dcl.ApplyOption) []SubnetworkLogConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]SubnetworkLogConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeSubnetworkLogConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]SubnetworkLogConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeSubnetworkLogConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewSubnetworkLogConfig(c *Client, des, nw *SubnetworkLogConfig) *SubnetworkLogConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for SubnetworkLogConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.IsZeroValue(nw.AggregationInterval) {
		nw.AggregationInterval = SubnetworkLogConfigAggregationIntervalEnumRef("INTERVAL_5_SEC")
	}

	if dcl.IsZeroValue(nw.FlowSampling) {
		nw.FlowSampling = dcl.Float64(0.5)
	}

	if dcl.IsZeroValue(nw.Metadata) {
		nw.Metadata = SubnetworkLogConfigMetadataEnumRef("INCLUDE_ALL_METADATA")
	}

	return nw
}

func canonicalizeNewSubnetworkLogConfigSet(c *Client, des, nw []SubnetworkLogConfig) []SubnetworkLogConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []SubnetworkLogConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareSubnetworkLogConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewSubnetworkLogConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewSubnetworkLogConfigSlice(c *Client, des, nw []SubnetworkLogConfig) []SubnetworkLogConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []SubnetworkLogConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewSubnetworkLogConfig(c, &d, &n))
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
func diffSubnetwork(c *Client, desired, actual *Subnetwork, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
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

	if ds, err := dcl.Diff(desired.GatewayAddress, actual.GatewayAddress, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GatewayAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPCidrRange, actual.IPCidrRange, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateSubnetworkExpandIpCidrRangeOperation")}, fn.AddNest("IpCidrRange")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Network, actual.Network, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Network")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Fingerprint, actual.Fingerprint, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Fingerprint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Purpose, actual.Purpose, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Purpose")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Role, actual.Role, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateSubnetworkUpdateOperation")}, fn.AddNest("Role")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecondaryIPRanges, actual.SecondaryIPRanges, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareSubnetworkSecondaryIPRangesNewStyle, EmptyObject: EmptySubnetworkSecondaryIPRanges, OperationSelector: dcl.TriggersOperation("updateSubnetworkUpdateOperation")}, fn.AddNest("SecondaryIpRanges")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrivateIPGoogleAccess, actual.PrivateIPGoogleAccess, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateSubnetworkSetPrivateIpGoogleAccessOperation")}, fn.AddNest("PrivateIpGoogleAccess")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Region, actual.Region, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Region")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LogConfig, actual.LogConfig, dcl.DiffInfo{ObjectFunction: compareSubnetworkLogConfigNewStyle, EmptyObject: EmptySubnetworkLogConfig, OperationSelector: dcl.TriggersOperation("updateSubnetworkUpdateOperation")}, fn.AddNest("LogConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.SelfLink, actual.SelfLink, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SelfLink")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableFlowLogs, actual.EnableFlowLogs, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableFlowLogs")); len(ds) != 0 || err != nil {
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
func compareSubnetworkSecondaryIPRangesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*SubnetworkSecondaryIPRanges)
	if !ok {
		desiredNotPointer, ok := d.(SubnetworkSecondaryIPRanges)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SubnetworkSecondaryIPRanges or *SubnetworkSecondaryIPRanges", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*SubnetworkSecondaryIPRanges)
	if !ok {
		actualNotPointer, ok := a.(SubnetworkSecondaryIPRanges)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SubnetworkSecondaryIPRanges", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.RangeName, actual.RangeName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RangeName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPCidrRange, actual.IPCidrRange, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IpCidrRange")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareSubnetworkLogConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*SubnetworkLogConfig)
	if !ok {
		desiredNotPointer, ok := d.(SubnetworkLogConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SubnetworkLogConfig or *SubnetworkLogConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*SubnetworkLogConfig)
	if !ok {
		actualNotPointer, ok := a.(SubnetworkLogConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SubnetworkLogConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AggregationInterval, actual.AggregationInterval, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateSubnetworkUpdateOperation")}, fn.AddNest("AggregationInterval")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FlowSampling, actual.FlowSampling, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateSubnetworkUpdateOperation")}, fn.AddNest("FlowSampling")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metadata, actual.Metadata, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateSubnetworkUpdateOperation")}, fn.AddNest("Metadata")); len(ds) != 0 || err != nil {
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
func (r *Subnetwork) urlNormalized() *Subnetwork {
	normalized := dcl.Copy(*r).(Subnetwork)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.GatewayAddress = dcl.SelfLinkToName(r.GatewayAddress)
	normalized.IPCidrRange = dcl.SelfLinkToName(r.IPCidrRange)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Network = dcl.SelfLinkToName(r.Network)
	normalized.Fingerprint = dcl.SelfLinkToName(r.Fingerprint)
	normalized.Region = dcl.SelfLinkToName(r.Region)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	return &normalized
}

func (r *Subnetwork) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "expandIpCidrRange" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"region":  dcl.ValueOrEmptyString(nr.Region),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/regions/{{region}}/subnetworks/{{name}}/expandIpCidrRange", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "setPrivateIpGoogleAccess" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"region":  dcl.ValueOrEmptyString(nr.Region),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/regions/{{region}}/subnetworks/{{name}}/setPrivateIpGoogleAccess", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "update" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"region":  dcl.ValueOrEmptyString(nr.Region),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/regions/{{region}}/subnetworks/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Subnetwork resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Subnetwork) marshal(c *Client) ([]byte, error) {
	m, err := expandSubnetwork(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Subnetwork: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalSubnetwork decodes JSON responses into the Subnetwork resource schema.
func unmarshalSubnetwork(b []byte, c *Client, res *Subnetwork) (*Subnetwork, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapSubnetwork(m, c, res)
}

func unmarshalMapSubnetwork(m map[string]interface{}, c *Client, res *Subnetwork) (*Subnetwork, error) {

	flattened := flattenSubnetwork(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandSubnetwork expands Subnetwork into a JSON request object.
func expandSubnetwork(c *Client, f *Subnetwork) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.IPCidrRange; dcl.ValueShouldBeSent(v) {
		m["ipCidrRange"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Network; dcl.ValueShouldBeSent(v) {
		m["network"] = v
	}
	if v := f.Purpose; dcl.ValueShouldBeSent(v) {
		m["purpose"] = v
	}
	if v := f.Role; dcl.ValueShouldBeSent(v) {
		m["role"] = v
	}
	if v, err := expandSubnetworkSecondaryIPRangesSlice(c, f.SecondaryIPRanges, res); err != nil {
		return nil, fmt.Errorf("error expanding SecondaryIPRanges into secondaryIpRanges: %w", err)
	} else if v != nil {
		m["secondaryIpRanges"] = v
	}
	if v := f.PrivateIPGoogleAccess; dcl.ValueShouldBeSent(v) {
		m["privateIpGoogleAccess"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Region into region: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["region"] = v
	}
	if v, err := expandSubnetworkLogConfig(c, f.LogConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LogConfig into logConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["logConfig"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}
	if v := f.EnableFlowLogs; dcl.ValueShouldBeSent(v) {
		m["enableFlowLogs"] = v
	}

	return m, nil
}

// flattenSubnetwork flattens Subnetwork from a JSON request object into the
// Subnetwork type.
func flattenSubnetwork(c *Client, i interface{}, res *Subnetwork) *Subnetwork {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Subnetwork{}
	resultRes.CreationTimestamp = dcl.FlattenString(m["creationTimestamp"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.GatewayAddress = dcl.FlattenString(m["gatewayAddress"])
	resultRes.IPCidrRange = dcl.FlattenString(m["ipCidrRange"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Network = dcl.FlattenString(m["network"])
	resultRes.Fingerprint = dcl.FlattenString(m["fingerprint"])
	resultRes.Purpose = flattenSubnetworkPurposeEnum(m["purpose"])
	resultRes.Role = flattenSubnetworkRoleEnum(m["role"])
	resultRes.SecondaryIPRanges = flattenSubnetworkSecondaryIPRangesSlice(c, m["secondaryIpRanges"], res)
	resultRes.PrivateIPGoogleAccess = dcl.FlattenBool(m["privateIpGoogleAccess"])
	resultRes.Region = dcl.FlattenString(m["region"])
	resultRes.LogConfig = flattenSubnetworkLogConfig(c, m["logConfig"], res)
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.EnableFlowLogs = dcl.FlattenBool(m["enableFlowLogs"])

	return resultRes
}

// expandSubnetworkSecondaryIPRangesMap expands the contents of SubnetworkSecondaryIPRanges into a JSON
// request object.
func expandSubnetworkSecondaryIPRangesMap(c *Client, f map[string]SubnetworkSecondaryIPRanges, res *Subnetwork) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandSubnetworkSecondaryIPRanges(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandSubnetworkSecondaryIPRangesSlice expands the contents of SubnetworkSecondaryIPRanges into a JSON
// request object.
func expandSubnetworkSecondaryIPRangesSlice(c *Client, f []SubnetworkSecondaryIPRanges, res *Subnetwork) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandSubnetworkSecondaryIPRanges(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenSubnetworkSecondaryIPRangesMap flattens the contents of SubnetworkSecondaryIPRanges from a JSON
// response object.
func flattenSubnetworkSecondaryIPRangesMap(c *Client, i interface{}, res *Subnetwork) map[string]SubnetworkSecondaryIPRanges {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SubnetworkSecondaryIPRanges{}
	}

	if len(a) == 0 {
		return map[string]SubnetworkSecondaryIPRanges{}
	}

	items := make(map[string]SubnetworkSecondaryIPRanges)
	for k, item := range a {
		items[k] = *flattenSubnetworkSecondaryIPRanges(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenSubnetworkSecondaryIPRangesSlice flattens the contents of SubnetworkSecondaryIPRanges from a JSON
// response object.
func flattenSubnetworkSecondaryIPRangesSlice(c *Client, i interface{}, res *Subnetwork) []SubnetworkSecondaryIPRanges {
	a, ok := i.([]interface{})
	if !ok {
		return []SubnetworkSecondaryIPRanges{}
	}

	if len(a) == 0 {
		return []SubnetworkSecondaryIPRanges{}
	}

	items := make([]SubnetworkSecondaryIPRanges, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSubnetworkSecondaryIPRanges(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandSubnetworkSecondaryIPRanges expands an instance of SubnetworkSecondaryIPRanges into a JSON
// request object.
func expandSubnetworkSecondaryIPRanges(c *Client, f *SubnetworkSecondaryIPRanges, res *Subnetwork) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.RangeName; !dcl.IsEmptyValueIndirect(v) {
		m["rangeName"] = v
	}
	if v := f.IPCidrRange; !dcl.IsEmptyValueIndirect(v) {
		m["ipCidrRange"] = v
	}

	return m, nil
}

// flattenSubnetworkSecondaryIPRanges flattens an instance of SubnetworkSecondaryIPRanges from a JSON
// response object.
func flattenSubnetworkSecondaryIPRanges(c *Client, i interface{}, res *Subnetwork) *SubnetworkSecondaryIPRanges {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &SubnetworkSecondaryIPRanges{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptySubnetworkSecondaryIPRanges
	}
	r.RangeName = dcl.FlattenString(m["rangeName"])
	r.IPCidrRange = dcl.FlattenString(m["ipCidrRange"])

	return r
}

// expandSubnetworkLogConfigMap expands the contents of SubnetworkLogConfig into a JSON
// request object.
func expandSubnetworkLogConfigMap(c *Client, f map[string]SubnetworkLogConfig, res *Subnetwork) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandSubnetworkLogConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandSubnetworkLogConfigSlice expands the contents of SubnetworkLogConfig into a JSON
// request object.
func expandSubnetworkLogConfigSlice(c *Client, f []SubnetworkLogConfig, res *Subnetwork) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandSubnetworkLogConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenSubnetworkLogConfigMap flattens the contents of SubnetworkLogConfig from a JSON
// response object.
func flattenSubnetworkLogConfigMap(c *Client, i interface{}, res *Subnetwork) map[string]SubnetworkLogConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SubnetworkLogConfig{}
	}

	if len(a) == 0 {
		return map[string]SubnetworkLogConfig{}
	}

	items := make(map[string]SubnetworkLogConfig)
	for k, item := range a {
		items[k] = *flattenSubnetworkLogConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenSubnetworkLogConfigSlice flattens the contents of SubnetworkLogConfig from a JSON
// response object.
func flattenSubnetworkLogConfigSlice(c *Client, i interface{}, res *Subnetwork) []SubnetworkLogConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []SubnetworkLogConfig{}
	}

	if len(a) == 0 {
		return []SubnetworkLogConfig{}
	}

	items := make([]SubnetworkLogConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSubnetworkLogConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandSubnetworkLogConfig expands an instance of SubnetworkLogConfig into a JSON
// request object.
func expandSubnetworkLogConfig(c *Client, f *SubnetworkLogConfig, res *Subnetwork) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AggregationInterval; !dcl.IsEmptyValueIndirect(v) {
		m["aggregationInterval"] = v
	}
	if v := f.FlowSampling; !dcl.IsEmptyValueIndirect(v) {
		m["flowSampling"] = v
	}
	if v := f.Metadata; !dcl.IsEmptyValueIndirect(v) {
		m["metadata"] = v
	}

	return m, nil
}

// flattenSubnetworkLogConfig flattens an instance of SubnetworkLogConfig from a JSON
// response object.
func flattenSubnetworkLogConfig(c *Client, i interface{}, res *Subnetwork) *SubnetworkLogConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &SubnetworkLogConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptySubnetworkLogConfig
	}
	r.AggregationInterval = flattenSubnetworkLogConfigAggregationIntervalEnum(m["aggregationInterval"])
	if dcl.IsEmptyValueIndirect(m["aggregationInterval"]) {
		c.Config.Logger.Info("Using default value for aggregationInterval.")
		r.AggregationInterval = SubnetworkLogConfigAggregationIntervalEnumRef("INTERVAL_5_SEC")
	}
	r.FlowSampling = dcl.FlattenDouble(m["flowSampling"])
	if dcl.IsEmptyValueIndirect(m["flowSampling"]) {
		c.Config.Logger.Info("Using default value for flowSampling.")
		r.FlowSampling = dcl.Float64(0.5)
	}
	r.Metadata = flattenSubnetworkLogConfigMetadataEnum(m["metadata"])
	if dcl.IsEmptyValueIndirect(m["metadata"]) {
		c.Config.Logger.Info("Using default value for metadata.")
		r.Metadata = SubnetworkLogConfigMetadataEnumRef("INCLUDE_ALL_METADATA")
	}

	return r
}

// flattenSubnetworkPurposeEnumMap flattens the contents of SubnetworkPurposeEnum from a JSON
// response object.
func flattenSubnetworkPurposeEnumMap(c *Client, i interface{}, res *Subnetwork) map[string]SubnetworkPurposeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SubnetworkPurposeEnum{}
	}

	if len(a) == 0 {
		return map[string]SubnetworkPurposeEnum{}
	}

	items := make(map[string]SubnetworkPurposeEnum)
	for k, item := range a {
		items[k] = *flattenSubnetworkPurposeEnum(item.(interface{}))
	}

	return items
}

// flattenSubnetworkPurposeEnumSlice flattens the contents of SubnetworkPurposeEnum from a JSON
// response object.
func flattenSubnetworkPurposeEnumSlice(c *Client, i interface{}, res *Subnetwork) []SubnetworkPurposeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []SubnetworkPurposeEnum{}
	}

	if len(a) == 0 {
		return []SubnetworkPurposeEnum{}
	}

	items := make([]SubnetworkPurposeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSubnetworkPurposeEnum(item.(interface{})))
	}

	return items
}

// flattenSubnetworkPurposeEnum asserts that an interface is a string, and returns a
// pointer to a *SubnetworkPurposeEnum with the same value as that string.
func flattenSubnetworkPurposeEnum(i interface{}) *SubnetworkPurposeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return SubnetworkPurposeEnumRef(s)
}

// flattenSubnetworkRoleEnumMap flattens the contents of SubnetworkRoleEnum from a JSON
// response object.
func flattenSubnetworkRoleEnumMap(c *Client, i interface{}, res *Subnetwork) map[string]SubnetworkRoleEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SubnetworkRoleEnum{}
	}

	if len(a) == 0 {
		return map[string]SubnetworkRoleEnum{}
	}

	items := make(map[string]SubnetworkRoleEnum)
	for k, item := range a {
		items[k] = *flattenSubnetworkRoleEnum(item.(interface{}))
	}

	return items
}

// flattenSubnetworkRoleEnumSlice flattens the contents of SubnetworkRoleEnum from a JSON
// response object.
func flattenSubnetworkRoleEnumSlice(c *Client, i interface{}, res *Subnetwork) []SubnetworkRoleEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []SubnetworkRoleEnum{}
	}

	if len(a) == 0 {
		return []SubnetworkRoleEnum{}
	}

	items := make([]SubnetworkRoleEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSubnetworkRoleEnum(item.(interface{})))
	}

	return items
}

// flattenSubnetworkRoleEnum asserts that an interface is a string, and returns a
// pointer to a *SubnetworkRoleEnum with the same value as that string.
func flattenSubnetworkRoleEnum(i interface{}) *SubnetworkRoleEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return SubnetworkRoleEnumRef(s)
}

// flattenSubnetworkLogConfigAggregationIntervalEnumMap flattens the contents of SubnetworkLogConfigAggregationIntervalEnum from a JSON
// response object.
func flattenSubnetworkLogConfigAggregationIntervalEnumMap(c *Client, i interface{}, res *Subnetwork) map[string]SubnetworkLogConfigAggregationIntervalEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SubnetworkLogConfigAggregationIntervalEnum{}
	}

	if len(a) == 0 {
		return map[string]SubnetworkLogConfigAggregationIntervalEnum{}
	}

	items := make(map[string]SubnetworkLogConfigAggregationIntervalEnum)
	for k, item := range a {
		items[k] = *flattenSubnetworkLogConfigAggregationIntervalEnum(item.(interface{}))
	}

	return items
}

// flattenSubnetworkLogConfigAggregationIntervalEnumSlice flattens the contents of SubnetworkLogConfigAggregationIntervalEnum from a JSON
// response object.
func flattenSubnetworkLogConfigAggregationIntervalEnumSlice(c *Client, i interface{}, res *Subnetwork) []SubnetworkLogConfigAggregationIntervalEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []SubnetworkLogConfigAggregationIntervalEnum{}
	}

	if len(a) == 0 {
		return []SubnetworkLogConfigAggregationIntervalEnum{}
	}

	items := make([]SubnetworkLogConfigAggregationIntervalEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSubnetworkLogConfigAggregationIntervalEnum(item.(interface{})))
	}

	return items
}

// flattenSubnetworkLogConfigAggregationIntervalEnum asserts that an interface is a string, and returns a
// pointer to a *SubnetworkLogConfigAggregationIntervalEnum with the same value as that string.
func flattenSubnetworkLogConfigAggregationIntervalEnum(i interface{}) *SubnetworkLogConfigAggregationIntervalEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return SubnetworkLogConfigAggregationIntervalEnumRef(s)
}

// flattenSubnetworkLogConfigMetadataEnumMap flattens the contents of SubnetworkLogConfigMetadataEnum from a JSON
// response object.
func flattenSubnetworkLogConfigMetadataEnumMap(c *Client, i interface{}, res *Subnetwork) map[string]SubnetworkLogConfigMetadataEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SubnetworkLogConfigMetadataEnum{}
	}

	if len(a) == 0 {
		return map[string]SubnetworkLogConfigMetadataEnum{}
	}

	items := make(map[string]SubnetworkLogConfigMetadataEnum)
	for k, item := range a {
		items[k] = *flattenSubnetworkLogConfigMetadataEnum(item.(interface{}))
	}

	return items
}

// flattenSubnetworkLogConfigMetadataEnumSlice flattens the contents of SubnetworkLogConfigMetadataEnum from a JSON
// response object.
func flattenSubnetworkLogConfigMetadataEnumSlice(c *Client, i interface{}, res *Subnetwork) []SubnetworkLogConfigMetadataEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []SubnetworkLogConfigMetadataEnum{}
	}

	if len(a) == 0 {
		return []SubnetworkLogConfigMetadataEnum{}
	}

	items := make([]SubnetworkLogConfigMetadataEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSubnetworkLogConfigMetadataEnum(item.(interface{})))
	}

	return items
}

// flattenSubnetworkLogConfigMetadataEnum asserts that an interface is a string, and returns a
// pointer to a *SubnetworkLogConfigMetadataEnum with the same value as that string.
func flattenSubnetworkLogConfigMetadataEnum(i interface{}) *SubnetworkLogConfigMetadataEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return SubnetworkLogConfigMetadataEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Subnetwork) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalSubnetwork(b, c, r)
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
		if nr.Region == nil && ncr.Region == nil {
			c.Config.Logger.Info("Both Region fields null - considering equal.")
		} else if nr.Region == nil || ncr.Region == nil {
			c.Config.Logger.Info("Only one Region field is null - considering unequal.")
			return false
		} else if *nr.Region != *ncr.Region {
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

type subnetworkDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         subnetworkApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToSubnetworkDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]subnetworkDiff, error) {
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
	var diffs []subnetworkDiff
	// For each operation name, create a subnetworkDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := subnetworkDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToSubnetworkApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToSubnetworkApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (subnetworkApiOperation, error) {
	switch opName {

	case "updateSubnetworkExpandIpCidrRangeOperation":
		return &updateSubnetworkExpandIpCidrRangeOperation{FieldDiffs: fieldDiffs}, nil

	case "updateSubnetworkSetPrivateIpGoogleAccessOperation":
		return &updateSubnetworkSetPrivateIpGoogleAccessOperation{FieldDiffs: fieldDiffs}, nil

	case "updateSubnetworkUpdateOperation":
		return &updateSubnetworkUpdateOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractSubnetworkFields(r *Subnetwork) error {
	vLogConfig := r.LogConfig
	if vLogConfig == nil {
		// note: explicitly not the empty object.
		vLogConfig = &SubnetworkLogConfig{}
	}
	if err := extractSubnetworkLogConfigFields(r, vLogConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLogConfig) {
		r.LogConfig = vLogConfig
	}
	return nil
}
func extractSubnetworkSecondaryIPRangesFields(r *Subnetwork, o *SubnetworkSecondaryIPRanges) error {
	return nil
}
func extractSubnetworkLogConfigFields(r *Subnetwork, o *SubnetworkLogConfig) error {
	return nil
}

func postReadExtractSubnetworkFields(r *Subnetwork) error {
	vLogConfig := r.LogConfig
	if vLogConfig == nil {
		// note: explicitly not the empty object.
		vLogConfig = &SubnetworkLogConfig{}
	}
	if err := postReadExtractSubnetworkLogConfigFields(r, vLogConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLogConfig) {
		r.LogConfig = vLogConfig
	}
	return nil
}
func postReadExtractSubnetworkSecondaryIPRangesFields(r *Subnetwork, o *SubnetworkSecondaryIPRanges) error {
	return nil
}
func postReadExtractSubnetworkLogConfigFields(r *Subnetwork, o *SubnetworkLogConfig) error {
	return nil
}
