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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *WorkerPool) validate() error {

	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"NetworkConfig", "PrivatePoolV1Config"}, r.NetworkConfig, r.PrivatePoolV1Config); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"WorkerConfig", "PrivatePoolV1Config"}, r.WorkerConfig, r.PrivatePoolV1Config); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.PrivatePoolV1Config) {
		if err := r.PrivatePoolV1Config.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.WorkerConfig) {
		if err := r.WorkerConfig.validate(); err != nil {
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
func (r *WorkerPoolPrivatePoolV1Config) validate() error {
	if !dcl.IsEmptyValueIndirect(r.WorkerConfig) {
		if err := r.WorkerConfig.validate(); err != nil {
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
func (r *WorkerPoolPrivatePoolV1ConfigWorkerConfig) validate() error {
	return nil
}
func (r *WorkerPoolPrivatePoolV1ConfigNetworkConfig) validate() error {
	return nil
}
func (r *WorkerPoolWorkerConfig) validate() error {
	return nil
}
func (r *WorkerPoolNetworkConfig) validate() error {
	if err := dcl.Required(r, "peeredNetwork"); err != nil {
		return err
	}
	return nil
}
func (r *WorkerPool) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://cloudbuild.googleapis.com/v1/", params)
}

func (r *WorkerPool) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/workerPools/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *WorkerPool) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/workerPools", nr.basePath(), userBasePath, params), nil

}

func (r *WorkerPool) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/workerPools?workerPoolId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *WorkerPool) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/workerPools/{{name}}", nr.basePath(), userBasePath, params), nil
}

// workerPoolApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type workerPoolApiOperation interface {
	do(context.Context, *WorkerPool, *Client) error
}

