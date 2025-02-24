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
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *Cluster) validate() error {

	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
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
	if !dcl.IsEmptyValueIndirect(r.VirtualClusterConfig) {
		if err := r.VirtualClusterConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterConfig) validate() error {
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
	if !dcl.IsEmptyValueIndirect(r.MetastoreConfig) {
		if err := r.MetastoreConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.DataprocMetricConfig) {
		if err := r.DataprocMetricConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterConfigGceClusterConfig) validate() error {
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
	if !dcl.IsEmptyValueIndirect(r.ShieldedInstanceConfig) {
		if err := r.ShieldedInstanceConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ConfidentialInstanceConfig) {
		if err := r.ConfidentialInstanceConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterConfigGceClusterConfigReservationAffinity) validate() error {
	return nil
}
func (r *ClusterConfigGceClusterConfigNodeGroupAffinity) validate() error {
	if err := dcl.Required(r, "nodeGroup"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterConfigGceClusterConfigShieldedInstanceConfig) validate() error {
	return nil
}
func (r *ClusterConfigGceClusterConfigConfidentialInstanceConfig) validate() error {
	return nil
}
func (r *ClusterConfigMasterConfig) validate() error {
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
func (r *ClusterConfigMasterConfigDiskConfig) validate() error {
	return nil
}
func (r *ClusterConfigMasterConfigManagedGroupConfig) validate() error {
	return nil
}
func (r *ClusterConfigMasterConfigAccelerators) validate() error {
	return nil
}
func (r *ClusterConfigMasterConfigInstanceReferences) validate() error {
	return nil
}
func (r *ClusterConfigWorkerConfig) validate() error {
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
func (r *ClusterConfigWorkerConfigDiskConfig) validate() error {
	return nil
}
func (r *ClusterConfigWorkerConfigManagedGroupConfig) validate() error {
	return nil
}
func (r *ClusterConfigWorkerConfigAccelerators) validate() error {
	return nil
}
func (r *ClusterConfigWorkerConfigInstanceReferences) validate() error {
	return nil
}
func (r *ClusterConfigSecondaryWorkerConfig) validate() error {
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
func (r *ClusterConfigSecondaryWorkerConfigDiskConfig) validate() error {
	return nil
}
func (r *ClusterConfigSecondaryWorkerConfigManagedGroupConfig) validate() error {
	return nil
}
func (r *ClusterConfigSecondaryWorkerConfigAccelerators) validate() error {
	return nil
}
func (r *ClusterConfigSecondaryWorkerConfigInstanceReferences) validate() error {
	return nil
}
func (r *ClusterConfigSoftwareConfig) validate() error {
	return nil
}
func (r *ClusterConfigInitializationActions) validate() error {
	if err := dcl.Required(r, "executableFile"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterConfigEncryptionConfig) validate() error {
	return nil
}
func (r *ClusterConfigAutoscalingConfig) validate() error {
	return nil
}
func (r *ClusterConfigSecurityConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.KerberosConfig) {
		if err := r.KerberosConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.IdentityConfig) {
		if err := r.IdentityConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterConfigSecurityConfigKerberosConfig) validate() error {
	return nil
}
func (r *ClusterConfigSecurityConfigIdentityConfig) validate() error {
	if err := dcl.Required(r, "userServiceAccountMapping"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterConfigLifecycleConfig) validate() error {
	return nil
}
func (r *ClusterConfigEndpointConfig) validate() error {
	return nil
}
func (r *ClusterConfigMetastoreConfig) validate() error {
	if err := dcl.Required(r, "dataprocMetastoreService"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterConfigDataprocMetricConfig) validate() error {
	if err := dcl.Required(r, "metrics"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterConfigDataprocMetricConfigMetrics) validate() error {
	if err := dcl.Required(r, "metricSource"); err != nil {
		return err
	}
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
func (r *ClusterVirtualClusterConfig) validate() error {
	if err := dcl.Required(r, "kubernetesClusterConfig"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.KubernetesClusterConfig) {
		if err := r.KubernetesClusterConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AuxiliaryServicesConfig) {
		if err := r.AuxiliaryServicesConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfig) validate() error {
	if err := dcl.Required(r, "gkeClusterConfig"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GkeClusterConfig) {
		if err := r.GkeClusterConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.KubernetesSoftwareConfig) {
		if err := r.KubernetesSoftwareConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) validate() error {
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) validate() error {
	if err := dcl.Required(r, "nodePool"); err != nil {
		return err
	}
	if err := dcl.Required(r, "roles"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.NodePoolConfig) {
		if err := r.NodePoolConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) validate() error {
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
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.EphemeralStorageConfig) {
		if err := r.EphemeralStorageConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) validate() error {
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) validate() error {
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) validate() error {
	return nil
}
func (r *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) validate() error {
	return nil
}
func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.MetastoreConfig) {
		if err := r.MetastoreConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SparkHistoryServerConfig) {
		if err := r.SparkHistoryServerConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) validate() error {
	if err := dcl.Required(r, "dataprocMetastoreService"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) validate() error {
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
	res := f
	_ = res

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
		res, err := unmarshalMapCluster(v, c, r)
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

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetCluster(ctx, r)
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

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractClusterFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

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
		rawDesired.Config = canonicalizeClusterConfig(rawDesired.Config, nil, opts...)
		rawDesired.Status = canonicalizeClusterStatus(rawDesired.Status, nil, opts...)
		rawDesired.Metrics = canonicalizeClusterMetrics(rawDesired.Metrics, nil, opts...)
		rawDesired.VirtualClusterConfig = canonicalizeClusterVirtualClusterConfig(rawDesired.VirtualClusterConfig, nil, opts...)

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
	canonicalDesired.Config = canonicalizeClusterConfig(rawDesired.Config, rawInitial.Config, opts...)
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.NameToSelfLink(rawDesired.Location, rawInitial.Location) {
		canonicalDesired.Location = rawInitial.Location
	} else {
		canonicalDesired.Location = rawDesired.Location
	}
	canonicalDesired.VirtualClusterConfig = canonicalizeClusterVirtualClusterConfig(rawDesired.VirtualClusterConfig, rawInitial.VirtualClusterConfig, opts...)
	return canonicalDesired, nil
}

func canonicalizeClusterNewState(c *Client, rawNew, rawDesired *Cluster) (*Cluster, error) {

	rawNew.Project = rawDesired.Project

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Config) && dcl.IsEmptyValueIndirect(rawDesired.Config) {
		rawNew.Config = rawDesired.Config
	} else {
		rawNew.Config = canonicalizeNewClusterConfig(c, rawDesired.Config, rawNew.Config)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Status) && dcl.IsEmptyValueIndirect(rawDesired.Status) {
		rawNew.Status = rawDesired.Status
	} else {
		rawNew.Status = canonicalizeNewClusterStatus(c, rawDesired.Status, rawNew.Status)
	}

	if dcl.IsEmptyValueIndirect(rawNew.StatusHistory) && dcl.IsEmptyValueIndirect(rawDesired.StatusHistory) {
		rawNew.StatusHistory = rawDesired.StatusHistory
	} else {
		rawNew.StatusHistory = canonicalizeNewClusterStatusHistorySlice(c, rawDesired.StatusHistory, rawNew.StatusHistory)
	}

	if dcl.IsEmptyValueIndirect(rawNew.ClusterUuid) && dcl.IsEmptyValueIndirect(rawDesired.ClusterUuid) {
		rawNew.ClusterUuid = rawDesired.ClusterUuid
	} else {
		if dcl.StringCanonicalize(rawDesired.ClusterUuid, rawNew.ClusterUuid) {
			rawNew.ClusterUuid = rawDesired.ClusterUuid
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Metrics) && dcl.IsEmptyValueIndirect(rawDesired.Metrics) {
		rawNew.Metrics = rawDesired.Metrics
	} else {
		rawNew.Metrics = canonicalizeNewClusterMetrics(c, rawDesired.Metrics, rawNew.Metrics)
	}

	rawNew.Location = rawDesired.Location

	if dcl.IsEmptyValueIndirect(rawNew.VirtualClusterConfig) && dcl.IsEmptyValueIndirect(rawDesired.VirtualClusterConfig) {
		rawNew.VirtualClusterConfig = rawDesired.VirtualClusterConfig
	} else {
		rawNew.VirtualClusterConfig = canonicalizeNewClusterVirtualClusterConfig(c, rawDesired.VirtualClusterConfig, rawNew.VirtualClusterConfig)
	}

	return rawNew, nil
}

func canonicalizeClusterConfig(des, initial *ClusterConfig, opts ...dcl.ApplyOption) *ClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfig{}

	if dcl.IsZeroValue(des.StagingBucket) || (dcl.IsEmptyValueIndirect(des.StagingBucket) && dcl.IsEmptyValueIndirect(initial.StagingBucket)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.StagingBucket = initial.StagingBucket
	} else {
		cDes.StagingBucket = des.StagingBucket
	}
	if dcl.IsZeroValue(des.TempBucket) || (dcl.IsEmptyValueIndirect(des.TempBucket) && dcl.IsEmptyValueIndirect(initial.TempBucket)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.TempBucket = initial.TempBucket
	} else {
		cDes.TempBucket = des.TempBucket
	}
	cDes.GceClusterConfig = canonicalizeClusterConfigGceClusterConfig(des.GceClusterConfig, initial.GceClusterConfig, opts...)
	cDes.MasterConfig = canonicalizeClusterConfigMasterConfig(des.MasterConfig, initial.MasterConfig, opts...)
	cDes.WorkerConfig = canonicalizeClusterConfigWorkerConfig(des.WorkerConfig, initial.WorkerConfig, opts...)
	cDes.SecondaryWorkerConfig = canonicalizeClusterConfigSecondaryWorkerConfig(des.SecondaryWorkerConfig, initial.SecondaryWorkerConfig, opts...)
	cDes.SoftwareConfig = canonicalizeClusterConfigSoftwareConfig(des.SoftwareConfig, initial.SoftwareConfig, opts...)
	cDes.InitializationActions = canonicalizeClusterConfigInitializationActionsSlice(des.InitializationActions, initial.InitializationActions, opts...)
	cDes.EncryptionConfig = canonicalizeClusterConfigEncryptionConfig(des.EncryptionConfig, initial.EncryptionConfig, opts...)
	cDes.AutoscalingConfig = canonicalizeClusterConfigAutoscalingConfig(des.AutoscalingConfig, initial.AutoscalingConfig, opts...)
	cDes.SecurityConfig = canonicalizeClusterConfigSecurityConfig(des.SecurityConfig, initial.SecurityConfig, opts...)
	cDes.LifecycleConfig = canonicalizeClusterConfigLifecycleConfig(des.LifecycleConfig, initial.LifecycleConfig, opts...)
	cDes.EndpointConfig = canonicalizeClusterConfigEndpointConfig(des.EndpointConfig, initial.EndpointConfig, opts...)
	cDes.MetastoreConfig = canonicalizeClusterConfigMetastoreConfig(des.MetastoreConfig, initial.MetastoreConfig, opts...)
	cDes.DataprocMetricConfig = canonicalizeClusterConfigDataprocMetricConfig(des.DataprocMetricConfig, initial.DataprocMetricConfig, opts...)

	return cDes
}

func canonicalizeClusterConfigSlice(des, initial []ClusterConfig, opts ...dcl.ApplyOption) []ClusterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfig(c *Client, des, nw *ClusterConfig) *ClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.GceClusterConfig = canonicalizeNewClusterConfigGceClusterConfig(c, des.GceClusterConfig, nw.GceClusterConfig)
	nw.MasterConfig = canonicalizeNewClusterConfigMasterConfig(c, des.MasterConfig, nw.MasterConfig)
	nw.WorkerConfig = canonicalizeNewClusterConfigWorkerConfig(c, des.WorkerConfig, nw.WorkerConfig)
	nw.SecondaryWorkerConfig = canonicalizeNewClusterConfigSecondaryWorkerConfig(c, des.SecondaryWorkerConfig, nw.SecondaryWorkerConfig)
	nw.SoftwareConfig = canonicalizeNewClusterConfigSoftwareConfig(c, des.SoftwareConfig, nw.SoftwareConfig)
	nw.InitializationActions = canonicalizeNewClusterConfigInitializationActionsSlice(c, des.InitializationActions, nw.InitializationActions)
	nw.EncryptionConfig = canonicalizeNewClusterConfigEncryptionConfig(c, des.EncryptionConfig, nw.EncryptionConfig)
	nw.AutoscalingConfig = canonicalizeNewClusterConfigAutoscalingConfig(c, des.AutoscalingConfig, nw.AutoscalingConfig)
	nw.SecurityConfig = canonicalizeNewClusterConfigSecurityConfig(c, des.SecurityConfig, nw.SecurityConfig)
	nw.LifecycleConfig = canonicalizeNewClusterConfigLifecycleConfig(c, des.LifecycleConfig, nw.LifecycleConfig)
	nw.EndpointConfig = canonicalizeNewClusterConfigEndpointConfig(c, des.EndpointConfig, nw.EndpointConfig)
	nw.MetastoreConfig = canonicalizeNewClusterConfigMetastoreConfig(c, des.MetastoreConfig, nw.MetastoreConfig)
	nw.DataprocMetricConfig = canonicalizeNewClusterConfigDataprocMetricConfig(c, des.DataprocMetricConfig, nw.DataprocMetricConfig)

	return nw
}

func canonicalizeNewClusterConfigSet(c *Client, des, nw []ClusterConfig) []ClusterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSlice(c *Client, des, nw []ClusterConfig) []ClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigGceClusterConfig(des, initial *ClusterConfigGceClusterConfig, opts ...dcl.ApplyOption) *ClusterConfigGceClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigGceClusterConfig{}

	if dcl.StringCanonicalize(des.Zone, initial.Zone) || dcl.IsZeroValue(des.Zone) {
		cDes.Zone = initial.Zone
	} else {
		cDes.Zone = des.Zone
	}
	if dcl.IsZeroValue(des.Network) || (dcl.IsEmptyValueIndirect(des.Network) && dcl.IsEmptyValueIndirect(initial.Network)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Network = initial.Network
	} else {
		cDes.Network = des.Network
	}
	if dcl.IsZeroValue(des.Subnetwork) || (dcl.IsEmptyValueIndirect(des.Subnetwork) && dcl.IsEmptyValueIndirect(initial.Subnetwork)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Subnetwork = initial.Subnetwork
	} else {
		cDes.Subnetwork = des.Subnetwork
	}
	if dcl.BoolCanonicalize(des.InternalIPOnly, initial.InternalIPOnly) || dcl.IsZeroValue(des.InternalIPOnly) {
		cDes.InternalIPOnly = initial.InternalIPOnly
	} else {
		cDes.InternalIPOnly = des.InternalIPOnly
	}
	if dcl.IsZeroValue(des.PrivateIPv6GoogleAccess) || (dcl.IsEmptyValueIndirect(des.PrivateIPv6GoogleAccess) && dcl.IsEmptyValueIndirect(initial.PrivateIPv6GoogleAccess)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.PrivateIPv6GoogleAccess = initial.PrivateIPv6GoogleAccess
	} else {
		cDes.PrivateIPv6GoogleAccess = des.PrivateIPv6GoogleAccess
	}
	if dcl.IsZeroValue(des.ServiceAccount) || (dcl.IsEmptyValueIndirect(des.ServiceAccount) && dcl.IsEmptyValueIndirect(initial.ServiceAccount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ServiceAccount = initial.ServiceAccount
	} else {
		cDes.ServiceAccount = des.ServiceAccount
	}
	if dcl.StringArrayCanonicalize(des.ServiceAccountScopes, initial.ServiceAccountScopes) {
		cDes.ServiceAccountScopes = initial.ServiceAccountScopes
	} else {
		cDes.ServiceAccountScopes = des.ServiceAccountScopes
	}
	if dcl.StringArrayCanonicalize(des.Tags, initial.Tags) {
		cDes.Tags = initial.Tags
	} else {
		cDes.Tags = des.Tags
	}
	if dcl.IsZeroValue(des.Metadata) || (dcl.IsEmptyValueIndirect(des.Metadata) && dcl.IsEmptyValueIndirect(initial.Metadata)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Metadata = initial.Metadata
	} else {
		cDes.Metadata = des.Metadata
	}
	cDes.ReservationAffinity = canonicalizeClusterConfigGceClusterConfigReservationAffinity(des.ReservationAffinity, initial.ReservationAffinity, opts...)
	cDes.NodeGroupAffinity = canonicalizeClusterConfigGceClusterConfigNodeGroupAffinity(des.NodeGroupAffinity, initial.NodeGroupAffinity, opts...)
	cDes.ShieldedInstanceConfig = canonicalizeClusterConfigGceClusterConfigShieldedInstanceConfig(des.ShieldedInstanceConfig, initial.ShieldedInstanceConfig, opts...)
	cDes.ConfidentialInstanceConfig = canonicalizeClusterConfigGceClusterConfigConfidentialInstanceConfig(des.ConfidentialInstanceConfig, initial.ConfidentialInstanceConfig, opts...)

	return cDes
}

func canonicalizeClusterConfigGceClusterConfigSlice(des, initial []ClusterConfigGceClusterConfig, opts ...dcl.ApplyOption) []ClusterConfigGceClusterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigGceClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigGceClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigGceClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigGceClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigGceClusterConfig(c *Client, des, nw *ClusterConfigGceClusterConfig) *ClusterConfigGceClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigGceClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Zone, nw.Zone) {
		nw.Zone = des.Zone
	}
	if dcl.BoolCanonicalize(des.InternalIPOnly, nw.InternalIPOnly) {
		nw.InternalIPOnly = des.InternalIPOnly
	}
	if dcl.StringArrayCanonicalize(des.ServiceAccountScopes, nw.ServiceAccountScopes) {
		nw.ServiceAccountScopes = des.ServiceAccountScopes
	}
	if dcl.StringArrayCanonicalize(des.Tags, nw.Tags) {
		nw.Tags = des.Tags
	}
	nw.ReservationAffinity = canonicalizeNewClusterConfigGceClusterConfigReservationAffinity(c, des.ReservationAffinity, nw.ReservationAffinity)
	nw.NodeGroupAffinity = canonicalizeNewClusterConfigGceClusterConfigNodeGroupAffinity(c, des.NodeGroupAffinity, nw.NodeGroupAffinity)
	nw.ShieldedInstanceConfig = canonicalizeNewClusterConfigGceClusterConfigShieldedInstanceConfig(c, des.ShieldedInstanceConfig, nw.ShieldedInstanceConfig)
	nw.ConfidentialInstanceConfig = canonicalizeNewClusterConfigGceClusterConfigConfidentialInstanceConfig(c, des.ConfidentialInstanceConfig, nw.ConfidentialInstanceConfig)

	return nw
}

func canonicalizeNewClusterConfigGceClusterConfigSet(c *Client, des, nw []ClusterConfigGceClusterConfig) []ClusterConfigGceClusterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigGceClusterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigGceClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigGceClusterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigGceClusterConfigSlice(c *Client, des, nw []ClusterConfigGceClusterConfig) []ClusterConfigGceClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigGceClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigGceClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigGceClusterConfigReservationAffinity(des, initial *ClusterConfigGceClusterConfigReservationAffinity, opts ...dcl.ApplyOption) *ClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigGceClusterConfigReservationAffinity{}

	if dcl.IsZeroValue(des.ConsumeReservationType) || (dcl.IsEmptyValueIndirect(des.ConsumeReservationType) && dcl.IsEmptyValueIndirect(initial.ConsumeReservationType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ConsumeReservationType = initial.ConsumeReservationType
	} else {
		cDes.ConsumeReservationType = des.ConsumeReservationType
	}
	if dcl.StringCanonicalize(des.Key, initial.Key) || dcl.IsZeroValue(des.Key) {
		cDes.Key = initial.Key
	} else {
		cDes.Key = des.Key
	}
	if dcl.StringArrayCanonicalize(des.Values, initial.Values) {
		cDes.Values = initial.Values
	} else {
		cDes.Values = des.Values
	}

	return cDes
}

func canonicalizeClusterConfigGceClusterConfigReservationAffinitySlice(des, initial []ClusterConfigGceClusterConfigReservationAffinity, opts ...dcl.ApplyOption) []ClusterConfigGceClusterConfigReservationAffinity {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigGceClusterConfigReservationAffinity, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigGceClusterConfigReservationAffinity(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigGceClusterConfigReservationAffinity, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigGceClusterConfigReservationAffinity(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigGceClusterConfigReservationAffinity(c *Client, des, nw *ClusterConfigGceClusterConfigReservationAffinity) *ClusterConfigGceClusterConfigReservationAffinity {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigGceClusterConfigReservationAffinity while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewClusterConfigGceClusterConfigReservationAffinitySet(c *Client, des, nw []ClusterConfigGceClusterConfigReservationAffinity) []ClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigGceClusterConfigReservationAffinity
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigGceClusterConfigReservationAffinityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigGceClusterConfigReservationAffinity(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, des, nw []ClusterConfigGceClusterConfigReservationAffinity) []ClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigGceClusterConfigReservationAffinity
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigGceClusterConfigReservationAffinity(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigGceClusterConfigNodeGroupAffinity(des, initial *ClusterConfigGceClusterConfigNodeGroupAffinity, opts ...dcl.ApplyOption) *ClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigGceClusterConfigNodeGroupAffinity{}

	if dcl.IsZeroValue(des.NodeGroup) || (dcl.IsEmptyValueIndirect(des.NodeGroup) && dcl.IsEmptyValueIndirect(initial.NodeGroup)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NodeGroup = initial.NodeGroup
	} else {
		cDes.NodeGroup = des.NodeGroup
	}

	return cDes
}

func canonicalizeClusterConfigGceClusterConfigNodeGroupAffinitySlice(des, initial []ClusterConfigGceClusterConfigNodeGroupAffinity, opts ...dcl.ApplyOption) []ClusterConfigGceClusterConfigNodeGroupAffinity {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigGceClusterConfigNodeGroupAffinity(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigGceClusterConfigNodeGroupAffinity(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, des, nw *ClusterConfigGceClusterConfigNodeGroupAffinity) *ClusterConfigGceClusterConfigNodeGroupAffinity {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigGceClusterConfigNodeGroupAffinity while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterConfigGceClusterConfigNodeGroupAffinitySet(c *Client, des, nw []ClusterConfigGceClusterConfigNodeGroupAffinity) []ClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigGceClusterConfigNodeGroupAffinity
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigGceClusterConfigNodeGroupAffinityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigGceClusterConfigNodeGroupAffinity(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, des, nw []ClusterConfigGceClusterConfigNodeGroupAffinity) []ClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigGceClusterConfigNodeGroupAffinity
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigGceClusterConfigNodeGroupAffinity(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigGceClusterConfigShieldedInstanceConfig(des, initial *ClusterConfigGceClusterConfigShieldedInstanceConfig, opts ...dcl.ApplyOption) *ClusterConfigGceClusterConfigShieldedInstanceConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigGceClusterConfigShieldedInstanceConfig{}

	if dcl.BoolCanonicalize(des.EnableSecureBoot, initial.EnableSecureBoot) || dcl.IsZeroValue(des.EnableSecureBoot) {
		cDes.EnableSecureBoot = initial.EnableSecureBoot
	} else {
		cDes.EnableSecureBoot = des.EnableSecureBoot
	}
	if dcl.BoolCanonicalize(des.EnableVtpm, initial.EnableVtpm) || dcl.IsZeroValue(des.EnableVtpm) {
		cDes.EnableVtpm = initial.EnableVtpm
	} else {
		cDes.EnableVtpm = des.EnableVtpm
	}
	if dcl.BoolCanonicalize(des.EnableIntegrityMonitoring, initial.EnableIntegrityMonitoring) || dcl.IsZeroValue(des.EnableIntegrityMonitoring) {
		cDes.EnableIntegrityMonitoring = initial.EnableIntegrityMonitoring
	} else {
		cDes.EnableIntegrityMonitoring = des.EnableIntegrityMonitoring
	}

	return cDes
}

func canonicalizeClusterConfigGceClusterConfigShieldedInstanceConfigSlice(des, initial []ClusterConfigGceClusterConfigShieldedInstanceConfig, opts ...dcl.ApplyOption) []ClusterConfigGceClusterConfigShieldedInstanceConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigGceClusterConfigShieldedInstanceConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigGceClusterConfigShieldedInstanceConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigGceClusterConfigShieldedInstanceConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigGceClusterConfigShieldedInstanceConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigGceClusterConfigShieldedInstanceConfig(c *Client, des, nw *ClusterConfigGceClusterConfigShieldedInstanceConfig) *ClusterConfigGceClusterConfigShieldedInstanceConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigGceClusterConfigShieldedInstanceConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.EnableSecureBoot, nw.EnableSecureBoot) {
		nw.EnableSecureBoot = des.EnableSecureBoot
	}
	if dcl.BoolCanonicalize(des.EnableVtpm, nw.EnableVtpm) {
		nw.EnableVtpm = des.EnableVtpm
	}
	if dcl.BoolCanonicalize(des.EnableIntegrityMonitoring, nw.EnableIntegrityMonitoring) {
		nw.EnableIntegrityMonitoring = des.EnableIntegrityMonitoring
	}

	return nw
}

func canonicalizeNewClusterConfigGceClusterConfigShieldedInstanceConfigSet(c *Client, des, nw []ClusterConfigGceClusterConfigShieldedInstanceConfig) []ClusterConfigGceClusterConfigShieldedInstanceConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigGceClusterConfigShieldedInstanceConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigGceClusterConfigShieldedInstanceConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigGceClusterConfigShieldedInstanceConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigGceClusterConfigShieldedInstanceConfigSlice(c *Client, des, nw []ClusterConfigGceClusterConfigShieldedInstanceConfig) []ClusterConfigGceClusterConfigShieldedInstanceConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigGceClusterConfigShieldedInstanceConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigGceClusterConfigShieldedInstanceConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigGceClusterConfigConfidentialInstanceConfig(des, initial *ClusterConfigGceClusterConfigConfidentialInstanceConfig, opts ...dcl.ApplyOption) *ClusterConfigGceClusterConfigConfidentialInstanceConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigGceClusterConfigConfidentialInstanceConfig{}

	if dcl.BoolCanonicalize(des.EnableConfidentialCompute, initial.EnableConfidentialCompute) || dcl.IsZeroValue(des.EnableConfidentialCompute) {
		cDes.EnableConfidentialCompute = initial.EnableConfidentialCompute
	} else {
		cDes.EnableConfidentialCompute = des.EnableConfidentialCompute
	}

	return cDes
}

func canonicalizeClusterConfigGceClusterConfigConfidentialInstanceConfigSlice(des, initial []ClusterConfigGceClusterConfigConfidentialInstanceConfig, opts ...dcl.ApplyOption) []ClusterConfigGceClusterConfigConfidentialInstanceConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigGceClusterConfigConfidentialInstanceConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigGceClusterConfigConfidentialInstanceConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigGceClusterConfigConfidentialInstanceConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigGceClusterConfigConfidentialInstanceConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigGceClusterConfigConfidentialInstanceConfig(c *Client, des, nw *ClusterConfigGceClusterConfigConfidentialInstanceConfig) *ClusterConfigGceClusterConfigConfidentialInstanceConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigGceClusterConfigConfidentialInstanceConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.EnableConfidentialCompute, nw.EnableConfidentialCompute) {
		nw.EnableConfidentialCompute = des.EnableConfidentialCompute
	}

	return nw
}

func canonicalizeNewClusterConfigGceClusterConfigConfidentialInstanceConfigSet(c *Client, des, nw []ClusterConfigGceClusterConfigConfidentialInstanceConfig) []ClusterConfigGceClusterConfigConfidentialInstanceConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigGceClusterConfigConfidentialInstanceConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigGceClusterConfigConfidentialInstanceConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigGceClusterConfigConfidentialInstanceConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigGceClusterConfigConfidentialInstanceConfigSlice(c *Client, des, nw []ClusterConfigGceClusterConfigConfidentialInstanceConfig) []ClusterConfigGceClusterConfigConfidentialInstanceConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigGceClusterConfigConfidentialInstanceConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigGceClusterConfigConfidentialInstanceConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigMasterConfig(des, initial *ClusterConfigMasterConfig, opts ...dcl.ApplyOption) *ClusterConfigMasterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigMasterConfig{}

	if dcl.IsZeroValue(des.NumInstances) || (dcl.IsEmptyValueIndirect(des.NumInstances) && dcl.IsEmptyValueIndirect(initial.NumInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumInstances = initial.NumInstances
	} else {
		cDes.NumInstances = des.NumInstances
	}
	if dcl.IsZeroValue(des.Image) || (dcl.IsEmptyValueIndirect(des.Image) && dcl.IsEmptyValueIndirect(initial.Image)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Image = initial.Image
	} else {
		cDes.Image = des.Image
	}
	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	cDes.DiskConfig = canonicalizeClusterConfigMasterConfigDiskConfig(des.DiskConfig, initial.DiskConfig, opts...)
	if dcl.IsZeroValue(des.Preemptibility) || (dcl.IsEmptyValueIndirect(des.Preemptibility) && dcl.IsEmptyValueIndirect(initial.Preemptibility)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Preemptibility = initial.Preemptibility
	} else {
		cDes.Preemptibility = des.Preemptibility
	}
	cDes.Accelerators = canonicalizeClusterConfigMasterConfigAcceleratorsSlice(des.Accelerators, initial.Accelerators, opts...)
	if dcl.StringCanonicalize(des.MinCpuPlatform, initial.MinCpuPlatform) || dcl.IsZeroValue(des.MinCpuPlatform) {
		cDes.MinCpuPlatform = initial.MinCpuPlatform
	} else {
		cDes.MinCpuPlatform = des.MinCpuPlatform
	}

	return cDes
}

func canonicalizeClusterConfigMasterConfigSlice(des, initial []ClusterConfigMasterConfig, opts ...dcl.ApplyOption) []ClusterConfigMasterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigMasterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigMasterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigMasterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigMasterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigMasterConfig(c *Client, des, nw *ClusterConfigMasterConfig) *ClusterConfigMasterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigMasterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.InstanceNames, nw.InstanceNames) {
		nw.InstanceNames = des.InstanceNames
	}
	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	nw.DiskConfig = canonicalizeNewClusterConfigMasterConfigDiskConfig(c, des.DiskConfig, nw.DiskConfig)
	if dcl.BoolCanonicalize(des.IsPreemptible, nw.IsPreemptible) {
		nw.IsPreemptible = des.IsPreemptible
	}
	nw.ManagedGroupConfig = canonicalizeNewClusterConfigMasterConfigManagedGroupConfig(c, des.ManagedGroupConfig, nw.ManagedGroupConfig)
	nw.Accelerators = canonicalizeNewClusterConfigMasterConfigAcceleratorsSlice(c, des.Accelerators, nw.Accelerators)
	if dcl.StringCanonicalize(des.MinCpuPlatform, nw.MinCpuPlatform) {
		nw.MinCpuPlatform = des.MinCpuPlatform
	}
	nw.InstanceReferences = canonicalizeNewClusterConfigMasterConfigInstanceReferencesSlice(c, des.InstanceReferences, nw.InstanceReferences)

	return nw
}

func canonicalizeNewClusterConfigMasterConfigSet(c *Client, des, nw []ClusterConfigMasterConfig) []ClusterConfigMasterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigMasterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigMasterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigMasterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigMasterConfigSlice(c *Client, des, nw []ClusterConfigMasterConfig) []ClusterConfigMasterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigMasterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigMasterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigMasterConfigDiskConfig(des, initial *ClusterConfigMasterConfigDiskConfig, opts ...dcl.ApplyOption) *ClusterConfigMasterConfigDiskConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigMasterConfigDiskConfig{}

	if dcl.StringCanonicalize(des.BootDiskType, initial.BootDiskType) || dcl.IsZeroValue(des.BootDiskType) {
		cDes.BootDiskType = initial.BootDiskType
	} else {
		cDes.BootDiskType = des.BootDiskType
	}
	if dcl.IsZeroValue(des.BootDiskSizeGb) || (dcl.IsEmptyValueIndirect(des.BootDiskSizeGb) && dcl.IsEmptyValueIndirect(initial.BootDiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.BootDiskSizeGb = initial.BootDiskSizeGb
	} else {
		cDes.BootDiskSizeGb = des.BootDiskSizeGb
	}
	if dcl.IsZeroValue(des.NumLocalSsds) || (dcl.IsEmptyValueIndirect(des.NumLocalSsds) && dcl.IsEmptyValueIndirect(initial.NumLocalSsds)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumLocalSsds = initial.NumLocalSsds
	} else {
		cDes.NumLocalSsds = des.NumLocalSsds
	}
	if dcl.StringCanonicalize(des.LocalSsdInterface, initial.LocalSsdInterface) || dcl.IsZeroValue(des.LocalSsdInterface) {
		cDes.LocalSsdInterface = initial.LocalSsdInterface
	} else {
		cDes.LocalSsdInterface = des.LocalSsdInterface
	}

	return cDes
}

func canonicalizeClusterConfigMasterConfigDiskConfigSlice(des, initial []ClusterConfigMasterConfigDiskConfig, opts ...dcl.ApplyOption) []ClusterConfigMasterConfigDiskConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigMasterConfigDiskConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigMasterConfigDiskConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigMasterConfigDiskConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigMasterConfigDiskConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigMasterConfigDiskConfig(c *Client, des, nw *ClusterConfigMasterConfigDiskConfig) *ClusterConfigMasterConfigDiskConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigMasterConfigDiskConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.BootDiskType, nw.BootDiskType) {
		nw.BootDiskType = des.BootDiskType
	}
	if dcl.StringCanonicalize(des.LocalSsdInterface, nw.LocalSsdInterface) {
		nw.LocalSsdInterface = des.LocalSsdInterface
	}

	return nw
}

func canonicalizeNewClusterConfigMasterConfigDiskConfigSet(c *Client, des, nw []ClusterConfigMasterConfigDiskConfig) []ClusterConfigMasterConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigMasterConfigDiskConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigMasterConfigDiskConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigMasterConfigDiskConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigMasterConfigDiskConfigSlice(c *Client, des, nw []ClusterConfigMasterConfigDiskConfig) []ClusterConfigMasterConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigMasterConfigDiskConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigMasterConfigDiskConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigMasterConfigManagedGroupConfig(des, initial *ClusterConfigMasterConfigManagedGroupConfig, opts ...dcl.ApplyOption) *ClusterConfigMasterConfigManagedGroupConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigMasterConfigManagedGroupConfig{}

	return cDes
}

func canonicalizeClusterConfigMasterConfigManagedGroupConfigSlice(des, initial []ClusterConfigMasterConfigManagedGroupConfig, opts ...dcl.ApplyOption) []ClusterConfigMasterConfigManagedGroupConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigMasterConfigManagedGroupConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigMasterConfigManagedGroupConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigMasterConfigManagedGroupConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigMasterConfigManagedGroupConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigMasterConfigManagedGroupConfig(c *Client, des, nw *ClusterConfigMasterConfigManagedGroupConfig) *ClusterConfigMasterConfigManagedGroupConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigMasterConfigManagedGroupConfig while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewClusterConfigMasterConfigManagedGroupConfigSet(c *Client, des, nw []ClusterConfigMasterConfigManagedGroupConfig) []ClusterConfigMasterConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigMasterConfigManagedGroupConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigMasterConfigManagedGroupConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigMasterConfigManagedGroupConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigMasterConfigManagedGroupConfigSlice(c *Client, des, nw []ClusterConfigMasterConfigManagedGroupConfig) []ClusterConfigMasterConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigMasterConfigManagedGroupConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigMasterConfigManagedGroupConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigMasterConfigAccelerators(des, initial *ClusterConfigMasterConfigAccelerators, opts ...dcl.ApplyOption) *ClusterConfigMasterConfigAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigMasterConfigAccelerators{}

	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}
	if dcl.IsZeroValue(des.AcceleratorCount) || (dcl.IsEmptyValueIndirect(des.AcceleratorCount) && dcl.IsEmptyValueIndirect(initial.AcceleratorCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}

	return cDes
}

func canonicalizeClusterConfigMasterConfigAcceleratorsSlice(des, initial []ClusterConfigMasterConfigAccelerators, opts ...dcl.ApplyOption) []ClusterConfigMasterConfigAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigMasterConfigAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigMasterConfigAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigMasterConfigAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigMasterConfigAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigMasterConfigAccelerators(c *Client, des, nw *ClusterConfigMasterConfigAccelerators) *ClusterConfigMasterConfigAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigMasterConfigAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}

	return nw
}

func canonicalizeNewClusterConfigMasterConfigAcceleratorsSet(c *Client, des, nw []ClusterConfigMasterConfigAccelerators) []ClusterConfigMasterConfigAccelerators {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigMasterConfigAccelerators
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigMasterConfigAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigMasterConfigAccelerators(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigMasterConfigAcceleratorsSlice(c *Client, des, nw []ClusterConfigMasterConfigAccelerators) []ClusterConfigMasterConfigAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigMasterConfigAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigMasterConfigAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigMasterConfigInstanceReferences(des, initial *ClusterConfigMasterConfigInstanceReferences, opts ...dcl.ApplyOption) *ClusterConfigMasterConfigInstanceReferences {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigMasterConfigInstanceReferences{}

	if dcl.StringCanonicalize(des.InstanceName, initial.InstanceName) || dcl.IsZeroValue(des.InstanceName) {
		cDes.InstanceName = initial.InstanceName
	} else {
		cDes.InstanceName = des.InstanceName
	}
	if dcl.StringCanonicalize(des.InstanceId, initial.InstanceId) || dcl.IsZeroValue(des.InstanceId) {
		cDes.InstanceId = initial.InstanceId
	} else {
		cDes.InstanceId = des.InstanceId
	}
	if dcl.StringCanonicalize(des.PublicKey, initial.PublicKey) || dcl.IsZeroValue(des.PublicKey) {
		cDes.PublicKey = initial.PublicKey
	} else {
		cDes.PublicKey = des.PublicKey
	}
	if dcl.StringCanonicalize(des.PublicEciesKey, initial.PublicEciesKey) || dcl.IsZeroValue(des.PublicEciesKey) {
		cDes.PublicEciesKey = initial.PublicEciesKey
	} else {
		cDes.PublicEciesKey = des.PublicEciesKey
	}

	return cDes
}

func canonicalizeClusterConfigMasterConfigInstanceReferencesSlice(des, initial []ClusterConfigMasterConfigInstanceReferences, opts ...dcl.ApplyOption) []ClusterConfigMasterConfigInstanceReferences {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigMasterConfigInstanceReferences, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigMasterConfigInstanceReferences(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigMasterConfigInstanceReferences, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigMasterConfigInstanceReferences(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigMasterConfigInstanceReferences(c *Client, des, nw *ClusterConfigMasterConfigInstanceReferences) *ClusterConfigMasterConfigInstanceReferences {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigMasterConfigInstanceReferences while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.InstanceName, nw.InstanceName) {
		nw.InstanceName = des.InstanceName
	}
	if dcl.StringCanonicalize(des.InstanceId, nw.InstanceId) {
		nw.InstanceId = des.InstanceId
	}
	if dcl.StringCanonicalize(des.PublicKey, nw.PublicKey) {
		nw.PublicKey = des.PublicKey
	}
	if dcl.StringCanonicalize(des.PublicEciesKey, nw.PublicEciesKey) {
		nw.PublicEciesKey = des.PublicEciesKey
	}

	return nw
}

func canonicalizeNewClusterConfigMasterConfigInstanceReferencesSet(c *Client, des, nw []ClusterConfigMasterConfigInstanceReferences) []ClusterConfigMasterConfigInstanceReferences {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigMasterConfigInstanceReferences
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigMasterConfigInstanceReferencesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigMasterConfigInstanceReferences(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigMasterConfigInstanceReferencesSlice(c *Client, des, nw []ClusterConfigMasterConfigInstanceReferences) []ClusterConfigMasterConfigInstanceReferences {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigMasterConfigInstanceReferences
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigMasterConfigInstanceReferences(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigWorkerConfig(des, initial *ClusterConfigWorkerConfig, opts ...dcl.ApplyOption) *ClusterConfigWorkerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigWorkerConfig{}

	if dcl.IsZeroValue(des.NumInstances) || (dcl.IsEmptyValueIndirect(des.NumInstances) && dcl.IsEmptyValueIndirect(initial.NumInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumInstances = initial.NumInstances
	} else {
		cDes.NumInstances = des.NumInstances
	}
	if dcl.IsZeroValue(des.Image) || (dcl.IsEmptyValueIndirect(des.Image) && dcl.IsEmptyValueIndirect(initial.Image)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Image = initial.Image
	} else {
		cDes.Image = des.Image
	}
	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	cDes.DiskConfig = canonicalizeClusterConfigWorkerConfigDiskConfig(des.DiskConfig, initial.DiskConfig, opts...)
	if dcl.IsZeroValue(des.Preemptibility) || (dcl.IsEmptyValueIndirect(des.Preemptibility) && dcl.IsEmptyValueIndirect(initial.Preemptibility)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Preemptibility = initial.Preemptibility
	} else {
		cDes.Preemptibility = des.Preemptibility
	}
	cDes.Accelerators = canonicalizeClusterConfigWorkerConfigAcceleratorsSlice(des.Accelerators, initial.Accelerators, opts...)
	if dcl.StringCanonicalize(des.MinCpuPlatform, initial.MinCpuPlatform) || dcl.IsZeroValue(des.MinCpuPlatform) {
		cDes.MinCpuPlatform = initial.MinCpuPlatform
	} else {
		cDes.MinCpuPlatform = des.MinCpuPlatform
	}

	return cDes
}

func canonicalizeClusterConfigWorkerConfigSlice(des, initial []ClusterConfigWorkerConfig, opts ...dcl.ApplyOption) []ClusterConfigWorkerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigWorkerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigWorkerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigWorkerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigWorkerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigWorkerConfig(c *Client, des, nw *ClusterConfigWorkerConfig) *ClusterConfigWorkerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigWorkerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.InstanceNames, nw.InstanceNames) {
		nw.InstanceNames = des.InstanceNames
	}
	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	nw.DiskConfig = canonicalizeNewClusterConfigWorkerConfigDiskConfig(c, des.DiskConfig, nw.DiskConfig)
	if dcl.BoolCanonicalize(des.IsPreemptible, nw.IsPreemptible) {
		nw.IsPreemptible = des.IsPreemptible
	}
	nw.ManagedGroupConfig = canonicalizeNewClusterConfigWorkerConfigManagedGroupConfig(c, des.ManagedGroupConfig, nw.ManagedGroupConfig)
	nw.Accelerators = canonicalizeNewClusterConfigWorkerConfigAcceleratorsSlice(c, des.Accelerators, nw.Accelerators)
	if dcl.StringCanonicalize(des.MinCpuPlatform, nw.MinCpuPlatform) {
		nw.MinCpuPlatform = des.MinCpuPlatform
	}
	nw.InstanceReferences = canonicalizeNewClusterConfigWorkerConfigInstanceReferencesSlice(c, des.InstanceReferences, nw.InstanceReferences)

	return nw
}

func canonicalizeNewClusterConfigWorkerConfigSet(c *Client, des, nw []ClusterConfigWorkerConfig) []ClusterConfigWorkerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigWorkerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigWorkerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigWorkerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigWorkerConfigSlice(c *Client, des, nw []ClusterConfigWorkerConfig) []ClusterConfigWorkerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigWorkerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigWorkerConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigWorkerConfigDiskConfig(des, initial *ClusterConfigWorkerConfigDiskConfig, opts ...dcl.ApplyOption) *ClusterConfigWorkerConfigDiskConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigWorkerConfigDiskConfig{}

	if dcl.StringCanonicalize(des.BootDiskType, initial.BootDiskType) || dcl.IsZeroValue(des.BootDiskType) {
		cDes.BootDiskType = initial.BootDiskType
	} else {
		cDes.BootDiskType = des.BootDiskType
	}
	if dcl.IsZeroValue(des.BootDiskSizeGb) || (dcl.IsEmptyValueIndirect(des.BootDiskSizeGb) && dcl.IsEmptyValueIndirect(initial.BootDiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.BootDiskSizeGb = initial.BootDiskSizeGb
	} else {
		cDes.BootDiskSizeGb = des.BootDiskSizeGb
	}
	if dcl.IsZeroValue(des.NumLocalSsds) || (dcl.IsEmptyValueIndirect(des.NumLocalSsds) && dcl.IsEmptyValueIndirect(initial.NumLocalSsds)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumLocalSsds = initial.NumLocalSsds
	} else {
		cDes.NumLocalSsds = des.NumLocalSsds
	}
	if dcl.StringCanonicalize(des.LocalSsdInterface, initial.LocalSsdInterface) || dcl.IsZeroValue(des.LocalSsdInterface) {
		cDes.LocalSsdInterface = initial.LocalSsdInterface
	} else {
		cDes.LocalSsdInterface = des.LocalSsdInterface
	}

	return cDes
}

func canonicalizeClusterConfigWorkerConfigDiskConfigSlice(des, initial []ClusterConfigWorkerConfigDiskConfig, opts ...dcl.ApplyOption) []ClusterConfigWorkerConfigDiskConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigWorkerConfigDiskConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigWorkerConfigDiskConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigWorkerConfigDiskConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigWorkerConfigDiskConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigWorkerConfigDiskConfig(c *Client, des, nw *ClusterConfigWorkerConfigDiskConfig) *ClusterConfigWorkerConfigDiskConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigWorkerConfigDiskConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.BootDiskType, nw.BootDiskType) {
		nw.BootDiskType = des.BootDiskType
	}
	if dcl.StringCanonicalize(des.LocalSsdInterface, nw.LocalSsdInterface) {
		nw.LocalSsdInterface = des.LocalSsdInterface
	}

	return nw
}

func canonicalizeNewClusterConfigWorkerConfigDiskConfigSet(c *Client, des, nw []ClusterConfigWorkerConfigDiskConfig) []ClusterConfigWorkerConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigWorkerConfigDiskConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigWorkerConfigDiskConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigWorkerConfigDiskConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigWorkerConfigDiskConfigSlice(c *Client, des, nw []ClusterConfigWorkerConfigDiskConfig) []ClusterConfigWorkerConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigWorkerConfigDiskConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigWorkerConfigDiskConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigWorkerConfigManagedGroupConfig(des, initial *ClusterConfigWorkerConfigManagedGroupConfig, opts ...dcl.ApplyOption) *ClusterConfigWorkerConfigManagedGroupConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigWorkerConfigManagedGroupConfig{}

	return cDes
}

func canonicalizeClusterConfigWorkerConfigManagedGroupConfigSlice(des, initial []ClusterConfigWorkerConfigManagedGroupConfig, opts ...dcl.ApplyOption) []ClusterConfigWorkerConfigManagedGroupConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigWorkerConfigManagedGroupConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigWorkerConfigManagedGroupConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigWorkerConfigManagedGroupConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigWorkerConfigManagedGroupConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigWorkerConfigManagedGroupConfig(c *Client, des, nw *ClusterConfigWorkerConfigManagedGroupConfig) *ClusterConfigWorkerConfigManagedGroupConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigWorkerConfigManagedGroupConfig while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewClusterConfigWorkerConfigManagedGroupConfigSet(c *Client, des, nw []ClusterConfigWorkerConfigManagedGroupConfig) []ClusterConfigWorkerConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigWorkerConfigManagedGroupConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigWorkerConfigManagedGroupConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigWorkerConfigManagedGroupConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigWorkerConfigManagedGroupConfigSlice(c *Client, des, nw []ClusterConfigWorkerConfigManagedGroupConfig) []ClusterConfigWorkerConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigWorkerConfigManagedGroupConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigWorkerConfigManagedGroupConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigWorkerConfigAccelerators(des, initial *ClusterConfigWorkerConfigAccelerators, opts ...dcl.ApplyOption) *ClusterConfigWorkerConfigAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigWorkerConfigAccelerators{}

	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}
	if dcl.IsZeroValue(des.AcceleratorCount) || (dcl.IsEmptyValueIndirect(des.AcceleratorCount) && dcl.IsEmptyValueIndirect(initial.AcceleratorCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}

	return cDes
}

func canonicalizeClusterConfigWorkerConfigAcceleratorsSlice(des, initial []ClusterConfigWorkerConfigAccelerators, opts ...dcl.ApplyOption) []ClusterConfigWorkerConfigAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigWorkerConfigAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigWorkerConfigAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigWorkerConfigAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigWorkerConfigAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigWorkerConfigAccelerators(c *Client, des, nw *ClusterConfigWorkerConfigAccelerators) *ClusterConfigWorkerConfigAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigWorkerConfigAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}

	return nw
}

func canonicalizeNewClusterConfigWorkerConfigAcceleratorsSet(c *Client, des, nw []ClusterConfigWorkerConfigAccelerators) []ClusterConfigWorkerConfigAccelerators {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigWorkerConfigAccelerators
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigWorkerConfigAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigWorkerConfigAccelerators(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigWorkerConfigAcceleratorsSlice(c *Client, des, nw []ClusterConfigWorkerConfigAccelerators) []ClusterConfigWorkerConfigAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigWorkerConfigAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigWorkerConfigAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigWorkerConfigInstanceReferences(des, initial *ClusterConfigWorkerConfigInstanceReferences, opts ...dcl.ApplyOption) *ClusterConfigWorkerConfigInstanceReferences {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigWorkerConfigInstanceReferences{}

	if dcl.StringCanonicalize(des.InstanceName, initial.InstanceName) || dcl.IsZeroValue(des.InstanceName) {
		cDes.InstanceName = initial.InstanceName
	} else {
		cDes.InstanceName = des.InstanceName
	}
	if dcl.StringCanonicalize(des.InstanceId, initial.InstanceId) || dcl.IsZeroValue(des.InstanceId) {
		cDes.InstanceId = initial.InstanceId
	} else {
		cDes.InstanceId = des.InstanceId
	}
	if dcl.StringCanonicalize(des.PublicKey, initial.PublicKey) || dcl.IsZeroValue(des.PublicKey) {
		cDes.PublicKey = initial.PublicKey
	} else {
		cDes.PublicKey = des.PublicKey
	}
	if dcl.StringCanonicalize(des.PublicEciesKey, initial.PublicEciesKey) || dcl.IsZeroValue(des.PublicEciesKey) {
		cDes.PublicEciesKey = initial.PublicEciesKey
	} else {
		cDes.PublicEciesKey = des.PublicEciesKey
	}

	return cDes
}

func canonicalizeClusterConfigWorkerConfigInstanceReferencesSlice(des, initial []ClusterConfigWorkerConfigInstanceReferences, opts ...dcl.ApplyOption) []ClusterConfigWorkerConfigInstanceReferences {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigWorkerConfigInstanceReferences, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigWorkerConfigInstanceReferences(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigWorkerConfigInstanceReferences, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigWorkerConfigInstanceReferences(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigWorkerConfigInstanceReferences(c *Client, des, nw *ClusterConfigWorkerConfigInstanceReferences) *ClusterConfigWorkerConfigInstanceReferences {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigWorkerConfigInstanceReferences while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.InstanceName, nw.InstanceName) {
		nw.InstanceName = des.InstanceName
	}
	if dcl.StringCanonicalize(des.InstanceId, nw.InstanceId) {
		nw.InstanceId = des.InstanceId
	}
	if dcl.StringCanonicalize(des.PublicKey, nw.PublicKey) {
		nw.PublicKey = des.PublicKey
	}
	if dcl.StringCanonicalize(des.PublicEciesKey, nw.PublicEciesKey) {
		nw.PublicEciesKey = des.PublicEciesKey
	}

	return nw
}

func canonicalizeNewClusterConfigWorkerConfigInstanceReferencesSet(c *Client, des, nw []ClusterConfigWorkerConfigInstanceReferences) []ClusterConfigWorkerConfigInstanceReferences {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigWorkerConfigInstanceReferences
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigWorkerConfigInstanceReferencesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigWorkerConfigInstanceReferences(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigWorkerConfigInstanceReferencesSlice(c *Client, des, nw []ClusterConfigWorkerConfigInstanceReferences) []ClusterConfigWorkerConfigInstanceReferences {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigWorkerConfigInstanceReferences
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigWorkerConfigInstanceReferences(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSecondaryWorkerConfig(des, initial *ClusterConfigSecondaryWorkerConfig, opts ...dcl.ApplyOption) *ClusterConfigSecondaryWorkerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSecondaryWorkerConfig{}

	if dcl.IsZeroValue(des.NumInstances) || (dcl.IsEmptyValueIndirect(des.NumInstances) && dcl.IsEmptyValueIndirect(initial.NumInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumInstances = initial.NumInstances
	} else {
		cDes.NumInstances = des.NumInstances
	}
	if dcl.IsZeroValue(des.Image) || (dcl.IsEmptyValueIndirect(des.Image) && dcl.IsEmptyValueIndirect(initial.Image)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Image = initial.Image
	} else {
		cDes.Image = des.Image
	}
	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	cDes.DiskConfig = canonicalizeClusterConfigSecondaryWorkerConfigDiskConfig(des.DiskConfig, initial.DiskConfig, opts...)
	if dcl.IsZeroValue(des.Preemptibility) || (dcl.IsEmptyValueIndirect(des.Preemptibility) && dcl.IsEmptyValueIndirect(initial.Preemptibility)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Preemptibility = initial.Preemptibility
	} else {
		cDes.Preemptibility = des.Preemptibility
	}
	cDes.Accelerators = canonicalizeClusterConfigSecondaryWorkerConfigAcceleratorsSlice(des.Accelerators, initial.Accelerators, opts...)
	if dcl.StringCanonicalize(des.MinCpuPlatform, initial.MinCpuPlatform) || dcl.IsZeroValue(des.MinCpuPlatform) {
		cDes.MinCpuPlatform = initial.MinCpuPlatform
	} else {
		cDes.MinCpuPlatform = des.MinCpuPlatform
	}

	return cDes
}

func canonicalizeClusterConfigSecondaryWorkerConfigSlice(des, initial []ClusterConfigSecondaryWorkerConfig, opts ...dcl.ApplyOption) []ClusterConfigSecondaryWorkerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSecondaryWorkerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSecondaryWorkerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSecondaryWorkerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSecondaryWorkerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSecondaryWorkerConfig(c *Client, des, nw *ClusterConfigSecondaryWorkerConfig) *ClusterConfigSecondaryWorkerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSecondaryWorkerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.InstanceNames, nw.InstanceNames) {
		nw.InstanceNames = des.InstanceNames
	}
	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	nw.DiskConfig = canonicalizeNewClusterConfigSecondaryWorkerConfigDiskConfig(c, des.DiskConfig, nw.DiskConfig)
	if dcl.BoolCanonicalize(des.IsPreemptible, nw.IsPreemptible) {
		nw.IsPreemptible = des.IsPreemptible
	}
	nw.ManagedGroupConfig = canonicalizeNewClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, des.ManagedGroupConfig, nw.ManagedGroupConfig)
	nw.Accelerators = canonicalizeNewClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c, des.Accelerators, nw.Accelerators)
	if dcl.StringCanonicalize(des.MinCpuPlatform, nw.MinCpuPlatform) {
		nw.MinCpuPlatform = des.MinCpuPlatform
	}
	nw.InstanceReferences = canonicalizeNewClusterConfigSecondaryWorkerConfigInstanceReferencesSlice(c, des.InstanceReferences, nw.InstanceReferences)

	return nw
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigSet(c *Client, des, nw []ClusterConfigSecondaryWorkerConfig) []ClusterConfigSecondaryWorkerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSecondaryWorkerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSecondaryWorkerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigSlice(c *Client, des, nw []ClusterConfigSecondaryWorkerConfig) []ClusterConfigSecondaryWorkerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSecondaryWorkerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSecondaryWorkerConfigDiskConfig(des, initial *ClusterConfigSecondaryWorkerConfigDiskConfig, opts ...dcl.ApplyOption) *ClusterConfigSecondaryWorkerConfigDiskConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSecondaryWorkerConfigDiskConfig{}

	if dcl.StringCanonicalize(des.BootDiskType, initial.BootDiskType) || dcl.IsZeroValue(des.BootDiskType) {
		cDes.BootDiskType = initial.BootDiskType
	} else {
		cDes.BootDiskType = des.BootDiskType
	}
	if dcl.IsZeroValue(des.BootDiskSizeGb) || (dcl.IsEmptyValueIndirect(des.BootDiskSizeGb) && dcl.IsEmptyValueIndirect(initial.BootDiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.BootDiskSizeGb = initial.BootDiskSizeGb
	} else {
		cDes.BootDiskSizeGb = des.BootDiskSizeGb
	}
	if dcl.IsZeroValue(des.NumLocalSsds) || (dcl.IsEmptyValueIndirect(des.NumLocalSsds) && dcl.IsEmptyValueIndirect(initial.NumLocalSsds)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumLocalSsds = initial.NumLocalSsds
	} else {
		cDes.NumLocalSsds = des.NumLocalSsds
	}
	if dcl.StringCanonicalize(des.LocalSsdInterface, initial.LocalSsdInterface) || dcl.IsZeroValue(des.LocalSsdInterface) {
		cDes.LocalSsdInterface = initial.LocalSsdInterface
	} else {
		cDes.LocalSsdInterface = des.LocalSsdInterface
	}

	return cDes
}

func canonicalizeClusterConfigSecondaryWorkerConfigDiskConfigSlice(des, initial []ClusterConfigSecondaryWorkerConfigDiskConfig, opts ...dcl.ApplyOption) []ClusterConfigSecondaryWorkerConfigDiskConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSecondaryWorkerConfigDiskConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSecondaryWorkerConfigDiskConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSecondaryWorkerConfigDiskConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSecondaryWorkerConfigDiskConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSecondaryWorkerConfigDiskConfig(c *Client, des, nw *ClusterConfigSecondaryWorkerConfigDiskConfig) *ClusterConfigSecondaryWorkerConfigDiskConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSecondaryWorkerConfigDiskConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.BootDiskType, nw.BootDiskType) {
		nw.BootDiskType = des.BootDiskType
	}
	if dcl.StringCanonicalize(des.LocalSsdInterface, nw.LocalSsdInterface) {
		nw.LocalSsdInterface = des.LocalSsdInterface
	}

	return nw
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigDiskConfigSet(c *Client, des, nw []ClusterConfigSecondaryWorkerConfigDiskConfig) []ClusterConfigSecondaryWorkerConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSecondaryWorkerConfigDiskConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSecondaryWorkerConfigDiskConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfigDiskConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigDiskConfigSlice(c *Client, des, nw []ClusterConfigSecondaryWorkerConfigDiskConfig) []ClusterConfigSecondaryWorkerConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSecondaryWorkerConfigDiskConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfigDiskConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSecondaryWorkerConfigManagedGroupConfig(des, initial *ClusterConfigSecondaryWorkerConfigManagedGroupConfig, opts ...dcl.ApplyOption) *ClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSecondaryWorkerConfigManagedGroupConfig{}

	return cDes
}

func canonicalizeClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice(des, initial []ClusterConfigSecondaryWorkerConfigManagedGroupConfig, opts ...dcl.ApplyOption) []ClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSecondaryWorkerConfigManagedGroupConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSecondaryWorkerConfigManagedGroupConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSecondaryWorkerConfigManagedGroupConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSecondaryWorkerConfigManagedGroupConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSecondaryWorkerConfigManagedGroupConfig(c *Client, des, nw *ClusterConfigSecondaryWorkerConfigManagedGroupConfig) *ClusterConfigSecondaryWorkerConfigManagedGroupConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSecondaryWorkerConfigManagedGroupConfig while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewClusterConfigSecondaryWorkerConfigManagedGroupConfigSet(c *Client, des, nw []ClusterConfigSecondaryWorkerConfigManagedGroupConfig) []ClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSecondaryWorkerConfigManagedGroupConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSecondaryWorkerConfigManagedGroupConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice(c *Client, des, nw []ClusterConfigSecondaryWorkerConfigManagedGroupConfig) []ClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSecondaryWorkerConfigManagedGroupConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSecondaryWorkerConfigAccelerators(des, initial *ClusterConfigSecondaryWorkerConfigAccelerators, opts ...dcl.ApplyOption) *ClusterConfigSecondaryWorkerConfigAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSecondaryWorkerConfigAccelerators{}

	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}
	if dcl.IsZeroValue(des.AcceleratorCount) || (dcl.IsEmptyValueIndirect(des.AcceleratorCount) && dcl.IsEmptyValueIndirect(initial.AcceleratorCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}

	return cDes
}

func canonicalizeClusterConfigSecondaryWorkerConfigAcceleratorsSlice(des, initial []ClusterConfigSecondaryWorkerConfigAccelerators, opts ...dcl.ApplyOption) []ClusterConfigSecondaryWorkerConfigAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSecondaryWorkerConfigAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSecondaryWorkerConfigAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSecondaryWorkerConfigAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSecondaryWorkerConfigAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSecondaryWorkerConfigAccelerators(c *Client, des, nw *ClusterConfigSecondaryWorkerConfigAccelerators) *ClusterConfigSecondaryWorkerConfigAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSecondaryWorkerConfigAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}

	return nw
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigAcceleratorsSet(c *Client, des, nw []ClusterConfigSecondaryWorkerConfigAccelerators) []ClusterConfigSecondaryWorkerConfigAccelerators {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSecondaryWorkerConfigAccelerators
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSecondaryWorkerConfigAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfigAccelerators(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c *Client, des, nw []ClusterConfigSecondaryWorkerConfigAccelerators) []ClusterConfigSecondaryWorkerConfigAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSecondaryWorkerConfigAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfigAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSecondaryWorkerConfigInstanceReferences(des, initial *ClusterConfigSecondaryWorkerConfigInstanceReferences, opts ...dcl.ApplyOption) *ClusterConfigSecondaryWorkerConfigInstanceReferences {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSecondaryWorkerConfigInstanceReferences{}

	if dcl.StringCanonicalize(des.InstanceName, initial.InstanceName) || dcl.IsZeroValue(des.InstanceName) {
		cDes.InstanceName = initial.InstanceName
	} else {
		cDes.InstanceName = des.InstanceName
	}
	if dcl.StringCanonicalize(des.InstanceId, initial.InstanceId) || dcl.IsZeroValue(des.InstanceId) {
		cDes.InstanceId = initial.InstanceId
	} else {
		cDes.InstanceId = des.InstanceId
	}
	if dcl.StringCanonicalize(des.PublicKey, initial.PublicKey) || dcl.IsZeroValue(des.PublicKey) {
		cDes.PublicKey = initial.PublicKey
	} else {
		cDes.PublicKey = des.PublicKey
	}
	if dcl.StringCanonicalize(des.PublicEciesKey, initial.PublicEciesKey) || dcl.IsZeroValue(des.PublicEciesKey) {
		cDes.PublicEciesKey = initial.PublicEciesKey
	} else {
		cDes.PublicEciesKey = des.PublicEciesKey
	}

	return cDes
}

func canonicalizeClusterConfigSecondaryWorkerConfigInstanceReferencesSlice(des, initial []ClusterConfigSecondaryWorkerConfigInstanceReferences, opts ...dcl.ApplyOption) []ClusterConfigSecondaryWorkerConfigInstanceReferences {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSecondaryWorkerConfigInstanceReferences, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSecondaryWorkerConfigInstanceReferences(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSecondaryWorkerConfigInstanceReferences, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSecondaryWorkerConfigInstanceReferences(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSecondaryWorkerConfigInstanceReferences(c *Client, des, nw *ClusterConfigSecondaryWorkerConfigInstanceReferences) *ClusterConfigSecondaryWorkerConfigInstanceReferences {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSecondaryWorkerConfigInstanceReferences while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.InstanceName, nw.InstanceName) {
		nw.InstanceName = des.InstanceName
	}
	if dcl.StringCanonicalize(des.InstanceId, nw.InstanceId) {
		nw.InstanceId = des.InstanceId
	}
	if dcl.StringCanonicalize(des.PublicKey, nw.PublicKey) {
		nw.PublicKey = des.PublicKey
	}
	if dcl.StringCanonicalize(des.PublicEciesKey, nw.PublicEciesKey) {
		nw.PublicEciesKey = des.PublicEciesKey
	}

	return nw
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigInstanceReferencesSet(c *Client, des, nw []ClusterConfigSecondaryWorkerConfigInstanceReferences) []ClusterConfigSecondaryWorkerConfigInstanceReferences {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSecondaryWorkerConfigInstanceReferences
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSecondaryWorkerConfigInstanceReferencesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfigInstanceReferences(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSecondaryWorkerConfigInstanceReferencesSlice(c *Client, des, nw []ClusterConfigSecondaryWorkerConfigInstanceReferences) []ClusterConfigSecondaryWorkerConfigInstanceReferences {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSecondaryWorkerConfigInstanceReferences
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSecondaryWorkerConfigInstanceReferences(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSoftwareConfig(des, initial *ClusterConfigSoftwareConfig, opts ...dcl.ApplyOption) *ClusterConfigSoftwareConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSoftwareConfig{}

	if dcl.StringCanonicalize(des.ImageVersion, initial.ImageVersion) || dcl.IsZeroValue(des.ImageVersion) {
		cDes.ImageVersion = initial.ImageVersion
	} else {
		cDes.ImageVersion = des.ImageVersion
	}
	if canonicalizeSoftwareConfigProperties(des.Properties, initial.Properties) || dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.IsZeroValue(des.OptionalComponents) || (dcl.IsEmptyValueIndirect(des.OptionalComponents) && dcl.IsEmptyValueIndirect(initial.OptionalComponents)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.OptionalComponents = initial.OptionalComponents
	} else {
		cDes.OptionalComponents = des.OptionalComponents
	}

	return cDes
}

func canonicalizeClusterConfigSoftwareConfigSlice(des, initial []ClusterConfigSoftwareConfig, opts ...dcl.ApplyOption) []ClusterConfigSoftwareConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSoftwareConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSoftwareConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSoftwareConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSoftwareConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSoftwareConfig(c *Client, des, nw *ClusterConfigSoftwareConfig) *ClusterConfigSoftwareConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSoftwareConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ImageVersion, nw.ImageVersion) {
		nw.ImageVersion = des.ImageVersion
	}
	if canonicalizeSoftwareConfigProperties(des.Properties, nw.Properties) {
		nw.Properties = des.Properties
	}

	return nw
}

func canonicalizeNewClusterConfigSoftwareConfigSet(c *Client, des, nw []ClusterConfigSoftwareConfig) []ClusterConfigSoftwareConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSoftwareConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSoftwareConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSoftwareConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSoftwareConfigSlice(c *Client, des, nw []ClusterConfigSoftwareConfig) []ClusterConfigSoftwareConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSoftwareConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSoftwareConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigInitializationActions(des, initial *ClusterConfigInitializationActions, opts ...dcl.ApplyOption) *ClusterConfigInitializationActions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigInitializationActions{}

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

func canonicalizeClusterConfigInitializationActionsSlice(des, initial []ClusterConfigInitializationActions, opts ...dcl.ApplyOption) []ClusterConfigInitializationActions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigInitializationActions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigInitializationActions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigInitializationActions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigInitializationActions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigInitializationActions(c *Client, des, nw *ClusterConfigInitializationActions) *ClusterConfigInitializationActions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigInitializationActions while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewClusterConfigInitializationActionsSet(c *Client, des, nw []ClusterConfigInitializationActions) []ClusterConfigInitializationActions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigInitializationActions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigInitializationActionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigInitializationActions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigInitializationActionsSlice(c *Client, des, nw []ClusterConfigInitializationActions) []ClusterConfigInitializationActions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigInitializationActions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigInitializationActions(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigEncryptionConfig(des, initial *ClusterConfigEncryptionConfig, opts ...dcl.ApplyOption) *ClusterConfigEncryptionConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigEncryptionConfig{}

	if dcl.IsZeroValue(des.GcePdKmsKeyName) || (dcl.IsEmptyValueIndirect(des.GcePdKmsKeyName) && dcl.IsEmptyValueIndirect(initial.GcePdKmsKeyName)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.GcePdKmsKeyName = initial.GcePdKmsKeyName
	} else {
		cDes.GcePdKmsKeyName = des.GcePdKmsKeyName
	}

	return cDes
}

func canonicalizeClusterConfigEncryptionConfigSlice(des, initial []ClusterConfigEncryptionConfig, opts ...dcl.ApplyOption) []ClusterConfigEncryptionConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigEncryptionConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigEncryptionConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigEncryptionConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigEncryptionConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigEncryptionConfig(c *Client, des, nw *ClusterConfigEncryptionConfig) *ClusterConfigEncryptionConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigEncryptionConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterConfigEncryptionConfigSet(c *Client, des, nw []ClusterConfigEncryptionConfig) []ClusterConfigEncryptionConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigEncryptionConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigEncryptionConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigEncryptionConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigEncryptionConfigSlice(c *Client, des, nw []ClusterConfigEncryptionConfig) []ClusterConfigEncryptionConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigEncryptionConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigEncryptionConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigAutoscalingConfig(des, initial *ClusterConfigAutoscalingConfig, opts ...dcl.ApplyOption) *ClusterConfigAutoscalingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigAutoscalingConfig{}

	if dcl.IsZeroValue(des.Policy) || (dcl.IsEmptyValueIndirect(des.Policy) && dcl.IsEmptyValueIndirect(initial.Policy)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Policy = initial.Policy
	} else {
		cDes.Policy = des.Policy
	}

	return cDes
}

func canonicalizeClusterConfigAutoscalingConfigSlice(des, initial []ClusterConfigAutoscalingConfig, opts ...dcl.ApplyOption) []ClusterConfigAutoscalingConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigAutoscalingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigAutoscalingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigAutoscalingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigAutoscalingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigAutoscalingConfig(c *Client, des, nw *ClusterConfigAutoscalingConfig) *ClusterConfigAutoscalingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigAutoscalingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterConfigAutoscalingConfigSet(c *Client, des, nw []ClusterConfigAutoscalingConfig) []ClusterConfigAutoscalingConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigAutoscalingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigAutoscalingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigAutoscalingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigAutoscalingConfigSlice(c *Client, des, nw []ClusterConfigAutoscalingConfig) []ClusterConfigAutoscalingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigAutoscalingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigAutoscalingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSecurityConfig(des, initial *ClusterConfigSecurityConfig, opts ...dcl.ApplyOption) *ClusterConfigSecurityConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSecurityConfig{}

	cDes.KerberosConfig = canonicalizeClusterConfigSecurityConfigKerberosConfig(des.KerberosConfig, initial.KerberosConfig, opts...)
	cDes.IdentityConfig = canonicalizeClusterConfigSecurityConfigIdentityConfig(des.IdentityConfig, initial.IdentityConfig, opts...)

	return cDes
}

func canonicalizeClusterConfigSecurityConfigSlice(des, initial []ClusterConfigSecurityConfig, opts ...dcl.ApplyOption) []ClusterConfigSecurityConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSecurityConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSecurityConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSecurityConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSecurityConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSecurityConfig(c *Client, des, nw *ClusterConfigSecurityConfig) *ClusterConfigSecurityConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSecurityConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KerberosConfig = canonicalizeNewClusterConfigSecurityConfigKerberosConfig(c, des.KerberosConfig, nw.KerberosConfig)
	nw.IdentityConfig = canonicalizeNewClusterConfigSecurityConfigIdentityConfig(c, des.IdentityConfig, nw.IdentityConfig)

	return nw
}

func canonicalizeNewClusterConfigSecurityConfigSet(c *Client, des, nw []ClusterConfigSecurityConfig) []ClusterConfigSecurityConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSecurityConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSecurityConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSecurityConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSecurityConfigSlice(c *Client, des, nw []ClusterConfigSecurityConfig) []ClusterConfigSecurityConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSecurityConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSecurityConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSecurityConfigKerberosConfig(des, initial *ClusterConfigSecurityConfigKerberosConfig, opts ...dcl.ApplyOption) *ClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSecurityConfigKerberosConfig{}

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
	if dcl.IsZeroValue(des.KmsKey) || (dcl.IsEmptyValueIndirect(des.KmsKey) && dcl.IsEmptyValueIndirect(initial.KmsKey)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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
	if dcl.IsZeroValue(des.TgtLifetimeHours) || (dcl.IsEmptyValueIndirect(des.TgtLifetimeHours) && dcl.IsEmptyValueIndirect(initial.TgtLifetimeHours)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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

func canonicalizeClusterConfigSecurityConfigKerberosConfigSlice(des, initial []ClusterConfigSecurityConfigKerberosConfig, opts ...dcl.ApplyOption) []ClusterConfigSecurityConfigKerberosConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSecurityConfigKerberosConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSecurityConfigKerberosConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSecurityConfigKerberosConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSecurityConfigKerberosConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSecurityConfigKerberosConfig(c *Client, des, nw *ClusterConfigSecurityConfigKerberosConfig) *ClusterConfigSecurityConfigKerberosConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSecurityConfigKerberosConfig while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewClusterConfigSecurityConfigKerberosConfigSet(c *Client, des, nw []ClusterConfigSecurityConfigKerberosConfig) []ClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSecurityConfigKerberosConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSecurityConfigKerberosConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSecurityConfigKerberosConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSecurityConfigKerberosConfigSlice(c *Client, des, nw []ClusterConfigSecurityConfigKerberosConfig) []ClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSecurityConfigKerberosConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSecurityConfigKerberosConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigSecurityConfigIdentityConfig(des, initial *ClusterConfigSecurityConfigIdentityConfig, opts ...dcl.ApplyOption) *ClusterConfigSecurityConfigIdentityConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigSecurityConfigIdentityConfig{}

	if dcl.IsZeroValue(des.UserServiceAccountMapping) || (dcl.IsEmptyValueIndirect(des.UserServiceAccountMapping) && dcl.IsEmptyValueIndirect(initial.UserServiceAccountMapping)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.UserServiceAccountMapping = initial.UserServiceAccountMapping
	} else {
		cDes.UserServiceAccountMapping = des.UserServiceAccountMapping
	}

	return cDes
}

func canonicalizeClusterConfigSecurityConfigIdentityConfigSlice(des, initial []ClusterConfigSecurityConfigIdentityConfig, opts ...dcl.ApplyOption) []ClusterConfigSecurityConfigIdentityConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigSecurityConfigIdentityConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigSecurityConfigIdentityConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigSecurityConfigIdentityConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigSecurityConfigIdentityConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigSecurityConfigIdentityConfig(c *Client, des, nw *ClusterConfigSecurityConfigIdentityConfig) *ClusterConfigSecurityConfigIdentityConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigSecurityConfigIdentityConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterConfigSecurityConfigIdentityConfigSet(c *Client, des, nw []ClusterConfigSecurityConfigIdentityConfig) []ClusterConfigSecurityConfigIdentityConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigSecurityConfigIdentityConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigSecurityConfigIdentityConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigSecurityConfigIdentityConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigSecurityConfigIdentityConfigSlice(c *Client, des, nw []ClusterConfigSecurityConfigIdentityConfig) []ClusterConfigSecurityConfigIdentityConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigSecurityConfigIdentityConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigSecurityConfigIdentityConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigLifecycleConfig(des, initial *ClusterConfigLifecycleConfig, opts ...dcl.ApplyOption) *ClusterConfigLifecycleConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigLifecycleConfig{}

	if dcl.StringCanonicalize(des.IdleDeleteTtl, initial.IdleDeleteTtl) || dcl.IsZeroValue(des.IdleDeleteTtl) {
		cDes.IdleDeleteTtl = initial.IdleDeleteTtl
	} else {
		cDes.IdleDeleteTtl = des.IdleDeleteTtl
	}
	if dcl.IsZeroValue(des.AutoDeleteTime) || (dcl.IsEmptyValueIndirect(des.AutoDeleteTime) && dcl.IsEmptyValueIndirect(initial.AutoDeleteTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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

func canonicalizeClusterConfigLifecycleConfigSlice(des, initial []ClusterConfigLifecycleConfig, opts ...dcl.ApplyOption) []ClusterConfigLifecycleConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigLifecycleConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigLifecycleConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigLifecycleConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigLifecycleConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigLifecycleConfig(c *Client, des, nw *ClusterConfigLifecycleConfig) *ClusterConfigLifecycleConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigLifecycleConfig while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewClusterConfigLifecycleConfigSet(c *Client, des, nw []ClusterConfigLifecycleConfig) []ClusterConfigLifecycleConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigLifecycleConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigLifecycleConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigLifecycleConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigLifecycleConfigSlice(c *Client, des, nw []ClusterConfigLifecycleConfig) []ClusterConfigLifecycleConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigLifecycleConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigLifecycleConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigEndpointConfig(des, initial *ClusterConfigEndpointConfig, opts ...dcl.ApplyOption) *ClusterConfigEndpointConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigEndpointConfig{}

	if dcl.BoolCanonicalize(des.EnableHttpPortAccess, initial.EnableHttpPortAccess) || dcl.IsZeroValue(des.EnableHttpPortAccess) {
		cDes.EnableHttpPortAccess = initial.EnableHttpPortAccess
	} else {
		cDes.EnableHttpPortAccess = des.EnableHttpPortAccess
	}

	return cDes
}

func canonicalizeClusterConfigEndpointConfigSlice(des, initial []ClusterConfigEndpointConfig, opts ...dcl.ApplyOption) []ClusterConfigEndpointConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigEndpointConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigEndpointConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigEndpointConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigEndpointConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigEndpointConfig(c *Client, des, nw *ClusterConfigEndpointConfig) *ClusterConfigEndpointConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigEndpointConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.EnableHttpPortAccess, nw.EnableHttpPortAccess) {
		nw.EnableHttpPortAccess = des.EnableHttpPortAccess
	}

	return nw
}

func canonicalizeNewClusterConfigEndpointConfigSet(c *Client, des, nw []ClusterConfigEndpointConfig) []ClusterConfigEndpointConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigEndpointConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigEndpointConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigEndpointConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigEndpointConfigSlice(c *Client, des, nw []ClusterConfigEndpointConfig) []ClusterConfigEndpointConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigEndpointConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigEndpointConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigMetastoreConfig(des, initial *ClusterConfigMetastoreConfig, opts ...dcl.ApplyOption) *ClusterConfigMetastoreConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigMetastoreConfig{}

	if dcl.IsZeroValue(des.DataprocMetastoreService) || (dcl.IsEmptyValueIndirect(des.DataprocMetastoreService) && dcl.IsEmptyValueIndirect(initial.DataprocMetastoreService)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DataprocMetastoreService = initial.DataprocMetastoreService
	} else {
		cDes.DataprocMetastoreService = des.DataprocMetastoreService
	}

	return cDes
}

func canonicalizeClusterConfigMetastoreConfigSlice(des, initial []ClusterConfigMetastoreConfig, opts ...dcl.ApplyOption) []ClusterConfigMetastoreConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigMetastoreConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigMetastoreConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigMetastoreConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigMetastoreConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigMetastoreConfig(c *Client, des, nw *ClusterConfigMetastoreConfig) *ClusterConfigMetastoreConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigMetastoreConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterConfigMetastoreConfigSet(c *Client, des, nw []ClusterConfigMetastoreConfig) []ClusterConfigMetastoreConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigMetastoreConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigMetastoreConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigMetastoreConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigMetastoreConfigSlice(c *Client, des, nw []ClusterConfigMetastoreConfig) []ClusterConfigMetastoreConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigMetastoreConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigMetastoreConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigDataprocMetricConfig(des, initial *ClusterConfigDataprocMetricConfig, opts ...dcl.ApplyOption) *ClusterConfigDataprocMetricConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigDataprocMetricConfig{}

	cDes.Metrics = canonicalizeClusterConfigDataprocMetricConfigMetricsSlice(des.Metrics, initial.Metrics, opts...)

	return cDes
}

func canonicalizeClusterConfigDataprocMetricConfigSlice(des, initial []ClusterConfigDataprocMetricConfig, opts ...dcl.ApplyOption) []ClusterConfigDataprocMetricConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigDataprocMetricConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigDataprocMetricConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigDataprocMetricConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigDataprocMetricConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigDataprocMetricConfig(c *Client, des, nw *ClusterConfigDataprocMetricConfig) *ClusterConfigDataprocMetricConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigDataprocMetricConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Metrics = canonicalizeNewClusterConfigDataprocMetricConfigMetricsSlice(c, des.Metrics, nw.Metrics)

	return nw
}

func canonicalizeNewClusterConfigDataprocMetricConfigSet(c *Client, des, nw []ClusterConfigDataprocMetricConfig) []ClusterConfigDataprocMetricConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigDataprocMetricConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigDataprocMetricConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigDataprocMetricConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigDataprocMetricConfigSlice(c *Client, des, nw []ClusterConfigDataprocMetricConfig) []ClusterConfigDataprocMetricConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigDataprocMetricConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigDataprocMetricConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterConfigDataprocMetricConfigMetrics(des, initial *ClusterConfigDataprocMetricConfigMetrics, opts ...dcl.ApplyOption) *ClusterConfigDataprocMetricConfigMetrics {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterConfigDataprocMetricConfigMetrics{}

	if dcl.IsZeroValue(des.MetricSource) || (dcl.IsEmptyValueIndirect(des.MetricSource) && dcl.IsEmptyValueIndirect(initial.MetricSource)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MetricSource = initial.MetricSource
	} else {
		cDes.MetricSource = des.MetricSource
	}
	if dcl.StringArrayCanonicalize(des.MetricOverrides, initial.MetricOverrides) {
		cDes.MetricOverrides = initial.MetricOverrides
	} else {
		cDes.MetricOverrides = des.MetricOverrides
	}

	return cDes
}

func canonicalizeClusterConfigDataprocMetricConfigMetricsSlice(des, initial []ClusterConfigDataprocMetricConfigMetrics, opts ...dcl.ApplyOption) []ClusterConfigDataprocMetricConfigMetrics {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterConfigDataprocMetricConfigMetrics, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterConfigDataprocMetricConfigMetrics(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterConfigDataprocMetricConfigMetrics, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterConfigDataprocMetricConfigMetrics(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterConfigDataprocMetricConfigMetrics(c *Client, des, nw *ClusterConfigDataprocMetricConfigMetrics) *ClusterConfigDataprocMetricConfigMetrics {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterConfigDataprocMetricConfigMetrics while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.MetricOverrides, nw.MetricOverrides) {
		nw.MetricOverrides = des.MetricOverrides
	}

	return nw
}

func canonicalizeNewClusterConfigDataprocMetricConfigMetricsSet(c *Client, des, nw []ClusterConfigDataprocMetricConfigMetrics) []ClusterConfigDataprocMetricConfigMetrics {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterConfigDataprocMetricConfigMetrics
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterConfigDataprocMetricConfigMetricsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterConfigDataprocMetricConfigMetrics(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterConfigDataprocMetricConfigMetricsSlice(c *Client, des, nw []ClusterConfigDataprocMetricConfigMetrics) []ClusterConfigDataprocMetricConfigMetrics {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterConfigDataprocMetricConfigMetrics
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterConfigDataprocMetricConfigMetrics(c, &d, &n))
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterStatus
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterStatusNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterStatus(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterStatusHistory
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterStatusHistoryNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterStatusHistory(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.HdfsMetrics) || (dcl.IsEmptyValueIndirect(des.HdfsMetrics) && dcl.IsEmptyValueIndirect(initial.HdfsMetrics)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.HdfsMetrics = initial.HdfsMetrics
	} else {
		cDes.HdfsMetrics = des.HdfsMetrics
	}
	if dcl.IsZeroValue(des.YarnMetrics) || (dcl.IsEmptyValueIndirect(des.YarnMetrics) && dcl.IsEmptyValueIndirect(initial.YarnMetrics)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.YarnMetrics = initial.YarnMetrics
	} else {
		cDes.YarnMetrics = des.YarnMetrics
	}

	return cDes
}

func canonicalizeClusterMetricsSlice(des, initial []ClusterMetrics, opts ...dcl.ApplyOption) []ClusterMetrics {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterMetrics
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterMetricsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterMetrics(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

func canonicalizeClusterVirtualClusterConfig(des, initial *ClusterVirtualClusterConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfig{}

	if dcl.IsZeroValue(des.StagingBucket) || (dcl.IsEmptyValueIndirect(des.StagingBucket) && dcl.IsEmptyValueIndirect(initial.StagingBucket)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.StagingBucket = initial.StagingBucket
	} else {
		cDes.StagingBucket = des.StagingBucket
	}
	cDes.KubernetesClusterConfig = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfig(des.KubernetesClusterConfig, initial.KubernetesClusterConfig, opts...)
	cDes.AuxiliaryServicesConfig = canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfig(des.AuxiliaryServicesConfig, initial.AuxiliaryServicesConfig, opts...)

	return cDes
}

func canonicalizeClusterVirtualClusterConfigSlice(des, initial []ClusterVirtualClusterConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfig(c *Client, des, nw *ClusterVirtualClusterConfig) *ClusterVirtualClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KubernetesClusterConfig = canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfig(c, des.KubernetesClusterConfig, nw.KubernetesClusterConfig)
	nw.AuxiliaryServicesConfig = canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfig(c, des.AuxiliaryServicesConfig, nw.AuxiliaryServicesConfig)

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigSet(c *Client, des, nw []ClusterVirtualClusterConfig) []ClusterVirtualClusterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfig) []ClusterVirtualClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfig(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfig{}

	if dcl.StringCanonicalize(des.KubernetesNamespace, initial.KubernetesNamespace) || dcl.IsZeroValue(des.KubernetesNamespace) {
		cDes.KubernetesNamespace = initial.KubernetesNamespace
	} else {
		cDes.KubernetesNamespace = des.KubernetesNamespace
	}
	cDes.GkeClusterConfig = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(des.GkeClusterConfig, initial.GkeClusterConfig, opts...)
	cDes.KubernetesSoftwareConfig = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(des.KubernetesSoftwareConfig, initial.KubernetesSoftwareConfig, opts...)

	return cDes
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfig(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfig) *ClusterVirtualClusterConfigKubernetesClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KubernetesNamespace, nw.KubernetesNamespace) {
		nw.KubernetesNamespace = des.KubernetesNamespace
	}
	nw.GkeClusterConfig = canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, des.GkeClusterConfig, nw.GkeClusterConfig)
	nw.KubernetesSoftwareConfig = canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, des.KubernetesSoftwareConfig, nw.KubernetesSoftwareConfig)

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfig) []ClusterVirtualClusterConfigKubernetesClusterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfig) []ClusterVirtualClusterConfigKubernetesClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{}

	if dcl.IsZeroValue(des.GkeClusterTarget) || (dcl.IsEmptyValueIndirect(des.GkeClusterTarget) && dcl.IsEmptyValueIndirect(initial.GkeClusterTarget)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.GkeClusterTarget = initial.GkeClusterTarget
	} else {
		cDes.GkeClusterTarget = des.GkeClusterTarget
	}
	cDes.NodePoolTarget = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice(des.NodePoolTarget, initial.NodePoolTarget, opts...)

	return cDes
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.NodePoolTarget = canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice(c, des.NodePoolTarget, nw.NodePoolTarget)

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget{}

	if dcl.IsZeroValue(des.NodePool) || (dcl.IsEmptyValueIndirect(des.NodePool) && dcl.IsEmptyValueIndirect(initial.NodePool)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NodePool = initial.NodePool
	} else {
		cDes.NodePool = des.NodePool
	}
	if dcl.IsZeroValue(des.Roles) || (dcl.IsEmptyValueIndirect(des.Roles) && dcl.IsEmptyValueIndirect(initial.Roles)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Roles = initial.Roles
	} else {
		cDes.Roles = des.Roles
	}
	cDes.NodePoolConfig = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(des.NodePoolConfig, initial.NodePoolConfig, opts...)

	return cDes
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.NodePoolConfig = des.NodePoolConfig

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{}

	cDes.Config = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(des.Config, initial.Config, opts...)
	if dcl.StringArrayCanonicalize(des.Locations, initial.Locations) {
		cDes.Locations = initial.Locations
	} else {
		cDes.Locations = des.Locations
	}
	cDes.Autoscaling = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(des.Autoscaling, initial.Autoscaling, opts...)

	return cDes
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Config = canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, des.Config, nw.Config)
	if dcl.StringArrayCanonicalize(des.Locations, nw.Locations) {
		nw.Locations = des.Locations
	}
	nw.Autoscaling = canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, des.Autoscaling, nw.Autoscaling)

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{}

	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	if dcl.IsZeroValue(des.LocalSsdCount) || (dcl.IsEmptyValueIndirect(des.LocalSsdCount) && dcl.IsEmptyValueIndirect(initial.LocalSsdCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.LocalSsdCount = initial.LocalSsdCount
	} else {
		cDes.LocalSsdCount = des.LocalSsdCount
	}
	if dcl.BoolCanonicalize(des.Preemptible, initial.Preemptible) || dcl.IsZeroValue(des.Preemptible) {
		cDes.Preemptible = initial.Preemptible
	} else {
		cDes.Preemptible = des.Preemptible
	}
	cDes.Accelerators = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice(des.Accelerators, initial.Accelerators, opts...)
	if dcl.StringCanonicalize(des.MinCpuPlatform, initial.MinCpuPlatform) || dcl.IsZeroValue(des.MinCpuPlatform) {
		cDes.MinCpuPlatform = initial.MinCpuPlatform
	} else {
		cDes.MinCpuPlatform = des.MinCpuPlatform
	}
	if dcl.StringCanonicalize(des.BootDiskKmsKey, initial.BootDiskKmsKey) || dcl.IsZeroValue(des.BootDiskKmsKey) {
		cDes.BootDiskKmsKey = initial.BootDiskKmsKey
	} else {
		cDes.BootDiskKmsKey = des.BootDiskKmsKey
	}
	cDes.EphemeralStorageConfig = canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(des.EphemeralStorageConfig, initial.EphemeralStorageConfig, opts...)
	if dcl.BoolCanonicalize(des.Spot, initial.Spot) || dcl.IsZeroValue(des.Spot) {
		cDes.Spot = initial.Spot
	} else {
		cDes.Spot = des.Spot
	}

	return cDes
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	if dcl.BoolCanonicalize(des.Preemptible, nw.Preemptible) {
		nw.Preemptible = des.Preemptible
	}
	nw.Accelerators = canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice(c, des.Accelerators, nw.Accelerators)
	if dcl.StringCanonicalize(des.MinCpuPlatform, nw.MinCpuPlatform) {
		nw.MinCpuPlatform = des.MinCpuPlatform
	}
	if dcl.StringCanonicalize(des.BootDiskKmsKey, nw.BootDiskKmsKey) {
		nw.BootDiskKmsKey = des.BootDiskKmsKey
	}
	nw.EphemeralStorageConfig = canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, des.EphemeralStorageConfig, nw.EphemeralStorageConfig)
	if dcl.BoolCanonicalize(des.Spot, nw.Spot) {
		nw.Spot = des.Spot
	}

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators{}

	if dcl.IsZeroValue(des.AcceleratorCount) || (dcl.IsEmptyValueIndirect(des.AcceleratorCount) && dcl.IsEmptyValueIndirect(initial.AcceleratorCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}
	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}
	if dcl.StringCanonicalize(des.GpuPartitionSize, initial.GpuPartitionSize) || dcl.IsZeroValue(des.GpuPartitionSize) {
		cDes.GpuPartitionSize = initial.GpuPartitionSize
	} else {
		cDes.GpuPartitionSize = des.GpuPartitionSize
	}

	return cDes
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}
	if dcl.StringCanonicalize(des.GpuPartitionSize, nw.GpuPartitionSize) {
		nw.GpuPartitionSize = des.GpuPartitionSize
	}

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{}

	if dcl.IsZeroValue(des.LocalSsdCount) || (dcl.IsEmptyValueIndirect(des.LocalSsdCount) && dcl.IsEmptyValueIndirect(initial.LocalSsdCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.LocalSsdCount = initial.LocalSsdCount
	} else {
		cDes.LocalSsdCount = des.LocalSsdCount
	}

	return cDes
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{}

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

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(des, initial *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{}

	if dcl.IsZeroValue(des.ComponentVersion) || (dcl.IsEmptyValueIndirect(des.ComponentVersion) && dcl.IsEmptyValueIndirect(initial.ComponentVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ComponentVersion = initial.ComponentVersion
	} else {
		cDes.ComponentVersion = des.ComponentVersion
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}

	return cDes
}

func canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigSlice(des, initial []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c *Client, des, nw *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfig(des, initial *ClusterVirtualClusterConfigAuxiliaryServicesConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigAuxiliaryServicesConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigAuxiliaryServicesConfig{}

	cDes.MetastoreConfig = canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(des.MetastoreConfig, initial.MetastoreConfig, opts...)
	cDes.SparkHistoryServerConfig = canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(des.SparkHistoryServerConfig, initial.SparkHistoryServerConfig, opts...)

	return cDes
}

func canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigSlice(des, initial []ClusterVirtualClusterConfigAuxiliaryServicesConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigAuxiliaryServicesConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfig(c *Client, des, nw *ClusterVirtualClusterConfigAuxiliaryServicesConfig) *ClusterVirtualClusterConfigAuxiliaryServicesConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigAuxiliaryServicesConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.MetastoreConfig = canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, des.MetastoreConfig, nw.MetastoreConfig)
	nw.SparkHistoryServerConfig = canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, des.SparkHistoryServerConfig, nw.SparkHistoryServerConfig)

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigAuxiliaryServicesConfig) []ClusterVirtualClusterConfigAuxiliaryServicesConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigAuxiliaryServicesConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigAuxiliaryServicesConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigAuxiliaryServicesConfig) []ClusterVirtualClusterConfigAuxiliaryServicesConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigAuxiliaryServicesConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(des, initial *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{}

	if dcl.IsZeroValue(des.DataprocMetastoreService) || (dcl.IsEmptyValueIndirect(des.DataprocMetastoreService) && dcl.IsEmptyValueIndirect(initial.DataprocMetastoreService)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DataprocMetastoreService = initial.DataprocMetastoreService
	} else {
		cDes.DataprocMetastoreService = des.DataprocMetastoreService
	}

	return cDes
}

func canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigSlice(des, initial []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c *Client, des, nw *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(des, initial *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, opts ...dcl.ApplyOption) *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{}

	if dcl.IsZeroValue(des.DataprocCluster) || (dcl.IsEmptyValueIndirect(des.DataprocCluster) && dcl.IsEmptyValueIndirect(initial.DataprocCluster)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DataprocCluster = initial.DataprocCluster
	} else {
		cDes.DataprocCluster = des.DataprocCluster
	}

	return cDes
}

func canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigSlice(des, initial []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, opts ...dcl.ApplyOption) []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c *Client, des, nw *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigSet(c *Client, des, nw []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigSlice(c *Client, des, nw []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, &d, &n))
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
	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ProjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.DiffInfo{ObjectFunction: compareClusterConfigNewStyle, EmptyObject: EmptyClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{IgnoredPrefixes: []string{"goog-dataproc-"}, OperationSelector: dcl.TriggersOperation("updateClusterUpdateClusterOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Status, actual.Status, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareClusterStatusNewStyle, EmptyObject: EmptyClusterStatus, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Status")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StatusHistory, actual.StatusHistory, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareClusterStatusHistoryNewStyle, EmptyObject: EmptyClusterStatusHistory, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StatusHistory")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClusterUuid, actual.ClusterUuid, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterUuid")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metrics, actual.Metrics, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareClusterMetricsNewStyle, EmptyObject: EmptyClusterMetrics, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Metrics")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.VirtualClusterConfig, actual.VirtualClusterConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("VirtualClusterConfig")); len(ds) != 0 || err != nil {
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
func compareClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfig or *ClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.StagingBucket, actual.StagingBucket, dcl.DiffInfo{ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConfigBucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TempBucket, actual.TempBucket, dcl.DiffInfo{ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TempBucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GceClusterConfig, actual.GceClusterConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigGceClusterConfigNewStyle, EmptyObject: EmptyClusterConfigGceClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GceClusterConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MasterConfig, actual.MasterConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigMasterConfigNewStyle, EmptyObject: EmptyClusterConfigMasterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MasterConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WorkerConfig, actual.WorkerConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigWorkerConfigNewStyle, EmptyObject: EmptyClusterConfigWorkerConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("WorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecondaryWorkerConfig, actual.SecondaryWorkerConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigSecondaryWorkerConfigNewStyle, EmptyObject: EmptyClusterConfigSecondaryWorkerConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SecondaryWorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SoftwareConfig, actual.SoftwareConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigSoftwareConfigNewStyle, EmptyObject: EmptyClusterConfigSoftwareConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SoftwareConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InitializationActions, actual.InitializationActions, dcl.DiffInfo{ObjectFunction: compareClusterConfigInitializationActionsNewStyle, EmptyObject: EmptyClusterConfigInitializationActions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InitializationActions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EncryptionConfig, actual.EncryptionConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigEncryptionConfigNewStyle, EmptyObject: EmptyClusterConfigEncryptionConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EncryptionConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoscalingConfig, actual.AutoscalingConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigAutoscalingConfigNewStyle, EmptyObject: EmptyClusterConfigAutoscalingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoscalingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecurityConfig, actual.SecurityConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigSecurityConfigNewStyle, EmptyObject: EmptyClusterConfigSecurityConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SecurityConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LifecycleConfig, actual.LifecycleConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigLifecycleConfigNewStyle, EmptyObject: EmptyClusterConfigLifecycleConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LifecycleConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EndpointConfig, actual.EndpointConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigEndpointConfigNewStyle, EmptyObject: EmptyClusterConfigEndpointConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EndpointConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetastoreConfig, actual.MetastoreConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigMetastoreConfigNewStyle, EmptyObject: EmptyClusterConfigMetastoreConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetastoreConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DataprocMetricConfig, actual.DataprocMetricConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigDataprocMetricConfigNewStyle, EmptyObject: EmptyClusterConfigDataprocMetricConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DataprocMetricConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigGceClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigGceClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigGceClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfig or *ClusterConfigGceClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigGceClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigGceClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Zone, actual.Zone, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ZoneUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Network, actual.Network, dcl.DiffInfo{ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Subnetwork, actual.Subnetwork, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubnetworkUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InternalIPOnly, actual.InternalIPOnly, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InternalIpOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrivateIPv6GoogleAccess, actual.PrivateIPv6GoogleAccess, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrivateIpv6GoogleAccess")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccount, actual.ServiceAccount, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAccount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccountScopes, actual.ServiceAccountScopes, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAccountScopes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tags, actual.Tags, dcl.DiffInfo{Type: "Set", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metadata, actual.Metadata, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Metadata")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ReservationAffinity, actual.ReservationAffinity, dcl.DiffInfo{ObjectFunction: compareClusterConfigGceClusterConfigReservationAffinityNewStyle, EmptyObject: EmptyClusterConfigGceClusterConfigReservationAffinity, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ReservationAffinity")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NodeGroupAffinity, actual.NodeGroupAffinity, dcl.DiffInfo{ObjectFunction: compareClusterConfigGceClusterConfigNodeGroupAffinityNewStyle, EmptyObject: EmptyClusterConfigGceClusterConfigNodeGroupAffinity, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodeGroupAffinity")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ShieldedInstanceConfig, actual.ShieldedInstanceConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigGceClusterConfigShieldedInstanceConfigNewStyle, EmptyObject: EmptyClusterConfigGceClusterConfigShieldedInstanceConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ShieldedInstanceConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ConfidentialInstanceConfig, actual.ConfidentialInstanceConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigGceClusterConfigConfidentialInstanceConfigNewStyle, EmptyObject: EmptyClusterConfigGceClusterConfigConfidentialInstanceConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConfidentialInstanceConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigGceClusterConfigReservationAffinityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigGceClusterConfigReservationAffinity)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigGceClusterConfigReservationAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfigReservationAffinity or *ClusterConfigGceClusterConfigReservationAffinity", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigGceClusterConfigReservationAffinity)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigGceClusterConfigReservationAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfigReservationAffinity", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ConsumeReservationType, actual.ConsumeReservationType, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConsumeReservationType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Key, actual.Key, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Key")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigGceClusterConfigNodeGroupAffinityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigGceClusterConfigNodeGroupAffinity)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigGceClusterConfigNodeGroupAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfigNodeGroupAffinity or *ClusterConfigGceClusterConfigNodeGroupAffinity", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigGceClusterConfigNodeGroupAffinity)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigGceClusterConfigNodeGroupAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfigNodeGroupAffinity", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NodeGroup, actual.NodeGroup, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodeGroupUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigGceClusterConfigShieldedInstanceConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigGceClusterConfigShieldedInstanceConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigGceClusterConfigShieldedInstanceConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfigShieldedInstanceConfig or *ClusterConfigGceClusterConfigShieldedInstanceConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigGceClusterConfigShieldedInstanceConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigGceClusterConfigShieldedInstanceConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfigShieldedInstanceConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.EnableSecureBoot, actual.EnableSecureBoot, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableSecureBoot")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableVtpm, actual.EnableVtpm, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableVtpm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableIntegrityMonitoring, actual.EnableIntegrityMonitoring, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableIntegrityMonitoring")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigGceClusterConfigConfidentialInstanceConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigGceClusterConfigConfidentialInstanceConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigGceClusterConfigConfidentialInstanceConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfigConfidentialInstanceConfig or *ClusterConfigGceClusterConfigConfidentialInstanceConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigGceClusterConfigConfidentialInstanceConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigGceClusterConfigConfidentialInstanceConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigGceClusterConfigConfidentialInstanceConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.EnableConfidentialCompute, actual.EnableConfidentialCompute, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableConfidentialCompute")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigMasterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigMasterConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigMasterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfig or *ClusterConfigMasterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigMasterConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigMasterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumInstances, actual.NumInstances, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceNames, actual.InstanceNames, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Image, actual.Image, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MachineTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskConfig, actual.DiskConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigMasterConfigDiskConfigNewStyle, EmptyObject: EmptyClusterConfigMasterConfigDiskConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IsPreemptible, actual.IsPreemptible, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsPreemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptibility, actual.Preemptibility, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptibility")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ManagedGroupConfig, actual.ManagedGroupConfig, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, ObjectFunction: compareClusterConfigMasterConfigManagedGroupConfigNewStyle, EmptyObject: EmptyClusterConfigMasterConfigManagedGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedGroupConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Accelerators, actual.Accelerators, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigMasterConfigAcceleratorsNewStyle, EmptyObject: EmptyClusterConfigMasterConfigAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Accelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceReferences, actual.InstanceReferences, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareClusterConfigMasterConfigInstanceReferencesNewStyle, EmptyObject: EmptyClusterConfigMasterConfigInstanceReferences, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceReferences")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigMasterConfigDiskConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigMasterConfigDiskConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigMasterConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfigDiskConfig or *ClusterConfigMasterConfigDiskConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigMasterConfigDiskConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigMasterConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfigDiskConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BootDiskType, actual.BootDiskType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BootDiskSizeGb, actual.BootDiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NumLocalSsds, actual.NumLocalSsds, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumLocalSsds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalSsdInterface, actual.LocalSsdInterface, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LocalSsdInterface")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigMasterConfigManagedGroupConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigMasterConfigManagedGroupConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigMasterConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfigManagedGroupConfig or *ClusterConfigMasterConfigManagedGroupConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigMasterConfigManagedGroupConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigMasterConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfigManagedGroupConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceTemplateName, actual.InstanceTemplateName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceTemplateName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceGroupManagerName, actual.InstanceGroupManagerName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceGroupManagerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigMasterConfigAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigMasterConfigAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigMasterConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfigAccelerators or *ClusterConfigMasterConfigAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigMasterConfigAccelerators)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigMasterConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfigAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigMasterConfigInstanceReferencesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigMasterConfigInstanceReferences)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigMasterConfigInstanceReferences)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfigInstanceReferences or *ClusterConfigMasterConfigInstanceReferences", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigMasterConfigInstanceReferences)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigMasterConfigInstanceReferences)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMasterConfigInstanceReferences", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceName, actual.InstanceName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceId, actual.InstanceId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicKey, actual.PublicKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicEciesKey, actual.PublicEciesKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicEciesKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigWorkerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigWorkerConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfig or *ClusterConfigWorkerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigWorkerConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumInstances, actual.NumInstances, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceNames, actual.InstanceNames, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Image, actual.Image, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MachineTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskConfig, actual.DiskConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigWorkerConfigDiskConfigNewStyle, EmptyObject: EmptyClusterConfigWorkerConfigDiskConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IsPreemptible, actual.IsPreemptible, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsPreemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptibility, actual.Preemptibility, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptibility")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ManagedGroupConfig, actual.ManagedGroupConfig, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, ObjectFunction: compareClusterConfigWorkerConfigManagedGroupConfigNewStyle, EmptyObject: EmptyClusterConfigWorkerConfigManagedGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedGroupConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Accelerators, actual.Accelerators, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigWorkerConfigAcceleratorsNewStyle, EmptyObject: EmptyClusterConfigWorkerConfigAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Accelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceReferences, actual.InstanceReferences, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareClusterConfigWorkerConfigInstanceReferencesNewStyle, EmptyObject: EmptyClusterConfigWorkerConfigInstanceReferences, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceReferences")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigWorkerConfigDiskConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigWorkerConfigDiskConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigWorkerConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfigDiskConfig or *ClusterConfigWorkerConfigDiskConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigWorkerConfigDiskConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigWorkerConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfigDiskConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BootDiskType, actual.BootDiskType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BootDiskSizeGb, actual.BootDiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NumLocalSsds, actual.NumLocalSsds, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumLocalSsds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalSsdInterface, actual.LocalSsdInterface, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LocalSsdInterface")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigWorkerConfigManagedGroupConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigWorkerConfigManagedGroupConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigWorkerConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfigManagedGroupConfig or *ClusterConfigWorkerConfigManagedGroupConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigWorkerConfigManagedGroupConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigWorkerConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfigManagedGroupConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceTemplateName, actual.InstanceTemplateName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceTemplateName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceGroupManagerName, actual.InstanceGroupManagerName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceGroupManagerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigWorkerConfigAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigWorkerConfigAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigWorkerConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfigAccelerators or *ClusterConfigWorkerConfigAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigWorkerConfigAccelerators)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigWorkerConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfigAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigWorkerConfigInstanceReferencesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigWorkerConfigInstanceReferences)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigWorkerConfigInstanceReferences)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfigInstanceReferences or *ClusterConfigWorkerConfigInstanceReferences", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigWorkerConfigInstanceReferences)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigWorkerConfigInstanceReferences)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigWorkerConfigInstanceReferences", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceName, actual.InstanceName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceId, actual.InstanceId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicKey, actual.PublicKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicEciesKey, actual.PublicEciesKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicEciesKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSecondaryWorkerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSecondaryWorkerConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSecondaryWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfig or *ClusterConfigSecondaryWorkerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSecondaryWorkerConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSecondaryWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumInstances, actual.NumInstances, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceNames, actual.InstanceNames, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Image, actual.Image, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MachineTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskConfig, actual.DiskConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigSecondaryWorkerConfigDiskConfigNewStyle, EmptyObject: EmptyClusterConfigSecondaryWorkerConfigDiskConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IsPreemptible, actual.IsPreemptible, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsPreemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptibility, actual.Preemptibility, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptibility")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ManagedGroupConfig, actual.ManagedGroupConfig, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, ObjectFunction: compareClusterConfigSecondaryWorkerConfigManagedGroupConfigNewStyle, EmptyObject: EmptyClusterConfigSecondaryWorkerConfigManagedGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedGroupConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Accelerators, actual.Accelerators, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterConfigSecondaryWorkerConfigAcceleratorsNewStyle, EmptyObject: EmptyClusterConfigSecondaryWorkerConfigAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Accelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceReferences, actual.InstanceReferences, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareClusterConfigSecondaryWorkerConfigInstanceReferencesNewStyle, EmptyObject: EmptyClusterConfigSecondaryWorkerConfigInstanceReferences, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceReferences")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSecondaryWorkerConfigDiskConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSecondaryWorkerConfigDiskConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSecondaryWorkerConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfigDiskConfig or *ClusterConfigSecondaryWorkerConfigDiskConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSecondaryWorkerConfigDiskConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSecondaryWorkerConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfigDiskConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BootDiskType, actual.BootDiskType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BootDiskSizeGb, actual.BootDiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NumLocalSsds, actual.NumLocalSsds, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumLocalSsds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalSsdInterface, actual.LocalSsdInterface, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LocalSsdInterface")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSecondaryWorkerConfigManagedGroupConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSecondaryWorkerConfigManagedGroupConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSecondaryWorkerConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfigManagedGroupConfig or *ClusterConfigSecondaryWorkerConfigManagedGroupConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSecondaryWorkerConfigManagedGroupConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSecondaryWorkerConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfigManagedGroupConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceTemplateName, actual.InstanceTemplateName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceTemplateName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceGroupManagerName, actual.InstanceGroupManagerName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceGroupManagerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSecondaryWorkerConfigAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSecondaryWorkerConfigAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSecondaryWorkerConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfigAccelerators or *ClusterConfigSecondaryWorkerConfigAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSecondaryWorkerConfigAccelerators)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSecondaryWorkerConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfigAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSecondaryWorkerConfigInstanceReferencesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSecondaryWorkerConfigInstanceReferences)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSecondaryWorkerConfigInstanceReferences)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfigInstanceReferences or *ClusterConfigSecondaryWorkerConfigInstanceReferences", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSecondaryWorkerConfigInstanceReferences)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSecondaryWorkerConfigInstanceReferences)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecondaryWorkerConfigInstanceReferences", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceName, actual.InstanceName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceId, actual.InstanceId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicKey, actual.PublicKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicEciesKey, actual.PublicEciesKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicEciesKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSoftwareConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSoftwareConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSoftwareConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSoftwareConfig or *ClusterConfigSoftwareConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSoftwareConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSoftwareConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSoftwareConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ImageVersion, actual.ImageVersion, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{CustomDiff: canonicalizeSoftwareConfigProperties, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OptionalComponents, actual.OptionalComponents, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OptionalComponents")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigInitializationActionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigInitializationActions)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigInitializationActions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigInitializationActions or *ClusterConfigInitializationActions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigInitializationActions)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigInitializationActions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigInitializationActions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ExecutableFile, actual.ExecutableFile, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExecutableFile")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExecutionTimeout, actual.ExecutionTimeout, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExecutionTimeout")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigEncryptionConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigEncryptionConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigEncryptionConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigEncryptionConfig or *ClusterConfigEncryptionConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigEncryptionConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigEncryptionConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigEncryptionConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GcePdKmsKeyName, actual.GcePdKmsKeyName, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GcePdKmsKeyName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigAutoscalingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigAutoscalingConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigAutoscalingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigAutoscalingConfig or *ClusterConfigAutoscalingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigAutoscalingConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigAutoscalingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigAutoscalingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Policy, actual.Policy, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PolicyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSecurityConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSecurityConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSecurityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecurityConfig or *ClusterConfigSecurityConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSecurityConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSecurityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecurityConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KerberosConfig, actual.KerberosConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigSecurityConfigKerberosConfigNewStyle, EmptyObject: EmptyClusterConfigSecurityConfigKerberosConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KerberosConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IdentityConfig, actual.IdentityConfig, dcl.DiffInfo{ObjectFunction: compareClusterConfigSecurityConfigIdentityConfigNewStyle, EmptyObject: EmptyClusterConfigSecurityConfigIdentityConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IdentityConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSecurityConfigKerberosConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSecurityConfigKerberosConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSecurityConfigKerberosConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecurityConfigKerberosConfig or *ClusterConfigSecurityConfigKerberosConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSecurityConfigKerberosConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSecurityConfigKerberosConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecurityConfigKerberosConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.EnableKerberos, actual.EnableKerberos, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableKerberos")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RootPrincipalPassword, actual.RootPrincipalPassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RootPrincipalPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KmsKey, actual.KmsKey, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KmsKeyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Keystore, actual.Keystore, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeystoreUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Truststore, actual.Truststore, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TruststoreUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeystorePassword, actual.KeystorePassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeystorePasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyPassword, actual.KeyPassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TruststorePassword, actual.TruststorePassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TruststorePasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustRealm, actual.CrossRealmTrustRealm, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustRealm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustKdc, actual.CrossRealmTrustKdc, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustKdc")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustAdminServer, actual.CrossRealmTrustAdminServer, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustAdminServer")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustSharedPassword, actual.CrossRealmTrustSharedPassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustSharedPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KdcDbKey, actual.KdcDbKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KdcDbKeyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TgtLifetimeHours, actual.TgtLifetimeHours, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TgtLifetimeHours")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Realm, actual.Realm, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Realm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigSecurityConfigIdentityConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigSecurityConfigIdentityConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigSecurityConfigIdentityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecurityConfigIdentityConfig or *ClusterConfigSecurityConfigIdentityConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigSecurityConfigIdentityConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigSecurityConfigIdentityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigSecurityConfigIdentityConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.UserServiceAccountMapping, actual.UserServiceAccountMapping, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UserServiceAccountMapping")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigLifecycleConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigLifecycleConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigLifecycleConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigLifecycleConfig or *ClusterConfigLifecycleConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigLifecycleConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigLifecycleConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigLifecycleConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IdleDeleteTtl, actual.IdleDeleteTtl, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IdleDeleteTtl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoDeleteTime, actual.AutoDeleteTime, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoDeleteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoDeleteTtl, actual.AutoDeleteTtl, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoDeleteTtl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IdleStartTime, actual.IdleStartTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IdleStartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigEndpointConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigEndpointConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigEndpointConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigEndpointConfig or *ClusterConfigEndpointConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigEndpointConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigEndpointConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigEndpointConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.HttpPorts, actual.HttpPorts, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HttpPorts")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableHttpPortAccess, actual.EnableHttpPortAccess, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableHttpPortAccess")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigMetastoreConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigMetastoreConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigMetastoreConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMetastoreConfig or *ClusterConfigMetastoreConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigMetastoreConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigMetastoreConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigMetastoreConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DataprocMetastoreService, actual.DataprocMetastoreService, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DataprocMetastoreService")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigDataprocMetricConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigDataprocMetricConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigDataprocMetricConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigDataprocMetricConfig or *ClusterConfigDataprocMetricConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigDataprocMetricConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigDataprocMetricConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigDataprocMetricConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Metrics, actual.Metrics, dcl.DiffInfo{ObjectFunction: compareClusterConfigDataprocMetricConfigMetricsNewStyle, EmptyObject: EmptyClusterConfigDataprocMetricConfigMetrics, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Metrics")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterConfigDataprocMetricConfigMetricsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterConfigDataprocMetricConfigMetrics)
	if !ok {
		desiredNotPointer, ok := d.(ClusterConfigDataprocMetricConfigMetrics)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigDataprocMetricConfigMetrics or *ClusterConfigDataprocMetricConfigMetrics", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterConfigDataprocMetricConfigMetrics)
	if !ok {
		actualNotPointer, ok := a.(ClusterConfigDataprocMetricConfigMetrics)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterConfigDataprocMetricConfigMetrics", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MetricSource, actual.MetricSource, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetricSource")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MetricOverrides, actual.MetricOverrides, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetricOverrides")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Detail, actual.Detail, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Detail")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StateStartTime, actual.StateStartTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StateStartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Substate, actual.Substate, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Substate")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Detail, actual.Detail, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Detail")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StateStartTime, actual.StateStartTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StateStartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Substate, actual.Substate, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Substate")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.HdfsMetrics, actual.HdfsMetrics, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HdfsMetrics")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.YarnMetrics, actual.YarnMetrics, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("YarnMetrics")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfig or *ClusterVirtualClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.StagingBucket, actual.StagingBucket, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StagingBucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KubernetesClusterConfig, actual.KubernetesClusterConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KubernetesClusterConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AuxiliaryServicesConfig, actual.AuxiliaryServicesConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigAuxiliaryServicesConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigAuxiliaryServicesConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AuxiliaryServicesConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfig or *ClusterVirtualClusterConfigKubernetesClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KubernetesNamespace, actual.KubernetesNamespace, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KubernetesNamespace")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GkeClusterConfig, actual.GkeClusterConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GkeClusterConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KubernetesSoftwareConfig, actual.KubernetesSoftwareConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KubernetesSoftwareConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig or *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GkeClusterTarget, actual.GkeClusterTarget, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GkeClusterTarget")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NodePoolTarget, actual.NodePoolTarget, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodePoolTarget")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget or *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NodePool, actual.NodePool, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodePool")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Roles, actual.Roles, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Roles")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NodePoolConfig, actual.NodePoolConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodePoolConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig or *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Locations, actual.Locations, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Locations")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Autoscaling, actual.Autoscaling, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Autoscaling")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig or *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MachineType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalSsdCount, actual.LocalSsdCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LocalSsdCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptible, actual.Preemptible, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Accelerators, actual.Accelerators, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Accelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BootDiskKmsKey, actual.BootDiskKmsKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskKmsKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EphemeralStorageConfig, actual.EphemeralStorageConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EphemeralStorageConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Spot, actual.Spot, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Spot")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators or *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GpuPartitionSize, actual.GpuPartitionSize, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GpuPartitionSize")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig or *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.LocalSsdCount, actual.LocalSsdCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LocalSsdCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling or *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MinNodeCount, actual.MinNodeCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinNodeCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxNodeCount, actual.MaxNodeCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxNodeCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig or *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ComponentVersion, actual.ComponentVersion, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ComponentVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigAuxiliaryServicesConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigAuxiliaryServicesConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigAuxiliaryServicesConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigAuxiliaryServicesConfig or *ClusterVirtualClusterConfigAuxiliaryServicesConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigAuxiliaryServicesConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigAuxiliaryServicesConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigAuxiliaryServicesConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MetastoreConfig, actual.MetastoreConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MetastoreConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SparkHistoryServerConfig, actual.SparkHistoryServerConfig, dcl.DiffInfo{ObjectFunction: compareClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigNewStyle, EmptyObject: EmptyClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SparkHistoryServerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig or *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DataprocMetastoreService, actual.DataprocMetastoreService, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DataprocMetastoreService")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig or *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DataprocCluster, actual.DataprocCluster, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DataprocCluster")); len(ds) != 0 || err != nil {
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
func unmarshalCluster(b []byte, c *Client, res *Cluster) (*Cluster, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapCluster(m, c, res)
}

