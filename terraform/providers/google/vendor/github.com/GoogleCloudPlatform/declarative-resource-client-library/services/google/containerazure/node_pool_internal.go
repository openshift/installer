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
package containerazure

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

func (r *NodePool) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "version"); err != nil {
		return err
	}
	if err := dcl.Required(r, "config"); err != nil {
		return err
	}
	if err := dcl.Required(r, "subnetId"); err != nil {
		return err
	}
	if err := dcl.Required(r, "autoscaling"); err != nil {
		return err
	}
	if err := dcl.Required(r, "maxPodsConstraint"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Cluster, "Cluster"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Config) {
		if err := r.Config.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Autoscaling) {
		if err := r.Autoscaling.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.MaxPodsConstraint) {
		if err := r.MaxPodsConstraint.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *NodePoolConfig) validate() error {
	if err := dcl.Required(r, "sshConfig"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.RootVolume) {
		if err := r.RootVolume.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SshConfig) {
		if err := r.SshConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ProxyConfig) {
		if err := r.ProxyConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *NodePoolConfigRootVolume) validate() error {
	return nil
}
func (r *NodePoolConfigSshConfig) validate() error {
	if err := dcl.Required(r, "authorizedKey"); err != nil {
		return err
	}
	return nil
}
func (r *NodePoolConfigProxyConfig) validate() error {
	if err := dcl.Required(r, "resourceGroupId"); err != nil {
		return err
	}
	if err := dcl.Required(r, "secretId"); err != nil {
		return err
	}
	return nil
}
func (r *NodePoolAutoscaling) validate() error {
	if err := dcl.Required(r, "minNodeCount"); err != nil {
		return err
	}
	if err := dcl.Required(r, "maxNodeCount"); err != nil {
		return err
	}
	return nil
}
func (r *NodePoolMaxPodsConstraint) validate() error {
	if err := dcl.Required(r, "maxPodsPerNode"); err != nil {
		return err
	}
	return nil
}
func (r *NodePool) basePath() string {
	params := map[string]interface{}{
		"location": dcl.ValueOrEmptyString(r.Location),
	}
	return dcl.Nprintf("https://{{location}}-gkemulticloud.googleapis.com/v1", params)
}

func (r *NodePool) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/azureClusters/{{cluster}}/azureNodePools/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *NodePool) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/azureClusters/{{cluster}}/azureNodePools", nr.basePath(), userBasePath, params), nil

}

func (r *NodePool) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/azureClusters/{{cluster}}/azureNodePools?azureNodePoolId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *NodePool) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/azureClusters/{{cluster}}/azureNodePools/{{name}}", nr.basePath(), userBasePath, params), nil
}

// nodePoolApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type nodePoolApiOperation interface {
	do(context.Context, *NodePool, *Client) error
}

