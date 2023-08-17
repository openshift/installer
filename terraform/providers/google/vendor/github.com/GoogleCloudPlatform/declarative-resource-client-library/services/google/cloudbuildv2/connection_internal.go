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
package cloudbuildv2

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

func (r *Connection) validate() error {

	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"GithubConfig", "GithubEnterpriseConfig", "GitlabConfig"}, r.GithubConfig, r.GithubEnterpriseConfig, r.GitlabConfig); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Name, "Name"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.GithubConfig) {
		if err := r.GithubConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.GithubEnterpriseConfig) {
		if err := r.GithubEnterpriseConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.GitlabConfig) {
		if err := r.GitlabConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.InstallationState) {
		if err := r.InstallationState.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ConnectionGithubConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.AuthorizerCredential) {
		if err := r.AuthorizerCredential.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ConnectionGithubConfigAuthorizerCredential) validate() error {
	return nil
}
func (r *ConnectionGithubEnterpriseConfig) validate() error {
	if err := dcl.Required(r, "hostUri"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.ServiceDirectoryConfig) {
		if err := r.ServiceDirectoryConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ConnectionGithubEnterpriseConfigServiceDirectoryConfig) validate() error {
	if err := dcl.Required(r, "service"); err != nil {
		return err
	}
	return nil
}
func (r *ConnectionGitlabConfig) validate() error {
	if err := dcl.Required(r, "webhookSecretSecretVersion"); err != nil {
		return err
	}
	if err := dcl.Required(r, "readAuthorizerCredential"); err != nil {
		return err
	}
	if err := dcl.Required(r, "authorizerCredential"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.ReadAuthorizerCredential) {
		if err := r.ReadAuthorizerCredential.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AuthorizerCredential) {
		if err := r.AuthorizerCredential.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ServiceDirectoryConfig) {
		if err := r.ServiceDirectoryConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *ConnectionGitlabConfigReadAuthorizerCredential) validate() error {
	if err := dcl.Required(r, "userTokenSecretVersion"); err != nil {
		return err
	}
	return nil
}
func (r *ConnectionGitlabConfigAuthorizerCredential) validate() error {
	if err := dcl.Required(r, "userTokenSecretVersion"); err != nil {
		return err
	}
	return nil
}
func (r *ConnectionGitlabConfigServiceDirectoryConfig) validate() error {
	if err := dcl.Required(r, "service"); err != nil {
		return err
	}
	return nil
}
func (r *ConnectionInstallationState) validate() error {
	return nil
}
func (r *Connection) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://cloudbuild.googleapis.com/v2/", params)
}

func (r *Connection) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/connections/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Connection) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/connections", nr.basePath(), userBasePath, params), nil

}

func (r *Connection) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/connections?connectionId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *Connection) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/connections/{{name}}", nr.basePath(), userBasePath, params), nil
}

// connectionApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type connectionApiOperation interface {
	do(context.Context, *Connection, *Client) error
}

// newUpdateConnectionUpdateConnectionRequest creates a request for an
// Connection resource's UpdateConnection update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateConnectionUpdateConnectionRequest(ctx context.Context, f *Connection, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := expandConnectionGithubConfig(c, f.GithubConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GithubConfig into githubConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["githubConfig"] = v
	}
	if v, err := expandConnectionGithubEnterpriseConfig(c, f.GithubEnterpriseConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GithubEnterpriseConfig into githubEnterpriseConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["githubEnterpriseConfig"] = v
	}
	if v, err := expandConnectionGitlabConfig(c, f.GitlabConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GitlabConfig into gitlabConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["gitlabConfig"] = v
	}
	if v := f.Disabled; !dcl.IsEmptyValueIndirect(v) {
		req["disabled"] = v
	}
	if v := f.Annotations; !dcl.IsEmptyValueIndirect(v) {
		req["annotations"] = v
	}
	b, err := c.getConnectionRaw(ctx, f)
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
	req["name"] = fmt.Sprintf("projects/%s/locations/%s/connections/%s", *f.Project, *f.Location, *f.Name)

	return req, nil
}

