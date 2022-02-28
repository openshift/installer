package tfe

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
)

// Compile-time proof of interface implementation.
var _ Workspaces = (*workspaces)(nil)

// Workspaces describes all the workspace related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://www.terraform.io/docs/enterprise/api/workspaces.html
type Workspaces interface {
	// List all the workspaces within an organization.
	List(ctx context.Context, organization string, options WorkspaceListOptions) (*WorkspaceList, error)

	// Create is used to create a new workspace.
	Create(ctx context.Context, organization string, options WorkspaceCreateOptions) (*Workspace, error)

	// Read a workspace by its name.
	Read(ctx context.Context, organization string, workspace string) (*Workspace, error)

	// Readme gets the readme of a workspace by its ID.
	Readme(ctx context.Context, workspaceID string) (io.Reader, error)

	// ReadByID reads a workspace by its ID.
	ReadByID(ctx context.Context, workspaceID string) (*Workspace, error)

	// Update settings of an existing workspace.
	Update(ctx context.Context, organization string, workspace string, options WorkspaceUpdateOptions) (*Workspace, error)

	// UpdateByID updates the settings of an existing workspace.
	UpdateByID(ctx context.Context, workspaceID string, options WorkspaceUpdateOptions) (*Workspace, error)

	// Delete a workspace by its name.
	Delete(ctx context.Context, organization string, workspace string) error

	// DeleteByID deletes a workspace by its ID.
	DeleteByID(ctx context.Context, workspaceID string) error

	// RemoveVCSConnection from a workspace.
	RemoveVCSConnection(ctx context.Context, organization, workspace string) (*Workspace, error)

	// RemoveVCSConnectionByID removes a VCS connection from a workspace.
	RemoveVCSConnectionByID(ctx context.Context, workspaceID string) (*Workspace, error)

	// Lock a workspace by its ID.
	Lock(ctx context.Context, workspaceID string, options WorkspaceLockOptions) (*Workspace, error)

	// Unlock a workspace by its ID.
	Unlock(ctx context.Context, workspaceID string) (*Workspace, error)

	// ForceUnlock a workspace by its ID.
	ForceUnlock(ctx context.Context, workspaceID string) (*Workspace, error)

	// AssignSSHKey to a workspace.
	AssignSSHKey(ctx context.Context, workspaceID string, options WorkspaceAssignSSHKeyOptions) (*Workspace, error)

	// UnassignSSHKey from a workspace.
	UnassignSSHKey(ctx context.Context, workspaceID string) (*Workspace, error)

	// RemoteStateConsumers reads the remote state consumers for a workspace.
	RemoteStateConsumers(ctx context.Context, workspaceID string) (*WorkspaceList, error)

	// AddRemoteStateConsumers adds remote state consumers to a workspace.
	AddRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceAddRemoteStateConsumersOptions) error

	// RemoveRemoteStateConsumers removes remote state consumers from a workspace.
	RemoveRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceRemoveRemoteStateConsumersOptions) error

	// UpdateRemoteStateConsumers updates all the remote state consumers for a workspace
	// to match the workspaces in the update options.
	UpdateRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceUpdateRemoteStateConsumersOptions) error
}

// workspaces implements Workspaces.
type workspaces struct {
	client *Client
}

// WorkspaceList represents a list of workspaces.
type WorkspaceList struct {
	*Pagination
	Items []*Workspace
}

