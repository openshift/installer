package tfe

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"
)

// Compile-time proof of interface implementation.
var _ StateVersions = (*stateVersions)(nil)

// StateVersions describes all the state version related methods that
// the Terraform Enterprise API supports.
//
// TFE API docs:
// https://www.terraform.io/docs/cloud/api/state-versions.html
type StateVersions interface {
	// List all the state versions for a given workspace.
	List(ctx context.Context, options *StateVersionListOptions) (*StateVersionList, error)

	// Create a new state version for the given workspace.
	Create(ctx context.Context, workspaceID string, options StateVersionCreateOptions) (*StateVersion, error)

	// Read a state version by its ID.
	Read(ctx context.Context, svID string) (*StateVersion, error)

	// ReadWithOptions reads a state version by its ID using the options supplied
	ReadWithOptions(ctx context.Context, svID string, options *StateVersionReadOptions) (*StateVersion, error)

	// ReadCurrent reads the latest available state from the given workspace.
	ReadCurrent(ctx context.Context, workspaceID string) (*StateVersion, error)

	// ReadCurrentWithOptions reads the latest available state from the given workspace using the options supplied
	ReadCurrentWithOptions(ctx context.Context, workspaceID string, options *StateVersionCurrentOptions) (*StateVersion, error)

	// Download retrieves the actual stored state of a state version
	Download(ctx context.Context, url string) ([]byte, error)

	// ListOutputs retrieves all the outputs of a state version by its ID.
	ListOutputs(ctx context.Context, svID string, options *StateVersionOutputsListOptions) (*StateVersionOutputsList, error)
}

// stateVersions implements StateVersions.
type stateVersions struct {
	client *Client
}

// StateVersionList represents a list of state versions.
type StateVersionList struct {
	*Pagination
	Items []*StateVersion
}

// StateVersion represents a Terraform Enterprise state version.
type StateVersion struct {
	ID                 string                   `jsonapi:"primary,state-versions"`
	CreatedAt          time.Time                `jsonapi:"attr,created-at,iso8601"`
	DownloadURL        string                   `jsonapi:"attr,hosted-state-download-url"`
	Serial             int64                    `jsonapi:"attr,serial"`
	VCSCommitSHA       string                   `jsonapi:"attr,vcs-commit-sha"`
	VCSCommitURL       string                   `jsonapi:"attr,vcs-commit-url"`
	ResourcesProcessed bool                     `jsonapi:"attr,resources-processed"`
	StateVersion       int                      `jsonapi:"attr,state-version"`
	TerraformVersion   string                   `jsonapi:"attr,terraform-version"`
	Modules            *StateVersionModules     `jsonapi:"attr,modules"`
	Providers          *StateVersionProviders   `jsonapi:"attr,providers"`
	Resources          []*StateVersionResources `jsonapi:"attr,resources"`

	// Relations
	Run     *Run                  `jsonapi:"relation,run"`
	Outputs []*StateVersionOutput `jsonapi:"relation,outputs"`
}

// StateVersionOutputsList represents a list of StateVersionOutput items.
type StateVersionOutputsList struct {
	*Pagination
	Items []*StateVersionOutput
}

// StateVersionListOptions represents the options for listing state versions.
type StateVersionListOptions struct {
	ListOptions
	Organization string `url:"filter[organization][name]"`
	Workspace    string `url:"filter[workspace][name]"`
}

// StateVersionIncludeOpt represents the available options for include query params.
// https://www.terraform.io/cloud-docs/api-docs/state-versions#available-related-resources
type StateVersionIncludeOpt string

const (
	SVcreatedby               StateVersionIncludeOpt = "created_by"
	SVrun                     StateVersionIncludeOpt = "run"
	SVrunCreatedBy            StateVersionIncludeOpt = "run.created_by"
	SVrunConfigurationVersion StateVersionIncludeOpt = "run.configuration_version"
	SVoutputs                 StateVersionIncludeOpt = "outputs"
)

// StateVersionReadOptions represents the options for reading state version.
type StateVersionReadOptions struct {
	// Optional: A list of relations to include. See available resources:
	// https://www.terraform.io/cloud-docs/api-docs/state-versions#available-related-resources
	Include []StateVersionIncludeOpt `url:"include,omitempty"`
}

// StateVersionOutputsListOptions represents the options for listing state
// version outputs.
type StateVersionOutputsListOptions struct {
	ListOptions
}

// StateVersionCurrentOptions represents the options for reading the current state version.
type StateVersionCurrentOptions struct {
	// Optional: A list of relations to include. See available resources:
	// https://www.terraform.io/cloud-docs/api-docs/state-versions#available-related-resources
	Include []StateVersionIncludeOpt `url:"include,omitempty"`
}

// StateVersionCreateOptions represents the options for creating a state version.
type StateVersionCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,state-versions"`

	// Optional: The lineage of the state.
	Lineage *string `jsonapi:"attr,lineage,omitempty"`

	// Required: The MD5 hash of the state version.
	MD5 *string `jsonapi:"attr,md5"`

	// Required: The serial of the state.
	Serial *int64 `jsonapi:"attr,serial"`

	// Required: The base64 encoded state.
	State *string `jsonapi:"attr,state"`

	// Optional: Force can be set to skip certain validations. Wrong use
	// of this flag can cause data loss, so USE WITH CAUTION!
	Force *bool `jsonapi:"attr,force,omitempty"`

	// Optional: Specifies the run to associate the state with.
	Run *Run `jsonapi:"relation,run,omitempty"`

	// Optional: The external, json representation of state data, base64 encoded.
	// https://www.terraform.io/internals/json-format#state-representation
	// Supplying this state representation can provide more details to the platform
	// about the current terraform state.
	//
	// **Note**: This field is in BETA, subject to change and not widely available yet.
	JSONState *string `jsonapi:"attr,json-state,omitempty"`
	// Optional: The external, json representation of state outputs, base64 encoded. Supplying this field
	// will provide more detailed output type information to TFE.
	// For more information on the contents of this field: https://www.terraform.io/internals/json-format#values-representation
	// about the current terraform state.
	//
	// **Note**: This field is in BETA, subject to change and not widely available yet.
	JSONStateOutputs *string `jsonapi:"attr,json-state-outputs,omitempty"`
}