// marshalUpdateConnectionUpdateConnectionRequest converts the update into
// the final JSON request body.
func marshalUpdateConnectionUpdateConnectionRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateConnectionUpdateConnectionOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateConnectionUpdateConnectionOperation) do(ctx context.Context, r *Connection, c *Client) error {
	_, err := c.GetConnection(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateConnection")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateConnectionUpdateConnectionRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateConnectionUpdateConnectionRequest(c, req)
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

func (c *Client) listConnectionRaw(ctx context.Context, r *Connection, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != ConnectionMaxPage {
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

type listConnectionOperation struct {
	Connections []map[string]interface{} `json:"connections"`
	Token       string                   `json:"nextPageToken"`
}

func (c *Client) listConnection(ctx context.Context, r *Connection, pageToken string, pageSize int32) ([]*Connection, string, error) {
	b, err := c.listConnectionRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listConnectionOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Connection
	for _, v := range m.Connections {
		res, err := unmarshalMapConnection(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllConnection(ctx context.Context, f func(*Connection) bool, resources []*Connection) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteConnection(ctx, res)
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

type deleteConnectionOperation struct{}

func (op *deleteConnectionOperation) do(ctx context.Context, r *Connection, c *Client) error {
	r, err := c.GetConnection(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Connection not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetConnection checking for existence. error: %v", err)
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
		_, err := c.GetConnection(ctx, r)
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
type createConnectionOperation struct {
	response map[string]interface{}
}

func (op *createConnectionOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createConnectionOperation) do(ctx context.Context, r *Connection, c *Client) error {
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

	if _, err := c.GetConnection(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getConnectionRaw(ctx context.Context, r *Connection) ([]byte, error) {

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

func (c *Client) connectionDiffsForRawDesired(ctx context.Context, rawDesired *Connection, opts ...dcl.ApplyOption) (initial, desired *Connection, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Connection
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Connection); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Connection, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetConnection(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Connection resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Connection resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Connection resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeConnectionDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Connection: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Connection: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractConnectionFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeConnectionInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Connection: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeConnectionDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Connection: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffConnection(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeConnectionInitialState(rawInitial, rawDesired *Connection) (*Connection, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.

	if !dcl.IsZeroValue(rawInitial.GithubConfig) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.GithubEnterpriseConfig, rawInitial.GitlabConfig) {
			rawInitial.GithubConfig = EmptyConnectionGithubConfig
		}
	}

	if !dcl.IsZeroValue(rawInitial.GithubEnterpriseConfig) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.GithubConfig, rawInitial.GitlabConfig) {
			rawInitial.GithubEnterpriseConfig = EmptyConnectionGithubEnterpriseConfig
		}
	}

	if !dcl.IsZeroValue(rawInitial.GitlabConfig) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.GithubConfig, rawInitial.GithubEnterpriseConfig) {
			rawInitial.GitlabConfig = EmptyConnectionGitlabConfig
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

func canonicalizeConnectionDesiredState(rawDesired, rawInitial *Connection, opts ...dcl.ApplyOption) (*Connection, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.GithubConfig = canonicalizeConnectionGithubConfig(rawDesired.GithubConfig, nil, opts...)
		rawDesired.GithubEnterpriseConfig = canonicalizeConnectionGithubEnterpriseConfig(rawDesired.GithubEnterpriseConfig, nil, opts...)
		rawDesired.GitlabConfig = canonicalizeConnectionGitlabConfig(rawDesired.GitlabConfig, nil, opts...)
		rawDesired.InstallationState = canonicalizeConnectionInstallationState(rawDesired.InstallationState, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Connection{}
	if dcl.NameToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	canonicalDesired.GithubConfig = canonicalizeConnectionGithubConfig(rawDesired.GithubConfig, rawInitial.GithubConfig, opts...)
	canonicalDesired.GithubEnterpriseConfig = canonicalizeConnectionGithubEnterpriseConfig(rawDesired.GithubEnterpriseConfig, rawInitial.GithubEnterpriseConfig, opts...)
	canonicalDesired.GitlabConfig = canonicalizeConnectionGitlabConfig(rawDesired.GitlabConfig, rawInitial.GitlabConfig, opts...)
	if dcl.BoolCanonicalize(rawDesired.Disabled, rawInitial.Disabled) {
		canonicalDesired.Disabled = rawInitial.Disabled
	} else {
		canonicalDesired.Disabled = rawDesired.Disabled
	}
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

	if canonicalDesired.GithubConfig != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.GithubEnterpriseConfig, rawDesired.GitlabConfig) {
			canonicalDesired.GithubConfig = EmptyConnectionGithubConfig
		}
	}

	if canonicalDesired.GithubEnterpriseConfig != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.GithubConfig, rawDesired.GitlabConfig) {
			canonicalDesired.GithubEnterpriseConfig = EmptyConnectionGithubEnterpriseConfig
		}
	}

	if canonicalDesired.GitlabConfig != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.GithubConfig, rawDesired.GithubEnterpriseConfig) {
			canonicalDesired.GitlabConfig = EmptyConnectionGitlabConfig
		}
	}

	return canonicalDesired, nil
}

func canonicalizeConnectionNewState(c *Client, rawNew, rawDesired *Connection) (*Connection, error) {

	rawNew.Name = rawDesired.Name

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.GithubConfig) && dcl.IsEmptyValueIndirect(rawDesired.GithubConfig) {
		rawNew.GithubConfig = rawDesired.GithubConfig
	} else {
		rawNew.GithubConfig = canonicalizeNewConnectionGithubConfig(c, rawDesired.GithubConfig, rawNew.GithubConfig)
	}

	if dcl.IsEmptyValueIndirect(rawNew.GithubEnterpriseConfig) && dcl.IsEmptyValueIndirect(rawDesired.GithubEnterpriseConfig) {
		rawNew.GithubEnterpriseConfig = rawDesired.GithubEnterpriseConfig
	} else {
		rawNew.GithubEnterpriseConfig = canonicalizeNewConnectionGithubEnterpriseConfig(c, rawDesired.GithubEnterpriseConfig, rawNew.GithubEnterpriseConfig)
	}

	if dcl.IsEmptyValueIndirect(rawNew.GitlabConfig) && dcl.IsEmptyValueIndirect(rawDesired.GitlabConfig) {
		rawNew.GitlabConfig = rawDesired.GitlabConfig
	} else {
		rawNew.GitlabConfig = canonicalizeNewConnectionGitlabConfig(c, rawDesired.GitlabConfig, rawNew.GitlabConfig)
	}

	if dcl.IsEmptyValueIndirect(rawNew.InstallationState) && dcl.IsEmptyValueIndirect(rawDesired.InstallationState) {
		rawNew.InstallationState = rawDesired.InstallationState
	} else {
		rawNew.InstallationState = canonicalizeNewConnectionInstallationState(c, rawDesired.InstallationState, rawNew.InstallationState)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Disabled) && dcl.IsEmptyValueIndirect(rawDesired.Disabled) {
		rawNew.Disabled = rawDesired.Disabled
	} else {
		if dcl.BoolCanonicalize(rawDesired.Disabled, rawNew.Disabled) {
			rawNew.Disabled = rawDesired.Disabled
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Reconciling) && dcl.IsEmptyValueIndirect(rawDesired.Reconciling) {
		rawNew.Reconciling = rawDesired.Reconciling
	} else {
		if dcl.BoolCanonicalize(rawDesired.Reconciling, rawNew.Reconciling) {
			rawNew.Reconciling = rawDesired.Reconciling
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Annotations) && dcl.IsEmptyValueIndirect(rawDesired.Annotations) {
		rawNew.Annotations = rawDesired.Annotations
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Etag) && dcl.IsEmptyValueIndirect(rawDesired.Etag) {
		rawNew.Etag = rawDesired.Etag
	} else {
		if dcl.StringCanonicalize(rawDesired.Etag, rawNew.Etag) {
			rawNew.Etag = rawDesired.Etag
		}
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeConnectionGithubConfig(des, initial *ConnectionGithubConfig, opts ...dcl.ApplyOption) *ConnectionGithubConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionGithubConfig{}

	cDes.AuthorizerCredential = canonicalizeConnectionGithubConfigAuthorizerCredential(des.AuthorizerCredential, initial.AuthorizerCredential, opts...)
	if dcl.IsZeroValue(des.AppInstallationId) || (dcl.IsEmptyValueIndirect(des.AppInstallationId) && dcl.IsEmptyValueIndirect(initial.AppInstallationId)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AppInstallationId = initial.AppInstallationId
	} else {
		cDes.AppInstallationId = des.AppInstallationId
	}

	return cDes
}

func canonicalizeConnectionGithubConfigSlice(des, initial []ConnectionGithubConfig, opts ...dcl.ApplyOption) []ConnectionGithubConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionGithubConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionGithubConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionGithubConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionGithubConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionGithubConfig(c *Client, des, nw *ConnectionGithubConfig) *ConnectionGithubConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionGithubConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.AuthorizerCredential = canonicalizeNewConnectionGithubConfigAuthorizerCredential(c, des.AuthorizerCredential, nw.AuthorizerCredential)

	return nw
}

func canonicalizeNewConnectionGithubConfigSet(c *Client, des, nw []ConnectionGithubConfig) []ConnectionGithubConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionGithubConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionGithubConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionGithubConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionGithubConfigSlice(c *Client, des, nw []ConnectionGithubConfig) []ConnectionGithubConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionGithubConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionGithubConfig(c, &d, &n))
	}

	return items
}

func canonicalizeConnectionGithubConfigAuthorizerCredential(des, initial *ConnectionGithubConfigAuthorizerCredential, opts ...dcl.ApplyOption) *ConnectionGithubConfigAuthorizerCredential {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionGithubConfigAuthorizerCredential{}

	if dcl.IsZeroValue(des.OAuthTokenSecretVersion) || (dcl.IsEmptyValueIndirect(des.OAuthTokenSecretVersion) && dcl.IsEmptyValueIndirect(initial.OAuthTokenSecretVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.OAuthTokenSecretVersion = initial.OAuthTokenSecretVersion
	} else {
		cDes.OAuthTokenSecretVersion = des.OAuthTokenSecretVersion
	}

	return cDes
}

func canonicalizeConnectionGithubConfigAuthorizerCredentialSlice(des, initial []ConnectionGithubConfigAuthorizerCredential, opts ...dcl.ApplyOption) []ConnectionGithubConfigAuthorizerCredential {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionGithubConfigAuthorizerCredential, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionGithubConfigAuthorizerCredential(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionGithubConfigAuthorizerCredential, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionGithubConfigAuthorizerCredential(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionGithubConfigAuthorizerCredential(c *Client, des, nw *ConnectionGithubConfigAuthorizerCredential) *ConnectionGithubConfigAuthorizerCredential {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionGithubConfigAuthorizerCredential while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Username, nw.Username) {
		nw.Username = des.Username
	}

	return nw
}

func canonicalizeNewConnectionGithubConfigAuthorizerCredentialSet(c *Client, des, nw []ConnectionGithubConfigAuthorizerCredential) []ConnectionGithubConfigAuthorizerCredential {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionGithubConfigAuthorizerCredential
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionGithubConfigAuthorizerCredentialNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionGithubConfigAuthorizerCredential(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionGithubConfigAuthorizerCredentialSlice(c *Client, des, nw []ConnectionGithubConfigAuthorizerCredential) []ConnectionGithubConfigAuthorizerCredential {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionGithubConfigAuthorizerCredential
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionGithubConfigAuthorizerCredential(c, &d, &n))
	}

	return items
}

func canonicalizeConnectionGithubEnterpriseConfig(des, initial *ConnectionGithubEnterpriseConfig, opts ...dcl.ApplyOption) *ConnectionGithubEnterpriseConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionGithubEnterpriseConfig{}

	if dcl.StringCanonicalize(des.HostUri, initial.HostUri) || dcl.IsZeroValue(des.HostUri) {
		cDes.HostUri = initial.HostUri
	} else {
		cDes.HostUri = des.HostUri
	}
	if dcl.IsZeroValue(des.AppId) || (dcl.IsEmptyValueIndirect(des.AppId) && dcl.IsEmptyValueIndirect(initial.AppId)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AppId = initial.AppId
	} else {
		cDes.AppId = des.AppId
	}
	if dcl.StringCanonicalize(des.AppSlug, initial.AppSlug) || dcl.IsZeroValue(des.AppSlug) {
		cDes.AppSlug = initial.AppSlug
	} else {
		cDes.AppSlug = des.AppSlug
	}
	if dcl.IsZeroValue(des.PrivateKeySecretVersion) || (dcl.IsEmptyValueIndirect(des.PrivateKeySecretVersion) && dcl.IsEmptyValueIndirect(initial.PrivateKeySecretVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.PrivateKeySecretVersion = initial.PrivateKeySecretVersion
	} else {
		cDes.PrivateKeySecretVersion = des.PrivateKeySecretVersion
	}
	if dcl.IsZeroValue(des.WebhookSecretSecretVersion) || (dcl.IsEmptyValueIndirect(des.WebhookSecretSecretVersion) && dcl.IsEmptyValueIndirect(initial.WebhookSecretSecretVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.WebhookSecretSecretVersion = initial.WebhookSecretSecretVersion
	} else {
		cDes.WebhookSecretSecretVersion = des.WebhookSecretSecretVersion
	}
	if dcl.IsZeroValue(des.AppInstallationId) || (dcl.IsEmptyValueIndirect(des.AppInstallationId) && dcl.IsEmptyValueIndirect(initial.AppInstallationId)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AppInstallationId = initial.AppInstallationId
	} else {
		cDes.AppInstallationId = des.AppInstallationId
	}
	cDes.ServiceDirectoryConfig = canonicalizeConnectionGithubEnterpriseConfigServiceDirectoryConfig(des.ServiceDirectoryConfig, initial.ServiceDirectoryConfig, opts...)
	if dcl.StringCanonicalize(des.SslCa, initial.SslCa) || dcl.IsZeroValue(des.SslCa) {
		cDes.SslCa = initial.SslCa
	} else {
		cDes.SslCa = des.SslCa
	}

	return cDes
}

func canonicalizeConnectionGithubEnterpriseConfigSlice(des, initial []ConnectionGithubEnterpriseConfig, opts ...dcl.ApplyOption) []ConnectionGithubEnterpriseConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionGithubEnterpriseConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionGithubEnterpriseConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionGithubEnterpriseConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionGithubEnterpriseConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionGithubEnterpriseConfig(c *Client, des, nw *ConnectionGithubEnterpriseConfig) *ConnectionGithubEnterpriseConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionGithubEnterpriseConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.HostUri, nw.HostUri) {
		nw.HostUri = des.HostUri
	}
	if dcl.StringCanonicalize(des.AppSlug, nw.AppSlug) {
		nw.AppSlug = des.AppSlug
	}
	nw.ServiceDirectoryConfig = canonicalizeNewConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, des.ServiceDirectoryConfig, nw.ServiceDirectoryConfig)
	if dcl.StringCanonicalize(des.SslCa, nw.SslCa) {
		nw.SslCa = des.SslCa
	}

	return nw
}

func canonicalizeNewConnectionGithubEnterpriseConfigSet(c *Client, des, nw []ConnectionGithubEnterpriseConfig) []ConnectionGithubEnterpriseConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionGithubEnterpriseConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionGithubEnterpriseConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionGithubEnterpriseConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionGithubEnterpriseConfigSlice(c *Client, des, nw []ConnectionGithubEnterpriseConfig) []ConnectionGithubEnterpriseConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionGithubEnterpriseConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionGithubEnterpriseConfig(c, &d, &n))
	}

	return items
}

func canonicalizeConnectionGithubEnterpriseConfigServiceDirectoryConfig(des, initial *ConnectionGithubEnterpriseConfigServiceDirectoryConfig, opts ...dcl.ApplyOption) *ConnectionGithubEnterpriseConfigServiceDirectoryConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionGithubEnterpriseConfigServiceDirectoryConfig{}

	if dcl.IsZeroValue(des.Service) || (dcl.IsEmptyValueIndirect(des.Service) && dcl.IsEmptyValueIndirect(initial.Service)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Service = initial.Service
	} else {
		cDes.Service = des.Service
	}

	return cDes
}

func canonicalizeConnectionGithubEnterpriseConfigServiceDirectoryConfigSlice(des, initial []ConnectionGithubEnterpriseConfigServiceDirectoryConfig, opts ...dcl.ApplyOption) []ConnectionGithubEnterpriseConfigServiceDirectoryConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionGithubEnterpriseConfigServiceDirectoryConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionGithubEnterpriseConfigServiceDirectoryConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionGithubEnterpriseConfigServiceDirectoryConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionGithubEnterpriseConfigServiceDirectoryConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionGithubEnterpriseConfigServiceDirectoryConfig(c *Client, des, nw *ConnectionGithubEnterpriseConfigServiceDirectoryConfig) *ConnectionGithubEnterpriseConfigServiceDirectoryConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionGithubEnterpriseConfigServiceDirectoryConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewConnectionGithubEnterpriseConfigServiceDirectoryConfigSet(c *Client, des, nw []ConnectionGithubEnterpriseConfigServiceDirectoryConfig) []ConnectionGithubEnterpriseConfigServiceDirectoryConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionGithubEnterpriseConfigServiceDirectoryConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionGithubEnterpriseConfigServiceDirectoryConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionGithubEnterpriseConfigServiceDirectoryConfigSlice(c *Client, des, nw []ConnectionGithubEnterpriseConfigServiceDirectoryConfig) []ConnectionGithubEnterpriseConfigServiceDirectoryConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionGithubEnterpriseConfigServiceDirectoryConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, &d, &n))
	}

	return items
}

