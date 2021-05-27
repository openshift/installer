package mccpv2

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//ErrCodeOrgDoesnotExist ...
var ErrCodeOrgDoesnotExist = "OrgDoesnotExist"

//Metadata ...
type Metadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//Resource ...
type Resource struct {
	Metadata Metadata
}

//OrgResource ...
type OrgResource struct {
	Resource
	Entity OrgEntity
}

//OrgEntity ...
type OrgEntity struct {
	Name                   string `json:"name"`
	Region                 string `json:"region"`
	BillingEnabled         bool   `json:"billing_enabled"`
	Status                 string `json:"status"`
	OrgQuotaDefinitionGUID string `json:"quota_definition_guid"`
}

//ToFields ..
func (resource OrgResource) ToFields() Organization {
	entity := resource.Entity

	return Organization{
		GUID:                   resource.Metadata.GUID,
		Name:                   entity.Name,
		Region:                 entity.Region,
		BillingEnabled:         entity.BillingEnabled,
		Status:                 entity.Status,
		OrgQuotaDefinitionGUID: entity.OrgQuotaDefinitionGUID,
	}
}

//OrgCreateRequest ...
type OrgCreateRequest struct {
	Name                   string `json:"name"`
	OrgQuotaDefinitionGUID string `json:"quota_definition_guid,omitempty"`
	Status                 string `json:"status,omitempty"`
}

//OrgUpdateRequest ...
type OrgUpdateRequest struct {
	Name                   *string `json:"name,omitempty"`
	OrgQuotaDefinitionGUID string  `json:"quota_definition_guid,omitempty"`
}

//Organization model
type Organization struct {
	GUID                   string
	Name                   string
	Region                 string
	BillingEnabled         bool
	Status                 string
	OrgQuotaDefinitionGUID string
}

//OrganizationFields ...
type OrganizationFields struct {
	Metadata Metadata
	Entity   OrgEntity
}

//OrgRole ...
type OrgRole struct {
	UserGUID string
	Admin    bool
	UserName string
}

//OrgRoleResource ...
type OrgRoleResource struct {
	Resource
	Entity OrgRoleEntity
}

//OrgRoleEntity ...
type OrgRoleEntity struct {
	UserGUID string `json:"guid"`
	Admin    bool   `json:"bool"`
	UserName string `json:"username"`
}

//ToFields ...
func (resource *OrgRoleResource) ToFields() OrgRole {
	entity := resource.Entity

	return OrgRole{
		UserGUID: resource.Metadata.GUID,
		Admin:    entity.Admin,
		UserName: entity.UserName,
	}
}

// OrgRegionInformation is the region information associated with an org
type OrgRegionInformation struct {
	ID          string `json:"id"`
	Domain      string `json:"domain"`
	Name        string `json:"name"`
	Region      string `json:"region"`
	DisplayName string `json:"display_name"`
	Customer    struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
	} `json:"customer"`
	Deployment struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
	} `json:"deployment"`
	Geo struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
	} `json:"geo"`
	Account struct {
		GUID       string   `json:"guid"`
		OwnerGUIDs []string `json:"owner_guids"`
	} `json:"account"`
	PublicRegionsByProximity []string `json:"public_regions_by_proximity"`
	ConsoleURL               string   `json:"console_url"`
	CFAPI                    string   `json:"cf_api"`
	MCCPAPI                  string   `json:"mccp_api"`
	Type                     string   `json:"type"`
	Home                     bool     `json:"home"`
	Stealth                  string   `json:"stealth"`
	Aliases                  []string `json:"aliases"`
	Settings                 struct {
		Devops struct {
			Enabled bool `json:"enabled"`
		} `json:"devops"`
		EnhancedAutoFix bool `json:"enhancedAutofix"`
	} `json:"settings"`
	OrgName string `json:"org_name"`
	OrgGUID string `json:"org_guid"`
}

//Organizations ...
type Organizations interface {
	Create(req OrgCreateRequest, opts ...bool) (*OrganizationFields, error)
	Get(orgGUID string) (*OrganizationFields, error)
	List(region string) ([]Organization, error)
	FindByName(orgName, region string) (*Organization, error)
	DeleteByRegion(guid string, region string, opts ...bool) error
	Delete(guid string, opts ...bool) error
	Update(guid string, req OrgUpdateRequest, opts ...bool) (*OrganizationFields, error)
	GetRegionInformation(orgGUID string) ([]OrgRegionInformation, error)

	AssociateBillingManager(orgGUID string, userMail string) (*OrganizationFields, error)
	AssociateAuditor(orgGUID string, userMail string) (*OrganizationFields, error)
	AssociateManager(orgGUID string, userMail string) (*OrganizationFields, error)
	AssociateUser(orgGUID string, userMail string) (*OrganizationFields, error)

	ListBillingManager(orgGUID string, filters ...string) ([]OrgRole, error)
	ListAuditors(orgGUID string, filters ...string) ([]OrgRole, error)
	ListManager(orgGUID string, filters ...string) ([]OrgRole, error)
	ListUsers(orgGUID string, filters ...string) ([]OrgRole, error)

	DisassociateBillingManager(orgGUID string, userMail string) error
	DisassociateManager(orgGUID string, userMail string) error
	DisassociateAuditor(orgGUID string, userMail string) error
	DisassociateUser(orgGUID string, userMail string) error
}