func unmarshalMapCluster(m map[string]interface{}, c *Client, res *Cluster) (*Cluster, error) {

	flattened := flattenCluster(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandCluster expands Cluster into a JSON request object.
func expandCluster(c *Client, f *Cluster) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into projectId: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["projectId"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["clusterName"] = v
	}
	if v, err := expandClusterConfig(c, f.Config, res); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["config"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Location into location: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}
	if v, err := expandClusterVirtualClusterConfig(c, f.VirtualClusterConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding VirtualClusterConfig into virtualClusterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["virtualClusterConfig"] = v
	}

	return m, nil
}

// flattenCluster flattens Cluster from a JSON request object into the
// Cluster type.
func flattenCluster(c *Client, i interface{}, res *Cluster) *Cluster {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Cluster{}
	resultRes.Project = dcl.FlattenString(m["projectId"])
	resultRes.Name = dcl.FlattenString(m["clusterName"])
	resultRes.Config = flattenClusterConfig(c, m["config"], res)
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Status = flattenClusterStatus(c, m["status"], res)
	resultRes.StatusHistory = flattenClusterStatusHistorySlice(c, m["statusHistory"], res)
	resultRes.ClusterUuid = dcl.FlattenString(m["clusterUuid"])
	resultRes.Metrics = flattenClusterMetrics(c, m["metrics"], res)
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.VirtualClusterConfig = flattenClusterVirtualClusterConfig(c, m["virtualClusterConfig"], res)

	return resultRes
}

// expandClusterConfigMap expands the contents of ClusterConfig into a JSON
// request object.
func expandClusterConfigMap(c *Client, f map[string]ClusterConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSlice expands the contents of ClusterConfig into a JSON
// request object.
func expandClusterConfigSlice(c *Client, f []ClusterConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigMap flattens the contents of ClusterConfig from a JSON
// response object.
func flattenClusterConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfig{}
	}

	items := make(map[string]ClusterConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSlice flattens the contents of ClusterConfig from a JSON
// response object.
func flattenClusterConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfig{}
	}

	items := make([]ClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfig expands an instance of ClusterConfig into a JSON
// request object.
func expandClusterConfig(c *Client, f *ClusterConfig, res *Cluster) (map[string]interface{}, error) {
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
	if v, err := expandClusterConfigGceClusterConfig(c, f.GceClusterConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GceClusterConfig into gceClusterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gceClusterConfig"] = v
	}
	if v, err := expandClusterConfigMasterConfig(c, f.MasterConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding MasterConfig into masterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["masterConfig"] = v
	}
	if v, err := expandClusterConfigWorkerConfig(c, f.WorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding WorkerConfig into workerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["workerConfig"] = v
	}
	if v, err := expandClusterConfigSecondaryWorkerConfig(c, f.SecondaryWorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SecondaryWorkerConfig into secondaryWorkerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["secondaryWorkerConfig"] = v
	}
	if v, err := expandClusterConfigSoftwareConfig(c, f.SoftwareConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SoftwareConfig into softwareConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["softwareConfig"] = v
	}
	if v, err := expandClusterConfigInitializationActionsSlice(c, f.InitializationActions, res); err != nil {
		return nil, fmt.Errorf("error expanding InitializationActions into initializationActions: %w", err)
	} else if v != nil {
		m["initializationActions"] = v
	}
	if v, err := expandClusterConfigEncryptionConfig(c, f.EncryptionConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding EncryptionConfig into encryptionConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["encryptionConfig"] = v
	}
	if v, err := expandClusterConfigAutoscalingConfig(c, f.AutoscalingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding AutoscalingConfig into autoscalingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["autoscalingConfig"] = v
	}
	if v, err := expandClusterConfigSecurityConfig(c, f.SecurityConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SecurityConfig into securityConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["securityConfig"] = v
	}
	if v, err := expandClusterConfigLifecycleConfig(c, f.LifecycleConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LifecycleConfig into lifecycleConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["lifecycleConfig"] = v
	}
	if v, err := expandClusterConfigEndpointConfig(c, f.EndpointConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding EndpointConfig into endpointConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["endpointConfig"] = v
	}
	if v, err := expandClusterConfigMetastoreConfig(c, f.MetastoreConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding MetastoreConfig into metastoreConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metastoreConfig"] = v
	}
	if v, err := expandClusterConfigDataprocMetricConfig(c, f.DataprocMetricConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding DataprocMetricConfig into dataprocMetricConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["dataprocMetricConfig"] = v
	}

	return m, nil
}

// flattenClusterConfig flattens an instance of ClusterConfig from a JSON
// response object.
func flattenClusterConfig(c *Client, i interface{}, res *Cluster) *ClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfig
	}
	r.StagingBucket = dcl.FlattenString(m["configBucket"])
	r.TempBucket = dcl.FlattenString(m["tempBucket"])
	r.GceClusterConfig = flattenClusterConfigGceClusterConfig(c, m["gceClusterConfig"], res)
	r.MasterConfig = flattenClusterConfigMasterConfig(c, m["masterConfig"], res)
	r.WorkerConfig = flattenClusterConfigWorkerConfig(c, m["workerConfig"], res)
	r.SecondaryWorkerConfig = flattenClusterConfigSecondaryWorkerConfig(c, m["secondaryWorkerConfig"], res)
	r.SoftwareConfig = flattenClusterConfigSoftwareConfig(c, m["softwareConfig"], res)
	r.InitializationActions = flattenClusterConfigInitializationActionsSlice(c, m["initializationActions"], res)
	r.EncryptionConfig = flattenClusterConfigEncryptionConfig(c, m["encryptionConfig"], res)
	r.AutoscalingConfig = flattenClusterConfigAutoscalingConfig(c, m["autoscalingConfig"], res)
	r.SecurityConfig = flattenClusterConfigSecurityConfig(c, m["securityConfig"], res)
	r.LifecycleConfig = flattenClusterConfigLifecycleConfig(c, m["lifecycleConfig"], res)
	r.EndpointConfig = flattenClusterConfigEndpointConfig(c, m["endpointConfig"], res)
	r.MetastoreConfig = flattenClusterConfigMetastoreConfig(c, m["metastoreConfig"], res)
	r.DataprocMetricConfig = flattenClusterConfigDataprocMetricConfig(c, m["dataprocMetricConfig"], res)

	return r
}

// expandClusterConfigGceClusterConfigMap expands the contents of ClusterConfigGceClusterConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfigMap(c *Client, f map[string]ClusterConfigGceClusterConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigGceClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigGceClusterConfigSlice expands the contents of ClusterConfigGceClusterConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfigSlice(c *Client, f []ClusterConfigGceClusterConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigGceClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigGceClusterConfigMap flattens the contents of ClusterConfigGceClusterConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigGceClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigGceClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigGceClusterConfig{}
	}

	items := make(map[string]ClusterConfigGceClusterConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigGceClusterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigGceClusterConfigSlice flattens the contents of ClusterConfigGceClusterConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigGceClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigGceClusterConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigGceClusterConfig{}
	}

	items := make([]ClusterConfigGceClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigGceClusterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigGceClusterConfig expands an instance of ClusterConfigGceClusterConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfig(c *Client, f *ClusterConfigGceClusterConfig, res *Cluster) (map[string]interface{}, error) {
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
	if v, err := expandClusterConfigGceClusterConfigReservationAffinity(c, f.ReservationAffinity, res); err != nil {
		return nil, fmt.Errorf("error expanding ReservationAffinity into reservationAffinity: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["reservationAffinity"] = v
	}
	if v, err := expandClusterConfigGceClusterConfigNodeGroupAffinity(c, f.NodeGroupAffinity, res); err != nil {
		return nil, fmt.Errorf("error expanding NodeGroupAffinity into nodeGroupAffinity: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["nodeGroupAffinity"] = v
	}
	if v, err := expandClusterConfigGceClusterConfigShieldedInstanceConfig(c, f.ShieldedInstanceConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ShieldedInstanceConfig into shieldedInstanceConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["shieldedInstanceConfig"] = v
	}
	if v, err := expandClusterConfigGceClusterConfigConfidentialInstanceConfig(c, f.ConfidentialInstanceConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ConfidentialInstanceConfig into confidentialInstanceConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["confidentialInstanceConfig"] = v
	}

	return m, nil
}

// flattenClusterConfigGceClusterConfig flattens an instance of ClusterConfigGceClusterConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigGceClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigGceClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigGceClusterConfig
	}
	r.Zone = dcl.FlattenString(m["zoneUri"])
	r.Network = dcl.FlattenString(m["networkUri"])
	r.Subnetwork = dcl.FlattenString(m["subnetworkUri"])
	r.InternalIPOnly = dcl.FlattenBool(m["internalIpOnly"])
	r.PrivateIPv6GoogleAccess = flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(m["privateIpv6GoogleAccess"])
	r.ServiceAccount = dcl.FlattenString(m["serviceAccount"])
	r.ServiceAccountScopes = dcl.FlattenStringSlice(m["serviceAccountScopes"])
	r.Tags = dcl.FlattenStringSlice(m["tags"])
	r.Metadata = dcl.FlattenKeyValuePairs(m["metadata"])
	r.ReservationAffinity = flattenClusterConfigGceClusterConfigReservationAffinity(c, m["reservationAffinity"], res)
	r.NodeGroupAffinity = flattenClusterConfigGceClusterConfigNodeGroupAffinity(c, m["nodeGroupAffinity"], res)
	r.ShieldedInstanceConfig = flattenClusterConfigGceClusterConfigShieldedInstanceConfig(c, m["shieldedInstanceConfig"], res)
	r.ConfidentialInstanceConfig = flattenClusterConfigGceClusterConfigConfidentialInstanceConfig(c, m["confidentialInstanceConfig"], res)

	return r
}

// expandClusterConfigGceClusterConfigReservationAffinityMap expands the contents of ClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandClusterConfigGceClusterConfigReservationAffinityMap(c *Client, f map[string]ClusterConfigGceClusterConfigReservationAffinity, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigGceClusterConfigReservationAffinity(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigGceClusterConfigReservationAffinitySlice expands the contents of ClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, f []ClusterConfigGceClusterConfigReservationAffinity, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigGceClusterConfigReservationAffinity(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigGceClusterConfigReservationAffinityMap flattens the contents of ClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenClusterConfigGceClusterConfigReservationAffinityMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigGceClusterConfigReservationAffinity {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigGceClusterConfigReservationAffinity{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigGceClusterConfigReservationAffinity{}
	}

	items := make(map[string]ClusterConfigGceClusterConfigReservationAffinity)
	for k, item := range a {
		items[k] = *flattenClusterConfigGceClusterConfigReservationAffinity(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigGceClusterConfigReservationAffinitySlice flattens the contents of ClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, i interface{}, res *Cluster) []ClusterConfigGceClusterConfigReservationAffinity {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigGceClusterConfigReservationAffinity{}
	}

	if len(a) == 0 {
		return []ClusterConfigGceClusterConfigReservationAffinity{}
	}

	items := make([]ClusterConfigGceClusterConfigReservationAffinity, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigGceClusterConfigReservationAffinity(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigGceClusterConfigReservationAffinity expands an instance of ClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandClusterConfigGceClusterConfigReservationAffinity(c *Client, f *ClusterConfigGceClusterConfigReservationAffinity, res *Cluster) (map[string]interface{}, error) {
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

// flattenClusterConfigGceClusterConfigReservationAffinity flattens an instance of ClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenClusterConfigGceClusterConfigReservationAffinity(c *Client, i interface{}, res *Cluster) *ClusterConfigGceClusterConfigReservationAffinity {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigGceClusterConfigReservationAffinity{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigGceClusterConfigReservationAffinity
	}
	r.ConsumeReservationType = flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(m["consumeReservationType"])
	r.Key = dcl.FlattenString(m["key"])
	r.Values = dcl.FlattenStringSlice(m["values"])

	return r
}

// expandClusterConfigGceClusterConfigNodeGroupAffinityMap expands the contents of ClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandClusterConfigGceClusterConfigNodeGroupAffinityMap(c *Client, f map[string]ClusterConfigGceClusterConfigNodeGroupAffinity, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigGceClusterConfigNodeGroupAffinity(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigGceClusterConfigNodeGroupAffinitySlice expands the contents of ClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, f []ClusterConfigGceClusterConfigNodeGroupAffinity, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigGceClusterConfigNodeGroupAffinity(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigGceClusterConfigNodeGroupAffinityMap flattens the contents of ClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenClusterConfigGceClusterConfigNodeGroupAffinityMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigGceClusterConfigNodeGroupAffinity {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	items := make(map[string]ClusterConfigGceClusterConfigNodeGroupAffinity)
	for k, item := range a {
		items[k] = *flattenClusterConfigGceClusterConfigNodeGroupAffinity(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigGceClusterConfigNodeGroupAffinitySlice flattens the contents of ClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, i interface{}, res *Cluster) []ClusterConfigGceClusterConfigNodeGroupAffinity {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	if len(a) == 0 {
		return []ClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	items := make([]ClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigGceClusterConfigNodeGroupAffinity(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigGceClusterConfigNodeGroupAffinity expands an instance of ClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, f *ClusterConfigGceClusterConfigNodeGroupAffinity, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NodeGroup; !dcl.IsEmptyValueIndirect(v) {
		m["nodeGroupUri"] = v
	}

	return m, nil
}

// flattenClusterConfigGceClusterConfigNodeGroupAffinity flattens an instance of ClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, i interface{}, res *Cluster) *ClusterConfigGceClusterConfigNodeGroupAffinity {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigGceClusterConfigNodeGroupAffinity{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigGceClusterConfigNodeGroupAffinity
	}
	r.NodeGroup = dcl.FlattenString(m["nodeGroupUri"])

	return r
}

// expandClusterConfigGceClusterConfigShieldedInstanceConfigMap expands the contents of ClusterConfigGceClusterConfigShieldedInstanceConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfigShieldedInstanceConfigMap(c *Client, f map[string]ClusterConfigGceClusterConfigShieldedInstanceConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigGceClusterConfigShieldedInstanceConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigGceClusterConfigShieldedInstanceConfigSlice expands the contents of ClusterConfigGceClusterConfigShieldedInstanceConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfigShieldedInstanceConfigSlice(c *Client, f []ClusterConfigGceClusterConfigShieldedInstanceConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigGceClusterConfigShieldedInstanceConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigGceClusterConfigShieldedInstanceConfigMap flattens the contents of ClusterConfigGceClusterConfigShieldedInstanceConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfigShieldedInstanceConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigGceClusterConfigShieldedInstanceConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}

	items := make(map[string]ClusterConfigGceClusterConfigShieldedInstanceConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigGceClusterConfigShieldedInstanceConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigGceClusterConfigShieldedInstanceConfigSlice flattens the contents of ClusterConfigGceClusterConfigShieldedInstanceConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfigShieldedInstanceConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigGceClusterConfigShieldedInstanceConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}

	items := make([]ClusterConfigGceClusterConfigShieldedInstanceConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigGceClusterConfigShieldedInstanceConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigGceClusterConfigShieldedInstanceConfig expands an instance of ClusterConfigGceClusterConfigShieldedInstanceConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfigShieldedInstanceConfig(c *Client, f *ClusterConfigGceClusterConfigShieldedInstanceConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.EnableSecureBoot; !dcl.IsEmptyValueIndirect(v) {
		m["enableSecureBoot"] = v
	}
	if v := f.EnableVtpm; !dcl.IsEmptyValueIndirect(v) {
		m["enableVtpm"] = v
	}
	if v := f.EnableIntegrityMonitoring; !dcl.IsEmptyValueIndirect(v) {
		m["enableIntegrityMonitoring"] = v
	}

	return m, nil
}

// flattenClusterConfigGceClusterConfigShieldedInstanceConfig flattens an instance of ClusterConfigGceClusterConfigShieldedInstanceConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfigShieldedInstanceConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigGceClusterConfigShieldedInstanceConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigGceClusterConfigShieldedInstanceConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigGceClusterConfigShieldedInstanceConfig
	}
	r.EnableSecureBoot = dcl.FlattenBool(m["enableSecureBoot"])
	r.EnableVtpm = dcl.FlattenBool(m["enableVtpm"])
	r.EnableIntegrityMonitoring = dcl.FlattenBool(m["enableIntegrityMonitoring"])

	return r
}

// expandClusterConfigGceClusterConfigConfidentialInstanceConfigMap expands the contents of ClusterConfigGceClusterConfigConfidentialInstanceConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfigConfidentialInstanceConfigMap(c *Client, f map[string]ClusterConfigGceClusterConfigConfidentialInstanceConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigGceClusterConfigConfidentialInstanceConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigGceClusterConfigConfidentialInstanceConfigSlice expands the contents of ClusterConfigGceClusterConfigConfidentialInstanceConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfigConfidentialInstanceConfigSlice(c *Client, f []ClusterConfigGceClusterConfigConfidentialInstanceConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigGceClusterConfigConfidentialInstanceConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigGceClusterConfigConfidentialInstanceConfigMap flattens the contents of ClusterConfigGceClusterConfigConfidentialInstanceConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfigConfidentialInstanceConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigGceClusterConfigConfidentialInstanceConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigGceClusterConfigConfidentialInstanceConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigGceClusterConfigConfidentialInstanceConfig{}
	}

	items := make(map[string]ClusterConfigGceClusterConfigConfidentialInstanceConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigGceClusterConfigConfidentialInstanceConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigGceClusterConfigConfidentialInstanceConfigSlice flattens the contents of ClusterConfigGceClusterConfigConfidentialInstanceConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfigConfidentialInstanceConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigGceClusterConfigConfidentialInstanceConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigGceClusterConfigConfidentialInstanceConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigGceClusterConfigConfidentialInstanceConfig{}
	}

	items := make([]ClusterConfigGceClusterConfigConfidentialInstanceConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigGceClusterConfigConfidentialInstanceConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigGceClusterConfigConfidentialInstanceConfig expands an instance of ClusterConfigGceClusterConfigConfidentialInstanceConfig into a JSON
// request object.
func expandClusterConfigGceClusterConfigConfidentialInstanceConfig(c *Client, f *ClusterConfigGceClusterConfigConfidentialInstanceConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.EnableConfidentialCompute; !dcl.IsEmptyValueIndirect(v) {
		m["enableConfidentialCompute"] = v
	}

	return m, nil
}

// flattenClusterConfigGceClusterConfigConfidentialInstanceConfig flattens an instance of ClusterConfigGceClusterConfigConfidentialInstanceConfig from a JSON
// response object.
func flattenClusterConfigGceClusterConfigConfidentialInstanceConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigGceClusterConfigConfidentialInstanceConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigGceClusterConfigConfidentialInstanceConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigGceClusterConfigConfidentialInstanceConfig
	}
	r.EnableConfidentialCompute = dcl.FlattenBool(m["enableConfidentialCompute"])

	return r
}

// expandClusterConfigMasterConfigMap expands the contents of ClusterConfigMasterConfig into a JSON
// request object.
func expandClusterConfigMasterConfigMap(c *Client, f map[string]ClusterConfigMasterConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigMasterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigMasterConfigSlice expands the contents of ClusterConfigMasterConfig into a JSON
// request object.
func expandClusterConfigMasterConfigSlice(c *Client, f []ClusterConfigMasterConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigMasterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigMasterConfigMap flattens the contents of ClusterConfigMasterConfig from a JSON
// response object.
func flattenClusterConfigMasterConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigMasterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigMasterConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigMasterConfig{}
	}

	items := make(map[string]ClusterConfigMasterConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigMasterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigMasterConfigSlice flattens the contents of ClusterConfigMasterConfig from a JSON
// response object.
func flattenClusterConfigMasterConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigMasterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigMasterConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigMasterConfig{}
	}

	items := make([]ClusterConfigMasterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigMasterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigMasterConfig expands an instance of ClusterConfigMasterConfig into a JSON
// request object.
func expandClusterConfigMasterConfig(c *Client, f *ClusterConfigMasterConfig, res *Cluster) (map[string]interface{}, error) {
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
	if v, err := expandClusterConfigMasterConfigDiskConfig(c, f.DiskConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding DiskConfig into diskConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["diskConfig"] = v
	}
	if v := f.Preemptibility; !dcl.IsEmptyValueIndirect(v) {
		m["preemptibility"] = v
	}
	if v, err := expandClusterConfigMasterConfigAcceleratorsSlice(c, f.Accelerators, res); err != nil {
		return nil, fmt.Errorf("error expanding Accelerators into accelerators: %w", err)
	} else if v != nil {
		m["accelerators"] = v
	}
	if v := f.MinCpuPlatform; !dcl.IsEmptyValueIndirect(v) {
		m["minCpuPlatform"] = v
	}

	return m, nil
}

// flattenClusterConfigMasterConfig flattens an instance of ClusterConfigMasterConfig from a JSON
// response object.
func flattenClusterConfigMasterConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigMasterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigMasterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigMasterConfig
	}
	r.NumInstances = dcl.FlattenInteger(m["numInstances"])
	r.InstanceNames = dcl.FlattenStringSlice(m["instanceNames"])
	r.Image = dcl.FlattenString(m["imageUri"])
	r.MachineType = dcl.FlattenString(m["machineTypeUri"])
	r.DiskConfig = flattenClusterConfigMasterConfigDiskConfig(c, m["diskConfig"], res)
	r.IsPreemptible = dcl.FlattenBool(m["isPreemptible"])
	r.Preemptibility = flattenClusterConfigMasterConfigPreemptibilityEnum(m["preemptibility"])
	r.ManagedGroupConfig = flattenClusterConfigMasterConfigManagedGroupConfig(c, m["managedGroupConfig"], res)
	r.Accelerators = flattenClusterConfigMasterConfigAcceleratorsSlice(c, m["accelerators"], res)
	r.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])
	r.InstanceReferences = flattenClusterConfigMasterConfigInstanceReferencesSlice(c, m["instanceReferences"], res)

	return r
}

