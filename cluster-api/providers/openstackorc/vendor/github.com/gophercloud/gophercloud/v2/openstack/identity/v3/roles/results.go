package roles

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Role grants permissions to a user.
type Role struct {
	// DomainID is the domain ID the role belongs to.
	DomainID string `json:"domain_id"`

	// ID is the unique ID of the role.
	ID string `json:"id"`

	// Links contains referencing links to the role.
	Links map[string]any `json:"links"`

	// Name is the role name
	Name string `json:"name"`

	// Description is the description of the role.
	Description string `json:"description"`

	// Extra is a collection of miscellaneous key/values.
	Extra map[string]any `json:"-"`

	// Options are a set of defined options that allow certain features for a role
	Options map[Option]any `json:"options"`
}

func (r *Role) UnmarshalJSON(b []byte) error {
	type tmp Role
	var s struct {
		tmp
		Extra map[string]any `json:"extra"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Role(s.tmp)

	// Collect other fields and bundle them into Extra
	// but only if a field titled "extra" wasn't sent.
	if s.Extra != nil {
		r.Extra = s.Extra
	} else {
		var result any
		err := json.Unmarshal(b, &result)
		if err != nil {
			return err
		}
		if resultMap, ok := result.(map[string]any); ok {
			r.Extra = gophercloud.RemainingKeys(Role{}, resultMap)

			// the following code is required for backward compatibility with the
			// old behavior, when description was in extra
			if description, ok := resultMap["description"]; ok {
				r.Extra["description"] = description
			}
		}
	}

	return err
}

type roleResult struct {
	gophercloud.Result
}

// GetResult is the response from a Get operation. Call its Extract method
// to interpret it as a Role.
type GetResult struct {
	roleResult
}

// CreateResult is the response from a Create operation. Call its Extract method
// to interpret it as a Role
type CreateResult struct {
	roleResult
}

// UpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a Role.
type UpdateResult struct {
	roleResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr to
// determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// RolePage is a single page of Role results.
type RolePage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of Roles contains any results.
func (r RolePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	roles, err := ExtractRoles(r)
	return len(roles) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r RolePage) NextPageURL() (string, error) {
	var s struct {
		Links struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Links.Next, err
}

// ExtractProjects returns a slice of Roles contained in a single page of
// results.
func ExtractRoles(r pagination.Page) ([]Role, error) {
	var s struct {
		Roles []Role `json:"roles"`
	}
	err := (r.(RolePage)).ExtractInto(&s)
	return s.Roles, err
}

// Extract interprets any roleResults as a Role.
func (r roleResult) Extract() (*Role, error) {
	var s struct {
		Role *Role `json:"role"`
	}
	err := r.ExtractInto(&s)
	return s.Role, err
}

// RoleAssignment is the result of a role assignments query.
type RoleAssignment struct {
	Role  AssignedRole `json:"role,omitempty"`
	Scope Scope        `json:"scope,omitempty"`
	User  User         `json:"user,omitempty"`
	Group Group        `json:"group,omitempty"`
}

// AssignedRole represents a Role in an assignment.
type AssignedRole struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Scope represents a scope in a Role assignment.
type Scope struct {
	Domain  Domain  `json:"domain,omitempty"`
	Project Project `json:"project,omitempty"`
}

// Domain represents a domain in a role assignment scope.
type Domain struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Project represents a project in a role assignment scope.
type Project struct {
	Domain Domain `json:"domain,omitempty"`
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
}

// User represents a user in a role assignment scope.
type User struct {
	Domain Domain `json:"domain,omitempty"`
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
}

// Group represents a group in a role assignment scope.
type Group struct {
	Domain Domain `json:"domain,omitempty"`
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
}

// RoleAssignmentPage is a single page of RoleAssignments results.
type RoleAssignmentPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if the RoleAssignmentPage contains no results.
func (r RoleAssignmentPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	roleAssignments, err := ExtractRoleAssignments(r)
	return len(roleAssignments) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to
// the next page of results.
func (r RoleAssignmentPage) NextPageURL() (string, error) {
	var s struct {
		Links struct {
			Next string `json:"next"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	return s.Links.Next, err
}

// ExtractRoleAssignments extracts a slice of RoleAssignments from a Collection
// acquired from List.
func ExtractRoleAssignments(r pagination.Page) ([]RoleAssignment, error) {
	var s struct {
		RoleAssignments []RoleAssignment `json:"role_assignments"`
	}
	err := (r.(RoleAssignmentPage)).ExtractInto(&s)
	return s.RoleAssignments, err
}

// AssignmentResult represents the result of an assign operation.
// Call ExtractErr method to determine if the request succeeded or failed.
type AssignmentResult struct {
	gophercloud.ErrResult
}

// UnassignmentResult represents the result of an unassign operation.
// Call ExtractErr method to determine if the request succeeded or failed.
type UnassignmentResult struct {
	gophercloud.ErrResult
}

type impliedRoleResult struct {
	gophercloud.Result
}

// ImpliedRoleResult is the result of an PUT request. Call its Extract method to
// interpret it as a roleInference.
type CreateImpliedRoleResult struct {
	impliedRoleResult
}

type GetImpliedRoleResult struct {
	impliedRoleResult
}
type PriorRole struct {
	// ID contains the ID of the role in a prior_role object.
	ID string `json:"id,omitempty"`
	// Name contains the name of a role in a prior_role object.
	Name string `json:"name,omitempty"`
	// Links contains referencing links to the  prior_role.
	Links map[string]any `json:"links"`
}

type ImpliedRole struct {
	// ID contains the ID of the role in an implied_role object.
	ID string `json:"id,omitempty"`
	// Name contains the name of role  in an implied_role.
	Name string `json:"name,omitempty"`
	// Links contains referencing links to the implied_role.
	Links map[string]any `json:"links"`
}

type RoleInference struct {
	// PriorRole is the role object that implies a list of implied_role objects.
	PriorRole PriorRole `json:"prior_role"`
	// Implies is an array of implied_role objects implied by a prior_role object.
	ImpliedRole ImpliedRole `json:"implies"`
}

type RoleInferenceRule struct {
	RoleInference RoleInference  `json:"role_inference"`
	Links         map[string]any `json:"links"`
}

func (r impliedRoleResult) Extract() (*RoleInferenceRule, error) {
	var s = &RoleInferenceRule{}
	err := r.ExtractInto(s)
	return s, err
}

type ListImpliedRolesResult struct {
	gophercloud.Result
}

type ImpliedRoleObject struct {
	// ID contains the ID of the role in an implied_role object.
	ID string `json:"id,omitempty"`
	// Name contains the name of role  in an implied_role.
	Name string `json:"name,omitempty"`
	// Name contains the name of role  in an implied_role.
	Description string `json:"description,omitempty"`
	// Links contains referencing links to the implied_role.
	Links map[string]any `json:"links"`
}

type PriorRoleObject struct {
	// ID contains the ID of the role in an implied_role object.
	ID string `json:"id,omitempty"`
	// Name contains the name of role  in an implied_role.
	Name string `json:"name,omitempty"`
	// Name contains the name of role  in an implied_role.
	Description string `json:"description,omitempty"`
	// Links contains referencing links to the implied_role.
	Links map[string]any `json:"links"`
}
type RoleInferenceRules struct {
	// PriorRole is the role object that implies a list of implied_role objects.
	PriorRole PriorRoleObject `json:"prior_role"`
	// Implies is an array of implied_role objects implied by a prior_role object.
	ImpliedRoles []ImpliedRoleObject `json:"implies"`
}

type RoleInferenceRuleList struct {
	RoleInferenceRuleList []RoleInferenceRules `json:"role_inferences"`
	Links                 map[string]any       `json:"links"`
}

func (r ListImpliedRolesResult) Extract() (*RoleInferenceRuleList, error) {
	var s = &RoleInferenceRuleList{}
	err := r.ExtractInto(s)
	return s, err
}

type DeleteImpliedRoleResult struct {
	gophercloud.ErrResult
}
