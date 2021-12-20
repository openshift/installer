package controller

import (
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type CreateServiceBindingRequest struct {
	SourceCRN  crn.CRN                `json:"source_crn"`
	TargetCRN  crn.CRN                `json:"target_crn"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type ResourceServiceBindingRepository interface {
	ListBindings(cb func(models.ServiceBinding) bool) error
	GetBinding(bindingID string) (models.ServiceBinding, error)
	CreateBinding(CreateServiceBindingRequest) (models.ServiceBinding, error)
	DeleteBinding(bindingID string) error
}

type serviceBindingRepository struct {
	client *client.Client
}

func newServiceBindingRepository(c *client.Client) ResourceServiceBindingRepository {
	return &serviceBindingRepository{
		client: c,
	}
}

func (r *serviceBindingRepository) ListBindings(cb func(models.ServiceBinding) bool) error {
	listRequest := rest.GetRequest("/v1/resource_bindings")
	req, err := listRequest.Build()
	if err != nil {
		return err
	}

	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewRCPaginatedResources(models.ServiceBinding{}),
		func(resource interface{}) bool {
			if binding, ok := resource.(models.ServiceBinding); ok {
				return cb(binding)
			}
			return false
		})

	return err
}

func (r *serviceBindingRepository) CreateBinding(createBindingRequest CreateServiceBindingRequest) (models.ServiceBinding, error) {
	resp := models.ServiceBinding{}
	_, err := r.client.Post("/v1/resource_bindings", createBindingRequest, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (r *serviceBindingRepository) GetBinding(bindingID string) (models.ServiceBinding, error) {
	resp := models.ServiceBinding{}
	request := rest.GetRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/resource_bindings/"+url.PathEscape(bindingID)))
	_, err := r.client.SendRequest(request, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil

}

func (r *serviceBindingRepository) DeleteBinding(bindingID string) error {
	request := rest.DeleteRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/resource_bindings/"+url.PathEscape(bindingID))).Query("id", url.PathEscape(bindingID))
	_, err := r.client.SendRequest(request, nil)
	return err
}