// Workspace represents a Terraform Enterprise workspace.
type Workspace struct {
	ID                   string                `jsonapi:"primary,workspaces"`
	Actions              *WorkspaceActions     `jsonapi:"attr,actions"`
	AgentPoolID          string                `jsonapi:"attr,agent-pool-id"`
	AllowDestroyPlan     bool                  `jsonapi:"attr,allow-destroy-plan"`
	AutoApply            bool                  `jsonapi:"attr,auto-apply"`
	CanQueueDestroyPlan  bool                  `jsonapi:"attr,can-queue-destroy-plan"`
	CreatedAt            time.Time             `jsonapi:"attr,created-at,iso8601"`
	Description          string                `jsonapi:"attr,description"`
	Environment          string                `jsonapi:"attr,environment"`
	ExecutionMode        string                `jsonapi:"attr,execution-mode"`
	FileTriggersEnabled  bool                  `jsonapi:"attr,file-triggers-enabled"`
	GlobalRemoteState    bool                  `jsonapi:"attr,global-remote-state"`
	Locked               bool                  `jsonapi:"attr,locked"`
	MigrationEnvironment string                `jsonapi:"attr,migration-environment"`
	Name                 string                `jsonapi:"attr,name"`
	Operations           bool                  `jsonapi:"attr,operations"`
	Permissions          *WorkspacePermissions `jsonapi:"attr,permissions"`
	QueueAllRuns         bool                  `jsonapi:"attr,queue-all-runs"`
	SpeculativeEnabled   bool                  `jsonapi:"attr,speculative-enabled"`
	TerraformVersion     string                `jsonapi:"attr,terraform-version"`
	TriggerPrefixes      []string              `jsonapi:"attr,trigger-prefixes"`
	VCSRepo              *VCSRepo              `jsonapi:"attr,vcs-repo"`
	WorkingDirectory     string                `jsonapi:"attr,working-directory"`
	UpdatedAt            time.Time             `jsonapi:"attr,updated-at,iso8601"`
	ResourceCount        int                   `jsonapi:"attr,resource-count"`
	ApplyDurationAverage time.Duration         `jsonapi:"attr,apply-duration-average"`
	PlanDurationAverage  time.Duration         `jsonapi:"attr,plan-duration-average"`
	PolicyCheckFailures  int                   `jsonapi:"attr,policy-check-failures"`
	RunFailures          int                   `jsonapi:"attr,run-failures"`
	RunsCount            int                   `jsonapi:"attr,workspace-kpis-runs-count"`

	// Relations
	AgentPool    *AgentPool    `jsonapi:"relation,agent-pool"`
	CurrentRun   *Run          `jsonapi:"relation,current-run"`
	Organization *Organization `jsonapi:"relation,organization"`
	SSHKey       *SSHKey       `jsonapi:"relation,ssh-key"`
}

// workspaceWithReadme is the same as a workspace but it has a readme.
type workspaceWithReadme struct {
	ID     string           `jsonapi:"primary,workspaces"`
	Readme *workspaceReadme `jsonapi:"relation,readme"`
}

// workspaceReadme contains the readme of the workspace.
type workspaceReadme struct {
	ID          string `jsonapi:"primary,workspace-readme"`
	RawMarkdown string `jsonapi:"attr,raw-markdown"`
}

// VCSRepo contains the configuration of a VCS integration.
type VCSRepo struct {
	Branch            string `jsonapi:"attr,branch"`
	DisplayIdentifier string `jsonapi:"attr,display-identifier"`
	Identifier        string `jsonapi:"attr,identifier"`
	IngressSubmodules bool   `jsonapi:"attr,ingress-submodules"`
	OAuthTokenID      string `jsonapi:"attr,oauth-token-id"`
	RepositoryHTTPURL string `jsonapi:"attr,repository-http-url"`
	ServiceProvider   string `jsonapi:"attr,service-provider"`
}

// WorkspaceActions represents the workspace actions.
type WorkspaceActions struct {
	IsDestroyable bool `jsonapi:"attr,is-destroyable"`
}

// WorkspacePermissions represents the workspace permissions.
type WorkspacePermissions struct {
	CanDestroy        bool `jsonapi:"attr,can-destroy"`
	CanForceUnlock    bool `jsonapi:"attr,can-force-unlock"`
	CanLock           bool `jsonapi:"attr,can-lock"`
	CanQueueApply     bool `jsonapi:"attr,can-queue-apply"`
	CanQueueDestroy   bool `jsonapi:"attr,can-queue-destroy"`
	CanQueueRun       bool `jsonapi:"attr,can-queue-run"`
	CanReadSettings   bool `jsonapi:"attr,can-read-settings"`
	CanUnlock         bool `jsonapi:"attr,can-unlock"`
	CanUpdate         bool `jsonapi:"attr,can-update"`
	CanUpdateVariable bool `jsonapi:"attr,can-update-variable"`
}

// WorkspaceListOptions represents the options for listing workspaces.
type WorkspaceListOptions struct {
	ListOptions

	// A search string (partial workspace name) used to filter the results.
	Search *string `url:"search[name],omitempty"`

	// A list of relations to include. See available resources https://www.terraform.io/docs/cloud/api/workspaces.html#available-related-resources
	Include *string `url:"include"`
}

