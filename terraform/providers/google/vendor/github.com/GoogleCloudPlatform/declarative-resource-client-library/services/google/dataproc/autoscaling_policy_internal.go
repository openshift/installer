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
package dataproc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func (r *AutoscalingPolicy) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "basicAlgorithm"); err != nil {
		return err
	}
	if err := dcl.Required(r, "workerConfig"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.BasicAlgorithm) {
		if err := r.BasicAlgorithm.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.WorkerConfig) {
		if err := r.WorkerConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SecondaryWorkerConfig) {
		if err := r.SecondaryWorkerConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *AutoscalingPolicyBasicAlgorithm) validate() error {
	if err := dcl.Required(r, "yarnConfig"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.YarnConfig) {
		if err := r.YarnConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *AutoscalingPolicyBasicAlgorithmYarnConfig) validate() error {
	if err := dcl.Required(r, "gracefulDecommissionTimeout"); err != nil {
		return err
	}
	if err := dcl.Required(r, "scaleUpFactor"); err != nil {
		return err
	}
	if err := dcl.Required(r, "scaleDownFactor"); err != nil {
		return err
	}
	return nil
}
func (r *AutoscalingPolicyWorkerConfig) validate() error {
	if err := dcl.Required(r, "maxInstances"); err != nil {
		return err
	}
	return nil
}
func (r *AutoscalingPolicySecondaryWorkerConfig) validate() error {
	return nil
}
func (r *AutoscalingPolicy) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://dataproc.googleapis.com/v1/", params)
}

func (r *AutoscalingPolicy) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/autoscalingPolicies/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *AutoscalingPolicy) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/autoscalingPolicies", nr.basePath(), userBasePath, params), nil

}

func (r *AutoscalingPolicy) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/autoscalingPolicies", nr.basePath(), userBasePath, params), nil

}

func (r *AutoscalingPolicy) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/autoscalingPolicies/{{name}}", nr.basePath(), userBasePath, params), nil
}

// autoscalingPolicyApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type autoscalingPolicyApiOperation interface {
	do(context.Context, *AutoscalingPolicy, *Client) error
}

// newUpdateAutoscalingPolicyUpdateAutoscalingPolicyRequest creates a request for an
// AutoscalingPolicy resource's UpdateAutoscalingPolicy update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateAutoscalingPolicyUpdateAutoscalingPolicyRequest(ctx context.Context, f *AutoscalingPolicy, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := expandAutoscalingPolicyBasicAlgorithm(c, f.BasicAlgorithm, res); err != nil {
		return nil, fmt.Errorf("error expanding BasicAlgorithm into basicAlgorithm: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["basicAlgorithm"] = v
	}
	if v, err := expandAutoscalingPolicyWorkerConfig(c, f.WorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding WorkerConfig into workerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["workerConfig"] = v
	}
	if v, err := expandAutoscalingPolicySecondaryWorkerConfig(c, f.SecondaryWorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SecondaryWorkerConfig into secondaryWorkerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["secondaryWorkerConfig"] = v
	}
	if v, err := dcl.DeriveField("%s", f.Name); err != nil {
		return nil, err
	} else {
		req["id"] = v
	}

	return req, nil
}