// expandClusterConfigMasterConfigDiskConfigMap expands the contents of ClusterConfigMasterConfigDiskConfig into a JSON
// request object.
func expandClusterConfigMasterConfigDiskConfigMap(c *Client, f map[string]ClusterConfigMasterConfigDiskConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigMasterConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigMasterConfigDiskConfigSlice expands the contents of ClusterConfigMasterConfigDiskConfig into a JSON
// request object.
func expandClusterConfigMasterConfigDiskConfigSlice(c *Client, f []ClusterConfigMasterConfigDiskConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigMasterConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigMasterConfigDiskConfigMap flattens the contents of ClusterConfigMasterConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigMasterConfigDiskConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigMasterConfigDiskConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigMasterConfigDiskConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigMasterConfigDiskConfig{}
	}

	items := make(map[string]ClusterConfigMasterConfigDiskConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigMasterConfigDiskConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigMasterConfigDiskConfigSlice flattens the contents of ClusterConfigMasterConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigMasterConfigDiskConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigMasterConfigDiskConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigMasterConfigDiskConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigMasterConfigDiskConfig{}
	}

	items := make([]ClusterConfigMasterConfigDiskConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigMasterConfigDiskConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigMasterConfigDiskConfig expands an instance of ClusterConfigMasterConfigDiskConfig into a JSON
// request object.
func expandClusterConfigMasterConfigDiskConfig(c *Client, f *ClusterConfigMasterConfigDiskConfig, res *Cluster) (map[string]interface{}, error) {
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
	if v := f.LocalSsdInterface; !dcl.IsEmptyValueIndirect(v) {
		m["localSsdInterface"] = v
	}

	return m, nil
}

// flattenClusterConfigMasterConfigDiskConfig flattens an instance of ClusterConfigMasterConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigMasterConfigDiskConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigMasterConfigDiskConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigMasterConfigDiskConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigMasterConfigDiskConfig
	}
	r.BootDiskType = dcl.FlattenString(m["bootDiskType"])
	r.BootDiskSizeGb = dcl.FlattenInteger(m["bootDiskSizeGb"])
	r.NumLocalSsds = dcl.FlattenInteger(m["numLocalSsds"])
	r.LocalSsdInterface = dcl.FlattenString(m["localSsdInterface"])

	return r
}