// List all the workspaces within an organization.
func (s *workspaces) List(ctx context.Context, organization string, options WorkspaceListOptions) (*WorkspaceList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/workspaces", url.QueryEscape(organization))
	req, err := s.client.newRequest("GET", u, &options)
	if err != nil {
		return nil, err
	}

	wl := &WorkspaceList{}
	err = s.client.do(ctx, req, wl)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

// WorkspaceCreateOptions represents the options for creating a new workspace.
type WorkspaceCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// Required when execution-mode is set to agent. The ID of the agent pool
	// belonging to the workspace's organization. This value must not be specified
	// if execution-mode is set to remote or local or if operations is set to true.
	AgentPoolID *string `jsonapi:"attr,agent-pool-id,omitempty"`

	// Whether destroy plans can be queued on the workspace.
	AllowDestroyPlan *bool `jsonapi:"attr,allow-destroy-plan,omitempty"`

	// Whether to automatically apply changes when a Terraform plan is successful.
	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`

	// A description for the workspace.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Which execution mode to use. Valid values are remote, local, and agent.
	// When set to local, the workspace will be used for state storage only.
	// This value must not be specified if operations is specified.
	// 'agent' execution mode is not available in Terraform Enterprise.
	ExecutionMode *string `jsonapi:"attr,execution-mode,omitempty"`

	// Whether to filter runs based on the changed files in a VCS push. If
	// enabled, the working directory and trigger prefixes describe a set of
	// paths which must contain changes for a VCS push to trigger a run. If
	// disabled, any push will trigger a run.
	FileTriggersEnabled *bool `jsonapi:"attr,file-triggers-enabled,omitempty"`

	GlobalRemoteState *bool `jsonapi:"attr,global-remote-state,omitempty"`

	// The legacy TFE environment to use as the source of the migration, in the
	// form organization/environment. Omit this unless you are migrating a legacy
	// environment.
	MigrationEnvironment *string `jsonapi:"attr,migration-environment,omitempty"`

	// The name of the workspace, which can only include letters, numbers, -,
	// and _. This will be used as an identifier and must be unique in the
	// organization.
	Name *string `jsonapi:"attr,name"`

	// DEPRECATED. Whether the workspace will use remote or local execution mode.
	// Use ExecutionMode instead.
	Operations *bool `jsonapi:"attr,operations,omitempty"`

	// Whether to queue all runs. Unless this is set to true, runs triggered by
	// a webhook will not be queued until at least one run is manually queued.
	QueueAllRuns *bool `jsonapi:"attr,queue-all-runs,omitempty"`

	// Whether this workspace allows speculative plans. Setting this to false
	// prevents Terraform Cloud or the Terraform Enterprise instance from
	// running plans on pull requests, which can improve security if the VCS
	// repository is public or includes untrusted contributors.
	SpeculativeEnabled *bool `jsonapi:"attr,speculative-enabled,omitempty"`

	// The version of Terraform to use for this workspace. Upon creating a
	// workspace, the latest version is selected unless otherwise specified.
	TerraformVersion *string `jsonapi:"attr,terraform-version,omitempty"`

	// List of repository-root-relative paths which list all locations to be
	// tracked for changes. See FileTriggersEnabled above for more details.
	TriggerPrefixes []string `jsonapi:"attr,trigger-prefixes,omitempty"`

	// Settings for the workspace's VCS repository. If omitted, the workspace is
	// created without a VCS repo. If included, you must specify at least the
	// oauth-token-id and identifier keys below.
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`

	// A relative path that Terraform will execute within. This defaults to the
	// root of your repository and is typically set to a subdirectory matching the
	// environment when multiple environments exist within the same repository.
	WorkingDirectory *string `jsonapi:"attr,working-directory,omitempty"`
}

// TODO: move this struct out. VCSRepoOptions is used by workspaces, policy sets, and registry modules
// VCSRepoOptions represents the configuration options of a VCS integration.
type VCSRepoOptions struct {
	Branch            *string `json:"branch,omitempty"`
	Identifier        *string `json:"identifier,omitempty"`
	IngressSubmodules *bool   `json:"ingress-submodules,omitempty"`
	OAuthTokenID      *string `json:"oauth-token-id,omitempty"`
}

