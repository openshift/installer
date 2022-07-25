package functions

import "github.com/go-openapi/strfmt"

// GetNamespacesOptions : The GetNamespaces options.
type GetNamespacesOptions struct {

	// The maximum number of namespaces to return. Default 100. Maximum 200.
	Limit *int64 `json:"limit,omitempty"`

	// The number of namespaces to skip. Default 0.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers to be GDPR compliant
	Headers map[string]string
}

// NamespaceResponse : NamespaceResponse - create/get response.
type NamespaceResponse struct {

	// Time the API key was activated.
	APIKeyCreated *strfmt.DateTime `json:"API_key_created,omitempty"`

	// ID of API key used by the namespace.
	APIKeyID *string `json:"API_key_id,omitempty"`

	// CF space GUID of classic namespace - present if it is or was a classic namespace.
	ClassicSpaceguid *string `json:"classic_spaceguid,omitempty"`

	// ClassicType <br/> This attribute will be absent for an IAM namespace, a namespace which is IAM-enabled and not
	// associated with any CF space. <br/> 1 : Classic - A namespace which is associated with a CF space.  <br/> Such
	// namespace is NOT IAM-enabled and can only be used by using the legacy API key ('entitlement key'). <br/> 2 : Classic
	// IAM enabled - A namespace which is associated with a CF space and which is IAM-enabled.  <br/> It accepts IMA token
	// and legacy API key ('entitlement key') for authorization.<br/> 3 : IAM migration complete - A namespace which was/is
	// associated with a CF space, which is IAM-enabled.  <br/> It accepts only an IAM token for authorization.<br/>.
	ClassicType *int64 `json:"classic_type,omitempty"`

	// CRN of namespace - absent if namespace is NOT IAM-enabled.
	Crn *string `json:"crn,omitempty"`

	// Description - absent if namespace is NOT IAM-enabled.
	Description *string `json:"description,omitempty"`

	// UUID of namespace.
	ID *string `json:"id" validate:"required"`

	// Location of the resource.
	Location *string `json:"location" validate:"required"`

	// Name - absent if namespace is NOT IAM-enabled.
	Name *string `json:"name,omitempty"`

	// Resourceplanid used - absent if namespace is NOT IAM-enabled.
	ResourcePlanID *string `json:"resource_plan_id,omitempty"`

	// Resourcegrpid used - absent if namespace is NOT IAM-enabled.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// Serviceid used by the namespace - absent if namespace is NOT IAM-enabled.
	ServiceID *string `json:"service_id,omitempty"`

	// Key used by the cf based namespace.
	Key string `json:"key,omitempty"`

	// UUID used by the cf based namespace.
	UUID string `json:"uuid,omitempty"`
}

// NamespaceResponseList : NamespaceResponseList -.
type NamespaceResponseList struct {

	// Maximum number of namespaces to return.
	Limit *int64 `json:"limit" validate:"required"`

	// List of namespaces.
	Namespaces []NamespaceResponse `json:"namespaces" validate:"required"`

	// Number of namespaces to skip.
	Offset *int64 `json:"offset" validate:"required"`

	// Total number of namespaces available.
	TotalCount *int64 `json:"total_count" validate:"required"`
}

// CreateNamespaceOptions : The CreateNamespace options.
type CreateNamespaceOptions struct {

	// Name.
	Name *string `json:"name" validate:"required"`

	// Resourcegroupid of resource group the namespace resource should be placed in. Use 'ibmcloud resource groups' to
	// query your resources groups and their ids.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Resourceplanid to use, e.g. 'functions-base-plan'.
	ResourcePlanID *string `json:"resource_plan_id" validate:"required"`

	// Description.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers to be GDPR compliant
	Headers map[string]string
}

// GetNamespaceOptions : The GetNamespace options.
type GetNamespaceOptions struct {

	// The id of the namespace to retrieve.
	ID *string `json:"id" validate:"required"`

	// Allows users to set headers to be GDPR compliant
	Headers map[string]string
}

// DeleteNamespaceOptions : The DeleteNamespace options.
type DeleteNamespaceOptions struct {

	// The id of the namespace to delete.
	ID *string `json:"id" validate:"required"`

	// Allows users to set headers to be GDPR compliant
	Headers map[string]string
}

// UpdateNamespaceOptions : The UpdateNamespace options.
type UpdateNamespaceOptions struct {

	// The id of the namespace to update.
	ID *string `json:"id" validate:"required"`

	// New description.
	Description *string `json:"description,omitempty"`

	// New name.
	Name *string `json:"name,omitempty"`

	// Allows users to set headers to be GDPR compliant
	Headers map[string]string
}

//NamespaceResource ..
type NamespaceResource interface {
	GetID() string
	GetLocation() string
	GetName() string
	GetUUID() string
	GetKey() string
	IsIamEnabled() bool
	IsCf() bool
}

//GetID ..
func (ns *NamespaceResponse) GetID() string {
	return *ns.ID
}

//GetName ..
func (ns *NamespaceResponse) GetName() string {
	// Classic support - if no name included in namespace obj return the ID (classic namespace name)
	if ns.Name != nil {
		return *ns.Name
	}
	return *ns.ID
}

//GetKey ..
func (ns *NamespaceResponse) GetKey() string {
	return ns.Key
}

//GetUUID ..
func (ns *NamespaceResponse) GetUUID() string {
	return ns.UUID
}

//GetLocation ..
func (ns *NamespaceResponse) GetLocation() string {
	return *ns.Location
}

//IsCf ..
func (ns *NamespaceResponse) IsCf() bool {
	var iscf bool = false
	if ns.ClassicType != nil {
		iscf = (*ns.ClassicType == NamespaceTypeCFBased)
	}
	return iscf
}

//IsIamEnabled ..
func (ns *NamespaceResponse) IsIamEnabled() bool {
	// IAM support - classic_type field is not included for new IAM namespaces so always return true if nil
	if ns.ClassicType == nil {
		return true
	}
	return false
}

//IsMigrated ..
func (ns *NamespaceResponse) IsMigrated() bool {
	if *ns.ClassicType == NamespaceTypeIamMigrated {
		return true
	}
	return false
}