// expandClusterConfigMasterConfigManagedGroupConfigMap expands the contents of ClusterConfigMasterConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigMasterConfigManagedGroupConfigMap(c *Client, f map[string]ClusterConfigMasterConfigManagedGroupConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigMasterConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigMasterConfigManagedGroupConfigSlice expands the contents of ClusterConfigMasterConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigMasterConfigManagedGroupConfigSlice(c *Client, f []ClusterConfigMasterConfigManagedGroupConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigMasterConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigMasterConfigManagedGroupConfigMap flattens the contents of ClusterConfigMasterConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigMasterConfigManagedGroupConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigMasterConfigManagedGroupConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigMasterConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigMasterConfigManagedGroupConfig{}
	}

	items := make(map[string]ClusterConfigMasterConfigManagedGroupConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigMasterConfigManagedGroupConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigMasterConfigManagedGroupConfigSlice flattens the contents of ClusterConfigMasterConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigMasterConfigManagedGroupConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigMasterConfigManagedGroupConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigMasterConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigMasterConfigManagedGroupConfig{}
	}

	items := make([]ClusterConfigMasterConfigManagedGroupConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigMasterConfigManagedGroupConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigMasterConfigManagedGroupConfig expands an instance of ClusterConfigMasterConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigMasterConfigManagedGroupConfig(c *Client, f *ClusterConfigMasterConfigManagedGroupConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenClusterConfigMasterConfigManagedGroupConfig flattens an instance of ClusterConfigMasterConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigMasterConfigManagedGroupConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigMasterConfigManagedGroupConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigMasterConfigManagedGroupConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigMasterConfigManagedGroupConfig
	}
	r.InstanceTemplateName = dcl.FlattenString(m["instanceTemplateName"])
	r.InstanceGroupManagerName = dcl.FlattenString(m["instanceGroupManagerName"])

	return r
}

