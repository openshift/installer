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
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Connection struct {
	Name                   *string                           `json:"name"`
	CreateTime             *string                           `json:"createTime"`
	UpdateTime             *string                           `json:"updateTime"`
	GithubConfig           *ConnectionGithubConfig           `json:"githubConfig"`
	GithubEnterpriseConfig *ConnectionGithubEnterpriseConfig `json:"githubEnterpriseConfig"`
	GitlabConfig           *ConnectionGitlabConfig           `json:"gitlabConfig"`
	InstallationState      *ConnectionInstallationState      `json:"installationState"`
	Disabled               *bool                             `json:"disabled"`
	Reconciling            *bool                             `json:"reconciling"`
	Annotations            map[string]string                 `json:"annotations"`
	Etag                   *string                           `json:"etag"`
	Project                *string                           `json:"project"`
	Location               *string                           `json:"location"`
}

func (r *Connection) String() string {
	return dcl.SprintResource(r)
}

// The enum ConnectionInstallationStateStageEnum.
type ConnectionInstallationStateStageEnum string

// ConnectionInstallationStateStageEnumRef returns a *ConnectionInstallationStateStageEnum with the value of string s
// If the empty string is provided, nil is returned.
func ConnectionInstallationStateStageEnumRef(s string) *ConnectionInstallationStateStageEnum {
	v := ConnectionInstallationStateStageEnum(s)
	return &v
}

func (v ConnectionInstallationStateStageEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STAGE_UNSPECIFIED", "PENDING_CREATE_APP", "PENDING_USER_OAUTH", "PENDING_INSTALL_APP", "COMPLETE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ConnectionInstallationStateStageEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type ConnectionGithubConfig struct {
	empty                bool                                        `json:"-"`
	AuthorizerCredential *ConnectionGithubConfigAuthorizerCredential `json:"authorizerCredential"`
	AppInstallationId    *int64                                      `json:"appInstallationId"`
}

type jsonConnectionGithubConfig ConnectionGithubConfig

func (r *ConnectionGithubConfig) UnmarshalJSON(data []byte) error {
	var res jsonConnectionGithubConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionGithubConfig
	} else {

		r.AuthorizerCredential = res.AuthorizerCredential

		r.AppInstallationId = res.AppInstallationId

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionGithubConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionGithubConfig *ConnectionGithubConfig = &ConnectionGithubConfig{empty: true}

func (r *ConnectionGithubConfig) Empty() bool {
	return r.empty
}

func (r *ConnectionGithubConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionGithubConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ConnectionGithubConfigAuthorizerCredential struct {
	empty                   bool    `json:"-"`
	OAuthTokenSecretVersion *string `json:"oauthTokenSecretVersion"`
	Username                *string `json:"username"`
}

type jsonConnectionGithubConfigAuthorizerCredential ConnectionGithubConfigAuthorizerCredential

func (r *ConnectionGithubConfigAuthorizerCredential) UnmarshalJSON(data []byte) error {
	var res jsonConnectionGithubConfigAuthorizerCredential
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionGithubConfigAuthorizerCredential
	} else {

		r.OAuthTokenSecretVersion = res.OAuthTokenSecretVersion

		r.Username = res.Username

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionGithubConfigAuthorizerCredential is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionGithubConfigAuthorizerCredential *ConnectionGithubConfigAuthorizerCredential = &ConnectionGithubConfigAuthorizerCredential{empty: true}

func (r *ConnectionGithubConfigAuthorizerCredential) Empty() bool {
	return r.empty
}

func (r *ConnectionGithubConfigAuthorizerCredential) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionGithubConfigAuthorizerCredential) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ConnectionGithubEnterpriseConfig struct {
	empty                      bool                                                    `json:"-"`
	HostUri                    *string                                                 `json:"hostUri"`
	AppId                      *int64                                                  `json:"appId"`
	AppSlug                    *string                                                 `json:"appSlug"`
	PrivateKeySecretVersion    *string                                                 `json:"privateKeySecretVersion"`
	WebhookSecretSecretVersion *string                                                 `json:"webhookSecretSecretVersion"`
	AppInstallationId          *int64                                                  `json:"appInstallationId"`
	ServiceDirectoryConfig     *ConnectionGithubEnterpriseConfigServiceDirectoryConfig `json:"serviceDirectoryConfig"`
	SslCa                      *string                                                 `json:"sslCa"`
}

type jsonConnectionGithubEnterpriseConfig ConnectionGithubEnterpriseConfig

func (r *ConnectionGithubEnterpriseConfig) UnmarshalJSON(data []byte) error {
	var res jsonConnectionGithubEnterpriseConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionGithubEnterpriseConfig
	} else {

		r.HostUri = res.HostUri

		r.AppId = res.AppId

		r.AppSlug = res.AppSlug

		r.PrivateKeySecretVersion = res.PrivateKeySecretVersion

		r.WebhookSecretSecretVersion = res.WebhookSecretSecretVersion

		r.AppInstallationId = res.AppInstallationId

		r.ServiceDirectoryConfig = res.ServiceDirectoryConfig

		r.SslCa = res.SslCa

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionGithubEnterpriseConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionGithubEnterpriseConfig *ConnectionGithubEnterpriseConfig = &ConnectionGithubEnterpriseConfig{empty: true}

func (r *ConnectionGithubEnterpriseConfig) Empty() bool {
	return r.empty
}

func (r *ConnectionGithubEnterpriseConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionGithubEnterpriseConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ConnectionGithubEnterpriseConfigServiceDirectoryConfig struct {
	empty   bool    `json:"-"`
	Service *string `json:"service"`
}

type jsonConnectionGithubEnterpriseConfigServiceDirectoryConfig ConnectionGithubEnterpriseConfigServiceDirectoryConfig

func (r *ConnectionGithubEnterpriseConfigServiceDirectoryConfig) UnmarshalJSON(data []byte) error {
	var res jsonConnectionGithubEnterpriseConfigServiceDirectoryConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionGithubEnterpriseConfigServiceDirectoryConfig
	} else {

		r.Service = res.Service

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionGithubEnterpriseConfigServiceDirectoryConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionGithubEnterpriseConfigServiceDirectoryConfig *ConnectionGithubEnterpriseConfigServiceDirectoryConfig = &ConnectionGithubEnterpriseConfigServiceDirectoryConfig{empty: true}

func (r *ConnectionGithubEnterpriseConfigServiceDirectoryConfig) Empty() bool {
	return r.empty
}

func (r *ConnectionGithubEnterpriseConfigServiceDirectoryConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionGithubEnterpriseConfigServiceDirectoryConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ConnectionGitlabConfig struct {
	empty                      bool                                            `json:"-"`
	HostUri                    *string                                         `json:"hostUri"`
	WebhookSecretSecretVersion *string                                         `json:"webhookSecretSecretVersion"`
	ReadAuthorizerCredential   *ConnectionGitlabConfigReadAuthorizerCredential `json:"readAuthorizerCredential"`
	AuthorizerCredential       *ConnectionGitlabConfigAuthorizerCredential     `json:"authorizerCredential"`
	ServiceDirectoryConfig     *ConnectionGitlabConfigServiceDirectoryConfig   `json:"serviceDirectoryConfig"`
	SslCa                      *string                                         `json:"sslCa"`
	ServerVersion              *string                                         `json:"serverVersion"`
}

type jsonConnectionGitlabConfig ConnectionGitlabConfig

func (r *ConnectionGitlabConfig) UnmarshalJSON(data []byte) error {
	var res jsonConnectionGitlabConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionGitlabConfig
	} else {

		r.HostUri = res.HostUri

		r.WebhookSecretSecretVersion = res.WebhookSecretSecretVersion

		r.ReadAuthorizerCredential = res.ReadAuthorizerCredential

		r.AuthorizerCredential = res.AuthorizerCredential

		r.ServiceDirectoryConfig = res.ServiceDirectoryConfig

		r.SslCa = res.SslCa

		r.ServerVersion = res.ServerVersion

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionGitlabConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionGitlabConfig *ConnectionGitlabConfig = &ConnectionGitlabConfig{empty: true}

func (r *ConnectionGitlabConfig) Empty() bool {
	return r.empty
}

func (r *ConnectionGitlabConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionGitlabConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ConnectionGitlabConfigReadAuthorizerCredential struct {
	empty                  bool    `json:"-"`
	UserTokenSecretVersion *string `json:"userTokenSecretVersion"`
	Username               *string `json:"username"`
}

type jsonConnectionGitlabConfigReadAuthorizerCredential ConnectionGitlabConfigReadAuthorizerCredential

func (r *ConnectionGitlabConfigReadAuthorizerCredential) UnmarshalJSON(data []byte) error {
	var res jsonConnectionGitlabConfigReadAuthorizerCredential
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionGitlabConfigReadAuthorizerCredential
	} else {

		r.UserTokenSecretVersion = res.UserTokenSecretVersion

		r.Username = res.Username

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionGitlabConfigReadAuthorizerCredential is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionGitlabConfigReadAuthorizerCredential *ConnectionGitlabConfigReadAuthorizerCredential = &ConnectionGitlabConfigReadAuthorizerCredential{empty: true}

func (r *ConnectionGitlabConfigReadAuthorizerCredential) Empty() bool {
	return r.empty
}

func (r *ConnectionGitlabConfigReadAuthorizerCredential) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionGitlabConfigReadAuthorizerCredential) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ConnectionGitlabConfigAuthorizerCredential struct {
	empty                  bool    `json:"-"`
	UserTokenSecretVersion *string `json:"userTokenSecretVersion"`
	Username               *string `json:"username"`
}

type jsonConnectionGitlabConfigAuthorizerCredential ConnectionGitlabConfigAuthorizerCredential

func (r *ConnectionGitlabConfigAuthorizerCredential) UnmarshalJSON(data []byte) error {
	var res jsonConnectionGitlabConfigAuthorizerCredential
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionGitlabConfigAuthorizerCredential
	} else {

		r.UserTokenSecretVersion = res.UserTokenSecretVersion

		r.Username = res.Username

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionGitlabConfigAuthorizerCredential is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionGitlabConfigAuthorizerCredential *ConnectionGitlabConfigAuthorizerCredential = &ConnectionGitlabConfigAuthorizerCredential{empty: true}

func (r *ConnectionGitlabConfigAuthorizerCredential) Empty() bool {
	return r.empty
}

func (r *ConnectionGitlabConfigAuthorizerCredential) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionGitlabConfigAuthorizerCredential) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ConnectionGitlabConfigServiceDirectoryConfig struct {
	empty   bool    `json:"-"`
	Service *string `json:"service"`
}

type jsonConnectionGitlabConfigServiceDirectoryConfig ConnectionGitlabConfigServiceDirectoryConfig

func (r *ConnectionGitlabConfigServiceDirectoryConfig) UnmarshalJSON(data []byte) error {
	var res jsonConnectionGitlabConfigServiceDirectoryConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionGitlabConfigServiceDirectoryConfig
	} else {

		r.Service = res.Service

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionGitlabConfigServiceDirectoryConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionGitlabConfigServiceDirectoryConfig *ConnectionGitlabConfigServiceDirectoryConfig = &ConnectionGitlabConfigServiceDirectoryConfig{empty: true}

func (r *ConnectionGitlabConfigServiceDirectoryConfig) Empty() bool {
	return r.empty
}

func (r *ConnectionGitlabConfigServiceDirectoryConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionGitlabConfigServiceDirectoryConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ConnectionInstallationState struct {
	empty     bool                                  `json:"-"`
	Stage     *ConnectionInstallationStateStageEnum `json:"stage"`
	Message   *string                               `json:"message"`
	ActionUri *string                               `json:"actionUri"`
}

type jsonConnectionInstallationState ConnectionInstallationState

func (r *ConnectionInstallationState) UnmarshalJSON(data []byte) error {
	var res jsonConnectionInstallationState
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyConnectionInstallationState
	} else {

		r.Stage = res.Stage

		r.Message = res.Message

		r.ActionUri = res.ActionUri

	}
	return nil
}

// This object is used to assert a desired state where this ConnectionInstallationState is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyConnectionInstallationState *ConnectionInstallationState = &ConnectionInstallationState{empty: true}

func (r *ConnectionInstallationState) Empty() bool {
	return r.empty
}

func (r *ConnectionInstallationState) String() string {
	return dcl.SprintResource(r)
}

func (r *ConnectionInstallationState) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Connection) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "cloudbuildv2",
		Type:    "Connection",
		Version: "cloudbuildv2",
	}
}

func (r *Connection) ID() (string, error) {
	if err := extractConnectionFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                     dcl.ValueOrEmptyString(nr.Name),
		"create_time":              dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":              dcl.ValueOrEmptyString(nr.UpdateTime),
		"github_config":            dcl.ValueOrEmptyString(nr.GithubConfig),
		"github_enterprise_config": dcl.ValueOrEmptyString(nr.GithubEnterpriseConfig),
		"gitlab_config":            dcl.ValueOrEmptyString(nr.GitlabConfig),
		"installation_state":       dcl.ValueOrEmptyString(nr.InstallationState),
		"disabled":                 dcl.ValueOrEmptyString(nr.Disabled),
		"reconciling":              dcl.ValueOrEmptyString(nr.Reconciling),
		"annotations":              dcl.ValueOrEmptyString(nr.Annotations),
		"etag":                     dcl.ValueOrEmptyString(nr.Etag),
		"project":                  dcl.ValueOrEmptyString(nr.Project),
		"location":                 dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/connections/{{name}}", params), nil
}

const ConnectionMaxPage = -1

type ConnectionList struct {
	Items []*Connection

	nextToken string

	pageSize int32

	resource *Connection
}

func (l *ConnectionList) HasNext() bool {
	return l.nextToken != ""
}

func (l *ConnectionList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listConnection(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListConnection(ctx context.Context, project, location string) (*ConnectionList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListConnectionWithMaxResults(ctx, project, location, ConnectionMaxPage)

}

func (c *Client) ListConnectionWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*ConnectionList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Connection{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listConnection(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &ConnectionList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetConnection(ctx context.Context, r *Connection) (*Connection, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractConnectionFields(r)

	b, err := c.getConnectionRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalConnection(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeConnectionNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractConnectionFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteConnection(ctx context.Context, r *Connection) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Connection resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Connection...")
	deleteOp := deleteConnectionOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllConnection deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllConnection(ctx context.Context, project, location string, filter func(*Connection) bool) error {
	listObj, err := c.ListConnection(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllConnection(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllConnection(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyConnection(ctx context.Context, rawDesired *Connection, opts ...dcl.ApplyOption) (*Connection, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Connection
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyConnectionHelper(c, ctx, rawDesired, opts...)
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

func applyConnectionHelper(c *Client, ctx context.Context, rawDesired *Connection, opts ...dcl.ApplyOption) (*Connection, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyConnection...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractConnectionFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.connectionDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToConnectionDiffs(c.Config, fieldDiffs, opts)
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
	var ops []connectionApiOperation
	if create {
		ops = append(ops, &createConnectionOperation{})
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
	return applyConnectionDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyConnectionDiff(c *Client, ctx context.Context, desired *Connection, rawDesired *Connection, ops []connectionApiOperation, opts ...dcl.ApplyOption) (*Connection, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetConnection(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createConnectionOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapConnection(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeConnectionNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeConnectionNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeConnectionDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractConnectionFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractConnectionFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffConnection(c, newDesired, newState)
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
