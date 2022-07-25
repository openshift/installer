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
package containeraws

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
	if err := dcl.Required(r, "autoscaling"); err != nil {
		return err
	}
	if err := dcl.Required(r, "subnetId"); err != nil {
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
	if err := dcl.Required(r, "iamInstanceProfile"); err != nil {
		return err
	}
	if err := dcl.Required(r, "configEncryption"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.RootVolume) {
		if err := r.RootVolume.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ConfigEncryption) {
		if err := r.ConfigEncryption.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SshConfig) {
		if err := r.SshConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *NodePoolConfigRootVolume) validate() error {
	return nil
}
func (r *NodePoolConfigTaints) validate() error {
	if err := dcl.Required(r, "key"); err != nil {
		return err
	}
	if err := dcl.Required(r, "value"); err != nil {
		return err
	}
	if err := dcl.Required(r, "effect"); err != nil {
		return err
	}
	return nil
}
func (r *NodePoolConfigConfigEncryption) validate() error {
	if err := dcl.Required(r, "kmsKeyArn"); err != nil {
		return err
	}
	return nil
}
func (r *NodePoolConfigSshConfig) validate() error {
	if err := dcl.Required(r, "ec2KeyPair"); err != nil {
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
	return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *NodePool) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools", nr.basePath(), userBasePath, params), nil

}

func (r *NodePool) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools?awsNodePoolId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *NodePool) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools/{{name}}", nr.basePath(), userBasePath, params), nil
}

// nodePoolApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type nodePoolApiOperation interface {
	do(context.Context, *NodePool, *Client) error
}

// newUpdateNodePoolUpdateAwsNodePoolRequest creates a request for an
// NodePool resource's UpdateAwsNodePool update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateNodePoolUpdateAwsNodePoolRequest(ctx context.Context, f *NodePool, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}

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