// expandClusterConfigMasterConfigAcceleratorsMap expands the contents of ClusterConfigMasterConfigAccelerators into a JSON
// request object.
func expandClusterConfigMasterConfigAcceleratorsMap(c *Client, f map[string]ClusterConfigMasterConfigAccelerators, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigMasterConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigMasterConfigAcceleratorsSlice expands the contents of ClusterConfigMasterConfigAccelerators into a JSON
// request object.
func expandClusterConfigMasterConfigAcceleratorsSlice(c *Client, f []ClusterConfigMasterConfigAccelerators, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigMasterConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigMasterConfigAcceleratorsMap flattens the contents of ClusterConfigMasterConfigAccelerators from a JSON
// response object.
func flattenClusterConfigMasterConfigAcceleratorsMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigMasterConfigAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigMasterConfigAccelerators{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigMasterConfigAccelerators{}
	}

	items := make(map[string]ClusterConfigMasterConfigAccelerators)
	for k, item := range a {
		items[k] = *flattenClusterConfigMasterConfigAccelerators(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigMasterConfigAcceleratorsSlice flattens the contents of ClusterConfigMasterConfigAccelerators from a JSON
// response object.
func flattenClusterConfigMasterConfigAcceleratorsSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigMasterConfigAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigMasterConfigAccelerators{}
	}

	if len(a) == 0 {
		return []ClusterConfigMasterConfigAccelerators{}
	}

	items := make([]ClusterConfigMasterConfigAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigMasterConfigAccelerators(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigMasterConfigAccelerators expands an instance of ClusterConfigMasterConfigAccelerators into a JSON
// request object.
func expandClusterConfigMasterConfigAccelerators(c *Client, f *ClusterConfigMasterConfigAccelerators, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
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

// flattenClusterConfigMasterConfigAccelerators flattens an instance of ClusterConfigMasterConfigAccelerators from a JSON
// response object.
func flattenClusterConfigMasterConfigAccelerators(c *Client, i interface{}, res *Cluster) *ClusterConfigMasterConfigAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigMasterConfigAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigMasterConfigAccelerators
	}
	r.AcceleratorType = dcl.FlattenString(m["acceleratorTypeUri"])
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])

	return r
}

// expandClusterConfigMasterConfigInstanceReferencesMap expands the contents of ClusterConfigMasterConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigMasterConfigInstanceReferencesMap(c *Client, f map[string]ClusterConfigMasterConfigInstanceReferences, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigMasterConfigInstanceReferences(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigMasterConfigInstanceReferencesSlice expands the contents of ClusterConfigMasterConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigMasterConfigInstanceReferencesSlice(c *Client, f []ClusterConfigMasterConfigInstanceReferences, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigMasterConfigInstanceReferences(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigMasterConfigInstanceReferencesMap flattens the contents of ClusterConfigMasterConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigMasterConfigInstanceReferencesMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigMasterConfigInstanceReferences {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigMasterConfigInstanceReferences{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigMasterConfigInstanceReferences{}
	}

	items := make(map[string]ClusterConfigMasterConfigInstanceReferences)
	for k, item := range a {
		items[k] = *flattenClusterConfigMasterConfigInstanceReferences(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigMasterConfigInstanceReferencesSlice flattens the contents of ClusterConfigMasterConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigMasterConfigInstanceReferencesSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigMasterConfigInstanceReferences {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigMasterConfigInstanceReferences{}
	}

	if len(a) == 0 {
		return []ClusterConfigMasterConfigInstanceReferences{}
	}

	items := make([]ClusterConfigMasterConfigInstanceReferences, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigMasterConfigInstanceReferences(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigMasterConfigInstanceReferences expands an instance of ClusterConfigMasterConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigMasterConfigInstanceReferences(c *Client, f *ClusterConfigMasterConfigInstanceReferences, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.InstanceName; !dcl.IsEmptyValueIndirect(v) {
		m["instanceName"] = v
	}
	if v := f.InstanceId; !dcl.IsEmptyValueIndirect(v) {
		m["instanceId"] = v
	}
	if v := f.PublicKey; !dcl.IsEmptyValueIndirect(v) {
		m["publicKey"] = v
	}
	if v := f.PublicEciesKey; !dcl.IsEmptyValueIndirect(v) {
		m["publicEciesKey"] = v
	}

	return m, nil
}

// flattenClusterConfigMasterConfigInstanceReferences flattens an instance of ClusterConfigMasterConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigMasterConfigInstanceReferences(c *Client, i interface{}, res *Cluster) *ClusterConfigMasterConfigInstanceReferences {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigMasterConfigInstanceReferences{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigMasterConfigInstanceReferences
	}
	r.InstanceName = dcl.FlattenString(m["instanceName"])
	r.InstanceId = dcl.FlattenString(m["instanceId"])
	r.PublicKey = dcl.FlattenString(m["publicKey"])
	r.PublicEciesKey = dcl.FlattenString(m["publicEciesKey"])

	return r
}

// expandClusterConfigWorkerConfigMap expands the contents of ClusterConfigWorkerConfig into a JSON
// request object.
func expandClusterConfigWorkerConfigMap(c *Client, f map[string]ClusterConfigWorkerConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigWorkerConfigSlice expands the contents of ClusterConfigWorkerConfig into a JSON
// request object.
func expandClusterConfigWorkerConfigSlice(c *Client, f []ClusterConfigWorkerConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigWorkerConfigMap flattens the contents of ClusterConfigWorkerConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigWorkerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigWorkerConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigWorkerConfig{}
	}

	items := make(map[string]ClusterConfigWorkerConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigWorkerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigWorkerConfigSlice flattens the contents of ClusterConfigWorkerConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigWorkerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigWorkerConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigWorkerConfig{}
	}

	items := make([]ClusterConfigWorkerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigWorkerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigWorkerConfig expands an instance of ClusterConfigWorkerConfig into a JSON
// request object.
func expandClusterConfigWorkerConfig(c *Client, f *ClusterConfigWorkerConfig, res *Cluster) (map[string]interface{}, error) {
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
	if v, err := expandClusterConfigWorkerConfigDiskConfig(c, f.DiskConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding DiskConfig into diskConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["diskConfig"] = v
	}
	if v := f.Preemptibility; !dcl.IsEmptyValueIndirect(v) {
		m["preemptibility"] = v
	}
	if v, err := expandClusterConfigWorkerConfigAcceleratorsSlice(c, f.Accelerators, res); err != nil {
		return nil, fmt.Errorf("error expanding Accelerators into accelerators: %w", err)
	} else if v != nil {
		m["accelerators"] = v
	}
	if v := f.MinCpuPlatform; !dcl.IsEmptyValueIndirect(v) {
		m["minCpuPlatform"] = v
	}

	return m, nil
}

// flattenClusterConfigWorkerConfig flattens an instance of ClusterConfigWorkerConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigWorkerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigWorkerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigWorkerConfig
	}
	r.NumInstances = dcl.FlattenInteger(m["numInstances"])
	r.InstanceNames = dcl.FlattenStringSlice(m["instanceNames"])
	r.Image = dcl.FlattenString(m["imageUri"])
	r.MachineType = dcl.FlattenString(m["machineTypeUri"])
	r.DiskConfig = flattenClusterConfigWorkerConfigDiskConfig(c, m["diskConfig"], res)
	r.IsPreemptible = dcl.FlattenBool(m["isPreemptible"])
	r.Preemptibility = flattenClusterConfigWorkerConfigPreemptibilityEnum(m["preemptibility"])
	r.ManagedGroupConfig = flattenClusterConfigWorkerConfigManagedGroupConfig(c, m["managedGroupConfig"], res)
	r.Accelerators = flattenClusterConfigWorkerConfigAcceleratorsSlice(c, m["accelerators"], res)
	r.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])
	r.InstanceReferences = flattenClusterConfigWorkerConfigInstanceReferencesSlice(c, m["instanceReferences"], res)

	return r
}

// expandClusterConfigWorkerConfigDiskConfigMap expands the contents of ClusterConfigWorkerConfigDiskConfig into a JSON
// request object.
func expandClusterConfigWorkerConfigDiskConfigMap(c *Client, f map[string]ClusterConfigWorkerConfigDiskConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigWorkerConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigWorkerConfigDiskConfigSlice expands the contents of ClusterConfigWorkerConfigDiskConfig into a JSON
// request object.
func expandClusterConfigWorkerConfigDiskConfigSlice(c *Client, f []ClusterConfigWorkerConfigDiskConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigWorkerConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigWorkerConfigDiskConfigMap flattens the contents of ClusterConfigWorkerConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfigDiskConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigWorkerConfigDiskConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigWorkerConfigDiskConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigWorkerConfigDiskConfig{}
	}

	items := make(map[string]ClusterConfigWorkerConfigDiskConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigWorkerConfigDiskConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigWorkerConfigDiskConfigSlice flattens the contents of ClusterConfigWorkerConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfigDiskConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigWorkerConfigDiskConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigWorkerConfigDiskConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigWorkerConfigDiskConfig{}
	}

	items := make([]ClusterConfigWorkerConfigDiskConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigWorkerConfigDiskConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigWorkerConfigDiskConfig expands an instance of ClusterConfigWorkerConfigDiskConfig into a JSON
// request object.
func expandClusterConfigWorkerConfigDiskConfig(c *Client, f *ClusterConfigWorkerConfigDiskConfig, res *Cluster) (map[string]interface{}, error) {
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
	if v := f.LocalSsdInterface; !dcl.IsEmptyValueIndirect(v) {
		m["localSsdInterface"] = v
	}

	return m, nil
}

// flattenClusterConfigWorkerConfigDiskConfig flattens an instance of ClusterConfigWorkerConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfigDiskConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigWorkerConfigDiskConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigWorkerConfigDiskConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigWorkerConfigDiskConfig
	}
	r.BootDiskType = dcl.FlattenString(m["bootDiskType"])
	r.BootDiskSizeGb = dcl.FlattenInteger(m["bootDiskSizeGb"])
	r.NumLocalSsds = dcl.FlattenInteger(m["numLocalSsds"])
	r.LocalSsdInterface = dcl.FlattenString(m["localSsdInterface"])

	return r
}

// expandClusterConfigWorkerConfigManagedGroupConfigMap expands the contents of ClusterConfigWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigWorkerConfigManagedGroupConfigMap(c *Client, f map[string]ClusterConfigWorkerConfigManagedGroupConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigWorkerConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigWorkerConfigManagedGroupConfigSlice expands the contents of ClusterConfigWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigWorkerConfigManagedGroupConfigSlice(c *Client, f []ClusterConfigWorkerConfigManagedGroupConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigWorkerConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigWorkerConfigManagedGroupConfigMap flattens the contents of ClusterConfigWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfigManagedGroupConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigWorkerConfigManagedGroupConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigWorkerConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigWorkerConfigManagedGroupConfig{}
	}

	items := make(map[string]ClusterConfigWorkerConfigManagedGroupConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigWorkerConfigManagedGroupConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigWorkerConfigManagedGroupConfigSlice flattens the contents of ClusterConfigWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfigManagedGroupConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigWorkerConfigManagedGroupConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigWorkerConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigWorkerConfigManagedGroupConfig{}
	}

	items := make([]ClusterConfigWorkerConfigManagedGroupConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigWorkerConfigManagedGroupConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigWorkerConfigManagedGroupConfig expands an instance of ClusterConfigWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigWorkerConfigManagedGroupConfig(c *Client, f *ClusterConfigWorkerConfigManagedGroupConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenClusterConfigWorkerConfigManagedGroupConfig flattens an instance of ClusterConfigWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigWorkerConfigManagedGroupConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigWorkerConfigManagedGroupConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigWorkerConfigManagedGroupConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigWorkerConfigManagedGroupConfig
	}
	r.InstanceTemplateName = dcl.FlattenString(m["instanceTemplateName"])
	r.InstanceGroupManagerName = dcl.FlattenString(m["instanceGroupManagerName"])

	return r
}

