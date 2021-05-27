package mccpv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//ErrCodeRouteDoesnotExist ...
var ErrCodeRouteDoesnotExist = "RouteDoesnotExist"

//RouteRequest ...
type RouteRequest struct {
	Host       string `json:"host,omitempty"`
	SpaceGUID  string `json:"space_guid"`
	DomainGUID string `json:"domain_guid,omitempty"`
	Path       string `json:"path,omitempty"`
	Port       *int   `json:"port,omitempty"`
}

//RouteUpdateRequest ...
type RouteUpdateRequest struct {
	Host *string `json:"host,omitempty"`
	Path *string `json:"path,omitempty"`
	Port *int    `json:"port,omitempty"`
}

//RouteMetadata ...
type RouteMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//RouteEntity ...
type RouteEntity struct {
	Host                string `json:"host"`
	Path                string `json:"path"`
	DomainGUID          string `json:"domain_guid"`
	SpaceGUID           string `json:"space_guid"`
	ServiceInstanceGUID string `json:"service_instance_guid"`
	Port                *int   `json:"port"`
	DomainURL           string `json:"domain_url"`
	SpaceURL            string `json:"space_url"`
	AppsURL             string `json:"apps_url"`
	RouteMappingURL     string `json:"route_mapping_url"`
}

//RouteResource ...
type RouteResource struct {
	Resource
	Entity RouteEntity
}

//RouteFields ...
type RouteFields struct {
	Metadata RouteMetadata
	Entity   RouteEntity
}

//ToFields ..
func (resource RouteResource) ToFields() Route {
	entity := resource.Entity

	return Route{
		GUID:                resource.Metadata.GUID,
		Host:                entity.Host,
		Path:                entity.Path,
		DomainGUID:          entity.DomainGUID,
		SpaceGUID:           entity.SpaceGUID,
		ServiceInstanceGUID: entity.ServiceInstanceGUID,
		Port:                entity.Port,
		DomainURL:           entity.DomainURL,
		SpaceURL:            entity.SpaceURL,
		AppsURL:             entity.AppsURL,
		RouteMappingURL:     entity.RouteMappingURL,
	}
}

//Route model
type Route struct {
	GUID                string
	Host                string
	Path                string
	DomainGUID          string
	SpaceGUID           string
	ServiceInstanceGUID string
	Port                *int
	DomainURL           string
	SpaceURL            string
	AppsURL             string
	RouteMappingURL     string
}

//Routes ...
type Routes interface {
	Find(hostname, domainGUID string) ([]Route, error)
	Create(req RouteRequest, opts ...bool) (*RouteFields, error)
	Get(routeGUID string) (*RouteFields, error)
	Update(routeGUID string, req RouteUpdateRequest, opts ...bool) (*RouteFields, error)
	Delete(routeGUID string, opts ...bool) error
}

type route struct {
	client *client.Client
}

func newRouteAPI(c *client.Client) Routes {
	return &route{
		client: c,
	}
}

func (r *route) Get(routeGUID string) (*RouteFields, error) {
	rawURL := fmt.Sprintf("/v2/routes/%s", routeGUID)
	routeFields := RouteFields{}
	_, err := r.client.Get(rawURL, &routeFields, nil)
	if err != nil {
		return nil, err
	}
	return &routeFields, nil
}

func (r *route) Find(hostname, domainGUID string) ([]Route, error) {
	rawURL := "/v2/routes?inline-relations-depth=1"
	req := rest.GetRequest(rawURL).Query("q", "host:"+hostname+";domain_guid:"+domainGUID)
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

// opts is list of boolean parametes
// opts[0] - async - Will run the create request in a background job. Recommended: 'true'. Default to 'true'.

func (r *route) Create(req RouteRequest, opts ...bool) (*RouteFields, error) {
	async := true
	if len(opts) > 0 {
		async = opts[0]
	}
	rawURL := fmt.Sprintf("/v2/routes?async=%t&inline-relations-depth=1", async)
	routeFields := RouteFields{}
	_, err := r.client.Post(rawURL, req, &routeFields)
	if err != nil {
		return nil, err
	}
	return &routeFields, nil
}

// opts is list of boolean parametes
// opts[0] - async - Will run the update request in a background job. Recommended: 'true'. Default to 'true'.

func (r *route) Update(routeGUID string, req RouteUpdateRequest, opts ...bool) (*RouteFields, error) {
	async := true
	if len(opts) > 0 {
		async = opts[0]
	}
	rawURL := fmt.Sprintf("/v2/routes/%s?async=%t", routeGUID, async)
	routeFields := RouteFields{}
	_, err := r.client.Put(rawURL, req, &routeFields)
	if err != nil {
		return nil, err
	}
	return &routeFields, nil
}

// opts is list of boolean parametes
// opts[0] - async - Will run the delete request in a background job. Recommended: 'true'. Default to 'true'.
// opts[1] - recursive - Will delete route service bindings and route mappings associated with the route. Default to 'false'.

func (r *route) Delete(routeGUID string, opts ...bool) error {
	async := true
	recursive := false
	if len(opts) > 0 {
		async = opts[0]
	}
	if len(opts) > 1 {
		recursive = opts[1]
	}
	rawURL := fmt.Sprintf("/v2/routes/%s?async=%t&recursive=%t", routeGUID, async, recursive)
	_, err := r.client.Delete(rawURL)
	return err
}

func listRouteWithPath(c *client.Client, path string) ([]Route, error) {
	var route []Route
	_, err := c.GetPaginated(path, NewCCPaginatedResources(RouteResource{}), func(resource interface{}) bool {
		if routeResource, ok := resource.(RouteResource); ok {
			route = append(route, routeResource.ToFields())
			return true
		}
		return false
	})
	return route, err
}