// marshalUpdateNodePoolUpdateAwsNodePoolRequest converts the update into
// the final JSON request body.
func marshalUpdateNodePoolUpdateAwsNodePoolRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateNodePoolUpdateAwsNodePoolOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateNodePoolUpdateAwsNodePoolOperation) do(ctx context.Context, r *NodePool, c *Client) error {
	_, err := c.GetNodePool(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateAwsNodePool")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateNodePoolUpdateAwsNodePoolRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateNodePoolUpdateAwsNodePoolRequest(c, req)
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
	AwsNodePools []map[string]interface{} `json:"awsNodePools"`
	Token        string                   `json:"nextPageToken"`
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
	for _, v := range m.AwsNodePools {
		res, err := unmarshalMapNodePool(v, c)
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

	// we saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// this is the reason we are adding retry to handle that case.
	maxRetry := 10
	for i := 1; i <= maxRetry; i++ {
		_, err = c.GetNodePool(ctx, r)
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
	canonicalDesired.Autoscaling = canonicalizeNodePoolAutoscaling(rawDesired.Autoscaling, rawInitial.Autoscaling, opts...)
	if dcl.StringCanonicalize(rawDesired.SubnetId, rawInitial.SubnetId) {
		canonicalDesired.SubnetId = rawInitial.SubnetId
	} else {
		canonicalDesired.SubnetId = rawDesired.SubnetId
	}
	if dcl.IsZeroValue(rawDesired.Annotations) {
		canonicalDesired.Annotations = rawInitial.Annotations
	} else {
		canonicalDesired.Annotations = rawDesired.Annotations
	}
	canonicalDesired.MaxPodsConstraint = canonicalizeNodePoolMaxPodsConstraint(rawDesired.MaxPodsConstraint, rawInitial.MaxPodsConstraint, opts...)
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

	if dcl.IsNotReturnedByServer(rawNew.Name) && dcl.IsNotReturnedByServer(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Version) && dcl.IsNotReturnedByServer(rawDesired.Version) {
		rawNew.Version = rawDesired.Version
	} else {
		if dcl.StringCanonicalize(rawDesired.Version, rawNew.Version) {
			rawNew.Version = rawDesired.Version
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Config) && dcl.IsNotReturnedByServer(rawDesired.Config) {
		rawNew.Config = rawDesired.Config
	} else {
		rawNew.Config = canonicalizeNewNodePoolConfig(c, rawDesired.Config, rawNew.Config)
	}

	if dcl.IsNotReturnedByServer(rawNew.Autoscaling) && dcl.IsNotReturnedByServer(rawDesired.Autoscaling) {
		rawNew.Autoscaling = rawDesired.Autoscaling
	} else {
		rawNew.Autoscaling = canonicalizeNewNodePoolAutoscaling(c, rawDesired.Autoscaling, rawNew.Autoscaling)
	}

	if dcl.IsNotReturnedByServer(rawNew.SubnetId) && dcl.IsNotReturnedByServer(rawDesired.SubnetId) {
		rawNew.SubnetId = rawDesired.SubnetId
	} else {
		if dcl.StringCanonicalize(rawDesired.SubnetId, rawNew.SubnetId) {
			rawNew.SubnetId = rawDesired.SubnetId
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.State) && dcl.IsNotReturnedByServer(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Uid) && dcl.IsNotReturnedByServer(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Reconciling) && dcl.IsNotReturnedByServer(rawDesired.Reconciling) {
		rawNew.Reconciling = rawDesired.Reconciling
	} else {
		if dcl.BoolCanonicalize(rawDesired.Reconciling, rawNew.Reconciling) {
			rawNew.Reconciling = rawDesired.Reconciling
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.CreateTime) && dcl.IsNotReturnedByServer(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.UpdateTime) && dcl.IsNotReturnedByServer(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Etag) && dcl.IsNotReturnedByServer(rawDesired.Etag) {
		rawNew.Etag = rawDesired.Etag
	} else {
		if dcl.StringCanonicalize(rawDesired.Etag, rawNew.Etag) {
			rawNew.Etag = rawDesired.Etag
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Annotations) && dcl.IsNotReturnedByServer(rawDesired.Annotations) {
		rawNew.Annotations = rawDesired.Annotations
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.MaxPodsConstraint) && dcl.IsNotReturnedByServer(rawDesired.MaxPodsConstraint) {
		rawNew.MaxPodsConstraint = rawDesired.MaxPodsConstraint
	} else {
		rawNew.MaxPodsConstraint = canonicalizeNewNodePoolMaxPodsConstraint(c, rawDesired.MaxPodsConstraint, rawNew.MaxPodsConstraint)
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

	if dcl.StringCanonicalize(des.InstanceType, initial.InstanceType) || dcl.IsZeroValue(des.InstanceType) {
		cDes.InstanceType = initial.InstanceType
	} else {
		cDes.InstanceType = des.InstanceType
	}
	cDes.RootVolume = canonicalizeNodePoolConfigRootVolume(des.RootVolume, initial.RootVolume, opts...)
	cDes.Taints = canonicalizeNodePoolConfigTaintsSlice(des.Taints, initial.Taints, opts...)
	if dcl.IsZeroValue(des.Labels) {
		cDes.Labels = initial.Labels
	} else {
		cDes.Labels = des.Labels
	}
	if dcl.IsZeroValue(des.Tags) {
		cDes.Tags = initial.Tags
	} else {
		cDes.Tags = des.Tags
	}
	if dcl.StringCanonicalize(des.IamInstanceProfile, initial.IamInstanceProfile) || dcl.IsZeroValue(des.IamInstanceProfile) {
		cDes.IamInstanceProfile = initial.IamInstanceProfile
	} else {
		cDes.IamInstanceProfile = des.IamInstanceProfile
	}
	cDes.ConfigEncryption = canonicalizeNodePoolConfigConfigEncryption(des.ConfigEncryption, initial.ConfigEncryption, opts...)
	cDes.SshConfig = canonicalizeNodePoolConfigSshConfig(des.SshConfig, initial.SshConfig, opts...)
	if dcl.StringArrayCanonicalize(des.SecurityGroupIds, initial.SecurityGroupIds) || dcl.IsZeroValue(des.SecurityGroupIds) {
		cDes.SecurityGroupIds = initial.SecurityGroupIds
	} else {
		cDes.SecurityGroupIds = des.SecurityGroupIds
	}

	return cDes
}

func canonicalizeNodePoolConfigSlice(des, initial []NodePoolConfig, opts ...dcl.ApplyOption) []NodePoolConfig {
	if des == nil {
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
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.InstanceType, nw.InstanceType) {
		nw.InstanceType = des.InstanceType
	}
	nw.RootVolume = canonicalizeNewNodePoolConfigRootVolume(c, des.RootVolume, nw.RootVolume)
	nw.Taints = canonicalizeNewNodePoolConfigTaintsSlice(c, des.Taints, nw.Taints)
	if dcl.StringCanonicalize(des.IamInstanceProfile, nw.IamInstanceProfile) {
		nw.IamInstanceProfile = des.IamInstanceProfile
	}
	nw.ConfigEncryption = canonicalizeNewNodePoolConfigConfigEncryption(c, des.ConfigEncryption, nw.ConfigEncryption)
	nw.SshConfig = canonicalizeNewNodePoolConfigSshConfig(c, des.SshConfig, nw.SshConfig)
	if dcl.StringArrayCanonicalize(des.SecurityGroupIds, nw.SecurityGroupIds) {
		nw.SecurityGroupIds = des.SecurityGroupIds
	}

	return nw
}

func canonicalizeNewNodePoolConfigSet(c *Client, des, nw []NodePoolConfig) []NodePoolConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []NodePoolConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareNodePoolConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

	if dcl.IsZeroValue(des.SizeGib) {
		cDes.SizeGib = initial.SizeGib
	} else {
		cDes.SizeGib = des.SizeGib
	}
	if dcl.IsZeroValue(des.VolumeType) {
		cDes.VolumeType = initial.VolumeType
	} else {
		cDes.VolumeType = des.VolumeType
	}
	if dcl.IsZeroValue(des.Iops) {
		cDes.Iops = initial.Iops
	} else {
		cDes.Iops = des.Iops
	}
	if dcl.StringCanonicalize(des.KmsKeyArn, initial.KmsKeyArn) || dcl.IsZeroValue(des.KmsKeyArn) {
		cDes.KmsKeyArn = initial.KmsKeyArn
	} else {
		cDes.KmsKeyArn = des.KmsKeyArn
	}

	return cDes
}

func canonicalizeNodePoolConfigRootVolumeSlice(des, initial []NodePoolConfigRootVolume, opts ...dcl.ApplyOption) []NodePoolConfigRootVolume {
	if des == nil {
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
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfigRootVolume while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KmsKeyArn, nw.KmsKeyArn) {
		nw.KmsKeyArn = des.KmsKeyArn
	}

	return nw
}

func canonicalizeNewNodePoolConfigRootVolumeSet(c *Client, des, nw []NodePoolConfigRootVolume) []NodePoolConfigRootVolume {
	if des == nil {
		return nw
	}
	var reorderedNew []NodePoolConfigRootVolume
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareNodePoolConfigRootVolumeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNodePoolConfigTaints(des, initial *NodePoolConfigTaints, opts ...dcl.ApplyOption) *NodePoolConfigTaints {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NodePoolConfigTaints{}

	if dcl.StringCanonicalize(des.Key, initial.Key) || dcl.IsZeroValue(des.Key) {
		cDes.Key = initial.Key
	} else {
		cDes.Key = des.Key
	}
	if dcl.StringCanonicalize(des.Value, initial.Value) || dcl.IsZeroValue(des.Value) {
		cDes.Value = initial.Value
	} else {
		cDes.Value = des.Value
	}
	if dcl.IsZeroValue(des.Effect) {
		cDes.Effect = initial.Effect
	} else {
		cDes.Effect = des.Effect
	}

	return cDes
}

func canonicalizeNodePoolConfigTaintsSlice(des, initial []NodePoolConfigTaints, opts ...dcl.ApplyOption) []NodePoolConfigTaints {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NodePoolConfigTaints, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNodePoolConfigTaints(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NodePoolConfigTaints, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNodePoolConfigTaints(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNodePoolConfigTaints(c *Client, des, nw *NodePoolConfigTaints) *NodePoolConfigTaints {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfigTaints while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewNodePoolConfigTaintsSet(c *Client, des, nw []NodePoolConfigTaints) []NodePoolConfigTaints {
	if des == nil {
		return nw
	}
	var reorderedNew []NodePoolConfigTaints
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareNodePoolConfigTaintsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewNodePoolConfigTaintsSlice(c *Client, des, nw []NodePoolConfigTaints) []NodePoolConfigTaints {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NodePoolConfigTaints
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNodePoolConfigTaints(c, &d, &n))
	}

	return items
}

func canonicalizeNodePoolConfigConfigEncryption(des, initial *NodePoolConfigConfigEncryption, opts ...dcl.ApplyOption) *NodePoolConfigConfigEncryption {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &NodePoolConfigConfigEncryption{}

	if dcl.StringCanonicalize(des.KmsKeyArn, initial.KmsKeyArn) || dcl.IsZeroValue(des.KmsKeyArn) {
		cDes.KmsKeyArn = initial.KmsKeyArn
	} else {
		cDes.KmsKeyArn = des.KmsKeyArn
	}

	return cDes
}

func canonicalizeNodePoolConfigConfigEncryptionSlice(des, initial []NodePoolConfigConfigEncryption, opts ...dcl.ApplyOption) []NodePoolConfigConfigEncryption {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]NodePoolConfigConfigEncryption, 0, len(des))
		for _, d := range des {
			cd := canonicalizeNodePoolConfigConfigEncryption(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]NodePoolConfigConfigEncryption, 0, len(des))
	for i, d := range des {
		cd := canonicalizeNodePoolConfigConfigEncryption(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewNodePoolConfigConfigEncryption(c *Client, des, nw *NodePoolConfigConfigEncryption) *NodePoolConfigConfigEncryption {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfigConfigEncryption while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KmsKeyArn, nw.KmsKeyArn) {
		nw.KmsKeyArn = des.KmsKeyArn
	}

	return nw
}

func canonicalizeNewNodePoolConfigConfigEncryptionSet(c *Client, des, nw []NodePoolConfigConfigEncryption) []NodePoolConfigConfigEncryption {
	if des == nil {
		return nw
	}
	var reorderedNew []NodePoolConfigConfigEncryption
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareNodePoolConfigConfigEncryptionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewNodePoolConfigConfigEncryptionSlice(c *Client, des, nw []NodePoolConfigConfigEncryption) []NodePoolConfigConfigEncryption {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []NodePoolConfigConfigEncryption
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewNodePoolConfigConfigEncryption(c, &d, &n))
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

	if dcl.StringCanonicalize(des.Ec2KeyPair, initial.Ec2KeyPair) || dcl.IsZeroValue(des.Ec2KeyPair) {
		cDes.Ec2KeyPair = initial.Ec2KeyPair
	} else {
		cDes.Ec2KeyPair = des.Ec2KeyPair
	}

	return cDes
}

func canonicalizeNodePoolConfigSshConfigSlice(des, initial []NodePoolConfigSshConfig, opts ...dcl.ApplyOption) []NodePoolConfigSshConfig {
	if des == nil {
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
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for NodePoolConfigSshConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Ec2KeyPair, nw.Ec2KeyPair) {
		nw.Ec2KeyPair = des.Ec2KeyPair
	}

	return nw
}

func canonicalizeNewNodePoolConfigSshConfigSet(c *Client, des, nw []NodePoolConfigSshConfig) []NodePoolConfigSshConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []NodePoolConfigSshConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareNodePoolConfigSshConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

	if dcl.IsZeroValue(des.MinNodeCount) {
		cDes.MinNodeCount = initial.MinNodeCount
	} else {
		cDes.MinNodeCount = des.MinNodeCount
	}
	if dcl.IsZeroValue(des.MaxNodeCount) {
		cDes.MaxNodeCount = initial.MaxNodeCount
	} else {
		cDes.MaxNodeCount = des.MaxNodeCount
	}

	return cDes
}

func canonicalizeNodePoolAutoscalingSlice(des, initial []NodePoolAutoscaling, opts ...dcl.ApplyOption) []NodePoolAutoscaling {
	if des == nil {
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
		if dcl.IsNotReturnedByServer(des) {
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
	var reorderedNew []NodePoolAutoscaling
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareNodePoolAutoscalingNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

	if dcl.IsZeroValue(des.MaxPodsPerNode) {
		cDes.MaxPodsPerNode = initial.MaxPodsPerNode
	} else {
		cDes.MaxPodsPerNode = des.MaxPodsPerNode
	}

	return cDes
}

func canonicalizeNodePoolMaxPodsConstraintSlice(des, initial []NodePoolMaxPodsConstraint, opts ...dcl.ApplyOption) []NodePoolMaxPodsConstraint {
	if des == nil {
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
		if dcl.IsNotReturnedByServer(des) {
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
	var reorderedNew []NodePoolMaxPodsConstraint
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareNodePoolMaxPodsConstraintNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Version, actual.Version, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Version")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.Info{ObjectFunction: compareNodePoolConfigNewStyle, EmptyObject: EmptyNodePoolConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Autoscaling, actual.Autoscaling, dcl.Info{ObjectFunction: compareNodePoolAutoscalingNewStyle, EmptyObject: EmptyNodePoolAutoscaling, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Autoscaling")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubnetId, actual.SubnetId, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubnetId")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Uid, actual.Uid, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uid")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Reconciling, actual.Reconciling, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Reconciling")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreateTime, actual.CreateTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UpdateTime, actual.UpdateTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UpdateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Etag, actual.Etag, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Etag")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Annotations, actual.Annotations, dcl.Info{OperationSelector: dcl.TriggersOperation("updateNodePoolUpdateAwsNodePoolOperation")}, fn.AddNest("Annotations")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxPodsConstraint, actual.MaxPodsConstraint, dcl.Info{ObjectFunction: compareNodePoolMaxPodsConstraintNewStyle, EmptyObject: EmptyNodePoolMaxPodsConstraint, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxPodsConstraint")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Cluster, actual.Cluster, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Cluster")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
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

	if ds, err := dcl.Diff(desired.InstanceType, actual.InstanceType, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RootVolume, actual.RootVolume, dcl.Info{ObjectFunction: compareNodePoolConfigRootVolumeNewStyle, EmptyObject: EmptyNodePoolConfigRootVolume, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RootVolume")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Taints, actual.Taints, dcl.Info{ObjectFunction: compareNodePoolConfigTaintsNewStyle, EmptyObject: EmptyNodePoolConfigTaints, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Taints")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tags, actual.Tags, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IamInstanceProfile, actual.IamInstanceProfile, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IamInstanceProfile")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ConfigEncryption, actual.ConfigEncryption, dcl.Info{ObjectFunction: compareNodePoolConfigConfigEncryptionNewStyle, EmptyObject: EmptyNodePoolConfigConfigEncryption, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConfigEncryption")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SshConfig, actual.SshConfig, dcl.Info{ObjectFunction: compareNodePoolConfigSshConfigNewStyle, EmptyObject: EmptyNodePoolConfigSshConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SshConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecurityGroupIds, actual.SecurityGroupIds, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SecurityGroupIds")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.SizeGib, actual.SizeGib, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SizeGib")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.VolumeType, actual.VolumeType, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("VolumeType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Iops, actual.Iops, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Iops")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KmsKeyArn, actual.KmsKeyArn, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KmsKeyArn")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNodePoolConfigTaintsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NodePoolConfigTaints)
	if !ok {
		desiredNotPointer, ok := d.(NodePoolConfigTaints)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigTaints or *NodePoolConfigTaints", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NodePoolConfigTaints)
	if !ok {
		actualNotPointer, ok := a.(NodePoolConfigTaints)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigTaints", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Key, actual.Key, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Key")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Value, actual.Value, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Value")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Effect, actual.Effect, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Effect")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareNodePoolConfigConfigEncryptionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*NodePoolConfigConfigEncryption)
	if !ok {
		desiredNotPointer, ok := d.(NodePoolConfigConfigEncryption)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigConfigEncryption or *NodePoolConfigConfigEncryption", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*NodePoolConfigConfigEncryption)
	if !ok {
		actualNotPointer, ok := a.(NodePoolConfigConfigEncryption)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a NodePoolConfigConfigEncryption", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KmsKeyArn, actual.KmsKeyArn, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KmsKeyArn")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Ec2KeyPair, actual.Ec2KeyPair, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Ec2KeyPair")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.MinNodeCount, actual.MinNodeCount, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinNodeCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxNodeCount, actual.MaxNodeCount, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxNodeCount")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.MaxPodsPerNode, actual.MaxPodsPerNode, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxPodsPerNode")); len(ds) != 0 || err != nil {
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
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.Cluster = dcl.SelfLinkToName(r.Cluster)
	return &normalized
}

func (r *NodePool) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateAwsNodePool" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"cluster":  dcl.ValueOrEmptyString(nr.Cluster),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools/{{name}}", nr.basePath(), userBasePath, fields), nil

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
func unmarshalNodePool(b []byte, c *Client) (*NodePool, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapNodePool(m, c)
}

func unmarshalMapNodePool(m map[string]interface{}, c *Client) (*NodePool, error) {

	flattened := flattenNodePool(c, m)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandNodePool expands NodePool into a JSON request object.
func expandNodePool(c *Client, f *NodePool) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if v, err := dcl.DeriveField("projects/%s/locations/%s/awsClusters/%s/awsNodePools/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Cluster), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if v != nil {
		m["name"] = v
	}
	if v := f.Version; dcl.ValueShouldBeSent(v) {
		m["version"] = v
	}
	if v, err := expandNodePoolConfig(c, f.Config); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if v != nil {
		m["config"] = v
	}
	if v, err := expandNodePoolAutoscaling(c, f.Autoscaling); err != nil {
		return nil, fmt.Errorf("error expanding Autoscaling into autoscaling: %w", err)
	} else if v != nil {
		m["autoscaling"] = v
	}
	if v := f.SubnetId; dcl.ValueShouldBeSent(v) {
		m["subnetId"] = v
	}
	if v := f.Annotations; dcl.ValueShouldBeSent(v) {
		m["annotations"] = v
	}
	if v, err := expandNodePoolMaxPodsConstraint(c, f.MaxPodsConstraint); err != nil {
		return nil, fmt.Errorf("error expanding MaxPodsConstraint into maxPodsConstraint: %w", err)
	} else if v != nil {
		m["maxPodsConstraint"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if v != nil {
		m["project"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if v != nil {
		m["location"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Cluster into cluster: %w", err)
	} else if v != nil {
		m["cluster"] = v
	}

	return m, nil
}

// flattenNodePool flattens NodePool from a JSON request object into the
// NodePool type.
func flattenNodePool(c *Client, i interface{}) *NodePool {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	res := &NodePool{}
	res.Name = dcl.FlattenString(m["name"])
	res.Version = dcl.FlattenString(m["version"])
	res.Config = flattenNodePoolConfig(c, m["config"])
	res.Autoscaling = flattenNodePoolAutoscaling(c, m["autoscaling"])
	res.SubnetId = dcl.FlattenString(m["subnetId"])
	res.State = flattenNodePoolStateEnum(m["state"])
	res.Uid = dcl.FlattenString(m["uid"])
	res.Reconciling = dcl.FlattenBool(m["reconciling"])
	res.CreateTime = dcl.FlattenString(m["createTime"])
	res.UpdateTime = dcl.FlattenString(m["updateTime"])
	res.Etag = dcl.FlattenString(m["etag"])
	res.Annotations = dcl.FlattenKeyValuePairs(m["annotations"])
	res.MaxPodsConstraint = flattenNodePoolMaxPodsConstraint(c, m["maxPodsConstraint"])
	res.Project = dcl.FlattenString(m["project"])
	res.Location = dcl.FlattenString(m["location"])
	res.Cluster = dcl.FlattenString(m["cluster"])

	return res
}

// expandNodePoolConfigMap expands the contents of NodePoolConfig into a JSON
// request object.
func expandNodePoolConfigMap(c *Client, f map[string]NodePoolConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfig(c, &item)
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
func expandNodePoolConfigSlice(c *Client, f []NodePoolConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigMap flattens the contents of NodePoolConfig from a JSON
// response object.
func flattenNodePoolConfigMap(c *Client, i interface{}) map[string]NodePoolConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfig{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfig{}
	}

	items := make(map[string]NodePoolConfig)
	for k, item := range a {
		items[k] = *flattenNodePoolConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenNodePoolConfigSlice flattens the contents of NodePoolConfig from a JSON
// response object.
func flattenNodePoolConfigSlice(c *Client, i interface{}) []NodePoolConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfig{}
	}

	if len(a) == 0 {
		return []NodePoolConfig{}
	}

	items := make([]NodePoolConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandNodePoolConfig expands an instance of NodePoolConfig into a JSON
// request object.
func expandNodePoolConfig(c *Client, f *NodePoolConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.InstanceType; !dcl.IsEmptyValueIndirect(v) {
		m["instanceType"] = v
	}
	if v, err := expandNodePoolConfigRootVolume(c, f.RootVolume); err != nil {
		return nil, fmt.Errorf("error expanding RootVolume into rootVolume: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["rootVolume"] = v
	}
	if v, err := expandNodePoolConfigTaintsSlice(c, f.Taints); err != nil {
		return nil, fmt.Errorf("error expanding Taints into taints: %w", err)
	} else if v != nil {
		m["taints"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		m["labels"] = v
	}
	if v := f.Tags; !dcl.IsEmptyValueIndirect(v) {
		m["tags"] = v
	}
	if v := f.IamInstanceProfile; !dcl.IsEmptyValueIndirect(v) {
		m["iamInstanceProfile"] = v
	}
	if v, err := expandNodePoolConfigConfigEncryption(c, f.ConfigEncryption); err != nil {
		return nil, fmt.Errorf("error expanding ConfigEncryption into configEncryption: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["configEncryption"] = v
	}
	if v, err := expandNodePoolConfigSshConfig(c, f.SshConfig); err != nil {
		return nil, fmt.Errorf("error expanding SshConfig into sshConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sshConfig"] = v
	}
	if v := f.SecurityGroupIds; v != nil {
		m["securityGroupIds"] = v
	}

	return m, nil
}

// flattenNodePoolConfig flattens an instance of NodePoolConfig from a JSON
// response object.
func flattenNodePoolConfig(c *Client, i interface{}) *NodePoolConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfig
	}
	r.InstanceType = dcl.FlattenString(m["instanceType"])
	r.RootVolume = flattenNodePoolConfigRootVolume(c, m["rootVolume"])
	r.Taints = flattenNodePoolConfigTaintsSlice(c, m["taints"])
	r.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	r.Tags = dcl.FlattenKeyValuePairs(m["tags"])
	r.IamInstanceProfile = dcl.FlattenString(m["iamInstanceProfile"])
	r.ConfigEncryption = flattenNodePoolConfigConfigEncryption(c, m["configEncryption"])
	r.SshConfig = flattenNodePoolConfigSshConfig(c, m["sshConfig"])
	r.SecurityGroupIds = dcl.FlattenStringSlice(m["securityGroupIds"])

	return r
}

// expandNodePoolConfigRootVolumeMap expands the contents of NodePoolConfigRootVolume into a JSON
// request object.
func expandNodePoolConfigRootVolumeMap(c *Client, f map[string]NodePoolConfigRootVolume) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfigRootVolume(c, &item)
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
func expandNodePoolConfigRootVolumeSlice(c *Client, f []NodePoolConfigRootVolume) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfigRootVolume(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigRootVolumeMap flattens the contents of NodePoolConfigRootVolume from a JSON
// response object.
func flattenNodePoolConfigRootVolumeMap(c *Client, i interface{}) map[string]NodePoolConfigRootVolume {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigRootVolume{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigRootVolume{}
	}

	items := make(map[string]NodePoolConfigRootVolume)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigRootVolume(c, item.(map[string]interface{}))
	}

	return items
}

// flattenNodePoolConfigRootVolumeSlice flattens the contents of NodePoolConfigRootVolume from a JSON
// response object.
func flattenNodePoolConfigRootVolumeSlice(c *Client, i interface{}) []NodePoolConfigRootVolume {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigRootVolume{}
	}

	if len(a) == 0 {
		return []NodePoolConfigRootVolume{}
	}

	items := make([]NodePoolConfigRootVolume, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigRootVolume(c, item.(map[string]interface{})))
	}

	return items
}

// expandNodePoolConfigRootVolume expands an instance of NodePoolConfigRootVolume into a JSON
// request object.
func expandNodePoolConfigRootVolume(c *Client, f *NodePoolConfigRootVolume) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.SizeGib; !dcl.IsEmptyValueIndirect(v) {
		m["sizeGib"] = v
	}
	if v := f.VolumeType; !dcl.IsEmptyValueIndirect(v) {
		m["volumeType"] = v
	}
	if v := f.Iops; !dcl.IsEmptyValueIndirect(v) {
		m["iops"] = v
	}
	if v := f.KmsKeyArn; !dcl.IsEmptyValueIndirect(v) {
		m["kmsKeyArn"] = v
	}

	return m, nil
}

// flattenNodePoolConfigRootVolume flattens an instance of NodePoolConfigRootVolume from a JSON
// response object.
func flattenNodePoolConfigRootVolume(c *Client, i interface{}) *NodePoolConfigRootVolume {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfigRootVolume{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfigRootVolume
	}
	r.SizeGib = dcl.FlattenInteger(m["sizeGib"])
	r.VolumeType = flattenNodePoolConfigRootVolumeVolumeTypeEnum(m["volumeType"])
	r.Iops = dcl.FlattenInteger(m["iops"])
	r.KmsKeyArn = dcl.FlattenString(m["kmsKeyArn"])

	return r
}

// expandNodePoolConfigTaintsMap expands the contents of NodePoolConfigTaints into a JSON
// request object.
func expandNodePoolConfigTaintsMap(c *Client, f map[string]NodePoolConfigTaints) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfigTaints(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNodePoolConfigTaintsSlice expands the contents of NodePoolConfigTaints into a JSON
// request object.
func expandNodePoolConfigTaintsSlice(c *Client, f []NodePoolConfigTaints) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfigTaints(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigTaintsMap flattens the contents of NodePoolConfigTaints from a JSON
// response object.
func flattenNodePoolConfigTaintsMap(c *Client, i interface{}) map[string]NodePoolConfigTaints {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigTaints{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigTaints{}
	}

	items := make(map[string]NodePoolConfigTaints)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigTaints(c, item.(map[string]interface{}))
	}

	return items
}

// flattenNodePoolConfigTaintsSlice flattens the contents of NodePoolConfigTaints from a JSON
// response object.
func flattenNodePoolConfigTaintsSlice(c *Client, i interface{}) []NodePoolConfigTaints {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigTaints{}
	}

	if len(a) == 0 {
		return []NodePoolConfigTaints{}
	}

	items := make([]NodePoolConfigTaints, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigTaints(c, item.(map[string]interface{})))
	}

	return items
}

// expandNodePoolConfigTaints expands an instance of NodePoolConfigTaints into a JSON
// request object.
func expandNodePoolConfigTaints(c *Client, f *NodePoolConfigTaints) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Key; !dcl.IsEmptyValueIndirect(v) {
		m["key"] = v
	}
	if v := f.Value; !dcl.IsEmptyValueIndirect(v) {
		m["value"] = v
	}
	if v := f.Effect; !dcl.IsEmptyValueIndirect(v) {
		m["effect"] = v
	}

	return m, nil
}

// flattenNodePoolConfigTaints flattens an instance of NodePoolConfigTaints from a JSON
// response object.
func flattenNodePoolConfigTaints(c *Client, i interface{}) *NodePoolConfigTaints {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfigTaints{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfigTaints
	}
	r.Key = dcl.FlattenString(m["key"])
	r.Value = dcl.FlattenString(m["value"])
	r.Effect = flattenNodePoolConfigTaintsEffectEnum(m["effect"])

	return r
}

// expandNodePoolConfigConfigEncryptionMap expands the contents of NodePoolConfigConfigEncryption into a JSON
// request object.
func expandNodePoolConfigConfigEncryptionMap(c *Client, f map[string]NodePoolConfigConfigEncryption) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfigConfigEncryption(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandNodePoolConfigConfigEncryptionSlice expands the contents of NodePoolConfigConfigEncryption into a JSON
// request object.
func expandNodePoolConfigConfigEncryptionSlice(c *Client, f []NodePoolConfigConfigEncryption) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfigConfigEncryption(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigConfigEncryptionMap flattens the contents of NodePoolConfigConfigEncryption from a JSON
// response object.
func flattenNodePoolConfigConfigEncryptionMap(c *Client, i interface{}) map[string]NodePoolConfigConfigEncryption {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigConfigEncryption{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigConfigEncryption{}
	}

	items := make(map[string]NodePoolConfigConfigEncryption)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigConfigEncryption(c, item.(map[string]interface{}))
	}

	return items
}

// flattenNodePoolConfigConfigEncryptionSlice flattens the contents of NodePoolConfigConfigEncryption from a JSON
// response object.
func flattenNodePoolConfigConfigEncryptionSlice(c *Client, i interface{}) []NodePoolConfigConfigEncryption {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigConfigEncryption{}
	}

	if len(a) == 0 {
		return []NodePoolConfigConfigEncryption{}
	}

	items := make([]NodePoolConfigConfigEncryption, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigConfigEncryption(c, item.(map[string]interface{})))
	}

	return items
}

// expandNodePoolConfigConfigEncryption expands an instance of NodePoolConfigConfigEncryption into a JSON
// request object.
func expandNodePoolConfigConfigEncryption(c *Client, f *NodePoolConfigConfigEncryption) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KmsKeyArn; !dcl.IsEmptyValueIndirect(v) {
		m["kmsKeyArn"] = v
	}

	return m, nil
}

// flattenNodePoolConfigConfigEncryption flattens an instance of NodePoolConfigConfigEncryption from a JSON
// response object.
func flattenNodePoolConfigConfigEncryption(c *Client, i interface{}) *NodePoolConfigConfigEncryption {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfigConfigEncryption{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfigConfigEncryption
	}
	r.KmsKeyArn = dcl.FlattenString(m["kmsKeyArn"])

	return r
}

// expandNodePoolConfigSshConfigMap expands the contents of NodePoolConfigSshConfig into a JSON
// request object.
func expandNodePoolConfigSshConfigMap(c *Client, f map[string]NodePoolConfigSshConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolConfigSshConfig(c, &item)
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
func expandNodePoolConfigSshConfigSlice(c *Client, f []NodePoolConfigSshConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolConfigSshConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolConfigSshConfigMap flattens the contents of NodePoolConfigSshConfig from a JSON
// response object.
func flattenNodePoolConfigSshConfigMap(c *Client, i interface{}) map[string]NodePoolConfigSshConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigSshConfig{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigSshConfig{}
	}

	items := make(map[string]NodePoolConfigSshConfig)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigSshConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenNodePoolConfigSshConfigSlice flattens the contents of NodePoolConfigSshConfig from a JSON
// response object.
func flattenNodePoolConfigSshConfigSlice(c *Client, i interface{}) []NodePoolConfigSshConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigSshConfig{}
	}

	if len(a) == 0 {
		return []NodePoolConfigSshConfig{}
	}

	items := make([]NodePoolConfigSshConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigSshConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandNodePoolConfigSshConfig expands an instance of NodePoolConfigSshConfig into a JSON
// request object.
func expandNodePoolConfigSshConfig(c *Client, f *NodePoolConfigSshConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Ec2KeyPair; !dcl.IsEmptyValueIndirect(v) {
		m["ec2KeyPair"] = v
	}

	return m, nil
}

// flattenNodePoolConfigSshConfig flattens an instance of NodePoolConfigSshConfig from a JSON
// response object.
func flattenNodePoolConfigSshConfig(c *Client, i interface{}) *NodePoolConfigSshConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &NodePoolConfigSshConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyNodePoolConfigSshConfig
	}
	r.Ec2KeyPair = dcl.FlattenString(m["ec2KeyPair"])

	return r
}

// expandNodePoolAutoscalingMap expands the contents of NodePoolAutoscaling into a JSON
// request object.
func expandNodePoolAutoscalingMap(c *Client, f map[string]NodePoolAutoscaling) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolAutoscaling(c, &item)
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
func expandNodePoolAutoscalingSlice(c *Client, f []NodePoolAutoscaling) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolAutoscaling(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolAutoscalingMap flattens the contents of NodePoolAutoscaling from a JSON
// response object.
func flattenNodePoolAutoscalingMap(c *Client, i interface{}) map[string]NodePoolAutoscaling {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolAutoscaling{}
	}

	if len(a) == 0 {
		return map[string]NodePoolAutoscaling{}
	}

	items := make(map[string]NodePoolAutoscaling)
	for k, item := range a {
		items[k] = *flattenNodePoolAutoscaling(c, item.(map[string]interface{}))
	}

	return items
}

// flattenNodePoolAutoscalingSlice flattens the contents of NodePoolAutoscaling from a JSON
// response object.
func flattenNodePoolAutoscalingSlice(c *Client, i interface{}) []NodePoolAutoscaling {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolAutoscaling{}
	}

	if len(a) == 0 {
		return []NodePoolAutoscaling{}
	}

	items := make([]NodePoolAutoscaling, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolAutoscaling(c, item.(map[string]interface{})))
	}

	return items
}

// expandNodePoolAutoscaling expands an instance of NodePoolAutoscaling into a JSON
// request object.
func expandNodePoolAutoscaling(c *Client, f *NodePoolAutoscaling) (map[string]interface{}, error) {
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
func flattenNodePoolAutoscaling(c *Client, i interface{}) *NodePoolAutoscaling {
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
func expandNodePoolMaxPodsConstraintMap(c *Client, f map[string]NodePoolMaxPodsConstraint) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandNodePoolMaxPodsConstraint(c, &item)
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
func expandNodePoolMaxPodsConstraintSlice(c *Client, f []NodePoolMaxPodsConstraint) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandNodePoolMaxPodsConstraint(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenNodePoolMaxPodsConstraintMap flattens the contents of NodePoolMaxPodsConstraint from a JSON
// response object.
func flattenNodePoolMaxPodsConstraintMap(c *Client, i interface{}) map[string]NodePoolMaxPodsConstraint {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolMaxPodsConstraint{}
	}

	if len(a) == 0 {
		return map[string]NodePoolMaxPodsConstraint{}
	}

	items := make(map[string]NodePoolMaxPodsConstraint)
	for k, item := range a {
		items[k] = *flattenNodePoolMaxPodsConstraint(c, item.(map[string]interface{}))
	}

	return items
}

// flattenNodePoolMaxPodsConstraintSlice flattens the contents of NodePoolMaxPodsConstraint from a JSON
// response object.
func flattenNodePoolMaxPodsConstraintSlice(c *Client, i interface{}) []NodePoolMaxPodsConstraint {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolMaxPodsConstraint{}
	}

	if len(a) == 0 {
		return []NodePoolMaxPodsConstraint{}
	}

	items := make([]NodePoolMaxPodsConstraint, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolMaxPodsConstraint(c, item.(map[string]interface{})))
	}

	return items
}

// expandNodePoolMaxPodsConstraint expands an instance of NodePoolMaxPodsConstraint into a JSON
// request object.
func expandNodePoolMaxPodsConstraint(c *Client, f *NodePoolMaxPodsConstraint) (map[string]interface{}, error) {
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
func flattenNodePoolMaxPodsConstraint(c *Client, i interface{}) *NodePoolMaxPodsConstraint {
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

// flattenNodePoolConfigRootVolumeVolumeTypeEnumMap flattens the contents of NodePoolConfigRootVolumeVolumeTypeEnum from a JSON
// response object.
func flattenNodePoolConfigRootVolumeVolumeTypeEnumMap(c *Client, i interface{}) map[string]NodePoolConfigRootVolumeVolumeTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigRootVolumeVolumeTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigRootVolumeVolumeTypeEnum{}
	}

	items := make(map[string]NodePoolConfigRootVolumeVolumeTypeEnum)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigRootVolumeVolumeTypeEnum(item.(interface{}))
	}

	return items
}

// flattenNodePoolConfigRootVolumeVolumeTypeEnumSlice flattens the contents of NodePoolConfigRootVolumeVolumeTypeEnum from a JSON
// response object.
func flattenNodePoolConfigRootVolumeVolumeTypeEnumSlice(c *Client, i interface{}) []NodePoolConfigRootVolumeVolumeTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigRootVolumeVolumeTypeEnum{}
	}

	if len(a) == 0 {
		return []NodePoolConfigRootVolumeVolumeTypeEnum{}
	}

	items := make([]NodePoolConfigRootVolumeVolumeTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigRootVolumeVolumeTypeEnum(item.(interface{})))
	}

	return items
}

// flattenNodePoolConfigRootVolumeVolumeTypeEnum asserts that an interface is a string, and returns a
// pointer to a *NodePoolConfigRootVolumeVolumeTypeEnum with the same value as that string.
func flattenNodePoolConfigRootVolumeVolumeTypeEnum(i interface{}) *NodePoolConfigRootVolumeVolumeTypeEnum {
	s, ok := i.(string)
	if !ok {
		return NodePoolConfigRootVolumeVolumeTypeEnumRef("")
	}

	return NodePoolConfigRootVolumeVolumeTypeEnumRef(s)
}

// flattenNodePoolConfigTaintsEffectEnumMap flattens the contents of NodePoolConfigTaintsEffectEnum from a JSON
// response object.
func flattenNodePoolConfigTaintsEffectEnumMap(c *Client, i interface{}) map[string]NodePoolConfigTaintsEffectEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]NodePoolConfigTaintsEffectEnum{}
	}

	if len(a) == 0 {
		return map[string]NodePoolConfigTaintsEffectEnum{}
	}

	items := make(map[string]NodePoolConfigTaintsEffectEnum)
	for k, item := range a {
		items[k] = *flattenNodePoolConfigTaintsEffectEnum(item.(interface{}))
	}

	return items
}

// flattenNodePoolConfigTaintsEffectEnumSlice flattens the contents of NodePoolConfigTaintsEffectEnum from a JSON
// response object.
func flattenNodePoolConfigTaintsEffectEnumSlice(c *Client, i interface{}) []NodePoolConfigTaintsEffectEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []NodePoolConfigTaintsEffectEnum{}
	}

	if len(a) == 0 {
		return []NodePoolConfigTaintsEffectEnum{}
	}

	items := make([]NodePoolConfigTaintsEffectEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenNodePoolConfigTaintsEffectEnum(item.(interface{})))
	}

	return items
}