type organization struct {
	client *client.Client
}

func newOrganizationAPI(c *client.Client) Organizations {
	return &organization{
		client: c,
	}
}

// opts is list of boolean parametes
// opts[0] - async - Will run the create request in a background job. Recommended: 'true'. Default to 'true'.

func (o *organization) Create(req OrgCreateRequest, opts ...bool) (*OrganizationFields, error) {
	async := true
	orgFields := OrganizationFields{}
	if len(opts) > 0 {
		async = opts[0]
	}
	rawURL := fmt.Sprintf("/v2/organizations?async=%t", async)
	_, err := o.client.Post(rawURL, req, &orgFields)
	if err != nil {
		return nil, err
	}
	return &orgFields, err
}

func (o *organization) Get(orgGUID string) (*OrganizationFields, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s", orgGUID)
	orgFields := OrganizationFields{}
	_, err := o.client.Get(rawURL, &orgFields)
	if err != nil {
		return nil, err
	}
	return &orgFields, err
}

// opts is list of boolean parametes
// opts[0] - async - Will run the update request in a background job. Recommended: 'true'. Default to 'true'.

func (o *organization) Update(guid string, req OrgUpdateRequest, opts ...bool) (*OrganizationFields, error) {
	async := true
	if len(opts) > 0 {
		async = opts[0]
	}
	orgFields := OrganizationFields{}
	rawURL := fmt.Sprintf("/v2/organizations/%s?async=%t", guid, async)

	_, err := o.client.Put(rawURL, req, &orgFields)
	return &orgFields, err
}

// opts is list of boolean parametes
// opts[0] - async - Will run the delete request in a background job. Recommended: 'true'. Default to 'true'.
// opts[1] - recursive - Will delete all spaces, apps, services, routes, and private domains associated with the org. Default to 'false'.
// Deprecated: Use DeleteByRegion instead.
func (o *organization) Delete(guid string, opts ...bool) error {
	async := true
	recursive := false
	if len(opts) > 0 {
		async = opts[0]
	}
	if len(opts) > 1 {
		recursive = opts[1]
	}
	rawURL := fmt.Sprintf("/v2/organizations/%s?async=%t&recursive=%t", guid, async, recursive)
	_, err := o.client.Delete(rawURL)
	return err
}

// opts is list of boolean parametes
// opts[0] - async - Will run the delete request in a background job. Recommended: 'true'. Default to 'true'.
// opts[1] - recursive - Will delete all spaces, apps, services, routes, and private domains associated with the org. Default to 'false'.
// region - specify the region where the org to be deleted. If org to be deleted in all region's pass the region as 'all'.
func (o *organization) DeleteByRegion(guid string, region string, opts ...bool) error {
	async := true
	recursive := false
	if len(opts) > 0 {
		async = opts[0]
	}
	if len(opts) > 1 {
		recursive = opts[1]
	}

	rawURL := fmt.Sprintf("/v2/organizations/%s?async=%t&recursive=%t&region=%s", guid, async, recursive, region)
	_, err := o.client.Delete(rawURL)
	return err
}

func (o *organization) List(region string) ([]Organization, error) {
	req := rest.GetRequest("/v2/organizations")
	if region != "" {
		req.Query("region", region)
	}
	path, err := o.url(req)
	if err != nil {
		return []Organization{}, err
	}

	var orgs []Organization
	err = o.listOrgResourcesWithPath(path, func(orgResource OrgResource) bool {
		orgs = append(orgs, orgResource.ToFields())
		return true
	})
	return orgs, err
}

//FindByName ...
func (o *organization) FindByName(name string, region string) (*Organization, error) {
	path, err := o.urlOfOrgWithName(name, region, false)
	if err != nil {
		return nil, err
	}

	var org Organization
	var found bool
	err = o.listOrgResourcesWithPath(path, func(orgResource OrgResource) bool {
		org = orgResource.ToFields()
		found = true
		return false
	})

	if err != nil {
		return nil, err
	}

	if found {
		return &org, err
	}

	//May not be found and no error
	return nil, bmxerror.New(ErrCodeOrgDoesnotExist,
		fmt.Sprintf("Given org %q doesn't exist in the given region %q", name, region))

}