// expandClusterConfigWorkerConfigAcceleratorsMap expands the contents of ClusterConfigWorkerConfigAccelerators into a JSON
// request object.
func expandClusterConfigWorkerConfigAcceleratorsMap(c *Client, f map[string]ClusterConfigWorkerConfigAccelerators, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigWorkerConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigWorkerConfigAcceleratorsSlice expands the contents of ClusterConfigWorkerConfigAccelerators into a JSON
// request object.
func expandClusterConfigWorkerConfigAcceleratorsSlice(c *Client, f []ClusterConfigWorkerConfigAccelerators, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigWorkerConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigWorkerConfigAcceleratorsMap flattens the contents of ClusterConfigWorkerConfigAccelerators from a JSON
// response object.
func flattenClusterConfigWorkerConfigAcceleratorsMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigWorkerConfigAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigWorkerConfigAccelerators{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigWorkerConfigAccelerators{}
	}

	items := make(map[string]ClusterConfigWorkerConfigAccelerators)
	for k, item := range a {
		items[k] = *flattenClusterConfigWorkerConfigAccelerators(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigWorkerConfigAcceleratorsSlice flattens the contents of ClusterConfigWorkerConfigAccelerators from a JSON
// response object.
func flattenClusterConfigWorkerConfigAcceleratorsSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigWorkerConfigAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigWorkerConfigAccelerators{}
	}

	if len(a) == 0 {
		return []ClusterConfigWorkerConfigAccelerators{}
	}

	items := make([]ClusterConfigWorkerConfigAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigWorkerConfigAccelerators(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigWorkerConfigAccelerators expands an instance of ClusterConfigWorkerConfigAccelerators into a JSON
// request object.
func expandClusterConfigWorkerConfigAccelerators(c *Client, f *ClusterConfigWorkerConfigAccelerators, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
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

// flattenClusterConfigWorkerConfigAccelerators flattens an instance of ClusterConfigWorkerConfigAccelerators from a JSON
// response object.
func flattenClusterConfigWorkerConfigAccelerators(c *Client, i interface{}, res *Cluster) *ClusterConfigWorkerConfigAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigWorkerConfigAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigWorkerConfigAccelerators
	}
	r.AcceleratorType = dcl.FlattenString(m["acceleratorTypeUri"])
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])

	return r
}

// expandClusterConfigWorkerConfigInstanceReferencesMap expands the contents of ClusterConfigWorkerConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigWorkerConfigInstanceReferencesMap(c *Client, f map[string]ClusterConfigWorkerConfigInstanceReferences, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigWorkerConfigInstanceReferences(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigWorkerConfigInstanceReferencesSlice expands the contents of ClusterConfigWorkerConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigWorkerConfigInstanceReferencesSlice(c *Client, f []ClusterConfigWorkerConfigInstanceReferences, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigWorkerConfigInstanceReferences(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigWorkerConfigInstanceReferencesMap flattens the contents of ClusterConfigWorkerConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigWorkerConfigInstanceReferencesMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigWorkerConfigInstanceReferences {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigWorkerConfigInstanceReferences{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigWorkerConfigInstanceReferences{}
	}

	items := make(map[string]ClusterConfigWorkerConfigInstanceReferences)
	for k, item := range a {
		items[k] = *flattenClusterConfigWorkerConfigInstanceReferences(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigWorkerConfigInstanceReferencesSlice flattens the contents of ClusterConfigWorkerConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigWorkerConfigInstanceReferencesSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigWorkerConfigInstanceReferences {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigWorkerConfigInstanceReferences{}
	}

	if len(a) == 0 {
		return []ClusterConfigWorkerConfigInstanceReferences{}
	}

	items := make([]ClusterConfigWorkerConfigInstanceReferences, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigWorkerConfigInstanceReferences(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigWorkerConfigInstanceReferences expands an instance of ClusterConfigWorkerConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigWorkerConfigInstanceReferences(c *Client, f *ClusterConfigWorkerConfigInstanceReferences, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.InstanceName; !dcl.IsEmptyValueIndirect(v) {
		m["instanceName"] = v
	}
	if v := f.InstanceId; !dcl.IsEmptyValueIndirect(v) {
		m["instanceId"] = v
	}
	if v := f.PublicKey; !dcl.IsEmptyValueIndirect(v) {
		m["publicKey"] = v
	}
	if v := f.PublicEciesKey; !dcl.IsEmptyValueIndirect(v) {
		m["publicEciesKey"] = v
	}

	return m, nil
}

// flattenClusterConfigWorkerConfigInstanceReferences flattens an instance of ClusterConfigWorkerConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigWorkerConfigInstanceReferences(c *Client, i interface{}, res *Cluster) *ClusterConfigWorkerConfigInstanceReferences {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigWorkerConfigInstanceReferences{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigWorkerConfigInstanceReferences
	}
	r.InstanceName = dcl.FlattenString(m["instanceName"])
	r.InstanceId = dcl.FlattenString(m["instanceId"])
	r.PublicKey = dcl.FlattenString(m["publicKey"])
	r.PublicEciesKey = dcl.FlattenString(m["publicEciesKey"])

	return r
}

// expandClusterConfigSecondaryWorkerConfigMap expands the contents of ClusterConfigSecondaryWorkerConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigMap(c *Client, f map[string]ClusterConfigSecondaryWorkerConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSecondaryWorkerConfigSlice expands the contents of ClusterConfigSecondaryWorkerConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigSlice(c *Client, f []ClusterConfigSecondaryWorkerConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSecondaryWorkerConfigMap flattens the contents of ClusterConfigSecondaryWorkerConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecondaryWorkerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecondaryWorkerConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecondaryWorkerConfig{}
	}

	items := make(map[string]ClusterConfigSecondaryWorkerConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecondaryWorkerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSecondaryWorkerConfigSlice flattens the contents of ClusterConfigSecondaryWorkerConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecondaryWorkerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecondaryWorkerConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecondaryWorkerConfig{}
	}

	items := make([]ClusterConfigSecondaryWorkerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecondaryWorkerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSecondaryWorkerConfig expands an instance of ClusterConfigSecondaryWorkerConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfig(c *Client, f *ClusterConfigSecondaryWorkerConfig, res *Cluster) (map[string]interface{}, error) {
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
	if v, err := expandClusterConfigSecondaryWorkerConfigDiskConfig(c, f.DiskConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding DiskConfig into diskConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["diskConfig"] = v
	}
	if v := f.Preemptibility; !dcl.IsEmptyValueIndirect(v) {
		m["preemptibility"] = v
	}
	if v, err := expandClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c, f.Accelerators, res); err != nil {
		return nil, fmt.Errorf("error expanding Accelerators into accelerators: %w", err)
	} else if v != nil {
		m["accelerators"] = v
	}
	if v := f.MinCpuPlatform; !dcl.IsEmptyValueIndirect(v) {
		m["minCpuPlatform"] = v
	}

	return m, nil
}

// flattenClusterConfigSecondaryWorkerConfig flattens an instance of ClusterConfigSecondaryWorkerConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigSecondaryWorkerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSecondaryWorkerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSecondaryWorkerConfig
	}
	r.NumInstances = dcl.FlattenInteger(m["numInstances"])
	r.InstanceNames = dcl.FlattenStringSlice(m["instanceNames"])
	r.Image = dcl.FlattenString(m["imageUri"])
	r.MachineType = dcl.FlattenString(m["machineTypeUri"])
	r.DiskConfig = flattenClusterConfigSecondaryWorkerConfigDiskConfig(c, m["diskConfig"], res)
	r.IsPreemptible = dcl.FlattenBool(m["isPreemptible"])
	r.Preemptibility = flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnum(m["preemptibility"])
	r.ManagedGroupConfig = flattenClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, m["managedGroupConfig"], res)
	r.Accelerators = flattenClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c, m["accelerators"], res)
	r.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])
	r.InstanceReferences = flattenClusterConfigSecondaryWorkerConfigInstanceReferencesSlice(c, m["instanceReferences"], res)

	return r
}

// expandClusterConfigSecondaryWorkerConfigDiskConfigMap expands the contents of ClusterConfigSecondaryWorkerConfigDiskConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigDiskConfigMap(c *Client, f map[string]ClusterConfigSecondaryWorkerConfigDiskConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSecondaryWorkerConfigDiskConfigSlice expands the contents of ClusterConfigSecondaryWorkerConfigDiskConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigDiskConfigSlice(c *Client, f []ClusterConfigSecondaryWorkerConfigDiskConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSecondaryWorkerConfigDiskConfigMap flattens the contents of ClusterConfigSecondaryWorkerConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigDiskConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecondaryWorkerConfigDiskConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecondaryWorkerConfigDiskConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecondaryWorkerConfigDiskConfig{}
	}

	items := make(map[string]ClusterConfigSecondaryWorkerConfigDiskConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecondaryWorkerConfigDiskConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSecondaryWorkerConfigDiskConfigSlice flattens the contents of ClusterConfigSecondaryWorkerConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigDiskConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecondaryWorkerConfigDiskConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecondaryWorkerConfigDiskConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecondaryWorkerConfigDiskConfig{}
	}

	items := make([]ClusterConfigSecondaryWorkerConfigDiskConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecondaryWorkerConfigDiskConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSecondaryWorkerConfigDiskConfig expands an instance of ClusterConfigSecondaryWorkerConfigDiskConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigDiskConfig(c *Client, f *ClusterConfigSecondaryWorkerConfigDiskConfig, res *Cluster) (map[string]interface{}, error) {
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
	if v := f.LocalSsdInterface; !dcl.IsEmptyValueIndirect(v) {
		m["localSsdInterface"] = v
	}

	return m, nil
}

// flattenClusterConfigSecondaryWorkerConfigDiskConfig flattens an instance of ClusterConfigSecondaryWorkerConfigDiskConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigDiskConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigSecondaryWorkerConfigDiskConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSecondaryWorkerConfigDiskConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSecondaryWorkerConfigDiskConfig
	}
	r.BootDiskType = dcl.FlattenString(m["bootDiskType"])
	r.BootDiskSizeGb = dcl.FlattenInteger(m["bootDiskSizeGb"])
	r.NumLocalSsds = dcl.FlattenInteger(m["numLocalSsds"])
	r.LocalSsdInterface = dcl.FlattenString(m["localSsdInterface"])

	return r
}

// expandClusterConfigSecondaryWorkerConfigManagedGroupConfigMap expands the contents of ClusterConfigSecondaryWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigManagedGroupConfigMap(c *Client, f map[string]ClusterConfigSecondaryWorkerConfigManagedGroupConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice expands the contents of ClusterConfigSecondaryWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice(c *Client, f []ClusterConfigSecondaryWorkerConfigManagedGroupConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSecondaryWorkerConfigManagedGroupConfigMap flattens the contents of ClusterConfigSecondaryWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigManagedGroupConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}

	items := make(map[string]ClusterConfigSecondaryWorkerConfigManagedGroupConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice flattens the contents of ClusterConfigSecondaryWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}

	items := make([]ClusterConfigSecondaryWorkerConfigManagedGroupConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSecondaryWorkerConfigManagedGroupConfig expands an instance of ClusterConfigSecondaryWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigManagedGroupConfig(c *Client, f *ClusterConfigSecondaryWorkerConfigManagedGroupConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenClusterConfigSecondaryWorkerConfigManagedGroupConfig flattens an instance of ClusterConfigSecondaryWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigManagedGroupConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSecondaryWorkerConfigManagedGroupConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSecondaryWorkerConfigManagedGroupConfig
	}
	r.InstanceTemplateName = dcl.FlattenString(m["instanceTemplateName"])
	r.InstanceGroupManagerName = dcl.FlattenString(m["instanceGroupManagerName"])

	return r
}

// expandClusterConfigSecondaryWorkerConfigAcceleratorsMap expands the contents of ClusterConfigSecondaryWorkerConfigAccelerators into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigAcceleratorsMap(c *Client, f map[string]ClusterConfigSecondaryWorkerConfigAccelerators, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSecondaryWorkerConfigAcceleratorsSlice expands the contents of ClusterConfigSecondaryWorkerConfigAccelerators into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c *Client, f []ClusterConfigSecondaryWorkerConfigAccelerators, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSecondaryWorkerConfigAcceleratorsMap flattens the contents of ClusterConfigSecondaryWorkerConfigAccelerators from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigAcceleratorsMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecondaryWorkerConfigAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecondaryWorkerConfigAccelerators{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecondaryWorkerConfigAccelerators{}
	}

	items := make(map[string]ClusterConfigSecondaryWorkerConfigAccelerators)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecondaryWorkerConfigAccelerators(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSecondaryWorkerConfigAcceleratorsSlice flattens the contents of ClusterConfigSecondaryWorkerConfigAccelerators from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecondaryWorkerConfigAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecondaryWorkerConfigAccelerators{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecondaryWorkerConfigAccelerators{}
	}

	items := make([]ClusterConfigSecondaryWorkerConfigAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecondaryWorkerConfigAccelerators(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSecondaryWorkerConfigAccelerators expands an instance of ClusterConfigSecondaryWorkerConfigAccelerators into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigAccelerators(c *Client, f *ClusterConfigSecondaryWorkerConfigAccelerators, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
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

// flattenClusterConfigSecondaryWorkerConfigAccelerators flattens an instance of ClusterConfigSecondaryWorkerConfigAccelerators from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigAccelerators(c *Client, i interface{}, res *Cluster) *ClusterConfigSecondaryWorkerConfigAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSecondaryWorkerConfigAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSecondaryWorkerConfigAccelerators
	}
	r.AcceleratorType = dcl.FlattenString(m["acceleratorTypeUri"])
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])

	return r
}

