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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type DeliveryPipeline struct {
	Name           *string                         `json:"name"`
	Uid            *string                         `json:"uid"`
	Description    *string                         `json:"description"`
	Annotations    map[string]string               `json:"annotations"`
	Labels         map[string]string               `json:"labels"`
	CreateTime     *string                         `json:"createTime"`
	UpdateTime     *string                         `json:"updateTime"`
	SerialPipeline *DeliveryPipelineSerialPipeline `json:"serialPipeline"`
	Condition      *DeliveryPipelineCondition      `json:"condition"`
	Etag           *string                         `json:"etag"`
	Project        *string                         `json:"project"`
	Location       *string                         `json:"location"`
	Suspended      *bool                           `json:"suspended"`
}

func (r *DeliveryPipeline) String() string {
	return dcl.SprintResource(r)
}

type DeliveryPipelineSerialPipeline struct {
	empty  bool                                   `json:"-"`
	Stages []DeliveryPipelineSerialPipelineStages `json:"stages"`
}

type jsonDeliveryPipelineSerialPipeline DeliveryPipelineSerialPipeline

func (r *DeliveryPipelineSerialPipeline) UnmarshalJSON(data []byte) error {
	var res jsonDeliveryPipelineSerialPipeline
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDeliveryPipelineSerialPipeline
	} else {

		r.Stages = res.Stages

	}
	return nil
}

// This object is used to assert a desired state where this DeliveryPipelineSerialPipeline is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDeliveryPipelineSerialPipeline *DeliveryPipelineSerialPipeline = &DeliveryPipelineSerialPipeline{empty: true}

func (r *DeliveryPipelineSerialPipeline) Empty() bool {
	return r.empty
}

func (r *DeliveryPipelineSerialPipeline) String() string {
	return dcl.SprintResource(r)
}

func (r *DeliveryPipelineSerialPipeline) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DeliveryPipelineSerialPipelineStages struct {
	empty    bool     `json:"-"`
	TargetId *string  `json:"targetId"`
	Profiles []string `json:"profiles"`
}

type jsonDeliveryPipelineSerialPipelineStages DeliveryPipelineSerialPipelineStages

func (r *DeliveryPipelineSerialPipelineStages) UnmarshalJSON(data []byte) error {
	var res jsonDeliveryPipelineSerialPipelineStages
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDeliveryPipelineSerialPipelineStages
	} else {

		r.TargetId = res.TargetId

		r.Profiles = res.Profiles

	}
	return nil
}

// This object is used to assert a desired state where this DeliveryPipelineSerialPipelineStages is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDeliveryPipelineSerialPipelineStages *DeliveryPipelineSerialPipelineStages = &DeliveryPipelineSerialPipelineStages{empty: true}

func (r *DeliveryPipelineSerialPipelineStages) Empty() bool {
	return r.empty
}

func (r *DeliveryPipelineSerialPipelineStages) String() string {
	return dcl.SprintResource(r)
}

func (r *DeliveryPipelineSerialPipelineStages) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DeliveryPipelineCondition struct {
	empty                   bool                                              `json:"-"`
	PipelineReadyCondition  *DeliveryPipelineConditionPipelineReadyCondition  `json:"pipelineReadyCondition"`
	TargetsPresentCondition *DeliveryPipelineConditionTargetsPresentCondition `json:"targetsPresentCondition"`
}

type jsonDeliveryPipelineCondition DeliveryPipelineCondition

func (r *DeliveryPipelineCondition) UnmarshalJSON(data []byte) error {
	var res jsonDeliveryPipelineCondition
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDeliveryPipelineCondition
	} else {

		r.PipelineReadyCondition = res.PipelineReadyCondition

		r.TargetsPresentCondition = res.TargetsPresentCondition

	}
	return nil
}

// This object is used to assert a desired state where this DeliveryPipelineCondition is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDeliveryPipelineCondition *DeliveryPipelineCondition = &DeliveryPipelineCondition{empty: true}

func (r *DeliveryPipelineCondition) Empty() bool {
	return r.empty
}

func (r *DeliveryPipelineCondition) String() string {
	return dcl.SprintResource(r)
}

func (r *DeliveryPipelineCondition) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DeliveryPipelineConditionPipelineReadyCondition struct {
	empty      bool    `json:"-"`
	Status     *bool   `json:"status"`
	UpdateTime *string `json:"updateTime"`
}

