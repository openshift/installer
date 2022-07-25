// Copyright 2021 Google LLC. All Rights Reserved.
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
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *Interconnect) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "location"); err != nil {
		return err
	}
	if err := dcl.Required(r, "linkType"); err != nil {
		return err
	}
	if err := dcl.Required(r, "interconnectType"); err != nil {
		return err
	}
	if err := dcl.Required(r, "customerName"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	return nil
}
func (r *InterconnectExpectedOutages) validate() error {
	return nil
}
func (r *InterconnectCircuitInfos) validate() error {
	return nil
}
func (r *Interconnect) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *Interconnect) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/global/interconnects/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Interconnect) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/global/interconnects", nr.basePath(), userBasePath, params), nil

}

func (r *Interconnect) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
	}
	return dcl.URL("projects/{{project}}/global/interconnects", nr.basePath(), userBasePath, params), nil

}

func (r *Interconnect) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/global/interconnects/{{name}}", nr.basePath(), userBasePath, params), nil
}

// interconnectApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type interconnectApiOperation interface {
	do(context.Context, *Interconnect, *Client) error
}

// newUpdateInterconnectPatchRequest creates a request for an
// Interconnect resource's Patch update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInterconnectPatchRequest(ctx context.Context, f *Interconnect, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		req["name"] = v
	}
	if v := f.Location; !dcl.IsEmptyValueIndirect(v) {
		req["location"] = v
	}
	if v := f.LinkType; !dcl.IsEmptyValueIndirect(v) {
		req["linkType"] = v
	}
	if v := f.RequestedLinkCount; !dcl.IsEmptyValueIndirect(v) {
		req["requestedLinkCount"] = v
	}
	if v := f.InterconnectType; !dcl.IsEmptyValueIndirect(v) {
		req["interconnectType"] = v
	}
	if v := f.AdminEnabled; !dcl.IsEmptyValueIndirect(v) {
		req["adminEnabled"] = v
	}
	if v := f.NocContactEmail; !dcl.IsEmptyValueIndirect(v) {
		req["nocContactEmail"] = v
	}
	if v := f.CustomerName; !dcl.IsEmptyValueIndirect(v) {
		req["customerName"] = v
	}
	return req, nil
}

