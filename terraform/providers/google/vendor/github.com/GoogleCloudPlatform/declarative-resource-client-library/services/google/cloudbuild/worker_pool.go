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
package cloudbuild

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type WorkerPool struct {
	Name                *string                        `json:"name"`
	DisplayName         *string                        `json:"displayName"`
	Uid                 *string                        `json:"uid"`
	Annotations         map[string]string              `json:"annotations"`
	CreateTime          *string                        `json:"createTime"`
	UpdateTime          *string                        `json:"updateTime"`
	DeleteTime          *string                        `json:"deleteTime"`
	State               *WorkerPoolStateEnum           `json:"state"`
	PrivatePoolV1Config *WorkerPoolPrivatePoolV1Config `json:"privatePoolV1Config"`
	Etag                *string                        `json:"etag"`
	WorkerConfig        *WorkerPoolWorkerConfig        `json:"workerConfig"`
	NetworkConfig       *WorkerPoolNetworkConfig       `json:"networkConfig"`
	Project             *string                        `json:"project"`
	Location            *string                        `json:"location"`
}

func (r *WorkerPool) String() string {
	return dcl.SprintResource(r)
}

// The enum WorkerPoolStateEnum.
type WorkerPoolStateEnum string

// WorkerPoolStateEnumRef returns a *WorkerPoolStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func WorkerPoolStateEnumRef(s string) *WorkerPoolStateEnum {
	v := WorkerPoolStateEnum(s)
	return &v
}

func (v WorkerPoolStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "PENDING", "APPROVED", "REJECTED", "CANCELLED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "WorkerPoolStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum.
type WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum string

// WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef returns a *WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum with the value of string s
// If the empty string is provided, nil is returned.
func WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef(s string) *WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum {
	v := WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum(s)
	return &v
}

func (v WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"EGRESS_OPTION_UNSPECIFIED", "NO_PUBLIC_EGRESS", "PUBLIC_EGRESS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type WorkerPoolPrivatePoolV1Config struct {
	empty         bool                                        `json:"-"`
	WorkerConfig  *WorkerPoolPrivatePoolV1ConfigWorkerConfig  `json:"workerConfig"`
	NetworkConfig *WorkerPoolPrivatePoolV1ConfigNetworkConfig `json:"networkConfig"`
}

type jsonWorkerPoolPrivatePoolV1Config WorkerPoolPrivatePoolV1Config

func (r *WorkerPoolPrivatePoolV1Config) UnmarshalJSON(data []byte) error {
	var res jsonWorkerPoolPrivatePoolV1Config
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkerPoolPrivatePoolV1Config
	} else {

		r.WorkerConfig = res.WorkerConfig

		r.NetworkConfig = res.NetworkConfig

	}
	return nil
}

// This object is used to assert a desired state where this WorkerPoolPrivatePoolV1Config is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyWorkerPoolPrivatePoolV1Config *WorkerPoolPrivatePoolV1Config = &WorkerPoolPrivatePoolV1Config{empty: true}

func (r *WorkerPoolPrivatePoolV1Config) Empty() bool {
	return r.empty
}

func (r *WorkerPoolPrivatePoolV1Config) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkerPoolPrivatePoolV1Config) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkerPoolPrivatePoolV1ConfigWorkerConfig struct {
	empty       bool    `json:"-"`
	MachineType *string `json:"machineType"`
	DiskSizeGb  *int64  `json:"diskSizeGb"`
}

type jsonWorkerPoolPrivatePoolV1ConfigWorkerConfig WorkerPoolPrivatePoolV1ConfigWorkerConfig

func (r *WorkerPoolPrivatePoolV1ConfigWorkerConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkerPoolPrivatePoolV1ConfigWorkerConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkerPoolPrivatePoolV1ConfigWorkerConfig
	} else {

		r.MachineType = res.MachineType

		r.DiskSizeGb = res.DiskSizeGb

	}
	return nil
}