func canonicalizeConnectionGitlabConfig(des, initial *ConnectionGitlabConfig, opts ...dcl.ApplyOption) *ConnectionGitlabConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionGitlabConfig{}

	if dcl.StringCanonicalize(des.HostUri, initial.HostUri) || dcl.IsZeroValue(des.HostUri) {
		cDes.HostUri = initial.HostUri
	} else {
		cDes.HostUri = des.HostUri
	}
	if dcl.IsZeroValue(des.WebhookSecretSecretVersion) || (dcl.IsEmptyValueIndirect(des.WebhookSecretSecretVersion) && dcl.IsEmptyValueIndirect(initial.WebhookSecretSecretVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.WebhookSecretSecretVersion = initial.WebhookSecretSecretVersion
	} else {
		cDes.WebhookSecretSecretVersion = des.WebhookSecretSecretVersion
	}
	cDes.ReadAuthorizerCredential = canonicalizeConnectionGitlabConfigReadAuthorizerCredential(des.ReadAuthorizerCredential, initial.ReadAuthorizerCredential, opts...)
	cDes.AuthorizerCredential = canonicalizeConnectionGitlabConfigAuthorizerCredential(des.AuthorizerCredential, initial.AuthorizerCredential, opts...)
	cDes.ServiceDirectoryConfig = canonicalizeConnectionGitlabConfigServiceDirectoryConfig(des.ServiceDirectoryConfig, initial.ServiceDirectoryConfig, opts...)
	if dcl.StringCanonicalize(des.SslCa, initial.SslCa) || dcl.IsZeroValue(des.SslCa) {
		cDes.SslCa = initial.SslCa
	} else {
		cDes.SslCa = des.SslCa
	}

	return cDes
}

func canonicalizeConnectionGitlabConfigSlice(des, initial []ConnectionGitlabConfig, opts ...dcl.ApplyOption) []ConnectionGitlabConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionGitlabConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionGitlabConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionGitlabConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionGitlabConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionGitlabConfig(c *Client, des, nw *ConnectionGitlabConfig) *ConnectionGitlabConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionGitlabConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.HostUri, nw.HostUri) {
		nw.HostUri = des.HostUri
	}
	nw.ReadAuthorizerCredential = canonicalizeNewConnectionGitlabConfigReadAuthorizerCredential(c, des.ReadAuthorizerCredential, nw.ReadAuthorizerCredential)
	nw.AuthorizerCredential = canonicalizeNewConnectionGitlabConfigAuthorizerCredential(c, des.AuthorizerCredential, nw.AuthorizerCredential)
	nw.ServiceDirectoryConfig = canonicalizeNewConnectionGitlabConfigServiceDirectoryConfig(c, des.ServiceDirectoryConfig, nw.ServiceDirectoryConfig)
	if dcl.StringCanonicalize(des.SslCa, nw.SslCa) {
		nw.SslCa = des.SslCa
	}
	if dcl.StringCanonicalize(des.ServerVersion, nw.ServerVersion) {
		nw.ServerVersion = des.ServerVersion
	}

	return nw
}

func canonicalizeNewConnectionGitlabConfigSet(c *Client, des, nw []ConnectionGitlabConfig) []ConnectionGitlabConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionGitlabConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionGitlabConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionGitlabConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionGitlabConfigSlice(c *Client, des, nw []ConnectionGitlabConfig) []ConnectionGitlabConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionGitlabConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionGitlabConfig(c, &d, &n))
	}

	return items
}

