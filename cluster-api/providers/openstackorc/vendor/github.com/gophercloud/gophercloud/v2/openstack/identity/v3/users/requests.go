package users

import (
	"context"
	"net/url"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Option is a specific option defined at the API to enable features
// on a user account.
type Option string

const (
	IgnoreChangePasswordUponFirstUse Option = "ignore_change_password_upon_first_use"
	IgnorePasswordExpiry             Option = "ignore_password_expiry"
	IgnoreLockoutFailureAttempts     Option = "ignore_lockout_failure_attempts"
	MultiFactorAuthRules             Option = "multi_factor_auth_rules"
	MultiFactorAuthEnabled           Option = "multi_factor_auth_enabled"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToUserListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// DomainID filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Enabled filters the response by enabled users.
	Enabled *bool `q:"enabled"`

	// IdpID filters the response by an Identity Provider ID.
	IdPID string `q:"idp_id"`

	// Name filters the response by username.
	Name string `q:"name"`

	// PasswordExpiresAt filters the response based on expiring passwords.
	PasswordExpiresAt string `q:"password_expires_at"`

	// ProtocolID filters the response by protocol ID.
	ProtocolID string `q:"protocol_id"`

	// UniqueID filters the response by unique ID.
	UniqueID string `q:"unique_id"`

	// Filters filters the response by custom filters such as
	// 'name__contains=foo'
	Filters map[string]string `q:"-"`
}

// ToUserListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToUserListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	params := q.Query()
	for k, v := range opts.Filters {
		i := strings.Index(k, "__")
		if i > 0 && i < len(k)-2 {
			params.Add(k, v)
		} else {
			return "", InvalidListFilter{FilterName: k}
		}
	}

	q = &url.URL{RawQuery: params.Encode()}
	return q.String(), err
}

// List enumerates the Users to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToUserListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single user, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToUserCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a user.
type CreateOpts struct {
	// Name is the name of the new user.
	Name string `json:"name" required:"true"`

	// DefaultProjectID is the ID of the default project of the user.
	DefaultProjectID string `json:"default_project_id,omitempty"`

	// Description is a description of the user.
	Description string `json:"description,omitempty"`

	// DomainID is the ID of the domain the user belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Enabled sets the user status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Extra is free-form extra key/value pairs to describe the user.
	Extra map[string]any `json:"-"`

	// Options are defined options in the API to enable certain features.
	Options map[Option]any `json:"options,omitempty"`

	// Password is the password of the new user.
	Password string `json:"password,omitempty"`
}

// ToUserCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToUserCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "user")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["user"].(map[string]any); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create creates a new User.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToUserCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToUserUpdateMap() (map[string]any, error)
}

// UpdateOpts provides options for updating a user account.
type UpdateOpts struct {
	// Name is the name of the new user.
	Name string `json:"name,omitempty"`

	// DefaultProjectID is the ID of the default project of the user.
	DefaultProjectID string `json:"default_project_id,omitempty"`

	// Description is a description of the user.
	Description *string `json:"description,omitempty"`

	// DomainID is the ID of the domain the user belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Enabled sets the user status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Extra is free-form extra key/value pairs to describe the user.
	Extra map[string]any `json:"-"`

	// Options are defined options in the API to enable certain features.
	Options map[Option]any `json:"options,omitempty"`

	// Password is the password of the new user.
	Password string `json:"password,omitempty"`
}

// ToUserUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToUserUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "user")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["user"].(map[string]any); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Update updates an existing User.
func Update(ctx context.Context, client *gophercloud.ServiceClient, userID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToUserUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, updateURL(client, userID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ChangePasswordOptsBuilder allows extensions to add additional parameters to
// the ChangePassword request.
type ChangePasswordOptsBuilder interface {
	ToUserChangePasswordMap() (map[string]any, error)
}

// ChangePasswordOpts provides options for changing password for a user.
type ChangePasswordOpts struct {
	// OriginalPassword is the original password of the user.
	OriginalPassword string `json:"original_password"`

	// Password is the new password of the user.
	Password string `json:"password"`
}

// ToUserChangePasswordMap formats a ChangePasswordOpts into a ChangePassword request.
func (opts ChangePasswordOpts) ToUserChangePasswordMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "user")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ChangePassword changes password for a user.
func ChangePassword(ctx context.Context, client *gophercloud.ServiceClient, userID string, opts ChangePasswordOptsBuilder) (r ChangePasswordResult) {
	b, err := opts.ToUserChangePasswordMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, changePasswordURL(client, userID), &b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a user.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, userID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, userID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListGroups enumerates groups user belongs to.
func ListGroups(client *gophercloud.ServiceClient, userID string) pagination.Pager {
	url := listGroupsURL(client, userID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return groups.GroupPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// AddToGroup adds a user to a group.
func AddToGroup(ctx context.Context, client *gophercloud.ServiceClient, groupID, userID string) (r AddToGroupResult) {
	url := addToGroupURL(client, groupID, userID)
	resp, err := client.Put(ctx, url, nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// IsMemberOfGroup checks whether a user belongs to a group.
func IsMemberOfGroup(ctx context.Context, client *gophercloud.ServiceClient, groupID, userID string) (r IsMemberOfGroupResult) {
	url := isMemberOfGroupURL(client, groupID, userID)
	resp, err := client.Head(ctx, url, &gophercloud.RequestOpts{
		OkCodes: []int{204, 404},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	if r.Err == nil {
		if resp.StatusCode == 204 {
			r.isMember = true
		}
	}
	return
}

// RemoveFromGroup removes a user from a group.
func RemoveFromGroup(ctx context.Context, client *gophercloud.ServiceClient, groupID, userID string) (r RemoveFromGroupResult) {
	url := removeFromGroupURL(client, groupID, userID)
	resp, err := client.Delete(ctx, url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListProjects enumerates groups user belongs to.
func ListProjects(client *gophercloud.ServiceClient, userID string) pagination.Pager {
	url := listProjectsURL(client, userID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return projects.ProjectPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListInGroup enumerates users that belong to a group.
func ListInGroup(client *gophercloud.ServiceClient, groupID string, opts ListOptsBuilder) pagination.Pager {
	url := listInGroupURL(client, groupID)
	if opts != nil {
		query, err := opts.ToUserListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return UserPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
