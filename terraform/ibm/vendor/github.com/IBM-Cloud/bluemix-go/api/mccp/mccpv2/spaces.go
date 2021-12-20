package mccpv2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//SpaceCreateRequest ...
type SpaceCreateRequest struct {
	Name           string `json:"name"`
	OrgGUID        string `json:"organization_guid"`
	SpaceQuotaGUID string `json:"space_quota_definition_guid,omitempty"`
}

//SpaceUpdateRequest ...
type SpaceUpdateRequest struct {
	Name *string `json:"name,omitempty"`
}

//Space ...
type Space struct {
	GUID           string
	Name           string
	OrgGUID        string
	SpaceQuotaGUID string
	AllowSSH       bool
}

//SpaceRole ...
type SpaceRole struct {
	UserGUID string
	Admin    bool
	UserName string
}

//SpaceFields ...
type SpaceFields struct {
	Metadata SpaceMetadata
	Entity   SpaceEntity
}

//SpaceMetadata ...
type SpaceMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//ErrCodeSpaceDoesnotExist ...
const ErrCodeSpaceDoesnotExist = "SpaceDoesnotExist"

//SpaceResource ...
type SpaceResource struct {
	Resource
	Entity SpaceEntity
}

//SpaceRoleResource ...
type SpaceRoleResource struct {
	Resource
	Entity SpaceRoleEntity
}

//SpaceRoleEntity ...
type SpaceRoleEntity struct {
	UserGUID string `json:"guid"`
	Admin    bool   `json:"bool"`
	UserName string `json:"username"`
}

//SpaceEntity ...
type SpaceEntity struct {
	Name           string `json:"name"`
	OrgGUID        string `json:"organization_guid"`
	SpaceQuotaGUID string `json:"space_quota_definition_guid"`
	AllowSSH       bool   `json:"allow_ssh"`
}

//ToFields ...
func (resource *SpaceResource) ToFields() Space {
	entity := resource.Entity

	return Space{
		GUID:           resource.Metadata.GUID,
		Name:           entity.Name,
		OrgGUID:        entity.OrgGUID,
		SpaceQuotaGUID: entity.SpaceQuotaGUID,
		AllowSSH:       entity.AllowSSH,
	}
}

//ToFields ...
func (resource *SpaceRoleResource) ToFields() SpaceRole {
	entity := resource.Entity

	return SpaceRole{
		UserGUID: resource.Metadata.GUID,
		Admin:    entity.Admin,
		UserName: entity.UserName,
	}
}

//RouteFilter ...
type RouteFilter struct {
	DomainGUID string
	Host       *string
	Path       *string
	Port       *int
}

//Spaces ...
type Spaces interface {
	ListSpacesInOrg(orgGUID, region string) ([]Space, error)
	FindByNameInOrg(orgGUID, name, region string) (*Space, error)
	Create(req SpaceCreateRequest, opts ...bool) (*SpaceFields, error)
	Update(spaceGUID string, req SpaceUpdateRequest, opts ...bool) (*SpaceFields, error)
	Delete(spaceGUID string, opts ...bool) error
	Get(spaceGUID string) (*SpaceFields, error)
	ListRoutes(spaceGUID string, req RouteFilter) ([]Route, error)
	AssociateAuditor(spaceGUID, userMail string) (*SpaceFields, error)
	AssociateDeveloper(spaceGUID, userMail string) (*SpaceFields, error)
	AssociateManager(spaceGUID, userMail string) (*SpaceFields, error)

	DisassociateAuditor(spaceGUID, userMail string) error
	DisassociateDeveloper(spaceGUID, userMail string) error
	DisassociateManager(spaceGUID, userMail string) error

	ListAuditors(spaceGUID string, filters ...string) ([]SpaceRole, error)
	ListDevelopers(spaceGUID string, filters ...string) ([]SpaceRole, error)
	ListManagers(spaceGUID string, filters ...string) ([]SpaceRole, error)
}