// This object is used to assert a desired state where this WorkerPoolPrivatePoolV1ConfigWorkerConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyWorkerPoolPrivatePoolV1ConfigWorkerConfig *WorkerPoolPrivatePoolV1ConfigWorkerConfig = &WorkerPoolPrivatePoolV1ConfigWorkerConfig{empty: true}

func (r *WorkerPoolPrivatePoolV1ConfigWorkerConfig) Empty() bool {
	return r.empty
}

func (r *WorkerPoolPrivatePoolV1ConfigWorkerConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkerPoolPrivatePoolV1ConfigWorkerConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkerPoolPrivatePoolV1ConfigNetworkConfig struct {
	empty                bool                                                        `json:"-"`
	PeeredNetwork        *string                                                     `json:"peeredNetwork"`
	PeeredNetworkIPRange *string                                                     `json:"peeredNetworkIPRange"`
	EgressOption         *WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum `json:"egressOption"`
}

type jsonWorkerPoolPrivatePoolV1ConfigNetworkConfig WorkerPoolPrivatePoolV1ConfigNetworkConfig

func (r *WorkerPoolPrivatePoolV1ConfigNetworkConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkerPoolPrivatePoolV1ConfigNetworkConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkerPoolPrivatePoolV1ConfigNetworkConfig
	} else {

		r.PeeredNetwork = res.PeeredNetwork

		r.PeeredNetworkIPRange = res.PeeredNetworkIPRange

		r.EgressOption = res.EgressOption

	}
	return nil
}

// This object is used to assert a desired state where this WorkerPoolPrivatePoolV1ConfigNetworkConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyWorkerPoolPrivatePoolV1ConfigNetworkConfig *WorkerPoolPrivatePoolV1ConfigNetworkConfig = &WorkerPoolPrivatePoolV1ConfigNetworkConfig{empty: true}

func (r *WorkerPoolPrivatePoolV1ConfigNetworkConfig) Empty() bool {
	return r.empty
}

func (r *WorkerPoolPrivatePoolV1ConfigNetworkConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkerPoolPrivatePoolV1ConfigNetworkConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkerPoolWorkerConfig struct {
	empty        bool    `json:"-"`
	MachineType  *string `json:"machineType"`
	DiskSizeGb   *int64  `json:"diskSizeGb"`
	NoExternalIP *bool   `json:"noExternalIP"`
}

type jsonWorkerPoolWorkerConfig WorkerPoolWorkerConfig

func (r *WorkerPoolWorkerConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkerPoolWorkerConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkerPoolWorkerConfig
	} else {

		r.MachineType = res.MachineType

		r.DiskSizeGb = res.DiskSizeGb

		r.NoExternalIP = res.NoExternalIP

	}
	return nil
}

// This object is used to assert a desired state where this WorkerPoolWorkerConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyWorkerPoolWorkerConfig *WorkerPoolWorkerConfig = &WorkerPoolWorkerConfig{empty: true}

func (r *WorkerPoolWorkerConfig) Empty() bool {
	return r.empty
}

func (r *WorkerPoolWorkerConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkerPoolWorkerConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type WorkerPoolNetworkConfig struct {
	empty                bool    `json:"-"`
	PeeredNetwork        *string `json:"peeredNetwork"`
	PeeredNetworkIPRange *string `json:"peeredNetworkIPRange"`
}

type jsonWorkerPoolNetworkConfig WorkerPoolNetworkConfig

func (r *WorkerPoolNetworkConfig) UnmarshalJSON(data []byte) error {
	var res jsonWorkerPoolNetworkConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyWorkerPoolNetworkConfig
	} else {

		r.PeeredNetwork = res.PeeredNetwork

		r.PeeredNetworkIPRange = res.PeeredNetworkIPRange

	}
	return nil
}

// This object is used to assert a desired state where this WorkerPoolNetworkConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyWorkerPoolNetworkConfig *WorkerPoolNetworkConfig = &WorkerPoolNetworkConfig{empty: true}

