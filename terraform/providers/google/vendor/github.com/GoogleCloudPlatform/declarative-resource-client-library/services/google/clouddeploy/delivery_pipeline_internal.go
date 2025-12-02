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
package clouddeploy

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

func (r *DeliveryPipeline) validate() error {

	if err := dcl.RequiredParameter(r.Name, "Name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.SerialPipeline) {
		if err := r.SerialPipeline.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Condition) {
		if err := r.Condition.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *DeliveryPipelineSerialPipeline) validate() error {
	return nil
}
func (r *DeliveryPipelineSerialPipelineStages) validate() error {
	return nil
}
func (r *DeliveryPipelineCondition) validate() error {
	if !dcl.IsEmptyValueIndirect(r.PipelineReadyCondition) {
		if err := r.PipelineReadyCondition.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.TargetsPresentCondition) {
		if err := r.TargetsPresentCondition.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *DeliveryPipelineConditionPipelineReadyCondition) validate() error {
	return nil
}
func (r *DeliveryPipelineConditionTargetsPresentCondition) validate() error {
	return nil
}
func (r *DeliveryPipeline) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://clouddeploy.googleapis.com/v1/", params)
}

func (r *DeliveryPipeline) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *DeliveryPipeline) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/deliveryPipelines", nr.basePath(), userBasePath, params), nil

}

func (r *DeliveryPipeline) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/deliveryPipelines?deliveryPipelineId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *DeliveryPipeline) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *DeliveryPipeline) SetPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{}
	return dcl.URL("", nr.basePath(), userBasePath, fields)
}

func (r *DeliveryPipeline) SetPolicyVerb() string {
	return ""
}

func (r *DeliveryPipeline) getPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{}
	return dcl.URL("", nr.basePath(), userBasePath, fields)
}

func (r *DeliveryPipeline) IAMPolicyVersion() int {
	return 3
}

// deliveryPipelineApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type deliveryPipelineApiOperation interface {
	do(context.Context, *DeliveryPipeline, *Client) error
}

