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
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *Trigger) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "matchingCriteria"); err != nil {
		return err
	}
	if err := dcl.Required(r, "destination"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Destination) {
		if err := r.Destination.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Transport) {
		if err := r.Transport.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *TriggerMatchingCriteria) validate() error {
	if err := dcl.Required(r, "attribute"); err != nil {
		return err
	}
	if err := dcl.Required(r, "value"); err != nil {
		return err
	}
	return nil
}
func (r *TriggerDestination) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"CloudRunService", "CloudFunction", "Gke", "Workflow"}, r.CloudRunService, r.CloudFunction, r.Gke, r.Workflow); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.CloudRunService) {
		if err := r.CloudRunService.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Gke) {
		if err := r.Gke.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.HttpEndpoint) {
		if err := r.HttpEndpoint.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.NetworkConfig) {
		if err := r.NetworkConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *TriggerDestinationCloudRunService) validate() error {
	if err := dcl.Required(r, "service"); err != nil {
		return err
	}
	if err := dcl.Required(r, "region"); err != nil {
		return err
	}
	return nil
}
func (r *TriggerDestinationGke) validate() error {
	if err := dcl.Required(r, "cluster"); err != nil {
		return err
	}
	if err := dcl.Required(r, "location"); err != nil {
		return err
	}
	if err := dcl.Required(r, "namespace"); err != nil {
		return err
	}
	if err := dcl.Required(r, "service"); err != nil {
		return err
	}
	return nil
}
func (r *TriggerDestinationHttpEndpoint) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *TriggerDestinationNetworkConfig) validate() error {
	if err := dcl.Required(r, "networkAttachment"); err != nil {
		return err
	}
	return nil
}
func (r *TriggerTransport) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Pubsub) {
		if err := r.Pubsub.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *TriggerTransportPubsub) validate() error {
	return nil
}
func (r *Trigger) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://eventarc.googleapis.com/v1/", params)
}

func (r *Trigger) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/triggers/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Trigger) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/triggers", nr.basePath(), userBasePath, params), nil

}

func (r *Trigger) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/triggers?triggerId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Trigger) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/triggers/{{name}}", nr.basePath(), userBasePath, params), nil
}

// triggerApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type triggerApiOperation interface {
	do(context.Context, *Trigger, *Client) error
}

