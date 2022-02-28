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
package dataproc

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

func (r *Cluster) validate() error {

	if err := dcl.Required(r, "project"); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Config) {
		if err := r.Config.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Status) {
		if err := r.Status.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Metrics) {
		if err := r.Metrics.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterClusterConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.GceClusterConfig) {
		if err := r.GceClusterConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.MasterConfig) {
		if err := r.MasterConfig.validate(); err != nil {
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
	if !dcl.IsEmptyValueIndirect(r.SoftwareConfig) {
		if err := r.SoftwareConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.EncryptionConfig) {
		if err := r.EncryptionConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AutoscalingConfig) {
		if err := r.AutoscalingConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SecurityConfig) {
		if err := r.SecurityConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.LifecycleConfig) {
		if err := r.LifecycleConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.EndpointConfig) {
		if err := r.EndpointConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterClusterConfigGceClusterConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.ReservationAffinity) {
		if err := r.ReservationAffinity.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.NodeGroupAffinity) {
		if err := r.NodeGroupAffinity.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterClusterConfigGceClusterConfigReservationAffinity) validate() error {
	return nil
}
func (r *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) validate() error {
	if err := dcl.Required(r, "nodeGroup"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterInstanceGroupConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.DiskConfig) {
		if err := r.DiskConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ManagedGroupConfig) {
		if err := r.ManagedGroupConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterInstanceGroupConfigDiskConfig) validate() error {
	return nil
}
func (r *ClusterInstanceGroupConfigManagedGroupConfig) validate() error {
	return nil
}
func (r *ClusterInstanceGroupConfigAccelerators) validate() error {
	return nil
}
func (r *ClusterClusterConfigSoftwareConfig) validate() error {
	return nil
}
func (r *ClusterClusterConfigInitializationActions) validate() error {
	if err := dcl.Required(r, "executableFile"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterClusterConfigEncryptionConfig) validate() error {
	return nil
}
func (r *ClusterClusterConfigAutoscalingConfig) validate() error {
	return nil
}
func (r *ClusterClusterConfigSecurityConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.KerberosConfig) {
		if err := r.KerberosConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterClusterConfigSecurityConfigKerberosConfig) validate() error {
	return nil
}
func (r *ClusterClusterConfigLifecycleConfig) validate() error {
	return nil
}
func (r *ClusterClusterConfigEndpointConfig) validate() error {
	return nil
}
func (r *ClusterStatus) validate() error {
	return nil
}
func (r *ClusterStatusHistory) validate() error {
	return nil
}
func (r *ClusterMetrics) validate() error {
	return nil
}
func (r *Cluster) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://dataproc.googleapis.com/v1/", params)
}

func (r *Cluster) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/clusters/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Cluster) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/clusters", nr.basePath(), userBasePath, params), nil

}

func (r *Cluster) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/clusters", nr.basePath(), userBasePath, params), nil

}

func (r *Cluster) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/clusters/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Cluster) SetPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{
		"project":  *nr.Project,
		"location": *nr.Location,
		"name":     *nr.Name,
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/clusters/{{name}}:setIamPolicy", nr.basePath(), userBasePath, fields)
}

func (r *Cluster) SetPolicyVerb() string {
	return "POST"
}

func (r *Cluster) getPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{
		"project":  *nr.Project,
		"location": *nr.Location,
		"name":     *nr.Name,
	}
	return dcl.URL("projects/{{project}}/regions/{{location}}/clusters/{{name}}:getIamPolicy", nr.basePath(), userBasePath, fields)
}

func (r *Cluster) IAMPolicyVersion() int {
	return 3
}

// clusterApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type clusterApiOperation interface {
	do(context.Context, *Cluster, *Client) error
}

// newUpdateClusterUpdateClusterRequest creates a request for an
// Cluster resource's UpdateCluster update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateClusterUpdateClusterRequest(ctx context.Context, f *Cluster, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}

	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	return req, nil
}

// marshalUpdateClusterUpdateClusterRequest converts the update into
// the final JSON request body.
func marshalUpdateClusterUpdateClusterRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateClusterUpdateClusterOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateClusterUpdateClusterOperation) do(ctx context.Context, r *Cluster, c *Client) error {
	_, err := c.GetCluster(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateCluster")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateClusterUpdateClusterRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateClusterUpdateClusterRequest(c, req)
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

func (c *Client) listClusterRaw(ctx context.Context, r *Cluster, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ClusterMaxPage {
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

type listClusterOperation struct {
	Clusters []map[string]interface{} `json:"clusters"`
	Token    string                   `json:"nextPageToken"`
}

func (c *Client) listCluster(ctx context.Context, r *Cluster, pageToken string, pageSize int32) ([]*Cluster, string, error) {
	b, err := c.listClusterRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listClusterOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Cluster
	for _, v := range m.Clusters {
		res, err := unmarshalMapCluster(v, c)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllCluster(ctx context.Context, f func(*Cluster) bool, resources []*Cluster) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteCluster(ctx, res)
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

type deleteClusterOperation struct{}

func (op *deleteClusterOperation) do(ctx context.Context, r *Cluster, c *Client) error {
	r, err := c.GetCluster(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Cluster not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetCluster checking for existence. error: %v", err)
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
		_, err = c.GetCluster(ctx, r)
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
type createClusterOperation struct {
	response map[string]interface{}
}

func (op *createClusterOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createClusterOperation) do(ctx context.Context, r *Cluster, c *Client) error {
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

	if _, err := c.GetCluster(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getClusterRaw(ctx context.Context, r *Cluster) ([]byte, error) {

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

func (c *Client) clusterDiffsForRawDesired(ctx context.Context, rawDesired *Cluster, opts ...dcl.ApplyOption) (initial, desired *Cluster, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Cluster
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Cluster); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Cluster, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetCluster(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Cluster resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Cluster resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Cluster resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeClusterDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Cluster: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Cluster: %v", rawDesired)

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeClusterInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Cluster: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeClusterDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Cluster: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffCluster(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeClusterInitialState(rawInitial, rawDesired *Cluster) (*Cluster, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeClusterDesiredState(rawDesired, rawInitial *Cluster, opts ...dcl.ApplyOption) (*Cluster, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Config = canonicalizeClusterClusterConfig(rawDesired.Config, nil, opts...)
		rawDesired.Status = canonicalizeClusterStatus(rawDesired.Status, nil, opts...)
		rawDesired.Metrics = canonicalizeClusterMetrics(rawDesired.Metrics, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Cluster{}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	canonicalDesired.Config = canonicalizeClusterClusterConfig(rawDesired.Config, rawInitial.Config, opts...)
	if dcl.IsZeroValue(rawDesired.Labels) {
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}

	return canonicalDesired, nil
}

func canonicalizeClusterNewState(c *Client, rawNew, rawDesired *Cluster) (*Cluster, error) {

	if dcl.IsNotReturnedByServer(rawNew.Project) && dcl.IsNotReturnedByServer(rawDesired.Project) {
		rawNew.Project = rawDesired.Project
	} else {
		if dcl.NameToSelfLink(rawDesired.Project, rawNew.Project) {
			rawNew.Project = rawDesired.Project
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Name) && dcl.IsNotReturnedByServer(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Config) && dcl.IsNotReturnedByServer(rawDesired.Config) {
		rawNew.Config = rawDesired.Config
	} else {
		rawNew.Config = canonicalizeNewClusterClusterConfig(c, rawDesired.Config, rawNew.Config)
	}

	if dcl.IsNotReturnedByServer(rawNew.Labels) && dcl.IsNotReturnedByServer(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Status) && dcl.IsNotReturnedByServer(rawDesired.Status) {
		rawNew.Status = rawDesired.Status
	} else {
		rawNew.Status = canonicalizeNewClusterStatus(c, rawDesired.Status, rawNew.Status)
	}

	if dcl.IsNotReturnedByServer(rawNew.StatusHistory) && dcl.IsNotReturnedByServer(rawDesired.StatusHistory) {
		rawNew.StatusHistory = rawDesired.StatusHistory
	} else {
		rawNew.StatusHistory = canonicalizeNewClusterStatusHistorySlice(c, rawDesired.StatusHistory, rawNew.StatusHistory)
	}

	if dcl.IsNotReturnedByServer(rawNew.ClusterUuid) && dcl.IsNotReturnedByServer(rawDesired.ClusterUuid) {
		rawNew.ClusterUuid = rawDesired.ClusterUuid
	} else {
		if dcl.StringCanonicalize(rawDesired.ClusterUuid, rawNew.ClusterUuid) {
			rawNew.ClusterUuid = rawDesired.ClusterUuid
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Metrics) && dcl.IsNotReturnedByServer(rawDesired.Metrics) {
		rawNew.Metrics = rawDesired.Metrics
	} else {
		rawNew.Metrics = canonicalizeNewClusterMetrics(c, rawDesired.Metrics, rawNew.Metrics)
	}

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeClusterClusterConfig(des, initial *ClusterClusterConfig, opts ...dcl.ApplyOption) *ClusterClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfig{}

	if dcl.NameToSelfLink(des.StagingBucket, initial.StagingBucket) || dcl.IsZeroValue(des.StagingBucket) {
		cDes.StagingBucket = initial.StagingBucket
	} else {
		cDes.StagingBucket = des.StagingBucket
	}
	if dcl.NameToSelfLink(des.TempBucket, initial.TempBucket) || dcl.IsZeroValue(des.TempBucket) {
		cDes.TempBucket = initial.TempBucket
	} else {
		cDes.TempBucket = des.TempBucket
	}
	cDes.GceClusterConfig = canonicalizeClusterClusterConfigGceClusterConfig(des.GceClusterConfig, initial.GceClusterConfig, opts...)
	cDes.MasterConfig = canonicalizeClusterInstanceGroupConfig(des.MasterConfig, initial.MasterConfig, opts...)
	cDes.WorkerConfig = canonicalizeClusterInstanceGroupConfig(des.WorkerConfig, initial.WorkerConfig, opts...)
	cDes.SecondaryWorkerConfig = canonicalizeClusterInstanceGroupConfig(des.SecondaryWorkerConfig, initial.SecondaryWorkerConfig, opts...)
	cDes.SoftwareConfig = canonicalizeClusterClusterConfigSoftwareConfig(des.SoftwareConfig, initial.SoftwareConfig, opts...)
	cDes.InitializationActions = canonicalizeClusterClusterConfigInitializationActionsSlice(des.InitializationActions, initial.InitializationActions, opts...)
	cDes.EncryptionConfig = canonicalizeClusterClusterConfigEncryptionConfig(des.EncryptionConfig, initial.EncryptionConfig, opts...)
	cDes.AutoscalingConfig = canonicalizeClusterClusterConfigAutoscalingConfig(des.AutoscalingConfig, initial.AutoscalingConfig, opts...)
	cDes.SecurityConfig = canonicalizeClusterClusterConfigSecurityConfig(des.SecurityConfig, initial.SecurityConfig, opts...)
	cDes.LifecycleConfig = canonicalizeClusterClusterConfigLifecycleConfig(des.LifecycleConfig, initial.LifecycleConfig, opts...)
	cDes.EndpointConfig = canonicalizeClusterClusterConfigEndpointConfig(des.EndpointConfig, initial.EndpointConfig, opts...)

	return cDes
}

func canonicalizeClusterClusterConfigSlice(des, initial []ClusterClusterConfig, opts ...dcl.ApplyOption) []ClusterClusterConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfig(c *Client, des, nw *ClusterClusterConfig) *ClusterClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.NameToSelfLink(des.StagingBucket, nw.StagingBucket) {
		nw.StagingBucket = des.StagingBucket
	}
	if dcl.NameToSelfLink(des.TempBucket, nw.TempBucket) {
		nw.TempBucket = des.TempBucket
	}
	nw.GceClusterConfig = canonicalizeNewClusterClusterConfigGceClusterConfig(c, des.GceClusterConfig, nw.GceClusterConfig)
	nw.MasterConfig = canonicalizeNewClusterInstanceGroupConfig(c, des.MasterConfig, nw.MasterConfig)
	nw.WorkerConfig = canonicalizeNewClusterInstanceGroupConfig(c, des.WorkerConfig, nw.WorkerConfig)
	nw.SecondaryWorkerConfig = canonicalizeNewClusterInstanceGroupConfig(c, des.SecondaryWorkerConfig, nw.SecondaryWorkerConfig)
	nw.SoftwareConfig = canonicalizeNewClusterClusterConfigSoftwareConfig(c, des.SoftwareConfig, nw.SoftwareConfig)
	nw.InitializationActions = canonicalizeNewClusterClusterConfigInitializationActionsSlice(c, des.InitializationActions, nw.InitializationActions)
	nw.EncryptionConfig = canonicalizeNewClusterClusterConfigEncryptionConfig(c, des.EncryptionConfig, nw.EncryptionConfig)
	nw.AutoscalingConfig = canonicalizeNewClusterClusterConfigAutoscalingConfig(c, des.AutoscalingConfig, nw.AutoscalingConfig)
	nw.SecurityConfig = canonicalizeNewClusterClusterConfigSecurityConfig(c, des.SecurityConfig, nw.SecurityConfig)
	nw.LifecycleConfig = canonicalizeNewClusterClusterConfigLifecycleConfig(c, des.LifecycleConfig, nw.LifecycleConfig)
	nw.EndpointConfig = canonicalizeNewClusterClusterConfigEndpointConfig(c, des.EndpointConfig, nw.EndpointConfig)

	return nw
}

func canonicalizeNewClusterClusterConfigSet(c *Client, des, nw []ClusterClusterConfig) []ClusterClusterConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigSlice(c *Client, des, nw []ClusterClusterConfig) []ClusterClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigGceClusterConfig(des, initial *ClusterClusterConfigGceClusterConfig, opts ...dcl.ApplyOption) *ClusterClusterConfigGceClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigGceClusterConfig{}

	if dcl.StringCanonicalize(des.Zone, initial.Zone) || dcl.IsZeroValue(des.Zone) {
		cDes.Zone = initial.Zone
	} else {
		cDes.Zone = des.Zone
	}
	if dcl.NameToSelfLink(des.Network, initial.Network) || dcl.IsZeroValue(des.Network) {
		cDes.Network = initial.Network
	} else {
		cDes.Network = des.Network
	}
	if dcl.NameToSelfLink(des.Subnetwork, initial.Subnetwork) || dcl.IsZeroValue(des.Subnetwork) {
		cDes.Subnetwork = initial.Subnetwork
	} else {
		cDes.Subnetwork = des.Subnetwork
	}
	if dcl.BoolCanonicalize(des.InternalIPOnly, initial.InternalIPOnly) || dcl.IsZeroValue(des.InternalIPOnly) {
		cDes.InternalIPOnly = initial.InternalIPOnly
	} else {
		cDes.InternalIPOnly = des.InternalIPOnly
	}
	if dcl.IsZeroValue(des.PrivateIPv6GoogleAccess) {
		cDes.PrivateIPv6GoogleAccess = initial.PrivateIPv6GoogleAccess
	} else {
		cDes.PrivateIPv6GoogleAccess = des.PrivateIPv6GoogleAccess
	}
	if dcl.NameToSelfLink(des.ServiceAccount, initial.ServiceAccount) || dcl.IsZeroValue(des.ServiceAccount) {
		cDes.ServiceAccount = initial.ServiceAccount
	} else {
		cDes.ServiceAccount = des.ServiceAccount
	}
	if dcl.StringArrayCanonicalize(des.ServiceAccountScopes, initial.ServiceAccountScopes) || dcl.IsZeroValue(des.ServiceAccountScopes) {
		cDes.ServiceAccountScopes = initial.ServiceAccountScopes
	} else {
		cDes.ServiceAccountScopes = des.ServiceAccountScopes
	}
	if dcl.StringArrayCanonicalize(des.Tags, initial.Tags) || dcl.IsZeroValue(des.Tags) {
		cDes.Tags = initial.Tags
	} else {
		cDes.Tags = des.Tags
	}
	if dcl.IsZeroValue(des.Metadata) {
		cDes.Metadata = initial.Metadata
	} else {
		cDes.Metadata = des.Metadata
	}
	cDes.ReservationAffinity = canonicalizeClusterClusterConfigGceClusterConfigReservationAffinity(des.ReservationAffinity, initial.ReservationAffinity, opts...)
	cDes.NodeGroupAffinity = canonicalizeClusterClusterConfigGceClusterConfigNodeGroupAffinity(des.NodeGroupAffinity, initial.NodeGroupAffinity, opts...)

	return cDes
}

func canonicalizeClusterClusterConfigGceClusterConfigSlice(des, initial []ClusterClusterConfigGceClusterConfig, opts ...dcl.ApplyOption) []ClusterClusterConfigGceClusterConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigGceClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigGceClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigGceClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigGceClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigGceClusterConfig(c *Client, des, nw *ClusterClusterConfigGceClusterConfig) *ClusterClusterConfigGceClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigGceClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Zone, nw.Zone) {
		nw.Zone = des.Zone
	}
	if dcl.NameToSelfLink(des.Network, nw.Network) {
		nw.Network = des.Network
	}
	if dcl.NameToSelfLink(des.Subnetwork, nw.Subnetwork) {
		nw.Subnetwork = des.Subnetwork
	}
	if dcl.BoolCanonicalize(des.InternalIPOnly, nw.InternalIPOnly) {
		nw.InternalIPOnly = des.InternalIPOnly
	}
	if dcl.NameToSelfLink(des.ServiceAccount, nw.ServiceAccount) {
		nw.ServiceAccount = des.ServiceAccount
	}
	if dcl.StringArrayCanonicalize(des.ServiceAccountScopes, nw.ServiceAccountScopes) {
		nw.ServiceAccountScopes = des.ServiceAccountScopes
	}
	if dcl.StringArrayCanonicalize(des.Tags, nw.Tags) {
		nw.Tags = des.Tags
	}
	nw.ReservationAffinity = canonicalizeNewClusterClusterConfigGceClusterConfigReservationAffinity(c, des.ReservationAffinity, nw.ReservationAffinity)
	nw.NodeGroupAffinity = canonicalizeNewClusterClusterConfigGceClusterConfigNodeGroupAffinity(c, des.NodeGroupAffinity, nw.NodeGroupAffinity)

	return nw
}

func canonicalizeNewClusterClusterConfigGceClusterConfigSet(c *Client, des, nw []ClusterClusterConfigGceClusterConfig) []ClusterClusterConfigGceClusterConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigGceClusterConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigGceClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigGceClusterConfigSlice(c *Client, des, nw []ClusterClusterConfigGceClusterConfig) []ClusterClusterConfigGceClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigGceClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigGceClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigGceClusterConfigReservationAffinity(des, initial *ClusterClusterConfigGceClusterConfigReservationAffinity, opts ...dcl.ApplyOption) *ClusterClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigGceClusterConfigReservationAffinity{}

	if dcl.IsZeroValue(des.ConsumeReservationType) {
		cDes.ConsumeReservationType = initial.ConsumeReservationType
	} else {
		cDes.ConsumeReservationType = des.ConsumeReservationType
	}
	if dcl.StringCanonicalize(des.Key, initial.Key) || dcl.IsZeroValue(des.Key) {
		cDes.Key = initial.Key
	} else {
		cDes.Key = des.Key
	}
	if dcl.StringArrayCanonicalize(des.Values, initial.Values) || dcl.IsZeroValue(des.Values) {
		cDes.Values = initial.Values
	} else {
		cDes.Values = des.Values
	}

	return cDes
}

func canonicalizeClusterClusterConfigGceClusterConfigReservationAffinitySlice(des, initial []ClusterClusterConfigGceClusterConfigReservationAffinity, opts ...dcl.ApplyOption) []ClusterClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigGceClusterConfigReservationAffinity, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigGceClusterConfigReservationAffinity(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigGceClusterConfigReservationAffinity, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigGceClusterConfigReservationAffinity(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigGceClusterConfigReservationAffinity(c *Client, des, nw *ClusterClusterConfigGceClusterConfigReservationAffinity) *ClusterClusterConfigGceClusterConfigReservationAffinity {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigGceClusterConfigReservationAffinity while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}
	if dcl.StringArrayCanonicalize(des.Values, nw.Values) {
		nw.Values = des.Values
	}

	return nw
}

func canonicalizeNewClusterClusterConfigGceClusterConfigReservationAffinitySet(c *Client, des, nw []ClusterClusterConfigGceClusterConfigReservationAffinity) []ClusterClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigGceClusterConfigReservationAffinity
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigGceClusterConfigReservationAffinityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, des, nw []ClusterClusterConfigGceClusterConfigReservationAffinity) []ClusterClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigGceClusterConfigReservationAffinity
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigGceClusterConfigReservationAffinity(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigGceClusterConfigNodeGroupAffinity(des, initial *ClusterClusterConfigGceClusterConfigNodeGroupAffinity, opts ...dcl.ApplyOption) *ClusterClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigGceClusterConfigNodeGroupAffinity{}

	if dcl.NameToSelfLink(des.NodeGroup, initial.NodeGroup) || dcl.IsZeroValue(des.NodeGroup) {
		cDes.NodeGroup = initial.NodeGroup
	} else {
		cDes.NodeGroup = des.NodeGroup
	}

	return cDes
}

func canonicalizeClusterClusterConfigGceClusterConfigNodeGroupAffinitySlice(des, initial []ClusterClusterConfigGceClusterConfigNodeGroupAffinity, opts ...dcl.ApplyOption) []ClusterClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigGceClusterConfigNodeGroupAffinity(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigGceClusterConfigNodeGroupAffinity(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, des, nw *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) *ClusterClusterConfigGceClusterConfigNodeGroupAffinity {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigGceClusterConfigNodeGroupAffinity while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.NameToSelfLink(des.NodeGroup, nw.NodeGroup) {
		nw.NodeGroup = des.NodeGroup
	}

	return nw
}

func canonicalizeNewClusterClusterConfigGceClusterConfigNodeGroupAffinitySet(c *Client, des, nw []ClusterClusterConfigGceClusterConfigNodeGroupAffinity) []ClusterClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigGceClusterConfigNodeGroupAffinity
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigGceClusterConfigNodeGroupAffinityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, des, nw []ClusterClusterConfigGceClusterConfigNodeGroupAffinity) []ClusterClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigGceClusterConfigNodeGroupAffinity
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigGceClusterConfigNodeGroupAffinity(c, &d, &n))
	}

	return items
}

func canonicalizeClusterInstanceGroupConfig(des, initial *ClusterInstanceGroupConfig, opts ...dcl.ApplyOption) *ClusterInstanceGroupConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterInstanceGroupConfig{}

	if dcl.IsZeroValue(des.NumInstances) {
		cDes.NumInstances = initial.NumInstances
	} else {
		cDes.NumInstances = des.NumInstances
	}
	if dcl.NameToSelfLink(des.Image, initial.Image) || dcl.IsZeroValue(des.Image) {
		cDes.Image = initial.Image
	} else {
		cDes.Image = des.Image
	}
	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	cDes.DiskConfig = canonicalizeClusterInstanceGroupConfigDiskConfig(des.DiskConfig, initial.DiskConfig, opts...)
	if dcl.IsZeroValue(des.Preemptibility) {
		cDes.Preemptibility = initial.Preemptibility
	} else {
		cDes.Preemptibility = des.Preemptibility
	}
	cDes.Accelerators = canonicalizeClusterInstanceGroupConfigAcceleratorsSlice(des.Accelerators, initial.Accelerators, opts...)
	if dcl.StringCanonicalize(des.MinCpuPlatform, initial.MinCpuPlatform) || dcl.IsZeroValue(des.MinCpuPlatform) {
		cDes.MinCpuPlatform = initial.MinCpuPlatform
	} else {
		cDes.MinCpuPlatform = des.MinCpuPlatform
	}

	return cDes
}

func canonicalizeClusterInstanceGroupConfigSlice(des, initial []ClusterInstanceGroupConfig, opts ...dcl.ApplyOption) []ClusterInstanceGroupConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterInstanceGroupConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterInstanceGroupConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterInstanceGroupConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterInstanceGroupConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterInstanceGroupConfig(c *Client, des, nw *ClusterInstanceGroupConfig) *ClusterInstanceGroupConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterInstanceGroupConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.InstanceNames, nw.InstanceNames) {
		nw.InstanceNames = des.InstanceNames
	}
	if dcl.NameToSelfLink(des.Image, nw.Image) {
		nw.Image = des.Image
	}
	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	nw.DiskConfig = canonicalizeNewClusterInstanceGroupConfigDiskConfig(c, des.DiskConfig, nw.DiskConfig)
	if dcl.BoolCanonicalize(des.IsPreemptible, nw.IsPreemptible) {
		nw.IsPreemptible = des.IsPreemptible
	}
	nw.ManagedGroupConfig = canonicalizeNewClusterInstanceGroupConfigManagedGroupConfig(c, des.ManagedGroupConfig, nw.ManagedGroupConfig)
	nw.Accelerators = canonicalizeNewClusterInstanceGroupConfigAcceleratorsSlice(c, des.Accelerators, nw.Accelerators)
	if dcl.StringCanonicalize(des.MinCpuPlatform, nw.MinCpuPlatform) {
		nw.MinCpuPlatform = des.MinCpuPlatform
	}

	return nw
}

func canonicalizeNewClusterInstanceGroupConfigSet(c *Client, des, nw []ClusterInstanceGroupConfig) []ClusterInstanceGroupConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterInstanceGroupConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterInstanceGroupConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterInstanceGroupConfigSlice(c *Client, des, nw []ClusterInstanceGroupConfig) []ClusterInstanceGroupConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterInstanceGroupConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterInstanceGroupConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterInstanceGroupConfigDiskConfig(des, initial *ClusterInstanceGroupConfigDiskConfig, opts ...dcl.ApplyOption) *ClusterInstanceGroupConfigDiskConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterInstanceGroupConfigDiskConfig{}

	if dcl.StringCanonicalize(des.BootDiskType, initial.BootDiskType) || dcl.IsZeroValue(des.BootDiskType) {
		cDes.BootDiskType = initial.BootDiskType
	} else {
		cDes.BootDiskType = des.BootDiskType
	}
	if dcl.IsZeroValue(des.BootDiskSizeGb) {
		cDes.BootDiskSizeGb = initial.BootDiskSizeGb
	} else {
		cDes.BootDiskSizeGb = des.BootDiskSizeGb
	}
	if dcl.IsZeroValue(des.NumLocalSsds) {
		cDes.NumLocalSsds = initial.NumLocalSsds
	} else {
		cDes.NumLocalSsds = des.NumLocalSsds
	}

	return cDes
}

func canonicalizeClusterInstanceGroupConfigDiskConfigSlice(des, initial []ClusterInstanceGroupConfigDiskConfig, opts ...dcl.ApplyOption) []ClusterInstanceGroupConfigDiskConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterInstanceGroupConfigDiskConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterInstanceGroupConfigDiskConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterInstanceGroupConfigDiskConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterInstanceGroupConfigDiskConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterInstanceGroupConfigDiskConfig(c *Client, des, nw *ClusterInstanceGroupConfigDiskConfig) *ClusterInstanceGroupConfigDiskConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterInstanceGroupConfigDiskConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.BootDiskType, nw.BootDiskType) {
		nw.BootDiskType = des.BootDiskType
	}

	return nw
}

func canonicalizeNewClusterInstanceGroupConfigDiskConfigSet(c *Client, des, nw []ClusterInstanceGroupConfigDiskConfig) []ClusterInstanceGroupConfigDiskConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterInstanceGroupConfigDiskConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterInstanceGroupConfigDiskConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterInstanceGroupConfigDiskConfigSlice(c *Client, des, nw []ClusterInstanceGroupConfigDiskConfig) []ClusterInstanceGroupConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterInstanceGroupConfigDiskConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterInstanceGroupConfigDiskConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterInstanceGroupConfigManagedGroupConfig(des, initial *ClusterInstanceGroupConfigManagedGroupConfig, opts ...dcl.ApplyOption) *ClusterInstanceGroupConfigManagedGroupConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterInstanceGroupConfigManagedGroupConfig{}

	return cDes
}

func canonicalizeClusterInstanceGroupConfigManagedGroupConfigSlice(des, initial []ClusterInstanceGroupConfigManagedGroupConfig, opts ...dcl.ApplyOption) []ClusterInstanceGroupConfigManagedGroupConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterInstanceGroupConfigManagedGroupConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterInstanceGroupConfigManagedGroupConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterInstanceGroupConfigManagedGroupConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterInstanceGroupConfigManagedGroupConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterInstanceGroupConfigManagedGroupConfig(c *Client, des, nw *ClusterInstanceGroupConfigManagedGroupConfig) *ClusterInstanceGroupConfigManagedGroupConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterInstanceGroupConfigManagedGroupConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.InstanceTemplateName, nw.InstanceTemplateName) {
		nw.InstanceTemplateName = des.InstanceTemplateName
	}
	if dcl.StringCanonicalize(des.InstanceGroupManagerName, nw.InstanceGroupManagerName) {
		nw.InstanceGroupManagerName = des.InstanceGroupManagerName
	}

	return nw
}

