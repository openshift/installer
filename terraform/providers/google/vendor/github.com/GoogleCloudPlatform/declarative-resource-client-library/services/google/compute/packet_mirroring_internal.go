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

func (r *PacketMirroring) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "network"); err != nil {
		return err
	}
	if err := dcl.Required(r, "collectorIlb"); err != nil {
		return err
	}
	if err := dcl.Required(r, "mirroredResources"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Network) {
		if err := r.Network.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.CollectorIlb) {
		if err := r.CollectorIlb.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.MirroredResources) {
		if err := r.MirroredResources.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Filter) {
		if err := r.Filter.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *PacketMirroringNetwork) validate() error {
	if err := dcl.Required(r, "url"); err != nil {
		return err
	}
	return nil
}
func (r *PacketMirroringCollectorIlb) validate() error {
	if err := dcl.Required(r, "url"); err != nil {
		return err
	}
	return nil
}
func (r *PacketMirroringMirroredResources) validate() error {
	return nil
}
func (r *PacketMirroringMirroredResourcesSubnetworks) validate() error {
	return nil
}
func (r *PacketMirroringMirroredResourcesInstances) validate() error {
	return nil
}
func (r *PacketMirroringFilter) validate() error {
	return nil
}
func (r *PacketMirroring) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *PacketMirroring) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/packetMirrorings/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *PacketMirroring) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/packetMirrorings", nr.basePath(), userBasePath, params), nil

}

func (r *PacketMirroring) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/packetMirrorings", nr.basePath(), userBasePath, params), nil

}

func (r *PacketMirroring) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/packetMirrorings/{{name}}", nr.basePath(), userBasePath, params), nil
}

// packetMirroringApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type packetMirroringApiOperation interface {
	do(context.Context, *PacketMirroring, *Client) error
}