// marshalUpdateInterconnectPatchRequest converts the update into
// the final JSON request body.
func marshalUpdateInterconnectPatchRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateInterconnectPatchOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInterconnectPatchOperation) do(ctx context.Context, r *Interconnect, c *Client) error {
	_, err := c.GetInterconnect(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "Patch")
	if err != nil {
		return err
	}

	req, err := newUpdateInterconnectPatchRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInterconnectPatchRequest(c, req)
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

func (c *Client) listInterconnectRaw(ctx context.Context, r *Interconnect, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != InterconnectMaxPage {
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

type listInterconnectOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listInterconnect(ctx context.Context, r *Interconnect, pageToken string, pageSize int32) ([]*Interconnect, string, error) {
	b, err := c.listInterconnectRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listInterconnectOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Interconnect
	for _, v := range m.Items {
		res, err := unmarshalMapInterconnect(v, c)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllInterconnect(ctx context.Context, f func(*Interconnect) bool, resources []*Interconnect) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteInterconnect(ctx, res)
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

type deleteInterconnectOperation struct{}

func (op *deleteInterconnectOperation) do(ctx context.Context, r *Interconnect, c *Client) error {
	r, err := c.GetInterconnect(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Interconnect not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetInterconnect checking for existence. error: %v", err)
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

	// we saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// this is the reason we are adding retry to handle that case.
	maxRetry := 10
	for i := 1; i <= maxRetry; i++ {
		_, err = c.GetInterconnect(ctx, r)
		if !dcl.IsNotFound(err) {
			if i == maxRetry {
				return dcl.NotDeletedError{ExistingResource: r}
			}
			time.Sleep(1000 * time.Millisecond)
		} else {
			break
		}
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createInterconnectOperation struct {
	response map[string]interface{}
}

func (op *createInterconnectOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createInterconnectOperation) do(ctx context.Context, r *Interconnect, c *Client) error {
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

	if _, err := c.GetInterconnect(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getInterconnectRaw(ctx context.Context, r *Interconnect) ([]byte, error) {

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

func (c *Client) interconnectDiffsForRawDesired(ctx context.Context, rawDesired *Interconnect, opts ...dcl.ApplyOption) (initial, desired *Interconnect, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Interconnect
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Interconnect); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Interconnect, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetInterconnect(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Interconnect resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Interconnect resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Interconnect resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeInterconnectDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Interconnect: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Interconnect: %v", rawDesired)

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeInterconnectInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Interconnect: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeInterconnectDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Interconnect: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffInterconnect(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeInterconnectInitialState(rawInitial, rawDesired *Interconnect) (*Interconnect, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeInterconnectDesiredState(rawDesired, rawInitial *Interconnect, opts ...dcl.ApplyOption) (*Interconnect, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.

		return rawDesired, nil
	}
	canonicalDesired := &Interconnect{}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	if dcl.IsZeroValue(rawDesired.LinkType) {
		canonicalDesired.LinkType = rawInitial.LinkType
	} else {
		canonicalDesired.LinkType = rawDesired.LinkType
	}
	if dcl.IsZeroValue(rawDesired.RequestedLinkCount) {
		canonicalDesired.RequestedLinkCount = rawInitial.RequestedLinkCount
	} else {
		canonicalDesired.RequestedLinkCount = rawDesired.RequestedLinkCount
	}
	if dcl.IsZeroValue(rawDesired.InterconnectType) {
		canonicalDesired.InterconnectType = rawInitial.InterconnectType
	} else {
		canonicalDesired.InterconnectType = rawDesired.InterconnectType
	}
	if dcl.BoolCanonicalize(rawDesired.AdminEnabled, rawInitial.AdminEnabled) {
		canonicalDesired.AdminEnabled = rawInitial.AdminEnabled
	} else {
		canonicalDesired.AdminEnabled = rawDesired.AdminEnabled
	}
	if dcl.StringCanonicalize(rawDesired.NocContactEmail, rawInitial.NocContactEmail) {
		canonicalDesired.NocContactEmail = rawInitial.NocContactEmail
	} else {
		canonicalDesired.NocContactEmail = rawDesired.NocContactEmail
	}
	if dcl.StringCanonicalize(rawDesired.CustomerName, rawInitial.CustomerName) {
		canonicalDesired.CustomerName = rawInitial.CustomerName
	} else {
		canonicalDesired.CustomerName = rawDesired.CustomerName
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}

	return canonicalDesired, nil
}

func canonicalizeInterconnectNewState(c *Client, rawNew, rawDesired *Interconnect) (*Interconnect, error) {

	if dcl.IsNotReturnedByServer(rawNew.Description) && dcl.IsNotReturnedByServer(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.SelfLink) && dcl.IsNotReturnedByServer(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Id) && dcl.IsNotReturnedByServer(rawDesired.Id) {
		rawNew.Id = rawDesired.Id
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Name) && dcl.IsNotReturnedByServer(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Location) && dcl.IsNotReturnedByServer(rawDesired.Location) {
		rawNew.Location = rawDesired.Location
	} else {
		if dcl.StringCanonicalize(rawDesired.Location, rawNew.Location) {
			rawNew.Location = rawDesired.Location
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.LinkType) && dcl.IsNotReturnedByServer(rawDesired.LinkType) {
		rawNew.LinkType = rawDesired.LinkType
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.RequestedLinkCount) && dcl.IsNotReturnedByServer(rawDesired.RequestedLinkCount) {
		rawNew.RequestedLinkCount = rawDesired.RequestedLinkCount
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.InterconnectType) && dcl.IsNotReturnedByServer(rawDesired.InterconnectType) {
		rawNew.InterconnectType = rawDesired.InterconnectType
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.AdminEnabled) && dcl.IsNotReturnedByServer(rawDesired.AdminEnabled) {
		rawNew.AdminEnabled = rawDesired.AdminEnabled
	} else {
		if dcl.BoolCanonicalize(rawDesired.AdminEnabled, rawNew.AdminEnabled) {
			rawNew.AdminEnabled = rawDesired.AdminEnabled
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.NocContactEmail) && dcl.IsNotReturnedByServer(rawDesired.NocContactEmail) {
		rawNew.NocContactEmail = rawDesired.NocContactEmail
	} else {
		if dcl.StringCanonicalize(rawDesired.NocContactEmail, rawNew.NocContactEmail) {
			rawNew.NocContactEmail = rawDesired.NocContactEmail
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.CustomerName) && dcl.IsNotReturnedByServer(rawDesired.CustomerName) {
		rawNew.CustomerName = rawDesired.CustomerName
	} else {
		if dcl.StringCanonicalize(rawDesired.CustomerName, rawNew.CustomerName) {
			rawNew.CustomerName = rawDesired.CustomerName
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.OperationalStatus) && dcl.IsNotReturnedByServer(rawDesired.OperationalStatus) {
		rawNew.OperationalStatus = rawDesired.OperationalStatus
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.ProvisionedLinkCount) && dcl.IsNotReturnedByServer(rawDesired.ProvisionedLinkCount) {
		rawNew.ProvisionedLinkCount = rawDesired.ProvisionedLinkCount
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.InterconnectAttachments) && dcl.IsNotReturnedByServer(rawDesired.InterconnectAttachments) {
		rawNew.InterconnectAttachments = rawDesired.InterconnectAttachments
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.InterconnectAttachments, rawNew.InterconnectAttachments) {
			rawNew.InterconnectAttachments = rawDesired.InterconnectAttachments
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.PeerIPAddress) && dcl.IsNotReturnedByServer(rawDesired.PeerIPAddress) {
		rawNew.PeerIPAddress = rawDesired.PeerIPAddress
	} else {
		if dcl.StringCanonicalize(rawDesired.PeerIPAddress, rawNew.PeerIPAddress) {
			rawNew.PeerIPAddress = rawDesired.PeerIPAddress
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.GoogleIPAddress) && dcl.IsNotReturnedByServer(rawDesired.GoogleIPAddress) {
		rawNew.GoogleIPAddress = rawDesired.GoogleIPAddress
	} else {
		if dcl.StringCanonicalize(rawDesired.GoogleIPAddress, rawNew.GoogleIPAddress) {
			rawNew.GoogleIPAddress = rawDesired.GoogleIPAddress
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.GoogleReferenceId) && dcl.IsNotReturnedByServer(rawDesired.GoogleReferenceId) {
		rawNew.GoogleReferenceId = rawDesired.GoogleReferenceId
	} else {
		if dcl.StringCanonicalize(rawDesired.GoogleReferenceId, rawNew.GoogleReferenceId) {
			rawNew.GoogleReferenceId = rawDesired.GoogleReferenceId
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.ExpectedOutages) && dcl.IsNotReturnedByServer(rawDesired.ExpectedOutages) {
		rawNew.ExpectedOutages = rawDesired.ExpectedOutages
	} else {
		rawNew.ExpectedOutages = canonicalizeNewInterconnectExpectedOutagesSlice(c, rawDesired.ExpectedOutages, rawNew.ExpectedOutages)
	}

	if dcl.IsNotReturnedByServer(rawNew.CircuitInfos) && dcl.IsNotReturnedByServer(rawDesired.CircuitInfos) {
		rawNew.CircuitInfos = rawDesired.CircuitInfos
	} else {
		rawNew.CircuitInfos = canonicalizeNewInterconnectCircuitInfosSlice(c, rawDesired.CircuitInfos, rawNew.CircuitInfos)
	}

	if dcl.IsNotReturnedByServer(rawNew.State) && dcl.IsNotReturnedByServer(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	rawNew.Project = rawDesired.Project

	return rawNew, nil
}

func canonicalizeInterconnectExpectedOutages(des, initial *InterconnectExpectedOutages, opts ...dcl.ApplyOption) *InterconnectExpectedOutages {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InterconnectExpectedOutages{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}
	if dcl.IsZeroValue(des.Source) {
		des.Source = initial.Source
	} else {
		cDes.Source = des.Source
	}
	if dcl.IsZeroValue(des.State) {
		des.State = initial.State
	} else {
		cDes.State = des.State
	}
	if dcl.IsZeroValue(des.IssueType) {
		des.IssueType = initial.IssueType
	} else {
		cDes.IssueType = des.IssueType
	}
	if dcl.StringArrayCanonicalize(des.AffectedCircuits, initial.AffectedCircuits) || dcl.IsZeroValue(des.AffectedCircuits) {
		cDes.AffectedCircuits = initial.AffectedCircuits
	} else {
		cDes.AffectedCircuits = des.AffectedCircuits
	}
	if dcl.IsZeroValue(des.StartTime) {
		des.StartTime = initial.StartTime
	} else {
		cDes.StartTime = des.StartTime
	}
	if dcl.IsZeroValue(des.EndTime) {
		des.EndTime = initial.EndTime
	} else {
		cDes.EndTime = des.EndTime
	}

	return cDes
}

func canonicalizeInterconnectExpectedOutagesSlice(des, initial []InterconnectExpectedOutages, opts ...dcl.ApplyOption) []InterconnectExpectedOutages {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InterconnectExpectedOutages, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInterconnectExpectedOutages(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InterconnectExpectedOutages, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInterconnectExpectedOutages(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInterconnectExpectedOutages(c *Client, des, nw *InterconnectExpectedOutages) *InterconnectExpectedOutages {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for InterconnectExpectedOutages while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}
	if dcl.StringArrayCanonicalize(des.AffectedCircuits, nw.AffectedCircuits) {
		nw.AffectedCircuits = des.AffectedCircuits
	}

	return nw
}

func canonicalizeNewInterconnectExpectedOutagesSet(c *Client, des, nw []InterconnectExpectedOutages) []InterconnectExpectedOutages {
	if des == nil {
		return nw
	}
	var reorderedNew []InterconnectExpectedOutages
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareInterconnectExpectedOutagesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewInterconnectExpectedOutagesSlice(c *Client, des, nw []InterconnectExpectedOutages) []InterconnectExpectedOutages {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InterconnectExpectedOutages
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInterconnectExpectedOutages(c, &d, &n))
	}

	return items
}

func canonicalizeInterconnectCircuitInfos(des, initial *InterconnectCircuitInfos, opts ...dcl.ApplyOption) *InterconnectCircuitInfos {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InterconnectCircuitInfos{}

	if dcl.StringCanonicalize(des.GoogleCircuitId, initial.GoogleCircuitId) || dcl.IsZeroValue(des.GoogleCircuitId) {
		cDes.GoogleCircuitId = initial.GoogleCircuitId
	} else {
		cDes.GoogleCircuitId = des.GoogleCircuitId
	}
	if dcl.StringCanonicalize(des.GoogleDemarcId, initial.GoogleDemarcId) || dcl.IsZeroValue(des.GoogleDemarcId) {
		cDes.GoogleDemarcId = initial.GoogleDemarcId
	} else {
		cDes.GoogleDemarcId = des.GoogleDemarcId
	}
	if dcl.StringCanonicalize(des.CustomerDemarcId, initial.CustomerDemarcId) || dcl.IsZeroValue(des.CustomerDemarcId) {
		cDes.CustomerDemarcId = initial.CustomerDemarcId
	} else {
		cDes.CustomerDemarcId = des.CustomerDemarcId
	}

	return cDes
}

func canonicalizeInterconnectCircuitInfosSlice(des, initial []InterconnectCircuitInfos, opts ...dcl.ApplyOption) []InterconnectCircuitInfos {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InterconnectCircuitInfos, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInterconnectCircuitInfos(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InterconnectCircuitInfos, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInterconnectCircuitInfos(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInterconnectCircuitInfos(c *Client, des, nw *InterconnectCircuitInfos) *InterconnectCircuitInfos {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for InterconnectCircuitInfos while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.GoogleCircuitId, nw.GoogleCircuitId) {
		nw.GoogleCircuitId = des.GoogleCircuitId
	}
	if dcl.StringCanonicalize(des.GoogleDemarcId, nw.GoogleDemarcId) {
		nw.GoogleDemarcId = des.GoogleDemarcId
	}
	if dcl.StringCanonicalize(des.CustomerDemarcId, nw.CustomerDemarcId) {
		nw.CustomerDemarcId = des.CustomerDemarcId
	}

	return nw
}

func canonicalizeNewInterconnectCircuitInfosSet(c *Client, des, nw []InterconnectCircuitInfos) []InterconnectCircuitInfos {
	if des == nil {
		return nw
	}
	var reorderedNew []InterconnectCircuitInfos
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareInterconnectCircuitInfosNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedNew = idx
				break
			}
		}
		if matchedNew != -1 {
			reorderedNew = append(reorderedNew, nw[matchedNew])
			nw = append(nw[:matchedNew], nw[matchedNew+1:]...)
		}
	}
	reorderedNew = append(reorderedNew, nw...)

	return reorderedNew
}

func canonicalizeNewInterconnectCircuitInfosSlice(c *Client, des, nw []InterconnectCircuitInfos) []InterconnectCircuitInfos {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InterconnectCircuitInfos
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInterconnectCircuitInfos(c, &d, &n))
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
func diffInterconnect(c *Client, desired, actual *Interconnect, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SelfLink, actual.SelfLink, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SelfLink")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.Info{OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LinkType, actual.LinkType, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("LinkType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RequestedLinkCount, actual.RequestedLinkCount, dcl.Info{OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("RequestedLinkCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InterconnectType, actual.InterconnectType, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("InterconnectType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdminEnabled, actual.AdminEnabled, dcl.Info{OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("AdminEnabled")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NocContactEmail, actual.NocContactEmail, dcl.Info{OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("NocContactEmail")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CustomerName, actual.CustomerName, dcl.Info{OperationSelector: dcl.TriggersOperation("updateInterconnectPatchOperation")}, fn.AddNest("CustomerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OperationalStatus, actual.OperationalStatus, dcl.Info{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OperationalStatus")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ProvisionedLinkCount, actual.ProvisionedLinkCount, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ProvisionedLinkCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InterconnectAttachments, actual.InterconnectAttachments, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InterconnectAttachments")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PeerIPAddress, actual.PeerIPAddress, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeerIpAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GoogleIPAddress, actual.GoogleIPAddress, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GoogleIpAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GoogleReferenceId, actual.GoogleReferenceId, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GoogleReferenceId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExpectedOutages, actual.ExpectedOutages, dcl.Info{OutputOnly: true, ObjectFunction: compareInterconnectExpectedOutagesNewStyle, EmptyObject: EmptyInterconnectExpectedOutages, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExpectedOutages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CircuitInfos, actual.CircuitInfos, dcl.Info{OutputOnly: true, ObjectFunction: compareInterconnectCircuitInfosNewStyle, EmptyObject: EmptyInterconnectCircuitInfos, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CircuitInfos")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.Info{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Project")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	return newDiffs, nil
}
func compareInterconnectExpectedOutagesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InterconnectExpectedOutages)
	if !ok {
		desiredNotPointer, ok := d.(InterconnectExpectedOutages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InterconnectExpectedOutages or *InterconnectExpectedOutages", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InterconnectExpectedOutages)
	if !ok {
		actualNotPointer, ok := a.(InterconnectExpectedOutages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InterconnectExpectedOutages", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IssueType, actual.IssueType, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IssueType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AffectedCircuits, actual.AffectedCircuits, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AffectedCircuits")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StartTime, actual.StartTime, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EndTime, actual.EndTime, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EndTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInterconnectCircuitInfosNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InterconnectCircuitInfos)
	if !ok {
		desiredNotPointer, ok := d.(InterconnectCircuitInfos)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InterconnectCircuitInfos or *InterconnectCircuitInfos", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InterconnectCircuitInfos)
	if !ok {
		actualNotPointer, ok := a.(InterconnectCircuitInfos)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InterconnectCircuitInfos", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GoogleCircuitId, actual.GoogleCircuitId, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GoogleCircuitId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GoogleDemarcId, actual.GoogleDemarcId, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GoogleDemarcId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CustomerDemarcId, actual.CustomerDemarcId, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CustomerDemarcId")); len(ds) != 0 || err != nil {
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
func (r *Interconnect) urlNormalized() *Interconnect {
	normalized := dcl.Copy(*r).(Interconnect)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.NocContactEmail = dcl.SelfLinkToName(r.NocContactEmail)
	normalized.CustomerName = dcl.SelfLinkToName(r.CustomerName)
	normalized.PeerIPAddress = dcl.SelfLinkToName(r.PeerIPAddress)
	normalized.GoogleIPAddress = dcl.SelfLinkToName(r.GoogleIPAddress)
	normalized.GoogleReferenceId = dcl.SelfLinkToName(r.GoogleReferenceId)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	return &normalized
}

func (r *Interconnect) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "Patch" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/global/interconnects/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Interconnect resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Interconnect) marshal(c *Client) ([]byte, error) {
	m, err := expandInterconnect(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Interconnect: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalInterconnect decodes JSON responses into the Interconnect resource schema.
func unmarshalInterconnect(b []byte, c *Client) (*Interconnect, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapInterconnect(m, c)
}

func unmarshalMapInterconnect(m map[string]interface{}, c *Client) (*Interconnect, error) {

	flattened := flattenInterconnect(c, m)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandInterconnect expands Interconnect into a JSON request object.
func expandInterconnect(c *Client, f *Interconnect) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v := f.Location; dcl.ValueShouldBeSent(v) {
		m["location"] = v
	}
	if v := f.LinkType; dcl.ValueShouldBeSent(v) {
		m["linkType"] = v
	}
	if v := f.RequestedLinkCount; dcl.ValueShouldBeSent(v) {
		m["requestedLinkCount"] = v
	}
	if v := f.InterconnectType; dcl.ValueShouldBeSent(v) {
		m["interconnectType"] = v
	}
	if v := f.AdminEnabled; dcl.ValueShouldBeSent(v) {
		m["adminEnabled"] = v
	}
	if v := f.NocContactEmail; dcl.ValueShouldBeSent(v) {
		m["nocContactEmail"] = v
	}
	if v := f.CustomerName; dcl.ValueShouldBeSent(v) {
		m["customerName"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if v != nil {
		m["project"] = v
	}

	return m, nil
}

// flattenInterconnect flattens Interconnect from a JSON request object into the
// Interconnect type.
func flattenInterconnect(c *Client, i interface{}) *Interconnect {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	res := &Interconnect{}
	res.Description = dcl.FlattenString(m["description"])
	res.SelfLink = dcl.FlattenString(m["selfLink"])
	res.Id = dcl.FlattenInteger(m["id"])
	res.Name = dcl.FlattenString(m["name"])
	res.Location = dcl.FlattenString(m["location"])
	res.LinkType = flattenInterconnectLinkTypeEnum(m["linkType"])
	res.RequestedLinkCount = dcl.FlattenInteger(m["requestedLinkCount"])
	res.InterconnectType = flattenInterconnectInterconnectTypeEnum(m["interconnectType"])
	res.AdminEnabled = dcl.FlattenBool(m["adminEnabled"])
	res.NocContactEmail = dcl.FlattenString(m["nocContactEmail"])
	res.CustomerName = dcl.FlattenString(m["customerName"])
	res.OperationalStatus = flattenInterconnectOperationalStatusEnum(m["operationalStatus"])
	res.ProvisionedLinkCount = dcl.FlattenInteger(m["provisionedLinkCount"])
	res.InterconnectAttachments = dcl.FlattenStringSlice(m["interconnectAttachments"])
	res.PeerIPAddress = dcl.FlattenString(m["peerIpAddress"])
	res.GoogleIPAddress = dcl.FlattenString(m["googleIpAddress"])
	res.GoogleReferenceId = dcl.FlattenString(m["googleReferenceId"])
	res.ExpectedOutages = flattenInterconnectExpectedOutagesSlice(c, m["expectedOutages"])
	res.CircuitInfos = flattenInterconnectCircuitInfosSlice(c, m["circuitInfos"])
	res.State = flattenInterconnectStateEnum(m["state"])
	res.Project = dcl.FlattenString(m["project"])

	return res
}

// expandInterconnectExpectedOutagesMap expands the contents of InterconnectExpectedOutages into a JSON
// request object.
func expandInterconnectExpectedOutagesMap(c *Client, f map[string]InterconnectExpectedOutages) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInterconnectExpectedOutages(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInterconnectExpectedOutagesSlice expands the contents of InterconnectExpectedOutages into a JSON
// request object.
func expandInterconnectExpectedOutagesSlice(c *Client, f []InterconnectExpectedOutages) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInterconnectExpectedOutages(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInterconnectExpectedOutagesMap flattens the contents of InterconnectExpectedOutages from a JSON
// response object.
func flattenInterconnectExpectedOutagesMap(c *Client, i interface{}) map[string]InterconnectExpectedOutages {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectExpectedOutages{}
	}

	if len(a) == 0 {
		return map[string]InterconnectExpectedOutages{}
	}

	items := make(map[string]InterconnectExpectedOutages)
	for k, item := range a {
		items[k] = *flattenInterconnectExpectedOutages(c, item.(map[string]interface{}))
	}

	return items
}

// flattenInterconnectExpectedOutagesSlice flattens the contents of InterconnectExpectedOutages from a JSON
// response object.
func flattenInterconnectExpectedOutagesSlice(c *Client, i interface{}) []InterconnectExpectedOutages {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectExpectedOutages{}
	}

	if len(a) == 0 {
		return []InterconnectExpectedOutages{}
	}

	items := make([]InterconnectExpectedOutages, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectExpectedOutages(c, item.(map[string]interface{})))
	}

	return items
}

// expandInterconnectExpectedOutages expands an instance of InterconnectExpectedOutages into a JSON
// request object.
func expandInterconnectExpectedOutages(c *Client, f *InterconnectExpectedOutages) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}
	if v := f.Source; !dcl.IsEmptyValueIndirect(v) {
		m["source"] = v
	}
	if v := f.State; !dcl.IsEmptyValueIndirect(v) {
		m["state"] = v
	}
	if v := f.IssueType; !dcl.IsEmptyValueIndirect(v) {
		m["issueType"] = v
	}
	if v := f.AffectedCircuits; v != nil {
		m["affectedCircuits"] = v
	}
	if v := f.StartTime; !dcl.IsEmptyValueIndirect(v) {
		m["startTime"] = v
	}
	if v := f.EndTime; !dcl.IsEmptyValueIndirect(v) {
		m["endTime"] = v
	}

	return m, nil
}

// flattenInterconnectExpectedOutages flattens an instance of InterconnectExpectedOutages from a JSON
// response object.
func flattenInterconnectExpectedOutages(c *Client, i interface{}) *InterconnectExpectedOutages {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InterconnectExpectedOutages{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInterconnectExpectedOutages
	}
	r.Name = dcl.FlattenString(m["name"])
	r.Description = dcl.FlattenString(m["description"])
	r.Source = flattenInterconnectExpectedOutagesSourceEnum(m["source"])
	r.State = flattenInterconnectExpectedOutagesStateEnum(m["state"])
	r.IssueType = flattenInterconnectExpectedOutagesIssueTypeEnum(m["issueType"])
	r.AffectedCircuits = dcl.FlattenStringSlice(m["affectedCircuits"])
	r.StartTime = dcl.FlattenInteger(m["startTime"])
	r.EndTime = dcl.FlattenInteger(m["endTime"])

	return r
}

// expandInterconnectCircuitInfosMap expands the contents of InterconnectCircuitInfos into a JSON
// request object.
func expandInterconnectCircuitInfosMap(c *Client, f map[string]InterconnectCircuitInfos) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInterconnectCircuitInfos(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInterconnectCircuitInfosSlice expands the contents of InterconnectCircuitInfos into a JSON
// request object.
func expandInterconnectCircuitInfosSlice(c *Client, f []InterconnectCircuitInfos) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInterconnectCircuitInfos(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInterconnectCircuitInfosMap flattens the contents of InterconnectCircuitInfos from a JSON
// response object.
func flattenInterconnectCircuitInfosMap(c *Client, i interface{}) map[string]InterconnectCircuitInfos {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectCircuitInfos{}
	}

	if len(a) == 0 {
		return map[string]InterconnectCircuitInfos{}
	}

	items := make(map[string]InterconnectCircuitInfos)
	for k, item := range a {
		items[k] = *flattenInterconnectCircuitInfos(c, item.(map[string]interface{}))
	}

	return items
}

// flattenInterconnectCircuitInfosSlice flattens the contents of InterconnectCircuitInfos from a JSON
// response object.
func flattenInterconnectCircuitInfosSlice(c *Client, i interface{}) []InterconnectCircuitInfos {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectCircuitInfos{}
	}

	if len(a) == 0 {
		return []InterconnectCircuitInfos{}
	}

	items := make([]InterconnectCircuitInfos, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectCircuitInfos(c, item.(map[string]interface{})))
	}

	return items
}

// expandInterconnectCircuitInfos expands an instance of InterconnectCircuitInfos into a JSON
// request object.
func expandInterconnectCircuitInfos(c *Client, f *InterconnectCircuitInfos) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GoogleCircuitId; !dcl.IsEmptyValueIndirect(v) {
		m["googleCircuitId"] = v
	}
	if v := f.GoogleDemarcId; !dcl.IsEmptyValueIndirect(v) {
		m["googleDemarcId"] = v
	}
	if v := f.CustomerDemarcId; !dcl.IsEmptyValueIndirect(v) {
		m["customerDemarcId"] = v
	}

	return m, nil
}

// flattenInterconnectCircuitInfos flattens an instance of InterconnectCircuitInfos from a JSON
// response object.
func flattenInterconnectCircuitInfos(c *Client, i interface{}) *InterconnectCircuitInfos {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InterconnectCircuitInfos{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInterconnectCircuitInfos
	}
	r.GoogleCircuitId = dcl.FlattenString(m["googleCircuitId"])
	r.GoogleDemarcId = dcl.FlattenString(m["googleDemarcId"])
	r.CustomerDemarcId = dcl.FlattenString(m["customerDemarcId"])

	return r
}

// flattenInterconnectLinkTypeEnumMap flattens the contents of InterconnectLinkTypeEnum from a JSON
// response object.
func flattenInterconnectLinkTypeEnumMap(c *Client, i interface{}) map[string]InterconnectLinkTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectLinkTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectLinkTypeEnum{}
	}

	items := make(map[string]InterconnectLinkTypeEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectLinkTypeEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectLinkTypeEnumSlice flattens the contents of InterconnectLinkTypeEnum from a JSON
// response object.
func flattenInterconnectLinkTypeEnumSlice(c *Client, i interface{}) []InterconnectLinkTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectLinkTypeEnum{}
	}

	if len(a) == 0 {
		return []InterconnectLinkTypeEnum{}
	}

	items := make([]InterconnectLinkTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectLinkTypeEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectLinkTypeEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectLinkTypeEnum with the same value as that string.
func flattenInterconnectLinkTypeEnum(i interface{}) *InterconnectLinkTypeEnum {
	s, ok := i.(string)
	if !ok {
		return InterconnectLinkTypeEnumRef("")
	}

	return InterconnectLinkTypeEnumRef(s)
}

// flattenInterconnectInterconnectTypeEnumMap flattens the contents of InterconnectInterconnectTypeEnum from a JSON
// response object.
func flattenInterconnectInterconnectTypeEnumMap(c *Client, i interface{}) map[string]InterconnectInterconnectTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectInterconnectTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectInterconnectTypeEnum{}
	}

	items := make(map[string]InterconnectInterconnectTypeEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectInterconnectTypeEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectInterconnectTypeEnumSlice flattens the contents of InterconnectInterconnectTypeEnum from a JSON
// response object.
func flattenInterconnectInterconnectTypeEnumSlice(c *Client, i interface{}) []InterconnectInterconnectTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectInterconnectTypeEnum{}
	}

	if len(a) == 0 {
		return []InterconnectInterconnectTypeEnum{}
	}

	items := make([]InterconnectInterconnectTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectInterconnectTypeEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectInterconnectTypeEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectInterconnectTypeEnum with the same value as that string.
func flattenInterconnectInterconnectTypeEnum(i interface{}) *InterconnectInterconnectTypeEnum {
	s, ok := i.(string)
	if !ok {
		return InterconnectInterconnectTypeEnumRef("")
	}

	return InterconnectInterconnectTypeEnumRef(s)
}

// flattenInterconnectOperationalStatusEnumMap flattens the contents of InterconnectOperationalStatusEnum from a JSON
// response object.
func flattenInterconnectOperationalStatusEnumMap(c *Client, i interface{}) map[string]InterconnectOperationalStatusEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectOperationalStatusEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectOperationalStatusEnum{}
	}

	items := make(map[string]InterconnectOperationalStatusEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectOperationalStatusEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectOperationalStatusEnumSlice flattens the contents of InterconnectOperationalStatusEnum from a JSON
// response object.
func flattenInterconnectOperationalStatusEnumSlice(c *Client, i interface{}) []InterconnectOperationalStatusEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectOperationalStatusEnum{}
	}

	if len(a) == 0 {
		return []InterconnectOperationalStatusEnum{}
	}

	items := make([]InterconnectOperationalStatusEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectOperationalStatusEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectOperationalStatusEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectOperationalStatusEnum with the same value as that string.
func flattenInterconnectOperationalStatusEnum(i interface{}) *InterconnectOperationalStatusEnum {
	s, ok := i.(string)
	if !ok {
		return InterconnectOperationalStatusEnumRef("")
	}

	return InterconnectOperationalStatusEnumRef(s)
}

// flattenInterconnectExpectedOutagesSourceEnumMap flattens the contents of InterconnectExpectedOutagesSourceEnum from a JSON
// response object.
func flattenInterconnectExpectedOutagesSourceEnumMap(c *Client, i interface{}) map[string]InterconnectExpectedOutagesSourceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectExpectedOutagesSourceEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectExpectedOutagesSourceEnum{}
	}

	items := make(map[string]InterconnectExpectedOutagesSourceEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectExpectedOutagesSourceEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectExpectedOutagesSourceEnumSlice flattens the contents of InterconnectExpectedOutagesSourceEnum from a JSON
// response object.
func flattenInterconnectExpectedOutagesSourceEnumSlice(c *Client, i interface{}) []InterconnectExpectedOutagesSourceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectExpectedOutagesSourceEnum{}
	}

	if len(a) == 0 {
		return []InterconnectExpectedOutagesSourceEnum{}
	}

	items := make([]InterconnectExpectedOutagesSourceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectExpectedOutagesSourceEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectExpectedOutagesSourceEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectExpectedOutagesSourceEnum with the same value as that string.
func flattenInterconnectExpectedOutagesSourceEnum(i interface{}) *InterconnectExpectedOutagesSourceEnum {
	s, ok := i.(string)
	if !ok {
		return InterconnectExpectedOutagesSourceEnumRef("")
	}

	return InterconnectExpectedOutagesSourceEnumRef(s)
}

// flattenInterconnectExpectedOutagesStateEnumMap flattens the contents of InterconnectExpectedOutagesStateEnum from a JSON
// response object.
func flattenInterconnectExpectedOutagesStateEnumMap(c *Client, i interface{}) map[string]InterconnectExpectedOutagesStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectExpectedOutagesStateEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectExpectedOutagesStateEnum{}
	}

	items := make(map[string]InterconnectExpectedOutagesStateEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectExpectedOutagesStateEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectExpectedOutagesStateEnumSlice flattens the contents of InterconnectExpectedOutagesStateEnum from a JSON
// response object.
func flattenInterconnectExpectedOutagesStateEnumSlice(c *Client, i interface{}) []InterconnectExpectedOutagesStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectExpectedOutagesStateEnum{}
	}

	if len(a) == 0 {
		return []InterconnectExpectedOutagesStateEnum{}
	}

	items := make([]InterconnectExpectedOutagesStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectExpectedOutagesStateEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectExpectedOutagesStateEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectExpectedOutagesStateEnum with the same value as that string.
func flattenInterconnectExpectedOutagesStateEnum(i interface{}) *InterconnectExpectedOutagesStateEnum {
	s, ok := i.(string)
	if !ok {
		return InterconnectExpectedOutagesStateEnumRef("")
	}

	return InterconnectExpectedOutagesStateEnumRef(s)
}

// flattenInterconnectExpectedOutagesIssueTypeEnumMap flattens the contents of InterconnectExpectedOutagesIssueTypeEnum from a JSON
// response object.
func flattenInterconnectExpectedOutagesIssueTypeEnumMap(c *Client, i interface{}) map[string]InterconnectExpectedOutagesIssueTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectExpectedOutagesIssueTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectExpectedOutagesIssueTypeEnum{}
	}

	items := make(map[string]InterconnectExpectedOutagesIssueTypeEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectExpectedOutagesIssueTypeEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectExpectedOutagesIssueTypeEnumSlice flattens the contents of InterconnectExpectedOutagesIssueTypeEnum from a JSON
// response object.
func flattenInterconnectExpectedOutagesIssueTypeEnumSlice(c *Client, i interface{}) []InterconnectExpectedOutagesIssueTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectExpectedOutagesIssueTypeEnum{}
	}

	if len(a) == 0 {
		return []InterconnectExpectedOutagesIssueTypeEnum{}
	}

	items := make([]InterconnectExpectedOutagesIssueTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectExpectedOutagesIssueTypeEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectExpectedOutagesIssueTypeEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectExpectedOutagesIssueTypeEnum with the same value as that string.
func flattenInterconnectExpectedOutagesIssueTypeEnum(i interface{}) *InterconnectExpectedOutagesIssueTypeEnum {
	s, ok := i.(string)
	if !ok {
		return InterconnectExpectedOutagesIssueTypeEnumRef("")
	}

	return InterconnectExpectedOutagesIssueTypeEnumRef(s)
}

// flattenInterconnectStateEnumMap flattens the contents of InterconnectStateEnum from a JSON
// response object.
func flattenInterconnectStateEnumMap(c *Client, i interface{}) map[string]InterconnectStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InterconnectStateEnum{}
	}

	if len(a) == 0 {
		return map[string]InterconnectStateEnum{}
	}

	items := make(map[string]InterconnectStateEnum)
	for k, item := range a {
		items[k] = *flattenInterconnectStateEnum(item.(interface{}))
	}

	return items
}

// flattenInterconnectStateEnumSlice flattens the contents of InterconnectStateEnum from a JSON
// response object.
func flattenInterconnectStateEnumSlice(c *Client, i interface{}) []InterconnectStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InterconnectStateEnum{}
	}

	if len(a) == 0 {
		return []InterconnectStateEnum{}
	}

	items := make([]InterconnectStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInterconnectStateEnum(item.(interface{})))
	}

	return items
}

// flattenInterconnectStateEnum asserts that an interface is a string, and returns a
// pointer to a *InterconnectStateEnum with the same value as that string.
func flattenInterconnectStateEnum(i interface{}) *InterconnectStateEnum {
	s, ok := i.(string)
	if !ok {
		return InterconnectStateEnumRef("")
	}

	return InterconnectStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Interconnect) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalInterconnect(b, c)
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

type interconnectDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         interconnectApiOperation
}

func convertFieldDiffsToInterconnectDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]interconnectDiff, error) {
	opNamesToFieldDiffs := make(map[string][]*dcl.FieldDiff)
	// Map each operation name to the field diffs associated with it.
	for _, fd := range fds {
		for _, ro := range fd.ResultingOperation {
			if fieldDiffs, ok := opNamesToFieldDiffs[ro]; ok {
				fieldDiffs = append(fieldDiffs, fd)
				opNamesToFieldDiffs[ro] = fieldDiffs
			} else {
				config.Logger.Infof("%s required due to diff in %q", ro, fd.FieldName)
				opNamesToFieldDiffs[ro] = []*dcl.FieldDiff{fd}
			}
		}
	}
	var diffs []interconnectDiff
	// For each operation name, create a interconnectDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		diff := interconnectDiff{}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToInterconnectApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToInterconnectApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (interconnectApiOperation, error) {
	switch opName {

	case "updateInterconnectPatchOperation":
		return &updateInterconnectPatchOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractInterconnectFields(r *Interconnect) error {
	return nil
}
func extractInterconnectExpectedOutagesFields(r *Interconnect, o *InterconnectExpectedOutages) error {
	return nil
}
func extractInterconnectCircuitInfosFields(r *Interconnect, o *InterconnectCircuitInfos) error {
	return nil
}

func postReadExtractInterconnectFields(r *Interconnect) error {
	return nil
}
func postReadExtractInterconnectExpectedOutagesFields(r *Interconnect, o *InterconnectExpectedOutages) error {
	return nil
}
func postReadExtractInterconnectCircuitInfosFields(r *Interconnect, o *InterconnectCircuitInfos) error {
	return nil
}