func canonicalizeNewClusterInstanceGroupConfigManagedGroupConfigSet(c *Client, des, nw []ClusterInstanceGroupConfigManagedGroupConfig) []ClusterInstanceGroupConfigManagedGroupConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterInstanceGroupConfigManagedGroupConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterInstanceGroupConfigManagedGroupConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterInstanceGroupConfigManagedGroupConfigSlice(c *Client, des, nw []ClusterInstanceGroupConfigManagedGroupConfig) []ClusterInstanceGroupConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterInstanceGroupConfigManagedGroupConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterInstanceGroupConfigManagedGroupConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterInstanceGroupConfigAccelerators(des, initial *ClusterInstanceGroupConfigAccelerators, opts ...dcl.ApplyOption) *ClusterInstanceGroupConfigAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterInstanceGroupConfigAccelerators{}

	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}
	if dcl.IsZeroValue(des.AcceleratorCount) {
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}

	return cDes
}

func canonicalizeClusterInstanceGroupConfigAcceleratorsSlice(des, initial []ClusterInstanceGroupConfigAccelerators, opts ...dcl.ApplyOption) []ClusterInstanceGroupConfigAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterInstanceGroupConfigAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterInstanceGroupConfigAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterInstanceGroupConfigAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterInstanceGroupConfigAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterInstanceGroupConfigAccelerators(c *Client, des, nw *ClusterInstanceGroupConfigAccelerators) *ClusterInstanceGroupConfigAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterInstanceGroupConfigAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}

	return nw
}

func canonicalizeNewClusterInstanceGroupConfigAcceleratorsSet(c *Client, des, nw []ClusterInstanceGroupConfigAccelerators) []ClusterInstanceGroupConfigAccelerators {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterInstanceGroupConfigAccelerators
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterInstanceGroupConfigAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterInstanceGroupConfigAcceleratorsSlice(c *Client, des, nw []ClusterInstanceGroupConfigAccelerators) []ClusterInstanceGroupConfigAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterInstanceGroupConfigAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterInstanceGroupConfigAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigSoftwareConfig(des, initial *ClusterClusterConfigSoftwareConfig, opts ...dcl.ApplyOption) *ClusterClusterConfigSoftwareConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigSoftwareConfig{}

	if dcl.StringCanonicalize(des.ImageVersion, initial.ImageVersion) || dcl.IsZeroValue(des.ImageVersion) {
		cDes.ImageVersion = initial.ImageVersion
	} else {
		cDes.ImageVersion = des.ImageVersion
	}
	if dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.IsZeroValue(des.OptionalComponents) {
		cDes.OptionalComponents = initial.OptionalComponents
	} else {
		cDes.OptionalComponents = des.OptionalComponents
	}

	return cDes
}

func canonicalizeClusterClusterConfigSoftwareConfigSlice(des, initial []ClusterClusterConfigSoftwareConfig, opts ...dcl.ApplyOption) []ClusterClusterConfigSoftwareConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigSoftwareConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigSoftwareConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigSoftwareConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigSoftwareConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigSoftwareConfig(c *Client, des, nw *ClusterClusterConfigSoftwareConfig) *ClusterClusterConfigSoftwareConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigSoftwareConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ImageVersion, nw.ImageVersion) {
		nw.ImageVersion = des.ImageVersion
	}

	return nw
}

func canonicalizeNewClusterClusterConfigSoftwareConfigSet(c *Client, des, nw []ClusterClusterConfigSoftwareConfig) []ClusterClusterConfigSoftwareConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigSoftwareConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigSoftwareConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigSoftwareConfigSlice(c *Client, des, nw []ClusterClusterConfigSoftwareConfig) []ClusterClusterConfigSoftwareConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigSoftwareConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigSoftwareConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigInitializationActions(des, initial *ClusterClusterConfigInitializationActions, opts ...dcl.ApplyOption) *ClusterClusterConfigInitializationActions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigInitializationActions{}

	if dcl.StringCanonicalize(des.ExecutableFile, initial.ExecutableFile) || dcl.IsZeroValue(des.ExecutableFile) {
		cDes.ExecutableFile = initial.ExecutableFile
	} else {
		cDes.ExecutableFile = des.ExecutableFile
	}
	if dcl.StringCanonicalize(des.ExecutionTimeout, initial.ExecutionTimeout) || dcl.IsZeroValue(des.ExecutionTimeout) {
		cDes.ExecutionTimeout = initial.ExecutionTimeout
	} else {
		cDes.ExecutionTimeout = des.ExecutionTimeout
	}

	return cDes
}

func canonicalizeClusterClusterConfigInitializationActionsSlice(des, initial []ClusterClusterConfigInitializationActions, opts ...dcl.ApplyOption) []ClusterClusterConfigInitializationActions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigInitializationActions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigInitializationActions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigInitializationActions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigInitializationActions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigInitializationActions(c *Client, des, nw *ClusterClusterConfigInitializationActions) *ClusterClusterConfigInitializationActions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigInitializationActions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ExecutableFile, nw.ExecutableFile) {
		nw.ExecutableFile = des.ExecutableFile
	}
	if dcl.StringCanonicalize(des.ExecutionTimeout, nw.ExecutionTimeout) {
		nw.ExecutionTimeout = des.ExecutionTimeout
	}

	return nw
}

func canonicalizeNewClusterClusterConfigInitializationActionsSet(c *Client, des, nw []ClusterClusterConfigInitializationActions) []ClusterClusterConfigInitializationActions {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigInitializationActions
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigInitializationActionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigInitializationActionsSlice(c *Client, des, nw []ClusterClusterConfigInitializationActions) []ClusterClusterConfigInitializationActions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigInitializationActions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigInitializationActions(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigEncryptionConfig(des, initial *ClusterClusterConfigEncryptionConfig, opts ...dcl.ApplyOption) *ClusterClusterConfigEncryptionConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigEncryptionConfig{}

	if dcl.NameToSelfLink(des.GcePdKmsKeyName, initial.GcePdKmsKeyName) || dcl.IsZeroValue(des.GcePdKmsKeyName) {
		cDes.GcePdKmsKeyName = initial.GcePdKmsKeyName
	} else {
		cDes.GcePdKmsKeyName = des.GcePdKmsKeyName
	}

	return cDes
}

func canonicalizeClusterClusterConfigEncryptionConfigSlice(des, initial []ClusterClusterConfigEncryptionConfig, opts ...dcl.ApplyOption) []ClusterClusterConfigEncryptionConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigEncryptionConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigEncryptionConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigEncryptionConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigEncryptionConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigEncryptionConfig(c *Client, des, nw *ClusterClusterConfigEncryptionConfig) *ClusterClusterConfigEncryptionConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigEncryptionConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.NameToSelfLink(des.GcePdKmsKeyName, nw.GcePdKmsKeyName) {
		nw.GcePdKmsKeyName = des.GcePdKmsKeyName
	}

	return nw
}

func canonicalizeNewClusterClusterConfigEncryptionConfigSet(c *Client, des, nw []ClusterClusterConfigEncryptionConfig) []ClusterClusterConfigEncryptionConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigEncryptionConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigEncryptionConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigEncryptionConfigSlice(c *Client, des, nw []ClusterClusterConfigEncryptionConfig) []ClusterClusterConfigEncryptionConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigEncryptionConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigEncryptionConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigAutoscalingConfig(des, initial *ClusterClusterConfigAutoscalingConfig, opts ...dcl.ApplyOption) *ClusterClusterConfigAutoscalingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigAutoscalingConfig{}

	if dcl.NameToSelfLink(des.Policy, initial.Policy) || dcl.IsZeroValue(des.Policy) {
		cDes.Policy = initial.Policy
	} else {
		cDes.Policy = des.Policy
	}

	return cDes
}

