package controllerv2

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/client"
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
	Type            string
	SubType         string
	Limit           string
	UpdatedFrom     string
	UpdatedTo       string
	Guid            string
}

//ErrCodeResourceServiceInstanceDoesnotExist ...
const ErrCodeResourceServiceInstanceDoesnotExist = "ResourceServiceInstanceDoesnotExist"

//ResourceServiceInstanceQuery ...
type ResourceServiceInstanceRepository interface {
	ListInstances(query ServiceInstanceQuery) ([]models.ServiceInstanceV2, error)
	GetInstance(serviceInstanceID string) (models.ServiceInstanceV2, error)
}

type resourceServiceInstance struct {
	client *client.Client
}

func newResourceServiceInstanceAPI(c *client.Client) ResourceServiceInstanceRepository {
	return &resourceServiceInstance{
		client: c,
	}
}

func (r *resourceServiceInstance) ListInstances(query ServiceInstanceQuery) ([]models.ServiceInstanceV2, error) {
	listRequest := rest.GetRequest("/v2/resource_instances").
		Query("resource_group_id", query.ResourceGroupID).
		Query("resource_id", query.ServiceID).
		Query("resource_plan_id", query.ServicePlanID).
		Query("type", query.Type).
		Query("sub_type", query.SubType).
		Query("limit", query.Limit).
		Query("updated_from", query.UpdatedFrom).
		Query("updated_to", query.UpdatedTo).
		Query("guid", query.Guid)

	req, err := listRequest.Build()
	if err != nil {
		return nil, err
	}

	var instances []models.ServiceInstanceV2
	_, err = r.client.GetPaginated(
		req.URL.String(),
		NewRCPaginatedResources(models.ServiceInstanceV2{}),
		func(resource interface{}) bool {
			if instance, ok := resource.(models.ServiceInstanceV2); ok {
				instances = append(instances, instance)
				return true
			}
			return false
		},
	)
	if err != nil {
		return []models.ServiceInstanceV2{}, err
	}

	if query.Name != "" {
		instances = filterInstancesByName(instances, query.Name)
	}
	return instances, nil
}

func (r *resourceServiceInstance) GetInstance(serviceInstanceID string) (models.ServiceInstanceV2, error) {
	var instance models.ServiceInstanceV2
	resp, err := r.client.Get("/v2/resource_instances/"+url.PathEscape(serviceInstanceID), &instance)
	if resp.StatusCode == http.StatusNotFound {
		return models.ServiceInstanceV2{}, bmxerror.New(ErrCodeResourceServiceInstanceDoesnotExist,
			fmt.Sprintf("Given service instance : %q doesn't exist", serviceInstanceID))
	}
	return instance, err
}

func filterInstancesByName(instances []models.ServiceInstanceV2, name string) []models.ServiceInstanceV2 {
	ret := []models.ServiceInstanceV2{}
	for _, instance := range instances {
		if instance.Name == name {
			ret = append(ret, instance)
		}
	}
	return ret
}