// expandClusterConfigSecondaryWorkerConfigInstanceReferencesMap expands the contents of ClusterConfigSecondaryWorkerConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigInstanceReferencesMap(c *Client, f map[string]ClusterConfigSecondaryWorkerConfigInstanceReferences, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfigInstanceReferences(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSecondaryWorkerConfigInstanceReferencesSlice expands the contents of ClusterConfigSecondaryWorkerConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigInstanceReferencesSlice(c *Client, f []ClusterConfigSecondaryWorkerConfigInstanceReferences, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSecondaryWorkerConfigInstanceReferences(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSecondaryWorkerConfigInstanceReferencesMap flattens the contents of ClusterConfigSecondaryWorkerConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigInstanceReferencesMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecondaryWorkerConfigInstanceReferences {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecondaryWorkerConfigInstanceReferences{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecondaryWorkerConfigInstanceReferences{}
	}

	items := make(map[string]ClusterConfigSecondaryWorkerConfigInstanceReferences)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecondaryWorkerConfigInstanceReferences(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSecondaryWorkerConfigInstanceReferencesSlice flattens the contents of ClusterConfigSecondaryWorkerConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigInstanceReferencesSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecondaryWorkerConfigInstanceReferences {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecondaryWorkerConfigInstanceReferences{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecondaryWorkerConfigInstanceReferences{}
	}

	items := make([]ClusterConfigSecondaryWorkerConfigInstanceReferences, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecondaryWorkerConfigInstanceReferences(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSecondaryWorkerConfigInstanceReferences expands an instance of ClusterConfigSecondaryWorkerConfigInstanceReferences into a JSON
// request object.
func expandClusterConfigSecondaryWorkerConfigInstanceReferences(c *Client, f *ClusterConfigSecondaryWorkerConfigInstanceReferences, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.InstanceName; !dcl.IsEmptyValueIndirect(v) {
		m["instanceName"] = v
	}
	if v := f.InstanceId; !dcl.IsEmptyValueIndirect(v) {
		m["instanceId"] = v
	}
	if v := f.PublicKey; !dcl.IsEmptyValueIndirect(v) {
		m["publicKey"] = v
	}
	if v := f.PublicEciesKey; !dcl.IsEmptyValueIndirect(v) {
		m["publicEciesKey"] = v
	}

	return m, nil
}

// flattenClusterConfigSecondaryWorkerConfigInstanceReferences flattens an instance of ClusterConfigSecondaryWorkerConfigInstanceReferences from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigInstanceReferences(c *Client, i interface{}, res *Cluster) *ClusterConfigSecondaryWorkerConfigInstanceReferences {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSecondaryWorkerConfigInstanceReferences{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSecondaryWorkerConfigInstanceReferences
	}
	r.InstanceName = dcl.FlattenString(m["instanceName"])
	r.InstanceId = dcl.FlattenString(m["instanceId"])
	r.PublicKey = dcl.FlattenString(m["publicKey"])
	r.PublicEciesKey = dcl.FlattenString(m["publicEciesKey"])

	return r
}

// expandClusterConfigSoftwareConfigMap expands the contents of ClusterConfigSoftwareConfig into a JSON
// request object.
func expandClusterConfigSoftwareConfigMap(c *Client, f map[string]ClusterConfigSoftwareConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSoftwareConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSoftwareConfigSlice expands the contents of ClusterConfigSoftwareConfig into a JSON
// request object.
func expandClusterConfigSoftwareConfigSlice(c *Client, f []ClusterConfigSoftwareConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSoftwareConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSoftwareConfigMap flattens the contents of ClusterConfigSoftwareConfig from a JSON
// response object.
func flattenClusterConfigSoftwareConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSoftwareConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSoftwareConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSoftwareConfig{}
	}

	items := make(map[string]ClusterConfigSoftwareConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigSoftwareConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSoftwareConfigSlice flattens the contents of ClusterConfigSoftwareConfig from a JSON
// response object.
func flattenClusterConfigSoftwareConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSoftwareConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSoftwareConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigSoftwareConfig{}
	}

	items := make([]ClusterConfigSoftwareConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSoftwareConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSoftwareConfig expands an instance of ClusterConfigSoftwareConfig into a JSON
// request object.
func expandClusterConfigSoftwareConfig(c *Client, f *ClusterConfigSoftwareConfig, res *Cluster) (map[string]interface{}, error) {
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

// flattenClusterConfigSoftwareConfig flattens an instance of ClusterConfigSoftwareConfig from a JSON
// response object.
func flattenClusterConfigSoftwareConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigSoftwareConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSoftwareConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSoftwareConfig
	}
	r.ImageVersion = dcl.FlattenString(m["imageVersion"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.OptionalComponents = flattenClusterConfigSoftwareConfigOptionalComponentsEnumSlice(c, m["optionalComponents"], res)

	return r
}

// expandClusterConfigInitializationActionsMap expands the contents of ClusterConfigInitializationActions into a JSON
// request object.
func expandClusterConfigInitializationActionsMap(c *Client, f map[string]ClusterConfigInitializationActions, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigInitializationActions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigInitializationActionsSlice expands the contents of ClusterConfigInitializationActions into a JSON
// request object.
func expandClusterConfigInitializationActionsSlice(c *Client, f []ClusterConfigInitializationActions, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigInitializationActions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigInitializationActionsMap flattens the contents of ClusterConfigInitializationActions from a JSON
// response object.
func flattenClusterConfigInitializationActionsMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigInitializationActions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigInitializationActions{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigInitializationActions{}
	}

	items := make(map[string]ClusterConfigInitializationActions)
	for k, item := range a {
		items[k] = *flattenClusterConfigInitializationActions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigInitializationActionsSlice flattens the contents of ClusterConfigInitializationActions from a JSON
// response object.
func flattenClusterConfigInitializationActionsSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigInitializationActions {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigInitializationActions{}
	}

	if len(a) == 0 {
		return []ClusterConfigInitializationActions{}
	}

	items := make([]ClusterConfigInitializationActions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigInitializationActions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigInitializationActions expands an instance of ClusterConfigInitializationActions into a JSON
// request object.
func expandClusterConfigInitializationActions(c *Client, f *ClusterConfigInitializationActions, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
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

// flattenClusterConfigInitializationActions flattens an instance of ClusterConfigInitializationActions from a JSON
// response object.
func flattenClusterConfigInitializationActions(c *Client, i interface{}, res *Cluster) *ClusterConfigInitializationActions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigInitializationActions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigInitializationActions
	}
	r.ExecutableFile = dcl.FlattenString(m["executableFile"])
	r.ExecutionTimeout = dcl.FlattenString(m["executionTimeout"])

	return r
}

// expandClusterConfigEncryptionConfigMap expands the contents of ClusterConfigEncryptionConfig into a JSON
// request object.
func expandClusterConfigEncryptionConfigMap(c *Client, f map[string]ClusterConfigEncryptionConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigEncryptionConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigEncryptionConfigSlice expands the contents of ClusterConfigEncryptionConfig into a JSON
// request object.
func expandClusterConfigEncryptionConfigSlice(c *Client, f []ClusterConfigEncryptionConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigEncryptionConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigEncryptionConfigMap flattens the contents of ClusterConfigEncryptionConfig from a JSON
// response object.
func flattenClusterConfigEncryptionConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigEncryptionConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigEncryptionConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigEncryptionConfig{}
	}

	items := make(map[string]ClusterConfigEncryptionConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigEncryptionConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigEncryptionConfigSlice flattens the contents of ClusterConfigEncryptionConfig from a JSON
// response object.
func flattenClusterConfigEncryptionConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigEncryptionConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigEncryptionConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigEncryptionConfig{}
	}

	items := make([]ClusterConfigEncryptionConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigEncryptionConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigEncryptionConfig expands an instance of ClusterConfigEncryptionConfig into a JSON
// request object.
func expandClusterConfigEncryptionConfig(c *Client, f *ClusterConfigEncryptionConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GcePdKmsKeyName; !dcl.IsEmptyValueIndirect(v) {
		m["gcePdKmsKeyName"] = v
	}

	return m, nil
}

// flattenClusterConfigEncryptionConfig flattens an instance of ClusterConfigEncryptionConfig from a JSON
// response object.
func flattenClusterConfigEncryptionConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigEncryptionConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigEncryptionConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigEncryptionConfig
	}
	r.GcePdKmsKeyName = dcl.FlattenString(m["gcePdKmsKeyName"])

	return r
}

// expandClusterConfigAutoscalingConfigMap expands the contents of ClusterConfigAutoscalingConfig into a JSON
// request object.
func expandClusterConfigAutoscalingConfigMap(c *Client, f map[string]ClusterConfigAutoscalingConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigAutoscalingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigAutoscalingConfigSlice expands the contents of ClusterConfigAutoscalingConfig into a JSON
// request object.
func expandClusterConfigAutoscalingConfigSlice(c *Client, f []ClusterConfigAutoscalingConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigAutoscalingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigAutoscalingConfigMap flattens the contents of ClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenClusterConfigAutoscalingConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigAutoscalingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigAutoscalingConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigAutoscalingConfig{}
	}

	items := make(map[string]ClusterConfigAutoscalingConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigAutoscalingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigAutoscalingConfigSlice flattens the contents of ClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenClusterConfigAutoscalingConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigAutoscalingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigAutoscalingConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigAutoscalingConfig{}
	}

	items := make([]ClusterConfigAutoscalingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigAutoscalingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigAutoscalingConfig expands an instance of ClusterConfigAutoscalingConfig into a JSON
// request object.
func expandClusterConfigAutoscalingConfig(c *Client, f *ClusterConfigAutoscalingConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Policy; !dcl.IsEmptyValueIndirect(v) {
		m["policyUri"] = v
	}

	return m, nil
}

// flattenClusterConfigAutoscalingConfig flattens an instance of ClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenClusterConfigAutoscalingConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigAutoscalingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigAutoscalingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigAutoscalingConfig
	}
	r.Policy = dcl.FlattenString(m["policyUri"])

	return r
}

// expandClusterConfigSecurityConfigMap expands the contents of ClusterConfigSecurityConfig into a JSON
// request object.
func expandClusterConfigSecurityConfigMap(c *Client, f map[string]ClusterConfigSecurityConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSecurityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSecurityConfigSlice expands the contents of ClusterConfigSecurityConfig into a JSON
// request object.
func expandClusterConfigSecurityConfigSlice(c *Client, f []ClusterConfigSecurityConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSecurityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSecurityConfigMap flattens the contents of ClusterConfigSecurityConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecurityConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecurityConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecurityConfig{}
	}

	items := make(map[string]ClusterConfigSecurityConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecurityConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSecurityConfigSlice flattens the contents of ClusterConfigSecurityConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecurityConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecurityConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecurityConfig{}
	}

	items := make([]ClusterConfigSecurityConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecurityConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSecurityConfig expands an instance of ClusterConfigSecurityConfig into a JSON
// request object.
func expandClusterConfigSecurityConfig(c *Client, f *ClusterConfigSecurityConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandClusterConfigSecurityConfigKerberosConfig(c, f.KerberosConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding KerberosConfig into kerberosConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["kerberosConfig"] = v
	}
	if v, err := expandClusterConfigSecurityConfigIdentityConfig(c, f.IdentityConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding IdentityConfig into identityConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["identityConfig"] = v
	}

	return m, nil
}

// flattenClusterConfigSecurityConfig flattens an instance of ClusterConfigSecurityConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigSecurityConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSecurityConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSecurityConfig
	}
	r.KerberosConfig = flattenClusterConfigSecurityConfigKerberosConfig(c, m["kerberosConfig"], res)
	r.IdentityConfig = flattenClusterConfigSecurityConfigIdentityConfig(c, m["identityConfig"], res)

	return r
}

// expandClusterConfigSecurityConfigKerberosConfigMap expands the contents of ClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandClusterConfigSecurityConfigKerberosConfigMap(c *Client, f map[string]ClusterConfigSecurityConfigKerberosConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSecurityConfigKerberosConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSecurityConfigKerberosConfigSlice expands the contents of ClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandClusterConfigSecurityConfigKerberosConfigSlice(c *Client, f []ClusterConfigSecurityConfigKerberosConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSecurityConfigKerberosConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSecurityConfigKerberosConfigMap flattens the contents of ClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfigKerberosConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecurityConfigKerberosConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecurityConfigKerberosConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecurityConfigKerberosConfig{}
	}

	items := make(map[string]ClusterConfigSecurityConfigKerberosConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecurityConfigKerberosConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSecurityConfigKerberosConfigSlice flattens the contents of ClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfigKerberosConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecurityConfigKerberosConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecurityConfigKerberosConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecurityConfigKerberosConfig{}
	}

	items := make([]ClusterConfigSecurityConfigKerberosConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecurityConfigKerberosConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSecurityConfigKerberosConfig expands an instance of ClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandClusterConfigSecurityConfigKerberosConfig(c *Client, f *ClusterConfigSecurityConfigKerberosConfig, res *Cluster) (map[string]interface{}, error) {
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

// flattenClusterConfigSecurityConfigKerberosConfig flattens an instance of ClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfigKerberosConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigSecurityConfigKerberosConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSecurityConfigKerberosConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSecurityConfigKerberosConfig
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

// expandClusterConfigSecurityConfigIdentityConfigMap expands the contents of ClusterConfigSecurityConfigIdentityConfig into a JSON
// request object.
func expandClusterConfigSecurityConfigIdentityConfigMap(c *Client, f map[string]ClusterConfigSecurityConfigIdentityConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigSecurityConfigIdentityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigSecurityConfigIdentityConfigSlice expands the contents of ClusterConfigSecurityConfigIdentityConfig into a JSON
// request object.
func expandClusterConfigSecurityConfigIdentityConfigSlice(c *Client, f []ClusterConfigSecurityConfigIdentityConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigSecurityConfigIdentityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigSecurityConfigIdentityConfigMap flattens the contents of ClusterConfigSecurityConfigIdentityConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfigIdentityConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecurityConfigIdentityConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecurityConfigIdentityConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecurityConfigIdentityConfig{}
	}

	items := make(map[string]ClusterConfigSecurityConfigIdentityConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecurityConfigIdentityConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigSecurityConfigIdentityConfigSlice flattens the contents of ClusterConfigSecurityConfigIdentityConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfigIdentityConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecurityConfigIdentityConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecurityConfigIdentityConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecurityConfigIdentityConfig{}
	}

	items := make([]ClusterConfigSecurityConfigIdentityConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecurityConfigIdentityConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigSecurityConfigIdentityConfig expands an instance of ClusterConfigSecurityConfigIdentityConfig into a JSON
// request object.
func expandClusterConfigSecurityConfigIdentityConfig(c *Client, f *ClusterConfigSecurityConfigIdentityConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.UserServiceAccountMapping; !dcl.IsEmptyValueIndirect(v) {
		m["userServiceAccountMapping"] = v
	}

	return m, nil
}

// flattenClusterConfigSecurityConfigIdentityConfig flattens an instance of ClusterConfigSecurityConfigIdentityConfig from a JSON
// response object.
func flattenClusterConfigSecurityConfigIdentityConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigSecurityConfigIdentityConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigSecurityConfigIdentityConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigSecurityConfigIdentityConfig
	}
	r.UserServiceAccountMapping = dcl.FlattenKeyValuePairs(m["userServiceAccountMapping"])

	return r
}

// expandClusterConfigLifecycleConfigMap expands the contents of ClusterConfigLifecycleConfig into a JSON
// request object.
func expandClusterConfigLifecycleConfigMap(c *Client, f map[string]ClusterConfigLifecycleConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigLifecycleConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigLifecycleConfigSlice expands the contents of ClusterConfigLifecycleConfig into a JSON
// request object.
func expandClusterConfigLifecycleConfigSlice(c *Client, f []ClusterConfigLifecycleConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigLifecycleConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigLifecycleConfigMap flattens the contents of ClusterConfigLifecycleConfig from a JSON
// response object.
func flattenClusterConfigLifecycleConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigLifecycleConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigLifecycleConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigLifecycleConfig{}
	}

	items := make(map[string]ClusterConfigLifecycleConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigLifecycleConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigLifecycleConfigSlice flattens the contents of ClusterConfigLifecycleConfig from a JSON
// response object.
func flattenClusterConfigLifecycleConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigLifecycleConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigLifecycleConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigLifecycleConfig{}
	}

	items := make([]ClusterConfigLifecycleConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigLifecycleConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigLifecycleConfig expands an instance of ClusterConfigLifecycleConfig into a JSON
// request object.
func expandClusterConfigLifecycleConfig(c *Client, f *ClusterConfigLifecycleConfig, res *Cluster) (map[string]interface{}, error) {
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

// flattenClusterConfigLifecycleConfig flattens an instance of ClusterConfigLifecycleConfig from a JSON
// response object.
func flattenClusterConfigLifecycleConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigLifecycleConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigLifecycleConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigLifecycleConfig
	}
	r.IdleDeleteTtl = dcl.FlattenString(m["idleDeleteTtl"])
	r.AutoDeleteTime = dcl.FlattenString(m["autoDeleteTime"])
	r.AutoDeleteTtl = dcl.FlattenString(m["autoDeleteTtl"])
	r.IdleStartTime = dcl.FlattenString(m["idleStartTime"])

	return r
}

// expandClusterConfigEndpointConfigMap expands the contents of ClusterConfigEndpointConfig into a JSON
// request object.
func expandClusterConfigEndpointConfigMap(c *Client, f map[string]ClusterConfigEndpointConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigEndpointConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigEndpointConfigSlice expands the contents of ClusterConfigEndpointConfig into a JSON
// request object.
func expandClusterConfigEndpointConfigSlice(c *Client, f []ClusterConfigEndpointConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigEndpointConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigEndpointConfigMap flattens the contents of ClusterConfigEndpointConfig from a JSON
// response object.
func flattenClusterConfigEndpointConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigEndpointConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigEndpointConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigEndpointConfig{}
	}

	items := make(map[string]ClusterConfigEndpointConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigEndpointConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigEndpointConfigSlice flattens the contents of ClusterConfigEndpointConfig from a JSON
// response object.
func flattenClusterConfigEndpointConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigEndpointConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigEndpointConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigEndpointConfig{}
	}

	items := make([]ClusterConfigEndpointConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigEndpointConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigEndpointConfig expands an instance of ClusterConfigEndpointConfig into a JSON
// request object.
func expandClusterConfigEndpointConfig(c *Client, f *ClusterConfigEndpointConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.EnableHttpPortAccess; !dcl.IsEmptyValueIndirect(v) {
		m["enableHttpPortAccess"] = v
	}

	return m, nil
}

// flattenClusterConfigEndpointConfig flattens an instance of ClusterConfigEndpointConfig from a JSON
// response object.
func flattenClusterConfigEndpointConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigEndpointConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigEndpointConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigEndpointConfig
	}
	r.HttpPorts = dcl.FlattenKeyValuePairs(m["httpPorts"])
	r.EnableHttpPortAccess = dcl.FlattenBool(m["enableHttpPortAccess"])

	return r
}

// expandClusterConfigMetastoreConfigMap expands the contents of ClusterConfigMetastoreConfig into a JSON
// request object.
func expandClusterConfigMetastoreConfigMap(c *Client, f map[string]ClusterConfigMetastoreConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigMetastoreConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigMetastoreConfigSlice expands the contents of ClusterConfigMetastoreConfig into a JSON
// request object.
func expandClusterConfigMetastoreConfigSlice(c *Client, f []ClusterConfigMetastoreConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigMetastoreConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigMetastoreConfigMap flattens the contents of ClusterConfigMetastoreConfig from a JSON
// response object.
func flattenClusterConfigMetastoreConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigMetastoreConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigMetastoreConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigMetastoreConfig{}
	}

	items := make(map[string]ClusterConfigMetastoreConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigMetastoreConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigMetastoreConfigSlice flattens the contents of ClusterConfigMetastoreConfig from a JSON
// response object.
func flattenClusterConfigMetastoreConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigMetastoreConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigMetastoreConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigMetastoreConfig{}
	}

	items := make([]ClusterConfigMetastoreConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigMetastoreConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigMetastoreConfig expands an instance of ClusterConfigMetastoreConfig into a JSON
// request object.
func expandClusterConfigMetastoreConfig(c *Client, f *ClusterConfigMetastoreConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DataprocMetastoreService; !dcl.IsEmptyValueIndirect(v) {
		m["dataprocMetastoreService"] = v
	}

	return m, nil
}

// flattenClusterConfigMetastoreConfig flattens an instance of ClusterConfigMetastoreConfig from a JSON
// response object.
func flattenClusterConfigMetastoreConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigMetastoreConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigMetastoreConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigMetastoreConfig
	}
	r.DataprocMetastoreService = dcl.FlattenString(m["dataprocMetastoreService"])

	return r
}

// expandClusterConfigDataprocMetricConfigMap expands the contents of ClusterConfigDataprocMetricConfig into a JSON
// request object.
func expandClusterConfigDataprocMetricConfigMap(c *Client, f map[string]ClusterConfigDataprocMetricConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigDataprocMetricConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigDataprocMetricConfigSlice expands the contents of ClusterConfigDataprocMetricConfig into a JSON
// request object.
func expandClusterConfigDataprocMetricConfigSlice(c *Client, f []ClusterConfigDataprocMetricConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigDataprocMetricConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigDataprocMetricConfigMap flattens the contents of ClusterConfigDataprocMetricConfig from a JSON
// response object.
func flattenClusterConfigDataprocMetricConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigDataprocMetricConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigDataprocMetricConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigDataprocMetricConfig{}
	}

	items := make(map[string]ClusterConfigDataprocMetricConfig)
	for k, item := range a {
		items[k] = *flattenClusterConfigDataprocMetricConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigDataprocMetricConfigSlice flattens the contents of ClusterConfigDataprocMetricConfig from a JSON
// response object.
func flattenClusterConfigDataprocMetricConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigDataprocMetricConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigDataprocMetricConfig{}
	}

	if len(a) == 0 {
		return []ClusterConfigDataprocMetricConfig{}
	}

	items := make([]ClusterConfigDataprocMetricConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigDataprocMetricConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigDataprocMetricConfig expands an instance of ClusterConfigDataprocMetricConfig into a JSON
// request object.
func expandClusterConfigDataprocMetricConfig(c *Client, f *ClusterConfigDataprocMetricConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandClusterConfigDataprocMetricConfigMetricsSlice(c, f.Metrics, res); err != nil {
		return nil, fmt.Errorf("error expanding Metrics into metrics: %w", err)
	} else if v != nil {
		m["metrics"] = v
	}

	return m, nil
}

// flattenClusterConfigDataprocMetricConfig flattens an instance of ClusterConfigDataprocMetricConfig from a JSON
// response object.
func flattenClusterConfigDataprocMetricConfig(c *Client, i interface{}, res *Cluster) *ClusterConfigDataprocMetricConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigDataprocMetricConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigDataprocMetricConfig
	}
	r.Metrics = flattenClusterConfigDataprocMetricConfigMetricsSlice(c, m["metrics"], res)

	return r
}

// expandClusterConfigDataprocMetricConfigMetricsMap expands the contents of ClusterConfigDataprocMetricConfigMetrics into a JSON
// request object.
func expandClusterConfigDataprocMetricConfigMetricsMap(c *Client, f map[string]ClusterConfigDataprocMetricConfigMetrics, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterConfigDataprocMetricConfigMetrics(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterConfigDataprocMetricConfigMetricsSlice expands the contents of ClusterConfigDataprocMetricConfigMetrics into a JSON
// request object.
func expandClusterConfigDataprocMetricConfigMetricsSlice(c *Client, f []ClusterConfigDataprocMetricConfigMetrics, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterConfigDataprocMetricConfigMetrics(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterConfigDataprocMetricConfigMetricsMap flattens the contents of ClusterConfigDataprocMetricConfigMetrics from a JSON
// response object.
func flattenClusterConfigDataprocMetricConfigMetricsMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigDataprocMetricConfigMetrics {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigDataprocMetricConfigMetrics{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigDataprocMetricConfigMetrics{}
	}

	items := make(map[string]ClusterConfigDataprocMetricConfigMetrics)
	for k, item := range a {
		items[k] = *flattenClusterConfigDataprocMetricConfigMetrics(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterConfigDataprocMetricConfigMetricsSlice flattens the contents of ClusterConfigDataprocMetricConfigMetrics from a JSON
// response object.
func flattenClusterConfigDataprocMetricConfigMetricsSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigDataprocMetricConfigMetrics {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigDataprocMetricConfigMetrics{}
	}

	if len(a) == 0 {
		return []ClusterConfigDataprocMetricConfigMetrics{}
	}

	items := make([]ClusterConfigDataprocMetricConfigMetrics, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigDataprocMetricConfigMetrics(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterConfigDataprocMetricConfigMetrics expands an instance of ClusterConfigDataprocMetricConfigMetrics into a JSON
// request object.
func expandClusterConfigDataprocMetricConfigMetrics(c *Client, f *ClusterConfigDataprocMetricConfigMetrics, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MetricSource; !dcl.IsEmptyValueIndirect(v) {
		m["metricSource"] = v
	}
	if v := f.MetricOverrides; v != nil {
		m["metricOverrides"] = v
	}

	return m, nil
}

// flattenClusterConfigDataprocMetricConfigMetrics flattens an instance of ClusterConfigDataprocMetricConfigMetrics from a JSON
// response object.
func flattenClusterConfigDataprocMetricConfigMetrics(c *Client, i interface{}, res *Cluster) *ClusterConfigDataprocMetricConfigMetrics {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterConfigDataprocMetricConfigMetrics{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterConfigDataprocMetricConfigMetrics
	}
	r.MetricSource = flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnum(m["metricSource"])
	r.MetricOverrides = dcl.FlattenStringSlice(m["metricOverrides"])

	return r
}

// expandClusterStatusMap expands the contents of ClusterStatus into a JSON
// request object.
func expandClusterStatusMap(c *Client, f map[string]ClusterStatus, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterStatus(c, &item, res)
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
func expandClusterStatusSlice(c *Client, f []ClusterStatus, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterStatus(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterStatusMap flattens the contents of ClusterStatus from a JSON
// response object.
func flattenClusterStatusMap(c *Client, i interface{}, res *Cluster) map[string]ClusterStatus {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStatus{}
	}

	if len(a) == 0 {
		return map[string]ClusterStatus{}
	}

	items := make(map[string]ClusterStatus)
	for k, item := range a {
		items[k] = *flattenClusterStatus(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterStatusSlice flattens the contents of ClusterStatus from a JSON
// response object.
func flattenClusterStatusSlice(c *Client, i interface{}, res *Cluster) []ClusterStatus {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStatus{}
	}

	if len(a) == 0 {
		return []ClusterStatus{}
	}

	items := make([]ClusterStatus, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStatus(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterStatus expands an instance of ClusterStatus into a JSON
// request object.
func expandClusterStatus(c *Client, f *ClusterStatus, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenClusterStatus flattens an instance of ClusterStatus from a JSON
// response object.
func flattenClusterStatus(c *Client, i interface{}, res *Cluster) *ClusterStatus {
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
func expandClusterStatusHistoryMap(c *Client, f map[string]ClusterStatusHistory, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterStatusHistory(c, &item, res)
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
func expandClusterStatusHistorySlice(c *Client, f []ClusterStatusHistory, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterStatusHistory(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterStatusHistoryMap flattens the contents of ClusterStatusHistory from a JSON
// response object.
func flattenClusterStatusHistoryMap(c *Client, i interface{}, res *Cluster) map[string]ClusterStatusHistory {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStatusHistory{}
	}

	if len(a) == 0 {
		return map[string]ClusterStatusHistory{}
	}

	items := make(map[string]ClusterStatusHistory)
	for k, item := range a {
		items[k] = *flattenClusterStatusHistory(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterStatusHistorySlice flattens the contents of ClusterStatusHistory from a JSON
// response object.
func flattenClusterStatusHistorySlice(c *Client, i interface{}, res *Cluster) []ClusterStatusHistory {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStatusHistory{}
	}

	if len(a) == 0 {
		return []ClusterStatusHistory{}
	}

	items := make([]ClusterStatusHistory, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStatusHistory(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterStatusHistory expands an instance of ClusterStatusHistory into a JSON
// request object.
func expandClusterStatusHistory(c *Client, f *ClusterStatusHistory, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenClusterStatusHistory flattens an instance of ClusterStatusHistory from a JSON
// response object.
func flattenClusterStatusHistory(c *Client, i interface{}, res *Cluster) *ClusterStatusHistory {
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
func expandClusterMetricsMap(c *Client, f map[string]ClusterMetrics, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterMetrics(c, &item, res)
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
func expandClusterMetricsSlice(c *Client, f []ClusterMetrics, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterMetrics(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterMetricsMap flattens the contents of ClusterMetrics from a JSON
// response object.
func flattenClusterMetricsMap(c *Client, i interface{}, res *Cluster) map[string]ClusterMetrics {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterMetrics{}
	}

	if len(a) == 0 {
		return map[string]ClusterMetrics{}
	}

	items := make(map[string]ClusterMetrics)
	for k, item := range a {
		items[k] = *flattenClusterMetrics(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterMetricsSlice flattens the contents of ClusterMetrics from a JSON
// response object.
func flattenClusterMetricsSlice(c *Client, i interface{}, res *Cluster) []ClusterMetrics {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterMetrics{}
	}

	if len(a) == 0 {
		return []ClusterMetrics{}
	}

	items := make([]ClusterMetrics, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterMetrics(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterMetrics expands an instance of ClusterMetrics into a JSON
// request object.
func expandClusterMetrics(c *Client, f *ClusterMetrics, res *Cluster) (map[string]interface{}, error) {
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
func flattenClusterMetrics(c *Client, i interface{}, res *Cluster) *ClusterMetrics {
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

// expandClusterVirtualClusterConfigMap expands the contents of ClusterVirtualClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigMap(c *Client, f map[string]ClusterVirtualClusterConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigSlice expands the contents of ClusterVirtualClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigSlice(c *Client, f []ClusterVirtualClusterConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigMap flattens the contents of ClusterVirtualClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigSlice flattens the contents of ClusterVirtualClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfig{}
	}

	items := make([]ClusterVirtualClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfig expands an instance of ClusterVirtualClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfig(c *Client, f *ClusterVirtualClusterConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.StagingBucket; !dcl.IsEmptyValueIndirect(v) {
		m["stagingBucket"] = v
	}
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfig(c, f.KubernetesClusterConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding KubernetesClusterConfig into kubernetesClusterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["kubernetesClusterConfig"] = v
	}
	if v, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfig(c, f.AuxiliaryServicesConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding AuxiliaryServicesConfig into auxiliaryServicesConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["auxiliaryServicesConfig"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfig flattens an instance of ClusterVirtualClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfig
	}
	r.StagingBucket = dcl.FlattenString(m["stagingBucket"])
	r.KubernetesClusterConfig = flattenClusterVirtualClusterConfigKubernetesClusterConfig(c, m["kubernetesClusterConfig"], res)
	r.AuxiliaryServicesConfig = flattenClusterVirtualClusterConfigAuxiliaryServicesConfig(c, m["auxiliaryServicesConfig"], res)

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfig{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfig expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfig(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KubernetesNamespace; !dcl.IsEmptyValueIndirect(v) {
		m["kubernetesNamespace"] = v
	}
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, f.GkeClusterConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GkeClusterConfig into gkeClusterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gkeClusterConfig"] = v
	}
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, f.KubernetesSoftwareConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding KubernetesSoftwareConfig into kubernetesSoftwareConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["kubernetesSoftwareConfig"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfig flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfig
	}
	r.KubernetesNamespace = dcl.FlattenString(m["kubernetesNamespace"])
	r.GkeClusterConfig = flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, m["gkeClusterConfig"], res)
	r.KubernetesSoftwareConfig = flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, m["kubernetesSoftwareConfig"], res)

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GkeClusterTarget; !dcl.IsEmptyValueIndirect(v) {
		m["gkeClusterTarget"] = v
	}
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice(c, f.NodePoolTarget, res); err != nil {
		return nil, fmt.Errorf("error expanding NodePoolTarget into nodePoolTarget: %w", err)
	} else if v != nil {
		m["nodePoolTarget"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig
	}
	r.GkeClusterTarget = dcl.FlattenString(m["gkeClusterTarget"])
	r.NodePoolTarget = flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice(c, m["nodePoolTarget"], res)

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NodePool; !dcl.IsEmptyValueIndirect(v) {
		m["nodePool"] = v
	}
	if v := f.Roles; v != nil {
		m["roles"] = v
	}
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c, f.NodePoolConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding NodePoolConfig into nodePoolConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["nodePoolConfig"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget
	}
	r.NodePool = dcl.FlattenString(m["nodePool"])
	r.Roles = flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnumSlice(c, m["roles"], res)
	r.NodePoolConfig = flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c, m["nodePoolConfig"], res)

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, f.Config, res); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["config"] = v
	}
	if v := f.Locations; v != nil {
		m["locations"] = v
	}
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, f.Autoscaling, res); err != nil {
		return nil, fmt.Errorf("error expanding Autoscaling into autoscaling: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["autoscaling"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig
	}
	r.Config = flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, m["config"], res)
	r.Locations = dcl.FlattenStringSlice(m["locations"])
	r.Autoscaling = flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, m["autoscaling"], res)

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MachineType; !dcl.IsEmptyValueIndirect(v) {
		m["machineType"] = v
	}
	if v := f.LocalSsdCount; !dcl.IsEmptyValueIndirect(v) {
		m["localSsdCount"] = v
	}
	if v := f.Preemptible; !dcl.IsEmptyValueIndirect(v) {
		m["preemptible"] = v
	}
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice(c, f.Accelerators, res); err != nil {
		return nil, fmt.Errorf("error expanding Accelerators into accelerators: %w", err)
	} else if v != nil {
		m["accelerators"] = v
	}
	if v := f.MinCpuPlatform; !dcl.IsEmptyValueIndirect(v) {
		m["minCpuPlatform"] = v
	}
	if v := f.BootDiskKmsKey; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskKmsKey"] = v
	}
	if v, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, f.EphemeralStorageConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding EphemeralStorageConfig into ephemeralStorageConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["ephemeralStorageConfig"] = v
	}
	if v := f.Spot; !dcl.IsEmptyValueIndirect(v) {
		m["spot"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig
	}
	r.MachineType = dcl.FlattenString(m["machineType"])
	r.LocalSsdCount = dcl.FlattenInteger(m["localSsdCount"])
	r.Preemptible = dcl.FlattenBool(m["preemptible"])
	r.Accelerators = flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice(c, m["accelerators"], res)
	r.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])
	r.BootDiskKmsKey = dcl.FlattenString(m["bootDiskKmsKey"])
	r.EphemeralStorageConfig = flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, m["ephemeralStorageConfig"], res)
	r.Spot = dcl.FlattenBool(m["spot"])

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AcceleratorCount; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorCount"] = v
	}
	if v := f.AcceleratorType; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorType"] = v
	}
	if v := f.GpuPartitionSize; !dcl.IsEmptyValueIndirect(v) {
		m["gpuPartitionSize"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators
	}
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])
	r.AcceleratorType = dcl.FlattenString(m["acceleratorType"])
	r.GpuPartitionSize = dcl.FlattenString(m["gpuPartitionSize"])

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.LocalSsdCount; !dcl.IsEmptyValueIndirect(v) {
		m["localSsdCount"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig
	}
	r.LocalSsdCount = dcl.FlattenInteger(m["localSsdCount"])

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling, res *Cluster) (map[string]interface{}, error) {
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

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling
	}
	r.MinNodeCount = dcl.FlattenInteger(m["minNodeCount"])
	r.MaxNodeCount = dcl.FlattenInteger(m["maxNodeCount"])

	return r
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigMap expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigSlice expands the contents of ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigSlice(c *Client, f []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig expands an instance of ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c *Client, f *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ComponentVersion; !dcl.IsEmptyValueIndirect(v) {
		m["componentVersion"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig flattens an instance of ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig
	}
	r.ComponentVersion = dcl.FlattenKeyValuePairs(m["componentVersion"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])

	return r
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfigMap expands the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfigSlice expands the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfigSlice(c *Client, f []ClusterVirtualClusterConfigAuxiliaryServicesConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMap flattens the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigAuxiliaryServicesConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSlice flattens the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigAuxiliaryServicesConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigAuxiliaryServicesConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigAuxiliaryServicesConfig{}
	}

	items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigAuxiliaryServicesConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfig expands an instance of ClusterVirtualClusterConfigAuxiliaryServicesConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfig(c *Client, f *ClusterVirtualClusterConfigAuxiliaryServicesConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, f.MetastoreConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding MetastoreConfig into metastoreConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metastoreConfig"] = v
	}
	if v, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, f.SparkHistoryServerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SparkHistoryServerConfig into sparkHistoryServerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sparkHistoryServerConfig"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfig flattens an instance of ClusterVirtualClusterConfigAuxiliaryServicesConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigAuxiliaryServicesConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigAuxiliaryServicesConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigAuxiliaryServicesConfig
	}
	r.MetastoreConfig = flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, m["metastoreConfig"], res)
	r.SparkHistoryServerConfig = flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, m["sparkHistoryServerConfig"], res)

	return r
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigMap expands the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigSlice expands the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigSlice(c *Client, f []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigMap flattens the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigSlice flattens the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{}
	}

	items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig expands an instance of ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c *Client, f *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DataprocMetastoreService; !dcl.IsEmptyValueIndirect(v) {
		m["dataprocMetastoreService"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig flattens an instance of ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig
	}
	r.DataprocMetastoreService = dcl.FlattenString(m["dataprocMetastoreService"])

	return r
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigMap expands the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigMap(c *Client, f map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigSlice expands the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigSlice(c *Client, f []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigMap flattens the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{}
	}

	items := make(map[string]ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigSlice flattens the contents of ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{}
	}

	items := make([]ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig expands an instance of ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig into a JSON
// request object.
func expandClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c *Client, f *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DataprocCluster; !dcl.IsEmptyValueIndirect(v) {
		m["dataprocCluster"] = v
	}

	return m, nil
}

// flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig flattens an instance of ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig from a JSON
// response object.
func flattenClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig(c *Client, i interface{}, res *Cluster) *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig
	}
	r.DataprocCluster = dcl.FlattenString(m["dataprocCluster"])

	return r
}

// flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumMap flattens the contents of ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum from a JSON
// response object.
func flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	items := make(map[string]ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum)
	for k, item := range a {
		items[k] = *flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(item.(interface{}))
	}

	return items
}

// flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumSlice flattens the contents of ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum from a JSON
// response object.
func flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	if len(a) == 0 {
		return []ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	items := make([]ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(item.(interface{})))
	}

	return items
}

// flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum with the same value as that string.
func flattenClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(i interface{}) *ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef(s)
}

// flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumMap flattens the contents of ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum from a JSON
// response object.
func flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	items := make(map[string]ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum)
	for k, item := range a {
		items[k] = *flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(item.(interface{}))
	}

	return items
}

// flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumSlice flattens the contents of ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum from a JSON
// response object.
func flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	if len(a) == 0 {
		return []ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	items := make([]ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(item.(interface{})))
	}

	return items
}

// flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum with the same value as that string.
func flattenClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(i interface{}) *ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef(s)
}

// flattenClusterConfigMasterConfigPreemptibilityEnumMap flattens the contents of ClusterConfigMasterConfigPreemptibilityEnum from a JSON
// response object.
func flattenClusterConfigMasterConfigPreemptibilityEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigMasterConfigPreemptibilityEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigMasterConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigMasterConfigPreemptibilityEnum{}
	}

	items := make(map[string]ClusterConfigMasterConfigPreemptibilityEnum)
	for k, item := range a {
		items[k] = *flattenClusterConfigMasterConfigPreemptibilityEnum(item.(interface{}))
	}

	return items
}

// flattenClusterConfigMasterConfigPreemptibilityEnumSlice flattens the contents of ClusterConfigMasterConfigPreemptibilityEnum from a JSON
// response object.
func flattenClusterConfigMasterConfigPreemptibilityEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigMasterConfigPreemptibilityEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigMasterConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return []ClusterConfigMasterConfigPreemptibilityEnum{}
	}

	items := make([]ClusterConfigMasterConfigPreemptibilityEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigMasterConfigPreemptibilityEnum(item.(interface{})))
	}

	return items
}

// flattenClusterConfigMasterConfigPreemptibilityEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterConfigMasterConfigPreemptibilityEnum with the same value as that string.
func flattenClusterConfigMasterConfigPreemptibilityEnum(i interface{}) *ClusterConfigMasterConfigPreemptibilityEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterConfigMasterConfigPreemptibilityEnumRef(s)
}

// flattenClusterConfigWorkerConfigPreemptibilityEnumMap flattens the contents of ClusterConfigWorkerConfigPreemptibilityEnum from a JSON
// response object.
func flattenClusterConfigWorkerConfigPreemptibilityEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigWorkerConfigPreemptibilityEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigWorkerConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigWorkerConfigPreemptibilityEnum{}
	}

	items := make(map[string]ClusterConfigWorkerConfigPreemptibilityEnum)
	for k, item := range a {
		items[k] = *flattenClusterConfigWorkerConfigPreemptibilityEnum(item.(interface{}))
	}

	return items
}

// flattenClusterConfigWorkerConfigPreemptibilityEnumSlice flattens the contents of ClusterConfigWorkerConfigPreemptibilityEnum from a JSON
// response object.
func flattenClusterConfigWorkerConfigPreemptibilityEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigWorkerConfigPreemptibilityEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigWorkerConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return []ClusterConfigWorkerConfigPreemptibilityEnum{}
	}

	items := make([]ClusterConfigWorkerConfigPreemptibilityEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigWorkerConfigPreemptibilityEnum(item.(interface{})))
	}

	return items
}

// flattenClusterConfigWorkerConfigPreemptibilityEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterConfigWorkerConfigPreemptibilityEnum with the same value as that string.
func flattenClusterConfigWorkerConfigPreemptibilityEnum(i interface{}) *ClusterConfigWorkerConfigPreemptibilityEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterConfigWorkerConfigPreemptibilityEnumRef(s)
}

// flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnumMap flattens the contents of ClusterConfigSecondaryWorkerConfigPreemptibilityEnum from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSecondaryWorkerConfigPreemptibilityEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSecondaryWorkerConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSecondaryWorkerConfigPreemptibilityEnum{}
	}

	items := make(map[string]ClusterConfigSecondaryWorkerConfigPreemptibilityEnum)
	for k, item := range a {
		items[k] = *flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnum(item.(interface{}))
	}

	return items
}

// flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnumSlice flattens the contents of ClusterConfigSecondaryWorkerConfigPreemptibilityEnum from a JSON
// response object.
func flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSecondaryWorkerConfigPreemptibilityEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSecondaryWorkerConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return []ClusterConfigSecondaryWorkerConfigPreemptibilityEnum{}
	}

	items := make([]ClusterConfigSecondaryWorkerConfigPreemptibilityEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnum(item.(interface{})))
	}

	return items
}

// flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterConfigSecondaryWorkerConfigPreemptibilityEnum with the same value as that string.
func flattenClusterConfigSecondaryWorkerConfigPreemptibilityEnum(i interface{}) *ClusterConfigSecondaryWorkerConfigPreemptibilityEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterConfigSecondaryWorkerConfigPreemptibilityEnumRef(s)
}

// flattenClusterConfigSoftwareConfigOptionalComponentsEnumMap flattens the contents of ClusterConfigSoftwareConfigOptionalComponentsEnum from a JSON
// response object.
func flattenClusterConfigSoftwareConfigOptionalComponentsEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigSoftwareConfigOptionalComponentsEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	items := make(map[string]ClusterConfigSoftwareConfigOptionalComponentsEnum)
	for k, item := range a {
		items[k] = *flattenClusterConfigSoftwareConfigOptionalComponentsEnum(item.(interface{}))
	}

	return items
}

// flattenClusterConfigSoftwareConfigOptionalComponentsEnumSlice flattens the contents of ClusterConfigSoftwareConfigOptionalComponentsEnum from a JSON
// response object.
func flattenClusterConfigSoftwareConfigOptionalComponentsEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigSoftwareConfigOptionalComponentsEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	if len(a) == 0 {
		return []ClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	items := make([]ClusterConfigSoftwareConfigOptionalComponentsEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigSoftwareConfigOptionalComponentsEnum(item.(interface{})))
	}

	return items
}

// flattenClusterConfigSoftwareConfigOptionalComponentsEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterConfigSoftwareConfigOptionalComponentsEnum with the same value as that string.
func flattenClusterConfigSoftwareConfigOptionalComponentsEnum(i interface{}) *ClusterConfigSoftwareConfigOptionalComponentsEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterConfigSoftwareConfigOptionalComponentsEnumRef(s)
}

// flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnumMap flattens the contents of ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum from a JSON
// response object.
func flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum{}
	}

	items := make(map[string]ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum)
	for k, item := range a {
		items[k] = *flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnum(item.(interface{}))
	}

	return items
}

// flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnumSlice flattens the contents of ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum from a JSON
// response object.
func flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum{}
	}

	if len(a) == 0 {
		return []ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum{}
	}

	items := make([]ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnum(item.(interface{})))
	}

	return items
}

// flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum with the same value as that string.
func flattenClusterConfigDataprocMetricConfigMetricsMetricSourceEnum(i interface{}) *ClusterConfigDataprocMetricConfigMetricsMetricSourceEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterConfigDataprocMetricConfigMetricsMetricSourceEnumRef(s)
}

// flattenClusterStatusStateEnumMap flattens the contents of ClusterStatusStateEnum from a JSON
// response object.
func flattenClusterStatusStateEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterStatusStateEnum {
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
func flattenClusterStatusStateEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterStatusStateEnum {
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
		return nil
	}

	return ClusterStatusStateEnumRef(s)
}

// flattenClusterStatusSubstateEnumMap flattens the contents of ClusterStatusSubstateEnum from a JSON
// response object.
func flattenClusterStatusSubstateEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterStatusSubstateEnum {
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
func flattenClusterStatusSubstateEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterStatusSubstateEnum {
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
		return nil
	}

	return ClusterStatusSubstateEnumRef(s)
}

// flattenClusterStatusHistoryStateEnumMap flattens the contents of ClusterStatusHistoryStateEnum from a JSON
// response object.
func flattenClusterStatusHistoryStateEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterStatusHistoryStateEnum {
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
func flattenClusterStatusHistoryStateEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterStatusHistoryStateEnum {
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
		return nil
	}

	return ClusterStatusHistoryStateEnumRef(s)
}

// flattenClusterStatusHistorySubstateEnumMap flattens the contents of ClusterStatusHistorySubstateEnum from a JSON
// response object.
func flattenClusterStatusHistorySubstateEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterStatusHistorySubstateEnum {
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
func flattenClusterStatusHistorySubstateEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterStatusHistorySubstateEnum {
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
		return nil
	}

	return ClusterStatusHistorySubstateEnumRef(s)
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnumMap flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum{}
	}

	items := make(map[string]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum)
	for k, item := range a {
		items[k] = *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum(item.(interface{}))
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnumSlice flattens the contents of ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum from a JSON
// response object.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum{}
	}

	if len(a) == 0 {
		return []ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum{}
	}

	items := make([]ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum(item.(interface{})))
	}

	return items
}

// flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum with the same value as that string.
func flattenClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum(i interface{}) *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetRolesEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Cluster) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalCluster(b, c, r)
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
	FieldName        string // used for error logging
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
		// Use the first field diff's field name for logging required recreate error.
		diff := clusterDiff{FieldName: fieldDiffs[0].FieldName}
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
	vConfig := r.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &ClusterConfig{}
	}
	if err := extractClusterConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		r.Config = vConfig
	}
	vStatus := r.Status
	if vStatus == nil {
		// note: explicitly not the empty object.
		vStatus = &ClusterStatus{}
	}
	if err := extractClusterStatusFields(r, vStatus); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vStatus) {
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
	if !dcl.IsEmptyValueIndirect(vMetrics) {
		r.Metrics = vMetrics
	}
	vVirtualClusterConfig := r.VirtualClusterConfig
	if vVirtualClusterConfig == nil {
		// note: explicitly not the empty object.
		vVirtualClusterConfig = &ClusterVirtualClusterConfig{}
	}
	if err := extractClusterVirtualClusterConfigFields(r, vVirtualClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vVirtualClusterConfig) {
		r.VirtualClusterConfig = vVirtualClusterConfig
	}
	return nil
}
func extractClusterConfigFields(r *Cluster, o *ClusterConfig) error {
	vGceClusterConfig := o.GceClusterConfig
	if vGceClusterConfig == nil {
		// note: explicitly not the empty object.
		vGceClusterConfig = &ClusterConfigGceClusterConfig{}
	}
	if err := extractClusterConfigGceClusterConfigFields(r, vGceClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGceClusterConfig) {
		o.GceClusterConfig = vGceClusterConfig
	}
	vMasterConfig := o.MasterConfig
	if vMasterConfig == nil {
		// note: explicitly not the empty object.
		vMasterConfig = &ClusterConfigMasterConfig{}
	}
	if err := extractClusterConfigMasterConfigFields(r, vMasterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMasterConfig) {
		o.MasterConfig = vMasterConfig
	}
	vWorkerConfig := o.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &ClusterConfigWorkerConfig{}
	}
	if err := extractClusterConfigWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		o.WorkerConfig = vWorkerConfig
	}
	vSecondaryWorkerConfig := o.SecondaryWorkerConfig
	if vSecondaryWorkerConfig == nil {
		// note: explicitly not the empty object.
		vSecondaryWorkerConfig = &ClusterConfigSecondaryWorkerConfig{}
	}
	if err := extractClusterConfigSecondaryWorkerConfigFields(r, vSecondaryWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecondaryWorkerConfig) {
		o.SecondaryWorkerConfig = vSecondaryWorkerConfig
	}
	vSoftwareConfig := o.SoftwareConfig
	if vSoftwareConfig == nil {
		// note: explicitly not the empty object.
		vSoftwareConfig = &ClusterConfigSoftwareConfig{}
	}
	if err := extractClusterConfigSoftwareConfigFields(r, vSoftwareConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSoftwareConfig) {
		o.SoftwareConfig = vSoftwareConfig
	}
	vEncryptionConfig := o.EncryptionConfig
	if vEncryptionConfig == nil {
		// note: explicitly not the empty object.
		vEncryptionConfig = &ClusterConfigEncryptionConfig{}
	}
	if err := extractClusterConfigEncryptionConfigFields(r, vEncryptionConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEncryptionConfig) {
		o.EncryptionConfig = vEncryptionConfig
	}
	vAutoscalingConfig := o.AutoscalingConfig
	if vAutoscalingConfig == nil {
		// note: explicitly not the empty object.
		vAutoscalingConfig = &ClusterConfigAutoscalingConfig{}
	}
	if err := extractClusterConfigAutoscalingConfigFields(r, vAutoscalingConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAutoscalingConfig) {
		o.AutoscalingConfig = vAutoscalingConfig
	}
	vSecurityConfig := o.SecurityConfig
	if vSecurityConfig == nil {
		// note: explicitly not the empty object.
		vSecurityConfig = &ClusterConfigSecurityConfig{}
	}
	if err := extractClusterConfigSecurityConfigFields(r, vSecurityConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecurityConfig) {
		o.SecurityConfig = vSecurityConfig
	}
	vLifecycleConfig := o.LifecycleConfig
	if vLifecycleConfig == nil {
		// note: explicitly not the empty object.
		vLifecycleConfig = &ClusterConfigLifecycleConfig{}
	}
	if err := extractClusterConfigLifecycleConfigFields(r, vLifecycleConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLifecycleConfig) {
		o.LifecycleConfig = vLifecycleConfig
	}
	vEndpointConfig := o.EndpointConfig
	if vEndpointConfig == nil {
		// note: explicitly not the empty object.
		vEndpointConfig = &ClusterConfigEndpointConfig{}
	}
	if err := extractClusterConfigEndpointConfigFields(r, vEndpointConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEndpointConfig) {
		o.EndpointConfig = vEndpointConfig
	}
	vMetastoreConfig := o.MetastoreConfig
	if vMetastoreConfig == nil {
		// note: explicitly not the empty object.
		vMetastoreConfig = &ClusterConfigMetastoreConfig{}
	}
	if err := extractClusterConfigMetastoreConfigFields(r, vMetastoreConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetastoreConfig) {
		o.MetastoreConfig = vMetastoreConfig
	}
	vDataprocMetricConfig := o.DataprocMetricConfig
	if vDataprocMetricConfig == nil {
		// note: explicitly not the empty object.
		vDataprocMetricConfig = &ClusterConfigDataprocMetricConfig{}
	}
	if err := extractClusterConfigDataprocMetricConfigFields(r, vDataprocMetricConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDataprocMetricConfig) {
		o.DataprocMetricConfig = vDataprocMetricConfig
	}
	return nil
}
func extractClusterConfigGceClusterConfigFields(r *Cluster, o *ClusterConfigGceClusterConfig) error {
	vReservationAffinity := o.ReservationAffinity
	if vReservationAffinity == nil {
		// note: explicitly not the empty object.
		vReservationAffinity = &ClusterConfigGceClusterConfigReservationAffinity{}
	}
	if err := extractClusterConfigGceClusterConfigReservationAffinityFields(r, vReservationAffinity); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vReservationAffinity) {
		o.ReservationAffinity = vReservationAffinity
	}
	vNodeGroupAffinity := o.NodeGroupAffinity
	if vNodeGroupAffinity == nil {
		// note: explicitly not the empty object.
		vNodeGroupAffinity = &ClusterConfigGceClusterConfigNodeGroupAffinity{}
	}
	if err := extractClusterConfigGceClusterConfigNodeGroupAffinityFields(r, vNodeGroupAffinity); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNodeGroupAffinity) {
		o.NodeGroupAffinity = vNodeGroupAffinity
	}
	vShieldedInstanceConfig := o.ShieldedInstanceConfig
	if vShieldedInstanceConfig == nil {
		// note: explicitly not the empty object.
		vShieldedInstanceConfig = &ClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}
	if err := extractClusterConfigGceClusterConfigShieldedInstanceConfigFields(r, vShieldedInstanceConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vShieldedInstanceConfig) {
		o.ShieldedInstanceConfig = vShieldedInstanceConfig
	}
	vConfidentialInstanceConfig := o.ConfidentialInstanceConfig
	if vConfidentialInstanceConfig == nil {
		// note: explicitly not the empty object.
		vConfidentialInstanceConfig = &ClusterConfigGceClusterConfigConfidentialInstanceConfig{}
	}
	if err := extractClusterConfigGceClusterConfigConfidentialInstanceConfigFields(r, vConfidentialInstanceConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfidentialInstanceConfig) {
		o.ConfidentialInstanceConfig = vConfidentialInstanceConfig
	}
	return nil
}
func extractClusterConfigGceClusterConfigReservationAffinityFields(r *Cluster, o *ClusterConfigGceClusterConfigReservationAffinity) error {
	return nil
}
func extractClusterConfigGceClusterConfigNodeGroupAffinityFields(r *Cluster, o *ClusterConfigGceClusterConfigNodeGroupAffinity) error {
	return nil
}
func extractClusterConfigGceClusterConfigShieldedInstanceConfigFields(r *Cluster, o *ClusterConfigGceClusterConfigShieldedInstanceConfig) error {
	return nil
}
func extractClusterConfigGceClusterConfigConfidentialInstanceConfigFields(r *Cluster, o *ClusterConfigGceClusterConfigConfidentialInstanceConfig) error {
	return nil
}
func extractClusterConfigMasterConfigFields(r *Cluster, o *ClusterConfigMasterConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &ClusterConfigMasterConfigDiskConfig{}
	}
	if err := extractClusterConfigMasterConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &ClusterConfigMasterConfigManagedGroupConfig{}
	}
	if err := extractClusterConfigMasterConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func extractClusterConfigMasterConfigDiskConfigFields(r *Cluster, o *ClusterConfigMasterConfigDiskConfig) error {
	return nil
}
func extractClusterConfigMasterConfigManagedGroupConfigFields(r *Cluster, o *ClusterConfigMasterConfigManagedGroupConfig) error {
	return nil
}
func extractClusterConfigMasterConfigAcceleratorsFields(r *Cluster, o *ClusterConfigMasterConfigAccelerators) error {
	return nil
}
func extractClusterConfigMasterConfigInstanceReferencesFields(r *Cluster, o *ClusterConfigMasterConfigInstanceReferences) error {
	return nil
}
func extractClusterConfigWorkerConfigFields(r *Cluster, o *ClusterConfigWorkerConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &ClusterConfigWorkerConfigDiskConfig{}
	}
	if err := extractClusterConfigWorkerConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &ClusterConfigWorkerConfigManagedGroupConfig{}
	}
	if err := extractClusterConfigWorkerConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func extractClusterConfigWorkerConfigDiskConfigFields(r *Cluster, o *ClusterConfigWorkerConfigDiskConfig) error {
	return nil
}
func extractClusterConfigWorkerConfigManagedGroupConfigFields(r *Cluster, o *ClusterConfigWorkerConfigManagedGroupConfig) error {
	return nil
}
func extractClusterConfigWorkerConfigAcceleratorsFields(r *Cluster, o *ClusterConfigWorkerConfigAccelerators) error {
	return nil
}
func extractClusterConfigWorkerConfigInstanceReferencesFields(r *Cluster, o *ClusterConfigWorkerConfigInstanceReferences) error {
	return nil
}
func extractClusterConfigSecondaryWorkerConfigFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &ClusterConfigSecondaryWorkerConfigDiskConfig{}
	}
	if err := extractClusterConfigSecondaryWorkerConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &ClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}
	if err := extractClusterConfigSecondaryWorkerConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func extractClusterConfigSecondaryWorkerConfigDiskConfigFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfigDiskConfig) error {
	return nil
}
func extractClusterConfigSecondaryWorkerConfigManagedGroupConfigFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfigManagedGroupConfig) error {
	return nil
}
func extractClusterConfigSecondaryWorkerConfigAcceleratorsFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfigAccelerators) error {
	return nil
}
func extractClusterConfigSecondaryWorkerConfigInstanceReferencesFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfigInstanceReferences) error {
	return nil
}
func extractClusterConfigSoftwareConfigFields(r *Cluster, o *ClusterConfigSoftwareConfig) error {
	return nil
}
func extractClusterConfigInitializationActionsFields(r *Cluster, o *ClusterConfigInitializationActions) error {
	return nil
}
func extractClusterConfigEncryptionConfigFields(r *Cluster, o *ClusterConfigEncryptionConfig) error {
	return nil
}
func extractClusterConfigAutoscalingConfigFields(r *Cluster, o *ClusterConfigAutoscalingConfig) error {
	return nil
}
func extractClusterConfigSecurityConfigFields(r *Cluster, o *ClusterConfigSecurityConfig) error {
	vKerberosConfig := o.KerberosConfig
	if vKerberosConfig == nil {
		// note: explicitly not the empty object.
		vKerberosConfig = &ClusterConfigSecurityConfigKerberosConfig{}
	}
	if err := extractClusterConfigSecurityConfigKerberosConfigFields(r, vKerberosConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKerberosConfig) {
		o.KerberosConfig = vKerberosConfig
	}
	vIdentityConfig := o.IdentityConfig
	if vIdentityConfig == nil {
		// note: explicitly not the empty object.
		vIdentityConfig = &ClusterConfigSecurityConfigIdentityConfig{}
	}
	if err := extractClusterConfigSecurityConfigIdentityConfigFields(r, vIdentityConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIdentityConfig) {
		o.IdentityConfig = vIdentityConfig
	}
	return nil
}
func extractClusterConfigSecurityConfigKerberosConfigFields(r *Cluster, o *ClusterConfigSecurityConfigKerberosConfig) error {
	return nil
}
func extractClusterConfigSecurityConfigIdentityConfigFields(r *Cluster, o *ClusterConfigSecurityConfigIdentityConfig) error {
	return nil
}
func extractClusterConfigLifecycleConfigFields(r *Cluster, o *ClusterConfigLifecycleConfig) error {
	return nil
}
func extractClusterConfigEndpointConfigFields(r *Cluster, o *ClusterConfigEndpointConfig) error {
	return nil
}
func extractClusterConfigMetastoreConfigFields(r *Cluster, o *ClusterConfigMetastoreConfig) error {
	return nil
}
func extractClusterConfigDataprocMetricConfigFields(r *Cluster, o *ClusterConfigDataprocMetricConfig) error {
	return nil
}
func extractClusterConfigDataprocMetricConfigMetricsFields(r *Cluster, o *ClusterConfigDataprocMetricConfigMetrics) error {
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
func extractClusterVirtualClusterConfigFields(r *Cluster, o *ClusterVirtualClusterConfig) error {
	vKubernetesClusterConfig := o.KubernetesClusterConfig
	if vKubernetesClusterConfig == nil {
		// note: explicitly not the empty object.
		vKubernetesClusterConfig = &ClusterVirtualClusterConfigKubernetesClusterConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigFields(r, vKubernetesClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKubernetesClusterConfig) {
		o.KubernetesClusterConfig = vKubernetesClusterConfig
	}
	vAuxiliaryServicesConfig := o.AuxiliaryServicesConfig
	if vAuxiliaryServicesConfig == nil {
		// note: explicitly not the empty object.
		vAuxiliaryServicesConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfig{}
	}
	if err := extractClusterVirtualClusterConfigAuxiliaryServicesConfigFields(r, vAuxiliaryServicesConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuxiliaryServicesConfig) {
		o.AuxiliaryServicesConfig = vAuxiliaryServicesConfig
	}
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfig) error {
	vGkeClusterConfig := o.GkeClusterConfig
	if vGkeClusterConfig == nil {
		// note: explicitly not the empty object.
		vGkeClusterConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigFields(r, vGkeClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGkeClusterConfig) {
		o.GkeClusterConfig = vGkeClusterConfig
	}
	vKubernetesSoftwareConfig := o.KubernetesSoftwareConfig
	if vKubernetesSoftwareConfig == nil {
		// note: explicitly not the empty object.
		vKubernetesSoftwareConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigFields(r, vKubernetesSoftwareConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKubernetesSoftwareConfig) {
		o.KubernetesSoftwareConfig = vKubernetesSoftwareConfig
	}
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) error {
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) error {
	vNodePoolConfig := o.NodePoolConfig
	if vNodePoolConfig == nil {
		// note: explicitly not the empty object.
		vNodePoolConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigFields(r, vNodePoolConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNodePoolConfig) {
		o.NodePoolConfig = vNodePoolConfig
	}
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) error {
	vConfig := o.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		o.Config = vConfig
	}
	vAutoscaling := o.Autoscaling
	if vAutoscaling == nil {
		// note: explicitly not the empty object.
		vAutoscaling = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingFields(r, vAutoscaling); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAutoscaling) {
		o.Autoscaling = vAutoscaling
	}
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) error {
	vEphemeralStorageConfig := o.EphemeralStorageConfig
	if vEphemeralStorageConfig == nil {
		// note: explicitly not the empty object.
		vEphemeralStorageConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigFields(r, vEphemeralStorageConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEphemeralStorageConfig) {
		o.EphemeralStorageConfig = vEphemeralStorageConfig
	}
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) error {
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) error {
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) error {
	return nil
}
func extractClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) error {
	return nil
}
func extractClusterVirtualClusterConfigAuxiliaryServicesConfigFields(r *Cluster, o *ClusterVirtualClusterConfigAuxiliaryServicesConfig) error {
	vMetastoreConfig := o.MetastoreConfig
	if vMetastoreConfig == nil {
		// note: explicitly not the empty object.
		vMetastoreConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{}
	}
	if err := extractClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigFields(r, vMetastoreConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetastoreConfig) {
		o.MetastoreConfig = vMetastoreConfig
	}
	vSparkHistoryServerConfig := o.SparkHistoryServerConfig
	if vSparkHistoryServerConfig == nil {
		// note: explicitly not the empty object.
		vSparkHistoryServerConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{}
	}
	if err := extractClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigFields(r, vSparkHistoryServerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSparkHistoryServerConfig) {
		o.SparkHistoryServerConfig = vSparkHistoryServerConfig
	}
	return nil
}
func extractClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigFields(r *Cluster, o *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) error {
	return nil
}
func extractClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigFields(r *Cluster, o *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) error {
	return nil
}

func postReadExtractClusterFields(r *Cluster) error {
	vConfig := r.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &ClusterConfig{}
	}
	if err := postReadExtractClusterConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		r.Config = vConfig
	}
	vStatus := r.Status
	if vStatus == nil {
		// note: explicitly not the empty object.
		vStatus = &ClusterStatus{}
	}
	if err := postReadExtractClusterStatusFields(r, vStatus); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vStatus) {
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
	if !dcl.IsEmptyValueIndirect(vMetrics) {
		r.Metrics = vMetrics
	}
	vVirtualClusterConfig := r.VirtualClusterConfig
	if vVirtualClusterConfig == nil {
		// note: explicitly not the empty object.
		vVirtualClusterConfig = &ClusterVirtualClusterConfig{}
	}
	if err := postReadExtractClusterVirtualClusterConfigFields(r, vVirtualClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vVirtualClusterConfig) {
		r.VirtualClusterConfig = vVirtualClusterConfig
	}
	return nil
}
func postReadExtractClusterConfigFields(r *Cluster, o *ClusterConfig) error {
	vGceClusterConfig := o.GceClusterConfig
	if vGceClusterConfig == nil {
		// note: explicitly not the empty object.
		vGceClusterConfig = &ClusterConfigGceClusterConfig{}
	}
	if err := extractClusterConfigGceClusterConfigFields(r, vGceClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGceClusterConfig) {
		o.GceClusterConfig = vGceClusterConfig
	}
	vMasterConfig := o.MasterConfig
	if vMasterConfig == nil {
		// note: explicitly not the empty object.
		vMasterConfig = &ClusterConfigMasterConfig{}
	}
	if err := extractClusterConfigMasterConfigFields(r, vMasterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMasterConfig) {
		o.MasterConfig = vMasterConfig
	}
	vWorkerConfig := o.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &ClusterConfigWorkerConfig{}
	}
	if err := extractClusterConfigWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		o.WorkerConfig = vWorkerConfig
	}
	vSecondaryWorkerConfig := o.SecondaryWorkerConfig
	if vSecondaryWorkerConfig == nil {
		// note: explicitly not the empty object.
		vSecondaryWorkerConfig = &ClusterConfigSecondaryWorkerConfig{}
	}
	if err := extractClusterConfigSecondaryWorkerConfigFields(r, vSecondaryWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecondaryWorkerConfig) {
		o.SecondaryWorkerConfig = vSecondaryWorkerConfig
	}
	vSoftwareConfig := o.SoftwareConfig
	if vSoftwareConfig == nil {
		// note: explicitly not the empty object.
		vSoftwareConfig = &ClusterConfigSoftwareConfig{}
	}
	if err := extractClusterConfigSoftwareConfigFields(r, vSoftwareConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSoftwareConfig) {
		o.SoftwareConfig = vSoftwareConfig
	}
	vEncryptionConfig := o.EncryptionConfig
	if vEncryptionConfig == nil {
		// note: explicitly not the empty object.
		vEncryptionConfig = &ClusterConfigEncryptionConfig{}
	}
	if err := extractClusterConfigEncryptionConfigFields(r, vEncryptionConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEncryptionConfig) {
		o.EncryptionConfig = vEncryptionConfig
	}
	vAutoscalingConfig := o.AutoscalingConfig
	if vAutoscalingConfig == nil {
		// note: explicitly not the empty object.
		vAutoscalingConfig = &ClusterConfigAutoscalingConfig{}
	}
	if err := extractClusterConfigAutoscalingConfigFields(r, vAutoscalingConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAutoscalingConfig) {
		o.AutoscalingConfig = vAutoscalingConfig
	}
	vSecurityConfig := o.SecurityConfig
	if vSecurityConfig == nil {
		// note: explicitly not the empty object.
		vSecurityConfig = &ClusterConfigSecurityConfig{}
	}
	if err := extractClusterConfigSecurityConfigFields(r, vSecurityConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecurityConfig) {
		o.SecurityConfig = vSecurityConfig
	}
	vLifecycleConfig := o.LifecycleConfig
	if vLifecycleConfig == nil {
		// note: explicitly not the empty object.
		vLifecycleConfig = &ClusterConfigLifecycleConfig{}
	}
	if err := extractClusterConfigLifecycleConfigFields(r, vLifecycleConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLifecycleConfig) {
		o.LifecycleConfig = vLifecycleConfig
	}
	vEndpointConfig := o.EndpointConfig
	if vEndpointConfig == nil {
		// note: explicitly not the empty object.
		vEndpointConfig = &ClusterConfigEndpointConfig{}
	}
	if err := extractClusterConfigEndpointConfigFields(r, vEndpointConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEndpointConfig) {
		o.EndpointConfig = vEndpointConfig
	}
	vMetastoreConfig := o.MetastoreConfig
	if vMetastoreConfig == nil {
		// note: explicitly not the empty object.
		vMetastoreConfig = &ClusterConfigMetastoreConfig{}
	}
	if err := extractClusterConfigMetastoreConfigFields(r, vMetastoreConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetastoreConfig) {
		o.MetastoreConfig = vMetastoreConfig
	}
	vDataprocMetricConfig := o.DataprocMetricConfig
	if vDataprocMetricConfig == nil {
		// note: explicitly not the empty object.
		vDataprocMetricConfig = &ClusterConfigDataprocMetricConfig{}
	}
	if err := extractClusterConfigDataprocMetricConfigFields(r, vDataprocMetricConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDataprocMetricConfig) {
		o.DataprocMetricConfig = vDataprocMetricConfig
	}
	return nil
}
func postReadExtractClusterConfigGceClusterConfigFields(r *Cluster, o *ClusterConfigGceClusterConfig) error {
	vReservationAffinity := o.ReservationAffinity
	if vReservationAffinity == nil {
		// note: explicitly not the empty object.
		vReservationAffinity = &ClusterConfigGceClusterConfigReservationAffinity{}
	}
	if err := extractClusterConfigGceClusterConfigReservationAffinityFields(r, vReservationAffinity); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vReservationAffinity) {
		o.ReservationAffinity = vReservationAffinity
	}
	vNodeGroupAffinity := o.NodeGroupAffinity
	if vNodeGroupAffinity == nil {
		// note: explicitly not the empty object.
		vNodeGroupAffinity = &ClusterConfigGceClusterConfigNodeGroupAffinity{}
	}
	if err := extractClusterConfigGceClusterConfigNodeGroupAffinityFields(r, vNodeGroupAffinity); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNodeGroupAffinity) {
		o.NodeGroupAffinity = vNodeGroupAffinity
	}
	vShieldedInstanceConfig := o.ShieldedInstanceConfig
	if vShieldedInstanceConfig == nil {
		// note: explicitly not the empty object.
		vShieldedInstanceConfig = &ClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}
	if err := extractClusterConfigGceClusterConfigShieldedInstanceConfigFields(r, vShieldedInstanceConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vShieldedInstanceConfig) {
		o.ShieldedInstanceConfig = vShieldedInstanceConfig
	}
	vConfidentialInstanceConfig := o.ConfidentialInstanceConfig
	if vConfidentialInstanceConfig == nil {
		// note: explicitly not the empty object.
		vConfidentialInstanceConfig = &ClusterConfigGceClusterConfigConfidentialInstanceConfig{}
	}
	if err := extractClusterConfigGceClusterConfigConfidentialInstanceConfigFields(r, vConfidentialInstanceConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfidentialInstanceConfig) {
		o.ConfidentialInstanceConfig = vConfidentialInstanceConfig
	}
	return nil
}
func postReadExtractClusterConfigGceClusterConfigReservationAffinityFields(r *Cluster, o *ClusterConfigGceClusterConfigReservationAffinity) error {
	return nil
}
func postReadExtractClusterConfigGceClusterConfigNodeGroupAffinityFields(r *Cluster, o *ClusterConfigGceClusterConfigNodeGroupAffinity) error {
	return nil
}
func postReadExtractClusterConfigGceClusterConfigShieldedInstanceConfigFields(r *Cluster, o *ClusterConfigGceClusterConfigShieldedInstanceConfig) error {
	return nil
}
func postReadExtractClusterConfigGceClusterConfigConfidentialInstanceConfigFields(r *Cluster, o *ClusterConfigGceClusterConfigConfidentialInstanceConfig) error {
	return nil
}
func postReadExtractClusterConfigMasterConfigFields(r *Cluster, o *ClusterConfigMasterConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &ClusterConfigMasterConfigDiskConfig{}
	}
	if err := extractClusterConfigMasterConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &ClusterConfigMasterConfigManagedGroupConfig{}
	}
	if err := extractClusterConfigMasterConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func postReadExtractClusterConfigMasterConfigDiskConfigFields(r *Cluster, o *ClusterConfigMasterConfigDiskConfig) error {
	return nil
}
func postReadExtractClusterConfigMasterConfigManagedGroupConfigFields(r *Cluster, o *ClusterConfigMasterConfigManagedGroupConfig) error {
	return nil
}
func postReadExtractClusterConfigMasterConfigAcceleratorsFields(r *Cluster, o *ClusterConfigMasterConfigAccelerators) error {
	return nil
}
func postReadExtractClusterConfigMasterConfigInstanceReferencesFields(r *Cluster, o *ClusterConfigMasterConfigInstanceReferences) error {
	return nil
}
func postReadExtractClusterConfigWorkerConfigFields(r *Cluster, o *ClusterConfigWorkerConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &ClusterConfigWorkerConfigDiskConfig{}
	}
	if err := extractClusterConfigWorkerConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &ClusterConfigWorkerConfigManagedGroupConfig{}
	}
	if err := extractClusterConfigWorkerConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func postReadExtractClusterConfigWorkerConfigDiskConfigFields(r *Cluster, o *ClusterConfigWorkerConfigDiskConfig) error {
	return nil
}
func postReadExtractClusterConfigWorkerConfigManagedGroupConfigFields(r *Cluster, o *ClusterConfigWorkerConfigManagedGroupConfig) error {
	return nil
}
func postReadExtractClusterConfigWorkerConfigAcceleratorsFields(r *Cluster, o *ClusterConfigWorkerConfigAccelerators) error {
	return nil
}
func postReadExtractClusterConfigWorkerConfigInstanceReferencesFields(r *Cluster, o *ClusterConfigWorkerConfigInstanceReferences) error {
	return nil
}
func postReadExtractClusterConfigSecondaryWorkerConfigFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &ClusterConfigSecondaryWorkerConfigDiskConfig{}
	}
	if err := extractClusterConfigSecondaryWorkerConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &ClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}
	if err := extractClusterConfigSecondaryWorkerConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func postReadExtractClusterConfigSecondaryWorkerConfigDiskConfigFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfigDiskConfig) error {
	return nil
}
func postReadExtractClusterConfigSecondaryWorkerConfigManagedGroupConfigFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfigManagedGroupConfig) error {
	return nil
}
func postReadExtractClusterConfigSecondaryWorkerConfigAcceleratorsFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfigAccelerators) error {
	return nil
}
func postReadExtractClusterConfigSecondaryWorkerConfigInstanceReferencesFields(r *Cluster, o *ClusterConfigSecondaryWorkerConfigInstanceReferences) error {
	return nil
}
func postReadExtractClusterConfigSoftwareConfigFields(r *Cluster, o *ClusterConfigSoftwareConfig) error {
	return nil
}
func postReadExtractClusterConfigInitializationActionsFields(r *Cluster, o *ClusterConfigInitializationActions) error {
	return nil
}
func postReadExtractClusterConfigEncryptionConfigFields(r *Cluster, o *ClusterConfigEncryptionConfig) error {
	return nil
}
func postReadExtractClusterConfigAutoscalingConfigFields(r *Cluster, o *ClusterConfigAutoscalingConfig) error {
	return nil
}
func postReadExtractClusterConfigSecurityConfigFields(r *Cluster, o *ClusterConfigSecurityConfig) error {
	vKerberosConfig := o.KerberosConfig
	if vKerberosConfig == nil {
		// note: explicitly not the empty object.
		vKerberosConfig = &ClusterConfigSecurityConfigKerberosConfig{}
	}
	if err := extractClusterConfigSecurityConfigKerberosConfigFields(r, vKerberosConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKerberosConfig) {
		o.KerberosConfig = vKerberosConfig
	}
	vIdentityConfig := o.IdentityConfig
	if vIdentityConfig == nil {
		// note: explicitly not the empty object.
		vIdentityConfig = &ClusterConfigSecurityConfigIdentityConfig{}
	}
	if err := extractClusterConfigSecurityConfigIdentityConfigFields(r, vIdentityConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIdentityConfig) {
		o.IdentityConfig = vIdentityConfig
	}
	return nil
}
func postReadExtractClusterConfigSecurityConfigKerberosConfigFields(r *Cluster, o *ClusterConfigSecurityConfigKerberosConfig) error {
	return nil
}
func postReadExtractClusterConfigSecurityConfigIdentityConfigFields(r *Cluster, o *ClusterConfigSecurityConfigIdentityConfig) error {
	return nil
}
func postReadExtractClusterConfigLifecycleConfigFields(r *Cluster, o *ClusterConfigLifecycleConfig) error {
	return nil
}
func postReadExtractClusterConfigEndpointConfigFields(r *Cluster, o *ClusterConfigEndpointConfig) error {
	return nil
}
func postReadExtractClusterConfigMetastoreConfigFields(r *Cluster, o *ClusterConfigMetastoreConfig) error {
	return nil
}
func postReadExtractClusterConfigDataprocMetricConfigFields(r *Cluster, o *ClusterConfigDataprocMetricConfig) error {
	return nil
}
func postReadExtractClusterConfigDataprocMetricConfigMetricsFields(r *Cluster, o *ClusterConfigDataprocMetricConfigMetrics) error {
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
func postReadExtractClusterVirtualClusterConfigFields(r *Cluster, o *ClusterVirtualClusterConfig) error {
	vKubernetesClusterConfig := o.KubernetesClusterConfig
	if vKubernetesClusterConfig == nil {
		// note: explicitly not the empty object.
		vKubernetesClusterConfig = &ClusterVirtualClusterConfigKubernetesClusterConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigFields(r, vKubernetesClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKubernetesClusterConfig) {
		o.KubernetesClusterConfig = vKubernetesClusterConfig
	}
	vAuxiliaryServicesConfig := o.AuxiliaryServicesConfig
	if vAuxiliaryServicesConfig == nil {
		// note: explicitly not the empty object.
		vAuxiliaryServicesConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfig{}
	}
	if err := extractClusterVirtualClusterConfigAuxiliaryServicesConfigFields(r, vAuxiliaryServicesConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuxiliaryServicesConfig) {
		o.AuxiliaryServicesConfig = vAuxiliaryServicesConfig
	}
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfig) error {
	vGkeClusterConfig := o.GkeClusterConfig
	if vGkeClusterConfig == nil {
		// note: explicitly not the empty object.
		vGkeClusterConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigFields(r, vGkeClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGkeClusterConfig) {
		o.GkeClusterConfig = vGkeClusterConfig
	}
	vKubernetesSoftwareConfig := o.KubernetesSoftwareConfig
	if vKubernetesSoftwareConfig == nil {
		// note: explicitly not the empty object.
		vKubernetesSoftwareConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigFields(r, vKubernetesSoftwareConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKubernetesSoftwareConfig) {
		o.KubernetesSoftwareConfig = vKubernetesSoftwareConfig
	}
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfig) error {
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTarget) error {
	vNodePoolConfig := o.NodePoolConfig
	if vNodePoolConfig == nil {
		// note: explicitly not the empty object.
		vNodePoolConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigFields(r, vNodePoolConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNodePoolConfig) {
		o.NodePoolConfig = vNodePoolConfig
	}
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfig) error {
	vConfig := o.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		o.Config = vConfig
	}
	vAutoscaling := o.Autoscaling
	if vAutoscaling == nil {
		// note: explicitly not the empty object.
		vAutoscaling = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingFields(r, vAutoscaling); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAutoscaling) {
		o.Autoscaling = vAutoscaling
	}
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfig) error {
	vEphemeralStorageConfig := o.EphemeralStorageConfig
	if vEphemeralStorageConfig == nil {
		// note: explicitly not the empty object.
		vEphemeralStorageConfig = &ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig{}
	}
	if err := extractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigFields(r, vEphemeralStorageConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEphemeralStorageConfig) {
		o.EphemeralStorageConfig = vEphemeralStorageConfig
	}
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAcceleratorsFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigAccelerators) error {
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigEphemeralStorageConfig) error {
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscaling) error {
	return nil
}
func postReadExtractClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigFields(r *Cluster, o *ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfig) error {
	return nil
}
func postReadExtractClusterVirtualClusterConfigAuxiliaryServicesConfigFields(r *Cluster, o *ClusterVirtualClusterConfigAuxiliaryServicesConfig) error {
	vMetastoreConfig := o.MetastoreConfig
	if vMetastoreConfig == nil {
		// note: explicitly not the empty object.
		vMetastoreConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig{}
	}
	if err := extractClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigFields(r, vMetastoreConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMetastoreConfig) {
		o.MetastoreConfig = vMetastoreConfig
	}
	vSparkHistoryServerConfig := o.SparkHistoryServerConfig
	if vSparkHistoryServerConfig == nil {
		// note: explicitly not the empty object.
		vSparkHistoryServerConfig = &ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig{}
	}
	if err := extractClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigFields(r, vSparkHistoryServerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSparkHistoryServerConfig) {
		o.SparkHistoryServerConfig = vSparkHistoryServerConfig
	}
	return nil
}
func postReadExtractClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigFields(r *Cluster, o *ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfig) error {
	return nil
}
func postReadExtractClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigFields(r *Cluster, o *ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfig) error {
	return nil
}