func canonicalizeClusterClusterConfigAutoscalingConfigSlice(des, initial []ClusterClusterConfigAutoscalingConfig, opts ...dcl.ApplyOption) []ClusterClusterConfigAutoscalingConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigAutoscalingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigAutoscalingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigAutoscalingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigAutoscalingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigAutoscalingConfig(c *Client, des, nw *ClusterClusterConfigAutoscalingConfig) *ClusterClusterConfigAutoscalingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigAutoscalingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.NameToSelfLink(des.Policy, nw.Policy) {
		nw.Policy = des.Policy
	}

	return nw
}

func canonicalizeNewClusterClusterConfigAutoscalingConfigSet(c *Client, des, nw []ClusterClusterConfigAutoscalingConfig) []ClusterClusterConfigAutoscalingConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigAutoscalingConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigAutoscalingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigAutoscalingConfigSlice(c *Client, des, nw []ClusterClusterConfigAutoscalingConfig) []ClusterClusterConfigAutoscalingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigAutoscalingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigAutoscalingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigSecurityConfig(des, initial *ClusterClusterConfigSecurityConfig, opts ...dcl.ApplyOption) *ClusterClusterConfigSecurityConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigSecurityConfig{}

	cDes.KerberosConfig = canonicalizeClusterClusterConfigSecurityConfigKerberosConfig(des.KerberosConfig, initial.KerberosConfig, opts...)

	return cDes
}

func canonicalizeClusterClusterConfigSecurityConfigSlice(des, initial []ClusterClusterConfigSecurityConfig, opts ...dcl.ApplyOption) []ClusterClusterConfigSecurityConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigSecurityConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigSecurityConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigSecurityConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigSecurityConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigSecurityConfig(c *Client, des, nw *ClusterClusterConfigSecurityConfig) *ClusterClusterConfigSecurityConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigSecurityConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KerberosConfig = canonicalizeNewClusterClusterConfigSecurityConfigKerberosConfig(c, des.KerberosConfig, nw.KerberosConfig)

	return nw
}

func canonicalizeNewClusterClusterConfigSecurityConfigSet(c *Client, des, nw []ClusterClusterConfigSecurityConfig) []ClusterClusterConfigSecurityConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigSecurityConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigSecurityConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigSecurityConfigSlice(c *Client, des, nw []ClusterClusterConfigSecurityConfig) []ClusterClusterConfigSecurityConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigSecurityConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigSecurityConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigSecurityConfigKerberosConfig(des, initial *ClusterClusterConfigSecurityConfigKerberosConfig, opts ...dcl.ApplyOption) *ClusterClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigSecurityConfigKerberosConfig{}

	if dcl.BoolCanonicalize(des.EnableKerberos, initial.EnableKerberos) || dcl.IsZeroValue(des.EnableKerberos) {
		cDes.EnableKerberos = initial.EnableKerberos
	} else {
		cDes.EnableKerberos = des.EnableKerberos
	}
	if dcl.StringCanonicalize(des.RootPrincipalPassword, initial.RootPrincipalPassword) || dcl.IsZeroValue(des.RootPrincipalPassword) {
		cDes.RootPrincipalPassword = initial.RootPrincipalPassword
	} else {
		cDes.RootPrincipalPassword = des.RootPrincipalPassword
	}
	if dcl.NameToSelfLink(des.KmsKey, initial.KmsKey) || dcl.IsZeroValue(des.KmsKey) {
		cDes.KmsKey = initial.KmsKey
	} else {
		cDes.KmsKey = des.KmsKey
	}
	if dcl.StringCanonicalize(des.Keystore, initial.Keystore) || dcl.IsZeroValue(des.Keystore) {
		cDes.Keystore = initial.Keystore
	} else {
		cDes.Keystore = des.Keystore
	}
	if dcl.StringCanonicalize(des.Truststore, initial.Truststore) || dcl.IsZeroValue(des.Truststore) {
		cDes.Truststore = initial.Truststore
	} else {
		cDes.Truststore = des.Truststore
	}
	if dcl.StringCanonicalize(des.KeystorePassword, initial.KeystorePassword) || dcl.IsZeroValue(des.KeystorePassword) {
		cDes.KeystorePassword = initial.KeystorePassword
	} else {
		cDes.KeystorePassword = des.KeystorePassword
	}
	if dcl.StringCanonicalize(des.KeyPassword, initial.KeyPassword) || dcl.IsZeroValue(des.KeyPassword) {
		cDes.KeyPassword = initial.KeyPassword
	} else {
		cDes.KeyPassword = des.KeyPassword
	}
	if dcl.StringCanonicalize(des.TruststorePassword, initial.TruststorePassword) || dcl.IsZeroValue(des.TruststorePassword) {
		cDes.TruststorePassword = initial.TruststorePassword
	} else {
		cDes.TruststorePassword = des.TruststorePassword
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustRealm, initial.CrossRealmTrustRealm) || dcl.IsZeroValue(des.CrossRealmTrustRealm) {
		cDes.CrossRealmTrustRealm = initial.CrossRealmTrustRealm
	} else {
		cDes.CrossRealmTrustRealm = des.CrossRealmTrustRealm
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustKdc, initial.CrossRealmTrustKdc) || dcl.IsZeroValue(des.CrossRealmTrustKdc) {
		cDes.CrossRealmTrustKdc = initial.CrossRealmTrustKdc
	} else {
		cDes.CrossRealmTrustKdc = des.CrossRealmTrustKdc
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustAdminServer, initial.CrossRealmTrustAdminServer) || dcl.IsZeroValue(des.CrossRealmTrustAdminServer) {
		cDes.CrossRealmTrustAdminServer = initial.CrossRealmTrustAdminServer
	} else {
		cDes.CrossRealmTrustAdminServer = des.CrossRealmTrustAdminServer
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustSharedPassword, initial.CrossRealmTrustSharedPassword) || dcl.IsZeroValue(des.CrossRealmTrustSharedPassword) {
		cDes.CrossRealmTrustSharedPassword = initial.CrossRealmTrustSharedPassword
	} else {
		cDes.CrossRealmTrustSharedPassword = des.CrossRealmTrustSharedPassword
	}
	if dcl.StringCanonicalize(des.KdcDbKey, initial.KdcDbKey) || dcl.IsZeroValue(des.KdcDbKey) {
		cDes.KdcDbKey = initial.KdcDbKey
	} else {
		cDes.KdcDbKey = des.KdcDbKey
	}
	if dcl.IsZeroValue(des.TgtLifetimeHours) {
		cDes.TgtLifetimeHours = initial.TgtLifetimeHours
	} else {
		cDes.TgtLifetimeHours = des.TgtLifetimeHours
	}
	if dcl.StringCanonicalize(des.Realm, initial.Realm) || dcl.IsZeroValue(des.Realm) {
		cDes.Realm = initial.Realm
	} else {
		cDes.Realm = des.Realm
	}

	return cDes
}

func canonicalizeClusterClusterConfigSecurityConfigKerberosConfigSlice(des, initial []ClusterClusterConfigSecurityConfigKerberosConfig, opts ...dcl.ApplyOption) []ClusterClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigSecurityConfigKerberosConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigSecurityConfigKerberosConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigSecurityConfigKerberosConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigSecurityConfigKerberosConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigSecurityConfigKerberosConfig(c *Client, des, nw *ClusterClusterConfigSecurityConfigKerberosConfig) *ClusterClusterConfigSecurityConfigKerberosConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigSecurityConfigKerberosConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.EnableKerberos, nw.EnableKerberos) {
		nw.EnableKerberos = des.EnableKerberos
	}
	if dcl.StringCanonicalize(des.RootPrincipalPassword, nw.RootPrincipalPassword) {
		nw.RootPrincipalPassword = des.RootPrincipalPassword
	}
	if dcl.NameToSelfLink(des.KmsKey, nw.KmsKey) {
		nw.KmsKey = des.KmsKey
	}
	if dcl.StringCanonicalize(des.Keystore, nw.Keystore) {
		nw.Keystore = des.Keystore
	}
	if dcl.StringCanonicalize(des.Truststore, nw.Truststore) {
		nw.Truststore = des.Truststore
	}
	if dcl.StringCanonicalize(des.KeystorePassword, nw.KeystorePassword) {
		nw.KeystorePassword = des.KeystorePassword
	}
	if dcl.StringCanonicalize(des.KeyPassword, nw.KeyPassword) {
		nw.KeyPassword = des.KeyPassword
	}
	if dcl.StringCanonicalize(des.TruststorePassword, nw.TruststorePassword) {
		nw.TruststorePassword = des.TruststorePassword
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustRealm, nw.CrossRealmTrustRealm) {
		nw.CrossRealmTrustRealm = des.CrossRealmTrustRealm
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustKdc, nw.CrossRealmTrustKdc) {
		nw.CrossRealmTrustKdc = des.CrossRealmTrustKdc
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustAdminServer, nw.CrossRealmTrustAdminServer) {
		nw.CrossRealmTrustAdminServer = des.CrossRealmTrustAdminServer
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustSharedPassword, nw.CrossRealmTrustSharedPassword) {
		nw.CrossRealmTrustSharedPassword = des.CrossRealmTrustSharedPassword
	}
	if dcl.StringCanonicalize(des.KdcDbKey, nw.KdcDbKey) {
		nw.KdcDbKey = des.KdcDbKey
	}
	if dcl.StringCanonicalize(des.Realm, nw.Realm) {
		nw.Realm = des.Realm
	}

	return nw
}

func canonicalizeNewClusterClusterConfigSecurityConfigKerberosConfigSet(c *Client, des, nw []ClusterClusterConfigSecurityConfigKerberosConfig) []ClusterClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigSecurityConfigKerberosConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigSecurityConfigKerberosConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigSecurityConfigKerberosConfigSlice(c *Client, des, nw []ClusterClusterConfigSecurityConfigKerberosConfig) []ClusterClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigSecurityConfigKerberosConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigSecurityConfigKerberosConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigLifecycleConfig(des, initial *ClusterClusterConfigLifecycleConfig, opts ...dcl.ApplyOption) *ClusterClusterConfigLifecycleConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigLifecycleConfig{}

	if dcl.StringCanonicalize(des.IdleDeleteTtl, initial.IdleDeleteTtl) || dcl.IsZeroValue(des.IdleDeleteTtl) {
		cDes.IdleDeleteTtl = initial.IdleDeleteTtl
	} else {
		cDes.IdleDeleteTtl = des.IdleDeleteTtl
	}
	if dcl.IsZeroValue(des.AutoDeleteTime) {
		cDes.AutoDeleteTime = initial.AutoDeleteTime
	} else {
		cDes.AutoDeleteTime = des.AutoDeleteTime
	}
	if dcl.StringCanonicalize(des.AutoDeleteTtl, initial.AutoDeleteTtl) || dcl.IsZeroValue(des.AutoDeleteTtl) {
		cDes.AutoDeleteTtl = initial.AutoDeleteTtl
	} else {
		cDes.AutoDeleteTtl = des.AutoDeleteTtl
	}

	return cDes
}

func canonicalizeClusterClusterConfigLifecycleConfigSlice(des, initial []ClusterClusterConfigLifecycleConfig, opts ...dcl.ApplyOption) []ClusterClusterConfigLifecycleConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigLifecycleConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigLifecycleConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigLifecycleConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigLifecycleConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigLifecycleConfig(c *Client, des, nw *ClusterClusterConfigLifecycleConfig) *ClusterClusterConfigLifecycleConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigLifecycleConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.IdleDeleteTtl, nw.IdleDeleteTtl) {
		nw.IdleDeleteTtl = des.IdleDeleteTtl
	}
	if dcl.StringCanonicalize(des.AutoDeleteTtl, nw.AutoDeleteTtl) {
		nw.AutoDeleteTtl = des.AutoDeleteTtl
	}

	return nw
}

func canonicalizeNewClusterClusterConfigLifecycleConfigSet(c *Client, des, nw []ClusterClusterConfigLifecycleConfig) []ClusterClusterConfigLifecycleConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigLifecycleConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigLifecycleConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigLifecycleConfigSlice(c *Client, des, nw []ClusterClusterConfigLifecycleConfig) []ClusterClusterConfigLifecycleConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigLifecycleConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigLifecycleConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterClusterConfigEndpointConfig(des, initial *ClusterClusterConfigEndpointConfig, opts ...dcl.ApplyOption) *ClusterClusterConfigEndpointConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterClusterConfigEndpointConfig{}

	if dcl.BoolCanonicalize(des.EnableHttpPortAccess, initial.EnableHttpPortAccess) || dcl.IsZeroValue(des.EnableHttpPortAccess) {
		cDes.EnableHttpPortAccess = initial.EnableHttpPortAccess
	} else {
		cDes.EnableHttpPortAccess = des.EnableHttpPortAccess
	}

	return cDes
}

func canonicalizeClusterClusterConfigEndpointConfigSlice(des, initial []ClusterClusterConfigEndpointConfig, opts ...dcl.ApplyOption) []ClusterClusterConfigEndpointConfig {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterClusterConfigEndpointConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterClusterConfigEndpointConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterClusterConfigEndpointConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterClusterConfigEndpointConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterClusterConfigEndpointConfig(c *Client, des, nw *ClusterClusterConfigEndpointConfig) *ClusterClusterConfigEndpointConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterClusterConfigEndpointConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.EnableHttpPortAccess, nw.EnableHttpPortAccess) {
		nw.EnableHttpPortAccess = des.EnableHttpPortAccess
	}

	return nw
}

