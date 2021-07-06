package mccpv2

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//ServiceBindingRequest ...
type ServiceBindingRequest struct {
	ServiceInstanceGUID string `json:"service_instance_guid"`
	AppGUID             string `json:"app_guid"`
	Parameters          string `json:"parameters,omitempty"`
}

//ServiceBindingMetadata ...
type ServiceBindingMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//ServiceBindingEntity ...
type ServiceBindingEntity struct {
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	AppGUID             string                 `json:"app_guid"`
	Credentials         map[string]interface{} `json:"credentials"`
}

//ServiceBindingResource ...
type ServiceBindingResource struct {
	Resource
	Entity ServiceBindingEntity
}

//ServiceBindingFields ...
type ServiceBindingFields struct {
	Metadata ServiceBindingMetadata
	Entity   ServiceBindingEntity
}

//ServiceBinding model
type ServiceBinding struct {
	GUID                string
	ServiceInstanceGUID string
	AppGUID             string
	Credentials         map[string]interface{}
}

//ToFields ..
func (resource ServiceBindingResource) ToFields() ServiceBinding {
	entity := resource.Entity

	return ServiceBinding{
		GUID:                resource.Metadata.GUID,
		ServiceInstanceGUID: entity.ServiceInstanceGUID,
		AppGUID:             entity.AppGUID,
		Credentials:         entity.Credentials,
	}
}

//ServiceBindings ...
type ServiceBindings interface {
	Create(req ServiceBindingRequest) (*ServiceBindingFields, error)
	Get(guid string) (*ServiceBindingFields, error)
	Delete(guid string, opts ...bool) error
	List(filters ...string) ([]ServiceBinding, error)
}

type serviceBinding struct {
	client *client.Client
}

func newServiceBindingAPI(c *client.Client) ServiceBindings {
	return &serviceBinding{
		client: c,
	}
}

func (r *serviceBinding) Get(sbGUID string) (*ServiceBindingFields, error) {
	rawURL := fmt.Sprintf("/v2/service_bindings/%s", sbGUID)
	sbFields := ServiceBindingFields{}
	_, err := r.client.Get(rawURL, &sbFields, nil)
	if err != nil {
		return nil, err
	}
	return &sbFields, nil
}

func (r *serviceBinding) Create(req ServiceBindingRequest) (*ServiceBindingFields, error) {
	rawURL := "/v2/service_bindings"
	sbFields := ServiceBindingFields{}
	_, err := r.client.Post(rawURL, req, &sbFields)
	if err != nil {
		return nil, err
	}
	return &sbFields, nil
}

// opts is list of boolean parametes
// opts[0] - async - Will run the delete request in a background job. Recommended: 'true'. Default to 'true'.

func (r *serviceBinding) Delete(guid string, opts ...bool) error {
	async := true
	if len(opts) > 0 {
		async = opts[0]
	}
	rawURL := fmt.Sprintf("/v2/service_bindings/%s?async=%t", guid, async)
	_, err := r.client.Delete(rawURL)
	return err
}

func (r *serviceBinding) List(filters ...string) ([]ServiceBinding, error) {
	rawURL := "/v2/service_bindings"
	req := rest.GetRequest(rawURL)
	if len(filters) > 0 {
		req.Query("q", strings.Join(filters, ""))
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	bindings, err := listServiceBindingWithPath(r.client, path)
	if err != nil {
		return nil, err
	}
	return bindings, nil
}

func listServiceBindingWithPath(c *client.Client, path string) ([]ServiceBinding, error) {
	var sb []ServiceBinding
	_, err := c.GetPaginated(path, NewCCPaginatedResources(ServiceBindingResource{}), func(resource interface{}) bool {
		if sbResource, ok := resource.(ServiceBindingResource); ok {
			sb = append(sb, sbResource.ToFields())
			return true
		}
		return false
	})
	return sb, err
}
