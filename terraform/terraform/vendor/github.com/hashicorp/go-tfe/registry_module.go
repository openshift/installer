package tfe

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
)

// Compile-time proof of interface implementation.
var _ RegistryModules = (*registryModules)(nil)

// RegistryModules describes all the registry module related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/modules.html
type RegistryModules interface {
	// List all the registory modules within an organization
	List(ctx context.Context, organization string, options *RegistryModuleListOptions) (*RegistryModuleList, error)

	// Create a registry module without a VCS repo
	Create(ctx context.Context, organization string, options RegistryModuleCreateOptions) (*RegistryModule, error)

	// Create a registry module version
	CreateVersion(ctx context.Context, moduleID RegistryModuleID, options RegistryModuleCreateVersionOptions) (*RegistryModuleVersion, error)

	// Create and publish a registry module with a VCS repo
	CreateWithVCSConnection(ctx context.Context, options RegistryModuleCreateWithVCSConnectionOptions) (*RegistryModule, error)

	// Read a registry module
	Read(ctx context.Context, moduleID RegistryModuleID) (*RegistryModule, error)

	// Delete a registry module
	Delete(ctx context.Context, organization string, name string) error

	// Delete a specific registry module provider
	DeleteProvider(ctx context.Context, moduleID RegistryModuleID) error

	// Delete a specific registry module version
	DeleteVersion(ctx context.Context, moduleID RegistryModuleID, version string) error

	// Upload Terraform configuration files for the provided registry module version. It
	// requires a path to the configuration files on disk, which will be packaged by
	// hashicorp/go-slug before being uploaded.
	Upload(ctx context.Context, rmv RegistryModuleVersion, path string) error
}

// registryModules implements RegistryModules.
type registryModules struct {
	client *Client
}

// RegistryModuleStatus represents the status of the registry module
type RegistryModuleStatus string

// List of available registry module statuses
const (
	RegistryModuleStatusPending       RegistryModuleStatus = "pending"
	RegistryModuleStatusNoVersionTags RegistryModuleStatus = "no_version_tags"
	RegistryModuleStatusSetupFailed   RegistryModuleStatus = "setup_failed"
	RegistryModuleStatusSetupComplete RegistryModuleStatus = "setup_complete"
)

// RegistryModuleVersionStatus represents the status of a specific version of a registry module
type RegistryModuleVersionStatus string

// List of available registry module version statuses
const (
	RegistryModuleVersionStatusPending             RegistryModuleVersionStatus = "pending"
	RegistryModuleVersionStatusCloning             RegistryModuleVersionStatus = "cloning"
	RegistryModuleVersionStatusCloneFailed         RegistryModuleVersionStatus = "clone_failed"
	RegistryModuleVersionStatusRegIngressReqFailed RegistryModuleVersionStatus = "reg_ingress_req_failed"
	RegistryModuleVersionStatusRegIngressing       RegistryModuleVersionStatus = "reg_ingressing"
	RegistryModuleVersionStatusRegIngressFailed    RegistryModuleVersionStatus = "reg_ingress_failed"
	RegistryModuleVersionStatusOk                  RegistryModuleVersionStatus = "ok"
)

// RegistryModuleID represents the set of IDs that identify a RegistryModule
type RegistryModuleID struct {
	// The organization the module belongs to, see RegistryModule.Organization.Name
	Organization string
	// The name of the module, see RegistryModule.Name
	Name string
	// The module's provider, see RegistryModule.Provider
	Provider string
	// The namespace of the module. For private modules this is the name of the organization that owns the module
	// Required for public modules
	Namespace string
	// Either public or private. If not provided, defaults to private
	RegistryName RegistryName
}

// RegistryModuleList represents a list of registry modules.
type RegistryModuleList struct {
	*Pagination
	Items []*RegistryModule
}