func canonicalizeNewClusterClusterConfigEndpointConfigSet(c *Client, des, nw []ClusterClusterConfigEndpointConfig) []ClusterClusterConfigEndpointConfig {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterClusterConfigEndpointConfig
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterClusterConfigEndpointConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterClusterConfigEndpointConfigSlice(c *Client, des, nw []ClusterClusterConfigEndpointConfig) []ClusterClusterConfigEndpointConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterClusterConfigEndpointConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterClusterConfigEndpointConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterStatus(des, initial *ClusterStatus, opts ...dcl.ApplyOption) *ClusterStatus {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterStatus{}

	return cDes
}

func canonicalizeClusterStatusSlice(des, initial []ClusterStatus, opts ...dcl.ApplyOption) []ClusterStatus {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterStatus, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterStatus(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterStatus, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterStatus(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterStatus(c *Client, des, nw *ClusterStatus) *ClusterStatus {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterStatus while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Detail, nw.Detail) {
		nw.Detail = des.Detail
	}

	return nw
}

func canonicalizeNewClusterStatusSet(c *Client, des, nw []ClusterStatus) []ClusterStatus {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterStatus
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterStatusNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterStatusSlice(c *Client, des, nw []ClusterStatus) []ClusterStatus {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterStatus
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterStatus(c, &d, &n))
	}

	return items
}

func canonicalizeClusterStatusHistory(des, initial *ClusterStatusHistory, opts ...dcl.ApplyOption) *ClusterStatusHistory {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterStatusHistory{}

	return cDes
}

func canonicalizeClusterStatusHistorySlice(des, initial []ClusterStatusHistory, opts ...dcl.ApplyOption) []ClusterStatusHistory {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterStatusHistory, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterStatusHistory(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterStatusHistory, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterStatusHistory(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterStatusHistory(c *Client, des, nw *ClusterStatusHistory) *ClusterStatusHistory {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterStatusHistory while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Detail, nw.Detail) {
		nw.Detail = des.Detail
	}

	return nw
}

func canonicalizeNewClusterStatusHistorySet(c *Client, des, nw []ClusterStatusHistory) []ClusterStatusHistory {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterStatusHistory
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterStatusHistoryNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterStatusHistorySlice(c *Client, des, nw []ClusterStatusHistory) []ClusterStatusHistory {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterStatusHistory
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterStatusHistory(c, &d, &n))
	}

	return items
}

func canonicalizeClusterMetrics(des, initial *ClusterMetrics, opts ...dcl.ApplyOption) *ClusterMetrics {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterMetrics{}

	if dcl.IsZeroValue(des.HdfsMetrics) {
		cDes.HdfsMetrics = initial.HdfsMetrics
	} else {
		cDes.HdfsMetrics = des.HdfsMetrics
	}
	if dcl.IsZeroValue(des.YarnMetrics) {
		cDes.YarnMetrics = initial.YarnMetrics
	} else {
		cDes.YarnMetrics = des.YarnMetrics
	}

	return cDes
}

func canonicalizeClusterMetricsSlice(des, initial []ClusterMetrics, opts ...dcl.ApplyOption) []ClusterMetrics {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterMetrics, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterMetrics(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterMetrics, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterMetrics(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterMetrics(c *Client, des, nw *ClusterMetrics) *ClusterMetrics {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterMetrics while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterMetricsSet(c *Client, des, nw []ClusterMetrics) []ClusterMetrics {
	if des == nil {
		return nw
	}
	var reorderedNew []ClusterMetrics
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareClusterMetricsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewClusterMetricsSlice(c *Client, des, nw []ClusterMetrics) []ClusterMetrics {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterMetrics
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterMetrics(c, &d, &n))
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
func diffCluster(c *Client, desired, actual *Cluster, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ProjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.Info{ObjectFunction: compareClusterClusterConfigNewStyle, EmptyObject: EmptyClusterClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{IgnoredPrefixes: []string{"goog-dataproc-"}, OperationSelector: dcl.TriggersOperation("updateClusterUpdateClusterOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Status, actual.Status, dcl.Info{OutputOnly: true, ObjectFunction: compareClusterStatusNewStyle, EmptyObject: EmptyClusterStatus, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Status")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StatusHistory, actual.StatusHistory, dcl.Info{OutputOnly: true, ObjectFunction: compareClusterStatusHistoryNewStyle, EmptyObject: EmptyClusterStatusHistory, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StatusHistory")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClusterUuid, actual.ClusterUuid, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterUuid")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metrics, actual.Metrics, dcl.Info{OutputOnly: true, ObjectFunction: compareClusterMetricsNewStyle, EmptyObject: EmptyClusterMetrics, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Metrics")); len(ds) != 0 || err != nil {
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

	return newDiffs, nil
}
func compareClusterClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfig or *ClusterClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.StagingBucket, actual.StagingBucket, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConfigBucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TempBucket, actual.TempBucket, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TempBucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GceClusterConfig, actual.GceClusterConfig, dcl.Info{ObjectFunction: compareClusterClusterConfigGceClusterConfigNewStyle, EmptyObject: EmptyClusterClusterConfigGceClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GceClusterConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MasterConfig, actual.MasterConfig, dcl.Info{ObjectFunction: compareClusterInstanceGroupConfigNewStyle, EmptyObject: EmptyClusterInstanceGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MasterConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WorkerConfig, actual.WorkerConfig, dcl.Info{ObjectFunction: compareClusterInstanceGroupConfigNewStyle, EmptyObject: EmptyClusterInstanceGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("WorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecondaryWorkerConfig, actual.SecondaryWorkerConfig, dcl.Info{ObjectFunction: compareClusterInstanceGroupConfigNewStyle, EmptyObject: EmptyClusterInstanceGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SecondaryWorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SoftwareConfig, actual.SoftwareConfig, dcl.Info{ObjectFunction: compareClusterClusterConfigSoftwareConfigNewStyle, EmptyObject: EmptyClusterClusterConfigSoftwareConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SoftwareConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InitializationActions, actual.InitializationActions, dcl.Info{ObjectFunction: compareClusterClusterConfigInitializationActionsNewStyle, EmptyObject: EmptyClusterClusterConfigInitializationActions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InitializationActions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EncryptionConfig, actual.EncryptionConfig, dcl.Info{ObjectFunction: compareClusterClusterConfigEncryptionConfigNewStyle, EmptyObject: EmptyClusterClusterConfigEncryptionConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EncryptionConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoscalingConfig, actual.AutoscalingConfig, dcl.Info{ObjectFunction: compareClusterClusterConfigAutoscalingConfigNewStyle, EmptyObject: EmptyClusterClusterConfigAutoscalingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoscalingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecurityConfig, actual.SecurityConfig, dcl.Info{ObjectFunction: compareClusterClusterConfigSecurityConfigNewStyle, EmptyObject: EmptyClusterClusterConfigSecurityConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SecurityConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LifecycleConfig, actual.LifecycleConfig, dcl.Info{ObjectFunction: compareClusterClusterConfigLifecycleConfigNewStyle, EmptyObject: EmptyClusterClusterConfigLifecycleConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LifecycleConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EndpointConfig, actual.EndpointConfig, dcl.Info{ObjectFunction: compareClusterClusterConfigEndpointConfigNewStyle, EmptyObject: EmptyClusterClusterConfigEndpointConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EndpointConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigGceClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigGceClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigGceClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigGceClusterConfig or *ClusterClusterConfigGceClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigGceClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigGceClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigGceClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Zone, actual.Zone, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ZoneUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Network, actual.Network, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Subnetwork, actual.Subnetwork, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubnetworkUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InternalIPOnly, actual.InternalIPOnly, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InternalIpOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrivateIPv6GoogleAccess, actual.PrivateIPv6GoogleAccess, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrivateIpv6GoogleAccess")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccount, actual.ServiceAccount, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAccount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccountScopes, actual.ServiceAccountScopes, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAccountScopes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tags, actual.Tags, dcl.Info{Type: "Set", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metadata, actual.Metadata, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Metadata")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ReservationAffinity, actual.ReservationAffinity, dcl.Info{ObjectFunction: compareClusterClusterConfigGceClusterConfigReservationAffinityNewStyle, EmptyObject: EmptyClusterClusterConfigGceClusterConfigReservationAffinity, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ReservationAffinity")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NodeGroupAffinity, actual.NodeGroupAffinity, dcl.Info{ObjectFunction: compareClusterClusterConfigGceClusterConfigNodeGroupAffinityNewStyle, EmptyObject: EmptyClusterClusterConfigGceClusterConfigNodeGroupAffinity, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodeGroupAffinity")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigGceClusterConfigReservationAffinityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigGceClusterConfigReservationAffinity)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigGceClusterConfigReservationAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigGceClusterConfigReservationAffinity or *ClusterClusterConfigGceClusterConfigReservationAffinity", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigGceClusterConfigReservationAffinity)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigGceClusterConfigReservationAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigGceClusterConfigReservationAffinity", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ConsumeReservationType, actual.ConsumeReservationType, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConsumeReservationType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Key, actual.Key, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Key")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigGceClusterConfigNodeGroupAffinityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigGceClusterConfigNodeGroupAffinity)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigGceClusterConfigNodeGroupAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigGceClusterConfigNodeGroupAffinity or *ClusterClusterConfigGceClusterConfigNodeGroupAffinity", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigGceClusterConfigNodeGroupAffinity)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigGceClusterConfigNodeGroupAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigGceClusterConfigNodeGroupAffinity", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NodeGroup, actual.NodeGroup, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodeGroupUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterInstanceGroupConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterInstanceGroupConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterInstanceGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterInstanceGroupConfig or *ClusterInstanceGroupConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterInstanceGroupConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterInstanceGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterInstanceGroupConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumInstances, actual.NumInstances, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceNames, actual.InstanceNames, dcl.Info{OutputOnly: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Image, actual.Image, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MachineTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskConfig, actual.DiskConfig, dcl.Info{ObjectFunction: compareClusterInstanceGroupConfigDiskConfigNewStyle, EmptyObject: EmptyClusterInstanceGroupConfigDiskConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IsPreemptible, actual.IsPreemptible, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsPreemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptibility, actual.Preemptibility, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptibility")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ManagedGroupConfig, actual.ManagedGroupConfig, dcl.Info{OutputOnly: true, ObjectFunction: compareClusterInstanceGroupConfigManagedGroupConfigNewStyle, EmptyObject: EmptyClusterInstanceGroupConfigManagedGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedGroupConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Accelerators, actual.Accelerators, dcl.Info{ObjectFunction: compareClusterInstanceGroupConfigAcceleratorsNewStyle, EmptyObject: EmptyClusterInstanceGroupConfigAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Accelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterInstanceGroupConfigDiskConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterInstanceGroupConfigDiskConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterInstanceGroupConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterInstanceGroupConfigDiskConfig or *ClusterInstanceGroupConfigDiskConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterInstanceGroupConfigDiskConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterInstanceGroupConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterInstanceGroupConfigDiskConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BootDiskType, actual.BootDiskType, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BootDiskSizeGb, actual.BootDiskSizeGb, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NumLocalSsds, actual.NumLocalSsds, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumLocalSsds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterInstanceGroupConfigManagedGroupConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterInstanceGroupConfigManagedGroupConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterInstanceGroupConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterInstanceGroupConfigManagedGroupConfig or *ClusterInstanceGroupConfigManagedGroupConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterInstanceGroupConfigManagedGroupConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterInstanceGroupConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterInstanceGroupConfigManagedGroupConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceTemplateName, actual.InstanceTemplateName, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceTemplateName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceGroupManagerName, actual.InstanceGroupManagerName, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceGroupManagerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterInstanceGroupConfigAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterInstanceGroupConfigAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(ClusterInstanceGroupConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterInstanceGroupConfigAccelerators or *ClusterInstanceGroupConfigAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterInstanceGroupConfigAccelerators)
	if !ok {
		actualNotPointer, ok := a.(ClusterInstanceGroupConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterInstanceGroupConfigAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigSoftwareConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigSoftwareConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigSoftwareConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigSoftwareConfig or *ClusterClusterConfigSoftwareConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigSoftwareConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigSoftwareConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigSoftwareConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ImageVersion, actual.ImageVersion, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OptionalComponents, actual.OptionalComponents, dcl.Info{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OptionalComponents")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigInitializationActionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigInitializationActions)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigInitializationActions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigInitializationActions or *ClusterClusterConfigInitializationActions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigInitializationActions)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigInitializationActions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigInitializationActions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ExecutableFile, actual.ExecutableFile, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExecutableFile")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExecutionTimeout, actual.ExecutionTimeout, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExecutionTimeout")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigEncryptionConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigEncryptionConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigEncryptionConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigEncryptionConfig or *ClusterClusterConfigEncryptionConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigEncryptionConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigEncryptionConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigEncryptionConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GcePdKmsKeyName, actual.GcePdKmsKeyName, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GcePdKmsKeyName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigAutoscalingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigAutoscalingConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigAutoscalingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigAutoscalingConfig or *ClusterClusterConfigAutoscalingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigAutoscalingConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigAutoscalingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigAutoscalingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Policy, actual.Policy, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PolicyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigSecurityConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigSecurityConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigSecurityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigSecurityConfig or *ClusterClusterConfigSecurityConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigSecurityConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigSecurityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigSecurityConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KerberosConfig, actual.KerberosConfig, dcl.Info{ObjectFunction: compareClusterClusterConfigSecurityConfigKerberosConfigNewStyle, EmptyObject: EmptyClusterClusterConfigSecurityConfigKerberosConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KerberosConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigSecurityConfigKerberosConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigSecurityConfigKerberosConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigSecurityConfigKerberosConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigSecurityConfigKerberosConfig or *ClusterClusterConfigSecurityConfigKerberosConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigSecurityConfigKerberosConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigSecurityConfigKerberosConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigSecurityConfigKerberosConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.EnableKerberos, actual.EnableKerberos, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableKerberos")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RootPrincipalPassword, actual.RootPrincipalPassword, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RootPrincipalPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KmsKey, actual.KmsKey, dcl.Info{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KmsKeyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Keystore, actual.Keystore, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeystoreUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Truststore, actual.Truststore, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TruststoreUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeystorePassword, actual.KeystorePassword, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeystorePasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyPassword, actual.KeyPassword, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TruststorePassword, actual.TruststorePassword, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TruststorePasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustRealm, actual.CrossRealmTrustRealm, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustRealm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustKdc, actual.CrossRealmTrustKdc, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustKdc")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustAdminServer, actual.CrossRealmTrustAdminServer, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustAdminServer")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustSharedPassword, actual.CrossRealmTrustSharedPassword, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustSharedPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KdcDbKey, actual.KdcDbKey, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KdcDbKeyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TgtLifetimeHours, actual.TgtLifetimeHours, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TgtLifetimeHours")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Realm, actual.Realm, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Realm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigLifecycleConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigLifecycleConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigLifecycleConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigLifecycleConfig or *ClusterClusterConfigLifecycleConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigLifecycleConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigLifecycleConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigLifecycleConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IdleDeleteTtl, actual.IdleDeleteTtl, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IdleDeleteTtl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoDeleteTime, actual.AutoDeleteTime, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoDeleteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoDeleteTtl, actual.AutoDeleteTtl, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoDeleteTtl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IdleStartTime, actual.IdleStartTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IdleStartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterClusterConfigEndpointConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterClusterConfigEndpointConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterClusterConfigEndpointConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigEndpointConfig or *ClusterClusterConfigEndpointConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterClusterConfigEndpointConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterClusterConfigEndpointConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterClusterConfigEndpointConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.HttpPorts, actual.HttpPorts, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HttpPorts")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableHttpPortAccess, actual.EnableHttpPortAccess, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableHttpPortAccess")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterStatusNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterStatus)
	if !ok {
		desiredNotPointer, ok := d.(ClusterStatus)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterStatus or *ClusterStatus", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterStatus)
	if !ok {
		actualNotPointer, ok := a.(ClusterStatus)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterStatus", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.Info{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Detail, actual.Detail, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Detail")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StateStartTime, actual.StateStartTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StateStartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Substate, actual.Substate, dcl.Info{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Substate")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterStatusHistoryNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterStatusHistory)
	if !ok {
		desiredNotPointer, ok := d.(ClusterStatusHistory)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterStatusHistory or *ClusterStatusHistory", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterStatusHistory)
	if !ok {
		actualNotPointer, ok := a.(ClusterStatusHistory)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterStatusHistory", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.Info{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Detail, actual.Detail, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Detail")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StateStartTime, actual.StateStartTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StateStartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Substate, actual.Substate, dcl.Info{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Substate")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterMetricsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterMetrics)
	if !ok {
		desiredNotPointer, ok := d.(ClusterMetrics)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterMetrics or *ClusterMetrics", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterMetrics)
	if !ok {
		actualNotPointer, ok := a.(ClusterMetrics)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterMetrics", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.HdfsMetrics, actual.HdfsMetrics, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HdfsMetrics")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.YarnMetrics, actual.YarnMetrics, dcl.Info{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("YarnMetrics")); len(ds) != 0 || err != nil {
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
func (r *Cluster) urlNormalized() *Cluster {
	normalized := dcl.Copy(*r).(Cluster)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.ClusterUuid = dcl.SelfLinkToName(r.ClusterUuid)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *Cluster) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateCluster" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/regions/{{location}}/clusters/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Cluster resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Cluster) marshal(c *Client) ([]byte, error) {
	m, err := expandCluster(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Cluster: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalCluster decodes JSON responses into the Cluster resource schema.
func unmarshalCluster(b []byte, c *Client) (*Cluster, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapCluster(m, c)
}

func unmarshalMapCluster(m map[string]interface{}, c *Client) (*Cluster, error) {

	flattened := flattenCluster(c, m)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandCluster expands Cluster into a JSON request object.
func expandCluster(c *Client, f *Cluster) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if v, err := expandClusterProject(f, f.Project); err != nil {
		return nil, fmt.Errorf("error expanding Project into projectId: %w", err)
	} else if v != nil {
		m["projectId"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["clusterName"] = v
	}
	if v, err := expandClusterClusterConfig(c, f.Config); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if v != nil {
		m["config"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if v != nil {
		m["location"] = v
	}

	return m, nil
}

// flattenCluster flattens Cluster from a JSON request object into the
// Cluster type.
func flattenCluster(c *Client, i interface{}) *Cluster {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	res := &Cluster{}
	res.Project = dcl.FlattenString(m["projectId"])
	res.Name = dcl.FlattenString(m["clusterName"])
	res.Config = flattenClusterClusterConfig(c, m["config"])
	res.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	res.Status = flattenClusterStatus(c, m["status"])
	res.StatusHistory = flattenClusterStatusHistorySlice(c, m["statusHistory"])
	res.ClusterUuid = dcl.FlattenString(m["clusterUuid"])
	res.Metrics = flattenClusterMetrics(c, m["metrics"])
	res.Location = dcl.FlattenString(m["location"])

	return res
}

// expandClusterClusterConfigMap expands the contents of ClusterClusterConfig into a JSON
// request object.
func expandClusterClusterConfigMap(c *Client, f map[string]ClusterClusterConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigSlice expands the contents of ClusterClusterConfig into a JSON
// request object.
func expandClusterClusterConfigSlice(c *Client, f []ClusterClusterConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigMap flattens the contents of ClusterClusterConfig from a JSON
// response object.
func flattenClusterClusterConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfig{}
	}

	items := make(map[string]ClusterClusterConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigSlice flattens the contents of ClusterClusterConfig from a JSON
// response object.
func flattenClusterClusterConfigSlice(c *Client, i interface{}) []ClusterClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfig{}
	}

	items := make([]ClusterClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfig expands an instance of ClusterClusterConfig into a JSON
// request object.
func expandClusterClusterConfig(c *Client, f *ClusterClusterConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.StagingBucket; !dcl.IsEmptyValueIndirect(v) {
		m["configBucket"] = v
	}
	if v := f.TempBucket; !dcl.IsEmptyValueIndirect(v) {
		m["tempBucket"] = v
	}
	if v, err := expandClusterClusterConfigGceClusterConfig(c, f.GceClusterConfig); err != nil {
		return nil, fmt.Errorf("error expanding GceClusterConfig into gceClusterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gceClusterConfig"] = v
	}
	if v, err := expandClusterInstanceGroupConfig(c, f.MasterConfig); err != nil {
		return nil, fmt.Errorf("error expanding MasterConfig into masterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["masterConfig"] = v
	}
	if v, err := expandClusterInstanceGroupConfig(c, f.WorkerConfig); err != nil {
		return nil, fmt.Errorf("error expanding WorkerConfig into workerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["workerConfig"] = v
	}
	if v, err := expandClusterInstanceGroupConfig(c, f.SecondaryWorkerConfig); err != nil {
		return nil, fmt.Errorf("error expanding SecondaryWorkerConfig into secondaryWorkerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["secondaryWorkerConfig"] = v
	}
	if v, err := expandClusterClusterConfigSoftwareConfig(c, f.SoftwareConfig); err != nil {
		return nil, fmt.Errorf("error expanding SoftwareConfig into softwareConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["softwareConfig"] = v
	}
	if v, err := expandClusterClusterConfigInitializationActionsSlice(c, f.InitializationActions); err != nil {
		return nil, fmt.Errorf("error expanding InitializationActions into initializationActions: %w", err)
	} else if v != nil {
		m["initializationActions"] = v
	}
	if v, err := expandClusterClusterConfigEncryptionConfig(c, f.EncryptionConfig); err != nil {
		return nil, fmt.Errorf("error expanding EncryptionConfig into encryptionConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["encryptionConfig"] = v
	}
	if v, err := expandClusterClusterConfigAutoscalingConfig(c, f.AutoscalingConfig); err != nil {
		return nil, fmt.Errorf("error expanding AutoscalingConfig into autoscalingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["autoscalingConfig"] = v
	}
	if v, err := expandClusterClusterConfigSecurityConfig(c, f.SecurityConfig); err != nil {
		return nil, fmt.Errorf("error expanding SecurityConfig into securityConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["securityConfig"] = v
	}
	if v, err := expandClusterClusterConfigLifecycleConfig(c, f.LifecycleConfig); err != nil {
		return nil, fmt.Errorf("error expanding LifecycleConfig into lifecycleConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["lifecycleConfig"] = v
	}
	if v, err := expandClusterClusterConfigEndpointConfig(c, f.EndpointConfig); err != nil {
		return nil, fmt.Errorf("error expanding EndpointConfig into endpointConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["endpointConfig"] = v
	}

	return m, nil
}

// flattenClusterClusterConfig flattens an instance of ClusterClusterConfig from a JSON
// response object.
func flattenClusterClusterConfig(c *Client, i interface{}) *ClusterClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfig
	}
	r.StagingBucket = dcl.FlattenString(m["configBucket"])
	r.TempBucket = dcl.FlattenString(m["tempBucket"])
	r.GceClusterConfig = flattenClusterClusterConfigGceClusterConfig(c, m["gceClusterConfig"])
	r.MasterConfig = flattenClusterInstanceGroupConfig(c, m["masterConfig"])
	r.WorkerConfig = flattenClusterInstanceGroupConfig(c, m["workerConfig"])
	r.SecondaryWorkerConfig = flattenClusterInstanceGroupConfig(c, m["secondaryWorkerConfig"])
	r.SoftwareConfig = flattenClusterClusterConfigSoftwareConfig(c, m["softwareConfig"])
	r.InitializationActions = flattenClusterClusterConfigInitializationActionsSlice(c, m["initializationActions"])
	r.EncryptionConfig = flattenClusterClusterConfigEncryptionConfig(c, m["encryptionConfig"])
	r.AutoscalingConfig = flattenClusterClusterConfigAutoscalingConfig(c, m["autoscalingConfig"])
	r.SecurityConfig = flattenClusterClusterConfigSecurityConfig(c, m["securityConfig"])
	r.LifecycleConfig = flattenClusterClusterConfigLifecycleConfig(c, m["lifecycleConfig"])
	r.EndpointConfig = flattenClusterClusterConfigEndpointConfig(c, m["endpointConfig"])

	return r
}

// expandClusterClusterConfigGceClusterConfigMap expands the contents of ClusterClusterConfigGceClusterConfig into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfigMap(c *Client, f map[string]ClusterClusterConfigGceClusterConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigGceClusterConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigGceClusterConfigSlice expands the contents of ClusterClusterConfigGceClusterConfig into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfigSlice(c *Client, f []ClusterClusterConfigGceClusterConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigGceClusterConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigGceClusterConfigMap flattens the contents of ClusterClusterConfigGceClusterConfig from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfigGceClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigGceClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigGceClusterConfig{}
	}

	items := make(map[string]ClusterClusterConfigGceClusterConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigGceClusterConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigGceClusterConfigSlice flattens the contents of ClusterClusterConfigGceClusterConfig from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigSlice(c *Client, i interface{}) []ClusterClusterConfigGceClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigGceClusterConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigGceClusterConfig{}
	}

	items := make([]ClusterClusterConfigGceClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigGceClusterConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigGceClusterConfig expands an instance of ClusterClusterConfigGceClusterConfig into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfig(c *Client, f *ClusterClusterConfigGceClusterConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Zone; !dcl.IsEmptyValueIndirect(v) {
		m["zoneUri"] = v
	}
	if v := f.Network; !dcl.IsEmptyValueIndirect(v) {
		m["networkUri"] = v
	}
	if v := f.Subnetwork; !dcl.IsEmptyValueIndirect(v) {
		m["subnetworkUri"] = v
	}
	if v := f.InternalIPOnly; !dcl.IsEmptyValueIndirect(v) {
		m["internalIpOnly"] = v
	}
	if v := f.PrivateIPv6GoogleAccess; !dcl.IsEmptyValueIndirect(v) {
		m["privateIpv6GoogleAccess"] = v
	}
	if v := f.ServiceAccount; !dcl.IsEmptyValueIndirect(v) {
		m["serviceAccount"] = v
	}
	if v := f.ServiceAccountScopes; v != nil {
		m["serviceAccountScopes"] = v
	}
	if v := f.Tags; v != nil {
		m["tags"] = v
	}
	if v := f.Metadata; !dcl.IsEmptyValueIndirect(v) {
		m["metadata"] = v
	}
	if v, err := expandClusterClusterConfigGceClusterConfigReservationAffinity(c, f.ReservationAffinity); err != nil {
		return nil, fmt.Errorf("error expanding ReservationAffinity into reservationAffinity: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["reservationAffinity"] = v
	}
	if v, err := expandClusterClusterConfigGceClusterConfigNodeGroupAffinity(c, f.NodeGroupAffinity); err != nil {
		return nil, fmt.Errorf("error expanding NodeGroupAffinity into nodeGroupAffinity: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["nodeGroupAffinity"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigGceClusterConfig flattens an instance of ClusterClusterConfigGceClusterConfig from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfig(c *Client, i interface{}) *ClusterClusterConfigGceClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigGceClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigGceClusterConfig
	}
	r.Zone = dcl.FlattenString(m["zoneUri"])
	r.Network = dcl.FlattenString(m["networkUri"])
	r.Subnetwork = dcl.FlattenString(m["subnetworkUri"])
	r.InternalIPOnly = dcl.FlattenBool(m["internalIpOnly"])
	r.PrivateIPv6GoogleAccess = flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(m["privateIpv6GoogleAccess"])
	r.ServiceAccount = dcl.FlattenString(m["serviceAccount"])
	r.ServiceAccountScopes = dcl.FlattenStringSlice(m["serviceAccountScopes"])
	r.Tags = dcl.FlattenStringSlice(m["tags"])
	r.Metadata = dcl.FlattenKeyValuePairs(m["metadata"])
	r.ReservationAffinity = flattenClusterClusterConfigGceClusterConfigReservationAffinity(c, m["reservationAffinity"])
	r.NodeGroupAffinity = flattenClusterClusterConfigGceClusterConfigNodeGroupAffinity(c, m["nodeGroupAffinity"])

	return r
}

// expandClusterClusterConfigGceClusterConfigReservationAffinityMap expands the contents of ClusterClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfigReservationAffinityMap(c *Client, f map[string]ClusterClusterConfigGceClusterConfigReservationAffinity) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigGceClusterConfigReservationAffinity(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigGceClusterConfigReservationAffinitySlice expands the contents of ClusterClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, f []ClusterClusterConfigGceClusterConfigReservationAffinity) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigGceClusterConfigReservationAffinity(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigGceClusterConfigReservationAffinityMap flattens the contents of ClusterClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigReservationAffinityMap(c *Client, i interface{}) map[string]ClusterClusterConfigGceClusterConfigReservationAffinity {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigGceClusterConfigReservationAffinity{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigGceClusterConfigReservationAffinity{}
	}

	items := make(map[string]ClusterClusterConfigGceClusterConfigReservationAffinity)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigGceClusterConfigReservationAffinity(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigGceClusterConfigReservationAffinitySlice flattens the contents of ClusterClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, i interface{}) []ClusterClusterConfigGceClusterConfigReservationAffinity {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigGceClusterConfigReservationAffinity{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigGceClusterConfigReservationAffinity{}
	}

	items := make([]ClusterClusterConfigGceClusterConfigReservationAffinity, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigGceClusterConfigReservationAffinity(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigGceClusterConfigReservationAffinity expands an instance of ClusterClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfigReservationAffinity(c *Client, f *ClusterClusterConfigGceClusterConfigReservationAffinity) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ConsumeReservationType; !dcl.IsEmptyValueIndirect(v) {
		m["consumeReservationType"] = v
	}
	if v := f.Key; !dcl.IsEmptyValueIndirect(v) {
		m["key"] = v
	}
	if v := f.Values; v != nil {
		m["values"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigGceClusterConfigReservationAffinity flattens an instance of ClusterClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigReservationAffinity(c *Client, i interface{}) *ClusterClusterConfigGceClusterConfigReservationAffinity {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigGceClusterConfigReservationAffinity{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigGceClusterConfigReservationAffinity
	}
	r.ConsumeReservationType = flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(m["consumeReservationType"])
	r.Key = dcl.FlattenString(m["key"])
	r.Values = dcl.FlattenStringSlice(m["values"])

	return r
}

// expandClusterClusterConfigGceClusterConfigNodeGroupAffinityMap expands the contents of ClusterClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfigNodeGroupAffinityMap(c *Client, f map[string]ClusterClusterConfigGceClusterConfigNodeGroupAffinity) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigGceClusterConfigNodeGroupAffinity(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigGceClusterConfigNodeGroupAffinitySlice expands the contents of ClusterClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, f []ClusterClusterConfigGceClusterConfigNodeGroupAffinity) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigGceClusterConfigNodeGroupAffinity(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigGceClusterConfigNodeGroupAffinityMap flattens the contents of ClusterClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigNodeGroupAffinityMap(c *Client, i interface{}) map[string]ClusterClusterConfigGceClusterConfigNodeGroupAffinity {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	items := make(map[string]ClusterClusterConfigGceClusterConfigNodeGroupAffinity)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigGceClusterConfigNodeGroupAffinity(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigGceClusterConfigNodeGroupAffinitySlice flattens the contents of ClusterClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, i interface{}) []ClusterClusterConfigGceClusterConfigNodeGroupAffinity {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	items := make([]ClusterClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigGceClusterConfigNodeGroupAffinity(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigGceClusterConfigNodeGroupAffinity expands an instance of ClusterClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandClusterClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, f *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NodeGroup; !dcl.IsEmptyValueIndirect(v) {
		m["nodeGroupUri"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigGceClusterConfigNodeGroupAffinity flattens an instance of ClusterClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, i interface{}) *ClusterClusterConfigGceClusterConfigNodeGroupAffinity {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigGceClusterConfigNodeGroupAffinity{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigGceClusterConfigNodeGroupAffinity
	}
	r.NodeGroup = dcl.FlattenString(m["nodeGroupUri"])

	return r
}

// expandClusterInstanceGroupConfigMap expands the contents of ClusterInstanceGroupConfig into a JSON
// request object.
func expandClusterInstanceGroupConfigMap(c *Client, f map[string]ClusterInstanceGroupConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterInstanceGroupConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterInstanceGroupConfigSlice expands the contents of ClusterInstanceGroupConfig into a JSON
// request object.
func expandClusterInstanceGroupConfigSlice(c *Client, f []ClusterInstanceGroupConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterInstanceGroupConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterInstanceGroupConfigMap flattens the contents of ClusterInstanceGroupConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfigMap(c *Client, i interface{}) map[string]ClusterInstanceGroupConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterInstanceGroupConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterInstanceGroupConfig{}
	}

	items := make(map[string]ClusterInstanceGroupConfig)
	for k, item := range a {
		items[k] = *flattenClusterInstanceGroupConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterInstanceGroupConfigSlice flattens the contents of ClusterInstanceGroupConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfigSlice(c *Client, i interface{}) []ClusterInstanceGroupConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterInstanceGroupConfig{}
	}

	if len(a) == 0 {
		return []ClusterInstanceGroupConfig{}
	}

	items := make([]ClusterInstanceGroupConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterInstanceGroupConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterInstanceGroupConfig expands an instance of ClusterInstanceGroupConfig into a JSON
// request object.
func expandClusterInstanceGroupConfig(c *Client, f *ClusterInstanceGroupConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NumInstances; !dcl.IsEmptyValueIndirect(v) {
		m["numInstances"] = v
	}
	if v := f.Image; !dcl.IsEmptyValueIndirect(v) {
		m["imageUri"] = v
	}
	if v := f.MachineType; !dcl.IsEmptyValueIndirect(v) {
		m["machineTypeUri"] = v
	}
	if v, err := expandClusterInstanceGroupConfigDiskConfig(c, f.DiskConfig); err != nil {
		return nil, fmt.Errorf("error expanding DiskConfig into diskConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["diskConfig"] = v
	}
	if v := f.Preemptibility; !dcl.IsEmptyValueIndirect(v) {
		m["preemptibility"] = v
	}
	if v, err := expandClusterInstanceGroupConfigAcceleratorsSlice(c, f.Accelerators); err != nil {
		return nil, fmt.Errorf("error expanding Accelerators into accelerators: %w", err)
	} else if v != nil {
		m["accelerators"] = v
	}
	if v := f.MinCpuPlatform; !dcl.IsEmptyValueIndirect(v) {
		m["minCpuPlatform"] = v
	}

	return m, nil
}

// flattenClusterInstanceGroupConfig flattens an instance of ClusterInstanceGroupConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfig(c *Client, i interface{}) *ClusterInstanceGroupConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterInstanceGroupConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterInstanceGroupConfig
	}
	r.NumInstances = dcl.FlattenInteger(m["numInstances"])
	r.InstanceNames = dcl.FlattenStringSlice(m["instanceNames"])
	r.Image = dcl.FlattenString(m["imageUri"])
	r.MachineType = dcl.FlattenString(m["machineTypeUri"])
	r.DiskConfig = flattenClusterInstanceGroupConfigDiskConfig(c, m["diskConfig"])
	r.IsPreemptible = dcl.FlattenBool(m["isPreemptible"])
	r.Preemptibility = flattenClusterInstanceGroupConfigPreemptibilityEnum(m["preemptibility"])
	r.ManagedGroupConfig = flattenClusterInstanceGroupConfigManagedGroupConfig(c, m["managedGroupConfig"])
	r.Accelerators = flattenClusterInstanceGroupConfigAcceleratorsSlice(c, m["accelerators"])
	r.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])

	return r
}

// expandClusterInstanceGroupConfigDiskConfigMap expands the contents of ClusterInstanceGroupConfigDiskConfig into a JSON
// request object.
func expandClusterInstanceGroupConfigDiskConfigMap(c *Client, f map[string]ClusterInstanceGroupConfigDiskConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterInstanceGroupConfigDiskConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterInstanceGroupConfigDiskConfigSlice expands the contents of ClusterInstanceGroupConfigDiskConfig into a JSON
// request object.
func expandClusterInstanceGroupConfigDiskConfigSlice(c *Client, f []ClusterInstanceGroupConfigDiskConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterInstanceGroupConfigDiskConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterInstanceGroupConfigDiskConfigMap flattens the contents of ClusterInstanceGroupConfigDiskConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfigDiskConfigMap(c *Client, i interface{}) map[string]ClusterInstanceGroupConfigDiskConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterInstanceGroupConfigDiskConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterInstanceGroupConfigDiskConfig{}
	}

	items := make(map[string]ClusterInstanceGroupConfigDiskConfig)
	for k, item := range a {
		items[k] = *flattenClusterInstanceGroupConfigDiskConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterInstanceGroupConfigDiskConfigSlice flattens the contents of ClusterInstanceGroupConfigDiskConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfigDiskConfigSlice(c *Client, i interface{}) []ClusterInstanceGroupConfigDiskConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterInstanceGroupConfigDiskConfig{}
	}

	if len(a) == 0 {
		return []ClusterInstanceGroupConfigDiskConfig{}
	}

	items := make([]ClusterInstanceGroupConfigDiskConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterInstanceGroupConfigDiskConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterInstanceGroupConfigDiskConfig expands an instance of ClusterInstanceGroupConfigDiskConfig into a JSON
// request object.
func expandClusterInstanceGroupConfigDiskConfig(c *Client, f *ClusterInstanceGroupConfigDiskConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.BootDiskType; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskType"] = v
	}
	if v := f.BootDiskSizeGb; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskSizeGb"] = v
	}
	if v := f.NumLocalSsds; !dcl.IsEmptyValueIndirect(v) {
		m["numLocalSsds"] = v
	}

	return m, nil
}

// flattenClusterInstanceGroupConfigDiskConfig flattens an instance of ClusterInstanceGroupConfigDiskConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfigDiskConfig(c *Client, i interface{}) *ClusterInstanceGroupConfigDiskConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterInstanceGroupConfigDiskConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterInstanceGroupConfigDiskConfig
	}
	r.BootDiskType = dcl.FlattenString(m["bootDiskType"])
	r.BootDiskSizeGb = dcl.FlattenInteger(m["bootDiskSizeGb"])
	r.NumLocalSsds = dcl.FlattenInteger(m["numLocalSsds"])

	return r
}

// expandClusterInstanceGroupConfigManagedGroupConfigMap expands the contents of ClusterInstanceGroupConfigManagedGroupConfig into a JSON
// request object.
func expandClusterInstanceGroupConfigManagedGroupConfigMap(c *Client, f map[string]ClusterInstanceGroupConfigManagedGroupConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterInstanceGroupConfigManagedGroupConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterInstanceGroupConfigManagedGroupConfigSlice expands the contents of ClusterInstanceGroupConfigManagedGroupConfig into a JSON
// request object.
func expandClusterInstanceGroupConfigManagedGroupConfigSlice(c *Client, f []ClusterInstanceGroupConfigManagedGroupConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterInstanceGroupConfigManagedGroupConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterInstanceGroupConfigManagedGroupConfigMap flattens the contents of ClusterInstanceGroupConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfigManagedGroupConfigMap(c *Client, i interface{}) map[string]ClusterInstanceGroupConfigManagedGroupConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterInstanceGroupConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterInstanceGroupConfigManagedGroupConfig{}
	}

	items := make(map[string]ClusterInstanceGroupConfigManagedGroupConfig)
	for k, item := range a {
		items[k] = *flattenClusterInstanceGroupConfigManagedGroupConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterInstanceGroupConfigManagedGroupConfigSlice flattens the contents of ClusterInstanceGroupConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfigManagedGroupConfigSlice(c *Client, i interface{}) []ClusterInstanceGroupConfigManagedGroupConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterInstanceGroupConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return []ClusterInstanceGroupConfigManagedGroupConfig{}
	}

	items := make([]ClusterInstanceGroupConfigManagedGroupConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterInstanceGroupConfigManagedGroupConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterInstanceGroupConfigManagedGroupConfig expands an instance of ClusterInstanceGroupConfigManagedGroupConfig into a JSON
// request object.
func expandClusterInstanceGroupConfigManagedGroupConfig(c *Client, f *ClusterInstanceGroupConfigManagedGroupConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenClusterInstanceGroupConfigManagedGroupConfig flattens an instance of ClusterInstanceGroupConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterInstanceGroupConfigManagedGroupConfig(c *Client, i interface{}) *ClusterInstanceGroupConfigManagedGroupConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterInstanceGroupConfigManagedGroupConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterInstanceGroupConfigManagedGroupConfig
	}
	r.InstanceTemplateName = dcl.FlattenString(m["instanceTemplateName"])
	r.InstanceGroupManagerName = dcl.FlattenString(m["instanceGroupManagerName"])

	return r
}

// expandClusterInstanceGroupConfigAcceleratorsMap expands the contents of ClusterInstanceGroupConfigAccelerators into a JSON
// request object.
func expandClusterInstanceGroupConfigAcceleratorsMap(c *Client, f map[string]ClusterInstanceGroupConfigAccelerators) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterInstanceGroupConfigAccelerators(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterInstanceGroupConfigAcceleratorsSlice expands the contents of ClusterInstanceGroupConfigAccelerators into a JSON
// request object.
func expandClusterInstanceGroupConfigAcceleratorsSlice(c *Client, f []ClusterInstanceGroupConfigAccelerators) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterInstanceGroupConfigAccelerators(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterInstanceGroupConfigAcceleratorsMap flattens the contents of ClusterInstanceGroupConfigAccelerators from a JSON
// response object.
func flattenClusterInstanceGroupConfigAcceleratorsMap(c *Client, i interface{}) map[string]ClusterInstanceGroupConfigAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterInstanceGroupConfigAccelerators{}
	}

	if len(a) == 0 {
		return map[string]ClusterInstanceGroupConfigAccelerators{}
	}

	items := make(map[string]ClusterInstanceGroupConfigAccelerators)
	for k, item := range a {
		items[k] = *flattenClusterInstanceGroupConfigAccelerators(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterInstanceGroupConfigAcceleratorsSlice flattens the contents of ClusterInstanceGroupConfigAccelerators from a JSON
// response object.
func flattenClusterInstanceGroupConfigAcceleratorsSlice(c *Client, i interface{}) []ClusterInstanceGroupConfigAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterInstanceGroupConfigAccelerators{}
	}

	if len(a) == 0 {
		return []ClusterInstanceGroupConfigAccelerators{}
	}

	items := make([]ClusterInstanceGroupConfigAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterInstanceGroupConfigAccelerators(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterInstanceGroupConfigAccelerators expands an instance of ClusterInstanceGroupConfigAccelerators into a JSON
// request object.
func expandClusterInstanceGroupConfigAccelerators(c *Client, f *ClusterInstanceGroupConfigAccelerators) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AcceleratorType; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorTypeUri"] = v
	}
	if v := f.AcceleratorCount; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorCount"] = v
	}

	return m, nil
}

// flattenClusterInstanceGroupConfigAccelerators flattens an instance of ClusterInstanceGroupConfigAccelerators from a JSON
// response object.
func flattenClusterInstanceGroupConfigAccelerators(c *Client, i interface{}) *ClusterInstanceGroupConfigAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterInstanceGroupConfigAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterInstanceGroupConfigAccelerators
	}
	r.AcceleratorType = dcl.FlattenString(m["acceleratorTypeUri"])
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])

	return r
}

// expandClusterClusterConfigSoftwareConfigMap expands the contents of ClusterClusterConfigSoftwareConfig into a JSON
// request object.
func expandClusterClusterConfigSoftwareConfigMap(c *Client, f map[string]ClusterClusterConfigSoftwareConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigSoftwareConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigSoftwareConfigSlice expands the contents of ClusterClusterConfigSoftwareConfig into a JSON
// request object.
func expandClusterClusterConfigSoftwareConfigSlice(c *Client, f []ClusterClusterConfigSoftwareConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigSoftwareConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigSoftwareConfigMap flattens the contents of ClusterClusterConfigSoftwareConfig from a JSON
// response object.
func flattenClusterClusterConfigSoftwareConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfigSoftwareConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigSoftwareConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigSoftwareConfig{}
	}

	items := make(map[string]ClusterClusterConfigSoftwareConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigSoftwareConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigSoftwareConfigSlice flattens the contents of ClusterClusterConfigSoftwareConfig from a JSON
// response object.
func flattenClusterClusterConfigSoftwareConfigSlice(c *Client, i interface{}) []ClusterClusterConfigSoftwareConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigSoftwareConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigSoftwareConfig{}
	}

	items := make([]ClusterClusterConfigSoftwareConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigSoftwareConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigSoftwareConfig expands an instance of ClusterClusterConfigSoftwareConfig into a JSON
// request object.
func expandClusterClusterConfigSoftwareConfig(c *Client, f *ClusterClusterConfigSoftwareConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ImageVersion; !dcl.IsEmptyValueIndirect(v) {
		m["imageVersion"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v := f.OptionalComponents; v != nil {
		m["optionalComponents"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigSoftwareConfig flattens an instance of ClusterClusterConfigSoftwareConfig from a JSON
// response object.
func flattenClusterClusterConfigSoftwareConfig(c *Client, i interface{}) *ClusterClusterConfigSoftwareConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigSoftwareConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigSoftwareConfig
	}
	r.ImageVersion = dcl.FlattenString(m["imageVersion"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.OptionalComponents = flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnumSlice(c, m["optionalComponents"])

	return r
}

// expandClusterClusterConfigInitializationActionsMap expands the contents of ClusterClusterConfigInitializationActions into a JSON
// request object.
func expandClusterClusterConfigInitializationActionsMap(c *Client, f map[string]ClusterClusterConfigInitializationActions) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigInitializationActions(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigInitializationActionsSlice expands the contents of ClusterClusterConfigInitializationActions into a JSON
// request object.
func expandClusterClusterConfigInitializationActionsSlice(c *Client, f []ClusterClusterConfigInitializationActions) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigInitializationActions(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigInitializationActionsMap flattens the contents of ClusterClusterConfigInitializationActions from a JSON
// response object.
func flattenClusterClusterConfigInitializationActionsMap(c *Client, i interface{}) map[string]ClusterClusterConfigInitializationActions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigInitializationActions{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigInitializationActions{}
	}

	items := make(map[string]ClusterClusterConfigInitializationActions)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigInitializationActions(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigInitializationActionsSlice flattens the contents of ClusterClusterConfigInitializationActions from a JSON
// response object.
func flattenClusterClusterConfigInitializationActionsSlice(c *Client, i interface{}) []ClusterClusterConfigInitializationActions {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigInitializationActions{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigInitializationActions{}
	}

	items := make([]ClusterClusterConfigInitializationActions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigInitializationActions(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigInitializationActions expands an instance of ClusterClusterConfigInitializationActions into a JSON
// request object.
func expandClusterClusterConfigInitializationActions(c *Client, f *ClusterClusterConfigInitializationActions) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ExecutableFile; !dcl.IsEmptyValueIndirect(v) {
		m["executableFile"] = v
	}
	if v := f.ExecutionTimeout; !dcl.IsEmptyValueIndirect(v) {
		m["executionTimeout"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigInitializationActions flattens an instance of ClusterClusterConfigInitializationActions from a JSON
// response object.
func flattenClusterClusterConfigInitializationActions(c *Client, i interface{}) *ClusterClusterConfigInitializationActions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigInitializationActions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigInitializationActions
	}
	r.ExecutableFile = dcl.FlattenString(m["executableFile"])
	r.ExecutionTimeout = dcl.FlattenString(m["executionTimeout"])

	return r
}

// expandClusterClusterConfigEncryptionConfigMap expands the contents of ClusterClusterConfigEncryptionConfig into a JSON
// request object.
func expandClusterClusterConfigEncryptionConfigMap(c *Client, f map[string]ClusterClusterConfigEncryptionConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigEncryptionConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigEncryptionConfigSlice expands the contents of ClusterClusterConfigEncryptionConfig into a JSON
// request object.
func expandClusterClusterConfigEncryptionConfigSlice(c *Client, f []ClusterClusterConfigEncryptionConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigEncryptionConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigEncryptionConfigMap flattens the contents of ClusterClusterConfigEncryptionConfig from a JSON
// response object.
func flattenClusterClusterConfigEncryptionConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfigEncryptionConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigEncryptionConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigEncryptionConfig{}
	}

	items := make(map[string]ClusterClusterConfigEncryptionConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigEncryptionConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigEncryptionConfigSlice flattens the contents of ClusterClusterConfigEncryptionConfig from a JSON
// response object.
func flattenClusterClusterConfigEncryptionConfigSlice(c *Client, i interface{}) []ClusterClusterConfigEncryptionConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigEncryptionConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigEncryptionConfig{}
	}

	items := make([]ClusterClusterConfigEncryptionConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigEncryptionConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigEncryptionConfig expands an instance of ClusterClusterConfigEncryptionConfig into a JSON
// request object.
func expandClusterClusterConfigEncryptionConfig(c *Client, f *ClusterClusterConfigEncryptionConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GcePdKmsKeyName; !dcl.IsEmptyValueIndirect(v) {
		m["gcePdKmsKeyName"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigEncryptionConfig flattens an instance of ClusterClusterConfigEncryptionConfig from a JSON
// response object.
func flattenClusterClusterConfigEncryptionConfig(c *Client, i interface{}) *ClusterClusterConfigEncryptionConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigEncryptionConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigEncryptionConfig
	}
	r.GcePdKmsKeyName = dcl.FlattenString(m["gcePdKmsKeyName"])

	return r
}

// expandClusterClusterConfigAutoscalingConfigMap expands the contents of ClusterClusterConfigAutoscalingConfig into a JSON
// request object.
func expandClusterClusterConfigAutoscalingConfigMap(c *Client, f map[string]ClusterClusterConfigAutoscalingConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigAutoscalingConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigAutoscalingConfigSlice expands the contents of ClusterClusterConfigAutoscalingConfig into a JSON
// request object.
func expandClusterClusterConfigAutoscalingConfigSlice(c *Client, f []ClusterClusterConfigAutoscalingConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigAutoscalingConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigAutoscalingConfigMap flattens the contents of ClusterClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenClusterClusterConfigAutoscalingConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfigAutoscalingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigAutoscalingConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigAutoscalingConfig{}
	}

	items := make(map[string]ClusterClusterConfigAutoscalingConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigAutoscalingConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigAutoscalingConfigSlice flattens the contents of ClusterClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenClusterClusterConfigAutoscalingConfigSlice(c *Client, i interface{}) []ClusterClusterConfigAutoscalingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigAutoscalingConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigAutoscalingConfig{}
	}

	items := make([]ClusterClusterConfigAutoscalingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigAutoscalingConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigAutoscalingConfig expands an instance of ClusterClusterConfigAutoscalingConfig into a JSON
// request object.
func expandClusterClusterConfigAutoscalingConfig(c *Client, f *ClusterClusterConfigAutoscalingConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Policy; !dcl.IsEmptyValueIndirect(v) {
		m["policyUri"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigAutoscalingConfig flattens an instance of ClusterClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenClusterClusterConfigAutoscalingConfig(c *Client, i interface{}) *ClusterClusterConfigAutoscalingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigAutoscalingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigAutoscalingConfig
	}
	r.Policy = dcl.FlattenString(m["policyUri"])

	return r
}

// expandClusterClusterConfigSecurityConfigMap expands the contents of ClusterClusterConfigSecurityConfig into a JSON
// request object.
func expandClusterClusterConfigSecurityConfigMap(c *Client, f map[string]ClusterClusterConfigSecurityConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigSecurityConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigSecurityConfigSlice expands the contents of ClusterClusterConfigSecurityConfig into a JSON
// request object.
func expandClusterClusterConfigSecurityConfigSlice(c *Client, f []ClusterClusterConfigSecurityConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigSecurityConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigSecurityConfigMap flattens the contents of ClusterClusterConfigSecurityConfig from a JSON
// response object.
func flattenClusterClusterConfigSecurityConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfigSecurityConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigSecurityConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigSecurityConfig{}
	}

	items := make(map[string]ClusterClusterConfigSecurityConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigSecurityConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigSecurityConfigSlice flattens the contents of ClusterClusterConfigSecurityConfig from a JSON
// response object.
func flattenClusterClusterConfigSecurityConfigSlice(c *Client, i interface{}) []ClusterClusterConfigSecurityConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigSecurityConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigSecurityConfig{}
	}

	items := make([]ClusterClusterConfigSecurityConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigSecurityConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigSecurityConfig expands an instance of ClusterClusterConfigSecurityConfig into a JSON
// request object.
func expandClusterClusterConfigSecurityConfig(c *Client, f *ClusterClusterConfigSecurityConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandClusterClusterConfigSecurityConfigKerberosConfig(c, f.KerberosConfig); err != nil {
		return nil, fmt.Errorf("error expanding KerberosConfig into kerberosConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["kerberosConfig"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigSecurityConfig flattens an instance of ClusterClusterConfigSecurityConfig from a JSON
// response object.
func flattenClusterClusterConfigSecurityConfig(c *Client, i interface{}) *ClusterClusterConfigSecurityConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigSecurityConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigSecurityConfig
	}
	r.KerberosConfig = flattenClusterClusterConfigSecurityConfigKerberosConfig(c, m["kerberosConfig"])

	return r
}

// expandClusterClusterConfigSecurityConfigKerberosConfigMap expands the contents of ClusterClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandClusterClusterConfigSecurityConfigKerberosConfigMap(c *Client, f map[string]ClusterClusterConfigSecurityConfigKerberosConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigSecurityConfigKerberosConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigSecurityConfigKerberosConfigSlice expands the contents of ClusterClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandClusterClusterConfigSecurityConfigKerberosConfigSlice(c *Client, f []ClusterClusterConfigSecurityConfigKerberosConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigSecurityConfigKerberosConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigSecurityConfigKerberosConfigMap flattens the contents of ClusterClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenClusterClusterConfigSecurityConfigKerberosConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfigSecurityConfigKerberosConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigSecurityConfigKerberosConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigSecurityConfigKerberosConfig{}
	}

	items := make(map[string]ClusterClusterConfigSecurityConfigKerberosConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigSecurityConfigKerberosConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigSecurityConfigKerberosConfigSlice flattens the contents of ClusterClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenClusterClusterConfigSecurityConfigKerberosConfigSlice(c *Client, i interface{}) []ClusterClusterConfigSecurityConfigKerberosConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigSecurityConfigKerberosConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigSecurityConfigKerberosConfig{}
	}

	items := make([]ClusterClusterConfigSecurityConfigKerberosConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigSecurityConfigKerberosConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigSecurityConfigKerberosConfig expands an instance of ClusterClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandClusterClusterConfigSecurityConfigKerberosConfig(c *Client, f *ClusterClusterConfigSecurityConfigKerberosConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.EnableKerberos; !dcl.IsEmptyValueIndirect(v) {
		m["enableKerberos"] = v
	}
	if v := f.RootPrincipalPassword; !dcl.IsEmptyValueIndirect(v) {
		m["rootPrincipalPasswordUri"] = v
	}
	if v := f.KmsKey; !dcl.IsEmptyValueIndirect(v) {
		m["kmsKeyUri"] = v
	}
	if v := f.Keystore; !dcl.IsEmptyValueIndirect(v) {
		m["keystoreUri"] = v
	}
	if v := f.Truststore; !dcl.IsEmptyValueIndirect(v) {
		m["truststoreUri"] = v
	}
	if v := f.KeystorePassword; !dcl.IsEmptyValueIndirect(v) {
		m["keystorePasswordUri"] = v
	}
	if v := f.KeyPassword; !dcl.IsEmptyValueIndirect(v) {
		m["keyPasswordUri"] = v
	}
	if v := f.TruststorePassword; !dcl.IsEmptyValueIndirect(v) {
		m["truststorePasswordUri"] = v
	}
	if v := f.CrossRealmTrustRealm; !dcl.IsEmptyValueIndirect(v) {
		m["crossRealmTrustRealm"] = v
	}
	if v := f.CrossRealmTrustKdc; !dcl.IsEmptyValueIndirect(v) {
		m["crossRealmTrustKdc"] = v
	}
	if v := f.CrossRealmTrustAdminServer; !dcl.IsEmptyValueIndirect(v) {
		m["crossRealmTrustAdminServer"] = v
	}
	if v := f.CrossRealmTrustSharedPassword; !dcl.IsEmptyValueIndirect(v) {
		m["crossRealmTrustSharedPasswordUri"] = v
	}
	if v := f.KdcDbKey; !dcl.IsEmptyValueIndirect(v) {
		m["kdcDbKeyUri"] = v
	}
	if v := f.TgtLifetimeHours; !dcl.IsEmptyValueIndirect(v) {
		m["tgtLifetimeHours"] = v
	}
	if v := f.Realm; !dcl.IsEmptyValueIndirect(v) {
		m["realm"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigSecurityConfigKerberosConfig flattens an instance of ClusterClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenClusterClusterConfigSecurityConfigKerberosConfig(c *Client, i interface{}) *ClusterClusterConfigSecurityConfigKerberosConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigSecurityConfigKerberosConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigSecurityConfigKerberosConfig
	}
	r.EnableKerberos = dcl.FlattenBool(m["enableKerberos"])
	r.RootPrincipalPassword = dcl.FlattenString(m["rootPrincipalPasswordUri"])
	r.KmsKey = dcl.FlattenString(m["kmsKeyUri"])
	r.Keystore = dcl.FlattenString(m["keystoreUri"])
	r.Truststore = dcl.FlattenString(m["truststoreUri"])
	r.KeystorePassword = dcl.FlattenString(m["keystorePasswordUri"])
	r.KeyPassword = dcl.FlattenString(m["keyPasswordUri"])
	r.TruststorePassword = dcl.FlattenString(m["truststorePasswordUri"])
	r.CrossRealmTrustRealm = dcl.FlattenString(m["crossRealmTrustRealm"])
	r.CrossRealmTrustKdc = dcl.FlattenString(m["crossRealmTrustKdc"])
	r.CrossRealmTrustAdminServer = dcl.FlattenString(m["crossRealmTrustAdminServer"])
	r.CrossRealmTrustSharedPassword = dcl.FlattenString(m["crossRealmTrustSharedPasswordUri"])
	r.KdcDbKey = dcl.FlattenString(m["kdcDbKeyUri"])
	r.TgtLifetimeHours = dcl.FlattenInteger(m["tgtLifetimeHours"])
	r.Realm = dcl.FlattenString(m["realm"])

	return r
}

// expandClusterClusterConfigLifecycleConfigMap expands the contents of ClusterClusterConfigLifecycleConfig into a JSON
// request object.
func expandClusterClusterConfigLifecycleConfigMap(c *Client, f map[string]ClusterClusterConfigLifecycleConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigLifecycleConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigLifecycleConfigSlice expands the contents of ClusterClusterConfigLifecycleConfig into a JSON
// request object.
func expandClusterClusterConfigLifecycleConfigSlice(c *Client, f []ClusterClusterConfigLifecycleConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigLifecycleConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigLifecycleConfigMap flattens the contents of ClusterClusterConfigLifecycleConfig from a JSON
// response object.
func flattenClusterClusterConfigLifecycleConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfigLifecycleConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigLifecycleConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigLifecycleConfig{}
	}

	items := make(map[string]ClusterClusterConfigLifecycleConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigLifecycleConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigLifecycleConfigSlice flattens the contents of ClusterClusterConfigLifecycleConfig from a JSON
// response object.
func flattenClusterClusterConfigLifecycleConfigSlice(c *Client, i interface{}) []ClusterClusterConfigLifecycleConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigLifecycleConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigLifecycleConfig{}
	}

	items := make([]ClusterClusterConfigLifecycleConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigLifecycleConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigLifecycleConfig expands an instance of ClusterClusterConfigLifecycleConfig into a JSON
// request object.
func expandClusterClusterConfigLifecycleConfig(c *Client, f *ClusterClusterConfigLifecycleConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.IdleDeleteTtl; !dcl.IsEmptyValueIndirect(v) {
		m["idleDeleteTtl"] = v
	}
	if v := f.AutoDeleteTime; !dcl.IsEmptyValueIndirect(v) {
		m["autoDeleteTime"] = v
	}
	if v := f.AutoDeleteTtl; !dcl.IsEmptyValueIndirect(v) {
		m["autoDeleteTtl"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigLifecycleConfig flattens an instance of ClusterClusterConfigLifecycleConfig from a JSON
// response object.
func flattenClusterClusterConfigLifecycleConfig(c *Client, i interface{}) *ClusterClusterConfigLifecycleConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigLifecycleConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigLifecycleConfig
	}
	r.IdleDeleteTtl = dcl.FlattenString(m["idleDeleteTtl"])
	r.AutoDeleteTime = dcl.FlattenString(m["autoDeleteTime"])
	r.AutoDeleteTtl = dcl.FlattenString(m["autoDeleteTtl"])
	r.IdleStartTime = dcl.FlattenString(m["idleStartTime"])

	return r
}

// expandClusterClusterConfigEndpointConfigMap expands the contents of ClusterClusterConfigEndpointConfig into a JSON
// request object.
func expandClusterClusterConfigEndpointConfigMap(c *Client, f map[string]ClusterClusterConfigEndpointConfig) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterClusterConfigEndpointConfig(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterClusterConfigEndpointConfigSlice expands the contents of ClusterClusterConfigEndpointConfig into a JSON
// request object.
func expandClusterClusterConfigEndpointConfigSlice(c *Client, f []ClusterClusterConfigEndpointConfig) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterClusterConfigEndpointConfig(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterClusterConfigEndpointConfigMap flattens the contents of ClusterClusterConfigEndpointConfig from a JSON
// response object.
func flattenClusterClusterConfigEndpointConfigMap(c *Client, i interface{}) map[string]ClusterClusterConfigEndpointConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigEndpointConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigEndpointConfig{}
	}

	items := make(map[string]ClusterClusterConfigEndpointConfig)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigEndpointConfig(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterClusterConfigEndpointConfigSlice flattens the contents of ClusterClusterConfigEndpointConfig from a JSON
// response object.
func flattenClusterClusterConfigEndpointConfigSlice(c *Client, i interface{}) []ClusterClusterConfigEndpointConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigEndpointConfig{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigEndpointConfig{}
	}

	items := make([]ClusterClusterConfigEndpointConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigEndpointConfig(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterClusterConfigEndpointConfig expands an instance of ClusterClusterConfigEndpointConfig into a JSON
// request object.
func expandClusterClusterConfigEndpointConfig(c *Client, f *ClusterClusterConfigEndpointConfig) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.EnableHttpPortAccess; !dcl.IsEmptyValueIndirect(v) {
		m["enableHttpPortAccess"] = v
	}

	return m, nil
}

// flattenClusterClusterConfigEndpointConfig flattens an instance of ClusterClusterConfigEndpointConfig from a JSON
// response object.
func flattenClusterClusterConfigEndpointConfig(c *Client, i interface{}) *ClusterClusterConfigEndpointConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterClusterConfigEndpointConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterClusterConfigEndpointConfig
	}
	r.HttpPorts = dcl.FlattenKeyValuePairs(m["httpPorts"])
	r.EnableHttpPortAccess = dcl.FlattenBool(m["enableHttpPortAccess"])

	return r
}

// expandClusterStatusMap expands the contents of ClusterStatus into a JSON
// request object.
func expandClusterStatusMap(c *Client, f map[string]ClusterStatus) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterStatus(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterStatusSlice expands the contents of ClusterStatus into a JSON
// request object.
func expandClusterStatusSlice(c *Client, f []ClusterStatus) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterStatus(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterStatusMap flattens the contents of ClusterStatus from a JSON
// response object.
func flattenClusterStatusMap(c *Client, i interface{}) map[string]ClusterStatus {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStatus{}
	}

	if len(a) == 0 {
		return map[string]ClusterStatus{}
	}

	items := make(map[string]ClusterStatus)
	for k, item := range a {
		items[k] = *flattenClusterStatus(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterStatusSlice flattens the contents of ClusterStatus from a JSON
// response object.
func flattenClusterStatusSlice(c *Client, i interface{}) []ClusterStatus {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStatus{}
	}

	if len(a) == 0 {
		return []ClusterStatus{}
	}

	items := make([]ClusterStatus, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStatus(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterStatus expands an instance of ClusterStatus into a JSON
// request object.
func expandClusterStatus(c *Client, f *ClusterStatus) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenClusterStatus flattens an instance of ClusterStatus from a JSON
// response object.
func flattenClusterStatus(c *Client, i interface{}) *ClusterStatus {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterStatus{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterStatus
	}
	r.State = flattenClusterStatusStateEnum(m["state"])
	r.Detail = dcl.FlattenString(m["detail"])
	r.StateStartTime = dcl.FlattenString(m["stateStartTime"])
	r.Substate = flattenClusterStatusSubstateEnum(m["substate"])

	return r
}

// expandClusterStatusHistoryMap expands the contents of ClusterStatusHistory into a JSON
// request object.
func expandClusterStatusHistoryMap(c *Client, f map[string]ClusterStatusHistory) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterStatusHistory(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterStatusHistorySlice expands the contents of ClusterStatusHistory into a JSON
// request object.
func expandClusterStatusHistorySlice(c *Client, f []ClusterStatusHistory) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterStatusHistory(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterStatusHistoryMap flattens the contents of ClusterStatusHistory from a JSON
// response object.
func flattenClusterStatusHistoryMap(c *Client, i interface{}) map[string]ClusterStatusHistory {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStatusHistory{}
	}

	if len(a) == 0 {
		return map[string]ClusterStatusHistory{}
	}

	items := make(map[string]ClusterStatusHistory)
	for k, item := range a {
		items[k] = *flattenClusterStatusHistory(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterStatusHistorySlice flattens the contents of ClusterStatusHistory from a JSON
// response object.
func flattenClusterStatusHistorySlice(c *Client, i interface{}) []ClusterStatusHistory {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStatusHistory{}
	}

	if len(a) == 0 {
		return []ClusterStatusHistory{}
	}

	items := make([]ClusterStatusHistory, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStatusHistory(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterStatusHistory expands an instance of ClusterStatusHistory into a JSON
// request object.
func expandClusterStatusHistory(c *Client, f *ClusterStatusHistory) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenClusterStatusHistory flattens an instance of ClusterStatusHistory from a JSON
// response object.
func flattenClusterStatusHistory(c *Client, i interface{}) *ClusterStatusHistory {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterStatusHistory{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterStatusHistory
	}
	r.State = flattenClusterStatusHistoryStateEnum(m["state"])
	r.Detail = dcl.FlattenString(m["detail"])
	r.StateStartTime = dcl.FlattenString(m["stateStartTime"])
	r.Substate = flattenClusterStatusHistorySubstateEnum(m["substate"])

	return r
}

// expandClusterMetricsMap expands the contents of ClusterMetrics into a JSON
// request object.
func expandClusterMetricsMap(c *Client, f map[string]ClusterMetrics) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterMetrics(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterMetricsSlice expands the contents of ClusterMetrics into a JSON
// request object.
func expandClusterMetricsSlice(c *Client, f []ClusterMetrics) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterMetrics(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterMetricsMap flattens the contents of ClusterMetrics from a JSON
// response object.
func flattenClusterMetricsMap(c *Client, i interface{}) map[string]ClusterMetrics {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterMetrics{}
	}

	if len(a) == 0 {
		return map[string]ClusterMetrics{}
	}

	items := make(map[string]ClusterMetrics)
	for k, item := range a {
		items[k] = *flattenClusterMetrics(c, item.(map[string]interface{}))
	}

	return items
}

// flattenClusterMetricsSlice flattens the contents of ClusterMetrics from a JSON
// response object.
func flattenClusterMetricsSlice(c *Client, i interface{}) []ClusterMetrics {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterMetrics{}
	}

	if len(a) == 0 {
		return []ClusterMetrics{}
	}

	items := make([]ClusterMetrics, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterMetrics(c, item.(map[string]interface{})))
	}

	return items
}

// expandClusterMetrics expands an instance of ClusterMetrics into a JSON
// request object.
func expandClusterMetrics(c *Client, f *ClusterMetrics) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.HdfsMetrics; !dcl.IsEmptyValueIndirect(v) {
		m["hdfsMetrics"] = v
	}
	if v := f.YarnMetrics; !dcl.IsEmptyValueIndirect(v) {
		m["yarnMetrics"] = v
	}

	return m, nil
}

// flattenClusterMetrics flattens an instance of ClusterMetrics from a JSON
// response object.
func flattenClusterMetrics(c *Client, i interface{}) *ClusterMetrics {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterMetrics{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterMetrics
	}
	r.HdfsMetrics = dcl.FlattenKeyValuePairs(m["hdfsMetrics"])
	r.YarnMetrics = dcl.FlattenKeyValuePairs(m["yarnMetrics"])

	return r
}

// flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumMap flattens the contents of ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumMap(c *Client, i interface{}) map[string]ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	items := make(map[string]ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(item.(interface{}))
	}

	return items
}

// flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumSlice flattens the contents of ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumSlice(c *Client, i interface{}) []ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	items := make([]ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(item.(interface{})))
	}

	return items
}

// flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum with the same value as that string.
func flattenClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(i interface{}) *ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	s, ok := i.(string)
	if !ok {
		return ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef("")
	}

	return ClusterClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef(s)
}

// flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumMap flattens the contents of ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumMap(c *Client, i interface{}) map[string]ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	items := make(map[string]ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(item.(interface{}))
	}

	return items
}

// flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumSlice flattens the contents of ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum from a JSON
// response object.
func flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumSlice(c *Client, i interface{}) []ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	items := make([]ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(item.(interface{})))
	}

	return items
}

// flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum with the same value as that string.
func flattenClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(i interface{}) *ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	s, ok := i.(string)
	if !ok {
		return ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef("")
	}

	return ClusterClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef(s)
}

// flattenClusterInstanceGroupConfigPreemptibilityEnumMap flattens the contents of ClusterInstanceGroupConfigPreemptibilityEnum from a JSON
// response object.
func flattenClusterInstanceGroupConfigPreemptibilityEnumMap(c *Client, i interface{}) map[string]ClusterInstanceGroupConfigPreemptibilityEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterInstanceGroupConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterInstanceGroupConfigPreemptibilityEnum{}
	}

	items := make(map[string]ClusterInstanceGroupConfigPreemptibilityEnum)
	for k, item := range a {
		items[k] = *flattenClusterInstanceGroupConfigPreemptibilityEnum(item.(interface{}))
	}

	return items
}

// flattenClusterInstanceGroupConfigPreemptibilityEnumSlice flattens the contents of ClusterInstanceGroupConfigPreemptibilityEnum from a JSON
// response object.
func flattenClusterInstanceGroupConfigPreemptibilityEnumSlice(c *Client, i interface{}) []ClusterInstanceGroupConfigPreemptibilityEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterInstanceGroupConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return []ClusterInstanceGroupConfigPreemptibilityEnum{}
	}

	items := make([]ClusterInstanceGroupConfigPreemptibilityEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterInstanceGroupConfigPreemptibilityEnum(item.(interface{})))
	}

	return items
}

// flattenClusterInstanceGroupConfigPreemptibilityEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterInstanceGroupConfigPreemptibilityEnum with the same value as that string.
func flattenClusterInstanceGroupConfigPreemptibilityEnum(i interface{}) *ClusterInstanceGroupConfigPreemptibilityEnum {
	s, ok := i.(string)
	if !ok {
		return ClusterInstanceGroupConfigPreemptibilityEnumRef("")
	}

	return ClusterInstanceGroupConfigPreemptibilityEnumRef(s)
}

// flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnumMap flattens the contents of ClusterClusterConfigSoftwareConfigOptionalComponentsEnum from a JSON
// response object.
func flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnumMap(c *Client, i interface{}) map[string]ClusterClusterConfigSoftwareConfigOptionalComponentsEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	items := make(map[string]ClusterClusterConfigSoftwareConfigOptionalComponentsEnum)
	for k, item := range a {
		items[k] = *flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnum(item.(interface{}))
	}

	return items
}

// flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnumSlice flattens the contents of ClusterClusterConfigSoftwareConfigOptionalComponentsEnum from a JSON
// response object.
func flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnumSlice(c *Client, i interface{}) []ClusterClusterConfigSoftwareConfigOptionalComponentsEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	if len(a) == 0 {
		return []ClusterClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	items := make([]ClusterClusterConfigSoftwareConfigOptionalComponentsEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnum(item.(interface{})))
	}

	return items
}

// flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterClusterConfigSoftwareConfigOptionalComponentsEnum with the same value as that string.
func flattenClusterClusterConfigSoftwareConfigOptionalComponentsEnum(i interface{}) *ClusterClusterConfigSoftwareConfigOptionalComponentsEnum {
	s, ok := i.(string)
	if !ok {
		return ClusterClusterConfigSoftwareConfigOptionalComponentsEnumRef("")
	}

	return ClusterClusterConfigSoftwareConfigOptionalComponentsEnumRef(s)
}

// flattenClusterStatusStateEnumMap flattens the contents of ClusterStatusStateEnum from a JSON
// response object.
func flattenClusterStatusStateEnumMap(c *Client, i interface{}) map[string]ClusterStatusStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStatusStateEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterStatusStateEnum{}
	}

	items := make(map[string]ClusterStatusStateEnum)
	for k, item := range a {
		items[k] = *flattenClusterStatusStateEnum(item.(interface{}))
	}

	return items
}

// flattenClusterStatusStateEnumSlice flattens the contents of ClusterStatusStateEnum from a JSON
// response object.
func flattenClusterStatusStateEnumSlice(c *Client, i interface{}) []ClusterStatusStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStatusStateEnum{}
	}

	if len(a) == 0 {
		return []ClusterStatusStateEnum{}
	}

	items := make([]ClusterStatusStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStatusStateEnum(item.(interface{})))
	}

	return items
}

// flattenClusterStatusStateEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterStatusStateEnum with the same value as that string.
func flattenClusterStatusStateEnum(i interface{}) *ClusterStatusStateEnum {
	s, ok := i.(string)
	if !ok {
		return ClusterStatusStateEnumRef("")
	}

	return ClusterStatusStateEnumRef(s)
}

// flattenClusterStatusSubstateEnumMap flattens the contents of ClusterStatusSubstateEnum from a JSON
// response object.
func flattenClusterStatusSubstateEnumMap(c *Client, i interface{}) map[string]ClusterStatusSubstateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStatusSubstateEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterStatusSubstateEnum{}
	}

	items := make(map[string]ClusterStatusSubstateEnum)
	for k, item := range a {
		items[k] = *flattenClusterStatusSubstateEnum(item.(interface{}))
	}

	return items
}

// flattenClusterStatusSubstateEnumSlice flattens the contents of ClusterStatusSubstateEnum from a JSON
// response object.
func flattenClusterStatusSubstateEnumSlice(c *Client, i interface{}) []ClusterStatusSubstateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStatusSubstateEnum{}
	}

	if len(a) == 0 {
		return []ClusterStatusSubstateEnum{}
	}

	items := make([]ClusterStatusSubstateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStatusSubstateEnum(item.(interface{})))
	}

	return items
}

// flattenClusterStatusSubstateEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterStatusSubstateEnum with the same value as that string.
func flattenClusterStatusSubstateEnum(i interface{}) *ClusterStatusSubstateEnum {
	s, ok := i.(string)
	if !ok {
		return ClusterStatusSubstateEnumRef("")
	}

	return ClusterStatusSubstateEnumRef(s)
}

// flattenClusterStatusHistoryStateEnumMap flattens the contents of ClusterStatusHistoryStateEnum from a JSON
// response object.
func flattenClusterStatusHistoryStateEnumMap(c *Client, i interface{}) map[string]ClusterStatusHistoryStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStatusHistoryStateEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterStatusHistoryStateEnum{}
	}

	items := make(map[string]ClusterStatusHistoryStateEnum)
	for k, item := range a {
		items[k] = *flattenClusterStatusHistoryStateEnum(item.(interface{}))
	}

	return items
}

// flattenClusterStatusHistoryStateEnumSlice flattens the contents of ClusterStatusHistoryStateEnum from a JSON
// response object.
func flattenClusterStatusHistoryStateEnumSlice(c *Client, i interface{}) []ClusterStatusHistoryStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStatusHistoryStateEnum{}
	}

	if len(a) == 0 {
		return []ClusterStatusHistoryStateEnum{}
	}

	items := make([]ClusterStatusHistoryStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStatusHistoryStateEnum(item.(interface{})))
	}

	return items
}

// flattenClusterStatusHistoryStateEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterStatusHistoryStateEnum with the same value as that string.
func flattenClusterStatusHistoryStateEnum(i interface{}) *ClusterStatusHistoryStateEnum {
	s, ok := i.(string)
	if !ok {
		return ClusterStatusHistoryStateEnumRef("")
	}

	return ClusterStatusHistoryStateEnumRef(s)
}

// flattenClusterStatusHistorySubstateEnumMap flattens the contents of ClusterStatusHistorySubstateEnum from a JSON
// response object.
func flattenClusterStatusHistorySubstateEnumMap(c *Client, i interface{}) map[string]ClusterStatusHistorySubstateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStatusHistorySubstateEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterStatusHistorySubstateEnum{}
	}

	items := make(map[string]ClusterStatusHistorySubstateEnum)
	for k, item := range a {
		items[k] = *flattenClusterStatusHistorySubstateEnum(item.(interface{}))
	}

	return items
}

// flattenClusterStatusHistorySubstateEnumSlice flattens the contents of ClusterStatusHistorySubstateEnum from a JSON
// response object.
func flattenClusterStatusHistorySubstateEnumSlice(c *Client, i interface{}) []ClusterStatusHistorySubstateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStatusHistorySubstateEnum{}
	}

	if len(a) == 0 {
		return []ClusterStatusHistorySubstateEnum{}
	}

	items := make([]ClusterStatusHistorySubstateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStatusHistorySubstateEnum(item.(interface{})))
	}

	return items
}

// flattenClusterStatusHistorySubstateEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterStatusHistorySubstateEnum with the same value as that string.
func flattenClusterStatusHistorySubstateEnum(i interface{}) *ClusterStatusHistorySubstateEnum {
	s, ok := i.(string)
	if !ok {
		return ClusterStatusHistorySubstateEnumRef("")
	}

	return ClusterStatusHistorySubstateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Cluster) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalCluster(b, c)
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

type clusterDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         clusterApiOperation
}

func convertFieldDiffsToClusterDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]clusterDiff, error) {
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
	var diffs []clusterDiff
	// For each operation name, create a clusterDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		diff := clusterDiff{}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToClusterApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToClusterApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (clusterApiOperation, error) {
	switch opName {

	case "updateClusterUpdateClusterOperation":
		return &updateClusterUpdateClusterOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractClusterFields(r *Cluster) error {
	// *ClusterClusterConfig is a reused type - that's not compatible with function extractors.

	vStatus := r.Status
	if vStatus == nil {
		// note: explicitly not the empty object.
		vStatus = &ClusterStatus{}
	}
	if err := extractClusterStatusFields(r, vStatus); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vStatus) {
		r.Status = vStatus
	}
	vMetrics := r.Metrics
	if vMetrics == nil {
		// note: explicitly not the empty object.
		vMetrics = &ClusterMetrics{}
	}
	if err := extractClusterMetricsFields(r, vMetrics); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vMetrics) {
		r.Metrics = vMetrics
	}
	return nil
}
func extractClusterClusterConfigFields(r *Cluster, o *ClusterClusterConfig) error {
	vGceClusterConfig := o.GceClusterConfig
	if vGceClusterConfig == nil {
		// note: explicitly not the empty object.
		vGceClusterConfig = &ClusterClusterConfigGceClusterConfig{}
	}
	if err := extractClusterClusterConfigGceClusterConfigFields(r, vGceClusterConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGceClusterConfig) {
		o.GceClusterConfig = vGceClusterConfig
	}
	// *ClusterInstanceGroupConfig is a reused type - that's not compatible with function extractors.

	// *ClusterInstanceGroupConfig is a reused type - that's not compatible with function extractors.

	// *ClusterInstanceGroupConfig is a reused type - that's not compatible with function extractors.

	vSoftwareConfig := o.SoftwareConfig
	if vSoftwareConfig == nil {
		// note: explicitly not the empty object.
		vSoftwareConfig = &ClusterClusterConfigSoftwareConfig{}
	}
	if err := extractClusterClusterConfigSoftwareConfigFields(r, vSoftwareConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSoftwareConfig) {
		o.SoftwareConfig = vSoftwareConfig
	}
	vEncryptionConfig := o.EncryptionConfig
	if vEncryptionConfig == nil {
		// note: explicitly not the empty object.
		vEncryptionConfig = &ClusterClusterConfigEncryptionConfig{}
	}
	if err := extractClusterClusterConfigEncryptionConfigFields(r, vEncryptionConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vEncryptionConfig) {
		o.EncryptionConfig = vEncryptionConfig
	}
	vAutoscalingConfig := o.AutoscalingConfig
	if vAutoscalingConfig == nil {
		// note: explicitly not the empty object.
		vAutoscalingConfig = &ClusterClusterConfigAutoscalingConfig{}
	}
	if err := extractClusterClusterConfigAutoscalingConfigFields(r, vAutoscalingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vAutoscalingConfig) {
		o.AutoscalingConfig = vAutoscalingConfig
	}
	vSecurityConfig := o.SecurityConfig
	if vSecurityConfig == nil {
		// note: explicitly not the empty object.
		vSecurityConfig = &ClusterClusterConfigSecurityConfig{}
	}
	if err := extractClusterClusterConfigSecurityConfigFields(r, vSecurityConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSecurityConfig) {
		o.SecurityConfig = vSecurityConfig
	}
	vLifecycleConfig := o.LifecycleConfig
	if vLifecycleConfig == nil {
		// note: explicitly not the empty object.
		vLifecycleConfig = &ClusterClusterConfigLifecycleConfig{}
	}
	if err := extractClusterClusterConfigLifecycleConfigFields(r, vLifecycleConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLifecycleConfig) {
		o.LifecycleConfig = vLifecycleConfig
	}
	vEndpointConfig := o.EndpointConfig
	if vEndpointConfig == nil {
		// note: explicitly not the empty object.
		vEndpointConfig = &ClusterClusterConfigEndpointConfig{}
	}
	if err := extractClusterClusterConfigEndpointConfigFields(r, vEndpointConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vEndpointConfig) {
		o.EndpointConfig = vEndpointConfig
	}
	return nil
}
func extractClusterClusterConfigGceClusterConfigFields(r *Cluster, o *ClusterClusterConfigGceClusterConfig) error {
	vReservationAffinity := o.ReservationAffinity
	if vReservationAffinity == nil {
		// note: explicitly not the empty object.
		vReservationAffinity = &ClusterClusterConfigGceClusterConfigReservationAffinity{}
	}
	if err := extractClusterClusterConfigGceClusterConfigReservationAffinityFields(r, vReservationAffinity); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vReservationAffinity) {
		o.ReservationAffinity = vReservationAffinity
	}
	vNodeGroupAffinity := o.NodeGroupAffinity
	if vNodeGroupAffinity == nil {
		// note: explicitly not the empty object.
		vNodeGroupAffinity = &ClusterClusterConfigGceClusterConfigNodeGroupAffinity{}
	}
	if err := extractClusterClusterConfigGceClusterConfigNodeGroupAffinityFields(r, vNodeGroupAffinity); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vNodeGroupAffinity) {
		o.NodeGroupAffinity = vNodeGroupAffinity
	}
	return nil
}
func extractClusterClusterConfigGceClusterConfigReservationAffinityFields(r *Cluster, o *ClusterClusterConfigGceClusterConfigReservationAffinity) error {
	return nil
}
func extractClusterClusterConfigGceClusterConfigNodeGroupAffinityFields(r *Cluster, o *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) error {
	return nil
}
func extractClusterInstanceGroupConfigFields(r *Cluster, o *ClusterInstanceGroupConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &ClusterInstanceGroupConfigDiskConfig{}
	}
	if err := extractClusterInstanceGroupConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &ClusterInstanceGroupConfigManagedGroupConfig{}
	}
	if err := extractClusterInstanceGroupConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func extractClusterInstanceGroupConfigDiskConfigFields(r *Cluster, o *ClusterInstanceGroupConfigDiskConfig) error {
	return nil
}
func extractClusterInstanceGroupConfigManagedGroupConfigFields(r *Cluster, o *ClusterInstanceGroupConfigManagedGroupConfig) error {
	return nil
}
func extractClusterInstanceGroupConfigAcceleratorsFields(r *Cluster, o *ClusterInstanceGroupConfigAccelerators) error {
	return nil
}
func extractClusterClusterConfigSoftwareConfigFields(r *Cluster, o *ClusterClusterConfigSoftwareConfig) error {
	return nil
}
func extractClusterClusterConfigInitializationActionsFields(r *Cluster, o *ClusterClusterConfigInitializationActions) error {
	return nil
}
func extractClusterClusterConfigEncryptionConfigFields(r *Cluster, o *ClusterClusterConfigEncryptionConfig) error {
	return nil
}
func extractClusterClusterConfigAutoscalingConfigFields(r *Cluster, o *ClusterClusterConfigAutoscalingConfig) error {
	return nil
}
func extractClusterClusterConfigSecurityConfigFields(r *Cluster, o *ClusterClusterConfigSecurityConfig) error {
	vKerberosConfig := o.KerberosConfig
	if vKerberosConfig == nil {
		// note: explicitly not the empty object.
		vKerberosConfig = &ClusterClusterConfigSecurityConfigKerberosConfig{}
	}
	if err := extractClusterClusterConfigSecurityConfigKerberosConfigFields(r, vKerberosConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vKerberosConfig) {
		o.KerberosConfig = vKerberosConfig
	}
	return nil
}
func extractClusterClusterConfigSecurityConfigKerberosConfigFields(r *Cluster, o *ClusterClusterConfigSecurityConfigKerberosConfig) error {
	return nil
}
func extractClusterClusterConfigLifecycleConfigFields(r *Cluster, o *ClusterClusterConfigLifecycleConfig) error {
	return nil
}
func extractClusterClusterConfigEndpointConfigFields(r *Cluster, o *ClusterClusterConfigEndpointConfig) error {
	return nil
}
func extractClusterStatusFields(r *Cluster, o *ClusterStatus) error {
	return nil
}
func extractClusterStatusHistoryFields(r *Cluster, o *ClusterStatusHistory) error {
	return nil
}
func extractClusterMetricsFields(r *Cluster, o *ClusterMetrics) error {
	return nil
}

func postReadExtractClusterFields(r *Cluster) error {
	// *ClusterClusterConfig is a reused type - that's not compatible with function extractors.

	vStatus := r.Status
	if vStatus == nil {
		// note: explicitly not the empty object.
		vStatus = &ClusterStatus{}
	}
	if err := postReadExtractClusterStatusFields(r, vStatus); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vStatus) {
		r.Status = vStatus
	}
	vMetrics := r.Metrics
	if vMetrics == nil {
		// note: explicitly not the empty object.
		vMetrics = &ClusterMetrics{}
	}
	if err := postReadExtractClusterMetricsFields(r, vMetrics); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vMetrics) {
		r.Metrics = vMetrics
	}
	return nil
}
func postReadExtractClusterClusterConfigFields(r *Cluster, o *ClusterClusterConfig) error {
	vGceClusterConfig := o.GceClusterConfig
	if vGceClusterConfig == nil {
		// note: explicitly not the empty object.
		vGceClusterConfig = &ClusterClusterConfigGceClusterConfig{}
	}
	if err := extractClusterClusterConfigGceClusterConfigFields(r, vGceClusterConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGceClusterConfig) {
		o.GceClusterConfig = vGceClusterConfig
	}
	// *ClusterInstanceGroupConfig is a reused type - that's not compatible with function extractors.

	// *ClusterInstanceGroupConfig is a reused type - that's not compatible with function extractors.

	// *ClusterInstanceGroupConfig is a reused type - that's not compatible with function extractors.

	vSoftwareConfig := o.SoftwareConfig
	if vSoftwareConfig == nil {
		// note: explicitly not the empty object.
		vSoftwareConfig = &ClusterClusterConfigSoftwareConfig{}
	}
	if err := extractClusterClusterConfigSoftwareConfigFields(r, vSoftwareConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSoftwareConfig) {
		o.SoftwareConfig = vSoftwareConfig
	}
	vEncryptionConfig := o.EncryptionConfig
	if vEncryptionConfig == nil {
		// note: explicitly not the empty object.
		vEncryptionConfig = &ClusterClusterConfigEncryptionConfig{}
	}
	if err := extractClusterClusterConfigEncryptionConfigFields(r, vEncryptionConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vEncryptionConfig) {
		o.EncryptionConfig = vEncryptionConfig
	}
	vAutoscalingConfig := o.AutoscalingConfig
	if vAutoscalingConfig == nil {
		// note: explicitly not the empty object.
		vAutoscalingConfig = &ClusterClusterConfigAutoscalingConfig{}
	}
	if err := extractClusterClusterConfigAutoscalingConfigFields(r, vAutoscalingConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vAutoscalingConfig) {
		o.AutoscalingConfig = vAutoscalingConfig
	}
	vSecurityConfig := o.SecurityConfig
	if vSecurityConfig == nil {
		// note: explicitly not the empty object.
		vSecurityConfig = &ClusterClusterConfigSecurityConfig{}
	}
	if err := extractClusterClusterConfigSecurityConfigFields(r, vSecurityConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vSecurityConfig) {
		o.SecurityConfig = vSecurityConfig
	}
	vLifecycleConfig := o.LifecycleConfig
	if vLifecycleConfig == nil {
		// note: explicitly not the empty object.
		vLifecycleConfig = &ClusterClusterConfigLifecycleConfig{}
	}
	if err := extractClusterClusterConfigLifecycleConfigFields(r, vLifecycleConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vLifecycleConfig) {
		o.LifecycleConfig = vLifecycleConfig
	}
	vEndpointConfig := o.EndpointConfig
	if vEndpointConfig == nil {
		// note: explicitly not the empty object.
		vEndpointConfig = &ClusterClusterConfigEndpointConfig{}
	}
	if err := extractClusterClusterConfigEndpointConfigFields(r, vEndpointConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vEndpointConfig) {
		o.EndpointConfig = vEndpointConfig
	}
	return nil
}
func postReadExtractClusterClusterConfigGceClusterConfigFields(r *Cluster, o *ClusterClusterConfigGceClusterConfig) error {
	vReservationAffinity := o.ReservationAffinity
	if vReservationAffinity == nil {
		// note: explicitly not the empty object.
		vReservationAffinity = &ClusterClusterConfigGceClusterConfigReservationAffinity{}
	}
	if err := extractClusterClusterConfigGceClusterConfigReservationAffinityFields(r, vReservationAffinity); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vReservationAffinity) {
		o.ReservationAffinity = vReservationAffinity
	}
	vNodeGroupAffinity := o.NodeGroupAffinity
	if vNodeGroupAffinity == nil {
		// note: explicitly not the empty object.
		vNodeGroupAffinity = &ClusterClusterConfigGceClusterConfigNodeGroupAffinity{}
	}
	if err := extractClusterClusterConfigGceClusterConfigNodeGroupAffinityFields(r, vNodeGroupAffinity); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vNodeGroupAffinity) {
		o.NodeGroupAffinity = vNodeGroupAffinity
	}
	return nil
}
func postReadExtractClusterClusterConfigGceClusterConfigReservationAffinityFields(r *Cluster, o *ClusterClusterConfigGceClusterConfigReservationAffinity) error {
	return nil
}
func postReadExtractClusterClusterConfigGceClusterConfigNodeGroupAffinityFields(r *Cluster, o *ClusterClusterConfigGceClusterConfigNodeGroupAffinity) error {
	return nil
}
func postReadExtractClusterInstanceGroupConfigFields(r *Cluster, o *ClusterInstanceGroupConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &ClusterInstanceGroupConfigDiskConfig{}
	}
	if err := extractClusterInstanceGroupConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &ClusterInstanceGroupConfigManagedGroupConfig{}
	}
	if err := extractClusterInstanceGroupConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func postReadExtractClusterInstanceGroupConfigDiskConfigFields(r *Cluster, o *ClusterInstanceGroupConfigDiskConfig) error {
	return nil
}
func postReadExtractClusterInstanceGroupConfigManagedGroupConfigFields(r *Cluster, o *ClusterInstanceGroupConfigManagedGroupConfig) error {
	return nil
}
func postReadExtractClusterInstanceGroupConfigAcceleratorsFields(r *Cluster, o *ClusterInstanceGroupConfigAccelerators) error {
	return nil
}
func postReadExtractClusterClusterConfigSoftwareConfigFields(r *Cluster, o *ClusterClusterConfigSoftwareConfig) error {
	return nil
}
func postReadExtractClusterClusterConfigInitializationActionsFields(r *Cluster, o *ClusterClusterConfigInitializationActions) error {
	return nil
}
func postReadExtractClusterClusterConfigEncryptionConfigFields(r *Cluster, o *ClusterClusterConfigEncryptionConfig) error {
	return nil
}
func postReadExtractClusterClusterConfigAutoscalingConfigFields(r *Cluster, o *ClusterClusterConfigAutoscalingConfig) error {
	return nil
}
func postReadExtractClusterClusterConfigSecurityConfigFields(r *Cluster, o *ClusterClusterConfigSecurityConfig) error {
	vKerberosConfig := o.KerberosConfig
	if vKerberosConfig == nil {
		// note: explicitly not the empty object.
		vKerberosConfig = &ClusterClusterConfigSecurityConfigKerberosConfig{}
	}
	if err := extractClusterClusterConfigSecurityConfigKerberosConfigFields(r, vKerberosConfig); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vKerberosConfig) {
		o.KerberosConfig = vKerberosConfig
	}
	return nil
}
func postReadExtractClusterClusterConfigSecurityConfigKerberosConfigFields(r *Cluster, o *ClusterClusterConfigSecurityConfigKerberosConfig) error {
	return nil
}
func postReadExtractClusterClusterConfigLifecycleConfigFields(r *Cluster, o *ClusterClusterConfigLifecycleConfig) error {
	return nil
}
func postReadExtractClusterClusterConfigEndpointConfigFields(r *Cluster, o *ClusterClusterConfigEndpointConfig) error {
	return nil
}
func postReadExtractClusterStatusFields(r *Cluster, o *ClusterStatus) error {
	return nil
}
func postReadExtractClusterStatusHistoryFields(r *Cluster, o *ClusterStatusHistory) error {
	return nil
}
func postReadExtractClusterMetricsFields(r *Cluster, o *ClusterMetrics) error {
	return nil
}