type jsonDeliveryPipelineConditionPipelineReadyCondition DeliveryPipelineConditionPipelineReadyCondition

func (r *DeliveryPipelineConditionPipelineReadyCondition) UnmarshalJSON(data []byte) error {
	var res jsonDeliveryPipelineConditionPipelineReadyCondition
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDeliveryPipelineConditionPipelineReadyCondition
	} else {

		r.Status = res.Status

		r.UpdateTime = res.UpdateTime

	}
	return nil
}

// This object is used to assert a desired state where this DeliveryPipelineConditionPipelineReadyCondition is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDeliveryPipelineConditionPipelineReadyCondition *DeliveryPipelineConditionPipelineReadyCondition = &DeliveryPipelineConditionPipelineReadyCondition{empty: true}

func (r *DeliveryPipelineConditionPipelineReadyCondition) Empty() bool {
	return r.empty
}

func (r *DeliveryPipelineConditionPipelineReadyCondition) String() string {
	return dcl.SprintResource(r)
}

func (r *DeliveryPipelineConditionPipelineReadyCondition) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type DeliveryPipelineConditionTargetsPresentCondition struct {
	empty          bool     `json:"-"`
	Status         *bool    `json:"status"`
	MissingTargets []string `json:"missingTargets"`
	UpdateTime     *string  `json:"updateTime"`
}

type jsonDeliveryPipelineConditionTargetsPresentCondition DeliveryPipelineConditionTargetsPresentCondition

func (r *DeliveryPipelineConditionTargetsPresentCondition) UnmarshalJSON(data []byte) error {
	var res jsonDeliveryPipelineConditionTargetsPresentCondition
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyDeliveryPipelineConditionTargetsPresentCondition
	} else {

		r.Status = res.Status

		r.MissingTargets = res.MissingTargets

		r.UpdateTime = res.UpdateTime

	}
	return nil
}

// This object is used to assert a desired state where this DeliveryPipelineConditionTargetsPresentCondition is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyDeliveryPipelineConditionTargetsPresentCondition *DeliveryPipelineConditionTargetsPresentCondition = &DeliveryPipelineConditionTargetsPresentCondition{empty: true}

func (r *DeliveryPipelineConditionTargetsPresentCondition) Empty() bool {
	return r.empty
}

func (r *DeliveryPipelineConditionTargetsPresentCondition) String() string {
	return dcl.SprintResource(r)
}

func (r *DeliveryPipelineConditionTargetsPresentCondition) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *DeliveryPipeline) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "clouddeploy",
		Type:    "DeliveryPipeline",
		Version: "clouddeploy",
	}
}

func (r *DeliveryPipeline) ID() (string, error) {
	if err := extractDeliveryPipelineFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":            dcl.ValueOrEmptyString(nr.Name),
		"uid":             dcl.ValueOrEmptyString(nr.Uid),
		"description":     dcl.ValueOrEmptyString(nr.Description),
		"annotations":     dcl.ValueOrEmptyString(nr.Annotations),
		"labels":          dcl.ValueOrEmptyString(nr.Labels),
		"create_time":     dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":     dcl.ValueOrEmptyString(nr.UpdateTime),
		"serial_pipeline": dcl.ValueOrEmptyString(nr.SerialPipeline),
		"condition":       dcl.ValueOrEmptyString(nr.Condition),
		"etag":            dcl.ValueOrEmptyString(nr.Etag),
		"project":         dcl.ValueOrEmptyString(nr.Project),
		"location":        dcl.ValueOrEmptyString(nr.Location),
		"suspended":       dcl.ValueOrEmptyString(nr.Suspended),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}", params), nil
}

const DeliveryPipelineMaxPage = -1

type DeliveryPipelineList struct {
	Items []*DeliveryPipeline

	nextToken string

	pageSize int32

	resource *DeliveryPipeline
}

func (l *DeliveryPipelineList) HasNext() bool {
	return l.nextToken != ""
}