// RegistryModule represents a registry module
type RegistryModule struct {
	ID              string                          `jsonapi:"primary,registry-modules"`
	Name            string                          `jsonapi:"attr,name"`
	Provider        string                          `jsonapi:"attr,provider"`
	RegistryName    RegistryName                    `jsonapi:"attr,registry-name"`
	Namespace       string                          `jsonapi:"attr,namespace"`
	Permissions     *RegistryModulePermissions      `jsonapi:"attr,permissions"`
	Status          RegistryModuleStatus            `jsonapi:"attr,status"`
	VCSRepo         *VCSRepo                        `jsonapi:"attr,vcs-repo"`
	VersionStatuses []RegistryModuleVersionStatuses `jsonapi:"attr,version-statuses"`
	CreatedAt       string                          `jsonapi:"attr,created-at"`
	UpdatedAt       string                          `jsonapi:"attr,updated-at"`

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`
}

// RegistryModuleVersion represents a registry module version
type RegistryModuleVersion struct {
	ID        string                      `jsonapi:"primary,registry-module-versions"`
	Source    string                      `jsonapi:"attr,source"`
	Status    RegistryModuleVersionStatus `jsonapi:"attr,status"`
	Version   string                      `jsonapi:"attr,version"`
	CreatedAt string                      `jsonapi:"attr,created-at"`
	UpdatedAt string                      `jsonapi:"attr,updated-at"`

	// Relations
	RegistryModule *RegistryModule `jsonapi:"relation,registry-module"`

	// Links
	Links map[string]interface{} `jsonapi:"links,omitempty"`
}

type RegistryModulePermissions struct {
	CanDelete bool `jsonapi:"attr,can-delete"`
	CanResync bool `jsonapi:"attr,can-resync"`
	CanRetry  bool `jsonapi:"attr,can-retry"`
}

type RegistryModuleVersionStatuses struct {
	Version string                      `jsonapi:"attr,version"`
	Status  RegistryModuleVersionStatus `jsonapi:"attr,status"`
	Error   string                      `jsonapi:"attr,error"`
}

// RegistryModuleListOptions represents the options for listing registry modules.
type RegistryModuleListOptions struct {
	ListOptions
}

// RegistryModuleCreateOptions is used when creating a registry module without a VCS repo
type RegistryModuleCreateOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,registry-modules"`
	// Required:
	Name *string `jsonapi:"attr,name"`
	// Required:
	Provider *string `jsonapi:"attr,provider"`
	// Optional: Whether this is a publicly maintained module or private. Must be either public or private.
	// Defaults to private if not specified
	RegistryName RegistryName `jsonapi:"attr,registry-name,omitempty"`
	// Optional: The namespace of this module. Required for public modules only.
	Namespace string `jsonapi:"attr,namespace,omitempty"`
}

// RegistryModuleCreateVersionOptions is used when creating a registry module version
type RegistryModuleCreateVersionOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,registry-module-versions"`

	Version *string `jsonapi:"attr,version"`
}

// RegistryModuleCreateWithVCSConnectionOptions is used when creating a registry module with a VCS repo
type RegistryModuleCreateWithVCSConnectionOptions struct {
	// Type is a public field utilized by JSON:API to
	// set the resource type via the field tag.
	// It is not a user-defined value and does not need to be set.
	// https://jsonapi.org/format/#crud-creating
	Type string `jsonapi:"primary,registry-modules"`

	// Required: VCS repository information
	VCSRepo *RegistryModuleVCSRepoOptions `jsonapi:"attr,vcs-repo"`
}

type RegistryModuleVCSRepoOptions struct {
	Identifier        *string `json:"identifier"`         // Required
	OAuthTokenID      *string `json:"oauth-token-id"`     // Required
	DisplayIdentifier *string `json:"display-identifier"` // Required
}

