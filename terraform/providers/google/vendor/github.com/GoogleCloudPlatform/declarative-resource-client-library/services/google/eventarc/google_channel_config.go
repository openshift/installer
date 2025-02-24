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
	"context"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type GoogleChannelConfig struct {
	Name          *string `json:"name"`
	UpdateTime    *string `json:"updateTime"`
	CryptoKeyName *string `json:"cryptoKeyName"`
	Project       *string `json:"project"`
	Location      *string `json:"location"`
}

func (r *GoogleChannelConfig) String() string {
	return dcl.SprintResource(r)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *GoogleChannelConfig) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "eventarc",
		Type:    "GoogleChannelConfig",
		Version: "eventarc",
	}
}

func (r *GoogleChannelConfig) ID() (string, error) {
	if err := extractGoogleChannelConfigFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":            dcl.ValueOrEmptyString(nr.Name),
		"update_time":     dcl.ValueOrEmptyString(nr.UpdateTime),
		"crypto_key_name": dcl.ValueOrEmptyString(nr.CryptoKeyName),
		"project":         dcl.ValueOrEmptyString(nr.Project),
		"location":        dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/googleChannelConfig", params), nil
}

const GoogleChannelConfigMaxPage = -1

type GoogleChannelConfigList struct {
	Items []*GoogleChannelConfig

	nextToken string

	resource *GoogleChannelConfig
}

func (c *Client) GetGoogleChannelConfig(ctx context.Context, r *GoogleChannelConfig) (*GoogleChannelConfig, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractGoogleChannelConfigFields(r)

	b, err := c.getGoogleChannelConfigRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalGoogleChannelConfig(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeGoogleChannelConfigNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractGoogleChannelConfigFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteGoogleChannelConfig(ctx context.Context, r *GoogleChannelConfig) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("GoogleChannelConfig resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting GoogleChannelConfig...")
	deleteOp := deleteGoogleChannelConfigOperation{}
	return deleteOp.do(ctx, r, c)
}

func (c *Client) ApplyGoogleChannelConfig(ctx context.Context, rawDesired *GoogleChannelConfig, opts ...dcl.ApplyOption) (*GoogleChannelConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *GoogleChannelConfig
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyGoogleChannelConfigHelper(c, ctx, rawDesired, opts...)
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

func applyGoogleChannelConfigHelper(c *Client, ctx context.Context, rawDesired *GoogleChannelConfig, opts ...dcl.ApplyOption) (*GoogleChannelConfig, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyGoogleChannelConfig...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractGoogleChannelConfigFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.googleChannelConfigDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToGoogleChannelConfigDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		return nil, dcl.ApplyInfeasibleError{Message: "No initial state found for singleton resource."}
	} else {
		for _, d := range diffs {
			if d.UpdateOp == nil {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) no update method found for field", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}
	var ops []googleChannelConfigApiOperation
	for _, d := range diffs {
		ops = append(ops, d.UpdateOp)
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
	return applyGoogleChannelConfigDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyGoogleChannelConfigDiff(c *Client, ctx context.Context, desired *GoogleChannelConfig, rawDesired *GoogleChannelConfig, ops []googleChannelConfigApiOperation, opts ...dcl.ApplyOption) (*GoogleChannelConfig, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetGoogleChannelConfig(ctx, desired)
	if err != nil {
		return nil, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeGoogleChannelConfigNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeGoogleChannelConfigDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractGoogleChannelConfigFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractGoogleChannelConfigFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffGoogleChannelConfig(c, newDesired, newState)
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
