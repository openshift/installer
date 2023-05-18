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
package containeraws

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

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "networking"); err != nil {
		return err
	}
	if err := dcl.Required(r, "awsRegion"); err != nil {
		return err
	}
	if err := dcl.Required(r, "controlPlane"); err != nil {
		return err
	}
	if err := dcl.Required(r, "authorization"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if err := dcl.Required(r, "fleet"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Networking) {
		if err := r.Networking.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ControlPlane) {
		if err := r.ControlPlane.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Authorization) {
		if err := r.Authorization.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.WorkloadIdentityConfig) {
		if err := r.WorkloadIdentityConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Fleet) {
		if err := r.Fleet.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ClusterNetworking) validate() error {
	if err := dcl.Required(r, "vpcId"); err != nil {
		return err
	}
	if err := dcl.Required(r, "podAddressCidrBlocks"); err != nil {
		return err
	}
	if err := dcl.Required(r, "serviceAddressCidrBlocks"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterControlPlane) validate() error {
	if err := dcl.Required(r, "version"); err != nil {
		return err
	}
	if err := dcl.Required(r, "subnetIds"); err != nil {
		return err
	}
	if err := dcl.Required(r, "configEncryption"); err != nil {
		return err
	}
	if err := dcl.Required(r, "iamInstanceProfile"); err != nil {
		return err
	}
	if err := dcl.Required(r, "databaseEncryption"); err != nil {
		return err
	}
	if err := dcl.Required(r, "awsServicesAuthentication"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.SshConfig) {
		if err := r.SshConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ConfigEncryption) {
		if err := r.ConfigEncryption.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.RootVolume) {
		if err := r.RootVolume.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.MainVolume) {
		if err := r.MainVolume.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.DatabaseEncryption) {
		if err := r.DatabaseEncryption.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AwsServicesAuthentication) {
		if err := r.AwsServicesAuthentication.validate(); err != nil {
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
func (r *ClusterControlPlaneSshConfig) validate() error {
	if err := dcl.Required(r, "ec2KeyPair"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterControlPlaneConfigEncryption) validate() error {
	if err := dcl.Required(r, "kmsKeyArn"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterControlPlaneRootVolume) validate() error {
	return nil
}
func (r *ClusterControlPlaneMainVolume) validate() error {
	return nil
}
func (r *ClusterControlPlaneDatabaseEncryption) validate() error {
	if err := dcl.Required(r, "kmsKeyArn"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterControlPlaneAwsServicesAuthentication) validate() error {
	if err := dcl.Required(r, "roleArn"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterControlPlaneProxyConfig) validate() error {
	if err := dcl.Required(r, "secretArn"); err != nil {
		return err
	}
	if err := dcl.Required(r, "secretVersion"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterAuthorization) validate() error {
	if err := dcl.Required(r, "adminUsers"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterAuthorizationAdminUsers) validate() error {
	if err := dcl.Required(r, "username"); err != nil {
		return err
	}
	return nil
}
func (r *ClusterWorkloadIdentityConfig) validate() error {
	return nil
}
func (r *ClusterFleet) validate() error {
	if err := dcl.Required(r, "project"); err != nil {
		return err
	}
	return nil
}
func (r *Cluster) basePath() string {
	params := map[string]interface{}{
		"location": dcl.ValueOrEmptyString(r.Location),
	}
	return dcl.Nprintf("https://{{location}}-gkemulticloud.googleapis.com/v1", params)
}

func (r *Cluster) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Cluster) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters", nr.basePath(), userBasePath, params), nil

}

func (r *Cluster) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters?awsClusterId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Cluster) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters/{{name}}", nr.basePath(), userBasePath, params), nil
}

// clusterApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type clusterApiOperation interface {
	do(context.Context, *Cluster, *Client) error
}

// newUpdateClusterUpdateAwsClusterRequest creates a request for an
// Cluster resource's UpdateAwsCluster update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateClusterUpdateAwsClusterRequest(ctx context.Context, f *Cluster, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v, err := expandClusterControlPlane(c, f.ControlPlane, res); err != nil {
		return nil, fmt.Errorf("error expanding ControlPlane into controlPlane: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["controlPlane"] = v
	}
	if v, err := expandClusterAuthorization(c, f.Authorization, res); err != nil {
		return nil, fmt.Errorf("error expanding Authorization into authorization: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["authorization"] = v
	}
	if v := f.Annotations; !dcl.IsEmptyValueIndirect(v) {
		req["annotations"] = v
	}
	b, err := c.getClusterRaw(ctx, f)
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

// marshalUpdateClusterUpdateAwsClusterRequest converts the update into
// the final JSON request body.
func marshalUpdateClusterUpdateAwsClusterRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateClusterUpdateAwsClusterOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateClusterUpdateAwsClusterOperation) do(ctx context.Context, r *Cluster, c *Client) error {
	_, err := c.GetCluster(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateAwsCluster")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateClusterUpdateAwsClusterRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateClusterUpdateAwsClusterRequest(c, req)
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
	AwsClusters []map[string]interface{} `json:"awsClusters"`
	Token       string                   `json:"nextPageToken"`
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
	for _, v := range m.AwsClusters {
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
		rawDesired.Networking = canonicalizeClusterNetworking(rawDesired.Networking, nil, opts...)
		rawDesired.ControlPlane = canonicalizeClusterControlPlane(rawDesired.ControlPlane, nil, opts...)
		rawDesired.Authorization = canonicalizeClusterAuthorization(rawDesired.Authorization, nil, opts...)
		rawDesired.WorkloadIdentityConfig = canonicalizeClusterWorkloadIdentityConfig(rawDesired.WorkloadIdentityConfig, nil, opts...)
		rawDesired.Fleet = canonicalizeClusterFleet(rawDesired.Fleet, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Cluster{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	canonicalDesired.Networking = canonicalizeClusterNetworking(rawDesired.Networking, rawInitial.Networking, opts...)
	if dcl.StringCanonicalize(rawDesired.AwsRegion, rawInitial.AwsRegion) {
		canonicalDesired.AwsRegion = rawInitial.AwsRegion
	} else {
		canonicalDesired.AwsRegion = rawDesired.AwsRegion
	}
	canonicalDesired.ControlPlane = canonicalizeClusterControlPlane(rawDesired.ControlPlane, rawInitial.ControlPlane, opts...)
	canonicalDesired.Authorization = canonicalizeClusterAuthorization(rawDesired.Authorization, rawInitial.Authorization, opts...)
	if dcl.IsZeroValue(rawDesired.Annotations) || (dcl.IsEmptyValueIndirect(rawDesired.Annotations) && dcl.IsEmptyValueIndirect(rawInitial.Annotations)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Annotations = rawInitial.Annotations
	} else {
		canonicalDesired.Annotations = rawDesired.Annotations
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
	canonicalDesired.Fleet = canonicalizeClusterFleet(rawDesired.Fleet, rawInitial.Fleet, opts...)
	return canonicalDesired, nil
}

func canonicalizeClusterNewState(c *Client, rawNew, rawDesired *Cluster) (*Cluster, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Networking) && dcl.IsEmptyValueIndirect(rawDesired.Networking) {
		rawNew.Networking = rawDesired.Networking
	} else {
		rawNew.Networking = canonicalizeNewClusterNetworking(c, rawDesired.Networking, rawNew.Networking)
	}

	if dcl.IsEmptyValueIndirect(rawNew.AwsRegion) && dcl.IsEmptyValueIndirect(rawDesired.AwsRegion) {
		rawNew.AwsRegion = rawDesired.AwsRegion
	} else {
		if dcl.StringCanonicalize(rawDesired.AwsRegion, rawNew.AwsRegion) {
			rawNew.AwsRegion = rawDesired.AwsRegion
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.ControlPlane) && dcl.IsEmptyValueIndirect(rawDesired.ControlPlane) {
		rawNew.ControlPlane = rawDesired.ControlPlane
	} else {
		rawNew.ControlPlane = canonicalizeNewClusterControlPlane(c, rawDesired.ControlPlane, rawNew.ControlPlane)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Authorization) && dcl.IsEmptyValueIndirect(rawDesired.Authorization) {
		rawNew.Authorization = rawDesired.Authorization
	} else {
		rawNew.Authorization = canonicalizeNewClusterAuthorization(c, rawDesired.Authorization, rawNew.Authorization)
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Endpoint) && dcl.IsEmptyValueIndirect(rawDesired.Endpoint) {
		rawNew.Endpoint = rawDesired.Endpoint
	} else {
		if dcl.StringCanonicalize(rawDesired.Endpoint, rawNew.Endpoint) {
			rawNew.Endpoint = rawDesired.Endpoint
		}
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

	if dcl.IsEmptyValueIndirect(rawNew.WorkloadIdentityConfig) && dcl.IsEmptyValueIndirect(rawDesired.WorkloadIdentityConfig) {
		rawNew.WorkloadIdentityConfig = rawDesired.WorkloadIdentityConfig
	} else {
		rawNew.WorkloadIdentityConfig = canonicalizeNewClusterWorkloadIdentityConfig(c, rawDesired.WorkloadIdentityConfig, rawNew.WorkloadIdentityConfig)
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	if dcl.IsEmptyValueIndirect(rawNew.Fleet) && dcl.IsEmptyValueIndirect(rawDesired.Fleet) {
		rawNew.Fleet = rawDesired.Fleet
	} else {
		rawNew.Fleet = canonicalizeNewClusterFleet(c, rawDesired.Fleet, rawNew.Fleet)
	}

	return rawNew, nil
}

func canonicalizeClusterNetworking(des, initial *ClusterNetworking, opts ...dcl.ApplyOption) *ClusterNetworking {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterNetworking{}

	if dcl.StringCanonicalize(des.VPCId, initial.VPCId) || dcl.IsZeroValue(des.VPCId) {
		cDes.VPCId = initial.VPCId
	} else {
		cDes.VPCId = des.VPCId
	}
	if dcl.StringArrayCanonicalize(des.PodAddressCidrBlocks, initial.PodAddressCidrBlocks) {
		cDes.PodAddressCidrBlocks = initial.PodAddressCidrBlocks
	} else {
		cDes.PodAddressCidrBlocks = des.PodAddressCidrBlocks
	}
	if dcl.StringArrayCanonicalize(des.ServiceAddressCidrBlocks, initial.ServiceAddressCidrBlocks) {
		cDes.ServiceAddressCidrBlocks = initial.ServiceAddressCidrBlocks
	} else {
		cDes.ServiceAddressCidrBlocks = des.ServiceAddressCidrBlocks
	}

	return cDes
}

func canonicalizeClusterNetworkingSlice(des, initial []ClusterNetworking, opts ...dcl.ApplyOption) []ClusterNetworking {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterNetworking, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterNetworking(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterNetworking, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterNetworking(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterNetworking(c *Client, des, nw *ClusterNetworking) *ClusterNetworking {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterNetworking while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.VPCId, nw.VPCId) {
		nw.VPCId = des.VPCId
	}
	if dcl.StringArrayCanonicalize(des.PodAddressCidrBlocks, nw.PodAddressCidrBlocks) {
		nw.PodAddressCidrBlocks = des.PodAddressCidrBlocks
	}
	if dcl.StringArrayCanonicalize(des.ServiceAddressCidrBlocks, nw.ServiceAddressCidrBlocks) {
		nw.ServiceAddressCidrBlocks = des.ServiceAddressCidrBlocks
	}

	return nw
}

func canonicalizeNewClusterNetworkingSet(c *Client, des, nw []ClusterNetworking) []ClusterNetworking {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterNetworking
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterNetworkingNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterNetworking(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterNetworkingSlice(c *Client, des, nw []ClusterNetworking) []ClusterNetworking {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterNetworking
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterNetworking(c, &d, &n))
	}

	return items
}

func canonicalizeClusterControlPlane(des, initial *ClusterControlPlane, opts ...dcl.ApplyOption) *ClusterControlPlane {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterControlPlane{}

	if dcl.StringCanonicalize(des.Version, initial.Version) || dcl.IsZeroValue(des.Version) {
		cDes.Version = initial.Version
	} else {
		cDes.Version = des.Version
	}
	if dcl.StringCanonicalize(des.InstanceType, initial.InstanceType) || dcl.IsZeroValue(des.InstanceType) {
		cDes.InstanceType = initial.InstanceType
	} else {
		cDes.InstanceType = des.InstanceType
	}
	cDes.SshConfig = canonicalizeClusterControlPlaneSshConfig(des.SshConfig, initial.SshConfig, opts...)
	if dcl.StringArrayCanonicalize(des.SubnetIds, initial.SubnetIds) {
		cDes.SubnetIds = initial.SubnetIds
	} else {
		cDes.SubnetIds = des.SubnetIds
	}
	cDes.ConfigEncryption = canonicalizeClusterControlPlaneConfigEncryption(des.ConfigEncryption, initial.ConfigEncryption, opts...)
	if dcl.StringArrayCanonicalize(des.SecurityGroupIds, initial.SecurityGroupIds) {
		cDes.SecurityGroupIds = initial.SecurityGroupIds
	} else {
		cDes.SecurityGroupIds = des.SecurityGroupIds
	}
	if dcl.StringCanonicalize(des.IamInstanceProfile, initial.IamInstanceProfile) || dcl.IsZeroValue(des.IamInstanceProfile) {
		cDes.IamInstanceProfile = initial.IamInstanceProfile
	} else {
		cDes.IamInstanceProfile = des.IamInstanceProfile
	}
	cDes.RootVolume = canonicalizeClusterControlPlaneRootVolume(des.RootVolume, initial.RootVolume, opts...)
	cDes.MainVolume = canonicalizeClusterControlPlaneMainVolume(des.MainVolume, initial.MainVolume, opts...)
	cDes.DatabaseEncryption = canonicalizeClusterControlPlaneDatabaseEncryption(des.DatabaseEncryption, initial.DatabaseEncryption, opts...)
	if dcl.IsZeroValue(des.Tags) || (dcl.IsEmptyValueIndirect(des.Tags) && dcl.IsEmptyValueIndirect(initial.Tags)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Tags = initial.Tags
	} else {
		cDes.Tags = des.Tags
	}
	cDes.AwsServicesAuthentication = canonicalizeClusterControlPlaneAwsServicesAuthentication(des.AwsServicesAuthentication, initial.AwsServicesAuthentication, opts...)
	cDes.ProxyConfig = canonicalizeClusterControlPlaneProxyConfig(des.ProxyConfig, initial.ProxyConfig, opts...)

	return cDes
}

func canonicalizeClusterControlPlaneSlice(des, initial []ClusterControlPlane, opts ...dcl.ApplyOption) []ClusterControlPlane {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterControlPlane, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterControlPlane(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterControlPlane, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterControlPlane(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterControlPlane(c *Client, des, nw *ClusterControlPlane) *ClusterControlPlane {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterControlPlane while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Version, nw.Version) {
		nw.Version = des.Version
	}
	if dcl.StringCanonicalize(des.InstanceType, nw.InstanceType) {
		nw.InstanceType = des.InstanceType
	}
	nw.SshConfig = canonicalizeNewClusterControlPlaneSshConfig(c, des.SshConfig, nw.SshConfig)
	if dcl.StringArrayCanonicalize(des.SubnetIds, nw.SubnetIds) {
		nw.SubnetIds = des.SubnetIds
	}
	nw.ConfigEncryption = canonicalizeNewClusterControlPlaneConfigEncryption(c, des.ConfigEncryption, nw.ConfigEncryption)
	if dcl.StringArrayCanonicalize(des.SecurityGroupIds, nw.SecurityGroupIds) {
		nw.SecurityGroupIds = des.SecurityGroupIds
	}
	if dcl.StringCanonicalize(des.IamInstanceProfile, nw.IamInstanceProfile) {
		nw.IamInstanceProfile = des.IamInstanceProfile
	}
	nw.RootVolume = canonicalizeNewClusterControlPlaneRootVolume(c, des.RootVolume, nw.RootVolume)
	nw.MainVolume = canonicalizeNewClusterControlPlaneMainVolume(c, des.MainVolume, nw.MainVolume)
	nw.DatabaseEncryption = canonicalizeNewClusterControlPlaneDatabaseEncryption(c, des.DatabaseEncryption, nw.DatabaseEncryption)
	nw.AwsServicesAuthentication = canonicalizeNewClusterControlPlaneAwsServicesAuthentication(c, des.AwsServicesAuthentication, nw.AwsServicesAuthentication)
	nw.ProxyConfig = canonicalizeNewClusterControlPlaneProxyConfig(c, des.ProxyConfig, nw.ProxyConfig)

	return nw
}

func canonicalizeNewClusterControlPlaneSet(c *Client, des, nw []ClusterControlPlane) []ClusterControlPlane {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterControlPlane
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterControlPlaneNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterControlPlane(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterControlPlaneSlice(c *Client, des, nw []ClusterControlPlane) []ClusterControlPlane {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterControlPlane
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterControlPlane(c, &d, &n))
	}

	return items
}

func canonicalizeClusterControlPlaneSshConfig(des, initial *ClusterControlPlaneSshConfig, opts ...dcl.ApplyOption) *ClusterControlPlaneSshConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterControlPlaneSshConfig{}

	if dcl.StringCanonicalize(des.Ec2KeyPair, initial.Ec2KeyPair) || dcl.IsZeroValue(des.Ec2KeyPair) {
		cDes.Ec2KeyPair = initial.Ec2KeyPair
	} else {
		cDes.Ec2KeyPair = des.Ec2KeyPair
	}

	return cDes
}

func canonicalizeClusterControlPlaneSshConfigSlice(des, initial []ClusterControlPlaneSshConfig, opts ...dcl.ApplyOption) []ClusterControlPlaneSshConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterControlPlaneSshConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterControlPlaneSshConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterControlPlaneSshConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterControlPlaneSshConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterControlPlaneSshConfig(c *Client, des, nw *ClusterControlPlaneSshConfig) *ClusterControlPlaneSshConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterControlPlaneSshConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Ec2KeyPair, nw.Ec2KeyPair) {
		nw.Ec2KeyPair = des.Ec2KeyPair
	}

	return nw
}

func canonicalizeNewClusterControlPlaneSshConfigSet(c *Client, des, nw []ClusterControlPlaneSshConfig) []ClusterControlPlaneSshConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterControlPlaneSshConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterControlPlaneSshConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterControlPlaneSshConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterControlPlaneSshConfigSlice(c *Client, des, nw []ClusterControlPlaneSshConfig) []ClusterControlPlaneSshConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterControlPlaneSshConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterControlPlaneSshConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterControlPlaneConfigEncryption(des, initial *ClusterControlPlaneConfigEncryption, opts ...dcl.ApplyOption) *ClusterControlPlaneConfigEncryption {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterControlPlaneConfigEncryption{}

	if dcl.StringCanonicalize(des.KmsKeyArn, initial.KmsKeyArn) || dcl.IsZeroValue(des.KmsKeyArn) {
		cDes.KmsKeyArn = initial.KmsKeyArn
	} else {
		cDes.KmsKeyArn = des.KmsKeyArn
	}

	return cDes
}

func canonicalizeClusterControlPlaneConfigEncryptionSlice(des, initial []ClusterControlPlaneConfigEncryption, opts ...dcl.ApplyOption) []ClusterControlPlaneConfigEncryption {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterControlPlaneConfigEncryption, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterControlPlaneConfigEncryption(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterControlPlaneConfigEncryption, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterControlPlaneConfigEncryption(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterControlPlaneConfigEncryption(c *Client, des, nw *ClusterControlPlaneConfigEncryption) *ClusterControlPlaneConfigEncryption {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterControlPlaneConfigEncryption while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KmsKeyArn, nw.KmsKeyArn) {
		nw.KmsKeyArn = des.KmsKeyArn
	}

	return nw
}

func canonicalizeNewClusterControlPlaneConfigEncryptionSet(c *Client, des, nw []ClusterControlPlaneConfigEncryption) []ClusterControlPlaneConfigEncryption {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterControlPlaneConfigEncryption
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterControlPlaneConfigEncryptionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterControlPlaneConfigEncryption(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterControlPlaneConfigEncryptionSlice(c *Client, des, nw []ClusterControlPlaneConfigEncryption) []ClusterControlPlaneConfigEncryption {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterControlPlaneConfigEncryption
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterControlPlaneConfigEncryption(c, &d, &n))
	}

	return items
}

func canonicalizeClusterControlPlaneRootVolume(des, initial *ClusterControlPlaneRootVolume, opts ...dcl.ApplyOption) *ClusterControlPlaneRootVolume {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterControlPlaneRootVolume{}

	if dcl.IsZeroValue(des.SizeGib) || (dcl.IsEmptyValueIndirect(des.SizeGib) && dcl.IsEmptyValueIndirect(initial.SizeGib)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.SizeGib = initial.SizeGib
	} else {
		cDes.SizeGib = des.SizeGib
	}
	if dcl.IsZeroValue(des.VolumeType) || (dcl.IsEmptyValueIndirect(des.VolumeType) && dcl.IsEmptyValueIndirect(initial.VolumeType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.VolumeType = initial.VolumeType
	} else {
		cDes.VolumeType = des.VolumeType
	}
	if dcl.IsZeroValue(des.Iops) || (dcl.IsEmptyValueIndirect(des.Iops) && dcl.IsEmptyValueIndirect(initial.Iops)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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

func canonicalizeClusterControlPlaneRootVolumeSlice(des, initial []ClusterControlPlaneRootVolume, opts ...dcl.ApplyOption) []ClusterControlPlaneRootVolume {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterControlPlaneRootVolume, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterControlPlaneRootVolume(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterControlPlaneRootVolume, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterControlPlaneRootVolume(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterControlPlaneRootVolume(c *Client, des, nw *ClusterControlPlaneRootVolume) *ClusterControlPlaneRootVolume {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterControlPlaneRootVolume while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KmsKeyArn, nw.KmsKeyArn) {
		nw.KmsKeyArn = des.KmsKeyArn
	}

	return nw
}

func canonicalizeNewClusterControlPlaneRootVolumeSet(c *Client, des, nw []ClusterControlPlaneRootVolume) []ClusterControlPlaneRootVolume {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterControlPlaneRootVolume
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterControlPlaneRootVolumeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterControlPlaneRootVolume(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterControlPlaneRootVolumeSlice(c *Client, des, nw []ClusterControlPlaneRootVolume) []ClusterControlPlaneRootVolume {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterControlPlaneRootVolume
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterControlPlaneRootVolume(c, &d, &n))
	}

	return items
}

func canonicalizeClusterControlPlaneMainVolume(des, initial *ClusterControlPlaneMainVolume, opts ...dcl.ApplyOption) *ClusterControlPlaneMainVolume {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterControlPlaneMainVolume{}

	if dcl.IsZeroValue(des.SizeGib) || (dcl.IsEmptyValueIndirect(des.SizeGib) && dcl.IsEmptyValueIndirect(initial.SizeGib)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.SizeGib = initial.SizeGib
	} else {
		cDes.SizeGib = des.SizeGib
	}
	if dcl.IsZeroValue(des.VolumeType) || (dcl.IsEmptyValueIndirect(des.VolumeType) && dcl.IsEmptyValueIndirect(initial.VolumeType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.VolumeType = initial.VolumeType
	} else {
		cDes.VolumeType = des.VolumeType
	}
	if dcl.IsZeroValue(des.Iops) || (dcl.IsEmptyValueIndirect(des.Iops) && dcl.IsEmptyValueIndirect(initial.Iops)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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

func canonicalizeClusterControlPlaneMainVolumeSlice(des, initial []ClusterControlPlaneMainVolume, opts ...dcl.ApplyOption) []ClusterControlPlaneMainVolume {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterControlPlaneMainVolume, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterControlPlaneMainVolume(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterControlPlaneMainVolume, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterControlPlaneMainVolume(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterControlPlaneMainVolume(c *Client, des, nw *ClusterControlPlaneMainVolume) *ClusterControlPlaneMainVolume {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterControlPlaneMainVolume while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KmsKeyArn, nw.KmsKeyArn) {
		nw.KmsKeyArn = des.KmsKeyArn
	}

	return nw
}

func canonicalizeNewClusterControlPlaneMainVolumeSet(c *Client, des, nw []ClusterControlPlaneMainVolume) []ClusterControlPlaneMainVolume {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterControlPlaneMainVolume
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterControlPlaneMainVolumeNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterControlPlaneMainVolume(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterControlPlaneMainVolumeSlice(c *Client, des, nw []ClusterControlPlaneMainVolume) []ClusterControlPlaneMainVolume {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterControlPlaneMainVolume
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterControlPlaneMainVolume(c, &d, &n))
	}

	return items
}

func canonicalizeClusterControlPlaneDatabaseEncryption(des, initial *ClusterControlPlaneDatabaseEncryption, opts ...dcl.ApplyOption) *ClusterControlPlaneDatabaseEncryption {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterControlPlaneDatabaseEncryption{}

	if dcl.StringCanonicalize(des.KmsKeyArn, initial.KmsKeyArn) || dcl.IsZeroValue(des.KmsKeyArn) {
		cDes.KmsKeyArn = initial.KmsKeyArn
	} else {
		cDes.KmsKeyArn = des.KmsKeyArn
	}

	return cDes
}

func canonicalizeClusterControlPlaneDatabaseEncryptionSlice(des, initial []ClusterControlPlaneDatabaseEncryption, opts ...dcl.ApplyOption) []ClusterControlPlaneDatabaseEncryption {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterControlPlaneDatabaseEncryption, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterControlPlaneDatabaseEncryption(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterControlPlaneDatabaseEncryption, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterControlPlaneDatabaseEncryption(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterControlPlaneDatabaseEncryption(c *Client, des, nw *ClusterControlPlaneDatabaseEncryption) *ClusterControlPlaneDatabaseEncryption {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterControlPlaneDatabaseEncryption while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KmsKeyArn, nw.KmsKeyArn) {
		nw.KmsKeyArn = des.KmsKeyArn
	}

	return nw
}

func canonicalizeNewClusterControlPlaneDatabaseEncryptionSet(c *Client, des, nw []ClusterControlPlaneDatabaseEncryption) []ClusterControlPlaneDatabaseEncryption {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterControlPlaneDatabaseEncryption
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterControlPlaneDatabaseEncryptionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterControlPlaneDatabaseEncryption(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterControlPlaneDatabaseEncryptionSlice(c *Client, des, nw []ClusterControlPlaneDatabaseEncryption) []ClusterControlPlaneDatabaseEncryption {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterControlPlaneDatabaseEncryption
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterControlPlaneDatabaseEncryption(c, &d, &n))
	}

	return items
}

func canonicalizeClusterControlPlaneAwsServicesAuthentication(des, initial *ClusterControlPlaneAwsServicesAuthentication, opts ...dcl.ApplyOption) *ClusterControlPlaneAwsServicesAuthentication {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterControlPlaneAwsServicesAuthentication{}

	if dcl.StringCanonicalize(des.RoleArn, initial.RoleArn) || dcl.IsZeroValue(des.RoleArn) {
		cDes.RoleArn = initial.RoleArn
	} else {
		cDes.RoleArn = des.RoleArn
	}
	if dcl.StringCanonicalize(des.RoleSessionName, initial.RoleSessionName) || dcl.IsZeroValue(des.RoleSessionName) {
		cDes.RoleSessionName = initial.RoleSessionName
	} else {
		cDes.RoleSessionName = des.RoleSessionName
	}

	return cDes
}

func canonicalizeClusterControlPlaneAwsServicesAuthenticationSlice(des, initial []ClusterControlPlaneAwsServicesAuthentication, opts ...dcl.ApplyOption) []ClusterControlPlaneAwsServicesAuthentication {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterControlPlaneAwsServicesAuthentication, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterControlPlaneAwsServicesAuthentication(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterControlPlaneAwsServicesAuthentication, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterControlPlaneAwsServicesAuthentication(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterControlPlaneAwsServicesAuthentication(c *Client, des, nw *ClusterControlPlaneAwsServicesAuthentication) *ClusterControlPlaneAwsServicesAuthentication {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterControlPlaneAwsServicesAuthentication while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.RoleArn, nw.RoleArn) {
		nw.RoleArn = des.RoleArn
	}
	if dcl.StringCanonicalize(des.RoleSessionName, nw.RoleSessionName) {
		nw.RoleSessionName = des.RoleSessionName
	}

	return nw
}

func canonicalizeNewClusterControlPlaneAwsServicesAuthenticationSet(c *Client, des, nw []ClusterControlPlaneAwsServicesAuthentication) []ClusterControlPlaneAwsServicesAuthentication {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterControlPlaneAwsServicesAuthentication
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterControlPlaneAwsServicesAuthenticationNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterControlPlaneAwsServicesAuthentication(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterControlPlaneAwsServicesAuthenticationSlice(c *Client, des, nw []ClusterControlPlaneAwsServicesAuthentication) []ClusterControlPlaneAwsServicesAuthentication {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterControlPlaneAwsServicesAuthentication
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterControlPlaneAwsServicesAuthentication(c, &d, &n))
	}

	return items
}

func canonicalizeClusterControlPlaneProxyConfig(des, initial *ClusterControlPlaneProxyConfig, opts ...dcl.ApplyOption) *ClusterControlPlaneProxyConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterControlPlaneProxyConfig{}

	if dcl.StringCanonicalize(des.SecretArn, initial.SecretArn) || dcl.IsZeroValue(des.SecretArn) {
		cDes.SecretArn = initial.SecretArn
	} else {
		cDes.SecretArn = des.SecretArn
	}
	if dcl.StringCanonicalize(des.SecretVersion, initial.SecretVersion) || dcl.IsZeroValue(des.SecretVersion) {
		cDes.SecretVersion = initial.SecretVersion
	} else {
		cDes.SecretVersion = des.SecretVersion
	}

	return cDes
}

func canonicalizeClusterControlPlaneProxyConfigSlice(des, initial []ClusterControlPlaneProxyConfig, opts ...dcl.ApplyOption) []ClusterControlPlaneProxyConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterControlPlaneProxyConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterControlPlaneProxyConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterControlPlaneProxyConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterControlPlaneProxyConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterControlPlaneProxyConfig(c *Client, des, nw *ClusterControlPlaneProxyConfig) *ClusterControlPlaneProxyConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterControlPlaneProxyConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.SecretArn, nw.SecretArn) {
		nw.SecretArn = des.SecretArn
	}
	if dcl.StringCanonicalize(des.SecretVersion, nw.SecretVersion) {
		nw.SecretVersion = des.SecretVersion
	}

	return nw
}

func canonicalizeNewClusterControlPlaneProxyConfigSet(c *Client, des, nw []ClusterControlPlaneProxyConfig) []ClusterControlPlaneProxyConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterControlPlaneProxyConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterControlPlaneProxyConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterControlPlaneProxyConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterControlPlaneProxyConfigSlice(c *Client, des, nw []ClusterControlPlaneProxyConfig) []ClusterControlPlaneProxyConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterControlPlaneProxyConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterControlPlaneProxyConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterAuthorization(des, initial *ClusterAuthorization, opts ...dcl.ApplyOption) *ClusterAuthorization {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterAuthorization{}

	cDes.AdminUsers = canonicalizeClusterAuthorizationAdminUsersSlice(des.AdminUsers, initial.AdminUsers, opts...)

	return cDes
}

func canonicalizeClusterAuthorizationSlice(des, initial []ClusterAuthorization, opts ...dcl.ApplyOption) []ClusterAuthorization {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterAuthorization, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterAuthorization(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterAuthorization, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterAuthorization(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterAuthorization(c *Client, des, nw *ClusterAuthorization) *ClusterAuthorization {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterAuthorization while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.AdminUsers = canonicalizeNewClusterAuthorizationAdminUsersSlice(c, des.AdminUsers, nw.AdminUsers)

	return nw
}

func canonicalizeNewClusterAuthorizationSet(c *Client, des, nw []ClusterAuthorization) []ClusterAuthorization {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterAuthorization
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterAuthorizationNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterAuthorization(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterAuthorizationSlice(c *Client, des, nw []ClusterAuthorization) []ClusterAuthorization {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterAuthorization
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterAuthorization(c, &d, &n))
	}

	return items
}

func canonicalizeClusterAuthorizationAdminUsers(des, initial *ClusterAuthorizationAdminUsers, opts ...dcl.ApplyOption) *ClusterAuthorizationAdminUsers {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterAuthorizationAdminUsers{}

	if dcl.StringCanonicalize(des.Username, initial.Username) || dcl.IsZeroValue(des.Username) {
		cDes.Username = initial.Username
	} else {
		cDes.Username = des.Username
	}

	return cDes
}

func canonicalizeClusterAuthorizationAdminUsersSlice(des, initial []ClusterAuthorizationAdminUsers, opts ...dcl.ApplyOption) []ClusterAuthorizationAdminUsers {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterAuthorizationAdminUsers, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterAuthorizationAdminUsers(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterAuthorizationAdminUsers, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterAuthorizationAdminUsers(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterAuthorizationAdminUsers(c *Client, des, nw *ClusterAuthorizationAdminUsers) *ClusterAuthorizationAdminUsers {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterAuthorizationAdminUsers while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Username, nw.Username) {
		nw.Username = des.Username
	}

	return nw
}

func canonicalizeNewClusterAuthorizationAdminUsersSet(c *Client, des, nw []ClusterAuthorizationAdminUsers) []ClusterAuthorizationAdminUsers {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterAuthorizationAdminUsers
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterAuthorizationAdminUsersNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterAuthorizationAdminUsers(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterAuthorizationAdminUsersSlice(c *Client, des, nw []ClusterAuthorizationAdminUsers) []ClusterAuthorizationAdminUsers {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterAuthorizationAdminUsers
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterAuthorizationAdminUsers(c, &d, &n))
	}

	return items
}

func canonicalizeClusterWorkloadIdentityConfig(des, initial *ClusterWorkloadIdentityConfig, opts ...dcl.ApplyOption) *ClusterWorkloadIdentityConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterWorkloadIdentityConfig{}

	if dcl.StringCanonicalize(des.IssuerUri, initial.IssuerUri) || dcl.IsZeroValue(des.IssuerUri) {
		cDes.IssuerUri = initial.IssuerUri
	} else {
		cDes.IssuerUri = des.IssuerUri
	}
	if dcl.StringCanonicalize(des.WorkloadPool, initial.WorkloadPool) || dcl.IsZeroValue(des.WorkloadPool) {
		cDes.WorkloadPool = initial.WorkloadPool
	} else {
		cDes.WorkloadPool = des.WorkloadPool
	}
	if dcl.StringCanonicalize(des.IdentityProvider, initial.IdentityProvider) || dcl.IsZeroValue(des.IdentityProvider) {
		cDes.IdentityProvider = initial.IdentityProvider
	} else {
		cDes.IdentityProvider = des.IdentityProvider
	}

	return cDes
}

func canonicalizeClusterWorkloadIdentityConfigSlice(des, initial []ClusterWorkloadIdentityConfig, opts ...dcl.ApplyOption) []ClusterWorkloadIdentityConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterWorkloadIdentityConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterWorkloadIdentityConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterWorkloadIdentityConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterWorkloadIdentityConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterWorkloadIdentityConfig(c *Client, des, nw *ClusterWorkloadIdentityConfig) *ClusterWorkloadIdentityConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterWorkloadIdentityConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.IssuerUri, nw.IssuerUri) {
		nw.IssuerUri = des.IssuerUri
	}
	if dcl.StringCanonicalize(des.WorkloadPool, nw.WorkloadPool) {
		nw.WorkloadPool = des.WorkloadPool
	}
	if dcl.StringCanonicalize(des.IdentityProvider, nw.IdentityProvider) {
		nw.IdentityProvider = des.IdentityProvider
	}

	return nw
}

func canonicalizeNewClusterWorkloadIdentityConfigSet(c *Client, des, nw []ClusterWorkloadIdentityConfig) []ClusterWorkloadIdentityConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterWorkloadIdentityConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterWorkloadIdentityConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterWorkloadIdentityConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterWorkloadIdentityConfigSlice(c *Client, des, nw []ClusterWorkloadIdentityConfig) []ClusterWorkloadIdentityConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterWorkloadIdentityConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterWorkloadIdentityConfig(c, &d, &n))
	}

	return items
}

func canonicalizeClusterFleet(des, initial *ClusterFleet, opts ...dcl.ApplyOption) *ClusterFleet {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ClusterFleet{}

	if dcl.PartialSelfLinkToSelfLink(des.Project, initial.Project) || dcl.IsZeroValue(des.Project) {
		cDes.Project = initial.Project
	} else {
		cDes.Project = des.Project
	}

	return cDes
}

func canonicalizeClusterFleetSlice(des, initial []ClusterFleet, opts ...dcl.ApplyOption) []ClusterFleet {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ClusterFleet, 0, len(des))
		for _, d := range des {
			cd := canonicalizeClusterFleet(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ClusterFleet, 0, len(des))
	for i, d := range des {
		cd := canonicalizeClusterFleet(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewClusterFleet(c *Client, des, nw *ClusterFleet) *ClusterFleet {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ClusterFleet while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.PartialSelfLinkToSelfLink(des.Project, nw.Project) {
		nw.Project = des.Project
	}
	if dcl.StringCanonicalize(des.Membership, nw.Membership) {
		nw.Membership = des.Membership
	}

	return nw
}

func canonicalizeNewClusterFleetSet(c *Client, des, nw []ClusterFleet) []ClusterFleet {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ClusterFleet
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareClusterFleetNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewClusterFleet(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewClusterFleetSlice(c *Client, des, nw []ClusterFleet) []ClusterFleet {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ClusterFleet
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewClusterFleet(c, &d, &n))
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
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Networking, actual.Networking, dcl.DiffInfo{ObjectFunction: compareClusterNetworkingNewStyle, EmptyObject: EmptyClusterNetworking, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Networking")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AwsRegion, actual.AwsRegion, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AwsRegion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ControlPlane, actual.ControlPlane, dcl.DiffInfo{ObjectFunction: compareClusterControlPlaneNewStyle, EmptyObject: EmptyClusterControlPlane, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ControlPlane")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Authorization, actual.Authorization, dcl.DiffInfo{ObjectFunction: compareClusterAuthorizationNewStyle, EmptyObject: EmptyClusterAuthorization, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Authorization")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Endpoint, actual.Endpoint, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Endpoint")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Annotations, actual.Annotations, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("Annotations")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WorkloadIdentityConfig, actual.WorkloadIdentityConfig, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareClusterWorkloadIdentityConfigNewStyle, EmptyObject: EmptyClusterWorkloadIdentityConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("WorkloadIdentityConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Fleet, actual.Fleet, dcl.DiffInfo{ObjectFunction: compareClusterFleetNewStyle, EmptyObject: EmptyClusterFleet, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Fleet")); len(ds) != 0 || err != nil {
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
func compareClusterNetworkingNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterNetworking)
	if !ok {
		desiredNotPointer, ok := d.(ClusterNetworking)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterNetworking or *ClusterNetworking", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterNetworking)
	if !ok {
		actualNotPointer, ok := a.(ClusterNetworking)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterNetworking", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.VPCId, actual.VPCId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("VpcId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PodAddressCidrBlocks, actual.PodAddressCidrBlocks, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PodAddressCidrBlocks")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAddressCidrBlocks, actual.ServiceAddressCidrBlocks, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAddressCidrBlocks")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterControlPlaneNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterControlPlane)
	if !ok {
		desiredNotPointer, ok := d.(ClusterControlPlane)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlane or *ClusterControlPlane", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterControlPlane)
	if !ok {
		actualNotPointer, ok := a.(ClusterControlPlane)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlane", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Version, actual.Version, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("Version")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceType, actual.InstanceType, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("InstanceType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SshConfig, actual.SshConfig, dcl.DiffInfo{ObjectFunction: compareClusterControlPlaneSshConfigNewStyle, EmptyObject: EmptyClusterControlPlaneSshConfig, OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("SshConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubnetIds, actual.SubnetIds, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubnetIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ConfigEncryption, actual.ConfigEncryption, dcl.DiffInfo{ObjectFunction: compareClusterControlPlaneConfigEncryptionNewStyle, EmptyObject: EmptyClusterControlPlaneConfigEncryption, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConfigEncryption")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecurityGroupIds, actual.SecurityGroupIds, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("SecurityGroupIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IamInstanceProfile, actual.IamInstanceProfile, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("IamInstanceProfile")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RootVolume, actual.RootVolume, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterControlPlaneRootVolumeNewStyle, EmptyObject: EmptyClusterControlPlaneRootVolume, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RootVolume")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MainVolume, actual.MainVolume, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareClusterControlPlaneMainVolumeNewStyle, EmptyObject: EmptyClusterControlPlaneMainVolume, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainVolume")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DatabaseEncryption, actual.DatabaseEncryption, dcl.DiffInfo{ObjectFunction: compareClusterControlPlaneDatabaseEncryptionNewStyle, EmptyObject: EmptyClusterControlPlaneDatabaseEncryption, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DatabaseEncryption")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tags, actual.Tags, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AwsServicesAuthentication, actual.AwsServicesAuthentication, dcl.DiffInfo{ObjectFunction: compareClusterControlPlaneAwsServicesAuthenticationNewStyle, EmptyObject: EmptyClusterControlPlaneAwsServicesAuthentication, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AwsServicesAuthentication")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ProxyConfig, actual.ProxyConfig, dcl.DiffInfo{ObjectFunction: compareClusterControlPlaneProxyConfigNewStyle, EmptyObject: EmptyClusterControlPlaneProxyConfig, OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("ProxyConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterControlPlaneSshConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterControlPlaneSshConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterControlPlaneSshConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneSshConfig or *ClusterControlPlaneSshConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterControlPlaneSshConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterControlPlaneSshConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneSshConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Ec2KeyPair, actual.Ec2KeyPair, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("Ec2KeyPair")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterControlPlaneConfigEncryptionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterControlPlaneConfigEncryption)
	if !ok {
		desiredNotPointer, ok := d.(ClusterControlPlaneConfigEncryption)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneConfigEncryption or *ClusterControlPlaneConfigEncryption", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterControlPlaneConfigEncryption)
	if !ok {
		actualNotPointer, ok := a.(ClusterControlPlaneConfigEncryption)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneConfigEncryption", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KmsKeyArn, actual.KmsKeyArn, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("KmsKeyArn")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterControlPlaneRootVolumeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterControlPlaneRootVolume)
	if !ok {
		desiredNotPointer, ok := d.(ClusterControlPlaneRootVolume)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneRootVolume or *ClusterControlPlaneRootVolume", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterControlPlaneRootVolume)
	if !ok {
		actualNotPointer, ok := a.(ClusterControlPlaneRootVolume)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneRootVolume", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SizeGib, actual.SizeGib, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("SizeGib")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.VolumeType, actual.VolumeType, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("VolumeType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Iops, actual.Iops, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("Iops")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KmsKeyArn, actual.KmsKeyArn, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("KmsKeyArn")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterControlPlaneMainVolumeNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterControlPlaneMainVolume)
	if !ok {
		desiredNotPointer, ok := d.(ClusterControlPlaneMainVolume)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneMainVolume or *ClusterControlPlaneMainVolume", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterControlPlaneMainVolume)
	if !ok {
		actualNotPointer, ok := a.(ClusterControlPlaneMainVolume)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneMainVolume", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SizeGib, actual.SizeGib, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SizeGib")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.VolumeType, actual.VolumeType, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("VolumeType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Iops, actual.Iops, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Iops")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KmsKeyArn, actual.KmsKeyArn, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KmsKeyArn")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterControlPlaneDatabaseEncryptionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterControlPlaneDatabaseEncryption)
	if !ok {
		desiredNotPointer, ok := d.(ClusterControlPlaneDatabaseEncryption)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneDatabaseEncryption or *ClusterControlPlaneDatabaseEncryption", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterControlPlaneDatabaseEncryption)
	if !ok {
		actualNotPointer, ok := a.(ClusterControlPlaneDatabaseEncryption)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneDatabaseEncryption", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KmsKeyArn, actual.KmsKeyArn, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KmsKeyArn")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterControlPlaneAwsServicesAuthenticationNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterControlPlaneAwsServicesAuthentication)
	if !ok {
		desiredNotPointer, ok := d.(ClusterControlPlaneAwsServicesAuthentication)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneAwsServicesAuthentication or *ClusterControlPlaneAwsServicesAuthentication", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterControlPlaneAwsServicesAuthentication)
	if !ok {
		actualNotPointer, ok := a.(ClusterControlPlaneAwsServicesAuthentication)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneAwsServicesAuthentication", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.RoleArn, actual.RoleArn, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("RoleArn")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RoleSessionName, actual.RoleSessionName, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("RoleSessionName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterControlPlaneProxyConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterControlPlaneProxyConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterControlPlaneProxyConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneProxyConfig or *ClusterControlPlaneProxyConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterControlPlaneProxyConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterControlPlaneProxyConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterControlPlaneProxyConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SecretArn, actual.SecretArn, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("SecretArn")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecretVersion, actual.SecretVersion, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("SecretVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterAuthorizationNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterAuthorization)
	if !ok {
		desiredNotPointer, ok := d.(ClusterAuthorization)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterAuthorization or *ClusterAuthorization", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterAuthorization)
	if !ok {
		actualNotPointer, ok := a.(ClusterAuthorization)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterAuthorization", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AdminUsers, actual.AdminUsers, dcl.DiffInfo{ObjectFunction: compareClusterAuthorizationAdminUsersNewStyle, EmptyObject: EmptyClusterAuthorizationAdminUsers, OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("AdminUsers")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterAuthorizationAdminUsersNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterAuthorizationAdminUsers)
	if !ok {
		desiredNotPointer, ok := d.(ClusterAuthorizationAdminUsers)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterAuthorizationAdminUsers or *ClusterAuthorizationAdminUsers", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterAuthorizationAdminUsers)
	if !ok {
		actualNotPointer, ok := a.(ClusterAuthorizationAdminUsers)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterAuthorizationAdminUsers", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Username, actual.Username, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateClusterUpdateAwsClusterOperation")}, fn.AddNest("Username")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterWorkloadIdentityConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterWorkloadIdentityConfig)
	if !ok {
		desiredNotPointer, ok := d.(ClusterWorkloadIdentityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterWorkloadIdentityConfig or *ClusterWorkloadIdentityConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterWorkloadIdentityConfig)
	if !ok {
		actualNotPointer, ok := a.(ClusterWorkloadIdentityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterWorkloadIdentityConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IssuerUri, actual.IssuerUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IssuerUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WorkloadPool, actual.WorkloadPool, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("WorkloadPool")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IdentityProvider, actual.IdentityProvider, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IdentityProvider")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareClusterFleetNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ClusterFleet)
	if !ok {
		desiredNotPointer, ok := d.(ClusterFleet)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterFleet or *ClusterFleet", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ClusterFleet)
	if !ok {
		actualNotPointer, ok := a.(ClusterFleet)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ClusterFleet", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Project, actual.Project, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Project")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Membership, actual.Membership, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Membership")); len(ds) != 0 || err != nil {
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
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.AwsRegion = dcl.SelfLinkToName(r.AwsRegion)
	normalized.Endpoint = dcl.SelfLinkToName(r.Endpoint)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.Etag = dcl.SelfLinkToName(r.Etag)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *Cluster) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateAwsCluster" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/awsClusters/{{name}}", nr.basePath(), userBasePath, fields), nil

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
	if v, err := dcl.DeriveField("projects/%s/locations/%s/awsClusters/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := expandClusterNetworking(c, f.Networking, res); err != nil {
		return nil, fmt.Errorf("error expanding Networking into networking: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["networking"] = v
	}
	if v := f.AwsRegion; dcl.ValueShouldBeSent(v) {
		m["awsRegion"] = v
	}
	if v, err := expandClusterControlPlane(c, f.ControlPlane, res); err != nil {
		return nil, fmt.Errorf("error expanding ControlPlane into controlPlane: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["controlPlane"] = v
	}
	if v, err := expandClusterAuthorization(c, f.Authorization, res); err != nil {
		return nil, fmt.Errorf("error expanding Authorization into authorization: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["authorization"] = v
	}
	if v := f.Annotations; dcl.ValueShouldBeSent(v) {
		m["annotations"] = v
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
	if v, err := expandClusterFleet(c, f.Fleet, res); err != nil {
		return nil, fmt.Errorf("error expanding Fleet into fleet: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["fleet"] = v
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
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Networking = flattenClusterNetworking(c, m["networking"], res)
	resultRes.AwsRegion = dcl.FlattenString(m["awsRegion"])
	resultRes.ControlPlane = flattenClusterControlPlane(c, m["controlPlane"], res)
	resultRes.Authorization = flattenClusterAuthorization(c, m["authorization"], res)
	resultRes.State = flattenClusterStateEnum(m["state"])
	resultRes.Endpoint = dcl.FlattenString(m["endpoint"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.Reconciling = dcl.FlattenBool(m["reconciling"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Etag = dcl.FlattenString(m["etag"])
	resultRes.Annotations = dcl.FlattenKeyValuePairs(m["annotations"])
	resultRes.WorkloadIdentityConfig = flattenClusterWorkloadIdentityConfig(c, m["workloadIdentityConfig"], res)
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.Fleet = flattenClusterFleet(c, m["fleet"], res)

	return resultRes
}

// expandClusterNetworkingMap expands the contents of ClusterNetworking into a JSON
// request object.
func expandClusterNetworkingMap(c *Client, f map[string]ClusterNetworking, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterNetworking(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterNetworkingSlice expands the contents of ClusterNetworking into a JSON
// request object.
func expandClusterNetworkingSlice(c *Client, f []ClusterNetworking, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterNetworking(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterNetworkingMap flattens the contents of ClusterNetworking from a JSON
// response object.
func flattenClusterNetworkingMap(c *Client, i interface{}, res *Cluster) map[string]ClusterNetworking {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterNetworking{}
	}

	if len(a) == 0 {
		return map[string]ClusterNetworking{}
	}

	items := make(map[string]ClusterNetworking)
	for k, item := range a {
		items[k] = *flattenClusterNetworking(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterNetworkingSlice flattens the contents of ClusterNetworking from a JSON
// response object.
func flattenClusterNetworkingSlice(c *Client, i interface{}, res *Cluster) []ClusterNetworking {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterNetworking{}
	}

	if len(a) == 0 {
		return []ClusterNetworking{}
	}

	items := make([]ClusterNetworking, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterNetworking(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterNetworking expands an instance of ClusterNetworking into a JSON
// request object.
func expandClusterNetworking(c *Client, f *ClusterNetworking, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.VPCId; !dcl.IsEmptyValueIndirect(v) {
		m["vpcId"] = v
	}
	if v := f.PodAddressCidrBlocks; v != nil {
		m["podAddressCidrBlocks"] = v
	}
	if v := f.ServiceAddressCidrBlocks; v != nil {
		m["serviceAddressCidrBlocks"] = v
	}

	return m, nil
}

// flattenClusterNetworking flattens an instance of ClusterNetworking from a JSON
// response object.
func flattenClusterNetworking(c *Client, i interface{}, res *Cluster) *ClusterNetworking {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterNetworking{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterNetworking
	}
	r.VPCId = dcl.FlattenString(m["vpcId"])
	r.PodAddressCidrBlocks = dcl.FlattenStringSlice(m["podAddressCidrBlocks"])
	r.ServiceAddressCidrBlocks = dcl.FlattenStringSlice(m["serviceAddressCidrBlocks"])

	return r
}

// expandClusterControlPlaneMap expands the contents of ClusterControlPlane into a JSON
// request object.
func expandClusterControlPlaneMap(c *Client, f map[string]ClusterControlPlane, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterControlPlane(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterControlPlaneSlice expands the contents of ClusterControlPlane into a JSON
// request object.
func expandClusterControlPlaneSlice(c *Client, f []ClusterControlPlane, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterControlPlane(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterControlPlaneMap flattens the contents of ClusterControlPlane from a JSON
// response object.
func flattenClusterControlPlaneMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlane {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlane{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlane{}
	}

	items := make(map[string]ClusterControlPlane)
	for k, item := range a {
		items[k] = *flattenClusterControlPlane(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterControlPlaneSlice flattens the contents of ClusterControlPlane from a JSON
// response object.
func flattenClusterControlPlaneSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlane {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlane{}
	}

	if len(a) == 0 {
		return []ClusterControlPlane{}
	}

	items := make([]ClusterControlPlane, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlane(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterControlPlane expands an instance of ClusterControlPlane into a JSON
// request object.
func expandClusterControlPlane(c *Client, f *ClusterControlPlane, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Version; !dcl.IsEmptyValueIndirect(v) {
		m["version"] = v
	}
	if v := f.InstanceType; !dcl.IsEmptyValueIndirect(v) {
		m["instanceType"] = v
	}
	if v, err := expandClusterControlPlaneSshConfig(c, f.SshConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SshConfig into sshConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sshConfig"] = v
	}
	if v := f.SubnetIds; v != nil {
		m["subnetIds"] = v
	}
	if v, err := expandClusterControlPlaneConfigEncryption(c, f.ConfigEncryption, res); err != nil {
		return nil, fmt.Errorf("error expanding ConfigEncryption into configEncryption: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["configEncryption"] = v
	}
	if v := f.SecurityGroupIds; v != nil {
		m["securityGroupIds"] = v
	}
	if v := f.IamInstanceProfile; !dcl.IsEmptyValueIndirect(v) {
		m["iamInstanceProfile"] = v
	}
	if v, err := expandClusterControlPlaneRootVolume(c, f.RootVolume, res); err != nil {
		return nil, fmt.Errorf("error expanding RootVolume into rootVolume: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["rootVolume"] = v
	}
	if v, err := expandClusterControlPlaneMainVolume(c, f.MainVolume, res); err != nil {
		return nil, fmt.Errorf("error expanding MainVolume into mainVolume: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["mainVolume"] = v
	}
	if v, err := expandClusterControlPlaneDatabaseEncryption(c, f.DatabaseEncryption, res); err != nil {
		return nil, fmt.Errorf("error expanding DatabaseEncryption into databaseEncryption: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["databaseEncryption"] = v
	}
	if v := f.Tags; !dcl.IsEmptyValueIndirect(v) {
		m["tags"] = v
	}
	if v, err := expandClusterControlPlaneAwsServicesAuthentication(c, f.AwsServicesAuthentication, res); err != nil {
		return nil, fmt.Errorf("error expanding AwsServicesAuthentication into awsServicesAuthentication: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["awsServicesAuthentication"] = v
	}
	if v, err := expandClusterControlPlaneProxyConfig(c, f.ProxyConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ProxyConfig into proxyConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["proxyConfig"] = v
	}

	return m, nil
}

// flattenClusterControlPlane flattens an instance of ClusterControlPlane from a JSON
// response object.
func flattenClusterControlPlane(c *Client, i interface{}, res *Cluster) *ClusterControlPlane {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterControlPlane{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterControlPlane
	}
	r.Version = dcl.FlattenString(m["version"])
	r.InstanceType = dcl.FlattenString(m["instanceType"])
	r.SshConfig = flattenClusterControlPlaneSshConfig(c, m["sshConfig"], res)
	r.SubnetIds = dcl.FlattenStringSlice(m["subnetIds"])
	r.ConfigEncryption = flattenClusterControlPlaneConfigEncryption(c, m["configEncryption"], res)
	r.SecurityGroupIds = dcl.FlattenStringSlice(m["securityGroupIds"])
	r.IamInstanceProfile = dcl.FlattenString(m["iamInstanceProfile"])
	r.RootVolume = flattenClusterControlPlaneRootVolume(c, m["rootVolume"], res)
	r.MainVolume = flattenClusterControlPlaneMainVolume(c, m["mainVolume"], res)
	r.DatabaseEncryption = flattenClusterControlPlaneDatabaseEncryption(c, m["databaseEncryption"], res)
	r.Tags = dcl.FlattenKeyValuePairs(m["tags"])
	r.AwsServicesAuthentication = flattenClusterControlPlaneAwsServicesAuthentication(c, m["awsServicesAuthentication"], res)
	r.ProxyConfig = flattenClusterControlPlaneProxyConfig(c, m["proxyConfig"], res)

	return r
}

// expandClusterControlPlaneSshConfigMap expands the contents of ClusterControlPlaneSshConfig into a JSON
// request object.
func expandClusterControlPlaneSshConfigMap(c *Client, f map[string]ClusterControlPlaneSshConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterControlPlaneSshConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterControlPlaneSshConfigSlice expands the contents of ClusterControlPlaneSshConfig into a JSON
// request object.
func expandClusterControlPlaneSshConfigSlice(c *Client, f []ClusterControlPlaneSshConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterControlPlaneSshConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterControlPlaneSshConfigMap flattens the contents of ClusterControlPlaneSshConfig from a JSON
// response object.
func flattenClusterControlPlaneSshConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneSshConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneSshConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneSshConfig{}
	}

	items := make(map[string]ClusterControlPlaneSshConfig)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneSshConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterControlPlaneSshConfigSlice flattens the contents of ClusterControlPlaneSshConfig from a JSON
// response object.
func flattenClusterControlPlaneSshConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneSshConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneSshConfig{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneSshConfig{}
	}

	items := make([]ClusterControlPlaneSshConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneSshConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterControlPlaneSshConfig expands an instance of ClusterControlPlaneSshConfig into a JSON
// request object.
func expandClusterControlPlaneSshConfig(c *Client, f *ClusterControlPlaneSshConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Ec2KeyPair; !dcl.IsEmptyValueIndirect(v) {
		m["ec2KeyPair"] = v
	}

	return m, nil
}

// flattenClusterControlPlaneSshConfig flattens an instance of ClusterControlPlaneSshConfig from a JSON
// response object.
func flattenClusterControlPlaneSshConfig(c *Client, i interface{}, res *Cluster) *ClusterControlPlaneSshConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterControlPlaneSshConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterControlPlaneSshConfig
	}
	r.Ec2KeyPair = dcl.FlattenString(m["ec2KeyPair"])

	return r
}

// expandClusterControlPlaneConfigEncryptionMap expands the contents of ClusterControlPlaneConfigEncryption into a JSON
// request object.
func expandClusterControlPlaneConfigEncryptionMap(c *Client, f map[string]ClusterControlPlaneConfigEncryption, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterControlPlaneConfigEncryption(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterControlPlaneConfigEncryptionSlice expands the contents of ClusterControlPlaneConfigEncryption into a JSON
// request object.
func expandClusterControlPlaneConfigEncryptionSlice(c *Client, f []ClusterControlPlaneConfigEncryption, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterControlPlaneConfigEncryption(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterControlPlaneConfigEncryptionMap flattens the contents of ClusterControlPlaneConfigEncryption from a JSON
// response object.
func flattenClusterControlPlaneConfigEncryptionMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneConfigEncryption {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneConfigEncryption{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneConfigEncryption{}
	}

	items := make(map[string]ClusterControlPlaneConfigEncryption)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneConfigEncryption(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterControlPlaneConfigEncryptionSlice flattens the contents of ClusterControlPlaneConfigEncryption from a JSON
// response object.
func flattenClusterControlPlaneConfigEncryptionSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneConfigEncryption {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneConfigEncryption{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneConfigEncryption{}
	}

	items := make([]ClusterControlPlaneConfigEncryption, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneConfigEncryption(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterControlPlaneConfigEncryption expands an instance of ClusterControlPlaneConfigEncryption into a JSON
// request object.
func expandClusterControlPlaneConfigEncryption(c *Client, f *ClusterControlPlaneConfigEncryption, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KmsKeyArn; !dcl.IsEmptyValueIndirect(v) {
		m["kmsKeyArn"] = v
	}

	return m, nil
}

// flattenClusterControlPlaneConfigEncryption flattens an instance of ClusterControlPlaneConfigEncryption from a JSON
// response object.
func flattenClusterControlPlaneConfigEncryption(c *Client, i interface{}, res *Cluster) *ClusterControlPlaneConfigEncryption {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterControlPlaneConfigEncryption{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterControlPlaneConfigEncryption
	}
	r.KmsKeyArn = dcl.FlattenString(m["kmsKeyArn"])

	return r
}

// expandClusterControlPlaneRootVolumeMap expands the contents of ClusterControlPlaneRootVolume into a JSON
// request object.
func expandClusterControlPlaneRootVolumeMap(c *Client, f map[string]ClusterControlPlaneRootVolume, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterControlPlaneRootVolume(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterControlPlaneRootVolumeSlice expands the contents of ClusterControlPlaneRootVolume into a JSON
// request object.
func expandClusterControlPlaneRootVolumeSlice(c *Client, f []ClusterControlPlaneRootVolume, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterControlPlaneRootVolume(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterControlPlaneRootVolumeMap flattens the contents of ClusterControlPlaneRootVolume from a JSON
// response object.
func flattenClusterControlPlaneRootVolumeMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneRootVolume {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneRootVolume{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneRootVolume{}
	}

	items := make(map[string]ClusterControlPlaneRootVolume)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneRootVolume(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterControlPlaneRootVolumeSlice flattens the contents of ClusterControlPlaneRootVolume from a JSON
// response object.
func flattenClusterControlPlaneRootVolumeSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneRootVolume {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneRootVolume{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneRootVolume{}
	}

	items := make([]ClusterControlPlaneRootVolume, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneRootVolume(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterControlPlaneRootVolume expands an instance of ClusterControlPlaneRootVolume into a JSON
// request object.
func expandClusterControlPlaneRootVolume(c *Client, f *ClusterControlPlaneRootVolume, res *Cluster) (map[string]interface{}, error) {
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

// flattenClusterControlPlaneRootVolume flattens an instance of ClusterControlPlaneRootVolume from a JSON
// response object.
func flattenClusterControlPlaneRootVolume(c *Client, i interface{}, res *Cluster) *ClusterControlPlaneRootVolume {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterControlPlaneRootVolume{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterControlPlaneRootVolume
	}
	r.SizeGib = dcl.FlattenInteger(m["sizeGib"])
	r.VolumeType = flattenClusterControlPlaneRootVolumeVolumeTypeEnum(m["volumeType"])
	r.Iops = dcl.FlattenInteger(m["iops"])
	r.KmsKeyArn = dcl.FlattenString(m["kmsKeyArn"])

	return r
}

// expandClusterControlPlaneMainVolumeMap expands the contents of ClusterControlPlaneMainVolume into a JSON
// request object.
func expandClusterControlPlaneMainVolumeMap(c *Client, f map[string]ClusterControlPlaneMainVolume, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterControlPlaneMainVolume(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterControlPlaneMainVolumeSlice expands the contents of ClusterControlPlaneMainVolume into a JSON
// request object.
func expandClusterControlPlaneMainVolumeSlice(c *Client, f []ClusterControlPlaneMainVolume, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterControlPlaneMainVolume(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterControlPlaneMainVolumeMap flattens the contents of ClusterControlPlaneMainVolume from a JSON
// response object.
func flattenClusterControlPlaneMainVolumeMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneMainVolume {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneMainVolume{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneMainVolume{}
	}

	items := make(map[string]ClusterControlPlaneMainVolume)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneMainVolume(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterControlPlaneMainVolumeSlice flattens the contents of ClusterControlPlaneMainVolume from a JSON
// response object.
func flattenClusterControlPlaneMainVolumeSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneMainVolume {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneMainVolume{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneMainVolume{}
	}

	items := make([]ClusterControlPlaneMainVolume, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneMainVolume(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterControlPlaneMainVolume expands an instance of ClusterControlPlaneMainVolume into a JSON
// request object.
func expandClusterControlPlaneMainVolume(c *Client, f *ClusterControlPlaneMainVolume, res *Cluster) (map[string]interface{}, error) {
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

// flattenClusterControlPlaneMainVolume flattens an instance of ClusterControlPlaneMainVolume from a JSON
// response object.
func flattenClusterControlPlaneMainVolume(c *Client, i interface{}, res *Cluster) *ClusterControlPlaneMainVolume {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterControlPlaneMainVolume{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterControlPlaneMainVolume
	}
	r.SizeGib = dcl.FlattenInteger(m["sizeGib"])
	r.VolumeType = flattenClusterControlPlaneMainVolumeVolumeTypeEnum(m["volumeType"])
	r.Iops = dcl.FlattenInteger(m["iops"])
	r.KmsKeyArn = dcl.FlattenString(m["kmsKeyArn"])

	return r
}

// expandClusterControlPlaneDatabaseEncryptionMap expands the contents of ClusterControlPlaneDatabaseEncryption into a JSON
// request object.
func expandClusterControlPlaneDatabaseEncryptionMap(c *Client, f map[string]ClusterControlPlaneDatabaseEncryption, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterControlPlaneDatabaseEncryption(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterControlPlaneDatabaseEncryptionSlice expands the contents of ClusterControlPlaneDatabaseEncryption into a JSON
// request object.
func expandClusterControlPlaneDatabaseEncryptionSlice(c *Client, f []ClusterControlPlaneDatabaseEncryption, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterControlPlaneDatabaseEncryption(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterControlPlaneDatabaseEncryptionMap flattens the contents of ClusterControlPlaneDatabaseEncryption from a JSON
// response object.
func flattenClusterControlPlaneDatabaseEncryptionMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneDatabaseEncryption {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneDatabaseEncryption{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneDatabaseEncryption{}
	}

	items := make(map[string]ClusterControlPlaneDatabaseEncryption)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneDatabaseEncryption(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterControlPlaneDatabaseEncryptionSlice flattens the contents of ClusterControlPlaneDatabaseEncryption from a JSON
// response object.
func flattenClusterControlPlaneDatabaseEncryptionSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneDatabaseEncryption {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneDatabaseEncryption{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneDatabaseEncryption{}
	}

	items := make([]ClusterControlPlaneDatabaseEncryption, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneDatabaseEncryption(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterControlPlaneDatabaseEncryption expands an instance of ClusterControlPlaneDatabaseEncryption into a JSON
// request object.
func expandClusterControlPlaneDatabaseEncryption(c *Client, f *ClusterControlPlaneDatabaseEncryption, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KmsKeyArn; !dcl.IsEmptyValueIndirect(v) {
		m["kmsKeyArn"] = v
	}

	return m, nil
}

// flattenClusterControlPlaneDatabaseEncryption flattens an instance of ClusterControlPlaneDatabaseEncryption from a JSON
// response object.
func flattenClusterControlPlaneDatabaseEncryption(c *Client, i interface{}, res *Cluster) *ClusterControlPlaneDatabaseEncryption {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterControlPlaneDatabaseEncryption{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterControlPlaneDatabaseEncryption
	}
	r.KmsKeyArn = dcl.FlattenString(m["kmsKeyArn"])

	return r
}

// expandClusterControlPlaneAwsServicesAuthenticationMap expands the contents of ClusterControlPlaneAwsServicesAuthentication into a JSON
// request object.
func expandClusterControlPlaneAwsServicesAuthenticationMap(c *Client, f map[string]ClusterControlPlaneAwsServicesAuthentication, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterControlPlaneAwsServicesAuthentication(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterControlPlaneAwsServicesAuthenticationSlice expands the contents of ClusterControlPlaneAwsServicesAuthentication into a JSON
// request object.
func expandClusterControlPlaneAwsServicesAuthenticationSlice(c *Client, f []ClusterControlPlaneAwsServicesAuthentication, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterControlPlaneAwsServicesAuthentication(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterControlPlaneAwsServicesAuthenticationMap flattens the contents of ClusterControlPlaneAwsServicesAuthentication from a JSON
// response object.
func flattenClusterControlPlaneAwsServicesAuthenticationMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneAwsServicesAuthentication {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneAwsServicesAuthentication{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneAwsServicesAuthentication{}
	}

	items := make(map[string]ClusterControlPlaneAwsServicesAuthentication)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneAwsServicesAuthentication(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterControlPlaneAwsServicesAuthenticationSlice flattens the contents of ClusterControlPlaneAwsServicesAuthentication from a JSON
// response object.
func flattenClusterControlPlaneAwsServicesAuthenticationSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneAwsServicesAuthentication {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneAwsServicesAuthentication{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneAwsServicesAuthentication{}
	}

	items := make([]ClusterControlPlaneAwsServicesAuthentication, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneAwsServicesAuthentication(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterControlPlaneAwsServicesAuthentication expands an instance of ClusterControlPlaneAwsServicesAuthentication into a JSON
// request object.
func expandClusterControlPlaneAwsServicesAuthentication(c *Client, f *ClusterControlPlaneAwsServicesAuthentication, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.RoleArn; !dcl.IsEmptyValueIndirect(v) {
		m["roleArn"] = v
	}
	if v := f.RoleSessionName; !dcl.IsEmptyValueIndirect(v) {
		m["roleSessionName"] = v
	}

	return m, nil
}

// flattenClusterControlPlaneAwsServicesAuthentication flattens an instance of ClusterControlPlaneAwsServicesAuthentication from a JSON
// response object.
func flattenClusterControlPlaneAwsServicesAuthentication(c *Client, i interface{}, res *Cluster) *ClusterControlPlaneAwsServicesAuthentication {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterControlPlaneAwsServicesAuthentication{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterControlPlaneAwsServicesAuthentication
	}
	r.RoleArn = dcl.FlattenString(m["roleArn"])
	r.RoleSessionName = dcl.FlattenString(m["roleSessionName"])

	return r
}

// expandClusterControlPlaneProxyConfigMap expands the contents of ClusterControlPlaneProxyConfig into a JSON
// request object.
func expandClusterControlPlaneProxyConfigMap(c *Client, f map[string]ClusterControlPlaneProxyConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterControlPlaneProxyConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterControlPlaneProxyConfigSlice expands the contents of ClusterControlPlaneProxyConfig into a JSON
// request object.
func expandClusterControlPlaneProxyConfigSlice(c *Client, f []ClusterControlPlaneProxyConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterControlPlaneProxyConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterControlPlaneProxyConfigMap flattens the contents of ClusterControlPlaneProxyConfig from a JSON
// response object.
func flattenClusterControlPlaneProxyConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneProxyConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneProxyConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneProxyConfig{}
	}

	items := make(map[string]ClusterControlPlaneProxyConfig)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneProxyConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterControlPlaneProxyConfigSlice flattens the contents of ClusterControlPlaneProxyConfig from a JSON
// response object.
func flattenClusterControlPlaneProxyConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneProxyConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneProxyConfig{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneProxyConfig{}
	}

	items := make([]ClusterControlPlaneProxyConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneProxyConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterControlPlaneProxyConfig expands an instance of ClusterControlPlaneProxyConfig into a JSON
// request object.
func expandClusterControlPlaneProxyConfig(c *Client, f *ClusterControlPlaneProxyConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.SecretArn; !dcl.IsEmptyValueIndirect(v) {
		m["secretArn"] = v
	}
	if v := f.SecretVersion; !dcl.IsEmptyValueIndirect(v) {
		m["secretVersion"] = v
	}

	return m, nil
}

// flattenClusterControlPlaneProxyConfig flattens an instance of ClusterControlPlaneProxyConfig from a JSON
// response object.
func flattenClusterControlPlaneProxyConfig(c *Client, i interface{}, res *Cluster) *ClusterControlPlaneProxyConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterControlPlaneProxyConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterControlPlaneProxyConfig
	}
	r.SecretArn = dcl.FlattenString(m["secretArn"])
	r.SecretVersion = dcl.FlattenString(m["secretVersion"])

	return r
}

// expandClusterAuthorizationMap expands the contents of ClusterAuthorization into a JSON
// request object.
func expandClusterAuthorizationMap(c *Client, f map[string]ClusterAuthorization, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterAuthorization(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterAuthorizationSlice expands the contents of ClusterAuthorization into a JSON
// request object.
func expandClusterAuthorizationSlice(c *Client, f []ClusterAuthorization, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterAuthorization(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterAuthorizationMap flattens the contents of ClusterAuthorization from a JSON
// response object.
func flattenClusterAuthorizationMap(c *Client, i interface{}, res *Cluster) map[string]ClusterAuthorization {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterAuthorization{}
	}

	if len(a) == 0 {
		return map[string]ClusterAuthorization{}
	}

	items := make(map[string]ClusterAuthorization)
	for k, item := range a {
		items[k] = *flattenClusterAuthorization(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterAuthorizationSlice flattens the contents of ClusterAuthorization from a JSON
// response object.
func flattenClusterAuthorizationSlice(c *Client, i interface{}, res *Cluster) []ClusterAuthorization {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterAuthorization{}
	}

	if len(a) == 0 {
		return []ClusterAuthorization{}
	}

	items := make([]ClusterAuthorization, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterAuthorization(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterAuthorization expands an instance of ClusterAuthorization into a JSON
// request object.
func expandClusterAuthorization(c *Client, f *ClusterAuthorization, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandClusterAuthorizationAdminUsersSlice(c, f.AdminUsers, res); err != nil {
		return nil, fmt.Errorf("error expanding AdminUsers into adminUsers: %w", err)
	} else if v != nil {
		m["adminUsers"] = v
	}

	return m, nil
}

// flattenClusterAuthorization flattens an instance of ClusterAuthorization from a JSON
// response object.
func flattenClusterAuthorization(c *Client, i interface{}, res *Cluster) *ClusterAuthorization {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterAuthorization{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterAuthorization
	}
	r.AdminUsers = flattenClusterAuthorizationAdminUsersSlice(c, m["adminUsers"], res)

	return r
}

// expandClusterAuthorizationAdminUsersMap expands the contents of ClusterAuthorizationAdminUsers into a JSON
// request object.
func expandClusterAuthorizationAdminUsersMap(c *Client, f map[string]ClusterAuthorizationAdminUsers, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterAuthorizationAdminUsers(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterAuthorizationAdminUsersSlice expands the contents of ClusterAuthorizationAdminUsers into a JSON
// request object.
func expandClusterAuthorizationAdminUsersSlice(c *Client, f []ClusterAuthorizationAdminUsers, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterAuthorizationAdminUsers(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterAuthorizationAdminUsersMap flattens the contents of ClusterAuthorizationAdminUsers from a JSON
// response object.
func flattenClusterAuthorizationAdminUsersMap(c *Client, i interface{}, res *Cluster) map[string]ClusterAuthorizationAdminUsers {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterAuthorizationAdminUsers{}
	}

	if len(a) == 0 {
		return map[string]ClusterAuthorizationAdminUsers{}
	}

	items := make(map[string]ClusterAuthorizationAdminUsers)
	for k, item := range a {
		items[k] = *flattenClusterAuthorizationAdminUsers(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterAuthorizationAdminUsersSlice flattens the contents of ClusterAuthorizationAdminUsers from a JSON
// response object.
func flattenClusterAuthorizationAdminUsersSlice(c *Client, i interface{}, res *Cluster) []ClusterAuthorizationAdminUsers {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterAuthorizationAdminUsers{}
	}

	if len(a) == 0 {
		return []ClusterAuthorizationAdminUsers{}
	}

	items := make([]ClusterAuthorizationAdminUsers, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterAuthorizationAdminUsers(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterAuthorizationAdminUsers expands an instance of ClusterAuthorizationAdminUsers into a JSON
// request object.
func expandClusterAuthorizationAdminUsers(c *Client, f *ClusterAuthorizationAdminUsers, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Username; !dcl.IsEmptyValueIndirect(v) {
		m["username"] = v
	}

	return m, nil
}

// flattenClusterAuthorizationAdminUsers flattens an instance of ClusterAuthorizationAdminUsers from a JSON
// response object.
func flattenClusterAuthorizationAdminUsers(c *Client, i interface{}, res *Cluster) *ClusterAuthorizationAdminUsers {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterAuthorizationAdminUsers{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterAuthorizationAdminUsers
	}
	r.Username = dcl.FlattenString(m["username"])

	return r
}

// expandClusterWorkloadIdentityConfigMap expands the contents of ClusterWorkloadIdentityConfig into a JSON
// request object.
func expandClusterWorkloadIdentityConfigMap(c *Client, f map[string]ClusterWorkloadIdentityConfig, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterWorkloadIdentityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterWorkloadIdentityConfigSlice expands the contents of ClusterWorkloadIdentityConfig into a JSON
// request object.
func expandClusterWorkloadIdentityConfigSlice(c *Client, f []ClusterWorkloadIdentityConfig, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterWorkloadIdentityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterWorkloadIdentityConfigMap flattens the contents of ClusterWorkloadIdentityConfig from a JSON
// response object.
func flattenClusterWorkloadIdentityConfigMap(c *Client, i interface{}, res *Cluster) map[string]ClusterWorkloadIdentityConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterWorkloadIdentityConfig{}
	}

	if len(a) == 0 {
		return map[string]ClusterWorkloadIdentityConfig{}
	}

	items := make(map[string]ClusterWorkloadIdentityConfig)
	for k, item := range a {
		items[k] = *flattenClusterWorkloadIdentityConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterWorkloadIdentityConfigSlice flattens the contents of ClusterWorkloadIdentityConfig from a JSON
// response object.
func flattenClusterWorkloadIdentityConfigSlice(c *Client, i interface{}, res *Cluster) []ClusterWorkloadIdentityConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterWorkloadIdentityConfig{}
	}

	if len(a) == 0 {
		return []ClusterWorkloadIdentityConfig{}
	}

	items := make([]ClusterWorkloadIdentityConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterWorkloadIdentityConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterWorkloadIdentityConfig expands an instance of ClusterWorkloadIdentityConfig into a JSON
// request object.
func expandClusterWorkloadIdentityConfig(c *Client, f *ClusterWorkloadIdentityConfig, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.IssuerUri; !dcl.IsEmptyValueIndirect(v) {
		m["issuerUri"] = v
	}
	if v := f.WorkloadPool; !dcl.IsEmptyValueIndirect(v) {
		m["workloadPool"] = v
	}
	if v := f.IdentityProvider; !dcl.IsEmptyValueIndirect(v) {
		m["identityProvider"] = v
	}

	return m, nil
}

// flattenClusterWorkloadIdentityConfig flattens an instance of ClusterWorkloadIdentityConfig from a JSON
// response object.
func flattenClusterWorkloadIdentityConfig(c *Client, i interface{}, res *Cluster) *ClusterWorkloadIdentityConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterWorkloadIdentityConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterWorkloadIdentityConfig
	}
	r.IssuerUri = dcl.FlattenString(m["issuerUri"])
	r.WorkloadPool = dcl.FlattenString(m["workloadPool"])
	r.IdentityProvider = dcl.FlattenString(m["identityProvider"])

	return r
}

// expandClusterFleetMap expands the contents of ClusterFleet into a JSON
// request object.
func expandClusterFleetMap(c *Client, f map[string]ClusterFleet, res *Cluster) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandClusterFleet(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandClusterFleetSlice expands the contents of ClusterFleet into a JSON
// request object.
func expandClusterFleetSlice(c *Client, f []ClusterFleet, res *Cluster) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandClusterFleet(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenClusterFleetMap flattens the contents of ClusterFleet from a JSON
// response object.
func flattenClusterFleetMap(c *Client, i interface{}, res *Cluster) map[string]ClusterFleet {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterFleet{}
	}

	if len(a) == 0 {
		return map[string]ClusterFleet{}
	}

	items := make(map[string]ClusterFleet)
	for k, item := range a {
		items[k] = *flattenClusterFleet(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenClusterFleetSlice flattens the contents of ClusterFleet from a JSON
// response object.
func flattenClusterFleetSlice(c *Client, i interface{}, res *Cluster) []ClusterFleet {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterFleet{}
	}

	if len(a) == 0 {
		return []ClusterFleet{}
	}

	items := make([]ClusterFleet, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterFleet(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandClusterFleet expands an instance of ClusterFleet into a JSON
// request object.
func expandClusterFleet(c *Client, f *ClusterFleet, res *Cluster) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := dcl.DeriveField("projects/%s", f.Project, dcl.SelfLinkToName(f.Project)); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenClusterFleet flattens an instance of ClusterFleet from a JSON
// response object.
func flattenClusterFleet(c *Client, i interface{}, res *Cluster) *ClusterFleet {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ClusterFleet{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyClusterFleet
	}
	r.Project = dcl.FlattenString(m["project"])
	r.Membership = dcl.FlattenString(m["membership"])

	return r
}

// flattenClusterControlPlaneRootVolumeVolumeTypeEnumMap flattens the contents of ClusterControlPlaneRootVolumeVolumeTypeEnum from a JSON
// response object.
func flattenClusterControlPlaneRootVolumeVolumeTypeEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneRootVolumeVolumeTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneRootVolumeVolumeTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneRootVolumeVolumeTypeEnum{}
	}

	items := make(map[string]ClusterControlPlaneRootVolumeVolumeTypeEnum)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneRootVolumeVolumeTypeEnum(item.(interface{}))
	}

	return items
}

// flattenClusterControlPlaneRootVolumeVolumeTypeEnumSlice flattens the contents of ClusterControlPlaneRootVolumeVolumeTypeEnum from a JSON
// response object.
func flattenClusterControlPlaneRootVolumeVolumeTypeEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneRootVolumeVolumeTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneRootVolumeVolumeTypeEnum{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneRootVolumeVolumeTypeEnum{}
	}

	items := make([]ClusterControlPlaneRootVolumeVolumeTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneRootVolumeVolumeTypeEnum(item.(interface{})))
	}

	return items
}

// flattenClusterControlPlaneRootVolumeVolumeTypeEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterControlPlaneRootVolumeVolumeTypeEnum with the same value as that string.
func flattenClusterControlPlaneRootVolumeVolumeTypeEnum(i interface{}) *ClusterControlPlaneRootVolumeVolumeTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterControlPlaneRootVolumeVolumeTypeEnumRef(s)
}

// flattenClusterControlPlaneMainVolumeVolumeTypeEnumMap flattens the contents of ClusterControlPlaneMainVolumeVolumeTypeEnum from a JSON
// response object.
func flattenClusterControlPlaneMainVolumeVolumeTypeEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterControlPlaneMainVolumeVolumeTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterControlPlaneMainVolumeVolumeTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterControlPlaneMainVolumeVolumeTypeEnum{}
	}

	items := make(map[string]ClusterControlPlaneMainVolumeVolumeTypeEnum)
	for k, item := range a {
		items[k] = *flattenClusterControlPlaneMainVolumeVolumeTypeEnum(item.(interface{}))
	}

	return items
}

// flattenClusterControlPlaneMainVolumeVolumeTypeEnumSlice flattens the contents of ClusterControlPlaneMainVolumeVolumeTypeEnum from a JSON
// response object.
func flattenClusterControlPlaneMainVolumeVolumeTypeEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterControlPlaneMainVolumeVolumeTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterControlPlaneMainVolumeVolumeTypeEnum{}
	}

	if len(a) == 0 {
		return []ClusterControlPlaneMainVolumeVolumeTypeEnum{}
	}

	items := make([]ClusterControlPlaneMainVolumeVolumeTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterControlPlaneMainVolumeVolumeTypeEnum(item.(interface{})))
	}

	return items
}

// flattenClusterControlPlaneMainVolumeVolumeTypeEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterControlPlaneMainVolumeVolumeTypeEnum with the same value as that string.
func flattenClusterControlPlaneMainVolumeVolumeTypeEnum(i interface{}) *ClusterControlPlaneMainVolumeVolumeTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterControlPlaneMainVolumeVolumeTypeEnumRef(s)
}

// flattenClusterStateEnumMap flattens the contents of ClusterStateEnum from a JSON
// response object.
func flattenClusterStateEnumMap(c *Client, i interface{}, res *Cluster) map[string]ClusterStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ClusterStateEnum{}
	}

	if len(a) == 0 {
		return map[string]ClusterStateEnum{}
	}

	items := make(map[string]ClusterStateEnum)
	for k, item := range a {
		items[k] = *flattenClusterStateEnum(item.(interface{}))
	}

	return items
}

// flattenClusterStateEnumSlice flattens the contents of ClusterStateEnum from a JSON
// response object.
func flattenClusterStateEnumSlice(c *Client, i interface{}, res *Cluster) []ClusterStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ClusterStateEnum{}
	}

	if len(a) == 0 {
		return []ClusterStateEnum{}
	}

	items := make([]ClusterStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenClusterStateEnum(item.(interface{})))
	}

	return items
}

// flattenClusterStateEnum asserts that an interface is a string, and returns a
// pointer to a *ClusterStateEnum with the same value as that string.
func flattenClusterStateEnum(i interface{}) *ClusterStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ClusterStateEnumRef(s)
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

	case "updateClusterUpdateAwsClusterOperation":
		return &updateClusterUpdateAwsClusterOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractClusterFields(r *Cluster) error {
	vNetworking := r.Networking
	if vNetworking == nil {
		// note: explicitly not the empty object.
		vNetworking = &ClusterNetworking{}
	}
	if err := extractClusterNetworkingFields(r, vNetworking); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetworking) {
		r.Networking = vNetworking
	}
	vControlPlane := r.ControlPlane
	if vControlPlane == nil {
		// note: explicitly not the empty object.
		vControlPlane = &ClusterControlPlane{}
	}
	if err := extractClusterControlPlaneFields(r, vControlPlane); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vControlPlane) {
		r.ControlPlane = vControlPlane
	}
	vAuthorization := r.Authorization
	if vAuthorization == nil {
		// note: explicitly not the empty object.
		vAuthorization = &ClusterAuthorization{}
	}
	if err := extractClusterAuthorizationFields(r, vAuthorization); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorization) {
		r.Authorization = vAuthorization
	}
	vWorkloadIdentityConfig := r.WorkloadIdentityConfig
	if vWorkloadIdentityConfig == nil {
		// note: explicitly not the empty object.
		vWorkloadIdentityConfig = &ClusterWorkloadIdentityConfig{}
	}
	if err := extractClusterWorkloadIdentityConfigFields(r, vWorkloadIdentityConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkloadIdentityConfig) {
		r.WorkloadIdentityConfig = vWorkloadIdentityConfig
	}
	vFleet := r.Fleet
	if vFleet == nil {
		// note: explicitly not the empty object.
		vFleet = &ClusterFleet{}
	}
	if err := extractClusterFleetFields(r, vFleet); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFleet) {
		r.Fleet = vFleet
	}
	return nil
}
func extractClusterNetworkingFields(r *Cluster, o *ClusterNetworking) error {
	return nil
}
func extractClusterControlPlaneFields(r *Cluster, o *ClusterControlPlane) error {
	vSshConfig := o.SshConfig
	if vSshConfig == nil {
		// note: explicitly not the empty object.
		vSshConfig = &ClusterControlPlaneSshConfig{}
	}
	if err := extractClusterControlPlaneSshConfigFields(r, vSshConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSshConfig) {
		o.SshConfig = vSshConfig
	}
	vConfigEncryption := o.ConfigEncryption
	if vConfigEncryption == nil {
		// note: explicitly not the empty object.
		vConfigEncryption = &ClusterControlPlaneConfigEncryption{}
	}
	if err := extractClusterControlPlaneConfigEncryptionFields(r, vConfigEncryption); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfigEncryption) {
		o.ConfigEncryption = vConfigEncryption
	}
	vRootVolume := o.RootVolume
	if vRootVolume == nil {
		// note: explicitly not the empty object.
		vRootVolume = &ClusterControlPlaneRootVolume{}
	}
	if err := extractClusterControlPlaneRootVolumeFields(r, vRootVolume); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRootVolume) {
		o.RootVolume = vRootVolume
	}
	vMainVolume := o.MainVolume
	if vMainVolume == nil {
		// note: explicitly not the empty object.
		vMainVolume = &ClusterControlPlaneMainVolume{}
	}
	if err := extractClusterControlPlaneMainVolumeFields(r, vMainVolume); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMainVolume) {
		o.MainVolume = vMainVolume
	}
	vDatabaseEncryption := o.DatabaseEncryption
	if vDatabaseEncryption == nil {
		// note: explicitly not the empty object.
		vDatabaseEncryption = &ClusterControlPlaneDatabaseEncryption{}
	}
	if err := extractClusterControlPlaneDatabaseEncryptionFields(r, vDatabaseEncryption); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDatabaseEncryption) {
		o.DatabaseEncryption = vDatabaseEncryption
	}
	vAwsServicesAuthentication := o.AwsServicesAuthentication
	if vAwsServicesAuthentication == nil {
		// note: explicitly not the empty object.
		vAwsServicesAuthentication = &ClusterControlPlaneAwsServicesAuthentication{}
	}
	if err := extractClusterControlPlaneAwsServicesAuthenticationFields(r, vAwsServicesAuthentication); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAwsServicesAuthentication) {
		o.AwsServicesAuthentication = vAwsServicesAuthentication
	}
	vProxyConfig := o.ProxyConfig
	if vProxyConfig == nil {
		// note: explicitly not the empty object.
		vProxyConfig = &ClusterControlPlaneProxyConfig{}
	}
	if err := extractClusterControlPlaneProxyConfigFields(r, vProxyConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vProxyConfig) {
		o.ProxyConfig = vProxyConfig
	}
	return nil
}
func extractClusterControlPlaneSshConfigFields(r *Cluster, o *ClusterControlPlaneSshConfig) error {
	return nil
}
func extractClusterControlPlaneConfigEncryptionFields(r *Cluster, o *ClusterControlPlaneConfigEncryption) error {
	return nil
}
func extractClusterControlPlaneRootVolumeFields(r *Cluster, o *ClusterControlPlaneRootVolume) error {
	return nil
}
func extractClusterControlPlaneMainVolumeFields(r *Cluster, o *ClusterControlPlaneMainVolume) error {
	return nil
}
func extractClusterControlPlaneDatabaseEncryptionFields(r *Cluster, o *ClusterControlPlaneDatabaseEncryption) error {
	return nil
}
func extractClusterControlPlaneAwsServicesAuthenticationFields(r *Cluster, o *ClusterControlPlaneAwsServicesAuthentication) error {
	return nil
}
func extractClusterControlPlaneProxyConfigFields(r *Cluster, o *ClusterControlPlaneProxyConfig) error {
	return nil
}
func extractClusterAuthorizationFields(r *Cluster, o *ClusterAuthorization) error {
	return nil
}
func extractClusterAuthorizationAdminUsersFields(r *Cluster, o *ClusterAuthorizationAdminUsers) error {
	return nil
}
func extractClusterWorkloadIdentityConfigFields(r *Cluster, o *ClusterWorkloadIdentityConfig) error {
	return nil
}
func extractClusterFleetFields(r *Cluster, o *ClusterFleet) error {
	return nil
}

func postReadExtractClusterFields(r *Cluster) error {
	vNetworking := r.Networking
	if vNetworking == nil {
		// note: explicitly not the empty object.
		vNetworking = &ClusterNetworking{}
	}
	if err := postReadExtractClusterNetworkingFields(r, vNetworking); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNetworking) {
		r.Networking = vNetworking
	}
	vControlPlane := r.ControlPlane
	if vControlPlane == nil {
		// note: explicitly not the empty object.
		vControlPlane = &ClusterControlPlane{}
	}
	if err := postReadExtractClusterControlPlaneFields(r, vControlPlane); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vControlPlane) {
		r.ControlPlane = vControlPlane
	}
	vAuthorization := r.Authorization
	if vAuthorization == nil {
		// note: explicitly not the empty object.
		vAuthorization = &ClusterAuthorization{}
	}
	if err := postReadExtractClusterAuthorizationFields(r, vAuthorization); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorization) {
		r.Authorization = vAuthorization
	}
	vWorkloadIdentityConfig := r.WorkloadIdentityConfig
	if vWorkloadIdentityConfig == nil {
		// note: explicitly not the empty object.
		vWorkloadIdentityConfig = &ClusterWorkloadIdentityConfig{}
	}
	if err := postReadExtractClusterWorkloadIdentityConfigFields(r, vWorkloadIdentityConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkloadIdentityConfig) {
		r.WorkloadIdentityConfig = vWorkloadIdentityConfig
	}
	vFleet := r.Fleet
	if vFleet == nil {
		// note: explicitly not the empty object.
		vFleet = &ClusterFleet{}
	}
	if err := postReadExtractClusterFleetFields(r, vFleet); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFleet) {
		r.Fleet = vFleet
	}
	return nil
}
func postReadExtractClusterNetworkingFields(r *Cluster, o *ClusterNetworking) error {
	return nil
}
func postReadExtractClusterControlPlaneFields(r *Cluster, o *ClusterControlPlane) error {
	vSshConfig := o.SshConfig
	if vSshConfig == nil {
		// note: explicitly not the empty object.
		vSshConfig = &ClusterControlPlaneSshConfig{}
	}
	if err := extractClusterControlPlaneSshConfigFields(r, vSshConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSshConfig) {
		o.SshConfig = vSshConfig
	}
	vConfigEncryption := o.ConfigEncryption
	if vConfigEncryption == nil {
		// note: explicitly not the empty object.
		vConfigEncryption = &ClusterControlPlaneConfigEncryption{}
	}
	if err := extractClusterControlPlaneConfigEncryptionFields(r, vConfigEncryption); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfigEncryption) {
		o.ConfigEncryption = vConfigEncryption
	}
	vRootVolume := o.RootVolume
	if vRootVolume == nil {
		// note: explicitly not the empty object.
		vRootVolume = &ClusterControlPlaneRootVolume{}
	}
	if err := extractClusterControlPlaneRootVolumeFields(r, vRootVolume); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRootVolume) {
		o.RootVolume = vRootVolume
	}
	vMainVolume := o.MainVolume
	if vMainVolume == nil {
		// note: explicitly not the empty object.
		vMainVolume = &ClusterControlPlaneMainVolume{}
	}
	if err := extractClusterControlPlaneMainVolumeFields(r, vMainVolume); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMainVolume) {
		o.MainVolume = vMainVolume
	}
	vDatabaseEncryption := o.DatabaseEncryption
	if vDatabaseEncryption == nil {
		// note: explicitly not the empty object.
		vDatabaseEncryption = &ClusterControlPlaneDatabaseEncryption{}
	}
	if err := extractClusterControlPlaneDatabaseEncryptionFields(r, vDatabaseEncryption); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDatabaseEncryption) {
		o.DatabaseEncryption = vDatabaseEncryption
	}
	vAwsServicesAuthentication := o.AwsServicesAuthentication
	if vAwsServicesAuthentication == nil {
		// note: explicitly not the empty object.
		vAwsServicesAuthentication = &ClusterControlPlaneAwsServicesAuthentication{}
	}
	if err := extractClusterControlPlaneAwsServicesAuthenticationFields(r, vAwsServicesAuthentication); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAwsServicesAuthentication) {
		o.AwsServicesAuthentication = vAwsServicesAuthentication
	}
	vProxyConfig := o.ProxyConfig
	if vProxyConfig == nil {
		// note: explicitly not the empty object.
		vProxyConfig = &ClusterControlPlaneProxyConfig{}
	}
	if err := extractClusterControlPlaneProxyConfigFields(r, vProxyConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vProxyConfig) {
		o.ProxyConfig = vProxyConfig
	}
	return nil
}
func postReadExtractClusterControlPlaneSshConfigFields(r *Cluster, o *ClusterControlPlaneSshConfig) error {
	return nil
}
func postReadExtractClusterControlPlaneConfigEncryptionFields(r *Cluster, o *ClusterControlPlaneConfigEncryption) error {
	return nil
}
func postReadExtractClusterControlPlaneRootVolumeFields(r *Cluster, o *ClusterControlPlaneRootVolume) error {
	return nil
}
func postReadExtractClusterControlPlaneMainVolumeFields(r *Cluster, o *ClusterControlPlaneMainVolume) error {
	return nil
}
func postReadExtractClusterControlPlaneDatabaseEncryptionFields(r *Cluster, o *ClusterControlPlaneDatabaseEncryption) error {
	return nil
}
func postReadExtractClusterControlPlaneAwsServicesAuthenticationFields(r *Cluster, o *ClusterControlPlaneAwsServicesAuthentication) error {
	return nil
}
func postReadExtractClusterControlPlaneProxyConfigFields(r *Cluster, o *ClusterControlPlaneProxyConfig) error {
	return nil
}
func postReadExtractClusterAuthorizationFields(r *Cluster, o *ClusterAuthorization) error {
	return nil
}
func postReadExtractClusterAuthorizationAdminUsersFields(r *Cluster, o *ClusterAuthorizationAdminUsers) error {
	return nil
}
func postReadExtractClusterWorkloadIdentityConfigFields(r *Cluster, o *ClusterWorkloadIdentityConfig) error {
	return nil
}
func postReadExtractClusterFleetFields(r *Cluster, o *ClusterFleet) error {
	return nil
}