func (l *DeliveryPipelineList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listDeliveryPipeline(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListDeliveryPipeline(ctx context.Context, project, location string) (*DeliveryPipelineList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListDeliveryPipelineWithMaxResults(ctx, project, location, DeliveryPipelineMaxPage)

}

func (c *Client) ListDeliveryPipelineWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*DeliveryPipelineList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &DeliveryPipeline{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listDeliveryPipeline(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &DeliveryPipelineList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetDeliveryPipeline(ctx context.Context, r *DeliveryPipeline) (*DeliveryPipeline, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractDeliveryPipelineFields(r)

	b, err := c.getDeliveryPipelineRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalDeliveryPipeline(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeDeliveryPipelineNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractDeliveryPipelineFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteDeliveryPipeline(ctx context.Context, r *DeliveryPipeline) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("DeliveryPipeline resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting DeliveryPipeline...")
	deleteOp := deleteDeliveryPipelineOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllDeliveryPipeline deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllDeliveryPipeline(ctx context.Context, project, location string, filter func(*DeliveryPipeline) bool) error {
	listObj, err := c.ListDeliveryPipeline(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllDeliveryPipeline(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllDeliveryPipeline(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyDeliveryPipeline(ctx context.Context, rawDesired *DeliveryPipeline, opts ...dcl.ApplyOption) (*DeliveryPipeline, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *DeliveryPipeline
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyDeliveryPipelineHelper(c, ctx, rawDesired, opts...)
		resultNewState = newState
		if err != nil {
			// If the error is 409, there is conflict in resource update.
			// Here we want to apply changes based on latest state.
			if dcl.IsConflictError(err) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		return nil, nil
	}, c.Config.RetryProvider)
	return resultNewState, err
}

func applyDeliveryPipelineHelper(c *Client, ctx context.Context, rawDesired *DeliveryPipeline, opts ...dcl.ApplyOption) (*DeliveryPipeline, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyDeliveryPipeline...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractDeliveryPipelineFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.deliveryPipelineDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToDeliveryPipelineDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	var create bool
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		if dcl.HasLifecycleParam(lp, dcl.BlockCreation) {
			return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Creation blocked by lifecycle params: %#v.", desired)}
		}
		create = true
	} else if dcl.HasLifecycleParam(lp, dcl.BlockAcquire) {
		return nil, dcl.ApplyInfeasibleError{
			Message: fmt.Sprintf("Resource already exists - apply blocked by lifecycle params: %#v.", initial),
		}
	} else {
		for _, d := range diffs {
			if d.RequiresRecreate {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) would require recreation", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}

	// 2.4 Imperative Request Planning
	var ops []deliveryPipelineApiOperation
	if create {
		ops = append(ops, &createDeliveryPipelineOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created plan: %#v", ops)

	// 2.5 Request Actuation
	for _, op := range ops {
		c.Config.Logger.InfoWithContextf(ctx, "Performing operation %T %+v", op, op)
		if err := op.do(ctx, desired, c); err != nil {
			c.Config.Logger.InfoWithContextf(ctx, "Failed operation %T %+v: %v", op, op, err)
			return nil, err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Finished operation %T %+v", op, op)
	}
	return applyDeliveryPipelineDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyDeliveryPipelineDiff(c *Client, ctx context.Context, desired *DeliveryPipeline, rawDesired *DeliveryPipeline, ops []deliveryPipelineApiOperation, opts ...dcl.ApplyOption) (*DeliveryPipeline, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetDeliveryPipeline(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createDeliveryPipelineOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapDeliveryPipeline(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeDeliveryPipelineNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeDeliveryPipelineNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeDeliveryPipelineDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractDeliveryPipelineFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractDeliveryPipelineFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffDeliveryPipeline(c, newDesired, newState)
	if err != nil {
		return newState, err
	}

	if len(newDiffs) == 0 {
		c.Config.Logger.InfoWithContext(ctx, "No diffs found. Apply was successful.")
	} else {
		c.Config.Logger.InfoWithContextf(ctx, "Found diffs: %v", newDiffs)
		diffMessages := make([]string, len(newDiffs))
		for i, d := range newDiffs {
			diffMessages[i] = fmt.Sprintf("%v", d)
		}
		return newState, dcl.DiffAfterApplyError{Diffs: diffMessages}
	}
	c.Config.Logger.InfoWithContext(ctx, "Done Apply.")
	return newState, nil
}

func (r *DeliveryPipeline) GetPolicy(basePath string) (string, string, *bytes.Buffer, error) {
	u := r.getPolicyURL(basePath)
	body := &bytes.Buffer{}
	u, err := dcl.AddQueryParams(u, map[string]string{"optionsRequestedPolicyVersion": fmt.Sprintf("%d", r.IAMPolicyVersion())})
	if err != nil {
		return "", "", nil, err
	}
	return u, "", body, nil
}