// newUpdateDeliveryPipelineUpdateDeliveryPipelineRequest creates a request for an
// DeliveryPipeline resource's UpdateDeliveryPipeline update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateDeliveryPipelineUpdateDeliveryPipelineRequest(ctx context.Context, f *DeliveryPipeline, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v := f.Annotations; !dcl.IsEmptyValueIndirect(v) {
		req["annotations"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	if v, err := expandDeliveryPipelineSerialPipeline(c, f.SerialPipeline, res); err != nil {
		return nil, fmt.Errorf("error expanding SerialPipeline into serialPipeline: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["serialPipeline"] = v
	}
	if v, err := expandDeliveryPipelineCondition(c, f.Condition, res); err != nil {
		return nil, fmt.Errorf("error expanding Condition into condition: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["condition"] = v
	}
	if v := f.Suspended; !dcl.IsEmptyValueIndirect(v) {
		req["suspended"] = v
	}
	b, err := c.getDeliveryPipelineRaw(ctx, f)
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
	req["name"] = fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s", *f.Project, *f.Location, *f.Name)

	return req, nil
}

// marshalUpdateDeliveryPipelineUpdateDeliveryPipelineRequest converts the update into
// the final JSON request body.
func marshalUpdateDeliveryPipelineUpdateDeliveryPipelineRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateDeliveryPipelineUpdateDeliveryPipelineOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateDeliveryPipelineUpdateDeliveryPipelineOperation) do(ctx context.Context, r *DeliveryPipeline, c *Client) error {
	_, err := c.GetDeliveryPipeline(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateDeliveryPipeline")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateDeliveryPipelineUpdateDeliveryPipelineRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateDeliveryPipelineUpdateDeliveryPipelineRequest(c, req)
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

func (c *Client) listDeliveryPipelineRaw(ctx context.Context, r *DeliveryPipeline, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != DeliveryPipelineMaxPage {
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

type listDeliveryPipelineOperation struct {
	DeliveryPipelines []map[string]interface{} `json:"deliveryPipelines"`
	Token             string                   `json:"nextPageToken"`
}

func (c *Client) listDeliveryPipeline(ctx context.Context, r *DeliveryPipeline, pageToken string, pageSize int32) ([]*DeliveryPipeline, string, error) {
	b, err := c.listDeliveryPipelineRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listDeliveryPipelineOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*DeliveryPipeline
	for _, v := range m.DeliveryPipelines {
		res, err := unmarshalMapDeliveryPipeline(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllDeliveryPipeline(ctx context.Context, f func(*DeliveryPipeline) bool, resources []*DeliveryPipeline) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteDeliveryPipeline(ctx, res)
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

type deleteDeliveryPipelineOperation struct{}

func (op *deleteDeliveryPipelineOperation) do(ctx context.Context, r *DeliveryPipeline, c *Client) error {
	r, err := c.GetDeliveryPipeline(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "DeliveryPipeline not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetDeliveryPipeline checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	u, err = dcl.AddQueryParams(u, map[string]string{"force": "true"})
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
		_, err := c.GetDeliveryPipeline(ctx, r)
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
type createDeliveryPipelineOperation struct {
	response map[string]interface{}
}

func (op *createDeliveryPipelineOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createDeliveryPipelineOperation) do(ctx context.Context, r *DeliveryPipeline, c *Client) error {
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

	if _, err := c.GetDeliveryPipeline(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getDeliveryPipelineRaw(ctx context.Context, r *DeliveryPipeline) ([]byte, error) {

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

func (c *Client) deliveryPipelineDiffsForRawDesired(ctx context.Context, rawDesired *DeliveryPipeline, opts ...dcl.ApplyOption) (initial, desired *DeliveryPipeline, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *DeliveryPipeline
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*DeliveryPipeline); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected DeliveryPipeline, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetDeliveryPipeline(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a DeliveryPipeline resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve DeliveryPipeline resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that DeliveryPipeline resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeDeliveryPipelineDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for DeliveryPipeline: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for DeliveryPipeline: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractDeliveryPipelineFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeDeliveryPipelineInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for DeliveryPipeline: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeDeliveryPipelineDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for DeliveryPipeline: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffDeliveryPipeline(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeDeliveryPipelineInitialState(rawInitial, rawDesired *DeliveryPipeline) (*DeliveryPipeline, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeDeliveryPipelineDesiredState(rawDesired, rawInitial *DeliveryPipeline, opts ...dcl.ApplyOption) (*DeliveryPipeline, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.SerialPipeline = canonicalizeDeliveryPipelineSerialPipeline(rawDesired.SerialPipeline, nil, opts...)
		rawDesired.Condition = canonicalizeDeliveryPipelineCondition(rawDesired.Condition, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &DeliveryPipeline{}
	if dcl.NameToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	if dcl.IsZeroValue(rawDesired.Annotations) || (dcl.IsEmptyValueIndirect(rawDesired.Annotations) && dcl.IsEmptyValueIndirect(rawInitial.Annotations)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Annotations = rawInitial.Annotations
	} else {
		canonicalDesired.Annotations = rawDesired.Annotations
	}
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	canonicalDesired.SerialPipeline = canonicalizeDeliveryPipelineSerialPipeline(rawDesired.SerialPipeline, rawInitial.SerialPipeline, opts...)
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
	if dcl.BoolCanonicalize(rawDesired.Suspended, rawInitial.Suspended) {
		canonicalDesired.Suspended = rawInitial.Suspended
	} else {
		canonicalDesired.Suspended = rawDesired.Suspended
	}
	return canonicalDesired, nil
}

func canonicalizeDeliveryPipelineNewState(c *Client, rawNew, rawDesired *DeliveryPipeline) (*DeliveryPipeline, error) {

	rawNew.Name = rawDesired.Name

	if dcl.IsEmptyValueIndirect(rawNew.Uid) && dcl.IsEmptyValueIndirect(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Annotations) && dcl.IsEmptyValueIndirect(rawDesired.Annotations) {
		rawNew.Annotations = rawDesired.Annotations
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.SerialPipeline) && dcl.IsEmptyValueIndirect(rawDesired.SerialPipeline) {
		rawNew.SerialPipeline = rawDesired.SerialPipeline
	} else {
		rawNew.SerialPipeline = canonicalizeNewDeliveryPipelineSerialPipeline(c, rawDesired.SerialPipeline, rawNew.SerialPipeline)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Condition) && dcl.IsEmptyValueIndirect(rawDesired.Condition) {
		rawNew.Condition = rawDesired.Condition
	} else {
		rawNew.Condition = canonicalizeNewDeliveryPipelineCondition(c, rawDesired.Condition, rawNew.Condition)
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

	if dcl.IsEmptyValueIndirect(rawNew.Suspended) && dcl.IsEmptyValueIndirect(rawDesired.Suspended) {
		rawNew.Suspended = rawDesired.Suspended
	} else {
		if dcl.BoolCanonicalize(rawDesired.Suspended, rawNew.Suspended) {
			rawNew.Suspended = rawDesired.Suspended
		}
	}

	return rawNew, nil
}

func canonicalizeDeliveryPipelineSerialPipeline(des, initial *DeliveryPipelineSerialPipeline, opts ...dcl.ApplyOption) *DeliveryPipelineSerialPipeline {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &DeliveryPipelineSerialPipeline{}

	cDes.Stages = canonicalizeDeliveryPipelineSerialPipelineStagesSlice(des.Stages, initial.Stages, opts...)

	return cDes
}

func canonicalizeDeliveryPipelineSerialPipelineSlice(des, initial []DeliveryPipelineSerialPipeline, opts ...dcl.ApplyOption) []DeliveryPipelineSerialPipeline {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]DeliveryPipelineSerialPipeline, 0, len(des))
		for _, d := range des {
			cd := canonicalizeDeliveryPipelineSerialPipeline(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]DeliveryPipelineSerialPipeline, 0, len(des))
	for i, d := range des {
		cd := canonicalizeDeliveryPipelineSerialPipeline(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewDeliveryPipelineSerialPipeline(c *Client, des, nw *DeliveryPipelineSerialPipeline) *DeliveryPipelineSerialPipeline {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for DeliveryPipelineSerialPipeline while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Stages = canonicalizeNewDeliveryPipelineSerialPipelineStagesSlice(c, des.Stages, nw.Stages)

	return nw
}

func canonicalizeNewDeliveryPipelineSerialPipelineSet(c *Client, des, nw []DeliveryPipelineSerialPipeline) []DeliveryPipelineSerialPipeline {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []DeliveryPipelineSerialPipeline
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareDeliveryPipelineSerialPipelineNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewDeliveryPipelineSerialPipeline(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewDeliveryPipelineSerialPipelineSlice(c *Client, des, nw []DeliveryPipelineSerialPipeline) []DeliveryPipelineSerialPipeline {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []DeliveryPipelineSerialPipeline
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewDeliveryPipelineSerialPipeline(c, &d, &n))
	}

	return items
}

func canonicalizeDeliveryPipelineSerialPipelineStages(des, initial *DeliveryPipelineSerialPipelineStages, opts ...dcl.ApplyOption) *DeliveryPipelineSerialPipelineStages {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &DeliveryPipelineSerialPipelineStages{}

	if dcl.StringCanonicalize(des.TargetId, initial.TargetId) || dcl.IsZeroValue(des.TargetId) {
		cDes.TargetId = initial.TargetId
	} else {
		cDes.TargetId = des.TargetId
	}
	if dcl.StringArrayCanonicalize(des.Profiles, initial.Profiles) {
		cDes.Profiles = initial.Profiles
	} else {
		cDes.Profiles = des.Profiles
	}

	return cDes
}

func canonicalizeDeliveryPipelineSerialPipelineStagesSlice(des, initial []DeliveryPipelineSerialPipelineStages, opts ...dcl.ApplyOption) []DeliveryPipelineSerialPipelineStages {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]DeliveryPipelineSerialPipelineStages, 0, len(des))
		for _, d := range des {
			cd := canonicalizeDeliveryPipelineSerialPipelineStages(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]DeliveryPipelineSerialPipelineStages, 0, len(des))
	for i, d := range des {
		cd := canonicalizeDeliveryPipelineSerialPipelineStages(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewDeliveryPipelineSerialPipelineStages(c *Client, des, nw *DeliveryPipelineSerialPipelineStages) *DeliveryPipelineSerialPipelineStages {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for DeliveryPipelineSerialPipelineStages while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.TargetId, nw.TargetId) {
		nw.TargetId = des.TargetId
	}
	if dcl.StringArrayCanonicalize(des.Profiles, nw.Profiles) {
		nw.Profiles = des.Profiles
	}

	return nw
}

func canonicalizeNewDeliveryPipelineSerialPipelineStagesSet(c *Client, des, nw []DeliveryPipelineSerialPipelineStages) []DeliveryPipelineSerialPipelineStages {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []DeliveryPipelineSerialPipelineStages
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareDeliveryPipelineSerialPipelineStagesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewDeliveryPipelineSerialPipelineStages(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewDeliveryPipelineSerialPipelineStagesSlice(c *Client, des, nw []DeliveryPipelineSerialPipelineStages) []DeliveryPipelineSerialPipelineStages {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []DeliveryPipelineSerialPipelineStages
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewDeliveryPipelineSerialPipelineStages(c, &d, &n))
	}

	return items
}

func canonicalizeDeliveryPipelineCondition(des, initial *DeliveryPipelineCondition, opts ...dcl.ApplyOption) *DeliveryPipelineCondition {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &DeliveryPipelineCondition{}

	cDes.PipelineReadyCondition = canonicalizeDeliveryPipelineConditionPipelineReadyCondition(des.PipelineReadyCondition, initial.PipelineReadyCondition, opts...)
	cDes.TargetsPresentCondition = canonicalizeDeliveryPipelineConditionTargetsPresentCondition(des.TargetsPresentCondition, initial.TargetsPresentCondition, opts...)

	return cDes
}

func canonicalizeDeliveryPipelineConditionSlice(des, initial []DeliveryPipelineCondition, opts ...dcl.ApplyOption) []DeliveryPipelineCondition {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]DeliveryPipelineCondition, 0, len(des))
		for _, d := range des {
			cd := canonicalizeDeliveryPipelineCondition(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]DeliveryPipelineCondition, 0, len(des))
	for i, d := range des {
		cd := canonicalizeDeliveryPipelineCondition(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewDeliveryPipelineCondition(c *Client, des, nw *DeliveryPipelineCondition) *DeliveryPipelineCondition {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for DeliveryPipelineCondition while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.PipelineReadyCondition = canonicalizeNewDeliveryPipelineConditionPipelineReadyCondition(c, des.PipelineReadyCondition, nw.PipelineReadyCondition)
	nw.TargetsPresentCondition = canonicalizeNewDeliveryPipelineConditionTargetsPresentCondition(c, des.TargetsPresentCondition, nw.TargetsPresentCondition)

	return nw
}

func canonicalizeNewDeliveryPipelineConditionSet(c *Client, des, nw []DeliveryPipelineCondition) []DeliveryPipelineCondition {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []DeliveryPipelineCondition
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareDeliveryPipelineConditionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewDeliveryPipelineCondition(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewDeliveryPipelineConditionSlice(c *Client, des, nw []DeliveryPipelineCondition) []DeliveryPipelineCondition {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []DeliveryPipelineCondition
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewDeliveryPipelineCondition(c, &d, &n))
	}

	return items
}

func canonicalizeDeliveryPipelineConditionPipelineReadyCondition(des, initial *DeliveryPipelineConditionPipelineReadyCondition, opts ...dcl.ApplyOption) *DeliveryPipelineConditionPipelineReadyCondition {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &DeliveryPipelineConditionPipelineReadyCondition{}

	if dcl.BoolCanonicalize(des.Status, initial.Status) || dcl.IsZeroValue(des.Status) {
		cDes.Status = initial.Status
	} else {
		cDes.Status = des.Status
	}
	if dcl.IsZeroValue(des.UpdateTime) || (dcl.IsEmptyValueIndirect(des.UpdateTime) && dcl.IsEmptyValueIndirect(initial.UpdateTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.UpdateTime = initial.UpdateTime
	} else {
		cDes.UpdateTime = des.UpdateTime
	}

	return cDes
}

func canonicalizeDeliveryPipelineConditionPipelineReadyConditionSlice(des, initial []DeliveryPipelineConditionPipelineReadyCondition, opts ...dcl.ApplyOption) []DeliveryPipelineConditionPipelineReadyCondition {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]DeliveryPipelineConditionPipelineReadyCondition, 0, len(des))
		for _, d := range des {
			cd := canonicalizeDeliveryPipelineConditionPipelineReadyCondition(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]DeliveryPipelineConditionPipelineReadyCondition, 0, len(des))
	for i, d := range des {
		cd := canonicalizeDeliveryPipelineConditionPipelineReadyCondition(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewDeliveryPipelineConditionPipelineReadyCondition(c *Client, des, nw *DeliveryPipelineConditionPipelineReadyCondition) *DeliveryPipelineConditionPipelineReadyCondition {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for DeliveryPipelineConditionPipelineReadyCondition while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.Status, nw.Status) {
		nw.Status = des.Status
	}

	return nw
}

func canonicalizeNewDeliveryPipelineConditionPipelineReadyConditionSet(c *Client, des, nw []DeliveryPipelineConditionPipelineReadyCondition) []DeliveryPipelineConditionPipelineReadyCondition {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []DeliveryPipelineConditionPipelineReadyCondition
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareDeliveryPipelineConditionPipelineReadyConditionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewDeliveryPipelineConditionPipelineReadyCondition(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewDeliveryPipelineConditionPipelineReadyConditionSlice(c *Client, des, nw []DeliveryPipelineConditionPipelineReadyCondition) []DeliveryPipelineConditionPipelineReadyCondition {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []DeliveryPipelineConditionPipelineReadyCondition
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewDeliveryPipelineConditionPipelineReadyCondition(c, &d, &n))
	}

	return items
}

func canonicalizeDeliveryPipelineConditionTargetsPresentCondition(des, initial *DeliveryPipelineConditionTargetsPresentCondition, opts ...dcl.ApplyOption) *DeliveryPipelineConditionTargetsPresentCondition {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &DeliveryPipelineConditionTargetsPresentCondition{}

	if dcl.BoolCanonicalize(des.Status, initial.Status) || dcl.IsZeroValue(des.Status) {
		cDes.Status = initial.Status
	} else {
		cDes.Status = des.Status
	}
	if dcl.StringArrayCanonicalize(des.MissingTargets, initial.MissingTargets) {
		cDes.MissingTargets = initial.MissingTargets
	} else {
		cDes.MissingTargets = des.MissingTargets
	}
	if dcl.IsZeroValue(des.UpdateTime) || (dcl.IsEmptyValueIndirect(des.UpdateTime) && dcl.IsEmptyValueIndirect(initial.UpdateTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.UpdateTime = initial.UpdateTime
	} else {
		cDes.UpdateTime = des.UpdateTime
	}

	return cDes
}

func canonicalizeDeliveryPipelineConditionTargetsPresentConditionSlice(des, initial []DeliveryPipelineConditionTargetsPresentCondition, opts ...dcl.ApplyOption) []DeliveryPipelineConditionTargetsPresentCondition {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]DeliveryPipelineConditionTargetsPresentCondition, 0, len(des))
		for _, d := range des {
			cd := canonicalizeDeliveryPipelineConditionTargetsPresentCondition(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]DeliveryPipelineConditionTargetsPresentCondition, 0, len(des))
	for i, d := range des {
		cd := canonicalizeDeliveryPipelineConditionTargetsPresentCondition(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewDeliveryPipelineConditionTargetsPresentCondition(c *Client, des, nw *DeliveryPipelineConditionTargetsPresentCondition) *DeliveryPipelineConditionTargetsPresentCondition {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for DeliveryPipelineConditionTargetsPresentCondition while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.Status, nw.Status) {
		nw.Status = des.Status
	}
	if dcl.StringArrayCanonicalize(des.MissingTargets, nw.MissingTargets) {
		nw.MissingTargets = des.MissingTargets
	}

	return nw
}

func canonicalizeNewDeliveryPipelineConditionTargetsPresentConditionSet(c *Client, des, nw []DeliveryPipelineConditionTargetsPresentCondition) []DeliveryPipelineConditionTargetsPresentCondition {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []DeliveryPipelineConditionTargetsPresentCondition
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareDeliveryPipelineConditionTargetsPresentConditionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewDeliveryPipelineConditionTargetsPresentCondition(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewDeliveryPipelineConditionTargetsPresentConditionSlice(c *Client, des, nw []DeliveryPipelineConditionTargetsPresentCondition) []DeliveryPipelineConditionTargetsPresentCondition {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []DeliveryPipelineConditionTargetsPresentCondition
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewDeliveryPipelineConditionTargetsPresentCondition(c, &d, &n))
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
func diffDeliveryPipeline(c *Client, desired, actual *DeliveryPipeline, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Annotations, actual.Annotations, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("Annotations")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.SerialPipeline, actual.SerialPipeline, dcl.DiffInfo{ObjectFunction: compareDeliveryPipelineSerialPipelineNewStyle, EmptyObject: EmptyDeliveryPipelineSerialPipeline, OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("SerialPipeline")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Condition, actual.Condition, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareDeliveryPipelineConditionNewStyle, EmptyObject: EmptyDeliveryPipelineCondition, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Condition")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Suspended, actual.Suspended, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("Suspended")); len(ds) != 0 || err != nil {
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
func compareDeliveryPipelineSerialPipelineNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*DeliveryPipelineSerialPipeline)
	if !ok {
		desiredNotPointer, ok := d.(DeliveryPipelineSerialPipeline)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineSerialPipeline or *DeliveryPipelineSerialPipeline", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*DeliveryPipelineSerialPipeline)
	if !ok {
		actualNotPointer, ok := a.(DeliveryPipelineSerialPipeline)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineSerialPipeline", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Stages, actual.Stages, dcl.DiffInfo{ObjectFunction: compareDeliveryPipelineSerialPipelineStagesNewStyle, EmptyObject: EmptyDeliveryPipelineSerialPipelineStages, OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("Stages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareDeliveryPipelineSerialPipelineStagesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*DeliveryPipelineSerialPipelineStages)
	if !ok {
		desiredNotPointer, ok := d.(DeliveryPipelineSerialPipelineStages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineSerialPipelineStages or *DeliveryPipelineSerialPipelineStages", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*DeliveryPipelineSerialPipelineStages)
	if !ok {
		actualNotPointer, ok := a.(DeliveryPipelineSerialPipelineStages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineSerialPipelineStages", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.TargetId, actual.TargetId, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("TargetId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Profiles, actual.Profiles, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("Profiles")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareDeliveryPipelineConditionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*DeliveryPipelineCondition)
	if !ok {
		desiredNotPointer, ok := d.(DeliveryPipelineCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineCondition or *DeliveryPipelineCondition", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*DeliveryPipelineCondition)
	if !ok {
		actualNotPointer, ok := a.(DeliveryPipelineCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineCondition", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.PipelineReadyCondition, actual.PipelineReadyCondition, dcl.DiffInfo{ObjectFunction: compareDeliveryPipelineConditionPipelineReadyConditionNewStyle, EmptyObject: EmptyDeliveryPipelineConditionPipelineReadyCondition, OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("PipelineReadyCondition")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TargetsPresentCondition, actual.TargetsPresentCondition, dcl.DiffInfo{ObjectFunction: compareDeliveryPipelineConditionTargetsPresentConditionNewStyle, EmptyObject: EmptyDeliveryPipelineConditionTargetsPresentCondition, OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("TargetsPresentCondition")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareDeliveryPipelineConditionPipelineReadyConditionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*DeliveryPipelineConditionPipelineReadyCondition)
	if !ok {
		desiredNotPointer, ok := d.(DeliveryPipelineConditionPipelineReadyCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineConditionPipelineReadyCondition or *DeliveryPipelineConditionPipelineReadyCondition", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*DeliveryPipelineConditionPipelineReadyCondition)
	if !ok {
		actualNotPointer, ok := a.(DeliveryPipelineConditionPipelineReadyCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineConditionPipelineReadyCondition", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Status, actual.Status, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("Status")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareDeliveryPipelineConditionTargetsPresentConditionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*DeliveryPipelineConditionTargetsPresentCondition)
	if !ok {
		desiredNotPointer, ok := d.(DeliveryPipelineConditionTargetsPresentCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineConditionTargetsPresentCondition or *DeliveryPipelineConditionTargetsPresentCondition", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*DeliveryPipelineConditionTargetsPresentCondition)
	if !ok {
		actualNotPointer, ok := a.(DeliveryPipelineConditionTargetsPresentCondition)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a DeliveryPipelineConditionTargetsPresentCondition", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Status, actual.Status, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("Status")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MissingTargets, actual.MissingTargets, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("MissingTargets")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateDeliveryPipelineUpdateDeliveryPipelineOperation")}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
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
func (r *DeliveryPipeline) urlNormalized() *DeliveryPipeline {
	normalized := dcl.Copy(*r).(DeliveryPipeline)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Etag = dcl.SelfLinkToName(r.Etag)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *DeliveryPipeline) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateDeliveryPipeline" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the DeliveryPipeline resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *DeliveryPipeline) marshal(c *Client) ([]byte, error) {
	m, err := expandDeliveryPipeline(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling DeliveryPipeline: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalDeliveryPipeline decodes JSON responses into the DeliveryPipeline resource schema.
func unmarshalDeliveryPipeline(b []byte, c *Client, res *DeliveryPipeline) (*DeliveryPipeline, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapDeliveryPipeline(m, c, res)
}

func unmarshalMapDeliveryPipeline(m map[string]interface{}, c *Client, res *DeliveryPipeline) (*DeliveryPipeline, error) {

	flattened := flattenDeliveryPipeline(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandDeliveryPipeline expands DeliveryPipeline into a JSON request object.
func expandDeliveryPipeline(c *Client, f *DeliveryPipeline) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v := f.Annotations; dcl.ValueShouldBeSent(v) {
		m["annotations"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v, err := expandDeliveryPipelineSerialPipeline(c, f.SerialPipeline, res); err != nil {
		return nil, fmt.Errorf("error expanding SerialPipeline into serialPipeline: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["serialPipeline"] = v
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
	if v := f.Suspended; dcl.ValueShouldBeSent(v) {
		m["suspended"] = v
	}

	return m, nil
}

// flattenDeliveryPipeline flattens DeliveryPipeline from a JSON request object into the
// DeliveryPipeline type.
func flattenDeliveryPipeline(c *Client, i interface{}, res *DeliveryPipeline) *DeliveryPipeline {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &DeliveryPipeline{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Annotations = dcl.FlattenKeyValuePairs(m["annotations"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.SerialPipeline = flattenDeliveryPipelineSerialPipeline(c, m["serialPipeline"], res)
	resultRes.Condition = flattenDeliveryPipelineCondition(c, m["condition"], res)
	resultRes.Etag = dcl.FlattenString(m["etag"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.Suspended = dcl.FlattenBool(m["suspended"])

	return resultRes
}

// expandDeliveryPipelineSerialPipelineMap expands the contents of DeliveryPipelineSerialPipeline into a JSON
// request object.
func expandDeliveryPipelineSerialPipelineMap(c *Client, f map[string]DeliveryPipelineSerialPipeline, res *DeliveryPipeline) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandDeliveryPipelineSerialPipeline(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandDeliveryPipelineSerialPipelineSlice expands the contents of DeliveryPipelineSerialPipeline into a JSON
// request object.
func expandDeliveryPipelineSerialPipelineSlice(c *Client, f []DeliveryPipelineSerialPipeline, res *DeliveryPipeline) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandDeliveryPipelineSerialPipeline(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenDeliveryPipelineSerialPipelineMap flattens the contents of DeliveryPipelineSerialPipeline from a JSON
// response object.
func flattenDeliveryPipelineSerialPipelineMap(c *Client, i interface{}, res *DeliveryPipeline) map[string]DeliveryPipelineSerialPipeline {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]DeliveryPipelineSerialPipeline{}
	}

	if len(a) == 0 {
		return map[string]DeliveryPipelineSerialPipeline{}
	}

	items := make(map[string]DeliveryPipelineSerialPipeline)
	for k, item := range a {
		items[k] = *flattenDeliveryPipelineSerialPipeline(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenDeliveryPipelineSerialPipelineSlice flattens the contents of DeliveryPipelineSerialPipeline from a JSON
// response object.
func flattenDeliveryPipelineSerialPipelineSlice(c *Client, i interface{}, res *DeliveryPipeline) []DeliveryPipelineSerialPipeline {
	a, ok := i.([]interface{})
	if !ok {
		return []DeliveryPipelineSerialPipeline{}
	}

	if len(a) == 0 {
		return []DeliveryPipelineSerialPipeline{}
	}

	items := make([]DeliveryPipelineSerialPipeline, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenDeliveryPipelineSerialPipeline(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandDeliveryPipelineSerialPipeline expands an instance of DeliveryPipelineSerialPipeline into a JSON
// request object.
func expandDeliveryPipelineSerialPipeline(c *Client, f *DeliveryPipelineSerialPipeline, res *DeliveryPipeline) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandDeliveryPipelineSerialPipelineStagesSlice(c, f.Stages, res); err != nil {
		return nil, fmt.Errorf("error expanding Stages into stages: %w", err)
	} else if v != nil {
		m["stages"] = v
	}

	return m, nil
}

// flattenDeliveryPipelineSerialPipeline flattens an instance of DeliveryPipelineSerialPipeline from a JSON
// response object.
func flattenDeliveryPipelineSerialPipeline(c *Client, i interface{}, res *DeliveryPipeline) *DeliveryPipelineSerialPipeline {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &DeliveryPipelineSerialPipeline{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyDeliveryPipelineSerialPipeline
	}
	r.Stages = flattenDeliveryPipelineSerialPipelineStagesSlice(c, m["stages"], res)

	return r
}

// expandDeliveryPipelineSerialPipelineStagesMap expands the contents of DeliveryPipelineSerialPipelineStages into a JSON
// request object.
func expandDeliveryPipelineSerialPipelineStagesMap(c *Client, f map[string]DeliveryPipelineSerialPipelineStages, res *DeliveryPipeline) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandDeliveryPipelineSerialPipelineStages(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandDeliveryPipelineSerialPipelineStagesSlice expands the contents of DeliveryPipelineSerialPipelineStages into a JSON
// request object.
func expandDeliveryPipelineSerialPipelineStagesSlice(c *Client, f []DeliveryPipelineSerialPipelineStages, res *DeliveryPipeline) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandDeliveryPipelineSerialPipelineStages(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenDeliveryPipelineSerialPipelineStagesMap flattens the contents of DeliveryPipelineSerialPipelineStages from a JSON
// response object.
func flattenDeliveryPipelineSerialPipelineStagesMap(c *Client, i interface{}, res *DeliveryPipeline) map[string]DeliveryPipelineSerialPipelineStages {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]DeliveryPipelineSerialPipelineStages{}
	}

	if len(a) == 0 {
		return map[string]DeliveryPipelineSerialPipelineStages{}
	}

	items := make(map[string]DeliveryPipelineSerialPipelineStages)
	for k, item := range a {
		items[k] = *flattenDeliveryPipelineSerialPipelineStages(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenDeliveryPipelineSerialPipelineStagesSlice flattens the contents of DeliveryPipelineSerialPipelineStages from a JSON
// response object.
func flattenDeliveryPipelineSerialPipelineStagesSlice(c *Client, i interface{}, res *DeliveryPipeline) []DeliveryPipelineSerialPipelineStages {
	a, ok := i.([]interface{})
	if !ok {
		return []DeliveryPipelineSerialPipelineStages{}
	}

	if len(a) == 0 {
		return []DeliveryPipelineSerialPipelineStages{}
	}

	items := make([]DeliveryPipelineSerialPipelineStages, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenDeliveryPipelineSerialPipelineStages(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandDeliveryPipelineSerialPipelineStages expands an instance of DeliveryPipelineSerialPipelineStages into a JSON
// request object.
func expandDeliveryPipelineSerialPipelineStages(c *Client, f *DeliveryPipelineSerialPipelineStages, res *DeliveryPipeline) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.TargetId; !dcl.IsEmptyValueIndirect(v) {
		m["targetId"] = v
	}
	if v := f.Profiles; v != nil {
		m["profiles"] = v
	}

	return m, nil
}

// flattenDeliveryPipelineSerialPipelineStages flattens an instance of DeliveryPipelineSerialPipelineStages from a JSON
// response object.
func flattenDeliveryPipelineSerialPipelineStages(c *Client, i interface{}, res *DeliveryPipeline) *DeliveryPipelineSerialPipelineStages {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &DeliveryPipelineSerialPipelineStages{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyDeliveryPipelineSerialPipelineStages
	}
	r.TargetId = dcl.FlattenString(m["targetId"])
	r.Profiles = dcl.FlattenStringSlice(m["profiles"])

	return r
}

// expandDeliveryPipelineConditionMap expands the contents of DeliveryPipelineCondition into a JSON
// request object.
func expandDeliveryPipelineConditionMap(c *Client, f map[string]DeliveryPipelineCondition, res *DeliveryPipeline) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandDeliveryPipelineCondition(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandDeliveryPipelineConditionSlice expands the contents of DeliveryPipelineCondition into a JSON
// request object.
func expandDeliveryPipelineConditionSlice(c *Client, f []DeliveryPipelineCondition, res *DeliveryPipeline) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandDeliveryPipelineCondition(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenDeliveryPipelineConditionMap flattens the contents of DeliveryPipelineCondition from a JSON
// response object.
func flattenDeliveryPipelineConditionMap(c *Client, i interface{}, res *DeliveryPipeline) map[string]DeliveryPipelineCondition {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]DeliveryPipelineCondition{}
	}

	if len(a) == 0 {
		return map[string]DeliveryPipelineCondition{}
	}

	items := make(map[string]DeliveryPipelineCondition)
	for k, item := range a {
		items[k] = *flattenDeliveryPipelineCondition(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenDeliveryPipelineConditionSlice flattens the contents of DeliveryPipelineCondition from a JSON
// response object.
func flattenDeliveryPipelineConditionSlice(c *Client, i interface{}, res *DeliveryPipeline) []DeliveryPipelineCondition {
	a, ok := i.([]interface{})
	if !ok {
		return []DeliveryPipelineCondition{}
	}

	if len(a) == 0 {
		return []DeliveryPipelineCondition{}
	}

	items := make([]DeliveryPipelineCondition, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenDeliveryPipelineCondition(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandDeliveryPipelineCondition expands an instance of DeliveryPipelineCondition into a JSON
// request object.
func expandDeliveryPipelineCondition(c *Client, f *DeliveryPipelineCondition, res *DeliveryPipeline) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandDeliveryPipelineConditionPipelineReadyCondition(c, f.PipelineReadyCondition, res); err != nil {
		return nil, fmt.Errorf("error expanding PipelineReadyCondition into pipelineReadyCondition: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pipelineReadyCondition"] = v
	}
	if v, err := expandDeliveryPipelineConditionTargetsPresentCondition(c, f.TargetsPresentCondition, res); err != nil {
		return nil, fmt.Errorf("error expanding TargetsPresentCondition into targetsPresentCondition: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["targetsPresentCondition"] = v
	}

	return m, nil
}

// flattenDeliveryPipelineCondition flattens an instance of DeliveryPipelineCondition from a JSON
// response object.
func flattenDeliveryPipelineCondition(c *Client, i interface{}, res *DeliveryPipeline) *DeliveryPipelineCondition {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &DeliveryPipelineCondition{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyDeliveryPipelineCondition
	}
	r.PipelineReadyCondition = flattenDeliveryPipelineConditionPipelineReadyCondition(c, m["pipelineReadyCondition"], res)
	r.TargetsPresentCondition = flattenDeliveryPipelineConditionTargetsPresentCondition(c, m["targetsPresentCondition"], res)

	return r
}

// expandDeliveryPipelineConditionPipelineReadyConditionMap expands the contents of DeliveryPipelineConditionPipelineReadyCondition into a JSON
// request object.
func expandDeliveryPipelineConditionPipelineReadyConditionMap(c *Client, f map[string]DeliveryPipelineConditionPipelineReadyCondition, res *DeliveryPipeline) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandDeliveryPipelineConditionPipelineReadyCondition(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandDeliveryPipelineConditionPipelineReadyConditionSlice expands the contents of DeliveryPipelineConditionPipelineReadyCondition into a JSON
// request object.
func expandDeliveryPipelineConditionPipelineReadyConditionSlice(c *Client, f []DeliveryPipelineConditionPipelineReadyCondition, res *DeliveryPipeline) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandDeliveryPipelineConditionPipelineReadyCondition(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenDeliveryPipelineConditionPipelineReadyConditionMap flattens the contents of DeliveryPipelineConditionPipelineReadyCondition from a JSON
// response object.
func flattenDeliveryPipelineConditionPipelineReadyConditionMap(c *Client, i interface{}, res *DeliveryPipeline) map[string]DeliveryPipelineConditionPipelineReadyCondition {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]DeliveryPipelineConditionPipelineReadyCondition{}
	}

	if len(a) == 0 {
		return map[string]DeliveryPipelineConditionPipelineReadyCondition{}
	}

	items := make(map[string]DeliveryPipelineConditionPipelineReadyCondition)
	for k, item := range a {
		items[k] = *flattenDeliveryPipelineConditionPipelineReadyCondition(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenDeliveryPipelineConditionPipelineReadyConditionSlice flattens the contents of DeliveryPipelineConditionPipelineReadyCondition from a JSON
// response object.
func flattenDeliveryPipelineConditionPipelineReadyConditionSlice(c *Client, i interface{}, res *DeliveryPipeline) []DeliveryPipelineConditionPipelineReadyCondition {
	a, ok := i.([]interface{})
	if !ok {
		return []DeliveryPipelineConditionPipelineReadyCondition{}
	}

	if len(a) == 0 {
		return []DeliveryPipelineConditionPipelineReadyCondition{}
	}

	items := make([]DeliveryPipelineConditionPipelineReadyCondition, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenDeliveryPipelineConditionPipelineReadyCondition(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandDeliveryPipelineConditionPipelineReadyCondition expands an instance of DeliveryPipelineConditionPipelineReadyCondition into a JSON
// request object.
func expandDeliveryPipelineConditionPipelineReadyCondition(c *Client, f *DeliveryPipelineConditionPipelineReadyCondition, res *DeliveryPipeline) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Status; !dcl.IsEmptyValueIndirect(v) {
		m["status"] = v
	}
	if v := f.UpdateTime; !dcl.IsEmptyValueIndirect(v) {
		m["updateTime"] = v
	}

	return m, nil
}

// flattenDeliveryPipelineConditionPipelineReadyCondition flattens an instance of DeliveryPipelineConditionPipelineReadyCondition from a JSON
// response object.
func flattenDeliveryPipelineConditionPipelineReadyCondition(c *Client, i interface{}, res *DeliveryPipeline) *DeliveryPipelineConditionPipelineReadyCondition {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &DeliveryPipelineConditionPipelineReadyCondition{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyDeliveryPipelineConditionPipelineReadyCondition
	}
	r.Status = dcl.FlattenBool(m["status"])
	r.UpdateTime = dcl.FlattenString(m["updateTime"])

	return r
}

// expandDeliveryPipelineConditionTargetsPresentConditionMap expands the contents of DeliveryPipelineConditionTargetsPresentCondition into a JSON
// request object.
func expandDeliveryPipelineConditionTargetsPresentConditionMap(c *Client, f map[string]DeliveryPipelineConditionTargetsPresentCondition, res *DeliveryPipeline) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandDeliveryPipelineConditionTargetsPresentCondition(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandDeliveryPipelineConditionTargetsPresentConditionSlice expands the contents of DeliveryPipelineConditionTargetsPresentCondition into a JSON
// request object.
func expandDeliveryPipelineConditionTargetsPresentConditionSlice(c *Client, f []DeliveryPipelineConditionTargetsPresentCondition, res *DeliveryPipeline) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandDeliveryPipelineConditionTargetsPresentCondition(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenDeliveryPipelineConditionTargetsPresentConditionMap flattens the contents of DeliveryPipelineConditionTargetsPresentCondition from a JSON
// response object.
func flattenDeliveryPipelineConditionTargetsPresentConditionMap(c *Client, i interface{}, res *DeliveryPipeline) map[string]DeliveryPipelineConditionTargetsPresentCondition {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]DeliveryPipelineConditionTargetsPresentCondition{}
	}

	if len(a) == 0 {
		return map[string]DeliveryPipelineConditionTargetsPresentCondition{}
	}

	items := make(map[string]DeliveryPipelineConditionTargetsPresentCondition)
	for k, item := range a {
		items[k] = *flattenDeliveryPipelineConditionTargetsPresentCondition(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenDeliveryPipelineConditionTargetsPresentConditionSlice flattens the contents of DeliveryPipelineConditionTargetsPresentCondition from a JSON
// response object.
func flattenDeliveryPipelineConditionTargetsPresentConditionSlice(c *Client, i interface{}, res *DeliveryPipeline) []DeliveryPipelineConditionTargetsPresentCondition {
	a, ok := i.([]interface{})
	if !ok {
		return []DeliveryPipelineConditionTargetsPresentCondition{}
	}

	if len(a) == 0 {
		return []DeliveryPipelineConditionTargetsPresentCondition{}
	}

	items := make([]DeliveryPipelineConditionTargetsPresentCondition, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenDeliveryPipelineConditionTargetsPresentCondition(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandDeliveryPipelineConditionTargetsPresentCondition expands an instance of DeliveryPipelineConditionTargetsPresentCondition into a JSON
// request object.
func expandDeliveryPipelineConditionTargetsPresentCondition(c *Client, f *DeliveryPipelineConditionTargetsPresentCondition, res *DeliveryPipeline) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Status; !dcl.IsEmptyValueIndirect(v) {
		m["status"] = v
	}
	if v := f.MissingTargets; v != nil {
		m["missingTargets"] = v
	}
	if v := f.UpdateTime; !dcl.IsEmptyValueIndirect(v) {
		m["updateTime"] = v
	}

	return m, nil
}

// flattenDeliveryPipelineConditionTargetsPresentCondition flattens an instance of DeliveryPipelineConditionTargetsPresentCondition from a JSON
// response object.
func flattenDeliveryPipelineConditionTargetsPresentCondition(c *Client, i interface{}, res *DeliveryPipeline) *DeliveryPipelineConditionTargetsPresentCondition {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &DeliveryPipelineConditionTargetsPresentCondition{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyDeliveryPipelineConditionTargetsPresentCondition
	}
	r.Status = dcl.FlattenBool(m["status"])
	r.MissingTargets = dcl.FlattenStringSlice(m["missingTargets"])
	r.UpdateTime = dcl.FlattenString(m["updateTime"])

	return r
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *DeliveryPipeline) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalDeliveryPipeline(b, c, r)
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

type deliveryPipelineDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         deliveryPipelineApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToDeliveryPipelineDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]deliveryPipelineDiff, error) {
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
	var diffs []deliveryPipelineDiff
	// For each operation name, create a deliveryPipelineDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := deliveryPipelineDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToDeliveryPipelineApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToDeliveryPipelineApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (deliveryPipelineApiOperation, error) {
	switch opName {

	case "updateDeliveryPipelineUpdateDeliveryPipelineOperation":
		return &updateDeliveryPipelineUpdateDeliveryPipelineOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractDeliveryPipelineFields(r *DeliveryPipeline) error {
	vSerialPipeline := r.SerialPipeline
	if vSerialPipeline == nil {
		// note: explicitly not the empty object.
		vSerialPipeline = &DeliveryPipelineSerialPipeline{}
	}
	if err := extractDeliveryPipelineSerialPipelineFields(r, vSerialPipeline); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSerialPipeline) {
		r.SerialPipeline = vSerialPipeline
	}
	vCondition := r.Condition
	if vCondition == nil {
		// note: explicitly not the empty object.
		vCondition = &DeliveryPipelineCondition{}
	}
	if err := extractDeliveryPipelineConditionFields(r, vCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCondition) {
		r.Condition = vCondition
	}
	return nil
}
func extractDeliveryPipelineSerialPipelineFields(r *DeliveryPipeline, o *DeliveryPipelineSerialPipeline) error {
	return nil
}
func extractDeliveryPipelineSerialPipelineStagesFields(r *DeliveryPipeline, o *DeliveryPipelineSerialPipelineStages) error {
	return nil
}
func extractDeliveryPipelineConditionFields(r *DeliveryPipeline, o *DeliveryPipelineCondition) error {
	vPipelineReadyCondition := o.PipelineReadyCondition
	if vPipelineReadyCondition == nil {
		// note: explicitly not the empty object.
		vPipelineReadyCondition = &DeliveryPipelineConditionPipelineReadyCondition{}
	}
	if err := extractDeliveryPipelineConditionPipelineReadyConditionFields(r, vPipelineReadyCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPipelineReadyCondition) {
		o.PipelineReadyCondition = vPipelineReadyCondition
	}
	vTargetsPresentCondition := o.TargetsPresentCondition
	if vTargetsPresentCondition == nil {
		// note: explicitly not the empty object.
		vTargetsPresentCondition = &DeliveryPipelineConditionTargetsPresentCondition{}
	}
	if err := extractDeliveryPipelineConditionTargetsPresentConditionFields(r, vTargetsPresentCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vTargetsPresentCondition) {
		o.TargetsPresentCondition = vTargetsPresentCondition
	}
	return nil
}
func extractDeliveryPipelineConditionPipelineReadyConditionFields(r *DeliveryPipeline, o *DeliveryPipelineConditionPipelineReadyCondition) error {
	return nil
}
func extractDeliveryPipelineConditionTargetsPresentConditionFields(r *DeliveryPipeline, o *DeliveryPipelineConditionTargetsPresentCondition) error {
	return nil
}

func postReadExtractDeliveryPipelineFields(r *DeliveryPipeline) error {
	vSerialPipeline := r.SerialPipeline
	if vSerialPipeline == nil {
		// note: explicitly not the empty object.
		vSerialPipeline = &DeliveryPipelineSerialPipeline{}
	}
	if err := postReadExtractDeliveryPipelineSerialPipelineFields(r, vSerialPipeline); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSerialPipeline) {
		r.SerialPipeline = vSerialPipeline
	}
	vCondition := r.Condition
	if vCondition == nil {
		// note: explicitly not the empty object.
		vCondition = &DeliveryPipelineCondition{}
	}
	if err := postReadExtractDeliveryPipelineConditionFields(r, vCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCondition) {
		r.Condition = vCondition
	}
	return nil
}
func postReadExtractDeliveryPipelineSerialPipelineFields(r *DeliveryPipeline, o *DeliveryPipelineSerialPipeline) error {
	return nil
}
func postReadExtractDeliveryPipelineSerialPipelineStagesFields(r *DeliveryPipeline, o *DeliveryPipelineSerialPipelineStages) error {
	return nil
}
func postReadExtractDeliveryPipelineConditionFields(r *DeliveryPipeline, o *DeliveryPipelineCondition) error {
	vPipelineReadyCondition := o.PipelineReadyCondition
	if vPipelineReadyCondition == nil {
		// note: explicitly not the empty object.
		vPipelineReadyCondition = &DeliveryPipelineConditionPipelineReadyCondition{}
	}
	if err := extractDeliveryPipelineConditionPipelineReadyConditionFields(r, vPipelineReadyCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPipelineReadyCondition) {
		o.PipelineReadyCondition = vPipelineReadyCondition
	}
	vTargetsPresentCondition := o.TargetsPresentCondition
	if vTargetsPresentCondition == nil {
		// note: explicitly not the empty object.
		vTargetsPresentCondition = &DeliveryPipelineConditionTargetsPresentCondition{}
	}
	if err := extractDeliveryPipelineConditionTargetsPresentConditionFields(r, vTargetsPresentCondition); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vTargetsPresentCondition) {
		o.TargetsPresentCondition = vTargetsPresentCondition
	}
	return nil
}
func postReadExtractDeliveryPipelineConditionPipelineReadyConditionFields(r *DeliveryPipeline, o *DeliveryPipelineConditionPipelineReadyCondition) error {
	return nil
}
func postReadExtractDeliveryPipelineConditionTargetsPresentConditionFields(r *DeliveryPipeline, o *DeliveryPipelineConditionTargetsPresentCondition) error {
	return nil
}