func (r *WorkerPoolNetworkConfig) Empty() bool {
	return r.empty
}

func (r *WorkerPoolNetworkConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *WorkerPoolNetworkConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *WorkerPool) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "cloud_build",
		Type:    "WorkerPool",
		Version: "cloudbuild",
	}
}

func (r *WorkerPool) ID() (string, error) {
	if err := extractWorkerPoolFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                   dcl.ValueOrEmptyString(nr.Name),
		"display_name":           dcl.ValueOrEmptyString(nr.DisplayName),
		"uid":                    dcl.ValueOrEmptyString(nr.Uid),
		"annotations":            dcl.ValueOrEmptyString(nr.Annotations),
		"create_time":            dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":            dcl.ValueOrEmptyString(nr.UpdateTime),
		"delete_time":            dcl.ValueOrEmptyString(nr.DeleteTime),
		"state":                  dcl.ValueOrEmptyString(nr.State),
		"private_pool_v1_config": dcl.ValueOrEmptyString(nr.PrivatePoolV1Config),
		"etag":                   dcl.ValueOrEmptyString(nr.Etag),
		"worker_config":          dcl.ValueOrEmptyString(nr.WorkerConfig),
		"network_config":         dcl.ValueOrEmptyString(nr.NetworkConfig),
		"project":                dcl.ValueOrEmptyString(nr.Project),
		"location":               dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/workerPools/{{name}}", params), nil
}

const WorkerPoolMaxPage = -1

type WorkerPoolList struct {
	Items []*WorkerPool

	nextToken string

	pageSize int32

	resource *WorkerPool
}

func (l *WorkerPoolList) HasNext() bool {
	return l.nextToken != ""
}

func (l *WorkerPoolList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listWorkerPool(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListWorkerPool(ctx context.Context, project, location string) (*WorkerPoolList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListWorkerPoolWithMaxResults(ctx, project, location, WorkerPoolMaxPage)

}

func (c *Client) ListWorkerPoolWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*WorkerPoolList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &WorkerPool{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listWorkerPool(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &WorkerPoolList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetWorkerPool(ctx context.Context, r *WorkerPool) (*WorkerPool, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractWorkerPoolFields(r)

	b, err := c.getWorkerPoolRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalWorkerPool(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeWorkerPoolNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractWorkerPoolFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteWorkerPool(ctx context.Context, r *WorkerPool) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("WorkerPool resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting WorkerPool...")
	deleteOp := deleteWorkerPoolOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllWorkerPool deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllWorkerPool(ctx context.Context, project, location string, filter func(*WorkerPool) bool) error {
	listObj, err := c.ListWorkerPool(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllWorkerPool(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllWorkerPool(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyWorkerPool(ctx context.Context, rawDesired *WorkerPool, opts ...dcl.ApplyOption) (*WorkerPool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *WorkerPool
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyWorkerPoolHelper(c, ctx, rawDesired, opts...)
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

func applyWorkerPoolHelper(c *Client, ctx context.Context, rawDesired *WorkerPool, opts ...dcl.ApplyOption) (*WorkerPool, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyWorkerPool...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractWorkerPoolFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.workerPoolDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToWorkerPoolDiffs(c.Config, fieldDiffs, opts)
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
	var ops []workerPoolApiOperation
	if create {
		ops = append(ops, &createWorkerPoolOperation{})
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
	return applyWorkerPoolDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyWorkerPoolDiff(c *Client, ctx context.Context, desired *WorkerPool, rawDesired *WorkerPool, ops []workerPoolApiOperation, opts ...dcl.ApplyOption) (*WorkerPool, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetWorkerPool(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createWorkerPoolOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapWorkerPool(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeWorkerPoolNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeWorkerPoolNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeWorkerPoolDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractWorkerPoolFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractWorkerPoolFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffWorkerPool(c, newDesired, newState)
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