// newUpdateTriggerUpdateTriggerRequest creates a request for an
// Trigger resource's UpdateTrigger update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateTriggerUpdateTriggerRequest(ctx context.Context, f *Trigger, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := expandTriggerMatchingCriteriaSlice(c, f.MatchingCriteria, res); err != nil {
		return nil, fmt.Errorf("error expanding MatchingCriteria into eventFilters: %w", err)
	} else if v != nil {
		req["eventFilters"] = v
	}
	if v := f.ServiceAccount; !dcl.IsEmptyValueIndirect(v) {
		req["serviceAccount"] = v
	}
	if v, err := expandTriggerDestination(c, f.Destination, res); err != nil {
		return nil, fmt.Errorf("error expanding Destination into destination: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["destination"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	if v := f.EventDataContentType; !dcl.IsEmptyValueIndirect(v) {
		req["eventDataContentType"] = v
	}
	b, err := c.getTriggerRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawEtag, err := dcl.GetMapEntry(
		m,
		[]string{"etag"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["etag"] = rawEtag.(string)
	}
	return req, nil
}

// marshalUpdateTriggerUpdateTriggerRequest converts the update into
// the final JSON request body.
func marshalUpdateTriggerUpdateTriggerRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateTriggerUpdateTriggerOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateTriggerUpdateTriggerOperation) do(ctx context.Context, r *Trigger, c *Client) error {
	_, err := c.GetTrigger(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateTrigger")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateTriggerUpdateTriggerRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateTriggerUpdateTriggerRequest(c, req)
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

func (c *Client) listTriggerRaw(ctx context.Context, r *Trigger, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != TriggerMaxPage {
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

type listTriggerOperation struct {
	Triggers []map[string]interface{} `json:"triggers"`
	Token    string                   `json:"nextPageToken"`
}

func (c *Client) listTrigger(ctx context.Context, r *Trigger, pageToken string, pageSize int32) ([]*Trigger, string, error) {
	b, err := c.listTriggerRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listTriggerOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Trigger
	for _, v := range m.Triggers {
		res, err := unmarshalMapTrigger(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllTrigger(ctx context.Context, f func(*Trigger) bool, resources []*Trigger) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteTrigger(ctx, res)
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

type deleteTriggerOperation struct{}

func (op *deleteTriggerOperation) do(ctx context.Context, r *Trigger, c *Client) error {
	r, err := c.GetTrigger(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Trigger not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetTrigger checking for existence. error: %v", err)
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
		_, err := c.GetTrigger(ctx, r)
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
type createTriggerOperation struct {
	response map[string]interface{}
}

func (op *createTriggerOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createTriggerOperation) do(ctx context.Context, r *Trigger, c *Client) error {
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

	if _, err := c.GetTrigger(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getTriggerRaw(ctx context.Context, r *Trigger) ([]byte, error) {

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

func (c *Client) triggerDiffsForRawDesired(ctx context.Context, rawDesired *Trigger, opts ...dcl.ApplyOption) (initial, desired *Trigger, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Trigger
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Trigger); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Trigger, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetTrigger(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Trigger resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Trigger resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Trigger resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeTriggerDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Trigger: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Trigger: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractTriggerFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeTriggerInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Trigger: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeTriggerDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Trigger: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffTrigger(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeTriggerInitialState(rawInitial, rawDesired *Trigger) (*Trigger, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeTriggerDesiredState(rawDesired, rawInitial *Trigger, opts ...dcl.ApplyOption) (*Trigger, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Destination = canonicalizeTriggerDestination(rawDesired.Destination, nil, opts...)
		rawDesired.Transport = canonicalizeTriggerTransport(rawDesired.Transport, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Trigger{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	canonicalDesired.MatchingCriteria = canonicalizeTriggerMatchingCriteriaSlice(rawDesired.MatchingCriteria, rawInitial.MatchingCriteria, opts...)
	if dcl.IsZeroValue(rawDesired.ServiceAccount) || (dcl.IsEmptyValueIndirect(rawDesired.ServiceAccount) && dcl.IsEmptyValueIndirect(rawInitial.ServiceAccount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.ServiceAccount = rawInitial.ServiceAccount
	} else {
		canonicalDesired.ServiceAccount = rawDesired.ServiceAccount
	}
	canonicalDesired.Destination = canonicalizeTriggerDestination(rawDesired.Destination, rawInitial.Destination, opts...)
	canonicalDesired.Transport = canonicalizeTriggerTransport(rawDesired.Transport, rawInitial.Transport, opts...)
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
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
	if dcl.IsZeroValue(rawDesired.Channel) || (dcl.IsEmptyValueIndirect(rawDesired.Channel) && dcl.IsEmptyValueIndirect(rawInitial.Channel)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Channel = rawInitial.Channel
	} else {
		canonicalDesired.Channel = rawDesired.Channel
	}
	if dcl.StringCanonicalize(rawDesired.EventDataContentType, rawInitial.EventDataContentType) {
		canonicalDesired.EventDataContentType = rawInitial.EventDataContentType
	} else {
		canonicalDesired.EventDataContentType = rawDesired.EventDataContentType
	}
	return canonicalDesired, nil
}

func canonicalizeTriggerNewState(c *Client, rawNew, rawDesired *Trigger) (*Trigger, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Uid) && dcl.IsEmptyValueIndirect(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
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

	if dcl.IsEmptyValueIndirect(rawNew.MatchingCriteria) && dcl.IsEmptyValueIndirect(rawDesired.MatchingCriteria) {
		rawNew.MatchingCriteria = rawDesired.MatchingCriteria
	} else {
		rawNew.MatchingCriteria = canonicalizeNewTriggerMatchingCriteriaSet(c, rawDesired.MatchingCriteria, rawNew.MatchingCriteria)
	}

	if dcl.IsEmptyValueIndirect(rawNew.ServiceAccount) && dcl.IsEmptyValueIndirect(rawDesired.ServiceAccount) {
		rawNew.ServiceAccount = rawDesired.ServiceAccount
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Destination) && dcl.IsEmptyValueIndirect(rawDesired.Destination) {
		rawNew.Destination = rawDesired.Destination
	} else {
		rawNew.Destination = canonicalizeNewTriggerDestination(c, rawDesired.Destination, rawNew.Destination)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Transport) && dcl.IsEmptyValueIndirect(rawDesired.Transport) {
		rawNew.Transport = rawDesired.Transport
	} else {
		rawNew.Transport = canonicalizeNewTriggerTransport(c, rawDesired.Transport, rawNew.Transport)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Etag) && dcl.IsEmptyValueIndirect(rawDesired.Etag) {
		rawNew.Etag = rawDesired.Etag
	} else {
		if dcl.StringCanonicalize(rawDesired.Etag, rawNew.Etag) {
			rawNew.Etag = rawDesired.Etag
		}
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	if dcl.IsEmptyValueIndirect(rawNew.Channel) && dcl.IsEmptyValueIndirect(rawDesired.Channel) {
		rawNew.Channel = rawDesired.Channel
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Conditions) && dcl.IsEmptyValueIndirect(rawDesired.Conditions) {
		rawNew.Conditions = rawDesired.Conditions
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.EventDataContentType) && dcl.IsEmptyValueIndirect(rawDesired.EventDataContentType) {
		rawNew.EventDataContentType = rawDesired.EventDataContentType
	} else {
		if dcl.StringCanonicalize(rawDesired.EventDataContentType, rawNew.EventDataContentType) {
			rawNew.EventDataContentType = rawDesired.EventDataContentType
		}
	}

	return rawNew, nil
}

func canonicalizeTriggerMatchingCriteria(des, initial *TriggerMatchingCriteria, opts ...dcl.ApplyOption) *TriggerMatchingCriteria {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &TriggerMatchingCriteria{}

	if dcl.StringCanonicalize(des.Attribute, initial.Attribute) || dcl.IsZeroValue(des.Attribute) {
		cDes.Attribute = initial.Attribute
	} else {
		cDes.Attribute = des.Attribute
	}
	if dcl.StringCanonicalize(des.Value, initial.Value) || dcl.IsZeroValue(des.Value) {
		cDes.Value = initial.Value
	} else {
		cDes.Value = des.Value
	}
	if dcl.StringCanonicalize(des.Operator, initial.Operator) || dcl.IsZeroValue(des.Operator) {
		cDes.Operator = initial.Operator
	} else {
		cDes.Operator = des.Operator
	}

	return cDes
}

func canonicalizeTriggerMatchingCriteriaSlice(des, initial []TriggerMatchingCriteria, opts ...dcl.ApplyOption) []TriggerMatchingCriteria {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]TriggerMatchingCriteria, 0, len(des))
		for _, d := range des {
			cd := canonicalizeTriggerMatchingCriteria(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]TriggerMatchingCriteria, 0, len(des))
	for i, d := range des {
		cd := canonicalizeTriggerMatchingCriteria(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewTriggerMatchingCriteria(c *Client, des, nw *TriggerMatchingCriteria) *TriggerMatchingCriteria {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for TriggerMatchingCriteria while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Attribute, nw.Attribute) {
		nw.Attribute = des.Attribute
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}
	if dcl.StringCanonicalize(des.Operator, nw.Operator) {
		nw.Operator = des.Operator
	}

	return nw
}

func canonicalizeNewTriggerMatchingCriteriaSet(c *Client, des, nw []TriggerMatchingCriteria) []TriggerMatchingCriteria {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []TriggerMatchingCriteria
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareTriggerMatchingCriteriaNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewTriggerMatchingCriteria(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewTriggerMatchingCriteriaSlice(c *Client, des, nw []TriggerMatchingCriteria) []TriggerMatchingCriteria {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []TriggerMatchingCriteria
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewTriggerMatchingCriteria(c, &d, &n))
	}

	return items
}

func canonicalizeTriggerDestination(des, initial *TriggerDestination, opts ...dcl.ApplyOption) *TriggerDestination {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.CloudRunService != nil || (initial != nil && initial.CloudRunService != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.CloudFunction, des.Gke, des.Workflow) {
			des.CloudRunService = nil
			if initial != nil {
				initial.CloudRunService = nil
			}
		}
	}

	if des.CloudFunction != nil || (initial != nil && initial.CloudFunction != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.CloudRunService, des.Gke, des.Workflow) {
			des.CloudFunction = nil
			if initial != nil {
				initial.CloudFunction = nil
			}
		}
	}

	if des.Gke != nil || (initial != nil && initial.Gke != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.CloudRunService, des.CloudFunction, des.Workflow) {
			des.Gke = nil
			if initial != nil {
				initial.Gke = nil
			}
		}
	}

	if des.Workflow != nil || (initial != nil && initial.Workflow != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.CloudRunService, des.CloudFunction, des.Gke) {
			des.Workflow = nil
			if initial != nil {
				initial.Workflow = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &TriggerDestination{}

	cDes.CloudRunService = canonicalizeTriggerDestinationCloudRunService(des.CloudRunService, initial.CloudRunService, opts...)
	if dcl.PartialSelfLinkToSelfLink(des.CloudFunction, initial.CloudFunction) || dcl.IsZeroValue(des.CloudFunction) {
		cDes.CloudFunction = initial.CloudFunction
	} else {
		cDes.CloudFunction = des.CloudFunction
	}
	cDes.Gke = canonicalizeTriggerDestinationGke(des.Gke, initial.Gke, opts...)
	if dcl.PartialSelfLinkToSelfLink(des.Workflow, initial.Workflow) || dcl.IsZeroValue(des.Workflow) {
		cDes.Workflow = initial.Workflow
	} else {
		cDes.Workflow = des.Workflow
	}
	cDes.HttpEndpoint = canonicalizeTriggerDestinationHttpEndpoint(des.HttpEndpoint, initial.HttpEndpoint, opts...)
	cDes.NetworkConfig = canonicalizeTriggerDestinationNetworkConfig(des.NetworkConfig, initial.NetworkConfig, opts...)

	return cDes
}

func canonicalizeTriggerDestinationSlice(des, initial []TriggerDestination, opts ...dcl.ApplyOption) []TriggerDestination {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]TriggerDestination, 0, len(des))
		for _, d := range des {
			cd := canonicalizeTriggerDestination(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]TriggerDestination, 0, len(des))
	for i, d := range des {
		cd := canonicalizeTriggerDestination(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewTriggerDestination(c *Client, des, nw *TriggerDestination) *TriggerDestination {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for TriggerDestination while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.CloudRunService = canonicalizeNewTriggerDestinationCloudRunService(c, des.CloudRunService, nw.CloudRunService)
	if dcl.PartialSelfLinkToSelfLink(des.CloudFunction, nw.CloudFunction) {
		nw.CloudFunction = des.CloudFunction
	}
	nw.Gke = canonicalizeNewTriggerDestinationGke(c, des.Gke, nw.Gke)
	if dcl.PartialSelfLinkToSelfLink(des.Workflow, nw.Workflow) {
		nw.Workflow = des.Workflow
	}
	nw.HttpEndpoint = canonicalizeNewTriggerDestinationHttpEndpoint(c, des.HttpEndpoint, nw.HttpEndpoint)
	nw.NetworkConfig = canonicalizeNewTriggerDestinationNetworkConfig(c, des.NetworkConfig, nw.NetworkConfig)

	return nw
}

func canonicalizeNewTriggerDestinationSet(c *Client, des, nw []TriggerDestination) []TriggerDestination {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []TriggerDestination
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareTriggerDestinationNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewTriggerDestination(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewTriggerDestinationSlice(c *Client, des, nw []TriggerDestination) []TriggerDestination {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []TriggerDestination
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewTriggerDestination(c, &d, &n))
	}

	return items
}

func canonicalizeTriggerDestinationCloudRunService(des, initial *TriggerDestinationCloudRunService, opts ...dcl.ApplyOption) *TriggerDestinationCloudRunService {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &TriggerDestinationCloudRunService{}

	if dcl.IsZeroValue(des.Service) || (dcl.IsEmptyValueIndirect(des.Service) && dcl.IsEmptyValueIndirect(initial.Service)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Service = initial.Service
	} else {
		cDes.Service = des.Service
	}
	if dcl.StringCanonicalize(des.Path, initial.Path) || dcl.IsZeroValue(des.Path) {
		cDes.Path = initial.Path
	} else {
		cDes.Path = des.Path
	}
	if dcl.StringCanonicalize(des.Region, initial.Region) || dcl.IsZeroValue(des.Region) {
		cDes.Region = initial.Region
	} else {
		cDes.Region = des.Region
	}

	return cDes
}

func canonicalizeTriggerDestinationCloudRunServiceSlice(des, initial []TriggerDestinationCloudRunService, opts ...dcl.ApplyOption) []TriggerDestinationCloudRunService {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]TriggerDestinationCloudRunService, 0, len(des))
		for _, d := range des {
			cd := canonicalizeTriggerDestinationCloudRunService(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]TriggerDestinationCloudRunService, 0, len(des))
	for i, d := range des {
		cd := canonicalizeTriggerDestinationCloudRunService(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewTriggerDestinationCloudRunService(c *Client, des, nw *TriggerDestinationCloudRunService) *TriggerDestinationCloudRunService {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for TriggerDestinationCloudRunService while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Path, nw.Path) {
		nw.Path = des.Path
	}
	if dcl.StringCanonicalize(des.Region, nw.Region) {
		nw.Region = des.Region
	}

	return nw
}

func canonicalizeNewTriggerDestinationCloudRunServiceSet(c *Client, des, nw []TriggerDestinationCloudRunService) []TriggerDestinationCloudRunService {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []TriggerDestinationCloudRunService
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareTriggerDestinationCloudRunServiceNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewTriggerDestinationCloudRunService(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewTriggerDestinationCloudRunServiceSlice(c *Client, des, nw []TriggerDestinationCloudRunService) []TriggerDestinationCloudRunService {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []TriggerDestinationCloudRunService
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewTriggerDestinationCloudRunService(c, &d, &n))
	}

	return items
}

func canonicalizeTriggerDestinationGke(des, initial *TriggerDestinationGke, opts ...dcl.ApplyOption) *TriggerDestinationGke {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &TriggerDestinationGke{}

	if dcl.IsZeroValue(des.Cluster) || (dcl.IsEmptyValueIndirect(des.Cluster) && dcl.IsEmptyValueIndirect(initial.Cluster)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Cluster = initial.Cluster
	} else {
		cDes.Cluster = des.Cluster
	}
	if dcl.StringCanonicalize(des.Location, initial.Location) || dcl.IsZeroValue(des.Location) {
		cDes.Location = initial.Location
	} else {
		cDes.Location = des.Location
	}
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
	if dcl.StringCanonicalize(des.Path, initial.Path) || dcl.IsZeroValue(des.Path) {
		cDes.Path = initial.Path
	} else {
		cDes.Path = des.Path
	}

	return cDes
}

func canonicalizeTriggerDestinationGkeSlice(des, initial []TriggerDestinationGke, opts ...dcl.ApplyOption) []TriggerDestinationGke {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]TriggerDestinationGke, 0, len(des))
		for _, d := range des {
			cd := canonicalizeTriggerDestinationGke(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]TriggerDestinationGke, 0, len(des))
	for i, d := range des {
		cd := canonicalizeTriggerDestinationGke(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewTriggerDestinationGke(c *Client, des, nw *TriggerDestinationGke) *TriggerDestinationGke {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for TriggerDestinationGke while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Location, nw.Location) {
		nw.Location = des.Location
	}
	if dcl.StringCanonicalize(des.Namespace, nw.Namespace) {
		nw.Namespace = des.Namespace
	}
	if dcl.StringCanonicalize(des.Service, nw.Service) {
		nw.Service = des.Service
	}
	if dcl.StringCanonicalize(des.Path, nw.Path) {
		nw.Path = des.Path
	}

	return nw
}

func canonicalizeNewTriggerDestinationGkeSet(c *Client, des, nw []TriggerDestinationGke) []TriggerDestinationGke {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []TriggerDestinationGke
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareTriggerDestinationGkeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewTriggerDestinationGke(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewTriggerDestinationGkeSlice(c *Client, des, nw []TriggerDestinationGke) []TriggerDestinationGke {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []TriggerDestinationGke
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewTriggerDestinationGke(c, &d, &n))
	}

	return items
}

func canonicalizeTriggerDestinationHttpEndpoint(des, initial *TriggerDestinationHttpEndpoint, opts ...dcl.ApplyOption) *TriggerDestinationHttpEndpoint {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &TriggerDestinationHttpEndpoint{}

	if dcl.StringCanonicalize(des.Uri, initial.Uri) || dcl.IsZeroValue(des.Uri) {
		cDes.Uri = initial.Uri
	} else {
		cDes.Uri = des.Uri
	}

	return cDes
}

func canonicalizeTriggerDestinationHttpEndpointSlice(des, initial []TriggerDestinationHttpEndpoint, opts ...dcl.ApplyOption) []TriggerDestinationHttpEndpoint {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]TriggerDestinationHttpEndpoint, 0, len(des))
		for _, d := range des {
			cd := canonicalizeTriggerDestinationHttpEndpoint(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]TriggerDestinationHttpEndpoint, 0, len(des))
	for i, d := range des {
		cd := canonicalizeTriggerDestinationHttpEndpoint(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewTriggerDestinationHttpEndpoint(c *Client, des, nw *TriggerDestinationHttpEndpoint) *TriggerDestinationHttpEndpoint {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for TriggerDestinationHttpEndpoint while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Uri, nw.Uri) {
		nw.Uri = des.Uri
	}

	return nw
}

func canonicalizeNewTriggerDestinationHttpEndpointSet(c *Client, des, nw []TriggerDestinationHttpEndpoint) []TriggerDestinationHttpEndpoint {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []TriggerDestinationHttpEndpoint
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareTriggerDestinationHttpEndpointNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewTriggerDestinationHttpEndpoint(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewTriggerDestinationHttpEndpointSlice(c *Client, des, nw []TriggerDestinationHttpEndpoint) []TriggerDestinationHttpEndpoint {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []TriggerDestinationHttpEndpoint
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewTriggerDestinationHttpEndpoint(c, &d, &n))
	}

	return items
}

func canonicalizeTriggerDestinationNetworkConfig(des, initial *TriggerDestinationNetworkConfig, opts ...dcl.ApplyOption) *TriggerDestinationNetworkConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &TriggerDestinationNetworkConfig{}

	if dcl.IsZeroValue(des.NetworkAttachment) || (dcl.IsEmptyValueIndirect(des.NetworkAttachment) && dcl.IsEmptyValueIndirect(initial.NetworkAttachment)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NetworkAttachment = initial.NetworkAttachment
	} else {
		cDes.NetworkAttachment = des.NetworkAttachment
	}

	return cDes
}

func canonicalizeTriggerDestinationNetworkConfigSlice(des, initial []TriggerDestinationNetworkConfig, opts ...dcl.ApplyOption) []TriggerDestinationNetworkConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]TriggerDestinationNetworkConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeTriggerDestinationNetworkConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]TriggerDestinationNetworkConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeTriggerDestinationNetworkConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewTriggerDestinationNetworkConfig(c *Client, des, nw *TriggerDestinationNetworkConfig) *TriggerDestinationNetworkConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for TriggerDestinationNetworkConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewTriggerDestinationNetworkConfigSet(c *Client, des, nw []TriggerDestinationNetworkConfig) []TriggerDestinationNetworkConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []TriggerDestinationNetworkConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareTriggerDestinationNetworkConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewTriggerDestinationNetworkConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewTriggerDestinationNetworkConfigSlice(c *Client, des, nw []TriggerDestinationNetworkConfig) []TriggerDestinationNetworkConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []TriggerDestinationNetworkConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewTriggerDestinationNetworkConfig(c, &d, &n))
	}

	return items
}

func canonicalizeTriggerTransport(des, initial *TriggerTransport, opts ...dcl.ApplyOption) *TriggerTransport {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &TriggerTransport{}

	cDes.Pubsub = canonicalizeTriggerTransportPubsub(des.Pubsub, initial.Pubsub, opts...)

	return cDes
}

func canonicalizeTriggerTransportSlice(des, initial []TriggerTransport, opts ...dcl.ApplyOption) []TriggerTransport {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]TriggerTransport, 0, len(des))
		for _, d := range des {
			cd := canonicalizeTriggerTransport(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]TriggerTransport, 0, len(des))
	for i, d := range des {
		cd := canonicalizeTriggerTransport(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewTriggerTransport(c *Client, des, nw *TriggerTransport) *TriggerTransport {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for TriggerTransport while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Pubsub = canonicalizeNewTriggerTransportPubsub(c, des.Pubsub, nw.Pubsub)

	return nw
}

func canonicalizeNewTriggerTransportSet(c *Client, des, nw []TriggerTransport) []TriggerTransport {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []TriggerTransport
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareTriggerTransportNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewTriggerTransport(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewTriggerTransportSlice(c *Client, des, nw []TriggerTransport) []TriggerTransport {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []TriggerTransport
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewTriggerTransport(c, &d, &n))
	}

	return items
}

func canonicalizeTriggerTransportPubsub(des, initial *TriggerTransportPubsub, opts ...dcl.ApplyOption) *TriggerTransportPubsub {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &TriggerTransportPubsub{}

	if dcl.PartialSelfLinkToSelfLink(des.Topic, initial.Topic) || dcl.IsZeroValue(des.Topic) {
		cDes.Topic = initial.Topic
	} else {
		cDes.Topic = des.Topic
	}

	return cDes
}

func canonicalizeTriggerTransportPubsubSlice(des, initial []TriggerTransportPubsub, opts ...dcl.ApplyOption) []TriggerTransportPubsub {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]TriggerTransportPubsub, 0, len(des))
		for _, d := range des {
			cd := canonicalizeTriggerTransportPubsub(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]TriggerTransportPubsub, 0, len(des))
	for i, d := range des {
		cd := canonicalizeTriggerTransportPubsub(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewTriggerTransportPubsub(c *Client, des, nw *TriggerTransportPubsub) *TriggerTransportPubsub {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for TriggerTransportPubsub while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.PartialSelfLinkToSelfLink(des.Topic, nw.Topic) {
		nw.Topic = des.Topic
	}
	if dcl.PartialSelfLinkToSelfLink(des.Subscription, nw.Subscription) {
		nw.Subscription = des.Subscription
	}

	return nw
}

func canonicalizeNewTriggerTransportPubsubSet(c *Client, des, nw []TriggerTransportPubsub) []TriggerTransportPubsub {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []TriggerTransportPubsub
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareTriggerTransportPubsubNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewTriggerTransportPubsub(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewTriggerTransportPubsubSlice(c *Client, des, nw []TriggerTransportPubsub) []TriggerTransportPubsub {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []TriggerTransportPubsub
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewTriggerTransportPubsub(c, &d, &n))
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
func diffTrigger(c *Client, desired, actual *Trigger, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Uid, actual.Uid, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uid")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.MatchingCriteria, actual.MatchingCriteria, dcl.DiffInfo{Type: "Set", ObjectFunction: compareTriggerMatchingCriteriaNewStyle, EmptyObject: EmptyTriggerMatchingCriteria, OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("EventFilters")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccount, actual.ServiceAccount, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("ServiceAccount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Destination, actual.Destination, dcl.DiffInfo{MergeNestedDiffs: true, ObjectFunction: compareTriggerDestinationNewStyle, EmptyObject: EmptyTriggerDestination, OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("Destination")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Transport, actual.Transport, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareTriggerTransportNewStyle, EmptyObject: EmptyTriggerTransport, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Transport")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Etag, actual.Etag, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Etag")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Channel, actual.Channel, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Channel")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Conditions, actual.Conditions, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Conditions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EventDataContentType, actual.EventDataContentType, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("EventDataContentType")); len(ds) != 0 || err != nil {
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
func compareTriggerMatchingCriteriaNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*TriggerMatchingCriteria)
	if !ok {
		desiredNotPointer, ok := d.(TriggerMatchingCriteria)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerMatchingCriteria or *TriggerMatchingCriteria", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*TriggerMatchingCriteria)
	if !ok {
		actualNotPointer, ok := a.(TriggerMatchingCriteria)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerMatchingCriteria", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Attribute, actual.Attribute, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Attribute")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Operator, actual.Operator, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("Operator")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareTriggerDestinationNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*TriggerDestination)
	if !ok {
		desiredNotPointer, ok := d.(TriggerDestination)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestination or *TriggerDestination", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*TriggerDestination)
	if !ok {
		actualNotPointer, ok := a.(TriggerDestination)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestination", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CloudRunService, actual.CloudRunService, dcl.DiffInfo{ObjectFunction: compareTriggerDestinationCloudRunServiceNewStyle, EmptyObject: EmptyTriggerDestinationCloudRunService, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CloudRun")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CloudFunction, actual.CloudFunction, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CloudFunction")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Gke, actual.Gke, dcl.DiffInfo{ObjectFunction: compareTriggerDestinationGkeNewStyle, EmptyObject: EmptyTriggerDestinationGke, OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("Gke")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Workflow, actual.Workflow, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("Workflow")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.HttpEndpoint, actual.HttpEndpoint, dcl.DiffInfo{ObjectFunction: compareTriggerDestinationHttpEndpointNewStyle, EmptyObject: EmptyTriggerDestinationHttpEndpoint, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HttpEndpoint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NetworkConfig, actual.NetworkConfig, dcl.DiffInfo{ObjectFunction: compareTriggerDestinationNetworkConfigNewStyle, EmptyObject: EmptyTriggerDestinationNetworkConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareTriggerDestinationCloudRunServiceNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*TriggerDestinationCloudRunService)
	if !ok {
		desiredNotPointer, ok := d.(TriggerDestinationCloudRunService)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestinationCloudRunService or *TriggerDestinationCloudRunService", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*TriggerDestinationCloudRunService)
	if !ok {
		actualNotPointer, ok := a.(TriggerDestinationCloudRunService)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestinationCloudRunService", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Service, actual.Service, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Service")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Path, actual.Path, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Path")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Region, actual.Region, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Region")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareTriggerDestinationGkeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*TriggerDestinationGke)
	if !ok {
		desiredNotPointer, ok := d.(TriggerDestinationGke)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestinationGke or *TriggerDestinationGke", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*TriggerDestinationGke)
	if !ok {
		actualNotPointer, ok := a.(TriggerDestinationGke)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestinationGke", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Cluster, actual.Cluster, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Cluster")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Namespace, actual.Namespace, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("Namespace")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Service, actual.Service, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("Service")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Path, actual.Path, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateTriggerUpdateTriggerOperation")}, fn.AddNest("Path")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareTriggerDestinationHttpEndpointNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*TriggerDestinationHttpEndpoint)
	if !ok {
		desiredNotPointer, ok := d.(TriggerDestinationHttpEndpoint)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestinationHttpEndpoint or *TriggerDestinationHttpEndpoint", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*TriggerDestinationHttpEndpoint)
	if !ok {
		actualNotPointer, ok := a.(TriggerDestinationHttpEndpoint)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestinationHttpEndpoint", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareTriggerDestinationNetworkConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*TriggerDestinationNetworkConfig)
	if !ok {
		desiredNotPointer, ok := d.(TriggerDestinationNetworkConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestinationNetworkConfig or *TriggerDestinationNetworkConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*TriggerDestinationNetworkConfig)
	if !ok {
		actualNotPointer, ok := a.(TriggerDestinationNetworkConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerDestinationNetworkConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NetworkAttachment, actual.NetworkAttachment, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkAttachment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareTriggerTransportNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*TriggerTransport)
	if !ok {
		desiredNotPointer, ok := d.(TriggerTransport)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerTransport or *TriggerTransport", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*TriggerTransport)
	if !ok {
		actualNotPointer, ok := a.(TriggerTransport)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerTransport", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Pubsub, actual.Pubsub, dcl.DiffInfo{ObjectFunction: compareTriggerTransportPubsubNewStyle, EmptyObject: EmptyTriggerTransportPubsub, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Pubsub")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareTriggerTransportPubsubNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*TriggerTransportPubsub)
	if !ok {
		desiredNotPointer, ok := d.(TriggerTransportPubsub)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerTransportPubsub or *TriggerTransportPubsub", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*TriggerTransportPubsub)
	if !ok {
		actualNotPointer, ok := a.(TriggerTransportPubsub)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a TriggerTransportPubsub", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Topic, actual.Topic, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Topic")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Subscription, actual.Subscription, dcl.DiffInfo{OutputOnly: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Subscription")); len(ds) != 0 || err != nil {
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
func (r *Trigger) urlNormalized() *Trigger {
	normalized := dcl.Copy(*r).(Trigger)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.ServiceAccount = dcl.SelfLinkToName(r.ServiceAccount)
	normalized.Etag = dcl.SelfLinkToName(r.Etag)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.Channel = dcl.SelfLinkToName(r.Channel)
	normalized.EventDataContentType = dcl.SelfLinkToName(r.EventDataContentType)
	return &normalized
}

func (r *Trigger) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateTrigger" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/triggers/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Trigger resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Trigger) marshal(c *Client) ([]byte, error) {
	m, err := expandTrigger(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Trigger: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalTrigger decodes JSON responses into the Trigger resource schema.
func unmarshalTrigger(b []byte, c *Client, res *Trigger) (*Trigger, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapTrigger(m, c, res)
}

func unmarshalMapTrigger(m map[string]interface{}, c *Client, res *Trigger) (*Trigger, error) {

	flattened := flattenTrigger(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandTrigger expands Trigger into a JSON request object.
func expandTrigger(c *Client, f *Trigger) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/triggers/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v, err := expandTriggerMatchingCriteriaSlice(c, f.MatchingCriteria, res); err != nil {
		return nil, fmt.Errorf("error expanding MatchingCriteria into eventFilters: %w", err)
	} else if v != nil {
		m["eventFilters"] = v
	}
	if v := f.ServiceAccount; dcl.ValueShouldBeSent(v) {
		m["serviceAccount"] = v
	}
	if v, err := expandTriggerDestination(c, f.Destination, res); err != nil {
		return nil, fmt.Errorf("error expanding Destination into destination: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["destination"] = v
	}
	if v, err := expandTriggerTransport(c, f.Transport, res); err != nil {
		return nil, fmt.Errorf("error expanding Transport into transport: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["transport"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
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
	if v := f.Channel; dcl.ValueShouldBeSent(v) {
		m["channel"] = v
	}
	if v := f.EventDataContentType; dcl.ValueShouldBeSent(v) {
		m["eventDataContentType"] = v
	}

	return m, nil
}

// flattenTrigger flattens Trigger from a JSON request object into the
// Trigger type.
func flattenTrigger(c *Client, i interface{}, res *Trigger) *Trigger {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Trigger{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.MatchingCriteria = flattenTriggerMatchingCriteriaSlice(c, m["eventFilters"], res)
	resultRes.ServiceAccount = dcl.FlattenString(m["serviceAccount"])
	resultRes.Destination = flattenTriggerDestination(c, m["destination"], res)
	resultRes.Transport = flattenTriggerTransport(c, m["transport"], res)
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Etag = dcl.FlattenString(m["etag"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.Channel = dcl.FlattenString(m["channel"])
	resultRes.Conditions = dcl.FlattenKeyValuePairs(m["conditions"])
	resultRes.EventDataContentType = dcl.FlattenString(m["eventDataContentType"])

	return resultRes
}

// expandTriggerMatchingCriteriaMap expands the contents of TriggerMatchingCriteria into a JSON
// request object.
func expandTriggerMatchingCriteriaMap(c *Client, f map[string]TriggerMatchingCriteria, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandTriggerMatchingCriteria(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandTriggerMatchingCriteriaSlice expands the contents of TriggerMatchingCriteria into a JSON
// request object.
func expandTriggerMatchingCriteriaSlice(c *Client, f []TriggerMatchingCriteria, res *Trigger) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandTriggerMatchingCriteria(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenTriggerMatchingCriteriaMap flattens the contents of TriggerMatchingCriteria from a JSON
// response object.
func flattenTriggerMatchingCriteriaMap(c *Client, i interface{}, res *Trigger) map[string]TriggerMatchingCriteria {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]TriggerMatchingCriteria{}
	}

	if len(a) == 0 {
		return map[string]TriggerMatchingCriteria{}
	}

	items := make(map[string]TriggerMatchingCriteria)
	for k, item := range a {
		items[k] = *flattenTriggerMatchingCriteria(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenTriggerMatchingCriteriaSlice flattens the contents of TriggerMatchingCriteria from a JSON
// response object.
func flattenTriggerMatchingCriteriaSlice(c *Client, i interface{}, res *Trigger) []TriggerMatchingCriteria {
	a, ok := i.([]interface{})
	if !ok {
		return []TriggerMatchingCriteria{}
	}

	if len(a) == 0 {
		return []TriggerMatchingCriteria{}
	}

	items := make([]TriggerMatchingCriteria, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenTriggerMatchingCriteria(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandTriggerMatchingCriteria expands an instance of TriggerMatchingCriteria into a JSON
// request object.
func expandTriggerMatchingCriteria(c *Client, f *TriggerMatchingCriteria, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Attribute; !dcl.IsEmptyValueIndirect(v) {
		m["attribute"] = v
	}
	if v := f.Value; !dcl.IsEmptyValueIndirect(v) {
		m["value"] = v
	}
	if v := f.Operator; !dcl.IsEmptyValueIndirect(v) {
		m["operator"] = v
	}

	return m, nil
}

// flattenTriggerMatchingCriteria flattens an instance of TriggerMatchingCriteria from a JSON
// response object.
func flattenTriggerMatchingCriteria(c *Client, i interface{}, res *Trigger) *TriggerMatchingCriteria {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &TriggerMatchingCriteria{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyTriggerMatchingCriteria
	}
	r.Attribute = dcl.FlattenString(m["attribute"])
	r.Value = dcl.FlattenString(m["value"])
	r.Operator = dcl.FlattenString(m["operator"])

	return r
}

// expandTriggerDestinationMap expands the contents of TriggerDestination into a JSON
// request object.
func expandTriggerDestinationMap(c *Client, f map[string]TriggerDestination, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandTriggerDestination(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandTriggerDestinationSlice expands the contents of TriggerDestination into a JSON
// request object.
func expandTriggerDestinationSlice(c *Client, f []TriggerDestination, res *Trigger) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandTriggerDestination(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenTriggerDestinationMap flattens the contents of TriggerDestination from a JSON
// response object.
func flattenTriggerDestinationMap(c *Client, i interface{}, res *Trigger) map[string]TriggerDestination {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]TriggerDestination{}
	}

	if len(a) == 0 {
		return map[string]TriggerDestination{}
	}

	items := make(map[string]TriggerDestination)
	for k, item := range a {
		items[k] = *flattenTriggerDestination(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenTriggerDestinationSlice flattens the contents of TriggerDestination from a JSON
// response object.
func flattenTriggerDestinationSlice(c *Client, i interface{}, res *Trigger) []TriggerDestination {
	a, ok := i.([]interface{})
	if !ok {
		return []TriggerDestination{}
	}

	if len(a) == 0 {
		return []TriggerDestination{}
	}

	items := make([]TriggerDestination, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenTriggerDestination(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandTriggerDestination expands an instance of TriggerDestination into a JSON
// request object.
func expandTriggerDestination(c *Client, f *TriggerDestination, res *Trigger) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandTriggerDestinationCloudRunService(c, f.CloudRunService, res); err != nil {
		return nil, fmt.Errorf("error expanding CloudRunService into cloudRun: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["cloudRun"] = v
	}
	if v, err := dcl.DeriveField("projects/%s/locations/%s/functions/%s", f.CloudFunction, dcl.SelfLinkToName(res.Project), dcl.SelfLinkToName(res.Location), dcl.SelfLinkToName(f.CloudFunction)); err != nil {
		return nil, fmt.Errorf("error expanding CloudFunction into cloudFunction: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["cloudFunction"] = v
	}
	if v, err := expandTriggerDestinationGke(c, f.Gke, res); err != nil {
		return nil, fmt.Errorf("error expanding Gke into gke: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gke"] = v
	}
	if v, err := dcl.DeriveField("projects/%s/locations/%s/workflows/%s", f.Workflow, dcl.SelfLinkToName(res.Project), dcl.SelfLinkToName(res.Location), dcl.SelfLinkToName(f.Workflow)); err != nil {
		return nil, fmt.Errorf("error expanding Workflow into workflow: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["workflow"] = v
	}
	if v, err := expandTriggerDestinationHttpEndpoint(c, f.HttpEndpoint, res); err != nil {
		return nil, fmt.Errorf("error expanding HttpEndpoint into httpEndpoint: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["httpEndpoint"] = v
	}
	if v, err := expandTriggerDestinationNetworkConfig(c, f.NetworkConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding NetworkConfig into networkConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["networkConfig"] = v
	}

	return m, nil
}

// flattenTriggerDestination flattens an instance of TriggerDestination from a JSON
// response object.
func flattenTriggerDestination(c *Client, i interface{}, res *Trigger) *TriggerDestination {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &TriggerDestination{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyTriggerDestination
	}
	r.CloudRunService = flattenTriggerDestinationCloudRunService(c, m["cloudRun"], res)
	r.CloudFunction = dcl.FlattenString(m["cloudFunction"])
	r.Gke = flattenTriggerDestinationGke(c, m["gke"], res)
	r.Workflow = dcl.FlattenString(m["workflow"])
	r.HttpEndpoint = flattenTriggerDestinationHttpEndpoint(c, m["httpEndpoint"], res)
	r.NetworkConfig = flattenTriggerDestinationNetworkConfig(c, m["networkConfig"], res)

	return r
}

// expandTriggerDestinationCloudRunServiceMap expands the contents of TriggerDestinationCloudRunService into a JSON
// request object.
func expandTriggerDestinationCloudRunServiceMap(c *Client, f map[string]TriggerDestinationCloudRunService, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandTriggerDestinationCloudRunService(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandTriggerDestinationCloudRunServiceSlice expands the contents of TriggerDestinationCloudRunService into a JSON
// request object.
func expandTriggerDestinationCloudRunServiceSlice(c *Client, f []TriggerDestinationCloudRunService, res *Trigger) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandTriggerDestinationCloudRunService(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenTriggerDestinationCloudRunServiceMap flattens the contents of TriggerDestinationCloudRunService from a JSON
// response object.
func flattenTriggerDestinationCloudRunServiceMap(c *Client, i interface{}, res *Trigger) map[string]TriggerDestinationCloudRunService {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]TriggerDestinationCloudRunService{}
	}

	if len(a) == 0 {
		return map[string]TriggerDestinationCloudRunService{}
	}

	items := make(map[string]TriggerDestinationCloudRunService)
	for k, item := range a {
		items[k] = *flattenTriggerDestinationCloudRunService(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenTriggerDestinationCloudRunServiceSlice flattens the contents of TriggerDestinationCloudRunService from a JSON
// response object.
func flattenTriggerDestinationCloudRunServiceSlice(c *Client, i interface{}, res *Trigger) []TriggerDestinationCloudRunService {
	a, ok := i.([]interface{})
	if !ok {
		return []TriggerDestinationCloudRunService{}
	}

	if len(a) == 0 {
		return []TriggerDestinationCloudRunService{}
	}

	items := make([]TriggerDestinationCloudRunService, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenTriggerDestinationCloudRunService(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandTriggerDestinationCloudRunService expands an instance of TriggerDestinationCloudRunService into a JSON
// request object.
func expandTriggerDestinationCloudRunService(c *Client, f *TriggerDestinationCloudRunService, res *Trigger) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := dcl.SelfLinkToNameExpander(f.Service); err != nil {
		return nil, fmt.Errorf("error expanding Service into service: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["service"] = v
	}
	if v := f.Path; !dcl.IsEmptyValueIndirect(v) {
		m["path"] = v
	}
	if v := f.Region; !dcl.IsEmptyValueIndirect(v) {
		m["region"] = v
	}

	return m, nil
}

// flattenTriggerDestinationCloudRunService flattens an instance of TriggerDestinationCloudRunService from a JSON
// response object.
func flattenTriggerDestinationCloudRunService(c *Client, i interface{}, res *Trigger) *TriggerDestinationCloudRunService {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &TriggerDestinationCloudRunService{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyTriggerDestinationCloudRunService
	}
	r.Service = dcl.FlattenString(m["service"])
	r.Path = dcl.FlattenString(m["path"])
	r.Region = dcl.FlattenString(m["region"])

	return r
}

// expandTriggerDestinationGkeMap expands the contents of TriggerDestinationGke into a JSON
// request object.
func expandTriggerDestinationGkeMap(c *Client, f map[string]TriggerDestinationGke, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandTriggerDestinationGke(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandTriggerDestinationGkeSlice expands the contents of TriggerDestinationGke into a JSON
// request object.
func expandTriggerDestinationGkeSlice(c *Client, f []TriggerDestinationGke, res *Trigger) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandTriggerDestinationGke(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenTriggerDestinationGkeMap flattens the contents of TriggerDestinationGke from a JSON
// response object.
func flattenTriggerDestinationGkeMap(c *Client, i interface{}, res *Trigger) map[string]TriggerDestinationGke {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]TriggerDestinationGke{}
	}

	if len(a) == 0 {
		return map[string]TriggerDestinationGke{}
	}

	items := make(map[string]TriggerDestinationGke)
	for k, item := range a {
		items[k] = *flattenTriggerDestinationGke(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenTriggerDestinationGkeSlice flattens the contents of TriggerDestinationGke from a JSON
// response object.
func flattenTriggerDestinationGkeSlice(c *Client, i interface{}, res *Trigger) []TriggerDestinationGke {
	a, ok := i.([]interface{})
	if !ok {
		return []TriggerDestinationGke{}
	}

	if len(a) == 0 {
		return []TriggerDestinationGke{}
	}

	items := make([]TriggerDestinationGke, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenTriggerDestinationGke(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandTriggerDestinationGke expands an instance of TriggerDestinationGke into a JSON
// request object.
func expandTriggerDestinationGke(c *Client, f *TriggerDestinationGke, res *Trigger) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := dcl.SelfLinkToNameExpander(f.Cluster); err != nil {
		return nil, fmt.Errorf("error expanding Cluster into cluster: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["cluster"] = v
	}
	if v := f.Location; !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}
	if v := f.Namespace; !dcl.IsEmptyValueIndirect(v) {
		m["namespace"] = v
	}
	if v := f.Service; !dcl.IsEmptyValueIndirect(v) {
		m["service"] = v
	}
	if v := f.Path; !dcl.IsEmptyValueIndirect(v) {
		m["path"] = v
	}

	return m, nil
}

// flattenTriggerDestinationGke flattens an instance of TriggerDestinationGke from a JSON
// response object.
func flattenTriggerDestinationGke(c *Client, i interface{}, res *Trigger) *TriggerDestinationGke {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &TriggerDestinationGke{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyTriggerDestinationGke
	}
	r.Cluster = dcl.FlattenString(m["cluster"])
	r.Location = dcl.FlattenString(m["location"])
	r.Namespace = dcl.FlattenString(m["namespace"])
	r.Service = dcl.FlattenString(m["service"])
	r.Path = dcl.FlattenString(m["path"])

	return r
}

// expandTriggerDestinationHttpEndpointMap expands the contents of TriggerDestinationHttpEndpoint into a JSON
// request object.
func expandTriggerDestinationHttpEndpointMap(c *Client, f map[string]TriggerDestinationHttpEndpoint, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandTriggerDestinationHttpEndpoint(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandTriggerDestinationHttpEndpointSlice expands the contents of TriggerDestinationHttpEndpoint into a JSON
// request object.
func expandTriggerDestinationHttpEndpointSlice(c *Client, f []TriggerDestinationHttpEndpoint, res *Trigger) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandTriggerDestinationHttpEndpoint(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenTriggerDestinationHttpEndpointMap flattens the contents of TriggerDestinationHttpEndpoint from a JSON
// response object.
func flattenTriggerDestinationHttpEndpointMap(c *Client, i interface{}, res *Trigger) map[string]TriggerDestinationHttpEndpoint {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]TriggerDestinationHttpEndpoint{}
	}

	if len(a) == 0 {
		return map[string]TriggerDestinationHttpEndpoint{}
	}

	items := make(map[string]TriggerDestinationHttpEndpoint)
	for k, item := range a {
		items[k] = *flattenTriggerDestinationHttpEndpoint(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenTriggerDestinationHttpEndpointSlice flattens the contents of TriggerDestinationHttpEndpoint from a JSON
// response object.
func flattenTriggerDestinationHttpEndpointSlice(c *Client, i interface{}, res *Trigger) []TriggerDestinationHttpEndpoint {
	a, ok := i.([]interface{})
	if !ok {
		return []TriggerDestinationHttpEndpoint{}
	}

	if len(a) == 0 {
		return []TriggerDestinationHttpEndpoint{}
	}

	items := make([]TriggerDestinationHttpEndpoint, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenTriggerDestinationHttpEndpoint(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandTriggerDestinationHttpEndpoint expands an instance of TriggerDestinationHttpEndpoint into a JSON
// request object.
func expandTriggerDestinationHttpEndpoint(c *Client, f *TriggerDestinationHttpEndpoint, res *Trigger) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Uri; !dcl.IsEmptyValueIndirect(v) {
		m["uri"] = v
	}

	return m, nil
}

// flattenTriggerDestinationHttpEndpoint flattens an instance of TriggerDestinationHttpEndpoint from a JSON
// response object.
func flattenTriggerDestinationHttpEndpoint(c *Client, i interface{}, res *Trigger) *TriggerDestinationHttpEndpoint {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &TriggerDestinationHttpEndpoint{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyTriggerDestinationHttpEndpoint
	}
	r.Uri = dcl.FlattenString(m["uri"])

	return r
}

// expandTriggerDestinationNetworkConfigMap expands the contents of TriggerDestinationNetworkConfig into a JSON
// request object.
func expandTriggerDestinationNetworkConfigMap(c *Client, f map[string]TriggerDestinationNetworkConfig, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandTriggerDestinationNetworkConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandTriggerDestinationNetworkConfigSlice expands the contents of TriggerDestinationNetworkConfig into a JSON
// request object.
func expandTriggerDestinationNetworkConfigSlice(c *Client, f []TriggerDestinationNetworkConfig, res *Trigger) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandTriggerDestinationNetworkConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenTriggerDestinationNetworkConfigMap flattens the contents of TriggerDestinationNetworkConfig from a JSON
// response object.
func flattenTriggerDestinationNetworkConfigMap(c *Client, i interface{}, res *Trigger) map[string]TriggerDestinationNetworkConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]TriggerDestinationNetworkConfig{}
	}

	if len(a) == 0 {
		return map[string]TriggerDestinationNetworkConfig{}
	}

	items := make(map[string]TriggerDestinationNetworkConfig)
	for k, item := range a {
		items[k] = *flattenTriggerDestinationNetworkConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenTriggerDestinationNetworkConfigSlice flattens the contents of TriggerDestinationNetworkConfig from a JSON
// response object.
func flattenTriggerDestinationNetworkConfigSlice(c *Client, i interface{}, res *Trigger) []TriggerDestinationNetworkConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []TriggerDestinationNetworkConfig{}
	}

	if len(a) == 0 {
		return []TriggerDestinationNetworkConfig{}
	}

	items := make([]TriggerDestinationNetworkConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenTriggerDestinationNetworkConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandTriggerDestinationNetworkConfig expands an instance of TriggerDestinationNetworkConfig into a JSON
// request object.
func expandTriggerDestinationNetworkConfig(c *Client, f *TriggerDestinationNetworkConfig, res *Trigger) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NetworkAttachment; !dcl.IsEmptyValueIndirect(v) {
		m["networkAttachment"] = v
	}

	return m, nil
}

// flattenTriggerDestinationNetworkConfig flattens an instance of TriggerDestinationNetworkConfig from a JSON
// response object.
func flattenTriggerDestinationNetworkConfig(c *Client, i interface{}, res *Trigger) *TriggerDestinationNetworkConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &TriggerDestinationNetworkConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyTriggerDestinationNetworkConfig
	}
	r.NetworkAttachment = dcl.FlattenString(m["networkAttachment"])

	return r
}

// expandTriggerTransportMap expands the contents of TriggerTransport into a JSON
// request object.
func expandTriggerTransportMap(c *Client, f map[string]TriggerTransport, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandTriggerTransport(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandTriggerTransportSlice expands the contents of TriggerTransport into a JSON
// request object.
func expandTriggerTransportSlice(c *Client, f []TriggerTransport, res *Trigger) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandTriggerTransport(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenTriggerTransportMap flattens the contents of TriggerTransport from a JSON
// response object.
func flattenTriggerTransportMap(c *Client, i interface{}, res *Trigger) map[string]TriggerTransport {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]TriggerTransport{}
	}

	if len(a) == 0 {
		return map[string]TriggerTransport{}
	}

	items := make(map[string]TriggerTransport)
	for k, item := range a {
		items[k] = *flattenTriggerTransport(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenTriggerTransportSlice flattens the contents of TriggerTransport from a JSON
// response object.
func flattenTriggerTransportSlice(c *Client, i interface{}, res *Trigger) []TriggerTransport {
	a, ok := i.([]interface{})
	if !ok {
		return []TriggerTransport{}
	}

	if len(a) == 0 {
		return []TriggerTransport{}
	}

	items := make([]TriggerTransport, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenTriggerTransport(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandTriggerTransport expands an instance of TriggerTransport into a JSON
// request object.
func expandTriggerTransport(c *Client, f *TriggerTransport, res *Trigger) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandTriggerTransportPubsub(c, f.Pubsub, res); err != nil {
		return nil, fmt.Errorf("error expanding Pubsub into pubsub: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pubsub"] = v
	}

	return m, nil
}

// flattenTriggerTransport flattens an instance of TriggerTransport from a JSON
// response object.
func flattenTriggerTransport(c *Client, i interface{}, res *Trigger) *TriggerTransport {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &TriggerTransport{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyTriggerTransport
	}
	r.Pubsub = flattenTriggerTransportPubsub(c, m["pubsub"], res)

	return r
}

// expandTriggerTransportPubsubMap expands the contents of TriggerTransportPubsub into a JSON
// request object.
func expandTriggerTransportPubsubMap(c *Client, f map[string]TriggerTransportPubsub, res *Trigger) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandTriggerTransportPubsub(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandTriggerTransportPubsubSlice expands the contents of TriggerTransportPubsub into a JSON
// request object.
func expandTriggerTransportPubsubSlice(c *Client, f []TriggerTransportPubsub, res *Trigger) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandTriggerTransportPubsub(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenTriggerTransportPubsubMap flattens the contents of TriggerTransportPubsub from a JSON
// response object.
func flattenTriggerTransportPubsubMap(c *Client, i interface{}, res *Trigger) map[string]TriggerTransportPubsub {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]TriggerTransportPubsub{}
	}

	if len(a) == 0 {
		return map[string]TriggerTransportPubsub{}
	}

	items := make(map[string]TriggerTransportPubsub)
	for k, item := range a {
		items[k] = *flattenTriggerTransportPubsub(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenTriggerTransportPubsubSlice flattens the contents of TriggerTransportPubsub from a JSON
// response object.
func flattenTriggerTransportPubsubSlice(c *Client, i interface{}, res *Trigger) []TriggerTransportPubsub {
	a, ok := i.([]interface{})
	if !ok {
		return []TriggerTransportPubsub{}
	}

	if len(a) == 0 {
		return []TriggerTransportPubsub{}
	}

	items := make([]TriggerTransportPubsub, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenTriggerTransportPubsub(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandTriggerTransportPubsub expands an instance of TriggerTransportPubsub into a JSON
// request object.
func expandTriggerTransportPubsub(c *Client, f *TriggerTransportPubsub, res *Trigger) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := dcl.DeriveField("projects/%s/topics/%s", f.Topic, dcl.SelfLinkToName(res.Project), dcl.SelfLinkToName(f.Topic)); err != nil {
		return nil, fmt.Errorf("error expanding Topic into topic: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["topic"] = v
	}

	return m, nil
}

// flattenTriggerTransportPubsub flattens an instance of TriggerTransportPubsub from a JSON
// response object.
func flattenTriggerTransportPubsub(c *Client, i interface{}, res *Trigger) *TriggerTransportPubsub {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &TriggerTransportPubsub{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyTriggerTransportPubsub
	}
	r.Topic = dcl.FlattenString(m["topic"])
	r.Subscription = dcl.FlattenString(m["subscription"])

	return r
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Trigger) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalTrigger(b, c, r)
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

type triggerDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         triggerApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToTriggerDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]triggerDiff, error) {
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
	var diffs []triggerDiff
	// For each operation name, create a triggerDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := triggerDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToTriggerApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToTriggerApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (triggerApiOperation, error) {
	switch opName {

	case "updateTriggerUpdateTriggerOperation":
		return &updateTriggerUpdateTriggerOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractTriggerFields(r *Trigger) error {
	vDestination := r.Destination
	if vDestination == nil {
		// note: explicitly not the empty object.
		vDestination = &TriggerDestination{}
	}
	if err := extractTriggerDestinationFields(r, vDestination); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDestination) {
		r.Destination = vDestination
	}
	vTransport := r.Transport
	if vTransport == nil {
		// note: explicitly not the empty object.
		vTransport = &TriggerTransport{}
	}
	if err := extractTriggerTransportFields(r, vTransport); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vTransport) {
		r.Transport = vTransport
	}
	return nil
}
func extractTriggerMatchingCriteriaFields(r *Trigger, o *TriggerMatchingCriteria) error {
	return nil
}
func extractTriggerDestinationFields(r *Trigger, o *TriggerDestination) error {
	vCloudRunService := o.CloudRunService
	if vCloudRunService == nil {
		// note: explicitly not the empty object.
		vCloudRunService = &TriggerDestinationCloudRunService{}
	}
	if err := extractTriggerDestinationCloudRunServiceFields(r, vCloudRunService); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCloudRunService) {
		o.CloudRunService = vCloudRunService
	}
	vGke := o.Gke
	if vGke == nil {
		// note: explicitly not the empty object.
		vGke = &TriggerDestinationGke{}
	}
	if err := extractTriggerDestinationGkeFields(r, vGke); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGke) {
		o.Gke = vGke
	}
	vHttpEndpoint := o.HttpEndpoint
	if vHttpEndpoint == nil {
		// note: explicitly not the empty object.
		vHttpEndpoint = &TriggerDestinationHttpEndpoint{}
	}
	if err := extractTriggerDestinationHttpEndpointFields(r, vHttpEndpoint); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vHttpEndpoint) {
		o.HttpEndpoint = vHttpEndpoint
	}
	vNetworkConfig := o.NetworkConfig
	if vNetworkConfig == nil {
		// note: explicitly not the empty object.
		vNetworkConfig = &TriggerDestinationNetworkConfig{}
	}
	if err := extractTriggerDestinationNetworkConfigFields(r, vNetworkConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetworkConfig) {
		o.NetworkConfig = vNetworkConfig
	}
	return nil
}
func extractTriggerDestinationCloudRunServiceFields(r *Trigger, o *TriggerDestinationCloudRunService) error {
	return nil
}
func extractTriggerDestinationGkeFields(r *Trigger, o *TriggerDestinationGke) error {
	return nil
}
func extractTriggerDestinationHttpEndpointFields(r *Trigger, o *TriggerDestinationHttpEndpoint) error {
	return nil
}
func extractTriggerDestinationNetworkConfigFields(r *Trigger, o *TriggerDestinationNetworkConfig) error {
	return nil
}
func extractTriggerTransportFields(r *Trigger, o *TriggerTransport) error {
	vPubsub := o.Pubsub
	if vPubsub == nil {
		// note: explicitly not the empty object.
		vPubsub = &TriggerTransportPubsub{}
	}
	if err := extractTriggerTransportPubsubFields(r, vPubsub); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPubsub) {
		o.Pubsub = vPubsub
	}
	return nil
}
func extractTriggerTransportPubsubFields(r *Trigger, o *TriggerTransportPubsub) error {
	return nil
}

func postReadExtractTriggerFields(r *Trigger) error {
	vDestination := r.Destination
	if vDestination == nil {
		// note: explicitly not the empty object.
		vDestination = &TriggerDestination{}
	}
	if err := postReadExtractTriggerDestinationFields(r, vDestination); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDestination) {
		r.Destination = vDestination
	}
	vTransport := r.Transport
	if vTransport == nil {
		// note: explicitly not the empty object.
		vTransport = &TriggerTransport{}
	}
	if err := postReadExtractTriggerTransportFields(r, vTransport); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vTransport) {
		r.Transport = vTransport
	}
	return nil
}
func postReadExtractTriggerMatchingCriteriaFields(r *Trigger, o *TriggerMatchingCriteria) error {
	return nil
}
func postReadExtractTriggerDestinationFields(r *Trigger, o *TriggerDestination) error {
	vCloudRunService := o.CloudRunService
	if vCloudRunService == nil {
		// note: explicitly not the empty object.
		vCloudRunService = &TriggerDestinationCloudRunService{}
	}
	if err := extractTriggerDestinationCloudRunServiceFields(r, vCloudRunService); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCloudRunService) {
		o.CloudRunService = vCloudRunService
	}
	vGke := o.Gke
	if vGke == nil {
		// note: explicitly not the empty object.
		vGke = &TriggerDestinationGke{}
	}
	if err := extractTriggerDestinationGkeFields(r, vGke); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGke) {
		o.Gke = vGke
	}
	vHttpEndpoint := o.HttpEndpoint
	if vHttpEndpoint == nil {
		// note: explicitly not the empty object.
		vHttpEndpoint = &TriggerDestinationHttpEndpoint{}
	}
	if err := extractTriggerDestinationHttpEndpointFields(r, vHttpEndpoint); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vHttpEndpoint) {
		o.HttpEndpoint = vHttpEndpoint
	}
	vNetworkConfig := o.NetworkConfig
	if vNetworkConfig == nil {
		// note: explicitly not the empty object.
		vNetworkConfig = &TriggerDestinationNetworkConfig{}
	}
	if err := extractTriggerDestinationNetworkConfigFields(r, vNetworkConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetworkConfig) {
		o.NetworkConfig = vNetworkConfig
	}
	return nil
}
func postReadExtractTriggerDestinationCloudRunServiceFields(r *Trigger, o *TriggerDestinationCloudRunService) error {
	return nil
}
func postReadExtractTriggerDestinationGkeFields(r *Trigger, o *TriggerDestinationGke) error {
	return nil
}
func postReadExtractTriggerDestinationHttpEndpointFields(r *Trigger, o *TriggerDestinationHttpEndpoint) error {
	return nil
}
func postReadExtractTriggerDestinationNetworkConfigFields(r *Trigger, o *TriggerDestinationNetworkConfig) error {
	return nil
}
func postReadExtractTriggerTransportFields(r *Trigger, o *TriggerTransport) error {
	vPubsub := o.Pubsub
	if vPubsub == nil {
		// note: explicitly not the empty object.
		vPubsub = &TriggerTransportPubsub{}
	}
	if err := extractTriggerTransportPubsubFields(r, vPubsub); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPubsub) {
		o.Pubsub = vPubsub
	}
	return nil
}
func postReadExtractTriggerTransportPubsubFields(r *Trigger, o *TriggerTransportPubsub) error {
	return nil
}