// marshalUpdateAutoscalingPolicyUpdateAutoscalingPolicyRequest converts the update into
// the final JSON request body.
func marshalUpdateAutoscalingPolicyUpdateAutoscalingPolicyRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateAutoscalingPolicyUpdateAutoscalingPolicyOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateAutoscalingPolicyUpdateAutoscalingPolicyOperation) do(ctx context.Context, r *AutoscalingPolicy, c *Client) error {
	_, err := c.GetAutoscalingPolicy(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateAutoscalingPolicy")
	if err != nil {
		return err
	}

	req, err := newUpdateAutoscalingPolicyUpdateAutoscalingPolicyRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateAutoscalingPolicyUpdateAutoscalingPolicyRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PUT", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listAutoscalingPolicyRaw(ctx context.Context, r *AutoscalingPolicy, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != AutoscalingPolicyMaxPage {
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

type listAutoscalingPolicyOperation struct {
	Policies []map[string]interface{} `json:"policies"`
	Token    string                   `json:"nextPageToken"`
}

func (c *Client) listAutoscalingPolicy(ctx context.Context, r *AutoscalingPolicy, pageToken string, pageSize int32) ([]*AutoscalingPolicy, string, error) {
	b, err := c.listAutoscalingPolicyRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listAutoscalingPolicyOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*AutoscalingPolicy
	for _, v := range m.Policies {
		res, err := unmarshalMapAutoscalingPolicy(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllAutoscalingPolicy(ctx context.Context, f func(*AutoscalingPolicy) bool, resources []*AutoscalingPolicy) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteAutoscalingPolicy(ctx, res)
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

type deleteAutoscalingPolicyOperation struct{}

func (op *deleteAutoscalingPolicyOperation) do(ctx context.Context, r *AutoscalingPolicy, c *Client) error {
	r, err := c.GetAutoscalingPolicy(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "AutoscalingPolicy not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetAutoscalingPolicy checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	_, err = dcl.SendRequest(ctx, c.Config, "DELETE", u, body, c.Config.RetryProvider)
	if err != nil {
		return fmt.Errorf("failed to delete AutoscalingPolicy: %w", err)
	}

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetAutoscalingPolicy(ctx, r)
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
type createAutoscalingPolicyOperation struct {
	response map[string]interface{}
}

func (op *createAutoscalingPolicyOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createAutoscalingPolicyOperation) do(ctx context.Context, r *AutoscalingPolicy, c *Client) error {
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

	o, err := dcl.ResponseBodyAsJSON(resp)
	if err != nil {
		return fmt.Errorf("error decoding response body into JSON: %w", err)
	}
	op.response = o

	if _, err := c.GetAutoscalingPolicy(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getAutoscalingPolicyRaw(ctx context.Context, r *AutoscalingPolicy) ([]byte, error) {

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

func (c *Client) autoscalingPolicyDiffsForRawDesired(ctx context.Context, rawDesired *AutoscalingPolicy, opts ...dcl.ApplyOption) (initial, desired *AutoscalingPolicy, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *AutoscalingPolicy
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*AutoscalingPolicy); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected AutoscalingPolicy, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetAutoscalingPolicy(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a AutoscalingPolicy resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve AutoscalingPolicy resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that AutoscalingPolicy resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeAutoscalingPolicyDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for AutoscalingPolicy: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for AutoscalingPolicy: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractAutoscalingPolicyFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeAutoscalingPolicyInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for AutoscalingPolicy: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeAutoscalingPolicyDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for AutoscalingPolicy: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffAutoscalingPolicy(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeAutoscalingPolicyInitialState(rawInitial, rawDesired *AutoscalingPolicy) (*AutoscalingPolicy, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeAutoscalingPolicyDesiredState(rawDesired, rawInitial *AutoscalingPolicy, opts ...dcl.ApplyOption) (*AutoscalingPolicy, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.BasicAlgorithm = canonicalizeAutoscalingPolicyBasicAlgorithm(rawDesired.BasicAlgorithm, nil, opts...)
		rawDesired.WorkerConfig = canonicalizeAutoscalingPolicyWorkerConfig(rawDesired.WorkerConfig, nil, opts...)
		rawDesired.SecondaryWorkerConfig = canonicalizeAutoscalingPolicySecondaryWorkerConfig(rawDesired.SecondaryWorkerConfig, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &AutoscalingPolicy{}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	canonicalDesired.BasicAlgorithm = canonicalizeAutoscalingPolicyBasicAlgorithm(rawDesired.BasicAlgorithm, rawInitial.BasicAlgorithm, opts...)
	canonicalDesired.WorkerConfig = canonicalizeAutoscalingPolicyWorkerConfig(rawDesired.WorkerConfig, rawInitial.WorkerConfig, opts...)
	canonicalDesired.SecondaryWorkerConfig = canonicalizeAutoscalingPolicySecondaryWorkerConfig(rawDesired.SecondaryWorkerConfig, rawInitial.SecondaryWorkerConfig, opts...)
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

func canonicalizeAutoscalingPolicyNewState(c *Client, rawNew, rawDesired *AutoscalingPolicy) (*AutoscalingPolicy, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.BasicAlgorithm) && dcl.IsEmptyValueIndirect(rawDesired.BasicAlgorithm) {
		rawNew.BasicAlgorithm = rawDesired.BasicAlgorithm
	} else {
		rawNew.BasicAlgorithm = canonicalizeNewAutoscalingPolicyBasicAlgorithm(c, rawDesired.BasicAlgorithm, rawNew.BasicAlgorithm)
	}

	if dcl.IsEmptyValueIndirect(rawNew.WorkerConfig) && dcl.IsEmptyValueIndirect(rawDesired.WorkerConfig) {
		rawNew.WorkerConfig = rawDesired.WorkerConfig
	} else {
		rawNew.WorkerConfig = canonicalizeNewAutoscalingPolicyWorkerConfig(c, rawDesired.WorkerConfig, rawNew.WorkerConfig)
	}

	if dcl.IsEmptyValueIndirect(rawNew.SecondaryWorkerConfig) && dcl.IsEmptyValueIndirect(rawDesired.SecondaryWorkerConfig) {
		rawNew.SecondaryWorkerConfig = rawDesired.SecondaryWorkerConfig
	} else {
		rawNew.SecondaryWorkerConfig = canonicalizeNewAutoscalingPolicySecondaryWorkerConfig(c, rawDesired.SecondaryWorkerConfig, rawNew.SecondaryWorkerConfig)
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeAutoscalingPolicyBasicAlgorithm(des, initial *AutoscalingPolicyBasicAlgorithm, opts ...dcl.ApplyOption) *AutoscalingPolicyBasicAlgorithm {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &AutoscalingPolicyBasicAlgorithm{}

	cDes.YarnConfig = canonicalizeAutoscalingPolicyBasicAlgorithmYarnConfig(des.YarnConfig, initial.YarnConfig, opts...)
	if dcl.StringCanonicalize(des.CooldownPeriod, initial.CooldownPeriod) || dcl.IsZeroValue(des.CooldownPeriod) {
		cDes.CooldownPeriod = initial.CooldownPeriod
	} else {
		cDes.CooldownPeriod = des.CooldownPeriod
	}

	return cDes
}

func canonicalizeAutoscalingPolicyBasicAlgorithmSlice(des, initial []AutoscalingPolicyBasicAlgorithm, opts ...dcl.ApplyOption) []AutoscalingPolicyBasicAlgorithm {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]AutoscalingPolicyBasicAlgorithm, 0, len(des))
		for _, d := range des {
			cd := canonicalizeAutoscalingPolicyBasicAlgorithm(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]AutoscalingPolicyBasicAlgorithm, 0, len(des))
	for i, d := range des {
		cd := canonicalizeAutoscalingPolicyBasicAlgorithm(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewAutoscalingPolicyBasicAlgorithm(c *Client, des, nw *AutoscalingPolicyBasicAlgorithm) *AutoscalingPolicyBasicAlgorithm {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for AutoscalingPolicyBasicAlgorithm while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.YarnConfig = canonicalizeNewAutoscalingPolicyBasicAlgorithmYarnConfig(c, des.YarnConfig, nw.YarnConfig)
	if dcl.StringCanonicalize(des.CooldownPeriod, nw.CooldownPeriod) {
		nw.CooldownPeriod = des.CooldownPeriod
	}

	return nw
}

func canonicalizeNewAutoscalingPolicyBasicAlgorithmSet(c *Client, des, nw []AutoscalingPolicyBasicAlgorithm) []AutoscalingPolicyBasicAlgorithm {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []AutoscalingPolicyBasicAlgorithm
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareAutoscalingPolicyBasicAlgorithmNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewAutoscalingPolicyBasicAlgorithm(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewAutoscalingPolicyBasicAlgorithmSlice(c *Client, des, nw []AutoscalingPolicyBasicAlgorithm) []AutoscalingPolicyBasicAlgorithm {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []AutoscalingPolicyBasicAlgorithm
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewAutoscalingPolicyBasicAlgorithm(c, &d, &n))
	}

	return items
}

func canonicalizeAutoscalingPolicyBasicAlgorithmYarnConfig(des, initial *AutoscalingPolicyBasicAlgorithmYarnConfig, opts ...dcl.ApplyOption) *AutoscalingPolicyBasicAlgorithmYarnConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &AutoscalingPolicyBasicAlgorithmYarnConfig{}

	if dcl.StringCanonicalize(des.GracefulDecommissionTimeout, initial.GracefulDecommissionTimeout) || dcl.IsZeroValue(des.GracefulDecommissionTimeout) {
		cDes.GracefulDecommissionTimeout = initial.GracefulDecommissionTimeout
	} else {
		cDes.GracefulDecommissionTimeout = des.GracefulDecommissionTimeout
	}
	if dcl.IsZeroValue(des.ScaleUpFactor) || (dcl.IsEmptyValueIndirect(des.ScaleUpFactor) && dcl.IsEmptyValueIndirect(initial.ScaleUpFactor)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ScaleUpFactor = initial.ScaleUpFactor
	} else {
		cDes.ScaleUpFactor = des.ScaleUpFactor
	}
	if dcl.IsZeroValue(des.ScaleDownFactor) || (dcl.IsEmptyValueIndirect(des.ScaleDownFactor) && dcl.IsEmptyValueIndirect(initial.ScaleDownFactor)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ScaleDownFactor = initial.ScaleDownFactor
	} else {
		cDes.ScaleDownFactor = des.ScaleDownFactor
	}
	if dcl.IsZeroValue(des.ScaleUpMinWorkerFraction) || (dcl.IsEmptyValueIndirect(des.ScaleUpMinWorkerFraction) && dcl.IsEmptyValueIndirect(initial.ScaleUpMinWorkerFraction)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ScaleUpMinWorkerFraction = initial.ScaleUpMinWorkerFraction
	} else {
		cDes.ScaleUpMinWorkerFraction = des.ScaleUpMinWorkerFraction
	}
	if dcl.IsZeroValue(des.ScaleDownMinWorkerFraction) || (dcl.IsEmptyValueIndirect(des.ScaleDownMinWorkerFraction) && dcl.IsEmptyValueIndirect(initial.ScaleDownMinWorkerFraction)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ScaleDownMinWorkerFraction = initial.ScaleDownMinWorkerFraction
	} else {
		cDes.ScaleDownMinWorkerFraction = des.ScaleDownMinWorkerFraction
	}

	return cDes
}

func canonicalizeAutoscalingPolicyBasicAlgorithmYarnConfigSlice(des, initial []AutoscalingPolicyBasicAlgorithmYarnConfig, opts ...dcl.ApplyOption) []AutoscalingPolicyBasicAlgorithmYarnConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]AutoscalingPolicyBasicAlgorithmYarnConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeAutoscalingPolicyBasicAlgorithmYarnConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]AutoscalingPolicyBasicAlgorithmYarnConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeAutoscalingPolicyBasicAlgorithmYarnConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewAutoscalingPolicyBasicAlgorithmYarnConfig(c *Client, des, nw *AutoscalingPolicyBasicAlgorithmYarnConfig) *AutoscalingPolicyBasicAlgorithmYarnConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for AutoscalingPolicyBasicAlgorithmYarnConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.GracefulDecommissionTimeout, nw.GracefulDecommissionTimeout) {
		nw.GracefulDecommissionTimeout = des.GracefulDecommissionTimeout
	}

	return nw
}

func canonicalizeNewAutoscalingPolicyBasicAlgorithmYarnConfigSet(c *Client, des, nw []AutoscalingPolicyBasicAlgorithmYarnConfig) []AutoscalingPolicyBasicAlgorithmYarnConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []AutoscalingPolicyBasicAlgorithmYarnConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareAutoscalingPolicyBasicAlgorithmYarnConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewAutoscalingPolicyBasicAlgorithmYarnConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewAutoscalingPolicyBasicAlgorithmYarnConfigSlice(c *Client, des, nw []AutoscalingPolicyBasicAlgorithmYarnConfig) []AutoscalingPolicyBasicAlgorithmYarnConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []AutoscalingPolicyBasicAlgorithmYarnConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewAutoscalingPolicyBasicAlgorithmYarnConfig(c, &d, &n))
	}

	return items
}

func canonicalizeAutoscalingPolicyWorkerConfig(des, initial *AutoscalingPolicyWorkerConfig, opts ...dcl.ApplyOption) *AutoscalingPolicyWorkerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &AutoscalingPolicyWorkerConfig{}

	if dcl.IsZeroValue(des.MinInstances) || (dcl.IsEmptyValueIndirect(des.MinInstances) && dcl.IsEmptyValueIndirect(initial.MinInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MinInstances = initial.MinInstances
	} else {
		cDes.MinInstances = des.MinInstances
	}
	if dcl.IsZeroValue(des.MaxInstances) || (dcl.IsEmptyValueIndirect(des.MaxInstances) && dcl.IsEmptyValueIndirect(initial.MaxInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MaxInstances = initial.MaxInstances
	} else {
		cDes.MaxInstances = des.MaxInstances
	}
	if dcl.IsZeroValue(des.Weight) || (dcl.IsEmptyValueIndirect(des.Weight) && dcl.IsEmptyValueIndirect(initial.Weight)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Weight = initial.Weight
	} else {
		cDes.Weight = des.Weight
	}

	return cDes
}

func canonicalizeAutoscalingPolicyWorkerConfigSlice(des, initial []AutoscalingPolicyWorkerConfig, opts ...dcl.ApplyOption) []AutoscalingPolicyWorkerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]AutoscalingPolicyWorkerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeAutoscalingPolicyWorkerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]AutoscalingPolicyWorkerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeAutoscalingPolicyWorkerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewAutoscalingPolicyWorkerConfig(c *Client, des, nw *AutoscalingPolicyWorkerConfig) *AutoscalingPolicyWorkerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for AutoscalingPolicyWorkerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewAutoscalingPolicyWorkerConfigSet(c *Client, des, nw []AutoscalingPolicyWorkerConfig) []AutoscalingPolicyWorkerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []AutoscalingPolicyWorkerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareAutoscalingPolicyWorkerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewAutoscalingPolicyWorkerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewAutoscalingPolicyWorkerConfigSlice(c *Client, des, nw []AutoscalingPolicyWorkerConfig) []AutoscalingPolicyWorkerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []AutoscalingPolicyWorkerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewAutoscalingPolicyWorkerConfig(c, &d, &n))
	}

	return items
}

func canonicalizeAutoscalingPolicySecondaryWorkerConfig(des, initial *AutoscalingPolicySecondaryWorkerConfig, opts ...dcl.ApplyOption) *AutoscalingPolicySecondaryWorkerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &AutoscalingPolicySecondaryWorkerConfig{}

	if dcl.IsZeroValue(des.MinInstances) || (dcl.IsEmptyValueIndirect(des.MinInstances) && dcl.IsEmptyValueIndirect(initial.MinInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MinInstances = initial.MinInstances
	} else {
		cDes.MinInstances = des.MinInstances
	}
	if dcl.IsZeroValue(des.MaxInstances) || (dcl.IsEmptyValueIndirect(des.MaxInstances) && dcl.IsEmptyValueIndirect(initial.MaxInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MaxInstances = initial.MaxInstances
	} else {
		cDes.MaxInstances = des.MaxInstances
	}
	if dcl.IsZeroValue(des.Weight) || (dcl.IsEmptyValueIndirect(des.Weight) && dcl.IsEmptyValueIndirect(initial.Weight)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Weight = initial.Weight
	} else {
		cDes.Weight = des.Weight
	}

	return cDes
}

func canonicalizeAutoscalingPolicySecondaryWorkerConfigSlice(des, initial []AutoscalingPolicySecondaryWorkerConfig, opts ...dcl.ApplyOption) []AutoscalingPolicySecondaryWorkerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]AutoscalingPolicySecondaryWorkerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeAutoscalingPolicySecondaryWorkerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]AutoscalingPolicySecondaryWorkerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeAutoscalingPolicySecondaryWorkerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewAutoscalingPolicySecondaryWorkerConfig(c *Client, des, nw *AutoscalingPolicySecondaryWorkerConfig) *AutoscalingPolicySecondaryWorkerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for AutoscalingPolicySecondaryWorkerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewAutoscalingPolicySecondaryWorkerConfigSet(c *Client, des, nw []AutoscalingPolicySecondaryWorkerConfig) []AutoscalingPolicySecondaryWorkerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []AutoscalingPolicySecondaryWorkerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareAutoscalingPolicySecondaryWorkerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewAutoscalingPolicySecondaryWorkerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewAutoscalingPolicySecondaryWorkerConfigSlice(c *Client, des, nw []AutoscalingPolicySecondaryWorkerConfig) []AutoscalingPolicySecondaryWorkerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []AutoscalingPolicySecondaryWorkerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewAutoscalingPolicySecondaryWorkerConfig(c, &d, &n))
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
func diffAutoscalingPolicy(c *Client, desired, actual *AutoscalingPolicy, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BasicAlgorithm, actual.BasicAlgorithm, dcl.DiffInfo{ObjectFunction: compareAutoscalingPolicyBasicAlgorithmNewStyle, EmptyObject: EmptyAutoscalingPolicyBasicAlgorithm, OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("BasicAlgorithm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WorkerConfig, actual.WorkerConfig, dcl.DiffInfo{ObjectFunction: compareAutoscalingPolicyWorkerConfigNewStyle, EmptyObject: EmptyAutoscalingPolicyWorkerConfig, OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("WorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecondaryWorkerConfig, actual.SecondaryWorkerConfig, dcl.DiffInfo{ObjectFunction: compareAutoscalingPolicySecondaryWorkerConfigNewStyle, EmptyObject: EmptyAutoscalingPolicySecondaryWorkerConfig, OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("SecondaryWorkerConfig")); len(ds) != 0 || err != nil {
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
func compareAutoscalingPolicyBasicAlgorithmNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*AutoscalingPolicyBasicAlgorithm)
	if !ok {
		desiredNotPointer, ok := d.(AutoscalingPolicyBasicAlgorithm)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a AutoscalingPolicyBasicAlgorithm or *AutoscalingPolicyBasicAlgorithm", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*AutoscalingPolicyBasicAlgorithm)
	if !ok {
		actualNotPointer, ok := a.(AutoscalingPolicyBasicAlgorithm)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a AutoscalingPolicyBasicAlgorithm", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.YarnConfig, actual.YarnConfig, dcl.DiffInfo{ObjectFunction: compareAutoscalingPolicyBasicAlgorithmYarnConfigNewStyle, EmptyObject: EmptyAutoscalingPolicyBasicAlgorithmYarnConfig, OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("YarnConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CooldownPeriod, actual.CooldownPeriod, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("CooldownPeriod")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareAutoscalingPolicyBasicAlgorithmYarnConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*AutoscalingPolicyBasicAlgorithmYarnConfig)
	if !ok {
		desiredNotPointer, ok := d.(AutoscalingPolicyBasicAlgorithmYarnConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a AutoscalingPolicyBasicAlgorithmYarnConfig or *AutoscalingPolicyBasicAlgorithmYarnConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*AutoscalingPolicyBasicAlgorithmYarnConfig)
	if !ok {
		actualNotPointer, ok := a.(AutoscalingPolicyBasicAlgorithmYarnConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a AutoscalingPolicyBasicAlgorithmYarnConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GracefulDecommissionTimeout, actual.GracefulDecommissionTimeout, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("GracefulDecommissionTimeout")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScaleUpFactor, actual.ScaleUpFactor, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("ScaleUpFactor")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScaleDownFactor, actual.ScaleDownFactor, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("ScaleDownFactor")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScaleUpMinWorkerFraction, actual.ScaleUpMinWorkerFraction, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("ScaleUpMinWorkerFraction")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScaleDownMinWorkerFraction, actual.ScaleDownMinWorkerFraction, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("ScaleDownMinWorkerFraction")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareAutoscalingPolicyWorkerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*AutoscalingPolicyWorkerConfig)
	if !ok {
		desiredNotPointer, ok := d.(AutoscalingPolicyWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a AutoscalingPolicyWorkerConfig or *AutoscalingPolicyWorkerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*AutoscalingPolicyWorkerConfig)
	if !ok {
		actualNotPointer, ok := a.(AutoscalingPolicyWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a AutoscalingPolicyWorkerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MinInstances, actual.MinInstances, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("MinInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxInstances, actual.MaxInstances, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("MaxInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Weight, actual.Weight, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("Weight")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareAutoscalingPolicySecondaryWorkerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*AutoscalingPolicySecondaryWorkerConfig)
	if !ok {
		desiredNotPointer, ok := d.(AutoscalingPolicySecondaryWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a AutoscalingPolicySecondaryWorkerConfig or *AutoscalingPolicySecondaryWorkerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*AutoscalingPolicySecondaryWorkerConfig)
	if !ok {
		actualNotPointer, ok := a.(AutoscalingPolicySecondaryWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a AutoscalingPolicySecondaryWorkerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MinInstances, actual.MinInstances, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("MinInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxInstances, actual.MaxInstances, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("MaxInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Weight, actual.Weight, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateAutoscalingPolicyUpdateAutoscalingPolicyOperation")}, fn.AddNest("Weight")); len(ds) != 0 || err != nil {
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
func (r *AutoscalingPolicy) urlNormalized() *AutoscalingPolicy {
	normalized := dcl.Copy(*r).(AutoscalingPolicy)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *AutoscalingPolicy) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateAutoscalingPolicy" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/autoscalingPolicies/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the AutoscalingPolicy resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *AutoscalingPolicy) marshal(c *Client) ([]byte, error) {
	m, err := expandAutoscalingPolicy(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling AutoscalingPolicy: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalAutoscalingPolicy decodes JSON responses into the AutoscalingPolicy resource schema.
func unmarshalAutoscalingPolicy(b []byte, c *Client, res *AutoscalingPolicy) (*AutoscalingPolicy, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapAutoscalingPolicy(m, c, res)
}

func unmarshalMapAutoscalingPolicy(m map[string]interface{}, c *Client, res *AutoscalingPolicy) (*AutoscalingPolicy, error) {

	flattened := flattenAutoscalingPolicy(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandAutoscalingPolicy expands AutoscalingPolicy into a JSON request object.
func expandAutoscalingPolicy(c *Client, f *AutoscalingPolicy) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["id"] = v
	}
	if v, err := expandAutoscalingPolicyBasicAlgorithm(c, f.BasicAlgorithm, res); err != nil {
		return nil, fmt.Errorf("error expanding BasicAlgorithm into basicAlgorithm: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["basicAlgorithm"] = v
	}
	if v, err := expandAutoscalingPolicyWorkerConfig(c, f.WorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding WorkerConfig into workerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["workerConfig"] = v
	}
	if v, err := expandAutoscalingPolicySecondaryWorkerConfig(c, f.SecondaryWorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SecondaryWorkerConfig into secondaryWorkerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["secondaryWorkerConfig"] = v
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

// flattenAutoscalingPolicy flattens AutoscalingPolicy from a JSON request object into the
// AutoscalingPolicy type.
func flattenAutoscalingPolicy(c *Client, i interface{}, res *AutoscalingPolicy) *AutoscalingPolicy {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &AutoscalingPolicy{}
	resultRes.Name = dcl.FlattenString(m["id"])
	resultRes.BasicAlgorithm = flattenAutoscalingPolicyBasicAlgorithm(c, m["basicAlgorithm"], res)
	resultRes.WorkerConfig = flattenAutoscalingPolicyWorkerConfig(c, m["workerConfig"], res)
	resultRes.SecondaryWorkerConfig = flattenAutoscalingPolicySecondaryWorkerConfig(c, m["secondaryWorkerConfig"], res)
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandAutoscalingPolicyBasicAlgorithmMap expands the contents of AutoscalingPolicyBasicAlgorithm into a JSON
// request object.
func expandAutoscalingPolicyBasicAlgorithmMap(c *Client, f map[string]AutoscalingPolicyBasicAlgorithm, res *AutoscalingPolicy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandAutoscalingPolicyBasicAlgorithm(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandAutoscalingPolicyBasicAlgorithmSlice expands the contents of AutoscalingPolicyBasicAlgorithm into a JSON
// request object.
func expandAutoscalingPolicyBasicAlgorithmSlice(c *Client, f []AutoscalingPolicyBasicAlgorithm, res *AutoscalingPolicy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandAutoscalingPolicyBasicAlgorithm(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenAutoscalingPolicyBasicAlgorithmMap flattens the contents of AutoscalingPolicyBasicAlgorithm from a JSON
// response object.
func flattenAutoscalingPolicyBasicAlgorithmMap(c *Client, i interface{}, res *AutoscalingPolicy) map[string]AutoscalingPolicyBasicAlgorithm {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]AutoscalingPolicyBasicAlgorithm{}
	}

	if len(a) == 0 {
		return map[string]AutoscalingPolicyBasicAlgorithm{}
	}

	items := make(map[string]AutoscalingPolicyBasicAlgorithm)
	for k, item := range a {
		items[k] = *flattenAutoscalingPolicyBasicAlgorithm(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenAutoscalingPolicyBasicAlgorithmSlice flattens the contents of AutoscalingPolicyBasicAlgorithm from a JSON
// response object.
func flattenAutoscalingPolicyBasicAlgorithmSlice(c *Client, i interface{}, res *AutoscalingPolicy) []AutoscalingPolicyBasicAlgorithm {
	a, ok := i.([]interface{})
	if !ok {
		return []AutoscalingPolicyBasicAlgorithm{}
	}

	if len(a) == 0 {
		return []AutoscalingPolicyBasicAlgorithm{}
	}

	items := make([]AutoscalingPolicyBasicAlgorithm, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenAutoscalingPolicyBasicAlgorithm(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandAutoscalingPolicyBasicAlgorithm expands an instance of AutoscalingPolicyBasicAlgorithm into a JSON
// request object.
func expandAutoscalingPolicyBasicAlgorithm(c *Client, f *AutoscalingPolicyBasicAlgorithm, res *AutoscalingPolicy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandAutoscalingPolicyBasicAlgorithmYarnConfig(c, f.YarnConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding YarnConfig into yarnConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["yarnConfig"] = v
	}
	if v := f.CooldownPeriod; !dcl.IsEmptyValueIndirect(v) {
		m["cooldownPeriod"] = v
	}

	return m, nil
}

// flattenAutoscalingPolicyBasicAlgorithm flattens an instance of AutoscalingPolicyBasicAlgorithm from a JSON
// response object.
func flattenAutoscalingPolicyBasicAlgorithm(c *Client, i interface{}, res *AutoscalingPolicy) *AutoscalingPolicyBasicAlgorithm {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &AutoscalingPolicyBasicAlgorithm{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyAutoscalingPolicyBasicAlgorithm
	}
	r.YarnConfig = flattenAutoscalingPolicyBasicAlgorithmYarnConfig(c, m["yarnConfig"], res)
	r.CooldownPeriod = dcl.FlattenString(m["cooldownPeriod"])

	return r
}

// expandAutoscalingPolicyBasicAlgorithmYarnConfigMap expands the contents of AutoscalingPolicyBasicAlgorithmYarnConfig into a JSON
// request object.
func expandAutoscalingPolicyBasicAlgorithmYarnConfigMap(c *Client, f map[string]AutoscalingPolicyBasicAlgorithmYarnConfig, res *AutoscalingPolicy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandAutoscalingPolicyBasicAlgorithmYarnConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandAutoscalingPolicyBasicAlgorithmYarnConfigSlice expands the contents of AutoscalingPolicyBasicAlgorithmYarnConfig into a JSON
// request object.
func expandAutoscalingPolicyBasicAlgorithmYarnConfigSlice(c *Client, f []AutoscalingPolicyBasicAlgorithmYarnConfig, res *AutoscalingPolicy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandAutoscalingPolicyBasicAlgorithmYarnConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenAutoscalingPolicyBasicAlgorithmYarnConfigMap flattens the contents of AutoscalingPolicyBasicAlgorithmYarnConfig from a JSON
// response object.
func flattenAutoscalingPolicyBasicAlgorithmYarnConfigMap(c *Client, i interface{}, res *AutoscalingPolicy) map[string]AutoscalingPolicyBasicAlgorithmYarnConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]AutoscalingPolicyBasicAlgorithmYarnConfig{}
	}

	if len(a) == 0 {
		return map[string]AutoscalingPolicyBasicAlgorithmYarnConfig{}
	}

	items := make(map[string]AutoscalingPolicyBasicAlgorithmYarnConfig)
	for k, item := range a {
		items[k] = *flattenAutoscalingPolicyBasicAlgorithmYarnConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenAutoscalingPolicyBasicAlgorithmYarnConfigSlice flattens the contents of AutoscalingPolicyBasicAlgorithmYarnConfig from a JSON
// response object.
func flattenAutoscalingPolicyBasicAlgorithmYarnConfigSlice(c *Client, i interface{}, res *AutoscalingPolicy) []AutoscalingPolicyBasicAlgorithmYarnConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []AutoscalingPolicyBasicAlgorithmYarnConfig{}
	}

	if len(a) == 0 {
		return []AutoscalingPolicyBasicAlgorithmYarnConfig{}
	}

	items := make([]AutoscalingPolicyBasicAlgorithmYarnConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenAutoscalingPolicyBasicAlgorithmYarnConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandAutoscalingPolicyBasicAlgorithmYarnConfig expands an instance of AutoscalingPolicyBasicAlgorithmYarnConfig into a JSON
// request object.
func expandAutoscalingPolicyBasicAlgorithmYarnConfig(c *Client, f *AutoscalingPolicyBasicAlgorithmYarnConfig, res *AutoscalingPolicy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GracefulDecommissionTimeout; !dcl.IsEmptyValueIndirect(v) {
		m["gracefulDecommissionTimeout"] = v
	}
	if v := f.ScaleUpFactor; !dcl.IsEmptyValueIndirect(v) {
		m["scaleUpFactor"] = v
	}
	if v := f.ScaleDownFactor; !dcl.IsEmptyValueIndirect(v) {
		m["scaleDownFactor"] = v
	}
	if v := f.ScaleUpMinWorkerFraction; !dcl.IsEmptyValueIndirect(v) {
		m["scaleUpMinWorkerFraction"] = v
	}
	if v := f.ScaleDownMinWorkerFraction; !dcl.IsEmptyValueIndirect(v) {
		m["scaleDownMinWorkerFraction"] = v
	}

	return m, nil
}

// flattenAutoscalingPolicyBasicAlgorithmYarnConfig flattens an instance of AutoscalingPolicyBasicAlgorithmYarnConfig from a JSON
// response object.
func flattenAutoscalingPolicyBasicAlgorithmYarnConfig(c *Client, i interface{}, res *AutoscalingPolicy) *AutoscalingPolicyBasicAlgorithmYarnConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &AutoscalingPolicyBasicAlgorithmYarnConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyAutoscalingPolicyBasicAlgorithmYarnConfig
	}
	r.GracefulDecommissionTimeout = dcl.FlattenString(m["gracefulDecommissionTimeout"])
	r.ScaleUpFactor = dcl.FlattenDouble(m["scaleUpFactor"])
	r.ScaleDownFactor = dcl.FlattenDouble(m["scaleDownFactor"])
	r.ScaleUpMinWorkerFraction = dcl.FlattenDouble(m["scaleUpMinWorkerFraction"])
	r.ScaleDownMinWorkerFraction = dcl.FlattenDouble(m["scaleDownMinWorkerFraction"])

	return r
}

// expandAutoscalingPolicyWorkerConfigMap expands the contents of AutoscalingPolicyWorkerConfig into a JSON
// request object.
func expandAutoscalingPolicyWorkerConfigMap(c *Client, f map[string]AutoscalingPolicyWorkerConfig, res *AutoscalingPolicy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandAutoscalingPolicyWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandAutoscalingPolicyWorkerConfigSlice expands the contents of AutoscalingPolicyWorkerConfig into a JSON
// request object.
func expandAutoscalingPolicyWorkerConfigSlice(c *Client, f []AutoscalingPolicyWorkerConfig, res *AutoscalingPolicy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandAutoscalingPolicyWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenAutoscalingPolicyWorkerConfigMap flattens the contents of AutoscalingPolicyWorkerConfig from a JSON
// response object.
func flattenAutoscalingPolicyWorkerConfigMap(c *Client, i interface{}, res *AutoscalingPolicy) map[string]AutoscalingPolicyWorkerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]AutoscalingPolicyWorkerConfig{}
	}

	if len(a) == 0 {
		return map[string]AutoscalingPolicyWorkerConfig{}
	}

	items := make(map[string]AutoscalingPolicyWorkerConfig)
	for k, item := range a {
		items[k] = *flattenAutoscalingPolicyWorkerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenAutoscalingPolicyWorkerConfigSlice flattens the contents of AutoscalingPolicyWorkerConfig from a JSON
// response object.
func flattenAutoscalingPolicyWorkerConfigSlice(c *Client, i interface{}, res *AutoscalingPolicy) []AutoscalingPolicyWorkerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []AutoscalingPolicyWorkerConfig{}
	}

	if len(a) == 0 {
		return []AutoscalingPolicyWorkerConfig{}
	}

	items := make([]AutoscalingPolicyWorkerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenAutoscalingPolicyWorkerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandAutoscalingPolicyWorkerConfig expands an instance of AutoscalingPolicyWorkerConfig into a JSON
// request object.
func expandAutoscalingPolicyWorkerConfig(c *Client, f *AutoscalingPolicyWorkerConfig, res *AutoscalingPolicy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MinInstances; !dcl.IsEmptyValueIndirect(v) {
		m["minInstances"] = v
	}
	if v := f.MaxInstances; !dcl.IsEmptyValueIndirect(v) {
		m["maxInstances"] = v
	}
	if v := f.Weight; !dcl.IsEmptyValueIndirect(v) {
		m["weight"] = v
	}

	return m, nil
}

// flattenAutoscalingPolicyWorkerConfig flattens an instance of AutoscalingPolicyWorkerConfig from a JSON
// response object.
func flattenAutoscalingPolicyWorkerConfig(c *Client, i interface{}, res *AutoscalingPolicy) *AutoscalingPolicyWorkerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &AutoscalingPolicyWorkerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyAutoscalingPolicyWorkerConfig
	}
	r.MinInstances = dcl.FlattenInteger(m["minInstances"])
	r.MaxInstances = dcl.FlattenInteger(m["maxInstances"])
	r.Weight = dcl.FlattenInteger(m["weight"])

	return r
}

// expandAutoscalingPolicySecondaryWorkerConfigMap expands the contents of AutoscalingPolicySecondaryWorkerConfig into a JSON
// request object.
func expandAutoscalingPolicySecondaryWorkerConfigMap(c *Client, f map[string]AutoscalingPolicySecondaryWorkerConfig, res *AutoscalingPolicy) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandAutoscalingPolicySecondaryWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandAutoscalingPolicySecondaryWorkerConfigSlice expands the contents of AutoscalingPolicySecondaryWorkerConfig into a JSON
// request object.
func expandAutoscalingPolicySecondaryWorkerConfigSlice(c *Client, f []AutoscalingPolicySecondaryWorkerConfig, res *AutoscalingPolicy) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandAutoscalingPolicySecondaryWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenAutoscalingPolicySecondaryWorkerConfigMap flattens the contents of AutoscalingPolicySecondaryWorkerConfig from a JSON
// response object.
func flattenAutoscalingPolicySecondaryWorkerConfigMap(c *Client, i interface{}, res *AutoscalingPolicy) map[string]AutoscalingPolicySecondaryWorkerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]AutoscalingPolicySecondaryWorkerConfig{}
	}

	if len(a) == 0 {
		return map[string]AutoscalingPolicySecondaryWorkerConfig{}
	}

	items := make(map[string]AutoscalingPolicySecondaryWorkerConfig)
	for k, item := range a {
		items[k] = *flattenAutoscalingPolicySecondaryWorkerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenAutoscalingPolicySecondaryWorkerConfigSlice flattens the contents of AutoscalingPolicySecondaryWorkerConfig from a JSON
// response object.
func flattenAutoscalingPolicySecondaryWorkerConfigSlice(c *Client, i interface{}, res *AutoscalingPolicy) []AutoscalingPolicySecondaryWorkerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []AutoscalingPolicySecondaryWorkerConfig{}
	}

	if len(a) == 0 {
		return []AutoscalingPolicySecondaryWorkerConfig{}
	}

	items := make([]AutoscalingPolicySecondaryWorkerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenAutoscalingPolicySecondaryWorkerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandAutoscalingPolicySecondaryWorkerConfig expands an instance of AutoscalingPolicySecondaryWorkerConfig into a JSON
// request object.
func expandAutoscalingPolicySecondaryWorkerConfig(c *Client, f *AutoscalingPolicySecondaryWorkerConfig, res *AutoscalingPolicy) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MinInstances; !dcl.IsEmptyValueIndirect(v) {
		m["minInstances"] = v
	}
	if v := f.MaxInstances; !dcl.IsEmptyValueIndirect(v) {
		m["maxInstances"] = v
	}
	if v := f.Weight; !dcl.IsEmptyValueIndirect(v) {
		m["weight"] = v
	}

	return m, nil
}

// flattenAutoscalingPolicySecondaryWorkerConfig flattens an instance of AutoscalingPolicySecondaryWorkerConfig from a JSON
// response object.
func flattenAutoscalingPolicySecondaryWorkerConfig(c *Client, i interface{}, res *AutoscalingPolicy) *AutoscalingPolicySecondaryWorkerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &AutoscalingPolicySecondaryWorkerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyAutoscalingPolicySecondaryWorkerConfig
	}
	r.MinInstances = dcl.FlattenInteger(m["minInstances"])
	r.MaxInstances = dcl.FlattenInteger(m["maxInstances"])
	r.Weight = dcl.FlattenInteger(m["weight"])

	return r
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *AutoscalingPolicy) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalAutoscalingPolicy(b, c, r)
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

type autoscalingPolicyDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         autoscalingPolicyApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToAutoscalingPolicyDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]autoscalingPolicyDiff, error) {
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
	var diffs []autoscalingPolicyDiff
	// For each operation name, create a autoscalingPolicyDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := autoscalingPolicyDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToAutoscalingPolicyApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToAutoscalingPolicyApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (autoscalingPolicyApiOperation, error) {
	switch opName {

	case "updateAutoscalingPolicyUpdateAutoscalingPolicyOperation":
		return &updateAutoscalingPolicyUpdateAutoscalingPolicyOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractAutoscalingPolicyFields(r *AutoscalingPolicy) error {
	vBasicAlgorithm := r.BasicAlgorithm
	if vBasicAlgorithm == nil {
		// note: explicitly not the empty object.
		vBasicAlgorithm = &AutoscalingPolicyBasicAlgorithm{}
	}
	if err := extractAutoscalingPolicyBasicAlgorithmFields(r, vBasicAlgorithm); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBasicAlgorithm) {
		r.BasicAlgorithm = vBasicAlgorithm
	}
	vWorkerConfig := r.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &AutoscalingPolicyWorkerConfig{}
	}
	if err := extractAutoscalingPolicyWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		r.WorkerConfig = vWorkerConfig
	}
	vSecondaryWorkerConfig := r.SecondaryWorkerConfig
	if vSecondaryWorkerConfig == nil {
		// note: explicitly not the empty object.
		vSecondaryWorkerConfig = &AutoscalingPolicySecondaryWorkerConfig{}
	}
	if err := extractAutoscalingPolicySecondaryWorkerConfigFields(r, vSecondaryWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecondaryWorkerConfig) {
		r.SecondaryWorkerConfig = vSecondaryWorkerConfig
	}
	return nil
}
func extractAutoscalingPolicyBasicAlgorithmFields(r *AutoscalingPolicy, o *AutoscalingPolicyBasicAlgorithm) error {
	vYarnConfig := o.YarnConfig
	if vYarnConfig == nil {
		// note: explicitly not the empty object.
		vYarnConfig = &AutoscalingPolicyBasicAlgorithmYarnConfig{}
	}
	if err := extractAutoscalingPolicyBasicAlgorithmYarnConfigFields(r, vYarnConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vYarnConfig) {
		o.YarnConfig = vYarnConfig
	}
	return nil
}
func extractAutoscalingPolicyBasicAlgorithmYarnConfigFields(r *AutoscalingPolicy, o *AutoscalingPolicyBasicAlgorithmYarnConfig) error {
	return nil
}
func extractAutoscalingPolicyWorkerConfigFields(r *AutoscalingPolicy, o *AutoscalingPolicyWorkerConfig) error {
	return nil
}
func extractAutoscalingPolicySecondaryWorkerConfigFields(r *AutoscalingPolicy, o *AutoscalingPolicySecondaryWorkerConfig) error {
	return nil
}

func postReadExtractAutoscalingPolicyFields(r *AutoscalingPolicy) error {
	vBasicAlgorithm := r.BasicAlgorithm
	if vBasicAlgorithm == nil {
		// note: explicitly not the empty object.
		vBasicAlgorithm = &AutoscalingPolicyBasicAlgorithm{}
	}
	if err := postReadExtractAutoscalingPolicyBasicAlgorithmFields(r, vBasicAlgorithm); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBasicAlgorithm) {
		r.BasicAlgorithm = vBasicAlgorithm
	}
	vWorkerConfig := r.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &AutoscalingPolicyWorkerConfig{}
	}
	if err := postReadExtractAutoscalingPolicyWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		r.WorkerConfig = vWorkerConfig
	}
	vSecondaryWorkerConfig := r.SecondaryWorkerConfig
	if vSecondaryWorkerConfig == nil {
		// note: explicitly not the empty object.
		vSecondaryWorkerConfig = &AutoscalingPolicySecondaryWorkerConfig{}
	}
	if err := postReadExtractAutoscalingPolicySecondaryWorkerConfigFields(r, vSecondaryWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecondaryWorkerConfig) {
		r.SecondaryWorkerConfig = vSecondaryWorkerConfig
	}
	return nil
}
func postReadExtractAutoscalingPolicyBasicAlgorithmFields(r *AutoscalingPolicy, o *AutoscalingPolicyBasicAlgorithm) error {
	vYarnConfig := o.YarnConfig
	if vYarnConfig == nil {
		// note: explicitly not the empty object.
		vYarnConfig = &AutoscalingPolicyBasicAlgorithmYarnConfig{}
	}
	if err := extractAutoscalingPolicyBasicAlgorithmYarnConfigFields(r, vYarnConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vYarnConfig) {
		o.YarnConfig = vYarnConfig
	}
	return nil
}
func postReadExtractAutoscalingPolicyBasicAlgorithmYarnConfigFields(r *AutoscalingPolicy, o *AutoscalingPolicyBasicAlgorithmYarnConfig) error {
	return nil
}
func postReadExtractAutoscalingPolicyWorkerConfigFields(r *AutoscalingPolicy, o *AutoscalingPolicyWorkerConfig) error {
	return nil
}
func postReadExtractAutoscalingPolicySecondaryWorkerConfigFields(r *AutoscalingPolicy, o *AutoscalingPolicySecondaryWorkerConfig) error {
	return nil
}
