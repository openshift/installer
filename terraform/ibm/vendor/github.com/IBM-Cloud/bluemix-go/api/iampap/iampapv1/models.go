package iampapv1

import (
	"github.com/IBM-Cloud/bluemix-go/api/iampap/iampapv2"
	"github.com/IBM-Cloud/bluemix-go/models"
)

// Policy is the model of IAM PAP policy
type Policy struct {
	ID               string     `json:"id,omitempty"`
	Type             string     `json:"type"`
	Subjects         []Subject  `json:"subjects"`
	Roles            []Role     `json:"roles"`
	Resources        []Resource `json:"resources"`
	Href             string     `json:"href,omitempty"`
	CreatedAt        string     `json:"created_at,omitempty"`
	CreatedByID      string     `json:"created_by_id,omitempty"`
	LastModifiedAt   string     `json:"last_modified_at,omitempty"`
	LastModifiedByID string     `json:"last_modified_by_id,omitempty"`
	Version          string     `json:"-"`
}

// Role is the role model used by policy
type Role struct {
	RoleID      string `json:"role_id"`
	Name        string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
}

func fromModel(role models.PolicyRole) Role {
	return Role{
		RoleID: role.ID.String(),
		// When create/update, "name" and "description" are not allowed
		// Name:        role.Name,
		// Description: role.Description,
	}
}

// ConvertRoleModels will transform role models returned from "/v1/roles" to the model used by policy
func ConvertRoleModels(roles []models.PolicyRole) []Role {
	results := make([]Role, len(roles))
	for i, r := range roles {
		results[i] = fromModel(r)
	}
	return results
}

// ConvertV2RoleModels will transform role models returned from "/v2/roles" to the model used by policy
func ConvertV2RoleModels(roles []iampapv2.Role) []Role {
	results := make([]Role, len(roles))
	for i, r := range roles {
		results[i] = Role{
			RoleID: r.Crn,
		}
	}
	return results
}

// Subject is the target to which is assigned policy
type Subject struct {
	Attributes []Attribute `json:"attributes"`
}

const (
	AccessGroupIDAttribute   = "accesGroupId"
	AccountIDAttribute       = "accountId"
	OrganizationIDAttribute  = "organizationId"
	SpaceIDAttribute         = "spaceId"
	RegionAttribute          = "region"
	ServiceTypeAttribute     = "serviceType"
	ServiceNameAttribute     = "serviceName"
	ServiceInstanceAttribute = "serviceInstance"
	ResourceTypeAttribute    = "resourceType"
	ResourceAttribute        = "resource"
	ResourceGroupIDAttribute = "resourceGroupId"
)

// GetAttribute returns an attribute of policy subject
func (s *Subject) GetAttribute(name string) string {
	for _, a := range s.Attributes {
		if a.Name == name {
			return a.Value
		}
	}
	return ""
}

// SetAttribute sets value of an attribute of policy subject
func (s *Subject) SetAttribute(name string, value string) {
	for _, a := range s.Attributes {
		if a.Name == name {
			a.Value = value
			return
		}
	}
	s.Attributes = append(s.Attributes, Attribute{
		Name:  name,
		Value: value,
	})
}

// AccessGroupID returns access group ID attribute of policy subject if exists
func (s *Subject) AccessGroupID() string {
	return s.GetAttribute("access_group_id")
}

// AccountID returns account ID attribute of policy subject if exists
func (s *Subject) AccountID() string {
	return s.GetAttribute("accountId")
}

// IAMID returns IAM ID attribute of policy subject if exists
func (s *Subject) IAMID() string {
	return s.GetAttribute("iam_id")
}

// ServiceName returns service name attribute of policy subject if exists
func (s *Subject) ServiceName() string {
	return s.GetAttribute("serviceName")
}

// ServiceInstance returns service instance attribute of policy subject if exists
func (s *Subject) ServiceInstance() string {
	return s.GetAttribute("serviceInstance")
}

// ResourceType returns resource type of the policy subject if exists
func (s *Subject) ResourceType() string {
	return s.GetAttribute("resourceType")
}

// ResourceGroupID returns resource group ID attribute of policy resource if exists
func (s *Subject) ResourceGroupID() string {
	return s.GetAttribute(ResourceGroupIDAttribute)
}

// SetAccessGroupID sets value of access group ID attribute of policy subject
func (s *Subject) SetAccessGroupID(value string) {
	s.SetAttribute("access_group_id", value)
}

// SetAccountID sets value of account ID attribute of policy subject
func (s *Subject) SetAccountID(value string) {
	s.SetAttribute("accountId", value)
}

// SetIAMID sets value of IAM ID attribute of policy subject
func (s *Subject) SetIAMID(value string) {
	s.SetAttribute("iam_id", value)
}

// SetServiceName sets value of service name attribute of policy subject
func (s *Subject) SetServiceName(value string) {
	s.SetAttribute("serviceName", value)
}

// SetServiceInstance sets value of service instance attribute of policy subject
func (s *Subject) SetServiceInstance(value string) {
	s.SetAttribute("serviceInstance", value)
}

// SetResourceType sets value of resource type attribute of policy subject
func (s *Subject) SetResourceType(value string) {
	s.SetAttribute("resourceType", value)
}

