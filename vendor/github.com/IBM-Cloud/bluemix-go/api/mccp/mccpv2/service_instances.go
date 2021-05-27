package mccpv2

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//ServiceInstanceCreateRequest ...
type ServiceInstanceCreateRequest struct {
	Name      string                 `json:"name"`
	SpaceGUID string                 `json:"space_guid"`
	PlanGUID  string                 `json:"service_plan_guid"`
	Params    map[string]interface{} `json:"parameters,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
}

//ServiceInstanceUpdateRequest ...
type ServiceInstanceUpdateRequest struct {
	Name     *string                `json:"name,omitempty"`
	PlanGUID *string                `json:"service_plan_guid,omitempty"`
	Params   map[string]interface{} `json:"parameters,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
}

//ServiceInstance ...
type ServiceInstance struct {
	GUID              string
	Name              string                 `json:"name"`
	Credentials       map[string]interface{} `json:"credentials"`
	ServicePlanGUID   string                 `json:"service_plan_guid"`
	SpaceGUID         string                 `json:"space_guid"`
	GatewayData       string                 `json:"gateway_data"`
	Type              string                 `json:"type"`
	DashboardURL      string                 `json:"dashboard_url"`
	LastOperation     LastOperationFields    `json:"last_operation"`
	RouteServiceURL   string                 `json:"routes_url"`
	Tags              []string               `json:"tags"`
	SpaceURL          string                 `json:"space_url"`
	ServicePlanURL    string                 `json:"service_plan_url"`
	ServiceBindingURL string                 `json:"service_bindings_url"`
	ServiceKeysURL    string                 `json:"service_keys_url"`
	ServiceKeys       []ServiceKeyFields     `json:"service_keys"`
	ServicePlan       ServicePlanFields      `json:"service_plan"`
}

//ServiceInstanceFields ...
type ServiceInstanceFields struct {
	Metadata ServiceInstanceMetadata
	Entity   ServiceInstance
}