// flattenNodePoolConfigTaintsEffectEnum asserts that an interface is a string, and returns a
// pointer to a *NodePoolConfigTaintsEffectEnum with the same value as that string.
func flattenNodePoolConfigTaintsEffectEnum(i interface{}) *NodePoolConfigTaintsEffectEnum {
	s, ok := i.(string)
	if !ok {
		return NodePoolConfigTaintsEffectEnumRef("")
	}

	return NodePoolConfigTaintsEffectEnumRef(s)
}

// flattenNodePoolStateEnumMap flattens the contents of NodePoolStateEnum from a JSON
// response object.
func flattenNodePoolStateEnumMap(c *Client, i interface{}) map[string]NodePoolStateEnum {
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
func flattenNodePoolStateEnumSlice(c *Client, i interface{}) []NodePoolStateEnum {
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
		return NodePoolStateEnumRef("")
	}

	return NodePoolStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *NodePool) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalNodePool(b, c)
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
		diff := nodePoolDiff{}
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

	case "updateNodePoolUpdateAwsNodePoolOperation":
		return &updateNodePoolUpdateAwsNodePoolOperation{FieldDiffs: fieldDiffs}, nil

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
	if !dcl.IsNotReturnedByServer(vConfig) {
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
	if !dcl.IsNotReturnedByServer(vAutoscaling) {
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
	if !dcl.IsNotReturnedByServer(vMaxPodsConstraint) {
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
	if !dcl.IsNotReturnedByServer(vRootVolume) {
		o.RootVolume = vRootVolume
	}
	vConfigEncryption := o.ConfigEncryption
	if vConfigEncryption == nil {
		// note: explicitly not the empty object.
		vConfigEncryption = &NodePoolConfigConfigEncryption{}
	}
	if err := extractNodePoolConfigConfigEncryptionFields(r, vConfigEncryption); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vConfigEncryption) {
		o.ConfigEncryption = vConfigEncryption
	}
	vSshConfig := o.SshConfig
	if vSshConfig == nil {
		// note: explicitly not the empty object.
		vSshConfig = &NodePoolConfigSshConfig{}
	}
	if err := extractNodePoolConfigSshConfigFields(r, vSshConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSshConfig) {
		o.SshConfig = vSshConfig
	}
	return nil
}
func extractNodePoolConfigRootVolumeFields(r *NodePool, o *NodePoolConfigRootVolume) error {
	return nil
}
func extractNodePoolConfigTaintsFields(r *NodePool, o *NodePoolConfigTaints) error {
	return nil
}
func extractNodePoolConfigConfigEncryptionFields(r *NodePool, o *NodePoolConfigConfigEncryption) error {
	return nil
}
func extractNodePoolConfigSshConfigFields(r *NodePool, o *NodePoolConfigSshConfig) error {
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
	if !dcl.IsNotReturnedByServer(vConfig) {
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
	if !dcl.IsNotReturnedByServer(vAutoscaling) {
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
	if !dcl.IsNotReturnedByServer(vMaxPodsConstraint) {
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
	if !dcl.IsNotReturnedByServer(vRootVolume) {
		o.RootVolume = vRootVolume
	}
	vConfigEncryption := o.ConfigEncryption
	if vConfigEncryption == nil {
		// note: explicitly not the empty object.
		vConfigEncryption = &NodePoolConfigConfigEncryption{}
	}
	if err := extractNodePoolConfigConfigEncryptionFields(r, vConfigEncryption); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vConfigEncryption) {
		o.ConfigEncryption = vConfigEncryption
	}
	vSshConfig := o.SshConfig
	if vSshConfig == nil {
		// note: explicitly not the empty object.
		vSshConfig = &NodePoolConfigSshConfig{}
	}
	if err := extractNodePoolConfigSshConfigFields(r, vSshConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSshConfig) {
		o.SshConfig = vSshConfig
	}
	return nil
}
func postReadExtractNodePoolConfigRootVolumeFields(r *NodePool, o *NodePoolConfigRootVolume) error {
	return nil
}
func postReadExtractNodePoolConfigTaintsFields(r *NodePool, o *NodePoolConfigTaints) error {
	return nil
}
func postReadExtractNodePoolConfigConfigEncryptionFields(r *NodePool, o *NodePoolConfigConfigEncryption) error {
	return nil
}
func postReadExtractNodePoolConfigSshConfigFields(r *NodePool, o *NodePoolConfigSshConfig) error {
	return nil
}
func postReadExtractNodePoolAutoscalingFields(r *NodePool, o *NodePoolAutoscaling) error {
	return nil
}
func postReadExtractNodePoolMaxPodsConstraintFields(r *NodePool, o *NodePoolMaxPodsConstraint) error {
	return nil
}