// SetResourceGroupID sets value of resource group ID attribute of policy resource
func (s *Subject) SetResourceGroupID(value string) {
	s.SetAttribute(ResourceGroupIDAttribute, value)
}

// Resource is the object controlled by the policy
type Resource struct {
	Attributes []Attribute `json:"attributes"`
}

// GetAttribute returns an attribute of policy resource
func (r *Resource) GetAttribute(name string) string {
	for _, a := range r.Attributes {
		if a.Name == name {
			return a.Value
		}
	}
	return ""
}

// SetAttribute sets value of an attribute of policy resource
func (r *Resource) SetAttribute(name string, value string) {
	for _, a := range r.Attributes {
		if a.Name == name {
			a.Value = value
			return
		}
	}
	r.Attributes = append(r.Attributes, Attribute{
		Name:  name,
		Value: value,
	})
}

// AccessGroupID returns access group ID attribute of policy resource if exists
func (r *Resource) AccessGroupID() string {
	return r.GetAttribute(AccessGroupIDAttribute)
}

// AccountID returns account ID attribute of policy resource if exists
func (r *Resource) AccountID() string {
	return r.GetAttribute(AccountIDAttribute)
}

// OrganizationID returns organization ID attribute of policy resource if exists
func (r *Resource) OrganizationID() string {
	return r.GetAttribute(OrganizationIDAttribute)
}

// Region returns region attribute of policy resource if exists
func (r *Resource) Region() string {
	return r.GetAttribute(RegionAttribute)
}

// Resource returns resource attribute of policy resource if exists
func (r *Resource) Resource() string {
	return r.GetAttribute(ResourceAttribute)
}

// ResourceType returns resource type attribute of policy resource if exists
func (r *Resource) ResourceType() string {
	return r.GetAttribute(ResourceTypeAttribute)
}

// ResourceGroupID returns resource group ID attribute of policy resource if exists
func (r *Resource) ResourceGroupID() string {
	return r.GetAttribute(ResourceGroupIDAttribute)
}

// ServiceName returns service name attribute of policy resource if exists
func (r *Resource) ServiceName() string {
	return r.GetAttribute(ServiceNameAttribute)
}

// ServiceInstance returns service instance attribute of policy resource if exists
func (r *Resource) ServiceInstance() string {
	return r.GetAttribute(ServiceInstanceAttribute)
}

// SpaceID returns space ID attribute of policy resource if exists
func (r *Resource) SpaceID() string {
	return r.GetAttribute(SpaceIDAttribute)
}

// ServiceType returns service type attribute of policy resource if exists
func (r *Resource) ServiceType() string {
	return r.GetAttribute(ServiceTypeAttribute)
}

// CustomAttributes will return all attributes which are not system defined
func (r *Resource) CustomAttributes() []Attribute {
	attributes := []Attribute{}
	for _, a := range r.Attributes {
		switch a.Name {
		case AccessGroupIDAttribute:
		case AccountIDAttribute:
		case OrganizationIDAttribute:
		case SpaceIDAttribute:
		case RegionAttribute:
		case ResourceAttribute:
		case ResourceTypeAttribute:
		case ResourceGroupIDAttribute:
		case ServiceTypeAttribute:
		case ServiceNameAttribute:
		case ServiceInstanceAttribute:
		default:
			attributes = append(attributes, a)
		}
	}
	return attributes
}

// SetAccessGroupID sets value of access group ID attribute of policy resource
func (r *Resource) SetAccessGroupID(value string) {
	r.SetAttribute(AccessGroupIDAttribute, value)
}

// SetAccountID sets value of account ID attribute of policy resource
func (r *Resource) SetAccountID(value string) {
	r.SetAttribute(AccountIDAttribute, value)
}

// SetOrganizationID sets value of organization ID attribute of policy resource
func (r *Resource) SetOrganizationID(value string) {
	r.SetAttribute(OrganizationIDAttribute, value)
}

// SetRegion sets value of region attribute of policy resource
func (r *Resource) SetRegion(value string) {
	r.SetAttribute(RegionAttribute, value)
}

// SetResource sets value of resource attribute of policy resource
func (r *Resource) SetResource(value string) {
	r.SetAttribute(ResourceAttribute, value)
}

// SetResourceType sets value of resource type attribute of policy resource
func (r *Resource) SetResourceType(value string) {
	r.SetAttribute(ResourceTypeAttribute, value)
}

// SetResourceGroupID sets value of resource group ID attribute of policy resource
func (r *Resource) SetResourceGroupID(value string) {
	r.SetAttribute(ResourceGroupIDAttribute, value)
}

// SetServiceName sets value of service name attribute of policy resource
func (r *Resource) SetServiceName(value string) {
	r.SetAttribute(ServiceNameAttribute, value)
}

// SetServiceInstance sets value of service instance attribute of policy resource
func (r *Resource) SetServiceInstance(value string) {
	r.SetAttribute(ServiceInstanceAttribute, value)
}

// SetSpaceID sets value of space ID attribute of policy resource
func (r *Resource) SetSpaceID(value string) {
	r.SetAttribute("spaceID", value)
}

// SetServiceType sets value of service type attribute of policy resource
func (r *Resource) SetServiceType(value string) {
	r.SetAttribute(ServiceTypeAttribute, value)
}

// Attribute is part of policy subject and resource
type Attribute struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Operator string `json:"operator,omitempty"`
}
