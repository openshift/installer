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
package networkconnectivity

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

func (r *Spoke) validate() error {

	if err := dcl.ValidateExactlyOneOfFieldsSet([]string{"LinkedVpnTunnels", "LinkedInterconnectAttachments", "LinkedRouterApplianceInstances", "LinkedVPCNetwork"}, r.LinkedVpnTunnels, r.LinkedInterconnectAttachments, r.LinkedRouterApplianceInstances, r.LinkedVPCNetwork); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "hub"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.LinkedVpnTunnels) {
		if err := r.LinkedVpnTunnels.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.LinkedInterconnectAttachments) {
		if err := r.LinkedInterconnectAttachments.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.LinkedRouterApplianceInstances) {
		if err := r.LinkedRouterApplianceInstances.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.LinkedVPCNetwork) {
		if err := r.LinkedVPCNetwork.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *SpokeLinkedVpnTunnels) validate() error {
	if err := dcl.Required(r, "uris"); err != nil {
		return err
	}
	if err := dcl.Required(r, "siteToSiteDataTransfer"); err != nil {
		return err
	}
	return nil
}
func (r *SpokeLinkedInterconnectAttachments) validate() error {
	if err := dcl.Required(r, "uris"); err != nil {
		return err
	}
	if err := dcl.Required(r, "siteToSiteDataTransfer"); err != nil {
		return err
	}
	return nil
}
func (r *SpokeLinkedRouterApplianceInstances) validate() error {
	if err := dcl.Required(r, "instances"); err != nil {
		return err
	}
	if err := dcl.Required(r, "siteToSiteDataTransfer"); err != nil {
		return err
	}
	return nil
}
func (r *SpokeLinkedRouterApplianceInstancesInstances) validate() error {
	return nil
}
func (r *SpokeLinkedVPCNetwork) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *Spoke) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://networkconnectivity.googleapis.com/v1/", params)
}

func (r *Spoke) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/spokes/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Spoke) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/spokes", nr.basePath(), userBasePath, params), nil

}

func (r *Spoke) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/spokes?spokeId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Spoke) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/spokes/{{name}}", nr.basePath(), userBasePath, params), nil
}

// spokeApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type spokeApiOperation interface {
	do(context.Context, *Spoke, *Client) error
}

// newUpdateSpokeUpdateSpokeRequest creates a request for an
// Spoke resource's UpdateSpoke update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateSpokeUpdateSpokeRequest(ctx context.Context, f *Spoke, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	return req, nil
}