func (o WorkspaceCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if !validStringID(o.Name) {
		return ErrInvalidName
	}
	if o.Operations != nil && o.ExecutionMode != nil {
		return errors.New("operations is deprecated and cannot be specified when execution mode is used")
	}
	if o.AgentPoolID != nil && (o.ExecutionMode == nil || *o.ExecutionMode != "agent") {
		return errors.New("specifying an agent pool ID requires 'agent' execution mode")
	}
	if o.AgentPoolID == nil && (o.ExecutionMode != nil && *o.ExecutionMode == "agent") {
		return errors.New("'agent' execution mode requires an agent pool ID to be specified")
	}

	return nil
}

// Create is used to create a new workspace.
func (s *workspaces) Create(ctx context.Context, organization string, options WorkspaceCreateOptions) (*Workspace, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("organizations/%s/workspaces", url.QueryEscape(organization))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Read a workspace by its name.
func (s *workspaces) Read(ctx context.Context, organization, workspace string) (*Workspace, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return nil, ErrInvalidWorkspaceValue
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s",
		url.QueryEscape(organization),
		url.QueryEscape(workspace),
	)
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	// durations come over in ms
	w.ApplyDurationAverage *= time.Millisecond
	w.PlanDurationAverage *= time.Millisecond

	return w, nil
}

// ReadByID reads a workspace by its ID.
func (s *workspaces) ReadByID(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	// durations come over in ms
	w.ApplyDurationAverage *= time.Millisecond
	w.PlanDurationAverage *= time.Millisecond

	return w, nil
}

// Readme gets the readme of a workspace by its ID.
func (s *workspaces) Readme(ctx context.Context, workspaceID string) (io.Reader, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s?include=readme", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	r := &workspaceWithReadme{}
	err = s.client.do(ctx, req, r)
	if err != nil {
		return nil, err
	}
	if r.Readme == nil {
		return nil, nil
	}

	return strings.NewReader(r.Readme.RawMarkdown), nil
}

// WorkspaceUpdateOptions represents the options for updating a workspace.
type WorkspaceUpdateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// Required when execution-mode is set to agent. The ID of the agent pool
	// belonging to the workspace's organization. This value must not be specified
	// if execution-mode is set to remote or local or if operations is set to true.
	AgentPoolID *string `jsonapi:"attr,agent-pool-id,omitempty"`

	// Whether destroy plans can be queued on the workspace.
	AllowDestroyPlan *bool `jsonapi:"attr,allow-destroy-plan,omitempty"`

	// Whether to automatically apply changes when a Terraform plan is successful.
	AutoApply *bool `jsonapi:"attr,auto-apply,omitempty"`

	// A new name for the workspace, which can only include letters, numbers, -,
	// and _. This will be used as an identifier and must be unique in the
	// organization. Warning: Changing a workspace's name changes its URL in the
	// API and UI.
	Name *string `jsonapi:"attr,name,omitempty"`

	// A description for the workspace.
	Description *string `jsonapi:"attr,description,omitempty"`

	// Which execution mode to use. Valid values are remote, local, and agent.
	// When set to local, the workspace will be used for state storage only.
	// This value must not be specified if operations is specified.
	// 'agent' execution mode is not available in Terraform Enterprise.
	ExecutionMode *string `jsonapi:"attr,execution-mode,omitempty"`

	// Whether to filter runs based on the changed files in a VCS push. If
	// enabled, the working directory and trigger prefixes describe a set of
	// paths which must contain changes for a VCS push to trigger a run. If
	// disabled, any push will trigger a run.
	FileTriggersEnabled *bool `jsonapi:"attr,file-triggers-enabled,omitempty"`

	GlobalRemoteState *bool `jsonapi:"attr,global-remote-state,omitempty"`

	// DEPRECATED. Whether the workspace will use remote or local execution mode.
	// Use ExecutionMode instead.
	Operations *bool `jsonapi:"attr,operations,omitempty"`

	// Whether to queue all runs. Unless this is set to true, runs triggered by
	// a webhook will not be queued until at least one run is manually queued.
	QueueAllRuns *bool `jsonapi:"attr,queue-all-runs,omitempty"`

	// Whether this workspace allows speculative plans. Setting this to false
	// prevents Terraform Cloud or the Terraform Enterprise instance from
	// running plans on pull requests, which can improve security if the VCS
	// repository is public or includes untrusted contributors.
	SpeculativeEnabled *bool `jsonapi:"attr,speculative-enabled,omitempty"`

	// The version of Terraform to use for this workspace.
	TerraformVersion *string `jsonapi:"attr,terraform-version,omitempty"`

	// List of repository-root-relative paths which list all locations to be
	// tracked for changes. See FileTriggersEnabled above for more details.
	TriggerPrefixes []string `jsonapi:"attr,trigger-prefixes,omitempty"`

	// To delete a workspace's existing VCS repo, specify null instead of an
	// object. To modify a workspace's existing VCS repo, include whichever of
	// the keys below you wish to modify. To add a new VCS repo to a workspace
	// that didn't previously have one, include at least the oauth-token-id and
	// identifier keys.
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo,omitempty"`

	// A relative path that Terraform will execute within. This defaults to the
	// root of your repository and is typically set to a subdirectory matching
	// the environment when multiple environments exist within the same
	// repository.
	WorkingDirectory *string `jsonapi:"attr,working-directory,omitempty"`
}

