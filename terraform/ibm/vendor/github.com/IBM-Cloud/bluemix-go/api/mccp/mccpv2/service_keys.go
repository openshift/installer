package mccpv2

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

//ErrCodeServiceKeyDoesNotExist ...
const ErrCodeServiceKeyDoesNotExist = "erviceKeyDoesNotExist"

//ServiceKeyRequest ...
type ServiceKeyRequest struct {
	Name                string                 `json:"name"`
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	Params              map[string]interface{} `json:"parameters,omitempty"`
}

//ServiceKey  model...
type ServiceKey struct {
	GUID                string
	Name                string                 `json:"name"`
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	ServiceInstanceURL  string                 `json:"service_instance_url"`
	Credentials         map[string]interface{} `json:"credentials"`
}

//ServiceKeyFields ...
type ServiceKeyFields struct {
	Metadata ServiceKeyMetadata
	Entity   ServiceKey
}

//ServiceKeyMetadata ...
type ServiceKeyMetadata struct {
	GUID      string `json:"guid"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

//ServiceKeyResource ...
type ServiceKeyResource struct {
	Resource
	Entity ServiceKeyEntity
}

//ServiceKeyEntity ...
type ServiceKeyEntity struct {
	Name                string                 `json:"name"`
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	ServiceInstanceURL  string                 `json:"service_instance_url"`
	Credentials         map[string]interface{} `json:"credentials"`
}

//ToModel ...
func (resource ServiceKeyResource) ToModel() ServiceKey {

	entity := resource.Entity

	return ServiceKey{
		GUID:                resource.Metadata.GUID,
		Name:                entity.Name,
		ServiceInstanceGUID: entity.ServiceInstanceGUID,
		ServiceInstanceURL:  entity.ServiceInstanceURL,
		Credentials:         entity.Credentials,
	}
}

//ServiceKeys ...
type ServiceKeys interface {
	Create(serviceInstanceGUID string, keyName string, params map[string]interface{}) (*ServiceKeyFields, error)
	FindByName(serviceInstanceGUID string, keyName string) (*ServiceKey, error)
	Get(serviceKeyGUID string) (*ServiceKeyFields, error)
	Delete(serviceKeyGUID string) error
	List(filters ...string) ([]ServiceKey, error)
}

type serviceKey struct {
	client *client.Client
}

func newServiceKeyAPI(c *client.Client) ServiceKeys {
	return &serviceKey{
		client: c,
	}
}

func (r *serviceKey) Create(serviceInstanceGUID string, keyName string, params map[string]interface{}) (*ServiceKeyFields, error) {
	serviceKeyFields := ServiceKeyFields{}
	reqParam := ServiceKeyRequest{
		ServiceInstanceGUID: serviceInstanceGUID,
		Name:                keyName,
		Params:              params,
	}
	_, err := r.client.Post("/v2/service_keys", reqParam, &serviceKeyFields)
	if err != nil {
		return nil, err
	}
	return &serviceKeyFields, nil
}

func (r *serviceKey) Delete(serviceKeyGUID string) error {
	rawURL := fmt.Sprintf("/v2/service_keys/%s", serviceKeyGUID)
	_, err := r.client.Delete(rawURL)
	return err
}

func (r *serviceKey) Get(guid string) (*ServiceKeyFields, error) {
	rawURL := fmt.Sprintf("/v2/service_keys/%s", guid)
	serviceKeyFields := ServiceKeyFields{}
	_, err := r.client.Get(rawURL, &serviceKeyFields)
	if err != nil {
		return nil, err
	}

	return &serviceKeyFields, err
}

func (r *serviceKey) FindByName(serviceInstanceGUID string, keyName string) (*ServiceKey, error) {
	rawURL := fmt.Sprintf("/v2/service_instances/%s/service_keys", serviceInstanceGUID)
	req := rest.GetRequest(rawURL)
	if keyName != "" {
		req.Query("q", "name:"+keyName)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	serviceKeys, err := r.listServiceKeysWithPath(path)
	if err != nil {
		return nil, err
	}
	if len(serviceKeys) == 0 {
		return nil, bmxerror.New(ErrCodeServiceKeyDoesNotExist,
			fmt.Sprintf("Given service key %q doesn't exist for the given service instance  %q", keyName, serviceInstanceGUID))
	}
	return &serviceKeys[0], nil
}

func (r *serviceKey) List(filters ...string) ([]ServiceKey, error) {
	rawURL := "/v2/service_keys"
	req := rest.GetRequest(rawURL)
	if len(filters) > 0 {
		req.Query("q", strings.Join(filters, ""))
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	keys, err := r.listServiceKeysWithPath(path)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *serviceKey) listServiceKeysWithPath(path string) ([]ServiceKey, error) {
	var serviceKeys []ServiceKey
	_, err := r.client.GetPaginated(path, NewCCPaginatedResources(ServiceKeyResource{}), func(resource interface{}) bool {
		if serviceKeyResource, ok := resource.(ServiceKeyResource); ok {
			serviceKeys = append(serviceKeys, serviceKeyResource.ToModel())
			return true
		}
		return false
	})
	return serviceKeys, err
}