// List all the registory modules within an organization.
func (s *registryModules) List(ctx context.Context, organization string, options *RegistryModuleListOptions) (*RegistryModuleList, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}

	u := fmt.Sprintf("organizations/%s/registry-modules", url.QueryEscape(organization))
	req, err := s.client.NewRequest("GET", u, options)
	if err != nil {
		return nil, err
	}

	ml := &RegistryModuleList{}
	err = req.Do(ctx, ml)
	if err != nil {
		return nil, err
	}

	return ml, nil
}

// Upload uploads Terraform configuration files for the provided registry module version. It
// requires a path to the configuration files on disk, which will be packaged by
// hashicorp/go-slug before being uploaded.
func (r *registryModules) Upload(ctx context.Context, rmv RegistryModuleVersion, path string) error {
	uploadURL, ok := rmv.Links["upload"].(string)
	if !ok {
		return fmt.Errorf("provided RegistryModuleVersion does not contain an upload link")
	}

	body, err := packContents(path)
	if err != nil {
		return err
	}

	req, err := r.client.NewRequest("PUT", uploadURL, body)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// Create a new registry module without a VCS repo
func (r *registryModules) Create(ctx context.Context, organization string, options RegistryModuleCreateOptions) (*RegistryModule, error) {
	if !validStringID(&organization) {
		return nil, ErrInvalidOrg
	}
	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf(
		"organizations/%s/registry-modules",
		url.QueryEscape(organization),
	)
	req, err := r.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// CreateVersion creates a new registry module version
func (r *registryModules) CreateVersion(ctx context.Context, moduleID RegistryModuleID, options RegistryModuleCreateVersionOptions) (*RegistryModuleVersion, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	if err := options.valid(); err != nil {
		return nil, err
	}

	u := fmt.Sprintf(
		"registry-modules/%s/%s/%s/versions",
		url.QueryEscape(moduleID.Organization),
		url.QueryEscape(moduleID.Name),
		url.QueryEscape(moduleID.Provider),
	)
	req, err := r.client.NewRequest("POST", u, &options)
	if err != nil {
		return nil, err
	}

	rmv := &RegistryModuleVersion{}
	err = req.Do(ctx, rmv)
	if err != nil {
		return nil, err
	}

	return rmv, nil
}

// CreateWithVCSConnection is used to create and publish a new registry module with a VCS repo
func (r *registryModules) CreateWithVCSConnection(ctx context.Context, options RegistryModuleCreateWithVCSConnectionOptions) (*RegistryModule, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := r.client.NewRequest("POST", "registry-modules", &options)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// Read a specific registry module
func (r *registryModules) Read(ctx context.Context, moduleID RegistryModuleID) (*RegistryModule, error) {
	if err := moduleID.valid(); err != nil {
		return nil, err
	}

	if moduleID.RegistryName == "" {
		log.Println("[WARN] Support for using the RegistryModuleID without RegistryName is deprecated as of release 1.5.0 and may be removed in a future version. The preferred method is to include the RegistryName in RegistryModuleID.")
		moduleID.RegistryName = PrivateRegistry
	}

	if moduleID.RegistryName == PrivateRegistry && strings.TrimSpace(moduleID.Namespace) == "" {
		log.Println("[WARN] Support for using the RegistryModuleID without Namespace is deprecated as of release 1.5.0 and may be removed in a future version. The preferred method is to include the Namespace in RegistryModuleID.")
		moduleID.Namespace = moduleID.Organization
	}

	u := fmt.Sprintf(
		"organizations/%s/registry-modules/%s/%s/%s/%s",
		url.QueryEscape(moduleID.Organization),
		url.QueryEscape(string(moduleID.RegistryName)),
		url.QueryEscape(moduleID.Namespace),
		url.QueryEscape(moduleID.Name),
		url.QueryEscape(moduleID.Provider),
	)

	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	rm := &RegistryModule{}
	err = req.Do(ctx, rm)
	if err != nil {
		return nil, err
	}

	return rm, nil
}

// Delete is used to delete the entire registry module
func (r *registryModules) Delete(ctx context.Context, organization, name string) error {
	if !validStringID(&organization) {
		return ErrInvalidOrg
	}
	if !validString(&name) {
		return ErrRequiredName
	}
	if !validStringID(&name) {
		return ErrInvalidName
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s",
		url.QueryEscape(organization),
		url.QueryEscape(name),
	)
	req, err := r.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// DeleteProvider is used to delete the specific registry module provider
func (r *registryModules) DeleteProvider(ctx context.Context, moduleID RegistryModuleID) error {
	if err := moduleID.valid(); err != nil {
		return err
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s/%s",
		url.QueryEscape(moduleID.Organization),
		url.QueryEscape(moduleID.Name),
		url.QueryEscape(moduleID.Provider),
	)
	req, err := r.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

// DeleteVersion is used to delete the specific registry module version
func (r *registryModules) DeleteVersion(ctx context.Context, moduleID RegistryModuleID, version string) error {
	if err := moduleID.valid(); err != nil {
		return err
	}
	if !validString(&version) {
		return ErrRequiredVersion
	}
	if !validStringID(&version) {
		return ErrInvalidVersion
	}

	u := fmt.Sprintf(
		"registry-modules/actions/delete/%s/%s/%s/%s",
		url.QueryEscape(moduleID.Organization),
		url.QueryEscape(moduleID.Name),
		url.QueryEscape(moduleID.Provider),
		url.QueryEscape(version),
	)
	req, err := r.client.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}

	return req.Do(ctx, nil)
}

func (o RegistryModuleID) valid() error {
	if !validStringID(&o.Organization) {
		return ErrInvalidOrg
	}

	if !validString(&o.Name) {
		return ErrRequiredName
	}

	if !validStringID(&o.Name) {
		return ErrInvalidName
	}

	if !validString(&o.Provider) {
		return ErrRequiredProvider
	}

	if !validStringID(&o.Provider) {
		return ErrInvalidProvider
	}

	switch o.RegistryName {
	case PublicRegistry:
		if !validString(&o.Namespace) {
			return ErrRequiredNamespace
		}
	case PrivateRegistry:
	case "":
		// no-op:  RegistryName is optional
	// for all other string
	default:
		return ErrInvalidRegistryName
	}

	return nil
}

func (o RegistryModuleCreateOptions) valid() error {
	if !validString(o.Name) {
		return ErrRequiredName
	}
	if !validStringID(o.Name) {
		return ErrInvalidName
	}
	if !validString(o.Provider) {
		return ErrRequiredProvider
	}
	if !validStringID(o.Provider) {
		return ErrInvalidProvider
	}

	switch o.RegistryName {
	case PublicRegistry:
		if !validString(&o.Namespace) {
			return ErrRequiredNamespace
		}
	case PrivateRegistry:
		if validString(&o.Namespace) {
			return ErrUnsupportedBothNamespaceAndPrivateRegistryName
		}
	case "":
		// no-op:  RegistryName is optional
	// for all other string
	default:
		return ErrInvalidRegistryName
	}
	return nil
}

func (o RegistryModuleCreateVersionOptions) valid() error {
	if !validString(o.Version) {
		return ErrRequiredVersion
	}
	if !validStringID(o.Version) {
		return ErrInvalidVersion
	}
	return nil
}

func (o RegistryModuleCreateWithVCSConnectionOptions) valid() error {
	if o.VCSRepo == nil {
		return ErrRequiredVCSRepo
	}
	return o.VCSRepo.valid()
}

func (o RegistryModuleVCSRepoOptions) valid() error {
	if !validString(o.Identifier) {
		return ErrRequiredIdentifier
	}
	if !validString(o.OAuthTokenID) {
		return ErrRequiredOauthTokenID
	}
	if !validString(o.DisplayIdentifier) {
		return ErrRequiredDisplayIdentifier
	}
	return nil
}