func (o WorkspaceUpdateOptions) valid() error {
	if o.Name != nil && !validStringID(o.Name) {
		return ErrInvalidName
	}
	if o.Operations != nil && o.ExecutionMode != nil {
		return errors.New("operations is deprecated and cannot be specified when execution mode is used")
	}
	if o.AgentPoolID == nil && (o.ExecutionMode != nil && *o.ExecutionMode == "agent") {
		return errors.New("'agent' execution mode requires an agent pool ID to be specified")
	}

	return nil
}

// Update settings of an existing workspace.
func (s *workspaces) Update(ctx context.Context, organization, workspace string, options WorkspaceUpdateOptions) (*Workspace, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return nil, ErrInvalidWorkspaceValue
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s",
		url.QueryEscape(organization),
		url.QueryEscape(workspace),
	)
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// UpdateByID updates the settings of an existing workspace.
func (s *workspaces) UpdateByID(ctx context.Context, workspaceID string, options WorkspaceUpdateOptions) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Delete a workspace by its name.
func (s *workspaces) Delete(ctx context.Context, organization, workspace string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return ErrInvalidWorkspaceValue
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s",
		url.QueryEscape(organization),
		url.QueryEscape(workspace),
	)
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// DeleteByID deletes a workspace by its ID.
func (s *workspaces) DeleteByID(ctx context.Context, workspaceID string) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// workspaceRemoveVCSConnectionOptions
type workspaceRemoveVCSConnectionOptions struct {
	ID      string          `jsonapi:"primary,workspaces"`
	VCSRepo *VCSRepoOptions `jsonapi:"attr,vcs-repo"`
}