// newUpdateWorkerPoolUpdateWorkerPoolRequest creates a request for an
// WorkerPool resource's UpdateWorkerPool update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateWorkerPoolUpdateWorkerPoolRequest(ctx context.Context, f *WorkerPool, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		req["displayName"] = v
	}
	if v := f.Annotations; !dcl.IsEmptyValueIndirect(v) {
		req["annotations"] = v
	}
	if v, err := expandWorkerPoolPrivatePoolV1Config(c, f.PrivatePoolV1Config, res); err != nil {
		return nil, fmt.Errorf("error expanding PrivatePoolV1Config into privatePoolV1Config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["privatePoolV1Config"] = v
	}
	if v, err := expandWorkerPoolWorkerConfig(c, f.WorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding WorkerConfig into workerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["workerConfig"] = v
	}
	b, err := c.getWorkerPoolRaw(ctx, f)
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

// marshalUpdateWorkerPoolUpdateWorkerPoolRequest converts the update into
// the final JSON request body.
func marshalUpdateWorkerPoolUpdateWorkerPoolRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateWorkerPoolUpdateWorkerPoolOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateWorkerPoolUpdateWorkerPoolOperation) do(ctx context.Context, r *WorkerPool, c *Client) error {
	_, err := c.GetWorkerPool(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateWorkerPool")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateWorkerPoolUpdateWorkerPoolRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateWorkerPoolUpdateWorkerPoolRequest(c, req)
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

func (c *Client) listWorkerPoolRaw(ctx context.Context, r *WorkerPool, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != WorkerPoolMaxPage {
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

type listWorkerPoolOperation struct {
	WorkerPools []map[string]interface{} `json:"workerPools"`
	Token       string                   `json:"nextPageToken"`
}

func (c *Client) listWorkerPool(ctx context.Context, r *WorkerPool, pageToken string, pageSize int32) ([]*WorkerPool, string, error) {
	b, err := c.listWorkerPoolRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listWorkerPoolOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*WorkerPool
	for _, v := range m.WorkerPools {
		res, err := unmarshalMapWorkerPool(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllWorkerPool(ctx context.Context, f func(*WorkerPool) bool, resources []*WorkerPool) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteWorkerPool(ctx, res)
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

type deleteWorkerPoolOperation struct{}

func (op *deleteWorkerPoolOperation) do(ctx context.Context, r *WorkerPool, c *Client) error {
	r, err := c.GetWorkerPool(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "WorkerPool not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetWorkerPool checking for existence. error: %v", err)
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
		_, err := c.GetWorkerPool(ctx, r)
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
type createWorkerPoolOperation struct {
	response map[string]interface{}
}

func (op *createWorkerPoolOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createWorkerPoolOperation) do(ctx context.Context, r *WorkerPool, c *Client) error {
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

	if _, err := c.GetWorkerPool(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getWorkerPoolRaw(ctx context.Context, r *WorkerPool) ([]byte, error) {

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

func (c *Client) workerPoolDiffsForRawDesired(ctx context.Context, rawDesired *WorkerPool, opts ...dcl.ApplyOption) (initial, desired *WorkerPool, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *WorkerPool
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*WorkerPool); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected WorkerPool, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetWorkerPool(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a WorkerPool resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve WorkerPool resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that WorkerPool resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeWorkerPoolDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for WorkerPool: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for WorkerPool: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractWorkerPoolFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeWorkerPoolInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for WorkerPool: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeWorkerPoolDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for WorkerPool: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffWorkerPool(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeWorkerPoolInitialState(rawInitial, rawDesired *WorkerPool) (*WorkerPool, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.

	if !dcl.IsZeroValue(rawInitial.NetworkConfig) {
		// Check if anything else is set.
		if dcl.AnySet() {
			rawInitial.NetworkConfig = EmptyWorkerPoolNetworkConfig
		}
	}

	if !dcl.IsZeroValue(rawInitial.PrivatePoolV1Config) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.NetworkConfig) {
			rawInitial.PrivatePoolV1Config = EmptyWorkerPoolPrivatePoolV1Config
		}
	}

	if !dcl.IsZeroValue(rawInitial.WorkerConfig) {
		// Check if anything else is set.
		if dcl.AnySet() {
			rawInitial.WorkerConfig = EmptyWorkerPoolWorkerConfig
		}
	}

	if !dcl.IsZeroValue(rawInitial.PrivatePoolV1Config) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.WorkerConfig) {
			rawInitial.PrivatePoolV1Config = EmptyWorkerPoolPrivatePoolV1Config
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

func canonicalizeWorkerPoolDesiredState(rawDesired, rawInitial *WorkerPool, opts ...dcl.ApplyOption) (*WorkerPool, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.PrivatePoolV1Config = canonicalizeWorkerPoolPrivatePoolV1Config(rawDesired.PrivatePoolV1Config, nil, opts...)
		rawDesired.WorkerConfig = canonicalizeWorkerPoolWorkerConfig(rawDesired.WorkerConfig, nil, opts...)
		rawDesired.NetworkConfig = canonicalizeWorkerPoolNetworkConfig(rawDesired.NetworkConfig, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &WorkerPool{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.DisplayName, rawInitial.DisplayName) {
		canonicalDesired.DisplayName = rawInitial.DisplayName
	} else {
		canonicalDesired.DisplayName = rawDesired.DisplayName
	}
	if dcl.IsZeroValue(rawDesired.Annotations) || (dcl.IsEmptyValueIndirect(rawDesired.Annotations) && dcl.IsEmptyValueIndirect(rawInitial.Annotations)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Annotations = rawInitial.Annotations
	} else {
		canonicalDesired.Annotations = rawDesired.Annotations
	}
	canonicalDesired.PrivatePoolV1Config = canonicalizeWorkerPoolPrivatePoolV1Config(rawDesired.PrivatePoolV1Config, rawInitial.PrivatePoolV1Config, opts...)
	canonicalDesired.WorkerConfig = canonicalizeWorkerPoolWorkerConfig(rawDesired.WorkerConfig, rawInitial.WorkerConfig, opts...)
	canonicalDesired.NetworkConfig = canonicalizeWorkerPoolNetworkConfig(rawDesired.NetworkConfig, rawInitial.NetworkConfig, opts...)
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

	if canonicalDesired.NetworkConfig != nil {
		// Check if anything else is set.
		if dcl.AnySet() {
			canonicalDesired.NetworkConfig = EmptyWorkerPoolNetworkConfig
		}
	}

	if canonicalDesired.PrivatePoolV1Config != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.NetworkConfig) {
			canonicalDesired.PrivatePoolV1Config = EmptyWorkerPoolPrivatePoolV1Config
		}
	}

	if canonicalDesired.WorkerConfig != nil {
		// Check if anything else is set.
		if dcl.AnySet() {
			canonicalDesired.WorkerConfig = EmptyWorkerPoolWorkerConfig
		}
	}

	if canonicalDesired.PrivatePoolV1Config != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.WorkerConfig) {
			canonicalDesired.PrivatePoolV1Config = EmptyWorkerPoolPrivatePoolV1Config
		}
	}

	return canonicalDesired, nil
}

func canonicalizeWorkerPoolNewState(c *Client, rawNew, rawDesired *WorkerPool) (*WorkerPool, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.DisplayName) && dcl.IsEmptyValueIndirect(rawDesired.DisplayName) {
		rawNew.DisplayName = rawDesired.DisplayName
	} else {
		if dcl.StringCanonicalize(rawDesired.DisplayName, rawNew.DisplayName) {
			rawNew.DisplayName = rawDesired.DisplayName
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Uid) && dcl.IsEmptyValueIndirect(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Annotations) && dcl.IsEmptyValueIndirect(rawDesired.Annotations) {
		rawNew.Annotations = rawDesired.Annotations
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

	if dcl.IsEmptyValueIndirect(rawNew.DeleteTime) && dcl.IsEmptyValueIndirect(rawDesired.DeleteTime) {
		rawNew.DeleteTime = rawDesired.DeleteTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PrivatePoolV1Config) && dcl.IsEmptyValueIndirect(rawDesired.PrivatePoolV1Config) {
		rawNew.PrivatePoolV1Config = rawDesired.PrivatePoolV1Config
	} else {
		rawNew.PrivatePoolV1Config = canonicalizeNewWorkerPoolPrivatePoolV1Config(c, rawDesired.PrivatePoolV1Config, rawNew.PrivatePoolV1Config)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Etag) && dcl.IsEmptyValueIndirect(rawDesired.Etag) {
		rawNew.Etag = rawDesired.Etag
	} else {
		if dcl.StringCanonicalize(rawDesired.Etag, rawNew.Etag) {
			rawNew.Etag = rawDesired.Etag
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.WorkerConfig) && dcl.IsEmptyValueIndirect(rawDesired.WorkerConfig) {
		rawNew.WorkerConfig = rawDesired.WorkerConfig
	} else {
		rawNew.WorkerConfig = canonicalizeNewWorkerPoolWorkerConfig(c, rawDesired.WorkerConfig, rawNew.WorkerConfig)
	}

	if dcl.IsEmptyValueIndirect(rawNew.NetworkConfig) && dcl.IsEmptyValueIndirect(rawDesired.NetworkConfig) {
		rawNew.NetworkConfig = rawDesired.NetworkConfig
	} else {
		rawNew.NetworkConfig = canonicalizeNewWorkerPoolNetworkConfig(c, rawDesired.NetworkConfig, rawNew.NetworkConfig)
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeWorkerPoolPrivatePoolV1Config(des, initial *WorkerPoolPrivatePoolV1Config, opts ...dcl.ApplyOption) *WorkerPoolPrivatePoolV1Config {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkerPoolPrivatePoolV1Config{}

	cDes.WorkerConfig = canonicalizeWorkerPoolPrivatePoolV1ConfigWorkerConfig(des.WorkerConfig, initial.WorkerConfig, opts...)
	cDes.NetworkConfig = canonicalizeWorkerPoolPrivatePoolV1ConfigNetworkConfig(des.NetworkConfig, initial.NetworkConfig, opts...)

	return cDes
}

func canonicalizeWorkerPoolPrivatePoolV1ConfigSlice(des, initial []WorkerPoolPrivatePoolV1Config, opts ...dcl.ApplyOption) []WorkerPoolPrivatePoolV1Config {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkerPoolPrivatePoolV1Config, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkerPoolPrivatePoolV1Config(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkerPoolPrivatePoolV1Config, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkerPoolPrivatePoolV1Config(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkerPoolPrivatePoolV1Config(c *Client, des, nw *WorkerPoolPrivatePoolV1Config) *WorkerPoolPrivatePoolV1Config {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkerPoolPrivatePoolV1Config while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.WorkerConfig = canonicalizeNewWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, des.WorkerConfig, nw.WorkerConfig)
	nw.NetworkConfig = canonicalizeNewWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, des.NetworkConfig, nw.NetworkConfig)

	return nw
}

func canonicalizeNewWorkerPoolPrivatePoolV1ConfigSet(c *Client, des, nw []WorkerPoolPrivatePoolV1Config) []WorkerPoolPrivatePoolV1Config {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkerPoolPrivatePoolV1Config
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkerPoolPrivatePoolV1ConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkerPoolPrivatePoolV1Config(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkerPoolPrivatePoolV1ConfigSlice(c *Client, des, nw []WorkerPoolPrivatePoolV1Config) []WorkerPoolPrivatePoolV1Config {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkerPoolPrivatePoolV1Config
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkerPoolPrivatePoolV1Config(c, &d, &n))
	}

	return items
}

func canonicalizeWorkerPoolPrivatePoolV1ConfigWorkerConfig(des, initial *WorkerPoolPrivatePoolV1ConfigWorkerConfig, opts ...dcl.ApplyOption) *WorkerPoolPrivatePoolV1ConfigWorkerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkerPoolPrivatePoolV1ConfigWorkerConfig{}

	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	if dcl.IsZeroValue(des.DiskSizeGb) || (dcl.IsEmptyValueIndirect(des.DiskSizeGb) && dcl.IsEmptyValueIndirect(initial.DiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DiskSizeGb = initial.DiskSizeGb
	} else {
		cDes.DiskSizeGb = des.DiskSizeGb
	}

	return cDes
}

func canonicalizeWorkerPoolPrivatePoolV1ConfigWorkerConfigSlice(des, initial []WorkerPoolPrivatePoolV1ConfigWorkerConfig, opts ...dcl.ApplyOption) []WorkerPoolPrivatePoolV1ConfigWorkerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkerPoolPrivatePoolV1ConfigWorkerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkerPoolPrivatePoolV1ConfigWorkerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkerPoolPrivatePoolV1ConfigWorkerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkerPoolPrivatePoolV1ConfigWorkerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkerPoolPrivatePoolV1ConfigWorkerConfig(c *Client, des, nw *WorkerPoolPrivatePoolV1ConfigWorkerConfig) *WorkerPoolPrivatePoolV1ConfigWorkerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkerPoolPrivatePoolV1ConfigWorkerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}

	return nw
}

func canonicalizeNewWorkerPoolPrivatePoolV1ConfigWorkerConfigSet(c *Client, des, nw []WorkerPoolPrivatePoolV1ConfigWorkerConfig) []WorkerPoolPrivatePoolV1ConfigWorkerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkerPoolPrivatePoolV1ConfigWorkerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkerPoolPrivatePoolV1ConfigWorkerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkerPoolPrivatePoolV1ConfigWorkerConfigSlice(c *Client, des, nw []WorkerPoolPrivatePoolV1ConfigWorkerConfig) []WorkerPoolPrivatePoolV1ConfigWorkerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkerPoolPrivatePoolV1ConfigWorkerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkerPoolPrivatePoolV1ConfigNetworkConfig(des, initial *WorkerPoolPrivatePoolV1ConfigNetworkConfig, opts ...dcl.ApplyOption) *WorkerPoolPrivatePoolV1ConfigNetworkConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkerPoolPrivatePoolV1ConfigNetworkConfig{}

	if dcl.IsZeroValue(des.PeeredNetwork) || (dcl.IsEmptyValueIndirect(des.PeeredNetwork) && dcl.IsEmptyValueIndirect(initial.PeeredNetwork)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.PeeredNetwork = initial.PeeredNetwork
	} else {
		cDes.PeeredNetwork = des.PeeredNetwork
	}
	if dcl.StringCanonicalize(des.PeeredNetworkIPRange, initial.PeeredNetworkIPRange) || dcl.IsZeroValue(des.PeeredNetworkIPRange) {
		cDes.PeeredNetworkIPRange = initial.PeeredNetworkIPRange
	} else {
		cDes.PeeredNetworkIPRange = des.PeeredNetworkIPRange
	}
	if dcl.IsZeroValue(des.EgressOption) || (dcl.IsEmptyValueIndirect(des.EgressOption) && dcl.IsEmptyValueIndirect(initial.EgressOption)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.EgressOption = initial.EgressOption
	} else {
		cDes.EgressOption = des.EgressOption
	}

	return cDes
}

func canonicalizeWorkerPoolPrivatePoolV1ConfigNetworkConfigSlice(des, initial []WorkerPoolPrivatePoolV1ConfigNetworkConfig, opts ...dcl.ApplyOption) []WorkerPoolPrivatePoolV1ConfigNetworkConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkerPoolPrivatePoolV1ConfigNetworkConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkerPoolPrivatePoolV1ConfigNetworkConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkerPoolPrivatePoolV1ConfigNetworkConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkerPoolPrivatePoolV1ConfigNetworkConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkerPoolPrivatePoolV1ConfigNetworkConfig(c *Client, des, nw *WorkerPoolPrivatePoolV1ConfigNetworkConfig) *WorkerPoolPrivatePoolV1ConfigNetworkConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkerPoolPrivatePoolV1ConfigNetworkConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.PeeredNetworkIPRange, nw.PeeredNetworkIPRange) {
		nw.PeeredNetworkIPRange = des.PeeredNetworkIPRange
	}

	return nw
}

func canonicalizeNewWorkerPoolPrivatePoolV1ConfigNetworkConfigSet(c *Client, des, nw []WorkerPoolPrivatePoolV1ConfigNetworkConfig) []WorkerPoolPrivatePoolV1ConfigNetworkConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkerPoolPrivatePoolV1ConfigNetworkConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkerPoolPrivatePoolV1ConfigNetworkConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkerPoolPrivatePoolV1ConfigNetworkConfigSlice(c *Client, des, nw []WorkerPoolPrivatePoolV1ConfigNetworkConfig) []WorkerPoolPrivatePoolV1ConfigNetworkConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkerPoolPrivatePoolV1ConfigNetworkConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkerPoolWorkerConfig(des, initial *WorkerPoolWorkerConfig, opts ...dcl.ApplyOption) *WorkerPoolWorkerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkerPoolWorkerConfig{}

	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	if dcl.IsZeroValue(des.DiskSizeGb) || (dcl.IsEmptyValueIndirect(des.DiskSizeGb) && dcl.IsEmptyValueIndirect(initial.DiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DiskSizeGb = initial.DiskSizeGb
	} else {
		cDes.DiskSizeGb = des.DiskSizeGb
	}
	if dcl.BoolCanonicalize(des.NoExternalIP, initial.NoExternalIP) || dcl.IsZeroValue(des.NoExternalIP) {
		cDes.NoExternalIP = initial.NoExternalIP
	} else {
		cDes.NoExternalIP = des.NoExternalIP
	}

	return cDes
}

func canonicalizeWorkerPoolWorkerConfigSlice(des, initial []WorkerPoolWorkerConfig, opts ...dcl.ApplyOption) []WorkerPoolWorkerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkerPoolWorkerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkerPoolWorkerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkerPoolWorkerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkerPoolWorkerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkerPoolWorkerConfig(c *Client, des, nw *WorkerPoolWorkerConfig) *WorkerPoolWorkerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkerPoolWorkerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	if dcl.BoolCanonicalize(des.NoExternalIP, nw.NoExternalIP) {
		nw.NoExternalIP = des.NoExternalIP
	}

	return nw
}

func canonicalizeNewWorkerPoolWorkerConfigSet(c *Client, des, nw []WorkerPoolWorkerConfig) []WorkerPoolWorkerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkerPoolWorkerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkerPoolWorkerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkerPoolWorkerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkerPoolWorkerConfigSlice(c *Client, des, nw []WorkerPoolWorkerConfig) []WorkerPoolWorkerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkerPoolWorkerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkerPoolWorkerConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkerPoolNetworkConfig(des, initial *WorkerPoolNetworkConfig, opts ...dcl.ApplyOption) *WorkerPoolNetworkConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkerPoolNetworkConfig{}

	if dcl.IsZeroValue(des.PeeredNetwork) || (dcl.IsEmptyValueIndirect(des.PeeredNetwork) && dcl.IsEmptyValueIndirect(initial.PeeredNetwork)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.PeeredNetwork = initial.PeeredNetwork
	} else {
		cDes.PeeredNetwork = des.PeeredNetwork
	}
	if dcl.StringCanonicalize(des.PeeredNetworkIPRange, initial.PeeredNetworkIPRange) || dcl.IsZeroValue(des.PeeredNetworkIPRange) {
		cDes.PeeredNetworkIPRange = initial.PeeredNetworkIPRange
	} else {
		cDes.PeeredNetworkIPRange = des.PeeredNetworkIPRange
	}

	return cDes
}

func canonicalizeWorkerPoolNetworkConfigSlice(des, initial []WorkerPoolNetworkConfig, opts ...dcl.ApplyOption) []WorkerPoolNetworkConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkerPoolNetworkConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkerPoolNetworkConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkerPoolNetworkConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkerPoolNetworkConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkerPoolNetworkConfig(c *Client, des, nw *WorkerPoolNetworkConfig) *WorkerPoolNetworkConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkerPoolNetworkConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.PeeredNetworkIPRange, nw.PeeredNetworkIPRange) {
		nw.PeeredNetworkIPRange = des.PeeredNetworkIPRange
	}

	return nw
}

func canonicalizeNewWorkerPoolNetworkConfigSet(c *Client, des, nw []WorkerPoolNetworkConfig) []WorkerPoolNetworkConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkerPoolNetworkConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkerPoolNetworkConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkerPoolNetworkConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkerPoolNetworkConfigSlice(c *Client, des, nw []WorkerPoolNetworkConfig) []WorkerPoolNetworkConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkerPoolNetworkConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkerPoolNetworkConfig(c, &d, &n))
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
func diffWorkerPool(c *Client, desired, actual *WorkerPool, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Annotations, actual.Annotations, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("Annotations")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DeleteTime, actual.DeleteTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DeleteTime")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.PrivatePoolV1Config, actual.PrivatePoolV1Config, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkerPoolPrivatePoolV1ConfigNewStyle, EmptyObject: EmptyWorkerPoolPrivatePoolV1Config, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrivatePoolV1Config")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.WorkerConfig, actual.WorkerConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkerPoolWorkerConfigNewStyle, EmptyObject: EmptyWorkerPoolWorkerConfig, OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("WorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NetworkConfig, actual.NetworkConfig, dcl.DiffInfo{ObjectFunction: compareWorkerPoolNetworkConfigNewStyle, EmptyObject: EmptyWorkerPoolNetworkConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkConfig")); len(ds) != 0 || err != nil {
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
func compareWorkerPoolPrivatePoolV1ConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkerPoolPrivatePoolV1Config)
	if !ok {
		desiredNotPointer, ok := d.(WorkerPoolPrivatePoolV1Config)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolPrivatePoolV1Config or *WorkerPoolPrivatePoolV1Config", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkerPoolPrivatePoolV1Config)
	if !ok {
		actualNotPointer, ok := a.(WorkerPoolPrivatePoolV1Config)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolPrivatePoolV1Config", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.WorkerConfig, actual.WorkerConfig, dcl.DiffInfo{ObjectFunction: compareWorkerPoolPrivatePoolV1ConfigWorkerConfigNewStyle, EmptyObject: EmptyWorkerPoolPrivatePoolV1ConfigWorkerConfig, OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("WorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NetworkConfig, actual.NetworkConfig, dcl.DiffInfo{ObjectFunction: compareWorkerPoolPrivatePoolV1ConfigNetworkConfigNewStyle, EmptyObject: EmptyWorkerPoolPrivatePoolV1ConfigNetworkConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkerPoolPrivatePoolV1ConfigWorkerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkerPoolPrivatePoolV1ConfigWorkerConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkerPoolPrivatePoolV1ConfigWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolPrivatePoolV1ConfigWorkerConfig or *WorkerPoolPrivatePoolV1ConfigWorkerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkerPoolPrivatePoolV1ConfigWorkerConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkerPoolPrivatePoolV1ConfigWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolPrivatePoolV1ConfigWorkerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("MachineType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskSizeGb, actual.DiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("DiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkerPoolPrivatePoolV1ConfigNetworkConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkerPoolPrivatePoolV1ConfigNetworkConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkerPoolPrivatePoolV1ConfigNetworkConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolPrivatePoolV1ConfigNetworkConfig or *WorkerPoolPrivatePoolV1ConfigNetworkConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkerPoolPrivatePoolV1ConfigNetworkConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkerPoolPrivatePoolV1ConfigNetworkConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolPrivatePoolV1ConfigNetworkConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.PeeredNetwork, actual.PeeredNetwork, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeeredNetwork")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PeeredNetworkIPRange, actual.PeeredNetworkIPRange, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeeredNetworkIpRange")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EgressOption, actual.EgressOption, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("EgressOption")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkerPoolWorkerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkerPoolWorkerConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkerPoolWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolWorkerConfig or *WorkerPoolWorkerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkerPoolWorkerConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkerPoolWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolWorkerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("MachineType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskSizeGb, actual.DiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("DiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NoExternalIP, actual.NoExternalIP, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateWorkerPoolUpdateWorkerPoolOperation")}, fn.AddNest("NoExternalIp")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkerPoolNetworkConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkerPoolNetworkConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkerPoolNetworkConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolNetworkConfig or *WorkerPoolNetworkConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkerPoolNetworkConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkerPoolNetworkConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkerPoolNetworkConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.PeeredNetwork, actual.PeeredNetwork, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeeredNetwork")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PeeredNetworkIPRange, actual.PeeredNetworkIPRange, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PeeredNetworkIpRange")); len(ds) != 0 || err != nil {
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
func (r *WorkerPool) urlNormalized() *WorkerPool {
	normalized := dcl.Copy(*r).(WorkerPool)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.DisplayName = dcl.SelfLinkToName(r.DisplayName)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.Etag = dcl.SelfLinkToName(r.Etag)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *WorkerPool) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateWorkerPool" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/workerPools/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the WorkerPool resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *WorkerPool) marshal(c *Client) ([]byte, error) {
	m, err := expandWorkerPool(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling WorkerPool: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalWorkerPool decodes JSON responses into the WorkerPool resource schema.
func unmarshalWorkerPool(b []byte, c *Client, res *WorkerPool) (*WorkerPool, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapWorkerPool(m, c, res)
}

func unmarshalMapWorkerPool(m map[string]interface{}, c *Client, res *WorkerPool) (*WorkerPool, error) {

	flattened := flattenWorkerPool(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandWorkerPool expands WorkerPool into a JSON request object.
func expandWorkerPool(c *Client, f *WorkerPool) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/workerPools/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.DisplayName; dcl.ValueShouldBeSent(v) {
		m["displayName"] = v
	}
	if v := f.Annotations; dcl.ValueShouldBeSent(v) {
		m["annotations"] = v
	}
	if v, err := expandWorkerPoolPrivatePoolV1Config(c, f.PrivatePoolV1Config, res); err != nil {
		return nil, fmt.Errorf("error expanding PrivatePoolV1Config into privatePoolV1Config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["privatePoolV1Config"] = v
	}
	if v, err := expandWorkerPoolWorkerConfig(c, f.WorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding WorkerConfig into workerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["workerConfig"] = v
	}
	if v, err := expandWorkerPoolNetworkConfig(c, f.NetworkConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding NetworkConfig into networkConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["networkConfig"] = v
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

// flattenWorkerPool flattens WorkerPool from a JSON request object into the
// WorkerPool type.
func flattenWorkerPool(c *Client, i interface{}, res *WorkerPool) *WorkerPool {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &WorkerPool{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.DisplayName = dcl.FlattenString(m["displayName"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.Annotations = dcl.FlattenKeyValuePairs(m["annotations"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.DeleteTime = dcl.FlattenString(m["deleteTime"])
	resultRes.State = flattenWorkerPoolStateEnum(m["state"])
	resultRes.PrivatePoolV1Config = flattenWorkerPoolPrivatePoolV1Config(c, m["privatePoolV1Config"], res)
	resultRes.Etag = dcl.FlattenString(m["etag"])
	resultRes.WorkerConfig = flattenWorkerPoolWorkerConfig(c, m["workerConfig"], res)
	resultRes.NetworkConfig = flattenWorkerPoolNetworkConfig(c, m["networkConfig"], res)
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandWorkerPoolPrivatePoolV1ConfigMap expands the contents of WorkerPoolPrivatePoolV1Config into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1ConfigMap(c *Client, f map[string]WorkerPoolPrivatePoolV1Config, res *WorkerPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkerPoolPrivatePoolV1Config(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkerPoolPrivatePoolV1ConfigSlice expands the contents of WorkerPoolPrivatePoolV1Config into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1ConfigSlice(c *Client, f []WorkerPoolPrivatePoolV1Config, res *WorkerPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkerPoolPrivatePoolV1Config(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkerPoolPrivatePoolV1ConfigMap flattens the contents of WorkerPoolPrivatePoolV1Config from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigMap(c *Client, i interface{}, res *WorkerPool) map[string]WorkerPoolPrivatePoolV1Config {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkerPoolPrivatePoolV1Config{}
	}

	if len(a) == 0 {
		return map[string]WorkerPoolPrivatePoolV1Config{}
	}

	items := make(map[string]WorkerPoolPrivatePoolV1Config)
	for k, item := range a {
		items[k] = *flattenWorkerPoolPrivatePoolV1Config(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkerPoolPrivatePoolV1ConfigSlice flattens the contents of WorkerPoolPrivatePoolV1Config from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigSlice(c *Client, i interface{}, res *WorkerPool) []WorkerPoolPrivatePoolV1Config {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkerPoolPrivatePoolV1Config{}
	}

	if len(a) == 0 {
		return []WorkerPoolPrivatePoolV1Config{}
	}

	items := make([]WorkerPoolPrivatePoolV1Config, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkerPoolPrivatePoolV1Config(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkerPoolPrivatePoolV1Config expands an instance of WorkerPoolPrivatePoolV1Config into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1Config(c *Client, f *WorkerPoolPrivatePoolV1Config, res *WorkerPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, f.WorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding WorkerConfig into workerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["workerConfig"] = v
	}
	if v, err := expandWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, f.NetworkConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding NetworkConfig into networkConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["networkConfig"] = v
	}

	return m, nil
}

// flattenWorkerPoolPrivatePoolV1Config flattens an instance of WorkerPoolPrivatePoolV1Config from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1Config(c *Client, i interface{}, res *WorkerPool) *WorkerPoolPrivatePoolV1Config {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkerPoolPrivatePoolV1Config{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkerPoolPrivatePoolV1Config
	}
	r.WorkerConfig = flattenWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, m["workerConfig"], res)
	r.NetworkConfig = flattenWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, m["networkConfig"], res)

	return r
}

// expandWorkerPoolPrivatePoolV1ConfigWorkerConfigMap expands the contents of WorkerPoolPrivatePoolV1ConfigWorkerConfig into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1ConfigWorkerConfigMap(c *Client, f map[string]WorkerPoolPrivatePoolV1ConfigWorkerConfig, res *WorkerPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkerPoolPrivatePoolV1ConfigWorkerConfigSlice expands the contents of WorkerPoolPrivatePoolV1ConfigWorkerConfig into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1ConfigWorkerConfigSlice(c *Client, f []WorkerPoolPrivatePoolV1ConfigWorkerConfig, res *WorkerPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkerPoolPrivatePoolV1ConfigWorkerConfigMap flattens the contents of WorkerPoolPrivatePoolV1ConfigWorkerConfig from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigWorkerConfigMap(c *Client, i interface{}, res *WorkerPool) map[string]WorkerPoolPrivatePoolV1ConfigWorkerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkerPoolPrivatePoolV1ConfigWorkerConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkerPoolPrivatePoolV1ConfigWorkerConfig{}
	}

	items := make(map[string]WorkerPoolPrivatePoolV1ConfigWorkerConfig)
	for k, item := range a {
		items[k] = *flattenWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkerPoolPrivatePoolV1ConfigWorkerConfigSlice flattens the contents of WorkerPoolPrivatePoolV1ConfigWorkerConfig from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigWorkerConfigSlice(c *Client, i interface{}, res *WorkerPool) []WorkerPoolPrivatePoolV1ConfigWorkerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkerPoolPrivatePoolV1ConfigWorkerConfig{}
	}

	if len(a) == 0 {
		return []WorkerPoolPrivatePoolV1ConfigWorkerConfig{}
	}

	items := make([]WorkerPoolPrivatePoolV1ConfigWorkerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkerPoolPrivatePoolV1ConfigWorkerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkerPoolPrivatePoolV1ConfigWorkerConfig expands an instance of WorkerPoolPrivatePoolV1ConfigWorkerConfig into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1ConfigWorkerConfig(c *Client, f *WorkerPoolPrivatePoolV1ConfigWorkerConfig, res *WorkerPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MachineType; !dcl.IsEmptyValueIndirect(v) {
		m["machineType"] = v
	}
	if v := f.DiskSizeGb; !dcl.IsEmptyValueIndirect(v) {
		m["diskSizeGb"] = v
	}

	return m, nil
}

// flattenWorkerPoolPrivatePoolV1ConfigWorkerConfig flattens an instance of WorkerPoolPrivatePoolV1ConfigWorkerConfig from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigWorkerConfig(c *Client, i interface{}, res *WorkerPool) *WorkerPoolPrivatePoolV1ConfigWorkerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkerPoolPrivatePoolV1ConfigWorkerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkerPoolPrivatePoolV1ConfigWorkerConfig
	}
	r.MachineType = dcl.FlattenString(m["machineType"])
	r.DiskSizeGb = dcl.FlattenInteger(m["diskSizeGb"])

	return r
}

// expandWorkerPoolPrivatePoolV1ConfigNetworkConfigMap expands the contents of WorkerPoolPrivatePoolV1ConfigNetworkConfig into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1ConfigNetworkConfigMap(c *Client, f map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfig, res *WorkerPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkerPoolPrivatePoolV1ConfigNetworkConfigSlice expands the contents of WorkerPoolPrivatePoolV1ConfigNetworkConfig into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1ConfigNetworkConfigSlice(c *Client, f []WorkerPoolPrivatePoolV1ConfigNetworkConfig, res *WorkerPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigMap flattens the contents of WorkerPoolPrivatePoolV1ConfigNetworkConfig from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigMap(c *Client, i interface{}, res *WorkerPool) map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfig{}
	}

	items := make(map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfig)
	for k, item := range a {
		items[k] = *flattenWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigSlice flattens the contents of WorkerPoolPrivatePoolV1ConfigNetworkConfig from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigSlice(c *Client, i interface{}, res *WorkerPool) []WorkerPoolPrivatePoolV1ConfigNetworkConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkerPoolPrivatePoolV1ConfigNetworkConfig{}
	}

	if len(a) == 0 {
		return []WorkerPoolPrivatePoolV1ConfigNetworkConfig{}
	}

	items := make([]WorkerPoolPrivatePoolV1ConfigNetworkConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkerPoolPrivatePoolV1ConfigNetworkConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkerPoolPrivatePoolV1ConfigNetworkConfig expands an instance of WorkerPoolPrivatePoolV1ConfigNetworkConfig into a JSON
// request object.
func expandWorkerPoolPrivatePoolV1ConfigNetworkConfig(c *Client, f *WorkerPoolPrivatePoolV1ConfigNetworkConfig, res *WorkerPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.PeeredNetwork; !dcl.IsEmptyValueIndirect(v) {
		m["peeredNetwork"] = v
	}
	if v := f.PeeredNetworkIPRange; !dcl.IsEmptyValueIndirect(v) {
		m["peeredNetworkIpRange"] = v
	}
	if v := f.EgressOption; !dcl.IsEmptyValueIndirect(v) {
		m["egressOption"] = v
	}

	return m, nil
}

// flattenWorkerPoolPrivatePoolV1ConfigNetworkConfig flattens an instance of WorkerPoolPrivatePoolV1ConfigNetworkConfig from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigNetworkConfig(c *Client, i interface{}, res *WorkerPool) *WorkerPoolPrivatePoolV1ConfigNetworkConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkerPoolPrivatePoolV1ConfigNetworkConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkerPoolPrivatePoolV1ConfigNetworkConfig
	}
	r.PeeredNetwork = dcl.FlattenString(m["peeredNetwork"])
	r.PeeredNetworkIPRange = dcl.FlattenString(m["peeredNetworkIpRange"])
	r.EgressOption = flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum(m["egressOption"])

	return r
}

// expandWorkerPoolWorkerConfigMap expands the contents of WorkerPoolWorkerConfig into a JSON
// request object.
func expandWorkerPoolWorkerConfigMap(c *Client, f map[string]WorkerPoolWorkerConfig, res *WorkerPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkerPoolWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkerPoolWorkerConfigSlice expands the contents of WorkerPoolWorkerConfig into a JSON
// request object.
func expandWorkerPoolWorkerConfigSlice(c *Client, f []WorkerPoolWorkerConfig, res *WorkerPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkerPoolWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkerPoolWorkerConfigMap flattens the contents of WorkerPoolWorkerConfig from a JSON
// response object.
func flattenWorkerPoolWorkerConfigMap(c *Client, i interface{}, res *WorkerPool) map[string]WorkerPoolWorkerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkerPoolWorkerConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkerPoolWorkerConfig{}
	}

	items := make(map[string]WorkerPoolWorkerConfig)
	for k, item := range a {
		items[k] = *flattenWorkerPoolWorkerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkerPoolWorkerConfigSlice flattens the contents of WorkerPoolWorkerConfig from a JSON
// response object.
func flattenWorkerPoolWorkerConfigSlice(c *Client, i interface{}, res *WorkerPool) []WorkerPoolWorkerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkerPoolWorkerConfig{}
	}

	if len(a) == 0 {
		return []WorkerPoolWorkerConfig{}
	}

	items := make([]WorkerPoolWorkerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkerPoolWorkerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkerPoolWorkerConfig expands an instance of WorkerPoolWorkerConfig into a JSON
// request object.
func expandWorkerPoolWorkerConfig(c *Client, f *WorkerPoolWorkerConfig, res *WorkerPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MachineType; !dcl.IsEmptyValueIndirect(v) {
		m["machineType"] = v
	}
	if v := f.DiskSizeGb; !dcl.IsEmptyValueIndirect(v) {
		m["diskSizeGb"] = v
	}
	if v := f.NoExternalIP; !dcl.IsEmptyValueIndirect(v) {
		m["noExternalIp"] = v
	}

	return m, nil
}

// flattenWorkerPoolWorkerConfig flattens an instance of WorkerPoolWorkerConfig from a JSON
// response object.
func flattenWorkerPoolWorkerConfig(c *Client, i interface{}, res *WorkerPool) *WorkerPoolWorkerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkerPoolWorkerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkerPoolWorkerConfig
	}
	r.MachineType = dcl.FlattenString(m["machineType"])
	r.DiskSizeGb = dcl.FlattenInteger(m["diskSizeGb"])
	r.NoExternalIP = dcl.FlattenBool(m["noExternalIp"])

	return r
}

// expandWorkerPoolNetworkConfigMap expands the contents of WorkerPoolNetworkConfig into a JSON
// request object.
func expandWorkerPoolNetworkConfigMap(c *Client, f map[string]WorkerPoolNetworkConfig, res *WorkerPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkerPoolNetworkConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkerPoolNetworkConfigSlice expands the contents of WorkerPoolNetworkConfig into a JSON
// request object.
func expandWorkerPoolNetworkConfigSlice(c *Client, f []WorkerPoolNetworkConfig, res *WorkerPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkerPoolNetworkConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkerPoolNetworkConfigMap flattens the contents of WorkerPoolNetworkConfig from a JSON
// response object.
func flattenWorkerPoolNetworkConfigMap(c *Client, i interface{}, res *WorkerPool) map[string]WorkerPoolNetworkConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkerPoolNetworkConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkerPoolNetworkConfig{}
	}

	items := make(map[string]WorkerPoolNetworkConfig)
	for k, item := range a {
		items[k] = *flattenWorkerPoolNetworkConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkerPoolNetworkConfigSlice flattens the contents of WorkerPoolNetworkConfig from a JSON
// response object.
func flattenWorkerPoolNetworkConfigSlice(c *Client, i interface{}, res *WorkerPool) []WorkerPoolNetworkConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkerPoolNetworkConfig{}
	}

	if len(a) == 0 {
		return []WorkerPoolNetworkConfig{}
	}

	items := make([]WorkerPoolNetworkConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkerPoolNetworkConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkerPoolNetworkConfig expands an instance of WorkerPoolNetworkConfig into a JSON
// request object.
func expandWorkerPoolNetworkConfig(c *Client, f *WorkerPoolNetworkConfig, res *WorkerPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.PeeredNetwork; !dcl.IsEmptyValueIndirect(v) {
		m["peeredNetwork"] = v
	}
	if v := f.PeeredNetworkIPRange; !dcl.IsEmptyValueIndirect(v) {
		m["peeredNetworkIpRange"] = v
	}

	return m, nil
}

// flattenWorkerPoolNetworkConfig flattens an instance of WorkerPoolNetworkConfig from a JSON
// response object.
func flattenWorkerPoolNetworkConfig(c *Client, i interface{}, res *WorkerPool) *WorkerPoolNetworkConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkerPoolNetworkConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkerPoolNetworkConfig
	}
	r.PeeredNetwork = dcl.FlattenString(m["peeredNetwork"])
	r.PeeredNetworkIPRange = dcl.FlattenString(m["peeredNetworkIpRange"])

	return r
}

// flattenWorkerPoolStateEnumMap flattens the contents of WorkerPoolStateEnum from a JSON
// response object.
func flattenWorkerPoolStateEnumMap(c *Client, i interface{}, res *WorkerPool) map[string]WorkerPoolStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkerPoolStateEnum{}
	}

	if len(a) == 0 {
		return map[string]WorkerPoolStateEnum{}
	}

	items := make(map[string]WorkerPoolStateEnum)
	for k, item := range a {
		items[k] = *flattenWorkerPoolStateEnum(item.(interface{}))
	}

	return items
}

// flattenWorkerPoolStateEnumSlice flattens the contents of WorkerPoolStateEnum from a JSON
// response object.
func flattenWorkerPoolStateEnumSlice(c *Client, i interface{}, res *WorkerPool) []WorkerPoolStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkerPoolStateEnum{}
	}

	if len(a) == 0 {
		return []WorkerPoolStateEnum{}
	}

	items := make([]WorkerPoolStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkerPoolStateEnum(item.(interface{})))
	}

	return items
}

// flattenWorkerPoolStateEnum asserts that an interface is a string, and returns a
// pointer to a *WorkerPoolStateEnum with the same value as that string.
func flattenWorkerPoolStateEnum(i interface{}) *WorkerPoolStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return WorkerPoolStateEnumRef(s)
}

// flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumMap flattens the contents of WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumMap(c *Client, i interface{}, res *WorkerPool) map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum{}
	}

	if len(a) == 0 {
		return map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum{}
	}

	items := make(map[string]WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum)
	for k, item := range a {
		items[k] = *flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum(item.(interface{}))
	}

	return items
}

// flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumSlice flattens the contents of WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum from a JSON
// response object.
func flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumSlice(c *Client, i interface{}, res *WorkerPool) []WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum{}
	}

	if len(a) == 0 {
		return []WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum{}
	}

	items := make([]WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum(item.(interface{})))
	}

	return items
}

// flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum asserts that an interface is a string, and returns a
// pointer to a *WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum with the same value as that string.
func flattenWorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum(i interface{}) *WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *WorkerPool) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalWorkerPool(b, c, r)
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

type workerPoolDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         workerPoolApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToWorkerPoolDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]workerPoolDiff, error) {
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
	var diffs []workerPoolDiff
	// For each operation name, create a workerPoolDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := workerPoolDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToWorkerPoolApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToWorkerPoolApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (workerPoolApiOperation, error) {
	switch opName {

	case "updateWorkerPoolUpdateWorkerPoolOperation":
		return &updateWorkerPoolUpdateWorkerPoolOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractWorkerPoolFields(r *WorkerPool) error {
	if dcl.IsEmptyValueIndirect(r.PrivatePoolV1Config) {
		r.PrivatePoolV1Config = betaToGaPrivatePool(r, r.PrivatePoolV1Config)
	}
	vPrivatePoolV1Config := r.PrivatePoolV1Config
	if vPrivatePoolV1Config == nil {
		// note: explicitly not the empty object.
		vPrivatePoolV1Config = &WorkerPoolPrivatePoolV1Config{}
	}
	if err := extractWorkerPoolPrivatePoolV1ConfigFields(r, vPrivatePoolV1Config); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPrivatePoolV1Config) {
		r.PrivatePoolV1Config = vPrivatePoolV1Config
	}
	vWorkerConfig := r.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &WorkerPoolWorkerConfig{}
	}
	if err := extractWorkerPoolWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		r.WorkerConfig = vWorkerConfig
	}
	vNetworkConfig := r.NetworkConfig
	if vNetworkConfig == nil {
		// note: explicitly not the empty object.
		vNetworkConfig = &WorkerPoolNetworkConfig{}
	}
	if err := extractWorkerPoolNetworkConfigFields(r, vNetworkConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetworkConfig) {
		r.NetworkConfig = vNetworkConfig
	}
	return nil
}
func extractWorkerPoolPrivatePoolV1ConfigFields(r *WorkerPool, o *WorkerPoolPrivatePoolV1Config) error {
	vWorkerConfig := o.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &WorkerPoolPrivatePoolV1ConfigWorkerConfig{}
	}
	if err := extractWorkerPoolPrivatePoolV1ConfigWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		o.WorkerConfig = vWorkerConfig
	}
	vNetworkConfig := o.NetworkConfig
	if vNetworkConfig == nil {
		// note: explicitly not the empty object.
		vNetworkConfig = &WorkerPoolPrivatePoolV1ConfigNetworkConfig{}
	}
	if err := extractWorkerPoolPrivatePoolV1ConfigNetworkConfigFields(r, vNetworkConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetworkConfig) {
		o.NetworkConfig = vNetworkConfig
	}
	return nil
}
func extractWorkerPoolPrivatePoolV1ConfigWorkerConfigFields(r *WorkerPool, o *WorkerPoolPrivatePoolV1ConfigWorkerConfig) error {
	return nil
}
func extractWorkerPoolPrivatePoolV1ConfigNetworkConfigFields(r *WorkerPool, o *WorkerPoolPrivatePoolV1ConfigNetworkConfig) error {
	return nil
}
func extractWorkerPoolWorkerConfigFields(r *WorkerPool, o *WorkerPoolWorkerConfig) error {
	return nil
}
func extractWorkerPoolNetworkConfigFields(r *WorkerPool, o *WorkerPoolNetworkConfig) error {
	return nil
}