type spaces struct {
	client *client.Client
}

func newSpacesAPI(c *client.Client) Spaces {
	return &spaces{
		client: c,
	}
}

func (r *spaces) FindByNameInOrg(orgGUID string, name string, region string) (*Space, error) {
	rawURL := fmt.Sprintf("/v2/organizations/%s/spaces", orgGUID)
	req := rest.GetRequest(rawURL).Query("q", "name:"+name)
	if region != "" {
		req.Query("region", region)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()

	spaces, err := r.listSpacesWithPath(path)

	if err != nil {
		return nil, err
	}
	if len(spaces) == 0 {
		return nil, bmxerror.New(ErrCodeSpaceDoesnotExist,
			fmt.Sprintf("Given space:  %q doesn't exist in given org: %q in the given region %q", name, orgGUID, region))

	}
	return &spaces[0], nil
}

func (r *spaces) ListSpacesInOrg(orgGUID string, region string) ([]Space, error) {
	rawURL := fmt.Sprintf("v2/organizations/%s/spaces", orgGUID)
	req := rest.GetRequest(rawURL)
	if region != "" {
		req.Query("region", region)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()

	return r.listSpacesWithPath(path)
}

func (r *spaces) listSpacesWithPath(path string) ([]Space, error) {
	var spaces []Space
	_, err := r.client.GetPaginated(path, NewCCPaginatedResources(SpaceResource{}), func(resource interface{}) bool {
		if spaceResource, ok := resource.(SpaceResource); ok {
			spaces = append(spaces, spaceResource.ToFields())
			return true
		}
		return false
	})
	return spaces, err
}

func (r *spaces) listSpaceRolesWithPath(path string) ([]SpaceRole, error) {
	var spaceRoles []SpaceRole
	_, err := r.client.GetPaginated(path, NewCCPaginatedResources(SpaceRoleResource{}), func(resource interface{}) bool {
		if spaceRoleResource, ok := resource.(SpaceRoleResource); ok {
			spaceRoles = append(spaceRoles, spaceRoleResource.ToFields())
			return true
		}
		return false
	})
	return spaceRoles, err
}

// opts is list of boolean parametes
// opts[0] - async - Will run the create request in a background job. Recommended: 'true'. Default to 'true'.

func (r *spaces) Create(req SpaceCreateRequest, opts ...bool) (*SpaceFields, error) {
	async := true
	if len(opts) > 0 {
		async = opts[0]
	}
	rawURL := fmt.Sprintf("/v2/spaces?async=%t", async)
	spaceFields := SpaceFields{}
	_, err := r.client.Post(rawURL, req, &spaceFields)
	if err != nil {
		return nil, err
	}
	return &spaceFields, nil
}

func (r *spaces) Get(spaceGUID string) (*SpaceFields, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s", spaceGUID)
	spaceFields := SpaceFields{}
	_, err := r.client.Get(rawURL, &spaceFields)
	if err != nil {
		return nil, err
	}

	return &spaceFields, err
}

// opts is list of boolean parametes
// opts[0] - async - Will run the update request in a background job. Recommended: 'true'. Default to 'true'.

func (r *spaces) Update(spaceGUID string, req SpaceUpdateRequest, opts ...bool) (*SpaceFields, error) {
	async := true
	if len(opts) > 0 {
		async = opts[0]
	}
	rawURL := fmt.Sprintf("/v2/spaces/%s?async=%t", spaceGUID, async)
	spaceFields := SpaceFields{}
	_, err := r.client.Put(rawURL, req, &spaceFields)
	if err != nil {
		return nil, err
	}
	return &spaceFields, nil
}

// opts is list of boolean parametes
// opts[0] - async - Will run the delete request in a background job. Recommended: 'true'. Default to 'true'.
// opts[1] - recursive - Will delete all apps, services, routes, and service brokers associated with the space. Default to 'false'.

func (r *spaces) Delete(spaceGUID string, opts ...bool) error {
	async := true
	recursive := false
	if len(opts) > 0 {
		async = opts[0]
	}
	if len(opts) > 1 {
		recursive = opts[1]
	}
	rawURL := fmt.Sprintf("/v2/spaces/%s?async=%t&recursive=%t", spaceGUID, async, recursive)
	_, err := r.client.Delete(rawURL)
	return err
}

func (r *spaces) associateRole(url, userMail string) (*SpaceFields, error) {
	spaceFields := SpaceFields{}
	_, err := r.client.Put(url, map[string]string{"username": userMail}, &spaceFields)
	if err != nil {
		return nil, err
	}
	return &spaceFields, nil
}

func (r *spaces) removeRole(url, userMail string) error {
	spaceFields := SpaceFields{}
	_, err := r.client.DeleteWithBody(url, map[string]string{"username": userMail}, &spaceFields)
	return err
}

func (r *spaces) AssociateManager(spaceGUID string, userMail string) (*SpaceFields, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/managers", spaceGUID)
	return r.associateRole(rawURL, userMail)
}
func (r *spaces) AssociateDeveloper(spaceGUID string, userMail string) (*SpaceFields, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/developers", spaceGUID)
	return r.associateRole(rawURL, userMail)
}
func (r *spaces) AssociateAuditor(spaceGUID string, userMail string) (*SpaceFields, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/auditors", spaceGUID)
	return r.associateRole(rawURL, userMail)
}

func (r *spaces) DisassociateManager(spaceGUID string, userMail string) error {
	rawURL := fmt.Sprintf("/v2/spaces/%s/managers", spaceGUID)
	return r.removeRole(rawURL, userMail)
}

func (r *spaces) DisassociateDeveloper(spaceGUID string, userMail string) error {
	rawURL := fmt.Sprintf("/v2/spaces/%s/developers", spaceGUID)
	return r.removeRole(rawURL, userMail)
}
func (r *spaces) DisassociateAuditor(spaceGUID string, userMail string) error {
	rawURL := fmt.Sprintf("/v2/spaces/%s/auditors", spaceGUID)
	return r.removeRole(rawURL, userMail)
}

func (r *spaces) listSpaceRoles(rawURL string, filters ...string) ([]SpaceRole, error) {
	req := rest.GetRequest(rawURL)
	if len(filters) > 0 {
		req.Query("q", strings.Join(filters, ""))
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	return r.listSpaceRolesWithPath(path)
}

func (r *spaces) ListAuditors(spaceGUID string, filters ...string) ([]SpaceRole, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/auditors", spaceGUID)
	return r.listSpaceRoles(rawURL, filters...)
}

func (r *spaces) ListManagers(spaceGUID string, filters ...string) ([]SpaceRole, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/managers", spaceGUID)
	return r.listSpaceRoles(rawURL, filters...)
}
func (r *spaces) ListDevelopers(spaceGUID string, filters ...string) ([]SpaceRole, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/developers", spaceGUID)
	return r.listSpaceRoles(rawURL, filters...)
}

func (r *spaces) ListRoutes(spaceGUID string, routeFilter RouteFilter) ([]Route, error) {
	rawURL := fmt.Sprintf("/v2/spaces/%s/routes", spaceGUID)
	req := rest.GetRequest(rawURL)
	var query string
	if routeFilter.DomainGUID != "" {
		query = "domain_guid:" + routeFilter.DomainGUID + ";"
	}
	if routeFilter.Host != nil {
		query += "host:" + *routeFilter.Host + ";"
	}
	if routeFilter.Path != nil {
		query += "path:" + *routeFilter.Path + ";"
	}
	if routeFilter.Port != nil {
		query += "port:" + strconv.Itoa(*routeFilter.Port) + ";"
	}

	if len(query) > 0 {
		req.Query("q", query)
	}

	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	route, err := listRouteWithPath(r.client, path)
	if err != nil {
		return nil, err
	}
	return route, nil
}