// RemoveVCSConnection from a workspace.
func (s *workspaces) RemoveVCSConnection(ctx context.Context, organization, workspace string) (*Workspace, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if !validStringID(&workspace) {
		return nil, ErrInvalidWorkspaceValue
	}

	u := fmt.Sprintf(
		"organizations/%s/workspaces/%s",
		url.QueryEscape(organization),
		url.QueryEscape(workspace),
	)

	req, err := s.client.newRequest("PATCH", u, &workspaceRemoveVCSConnectionOptions{})
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// RemoveVCSConnectionByID removes a VCS connection from a workspace.
func (s *workspaces) RemoveVCSConnectionByID(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s", url.QueryEscape(workspaceID))

	req, err := s.client.newRequest("PATCH", u, &workspaceRemoveVCSConnectionOptions{})
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// WorkspaceLockOptions represents the options for locking a workspace.
type WorkspaceLockOptions struct {
	// Specifies the reason for locking the workspace.
	Reason *string `jsonapi:"attr,reason,omitempty"`
}

// Lock a workspace by its ID.
func (s *workspaces) Lock(ctx context.Context, workspaceID string, options WorkspaceLockOptions) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/actions/lock", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Unlock a workspace by its ID.
func (s *workspaces) Unlock(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/actions/unlock", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// ForceUnlock a workspace by its ID.
func (s *workspaces) ForceUnlock(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/actions/force-unlock", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("POST", u, nil)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// WorkspaceAssignSSHKeyOptions represents the options to assign an SSH key to
// a workspace.
type WorkspaceAssignSSHKeyOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// The SSH key ID to assign.
	SSHKeyID *string `jsonapi:"attr,id"`
}

func (o WorkspaceAssignSSHKeyOptions) valid() error {
	if !validString(o.SSHKeyID) {
		return errors.New("SSH key ID is required")
	}
	if !validStringID(o.SSHKeyID) {
		return errors.New("invalid value for SSH key ID")
	}
	return nil
}

// AssignSSHKey to a workspace.
func (s *workspaces) AssignSSHKey(ctx context.Context, workspaceID string, options WorkspaceAssignSSHKeyOptions) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/ssh-key", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("PATCH", u, &options)
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// workspaceUnassignSSHKeyOptions represents the options to unassign an SSH key
// to a workspace.
type workspaceUnassignSSHKeyOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,workspaces"`

	// Must be nil to unset the currently assigned SSH key.
	SSHKeyID *string `jsonapi:"attr,id"`
}

// UnassignSSHKey from a workspace.
func (s *workspaces) UnassignSSHKey(ctx context.Context, workspaceID string) (*Workspace, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/relationships/ssh-key", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("PATCH", u, &workspaceUnassignSSHKeyOptions{})
	if err != nil {
		return nil, err
	}

	w := &Workspace{}
	err = s.client.do(ctx, req, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// RemoteStateConsumers returns the remote state consumers for a given workspace.
func (s *workspaces) RemoteStateConsumers(ctx context.Context, workspaceID string) (*WorkspaceList, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}

	u := fmt.Sprintf("workspaces/%s/relationships/remote-state-consumers", url.QueryEscape(workspaceID))

	req, err := s.client.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	wl := &WorkspaceList{}
	err = s.client.do(ctx, req, wl)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

// WorkspaceAddRemoteStateConsumersOptions represents the options for adding remote state consumers
// to a workspace.
type WorkspaceAddRemoteStateConsumersOptions struct {
	/// The workspaces to add as remote state consumers to the workspace.
	Workspaces []*Workspace
}

func (o WorkspaceAddRemoteStateConsumersOptions) valid() error {
	if o.Workspaces == nil {
		return ErrWorkspacesRequired
	}
	if len(o.Workspaces) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

// AddRemoteStateConsumere adds the remote state consumers to a given workspace.
func (s *workspaces) AddRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceAddRemoteStateConsumersOptions) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/remote-state-consumers", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("POST", u, options.Workspaces)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// WorkspaceRemoveRemoteStateConsumersOptions represents the options for removing remote state
// consumers from a workspace.
type WorkspaceRemoveRemoteStateConsumersOptions struct {
	/// The workspaces to remove as remote state consumers from the workspace.
	Workspaces []*Workspace
}

func (o WorkspaceRemoveRemoteStateConsumersOptions) valid() error {
	if o.Workspaces == nil {
		return ErrWorkspacesRequired
	}
	if len(o.Workspaces) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

// RemoveRemoteStateConsumers removes the remote state consumers for a given workspace.
func (s *workspaces) RemoveRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceRemoveRemoteStateConsumersOptions) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/remote-state-consumers", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("DELETE", u, options.Workspaces)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}

// WorkspaceUpdateRemoteStateConsumersOptions represents the options for
// updatintg remote state consumers from a workspace.
type WorkspaceUpdateRemoteStateConsumersOptions struct {
	/// The workspaces to update remote state consumers for the workspace.
	Workspaces []*Workspace
}

func (o WorkspaceUpdateRemoteStateConsumersOptions) valid() error {
	if o.Workspaces == nil {
		return ErrWorkspacesRequired
	}
	if len(o.Workspaces) == 0 {
		return ErrWorkspaceMinLimit
	}
	return nil
}

// UpdateRemoteStateConsumers removes the remote state consumers for a given workspace.
func (s *workspaces) UpdateRemoteStateConsumers(ctx context.Context, workspaceID string, options WorkspaceUpdateRemoteStateConsumersOptions) error {
	if !validStringID(&workspaceID) {
		return ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf("workspaces/%s/relationships/remote-state-consumers", url.QueryEscape(workspaceID))
	req, err := s.client.newRequest("PATCH", u, options.Workspaces)
	if err != nil {
		return err
	}

	return s.client.do(ctx, req, nil)
}