//ServiceInstanceMetadata ...
type ServiceInstanceMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//LastOperationFields ...
type LastOperationFields struct {
	Type        string `json:"type"`
	State       string `json:"state"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

//ServiceInstanceResource ...
type ServiceInstanceResource struct {
	Resource
	Entity ServiceInstanceEntity
}

//ServiceInstanceEntity ...
type ServiceInstanceEntity struct {
	Name              string                 `json:"name"`
	Credentials       map[string]interface{} `json:"credentials"`
	ServicePlanGUID   string                 `json:"service_plan_guid"`
	SpaceGUID         string                 `json:"space_guid"`
	GatewayData       string                 `json:"gateway_data"`
	Type              string                 `json:"type"`
	DashboardURL      string                 `json:"dashboard_url"`
	LastOperation     LastOperationFields    `json:"last_operation"`
	RouteServiceURL   string                 `json:"routes_url"`
	Tags              []string               `json:"tags"`
	SpaceURL          string                 `json:"space_url"`
	ServicePlanURL    string                 `json:"service_plan_url"`
	ServiceBindingURL string                 `json:"service_bindings_url"`
	ServiceKeysURL    string                 `json:"service_keys_url"`
	ServiceKeys       []ServiceKeyFields     `json:"service_keys"`
	ServicePlan       ServicePlanFields      `json:"service_plan"`
}

//ToModel ...
func (resource ServiceInstanceResource) ToModel() ServiceInstance {

	entity := resource.Entity

	return ServiceInstance{
		GUID:              resource.Metadata.GUID,
		Name:              entity.Name,
		Credentials:       entity.Credentials,
		ServicePlanGUID:   entity.ServicePlanGUID,
		SpaceGUID:         entity.SpaceGUID,
		GatewayData:       entity.GatewayData,
		Type:              entity.Type,
		LastOperation:     entity.LastOperation,
		RouteServiceURL:   entity.RouteServiceURL,
		DashboardURL:      entity.DashboardURL,
		Tags:              entity.Tags,
		SpaceURL:          entity.SpaceURL,
		ServicePlanURL:    entity.ServicePlanURL,
		ServiceBindingURL: entity.ServiceBindingURL,
		ServiceKeysURL:    entity.ServiceKeysURL,
	}
}

//ServiceInstances ...
type ServiceInstances interface {
	Create(req ServiceInstanceCreateRequest) (*ServiceInstanceFields, error)
	Update(instanceGUID string, req ServiceInstanceUpdateRequest) (*ServiceInstanceFields, error)
	Delete(instanceGUID string, opts ...bool) error
	FindByName(instanceName string) (*ServiceInstance, error)
	FindByNameInSpace(spaceGUID string, instanceName string) (*ServiceInstance, error)
	Get(instanceGUID string, depth ...int) (*ServiceInstanceFields, error)
	ListServiceBindings(instanceGUID string) ([]ServiceBinding, error)
}

type serviceInstance struct {
	client *client.Client
}

func newServiceInstanceAPI(c *client.Client) ServiceInstances {
	return &serviceInstance{
		client: c,
	}
}

func (s *serviceInstance) Create(req ServiceInstanceCreateRequest) (*ServiceInstanceFields, error) {
	rawURL := "/v2/service_instances?accepts_incomplete=true"
	serviceFields := ServiceInstanceFields{}
	_, err := s.client.Post(rawURL, req, &serviceFields)
	if err != nil {
		return nil, err
	}
	return &serviceFields, nil
}

func (s *serviceInstance) Get(instanceGUID string, depth ...int) (*ServiceInstanceFields, error) {
	rawURL := fmt.Sprintf("/v2/service_instances/%s", instanceGUID)
	req := rest.GetRequest(rawURL)
	if len(depth) > 0 {
		req.Query("inline-relations-depth", strconv.Itoa(depth[0]))
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()

	serviceFields := ServiceInstanceFields{}
	_, err = s.client.Get(path, &serviceFields)
	if err != nil {
		return nil, err
	}
	return &serviceFields, err
}

func (s *serviceInstance) FindByName(instanceName string) (*ServiceInstance, error) {
	req := rest.GetRequest("/v2/service_instances")
	req.Query("return_user_provided_service_instances", "true")
	if instanceName != "" {
		req.Query("q", "name:"+instanceName)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	services, err := listServicesWithPath(s.client, path)
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("Service instance:  %q doesn't exist", instanceName)
	}
	return &services[0], nil
}

func (s *serviceInstance) FindByNameInSpace(spaceGUID string, instanceName string) (*ServiceInstance, error) {
	req := rest.GetRequest(fmt.Sprintf("/v2/spaces/%s/service_instances", spaceGUID))
	req.Query("return_user_provided_service_instances", "true")
	if instanceName != "" {
		req.Query("q", "name:"+instanceName)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	services, err := listServicesWithPath(s.client, path)
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("Service instance:  %q doesn't exist in the space %s", instanceName, spaceGUID)
	}
	return &services[0], nil
}

// opts is list of boolean parametes
// opts[0] - async - Will run the delete request in a background job. Recommended: 'true'. Default to 'true'.
// opts[1] - recursive - Will delete service bindings, service keys, and routes associated with the service instance. Default to 'false'.

func (s *serviceInstance) Delete(instanceGUID string, opts ...bool) error {
	async := true
	recursive := false
	if len(opts) > 0 {
		async = opts[0]
	}
	if len(opts) > 1 {
		recursive = opts[1]
	}
	rawURL := fmt.Sprintf("/v2/service_instances/%s?accepts_incomplete=true&async=%t&recursive=%t", instanceGUID, async, recursive)
	_, err := s.client.Delete(rawURL)
	return err
}

func (s *serviceInstance) Update(instanceGUID string, req ServiceInstanceUpdateRequest) (*ServiceInstanceFields, error) {
	rawURL := fmt.Sprintf("/v2/service_instances/%s?accepts_incomplete=true", instanceGUID)
	serviceFields := ServiceInstanceFields{}
	_, err := s.client.Put(rawURL, req, &serviceFields)
	if err != nil {
		return nil, err
	}
	return &serviceFields, nil
}

func (s *serviceInstance) ListServiceBindings(instanceGUID string) ([]ServiceBinding, error) {
	rawURL := fmt.Sprintf("/v2/service_instances/%s/service_bindings", instanceGUID)
	req := rest.GetRequest(rawURL)
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	sb, err := listServiceBindingWithPath(s.client, path)
	if err != nil {
		return nil, err
	}
	return sb, nil
}

func listServicesWithPath(client *client.Client, path string) ([]ServiceInstance, error) {
	var services []ServiceInstance
	_, err := client.GetPaginated(path, NewCCPaginatedResources(ServiceInstanceResource{}), func(resource interface{}) bool {
		if serviceInstanceResource, ok := resource.(ServiceInstanceResource); ok {
			services = append(services, serviceInstanceResource.ToModel())
			return true
		}
		return false
	})
	return services, err
}
