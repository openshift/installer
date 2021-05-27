package controller

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

type CreateServiceInstanceRequest struct {
	Name            string                 `json:"name"`
	ServicePlanID   string                 `json:"resource_plan_id"`
	ResourceGroupID string                 `json:"resource_group_id"`
	Crn             string                 `json:"crn,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
	Parameters      map[string]interface{} `json:"parameters,omitempty"`
	TargetCrn       string                 `json:"target_crn"`
}

type UpdateServiceInstanceRequest struct {
	Name          string                 `json:"name,omitempty"`
	ServicePlanID string                 `json:"resource_plan_id,omitempty"`
	Tags          []string               `json:"tags,omitempty"`
	Parameters    map[string]interface{} `json:"parameters,omitempty"`
	UpdateTime    int64                  `json:"update_time,omitempty"`
}

type ServiceInstanceQuery struct {
	ResourceGroupID string
	ServiceID       string
	ServicePlanID   string
	Name            string
}

//ErrCodeResourceServiceInstanceDoesnotExist ...
const ErrCodeResourceServiceInstanceDoesnotExist = "ResourceServiceInstanceDoesnotExist"

//ResourceServiceInstanceQuery ...
type ResourceServiceInstanceRepository interface {
	ListInstances(query ServiceInstanceQuery) ([]models.ServiceInstance, error)
	GetInstance(serviceInstanceID string) (models.ServiceInstance, error)
	CreateInstance(serviceInstanceRequest CreateServiceInstanceRequest) (models.ServiceInstance, error)
	UpdateInstance(serviceInstanceID string, updateInstanceRequest UpdateServiceInstanceRequest) (models.ServiceInstance, error)
	DeleteInstance(serviceInstanceID string, recursive bool) error
	//GetBindings(serviceInstanceID string) ([]ServiceBinding, error)
}

type resourceServiceInstance struct {
	client *client.Client
}

func newResourceServiceInstanceAPI(c *client.Client) ResourceServiceInstanceRepository {
	return &resourceServiceInstance{
		client: c,
	}
}

func (r *resourceServiceInstance) ListInstances(query ServiceInstanceQuery) ([]models.ServiceInstance, error) {
	listRequest := rest.GetRequest("/v1/resource_instances").
		Query("resource_group_id", query.ResourceGroupID).
		Query("resource_id", query.ServiceID).
		Query("resource_plan_id", query.ServicePlanID)

	req, err := listRequest.Build()
	if err != nil {
		return nil, err
	}

	var instances []models.ServiceInstance
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewRCPaginatedResources(models.ServiceInstance{}),
		func(resource interface{}) bool {
			if instance, ok := resource.(models.ServiceInstance); ok {
				instances = append(instances, instance)
				return true
			}
			return false
		},
	)
	if err != nil {
		return []models.ServiceInstance{}, err
	}

	if query.Name != "" {
		instances = filterInstancesByName(instances, query.Name)
	}
	return instances, nil
}

func (r *resourceServiceInstance) CreateInstance(serviceInstanceRequest CreateServiceInstanceRequest) (models.ServiceInstance, error) {
	resp := models.ServiceInstance{}
	request := rest.PostRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/resource_instances"))
	_, err := r.client.SendRequest(request.Body(serviceInstanceRequest), &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
func (r *resourceServiceInstance) UpdateInstance(serviceInstanceID string, updateInstanceRequest UpdateServiceInstanceRequest) (models.ServiceInstance, error) {
	resp := models.ServiceInstance{}
	request := rest.PatchRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/resource_instances/"+url.PathEscape(serviceInstanceID)))
	_, err := r.client.SendRequest(request.Body(updateInstanceRequest), &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (r *resourceServiceInstance) GetInstance(serviceInstanceID string) (models.ServiceInstance, error) {
	var instance models.ServiceInstance
	resp, err := r.client.Get("/v1/resource_instances/"+url.PathEscape(serviceInstanceID), &instance)
	if resp.StatusCode == http.StatusNotFound {
		return models.ServiceInstance{}, bmxerror.New(ErrCodeResourceServiceInstanceDoesnotExist,
			fmt.Sprintf("Given service instance : %q doesn't exist", serviceInstanceID))
	}
	return instance, err
}

func (r *resourceServiceInstance) DeleteInstance(resourceInstanceID string, recursive bool) error {
	request := rest.DeleteRequest(helpers.GetFullURL(*r.client.Config.Endpoint, "/v1/resource_instances/"+url.PathEscape(resourceInstanceID)))
	if recursive {
		request = request.Query("recursive", "true")
	}
	_, err := r.client.SendRequest(request, nil)
	return err
}

/*func (r *resourceServiceInstance) GetBindings(serviceInstanceID string) ([]ServiceBinding, error) {
	listRequest := rest.GetRequest("/v1/resource_instances/" + url.PathEscape(serviceInstanceID) + "/resource_bindings")
	req, err := listRequest.Build()
	if err != nil {
		return nil, err
	}

	var bindings []ServiceBinding
	_, err = r.client.GetPaginated(
		req.URL.String(),
		ServiceBinding{},
		func(resource interface{}) bool {
			if binding, ok := resource.(ServiceBinding); ok {
				bindings = append(bindings, binding)
				return true
			}
			return false
		})
	if err != nil {
		return []ServiceBinding{}, err
	}
	return bindings, nil
}*/

func filterInstancesByName(instances []models.ServiceInstance, name string) []models.ServiceInstance {
	ret := []models.ServiceInstance{}
	for _, instance := range instances {
		if instance.Name == name {
			ret = append(ret, instance)
		}
	}
	return ret
}