// newUpdatePacketMirroringPatchRequest creates a request for an
// PacketMirroring resource's Patch update type by filling in the update
// fields based on the intended state of the resource.
func newUpdatePacketMirroringPatchRequest(ctx context.Context, f *PacketMirroring, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.Priority; !dcl.IsEmptyValueIndirect(v) {
		req["priority"] = v
	}
	if v, err := expandPacketMirroringCollectorIlb(c, f.CollectorIlb, res); err != nil {
		return nil, fmt.Errorf("error expanding CollectorIlb into collectorIlb: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["collectorIlb"] = v
	}
	if v, err := expandPacketMirroringMirroredResources(c, f.MirroredResources, res); err != nil {
		return nil, fmt.Errorf("error expanding MirroredResources into mirroredResources: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["mirroredResources"] = v
	}
	if v, err := expandPacketMirroringFilter(c, f.Filter, res); err != nil {
		return nil, fmt.Errorf("error expanding Filter into filter: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["filter"] = v
	}
	if v := f.Enable; !dcl.IsEmptyValueIndirect(v) {
		req["enable"] = v
	}
	return req, nil
}

// marshalUpdatePacketMirroringPatchRequest converts the update into
// the final JSON request body.
func marshalUpdatePacketMirroringPatchRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updatePacketMirroringPatchOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updatePacketMirroringPatchOperation) do(ctx context.Context, r *PacketMirroring, c *Client) error {
	_, err := c.GetPacketMirroring(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "Patch")
	if err != nil {
		return err
	}

	req, err := newUpdatePacketMirroringPatchRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdatePacketMirroringPatchRequest(c, req)
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

func (c *Client) listPacketMirroringRaw(ctx context.Context, r *PacketMirroring, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != PacketMirroringMaxPage {
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

type listPacketMirroringOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listPacketMirroring(ctx context.Context, r *PacketMirroring, pageToken string, pageSize int32) ([]*PacketMirroring, string, error) {
	b, err := c.listPacketMirroringRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listPacketMirroringOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*PacketMirroring
	for _, v := range m.Items {
		res, err := unmarshalMapPacketMirroring(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllPacketMirroring(ctx context.Context, f func(*PacketMirroring) bool, resources []*PacketMirroring) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeletePacketMirroring(ctx, res)
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

type deletePacketMirroringOperation struct{}

func (op *deletePacketMirroringOperation) do(ctx context.Context, r *PacketMirroring, c *Client) error {
	r, err := c.GetPacketMirroring(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "PacketMirroring not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetPacketMirroring checking for existence. error: %v", err)
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
		_, err := c.GetPacketMirroring(ctx, r)
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
type createPacketMirroringOperation struct {
	response map[string]interface{}
}

func (op *createPacketMirroringOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createPacketMirroringOperation) do(ctx context.Context, r *PacketMirroring, c *Client) error {
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

	if _, err := c.GetPacketMirroring(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getPacketMirroringRaw(ctx context.Context, r *PacketMirroring) ([]byte, error) {

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

func (c *Client) packetMirroringDiffsForRawDesired(ctx context.Context, rawDesired *PacketMirroring, opts ...dcl.ApplyOption) (initial, desired *PacketMirroring, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *PacketMirroring
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*PacketMirroring); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected PacketMirroring, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetPacketMirroring(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a PacketMirroring resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve PacketMirroring resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that PacketMirroring resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizePacketMirroringDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for PacketMirroring: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for PacketMirroring: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractPacketMirroringFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizePacketMirroringInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for PacketMirroring: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizePacketMirroringDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for PacketMirroring: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffPacketMirroring(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizePacketMirroringInitialState(rawInitial, rawDesired *PacketMirroring) (*PacketMirroring, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizePacketMirroringDesiredState(rawDesired, rawInitial *PacketMirroring, opts ...dcl.ApplyOption) (*PacketMirroring, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Network = canonicalizePacketMirroringNetwork(rawDesired.Network, nil, opts...)
		rawDesired.CollectorIlb = canonicalizePacketMirroringCollectorIlb(rawDesired.CollectorIlb, nil, opts...)
		rawDesired.MirroredResources = canonicalizePacketMirroringMirroredResources(rawDesired.MirroredResources, nil, opts...)
		rawDesired.Filter = canonicalizePacketMirroringFilter(rawDesired.Filter, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &PacketMirroring{}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	canonicalDesired.Network = canonicalizePacketMirroringNetwork(rawDesired.Network, rawInitial.Network, opts...)
	if dcl.IsZeroValue(rawDesired.Priority) || (dcl.IsEmptyValueIndirect(rawDesired.Priority) && dcl.IsEmptyValueIndirect(rawInitial.Priority)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Priority = rawInitial.Priority
	} else {
		canonicalDesired.Priority = rawDesired.Priority
	}
	canonicalDesired.CollectorIlb = canonicalizePacketMirroringCollectorIlb(rawDesired.CollectorIlb, rawInitial.CollectorIlb, opts...)
	canonicalDesired.MirroredResources = canonicalizePacketMirroringMirroredResources(rawDesired.MirroredResources, rawInitial.MirroredResources, opts...)
	canonicalDesired.Filter = canonicalizePacketMirroringFilter(rawDesired.Filter, rawInitial.Filter, opts...)
	if dcl.IsZeroValue(rawDesired.Enable) || (dcl.IsEmptyValueIndirect(rawDesired.Enable) && dcl.IsEmptyValueIndirect(rawInitial.Enable)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Enable = rawInitial.Enable
	} else {
		canonicalDesired.Enable = rawDesired.Enable
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

func canonicalizePacketMirroringNewState(c *Client, rawNew, rawDesired *PacketMirroring) (*PacketMirroring, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Id) && dcl.IsEmptyValueIndirect(rawDesired.Id) {
		rawNew.Id = rawDesired.Id
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
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

	if dcl.IsEmptyValueIndirect(rawNew.Region) && dcl.IsEmptyValueIndirect(rawDesired.Region) {
		rawNew.Region = rawDesired.Region
	} else {
		if dcl.StringCanonicalize(rawDesired.Region, rawNew.Region) {
			rawNew.Region = rawDesired.Region
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Network) && dcl.IsEmptyValueIndirect(rawDesired.Network) {
		rawNew.Network = rawDesired.Network
	} else {
		rawNew.Network = canonicalizeNewPacketMirroringNetwork(c, rawDesired.Network, rawNew.Network)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Priority) && dcl.IsEmptyValueIndirect(rawDesired.Priority) {
		rawNew.Priority = rawDesired.Priority
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.CollectorIlb) && dcl.IsEmptyValueIndirect(rawDesired.CollectorIlb) {
		rawNew.CollectorIlb = rawDesired.CollectorIlb
	} else {
		rawNew.CollectorIlb = canonicalizeNewPacketMirroringCollectorIlb(c, rawDesired.CollectorIlb, rawNew.CollectorIlb)
	}

	if dcl.IsEmptyValueIndirect(rawNew.MirroredResources) && dcl.IsEmptyValueIndirect(rawDesired.MirroredResources) {
		rawNew.MirroredResources = rawDesired.MirroredResources
	} else {
		rawNew.MirroredResources = canonicalizeNewPacketMirroringMirroredResources(c, rawDesired.MirroredResources, rawNew.MirroredResources)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Filter) && dcl.IsEmptyValueIndirect(rawDesired.Filter) {
		rawNew.Filter = rawDesired.Filter
	} else {
		rawNew.Filter = canonicalizeNewPacketMirroringFilter(c, rawDesired.Filter, rawNew.Filter)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Enable) && dcl.IsEmptyValueIndirect(rawDesired.Enable) {
		rawNew.Enable = rawDesired.Enable
	} else {
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizePacketMirroringNetwork(des, initial *PacketMirroringNetwork, opts ...dcl.ApplyOption) *PacketMirroringNetwork {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PacketMirroringNetwork{}

	if dcl.IsZeroValue(des.Url) || (dcl.IsEmptyValueIndirect(des.Url) && dcl.IsEmptyValueIndirect(initial.Url)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Url = initial.Url
	} else {
		cDes.Url = des.Url
	}

	return cDes
}

func canonicalizePacketMirroringNetworkSlice(des, initial []PacketMirroringNetwork, opts ...dcl.ApplyOption) []PacketMirroringNetwork {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PacketMirroringNetwork, 0, len(des))
		for _, d := range des {
			cd := canonicalizePacketMirroringNetwork(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PacketMirroringNetwork, 0, len(des))
	for i, d := range des {
		cd := canonicalizePacketMirroringNetwork(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPacketMirroringNetwork(c *Client, des, nw *PacketMirroringNetwork) *PacketMirroringNetwork {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PacketMirroringNetwork while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.CanonicalUrl, nw.CanonicalUrl) {
		nw.CanonicalUrl = des.CanonicalUrl
	}

	return nw
}

func canonicalizeNewPacketMirroringNetworkSet(c *Client, des, nw []PacketMirroringNetwork) []PacketMirroringNetwork {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PacketMirroringNetwork
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePacketMirroringNetworkNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPacketMirroringNetwork(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPacketMirroringNetworkSlice(c *Client, des, nw []PacketMirroringNetwork) []PacketMirroringNetwork {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PacketMirroringNetwork
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPacketMirroringNetwork(c, &d, &n))
	}

	return items
}

func canonicalizePacketMirroringCollectorIlb(des, initial *PacketMirroringCollectorIlb, opts ...dcl.ApplyOption) *PacketMirroringCollectorIlb {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PacketMirroringCollectorIlb{}

	if dcl.IsZeroValue(des.Url) || (dcl.IsEmptyValueIndirect(des.Url) && dcl.IsEmptyValueIndirect(initial.Url)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Url = initial.Url
	} else {
		cDes.Url = des.Url
	}

	return cDes
}

func canonicalizePacketMirroringCollectorIlbSlice(des, initial []PacketMirroringCollectorIlb, opts ...dcl.ApplyOption) []PacketMirroringCollectorIlb {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PacketMirroringCollectorIlb, 0, len(des))
		for _, d := range des {
			cd := canonicalizePacketMirroringCollectorIlb(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PacketMirroringCollectorIlb, 0, len(des))
	for i, d := range des {
		cd := canonicalizePacketMirroringCollectorIlb(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPacketMirroringCollectorIlb(c *Client, des, nw *PacketMirroringCollectorIlb) *PacketMirroringCollectorIlb {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PacketMirroringCollectorIlb while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.CanonicalUrl, nw.CanonicalUrl) {
		nw.CanonicalUrl = des.CanonicalUrl
	}

	return nw
}

func canonicalizeNewPacketMirroringCollectorIlbSet(c *Client, des, nw []PacketMirroringCollectorIlb) []PacketMirroringCollectorIlb {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PacketMirroringCollectorIlb
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePacketMirroringCollectorIlbNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPacketMirroringCollectorIlb(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPacketMirroringCollectorIlbSlice(c *Client, des, nw []PacketMirroringCollectorIlb) []PacketMirroringCollectorIlb {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PacketMirroringCollectorIlb
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPacketMirroringCollectorIlb(c, &d, &n))
	}

	return items
}

func canonicalizePacketMirroringMirroredResources(des, initial *PacketMirroringMirroredResources, opts ...dcl.ApplyOption) *PacketMirroringMirroredResources {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PacketMirroringMirroredResources{}

	cDes.Subnetworks = canonicalizePacketMirroringMirroredResourcesSubnetworksSlice(des.Subnetworks, initial.Subnetworks, opts...)
	cDes.Instances = canonicalizePacketMirroringMirroredResourcesInstancesSlice(des.Instances, initial.Instances, opts...)
	if dcl.StringArrayCanonicalize(des.Tags, initial.Tags) {
		cDes.Tags = initial.Tags
	} else {
		cDes.Tags = des.Tags
	}

	return cDes
}

func canonicalizePacketMirroringMirroredResourcesSlice(des, initial []PacketMirroringMirroredResources, opts ...dcl.ApplyOption) []PacketMirroringMirroredResources {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PacketMirroringMirroredResources, 0, len(des))
		for _, d := range des {
			cd := canonicalizePacketMirroringMirroredResources(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PacketMirroringMirroredResources, 0, len(des))
	for i, d := range des {
		cd := canonicalizePacketMirroringMirroredResources(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPacketMirroringMirroredResources(c *Client, des, nw *PacketMirroringMirroredResources) *PacketMirroringMirroredResources {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PacketMirroringMirroredResources while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Subnetworks = canonicalizeNewPacketMirroringMirroredResourcesSubnetworksSlice(c, des.Subnetworks, nw.Subnetworks)
	nw.Instances = canonicalizeNewPacketMirroringMirroredResourcesInstancesSlice(c, des.Instances, nw.Instances)
	if dcl.StringArrayCanonicalize(des.Tags, nw.Tags) {
		nw.Tags = des.Tags
	}

	return nw
}

func canonicalizeNewPacketMirroringMirroredResourcesSet(c *Client, des, nw []PacketMirroringMirroredResources) []PacketMirroringMirroredResources {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PacketMirroringMirroredResources
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePacketMirroringMirroredResourcesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPacketMirroringMirroredResources(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPacketMirroringMirroredResourcesSlice(c *Client, des, nw []PacketMirroringMirroredResources) []PacketMirroringMirroredResources {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PacketMirroringMirroredResources
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPacketMirroringMirroredResources(c, &d, &n))
	}

	return items
}

func canonicalizePacketMirroringMirroredResourcesSubnetworks(des, initial *PacketMirroringMirroredResourcesSubnetworks, opts ...dcl.ApplyOption) *PacketMirroringMirroredResourcesSubnetworks {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PacketMirroringMirroredResourcesSubnetworks{}

	if dcl.IsZeroValue(des.Url) || (dcl.IsEmptyValueIndirect(des.Url) && dcl.IsEmptyValueIndirect(initial.Url)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Url = initial.Url
	} else {
		cDes.Url = des.Url
	}

	return cDes
}

func canonicalizePacketMirroringMirroredResourcesSubnetworksSlice(des, initial []PacketMirroringMirroredResourcesSubnetworks, opts ...dcl.ApplyOption) []PacketMirroringMirroredResourcesSubnetworks {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PacketMirroringMirroredResourcesSubnetworks, 0, len(des))
		for _, d := range des {
			cd := canonicalizePacketMirroringMirroredResourcesSubnetworks(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PacketMirroringMirroredResourcesSubnetworks, 0, len(des))
	for i, d := range des {
		cd := canonicalizePacketMirroringMirroredResourcesSubnetworks(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPacketMirroringMirroredResourcesSubnetworks(c *Client, des, nw *PacketMirroringMirroredResourcesSubnetworks) *PacketMirroringMirroredResourcesSubnetworks {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PacketMirroringMirroredResourcesSubnetworks while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.CanonicalUrl, nw.CanonicalUrl) {
		nw.CanonicalUrl = des.CanonicalUrl
	}

	return nw
}

func canonicalizeNewPacketMirroringMirroredResourcesSubnetworksSet(c *Client, des, nw []PacketMirroringMirroredResourcesSubnetworks) []PacketMirroringMirroredResourcesSubnetworks {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PacketMirroringMirroredResourcesSubnetworks
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePacketMirroringMirroredResourcesSubnetworksNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPacketMirroringMirroredResourcesSubnetworks(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPacketMirroringMirroredResourcesSubnetworksSlice(c *Client, des, nw []PacketMirroringMirroredResourcesSubnetworks) []PacketMirroringMirroredResourcesSubnetworks {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PacketMirroringMirroredResourcesSubnetworks
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPacketMirroringMirroredResourcesSubnetworks(c, &d, &n))
	}

	return items
}

func canonicalizePacketMirroringMirroredResourcesInstances(des, initial *PacketMirroringMirroredResourcesInstances, opts ...dcl.ApplyOption) *PacketMirroringMirroredResourcesInstances {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PacketMirroringMirroredResourcesInstances{}

	if dcl.IsZeroValue(des.Url) || (dcl.IsEmptyValueIndirect(des.Url) && dcl.IsEmptyValueIndirect(initial.Url)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Url = initial.Url
	} else {
		cDes.Url = des.Url
	}

	return cDes
}

func canonicalizePacketMirroringMirroredResourcesInstancesSlice(des, initial []PacketMirroringMirroredResourcesInstances, opts ...dcl.ApplyOption) []PacketMirroringMirroredResourcesInstances {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PacketMirroringMirroredResourcesInstances, 0, len(des))
		for _, d := range des {
			cd := canonicalizePacketMirroringMirroredResourcesInstances(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PacketMirroringMirroredResourcesInstances, 0, len(des))
	for i, d := range des {
		cd := canonicalizePacketMirroringMirroredResourcesInstances(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPacketMirroringMirroredResourcesInstances(c *Client, des, nw *PacketMirroringMirroredResourcesInstances) *PacketMirroringMirroredResourcesInstances {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PacketMirroringMirroredResourcesInstances while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.CanonicalUrl, nw.CanonicalUrl) {
		nw.CanonicalUrl = des.CanonicalUrl
	}

	return nw
}

func canonicalizeNewPacketMirroringMirroredResourcesInstancesSet(c *Client, des, nw []PacketMirroringMirroredResourcesInstances) []PacketMirroringMirroredResourcesInstances {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PacketMirroringMirroredResourcesInstances
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePacketMirroringMirroredResourcesInstancesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPacketMirroringMirroredResourcesInstances(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPacketMirroringMirroredResourcesInstancesSlice(c *Client, des, nw []PacketMirroringMirroredResourcesInstances) []PacketMirroringMirroredResourcesInstances {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PacketMirroringMirroredResourcesInstances
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPacketMirroringMirroredResourcesInstances(c, &d, &n))
	}

	return items
}

func canonicalizePacketMirroringFilter(des, initial *PacketMirroringFilter, opts ...dcl.ApplyOption) *PacketMirroringFilter {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &PacketMirroringFilter{}

	if dcl.StringArrayCanonicalize(des.CidrRanges, initial.CidrRanges) {
		cDes.CidrRanges = initial.CidrRanges
	} else {
		cDes.CidrRanges = des.CidrRanges
	}
	if dcl.StringArrayCanonicalize(des.IPProtocols, initial.IPProtocols) {
		cDes.IPProtocols = initial.IPProtocols
	} else {
		cDes.IPProtocols = des.IPProtocols
	}
	if dcl.IsZeroValue(des.Direction) || (dcl.IsEmptyValueIndirect(des.Direction) && dcl.IsEmptyValueIndirect(initial.Direction)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Direction = initial.Direction
	} else {
		cDes.Direction = des.Direction
	}

	return cDes
}

func canonicalizePacketMirroringFilterSlice(des, initial []PacketMirroringFilter, opts ...dcl.ApplyOption) []PacketMirroringFilter {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]PacketMirroringFilter, 0, len(des))
		for _, d := range des {
			cd := canonicalizePacketMirroringFilter(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]PacketMirroringFilter, 0, len(des))
	for i, d := range des {
		cd := canonicalizePacketMirroringFilter(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewPacketMirroringFilter(c *Client, des, nw *PacketMirroringFilter) *PacketMirroringFilter {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for PacketMirroringFilter while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.CidrRanges, nw.CidrRanges) {
		nw.CidrRanges = des.CidrRanges
	}
	if dcl.StringArrayCanonicalize(des.IPProtocols, nw.IPProtocols) {
		nw.IPProtocols = des.IPProtocols
	}

	return nw
}

func canonicalizeNewPacketMirroringFilterSet(c *Client, des, nw []PacketMirroringFilter) []PacketMirroringFilter {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []PacketMirroringFilter
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := comparePacketMirroringFilterNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewPacketMirroringFilter(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewPacketMirroringFilterSlice(c *Client, des, nw []PacketMirroringFilter) []PacketMirroringFilter {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []PacketMirroringFilter
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewPacketMirroringFilter(c, &d, &n))
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
func diffPacketMirroring(c *Client, desired, actual *PacketMirroring, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Network, actual.Network, dcl.DiffInfo{ObjectFunction: comparePacketMirroringNetworkNewStyle, EmptyObject: EmptyPacketMirroringNetwork, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Network")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Priority, actual.Priority, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Priority")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CollectorIlb, actual.CollectorIlb, dcl.DiffInfo{ObjectFunction: comparePacketMirroringCollectorIlbNewStyle, EmptyObject: EmptyPacketMirroringCollectorIlb, OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("CollectorIlb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MirroredResources, actual.MirroredResources, dcl.DiffInfo{ObjectFunction: comparePacketMirroringMirroredResourcesNewStyle, EmptyObject: EmptyPacketMirroringMirroredResources, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MirroredResources")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Filter, actual.Filter, dcl.DiffInfo{ServerDefault: true, ObjectFunction: comparePacketMirroringFilterNewStyle, EmptyObject: EmptyPacketMirroringFilter, OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Filter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Enable, actual.Enable, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Enable")); len(ds) != 0 || err != nil {
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
func comparePacketMirroringNetworkNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PacketMirroringNetwork)
	if !ok {
		desiredNotPointer, ok := d.(PacketMirroringNetwork)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringNetwork or *PacketMirroringNetwork", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PacketMirroringNetwork)
	if !ok {
		actualNotPointer, ok := a.(PacketMirroringNetwork)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringNetwork", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Url, actual.Url, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Url")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CanonicalUrl, actual.CanonicalUrl, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CanonicalUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePacketMirroringCollectorIlbNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PacketMirroringCollectorIlb)
	if !ok {
		desiredNotPointer, ok := d.(PacketMirroringCollectorIlb)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringCollectorIlb or *PacketMirroringCollectorIlb", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PacketMirroringCollectorIlb)
	if !ok {
		actualNotPointer, ok := a.(PacketMirroringCollectorIlb)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringCollectorIlb", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Url, actual.Url, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Url")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CanonicalUrl, actual.CanonicalUrl, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CanonicalUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePacketMirroringMirroredResourcesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PacketMirroringMirroredResources)
	if !ok {
		desiredNotPointer, ok := d.(PacketMirroringMirroredResources)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringMirroredResources or *PacketMirroringMirroredResources", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PacketMirroringMirroredResources)
	if !ok {
		actualNotPointer, ok := a.(PacketMirroringMirroredResources)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringMirroredResources", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Subnetworks, actual.Subnetworks, dcl.DiffInfo{ObjectFunction: comparePacketMirroringMirroredResourcesSubnetworksNewStyle, EmptyObject: EmptyPacketMirroringMirroredResourcesSubnetworks, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Subnetworks")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Instances, actual.Instances, dcl.DiffInfo{ObjectFunction: comparePacketMirroringMirroredResourcesInstancesNewStyle, EmptyObject: EmptyPacketMirroringMirroredResourcesInstances, OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Instances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tags, actual.Tags, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePacketMirroringMirroredResourcesSubnetworksNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PacketMirroringMirroredResourcesSubnetworks)
	if !ok {
		desiredNotPointer, ok := d.(PacketMirroringMirroredResourcesSubnetworks)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringMirroredResourcesSubnetworks or *PacketMirroringMirroredResourcesSubnetworks", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PacketMirroringMirroredResourcesSubnetworks)
	if !ok {
		actualNotPointer, ok := a.(PacketMirroringMirroredResourcesSubnetworks)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringMirroredResourcesSubnetworks", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Url, actual.Url, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Url")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CanonicalUrl, actual.CanonicalUrl, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CanonicalUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePacketMirroringMirroredResourcesInstancesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PacketMirroringMirroredResourcesInstances)
	if !ok {
		desiredNotPointer, ok := d.(PacketMirroringMirroredResourcesInstances)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringMirroredResourcesInstances or *PacketMirroringMirroredResourcesInstances", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PacketMirroringMirroredResourcesInstances)
	if !ok {
		actualNotPointer, ok := a.(PacketMirroringMirroredResourcesInstances)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringMirroredResourcesInstances", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Url, actual.Url, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Url")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CanonicalUrl, actual.CanonicalUrl, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CanonicalUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func comparePacketMirroringFilterNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*PacketMirroringFilter)
	if !ok {
		desiredNotPointer, ok := d.(PacketMirroringFilter)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringFilter or *PacketMirroringFilter", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*PacketMirroringFilter)
	if !ok {
		actualNotPointer, ok := a.(PacketMirroringFilter)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a PacketMirroringFilter", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CidrRanges, actual.CidrRanges, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("CidrRanges")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPProtocols, actual.IPProtocols, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("IPProtocols")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Direction, actual.Direction, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updatePacketMirroringPatchOperation")}, fn.AddNest("Direction")); len(ds) != 0 || err != nil {
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
func (r *PacketMirroring) urlNormalized() *PacketMirroring {
	normalized := dcl.Copy(*r).(PacketMirroring)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Region = dcl.SelfLinkToName(r.Region)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *PacketMirroring) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "Patch" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/regions/{{location}}/packetMirrorings/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the PacketMirroring resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *PacketMirroring) marshal(c *Client) ([]byte, error) {
	m, err := expandPacketMirroring(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling PacketMirroring: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalPacketMirroring decodes JSON responses into the PacketMirroring resource schema.
func unmarshalPacketMirroring(b []byte, c *Client, res *PacketMirroring) (*PacketMirroring, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapPacketMirroring(m, c, res)
}

func unmarshalMapPacketMirroring(m map[string]interface{}, c *Client, res *PacketMirroring) (*PacketMirroring, error) {

	flattened := flattenPacketMirroring(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandPacketMirroring expands PacketMirroring into a JSON request object.
func expandPacketMirroring(c *Client, f *PacketMirroring) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := expandPacketMirroringNetwork(c, f.Network, res); err != nil {
		return nil, fmt.Errorf("error expanding Network into network: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["network"] = v
	}
	if v := f.Priority; dcl.ValueShouldBeSent(v) {
		m["priority"] = v
	}
	if v, err := expandPacketMirroringCollectorIlb(c, f.CollectorIlb, res); err != nil {
		return nil, fmt.Errorf("error expanding CollectorIlb into collectorIlb: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["collectorIlb"] = v
	}
	if v, err := expandPacketMirroringMirroredResources(c, f.MirroredResources, res); err != nil {
		return nil, fmt.Errorf("error expanding MirroredResources into mirroredResources: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["mirroredResources"] = v
	}
	if v, err := expandPacketMirroringFilter(c, f.Filter, res); err != nil {
		return nil, fmt.Errorf("error expanding Filter into filter: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["filter"] = v
	}
	if v := f.Enable; dcl.ValueShouldBeSent(v) {
		m["enable"] = v
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

// flattenPacketMirroring flattens PacketMirroring from a JSON request object into the
// PacketMirroring type.
func flattenPacketMirroring(c *Client, i interface{}, res *PacketMirroring) *PacketMirroring {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &PacketMirroring{}
	resultRes.Id = dcl.FlattenInteger(m["id"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Region = dcl.FlattenString(m["region"])
	resultRes.Network = flattenPacketMirroringNetwork(c, m["network"], res)
	resultRes.Priority = dcl.FlattenInteger(m["priority"])
	resultRes.CollectorIlb = flattenPacketMirroringCollectorIlb(c, m["collectorIlb"], res)
	resultRes.MirroredResources = flattenPacketMirroringMirroredResources(c, m["mirroredResources"], res)
	resultRes.Filter = flattenPacketMirroringFilter(c, m["filter"], res)
	resultRes.Enable = flattenPacketMirroringEnableEnum(m["enable"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandPacketMirroringNetworkMap expands the contents of PacketMirroringNetwork into a JSON
// request object.
func expandPacketMirroringNetworkMap(c *Client, f map[string]PacketMirroringNetwork, res *PacketMirroring) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPacketMirroringNetwork(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPacketMirroringNetworkSlice expands the contents of PacketMirroringNetwork into a JSON
// request object.
func expandPacketMirroringNetworkSlice(c *Client, f []PacketMirroringNetwork, res *PacketMirroring) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPacketMirroringNetwork(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPacketMirroringNetworkMap flattens the contents of PacketMirroringNetwork from a JSON
// response object.
func flattenPacketMirroringNetworkMap(c *Client, i interface{}, res *PacketMirroring) map[string]PacketMirroringNetwork {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PacketMirroringNetwork{}
	}

	if len(a) == 0 {
		return map[string]PacketMirroringNetwork{}
	}

	items := make(map[string]PacketMirroringNetwork)
	for k, item := range a {
		items[k] = *flattenPacketMirroringNetwork(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPacketMirroringNetworkSlice flattens the contents of PacketMirroringNetwork from a JSON
// response object.
func flattenPacketMirroringNetworkSlice(c *Client, i interface{}, res *PacketMirroring) []PacketMirroringNetwork {
	a, ok := i.([]interface{})
	if !ok {
		return []PacketMirroringNetwork{}
	}

	if len(a) == 0 {
		return []PacketMirroringNetwork{}
	}

	items := make([]PacketMirroringNetwork, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPacketMirroringNetwork(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPacketMirroringNetwork expands an instance of PacketMirroringNetwork into a JSON
// request object.
func expandPacketMirroringNetwork(c *Client, f *PacketMirroringNetwork, res *PacketMirroring) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Url; !dcl.IsEmptyValueIndirect(v) {
		m["url"] = v
	}

	return m, nil
}

// flattenPacketMirroringNetwork flattens an instance of PacketMirroringNetwork from a JSON
// response object.
func flattenPacketMirroringNetwork(c *Client, i interface{}, res *PacketMirroring) *PacketMirroringNetwork {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PacketMirroringNetwork{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPacketMirroringNetwork
	}
	r.Url = dcl.FlattenString(m["url"])
	r.CanonicalUrl = dcl.FlattenString(m["canonicalUrl"])

	return r
}

// expandPacketMirroringCollectorIlbMap expands the contents of PacketMirroringCollectorIlb into a JSON
// request object.
func expandPacketMirroringCollectorIlbMap(c *Client, f map[string]PacketMirroringCollectorIlb, res *PacketMirroring) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPacketMirroringCollectorIlb(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPacketMirroringCollectorIlbSlice expands the contents of PacketMirroringCollectorIlb into a JSON
// request object.
func expandPacketMirroringCollectorIlbSlice(c *Client, f []PacketMirroringCollectorIlb, res *PacketMirroring) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPacketMirroringCollectorIlb(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPacketMirroringCollectorIlbMap flattens the contents of PacketMirroringCollectorIlb from a JSON
// response object.
func flattenPacketMirroringCollectorIlbMap(c *Client, i interface{}, res *PacketMirroring) map[string]PacketMirroringCollectorIlb {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PacketMirroringCollectorIlb{}
	}

	if len(a) == 0 {
		return map[string]PacketMirroringCollectorIlb{}
	}

	items := make(map[string]PacketMirroringCollectorIlb)
	for k, item := range a {
		items[k] = *flattenPacketMirroringCollectorIlb(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPacketMirroringCollectorIlbSlice flattens the contents of PacketMirroringCollectorIlb from a JSON
// response object.
func flattenPacketMirroringCollectorIlbSlice(c *Client, i interface{}, res *PacketMirroring) []PacketMirroringCollectorIlb {
	a, ok := i.([]interface{})
	if !ok {
		return []PacketMirroringCollectorIlb{}
	}

	if len(a) == 0 {
		return []PacketMirroringCollectorIlb{}
	}

	items := make([]PacketMirroringCollectorIlb, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPacketMirroringCollectorIlb(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPacketMirroringCollectorIlb expands an instance of PacketMirroringCollectorIlb into a JSON
// request object.
func expandPacketMirroringCollectorIlb(c *Client, f *PacketMirroringCollectorIlb, res *PacketMirroring) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Url; !dcl.IsEmptyValueIndirect(v) {
		m["url"] = v
	}

	return m, nil
}

// flattenPacketMirroringCollectorIlb flattens an instance of PacketMirroringCollectorIlb from a JSON
// response object.
func flattenPacketMirroringCollectorIlb(c *Client, i interface{}, res *PacketMirroring) *PacketMirroringCollectorIlb {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PacketMirroringCollectorIlb{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPacketMirroringCollectorIlb
	}
	r.Url = dcl.FlattenString(m["url"])
	r.CanonicalUrl = dcl.FlattenString(m["canonicalUrl"])

	return r
}

// expandPacketMirroringMirroredResourcesMap expands the contents of PacketMirroringMirroredResources into a JSON
// request object.
func expandPacketMirroringMirroredResourcesMap(c *Client, f map[string]PacketMirroringMirroredResources, res *PacketMirroring) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPacketMirroringMirroredResources(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPacketMirroringMirroredResourcesSlice expands the contents of PacketMirroringMirroredResources into a JSON
// request object.
func expandPacketMirroringMirroredResourcesSlice(c *Client, f []PacketMirroringMirroredResources, res *PacketMirroring) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPacketMirroringMirroredResources(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPacketMirroringMirroredResourcesMap flattens the contents of PacketMirroringMirroredResources from a JSON
// response object.
func flattenPacketMirroringMirroredResourcesMap(c *Client, i interface{}, res *PacketMirroring) map[string]PacketMirroringMirroredResources {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PacketMirroringMirroredResources{}
	}

	if len(a) == 0 {
		return map[string]PacketMirroringMirroredResources{}
	}

	items := make(map[string]PacketMirroringMirroredResources)
	for k, item := range a {
		items[k] = *flattenPacketMirroringMirroredResources(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPacketMirroringMirroredResourcesSlice flattens the contents of PacketMirroringMirroredResources from a JSON
// response object.
func flattenPacketMirroringMirroredResourcesSlice(c *Client, i interface{}, res *PacketMirroring) []PacketMirroringMirroredResources {
	a, ok := i.([]interface{})
	if !ok {
		return []PacketMirroringMirroredResources{}
	}

	if len(a) == 0 {
		return []PacketMirroringMirroredResources{}
	}

	items := make([]PacketMirroringMirroredResources, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPacketMirroringMirroredResources(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPacketMirroringMirroredResources expands an instance of PacketMirroringMirroredResources into a JSON
// request object.
func expandPacketMirroringMirroredResources(c *Client, f *PacketMirroringMirroredResources, res *PacketMirroring) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandPacketMirroringMirroredResourcesSubnetworksSlice(c, f.Subnetworks, res); err != nil {
		return nil, fmt.Errorf("error expanding Subnetworks into subnetworks: %w", err)
	} else if v != nil {
		m["subnetworks"] = v
	}
	if v, err := expandPacketMirroringMirroredResourcesInstancesSlice(c, f.Instances, res); err != nil {
		return nil, fmt.Errorf("error expanding Instances into instances: %w", err)
	} else if v != nil {
		m["instances"] = v
	}
	if v := f.Tags; v != nil {
		m["tags"] = v
	}

	return m, nil
}

// flattenPacketMirroringMirroredResources flattens an instance of PacketMirroringMirroredResources from a JSON
// response object.
func flattenPacketMirroringMirroredResources(c *Client, i interface{}, res *PacketMirroring) *PacketMirroringMirroredResources {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PacketMirroringMirroredResources{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPacketMirroringMirroredResources
	}
	r.Subnetworks = flattenPacketMirroringMirroredResourcesSubnetworksSlice(c, m["subnetworks"], res)
	r.Instances = flattenPacketMirroringMirroredResourcesInstancesSlice(c, m["instances"], res)
	r.Tags = dcl.FlattenStringSlice(m["tags"])

	return r
}

// expandPacketMirroringMirroredResourcesSubnetworksMap expands the contents of PacketMirroringMirroredResourcesSubnetworks into a JSON
// request object.
func expandPacketMirroringMirroredResourcesSubnetworksMap(c *Client, f map[string]PacketMirroringMirroredResourcesSubnetworks, res *PacketMirroring) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPacketMirroringMirroredResourcesSubnetworks(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPacketMirroringMirroredResourcesSubnetworksSlice expands the contents of PacketMirroringMirroredResourcesSubnetworks into a JSON
// request object.
func expandPacketMirroringMirroredResourcesSubnetworksSlice(c *Client, f []PacketMirroringMirroredResourcesSubnetworks, res *PacketMirroring) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPacketMirroringMirroredResourcesSubnetworks(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPacketMirroringMirroredResourcesSubnetworksMap flattens the contents of PacketMirroringMirroredResourcesSubnetworks from a JSON
// response object.
func flattenPacketMirroringMirroredResourcesSubnetworksMap(c *Client, i interface{}, res *PacketMirroring) map[string]PacketMirroringMirroredResourcesSubnetworks {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PacketMirroringMirroredResourcesSubnetworks{}
	}

	if len(a) == 0 {
		return map[string]PacketMirroringMirroredResourcesSubnetworks{}
	}

	items := make(map[string]PacketMirroringMirroredResourcesSubnetworks)
	for k, item := range a {
		items[k] = *flattenPacketMirroringMirroredResourcesSubnetworks(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPacketMirroringMirroredResourcesSubnetworksSlice flattens the contents of PacketMirroringMirroredResourcesSubnetworks from a JSON
// response object.
func flattenPacketMirroringMirroredResourcesSubnetworksSlice(c *Client, i interface{}, res *PacketMirroring) []PacketMirroringMirroredResourcesSubnetworks {
	a, ok := i.([]interface{})
	if !ok {
		return []PacketMirroringMirroredResourcesSubnetworks{}
	}

	if len(a) == 0 {
		return []PacketMirroringMirroredResourcesSubnetworks{}
	}

	items := make([]PacketMirroringMirroredResourcesSubnetworks, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPacketMirroringMirroredResourcesSubnetworks(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPacketMirroringMirroredResourcesSubnetworks expands an instance of PacketMirroringMirroredResourcesSubnetworks into a JSON
// request object.
func expandPacketMirroringMirroredResourcesSubnetworks(c *Client, f *PacketMirroringMirroredResourcesSubnetworks, res *PacketMirroring) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Url; !dcl.IsEmptyValueIndirect(v) {
		m["url"] = v
	}

	return m, nil
}

// flattenPacketMirroringMirroredResourcesSubnetworks flattens an instance of PacketMirroringMirroredResourcesSubnetworks from a JSON
// response object.
func flattenPacketMirroringMirroredResourcesSubnetworks(c *Client, i interface{}, res *PacketMirroring) *PacketMirroringMirroredResourcesSubnetworks {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PacketMirroringMirroredResourcesSubnetworks{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPacketMirroringMirroredResourcesSubnetworks
	}
	r.Url = dcl.FlattenString(m["url"])
	r.CanonicalUrl = dcl.FlattenString(m["canonicalUrl"])

	return r
}

// expandPacketMirroringMirroredResourcesInstancesMap expands the contents of PacketMirroringMirroredResourcesInstances into a JSON
// request object.
func expandPacketMirroringMirroredResourcesInstancesMap(c *Client, f map[string]PacketMirroringMirroredResourcesInstances, res *PacketMirroring) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPacketMirroringMirroredResourcesInstances(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPacketMirroringMirroredResourcesInstancesSlice expands the contents of PacketMirroringMirroredResourcesInstances into a JSON
// request object.
func expandPacketMirroringMirroredResourcesInstancesSlice(c *Client, f []PacketMirroringMirroredResourcesInstances, res *PacketMirroring) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPacketMirroringMirroredResourcesInstances(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPacketMirroringMirroredResourcesInstancesMap flattens the contents of PacketMirroringMirroredResourcesInstances from a JSON
// response object.
func flattenPacketMirroringMirroredResourcesInstancesMap(c *Client, i interface{}, res *PacketMirroring) map[string]PacketMirroringMirroredResourcesInstances {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PacketMirroringMirroredResourcesInstances{}
	}

	if len(a) == 0 {
		return map[string]PacketMirroringMirroredResourcesInstances{}
	}

	items := make(map[string]PacketMirroringMirroredResourcesInstances)
	for k, item := range a {
		items[k] = *flattenPacketMirroringMirroredResourcesInstances(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPacketMirroringMirroredResourcesInstancesSlice flattens the contents of PacketMirroringMirroredResourcesInstances from a JSON
// response object.
func flattenPacketMirroringMirroredResourcesInstancesSlice(c *Client, i interface{}, res *PacketMirroring) []PacketMirroringMirroredResourcesInstances {
	a, ok := i.([]interface{})
	if !ok {
		return []PacketMirroringMirroredResourcesInstances{}
	}

	if len(a) == 0 {
		return []PacketMirroringMirroredResourcesInstances{}
	}

	items := make([]PacketMirroringMirroredResourcesInstances, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPacketMirroringMirroredResourcesInstances(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPacketMirroringMirroredResourcesInstances expands an instance of PacketMirroringMirroredResourcesInstances into a JSON
// request object.
func expandPacketMirroringMirroredResourcesInstances(c *Client, f *PacketMirroringMirroredResourcesInstances, res *PacketMirroring) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Url; !dcl.IsEmptyValueIndirect(v) {
		m["url"] = v
	}

	return m, nil
}

// flattenPacketMirroringMirroredResourcesInstances flattens an instance of PacketMirroringMirroredResourcesInstances from a JSON
// response object.
func flattenPacketMirroringMirroredResourcesInstances(c *Client, i interface{}, res *PacketMirroring) *PacketMirroringMirroredResourcesInstances {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PacketMirroringMirroredResourcesInstances{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPacketMirroringMirroredResourcesInstances
	}
	r.Url = dcl.FlattenString(m["url"])
	r.CanonicalUrl = dcl.FlattenString(m["canonicalUrl"])

	return r
}

// expandPacketMirroringFilterMap expands the contents of PacketMirroringFilter into a JSON
// request object.
func expandPacketMirroringFilterMap(c *Client, f map[string]PacketMirroringFilter, res *PacketMirroring) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandPacketMirroringFilter(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandPacketMirroringFilterSlice expands the contents of PacketMirroringFilter into a JSON
// request object.
func expandPacketMirroringFilterSlice(c *Client, f []PacketMirroringFilter, res *PacketMirroring) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandPacketMirroringFilter(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenPacketMirroringFilterMap flattens the contents of PacketMirroringFilter from a JSON
// response object.
func flattenPacketMirroringFilterMap(c *Client, i interface{}, res *PacketMirroring) map[string]PacketMirroringFilter {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PacketMirroringFilter{}
	}

	if len(a) == 0 {
		return map[string]PacketMirroringFilter{}
	}

	items := make(map[string]PacketMirroringFilter)
	for k, item := range a {
		items[k] = *flattenPacketMirroringFilter(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenPacketMirroringFilterSlice flattens the contents of PacketMirroringFilter from a JSON
// response object.
func flattenPacketMirroringFilterSlice(c *Client, i interface{}, res *PacketMirroring) []PacketMirroringFilter {
	a, ok := i.([]interface{})
	if !ok {
		return []PacketMirroringFilter{}
	}

	if len(a) == 0 {
		return []PacketMirroringFilter{}
	}

	items := make([]PacketMirroringFilter, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPacketMirroringFilter(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandPacketMirroringFilter expands an instance of PacketMirroringFilter into a JSON
// request object.
func expandPacketMirroringFilter(c *Client, f *PacketMirroringFilter, res *PacketMirroring) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.CidrRanges; v != nil {
		m["cidrRanges"] = v
	}
	if v := f.IPProtocols; v != nil {
		m["IPProtocols"] = v
	}
	if v := f.Direction; !dcl.IsEmptyValueIndirect(v) {
		m["direction"] = v
	}

	return m, nil
}

// flattenPacketMirroringFilter flattens an instance of PacketMirroringFilter from a JSON
// response object.
func flattenPacketMirroringFilter(c *Client, i interface{}, res *PacketMirroring) *PacketMirroringFilter {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &PacketMirroringFilter{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyPacketMirroringFilter
	}
	r.CidrRanges = dcl.FlattenStringSlice(m["cidrRanges"])
	r.IPProtocols = dcl.FlattenStringSlice(m["IPProtocols"])
	r.Direction = flattenPacketMirroringFilterDirectionEnum(m["direction"])

	return r
}

// flattenPacketMirroringFilterDirectionEnumMap flattens the contents of PacketMirroringFilterDirectionEnum from a JSON
// response object.
func flattenPacketMirroringFilterDirectionEnumMap(c *Client, i interface{}, res *PacketMirroring) map[string]PacketMirroringFilterDirectionEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PacketMirroringFilterDirectionEnum{}
	}

	if len(a) == 0 {
		return map[string]PacketMirroringFilterDirectionEnum{}
	}

	items := make(map[string]PacketMirroringFilterDirectionEnum)
	for k, item := range a {
		items[k] = *flattenPacketMirroringFilterDirectionEnum(item.(interface{}))
	}

	return items
}

// flattenPacketMirroringFilterDirectionEnumSlice flattens the contents of PacketMirroringFilterDirectionEnum from a JSON
// response object.
func flattenPacketMirroringFilterDirectionEnumSlice(c *Client, i interface{}, res *PacketMirroring) []PacketMirroringFilterDirectionEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PacketMirroringFilterDirectionEnum{}
	}

	if len(a) == 0 {
		return []PacketMirroringFilterDirectionEnum{}
	}

	items := make([]PacketMirroringFilterDirectionEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPacketMirroringFilterDirectionEnum(item.(interface{})))
	}

	return items
}

// flattenPacketMirroringFilterDirectionEnum asserts that an interface is a string, and returns a
// pointer to a *PacketMirroringFilterDirectionEnum with the same value as that string.
func flattenPacketMirroringFilterDirectionEnum(i interface{}) *PacketMirroringFilterDirectionEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return PacketMirroringFilterDirectionEnumRef(s)
}

// flattenPacketMirroringEnableEnumMap flattens the contents of PacketMirroringEnableEnum from a JSON
// response object.
func flattenPacketMirroringEnableEnumMap(c *Client, i interface{}, res *PacketMirroring) map[string]PacketMirroringEnableEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]PacketMirroringEnableEnum{}
	}

	if len(a) == 0 {
		return map[string]PacketMirroringEnableEnum{}
	}

	items := make(map[string]PacketMirroringEnableEnum)
	for k, item := range a {
		items[k] = *flattenPacketMirroringEnableEnum(item.(interface{}))
	}

	return items
}

// flattenPacketMirroringEnableEnumSlice flattens the contents of PacketMirroringEnableEnum from a JSON
// response object.
func flattenPacketMirroringEnableEnumSlice(c *Client, i interface{}, res *PacketMirroring) []PacketMirroringEnableEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []PacketMirroringEnableEnum{}
	}

	if len(a) == 0 {
		return []PacketMirroringEnableEnum{}
	}

	items := make([]PacketMirroringEnableEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenPacketMirroringEnableEnum(item.(interface{})))
	}

	return items
}

// flattenPacketMirroringEnableEnum asserts that an interface is a string, and returns a
// pointer to a *PacketMirroringEnableEnum with the same value as that string.
func flattenPacketMirroringEnableEnum(i interface{}) *PacketMirroringEnableEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return PacketMirroringEnableEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *PacketMirroring) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalPacketMirroring(b, c, r)
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

type packetMirroringDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         packetMirroringApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToPacketMirroringDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]packetMirroringDiff, error) {
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
	var diffs []packetMirroringDiff
	// For each operation name, create a packetMirroringDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := packetMirroringDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToPacketMirroringApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToPacketMirroringApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (packetMirroringApiOperation, error) {
	switch opName {

	case "updatePacketMirroringPatchOperation":
		return &updatePacketMirroringPatchOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractPacketMirroringFields(r *PacketMirroring) error {
	vNetwork := r.Network
	if vNetwork == nil {
		// note: explicitly not the empty object.
		vNetwork = &PacketMirroringNetwork{}
	}
	if err := extractPacketMirroringNetworkFields(r, vNetwork); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetwork) {
		r.Network = vNetwork
	}
	vCollectorIlb := r.CollectorIlb
	if vCollectorIlb == nil {
		// note: explicitly not the empty object.
		vCollectorIlb = &PacketMirroringCollectorIlb{}
	}
	if err := extractPacketMirroringCollectorIlbFields(r, vCollectorIlb); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCollectorIlb) {
		r.CollectorIlb = vCollectorIlb
	}
	vMirroredResources := r.MirroredResources
	if vMirroredResources == nil {
		// note: explicitly not the empty object.
		vMirroredResources = &PacketMirroringMirroredResources{}
	}
	if err := extractPacketMirroringMirroredResourcesFields(r, vMirroredResources); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMirroredResources) {
		r.MirroredResources = vMirroredResources
	}
	vFilter := r.Filter
	if vFilter == nil {
		// note: explicitly not the empty object.
		vFilter = &PacketMirroringFilter{}
	}
	if err := extractPacketMirroringFilterFields(r, vFilter); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFilter) {
		r.Filter = vFilter
	}
	return nil
}
func extractPacketMirroringNetworkFields(r *PacketMirroring, o *PacketMirroringNetwork) error {
	return nil
}
func extractPacketMirroringCollectorIlbFields(r *PacketMirroring, o *PacketMirroringCollectorIlb) error {
	return nil
}
func extractPacketMirroringMirroredResourcesFields(r *PacketMirroring, o *PacketMirroringMirroredResources) error {
	return nil
}
func extractPacketMirroringMirroredResourcesSubnetworksFields(r *PacketMirroring, o *PacketMirroringMirroredResourcesSubnetworks) error {
	return nil
}
func extractPacketMirroringMirroredResourcesInstancesFields(r *PacketMirroring, o *PacketMirroringMirroredResourcesInstances) error {
	return nil
}
func extractPacketMirroringFilterFields(r *PacketMirroring, o *PacketMirroringFilter) error {
	return nil
}

func postReadExtractPacketMirroringFields(r *PacketMirroring) error {
	vNetwork := r.Network
	if vNetwork == nil {
		// note: explicitly not the empty object.
		vNetwork = &PacketMirroringNetwork{}
	}
	if err := postReadExtractPacketMirroringNetworkFields(r, vNetwork); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetwork) {
		r.Network = vNetwork
	}
	vCollectorIlb := r.CollectorIlb
	if vCollectorIlb == nil {
		// note: explicitly not the empty object.
		vCollectorIlb = &PacketMirroringCollectorIlb{}
	}
	if err := postReadExtractPacketMirroringCollectorIlbFields(r, vCollectorIlb); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCollectorIlb) {
		r.CollectorIlb = vCollectorIlb
	}
	vMirroredResources := r.MirroredResources
	if vMirroredResources == nil {
		// note: explicitly not the empty object.
		vMirroredResources = &PacketMirroringMirroredResources{}
	}
	if err := postReadExtractPacketMirroringMirroredResourcesFields(r, vMirroredResources); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMirroredResources) {
		r.MirroredResources = vMirroredResources
	}
	vFilter := r.Filter
	if vFilter == nil {
		// note: explicitly not the empty object.
		vFilter = &PacketMirroringFilter{}
	}
	if err := postReadExtractPacketMirroringFilterFields(r, vFilter); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFilter) {
		r.Filter = vFilter
	}
	return nil
}
func postReadExtractPacketMirroringNetworkFields(r *PacketMirroring, o *PacketMirroringNetwork) error {
	return nil
}
func postReadExtractPacketMirroringCollectorIlbFields(r *PacketMirroring, o *PacketMirroringCollectorIlb) error {
	return nil
}
func postReadExtractPacketMirroringMirroredResourcesFields(r *PacketMirroring, o *PacketMirroringMirroredResources) error {
	return nil
}
func postReadExtractPacketMirroringMirroredResourcesSubnetworksFields(r *PacketMirroring, o *PacketMirroringMirroredResourcesSubnetworks) error {
	return nil
}
func postReadExtractPacketMirroringMirroredResourcesInstancesFields(r *PacketMirroring, o *PacketMirroringMirroredResourcesInstances) error {
	return nil
}
func postReadExtractPacketMirroringFilterFields(r *PacketMirroring, o *PacketMirroringFilter) error {
	return nil
}