// marshalUpdateSpokeUpdateSpokeRequest converts the update into
// the final JSON request body.
func marshalUpdateSpokeUpdateSpokeRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateSpokeUpdateSpokeOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateSpokeUpdateSpokeOperation) do(ctx context.Context, r *Spoke, c *Client) error {
	_, err := c.GetSpoke(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateSpoke")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateSpokeUpdateSpokeRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateSpokeUpdateSpokeRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	err = o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET")

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listSpokeRaw(ctx context.Context, r *Spoke, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != SpokeMaxPage {
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

type listSpokeOperation struct {
	Spokes []map[string]interface{} `json:"spokes"`
	Token  string                   `json:"nextPageToken"`
}

func (c *Client) listSpoke(ctx context.Context, r *Spoke, pageToken string, pageSize int32) ([]*Spoke, string, error) {
	b, err := c.listSpokeRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listSpokeOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Spoke
	for _, v := range m.Spokes {
		res, err := unmarshalMapSpoke(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllSpoke(ctx context.Context, f func(*Spoke) bool, resources []*Spoke) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteSpoke(ctx, res)
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

type deleteSpokeOperation struct{}

func (op *deleteSpokeOperation) do(ctx context.Context, r *Spoke, c *Client) error {
	r, err := c.GetSpoke(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Spoke not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetSpoke checking for existence. error: %v", err)
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
	var o operations.StandardGCPOperation
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
		_, err := c.GetSpoke(ctx, r)
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
type createSpokeOperation struct {
	response map[string]interface{}
}

func (op *createSpokeOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createSpokeOperation) do(ctx context.Context, r *Spoke, c *Client) error {
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
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(context.WithValue(ctx, dcl.DoNotLogRequestsKey, true), c.Config, r.basePath(), "GET"); err != nil {
		c.Config.Logger.Warningf("Creation failed after waiting for operation: %v", err)
		return err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Successfully waited for operation")
	op.response, _ = o.FirstResponse()

	if _, err := c.GetSpoke(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getSpokeRaw(ctx context.Context, r *Spoke) ([]byte, error) {

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

func (c *Client) spokeDiffsForRawDesired(ctx context.Context, rawDesired *Spoke, opts ...dcl.ApplyOption) (initial, desired *Spoke, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Spoke
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Spoke); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Spoke, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetSpoke(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Spoke resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Spoke resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Spoke resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeSpokeDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Spoke: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Spoke: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractSpokeFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeSpokeInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Spoke: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeSpokeDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Spoke: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffSpoke(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeSpokeInitialState(rawInitial, rawDesired *Spoke) (*Spoke, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.

	if !dcl.IsZeroValue(rawInitial.LinkedVpnTunnels) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.LinkedInterconnectAttachments, rawInitial.LinkedRouterApplianceInstances, rawInitial.LinkedVPCNetwork) {
			rawInitial.LinkedVpnTunnels = EmptySpokeLinkedVpnTunnels
		}
	}

	if !dcl.IsZeroValue(rawInitial.LinkedInterconnectAttachments) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.LinkedVpnTunnels, rawInitial.LinkedRouterApplianceInstances, rawInitial.LinkedVPCNetwork) {
			rawInitial.LinkedInterconnectAttachments = EmptySpokeLinkedInterconnectAttachments
		}
	}

	if !dcl.IsZeroValue(rawInitial.LinkedRouterApplianceInstances) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.LinkedVpnTunnels, rawInitial.LinkedInterconnectAttachments, rawInitial.LinkedVPCNetwork) {
			rawInitial.LinkedRouterApplianceInstances = EmptySpokeLinkedRouterApplianceInstances
		}
	}

	if !dcl.IsZeroValue(rawInitial.LinkedVPCNetwork) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.LinkedVpnTunnels, rawInitial.LinkedInterconnectAttachments, rawInitial.LinkedRouterApplianceInstances) {
			rawInitial.LinkedVPCNetwork = EmptySpokeLinkedVPCNetwork
		}
	}

	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeSpokeDesiredState(rawDesired, rawInitial *Spoke, opts ...dcl.ApplyOption) (*Spoke, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.LinkedVpnTunnels = canonicalizeSpokeLinkedVpnTunnels(rawDesired.LinkedVpnTunnels, nil, opts...)
		rawDesired.LinkedInterconnectAttachments = canonicalizeSpokeLinkedInterconnectAttachments(rawDesired.LinkedInterconnectAttachments, nil, opts...)
		rawDesired.LinkedRouterApplianceInstances = canonicalizeSpokeLinkedRouterApplianceInstances(rawDesired.LinkedRouterApplianceInstances, nil, opts...)
		rawDesired.LinkedVPCNetwork = canonicalizeSpokeLinkedVPCNetwork(rawDesired.LinkedVPCNetwork, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Spoke{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.IsZeroValue(rawDesired.Hub) || (dcl.IsEmptyValueIndirect(rawDesired.Hub) && dcl.IsEmptyValueIndirect(rawInitial.Hub)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Hub = rawInitial.Hub
	} else {
		canonicalDesired.Hub = rawDesired.Hub
	}
	canonicalDesired.LinkedVpnTunnels = canonicalizeSpokeLinkedVpnTunnels(rawDesired.LinkedVpnTunnels, rawInitial.LinkedVpnTunnels, opts...)
	canonicalDesired.LinkedInterconnectAttachments = canonicalizeSpokeLinkedInterconnectAttachments(rawDesired.LinkedInterconnectAttachments, rawInitial.LinkedInterconnectAttachments, opts...)
	canonicalDesired.LinkedRouterApplianceInstances = canonicalizeSpokeLinkedRouterApplianceInstances(rawDesired.LinkedRouterApplianceInstances, rawInitial.LinkedRouterApplianceInstances, opts...)
	canonicalDesired.LinkedVPCNetwork = canonicalizeSpokeLinkedVPCNetwork(rawDesired.LinkedVPCNetwork, rawInitial.LinkedVPCNetwork, opts...)
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

	if canonicalDesired.LinkedVpnTunnels != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.LinkedInterconnectAttachments, rawDesired.LinkedRouterApplianceInstances, rawDesired.LinkedVPCNetwork) {
			canonicalDesired.LinkedVpnTunnels = EmptySpokeLinkedVpnTunnels
		}
	}

	if canonicalDesired.LinkedInterconnectAttachments != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.LinkedVpnTunnels, rawDesired.LinkedRouterApplianceInstances, rawDesired.LinkedVPCNetwork) {
			canonicalDesired.LinkedInterconnectAttachments = EmptySpokeLinkedInterconnectAttachments
		}
	}

	if canonicalDesired.LinkedRouterApplianceInstances != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.LinkedVpnTunnels, rawDesired.LinkedInterconnectAttachments, rawDesired.LinkedVPCNetwork) {
			canonicalDesired.LinkedRouterApplianceInstances = EmptySpokeLinkedRouterApplianceInstances
		}
	}

	if canonicalDesired.LinkedVPCNetwork != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.LinkedVpnTunnels, rawDesired.LinkedInterconnectAttachments, rawDesired.LinkedRouterApplianceInstances) {
			canonicalDesired.LinkedVPCNetwork = EmptySpokeLinkedVPCNetwork
		}
	}

	return canonicalDesired, nil
}

func canonicalizeSpokeNewState(c *Client, rawNew, rawDesired *Spoke) (*Spoke, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
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

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Hub) && dcl.IsEmptyValueIndirect(rawDesired.Hub) {
		rawNew.Hub = rawDesired.Hub
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.LinkedVpnTunnels) && dcl.IsEmptyValueIndirect(rawDesired.LinkedVpnTunnels) {
		rawNew.LinkedVpnTunnels = rawDesired.LinkedVpnTunnels
	} else {
		rawNew.LinkedVpnTunnels = canonicalizeNewSpokeLinkedVpnTunnels(c, rawDesired.LinkedVpnTunnels, rawNew.LinkedVpnTunnels)
	}

	if dcl.IsEmptyValueIndirect(rawNew.LinkedInterconnectAttachments) && dcl.IsEmptyValueIndirect(rawDesired.LinkedInterconnectAttachments) {
		rawNew.LinkedInterconnectAttachments = rawDesired.LinkedInterconnectAttachments
	} else {
		rawNew.LinkedInterconnectAttachments = canonicalizeNewSpokeLinkedInterconnectAttachments(c, rawDesired.LinkedInterconnectAttachments, rawNew.LinkedInterconnectAttachments)
	}

	if dcl.IsEmptyValueIndirect(rawNew.LinkedRouterApplianceInstances) && dcl.IsEmptyValueIndirect(rawDesired.LinkedRouterApplianceInstances) {
		rawNew.LinkedRouterApplianceInstances = rawDesired.LinkedRouterApplianceInstances
	} else {
		rawNew.LinkedRouterApplianceInstances = canonicalizeNewSpokeLinkedRouterApplianceInstances(c, rawDesired.LinkedRouterApplianceInstances, rawNew.LinkedRouterApplianceInstances)
	}

	if dcl.IsEmptyValueIndirect(rawNew.LinkedVPCNetwork) && dcl.IsEmptyValueIndirect(rawDesired.LinkedVPCNetwork) {
		rawNew.LinkedVPCNetwork = rawDesired.LinkedVPCNetwork
	} else {
		rawNew.LinkedVPCNetwork = canonicalizeNewSpokeLinkedVPCNetwork(c, rawDesired.LinkedVPCNetwork, rawNew.LinkedVPCNetwork)
	}

	if dcl.IsEmptyValueIndirect(rawNew.UniqueId) && dcl.IsEmptyValueIndirect(rawDesired.UniqueId) {
		rawNew.UniqueId = rawDesired.UniqueId
	} else {
		if dcl.StringCanonicalize(rawDesired.UniqueId, rawNew.UniqueId) {
			rawNew.UniqueId = rawDesired.UniqueId
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeSpokeLinkedVpnTunnels(des, initial *SpokeLinkedVpnTunnels, opts ...dcl.ApplyOption) *SpokeLinkedVpnTunnels {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &SpokeLinkedVpnTunnels{}

	if dcl.StringArrayCanonicalize(des.Uris, initial.Uris) {
		cDes.Uris = initial.Uris
	} else {
		cDes.Uris = des.Uris
	}
	if dcl.BoolCanonicalize(des.SiteToSiteDataTransfer, initial.SiteToSiteDataTransfer) || dcl.IsZeroValue(des.SiteToSiteDataTransfer) {
		cDes.SiteToSiteDataTransfer = initial.SiteToSiteDataTransfer
	} else {
		cDes.SiteToSiteDataTransfer = des.SiteToSiteDataTransfer
	}

	return cDes
}

func canonicalizeSpokeLinkedVpnTunnelsSlice(des, initial []SpokeLinkedVpnTunnels, opts ...dcl.ApplyOption) []SpokeLinkedVpnTunnels {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]SpokeLinkedVpnTunnels, 0, len(des))
		for _, d := range des {
			cd := canonicalizeSpokeLinkedVpnTunnels(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]SpokeLinkedVpnTunnels, 0, len(des))
	for i, d := range des {
		cd := canonicalizeSpokeLinkedVpnTunnels(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewSpokeLinkedVpnTunnels(c *Client, des, nw *SpokeLinkedVpnTunnels) *SpokeLinkedVpnTunnels {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for SpokeLinkedVpnTunnels while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Uris, nw.Uris) {
		nw.Uris = des.Uris
	}
	if dcl.BoolCanonicalize(des.SiteToSiteDataTransfer, nw.SiteToSiteDataTransfer) {
		nw.SiteToSiteDataTransfer = des.SiteToSiteDataTransfer
	}

	return nw
}

func canonicalizeNewSpokeLinkedVpnTunnelsSet(c *Client, des, nw []SpokeLinkedVpnTunnels) []SpokeLinkedVpnTunnels {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []SpokeLinkedVpnTunnels
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareSpokeLinkedVpnTunnelsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewSpokeLinkedVpnTunnels(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewSpokeLinkedVpnTunnelsSlice(c *Client, des, nw []SpokeLinkedVpnTunnels) []SpokeLinkedVpnTunnels {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []SpokeLinkedVpnTunnels
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewSpokeLinkedVpnTunnels(c, &d, &n))
	}

	return items
}

func canonicalizeSpokeLinkedInterconnectAttachments(des, initial *SpokeLinkedInterconnectAttachments, opts ...dcl.ApplyOption) *SpokeLinkedInterconnectAttachments {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &SpokeLinkedInterconnectAttachments{}

	if dcl.StringArrayCanonicalize(des.Uris, initial.Uris) {
		cDes.Uris = initial.Uris
	} else {
		cDes.Uris = des.Uris
	}
	if dcl.BoolCanonicalize(des.SiteToSiteDataTransfer, initial.SiteToSiteDataTransfer) || dcl.IsZeroValue(des.SiteToSiteDataTransfer) {
		cDes.SiteToSiteDataTransfer = initial.SiteToSiteDataTransfer
	} else {
		cDes.SiteToSiteDataTransfer = des.SiteToSiteDataTransfer
	}

	return cDes
}

func canonicalizeSpokeLinkedInterconnectAttachmentsSlice(des, initial []SpokeLinkedInterconnectAttachments, opts ...dcl.ApplyOption) []SpokeLinkedInterconnectAttachments {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]SpokeLinkedInterconnectAttachments, 0, len(des))
		for _, d := range des {
			cd := canonicalizeSpokeLinkedInterconnectAttachments(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]SpokeLinkedInterconnectAttachments, 0, len(des))
	for i, d := range des {
		cd := canonicalizeSpokeLinkedInterconnectAttachments(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewSpokeLinkedInterconnectAttachments(c *Client, des, nw *SpokeLinkedInterconnectAttachments) *SpokeLinkedInterconnectAttachments {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for SpokeLinkedInterconnectAttachments while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.Uris, nw.Uris) {
		nw.Uris = des.Uris
	}
	if dcl.BoolCanonicalize(des.SiteToSiteDataTransfer, nw.SiteToSiteDataTransfer) {
		nw.SiteToSiteDataTransfer = des.SiteToSiteDataTransfer
	}

	return nw
}

func canonicalizeNewSpokeLinkedInterconnectAttachmentsSet(c *Client, des, nw []SpokeLinkedInterconnectAttachments) []SpokeLinkedInterconnectAttachments {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []SpokeLinkedInterconnectAttachments
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareSpokeLinkedInterconnectAttachmentsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewSpokeLinkedInterconnectAttachments(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewSpokeLinkedInterconnectAttachmentsSlice(c *Client, des, nw []SpokeLinkedInterconnectAttachments) []SpokeLinkedInterconnectAttachments {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []SpokeLinkedInterconnectAttachments
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewSpokeLinkedInterconnectAttachments(c, &d, &n))
	}

	return items
}

func canonicalizeSpokeLinkedRouterApplianceInstances(des, initial *SpokeLinkedRouterApplianceInstances, opts ...dcl.ApplyOption) *SpokeLinkedRouterApplianceInstances {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &SpokeLinkedRouterApplianceInstances{}

	cDes.Instances = canonicalizeSpokeLinkedRouterApplianceInstancesInstancesSlice(des.Instances, initial.Instances, opts...)
	if dcl.BoolCanonicalize(des.SiteToSiteDataTransfer, initial.SiteToSiteDataTransfer) || dcl.IsZeroValue(des.SiteToSiteDataTransfer) {
		cDes.SiteToSiteDataTransfer = initial.SiteToSiteDataTransfer
	} else {
		cDes.SiteToSiteDataTransfer = des.SiteToSiteDataTransfer
	}

	return cDes
}

func canonicalizeSpokeLinkedRouterApplianceInstancesSlice(des, initial []SpokeLinkedRouterApplianceInstances, opts ...dcl.ApplyOption) []SpokeLinkedRouterApplianceInstances {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]SpokeLinkedRouterApplianceInstances, 0, len(des))
		for _, d := range des {
			cd := canonicalizeSpokeLinkedRouterApplianceInstances(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]SpokeLinkedRouterApplianceInstances, 0, len(des))
	for i, d := range des {
		cd := canonicalizeSpokeLinkedRouterApplianceInstances(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewSpokeLinkedRouterApplianceInstances(c *Client, des, nw *SpokeLinkedRouterApplianceInstances) *SpokeLinkedRouterApplianceInstances {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for SpokeLinkedRouterApplianceInstances while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Instances = canonicalizeNewSpokeLinkedRouterApplianceInstancesInstancesSlice(c, des.Instances, nw.Instances)
	if dcl.BoolCanonicalize(des.SiteToSiteDataTransfer, nw.SiteToSiteDataTransfer) {
		nw.SiteToSiteDataTransfer = des.SiteToSiteDataTransfer
	}

	return nw
}

func canonicalizeNewSpokeLinkedRouterApplianceInstancesSet(c *Client, des, nw []SpokeLinkedRouterApplianceInstances) []SpokeLinkedRouterApplianceInstances {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []SpokeLinkedRouterApplianceInstances
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareSpokeLinkedRouterApplianceInstancesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewSpokeLinkedRouterApplianceInstances(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewSpokeLinkedRouterApplianceInstancesSlice(c *Client, des, nw []SpokeLinkedRouterApplianceInstances) []SpokeLinkedRouterApplianceInstances {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []SpokeLinkedRouterApplianceInstances
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewSpokeLinkedRouterApplianceInstances(c, &d, &n))
	}

	return items
}

func canonicalizeSpokeLinkedRouterApplianceInstancesInstances(des, initial *SpokeLinkedRouterApplianceInstancesInstances, opts ...dcl.ApplyOption) *SpokeLinkedRouterApplianceInstancesInstances {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &SpokeLinkedRouterApplianceInstancesInstances{}

	if dcl.IsZeroValue(des.VirtualMachine) || (dcl.IsEmptyValueIndirect(des.VirtualMachine) && dcl.IsEmptyValueIndirect(initial.VirtualMachine)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.VirtualMachine = initial.VirtualMachine
	} else {
		cDes.VirtualMachine = des.VirtualMachine
	}
	if dcl.StringCanonicalize(des.IPAddress, initial.IPAddress) || dcl.IsZeroValue(des.IPAddress) {
		cDes.IPAddress = initial.IPAddress
	} else {
		cDes.IPAddress = des.IPAddress
	}

	return cDes
}

func canonicalizeSpokeLinkedRouterApplianceInstancesInstancesSlice(des, initial []SpokeLinkedRouterApplianceInstancesInstances, opts ...dcl.ApplyOption) []SpokeLinkedRouterApplianceInstancesInstances {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]SpokeLinkedRouterApplianceInstancesInstances, 0, len(des))
		for _, d := range des {
			cd := canonicalizeSpokeLinkedRouterApplianceInstancesInstances(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]SpokeLinkedRouterApplianceInstancesInstances, 0, len(des))
	for i, d := range des {
		cd := canonicalizeSpokeLinkedRouterApplianceInstancesInstances(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewSpokeLinkedRouterApplianceInstancesInstances(c *Client, des, nw *SpokeLinkedRouterApplianceInstancesInstances) *SpokeLinkedRouterApplianceInstancesInstances {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for SpokeLinkedRouterApplianceInstancesInstances while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.IPAddress, nw.IPAddress) {
		nw.IPAddress = des.IPAddress
	}

	return nw
}

func canonicalizeNewSpokeLinkedRouterApplianceInstancesInstancesSet(c *Client, des, nw []SpokeLinkedRouterApplianceInstancesInstances) []SpokeLinkedRouterApplianceInstancesInstances {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []SpokeLinkedRouterApplianceInstancesInstances
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareSpokeLinkedRouterApplianceInstancesInstancesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewSpokeLinkedRouterApplianceInstancesInstances(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewSpokeLinkedRouterApplianceInstancesInstancesSlice(c *Client, des, nw []SpokeLinkedRouterApplianceInstancesInstances) []SpokeLinkedRouterApplianceInstancesInstances {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []SpokeLinkedRouterApplianceInstancesInstances
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewSpokeLinkedRouterApplianceInstancesInstances(c, &d, &n))
	}

	return items
}

func canonicalizeSpokeLinkedVPCNetwork(des, initial *SpokeLinkedVPCNetwork, opts ...dcl.ApplyOption) *SpokeLinkedVPCNetwork {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &SpokeLinkedVPCNetwork{}

	if dcl.IsZeroValue(des.Uri) || (dcl.IsEmptyValueIndirect(des.Uri) && dcl.IsEmptyValueIndirect(initial.Uri)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Uri = initial.Uri
	} else {
		cDes.Uri = des.Uri
	}
	if dcl.StringArrayCanonicalize(des.ExcludeExportRanges, initial.ExcludeExportRanges) {
		cDes.ExcludeExportRanges = initial.ExcludeExportRanges
	} else {
		cDes.ExcludeExportRanges = des.ExcludeExportRanges
	}

	return cDes
}

func canonicalizeSpokeLinkedVPCNetworkSlice(des, initial []SpokeLinkedVPCNetwork, opts ...dcl.ApplyOption) []SpokeLinkedVPCNetwork {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]SpokeLinkedVPCNetwork, 0, len(des))
		for _, d := range des {
			cd := canonicalizeSpokeLinkedVPCNetwork(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]SpokeLinkedVPCNetwork, 0, len(des))
	for i, d := range des {
		cd := canonicalizeSpokeLinkedVPCNetwork(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewSpokeLinkedVPCNetwork(c *Client, des, nw *SpokeLinkedVPCNetwork) *SpokeLinkedVPCNetwork {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for SpokeLinkedVPCNetwork while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.ExcludeExportRanges, nw.ExcludeExportRanges) {
		nw.ExcludeExportRanges = des.ExcludeExportRanges
	}

	return nw
}

func canonicalizeNewSpokeLinkedVPCNetworkSet(c *Client, des, nw []SpokeLinkedVPCNetwork) []SpokeLinkedVPCNetwork {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []SpokeLinkedVPCNetwork
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareSpokeLinkedVPCNetworkNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewSpokeLinkedVPCNetwork(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewSpokeLinkedVPCNetworkSlice(c *Client, des, nw []SpokeLinkedVPCNetwork) []SpokeLinkedVPCNetwork {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []SpokeLinkedVPCNetwork
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewSpokeLinkedVPCNetwork(c, &d, &n))
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
func diffSpoke(c *Client, desired, actual *Spoke, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateSpokeUpdateSpokeOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateSpokeUpdateSpokeOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Hub, actual.Hub, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Hub")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LinkedVpnTunnels, actual.LinkedVpnTunnels, dcl.DiffInfo{ObjectFunction: compareSpokeLinkedVpnTunnelsNewStyle, EmptyObject: EmptySpokeLinkedVpnTunnels, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LinkedVpnTunnels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LinkedInterconnectAttachments, actual.LinkedInterconnectAttachments, dcl.DiffInfo{ObjectFunction: compareSpokeLinkedInterconnectAttachmentsNewStyle, EmptyObject: EmptySpokeLinkedInterconnectAttachments, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LinkedInterconnectAttachments")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LinkedRouterApplianceInstances, actual.LinkedRouterApplianceInstances, dcl.DiffInfo{ObjectFunction: compareSpokeLinkedRouterApplianceInstancesNewStyle, EmptyObject: EmptySpokeLinkedRouterApplianceInstances, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LinkedRouterApplianceInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LinkedVPCNetwork, actual.LinkedVPCNetwork, dcl.DiffInfo{ObjectFunction: compareSpokeLinkedVPCNetworkNewStyle, EmptyObject: EmptySpokeLinkedVPCNetwork, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LinkedVpcNetwork")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UniqueId, actual.UniqueId, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UniqueId")); len(ds) != 0 || err != nil {
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

	if len(newDiffs) > 0 {
		c.Config.Logger.Infof("Diff function found diffs: %v", newDiffs)
	}
	return newDiffs, nil
}
func compareSpokeLinkedVpnTunnelsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*SpokeLinkedVpnTunnels)
	if !ok {
		desiredNotPointer, ok := d.(SpokeLinkedVpnTunnels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedVpnTunnels or *SpokeLinkedVpnTunnels", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*SpokeLinkedVpnTunnels)
	if !ok {
		actualNotPointer, ok := a.(SpokeLinkedVpnTunnels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedVpnTunnels", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uris, actual.Uris, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SiteToSiteDataTransfer, actual.SiteToSiteDataTransfer, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SiteToSiteDataTransfer")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareSpokeLinkedInterconnectAttachmentsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*SpokeLinkedInterconnectAttachments)
	if !ok {
		desiredNotPointer, ok := d.(SpokeLinkedInterconnectAttachments)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedInterconnectAttachments or *SpokeLinkedInterconnectAttachments", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*SpokeLinkedInterconnectAttachments)
	if !ok {
		actualNotPointer, ok := a.(SpokeLinkedInterconnectAttachments)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedInterconnectAttachments", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uris, actual.Uris, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SiteToSiteDataTransfer, actual.SiteToSiteDataTransfer, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SiteToSiteDataTransfer")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareSpokeLinkedRouterApplianceInstancesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*SpokeLinkedRouterApplianceInstances)
	if !ok {
		desiredNotPointer, ok := d.(SpokeLinkedRouterApplianceInstances)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedRouterApplianceInstances or *SpokeLinkedRouterApplianceInstances", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*SpokeLinkedRouterApplianceInstances)
	if !ok {
		actualNotPointer, ok := a.(SpokeLinkedRouterApplianceInstances)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedRouterApplianceInstances", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Instances, actual.Instances, dcl.DiffInfo{ObjectFunction: compareSpokeLinkedRouterApplianceInstancesInstancesNewStyle, EmptyObject: EmptySpokeLinkedRouterApplianceInstancesInstances, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Instances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SiteToSiteDataTransfer, actual.SiteToSiteDataTransfer, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SiteToSiteDataTransfer")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareSpokeLinkedRouterApplianceInstancesInstancesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*SpokeLinkedRouterApplianceInstancesInstances)
	if !ok {
		desiredNotPointer, ok := d.(SpokeLinkedRouterApplianceInstancesInstances)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedRouterApplianceInstancesInstances or *SpokeLinkedRouterApplianceInstancesInstances", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*SpokeLinkedRouterApplianceInstancesInstances)
	if !ok {
		actualNotPointer, ok := a.(SpokeLinkedRouterApplianceInstancesInstances)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedRouterApplianceInstancesInstances", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.VirtualMachine, actual.VirtualMachine, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("VirtualMachine")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPAddress, actual.IPAddress, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IpAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareSpokeLinkedVPCNetworkNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*SpokeLinkedVPCNetwork)
	if !ok {
		desiredNotPointer, ok := d.(SpokeLinkedVPCNetwork)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedVPCNetwork or *SpokeLinkedVPCNetwork", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*SpokeLinkedVPCNetwork)
	if !ok {
		actualNotPointer, ok := a.(SpokeLinkedVPCNetwork)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a SpokeLinkedVPCNetwork", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExcludeExportRanges, actual.ExcludeExportRanges, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExcludeExportRanges")); len(ds) != 0 || err != nil {
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
func (r *Spoke) urlNormalized() *Spoke {
	normalized := dcl.Copy(*r).(Spoke)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Hub = dcl.SelfLinkToName(r.Hub)
	normalized.UniqueId = dcl.SelfLinkToName(r.UniqueId)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *Spoke) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateSpoke" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/spokes/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Spoke resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Spoke) marshal(c *Client) ([]byte, error) {
	m, err := expandSpoke(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Spoke: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalSpoke decodes JSON responses into the Spoke resource schema.
func unmarshalSpoke(b []byte, c *Client, res *Spoke) (*Spoke, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapSpoke(m, c, res)
}

func unmarshalMapSpoke(m map[string]interface{}, c *Client, res *Spoke) (*Spoke, error) {

	flattened := flattenSpoke(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandSpoke expands Spoke into a JSON request object.
func expandSpoke(c *Client, f *Spoke) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/spokes/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Hub; dcl.ValueShouldBeSent(v) {
		m["hub"] = v
	}
	if v, err := expandSpokeLinkedVpnTunnels(c, f.LinkedVpnTunnels, res); err != nil {
		return nil, fmt.Errorf("error expanding LinkedVpnTunnels into linkedVpnTunnels: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["linkedVpnTunnels"] = v
	}
	if v, err := expandSpokeLinkedInterconnectAttachments(c, f.LinkedInterconnectAttachments, res); err != nil {
		return nil, fmt.Errorf("error expanding LinkedInterconnectAttachments into linkedInterconnectAttachments: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["linkedInterconnectAttachments"] = v
	}
	if v, err := expandSpokeLinkedRouterApplianceInstances(c, f.LinkedRouterApplianceInstances, res); err != nil {
		return nil, fmt.Errorf("error expanding LinkedRouterApplianceInstances into linkedRouterApplianceInstances: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["linkedRouterApplianceInstances"] = v
	}
	if v, err := expandSpokeLinkedVPCNetwork(c, f.LinkedVPCNetwork, res); err != nil {
		return nil, fmt.Errorf("error expanding LinkedVPCNetwork into linkedVpcNetwork: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["linkedVpcNetwork"] = v
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

// flattenSpoke flattens Spoke from a JSON request object into the
// Spoke type.
func flattenSpoke(c *Client, i interface{}, res *Spoke) *Spoke {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Spoke{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Hub = dcl.FlattenString(m["hub"])
	resultRes.LinkedVpnTunnels = flattenSpokeLinkedVpnTunnels(c, m["linkedVpnTunnels"], res)
	resultRes.LinkedInterconnectAttachments = flattenSpokeLinkedInterconnectAttachments(c, m["linkedInterconnectAttachments"], res)
	resultRes.LinkedRouterApplianceInstances = flattenSpokeLinkedRouterApplianceInstances(c, m["linkedRouterApplianceInstances"], res)
	resultRes.LinkedVPCNetwork = flattenSpokeLinkedVPCNetwork(c, m["linkedVpcNetwork"], res)
	resultRes.UniqueId = dcl.FlattenString(m["uniqueId"])
	resultRes.State = flattenSpokeStateEnum(m["state"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandSpokeLinkedVpnTunnelsMap expands the contents of SpokeLinkedVpnTunnels into a JSON
// request object.
func expandSpokeLinkedVpnTunnelsMap(c *Client, f map[string]SpokeLinkedVpnTunnels, res *Spoke) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandSpokeLinkedVpnTunnels(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandSpokeLinkedVpnTunnelsSlice expands the contents of SpokeLinkedVpnTunnels into a JSON
// request object.
func expandSpokeLinkedVpnTunnelsSlice(c *Client, f []SpokeLinkedVpnTunnels, res *Spoke) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandSpokeLinkedVpnTunnels(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenSpokeLinkedVpnTunnelsMap flattens the contents of SpokeLinkedVpnTunnels from a JSON
// response object.
func flattenSpokeLinkedVpnTunnelsMap(c *Client, i interface{}, res *Spoke) map[string]SpokeLinkedVpnTunnels {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SpokeLinkedVpnTunnels{}
	}

	if len(a) == 0 {
		return map[string]SpokeLinkedVpnTunnels{}
	}

	items := make(map[string]SpokeLinkedVpnTunnels)
	for k, item := range a {
		items[k] = *flattenSpokeLinkedVpnTunnels(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenSpokeLinkedVpnTunnelsSlice flattens the contents of SpokeLinkedVpnTunnels from a JSON
// response object.
func flattenSpokeLinkedVpnTunnelsSlice(c *Client, i interface{}, res *Spoke) []SpokeLinkedVpnTunnels {
	a, ok := i.([]interface{})
	if !ok {
		return []SpokeLinkedVpnTunnels{}
	}

	if len(a) == 0 {
		return []SpokeLinkedVpnTunnels{}
	}

	items := make([]SpokeLinkedVpnTunnels, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSpokeLinkedVpnTunnels(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandSpokeLinkedVpnTunnels expands an instance of SpokeLinkedVpnTunnels into a JSON
// request object.
func expandSpokeLinkedVpnTunnels(c *Client, f *SpokeLinkedVpnTunnels, res *Spoke) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Uris; v != nil {
		m["uris"] = v
	}
	if v := f.SiteToSiteDataTransfer; !dcl.IsEmptyValueIndirect(v) {
		m["siteToSiteDataTransfer"] = v
	}

	return m, nil
}

// flattenSpokeLinkedVpnTunnels flattens an instance of SpokeLinkedVpnTunnels from a JSON
// response object.
func flattenSpokeLinkedVpnTunnels(c *Client, i interface{}, res *Spoke) *SpokeLinkedVpnTunnels {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &SpokeLinkedVpnTunnels{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptySpokeLinkedVpnTunnels
	}
	r.Uris = dcl.FlattenStringSlice(m["uris"])
	r.SiteToSiteDataTransfer = dcl.FlattenBool(m["siteToSiteDataTransfer"])

	return r
}

// expandSpokeLinkedInterconnectAttachmentsMap expands the contents of SpokeLinkedInterconnectAttachments into a JSON
// request object.
func expandSpokeLinkedInterconnectAttachmentsMap(c *Client, f map[string]SpokeLinkedInterconnectAttachments, res *Spoke) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandSpokeLinkedInterconnectAttachments(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandSpokeLinkedInterconnectAttachmentsSlice expands the contents of SpokeLinkedInterconnectAttachments into a JSON
// request object.
func expandSpokeLinkedInterconnectAttachmentsSlice(c *Client, f []SpokeLinkedInterconnectAttachments, res *Spoke) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandSpokeLinkedInterconnectAttachments(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenSpokeLinkedInterconnectAttachmentsMap flattens the contents of SpokeLinkedInterconnectAttachments from a JSON
// response object.
func flattenSpokeLinkedInterconnectAttachmentsMap(c *Client, i interface{}, res *Spoke) map[string]SpokeLinkedInterconnectAttachments {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SpokeLinkedInterconnectAttachments{}
	}

	if len(a) == 0 {
		return map[string]SpokeLinkedInterconnectAttachments{}
	}

	items := make(map[string]SpokeLinkedInterconnectAttachments)
	for k, item := range a {
		items[k] = *flattenSpokeLinkedInterconnectAttachments(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenSpokeLinkedInterconnectAttachmentsSlice flattens the contents of SpokeLinkedInterconnectAttachments from a JSON
// response object.
func flattenSpokeLinkedInterconnectAttachmentsSlice(c *Client, i interface{}, res *Spoke) []SpokeLinkedInterconnectAttachments {
	a, ok := i.([]interface{})
	if !ok {
		return []SpokeLinkedInterconnectAttachments{}
	}

	if len(a) == 0 {
		return []SpokeLinkedInterconnectAttachments{}
	}

	items := make([]SpokeLinkedInterconnectAttachments, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSpokeLinkedInterconnectAttachments(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandSpokeLinkedInterconnectAttachments expands an instance of SpokeLinkedInterconnectAttachments into a JSON
// request object.
func expandSpokeLinkedInterconnectAttachments(c *Client, f *SpokeLinkedInterconnectAttachments, res *Spoke) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Uris; v != nil {
		m["uris"] = v
	}
	if v := f.SiteToSiteDataTransfer; !dcl.IsEmptyValueIndirect(v) {
		m["siteToSiteDataTransfer"] = v
	}

	return m, nil
}

// flattenSpokeLinkedInterconnectAttachments flattens an instance of SpokeLinkedInterconnectAttachments from a JSON
// response object.
func flattenSpokeLinkedInterconnectAttachments(c *Client, i interface{}, res *Spoke) *SpokeLinkedInterconnectAttachments {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &SpokeLinkedInterconnectAttachments{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptySpokeLinkedInterconnectAttachments
	}
	r.Uris = dcl.FlattenStringSlice(m["uris"])
	r.SiteToSiteDataTransfer = dcl.FlattenBool(m["siteToSiteDataTransfer"])

	return r
}

// expandSpokeLinkedRouterApplianceInstancesMap expands the contents of SpokeLinkedRouterApplianceInstances into a JSON
// request object.
func expandSpokeLinkedRouterApplianceInstancesMap(c *Client, f map[string]SpokeLinkedRouterApplianceInstances, res *Spoke) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandSpokeLinkedRouterApplianceInstances(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandSpokeLinkedRouterApplianceInstancesSlice expands the contents of SpokeLinkedRouterApplianceInstances into a JSON
// request object.
func expandSpokeLinkedRouterApplianceInstancesSlice(c *Client, f []SpokeLinkedRouterApplianceInstances, res *Spoke) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandSpokeLinkedRouterApplianceInstances(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenSpokeLinkedRouterApplianceInstancesMap flattens the contents of SpokeLinkedRouterApplianceInstances from a JSON
// response object.
func flattenSpokeLinkedRouterApplianceInstancesMap(c *Client, i interface{}, res *Spoke) map[string]SpokeLinkedRouterApplianceInstances {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SpokeLinkedRouterApplianceInstances{}
	}

	if len(a) == 0 {
		return map[string]SpokeLinkedRouterApplianceInstances{}
	}

	items := make(map[string]SpokeLinkedRouterApplianceInstances)
	for k, item := range a {
		items[k] = *flattenSpokeLinkedRouterApplianceInstances(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenSpokeLinkedRouterApplianceInstancesSlice flattens the contents of SpokeLinkedRouterApplianceInstances from a JSON
// response object.
func flattenSpokeLinkedRouterApplianceInstancesSlice(c *Client, i interface{}, res *Spoke) []SpokeLinkedRouterApplianceInstances {
	a, ok := i.([]interface{})
	if !ok {
		return []SpokeLinkedRouterApplianceInstances{}
	}

	if len(a) == 0 {
		return []SpokeLinkedRouterApplianceInstances{}
	}

	items := make([]SpokeLinkedRouterApplianceInstances, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSpokeLinkedRouterApplianceInstances(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandSpokeLinkedRouterApplianceInstances expands an instance of SpokeLinkedRouterApplianceInstances into a JSON
// request object.
func expandSpokeLinkedRouterApplianceInstances(c *Client, f *SpokeLinkedRouterApplianceInstances, res *Spoke) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandSpokeLinkedRouterApplianceInstancesInstancesSlice(c, f.Instances, res); err != nil {
		return nil, fmt.Errorf("error expanding Instances into instances: %w", err)
	} else if v != nil {
		m["instances"] = v
	}
	if v := f.SiteToSiteDataTransfer; !dcl.IsEmptyValueIndirect(v) {
		m["siteToSiteDataTransfer"] = v
	}

	return m, nil
}

// flattenSpokeLinkedRouterApplianceInstances flattens an instance of SpokeLinkedRouterApplianceInstances from a JSON
// response object.
func flattenSpokeLinkedRouterApplianceInstances(c *Client, i interface{}, res *Spoke) *SpokeLinkedRouterApplianceInstances {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &SpokeLinkedRouterApplianceInstances{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptySpokeLinkedRouterApplianceInstances
	}
	r.Instances = flattenSpokeLinkedRouterApplianceInstancesInstancesSlice(c, m["instances"], res)
	r.SiteToSiteDataTransfer = dcl.FlattenBool(m["siteToSiteDataTransfer"])

	return r
}

// expandSpokeLinkedRouterApplianceInstancesInstancesMap expands the contents of SpokeLinkedRouterApplianceInstancesInstances into a JSON
// request object.
func expandSpokeLinkedRouterApplianceInstancesInstancesMap(c *Client, f map[string]SpokeLinkedRouterApplianceInstancesInstances, res *Spoke) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandSpokeLinkedRouterApplianceInstancesInstances(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandSpokeLinkedRouterApplianceInstancesInstancesSlice expands the contents of SpokeLinkedRouterApplianceInstancesInstances into a JSON
// request object.
func expandSpokeLinkedRouterApplianceInstancesInstancesSlice(c *Client, f []SpokeLinkedRouterApplianceInstancesInstances, res *Spoke) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandSpokeLinkedRouterApplianceInstancesInstances(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenSpokeLinkedRouterApplianceInstancesInstancesMap flattens the contents of SpokeLinkedRouterApplianceInstancesInstances from a JSON
// response object.
func flattenSpokeLinkedRouterApplianceInstancesInstancesMap(c *Client, i interface{}, res *Spoke) map[string]SpokeLinkedRouterApplianceInstancesInstances {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SpokeLinkedRouterApplianceInstancesInstances{}
	}

	if len(a) == 0 {
		return map[string]SpokeLinkedRouterApplianceInstancesInstances{}
	}

	items := make(map[string]SpokeLinkedRouterApplianceInstancesInstances)
	for k, item := range a {
		items[k] = *flattenSpokeLinkedRouterApplianceInstancesInstances(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenSpokeLinkedRouterApplianceInstancesInstancesSlice flattens the contents of SpokeLinkedRouterApplianceInstancesInstances from a JSON
// response object.
func flattenSpokeLinkedRouterApplianceInstancesInstancesSlice(c *Client, i interface{}, res *Spoke) []SpokeLinkedRouterApplianceInstancesInstances {
	a, ok := i.([]interface{})
	if !ok {
		return []SpokeLinkedRouterApplianceInstancesInstances{}
	}

	if len(a) == 0 {
		return []SpokeLinkedRouterApplianceInstancesInstances{}
	}

	items := make([]SpokeLinkedRouterApplianceInstancesInstances, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSpokeLinkedRouterApplianceInstancesInstances(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandSpokeLinkedRouterApplianceInstancesInstances expands an instance of SpokeLinkedRouterApplianceInstancesInstances into a JSON
// request object.
func expandSpokeLinkedRouterApplianceInstancesInstances(c *Client, f *SpokeLinkedRouterApplianceInstancesInstances, res *Spoke) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.VirtualMachine; !dcl.IsEmptyValueIndirect(v) {
		m["virtualMachine"] = v
	}
	if v := f.IPAddress; !dcl.IsEmptyValueIndirect(v) {
		m["ipAddress"] = v
	}

	return m, nil
}

// flattenSpokeLinkedRouterApplianceInstancesInstances flattens an instance of SpokeLinkedRouterApplianceInstancesInstances from a JSON
// response object.
func flattenSpokeLinkedRouterApplianceInstancesInstances(c *Client, i interface{}, res *Spoke) *SpokeLinkedRouterApplianceInstancesInstances {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &SpokeLinkedRouterApplianceInstancesInstances{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptySpokeLinkedRouterApplianceInstancesInstances
	}
	r.VirtualMachine = dcl.FlattenString(m["virtualMachine"])
	r.IPAddress = dcl.FlattenString(m["ipAddress"])

	return r
}

// expandSpokeLinkedVPCNetworkMap expands the contents of SpokeLinkedVPCNetwork into a JSON
// request object.
func expandSpokeLinkedVPCNetworkMap(c *Client, f map[string]SpokeLinkedVPCNetwork, res *Spoke) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandSpokeLinkedVPCNetwork(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandSpokeLinkedVPCNetworkSlice expands the contents of SpokeLinkedVPCNetwork into a JSON
// request object.
func expandSpokeLinkedVPCNetworkSlice(c *Client, f []SpokeLinkedVPCNetwork, res *Spoke) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandSpokeLinkedVPCNetwork(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenSpokeLinkedVPCNetworkMap flattens the contents of SpokeLinkedVPCNetwork from a JSON
// response object.
func flattenSpokeLinkedVPCNetworkMap(c *Client, i interface{}, res *Spoke) map[string]SpokeLinkedVPCNetwork {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SpokeLinkedVPCNetwork{}
	}

	if len(a) == 0 {
		return map[string]SpokeLinkedVPCNetwork{}
	}

	items := make(map[string]SpokeLinkedVPCNetwork)
	for k, item := range a {
		items[k] = *flattenSpokeLinkedVPCNetwork(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenSpokeLinkedVPCNetworkSlice flattens the contents of SpokeLinkedVPCNetwork from a JSON
// response object.
func flattenSpokeLinkedVPCNetworkSlice(c *Client, i interface{}, res *Spoke) []SpokeLinkedVPCNetwork {
	a, ok := i.([]interface{})
	if !ok {
		return []SpokeLinkedVPCNetwork{}
	}

	if len(a) == 0 {
		return []SpokeLinkedVPCNetwork{}
	}

	items := make([]SpokeLinkedVPCNetwork, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSpokeLinkedVPCNetwork(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandSpokeLinkedVPCNetwork expands an instance of SpokeLinkedVPCNetwork into a JSON
// request object.
func expandSpokeLinkedVPCNetwork(c *Client, f *SpokeLinkedVPCNetwork, res *Spoke) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Uri; !dcl.IsEmptyValueIndirect(v) {
		m["uri"] = v
	}
	if v := f.ExcludeExportRanges; v != nil {
		m["excludeExportRanges"] = v
	}

	return m, nil
}

// flattenSpokeLinkedVPCNetwork flattens an instance of SpokeLinkedVPCNetwork from a JSON
// response object.
func flattenSpokeLinkedVPCNetwork(c *Client, i interface{}, res *Spoke) *SpokeLinkedVPCNetwork {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &SpokeLinkedVPCNetwork{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptySpokeLinkedVPCNetwork
	}
	r.Uri = dcl.FlattenString(m["uri"])
	r.ExcludeExportRanges = dcl.FlattenStringSlice(m["excludeExportRanges"])

	return r
}

// flattenSpokeStateEnumMap flattens the contents of SpokeStateEnum from a JSON
// response object.
func flattenSpokeStateEnumMap(c *Client, i interface{}, res *Spoke) map[string]SpokeStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]SpokeStateEnum{}
	}

	if len(a) == 0 {
		return map[string]SpokeStateEnum{}
	}

	items := make(map[string]SpokeStateEnum)
	for k, item := range a {
		items[k] = *flattenSpokeStateEnum(item.(interface{}))
	}

	return items
}

// flattenSpokeStateEnumSlice flattens the contents of SpokeStateEnum from a JSON
// response object.
func flattenSpokeStateEnumSlice(c *Client, i interface{}, res *Spoke) []SpokeStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []SpokeStateEnum{}
	}

	if len(a) == 0 {
		return []SpokeStateEnum{}
	}

	items := make([]SpokeStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenSpokeStateEnum(item.(interface{})))
	}

	return items
}

// flattenSpokeStateEnum asserts that an interface is a string, and returns a
// pointer to a *SpokeStateEnum with the same value as that string.
func flattenSpokeStateEnum(i interface{}) *SpokeStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return SpokeStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Spoke) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalSpoke(b, c, r)
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

type spokeDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         spokeApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToSpokeDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]spokeDiff, error) {
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
	var diffs []spokeDiff
	// For each operation name, create a spokeDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := spokeDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToSpokeApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToSpokeApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (spokeApiOperation, error) {
	switch opName {

	case "updateSpokeUpdateSpokeOperation":
		return &updateSpokeUpdateSpokeOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractSpokeFields(r *Spoke) error {
	vLinkedVpnTunnels := r.LinkedVpnTunnels
	if vLinkedVpnTunnels == nil {
		// note: explicitly not the empty object.
		vLinkedVpnTunnels = &SpokeLinkedVpnTunnels{}
	}
	if err := extractSpokeLinkedVpnTunnelsFields(r, vLinkedVpnTunnels); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinkedVpnTunnels) {
		r.LinkedVpnTunnels = vLinkedVpnTunnels
	}
	vLinkedInterconnectAttachments := r.LinkedInterconnectAttachments
	if vLinkedInterconnectAttachments == nil {
		// note: explicitly not the empty object.
		vLinkedInterconnectAttachments = &SpokeLinkedInterconnectAttachments{}
	}
	if err := extractSpokeLinkedInterconnectAttachmentsFields(r, vLinkedInterconnectAttachments); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinkedInterconnectAttachments) {
		r.LinkedInterconnectAttachments = vLinkedInterconnectAttachments
	}
	vLinkedRouterApplianceInstances := r.LinkedRouterApplianceInstances
	if vLinkedRouterApplianceInstances == nil {
		// note: explicitly not the empty object.
		vLinkedRouterApplianceInstances = &SpokeLinkedRouterApplianceInstances{}
	}
	if err := extractSpokeLinkedRouterApplianceInstancesFields(r, vLinkedRouterApplianceInstances); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinkedRouterApplianceInstances) {
		r.LinkedRouterApplianceInstances = vLinkedRouterApplianceInstances
	}
	vLinkedVPCNetwork := r.LinkedVPCNetwork
	if vLinkedVPCNetwork == nil {
		// note: explicitly not the empty object.
		vLinkedVPCNetwork = &SpokeLinkedVPCNetwork{}
	}
	if err := extractSpokeLinkedVPCNetworkFields(r, vLinkedVPCNetwork); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinkedVPCNetwork) {
		r.LinkedVPCNetwork = vLinkedVPCNetwork
	}
	return nil
}
func extractSpokeLinkedVpnTunnelsFields(r *Spoke, o *SpokeLinkedVpnTunnels) error {
	return nil
}
func extractSpokeLinkedInterconnectAttachmentsFields(r *Spoke, o *SpokeLinkedInterconnectAttachments) error {
	return nil
}
func extractSpokeLinkedRouterApplianceInstancesFields(r *Spoke, o *SpokeLinkedRouterApplianceInstances) error {
	return nil
}
func extractSpokeLinkedRouterApplianceInstancesInstancesFields(r *Spoke, o *SpokeLinkedRouterApplianceInstancesInstances) error {
	return nil
}
func extractSpokeLinkedVPCNetworkFields(r *Spoke, o *SpokeLinkedVPCNetwork) error {
	return nil
}

func postReadExtractSpokeFields(r *Spoke) error {
	vLinkedVpnTunnels := r.LinkedVpnTunnels
	if vLinkedVpnTunnels == nil {
		// note: explicitly not the empty object.
		vLinkedVpnTunnels = &SpokeLinkedVpnTunnels{}
	}
	if err := postReadExtractSpokeLinkedVpnTunnelsFields(r, vLinkedVpnTunnels); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinkedVpnTunnels) {
		r.LinkedVpnTunnels = vLinkedVpnTunnels
	}
	vLinkedInterconnectAttachments := r.LinkedInterconnectAttachments
	if vLinkedInterconnectAttachments == nil {
		// note: explicitly not the empty object.
		vLinkedInterconnectAttachments = &SpokeLinkedInterconnectAttachments{}
	}
	if err := postReadExtractSpokeLinkedInterconnectAttachmentsFields(r, vLinkedInterconnectAttachments); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinkedInterconnectAttachments) {
		r.LinkedInterconnectAttachments = vLinkedInterconnectAttachments
	}
	vLinkedRouterApplianceInstances := r.LinkedRouterApplianceInstances
	if vLinkedRouterApplianceInstances == nil {
		// note: explicitly not the empty object.
		vLinkedRouterApplianceInstances = &SpokeLinkedRouterApplianceInstances{}
	}
	if err := postReadExtractSpokeLinkedRouterApplianceInstancesFields(r, vLinkedRouterApplianceInstances); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinkedRouterApplianceInstances) {
		r.LinkedRouterApplianceInstances = vLinkedRouterApplianceInstances
	}
	vLinkedVPCNetwork := r.LinkedVPCNetwork
	if vLinkedVPCNetwork == nil {
		// note: explicitly not the empty object.
		vLinkedVPCNetwork = &SpokeLinkedVPCNetwork{}
	}
	if err := postReadExtractSpokeLinkedVPCNetworkFields(r, vLinkedVPCNetwork); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLinkedVPCNetwork) {
		r.LinkedVPCNetwork = vLinkedVPCNetwork
	}
	return nil
}
func postReadExtractSpokeLinkedVpnTunnelsFields(r *Spoke, o *SpokeLinkedVpnTunnels) error {
	return nil
}
func postReadExtractSpokeLinkedInterconnectAttachmentsFields(r *Spoke, o *SpokeLinkedInterconnectAttachments) error {
	return nil
}
func postReadExtractSpokeLinkedRouterApplianceInstancesFields(r *Spoke, o *SpokeLinkedRouterApplianceInstances) error {
	return nil
}
func postReadExtractSpokeLinkedRouterApplianceInstancesInstancesFields(r *Spoke, o *SpokeLinkedRouterApplianceInstancesInstances) error {
	return nil
}
func postReadExtractSpokeLinkedVPCNetworkFields(r *Spoke, o *SpokeLinkedVPCNetwork) error {
	return nil
}