// newUpdateNodePoolUpdateAzureNodePoolRequest creates a request for an
// NodePool resource's UpdateAzureNodePool update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateNodePoolUpdateAzureNodePoolRequest(ctx context.Context, f *NodePool, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Version; !dcl.IsEmptyValueIndirect(v) {
		req["version"] = v
	}
	if v, err := expandNodePoolConfig(c, f.Config, res); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["config"] = v
	}
	if v, err := expandNodePoolAutoscaling(c, f.Autoscaling, res); err != nil {
		return nil, fmt.Errorf("error expanding Autoscaling into autoscaling: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["autoscaling"] = v
	}
	if v := f.Annotations; !dcl.IsEmptyValueIndirect(v) {
		req["annotations"] = v
	}
	b, err := c.getNodePoolRaw(ctx, f)
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

// marshalUpdateNodePoolUpdateAzureNodePoolRequest converts the update into
// the final JSON request body.
func marshalUpdateNodePoolUpdateAzureNodePoolRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateNodePoolUpdateAzureNodePoolOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateNodePoolUpdateAzureNodePoolOperation) do(ctx context.Context, r *NodePool, c *Client) error {
	_, err := c.GetNodePool(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateAzureNodePool")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateNodePoolUpdateAzureNodePoolRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateNodePoolUpdateAzureNodePoolRequest(c, req)
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

func (c *Client) listNodePoolRaw(ctx context.Context, r *NodePool, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != NodePoolMaxPage {
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

type listNodePoolOperation struct {
	AzureNodePools []map[string]interface{} `json:"azureNodePools"`
	Token          string                   `json:"nextPageToken"`
}

func (c *Client) listNodePool(ctx context.Context, r *NodePool, pageToken string, pageSize int32) ([]*NodePool, string, error) {
	b, err := c.listNodePoolRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listNodePoolOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*NodePool
	for _, v := range m.AzureNodePools {
		res, err := unmarshalMapNodePool(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		res.Cluster = r.Cluster
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllNodePool(ctx context.Context, f func(*NodePool) bool, resources []*NodePool) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteNodePool(ctx, res)
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

type deleteNodePoolOperation struct{}

func (op *deleteNodePoolOperation) do(ctx context.Context, r *NodePool, c *Client) error {
	r, err := c.GetNodePool(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "NodePool not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetNodePool checking for existence. error: %v", err)
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
		_, err := c.GetNodePool(ctx, r)
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
type createNodePoolOperation struct {
	response map[string]interface{}
}

func (op *createNodePoolOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createNodePoolOperation) do(ctx context.Context, r *NodePool, c *Client) error {
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

	if _, err := c.GetNodePool(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getNodePoolRaw(ctx context.Context, r *NodePool) ([]byte, error) {

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

func (c *Client) nodePoolDiffsForRawDesired(ctx context.Context, rawDesired *NodePool, opts ...dcl.ApplyOption) (initial, desired *NodePool, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *NodePool
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*NodePool); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected NodePool, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetNodePool(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a NodePool resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve NodePool resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that NodePool resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeNodePoolDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for NodePool: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for NodePool: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractNodePoolFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeNodePoolInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for NodePool: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeNodePoolDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for NodePool: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffNodePool(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeNodePoolInitialState(rawInitial, rawDesired *NodePool) (*NodePool, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeNodePoolDesiredState(rawDesired, rawInitial *NodePool, opts ...dcl.ApplyOption) (*NodePool, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Config = canonicalizeNodePoolConfig(rawDesired.Config, nil, opts...)
		rawDesired.Autoscaling = canonicalizeNodePoolAutoscaling(rawDesired.Autoscaling, nil, opts...)
		rawDesired.MaxPodsConstraint = canonicalizeNodePoolMaxPodsConstraint(rawDesired.MaxPodsConstraint, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &NodePool{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Version, rawInitial.Version) {
		canonicalDesired.Version = rawInitial.Version
	} else {
		canonicalDesired.Version = rawDesired.Version
	}
	canonicalDesired.Config = canonicalizeNodePoolConfig(rawDesired.Config, rawInitial.Config, opts...)
	if dcl.StringCanonicalize(rawDesired.SubnetId, rawInitial.SubnetId) {
		canonicalDesired.SubnetId = rawInitial.SubnetId
	} else {
		canonicalDesired.SubnetId = rawDesired.SubnetId
	}
	canonicalDesired.Autoscaling = canonicalizeNodePoolAutoscaling(rawDesired.Autoscaling, rawInitial.Autoscaling, opts...)
	if dcl.IsZeroValue(rawDesired.Annotations) || (dcl.IsEmptyValueIndirect(rawDesired.Annotations) && dcl.IsEmptyValueIndirect(rawInitial.Annotations)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Annotations = rawInitial.Annotations
	} else {
		canonicalDesired.Annotations = rawDesired.Annotations
	}
	canonicalDesired.MaxPodsConstraint = canonicalizeNodePoolMaxPodsConstraint(rawDesired.MaxPodsConstraint, rawInitial.MaxPodsConstraint, opts...)
	if dcl.StringCanonicalize(rawDesired.AzureAvailabilityZone, rawInitial.AzureAvailabilityZone) {
		canonicalDesired.AzureAvailabilityZone = rawInitial.AzureAvailabilityZone
	} else {
		canonicalDesired.AzureAvailabilityZone = rawDesired.AzureAvailabilityZone
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
	if dcl.NameToSelfLink(rawDesired.Cluster, rawInitial.Cluster) {
		canonicalDesired.Cluster = rawInitial.Cluster
	} else {
		canonicalDesired.Cluster = rawDesired.Cluster
	}
	return canonicalDesired, nil
}

func canonicalizeNodePoolNewState(c *Client, rawNew, rawDesired *NodePool) (*NodePool, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Version) && dcl.IsEmptyValueIndirect(rawDesired.Version) {
		rawNew.Version = rawDesired.Version
	} else {
		if dcl.StringCanonicalize(rawDesired.Version, rawNew.Version) {
			rawNew.Version = rawDesired.Version
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Config) && dcl.IsEmptyValueIndirect(rawDesired.Config) {
		rawNew.Config = rawDesired.Config
	} else {
		rawNew.Config = canonicalizeNewNodePoolConfig(c, rawDesired.Config, rawNew.Config)
	}

	if dcl.IsEmptyValueIndirect(rawNew.SubnetId) && dcl.IsEmptyValueIndirect(rawDesired.SubnetId) {
		rawNew.SubnetId = rawDesired.SubnetId
	} else {
		if dcl.StringCanonicalize(rawDesired.SubnetId, rawNew.SubnetId) {
			rawNew.SubnetId = rawDesired.SubnetId
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Autoscaling) && dcl.IsEmptyValueIndirect(rawDesired.Autoscaling) {
		rawNew.Autoscaling = rawDesired.Autoscaling
	} else {
		rawNew.Autoscaling = canonicalizeNewNodePoolAutoscaling(c, rawDesired.Autoscaling, rawNew.Autoscaling)
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Uid) && dcl.IsEmptyValueIndirect(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Reconciling) && dcl.IsEmptyValueIndirect(rawDesired.Reconciling) {
		rawNew.Reconciling = rawDesired.Reconciling
	} else {
		if dcl.BoolCanonicalize(rawDesired.Reconciling, rawNew.Reconciling) {
			rawNew.Reconciling = rawDesired.Reconciling
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

	if dcl.IsEmptyValueIndirect(rawNew.Etag) && dcl.IsEmptyValueIndirect(rawDesired.Etag) {
		rawNew.Etag = rawDesired.Etag
	} else {
		if dcl.StringCanonicalize(rawDesired.Etag, rawNew.Etag) {
			rawNew.Etag = rawDesired.Etag
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Annotations) && dcl.IsEmptyValueIndirect(rawDesired.Annotations) {
		rawNew.Annotations = rawDesired.Annotations
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.MaxPodsConstraint) && dcl.IsEmptyValueIndirect(rawDesired.MaxPodsConstraint) {
		rawNew.MaxPodsConstraint = rawDesired.MaxPodsConstraint
	} else {
		rawNew.MaxPodsConstraint = canonicalizeNewNodePoolMaxPodsConstraint(c, rawDesired.MaxPodsConstraint, rawNew.MaxPodsConstraint)
	}

	if dcl.IsEmptyValueIndirect(rawNew.AzureAvailabilityZone) && dcl.IsEmptyValueIndirect(rawDesired.AzureAvailabilityZone) {
		rawNew.AzureAvailabilityZone = rawDesired.AzureAvailabilityZone
	} else {
		if dcl.StringCanonicalize(rawDesired.AzureAvailabilityZone, rawNew.AzureAvailabilityZone) {
			rawNew.AzureAvailabilityZone = rawDesired.AzureAvailabilityZone
		}
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	rawNew.Cluster = rawDesired.Cluster

	return rawNew, nil
}

func canonicalizeNodePoolConfig(des, initial *NodePoolConfig, opts ...dcl.ApplyOption) *NodePoolConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NodePoolConfig{}

	if dcl.StringCanonicalize(des.VmSize, initial.VmSize) || dcl.IsZeroValue(des.VmSize) {
		cDes.VmSize = initial.VmSize
	} else {
		cDes.VmSize = des.VmSize
	}
	cDes.RootVolume = canonicalizeNodePoolConfigRootVolume(des.RootVolume, initial.RootVolume, opts...)
	if dcl.IsZeroValue(des.Tags) || (dcl.IsEmptyValueIndirect(des.Tags) && dcl.IsEmptyValueIndirect(initial.Tags)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Tags = initial.Tags
	} else {
		cDes.Tags = des.Tags
	}
	cDes.SshConfig = canonicalizeNodePoolConfigSshConfig(des.SshConfig, initial.SshConfig, opts...)
	cDes.ProxyConfig = canonicalizeNodePoolConfigProxyConfig(des.ProxyConfig, initial.ProxyConfig, opts...)

	return cDes
}

func canonicalizeNodePoolConfigSlice(des, initial []NodePoolConfig, opts ...dcl.ApplyOption) []NodePoolConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NodePoolConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNodePoolConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NodePoolConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNodePoolConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNodePoolConfig(c *Client, des, nw *NodePoolConfig) *NodePoolConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.VmSize, nw.VmSize) {
		nw.VmSize = des.VmSize
	}
	nw.RootVolume = canonicalizeNewNodePoolConfigRootVolume(c, des.RootVolume, nw.RootVolume)
	nw.SshConfig = canonicalizeNewNodePoolConfigSshConfig(c, des.SshConfig, nw.SshConfig)
	nw.ProxyConfig = canonicalizeNewNodePoolConfigProxyConfig(c, des.ProxyConfig, nw.ProxyConfig)

	return nw
}

func canonicalizeNewNodePoolConfigSet(c *Client, des, nw []NodePoolConfig) []NodePoolConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NodePoolConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNodePoolConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNodePoolConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNodePoolConfigSlice(c *Client, des, nw []NodePoolConfig) []NodePoolConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NodePoolConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNodePoolConfig(c, &d, &n))
	}

	return items
}

func canonicalizeNodePoolConfigRootVolume(des, initial *NodePoolConfigRootVolume, opts ...dcl.ApplyOption) *NodePoolConfigRootVolume {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NodePoolConfigRootVolume{}

	if dcl.IsZeroValue(des.SizeGib) || (dcl.IsEmptyValueIndirect(des.SizeGib) && dcl.IsEmptyValueIndirect(initial.SizeGib)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.SizeGib = initial.SizeGib
	} else {
		cDes.SizeGib = des.SizeGib
	}

	return cDes
}

func canonicalizeNodePoolConfigRootVolumeSlice(des, initial []NodePoolConfigRootVolume, opts ...dcl.ApplyOption) []NodePoolConfigRootVolume {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NodePoolConfigRootVolume, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNodePoolConfigRootVolume(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NodePoolConfigRootVolume, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNodePoolConfigRootVolume(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNodePoolConfigRootVolume(c *Client, des, nw *NodePoolConfigRootVolume) *NodePoolConfigRootVolume {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfigRootVolume while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewNodePoolConfigRootVolumeSet(c *Client, des, nw []NodePoolConfigRootVolume) []NodePoolConfigRootVolume {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NodePoolConfigRootVolume
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNodePoolConfigRootVolumeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNodePoolConfigRootVolume(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNodePoolConfigRootVolumeSlice(c *Client, des, nw []NodePoolConfigRootVolume) []NodePoolConfigRootVolume {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NodePoolConfigRootVolume
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNodePoolConfigRootVolume(c, &d, &n))
	}

	return items
}

func canonicalizeNodePoolConfigSshConfig(des, initial *NodePoolConfigSshConfig, opts ...dcl.ApplyOption) *NodePoolConfigSshConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NodePoolConfigSshConfig{}

	if dcl.StringCanonicalize(des.AuthorizedKey, initial.AuthorizedKey) || dcl.IsZeroValue(des.AuthorizedKey) {
		cDes.AuthorizedKey = initial.AuthorizedKey
	} else {
		cDes.AuthorizedKey = des.AuthorizedKey
	}

	return cDes
}

func canonicalizeNodePoolConfigSshConfigSlice(des, initial []NodePoolConfigSshConfig, opts ...dcl.ApplyOption) []NodePoolConfigSshConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NodePoolConfigSshConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNodePoolConfigSshConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NodePoolConfigSshConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNodePoolConfigSshConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNodePoolConfigSshConfig(c *Client, des, nw *NodePoolConfigSshConfig) *NodePoolConfigSshConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfigSshConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AuthorizedKey, nw.AuthorizedKey) {
		nw.AuthorizedKey = des.AuthorizedKey
	}

	return nw
}

func canonicalizeNewNodePoolConfigSshConfigSet(c *Client, des, nw []NodePoolConfigSshConfig) []NodePoolConfigSshConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NodePoolConfigSshConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNodePoolConfigSshConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNodePoolConfigSshConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNodePoolConfigSshConfigSlice(c *Client, des, nw []NodePoolConfigSshConfig) []NodePoolConfigSshConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NodePoolConfigSshConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNodePoolConfigSshConfig(c, &d, &n))
	}

	return items
}

func canonicalizeNodePoolConfigProxyConfig(des, initial *NodePoolConfigProxyConfig, opts ...dcl.ApplyOption) *NodePoolConfigProxyConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NodePoolConfigProxyConfig{}

	if dcl.StringCanonicalize(des.ResourceGroupId, initial.ResourceGroupId) || dcl.IsZeroValue(des.ResourceGroupId) {
		cDes.ResourceGroupId = initial.ResourceGroupId
	} else {
		cDes.ResourceGroupId = des.ResourceGroupId
	}
	if dcl.StringCanonicalize(des.SecretId, initial.SecretId) || dcl.IsZeroValue(des.SecretId) {
		cDes.SecretId = initial.SecretId
	} else {
		cDes.SecretId = des.SecretId
	}

	return cDes
}

func canonicalizeNodePoolConfigProxyConfigSlice(des, initial []NodePoolConfigProxyConfig, opts ...dcl.ApplyOption) []NodePoolConfigProxyConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NodePoolConfigProxyConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNodePoolConfigProxyConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NodePoolConfigProxyConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNodePoolConfigProxyConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNodePoolConfigProxyConfig(c *Client, des, nw *NodePoolConfigProxyConfig) *NodePoolConfigProxyConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfigProxyConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ResourceGroupId, nw.ResourceGroupId) {
		nw.ResourceGroupId = des.ResourceGroupId
	}
	if dcl.StringCanonicalize(des.SecretId, nw.SecretId) {
		nw.SecretId = des.SecretId
	}

	return nw
}

func canonicalizeNewNodePoolConfigProxyConfigSet(c *Client, des, nw []NodePoolConfigProxyConfig) []NodePoolConfigProxyConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NodePoolConfigProxyConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNodePoolConfigProxyConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNodePoolConfigProxyConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNodePoolConfigProxyConfigSlice(c *Client, des, nw []NodePoolConfigProxyConfig) []NodePoolConfigProxyConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NodePoolConfigProxyConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNodePoolConfigProxyConfig(c, &d, &n))
	}

	return items
}

func canonicalizeNodePoolAutoscaling(des, initial *NodePoolAutoscaling, opts ...dcl.ApplyOption) *NodePoolAutoscaling {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NodePoolAutoscaling{}

	if dcl.IsZeroValue(des.MinNodeCount) || (dcl.IsEmptyValueIndirect(des.MinNodeCount) && dcl.IsEmptyValueIndirect(initial.MinNodeCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MinNodeCount = initial.MinNodeCount
	} else {
		cDes.MinNodeCount = des.MinNodeCount
	}
	if dcl.IsZeroValue(des.MaxNodeCount) || (dcl.IsEmptyValueIndirect(des.MaxNodeCount) && dcl.IsEmptyValueIndirect(initial.MaxNodeCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MaxNodeCount = initial.MaxNodeCount
	} else {
		cDes.MaxNodeCount = des.MaxNodeCount
	}

	return cDes
}

func canonicalizeNodePoolAutoscalingSlice(des, initial []NodePoolAutoscaling, opts ...dcl.ApplyOption) []NodePoolAutoscaling {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NodePoolAutoscaling, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNodePoolAutoscaling(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NodePoolAutoscaling, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNodePoolAutoscaling(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNodePoolAutoscaling(c *Client, des, nw *NodePoolAutoscaling) *NodePoolAutoscaling {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolAutoscaling while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewNodePoolAutoscalingSet(c *Client, des, nw []NodePoolAutoscaling) []NodePoolAutoscaling {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NodePoolAutoscaling
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNodePoolAutoscalingNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNodePoolAutoscaling(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNodePoolAutoscalingSlice(c *Client, des, nw []NodePoolAutoscaling) []NodePoolAutoscaling {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NodePoolAutoscaling
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNodePoolAutoscaling(c, &d, &n))
	}

	return items
}

func canonicalizeNodePoolMaxPodsConstraint(des, initial *NodePoolMaxPodsConstraint, opts ...dcl.ApplyOption) *NodePoolMaxPodsConstraint {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NodePoolMaxPodsConstraint{}

	if dcl.IsZeroValue(des.MaxPodsPerNode) || (dcl.IsEmptyValueIndirect(des.MaxPodsPerNode) && dcl.IsEmptyValueIndirect(initial.MaxPodsPerNode)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MaxPodsPerNode = initial.MaxPodsPerNode
	} else {
		cDes.MaxPodsPerNode = des.MaxPodsPerNode
	}

	return cDes
}

func canonicalizeNodePoolMaxPodsConstraintSlice(des, initial []NodePoolMaxPodsConstraint, opts ...dcl.ApplyOption) []NodePoolMaxPodsConstraint {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NodePoolMaxPodsConstraint, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNodePoolMaxPodsConstraint(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NodePoolMaxPodsConstraint, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNodePoolMaxPodsConstraint(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNodePoolMaxPodsConstraint(c *Client, des, nw *NodePoolMaxPodsConstraint) *NodePoolMaxPodsConstraint {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolMaxPodsConstraint while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewNodePoolMaxPodsConstraintSet(c *Client, des, nw []NodePoolMaxPodsConstraint) []NodePoolMaxPodsConstraint {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []NodePoolMaxPodsConstraint
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareNodePoolMaxPodsConstraintNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewNodePoolMaxPodsConstraint(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewNodePoolMaxPodsConstraintSlice(c *Client, des, nw []NodePoolMaxPodsConstraint) []NodePoolMaxPodsConstraint {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NodePoolMaxPodsConstraint
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNodePoolMaxPodsConstraint(c, &d, &n))
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
func diffNodePool(c *Client, desired, actual *NodePool, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Version, actual.Version, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNodePoolUpdateAzureNodePoolOperation")}, fn.AddNest("Version")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.DiffInfo{ObjectFunction: compareNodePoolConfigNewStyle, EmptyObject: EmptyNodePoolConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubnetId, actual.SubnetId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubnetId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Autoscaling, actual.Autoscaling, dcl.DiffInfo{ObjectFunction: compareNodePoolAutoscalingNewStyle, EmptyObject: EmptyNodePoolAutoscaling, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Autoscaling")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Uid, actual.Uid, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uid")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Reconciling, actual.Reconciling, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Reconciling")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Etag, actual.Etag, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Etag")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Annotations, actual.Annotations, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNodePoolUpdateAzureNodePoolOperation")}, fn.AddNest("Annotations")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxPodsConstraint, actual.MaxPodsConstraint, dcl.DiffInfo{ObjectFunction: compareNodePoolMaxPodsConstraintNewStyle, EmptyObject: EmptyNodePoolMaxPodsConstraint, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxPodsConstraint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AzureAvailabilityZone, actual.AzureAvailabilityZone, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AzureAvailabilityZone")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Cluster, actual.Cluster, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Cluster")); len(ds) != 0 || err != nil {
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
func compareNodePoolConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NodePoolConfig)
	if !ok {
		desiredNotPointer, ok := d.(NodePoolConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfig or *NodePoolConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NodePoolConfig)
	if !ok {
		actualNotPointer, ok := a.(NodePoolConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.VmSize, actual.VmSize, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("VmSize")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RootVolume, actual.RootVolume, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareNodePoolConfigRootVolumeNewStyle, EmptyObject: EmptyNodePoolConfigRootVolume, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RootVolume")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tags, actual.Tags, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SshConfig, actual.SshConfig, dcl.DiffInfo{ObjectFunction: compareNodePoolConfigSshConfigNewStyle, EmptyObject: EmptyNodePoolConfigSshConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SshConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ProxyConfig, actual.ProxyConfig, dcl.DiffInfo{ObjectFunction: compareNodePoolConfigProxyConfigNewStyle, EmptyObject: EmptyNodePoolConfigProxyConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ProxyConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNodePoolConfigRootVolumeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NodePoolConfigRootVolume)
	if !ok {
		desiredNotPointer, ok := d.(NodePoolConfigRootVolume)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigRootVolume or *NodePoolConfigRootVolume", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NodePoolConfigRootVolume)
	if !ok {
		actualNotPointer, ok := a.(NodePoolConfigRootVolume)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigRootVolume", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SizeGib, actual.SizeGib, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SizeGib")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNodePoolConfigSshConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NodePoolConfigSshConfig)
	if !ok {
		desiredNotPointer, ok := d.(NodePoolConfigSshConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigSshConfig or *NodePoolConfigSshConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NodePoolConfigSshConfig)
	if !ok {
		actualNotPointer, ok := a.(NodePoolConfigSshConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigSshConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AuthorizedKey, actual.AuthorizedKey, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNodePoolUpdateAzureNodePoolOperation")}, fn.AddNest("AuthorizedKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNodePoolConfigProxyConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NodePoolConfigProxyConfig)
	if !ok {
		desiredNotPointer, ok := d.(NodePoolConfigProxyConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigProxyConfig or *NodePoolConfigProxyConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NodePoolConfigProxyConfig)
	if !ok {
		actualNotPointer, ok := a.(NodePoolConfigProxyConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigProxyConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ResourceGroupId, actual.ResourceGroupId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ResourceGroupId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecretId, actual.SecretId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SecretId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNodePoolAutoscalingNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NodePoolAutoscaling)
	if !ok {
		desiredNotPointer, ok := d.(NodePoolAutoscaling)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolAutoscaling or *NodePoolAutoscaling", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NodePoolAutoscaling)
	if !ok {
		actualNotPointer, ok := a.(NodePoolAutoscaling)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolAutoscaling", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MinNodeCount, actual.MinNodeCount, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNodePoolUpdateAzureNodePoolOperation")}, fn.AddNest("MinNodeCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxNodeCount, actual.MaxNodeCount, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateNodePoolUpdateAzureNodePoolOperation")}, fn.AddNest("MaxNodeCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNodePoolMaxPodsConstraintNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NodePoolMaxPodsConstraint)
	if !ok {
		desiredNotPointer, ok := d.(NodePoolMaxPodsConstraint)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolMaxPodsConstraint or *NodePoolMaxPodsConstraint", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NodePoolMaxPodsConstraint)
	if !ok {
		actualNotPointer, ok := a.(NodePoolMaxPodsConstraint)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolMaxPodsConstraint", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MaxPodsPerNode, actual.MaxPodsPerNode, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxPodsPerNode")); len(ds) != 0 || err != nil {
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
func (r *NodePool) urlNormalized() *NodePool {
	normalized := dcl.Copy(*r).(NodePool)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Version = dcl.SelfLinkToName(r.Version)
	normalized.SubnetId = dcl.SelfLinkToName(r.SubnetId)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.Etag = dcl.SelfLinkToName(r.Etag)
	normalized.AzureAvailabilityZone = dcl.SelfLinkToName(r.AzureAvailabilityZone)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.Cluster = dcl.SelfLinkToName(r.Cluster)
	return &normalized
}

func (r *NodePool) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateAzureNodePool" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/azureClusters/{{cluster}}/azureNodePools/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the NodePool resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *NodePool) marshal(c *Client) ([]byte, error) {
	m, err := expandNodePool(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling NodePool: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalNodePool decodes JSON responses into the NodePool resource schema.
func unmarshalNodePool(b []byte, c *Client, res *NodePool) (*NodePool, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapNodePool(m, c, res)
}

func unmarshalMapNodePool(m map[string]interface{}, c *Client, res *NodePool) (*NodePool, error) {

	flattened := flattenNodePool(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandNodePool expands NodePool into a JSON request object.
func expandNodePool(c *Client, f *NodePool) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/azureClusters/%s/azureNodePools/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Cluster), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Version; dcl.ValueShouldBeSent(v) {
		m["version"] = v
	}
	if v, err := expandNodePoolConfig(c, f.Config, res); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["config"] = v
	}
	if v := f.SubnetId; dcl.ValueShouldBeSent(v) {
		m["subnetId"] = v
	}
	if v, err := expandNodePoolAutoscaling(c, f.Autoscaling, res); err != nil {
		return nil, fmt.Errorf("error expanding Autoscaling into autoscaling: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["autoscaling"] = v
	}
	if v := f.Annotations; dcl.ValueShouldBeSent(v) {
		m["annotations"] = v
	}
	if v, err := expandNodePoolMaxPodsConstraint(c, f.MaxPodsConstraint, res); err != nil {
		return nil, fmt.Errorf("error expanding MaxPodsConstraint into maxPodsConstraint: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["maxPodsConstraint"] = v
	}
	if v := f.AzureAvailabilityZone; dcl.ValueShouldBeSent(v) {
		m["azureAvailabilityZone"] = v
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
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Cluster into cluster: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["cluster"] = v
	}

	return m, nil
}

// flattenNodePool flattens NodePool from a JSON request object into the
// NodePool type.
func flattenNodePool(c *Client, i interface{}, res *NodePool) *NodePool {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &NodePool{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Version = dcl.FlattenString(m["version"])
	resultRes.Config = flattenNodePoolConfig(c, m["config"], res)
	resultRes.SubnetId = dcl.FlattenString(m["subnetId"])
	resultRes.Autoscaling = flattenNodePoolAutoscaling(c, m["autoscaling"], res)
	resultRes.State = flattenNodePoolStateEnum(m["state"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.Reconciling = dcl.FlattenBool(m["reconciling"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Etag = dcl.FlattenString(m["etag"])
	resultRes.Annotations = dcl.FlattenKeyValuePairs(m["annotations"])
	resultRes.MaxPodsConstraint = flattenNodePoolMaxPodsConstraint(c, m["maxPodsConstraint"], res)
	resultRes.AzureAvailabilityZone = dcl.FlattenString(m["azureAvailabilityZone"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.Cluster = dcl.FlattenString(m["cluster"])

	return resultRes
}

// expandNodePoolConfigMap expands the contents of NodePoolConfig into a JSON
// request object.
func expandNodePoolConfigMap(c *Client, f map[string]NodePoolConfig, res *NodePool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNodePoolConfigSlice expands the contents of NodePoolConfig into a JSON
// request object.
func expandNodePoolConfigSlice(c *Client, f []NodePoolConfig, res *NodePool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigMap flattens the contents of NodePoolConfig from a JSON
// response object.
func flattenNodePoolConfigMap(c *Client, i interface{}, res *NodePool) map[string]NodePoolConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfig{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfig{}
	}

	items := make(map[string]NodePoolConfig)
	for k, item := range a {
		items[k] = *flattenNodePoolConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNodePoolConfigSlice flattens the contents of NodePoolConfig from a JSON
// response object.
func flattenNodePoolConfigSlice(c *Client, i interface{}, res *NodePool) []NodePoolConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfig{}
	}

	if len(a) == 0 {
		return []NodePoolConfig{}
	}

	items := make([]NodePoolConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNodePoolConfig expands an instance of NodePoolConfig into a JSON
// request object.
func expandNodePoolConfig(c *Client, f *NodePoolConfig, res *NodePool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.VmSize; !dcl.IsEmptyValueIndirect(v) {
		m["vmSize"] = v
	}
	if v, err := expandNodePoolConfigRootVolume(c, f.RootVolume, res); err != nil {
		return nil, fmt.Errorf("error expanding RootVolume into rootVolume: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["rootVolume"] = v
	}
	if v := f.Tags; !dcl.IsEmptyValueIndirect(v) {
		m["tags"] = v
	}
	if v, err := expandNodePoolConfigSshConfig(c, f.SshConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SshConfig into sshConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sshConfig"] = v
	}
	if v, err := expandNodePoolConfigProxyConfig(c, f.ProxyConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ProxyConfig into proxyConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["proxyConfig"] = v
	}

	return m, nil
}

// flattenNodePoolConfig flattens an instance of NodePoolConfig from a JSON
// response object.
func flattenNodePoolConfig(c *Client, i interface{}, res *NodePool) *NodePoolConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfig
	}
	r.VmSize = dcl.FlattenString(m["vmSize"])
	r.RootVolume = flattenNodePoolConfigRootVolume(c, m["rootVolume"], res)
	r.Tags = dcl.FlattenKeyValuePairs(m["tags"])
	r.SshConfig = flattenNodePoolConfigSshConfig(c, m["sshConfig"], res)
	r.ProxyConfig = flattenNodePoolConfigProxyConfig(c, m["proxyConfig"], res)

	return r
}

// expandNodePoolConfigRootVolumeMap expands the contents of NodePoolConfigRootVolume into a JSON
// request object.
func expandNodePoolConfigRootVolumeMap(c *Client, f map[string]NodePoolConfigRootVolume, res *NodePool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfigRootVolume(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNodePoolConfigRootVolumeSlice expands the contents of NodePoolConfigRootVolume into a JSON
// request object.
func expandNodePoolConfigRootVolumeSlice(c *Client, f []NodePoolConfigRootVolume, res *NodePool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfigRootVolume(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigRootVolumeMap flattens the contents of NodePoolConfigRootVolume from a JSON
// response object.
func flattenNodePoolConfigRootVolumeMap(c *Client, i interface{}, res *NodePool) map[string]NodePoolConfigRootVolume {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigRootVolume{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigRootVolume{}
	}

	items := make(map[string]NodePoolConfigRootVolume)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigRootVolume(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNodePoolConfigRootVolumeSlice flattens the contents of NodePoolConfigRootVolume from a JSON
// response object.
func flattenNodePoolConfigRootVolumeSlice(c *Client, i interface{}, res *NodePool) []NodePoolConfigRootVolume {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigRootVolume{}
	}

	if len(a) == 0 {
		return []NodePoolConfigRootVolume{}
	}

	items := make([]NodePoolConfigRootVolume, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigRootVolume(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNodePoolConfigRootVolume expands an instance of NodePoolConfigRootVolume into a JSON
// request object.
func expandNodePoolConfigRootVolume(c *Client, f *NodePoolConfigRootVolume, res *NodePool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.SizeGib; !dcl.IsEmptyValueIndirect(v) {
		m["sizeGib"] = v
	}

	return m, nil
}

// flattenNodePoolConfigRootVolume flattens an instance of NodePoolConfigRootVolume from a JSON
// response object.
func flattenNodePoolConfigRootVolume(c *Client, i interface{}, res *NodePool) *NodePoolConfigRootVolume {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfigRootVolume{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfigRootVolume
	}
	r.SizeGib = dcl.FlattenInteger(m["sizeGib"])

	return r
}

// expandNodePoolConfigSshConfigMap expands the contents of NodePoolConfigSshConfig into a JSON
// request object.
func expandNodePoolConfigSshConfigMap(c *Client, f map[string]NodePoolConfigSshConfig, res *NodePool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfigSshConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNodePoolConfigSshConfigSlice expands the contents of NodePoolConfigSshConfig into a JSON
// request object.
func expandNodePoolConfigSshConfigSlice(c *Client, f []NodePoolConfigSshConfig, res *NodePool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfigSshConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigSshConfigMap flattens the contents of NodePoolConfigSshConfig from a JSON
// response object.
func flattenNodePoolConfigSshConfigMap(c *Client, i interface{}, res *NodePool) map[string]NodePoolConfigSshConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigSshConfig{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigSshConfig{}
	}

	items := make(map[string]NodePoolConfigSshConfig)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigSshConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNodePoolConfigSshConfigSlice flattens the contents of NodePoolConfigSshConfig from a JSON
// response object.
func flattenNodePoolConfigSshConfigSlice(c *Client, i interface{}, res *NodePool) []NodePoolConfigSshConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigSshConfig{}
	}

	if len(a) == 0 {
		return []NodePoolConfigSshConfig{}
	}

	items := make([]NodePoolConfigSshConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigSshConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNodePoolConfigSshConfig expands an instance of NodePoolConfigSshConfig into a JSON
// request object.
func expandNodePoolConfigSshConfig(c *Client, f *NodePoolConfigSshConfig, res *NodePool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AuthorizedKey; !dcl.IsEmptyValueIndirect(v) {
		m["authorizedKey"] = v
	}

	return m, nil
}

// flattenNodePoolConfigSshConfig flattens an instance of NodePoolConfigSshConfig from a JSON
// response object.
func flattenNodePoolConfigSshConfig(c *Client, i interface{}, res *NodePool) *NodePoolConfigSshConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfigSshConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfigSshConfig
	}
	r.AuthorizedKey = dcl.FlattenString(m["authorizedKey"])

	return r
}

// expandNodePoolConfigProxyConfigMap expands the contents of NodePoolConfigProxyConfig into a JSON
// request object.
func expandNodePoolConfigProxyConfigMap(c *Client, f map[string]NodePoolConfigProxyConfig, res *NodePool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfigProxyConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNodePoolConfigProxyConfigSlice expands the contents of NodePoolConfigProxyConfig into a JSON
// request object.
func expandNodePoolConfigProxyConfigSlice(c *Client, f []NodePoolConfigProxyConfig, res *NodePool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfigProxyConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigProxyConfigMap flattens the contents of NodePoolConfigProxyConfig from a JSON
// response object.
func flattenNodePoolConfigProxyConfigMap(c *Client, i interface{}, res *NodePool) map[string]NodePoolConfigProxyConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigProxyConfig{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigProxyConfig{}
	}

	items := make(map[string]NodePoolConfigProxyConfig)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigProxyConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNodePoolConfigProxyConfigSlice flattens the contents of NodePoolConfigProxyConfig from a JSON
// response object.
func flattenNodePoolConfigProxyConfigSlice(c *Client, i interface{}, res *NodePool) []NodePoolConfigProxyConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigProxyConfig{}
	}

	if len(a) == 0 {
		return []NodePoolConfigProxyConfig{}
	}

	items := make([]NodePoolConfigProxyConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigProxyConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNodePoolConfigProxyConfig expands an instance of NodePoolConfigProxyConfig into a JSON
// request object.
func expandNodePoolConfigProxyConfig(c *Client, f *NodePoolConfigProxyConfig, res *NodePool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ResourceGroupId; !dcl.IsEmptyValueIndirect(v) {
		m["resourceGroupId"] = v
	}
	if v := f.SecretId; !dcl.IsEmptyValueIndirect(v) {
		m["secretId"] = v
	}

	return m, nil
}

// flattenNodePoolConfigProxyConfig flattens an instance of NodePoolConfigProxyConfig from a JSON
// response object.
func flattenNodePoolConfigProxyConfig(c *Client, i interface{}, res *NodePool) *NodePoolConfigProxyConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfigProxyConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfigProxyConfig
	}
	r.ResourceGroupId = dcl.FlattenString(m["resourceGroupId"])
	r.SecretId = dcl.FlattenString(m["secretId"])

	return r
}

// expandNodePoolAutoscalingMap expands the contents of NodePoolAutoscaling into a JSON
// request object.
func expandNodePoolAutoscalingMap(c *Client, f map[string]NodePoolAutoscaling, res *NodePool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolAutoscaling(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNodePoolAutoscalingSlice expands the contents of NodePoolAutoscaling into a JSON
// request object.
func expandNodePoolAutoscalingSlice(c *Client, f []NodePoolAutoscaling, res *NodePool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolAutoscaling(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolAutoscalingMap flattens the contents of NodePoolAutoscaling from a JSON
// response object.
func flattenNodePoolAutoscalingMap(c *Client, i interface{}, res *NodePool) map[string]NodePoolAutoscaling {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolAutoscaling{}
	}

	if len(a) == 0 {
		return map[string]NodePoolAutoscaling{}
	}

	items := make(map[string]NodePoolAutoscaling)
	for k, item := range a {
		items[k] = *flattenNodePoolAutoscaling(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNodePoolAutoscalingSlice flattens the contents of NodePoolAutoscaling from a JSON
// response object.
func flattenNodePoolAutoscalingSlice(c *Client, i interface{}, res *NodePool) []NodePoolAutoscaling {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolAutoscaling{}
	}

	if len(a) == 0 {
		return []NodePoolAutoscaling{}
	}

	items := make([]NodePoolAutoscaling, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolAutoscaling(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNodePoolAutoscaling expands an instance of NodePoolAutoscaling into a JSON
// request object.
func expandNodePoolAutoscaling(c *Client, f *NodePoolAutoscaling, res *NodePool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MinNodeCount; !dcl.IsEmptyValueIndirect(v) {
		m["minNodeCount"] = v
	}
	if v := f.MaxNodeCount; !dcl.IsEmptyValueIndirect(v) {
		m["maxNodeCount"] = v
	}

	return m, nil
}

// flattenNodePoolAutoscaling flattens an instance of NodePoolAutoscaling from a JSON
// response object.
func flattenNodePoolAutoscaling(c *Client, i interface{}, res *NodePool) *NodePoolAutoscaling {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolAutoscaling{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolAutoscaling
	}
	r.MinNodeCount = dcl.FlattenInteger(m["minNodeCount"])
	r.MaxNodeCount = dcl.FlattenInteger(m["maxNodeCount"])

	return r
}

// expandNodePoolMaxPodsConstraintMap expands the contents of NodePoolMaxPodsConstraint into a JSON
// request object.
func expandNodePoolMaxPodsConstraintMap(c *Client, f map[string]NodePoolMaxPodsConstraint, res *NodePool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolMaxPodsConstraint(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNodePoolMaxPodsConstraintSlice expands the contents of NodePoolMaxPodsConstraint into a JSON
// request object.
func expandNodePoolMaxPodsConstraintSlice(c *Client, f []NodePoolMaxPodsConstraint, res *NodePool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolMaxPodsConstraint(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolMaxPodsConstraintMap flattens the contents of NodePoolMaxPodsConstraint from a JSON
// response object.
func flattenNodePoolMaxPodsConstraintMap(c *Client, i interface{}, res *NodePool) map[string]NodePoolMaxPodsConstraint {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolMaxPodsConstraint{}
	}

	if len(a) == 0 {
		return map[string]NodePoolMaxPodsConstraint{}
	}

	items := make(map[string]NodePoolMaxPodsConstraint)
	for k, item := range a {
		items[k] = *flattenNodePoolMaxPodsConstraint(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenNodePoolMaxPodsConstraintSlice flattens the contents of NodePoolMaxPodsConstraint from a JSON
// response object.
func flattenNodePoolMaxPodsConstraintSlice(c *Client, i interface{}, res *NodePool) []NodePoolMaxPodsConstraint {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolMaxPodsConstraint{}
	}

	if len(a) == 0 {
		return []NodePoolMaxPodsConstraint{}
	}

	items := make([]NodePoolMaxPodsConstraint, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolMaxPodsConstraint(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandNodePoolMaxPodsConstraint expands an instance of NodePoolMaxPodsConstraint into a JSON
// request object.
func expandNodePoolMaxPodsConstraint(c *Client, f *NodePoolMaxPodsConstraint, res *NodePool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MaxPodsPerNode; !dcl.IsEmptyValueIndirect(v) {
		m["maxPodsPerNode"] = v
	}

	return m, nil
}

// flattenNodePoolMaxPodsConstraint flattens an instance of NodePoolMaxPodsConstraint from a JSON
// response object.
func flattenNodePoolMaxPodsConstraint(c *Client, i interface{}, res *NodePool) *NodePoolMaxPodsConstraint {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolMaxPodsConstraint{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolMaxPodsConstraint
	}
	r.MaxPodsPerNode = dcl.FlattenInteger(m["maxPodsPerNode"])

	return r
}

// flattenNodePoolStateEnumMap flattens the contents of NodePoolStateEnum from a JSON
// response object.
func flattenNodePoolStateEnumMap(c *Client, i interface{}, res *NodePool) map[string]NodePoolStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolStateEnum{}
	}

	if len(a) == 0 {
		return map[string]NodePoolStateEnum{}
	}

	items := make(map[string]NodePoolStateEnum)
	for k, item := range a {
		items[k] = *flattenNodePoolStateEnum(item.(interface{}))
	}

	return items
}

// flattenNodePoolStateEnumSlice flattens the contents of NodePoolStateEnum from a JSON
// response object.
func flattenNodePoolStateEnumSlice(c *Client, i interface{}, res *NodePool) []NodePoolStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolStateEnum{}
	}

	if len(a) == 0 {
		return []NodePoolStateEnum{}
	}

	items := make([]NodePoolStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolStateEnum(item.(interface{})))
	}

	return items
}

// flattenNodePoolStateEnum asserts that an interface is a string, and returns a
// pointer to a *NodePoolStateEnum with the same value as that string.
func flattenNodePoolStateEnum(i interface{}) *NodePoolStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return NodePoolStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *NodePool) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalNodePool(b, c, r)
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
		if nr.Cluster == nil && ncr.Cluster == nil {
			c.Config.Logger.Info("Both Cluster fields null - considering equal.")
		} else if nr.Cluster == nil || ncr.Cluster == nil {
			c.Config.Logger.Info("Only one Cluster field is null - considering unequal.")
			return false
		} else if *nr.Cluster != *ncr.Cluster {
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

type nodePoolDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         nodePoolApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToNodePoolDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]nodePoolDiff, error) {
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
	var diffs []nodePoolDiff
	// For each operation name, create a nodePoolDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := nodePoolDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToNodePoolApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToNodePoolApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (nodePoolApiOperation, error) {
	switch opName {

	case "updateNodePoolUpdateAzureNodePoolOperation":
		return &updateNodePoolUpdateAzureNodePoolOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractNodePoolFields(r *NodePool) error {
	vConfig := r.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &NodePoolConfig{}
	}
	if err := extractNodePoolConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		r.Config = vConfig
	}
	vAutoscaling := r.Autoscaling
	if vAutoscaling == nil {
		// note: explicitly not the empty object.
		vAutoscaling = &NodePoolAutoscaling{}
	}
	if err := extractNodePoolAutoscalingFields(r, vAutoscaling); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAutoscaling) {
		r.Autoscaling = vAutoscaling
	}
	vMaxPodsConstraint := r.MaxPodsConstraint
	if vMaxPodsConstraint == nil {
		// note: explicitly not the empty object.
		vMaxPodsConstraint = &NodePoolMaxPodsConstraint{}
	}
	if err := extractNodePoolMaxPodsConstraintFields(r, vMaxPodsConstraint); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMaxPodsConstraint) {
		r.MaxPodsConstraint = vMaxPodsConstraint
	}
	return nil
}
func extractNodePoolConfigFields(r *NodePool, o *NodePoolConfig) error {
	vRootVolume := o.RootVolume
	if vRootVolume == nil {
		// note: explicitly not the empty object.
		vRootVolume = &NodePoolConfigRootVolume{}
	}
	if err := extractNodePoolConfigRootVolumeFields(r, vRootVolume); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRootVolume) {
		o.RootVolume = vRootVolume
	}
	vSshConfig := o.SshConfig
	if vSshConfig == nil {
		// note: explicitly not the empty object.
		vSshConfig = &NodePoolConfigSshConfig{}
	}
	if err := extractNodePoolConfigSshConfigFields(r, vSshConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSshConfig) {
		o.SshConfig = vSshConfig
	}
	vProxyConfig := o.ProxyConfig
	if vProxyConfig == nil {
		// note: explicitly not the empty object.
		vProxyConfig = &NodePoolConfigProxyConfig{}
	}
	if err := extractNodePoolConfigProxyConfigFields(r, vProxyConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vProxyConfig) {
		o.ProxyConfig = vProxyConfig
	}
	return nil
}
func extractNodePoolConfigRootVolumeFields(r *NodePool, o *NodePoolConfigRootVolume) error {
	return nil
}
func extractNodePoolConfigSshConfigFields(r *NodePool, o *NodePoolConfigSshConfig) error {
	return nil
}
func extractNodePoolConfigProxyConfigFields(r *NodePool, o *NodePoolConfigProxyConfig) error {
	return nil
}
func extractNodePoolAutoscalingFields(r *NodePool, o *NodePoolAutoscaling) error {
	return nil
}
func extractNodePoolMaxPodsConstraintFields(r *NodePool, o *NodePoolMaxPodsConstraint) error {
	return nil
}

func postReadExtractNodePoolFields(r *NodePool) error {
	vConfig := r.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &NodePoolConfig{}
	}
	if err := postReadExtractNodePoolConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		r.Config = vConfig
	}
	vAutoscaling := r.Autoscaling
	if vAutoscaling == nil {
		// note: explicitly not the empty object.
		vAutoscaling = &NodePoolAutoscaling{}
	}
	if err := postReadExtractNodePoolAutoscalingFields(r, vAutoscaling); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAutoscaling) {
		r.Autoscaling = vAutoscaling
	}
	vMaxPodsConstraint := r.MaxPodsConstraint
	if vMaxPodsConstraint == nil {
		// note: explicitly not the empty object.
		vMaxPodsConstraint = &NodePoolMaxPodsConstraint{}
	}
	if err := postReadExtractNodePoolMaxPodsConstraintFields(r, vMaxPodsConstraint); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMaxPodsConstraint) {
		r.MaxPodsConstraint = vMaxPodsConstraint
	}
	return nil
}
func postReadExtractNodePoolConfigFields(r *NodePool, o *NodePoolConfig) error {
	vRootVolume := o.RootVolume
	if vRootVolume == nil {
		// note: explicitly not the empty object.
		vRootVolume = &NodePoolConfigRootVolume{}
	}
	if err := extractNodePoolConfigRootVolumeFields(r, vRootVolume); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRootVolume) {
		o.RootVolume = vRootVolume
	}
	vSshConfig := o.SshConfig
	if vSshConfig == nil {
		// note: explicitly not the empty object.
		vSshConfig = &NodePoolConfigSshConfig{}
	}
	if err := extractNodePoolConfigSshConfigFields(r, vSshConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSshConfig) {
		o.SshConfig = vSshConfig
	}
	vProxyConfig := o.ProxyConfig
	if vProxyConfig == nil {
		// note: explicitly not the empty object.
		vProxyConfig = &NodePoolConfigProxyConfig{}
	}
	if err := extractNodePoolConfigProxyConfigFields(r, vProxyConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vProxyConfig) {
		o.ProxyConfig = vProxyConfig
	}
	return nil
}
func postReadExtractNodePoolConfigRootVolumeFields(r *NodePool, o *NodePoolConfigRootVolume) error {
	return nil
}
func postReadExtractNodePoolConfigSshConfigFields(r *NodePool, o *NodePoolConfigSshConfig) error {
	return nil
}
func postReadExtractNodePoolConfigProxyConfigFields(r *NodePool, o *NodePoolConfigProxyConfig) error {
	return nil
}
func postReadExtractNodePoolAutoscalingFields(r *NodePool, o *NodePoolAutoscaling) error {
	return nil
}
func postReadExtractNodePoolMaxPodsConstraintFields(r *NodePool, o *NodePoolMaxPodsConstraint) error {
	return nil
}