type StateVersionModules struct {
	Root StateVersionModuleRoot `jsonapi:"attr,root"`
}

type StateVersionModuleRoot struct {
	NullResource         int `jsonapi:"attr,null-resource"`
	TerraformRemoteState int `jsonapi:"attr,data.terraform-remote-state"`
}

type StateVersionProviders struct {
	Data ProviderData `jsonapi:"attr,provider[map]string"`
}

type ProviderData struct {
	NullResource         int `json:"null-resource"`
	TerraformRemoteState int `json:"data.terraform-remote-state"`
}

type StateVersionResources struct {
	Name     string `jsonapi:"attr,name"`
	Count    string `jsonapi:"attr,count"`
	Type     int    `jsonapi:"attr,type"`
	Module   string `jsonapi:"attr,module"`
	Provider string `jsonapi:"attr,provider"`
}

// List all the state versions for a given workspace.
func (s *stateVersions) List(ctx context.Context, options *StateVersionListOptions) (*StateVersionList, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", "state-versions", options)
	if err != nil {
		return nil, err
	}

	svl := &StateVersionList{}
	err = req.Do(ctx, svl)
	if err != nil {
		return nil, err
	}

	return svl, nil
}

// Create a new state version for the given workspace.
func (s *stateVersions) Create(ctx context.Context, workspaceID string, options StateVersionCreateOptions) (*StateVersion, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/state-versions", url.QueryEscape(workspaceID))
	req, err := s.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	sv := &StateVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// Read a state version by its ID.
func (s *stateVersions) ReadWithOptions(ctx context.Context, svID string, options *StateVersionReadOptions) (*StateVersion, error) {
	if !validStringID(&svID) {
		return nil, ErrInvalidStateVerID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("state-versions/%s", url.QueryEscape(svID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	sv := &StateVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// Read a state version by its ID.
func (s *stateVersions) Read(ctx context.Context, svID string) (*StateVersion, error) {
	return s.ReadWithOptions(ctx, svID, nil)
}

// ReadCurrentWithOptions reads the latest available state from the given workspace using the options supplied.
func (s *stateVersions) ReadCurrentWithOptions(ctx context.Context, workspaceID string, options *StateVersionCurrentOptions) (*StateVersion, error) {
	if !validStringID(&workspaceID) {
		return nil, ErrInvalidWorkspaceID
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("workspaces/%s/current-state-version", url.QueryEscape(workspaceID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	sv := &StateVersion{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// ReadCurrent reads the latest available state from the given workspace.
func (s *stateVersions) ReadCurrent(ctx context.Context, workspaceID string) (*StateVersion, error) {
	return s.ReadCurrentWithOptions(ctx, workspaceID, nil)
}

// Download retrieves the actual stored state of a state version
func (s *stateVersions) Download(ctx context.Context, u string) ([]byte, error) {
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	var buf bytes.Buffer
	err = req.Do(ctx, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ListOutputs retrieves all the outputs of a state version by its ID.
func (s *stateVersions) ListOutputs(ctx context.Context, svID string, options *StateVersionOutputsListOptions) (*StateVersionOutputsList, error) {
	if !validStringID(&svID) {
		return nil, ErrInvalidStateVerID
	}

	u := fmt.Sprintf("state-versions/%s/outputs", url.QueryEscape(svID))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	sv := &StateVersionOutputsList{}
	err = req.Do(ctx, sv)
	if err != nil {
		return nil, err
	}

	return sv, nil
}

// check that StateVersionListOptions fields had valid values
func (o *StateVersionListOptions) valid() error {
	if o == nil {
		return ErrRequiredStateVerListOps
	}
	if !validString(&o.Organization) {
		return ErrRequiredOrg
	}
	if !validString(&o.Workspace) {
		return ErrRequiredWorkspace
	}
	return nil
}

func (o StateVersionCreateOptions) valid() error {
	if !validString(o.MD5) {
		return ErrRequiredM5
	}
	if o.Serial == nil {
		return ErrRequiredSerial
	}
	if !validString(o.State) {
		return ErrRequiredState
	}
	return nil
}

func (o *StateVersionReadOptions) valid() error {
	if o == nil {
		return nil // nothing to validate
	}

	if err := validateStateVerIncludeParams(o.Include); err != nil {
		return err
	}
	return nil
}
func (o *StateVersionCurrentOptions) valid() error {
	if o == nil {
		return nil // nothing to validate
	}

	if err := validateStateVerIncludeParams(o.Include); err != nil {
		return err
	}
	return nil
}

func validateStateVerIncludeParams(params []StateVersionIncludeOpt) error {
	for _, p := range params {
		switch p {
		case SVcreatedby, SVrun, SVrunCreatedBy, SVrunConfigurationVersion, SVoutputs:
			// do nothing
		default:
			return ErrInvalidIncludeValue
		}
	}

	return nil
}