// GetRegionInformation get the region information associated with this org.
func (o *organization) GetRegionInformation(orgGUID string) ([]OrgRegionInformation, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/regions", orgGUID)
	var regionOrgInfo []OrgRegionInformation
	_, err := o.client.Get(rawURL, &regionOrgInfo)
	if err != nil {
		return nil, err
	}
	return regionOrgInfo, nil
}

func (o *organization) listOrgResourcesWithPath(path string, cb func(OrgResource) bool) error {
	_, err := o.client.GetPaginated(path, NewCCPaginatedResources(OrgResource{}), func(resource interface{}) bool {
		if orgResource, ok := resource.(OrgResource); ok {
			return cb(orgResource)
		}
		return false
	})
	return err
}

func (o *organization) urlOfOrgWithName(name string, region string, inline bool) (string, error) {
	req := rest.GetRequest("/v2/organizations").Query("q", fmt.Sprintf("name:%s", name))
	if region != "" {
		req.Query("region", region)
	}
	if inline {
		req.Query("inline-relations-depth", "1")
	}
	return o.url(req)
}

func (o *organization) url(req *rest.Request) (string, error) {
	httpReq, err := req.Build()
	if err != nil {
		return "", err
	}
	return httpReq.URL.String(), nil
}

func (o *organization) associateOrgRole(url, userMail string) (*OrganizationFields, error) {
	orgFields := OrganizationFields{}
	_, err := o.client.Put(url, map[string]string{"username": userMail}, &orgFields)
	if err != nil {
		return nil, err
	}
	return &orgFields, nil
}

func (o *organization) removeOrgRole(url, userMail string) error {
	orgFields := OrganizationFields{}
	_, err := o.client.DeleteWithBody(url, map[string]string{"username": userMail}, &orgFields)
	return err
}
func (o *organization) AssociateBillingManager(orgGUID string, userMail string) (*OrganizationFields, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/billing_managers", orgGUID)
	return o.associateOrgRole(rawURL, userMail)

}
func (o *organization) AssociateAuditor(orgGUID string, userMail string) (*OrganizationFields, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/auditors", orgGUID)
	return o.associateOrgRole(rawURL, userMail)
}
func (o *organization) AssociateManager(orgGUID string, userMail string) (*OrganizationFields, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/managers", orgGUID)
	return o.associateOrgRole(rawURL, userMail)
}

func (o *organization) AssociateUser(orgGUID string, userMail string) (*OrganizationFields, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/users", orgGUID)
	return o.associateOrgRole(rawURL, userMail)
}

func (o *organization) DisassociateBillingManager(orgGUID string, userMail string) error {
	rawURL := fmt.Sprintf("/v2/organizations/%s/billing_managers", orgGUID)
	return o.removeOrgRole(rawURL, userMail)

}
func (o *organization) DisassociateAuditor(orgGUID string, userMail string) error {
	rawURL := fmt.Sprintf("/v2/organizations/%s/auditors", orgGUID)
	return o.removeOrgRole(rawURL, userMail)
}
func (o *organization) DisassociateManager(orgGUID string, userMail string) error {
	rawURL := fmt.Sprintf("/v2/organizations/%s/managers", orgGUID)
	return o.removeOrgRole(rawURL, userMail)
}

func (o *organization) DisassociateUser(orgGUID string, userMail string) error {
	rawURL := fmt.Sprintf("/v2/organizations/%s/users", orgGUID)
	return o.removeOrgRole(rawURL, userMail)
}

func (o *organization) listOrgRolesWithPath(path string) ([]OrgRole, error) {
	var orgRoles []OrgRole
	_, err := o.client.GetPaginated(path, NewCCPaginatedResources(OrgRoleResource{}), func(resource interface{}) bool {
		if orgRoleResource, ok := resource.(OrgRoleResource); ok {
			orgRoles = append(orgRoles, orgRoleResource.ToFields())
			return true
		}
		return false
	})
	return orgRoles, err
}
func (o *organization) listOrgRoles(rawURL string, filters ...string) ([]OrgRole, error) {
	req := rest.GetRequest(rawURL)
	if len(filters) > 0 {
		req.Query("q", strings.Join(filters, ""))
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	return o.listOrgRolesWithPath(path)
}

func (o *organization) ListBillingManager(orgGUID string, filters ...string) ([]OrgRole, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/billing_managers", orgGUID)
	return o.listOrgRoles(rawURL, filters...)
}

func (o *organization) ListManager(orgGUID string, filters ...string) ([]OrgRole, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/managers", orgGUID)
	return o.listOrgRoles(rawURL, filters...)
}

func (o *organization) ListAuditors(orgGUID string, filters ...string) ([]OrgRole, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/auditors", orgGUID)
	return o.listOrgRoles(rawURL, filters...)
}

func (o *organization) ListUsers(orgGUID string, filters ...string) ([]OrgRole, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/users", orgGUID)
	return o.listOrgRoles(rawURL, filters...)
}