func postReadExtractWorkerPoolFields(r *WorkerPool) error {

	r.PrivatePoolV1Config = gaToBetaPrivatePool(r, r.PrivatePoolV1Config)
	vPrivatePoolV1Config := r.PrivatePoolV1Config
	if vPrivatePoolV1Config == nil {
		// note: explicitly not the empty object.
		vPrivatePoolV1Config = &WorkerPoolPrivatePoolV1Config{}
	}
	if err := postReadExtractWorkerPoolPrivatePoolV1ConfigFields(r, vPrivatePoolV1Config); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPrivatePoolV1Config) {
		r.PrivatePoolV1Config = vPrivatePoolV1Config
	}
	vWorkerConfig := r.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &WorkerPoolWorkerConfig{}
	}
	if err := postReadExtractWorkerPoolWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		r.WorkerConfig = vWorkerConfig
	}
	vNetworkConfig := r.NetworkConfig
	if vNetworkConfig == nil {
		// note: explicitly not the empty object.
		vNetworkConfig = &WorkerPoolNetworkConfig{}
	}
	if err := postReadExtractWorkerPoolNetworkConfigFields(r, vNetworkConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetworkConfig) {
		r.NetworkConfig = vNetworkConfig
	}
	return nil
}
func postReadExtractWorkerPoolPrivatePoolV1ConfigFields(r *WorkerPool, o *WorkerPoolPrivatePoolV1Config) error {
	vWorkerConfig := o.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &WorkerPoolPrivatePoolV1ConfigWorkerConfig{}
	}
	if err := extractWorkerPoolPrivatePoolV1ConfigWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		o.WorkerConfig = vWorkerConfig
	}
	vNetworkConfig := o.NetworkConfig
	if vNetworkConfig == nil {
		// note: explicitly not the empty object.
		vNetworkConfig = &WorkerPoolPrivatePoolV1ConfigNetworkConfig{}
	}
	if err := extractWorkerPoolPrivatePoolV1ConfigNetworkConfigFields(r, vNetworkConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetworkConfig) {
		o.NetworkConfig = vNetworkConfig
	}
	return nil
}
func postReadExtractWorkerPoolPrivatePoolV1ConfigWorkerConfigFields(r *WorkerPool, o *WorkerPoolPrivatePoolV1ConfigWorkerConfig) error {
	return nil
}
func postReadExtractWorkerPoolPrivatePoolV1ConfigNetworkConfigFields(r *WorkerPool, o *WorkerPoolPrivatePoolV1ConfigNetworkConfig) error {
	return nil
}
func postReadExtractWorkerPoolWorkerConfigFields(r *WorkerPool, o *WorkerPoolWorkerConfig) error {
	return nil
}
func postReadExtractWorkerPoolNetworkConfigFields(r *WorkerPool, o *WorkerPoolNetworkConfig) error {
	return nil
}