func canonicalizeConnectionGitlabConfigReadAuthorizerCredential(des, initial *ConnectionGitlabConfigReadAuthorizerCredential, opts ...dcl.ApplyOption) *ConnectionGitlabConfigReadAuthorizerCredential {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionGitlabConfigReadAuthorizerCredential{}

	if dcl.IsZeroValue(des.UserTokenSecretVersion) || (dcl.IsEmptyValueIndirect(des.UserTokenSecretVersion) && dcl.IsEmptyValueIndirect(initial.UserTokenSecretVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.UserTokenSecretVersion = initial.UserTokenSecretVersion
	} else {
		cDes.UserTokenSecretVersion = des.UserTokenSecretVersion
	}

	return cDes
}

func canonicalizeConnectionGitlabConfigReadAuthorizerCredentialSlice(des, initial []ConnectionGitlabConfigReadAuthorizerCredential, opts ...dcl.ApplyOption) []ConnectionGitlabConfigReadAuthorizerCredential {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionGitlabConfigReadAuthorizerCredential, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionGitlabConfigReadAuthorizerCredential(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionGitlabConfigReadAuthorizerCredential, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionGitlabConfigReadAuthorizerCredential(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionGitlabConfigReadAuthorizerCredential(c *Client, des, nw *ConnectionGitlabConfigReadAuthorizerCredential) *ConnectionGitlabConfigReadAuthorizerCredential {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionGitlabConfigReadAuthorizerCredential while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Username, nw.Username) {
		nw.Username = des.Username
	}

	return nw
}

func canonicalizeNewConnectionGitlabConfigReadAuthorizerCredentialSet(c *Client, des, nw []ConnectionGitlabConfigReadAuthorizerCredential) []ConnectionGitlabConfigReadAuthorizerCredential {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionGitlabConfigReadAuthorizerCredential
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionGitlabConfigReadAuthorizerCredentialNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionGitlabConfigReadAuthorizerCredential(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionGitlabConfigReadAuthorizerCredentialSlice(c *Client, des, nw []ConnectionGitlabConfigReadAuthorizerCredential) []ConnectionGitlabConfigReadAuthorizerCredential {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionGitlabConfigReadAuthorizerCredential
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionGitlabConfigReadAuthorizerCredential(c, &d, &n))
	}

	return items
}

func canonicalizeConnectionGitlabConfigAuthorizerCredential(des, initial *ConnectionGitlabConfigAuthorizerCredential, opts ...dcl.ApplyOption) *ConnectionGitlabConfigAuthorizerCredential {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionGitlabConfigAuthorizerCredential{}

	if dcl.IsZeroValue(des.UserTokenSecretVersion) || (dcl.IsEmptyValueIndirect(des.UserTokenSecretVersion) && dcl.IsEmptyValueIndirect(initial.UserTokenSecretVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.UserTokenSecretVersion = initial.UserTokenSecretVersion
	} else {
		cDes.UserTokenSecretVersion = des.UserTokenSecretVersion
	}

	return cDes
}

func canonicalizeConnectionGitlabConfigAuthorizerCredentialSlice(des, initial []ConnectionGitlabConfigAuthorizerCredential, opts ...dcl.ApplyOption) []ConnectionGitlabConfigAuthorizerCredential {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionGitlabConfigAuthorizerCredential, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionGitlabConfigAuthorizerCredential(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionGitlabConfigAuthorizerCredential, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionGitlabConfigAuthorizerCredential(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionGitlabConfigAuthorizerCredential(c *Client, des, nw *ConnectionGitlabConfigAuthorizerCredential) *ConnectionGitlabConfigAuthorizerCredential {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionGitlabConfigAuthorizerCredential while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Username, nw.Username) {
		nw.Username = des.Username
	}

	return nw
}

func canonicalizeNewConnectionGitlabConfigAuthorizerCredentialSet(c *Client, des, nw []ConnectionGitlabConfigAuthorizerCredential) []ConnectionGitlabConfigAuthorizerCredential {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionGitlabConfigAuthorizerCredential
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionGitlabConfigAuthorizerCredentialNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionGitlabConfigAuthorizerCredential(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionGitlabConfigAuthorizerCredentialSlice(c *Client, des, nw []ConnectionGitlabConfigAuthorizerCredential) []ConnectionGitlabConfigAuthorizerCredential {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionGitlabConfigAuthorizerCredential
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionGitlabConfigAuthorizerCredential(c, &d, &n))
	}

	return items
}

func canonicalizeConnectionGitlabConfigServiceDirectoryConfig(des, initial *ConnectionGitlabConfigServiceDirectoryConfig, opts ...dcl.ApplyOption) *ConnectionGitlabConfigServiceDirectoryConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionGitlabConfigServiceDirectoryConfig{}

	if dcl.IsZeroValue(des.Service) || (dcl.IsEmptyValueIndirect(des.Service) && dcl.IsEmptyValueIndirect(initial.Service)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Service = initial.Service
	} else {
		cDes.Service = des.Service
	}

	return cDes
}

func canonicalizeConnectionGitlabConfigServiceDirectoryConfigSlice(des, initial []ConnectionGitlabConfigServiceDirectoryConfig, opts ...dcl.ApplyOption) []ConnectionGitlabConfigServiceDirectoryConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionGitlabConfigServiceDirectoryConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionGitlabConfigServiceDirectoryConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionGitlabConfigServiceDirectoryConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionGitlabConfigServiceDirectoryConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionGitlabConfigServiceDirectoryConfig(c *Client, des, nw *ConnectionGitlabConfigServiceDirectoryConfig) *ConnectionGitlabConfigServiceDirectoryConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionGitlabConfigServiceDirectoryConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewConnectionGitlabConfigServiceDirectoryConfigSet(c *Client, des, nw []ConnectionGitlabConfigServiceDirectoryConfig) []ConnectionGitlabConfigServiceDirectoryConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionGitlabConfigServiceDirectoryConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionGitlabConfigServiceDirectoryConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionGitlabConfigServiceDirectoryConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionGitlabConfigServiceDirectoryConfigSlice(c *Client, des, nw []ConnectionGitlabConfigServiceDirectoryConfig) []ConnectionGitlabConfigServiceDirectoryConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionGitlabConfigServiceDirectoryConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionGitlabConfigServiceDirectoryConfig(c, &d, &n))
	}

	return items
}

func canonicalizeConnectionInstallationState(des, initial *ConnectionInstallationState, opts ...dcl.ApplyOption) *ConnectionInstallationState {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &ConnectionInstallationState{}

	return cDes
}

func canonicalizeConnectionInstallationStateSlice(des, initial []ConnectionInstallationState, opts ...dcl.ApplyOption) []ConnectionInstallationState {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]ConnectionInstallationState, 0, len(des))
		for _, d := range des {
			cd := canonicalizeConnectionInstallationState(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]ConnectionInstallationState, 0, len(des))
	for i, d := range des {
		cd := canonicalizeConnectionInstallationState(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewConnectionInstallationState(c *Client, des, nw *ConnectionInstallationState) *ConnectionInstallationState {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for ConnectionInstallationState while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Message, nw.Message) {
		nw.Message = des.Message
	}
	if dcl.StringCanonicalize(des.ActionUri, nw.ActionUri) {
		nw.ActionUri = des.ActionUri
	}

	return nw
}

func canonicalizeNewConnectionInstallationStateSet(c *Client, des, nw []ConnectionInstallationState) []ConnectionInstallationState {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []ConnectionInstallationState
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareConnectionInstallationStateNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewConnectionInstallationState(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewConnectionInstallationStateSlice(c *Client, des, nw []ConnectionInstallationState) []ConnectionInstallationState {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []ConnectionInstallationState
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewConnectionInstallationState(c, &d, &n))
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
func diffConnection(c *Client, desired, actual *Connection, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.GithubConfig, actual.GithubConfig, dcl.DiffInfo{ObjectFunction: compareConnectionGithubConfigNewStyle, EmptyObject: EmptyConnectionGithubConfig, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("GithubConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GithubEnterpriseConfig, actual.GithubEnterpriseConfig, dcl.DiffInfo{ObjectFunction: compareConnectionGithubEnterpriseConfigNewStyle, EmptyObject: EmptyConnectionGithubEnterpriseConfig, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("GithubEnterpriseConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GitlabConfig, actual.GitlabConfig, dcl.DiffInfo{ObjectFunction: compareConnectionGitlabConfigNewStyle, EmptyObject: EmptyConnectionGitlabConfig, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("GitlabConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstallationState, actual.InstallationState, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareConnectionInstallationStateNewStyle, EmptyObject: EmptyConnectionInstallationState, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstallationState")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Disabled, actual.Disabled, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("Disabled")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Annotations, actual.Annotations, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("Annotations")); len(ds) != 0 || err != nil {
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
func compareConnectionGithubConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionGithubConfig)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionGithubConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGithubConfig or *ConnectionGithubConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionGithubConfig)
	if !ok {
		actualNotPointer, ok := a.(ConnectionGithubConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGithubConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AuthorizerCredential, actual.AuthorizerCredential, dcl.DiffInfo{ObjectFunction: compareConnectionGithubConfigAuthorizerCredentialNewStyle, EmptyObject: EmptyConnectionGithubConfigAuthorizerCredential, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("AuthorizerCredential")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AppInstallationId, actual.AppInstallationId, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("AppInstallationId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareConnectionGithubConfigAuthorizerCredentialNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionGithubConfigAuthorizerCredential)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionGithubConfigAuthorizerCredential)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGithubConfigAuthorizerCredential or *ConnectionGithubConfigAuthorizerCredential", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionGithubConfigAuthorizerCredential)
	if !ok {
		actualNotPointer, ok := a.(ConnectionGithubConfigAuthorizerCredential)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGithubConfigAuthorizerCredential", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.OAuthTokenSecretVersion, actual.OAuthTokenSecretVersion, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("OauthTokenSecretVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Username, actual.Username, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Username")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareConnectionGithubEnterpriseConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionGithubEnterpriseConfig)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionGithubEnterpriseConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGithubEnterpriseConfig or *ConnectionGithubEnterpriseConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionGithubEnterpriseConfig)
	if !ok {
		actualNotPointer, ok := a.(ConnectionGithubEnterpriseConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGithubEnterpriseConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.HostUri, actual.HostUri, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("HostUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AppId, actual.AppId, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("AppId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AppSlug, actual.AppSlug, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("AppSlug")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrivateKeySecretVersion, actual.PrivateKeySecretVersion, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("PrivateKeySecretVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WebhookSecretSecretVersion, actual.WebhookSecretSecretVersion, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("WebhookSecretSecretVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AppInstallationId, actual.AppInstallationId, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("AppInstallationId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceDirectoryConfig, actual.ServiceDirectoryConfig, dcl.DiffInfo{ObjectFunction: compareConnectionGithubEnterpriseConfigServiceDirectoryConfigNewStyle, EmptyObject: EmptyConnectionGithubEnterpriseConfigServiceDirectoryConfig, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("ServiceDirectoryConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SslCa, actual.SslCa, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("SslCa")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareConnectionGithubEnterpriseConfigServiceDirectoryConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionGithubEnterpriseConfigServiceDirectoryConfig)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionGithubEnterpriseConfigServiceDirectoryConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGithubEnterpriseConfigServiceDirectoryConfig or *ConnectionGithubEnterpriseConfigServiceDirectoryConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionGithubEnterpriseConfigServiceDirectoryConfig)
	if !ok {
		actualNotPointer, ok := a.(ConnectionGithubEnterpriseConfigServiceDirectoryConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGithubEnterpriseConfigServiceDirectoryConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Service, actual.Service, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("Service")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareConnectionGitlabConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionGitlabConfig)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionGitlabConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGitlabConfig or *ConnectionGitlabConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionGitlabConfig)
	if !ok {
		actualNotPointer, ok := a.(ConnectionGitlabConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGitlabConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.HostUri, actual.HostUri, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("HostUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WebhookSecretSecretVersion, actual.WebhookSecretSecretVersion, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("WebhookSecretSecretVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ReadAuthorizerCredential, actual.ReadAuthorizerCredential, dcl.DiffInfo{ObjectFunction: compareConnectionGitlabConfigReadAuthorizerCredentialNewStyle, EmptyObject: EmptyConnectionGitlabConfigReadAuthorizerCredential, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("ReadAuthorizerCredential")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AuthorizerCredential, actual.AuthorizerCredential, dcl.DiffInfo{ObjectFunction: compareConnectionGitlabConfigAuthorizerCredentialNewStyle, EmptyObject: EmptyConnectionGitlabConfigAuthorizerCredential, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("AuthorizerCredential")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceDirectoryConfig, actual.ServiceDirectoryConfig, dcl.DiffInfo{ObjectFunction: compareConnectionGitlabConfigServiceDirectoryConfigNewStyle, EmptyObject: EmptyConnectionGitlabConfigServiceDirectoryConfig, OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("ServiceDirectoryConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SslCa, actual.SslCa, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("SslCa")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServerVersion, actual.ServerVersion, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServerVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareConnectionGitlabConfigReadAuthorizerCredentialNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionGitlabConfigReadAuthorizerCredential)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionGitlabConfigReadAuthorizerCredential)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGitlabConfigReadAuthorizerCredential or *ConnectionGitlabConfigReadAuthorizerCredential", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionGitlabConfigReadAuthorizerCredential)
	if !ok {
		actualNotPointer, ok := a.(ConnectionGitlabConfigReadAuthorizerCredential)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGitlabConfigReadAuthorizerCredential", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.UserTokenSecretVersion, actual.UserTokenSecretVersion, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("UserTokenSecretVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Username, actual.Username, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Username")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareConnectionGitlabConfigAuthorizerCredentialNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionGitlabConfigAuthorizerCredential)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionGitlabConfigAuthorizerCredential)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGitlabConfigAuthorizerCredential or *ConnectionGitlabConfigAuthorizerCredential", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionGitlabConfigAuthorizerCredential)
	if !ok {
		actualNotPointer, ok := a.(ConnectionGitlabConfigAuthorizerCredential)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGitlabConfigAuthorizerCredential", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.UserTokenSecretVersion, actual.UserTokenSecretVersion, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("UserTokenSecretVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Username, actual.Username, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Username")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareConnectionGitlabConfigServiceDirectoryConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionGitlabConfigServiceDirectoryConfig)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionGitlabConfigServiceDirectoryConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGitlabConfigServiceDirectoryConfig or *ConnectionGitlabConfigServiceDirectoryConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionGitlabConfigServiceDirectoryConfig)
	if !ok {
		actualNotPointer, ok := a.(ConnectionGitlabConfigServiceDirectoryConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionGitlabConfigServiceDirectoryConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Service, actual.Service, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.TriggersOperation("updateConnectionUpdateConnectionOperation")}, fn.AddNest("Service")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareConnectionInstallationStateNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*ConnectionInstallationState)
	if !ok {
		desiredNotPointer, ok := d.(ConnectionInstallationState)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionInstallationState or *ConnectionInstallationState", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*ConnectionInstallationState)
	if !ok {
		actualNotPointer, ok := a.(ConnectionInstallationState)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a ConnectionInstallationState", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Stage, actual.Stage, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Stage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Message, actual.Message, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Message")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ActionUri, actual.ActionUri, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ActionUri")); len(ds) != 0 || err != nil {
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
func (r *Connection) urlNormalized() *Connection {
	normalized := dcl.Copy(*r).(Connection)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Etag = dcl.SelfLinkToName(r.Etag)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *Connection) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateConnection" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/connections/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Connection resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Connection) marshal(c *Client) ([]byte, error) {
	m, err := expandConnection(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Connection: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalConnection decodes JSON responses into the Connection resource schema.
func unmarshalConnection(b []byte, c *Client, res *Connection) (*Connection, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapConnection(m, c, res)
}

func unmarshalMapConnection(m map[string]interface{}, c *Client, res *Connection) (*Connection, error) {

	flattened := flattenConnection(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandConnection expands Connection into a JSON request object.
func expandConnection(c *Client, f *Connection) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v, err := expandConnectionGithubConfig(c, f.GithubConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GithubConfig into githubConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["githubConfig"] = v
	}
	if v, err := expandConnectionGithubEnterpriseConfig(c, f.GithubEnterpriseConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GithubEnterpriseConfig into githubEnterpriseConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["githubEnterpriseConfig"] = v
	}
	if v, err := expandConnectionGitlabConfig(c, f.GitlabConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GitlabConfig into gitlabConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gitlabConfig"] = v
	}
	if v := f.Disabled; dcl.ValueShouldBeSent(v) {
		m["disabled"] = v
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

	return m, nil
}

// flattenConnection flattens Connection from a JSON request object into the
// Connection type.
func flattenConnection(c *Client, i interface{}, res *Connection) *Connection {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Connection{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.GithubConfig = flattenConnectionGithubConfig(c, m["githubConfig"], res)
	resultRes.GithubEnterpriseConfig = flattenConnectionGithubEnterpriseConfig(c, m["githubEnterpriseConfig"], res)
	resultRes.GitlabConfig = flattenConnectionGitlabConfig(c, m["gitlabConfig"], res)
	resultRes.InstallationState = flattenConnectionInstallationState(c, m["installationState"], res)
	resultRes.Disabled = dcl.FlattenBool(m["disabled"])
	resultRes.Reconciling = dcl.FlattenBool(m["reconciling"])
	resultRes.Annotations = dcl.FlattenKeyValuePairs(m["annotations"])
	resultRes.Etag = dcl.FlattenString(m["etag"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandConnectionGithubConfigMap expands the contents of ConnectionGithubConfig into a JSON
// request object.
func expandConnectionGithubConfigMap(c *Client, f map[string]ConnectionGithubConfig, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionGithubConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionGithubConfigSlice expands the contents of ConnectionGithubConfig into a JSON
// request object.
func expandConnectionGithubConfigSlice(c *Client, f []ConnectionGithubConfig, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionGithubConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionGithubConfigMap flattens the contents of ConnectionGithubConfig from a JSON
// response object.
func flattenConnectionGithubConfigMap(c *Client, i interface{}, res *Connection) map[string]ConnectionGithubConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionGithubConfig{}
	}

	if len(a) == 0 {
		return map[string]ConnectionGithubConfig{}
	}

	items := make(map[string]ConnectionGithubConfig)
	for k, item := range a {
		items[k] = *flattenConnectionGithubConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionGithubConfigSlice flattens the contents of ConnectionGithubConfig from a JSON
// response object.
func flattenConnectionGithubConfigSlice(c *Client, i interface{}, res *Connection) []ConnectionGithubConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionGithubConfig{}
	}

	if len(a) == 0 {
		return []ConnectionGithubConfig{}
	}

	items := make([]ConnectionGithubConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionGithubConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionGithubConfig expands an instance of ConnectionGithubConfig into a JSON
// request object.
func expandConnectionGithubConfig(c *Client, f *ConnectionGithubConfig, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandConnectionGithubConfigAuthorizerCredential(c, f.AuthorizerCredential, res); err != nil {
		return nil, fmt.Errorf("error expanding AuthorizerCredential into authorizerCredential: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["authorizerCredential"] = v
	}
	if v := f.AppInstallationId; !dcl.IsEmptyValueIndirect(v) {
		m["appInstallationId"] = v
	}

	return m, nil
}

// flattenConnectionGithubConfig flattens an instance of ConnectionGithubConfig from a JSON
// response object.
func flattenConnectionGithubConfig(c *Client, i interface{}, res *Connection) *ConnectionGithubConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionGithubConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionGithubConfig
	}
	r.AuthorizerCredential = flattenConnectionGithubConfigAuthorizerCredential(c, m["authorizerCredential"], res)
	r.AppInstallationId = dcl.FlattenInteger(m["appInstallationId"])

	return r
}

// expandConnectionGithubConfigAuthorizerCredentialMap expands the contents of ConnectionGithubConfigAuthorizerCredential into a JSON
// request object.
func expandConnectionGithubConfigAuthorizerCredentialMap(c *Client, f map[string]ConnectionGithubConfigAuthorizerCredential, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionGithubConfigAuthorizerCredential(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionGithubConfigAuthorizerCredentialSlice expands the contents of ConnectionGithubConfigAuthorizerCredential into a JSON
// request object.
func expandConnectionGithubConfigAuthorizerCredentialSlice(c *Client, f []ConnectionGithubConfigAuthorizerCredential, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionGithubConfigAuthorizerCredential(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionGithubConfigAuthorizerCredentialMap flattens the contents of ConnectionGithubConfigAuthorizerCredential from a JSON
// response object.
func flattenConnectionGithubConfigAuthorizerCredentialMap(c *Client, i interface{}, res *Connection) map[string]ConnectionGithubConfigAuthorizerCredential {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionGithubConfigAuthorizerCredential{}
	}

	if len(a) == 0 {
		return map[string]ConnectionGithubConfigAuthorizerCredential{}
	}

	items := make(map[string]ConnectionGithubConfigAuthorizerCredential)
	for k, item := range a {
		items[k] = *flattenConnectionGithubConfigAuthorizerCredential(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionGithubConfigAuthorizerCredentialSlice flattens the contents of ConnectionGithubConfigAuthorizerCredential from a JSON
// response object.
func flattenConnectionGithubConfigAuthorizerCredentialSlice(c *Client, i interface{}, res *Connection) []ConnectionGithubConfigAuthorizerCredential {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionGithubConfigAuthorizerCredential{}
	}

	if len(a) == 0 {
		return []ConnectionGithubConfigAuthorizerCredential{}
	}

	items := make([]ConnectionGithubConfigAuthorizerCredential, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionGithubConfigAuthorizerCredential(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionGithubConfigAuthorizerCredential expands an instance of ConnectionGithubConfigAuthorizerCredential into a JSON
// request object.
func expandConnectionGithubConfigAuthorizerCredential(c *Client, f *ConnectionGithubConfigAuthorizerCredential, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.OAuthTokenSecretVersion; !dcl.IsEmptyValueIndirect(v) {
		m["oauthTokenSecretVersion"] = v
	}

	return m, nil
}

// flattenConnectionGithubConfigAuthorizerCredential flattens an instance of ConnectionGithubConfigAuthorizerCredential from a JSON
// response object.
func flattenConnectionGithubConfigAuthorizerCredential(c *Client, i interface{}, res *Connection) *ConnectionGithubConfigAuthorizerCredential {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionGithubConfigAuthorizerCredential{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionGithubConfigAuthorizerCredential
	}
	r.OAuthTokenSecretVersion = dcl.FlattenString(m["oauthTokenSecretVersion"])
	r.Username = dcl.FlattenString(m["username"])

	return r
}

// expandConnectionGithubEnterpriseConfigMap expands the contents of ConnectionGithubEnterpriseConfig into a JSON
// request object.
func expandConnectionGithubEnterpriseConfigMap(c *Client, f map[string]ConnectionGithubEnterpriseConfig, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionGithubEnterpriseConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionGithubEnterpriseConfigSlice expands the contents of ConnectionGithubEnterpriseConfig into a JSON
// request object.
func expandConnectionGithubEnterpriseConfigSlice(c *Client, f []ConnectionGithubEnterpriseConfig, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionGithubEnterpriseConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionGithubEnterpriseConfigMap flattens the contents of ConnectionGithubEnterpriseConfig from a JSON
// response object.
func flattenConnectionGithubEnterpriseConfigMap(c *Client, i interface{}, res *Connection) map[string]ConnectionGithubEnterpriseConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionGithubEnterpriseConfig{}
	}

	if len(a) == 0 {
		return map[string]ConnectionGithubEnterpriseConfig{}
	}

	items := make(map[string]ConnectionGithubEnterpriseConfig)
	for k, item := range a {
		items[k] = *flattenConnectionGithubEnterpriseConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionGithubEnterpriseConfigSlice flattens the contents of ConnectionGithubEnterpriseConfig from a JSON
// response object.
func flattenConnectionGithubEnterpriseConfigSlice(c *Client, i interface{}, res *Connection) []ConnectionGithubEnterpriseConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionGithubEnterpriseConfig{}
	}

	if len(a) == 0 {
		return []ConnectionGithubEnterpriseConfig{}
	}

	items := make([]ConnectionGithubEnterpriseConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionGithubEnterpriseConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionGithubEnterpriseConfig expands an instance of ConnectionGithubEnterpriseConfig into a JSON
// request object.
func expandConnectionGithubEnterpriseConfig(c *Client, f *ConnectionGithubEnterpriseConfig, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.HostUri; !dcl.IsEmptyValueIndirect(v) {
		m["hostUri"] = v
	}
	if v := f.AppId; !dcl.IsEmptyValueIndirect(v) {
		m["appId"] = v
	}
	if v := f.AppSlug; !dcl.IsEmptyValueIndirect(v) {
		m["appSlug"] = v
	}
	if v := f.PrivateKeySecretVersion; !dcl.IsEmptyValueIndirect(v) {
		m["privateKeySecretVersion"] = v
	}
	if v := f.WebhookSecretSecretVersion; !dcl.IsEmptyValueIndirect(v) {
		m["webhookSecretSecretVersion"] = v
	}
	if v := f.AppInstallationId; !dcl.IsEmptyValueIndirect(v) {
		m["appInstallationId"] = v
	}
	if v, err := expandConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, f.ServiceDirectoryConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ServiceDirectoryConfig into serviceDirectoryConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["serviceDirectoryConfig"] = v
	}
	if v := f.SslCa; !dcl.IsEmptyValueIndirect(v) {
		m["sslCa"] = v
	}

	return m, nil
}

// flattenConnectionGithubEnterpriseConfig flattens an instance of ConnectionGithubEnterpriseConfig from a JSON
// response object.
func flattenConnectionGithubEnterpriseConfig(c *Client, i interface{}, res *Connection) *ConnectionGithubEnterpriseConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionGithubEnterpriseConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionGithubEnterpriseConfig
	}
	r.HostUri = dcl.FlattenString(m["hostUri"])
	r.AppId = dcl.FlattenInteger(m["appId"])
	r.AppSlug = dcl.FlattenString(m["appSlug"])
	r.PrivateKeySecretVersion = dcl.FlattenString(m["privateKeySecretVersion"])
	r.WebhookSecretSecretVersion = dcl.FlattenString(m["webhookSecretSecretVersion"])
	r.AppInstallationId = dcl.FlattenInteger(m["appInstallationId"])
	r.ServiceDirectoryConfig = flattenConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, m["serviceDirectoryConfig"], res)
	r.SslCa = dcl.FlattenString(m["sslCa"])

	return r
}

// expandConnectionGithubEnterpriseConfigServiceDirectoryConfigMap expands the contents of ConnectionGithubEnterpriseConfigServiceDirectoryConfig into a JSON
// request object.
func expandConnectionGithubEnterpriseConfigServiceDirectoryConfigMap(c *Client, f map[string]ConnectionGithubEnterpriseConfigServiceDirectoryConfig, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionGithubEnterpriseConfigServiceDirectoryConfigSlice expands the contents of ConnectionGithubEnterpriseConfigServiceDirectoryConfig into a JSON
// request object.
func expandConnectionGithubEnterpriseConfigServiceDirectoryConfigSlice(c *Client, f []ConnectionGithubEnterpriseConfigServiceDirectoryConfig, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionGithubEnterpriseConfigServiceDirectoryConfigMap flattens the contents of ConnectionGithubEnterpriseConfigServiceDirectoryConfig from a JSON
// response object.
func flattenConnectionGithubEnterpriseConfigServiceDirectoryConfigMap(c *Client, i interface{}, res *Connection) map[string]ConnectionGithubEnterpriseConfigServiceDirectoryConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionGithubEnterpriseConfigServiceDirectoryConfig{}
	}

	if len(a) == 0 {
		return map[string]ConnectionGithubEnterpriseConfigServiceDirectoryConfig{}
	}

	items := make(map[string]ConnectionGithubEnterpriseConfigServiceDirectoryConfig)
	for k, item := range a {
		items[k] = *flattenConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionGithubEnterpriseConfigServiceDirectoryConfigSlice flattens the contents of ConnectionGithubEnterpriseConfigServiceDirectoryConfig from a JSON
// response object.
func flattenConnectionGithubEnterpriseConfigServiceDirectoryConfigSlice(c *Client, i interface{}, res *Connection) []ConnectionGithubEnterpriseConfigServiceDirectoryConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionGithubEnterpriseConfigServiceDirectoryConfig{}
	}

	if len(a) == 0 {
		return []ConnectionGithubEnterpriseConfigServiceDirectoryConfig{}
	}

	items := make([]ConnectionGithubEnterpriseConfigServiceDirectoryConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionGithubEnterpriseConfigServiceDirectoryConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionGithubEnterpriseConfigServiceDirectoryConfig expands an instance of ConnectionGithubEnterpriseConfigServiceDirectoryConfig into a JSON
// request object.
func expandConnectionGithubEnterpriseConfigServiceDirectoryConfig(c *Client, f *ConnectionGithubEnterpriseConfigServiceDirectoryConfig, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Service; !dcl.IsEmptyValueIndirect(v) {
		m["service"] = v
	}

	return m, nil
}

// flattenConnectionGithubEnterpriseConfigServiceDirectoryConfig flattens an instance of ConnectionGithubEnterpriseConfigServiceDirectoryConfig from a JSON
// response object.
func flattenConnectionGithubEnterpriseConfigServiceDirectoryConfig(c *Client, i interface{}, res *Connection) *ConnectionGithubEnterpriseConfigServiceDirectoryConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionGithubEnterpriseConfigServiceDirectoryConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionGithubEnterpriseConfigServiceDirectoryConfig
	}
	r.Service = dcl.FlattenString(m["service"])

	return r
}

// expandConnectionGitlabConfigMap expands the contents of ConnectionGitlabConfig into a JSON
// request object.
func expandConnectionGitlabConfigMap(c *Client, f map[string]ConnectionGitlabConfig, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionGitlabConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionGitlabConfigSlice expands the contents of ConnectionGitlabConfig into a JSON
// request object.
func expandConnectionGitlabConfigSlice(c *Client, f []ConnectionGitlabConfig, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionGitlabConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionGitlabConfigMap flattens the contents of ConnectionGitlabConfig from a JSON
// response object.
func flattenConnectionGitlabConfigMap(c *Client, i interface{}, res *Connection) map[string]ConnectionGitlabConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionGitlabConfig{}
	}

	if len(a) == 0 {
		return map[string]ConnectionGitlabConfig{}
	}

	items := make(map[string]ConnectionGitlabConfig)
	for k, item := range a {
		items[k] = *flattenConnectionGitlabConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionGitlabConfigSlice flattens the contents of ConnectionGitlabConfig from a JSON
// response object.
func flattenConnectionGitlabConfigSlice(c *Client, i interface{}, res *Connection) []ConnectionGitlabConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionGitlabConfig{}
	}

	if len(a) == 0 {
		return []ConnectionGitlabConfig{}
	}

	items := make([]ConnectionGitlabConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionGitlabConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionGitlabConfig expands an instance of ConnectionGitlabConfig into a JSON
// request object.
func expandConnectionGitlabConfig(c *Client, f *ConnectionGitlabConfig, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.HostUri; !dcl.IsEmptyValueIndirect(v) {
		m["hostUri"] = v
	}
	if v := f.WebhookSecretSecretVersion; !dcl.IsEmptyValueIndirect(v) {
		m["webhookSecretSecretVersion"] = v
	}
	if v, err := expandConnectionGitlabConfigReadAuthorizerCredential(c, f.ReadAuthorizerCredential, res); err != nil {
		return nil, fmt.Errorf("error expanding ReadAuthorizerCredential into readAuthorizerCredential: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["readAuthorizerCredential"] = v
	}
	if v, err := expandConnectionGitlabConfigAuthorizerCredential(c, f.AuthorizerCredential, res); err != nil {
		return nil, fmt.Errorf("error expanding AuthorizerCredential into authorizerCredential: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["authorizerCredential"] = v
	}
	if v, err := expandConnectionGitlabConfigServiceDirectoryConfig(c, f.ServiceDirectoryConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ServiceDirectoryConfig into serviceDirectoryConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["serviceDirectoryConfig"] = v
	}
	if v := f.SslCa; !dcl.IsEmptyValueIndirect(v) {
		m["sslCa"] = v
	}

	return m, nil
}

// flattenConnectionGitlabConfig flattens an instance of ConnectionGitlabConfig from a JSON
// response object.
func flattenConnectionGitlabConfig(c *Client, i interface{}, res *Connection) *ConnectionGitlabConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionGitlabConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionGitlabConfig
	}
	r.HostUri = dcl.FlattenString(m["hostUri"])
	r.WebhookSecretSecretVersion = dcl.FlattenString(m["webhookSecretSecretVersion"])
	r.ReadAuthorizerCredential = flattenConnectionGitlabConfigReadAuthorizerCredential(c, m["readAuthorizerCredential"], res)
	r.AuthorizerCredential = flattenConnectionGitlabConfigAuthorizerCredential(c, m["authorizerCredential"], res)
	r.ServiceDirectoryConfig = flattenConnectionGitlabConfigServiceDirectoryConfig(c, m["serviceDirectoryConfig"], res)
	r.SslCa = dcl.FlattenString(m["sslCa"])
	r.ServerVersion = dcl.FlattenString(m["serverVersion"])

	return r
}

// expandConnectionGitlabConfigReadAuthorizerCredentialMap expands the contents of ConnectionGitlabConfigReadAuthorizerCredential into a JSON
// request object.
func expandConnectionGitlabConfigReadAuthorizerCredentialMap(c *Client, f map[string]ConnectionGitlabConfigReadAuthorizerCredential, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionGitlabConfigReadAuthorizerCredential(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionGitlabConfigReadAuthorizerCredentialSlice expands the contents of ConnectionGitlabConfigReadAuthorizerCredential into a JSON
// request object.
func expandConnectionGitlabConfigReadAuthorizerCredentialSlice(c *Client, f []ConnectionGitlabConfigReadAuthorizerCredential, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionGitlabConfigReadAuthorizerCredential(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionGitlabConfigReadAuthorizerCredentialMap flattens the contents of ConnectionGitlabConfigReadAuthorizerCredential from a JSON
// response object.
func flattenConnectionGitlabConfigReadAuthorizerCredentialMap(c *Client, i interface{}, res *Connection) map[string]ConnectionGitlabConfigReadAuthorizerCredential {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionGitlabConfigReadAuthorizerCredential{}
	}

	if len(a) == 0 {
		return map[string]ConnectionGitlabConfigReadAuthorizerCredential{}
	}

	items := make(map[string]ConnectionGitlabConfigReadAuthorizerCredential)
	for k, item := range a {
		items[k] = *flattenConnectionGitlabConfigReadAuthorizerCredential(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionGitlabConfigReadAuthorizerCredentialSlice flattens the contents of ConnectionGitlabConfigReadAuthorizerCredential from a JSON
// response object.
func flattenConnectionGitlabConfigReadAuthorizerCredentialSlice(c *Client, i interface{}, res *Connection) []ConnectionGitlabConfigReadAuthorizerCredential {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionGitlabConfigReadAuthorizerCredential{}
	}

	if len(a) == 0 {
		return []ConnectionGitlabConfigReadAuthorizerCredential{}
	}

	items := make([]ConnectionGitlabConfigReadAuthorizerCredential, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionGitlabConfigReadAuthorizerCredential(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionGitlabConfigReadAuthorizerCredential expands an instance of ConnectionGitlabConfigReadAuthorizerCredential into a JSON
// request object.
func expandConnectionGitlabConfigReadAuthorizerCredential(c *Client, f *ConnectionGitlabConfigReadAuthorizerCredential, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.UserTokenSecretVersion; !dcl.IsEmptyValueIndirect(v) {
		m["userTokenSecretVersion"] = v
	}

	return m, nil
}

// flattenConnectionGitlabConfigReadAuthorizerCredential flattens an instance of ConnectionGitlabConfigReadAuthorizerCredential from a JSON
// response object.
func flattenConnectionGitlabConfigReadAuthorizerCredential(c *Client, i interface{}, res *Connection) *ConnectionGitlabConfigReadAuthorizerCredential {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionGitlabConfigReadAuthorizerCredential{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionGitlabConfigReadAuthorizerCredential
	}
	r.UserTokenSecretVersion = dcl.FlattenString(m["userTokenSecretVersion"])
	r.Username = dcl.FlattenString(m["username"])

	return r
}

// expandConnectionGitlabConfigAuthorizerCredentialMap expands the contents of ConnectionGitlabConfigAuthorizerCredential into a JSON
// request object.
func expandConnectionGitlabConfigAuthorizerCredentialMap(c *Client, f map[string]ConnectionGitlabConfigAuthorizerCredential, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionGitlabConfigAuthorizerCredential(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionGitlabConfigAuthorizerCredentialSlice expands the contents of ConnectionGitlabConfigAuthorizerCredential into a JSON
// request object.
func expandConnectionGitlabConfigAuthorizerCredentialSlice(c *Client, f []ConnectionGitlabConfigAuthorizerCredential, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionGitlabConfigAuthorizerCredential(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionGitlabConfigAuthorizerCredentialMap flattens the contents of ConnectionGitlabConfigAuthorizerCredential from a JSON
// response object.
func flattenConnectionGitlabConfigAuthorizerCredentialMap(c *Client, i interface{}, res *Connection) map[string]ConnectionGitlabConfigAuthorizerCredential {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionGitlabConfigAuthorizerCredential{}
	}

	if len(a) == 0 {
		return map[string]ConnectionGitlabConfigAuthorizerCredential{}
	}

	items := make(map[string]ConnectionGitlabConfigAuthorizerCredential)
	for k, item := range a {
		items[k] = *flattenConnectionGitlabConfigAuthorizerCredential(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionGitlabConfigAuthorizerCredentialSlice flattens the contents of ConnectionGitlabConfigAuthorizerCredential from a JSON
// response object.
func flattenConnectionGitlabConfigAuthorizerCredentialSlice(c *Client, i interface{}, res *Connection) []ConnectionGitlabConfigAuthorizerCredential {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionGitlabConfigAuthorizerCredential{}
	}

	if len(a) == 0 {
		return []ConnectionGitlabConfigAuthorizerCredential{}
	}

	items := make([]ConnectionGitlabConfigAuthorizerCredential, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionGitlabConfigAuthorizerCredential(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionGitlabConfigAuthorizerCredential expands an instance of ConnectionGitlabConfigAuthorizerCredential into a JSON
// request object.
func expandConnectionGitlabConfigAuthorizerCredential(c *Client, f *ConnectionGitlabConfigAuthorizerCredential, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.UserTokenSecretVersion; !dcl.IsEmptyValueIndirect(v) {
		m["userTokenSecretVersion"] = v
	}

	return m, nil
}

// flattenConnectionGitlabConfigAuthorizerCredential flattens an instance of ConnectionGitlabConfigAuthorizerCredential from a JSON
// response object.
func flattenConnectionGitlabConfigAuthorizerCredential(c *Client, i interface{}, res *Connection) *ConnectionGitlabConfigAuthorizerCredential {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionGitlabConfigAuthorizerCredential{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionGitlabConfigAuthorizerCredential
	}
	r.UserTokenSecretVersion = dcl.FlattenString(m["userTokenSecretVersion"])
	r.Username = dcl.FlattenString(m["username"])

	return r
}

// expandConnectionGitlabConfigServiceDirectoryConfigMap expands the contents of ConnectionGitlabConfigServiceDirectoryConfig into a JSON
// request object.
func expandConnectionGitlabConfigServiceDirectoryConfigMap(c *Client, f map[string]ConnectionGitlabConfigServiceDirectoryConfig, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionGitlabConfigServiceDirectoryConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionGitlabConfigServiceDirectoryConfigSlice expands the contents of ConnectionGitlabConfigServiceDirectoryConfig into a JSON
// request object.
func expandConnectionGitlabConfigServiceDirectoryConfigSlice(c *Client, f []ConnectionGitlabConfigServiceDirectoryConfig, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionGitlabConfigServiceDirectoryConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionGitlabConfigServiceDirectoryConfigMap flattens the contents of ConnectionGitlabConfigServiceDirectoryConfig from a JSON
// response object.
func flattenConnectionGitlabConfigServiceDirectoryConfigMap(c *Client, i interface{}, res *Connection) map[string]ConnectionGitlabConfigServiceDirectoryConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionGitlabConfigServiceDirectoryConfig{}
	}

	if len(a) == 0 {
		return map[string]ConnectionGitlabConfigServiceDirectoryConfig{}
	}

	items := make(map[string]ConnectionGitlabConfigServiceDirectoryConfig)
	for k, item := range a {
		items[k] = *flattenConnectionGitlabConfigServiceDirectoryConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionGitlabConfigServiceDirectoryConfigSlice flattens the contents of ConnectionGitlabConfigServiceDirectoryConfig from a JSON
// response object.
func flattenConnectionGitlabConfigServiceDirectoryConfigSlice(c *Client, i interface{}, res *Connection) []ConnectionGitlabConfigServiceDirectoryConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionGitlabConfigServiceDirectoryConfig{}
	}

	if len(a) == 0 {
		return []ConnectionGitlabConfigServiceDirectoryConfig{}
	}

	items := make([]ConnectionGitlabConfigServiceDirectoryConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionGitlabConfigServiceDirectoryConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionGitlabConfigServiceDirectoryConfig expands an instance of ConnectionGitlabConfigServiceDirectoryConfig into a JSON
// request object.
func expandConnectionGitlabConfigServiceDirectoryConfig(c *Client, f *ConnectionGitlabConfigServiceDirectoryConfig, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Service; !dcl.IsEmptyValueIndirect(v) {
		m["service"] = v
	}

	return m, nil
}

// flattenConnectionGitlabConfigServiceDirectoryConfig flattens an instance of ConnectionGitlabConfigServiceDirectoryConfig from a JSON
// response object.
func flattenConnectionGitlabConfigServiceDirectoryConfig(c *Client, i interface{}, res *Connection) *ConnectionGitlabConfigServiceDirectoryConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionGitlabConfigServiceDirectoryConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionGitlabConfigServiceDirectoryConfig
	}
	r.Service = dcl.FlattenString(m["service"])

	return r
}

// expandConnectionInstallationStateMap expands the contents of ConnectionInstallationState into a JSON
// request object.
func expandConnectionInstallationStateMap(c *Client, f map[string]ConnectionInstallationState, res *Connection) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandConnectionInstallationState(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandConnectionInstallationStateSlice expands the contents of ConnectionInstallationState into a JSON
// request object.
func expandConnectionInstallationStateSlice(c *Client, f []ConnectionInstallationState, res *Connection) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandConnectionInstallationState(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenConnectionInstallationStateMap flattens the contents of ConnectionInstallationState from a JSON
// response object.
func flattenConnectionInstallationStateMap(c *Client, i interface{}, res *Connection) map[string]ConnectionInstallationState {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionInstallationState{}
	}

	if len(a) == 0 {
		return map[string]ConnectionInstallationState{}
	}

	items := make(map[string]ConnectionInstallationState)
	for k, item := range a {
		items[k] = *flattenConnectionInstallationState(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenConnectionInstallationStateSlice flattens the contents of ConnectionInstallationState from a JSON
// response object.
func flattenConnectionInstallationStateSlice(c *Client, i interface{}, res *Connection) []ConnectionInstallationState {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionInstallationState{}
	}

	if len(a) == 0 {
		return []ConnectionInstallationState{}
	}

	items := make([]ConnectionInstallationState, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionInstallationState(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandConnectionInstallationState expands an instance of ConnectionInstallationState into a JSON
// request object.
func expandConnectionInstallationState(c *Client, f *ConnectionInstallationState, res *Connection) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenConnectionInstallationState flattens an instance of ConnectionInstallationState from a JSON
// response object.
func flattenConnectionInstallationState(c *Client, i interface{}, res *Connection) *ConnectionInstallationState {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &ConnectionInstallationState{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyConnectionInstallationState
	}
	r.Stage = flattenConnectionInstallationStateStageEnum(m["stage"])
	r.Message = dcl.FlattenString(m["message"])
	r.ActionUri = dcl.FlattenString(m["actionUri"])

	return r
}

// flattenConnectionInstallationStateStageEnumMap flattens the contents of ConnectionInstallationStateStageEnum from a JSON
// response object.
func flattenConnectionInstallationStateStageEnumMap(c *Client, i interface{}, res *Connection) map[string]ConnectionInstallationStateStageEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]ConnectionInstallationStateStageEnum{}
	}

	if len(a) == 0 {
		return map[string]ConnectionInstallationStateStageEnum{}
	}

	items := make(map[string]ConnectionInstallationStateStageEnum)
	for k, item := range a {
		items[k] = *flattenConnectionInstallationStateStageEnum(item.(interface{}))
	}

	return items
}

// flattenConnectionInstallationStateStageEnumSlice flattens the contents of ConnectionInstallationStateStageEnum from a JSON
// response object.
func flattenConnectionInstallationStateStageEnumSlice(c *Client, i interface{}, res *Connection) []ConnectionInstallationStateStageEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []ConnectionInstallationStateStageEnum{}
	}

	if len(a) == 0 {
		return []ConnectionInstallationStateStageEnum{}
	}

	items := make([]ConnectionInstallationStateStageEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenConnectionInstallationStateStageEnum(item.(interface{})))
	}

	return items
}

// flattenConnectionInstallationStateStageEnum asserts that an interface is a string, and returns a
// pointer to a *ConnectionInstallationStateStageEnum with the same value as that string.
func flattenConnectionInstallationStateStageEnum(i interface{}) *ConnectionInstallationStateStageEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return ConnectionInstallationStateStageEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Connection) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalConnection(b, c, r)
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

type connectionDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         connectionApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToConnectionDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]connectionDiff, error) {
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
	var diffs []connectionDiff
	// For each operation name, create a connectionDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := connectionDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToConnectionApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToConnectionApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (connectionApiOperation, error) {
	switch opName {

	case "updateConnectionUpdateConnectionOperation":
		return &updateConnectionUpdateConnectionOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractConnectionFields(r *Connection) error {
	vGithubConfig := r.GithubConfig
	if vGithubConfig == nil {
		// note: explicitly not the empty object.
		vGithubConfig = &ConnectionGithubConfig{}
	}
	if err := extractConnectionGithubConfigFields(r, vGithubConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGithubConfig) {
		r.GithubConfig = vGithubConfig
	}
	vGithubEnterpriseConfig := r.GithubEnterpriseConfig
	if vGithubEnterpriseConfig == nil {
		// note: explicitly not the empty object.
		vGithubEnterpriseConfig = &ConnectionGithubEnterpriseConfig{}
	}
	if err := extractConnectionGithubEnterpriseConfigFields(r, vGithubEnterpriseConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGithubEnterpriseConfig) {
		r.GithubEnterpriseConfig = vGithubEnterpriseConfig
	}
	vGitlabConfig := r.GitlabConfig
	if vGitlabConfig == nil {
		// note: explicitly not the empty object.
		vGitlabConfig = &ConnectionGitlabConfig{}
	}
	if err := extractConnectionGitlabConfigFields(r, vGitlabConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGitlabConfig) {
		r.GitlabConfig = vGitlabConfig
	}
	vInstallationState := r.InstallationState
	if vInstallationState == nil {
		// note: explicitly not the empty object.
		vInstallationState = &ConnectionInstallationState{}
	}
	if err := extractConnectionInstallationStateFields(r, vInstallationState); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vInstallationState) {
		r.InstallationState = vInstallationState
	}
	return nil
}
func extractConnectionGithubConfigFields(r *Connection, o *ConnectionGithubConfig) error {
	vAuthorizerCredential := o.AuthorizerCredential
	if vAuthorizerCredential == nil {
		// note: explicitly not the empty object.
		vAuthorizerCredential = &ConnectionGithubConfigAuthorizerCredential{}
	}
	if err := extractConnectionGithubConfigAuthorizerCredentialFields(r, vAuthorizerCredential); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorizerCredential) {
		o.AuthorizerCredential = vAuthorizerCredential
	}
	return nil
}
func extractConnectionGithubConfigAuthorizerCredentialFields(r *Connection, o *ConnectionGithubConfigAuthorizerCredential) error {
	return nil
}
func extractConnectionGithubEnterpriseConfigFields(r *Connection, o *ConnectionGithubEnterpriseConfig) error {
	vServiceDirectoryConfig := o.ServiceDirectoryConfig
	if vServiceDirectoryConfig == nil {
		// note: explicitly not the empty object.
		vServiceDirectoryConfig = &ConnectionGithubEnterpriseConfigServiceDirectoryConfig{}
	}
	if err := extractConnectionGithubEnterpriseConfigServiceDirectoryConfigFields(r, vServiceDirectoryConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vServiceDirectoryConfig) {
		o.ServiceDirectoryConfig = vServiceDirectoryConfig
	}
	return nil
}
func extractConnectionGithubEnterpriseConfigServiceDirectoryConfigFields(r *Connection, o *ConnectionGithubEnterpriseConfigServiceDirectoryConfig) error {
	return nil
}
func extractConnectionGitlabConfigFields(r *Connection, o *ConnectionGitlabConfig) error {
	vReadAuthorizerCredential := o.ReadAuthorizerCredential
	if vReadAuthorizerCredential == nil {
		// note: explicitly not the empty object.
		vReadAuthorizerCredential = &ConnectionGitlabConfigReadAuthorizerCredential{}
	}
	if err := extractConnectionGitlabConfigReadAuthorizerCredentialFields(r, vReadAuthorizerCredential); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vReadAuthorizerCredential) {
		o.ReadAuthorizerCredential = vReadAuthorizerCredential
	}
	vAuthorizerCredential := o.AuthorizerCredential
	if vAuthorizerCredential == nil {
		// note: explicitly not the empty object.
		vAuthorizerCredential = &ConnectionGitlabConfigAuthorizerCredential{}
	}
	if err := extractConnectionGitlabConfigAuthorizerCredentialFields(r, vAuthorizerCredential); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorizerCredential) {
		o.AuthorizerCredential = vAuthorizerCredential
	}
	vServiceDirectoryConfig := o.ServiceDirectoryConfig
	if vServiceDirectoryConfig == nil {
		// note: explicitly not the empty object.
		vServiceDirectoryConfig = &ConnectionGitlabConfigServiceDirectoryConfig{}
	}
	if err := extractConnectionGitlabConfigServiceDirectoryConfigFields(r, vServiceDirectoryConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vServiceDirectoryConfig) {
		o.ServiceDirectoryConfig = vServiceDirectoryConfig
	}
	return nil
}
func extractConnectionGitlabConfigReadAuthorizerCredentialFields(r *Connection, o *ConnectionGitlabConfigReadAuthorizerCredential) error {
	return nil
}
func extractConnectionGitlabConfigAuthorizerCredentialFields(r *Connection, o *ConnectionGitlabConfigAuthorizerCredential) error {
	return nil
}
func extractConnectionGitlabConfigServiceDirectoryConfigFields(r *Connection, o *ConnectionGitlabConfigServiceDirectoryConfig) error {
	return nil
}
func extractConnectionInstallationStateFields(r *Connection, o *ConnectionInstallationState) error {
	return nil
}

func postReadExtractConnectionFields(r *Connection) error {
	vGithubConfig := r.GithubConfig
	if vGithubConfig == nil {
		// note: explicitly not the empty object.
		vGithubConfig = &ConnectionGithubConfig{}
	}
	if err := postReadExtractConnectionGithubConfigFields(r, vGithubConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGithubConfig) {
		r.GithubConfig = vGithubConfig
	}
	vGithubEnterpriseConfig := r.GithubEnterpriseConfig
	if vGithubEnterpriseConfig == nil {
		// note: explicitly not the empty object.
		vGithubEnterpriseConfig = &ConnectionGithubEnterpriseConfig{}
	}
	if err := postReadExtractConnectionGithubEnterpriseConfigFields(r, vGithubEnterpriseConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGithubEnterpriseConfig) {
		r.GithubEnterpriseConfig = vGithubEnterpriseConfig
	}
	vGitlabConfig := r.GitlabConfig
	if vGitlabConfig == nil {
		// note: explicitly not the empty object.
		vGitlabConfig = &ConnectionGitlabConfig{}
	}
	if err := postReadExtractConnectionGitlabConfigFields(r, vGitlabConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGitlabConfig) {
		r.GitlabConfig = vGitlabConfig
	}
	vInstallationState := r.InstallationState
	if vInstallationState == nil {
		// note: explicitly not the empty object.
		vInstallationState = &ConnectionInstallationState{}
	}
	if err := postReadExtractConnectionInstallationStateFields(r, vInstallationState); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vInstallationState) {
		r.InstallationState = vInstallationState
	}
	return nil
}
func postReadExtractConnectionGithubConfigFields(r *Connection, o *ConnectionGithubConfig) error {
	vAuthorizerCredential := o.AuthorizerCredential
	if vAuthorizerCredential == nil {
		// note: explicitly not the empty object.
		vAuthorizerCredential = &ConnectionGithubConfigAuthorizerCredential{}
	}
	if err := extractConnectionGithubConfigAuthorizerCredentialFields(r, vAuthorizerCredential); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorizerCredential) {
		o.AuthorizerCredential = vAuthorizerCredential
	}
	return nil
}
func postReadExtractConnectionGithubConfigAuthorizerCredentialFields(r *Connection, o *ConnectionGithubConfigAuthorizerCredential) error {
	return nil
}
func postReadExtractConnectionGithubEnterpriseConfigFields(r *Connection, o *ConnectionGithubEnterpriseConfig) error {
	vServiceDirectoryConfig := o.ServiceDirectoryConfig
	if vServiceDirectoryConfig == nil {
		// note: explicitly not the empty object.
		vServiceDirectoryConfig = &ConnectionGithubEnterpriseConfigServiceDirectoryConfig{}
	}
	if err := extractConnectionGithubEnterpriseConfigServiceDirectoryConfigFields(r, vServiceDirectoryConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vServiceDirectoryConfig) {
		o.ServiceDirectoryConfig = vServiceDirectoryConfig
	}
	return nil
}
func postReadExtractConnectionGithubEnterpriseConfigServiceDirectoryConfigFields(r *Connection, o *ConnectionGithubEnterpriseConfigServiceDirectoryConfig) error {
	return nil
}
func postReadExtractConnectionGitlabConfigFields(r *Connection, o *ConnectionGitlabConfig) error {
	vReadAuthorizerCredential := o.ReadAuthorizerCredential
	if vReadAuthorizerCredential == nil {
		// note: explicitly not the empty object.
		vReadAuthorizerCredential = &ConnectionGitlabConfigReadAuthorizerCredential{}
	}
	if err := extractConnectionGitlabConfigReadAuthorizerCredentialFields(r, vReadAuthorizerCredential); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vReadAuthorizerCredential) {
		o.ReadAuthorizerCredential = vReadAuthorizerCredential
	}
	vAuthorizerCredential := o.AuthorizerCredential
	if vAuthorizerCredential == nil {
		// note: explicitly not the empty object.
		vAuthorizerCredential = &ConnectionGitlabConfigAuthorizerCredential{}
	}
	if err := extractConnectionGitlabConfigAuthorizerCredentialFields(r, vAuthorizerCredential); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorizerCredential) {
		o.AuthorizerCredential = vAuthorizerCredential
	}
	vServiceDirectoryConfig := o.ServiceDirectoryConfig
	if vServiceDirectoryConfig == nil {
		// note: explicitly not the empty object.
		vServiceDirectoryConfig = &ConnectionGitlabConfigServiceDirectoryConfig{}
	}
	if err := extractConnectionGitlabConfigServiceDirectoryConfigFields(r, vServiceDirectoryConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vServiceDirectoryConfig) {
		o.ServiceDirectoryConfig = vServiceDirectoryConfig
	}
	return nil
}
func postReadExtractConnectionGitlabConfigReadAuthorizerCredentialFields(r *Connection, o *ConnectionGitlabConfigReadAuthorizerCredential) error {
	return nil
}
func postReadExtractConnectionGitlabConfigAuthorizerCredentialFields(r *Connection, o *ConnectionGitlabConfigAuthorizerCredential) error {
	return nil
}
func postReadExtractConnectionGitlabConfigServiceDirectoryConfigFields(r *Connection, o *ConnectionGitlabConfigServiceDirectoryConfig) error {
	return nil
}
func postReadExtractConnectionInstallationStateFields(r *Connection, o *ConnectionInstallationState) error {
	return nil
}
